package pgx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgconn"
)

// TxIsoLevel is the transaction isolation level (serializable, repeatable read, read committed or read uncommitted)
type TxIsoLevel string

// Transaction isolation levels
const (
	Serializable    TxIsoLevel = "serializable"
	RepeatableRead  TxIsoLevel = "repeatable read"
	ReadCommitted   TxIsoLevel = "read committed"
	ReadUncommitted TxIsoLevel = "read uncommitted"
)

// TxAccessMode is the transaction access mode (read write or read only)
type TxAccessMode string

// Transaction access modes
const (
	ReadWrite TxAccessMode = "read write"
	ReadOnly  TxAccessMode = "read only"
)

// TxDeferrableMode is the transaction deferrable mode (deferrable or not deferrable)
type TxDeferrableMode string

// Transaction deferrable modes
const (
	Deferrable    TxDeferrableMode = "deferrable"
	NotDeferrable TxDeferrableMode = "not deferrable"
)

// TxOptions are transaction modes within a transaction block
type TxOptions struct {
	IsoLevel       TxIsoLevel
	AccessMode     TxAccessMode
	DeferrableMode TxDeferrableMode
}

var emptyTxOptions TxOptions

func (txOptions TxOptions) beginSQL() string {
	if txOptions == emptyTxOptions {
		return "begin"
	}
	buf := &bytes.Buffer{}
	buf.WriteString("begin")
	if txOptions.IsoLevel != "" {
		fmt.Fprintf(buf, " isolation level %s", txOptions.IsoLevel)
	}
	if txOptions.AccessMode != "" {
		fmt.Fprintf(buf, " %s", txOptions.AccessMode)
	}
	if txOptions.DeferrableMode != "" {
		fmt.Fprintf(buf, " %s", txOptions.DeferrableMode)
	}

	return buf.String()
}

var ErrTxClosed = errors.New("tx is closed")

// ErrTxCommitRollback occurs when an error has occurred in a transaction and
// Commit() is called. PostgreSQL accepts COMMIT on aborted transactions, but
// it is treated as ROLLBACK.
var ErrTxCommitRollback = errors.New("commit unexpectedly resulted in rollback")

// Begin starts a transaction. Unlike database/sql, the context only affects the begin command. i.e. there is no
// auto-rollback on context cancellation.
func (c *Conn) Begin(ctx context.Context) (Tx, error) {
	return c.BeginTx(ctx, TxOptions{})
}

// BeginTx starts a transaction with txOptions determining the transaction mode. Unlike database/sql, the context only
// affects the begin command. i.e. there is no auto-rollback on context cancellation.
func (c *Conn) BeginTx(ctx context.Context, txOptions TxOptions) (Tx, error) {
	_, err := c.Exec(ctx, txOptions.beginSQL())
	if err != nil {
		// begin should never fail unless there is an underlying connection issue or
		// a context timeout. In either case, the connection is possibly broken.
		c.die(errors.New("failed to begin transaction"))
		return nil, err
	}

	return &dbTx{conn: c}, nil
}

// BeginFunc starts a transaction and calls f. If f does not return an error the transaction is committed. If f returns
// an error the transaction is rolled back. The context will be used when executing the transaction control statements
// (BEGIN, ROLLBACK, and COMMIT) but does not otherwise affect the execution of f.
func (c *Conn) BeginFunc(ctx context.Context, f func(Tx) error) (err error) {
	return c.BeginTxFunc(ctx, TxOptions{}, f)
}

// BeginTxFunc starts a transaction with txOptions determining the transaction mode and calls f. If f does not return
// an error the transaction is committed. If f returns an error the transaction is rolled back. The context will be
// used when executing the transaction control statements (BEGIN, ROLLBACK, and COMMIT) but does not otherwise affect
// the execution of f.
func (c *Conn) BeginTxFunc(ctx context.Context, txOptions TxOptions, f func(Tx) error) (err error) {
	var tx Tx
	tx, err = c.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}
	defer func() {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, ErrTxClosed) {
			err = rollbackErr
		}
	}()

	fErr := f(tx)
	if fErr != nil {
		_ = tx.Rollback(ctx) // ignore rollback error as there is already an error to return
		return fErr
	}

	return tx.Commit(ctx)
}

