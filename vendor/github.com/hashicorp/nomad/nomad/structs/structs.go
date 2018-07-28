package structs

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/blake2b"

	"github.com/gorhill/cronexpr"
	"github.com/hashicorp/consul/api"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/nomad/acl"
	"github.com/hashicorp/nomad/helper"
	"github.com/hashicorp/nomad/helper/args"
	"github.com/hashicorp/nomad/helper/uuid"
	"github.com/mitchellh/copystructure"
	"github.com/ugorji/go/codec"

	"math"

	hcodec "github.com/hashicorp/go-msgpack/codec"
)

var (
	// validPolicyName is used to validate a policy name
	validPolicyName = regexp.MustCompile("^[a-zA-Z0-9-]{1,128}$")

	// b32 is a lowercase base32 encoding for use in URL friendly service hashes
	b32 = base32.NewEncoding(strings.ToLower("abcdefghijklmnopqrstuvwxyz234567"))
)

type MessageType uint8

const (
	NodeRegisterRequestType MessageType = iota
	NodeDeregisterRequestType
	NodeUpdateStatusRequestType
	NodeUpdateDrainRequestType
	JobRegisterRequestType
	JobDeregisterRequestType
	EvalUpdateRequestType
	EvalDeleteRequestType
	AllocUpdateRequestType
	AllocClientUpdateRequestType
	ReconcileJobSummariesRequestType
	VaultAccessorRegisterRequestType
	VaultAccessorDeregisterRequestType
	ApplyPlanResultsRequestType
	DeploymentStatusUpdateRequestType
	DeploymentPromoteRequestType
	DeploymentAllocHealthRequestType
	DeploymentDeleteRequestType
	JobStabilityRequestType
	ACLPolicyUpsertRequestType
	ACLPolicyDeleteRequestType
	ACLTokenUpsertRequestType
	ACLTokenDeleteRequestType
	ACLTokenBootstrapRequestType
	AutopilotRequestType
	UpsertNodeEventsType
	JobBatchDeregisterRequestType
	AllocUpdateDesiredTransitionRequestType
	NodeUpdateEligibilityRequestType
	BatchNodeUpdateDrainRequestType
)

const (
	// IgnoreUnknownTypeFlag is set along with a MessageType
	// to indicate that the message type can be safely ignored
	// if it is not recognized. This is for future proofing, so
	// that new commands can be added in a way that won't cause
	// old servers to crash when the FSM attempts to process them.
	IgnoreUnknownTypeFlag MessageType = 128

	// ApiMajorVersion is returned as part of the Status.Version request.
	// It should be incremented anytime the APIs are changed in a way
	// that would break clients for sane client versioning.
	ApiMajorVersion = 1

	// ApiMinorVersion is returned as part of the Status.Version request.
	// It should be incremented anytime the APIs are changed to allow
	// for sane client versioning. Minor changes should be compatible
	// within the major version.
	ApiMinorVersion = 1

	ProtocolVersion = "protocol"
	APIMajorVersion = "api.major"
	APIMinorVersion = "api.minor"

	GetterModeAny  = "any"
	GetterModeFile = "file"
	GetterModeDir  = "dir"

	// maxPolicyDescriptionLength limits a policy description length
	maxPolicyDescriptionLength = 256

	// maxTokenNameLength limits a ACL token name length
	maxTokenNameLength = 256

	// ACLClientToken and ACLManagementToken are the only types of tokens
	ACLClientToken     = "client"
	ACLManagementToken = "management"

	// DefaultNamespace is the default namespace.
	DefaultNamespace            = "default"
	DefaultNamespaceDescription = "Default shared namespace"

	// JitterFraction is a the limit to the amount of jitter we apply
	// to a user specified MaxQueryTime. We divide the specified time by
	// the fraction. So 16 == 6.25% limit of jitter. This jitter is also
	// applied to RPCHoldTimeout.
	JitterFraction = 16

	// MaxRetainedNodeEvents is the maximum number of node events that will be
	// retained for a single node
	MaxRetainedNodeEvents = 10
)

// Context defines the scope in which a search for Nomad object operates, and
// is also used to query the matching index value for this context
type Context string

const (
	Allocs      Context = "allocs"
	Deployments Context = "deployment"
	Evals       Context = "evals"
	Jobs        Context = "jobs"
	Nodes       Context = "nodes"
	Namespaces  Context = "namespaces"
	Quotas      Context = "quotas"
	All         Context = "all"
)

// NamespacedID is a tuple of an ID and a namespace
type NamespacedID struct {
	ID        string
	Namespace string
}

func (n NamespacedID) String() string {
	return fmt.Sprintf("<ns: %q, id: %q>", n.Namespace, n.ID)
}

// RPCInfo is used to describe common information about query
type RPCInfo interface {
	RequestRegion() string
	IsRead() bool
	AllowStaleRead() bool
	IsForwarded() bool
	SetForwarded()
}

// InternalRpcInfo allows adding internal RPC metadata to an RPC. This struct
// should NOT be replicated in the API package as it is internal only.
type InternalRpcInfo struct {
	// Forwarded marks whether the RPC has been forwarded.
	Forwarded bool
}

// IsForwarded returns whether the RPC is forwarded from another server.
func (i *InternalRpcInfo) IsForwarded() bool {
	return i.Forwarded
}

// SetForwarded marks that the RPC is being forwarded from another server.
func (i *InternalRpcInfo) SetForwarded() {
	i.Forwarded = true
}

// QueryOptions is used to specify various flags for read queries
type QueryOptions struct {
	// The target region for this query
	Region string

	// Namespace is the target namespace for the query.
	Namespace string

	// If set, wait until query exceeds given index. Must be provided
	// with MaxQueryTime.
	MinQueryIndex uint64

	// Provided with MinQueryIndex to wait for change.
	MaxQueryTime time.Duration

	// If set, any follower can service the request. Results
	// may be arbitrarily stale.
	AllowStale bool

	// If set, used as prefix for resource list searches
	Prefix string

	// AuthToken is secret portion of the ACL token used for the request
	AuthToken string

	InternalRpcInfo
}

func (q QueryOptions) RequestRegion() string {
	return q.Region
}

func (q QueryOptions) RequestNamespace() string {
	if q.Namespace == "" {
		return DefaultNamespace
	}
	return q.Namespace
}

// QueryOption only applies to reads, so always true
func (q QueryOptions) IsRead() bool {
	return true
}

func (q QueryOptions) AllowStaleRead() bool {
	return q.AllowStale
}

type WriteRequest struct {
	// The target region for this write
	Region string

	// Namespace is the target namespace for the write.
	Namespace string

	// AuthToken is secret portion of the ACL token used for the request
	AuthToken string

	InternalRpcInfo
}

func (w WriteRequest) RequestRegion() string {
	// The target region for this request
	return w.Region
}

func (w WriteRequest) RequestNamespace() string {
	if w.Namespace == "" {
		return DefaultNamespace
	}
	return w.Namespace
}

// WriteRequest only applies to writes, always false
func (w WriteRequest) IsRead() bool {
	return false
}

func (w WriteRequest) AllowStaleRead() bool {
	return false
}

// QueryMeta allows a query response to include potentially
// useful metadata about a query
type QueryMeta struct {
	// This is the index associated with the read
	Index uint64

	// If AllowStale is used, this is time elapsed since
	// last contact between the follower and leader. This
	// can be used to gauge staleness.
	LastContact time.Duration

	// Used to indicate if there is a known leader node
	KnownLeader bool
}

// WriteMeta allows a write response to include potentially
// useful metadata about the write
type WriteMeta struct {
	// This is the index associated with the write
	Index uint64
}

// NodeRegisterRequest is used for Node.Register endpoint
// to register a node as being a schedulable entity.
type NodeRegisterRequest struct {
	Node      *Node
	NodeEvent *NodeEvent
	WriteRequest
}

// NodeDeregisterRequest is used for Node.Deregister endpoint
// to deregister a node as being a schedulable entity.
type NodeDeregisterRequest struct {
	NodeID string
	WriteRequest
}

// NodeServerInfo is used to in NodeUpdateResponse to return Nomad server
// information used in RPC server lists.
type NodeServerInfo struct {
	// RPCAdvertiseAddr is the IP endpoint that a Nomad Server wishes to
	// be contacted at for RPCs.
	RPCAdvertiseAddr string

	// RpcMajorVersion is the major version number the Nomad Server
	// supports
	RPCMajorVersion int32

	// RpcMinorVersion is the minor version number the Nomad Server
	// supports
	RPCMinorVersion int32

	// Datacenter is the datacenter that a Nomad server belongs to
	Datacenter string
}

// NodeUpdateStatusRequest is used for Node.UpdateStatus endpoint
// to update the status of a node.
type NodeUpdateStatusRequest struct {
	NodeID    string
	Status    string
	NodeEvent *NodeEvent
	WriteRequest
}

// NodeUpdateDrainRequest is used for updating the drain strategy
type NodeUpdateDrainRequest struct {
	NodeID        string
	DrainStrategy *DrainStrategy

	// COMPAT Remove in version 0.10
	// As part of Nomad 0.8 we have deprecated the drain boolean in favor of a
	// drain strategy but we need to handle the upgrade path where the Raft log
	// contains drain updates with just the drain boolean being manipulated.
	Drain bool

	// MarkEligible marks the node as eligible if removing the drain strategy.
	MarkEligible bool

	// NodeEvent is the event added to the node
	NodeEvent *NodeEvent

	WriteRequest
}

// BatchNodeUpdateDrainRequest is used for updating the drain strategy for a
// batch of nodes
type BatchNodeUpdateDrainRequest struct {
	// Updates is a mapping of nodes to their updated drain strategy
	Updates map[string]*DrainUpdate

	// NodeEvents is a mapping of the node to the event to add to the node
	NodeEvents map[string]*NodeEvent

	WriteRequest
}

// DrainUpdate is used to update the drain of a node
type DrainUpdate struct {
	// DrainStrategy is the new strategy for the node
	DrainStrategy *DrainStrategy

	// MarkEligible marks the node as eligible if removing the drain strategy.
	MarkEligible bool
}

// NodeUpdateEligibilityRequest is used for updating the scheduling	eligibility
type NodeUpdateEligibilityRequest struct {
	NodeID      string
	Eligibility string

	// NodeEvent is the event added to the node
	NodeEvent *NodeEvent

	WriteRequest
}

// NodeEvaluateRequest is used to re-evaluate the node
type NodeEvaluateRequest struct {
	NodeID string
	WriteRequest
}

// NodeSpecificRequest is used when we just need to specify a target node
type NodeSpecificRequest struct {
	NodeID   string
	SecretID string
	QueryOptions
}

// SearchResponse is used to return matches and information about whether
// the match list is truncated specific to each type of context.
type SearchResponse struct {
	// Map of context types to ids which match a specified prefix
	Matches map[Context][]string

	// Truncations indicates whether the matches for a particular context have
	// been truncated
	Truncations map[Context]bool

	QueryMeta
}

// SearchRequest is used to parameterize a request, and returns a
// list of matches made up of jobs, allocations, evaluations, and/or nodes,
// along with whether or not the information returned is truncated.
type SearchRequest struct {
	// Prefix is what ids are matched to. I.e, if the given prefix were
	// "a", potential matches might be "abcd" or "aabb"
	Prefix string

	// Context is the type that can be matched against. A context can be a job,
	// node, evaluation, allocation, or empty (indicated every context should be
	// matched)
	Context Context

	QueryOptions
}

// JobRegisterRequest is used for Job.Register endpoint
// to register a job as being a schedulable entity.
type JobRegisterRequest struct {
	Job *Job

	// If EnforceIndex is set then the job will only be registered if the passed
	// JobModifyIndex matches the current Jobs index. If the index is zero, the
	// register only occurs if the job is new.
	EnforceIndex   bool
	JobModifyIndex uint64

	// PolicyOverride is set when the user is attempting to override any policies
	PolicyOverride bool

	WriteRequest
}

// JobDeregisterRequest is used for Job.Deregister endpoint
// to deregister a job as being a schedulable entity.
type JobDeregisterRequest struct {
	JobID string

	// Purge controls whether the deregister purges the job from the system or
	// whether the job is just marked as stopped and will be removed by the
	// garbage collector
	Purge bool

	WriteRequest
}

// JobBatchDeregisterRequest is used to batch deregister jobs and upsert
// evaluations.
type JobBatchDeregisterRequest struct {
	// Jobs is the set of jobs to deregister
	Jobs map[NamespacedID]*JobDeregisterOptions

	// Evals is the set of evaluations to create.
	Evals []*Evaluation

	WriteRequest
}

// JobDeregisterOptions configures how a job is deregistered.
type JobDeregisterOptions struct {
	// Purge controls whether the deregister purges the job from the system or
	// whether the job is just marked as stopped and will be removed by the
	// garbage collector
	Purge bool
}

// JobEvaluateRequest is used when we just need to re-evaluate a target job
type JobEvaluateRequest struct {
	JobID       string
	EvalOptions EvalOptions
	WriteRequest
}

// EvalOptions is used to encapsulate options when forcing a job evaluation
type EvalOptions struct {
	ForceReschedule bool
}

// JobSpecificRequest is used when we just need to specify a target job
type JobSpecificRequest struct {
	JobID     string
	AllAllocs bool
	QueryOptions
}

// JobListRequest is used to parameterize a list request
type JobListRequest struct {
	QueryOptions
}

// JobPlanRequest is used for the Job.Plan endpoint to trigger a dry-run
// evaluation of the Job.
type JobPlanRequest struct {
	Job  *Job
	Diff bool // Toggles an annotated diff
	// PolicyOverride is set when the user is attempting to override any policies
	PolicyOverride bool
	WriteRequest
}

// JobSummaryRequest is used when we just need to get a specific job summary
type JobSummaryRequest struct {
	JobID string
	QueryOptions
}

// JobDispatchRequest is used to dispatch a job based on a parameterized job
type JobDispatchRequest struct {
	JobID   string
	Payload []byte
	Meta    map[string]string
	WriteRequest
}

// JobValidateRequest is used to validate a job
type JobValidateRequest struct {
	Job *Job
	WriteRequest
}

// JobRevertRequest is used to revert a job to a prior version.
type JobRevertRequest struct {
	// JobID is the ID of the job  being reverted
	JobID string

	// JobVersion the version to revert to.
	JobVersion uint64

	// EnforcePriorVersion if set will enforce that the job is at the given
	// version before reverting.
	EnforcePriorVersion *uint64

	WriteRequest
}

// JobStabilityRequest is used to marked a job as stable.
type JobStabilityRequest struct {
	// Job to set the stability on
	JobID      string
	JobVersion uint64

	// Set the stability
	Stable bool
	WriteRequest
}

// JobStabilityResponse is the response when marking a job as stable.
type JobStabilityResponse struct {
	WriteMeta
}

// NodeListRequest is used to parameterize a list request
type NodeListRequest struct {
	QueryOptions
}

// EvalUpdateRequest is used for upserting evaluations.
type EvalUpdateRequest struct {
	Evals     []*Evaluation
	EvalToken string
	WriteRequest
}

// EvalDeleteRequest is used for deleting an evaluation.
type EvalDeleteRequest struct {
	Evals  []string
	Allocs []string
	WriteRequest
}

// EvalSpecificRequest is used when we just need to specify a target evaluation
type EvalSpecificRequest struct {
	EvalID string
	QueryOptions
}

// EvalAckRequest is used to Ack/Nack a specific evaluation
type EvalAckRequest struct {
	EvalID string
	Token  string
	WriteRequest
}

// EvalDequeueRequest is used when we want to dequeue an evaluation
type EvalDequeueRequest struct {
	Schedulers       []string
	Timeout          time.Duration
	SchedulerVersion uint16
	WriteRequest
}

// EvalListRequest is used to list the evaluations
type EvalListRequest struct {
	QueryOptions
}

// PlanRequest is used to submit an allocation plan to the leader
type PlanRequest struct {
	Plan *Plan
	WriteRequest
}

// ApplyPlanResultsRequest is used by the planner to apply a Raft transaction
// committing the result of a plan.
type ApplyPlanResultsRequest struct {
	// AllocUpdateRequest holds the allocation updates to be made by the
	// scheduler.
	AllocUpdateRequest

	// Deployment is the deployment created or updated as a result of a
	// scheduling event.
	Deployment *Deployment

	// DeploymentUpdates is a set of status updates to apply to the given
	// deployments. This allows the scheduler to cancel any unneeded deployment
	// because the job is stopped or the update block is removed.
	DeploymentUpdates []*DeploymentStatusUpdate

	// EvalID is the eval ID of the plan being applied. The modify index of the
	// evaluation is updated as part of applying the plan to ensure that subsequent
	// scheduling events for the same job will wait for the index that last produced
	// state changes. This is necessary for blocked evaluations since they can be
	// processed many times, potentially making state updates, without the state of
	// the evaluation itself being updated.
	EvalID string
}

// AllocUpdateRequest is used to submit changes to allocations, either
// to cause evictions or to assign new allocations. Both can be done
// within a single transaction
type AllocUpdateRequest struct {
	// Alloc is the list of new allocations to assign
	Alloc []*Allocation

	// Evals is the list of new evaluations to create
	// Evals are valid only when used in the Raft RPC
	Evals []*Evaluation

	// Job is the shared parent job of the allocations.
	// It is pulled out since it is common to reduce payload size.
	Job *Job

	WriteRequest
}

// AllocUpdateDesiredTransitionRequest is used to submit changes to allocations
// desired transition state.
type AllocUpdateDesiredTransitionRequest struct {
	// Allocs is the mapping of allocation ids to their desired state
	// transition
	Allocs map[string]*DesiredTransition

	// Evals is the set of evaluations to create
	Evals []*Evaluation

	WriteRequest
}

// AllocListRequest is used to request a list of allocations
type AllocListRequest struct {
	QueryOptions
}

// AllocSpecificRequest is used to query a specific allocation
type AllocSpecificRequest struct {
	AllocID string
	QueryOptions
}

// AllocsGetRequest is used to query a set of allocations
type AllocsGetRequest struct {
	AllocIDs []string
	QueryOptions
}

// PeriodicForceRequest is used to force a specific periodic job.
type PeriodicForceRequest struct {
	JobID string
	WriteRequest
}

// ServerMembersResponse has the list of servers in a cluster
type ServerMembersResponse struct {
	ServerName   string
	ServerRegion string
	ServerDC     string
	Members      []*ServerMember
}

// ServerMember holds information about a Nomad server agent in a cluster
type ServerMember struct {
	Name        string
	Addr        net.IP
	Port        uint16
	Tags        map[string]string
	Status      string
	ProtocolMin uint8
	ProtocolMax uint8
	ProtocolCur uint8
	DelegateMin uint8
	DelegateMax uint8
	DelegateCur uint8
}

// DeriveVaultTokenRequest is used to request wrapped Vault tokens for the
// following tasks in the given allocation
type DeriveVaultTokenRequest struct {
	NodeID   string
	SecretID string
	AllocID  string
	Tasks    []string
	QueryOptions
}

// VaultAccessorsRequest is used to operate on a set of Vault accessors
type VaultAccessorsRequest struct {
	Accessors []*VaultAccessor
}

// VaultAccessor is a reference to a created Vault token on behalf of
// an allocation's task.
type VaultAccessor struct {
	AllocID     string
	Task        string
	NodeID      string
	Accessor    string
	CreationTTL int

	// Raft Indexes
	CreateIndex uint64
}

// DeriveVaultTokenResponse returns the wrapped tokens for each requested task
type DeriveVaultTokenResponse struct {
	// Tasks is a mapping between the task name and the wrapped token
	Tasks map[string]string

	// Error stores any error that occurred. Errors are stored here so we can
	// communicate whether it is retriable
	Error *RecoverableError

	QueryMeta
}

// GenericRequest is used to request where no
// specific information is needed.
type GenericRequest struct {
	QueryOptions
}

// DeploymentListRequest is used to list the deployments
type DeploymentListRequest struct {
	QueryOptions
}

// DeploymentDeleteRequest is used for deleting deployments.
type DeploymentDeleteRequest struct {
	Deployments []string
	WriteRequest
}

// DeploymentStatusUpdateRequest is used to update the status of a deployment as
// well as optionally creating an evaluation atomically.
type DeploymentStatusUpdateRequest struct {
	// Eval, if set, is used to create an evaluation at the same time as
	// updating the status of a deployment.
	Eval *Evaluation

	// DeploymentUpdate is a status update to apply to the given
	// deployment.
	DeploymentUpdate *DeploymentStatusUpdate

	// Job is used to optionally upsert a job. This is used when setting the
	// allocation health results in a deployment failure and the deployment
	// auto-reverts to the latest stable job.
	Job *Job
}

// DeploymentAllocHealthRequest is used to set the health of a set of
// allocations as part of a deployment.
type DeploymentAllocHealthRequest struct {
	DeploymentID string

	// Marks these allocations as healthy, allow further allocations
	// to be rolled.
	HealthyAllocationIDs []string

	// Any unhealthy allocations fail the deployment
	UnhealthyAllocationIDs []string

	WriteRequest
}

// ApplyDeploymentAllocHealthRequest is used to apply an alloc health request via Raft
type ApplyDeploymentAllocHealthRequest struct {
	DeploymentAllocHealthRequest

	// Timestamp is the timestamp to use when setting the allocations health.
	Timestamp time.Time

	// An optional field to update the status of a deployment
	DeploymentUpdate *DeploymentStatusUpdate

	// Job is used to optionally upsert a job. This is used when setting the
	// allocation health results in a deployment failure and the deployment
	// auto-reverts to the latest stable job.
	Job *Job

	// An optional evaluation to create after promoting the canaries
	Eval *Evaluation
}

// DeploymentPromoteRequest is used to promote task groups in a deployment
type DeploymentPromoteRequest struct {
	DeploymentID string

	// All is to promote all task groups
	All bool

	// Groups is used to set the promotion status per task group
	Groups []string

	WriteRequest
}

// ApplyDeploymentPromoteRequest is used to apply a promotion request via Raft
type ApplyDeploymentPromoteRequest struct {
	DeploymentPromoteRequest

	// An optional evaluation to create after promoting the canaries
	Eval *Evaluation
}

// DeploymentPauseRequest is used to pause a deployment
type DeploymentPauseRequest struct {
	DeploymentID string

	// Pause sets the pause status
	Pause bool

	WriteRequest
}

// DeploymentSpecificRequest is used to make a request specific to a particular
// deployment
type DeploymentSpecificRequest struct {
	DeploymentID string
	QueryOptions
}

// DeploymentFailRequest is used to fail a particular deployment
type DeploymentFailRequest struct {
	DeploymentID string
	WriteRequest
}

// SingleDeploymentResponse is used to respond with a single deployment
type SingleDeploymentResponse struct {
	Deployment *Deployment
	QueryMeta
}

// GenericResponse is used to respond to a request where no
// specific response information is needed.
type GenericResponse struct {
	WriteMeta
}

// VersionResponse is used for the Status.Version response
type VersionResponse struct {
	Build    string
	Versions map[string]int
	QueryMeta
}

// JobRegisterResponse is used to respond to a job registration
type JobRegisterResponse struct {
	EvalID          string
	EvalCreateIndex uint64
	JobModifyIndex  uint64

	// Warnings contains any warnings about the given job. These may include
	// deprecation warnings.
	Warnings string

	QueryMeta
}

// JobDeregisterResponse is used to respond to a job deregistration
type JobDeregisterResponse struct {
	EvalID          string
	EvalCreateIndex uint64
	JobModifyIndex  uint64
	QueryMeta
}

// JobBatchDeregisterResponse is used to respond to a batch job deregistration
type JobBatchDeregisterResponse struct {
	// JobEvals maps the job to its created evaluation
	JobEvals map[NamespacedID]string
	QueryMeta
}

// JobValidateResponse is the response from validate request
type JobValidateResponse struct {
	// DriverConfigValidated indicates whether the agent validated the driver
	// config
	DriverConfigValidated bool

	// ValidationErrors is a list of validation errors
	ValidationErrors []string

	// Error is a string version of any error that may have occurred
	Error string

	// Warnings contains any warnings about the given job. These may include
	// deprecation warnings.
	Warnings string
}

// NodeUpdateResponse is used to respond to a node update
type NodeUpdateResponse struct {
	HeartbeatTTL    time.Duration
	EvalIDs         []string
	EvalCreateIndex uint64
	NodeModifyIndex uint64

	// LeaderRPCAddr is the RPC address of the current Raft Leader.  If
	// empty, the current Nomad Server is in the minority of a partition.
	LeaderRPCAddr string

	// NumNodes is the number of Nomad nodes attached to this quorum of
	// Nomad Servers at the time of the response.  This value can
	// fluctuate based on the health of the cluster between heartbeats.
	NumNodes int32

	// Servers is the full list of known Nomad servers in the local
	// region.
	Servers []*NodeServerInfo

	QueryMeta
}

// NodeDrainUpdateResponse is used to respond to a node drain update
type NodeDrainUpdateResponse struct {
	NodeModifyIndex uint64
	EvalIDs         []string
	EvalCreateIndex uint64
	WriteMeta
}

// NodeEligibilityUpdateResponse is used to respond to a node eligibility update
type NodeEligibilityUpdateResponse struct {
	NodeModifyIndex uint64
	EvalIDs         []string
	EvalCreateIndex uint64
	WriteMeta
}

// NodeAllocsResponse is used to return allocs for a single node
type NodeAllocsResponse struct {
	Allocs []*Allocation
	QueryMeta
}

// NodeClientAllocsResponse is used to return allocs meta data for a single node
type NodeClientAllocsResponse struct {
	Allocs map[string]uint64

	// MigrateTokens are used when ACLs are enabled to allow cross node,
	// authenticated access to sticky volumes
	MigrateTokens map[string]string

	QueryMeta
}

// SingleNodeResponse is used to return a single node
type SingleNodeResponse struct {
	Node *Node
	QueryMeta
}

// NodeListResponse is used for a list request
type NodeListResponse struct {
	Nodes []*NodeListStub
	QueryMeta
}

// SingleJobResponse is used to return a single job
type SingleJobResponse struct {
	Job *Job
	QueryMeta
}

// JobSummaryResponse is used to return a single job summary
type JobSummaryResponse struct {
	JobSummary *JobSummary
	QueryMeta
}

type JobDispatchResponse struct {
	DispatchedJobID string
	EvalID          string
	EvalCreateIndex uint64
	JobCreateIndex  uint64
	WriteMeta
}

// JobListResponse is used for a list request
type JobListResponse struct {
	Jobs []*JobListStub
	QueryMeta
}

// JobVersionsRequest is used to get a jobs versions
type JobVersionsRequest struct {
	JobID string
	Diffs bool
	QueryOptions
}

// JobVersionsResponse is used for a job get versions request
type JobVersionsResponse struct {
	Versions []*Job
	Diffs    []*JobDiff
	QueryMeta
}

