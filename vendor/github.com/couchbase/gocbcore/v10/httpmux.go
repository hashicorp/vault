package gocbcore

import (
	"bytes"
	"fmt"
	"net/url"
	"sync/atomic"
	"unsafe"
)

type httpMux struct {
	muxPtr        unsafe.Pointer
	breakerCfg    CircuitBreakerConfig
	cfgMgr        configManager
	noSeedNodeTLS bool
}

func newHTTPMux(breakerCfg CircuitBreakerConfig, cfgMgr configManager, muxState *httpClientMux, noSeedNodeTLS bool) *httpMux {
	mux := &httpMux{
		breakerCfg:    breakerCfg,
		cfgMgr:        cfgMgr,
		muxPtr:        unsafe.Pointer(muxState),
		noSeedNodeTLS: noSeedNodeTLS,
	}

	cfgMgr.AddConfigWatcher(mux)

	return mux
}

func (mux *httpMux) Get() *httpClientMux {
	muxCfg := atomic.LoadPointer(&mux.muxPtr)
	if muxCfg == nil {
		return nil
	}
	return (*httpClientMux)(muxCfg)
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
	if oldHTTPMux == nil {
		logWarnf("HTTP mux received new route config after shutdown")
		return
	}

	endpoints := mux.buildEndpoints(cfg, oldHTTPMux.tlsConfig != nil)

	var buffer bytes.Buffer
	addEps := func(title string, eps []routeEndpoint) {
		fmt.Fprintf(&buffer, "%s Eps:\n", title)
		for _, ep := range eps {
			fmt.Fprintf(&buffer, "  - %s\n", ep.Address)
		}
	}

	buffer.WriteString(fmt.Sprintln("HTTP muxer applying endpoints:"))
	buffer.WriteString(fmt.Sprintf("Bucket: %s\n", cfg.name))
	addEps("Capi", endpoints.capiEpList)
	addEps("Mgmt", endpoints.mgmtEpList)
	addEps("N1ql", endpoints.n1qlEpList)
	addEps("FTS", endpoints.ftsEpList)
	addEps("CBAS", endpoints.cbasEpList)
	addEps("Eventing", endpoints.eventingEpList)
	addEps("GSI", endpoints.gsiEpList)
	addEps("Backup", endpoints.backupEpList)

	logDebugf(buffer.String())

	newHTTPMux := newHTTPClientMux(cfg, endpoints, oldHTTPMux.tlsConfig, oldHTTPMux.auth, mux.breakerCfg)

	if !mux.Update(oldHTTPMux, newHTTPMux) {
		logDebugf("Failed to update HTTP mux")
	}
}

func (mux *httpMux) UpdateTLS(tlsConfig *dynTLSConfig, auth AuthProvider) {
	oldMux := mux.Get()
	if oldMux == nil {
		logWarnf("HTTP mux received TLS update after shutdown")
		return
	}

	endpoints := mux.buildEndpoints(&oldMux.srcConfig, tlsConfig != nil)

	newMux := newHTTPClientMux(&oldMux.srcConfig, endpoints, tlsConfig, auth, oldMux.breakerCfg)
	if !atomic.CompareAndSwapPointer(&mux.muxPtr, unsafe.Pointer(oldMux), unsafe.Pointer(newMux)) {
		// A new config must have come in so let's try again.
		mux.UpdateTLS(tlsConfig, auth)
	}
}

func makeEpList(endpoints []routeEndpoint) []string {
	var epList []string
	for _, ep := range endpoints {
		epList = append(epList, ep.Address)
	}

	return epList
}

// CapiEps returns the capi endpoints with the path escaped bucket name appended.
func (mux *httpMux) CapiEps() []string {
	clientMux := mux.Get()
	if clientMux == nil {
		return nil
	}

	var epList []string
	for _, ep := range clientMux.capiEpList {
		epList = append(epList, ep.Address+"/"+url.PathEscape(clientMux.bucket))
	}

	return epList
}

func (mux *httpMux) MgmtEps() []string {
	clientMux := mux.Get()
	if clientMux == nil {
		return nil
	}

	return makeEpList(clientMux.mgmtEpList)
}

