package gocbcore

// Cas represents a unique revision of a document.  This can be used
// to perform optimistic locking.
type Cas uint64

// VbUUID represents a unique identifier for a particular vbucket history.
type VbUUID uint64

// SeqNo is a sequential mutation number indicating the order and precise
// position of a write that has occurred.
type SeqNo uint64

// MutationToken represents a particular mutation within the cluster.
type MutationToken struct {
	VbID   uint16
	VbUUID VbUUID
	SeqNo  SeqNo
}