// JobPlanResponse is used to respond to a job plan request
type JobPlanResponse struct {
	// Annotations stores annotations explaining decisions the scheduler made.
	Annotations *PlanAnnotations

	// FailedTGAllocs is the placement failures per task group.
	FailedTGAllocs map[string]*AllocMetric

	// JobModifyIndex is the modification index of the job. The value can be
	// used when running `nomad run` to ensure that the Job wasn’t modified
	// since the last plan. If the job is being created, the value is zero.
	JobModifyIndex uint64

	// CreatedEvals is the set of evaluations created by the scheduler. The
	// reasons for this can be rolling-updates or blocked evals.
	CreatedEvals []*Evaluation

	// Diff contains the diff of the job and annotations on whether the change
	// causes an in-place update or create/destroy
	Diff *JobDiff

	// NextPeriodicLaunch is the time duration till the job would be launched if
	// submitted.
	NextPeriodicLaunch time.Time

	// Warnings contains any warnings about the given job. These may include
	// deprecation warnings.
	Warnings string

	WriteMeta
}

// SingleAllocResponse is used to return a single allocation
type SingleAllocResponse struct {
	Alloc *Allocation
	QueryMeta
}

// AllocsGetResponse is used to return a set of allocations
type AllocsGetResponse struct {
	Allocs []*Allocation
	QueryMeta
}

// JobAllocationsResponse is used to return the allocations for a job
type JobAllocationsResponse struct {
	Allocations []*AllocListStub
	QueryMeta
}

// JobEvaluationsResponse is used to return the evaluations for a job
type JobEvaluationsResponse struct {
	Evaluations []*Evaluation
	QueryMeta
}

// SingleEvalResponse is used to return a single evaluation
type SingleEvalResponse struct {
	Eval *Evaluation
	QueryMeta
}

// EvalDequeueResponse is used to return from a dequeue
type EvalDequeueResponse struct {
	Eval  *Evaluation
	Token string

	// WaitIndex is the Raft index the worker should wait until invoking the
	// scheduler.
	WaitIndex uint64

	QueryMeta
}

// GetWaitIndex is used to retrieve the Raft index in which state should be at
// or beyond before invoking the scheduler.
func (e *EvalDequeueResponse) GetWaitIndex() uint64 {
	// Prefer the wait index sent. This will be populated on all responses from
	// 0.7.0 and above
	if e.WaitIndex != 0 {
		return e.WaitIndex
	} else if e.Eval != nil {
		return e.Eval.ModifyIndex
	}

	// This should never happen
	return 1
}

// PlanResponse is used to return from a PlanRequest
type PlanResponse struct {
	Result *PlanResult
	WriteMeta
}

// AllocListResponse is used for a list request
type AllocListResponse struct {
	Allocations []*AllocListStub
	QueryMeta
}

// DeploymentListResponse is used for a list request
type DeploymentListResponse struct {
	Deployments []*Deployment
	QueryMeta
}

// EvalListResponse is used for a list request
type EvalListResponse struct {
	Evaluations []*Evaluation
	QueryMeta
}

// EvalAllocationsResponse is used to return the allocations for an evaluation
type EvalAllocationsResponse struct {
	Allocations []*AllocListStub
	QueryMeta
}

// PeriodicForceResponse is used to respond to a periodic job force launch
type PeriodicForceResponse struct {
	EvalID          string
	EvalCreateIndex uint64
	WriteMeta
}

// DeploymentUpdateResponse is used to respond to a deployment change. The
// response will include the modify index of the deployment as well as details
// of any triggered evaluation.
type DeploymentUpdateResponse struct {
	EvalID                string
	EvalCreateIndex       uint64
	DeploymentModifyIndex uint64

	// RevertedJobVersion is the version the job was reverted to. If unset, the
	// job wasn't reverted
	RevertedJobVersion *uint64

	WriteMeta
}

// NodeConnQueryResponse is used to respond to a query of whether a server has
// a connection to a specific Node
type NodeConnQueryResponse struct {
	// Connected indicates whether a connection to the Client exists
	Connected bool

	// Established marks the time at which the connection was established
	Established time.Time

	QueryMeta
}

// EmitNodeEventsRequest is a request to update the node events source
// with a new client-side event
type EmitNodeEventsRequest struct {
	// NodeEvents are a map where the key is a node id, and value is a list of
	// events for that node
	NodeEvents map[string][]*NodeEvent

	WriteRequest
}

// EmitNodeEventsResponse is a response to the client about the status of
// the node event source update.
type EmitNodeEventsResponse struct {
	Index uint64
	WriteMeta
}

const (
	NodeEventSubsystemDrain     = "Drain"
	NodeEventSubsystemDriver    = "Driver"
	NodeEventSubsystemHeartbeat = "Heartbeat"
	NodeEventSubsystemCluster   = "Cluster"
)

// NodeEvent is a single unit representing a node’s state change
type NodeEvent struct {
	Message     string
	Subsystem   string
	Details     map[string]string
	Timestamp   time.Time
	CreateIndex uint64
}

func (ne *NodeEvent) String() string {
	var details []string
	for k, v := range ne.Details {
		details = append(details, fmt.Sprintf("%s: %s", k, v))
	}

	return fmt.Sprintf("Message: %s, Subsystem: %s, Details: %s, Timestamp: %s", ne.Message, ne.Subsystem, strings.Join(details, ","), ne.Timestamp.String())
}

func (ne *NodeEvent) Copy() *NodeEvent {
	c := new(NodeEvent)
	*c = *ne
	c.Details = helper.CopyMapStringString(ne.Details)
	return c
}

// NewNodeEvent generates a new node event storing the current time as the
// timestamp
func NewNodeEvent() *NodeEvent {
	return &NodeEvent{Timestamp: time.Now()}
}

// SetMessage is used to set the message on the node event
func (ne *NodeEvent) SetMessage(msg string) *NodeEvent {
	ne.Message = msg
	return ne
}

// SetSubsystem is used to set the subsystem on the node event
func (ne *NodeEvent) SetSubsystem(sys string) *NodeEvent {
	ne.Subsystem = sys
	return ne
}

// SetTimestamp is used to set the timestamp on the node event
func (ne *NodeEvent) SetTimestamp(ts time.Time) *NodeEvent {
	ne.Timestamp = ts
	return ne
}

// AddDetail is used to add a detail to the node event
func (ne *NodeEvent) AddDetail(k, v string) *NodeEvent {
	if ne.Details == nil {
		ne.Details = make(map[string]string, 1)
	}
	ne.Details[k] = v
	return ne
}

const (
	NodeStatusInit  = "initializing"
	NodeStatusReady = "ready"
	NodeStatusDown  = "down"
)

// ShouldDrainNode checks if a given node status should trigger an
// evaluation. Some states don't require any further action.
func ShouldDrainNode(status string) bool {
	switch status {
	case NodeStatusInit, NodeStatusReady:
		return false
	case NodeStatusDown:
		return true
	default:
		panic(fmt.Sprintf("unhandled node status %s", status))
	}
}

// ValidNodeStatus is used to check if a node status is valid
func ValidNodeStatus(status string) bool {
	switch status {
	case NodeStatusInit, NodeStatusReady, NodeStatusDown:
		return true
	default:
		return false
	}
}

const (
	// NodeSchedulingEligible and Ineligible marks the node as eligible or not,
	// respectively, for receiving allocations. This is orthoginal to the node
	// status being ready.
	NodeSchedulingEligible   = "eligible"
	NodeSchedulingIneligible = "ineligible"
)

// DrainSpec describes a Node's desired drain behavior.
type DrainSpec struct {
	// Deadline is the duration after StartTime when the remaining
	// allocations on a draining Node should be told to stop.
	Deadline time.Duration

	// IgnoreSystemJobs allows systems jobs to remain on the node even though it
	// has been marked for draining.
	IgnoreSystemJobs bool
}

// DrainStrategy describes a Node's drain behavior.
type DrainStrategy struct {
	// DrainSpec is the user declared drain specification
	DrainSpec

	// ForceDeadline is the deadline time for the drain after which drains will
	// be forced
	ForceDeadline time.Time
}

func (d *DrainStrategy) Copy() *DrainStrategy {
	if d == nil {
		return nil
	}

	nd := new(DrainStrategy)
	*nd = *d
	return nd
}

// DeadlineTime returns a boolean whether the drain strategy allows an infinite
// duration or otherwise the deadline time. The force drain is captured by the
// deadline time being in the past.
func (d *DrainStrategy) DeadlineTime() (infinite bool, deadline time.Time) {
	// Treat the nil case as a force drain so during an upgrade where a node may
	// not have a drain strategy but has Drain set to true, it is treated as a
	// force to mimick old behavior.
	if d == nil {
		return false, time.Time{}
	}

	ns := d.Deadline.Nanoseconds()
	switch {
	case ns < 0: // Force
		return false, time.Time{}
	case ns == 0: // Infinite
		return true, time.Time{}
	default:
		return false, d.ForceDeadline
	}
}

func (d *DrainStrategy) Equal(o *DrainStrategy) bool {
	if d == nil && o == nil {
		return true
	} else if o != nil && d == nil {
		return false
	} else if d != nil && o == nil {
		return false
	}

	// Compare values
	if d.ForceDeadline != o.ForceDeadline {
		return false
	} else if d.Deadline != o.Deadline {
		return false
	} else if d.IgnoreSystemJobs != o.IgnoreSystemJobs {
		return false
	}

	return true
}

// Node is a representation of a schedulable client node
type Node struct {
	// ID is a unique identifier for the node. It can be constructed
	// by doing a concatenation of the Name and Datacenter as a simple
	// approach. Alternatively a UUID may be used.
	ID string

	// SecretID is an ID that is only known by the Node and the set of Servers.
	// It is not accessible via the API and is used to authenticate nodes
	// conducting privileged activities.
	SecretID string

	// Datacenter for this node
	Datacenter string

	// Node name
	Name string

	// HTTPAddr is the address on which the Nomad client is listening for http
	// requests
	HTTPAddr string

	// TLSEnabled indicates if the Agent has TLS enabled for the HTTP API
	TLSEnabled bool

	// Attributes is an arbitrary set of key/value
	// data that can be used for constraints. Examples
	// include "kernel.name=linux", "arch=386", "driver.docker=1",
	// "docker.runtime=1.8.3"
	Attributes map[string]string

	// Resources is the available resources on the client.
	// For example 'cpu=2' 'memory=2048'
	Resources *Resources

	// Reserved is the set of resources that are reserved,
	// and should be subtracted from the total resources for
	// the purposes of scheduling. This may be provide certain
	// high-watermark tolerances or because of external schedulers
	// consuming resources.
	Reserved *Resources

	// Links are used to 'link' this client to external
	// systems. For example 'consul=foo.dc1' 'aws=i-83212'
	// 'ami=ami-123'
	Links map[string]string

	// Meta is used to associate arbitrary metadata with this
	// client. This is opaque to Nomad.
	Meta map[string]string

	// NodeClass is an opaque identifier used to group nodes
	// together for the purpose of determining scheduling pressure.
	NodeClass string

	// ComputedClass is a unique id that identifies nodes with a common set of
	// attributes and capabilities.
	ComputedClass string

	// COMPAT: Remove in Nomad 0.9
	// Drain is controlled by the servers, and not the client.
	// If true, no jobs will be scheduled to this node, and existing
	// allocations will be drained. Superceded by DrainStrategy in Nomad
	// 0.8 but kept for backward compat.
	Drain bool

	// DrainStrategy determines the node's draining behavior. Will be nil
	// when Drain=false.
	DrainStrategy *DrainStrategy

	// SchedulingEligibility determines whether this node will receive new
	// placements.
	SchedulingEligibility string

	// Status of this node
	Status string

	// StatusDescription is meant to provide more human useful information
	StatusDescription string

	// StatusUpdatedAt is the time stamp at which the state of the node was
	// updated
	StatusUpdatedAt int64

	// Events is the most recent set of events generated for the node,
	// retaining only MaxRetainedNodeEvents number at a time
	Events []*NodeEvent

	// Drivers is a map of driver names to current driver information
	Drivers map[string]*DriverInfo

	// Raft Indexes
	CreateIndex uint64
	ModifyIndex uint64
}

// Ready returns true if the node is ready for running allocations
func (n *Node) Ready() bool {
	// Drain is checked directly to support pre-0.8 Node data
	return n.Status == NodeStatusReady && !n.Drain && n.SchedulingEligibility == NodeSchedulingEligible
}

func (n *Node) Canonicalize() {
	if n == nil {
		return
	}

	// COMPAT Remove in 0.10
	// In v0.8.0 we introduced scheduling eligibility, so we need to set it for
	// upgrading nodes
	if n.SchedulingEligibility == "" {
		if n.Drain {
			n.SchedulingEligibility = NodeSchedulingIneligible
		} else {
			n.SchedulingEligibility = NodeSchedulingEligible
		}
	}
}

func (n *Node) Copy() *Node {
	if n == nil {
		return nil
	}
	nn := new(Node)
	*nn = *n
	nn.Attributes = helper.CopyMapStringString(nn.Attributes)
	nn.Resources = nn.Resources.Copy()
	nn.Reserved = nn.Reserved.Copy()
	nn.Links = helper.CopyMapStringString(nn.Links)
	nn.Meta = helper.CopyMapStringString(nn.Meta)
	nn.Events = copyNodeEvents(n.Events)
	nn.DrainStrategy = nn.DrainStrategy.Copy()
	nn.Drivers = copyNodeDrivers(n.Drivers)
	return nn
}

// copyNodeEvents is a helper to copy a list of NodeEvent's
func copyNodeEvents(events []*NodeEvent) []*NodeEvent {
	l := len(events)
	if l == 0 {
		return nil
	}

	c := make([]*NodeEvent, l)
	for i, event := range events {
		c[i] = event.Copy()
	}
	return c
}

// copyNodeDrivers is a helper to copy a map of DriverInfo
func copyNodeDrivers(drivers map[string]*DriverInfo) map[string]*DriverInfo {
	l := len(drivers)
	if l == 0 {
		return nil
	}

	c := make(map[string]*DriverInfo, l)
	for driver, info := range drivers {
		c[driver] = info.Copy()
	}
	return c
}

// TerminalStatus returns if the current status is terminal and
// will no longer transition.
func (n *Node) TerminalStatus() bool {
	switch n.Status {
	case NodeStatusDown:
		return true
	default:
		return false
	}
}

// Stub returns a summarized version of the node
func (n *Node) Stub() *NodeListStub {

	addr, _, _ := net.SplitHostPort(n.HTTPAddr)

	return &NodeListStub{
		Address:    addr,
		ID:         n.ID,
		Datacenter: n.Datacenter,
		Name:       n.Name,
		NodeClass:  n.NodeClass,
		Version:    n.Attributes["nomad.version"],
		Drain:      n.Drain,
		SchedulingEligibility: n.SchedulingEligibility,
		Status:                n.Status,
		StatusDescription:     n.StatusDescription,
		Drivers:               n.Drivers,
		CreateIndex:           n.CreateIndex,
		ModifyIndex:           n.ModifyIndex,
	}
}

// NodeListStub is used to return a subset of job information
// for the job list
type NodeListStub struct {
	Address               string
	ID                    string
	Datacenter            string
	Name                  string
	NodeClass             string
	Version               string
	Drain                 bool
	SchedulingEligibility string
	Status                string
	StatusDescription     string
	Drivers               map[string]*DriverInfo
	CreateIndex           uint64
	ModifyIndex           uint64
}

// Networks defined for a task on the Resources struct.
type Networks []*NetworkResource

// Port assignment and IP for the given label or empty values.
func (ns Networks) Port(label string) (string, int) {
	for _, n := range ns {
		for _, p := range n.ReservedPorts {
			if p.Label == label {
				return n.IP, p.Value
			}
		}
		for _, p := range n.DynamicPorts {
			if p.Label == label {
				return n.IP, p.Value
			}
		}
	}
	return "", 0
}

// Resources is used to define the resources available
// on a client
type Resources struct {
	CPU      int
	MemoryMB int
	DiskMB   int
	IOPS     int
	Networks Networks
}

const (
	BytesInMegabyte = 1024 * 1024
)

// DefaultResources is a small resources object that contains the
// default resources requests that we will provide to an object.
// ---  THIS FUNCTION IS REPLICATED IN api/resources.go and should
// be kept in sync.
func DefaultResources() *Resources {
	return &Resources{
		CPU:      100,
		MemoryMB: 300,
		IOPS:     0,
	}
}

// MinResources is a small resources object that contains the
// absolute minimum resources that we will provide to an object.
// This should not be confused with the defaults which are
// provided in Canonicalize() ---  THIS FUNCTION IS REPLICATED IN
// api/resources.go and should be kept in sync.
func MinResources() *Resources {
	return &Resources{
		CPU:      20,
		MemoryMB: 10,
		IOPS:     0,
	}
}

// DiskInBytes returns the amount of disk resources in bytes.
func (r *Resources) DiskInBytes() int64 {
	return int64(r.DiskMB * BytesInMegabyte)
}

// Merge merges this resource with another resource.
func (r *Resources) Merge(other *Resources) {
	if other.CPU != 0 {
		r.CPU = other.CPU
	}
	if other.MemoryMB != 0 {
		r.MemoryMB = other.MemoryMB
	}
	if other.DiskMB != 0 {
		r.DiskMB = other.DiskMB
	}
	if other.IOPS != 0 {
		r.IOPS = other.IOPS
	}
	if len(other.Networks) != 0 {
		r.Networks = other.Networks
	}
}

func (r *Resources) Canonicalize() {
	// Ensure that an empty and nil slices are treated the same to avoid scheduling
	// problems since we use reflect DeepEquals.
	if len(r.Networks) == 0 {
		r.Networks = nil
	}

	for _, n := range r.Networks {
		n.Canonicalize()
	}
}

// MeetsMinResources returns an error if the resources specified are less than
// the minimum allowed.
// This is based on the minimums defined in the Resources type
func (r *Resources) MeetsMinResources() error {
	var mErr multierror.Error
	minResources := MinResources()
	if r.CPU < minResources.CPU {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("minimum CPU value is %d; got %d", minResources.CPU, r.CPU))
	}
	if r.MemoryMB < minResources.MemoryMB {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("minimum MemoryMB value is %d; got %d", minResources.MemoryMB, r.MemoryMB))
	}
	if r.IOPS < minResources.IOPS {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("minimum IOPS value is %d; got %d", minResources.IOPS, r.IOPS))
	}
	for i, n := range r.Networks {
		if err := n.MeetsMinResources(); err != nil {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("network resource at index %d failed: %v", i, err))
		}
	}

	return mErr.ErrorOrNil()
}

// Copy returns a deep copy of the resources
func (r *Resources) Copy() *Resources {
	if r == nil {
		return nil
	}
	newR := new(Resources)
	*newR = *r
	if r.Networks != nil {
		n := len(r.Networks)
		newR.Networks = make([]*NetworkResource, n)
		for i := 0; i < n; i++ {
			newR.Networks[i] = r.Networks[i].Copy()
		}
	}
	return newR
}

// NetIndex finds the matching net index using device name
func (r *Resources) NetIndex(n *NetworkResource) int {
	for idx, net := range r.Networks {
		if net.Device == n.Device {
			return idx
		}
	}
	return -1
}

// Superset checks if one set of resources is a superset
// of another. This ignores network resources, and the NetworkIndex
// should be used for that.
func (r *Resources) Superset(other *Resources) (bool, string) {
	if r.CPU < other.CPU {
		return false, "cpu"
	}
	if r.MemoryMB < other.MemoryMB {
		return false, "memory"
	}
	if r.DiskMB < other.DiskMB {
		return false, "disk"
	}
	if r.IOPS < other.IOPS {
		return false, "iops"
	}
	return true, ""
}

// Add adds the resources of the delta to this, potentially
// returning an error if not possible.
func (r *Resources) Add(delta *Resources) error {
	if delta == nil {
		return nil
	}
	r.CPU += delta.CPU
	r.MemoryMB += delta.MemoryMB
	r.DiskMB += delta.DiskMB
	r.IOPS += delta.IOPS

	for _, n := range delta.Networks {
		// Find the matching interface by IP or CIDR
		idx := r.NetIndex(n)
		if idx == -1 {
			r.Networks = append(r.Networks, n.Copy())
		} else {
			r.Networks[idx].Add(n)
		}
	}
	return nil
}

func (r *Resources) GoString() string {
	return fmt.Sprintf("*%#v", *r)
}

type Port struct {
	Label string
	Value int
}

// NetworkResource is used to represent available network
// resources
type NetworkResource struct {
	Device        string // Name of the device
	CIDR          string // CIDR block of addresses
	IP            string // Host IP address
	MBits         int    // Throughput
	ReservedPorts []Port // Host Reserved ports
	DynamicPorts  []Port // Host Dynamically assigned ports
}

func (nr *NetworkResource) Equals(other *NetworkResource) bool {
	if nr.Device != other.Device {
		return false
	}

	if nr.CIDR != other.CIDR {
		return false
	}

	if nr.IP != other.IP {
		return false
	}

	if nr.MBits != other.MBits {
		return false
	}

	if len(nr.ReservedPorts) != len(other.ReservedPorts) {
		return false
	}

	for i, port := range nr.ReservedPorts {
		if len(other.ReservedPorts) <= i {
			return false
		}
		if port != other.ReservedPorts[i] {
			return false
		}
	}

	if len(nr.DynamicPorts) != len(other.DynamicPorts) {
		return false
	}
	for i, port := range nr.DynamicPorts {
		if len(other.DynamicPorts) <= i {
			return false
		}
		if port != other.DynamicPorts[i] {
			return false
		}
	}
	return true
}

func (n *NetworkResource) Canonicalize() {
	// Ensure that an empty and nil slices are treated the same to avoid scheduling
	// problems since we use reflect DeepEquals.
	if len(n.ReservedPorts) == 0 {
		n.ReservedPorts = nil
	}
	if len(n.DynamicPorts) == 0 {
		n.DynamicPorts = nil
	}
}

// MeetsMinResources returns an error if the resources specified are less than
// the minimum allowed.
func (n *NetworkResource) MeetsMinResources() error {
	var mErr multierror.Error
	if n.MBits < 1 {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("minimum MBits value is 1; got %d", n.MBits))
	}
	return mErr.ErrorOrNil()
}

// Copy returns a deep copy of the network resource
func (n *NetworkResource) Copy() *NetworkResource {
	if n == nil {
		return nil
	}
	newR := new(NetworkResource)
	*newR = *n
	if n.ReservedPorts != nil {
		newR.ReservedPorts = make([]Port, len(n.ReservedPorts))
		copy(newR.ReservedPorts, n.ReservedPorts)
	}
	if n.DynamicPorts != nil {
		newR.DynamicPorts = make([]Port, len(n.DynamicPorts))
		copy(newR.DynamicPorts, n.DynamicPorts)
	}
	return newR
}

// Add adds the resources of the delta to this, potentially
// returning an error if not possible.
func (n *NetworkResource) Add(delta *NetworkResource) {
	if len(delta.ReservedPorts) > 0 {
		n.ReservedPorts = append(n.ReservedPorts, delta.ReservedPorts...)
	}
	n.MBits += delta.MBits
	n.DynamicPorts = append(n.DynamicPorts, delta.DynamicPorts...)
}

func (n *NetworkResource) GoString() string {
	return fmt.Sprintf("*%#v", *n)
}

// PortLabels returns a map of port labels to their assigned host ports.
func (n *NetworkResource) PortLabels() map[string]int {
	num := len(n.ReservedPorts) + len(n.DynamicPorts)
	labelValues := make(map[string]int, num)
	for _, port := range n.ReservedPorts {
		labelValues[port.Label] = port.Value
	}
	for _, port := range n.DynamicPorts {
		labelValues[port.Label] = port.Value
	}
	return labelValues
}

const (
	// JobTypeNomad is reserved for internal system tasks and is
	// always handled by the CoreScheduler.
	JobTypeCore    = "_core"
	JobTypeService = "service"
	JobTypeBatch   = "batch"
	JobTypeSystem  = "system"
)

const (
	JobStatusPending = "pending" // Pending means the job is waiting on scheduling
	JobStatusRunning = "running" // Running means the job has non-terminal allocations
	JobStatusDead    = "dead"    // Dead means all evaluation's and allocations are terminal
)

const (
	// JobMinPriority is the minimum allowed priority
	JobMinPriority = 1

	// JobDefaultPriority is the default priority if not
	// not specified.
	JobDefaultPriority = 50

	// JobMaxPriority is the maximum allowed priority
	JobMaxPriority = 100

	// Ensure CoreJobPriority is higher than any user
	// specified job so that it gets priority. This is important
	// for the system to remain healthy.
	CoreJobPriority = JobMaxPriority * 2

	// JobTrackedVersions is the number of historic job versions that are
	// kept.
	JobTrackedVersions = 6
)

// Job is the scope of a scheduling request to Nomad. It is the largest
// scoped object, and is a named collection of task groups. Each task group
// is further composed of tasks. A task group (TG) is the unit of scheduling
// however.
type Job struct {
	// Stop marks whether the user has stopped the job. A stopped job will
	// have all created allocations stopped and acts as a way to stop a job
	// without purging it from the system. This allows existing allocs to be
	// queried and the job to be inspected as it is being killed.
	Stop bool

	// Region is the Nomad region that handles scheduling this job
	Region string

	// Namespace is the namespace the job is submitted into.
	Namespace string

	// ID is a unique identifier for the job per region. It can be
	// specified hierarchically like LineOfBiz/OrgName/Team/Project
	ID string

	// ParentID is the unique identifier of the job that spawned this job.
	ParentID string

	// Name is the logical name of the job used to refer to it. This is unique
	// per region, but not unique globally.
	Name string

	// Type is used to control various behaviors about the job. Most jobs
	// are service jobs, meaning they are expected to be long lived.
	// Some jobs are batch oriented meaning they run and then terminate.
	// This can be extended in the future to support custom schedulers.
	Type string

	// Priority is used to control scheduling importance and if this job
	// can preempt other jobs.
	Priority int

	// AllAtOnce is used to control if incremental scheduling of task groups
	// is allowed or if we must do a gang scheduling of the entire job. This
	// can slow down larger jobs if resources are not available.
	AllAtOnce bool

	// Datacenters contains all the datacenters this job is allowed to span
	Datacenters []string

	// Constraints can be specified at a job level and apply to
	// all the task groups and tasks.
	Constraints []*Constraint

	// TaskGroups are the collections of task groups that this job needs
	// to run. Each task group is an atomic unit of scheduling and placement.
	TaskGroups []*TaskGroup

	// COMPAT: Remove in 0.7.0. Stagger is deprecated in 0.6.0.
	Update UpdateStrategy

	// Periodic is used to define the interval the job is run at.
	Periodic *PeriodicConfig

	// ParameterizedJob is used to specify the job as a parameterized job
	// for dispatching.
	ParameterizedJob *ParameterizedJobConfig

	// Dispatched is used to identify if the Job has been dispatched from a
	// parameterized job.
	Dispatched bool

	// Payload is the payload supplied when the job was dispatched.
	Payload []byte

	// Meta is used to associate arbitrary metadata with this
	// job. This is opaque to Nomad.
	Meta map[string]string

	// VaultToken is the Vault token that proves the submitter of the job has
	// access to the specified Vault policies. This field is only used to
	// transfer the token and is not stored after Job submission.
	VaultToken string

	// Job status
	Status string

	// StatusDescription is meant to provide more human useful information
	StatusDescription string

	// Stable marks a job as stable. Stability is only defined on "service" and
	// "system" jobs. The stability of a job will be set automatically as part
	// of a deployment and can be manually set via APIs.
	Stable bool

	// Version is a monotonically increasing version number that is incremented
	// on each job register.
	Version uint64

	// SubmitTime is the time at which the job was submitted as a UnixNano in
	// UTC
	SubmitTime int64

	// Raft Indexes
	CreateIndex    uint64
	ModifyIndex    uint64
	JobModifyIndex uint64
}

