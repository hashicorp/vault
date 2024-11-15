package gocbcore

const (
	goCbCoreVersionStr = "v10.5.2"
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

	// EventingService represents the eventing service.
	EventingService = ServiceType(7)

	// GSIService represents the indexing service.
	GSIService = ServiceType(8)

	// BackupService represents the backup service.
	BackupService = ServiceType(9)
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

type BucketCapability uint32

const (
	BucketCapabilityDurableWrites        BucketCapability = 0x00
	BucketCapabilityCreateAsDeleted      BucketCapability = 0x01
	BucketCapabilityReplaceBodyWithXattr BucketCapability = 0x02
	BucketCapabilityRangeScan            BucketCapability = 0x03
	BucketCapabilityReplicaRead          BucketCapability = 0x04
	BucketCapabilityNonDedupedHistory    BucketCapability = 0x05
)

type CapabilityStatus uint32

const (
	CapabilityStatusUnknown CapabilityStatus = iota
	CapabilityStatusSupported
	CapabilityStatusUnsupported
)

// Deprecated: Use CapabilityStatus instead
type BucketCapabilityStatus = CapabilityStatus

const (
	// Deprecated: Use CapabilityStatusUnknown instead
	BucketCapabilityStatusUnknown = CapabilityStatusUnknown

	// Deprecated: Use CapabilityStatusSupported instead
	BucketCapabilityStatusSupported = CapabilityStatusSupported

	// Deprecated: Use CapabilityStatusUnsupported instead
	BucketCapabilityStatusUnsupported = CapabilityStatusUnsupported
)

// DCPBackfillOrder represents the order in which vBuckets will be backfilled by the cluster.
type DCPBackfillOrder uint8

const (
	// DCPBackfillOrderRoundRobin means that all the requested vBuckets will be backfilled together where each vBucket
	// has some data backfilled before moving on to the next. This is the default behaviour.
	DCPBackfillOrderRoundRobin DCPBackfillOrder = iota + 1

	// DCPBackfillOrderSequential means that all the data for the first vBucket will be streamed before advancing onto
	// the next vBucket.
	DCPBackfillOrderSequential
)

const (
	spanNameDispatchToServer    = "dispatch_to_server"
	spanAttribDBSystemKey       = "db.system"
	spanAttribDBSystemValue     = "couchbase"
	spanAttribNetTransportKey   = "net.transport"
	spanAttribNetTransportValue = "IP.TCP"
	spanAttribOperationIDKey    = "db.couchbase.operation_id"
	spanAttribLocalIDKey        = "db.couchbase.local_id"
	spanAttribNetHostNameKey    = "net.host.name"
	spanAttribNetHostPortKey    = "net.host.port"
	spanAttribNetPeerNameKey    = "net.peer.name"
	spanAttribNetPeerPortKey    = "net.peer.port"
	spanAttribServerDurationKey = "db.couchbase.server_duration"
	spanAttribNumRetries        = "db.couchbase.retries"
)

const (
	metricAttribServiceKey           = "db.couchbase.service"
	metricAttribOperationKey         = "db.operation"
	meterNameCBOperations            = "db.couchbase.operations"
	metricValueServiceKeyValue       = "kv"
	metricValueServiceQueryValue     = "n1ql"
	metricValueServiceSearchValue    = "fts"
	metricValueServiceAnalyticsValue = "cbas"
	metricValueServiceViewsValue     = "capi"
	metricValueServiceHTTPValue      = "http"
)

type SpanStatus string

const (
	SpanStatusOK    SpanStatus = "Ok"
	SpanStatusError SpanStatus = "Error"
)

type statusClass uint8

const (
	statusClassOK statusClass = iota
	statusClassError
)

var crc32cMacro = []byte("\"${Mutation.value_crc32c}\"")
var revidMacro = []byte("\"${$document.revid}\"")
var exptimeMacro = []byte("\"${$document.exptime}\"")
var casMacro = []byte("\"${$document.CAS}\"")
var hlcMacro = "$vbucket.HLC"
