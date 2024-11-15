package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	p "github.com/SAP/go-hdb/driver/internal/protocol"
	"github.com/SAP/go-hdb/driver/internal/protocol/auth"
)

// ErrUnsupportedIsolationLevel is the error raised if a transaction is started with a not supported isolation level.
var ErrUnsupportedIsolationLevel = errors.New("unsupported isolation level")

// ErrNestedTransaction is the error raised if a transaction is created within a transaction as this is not supported by hdb.
var ErrNestedTransaction = errors.New("nested transactions are not supported")

// ErrNestedQuery is the error raised if a new sql statement is sent to the database server before the resultset
// processing of a previous sql query statement is finalized.
// Currently this only can happen if connections are used concurrently and if stream enabled fields (LOBs) are part
// of the resultset.
// This error can be avoided in whether using a transaction or a dedicated connection (sql.Tx or sql.Conn).
var ErrNestedQuery = errors.New("nested sql queries are not supported")

// queries.
const (
	pingQuery                       = "select 1 from dummy"
	setIsolationLevelReadCommitted  = "set transaction isolation level read committed"
	setIsolationLevelRepeatableRead = "set transaction isolation level repeatable read"
	setIsolationLevelSerializable   = "set transaction isolation level serializable"
	setAccessModeReadOnly           = "set transaction read only"
	setAccessModeReadWrite          = "set transaction read write"
)

var (
	// register as var to execute even before init() funcs are called.
	_ = p.RegisterScanType(p.DtBytes, reflect.TypeFor[[]byte](), reflect.TypeFor[NullBytes]())
	_ = p.RegisterScanType(p.DtDecimal, reflect.TypeFor[Decimal](), reflect.TypeFor[NullDecimal]())
	_ = p.RegisterScanType(p.DtLob, reflect.TypeFor[Lob](), reflect.TypeFor[NullLob]())
)

// check if conn implements all required interfaces.
var (
	_ driver.Conn               = (*conn)(nil)
	_ driver.ConnPrepareContext = (*conn)(nil)
	_ driver.Pinger             = (*conn)(nil)
	_ driver.ConnBeginTx        = (*conn)(nil)
	_ driver.ExecerContext      = (*conn)(nil)
	_ driver.QueryerContext     = (*conn)(nil)
	_ driver.NamedValueChecker  = (*conn)(nil)
	_ driver.SessionResetter    = (*conn)(nil)
	_ driver.Validator          = (*conn)(nil)
	_ Conn                      = (*conn)(nil) // go-hdb enhancements
)

// connection hook for testing.
// use unexported type to avoid key collisions.
type connHookCtxKeyType struct{}

var connHookCtxKey connHookCtxKeyType

// ...connection hook operations.
const (
	choNone = iota
	choStmtExec
)

// ...connection hook function.
type connHookFn func(op int)

func withConnHook(ctx context.Context, fn connHookFn) context.Context {
	return context.WithValue(ctx, connHookCtxKey, fn)
}

// Conn enhances a connection with go-hdb specific connection functions.
type Conn interface {
	HDBVersion() *Version
	DatabaseName() string
	DBConnectInfo(ctx context.Context, databaseName string) (*DBConnectInfo, error)
}

var stdConnTracker = &connTracker{}

type connTracker struct {
	mu      sync.Mutex
	_callDB *sql.DB
	numConn int64
}

func (t *connTracker) add() { t.mu.Lock(); t.numConn++; t.mu.Unlock() }

func (t *connTracker) remove() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.numConn--
	if t.numConn > 0 {
		return
	}
	t.numConn = 0
	if t._callDB != nil {
		t._callDB.Close()
		t._callDB = nil
	}
}

func (t *connTracker) callDB() *sql.DB {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t._callDB == nil {
		t._callDB = sql.OpenDB(new(callConnector))
	}
	return t._callDB
}

// Conn is the implementation of the database/sql/driver Conn interface.
type conn struct {
	attrs   *connAttrs
	metrics *metrics
	logger  *slog.Logger
	session *session
}

// isAuthError returns true in case of X509 certificate validation errrors or hdb authentication errors, else otherwise.
func isAuthError(err error) bool {
	var certValidationError *auth.CertValidationError
	if errors.As(err, &certValidationError) {
		return true
	}
	var hdbErrors *p.HdbErrors
	if !errors.As(err, &hdbErrors) {
		return false
	}
	return hdbErrors.Code() == p.HdbErrAuthenticationFailed
}

