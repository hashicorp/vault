// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"bufio"
	"context"
	"database/sql/driver"
	"fmt"
	"io"

	"golang.org/x/text/transform"

	"github.com/SAP/go-hdb/driver/common"
	"github.com/SAP/go-hdb/internal/unicode"
	"github.com/SAP/go-hdb/internal/unicode/cesu8"
)

//padding
const padding = 8

func padBytes(size int) int {
	if r := size % padding; r != 0 {
		return padding - r
	}
	return 0
}

const (
	dfvLevel1        = 1
	defaultSessionID = -1
)

// Session represents a HDB session.
type Session struct {
	cfg *SessionConfig

	sessionID     int64
	serverOptions connectOptions
	serverVersion *common.HDBVersion

	pr *protocolReader
	pw *protocolWriter
}

// NewSession creates a new database session.
func NewSession(ctx context.Context, rw *bufio.ReadWriter, cfg *SessionConfig) (*Session, error) {
	pw := newProtocolWriter(rw.Writer, cfg.SessionVariables) // write upstream
	if err := pw.writeProlog(); err != nil {
		return nil, err
	}

	pr := newProtocolReader(false, rw.Reader) // read downstream
	if err := pr.readProlog(); err != nil {
		return nil, err
	}

	s := &Session{cfg: cfg, sessionID: defaultSessionID, pr: pr, pw: pw}

	authStepper := newAuth(cfg.Username, cfg.Password)
	var err error
	if s.sessionID, s.serverOptions, err = s.authenticate(authStepper); err != nil {
		return nil, err
	}

	if s.sessionID <= 0 {
		return nil, fmt.Errorf("invalid session id %d", s.sessionID)
	}

	s.serverVersion = common.ParseHDBVersion(s.serverOptions.fullVersionString())
	return s, nil
}

// SessionID returns the session id of the hdb connection.
func (s *Session) SessionID() int64 { return s.sessionID }

// ServerInfo returnsinformation reported by hdb server.
func (s *Session) ServerInfo() *common.ServerInfo {
	return &common.ServerInfo{
		Version: s.serverVersion,
	}
}

func (s *Session) defaultClientOptions() connectOptions {
	co := connectOptions{
		int8(coDistributionProtocolVersion): optBooleanType(false),
		int8(coSelectForUpdateSupported):    optBooleanType(false),
		int8(coSplitBatchCommands):          optBooleanType(true),
		int8(coDataFormatVersion2):          optIntType(s.cfg.Dfv),
		int8(coCompleteArrayExecution):      optBooleanType(true),
		int8(coClientDistributionMode):      cdmOff,
		// int8(coImplicitLobStreaming):        optBooleanType(true),
	}
	if s.cfg.Locale != "" {
		co[int8(coClientLocale)] = optStringType(s.cfg.Locale)
	}
	return co
}

func (s *Session) authenticate(stepper authStepper) (int64, connectOptions, error) {
	var auth partReadWriter
	var err error

	// client context
	clientContext := clientContext(plainOptions{
		int8(ccoClientVersion):            optStringType(s.cfg.DriverVersion),
		int8(ccoClientType):               optStringType(s.cfg.DriverName),
		int8(ccoClientApplicationProgram): optStringType(s.cfg.ApplicationName),
	})

	if auth, err = stepper.next(); err != nil {
		return 0, nil, err
	}
	if err := s.pw.write(s.sessionID, mtAuthenticate, false, clientContext, auth); err != nil {
		return 0, nil, err
	}

	if auth, err = stepper.next(); err != nil {
		return 0, nil, err
	}
	if err := s.pr.iterateParts(func(ph *partHeader) {
		if ph.partKind == pkAuthentication {
			s.pr.read(auth)
		}
	}); err != nil {
		return 0, nil, err
	}

	if auth, err = stepper.next(); err != nil {
		return 0, nil, err
	}
	id := newClientID()
	co := s.defaultClientOptions()
	if err := s.pw.write(s.sessionID, mtConnect, false, auth, id, co); err != nil {
		return 0, nil, err
	}

	if auth, err = stepper.next(); err != nil {
		return 0, nil, err
	}
	if err := s.pr.iterateParts(func(ph *partHeader) {
		switch ph.partKind {
		case pkAuthentication:
			s.pr.read(auth)
		case pkConnectOptions:
			s.pr.read(&co)
			// set data format version
			// TODO generalize for sniffer
			s.pr.setDfv(int(co[int8(coDataFormatVersion2)].(optIntType)))
		}
	}); err != nil {
		return 0, nil, err
	}

	return s.pr.sessionID(), co, nil
}

