package gocb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// BucketType specifies the kind of bucket.
type BucketType string

const (
	// CouchbaseBucketType indicates a Couchbase bucket type.
	CouchbaseBucketType BucketType = "membase"

	// MemcachedBucketType indicates a Memcached bucket type.
	MemcachedBucketType BucketType = "memcached"

	// EphemeralBucketType indicates an Ephemeral bucket type.
	EphemeralBucketType BucketType = "ephemeral"
)

// ConflictResolutionType specifies the kind of conflict resolution to use for a bucket.
type ConflictResolutionType string

const (
	// ConflictResolutionTypeTimestamp specifies to use timestamp conflict resolution on the bucket.
	ConflictResolutionTypeTimestamp ConflictResolutionType = "lww"

	// ConflictResolutionTypeSequenceNumber specifies to use sequence number conflict resolution on the bucket.
	ConflictResolutionTypeSequenceNumber ConflictResolutionType = "seqno"
)

// EvictionPolicyType specifies the kind of eviction policy to use for a bucket.
type EvictionPolicyType string

const (
	// EvictionPolicyTypeFull specifies to use full eviction for a couchbase bucket.
	EvictionPolicyTypeFull EvictionPolicyType = "fullEviction"

	// EvictionPolicyTypeValueOnly specifies to use value only eviction for a couchbase bucket.
	EvictionPolicyTypeValueOnly EvictionPolicyType = "valueOnly"

	// EvictionPolicyTypeNotRecentlyUsed specifies to use not recently used (nru) eviction for an ephemeral bucket.
	// UNCOMMITTED: This API may change in the future.
	EvictionPolicyTypeNotRecentlyUsed EvictionPolicyType = "nruEviction"

	// EvictionPolicyTypeNRU specifies to use no eviction for an ephemeral bucket.
	// UNCOMMITTED: This API may change in the future.
	EvictionPolicyTypeNoEviction EvictionPolicyType = "noEviction"
)

// CompressionMode specifies the kind of compression to use for a bucket.
type CompressionMode string

const (
	// CompressionModeOff specifies to use no compression for a bucket.
	CompressionModeOff CompressionMode = "off"

	// CompressionModePassive specifies to use passive compression for a bucket.
	CompressionModePassive CompressionMode = "passive"

	// CompressionModeActive specifies to use active compression for a bucket.
	CompressionModeActive CompressionMode = "active"
)

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
	ReplicaNumber          uint32 `json:"replicaNumber"`
	BucketType             string `json:"bucketType"`
	ConflictResolutionType string `json:"conflictResolutionType"`
	EvictionPolicy         string `json:"evictionPolicy"`
	MaxTTL                 uint32 `json:"maxTTL"`
	CompressionMode        string `json:"compressionMode"`
}

// BucketSettings holds information about the settings for a bucket.
type BucketSettings struct {
	Name                 string
	FlushEnabled         bool
	ReplicaIndexDisabled bool // inverted so that zero value matches server default.
	RAMQuotaMB           uint64
	NumReplicas          uint32     // NOTE: If not set this will set 0 replicas.
	BucketType           BucketType // Defaults to CouchbaseBucketType.
	EvictionPolicy       EvictionPolicyType
	MaxTTL               time.Duration
	CompressionMode      CompressionMode
}

func (bs *BucketSettings) fromData(data jsonBucketSettings) error {
	bs.Name = data.Name
	bs.FlushEnabled = data.Controllers.Flush != ""
	bs.ReplicaIndexDisabled = !data.ReplicaIndex
	bs.RAMQuotaMB = data.Quota.RawRAM / 1024 / 1024
	bs.NumReplicas = data.ReplicaNumber
	bs.EvictionPolicy = EvictionPolicyType(data.EvictionPolicy)
	bs.MaxTTL = time.Duration(data.MaxTTL) * time.Second
	bs.CompressionMode = CompressionMode(data.CompressionMode)

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

func (bm *BucketManager) tryParseErrorMessage(req *mgmtRequest, resp *mgmtResponse) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logDebugf("Failed to read bucket manager response body: %s", err)
		return nil
	}

	if resp.StatusCode == 404 {
		// If it was a 404 then there's no chance of the response body containing any structure
		if strings.Contains(strings.ToLower(string(b)), "resource not found") {
			return makeGenericMgmtError(ErrBucketNotFound, req, resp)
		}

		return makeGenericMgmtError(errors.New(string(b)), req, resp)
	}

	var mgrErr bucketMgrErrorResp
	err = json.Unmarshal(b, &mgrErr)
	if err != nil {
		logDebugf("Failed to unmarshal error body: %s", err)
		return makeGenericMgmtError(errors.New(string(b)), req, resp)
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

	return makeGenericMgmtError(bodyErr, req, resp)
}

