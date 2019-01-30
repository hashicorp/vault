/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"
	"time"

	"github.com/SAP/go-hdb/driver/sqltrace"

	p "github.com/SAP/go-hdb/internal/protocol"
)

// DriverVersion is the version number of the hdb driver.
const DriverVersion = "0.13.1"

// DriverName is the driver name to use with sql.Open for hdb databases.
const DriverName = "hdb"

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
var ErrUnsupportedIsolationLevel = errors.New("Unsupported isolation level")

// ErrNestedTransaction is the error raised if a tranasction is created within a transaction as this is not supported by hdb.
var ErrNestedTransaction = errors.New("Nested transactions are not supported")

// needed for testing
const driverDataFormatVersion = 1

// queries
const (
	pingQuery          = "select 1 from dummy"
	isolationLevelStmt = "set transaction isolation level %s"
	accessModeStmt     = "set transaction %s"
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

var drv = &hdbDrv{}

func init() {
	sql.Register(DriverName, drv)
}

// driver

//  check if driver implements all required interfaces
var (
	_ driver.Driver        = (*hdbDrv)(nil)
	_ driver.DriverContext = (*hdbDrv)(nil)
)

type hdbDrv struct{}

func (d *hdbDrv) Open(dsn string) (driver.Conn, error) {
	connector, err := NewDSNConnector(dsn)
	if err != nil {
		return nil, err
	}
	return connector.Connect(context.Background())
}

func (d *hdbDrv) OpenConnector(dsn string) (driver.Connector, error) {
	return NewDSNConnector(dsn)
}

// database connection

//  check if conn implements all required interfaces
var (
	_ driver.Conn               = (*conn)(nil)
	_ driver.ConnPrepareContext = (*conn)(nil)
	_ driver.Pinger             = (*conn)(nil)
	_ driver.ConnBeginTx        = (*conn)(nil)
	_ driver.ExecerContext      = (*conn)(nil)
	//go 1.9 issue (ExecerContext is only called if Execer is implemented)
	_ driver.Execer         = (*conn)(nil)
	_ driver.QueryerContext = (*conn)(nil)
	//go 1.9 issue (QueryerContext is only called if Queryer is implemented)
	// QueryContext is needed for stored procedures with table output parameters.
	_ driver.Queryer           = (*conn)(nil)
	_ driver.NamedValueChecker = (*conn)(nil)
)

type conn struct {
	session *p.Session
}

func newConn(ctx context.Context, c *Connector) (driver.Conn, error) {
	session, err := p.NewSession(ctx, c)
	if err != nil {
		return nil, err
	}
	return &conn{session: session}, nil
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	panic("deprecated")
}

func (c *conn) Close() error {
	c.session.Close()
	return nil
}

func (c *conn) Begin() (driver.Tx, error) {
	panic("deprecated")
}

func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {

	if c.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	if c.session.InTx() {
		return nil, ErrNestedTransaction
	}

	level, ok := isolationLevel[opts.Isolation]
	if !ok {
		return nil, ErrUnsupportedIsolationLevel
	}

	done := make(chan struct{})
	go func() {
		// set isolation level
		if _, err = c.ExecContext(ctx, fmt.Sprintf(isolationLevelStmt, level), nil); err != nil {
			goto done
		}
		// set access mode
		if _, err = c.ExecContext(ctx, fmt.Sprintf(accessModeStmt, readOnly[opts.ReadOnly]), nil); err != nil {
			goto done
		}
		c.session.SetInTx(true)
		tx = newTx(c.session)
	done:
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return tx, err
	}
}

// Exec implements the database/sql/driver/Execer interface.
// delete after go 1.9 compatibility is given up.
func (c *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	panic("deprecated")
}

// ExecContext implements the database/sql/driver/ExecerContext interface.
func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (r driver.Result, err error) {
	if c.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	if len(args) != 0 {
		return nil, driver.ErrSkip //fast path not possible (prepare needed)
	}

	sqltrace.Traceln(query)

	done := make(chan struct{})
	go func() {
		r, err = c.session.ExecDirect(query)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return r, err
	}
}

// Queryer implements the database/sql/driver/Queryer interface.
// delete after go 1.9 compatibility is given up.
func (c *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	panic("deprecated")
}

func (c *conn) Ping(ctx context.Context) (err error) {
	if c.session.IsBad() {
		return driver.ErrBadConn
	}

	done := make(chan struct{})
	go func() {
		_, err = c.QueryContext(ctx, pingQuery, nil)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return err
	}
}

// CheckNamedValue implements NamedValueChecker interface.
// implemented for conn:
// if querier or execer is called, sql checks parameters before
// in case of parameters the method can be 'skipped' and force the prepare path
// --> guarantee that a valid driver value is returned
// --> if not implemented, Lob need to have a pseudo Value method to return a valid driver value
func (c *conn) CheckNamedValue(nv *driver.NamedValue) error {
	switch nv.Value.(type) {
	case Lob, *Lob:
		nv.Value = nil
	}
	return nil
}

//transaction

//  check if tx implements all required interfaces
var (
	_ driver.Tx = (*tx)(nil)
)

type tx struct {
	session *p.Session
}

func newTx(session *p.Session) *tx {
	return &tx{
		session: session,
	}
}

func (t *tx) Commit() error {
	if t.session.IsBad() {
		return driver.ErrBadConn
	}

	return t.session.Commit()
}

func (t *tx) Rollback() error {
	if t.session.IsBad() {
		return driver.ErrBadConn
	}

	return t.session.Rollback()
}

//statement

var argsPool = sync.Pool{}

//  check if stmt implements all required interfaces
var (
	_ driver.Stmt              = (*stmt)(nil)
	_ driver.StmtExecContext   = (*stmt)(nil)
	_ driver.StmtQueryContext  = (*stmt)(nil)
	_ driver.NamedValueChecker = (*stmt)(nil)
)

type stmt struct {
	qt             p.QueryType
	session        *p.Session
	query          string
	id             uint64
	prmFieldSet    *p.ParameterFieldSet
	resultFieldSet *p.ResultFieldSet
	bulk, noFlush  bool
	numArg         int
	args           []driver.NamedValue
}

func newStmt(qt p.QueryType, session *p.Session, query string, id uint64, prmFieldSet *p.ParameterFieldSet, resultFieldSet *p.ResultFieldSet) (*stmt, error) {
	return &stmt{qt: qt, session: session, query: query, id: id, prmFieldSet: prmFieldSet, resultFieldSet: resultFieldSet}, nil
}

func (s *stmt) Close() error {
	if s.args != nil {
		if len(s.args) != 0 {
			sqltrace.Tracef("close: %s - not flushed records: %d)", s.query, int(len(s.args)/s.NumInput()))
		}
		argsPool.Put(s.args)
		s.args = nil
	}
	return s.session.DropStatementID(s.id)
}

func (s *stmt) NumInput() int {
	return s.prmFieldSet.NumInputField()
}

func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	panic("deprecated")
}