func (mux *httpMux) N1qlEps() []string {
	clientMux := mux.Get()
	if clientMux == nil {
		return nil
	}

	return makeEpList(clientMux.n1qlEpList)
}

func (mux *httpMux) CbasEps() []string {
	clientMux := mux.Get()
	if clientMux == nil {
		return nil
	}

	return makeEpList(clientMux.cbasEpList)
}

func (mux *httpMux) FtsEps() []string {
	clientMux := mux.Get()
	if clientMux == nil {
		return nil
	}

	return makeEpList(clientMux.ftsEpList)
}

func (mux *httpMux) EventingEps() []string {
	if cMux := mux.Get(); cMux != nil {
		return makeEpList(cMux.eventingEpList)
	}

	return nil
}

func (mux *httpMux) GSIEps() []string {
	if cMux := mux.Get(); cMux != nil {
		return makeEpList(cMux.gsiEpList)
	}

	return nil
}

func (mux *httpMux) BackupEps() []string {
	if cMux := mux.Get(); cMux != nil {
		return makeEpList(cMux.backupEpList)
	}

	return nil
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

func (mux *httpMux) Auth() AuthProvider {
	clientMux := mux.Get()
	if clientMux == nil {
		return nil
	}

	return clientMux.auth
}

func (mux *httpMux) buildEndpoints(config *routeConfig, useTLS bool) httpClientMuxEndpoints {
	var endpoints httpClientMuxEndpoints
	if useTLS {
		if mux.noSeedNodeTLS {
			endpoints = httpClientMuxEndpoints{
				capiEpList:     mux.buildSSLEpListWithNoSSLSeed(config.capiEpList),
				mgmtEpList:     mux.buildSSLEpListWithNoSSLSeed(config.mgmtEpList),
				n1qlEpList:     mux.buildSSLEpListWithNoSSLSeed(config.n1qlEpList),
				ftsEpList:      mux.buildSSLEpListWithNoSSLSeed(config.ftsEpList),
				cbasEpList:     mux.buildSSLEpListWithNoSSLSeed(config.cbasEpList),
				eventingEpList: mux.buildSSLEpListWithNoSSLSeed(config.eventingEpList),
				gsiEpList:      mux.buildSSLEpListWithNoSSLSeed(config.gsiEpList),
				backupEpList:   mux.buildSSLEpListWithNoSSLSeed(config.backupEpList),
			}
		} else {
			endpoints = httpClientMuxEndpoints{
				capiEpList:     config.capiEpList.SSLEndpoints,
				mgmtEpList:     config.mgmtEpList.SSLEndpoints,
				n1qlEpList:     config.n1qlEpList.SSLEndpoints,
				ftsEpList:      config.ftsEpList.SSLEndpoints,
				cbasEpList:     config.cbasEpList.SSLEndpoints,
				eventingEpList: config.eventingEpList.SSLEndpoints,
				gsiEpList:      config.gsiEpList.SSLEndpoints,
				backupEpList:   config.backupEpList.SSLEndpoints,
			}
		}
	} else {
		endpoints = httpClientMuxEndpoints{
			capiEpList:     config.capiEpList.NonSSLEndpoints,
			mgmtEpList:     config.mgmtEpList.NonSSLEndpoints,
			n1qlEpList:     config.n1qlEpList.NonSSLEndpoints,
			ftsEpList:      config.ftsEpList.NonSSLEndpoints,
			cbasEpList:     config.cbasEpList.NonSSLEndpoints,
			eventingEpList: config.eventingEpList.NonSSLEndpoints,
			gsiEpList:      config.gsiEpList.NonSSLEndpoints,
			backupEpList:   config.backupEpList.NonSSLEndpoints,
		}
	}

	return endpoints
}

func (mux *httpMux) buildSSLEpListWithNoSSLSeed(list routeEndpoints) []routeEndpoint {
	var newlist []routeEndpoint
	for _, ep := range list.SSLEndpoints {
		if !ep.IsSeedNode {
			newlist = append(newlist, ep)
		}
	}
	for _, ep := range list.NonSSLEndpoints {
		if ep.IsSeedNode {
			newlist = append(newlist, ep)
			break
		}
	}

	return newlist
}