// Tx represents a database transaction.
//
// Tx is an interface instead of a struct to enable connection pools to be implemented without relying on internal pgx
// state, to support pseudo-nested transactions with savepoints, and to allow tests to mock transactions. However,
// adding a method to an interface is technically a breaking change. If new methods are added to Conn it may be
// desirable to add them to Tx as well. Because of this the Tx interface is partially excluded from semantic version
// requirements. Methods will not be removed or changed, but new methods may be added.
type Tx interface {
	// Begin starts a pseudo nested transaction.
	Begin(ctx context.Context) (Tx, error)

	// BeginFunc starts a pseudo nested transaction and executes f. If f does not return an err the pseudo nested
	// transaction will be committed. If it does then it will be rolled back.
	BeginFunc(ctx context.Context, f func(Tx) error) (err error)

	// Commit commits the transaction if this is a real transaction or releases the savepoint if this is a pseudo nested
	// transaction. Commit will return ErrTxClosed if the Tx is already closed, but is otherwise safe to call multiple
	// times. If the commit fails with a rollback status (e.g. the transaction was already in a broken state) then
	// ErrTxCommitRollback will be returned.
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction if this is a real transaction or rolls back to the savepoint if this is a
	// pseudo nested transaction. Rollback will return ErrTxClosed if the Tx is already closed, but is otherwise safe to
	// call multiple times. Hence, a defer tx.Rollback() is safe even if tx.Commit() will be called first in a non-error
	// condition. Any other failure of a real transaction will result in the connection being closed.
	Rollback(ctx context.Context) error

	CopyFrom(ctx context.Context, tableName Identifier, columnNames []string, rowSrc CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *Batch) BatchResults
	LargeObjects() LargeObjects

	Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error)

	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(QueryFuncRow) error) (pgconn.CommandTag, error)

	// Conn returns the underlying *Conn that on which this transaction is executing.
	Conn() *Conn
}

// dbTx represents a database transaction.
//
// All dbTx methods return ErrTxClosed if Commit or Rollback has already been
// called on the dbTx.
type dbTx struct {
	conn         *Conn
	err          error
	savepointNum int64
	closed       bool
}

// Begin starts a pseudo nested transaction implemented with a savepoint.
func (tx *dbTx) Begin(ctx context.Context) (Tx, error) {
	if tx.closed {
		return nil, ErrTxClosed
	}

	tx.savepointNum++
	_, err := tx.conn.Exec(ctx, "savepoint sp_"+strconv.FormatInt(tx.savepointNum, 10))
	if err != nil {
		return nil, err
	}

	return &dbSimulatedNestedTx{tx: tx, savepointNum: tx.savepointNum}, nil
}

func (tx *dbTx) BeginFunc(ctx context.Context, f func(Tx) error) (err error) {
	if tx.closed {
		return ErrTxClosed
	}

	var savepoint Tx
	savepoint, err = tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		rollbackErr := savepoint.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, ErrTxClosed) {
			err = rollbackErr
		}
	}()

	fErr := f(savepoint)
	if fErr != nil {
		_ = savepoint.Rollback(ctx) // ignore rollback error as there is already an error to return
		return fErr
	}

	return savepoint.Commit(ctx)
}

// Commit commits the transaction.
func (tx *dbTx) Commit(ctx context.Context) error {
	if tx.closed {
		return ErrTxClosed
	}

	commandTag, err := tx.conn.Exec(ctx, "commit")
	tx.closed = true
	if err != nil {
		if tx.conn.PgConn().TxStatus() != 'I' {
			_ = tx.conn.Close(ctx) // already have error to return
		}
		return err
	}
	if string(commandTag) == "ROLLBACK" {
		return ErrTxCommitRollback
	}

	return nil
}

// Rollback rolls back the transaction. Rollback will return ErrTxClosed if the
// Tx is already closed, but is otherwise safe to call multiple times. Hence, a
// defer tx.Rollback() is safe even if tx.Commit() will be called first in a
// non-error condition.
func (tx *dbTx) Rollback(ctx context.Context) error {
	if tx.closed {
		return ErrTxClosed
	}

	_, err := tx.conn.Exec(ctx, "rollback")
	tx.closed = true
	if err != nil {
		// A rollback failure leaves the connection in an undefined state
		tx.conn.die(fmt.Errorf("rollback failed: %w", err))
		return err
	}

	return nil
}

// Exec delegates to the underlying *Conn
func (tx *dbTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error) {
	if tx.closed {
		return pgconn.CommandTag{}, ErrTxClosed
	}

	return tx.conn.Exec(ctx, sql, arguments...)
}

// Prepare delegates to the underlying *Conn
func (tx *dbTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	if tx.closed {
		return nil, ErrTxClosed
	}

	return tx.conn.Prepare(ctx, name, sql)
}

// Query delegates to the underlying *Conn
func (tx *dbTx) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	if tx.closed {
		// Because checking for errors can be deferred to the *Rows, build one with the error
		err := ErrTxClosed
		return &connRows{closed: true, err: err}, err
	}

	return tx.conn.Query(ctx, sql, args...)
}

// QueryRow delegates to the underlying *Conn
func (tx *dbTx) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	rows, _ := tx.Query(ctx, sql, args...)
	return (*connRow)(rows.(*connRows))
}

