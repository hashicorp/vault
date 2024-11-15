// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apache/arrow/go/v15/arrow/ipc"
)

const (
	httpHeaderContentType      = "Content-Type"
	httpHeaderAccept           = "accept"
	httpHeaderUserAgent        = "User-Agent"
	httpHeaderServiceName      = "X-Snowflake-Service"
	httpHeaderContentLength    = "Content-Length"
	httpHeaderHost             = "Host"
	httpHeaderValueOctetStream = "application/octet-stream"
	httpHeaderContentEncoding  = "Content-Encoding"
	httpClientAppID            = "CLIENT_APP_ID"
	httpClientAppVersion       = "CLIENT_APP_VERSION"
)

const (
	statementTypeIDSelect           = int64(0x1000)
	statementTypeIDDml              = int64(0x3000)
	statementTypeIDMultiTableInsert = statementTypeIDDml + int64(0x500)
	statementTypeIDMultistatement   = int64(0xA000)
)

const (
	sessionClientSessionKeepAlive          = "client_session_keep_alive"
	sessionClientValidateDefaultParameters = "CLIENT_VALIDATE_DEFAULT_PARAMETERS"
	sessionArrayBindStageThreshold         = "client_stage_array_binding_threshold"
	serviceName                            = "service_name"
)

type resultType string

const (
	snowflakeResultType contextKey = "snowflakeResultType"
	execResultType      resultType = "exec"
	queryResultType     resultType = "query"
)

type execKey string

const (
	executionType          execKey = "executionType"
	executionTypeStatement string  = "statement"
)

// snowflakeConn manages its own context.
// External cancellation should not be supported because the connection
// may be reused after the original query/request has completed.
type snowflakeConn struct {
	ctx                 context.Context
	cfg                 *Config
	rest                *snowflakeRestful
	SequenceCounter     uint64
	telemetry           *snowflakeTelemetry
	internal            InternalClient
	queryContextCache   *queryContextCache
	currentTimeProvider currentTimeProvider
}

var (
	queryIDPattern = `[\w\-_]+`
	queryIDRegexp  = regexp.MustCompile(queryIDPattern)
)

func (sc *snowflakeConn) exec(
	ctx context.Context,
	query string,
	noResult bool,
	isInternal bool,
	describeOnly bool,
	bindings []driver.NamedValue) (
	*execResponse, error) {
	var err error
	counter := atomic.AddUint64(&sc.SequenceCounter, 1) // query sequence counter

	queryContext, err := buildQueryContext(sc.queryContextCache)
	if err != nil {
		logger.WithContext(ctx).Errorf("error while building query context: %v", err)
	}
	req := execRequest{
		SQLText:      query,
		AsyncExec:    noResult,
		Parameters:   map[string]interface{}{},
		IsInternal:   isInternal,
		DescribeOnly: describeOnly,
		SequenceID:   counter,
		QueryContext: queryContext,
	}
	if key := ctx.Value(multiStatementCount); key != nil {
		req.Parameters[string(multiStatementCount)] = key
	}
	if tag := ctx.Value(queryTag); tag != nil {
		req.Parameters[string(queryTag)] = tag
	}
	logger.WithContext(ctx).Infof("parameters: %v", req.Parameters)

	// handle bindings, if required
	requestID := getOrGenerateRequestIDFromContext(ctx)
	if len(bindings) > 0 {
		if err = sc.processBindings(ctx, bindings, describeOnly, requestID, &req); err != nil {
			return nil, err
		}
	}
	logger.WithContext(ctx).Infof("bindings: %v", req.Bindings)

	// populate headers
	headers := getHeaders()
	if isFileTransfer(query) {
		headers[httpHeaderAccept] = headerContentTypeApplicationJSON
	}
	paramsMutex.Lock()
	if serviceName, ok := sc.cfg.Params[serviceName]; ok {
		headers[httpHeaderServiceName] = *serviceName
	}
	paramsMutex.Unlock()

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	data, err := sc.rest.FuncPostQuery(ctx, sc.rest, &url.Values{}, headers,
		jsonBody, sc.rest.RequestTimeout, requestID, sc.cfg)
	if err != nil {
		return data, err
	}
	code := -1
	if data.Code != "" {
		code, err = strconv.Atoi(data.Code)
		if err != nil {
			return data, err
		}
	}
	logger.WithContext(ctx).Infof("Success: %v, Code: %v", data.Success, code)
	if !data.Success {
		err = (populateErrorFields(code, data)).exceptionTelemetry(sc)
		return nil, err
	}

	if !sc.cfg.DisableQueryContextCache && data.Data.QueryContext != nil {
		queryContext, err := extractQueryContext(data)
		if err != nil {
			logger.WithContext(ctx).Errorf("error while decoding query context: %v", err)
		} else {
			sc.queryContextCache.add(sc, queryContext.Entries...)
		}
	}

	// handle PUT/GET commands
	fileTransferChan := make(chan error, 1)
	if isFileTransfer(query) {
		go func() {
			data, err = sc.processFileTransfer(ctx, data, query, isInternal)
			fileTransferChan <- err
		}()

		select {
		case <-ctx.Done():
			logger.WithContext(ctx).Info("File transfer has been cancelled")
			return nil, ctx.Err()
		case err := <-fileTransferChan:
			if err != nil {
				return nil, err
			}
		}
	}

	logger.WithContext(ctx).Infof("Exec/Query SUCCESS with total=%v, returned=%v", data.Data.Total, data.Data.Returned)
	if data.Data.FinalDatabaseName != "" {
		sc.cfg.Database = data.Data.FinalDatabaseName
	}
	if data.Data.FinalSchemaName != "" {
		sc.cfg.Schema = data.Data.FinalSchemaName
	}
	if data.Data.FinalWarehouseName != "" {
		sc.cfg.Warehouse = data.Data.FinalWarehouseName
	}
	if data.Data.FinalRoleName != "" {
		sc.cfg.Role = data.Data.FinalRoleName
	}
	sc.populateSessionParameters(data.Data.Parameters)
	return data, err
}

