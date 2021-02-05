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

	bucketCapabilities   map[BucketCapability]BucketCapabilityStatus
	collectionsSupported bool
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

		bucketCapabilities: map[BucketCapability]BucketCapabilityStatus{
			BucketCapabilityDurableWrites:        BucketCapabilityStatusUnknown,
			BucketCapabilityCreateAsDeleted:      BucketCapabilityStatusUnknown,
			BucketCapabilityReplaceBodyWithXattr: BucketCapabilityStatusUnknown,
		},

		collectionsSupported: cfg.ContainsBucketCapability("collections"),
	}

	// We setup with a fake config, this means that durability support is still unknown.
	if cfg.revID > -1 {
		if cfg.ContainsBucketCapability("durableWrite") {
			mux.bucketCapabilities[BucketCapabilityDurableWrites] = BucketCapabilityStatusSupported
		} else {
			mux.bucketCapabilities[BucketCapabilityDurableWrites] = BucketCapabilityStatusUnsupported
		}

		if cfg.ContainsBucketCapability("tombstonedUserXAttrs") {
			mux.bucketCapabilities[BucketCapabilityCreateAsDeleted] = BucketCapabilityStatusSupported
		} else {
			mux.bucketCapabilities[BucketCapabilityCreateAsDeleted] = BucketCapabilityStatusUnsupported
		}

		if cfg.ContainsBucketCapability("subdoc.ReplaceBodyWithXattr") {
			mux.bucketCapabilities[BucketCapabilityReplaceBodyWithXattr] = BucketCapabilityStatusSupported
		} else {
			mux.bucketCapabilities[BucketCapabilityReplaceBodyWithXattr] = BucketCapabilityStatusUnsupported
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

func (mux *kvMuxState) HasBucketCapabilityStatus(cap BucketCapability, status BucketCapabilityStatus) bool {
	st, ok := mux.bucketCapabilities[cap]
	if !ok {
		return status == BucketCapabilityStatusUnsupported
	}

	return st == status
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
