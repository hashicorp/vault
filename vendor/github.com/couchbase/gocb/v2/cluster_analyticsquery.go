package gocb

import (
	"encoding/json"
	"errors"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v9"
)

type jsonAnalyticsMetrics struct {
	ElapsedTime      string `json:"elapsedTime"`
	ExecutionTime    string `json:"executionTime"`
	ResultCount      uint64 `json:"resultCount"`
	ResultSize       uint64 `json:"resultSize"`
	MutationCount    uint64 `json:"mutationCount,omitempty"`
	SortCount        uint64 `json:"sortCount,omitempty"`
	ErrorCount       uint64 `json:"errorCount,omitempty"`
	WarningCount     uint64 `json:"warningCount,omitempty"`
	ProcessedObjects uint64 `json:"processedObjects,omitempty"`
}

type jsonAnalyticsWarning struct {
	Code    uint32 `json:"code"`
	Message string `json:"msg"`
}

type jsonAnalyticsResponse struct {
	RequestID       string                 `json:"requestID"`
	ClientContextID string                 `json:"clientContextID"`
	Status          string                 `json:"status"`
	Warnings        []jsonAnalyticsWarning `json:"warnings"`
	Metrics         jsonAnalyticsMetrics   `json:"metrics"`
	Signature       interface{}            `json:"signature"`
}

// AnalyticsMetrics encapsulates various metrics gathered during a queries execution.
type AnalyticsMetrics struct {
	ElapsedTime      time.Duration
	ExecutionTime    time.Duration
	ResultCount      uint64
	ResultSize       uint64
	MutationCount    uint64
	SortCount        uint64
	ErrorCount       uint64
	WarningCount     uint64
	ProcessedObjects uint64
}

func (metrics *AnalyticsMetrics) fromData(data jsonAnalyticsMetrics) error {
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
	metrics.ProcessedObjects = data.ProcessedObjects

	return nil
}

// AnalyticsWarning encapsulates any warnings returned by a query.
type AnalyticsWarning struct {
	Code    uint32
	Message string
}

func (warning *AnalyticsWarning) fromData(data jsonAnalyticsWarning) error {
	warning.Code = data.Code
	warning.Message = data.Message

	return nil
}

// AnalyticsMetaData provides access to the meta-data properties of a query result.
type AnalyticsMetaData struct {
	RequestID       string
	ClientContextID string
	Metrics         AnalyticsMetrics
	Signature       interface{}
	Warnings        []AnalyticsWarning
}

func (meta *AnalyticsMetaData) fromData(data jsonAnalyticsResponse) error {
	metrics := AnalyticsMetrics{}
	if err := metrics.fromData(data.Metrics); err != nil {
		return err
	}

	warnings := make([]AnalyticsWarning, len(data.Warnings))
	for wIdx, jsonWarning := range data.Warnings {
		err := warnings[wIdx].fromData(jsonWarning)
		if err != nil {
			return err
		}
	}

	meta.RequestID = data.RequestID
	meta.ClientContextID = data.ClientContextID
	meta.Metrics = metrics
	meta.Signature = data.Signature
	meta.Warnings = warnings

	return nil
}

// AnalyticsResultRaw provides raw access to analytics query data.
// VOLATILE: This API is subject to change at any time.
type AnalyticsResultRaw struct {
	reader analyticsRowReader
}

// NextBytes returns the next row as bytes.
func (arr *AnalyticsResultRaw) NextBytes() []byte {
	return arr.reader.NextRow()
}

// Err returns any errors that have occurred on the stream
func (arr *AnalyticsResultRaw) Err() error {
	return arr.reader.Err()
}

// Close marks the results as closed, returning any errors that occurred during reading the results.
func (arr *AnalyticsResultRaw) Close() error {
	return arr.reader.Close()
}

// MetaData returns any meta-data that was available from this query as bytes.
func (arr *AnalyticsResultRaw) MetaData() ([]byte, error) {
	return arr.reader.MetaData()
}

// AnalyticsResult allows access to the results of a query.
type AnalyticsResult struct {
	reader analyticsRowReader

	rowBytes []byte
}

func newAnalyticsResult(reader analyticsRowReader) *AnalyticsResult {
	return &AnalyticsResult{
		reader: reader,
	}
}

type analyticsRowReader interface {
	NextRow() []byte
	Err() error
	MetaData() ([]byte, error)
	Close() error
}

// Raw returns a AnalyticsResult which can be used to access the raw byte data from search queries.
// Calling this function invalidates the underlying AnalyticsResult which will no longer be able to be used.
// VOLATILE: This API is subject to change at any time.
func (r *AnalyticsResult) Raw() *AnalyticsResultRaw {
	vr := &AnalyticsResultRaw{
		reader: r.reader,
	}

	r.reader = nil
	return vr
}

