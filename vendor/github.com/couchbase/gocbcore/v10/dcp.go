package gocbcore

// OpenStreamFilterOptions are the filtering options available to the OpenStream operation.
type OpenStreamFilterOptions struct {
	ScopeID       uint32
	CollectionIDs []uint32
}

// OpenStreamStreamOptions are the stream options available to the OpenStream operation.
type OpenStreamStreamOptions struct {
	StreamID uint16
}

// OpenStreamManifestOptions are the manifest options available to the OpenStream operation.
type OpenStreamManifestOptions struct {
	ManifestUID uint64
}

// OpenStreamOptions are the options available to the OpenStream operation.
type OpenStreamOptions struct {
	FilterOptions   *OpenStreamFilterOptions
	StreamOptions   *OpenStreamStreamOptions
	ManifestOptions *OpenStreamManifestOptions
}

// GetVbucketSeqnoFilterOptions are the filter options available to the GetVbucketSeqno operation.
type GetVbucketSeqnoFilterOptions struct {
	CollectionID uint32
}

// GetVbucketSeqnoOptions are the options available to the GetVbucketSeqno operation.
type GetVbucketSeqnoOptions struct {
	FilterOptions *GetVbucketSeqnoFilterOptions
}

// CloseStreamStreamOptions are the stream options available to the CloseStream operation.
type CloseStreamStreamOptions struct {
	StreamID uint16
}

// CloseStreamOptions are the options available to the CloseStream operation.
type CloseStreamOptions struct {
	StreamOptions *CloseStreamStreamOptions
}

// SnapshotState represents the state of a particular cluster snapshot.
type SnapshotState uint32

// HasInMemory returns whether this snapshot is available in memory.
func (s SnapshotState) HasInMemory() bool {
	return uint32(s)&1 != 0
}

// HasOnDisk returns whether this snapshot is available on disk.
func (s SnapshotState) HasOnDisk() bool {
	return uint32(s)&2 != 0
}

// HasHistory returns whether this snapshot represents a view of history.
func (s SnapshotState) HasHistory() bool {
	return uint32(s)&16 != 0
}

// HasMayDuplicateKeys returns whether this snapshot may contain duplicate keys.
func (s SnapshotState) HasMayDuplicateKeys() bool {
	return uint32(s)&32 != 0
}

// FailoverEntry represents a single entry in the server fail-over log.
type FailoverEntry struct {
	VbUUID VbUUID
	SeqNo  SeqNo
}

// DcpSnapshotMarker represents a single response from the server
type DcpSnapshotMarker struct {
	StartSeqNo, EndSeqNo                                   uint64
	VbID, StreamID                                         uint16
	SnapshotType                                           SnapshotState
	MaxVisibleSeqNo, HighCompletedSeqNo, SnapshotTimeStamp uint64
}

// DcpMutation represents a single DCP mutation from the server
type DcpMutation struct {
	SeqNo, RevNo            uint64
	Cas                     uint64
	Flags, Expiry, LockTime uint32
	CollectionID            uint32
	VbID                    uint16
	StreamID                uint16
	Datatype                uint8
	Key, Value              []byte
}

// DcpDeletion represents a single DCP deletion from the server
type DcpDeletion struct {
	SeqNo, RevNo uint64
	Cas          uint64
	DeleteTime   uint32
	CollectionID uint32
	VbID         uint16
	StreamID     uint16
	Datatype     uint8
	Key, Value   []byte
}

// DcpExpiration represents a single DCP expiration from the server
type DcpExpiration struct {
	SeqNo, RevNo uint64
	Cas          uint64
	DeleteTime   uint32
	CollectionID uint32
	VbID         uint16
	StreamID     uint16
	Key          []byte
}

// DcpCollectionCreation represents a collection create DCP event from the server
type DcpCollectionCreation struct {
	SeqNo        uint64
	Version      uint8
	VbID         uint16
	ManifestUID  uint64
	ScopeID      uint32
	CollectionID uint32
	Ttl          uint32
	StreamID     uint16
	Key          []byte
}

