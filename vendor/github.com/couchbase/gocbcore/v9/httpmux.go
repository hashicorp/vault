package gocbcore

import (
	"sync/atomic"
	"unsafe"
)

type httpMux struct {
	muxPtr     unsafe.Pointer
	breakerCfg CircuitBreakerConfig
	cfgMgr     configManager
}

func newHTTPMux(breakerCfg CircuitBreakerConfig, cfgMgr configManager) *httpMux {
	mux := &httpMux{
		breakerCfg: breakerCfg,
		cfgMgr:     cfgMgr,
	}

	cfgMgr.AddConfigWatcher(mux)

	return mux
}

func (mux *httpMux) Get() *httpClientMux {
	return (*httpClientMux)(atomic.LoadPointer(&mux.muxPtr))
}

func (mux *httpMux) Update(old, new *httpClientMux) bool {
	if new == nil {
		logErrorf("Attempted to update to nil httpClientMux")
		return false
	}

	if old != nil {
		return atomic.CompareAndSwapPointer(&mux.muxPtr, unsafe.Pointer(old), unsafe.Pointer(new))
	}

	if atomic.SwapPointer(&mux.muxPtr, unsafe.Pointer(new)) != nil {
		logErrorf("Updated from nil attempted on initialized httpClientMux")
		return false
	}

	return true
}

func (mux *httpMux) Clear() *httpClientMux {
	val := atomic.SwapPointer(&mux.muxPtr, nil)
	return (*httpClientMux)(val)
}

func (mux *httpMux) OnNewRouteConfig(cfg *routeConfig) {
	oldHTTPMux := mux.Get()

	newHTTPMux := newHTTPClientMux(cfg, mux.breakerCfg)

	mux.Update(oldHTTPMux, newHTTPMux)
}

func (mux *httpMux) CapiEps() []string {
	clientMux := mux.Get()
	if clientMux == nil {
		return nil
	}

	return clientMux.capiEpList
}

func (mux *httpMux) MgmtEps() []string {
	clientMux := mux.Get()
	if clientMux == nil {
		return nil
	}

	return clientMux.mgmtEpList
}

func (mux *httpMux) N1qlEps() []string {
	clientMux := mux.Get()
	if clientMux == nil {
		return nil
	}

	return clientMux.n1qlEpList
}

func (mux *httpMux) CbasEps() []string {
	clientMux := mux.Get()
	if clientMux == nil {
		return nil
	}

	return clientMux.cbasEpList
}

func (mux *httpMux) FtsEps() []string {
	clientMux := mux.Get()
	if clientMux == nil {
		return nil
	}

	return clientMux.ftsEpList
}

func (mux *httpMux) ConfigRev() (int64, error) {
	clientMux := mux.Get()
	if clientMux == nil {
		return 0, errShutdown
	}

	return clientMux.revID, nil
}

func (mux *httpMux) Close() error {
	mux.cfgMgr.RemoveConfigWatcher(mux)
	mux.Clear()
	return nil
}
