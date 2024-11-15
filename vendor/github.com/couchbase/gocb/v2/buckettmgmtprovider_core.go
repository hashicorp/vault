package gocb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type bucketManagementProviderCore struct {
	mgmtProvider mgmtProvider
	tracer       RequestTracer
	meter        *meterWrapper
}

type jsonBucketSettings struct {
	Name        string `json:"name"`
	Controllers struct {
		Flush string `json:"flush"`
	} `json:"controllers"`
	ReplicaIndex bool `json:"replicaIndex"`
	Quota        struct {
		RAM    uint64 `json:"ram"`
		RawRAM uint64 `json:"rawRAM"`
	} `json:"quota"`
	ReplicaNumber                     uint32 `json:"replicaNumber"`
	BucketType                        string `json:"bucketType"`
	ConflictResolutionType            string `json:"conflictResolutionType"`
	EvictionPolicy                    string `json:"evictionPolicy"`
	MaxTTL                            uint32 `json:"maxTTL"`
	CompressionMode                   string `json:"compressionMode"`
	MinimumDurabilityLevel            string `json:"durabilityMinLevel"`
	StorageBackend                    string `json:"storageBackend"`
	HistoryRetentionCollectionDefault *bool  `json:"historyRetentionCollectionDefault"`
	HistoryRetentionBytes             uint64 `json:"historyRetentionBytes"`
	HistoryRetentionSeconds           int    `json:"historyRetentionSeconds"`
}

func (bs *BucketSettings) fromData(data jsonBucketSettings) error {
	bs.Name = data.Name
	bs.FlushEnabled = data.Controllers.Flush != ""
	bs.ReplicaIndexDisabled = !data.ReplicaIndex
	bs.RAMQuotaMB = data.Quota.RawRAM / 1024 / 1024
	bs.NumReplicas = data.ReplicaNumber
	bs.EvictionPolicy = EvictionPolicyType(data.EvictionPolicy)
	bs.MaxTTL = time.Duration(data.MaxTTL) * time.Second
	bs.MaxExpiry = time.Duration(data.MaxTTL) * time.Second
	bs.CompressionMode = CompressionMode(data.CompressionMode)
	bs.MinimumDurabilityLevel = durabilityLevelFromManagementAPI(data.MinimumDurabilityLevel)
	bs.StorageBackend = StorageBackend(data.StorageBackend)
	bs.HistoryRetentionBytes = data.HistoryRetentionBytes
	bs.HistoryRetentionDuration = time.Duration(data.HistoryRetentionSeconds) * time.Second

	if data.HistoryRetentionCollectionDefault != nil {
		if *data.HistoryRetentionCollectionDefault {
			bs.HistoryRetentionCollectionDefault = HistoryRetentionCollectionDefaultEnabled
		} else {
			bs.HistoryRetentionCollectionDefault = HistoryRetentionCollectionDefaultDisabled
		}
	}

	switch data.BucketType {
	case "membase":
		bs.BucketType = CouchbaseBucketType
	case "memcached":
		bs.BucketType = MemcachedBucketType
	case "ephemeral":
		bs.BucketType = EphemeralBucketType
	default:
		return errors.New("unrecognized bucket type string")
	}

	return nil
}

type bucketMgrErrorResp struct {
	Errors map[string]string `json:"errors"`
}

func (bm *bucketManagementProviderCore) tryParseErrorMessage(req *mgmtRequest, resp *mgmtResponse) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logDebugf("Failed to read bucket manager response body: %s", err)
		return nil
	}

	if resp.StatusCode == 404 {
		// If it was a 404 then there's no chance of the response body containing any structure
		// The server will not return 'non existent bucket' after MB-60487, will remove soon as it will not be returned from a released server version
		if strings.Contains(strings.ToLower(string(b)), "resource not found") ||
			strings.Contains(strings.ToLower(string(b)), "non existent bucket") {

			return makeGenericMgmtError(ErrBucketNotFound, req, resp, string(b))
		}

		return makeGenericMgmtError(errors.New(string(b)), req, resp, string(b))
	}

	if err := checkForRateLimitError(resp.StatusCode, string(b)); err != nil {
		return makeGenericMgmtError(err, req, resp, string(b))
	}

	var mgrErr bucketMgrErrorResp
	err = json.Unmarshal(b, &mgrErr)
	if err != nil {
		logDebugf("Failed to unmarshal error body: %s", err)
		return makeGenericMgmtError(errors.New(string(b)), req, resp, string(b))
	}

	var bodyErr error
	var firstErr string
	for _, err := range mgrErr.Errors {
		firstErr = strings.ToLower(err)
		break
	}

	if strings.Contains(firstErr, "bucket with given name already exists") {
		bodyErr = ErrBucketExists
	} else {
		bodyErr = errors.New(firstErr)
	}

	return makeGenericMgmtError(bodyErr, req, resp, string(b))
}

