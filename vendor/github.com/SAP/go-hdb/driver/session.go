package driver

import (
	"bufio"
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	p "github.com/SAP/go-hdb/driver/internal/protocol"
	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
	"golang.org/x/text/transform"
)

// SessionUser provides the fields for a hdb 'connect' (switch user) statement.
type SessionUser struct {
	Username, Password string
	Schema             string
}

func (u *SessionUser) equal(cmp *SessionUser) bool {
	if cmp == nil {
		return false
	}
	return u.Username == cmp.Username && u.Password == cmp.Password
}

func (u *SessionUser) clone() *SessionUser {
	return &SessionUser{Username: u.Username, Password: u.Password}
}

// use unexported type to avoid key collisions.
type switchUserCtxKeyType struct{}

var switchUserCtxKey switchUserCtxKeyType

// WithUserSwitch can be used to switch a user on a new or an existing connection
// (see https://help.sap.com/docs/hana-cloud-database/sap-hana-cloud-sap-hana-database-sql-reference-guide/connect-statement-session-management).
func WithUserSwitch(ctx context.Context, u *SessionUser) context.Context {
	return context.WithValue(ctx, switchUserCtxKey, u)
}

type session struct {
	dbConn  dbConn
	metrics *metrics
	attrs   *connAttrs

	prd *p.Reader
	pwr *p.Writer

	cesu8Encoder transform.Transformer

	hdbVersion   *Version
	databaseName string

	user *SessionUser // session user

	inTx bool

	sqlTracer *sqlTracer
}

func newSession(ctx context.Context, host string, logger *slog.Logger, metrics *metrics, attrs *connAttrs, authHnd *p.AuthHnd) (*session, error) {
	dbConn, err := newDBConn(ctx, logger, host, metrics, attrs)
	if err != nil {
		return nil, err
	}

	rd := bufio.NewReaderSize(dbConn, attrs._bufferSize)
	wr := bufio.NewWriterSize(dbConn, attrs._bufferSize)

	cesu8Encoder := attrs._cesu8Encoder() // call function only once.

	dec := encoding.NewDecoder(rd, attrs._cesu8Decoder(), attrs._emptyDateAsNull)
	enc := encoding.NewEncoder(wr, cesu8Encoder)

	protTrace := protTrace.Load()

	prd := p.NewDBReader(dec, protTrace, logger, attrs._lobChunkSize)
	pwr := p.NewWriter(wr, enc, protTrace, logger, attrs._sessionVariables)

	// prolog
	if err := pwr.WriteProlog(ctx); err != nil {
		dbConn.Close()
		return nil, err
	}
	if err := prd.ReadProlog(ctx); err != nil {
		dbConn.Close()
		return nil, err
	}

	var sqlTracer *sqlTracer
	if sqlTrace.Load() {
		sqlTracer = newSQLTracer(logger, 0)
	}
	s := &session{dbConn: dbConn, metrics: metrics, attrs: attrs, prd: prd, pwr: pwr, cesu8Encoder: cesu8Encoder, sqlTracer: sqlTracer}

	if authHnd != nil { // authenticate
		serverOptions, err := s.authenticate(ctx, authHnd, attrs)
		if err != nil {
			dbConn.Close()
			return nil, err
		}
		s.hdbVersion = parseVersion(serverOptions.FullVersionOrZero())
		s.databaseName = serverOptions.DatabaseNameOrZero()

		dec.SetAlphanumDfv1(serverOptions.DataFormatVersion2OrZero() == p.DfvLevel1)

		if err := s.setSchema(ctx); err != nil {
			dbConn.Close()
			return nil, err
		}
	}
	return s, nil
}

// we cannot work with nested errors containing driver.ErrBadConn
// as go sql retries these statements.
func (s *session) isBad() bool { return s.pwr.CancelledOrError() || s.prd.Cancelled() }

func (s *session) close() error {
	// do not disconnect if isBad.
	var disconnectErr error
	if !s.isBad() {
		disconnectErr = s.disconnect(context.Background())
	}
	closeErr := s.dbConn.Close()
	return errors.Join(disconnectErr, closeErr)
}

