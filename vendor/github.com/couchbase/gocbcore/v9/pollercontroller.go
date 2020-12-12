package gocbcore

import (
	"errors"
	"sync"
)

type pollerController struct {
	activeController configPollerController
	controllerLock   sync.Mutex
	stopped          bool

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

// We listen out for the first config that comes in so that we (re)start the cccp if applicable.
func (pc *pollerController) OnNewRouteConfig(cfg *routeConfig) {
	if cfg.bktType == bktTypeCouchbase || cfg.bktType == bktTypeMemcached {
		pc.cfgMgr.RemoveConfigWatcher(pc)
	}

	pc.controllerLock.Lock()
	if cfg.bktType == bktTypeCouchbase && pc.activeController == pc.httpPoller {
		logDebugf("Found couchbase bucket and HTTP poller in use. Resetting pollers to start cccp.")
		pc.activeController = nil
		pc.controllerLock.Unlock()
		go func() {
			pc.httpPoller.Stop()
			pollerCh := pc.httpPoller.Done()
			if pollerCh != nil {
				<-pollerCh
			}
			pc.httpPoller.Reset()
			pc.cccpPoller.Reset()
			pc.Start()
		}()
	} else {
		pc.controllerLock.Unlock()
	}
}

func (pc *pollerController) Start() {
	pc.controllerLock.Lock()
	if pc.stopped {
		pc.controllerLock.Unlock()
		return
	}

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

func isPollingFallbackError(err error) bool {
	return errors.Is(err, ErrDocumentNotFound) || errors.Is(err, ErrUnsupportedOperation) ||
		errors.Is(err, errNoCCCPHosts)
}
