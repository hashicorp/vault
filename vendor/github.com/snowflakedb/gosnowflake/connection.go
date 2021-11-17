// Copyright (c) 2017-2021 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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
)

const (
	statementTypeIDMulti            = int64(0x1000)
	statementTypeIDDml              = int64(0x3000)
	statementTypeIDMultiTableInsert = statementTypeIDDml + int64(0x500)
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

type snowflakeConn struct {
	ctx             context.Context
	cfg             *Config
	rest            *snowflakeRestful
	SequenceCounter uint64
	QueryID         string
	SQLState        string
	telemetry       *snowflakeTelemetry
	internal        InternalClient
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

	req := execRequest{
		SQLText:      query,
		AsyncExec:    noResult,
		Parameters:   map[string]interface{}{},
		IsInternal:   isInternal,
		DescribeOnly: describeOnly,
		SequenceID:   counter,
	}
	if key := ctx.Value(multiStatementCount); key != nil {
		req.Parameters[string(multiStatementCount)] = key
	}
	logger.WithContext(ctx).Infof("parameters: %v", req.Parameters)

	requestID := getOrGenerateRequestIDFromContext(ctx)
	if len(bindings) > 0 {
		arrayBindThreshold := sc.getArrayBindStageThreshold()
		numBinds := arrayBindValueCount(bindings)
		if 0 < arrayBindThreshold && arrayBindThreshold <= numBinds &&
			!describeOnly && isArrayBind(bindings) {
			// bulk array insert binding
			uploader := bindUploader{
				sc:        sc,
				ctx:       ctx,
				stagePath: "@" + bindStageName + "/" + requestID.String(),
			}
			uploader.upload(bindings)
			req.Bindings = nil
			req.BindStage = uploader.stagePath
		} else {
			// variable or array binding
			req.Bindings, err = getBindValues(bindings)
			if err != nil {
				return nil, err
			}
			req.BindStage = ""
		}
	}
	logger.WithContext(ctx).Infof("bindings: %v", req.Bindings)

	headers := getHeaders()
	if isFileTransfer(query) {
		headers[httpHeaderAccept] = headerContentTypeApplicationJSON
	}
	if serviceName, ok := sc.cfg.Params[serviceName]; ok {
		headers[httpHeaderServiceName] = *serviceName
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var data *execResponse
	data, err = sc.rest.FuncPostQuery(ctx, sc.rest, &url.Values{}, headers,
		jsonBody, sc.rest.RequestTimeout, requestID, sc.cfg)
	if err != nil {
		return data, err
	}
	var code int
	if data.Code != "" {
		code, err = strconv.Atoi(data.Code)
		if err != nil {
			code = -1
			return data, err
		}
	} else {
		code = -1
	}
	logger.WithContext(ctx).Infof("Success: %v, Code: %v", data.Success, code)
	if !data.Success {
		return nil, (&SnowflakeError{
			Number:   code,
			SQLState: data.Data.SQLState,
			Message:  data.Message,
			QueryID:  data.Data.QueryID,
		}).exceptionTelemetry(sc)
	}
	if isFileTransfer(query) {
		data, err = sc.processFileTransfer(ctx, data, query, isInternal)
		if err != nil {
			return nil, err
		}
	}

	logger.WithContext(ctx).Info("Exec/Query SUCCESS")
	sc.cfg.Database = data.Data.FinalDatabaseName
	sc.cfg.Schema = data.Data.FinalSchemaName
	sc.cfg.Role = data.Data.FinalRoleName
	sc.cfg.Warehouse = data.Data.FinalWarehouseName
	sc.QueryID = data.Data.QueryID
	sc.SQLState = data.Data.SQLState
	sc.populateSessionParameters(data.Data.Parameters)
	return data, err
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
	return &snowflakeTx{sc}, nil
}

func (sc *snowflakeConn) cleanup() {
	// must flush log buffer while the process is running.
	sc.rest = nil
	sc.cfg = nil
}

func (sc *snowflakeConn) Close() (err error) {
	logger.WithContext(sc.ctx).Infoln("Close")
	sc.telemetry.sendBatch()
	sc.stopHeartBeat()

	if !sc.cfg.KeepSessionAlive {
		if err = sc.rest.FuncCloseSession(sc.ctx, sc.rest, sc.rest.RequestTimeout); err != nil {
			logger.Error(err)
		}
	}
	sc.cleanup()
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
	noResult := isAsyncMode(ctx)
	data, err := sc.exec(ctx, query, noResult, false, /* isInternal */
		true /* describeOnly */, []driver.NamedValue{})
	if err != nil {
		if data != nil {
			code, err := strconv.Atoi(data.Code)
			if err != nil {
				return nil, err
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
			code, err := strconv.Atoi(data.Code)
			if err != nil {
				return nil, err
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
			queryID:      sc.QueryID,
		}, nil // last insert id is not supported by Snowflake
	} else if isMultiStmt(&data.Data) {
		return sc.handleMultiExec(ctx, data.Data)
	}
	logger.Debug("DDL")
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
	if err == nil || (err != nil && err.(*SnowflakeError).Number == ErrQueryIsRunning) {
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
			code, err := strconv.Atoi(data.Code)
			if err != nil {
				return nil, err
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
	rows.queryID = sc.QueryID

	if isMultiStmt(&data.Data) {
		// handleMultiQuery is responsible to fill rows with childResults
		if err = sc.handleMultiQuery(ctx, data.Data, rows); err != nil {
			return nil, err
		}
	} else {
		rows.addDownloader(populateChunkDownloader(ctx, sc, data.Data))
	}

	rows.ChunkDownloader.start()
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
	_, err := sc.exec(ctx, "SELECT 1", noResult, false, /* isInternal */
		isDesc, []driver.NamedValue{})
	return err
}

// CheckNamedValue determines which types are handled by this driver aside from
// the instances captured by driver.Value
func (sc *snowflakeConn) CheckNamedValue(nv *driver.NamedValue) error {
	if supported := supportedArrayBind(nv); !supported {
		return driver.ErrSkip
	}
	return nil
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

func (sc *snowflakeConn) handleMultiExec(
	ctx context.Context,
	data execResponseData) (
	driver.Result, error) {
	var updatedRows int64
	childResults := getChildResults(data.ResultIDs, data.ResultTypes)
	for _, child := range childResults {
		resultPath := fmt.Sprintf(urlQueriesResultFmt, child.id)
		childData, err := sc.getQueryResultResp(ctx, resultPath)
		if err != nil {
			logger.Errorf("error: %v", err)
			code, err := strconv.Atoi(childData.Code)
			if err != nil {
				return nil, err
			}
			if childData != nil {
				return nil, (&SnowflakeError{
					Number:   code,
					SQLState: childData.Data.SQLState,
					Message:  err.Error(),
					QueryID:  childData.Data.QueryID,
				}).exceptionTelemetry(sc)
			}
			return nil, err
		}
		if isDml(childData.Data.StatementTypeID) {
			count, err := updateRows(childData.Data)
			if err != nil {
				logger.WithContext(ctx).Errorf("error: %v", err)
				if childData != nil {
					code, err := strconv.Atoi(childData.Code)
					if err != nil {
						return nil, err
					}
					return nil, (&SnowflakeError{
						Number:   code,
						SQLState: childData.Data.SQLState,
						Message:  err.Error(),
						QueryID:  childData.Data.QueryID,
					}).exceptionTelemetry(sc)
				}
				return nil, err
			}
			updatedRows += count
		}
	}
	logger.WithContext(ctx).Infof("number of updated rows: %#v", updatedRows)
	return &snowflakeResult{
		affectedRows: updatedRows,
		insertID:     -1,
		queryID:      sc.QueryID,
	}, nil
}

// Fill the corresponding rows and add chunk downloader into the rows when
// iterating across the childResults
func (sc *snowflakeConn) handleMultiQuery(
	ctx context.Context,
	data execResponseData,
	rows *snowflakeRows) error {
	childResults := getChildResults(data.ResultIDs, data.ResultTypes)
	for _, child := range childResults {
		if err := sc.rowsForRunningQuery(ctx, child.id, rows); err != nil {
			return err
		}
	}
	return nil
}

func getAsync(
	ctx context.Context,
	sr *snowflakeRestful,
	headers map[string]string,
	URL *url.URL,
	timeout time.Duration,
	res *snowflakeResult,
	rows *snowflakeRows,
	cfg *Config,
) {
	resType := getResultType(ctx)
	var errChannel chan error
	sfError := &SnowflakeError{
		Number: -1,
	}
	if resType == execResultType {
		errChannel = res.errChannel
		sfError.QueryID = res.queryID
	} else {
		errChannel = rows.errChannel
		sfError.QueryID = rows.queryID
	}
	defer close(errChannel)
	token, _, _ := sr.TokenAccessor.GetTokens()
	headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, token)
	resp, err := sr.FuncGet(ctx, sr, URL, headers, timeout)
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to get response. err: %v", err)
		sfError.Message = err.Error()
		errChannel <- sfError
		close(errChannel)
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	respd := execResponse{}
	err = json.NewDecoder(resp.Body).Decode(&respd)
	resp.Body.Close()
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to decode JSON. err: %v", err)
		sfError.Message = err.Error()
		errChannel <- sfError
		close(errChannel)
		return
	}

	sc := &snowflakeConn{rest: sr, cfg: cfg}
	if respd.Success {
		if resType == execResultType {
			res.insertID = -1
			if isDml(respd.Data.StatementTypeID) {
				res.affectedRows, _ = updateRows(respd.Data)
			} else if isMultiStmt(&respd.Data) {
				r, err := sc.handleMultiExec(ctx, respd.Data)
				if err != nil {
					res.errChannel <- err
					close(errChannel)
					return
				}
				res.affectedRows, err = r.RowsAffected()
				if err != nil {
					res.errChannel <- err
					close(errChannel)
					return
				}
			}
			res.queryID = respd.Data.QueryID
			res.errChannel <- nil // mark exec status complete
		} else {
			rows.sc = sc
			rows.queryID = respd.Data.QueryID
			if isMultiStmt(&respd.Data) {
				err = sc.handleMultiQuery(ctx, respd.Data, rows)
				if err != nil {
					rows.errChannel <- err
					close(errChannel)
					return
				}
			} else {
				rows.addDownloader(populateChunkDownloader(ctx, sc, respd.Data))
			}
			rows.ChunkDownloader.start()
			rows.errChannel <- nil // mark query status complete
		}
	} else {
		var code int
		if respd.Code != "" {
			code, err = strconv.Atoi(respd.Code)
			if err != nil {
				code = -1
			}
		} else {
			code = -1
		}
		errChannel <- &SnowflakeError{
			Number:   code,
			SQLState: respd.Data.SQLState,
			Message:  respd.Message,
			QueryID:  respd.Data.QueryID,
		}
	}
}

func buildSnowflakeConn(ctx context.Context, config Config) (*snowflakeConn, error) {
	sc := &snowflakeConn{
		SequenceCounter: 0,
		ctx:             ctx,
		cfg:             &config,
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
	if strings.HasSuffix(sc.cfg.Host, "privatelink.snowflakecomputing.com") {
		if err := sc.setupOCSPPrivatelink(sc.cfg.Application, sc.cfg.Host); err != nil {
			return nil, err
		}
	} else {
		if _, set := os.LookupEnv(cacheServerURLEnv); set {
			os.Unsetenv(cacheServerURLEnv)
		}
	}
	var tokenAccessor TokenAccessor
	if sc.cfg.TokenAccessor != nil {
		tokenAccessor = sc.cfg.TokenAccessor
	} else {
		tokenAccessor = getSimpleTokenAccessor()
	}
	if sc.cfg.DisableTelemetry {
		sc.telemetry.enabled = false
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
		TokenAccessor:       tokenAccessor,
		LoginTimeout:        sc.cfg.LoginTimeout,
		RequestTimeout:      sc.cfg.RequestTimeout,
		FuncPost:            postRestful,
		FuncGet:             getRestful,
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
	sc.telemetry = &snowflakeTelemetry{
		flushSize: defaultFlushSize,
		sr:        sc.rest,
		mutex:     &sync.Mutex{},
		enabled:   true,
	}
	return sc, nil
}