func fetchRedirectHost(ctx context.Context, host, databaseName string, metrics *metrics, attrs *connAttrs) (string, error) {
	c, err := newConn(ctx, host, metrics, attrs, nil)
	if err != nil {
		return "", err
	}
	defer c.Close()
	dbi, err := c.session.dbConnectInfo(ctx, databaseName)
	if err != nil {
		return "", err
	}
	if dbi.IsConnected { // if databaseName == "SYSTEMDB" and isConnected == true host and port are initial
		return host, nil
	}
	return net.JoinHostPort(dbi.Host, strconv.Itoa(dbi.Port)), nil
}

func connect(ctx context.Context, host string, metrics *metrics, connAttrs *connAttrs, authAttrs *authAttrs) (driver.Conn, error) {
	// can we connect via cookie?
	if auth := authAttrs.cookieAuth(); auth != nil {
		conn, err := newConn(ctx, host, metrics, connAttrs, auth)
		if err == nil {
			return conn, nil
		}
		if !isAuthError(err) {
			return nil, err
		}
		authAttrs.invalidateCookie() // cookie auth was not successful - do not try again with the same data
	}

	lastVersion := authAttrs.version.Load()
	for {
		authHnd := authAttrs.authHnd()

		conn, err := newConn(ctx, host, metrics, connAttrs, authHnd)
		if err == nil {
			if method, ok := authHnd.Selected().(auth.CookieGetter); ok {
				authAttrs.setCookie(method.Cookie())
			}
			return conn, nil
		}
		if !isAuthError(err) {
			return nil, err
		}

		if err := authAttrs.refresh(); err != nil {
			return nil, err
		}

		version := authAttrs.version.Load()
		if version == lastVersion { // no connection retry in case no new version available
			return nil, err
		}
		lastVersion = version
	}
}

// unique connection number.
var connNo atomic.Uint64

func newConn(ctx context.Context, host string, metrics *metrics, attrs *connAttrs, authHnd *p.AuthHnd) (*conn, error) {
	logger := attrs._logger.With(slog.Uint64("conn", connNo.Add(1)))

	metrics.lazyInit()

	session, err := newSession(ctx, host, logger, metrics, attrs, authHnd)
	if err != nil {
		return nil, err
	}

	stdConnTracker.add()
	metrics.msgCh <- gaugeMsg{idx: gaugeConn, v: 1} // increment open connections.

	return &conn{attrs: attrs, metrics: metrics, logger: logger, session: session}, nil
}

// Close implements the driver.Conn interface.
func (c *conn) Close() error {
	c.metrics.msgCh <- gaugeMsg{idx: gaugeConn, v: -1} // decrement open connections.
	stdConnTracker.remove()
	return c.session.close()
}

// ResetSession implements the driver.SessionResetter interface.
func (c *conn) ResetSession(ctx context.Context) error {
	if c.session.isBad() {
		return driver.ErrBadConn
	}

	lastRead := c.session.dbConn.lastRead()

	if c.attrs._pingInterval == 0 || lastRead.IsZero() || time.Since(lastRead) < c.attrs._pingInterval {
		return nil
	}

	if _, err := c.session.queryDirect(ctx, pingQuery); err != nil {
		return fmt.Errorf("%w: %w", driver.ErrBadConn, err)
	}
	return nil
}

// IsValid implements the driver.Validator interface.
func (c *conn) IsValid() bool { return !c.session.isBad() }

// Ping implements the driver.Pinger interface.
func (c *conn) Ping(ctx context.Context) error {
	_, err := c.session.queryDirect(ctx, pingQuery)
	return err
}

// PrepareContext implements the driver.ConnPrepareContext interface.
func (c *conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if err := c.session.switchUser(ctx); err != nil {
		return nil, err
	}

	pr, err := c.session.prepare(ctx, query)
	if err != nil {
		return nil, err
	}

	stmt := newStmt(c.session, c.attrs, c.metrics, query, pr)
	if stmtMetadata, ok := ctx.Value(stmtMetadataCtxKey).(*StmtMetadata); ok {
		*stmtMetadata = pr
	}
	return stmt, nil
}