// NamespacedID returns the namespaced id useful for logging
func (j *Job) NamespacedID() *NamespacedID {
	return &NamespacedID{
		ID:        j.ID,
		Namespace: j.Namespace,
	}
}

// Canonicalize is used to canonicalize fields in the Job. This should be called
// when registering a Job. A set of warnings are returned if the job was changed
// in anyway that the user should be made aware of.
func (j *Job) Canonicalize() (warnings error) {
	if j == nil {
		return nil
	}

	var mErr multierror.Error
	// Ensure that an empty and nil map are treated the same to avoid scheduling
	// problems since we use reflect DeepEquals.
	if len(j.Meta) == 0 {
		j.Meta = nil
	}

	// Ensure the job is in a namespace.
	if j.Namespace == "" {
		j.Namespace = DefaultNamespace
	}

	for _, tg := range j.TaskGroups {
		tg.Canonicalize(j)
	}

	if j.ParameterizedJob != nil {
		j.ParameterizedJob.Canonicalize()
	}

	if j.Periodic != nil {
		j.Periodic.Canonicalize()
	}

	return mErr.ErrorOrNil()
}

// Copy returns a deep copy of the Job. It is expected that callers use recover.
// This job can panic if the deep copy failed as it uses reflection.
func (j *Job) Copy() *Job {
	if j == nil {
		return nil
	}
	nj := new(Job)
	*nj = *j
	nj.Datacenters = helper.CopySliceString(nj.Datacenters)
	nj.Constraints = CopySliceConstraints(nj.Constraints)

	if j.TaskGroups != nil {
		tgs := make([]*TaskGroup, len(nj.TaskGroups))
		for i, tg := range nj.TaskGroups {
			tgs[i] = tg.Copy()
		}
		nj.TaskGroups = tgs
	}

	nj.Periodic = nj.Periodic.Copy()
	nj.Meta = helper.CopyMapStringString(nj.Meta)
	nj.ParameterizedJob = nj.ParameterizedJob.Copy()
	return nj
}

// Validate is used to sanity check a job input
func (j *Job) Validate() error {
	var mErr multierror.Error

	if j.Region == "" {
		mErr.Errors = append(mErr.Errors, errors.New("Missing job region"))
	}
	if j.ID == "" {
		mErr.Errors = append(mErr.Errors, errors.New("Missing job ID"))
	} else if strings.Contains(j.ID, " ") {
		mErr.Errors = append(mErr.Errors, errors.New("Job ID contains a space"))
	}
	if j.Name == "" {
		mErr.Errors = append(mErr.Errors, errors.New("Missing job name"))
	}
	if j.Namespace == "" {
		mErr.Errors = append(mErr.Errors, errors.New("Job must be in a namespace"))
	}
	switch j.Type {
	case JobTypeCore, JobTypeService, JobTypeBatch, JobTypeSystem:
	case "":
		mErr.Errors = append(mErr.Errors, errors.New("Missing job type"))
	default:
		mErr.Errors = append(mErr.Errors, fmt.Errorf("Invalid job type: %q", j.Type))
	}
	if j.Priority < JobMinPriority || j.Priority > JobMaxPriority {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("Job priority must be between [%d, %d]", JobMinPriority, JobMaxPriority))
	}
	if len(j.Datacenters) == 0 {
		mErr.Errors = append(mErr.Errors, errors.New("Missing job datacenters"))
	}
	if len(j.TaskGroups) == 0 {
		mErr.Errors = append(mErr.Errors, errors.New("Missing job task groups"))
	}
	for idx, constr := range j.Constraints {
		if err := constr.Validate(); err != nil {
			outer := fmt.Errorf("Constraint %d validation failed: %s", idx+1, err)
			mErr.Errors = append(mErr.Errors, outer)
		}
	}

	// Check for duplicate task groups
	taskGroups := make(map[string]int)
	for idx, tg := range j.TaskGroups {
		if tg.Name == "" {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Job task group %d missing name", idx+1))
		} else if existing, ok := taskGroups[tg.Name]; ok {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Job task group %d redefines '%s' from group %d", idx+1, tg.Name, existing+1))
		} else {
			taskGroups[tg.Name] = idx
		}

		if j.Type == "system" && tg.Count > 1 {
			mErr.Errors = append(mErr.Errors,
				fmt.Errorf("Job task group %s has count %d. Count cannot exceed 1 with system scheduler",
					tg.Name, tg.Count))
		}
	}

	// Validate the task group
	for _, tg := range j.TaskGroups {
		if err := tg.Validate(j); err != nil {
			outer := fmt.Errorf("Task group %s validation failed: %v", tg.Name, err)
			mErr.Errors = append(mErr.Errors, outer)
		}
	}

	// Validate periodic is only used with batch jobs.
	if j.IsPeriodic() && j.Periodic.Enabled {
		if j.Type != JobTypeBatch {
			mErr.Errors = append(mErr.Errors,
				fmt.Errorf("Periodic can only be used with %q scheduler", JobTypeBatch))
		}

		if err := j.Periodic.Validate(); err != nil {
			mErr.Errors = append(mErr.Errors, err)
		}
	}

	if j.IsParameterized() {
		if j.Type != JobTypeBatch {
			mErr.Errors = append(mErr.Errors,
				fmt.Errorf("Parameterized job can only be used with %q scheduler", JobTypeBatch))
		}

		if err := j.ParameterizedJob.Validate(); err != nil {
			mErr.Errors = append(mErr.Errors, err)
		}
	}

	return mErr.ErrorOrNil()
}

// Warnings returns a list of warnings that may be from dubious settings or
// deprecation warnings.
func (j *Job) Warnings() error {
	var mErr multierror.Error

	// Check the groups
	for _, tg := range j.TaskGroups {
		if err := tg.Warnings(j); err != nil {
			outer := fmt.Errorf("Group %q has warnings: %v", tg.Name, err)
			mErr.Errors = append(mErr.Errors, outer)
		}
	}

	return mErr.ErrorOrNil()
}

// LookupTaskGroup finds a task group by name
func (j *Job) LookupTaskGroup(name string) *TaskGroup {
	for _, tg := range j.TaskGroups {
		if tg.Name == name {
			return tg
		}
	}
	return nil
}

// CombinedTaskMeta takes a TaskGroup and Task name and returns the combined
// meta data for the task. When joining Job, Group and Task Meta, the precedence
// is by deepest scope (Task > Group > Job).
func (j *Job) CombinedTaskMeta(groupName, taskName string) map[string]string {
	group := j.LookupTaskGroup(groupName)
	if group == nil {
		return nil
	}

	task := group.LookupTask(taskName)
	if task == nil {
		return nil
	}

	meta := helper.CopyMapStringString(task.Meta)
	if meta == nil {
		meta = make(map[string]string, len(group.Meta)+len(j.Meta))
	}

	// Add the group specific meta
	for k, v := range group.Meta {
		if _, ok := meta[k]; !ok {
			meta[k] = v
		}
	}

	// Add the job specific meta
	for k, v := range j.Meta {
		if _, ok := meta[k]; !ok {
			meta[k] = v
		}
	}

	return meta
}

// Stopped returns if a job is stopped.
func (j *Job) Stopped() bool {
	return j == nil || j.Stop
}

// HasUpdateStrategy returns if any task group in the job has an update strategy
func (j *Job) HasUpdateStrategy() bool {
	for _, tg := range j.TaskGroups {
		if tg.Update != nil {
			return true
		}
	}

	return false
}

// Stub is used to return a summary of the job
func (j *Job) Stub(summary *JobSummary) *JobListStub {
	return &JobListStub{
		ID:                j.ID,
		ParentID:          j.ParentID,
		Name:              j.Name,
		Type:              j.Type,
		Priority:          j.Priority,
		Periodic:          j.IsPeriodic(),
		ParameterizedJob:  j.IsParameterized(),
		Stop:              j.Stop,
		Status:            j.Status,
		StatusDescription: j.StatusDescription,
		CreateIndex:       j.CreateIndex,
		ModifyIndex:       j.ModifyIndex,
		JobModifyIndex:    j.JobModifyIndex,
		SubmitTime:        j.SubmitTime,
		JobSummary:        summary,
	}
}

// IsPeriodic returns whether a job is periodic.
func (j *Job) IsPeriodic() bool {
	return j.Periodic != nil
}

// IsPeriodicActive returns whether the job is an active periodic job that will
// create child jobs
func (j *Job) IsPeriodicActive() bool {
	return j.IsPeriodic() && j.Periodic.Enabled && !j.Stopped() && !j.IsParameterized()
}

// IsParameterized returns whether a job is parameterized job.
func (j *Job) IsParameterized() bool {
	return j.ParameterizedJob != nil && !j.Dispatched
}

// VaultPolicies returns the set of Vault policies per task group, per task
func (j *Job) VaultPolicies() map[string]map[string]*Vault {
	policies := make(map[string]map[string]*Vault, len(j.TaskGroups))

	for _, tg := range j.TaskGroups {
		tgPolicies := make(map[string]*Vault, len(tg.Tasks))

		for _, task := range tg.Tasks {
			if task.Vault == nil {
				continue
			}

			tgPolicies[task.Name] = task.Vault
		}

		if len(tgPolicies) != 0 {
			policies[tg.Name] = tgPolicies
		}
	}

	return policies
}

// RequiredSignals returns a mapping of task groups to tasks to their required
// set of signals
func (j *Job) RequiredSignals() map[string]map[string][]string {
	signals := make(map[string]map[string][]string)

	for _, tg := range j.TaskGroups {
		for _, task := range tg.Tasks {
			// Use this local one as a set
			taskSignals := make(map[string]struct{})

			// Check if the Vault change mode uses signals
			if task.Vault != nil && task.Vault.ChangeMode == VaultChangeModeSignal {
				taskSignals[task.Vault.ChangeSignal] = struct{}{}
			}

			// If a user has specified a KillSignal, add it to required signals
			if task.KillSignal != "" {
				taskSignals[task.KillSignal] = struct{}{}
			}

			// Check if any template change mode uses signals
			for _, t := range task.Templates {
				if t.ChangeMode != TemplateChangeModeSignal {
					continue
				}

				taskSignals[t.ChangeSignal] = struct{}{}
			}

			// Flatten and sort the signals
			l := len(taskSignals)
			if l == 0 {
				continue
			}

			flat := make([]string, 0, l)
			for sig := range taskSignals {
				flat = append(flat, sig)
			}

			sort.Strings(flat)
			tgSignals, ok := signals[tg.Name]
			if !ok {
				tgSignals = make(map[string][]string)
				signals[tg.Name] = tgSignals
			}
			tgSignals[task.Name] = flat
		}

	}

	return signals
}

// SpecChanged determines if the functional specification has changed between
// two job versions.
func (j *Job) SpecChanged(new *Job) bool {
	if j == nil {
		return new != nil
	}

	// Create a copy of the new job
	c := new.Copy()

	// Update the new job so we can do a reflect
	c.Status = j.Status
	c.StatusDescription = j.StatusDescription
	c.Stable = j.Stable
	c.Version = j.Version
	c.CreateIndex = j.CreateIndex
	c.ModifyIndex = j.ModifyIndex
	c.JobModifyIndex = j.JobModifyIndex
	c.SubmitTime = j.SubmitTime

	// Deep equals the jobs
	return !reflect.DeepEqual(j, c)
}

func (j *Job) SetSubmitTime() {
	j.SubmitTime = time.Now().UTC().UnixNano()
}

// JobListStub is used to return a subset of job information
// for the job list
type JobListStub struct {
	ID                string
	ParentID          string
	Name              string
	Type              string
	Priority          int
	Periodic          bool
	ParameterizedJob  bool
	Stop              bool
	Status            string
	StatusDescription string
	JobSummary        *JobSummary
	CreateIndex       uint64
	ModifyIndex       uint64
	JobModifyIndex    uint64
	SubmitTime        int64
}

// JobSummary summarizes the state of the allocations of a job
type JobSummary struct {
	// JobID is the ID of the job the summary is for
	JobID string

	// Namespace is the namespace of the job and its summary
	Namespace string

	// Summary contains the summary per task group for the Job
	Summary map[string]TaskGroupSummary

	// Children contains a summary for the children of this job.
	Children *JobChildrenSummary

	// Raft Indexes
	CreateIndex uint64
	ModifyIndex uint64
}

// Copy returns a new copy of JobSummary
func (js *JobSummary) Copy() *JobSummary {
	newJobSummary := new(JobSummary)
	*newJobSummary = *js
	newTGSummary := make(map[string]TaskGroupSummary, len(js.Summary))
	for k, v := range js.Summary {
		newTGSummary[k] = v
	}
	newJobSummary.Summary = newTGSummary
	newJobSummary.Children = newJobSummary.Children.Copy()
	return newJobSummary
}

// JobChildrenSummary contains the summary of children job statuses
type JobChildrenSummary struct {
	Pending int64
	Running int64
	Dead    int64
}

// Copy returns a new copy of a JobChildrenSummary
func (jc *JobChildrenSummary) Copy() *JobChildrenSummary {
	if jc == nil {
		return nil
	}

	njc := new(JobChildrenSummary)
	*njc = *jc
	return njc
}

// TaskGroup summarizes the state of all the allocations of a particular
// TaskGroup
type TaskGroupSummary struct {
	Queued   int
	Complete int
	Failed   int
	Running  int
	Starting int
	Lost     int
}

const (
	// Checks uses any registered health check state in combination with task
	// states to determine if a allocation is healthy.
	UpdateStrategyHealthCheck_Checks = "checks"

	// TaskStates uses the task states of an allocation to determine if the
	// allocation is healthy.
	UpdateStrategyHealthCheck_TaskStates = "task_states"

	// Manual allows the operator to manually signal to Nomad when an
	// allocations is healthy. This allows more advanced health checking that is
	// outside of the scope of Nomad.
	UpdateStrategyHealthCheck_Manual = "manual"
)

var (
	// DefaultUpdateStrategy provides a baseline that can be used to upgrade
	// jobs with the old policy or for populating field defaults.
	DefaultUpdateStrategy = &UpdateStrategy{
		Stagger:          30 * time.Second,
		MaxParallel:      1,
		HealthCheck:      UpdateStrategyHealthCheck_Checks,
		MinHealthyTime:   10 * time.Second,
		HealthyDeadline:  5 * time.Minute,
		ProgressDeadline: 10 * time.Minute,
		AutoRevert:       false,
		Canary:           0,
	}
)

// UpdateStrategy is used to modify how updates are done
type UpdateStrategy struct {
	// Stagger is used to determine the rate at which allocations are migrated
	// due to down or draining nodes.
	Stagger time.Duration

	// MaxParallel is how many updates can be done in parallel
	MaxParallel int

	// HealthCheck specifies the mechanism in which allocations are marked
	// healthy or unhealthy as part of a deployment.
	HealthCheck string

	// MinHealthyTime is the minimum time an allocation must be in the healthy
	// state before it is marked as healthy, unblocking more allocations to be
	// rolled.
	MinHealthyTime time.Duration

	// HealthyDeadline is the time in which an allocation must be marked as
	// healthy before it is automatically transitioned to unhealthy. This time
	// period doesn't count against the MinHealthyTime.
	HealthyDeadline time.Duration

	// ProgressDeadline is the time in which an allocation as part of the
	// deployment must transition to healthy. If no allocation becomes healthy
	// after the deadline, the deployment is marked as failed. If the deadline
	// is zero, the first failure causes the deployment to fail.
	ProgressDeadline time.Duration

	// AutoRevert declares that if a deployment fails because of unhealthy
	// allocations, there should be an attempt to auto-revert the job to a
	// stable version.
	AutoRevert bool

	// Canary is the number of canaries to deploy when a change to the task
	// group is detected.
	Canary int
}

func (u *UpdateStrategy) Copy() *UpdateStrategy {
	if u == nil {
		return nil
	}

	copy := new(UpdateStrategy)
	*copy = *u
	return copy
}

func (u *UpdateStrategy) Validate() error {
	if u == nil {
		return nil
	}

	var mErr multierror.Error
	switch u.HealthCheck {
	case UpdateStrategyHealthCheck_Checks, UpdateStrategyHealthCheck_TaskStates, UpdateStrategyHealthCheck_Manual:
	default:
		multierror.Append(&mErr, fmt.Errorf("Invalid health check given: %q", u.HealthCheck))
	}

	if u.MaxParallel < 1 {
		multierror.Append(&mErr, fmt.Errorf("Max parallel can not be less than one: %d < 1", u.MaxParallel))
	}
	if u.Canary < 0 {
		multierror.Append(&mErr, fmt.Errorf("Canary count can not be less than zero: %d < 0", u.Canary))
	}
	if u.MinHealthyTime < 0 {
		multierror.Append(&mErr, fmt.Errorf("Minimum healthy time may not be less than zero: %v", u.MinHealthyTime))
	}
	if u.HealthyDeadline <= 0 {
		multierror.Append(&mErr, fmt.Errorf("Healthy deadline must be greater than zero: %v", u.HealthyDeadline))
	}
	if u.ProgressDeadline < 0 {
		multierror.Append(&mErr, fmt.Errorf("Progress deadline must be zero or greater: %v", u.ProgressDeadline))
	}
	if u.MinHealthyTime >= u.HealthyDeadline {
		multierror.Append(&mErr, fmt.Errorf("Minimum healthy time must be less than healthy deadline: %v > %v", u.MinHealthyTime, u.HealthyDeadline))
	}
	if u.ProgressDeadline != 0 && u.HealthyDeadline >= u.ProgressDeadline {
		multierror.Append(&mErr, fmt.Errorf("Healthy deadline must be less than progress deadline: %v > %v", u.HealthyDeadline, u.ProgressDeadline))
	}
	if u.Stagger <= 0 {
		multierror.Append(&mErr, fmt.Errorf("Stagger must be greater than zero: %v", u.Stagger))
	}

	return mErr.ErrorOrNil()
}

// TODO(alexdadgar): Remove once no longer used by the scheduler.
// Rolling returns if a rolling strategy should be used
func (u *UpdateStrategy) Rolling() bool {
	return u.Stagger > 0 && u.MaxParallel > 0
}

const (
	// PeriodicSpecCron is used for a cron spec.
	PeriodicSpecCron = "cron"

	// PeriodicSpecTest is only used by unit tests. It is a sorted, comma
	// separated list of unix timestamps at which to launch.
	PeriodicSpecTest = "_internal_test"
)

// Periodic defines the interval a job should be run at.
type PeriodicConfig struct {
	// Enabled determines if the job should be run periodically.
	Enabled bool

	// Spec specifies the interval the job should be run as. It is parsed based
	// on the SpecType.
	Spec string

	// SpecType defines the format of the spec.
	SpecType string

	// ProhibitOverlap enforces that spawned jobs do not run in parallel.
	ProhibitOverlap bool

	// TimeZone is the user specified string that determines the time zone to
	// launch against. The time zones must be specified from IANA Time Zone
	// database, such as "America/New_York".
	// Reference: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
	// Reference: https://www.iana.org/time-zones
	TimeZone string

	// location is the time zone to evaluate the launch time against
	location *time.Location
}

func (p *PeriodicConfig) Copy() *PeriodicConfig {
	if p == nil {
		return nil
	}
	np := new(PeriodicConfig)
	*np = *p
	return np
}

func (p *PeriodicConfig) Validate() error {
	if !p.Enabled {
		return nil
	}

	var mErr multierror.Error
	if p.Spec == "" {
		multierror.Append(&mErr, fmt.Errorf("Must specify a spec"))
	}

	// Check if we got a valid time zone
	if p.TimeZone != "" {
		if _, err := time.LoadLocation(p.TimeZone); err != nil {
			multierror.Append(&mErr, fmt.Errorf("Invalid time zone %q: %v", p.TimeZone, err))
		}
	}

	switch p.SpecType {
	case PeriodicSpecCron:
		// Validate the cron spec
		if _, err := cronexpr.Parse(p.Spec); err != nil {
			multierror.Append(&mErr, fmt.Errorf("Invalid cron spec %q: %v", p.Spec, err))
		}
	case PeriodicSpecTest:
		// No-op
	default:
		multierror.Append(&mErr, fmt.Errorf("Unknown periodic specification type %q", p.SpecType))
	}

	return mErr.ErrorOrNil()
}

func (p *PeriodicConfig) Canonicalize() {
	// Load the location
	l, err := time.LoadLocation(p.TimeZone)
	if err != nil {
		p.location = time.UTC
	}

	p.location = l
}

// CronParseNext is a helper that parses the next time for the given expression
// but captures any panic that may occur in the underlying library.
func CronParseNext(e *cronexpr.Expression, fromTime time.Time, spec string) (t time.Time, err error) {
	defer func() {
		if recover() != nil {
			t = time.Time{}
			err = fmt.Errorf("failed parsing cron expression: %q", spec)
		}
	}()

	return e.Next(fromTime), nil
}

// Next returns the closest time instant matching the spec that is after the
// passed time. If no matching instance exists, the zero value of time.Time is
// returned. The `time.Location` of the returned value matches that of the
// passed time.
func (p *PeriodicConfig) Next(fromTime time.Time) (time.Time, error) {
	switch p.SpecType {
	case PeriodicSpecCron:
		if e, err := cronexpr.Parse(p.Spec); err == nil {
			return CronParseNext(e, fromTime, p.Spec)
		}
	case PeriodicSpecTest:
		split := strings.Split(p.Spec, ",")
		if len(split) == 1 && split[0] == "" {
			return time.Time{}, nil
		}

		// Parse the times
		times := make([]time.Time, len(split))
		for i, s := range split {
			unix, err := strconv.Atoi(s)
			if err != nil {
				return time.Time{}, nil
			}

			times[i] = time.Unix(int64(unix), 0)
		}

		// Find the next match
		for _, next := range times {
			if fromTime.Before(next) {
				return next, nil
			}
		}
	}

	return time.Time{}, nil
}

// GetLocation returns the location to use for determining the time zone to run
// the periodic job against.
func (p *PeriodicConfig) GetLocation() *time.Location {
	// Jobs pre 0.5.5 will not have this
	if p.location != nil {
		return p.location
	}

	return time.UTC
}

const (
	// PeriodicLaunchSuffix is the string appended to the periodic jobs ID
	// when launching derived instances of it.
	PeriodicLaunchSuffix = "/periodic-"
)

// PeriodicLaunch tracks the last launch time of a periodic job.
type PeriodicLaunch struct {
	ID        string    // ID of the periodic job.
	Namespace string    // Namespace of the periodic job
	Launch    time.Time // The last launch time.

	// Raft Indexes
	CreateIndex uint64
	ModifyIndex uint64
}

const (
	DispatchPayloadForbidden = "forbidden"
	DispatchPayloadOptional  = "optional"
	DispatchPayloadRequired  = "required"

	// DispatchLaunchSuffix is the string appended to the parameterized job's ID
	// when dispatching instances of it.
	DispatchLaunchSuffix = "/dispatch-"
)

// ParameterizedJobConfig is used to configure the parameterized job
type ParameterizedJobConfig struct {
	// Payload configure the payload requirements
	Payload string

	// MetaRequired is metadata keys that must be specified by the dispatcher
	MetaRequired []string

	// MetaOptional is metadata keys that may be specified by the dispatcher
	MetaOptional []string
}

func (d *ParameterizedJobConfig) Validate() error {
	var mErr multierror.Error
	switch d.Payload {
	case DispatchPayloadOptional, DispatchPayloadRequired, DispatchPayloadForbidden:
	default:
		multierror.Append(&mErr, fmt.Errorf("Unknown payload requirement: %q", d.Payload))
	}

	// Check that the meta configurations are disjoint sets
	disjoint, offending := helper.SliceSetDisjoint(d.MetaRequired, d.MetaOptional)
	if !disjoint {
		multierror.Append(&mErr, fmt.Errorf("Required and optional meta keys should be disjoint. Following keys exist in both: %v", offending))
	}

	return mErr.ErrorOrNil()
}

func (d *ParameterizedJobConfig) Canonicalize() {
	if d.Payload == "" {
		d.Payload = DispatchPayloadOptional
	}
}

func (d *ParameterizedJobConfig) Copy() *ParameterizedJobConfig {
	if d == nil {
		return nil
	}
	nd := new(ParameterizedJobConfig)
	*nd = *d
	nd.MetaOptional = helper.CopySliceString(nd.MetaOptional)
	nd.MetaRequired = helper.CopySliceString(nd.MetaRequired)
	return nd
}

// DispatchedID returns an ID appropriate for a job dispatched against a
// particular parameterized job
func DispatchedID(templateID string, t time.Time) string {
	u := uuid.Generate()[:8]
	return fmt.Sprintf("%s%s%d-%s", templateID, DispatchLaunchSuffix, t.Unix(), u)
}

// DispatchPayloadConfig configures how a task gets its input from a job dispatch
type DispatchPayloadConfig struct {
	// File specifies a relative path to where the input data should be written
	File string
}

func (d *DispatchPayloadConfig) Copy() *DispatchPayloadConfig {
	if d == nil {
		return nil
	}
	nd := new(DispatchPayloadConfig)
	*nd = *d
	return nd
}

func (d *DispatchPayloadConfig) Validate() error {
	// Verify the destination doesn't escape
	escaped, err := PathEscapesAllocDir("task/local/", d.File)
	if err != nil {
		return fmt.Errorf("invalid destination path: %v", err)
	} else if escaped {
		return fmt.Errorf("destination escapes allocation directory")
	}

	return nil
}

var (
	DefaultServiceJobRestartPolicy = RestartPolicy{
		Delay:    15 * time.Second,
		Attempts: 2,
		Interval: 30 * time.Minute,
		Mode:     RestartPolicyModeFail,
	}
	DefaultBatchJobRestartPolicy = RestartPolicy{
		Delay:    15 * time.Second,
		Attempts: 3,
		Interval: 24 * time.Hour,
		Mode:     RestartPolicyModeFail,
	}
)

var (
	DefaultServiceJobReschedulePolicy = ReschedulePolicy{
		Delay:         30 * time.Second,
		DelayFunction: "exponential",
		MaxDelay:      1 * time.Hour,
		Unlimited:     true,
	}
	DefaultBatchJobReschedulePolicy = ReschedulePolicy{
		Attempts:      1,
		Interval:      24 * time.Hour,
		Delay:         5 * time.Second,
		DelayFunction: "constant",
	}
)

const (
	// RestartPolicyModeDelay causes an artificial delay till the next interval is
	// reached when the specified attempts have been reached in the interval.
	RestartPolicyModeDelay = "delay"

	// RestartPolicyModeFail causes a job to fail if the specified number of
	// attempts are reached within an interval.
	RestartPolicyModeFail = "fail"

	// RestartPolicyMinInterval is the minimum interval that is accepted for a
	// restart policy.
	RestartPolicyMinInterval = 5 * time.Second

	// ReasonWithinPolicy describes restart events that are within policy
	ReasonWithinPolicy = "Restart within policy"
)

