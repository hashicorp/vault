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

package protocol

import (
	"database/sql/driver"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"

	"github.com/SAP/go-hdb/internal/bufio"
	"github.com/SAP/go-hdb/internal/unicode"
	"github.com/SAP/go-hdb/internal/unicode/cesu8"

	"github.com/SAP/go-hdb/driver/sqltrace"
)

const (
	mnSCRAMSHA256 = "SCRAMSHA256"
	mnGSS         = "GSS"
	mnSAML        = "SAML"
)

var trace bool

func init() {
	flag.BoolVar(&trace, "hdb.protocol.trace", false, "enabling hdb protocol trace")
}

var (
	outLogger = log.New(os.Stdout, "hdb.protocol ", log.Ldate|log.Ltime|log.Lshortfile)
	errLogger = log.New(os.Stderr, "hdb.protocol ", log.Ldate|log.Ltime|log.Lshortfile)
)

//padding
const (
	padding = 8
)

func padBytes(size int) int {
	if r := size % padding; r != 0 {
		return padding - r
	}
	return 0
}

// SessionConn wraps the database tcp connection. It sets timeouts and handles driver ErrBadConn behavior.
type sessionConn struct {
	addr            string
	timeoutDuration time.Duration
	conn            net.Conn
	isBad           bool  // bad connection
	badError        error // error cause for session bad state
	inTx            bool  // in transaction
}

func newSessionConn(addr string, timeout int) (*sessionConn, error) {
	timeoutDuration := time.Duration(timeout) * time.Second
	conn, err := net.DialTimeout("tcp", addr, timeoutDuration)
	if err != nil {
		return nil, err
	}

	return &sessionConn{
		addr:            addr,
		timeoutDuration: timeoutDuration,
		conn:            conn,
	}, nil
}

func (c *sessionConn) close() error {
	return c.conn.Close()
}

// Read implements the io.Reader interface.
func (c *sessionConn) Read(b []byte) (int, error) {
	//set timeout
	if err := c.conn.SetReadDeadline(time.Now().Add(c.timeoutDuration)); err != nil {
		return 0, err
	}
	n, err := c.conn.Read(b)
	if err != nil {
		errLogger.Printf("Connection read error local address %s remote address %s: %s", c.conn.LocalAddr(), c.conn.RemoteAddr(), err)
		c.isBad = true
		c.badError = err
		return n, driver.ErrBadConn
	}
	return n, nil
}

// Write implements the io.Writer interface.
func (c *sessionConn) Write(b []byte) (int, error) {
	//set timeout
	if err := c.conn.SetWriteDeadline(time.Now().Add(c.timeoutDuration)); err != nil {
		return 0, err
	}
	n, err := c.conn.Write(b)
	if err != nil {
		errLogger.Printf("Connection write error local address %s remote address %s: %s", c.conn.LocalAddr(), c.conn.RemoteAddr(), err)
		c.isBad = true
		c.badError = err
		return n, driver.ErrBadConn
	}
	return n, nil
}

type providePart func(pk partKind) replyPart
type beforeRead func(p replyPart)

// Session represents a HDB session.
type Session struct {
	prm *SessionPrm

	conn *sessionConn
	rd   *bufio.Reader
	wr   *bufio.Writer

	// reuse header
	mh *messageHeader
	sh *segmentHeader
	ph *partHeader

	//reuse request / reply parts
	rowsAffected      *rowsAffected
	statementID       *statementID
	resultMetadata    *resultMetadata
	resultsetID       *resultsetID
	resultset         *resultset
	parameterMetadata *parameterMetadata
	outputParameters  *outputParameters
	readLobRequest    *readLobRequest
	readLobReply      *readLobReply

	//standard replies
	stmtCtx   *statementContext
	txFlags   *transactionFlags
	lastError *hdbError
}