// Flush doesn't use the same body format as anything else...
func (bm *bucketManagementProviderCore) tryParseFlushErrorMessage(req *mgmtRequest, resp *mgmtResponse) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logDebugf("Failed to read bucket manager response body: %s", err)
		return makeMgmtBadStatusError("failed to flush bucket", req, resp)
	}

	if resp.StatusCode == 404 {
		// If it was a 404 then there's no chance of the response body containing any structure
		// The server will not return 'non existent bucket' after MB-60487, will remove soon as it will not be returned from a released server version
		if strings.Contains(strings.ToLower(string(b)), "resource not found") ||
			strings.Contains(strings.ToLower(string(b)), "non existent bucket") {

			return makeGenericMgmtError(ErrBucketNotFound, req, resp, string(b))
		}

		return makeGenericMgmtError(errors.New(string(b)), req, resp, string(b))
	}

	var bodyErrMsgs map[string]string
	err = json.Unmarshal(b, &bodyErrMsgs)
	if err != nil {
		return errors.New(string(b))
	}

	if errMsg, ok := bodyErrMsgs["_"]; ok {
		if strings.Contains(strings.ToLower(errMsg), "flush is disabled") {
			return makeGenericMgmtError(ErrBucketNotFlushable, req, resp, string(b))
		}
	}

	return errors.New(string(b))
}

// GetBucket returns settings for a bucket on the cluster.
func (bm *bucketManagementProviderCore) GetBucket(bucketName string, opts *GetBucketOptions) (*BucketSettings, error) {
	start := time.Now()
	defer bm.meter.ValueRecord(meterValueServiceManagement, "manager_bucket_get_bucket", start)

	path := fmt.Sprintf("/pools/default/buckets/%s", url.PathEscape(bucketName))
	span := createSpan(bm.tracer, opts.ParentSpan, "manager_bucket_get_bucket", "management")
	span.SetAttribute("db.name", bucketName)
	span.SetAttribute("db.operation", "GET "+path)
	defer span.End()

	return bm.get(opts.Context, span.Context(), path, opts.RetryStrategy, opts.Timeout)
}

func (bm *bucketManagementProviderCore) get(ctx context.Context, tracectx RequestSpanContext, path string,
	strategy RetryStrategy, timeout time.Duration) (*BucketSettings, error) {

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          path,
		Method:        "GET",
		IsIdempotent:  true,
		RetryStrategy: strategy,
		UniqueID:      uuid.New().String(),
		Timeout:       timeout,
		parentSpanCtx: tracectx,
	}

	resp, err := bm.mgmtProvider.executeMgmtRequest(ctx, req)
	if err != nil {
		return nil, makeGenericMgmtError(err, &req, resp, "")
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		bktErr := bm.tryParseErrorMessage(&req, resp)
		if bktErr != nil {
			return nil, bktErr
		}

		return nil, makeMgmtBadStatusError("failed to get bucket", &req, resp)
	}

	var bucketData jsonBucketSettings
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&bucketData)
	if err != nil {
		return nil, err
	}

	var settings BucketSettings
	err = settings.fromData(bucketData)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// GetAllBuckets returns a list of all active buckets on the cluster.
func (bm *bucketManagementProviderCore) GetAllBuckets(opts *GetAllBucketsOptions) (map[string]BucketSettings, error) {
	start := time.Now()
	defer bm.meter.ValueRecord(meterValueServiceManagement, "manager_bucket_get_all_buckets", start)

	span := createSpan(bm.tracer, opts.ParentSpan, "manager_bucket_get_all_buckets", "management")
	span.SetAttribute("db.operation", "GET /pools/default/buckets")
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          "/pools/default/buckets",
		Method:        "GET",
		IsIdempotent:  true,
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := bm.mgmtProvider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, makeGenericMgmtError(err, &req, resp, "")
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		bktErr := bm.tryParseErrorMessage(&req, resp)
		if bktErr != nil {
			return nil, bktErr
		}

		return nil, makeMgmtBadStatusError("failed to get all buckets", &req, resp)
	}

	var bucketsData []*jsonBucketSettings
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&bucketsData)
	if err != nil {
		return nil, err
	}

	buckets := make(map[string]BucketSettings, len(bucketsData))
	for _, bucketData := range bucketsData {
		var bucket BucketSettings
		err := bucket.fromData(*bucketData)
		if err != nil {
			return nil, err
		}

		buckets[bucket.Name] = bucket
	}

	return buckets, nil
}

