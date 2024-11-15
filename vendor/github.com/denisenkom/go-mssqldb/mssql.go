package mssql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/denisenkom/go-mssqldb/internal/querytext"
	"github.com/denisenkom/go-mssqldb/msdsn"
	"github.com/golang-sql/sqlexp"
)

// ReturnStatus may be used to return the return value from a proc.
//
//   var rs mssql.ReturnStatus
//   _, err := db.Exec("theproc", &rs)
//   log.Printf("return status = %d", rs)
type ReturnStatus int32

var driverInstance = &Driver{processQueryText: true}
var driverInstanceNoProcess = &Driver{processQueryText: false}

func init() {
	sql.Register("mssql", driverInstance)
	sql.Register("sqlserver", driverInstanceNoProcess)
	createDialer = func(p *msdsn.Config) Dialer {
		ka := p.KeepAlive
		if ka == 0 {
			ka = 30 * time.Second
		}
		return netDialer{&net.Dialer{KeepAlive: ka}}
	}
}

var createDialer func(p *msdsn.Config) Dialer

type netDialer struct {
	nd *net.Dialer
}

func (d netDialer) DialContext(ctx context.Context, network string, addr string) (net.Conn, error) {
	return d.nd.DialContext(ctx, network, addr)
}

type Driver struct {
	logger optionalLogger

	processQueryText bool
}

// OpenConnector opens a new connector. Useful to dial with a context.
func (d *Driver) OpenConnector(dsn string) (*Connector, error) {
	params, _, err := msdsn.Parse(dsn)
	if err != nil {
		return nil, err
	}

	return &Connector{
		params: params,
		driver: d,
	}, nil
}

func (d *Driver) Open(dsn string) (driver.Conn, error) {
	return d.open(context.Background(), dsn)
}

// SetLogger sets a Logger for both driver instances ("mssql" and "sqlserver").
// Use this to have go-msqldb log additional information in a format it picks.
// You can set either a Logger or a ContextLogger, but not both. Calling SetLogger
// will overwrite any ContextLogger you set with SetContextLogger.
func SetLogger(logger Logger) {
	driverInstance.SetLogger(logger)
	driverInstanceNoProcess.SetLogger(logger)
}

// SetLogger sets a Logger for the driver instance on which you call it.
// Use this to have go-msqldb log additional information in a format it picks.
// You can set either a Logger or a ContextLogger, but not both. Calling SetLogger
// will overwrite any ContextLogger you set with SetContextLogger.
func (d *Driver) SetLogger(logger Logger) {
	d.logger = optionalLogger{loggerAdapter{logger}}
}

// SetContextLogger sets a ContextLogger for both driver instances ("mssql" and "sqlserver").
// Use this to get callbacks from go-mssqldb with additional information and extra details
// that you can log in the format of your choice.
// You can set either a ContextLogger or a Logger, but not both. Calling SetContextLogger
// will overwrite any Logger you set with SetLogger.
func SetContextLogger(ctxLogger ContextLogger) {
	driverInstance.SetContextLogger(ctxLogger)
	driverInstanceNoProcess.SetContextLogger(ctxLogger)
}

// SetContextLogger sets a ContextLogger for the driver instance on which you call it.
// Use this to get callbacks from go-mssqldb with additional information and extra details
// that you can log in the format of your choice.
// You can set either a ContextLogger or a Logger, but not both. Calling SetContextLogger
// will overwrite any Logger you set with SetLogger.
func (d *Driver) SetContextLogger(ctxLogger ContextLogger) {
	d.logger = optionalLogger{ctxLogger}
}

// NewConnector creates a new connector from a DSN.
// The returned connector may be used with sql.OpenDB.
func NewConnector(dsn string) (*Connector, error) {
	params, _, err := msdsn.Parse(dsn)
	if err != nil {
		return nil, err
	}
	c := &Connector{
		params: params,
		driver: driverInstanceNoProcess,
	}
	return c, nil
}

// NewConnectorConfig creates a new Connector for a DSN Config struct.
// The returned connector may be used with sql.OpenDB.
func NewConnectorConfig(config msdsn.Config) *Connector {
	return &Connector{
		params: config,
		driver: driverInstanceNoProcess,
	}
}