// RestartPolicy configures how Tasks are restarted when they crash or fail.
type RestartPolicy struct {
	// Attempts is the number of restart that will occur in an interval.
	Attempts int

	// Interval is a duration in which we can limit the number of restarts
	// within.
	Interval time.Duration

	// Delay is the time between a failure and a restart.
	Delay time.Duration

	// Mode controls what happens when the task restarts more than attempt times
	// in an interval.
	Mode string
}

func (r *RestartPolicy) Copy() *RestartPolicy {
	if r == nil {
		return nil
	}
	nrp := new(RestartPolicy)
	*nrp = *r
	return nrp
}

func (r *RestartPolicy) Validate() error {
	var mErr multierror.Error
	switch r.Mode {
	case RestartPolicyModeDelay, RestartPolicyModeFail:
	default:
		multierror.Append(&mErr, fmt.Errorf("Unsupported restart mode: %q", r.Mode))
	}

	// Check for ambiguous/confusing settings
	if r.Attempts == 0 && r.Mode != RestartPolicyModeFail {
		multierror.Append(&mErr, fmt.Errorf("Restart policy %q with %d attempts is ambiguous", r.Mode, r.Attempts))
	}

	if r.Interval.Nanoseconds() < RestartPolicyMinInterval.Nanoseconds() {
		multierror.Append(&mErr, fmt.Errorf("Interval can not be less than %v (got %v)", RestartPolicyMinInterval, r.Interval))
	}
	if time.Duration(r.Attempts)*r.Delay > r.Interval {
		multierror.Append(&mErr,
			fmt.Errorf("Nomad can't restart the TaskGroup %v times in an interval of %v with a delay of %v", r.Attempts, r.Interval, r.Delay))
	}
	return mErr.ErrorOrNil()
}

func NewRestartPolicy(jobType string) *RestartPolicy {
	switch jobType {
	case JobTypeService, JobTypeSystem:
		rp := DefaultServiceJobRestartPolicy
		return &rp
	case JobTypeBatch:
		rp := DefaultBatchJobRestartPolicy
		return &rp
	}
	return nil
}

const ReschedulePolicyMinInterval = 15 * time.Second
const ReschedulePolicyMinDelay = 5 * time.Second

var RescheduleDelayFunctions = [...]string{"constant", "exponential", "fibonacci"}

// ReschedulePolicy configures how Tasks are rescheduled  when they crash or fail.
type ReschedulePolicy struct {
	// Attempts limits the number of rescheduling attempts that can occur in an interval.
	Attempts int

	// Interval is a duration in which we can limit the number of reschedule attempts.
	Interval time.Duration

	// Delay is a minimum duration to wait between reschedule attempts.
	// The delay function determines how much subsequent reschedule attempts are delayed by.
	Delay time.Duration

	// DelayFunction determines how the delay progressively changes on subsequent reschedule
	// attempts. Valid values are "exponential", "constant", and "fibonacci".
	DelayFunction string

	// MaxDelay is an upper bound on the delay.
	MaxDelay time.Duration

	// Unlimited allows infinite rescheduling attempts. Only allowed when delay is set
	// between reschedule attempts.
	Unlimited bool
}

func (r *ReschedulePolicy) Copy() *ReschedulePolicy {
	if r == nil {
		return nil
	}
	nrp := new(ReschedulePolicy)
	*nrp = *r
	return nrp
}

func (r *ReschedulePolicy) Enabled() bool {
	enabled := r != nil && (r.Attempts > 0 || r.Unlimited)
	return enabled
}

// Validate uses different criteria to validate the reschedule policy
// Delay must be a minimum of 5 seconds
// Delay Ceiling is ignored if Delay Function is "constant"
// Number of possible attempts is validated, given the interval, delay and delay function
func (r *ReschedulePolicy) Validate() error {
	if !r.Enabled() {
		return nil
	}
	var mErr multierror.Error
	// Check for ambiguous/confusing settings
	if r.Attempts > 0 {
		if r.Interval <= 0 {
			multierror.Append(&mErr, fmt.Errorf("Interval must be a non zero value if Attempts > 0"))
		}
		if r.Unlimited {
			multierror.Append(&mErr, fmt.Errorf("Reschedule Policy with Attempts = %v, Interval = %v, "+
				"and Unlimited = %v is ambiguous", r.Attempts, r.Interval, r.Unlimited))
			multierror.Append(&mErr, errors.New("If Attempts >0, Unlimited cannot also be set to true"))
		}
	}

	delayPreCheck := true
	// Delay should be bigger than the default
	if r.Delay.Nanoseconds() < ReschedulePolicyMinDelay.Nanoseconds() {
		multierror.Append(&mErr, fmt.Errorf("Delay cannot be less than %v (got %v)", ReschedulePolicyMinDelay, r.Delay))
		delayPreCheck = false
	}

	// Must use a valid delay function
	if !isValidDelayFunction(r.DelayFunction) {
		multierror.Append(&mErr, fmt.Errorf("Invalid delay function %q, must be one of %q", r.DelayFunction, RescheduleDelayFunctions))
		delayPreCheck = false
	}

	// Validate MaxDelay if not using linear delay progression
	if r.DelayFunction != "constant" {
		if r.MaxDelay.Nanoseconds() < ReschedulePolicyMinDelay.Nanoseconds() {
			multierror.Append(&mErr, fmt.Errorf("Max Delay cannot be less than %v (got %v)", ReschedulePolicyMinDelay, r.Delay))
			delayPreCheck = false
		}
		if r.MaxDelay < r.Delay {
			multierror.Append(&mErr, fmt.Errorf("Max Delay cannot be less than Delay %v (got %v)", r.Delay, r.MaxDelay))
			delayPreCheck = false
		}

	}

	// Validate Interval and other delay parameters if attempts are limited
	if !r.Unlimited {
		if r.Interval.Nanoseconds() < ReschedulePolicyMinInterval.Nanoseconds() {
			multierror.Append(&mErr, fmt.Errorf("Interval cannot be less than %v (got %v)", ReschedulePolicyMinInterval, r.Interval))
		}
		if !delayPreCheck {
			// We can't cross validate the rest of the delay params if delayPreCheck fails, so return early
			return mErr.ErrorOrNil()
		}
		crossValidationErr := r.validateDelayParams()
		if crossValidationErr != nil {
			multierror.Append(&mErr, crossValidationErr)
		}
	}
	return mErr.ErrorOrNil()
}

func isValidDelayFunction(delayFunc string) bool {
	for _, value := range RescheduleDelayFunctions {
		if value == delayFunc {
			return true
		}
	}
	return false
}

func (r *ReschedulePolicy) validateDelayParams() error {
	ok, possibleAttempts, recommendedInterval := r.viableAttempts()
	if ok {
		return nil
	}
	var mErr multierror.Error
	if r.DelayFunction == "constant" {
		multierror.Append(&mErr, fmt.Errorf("Nomad can only make %v attempts in %v with initial delay %v and "+
			"delay function %q", possibleAttempts, r.Interval, r.Delay, r.DelayFunction))
	} else {
		multierror.Append(&mErr, fmt.Errorf("Nomad can only make %v attempts in %v with initial delay %v, "+
			"delay function %q, and delay ceiling %v", possibleAttempts, r.Interval, r.Delay, r.DelayFunction, r.MaxDelay))
	}
	multierror.Append(&mErr, fmt.Errorf("Set the interval to at least %v to accommodate %v attempts", recommendedInterval.Round(time.Second), r.Attempts))
	return mErr.ErrorOrNil()
}

func (r *ReschedulePolicy) viableAttempts() (bool, int, time.Duration) {
	var possibleAttempts int
	var recommendedInterval time.Duration
	valid := true
	switch r.DelayFunction {
	case "constant":
		recommendedInterval = time.Duration(r.Attempts) * r.Delay
		if r.Interval < recommendedInterval {
			possibleAttempts = int(r.Interval / r.Delay)
			valid = false
		}
	case "exponential":
		for i := 0; i < r.Attempts; i++ {
			nextDelay := time.Duration(math.Pow(2, float64(i))) * r.Delay
			if nextDelay > r.MaxDelay {
				nextDelay = r.MaxDelay
				recommendedInterval += nextDelay
			} else {
				recommendedInterval = nextDelay
			}
			if recommendedInterval < r.Interval {
				possibleAttempts++
			}
		}
		if possibleAttempts < r.Attempts {
			valid = false
		}
	case "fibonacci":
		var slots []time.Duration
		slots = append(slots, r.Delay)
		slots = append(slots, r.Delay)
		reachedCeiling := false
		for i := 2; i < r.Attempts; i++ {
			var nextDelay time.Duration
			if reachedCeiling {
				//switch to linear
				nextDelay = slots[i-1] + r.MaxDelay
			} else {
				nextDelay = slots[i-1] + slots[i-2]
				if nextDelay > r.MaxDelay {
					nextDelay = r.MaxDelay
					reachedCeiling = true
				}
			}
			slots = append(slots, nextDelay)
		}
		recommendedInterval = slots[len(slots)-1]
		if r.Interval < recommendedInterval {
			valid = false
			// calculate possible attempts
			for i := 0; i < len(slots); i++ {
				if slots[i] > r.Interval {
					possibleAttempts = i
					break
				}
			}
		}
	default:
		return false, 0, 0
	}
	if possibleAttempts < 0 { // can happen if delay is bigger than interval
		possibleAttempts = 0
	}
	return valid, possibleAttempts, recommendedInterval
}

func NewReschedulePolicy(jobType string) *ReschedulePolicy {
	switch jobType {
	case JobTypeService:
		rp := DefaultServiceJobReschedulePolicy
		return &rp
	case JobTypeBatch:
		rp := DefaultBatchJobReschedulePolicy
		return &rp
	}
	return nil
}

const (
	MigrateStrategyHealthChecks = "checks"
	MigrateStrategyHealthStates = "task_states"
)

type MigrateStrategy struct {
	MaxParallel     int
	HealthCheck     string
	MinHealthyTime  time.Duration
	HealthyDeadline time.Duration
}

// DefaultMigrateStrategy is used for backwards compat with pre-0.8 Allocations
// that lack an update strategy.
//
// This function should match its counterpart in api/tasks.go
func DefaultMigrateStrategy() *MigrateStrategy {
	return &MigrateStrategy{
		MaxParallel:     1,
		HealthCheck:     MigrateStrategyHealthChecks,
		MinHealthyTime:  10 * time.Second,
		HealthyDeadline: 5 * time.Minute,
	}
}

func (m *MigrateStrategy) Validate() error {
	var mErr multierror.Error

	if m.MaxParallel < 0 {
		multierror.Append(&mErr, fmt.Errorf("MaxParallel must be >= 0 but found %d", m.MaxParallel))
	}

	switch m.HealthCheck {
	case MigrateStrategyHealthChecks, MigrateStrategyHealthStates:
		// ok
	case "":
		if m.MaxParallel > 0 {
			multierror.Append(&mErr, fmt.Errorf("Missing HealthCheck"))
		}
	default:
		multierror.Append(&mErr, fmt.Errorf("Invalid HealthCheck: %q", m.HealthCheck))
	}

	if m.MinHealthyTime < 0 {
		multierror.Append(&mErr, fmt.Errorf("MinHealthyTime is %s and must be >= 0", m.MinHealthyTime))
	}

	if m.HealthyDeadline < 0 {
		multierror.Append(&mErr, fmt.Errorf("HealthyDeadline is %s and must be >= 0", m.HealthyDeadline))
	}

	if m.MinHealthyTime > m.HealthyDeadline {
		multierror.Append(&mErr, fmt.Errorf("MinHealthyTime must be less than HealthyDeadline"))
	}

	return mErr.ErrorOrNil()
}

// TaskGroup is an atomic unit of placement. Each task group belongs to
// a job and may contain any number of tasks. A task group support running
// in many replicas using the same configuration..
type TaskGroup struct {
	// Name of the task group
	Name string

	// Count is the number of replicas of this task group that should
	// be scheduled.
	Count int

	// Update is used to control the update strategy for this task group
	Update *UpdateStrategy

	// Migrate is used to control the migration strategy for this task group
	Migrate *MigrateStrategy

	// Constraints can be specified at a task group level and apply to
	// all the tasks contained.
	Constraints []*Constraint

	//RestartPolicy of a TaskGroup
	RestartPolicy *RestartPolicy

	// Tasks are the collection of tasks that this task group needs to run
	Tasks []*Task

	// EphemeralDisk is the disk resources that the task group requests
	EphemeralDisk *EphemeralDisk

	// Meta is used to associate arbitrary metadata with this
	// task group. This is opaque to Nomad.
	Meta map[string]string

	// ReschedulePolicy is used to configure how the scheduler should
	// retry failed allocations.
	ReschedulePolicy *ReschedulePolicy
}

func (tg *TaskGroup) Copy() *TaskGroup {
	if tg == nil {
		return nil
	}
	ntg := new(TaskGroup)
	*ntg = *tg
	ntg.Update = ntg.Update.Copy()
	ntg.Constraints = CopySliceConstraints(ntg.Constraints)
	ntg.RestartPolicy = ntg.RestartPolicy.Copy()
	ntg.ReschedulePolicy = ntg.ReschedulePolicy.Copy()

	if tg.Tasks != nil {
		tasks := make([]*Task, len(ntg.Tasks))
		for i, t := range ntg.Tasks {
			tasks[i] = t.Copy()
		}
		ntg.Tasks = tasks
	}

	ntg.Meta = helper.CopyMapStringString(ntg.Meta)

	if tg.EphemeralDisk != nil {
		ntg.EphemeralDisk = tg.EphemeralDisk.Copy()
	}
	return ntg
}

// Canonicalize is used to canonicalize fields in the TaskGroup.
func (tg *TaskGroup) Canonicalize(job *Job) {
	// Ensure that an empty and nil map are treated the same to avoid scheduling
	// problems since we use reflect DeepEquals.
	if len(tg.Meta) == 0 {
		tg.Meta = nil
	}

	// Set the default restart policy.
	if tg.RestartPolicy == nil {
		tg.RestartPolicy = NewRestartPolicy(job.Type)
	}

	if tg.ReschedulePolicy == nil {
		tg.ReschedulePolicy = NewReschedulePolicy(job.Type)
	}

	// Canonicalize Migrate for service jobs
	if job.Type == JobTypeService && tg.Migrate == nil {
		tg.Migrate = DefaultMigrateStrategy()
	}

	// Set a default ephemeral disk object if the user has not requested for one
	if tg.EphemeralDisk == nil {
		tg.EphemeralDisk = DefaultEphemeralDisk()
	}

	for _, task := range tg.Tasks {
		task.Canonicalize(job, tg)
	}

	// Add up the disk resources to EphemeralDisk. This is done so that users
	// are not required to move their disk attribute from resources to
	// EphemeralDisk section of the job spec in Nomad 0.5
	// COMPAT 0.4.1 -> 0.5
	// Remove in 0.6
	var diskMB int
	for _, task := range tg.Tasks {
		diskMB += task.Resources.DiskMB
	}
	if diskMB > 0 {
		tg.EphemeralDisk.SizeMB = diskMB
	}
}

// Validate is used to sanity check a task group
func (tg *TaskGroup) Validate(j *Job) error {
	var mErr multierror.Error
	if tg.Name == "" {
		mErr.Errors = append(mErr.Errors, errors.New("Missing task group name"))
	}
	if tg.Count < 0 {
		mErr.Errors = append(mErr.Errors, errors.New("Task group count can't be negative"))
	}
	if len(tg.Tasks) == 0 {
		mErr.Errors = append(mErr.Errors, errors.New("Missing tasks for task group"))
	}
	for idx, constr := range tg.Constraints {
		if err := constr.Validate(); err != nil {
			outer := fmt.Errorf("Constraint %d validation failed: %s", idx+1, err)
			mErr.Errors = append(mErr.Errors, outer)
		}
	}

	if tg.RestartPolicy != nil {
		if err := tg.RestartPolicy.Validate(); err != nil {
			mErr.Errors = append(mErr.Errors, err)
		}
	} else {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("Task Group %v should have a restart policy", tg.Name))
	}

	if j.Type == JobTypeSystem {
		if tg.ReschedulePolicy != nil {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("System jobs should not have a reschedule policy"))
		}
	} else {
		if tg.ReschedulePolicy != nil {
			if err := tg.ReschedulePolicy.Validate(); err != nil {
				mErr.Errors = append(mErr.Errors, err)
			}
		} else {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Task Group %v should have a reschedule policy", tg.Name))
		}
	}

	if tg.EphemeralDisk != nil {
		if err := tg.EphemeralDisk.Validate(); err != nil {
			mErr.Errors = append(mErr.Errors, err)
		}
	} else {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("Task Group %v should have an ephemeral disk object", tg.Name))
	}

	// Validate the update strategy
	if u := tg.Update; u != nil {
		switch j.Type {
		case JobTypeService, JobTypeSystem:
		default:
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Job type %q does not allow update block", j.Type))
		}
		if err := u.Validate(); err != nil {
			mErr.Errors = append(mErr.Errors, err)
		}
	}

	// Validate the migration strategy
	switch j.Type {
	case JobTypeService:
		if tg.Migrate != nil {
			if err := tg.Migrate.Validate(); err != nil {
				mErr.Errors = append(mErr.Errors, err)
			}
		}
	default:
		if tg.Migrate != nil {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Job type %q does not allow migrate block", j.Type))
		}
	}

	// Check for duplicate tasks, that there is only leader task if any,
	// and no duplicated static ports
	tasks := make(map[string]int)
	staticPorts := make(map[int]string)
	leaderTasks := 0
	for idx, task := range tg.Tasks {
		if task.Name == "" {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Task %d missing name", idx+1))
		} else if existing, ok := tasks[task.Name]; ok {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Task %d redefines '%s' from task %d", idx+1, task.Name, existing+1))
		} else {
			tasks[task.Name] = idx
		}

		if task.Leader {
			leaderTasks++
		}

		if task.Resources == nil {
			continue
		}

		for _, net := range task.Resources.Networks {
			for _, port := range net.ReservedPorts {
				if other, ok := staticPorts[port.Value]; ok {
					err := fmt.Errorf("Static port %d already reserved by %s", port.Value, other)
					mErr.Errors = append(mErr.Errors, err)
				} else {
					staticPorts[port.Value] = fmt.Sprintf("%s:%s", task.Name, port.Label)
				}
			}
		}
	}

	if leaderTasks > 1 {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("Only one task may be marked as leader"))
	}

	// Validate the tasks
	for _, task := range tg.Tasks {
		if err := task.Validate(tg.EphemeralDisk); err != nil {
			outer := fmt.Errorf("Task %s validation failed: %v", task.Name, err)
			mErr.Errors = append(mErr.Errors, outer)
		}
	}
	return mErr.ErrorOrNil()
}

// Warnings returns a list of warnings that may be from dubious settings or
// deprecation warnings.
func (tg *TaskGroup) Warnings(j *Job) error {
	var mErr multierror.Error

	// Validate the update strategy
	if u := tg.Update; u != nil {
		// Check the counts are appropriate
		if u.MaxParallel > tg.Count {
			mErr.Errors = append(mErr.Errors,
				fmt.Errorf("Update max parallel count is greater than task group count (%d > %d). "+
					"A destructive change would result in the simultaneous replacement of all allocations.", u.MaxParallel, tg.Count))
		}
	}

	return mErr.ErrorOrNil()
}

// LookupTask finds a task by name
func (tg *TaskGroup) LookupTask(name string) *Task {
	for _, t := range tg.Tasks {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func (tg *TaskGroup) GoString() string {
	return fmt.Sprintf("*%#v", *tg)
}

// CombinedResources returns the combined resources for the task group
func (tg *TaskGroup) CombinedResources() *Resources {
	r := &Resources{
		DiskMB: tg.EphemeralDisk.SizeMB,
	}
	for _, task := range tg.Tasks {
		r.Add(task.Resources)
	}
	return r
}

// CheckRestart describes if and when a task should be restarted based on
// failing health checks.
type CheckRestart struct {
	Limit          int           // Restart task after this many unhealthy intervals
	Grace          time.Duration // Grace time to give tasks after starting to get healthy
	IgnoreWarnings bool          // If true treat checks in `warning` as passing
}

func (c *CheckRestart) Copy() *CheckRestart {
	if c == nil {
		return nil
	}

	nc := new(CheckRestart)
	*nc = *c
	return nc
}

func (c *CheckRestart) Validate() error {
	if c == nil {
		return nil
	}

	var mErr multierror.Error
	if c.Limit < 0 {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("limit must be greater than or equal to 0 but found %d", c.Limit))
	}

	if c.Grace < 0 {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("grace period must be greater than or equal to 0 but found %d", c.Grace))
	}

	return mErr.ErrorOrNil()
}

const (
	ServiceCheckHTTP   = "http"
	ServiceCheckTCP    = "tcp"
	ServiceCheckScript = "script"
	ServiceCheckGRPC   = "grpc"

	// minCheckInterval is the minimum check interval permitted.  Consul
	// currently has its MinInterval set to 1s.  Mirror that here for
	// consistency.
	minCheckInterval = 1 * time.Second

	// minCheckTimeout is the minimum check timeout permitted for Consul
	// script TTL checks.
	minCheckTimeout = 1 * time.Second
)

// The ServiceCheck data model represents the consul health check that
// Nomad registers for a Task
type ServiceCheck struct {
	Name          string              // Name of the check, defaults to id
	Type          string              // Type of the check - tcp, http, docker and script
	Command       string              // Command is the command to run for script checks
	Args          []string            // Args is a list of arguments for script checks
	Path          string              // path of the health check url for http type check
	Protocol      string              // Protocol to use if check is http, defaults to http
	PortLabel     string              // The port to use for tcp/http checks
	AddressMode   string              // 'host' to use host ip:port or 'driver' to use driver's
	Interval      time.Duration       // Interval of the check
	Timeout       time.Duration       // Timeout of the response from the check before consul fails the check
	InitialStatus string              // Initial status of the check
	TLSSkipVerify bool                // Skip TLS verification when Protocol=https
	Method        string              // HTTP Method to use (GET by default)
	Header        map[string][]string // HTTP Headers for Consul to set when making HTTP checks
	CheckRestart  *CheckRestart       // If and when a task should be restarted based on checks
	GRPCService   string              // Service for GRPC checks
	GRPCUseTLS    bool                // Whether or not to use TLS for GRPC checks
}

func (sc *ServiceCheck) Copy() *ServiceCheck {
	if sc == nil {
		return nil
	}
	nsc := new(ServiceCheck)
	*nsc = *sc
	nsc.Args = helper.CopySliceString(sc.Args)
	nsc.Header = helper.CopyMapStringSliceString(sc.Header)
	nsc.CheckRestart = sc.CheckRestart.Copy()
	return nsc
}

func (sc *ServiceCheck) Canonicalize(serviceName string) {
	// Ensure empty maps/slices are treated as null to avoid scheduling
	// issues when using DeepEquals.
	if len(sc.Args) == 0 {
		sc.Args = nil
	}

	if len(sc.Header) == 0 {
		sc.Header = nil
	} else {
		for k, v := range sc.Header {
			if len(v) == 0 {
				sc.Header[k] = nil
			}
		}
	}

	if sc.Name == "" {
		sc.Name = fmt.Sprintf("service: %q check", serviceName)
	}
}

// validate a Service's ServiceCheck
func (sc *ServiceCheck) validate() error {
	// Validate Type
	switch strings.ToLower(sc.Type) {
	case ServiceCheckGRPC:
	case ServiceCheckTCP:
	case ServiceCheckHTTP:
		if sc.Path == "" {
			return fmt.Errorf("http type must have a valid http path")
		}
		url, err := url.Parse(sc.Path)
		if err != nil {
			return fmt.Errorf("http type must have a valid http path")
		}
		if url.IsAbs() {
			return fmt.Errorf("http type must have a relative http path")
		}

	case ServiceCheckScript:
		if sc.Command == "" {
			return fmt.Errorf("script type must have a valid script path")
		}

	default:
		return fmt.Errorf(`invalid type (%+q), must be one of "http", "tcp", or "script" type`, sc.Type)
	}

	// Validate interval and timeout
	if sc.Interval == 0 {
		return fmt.Errorf("missing required value interval. Interval cannot be less than %v", minCheckInterval)
	} else if sc.Interval < minCheckInterval {
		return fmt.Errorf("interval (%v) cannot be lower than %v", sc.Interval, minCheckInterval)
	}

	if sc.Timeout == 0 {
		return fmt.Errorf("missing required value timeout. Timeout cannot be less than %v", minCheckInterval)
	} else if sc.Timeout < minCheckTimeout {
		return fmt.Errorf("timeout (%v) is lower than required minimum timeout %v", sc.Timeout, minCheckInterval)
	}

	// Validate InitialStatus
	switch sc.InitialStatus {
	case "":
	case api.HealthPassing:
	case api.HealthWarning:
	case api.HealthCritical:
	default:
		return fmt.Errorf(`invalid initial check state (%s), must be one of %q, %q, %q or empty`, sc.InitialStatus, api.HealthPassing, api.HealthWarning, api.HealthCritical)

	}

	// Validate AddressMode
	switch sc.AddressMode {
	case "", AddressModeHost, AddressModeDriver:
		// Ok
	case AddressModeAuto:
		return fmt.Errorf("invalid address_mode %q - %s only valid for services", sc.AddressMode, AddressModeAuto)
	default:
		return fmt.Errorf("invalid address_mode %q", sc.AddressMode)
	}

	return sc.CheckRestart.Validate()
}

// RequiresPort returns whether the service check requires the task has a port.
func (sc *ServiceCheck) RequiresPort() bool {
	switch sc.Type {
	case ServiceCheckGRPC, ServiceCheckHTTP, ServiceCheckTCP:
		return true
	default:
		return false
	}
}

// TriggersRestarts returns true if this check should be watched and trigger a restart
// on failure.
func (sc *ServiceCheck) TriggersRestarts() bool {
	return sc.CheckRestart != nil && sc.CheckRestart.Limit > 0
}

