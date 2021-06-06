package memd

import "fmt"

// CmdMagic represents the magic number that begins the header
// of every packet and informs the rest of the header format.
type CmdMagic uint8

const (
	// CmdMagicReq indicates that the packet is a request.
	CmdMagicReq = CmdMagic(0x80)

	// CmdMagicRes indicates that the packet is a response.
	CmdMagicRes = CmdMagic(0x81)

	// These are private rather than public as the library will automatically
	// switch to and from these magics based on the use of frames within a packet.
	cmdMagicReqExt = CmdMagic(0x08)
	cmdMagicResExt = CmdMagic(0x18)
)

// frameType specifies which kind of frame extra a particular block belongs to.
// This is a private type since we automatically encode this internally based on
// whether the specific frame block is attached to the packet.
type frameType uint8

const (
	frameTypeReqBarrier        = frameType(0)
	frameTypeReqSyncDurability = frameType(1)
	frameTypeReqStreamID       = frameType(2)
	frameTypeReqOpenTracing    = frameType(3)
	frameTypeResSrvDuration    = frameType(0)
)

// HelloFeature represents a feature code included in a memcached
// HELLO operation.
type HelloFeature uint16

const (
	// FeatureDatatype indicates support for Datatype fields.
	FeatureDatatype = HelloFeature(0x01)

	// FeatureTLS indicates support for TLS
	FeatureTLS = HelloFeature(0x02)

	// FeatureTCPNoDelay indicates support for TCP no-delay.
	FeatureTCPNoDelay = HelloFeature(0x03)

	// FeatureSeqNo indicates support for mutation tokens.
	FeatureSeqNo = HelloFeature(0x04)

	// FeatureTCPDelay indicates support for TCP delay.
	FeatureTCPDelay = HelloFeature(0x05)

	// FeatureXattr indicates support for document xattrs.
	FeatureXattr = HelloFeature(0x06)

	// FeatureXerror indicates support for extended errors.
	FeatureXerror = HelloFeature(0x07)

	// FeatureSelectBucket indicates support for the SelectBucket operation.
	FeatureSelectBucket = HelloFeature(0x08)

	// Feature 0x09 is reserved and cannot be used.

	// FeatureSnappy indicates support for snappy compressed documents.
	FeatureSnappy = HelloFeature(0x0a)

	// FeatureJSON indicates support for JSON datatype data.
	FeatureJSON = HelloFeature(0x0b)

	// FeatureDuplex indicates support for duplex communications.
	FeatureDuplex = HelloFeature(0x0c)

	// FeatureClusterMapNotif indicates support for cluster-map update notifications.
	FeatureClusterMapNotif = HelloFeature(0x0d)

	// FeatureUnorderedExec indicates support for unordered execution of operations.
	FeatureUnorderedExec = HelloFeature(0x0e)

	// FeatureDurations indicates support for server durations.
	FeatureDurations = HelloFeature(0xf)

	// FeatureAltRequests indicates support for requests with flexible frame extras.
	FeatureAltRequests = HelloFeature(0x10)

	// FeatureSyncReplication indicates support for requests synchronous durability requirements.
	FeatureSyncReplication = HelloFeature(0x11)

	// FeatureCollections indicates support for collections.
	FeatureCollections = HelloFeature(0x12)

	// FeatureOpenTracing indicates support for OpenTracing.
	FeatureOpenTracing = HelloFeature(0x13)

	// FeatureCreateAsDeleted indicates support for the create as deleted feature.
	FeatureCreateAsDeleted = HelloFeature(0x17)
)

// StreamEndStatus represents the reason for a DCP stream ending
type StreamEndStatus uint32