func (s *session) authenticate(ctx context.Context, authHnd *p.AuthHnd, attrs *connAttrs) (*p.ConnectOptions, error) {
	defer metricsAddTimeValue(s.metrics, time.Now(), timeAuth)

	// client context
	clientContext := &p.ClientContext{}
	clientContext.SetVersion(DriverVersion)
	clientContext.SetType(clientType)
	clientContext.SetApplicationProgram(attrs._applicationName)

	initRequest, err := authHnd.InitRequest()
	if err != nil {
		return nil, err
	}
	if err := s.pwr.Write(ctx, p.MtAuthenticate, false, clientContext, initRequest); err != nil {
		return nil, err
	}

	initReply, err := authHnd.InitReply()
	if err != nil {
		return nil, err
	}
	if _, err := s.prd.IterateParts(ctx, 0, func(kind p.PartKind, attrs p.PartAttributes) error {
		if kind == p.PkAuthentication {
			return s.prd.ReadPart(ctx, initReply, nil)
		}
		return p.ErrSkipped
	}); err != nil {
		return nil, err
	}

	finalRequest, err := authHnd.FinalRequest()
	if err != nil {
		return nil, err
	}

	co := &p.ConnectOptions{}
	co.SetDataFormatVersion2(attrs._dfv)
	co.SetClientDistributionMode(p.CdmOff)
	// co.SetClientDistributionMode(p.CdmConnectionStatement)
	// co.SetSelectForUpdateSupported(true) // doesn't seem to make a difference
	/*
		p.CoSplitBatchCommands:          true,
		p.CoCompleteArrayExecution:      true,
	*/

	if attrs._locale != "" {
		co.SetClientLocale(attrs._locale)
	}

	if err := s.pwr.Write(ctx, p.MtConnect, false, finalRequest, p.ClientID(clientID), co); err != nil {
		return nil, err
	}

	finalReply, err := authHnd.FinalReply()
	if err != nil {
		return nil, err
	}

	ti := new(p.TopologyInformation)

	if _, err := s.prd.IterateParts(ctx, 0, func(kind p.PartKind, attrs p.PartAttributes) error {
		switch kind {
		case p.PkAuthentication:
			return s.prd.ReadPart(ctx, finalReply, nil)
		case p.PkConnectOptions:
			return s.prd.ReadPart(ctx, co, nil)
		case p.PkTopologyInformation:
			return s.prd.ReadPart(ctx, ti, nil)
		default:
			return p.ErrSkipped
		}
	}); err != nil {
		return nil, err
	}

	sessionID := s.prd.SessionID()
	if sessionID <= 0 {
		return nil, fmt.Errorf("invalid session id %d", sessionID)
	}
	s.pwr.SetSessionID(sessionID)
	// log.Printf("co: %s", co)
	// log.Printf("ti: %s", ti)
	return co, nil
}

func (s *session) setSchema(ctx context.Context) error {
	switch {
	case s.user != nil && s.user.Schema != "":
		_, err := s.execDirect(ctx, "set schema "+Identifier(s.user.Schema).String())
		return err
	case s.attrs._defaultSchema != "":
		_, err := s.execDirect(ctx, "set schema "+Identifier(s.attrs._defaultSchema).String())
		return err
	default:
		return nil
	}
}

// ErrSwitchUser is the error raised if a switch user is requested in a not allowed context.
var ErrSwitchUser = errors.New("switch user inside transaction or in statement scope (prepared query) is not allowed")

func (s *session) switchUser(ctx context.Context) error {
	user, ok := ctx.Value(switchUserCtxKey).(*SessionUser)
	if !ok || user.equal(s.user) {
		return nil
	}
	if s.inTx {
		return ErrSwitchUser
	}
	s.user = user.clone()
	if _, err := s.execDirect(ctx, "connect "+user.Username+" password "+user.Password); err != nil {
		return err
	}
	s.metrics.msgCh <- counterMsg{idx: counterSessionConnects, v: uint64(1)}
	return s.setSchema(ctx)
}

func (s *session) preventSwitchUser(ctx context.Context) error {
	user, ok := ctx.Value(switchUserCtxKey).(*SessionUser)
	if !ok || user.equal(s.user) {
		return nil
	}
	return ErrSwitchUser
}

func (s *session) dbConnectInfo(ctx context.Context, databaseName string) (*DBConnectInfo, error) {
	ci := &p.DBConnectInfo{}
	ci.SetDatabaseName(databaseName)
	if err := s.pwr.Write(ctx, p.MtDBConnectInfo, false, ci); err != nil {
		return nil, err
	}

	if _, err := s.prd.IterateParts(ctx, 0, func(kind p.PartKind, attrs p.PartAttributes) error {
		if kind == p.PkDBConnectInfo {
			return s.prd.ReadPart(ctx, ci, nil)
		}
		return p.ErrSkipped
	}); err != nil {
		return nil, err
	}

	return &DBConnectInfo{
		DatabaseName: databaseName,
		Host:         ci.HostOrZero(),
		Port:         ci.PortOrZero(),
		IsConnected:  ci.IsConnectedOrZero(),
	}, nil
}