// Next assigns the next result from the results into the value pointer, returning whether the read was successful.
func (r *AnalyticsResult) Next() bool {
	if r.reader == nil {
		return false
	}

	rowBytes := r.reader.NextRow()
	if rowBytes == nil {
		return false
	}

	r.rowBytes = rowBytes
	return true
}

// Row returns the value of the current row
func (r *AnalyticsResult) Row(valuePtr interface{}) error {
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
func (r *AnalyticsResult) Err() error {
	if r.reader == nil {
		return errors.New("result object is no longer valid")
	}

	return r.reader.Err()
}

// Close marks the results as closed, returning any errors that occurred during reading the results.
func (r *AnalyticsResult) Close() error {
	if r.reader == nil {
		return r.Err()
	}

	return r.reader.Close()
}

// One assigns the first value from the results into the value pointer.
// It will close the results but not before iterating through all remaining
// results, as such this should only be used for very small resultsets - ideally
// of, at most, length 1.
func (r *AnalyticsResult) One(valuePtr interface{}) error {
	if r.reader == nil {
		return r.Err()
	}

	// Read the bytes from the first row
	valueBytes := r.reader.NextRow()
	if valueBytes == nil {
		return ErrNoResult
	}

	// Skip through the remaining rows
	for r.reader.NextRow() != nil {
		// do nothing with the row
	}

	return json.Unmarshal(valueBytes, valuePtr)
}

// MetaData returns any meta-data that was available from this query.  Note that
// the meta-data will only be available once the object has been closed (either
// implicitly or explicitly).
func (r *AnalyticsResult) MetaData() (*AnalyticsMetaData, error) {
	if r.reader == nil {
		return nil, r.Err()
	}

	metaDataBytes, err := r.reader.MetaData()
	if err != nil {
		return nil, err
	}

	var jsonResp jsonAnalyticsResponse
	err = json.Unmarshal(metaDataBytes, &jsonResp)
	if err != nil {
		return nil, err
	}

	var metaData AnalyticsMetaData
	err = metaData.fromData(jsonResp)
	if err != nil {
		return nil, err
	}

	return &metaData, nil
}

// AnalyticsQuery executes the analytics query statement on the server.
func (c *Cluster) AnalyticsQuery(statement string, opts *AnalyticsOptions) (*AnalyticsResult, error) {
	if opts == nil {
		opts = &AnalyticsOptions{}
	}

	span := c.tracer.StartSpan("Query", opts.parentSpan).
		SetTag("couchbase.service", "analytics")
	defer span.Finish()

	timeout := opts.Timeout
	if opts.Timeout == 0 {
		timeout = c.timeoutsConfig.AnalyticsTimeout
	}
	deadline := time.Now().Add(timeout)

	retryStrategy := c.retryStrategyWrapper
	if opts.RetryStrategy != nil {
		retryStrategy = newRetryStrategyWrapper(opts.RetryStrategy)
	}

	queryOpts, err := opts.toMap()
	if err != nil {
		return nil, AnalyticsError{
			InnerError:      wrapError(err, "failed to generate query options"),
			Statement:       statement,
			ClientContextID: opts.ClientContextID,
		}
	}

	var priorityInt int32
	if opts.Priority {
		priorityInt = -1
	}

	queryOpts["statement"] = statement

	provider, err := c.getAnalyticsProvider()
	if err != nil {
		return nil, AnalyticsError{
			InnerError:      wrapError(err, "failed to get query provider"),
			Statement:       statement,
			ClientContextID: maybeGetAnalyticsOption(queryOpts, "client_context_id"),
		}
	}

	return execAnalyticsQuery(span, queryOpts, priorityInt, deadline, retryStrategy, provider, c.tracer)
}

func maybeGetAnalyticsOption(options map[string]interface{}, name string) string {
	if value, ok := options[name].(string); ok {
		return value
	}
	return ""
}

func execAnalyticsQuery(
	span requestSpan,
	options map[string]interface{},
	priority int32,
	deadline time.Time,
	retryStrategy *retryStrategyWrapper,
	provider analyticsProvider,
	tracer requestTracer,
) (*AnalyticsResult, error) {
	eSpan := tracer.StartSpan("request_encoding", span.Context())
	reqBytes, err := json.Marshal(options)
	eSpan.Finish()
	if err != nil {
		return nil, AnalyticsError{
			InnerError:      wrapError(err, "failed to marshall query body"),
			Statement:       maybeGetAnalyticsOption(options, "statement"),
			ClientContextID: maybeGetAnalyticsOption(options, "client_context_id"),
		}
	}

	res, err := provider.AnalyticsQuery(gocbcore.AnalyticsQueryOptions{
		Payload:       reqBytes,
		Priority:      int(priority),
		RetryStrategy: retryStrategy,
		Deadline:      deadline,
		TraceContext:  span.Context(),
	})
	if err != nil {
		return nil, maybeEnhanceAnalyticsError(err)
	}

	return newAnalyticsResult(res), nil
}