const (
	// StreamEndOK represents that the stream ended successfully.
	StreamEndOK = StreamEndStatus(0x00)

	// StreamEndClosed represents that the stream was forcefully closed.
	StreamEndClosed = StreamEndStatus(0x01)

	// StreamEndStateChanged represents that the stream was closed due to a state change.
	StreamEndStateChanged = StreamEndStatus(0x02)

	// StreamEndDisconnected represents that the stream was closed due to disconnection.
	StreamEndDisconnected = StreamEndStatus(0x03)

	// StreamEndTooSlow represents that the stream was closed due to the stream being too slow.
	StreamEndTooSlow = StreamEndStatus(0x04)

	// StreamEndBackfillFailed represents that the stream was closed due to backfill failing.
	StreamEndBackfillFailed = StreamEndStatus(0x05)

	// StreamEndFilterEmpty represents that the stream was closed due to the filter being empty.
	StreamEndFilterEmpty = StreamEndStatus(0x07)
)

// KVText returns the textual representation of this StreamEndStatus.
func (code StreamEndStatus) KVText() string {
	switch code {
	case StreamEndOK:
		return "success"
	case StreamEndClosed:
		return "stream closed"
	case StreamEndStateChanged:
		return "state changed"
	case StreamEndDisconnected:
		return "disconnected"
	case StreamEndTooSlow:
		return "too slow"
	case StreamEndFilterEmpty:
		return "filter empty"
	case StreamEndBackfillFailed:
		return "backfill failed"
	default:
		return fmt.Sprintf("unknown stream close reason (%d)", code)
	}
}

// StreamEventCode is the code for a DCP Stream event
type StreamEventCode uint32

const (
	// StreamEventCollectionCreate is the StreamEventCode for a collection create event
	StreamEventCollectionCreate = StreamEventCode(0x00)

	// StreamEventCollectionDelete is the StreamEventCode for a collection delete event
	StreamEventCollectionDelete = StreamEventCode(0x01)

	// StreamEventCollectionFlush is the StreamEventCode for a collection flush event
	StreamEventCollectionFlush = StreamEventCode(0x02)

	// StreamEventScopeCreate is the StreamEventCode for a scope create event
	StreamEventScopeCreate = StreamEventCode(0x03)

	// StreamEventScopeDelete is the StreamEventCode for a scope delete event
	StreamEventScopeDelete = StreamEventCode(0x04)

	// StreamEventCollectionChanged is the StreamEventCode for a collection changed event
	StreamEventCollectionChanged = StreamEventCode(0x05)
)

// VbucketState represents the state of a particular vbucket on a particular server.
type VbucketState uint32

const (
	// VbucketStateActive indicates the vbucket is active on this server
	VbucketStateActive = VbucketState(0x01)

	// VbucketStateReplica indicates the vbucket is a replica on this server
	VbucketStateReplica = VbucketState(0x02)

	// VbucketStatePending indicates the vbucket is preparing to become active on this server.
	VbucketStatePending = VbucketState(0x03)

	// VbucketStateDead indicates the vbucket is no longer valid on this server.
	VbucketStateDead = VbucketState(0x04)
)

// SetMetaOption represents possible option values for a SetMeta operation.
type SetMetaOption uint32

const (
	// ForceMetaOp disables conflict resolution for the document and allows the
	// operation to be applied to an active, pending, or replica vbucket.
	ForceMetaOp = SetMetaOption(0x01)

	// UseLwwConflictResolution switches to Last-Write-Wins conflict resolution
	// for the document.
	UseLwwConflictResolution = SetMetaOption(0x02)

	// RegenerateCas causes the server to invalidate the current CAS value for
	// a document, and to generate a new one.
	RegenerateCas = SetMetaOption(0x04)

	// SkipConflictResolution disables conflict resolution for the document.
	SkipConflictResolution = SetMetaOption(0x08)

	// IsExpiration indicates that the message is for an expired document.
	IsExpiration = SetMetaOption(0x10)
)

// KeyState represents the various storage states of a key on the server.
type KeyState uint8

const (
	// KeyStateNotPersisted indicates the key is in memory, but not yet written to disk.
	KeyStateNotPersisted = KeyState(0x00)

	// KeyStatePersisted indicates that the key has been written to disk.
	KeyStatePersisted = KeyState(0x01)

	// KeyStateNotFound indicates that the key is not found in memory or on disk.
	KeyStateNotFound = KeyState(0x80)

	// KeyStateDeleted indicates that the key has been written to disk as deleted.
	KeyStateDeleted = KeyState(0x81)
)