// QueryDirect executes a query without query parameters.
func (s *Session) QueryDirect(query string, commit bool) (driver.Rows, error) {
	// allow e.g inserts as query -> handle commit like in ExecDirect
	if err := s.pw.write(s.sessionID, mtExecuteDirect, commit, command(query)); err != nil {
		return nil, err
	}

	qr := &queryResult{session: s}
	meta := &resultMetadata{}
	resSet := &resultset{}

	if err := s.pr.iterateParts(func(ph *partHeader) {
		switch ph.partKind {
		case pkResultMetadata:
			s.pr.read(meta)
			qr.fields = meta.resultFields
		case pkResultsetID:
			s.pr.read((*resultsetID)(&qr.rsID))
		case pkResultset:
			resSet.resultFields = qr.fields
			s.pr.read(resSet)
			qr.fieldValues = resSet.fieldValues
			qr.attributes = ph.partAttributes
		}
	}); err != nil {
		return nil, err
	}
	if qr.rsID == 0 { // non select query
		return noResult, nil
	}
	return qr, nil
}

// ExecDirect executes a sql statement without statement parameters.
func (s *Session) ExecDirect(query string, commit bool) (driver.Result, error) {
	if err := s.pw.write(s.sessionID, mtExecuteDirect, commit, command(query)); err != nil {
		return nil, err
	}

	rows := &rowsAffected{}
	var numRow int64
	if err := s.pr.iterateParts(func(ph *partHeader) {
		if ph.partKind == pkRowsAffected {
			s.pr.read(rows)
			numRow = rows.total()
		}
	}); err != nil {
		return nil, err
	}
	if s.pr.functionCode() == fcDDL {
		return driver.ResultNoRows, nil
	}
	return driver.RowsAffected(numRow), nil
}

// Prepare prepares a sql statement.
func (s *Session) Prepare(query string) (*PrepareResult, error) {
	if err := s.pw.write(s.sessionID, mtPrepare, false, command(query)); err != nil {
		return nil, err
	}

	pr := &PrepareResult{}
	resMeta := &resultMetadata{}
	prmMeta := &parameterMetadata{}

	if err := s.pr.iterateParts(func(ph *partHeader) {
		switch ph.partKind {
		case pkStatementID:
			s.pr.read((*statementID)(&pr.stmtID))
		case pkResultMetadata:
			s.pr.read(resMeta)
			pr.resultFields = resMeta.resultFields
		case pkParameterMetadata:
			s.pr.read(prmMeta)
			pr.parameterFields = prmMeta.parameterFields
		}
	}); err != nil {
		return nil, err
	}
	pr.fc = s.pr.functionCode()
	return pr, nil
}

// Exec executes a sql statement.
func (s *Session) Exec(pr *PrepareResult, args []interface{}, commit bool) (driver.Result, error) {
	if err := s.pw.write(s.sessionID, mtExecute, commit, statementID(pr.stmtID), newInputParameters(pr.parameterFields, args)); err != nil {
		return nil, err
	}

	rows := &rowsAffected{}
	var ids []locatorID
	lobReply := &writeLobReply{}
	var numRow int64

	if err := s.pr.iterateParts(func(ph *partHeader) {
		switch ph.partKind {
		case pkRowsAffected:
			s.pr.read(rows)
			numRow = rows.total()
		case pkWriteLobReply:
			s.pr.read(lobReply)
			ids = lobReply.ids
		}
	}); err != nil {
		return nil, err
	}
	fc := s.pr.functionCode()

	if len(ids) != 0 {
		/*
			writeLobParameters:
			- chunkReaders
			- nil (no callResult, exec does not have output parameters)
		*/
		if err := s.encodeLobs(nil, ids, pr.parameterFields, args); err != nil {
			return nil, err
		}
	}

	if fc == fcDDL {
		return driver.ResultNoRows, nil
	}
	return driver.RowsAffected(numRow), nil
}

