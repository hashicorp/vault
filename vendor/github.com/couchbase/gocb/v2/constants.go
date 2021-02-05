package gocb

import (
	"fmt"
	gocbcore "github.com/couchbase/gocbcore/v9"
	"github.com/couchbase/gocbcore/v9/memd"
	"time"
)

const (
	goCbVersionStr = "v2.2.0"

	durabilityTimeoutFloor = 1500 * time.Millisecond
)

// QueryIndexType provides information on the type of indexer used for an index.
type QueryIndexType string

const (
	// QueryIndexTypeGsi indicates that GSI was used to build the index.
	QueryIndexTypeGsi QueryIndexType = "gsi"

	// QueryIndexTypeView indicates that views were used to build the index.
	QueryIndexTypeView QueryIndexType = "views"
)

// QueryStatus provides information about the current status of a query.
type QueryStatus string

const (
	// QueryStatusRunning indicates the query is still running
	QueryStatusRunning QueryStatus = "running"

	// QueryStatusSuccess indicates the query was successful.
	QueryStatusSuccess QueryStatus = "success"

	// QueryStatusErrors indicates a query completed with errors.
	QueryStatusErrors QueryStatus = "errors"

	// QueryStatusCompleted indicates a query has completed.
	QueryStatusCompleted QueryStatus = "completed"

	// QueryStatusStopped indicates a query has been stopped.
	QueryStatusStopped QueryStatus = "stopped"

	// QueryStatusTimeout indicates a query timed out.
	QueryStatusTimeout QueryStatus = "timeout"

	// QueryStatusClosed indicates that a query was closed.
	QueryStatusClosed QueryStatus = "closed"

	// QueryStatusFatal indicates that a query ended with a fatal error.
	QueryStatusFatal QueryStatus = "fatal"

	// QueryStatusAborted indicates that a query was aborted.
	QueryStatusAborted QueryStatus = "aborted"

	// QueryStatusUnknown indicates that the query status is unknown.
	QueryStatusUnknown QueryStatus = "unknown"
)

// ServiceType specifies a particular Couchbase service type.
type ServiceType gocbcore.ServiceType

const (
	// ServiceTypeManagement represents a management service.
	ServiceTypeManagement ServiceType = ServiceType(gocbcore.MgmtService)

	// ServiceTypeKeyValue represents a memcached service.
	ServiceTypeKeyValue ServiceType = ServiceType(gocbcore.MemdService)

	// ServiceTypeViews represents a views service.
	ServiceTypeViews ServiceType = ServiceType(gocbcore.CapiService)

	// ServiceTypeQuery represents a query service.
	ServiceTypeQuery ServiceType = ServiceType(gocbcore.N1qlService)

	// ServiceTypeSearch represents a full-text-search service.
	ServiceTypeSearch ServiceType = ServiceType(gocbcore.FtsService)

	// ServiceTypeAnalytics represents an analytics service.
	ServiceTypeAnalytics ServiceType = ServiceType(gocbcore.CbasService)
)

// QueryProfileMode specifies the profiling mode to use during a query.
type QueryProfileMode string

const (
	// QueryProfileModeNone disables query profiling
	QueryProfileModeNone QueryProfileMode = "off"

	// QueryProfileModePhases includes phase profiling information in the query response
	QueryProfileModePhases QueryProfileMode = "phases"

	// QueryProfileModeTimings includes timing profiling information in the query response
	QueryProfileModeTimings QueryProfileMode = "timings"
)

// SubdocFlag provides special handling flags for sub-document operations
type SubdocFlag memd.SubdocFlag

const (
	// SubdocFlagNone indicates no special behaviours
	SubdocFlagNone SubdocFlag = SubdocFlag(memd.SubdocFlagNone)

	// SubdocFlagCreatePath indicates you wish to recursively create the tree of paths
	// if it does not already exist within the document.
	SubdocFlagCreatePath SubdocFlag = SubdocFlag(memd.SubdocFlagMkDirP)

	// SubdocFlagXattr indicates your path refers to an extended attribute rather than the document.
	SubdocFlagXattr SubdocFlag = SubdocFlag(memd.SubdocFlagXattrPath)

	// SubdocFlagUseMacros indicates that you wish macro substitution to occur on the value
	SubdocFlagUseMacros SubdocFlag = SubdocFlag(memd.SubdocFlagExpandMacros)
)

// SubdocDocFlag specifies document-level flags for a sub-document operation.
type SubdocDocFlag memd.SubdocDocFlag

const (
	// SubdocDocFlagNone indicates no special behaviours
	SubdocDocFlagNone SubdocDocFlag = SubdocDocFlag(memd.SubdocDocFlagNone)

	// SubdocDocFlagMkDoc indicates that the document should be created if it does not already exist.
	SubdocDocFlagMkDoc SubdocDocFlag = SubdocDocFlag(memd.SubdocDocFlagMkDoc)

	// SubdocDocFlagAddDoc indices that the document should be created only if it does not already exist.
	SubdocDocFlagAddDoc SubdocDocFlag = SubdocDocFlag(memd.SubdocDocFlagAddDoc)

	// SubdocDocFlagAccessDeleted indicates that you wish to receive soft-deleted documents.
	SubdocDocFlagAccessDeleted SubdocDocFlag = SubdocDocFlag(memd.SubdocDocFlagAccessDeleted)
)

