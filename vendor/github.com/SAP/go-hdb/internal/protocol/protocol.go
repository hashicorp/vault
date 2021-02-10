// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"bufio"
	"database/sql/driver"
	"fmt"
	"io"
	"math"
	"reflect"

	"github.com/SAP/go-hdb/driver/sqltrace"
	"github.com/SAP/go-hdb/internal/container/vermap"
	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

// A PrepareResult represents the result of a prepare statement.
type PrepareResult struct {
	fc              functionCode
	stmtID          uint64
	parameterFields []*ParameterField
	resultFields    []*resultField
}

// Check checks consistency of the prepare result.
func (pr *PrepareResult) Check(qd *QueryDescr) error {
	call := qd.kind == QkCall
	if call != pr.fc.isProcedureCall() {
		return fmt.Errorf("function code mismatch: query descriptor %s - function code %s", qd.kind, pr.fc)
	}

	if !call {
		// only input parameters allowed
		for _, f := range pr.parameterFields {
			if f.Out() {
				return fmt.Errorf("invalid parameter %s", f)
			}
		}
	}
	return nil
}

// StmtID returns the statement id.
func (pr *PrepareResult) StmtID() uint64 {
	return pr.stmtID
}

// IsProcedureCall returns true if the statement is a call statement.
func (pr *PrepareResult) IsProcedureCall() bool {
	return pr.fc.isProcedureCall()
}

// NumField returns the number of parameter fields in a database statement.
func (pr *PrepareResult) NumField() int {
	return len(pr.parameterFields)
}

// NumInputField returns the number of input fields in a database statement.
func (pr *PrepareResult) NumInputField() int {
	if !pr.fc.isProcedureCall() {
		return len(pr.parameterFields) // only input fields
	}
	numField := 0
	for _, f := range pr.parameterFields {
		if f.In() {
			numField++
		}
	}
	return numField
}

// ParameterField returns the parameter field at index idx.
func (pr *PrepareResult) ParameterField(idx int) *ParameterField {
	return pr.parameterFields[idx]
}

// OnCloser defines getter and setter for a function which should be called when closing.
type OnCloser interface {
	OnClose() func()
	SetOnClose(func())
}

// NoResult is the driver.Rows drop-in replacement if driver Query or QueryRow is used for statements that do not return rows.
var noResult = new(noResultType)

//  check if noResultType implements all required interfaces
var (
	_ driver.Rows = (*noResultType)(nil)
)

var noColumns = []string{}

type noResultType struct{}

func (r *noResultType) Columns() []string              { return noColumns }
func (r *noResultType) Close() error                   { return nil }
func (r *noResultType) Next(dest []driver.Value) error { return io.EOF }

// check if queryResult does implement all driver row interfaces.
var (
	_ driver.Rows                           = (*queryResult)(nil)
	_ driver.RowsColumnTypeDatabaseTypeName = (*queryResult)(nil)
	_ driver.RowsColumnTypeLength           = (*queryResult)(nil)
	_ driver.RowsColumnTypeNullable         = (*queryResult)(nil)
	_ driver.RowsColumnTypePrecisionScale   = (*queryResult)(nil)
	_ driver.RowsColumnTypeScanType         = (*queryResult)(nil)
	_ driver.RowsNextResultSet              = (*queryResult)(nil)
)

// A QueryResult represents the resultset of a query.
type queryResult struct {
	// field alignment
	fields      []*resultField
	fieldValues []driver.Value
	_columns    []string
	lastErr     error
	session     *Session
	rsID        uint64
	pos         int
	onClose     func()
	attributes  partAttributes
	closed      bool
}

// OnClose implements the OnCloser interface
func (qr *queryResult) OnClose() func() { return qr.onClose }

// SetOnClose implements the OnCloser interface
func (qr *queryResult) SetOnClose(f func()) { qr.onClose = f }

