package gocbcore

type httpConfigController struct {
	muxer     *httpMux
	seenNodes map[string]uint64
	*baseHTTPConfigController
}

func newHTTPConfigController(bucketName string, props httpPollerProperties, muxer *httpMux,
	cfgMgr *configManagementComponent) *httpConfigController {
	ctrlr := &httpConfigController{
		muxer:     muxer,
		seenNodes: make(map[string]uint64),
	}

	ctrlr.baseHTTPConfigController = newBaseHTTPConfigController(bucketName, props, cfgMgr, ctrlr.GetEndpoint)

	return ctrlr
}

func (hcc *httpConfigController) GetEndpoint(iterNum uint64) string {
	var pickedSrv string
	for _, srv := range hcc.muxer.MgmtEps() {
		if hcc.seenNodes[srv] >= iterNum {
			continue
		}
		pickedSrv = srv
		break
	}

	if pickedSrv != "" {
		hcc.seenNodes[pickedSrv] = iterNum
	}

	return pickedSrv
}

func (hcc *httpConfigController) CanPoll() bool {
	return len(hcc.muxer.MgmtEps()) > 0
}