// NewSession creates a new database session.
func NewSession(prm *SessionPrm) (*Session, error) {

	if trace {
		outLogger.Printf("%s", prm)
	}

	conn, err := newSessionConn(prm.Host, prm.Timeout)
	if err != nil {
		return nil, err
	}

	var rd *bufio.Reader
	var wr *bufio.Writer
	if prm.BufferSize > 0 {
		rd = bufio.NewReaderSize(conn, prm.BufferSize)
		wr = bufio.NewWriterSize(conn, prm.BufferSize)
	} else {
		rd = bufio.NewReader(conn)
		wr = bufio.NewWriter(conn)
	}

	s := &Session{
		prm:               prm,
		conn:              conn,
		rd:                rd,
		wr:                wr,
		mh:                new(messageHeader),
		sh:                new(segmentHeader),
		ph:                new(partHeader),
		rowsAffected:      new(rowsAffected),
		statementID:       new(statementID),
		resultMetadata:    new(resultMetadata),
		resultsetID:       new(resultsetID),
		resultset:         new(resultset),
		parameterMetadata: new(parameterMetadata),
		outputParameters:  new(outputParameters),
		readLobRequest:    new(readLobRequest),
		readLobReply:      new(readLobReply),
		stmtCtx:           newStatementContext(),
		txFlags:           newTransactionFlags(),
		lastError:         newHdbError(),
	}

	if err = s.init(); err != nil {
		return nil, err
	}

	return s, nil
}

// Close closes the session.
func (s *Session) Close() error {
	return s.conn.close()
}

func (s *Session) sessionID() int64 {
	return s.mh.sessionID
}

// InTx indicates, if the session is in transaction mode.
func (s *Session) InTx() bool {
	return s.conn.inTx
}

// SetInTx sets session in transaction mode.
func (s *Session) SetInTx(v bool) {
	s.conn.inTx = v
}

// IsBad indicates, that the session is in bad state.
func (s *Session) IsBad() bool {
	return s.conn.isBad
}

// BadErr returns the error, that caused the bad session state.
func (s *Session) BadErr() error {
	return s.conn.badError
}

func (s *Session) init() error {

	if err := s.initRequest(); err != nil {
		return err
	}

	// TODO: detect authentication method
	// - actually only basic authetication supported

	authentication := mnSCRAMSHA256

	switch authentication {
	default:
		return fmt.Errorf("invalid authentication %s", authentication)

	case mnSCRAMSHA256:
		if err := s.authenticateScramsha256(); err != nil {
			return err
		}
	case mnGSS:
		panic("not implemented error")
	case mnSAML:
		panic("not implemented error")
	}

	id := s.sessionID()
	if id <= 0 {
		return fmt.Errorf("invalid session id %d", id)
	}

	if trace {
		outLogger.Printf("sessionId %d", id)
	}

	return nil
}

func (s *Session) authenticateScramsha256() error {
	tr := unicode.Utf8ToCesu8Transformer
	tr.Reset()

	username := make([]byte, cesu8.StringSize(s.prm.Username))
	if _, _, err := tr.Transform(username, []byte(s.prm.Username), true); err != nil {
		return err // should never happen
	}

	password := make([]byte, cesu8.StringSize(s.prm.Password))
	if _, _, err := tr.Transform(password, []byte(s.prm.Password), true); err != nil {
		return err //should never happen
	}

	clientChallenge := clientChallenge()

	//initial request
	ireq := newScramsha256InitialRequest()
	ireq.username = username
	ireq.clientChallenge = clientChallenge

	if err := s.writeRequest(mtAuthenticate, false, ireq); err != nil {
		return err
	}

	irep := newScramsha256InitialReply()

	f := func(pk partKind) replyPart {
		switch pk {
		case pkAuthentication:
			return irep
		default:
			return nil
		}
	}

	if err := s.readReply(f, nil); err != nil {
		return err
	}

	//final request
	freq := newScramsha256FinalRequest()
	freq.username = username
	freq.clientProof = clientProof(irep.salt, irep.serverChallenge, clientChallenge, password)

	id := newClientID()

	co := newConnectOptions()
	co.set(coDistributionProtocolVersion, booleanType(false))
	co.set(coSelectForUpdateSupported, booleanType(false))
	co.set(coSplitBatchCommands, booleanType(true))
	co.set(coDataFormatVersion, dfvBaseline)
	co.set(coDataFormatVersion2, dfvBaseline)
	co.set(coCompleteArrayExecution, booleanType(true))
	co.set(coClientLocale, stringType(s.prm.Locale))
	co.set(coClientDistributionMode, cdmOff)

	if err := s.writeRequest(mtConnect, false, freq, id, co); err != nil {
		return err
	}

	frep := newScramsha256FinalReply()
	topo := newTopologyOptions()

	f = func(pk partKind) replyPart {
		switch pk {
		case pkAuthentication:
			return frep
		case pkTopologyInformation:
			return topo
		case pkConnectOptions:
			return co
		default:
			return nil
		}
	}

	if err := s.readReply(f, nil); err != nil {
		return err
	}

	return nil
}

