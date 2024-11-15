package gocb

import (
	"encoding/json"
	"errors"
	"time"
)

type jsonQueryMetrics struct {
	ElapsedTime   string `json:"elapsedTime"`
	ExecutionTime string `json:"executionTime"`
	ResultCount   uint64 `json:"resultCount"`
	ResultSize    uint64 `json:"resultSize"`
	MutationCount uint64 `json:"mutationCount,omitempty"`
	SortCount     uint64 `json:"sortCount,omitempty"`
	ErrorCount    uint64 `json:"errorCount,omitempty"`
	WarningCount  uint64 `json:"warningCount,omitempty"`
}

type jsonQueryWarning struct {
	Code    uint32 `json:"code"`
	Message string `json:"msg"`
}

type jsonQueryResponse struct {
	RequestID       string             `json:"requestID"`
	ClientContextID string             `json:"clientContextID"`
	Status          QueryStatus        `json:"status"`
	Warnings        []jsonQueryWarning `json:"warnings"`
	Metrics         *jsonQueryMetrics  `json:"metrics,omitempty"`
	Profile         interface{}        `json:"profile"`
	Signature       interface{}        `json:"signature"`
	Prepared        string             `json:"prepared"`
}

// QueryMetrics encapsulates various metrics gathered during a queries execution.
type QueryMetrics struct {
	ElapsedTime   time.Duration
	ExecutionTime time.Duration
	ResultCount   uint64
	ResultSize    uint64
	MutationCount uint64
	SortCount     uint64
	ErrorCount    uint64
	WarningCount  uint64
}

func (metrics *QueryMetrics) fromData(data *jsonQueryMetrics) error {
	elapsedTime, err := time.ParseDuration(data.ElapsedTime)
	if err != nil {
		logDebugf("Failed to parse query metrics elapsed time: %s", err)
	}

	executionTime, err := time.ParseDuration(data.ExecutionTime)
	if err != nil {
		logDebugf("Failed to parse query metrics execution time: %s", err)
	}

	metrics.ElapsedTime = elapsedTime
	metrics.ExecutionTime = executionTime
	metrics.ResultCount = data.ResultCount
	metrics.ResultSize = data.ResultSize
	metrics.MutationCount = data.MutationCount
	metrics.SortCount = data.SortCount
	metrics.ErrorCount = data.ErrorCount
	metrics.WarningCount = data.WarningCount

	return nil
}

// QueryWarning encapsulates any warnings returned by a query.
type QueryWarning struct {
	Code    uint32
	Message string
}

func (warning *QueryWarning) fromData(data jsonQueryWarning) error {
	warning.Code = data.Code
	warning.Message = data.Message

	return nil
}

// QueryMetaData provides access to the meta-data properties of a query result.
type QueryMetaData struct {
	RequestID       string
	ClientContextID string
	Status          QueryStatus
	Metrics         QueryMetrics
	Signature       interface{}
	Warnings        []QueryWarning
	Profile         interface{}

	preparedName string
}

func (meta *QueryMetaData) fromData(data jsonQueryResponse) error {
	metrics := QueryMetrics{}
	if data.Metrics != nil {
		if err := metrics.fromData(data.Metrics); err != nil {
			return err
		}
	}

	warnings := make([]QueryWarning, len(data.Warnings))
	for wIdx, jsonWarning := range data.Warnings {
		err := warnings[wIdx].fromData(jsonWarning)
		if err != nil {
			return err
		}
	}

	meta.RequestID = data.RequestID
	meta.ClientContextID = data.ClientContextID
	meta.Status = data.Status
	meta.Metrics = metrics
	meta.Signature = data.Signature
	meta.Warnings = warnings
	meta.Profile = data.Profile
	meta.preparedName = data.Prepared

	return nil
}

// QueryResultRaw provides raw access to query data.
// VOLATILE: This API is subject to change at any time.
type QueryResultRaw struct {
	reader        queryRowReader
	transactionID string
	nextRowBytes  []byte
}

// NextBytes returns the next row as bytes.
func (qrr *QueryResultRaw) NextBytes() []byte {
	rowBytes := qrr.nextRowBytes
	qrr.nextRowBytes = qrr.reader.NextRow()

	return rowBytes
}

// Err returns any errors that have occurred on the stream
func (qrr *QueryResultRaw) Err() error {
	err := qrr.reader.Err()
	if err != nil {
		if qrr.transactionID != "" {
			return singleQueryErrToTransactionError(err, qrr.transactionID)
		}

		return err
	}

	return nil
}

// Close marks the results as closed, returning any errors that occurred during reading the results.
func (qrr *QueryResultRaw) Close() error {
	err := qrr.reader.Close()
	if err != nil {
		if qrr.transactionID != "" {
			return singleQueryErrToTransactionError(err, qrr.transactionID)
		}

		return err
	}

	return nil
}

