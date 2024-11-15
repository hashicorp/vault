package gocbcore

type seedConfigController struct {
	*baseHTTPConfigController
	seed       string
	iterNum    uint64
	stoppedSig chan struct{}
}

func newSeedConfigController(seed, bucketName string, props httpPollerProperties,
	cfgMgr *configManagementComponent) *seedConfigController {
	scc := &seedConfigController{
		seed:       seed,
		stoppedSig: make(chan struct{}),
	}
	scc.baseHTTPConfigController = newBaseHTTPConfigController(bucketName, props, cfgMgr, scc.GetEndpoint)

	return scc
}

func (scc *seedConfigController) GetEndpoint(iterNum uint64) string {
	if scc.iterNum == iterNum {
		return ""
	}

	scc.iterNum = iterNum
	return scc.seed
}

func (scc *seedConfigController) Stop() {
	logInfof("Seed poller stopping.")
	scc.baseHTTPConfigController.Stop()
	<-scc.stoppedSig
}

func (scc *seedConfigController) Run() {
	scc.DoLoop()
	close(scc.stoppedSig)
}

func (scc *seedConfigController) PollerError() error {
	return scc.Error()
}

// We're already a http poller so do nothing
func (scc *seedConfigController) ForceHTTPPoller() {
}