// BeginTx implements the driver.ConnBeginTx interface.
func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if c.session.inTx {
		return nil, ErrNestedTransaction
	}

	if err := c.session.switchUser(ctx); err != nil {
		return nil, err
	}

	var isolationLevelQuery string
	switch sql.IsolationLevel(opts.Isolation) {
	case sql.LevelDefault, sql.LevelReadCommitted:
		isolationLevelQuery = setIsolationLevelReadCommitted
	case sql.LevelRepeatableRead:
		isolationLevelQuery = setIsolationLevelRepeatableRead
	case sql.LevelSerializable:
		isolationLevelQuery = setIsolationLevelSerializable
	default:
		return nil, ErrUnsupportedIsolationLevel
	}

	var accessModeQuery string
	if opts.ReadOnly {
		accessModeQuery = setAccessModeReadOnly
	} else {
		accessModeQuery = setAccessModeReadWrite
	}

	// set isolation level
	if _, err := c.session.execDirect(ctx, isolationLevelQuery); err != nil {
		return nil, err
	}
	// set access mode
	if _, err := c.session.execDirect(ctx, accessModeQuery); err != nil {
		return nil, err
	}

	tx := newTx(c)
	c.session.inTx = true

	return tx, nil
}

// QueryContext implements the driver.QueryerContext interface.
func (c *conn) QueryContext(ctx context.Context, query string, nvargs []driver.NamedValue) (driver.Rows, error) {
	// accepts stored procedures (call) without parameters to avoid parsing
	// the query string which might have comments, etc.
	if len(nvargs) != 0 {
		return nil, driver.ErrSkip // fast path not possible (prepare needed)
	}

	if err := c.session.switchUser(ctx); err != nil {
		return nil, err
	}
	return c.session.queryDirect(ctx, query)
}

// ExecContext implements the driver.ExecerContext interface.
func (c *conn) ExecContext(ctx context.Context, query string, nvargs []driver.NamedValue) (driver.Result, error) {
	if len(nvargs) != 0 {
		return nil, driver.ErrSkip // fast path not possible (prepare needed)
	}

	if err := c.session.switchUser(ctx); err != nil {
		return nil, err
	}
	// handle procedure call without parameters here as well
	return c.session.execDirect(ctx, query)
}

// CheckNamedValue implements the NamedValueChecker interface.
func (c *conn) CheckNamedValue(nv *driver.NamedValue) error {
	// - called by sql driver for ExecContext and QueryContext
	// - no check needs to be performed as ExecContext and QueryContext provided
	//   with parameters will force the 'prepare way' (driver.ErrSkip)
	// - Anyway, CheckNamedValue must be implemented to avoid default sql driver checks
	//   which would fail for custom arg types like Lob
	return nil
}

// Conn Raw access methods

// HDBVersion implements the Conn interface.
func (c *conn) HDBVersion() *Version { return c.session.hdbVersion }

// DatabaseName implements the Conn interface.
func (c *conn) DatabaseName() string { return c.session.databaseName }

// DBConnectInfo implements the Conn interface.
func (c *conn) DBConnectInfo(ctx context.Context, databaseName string) (*DBConnectInfo, error) {
	return c.session.dbConnectInfo(ctx, databaseName)
}

// transaction.

// check if tx implements all required interfaces.
var (
	_ driver.Tx = (*tx)(nil)
)

type tx struct {
	conn   *conn
	closed atomic.Bool
}

func newTx(conn *conn) *tx {
	conn.metrics.msgCh <- gaugeMsg{idx: gaugeTx, v: 1} // increment number of transactions.
	return &tx{conn: conn}
}

func (t *tx) Commit() error   { return t.close(false) }
func (t *tx) Rollback() error { return t.close(true) }

func (t *tx) close(rollback bool) error {
	c := t.conn

	c.metrics.msgCh <- gaugeMsg{idx: gaugeTx, v: -1} // decrement number of transactions.

	defer func() {
		c.session.inTx = false
	}()

	if c.session.isBad() {
		return driver.ErrBadConn
	}
	if closed := t.closed.Swap(true); closed {
		return nil
	}

	if rollback {
		return c.session.rollback(context.Background())
	}
	return c.session.commit(context.Background())
}
