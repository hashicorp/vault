// Copyright (c) 2021-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const urlQueriesResultFmt = "/queries/%s/result"

// queryResultStatus is status returned from server
type queryResultStatus int

// Query Status defined at server side
const (
	SFQueryRunning queryResultStatus = iota
	SFQueryAborting
	SFQuerySuccess
	SFQueryFailedWithError
	SFQueryAborted
	SFQueryQueued
	SFQueryFailedWithIncident
	SFQueryDisconnected
	SFQueryResumingWarehouse
	// SFQueryQueueRepairingWarehouse present in QueryDTO.java.
	SFQueryQueueRepairingWarehouse
	SFQueryRestarted
	// SFQueryBlocked is when a statement is waiting on a lock on resource held
	// by another statement.
	SFQueryBlocked
	SFQueryNoData
)

func (qs queryResultStatus) String() string {
	return [...]string{"RUNNING", "ABORTING", "SUCCESS", "FAILED_WITH_ERROR",
		"ABORTED", "QUEUED", "FAILED_WITH_INCIDENT", "DISCONNECTED",
		"RESUMING_WAREHOUSE", "QUEUED_REPAIRING_WAREHOUSE", "RESTARTED",
		"BLOCKED", "NO_DATA"}[qs]
}

func (qs queryResultStatus) isRunning() bool {
	switch qs {
	case SFQueryRunning, SFQueryResumingWarehouse, SFQueryQueued,
		SFQueryQueueRepairingWarehouse, SFQueryNoData:
		return true
	default:
		return false
	}
}

func (qs queryResultStatus) isError() bool {
	switch qs {
	case SFQueryAborting, SFQueryFailedWithError, SFQueryAborted,
		SFQueryFailedWithIncident, SFQueryDisconnected, SFQueryBlocked:
		return true
	default:
		return false
	}
}

var strQueryStatusMap = map[string]queryResultStatus{"RUNNING": SFQueryRunning,
	"ABORTING": SFQueryAborting, "SUCCESS": SFQuerySuccess,
	"FAILED_WITH_ERROR": SFQueryFailedWithError, "ABORTED": SFQueryAborted,
	"QUEUED": SFQueryQueued, "FAILED_WITH_INCIDENT": SFQueryFailedWithIncident,
	"DISCONNECTED":               SFQueryDisconnected,
	"RESUMING_WAREHOUSE":         SFQueryResumingWarehouse,
	"QUEUED_REPAIRING_WAREHOUSE": SFQueryQueueRepairingWarehouse,
	"RESTARTED":                  SFQueryRestarted,
	"BLOCKED":                    SFQueryBlocked, "NO_DATA": SFQueryNoData}

type retStatus struct {
	Status       string   `json:"status"`
	SQLText      string   `json:"sqlText"`
	StartTime    int64    `json:"startTime"`
	EndTime      int64    `json:"endTime"`
	ErrorCode    string   `json:"errorCode"`
	ErrorMessage string   `json:"errorMessage"`
	Stats        retStats `json:"stats"`
}

type retStats struct {
	ScanBytes    int64 `json:"scanBytes"`
	ProducedRows int64 `json:"producedRows"`
}

type statusResponse struct {
	Data struct {
		Queries []retStatus `json:"queries"`
	} `json:"data"`
	Message string `json:"message"`
	Code    string `json:"code"`
	Success bool   `json:"success"`
}

func strToQueryStatus(in string) queryResultStatus {
	return strQueryStatusMap[in]
}

// SnowflakeQueryStatus is the query status metadata of a snowflake query
type SnowflakeQueryStatus struct {
	SQLText      string
	StartTime    int64
	EndTime      int64
	ErrorCode    string
	ErrorMessage string
	ScanBytes    int64
	ProducedRows int64
}

// SnowflakeConnection is a wrapper to snowflakeConn that exposes API functions
type SnowflakeConnection interface {
	GetQueryStatus(ctx context.Context, queryID string) (*SnowflakeQueryStatus, error)
}