// QueryFunc delegates to the underlying *Conn.
func (tx *dbTx) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(QueryFuncRow) error) (pgconn.CommandTag, error) {
	if tx.closed {
		return nil, ErrTxClosed
	}

	return tx.conn.QueryFunc(ctx, sql, args, scans, f)
}

// CopyFrom delegates to the underlying *Conn
func (tx *dbTx) CopyFrom(ctx context.Context, tableName Identifier, columnNames []string, rowSrc CopyFromSource) (int64, error) {
	if tx.closed {
		return 0, ErrTxClosed
	}

	return tx.conn.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

// SendBatch delegates to the underlying *Conn
func (tx *dbTx) SendBatch(ctx context.Context, b *Batch) BatchResults {
	if tx.closed {
		return &batchResults{err: ErrTxClosed}
	}

	return tx.conn.SendBatch(ctx, b)
}

// LargeObjects returns a LargeObjects instance for the transaction.
func (tx *dbTx) LargeObjects() LargeObjects {
	return LargeObjects{tx: tx}
}

func (tx *dbTx) Conn() *Conn {
	return tx.conn
}

// dbSimulatedNestedTx represents a simulated nested transaction implemented by a savepoint.
type dbSimulatedNestedTx struct {
	tx           Tx
	savepointNum int64
	closed       bool
}

// Begin starts a pseudo nested transaction implemented with a savepoint.
func (sp *dbSimulatedNestedTx) Begin(ctx context.Context) (Tx, error) {
	if sp.closed {
		return nil, ErrTxClosed
	}

	return sp.tx.Begin(ctx)
}

func (sp *dbSimulatedNestedTx) BeginFunc(ctx context.Context, f func(Tx) error) (err error) {
	if sp.closed {
		return ErrTxClosed
	}

	return sp.tx.BeginFunc(ctx, f)
}

// Commit releases the savepoint essentially committing the pseudo nested transaction.
func (sp *dbSimulatedNestedTx) Commit(ctx context.Context) error {
	if sp.closed {
		return ErrTxClosed
	}

	_, err := sp.Exec(ctx, "release savepoint sp_"+strconv.FormatInt(sp.savepointNum, 10))
	sp.closed = true
	return err
}

// Rollback rolls back to the savepoint essentially rolling back the pseudo nested transaction. Rollback will return
// ErrTxClosed if the dbSavepoint is already closed, but is otherwise safe to call multiple times. Hence, a defer sp.Rollback()
// is safe even if sp.Commit() will be called first in a non-error condition.
func (sp *dbSimulatedNestedTx) Rollback(ctx context.Context) error {
	if sp.closed {
		return ErrTxClosed
	}

	_, err := sp.Exec(ctx, "rollback to savepoint sp_"+strconv.FormatInt(sp.savepointNum, 10))
	sp.closed = true
	return err
}

// Exec delegates to the underlying Tx
func (sp *dbSimulatedNestedTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error) {
	if sp.closed {
		return nil, ErrTxClosed
	}

	return sp.tx.Exec(ctx, sql, arguments...)
}

// Prepare delegates to the underlying Tx
func (sp *dbSimulatedNestedTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	if sp.closed {
		return nil, ErrTxClosed
	}

	return sp.tx.Prepare(ctx, name, sql)
}

// Query delegates to the underlying Tx
func (sp *dbSimulatedNestedTx) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	if sp.closed {
		// Because checking for errors can be deferred to the *Rows, build one with the error
		err := ErrTxClosed
		return &connRows{closed: true, err: err}, err
	}

	return sp.tx.Query(ctx, sql, args...)
}

// QueryRow delegates to the underlying Tx
func (sp *dbSimulatedNestedTx) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	rows, _ := sp.Query(ctx, sql, args...)
	return (*connRow)(rows.(*connRows))
}

// QueryFunc delegates to the underlying Tx.
func (sp *dbSimulatedNestedTx) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(QueryFuncRow) error) (pgconn.CommandTag, error) {
	if sp.closed {
		return nil, ErrTxClosed
	}

	return sp.tx.QueryFunc(ctx, sql, args, scans, f)
}

// CopyFrom delegates to the underlying *Conn
func (sp *dbSimulatedNestedTx) CopyFrom(ctx context.Context, tableName Identifier, columnNames []string, rowSrc CopyFromSource) (int64, error) {
	if sp.closed {
		return 0, ErrTxClosed
	}

	return sp.tx.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

// SendBatch delegates to the underlying *Conn
func (sp *dbSimulatedNestedTx) SendBatch(ctx context.Context, b *Batch) BatchResults {
	if sp.closed {
		return &batchResults{err: ErrTxClosed}
	}

	return sp.tx.SendBatch(ctx, b)
}

func (sp *dbSimulatedNestedTx) LargeObjects() LargeObjects {
	return LargeObjects{tx: sp}
}

func (sp *dbSimulatedNestedTx) Conn() *Conn {
	return sp.tx.Conn()
}
