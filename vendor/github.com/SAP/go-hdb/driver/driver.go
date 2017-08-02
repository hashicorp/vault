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
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sync"

	"github.com/SAP/go-hdb/driver/sqltrace"

	p "github.com/SAP/go-hdb/internal/protocol"
)

// DriverVersion is the version number of the hdb driver.
const DriverVersion = "0.9"

// DriverName is the driver name to use with sql.Open for hdb databases.
const DriverName = "hdb"

func init() {
	sql.Register(DriverName, &drv{})
}

var reBulk = regexp.MustCompile("(?i)^(\\s)*(bulk +)(.*)")

func checkBulkInsert(sql string) (string, bool) {
	if reBulk.MatchString(sql) {
		return reBulk.ReplaceAllString(sql, "${3}"), true
	}
	return sql, false
}

var reCall = regexp.MustCompile("(?i)^(\\s)*(call +)(.*)")

func checkCallProcedure(sql string) bool {
	return reCall.MatchString(sql)
}

var errProcTableQuery = errors.New("Invalid procedure table query")

// driver
type drv struct{}

func (d *drv) Open(dsn string) (driver.Conn, error) {
	return newConn(dsn)
}

// database connection
type conn struct {
	session *p.Session
}

func newConn(dsn string) (driver.Conn, error) {
	sessionPrm, err := parseDSN(dsn)
	if err != nil {
		return nil, err
	}
	session, err := p.NewSession(sessionPrm)
	if err != nil {
		return nil, err
	}
	return &conn{session: session}, nil
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	if c.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	prepareQuery, bulkInsert := checkBulkInsert(query)

	qt, id, parameterFieldSet, resultFieldSet, err := c.session.Prepare(prepareQuery)
	if err != nil {
		return nil, err
	}

	if bulkInsert {
		return newBulkInsertStmt(c.session, prepareQuery, id, parameterFieldSet)
	}
	return newStmt(qt, c.session, prepareQuery, id, parameterFieldSet, resultFieldSet)
}

func (c *conn) Close() error {
	c.session.Close()
	return nil
}

func (c *conn) Begin() (driver.Tx, error) {
	if c.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	if c.session.InTx() {
		return nil, fmt.Errorf("nested transactions are not supported")
	}

	c.session.SetInTx(true)
	return newTx(c.session), nil
}

// Exec implements the database/sql/driver/Execer interface.
func (c *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	if c.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	if len(args) != 0 {
		return nil, driver.ErrSkip //fast path not possible (prepare needed)
	}

	sqltrace.Traceln(query)

	return c.session.ExecDirect(query)
}

// bug?: check args is performed indepently of queryer raising ErrSkip or not
// - leads to different behavior to prepare - stmt - execute default logic
// - seems to be the same for Execer interface

// Queryer implements the database/sql/driver/Queryer interface.
func (c *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	if c.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	if len(args) != 0 {
		return nil, driver.ErrSkip //fast path not possible (prepare needed)
	}

	// direct execution of call procedure
	// - returns no parameter metadata (sps 82) but only field values
	// --> let's take the 'prepare way' for stored procedures
	if checkCallProcedure(query) {
		return nil, driver.ErrSkip
	}

	sqltrace.Traceln(query)

	id, idx, ok := decodeTableQuery(query)
	if ok {
		r := procedureCallResultStore.get(id)
		if r == nil {
			return nil, fmt.Errorf("invalid procedure table query %s", query)
		}
		return r.tableRows(int(idx))
	}

	id, meta, values, attributes, err := c.session.QueryDirect(query)
	if err != nil {
		return nil, err
	}
	if id == 0 { // non select query
		return noResult, nil
	}
	return newQueryResult(c.session, id, meta, values, attributes)
}

//transaction
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
type stmt struct {
	qt             p.QueryType
	session        *p.Session
	query          string
	id             uint64
	prmFieldSet    *p.FieldSet
	resultFieldSet *p.FieldSet
}

func newStmt(qt p.QueryType, session *p.Session, query string, id uint64, prmFieldSet *p.FieldSet, resultFieldSet *p.FieldSet) (*stmt, error) {
	return &stmt{qt: qt, session: session, query: query, id: id, prmFieldSet: prmFieldSet, resultFieldSet: resultFieldSet}, nil
}

func (s *stmt) Close() error {
	return s.session.DropStatementID(s.id)
}

func (s *stmt) NumInput() int {
	return s.prmFieldSet.NumInputField()
}

