// Copyright (c) 2017-2020 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	statementTypeIDMulti = int64(0x1000)

	statementTypeIDDml              = int64(0x3000)
	statementTypeIDInsert           = statementTypeIDDml + int64(0x100)
	statementTypeIDUpdate           = statementTypeIDDml + int64(0x200)
	statementTypeIDDelete           = statementTypeIDDml + int64(0x300)
	statementTypeIDMerge            = statementTypeIDDml + int64(0x400)
	statementTypeIDMultiTableInsert = statementTypeIDDml + int64(0x500)
)

const (
	sessionClientSessionKeepAlive          = "client_session_keep_alive"
	sessionClientValidateDefaultParameters = "CLIENT_VALIDATE_DEFAULT_PARAMETERS"
	serviceName                            = "service_name"
)

type snowflakeConn struct {
	cfg             *Config
	rest            *snowflakeRestful
	SequenceCounter uint64
	QueryID         string
	SQLState        string
}

// isDml returns true if the statement type code is in the range of DML.
func (sc *snowflakeConn) isDml(v int64) bool {
	switch v {
	case statementTypeIDDml, statementTypeIDInsert,
		statementTypeIDUpdate, statementTypeIDDelete,
		statementTypeIDMerge, statementTypeIDMultiTableInsert:
		return true
	}
	return false
}

// isMultiStmt returns true if the statement type code is of type multistatement
// Note that the statement type code is also equivalent to type INSERT, so an additional check of the name is required
func (sc *snowflakeConn) isMultiStmt(data execResponseData) bool {
	return data.StatementTypeID == statementTypeIDMulti && data.RowType[0].Name == "multiple statement execution"
}