func extractQueryContext(data *execResponse) (queryContext, error) {
	var queryContext queryContext
	err := json.Unmarshal(data.Data.QueryContext, &queryContext)
	return queryContext, err
}

func buildQueryContext(qcc *queryContextCache) (requestQueryContext, error) {
	rqc := requestQueryContext{}
	if qcc == nil || len(qcc.entries) == 0 {
		logger.Debugf("empty qcc")
		return rqc, nil
	}
	for _, qce := range qcc.entries {
		contextData := contextData{}
		if qce.Context == "" {
			contextData.Base64Data = qce.Context
		}
		rqc.Entries = append(rqc.Entries, requestQueryContextEntry{
			ID:        qce.ID,
			Priority:  qce.Priority,
			Timestamp: qce.Timestamp,
			Context:   contextData,
		})
	}
	return rqc, nil
}

func (sc *snowflakeConn) Begin() (driver.Tx, error) {
	return sc.BeginTx(sc.ctx, driver.TxOptions{})
}

func (sc *snowflakeConn) BeginTx(
	ctx context.Context,
	opts driver.TxOptions) (
	driver.Tx, error) {
	logger.WithContext(ctx).Info("BeginTx")
	if opts.ReadOnly {
		return nil, (&SnowflakeError{
			Number:   ErrNoReadOnlyTransaction,
			SQLState: SQLStateFeatureNotSupported,
			Message:  errMsgNoReadOnlyTransaction,
		}).exceptionTelemetry(sc)
	}
	if int(opts.Isolation) != int(sql.LevelDefault) {
		return nil, (&SnowflakeError{
			Number:   ErrNoDefaultTransactionIsolationLevel,
			SQLState: SQLStateFeatureNotSupported,
			Message:  errMsgNoDefaultTransactionIsolationLevel,
		}).exceptionTelemetry(sc)
	}
	if sc.rest == nil {
		return nil, driver.ErrBadConn
	}
	isDesc := isDescribeOnly(ctx)
	if _, err := sc.exec(ctx, "BEGIN", false, /* noResult */
		false /* isInternal */, isDesc, nil); err != nil {
		return nil, err
	}
	return &snowflakeTx{sc, ctx}, nil
}