func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	if s.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	numField := s.prmFieldSet.NumInputField()
	if len(args) != numField {
		return nil, fmt.Errorf("invalid number of arguments %d - %d expected", len(args), numField)
	}

	sqltrace.Tracef("%s %v", s.query, args)

	return s.session.Exec(s.id, s.prmFieldSet, args)
}

func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	if s.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	switch s.qt {
	default:
		rows, err := s.defaultQuery(args)
		return rows, err
	case p.QtProcedureCall:
		rows, err := s.procedureCall(args)
		return rows, err
	}
}

func (s *stmt) defaultQuery(args []driver.Value) (driver.Rows, error) {

	sqltrace.Tracef("%s %v", s.query, args)

	rid, values, attributes, err := s.session.Query(s.id, s.prmFieldSet, s.resultFieldSet, args)
	if err != nil {
		return nil, err
	}
	if rid == 0 { // non select query
		return noResult, nil
	}
	return newQueryResult(s.session, rid, s.resultFieldSet, values, attributes)
}

func (s *stmt) procedureCall(args []driver.Value) (driver.Rows, error) {

	sqltrace.Tracef("%s %v", s.query, args)

	fieldValues, tableResults, err := s.session.Call(s.id, s.prmFieldSet, args)
	if err != nil {
		return nil, err
	}

	return newProcedureCallResult(s.session, s.prmFieldSet, fieldValues, tableResults)
}

func (s *stmt) ColumnConverter(idx int) driver.ValueConverter {
	return columnConverter(s.prmFieldSet.DataType(idx))
}

// bulk insert statement
type bulkInsertStmt struct {
	session           *p.Session
	query             string
	id                uint64
	parameterFieldSet *p.FieldSet
	numArg            int
	args              []driver.Value
}

func newBulkInsertStmt(session *p.Session, query string, id uint64, parameterFieldSet *p.FieldSet) (*bulkInsertStmt, error) {
	return &bulkInsertStmt{session: session, query: query, id: id, parameterFieldSet: parameterFieldSet, args: make([]driver.Value, 0)}, nil
}

func (s *bulkInsertStmt) Close() error {
	return s.session.DropStatementID(s.id)
}

func (s *bulkInsertStmt) NumInput() int {
	return -1
}

func (s *bulkInsertStmt) Exec(args []driver.Value) (driver.Result, error) {
	if s.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	sqltrace.Tracef("%s %v", s.query, args)

	if args == nil || len(args) == 0 {
		return s.execFlush()
	}
	return s.execBuffer(args)
}

func (s *bulkInsertStmt) execFlush() (driver.Result, error) {

	if s.numArg == 0 {
		return driver.ResultNoRows, nil
	}

	result, err := s.session.Exec(s.id, s.parameterFieldSet, s.args)
	s.args = s.args[:0]
	s.numArg = 0
	return result, err
}

func (s *bulkInsertStmt) execBuffer(args []driver.Value) (driver.Result, error) {

	numField := s.parameterFieldSet.NumInputField()
	if len(args) != numField {
		return nil, fmt.Errorf("invalid number of arguments %d - %d expected", len(args), numField)
	}

	var result driver.Result = driver.ResultNoRows
	var err error

	if s.numArg == maxSmallint { // TODO: check why bigArgument count does not work
		result, err = s.execFlush()
	}

	s.args = append(s.args, args...)
	s.numArg++

	return result, err
}

func (s *bulkInsertStmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, fmt.Errorf("query not allowed in context of bulk insert statement %s", s.query)
}

func (s *bulkInsertStmt) ColumnConverter(idx int) driver.ValueConverter {
	return columnConverter(s.parameterFieldSet.DataType(idx))
}

// driver.Rows drop-in replacement if driver Query or QueryRow is used for statements that doesn't return rows
var noColumns = []string{}
var noResult = new(noResultType)

type noResultType struct{}

func (r *noResultType) Columns() []string              { return noColumns }
func (r *noResultType) Close() error                   { return nil }
func (r *noResultType) Next(dest []driver.Value) error { return io.EOF }

// query result
type queryResult struct {
	session     *p.Session
	id          uint64
	fieldSet    *p.FieldSet
	fieldValues *p.FieldValues
	pos         int
	attrs       p.PartAttributes
	columns     []string
	lastErr     error
}