// Hash all ServiceCheck fields and the check's corresponding service ID to
// create an identifier. The identifier is not guaranteed to be unique as if
// the PortLabel is blank, the Service's PortLabel will be used after Hash is
// called.
func (sc *ServiceCheck) Hash(serviceID string) string {
	h := sha1.New()
	io.WriteString(h, serviceID)
	io.WriteString(h, sc.Name)
	io.WriteString(h, sc.Type)
	io.WriteString(h, sc.Command)
	io.WriteString(h, strings.Join(sc.Args, ""))
	io.WriteString(h, sc.Path)
	io.WriteString(h, sc.Protocol)
	io.WriteString(h, sc.PortLabel)
	io.WriteString(h, sc.Interval.String())
	io.WriteString(h, sc.Timeout.String())
	io.WriteString(h, sc.Method)
	// Only include TLSSkipVerify if set to maintain ID stability with Nomad <0.6
	if sc.TLSSkipVerify {
		io.WriteString(h, "true")
	}

	// Since map iteration order isn't stable we need to write k/v pairs to
	// a slice and sort it before hashing.
	if len(sc.Header) > 0 {
		headers := make([]string, 0, len(sc.Header))
		for k, v := range sc.Header {
			headers = append(headers, k+strings.Join(v, ""))
		}
		sort.Strings(headers)
		io.WriteString(h, strings.Join(headers, ""))
	}

	// Only include AddressMode if set to maintain ID stability with Nomad <0.7.1
	if len(sc.AddressMode) > 0 {
		io.WriteString(h, sc.AddressMode)
	}

	// Only include GRPC if set to maintain ID stability with Nomad <0.8.4
	if sc.GRPCService != "" {
		io.WriteString(h, sc.GRPCService)
	}
	if sc.GRPCUseTLS {
		io.WriteString(h, "true")
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

const (
	AddressModeAuto   = "auto"
	AddressModeHost   = "host"
	AddressModeDriver = "driver"
)

// Service represents a Consul service definition in Nomad
type Service struct {
	// Name of the service registered with Consul. Consul defaults the
	// Name to ServiceID if not specified.  The Name if specified is used
	// as one of the seed values when generating a Consul ServiceID.
	Name string

	// PortLabel is either the numeric port number or the `host:port`.
	// To specify the port number using the host's Consul Advertise
	// address, specify an empty host in the PortLabel (e.g. `:port`).
	PortLabel string

	// AddressMode specifies whether or not to use the host ip:port for
	// this service.
	AddressMode string

	Tags       []string        // List of tags for the service
	CanaryTags []string        // List of tags for the service when it is a canary
	Checks     []*ServiceCheck // List of checks associated with the service
}

func (s *Service) Copy() *Service {
	if s == nil {
		return nil
	}
	ns := new(Service)
	*ns = *s
	ns.Tags = helper.CopySliceString(ns.Tags)
	ns.CanaryTags = helper.CopySliceString(ns.CanaryTags)

	if s.Checks != nil {
		checks := make([]*ServiceCheck, len(ns.Checks))
		for i, c := range ns.Checks {
			checks[i] = c.Copy()
		}
		ns.Checks = checks
	}

	return ns
}

// Canonicalize interpolates values of Job, Task Group and Task in the Service
// Name. This also generates check names, service id and check ids.
func (s *Service) Canonicalize(job string, taskGroup string, task string) {
	// Ensure empty lists are treated as null to avoid scheduler issues when
	// using DeepEquals
	if len(s.Tags) == 0 {
		s.Tags = nil
	}
	if len(s.CanaryTags) == 0 {
		s.CanaryTags = nil
	}
	if len(s.Checks) == 0 {
		s.Checks = nil
	}

	s.Name = args.ReplaceEnv(s.Name, map[string]string{
		"JOB":       job,
		"TASKGROUP": taskGroup,
		"TASK":      task,
		"BASE":      fmt.Sprintf("%s-%s-%s", job, taskGroup, task),
	},
	)

	for _, check := range s.Checks {
		check.Canonicalize(s.Name)
	}
}

// Validate checks if the Check definition is valid
func (s *Service) Validate() error {
	var mErr multierror.Error

	// Ensure the service name is valid per the below RFCs but make an exception
	// for our interpolation syntax by first stripping any environment variables from the name

	serviceNameStripped := args.ReplaceEnvWithPlaceHolder(s.Name, "ENV-VAR")

	if err := s.ValidateName(serviceNameStripped); err != nil {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("service name must be valid per RFC 1123 and can contain only alphanumeric characters or dashes: %q", s.Name))
	}

	switch s.AddressMode {
	case "", AddressModeAuto, AddressModeHost, AddressModeDriver:
		// OK
	default:
		mErr.Errors = append(mErr.Errors, fmt.Errorf("service address_mode must be %q, %q, or %q; not %q", AddressModeAuto, AddressModeHost, AddressModeDriver, s.AddressMode))
	}

	for _, c := range s.Checks {
		if s.PortLabel == "" && c.PortLabel == "" && c.RequiresPort() {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("check %s invalid: check requires a port but neither check nor service %+q have a port", c.Name, s.Name))
			continue
		}

		if err := c.validate(); err != nil {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("check %s invalid: %v", c.Name, err))
		}
	}

	return mErr.ErrorOrNil()
}

// ValidateName checks if the services Name is valid and should be called after
// the name has been interpolated
func (s *Service) ValidateName(name string) error {
	// Ensure the service name is valid per RFC-952 §1
	// (https://tools.ietf.org/html/rfc952), RFC-1123 §2.1
	// (https://tools.ietf.org/html/rfc1123), and RFC-2782
	// (https://tools.ietf.org/html/rfc2782).
	re := regexp.MustCompile(`^(?i:[a-z0-9]|[a-z0-9][a-z0-9\-]{0,61}[a-z0-9])$`)
	if !re.MatchString(name) {
		return fmt.Errorf("service name must be valid per RFC 1123 and can contain only alphanumeric characters or dashes and must be no longer than 63 characters: %q", name)
	}
	return nil
}

// Hash returns a base32 encoded hash of a Service's contents excluding checks
// as they're hashed independently.
func (s *Service) Hash(allocID, taskName string, canary bool) string {
	h := sha1.New()
	io.WriteString(h, allocID)
	io.WriteString(h, taskName)
	io.WriteString(h, s.Name)
	io.WriteString(h, s.PortLabel)
	io.WriteString(h, s.AddressMode)
	for _, tag := range s.Tags {
		io.WriteString(h, tag)
	}
	for _, tag := range s.CanaryTags {
		io.WriteString(h, tag)
	}

	// Vary ID on whether or not CanaryTags will be used
	if canary {
		h.Write([]byte("Canary"))
	}

	// Base32 is used for encoding the hash as sha1 hashes can always be
	// encoded without padding, only 4 bytes larger than base64, and saves
	// 8 bytes vs hex. Since these hashes are used in Consul URLs it's nice
	// to have a reasonably compact URL-safe representation.
	return b32.EncodeToString(h.Sum(nil))
}

const (
	// DefaultKillTimeout is the default timeout between signaling a task it
	// will be killed and killing it.
	DefaultKillTimeout = 5 * time.Second
)

// LogConfig provides configuration for log rotation
type LogConfig struct {
	MaxFiles      int
	MaxFileSizeMB int
}

// DefaultLogConfig returns the default LogConfig values.
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		MaxFiles:      10,
		MaxFileSizeMB: 10,
	}
}

// Validate returns an error if the log config specified are less than
// the minimum allowed.
func (l *LogConfig) Validate() error {
	var mErr multierror.Error
	if l.MaxFiles < 1 {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("minimum number of files is 1; got %d", l.MaxFiles))
	}
	if l.MaxFileSizeMB < 1 {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("minimum file size is 1MB; got %d", l.MaxFileSizeMB))
	}
	return mErr.ErrorOrNil()
}

// Task is a single process typically that is executed as part of a task group.
type Task struct {
	// Name of the task
	Name string

	// Driver is used to control which driver is used
	Driver string

	// User is used to determine which user will run the task. It defaults to
	// the same user the Nomad client is being run as.
	User string

	// Config is provided to the driver to initialize
	Config map[string]interface{}

	// Map of environment variables to be used by the driver
	Env map[string]string

	// List of service definitions exposed by the Task
	Services []*Service

	// Vault is used to define the set of Vault policies that this task should
	// have access to.
	Vault *Vault

	// Templates are the set of templates to be rendered for the task.
	Templates []*Template

	// Constraints can be specified at a task level and apply only to
	// the particular task.
	Constraints []*Constraint

	// Resources is the resources needed by this task
	Resources *Resources

	// DispatchPayload configures how the task retrieves its input from a dispatch
	DispatchPayload *DispatchPayloadConfig

	// Meta is used to associate arbitrary metadata with this
	// task. This is opaque to Nomad.
	Meta map[string]string

	// KillTimeout is the time between signaling a task that it will be
	// killed and killing it.
	KillTimeout time.Duration

	// LogConfig provides configuration for log rotation
	LogConfig *LogConfig

	// Artifacts is a list of artifacts to download and extract before running
	// the task.
	Artifacts []*TaskArtifact

	// Leader marks the task as the leader within the group. When the leader
	// task exits, other tasks will be gracefully terminated.
	Leader bool

	// ShutdownDelay is the duration of the delay between deregistering a
	// task from Consul and sending it a signal to shutdown. See #2441
	ShutdownDelay time.Duration

	// The kill signal to use for the task. This is an optional specification,

	// KillSignal is the kill signal to use for the task. This is an optional
	// specification and defaults to SIGINT
	KillSignal string
}

func (t *Task) Copy() *Task {
	if t == nil {
		return nil
	}
	nt := new(Task)
	*nt = *t
	nt.Env = helper.CopyMapStringString(nt.Env)

	if t.Services != nil {
		services := make([]*Service, len(nt.Services))
		for i, s := range nt.Services {
			services[i] = s.Copy()
		}
		nt.Services = services
	}

	nt.Constraints = CopySliceConstraints(nt.Constraints)

	nt.Vault = nt.Vault.Copy()
	nt.Resources = nt.Resources.Copy()
	nt.Meta = helper.CopyMapStringString(nt.Meta)
	nt.DispatchPayload = nt.DispatchPayload.Copy()

	if t.Artifacts != nil {
		artifacts := make([]*TaskArtifact, 0, len(t.Artifacts))
		for _, a := range nt.Artifacts {
			artifacts = append(artifacts, a.Copy())
		}
		nt.Artifacts = artifacts
	}

	if i, err := copystructure.Copy(nt.Config); err != nil {
		panic(err.Error())
	} else {
		nt.Config = i.(map[string]interface{})
	}

	if t.Templates != nil {
		templates := make([]*Template, len(t.Templates))
		for i, tmpl := range nt.Templates {
			templates[i] = tmpl.Copy()
		}
		nt.Templates = templates
	}

	return nt
}

// Canonicalize canonicalizes fields in the task.
func (t *Task) Canonicalize(job *Job, tg *TaskGroup) {
	// Ensure that an empty and nil map are treated the same to avoid scheduling
	// problems since we use reflect DeepEquals.
	if len(t.Meta) == 0 {
		t.Meta = nil
	}
	if len(t.Config) == 0 {
		t.Config = nil
	}
	if len(t.Env) == 0 {
		t.Env = nil
	}

	for _, service := range t.Services {
		service.Canonicalize(job.Name, tg.Name, t.Name)
	}

	// If Resources are nil initialize them to defaults, otherwise canonicalize
	if t.Resources == nil {
		t.Resources = DefaultResources()
	} else {
		t.Resources.Canonicalize()
	}

	// Set the default timeout if it is not specified.
	if t.KillTimeout == 0 {
		t.KillTimeout = DefaultKillTimeout
	}

	if t.Vault != nil {
		t.Vault.Canonicalize()
	}

	for _, template := range t.Templates {
		template.Canonicalize()
	}
}

func (t *Task) GoString() string {
	return fmt.Sprintf("*%#v", *t)
}

