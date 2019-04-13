package pgx

import (
	"context"

	"github.com/jackc/pgx/pgproto3"
	"github.com/jackc/pgx/pgtype"
)

type batchItem struct {
	query             string
	arguments         []interface{}
	parameterOIDs     []pgtype.OID
	resultFormatCodes []int16
}

// Batch queries are a way of bundling multiple queries together to avoid
// unnecessary network round trips.
type Batch struct {
	conn                   *Conn
	connPool               *ConnPool
	items                  []*batchItem
	resultsRead            int
	pendingCommandComplete bool
	ctx                    context.Context
	err                    error
	inTx                   bool
}

// BeginBatch returns a *Batch query for c.
func (c *Conn) BeginBatch() *Batch {
	return &Batch{conn: c}
}

// BeginBatch returns a *Batch query for tx. Since this *Batch is already part
// of a transaction it will not automatically be wrapped in a transaction.
func (tx *Tx) BeginBatch() *Batch {
	return &Batch{conn: tx.conn, inTx: true}
}

// Conn returns the underlying connection that b will or was performed on.
func (b *Batch) Conn() *Conn {
	return b.conn
}

// Queue queues a query to batch b. parameterOIDs are required if there are
// parameters and query is not the name of a prepared statement.
// resultFormatCodes are required if there is a result.
func (b *Batch) Queue(query string, arguments []interface{}, parameterOIDs []pgtype.OID, resultFormatCodes []int16) {
	b.items = append(b.items, &batchItem{
		query:             query,
		arguments:         arguments,
		parameterOIDs:     parameterOIDs,
		resultFormatCodes: resultFormatCodes,
	})
}

// Send sends all queued queries to the server at once.
// If the batch is created from a conn Object then All queries are wrapped
// in a transaction. The transaction can optionally be configured with
// txOptions. The context is in effect until the Batch is closed.
//
// Warning: Send writes all queued queries before reading any results. This can
// cause a deadlock if an excessive number of queries are queued. It is highly
// advisable to use a timeout context to protect against this possibility.
// Unfortunately, this excessive number can vary based on operating system,
// connection type (TCP or Unix domain socket), and type of query. Unix domain
// sockets seem to be much more susceptible to this issue than TCP connections.
// However, it usually is at least several thousand.
//
// The deadlock occurs when the batched queries to be sent are so large that the
// PostgreSQL server cannot receive it all at once. PostgreSQL received some of
// the queued queries and starts executing them. As PostgreSQL executes the
// queries it sends responses back. pgx will not read any of these responses
// until it has finished sending. Therefore, if all network buffers are full pgx
// will not be able to finish sending the queries and PostgreSQL will not be
// able to finish sending the responses.
//
// See https://github.com/jackc/pgx/issues/374.
func (b *Batch) Send(ctx context.Context, txOptions *TxOptions) error {
	if b.err != nil {
		return b.err
	}

	b.ctx = ctx

	err := b.conn.waitForPreviousCancelQuery(ctx)
	if err != nil {
		return err
	}

	if err := b.conn.ensureConnectionReadyForQuery(); err != nil {
		return err
	}

	buf := b.conn.wbuf
	if !b.inTx {
		buf = appendQuery(buf, txOptions.beginSQL())
	}

	err = b.conn.initContext(ctx)
	if err != nil {
		return err
	}

	for _, bi := range b.items {
		var psName string
		var psParameterOIDs []pgtype.OID

		if ps, ok := b.conn.preparedStatements[bi.query]; ok {
			psName = ps.Name
			psParameterOIDs = ps.ParameterOIDs
		} else {
			psParameterOIDs = bi.parameterOIDs
			buf = appendParse(buf, "", bi.query, psParameterOIDs)
		}

		var err error
		buf, err = appendBind(buf, "", psName, b.conn.ConnInfo, psParameterOIDs, bi.arguments, bi.resultFormatCodes)
		if err != nil {
			return err
		}

		buf = appendDescribe(buf, 'P', "")
		buf = appendExecute(buf, "", 0)
	}

	buf = appendSync(buf)
	b.conn.pendingReadyForQueryCount++

	if !b.inTx {
		buf = appendQuery(buf, "commit")
		b.conn.pendingReadyForQueryCount++
	}

	n, err := b.conn.conn.Write(buf)
	if err != nil {
		if fatalWriteErr(n, err) {
			b.conn.die(err)
		}
		return err
	}

	for !b.inTx {
		msg, err := b.conn.rxMsg()
		if err != nil {
			return err
		}

		switch msg := msg.(type) {
		case *pgproto3.ReadyForQuery:
			return nil
		default:
			if err := b.conn.processContextFreeMsg(msg); err != nil {
				return err
			}
		}
	}

	return nil
}

