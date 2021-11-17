package gocb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

// EventingFunctionManager provides methods for performing eventing function management operations.
// This manager is designed to work only against Couchbase Server 7.0+, it might work against earlier server
// versions but that is not tested and is not supported.
// Volatile: This API is subject to change at any time.
type EventingFunctionManager struct {
	mgmtProvider mgmtProvider
	tracer       RequestTracer
	meter        *meterWrapper
}

func (efm *EventingFunctionManager) doMgmtRequest(ctx context.Context, req mgmtRequest) (*mgmtResponse, error) {
	resp, err := efm.mgmtProvider.executeMgmtRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (efm *EventingFunctionManager) tryParseErrorMessage(req *mgmtRequest, resp *mgmtResponse) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logDebugf("Failed to read eventing function response body: %s", err)
		return nil
	}

	var baseErr error
	strBody := string(b)
	if strings.Contains(strBody, "ERR_APP_NOT_FOUND_TS") {
		baseErr = ErrEventingFunctionNotFound
	} else if strings.Contains(strBody, "ERR_APP_NOT_DEPLOYED") {
		baseErr = ErrEventingFunctionNotDeployed
	} else if strings.Contains(strBody, "ERR_HANDLER_COMPILATION") {
		baseErr = ErrEventingFunctionCompilationFailure
	} else if strings.Contains(strBody, "ERR_SRC_MB_SAME") {
		baseErr = ErrEventingFunctionIdenticalKeyspace
	} else if strings.Contains(strBody, "ERR_APP_NOT_BOOTSTRAPPED") {
		baseErr = ErrEventingFunctionNotBootstrapped
	} else if strings.Contains(strBody, "ERR_APP_NOT_UNDEPLOYED") {
		baseErr = ErrEventingFunctionDeployed
	} else if strings.Contains(strBody, "ERR_COLLECTION_MISSING") {
		baseErr = ErrCollectionNotFound
	} else if strings.Contains(strBody, "ERR_BUCKET_MISSING") {
		baseErr = ErrBucketNotFound
	} else {
		baseErr = errors.New(string(b))
	}

	return makeGenericMgmtError(baseErr, req, resp)
}

type jsonEventingFunction struct {
	Name               string                               `json:"appname"`
	Code               string                               `json:"appcode"`
	Version            string                               `json:"version"`
	EnforceSchema      bool                                 `json:"enforce_schema,omitempty"`
	HandlerUUID        int                                  `json:"handleruuid,omitempty"`
	FunctionInstanceID string                               `json:"function_instance_id,omitempty"`
	Settings           jsonEventingFunctionSettings         `json:"settings"`
	DeploymentConfig   jsonEventingFunctionDeploymentConfig `json:"depcfg"`
}

type jsonEventingFunctionSettings struct {
	CPPWorkerThreadCount   int                                   `json:"cpp_worker_thread_count,omitempty"`
	DCPStreamBoundary      EventingFunctionDCPBoundary           `json:"dcp_stream_boundary,omitempty"`
	Description            string                                `json:"description,omitempty"`
	DeploymentStatus       EventingFunctionDeploymentStatus      `json:"deployment_status"`
	ProcessingStatus       EventingFunctionProcessingStatus      `json:"processing_status"`
	LanguageCompatibility  EventingFunctionLanguageCompatibility `json:"language_compatibility,omitempty"`
	LogLevel               EventingFunctionLogLevel              `json:"log_level,omitempty"`
	ExecutionTimeout       int                                   `json:"execution_timeout,omitempty"`
	LCBInstCapacity        int                                   `json:"lcb_inst_capacity,omitempty"`
	LCBRetryCount          int                                   `json:"lcb_retry_count,omitempty"`
	LCBTimeout             int                                   `json:"lcb_timeout,omitempty"`
	QueryConsistency       QueryScanConsistency                  `json:"n1ql_consistency,omitempty"`
	NumTimerPartitions     int                                   `json:"num_timer_partitions,omitempty"`
	SockBatchSize          int                                   `json:"sock_batch_size,omitempty"`
	TickDuration           int                                   `json:"tick_duration,omitempty"`
	TimerContextSize       int                                   `json:"timer_context_size,omitempty"`
	UserPrefix             string                                `json:"user_prefix,omitempty"`
	BucketCacheSize        int                                   `json:"bucket_cache_size,omitempty"`
	BucketCacheAge         int                                   `json:"bucket_cache_age,omitempty"`
	CurlMaxAllowedRespSize int                                   `json:"curl_max_allowed_resp_size,omitempty"`
	QueryPrepareAll        bool                                  `json:"n1ql_prepare_all,omitempty"`
	WorkerCount            int                                   `json:"worker_count,omitempty"`
	HandlerHeaders         []string                              `json:"handler_headers,omitempty"`
	HandlerFooters         []string                              `json:"handler_footers,omitempty"`
	EnableAppLogRotation   bool                                  `json:"enable_applog_rotation,omitempty"`
	AppLogDir              string                                `json:"app_log_dir,omitempty"`
	AppLogMaxSize          int                                   `json:"app_log_max_size,omitempty"`
	AppLogMaxFiles         int                                   `json:"app_log_max_files,omitempty"`
	CheckpointInterval     int                                   `json:"checkpoint_interval,omitempty"`
}