// Connector holds the parsed DSN and is ready to make a new connection
// at any time.
//
// In the future, settings that cannot be passed through a string DSN
// may be set directly on the connector.
type Connector struct {
	params msdsn.Config
	driver *Driver

	fedAuthRequired     bool
	fedAuthLibrary      int
	fedAuthADALWorkflow byte

	// callback that can provide a security token during login
	securityTokenProvider func(ctx context.Context) (string, error)

	// callback that can provide a security token during ADAL login
	adalTokenProvider func(ctx context.Context, serverSPN, stsURL string) (string, error)

	// SessionInitSQL is executed after marking a given session to be reset.
	// When not present, the next query will still reset the session to the
	// database defaults.
	//
	// When present the connection will immediately mark the session to
	// be reset, then execute the SessionInitSQL text to setup the session
	// that may be different from the base database defaults.
	//
	// For Example, the application relies on the following defaults
	// but is not allowed to set them at the database system level.
	//
	//    SET XACT_ABORT ON;
	//    SET TEXTSIZE -1;
	//    SET ANSI_NULLS ON;
	//    SET LOCK_TIMEOUT 10000;
	//
	// SessionInitSQL should not attempt to manually call sp_reset_connection.
	// This will happen at the TDS layer.
	//
	// SessionInitSQL is optional. The session will be reset even if
	// SessionInitSQL is empty.
	SessionInitSQL string

	// Dialer sets a custom dialer for all network operations.
	// If Dialer is not set, normal net dialers are used.
	Dialer Dialer
}

type Dialer interface {
	DialContext(ctx context.Context, network string, addr string) (net.Conn, error)
}

func (c *Connector) getDialer(p *msdsn.Config) Dialer {
	if c != nil && c.Dialer != nil {
		return c.Dialer
	}
	return createDialer(p)
}

type Conn struct {
	connector      *Connector
	sess           *tdsSession
	transactionCtx context.Context
	resetSession   bool

	processQueryText bool
	connectionGood   bool

	outs outputs
}

type outputs struct {
	params       map[string]interface{}
	returnStatus *ReturnStatus
	msgq         *sqlexp.ReturnMessage
}

// IsValid satisfies the driver.Validator interface.
func (c *Conn) IsValid() bool {
	return c.connectionGood
}

// checkBadConn marks the connection as bad based on the characteristics
// of the supplied error. Bad connections will be dropped from the connection
// pool rather than reused.
//
// If bad connection retry is enabled and the error + connection state permits
// retrying, checkBadConn will return a RetryableError that allows database/sql
// to automatically retry the query with another connection.
func (c *Conn) checkBadConn(ctx context.Context, err error, mayRetry bool) error {
	switch err {
	case nil:
		return nil
	case io.EOF:
		c.connectionGood = false
	case driver.ErrBadConn:
		// It is an internal programming error if driver.ErrBadConn
		// is ever passed to this function. driver.ErrBadConn should
		// only ever be returned in response to a *mssql.Conn.connectionGood == false
		// check in the external facing API.
		panic("driver.ErrBadConn in checkBadConn. This should not happen.")
	}

	switch err.(type) {
	case net.Error:
		c.connectionGood = false
	case StreamError:
		c.connectionGood = false
	case ServerError:
		c.connectionGood = false
	}

	if !c.connectionGood && mayRetry && !c.connector.params.DisableRetry {
		if c.sess.logFlags&logRetries != 0 {
			c.sess.logger.Log(ctx, msdsn.LogRetries, err.Error())
		}
		return newRetryableError(err)
	}

	return err
}

func (c *Conn) clearOuts() {
	c.outs = outputs{}
}

func (c *Conn) simpleProcessResp(ctx context.Context) error {
	reader := startReading(c.sess, ctx, c.outs)
	c.clearOuts()

	var resultError error
	err := reader.iterateResponse()
	if err != nil {
		return c.checkBadConn(ctx, err, false)
	}
	return resultError
}

func (c *Conn) Commit() error {
	if !c.connectionGood {
		return driver.ErrBadConn
	}
	if err := c.sendCommitRequest(); err != nil {
		return c.checkBadConn(c.transactionCtx, err, true)
	}
	return c.simpleProcessResp(c.transactionCtx)
}

func (c *Conn) sendCommitRequest() error {
	headers := []headerStruct{
		{hdrtype: dataStmHdrTransDescr,
			data: transDescrHdr{c.sess.tranid, 1}.pack()},
	}
	reset := c.resetSession
	c.resetSession = false
	if err := sendCommitXact(c.sess.buf, headers, "", 0, 0, "", reset); err != nil {
		if c.sess.logFlags&logErrors != 0 {
			c.sess.logger.Log(c.transactionCtx, msdsn.LogErrors, fmt.Sprintf("Failed to send CommitXact with %v", err))
		}
		c.connectionGood = false
		return fmt.Errorf("faild to send CommitXact: %v", err)
	}
	return nil
}

func (c *Conn) Rollback() error {
	if !c.connectionGood {
		return driver.ErrBadConn
	}
	if err := c.sendRollbackRequest(); err != nil {
		return c.checkBadConn(c.transactionCtx, err, true)
	}
	return c.simpleProcessResp(c.transactionCtx)
}

