package driver

import (
	"context"
	"database/sql/driver"
	"errors"
	"io"
	"reflect"

	p "github.com/SAP/go-hdb/driver/internal/protocol"
)

// check if rows types do implement all driver row interfaces.
var (
	// queryResult.
	_ driver.Rows                           = (*queryResult)(nil)
	_ driver.RowsColumnTypeDatabaseTypeName = (*queryResult)(nil)
	_ driver.RowsColumnTypeLength           = (*queryResult)(nil)
	_ driver.RowsColumnTypeNullable         = (*queryResult)(nil)
	_ driver.RowsColumnTypePrecisionScale   = (*queryResult)(nil)
	_ driver.RowsColumnTypeScanType         = (*queryResult)(nil)
	//	currently not used
	//	could be implemented as pointer to next queryResult (advancing by copying data from next)
	//	_ driver.RowsNextResultSet = (*queryResult)(nil)

	// noResultType.
	_ driver.Rows = (*noResultType)(nil)
	// callResult.
	_ driver.Rows = (*callResult)(nil)
)

type prepareResult struct {
	fc              p.FunctionCode
	stmtID          uint64
	parameterFields []*p.ParameterField
	resultFields    []*p.ResultField
}

// ParameterTypes implements the PrepareMetadata interface.
func (pr *prepareResult) ParameterTypes() []ParameterType {
	parameterTypes := make([]ParameterType, len(pr.parameterFields))
	for i, f := range pr.parameterFields {
		parameterTypes[i] = f
	}
	return parameterTypes
}

func (pr *prepareResult) columnTypes() []ColumnType {
	columnTypes := make([]ColumnType, len(pr.resultFields))
	for i, f := range pr.resultFields {
		columnTypes[i] = f
	}
	return columnTypes
}

func (pr *prepareResult) procedureCallColumnTypes() []ColumnType {
	var columnTypes []ColumnType
	for _, f := range pr.parameterFields {
		if f.InOut() || f.Out() {
			columnTypes = append(columnTypes, f)
		}
	}
	return columnTypes
}

// ColumnTypes implements the PrepareMetadata interface.
func (pr *prepareResult) ColumnTypes() []ColumnType {
	if pr.isProcedureCall() {
		return pr.procedureCallColumnTypes()
	}
	return pr.columnTypes()
}

// isProcedureCall returns true if the statement is a call statement.
func (pr *prepareResult) isProcedureCall() bool { return pr.fc.IsProcedureCall() }

// numField returns the number of parameter fields in a database statement.
func (pr *prepareResult) numField() int { return len(pr.parameterFields) }

// NoResult is the driver.Rows drop-in replacement if driver Query or QueryRow is used for statements that do not return rows.
var noResult = new(noResultType)

var noColumns = []string{}

type noResultType struct{}

func (r *noResultType) Columns() []string              { return noColumns }
func (r *noResultType) Close() error                   { return nil }
func (r *noResultType) Next(dest []driver.Value) error { return io.EOF }

// queryResult represents the resultset of a query.
type queryResult struct {
	// field alignment
	fields       []*p.ResultField
	fieldValues  []driver.Value
	decodeErrors p.DecodeErrors
	_columns     []string
	lastErr      error
	session      *session
	rsID         uint64
	pos          int
	attrs        p.PartAttributes
	closed       bool
}

// ErrScanOnClosedResultset is the error raised in case a scan is executed on a closed resultset.
var ErrScanOnClosedResultset = errors.New("scan on closed resultset")

// Columns implements the driver.Rows interface.
func (qr *queryResult) Columns() []string {
	if qr._columns != nil {
		return qr._columns
	}
	qr._columns = make([]string, len(qr.fields))
	for i, f := range qr.fields {
		qr._columns[i] = f.Name()
	}
	return qr._columns
}

// Close implements the driver.Rows interface.
func (qr *queryResult) Close() error {
	qr.closed = true
	if qr.attrs.ResultsetClosed() {
		return nil
	}
	// if lastError is set, attrs are nil
	if qr.lastErr != nil {
		return qr.lastErr
	}
	return qr.session.closeResultsetID(context.Background(), qr.rsID)
}

func (qr *queryResult) numRow() int {
	if len(qr.fieldValues) == 0 {
		return 0
	}
	return len(qr.fieldValues) / len(qr.fields)
}

// Next implements the driver.Rows interface.
func (qr *queryResult) Next(dest []driver.Value) error {
	if qr.pos >= qr.numRow() {
		if qr.attrs.LastPacket() {
			return io.EOF
		}
		if err := qr.session.fetchNext(context.Background(), qr); err != nil {
			qr.lastErr = err // fieldValues and attrs are nil
			return err
		}
		if qr.numRow() == 0 {
			return io.EOF
		}
		qr.pos = 0
	}

	// copy row.
	cols := len(qr.fields)
	copy(dest, qr.fieldValues[qr.pos*cols:(qr.pos+1)*cols])
	err := qr.decodeErrors.RowErrors(qr.pos)
	qr.pos++
	return err
}

// ColumnTypeDatabaseTypeName implements the driver.RowsColumnTypeDatabaseTypeName interface.
func (qr *queryResult) ColumnTypeDatabaseTypeName(idx int) string {
	return qr.fields[idx].DatabaseTypeName()
}

// ColumnTypeLength implements the driver.RowsColumnTypeLength interface.
func (qr *queryResult) ColumnTypeLength(idx int) (int64, bool) { return qr.fields[idx].Length() }

// ColumnTypeNullable implements the driver.RowsColumnTypeNullable interface.
func (qr *queryResult) ColumnTypeNullable(idx int) (bool, bool) { return qr.fields[idx].Nullable() }

// ColumnTypePrecisionScale implements the driver.RowsColumnTypePrecisionScale interface.
func (qr *queryResult) ColumnTypePrecisionScale(idx int) (int64, int64, bool) {
	return qr.fields[idx].DecimalSize()
}

// ColumnTypeScanType implements the driver.RowsColumnTypeScanType interface.
func (qr *queryResult) ColumnTypeScanType(idx int) reflect.Type { return qr.fields[idx].ScanType() }

// ReadLob used by protocol LobReader.
func (qr *queryResult) ReadLob(request *p.ReadLobRequest, reply *p.ReadLobReply) error {
	if qr.closed {
		return ErrScanOnClosedResultset
	}
	return qr.session.readLob(context.Background(), request, reply)
}

type callResult struct { // call output parameters
	session      *session
	outFields    []*p.ParameterField
	fieldValues  []driver.Value
	decodeErrors p.DecodeErrors
	_columns     []string
	eof          bool
	closed       bool
}

// Columns implements the driver.Rows interface.
func (cr *callResult) Columns() []string {
	if cr._columns != nil {
		return cr._columns
	}
	cr._columns = make([]string, len(cr.outFields))
	for i, f := range cr.outFields {
		cr._columns[i] = f.Name()
	}
	return cr._columns
}

// Next implements the driver.Rows interface.
func (cr *callResult) Next(dest []driver.Value) error {
	if len(cr.fieldValues) == 0 || cr.eof {
		return io.EOF
	}

	cr.eof = true
	copy(dest, cr.fieldValues)
	return cr.decodeErrors.RowErrors(0)
}

// Close implements the driver.Rows interface.
func (cr *callResult) Close() error { cr.closed = true; return nil }

// ReadLob used by protocol LobReader.
func (cr *callResult) ReadLob(request *p.ReadLobRequest, reply *p.ReadLobReply) error {
	if cr.closed {
		return ErrScanOnClosedResultset
	}
	return cr.session.readLob(context.Background(), request, reply)
}
