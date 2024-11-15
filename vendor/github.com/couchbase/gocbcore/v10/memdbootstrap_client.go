package gocbcore

import (
	"encoding/binary"
	"errors"
	"strings"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

type bootstrapableClient interface {
	SendRequest(*memdQRequest) error
	Address() string
	ConnID() string
	SupportsFeature(feature memd.HelloFeature) bool
	Features([]memd.HelloFeature)
	loggerID() string
}

type bootstrapClient interface {
	Address() string
	ConnID() string
	Features(features []memd.HelloFeature)
	SupportsFeature(feature memd.HelloFeature) bool
	SaslAuth(k, v []byte, deadline time.Time, cb func(b []byte, err error)) error
	SaslStep(k, v []byte, deadline time.Time, cb func(err error)) error
	ExecSelectBucket(b []byte, deadline time.Time) (chan error, error)
	ExecGetErrorMap(version uint16, deadline time.Time) (chan errorMapResponse, error)
	SaslListMechs(deadline time.Time, cb func(mechs []AuthMechanism, err error)) error
	ExecHello(clientID string, features []memd.HelloFeature, deadline time.Time) (chan ExecHelloResponse, error)
	ExecGetConfig(deadline time.Time) (chan getConfigResponse, error)
	LoggerID() string
}

// Due to AuthProvider we are currently tied to bootstrapping passing around a deadline and the bootstrap
// "owner" has to hold onto a cancel sig for use at request time.
// In the future we can combine deadline and cancellation into a context.Context and pass that everywhere as a parameter,
// we will then able to expose utility functions to allow user to build their own bootstrap from existing building
// blocks.
func newMemdBootstrapClient(client bootstrapableClient, cancelSig <-chan struct{}) *memdBootstrapClient {
	return &memdBootstrapClient{
		cancelSig: cancelSig,
		client:    client,
	}
}

type memdBootstrapClient struct {
	client    bootstrapableClient
	cancelSig <-chan struct{}
}

func (bc *memdBootstrapClient) Address() string {
	return bc.client.Address()
}

func (bc *memdBootstrapClient) ConnID() string {
	return bc.client.ConnID()
}

func (bc *memdBootstrapClient) Features(features []memd.HelloFeature) {
	bc.client.Features(features)
}

func (bc *memdBootstrapClient) SupportsFeature(feature memd.HelloFeature) bool {
	return bc.client.SupportsFeature(feature)
}

func (client *memdBootstrapClient) LoggerID() string {
	return client.client.loggerID()
}

func (bc *memdBootstrapClient) SaslAuth(k, v []byte, deadline time.Time, cb func(b []byte, err error)) error {
	err := bc.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdSASLAuth,
				Key:     k,
				Value:   v,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				// Auth is special, auth continue is surfaced as an error
				var val []byte
				if resp != nil {
					val = resp.Value
				}

				cb(val, err)
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return err
	}

	return nil
}

func (bc *memdBootstrapClient) SaslStep(k, v []byte, deadline time.Time, cb func(err error)) error {
	err := bc.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdSASLStep,
				Key:     k,
				Value:   v,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				if err != nil {
					cb(err)
					return
				}

				cb(nil)
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return err
	}

	return nil
}

func (bc *memdBootstrapClient) ExecSelectBucket(b []byte, deadline time.Time) (chan error, error) {
	completedCh := make(chan error, 1)
	err := bc.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdSelectBucket,
				Key:     b,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				if err != nil {
					if errors.Is(err, ErrDocumentNotFound) {
						// Bucket not found means that the user has privileges to access the bucket but that the bucket
						// is in some way not existing right now (e.g. in warmup).
						err = errBucketNotFound
					}
					completedCh <- err
					return
				}

				completedCh <- nil
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return nil, err
	}

	return completedCh, nil
}

type errorMapResponse struct {
	Err   error
	Bytes []byte
}

func (bc *memdBootstrapClient) ExecGetErrorMap(version uint16, deadline time.Time) (chan errorMapResponse, error) {
	completedCh := make(chan errorMapResponse, 1)
	valueBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(valueBuf, version)

	err := bc.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdGetErrorMap,
				Value:   valueBuf,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				if err != nil {
					completedCh <- errorMapResponse{
						Err: err,
					}
					return
				}

				completedCh <- errorMapResponse{
					Bytes: resp.Value,
				}
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return nil, err
	}

	return completedCh, nil
}

