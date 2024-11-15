package stmtcache

import (
	"container/list"
	"context"
	"fmt"
	"sync/atomic"

	"github.com/jackc/pgconn"
)

var lruCount uint64

// LRU implements Cache with a Least Recently Used (LRU) cache.
type LRU struct {
	conn         *pgconn.PgConn
	mode         int
	cap          int
	prepareCount int
	m            map[string]*list.Element
	l            *list.List
	psNamePrefix string
	stmtsToClear []string
}

// NewLRU creates a new LRU. mode is either ModePrepare or ModeDescribe. cap is the maximum size of the cache.
func NewLRU(conn *pgconn.PgConn, mode int, cap int) *LRU {
	mustBeValidMode(mode)
	mustBeValidCap(cap)

	n := atomic.AddUint64(&lruCount, 1)

	return &LRU{
		conn:         conn,
		mode:         mode,
		cap:          cap,
		m:            make(map[string]*list.Element),
		l:            list.New(),
		psNamePrefix: fmt.Sprintf("lrupsc_%d", n),
	}
}

// Get returns the prepared statement description for sql preparing or describing the sql on the server as needed.
func (c *LRU) Get(ctx context.Context, sql string) (*pgconn.StatementDescription, error) {
	if ctx != context.Background() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	// flush an outstanding bad statements
	txStatus := c.conn.TxStatus()
	if (txStatus == 'I' || txStatus == 'T') && len(c.stmtsToClear) > 0 {
		for _, stmt := range c.stmtsToClear {
			err := c.clearStmt(ctx, stmt)
			if err != nil {
				return nil, err
			}
		}
	}

	if el, ok := c.m[sql]; ok {
		c.l.MoveToFront(el)
		return el.Value.(*pgconn.StatementDescription), nil
	}

	if c.l.Len() == c.cap {
		err := c.removeOldest(ctx)
		if err != nil {
			return nil, err
		}
	}

	psd, err := c.prepare(ctx, sql)
	if err != nil {
		return nil, err
	}

	el := c.l.PushFront(psd)
	c.m[sql] = el

	return psd, nil
}

// Clear removes all entries in the cache. Any prepared statements will be deallocated from the PostgreSQL session.
func (c *LRU) Clear(ctx context.Context) error {
	for c.l.Len() > 0 {
		err := c.removeOldest(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *LRU) StatementErrored(sql string, err error) {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return
	}

	// https://github.com/jackc/pgx/issues/1162
	//
	// We used to look for the message "cached plan must not change result type". However, that message can be localized.
	// Unfortunately, error code "0A000" - "FEATURE NOT SUPPORTED" is used for many different errors and the only way to
	// tell the difference is by the message. But all that happens is we clear a statement that we otherwise wouldn't
	// have so it should be safe.
	possibleInvalidCachedPlanError := pgErr.Code == "0A000"
	if possibleInvalidCachedPlanError {
		c.stmtsToClear = append(c.stmtsToClear, sql)
	}
}

func (c *LRU) clearStmt(ctx context.Context, sql string) error {
	elem, inMap := c.m[sql]
	if !inMap {
		// The statement probably fell off the back of the list. In that case, we've
		// ensured that it isn't in the cache, so we can declare victory.
		return nil
	}

	c.l.Remove(elem)

	psd := elem.Value.(*pgconn.StatementDescription)
	delete(c.m, psd.SQL)
	if c.mode == ModePrepare {
		return c.conn.Exec(ctx, fmt.Sprintf("deallocate %s", psd.Name)).Close()
	}
	return nil
}

// Len returns the number of cached prepared statement descriptions.
func (c *LRU) Len() int {
	return c.l.Len()
}

// Cap returns the maximum number of cached prepared statement descriptions.
func (c *LRU) Cap() int {
	return c.cap
}

// Mode returns the mode of the cache (ModePrepare or ModeDescribe)
func (c *LRU) Mode() int {
	return c.mode
}

func (c *LRU) prepare(ctx context.Context, sql string) (*pgconn.StatementDescription, error) {
	var name string
	if c.mode == ModePrepare {
		name = fmt.Sprintf("%s_%d", c.psNamePrefix, c.prepareCount)
		c.prepareCount += 1
	}

	return c.conn.Prepare(ctx, name, sql, nil)
}

func (c *LRU) removeOldest(ctx context.Context) error {
	oldest := c.l.Back()
	c.l.Remove(oldest)
	psd := oldest.Value.(*pgconn.StatementDescription)
	delete(c.m, psd.SQL)
	if c.mode == ModePrepare {
		return c.conn.Exec(ctx, fmt.Sprintf("deallocate %s", psd.Name)).Close()
	}
	return nil
}
