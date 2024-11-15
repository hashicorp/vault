package gocb

import (
	"context"
	"encoding/json"
	"time"
)

// EventingFunctionManager provides methods for performing eventing function management operations.
// This manager is designed to work only against Couchbase Server 7.0+, it might work against earlier server
// versions but that is not tested and is not supported.
//
// # UNCOMMITTED
//
// This API is UNCOMMITTED and may change in the future.
type EventingFunctionManager struct {
	controller *providerController[eventingManagementProvider]
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

	// EventingFunctionLanguageCompatibilityVersion720 represents the eventing function language compatibility 7.2.0
	EventingFunctionLanguageCompatibilityVersion720 EventingFunctionLanguageCompatibility = "7.2.0"
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

func (es *EventingStatus) UnmarshalJSON(b []byte) error {
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

	es.NumEventingNodes = jf.NumEventingNodes
	es.Functions = funcs

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

func (ef *EventingFunction) MarshalJSON() ([]byte, error) {
	jsonSettings := ef.toJSONEventingFunction()
	return json.Marshal(jsonSettings)
}

func (ef *EventingFunction) UnmarshalJSON(b []byte) error {
	var jf jsonEventingFunction
	err := json.Unmarshal(b, &jf)
	if err != nil {
		return err
	}
	ef.fromJSONEventingFunction(jf)
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
	return autoOpControlErrorOnly(efm.controller, func(provider eventingManagementProvider) error {
		if opts == nil {
			opts = &UpsertEventingFunctionOptions{}
		}

		return provider.UpsertFunction(nil, function, opts)
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
	return autoOpControlErrorOnly(efm.controller, func(provider eventingManagementProvider) error {
		if opts == nil {
			opts = &DropEventingFunctionOptions{}
		}

		return provider.DropFunction(nil, name, opts)
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
	return autoOpControlErrorOnly(efm.controller, func(provider eventingManagementProvider) error {
		if opts == nil {
			opts = &DeployEventingFunctionOptions{}
		}

		return provider.DeployFunction(nil, name, opts)
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
	return autoOpControlErrorOnly(efm.controller, func(provider eventingManagementProvider) error {
		if opts == nil {
			opts = &UndeployEventingFunctionOptions{}
		}

		return provider.UndeployFunction(nil, name, opts)
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

// GetAllFunctions fetches all the eventing functions.
func (efm *EventingFunctionManager) GetAllFunctions(opts *GetAllEventingFunctionsOptions) ([]EventingFunction, error) {
	return autoOpControl(efm.controller, func(provider eventingManagementProvider) ([]EventingFunction, error) {
		if opts == nil {
			opts = &GetAllEventingFunctionsOptions{}
		}

		return provider.GetAllFunctions(nil, opts)
	})
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
	return autoOpControl(efm.controller, func(provider eventingManagementProvider) (*EventingFunction, error) {
		if opts == nil {
			opts = &GetEventingFunctionOptions{}
		}

		return provider.GetFunction(nil, name, opts)
	})
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
	return autoOpControlErrorOnly(efm.controller, func(provider eventingManagementProvider) error {
		if opts == nil {
			opts = &PauseEventingFunctionOptions{}
		}

		return provider.PauseFunction(nil, name, opts)
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
	return autoOpControlErrorOnly(efm.controller, func(provider eventingManagementProvider) error {
		if opts == nil {
			opts = &ResumeEventingFunctionOptions{}
		}

		return provider.ResumeFunction(nil, name, opts)
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
	return autoOpControl(efm.controller, func(provider eventingManagementProvider) (*EventingStatus, error) {
		if opts == nil {
			opts = &EventingFunctionsStatusOptions{}
		}

		return provider.FunctionsStatus(nil, opts)
	})
}