func (bc *memdBootstrapClient) SaslListMechs(deadline time.Time, cb func(mechs []AuthMechanism, err error)) error {
	err := bc.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdSASLListMechs,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				if err != nil {
					cb(nil, err)
					return
				}

				mechs := strings.Split(string(resp.Value), " ")
				var authMechs []AuthMechanism
				for _, mech := range mechs {
					authMechs = append(authMechs, AuthMechanism(mech))
				}

				cb(authMechs, nil)
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return err
	}

	return nil
}

// ExecHelloResponse contains the features and/or error from an ExecHello operation.
type ExecHelloResponse struct {
	SrvFeatures []memd.HelloFeature
	Err         error
}

func (bc *memdBootstrapClient) ExecHello(clientID string, features []memd.HelloFeature, deadline time.Time) (chan ExecHelloResponse, error) {
	appendFeatureCode := func(bytes []byte, feature memd.HelloFeature) []byte {
		bytes = append(bytes, 0, 0)
		binary.BigEndian.PutUint16(bytes[len(bytes)-2:], uint16(feature))
		return bytes
	}

	var featureBytes []byte
	for _, feature := range features {
		featureBytes = appendFeatureCode(featureBytes, feature)
	}

	completedCh := make(chan ExecHelloResponse, 1)
	err := bc.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdHello,
				Key:     []byte(clientID),
				Value:   featureBytes,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				if err != nil {
					completedCh <- ExecHelloResponse{
						Err: err,
					}
					return
				}

				var srvFeatures []memd.HelloFeature
				for i := 0; i < len(resp.Value); i += 2 {
					feature := binary.BigEndian.Uint16(resp.Value[i:])
					srvFeatures = append(srvFeatures, memd.HelloFeature(feature))
				}

				completedCh <- ExecHelloResponse{
					SrvFeatures: srvFeatures,
				}
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return nil, err
	}

	return completedCh, nil
}

type getConfigResponse struct {
	Err    error
	Config *cfgBucket
}

func (bc *memdBootstrapClient) ExecGetConfig(deadline time.Time) (chan getConfigResponse, error) {
	completedCh := make(chan getConfigResponse, 1)
	// Note that the revid/revepoch do not matter here, we only send GetConfig on bootstrap if we haven't actually
	// seen a config yet.
	err := bc.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdGetClusterConfig,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				if err != nil {
					completedCh <- getConfigResponse{
						Err: err,
					}
					return
				}

				hostName, err := hostFromHostPort(bc.Address())
				if err != nil {
					logWarnf("Boostrap client: Failed to parse source address. %s", err)
					completedCh <- getConfigResponse{
						Err: err,
					}
					return
				}

				bk, err := parseConfig(resp.Value, hostName)
				if err != nil {
					logWarnf("Boostrap client: Failed to parse CCCP config. %v", err)
					completedCh <- getConfigResponse{
						Err: err,
					}
					return
				}

				completedCh <- getConfigResponse{
					Config: bk,
				}
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return nil, err
	}

	return completedCh, nil
}

func (bc *memdBootstrapClient) doBootstrapRequest(req *memdQRequest, deadline time.Time) error {
	origCb := req.Callback
	doneCh := make(chan struct{})
	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		close(doneCh)
		origCb(resp, req, err)
	}

	req.Callback = handler

	err := bc.client.SendRequest(req)
	if err != nil {
		return err
	}

	start := time.Now()
	req.SetTimer(time.AfterFunc(deadline.Sub(start), func() {
		connInfo := req.ConnectionInfo()
		count, reasons := req.Retries()
		req.cancelWithCallback(&TimeoutError{
			InnerError:         errAmbiguousTimeout,
			OperationID:        req.Command.Name(),
			Opaque:             req.Identifier(),
			TimeObserved:       time.Since(start),
			RetryReasons:       reasons,
			RetryAttempts:      count,
			LastDispatchedTo:   connInfo.lastDispatchedTo,
			LastDispatchedFrom: connInfo.lastDispatchedFrom,
			LastConnectionID:   connInfo.lastConnectionID,
		})
	}))

	go func() {
		select {
		case <-doneCh:
			return
		case <-bc.cancelSig:
			req.Cancel()
			<-doneCh
			return
		}
	}()

	return nil
}
