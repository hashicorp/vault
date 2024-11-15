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
	cccpFetcher        *cccpConfigFetcher

	looperStopSig chan struct{}

	fetchErr error
	errLock  sync.Mutex

	isFallbackErrorFn func(error) bool
	noConfigFoundFn   func(error)
}

func newCCCPConfigController(props cccpPollerProperties, muxer dispatcher, cfgMgr *configManagementComponent,
	isFallbackErrorFn func(error) bool, noConfigFoundFn func(error)) *cccpConfigController {
	return &cccpConfigController{
		muxer:              muxer,
		cfgMgr:             cfgMgr,
		confCccpPollPeriod: props.confCccpPollPeriod,
		cccpFetcher:        props.cccpConfigFetcher,

		looperStopSig: make(chan struct{}),

		isFallbackErrorFn: isFallbackErrorFn,
		noConfigFoundFn:   noConfigFoundFn,
	}
}

type cccpPollerProperties struct {
	confCccpPollPeriod time.Duration
	cccpConfigFetcher  *cccpConfigFetcher
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

func (ccc *cccpConfigController) Stop() {
	logInfof("CCCP Looper stopping")
	close(ccc.looperStopSig)
}

// Reset must never be called concurrently with the Stop or whilst the poll loop is running.
func (ccc *cccpConfigController) Reset() {
	ccc.looperStopSig = make(chan struct{})
}

func (ccc *cccpConfigController) DoLoop() error {
	if err := ccc.doLoop(); err != nil {
		logInfof("CCCP Looper errored")

		return err
	}

	logInfof("CCCP Looper stopped")
	return nil
}

func (ccc *cccpConfigController) doLoop() error {
	tickTime := ccc.confCccpPollPeriod

	logInfof("CCCP Looper starting.")
	nodeIdx := -1
	// The first time that we loop we want to skip any sleep so that we can try get a config and bootstrapped ASAP.
	firstLoop := true

	for {
		if !firstLoop {
			// Wait for either the agent to be shut down, or our tick time to expire
			select {
			case <-ccc.looperStopSig:
				return nil
			case <-time.After(tickTime):
			}
		}
		firstLoop = false

		iter, err := ccc.muxer.PipelineSnapshot()
		if err != nil {
			// If we have an error it indicates the client is shut down.
			break
		}

		numNodes := iter.NumPipelines()
		if numNodes == 0 {
			logInfof("CCCPPOLL: No nodes available to poll, returning upstream")
			return errNoCCCPHosts
		}

		if nodeIdx < 0 || nodeIdx > numNodes {
			nodeIdx = rand.Intn(numNodes) // #nosec G404
		}

		var foundConfig *cfgBucket
		var configAlreadyLatest bool
		var fallbackErr error
		var wasCancelled bool
		var numNodesSupportNotifs int
		iter.Iterate(nodeIdx, func(pipeline *memdPipeline) bool {
			nodeIdx = (nodeIdx + 1) % numNodes
			if pipeline.SupportsFeature(memd.FeatureClustermapChangeNotificationBrief) {
				numNodesSupportNotifs++
				return false
			}

			cccpBytes, err := ccc.getClusterConfig(pipeline)
			if err != nil {
				if ccc.isFallbackErrorFn(err) {
					fallbackErr = err
					return false
				}

				// Only log the error at warn if it's unexpected.
				// If we cancelled the request or we're shutting down the connection then it's not really unexpected.
				if errors.Is(err, ErrRequestCanceled) || errors.Is(err, ErrShutdown) {
					wasCancelled = true
					logDebugf("CCCPPOLL: CCCP request was cancelled or connection was shutdown: %v", err)
					return true
				}

				// This error is checked by WaitUntilReady when no config has been seen.
				ccc.setError(err)

				logWarnf("CCCPPOLL: Failed to retrieve CCCP config. %s", err)
				return false
			}
			fallbackErr = nil
			ccc.setError(nil)

			if len(cccpBytes) > 0 {
				logDebugf("CCCPPOLL: Got Block: %s", string(cccpBytes))

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
			} else {
				configAlreadyLatest = true
			}
			return true
		})
		if fallbackErr != nil {
			// This error is indicative of a memcached bucket which we can't handle so return the error.
			logInfof("CCCPPOLL: CCCP not supported, returning error upstream.")
			return fallbackErr
		}

		if numNodesSupportNotifs == numNodes {
			continue
		}

		if configAlreadyLatest {
			logDebugf("CCCPPOLL: Received empty config")
			continue
		}

		if foundConfig == nil {
			// Only log the error at warn if it's unexpected.
			// If we cancelled the request then we're shutting down or request was requeued and this isn't unexpected.
			if wasCancelled {
				logDebugf("CCCPPOLL: CCCP request was cancelled.")
			} else {
				logWarnf("CCCPPOLL: Failed to retrieve config from any node.")
				ccc.noConfigFoundFn(err)
			}
			continue
		}

		logDebugf("CCCPPOLL: Received new config")
		ccc.cfgMgr.OnNewConfig(foundConfig)

	}

	return nil
}

func (ccc *cccpConfigController) getClusterConfig(pipeline *memdPipeline) ([]byte, error) {
	revID, revEpoch := ccc.cfgMgr.CurrentRev()
	cfg, err := ccc.cccpFetcher.GetClusterConfig(pipeline, revID, revEpoch, ccc.looperStopSig)
	if err != nil {
		if errors.Is(err, ErrTimeout) {
			// We've timed out so lets check underlying connections to see if they're responsible.
			clients := pipeline.Clients()
			for _, cli := range clients {
				err := cli.Error()
				if err != nil {
					logDebugf("Found error in pipeline client %p/%s: %v", cli, cli.address, err)
					return nil, err
				}
			}
		}

		return nil, err
	}

	return cfg, nil
}
