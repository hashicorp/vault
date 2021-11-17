package gocbcore

import (
	"fmt"
)

type kvMuxState struct {
	pipelines []*memdPipeline
	deadPipe  *memdPipeline

	routeCfg *routeConfig

	bucketCapabilities   map[BucketCapability]BucketCapabilityStatus
	collectionsSupported bool
}

func newKVMuxState(cfg *routeConfig, pipelines []*memdPipeline, deadpipe *memdPipeline) *kvMuxState {
	mux := &kvMuxState{
		pipelines: pipelines,
		deadPipe:  deadpipe,

		routeCfg: cfg,

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

func (mux *kvMuxState) RouteConfig() *routeConfig {
	return mux.routeCfg
}

func (mux *kvMuxState) RevID() int64 {
	return mux.routeCfg.revID
}

func (mux *kvMuxState) VBMap() *vbucketMap {
	return mux.routeCfg.vbMap
}

func (mux *kvMuxState) UUID() string {
	return mux.routeCfg.uuid
}

func (mux *kvMuxState) KetamaMap() *ketamaContinuum {
	return mux.routeCfg.ketamaMap
}

func (mux *kvMuxState) BucketType() bucketType {
	return mux.routeCfg.bktType
}

func (mux *kvMuxState) KVEps() []string {
	return mux.routeCfg.kvServerList
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

func (mux *kvMuxState) BucketCapabilityStatus(cap BucketCapability) BucketCapabilityStatus {
	st, ok := mux.bucketCapabilities[cap]
	if !ok {
		return BucketCapabilityStatusUnsupported
	}

	return st
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
