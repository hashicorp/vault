package gocbcore

import (
	"fmt"
)

type kvMuxState struct {
	pipelines []*memdPipeline
	deadPipe  *memdPipeline

	routeCfg routeConfig

	expectedBucketName   string
	bucketCapabilities   map[BucketCapability]CapabilityStatus
	collectionsSupported bool

	kvServerList   []routeEndpoint
	tlsConfig      *dynTLSConfig
	authMechanisms []AuthMechanism
	auth           AuthProvider
}

func newKVMuxState(cfg *routeConfig, kvServerList []routeEndpoint, tlsConfig *dynTLSConfig,
	authMechanisms []AuthMechanism, auth AuthProvider, expectedBucketName string, pipelines []*memdPipeline, deadpipe *memdPipeline) *kvMuxState {
	mux := &kvMuxState{
		pipelines: pipelines,
		deadPipe:  deadpipe,

		routeCfg: *cfg,

		expectedBucketName: expectedBucketName,
		bucketCapabilities: map[BucketCapability]CapabilityStatus{
			BucketCapabilityDurableWrites:        CapabilityStatusUnknown,
			BucketCapabilityCreateAsDeleted:      CapabilityStatusUnknown,
			BucketCapabilityReplaceBodyWithXattr: CapabilityStatusUnknown,
			BucketCapabilityRangeScan:            CapabilityStatusUnknown,
			BucketCapabilityReplicaRead:          CapabilityStatusUnknown,
			BucketCapabilityNonDedupedHistory:    CapabilityStatusUnknown,
		},

		collectionsSupported: cfg.ContainsBucketCapability("collections"),

		kvServerList:   kvServerList,
		tlsConfig:      tlsConfig,
		authMechanisms: authMechanisms,
		auth:           auth,
	}

	// We setup with a fake config, this means that durability support is still unknown.
	// We only want to update bucket capabilities once we actually have a bucket config.
	if cfg.revID > -1 && cfg.name == expectedBucketName {
		if cfg.ContainsBucketCapability("durableWrite") {
			mux.bucketCapabilities[BucketCapabilityDurableWrites] = CapabilityStatusSupported
		} else {
			mux.bucketCapabilities[BucketCapabilityDurableWrites] = CapabilityStatusUnsupported
		}

		if cfg.ContainsBucketCapability("tombstonedUserXAttrs") {
			mux.bucketCapabilities[BucketCapabilityCreateAsDeleted] = CapabilityStatusSupported
		} else {
			mux.bucketCapabilities[BucketCapabilityCreateAsDeleted] = CapabilityStatusUnsupported
		}

		if cfg.ContainsBucketCapability("subdoc.ReplaceBodyWithXattr") {
			mux.bucketCapabilities[BucketCapabilityReplaceBodyWithXattr] = CapabilityStatusSupported
		} else {
			mux.bucketCapabilities[BucketCapabilityReplaceBodyWithXattr] = CapabilityStatusUnsupported
		}

		if cfg.ContainsBucketCapability("rangeScan") {
			mux.bucketCapabilities[BucketCapabilityRangeScan] = CapabilityStatusSupported
		} else {
			mux.bucketCapabilities[BucketCapabilityRangeScan] = CapabilityStatusUnsupported
		}

		if cfg.ContainsBucketCapability("subdoc.ReplicaRead") {
			mux.bucketCapabilities[BucketCapabilityReplicaRead] = CapabilityStatusSupported
		} else {
			mux.bucketCapabilities[BucketCapabilityReplicaRead] = CapabilityStatusUnsupported
		}

		if cfg.ContainsBucketCapability("nonDedupedHistory") {
			mux.bucketCapabilities[BucketCapabilityNonDedupedHistory] = CapabilityStatusSupported
		} else {
			mux.bucketCapabilities[BucketCapabilityNonDedupedHistory] = CapabilityStatusUnsupported
		}
	}

	return mux
}

func (mux *kvMuxState) RouteConfig() *routeConfig {
	return &mux.routeCfg
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
	var epList []string
	for _, s := range mux.kvServerList {
		epList = append(epList, s.Address)
	}
	return epList
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

func (mux *kvMuxState) HasBucketCapabilityStatus(cap BucketCapability, status CapabilityStatus) bool {
	st, ok := mux.bucketCapabilities[cap]
	if !ok {
		return status == CapabilityStatusUnsupported
	}

	return st == status
}

func (mux *kvMuxState) BucketCapabilityStatus(cap BucketCapability) CapabilityStatus {
	st, ok := mux.bucketCapabilities[cap]
	if !ok {
		return CapabilityStatusUnsupported
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