// QueryDirect executes a query without query parameters.
func (s *Session) QueryDirect(query string) (uint64, *FieldSet, *FieldValues, PartAttributes, error) {

	if err := s.writeRequest(mtExecuteDirect, false, command(query)); err != nil {
		return 0, nil, nil, nil, err
	}

	var id uint64
	var fieldSet *FieldSet
	fieldValues := newFieldValues(s)

	f := func(p replyPart) {

		switch p := p.(type) {

		case *resultsetID:
			p.id = &id
		case *resultMetadata:
			fieldSet = newFieldSet(p.numArg)
			p.fieldSet = fieldSet
		case *resultset:
			p.fieldSet = fieldSet
			p.fieldValues = fieldValues
		}
	}

	if err := s.readReply(nil, f); err != nil {
		return 0, nil, nil, nil, err
	}

	attrs := s.ph.partAttributes

	return id, fieldSet, fieldValues, attrs, nil
}

// ExecDirect executes a sql statement without statement parameters.
func (s *Session) ExecDirect(query string) (driver.Result, error) {

	if err := s.writeRequest(mtExecuteDirect, !s.conn.inTx, command(query)); err != nil {
		return nil, err
	}

	if err := s.readReply(nil, nil); err != nil {
		return nil, err
	}

	if s.sh.functionCode == fcDDL {
		return driver.ResultNoRows, nil
	}
	return driver.RowsAffected(s.rowsAffected.total()), nil
}

// Prepare prepares a sql statement.
func (s *Session) Prepare(query string) (QueryType, uint64, *FieldSet, *FieldSet, error) {

	if err := s.writeRequest(mtPrepare, false, command(query)); err != nil {
		return QtNone, 0, nil, nil, err
	}

	var id uint64
	var prmFieldSet *FieldSet
	var resultFieldSet *FieldSet

	f := func(p replyPart) {

		switch p := p.(type) {

		case *statementID:
			p.id = &id
		case *parameterMetadata:
			prmFieldSet = newFieldSet(p.numArg)
			p.fieldSet = prmFieldSet
		case *resultMetadata:
			resultFieldSet = newFieldSet(p.numArg)
			p.fieldSet = resultFieldSet
		}
	}

	if err := s.readReply(nil, f); err != nil {
		return QtNone, 0, nil, nil, err
	}

	return s.sh.functionCode.queryType(), id, prmFieldSet, resultFieldSet, nil
}

