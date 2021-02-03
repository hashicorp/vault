// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"bufio"
	"context"
	"crypto/tls"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/SAP/go-hdb/driver/common"
	"github.com/SAP/go-hdb/driver/dial"
	"github.com/SAP/go-hdb/driver/sqltrace"
	p "github.com/SAP/go-hdb/internal/protocol"
	"github.com/SAP/go-hdb/internal/protocol/scanner"
)

// Transaction isolation levels supported by hdb.
const (
	LevelReadCommitted  = "READ COMMITTED"
	LevelRepeatableRead = "REPEATABLE READ"
	LevelSerializable   = "SERIALIZABLE"
)

// Access modes supported by hdb.
const (
	modeReadOnly  = "READ ONLY"
	modeReadWrite = "READ WRITE"
)

// map sql isolation level to hdb isolation level.
var isolationLevel = map[driver.IsolationLevel]string{
	driver.IsolationLevel(sql.LevelDefault):        LevelReadCommitted,
	driver.IsolationLevel(sql.LevelReadCommitted):  LevelReadCommitted,
	driver.IsolationLevel(sql.LevelRepeatableRead): LevelRepeatableRead,
	driver.IsolationLevel(sql.LevelSerializable):   LevelSerializable,
}

// map sql read only flag to hdb access mode.
var readOnly = map[bool]string{
	true:  modeReadOnly,
	false: modeReadWrite,
}

// ErrUnsupportedIsolationLevel is the error raised if a transaction is started with a not supported isolation level.
var ErrUnsupportedIsolationLevel = errors.New("unsupported isolation level")

// ErrNestedTransaction is the error raised if a transaction is created within a transaction as this is not supported by hdb.
var ErrNestedTransaction = errors.New("nested transactions are not supported")

// ErrNestedQuery is the error raised if a sql statement is executed before an "active" statement is closed.
// Example: execute sql statement before rows of previous select statement are closed.
var ErrNestedQuery = errors.New("nested sql queries are not supported")

// queries
const (
	dummyQuery        = "select 1 from dummy"
	setIsolationLevel = "set transaction isolation level %s"
	setAccessMode     = "set transaction %s"
	setDefaultSchema  = "set schema %s"
)

// bulk statement
const (
	bulk = "b$"
)

var (
	flushTok   = new(struct{})
	noFlushTok = new(struct{})
)

var (
	// NoFlush is to be used as parameter in bulk statements to delay execution.
	NoFlush = sql.Named(bulk, &noFlushTok)
	// Flush can be used as optional parameter in bulk statements but is not required to trigger execution.
	Flush = sql.Named(bulk, &flushTok)
)

const (
	maxNumTraceArg = 20
)

func init() {
	p.RegisterScanType(p.DtDecimal, reflect.TypeOf((*Decimal)(nil)).Elem())
	p.RegisterScanType(p.DtLob, reflect.TypeOf((*Lob)(nil)).Elem())
}

// dbConn wraps the database tcp connection. It sets timeouts and handles driver ErrBadConn behavior.
type dbConn struct {
	// atomic access - alignment
	canceled  int32
	conn      net.Conn
	timeout   time.Duration
	lastError error // error bad connection
	closed    bool
}

func (c *dbConn) isBad() bool {
	return c.lastError != nil
}

func (c *dbConn) deadline() (deadline time.Time) {
	if c.timeout == 0 {
		return
	}
	return time.Now().Add(c.timeout)
}

var (
	errCancelled = errors.New("db connection is canceled")
	errClosed    = errors.New("db connection is closed")
)

func (c *dbConn) cancel() {
	atomic.StoreInt32(&c.canceled, 1)
	c.lastError = errCancelled
}

func (c *dbConn) close() error {
	c.closed = true
	c.lastError = errClosed
	return c.conn.Close()
}

// Read implements the io.Reader interface.
func (c *dbConn) Read(b []byte) (n int, err error) {
	// check if killed
	if atomic.LoadInt32(&c.canceled) == 1 {
		return 0, driver.ErrBadConn
	}
	//set timeout
	if err = c.conn.SetReadDeadline(c.deadline()); err != nil {
		goto retError
	}
	if n, err = c.conn.Read(b); err != nil {
		goto retError
	}
	return
retError:
	dlog.Printf("Connection read error local address %s remote address %s: %s", c.conn.LocalAddr(), c.conn.RemoteAddr(), err)
	c.lastError = err
	return n, driver.ErrBadConn
}

