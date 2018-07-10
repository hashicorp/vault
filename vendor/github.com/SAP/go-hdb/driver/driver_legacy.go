// +build !future

/*
Copyright 2018 SAP SE

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
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"io"
	"regexp"
	"sync"

	"github.com/SAP/go-hdb/driver/sqltrace"

	p "github.com/SAP/go-hdb/internal/protocol"
)

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

// database connection

func (c *conn) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	if c.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	done := make(chan struct{})
	go func() {
		prepareQuery, bulkInsert := checkBulkInsert(query)
		var (
			qt             p.QueryType
			id             uint64
			prmFieldSet    *p.ParameterFieldSet
			resultFieldSet *p.ResultFieldSet
		)
		qt, id, prmFieldSet, resultFieldSet, err = c.session.Prepare(prepareQuery)
		if err != nil {
			goto done
		}
		select {
		default:
		case <-ctx.Done():
			return
		}
		if bulkInsert {
			stmt, err = newBulkInsertStmt(c.session, prepareQuery, id, prmFieldSet)
		} else {
			stmt, err = newStmt(qt, c.session, prepareQuery, id, prmFieldSet, resultFieldSet)
		}
	done:
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return stmt, err
	}
}

// QueryContext implements the database/sql/driver/QueryerContext interface.
func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
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

	done := make(chan struct{})
	go func() {
		var (
			id             uint64
			resultFieldSet *p.ResultFieldSet
			fieldValues    *p.FieldValues
			attributes     p.PartAttributes
		)
		id, resultFieldSet, fieldValues, attributes, err = c.session.QueryDirect(query)
		if err != nil {
			goto done
		}
		select {
		default:
		case <-ctx.Done():
			return
		}
		if id == 0 { // non select query
			rows = noResult
		} else {
			rows, err = newQueryResult(c.session, id, resultFieldSet, fieldValues, attributes)
		}
	done:
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return rows, err
	}
}

//statement

func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {

	if s.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	done := make(chan struct{})
	go func() {
		switch s.qt {
		default:
			rows, err = s.defaultQuery(ctx, args)
		case p.QtProcedureCall:
			rows, err = s.procedureCall(ctx, args)
		}
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return rows, err
	}
}

func (s *stmt) defaultQuery(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {

	sqltrace.Tracef("%s %v", s.query, args)

	rid, values, attributes, err := s.session.Query(s.id, s.prmFieldSet, s.resultFieldSet, args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if rid == 0 { // non select query
		return noResult, nil
	}
	return newQueryResult(s.session, rid, s.resultFieldSet, values, attributes)
}

func (s *stmt) procedureCall(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {

	sqltrace.Tracef("%s %v", s.query, args)

	fieldValues, tableResults, err := s.session.Call(s.id, s.prmFieldSet, args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return newProcedureCallResult(s.session, s.prmFieldSet, fieldValues, tableResults)
}

// bulk insert statement

//  check if bulkInsertStmt implements all required interfaces
var (
	_ driver.Stmt              = (*bulkInsertStmt)(nil)
	_ driver.StmtExecContext   = (*bulkInsertStmt)(nil)
	_ driver.StmtQueryContext  = (*bulkInsertStmt)(nil)
	_ driver.NamedValueChecker = (*bulkInsertStmt)(nil)
)

type bulkInsertStmt struct {
	session     *p.Session
	query       string
	id          uint64
	prmFieldSet *p.ParameterFieldSet
	numArg      int
	args        []driver.NamedValue
}

func newBulkInsertStmt(session *p.Session, query string, id uint64, prmFieldSet *p.ParameterFieldSet) (*bulkInsertStmt, error) {
	return &bulkInsertStmt{session: session, query: query, id: id, prmFieldSet: prmFieldSet, args: make([]driver.NamedValue, 0)}, nil
}

func (s *bulkInsertStmt) Close() error {
	return s.session.DropStatementID(s.id)
}

func (s *bulkInsertStmt) NumInput() int {
	return -1
}

func (s *bulkInsertStmt) Exec(args []driver.Value) (driver.Result, error) {
	panic("deprecated")
}

func (s *bulkInsertStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (r driver.Result, err error) {

	if s.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	sqltrace.Tracef("%s %v", s.query, args)

	done := make(chan struct{})
	go func() {
		if args == nil || len(args) == 0 {
			r, err = s.execFlush()
		} else {
			r, err = s.execBuffer(args)
		}
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return r, err
	}
}

func (s *bulkInsertStmt) execFlush() (driver.Result, error) {

	if s.numArg == 0 {
		return driver.ResultNoRows, nil
	}

	sqltrace.Traceln("execFlush")

	result, err := s.session.Exec(s.id, s.prmFieldSet, s.args)
	s.args = s.args[:0]
	s.numArg = 0
	return result, err
}

func (s *bulkInsertStmt) execBuffer(args []driver.NamedValue) (driver.Result, error) {

	numField := s.prmFieldSet.NumInputField()
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
	panic("deprecated")
}

func (s *bulkInsertStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	return nil, fmt.Errorf("query not allowed in context of bulk insert statement %s", s.query)
}

// Deprecated: see NamedValueChecker.
//func (s *bulkInsertStmt) ColumnConverter(idx int) driver.ValueConverter {
//}

// CheckNamedValue implements NamedValueChecker interface.
func (s *bulkInsertStmt) CheckNamedValue(nv *driver.NamedValue) error {
	return checkNamedValue(s.prmFieldSet, nv)
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

//  check if procedureCallResult implements all required interfaces
var _ driver.Rows = (*procedureCallResult)(nil)

type procedureCallResult struct {
	id          uint64
	session     *p.Session
	prmFieldSet *p.ParameterFieldSet
	fieldValues *p.FieldValues
	_tableRows  []driver.Rows
	columns     []string
	eof         error
}

func newProcedureCallResult(session *p.Session, prmFieldSet *p.ParameterFieldSet, fieldValues *p.FieldValues, tableResults []*p.TableResult) (driver.Rows, error) {

	fieldIdx := prmFieldSet.NumOutputField()
	columns := make([]string, fieldIdx+len(tableResults))

	for i := 0; i < fieldIdx; i++ {
		columns[i] = prmFieldSet.OutputField(i).Name()
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
		prmFieldSet: prmFieldSet,
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

	i := r.prmFieldSet.NumOutputField()
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