// Columns implements the driver.Rows interface.
func (qr *queryResult) Columns() []string {
	if qr._columns == nil {
		numField := len(qr.fields)
		qr._columns = make([]string, numField)
		for i := 0; i < numField; i++ {
			qr._columns[i] = qr.fields[i].name()
		}
	}
	return qr._columns
}

// Close implements the driver.Rows interface.
func (qr *queryResult) Close() error {
	if !qr.closed && qr.onClose != nil {
		defer qr.onClose()
	}
	qr.closed = true

	if qr.attributes.ResultsetClosed() {
		return nil
	}
	// if lastError is set, attrs are nil
	if qr.lastErr != nil {
		return qr.lastErr
	}
	return qr.session.CloseResultsetID(qr.rsID)
}

func (qr *queryResult) numRow() int {
	if len(qr.fieldValues) == 0 {
		return 0
	}
	return len(qr.fieldValues) / len(qr.fields)
}

func (qr *queryResult) copyRow(idx int, dest []driver.Value) {
	cols := len(qr.fields)
	copy(dest, qr.fieldValues[idx*cols:(idx+1)*cols])
}

// Next implements the driver.Rows interface.
func (qr *queryResult) Next(dest []driver.Value) error {
	if qr.pos >= qr.numRow() {
		if qr.attributes.LastPacket() {
			return io.EOF
		}
		if err := qr.session.fetchNext(qr); err != nil {
			qr.lastErr = err //fieldValues and attrs are nil
			return err
		}
		if qr.numRow() == 0 {
			return io.EOF
		}
		qr.pos = 0
	}

	qr.copyRow(qr.pos, dest)
	qr.pos++

	// TODO eliminate
	for _, v := range dest {
		if v, ok := v.(sessionSetter); ok {
			v.setSession(qr.session)
		}
	}
	return nil
}

// ColumnTypeDatabaseTypeName implements the driver.RowsColumnTypeDatabaseTypeName interface.
func (qr *queryResult) ColumnTypeDatabaseTypeName(idx int) string { return qr.fields[idx].typeName() }

// ColumnTypeLength implements the driver.RowsColumnTypeLength interface.
func (qr *queryResult) ColumnTypeLength(idx int) (int64, bool) { return qr.fields[idx].typeLength() }

// ColumnTypeNullable implements the driver.RowsColumnTypeNullable interface.
func (qr *queryResult) ColumnTypeNullable(idx int) (bool, bool) {
	return qr.fields[idx].nullable(), true
}

// ColumnTypePrecisionScale implements the driver.RowsColumnTypePrecisionScale interface.
func (qr *queryResult) ColumnTypePrecisionScale(idx int) (int64, int64, bool) {
	return qr.fields[idx].typePrecisionScale()
}

// ColumnTypeScanType implements the driver.RowsColumnTypeScanType interface.
func (qr *queryResult) ColumnTypeScanType(idx int) reflect.Type {
	return scanTypeMap[qr.fields[idx].scanType()]
}

/*
driver.RowsNextResultSet:
- currently not used
- could be implemented as pointer to next queryResult (advancing by copying data from next)
*/

// HasNextResultSet implements the driver.RowsNextResultSet interface.
func (qr *queryResult) HasNextResultSet() bool { return false }

// NextResultSet implements the driver.RowsNextResultSet interface.
func (qr *queryResult) NextResultSet() error { return io.EOF }

// check if callResult does implement all driver row interfaces.
var (
	_ driver.Rows                           = (*callResult)(nil)
	_ driver.RowsColumnTypeDatabaseTypeName = (*callResult)(nil)
	_ driver.RowsColumnTypeLength           = (*callResult)(nil)
	_ driver.RowsColumnTypeNullable         = (*callResult)(nil)
	_ driver.RowsColumnTypePrecisionScale   = (*callResult)(nil)
	_ driver.RowsColumnTypeScanType         = (*callResult)(nil)
	_ driver.RowsNextResultSet              = (*callResult)(nil)
)