func (sc *snowflakeConn) cleanup() {
	// must flush log buffer while the process is running.
	logger.WithContext(sc.ctx).Debugln("Snowflake connection closing.")
	if sc.rest != nil && sc.rest.Client != nil {
		sc.rest.Client.CloseIdleConnections()
	}
	sc.rest = nil
	sc.cfg = nil
}

func (sc *snowflakeConn) Close() (err error) {
	logger.WithContext(sc.ctx).Infoln("Close")
	sc.telemetry.sendBatch()
	sc.stopHeartBeat()
	defer sc.cleanup()

	if sc.cfg != nil && !sc.cfg.KeepSessionAlive {
		if err = sc.rest.FuncCloseSession(sc.ctx, sc.rest, sc.rest.RequestTimeout); err != nil {
			logger.WithContext(sc.ctx).Error(err)
		}
	}
	return nil
}

func (sc *snowflakeConn) PrepareContext(
	ctx context.Context,
	query string) (
	driver.Stmt, error) {
	logger.WithContext(sc.ctx).Infoln("Prepare")
	if sc.rest == nil {
		return nil, driver.ErrBadConn
	}
	stmt := &snowflakeStmt{
		sc:    sc,
		query: query,
	}
	return stmt, nil
}

func (sc *snowflakeConn) ExecContext(
	ctx context.Context,
	query string,
	args []driver.NamedValue) (
	driver.Result, error) {
	logger.WithContext(ctx).Infof("Exec: %#v, %v", query, args)
	if sc.rest == nil {
		return nil, driver.ErrBadConn
	}
	noResult := isAsyncMode(ctx)
	isDesc := isDescribeOnly(ctx)
	// TODO handle isInternal
	ctx = setResultType(ctx, execResultType)
	data, err := sc.exec(ctx, query, noResult, false /* isInternal */, isDesc, args)
	if err != nil {
		logger.WithContext(ctx).Infof("error: %v", err)
		if data != nil {
			code, e := strconv.Atoi(data.Code)
			if e != nil {
				return nil, e
			}
			return nil, (&SnowflakeError{
				Number:   code,
				SQLState: data.Data.SQLState,
				Message:  err.Error(),
				QueryID:  data.Data.QueryID,
			}).exceptionTelemetry(sc)
		}
		return nil, err
	}

	// if async exec, return result object right away
	if noResult {
		return data.Data.AsyncResult, nil
	}

	if isDml(data.Data.StatementTypeID) {
		// collects all values from the returned row sets
		updatedRows, err := updateRows(data.Data)
		if err != nil {
			return nil, err
		}
		logger.WithContext(ctx).Debugf("number of updated rows: %#v", updatedRows)
		return &snowflakeResult{
			affectedRows: updatedRows,
			insertID:     -1,
			queryID:      data.Data.QueryID,
		}, nil // last insert id is not supported by Snowflake
	} else if isMultiStmt(&data.Data) {
		return sc.handleMultiExec(ctx, data.Data)
	} else if isDql(&data.Data) {
		logger.WithContext(ctx).Debugf("DQL")
		if isStatementContext(ctx) {
			return &snowflakeResultNoRows{queryID: data.Data.QueryID}, nil
		}
		return driver.ResultNoRows, nil
	}
	logger.WithContext(ctx).Debug("DDL")
	if isStatementContext(ctx) {
		return &snowflakeResultNoRows{queryID: data.Data.QueryID}, nil
	}
	return driver.ResultNoRows, nil
}

func (sc *snowflakeConn) QueryContext(
	ctx context.Context,
	query string,
	args []driver.NamedValue) (
	driver.Rows, error) {
	qid, err := getResumeQueryID(ctx)
	if err != nil {
		return nil, err
	}
	if qid == "" {
		return sc.queryContextInternal(ctx, query, args)
	}

	// check the query status to find out if there is a result to fetch
	_, err = sc.checkQueryStatus(ctx, qid)
	snowflakeErr, isSnowflakeError := err.(*SnowflakeError)
	if err == nil || (isSnowflakeError && snowflakeErr.Number == ErrQueryIsRunning) {
		// the query is running. Rows object will be returned from here.
		return sc.buildRowsForRunningQuery(ctx, qid)
	}
	return nil, err
}