// SubDocOpType specifies the type of a sub-document operation.
type SubDocOpType uint8

const (
	// SubDocOpGet indicates the operation is a sub-document `Get` operation.
	SubDocOpGet = SubDocOpType(CmdSubDocGet)

	// SubDocOpExists indicates the operation is a sub-document `Exists` operation.
	SubDocOpExists = SubDocOpType(CmdSubDocExists)

	// SubDocOpGetCount indicates the operation is a sub-document `GetCount` operation.
	SubDocOpGetCount = SubDocOpType(CmdSubDocGetCount)

	// SubDocOpDictAdd indicates the operation is a sub-document `Add` operation.
	SubDocOpDictAdd = SubDocOpType(CmdSubDocDictAdd)

	// SubDocOpDictSet indicates the operation is a sub-document `Set` operation.
	SubDocOpDictSet = SubDocOpType(CmdSubDocDictSet)

	// SubDocOpDelete indicates the operation is a sub-document `Remove` operation.
	SubDocOpDelete = SubDocOpType(CmdSubDocDelete)

	// SubDocOpReplace indicates the operation is a sub-document `Replace` operation.
	SubDocOpReplace = SubDocOpType(CmdSubDocReplace)

	// SubDocOpArrayPushLast indicates the operation is a sub-document `ArrayPushLast` operation.
	SubDocOpArrayPushLast = SubDocOpType(CmdSubDocArrayPushLast)

	// SubDocOpArrayPushFirst indicates the operation is a sub-document `ArrayPushFirst` operation.
	SubDocOpArrayPushFirst = SubDocOpType(CmdSubDocArrayPushFirst)

	// SubDocOpArrayInsert indicates the operation is a sub-document `ArrayInsert` operation.
	SubDocOpArrayInsert = SubDocOpType(CmdSubDocArrayInsert)

	// SubDocOpArrayAddUnique indicates the operation is a sub-document `ArrayAddUnique` operation.
	SubDocOpArrayAddUnique = SubDocOpType(CmdSubDocArrayAddUnique)

	// SubDocOpCounter indicates the operation is a sub-document `Counter` operation.
	SubDocOpCounter = SubDocOpType(CmdSubDocCounter)

	// SubDocOpGetDoc represents a full document retrieval, for use with extended attribute ops.
	SubDocOpGetDoc = SubDocOpType(CmdGet)

	// SubDocOpSetDoc represents a full document set, for use with extended attribute ops.
	SubDocOpSetDoc = SubDocOpType(CmdSet)

	// SubDocOpAddDoc represents a full document add, for use with extended attribute ops.
	SubDocOpAddDoc = SubDocOpType(CmdAdd)

	// SubDocOpDeleteDoc represents a full document delete, for use with extended attribute ops.
	SubDocOpDeleteDoc = SubDocOpType(CmdDelete)
)

// DcpOpenFlag specifies flags for DCP connections configured when the stream is opened.
type DcpOpenFlag uint32

const (
	// DcpOpenFlagProducer indicates this connection wants the other end to be a producer.
	DcpOpenFlagProducer = DcpOpenFlag(0x01)

	// DcpOpenFlagNotifier indicates this connection wants the other end to be a notifier.
	DcpOpenFlagNotifier = DcpOpenFlag(0x02)

	// DcpOpenFlagIncludeXattrs indicates the client wishes to receive extended attributes.
	DcpOpenFlagIncludeXattrs = DcpOpenFlag(0x04)

	// DcpOpenFlagNoValue indicates the client does not wish to receive mutation values.
	DcpOpenFlagNoValue = DcpOpenFlag(0x08)

	// DcpOpenFlagIncludeDeleteTimes indicates the client wishes to receive delete times.
	DcpOpenFlagIncludeDeleteTimes = DcpOpenFlag(0x20)
)

// DcpStreamAddFlag specifies flags for DCP streams configured when the stream is opened.
type DcpStreamAddFlag uint32

