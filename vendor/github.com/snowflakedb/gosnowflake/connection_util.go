// Copyright (c) 2021 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

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

func (sc *snowflakeConn) getArrayBindStageThreshold() int {
	v, ok := sc.cfg.Params[sessionArrayBindStageThreshold]
	if !ok {
		return 0
	}
	num, err := strconv.Atoi(*v)
	if err != nil {
		return 0
	}
	return num
}

func (sc *snowflakeConn) connectionTelemetry(cfg *Config) {
	data := &telemetryData{
		Message: map[string]string{
			typeKey:          connectionParameters,
			driverTypeKey:    "Go",
			driverVersionKey: SnowflakeGoDriverVersion,
		},
		Timestamp: time.Now().UnixNano(),
	}
	for k, v := range cfg.Params {
		data.Message[k] = *v
	}
	sc.telemetry.addLog(data)
	sc.telemetry.sendBatch()
}

// processFileTransfer creates a snowflakeFileTransferAgent object to process
// any PUT/GET commands with their specified options
func (sc *snowflakeConn) processFileTransfer(
	ctx context.Context,
	data *execResponse,
	query string,
	isInternal bool) (
	*execResponse, error) {
	sfa := snowflakeFileTransferAgent{
		sc:      sc,
		data:    &data.Data,
		command: query,
		options: new(SnowflakeFileTransferOptions),
	}
	if fs := getFileStream(ctx); fs != nil {
		sfa.sourceStream = fs
		if isInternal {
			sfa.data.AutoCompress = false
		}
	}
	if op := getFileTransferOptions(ctx); op != nil {
		sfa.options = op
	}
	if sfa.options.MultiPartThreshold == 0 {
		sfa.options.MultiPartThreshold = dataSizeThreshold
	}
	if err := sfa.execute(); err != nil {
		return nil, err
	}
	data, err := sfa.result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getFileStream(ctx context.Context) *bytes.Buffer {
	s := ctx.Value(fileStreamFile)
	r, ok := s.(io.Reader)
	if !ok {
		return nil
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf
}

func getFileTransferOptions(ctx context.Context) *SnowflakeFileTransferOptions {
	v := ctx.Value(fileTransferOptions)
	if v == nil {
		return nil
	}
	o, ok := v.(*SnowflakeFileTransferOptions)
	if !ok {
		return nil
	}
	return o
}

func (sc *snowflakeConn) populateSessionParameters(parameters []nameValueParameter) {
	// other session parameters (not all)
	logger.WithContext(sc.ctx).Infof("params: %#v", parameters)
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
		logger.Debugf("parameter. name: %v, value: %v", param.Name, v)
		sc.cfg.Params[strings.ToLower(param.Name)] = &v
	}
}

func isAsyncMode(ctx context.Context) bool {
	val := ctx.Value(asyncMode)
	if val == nil {
		return false
	}
	a, ok := val.(bool)
	return ok && a
}

func isDescribeOnly(ctx context.Context) bool {
	v := ctx.Value(describeOnly)
	if v == nil {
		return false
	}
	d, ok := v.(bool)
	return ok && d
}

func setResultType(ctx context.Context, resType resultType) context.Context {
	return context.WithValue(ctx, snowflakeResultType, resType)
}

func getResultType(ctx context.Context) resultType {
	return ctx.Value(snowflakeResultType).(resultType)
}

// isDml returns true if the statement type code is in the range of DML.
func isDml(v int64) bool {
	return statementTypeIDDml <= v && v <= statementTypeIDMultiTableInsert
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

// isMultiStmt returns true if the statement code is of type multistatement
// Note that the statement type code is also equivalent to type INSERT, so an
// additional check of the name is required
func isMultiStmt(data *execResponseData) bool {
	return data.StatementTypeID == statementTypeIDMulti &&
		data.RowType[0].Name == "multiple statement execution"
}

func getResumeQueryID(ctx context.Context) (string, error) {
	val := ctx.Value(fetchResultByID)
	if val == nil {
		return "", nil
	}
	strVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("failed to cast val %+v to string", val)
	}
	// so there is a queryID in context for which we want to fetch the result
	if !queryIDRegexp.MatchString(strVal) {
		return strVal, &SnowflakeError{
			Number:  ErrQueryIDFormat,
			Message: "Invalid QID",
			QueryID: strVal}
	}
	return strVal, nil
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

// returns snowflake chunk downloader by default or stream based chunk
// downloader if option provided through context
func populateChunkDownloader(
	ctx context.Context,
	sc *snowflakeConn,
	data execResponseData) chunkDownloader {
	if useStreamDownloader(ctx) && data.QueryResultFormat == "json" {
		// stream chunk downloading only works for row based data formats, i.e. json
		fetcher := &httpStreamChunkFetcher{
			ctx:      ctx,
			client:   sc.rest.Client,
			clientIP: sc.cfg.ClientIP,
			headers:  data.ChunkHeaders,
			qrmk:     data.Qrmk,
		}
		return newStreamChunkDownloader(ctx, fetcher, data.Total, data.RowType,
			data.RowSet, data.Chunks)
	}

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
		RowSet: rowSetType{
			RowType:      data.RowType,
			JSON:         data.RowSet,
			RowSetBase64: data.RowSetBase64,
		},
	}
}

func (sc *snowflakeConn) setupOCSPPrivatelink(app string, host string) error {
	ocspCacheServer := fmt.Sprintf("http://ocsp.%v/ocsp_response_cache.json", host)
	if err := os.Setenv(cacheServerURLEnv, ocspCacheServer); err != nil {
		return err
	}
	ocspRetryHost := fmt.Sprintf("http://ocsp.%v/retry/", host) + "%v/%v"
	if err := os.Setenv(ocspRetryURLEnv, ocspRetryHost); err != nil {
		return err
	}
	return nil
}