// A CallResult represents the result (output parameters and values) of a call statement.
type callResult struct { // call output parameters
	session      *Session
	outputFields []*ParameterField
	fieldValues  []driver.Value
	_columns     []string
	qrs          []*queryResult // table output parameters
	eof          bool
	closed       bool
	onClose      func()
}

// OnClose implements the OnCloser interface
func (cr *callResult) OnClose() func() { return cr.onClose }

// SetOnClose implements the OnCloser interface
func (cr *callResult) SetOnClose(f func()) { cr.onClose = f }

// Columns implements the driver.Rows interface.
func (cr *callResult) Columns() []string {
	if cr._columns == nil {
		numField := len(cr.outputFields)
		cr._columns = make([]string, numField)
		for i := 0; i < numField; i++ {
			cr._columns[i] = cr.outputFields[i].name()
		}
	}
	return cr._columns
}

/// Next implements the driver.Rows interface.
func (cr *callResult) Next(dest []driver.Value) error {
	if len(cr.fieldValues) == 0 || cr.eof {
		return io.EOF
	}

	copy(dest, cr.fieldValues)
	cr.eof = true
	// TODO eliminate
	for _, v := range dest {
		if v, ok := v.(sessionSetter); ok {
			v.setSession(cr.session)
		}
	}
	return nil
}

// Close implements the driver.Rows interface.
func (cr *callResult) Close() error {
	if !cr.closed && cr.onClose != nil {
		cr.onClose()
	}
	cr.closed = true
	return nil
}

// ColumnTypeDatabaseTypeName implements the driver.RowsColumnTypeDatabaseTypeName interface.
func (cr *callResult) ColumnTypeDatabaseTypeName(idx int) string {
	return cr.outputFields[idx].typeName()
}

// ColumnTypeLength implements the driver.RowsColumnTypeLength interface.
func (cr *callResult) ColumnTypeLength(idx int) (int64, bool) {
	return cr.outputFields[idx].typeLength()
}

// ColumnTypeNullable implements the driver.RowsColumnTypeNullable interface.
func (cr *callResult) ColumnTypeNullable(idx int) (bool, bool) {
	return cr.outputFields[idx].nullable(), true
}

// ColumnTypePrecisionScale implements the driver.RowsColumnTypePrecisionScale interface.
func (cr *callResult) ColumnTypePrecisionScale(idx int) (int64, int64, bool) {
	return cr.outputFields[idx].typePrecisionScale()
}

// ColumnTypeScanType implements the driver.RowsColumnTypeScanType interface.
func (cr *callResult) ColumnTypeScanType(idx int) reflect.Type {
	return scanTypeMap[cr.outputFields[idx].scanType()]
}

/*
driver.RowsNextResultSet:
- currently not used
- could be implemented as pointer to next queryResult (advancing by copying data from next)
*/

// HasNextResultSet implements the driver.RowsNextResultSet interface.
func (cr *callResult) HasNextResultSet() bool { return false }

// NextResultSet implements the driver.RowsNextResultSet interface.
func (cr *callResult) NextResultSet() error { return io.EOF }

//
func (cr *callResult) appendTableRefFields() {
	for i, qr := range cr.qrs {
		cr.outputFields = append(cr.outputFields, &ParameterField{fieldName: fmt.Sprintf("table %d", i), tc: tcTableRef, mode: pmOut, offset: 0})
		cr.fieldValues = append(cr.fieldValues, encodeID(qr.rsID))
	}
}

func (cr *callResult) appendTableRowsFields() {
	for i, qr := range cr.qrs {
		cr.outputFields = append(cr.outputFields, &ParameterField{fieldName: fmt.Sprintf("table %d", i), tc: tcTableRows, mode: pmOut, offset: 0})
		cr.fieldValues = append(cr.fieldValues, qr)
	}
}

