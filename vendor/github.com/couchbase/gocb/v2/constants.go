package gocb

import (
	gocbcore "github.com/couchbase/gocbcore/v9"
	"github.com/couchbase/gocbcore/v9/memd"
)

const (
	goCbVersionStr = "v2.1.1"

	persistenceTimeoutFloor = 1500
)

// QueryIndexType provides information on the type of indexer used for an index.
type QueryIndexType string

const (
	// QueryIndexTypeGsi indicates that GSI was used to build the index.
	QueryIndexTypeGsi = QueryIndexType("gsi")

	// QueryIndexTypeView indicates that views were used to build the index.
	QueryIndexTypeView = QueryIndexType("views")
)

// QueryStatus provides information about the current status of a query.
type QueryStatus string

const (
	// QueryStatusRunning indicates the query is still running
	QueryStatusRunning = QueryStatus("running")

	// QueryStatusSuccess indicates the query was successful.
	QueryStatusSuccess = QueryStatus("success")

	// QueryStatusErrors indicates a query completed with errors.
	QueryStatusErrors = QueryStatus("errors")

	// QueryStatusCompleted indicates a query has completed.
	QueryStatusCompleted = QueryStatus("completed")

	// QueryStatusStopped indicates a query has been stopped.
	QueryStatusStopped = QueryStatus("stopped")

	// QueryStatusTimeout indicates a query timed out.
	QueryStatusTimeout = QueryStatus("timeout")

	// QueryStatusClosed indicates that a query was closed.
	QueryStatusClosed = QueryStatus("closed")

	// QueryStatusFatal indicates that a query ended with a fatal error.
	QueryStatusFatal = QueryStatus("fatal")

	// QueryStatusAborted indicates that a query was aborted.
	QueryStatusAborted = QueryStatus("aborted")

	// QueryStatusUnknown indicates that the query status is unknown.
	QueryStatusUnknown = QueryStatus("unknown")
)

// ServiceType specifies a particular Couchbase service type.
type ServiceType gocbcore.ServiceType

const (
	// ServiceTypeManagement represents a management service.
	ServiceTypeManagement = ServiceType(gocbcore.MgmtService)

	// ServiceTypeKeyValue represents a memcached service.
	ServiceTypeKeyValue = ServiceType(gocbcore.MemdService)

	// ServiceTypeViews represents a views service.
	ServiceTypeViews = ServiceType(gocbcore.CapiService)

	// ServiceTypeQuery represents a query service.
	ServiceTypeQuery = ServiceType(gocbcore.N1qlService)

	// ServiceTypeSearch represents a full-text-search service.
	ServiceTypeSearch = ServiceType(gocbcore.FtsService)

	// ServiceTypeAnalytics represents an analytics service.
	ServiceTypeAnalytics = ServiceType(gocbcore.CbasService)
)

// QueryProfileMode specifies the profiling mode to use during a query.
type QueryProfileMode string

const (
	// QueryProfileModeNone disables query profiling
	QueryProfileModeNone = QueryProfileMode("off")

	// QueryProfileModePhases includes phase profiling information in the query response
	QueryProfileModePhases = QueryProfileMode("phases")

	// QueryProfileModeTimings includes timing profiling information in the query response
	QueryProfileModeTimings = QueryProfileMode("timings")
)

// SubdocFlag provides special handling flags for sub-document operations
type SubdocFlag memd.SubdocFlag

const (
	// SubdocFlagNone indicates no special behaviours
	SubdocFlagNone = SubdocFlag(memd.SubdocFlagNone)

	// SubdocFlagCreatePath indicates you wish to recursively create the tree of paths
	// if it does not already exist within the document.
	SubdocFlagCreatePath = SubdocFlag(memd.SubdocFlagMkDirP)

	// SubdocFlagXattr indicates your path refers to an extended attribute rather than the document.
	SubdocFlagXattr = SubdocFlag(memd.SubdocFlagXattrPath)

	// SubdocFlagUseMacros indicates that you wish macro substitution to occur on the value
	SubdocFlagUseMacros = SubdocFlag(memd.SubdocFlagExpandMacros)
)