func (s *session) queryDirect(ctx context.Context, query string) (driver.Rows, error) {
	t := time.Now()
	defer metricsAddSQLTimeValue(s.metrics, time.Now(), sqlTimeQuery)

	// allow e.g inserts as query -> handle commit like in _execDirect
	if err := s.pwr.Write(ctx, p.MtExecuteDirect, !s.inTx, p.Command(query)); err != nil {
		return nil, err
	}

	qr := &queryResult{session: s}
	meta := &p.ResultMetadata{}
	resSet := &p.Resultset{}

	if _, err := s.prd.IterateParts(ctx, 0, func(kind p.PartKind, attrs p.PartAttributes) error {
		switch kind {
		case p.PkResultMetadata:
			if err := s.prd.ReadPart(ctx, meta, nil); err != nil {
				return err
			}
			qr.fields = meta.ResultFields
			return nil
		case p.PkResultsetID:
			return s.prd.ReadPart(ctx, (*p.ResultsetID)(&qr.rsID), nil)
		case p.PkResultset:
			resSet.ResultFields = qr.fields
			if err := s.prd.ReadPart(ctx, resSet, qr); err != nil {
				return err
			}
			qr.fieldValues = resSet.FieldValues
			qr.decodeErrors = resSet.DecodeErrors
			qr.attrs = attrs
			return nil
		default:
			return p.ErrSkipped
		}
	}); err != nil {
		return nil, err
	}
	if s.sqlTracer != nil {
		s.sqlTracer.log(ctx, t, traceQueryLogKind(query), query)
	}
	if qr.rsID == 0 { // non select query
		return noResult, nil
	}
	return qr, nil
}

func (s *session) execDirect(ctx context.Context, query string) (driver.Result, error) {
	t := time.Now()
	defer metricsAddSQLTimeValue(s.metrics, time.Now(), sqlTimeExec)

	if err := s.pwr.Write(ctx, p.MtExecuteDirect, !s.inTx, p.Command(query)); err != nil {
		return nil, err
	}

	numRow, err := s.prd.IterateParts(ctx, 0, nil)
	if err != nil {
		return nil, err
	}
	if s.sqlTracer != nil {
		s.sqlTracer.log(ctx, t, traceExec, query)
	}
	if s.prd.FunctionCode() == p.FcDDL {
		return driver.ResultNoRows, nil
	}
	return driver.RowsAffected(numRow), nil
}

func (s *session) prepare(ctx context.Context, query string) (*prepareResult, error) {
	t := time.Now()
	defer metricsAddSQLTimeValue(s.metrics, time.Now(), sqlTimePrepare)

	if err := s.pwr.Write(ctx, p.MtPrepare, false, p.Command(query)); err != nil {
		return nil, err
	}

	pr := &prepareResult{}
	resMeta := &p.ResultMetadata{}
	prmMeta := &p.ParameterMetadata{}

	if _, err := s.prd.IterateParts(ctx, 0, func(kind p.PartKind, attrs p.PartAttributes) error {
		switch kind {
		case p.PkStatementID:
			return s.prd.ReadPart(ctx, (*p.StatementID)(&pr.stmtID), nil)
		case p.PkResultMetadata:
			if err := s.prd.ReadPart(ctx, resMeta, nil); err != nil {
				return err
			}
			pr.resultFields = resMeta.ResultFields
			return nil
		case p.PkParameterMetadata:
			if err := s.prd.ReadPart(ctx, prmMeta, nil); err != nil {
				return err
			}
			pr.parameterFields = prmMeta.ParameterFields
			return nil
		default:
			return p.ErrSkipped
		}
	}); err != nil {
		return nil, err
	}
	pr.fc = s.prd.FunctionCode()
	if s.sqlTracer != nil {
		s.sqlTracer.log(ctx, t, tracePrepare, query)
	}
	return pr, nil
}