type jsonEventingFunctionDeploymentConfig struct {
	MetadataBucket     string                                `json:"metadata_bucket"`
	MetadataScope      string                                `json:"metadata_scope,omitempty"`
	MetadataCollection string                                `json:"metadata_collection,omitempty"`
	SourceBucket       string                                `json:"source_bucket"`
	SourceScope        string                                `json:"source_scope,omitempty"`
	SourceCollection   string                                `json:"source_collection,omitempty"`
	BucketBindings     []jsonEventingFunctionBucketBinding   `json:"buckets,omitempty"`
	UrlBindings        []jsonEventingFunctionUrlBinding      `json:"curl,omitempty"`
	ConstantBindings   []jsonEventingFunctionConstantBinding `json:"constants,omitempty"`
}

type jsonEventingFunctionBucketBinding struct {
	Alias      string                       `json:"alias"`
	Bucket     string                       `json:"bucket_name"`
	Scope      string                       `json:"scope_name,omitempty"`
	Collection string                       `json:"collection_name,omitempty"`
	Access     EventingFunctionBucketAccess `json:"access"`
}

type jsonEventingFunctionUrlBinding struct {
	Hostname               string `json:"hostname"`
	Alias                  string `json:"value"`
	AuthType               string `json:"auth_type"`
	AllowCookies           bool   `json:"allow_cookies"`
	ValidateSSLCertificate bool   `json:"validate_ssl_certificate"`
	Username               string `json:"username,omitempty"`
	Password               string `json:"password,omitempty"`
	BearerKey              string `json:"bearer_key,omitempty"`
}

type jsonEventingFunctionConstantBinding struct {
	Alias   string `json:"value"`
	Literal string `json:"literal"`
}

type jsonEventingFunctionsStatusApp struct {
	CompositeStatus       EventingFunctionStatus           `json:"composite_status"`
	Name                  string                           `json:"name"`
	NumBootstrappingNodes int                              `json:"num_bootstrapping_nodes"`
	NumDeployedNodes      int                              `json:"num_deployed_nodes"`
	DeploymentStatus      EventingFunctionDeploymentStatus `json:"deployment_status"`
	ProcessingStatus      EventingFunctionProcessingStatus `json:"processing_status"`
}

type jsonEventingFunctionsStatus struct {
	Apps             []jsonEventingFunctionsStatusApp `json:"apps"`
	NumEventingNodes int                              `json:"num_eventing_nodes"`
}

// EventingFunctionStatus describes the current state of an eventing function.
type EventingFunctionStatus string

var (
	// EventingFunctionStateUndeployed represents that the eventing function is undeployed.
	EventingFunctionStateUndeployed EventingFunctionStatus = "undeployed"

	// EventingFunctionStateDeploying represents that the eventing function is deploying.
	EventingFunctionStateDeploying EventingFunctionStatus = "deploying"

	// EventingFunctionStateDeployed represents that the eventing function is deployed.
	EventingFunctionStateDeployed EventingFunctionStatus = "deployed"

	// EventingFunctionStateUndeploying represents that the eventing function is undeploying.
	EventingFunctionStateUndeploying EventingFunctionStatus = "undeploying"

	// EventingFunctionStatePaused represents that the eventing function is paused.
	EventingFunctionStatePaused EventingFunctionStatus = "paused"

	// EventingFunctionStatePausing represents that the eventing function is pausing.
	EventingFunctionStatePausing EventingFunctionStatus = "pausing"
)