// Write implements the io.Writer interface.
func (c *dbConn) Write(b []byte) (n int, err error) {
	// check if killed
	if atomic.LoadInt32(&c.canceled) == 1 {
		return 0, driver.ErrBadConn
	}
	//set timeout
	if err = c.conn.SetWriteDeadline(c.deadline()); err != nil {
		goto retError
	}
	if n, err = c.conn.Write(b); err != nil {
		goto retError
	}
	return
retError:
	dlog.Printf("Connection write error local address %s remote address %s: %s", c.conn.LocalAddr(), c.conn.RemoteAddr(), err)
	c.lastError = err
	return n, driver.ErrBadConn
}

const (
	lrNestedQuery = 1
)

type connLock struct {
	// 64 bit alignment
	lockReason int64 // atomic access

	mu     sync.Mutex // tryLock mutex
	connMu sync.Mutex // connection mutex
}

func (l *connLock) tryLock(lockReason int64) error {
	l.mu.Lock()
	if atomic.LoadInt64(&l.lockReason) == lrNestedQuery {
		l.mu.Unlock()
		return ErrNestedQuery
	}
	l.connMu.Lock()
	atomic.StoreInt64(&l.lockReason, lockReason)
	l.mu.Unlock()
	return nil
}

func (l *connLock) lock() { l.connMu.Lock() }

func (l *connLock) unlock() {
	atomic.StoreInt64(&l.lockReason, 0)
	l.connMu.Unlock()
}

//  check if conn implements all required interfaces
var (
	_ driver.Conn               = (*Conn)(nil)
	_ driver.ConnPrepareContext = (*Conn)(nil)
	_ driver.Pinger             = (*Conn)(nil)
	_ driver.ConnBeginTx        = (*Conn)(nil)
	_ driver.ExecerContext      = (*Conn)(nil)
	_ driver.Execer             = (*Conn)(nil) //go 1.9 issue (ExecerContext is only called if Execer is implemented)
	_ driver.QueryerContext     = (*Conn)(nil)
	_ driver.Queryer            = (*Conn)(nil) //go 1.9 issue (QueryerContext is only called if Queryer is implemented)
	_ driver.NamedValueChecker  = (*Conn)(nil)
	_ driver.SessionResetter    = (*Conn)(nil)
)

// connHook is a hook for testing.
var connHook func(c *Conn, op int)

// connection hook operations
const (
	choNone = iota
	choStmtExec
)

// Conn is the implementation of the database/sql/driver Conn interface.
type Conn struct {
	// Holding connection lock in QueryResultSet (see rows.onClose)
	/*
		As long as a session is in query mode no other sql statement must be executed.
		Example:
		- pinger is active
		- select with blob fields is executed
		- scan is hitting the database again (blob streaming)
		- if in between a ping gets executed (ping selects db) hdb raises error
		  "SQL Error 1033 - error while parsing protocol: invalid lob locator id (piecewise lob reading)"
	*/
	connLock

	ctr     *Connector
	dbConn  *dbConn
	session *p.Session
	scanner *scanner.Scanner
	closed  chan struct{}

	inTx bool // in transaction
}