// DcpCollectionDeleteion represents a collection delete DCP event from the server
type DcpCollectionDeletion struct {
	SeqNo        uint64
	ManifestUID  uint64
	ScopeID      uint32
	CollectionID uint32
	StreamID     uint16
	VbID         uint16
	Version      uint8
}

// DcpCollectionFlush represents a collection flush DCP event from the server
type DcpCollectionFlush struct {
	SeqNo        uint64
	Version      uint8
	VbID         uint16
	ManifestUID  uint64
	CollectionID uint32
	StreamID     uint16
}

// DcpScopeCreation represents a scope creation DCP event from the server
type DcpScopeCreation struct {
	SeqNo       uint64
	Version     uint8
	VbID        uint16
	ManifestUID uint64
	ScopeID     uint32
	StreamID    uint16
	Key         []byte
}

// DcpScopeDeletion represents a scope Deletion DCP event from the server
type DcpScopeDeletion struct {
	SeqNo       uint64
	Version     uint8
	VbID        uint16
	ManifestUID uint64
	ScopeID     uint32
	StreamID    uint16
}

// DcpCollectionModification represents a DCP collection modify event from the server
type DcpCollectionModification struct {
	SeqNo        uint64
	ManifestUID  uint64
	CollectionID uint32
	Ttl          uint32
	VbID         uint16
	StreamID     uint16
	Version      uint8
}

// DcpOSOSnapshot reprensents a DCP OSSSnapshot from the server
type DcpOSOSnapshot struct {
	SnapshotType uint32
	VbID         uint16
	StreamID     uint16
}

// DcpSeqNoAdvanced represents a DCP SeqNoAdvanced from the server
type DcpSeqNoAdvanced struct {
	SeqNo    uint64
	VbID     uint16
	StreamID uint16
}

// DcpStreamEnd represnets a DCP stream end from the server
type DcpStreamEnd struct {
	VbID     uint16
	StreamID uint16
}

// StreamObserver provides an interface to receive events from a running DCP stream.
type StreamObserver interface {
	SnapshotMarker(snapshotMarker DcpSnapshotMarker)
	Mutation(mutation DcpMutation)
	Deletion(deletion DcpDeletion)
	Expiration(expiration DcpExpiration)
	End(end DcpStreamEnd, err error)
	CreateCollection(creation DcpCollectionCreation)
	DeleteCollection(deletion DcpCollectionDeletion)
	FlushCollection(flush DcpCollectionFlush)
	CreateScope(creation DcpScopeCreation)
	DeleteScope(deletion DcpScopeDeletion)
	ModifyCollection(modification DcpCollectionModification)
	OSOSnapshot(snapshot DcpOSOSnapshot)
	SeqNoAdvanced(seqNoAdvanced DcpSeqNoAdvanced)
}

type streamFilter struct {
	ManifestUID string   `json:"uid,omitempty"`
	Collections []string `json:"collections,omitempty"`
	Scope       string   `json:"scope,omitempty"`
	StreamID    uint16   `json:"sid,omitempty"`
}

// OpenStreamCallback is invoked with the results of `OpenStream` operations.
type OpenStreamCallback func([]FailoverEntry, error)

// CloseStreamCallback is invoked with the results of `CloseStream` operations.
type CloseStreamCallback func(error)

// GetFailoverLogCallback is invoked with the results of `GetFailoverLog` operations.
type GetFailoverLogCallback func([]FailoverEntry, error)

// VbSeqNoEntry represents a single GetVbucketSeqnos sequence number entry.
type VbSeqNoEntry struct {
	VbID  uint16
	SeqNo SeqNo
}

// GetVBucketSeqnosCallback is invoked with the results of `GetVBucketSeqnos` operations.
type GetVBucketSeqnosCallback func([]VbSeqNoEntry, error)