func (c *Conn) sendRollbackRequest() error {
	headers := []headerStruct{
		{hdrtype: dataStmHdrTransDescr,
			data: transDescrHdr{c.sess.tranid, 1}.pack()},
	}
	reset := c.resetSession
	c.resetSession = false
	if err := sendRollbackXact(c.sess.buf, headers, "", 0, 0, "", reset); err != nil {
		if c.sess.logFlags&logErrors != 0 {
			c.sess.logger.Log(c.transactionCtx, msdsn.LogErrors, fmt.Sprintf("Failed to send RollbackXact with %v", err))
		}
		c.connectionGood = false
		return fmt.Errorf("failed to send RollbackXact: %v", err)
	}
	return nil
}

func (c *Conn) Begin() (driver.Tx, error) {
	return c.begin(context.Background(), isolationUseCurrent)
}

func (c *Conn) begin(ctx context.Context, tdsIsolation isoLevel) (tx driver.Tx, err error) {
	if !c.connectionGood {
		return nil, driver.ErrBadConn
	}
	err = c.sendBeginRequest(ctx, tdsIsolation)
	if err != nil {
		return nil, c.checkBadConn(ctx, err, true)
	}
	tx, err = c.processBeginResponse(ctx)
	if err != nil {
		return nil, err
	}
	return
}

func (c *Conn) sendBeginRequest(ctx context.Context, tdsIsolation isoLevel) error {
	c.transactionCtx = ctx
	headers := []headerStruct{
		{hdrtype: dataStmHdrTransDescr,
			data: transDescrHdr{0, 1}.pack()},
	}
	reset := c.resetSession
	c.resetSession = false
	if err := sendBeginXact(c.sess.buf, headers, tdsIsolation, "", reset); err != nil {
		if c.sess.logFlags&logErrors != 0 {
			c.sess.logger.Log(ctx, msdsn.LogErrors, fmt.Sprintf("Failed to send BeginXact with %v", err))
		}
		c.connectionGood = false
		return fmt.Errorf("failed to send BeginXact: %v", err)
	}
	return nil
}

func (c *Conn) processBeginResponse(ctx context.Context) (driver.Tx, error) {
	if err := c.simpleProcessResp(ctx); err != nil {
		return nil, err
	}
	// successful BEGINXACT request will return sess.tranid
	// for started transaction
	return c, nil
}

func (d *Driver) open(ctx context.Context, dsn string) (*Conn, error) {
	params, _, err := msdsn.Parse(dsn)
	if err != nil {
		return nil, err
	}
	c := &Connector{params: params}
	return d.connect(ctx, c, params)
}

// connect to the server, using the provided context for dialing only.
func (d *Driver) connect(ctx context.Context, c *Connector, params msdsn.Config) (*Conn, error) {
	sess, err := connect(ctx, c, d.logger, params)
	if err != nil {
		// main server failed, try fail-over partner
		if params.FailOverPartner == "" {
			return nil, err
		}

		params.Host = params.FailOverPartner
		if params.FailOverPort != 0 {
			params.Port = params.FailOverPort
		}

		sess, err = connect(ctx, c, d.logger, params)
		if err != nil {
			// fail-over partner also failed, now fail
			return nil, err
		}
	}

	conn := &Conn{
		connector:        c,
		sess:             sess,
		transactionCtx:   context.Background(),
		processQueryText: d.processQueryText,
		connectionGood:   true,
	}

	return conn, nil
}

func (c *Conn) Close() error {
	return c.sess.buf.transport.Close()
}

type Stmt struct {
	c          *Conn
	query      string
	paramCount int
	notifSub   *queryNotifSub
}

type queryNotifSub struct {
	msgText string
	options string
	timeout uint32
}

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	if !c.connectionGood {
		return nil, driver.ErrBadConn
	}
	if len(query) > 10 && strings.EqualFold(query[:10], "INSERTBULK") {
		return c.prepareCopyIn(context.Background(), query)
	}
	return c.prepareContext(context.Background(), query)
}

func (c *Conn) prepareContext(ctx context.Context, query string) (*Stmt, error) {
	paramCount := -1
	if c.processQueryText {
		query, paramCount = querytext.ParseParams(query)
	}
	return &Stmt{c, query, paramCount, nil}, nil
}

func (s *Stmt) Close() error {
	return nil
}

func (s *Stmt) SetQueryNotification(id, options string, timeout time.Duration) {
	// 2.2.5.3.1 Query Notifications Header
	// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-tds/e168d373-a7b7-41aa-b6ca-25985466a7e0
	// Timeout in milliseconds in TDS protocol.
	to := uint32(timeout / time.Millisecond)
	if to < 1 {
		to = 1
	}
	s.notifSub = &queryNotifSub{id, options, to}
}

func (s *Stmt) NumInput() int {
	return s.paramCount
}