func newConn(ctx context.Context, ctr *Connector) (driver.Conn, error) {

	ctr.mu.RLock() // lock connector
	defer ctr.mu.RUnlock()

	conn, err := ctr.dialer.DialContext(ctx, ctr.host, dial.DialerOptions{Timeout: ctr.timeout, TCPKeepAlive: ctr.tcpKeepAlive})
	if err != nil {
		return nil, err
	}

	// is TLS connection requested?
	if ctr.tlsConfig != nil {
		conn = tls.Client(conn, ctr.tlsConfig)
	}

	dbConn := &dbConn{conn: conn, timeout: ctr.timeout}
	// buffer connection
	rw := bufio.NewReadWriter(bufio.NewReaderSize(dbConn, ctr.bufferSize), bufio.NewWriterSize(dbConn, ctr.bufferSize))

	session, err := p.NewSession(ctx, rw,
		&p.SessionConfig{
			DriverVersion:    DriverVersion,
			DriverName:       DriverName,
			ApplicationName:  ctr.applicationName,
			Username:         ctr.username,
			Password:         ctr.password,
			SessionVariables: ctr.sessionVariables,
			Locale:           ctr.locale,
			FetchSize:        ctr.fetchSize,
			LobChunkSize:     ctr.lobChunkSize,
			Dfv:              ctr.dfv,
			Legacy:           ctr.legacy,
		},
	)
	if err != nil {
		return nil, err
	}

	c := &Conn{ctr: ctr, dbConn: dbConn, session: session, scanner: &scanner.Scanner{}, closed: make(chan struct{})}
	if ctr.defaultSchema != "" {
		if _, err := c.ExecContext(ctx, fmt.Sprintf(setDefaultSchema, Identifier(ctr.defaultSchema)), nil); err != nil {
			return nil, err
		}
	}

	if ctr.pingInterval != 0 {
		go c.pinger(ctr.pingInterval, c.closed)
	}

	hdbDriver.addConn(1) // increment open connections.

	return c, nil
}

func (c *Conn) pinger(d time.Duration, done <-chan struct{}) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	ctx := context.Background()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			c.Ping(ctx)
		}
	}
}

// Ping implements the driver.Pinger interface.
func (c *Conn) Ping(ctx context.Context) (err error) {
	if err := c.tryLock(0); err != nil {
		return err
	}
	defer c.unlock()

	if c.dbConn.isBad() {
		return driver.ErrBadConn
	}

	done := make(chan struct{})
	go func() {
		_, err = c.session.QueryDirect(dummyQuery, !c.inTx)
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.dbConn.cancel()
		return ctx.Err()
	case <-done:
		return err
	}
}

// ResetSession implements the driver.SessionResetter interface.
func (c *Conn) ResetSession(ctx context.Context) error {
	c.lock()
	defer c.unlock()

	p.QueryResultCache.Cleanup(c.session)

	if c.dbConn.isBad() {
		return driver.ErrBadConn
	}
	return nil
}

// PrepareContext implements the driver.ConnPrepareContext interface.
func (c *Conn) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	if err := c.tryLock(0); err != nil {
		return nil, err
	}
	defer c.unlock()

	if c.dbConn.isBad() {
		return nil, driver.ErrBadConn
	}

	done := make(chan struct{})
	go func() {
		var (
			qd *p.QueryDescr
			pr *p.PrepareResult
		)

		qd, err = p.NewQueryDescr(query, c.scanner)
		if err != nil {
			goto done
		}
		pr, err = c.session.Prepare(qd.Query())
		if err != nil {
			goto done
		}

		if err = pr.Check(qd); err != nil {
			goto done
		}

		select {
		default:
		case <-ctx.Done():
			return
		}

		if pr.IsProcedureCall() {
			stmt = newCallStmt(c, qd.Query(), pr)
		} else {
			stmt = newStmt(c, qd.Query(), qd.IsBulk(), c.ctr.BulkSize(), pr) //take latest connector bulk size
		}

	done:
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.dbConn.cancel()
		return nil, ctx.Err()
	case <-done:
		hdbDriver.addStmt(1) // increment number of statements.
		return stmt, err
	}
}

// Close implements the driver.Conn interface.
func (c *Conn) Close() error {
	c.lock()
	defer c.unlock()

	hdbDriver.addConn(-1) // decrement open connections.
	close(c.closed)       // signal connection close

	// cleanup query cache
	p.QueryResultCache.Cleanup(c.session)

	// if isBad do not disconnect
	if !c.dbConn.isBad() {
		c.session.Disconnect() // ignore error
	}
	return c.dbConn.close()
}

