package raft

// LogType describes various types of log entries.
type LogType uint8

const (
	// LogCommand is applied to a user FSM.
	LogCommand LogType = iota

	// LogNoop is used to assert leadership.
	LogNoop

	// LogAddPeer is used to add a new peer. This should only be used with
	// older protocol versions designed to be compatible with unversioned
	// Raft servers. See comments in config.go for details.
	LogAddPeerDeprecated

	// LogRemovePeer is used to remove an existing peer. This should only be
	// used with older protocol versions designed to be compatible with
	// unversioned Raft servers. See comments in config.go for details.
	LogRemovePeerDeprecated

	// LogBarrier is used to ensure all preceding operations have been
	// applied to the FSM. It is similar to LogNoop, but instead of returning
	// once committed, it only returns once the FSM manager acks it. Otherwise
	// it is possible there are operations committed but not yet applied to
	// the FSM.
	LogBarrier

	// LogConfiguration establishes a membership change configuration. It is
	// created when a server is added, removed, promoted, etc. Only used
	// when protocol version 1 or greater is in use.
	LogConfiguration
)

// Log entries are replicated to all members of the Raft cluster
// and form the heart of the replicated state machine.
type Log struct {
	// Index holds the index of the log entry.
	Index uint64

	// Term holds the election term of the log entry.
	Term uint64

	// Type holds the type of the log entry.
	Type LogType

	// Data holds the log entry's type-specific data.
	Data []byte

	// Extensions holds an opaque byte slice of information for middleware. It
	// is up to the client of the library to properly modify this as it adds
	// layers and remove those layers when appropriate. This value is a part of
	// the log, so very large values could cause timing issues.
	//
	// N.B. It is _up to the client_ to handle upgrade paths. For instance if
	// using this with go-raftchunking, the client should ensure that all Raft
	// peers are using a version that can handle that extension before ever
	// actually triggering chunking behavior. It is sometimes sufficient to
	// ensure that non-leaders are upgraded first, then the current leader is
	// upgraded, but a leader changeover during this process could lead to
	// trouble, so gating extension behavior via some flag in the client
	// program is also a good idea.
	Extensions []byte
}

// LogStore is used to provide an interface for storing
// and retrieving logs in a durable fashion.
type LogStore interface {
	// FirstIndex returns the first index written. 0 for no entries.
	FirstIndex() (uint64, error)

	// LastIndex returns the last index written. 0 for no entries.
	LastIndex() (uint64, error)

	// GetLog gets a log entry at a given index.
	GetLog(index uint64, log *Log) error

	// StoreLog stores a log entry.
	StoreLog(log *Log) error

	// StoreLogs stores multiple log entries.
	StoreLogs(logs []*Log) error

	// DeleteRange deletes a range of log entries. The range is inclusive.
	DeleteRange(min, max uint64) error
}