func (s *Stmt) sendQuery(ctx context.Context, args []namedValue) (err error) {
	headers := []headerStruct{
		{hdrtype: dataStmHdrTransDescr,
			data: transDescrHdr{s.c.sess.tranid, 1}.pack()},
	}

	if s.notifSub != nil {
		headers = append(headers,
			headerStruct{
				hdrtype: dataStmHdrQueryNotif,
				data: queryNotifHdr{
					s.notifSub.msgText,
					s.notifSub.options,
					s.notifSub.timeout,
				}.pack(),
			})
	}

	conn := s.c

	// no need to check number of parameters here, it is checked by database/sql
	if conn.sess.logFlags&logSQL != 0 {
		conn.sess.logger.Log(ctx, msdsn.LogSQL, s.query)
	}
	if conn.sess.logFlags&logParams != 0 && len(args) > 0 {
		for i := 0; i < len(args); i++ {
			if len(args[i].Name) > 0 {
				s.c.sess.logger.Log(ctx, msdsn.LogParams, fmt.Sprintf("\t@%s\t%v", args[i].Name, args[i].Value))
			} else {
				s.c.sess.logger.Log(ctx, msdsn.LogParams, fmt.Sprintf("\t@p%d\t%v", i+1, args[i].Value))
			}
		}
	}

	reset := conn.resetSession
	conn.resetSession = false
	isProc := isProc(s.query)
	if len(args) == 0 && !isProc {
		if err = sendSqlBatch72(conn.sess.buf, s.query, headers, reset); err != nil {
			if conn.sess.logFlags&logErrors != 0 {
				conn.sess.logger.Log(ctx, msdsn.LogErrors, fmt.Sprintf("Failed to send SqlBatch with %v", err))
			}
			conn.connectionGood = false
			return fmt.Errorf("failed to send SQL Batch: %v", err)
		}
	} else {
		proc := sp_ExecuteSql
		var params []param
		if isProc {
			proc.name = s.query
			params, _, err = s.makeRPCParams(args, true)
			if err != nil {
				return
			}
		} else {
			var decls []string
			params, decls, err = s.makeRPCParams(args, false)
			if err != nil {
				return
			}
			params[0] = makeStrParam(s.query)
			params[1] = makeStrParam(strings.Join(decls, ","))
		}
		if err = sendRpc(conn.sess.buf, headers, proc, 0, params, reset); err != nil {
			if conn.sess.logFlags&logErrors != 0 {
				conn.sess.logger.Log(ctx, msdsn.LogErrors, fmt.Sprintf("Failed to send Rpc with %v", err))
			}
			conn.connectionGood = false
			return fmt.Errorf("failed to send RPC: %v", err)
		}
	}
	return
}

// isProc takes the query text in s and determines if it is a stored proc name
// or SQL text.
func isProc(s string) bool {
	if len(s) == 0 {
		return false
	}
	const (
		outside = iota
		text
		escaped
	)
	st := outside
	var rn1, rPrev rune
	for _, r := range s {
		rPrev = rn1
		rn1 = r
		if st != escaped {
			switch r {
			// No newlines or string sequences.
			case '\n', '\r', '\'', ';':
				return false
			}
		}
		switch st {
		case outside:
			switch {
			case r == '[':
				st = escaped
			case r == ']' && rPrev == ']':
				st = escaped
			case unicode.IsLetter(r):
				st = text
			case r == '_':
				st = text
			case r == '#':
				st = text
			case r == '.':
			default:
				return false
			}
		case text:
			switch {
			case r == '.':
				st = outside
			case r == '[':
				return false
			case r == '(':
				return false
			case unicode.IsSpace(r):
				return false
			}
		case escaped:
			switch {
			case r == ']':
				st = outside
			}
		}
	}
	return true
}

func (s *Stmt) makeRPCParams(args []namedValue, isProc bool) ([]param, []string, error) {
	var err error
	var offset int
	if !isProc {
		offset = 2
	}
	params := make([]param, len(args)+offset)
	decls := make([]string, len(args))
	for i, val := range args {
		params[i+offset], err = s.makeParam(val.Value)
		if err != nil {
			return nil, nil, err
		}
		var name string
		if len(val.Name) > 0 {
			name = "@" + val.Name
		} else if !isProc {
			name = fmt.Sprintf("@p%d", val.Ordinal)
		}
		params[i+offset].Name = name
		const outputSuffix = " output"
		var output string
		if isOutputValue(val.Value) {
			output = outputSuffix
		}
		decls[i] = fmt.Sprintf("%s %s%s", name, makeDecl(params[i+offset].ti), output)

	}
	return params, decls, nil
}

type namedValue struct {
	Name    string
	Ordinal int
	Value   driver.Value
}