type protocolReader struct {
	upStream bool

	step int // authentication

	dec    *encoding.Decoder
	tracer traceLogger

	mh *messageHeader
	sh *segmentHeader
	ph *partHeader

	msgSize  int64
	numPart  int
	cntPart  int
	partRead bool

	partReaderCache map[partKind]partReader

	lastErrors       *hdbErrors
	lastRowsAffected *rowsAffected

	// partReader read errors could be
	// - read buffer errors -> buffer Error() and ResetError()
	// - plus other errors (which cannot be ignored, e.g. Lob reader)
	err error
}

func newProtocolReader(upStream bool, rd io.Reader) *protocolReader {
	return &protocolReader{
		upStream:        upStream,
		dec:             encoding.NewDecoder(rd),
		tracer:          newTraceLogger(upStream),
		partReaderCache: map[partKind]partReader{},
		mh:              &messageHeader{},
		sh:              &segmentHeader{},
		ph:              &partHeader{},
	}
}

func (r *protocolReader) setDfv(dfv int) {
	r.dec.SetDfv(dfv)
}

func (r *protocolReader) readSkip() error            { return r.iterateParts(nil) }
func (r *protocolReader) sessionID() int64           { return r.mh.sessionID }
func (r *protocolReader) functionCode() functionCode { return r.sh.functionCode }

func (r *protocolReader) readInitRequest() error {
	req := &initRequest{}
	if err := req.decode(r.dec); err != nil {
		return err
	}
	r.tracer.Log(req)
	return nil
}

func (r *protocolReader) readInitReply() error {
	rep := &initReply{}
	if err := rep.decode(r.dec); err != nil {
		return err
	}
	r.tracer.Log(rep)
	return nil
}

func (r *protocolReader) readProlog() error {
	if r.upStream {
		return r.readInitRequest()
	}
	return r.readInitReply()
}

func (r *protocolReader) checkError() error {
	defer func() { // init readFlags
		r.lastErrors = nil
		r.lastRowsAffected = nil
		r.err = nil
		r.dec.ResetError()
	}()

	if r.err != nil {
		return r.err
	}

	if err := r.dec.Error(); err != nil {
		return err
	}

	if r.lastErrors == nil {
		return nil
	}

	if r.lastRowsAffected != nil { // link statement to error
		j := 0
		for i, rows := range *r.lastRowsAffected {
			if rows == raExecutionFailed {
				r.lastErrors.setStmtNo(j, i)
				j++
			}
		}
	}

	if r.lastErrors.isWarnings() {
		if sqltrace.On() {
			for _, e := range r.lastErrors.errors {
				sqltrace.Traceln(e)
			}
		}
		return nil
	}

	return r.lastErrors
}

func (r *protocolReader) canSkip(pk partKind) bool {
	// errors and rowsAffected needs always to be read
	if pk == pkError || pk == pkRowsAffected {
		return false
	}
	if debug {
		return false
	}
	return true
}

func (r *protocolReader) read(part partReader) error {
	r.partRead = true

	err := r.readPart(part)
	if err != nil {
		r.err = err
	}

	switch part := part.(type) {
	case *hdbErrors:
		r.lastErrors = part
	case *rowsAffected:
		r.lastRowsAffected = part
	}
	return err
}

func (r *protocolReader) authPart() partReader {
	defer func() { r.step++ }()

	switch {
	case r.upStream && r.step == 0:
		return &authInitReq{}
	case r.upStream:
		return &authFinalReq{}
	case !r.upStream && r.step == 0:
		return &authInitRep{}
	case !r.upStream:
		return &authFinalRep{}
	default:
		panic(fmt.Errorf("invalid auth step in protocol reader %d", r.step))
	}
}

func (r *protocolReader) defaultPart(pk partKind) (partReader, error) {
	part, ok := r.partReaderCache[pk]
	if !ok {
		var err error
		part, err = newPartReader(pk)
		if err != nil {
			return nil, err
		}
		r.partReaderCache[pk] = part
	}
	return part, nil
}

