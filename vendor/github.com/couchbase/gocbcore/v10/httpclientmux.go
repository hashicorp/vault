package gocbcore

type httpClientMuxEndpoints struct {
	capiEpList     []routeEndpoint
	mgmtEpList     []routeEndpoint
	n1qlEpList     []routeEndpoint
	ftsEpList      []routeEndpoint
	cbasEpList     []routeEndpoint
	eventingEpList []routeEndpoint
	gsiEpList      []routeEndpoint
	backupEpList   []routeEndpoint
}

type httpClientMux struct {
	capiEpList     []routeEndpoint
	mgmtEpList     []routeEndpoint
	n1qlEpList     []routeEndpoint
	ftsEpList      []routeEndpoint
	cbasEpList     []routeEndpoint
	eventingEpList []routeEndpoint
	gsiEpList      []routeEndpoint
	backupEpList   []routeEndpoint

	bucket string

	uuid       string
	revID      int64
	breakerCfg CircuitBreakerConfig

	srcConfig routeConfig

	tlsConfig *dynTLSConfig
	auth      AuthProvider
}

func newHTTPClientMux(cfg *routeConfig, endpoints httpClientMuxEndpoints, tlsConfig *dynTLSConfig, auth AuthProvider,
	breakerCfg CircuitBreakerConfig) *httpClientMux {
	return &httpClientMux{
		capiEpList:     endpoints.capiEpList,
		mgmtEpList:     endpoints.mgmtEpList,
		n1qlEpList:     endpoints.n1qlEpList,
		ftsEpList:      endpoints.ftsEpList,
		cbasEpList:     endpoints.cbasEpList,
		eventingEpList: endpoints.eventingEpList,
		gsiEpList:      endpoints.gsiEpList,
		backupEpList:   endpoints.backupEpList,

		bucket: cfg.name,

		uuid:       cfg.uuid,
		revID:      cfg.revID,
		breakerCfg: breakerCfg,

		srcConfig: *cfg,

		tlsConfig: tlsConfig,
		auth:      auth,
	}
}