// EventingFunctionDCPBoundary sets what data mutations to deploy the eventing function for.
type EventingFunctionDCPBoundary string

var (
	// EventingFunctionDCPBoundaryEverything will deploy the eventing function for all data mutations.
	EventingFunctionDCPBoundaryEverything EventingFunctionDCPBoundary = "everything"

	// EventingFunctionDCPBoundaryFromNow will deploy the eventing function for only data mutations occurring post deployment.
	EventingFunctionDCPBoundaryFromNow EventingFunctionDCPBoundary = "from_now"
)

// EventingFunctionDeploymentStatus represents the current deployment status for the eventing function.
type EventingFunctionDeploymentStatus bool

var (
	// EventingFunctionDeploymentStatusDeployed represents that the eventing function is currently deployed.
	EventingFunctionDeploymentStatusDeployed EventingFunctionDeploymentStatus = true

	// EventingFunctionDeploymentStatusUndeployed represents that the eventing function is currently undeployed.
	EventingFunctionDeploymentStatusUndeployed EventingFunctionDeploymentStatus = false
)

// EventingFunctionProcessingStatus represents the current processing status for the eventing function.
type EventingFunctionProcessingStatus bool

var (
	// EventingFunctionProcessingStatusRunning represents that the eventing function is currently running.
	EventingFunctionProcessingStatusRunning EventingFunctionProcessingStatus = true

	// EventingFunctionProcessingStatusPaused represents that the eventing function is currently paused.
	EventingFunctionProcessingStatusPaused EventingFunctionProcessingStatus = false
)

// EventingFunctionLanguageCompatibility represents the eventing function language compatibility for backward compatibility.
type EventingFunctionLanguageCompatibility string

var (
	// EventingFunctionLanguageCompatibilityVersion600 represents the eventing function language compatibility 6.0.0.
	EventingFunctionLanguageCompatibilityVersion600 EventingFunctionLanguageCompatibility = "6.0.0"

	// EventingFunctionLanguageCompatibilityVersion650 represents the eventing function language compatibility 6.5.0.
	EventingFunctionLanguageCompatibilityVersion650 EventingFunctionLanguageCompatibility = "6.5.0"

	// EventingFunctionLanguageCompatibilityVersion662 represents the eventing function language compatibility 6.6.2.
	EventingFunctionLanguageCompatibilityVersion662 EventingFunctionLanguageCompatibility = "6.6.2"
)

// EventingFunctionLogLevel represents the granularity at which to log messages for the eventing function.
type EventingFunctionLogLevel string

var (
	// EventingFunctionLogLevelInfo represents to log messages at INFO for the eventing function.
	EventingFunctionLogLevelInfo EventingFunctionLogLevel = "INFO"

	// EventingFunctionLogLevelError represents to log messages at ERROR for the eventing function.
	EventingFunctionLogLevelError EventingFunctionLogLevel = "ERROR"

	// EventingFunctionLogLevelWarning represents to log messages at WARNING for the eventing function.
	EventingFunctionLogLevelWarning EventingFunctionLogLevel = "WARNING"

	// EventingFunctionLogLevelDebug represents to log messages at DEBUG for the eventing function.
	EventingFunctionLogLevelDebug EventingFunctionLogLevel = "DEBUG"

	// EventingFunctionLogLevelTrace represents to log messages at TRACE for the eventing function.
	EventingFunctionLogLevelTrace EventingFunctionLogLevel = "TRACE"
)

// EventingFunctionSettings are the settings for an EventingFunction.
type EventingFunctionSettings struct {
	CPPWorkerThreadCount   int
	DCPStreamBoundary      EventingFunctionDCPBoundary
	Description            string
	DeploymentStatus       EventingFunctionDeploymentStatus
	ProcessingStatus       EventingFunctionProcessingStatus
	LanguageCompatibility  EventingFunctionLanguageCompatibility
	LogLevel               EventingFunctionLogLevel
	ExecutionTimeout       time.Duration
	LCBInstCapacity        int
	LCBRetryCount          int
	LCBTimeout             time.Duration
	QueryConsistency       QueryScanConsistency
	NumTimerPartitions     int
	SockBatchSize          int
	TickDuration           time.Duration
	TimerContextSize       int
	UserPrefix             string
	BucketCacheSize        int
	BucketCacheAge         int
	CurlMaxAllowedRespSize int
	QueryPrepareAll        bool
	WorkerCount            int
	HandlerHeaders         []string
	HandlerFooters         []string
	EnableAppLogRotation   bool
	AppLogDir              string
	AppLogMaxSize          int
	AppLogMaxFiles         int
	CheckpointInterval     time.Duration
}

