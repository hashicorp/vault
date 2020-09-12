package gocbcore

import (
	"fmt"
)

type kvMuxState struct {
	pipelines []*memdPipeline
	deadPipe  *memdPipeline

	kvServerList []string
	bktType      bucketType
	vbMap        *vbucketMap
	ketamaMap    *ketamaContinuum
	uuid         string
	revID        int64

	durabilityLevelStatus durabilityLevelStatus
	createAsDeletedStatus createAsDeletedStatus
	collectionsSupported  bool
}

func newKVMuxState(cfg *routeConfig, pipelines []*memdPipeline, deadpipe *memdPipeline) *kvMuxState {
	mux := &kvMuxState{
		pipelines: pipelines,
		deadPipe:  deadpipe,

		kvServerList: cfg.kvServerList,
		bktType:      cfg.bktType,
		vbMap:        cfg.vbMap,
		ketamaMap:    cfg.ketamaMap,
		uuid:         cfg.uuid,
		revID:        cfg.revID,

		durabilityLevelStatus: durabilityLevelStatusUnknown,
		createAsDeletedStatus: createAsDeletedStatusUnknown,

		collectionsSupported: cfg.ContainsBucketCapability("collections"),
	}

	// We setup with a fake config, this means that durability support is still unknown.
	if cfg.revID > -1 {
		if cfg.ContainsBucketCapability("durableWrite") {
			mux.durabilityLevelStatus = durabilityLevelStatusSupported
		} else {
			mux.durabilityLevelStatus = durabilityLevelStatusUnsupported
		}

		if cfg.ContainsBucketCapability("tombstonedUserXAttrs") {
			mux.createAsDeletedStatus = createAsDeletedStatusSupported
		} else {
			mux.createAsDeletedStatus = createAsDeletedStatusUnsupported
		}
	}

	return mux
}

func (mux *kvMuxState) BucketType() bucketType {
	return mux.bktType
}

func (mux *kvMuxState) NumPipelines() int {
	return len(mux.pipelines)
}

func (mux *kvMuxState) GetPipeline(index int) *memdPipeline {
	if index < 0 || index >= len(mux.pipelines) {
		return mux.deadPipe
	}
	return mux.pipelines[index]
}

// nolint: unused
func (mux *kvMuxState) debugString() string {
	var outStr string

	for i, n := range mux.pipelines {
		outStr += fmt.Sprintf("Pipeline %d:\n", i)
		outStr += reindentLog("  ", n.debugString()) + "\n"
	}

	outStr += "Dead Pipeline:\n"
	if mux.deadPipe != nil {
		outStr += reindentLog("  ", mux.deadPipe.debugString()) + "\n"
	} else {
		outStr += "  Disabled\n"
	}

	return outStr
}