// MetaData returns any meta-data that was available from this query as bytes.
func (qrr *QueryResultRaw) MetaData() ([]byte, error) {
	return qrr.reader.MetaData()
}

// QueryResult allows access to the results of a query.
type QueryResult struct {
	reader queryRowReader

	transactionID string
	nextRowBytes  []byte
	rowBytes      []byte
	endpoint      string
}

func newQueryResult(reader queryRowReader) *QueryResult {
	return &QueryResult{
		reader:       reader,
		endpoint:     reader.Endpoint(),
		nextRowBytes: reader.NextRow(),
	}
}

// Raw returns a QueryResultRaw which can be used to access the raw byte data from search queries.
// Calling this function invalidates the underlying QueryResult which will no longer be able to be used.
// VOLATILE: This API is subject to change at any time.
func (r *QueryResult) Raw() *QueryResultRaw {
	vr := &QueryResultRaw{
		reader:        r.reader,
		transactionID: r.transactionID,
		nextRowBytes:  r.nextRowBytes,
	}

	r.reader = nil
	r.transactionID = ""
	r.nextRowBytes = nil
	return vr
}

func (r *QueryResult) peekNext() []byte {
	return r.nextRowBytes
}

// Next assigns the next result from the results into the value pointer, returning whether the read was successful.
func (r *QueryResult) Next() bool {
	if r.reader == nil {
		return false
	}

	if len(r.nextRowBytes) == 0 {
		return false
	}

	r.rowBytes = r.nextRowBytes
	r.nextRowBytes = r.reader.NextRow()

	return true
}

// Row returns the contents of the current row
func (r *QueryResult) Row(valuePtr interface{}) error {
	if r.reader == nil {
		return r.Err()
	}

	if r.rowBytes == nil {
		return ErrNoResult
	}

	if bytesPtr, ok := valuePtr.(*json.RawMessage); ok {
		*bytesPtr = r.rowBytes
		return nil
	}

	return json.Unmarshal(r.rowBytes, valuePtr)
}

// Err returns any errors that have occurred on the stream
func (r *QueryResult) Err() error {
	if r.reader == nil {
		return errors.New("result object is no longer valid")
	}

	err := r.reader.Err()
	if err != nil {
		if r.transactionID != "" {
			return singleQueryErrToTransactionError(err, r.transactionID)
		}

		return err
	}

	return nil
}

// Close marks the results as closed, returning any errors that occurred during reading the results.
func (r *QueryResult) Close() error {
	if r.reader == nil {
		return r.Err()
	}

	err := r.reader.Close()
	if err != nil {
		if r.transactionID != "" {
			return singleQueryErrToTransactionError(err, r.transactionID)
		}

		return err
	}

	return nil
}

// One assigns the first value from the results into the value pointer.
// It will Close the results but not before iterating through all remaining
// results, as such this should only be used for very small resultsets - ideally
// of, at most, length 1.
func (r *QueryResult) One(valuePtr interface{}) error {
	if r.reader == nil {
		return r.Err()
	}

	// Read the bytes from the first row
	valueBytes := r.nextRowBytes
	if valueBytes == nil {
		return ErrNoResult
	}

	// Skip through the remaining rows
	for r.reader.NextRow() != nil {
		// do nothing with the row
	}
	r.nextRowBytes = nil

	return json.Unmarshal(valueBytes, valuePtr)
}

// MetaData returns any meta-data that was available from this query.  Note that
// the meta-data will only be available once the object has been closed (either
// implicitly or explicitly).
func (r *QueryResult) MetaData() (*QueryMetaData, error) {
	if r.reader == nil {
		return nil, r.Err()
	}

	metaDataBytes, err := r.reader.MetaData()
	if err != nil {
		return nil, err
	}

	var jsonResp jsonQueryResponse
	err = json.Unmarshal(metaDataBytes, &jsonResp)
	if err != nil {
		return nil, err
	}

	var metaData QueryMetaData
	err = metaData.fromData(jsonResp)
	if err != nil {
		return nil, err
	}

	return &metaData, nil
}

// QueryResultInternal provides access to internal only functionality.
// Internal: This should never be used and is not supported.
type QueryResultInternal struct {
	endpoint string
}

// Internal provides access to internal only functionality.
// Internal: This should never be used and is not supported.
func (r *QueryResult) Internal() *QueryResultInternal {
	return &QueryResultInternal{
		endpoint: r.endpoint,
	}
}

// Endpoint returns the endpoint that this query was sent to.
func (r *QueryResultInternal) Endpoint() string {
	return r.endpoint
}

// Query executes the query statement on the server.
func (c *Cluster) Query(statement string, opts *QueryOptions) (*QueryResult, error) {
	return autoOpControl(c.queryController(), func(provider queryProvider) (*QueryResult, error) {
		if opts == nil {
			opts = &QueryOptions{}
		}

		if opts.AsTransaction != nil {
			return c.Transactions().singleQuery(statement, nil, *opts)
		}

		return provider.Query(statement, nil, opts)
	})
}