// EventingFunctionBucketAccess represents the level of access an eventing function has to a bucket.
type EventingFunctionBucketAccess string

var (
	// EventingFunctionBucketAccessReadOnly represents readonly access to a bucket for an eventing function.
	EventingFunctionBucketAccessReadOnly EventingFunctionBucketAccess = "r"

	// EventingFunctionBucketAccessReadWrite represents readwrite access to a bucket for an eventing function.
	EventingFunctionBucketAccessReadWrite EventingFunctionBucketAccess = "rw"
)

// EventingFunctionUrlAuth represents an authentication method for EventingFunctionUrlBinding for an eventing function.
type EventingFunctionUrlAuth interface {
	Method() string
	Username() string
	Password() string
	Key() string
}

// EventingFunctionUrlNoAuth specifies that no authentication is used for the EventingFunctionUrlBinding.
type EventingFunctionUrlNoAuth struct{}

func (ua EventingFunctionUrlNoAuth) Method() string {
	return "no-auth"
}
func (ua EventingFunctionUrlNoAuth) Username() string {
	return ""
}
func (ua EventingFunctionUrlNoAuth) Password() string {
	return ""
}
func (ua EventingFunctionUrlNoAuth) Key() string {
	return ""
}

// EventingFunctionUrlAuthBasic specifies that basic authentication is used for the EventingFunctionUrlBinding.
type EventingFunctionUrlAuthBasic struct {
	User string
	Pass string
}

func (ua EventingFunctionUrlAuthBasic) Method() string {
	return "basic"
}
func (ua EventingFunctionUrlAuthBasic) Username() string {
	return ua.User
}
func (ua EventingFunctionUrlAuthBasic) Password() string {
	return ua.Pass
}
func (ua EventingFunctionUrlAuthBasic) Key() string {
	return ""
}

// EventingFunctionUrlAuthDigest specifies that digest authentication is used for the EventingFunctionUrlBinding.
type EventingFunctionUrlAuthDigest struct {
	User string
	Pass string
}

func (ua EventingFunctionUrlAuthDigest) Method() string {
	return "digest"
}
func (ua EventingFunctionUrlAuthDigest) Username() string {
	return ua.User
}
func (ua EventingFunctionUrlAuthDigest) Password() string {
	return ua.Pass
}
func (ua EventingFunctionUrlAuthDigest) Key() string {
	return ""
}

// EventingFunctionUrlAuthBearer specifies that bearer token authentication is used for the EventingFunctionUrlBinding.
type EventingFunctionUrlAuthBearer struct {
	BearerKey string
}

func (ua EventingFunctionUrlAuthBearer) Method() string {
	return "bearer"
}
func (ua EventingFunctionUrlAuthBearer) Username() string {
	return ""
}
func (ua EventingFunctionUrlAuthBearer) Password() string {
	return ""
}
func (ua EventingFunctionUrlAuthBearer) Key() string {
	return ua.BearerKey
}

// EventingFunctionBucketBinding represents an eventing function binding allowing the function access to buckets,
// scopes, and collections.
type EventingFunctionBucketBinding struct {
	Alias  string
	Name   EventingFunctionKeyspace
	Access EventingFunctionBucketAccess
}

// EventingFunctionUrlBinding represents an eventing function binding allowing the function access external resources
// via cURL.
type EventingFunctionUrlBinding struct {
	Hostname               string
	Alias                  string
	Auth                   EventingFunctionUrlAuth
	AllowCookies           bool
	ValidateSSLCertificate bool
}

// EventingFunctionConstantBinding represents an eventing function binding allowing the function to utilize global variables.
type EventingFunctionConstantBinding struct {
	Alias   string
	Literal string
}

// EventingFunctionKeyspace represents a triple of bucket, collection, and scope names.
type EventingFunctionKeyspace struct {
	Bucket     string
	Scope      string
	Collection string
}

// EventingStatus represents the current state of all eventing functions.
type EventingStatus struct {
	NumEventingNodes int
	Functions        []EventingFunctionState
}

