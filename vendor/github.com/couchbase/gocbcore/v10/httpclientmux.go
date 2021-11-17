package gocbcore

type httpClientMux struct {
	capiEpList     []string
	mgmtEpList     []string
	n1qlEpList     []string
	ftsEpList      []string
	cbasEpList     []string
	eventingEpList []string
	gsiEpList      []string
	backupEpList   []string

	bucket string

	uuid       string
	revID      int64
	breakerCfg CircuitBreakerConfig
}

func newHTTPClientMux(cfg *routeConfig, breakerCfg CircuitBreakerConfig) *httpClientMux {
	return &httpClientMux{
		capiEpList:     cfg.capiEpList,
		mgmtEpList:     cfg.mgmtEpList,
		n1qlEpList:     cfg.n1qlEpList,
		ftsEpList:      cfg.ftsEpList,
		cbasEpList:     cfg.cbasEpList,
		eventingEpList: cfg.eventingEpList,
		gsiEpList:      cfg.gsiEpList,
		backupEpList:   cfg.backupEpList,

		bucket: cfg.name,

		uuid:       cfg.uuid,
		revID:      cfg.revID,
		breakerCfg: breakerCfg,
	}
}