func convertOldArgs(args []driver.Value) []namedValue {
	list := make([]namedValue, len(args))
	for i, v := range args {
		list[i] = namedValue{
			Ordinal: i + 1,
			Value:   v,
		}
	}
	return list
}

func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	defer s.c.clearOuts()

	return s.queryContext(context.Background(), convertOldArgs(args))
}

func (s *Stmt) queryContext(ctx context.Context, args []namedValue) (rows driver.Rows, err error) {
	if !s.c.connectionGood {
		return nil, driver.ErrBadConn
	}
	if err = s.sendQuery(ctx, args); err != nil {
		return nil, s.c.checkBadConn(ctx, err, true)
	}
	return s.processQueryResponse(ctx)
}

func (s *Stmt) processQueryResponse(ctx context.Context) (res driver.Rows, err error) {
	ctx, cancel := context.WithCancel(ctx)
	reader := startReading(s.c.sess, ctx, s.c.outs)
	s.c.clearOuts()
	// For apps using a message queue, return right away and let Rowsq do all the work
	if reader.outs.msgq != nil {
		res = &Rowsq{stmt: s, reader: reader, cols: nil, cancel: cancel}
		return res, nil
	}
	// process metadata
	var cols []columnStruct
loop:
	for {
		tok, err := reader.nextToken()
		if err == nil {
			if tok == nil {
				break
			} else {
				switch token := tok.(type) {
				// By ignoring DONE token we effectively
				// skip empty result-sets.
				// This improves results in queries like that:
				// set nocount on; select 1
				// see TestIgnoreEmptyResults test
				//case doneStruct:
				//break loop
				case []columnStruct:
					cols = token
					break loop
				case doneStruct:
					if token.isError() {
						// need to cleanup cancellable context
						cancel()
						return nil, s.c.checkBadConn(ctx, token.getError(), false)
					}
				case ReturnStatus:
					if reader.outs.returnStatus != nil {
						*reader.outs.returnStatus = token
					}
				}
			}
		} else {
			// need to cleanup cancellable context
			cancel()
			return nil, s.c.checkBadConn(ctx, err, false)
		}
	}
	res = &Rows{stmt: s, reader: reader, cols: cols, cancel: cancel}
	return
}

func (s *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	defer s.c.clearOuts()

	return s.exec(context.Background(), convertOldArgs(args))
}

func (s *Stmt) exec(ctx context.Context, args []namedValue) (res driver.Result, err error) {
	if !s.c.connectionGood {
		return nil, driver.ErrBadConn
	}
	if err = s.sendQuery(ctx, args); err != nil {
		return nil, s.c.checkBadConn(ctx, err, true)
	}
	if res, err = s.processExec(ctx); err != nil {
		return nil, err
	}
	return
}

func (s *Stmt) processExec(ctx context.Context) (res driver.Result, err error) {
	reader := startReading(s.c.sess, ctx, s.c.outs)
	s.c.clearOuts()
	err = reader.iterateResponse()
	if err != nil {
		return nil, s.c.checkBadConn(ctx, err, false)
	}
	return &Result{s.c, reader.rowCount}, nil
}

// Rows represents the non-experimental data/sql model for Query and QueryContext
type Rows struct {
	stmt     *Stmt
	cols     []columnStruct
	reader   *tokenProcessor
	nextCols []columnStruct
	cancel   func()
}

func (rc *Rows) Close() error {
	// need to add a test which returns lots of rows
	// and check closing after reading only few rows
	rc.cancel()

	for {
		tok, err := rc.reader.nextToken()
		if err == nil {
			if tok == nil {
				return nil
			} else {
				// continue consuming tokens
				continue
			}
		} else {
			if err == rc.reader.ctx.Err() {
				return nil
			} else {
				return err
			}
		}
	}
}

func (rc *Rows) Columns() (res []string) {

	res = make([]string, len(rc.cols))
	for i, col := range rc.cols {
		res[i] = col.ColName
	}
	return
}

func (rc *Rows) Next(dest []driver.Value) error {
	if !rc.stmt.c.connectionGood {
		return driver.ErrBadConn
	}
	if rc.nextCols != nil {
		return io.EOF
	}
	for {
		tok, err := rc.reader.nextToken()
		if err == nil {
			if tok == nil {
				return io.EOF
			} else {
				switch tokdata := tok.(type) {
				// processQueryResponse may have delegated all the token reading to us
				case []columnStruct:
					rc.nextCols = tokdata
					return io.EOF
				case []interface{}:
					for i := range dest {
						dest[i] = tokdata[i]
					}
					return nil
				case doneStruct:
					if tokdata.isError() {
						return rc.stmt.c.checkBadConn(rc.reader.ctx, tokdata.getError(), false)
					}
				case ReturnStatus:
					if rc.reader.outs.returnStatus != nil {
						*rc.reader.outs.returnStatus = tokdata
					}
				}
			}

		} else {
			return rc.stmt.c.checkBadConn(rc.reader.ctx, err, false)
		}
	}
}

