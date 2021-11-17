package driver

// LegacyOperationKind indicates if an operation is a legacy find, getMore, or killCursors. This is used
// in Operation.Execute, which will create legacy OP_QUERY, OP_GET_MORE, or OP_KILL_CURSORS instead
// of sending them as a command.
type LegacyOperationKind uint

// These constants represent the three different kinds of legacy operations.
const (
	LegacyNone LegacyOperationKind = iota
	LegacyFind
	LegacyGetMore
	LegacyKillCursors
	LegacyListCollections
	LegacyListIndexes
)
