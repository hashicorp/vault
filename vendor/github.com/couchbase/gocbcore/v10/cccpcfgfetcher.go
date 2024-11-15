package gocbcore

import (
	"encoding/binary"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

type cccpConfigFetcher struct {
	requestTimeout time.Duration
}

func newCCCPConfigFetcher(reqTimeout time.Duration) *cccpConfigFetcher {
	return &cccpConfigFetcher{
		requestTimeout: reqTimeout,
	}
}

func (ccc *cccpConfigFetcher) GetClusterConfig(pipeline *memdPipeline, revID, revEpoch int64, cancelSig chan struct{}) ([]byte, error) {
	var extras []byte
	if revID > 0 && pipeline.SupportsFeature(memd.FeatureClusterMapKnownVersion) {
		extras = make([]byte, 16)
		binary.BigEndian.PutUint64(extras[0:], uint64(revEpoch))
		binary.BigEndian.PutUint64(extras[8:], uint64(revID))
	}
	var cfgOut []byte
	var errOut error
	signal := make(chan struct{}, 1)

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:   memd.CmdMagicReq,
			Command: memd.CmdGetClusterConfig,
			Extras:  extras,
		},
		Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
			if resp != nil {
				cfgOut = resp.Packet.Value
			}
			errOut = err
			signal <- struct{}{}
		},
		RetryStrategy: newFailFastRetryStrategy(),
	}
	err := pipeline.SendRequest(req)
	if err != nil {
		return nil, err
	}

	timeoutTmr := AcquireTimer(ccc.requestTimeout)
	select {
	case <-signal:
		ReleaseTimer(timeoutTmr, false)
		return cfgOut, errOut
	case <-timeoutTmr.C:
		ReleaseTimer(timeoutTmr, true)

		req.cancelWithCallback(errUnambiguousTimeout)
		<-signal
		return cfgOut, errOut
	case <-cancelSig:
		ReleaseTimer(timeoutTmr, false)
		req.cancelWithCallback(errRequestCanceled)
		<-signal
		return nil, errOut
	}
}