func (s *session) query(ctx context.Context, query string, pr *prepareResult, nvargs []driver.NamedValue) (driver.Rows, error) {
	t := time.Now()
	defer metricsAddSQLTimeValue(s.metrics, time.Now(), sqlTimeQuery)

	// allow e.g inserts as query -> handle commit like in exec

	if err := convertQueryArgs(pr.parameterFields, nvargs, s.cesu8Encoder, s.attrs._lobChunkSize); err != nil {
		return nil, err
	}
	inputParameters, err := p.NewInputParameters(pr.parameterFields, nvargs)
	if err != nil {
		return nil, err
	}
	if err := s.pwr.Write(ctx, p.MtExecute, !s.inTx, p.StatementID(pr.stmtID), inputParameters); err != nil {
		return nil, err
	}

	qr := &queryResult{session: s, fields: pr.resultFields}
	resSet := &p.Resultset{}

	if _, err := s.prd.IterateParts(ctx, 0, func(kind p.PartKind, attrs p.PartAttributes) error {
		switch kind {
		case p.PkResultsetID:
			return s.prd.ReadPart(ctx, (*p.ResultsetID)(&qr.rsID), nil)
		case p.PkResultset:
			resSet.ResultFields = qr.fields
			if err := s.prd.ReadPart(ctx, resSet, qr); err != nil {
				return err
			}
			qr.fieldValues = resSet.FieldValues
			qr.decodeErrors = resSet.DecodeErrors
			qr.attrs = attrs
			return nil
		default:
			return p.ErrSkipped
		}
	}); err != nil {
		return nil, err
	}
	if s.sqlTracer != nil {
		s.sqlTracer.log(ctx, t, traceQuery, query, nvargs...)
	}
	if qr.rsID == 0 { // non select query
		return noResult, nil
	}
	return qr, nil
}

func (s *session) exec(ctx context.Context, query string, pr *prepareResult, nvargs []driver.NamedValue, offset int) (driver.Result, error) {
	t := time.Now()
	defer metricsAddSQLTimeValue(s.metrics, time.Now(), sqlTimeExec)

	inputParameters, err := p.NewInputParameters(pr.parameterFields, nvargs)
	if err != nil {
		return nil, err
	}
	if err := s.pwr.Write(ctx, p.MtExecute, !s.inTx, p.StatementID(pr.stmtID), inputParameters); err != nil {
		return nil, err
	}

	var ids []p.LocatorID
	lobReply := &p.WriteLobReply{}

	numRow, err := s.prd.IterateParts(ctx, offset, func(kind p.PartKind, attrs p.PartAttributes) error {
		switch kind {
		case p.PkWriteLobReply:
			if err := s.prd.ReadPart(ctx, lobReply, nil); err != nil {
				return err
			}
			ids = lobReply.IDs
			return nil
		default:
			return p.ErrSkipped
		}
	})
	if err != nil {
		return nil, err
	}
	fc := s.prd.FunctionCode()

	if len(ids) != 0 {
		/*
			writeLobParameters:
			- chunkReaders
			- nil (no callResult, exec does not have output parameters)
		*/

		/*
			write lob data only for the last record as lob streaming is only available for the last one
		*/
		startLastRec := len(nvargs) - len(pr.parameterFields)
		if err := s.writeLobs(ctx, nil, ids, pr.parameterFields, nvargs[startLastRec:]); err != nil {
			return nil, err
		}
	}
	if s.sqlTracer != nil {
		s.sqlTracer.log(ctx, t, traceExec, query, nvargs...)
	}
	if fc == p.FcDDL {
		return driver.ResultNoRows, nil
	}
	return driver.RowsAffected(numRow), nil
}

