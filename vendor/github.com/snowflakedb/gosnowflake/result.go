// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import "errors"

type queryStatus string

const (
	// QueryStatusInProgress denotes a query execution in progress
	QueryStatusInProgress queryStatus = "queryStatusInProgress"
	// QueryStatusComplete denotes a completed query execution
	QueryStatusComplete queryStatus = "queryStatusComplete"
	// QueryFailed denotes a failed query
	QueryFailed queryStatus = "queryFailed"
)

// SnowflakeResult provides an API for methods exposed to the clients
type SnowflakeResult interface {
	GetQueryID() string
	GetStatus() queryStatus
	GetArrowBatches() ([]*ArrowBatch, error)
}

type snowflakeResult struct {
	affectedRows int64
	insertID     int64 // Snowflake doesn't support last insert id
	queryID      string
	status       queryStatus
	err          error
	errChannel   chan error
}

func (res *snowflakeResult) LastInsertId() (int64, error) {
	if err := res.waitForAsyncExecStatus(); err != nil {
		return -1, err
	}
	return res.insertID, nil
}

func (res *snowflakeResult) RowsAffected() (int64, error) {
	if err := res.waitForAsyncExecStatus(); err != nil {
		return -1, err
	}
	return res.affectedRows, nil
}

func (res *snowflakeResult) GetQueryID() string {
	return res.queryID
}

func (res *snowflakeResult) GetStatus() queryStatus {
	return res.status
}

func (res *snowflakeResult) GetArrowBatches() ([]*ArrowBatch, error) {
	return nil, &SnowflakeError{
		Number:  ErrNotImplemented,
		Message: errMsgNotImplemented,
	}
}

func (res *snowflakeResult) waitForAsyncExecStatus() error {
	// if async exec, block until execution is finished
	if res.status == QueryStatusInProgress {
		err := <-res.errChannel
		res.status = QueryStatusComplete
		if err != nil {
			res.status = QueryFailed
			res.err = err
			return err
		}
	} else if res.status == QueryFailed {
		return res.err
	}
	return nil
}

type snowflakeResultNoRows struct {
	queryID string
}

func (*snowflakeResultNoRows) LastInsertId() (int64, error) {
	return 0, errors.New("no LastInsertId available")
}

func (*snowflakeResultNoRows) RowsAffected() (int64, error) {
	return 0, errors.New("no RowsAffected available")
}

func (rnr *snowflakeResultNoRows) GetQueryID() string {
	return rnr.queryID
}
