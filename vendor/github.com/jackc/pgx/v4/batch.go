package pgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)

type batchItem struct {
	query     string
	arguments []interface{}
}

// Batch queries are a way of bundling multiple queries together to avoid
// unnecessary network round trips.
type Batch struct {
	items []*batchItem
}

// Queue queues a query to batch b. query can be an SQL query or the name of a prepared statement.
func (b *Batch) Queue(query string, arguments ...interface{}) {
	b.items = append(b.items, &batchItem{
		query:     query,
		arguments: arguments,
	})
}

// Len returns number of queries that have been queued so far.
func (b *Batch) Len() int {
	return len(b.items)
}

type BatchResults interface {
	// Exec reads the results from the next query in the batch as if the query has been sent with Conn.Exec.
	Exec() (pgconn.CommandTag, error)

	// Query reads the results from the next query in the batch as if the query has been sent with Conn.Query.
	Query() (Rows, error)

	// QueryRow reads the results from the next query in the batch as if the query has been sent with Conn.QueryRow.
	QueryRow() Row

	// QueryFunc reads the results from the next query in the batch as if the query has been sent with Conn.QueryFunc.
	QueryFunc(scans []interface{}, f func(QueryFuncRow) error) (pgconn.CommandTag, error)

	// Close closes the batch operation. This must be called before the underlying connection can be used again. Any error
	// that occurred during a batch operation may have made it impossible to resyncronize the connection with the server.
	// In this case the underlying connection will have been closed. Close is safe to call multiple times.
	Close() error
}

type batchResults struct {
	ctx    context.Context
	conn   *Conn
	mrr    *pgconn.MultiResultReader
	err    error
	b      *Batch
	ix     int
	closed bool
}

// Exec reads the results from the next query in the batch as if the query has been sent with Exec.
func (br *batchResults) Exec() (pgconn.CommandTag, error) {
	if br.err != nil {
		return nil, br.err
	}
	if br.closed {
		return nil, fmt.Errorf("batch already closed")
	}

	query, arguments, _ := br.nextQueryAndArgs()

	if !br.mrr.NextResult() {
		err := br.mrr.Close()
		if err == nil {
			err = errors.New("no result")
		}
		if br.conn.shouldLog(LogLevelError) {
			br.conn.log(br.ctx, LogLevelError, "BatchResult.Exec", map[string]interface{}{
				"sql":  query,
				"args": logQueryArgs(arguments),
				"err":  err,
			})
		}
		return nil, err
	}

	commandTag, err := br.mrr.ResultReader().Close()

	if err != nil {
		if br.conn.shouldLog(LogLevelError) {
			br.conn.log(br.ctx, LogLevelError, "BatchResult.Exec", map[string]interface{}{
				"sql":  query,
				"args": logQueryArgs(arguments),
				"err":  err,
			})
		}
	} else if br.conn.shouldLog(LogLevelInfo) {
		br.conn.log(br.ctx, LogLevelInfo, "BatchResult.Exec", map[string]interface{}{
			"sql":        query,
			"args":       logQueryArgs(arguments),
			"commandTag": commandTag,
		})
	}

	return commandTag, err
}

// Query reads the results from the next query in the batch as if the query has been sent with Query.
func (br *batchResults) Query() (Rows, error) {
	query, arguments, ok := br.nextQueryAndArgs()
	if !ok {
		query = "batch query"
	}

	if br.err != nil {
		return &connRows{err: br.err, closed: true}, br.err
	}

	if br.closed {
		alreadyClosedErr := fmt.Errorf("batch already closed")
		return &connRows{err: alreadyClosedErr, closed: true}, alreadyClosedErr
	}

	rows := br.conn.getRows(br.ctx, query, arguments)

	if !br.mrr.NextResult() {
		rows.err = br.mrr.Close()
		if rows.err == nil {
			rows.err = errors.New("no result")
		}
		rows.closed = true

		if br.conn.shouldLog(LogLevelError) {
			br.conn.log(br.ctx, LogLevelError, "BatchResult.Query", map[string]interface{}{
				"sql":  query,
				"args": logQueryArgs(arguments),
				"err":  rows.err,
			})
		}

		return rows, rows.err
	}

	rows.resultReader = br.mrr.ResultReader()
	return rows, nil
}

// QueryFunc reads the results from the next query in the batch as if the query has been sent with Conn.QueryFunc.
func (br *batchResults) QueryFunc(scans []interface{}, f func(QueryFuncRow) error) (pgconn.CommandTag, error) {
	if br.closed {
		return nil, fmt.Errorf("batch already closed")
	}

	rows, err := br.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}

		err = f(rows)
		if err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rows.CommandTag(), nil
}

// QueryRow reads the results from the next query in the batch as if the query has been sent with QueryRow.
func (br *batchResults) QueryRow() Row {
	rows, _ := br.Query()
	return (*connRow)(rows.(*connRows))

}

// Close closes the batch operation. Any error that occurred during a batch operation may have made it impossible to
// resyncronize the connection with the server. In this case the underlying connection will have been closed.
func (br *batchResults) Close() error {
	if br.err != nil {
		return br.err
	}

	if br.closed {
		return nil
	}
	br.closed = true

	// log any queries that haven't yet been logged by Exec or Query
	for {
		query, args, ok := br.nextQueryAndArgs()
		if !ok {
			break
		}

		if br.conn.shouldLog(LogLevelInfo) {
			br.conn.log(br.ctx, LogLevelInfo, "BatchResult.Close", map[string]interface{}{
				"sql":  query,
				"args": logQueryArgs(args),
			})
		}
	}

	return br.mrr.Close()
}

func (br *batchResults) nextQueryAndArgs() (query string, args []interface{}, ok bool) {
	if br.b != nil && br.ix < len(br.b.items) {
		bi := br.b.items[br.ix]
		query = bi.query
		args = bi.arguments
		ok = true
		br.ix++
	}
	return
}