// Flush doesn't use the same body format as anything else...
func (bm *BucketManager) tryParseFlushErrorMessage(req *mgmtRequest, resp *mgmtResponse) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logDebugf("Failed to read bucket manager response body: %s", err)
		return makeMgmtBadStatusError("failed to flush bucket", req, resp)
	}

	var bodyErrMsgs map[string]string
	err = json.Unmarshal(b, &bodyErrMsgs)
	if err != nil {
		return errors.New(string(b))
	}

	if errMsg, ok := bodyErrMsgs["_"]; ok {
		if strings.Contains(strings.ToLower(errMsg), "flush is disabled") {
			return ErrBucketNotFlushable
		}
	}

	return errors.New(string(b))
}

// BucketManager provides methods for performing bucket management operations.
// See BucketManager for methods that allow creating and removing buckets themselves.
type BucketManager struct {
	provider mgmtProvider
	tracer   requestTracer
}

// GetBucketOptions is the set of options available to the bucket manager GetBucket operation.
type GetBucketOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// GetBucket returns settings for a bucket on the cluster.
func (bm *BucketManager) GetBucket(bucketName string, opts *GetBucketOptions) (*BucketSettings, error) {
	if opts == nil {
		opts = &GetBucketOptions{}
	}

	span := bm.tracer.StartSpan("GetBucket", nil).
		SetTag("couchbase.service", "mgmt")
	defer span.Finish()

	return bm.get(span.Context(), bucketName, opts.RetryStrategy, opts.Timeout)
}

