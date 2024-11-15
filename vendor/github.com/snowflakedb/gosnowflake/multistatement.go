// Copyright (c) 2021-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

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

func (sc *snowflakeConn) handleMultiExec(
	ctx context.Context,
	data execResponseData) (
	driver.Result, error) {
	if data.ResultIDs == "" {
		return nil, (&SnowflakeError{
			Number:   ErrNoResultIDs,
			SQLState: data.SQLState,
			Message:  errMsgNoResultIDs,
			QueryID:  data.QueryID,
		}).exceptionTelemetry(sc)
	}
	var updatedRows int64
	childResults := getChildResults(data.ResultIDs, data.ResultTypes)
	for _, child := range childResults {
		resultPath := fmt.Sprintf(urlQueriesResultFmt, child.id)
		childResultType, err := strconv.ParseInt(child.typ, 10, 64)
		if err != nil {
			return nil, err
		}
		if isDml(childResultType) {
			childData, err := sc.getQueryResultResp(ctx, resultPath)
			if err != nil {
				logger.WithContext(ctx).Errorf("error: %v", err)
				return nil, err
			}
			if childData != nil && !childData.Success {
				code, err := strconv.Atoi(childData.Code)
				if err != nil {
					return nil, err
				}
				return nil, (&SnowflakeError{
					Number:   code,
					SQLState: childData.Data.SQLState,
					Message:  childData.Message,
					QueryID:  childData.Data.QueryID,
				}).exceptionTelemetry(sc)
			}
			count, err := updateRows(childData.Data)
			if err != nil {
				logger.WithContext(ctx).Errorf("error: %v", err)
				return nil, err
			}
			updatedRows += count
		}
	}
	logger.WithContext(ctx).Infof("number of updated rows: %#v", updatedRows)
	return &snowflakeResult{
		affectedRows: updatedRows,
		insertID:     -1,
		queryID:      data.QueryID,
	}, nil
}

// Fill the corresponding rows and add chunk downloader into the rows when
// iterating across the childResults
func (sc *snowflakeConn) handleMultiQuery(
	ctx context.Context,
	data execResponseData,
	rows *snowflakeRows) error {
	if data.ResultIDs == "" {
		return (&SnowflakeError{
			Number:   ErrNoResultIDs,
			SQLState: data.SQLState,
			Message:  errMsgNoResultIDs,
			QueryID:  data.QueryID,
		}).exceptionTelemetry(sc)
	}
	childResults := getChildResults(data.ResultIDs, data.ResultTypes)
	for _, child := range childResults {
		if err := sc.rowsForRunningQuery(ctx, child.id, rows); err != nil {
			return err
		}
	}
	return nil
}
