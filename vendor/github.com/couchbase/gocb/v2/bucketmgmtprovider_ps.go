package gocb

import (
	"fmt"
	"time"

	"github.com/couchbase/goprotostellar/genproto/admin_bucket_v1"
	"github.com/couchbase/goprotostellar/genproto/kv_v1"
)

type bucketManagementProviderPs struct {
	provider admin_bucket_v1.BucketAdminServiceClient

	managerProvider *psOpManagerProvider
}

func (bm bucketManagementProviderPs) newOpManager(parentSpan RequestSpan, opName string, attribs map[string]interface{}) *psOpManagerDefault {
	return bm.managerProvider.NewManager(parentSpan, opName, attribs)
}

func (bm bucketManagementProviderPs) GetBucket(bucketName string, opts *GetBucketOptions) (*BucketSettings, error) {
	manager := bm.newOpManager(opts.ParentSpan, "manager_bucket_get_bucket", map[string]interface{}{
		"db.name":      bucketName,
		"db.operation": "ListBuckets",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(true)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return nil, err
	}

	req := &admin_bucket_v1.ListBucketsRequest{}
	resp, err := wrapPSOp(manager, req, bm.provider.ListBuckets)
	if err != nil {
		return nil, err
	}

	for _, source := range resp.Buckets {
		if source.BucketName == bucketName {
			bucket, err := bm.psBucketToBucket(source)
			if err != nil {
				return nil, makeGenericError(err, nil)
			}

			return bucket, nil
		}
	}

	return nil, makeGenericError(ErrBucketNotFound, nil)
}

func (bm bucketManagementProviderPs) GetAllBuckets(opts *GetAllBucketsOptions) (map[string]BucketSettings, error) {
	manager := bm.newOpManager(opts.ParentSpan, "manager_bucket_get_all_buckets", map[string]interface{}{
		"db.operation": "ListBuckets",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(true)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return nil, err
	}

	req := &admin_bucket_v1.ListBucketsRequest{}
	resp, err := wrapPSOp(manager, req, bm.provider.ListBuckets)
	if err != nil {
		return nil, err
	}

	buckets := make(map[string]BucketSettings)
	for _, source := range resp.Buckets {
		bucket, err := bm.psBucketToBucket(source)
		if err != nil {
			return nil, makeGenericError(err, nil)
		}

		buckets[bucket.Name] = *bucket
	}

	return buckets, nil
}

func (bm bucketManagementProviderPs) CreateBucket(settings CreateBucketSettings, opts *CreateBucketOptions) error {
	manager := bm.newOpManager(opts.ParentSpan, "manager_bucket_create_bucket", map[string]interface{}{
		"db.name":      settings.Name,
		"db.operation": "CreateBucket",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return err
	}

	req, err := bm.settingsToCreateReq(settings)
	if err != nil {
		return makeGenericError(err, nil)
	}

	_, err = wrapPSOp(manager, req, bm.provider.CreateBucket)
	if err != nil {
		return err
	}

	return nil
}

func (bm bucketManagementProviderPs) UpdateBucket(settings BucketSettings, opts *UpdateBucketOptions) error {
	manager := bm.newOpManager(opts.ParentSpan, "manager_bucket_update_bucket", map[string]interface{}{
		"db.name":      settings.Name,
		"db.operation": "UpdateBucket",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return err
	}

	req, err := bm.settingsToUpdateReq(settings)
	if err != nil {
		return makeGenericError(err, nil)
	}

	_, err = wrapPSOp(manager, req, bm.provider.UpdateBucket)
	if err != nil {
		return err
	}

	return nil
}

func (bm bucketManagementProviderPs) DropBucket(name string, opts *DropBucketOptions) error {
	manager := bm.newOpManager(opts.ParentSpan, "manager_bucket_drop_bucket", map[string]interface{}{
		"db.name":      name,
		"db.operation": "DeleteBucket",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return err
	}

	req := &admin_bucket_v1.DeleteBucketRequest{BucketName: name}

	_, err := wrapPSOp(manager, req, bm.provider.DeleteBucket)
	if err != nil {
		return err
	}

	return nil
}

func (bm bucketManagementProviderPs) FlushBucket(name string, opts *FlushBucketOptions) error {
	manager := bm.newOpManager(opts.ParentSpan, "manager_bucket_flush_bucket", map[string]interface{}{
		"db.name":      name,
		"db.operation": "FlushBucket",
	})
	defer manager.Finish(false)

	manager.SetContext(opts.Context)
	manager.SetIsIdempotent(false)
	manager.SetRetryStrategy(opts.RetryStrategy)
	manager.SetTimeout(opts.Timeout)

	if err := manager.CheckReadyForOp(); err != nil {
		return err
	}

	req := &admin_bucket_v1.FlushBucketRequest{BucketName: name}

	_, err := wrapPSOp(manager, req, bm.provider.FlushBucket)
	if err != nil {
		return err
	}

	return nil
}

func (bm bucketManagementProviderPs) psBucketToBucket(source *admin_bucket_v1.ListBucketsResponse_Bucket) (*BucketSettings, error) {
	bucket := &BucketSettings{
		Name:                 source.BucketName,
		FlushEnabled:         source.FlushEnabled,
		ReplicaIndexDisabled: !source.ReplicaIndexes,
		RAMQuotaMB:           source.RamQuotaMb,
		NumReplicas:          source.NumReplicas,
		MaxExpiry:            time.Duration(source.MaxExpirySecs) * time.Second,
	}

	// memcached buckets are not supported by couchbase2 so we shouldn't be receiving buckets of that type.
	switch source.BucketType {
	case admin_bucket_v1.BucketType_BUCKET_TYPE_COUCHBASE:
		bucket.BucketType = CouchbaseBucketType
	case admin_bucket_v1.BucketType_BUCKET_TYPE_EPHEMERAL:
		bucket.BucketType = EphemeralBucketType
	default:
		return nil, fmt.Errorf("unrecognized bucket type %s", source.BucketType)
	}

	switch source.EvictionMode {
	case admin_bucket_v1.EvictionMode_EVICTION_MODE_FULL:
		bucket.EvictionPolicy = EvictionPolicyTypeFull
	case admin_bucket_v1.EvictionMode_EVICTION_MODE_NOT_RECENTLY_USED:
		bucket.EvictionPolicy = EvictionPolicyTypeNotRecentlyUsed
	case admin_bucket_v1.EvictionMode_EVICTION_MODE_VALUE_ONLY:
		bucket.EvictionPolicy = EvictionPolicyTypeValueOnly
	case admin_bucket_v1.EvictionMode_EVICTION_MODE_NONE:
		bucket.EvictionPolicy = EvictionPolicyTypeNoEviction
	}

	switch source.CompressionMode {
	case admin_bucket_v1.CompressionMode_COMPRESSION_MODE_OFF:
		bucket.CompressionMode = CompressionModeOff
	case admin_bucket_v1.CompressionMode_COMPRESSION_MODE_PASSIVE:
		bucket.CompressionMode = CompressionModePassive
	case admin_bucket_v1.CompressionMode_COMPRESSION_MODE_ACTIVE:
		bucket.CompressionMode = CompressionModeActive
	}

	if source.MinimumDurabilityLevel == nil {
		bucket.MinimumDurabilityLevel = DurabilityLevelNone
	} else {
		switch source.GetMinimumDurabilityLevel() {
		case kv_v1.DurabilityLevel_DURABILITY_LEVEL_MAJORITY:
			bucket.MinimumDurabilityLevel = DurabilityLevelMajority
		case kv_v1.DurabilityLevel_DURABILITY_LEVEL_MAJORITY_AND_PERSIST_TO_ACTIVE:
			bucket.MinimumDurabilityLevel = DurabilityLevelMajorityAndPersistOnMaster
		case kv_v1.DurabilityLevel_DURABILITY_LEVEL_PERSIST_TO_MAJORITY:
			bucket.MinimumDurabilityLevel = DurabilityLevelPersistToMajority
		}
	}

	if source.StorageBackend != nil {
		switch *source.StorageBackend {
		case admin_bucket_v1.StorageBackend_STORAGE_BACKEND_COUCHSTORE:
			bucket.StorageBackend = StorageBackendCouchstore
		case admin_bucket_v1.StorageBackend_STORAGE_BACKEND_MAGMA:
			bucket.StorageBackend = StorageBackendMagma
		}
	}

	if source.HistoryRetentionCollectionDefault != nil {
		if *source.HistoryRetentionCollectionDefault {
			bucket.HistoryRetentionCollectionDefault = HistoryRetentionCollectionDefaultEnabled
		} else {
			bucket.HistoryRetentionCollectionDefault = HistoryRetentionCollectionDefaultDisabled
		}
	}

	if source.HistoryRetentionDurationSecs != nil && *source.HistoryRetentionDurationSecs > 0 {
		bucket.HistoryRetentionDuration = time.Duration(*source.HistoryRetentionDurationSecs) * time.Second
	}

	if source.HistoryRetentionBytes != nil && *source.HistoryRetentionBytes > 0 {
		bucket.HistoryRetentionBytes = *source.HistoryRetentionBytes
	}

	return bucket, nil
}

func (bm *bucketManagementProviderPs) settingsToCreateReq(settings CreateBucketSettings) (*admin_bucket_v1.CreateBucketRequest, error) {
	request := &admin_bucket_v1.CreateBucketRequest{}

	err := bm.validateSettings(settings.BucketSettings)
	if err != nil {
		return nil, err
	}

	request.BucketName = settings.Name
	request.NumReplicas = &settings.NumReplicas
	request.RamQuotaMb = &settings.RAMQuotaMB

	if settings.FlushEnabled {
		request.FlushEnabled = &settings.FlushEnabled
	}

	if settings.ReplicaIndexDisabled {
		replicasEnabled := false
		request.ReplicaIndexes = &replicasEnabled
	}

	request.BucketType, err = bm.bucketTypeToPS(settings.BucketType)
	if err != nil {
		return nil, err
	}

	if settings.EvictionPolicy != "" {
		request.EvictionMode, err = bm.evictionPolicyToPS(settings.EvictionPolicy, settings.BucketType)
		if err != nil {
			return nil, err
		}
	}

	if settings.MaxTTL > 0 {
		expiry := uint32(settings.MaxTTL.Seconds())
		request.MaxExpirySecs = &expiry
	}
	if settings.MaxExpiry > 0 {
		expiry := uint32(settings.MaxExpiry.Seconds())
		request.MaxExpirySecs = &expiry
	}

	if settings.CompressionMode != "" {
		request.CompressionMode, err = bm.compressionModeToPS(settings.CompressionMode)
		if err != nil {
			return nil, err
		}
	}

	if settings.MinimumDurabilityLevel > DurabilityLevelNone {
		request.MinimumDurabilityLevel, err = bm.durabilityLevelToPS(settings.MinimumDurabilityLevel)
		if err != nil {
			return nil, err
		}
	}

	if settings.StorageBackend != "" {
		request.StorageBackend, err = bm.storageBackendToPS(settings.StorageBackend)
		if err != nil {
			return nil, err
		}
	}

	if settings.ConflictResolutionType != "" {
		var conflictRes admin_bucket_v1.ConflictResolutionType
		switch settings.ConflictResolutionType {
		case ConflictResolutionTypeTimestamp:
			conflictRes = admin_bucket_v1.ConflictResolutionType_CONFLICT_RESOLUTION_TYPE_TIMESTAMP
		case ConflictResolutionTypeSequenceNumber:
			conflictRes = admin_bucket_v1.ConflictResolutionType_CONFLICT_RESOLUTION_TYPE_SEQUENCE_NUMBER
		case ConflictResolutionTypeCustom:
			conflictRes = admin_bucket_v1.ConflictResolutionType_CONFLICT_RESOLUTION_TYPE_CUSTOM
		default:
			return nil, makeInvalidArgumentsError(fmt.Sprintf("unrecognized conflict resolution type %s", settings.ConflictResolutionType))
		}
		request.ConflictResolutionType = &conflictRes
	}

	if settings.HistoryRetentionCollectionDefault != HistoryRetentionCollectionDefaultUnset {
		var val bool
		if settings.HistoryRetentionCollectionDefault == HistoryRetentionCollectionDefaultEnabled {
			val = true
		} else if settings.HistoryRetentionCollectionDefault == HistoryRetentionCollectionDefaultDisabled {
			val = false
		} else {
			return nil, makeInvalidArgumentsError("unrecognized history retention collection default value")
		}
		request.HistoryRetentionCollectionDefault = &val
	}

	if settings.HistoryRetentionDuration > 0 {
		duration := uint32(settings.HistoryRetentionDuration / time.Second)
		request.HistoryRetentionDurationSecs = &duration
	}
	if settings.HistoryRetentionBytes > 0 {
		bytes := settings.HistoryRetentionBytes
		request.HistoryRetentionBytes = &bytes
	}

	return request, nil
}

func (bm *bucketManagementProviderPs) settingsToUpdateReq(settings BucketSettings) (*admin_bucket_v1.UpdateBucketRequest, error) {
	request := &admin_bucket_v1.UpdateBucketRequest{}

	err := bm.validateSettings(settings)
	if err != nil {
		return nil, err
	}

	request.BucketName = settings.Name
	request.NumReplicas = &settings.NumReplicas
	request.RamQuotaMb = &settings.RAMQuotaMB

	if settings.FlushEnabled {
		request.FlushEnabled = &settings.FlushEnabled
	}

	if settings.EvictionPolicy != "" {
		request.EvictionMode, err = bm.evictionPolicyToPS(settings.EvictionPolicy, settings.BucketType)
		if err != nil {
			return nil, err
		}
	}

	if settings.MaxTTL > 0 {
		expiry := uint32(settings.MaxTTL.Seconds())
		request.MaxExpirySecs = &expiry
	}
	if settings.MaxExpiry > 0 {
		expiry := uint32(settings.MaxExpiry.Seconds())
		request.MaxExpirySecs = &expiry
	}

	if settings.CompressionMode != "" {
		request.CompressionMode, err = bm.compressionModeToPS(settings.CompressionMode)
		if err != nil {
			return nil, err
		}
	}

	if settings.MinimumDurabilityLevel > DurabilityLevelNone {
		request.MinimumDurabilityLevel, err = bm.durabilityLevelToPS(settings.MinimumDurabilityLevel)
		if err != nil {
			return nil, err
		}
	}

	if settings.HistoryRetentionCollectionDefault != HistoryRetentionCollectionDefaultUnset {
		var val bool
		if settings.HistoryRetentionCollectionDefault == HistoryRetentionCollectionDefaultEnabled {
			val = true
		} else if settings.HistoryRetentionCollectionDefault == HistoryRetentionCollectionDefaultDisabled {
			val = false
		} else {
			return nil, makeInvalidArgumentsError("unrecognized history retention collection default value")
		}
		request.HistoryRetentionCollectionDefault = &val
	}

	if settings.HistoryRetentionDuration > 0 {
		duration := uint32(settings.HistoryRetentionDuration / time.Second)
		request.HistoryRetentionDurationSecs = &duration
	}

	if settings.HistoryRetentionBytes > 0 {
		bytes := settings.HistoryRetentionBytes
		request.HistoryRetentionBytes = &bytes
	}

	return request, nil
}

func (bm *bucketManagementProviderPs) validateSettings(settings BucketSettings) error {
	if settings.Name == "" {
		return makeInvalidArgumentsError("Name invalid, must be set.")
	}
	if settings.RAMQuotaMB < 100 {
		return makeInvalidArgumentsError("Memory quota invalid, must be greater than 100MB")
	}
	if (settings.MaxTTL > 0 || settings.MaxExpiry > 0) && settings.BucketType == MemcachedBucketType {
		return makeInvalidArgumentsError("maxExpiry is not supported for memcached buckets")
	}
	if settings.BucketType == MemcachedBucketType && settings.NumReplicas > 0 {
		return makeInvalidArgumentsError("replicas cannot be used with memcached buckets")
	}

	return nil
}

func (bm *bucketManagementProviderPs) bucketTypeToPS(bucketType BucketType) (admin_bucket_v1.BucketType, error) {
	switch bucketType {
	case CouchbaseBucketType:
		return admin_bucket_v1.BucketType_BUCKET_TYPE_COUCHBASE, nil
	case MemcachedBucketType:
		return 0, makeInvalidArgumentsError("memcached bucket type is not supported by the couchbase2 protocol")
	case EphemeralBucketType:
		return admin_bucket_v1.BucketType_BUCKET_TYPE_EPHEMERAL, nil
	default:
		return 0, makeInvalidArgumentsError(fmt.Sprintf("unrecognized bucket type %s", bucketType))
	}
}

func (bm *bucketManagementProviderPs) evictionPolicyToPS(evictionPolicy EvictionPolicyType,
	bucketType BucketType) (*admin_bucket_v1.EvictionMode, error) {
	var policy admin_bucket_v1.EvictionMode
	switch bucketType {
	case MemcachedBucketType:
		return nil, makeInvalidArgumentsError("eviction policy is not valid for memcached buckets")
	case CouchbaseBucketType:
		switch evictionPolicy {
		case EvictionPolicyTypeNoEviction:
			return nil, makeInvalidArgumentsError("eviction policy is not valid for couchbase buckets")
		case EvictionPolicyTypeNotRecentlyUsed:
			return nil, makeInvalidArgumentsError("eviction policy is not valid for couchbase buckets")
		case EvictionPolicyTypeValueOnly:
			policy = admin_bucket_v1.EvictionMode_EVICTION_MODE_VALUE_ONLY
		case EvictionPolicyTypeFull:
			policy = admin_bucket_v1.EvictionMode_EVICTION_MODE_FULL
		default:
			return nil, makeInvalidArgumentsError(fmt.Sprintf("unrecognized eviction policy %s", evictionPolicy))

		}
	case EphemeralBucketType:
		switch evictionPolicy {
		case EvictionPolicyTypeNoEviction:
			policy = admin_bucket_v1.EvictionMode_EVICTION_MODE_NONE
		case EvictionPolicyTypeNotRecentlyUsed:
			policy = admin_bucket_v1.EvictionMode_EVICTION_MODE_NOT_RECENTLY_USED
		case EvictionPolicyTypeValueOnly:
			return nil, makeInvalidArgumentsError("eviction policy is not valid for ephemeral buckets")
		case EvictionPolicyTypeFull:
			return nil, makeInvalidArgumentsError("eviction policy is not valid for ephemeral buckets")
		default:
			return nil, makeInvalidArgumentsError(fmt.Sprintf("unrecognized eviction policy %s", evictionPolicy))
		}
	}
	return &policy, nil
}

func (bm *bucketManagementProviderPs) compressionModeToPS(mode CompressionMode) (*admin_bucket_v1.CompressionMode, error) {
	var compressionMode admin_bucket_v1.CompressionMode
	switch mode {
	case CompressionModeOff:
		compressionMode = admin_bucket_v1.CompressionMode_COMPRESSION_MODE_OFF
	case CompressionModePassive:
		compressionMode = admin_bucket_v1.CompressionMode_COMPRESSION_MODE_PASSIVE
	case CompressionModeActive:
		compressionMode = admin_bucket_v1.CompressionMode_COMPRESSION_MODE_ACTIVE
	default:
		return nil, makeInvalidArgumentsError(fmt.Sprintf("unrecognized compression mode %s", compressionMode))
	}
	return &compressionMode, nil
}

func (bm *bucketManagementProviderPs) durabilityLevelToPS(level DurabilityLevel) (*kv_v1.DurabilityLevel, error) {
	var duraLevel kv_v1.DurabilityLevel
	switch level {
	case DurabilityLevelMajority:
		duraLevel = kv_v1.DurabilityLevel_DURABILITY_LEVEL_MAJORITY
	case DurabilityLevelMajorityAndPersistOnMaster:
		duraLevel = kv_v1.DurabilityLevel_DURABILITY_LEVEL_MAJORITY_AND_PERSIST_TO_ACTIVE
	case DurabilityLevelPersistToMajority:
		duraLevel = kv_v1.DurabilityLevel_DURABILITY_LEVEL_PERSIST_TO_MAJORITY
	default:
		return nil, makeInvalidArgumentsError(fmt.Sprintf("unrecognized durability level %d", level))
	}

	return &duraLevel, nil
}

func (bm *bucketManagementProviderPs) storageBackendToPS(storage StorageBackend) (*admin_bucket_v1.StorageBackend, error) {
	var backend admin_bucket_v1.StorageBackend
	switch storage {
	case StorageBackendCouchstore:
		backend = admin_bucket_v1.StorageBackend_STORAGE_BACKEND_COUCHSTORE
	case StorageBackendMagma:
		backend = admin_bucket_v1.StorageBackend_STORAGE_BACKEND_MAGMA
	default:
		return nil, makeInvalidArgumentsError(fmt.Sprintf("unrecognized storage backend %s", storage))
	}

	return &backend, nil
}