func (s *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (r driver.Result, err error) {
	if s.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	numField := s.prmFieldSet.NumInputField()
	if len(args) != numField {
		return nil, fmt.Errorf("invalid number of arguments %d - %d expected", len(args), numField)
	}

	sqltrace.Tracef("%s %v", s.query, args)

	// init noFlush
	noFlush := s.noFlush
	s.noFlush = false

	var _args []driver.NamedValue

	done := make(chan struct{})

	if !s.bulk {
		go func() {
			r, err = s.session.Exec(s.id, s.prmFieldSet, args)
			close(done)
		}()
		goto done
	}

	if s.args == nil {
		s.args, _ = argsPool.Get().([]driver.NamedValue)
		if s.args == nil {
			s.args = make([]driver.NamedValue, 0, len(args)*1000)
		}
		s.args = s.args[:0]
	}

	s.args = append(s.args, args...)
	s.numArg++

	if noFlush && s.numArg < maxSmallint { //TODO: check why bigArgument count does not work
		return driver.ResultNoRows, nil
	}

	_args, _ = argsPool.Get().([]driver.NamedValue)
	if _args == nil || cap(_args) < len(s.args) {
		_args = make([]driver.NamedValue, len(s.args))
	}
	_args = _args[:len(s.args)]

	copy(_args, s.args)
	s.args = s.args[:0]
	s.numArg = 0

	go func() {
		r, err = s.session.Exec(s.id, s.prmFieldSet, _args)
		argsPool.Put(_args)
		close(done)
	}()

done:
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return r, err
	}
}

func (s *stmt) Query(args []driver.Value) (rows driver.Rows, err error) {
	panic("deprecated")
}

// Deprecated: see NamedValueChecker.
//func (s *stmt) ColumnConverter(idx int) driver.ValueConverter {
//}

// CheckNamedValue implements NamedValueChecker interface.
func (s *stmt) CheckNamedValue(nv *driver.NamedValue) error {
	if nv.Name == bulk {
		if ptr, ok := nv.Value.(**struct{}); ok {
			switch ptr {
			case &noFlushTok:
				s.bulk, s.noFlush = true, true
				return driver.ErrRemoveArgument
			case &flushTok:
				return driver.ErrRemoveArgument
			}
		}
	}
	return checkNamedValue(s.prmFieldSet, nv)
}

// driver.Rows drop-in replacement if driver Query or QueryRow is used for statements that doesn't return rows
var noColumns = []string{}
var noResult = new(noResultType)

//  check if noResultType implements all required interfaces
var (
	_ driver.Rows = (*noResultType)(nil)
)

type noResultType struct{}