// SubdocDocFlag specifies document-level flags for a sub-document operation.
type SubdocDocFlag memd.SubdocDocFlag

const (
	// SubdocDocFlagNone indicates no special behaviours
	SubdocDocFlagNone = SubdocDocFlag(memd.SubdocDocFlagNone)

	// SubdocDocFlagMkDoc indicates that the document should be created if it does not already exist.
	SubdocDocFlagMkDoc = SubdocDocFlag(memd.SubdocDocFlagMkDoc)

	// SubdocDocFlagAddDoc indices that the document should be created only if it does not already exist.
	SubdocDocFlagAddDoc = SubdocDocFlag(memd.SubdocDocFlagAddDoc)

	// SubdocDocFlagAccessDeleted indicates that you wish to receive soft-deleted documents.
	SubdocDocFlagAccessDeleted = SubdocDocFlag(memd.SubdocDocFlagAccessDeleted)
)

// DurabilityLevel specifies the level of synchronous replication to use.
type DurabilityLevel uint8

const (
	// DurabilityLevelMajority specifies that a mutation must be replicated (held in memory) to a majority of nodes.
	DurabilityLevelMajority = DurabilityLevel(1)

	// DurabilityLevelMajorityAndPersistOnMaster specifies that a mutation must be replicated (held in memory) to a
	// majority of nodes and also persisted (written to disk) on the active node.
	DurabilityLevelMajorityAndPersistOnMaster = DurabilityLevel(2)

	// DurabilityLevelPersistToMajority specifies that a mutation must be persisted (written to disk) to a majority
	// of nodes.
	DurabilityLevelPersistToMajority = DurabilityLevel(3)
)

// MutationMacro can be supplied to MutateIn operations to perform ExpandMacros operations.
type MutationMacro string

const (
	// MutationMacroCAS can be used to tell the server to use the CAS macro.
	MutationMacroCAS = MutationMacro("\"${Mutation.CAS}\"")

	// MutationMacroSeqNo can be used to tell the server to use the seqno macro.
	MutationMacroSeqNo = MutationMacro("\"${Mutation.seqno}\"")

	// MutationMacroValueCRC32c can be used to tell the server to use the value_crc32c macro.
	MutationMacroValueCRC32c = MutationMacro("\"${Mutation.value_crc32c}\"")
)

// ClusterState specifies the current state of the cluster
type ClusterState uint

const (
	// ClusterStateOnline indicates that all nodes are online and reachable.
	ClusterStateOnline = ClusterState(1)

	// ClusterStateDegraded indicates that all services will function, but possibly not optimally.
	ClusterStateDegraded = ClusterState(2)

	// ClusterStateOffline indicates that no nodes were reachable.
	ClusterStateOffline = ClusterState(3)
)

// EndpointState specifies the current state of an endpoint.
type EndpointState uint

const (
	// EndpointStateDisconnected indicates the endpoint socket is unreachable.
	EndpointStateDisconnected = EndpointState(1)

	// EndpointStateConnecting indicates the endpoint socket is connecting.
	EndpointStateConnecting = EndpointState(2)

	// EndpointStateConnected indicates the endpoint socket is connected and ready.
	EndpointStateConnected = EndpointState(3)

	// EndpointStateDisconnecting indicates the endpoint socket is disconnecting.
	EndpointStateDisconnecting = EndpointState(4)
)

// PingState specifies the result of the ping operation
type PingState uint

const (
	// PingStateOk indicates that the ping operation was successful.
	PingStateOk = PingState(1)

	// PingStateTimeout indicates that the ping operation timed out.
	PingStateTimeout = PingState(2)

	// PingStateError indicates that the ping operation failed.
	PingStateError = PingState(3)
)