// BeginTx implements the driver.ConnBeginTx interface.
func (c *Conn) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {
	if err := c.tryLock(0); err != nil {
		return nil, err
	}
	defer c.unlock()

	if c.dbConn.isBad() {
		return nil, driver.ErrBadConn
	}

	if c.inTx {
		return nil, ErrNestedTransaction
	}

	level, ok := isolationLevel[opts.Isolation]
	if !ok {
		return nil, ErrUnsupportedIsolationLevel
	}

	done := make(chan struct{})
	go func() {
		// set isolation level
		if _, err = c.session.ExecDirect(fmt.Sprintf(setIsolationLevel, level), !c.inTx); err != nil {
			goto done
		}
		// set access mode
		if _, err = c.session.ExecDirect(fmt.Sprintf(setAccessMode, readOnly[opts.ReadOnly]), !c.inTx); err != nil {
			goto done
		}
		c.inTx = true
		tx = newTx(c)
	done:
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.dbConn.cancel()
		return nil, ctx.Err()
	case <-done:
		hdbDriver.addTx(1) // increment number of transactions.
		return tx, err
	}
}

// QueryContext implements the driver.QueryerContext interface.
func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	if err := c.tryLock(lrNestedQuery); err != nil {
		return nil, err
	}
	hasRowsCloser := false
	defer func() {
		// unlock connection if rows will not do it
		if !hasRowsCloser {
			c.unlock()
		}
	}()

	if c.dbConn.isBad() {
		return nil, driver.ErrBadConn
	}

	if len(args) != 0 {
		return nil, driver.ErrSkip //fast path not possible (prepare needed)
	}

	qd, err := p.NewQueryDescr(query, c.scanner)
	if err != nil {
		return nil, err
	}
	switch qd.Kind() {
	case p.QkCall:
		// direct execution of call procedure
		// - returns no parameter metadata (sps 82) but only field values
		// --> let's take the 'prepare way' for stored procedures
		return nil, driver.ErrSkip
	case p.QkID:
		// query call table result
		rows, ok := p.QueryResultCache.Get(qd.ID())
		if !ok {
			return nil, fmt.Errorf("invalid result set id %s", query)
		}
		if onCloser, ok := rows.(p.OnCloser); ok {
			onCloser.SetOnClose(c.unlock)
			hasRowsCloser = true
		}
		return rows, nil
	}

	if sqltrace.On() {
		sqltrace.Traceln(query)
	}

	done := make(chan struct{})
	go func() {
		rows, err = c.session.QueryDirect(query, !c.inTx)
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.dbConn.cancel()
		return nil, ctx.Err()
	case <-done:
		if onCloser, ok := rows.(p.OnCloser); ok {
			onCloser.SetOnClose(c.unlock)
			hasRowsCloser = true
		}
		return rows, err
	}
}