func newQueryResult(session *p.Session, id uint64, fieldSet *p.FieldSet, fieldValues *p.FieldValues, attrs p.PartAttributes) (driver.Rows, error) {
	columns := make([]string, fieldSet.NumOutputField())
	if err := fieldSet.OutputNames(columns); err != nil {
		return nil, err
	}

	return &queryResult{
		session:     session,
		id:          id,
		fieldSet:    fieldSet,
		fieldValues: fieldValues,
		attrs:       attrs,
		columns:     columns,
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

		if r.fieldValues, r.attrs, err = r.session.FetchNext(r.id, r.fieldSet); err != nil {
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

//call result store
type callResultStore struct {
	mu    sync.RWMutex
	store map[uint64]*procedureCallResult
	cnt   uint64
	free  []uint64
}

func (s *callResultStore) get(k uint64) *procedureCallResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if r, ok := s.store[k]; ok {
		return r
	}
	return nil
}

func (s *callResultStore) add(v *procedureCallResult) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	var k uint64

	if s.free == nil || len(s.free) == 0 {
		s.cnt++
		k = s.cnt
	} else {
		size := len(s.free)
		k = s.free[size-1]
		s.free = s.free[:size-1]
	}

	if s.store == nil {
		s.store = make(map[uint64]*procedureCallResult)
	}

	s.store[k] = v

	return k
}

func (s *callResultStore) del(k uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, k)

	if s.free == nil {
		s.free = []uint64{k}
	} else {
		s.free = append(s.free, k)
	}
}

var procedureCallResultStore = new(callResultStore)

//procedure call result
type procedureCallResult struct {
	id          uint64
	session     *p.Session
	fieldSet    *p.FieldSet
	fieldValues *p.FieldValues
	_tableRows  []driver.Rows
	columns     []string
	eof         error
}

func newProcedureCallResult(session *p.Session, fieldSet *p.FieldSet, fieldValues *p.FieldValues, tableResults []*p.TableResult) (driver.Rows, error) {

	fieldIdx := fieldSet.NumOutputField()
	columns := make([]string, fieldIdx+len(tableResults))
	if err := fieldSet.OutputNames(columns); err != nil {
		return nil, err
	}

	tableRows := make([]driver.Rows, len(tableResults))
	for i, tableResult := range tableResults {
		var err error

		if tableRows[i], err = newQueryResult(session, tableResult.ID(), tableResult.FieldSet(), tableResult.FieldValues(), tableResult.Attrs()); err != nil {
			return nil, err
		}

		columns[fieldIdx] = fmt.Sprintf("table %d", i)

		fieldIdx++

	}

	result := &procedureCallResult{
		session:     session,
		fieldSet:    fieldSet,
		fieldValues: fieldValues,
		_tableRows:  tableRows,
		columns:     columns,
	}
	id := procedureCallResultStore.add(result)
	result.id = id
	return result, nil
}

func (r *procedureCallResult) Columns() []string {
	return r.columns
}

func (r *procedureCallResult) Close() error {
	procedureCallResultStore.del(r.id)
	return nil
}

func (r *procedureCallResult) Next(dest []driver.Value) error {
	if r.session.IsBad() {
		return driver.ErrBadConn
	}

	if r.eof != nil {
		return r.eof
	}

	if r.fieldValues.NumRow() == 0 && len(r._tableRows) == 0 {
		r.eof = io.EOF
		return r.eof
	}

	if r.fieldValues.NumRow() != 0 {
		r.fieldValues.Row(0, dest)
	}

	i := r.fieldSet.NumOutputField()
	for j := range r._tableRows {
		dest[i] = encodeTableQuery(r.id, uint64(j))
		i++
	}

	r.eof = io.EOF
	return nil
}

func (r *procedureCallResult) tableRows(idx int) (driver.Rows, error) {
	if idx >= len(r._tableRows) {
		return nil, fmt.Errorf("table row index %d exceeds maximun %d", idx, len(r._tableRows)-1)
	}
	return r._tableRows[idx], nil
}

// helper
const tableQueryPrefix = "@tq"

func encodeTableQuery(id, idx uint64) string {
	start := len(tableQueryPrefix)
	b := make([]byte, start+8+8)
	copy(b, tableQueryPrefix)
	binary.LittleEndian.PutUint64(b[start:start+8], id)
	binary.LittleEndian.PutUint64(b[start+8:start+8+8], idx)
	return string(b)
}

func decodeTableQuery(query string) (uint64, uint64, bool) {
	size := len(query)
	start := len(tableQueryPrefix)
	if size != start+8+8 {
		return 0, 0, false
	}
	if query[:start] != tableQueryPrefix {
		return 0, 0, false
	}
	id := binary.LittleEndian.Uint64([]byte(query[start : start+8]))
	idx := binary.LittleEndian.Uint64([]byte(query[start+8 : start+8+8]))
	return id, idx, true
}