func (sc *snowflakeConn) queryContextInternal(
	ctx context.Context,
	query string,
	args []driver.NamedValue) (
	driver.Rows, error) {
	logger.WithContext(ctx).Infof("Query: %#v, %v", query, args)
	if sc.rest == nil {
		return nil, driver.ErrBadConn
	}

	noResult := isAsyncMode(ctx)
	isDesc := isDescribeOnly(ctx)
	ctx = setResultType(ctx, queryResultType)
	// TODO: handle isInternal
	data, err := sc.exec(ctx, query, noResult, false /* isInternal */, isDesc, args)
	if err != nil {
		logger.WithContext(ctx).Errorf("error: %v", err)
		if data != nil {
			code, e := strconv.Atoi(data.Code)
			if e != nil {
				return nil, e
			}
			return nil, (&SnowflakeError{
				Number:   code,
				SQLState: data.Data.SQLState,
				Message:  err.Error(),
				QueryID:  data.Data.QueryID,
			}).exceptionTelemetry(sc)
		}
		return nil, err
	}

	// if async query, return row object right away
	if noResult {
		return data.Data.AsyncRows, nil
	}

	rows := new(snowflakeRows)
	rows.sc = sc
	rows.queryID = data.Data.QueryID
	rows.ctx = ctx

	if isMultiStmt(&data.Data) {
		// handleMultiQuery is responsible to fill rows with childResults
		if err = sc.handleMultiQuery(ctx, data.Data, rows); err != nil {
			return nil, err
		}
	} else {
		rows.addDownloader(populateChunkDownloader(ctx, sc, data.Data))
	}

	err = rows.ChunkDownloader.start()
	return rows, err
}

func (sc *snowflakeConn) Prepare(query string) (driver.Stmt, error) {
	return sc.PrepareContext(sc.ctx, query)
}

func (sc *snowflakeConn) Exec(
	query string,
	args []driver.Value) (
	driver.Result, error) {
	return sc.ExecContext(sc.ctx, query, toNamedValues(args))
}

func (sc *snowflakeConn) Query(
	query string,
	args []driver.Value) (
	driver.Rows, error) {
	return sc.QueryContext(sc.ctx, query, toNamedValues(args))
}

func (sc *snowflakeConn) Ping(ctx context.Context) error {
	logger.WithContext(ctx).Infoln("Ping")
	if sc.rest == nil {
		return driver.ErrBadConn
	}
	noResult := isAsyncMode(ctx)
	isDesc := isDescribeOnly(ctx)
	// TODO: handle isInternal
	ctx = setResultType(ctx, execResultType)
	_, err := sc.exec(ctx, "SELECT 1", noResult, false, /* isInternal */
		isDesc, []driver.NamedValue{})
	return err
}

// CheckNamedValue determines which types are handled by this driver aside from
// the instances captured by driver.Value
func (sc *snowflakeConn) CheckNamedValue(nv *driver.NamedValue) error {
	if supportedNullBind(nv) || supportedArrayBind(nv) || supportedStructuredObjectWriterBind(nv) || supportedStructuredArrayBind(nv) || supportedStructuredMapBind(nv) {
		return nil
	}
	return driver.ErrSkip
}

func (sc *snowflakeConn) GetQueryStatus(
	ctx context.Context,
	queryID string) (
	*SnowflakeQueryStatus, error) {
	queryRet, err := sc.checkQueryStatus(ctx, queryID)
	if err != nil {
		return nil, err
	}
	return &SnowflakeQueryStatus{
		queryRet.SQLText,
		queryRet.StartTime,
		queryRet.EndTime,
		queryRet.ErrorCode,
		queryRet.ErrorMessage,
		queryRet.Stats.ScanBytes,
		queryRet.Stats.ProducedRows,
	}, nil
}