func (rc *Rows) HasNextResultSet() bool {
	return rc.nextCols != nil
}

func (rc *Rows) NextResultSet() error {
	rc.cols = rc.nextCols
	rc.nextCols = nil
	if rc.cols == nil {
		return io.EOF
	}
	return nil
}

// It should return
// the value type that can be used to scan types into. For example, the database
// column type "bigint" this should return "reflect.TypeOf(int64(0))".
func (r *Rows) ColumnTypeScanType(index int) reflect.Type {
	return makeGoLangScanType(r.cols[index].ti)
}

// RowsColumnTypeDatabaseTypeName may be implemented by Rows. It should return the
// database system type name without the length. Type names should be uppercase.
// Examples of returned types: "VARCHAR", "NVARCHAR", "VARCHAR2", "CHAR", "TEXT",
// "DECIMAL", "SMALLINT", "INT", "BIGINT", "BOOL", "[]BIGINT", "JSONB", "XML",
// "TIMESTAMP".
func (r *Rows) ColumnTypeDatabaseTypeName(index int) string {
	return makeGoLangTypeName(r.cols[index].ti)
}

// RowsColumnTypeLength may be implemented by Rows. It should return the length
// of the column type if the column is a variable length type. If the column is
// not a variable length type ok should return false.
// If length is not limited other than system limits, it should return math.MaxInt64.
// The following are examples of returned values for various types:
//   TEXT          (math.MaxInt64, true)
//   varchar(10)   (10, true)
//   nvarchar(10)  (10, true)
//   decimal       (0, false)
//   int           (0, false)
//   bytea(30)     (30, true)
func (r *Rows) ColumnTypeLength(index int) (int64, bool) {
	return makeGoLangTypeLength(r.cols[index].ti)
}

// It should return
// the precision and scale for decimal types. If not applicable, ok should be false.
// The following are examples of returned values for various types:
//   decimal(38, 4)    (38, 4, true)
//   int               (0, 0, false)
//   decimal           (math.MaxInt64, math.MaxInt64, true)
func (r *Rows) ColumnTypePrecisionScale(index int) (int64, int64, bool) {
	return makeGoLangTypePrecisionScale(r.cols[index].ti)
}

// The nullable value should
// be true if it is known the column may be null, or false if the column is known
// to be not nullable.
// If the column nullability is unknown, ok should be false.
func (r *Rows) ColumnTypeNullable(index int) (nullable, ok bool) {
	nullable = r.cols[index].Flags&colFlagNullable != 0
	ok = true
	return
}

func makeStrParam(val string) (res param) {
	res.ti.TypeId = typeNVarChar
	res.buffer = str2ucs2(val)
	res.ti.Size = len(res.buffer)
	return
}

func (s *Stmt) makeParam(val driver.Value) (res param, err error) {
	if val == nil {
		res.ti.TypeId = typeNull
		res.buffer = nil
		res.ti.Size = 0
		return
	}
	switch val := val.(type) {
	case int64:
		res.ti.TypeId = typeIntN
		res.buffer = make([]byte, 8)
		res.ti.Size = 8
		binary.LittleEndian.PutUint64(res.buffer, uint64(val))
	case sql.NullInt64:
		// only null values should be getting here
		res.ti.TypeId = typeIntN
		res.ti.Size = 8
		res.buffer = []byte{}

	case float64:
		res.ti.TypeId = typeFltN
		res.ti.Size = 8
		res.buffer = make([]byte, 8)
		binary.LittleEndian.PutUint64(res.buffer, math.Float64bits(val))
	case sql.NullFloat64:
		// only null values should be getting here
		res.ti.TypeId = typeFltN
		res.ti.Size = 8
		res.buffer = []byte{}

	case []byte:
		res.ti.TypeId = typeBigVarBin
		res.ti.Size = len(val)
		res.buffer = val
	case string:
		res = makeStrParam(val)
	case sql.NullString:
		// only null values should be getting here
		res.ti.TypeId = typeNVarChar
		res.buffer = nil
		res.ti.Size = 8000
	case bool:
		res.ti.TypeId = typeBitN
		res.ti.Size = 1
		res.buffer = make([]byte, 1)
		if val {
			res.buffer[0] = 1
		}
	case sql.NullBool:
		// only null values should be getting here
		res.ti.TypeId = typeBitN
		res.ti.Size = 1
		res.buffer = []byte{}

	case time.Time:
		if s.c.sess.loginAck.TDSVersion >= verTDS73 {
			res.ti.TypeId = typeDateTimeOffsetN
			res.ti.Scale = 7
			res.buffer = encodeDateTimeOffset(val, int(res.ti.Scale))
			res.ti.Size = len(res.buffer)
		} else {
			res.ti.TypeId = typeDateTimeN
			res.buffer = encodeDateTime(val)
			res.ti.Size = len(res.buffer)
		}
	default:
		return s.makeParamExtra(val)
	}
	return
}