// DurabilityLevel specifies the level of synchronous replication to use.
type DurabilityLevel uint8

const (
	// DurabilityLevelNone specifies that no durability level should be applied.
	DurabilityLevelNone DurabilityLevel = iota

	// DurabilityLevelMajority specifies that a mutation must be replicated (held in memory) to a majority of nodes.
	DurabilityLevelMajority

	// DurabilityLevelMajorityAndPersistOnMaster specifies that a mutation must be replicated (held in memory) to a
	// majority of nodes and also persisted (written to disk) on the active node.
	DurabilityLevelMajorityAndPersistOnMaster

	// DurabilityLevelPersistToMajority specifies that a mutation must be persisted (written to disk) to a majority
	// of nodes.
	DurabilityLevelPersistToMajority
)

func (dl DurabilityLevel) toManagementAPI() (string, error) {
	switch dl {
	case DurabilityLevelNone:
		return "none", nil
	case DurabilityLevelMajority:
		return "majority", nil
	case DurabilityLevelMajorityAndPersistOnMaster:
		return "majorityAndPersistActive", nil
	case DurabilityLevelPersistToMajority:
		return "persistToMajority", nil
	default:
		return "", invalidArgumentsError{
			message: fmt.Sprintf("unknown durability level: %d", dl),
		}
	}
}

func durabilityLevelFromManagementAPI(level string) DurabilityLevel {
	switch level {
	case "majority":
		return DurabilityLevelMajority
	case "majorityAndPersistActive":
		return DurabilityLevelMajorityAndPersistOnMaster
	case "persistToMajority":
		return DurabilityLevelPersistToMajority
	default:
		return DurabilityLevelNone
	}
}

// MutationMacro can be supplied to MutateIn operations to perform ExpandMacros operations.
type MutationMacro string

const (
	// MutationMacroCAS can be used to tell the server to use the CAS macro.
	MutationMacroCAS MutationMacro = "\"${Mutation.CAS}\""

	// MutationMacroSeqNo can be used to tell the server to use the seqno macro.
	MutationMacroSeqNo MutationMacro = "\"${Mutation.seqno}\""

	// MutationMacroValueCRC32c can be used to tell the server to use the value_crc32c macro.
	MutationMacroValueCRC32c MutationMacro = "\"${Mutation.value_crc32c}\""
)

// ClusterState specifies the current state of the cluster
type ClusterState uint

const (
	// ClusterStateOnline indicates that all nodes are online and reachable.
	ClusterStateOnline ClusterState = iota + 1

	// ClusterStateDegraded indicates that all services will function, but possibly not optimally.
	ClusterStateDegraded

	// ClusterStateOffline indicates that no nodes were reachable.
	ClusterStateOffline
)

// EndpointState specifies the current state of an endpoint.
type EndpointState uint

const (
	// EndpointStateDisconnected indicates the endpoint socket is unreachable.
	EndpointStateDisconnected EndpointState = iota + 1

	// EndpointStateConnecting indicates the endpoint socket is connecting.
	EndpointStateConnecting

	// EndpointStateConnected indicates the endpoint socket is connected and ready.
	EndpointStateConnected

	// EndpointStateDisconnecting indicates the endpoint socket is disconnecting.
	EndpointStateDisconnecting
)

// PingState specifies the result of the ping operation
type PingState uint

const (
	// PingStateOk indicates that the ping operation was successful.
	PingStateOk PingState = iota + 1

	// PingStateTimeout indicates that the ping operation timed out.
	PingStateTimeout

	// PingStateError indicates that the ping operation failed.
	PingStateError
)

// SaslMechanism represents a type of auth that can be performed.
type SaslMechanism string

const (
	// PlainSaslMechanism represents that PLAIN auth should be performed.
	PlainSaslMechanism SaslMechanism = SaslMechanism(gocbcore.PlainAuthMechanism)

	// ScramSha1SaslMechanism represents that SCRAM SHA1 auth should be performed.
	ScramSha1SaslMechanism SaslMechanism = SaslMechanism(gocbcore.ScramSha1AuthMechanism)

	// ScramSha256SaslMechanism represents that SCRAM SHA256 auth should be performed.
	ScramSha256SaslMechanism SaslMechanism = SaslMechanism(gocbcore.ScramSha256AuthMechanism)

	// ScramSha512SaslMechanism represents that SCRAM SHA512 auth should be performed.
	ScramSha512SaslMechanism SaslMechanism = SaslMechanism(gocbcore.ScramSha512AuthMechanism)
)
