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
}

type configPollerController interface {
	Pause(paused bool)
	Done() chan struct{}
	Stop()
}

func newPollerController(cccpPoller *cccpConfigController, httpPoller *httpConfigController) *pollerController {
	return &pollerController{
		cccpPoller: cccpPoller,
		httpPoller: httpPoller,
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

func isPollingFallbackError(err error) bool {
	return errors.Is(err, ErrDocumentNotFound) || errors.Is(err, ErrUnsupportedOperation)
}
