// Copyright (c) 2017-2021 Snowflake Computing Inc. All right reserved.

package gosnowflake

type queryStatus string

const (
	// QueryStatusWaiting denotes a query execution waiting to happen
	QueryStatusWaiting queryStatus = "queryStatusWaiting"
	// QueryStatusInProgress denotes a query execution in progress
	QueryStatusInProgress queryStatus = "queryStatusInProgress"
	// QueryStatusComplete denotes a completed query execution
	QueryStatusComplete queryStatus = "queryStatusComplete"
	// QueryFailed denotes a failed query
	QueryFailed queryStatus = "queryFailed"
)

// SnowflakeResult provides the associated query ID
type SnowflakeResult interface {
	GetQueryID() string
	GetStatus() queryStatus
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