func (r *protocolReader) skip() error {
	pk := r.ph.partKind
	if r.canSkip(pk) {
		return r.skipPart()
	}

	var part partReader
	var err error
	if pk == pkAuthentication {
		part = r.authPart()
	} else {
		part, err = r.defaultPart(pk)
	}
	if err != nil {
		return r.skipPart()
	}
	return r.read(part)
}

func (r *protocolReader) skipPart() error {
	r.dec.Skip(int(r.ph.bufferLength))
	r.tracer.Log("*skipped")

	/*
		hdb protocol
		- in general padding but
		- in some messages the last record sent is not padded
		  - message header varPartLength < segment header segmentLength
		    - msgSize == 0: mh.varPartLength == sh.segmentLength
			- msgSize < 0 : mh.varPartLength < sh.segmentLength
	*/
	if r.cntPart != r.numPart || r.msgSize == 0 {
		r.dec.Skip(padBytes(int(r.ph.bufferLength)))
	}
	return nil
}

func (r *protocolReader) readPart(part partReader) error {

	r.dec.ResetCnt()
	if err := part.decode(r.dec, r.ph); err != nil {
		return err // do not ignore partReader errros
	}
	cnt := r.dec.Cnt()
	r.tracer.Log(part)

	bufferLen := int(r.ph.bufferLength)
	switch {
	case cnt < bufferLen: // protocol buffer length > read bytes -> skip the unread bytes

		// TODO enable for debug
		// b := make([]byte, bufferLen-cnt)
		// p.rd.ReadFull(b)
		// println(fmt.Sprintf("%x", b))
		// println(string(b))

		r.dec.Skip(bufferLen - cnt)

	case cnt > bufferLen: // read bytes > protocol buffer length -> should never happen
		panic(fmt.Errorf("protocol error: read bytes %d > buffer length %d", cnt, bufferLen))
	}

	/*
		hdb protocol
		- in general padding but
		- in some messages the last record sent is not padded
		  - message header varPartLength < segment header segmentLength
		    - msgSize == 0: mh.varPartLength == sh.segmentLength
			- msgSize < 0 : mh.varPartLength < sh.segmentLength
	*/
	if r.cntPart != r.numPart || r.msgSize == 0 {
		r.dec.Skip(padBytes(int(r.ph.bufferLength)))
	}
	return nil
}

func (r *protocolReader) iterateParts(partCb func(ph *partHeader)) error {
	if err := r.mh.decode(r.dec); err != nil {
		return err
	}
	r.tracer.Log(r.mh)

	r.msgSize = int64(r.mh.varPartLength)

	for i := 0; i < int(r.mh.noOfSegm); i++ {

		if err := r.sh.decode(r.dec); err != nil {
			return err
		}
		r.tracer.Log(r.sh)

		r.msgSize -= int64(r.sh.segmentLength)
		r.numPart = int(r.sh.noOfParts)
		r.cntPart = 0

		for j := 0; j < int(r.sh.noOfParts); j++ {

			if err := r.ph.decode(r.dec); err != nil {
				return err
			}
			r.tracer.Log(r.ph)

			r.cntPart++

			r.partRead = false
			if partCb != nil {
				partCb(r.ph)
			}
			if !r.partRead {
				r.skip()
			}
			if r.err != nil {
				return r.err
			}
		}
	}
	return r.checkError()
}

// protocol writer
type protocolWriter struct {
	wr  *bufio.Writer
	enc *encoding.Encoder

	sv *vermap.VerMap // link to session variables
	// last session variables snapshot
	lastSVVersion int64
	lastSV        map[string]string

	tracer traceLogger

	// reuse header
	mh *messageHeader
	sh *segmentHeader
	ph *partHeader
}