// EventingFunctionState represents the current state of an eventing function.
type EventingFunctionState struct {
	Name                  string
	Status                EventingFunctionStatus
	NumBootstrappingNodes int
	NumDeployedNodes      int
	DeploymentStatus      EventingFunctionDeploymentStatus
	ProcessingStatus      EventingFunctionProcessingStatus
}

func (ef *EventingStatus) UnmarshalJSON(b []byte) error {
	var jf jsonEventingFunctionsStatus
	err := json.Unmarshal(b, &jf)
	if err != nil {
		return err
	}

	var funcs []EventingFunctionState
	for _, f := range jf.Apps {
		funcs = append(funcs, EventingFunctionState{
			Name:                  f.Name,
			Status:                f.CompositeStatus,
			NumBootstrappingNodes: f.NumBootstrappingNodes,
			NumDeployedNodes:      f.NumDeployedNodes,
			DeploymentStatus:      f.DeploymentStatus,
			ProcessingStatus:      f.ProcessingStatus,
		})
	}

	ef.NumEventingNodes = jf.NumEventingNodes
	ef.Functions = funcs

	return nil
}

// EventingFunction represents an eventing function.
type EventingFunction struct {
	Name               string // Required
	Code               string // Required
	Version            string
	EnforceSchema      bool
	HandlerUUID        int
	FunctionInstanceID string
	MetadataKeyspace   EventingFunctionKeyspace // Required
	SourceKeyspace     EventingFunctionKeyspace // Required
	BucketBindings     []EventingFunctionBucketBinding
	UrlBindings        []EventingFunctionUrlBinding
	ConstantBindings   []EventingFunctionConstantBinding
	Settings           EventingFunctionSettings
}

func (ef EventingFunction) MarshalJSON() ([]byte, error) {
	var bucketBindings []jsonEventingFunctionBucketBinding
	for _, b := range ef.BucketBindings {
		bucketBindings = append(bucketBindings, jsonEventingFunctionBucketBinding{
			Alias:      b.Alias,
			Bucket:     b.Name.Bucket,
			Scope:      b.Name.Scope,
			Collection: b.Name.Collection,
			Access:     b.Access,
		})
	}
	var urlBindings []jsonEventingFunctionUrlBinding
	for _, b := range ef.UrlBindings {
		urlBindings = append(urlBindings, jsonEventingFunctionUrlBinding{
			Hostname:               b.Hostname,
			Alias:                  b.Alias,
			AuthType:               b.Auth.Method(),
			AllowCookies:           b.AllowCookies,
			ValidateSSLCertificate: b.ValidateSSLCertificate,
			Username:               b.Auth.Username(),
			Password:               b.Auth.Password(),
			BearerKey:              b.Auth.Key(),
		})
	}
	var constantBindings []jsonEventingFunctionConstantBinding
	for _, b := range ef.ConstantBindings {
		constantBindings = append(constantBindings, jsonEventingFunctionConstantBinding(b))
	}

	jsonSettings := jsonEventingFunction{
		Name:               ef.Name,
		Code:               ef.Code,
		Version:            ef.Version,
		EnforceSchema:      ef.EnforceSchema,
		HandlerUUID:        ef.HandlerUUID,
		FunctionInstanceID: ef.FunctionInstanceID,
		Settings: jsonEventingFunctionSettings{
			CPPWorkerThreadCount:   ef.Settings.CPPWorkerThreadCount,
			DCPStreamBoundary:      ef.Settings.DCPStreamBoundary,
			Description:            ef.Settings.Description,
			DeploymentStatus:       ef.Settings.DeploymentStatus,
			ProcessingStatus:       ef.Settings.ProcessingStatus,
			LanguageCompatibility:  ef.Settings.LanguageCompatibility,
			LogLevel:               ef.Settings.LogLevel,
			ExecutionTimeout:       int(ef.Settings.ExecutionTimeout.Seconds()),
			LCBInstCapacity:        ef.Settings.LCBInstCapacity,
			LCBRetryCount:          ef.Settings.LCBRetryCount,
			LCBTimeout:             int(ef.Settings.LCBTimeout.Seconds()),
			QueryConsistency:       ef.Settings.QueryConsistency,
			NumTimerPartitions:     ef.Settings.NumTimerPartitions,
			SockBatchSize:          ef.Settings.SockBatchSize,
			TickDuration:           int(ef.Settings.TickDuration.Milliseconds()),
			TimerContextSize:       ef.Settings.TimerContextSize,
			UserPrefix:             ef.Settings.UserPrefix,
			BucketCacheSize:        ef.Settings.BucketCacheSize,
			BucketCacheAge:         ef.Settings.BucketCacheAge,
			CurlMaxAllowedRespSize: ef.Settings.CurlMaxAllowedRespSize,
			QueryPrepareAll:        ef.Settings.QueryPrepareAll,
			WorkerCount:            ef.Settings.WorkerCount,
			HandlerHeaders:         ef.Settings.HandlerHeaders,
			HandlerFooters:         ef.Settings.HandlerFooters,
			EnableAppLogRotation:   ef.Settings.EnableAppLogRotation,
			AppLogDir:              ef.Settings.AppLogDir,
			AppLogMaxSize:          ef.Settings.AppLogMaxSize,
			AppLogMaxFiles:         ef.Settings.AppLogMaxFiles,
			CheckpointInterval:     int(ef.Settings.CheckpointInterval.Seconds()),
		},
		DeploymentConfig: jsonEventingFunctionDeploymentConfig{
			MetadataBucket:     ef.MetadataKeyspace.Bucket,
			MetadataScope:      ef.MetadataKeyspace.Scope,
			MetadataCollection: ef.MetadataKeyspace.Collection,
			SourceBucket:       ef.SourceKeyspace.Bucket,
			SourceScope:        ef.SourceKeyspace.Scope,
			SourceCollection:   ef.SourceKeyspace.Collection,
			BucketBindings:     bucketBindings,
			UrlBindings:        urlBindings,
			ConstantBindings:   constantBindings,
		},
	}

	return json.Marshal(jsonSettings)
}

