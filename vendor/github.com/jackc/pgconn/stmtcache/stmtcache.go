// Package stmtcache is a cache that can be used to implement lazy prepared statements.
package stmtcache

import (
	"context"

	"github.com/jackc/pgconn"
)

const (
	ModePrepare  = iota // Cache should prepare named statements.
	ModeDescribe        // Cache should prepare the anonymous prepared statement to only fetch the description of the statement.
)

// Cache prepares and caches prepared statement descriptions.
type Cache interface {
	// Get returns the prepared statement description for sql preparing or describing the sql on the server as needed.
	Get(ctx context.Context, sql string) (*pgconn.StatementDescription, error)

	// Clear removes all entries in the cache. Any prepared statements will be deallocated from the PostgreSQL session.
	Clear(ctx context.Context) error

	// StatementErrored informs the cache that the given statement resulted in an error when it
	// was last used against the database. In some cases, this will cause the cache to maer that
	// statement as bad. The bad statement will instead be flushed during the next call to Get
	// that occurs outside of a failed transaction.
	StatementErrored(sql string, err error)

	// Len returns the number of cached prepared statement descriptions.
	Len() int

	// Cap returns the maximum number of cached prepared statement descriptions.
	Cap() int

	// Mode returns the mode of the cache (ModePrepare or ModeDescribe)
	Mode() int
}

// New returns the preferred cache implementation for mode and cap. mode is either ModePrepare or ModeDescribe. cap is
// the maximum size of the cache.
func New(conn *pgconn.PgConn, mode int, cap int) Cache {
	mustBeValidMode(mode)
	mustBeValidCap(cap)

	return NewLRU(conn, mode, cap)
}

func mustBeValidMode(mode int) {
	if mode != ModePrepare && mode != ModeDescribe {
		panic("mode must be ModePrepare or ModeDescribe")
	}
}

func mustBeValidCap(cap int) {
	if cap < 1 {
		panic("cache must have cap of >= 1")
	}
}