// checkQueryStatus returns the status given the query ID. If successful,
// the error will be nil, indicating there is a complete query result to fetch.
// Other than nil, there are three error types that can be returned:
// 1. ErrQueryStatus, if GS cannot return any kind of status due to any reason,
// i.e. connection, permission, if a query was just submitted, etc.
// 2, ErrQueryReportedError, if the requested query was terminated or aborted
// and GS returned an error status included in query. SFQueryFailedWithError
// 3, ErrQueryIsRunning, if the requested query is still running and might have
// a complete result later, these statuses were listed in query. SFQueryRunning
func (sc *snowflakeConn) checkQueryStatus(
	ctx context.Context,
	qid string) (
	*retStatus, error) {
	headers := make(map[string]string)
	param := make(url.Values)
	param.Set(requestGUIDKey, NewUUID().String())
	if tok, _, _ := sc.rest.TokenAccessor.GetTokens(); tok != "" {
		headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, tok)
	}
	resultPath := fmt.Sprintf("%s/%s", monitoringQueriesPath, qid)
	url := sc.rest.getFullURL(resultPath, &param)

	res, err := sc.rest.FuncGet(ctx, sc.rest, url, headers, sc.rest.RequestTimeout)
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to get response. err: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	var statusResp = statusResponse{}
	if err = json.NewDecoder(res.Body).Decode(&statusResp); err != nil {
		logger.WithContext(ctx).Errorf("failed to decode JSON. err: %v", err)
		return nil, err
	}

	if !statusResp.Success || len(statusResp.Data.Queries) == 0 {
		logger.WithContext(ctx).Errorf("status query returned not-success or no status returned.")
		return nil, (&SnowflakeError{
			Number:  ErrQueryStatus,
			Message: "status query returned not-success or no status returned. Please retry",
		}).exceptionTelemetry(sc)
	}

	queryRet := statusResp.Data.Queries[0]
	if queryRet.ErrorCode != "" {
		return &queryRet, (&SnowflakeError{
			Number:         ErrQueryStatus,
			Message:        errMsgQueryStatus,
			MessageArgs:    []interface{}{queryRet.ErrorCode, queryRet.ErrorMessage},
			IncludeQueryID: true,
			QueryID:        qid,
		}).exceptionTelemetry(sc)
	}

	// returned errorCode is 0. Now check what is the returned status of the query.
	qStatus := strToQueryStatus(queryRet.Status)
	if qStatus.isError() {
		return &queryRet, (&SnowflakeError{
			Number: ErrQueryReportedError,
			Message: fmt.Sprintf("%s: status from server: [%s]",
				queryRet.ErrorMessage, queryRet.Status),
			IncludeQueryID: true,
			QueryID:        qid,
		}).exceptionTelemetry(sc)
	}

	if qStatus.isRunning() {
		return &queryRet, (&SnowflakeError{
			Number: ErrQueryIsRunning,
			Message: fmt.Sprintf("%s: status from server: [%s]",
				queryRet.ErrorMessage, queryRet.Status),
			IncludeQueryID: true,
			QueryID:        qid,
		}).exceptionTelemetry(sc)
	}
	//success
	return &queryRet, nil
}

func (sc *snowflakeConn) getQueryResultResp(
	ctx context.Context,
	resultPath string) (
	*execResponse, error) {
	headers := getHeaders()
	paramsMutex.Lock()
	if serviceName, ok := sc.cfg.Params[serviceName]; ok {
		headers[httpHeaderServiceName] = *serviceName
	}
	paramsMutex.Unlock()
	param := make(url.Values)
	param.Set(requestIDKey, getOrGenerateRequestIDFromContext(ctx).String())
	param.Set("clientStartTime", strconv.FormatInt(sc.currentTimeProvider.currentTime(), 10))
	param.Set(requestGUIDKey, NewUUID().String())
	token, _, _ := sc.rest.TokenAccessor.GetTokens()
	if token != "" {
		headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, token)
	}
	url := sc.rest.getFullURL(resultPath, &param)

	respd, err := getQueryResultWithRetriesForAsyncMode(ctx, sc.rest, url, headers, sc.rest.RequestTimeout)
	if err != nil {
		logger.WithContext(ctx).Errorf("error: %v", err)
		return nil, err
	}
	return respd, nil
}

// Fetch query result for a query id from /queries/<qid>/result endpoint.
func (sc *snowflakeConn) rowsForRunningQuery(
	ctx context.Context, qid string,
	rows *snowflakeRows) error {
	resultPath := fmt.Sprintf(urlQueriesResultFmt, qid)
	resp, err := sc.getQueryResultResp(ctx, resultPath)
	if err != nil {
		logger.WithContext(ctx).Errorf("error: %v", err)
		return err
	}

	if !resp.Success {
		code, err := strconv.Atoi(resp.Code)
		if err != nil {
			return err
		}
		return (&SnowflakeError{
			Number:   code,
			SQLState: resp.Data.SQLState,
			Message:  resp.Message,
			QueryID:  resp.Data.QueryID,
		}).exceptionTelemetry(sc)
	}
	rows.addDownloader(populateChunkDownloader(ctx, sc, resp.Data))
	return nil
}

// prepare a Rows object to return for query of 'qid'
func (sc *snowflakeConn) buildRowsForRunningQuery(
	ctx context.Context,
	qid string) (
	driver.Rows, error) {
	rows := new(snowflakeRows)
	rows.sc = sc
	rows.queryID = qid
	rows.ctx = ctx
	if err := sc.rowsForRunningQuery(ctx, qid, rows); err != nil {
		return nil, err
	}
	err := rows.ChunkDownloader.start()
	return rows, err
}