// Validate is used to sanity check a task
func (t *Task) Validate(ephemeralDisk *EphemeralDisk) error {
	var mErr multierror.Error
	if t.Name == "" {
		mErr.Errors = append(mErr.Errors, errors.New("Missing task name"))
	}
	if strings.ContainsAny(t.Name, `/\`) {
		// We enforce this so that when creating the directory on disk it will
		// not have any slashes.
		mErr.Errors = append(mErr.Errors, errors.New("Task name cannot include slashes"))
	}
	if t.Driver == "" {
		mErr.Errors = append(mErr.Errors, errors.New("Missing task driver"))
	}
	if t.KillTimeout < 0 {
		mErr.Errors = append(mErr.Errors, errors.New("KillTimeout must be a positive value"))
	}
	if t.ShutdownDelay < 0 {
		mErr.Errors = append(mErr.Errors, errors.New("ShutdownDelay must be a positive value"))
	}

	// Validate the resources.
	if t.Resources == nil {
		mErr.Errors = append(mErr.Errors, errors.New("Missing task resources"))
	} else {
		if err := t.Resources.MeetsMinResources(); err != nil {
			mErr.Errors = append(mErr.Errors, err)
		}

		// Ensure the task isn't asking for disk resources
		if t.Resources.DiskMB > 0 {
			mErr.Errors = append(mErr.Errors, errors.New("Task can't ask for disk resources, they have to be specified at the task group level."))
		}
	}

	// Validate the log config
	if t.LogConfig == nil {
		mErr.Errors = append(mErr.Errors, errors.New("Missing Log Config"))
	} else if err := t.LogConfig.Validate(); err != nil {
		mErr.Errors = append(mErr.Errors, err)
	}

	for idx, constr := range t.Constraints {
		if err := constr.Validate(); err != nil {
			outer := fmt.Errorf("Constraint %d validation failed: %s", idx+1, err)
			mErr.Errors = append(mErr.Errors, outer)
		}

		switch constr.Operand {
		case ConstraintDistinctHosts, ConstraintDistinctProperty:
			outer := fmt.Errorf("Constraint %d has disallowed Operand at task level: %s", idx+1, constr.Operand)
			mErr.Errors = append(mErr.Errors, outer)
		}
	}

	// Validate Services
	if err := validateServices(t); err != nil {
		mErr.Errors = append(mErr.Errors, err)
	}

	if t.LogConfig != nil && ephemeralDisk != nil {
		logUsage := (t.LogConfig.MaxFiles * t.LogConfig.MaxFileSizeMB)
		if ephemeralDisk.SizeMB <= logUsage {
			mErr.Errors = append(mErr.Errors,
				fmt.Errorf("log storage (%d MB) must be less than requested disk capacity (%d MB)",
					logUsage, ephemeralDisk.SizeMB))
		}
	}

	for idx, artifact := range t.Artifacts {
		if err := artifact.Validate(); err != nil {
			outer := fmt.Errorf("Artifact %d validation failed: %v", idx+1, err)
			mErr.Errors = append(mErr.Errors, outer)
		}
	}

	if t.Vault != nil {
		if err := t.Vault.Validate(); err != nil {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Vault validation failed: %v", err))
		}
	}

	destinations := make(map[string]int, len(t.Templates))
	for idx, tmpl := range t.Templates {
		if err := tmpl.Validate(); err != nil {
			outer := fmt.Errorf("Template %d validation failed: %s", idx+1, err)
			mErr.Errors = append(mErr.Errors, outer)
		}

		if other, ok := destinations[tmpl.DestPath]; ok {
			outer := fmt.Errorf("Template %d has same destination as %d", idx+1, other)
			mErr.Errors = append(mErr.Errors, outer)
		} else {
			destinations[tmpl.DestPath] = idx + 1
		}
	}

	// Validate the dispatch payload block if there
	if t.DispatchPayload != nil {
		if err := t.DispatchPayload.Validate(); err != nil {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Dispatch Payload validation failed: %v", err))
		}
	}

	return mErr.ErrorOrNil()
}

// validateServices takes a task and validates the services within it are valid
// and reference ports that exist.
func validateServices(t *Task) error {
	var mErr multierror.Error

	// Ensure that services don't ask for nonexistent ports and their names are
	// unique.
	servicePorts := make(map[string]map[string]struct{})
	addServicePort := func(label, service string) {
		if _, ok := servicePorts[label]; !ok {
			servicePorts[label] = map[string]struct{}{}
		}
		servicePorts[label][service] = struct{}{}
	}
	knownServices := make(map[string]struct{})
	for i, service := range t.Services {
		if err := service.Validate(); err != nil {
			outer := fmt.Errorf("service[%d] %+q validation failed: %s", i, service.Name, err)
			mErr.Errors = append(mErr.Errors, outer)
		}

		// Ensure that services with the same name are not being registered for
		// the same port
		if _, ok := knownServices[service.Name+service.PortLabel]; ok {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("service %q is duplicate", service.Name))
		}
		knownServices[service.Name+service.PortLabel] = struct{}{}

		if service.PortLabel != "" {
			if service.AddressMode == "driver" {
				// Numeric port labels are valid for address_mode=driver
				_, err := strconv.Atoi(service.PortLabel)
				if err != nil {
					// Not a numeric port label, add it to list to check
					addServicePort(service.PortLabel, service.Name)
				}
			} else {
				addServicePort(service.PortLabel, service.Name)
			}
		}

		// Ensure that check names are unique and have valid ports
		knownChecks := make(map[string]struct{})
		for _, check := range service.Checks {
			if _, ok := knownChecks[check.Name]; ok {
				mErr.Errors = append(mErr.Errors, fmt.Errorf("check %q is duplicate", check.Name))
			}
			knownChecks[check.Name] = struct{}{}

			if !check.RequiresPort() {
				// No need to continue validating check if it doesn't need a port
				continue
			}

			effectivePort := check.PortLabel
			if effectivePort == "" {
				// Inherits from service
				effectivePort = service.PortLabel
			}

			if effectivePort == "" {
				mErr.Errors = append(mErr.Errors, fmt.Errorf("check %q is missing a port", check.Name))
				continue
			}

			isNumeric := false
			portNumber, err := strconv.Atoi(effectivePort)
			if err == nil {
				isNumeric = true
			}

			// Numeric ports are fine for address_mode = "driver"
			if check.AddressMode == "driver" && isNumeric {
				if portNumber <= 0 {
					mErr.Errors = append(mErr.Errors, fmt.Errorf("check %q has invalid numeric port %d", check.Name, portNumber))
				}
				continue
			}

			if isNumeric {
				mErr.Errors = append(mErr.Errors, fmt.Errorf(`check %q cannot use a numeric port %d without setting address_mode="driver"`, check.Name, portNumber))
				continue
			}

			// PortLabel must exist, report errors by its parent service
			addServicePort(effectivePort, service.Name)
		}
	}

	// Get the set of port labels.
	portLabels := make(map[string]struct{})
	if t.Resources != nil {
		for _, network := range t.Resources.Networks {
			ports := network.PortLabels()
			for portLabel := range ports {
				portLabels[portLabel] = struct{}{}
			}
		}
	}

	// Iterate over a sorted list of keys to make error listings stable
	keys := make([]string, 0, len(servicePorts))
	for p := range servicePorts {
		keys = append(keys, p)
	}
	sort.Strings(keys)

	// Ensure all ports referenced in services exist.
	for _, servicePort := range keys {
		services := servicePorts[servicePort]
		_, ok := portLabels[servicePort]
		if !ok {
			names := make([]string, 0, len(services))
			for name := range services {
				names = append(names, name)
			}

			// Keep order deterministic
			sort.Strings(names)
			joined := strings.Join(names, ", ")
			err := fmt.Errorf("port label %q referenced by services %v does not exist", servicePort, joined)
			mErr.Errors = append(mErr.Errors, err)
		}
	}

	// Ensure address mode is valid
	return mErr.ErrorOrNil()
}

const (
	// TemplateChangeModeNoop marks that no action should be taken if the
	// template is re-rendered
	TemplateChangeModeNoop = "noop"

	// TemplateChangeModeSignal marks that the task should be signaled if the
	// template is re-rendered
	TemplateChangeModeSignal = "signal"

	// TemplateChangeModeRestart marks that the task should be restarted if the
	// template is re-rendered
	TemplateChangeModeRestart = "restart"
)

var (
	// TemplateChangeModeInvalidError is the error for when an invalid change
	// mode is given
	TemplateChangeModeInvalidError = errors.New("Invalid change mode. Must be one of the following: noop, signal, restart")
)

// Template represents a template configuration to be rendered for a given task
type Template struct {
	// SourcePath is the path to the template to be rendered
	SourcePath string

	// DestPath is the path to where the template should be rendered
	DestPath string

	// EmbeddedTmpl store the raw template. This is useful for smaller templates
	// where they are embedded in the job file rather than sent as an artifact
	EmbeddedTmpl string

	// ChangeMode indicates what should be done if the template is re-rendered
	ChangeMode string

	// ChangeSignal is the signal that should be sent if the change mode
	// requires it.
	ChangeSignal string

	// Splay is used to avoid coordinated restarts of processes by applying a
	// random wait between 0 and the given splay value before signalling the
	// application of a change
	Splay time.Duration

	// Perms is the permission the file should be written out with.
	Perms string

	// LeftDelim and RightDelim are optional configurations to control what
	// delimiter is utilized when parsing the template.
	LeftDelim  string
	RightDelim string

	// Envvars enables exposing the template as environment variables
	// instead of as a file. The template must be of the form:
	//
	//	VAR_NAME_1={{ key service/my-key }}
	//	VAR_NAME_2=raw string and {{ env "attr.kernel.name" }}
	//
	// Lines will be split on the initial "=" with the first part being the
	// key name and the second part the value.
	// Empty lines and lines starting with # will be ignored, but to avoid
	// escaping issues #s within lines will not be treated as comments.
	Envvars bool

	// VaultGrace is the grace duration between lease renewal and reacquiring a
	// secret. If the lease of a secret is less than the grace, a new secret is
	// acquired.
	VaultGrace time.Duration
}

// DefaultTemplate returns a default template.
func DefaultTemplate() *Template {
	return &Template{
		ChangeMode: TemplateChangeModeRestart,
		Splay:      5 * time.Second,
		Perms:      "0644",
	}
}

func (t *Template) Copy() *Template {
	if t == nil {
		return nil
	}
	copy := new(Template)
	*copy = *t
	return copy
}

func (t *Template) Canonicalize() {
	if t.ChangeSignal != "" {
		t.ChangeSignal = strings.ToUpper(t.ChangeSignal)
	}
}

func (t *Template) Validate() error {
	var mErr multierror.Error

	// Verify we have something to render
	if t.SourcePath == "" && t.EmbeddedTmpl == "" {
		multierror.Append(&mErr, fmt.Errorf("Must specify a source path or have an embedded template"))
	}

	// Verify we can render somewhere
	if t.DestPath == "" {
		multierror.Append(&mErr, fmt.Errorf("Must specify a destination for the template"))
	}

	// Verify the destination doesn't escape
	escaped, err := PathEscapesAllocDir("task", t.DestPath)
	if err != nil {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("invalid destination path: %v", err))
	} else if escaped {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("destination escapes allocation directory"))
	}

	// Verify a proper change mode
	switch t.ChangeMode {
	case TemplateChangeModeNoop, TemplateChangeModeRestart:
	case TemplateChangeModeSignal:
		if t.ChangeSignal == "" {
			multierror.Append(&mErr, fmt.Errorf("Must specify signal value when change mode is signal"))
		}
		if t.Envvars {
			multierror.Append(&mErr, fmt.Errorf("cannot use signals with env var templates"))
		}
	default:
		multierror.Append(&mErr, TemplateChangeModeInvalidError)
	}

	// Verify the splay is positive
	if t.Splay < 0 {
		multierror.Append(&mErr, fmt.Errorf("Must specify positive splay value"))
	}

	// Verify the permissions
	if t.Perms != "" {
		if _, err := strconv.ParseUint(t.Perms, 8, 12); err != nil {
			multierror.Append(&mErr, fmt.Errorf("Failed to parse %q as octal: %v", t.Perms, err))
		}
	}

	if t.VaultGrace.Nanoseconds() < 0 {
		multierror.Append(&mErr, fmt.Errorf("Vault grace must be greater than zero: %v < 0", t.VaultGrace))
	}

	return mErr.ErrorOrNil()
}

// Set of possible states for a task.
const (
	TaskStatePending = "pending" // The task is waiting to be run.
	TaskStateRunning = "running" // The task is currently running.
	TaskStateDead    = "dead"    // Terminal state of task.
)

// TaskState tracks the current state of a task and events that caused state
// transitions.
type TaskState struct {
	// The current state of the task.
	State string

	// Failed marks a task as having failed
	Failed bool

	// Restarts is the number of times the task has restarted
	Restarts uint64

	// LastRestart is the time the task last restarted. It is updated each time the
	// task restarts
	LastRestart time.Time

	// StartedAt is the time the task is started. It is updated each time the
	// task starts
	StartedAt time.Time

	// FinishedAt is the time at which the task transitioned to dead and will
	// not be started again.
	FinishedAt time.Time

	// Series of task events that transition the state of the task.
	Events []*TaskEvent
}

func (ts *TaskState) Copy() *TaskState {
	if ts == nil {
		return nil
	}
	copy := new(TaskState)
	*copy = *ts

	if ts.Events != nil {
		copy.Events = make([]*TaskEvent, len(ts.Events))
		for i, e := range ts.Events {
			copy.Events[i] = e.Copy()
		}
	}
	return copy
}

// Successful returns whether a task finished successfully. This doesn't really
// have meaning on a non-batch allocation because a service and system
// allocation should not finish.
func (ts *TaskState) Successful() bool {
	l := len(ts.Events)
	if ts.State != TaskStateDead || l == 0 {
		return false
	}

	e := ts.Events[l-1]
	if e.Type != TaskTerminated {
		return false
	}

	return e.ExitCode == 0
}

const (
	// TaskSetupFailure indicates that the task could not be started due to a
	// a setup failure.
	TaskSetupFailure = "Setup Failure"

	// TaskDriveFailure indicates that the task could not be started due to a
	// failure in the driver.
	TaskDriverFailure = "Driver Failure"

	// TaskReceived signals that the task has been pulled by the client at the
	// given timestamp.
	TaskReceived = "Received"

	// TaskFailedValidation indicates the task was invalid and as such was not
	// run.
	TaskFailedValidation = "Failed Validation"

	// TaskStarted signals that the task was started and its timestamp can be
	// used to determine the running length of the task.
	TaskStarted = "Started"

	// TaskTerminated indicates that the task was started and exited.
	TaskTerminated = "Terminated"

	// TaskKilling indicates a kill signal has been sent to the task.
	TaskKilling = "Killing"

	// TaskKilled indicates a user has killed the task.
	TaskKilled = "Killed"

	// TaskRestarting indicates that task terminated and is being restarted.
	TaskRestarting = "Restarting"

	// TaskNotRestarting indicates that the task has failed and is not being
	// restarted because it has exceeded its restart policy.
	TaskNotRestarting = "Not Restarting"

	// TaskRestartSignal indicates that the task has been signalled to be
	// restarted
	TaskRestartSignal = "Restart Signaled"

	// TaskSignaling indicates that the task is being signalled.
	TaskSignaling = "Signaling"

	// TaskDownloadingArtifacts means the task is downloading the artifacts
	// specified in the task.
	TaskDownloadingArtifacts = "Downloading Artifacts"

	// TaskArtifactDownloadFailed indicates that downloading the artifacts
	// failed.
	TaskArtifactDownloadFailed = "Failed Artifact Download"

	// TaskBuildingTaskDir indicates that the task directory/chroot is being
	// built.
	TaskBuildingTaskDir = "Building Task Directory"

	// TaskSetup indicates the task runner is setting up the task environment
	TaskSetup = "Task Setup"

	// TaskDiskExceeded indicates that one of the tasks in a taskgroup has
	// exceeded the requested disk resources.
	TaskDiskExceeded = "Disk Resources Exceeded"

	// TaskSiblingFailed indicates that a sibling task in the task group has
	// failed.
	TaskSiblingFailed = "Sibling Task Failed"

	// TaskDriverMessage is an informational event message emitted by
	// drivers such as when they're performing a long running action like
	// downloading an image.
	TaskDriverMessage = "Driver"

	// TaskLeaderDead indicates that the leader task within the has finished.
	TaskLeaderDead = "Leader Task Dead"
)

// TaskEvent is an event that effects the state of a task and contains meta-data
// appropriate to the events type.
type TaskEvent struct {
	Type string
	Time int64 // Unix Nanosecond timestamp

	Message string // A possible message explaining the termination of the task.

	// DisplayMessage is a human friendly message about the event
	DisplayMessage string

	// Details is a map with annotated info about the event
	Details map[string]string

	// DEPRECATION NOTICE: The following fields are deprecated and will be removed
	// in a future release. Field values are available in the Details map.

	// FailsTask marks whether this event fails the task.
	// Deprecated, use Details["fails_task"] to access this.
	FailsTask bool

	// Restart fields.
	// Deprecated, use Details["restart_reason"] to access this.
	RestartReason string

	// Setup Failure fields.
	// Deprecated, use Details["setup_error"] to access this.
	SetupError string

	// Driver Failure fields.
	// Deprecated, use Details["driver_error"] to access this.
	DriverError string // A driver error occurred while starting the task.

	// Task Terminated Fields.

	// Deprecated, use Details["exit_code"] to access this.
	ExitCode int // The exit code of the task.

	// Deprecated, use Details["signal"] to access this.
	Signal int // The signal that terminated the task.

	// Killing fields
	// Deprecated, use Details["kill_timeout"] to access this.
	KillTimeout time.Duration

	// Task Killed Fields.
	// Deprecated, use Details["kill_error"] to access this.
	KillError string // Error killing the task.

	// KillReason is the reason the task was killed
	// Deprecated, use Details["kill_reason"] to access this.
	KillReason string

	// TaskRestarting fields.
	// Deprecated, use Details["start_delay"] to access this.
	StartDelay int64 // The sleep period before restarting the task in unix nanoseconds.

	// Artifact Download fields
	// Deprecated, use Details["download_error"] to access this.
	DownloadError string // Error downloading artifacts

	// Validation fields
	// Deprecated, use Details["validation_error"] to access this.
	ValidationError string // Validation error

	// The maximum allowed task disk size.
	// Deprecated, use Details["disk_limit"] to access this.
	DiskLimit int64

	// Name of the sibling task that caused termination of the task that
	// the TaskEvent refers to.
	// Deprecated, use Details["failed_sibling"] to access this.
	FailedSibling string

	// VaultError is the error from token renewal
	// Deprecated, use Details["vault_renewal_error"] to access this.
	VaultError string

	// TaskSignalReason indicates the reason the task is being signalled.
	// Deprecated, use Details["task_signal_reason"] to access this.
	TaskSignalReason string

	// TaskSignal is the signal that was sent to the task
	// Deprecated, use Details["task_signal"] to access this.
	TaskSignal string

	// DriverMessage indicates a driver action being taken.
	// Deprecated, use Details["driver_message"] to access this.
	DriverMessage string

	// GenericSource is the source of a message.
	// Deprecated, is redundant with event type.
	GenericSource string
}

func (event *TaskEvent) PopulateEventDisplayMessage() {
	// Build up the description based on the event type.
	if event == nil { //TODO(preetha) needs investigation alloc_runner's Run method sends a nil event when sigterming nomad. Why?
		return
	}

	if event.DisplayMessage != "" {
		return
	}

	var desc string
	switch event.Type {
	case TaskSetup:
		desc = event.Message
	case TaskStarted:
		desc = "Task started by client"
	case TaskReceived:
		desc = "Task received by client"
	case TaskFailedValidation:
		if event.ValidationError != "" {
			desc = event.ValidationError
		} else {
			desc = "Validation of task failed"
		}
	case TaskSetupFailure:
		if event.SetupError != "" {
			desc = event.SetupError
		} else {
			desc = "Task setup failed"
		}
	case TaskDriverFailure:
		if event.DriverError != "" {
			desc = event.DriverError
		} else {
			desc = "Failed to start task"
		}
	case TaskDownloadingArtifacts:
		desc = "Client is downloading artifacts"
	case TaskArtifactDownloadFailed:
		if event.DownloadError != "" {
			desc = event.DownloadError
		} else {
			desc = "Failed to download artifacts"
		}
	case TaskKilling:
		if event.KillReason != "" {
			desc = event.KillReason
		} else if event.KillTimeout != 0 {
			desc = fmt.Sprintf("Sent interrupt. Waiting %v before force killing", event.KillTimeout)
		} else {
			desc = "Sent interrupt"
		}
	case TaskKilled:
		if event.KillError != "" {
			desc = event.KillError
		} else {
			desc = "Task successfully killed"
		}
	case TaskTerminated:
		var parts []string
		parts = append(parts, fmt.Sprintf("Exit Code: %d", event.ExitCode))

		if event.Signal != 0 {
			parts = append(parts, fmt.Sprintf("Signal: %d", event.Signal))
		}

		if event.Message != "" {
			parts = append(parts, fmt.Sprintf("Exit Message: %q", event.Message))
		}
		desc = strings.Join(parts, ", ")
	case TaskRestarting:
		in := fmt.Sprintf("Task restarting in %v", time.Duration(event.StartDelay))
		if event.RestartReason != "" && event.RestartReason != ReasonWithinPolicy {
			desc = fmt.Sprintf("%s - %s", event.RestartReason, in)
		} else {
			desc = in
		}
	case TaskNotRestarting:
		if event.RestartReason != "" {
			desc = event.RestartReason
		} else {
			desc = "Task exceeded restart policy"
		}
	case TaskSiblingFailed:
		if event.FailedSibling != "" {
			desc = fmt.Sprintf("Task's sibling %q failed", event.FailedSibling)
		} else {
			desc = "Task's sibling failed"
		}
	case TaskSignaling:
		sig := event.TaskSignal
		reason := event.TaskSignalReason

		if sig == "" && reason == "" {
			desc = "Task being sent a signal"
		} else if sig == "" {
			desc = reason
		} else if reason == "" {
			desc = fmt.Sprintf("Task being sent signal %v", sig)
		} else {
			desc = fmt.Sprintf("Task being sent signal %v: %v", sig, reason)
		}
	case TaskRestartSignal:
		if event.RestartReason != "" {
			desc = event.RestartReason
		} else {
			desc = "Task signaled to restart"
		}
	case TaskDriverMessage:
		desc = event.DriverMessage
	case TaskLeaderDead:
		desc = "Leader Task in Group dead"
	default:
		desc = event.Message
	}

	event.DisplayMessage = desc
}

func (te *TaskEvent) GoString() string {
	return fmt.Sprintf("%v - %v", te.Time, te.Type)
}

// SetMessage sets the message of TaskEvent
func (te *TaskEvent) SetMessage(msg string) *TaskEvent {
	te.Message = msg
	te.Details["message"] = msg
	return te
}

func (te *TaskEvent) Copy() *TaskEvent {
	if te == nil {
		return nil
	}
	copy := new(TaskEvent)
	*copy = *te
	return copy
}

func NewTaskEvent(event string) *TaskEvent {
	return &TaskEvent{
		Type:    event,
		Time:    time.Now().UnixNano(),
		Details: make(map[string]string),
	}
}

// SetSetupError is used to store an error that occurred while setting up the
// task
func (e *TaskEvent) SetSetupError(err error) *TaskEvent {
	if err != nil {
		e.SetupError = err.Error()
		e.Details["setup_error"] = err.Error()
	}
	return e
}

func (e *TaskEvent) SetFailsTask() *TaskEvent {
	e.FailsTask = true
	e.Details["fails_task"] = "true"
	return e
}

func (e *TaskEvent) SetDriverError(err error) *TaskEvent {
	if err != nil {
		e.DriverError = err.Error()
		e.Details["driver_error"] = err.Error()
	}
	return e
}

func (e *TaskEvent) SetExitCode(c int) *TaskEvent {
	e.ExitCode = c
	e.Details["exit_code"] = fmt.Sprintf("%d", c)
	return e
}

func (e *TaskEvent) SetSignal(s int) *TaskEvent {
	e.Signal = s
	e.Details["signal"] = fmt.Sprintf("%d", s)
	return e
}

func (e *TaskEvent) SetExitMessage(err error) *TaskEvent {
	if err != nil {
		e.Message = err.Error()
		e.Details["exit_message"] = err.Error()
	}
	return e
}

func (e *TaskEvent) SetKillError(err error) *TaskEvent {
	if err != nil {
		e.KillError = err.Error()
		e.Details["kill_error"] = err.Error()
	}
	return e
}

func (e *TaskEvent) SetKillReason(r string) *TaskEvent {
	e.KillReason = r
	e.Details["kill_reason"] = r
	return e
}

func (e *TaskEvent) SetRestartDelay(delay time.Duration) *TaskEvent {
	e.StartDelay = int64(delay)
	e.Details["start_delay"] = fmt.Sprintf("%d", delay)
	return e
}

func (e *TaskEvent) SetRestartReason(reason string) *TaskEvent {
	e.RestartReason = reason
	e.Details["restart_reason"] = reason
	return e
}

func (e *TaskEvent) SetTaskSignalReason(r string) *TaskEvent {
	e.TaskSignalReason = r
	e.Details["task_signal_reason"] = r
	return e
}

func (e *TaskEvent) SetTaskSignal(s os.Signal) *TaskEvent {
	e.TaskSignal = s.String()
	e.Details["task_signal"] = s.String()
	return e
}

func (e *TaskEvent) SetDownloadError(err error) *TaskEvent {
	if err != nil {
		e.DownloadError = err.Error()
		e.Details["download_error"] = err.Error()
	}
	return e
}

func (e *TaskEvent) SetValidationError(err error) *TaskEvent {
	if err != nil {
		e.ValidationError = err.Error()
		e.Details["validation_error"] = err.Error()
	}
	return e
}

func (e *TaskEvent) SetKillTimeout(timeout time.Duration) *TaskEvent {
	e.KillTimeout = timeout
	e.Details["kill_timeout"] = timeout.String()
	return e
}

func (e *TaskEvent) SetDiskLimit(limit int64) *TaskEvent {
	e.DiskLimit = limit
	e.Details["disk_limit"] = fmt.Sprintf("%d", limit)
	return e
}

func (e *TaskEvent) SetFailedSibling(sibling string) *TaskEvent {
	e.FailedSibling = sibling
	e.Details["failed_sibling"] = sibling
	return e
}

func (e *TaskEvent) SetVaultRenewalError(err error) *TaskEvent {
	if err != nil {
		e.VaultError = err.Error()
		e.Details["vault_renewal_error"] = err.Error()
	}
	return e
}

func (e *TaskEvent) SetDriverMessage(m string) *TaskEvent {
	e.DriverMessage = m
	e.Details["driver_message"] = m
	return e
}

// TaskArtifact is an artifact to download before running the task.
type TaskArtifact struct {
	// GetterSource is the source to download an artifact using go-getter
	GetterSource string

	// GetterOptions are options to use when downloading the artifact using
	// go-getter.
	GetterOptions map[string]string

	// GetterMode is the go-getter.ClientMode for fetching resources.
	// Defaults to "any" but can be set to "file" or "dir".
	GetterMode string

	// RelativeDest is the download destination given relative to the task's
	// directory.
	RelativeDest string
}

func (ta *TaskArtifact) Copy() *TaskArtifact {
	if ta == nil {
		return nil
	}
	nta := new(TaskArtifact)
	*nta = *ta
	nta.GetterOptions = helper.CopyMapStringString(ta.GetterOptions)
	return nta
}

func (ta *TaskArtifact) GoString() string {
	return fmt.Sprintf("%+v", ta)
}

// PathEscapesAllocDir returns if the given path escapes the allocation
// directory. The prefix allows adding a prefix if the path will be joined, for
// example a "task/local" prefix may be provided if the path will be joined
// against that prefix.
func PathEscapesAllocDir(prefix, path string) (bool, error) {
	// Verify the destination doesn't escape the tasks directory
	alloc, err := filepath.Abs(filepath.Join("/", "alloc-dir/", "alloc-id/"))
	if err != nil {
		return false, err
	}
	abs, err := filepath.Abs(filepath.Join(alloc, prefix, path))
	if err != nil {
		return false, err
	}
	rel, err := filepath.Rel(alloc, abs)
	if err != nil {
		return false, err
	}

	return strings.HasPrefix(rel, ".."), nil
}

func (ta *TaskArtifact) Validate() error {
	// Verify the source
	var mErr multierror.Error
	if ta.GetterSource == "" {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("source must be specified"))
	}

	switch ta.GetterMode {
	case "":
		// Default to any
		ta.GetterMode = GetterModeAny
	case GetterModeAny, GetterModeFile, GetterModeDir:
		// Ok
	default:
		mErr.Errors = append(mErr.Errors, fmt.Errorf("invalid artifact mode %q; must be one of: %s, %s, %s",
			ta.GetterMode, GetterModeAny, GetterModeFile, GetterModeDir))
	}

	escaped, err := PathEscapesAllocDir("task", ta.RelativeDest)
	if err != nil {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("invalid destination path: %v", err))
	} else if escaped {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("destination escapes allocation directory"))
	}

	// Verify the checksum
	if check, ok := ta.GetterOptions["checksum"]; ok {
		check = strings.TrimSpace(check)
		if check == "" {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("checksum value cannot be empty"))
			return mErr.ErrorOrNil()
		}

		parts := strings.Split(check, ":")
		if l := len(parts); l != 2 {
			mErr.Errors = append(mErr.Errors, fmt.Errorf(`checksum must be given as "type:value"; got %q`, check))
			return mErr.ErrorOrNil()
		}

		checksumVal := parts[1]
		checksumBytes, err := hex.DecodeString(checksumVal)
		if err != nil {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("invalid checksum: %v", err))
			return mErr.ErrorOrNil()
		}

		checksumType := parts[0]
		expectedLength := 0
		switch checksumType {
		case "md5":
			expectedLength = md5.Size
		case "sha1":
			expectedLength = sha1.Size
		case "sha256":
			expectedLength = sha256.Size
		case "sha512":
			expectedLength = sha512.Size
		default:
			mErr.Errors = append(mErr.Errors, fmt.Errorf("unsupported checksum type: %s", checksumType))
			return mErr.ErrorOrNil()
		}

		if len(checksumBytes) != expectedLength {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("invalid %s checksum: %v", checksumType, checksumVal))
			return mErr.ErrorOrNil()
		}
	}

	return mErr.ErrorOrNil()
}

const (
	ConstraintDistinctProperty = "distinct_property"
	ConstraintDistinctHosts    = "distinct_hosts"
	ConstraintRegex            = "regexp"
	ConstraintVersion          = "version"
	ConstraintSetContains      = "set_contains"
)

// Constraints are used to restrict placement options.
type Constraint struct {
	LTarget string // Left-hand target
	RTarget string // Right-hand target
	Operand string // Constraint operand (<=, <, =, !=, >, >=), contains, near
	str     string // Memoized string
}

// Equal checks if two constraints are equal
func (c *Constraint) Equal(o *Constraint) bool {
	return c.LTarget == o.LTarget &&
		c.RTarget == o.RTarget &&
		c.Operand == o.Operand
}

func (c *Constraint) Copy() *Constraint {
	if c == nil {
		return nil
	}
	nc := new(Constraint)
	*nc = *c
	return nc
}

func (c *Constraint) String() string {
	if c.str != "" {
		return c.str
	}
	c.str = fmt.Sprintf("%s %s %s", c.LTarget, c.Operand, c.RTarget)
	return c.str
}

func (c *Constraint) Validate() error {
	var mErr multierror.Error
	if c.Operand == "" {
		mErr.Errors = append(mErr.Errors, errors.New("Missing constraint operand"))
	}

	// requireLtarget specifies whether the constraint requires an LTarget to be
	// provided.
	requireLtarget := true

	// Perform additional validation based on operand
	switch c.Operand {
	case ConstraintDistinctHosts:
		requireLtarget = false
	case ConstraintSetContains:
		if c.RTarget == "" {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Set contains constraint requires an RTarget"))
		}
	case ConstraintRegex:
		if _, err := regexp.Compile(c.RTarget); err != nil {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Regular expression failed to compile: %v", err))
		}
	case ConstraintVersion:
		if _, err := version.NewConstraint(c.RTarget); err != nil {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Version constraint is invalid: %v", err))
		}
	case ConstraintDistinctProperty:
		// If a count is set, make sure it is convertible to a uint64
		if c.RTarget != "" {
			count, err := strconv.ParseUint(c.RTarget, 10, 64)
			if err != nil {
				mErr.Errors = append(mErr.Errors, fmt.Errorf("Failed to convert RTarget %q to uint64: %v", c.RTarget, err))
			} else if count < 1 {
				mErr.Errors = append(mErr.Errors, fmt.Errorf("Distinct Property must have an allowed count of 1 or greater: %d < 1", count))
			}
		}
	case "=", "==", "is", "!=", "not", "<", "<=", ">", ">=":
		if c.RTarget == "" {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Operator %q requires an RTarget", c.Operand))
		}
	default:
		mErr.Errors = append(mErr.Errors, fmt.Errorf("Unknown constraint type %q", c.Operand))
	}

	// Ensure we have an LTarget for the constraints that need one
	if requireLtarget && c.LTarget == "" {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("No LTarget provided but is required by constraint"))
	}

	return mErr.ErrorOrNil()
}

// EphemeralDisk is an ephemeral disk object
type EphemeralDisk struct {
	// Sticky indicates whether the allocation is sticky to a node
	Sticky bool

	// SizeMB is the size of the local disk
	SizeMB int

	// Migrate determines if Nomad client should migrate the allocation dir for
	// sticky allocations
	Migrate bool
}

// DefaultEphemeralDisk returns a EphemeralDisk with default configurations
func DefaultEphemeralDisk() *EphemeralDisk {
	return &EphemeralDisk{
		SizeMB: 300,
	}
}

// Validate validates EphemeralDisk
func (d *EphemeralDisk) Validate() error {
	if d.SizeMB < 10 {
		return fmt.Errorf("minimum DiskMB value is 10; got %d", d.SizeMB)
	}
	return nil
}

// Copy copies the EphemeralDisk struct and returns a new one
func (d *EphemeralDisk) Copy() *EphemeralDisk {
	ld := new(EphemeralDisk)
	*ld = *d
	return ld
}

var (
	// VaultUnrecoverableError matches unrecoverable errors returned by a Vault
	// server
	VaultUnrecoverableError = regexp.MustCompile(`Code:\s+40(0|3|4)`)
)

const (
	// VaultChangeModeNoop takes no action when a new token is retrieved.
	VaultChangeModeNoop = "noop"

	// VaultChangeModeSignal signals the task when a new token is retrieved.
	VaultChangeModeSignal = "signal"

	// VaultChangeModeRestart restarts the task when a new token is retrieved.
	VaultChangeModeRestart = "restart"
)

// Vault stores the set of permissions a task needs access to from Vault.
type Vault struct {
	// Policies is the set of policies that the task needs access to
	Policies []string

	// Env marks whether the Vault Token should be exposed as an environment
	// variable
	Env bool

	// ChangeMode is used to configure the task's behavior when the Vault
	// token changes because the original token could not be renewed in time.
	ChangeMode string

	// ChangeSignal is the signal sent to the task when a new token is
	// retrieved. This is only valid when using the signal change mode.
	ChangeSignal string
}

func DefaultVaultBlock() *Vault {
	return &Vault{
		Env:        true,
		ChangeMode: VaultChangeModeRestart,
	}
}

// Copy returns a copy of this Vault block.
func (v *Vault) Copy() *Vault {
	if v == nil {
		return nil
	}

	nv := new(Vault)
	*nv = *v
	return nv
}

func (v *Vault) Canonicalize() {
	if v.ChangeSignal != "" {
		v.ChangeSignal = strings.ToUpper(v.ChangeSignal)
	}
}

// Validate returns if the Vault block is valid.
func (v *Vault) Validate() error {
	if v == nil {
		return nil
	}

	var mErr multierror.Error
	if len(v.Policies) == 0 {
		multierror.Append(&mErr, fmt.Errorf("Policy list cannot be empty"))
	}

	for _, p := range v.Policies {
		if p == "root" {
			multierror.Append(&mErr, fmt.Errorf("Can not specify \"root\" policy"))
		}
	}

	switch v.ChangeMode {
	case VaultChangeModeSignal:
		if v.ChangeSignal == "" {
			multierror.Append(&mErr, fmt.Errorf("Signal must be specified when using change mode %q", VaultChangeModeSignal))
		}
	case VaultChangeModeNoop, VaultChangeModeRestart:
	default:
		multierror.Append(&mErr, fmt.Errorf("Unknown change mode %q", v.ChangeMode))
	}

	return mErr.ErrorOrNil()
}

const (
	// DeploymentStatuses are the various states a deployment can be be in
	DeploymentStatusRunning    = "running"
	DeploymentStatusPaused     = "paused"
	DeploymentStatusFailed     = "failed"
	DeploymentStatusSuccessful = "successful"
	DeploymentStatusCancelled  = "cancelled"

	// DeploymentStatusDescriptions are the various descriptions of the states a
	// deployment can be in.
	DeploymentStatusDescriptionRunning               = "Deployment is running"
	DeploymentStatusDescriptionRunningNeedsPromotion = "Deployment is running but requires promotion"
	DeploymentStatusDescriptionPaused                = "Deployment is paused"
	DeploymentStatusDescriptionSuccessful            = "Deployment completed successfully"
	DeploymentStatusDescriptionStoppedJob            = "Cancelled because job is stopped"
	DeploymentStatusDescriptionNewerJob              = "Cancelled due to newer version of job"
	DeploymentStatusDescriptionFailedAllocations     = "Failed due to unhealthy allocations"
	DeploymentStatusDescriptionProgressDeadline      = "Failed due to progress deadline"
	DeploymentStatusDescriptionFailedByUser          = "Deployment marked as failed"
)

// DeploymentStatusDescriptionRollback is used to get the status description of
// a deployment when rolling back to an older job.
func DeploymentStatusDescriptionRollback(baseDescription string, jobVersion uint64) string {
	return fmt.Sprintf("%s - rolling back to job version %d", baseDescription, jobVersion)
}

// DeploymentStatusDescriptionRollbackNoop is used to get the status description of
// a deployment when rolling back is not possible because it has the same specification
func DeploymentStatusDescriptionRollbackNoop(baseDescription string, jobVersion uint64) string {
	return fmt.Sprintf("%s - not rolling back to stable job version %d as current job has same specification", baseDescription, jobVersion)
}

// DeploymentStatusDescriptionNoRollbackTarget is used to get the status description of
// a deployment when there is no target to rollback to but autorevert is desired.
func DeploymentStatusDescriptionNoRollbackTarget(baseDescription string) string {
	return fmt.Sprintf("%s - no stable job version to auto revert to", baseDescription)
}

// Deployment is the object that represents a job deployment which is used to
// transition a job between versions.
type Deployment struct {
	// ID is a generated UUID for the deployment
	ID string

	// Namespace is the namespace the deployment is created in
	Namespace string

	// JobID is the job the deployment is created for
	JobID string

	// JobVersion is the version of the job at which the deployment is tracking
	JobVersion uint64

	// JobModifyIndex is the ModifyIndex of the job which the deployment is
	// tracking.
	JobModifyIndex uint64

	// JobSpecModifyIndex is the JobModifyIndex of the job which the
	// deployment is tracking.
	JobSpecModifyIndex uint64

	// JobCreateIndex is the create index of the job which the deployment is
	// tracking. It is needed so that if the job gets stopped and reran we can
	// present the correct list of deployments for the job and not old ones.
	JobCreateIndex uint64

	// TaskGroups is the set of task groups effected by the deployment and their
	// current deployment status.
	TaskGroups map[string]*DeploymentState

	// The status of the deployment
	Status string

	// StatusDescription allows a human readable description of the deployment
	// status.
	StatusDescription string

	CreateIndex uint64
	ModifyIndex uint64
}

// NewDeployment creates a new deployment given the job.
func NewDeployment(job *Job) *Deployment {
	return &Deployment{
		ID:                 uuid.Generate(),
		Namespace:          job.Namespace,
		JobID:              job.ID,
		JobVersion:         job.Version,
		JobModifyIndex:     job.ModifyIndex,
		JobSpecModifyIndex: job.JobModifyIndex,
		JobCreateIndex:     job.CreateIndex,
		Status:             DeploymentStatusRunning,
		StatusDescription:  DeploymentStatusDescriptionRunning,
		TaskGroups:         make(map[string]*DeploymentState, len(job.TaskGroups)),
	}
}

func (d *Deployment) Copy() *Deployment {
	if d == nil {
		return nil
	}

	c := &Deployment{}
	*c = *d

	c.TaskGroups = nil
	if l := len(d.TaskGroups); d.TaskGroups != nil {
		c.TaskGroups = make(map[string]*DeploymentState, l)
		for tg, s := range d.TaskGroups {
			c.TaskGroups[tg] = s.Copy()
		}
	}

	return c
}

// Active returns whether the deployment is active or terminal.
func (d *Deployment) Active() bool {
	switch d.Status {
	case DeploymentStatusRunning, DeploymentStatusPaused:
		return true
	default:
		return false
	}
}

// GetID is a helper for getting the ID when the object may be nil
func (d *Deployment) GetID() string {
	if d == nil {
		return ""
	}
	return d.ID
}

// HasPlacedCanaries returns whether the deployment has placed canaries
func (d *Deployment) HasPlacedCanaries() bool {
	if d == nil || len(d.TaskGroups) == 0 {
		return false
	}
	for _, group := range d.TaskGroups {
		if len(group.PlacedCanaries) != 0 {
			return true
		}
	}
	return false
}

// RequiresPromotion returns whether the deployment requires promotion to
// continue
func (d *Deployment) RequiresPromotion() bool {
	if d == nil || len(d.TaskGroups) == 0 || d.Status != DeploymentStatusRunning {
		return false
	}
	for _, group := range d.TaskGroups {
		if group.DesiredCanaries > 0 && !group.Promoted {
			return true
		}
	}
	return false
}

func (d *Deployment) GoString() string {
	base := fmt.Sprintf("Deployment ID %q for job %q has status %q (%v):", d.ID, d.JobID, d.Status, d.StatusDescription)
	for group, state := range d.TaskGroups {
		base += fmt.Sprintf("\nTask Group %q has state:\n%#v", group, state)
	}
	return base
}

// DeploymentState tracks the state of a deployment for a given task group.
type DeploymentState struct {
	// AutoRevert marks whether the task group has indicated the job should be
	// reverted on failure
	AutoRevert bool

	// ProgressDeadline is the deadline by which an allocation must transition
	// to healthy before the deployment is considered failed.
	ProgressDeadline time.Duration

	// RequireProgressBy is the time by which an allocation must transition
	// to healthy before the deployment is considered failed.
	RequireProgressBy time.Time

	// Promoted marks whether the canaries have been promoted
	Promoted bool

	// PlacedCanaries is the set of placed canary allocations
	PlacedCanaries []string

	// DesiredCanaries is the number of canaries that should be created.
	DesiredCanaries int

	// DesiredTotal is the total number of allocations that should be created as
	// part of the deployment.
	DesiredTotal int

	// PlacedAllocs is the number of allocations that have been placed
	PlacedAllocs int

	// HealthyAllocs is the number of allocations that have been marked healthy.
	HealthyAllocs int

	// UnhealthyAllocs are allocations that have been marked as unhealthy.
	UnhealthyAllocs int
}

func (d *DeploymentState) GoString() string {
	base := fmt.Sprintf("\tDesired Total: %d", d.DesiredTotal)
	base += fmt.Sprintf("\n\tDesired Canaries: %d", d.DesiredCanaries)
	base += fmt.Sprintf("\n\tPlaced Canaries: %#v", d.PlacedCanaries)
	base += fmt.Sprintf("\n\tPromoted: %v", d.Promoted)
	base += fmt.Sprintf("\n\tPlaced: %d", d.PlacedAllocs)
	base += fmt.Sprintf("\n\tHealthy: %d", d.HealthyAllocs)
	base += fmt.Sprintf("\n\tUnhealthy: %d", d.UnhealthyAllocs)
	base += fmt.Sprintf("\n\tAutoRevert: %v", d.AutoRevert)
	return base
}

func (d *DeploymentState) Copy() *DeploymentState {
	c := &DeploymentState{}
	*c = *d
	c.PlacedCanaries = helper.CopySliceString(d.PlacedCanaries)
	return c
}

// DeploymentStatusUpdate is used to update the status of a given deployment
type DeploymentStatusUpdate struct {
	// DeploymentID is the ID of the deployment to update
	DeploymentID string

	// Status is the new status of the deployment.
	Status string

	// StatusDescription is the new status description of the deployment.
	StatusDescription string
}

// RescheduleTracker encapsulates previous reschedule events
type RescheduleTracker struct {
	Events []*RescheduleEvent
}

func (rt *RescheduleTracker) Copy() *RescheduleTracker {
	if rt == nil {
		return nil
	}
	nt := &RescheduleTracker{}
	*nt = *rt
	rescheduleEvents := make([]*RescheduleEvent, 0, len(rt.Events))
	for _, tracker := range rt.Events {
		rescheduleEvents = append(rescheduleEvents, tracker.Copy())
	}
	nt.Events = rescheduleEvents
	return nt
}

// RescheduleEvent is used to keep track of previous attempts at rescheduling an allocation
type RescheduleEvent struct {
	// RescheduleTime is the timestamp of a reschedule attempt
	RescheduleTime int64

	// PrevAllocID is the ID of the previous allocation being restarted
	PrevAllocID string

	// PrevNodeID is the node ID of the previous allocation
	PrevNodeID string

	// Delay is the reschedule delay associated with the attempt
	Delay time.Duration
}

func NewRescheduleEvent(rescheduleTime int64, prevAllocID string, prevNodeID string, delay time.Duration) *RescheduleEvent {
	return &RescheduleEvent{RescheduleTime: rescheduleTime,
		PrevAllocID: prevAllocID,
		PrevNodeID:  prevNodeID,
		Delay:       delay}
}

func (re *RescheduleEvent) Copy() *RescheduleEvent {
	if re == nil {
		return nil
	}
	copy := new(RescheduleEvent)
	*copy = *re
	return copy
}

// DesiredTransition is used to mark an allocation as having a desired state
// transition. This information can be used by the scheduler to make the
// correct decision.
type DesiredTransition struct {
	// Migrate is used to indicate that this allocation should be stopped and
	// migrated to another node.
	Migrate *bool

	// Reschedule is used to indicate that this allocation is eligible to be
	// rescheduled. Most allocations are automatically eligible for
	// rescheduling, so this field is only required when an allocation is not
	// automatically eligible. An example is an allocation that is part of a
	// deployment.
	Reschedule *bool

	// ForceReschedule is used to indicate that this allocation must be rescheduled.
	// This field is only used when operators want to force a placement even if
	// a failed allocation is not eligible to be rescheduled
	ForceReschedule *bool
}

// Merge merges the two desired transitions, preferring the values from the
// passed in object.
func (d *DesiredTransition) Merge(o *DesiredTransition) {
	if o.Migrate != nil {
		d.Migrate = o.Migrate
	}

	if o.Reschedule != nil {
		d.Reschedule = o.Reschedule
	}

	if o.ForceReschedule != nil {
		d.ForceReschedule = o.ForceReschedule
	}
}

// ShouldMigrate returns whether the transition object dictates a migration.
func (d *DesiredTransition) ShouldMigrate() bool {
	return d.Migrate != nil && *d.Migrate
}

// ShouldReschedule returns whether the transition object dictates a
// rescheduling.
func (d *DesiredTransition) ShouldReschedule() bool {
	return d.Reschedule != nil && *d.Reschedule
}

// ShouldForceReschedule returns whether the transition object dictates a
// forced rescheduling.
func (d *DesiredTransition) ShouldForceReschedule() bool {
	if d == nil {
		return false
	}
	return d.ForceReschedule != nil && *d.ForceReschedule
}

const (
	AllocDesiredStatusRun   = "run"   // Allocation should run
	AllocDesiredStatusStop  = "stop"  // Allocation should stop
	AllocDesiredStatusEvict = "evict" // Allocation should stop, and was evicted
)

const (
	AllocClientStatusPending  = "pending"
	AllocClientStatusRunning  = "running"
	AllocClientStatusComplete = "complete"
	AllocClientStatusFailed   = "failed"
	AllocClientStatusLost     = "lost"
)

// Allocation is used to allocate the placement of a task group to a node.
type Allocation struct {
	// ID of the allocation (UUID)
	ID string

	// Namespace is the namespace the allocation is created in
	Namespace string

	// ID of the evaluation that generated this allocation
	EvalID string

	// Name is a logical name of the allocation.
	Name string

	// NodeID is the node this is being placed on
	NodeID string

	// Job is the parent job of the task group being allocated.
	// This is copied at allocation time to avoid issues if the job
	// definition is updated.
	JobID string
	Job   *Job

	// TaskGroup is the name of the task group that should be run
	TaskGroup string

	// Resources is the total set of resources allocated as part
	// of this allocation of the task group.
	Resources *Resources

	// SharedResources are the resources that are shared by all the tasks in an
	// allocation
	SharedResources *Resources

	// TaskResources is the set of resources allocated to each
	// task. These should sum to the total Resources.
	TaskResources map[string]*Resources

	// Metrics associated with this allocation
	Metrics *AllocMetric

	// Desired Status of the allocation on the client
	DesiredStatus string

	// DesiredStatusDescription is meant to provide more human useful information
	DesiredDescription string

	// DesiredTransition is used to indicate that a state transition
	// is desired for a given reason.
	DesiredTransition DesiredTransition

	// Status of the allocation on the client
	ClientStatus string

	// ClientStatusDescription is meant to provide more human useful information
	ClientDescription string

	// TaskStates stores the state of each task,
	TaskStates map[string]*TaskState

	// PreviousAllocation is the allocation that this allocation is replacing
	PreviousAllocation string

	// NextAllocation is the allocation that this allocation is being replaced by
	NextAllocation string

	// DeploymentID identifies an allocation as being created from a
	// particular deployment
	DeploymentID string

	// DeploymentStatus captures the status of the allocation as part of the
	// given deployment
	DeploymentStatus *AllocDeploymentStatus

	// RescheduleTrackers captures details of previous reschedule attempts of the allocation
	RescheduleTracker *RescheduleTracker

	// FollowupEvalID captures a follow up evaluation created to handle a failed allocation
	// that can be rescheduled in the future
	FollowupEvalID string

	// Raft Indexes
	CreateIndex uint64
	ModifyIndex uint64

	// AllocModifyIndex is not updated when the client updates allocations. This
	// lets the client pull only the allocs updated by the server.
	AllocModifyIndex uint64

	// CreateTime is the time the allocation has finished scheduling and been
	// verified by the plan applier.
	CreateTime int64

	// ModifyTime is the time the allocation was last updated.
	ModifyTime int64
}

// Index returns the index of the allocation. If the allocation is from a task
// group with count greater than 1, there will be multiple allocations for it.
func (a *Allocation) Index() uint {
	l := len(a.Name)
	prefix := len(a.JobID) + len(a.TaskGroup) + 2
	if l <= 3 || l <= prefix {
		return uint(0)
	}

	strNum := a.Name[prefix : len(a.Name)-1]
	num, _ := strconv.Atoi(strNum)
	return uint(num)
}

func (a *Allocation) Copy() *Allocation {
	return a.copyImpl(true)
}

// Copy provides a copy of the allocation but doesn't deep copy the job
func (a *Allocation) CopySkipJob() *Allocation {
	return a.copyImpl(false)
}

func (a *Allocation) copyImpl(job bool) *Allocation {
	if a == nil {
		return nil
	}
	na := new(Allocation)
	*na = *a

	if job {
		na.Job = na.Job.Copy()
	}

	na.Resources = na.Resources.Copy()
	na.SharedResources = na.SharedResources.Copy()

	if a.TaskResources != nil {
		tr := make(map[string]*Resources, len(na.TaskResources))
		for task, resource := range na.TaskResources {
			tr[task] = resource.Copy()
		}
		na.TaskResources = tr
	}

	na.Metrics = na.Metrics.Copy()
	na.DeploymentStatus = na.DeploymentStatus.Copy()

	if a.TaskStates != nil {
		ts := make(map[string]*TaskState, len(na.TaskStates))
		for task, state := range na.TaskStates {
			ts[task] = state.Copy()
		}
		na.TaskStates = ts
	}

	na.RescheduleTracker = a.RescheduleTracker.Copy()
	return na
}

// TerminalStatus returns if the desired or actual status is terminal and
// will no longer transition.
func (a *Allocation) TerminalStatus() bool {
	// First check the desired state and if that isn't terminal, check client
	// state.
	switch a.DesiredStatus {
	case AllocDesiredStatusStop, AllocDesiredStatusEvict:
		return true
	default:
	}

	return a.ClientTerminalStatus()
}

// ClientTerminalStatus returns if the client status is terminal and will no longer transition
func (a *Allocation) ClientTerminalStatus() bool {
	switch a.ClientStatus {
	case AllocClientStatusComplete, AllocClientStatusFailed, AllocClientStatusLost:
		return true
	default:
		return false
	}
}

// ShouldReschedule returns if the allocation is eligible to be rescheduled according
// to its status and ReschedulePolicy given its failure time
func (a *Allocation) ShouldReschedule(reschedulePolicy *ReschedulePolicy, failTime time.Time) bool {
	// First check the desired state
	switch a.DesiredStatus {
	case AllocDesiredStatusStop, AllocDesiredStatusEvict:
		return false
	default:
	}
	switch a.ClientStatus {
	case AllocClientStatusFailed:
		return a.RescheduleEligible(reschedulePolicy, failTime)
	default:
		return false
	}
}

// RescheduleEligible returns if the allocation is eligible to be rescheduled according
// to its ReschedulePolicy and the current state of its reschedule trackers
func (a *Allocation) RescheduleEligible(reschedulePolicy *ReschedulePolicy, failTime time.Time) bool {
	if reschedulePolicy == nil {
		return false
	}
	attempts := reschedulePolicy.Attempts
	interval := reschedulePolicy.Interval
	enabled := attempts > 0 || reschedulePolicy.Unlimited
	if !enabled {
		return false
	}
	if reschedulePolicy.Unlimited {
		return true
	}
	// Early return true if there are no attempts yet and the number of allowed attempts is > 0
	if (a.RescheduleTracker == nil || len(a.RescheduleTracker.Events) == 0) && attempts > 0 {
		return true
	}
	attempted := 0
	for j := len(a.RescheduleTracker.Events) - 1; j >= 0; j-- {
		lastAttempt := a.RescheduleTracker.Events[j].RescheduleTime
		timeDiff := failTime.UTC().UnixNano() - lastAttempt
		if timeDiff < interval.Nanoseconds() {
			attempted += 1
		}
	}
	return attempted < attempts
}

// LastEventTime is the time of the last task event in the allocation.
// It is used to determine allocation failure time. If the FinishedAt field
// is not set, the alloc's modify time is used
func (a *Allocation) LastEventTime() time.Time {
	var lastEventTime time.Time
	if a.TaskStates != nil {
		for _, s := range a.TaskStates {
			if lastEventTime.IsZero() || s.FinishedAt.After(lastEventTime) {
				lastEventTime = s.FinishedAt
			}
		}
	}

	if lastEventTime.IsZero() {
		return time.Unix(0, a.ModifyTime).UTC()
	}
	return lastEventTime
}

// ReschedulePolicy returns the reschedule policy based on the task group
func (a *Allocation) ReschedulePolicy() *ReschedulePolicy {
	tg := a.Job.LookupTaskGroup(a.TaskGroup)
	if tg == nil {
		return nil
	}
	return tg.ReschedulePolicy
}

// NextRescheduleTime returns a time on or after which the allocation is eligible to be rescheduled,
// and whether the next reschedule time is within policy's interval if the policy doesn't allow unlimited reschedules
func (a *Allocation) NextRescheduleTime() (time.Time, bool) {
	failTime := a.LastEventTime()
	reschedulePolicy := a.ReschedulePolicy()
	if a.DesiredStatus == AllocDesiredStatusStop || a.ClientStatus != AllocClientStatusFailed || failTime.IsZero() || reschedulePolicy == nil {
		return time.Time{}, false
	}

	nextDelay := a.NextDelay()
	nextRescheduleTime := failTime.Add(nextDelay)
	rescheduleEligible := reschedulePolicy.Unlimited || (reschedulePolicy.Attempts > 0 && a.RescheduleTracker == nil)
	if reschedulePolicy.Attempts > 0 && a.RescheduleTracker != nil && a.RescheduleTracker.Events != nil {
		// Check for eligibility based on the interval if max attempts is set
		attempted := 0
		for j := len(a.RescheduleTracker.Events) - 1; j >= 0; j-- {
			lastAttempt := a.RescheduleTracker.Events[j].RescheduleTime
			timeDiff := failTime.UTC().UnixNano() - lastAttempt
			if timeDiff < reschedulePolicy.Interval.Nanoseconds() {
				attempted += 1
			}
		}
		rescheduleEligible = attempted < reschedulePolicy.Attempts && nextDelay < reschedulePolicy.Interval
	}
	return nextRescheduleTime, rescheduleEligible
}

// NextDelay returns a duration after which the allocation can be rescheduled.
// It is calculated according to the delay function and previous reschedule attempts.
func (a *Allocation) NextDelay() time.Duration {
	policy := a.ReschedulePolicy()
	delayDur := policy.Delay
	if a.RescheduleTracker == nil || a.RescheduleTracker.Events == nil || len(a.RescheduleTracker.Events) == 0 {
		return delayDur
	}
	events := a.RescheduleTracker.Events
	switch policy.DelayFunction {
	case "exponential":
		delayDur = a.RescheduleTracker.Events[len(a.RescheduleTracker.Events)-1].Delay * 2
	case "fibonacci":
		if len(events) >= 2 {
			fibN1Delay := events[len(events)-1].Delay
			fibN2Delay := events[len(events)-2].Delay
			// Handle reset of delay ceiling which should cause
			// a new series to start
			if fibN2Delay == policy.MaxDelay && fibN1Delay == policy.Delay {
				delayDur = fibN1Delay
			} else {
				delayDur = fibN1Delay + fibN2Delay
			}
		}
	default:
		return delayDur
	}
	if policy.MaxDelay > 0 && delayDur > policy.MaxDelay {
		delayDur = policy.MaxDelay
		// check if delay needs to be reset

		lastRescheduleEvent := a.RescheduleTracker.Events[len(a.RescheduleTracker.Events)-1]
		timeDiff := a.LastEventTime().UTC().UnixNano() - lastRescheduleEvent.RescheduleTime
		if timeDiff > delayDur.Nanoseconds() {
			delayDur = policy.Delay
		}

	}

	return delayDur
}

// Terminated returns if the allocation is in a terminal state on a client.
func (a *Allocation) Terminated() bool {
	if a.ClientStatus == AllocClientStatusFailed ||
		a.ClientStatus == AllocClientStatusComplete ||
		a.ClientStatus == AllocClientStatusLost {
		return true
	}
	return false
}

// RanSuccessfully returns whether the client has ran the allocation and all
// tasks finished successfully. Critically this function returns whether the
// allocation has ran to completion and not just that the alloc has converged to
// its desired state. That is to say that a batch allocation must have finished
// with exit code 0 on all task groups. This doesn't really have meaning on a
// non-batch allocation because a service and system allocation should not
// finish.
func (a *Allocation) RanSuccessfully() bool {
	// Handle the case the client hasn't started the allocation.
	if len(a.TaskStates) == 0 {
		return false
	}

	// Check to see if all the tasks finished successfully in the allocation
	allSuccess := true
	for _, state := range a.TaskStates {
		allSuccess = allSuccess && state.Successful()
	}

	return allSuccess
}

// ShouldMigrate returns if the allocation needs data migration
func (a *Allocation) ShouldMigrate() bool {
	if a.PreviousAllocation == "" {
		return false
	}

	if a.DesiredStatus == AllocDesiredStatusStop || a.DesiredStatus == AllocDesiredStatusEvict {
		return false
	}

	tg := a.Job.LookupTaskGroup(a.TaskGroup)

	// if the task group is nil or the ephemeral disk block isn't present then
	// we won't migrate
	if tg == nil || tg.EphemeralDisk == nil {
		return false
	}

	// We won't migrate any data is the user hasn't enabled migration or the
	// disk is not marked as sticky
	if !tg.EphemeralDisk.Migrate || !tg.EphemeralDisk.Sticky {
		return false
	}

	return true
}

// SetEventDisplayMessage populates the display message if its not already set,
// a temporary fix to handle old allocations that don't have it.
// This method will be removed in a future release.
func (a *Allocation) SetEventDisplayMessages() {
	setDisplayMsg(a.TaskStates)
}

// Stub returns a list stub for the allocation
func (a *Allocation) Stub() *AllocListStub {
	return &AllocListStub{
		ID:                 a.ID,
		EvalID:             a.EvalID,
		Name:               a.Name,
		NodeID:             a.NodeID,
		JobID:              a.JobID,
		JobVersion:         a.Job.Version,
		TaskGroup:          a.TaskGroup,
		DesiredStatus:      a.DesiredStatus,
		DesiredDescription: a.DesiredDescription,
		ClientStatus:       a.ClientStatus,
		ClientDescription:  a.ClientDescription,
		DesiredTransition:  a.DesiredTransition,
		TaskStates:         a.TaskStates,
		DeploymentStatus:   a.DeploymentStatus,
		FollowupEvalID:     a.FollowupEvalID,
		RescheduleTracker:  a.RescheduleTracker,
		CreateIndex:        a.CreateIndex,
		ModifyIndex:        a.ModifyIndex,
		CreateTime:         a.CreateTime,
		ModifyTime:         a.ModifyTime,
	}
}

// AllocListStub is used to return a subset of alloc information
type AllocListStub struct {
	ID                 string
	EvalID             string
	Name               string
	NodeID             string
	JobID              string
	JobVersion         uint64
	TaskGroup          string
	DesiredStatus      string
	DesiredDescription string
	ClientStatus       string
	ClientDescription  string
	DesiredTransition  DesiredTransition
	TaskStates         map[string]*TaskState
	DeploymentStatus   *AllocDeploymentStatus
	FollowupEvalID     string
	RescheduleTracker  *RescheduleTracker
	CreateIndex        uint64
	ModifyIndex        uint64
	CreateTime         int64
	ModifyTime         int64
}

// SetEventDisplayMessage populates the display message if its not already set,
// a temporary fix to handle old allocations that don't have it.
// This method will be removed in a future release.
func (a *AllocListStub) SetEventDisplayMessages() {
	setDisplayMsg(a.TaskStates)
}

func setDisplayMsg(taskStates map[string]*TaskState) {
	if taskStates != nil {
		for _, taskState := range taskStates {
			for _, event := range taskState.Events {
				event.PopulateEventDisplayMessage()
			}
		}
	}
}

// AllocMetric is used to track various metrics while attempting
// to make an allocation. These are used to debug a job, or to better
// understand the pressure within the system.
type AllocMetric struct {
	// NodesEvaluated is the number of nodes that were evaluated
	NodesEvaluated int

	// NodesFiltered is the number of nodes filtered due to a constraint
	NodesFiltered int

	// NodesAvailable is the number of nodes available for evaluation per DC.
	NodesAvailable map[string]int

	// ClassFiltered is the number of nodes filtered by class
	ClassFiltered map[string]int

	// ConstraintFiltered is the number of failures caused by constraint
	ConstraintFiltered map[string]int

	// NodesExhausted is the number of nodes skipped due to being
	// exhausted of at least one resource
	NodesExhausted int

	// ClassExhausted is the number of nodes exhausted by class
	ClassExhausted map[string]int

	// DimensionExhausted provides the count by dimension or reason
	DimensionExhausted map[string]int

	// QuotaExhausted provides the exhausted dimensions
	QuotaExhausted []string

	// Scores is the scores of the final few nodes remaining
	// for placement. The top score is typically selected.
	Scores map[string]float64

	// AllocationTime is a measure of how long the allocation
	// attempt took. This can affect performance and SLAs.
	AllocationTime time.Duration

	// CoalescedFailures indicates the number of other
	// allocations that were coalesced into this failed allocation.
	// This is to prevent creating many failed allocations for a
	// single task group.
	CoalescedFailures int
}

func (a *AllocMetric) Copy() *AllocMetric {
	if a == nil {
		return nil
	}
	na := new(AllocMetric)
	*na = *a
	na.NodesAvailable = helper.CopyMapStringInt(na.NodesAvailable)
	na.ClassFiltered = helper.CopyMapStringInt(na.ClassFiltered)
	na.ConstraintFiltered = helper.CopyMapStringInt(na.ConstraintFiltered)
	na.ClassExhausted = helper.CopyMapStringInt(na.ClassExhausted)
	na.DimensionExhausted = helper.CopyMapStringInt(na.DimensionExhausted)
	na.QuotaExhausted = helper.CopySliceString(na.QuotaExhausted)
	na.Scores = helper.CopyMapStringFloat64(na.Scores)
	return na
}

func (a *AllocMetric) EvaluateNode() {
	a.NodesEvaluated += 1
}

func (a *AllocMetric) FilterNode(node *Node, constraint string) {
	a.NodesFiltered += 1
	if node != nil && node.NodeClass != "" {
		if a.ClassFiltered == nil {
			a.ClassFiltered = make(map[string]int)
		}
		a.ClassFiltered[node.NodeClass] += 1
	}
	if constraint != "" {
		if a.ConstraintFiltered == nil {
			a.ConstraintFiltered = make(map[string]int)
		}
		a.ConstraintFiltered[constraint] += 1
	}
}

func (a *AllocMetric) ExhaustedNode(node *Node, dimension string) {
	a.NodesExhausted += 1
	if node != nil && node.NodeClass != "" {
		if a.ClassExhausted == nil {
			a.ClassExhausted = make(map[string]int)
		}
		a.ClassExhausted[node.NodeClass] += 1
	}
	if dimension != "" {
		if a.DimensionExhausted == nil {
			a.DimensionExhausted = make(map[string]int)
		}
		a.DimensionExhausted[dimension] += 1
	}
}

func (a *AllocMetric) ExhaustQuota(dimensions []string) {
	if a.QuotaExhausted == nil {
		a.QuotaExhausted = make([]string, 0, len(dimensions))
	}

	a.QuotaExhausted = append(a.QuotaExhausted, dimensions...)
}

func (a *AllocMetric) ScoreNode(node *Node, name string, score float64) {
	if a.Scores == nil {
		a.Scores = make(map[string]float64)
	}
	key := fmt.Sprintf("%s.%s", node.ID, name)
	a.Scores[key] = score
}

// AllocDeploymentStatus captures the status of the allocation as part of the
// deployment. This can include things like if the allocation has been marked as
// healthy.
type AllocDeploymentStatus struct {
	// Healthy marks whether the allocation has been marked healthy or unhealthy
	// as part of a deployment. It can be unset if it has neither been marked
	// healthy or unhealthy.
	Healthy *bool

	// Timestamp is the time at which the health status was set.
	Timestamp time.Time

	// Canary marks whether the allocation is a canary or not. A canary that has
	// been promoted will have this field set to false.
	Canary bool

	// ModifyIndex is the raft index in which the deployment status was last
	// changed.
	ModifyIndex uint64
}

// HasHealth returns true if the allocation has its health set.
func (a *AllocDeploymentStatus) HasHealth() bool {
	return a != nil && a.Healthy != nil
}

// IsHealthy returns if the allocation is marked as healthy as part of a
// deployment
func (a *AllocDeploymentStatus) IsHealthy() bool {
	if a == nil {
		return false
	}

	return a.Healthy != nil && *a.Healthy
}

// IsUnhealthy returns if the allocation is marked as unhealthy as part of a
// deployment
func (a *AllocDeploymentStatus) IsUnhealthy() bool {
	if a == nil {
		return false
	}

	return a.Healthy != nil && !*a.Healthy
}

// IsCanary returns if the allocation is marked as a canary
func (a *AllocDeploymentStatus) IsCanary() bool {
	if a == nil {
		return false
	}

	return a.Canary
}

func (a *AllocDeploymentStatus) Copy() *AllocDeploymentStatus {
	if a == nil {
		return nil
	}

	c := new(AllocDeploymentStatus)
	*c = *a

	if a.Healthy != nil {
		c.Healthy = helper.BoolToPtr(*a.Healthy)
	}

	return c
}

const (
	EvalStatusBlocked   = "blocked"
	EvalStatusPending   = "pending"
	EvalStatusComplete  = "complete"
	EvalStatusFailed    = "failed"
	EvalStatusCancelled = "canceled"
)

const (
	EvalTriggerJobRegister       = "job-register"
	EvalTriggerJobDeregister     = "job-deregister"
	EvalTriggerPeriodicJob       = "periodic-job"
	EvalTriggerNodeDrain         = "node-drain"
	EvalTriggerNodeUpdate        = "node-update"
	EvalTriggerScheduled         = "scheduled"
	EvalTriggerRollingUpdate     = "rolling-update"
	EvalTriggerDeploymentWatcher = "deployment-watcher"
	EvalTriggerFailedFollowUp    = "failed-follow-up"
	EvalTriggerMaxPlans          = "max-plan-attempts"
	EvalTriggerRetryFailedAlloc  = "alloc-failure"
)

const (
	// CoreJobEvalGC is used for the garbage collection of evaluations
	// and allocations. We periodically scan evaluations in a terminal state,
	// in which all the corresponding allocations are also terminal. We
	// delete these out of the system to bound the state.
	CoreJobEvalGC = "eval-gc"

	// CoreJobNodeGC is used for the garbage collection of failed nodes.
	// We periodically scan nodes in a terminal state, and if they have no
	// corresponding allocations we delete these out of the system.
	CoreJobNodeGC = "node-gc"

	// CoreJobJobGC is used for the garbage collection of eligible jobs. We
	// periodically scan garbage collectible jobs and check if both their
	// evaluations and allocations are terminal. If so, we delete these out of
	// the system.
	CoreJobJobGC = "job-gc"

	// CoreJobDeploymentGC is used for the garbage collection of eligible
	// deployments. We periodically scan garbage collectible deployments and
	// check if they are terminal. If so, we delete these out of the system.
	CoreJobDeploymentGC = "deployment-gc"

	// CoreJobForceGC is used to force garbage collection of all GCable objects.
	CoreJobForceGC = "force-gc"
)

// Evaluation is used anytime we need to apply business logic as a result
// of a change to our desired state (job specification) or the emergent state
// (registered nodes). When the inputs change, we need to "evaluate" them,
// potentially taking action (allocation of work) or doing nothing if the state
// of the world does not require it.
type Evaluation struct {
	// ID is a randomly generated UUID used for this evaluation. This
	// is assigned upon the creation of the evaluation.
	ID string

	// Namespace is the namespace the evaluation is created in
	Namespace string

	// Priority is used to control scheduling importance and if this job
	// can preempt other jobs.
	Priority int

	// Type is used to control which schedulers are available to handle
	// this evaluation.
	Type string

	// TriggeredBy is used to give some insight into why this Eval
	// was created. (Job change, node failure, alloc failure, etc).
	TriggeredBy string

	// JobID is the job this evaluation is scoped to. Evaluations cannot
	// be run in parallel for a given JobID, so we serialize on this.
	JobID string

	// JobModifyIndex is the modify index of the job at the time
	// the evaluation was created
	JobModifyIndex uint64

	// NodeID is the node that was affected triggering the evaluation.
	NodeID string

	// NodeModifyIndex is the modify index of the node at the time
	// the evaluation was created
	NodeModifyIndex uint64

	// DeploymentID is the ID of the deployment that triggered the evaluation.
	DeploymentID string

	// Status of the evaluation
	Status string

	// StatusDescription is meant to provide more human useful information
	StatusDescription string

	// Wait is a minimum wait time for running the eval. This is used to
	// support a rolling upgrade in versions prior to 0.7.0
	// Deprecated
	Wait time.Duration

	// WaitUntil is the time when this eval should be run. This is used to
	// supported delayed rescheduling of failed allocations
	WaitUntil time.Time

	// NextEval is the evaluation ID for the eval created to do a followup.
	// This is used to support rolling upgrades, where we need a chain of evaluations.
	NextEval string

	// PreviousEval is the evaluation ID for the eval creating this one to do a followup.
	// This is used to support rolling upgrades, where we need a chain of evaluations.
	PreviousEval string

	// BlockedEval is the evaluation ID for a created blocked eval. A
	// blocked eval will be created if all allocations could not be placed due
	// to constraints or lacking resources.
	BlockedEval string

	// FailedTGAllocs are task groups which have allocations that could not be
	// made, but the metrics are persisted so that the user can use the feedback
	// to determine the cause.
	FailedTGAllocs map[string]*AllocMetric

	// ClassEligibility tracks computed node classes that have been explicitly
	// marked as eligible or ineligible.
	ClassEligibility map[string]bool

	// QuotaLimitReached marks whether a quota limit was reached for the
	// evaluation.
	QuotaLimitReached string

	// EscapedComputedClass marks whether the job has constraints that are not
	// captured by computed node classes.
	EscapedComputedClass bool

	// AnnotatePlan triggers the scheduler to provide additional annotations
	// during the evaluation. This should not be set during normal operations.
	AnnotatePlan bool

	// QueuedAllocations is the number of unplaced allocations at the time the
	// evaluation was processed. The map is keyed by Task Group names.
	QueuedAllocations map[string]int

	// LeaderACL provides the ACL token to when issuing RPCs back to the
	// leader. This will be a valid management token as long as the leader is
	// active. This should not ever be exposed via the API.
	LeaderACL string

	// SnapshotIndex is the Raft index of the snapshot used to process the
	// evaluation. As such it will only be set once it has gone through the
	// scheduler.
	SnapshotIndex uint64

	// Raft Indexes
	CreateIndex uint64
	ModifyIndex uint64
}

// TerminalStatus returns if the current status is terminal and
// will no longer transition.
func (e *Evaluation) TerminalStatus() bool {
	switch e.Status {
	case EvalStatusComplete, EvalStatusFailed, EvalStatusCancelled:
		return true
	default:
		return false
	}
}

func (e *Evaluation) GoString() string {
	return fmt.Sprintf("<Eval %q JobID: %q Namespace: %q>", e.ID, e.JobID, e.Namespace)
}

func (e *Evaluation) Copy() *Evaluation {
	if e == nil {
		return nil
	}
	ne := new(Evaluation)
	*ne = *e

	// Copy ClassEligibility
	if e.ClassEligibility != nil {
		classes := make(map[string]bool, len(e.ClassEligibility))
		for class, elig := range e.ClassEligibility {
			classes[class] = elig
		}
		ne.ClassEligibility = classes
	}

	// Copy FailedTGAllocs
	if e.FailedTGAllocs != nil {
		failedTGs := make(map[string]*AllocMetric, len(e.FailedTGAllocs))
		for tg, metric := range e.FailedTGAllocs {
			failedTGs[tg] = metric.Copy()
		}
		ne.FailedTGAllocs = failedTGs
	}

	// Copy queued allocations
	if e.QueuedAllocations != nil {
		queuedAllocations := make(map[string]int, len(e.QueuedAllocations))
		for tg, num := range e.QueuedAllocations {
			queuedAllocations[tg] = num
		}
		ne.QueuedAllocations = queuedAllocations
	}

	return ne
}

// ShouldEnqueue checks if a given evaluation should be enqueued into the
// eval_broker
func (e *Evaluation) ShouldEnqueue() bool {
	switch e.Status {
	case EvalStatusPending:
		return true
	case EvalStatusComplete, EvalStatusFailed, EvalStatusBlocked, EvalStatusCancelled:
		return false
	default:
		panic(fmt.Sprintf("unhandled evaluation (%s) status %s", e.ID, e.Status))
	}
}

// ShouldBlock checks if a given evaluation should be entered into the blocked
// eval tracker.
func (e *Evaluation) ShouldBlock() bool {
	switch e.Status {
	case EvalStatusBlocked:
		return true
	case EvalStatusComplete, EvalStatusFailed, EvalStatusPending, EvalStatusCancelled:
		return false
	default:
		panic(fmt.Sprintf("unhandled evaluation (%s) status %s", e.ID, e.Status))
	}
}

// MakePlan is used to make a plan from the given evaluation
// for a given Job
func (e *Evaluation) MakePlan(j *Job) *Plan {
	p := &Plan{
		EvalID:         e.ID,
		Priority:       e.Priority,
		Job:            j,
		NodeUpdate:     make(map[string][]*Allocation),
		NodeAllocation: make(map[string][]*Allocation),
	}
	if j != nil {
		p.AllAtOnce = j.AllAtOnce
	}
	return p
}

// NextRollingEval creates an evaluation to followup this eval for rolling updates
func (e *Evaluation) NextRollingEval(wait time.Duration) *Evaluation {
	return &Evaluation{
		ID:             uuid.Generate(),
		Namespace:      e.Namespace,
		Priority:       e.Priority,
		Type:           e.Type,
		TriggeredBy:    EvalTriggerRollingUpdate,
		JobID:          e.JobID,
		JobModifyIndex: e.JobModifyIndex,
		Status:         EvalStatusPending,
		Wait:           wait,
		PreviousEval:   e.ID,
	}
}

// CreateBlockedEval creates a blocked evaluation to followup this eval to place any
// failed allocations. It takes the classes marked explicitly eligible or
// ineligible, whether the job has escaped computed node classes and whether the
// quota limit was reached.
func (e *Evaluation) CreateBlockedEval(classEligibility map[string]bool,
	escaped bool, quotaReached string) *Evaluation {

	return &Evaluation{
		ID:                   uuid.Generate(),
		Namespace:            e.Namespace,
		Priority:             e.Priority,
		Type:                 e.Type,
		TriggeredBy:          e.TriggeredBy,
		JobID:                e.JobID,
		JobModifyIndex:       e.JobModifyIndex,
		Status:               EvalStatusBlocked,
		PreviousEval:         e.ID,
		ClassEligibility:     classEligibility,
		EscapedComputedClass: escaped,
		QuotaLimitReached:    quotaReached,
	}
}

// CreateFailedFollowUpEval creates a follow up evaluation when the current one
// has been marked as failed because it has hit the delivery limit and will not
// be retried by the eval_broker.
func (e *Evaluation) CreateFailedFollowUpEval(wait time.Duration) *Evaluation {
	return &Evaluation{
		ID:             uuid.Generate(),
		Namespace:      e.Namespace,
		Priority:       e.Priority,
		Type:           e.Type,
		TriggeredBy:    EvalTriggerFailedFollowUp,
		JobID:          e.JobID,
		JobModifyIndex: e.JobModifyIndex,
		Status:         EvalStatusPending,
		Wait:           wait,
		PreviousEval:   e.ID,
	}
}

// Plan is used to submit a commit plan for task allocations. These
// are submitted to the leader which verifies that resources have
// not been overcommitted before admitting the plan.
type Plan struct {
	// EvalID is the evaluation ID this plan is associated with
	EvalID string

	// EvalToken is used to prevent a split-brain processing of
	// an evaluation. There should only be a single scheduler running
	// an Eval at a time, but this could be violated after a leadership
	// transition. This unique token is used to reject plans that are
	// being submitted from a different leader.
	EvalToken string

	// Priority is the priority of the upstream job
	Priority int

	// AllAtOnce is used to control if incremental scheduling of task groups
	// is allowed or if we must do a gang scheduling of the entire job.
	// If this is false, a plan may be partially applied. Otherwise, the
	// entire plan must be able to make progress.
	AllAtOnce bool

	// Job is the parent job of all the allocations in the Plan.
	// Since a Plan only involves a single Job, we can reduce the size
	// of the plan by only including it once.
	Job *Job

	// NodeUpdate contains all the allocations for each node. For each node,
	// this is a list of the allocations to update to either stop or evict.
	NodeUpdate map[string][]*Allocation

	// NodeAllocation contains all the allocations for each node.
	// The evicts must be considered prior to the allocations.
	NodeAllocation map[string][]*Allocation

	// Annotations contains annotations by the scheduler to be used by operators
	// to understand the decisions made by the scheduler.
	Annotations *PlanAnnotations

	// Deployment is the deployment created or updated by the scheduler that
	// should be applied by the planner.
	Deployment *Deployment

	// DeploymentUpdates is a set of status updates to apply to the given
	// deployments. This allows the scheduler to cancel any unneeded deployment
	// because the job is stopped or the update block is removed.
	DeploymentUpdates []*DeploymentStatusUpdate
}

// AppendUpdate marks the allocation for eviction. The clientStatus of the
// allocation may be optionally set by passing in a non-empty value.
func (p *Plan) AppendUpdate(alloc *Allocation, desiredStatus, desiredDesc, clientStatus string) {
	newAlloc := new(Allocation)
	*newAlloc = *alloc

	// If the job is not set in the plan we are deregistering a job so we
	// extract the job from the allocation.
	if p.Job == nil && newAlloc.Job != nil {
		p.Job = newAlloc.Job
	}

	// Normalize the job
	newAlloc.Job = nil

	// Strip the resources as it can be rebuilt.
	newAlloc.Resources = nil

	newAlloc.DesiredStatus = desiredStatus
	newAlloc.DesiredDescription = desiredDesc

	if clientStatus != "" {
		newAlloc.ClientStatus = clientStatus
	}

	node := alloc.NodeID
	existing := p.NodeUpdate[node]
	p.NodeUpdate[node] = append(existing, newAlloc)
}

func (p *Plan) PopUpdate(alloc *Allocation) {
	existing := p.NodeUpdate[alloc.NodeID]
	n := len(existing)
	if n > 0 && existing[n-1].ID == alloc.ID {
		existing = existing[:n-1]
		if len(existing) > 0 {
			p.NodeUpdate[alloc.NodeID] = existing
		} else {
			delete(p.NodeUpdate, alloc.NodeID)
		}
	}
}

func (p *Plan) AppendAlloc(alloc *Allocation) {
	node := alloc.NodeID
	existing := p.NodeAllocation[node]
	p.NodeAllocation[node] = append(existing, alloc)
}

// IsNoOp checks if this plan would do nothing
func (p *Plan) IsNoOp() bool {
	return len(p.NodeUpdate) == 0 &&
		len(p.NodeAllocation) == 0 &&
		p.Deployment == nil &&
		len(p.DeploymentUpdates) == 0
}

// PlanResult is the result of a plan submitted to the leader.
type PlanResult struct {
	// NodeUpdate contains all the updates that were committed.
	NodeUpdate map[string][]*Allocation

	// NodeAllocation contains all the allocations that were committed.
	NodeAllocation map[string][]*Allocation

	// Deployment is the deployment that was committed.
	Deployment *Deployment

	// DeploymentUpdates is the set of deployment updates that were committed.
	DeploymentUpdates []*DeploymentStatusUpdate

	// RefreshIndex is the index the worker should refresh state up to.
	// This allows all evictions and allocations to be materialized.
	// If any allocations were rejected due to stale data (node state,
	// over committed) this can be used to force a worker refresh.
	RefreshIndex uint64

	// AllocIndex is the Raft index in which the evictions and
	// allocations took place. This is used for the write index.
	AllocIndex uint64
}

// IsNoOp checks if this plan result would do nothing
func (p *PlanResult) IsNoOp() bool {
	return len(p.NodeUpdate) == 0 && len(p.NodeAllocation) == 0 &&
		len(p.DeploymentUpdates) == 0 && p.Deployment == nil
}

// FullCommit is used to check if all the allocations in a plan
// were committed as part of the result. Returns if there was
// a match, and the number of expected and actual allocations.
func (p *PlanResult) FullCommit(plan *Plan) (bool, int, int) {
	expected := 0
	actual := 0
	for name, allocList := range plan.NodeAllocation {
		didAlloc, _ := p.NodeAllocation[name]
		expected += len(allocList)
		actual += len(didAlloc)
	}
	return actual == expected, expected, actual
}

// PlanAnnotations holds annotations made by the scheduler to give further debug
// information to operators.
type PlanAnnotations struct {
	// DesiredTGUpdates is the set of desired updates per task group.
	DesiredTGUpdates map[string]*DesiredUpdates
}

// DesiredUpdates is the set of changes the scheduler would like to make given
// sufficient resources and cluster capacity.
type DesiredUpdates struct {
	Ignore            uint64
	Place             uint64
	Migrate           uint64
	Stop              uint64
	InPlaceUpdate     uint64
	DestructiveUpdate uint64
	Canary            uint64
}

func (d *DesiredUpdates) GoString() string {
	return fmt.Sprintf("(place %d) (inplace %d) (destructive %d) (stop %d) (migrate %d) (ignore %d) (canary %d)",
		d.Place, d.InPlaceUpdate, d.DestructiveUpdate, d.Stop, d.Migrate, d.Ignore, d.Canary)
}

// msgpackHandle is a shared handle for encoding/decoding of structs
var MsgpackHandle = func() *codec.MsgpackHandle {
	h := &codec.MsgpackHandle{RawToString: true}

	// Sets the default type for decoding a map into a nil interface{}.
	// This is necessary in particular because we store the driver configs as a
	// nil interface{}.
	h.MapType = reflect.TypeOf(map[string]interface{}(nil))
	return h
}()

var (
	// JsonHandle and JsonHandlePretty are the codec handles to JSON encode
	// structs. The pretty handle will add indents for easier human consumption.
	JsonHandle = &codec.JsonHandle{
		HTMLCharsAsIs: true,
	}
	JsonHandlePretty = &codec.JsonHandle{
		HTMLCharsAsIs: true,
		Indent:        4,
	}
)

// TODO Figure out if we can remove this. This is our fork that is just way
// behind. I feel like its original purpose was to pin at a stable version but
// now we can accomplish this with vendoring.
var HashiMsgpackHandle = func() *hcodec.MsgpackHandle {
	h := &hcodec.MsgpackHandle{RawToString: true}

	// Sets the default type for decoding a map into a nil interface{}.
	// This is necessary in particular because we store the driver configs as a
	// nil interface{}.
	h.MapType = reflect.TypeOf(map[string]interface{}(nil))
	return h
}()

// Decode is used to decode a MsgPack encoded object
func Decode(buf []byte, out interface{}) error {
	return codec.NewDecoder(bytes.NewReader(buf), MsgpackHandle).Decode(out)
}

// Encode is used to encode a MsgPack object with type prefix
func Encode(t MessageType, msg interface{}) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(uint8(t))
	err := codec.NewEncoder(&buf, MsgpackHandle).Encode(msg)
	return buf.Bytes(), err
}

// KeyringResponse is a unified key response and can be used for install,
// remove, use, as well as listing key queries.
type KeyringResponse struct {
	Messages map[string]string
	Keys     map[string]int
	NumNodes int
}

// KeyringRequest is request objects for serf key operations.
type KeyringRequest struct {
	Key string
}

// RecoverableError wraps an error and marks whether it is recoverable and could
// be retried or it is fatal.
type RecoverableError struct {
	Err         string
	Recoverable bool
}

// NewRecoverableError is used to wrap an error and mark it as recoverable or
// not.
func NewRecoverableError(e error, recoverable bool) error {
	if e == nil {
		return nil
	}

	return &RecoverableError{
		Err:         e.Error(),
		Recoverable: recoverable,
	}
}

// WrapRecoverable wraps an existing error in a new RecoverableError with a new
// message. If the error was recoverable before the returned error is as well;
// otherwise it is unrecoverable.
func WrapRecoverable(msg string, err error) error {
	return &RecoverableError{Err: msg, Recoverable: IsRecoverable(err)}
}

func (r *RecoverableError) Error() string {
	return r.Err
}

func (r *RecoverableError) IsRecoverable() bool {
	return r.Recoverable
}

// Recoverable is an interface for errors to implement to indicate whether or
// not they are fatal or recoverable.
type Recoverable interface {
	error
	IsRecoverable() bool
}

// IsRecoverable returns true if error is a RecoverableError with
// Recoverable=true. Otherwise false is returned.
func IsRecoverable(e error) bool {
	if re, ok := e.(Recoverable); ok {
		return re.IsRecoverable()
	}
	return false
}

// WrappedServerError wraps an error and satisfies
// both the Recoverable and the ServerSideError interfaces
type WrappedServerError struct {
	Err error
}

// NewWrappedServerError is used to create a wrapped server side error
func NewWrappedServerError(e error) error {
	return &WrappedServerError{
		Err: e,
	}
}

func (r *WrappedServerError) IsRecoverable() bool {
	return IsRecoverable(r.Err)
}

func (r *WrappedServerError) Error() string {
	return r.Err.Error()
}

func (r *WrappedServerError) IsServerSide() bool {
	return true
}

// ServerSideError is an interface for errors to implement to indicate
// errors occurring after the request makes it to a server
type ServerSideError interface {
	error
	IsServerSide() bool
}

// IsServerSide returns true if error is a wrapped
// server side error
func IsServerSide(e error) bool {
	if se, ok := e.(ServerSideError); ok {
		return se.IsServerSide()
	}
	return false
}

// ACLPolicy is used to represent an ACL policy
type ACLPolicy struct {
	Name        string // Unique name
	Description string // Human readable
	Rules       string // HCL or JSON format
	Hash        []byte
	CreateIndex uint64
	ModifyIndex uint64
}

// SetHash is used to compute and set the hash of the ACL policy
func (c *ACLPolicy) SetHash() []byte {
	// Initialize a 256bit Blake2 hash (32 bytes)
	hash, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}

	// Write all the user set fields
	hash.Write([]byte(c.Name))
	hash.Write([]byte(c.Description))
	hash.Write([]byte(c.Rules))

	// Finalize the hash
	hashVal := hash.Sum(nil)

	// Set and return the hash
	c.Hash = hashVal
	return hashVal
}

func (a *ACLPolicy) Stub() *ACLPolicyListStub {
	return &ACLPolicyListStub{
		Name:        a.Name,
		Description: a.Description,
		Hash:        a.Hash,
		CreateIndex: a.CreateIndex,
		ModifyIndex: a.ModifyIndex,
	}
}

func (a *ACLPolicy) Validate() error {
	var mErr multierror.Error
	if !validPolicyName.MatchString(a.Name) {
		err := fmt.Errorf("invalid name '%s'", a.Name)
		mErr.Errors = append(mErr.Errors, err)
	}
	if _, err := acl.Parse(a.Rules); err != nil {
		err = fmt.Errorf("failed to parse rules: %v", err)
		mErr.Errors = append(mErr.Errors, err)
	}
	if len(a.Description) > maxPolicyDescriptionLength {
		err := fmt.Errorf("description longer than %d", maxPolicyDescriptionLength)
		mErr.Errors = append(mErr.Errors, err)
	}
	return mErr.ErrorOrNil()
}

// ACLPolicyListStub is used to for listing ACL policies
type ACLPolicyListStub struct {
	Name        string
	Description string
	Hash        []byte
	CreateIndex uint64
	ModifyIndex uint64
}

// ACLPolicyListRequest is used to request a list of policies
type ACLPolicyListRequest struct {
	QueryOptions
}

// ACLPolicySpecificRequest is used to query a specific policy
type ACLPolicySpecificRequest struct {
	Name string
	QueryOptions
}

// ACLPolicySetRequest is used to query a set of policies
type ACLPolicySetRequest struct {
	Names []string
	QueryOptions
}

// ACLPolicyListResponse is used for a list request
type ACLPolicyListResponse struct {
	Policies []*ACLPolicyListStub
	QueryMeta
}

// SingleACLPolicyResponse is used to return a single policy
type SingleACLPolicyResponse struct {
	Policy *ACLPolicy
	QueryMeta
}

// ACLPolicySetResponse is used to return a set of policies
type ACLPolicySetResponse struct {
	Policies map[string]*ACLPolicy
	QueryMeta
}

// ACLPolicyDeleteRequest is used to delete a set of policies
type ACLPolicyDeleteRequest struct {
	Names []string
	WriteRequest
}

// ACLPolicyUpsertRequest is used to upsert a set of policies
type ACLPolicyUpsertRequest struct {
	Policies []*ACLPolicy
	WriteRequest
}

// ACLToken represents a client token which is used to Authenticate
type ACLToken struct {
	AccessorID  string   // Public Accessor ID (UUID)
	SecretID    string   // Secret ID, private (UUID)
	Name        string   // Human friendly name
	Type        string   // Client or Management
	Policies    []string // Policies this token ties to
	Global      bool     // Global or Region local
	Hash        []byte
	CreateTime  time.Time // Time of creation
	CreateIndex uint64
	ModifyIndex uint64
}

var (
	// AnonymousACLToken is used no SecretID is provided, and the
	// request is made anonymously.
	AnonymousACLToken = &ACLToken{
		AccessorID: "anonymous",
		Name:       "Anonymous Token",
		Type:       ACLClientToken,
		Policies:   []string{"anonymous"},
		Global:     false,
	}
)

type ACLTokenListStub struct {
	AccessorID  string
	Name        string
	Type        string
	Policies    []string
	Global      bool
	Hash        []byte
	CreateTime  time.Time
	CreateIndex uint64
	ModifyIndex uint64
}

// SetHash is used to compute and set the hash of the ACL token
func (a *ACLToken) SetHash() []byte {
	// Initialize a 256bit Blake2 hash (32 bytes)
	hash, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}

	// Write all the user set fields
	hash.Write([]byte(a.Name))
	hash.Write([]byte(a.Type))
	for _, policyName := range a.Policies {
		hash.Write([]byte(policyName))
	}
	if a.Global {
		hash.Write([]byte("global"))
	} else {
		hash.Write([]byte("local"))
	}

	// Finalize the hash
	hashVal := hash.Sum(nil)

	// Set and return the hash
	a.Hash = hashVal
	return hashVal
}

func (a *ACLToken) Stub() *ACLTokenListStub {
	return &ACLTokenListStub{
		AccessorID:  a.AccessorID,
		Name:        a.Name,
		Type:        a.Type,
		Policies:    a.Policies,
		Global:      a.Global,
		Hash:        a.Hash,
		CreateTime:  a.CreateTime,
		CreateIndex: a.CreateIndex,
		ModifyIndex: a.ModifyIndex,
	}
}

// Validate is used to sanity check a token
func (a *ACLToken) Validate() error {
	var mErr multierror.Error
	if len(a.Name) > maxTokenNameLength {
		mErr.Errors = append(mErr.Errors, fmt.Errorf("token name too long"))
	}
	switch a.Type {
	case ACLClientToken:
		if len(a.Policies) == 0 {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("client token missing policies"))
		}
	case ACLManagementToken:
		if len(a.Policies) != 0 {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("management token cannot be associated with policies"))
		}
	default:
		mErr.Errors = append(mErr.Errors, fmt.Errorf("token type must be client or management"))
	}
	return mErr.ErrorOrNil()
}

// PolicySubset checks if a given set of policies is a subset of the token
func (a *ACLToken) PolicySubset(policies []string) bool {
	// Hot-path the management tokens, superset of all policies.
	if a.Type == ACLManagementToken {
		return true
	}
	associatedPolicies := make(map[string]struct{}, len(a.Policies))
	for _, policy := range a.Policies {
		associatedPolicies[policy] = struct{}{}
	}
	for _, policy := range policies {
		if _, ok := associatedPolicies[policy]; !ok {
			return false
		}
	}
	return true
}

// ACLTokenListRequest is used to request a list of tokens
type ACLTokenListRequest struct {
	GlobalOnly bool
	QueryOptions
}

// ACLTokenSpecificRequest is used to query a specific token
type ACLTokenSpecificRequest struct {
	AccessorID string
	QueryOptions
}

// ACLTokenSetRequest is used to query a set of tokens
type ACLTokenSetRequest struct {
	AccessorIDS []string
	QueryOptions
}

// ACLTokenListResponse is used for a list request
type ACLTokenListResponse struct {
	Tokens []*ACLTokenListStub
	QueryMeta
}

// SingleACLTokenResponse is used to return a single token
type SingleACLTokenResponse struct {
	Token *ACLToken
	QueryMeta
}

// ACLTokenSetResponse is used to return a set of token
type ACLTokenSetResponse struct {
	Tokens map[string]*ACLToken // Keyed by Accessor ID
	QueryMeta
}

// ResolveACLTokenRequest is used to resolve a specific token
type ResolveACLTokenRequest struct {
	SecretID string
	QueryOptions
}

// ResolveACLTokenResponse is used to resolve a single token
type ResolveACLTokenResponse struct {
	Token *ACLToken
	QueryMeta
}

// ACLTokenDeleteRequest is used to delete a set of tokens
type ACLTokenDeleteRequest struct {
	AccessorIDs []string
	WriteRequest
}

// ACLTokenBootstrapRequest is used to bootstrap ACLs
type ACLTokenBootstrapRequest struct {
	Token      *ACLToken // Not client specifiable
	ResetIndex uint64    // Reset index is used to clear the bootstrap token
	WriteRequest
}

// ACLTokenUpsertRequest is used to upsert a set of tokens
type ACLTokenUpsertRequest struct {
	Tokens []*ACLToken
	WriteRequest
}

// ACLTokenUpsertResponse is used to return from an ACLTokenUpsertRequest
type ACLTokenUpsertResponse struct {
	Tokens []*ACLToken
	WriteMeta
}