// QueryArrowStream returns batches which can be queried for their raw arrow
// ipc stream of bytes. This way consumers don't need to be using the exact
// same version of Arrow as the connection is using internally in order
// to consume Arrow data.
func (sc *snowflakeConn) QueryArrowStream(ctx context.Context, query string, bindings ...driver.NamedValue) (ArrowStreamLoader, error) {
	ctx = WithArrowBatches(context.WithValue(ctx, asyncMode, false))
	ctx = setResultType(ctx, queryResultType)
	isDesc := isDescribeOnly(ctx)
	data, err := sc.exec(ctx, query, false, false /* isinternal */, isDesc, bindings)
	if err != nil {
		logger.WithContext(ctx).Errorf("error: %v", err)
		if data != nil {
			code, e := strconv.Atoi(data.Code)
			if e != nil {
				return nil, e
			}
			return nil, (&SnowflakeError{
				Number:   code,
				SQLState: data.Data.SQLState,
				Message:  err.Error(),
				QueryID:  data.Data.QueryID,
			}).exceptionTelemetry(sc)
		}
		return nil, err
	}

	return &snowflakeArrowStreamChunkDownloader{
		sc:          sc,
		ChunkMetas:  data.Data.Chunks,
		Total:       data.Data.Total,
		Qrmk:        data.Data.Qrmk,
		ChunkHeader: data.Data.ChunkHeaders,
		FuncGet:     getChunk,
		RowSet: rowSetType{
			RowType:      data.Data.RowType,
			JSON:         data.Data.RowSet,
			RowSetBase64: data.Data.RowSetBase64,
		},
	}, nil
}

// ArrowStreamBatch is a type describing a potentially yet-to-be-downloaded
// Arrow IPC stream. Call `GetStream` to download and retrieve an io.Reader
// that can be used with ipc.NewReader to get record batch results.
type ArrowStreamBatch struct {
	idx     int
	numrows int64
	scd     *snowflakeArrowStreamChunkDownloader
	Loc     *time.Location
	rr      io.ReadCloser
}

// NumRows returns the total number of rows that the metadata stated should
// be in this stream of record batches.
func (asb *ArrowStreamBatch) NumRows() int64 { return asb.numrows }

// gzip.Reader.Close does NOT close the underlying reader, so we
// need to wrap with wrapReader so that closing will close the
// response body (or any other reader that we want to gzip uncompress)
type wrapReader struct {
	io.Reader
	wrapped io.ReadCloser
}

func (w *wrapReader) Close() error {
	if cl, ok := w.Reader.(io.ReadCloser); ok {
		if err := cl.Close(); err != nil {
			return err
		}
	}
	return w.wrapped.Close()
}

func (asb *ArrowStreamBatch) downloadChunkStreamHelper(ctx context.Context) error {
	headers := make(map[string]string)
	if len(asb.scd.ChunkHeader) > 0 {
		logger.WithContext(ctx).Debug("chunk header is provided")
		for k, v := range asb.scd.ChunkHeader {
			logger.Debugf("adding header: %v, value: %v", k, v)

			headers[k] = v
		}
	} else {
		headers[headerSseCAlgorithm] = headerSseCAes
		headers[headerSseCKey] = asb.scd.Qrmk
	}

	resp, err := asb.scd.FuncGet(ctx, asb.scd.sc, asb.scd.ChunkMetas[asb.idx].URL, headers, asb.scd.sc.rest.RequestTimeout)
	if err != nil {
		return err
	}
	logger.WithContext(ctx).Debugf("response returned chunk: %v for URL: %v", asb.idx+1, asb.scd.ChunkMetas[asb.idx].URL)
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		logger.WithContext(ctx).Infof("HTTP: %v, URL: %v, Body: %v", resp.StatusCode, asb.scd.ChunkMetas[asb.idx].URL, b)
		logger.WithContext(ctx).Infof("Header: %v", resp.Header)
		return &SnowflakeError{
			Number:      ErrFailedToGetChunk,
			SQLState:    SQLStateConnectionFailure,
			Message:     errMsgFailedToGetChunk,
			MessageArgs: []interface{}{asb.idx},
		}
	}

	defer func() {
		if asb.rr == nil {
			resp.Body.Close()
		}
	}()

	bufStream := bufio.NewReader(resp.Body)
	gzipMagic, err := bufStream.Peek(2)
	if err != nil {
		return err
	}

	if gzipMagic[0] == 0x1f && gzipMagic[1] == 0x8b {
		// detect and uncompress gzip
		bufStream0, err := gzip.NewReader(bufStream)
		if err != nil {
			return err
		}
		// gzip.Reader.Close() does NOT close the underlying
		// reader, so we need to wrap it and ensure close will
		// close the response body. Otherwise we'll leak it.
		asb.rr = &wrapReader{Reader: bufStream0, wrapped: resp.Body}
	} else {
		asb.rr = &wrapReader{Reader: bufStream, wrapped: resp.Body}
	}
	return nil
}