func (s *session) execCall(ctx context.Context, query string, pr *prepareResult, nvargs []driver.NamedValue) (*callResult, *callArgs, int64, error) {
	t := time.Now()
	defer metricsAddSQLTimeValue(s.metrics, time.Now(), sqlTimeCall)

	callArgs, err := convertCallArgs(pr.parameterFields, nvargs, s.cesu8Encoder, s.attrs._lobChunkSize)
	if err != nil {
		return nil, nil, 0, err
	}
	inputParameters, err := p.NewInputParameters(callArgs.inFields, callArgs.inArgs)
	if err != nil {
		return nil, nil, 0, err
	}

	if err := s.pwr.Write(ctx, p.MtExecute, !s.inTx, (*p.StatementID)(&pr.stmtID), inputParameters); err != nil {
		return nil, nil, 0, err
	}

	cr := &callResult{session: s, outFields: callArgs.outFields}

	var qr *queryResult
	var ids []p.LocatorID
	outPrms := &p.OutputParameters{}
	meta := &p.ResultMetadata{}
	resSet := &p.Resultset{}
	lobReply := &p.WriteLobReply{}
	tableRowIdx := 0

	numRow, err := s.prd.IterateParts(ctx, 0, func(kind p.PartKind, attrs p.PartAttributes) error {
		switch kind {
		case p.PkOutputParameters:
			outPrms.OutputFields = cr.outFields
			if err := s.prd.ReadPart(ctx, outPrms, cr); err != nil {
				return err
			}
			cr.fieldValues = outPrms.FieldValues
			cr.decodeErrors = outPrms.DecodeErrors
			return nil
		case p.PkResultMetadata:
			/*
				procedure call with table parameters does return metadata for each table
				sequence: metadata, resultsetID, resultset
				but:
				- resultset might not be provided for all tables
				- so, 'additional' query result is detected by new metadata part
			*/
			qr = &queryResult{session: s}
			cr.outFields = append(cr.outFields, p.NewTableRowsParameterField(tableRowIdx))
			cr.fieldValues = append(cr.fieldValues, qr)
			tableRowIdx++
			if err := s.prd.ReadPart(ctx, meta, nil); err != nil {
				return err
			}
			qr.fields = meta.ResultFields
			return nil
		case p.PkResultset:
			resSet.ResultFields = qr.fields
			if err := s.prd.ReadPart(ctx, resSet, qr); err != nil {
				return err
			}
			qr.fieldValues = resSet.FieldValues
			qr.decodeErrors = resSet.DecodeErrors
			qr.attrs = attrs
			return nil
		case p.PkResultsetID:
			return s.prd.ReadPart(ctx, (*p.ResultsetID)(&qr.rsID), nil)
		case p.PkWriteLobReply:
			if err := s.prd.ReadPart(ctx, lobReply, nil); err != nil {
				return err
			}
			ids = lobReply.IDs
			return nil
		default:
			return p.ErrSkipped
		}
	})
	if err != nil {
		return nil, nil, 0, err
	}

	if len(ids) != 0 {
		/*
			writeLobParameters:
			- chunkReaders
			- cr (callResult output parameters are set after all lob input parameters are written)
		*/
		if err := s.writeLobs(ctx, cr, ids, callArgs.inFields, callArgs.inArgs); err != nil {
			return nil, nil, 0, err
		}
	}
	if s.sqlTracer != nil {
		s.sqlTracer.log(ctx, t, traceExecCall, query, nvargs...)
	}
	return cr, callArgs, numRow, nil
}

func (s *session) fetchNext(ctx context.Context, qr *queryResult) error {
	defer metricsAddSQLTimeValue(s.metrics, time.Now(), sqlTimeFetch)

	if err := s.pwr.Write(ctx, p.MtFetchNext, false, p.ResultsetID(qr.rsID), p.Fetchsize(s.attrs._fetchSize)); err != nil { //nolint: gosec
		return err
	}

	resSet := &p.Resultset{ResultFields: qr.fields, FieldValues: qr.fieldValues} // reuse field values

	_, err := s.prd.IterateParts(ctx, 0, func(kind p.PartKind, attrs p.PartAttributes) error {
		switch kind {
		case p.PkResultset:
			if err := s.prd.ReadPart(ctx, resSet, qr); err != nil {
				return err
			}
			qr.fieldValues = resSet.FieldValues
			qr.decodeErrors = resSet.DecodeErrors
			qr.attrs = attrs
			return nil
		default:
			return p.ErrSkipped
		}
	})
	return err
}

func (s *session) dropStatementID(ctx context.Context, id uint64) error {
	if err := s.pwr.Write(ctx, p.MtDropStatementID, false, p.StatementID(id)); err != nil {
		return err
	}
	return s.prd.SkipParts(ctx)
}

func (s *session) closeResultsetID(ctx context.Context, id uint64) error {
	if err := s.pwr.Write(ctx, p.MtCloseResultset, false, p.ResultsetID(id)); err != nil {
		return err
	}
	return s.prd.SkipParts(ctx)
}

func (s *session) commit(ctx context.Context) error {
	defer metricsAddSQLTimeValue(s.metrics, time.Now(), sqlTimeCommit)

	if err := s.pwr.Write(ctx, p.MtCommit, false); err != nil {
		return err
	}
	if err := s.prd.SkipParts(ctx); err != nil {
		return err
	}
	return nil
}