func (bm *BucketManager) get(tracectx requestSpanContext, bucketName string,
	strategy RetryStrategy, timeout time.Duration) (*BucketSettings, error) {

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          fmt.Sprintf("/pools/default/buckets/%s", bucketName),
		Method:        "GET",
		IsIdempotent:  true,
		RetryStrategy: strategy,
		UniqueID:      uuid.New().String(),
		Timeout:       timeout,
		parentSpan:    tracectx,
	}

	resp, err := bm.provider.executeMgmtRequest(req)
	if err != nil {
		return nil, makeGenericMgmtError(err, &req, resp)
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

// GetAllBucketsOptions is the set of options available to the bucket manager GetAll operation.
type GetAllBucketsOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// GetAllBuckets returns a list of all active buckets on the cluster.
func (bm *BucketManager) GetAllBuckets(opts *GetAllBucketsOptions) (map[string]BucketSettings, error) {
	if opts == nil {
		opts = &GetAllBucketsOptions{}
	}

	span := bm.tracer.StartSpan("GetAllBuckets", nil).
		SetTag("couchbase.service", "mgmt")
	defer span.Finish()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          "/pools/default/buckets",
		Method:        "GET",
		IsIdempotent:  true,
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpan:    span.Context(),
	}

	resp, err := bm.provider.executeMgmtRequest(req)
	if err != nil {
		return nil, makeGenericMgmtError(err, &req, resp)
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

// CreateBucketSettings are the settings available when creating a bucket.
type CreateBucketSettings struct {
	BucketSettings
	ConflictResolutionType ConflictResolutionType
}

// CreateBucketOptions is the set of options available to the bucket manager CreateBucket operation.
type CreateBucketOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// CreateBucket creates a bucket on the cluster.
func (bm *BucketManager) CreateBucket(settings CreateBucketSettings, opts *CreateBucketOptions) error {
	if opts == nil {
		opts = &CreateBucketOptions{}
	}

	span := bm.tracer.StartSpan("CreateBucket", nil).
		SetTag("couchbase.service", "mgmt")
	defer span.Finish()

	posts, err := bm.settingsToPostData(&settings.BucketSettings)
	if err != nil {
		return err
	}

	if settings.ConflictResolutionType != "" {
		posts.Add("conflictResolutionType", string(settings.ConflictResolutionType))
	}

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          "/pools/default/buckets",
		Method:        "POST",
		Body:          []byte(posts.Encode()),
		ContentType:   "application/x-www-form-urlencoded",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpan:    span.Context(),
	}

	resp, err := bm.provider.executeMgmtRequest(req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp)
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

// UpdateBucketOptions is the set of options available to the bucket manager UpdateBucket operation.
type UpdateBucketOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// UpdateBucket updates a bucket on the cluster.
func (bm *BucketManager) UpdateBucket(settings BucketSettings, opts *UpdateBucketOptions) error {
	if opts == nil {
		opts = &UpdateBucketOptions{}
	}

	span := bm.tracer.StartSpan("UpdateBucket", nil).
		SetTag("couchbase.service", "mgmt")
	defer span.Finish()

	posts, err := bm.settingsToPostData(&settings)
	if err != nil {
		return err
	}

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          fmt.Sprintf("/pools/default/buckets/%s", settings.Name),
		Method:        "POST",
		Body:          []byte(posts.Encode()),
		ContentType:   "application/x-www-form-urlencoded",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpan:    span.Context(),
	}

	resp, err := bm.provider.executeMgmtRequest(req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp)
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

// DropBucketOptions is the set of options available to the bucket manager DropBucket operation.
type DropBucketOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// DropBucket will delete a bucket from the cluster by name.
func (bm *BucketManager) DropBucket(name string, opts *DropBucketOptions) error {
	if opts == nil {
		opts = &DropBucketOptions{}
	}

	span := bm.tracer.StartSpan("DropBucket", nil).
		SetTag("couchbase.service", "mgmt")
	defer span.Finish()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          fmt.Sprintf("/pools/default/buckets/%s", name),
		Method:        "DELETE",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpan:    span.Context(),
	}

	resp, err := bm.provider.executeMgmtRequest(req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp)
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

// FlushBucketOptions is the set of options available to the bucket manager FlushBucket operation.
type FlushBucketOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// FlushBucket will delete all the of the data from a bucket.
// Keep in mind that you must have flushing enabled in the buckets configuration.
func (bm *BucketManager) FlushBucket(name string, opts *FlushBucketOptions) error {
	if opts == nil {
		opts = &FlushBucketOptions{}
	}

	span := bm.tracer.StartSpan("FlushBucket", nil).
		SetTag("couchbase.service", "mgmt")
	defer span.Finish()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          fmt.Sprintf("/pools/default/buckets/%s/controller/doFlush", name),
		Method:        "POST",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpan:    span.Context(),
	}

	resp, err := bm.provider.executeMgmtRequest(req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		return bm.tryParseFlushErrorMessage(&req, resp)
	}

	return nil
}

func (bm *BucketManager) settingsToPostData(settings *BucketSettings) (url.Values, error) {
	posts := url.Values{}

	if settings.Name == "" {
		return nil, makeInvalidArgumentsError("Name invalid, must be set.")
	}

	if settings.RAMQuotaMB < 100 {
		return nil, makeInvalidArgumentsError("Memory quota invalid, must be greater than 100MB")
	}

	if settings.MaxTTL > 0 && settings.BucketType == MemcachedBucketType {
		return nil, makeInvalidArgumentsError("maxTTL is not supported for memcached buckets")
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
		return nil, makeInvalidArgumentsError("Unrecognized bucket type")
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

	if settings.CompressionMode != "" {
		posts.Add("compressionMode", string(settings.CompressionMode))
	}

	return posts, nil
}