// QueryCall executes a stored procecure (by Query).
func (s *Session) QueryCall(pr *PrepareResult, args []interface{}) (driver.Rows, error) {
	/*
		only in args
		invariant: #inPrmFields == #args
	*/
	var inPrmFields, outPrmFields []*ParameterField
	for _, f := range pr.parameterFields {
		if f.In() {
			inPrmFields = append(inPrmFields, f)
		}
		if f.Out() {
			outPrmFields = append(outPrmFields, f)
		}
	}

	if err := s.pw.write(s.sessionID, mtExecute, false, statementID(pr.stmtID), newInputParameters(inPrmFields, args)); err != nil {
		return nil, err
	}

	/*
		call without lob input parameters:
		--> callResult output parameter values are set after read call
		call with lob input parameters:
		--> callResult output parameter values are set after last lob input write
	*/

	cr, ids, _, err := s.readCall(outPrmFields) // ignore numRow
	if err != nil {
		return nil, err
	}

	if len(ids) != 0 {
		/*
			writeLobParameters:
			- chunkReaders
			- cr (callResult output parameters are set after all lob input parameters are written)
		*/
		if err := s.encodeLobs(cr, ids, inPrmFields, args); err != nil {
			return nil, err
		}
	}

	// legacy mode?
	if s.cfg.Legacy {
		cr.appendTableRefFields() // TODO review
		for _, qr := range cr.qrs {
			// add to cache
			QueryResultCache.set(qr.rsID, qr)
		}
	} else {
		cr.appendTableRowsFields()
	}
	return cr, nil
}

// ExecCall executes a stored procecure (by Exec).
func (s *Session) ExecCall(pr *PrepareResult, args []interface{}) (driver.Result, error) {
	/*
		in,- and output args
		invariant: #prmFields == #args
	*/
	var inPrmFields, outPrmFields []*ParameterField
	var inArgs, outArgs []interface{}
	for i, f := range pr.parameterFields {
		if f.In() {
			inPrmFields = append(inPrmFields, f)
			inArgs = append(inArgs, args[i])
		}
		if f.Out() {
			outPrmFields = append(outPrmFields, f)
			outArgs = append(outArgs, args[i])
		}
	}

	// TODO release v1.0.0 - assign output parameters
	if len(outPrmFields) != 0 {
		return nil, fmt.Errorf("stmt.Exec: support of output parameters not implemented yet")
	}

	if err := s.pw.write(s.sessionID, mtExecute, false, statementID(pr.stmtID), newInputParameters(inPrmFields, inArgs)); err != nil {
		return nil, err
	}

	/*
		call without lob input parameters:
		--> callResult output parameter values are set after read call
		call with lob output parameters:
		--> callResult output parameter values are set after last lob input write
	*/

	cr, ids, numRow, err := s.readCall(outPrmFields)
	if err != nil {
		return nil, err
	}

	if len(ids) != 0 {
		/*
			writeLobParameters:
			- chunkReaders
			- cr (callResult output parameters are set after all lob input parameters are written)
		*/
		if err := s.encodeLobs(cr, ids, inPrmFields, inArgs); err != nil {
			return nil, err
		}
	}
	return driver.RowsAffected(numRow), nil
}

func (s *Session) readCall(outputFields []*ParameterField) (*callResult, []locatorID, int64, error) {
	cr := &callResult{session: s, outputFields: outputFields}

	//var qrs []*QueryResult
	var qr *queryResult
	rows := &rowsAffected{}
	var ids []locatorID
	outPrms := &outputParameters{}
	meta := &resultMetadata{}
	resSet := &resultset{}
	lobReply := &writeLobReply{}
	var numRow int64

	if err := s.pr.iterateParts(func(ph *partHeader) {
		switch ph.partKind {
		case pkRowsAffected:
			s.pr.read(rows)
			numRow = rows.total()
		case pkOutputParameters:
			outPrms.outputFields = cr.outputFields
			s.pr.read(outPrms)
			cr.fieldValues = outPrms.fieldValues
		case pkResultMetadata:
			/*
				procedure call with table parameters does return metadata for each table
				sequence: metadata, resultsetID, resultset
				but:
				- resultset might not be provided for all tables
				- so, 'additional' query result is detected by new metadata part
			*/
			qr = &queryResult{session: s}
			cr.qrs = append(cr.qrs, qr)
			s.pr.read(meta)
			qr.fields = meta.resultFields
		case pkResultset:
			resSet.resultFields = qr.fields
			s.pr.read(resSet)
			qr.fieldValues = resSet.fieldValues
			qr.attributes = ph.partAttributes
		case pkResultsetID:
			s.pr.read((*resultsetID)(&qr.rsID))
		case pkWriteLobReply:
			s.pr.read(lobReply)
			ids = lobReply.ids
		}
	}); err != nil {
		return nil, nil, 0, err
	}
	return cr, ids, numRow, nil
}