func newProtocolWriter(wr *bufio.Writer, sv *vermap.VerMap) *protocolWriter {
	return &protocolWriter{
		wr:     wr,
		sv:     sv,
		enc:    encoding.NewEncoder(wr),
		tracer: newTraceLogger(true),
		mh:     new(messageHeader),
		sh:     new(segmentHeader),
		ph:     new(partHeader),
	}
}

const (
	productVersionMajor  = 4
	productVersionMinor  = 20
	protocolVersionMajor = 4
	protocolVersionMinor = 1
)

func (w *protocolWriter) writeProlog() error {
	req := &initRequest{}
	req.product.major = productVersionMajor
	req.product.minor = productVersionMinor
	req.protocol.major = protocolVersionMajor
	req.protocol.minor = protocolVersionMinor
	req.numOptions = 1
	req.endianess = littleEndian
	if err := req.encode(w.enc); err != nil {
		return err
	}
	w.tracer.Log(req)
	return w.wr.Flush()
}

func (w *protocolWriter) write(sessionID int64, messageType messageType, commit bool, writers ...partWriter) error {
	// check on session variables to be send as ClientInfo
	if messageType.clientInfoSupported() && w.sv.Version() != w.lastSVVersion {
		var upd map[string]string
		var del map[string]bool

		w.sv.WithRLock(func() {
			upd, del = w.sv.CompareWithRLock(w.lastSV)
			w.lastSVVersion = w.sv.Version()
			w.lastSV = w.sv.LoadWithRLock()
		})

		// TODO: how to delete session variables via clientInfo
		// ...for the time being we set the value to <space>...
		for k := range del {
			upd[k] = ""
		}
		writers = append([]partWriter{clientInfo(upd)}, writers...)
	}

	numWriters := len(writers)
	partSize := make([]int, numWriters)
	size := int64(segmentHeaderSize + numWriters*partHeaderSize) //int64 to hold MaxUInt32 in 32bit OS

	for i, part := range writers {
		s := part.size()
		size += int64(s + padBytes(s))
		partSize[i] = s // buffer size (expensive calculation)
	}

	if size > math.MaxUint32 {
		return fmt.Errorf("message size %d exceeds maximum message header value %d", size, int64(math.MaxUint32)) //int64: without cast overflow error in 32bit OS
	}

	bufferSize := size

	w.mh.sessionID = sessionID
	w.mh.varPartLength = uint32(size)
	w.mh.varPartSize = uint32(bufferSize)
	w.mh.noOfSegm = 1

	if err := w.mh.encode(w.enc); err != nil {
		return err
	}
	w.tracer.Log(w.mh)

	if size > math.MaxInt32 {
		return fmt.Errorf("message size %d exceeds maximum part header value %d", size, math.MaxInt32)
	}

	w.sh.messageType = messageType
	w.sh.commit = commit
	w.sh.segmentKind = skRequest
	w.sh.segmentLength = int32(size)
	w.sh.segmentOfs = 0
	w.sh.noOfParts = int16(numWriters)
	w.sh.segmentNo = 1

	if err := w.sh.encode(w.enc); err != nil {
		return err
	}
	w.tracer.Log(w.sh)

	bufferSize -= segmentHeaderSize

	for i, part := range writers {

		size := partSize[i]
		pad := padBytes(size)

		w.ph.partKind = part.kind()
		if err := w.ph.setNumArg(part.numArg()); err != nil {
			return err
		}
		w.ph.bufferLength = int32(size)
		w.ph.bufferSize = int32(bufferSize)

		if err := w.ph.encode(w.enc); err != nil {
			return err
		}
		w.tracer.Log(w.ph)

		if err := part.encode(w.enc); err != nil {
			return err
		}
		w.tracer.Log(part)

		w.enc.Zeroes(pad)

		bufferSize -= int64(partHeaderSize + size + pad)
	}
	return w.wr.Flush()
}