func (sc *snowflakeConn) exec(
	ctx context.Context,
	query string,
	noResult bool,
	isInternal bool,
	bindings []driver.NamedValue) (
	*execResponse, error) {
	var err error
	counter := atomic.AddUint64(&sc.SequenceCounter, 1) // query sequence counter

	req := execRequest{
		SQLText:    query,
		AsyncExec:  noResult,
		SequenceID: counter,
	}
	req.IsInternal = isInternal
	tsmode := "TIMESTAMP_NTZ"
	idx := 1
	if len(bindings) > 0 {
		req.Bindings = make(map[string]execBindParameter, len(bindings))
		for i, n := 0, len(bindings); i < n; i++ {
			t := goTypeToSnowflake(bindings[i].Value, tsmode)
			glog.V(2).Infof("tmode: %v\n", t)
			if t == "CHANGE_TYPE" {
				tsmode, err = dataTypeMode(bindings[i].Value)
				if err != nil {
					return nil, err
				}
			} else {
				var v1 interface{}
				if t == "ARRAY" {
					t, v1 = arrayToString(bindings[i].Value)
				} else {
					v1, err = valueToString(bindings[i].Value, tsmode)
				}
				if err != nil {
					return nil, err
				}
				req.Bindings[strconv.Itoa(idx)] = execBindParameter{
					Type:  t,
					Value: v1,
				}
				idx++
			}
		}
	}
	multiCount := ctx.Value(MultiStatementCount)
	if multiCount != nil {
		req.Parameters = map[string]interface{}{string(MultiStatementCount): multiCount}
	}
	glog.V(2).Infof("bindings: %v", req.Bindings)
	glog.V(2).Infof("parameters: %v", req.Parameters)

	headers := make(map[string]string)
	headers["Content-Type"] = headerContentTypeApplicationJSON
	headers["accept"] = headerAcceptTypeApplicationSnowflake // TODO v1.1: change to JSON in case of PUT/GET
	headers["User-Agent"] = userAgent
	if serviceName, ok := sc.cfg.Params[serviceName]; ok {
		headers["X-Snowflake-Service"] = *serviceName
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var data *execResponse

	requestID := uuid.New()
	data, err = sc.rest.FuncPostQuery(ctx, sc.rest, &url.Values{}, headers, jsonBody, sc.rest.RequestTimeout, &requestID)
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
	glog.V(2).Infof("Success: %v, Code: %v", data.Success, code)
	if !data.Success {
		return nil, &SnowflakeError{
			Number:   code,
			SQLState: data.Data.SQLState,
			Message:  data.Message,
			QueryID:  data.Data.QueryID,
		}
	}
	glog.V(2).Info("Exec/Query SUCCESS")
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
	return sc.BeginTx(context.TODO(), driver.TxOptions{})
}

func (sc *snowflakeConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	glog.V(2).Info("BeginTx")
	if opts.ReadOnly {
		return nil, &SnowflakeError{
			Number:   ErrNoReadOnlyTransaction,
			SQLState: SQLStateFeatureNotSupported,
			Message:  errMsgNoReadOnlyTransaction,
		}
	}
	if int(opts.Isolation) != int(sql.LevelDefault) {
		return nil, &SnowflakeError{
			Number:   ErrNoDefaultTransactionIsolationLevel,
			SQLState: SQLStateFeatureNotSupported,
			Message:  errMsgNoDefaultTransactionIsolationLevel,
		}
	}
	if sc.rest == nil {
		return nil, driver.ErrBadConn
	}
	_, err := sc.exec(ctx, "BEGIN", false, false, nil)
	if err != nil {
		return nil, err
	}
	return &snowflakeTx{sc}, err
}

func (sc *snowflakeConn) cleanup() {
	glog.Flush() // must flush log buffer while the process is running.
	sc.rest = nil
	sc.cfg = nil
}

func (sc *snowflakeConn) Close() (err error) {
	glog.V(2).Infoln("Close")
	sc.stopHeartBeat()

	err = sc.rest.FuncCloseSession(context.TODO(), sc.rest, sc.rest.RequestTimeout)
	if err != nil {
		glog.V(2).Info(err)
	}
	sc.cleanup()
	return nil
}

func (sc *snowflakeConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	glog.V(2).Infoln("Prepare")
	if sc.rest == nil {
		return nil, driver.ErrBadConn
	}
	stmt := &snowflakeStmt{
		sc:    sc,
		query: query,
	}
	return stmt, nil
}

func (sc *snowflakeConn) Prepare(query string) (driver.Stmt, error) {
	return sc.PrepareContext(context.TODO(), query)
}

func (sc *snowflakeConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	glog.V(2).Infof("Exec: %#v, %v", query, args)
	if sc.rest == nil {
		return nil, driver.ErrBadConn
	}
	// TODO: handle noResult and isInternal
	data, err := sc.exec(ctx, query, false, false, args)
	if err != nil {
		glog.V(2).Infof("error: %v", err)
		if data != nil {
			code, err := strconv.Atoi(data.Code)
			if err != nil {
				return nil, err
			}
			return nil, &SnowflakeError{
				Number:   code,
				SQLState: data.Data.SQLState,
				Message:  err.Error(),
				QueryID:  data.Data.QueryID}
		}
		return nil, err
	}

	var updatedRows int64
	if sc.isDml(data.Data.StatementTypeID) {
		// collects all values from the returned row sets
		updatedRows, err = updateRows(data.Data)
		if err != nil {
			return nil, err
		}
		glog.V(2).Infof("number of updated rows: %#v", updatedRows)
		return &snowflakeResult{
			affectedRows: updatedRows,
			insertID:     -1,
			queryID:      sc.QueryID,
		}, nil // last insert id is not supported by Snowflake
	} else if sc.isMultiStmt(data.Data) {
		childResults := getChildResults(data.Data.ResultIDs, data.Data.ResultTypes)
		for _, child := range childResults {
			resultPath := fmt.Sprintf("/queries/%s/result", child.id)
			childData, err := sc.getQueryResult(ctx, resultPath)
			if err != nil {
				glog.V(2).Infof("error: %v", err)
				code, err := strconv.Atoi(childData.Code)
				if err != nil {
					return nil, err
				}
				if childData != nil {
					return nil, &SnowflakeError{
						Number:   code,
						SQLState: childData.Data.SQLState,
						Message:  err.Error(),
						QueryID:  childData.Data.QueryID}
				}
				return nil, err
			}
			if sc.isDml(childData.Data.StatementTypeID) {
				count, err := updateRows(childData.Data)
				if err != nil {
					glog.V(2).Infof("error: %v", err)
					if childData != nil {
						code, err := strconv.Atoi(childData.Code)
						if err != nil {
							return nil, err
						}
						return nil, &SnowflakeError{
							Number:   code,
							SQLState: childData.Data.SQLState,
							Message:  err.Error(),
							QueryID:  childData.Data.QueryID}
					}
					return nil, err
				}
				updatedRows += count
			}
		}
		glog.V(2).Infof("number of updated rows: %#v", updatedRows)
		return &snowflakeResult{
			affectedRows: updatedRows,
			insertID:     -1,
			queryID:      sc.QueryID,
		}, nil
	}
	glog.V(2).Info("DDL")
	return driver.ResultNoRows, nil
}

func (sc *snowflakeConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	glog.V(2).Infof("Query: %#v, %v", query, args)
	if sc.rest == nil {
		return nil, driver.ErrBadConn
	}
	// TODO: handle noResult and isInternal
	data, err := sc.exec(ctx, query, false, false, args)
	if err != nil {
		glog.V(2).Infof("error: %v", err)
		if data != nil {
			code, err := strconv.Atoi(data.Code)
			if err != nil {
				return nil, err
			}
			return nil, &SnowflakeError{
				Number:   code,
				SQLState: data.Data.SQLState,
				Message:  err.Error(),
				QueryID:  data.Data.QueryID}
		}
		return nil, err
	}

	rows := new(snowflakeRows)
	rows.sc = sc
	rows.RowType = data.Data.RowType
	rows.ChunkDownloader = &snowflakeChunkDownloader{
		sc:                 sc,
		ctx:                ctx,
		CurrentChunk:       make([]chunkRowType, len(data.Data.RowSet)),
		ChunkMetas:         data.Data.Chunks,
		Total:              data.Data.Total,
		TotalRowIndex:      int64(-1),
		CellCount:          len(data.Data.RowType),
		Qrmk:               data.Data.Qrmk,
		QueryResultFormat:  data.Data.QueryResultFormat,
		ChunkHeader:        data.Data.ChunkHeaders,
		FuncDownload:       downloadChunk,
		FuncDownloadHelper: downloadChunkHelper,
		FuncGet:            getChunk,
		RowSet: rowSetType{RowType: data.Data.RowType,
			JSON:         data.Data.RowSet,
			RowSetBase64: data.Data.RowSetBase64,
		},
	}
	rows.queryID = sc.QueryID

	if sc.isMultiStmt(data.Data) {
		childResults := getChildResults(data.Data.ResultIDs, data.Data.ResultTypes)
		var nextChunkDownloader *snowflakeChunkDownloader
		firstResultSet := false

		for _, child := range childResults {
			resultPath := fmt.Sprintf("/queries/%s/result", child.id)
			childData, err := sc.getQueryResult(ctx, resultPath)
			if err != nil {
				glog.V(2).Infof("error: %v", err)
				if childData != nil {
					code, err := strconv.Atoi(childData.Code)
					if err != nil {
						return nil, err
					}
					return nil, &SnowflakeError{
						Number:   code,
						SQLState: childData.Data.SQLState,
						Message:  err.Error(),
						QueryID:  childData.Data.QueryID}
				}
				return nil, err
			}
			if !firstResultSet {
				// populate rows.ChunkDownloader with the first child
				rows.ChunkDownloader = populateChunkDownloader(ctx, sc, childData.Data)
				nextChunkDownloader = rows.ChunkDownloader
				firstResultSet = true
			} else {
				nextChunkDownloader.NextDownloader = populateChunkDownloader(ctx, sc, childData.Data)
				nextChunkDownloader = nextChunkDownloader.NextDownloader
			}
		}
	}

	rows.ChunkDownloader.start()
	return rows, err
}

func (sc *snowflakeConn) Exec(
	query string,
	args []driver.Value) (
	driver.Result, error) {
	return sc.ExecContext(context.TODO(), query, toNamedValues(args))
}

func (sc *snowflakeConn) Query(
	query string,
	args []driver.Value) (
	driver.Rows, error) {
	return sc.QueryContext(context.TODO(), query, toNamedValues(args))
}

func (sc *snowflakeConn) Ping(ctx context.Context) error {
	glog.V(2).Infoln("Ping")
	if sc.rest == nil {
		return driver.ErrBadConn
	}
	// TODO: handle noResult and isInternal
	_, err := sc.exec(ctx, "SELECT 1", false, false, []driver.NamedValue{})
	return err
}

func (sc *snowflakeConn) CheckNamedValue(nv *driver.NamedValue) error {
	switch reflect.TypeOf(nv.Value) {
	case reflect.TypeOf([]int{0}), reflect.TypeOf([]int64{0}), reflect.TypeOf([]float64{0}),
		reflect.TypeOf([]bool{false}), reflect.TypeOf([]string{""}):
		return nil
	default:
		return driver.ErrSkip
	}
}

func (sc *snowflakeConn) populateSessionParameters(parameters []nameValueParameter) {
	// other session parameters (not all)
	glog.V(2).Infof("params: %#v", parameters)
	for _, param := range parameters {
		v := ""
		switch param.Value.(type) {
		case int64:
			if vv, ok := param.Value.(int64); ok {
				v = strconv.FormatInt(vv, 10)
			}
		case float64:
			if vv, ok := param.Value.(float64); ok {
				v = strconv.FormatFloat(vv, 'g', -1, 64)
			}
		case bool:
			if vv, ok := param.Value.(bool); ok {
				v = strconv.FormatBool(vv)
			}
		default:
			if vv, ok := param.Value.(string); ok {
				v = vv
			}
		}
		glog.V(3).Infof("parameter. name: %v, value: %v", param.Name, v)
		sc.cfg.Params[strings.ToLower(param.Name)] = &v
	}
}

func (sc *snowflakeConn) isClientSessionKeepAliveEnabled() bool {
	v, ok := sc.cfg.Params[sessionClientSessionKeepAlive]
	if !ok {
		return false
	}
	return strings.Compare(*v, "true") == 0
}

func (sc *snowflakeConn) startHeartBeat() {
	if !sc.isClientSessionKeepAliveEnabled() {
		return
	}
	sc.rest.HeartBeat = &heartbeat{
		restful: sc.rest,
	}
	sc.rest.HeartBeat.start()
}

func (sc *snowflakeConn) stopHeartBeat() {
	if !sc.isClientSessionKeepAliveEnabled() {
		return
	}
	sc.rest.HeartBeat.stop()
}

func updateRows(data execResponseData) (int64, error) {
	var count int64
	for i, n := 0, len(data.RowType); i < n; i++ {
		v, err := strconv.ParseInt(*data.RowSet[0][i], 10, 64)
		if err != nil {
			return -1, err
		}
		count += v
	}
	return count, nil
}

type childResult struct {
	id  string
	typ string
}

func getChildResults(IDs string, types string) []childResult {
	if IDs == "" {
		return nil
	}
	queryIDs := strings.Split(IDs, ",")
	resultTypes := strings.Split(types, ",")
	res := make([]childResult, len(queryIDs))
	for i, id := range queryIDs {
		res[i] = childResult{id, resultTypes[i]}
	}
	return res
}

func (sc *snowflakeConn) getQueryResult(ctx context.Context, resultPath string) (*execResponse, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = headerContentTypeApplicationJSON
	headers["accept"] = headerAcceptTypeApplicationSnowflake
	headers["User-Agent"] = userAgent
	if serviceName, ok := sc.cfg.Params[serviceName]; ok {
		headers["X-Snowflake-Service"] = *serviceName
	}
	param := make(url.Values)
	param.Add(requestIDKey, uuid.New().String())
	param.Add("clientStartTime", strconv.FormatInt(time.Now().Unix(), 10))
	param.Add(requestGUIDKey, uuid.New().String())
	if sc.rest.Token != "" {
		headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, sc.rest.Token)
	}
	url := sc.rest.getFullURL(resultPath, &param)
	res, err := sc.rest.FuncGet(ctx, sc.rest, url, headers, sc.rest.RequestTimeout)
	if err != nil {
		glog.V(1).Infof("failed to get response. err: %v", err)
		glog.Flush()
		return nil, err
	}
	var respd *execResponse
	err = json.NewDecoder(res.Body).Decode(&respd)
	if err != nil {
		glog.V(1).Infof("failed to decode JSON. err: %v", err)
		glog.Flush()
		return nil, err
	}
	return respd, nil
}

func populateChunkDownloader(ctx context.Context, sc *snowflakeConn, data execResponseData) *snowflakeChunkDownloader {
	return &snowflakeChunkDownloader{
		sc:                 sc,
		ctx:                ctx,
		CurrentChunk:       make([]chunkRowType, len(data.RowSet)),
		ChunkMetas:         data.Chunks,
		Total:              data.Total,
		TotalRowIndex:      int64(-1),
		CellCount:          len(data.RowType),
		Qrmk:               data.Qrmk,
		QueryResultFormat:  data.QueryResultFormat,
		ChunkHeader:        data.ChunkHeaders,
		FuncDownload:       downloadChunk,
		FuncDownloadHelper: downloadChunkHelper,
		FuncGet:            getChunk,
		RowSet: rowSetType{RowType: data.RowType,
			JSON:         data.RowSet,
			RowSetBase64: data.RowSetBase64,
		},
	}
}