// GetStream returns a stream of bytes consisting of an Arrow IPC Record
// batch stream. Close should be called on the returned stream when done
// to ensure no leaked memory.
func (asb *ArrowStreamBatch) GetStream(ctx context.Context) (io.ReadCloser, error) {
	if asb.rr == nil {
		if err := asb.downloadChunkStreamHelper(ctx); err != nil {
			return nil, err
		}
	}

	return asb.rr, nil
}

// ArrowStreamLoader is a convenience interface for downloading
// Snowflake results via multiple Arrow Record Batch streams.
//
// Some queries from Snowflake do not return Arrow data regardless
// of the settings, such as "SHOW WAREHOUSES". In these cases,
// you'll find TotalRows() > 0 but GetBatches returns no batches
// and no errors. In this case, the data is accessible via JSONData
// with the actual types matching up to the metadata in RowTypes.
type ArrowStreamLoader interface {
	GetBatches() ([]ArrowStreamBatch, error)
	TotalRows() int64
	RowTypes() []execResponseRowType
	Location() *time.Location
	JSONData() [][]*string
}

type snowflakeArrowStreamChunkDownloader struct {
	sc          *snowflakeConn
	ChunkMetas  []execResponseChunk
	Total       int64
	Qrmk        string
	ChunkHeader map[string]string
	FuncGet     func(context.Context, *snowflakeConn, string, map[string]string, time.Duration) (*http.Response, error)
	RowSet      rowSetType
}

func (scd *snowflakeArrowStreamChunkDownloader) Location() *time.Location {
	if scd.sc != nil && scd.sc.cfg != nil {
		return getCurrentLocation(scd.sc.cfg.Params)
	}
	return nil
}
func (scd *snowflakeArrowStreamChunkDownloader) TotalRows() int64 { return scd.Total }
func (scd *snowflakeArrowStreamChunkDownloader) RowTypes() []execResponseRowType {
	return scd.RowSet.RowType
}
func (scd *snowflakeArrowStreamChunkDownloader) JSONData() [][]*string {
	return scd.RowSet.JSON
}

// the server might have had an empty first batch, check if we can decode
// that first batch, if not we skip it.
func (scd *snowflakeArrowStreamChunkDownloader) maybeFirstBatch() ([]byte, error) {
	if scd.RowSet.RowSetBase64 == "" {
		return nil, nil
	}

	// first batch
	rowSetBytes, err := base64.StdEncoding.DecodeString(scd.RowSet.RowSetBase64)
	if err != nil {
		// match logic in buildFirstArrowChunk
		// assume there's no first chunk if we can't decode the base64 string
		logger.Warnf("skipping first batch as it is not a valid base64 response. %v", err)
		return nil, err
	}

	// verify it's a valid ipc stream, otherwise skip it
	rr, err := ipc.NewReader(bytes.NewReader(rowSetBytes))
	if err != nil {
		logger.Warnf("skipping first batch as it is not a valid IPC stream. %v", err)
		return nil, err
	}
	rr.Release()

	return rowSetBytes, nil
}

