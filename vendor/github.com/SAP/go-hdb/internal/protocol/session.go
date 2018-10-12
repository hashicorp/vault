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
	"context"
	"crypto/tls"
	"database/sql/driver"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"sync"
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
	addr     string
	timeout  time.Duration
	conn     net.Conn
	isBad    bool  // bad connection
	badError error // error cause for session bad state
	inTx     bool  // in transaction
}

func newSessionConn(ctx context.Context, addr string, timeoutSec int, config *tls.Config) (*sessionConn, error) {
	timeout := time.Duration(timeoutSec) * time.Second
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	// is TLS connection requested?
	if config != nil {
		conn = tls.Client(conn, config)
	}

	return &sessionConn{addr: addr, timeout: timeout, conn: conn}, nil
}

func (c *sessionConn) close() error {
	return c.conn.Close()
}

// Read implements the io.Reader interface.
func (c *sessionConn) Read(b []byte) (int, error) {
	//set timeout
	if err := c.conn.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
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
	if err := c.conn.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
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

type beforeRead func(p replyPart)

// session parameter
type sessionPrm interface {
	Host() string
	Username() string
	Password() string
	Locale() string
	FetchSize() int
	Timeout() int
	TLSConfig() *tls.Config
}

// Session represents a HDB session.
type Session struct {
	prm sessionPrm

	conn *sessionConn
	rd   *bufio.Reader
	wr   *bufio.Writer

	// reuse header
	mh *messageHeader
	sh *segmentHeader
	ph *partHeader

	//reuse request / reply parts
	scramsha256InitialRequest *scramsha256InitialRequest
	scramsha256InitialReply   *scramsha256InitialReply
	scramsha256FinalRequest   *scramsha256FinalRequest
	scramsha256FinalReply     *scramsha256FinalReply
	topologyInformation       *topologyInformation
	connectOptions            *connectOptions
	rowsAffected              *rowsAffected
	statementID               *statementID
	resultMetadata            *resultMetadata
	resultsetID               *resultsetID
	resultset                 *resultset
	parameterMetadata         *parameterMetadata
	outputParameters          *outputParameters
	writeLobRequest           *writeLobRequest
	readLobRequest            *readLobRequest
	writeLobReply             *writeLobReply
	readLobReply              *readLobReply

	//standard replies
	stmtCtx   *statementContext
	txFlags   *transactionFlags
	lastError *hdbErrors

	//serialize write request - read reply
	//supports calling session methods in go routines (driver methods with context cancellation)
	mu sync.Mutex
}

// NewSession creates a new database session.
func NewSession(ctx context.Context, prm sessionPrm) (*Session, error) {

	if trace {
		outLogger.Printf("%s", prm)
	}

	conn, err := newSessionConn(ctx, prm.Host(), prm.Timeout(), prm.TLSConfig())
	if err != nil {
		return nil, err
	}

	rd := bufio.NewReader(conn)
	wr := bufio.NewWriter(conn)

	s := &Session{
		prm:                       prm,
		conn:                      conn,
		rd:                        rd,
		wr:                        wr,
		mh:                        new(messageHeader),
		sh:                        new(segmentHeader),
		ph:                        new(partHeader),
		scramsha256InitialRequest: new(scramsha256InitialRequest),
		scramsha256InitialReply:   new(scramsha256InitialReply),
		scramsha256FinalRequest:   new(scramsha256FinalRequest),
		scramsha256FinalReply:     new(scramsha256FinalReply),
		topologyInformation:       newTopologyInformation(),
		connectOptions:            newConnectOptions(),
		rowsAffected:              new(rowsAffected),
		statementID:               new(statementID),
		resultMetadata:            new(resultMetadata),
		resultsetID:               new(resultsetID),
		resultset:                 new(resultset),
		parameterMetadata:         new(parameterMetadata),
		outputParameters:          new(outputParameters),
		writeLobRequest:           new(writeLobRequest),
		readLobRequest:            new(readLobRequest),
		writeLobReply:             new(writeLobReply),
		readLobReply:              new(readLobReply),
		stmtCtx:                   newStatementContext(),
		txFlags:                   newTransactionFlags(),
		lastError:                 new(hdbErrors),
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

	username := make([]byte, cesu8.StringSize(s.prm.Username()))
	if _, _, err := tr.Transform(username, []byte(s.prm.Username()), true); err != nil {
		return err // should never happen
	}

	password := make([]byte, cesu8.StringSize(s.prm.Password()))
	if _, _, err := tr.Transform(password, []byte(s.prm.Password()), true); err != nil {
		return err //should never happen
	}

	clientChallenge := clientChallenge()

	//initial request
	s.scramsha256InitialRequest.username = username
	s.scramsha256InitialRequest.clientChallenge = clientChallenge

	if err := s.writeRequest(mtAuthenticate, false, s.scramsha256InitialRequest); err != nil {
		return err
	}

	if err := s.readReply(nil); err != nil {
		return err
	}

	//final request
	s.scramsha256FinalRequest.username = username
	s.scramsha256FinalRequest.clientProof = clientProof(s.scramsha256InitialReply.salt, s.scramsha256InitialReply.serverChallenge, clientChallenge, password)

	s.scramsha256InitialReply = nil // !!! next time readReply uses FinalReply

	id := newClientID()

	co := newConnectOptions()
	co.set(coDistributionProtocolVersion, booleanType(false))
	co.set(coSelectForUpdateSupported, booleanType(false))
	co.set(coSplitBatchCommands, booleanType(true))
	// cannot use due to HDB protocol error with secondtime datatype
	//co.set(coDataFormatVersion2, dfvSPS06)
	co.set(coDataFormatVersion2, dfvBaseline)
	co.set(coCompleteArrayExecution, booleanType(true))
	if s.prm.Locale() != "" {
		co.set(coClientLocale, stringType(s.prm.Locale()))
	}
	co.set(coClientDistributionMode, cdmOff)
	// setting this option has no effect
	//co.set(coImplicitLobStreaming, booleanType(true))

	if err := s.writeRequest(mtConnect, false, s.scramsha256FinalRequest, id, co); err != nil {
		return err
	}

	if err := s.readReply(nil); err != nil {
		return err
	}

	return nil
}

// QueryDirect executes a query without query parameters.
func (s *Session) QueryDirect(query string) (uint64, *ResultFieldSet, *FieldValues, PartAttributes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.writeRequest(mtExecuteDirect, false, command(query)); err != nil {
		return 0, nil, nil, nil, err
	}

	var id uint64
	var resultFieldSet *ResultFieldSet
	fieldValues := newFieldValues()

	f := func(p replyPart) {

		switch p := p.(type) {

		case *resultsetID:
			p.id = &id
		case *resultMetadata:
			resultFieldSet = newResultFieldSet(p.numArg)
			p.resultFieldSet = resultFieldSet
		case *resultset:
			p.s = s
			p.resultFieldSet = resultFieldSet
			p.fieldValues = fieldValues
		}
	}

	if err := s.readReply(f); err != nil {
		return 0, nil, nil, nil, err
	}

	return id, resultFieldSet, fieldValues, s.ph.partAttributes, nil
}

// ExecDirect executes a sql statement without statement parameters.
func (s *Session) ExecDirect(query string) (driver.Result, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.writeRequest(mtExecuteDirect, !s.conn.inTx, command(query)); err != nil {
		return nil, err
	}

	if err := s.readReply(nil); err != nil {
		return nil, err
	}

	if s.sh.functionCode == fcDDL {
		return driver.ResultNoRows, nil
	}
	return driver.RowsAffected(s.rowsAffected.total()), nil
}

// Prepare prepares a sql statement.
func (s *Session) Prepare(query string) (QueryType, uint64, *ParameterFieldSet, *ResultFieldSet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.writeRequest(mtPrepare, false, command(query)); err != nil {
		return QtNone, 0, nil, nil, err
	}

	var id uint64
	var prmFieldSet *ParameterFieldSet
	var resultFieldSet *ResultFieldSet

	f := func(p replyPart) {

		switch p := p.(type) {

		case *statementID:
			p.id = &id
		case *parameterMetadata:
			prmFieldSet = newParameterFieldSet(p.numArg)
			p.prmFieldSet = prmFieldSet
		case *resultMetadata:
			resultFieldSet = newResultFieldSet(p.numArg)
			p.resultFieldSet = resultFieldSet
		}
	}

	if err := s.readReply(f); err != nil {
		return QtNone, 0, nil, nil, err
	}

	return s.sh.functionCode.queryType(), id, prmFieldSet, resultFieldSet, nil
}

// Exec executes a sql statement.
func (s *Session) Exec(id uint64, prmFieldSet *ParameterFieldSet, args []driver.NamedValue) (driver.Result, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.statementID.id = &id
	if err := s.writeRequest(mtExecute, !s.conn.inTx, s.statementID, newInputParameters(prmFieldSet.inputFields(), args)); err != nil {
		return nil, err
	}

	if err := s.readReply(nil); err != nil {
		return nil, err
	}

	var result driver.Result
	if s.sh.functionCode == fcDDL {
		result = driver.ResultNoRows
	} else {
		result = driver.RowsAffected(s.rowsAffected.total())
	}

	if err := s.writeLobStream(prmFieldSet, nil, args); err != nil {
		return nil, err
	}

	return result, nil
}

// DropStatementID releases the hdb statement handle.
func (s *Session) DropStatementID(id uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.statementID.id = &id
	if err := s.writeRequest(mtDropStatementID, false, s.statementID); err != nil {
		return err
	}

	if err := s.readReply(nil); err != nil {
		return err
	}

	return nil
}

// Call executes a stored procedure.
func (s *Session) Call(id uint64, prmFieldSet *ParameterFieldSet, args []driver.NamedValue) (*FieldValues, []*TableResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.statementID.id = &id
	if err := s.writeRequest(mtExecute, false, s.statementID, newInputParameters(prmFieldSet.inputFields(), args)); err != nil {
		return nil, nil, err
	}

	prmFieldValues := newFieldValues()
	var tableResults []*TableResult
	var tableResult *TableResult

	f := func(p replyPart) {

		switch p := p.(type) {

		case *outputParameters:
			p.s = s
			p.outputFields = prmFieldSet.outputFields()
			p.fieldValues = prmFieldValues

		// table output parameters: meta, id, result (only first param?)
		case *resultMetadata:
			tableResult = newTableResult(s, p.numArg)
			tableResults = append(tableResults, tableResult)
			p.resultFieldSet = tableResult.resultFieldSet
		case *resultsetID:
			p.id = &(tableResult.id)
		case *resultset:
			p.s = s
			tableResult.attrs = s.ph.partAttributes
			p.resultFieldSet = tableResult.resultFieldSet
			p.fieldValues = tableResult.fieldValues
		}
	}

	if err := s.readReply(f); err != nil {
		return nil, nil, err
	}

	if err := s.writeLobStream(prmFieldSet, prmFieldValues, args); err != nil {
		return nil, nil, err
	}

	return prmFieldValues, tableResults, nil
}

// Query executes a query.
func (s *Session) Query(stmtID uint64, prmFieldSet *ParameterFieldSet, resultFieldSet *ResultFieldSet, args []driver.NamedValue) (uint64, *FieldValues, PartAttributes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.statementID.id = &stmtID
	if err := s.writeRequest(mtExecute, false, s.statementID, newInputParameters(prmFieldSet.inputFields(), args)); err != nil {
		return 0, nil, nil, err
	}

	var rsetID uint64
	fieldValues := newFieldValues()

	f := func(p replyPart) {

		switch p := p.(type) {

		case *resultsetID:
			p.id = &rsetID
		case *resultset:
			p.s = s
			p.resultFieldSet = resultFieldSet
			p.fieldValues = fieldValues
		}
	}

	if err := s.readReply(f); err != nil {
		return 0, nil, nil, err
	}

	return rsetID, fieldValues, s.ph.partAttributes, nil
}

// FetchNext fetches next chunk in query result set.
func (s *Session) FetchNext(id uint64, resultFieldSet *ResultFieldSet, fieldValues *FieldValues) (PartAttributes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.resultsetID.id = &id
	if err := s.writeRequest(mtFetchNext, false, s.resultsetID, fetchsize(s.prm.FetchSize())); err != nil {
		return nil, err
	}

	f := func(p replyPart) {

		switch p := p.(type) {

		case *resultset:
			p.s = s
			p.resultFieldSet = resultFieldSet
			p.fieldValues = fieldValues
		}
	}

	if err := s.readReply(f); err != nil {
		return nil, err
	}

	return s.ph.partAttributes, nil
}

// CloseResultsetID releases the hdb resultset handle.
func (s *Session) CloseResultsetID(id uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.resultsetID.id = &id
	if err := s.writeRequest(mtCloseResultset, false, s.resultsetID); err != nil {
		return err
	}

	if err := s.readReply(nil); err != nil {
		return err
	}

	return nil
}

// Commit executes a database commit.
func (s *Session) Commit() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.writeRequest(mtCommit, false); err != nil {
		return err
	}

	if err := s.readReply(nil); err != nil {
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
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.writeRequest(mtRollback, false); err != nil {
		return err
	}

	if err := s.readReply(nil); err != nil {
		return err
	}

	if trace {
		outLogger.Printf("transaction flags: %s", s.txFlags)
	}

	s.conn.inTx = false
	return nil
}

//

func (s *Session) readLobStream(w lobChunkWriter) error {

	s.readLobRequest.w = w
	s.readLobReply.w = w

	for !w.eof() {

		if err := s.writeRequest(mtWriteLob, false, s.readLobRequest); err != nil {
			return err
		}
		if err := s.readReply(nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) writeLobStream(prmFieldSet *ParameterFieldSet, prmFieldValues *FieldValues, args []driver.NamedValue) error {

	if s.writeLobReply.numArg == 0 {
		return nil
	}

	lobPrmFields := make([]*ParameterField, s.writeLobReply.numArg)

	j := 0
	for _, f := range prmFieldSet.fields {
		if f.TypeCode().isLob() && f.In() && f.chunkReader != nil {
			f.lobLocatorID = s.writeLobReply.ids[j]
			lobPrmFields[j] = f
			j++
		}
	}
	if j != s.writeLobReply.numArg {
		return fmt.Errorf("protocol error: invalid number of lob parameter ids %d - expected %d", j, s.writeLobReply.numArg)
	}

	s.writeLobRequest.lobPrmFields = lobPrmFields

	f := func(p replyPart) {
		if p, ok := p.(*outputParameters); ok {
			p.s = s
			p.outputFields = prmFieldSet.outputFields()
			p.fieldValues = prmFieldValues
		}
	}

	for s.writeLobReply.numArg != 0 {
		if err := s.writeRequest(mtReadLob, false, s.writeLobRequest); err != nil {
			return err
		}

		if err := s.readReply(f); err != nil {
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
	req.endianess = littleEndian
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

		s.wr.WriteZeroes(pad)

		bufferSize -= int64(partHeaderSize + size + pad)

	}

	return s.wr.Flush()

}

func (s *Session) readReply(beforeRead beforeRead) error {

	replyRowsAffected := false
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
	lastPart := noOfParts - 1

	for i := 0; i < noOfParts; i++ {

		if err := s.ph.read(s.rd); err != nil {
			return err
		}

		numArg := int(s.ph.argumentCount)

		var part replyPart

		switch s.ph.partKind {

		case pkAuthentication:
			if s.scramsha256InitialReply != nil { // first call: initial reply
				part = s.scramsha256InitialReply
			} else { // second call: final reply
				part = s.scramsha256FinalReply
			}
		case pkTopologyInformation:
			part = s.topologyInformation
		case pkConnectOptions:
			part = s.connectOptions
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
			replyRowsAffected = true
			part = s.rowsAffected
		case pkReadLobReply:
			part = s.readLobReply
		case pkWriteLobReply:
			part = s.writeLobReply
		default:
			return fmt.Errorf("read not expected part kind %s", s.ph.partKind)
		}

		part.setNumArg(numArg)

		if beforeRead != nil {
			beforeRead(part)
		}

		s.rd.ResetCnt()
		if err := part.read(s.rd); err != nil {
			return err
		}
		cnt := s.rd.Cnt()

		if cnt != int(s.ph.bufferLength) {
			outLogger.Printf("+++ partLenght: %d - not equal read byte amount: %d", s.ph.bufferLength, cnt)
		}

		if i != lastPart { // not last part
			s.rd.Skip(padBytes(int(s.ph.bufferLength)))
		}
	}

	// last part
	// TODO: workaround (see *)
	if diff == 0 {
		s.rd.Skip(padBytes(int(s.ph.bufferLength)))
	}

	if err := s.rd.GetError(); err != nil {
		return err
	}

	if replyError {
		if replyRowsAffected { //link statement to error
			j := 0
			for i, rows := range s.rowsAffected.rows {
				if rows == raExecutionFailed {
					s.lastError.setStmtNo(j, i)
					j++
				}
			}
		}
		if s.lastError.isWarnings() {
			for _, _error := range s.lastError.errors {
				sqltrace.Traceln(_error)
			}
			return nil
		}
		return s.lastError
	}
	return nil
}