// CreateBucket creates a bucket on the cluster.
func (bm *bucketManagementProviderCore) CreateBucket(settings CreateBucketSettings, opts *CreateBucketOptions) error {
	start := time.Now()
	defer bm.meter.ValueRecord(meterValueServiceManagement, "manager_bucket_create_bucket", start)

	span := createSpan(bm.tracer, opts.ParentSpan, "manager_bucket_create_bucket", "management")
	span.SetAttribute("db.name", settings.Name)
	span.SetAttribute("db.operation", "POST /pools/default/buckets")
	defer span.End()

	posts, err := bm.settingsToPostData(&settings.BucketSettings)
	if err != nil {
		return err
	}

	if settings.ConflictResolutionType != "" {
		posts.Add("conflictResolutionType", string(settings.ConflictResolutionType))
	}

	eSpan := createSpan(bm.tracer, span, "request_encoding", "")
	d := posts.Encode()
	eSpan.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          "/pools/default/buckets",
		Method:        "POST",
		Body:          []byte(d),
		ContentType:   "application/x-www-form-urlencoded",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := bm.mgmtProvider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp, "")
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 202 {
		bktErr := bm.tryParseErrorMessage(&req, resp)
		if bktErr != nil {
			return bktErr
		}

		return makeMgmtBadStatusError("failed to create bucket", &req, resp)
	}

	return nil
}

// UpdateBucket updates a bucket on the cluster.
func (bm *bucketManagementProviderCore) UpdateBucket(settings BucketSettings, opts *UpdateBucketOptions) error {
	start := time.Now()
	defer bm.meter.ValueRecord(meterValueServiceManagement, "manager_bucket_update_bucket", start)

	path := fmt.Sprintf("/pools/default/buckets/%s", url.PathEscape(settings.Name))
	span := createSpan(bm.tracer, opts.ParentSpan, "manager_bucket_update_bucket", "management")
	span.SetAttribute("db.name", settings.Name)
	span.SetAttribute("db.operation", "POST "+path)
	defer span.End()

	posts, err := bm.settingsToPostData(&settings)
	if err != nil {
		return err
	}

	eSpan := createSpan(bm.tracer, span, "request_encoding", "")
	d := posts.Encode()
	eSpan.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          path,
		Method:        "POST",
		Body:          []byte(d),
		ContentType:   "application/x-www-form-urlencoded",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := bm.mgmtProvider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp, "")
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		bktErr := bm.tryParseErrorMessage(&req, resp)
		if bktErr != nil {
			return bktErr
		}

		return makeMgmtBadStatusError("failed to update bucket", &req, resp)
	}

	return nil
}

// DropBucket will delete a bucket from the cluster by name.
func (bm *bucketManagementProviderCore) DropBucket(name string, opts *DropBucketOptions) error {
	start := time.Now()
	defer bm.meter.ValueRecord(meterValueServiceManagement, "manager_bucket_drop_bucket", start)

	path := fmt.Sprintf("/pools/default/buckets/%s", url.PathEscape(name))
	span := createSpan(bm.tracer, opts.ParentSpan, "manager_bucket_drop_bucket", "management")
	span.SetAttribute("db.name", name)
	span.SetAttribute("db.operation", "DELETE "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          path,
		Method:        "DELETE",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := bm.mgmtProvider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp, "")
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		bktErr := bm.tryParseErrorMessage(&req, resp)
		if bktErr != nil {
			return bktErr
		}

		return makeMgmtBadStatusError("failed to drop bucket", &req, resp)
	}

	return nil
}

// FlushBucket will delete all the of the data from a bucket.
// Keep in mind that you must have flushing enabled in the buckets configuration.
func (bm *bucketManagementProviderCore) FlushBucket(name string, opts *FlushBucketOptions) error {
	start := time.Now()
	defer bm.meter.ValueRecord(meterValueServiceManagement, "manager_bucket_flush_bucket", start)

	path := fmt.Sprintf("/pools/default/buckets/%s/controller/doFlush", url.PathEscape(name))
	span := createSpan(bm.tracer, opts.ParentSpan, "manager_bucket_flush_bucket", "management")
	span.SetAttribute("db.name", name)
	span.SetAttribute("db.operation", "POST "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          path,
		Method:        "POST",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := bm.mgmtProvider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp, "")
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		return bm.tryParseFlushErrorMessage(&req, resp)
	}

	return nil
}