// Query executes a query.
func (s *Session) Query(pr *PrepareResult, args []interface{}, commit bool) (driver.Rows, error) {
	// allow e.g inserts as query -> handle commit like in exec
	if err := s.pw.write(s.sessionID, mtExecute, commit, statementID(pr.stmtID), newInputParameters(pr.parameterFields, args)); err != nil {
		return nil, err
	}

	qr := &queryResult{session: s, fields: pr.resultFields}
	resSet := &resultset{}

	if err := s.pr.iterateParts(func(ph *partHeader) {
		switch ph.partKind {
		case pkResultsetID:
			s.pr.read((*resultsetID)(&qr.rsID))
		case pkResultset:
			resSet.resultFields = qr.fields
			s.pr.read(resSet)
			qr.fieldValues = resSet.fieldValues
			qr.attributes = ph.partAttributes
		}
	}); err != nil {
		return nil, err
	}
	if qr.rsID == 0 { // non select query
		return noResult, nil
	}
	return qr, nil
}

// FetchNext fetches next chunk in query result set.
func (s *Session) fetchNext(qr *queryResult) error {
	if err := s.pw.write(s.sessionID, mtFetchNext, false, resultsetID(qr.rsID), fetchsize(s.cfg.FetchSize)); err != nil {
		return err
	}

	resSet := &resultset{resultFields: qr.fields, fieldValues: qr.fieldValues} // reuse field values

	return s.pr.iterateParts(func(ph *partHeader) {
		if ph.partKind == pkResultset {
			s.pr.read(resSet)
			qr.attributes = ph.partAttributes
			qr.fieldValues = resSet.fieldValues
		}
	})
}

// DropStatementID releases the hdb statement handle.
func (s *Session) DropStatementID(id uint64) error {
	if err := s.pw.write(s.sessionID, mtDropStatementID, false, statementID(id)); err != nil {
		return err
	}
	return s.pr.readSkip()
}

// CloseResultsetID releases the hdb resultset handle.
func (s *Session) CloseResultsetID(id uint64) error {
	if err := s.pw.write(s.sessionID, mtCloseResultset, false, resultsetID(id)); err != nil {
		return err
	}
	return s.pr.readSkip()
}

// Commit executes a database commit.
func (s *Session) Commit() error {
	if err := s.pw.write(s.sessionID, mtCommit, false); err != nil {
		return err
	}
	if err := s.pr.readSkip(); err != nil {
		return err
	}
	return nil
}

// Rollback executes a database rollback.
func (s *Session) Rollback() error {
	if err := s.pw.write(s.sessionID, mtRollback, false); err != nil {
		return err
	}
	if err := s.pr.readSkip(); err != nil {
		return err
	}
	return nil
}

// Disconnect disconnects the session.
func (s *Session) Disconnect() error {
	if err := s.pw.write(s.sessionID, mtDisconnect, false); err != nil {
		return err
	}
	/*
		Do not read server reply as on slow connections the TCP/IP connection is closed (by Server)
		before the reply can be read completely.

		// if err := s.pr.readSkip(); err != nil {
		// 	return err
		// }

	*/
	return nil
}

// decodeLobs decodes (reads from db) output lob or result lob parameters.

// read lob reply
// - seems like readLobreply returns only a result for one lob - even if more then one is requested
// --> read single lobs
func (s *Session) decodeLobs(descr *lobOutDescr, wr io.Writer) error {
	var err error

	if descr.isCharBased {
		wrcl := transform.NewWriter(wr, unicode.Cesu8ToUtf8Transformer) // CESU8 transformer
		err = s._decodeLobs(descr, wrcl, func(b []byte) (int64, error) {
			// Caution: hdb counts 4 byte utf-8 encodings (cesu-8 6 bytes) as 2 (3 byte) chars
			numChars := int64(0)
			for len(b) > 0 {
				if !cesu8.FullRune(b) { //
					return 0, fmt.Errorf("lob chunk consists of incomplete CESU-8 runes")
				}
				_, size := cesu8.DecodeRune(b)
				b = b[size:]
				numChars++
				if size == cesu8.CESUMax {
					numChars++
				}
			}
			return numChars, nil
		})
	} else {
		err = s._decodeLobs(descr, wr, func(b []byte) (int64, error) { return int64(len(b)), nil })
	}

	if pw, ok := wr.(*io.PipeWriter); ok { // if the writer is a pipe-end -> close at the end
		if err != nil {
			pw.CloseWithError(err)
		} else {
			pw.Close()
		}
	}
	return err
}

