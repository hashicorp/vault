package canoe

// SnapshotData defines what a snapshot should look like
type SnapshotData []byte

// FSM is an interface for what your state machine needs to define in order to be compatible with canoe
type FSM interface {
	// Be sparing with errors for all the following.
	// Err only if it results in bad state.
	// Because it will halt all the things

	// Apply is called whenever a new log is committed to raft.
	// The FSM is responsible for applying it in an atomic fashion
	Apply(entry LogData) error

	// Snapshot should return a snapshot in the form of restorable info for the entire FSM
	Snapshot() (SnapshotData, error)

	// Restore should take a snapshot, and use it to restore the state of the FSM
	Restore(snap SnapshotData) error
}
