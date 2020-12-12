package gocbcore

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/couchbase/gocbcore/v9/memd"
)

type memdSenderClient interface {
	SupportsFeature(memd.HelloFeature) bool
	Address() string
	SendRequest(*memdQRequest) error
}

type syncClient struct {
	client memdSenderClient
}

func (client *syncClient) SupportsFeature(feature memd.HelloFeature) bool {
	return client.client.SupportsFeature(feature)
}

func (client *syncClient) Address() string {
	return client.client.Address()
}

func (client *syncClient) doRequest(req *memd.Packet, deadline time.Time) (respOut *memd.Packet, errOut error) {
	signal := make(chan bool, 1)

	qreq := memdQRequest{
		Packet: *req,
		Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
			if resp != nil {
				respOut = &resp.Packet
			}
			errOut = err
			signal <- true
		},
		RetryStrategy: newFailFastRetryStrategy(),
	}

	err := client.client.SendRequest(&qreq)
	if err != nil {
		return nil, err
	}

	timeoutTmr := AcquireTimer(time.Until(deadline))
	select {
	case <-signal:
		ReleaseTimer(timeoutTmr, false)
		return
	case <-timeoutTmr.C:
		ReleaseTimer(timeoutTmr, true)
		qreq.cancelWithCallback(errAmbiguousTimeout)
		<-signal
		return
	}
}

func (client *syncClient) doBasicOp(cmd memd.CmdCode, k, v, e []byte, deadline time.Time) ([]byte, error) {
	resp, err := client.doRequest(
		&memd.Packet{
			Magic:   memd.CmdMagicReq,
			Command: cmd,
			Key:     k,
			Value:   v,
			Extras:  e,
		},
		deadline,
	)

	// We do it this way as the response value could still be useful even if an
	// error status code is returned.  For instance, StatusAuthContinue still
	// contains authentication stepping information.
	if resp == nil {
		return nil, err
	}

	return resp.Value, err
}

func (client *syncClient) ExecDcpControl(key string, value string, deadline time.Time) error {
	_, err := client.doBasicOp(memd.CmdDcpControl, []byte(key), []byte(value), nil, deadline)
	return err
}

func (client *syncClient) ExecGetClusterConfig(deadline time.Time) ([]byte, error) {
	return client.doBasicOp(memd.CmdGetClusterConfig, nil, nil, nil, deadline)
}

func (client *syncClient) ExecOpenDcpConsumer(streamName string, openFlags memd.DcpOpenFlag, deadline time.Time) error {
	_, ok := client.client.(*memdClient)
	if !ok {
		return errCliInternalError
	}

	extraBuf := make([]byte, 8)
	binary.BigEndian.PutUint32(extraBuf[0:], 0)
	binary.BigEndian.PutUint32(extraBuf[4:], uint32((openFlags & ^memd.DcpOpenFlag(3))|memd.DcpOpenFlagProducer))
	_, err := client.doBasicOp(memd.CmdDcpOpenConnection, []byte(streamName), nil, extraBuf, deadline)
	return err
}

func (client *syncClient) ExecEnableDcpNoop(period time.Duration, deadline time.Time) error {
	_, ok := client.client.(*memdClient)
	if !ok {
		return errCliInternalError
	}
	// The client will always reply to No-Op's.  No need to enable it

	err := client.ExecDcpControl("enable_noop", "true", deadline)
	if err != nil {
		return err
	}

	periodStr := fmt.Sprintf("%d", period/time.Second)
	err = client.ExecDcpControl("set_noop_interval", periodStr, deadline)
	if err != nil {
		return err
	}

	return nil
}

func (client *syncClient) ExecEnableDcpClientEnd(deadline time.Time) error {
	memcli, ok := client.client.(*memdClient)
	if !ok {
		return errCliInternalError
	}

	err := client.ExecDcpControl("send_stream_end_on_client_close_stream", "true", deadline)
	if err != nil {
		memcli.streamEndNotSupported = true
	}

	return nil
}

func (client *syncClient) ExecEnableDcpBufferAck(bufferSize int, deadline time.Time) error {
	mclient, ok := client.client.(*memdClient)
	if !ok {
		return errCliInternalError
	}

	// Enable buffer acknowledgment on the client
	mclient.EnableDcpBufferAck(bufferSize / 2)

	bufferSizeStr := fmt.Sprintf("%d", bufferSize)
	err := client.ExecDcpControl("connection_buffer_size", bufferSizeStr, deadline)
	if err != nil {
		return err
	}

	return nil
}
