package gocbcore

const (
	goCbCoreVersionStr = "v9.0.1"
)

type bucketType int

const (
	bktTypeNone                 = -1
	bktTypeInvalid   bucketType = 0
	bktTypeCouchbase            = iota
	bktTypeMemcached            = iota
)

// ServiceType specifies a particular Couchbase service type.
type ServiceType int

const (
	// MemdService represents a memcached service.
	MemdService = ServiceType(1)

	// MgmtService represents a management service (typically ns_server).
	MgmtService = ServiceType(2)

	// CapiService represents a CouchAPI service (typically for views).
	CapiService = ServiceType(3)

	// N1qlService represents a N1QL service (typically for query).
	N1qlService = ServiceType(4)

	// FtsService represents a full-text-search service.
	FtsService = ServiceType(5)

	// CbasService represents an analytics service.
	CbasService = ServiceType(6)
)

// DcpAgentPriority specifies the priority level for a dcp stream
type DcpAgentPriority uint8

const (
	// DcpAgentPriorityLow sets the priority for the dcp stream to low
	DcpAgentPriorityLow = DcpAgentPriority(0)

	// DcpAgentPriorityMed sets the priority for the dcp stream to medium
	DcpAgentPriorityMed = DcpAgentPriority(1)

	// DcpAgentPriorityHigh sets the priority for the dcp stream to high
	DcpAgentPriorityHigh = DcpAgentPriority(2)
)

type durabilityLevelStatus uint32

const (
	durabilityLevelStatusUnknown     = durabilityLevelStatus(0x00)
	durabilityLevelStatusSupported   = durabilityLevelStatus(0x01)
	durabilityLevelStatusUnsupported = durabilityLevelStatus(0x02)
)

// ClusterCapability represents a capability that the cluster supports
type ClusterCapability uint32

const (
	// ClusterCapabilityEnhancedPreparedStatements represents that the cluster supports enhanced prepared statements.
	ClusterCapabilityEnhancedPreparedStatements = ClusterCapability(0x01)
)