func (r *noResultType) Columns() []string              { return noColumns }
func (r *noResultType) Close() error                   { return nil }
func (r *noResultType) Next(dest []driver.Value) error { return io.EOF }

// rows
type rows struct {
}

// query result

//  check if queryResult implements all required interfaces
var (
	_ driver.Rows                           = (*queryResult)(nil)
	_ driver.RowsColumnTypeDatabaseTypeName = (*queryResult)(nil) // go 1.8
	_ driver.RowsColumnTypeLength           = (*queryResult)(nil) // go 1.8
	_ driver.RowsColumnTypeNullable         = (*queryResult)(nil) // go 1.8
	_ driver.RowsColumnTypePrecisionScale   = (*queryResult)(nil) // go 1.8
	_ driver.RowsColumnTypeScanType         = (*queryResult)(nil) // go 1.8
)

type queryResult struct {
	session        *p.Session
	id             uint64
	resultFieldSet *p.ResultFieldSet
	fieldValues    *p.FieldValues
	pos            int
	attrs          p.PartAttributes
	columns        []string
	lastErr        error
}

func newQueryResult(session *p.Session, id uint64, resultFieldSet *p.ResultFieldSet, fieldValues *p.FieldValues, attrs p.PartAttributes) (driver.Rows, error) {
	columns := make([]string, resultFieldSet.NumField())
	for i := 0; i < len(columns); i++ {
		columns[i] = resultFieldSet.Field(i).Name()
	}

	return &queryResult{
		session:        session,
		id:             id,
		resultFieldSet: resultFieldSet,
		fieldValues:    fieldValues,
		attrs:          attrs,
		columns:        columns,
	}, nil
}

func (r *queryResult) Columns() []string {
	return r.columns
}

func (r *queryResult) Close() error {
	// if lastError is set, attrs are nil
	if r.lastErr != nil {
		return r.lastErr
	}

	if !r.attrs.ResultsetClosed() {
		return r.session.CloseResultsetID(r.id)
	}
	return nil
}

func (r *queryResult) Next(dest []driver.Value) error {
	if r.session.IsBad() {
		return driver.ErrBadConn
	}

	if r.pos >= r.fieldValues.NumRow() {
		if r.attrs.LastPacket() {
			return io.EOF
		}

		var err error

		if r.attrs, err = r.session.FetchNext(r.id, r.resultFieldSet, r.fieldValues); err != nil {
			r.lastErr = err //fieldValues and attrs are nil
			return err
		}

		if r.attrs.NoRows() {
			return io.EOF
		}

		r.pos = 0

	}

	r.fieldValues.Row(r.pos, dest)
	r.pos++

	return nil
}

func (r *queryResult) ColumnTypeDatabaseTypeName(idx int) string {
	return r.resultFieldSet.Field(idx).TypeCode().TypeName()
}

func (r *queryResult) ColumnTypeLength(idx int) (int64, bool) {
	return r.resultFieldSet.Field(idx).TypeLength()
}

func (r *queryResult) ColumnTypePrecisionScale(idx int) (int64, int64, bool) {
	return r.resultFieldSet.Field(idx).TypePrecisionScale()
}

func (r *queryResult) ColumnTypeNullable(idx int) (bool, bool) {
	return r.resultFieldSet.Field(idx).Nullable(), true
}

var (
	scanTypeUnknown  = reflect.TypeOf(new(interface{})).Elem()
	scanTypeTinyint  = reflect.TypeOf(uint8(0))
	scanTypeSmallint = reflect.TypeOf(int16(0))
	scanTypeInteger  = reflect.TypeOf(int32(0))
	scanTypeBigint   = reflect.TypeOf(int64(0))
	scanTypeReal     = reflect.TypeOf(float32(0.0))
	scanTypeDouble   = reflect.TypeOf(float64(0.0))
	scanTypeTime     = reflect.TypeOf(time.Time{})
	scanTypeString   = reflect.TypeOf(string(""))
	scanTypeBytes    = reflect.TypeOf([]byte{})
	scanTypeDecimal  = reflect.TypeOf(Decimal{})
	scanTypeLob      = reflect.TypeOf(Lob{})
)

func (r *queryResult) ColumnTypeScanType(idx int) reflect.Type {
	switch r.resultFieldSet.Field(idx).TypeCode().DataType() {
	default:
		return scanTypeUnknown
	case p.DtTinyint:
		return scanTypeTinyint
	case p.DtSmallint:
		return scanTypeSmallint
	case p.DtInteger:
		return scanTypeInteger
	case p.DtBigint:
		return scanTypeBigint
	case p.DtReal:
		return scanTypeReal
	case p.DtDouble:
		return scanTypeDouble
	case p.DtTime:
		return scanTypeTime
	case p.DtDecimal:
		return scanTypeDecimal
	case p.DtString:
		return scanTypeString
	case p.DtBytes:
		return scanTypeBytes
	case p.DtLob:
		return scanTypeLob
	}
}