// Exec executes a sql statement.
func (s *Session) Exec(id uint64, parameterFieldSet *FieldSet, args []driver.Value) (driver.Result, error) {

	s.statementID.id = &id
	if err := s.writeRequest(mtExecute, !s.conn.inTx, s.statementID, newParameters(parameterFieldSet, args)); err != nil {
		return nil, err
	}

	wlr := newWriteLobReply() //lob streaming

	f := func(pk partKind) replyPart {
		switch pk {
		case pkWriteLobReply:
			return wlr
		default:
			return nil
		}
	}

	if err := s.readReply(f, nil); err != nil {
		return nil, err
	}

	var result driver.Result
	if s.sh.functionCode == fcDDL {
		result = driver.ResultNoRows
	} else {
		result = driver.RowsAffected(s.rowsAffected.total())
	}

	if wlr.numArg > 0 {
		if err := s.writeLobStream(parameterFieldSet, nil, args, wlr); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// DropStatementID releases the hdb statement handle.
func (s *Session) DropStatementID(id uint64) error {

	s.statementID.id = &id
	if err := s.writeRequest(mtDropStatementID, false, s.statementID); err != nil {
		return err
	}

	if err := s.readReply(nil, nil); err != nil {
		return err
	}

	return nil
}

// Call executes a stored procedure.
func (s *Session) Call(id uint64, prmFieldSet *FieldSet, args []driver.Value) (*FieldValues, []*TableResult, error) {

	s.statementID.id = &id
	if err := s.writeRequest(mtExecute, false, s.statementID, newParameters(prmFieldSet, args)); err != nil {
		return nil, nil, err
	}

	wlr := newWriteLobReply() //lob streaming

	f := func(pk partKind) replyPart {
		switch pk {
		case pkWriteLobReply:
			return wlr
		default:
			return nil
		}
	}

	prmFieldValues := newFieldValues(s)
	var tableResults []*TableResult
	var tableResult *TableResult

	g := func(p replyPart) {

		switch p := p.(type) {

		case *outputParameters:
			p.fieldSet = prmFieldSet
			p.fieldValues = prmFieldValues

		// table output parameters: meta, id, result (only first param?)
		case *resultMetadata:
			tableResult = newTableResult(s, p.numArg)
			tableResults = append(tableResults, tableResult)
			p.fieldSet = tableResult.fieldSet
		case *resultsetID:
			p.id = &(tableResult.id)
		case *resultset:
			tableResult.attrs = s.ph.partAttributes
			p.fieldSet = tableResult.fieldSet
			p.fieldValues = tableResult.fieldValues
		}
	}

	if err := s.readReply(f, g); err != nil {
		return nil, nil, err
	}

	if wlr.numArg > 0 {
		if err := s.writeLobStream(prmFieldSet, prmFieldValues, args, wlr); err != nil {
			return nil, nil, err
		}
	}

	return prmFieldValues, tableResults, nil
}

// Query executes a query.
func (s *Session) Query(stmtID uint64, parameterFieldSet *FieldSet, resultFieldSet *FieldSet, args []driver.Value) (uint64, *FieldValues, PartAttributes, error) {

	s.statementID.id = &stmtID
	if err := s.writeRequest(mtExecute, false, s.statementID, newParameters(parameterFieldSet, args)); err != nil {
		return 0, nil, nil, err
	}

	var rsetID uint64
	fieldValues := newFieldValues(s)

	f := func(p replyPart) {

		switch p := p.(type) {

		case *resultsetID:
			p.id = &rsetID
		case *resultset:
			p.fieldSet = resultFieldSet
			p.fieldValues = fieldValues
		}
	}

	if err := s.readReply(nil, f); err != nil {
		return 0, nil, nil, err
	}

	attrs := s.ph.partAttributes

	return rsetID, fieldValues, attrs, nil
}

// FetchNext fetches next chunk in query result set.
func (s *Session) FetchNext(id uint64, resultFieldSet *FieldSet) (*FieldValues, PartAttributes, error) {
	s.resultsetID.id = &id
	if err := s.writeRequest(mtFetchNext, false, s.resultsetID, fetchsize(s.prm.FetchSize)); err != nil {
		return nil, nil, err
	}

	fieldValues := newFieldValues(s)

	f := func(p replyPart) {

		switch p := p.(type) {

		case *resultset:
			p.fieldSet = resultFieldSet
			p.fieldValues = fieldValues
		}
	}

	if err := s.readReply(nil, f); err != nil {
		return nil, nil, err
	}

	attrs := s.ph.partAttributes

	return fieldValues, attrs, nil
}

// CloseResultsetID releases the hdb resultset handle.
func (s *Session) CloseResultsetID(id uint64) error {

	s.resultsetID.id = &id
	if err := s.writeRequest(mtCloseResultset, false, s.resultsetID); err != nil {
		return err
	}

	if err := s.readReply(nil, nil); err != nil {
		return err
	}

	return nil
}

// Commit executes a database commit.
func (s *Session) Commit() error {

	if err := s.writeRequest(mtCommit, false); err != nil {
		return err
	}

	if err := s.readReply(nil, nil); err != nil {
		return err
	}

	if trace {
		outLogger.Printf("transaction flags: %s", s.txFlags)
	}

	s.conn.inTx = false
	return nil
}

// Rollback executes a database rollback.
func (s *Session) Rollback() error {

	if err := s.writeRequest(mtRollback, false); err != nil {
		return err
	}

	if err := s.readReply(nil, nil); err != nil {
		return err
	}

	if trace {
		outLogger.Printf("transaction flags: %s", s.txFlags)
	}

	s.conn.inTx = false
	return nil
}

// helper
func readLobStreamDone(writers []lobWriter) bool {
	for _, writer := range writers {
		if !writer.eof() {
			return false
		}
	}
	return true
}

//

func (s *Session) readLobStream(writers []lobWriter) error {

	f := func(pk partKind) replyPart {
		switch pk {
		case pkReadLobReply:
			return s.readLobReply
		default:
			return nil
		}
	}

	for !readLobStreamDone(writers) {

		s.readLobRequest.writers = writers
		s.readLobReply.writers = writers

		if err := s.writeRequest(mtWriteLob, false, s.readLobRequest); err != nil {
			return err
		}
		if err := s.readReply(f, nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) writeLobStream(prmFieldSet *FieldSet, prmFieldValues *FieldValues, args []driver.Value, reply *writeLobReply) error {

	num := reply.numArg
	readers := make([]lobReader, num)

	request := newWriteLobRequest(readers)

	j := 0
	for i, field := range prmFieldSet.fields {

		if field.typeCode().isLob() && field.in() {

			ptr, ok := args[i].(int64)
			if !ok {
				return fmt.Errorf("protocol error: invalid lob driver value type %T", args[i])
			}

			descr := pointerToLobWriteDescr(ptr)
			if descr.r == nil {
				return fmt.Errorf("protocol error: lob reader %d initial", ptr)
			}

			if j >= num {
				return fmt.Errorf("protocol error: invalid number of lob parameter ids %d", num)
			}

			if field.typeCode().isCharBased() {
				readers[j] = newCharLobReader(descr.r, reply.ids[j])
			} else {
				readers[j] = newBinaryLobReader(descr.r, reply.ids[j])
			}

			j++
		}
	}

	f := func(pk partKind) replyPart {
		switch pk {
		case pkWriteLobReply:
			return reply
		default:
			return nil
		}
	}

	g := func(p replyPart) {
		if p, ok := p.(*outputParameters); ok {
			p.fieldSet = prmFieldSet
			p.fieldValues = prmFieldValues
		}
	}

	for reply.numArg != 0 {
		if err := s.writeRequest(mtReadLob, false, request); err != nil {
			return err
		}

		if err := s.readReply(f, g); err != nil {
			return err
		}
	}

	return nil
}

//

func (s *Session) initRequest() error {

	// init
	s.mh.sessionID = -1

	// handshake
	req := newInitRequest()
	// TODO: constants
	req.product.major = 4
	req.product.minor = 20
	req.protocol.major = 4
	req.protocol.minor = 1
	req.numOptions = 1
	req.endianess = archEndian
	if err := req.write(s.wr); err != nil {
		return err
	}

	rep := newInitReply()
	if err := rep.read(s.rd); err != nil {
		return err
	}
	return nil
}

func (s *Session) writeRequest(messageType messageType, commit bool, requests ...requestPart) error {

	partSize := make([]int, len(requests))

	size := int64(segmentHeaderSize + len(requests)*partHeaderSize) //int64 to hold MaxUInt32 in 32bit OS

	for i, part := range requests {
		s, err := part.size()
		if err != nil {
			return err
		}
		size += int64(s + padBytes(s))
		partSize[i] = s // buffer size (expensive calculation)
	}

	if size > math.MaxUint32 {
		return fmt.Errorf("message size %d exceeds maximum message header value %d", size, int64(math.MaxUint32)) //int64: without cast overflow error in 32bit OS
	}

	bufferSize := size

	s.mh.varPartLength = uint32(size)
	s.mh.varPartSize = uint32(bufferSize)
	s.mh.noOfSegm = 1

	if err := s.mh.write(s.wr); err != nil {
		return err
	}

	if size > math.MaxInt32 {
		return fmt.Errorf("message size %d exceeds maximum part header value %d", size, math.MaxInt32)
	}

	s.sh.messageType = messageType
	s.sh.commit = commit
	s.sh.segmentKind = skRequest
	s.sh.segmentLength = int32(size)
	s.sh.segmentOfs = 0
	s.sh.noOfParts = int16(len(requests))
	s.sh.segmentNo = 1

	if err := s.sh.write(s.wr); err != nil {
		return err
	}

	bufferSize -= segmentHeaderSize

	for i, part := range requests {

		size := partSize[i]
		pad := padBytes(size)

		s.ph.partKind = part.kind()
		numArg := part.numArg()
		switch {
		default:
			return fmt.Errorf("maximum number of arguments %d exceeded", numArg)
		case numArg <= math.MaxInt16:
			s.ph.argumentCount = int16(numArg)
			s.ph.bigArgumentCount = 0

		// TODO: seems not to work: see bulk insert test
		case numArg <= math.MaxInt32:
			s.ph.argumentCount = 0
			s.ph.bigArgumentCount = int32(numArg)
		}

		s.ph.bufferLength = int32(size)
		s.ph.bufferSize = int32(bufferSize)

		if err := s.ph.write(s.wr); err != nil {
			return err
		}

		if err := part.write(s.wr); err != nil {
			return err
		}

		if err := s.wr.WriteZeroes(pad); err != nil {
			return err
		}

		bufferSize -= int64(partHeaderSize + size + pad)

	}

	if err := s.wr.Flush(); err != nil {
		return err
	}

	return nil
}

func (s *Session) readReply(providePart providePart, beforeRead beforeRead) error {

	replyError := false

	if err := s.mh.read(s.rd); err != nil {
		return err
	}
	if s.mh.noOfSegm != 1 {
		return fmt.Errorf("simple message: no of segments %d - expected 1", s.mh.noOfSegm)
	}
	if err := s.sh.read(s.rd); err != nil {
		return err
	}

	// TODO: protocol error (sps 82)?: message header varPartLength < segment header segmentLength (*1)
	diff := int(s.mh.varPartLength) - int(s.sh.segmentLength)
	if trace && diff != 0 {
		outLogger.Printf("+++++diff %d", diff)
	}

	noOfParts := int(s.sh.noOfParts)

	for i := 0; i < noOfParts; i++ {

		if err := s.ph.read(s.rd); err != nil {
			return err
		}

		numArg := int(s.ph.argumentCount)

		var part replyPart

		if providePart != nil {
			part = providePart(s.ph.partKind)
		} else {
			part = nil
		}

		if part == nil { // use pre defined parts

			switch s.ph.partKind {

			case pkStatementID:
				part = s.statementID
			case pkResultMetadata:
				part = s.resultMetadata
			case pkResultsetID:
				part = s.resultsetID
			case pkResultset:
				part = s.resultset
			case pkParameterMetadata:
				part = s.parameterMetadata
			case pkOutputParameters:
				part = s.outputParameters
			case pkError:
				replyError = true
				part = s.lastError
			case pkStatementContext:
				part = s.stmtCtx
			case pkTransactionFlags:
				part = s.txFlags
			case pkRowsAffected:
				part = s.rowsAffected
			default:
				return fmt.Errorf("read not expected part kind %s", s.ph.partKind)
			}
		}

		part.setNumArg(numArg)

		if beforeRead != nil {
			beforeRead(part)
		}

		if err := part.read(s.rd); err != nil {
			return err
		}

		// TODO: workaround (see *)
		if i != (noOfParts-1) || (i == (noOfParts-1) && diff == 0) {
			if err := s.rd.Skip(padBytes(int(s.ph.bufferLength))); err != nil {
				return err
			}
		}
	}

	if replyError {
		if s.lastError.IsWarning() {
			sqltrace.Traceln(s.lastError)
		} else {
			return s.lastError
		}
	}
	return nil
}