type Result struct {
	c            *Conn
	rowsAffected int64
}

func (r *Result) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

var _ driver.Pinger = &Conn{}

// Ping is used to check if the remote server is available and satisfies the Pinger interface.
func (c *Conn) Ping(ctx context.Context) error {
	if !c.connectionGood {
		return driver.ErrBadConn
	}
	stmt := &Stmt{c, `select 1;`, 0, nil}
	_, err := stmt.ExecContext(ctx, nil)
	return err
}

var _ driver.ConnBeginTx = &Conn{}

func convertIsolationLevel(level sql.IsolationLevel) (isoLevel, error) {
	switch level {
	case sql.LevelDefault:
		return isolationUseCurrent, nil
	case sql.LevelReadUncommitted:
		return isolationReadUncommited, nil
	case sql.LevelReadCommitted:
		return isolationReadCommited, nil
	case sql.LevelWriteCommitted:
		return isolationUseCurrent, errors.New("LevelWriteCommitted isolation level is not supported")
	case sql.LevelRepeatableRead:
		return isolationRepeatableRead, nil
	case sql.LevelSnapshot:
		return isolationSnapshot, nil
	case sql.LevelSerializable:
		return isolationSerializable, nil
	case sql.LevelLinearizable:
		return isolationUseCurrent, errors.New("LevelLinearizable isolation level is not supported")
	default:
		return isolationUseCurrent, errors.New("isolation level is not supported or unknown")
	}
}

// BeginTx satisfies ConnBeginTx.
func (c *Conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if !c.connectionGood {
		return nil, driver.ErrBadConn
	}
	if opts.ReadOnly {
		return nil, errors.New("read-only transactions are not supported")
	}

	tdsIsolation, err := convertIsolationLevel(sql.IsolationLevel(opts.Isolation))
	if err != nil {
		return nil, err
	}
	return c.begin(ctx, tdsIsolation)
}

func (c *Conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if !c.connectionGood {
		return nil, driver.ErrBadConn
	}
	if len(query) > 10 && strings.EqualFold(query[:10], "INSERTBULK") {
		return c.prepareCopyIn(ctx, query)
	}

	return c.prepareContext(ctx, query)
}

func (s *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	defer s.c.clearOuts()

	if !s.c.connectionGood {
		return nil, driver.ErrBadConn
	}
	list := make([]namedValue, len(args))
	for i, nv := range args {
		list[i] = namedValue(nv)
	}
	return s.queryContext(ctx, list)
}

func (s *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	defer s.c.clearOuts()

	if !s.c.connectionGood {
		return nil, driver.ErrBadConn
	}
	list := make([]namedValue, len(args))
	for i, nv := range args {
		list[i] = namedValue(nv)
	}
	return s.exec(ctx, list)
}

// Rowsq implements the sqlexp messages model for Query and QueryContext
// Theory: We could also implement the non-experimental model this way
type Rowsq struct {
	stmt        *Stmt
	cols        []columnStruct
	reader      *tokenProcessor
	nextCols    []columnStruct
	cancel      func()
	requestDone bool
	inResultSet bool
}

func (rc *Rowsq) Close() error {
	rc.cancel()

	for {
		tok, err := rc.reader.nextToken()
		if err == nil {
			if tok == nil {
				return nil
			} else {
				// continue consuming tokens
				continue
			}
		} else {
			if err == rc.reader.ctx.Err() {
				return nil
			} else {
				return err
			}
		}
	}
}

// data/sql calls Columns during the app's call to Next
func (rc *Rowsq) Columns() (res []string) {
	if rc.cols == nil {
	scan:
		for {
			tok, err := rc.reader.nextToken()
			if err == nil {
				if rc.reader.sess.logFlags&logDebug != 0 {
					rc.reader.sess.logger.Log(rc.reader.ctx, msdsn.LogDebug, fmt.Sprintf("Columns() token type:%v", reflect.TypeOf(tok)))
				}
				if tok == nil {
					return []string{}
				} else {
					switch tokdata := tok.(type) {
					case []columnStruct:
						rc.cols = tokdata
						rc.inResultSet = true
						break scan
					}
				}
			}
		}
	}
	res = make([]string, len(rc.cols))
	for i, col := range rc.cols {
		res[i] = col.ColName
	}
	return
}