func (ef *EventingFunction) UnmarshalJSON(b []byte) error {
	var jf jsonEventingFunction
	err := json.Unmarshal(b, &jf)
	if err != nil {
		return err
	}

	var bucketBindings []EventingFunctionBucketBinding
	for _, b := range jf.DeploymentConfig.BucketBindings {
		bucketBindings = append(bucketBindings, EventingFunctionBucketBinding{
			Alias: b.Alias,
			Name: EventingFunctionKeyspace{
				Bucket:     b.Bucket,
				Scope:      b.Scope,
				Collection: b.Collection,
			},
			Access: b.Access,
		})
	}
	var urlBindings []EventingFunctionUrlBinding
	for _, b := range jf.DeploymentConfig.UrlBindings {
		var auth EventingFunctionUrlAuth
		switch b.AuthType {
		case "no-auth":
			auth = EventingFunctionUrlNoAuth{}
		case "basic":
			auth = EventingFunctionUrlAuthBasic{
				User: b.Username,
			}
		case "digest":
			auth = EventingFunctionUrlAuthDigest{
				User: b.Username,
			}
		case "bearer":
			auth = EventingFunctionUrlAuthBearer{}
		}

		urlBindings = append(urlBindings, EventingFunctionUrlBinding{
			Hostname:               b.Hostname,
			Alias:                  b.Alias,
			Auth:                   auth,
			AllowCookies:           b.AllowCookies,
			ValidateSSLCertificate: b.ValidateSSLCertificate,
		})
	}
	var constantBindings []EventingFunctionConstantBinding
	for _, b := range jf.DeploymentConfig.ConstantBindings {
		constantBindings = append(constantBindings, EventingFunctionConstantBinding(b))
	}

	ef.Name = jf.Name
	ef.Code = jf.Code
	ef.Version = jf.Version
	ef.EnforceSchema = jf.EnforceSchema
	ef.HandlerUUID = jf.HandlerUUID
	ef.FunctionInstanceID = jf.FunctionInstanceID
	ef.Settings = EventingFunctionSettings{
		CPPWorkerThreadCount:   jf.Settings.CPPWorkerThreadCount,
		DCPStreamBoundary:      jf.Settings.DCPStreamBoundary,
		Description:            jf.Settings.Description,
		DeploymentStatus:       jf.Settings.DeploymentStatus,
		ProcessingStatus:       jf.Settings.ProcessingStatus,
		LanguageCompatibility:  jf.Settings.LanguageCompatibility,
		LogLevel:               jf.Settings.LogLevel,
		ExecutionTimeout:       time.Duration(jf.Settings.ExecutionTimeout) * time.Second,
		LCBInstCapacity:        jf.Settings.LCBInstCapacity,
		LCBRetryCount:          jf.Settings.LCBRetryCount,
		LCBTimeout:             time.Duration(jf.Settings.LCBTimeout) * time.Second,
		QueryConsistency:       jf.Settings.QueryConsistency,
		NumTimerPartitions:     jf.Settings.NumTimerPartitions,
		SockBatchSize:          jf.Settings.SockBatchSize,
		TickDuration:           time.Duration(jf.Settings.TickDuration) * time.Millisecond,
		TimerContextSize:       jf.Settings.TimerContextSize,
		UserPrefix:             jf.Settings.UserPrefix,
		BucketCacheSize:        jf.Settings.BucketCacheSize,
		BucketCacheAge:         jf.Settings.BucketCacheAge,
		CurlMaxAllowedRespSize: jf.Settings.CurlMaxAllowedRespSize,
		QueryPrepareAll:        jf.Settings.QueryPrepareAll,
		WorkerCount:            jf.Settings.WorkerCount,
		HandlerHeaders:         jf.Settings.HandlerHeaders,
		HandlerFooters:         jf.Settings.HandlerFooters,
		EnableAppLogRotation:   jf.Settings.EnableAppLogRotation,
		AppLogDir:              jf.Settings.AppLogDir,
		AppLogMaxSize:          jf.Settings.AppLogMaxSize,
		AppLogMaxFiles:         jf.Settings.AppLogMaxFiles,
		CheckpointInterval:     time.Duration(jf.Settings.CheckpointInterval) * time.Second,
	}
	ef.MetadataKeyspace = EventingFunctionKeyspace{
		Bucket:     jf.DeploymentConfig.MetadataBucket,
		Scope:      jf.DeploymentConfig.MetadataScope,
		Collection: jf.DeploymentConfig.MetadataCollection,
	}
	ef.SourceKeyspace = EventingFunctionKeyspace{
		Bucket:     jf.DeploymentConfig.SourceBucket,
		Scope:      jf.DeploymentConfig.SourceScope,
		Collection: jf.DeploymentConfig.SourceCollection,
	}
	ef.BucketBindings = bucketBindings
	ef.UrlBindings = urlBindings
	ef.ConstantBindings = constantBindings

	return nil
}

type eventingRequestOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan
	Context       context.Context
}

func (efm *EventingFunctionManager) doRequest(path string, method string, opName string, function *EventingFunction,
	target interface{}, opts eventingRequestOptions) error {
	start := time.Now()
	defer efm.meter.ValueRecord(meterValueServiceManagement, opName, start)

	op := "manager_eventing_" + opName
	span := createSpan(efm.tracer, opts.ParentSpan, op, "management")
	span.SetAttribute("db.operation", method+" "+path)
	defer span.End()

	var b []byte
	if function != nil {
		var err error
		b, err = json.Marshal(function)
		if err != nil {
			return err
		}
	}

	req := mgmtRequest{
		Service:       ServiceTypeEventing,
		Method:        method,
		Path:          path,
		RetryStrategy: opts.RetryStrategy,
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
		Body:          b,
	}
	resp, err := efm.doMgmtRequest(opts.Context, req)
	if err != nil {
		return err
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		idxErr := efm.tryParseErrorMessage(&req, resp)
		if idxErr != nil {
			return idxErr
		}

		return makeMgmtBadStatusError("failed eventing "+opName, &req, resp)
	}

	if target != nil {
		jsonDec := json.NewDecoder(resp.Body)
		err = jsonDec.Decode(target)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpsertEventingFunctionOptions are the options available when using the UpsertFunction operation.
type UpsertEventingFunctionOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// UpsertFunction inserts or updates an eventing function.
func (efm *EventingFunctionManager) UpsertFunction(function EventingFunction, opts *UpsertEventingFunctionOptions) error {
	if opts == nil {
		opts = &UpsertEventingFunctionOptions{}
	}

	return efm.doRequest(fmt.Sprintf("/api/v1/functions/%s", function.Name), "POST",
		"upsert_function", &function, nil, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
}

// DropEventingFunctionOptions are the options available when using the DropFunction operation.
type DropEventingFunctionOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropFunction drops an eventing function.
func (efm *EventingFunctionManager) DropFunction(name string, opts *DropEventingFunctionOptions) error {
	if opts == nil {
		opts = &DropEventingFunctionOptions{}
	}

	return efm.doRequest(fmt.Sprintf("/api/v1/functions/%s", name), "DELETE",
		"drop_function", nil, nil, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
}

// DeployEventingFunctionOptions are the options available when using the DeployFunction operation.
type DeployEventingFunctionOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DeployFunction deploys an eventing function.
func (efm *EventingFunctionManager) DeployFunction(name string, opts *DeployEventingFunctionOptions) error {
	if opts == nil {
		opts = &DeployEventingFunctionOptions{}
	}

	return efm.doRequest(fmt.Sprintf("/api/v1/functions/%s/deploy", name), "POST",
		"deploy_function", nil, nil, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
}

// UndeployEventingFunctionOptions are the options available when using the UndeployFunction operation.
type UndeployEventingFunctionOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// UndeployFunction undeploys an eventing function.
func (efm *EventingFunctionManager) UndeployFunction(name string, opts *UndeployEventingFunctionOptions) error {
	if opts == nil {
		opts = &UndeployEventingFunctionOptions{}
	}

	return efm.doRequest(fmt.Sprintf("/api/v1/functions/%s/undeploy", name), "POST",
		"undeploy_function", nil, nil, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
}

// GetAllEventingFunctionsOptions are the options available when using the GetAllFunctions operation.
type GetAllEventingFunctionsOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetAllFunctions fetches all of the eventing functions.
func (efm *EventingFunctionManager) GetAllFunctions(opts *GetAllEventingFunctionsOptions) ([]EventingFunction, error) {
	if opts == nil {
		opts = &GetAllEventingFunctionsOptions{}
	}

	var functions []EventingFunction
	err := efm.doRequest("/api/v1/functions", "GET",
		"get_all_functions", nil, &functions, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
	if err != nil {
		return nil, err
	}

	return functions, nil
}

// GetEventingFunctionOptions are the options available when using the GetFunction operation.
type GetEventingFunctionOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetFunction fetches an eventing function.
func (efm *EventingFunctionManager) GetFunction(name string, opts *GetEventingFunctionOptions) (*EventingFunction, error) {
	if opts == nil {
		opts = &GetEventingFunctionOptions{}
	}

	var function *EventingFunction
	err := efm.doRequest(fmt.Sprintf("/api/v1/functions/%s", name), "GET",
		"get_function", nil, &function, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
	if err != nil {
		return nil, err
	}

	return function, nil
}

// PauseEventingFunctionOptions are the options available when using the PauseFunction operation.
type PauseEventingFunctionOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// PauseFunction pauses an eventing function.
func (efm *EventingFunctionManager) PauseFunction(name string, opts *PauseEventingFunctionOptions) error {
	if opts == nil {
		opts = &PauseEventingFunctionOptions{}
	}

	return efm.doRequest(fmt.Sprintf("/api/v1/functions/%s/pause", name), "POST",
		"pause_function", nil, nil, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
}

// ResumeEventingFunctionOptions are the options available when using the ResumeFunction operation.
type ResumeEventingFunctionOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// ResumeFunction resumes an eventing function.
func (efm *EventingFunctionManager) ResumeFunction(name string, opts *ResumeEventingFunctionOptions) error {
	if opts == nil {
		opts = &ResumeEventingFunctionOptions{}
	}

	return efm.doRequest(fmt.Sprintf("/api/v1/functions/%s/resume", name), "POST",
		"resume_function", nil, nil, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
}

// EventingFunctionsStatusOptions are the options available when using the FunctionsStatus operation.
type EventingFunctionsStatusOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// FunctionsStatus fetches the current status of all eventing functions.
func (efm *EventingFunctionManager) FunctionsStatus(opts *EventingFunctionsStatusOptions) (*EventingStatus, error) {
	if opts == nil {
		opts = &EventingFunctionsStatusOptions{}
	}

	var functions *EventingStatus
	err := efm.doRequest("/api/v1/status", "GET",
		"functions_status", nil, &functions, eventingRequestOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
			ParentSpan:    opts.ParentSpan,
			Context:       opts.Context,
		})
	if err != nil {
		return nil, err
	}

	return functions, nil
}