func (bm *bucketManagementProviderCore) settingsToPostData(settings *BucketSettings) (url.Values, error) {
	posts := url.Values{}

	if settings.Name == "" {
		return nil, makeInvalidArgumentsError("Name invalid, must be set.")
	}

	if settings.RAMQuotaMB < 100 {
		return nil, makeInvalidArgumentsError("Memory quota invalid, must be greater than 100MB")
	}

	if (settings.MaxTTL > 0 || settings.MaxExpiry > 0) && settings.BucketType == MemcachedBucketType {
		return nil, makeInvalidArgumentsError("maxExpiry is not supported for memcached buckets")
	}

	posts.Add("name", settings.Name)
	// posts.Add("saslPassword", settings.Password)

	if settings.FlushEnabled {
		posts.Add("flushEnabled", "1")
	} else {
		posts.Add("flushEnabled", "0")
	}

	// replicaIndex can't be set at all on ephemeral buckets.
	if settings.BucketType != EphemeralBucketType {
		if settings.ReplicaIndexDisabled {
			posts.Add("replicaIndex", "0")
		} else {
			posts.Add("replicaIndex", "1")
		}
	}

	switch settings.BucketType {
	case CouchbaseBucketType:
		posts.Add("bucketType", string(settings.BucketType))
		posts.Add("replicaNumber", fmt.Sprintf("%d", settings.NumReplicas))
	case MemcachedBucketType:
		posts.Add("bucketType", string(settings.BucketType))
		if settings.NumReplicas > 0 {
			return nil, makeInvalidArgumentsError("replicas cannot be used with memcached buckets")
		}
	case EphemeralBucketType:
		posts.Add("bucketType", string(settings.BucketType))
		posts.Add("replicaNumber", fmt.Sprintf("%d", settings.NumReplicas))
	default:
		return nil, makeInvalidArgumentsError("unrecognized bucket type")
	}

	posts.Add("ramQuotaMB", fmt.Sprintf("%d", settings.RAMQuotaMB))

	if settings.EvictionPolicy != "" {
		switch settings.BucketType {
		case MemcachedBucketType:
			return nil, makeInvalidArgumentsError("eviction policy is not valid for memcached buckets")
		case CouchbaseBucketType:
			if settings.EvictionPolicy == EvictionPolicyTypeNoEviction || settings.EvictionPolicy == EvictionPolicyTypeNotRecentlyUsed {
				return nil, makeInvalidArgumentsError("eviction policy is not valid for couchbase buckets")
			}
		case EphemeralBucketType:
			if settings.EvictionPolicy == EvictionPolicyTypeFull || settings.EvictionPolicy == EvictionPolicyTypeValueOnly {
				return nil, makeInvalidArgumentsError("eviction policy is not valid for ephemeral buckets")
			}
		}
		posts.Add("evictionPolicy", string(settings.EvictionPolicy))
	}

	if settings.MaxTTL > 0 {
		posts.Add("maxTTL", fmt.Sprintf("%d", settings.MaxTTL/time.Second))
	}

	if settings.MaxExpiry > 0 {
		posts.Add("maxTTL", fmt.Sprintf("%d", settings.MaxExpiry/time.Second))
	}

	if settings.CompressionMode != "" {
		posts.Add("compressionMode", string(settings.CompressionMode))
	}

	if settings.MinimumDurabilityLevel > DurabilityLevelNone {
		level, err := settings.MinimumDurabilityLevel.toManagementAPI()
		if err != nil {
			return nil, err
		}
		posts.Add("durabilityMinLevel", level)
	}

	if settings.StorageBackend != "" {
		posts.Add("storageBackend", string(settings.StorageBackend))
	}

	if settings.HistoryRetentionCollectionDefault != HistoryRetentionCollectionDefaultUnset {
		switch settings.HistoryRetentionCollectionDefault {
		case HistoryRetentionCollectionDefaultEnabled:
			posts.Add("historyRetentionCollectionDefault", "true")
		case HistoryRetentionCollectionDefaultDisabled:
			posts.Add("historyRetentionCollectionDefault", "false")
		default:
			return nil, makeInvalidArgumentsError("unrecognized history retention collection default value")
		}
	}
	if settings.HistoryRetentionDuration > 0 {
		posts.Add("historyRetentionSeconds", fmt.Sprintf("%d", settings.HistoryRetentionDuration/time.Second))
	}
	if settings.HistoryRetentionBytes > 0 {
		posts.Add("historyRetentionBytes", fmt.Sprintf("%d", settings.HistoryRetentionBytes))
	}

	return posts, nil
}