func (rc *Rowsq) Next(dest []driver.Value) error {
	if !rc.stmt.c.connectionGood {
		return driver.ErrBadConn
	}
	for {
		tok, err := rc.reader.nextToken()
		if rc.reader.sess.logFlags&logDebug != 0 {
			rc.reader.sess.logger.Log(rc.reader.ctx, msdsn.LogDebug, fmt.Sprintf("Next() token type:%v", reflect.TypeOf(tok)))
		}
		if err == nil {
			if tok == nil {
				return io.EOF
			} else {
				switch tokdata := tok.(type) {
				case []interface{}:
					for i := range dest {
						dest[i] = tokdata[i]
					}
					return nil
				case doneStruct:
					if tokdata.Status&doneMore == 0 {
						rc.requestDone = true
					}
					if tokdata.isError() {
						e := rc.stmt.c.checkBadConn(rc.reader.ctx, tokdata.getError(), false)
						switch e.(type) {
						case Error:
							// Ignore non-fatal server errors. Fatal errors are of type ServerError
						default:
							return e
						}
					}
					if rc.inResultSet {
						rc.inResultSet = false
						return io.EOF
					}
				case ReturnStatus:
					if rc.reader.outs.returnStatus != nil {
						*rc.reader.outs.returnStatus = tokdata
					}
				}
			}

		} else {
			return rc.stmt.c.checkBadConn(rc.reader.ctx, err, false)
		}
	}
}

// In Message Queue mode, we always claim another resultset could be on the way
// to avoid Rows being closed prematurely
func (rc *Rowsq) HasNextResultSet() bool {
	return !rc.requestDone
}

// Scans to the next set of columns in the stream
// Note that the caller may not have read all the rows in the prior set
func (rc *Rowsq) NextResultSet() error {
	if rc.requestDone {
		return io.EOF
	}
scan:
	for {
		// we should have a columns token in the channel if we aren't at the end
		tok, err := rc.reader.nextToken()
		if rc.reader.sess.logFlags&logDebug != 0 {
			rc.reader.sess.logger.Log(rc.reader.ctx, msdsn.LogDebug, fmt.Sprintf("NextResultSet() token type:%v", reflect.TypeOf(tok)))
		}

		if err != nil {
			return err
		}
		if tok == nil {
			return io.EOF
		}
		switch tokdata := tok.(type) {
		case []columnStruct:
			rc.nextCols = tokdata
			rc.inResultSet = true
			break scan
		case doneStruct:
			if tokdata.Status&doneMore == 0 {
				rc.nextCols = nil
				rc.requestDone = true
				break scan
			}
		}
	}
	rc.cols = rc.nextCols
	rc.nextCols = nil
	if rc.cols == nil {
		return io.EOF
	}
	return nil
}

// It should return
// the value type that can be used to scan types into. For example, the database
// column type "bigint" this should return "reflect.TypeOf(int64(0))".
func (r *Rowsq) ColumnTypeScanType(index int) reflect.Type {
	return makeGoLangScanType(r.cols[index].ti)
}

// RowsColumnTypeDatabaseTypeName may be implemented by Rows. It should return the
// database system type name without the length. Type names should be uppercase.
// Examples of returned types: "VARCHAR", "NVARCHAR", "VARCHAR2", "CHAR", "TEXT",
// "DECIMAL", "SMALLINT", "INT", "BIGINT", "BOOL", "[]BIGINT", "JSONB", "XML",
// "TIMESTAMP".
func (r *Rowsq) ColumnTypeDatabaseTypeName(index int) string {
	return makeGoLangTypeName(r.cols[index].ti)
}

// RowsColumnTypeLength may be implemented by Rows. It should return the length
// of the column type if the column is a variable length type. If the column is
// not a variable length type ok should return false.
// If length is not limited other than system limits, it should return math.MaxInt64.
// The following are examples of returned values for various types:
//   TEXT          (math.MaxInt64, true)
//   varchar(10)   (10, true)
//   nvarchar(10)  (10, true)
//   decimal       (0, false)
//   int           (0, false)
//   bytea(30)     (30, true)
func (r *Rowsq) ColumnTypeLength(index int) (int64, bool) {
	return makeGoLangTypeLength(r.cols[index].ti)
}

// It should return
// the precision and scale for decimal types. If not applicable, ok should be false.
// The following are examples of returned values for various types:
//   decimal(38, 4)    (38, 4, true)
//   int               (0, 0, false)
//   decimal           (math.MaxInt64, math.MaxInt64, true)
func (r *Rowsq) ColumnTypePrecisionScale(index int) (int64, int64, bool) {
	return makeGoLangTypePrecisionScale(r.cols[index].ti)
}

// The nullable value should
// be true if it is known the column may be null, or false if the column is known
// to be not nullable.
// If the column nullability is unknown, ok should be false.
func (r *Rowsq) ColumnTypeNullable(index int) (nullable, ok bool) {
	nullable = r.cols[index].Flags&colFlagNullable != 0
	ok = true
	return
}