const (
	//DcpStreamAddFlagDiskOnly indicates that stream should only send items if they are on disk
	DcpStreamAddFlagDiskOnly = DcpStreamAddFlag(0x02)

	// DcpStreamAddFlagLatest indicates this stream wants to get data up to the latest seqno.
	DcpStreamAddFlagLatest = DcpStreamAddFlag(0x04)

	// DcpStreamAddFlagActiveOnly indicates this stream should only connect to an active vbucket.
	DcpStreamAddFlagActiveOnly = DcpStreamAddFlag(0x10)

	// DcpStreamAddFlagStrictVBUUID indicates the vbuuid must match unless the start seqno
	// is 0 and the vbuuid is also 0.
	DcpStreamAddFlagStrictVBUUID = DcpStreamAddFlag(0x20)
)

// DatatypeFlag specifies data flags for the value of a document.
type DatatypeFlag uint8

const (
	// DatatypeFlagJSON indicates the server believes the value payload to be JSON.
	DatatypeFlagJSON = DatatypeFlag(0x01)

	// DatatypeFlagCompressed indicates the value payload is compressed.
	DatatypeFlagCompressed = DatatypeFlag(0x02)

	// DatatypeFlagXattrs indicates the inclusion of xattr data in the value payload.
	DatatypeFlagXattrs = DatatypeFlag(0x04)
)

// SubdocFlag specifies flags for a sub-document operation.
type SubdocFlag uint8

const (
	// SubdocFlagNone indicates no special treatment for this operation.
	SubdocFlagNone = SubdocFlag(0x00)

	// SubdocFlagMkDirP indicates that the path should be created if it does not already exist.
	SubdocFlagMkDirP = SubdocFlag(0x01)

	// 0x02 is unused, formally SubdocFlagMkDoc

	// SubdocFlagXattrPath indicates that the path refers to an Xattr rather than the document body.
	SubdocFlagXattrPath = SubdocFlag(0x04)

	// 0x08 is unused, formally SubdocFlagAccessDeleted

	// SubdocFlagExpandMacros indicates that the value portion of any sub-document mutations
	// should be expanded if they contain macros such as ${Mutation.CAS}.
	SubdocFlagExpandMacros = SubdocFlag(0x10)
)

// SubdocDocFlag specifies document-level flags for a sub-document operation.
type SubdocDocFlag uint8

const (
	// SubdocDocFlagNone indicates no special treatment for this operation.
	SubdocDocFlagNone = SubdocDocFlag(0x00)

	// SubdocDocFlagMkDoc indicates that the document should be created if it does not already exist.
	SubdocDocFlagMkDoc = SubdocDocFlag(0x01)

	// SubdocDocFlagAddDoc indices that this operation should be an add rather than set.
	SubdocDocFlagAddDoc = SubdocDocFlag(0x02)

	// SubdocDocFlagAccessDeleted indicates that you wish to receive soft-deleted documents.
	// Internal: This should never be used and is not supported.
	SubdocDocFlagAccessDeleted = SubdocDocFlag(0x04)

	// SubdocDocFlagCreateAsDeleted indicates that the document should be created as deleted.
	// That is, to create a tombstone only.
	// Internal: This should never be used and is not supported.
	SubdocDocFlagCreateAsDeleted = SubdocDocFlag(0x08)
)

// DurabilityLevel specifies the level to use for enhanced durability requirements.
type DurabilityLevel uint8

const (
	// DurabilityLevelMajority specifies that a change must be replicated to (held in memory)
	// a majority of the nodes for the bucket.
	DurabilityLevelMajority = DurabilityLevel(0x01)

	// DurabilityLevelMajorityAndPersistOnMaster specifies that a change must be replicated to (held in memory)
	// a majority of the nodes for the bucket and additionally persisted to disk on the active node.
	DurabilityLevelMajorityAndPersistOnMaster = DurabilityLevel(0x02)

	// DurabilityLevelPersistToMajority specifies that a change must be persisted to (written to disk)
	// a majority for the bucket.
	DurabilityLevelPersistToMajority = DurabilityLevel(0x03)
)
