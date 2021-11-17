package gocbcore

import (
	"errors"
	"sync"
	"sync/atomic"
)

type pollerController struct {
	activeController configPollerController
	controllerLock   sync.Mutex
	stopped          bool
	bucketConfigSeen uint32

	cccpPoller *cccpConfigController
	httpPoller *httpConfigController
	cfgMgr     configManager
}

type configPollerController interface {
	Pause(paused bool)
	Done() chan struct{}
	Stop()
	Reset()
	Error() error
}

func newPollerController(cccpPoller *cccpConfigController, httpPoller *httpConfigController, cfgMgr configManager) *pollerController {
	pc := &pollerController{
		cccpPoller: cccpPoller,
		httpPoller: httpPoller,
		cfgMgr:     cfgMgr,
	}
	cfgMgr.AddConfigWatcher(pc)

	return pc
}

// We listen out for every config that comes in so that we (re)start the cccp if applicable.
func (pc *pollerController) OnNewRouteConfig(cfg *routeConfig) {
	if cfg.bktType != bktTypeCouchbase && cfg.bktType != bktTypeMemcached {
		return
	}
	atomic.SwapUint32(&pc.bucketConfigSeen, 1)

	if cfg.bktType == bktTypeMemcached {
		pc.cfgMgr.RemoveConfigWatcher(pc)
		return
	}

	go func() {
		pc.controllerLock.Lock()
		if pc.stopped {
			pc.controllerLock.Unlock()
			return
		}
		if pc.activeController == pc.httpPoller {
			logInfof("Found couchbase bucket and HTTP poller in use. Resetting pollers to start cccp.")
			pc.activeController = nil
			pc.controllerLock.Unlock()

			pc.httpPoller.Stop()
			pollerCh := pc.httpPoller.Done()
			if pollerCh != nil {
				<-pollerCh
			}
			pc.httpPoller.Reset()
			pc.cccpPoller.Reset()
			pc.Start()
		} else {
			pc.controllerLock.Unlock()
		}
	}()
}

func (pc *pollerController) Start() {
	pc.controllerLock.Lock()
	if pc.stopped {
		pc.controllerLock.Unlock()
		return
	}

	atomic.SwapUint32(&pc.bucketConfigSeen, 0)
	if pc.cccpPoller == nil {
		pc.activeController = pc.httpPoller
		pc.controllerLock.Unlock()
		pc.httpPoller.DoLoop()
		return
	}
	pc.activeController = pc.cccpPoller
	pc.controllerLock.Unlock()
	err := pc.cccpPoller.DoLoop()
	if err != nil {
		if pc.httpPoller == nil {
			logErrorf("CCCP poller has exited for http fallback but no http poller is configured")
			return
		}
		if isPollingFallbackError(err) {
			pc.controllerLock.Lock()
			// We can get into a weird race where the poller controller sent stop to the active controller but we then
			// swap to a different one and so the Done() function never completes.
			if pc.stopped {
				pc.activeController = nil
				pc.controllerLock.Unlock()
			} else {
				pc.activeController = pc.httpPoller
				pc.controllerLock.Unlock()
				pc.httpPoller.DoLoop()
			}
		}
	}
}

func (pc *pollerController) Pause(paused bool) {
	pc.controllerLock.Lock()
	controller := pc.activeController
	pc.controllerLock.Unlock()
	if controller != nil {
		controller.Pause(paused)
	}
}

// Stop should never be called more than once.
func (pc *pollerController) Stop() {
	pc.controllerLock.Lock()
	pc.stopped = true
	controller := pc.activeController
	pc.controllerLock.Unlock()

	if controller != nil {
		controller.Stop()
	}
}

func (pc *pollerController) Done() chan struct{} {
	pc.controllerLock.Lock()
	controller := pc.activeController
	pc.controllerLock.Unlock()

	if controller == nil {
		return nil
	}
	return controller.Done()
}

type pollerErrorProvider interface {
	PollerError() error
}

// If the underlying poller is currently in an error state then this will surface that error.
func (pc *pollerController) PollerError() error {
	pc.controllerLock.Lock()
	controller := pc.activeController
	pc.controllerLock.Unlock()

	if controller == nil {
		return nil
	}

	return controller.Error()
}

func (pc *pollerController) ForceHTTPPoller() {
	go func() {
		if atomic.LoadUint32(&pc.bucketConfigSeen) == 1 {
			logInfof("Config already seen, not forcing HTTP")
			// If we've seen a config already then either cccp or http polling have managed to fetch a config and
			// bucket type can't have changed so there's no reason to fallback.
			return
		}

		pc.controllerLock.Lock()
		if pc.stopped || pc.activeController == nil {
			// If active controller is nil at this point then something strange is happening, we're trying to force
			// http polling at the same time as we've received a config via http polling and are attempting to reset to
			// use cccp polling (which means that the server must support cccp. If this happens let's just let
			// cccp start up.
			pc.controllerLock.Unlock()
			return
		}
		if pc.activeController == pc.cccpPoller {
			logInfof("Stopping CCCP poller for HTTP polling takeover")
			pc.cccpPoller.Stop()
			pollerCh := pc.cccpPoller.Done()
			if pollerCh != nil {
				<-pollerCh
			}
			pc.httpPoller.Reset()
			pc.cccpPoller.Reset()
			if atomic.LoadUint32(&pc.bucketConfigSeen) == 1 {
				pc.controllerLock.Unlock()
				logInfof("Config seen whilst waiting for CCCP poller to stop, restarting CCCP poller.")
				// CCCP managed to fetch a config whilst we were waiting for shutdown, in this case we want to just
				// start CCCP again as the bucket must exist and be a couchbase bucket.
				pc.Start()
				return
			}
		} else if pc.activeController == pc.httpPoller {
			pc.controllerLock.Unlock()
			return
		}

		pc.activeController = pc.httpPoller
		pc.controllerLock.Unlock()

		pc.httpPoller.DoLoop()
	}()
}

func isPollingFallbackError(err error) bool {
	return errors.Is(err, ErrDocumentNotFound) || errors.Is(err, ErrUnsupportedOperation) ||
		errors.Is(err, errNoCCCPHosts) || errors.Is(err, ErrBucketNotFound)
}