// ExecResults reads the results from the next query in the batch as if the
// query has been sent with Exec.
func (b *Batch) ExecResults() (CommandTag, error) {
	if b.err != nil {
		return "", b.err
	}

	select {
	case <-b.ctx.Done():
		b.die(b.ctx.Err())
		return "", b.ctx.Err()
	default:
	}

	if err := b.ensureCommandComplete(); err != nil {
		b.die(err)
		return "", err
	}

	b.resultsRead++

	b.pendingCommandComplete = true

	for {
		msg, err := b.conn.rxMsg()
		if err != nil {
			return "", err
		}

		switch msg := msg.(type) {
		case *pgproto3.CommandComplete:
			b.pendingCommandComplete = false
			return CommandTag(msg.CommandTag), nil
		default:
			if err := b.conn.processContextFreeMsg(msg); err != nil {
				return "", err
			}
		}
	}
}

// QueryResults reads the results from the next query in the batch as if the
// query has been sent with Query.
func (b *Batch) QueryResults() (*Rows, error) {
	rows := b.conn.getRows("batch query", nil)

	if b.err != nil {
		rows.fatal(b.err)
		return rows, b.err
	}

	select {
	case <-b.ctx.Done():
		b.die(b.ctx.Err())
		rows.fatal(b.err)
		return rows, b.ctx.Err()
	default:
	}

	if err := b.ensureCommandComplete(); err != nil {
		b.die(err)
		rows.fatal(err)
		return rows, err
	}

	b.resultsRead++

	b.pendingCommandComplete = true

	fieldDescriptions, err := b.conn.readUntilRowDescription()
	if err != nil {
		b.die(err)
		rows.fatal(b.err)
		return rows, err
	}

	rows.batch = b
	rows.fields = fieldDescriptions
	return rows, nil
}

// QueryRowResults reads the results from the next query in the batch as if the
// query has been sent with QueryRow.
func (b *Batch) QueryRowResults() *Row {
	rows, _ := b.QueryResults()
	return (*Row)(rows)

}

// Close closes the batch operation. Any error that occured during a batch
// operation may have made it impossible to resyncronize the connection with the
// server. In this case the underlying connection will have been closed.
func (b *Batch) Close() (err error) {
	if b.err != nil {
		return b.err
	}

	defer func() {
		err = b.conn.termContext(err)
		if b.conn != nil && b.connPool != nil {
			b.connPool.Release(b.conn)
		}
	}()

	for i := b.resultsRead; i < len(b.items); i++ {
		if _, err = b.ExecResults(); err != nil {
			return err
		}
	}

	if err = b.conn.ensureConnectionReadyForQuery(); err != nil {
		return err
	}

	return nil
}

func (b *Batch) die(err error) {
	if b.err != nil {
		return
	}

	b.err = err
	b.conn.die(err)

	if b.conn != nil && b.connPool != nil {
		b.connPool.Release(b.conn)
	}
}

func (b *Batch) ensureCommandComplete() error {
	for b.pendingCommandComplete {
		msg, err := b.conn.rxMsg()
		if err != nil {
			return err
		}

		switch msg := msg.(type) {
		case *pgproto3.CommandComplete:
			b.pendingCommandComplete = false
			return nil
		default:
			err = b.conn.processContextFreeMsg(msg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