// ExecContext implements the driver.ExecerContext interface.
func (c *Conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (r driver.Result, err error) {
	if err := c.tryLock(0); err != nil {
		return nil, err
	}
	defer c.unlock()

	if c.dbConn.isBad() {
		return nil, driver.ErrBadConn
	}

	if len(args) != 0 {
		return nil, driver.ErrSkip //fast path not possible (prepare needed)
	}

	qd, err := p.NewQueryDescr(query, c.scanner)
	if err != nil {
		return nil, err
	}

	if sqltrace.On() {
		sqltrace.Traceln(query)
	}

	done := make(chan struct{})
	go func() {
		/*
			handle call procedure (qd.Kind() == p.QkCall) without parameters here as well
		*/
		r, err = c.session.ExecDirect(qd.Query(), !c.inTx)
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.dbConn.cancel()
		return nil, ctx.Err()
	case <-done:
		return r, err
	}
}

// CheckNamedValue implements the NamedValueChecker interface.
func (c *Conn) CheckNamedValue(nv *driver.NamedValue) error {
	// - called by sql driver for ExecContext and QueryContext
	// - no check needs to be performed as ExecContext and QueryContext provided
	//   with parameters will force the 'prepare way' (driver.ErrSkip)
	// - Anyway, CheckNamedValue must be implemented to avoid default sql driver checks
	//   which would fail for custom arg types like Lob
	return nil
}

// Conn Raw access methods

// ServerInfo returns parameters reported by hdb server.
func (c *Conn) ServerInfo() *common.ServerInfo {
	return c.session.ServerInfo()
}

//transaction

//  check if tx implements all required interfaces
var (
	_ driver.Tx = (*tx)(nil)
)

type tx struct {
	conn   *Conn
	closed bool
}

func newTx(conn *Conn) *tx { return &tx{conn: conn} }

func (t *tx) Commit() error   { return t.close(false) }
func (t *tx) Rollback() error { return t.close(true) }

func (t *tx) close(rollback bool) error {
	c := t.conn

	c.lock()
	defer c.unlock()

	if t.closed {
		return nil
	}
	t.closed = true

	if c.dbConn.isBad() {
		return driver.ErrBadConn
	}

	c.inTx = false

	hdbDriver.addTx(-1) // decrement number of transactions.

	if rollback {
		return c.session.Rollback()
	}
	return c.session.Commit()
}

/*
statements

args interface to session
. []interface{} (args) instead of []driver.NamedValue (nvargs) is used as
  . bulk / many operations would have a huge allocation effort / overhead
    converting args to nvargs
  . drawback: nvargs for simply query and exec stmts need to convert nvargs to args

nvargs
. check support (v1.0.0)
  . call (most probably as HANA does support parameter names)
  . query input parameters (most probably not, as HANA does not support them)
  . exec input parameters (could be done (map to table field name) but is it worth the effort?
*/

// TODO handling of nvargs when real named args are supported (v1.0.0)

//  check if statements implements all required interfaces
var (
	_ driver.Stmt              = (*stmt)(nil)
	_ driver.StmtExecContext   = (*stmt)(nil)
	_ driver.StmtQueryContext  = (*stmt)(nil)
	_ driver.NamedValueChecker = (*stmt)(nil)

	_ driver.Stmt              = (*callStmt)(nil)
	_ driver.StmtExecContext   = (*callStmt)(nil)
	_ driver.StmtQueryContext  = (*callStmt)(nil)
	_ driver.NamedValueChecker = (*callStmt)(nil)
)

type argsPool struct {
	sync.Pool
}

func (ap *argsPool) put(v []interface{}) { ap.Put(v) }

func (ap *argsPool) getSize(size int) []interface{} {
	v := ap.Get()
	if v == nil || cap(v.([]interface{})) < size {
		return make([]interface{}, size)
	}
	return v.([]interface{})[0:size]
}

func (ap *argsPool) getNVArgs(nvargs []driver.NamedValue) []interface{} {
	v := ap.getSize(len(nvargs))
	for i, nv := range nvargs {
		v[i] = nv.Value
	}
	return v
}

var smallArgsPool = argsPool{} // rather small slices

type stmt struct {
	conn              *Conn
	query             string
	pr                *p.PrepareResult
	bulk, flush, many bool
	bulkSize, numBulk int
	trace             bool          // store flag for performance reasons (especially bulk inserts)
	args              []interface{} // bulk or many
}

func newStmt(conn *Conn, query string, bulk bool, bulkSize int, pr *p.PrepareResult) *stmt {
	return &stmt{conn: conn, query: query, pr: pr, bulk: bulk, bulkSize: bulkSize, trace: sqltrace.On()}
}

type callStmt struct {
	conn  *Conn
	query string
	pr    *p.PrepareResult
}

func newCallStmt(conn *Conn, query string, pr *p.PrepareResult) *callStmt {
	return &callStmt{conn: conn, query: query, pr: pr}
}

/*
	NumInput differs dependent on statement (check is done in QueryContext and ExecContext):
	- #args == #param (only in params):    query, exec, exec bulk (non control query)
	- #args == #param (in and out params): exec call
	- #args == 0:                          exec bulk (control query)
	- #args == #input param:               query call
*/
func (s *stmt) NumInput() int     { return -1 }
func (s *callStmt) NumInput() int { return -1 }

// stmt methods

/*
reset args
- keep slice to avoid additional allocations but
- free elements (GC)
*/
func (s *stmt) resetArgs() {
	for i := 0; i < len(s.args); i++ {
		s.args[i] = nil
	}
	s.args = s.args[:0]
}

func (s *stmt) Close() error {
	c := s.conn

	c.lock()
	defer c.unlock()

	if c.dbConn.isBad() {
		return driver.ErrBadConn
	}

	hdbDriver.addStmt(-1) // decrement number of statements.

	if s.args != nil {
		if len(s.args) != 0 { // log always //TODO: Fatal?
			dlog.Printf("close: %s - not flushed records: %d)", s.query, len(s.args)/s.pr.NumField())
		}
		s.args = nil
	}

	return c.session.DropStatementID(s.pr.StmtID())
}

func (s *stmt) QueryContext(ctx context.Context, nvargs []driver.NamedValue) (rows driver.Rows, err error) {
	c := s.conn

	if err := c.tryLock(lrNestedQuery); err != nil {
		return nil, err
	}
	hasRowsCloser := false
	defer func() {
		// unlock connection if rows will not do it
		if !hasRowsCloser {
			c.unlock()
		}
	}()

	if c.dbConn.isBad() {
		return nil, driver.ErrBadConn
	}

	args := smallArgsPool.getNVArgs(nvargs)
	defer smallArgsPool.put(args)

	if s.trace {
		sqltrace.Tracef("%s %v", s.query, args)
	}

	if len(args) != s.pr.NumField() { // all fields needs to be input fields
		return nil, fmt.Errorf("invalid number of arguments %d - %d expected", len(args), s.pr.NumField())
	}

	done := make(chan struct{})
	go func() {
		rows, err = c.session.Query(s.pr, args, !c.inTx)
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.dbConn.cancel()
		return nil, ctx.Err()
	case <-done:
		if onCloser, ok := rows.(p.OnCloser); ok {
			onCloser.SetOnClose(c.unlock)
			hasRowsCloser = true
		}
		return rows, err
	}
}

func (s *stmt) ExecContext(ctx context.Context, nvargs []driver.NamedValue) (driver.Result, error) {
	numArg := len(nvargs)
	switch {
	case s.bulk:
		flush := s.flush
		s.flush = false
		if numArg != 0 && numArg != s.pr.NumField() {
			return nil, fmt.Errorf("invalid number of arguments %d - %d expected", numArg, s.pr.NumField())
		}
		return s.execBulk(ctx, nvargs, flush)
	case s.many:
		s.many = false
		if numArg != 1 {
			return nil, fmt.Errorf("invalid argument of arguments %d when using composite arguments - 1 expected", numArg)
		}
		return s.execMany(ctx, &nvargs[0])
	default:
		if numArg != s.pr.NumField() {
			return nil, fmt.Errorf("invalid number of arguments %d - %d expected", numArg, s.pr.NumField())
		}
		args := smallArgsPool.getNVArgs(nvargs)
		defer smallArgsPool.put(args)
		return s.exec(ctx, args)
	}
}

func (s *stmt) exec(ctx context.Context, args []interface{}) (r driver.Result, err error) {
	c := s.conn

	if err := c.tryLock(0); err != nil {
		return nil, err
	}
	defer c.unlock()

	if c.dbConn.isBad() {
		return nil, driver.ErrBadConn
	}

	if connHook != nil {
		connHook(c, choStmtExec)
	}

	if s.trace {
		if len(args) > maxNumTraceArg {
			sqltrace.Tracef("%s first %d arguments: %v", s.query, maxNumTraceArg, args[:maxNumTraceArg])
		} else {
			sqltrace.Tracef("%s %v", s.query, args)
		}
	}

	done := make(chan struct{})
	go func() {
		r, err = c.session.Exec(s.pr, args, !c.inTx)
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.dbConn.cancel()
		return nil, ctx.Err()
	case <-done:
		return r, err
	}
}

func (s *stmt) execBulk(ctx context.Context, nvargs []driver.NamedValue, flush bool) (r driver.Result, err error) {
	numArg := len(nvargs)

	switch numArg {
	case 0: // exec without args --> flush
		flush = true
	default: // add to argument buffer
		if s.args == nil {
			s.args = make([]interface{}, 0, s.pr.NumField()*s.bulkSize)
		}
		for _, nv := range nvargs {
			s.args = append(s.args, nv.Value)
		}
		s.numBulk++
		if s.numBulk >= s.bulkSize {
			flush = true
		}
	}

	if !flush || s.numBulk == 0 { // done: no flush
		return driver.ResultNoRows, nil
	}

	// flush
	r, err = s.exec(ctx, s.args)
	s.resetArgs()
	s.numBulk = 0
	return
}

/*
execMany variants
*/

type execManyer interface {
	numRow() int
	fill(pr *p.PrepareResult, startRow, endRow int, args []interface{}) error
}

type execManyIntfList []interface{}
type execManyIntfMatrix [][]interface{}
type execManyGenList reflect.Value
type execManyGenMatrix reflect.Value

func (em execManyIntfList) numRow() int   { return len(em) }
func (em execManyIntfMatrix) numRow() int { return len(em) }
func (em execManyGenList) numRow() int    { return reflect.Value(em).Len() }
func (em execManyGenMatrix) numRow() int  { return reflect.Value(em).Len() }

func (em execManyIntfList) fill(pr *p.PrepareResult, startRow, endRow int, args []interface{}) error {
	f := pr.ParameterField(0)
	rows := em[startRow:endRow]
	for i, row := range rows {
		row, err := convertValue(f, row)
		if err != nil {
			return err
		}
		args[i] = row
	}
	return nil
}

func (em execManyGenList) fill(pr *p.PrepareResult, startRow, endRow int, args []interface{}) error {
	f := pr.ParameterField(0)
	cnt := 0
	for i := startRow; i < endRow; i++ {
		row, err := convertValue(f, reflect.Value(em).Index(i).Interface())
		if err != nil {
			return err
		}
		args[cnt] = row
		cnt++
	}
	return nil
}

func (em execManyIntfMatrix) fill(pr *p.PrepareResult, startRow, endRow int, args []interface{}) error {
	numField := pr.NumField()
	rows := em[startRow:endRow]
	cnt := 0
	for i, row := range rows {
		if len(row) != numField {
			return fmt.Errorf("invalid number of fields in row %d - got %d - expected %d", i, len(row), numField)
		}
		for j, col := range row {
			f := pr.ParameterField(j)
			col, err := convertValue(f, col)
			if err != nil {
				return err
			}
			args[cnt] = col
			cnt++
		}
	}
	return nil
}

func (em execManyGenMatrix) fill(pr *p.PrepareResult, startRow, endRow int, args []interface{}) error {
	numField := pr.NumField()
	cnt := 0
	for i := startRow; i < endRow; i++ {
		v, err := convertMany(reflect.Value(em).Index(i).Interface())
		if err != nil {
			return err
		}
		row := reflect.ValueOf(v) // need to be array or slice
		if row.Len() != numField {
			return fmt.Errorf("invalid number of fields in row %d - got %d - expected %d", i, row.Len(), numField)
		}
		for j := 0; j < numField; j++ {
			col := row.Index(j).Interface()
			f := pr.ParameterField(j)
			col, err := convertValue(f, col)
			if err != nil {
				return err
			}
			args[cnt] = col
			cnt++
		}
	}
	return nil
}

func (s *stmt) newExecManyVariant(numField int, v interface{}) execManyer {
	if numField == 1 {
		if v, ok := v.([]interface{}); ok {
			return execManyIntfList(v)
		}
		return execManyGenList(reflect.ValueOf(v))
	}
	if v, ok := v.([][]interface{}); ok {
		return execManyIntfMatrix(v)
	}
	return execManyGenMatrix(reflect.ValueOf(v))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/*
Non 'atomic' (transactional) operation due to the split in packages (maxBulkSize),
execMany data might only be written partially to the database in case of hdb stmt errors.
*/
func (s *stmt) execMany(ctx context.Context, nvarg *driver.NamedValue) (driver.Result, error) {

	if len(s.args) != 0 {
		return driver.ResultNoRows, fmt.Errorf("execMany: not flushed entries: %d)", len(s.args))
	}

	numField := s.pr.NumField()

	defer func() { s.resetArgs() }() // reset args

	var totalRowsAffected int64

	variant := s.newExecManyVariant(numField, nvarg.Value)
	numRow := variant.numRow()

	size := min(numRow*numField, s.bulkSize*numField)
	if s.args == nil || cap(s.args) < size {
		s.args = make([]interface{}, size)
	} else {
		s.args = s.args[:size]
	}

	numPack := numRow / s.bulkSize
	if numRow%s.bulkSize != 0 {
		numPack++
	}

	for p := 0; p < numPack; p++ {

		startRow := p * s.bulkSize
		endRow := min(startRow+s.bulkSize, numRow)

		args := s.args[0 : (endRow-startRow)*numField]

		if err := variant.fill(s.pr, startRow, endRow, args); err != nil {
			return driver.RowsAffected(totalRowsAffected), err
		}

		// flush
		r, err := s.exec(ctx, args)
		if err != nil {
			return driver.RowsAffected(totalRowsAffected), err
		}
		n, err := r.RowsAffected()
		totalRowsAffected += n
		if err != nil {
			return driver.RowsAffected(totalRowsAffected), err
		}
	}

	return driver.RowsAffected(totalRowsAffected), nil
}

// CheckNamedValue implements NamedValueChecker interface.
func (s *stmt) CheckNamedValue(nv *driver.NamedValue) error {
	// check on bulk args
	if nv.Name == bulk {
		if ptr, ok := nv.Value.(**struct{}); ok {
			switch ptr {
			case &noFlushTok:
				s.bulk = true
				return driver.ErrRemoveArgument
			case &flushTok:
				s.flush = true
				return driver.ErrRemoveArgument
			}
		}
	}

	// check on standard value
	err := convertNamedValue(s.pr, nv)
	if err == nil || s.bulk || nv.Ordinal != 1 {
		return err
	}

	// check first argument if 'composite'
	if nv.Value, err = convertMany(nv.Value); err != nil {
		return err
	}

	s.many = true
	return nil

}

// callStmt methods

func (s *callStmt) Close() error {
	c := s.conn

	c.lock()
	defer c.unlock()

	if c.dbConn.isBad() {
		return driver.ErrBadConn
	}

	hdbDriver.addStmt(-1) // decrement number of statements.

	return c.session.DropStatementID(s.pr.StmtID())
}

func (s *callStmt) QueryContext(ctx context.Context, nvargs []driver.NamedValue) (rows driver.Rows, err error) {
	c := s.conn

	if err := c.tryLock(lrNestedQuery); err != nil {
		return nil, err
	}
	hasRowsCloser := false
	defer func() {
		// unlock connection if rows will not do it
		if !hasRowsCloser {
			c.unlock()
		}
	}()

	if c.dbConn.isBad() {
		return nil, driver.ErrBadConn
	}

	args := smallArgsPool.getNVArgs(nvargs)
	defer smallArgsPool.put(args)

	if sqltrace.On() {
		sqltrace.Tracef("%s %v", s.query, args)
	}

	if len(args) != s.pr.NumInputField() { // input fields only
		return nil, fmt.Errorf("invalid number of arguments %d - %d expected", len(args), s.pr.NumInputField())
	}

	done := make(chan struct{})
	go func() {
		rows, err = c.session.QueryCall(s.pr, args)
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.dbConn.cancel()
		return nil, ctx.Err()
	case <-done:
		if onCloser, ok := rows.(p.OnCloser); ok {
			onCloser.SetOnClose(c.unlock)
			hasRowsCloser = true
		}
		return rows, err
	}
}

func (s *callStmt) ExecContext(ctx context.Context, nvargs []driver.NamedValue) (r driver.Result, err error) {
	c := s.conn

	if err := c.tryLock(0); err != nil {
		return nil, err
	}
	defer c.unlock()

	if c.dbConn.isBad() {
		return nil, driver.ErrBadConn
	}

	args := smallArgsPool.getNVArgs(nvargs)
	defer smallArgsPool.put(args)

	if sqltrace.On() {
		sqltrace.Tracef("%s %v", s.query, args)
	}

	if len(args) != s.pr.NumField() {
		return nil, fmt.Errorf("invalid number of arguments %d - %d expected", len(args), s.pr.NumField())
	}

	done := make(chan struct{})
	go func() {
		r, err = c.session.ExecCall(s.pr, args)
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.dbConn.cancel()
		return nil, ctx.Err()
	case <-done:
		return r, err
	}
}

// CheckNamedValue implements NamedValueChecker interface.
func (s *callStmt) CheckNamedValue(nv *driver.NamedValue) error {
	return convertNamedValue(s.pr, nv)
}