func (s *Session) _decodeLobs(descr *lobOutDescr, wr io.Writer, countChars func(b []byte) (int64, error)) error {
	lobChunkSize := int64(s.cfg.LobChunkSize)

	chunkSize := func(numChar, ofs int64) int32 {
		chunkSize := numChar - ofs
		if chunkSize > lobChunkSize {
			return int32(lobChunkSize)
		}
		return int32(chunkSize)
	}

	if _, err := wr.Write(descr.b); err != nil {
		return err
	}

	lobRequest := &readLobRequest{}
	lobRequest.id = descr.id

	lobReply := &readLobReply{}

	eof := descr.opt.isLastData()

	ofs, err := countChars(descr.b)
	if err != nil {
		return err
	}

	for !eof {

		lobRequest.ofs += ofs
		lobRequest.chunkSize = chunkSize(descr.numChar, ofs)

		if err := s.pw.write(s.sessionID, mtWriteLob, false, lobRequest); err != nil {
			return err
		}

		if err := s.pr.iterateParts(func(ph *partHeader) {
			if ph.partKind == pkReadLobReply {
				s.pr.read(lobReply)
			}
		}); err != nil {
			return err
		}

		if lobReply.id != lobRequest.id {
			return fmt.Errorf("internal error: invalid lob locator %d - expected %d", lobReply.id, lobRequest.id)
		}

		if _, err := wr.Write(lobReply.b); err != nil {
			return err
		}

		ofs, err = countChars(lobReply.b)
		if err != nil {
			return err
		}
		eof = lobReply.opt.isLastData()
	}
	return nil
}

// encodeLobs encodes (write to db) input lob parameters.
func (s *Session) encodeLobs(cr *callResult, ids []locatorID, inPrmFields []*ParameterField, args []interface{}) error {
	chunkSize := s.cfg.LobChunkSize

	readers := make([]io.Reader, 0, len(ids))
	descrs := make([]*writeLobDescr, 0, len(ids))

	numInPrmField := len(inPrmFields)

	j := 0
	for i, arg := range args { // range over args (mass / bulk operation)
		f := inPrmFields[i%numInPrmField]
		if f.tc.isLob() {
			rd, ok := arg.(io.Reader)
			if !ok {
				return fmt.Errorf("protocol error: invalid lob parameter %[1]T %[1]v - io.Reader expected", arg)
			}
			if f.tc.isCharBased() {
				rd = transform.NewReader(rd, unicode.Utf8ToCesu8Transformer) // CESU8 transformer
			}
			if j >= len(ids) {
				return fmt.Errorf("protocol error: invalid number of lob parameter ids %d", len(ids))
			}
			readers = append(readers, rd)
			descrs = append(descrs, &writeLobDescr{id: ids[j]})
			j++
		}
	}

	writeLobRequest := &writeLobRequest{}

	for len(descrs) != 0 {

		if len(descrs) != len(ids) {
			return fmt.Errorf("protocol error: invalid number of lob parameter ids %d - expected %d", len(descrs), len(ids))
		}
		for i, descr := range descrs { // check if ids and descrs are in sync
			if descr.id != ids[i] {
				return fmt.Errorf("protocol error: lob parameter id mismatch %d - expected %d", descr.id, ids[i])
			}
		}

		// TODO check total size limit
		for i, descr := range descrs {
			descr.b = make([]byte, chunkSize)
			size, err := readers[i].Read(descr.b)
			descr.b = descr.b[:size]
			if err != nil && err != io.EOF {
				return err
			}
			descr.ofs = -1 //offset (-1 := append)
			descr.opt = loDataincluded
			if err == io.EOF {
				descr.opt |= loLastdata
			}
		}

		writeLobRequest.descrs = descrs

		if err := s.pw.write(s.sessionID, mtReadLob, false, writeLobRequest); err != nil {
			return err
		}

		lobReply := &writeLobReply{}
		outPrms := &outputParameters{}

		if err := s.pr.iterateParts(func(ph *partHeader) {
			switch ph.partKind {
			case pkOutputParameters:
				outPrms.outputFields = cr.outputFields
				s.pr.read(outPrms)
				cr.fieldValues = outPrms.fieldValues
			case pkWriteLobReply:
				s.pr.read(lobReply)
				ids = lobReply.ids
			}
		}); err != nil {
			return err
		}

		// remove done descr and readers
		j := 0
		for i, descr := range descrs {
			if !descr.opt.isLastData() {
				descrs[j] = descr
				readers[j] = readers[i]
				j++
			}
		}
		descrs = descrs[:j]
		readers = readers[:j]
	}
	return nil
}
