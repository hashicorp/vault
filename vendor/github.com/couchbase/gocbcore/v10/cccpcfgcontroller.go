package gocbcore

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

type cccpConfigController struct {
	muxer              dispatcher
	cfgMgr             *configManagementComponent
	confCccpPollPeriod time.Duration
	confCccpMaxWait    time.Duration

	// Used exclusively for testing to overcome GOCBC-780. It allows a test to pause the cccp looper preventing
	// unwanted requests from being sent to the mock once it has been setup for error map testing.
	looperPauseSig chan bool

	looperStopSig chan struct{}
	looperDoneSig chan struct{}

	fetchErr error
	errLock  sync.Mutex
}

func newCCCPConfigController(props cccpPollerProperties, muxer dispatcher, cfgMgr *configManagementComponent) *cccpConfigController {
	return &cccpConfigController{
		muxer:              muxer,
		cfgMgr:             cfgMgr,
		confCccpPollPeriod: props.confCccpPollPeriod,
		confCccpMaxWait:    props.confCccpMaxWait,

		looperPauseSig: make(chan bool),
		looperStopSig:  make(chan struct{}),
		looperDoneSig:  make(chan struct{}),
	}
}

type cccpPollerProperties struct {
	confCccpPollPeriod time.Duration
	confCccpMaxWait    time.Duration
}

func (ccc *cccpConfigController) Error() error {
	ccc.errLock.Lock()
	defer ccc.errLock.Unlock()
	return ccc.fetchErr
}

func (ccc *cccpConfigController) setError(err error) {
	ccc.errLock.Lock()
	ccc.fetchErr = err
	ccc.errLock.Unlock()
}

func (ccc *cccpConfigController) Pause(paused bool) {
	ccc.looperPauseSig <- paused
}

func (ccc *cccpConfigController) Stop() {
	close(ccc.looperStopSig)
}

func (ccc *cccpConfigController) Done() chan struct{} {
	return ccc.looperDoneSig
}

// Reset must never be called concurrently with the Stop or whilst the poll loop is running.
func (ccc *cccpConfigController) Reset() {
	ccc.looperStopSig = make(chan struct{})
	ccc.looperDoneSig = make(chan struct{})
}

func (ccc *cccpConfigController) DoLoop() error {
	tickTime := ccc.confCccpPollPeriod
	paused := false

	logInfof("CCCP Looper starting.")
	nodeIdx := -1
	// The first time that we loop we want to skip any sleep so that we can try get a config and bootstrapped ASAP.
	firstLoop := true

Looper:
	for {
		if !firstLoop {
			// Wait for either the agent to be shut down, or our tick time to expire
			select {
			case <-ccc.looperStopSig:
				break Looper
			case pause := <-ccc.looperPauseSig:
				paused = pause
			case <-time.After(tickTime):
			}
		}
		firstLoop = false

		if paused {
			continue
		}

		iter, err := ccc.muxer.PipelineSnapshot()
		if err != nil {
			// If we have an error it indicates the client is shut down.
			break
		}

		numNodes := iter.NumPipelines()
		if numNodes == 0 {
			logInfof("CCCPPOLL: No nodes available to poll, return upstream")
			return errNoCCCPHosts
		}

		if nodeIdx < 0 || nodeIdx > numNodes {
			nodeIdx = rand.Intn(numNodes) // #nosec G404
		}

		var foundConfig *cfgBucket
		var foundErr error
		iter.Iterate(nodeIdx, func(pipeline *memdPipeline) bool {
			nodeIdx = (nodeIdx + 1) % numNodes
			cccpBytes, err := ccc.getClusterConfig(pipeline)
			if err != nil {
				if isPollingFallbackError(err) {
					// This error is indicative of a memcached bucket which we can't handle so return the error.
					logInfof("CCCPPOLL: CCCP not supported, returning error upstream.")
					foundErr = err
					return true
				}

				// Only log the error at warn if it's unexpected.
				// If we cancelled the request or we're shutting down the connection then it's not really unexpected.
				ccc.setError(err)
				if errors.Is(err, ErrRequestCanceled) || errors.Is(err, ErrShutdown) {
					logDebugf("CCCPPOLL: CCCP request was cancelled or connection was shutdown: %v", err)
					return true
				}

				logWarnf("CCCPPOLL: Failed to retrieve CCCP config. %s", err)
				return false
			}
			ccc.setError(nil)

			logDebugf("CCCPPOLL: Got Block: %v", string(cccpBytes))

			hostName, err := hostFromHostPort(pipeline.Address())
			if err != nil {
				logWarnf("CCCPPOLL: Failed to parse source address. %s", err)
				return false
			}

			bk, err := parseConfig(cccpBytes, hostName)
			if err != nil {
				logWarnf("CCCPPOLL: Failed to parse CCCP config. %v", err)
				return false
			}

			foundConfig = bk
			return true
		})
		if foundErr != nil {
			return foundErr
		}

		if foundConfig == nil {
			// Only log the error at warn if it's unexpected.
			// If we cancelled the request then we're shutting down and this isn't unexpected.
			if errors.Is(ccc.Error(), ErrRequestCanceled) || errors.Is(ccc.Error(), ErrShutdown) {
				logDebugf("CCCPPOLL: CCCP request was cancelled.")
			} else {
				logWarnf("CCCPPOLL: Failed to retrieve config from any node.")
			}
			continue
		}

		logDebugf("CCCPPOLL: Received new config")
		ccc.cfgMgr.OnNewConfig(foundConfig)
	}

	close(ccc.looperDoneSig)
	return nil
}

func (ccc *cccpConfigController) getClusterConfig(pipeline *memdPipeline) (cfgOut []byte, errOut error) {
	signal := make(chan struct{}, 1)
	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:   memd.CmdMagicReq,
			Command: memd.CmdGetClusterConfig,
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

	timeoutTmr := AcquireTimer(ccc.confCccpMaxWait)
	select {
	case <-signal:
		ReleaseTimer(timeoutTmr, false)
		return
	case <-timeoutTmr.C:
		ReleaseTimer(timeoutTmr, true)

		// We've timed out so lets check underlying connections to see if they're responsible.
		clients := pipeline.Clients()
		for _, cli := range clients {
			err := cli.Error()
			if err != nil {
				req.cancelWithCallback(err)
				<-signal
				return
			}
		}
		req.cancelWithCallback(errAmbiguousTimeout)
		<-signal
		return
	case <-ccc.looperStopSig:
		ReleaseTimer(timeoutTmr, false)
		req.cancelWithCallback(errRequestCanceled)
		<-signal
		return

	}
}