func (s *session) rollback(ctx context.Context) error {
	defer metricsAddSQLTimeValue(s.metrics, time.Now(), sqlTimeRollback)

	if err := s.pwr.Write(ctx, p.MtRollback, false); err != nil {
		return err
	}
	if err := s.prd.SkipParts(ctx); err != nil {
		return err
	}
	return nil
}

func (s *session) disconnect(ctx context.Context) error {
	if err := s.pwr.Write(ctx, p.MtDisconnect, false); err != nil {
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

/*
readLob reads output lob or result lob parameters from db.

read lob reply
  - seems like readLobreply returns only a result for one lob - even if more then one is requested
    --> read single lobs
*/
func (s *session) readLob(ctx context.Context, request *p.ReadLobRequest, reply *p.ReadLobReply) error {
	defer metricsAddSQLTimeValue(s.metrics, time.Now(), sqlTimeFetchLob)

	var err error
	for err != io.EOF { //nolint: errorlint
		if err = s.pwr.Write(ctx, p.MtWriteLob, false, request); err != nil {
			return err
		}

		if _, err = s.prd.IterateParts(ctx, 0, func(kind p.PartKind, attrs p.PartAttributes) error {
			if kind == p.PkReadLobReply {
				return s.prd.ReadPart(ctx, reply, nil)
			}
			return p.ErrSkipped
		}); err != nil {
			return err
		}

		_, err = reply.Write()
		if err != nil && err != io.EOF { //nolint: errorlint
			return err
		}
	}
	return nil
}

// writeLobs writes input lob parameters to db.
func (s *session) writeLobs(ctx context.Context, cr *callResult, ids []p.LocatorID, inPrmFields []*p.ParameterField, nvargs []driver.NamedValue) error {
	if len(inPrmFields) != len(nvargs) {
		panic("lob streaming can only be done for one (the last) record")
	}
	descrs := make([]*p.WriteLobDescr, 0, len(ids))
	j := 0
	for i, f := range inPrmFields {
		if f.IsLob() {
			lobInDescr, ok := nvargs[i].Value.(*p.LobInDescr)
			if !ok {
				return fmt.Errorf("protocol error: invalid lob parameter %[1]T %[1]v - *lobInDescr expected", nvargs[i])
			}
			if j > len(ids) {
				return fmt.Errorf("protocol error: invalid number of lob parameter ids %d", len(ids))
			}
			if !lobInDescr.IsLastData() {
				descrs = append(descrs, &p.WriteLobDescr{LobInDescr: lobInDescr, ID: ids[j]})
				j++
			}
		}
	}

	writeLobRequest := &p.WriteLobRequest{}
	for len(descrs) != 0 {

		if len(descrs) != len(ids) {
			return fmt.Errorf("protocol error: invalid number of lob parameter ids %d - expected %d", len(descrs), len(ids))
		}
		for i, descr := range descrs { // check if ids and descrs are in sync
			if descr.ID != ids[i] {
				return fmt.Errorf("protocol error: lob parameter id mismatch %d - expected %d", descr.ID, ids[i])
			}
		}

		// TODO check total size limit
		for _, descr := range descrs {
			if err := descr.FetchNext(s.attrs._lobChunkSize); err != nil {
				return err
			}
		}

		writeLobRequest.Descrs = descrs

		if err := s.pwr.Write(ctx, p.MtReadLob, false, writeLobRequest); err != nil {
			return err
		}

		lobReply := &p.WriteLobReply{}
		outPrms := &p.OutputParameters{}

		if _, err := s.prd.IterateParts(ctx, 0, func(kind p.PartKind, attrs p.PartAttributes) error {
			switch kind {
			case p.PkOutputParameters:
				outPrms.OutputFields = cr.outFields
				if err := s.prd.ReadPart(ctx, outPrms, nil); err != nil {
					return err
				}
				cr.fieldValues = outPrms.FieldValues
				cr.decodeErrors = outPrms.DecodeErrors
				return nil
			case p.PkWriteLobReply:
				if err := s.prd.ReadPart(ctx, lobReply, nil); err != nil {
					return err
				}
				ids = lobReply.IDs
				return nil
			default:
				return p.ErrSkipped
			}
		}); err != nil {
			return err
		}

		// remove done descr
		j := 0
		for _, descr := range descrs {
			if !descr.IsLastData() {
				descrs[j] = descr
				j++
			}
		}
		descrs = descrs[:j]
	}
	return nil
}