func (scd *snowflakeArrowStreamChunkDownloader) GetBatches() (out []ArrowStreamBatch, err error) {
	chunkMetaLen := len(scd.ChunkMetas)
	loc := scd.Location()

	out = make([]ArrowStreamBatch, chunkMetaLen, chunkMetaLen+1)
	toFill := out
	rowSetBytes, err := scd.maybeFirstBatch()
	if err != nil {
		return nil, err
	}
	// if there was no first batch in the response from the server,
	// skip it and move on. toFill == out
	// otherwise expand out by one to account for the first batch
	// and fill it in. have toFill refer to the slice of out excluding
	// the first batch.
	if len(rowSetBytes) > 0 {
		out = out[:chunkMetaLen+1]
		out[0] = ArrowStreamBatch{
			scd: scd,
			Loc: loc,
			rr:  io.NopCloser(bytes.NewReader(rowSetBytes)),
		}
		toFill = out[1:]
	}

	var totalCounted int64
	for i := range toFill {
		toFill[i] = ArrowStreamBatch{
			idx:     i,
			numrows: int64(scd.ChunkMetas[i].RowCount),
			Loc:     loc,
			scd:     scd,
		}
		logger.Debugf("batch %v, numrows: %v", i, toFill[i].numrows)
		totalCounted += int64(scd.ChunkMetas[i].RowCount)
	}

	if len(rowSetBytes) > 0 {
		// if we had a first batch, fill in the numrows
		out[0].numrows = scd.Total - totalCounted
		logger.Debugf("first batch, numrows: %v", out[0].numrows)
	}
	return
}

// buildSnowflakeConn creates a new snowflakeConn.
// The provided context is used only for establishing the initial connection.
func buildSnowflakeConn(ctx context.Context, config Config) (*snowflakeConn, error) {
	sc := &snowflakeConn{
		SequenceCounter:     0,
		ctx:                 context.Background(),
		cfg:                 &config,
		queryContextCache:   (&queryContextCache{}).init(),
		currentTimeProvider: defaultTimeProvider,
	}
	err := initEasyLogging(config.ClientConfigFile)
	if err != nil {
		return nil, err
	}
	var st http.RoundTripper = SnowflakeTransport
	if sc.cfg.Transporter == nil {
		if sc.cfg.InsecureMode {
			// no revocation check with OCSP. Think twice when you want to enable this option.
			st = snowflakeInsecureTransport
		} else {
			// set OCSP fail open mode
			ocspResponseCacheLock.Lock()
			atomic.StoreUint32((*uint32)(&ocspFailOpen), uint32(sc.cfg.OCSPFailOpen))
			ocspResponseCacheLock.Unlock()
		}
	} else {
		// use the custom transport
		st = sc.cfg.Transporter
	}
	if err = setupOCSPEnvVars(ctx, sc.cfg.Host); err != nil {
		return nil, err
	}
	var tokenAccessor TokenAccessor
	if sc.cfg.TokenAccessor != nil {
		tokenAccessor = sc.cfg.TokenAccessor
	} else {
		tokenAccessor = getSimpleTokenAccessor()
	}

	// authenticate
	sc.rest = &snowflakeRestful{
		Host:     sc.cfg.Host,
		Port:     sc.cfg.Port,
		Protocol: sc.cfg.Protocol,
		Client: &http.Client{
			// request timeout including reading response body
			Timeout:   sc.cfg.ClientTimeout,
			Transport: st,
		},
		JWTClient: &http.Client{
			Timeout:   sc.cfg.JWTClientTimeout,
			Transport: st,
		},
		TokenAccessor:       tokenAccessor,
		LoginTimeout:        sc.cfg.LoginTimeout,
		RequestTimeout:      sc.cfg.RequestTimeout,
		MaxRetryCount:       sc.cfg.MaxRetryCount,
		FuncPost:            postRestful,
		FuncGet:             getRestful,
		FuncAuthPost:        postAuthRestful,
		FuncPostQuery:       postRestfulQuery,
		FuncPostQueryHelper: postRestfulQueryHelper,
		FuncRenewSession:    renewRestfulSession,
		FuncPostAuth:        postAuth,
		FuncCloseSession:    closeSession,
		FuncCancelQuery:     cancelQuery,
		FuncPostAuthSAML:    postAuthSAML,
		FuncPostAuthOKTA:    postAuthOKTA,
		FuncGetSSO:          getSSO,
	}

	if sc.cfg.DisableTelemetry {
		sc.telemetry = &snowflakeTelemetry{enabled: false}
	} else {
		sc.telemetry = &snowflakeTelemetry{
			flushSize: defaultFlushSize,
			sr:        sc.rest,
			mutex:     &sync.Mutex{},
			enabled:   true,
		}
	}

	return sc, nil
}
