package gocb

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
)

// QueryIndexManager provides methods for performing Couchbase query index management.
type QueryIndexManager struct {
	provider queryIndexQueryProvider

	globalTimeout time.Duration
	tracer        requestTracer
}

type queryIndexQueryProvider interface {
	Query(statement string, opts *QueryOptions) (*QueryResult, error)
}

func (qm *QueryIndexManager) tryParseErrorMessage(err error) error {
	var qErr *QueryError
	if !errors.As(err, &qErr) {
		return err
	}

	if len(qErr.Errors) == 0 {
		return err
	}

	firstErr := qErr.Errors[0]
	var innerErr error
	// The server doesn't return meaningful error codes when it comes to index management so we need to go spelunking.
	msg := strings.ToLower(firstErr.Message)
	if match, err := regexp.MatchString(".*?ndex .*? not found.*", msg); err == nil && match {
		innerErr = ErrIndexNotFound
	} else if match, err := regexp.MatchString(".*?ndex .*? already exists.*", msg); err == nil && match {
		innerErr = ErrIndexExists
	}

	if innerErr == nil {
		return err
	}

	return QueryError{
		InnerError:      innerErr,
		Statement:       qErr.Statement,
		ClientContextID: qErr.ClientContextID,
		Errors:          qErr.Errors,
		Endpoint:        qErr.Endpoint,
		RetryReasons:    qErr.RetryReasons,
		RetryAttempts:   qErr.RetryAttempts,
	}
}

func (qm *QueryIndexManager) doQuery(q string, opts *QueryOptions) ([][]byte, error) {
	if opts.Timeout == 0 {
		opts.Timeout = qm.globalTimeout
	}

	result, err := qm.provider.Query(q, opts)
	if err != nil {
		return nil, qm.tryParseErrorMessage(err)
	}

	var rows [][]byte
	for result.Next() {
		var row json.RawMessage
		err := result.Row(&row)
		if err != nil {
			logWarnf("management operation failed to read row: %s", err)
		} else {
			rows = append(rows, row)
		}
	}
	err = result.Err()
	if err != nil {
		return nil, qm.tryParseErrorMessage(err)
	}

	return rows, nil
}

type jsonQueryIndex struct {
	Name      string         `json:"name"`
	IsPrimary bool           `json:"is_primary"`
	Type      QueryIndexType `json:"using"`
	State     string         `json:"state"`
	Keyspace  string         `json:"keyspace_id"`
	Namespace string         `json:"namespace_id"`
	IndexKey  []string       `json:"index_key"`
	Condition string         `json:"condition"`
}

// QueryIndex represents a Couchbase GSI index.
type QueryIndex struct {
	Name      string
	IsPrimary bool
	Type      QueryIndexType
	State     string
	Keyspace  string
	Namespace string
	IndexKey  []string
	Condition string
}

func (index *QueryIndex) fromData(data jsonQueryIndex) error {
	index.Name = data.Name
	index.IsPrimary = data.IsPrimary
	index.Type = data.Type
	index.State = data.State
	index.Keyspace = data.Keyspace
	index.Namespace = data.Namespace
	index.IndexKey = data.IndexKey
	index.Condition = data.Condition

	return nil
}

type createQueryIndexOptions struct {
	IgnoreIfExists bool
	Deferred       bool

	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

func (qm *QueryIndexManager) createIndex(
	tracectx requestSpanContext,
	bucketName, indexName string,
	fields []string,
	opts createQueryIndexOptions,
) error {
	var qs string

	if len(fields) == 0 {
		qs += "CREATE PRIMARY INDEX"
	} else {
		qs += "CREATE INDEX"
	}
	if indexName != "" {
		qs += " `" + indexName + "`"
	}
	qs += " ON `" + bucketName + "`"
	if len(fields) > 0 {
		qs += " ("
		for i := 0; i < len(fields); i++ {
			if i > 0 {
				qs += ", "
			}
			qs += "`" + fields[i] + "`"
		}
		qs += ")"
	}
	if opts.Deferred {
		qs += " WITH {\"defer_build\": true}"
	}

	_, err := qm.doQuery(qs, &QueryOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		Adhoc:         true,
		parentSpan:    tracectx,
	})
	if err == nil {
		return nil
	}

	if opts.IgnoreIfExists && errors.Is(err, ErrIndexExists) {
		return nil
	}

	return err
}

// CreateQueryIndexOptions is the set of options available to the query indexes CreateIndex operation.
type CreateQueryIndexOptions struct {
	IgnoreIfExists bool
	Deferred       bool

	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// CreateIndex creates an index over the specified fields.
func (qm *QueryIndexManager) CreateIndex(bucketName, indexName string, fields []string, opts *CreateQueryIndexOptions) error {
	if opts == nil {
		opts = &CreateQueryIndexOptions{}
	}

	if indexName == "" {
		return invalidArgumentsError{
			message: "an invalid index name was specified",
		}
	}
	if len(fields) <= 0 {
		return invalidArgumentsError{
			message: "you must specify at least one field to index",
		}
	}

	span := qm.tracer.StartSpan("CreateIndex", nil).
		SetTag("couchbase.service", "query")
	defer span.Finish()

	return qm.createIndex(span.Context(), bucketName, indexName, fields, createQueryIndexOptions{
		IgnoreIfExists: opts.IgnoreIfExists,
		Deferred:       opts.Deferred,
		Timeout:        opts.Timeout,
		RetryStrategy:  opts.RetryStrategy,
	})
}

// CreatePrimaryQueryIndexOptions is the set of options available to the query indexes CreatePrimaryIndex operation.
type CreatePrimaryQueryIndexOptions struct {
	IgnoreIfExists bool
	Deferred       bool
	CustomName     string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// CreatePrimaryIndex creates a primary index.  An empty customName uses the default naming.
func (qm *QueryIndexManager) CreatePrimaryIndex(bucketName string, opts *CreatePrimaryQueryIndexOptions) error {
	if opts == nil {
		opts = &CreatePrimaryQueryIndexOptions{}
	}

	span := qm.tracer.StartSpan("CreatePrimaryIndex", nil).
		SetTag("couchbase.service", "query")
	defer span.Finish()

	return qm.createIndex(
		span.Context(),
		bucketName,
		opts.CustomName,
		nil,
		createQueryIndexOptions{
			IgnoreIfExists: opts.IgnoreIfExists,
			Deferred:       opts.Deferred,
			Timeout:        opts.Timeout,
			RetryStrategy:  opts.RetryStrategy,
		})
}

type dropQueryIndexOptions struct {
	IgnoreIfNotExists bool

	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

func (qm *QueryIndexManager) dropIndex(
	tracectx requestSpanContext,
	bucketName, indexName string,
	opts dropQueryIndexOptions,
) error {
	var qs string

	if indexName == "" {
		qs += "DROP PRIMARY INDEX ON `" + bucketName + "`"
	} else {
		qs += "DROP INDEX `" + bucketName + "`.`" + indexName + "`"
	}

	_, err := qm.doQuery(qs, &QueryOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		Adhoc:         true,
		parentSpan:    tracectx,
	})
	if err == nil {
		return nil
	}

	if opts.IgnoreIfNotExists && errors.Is(err, ErrIndexNotFound) {
		return nil
	}

	return err
}

// DropQueryIndexOptions is the set of options available to the query indexes DropIndex operation.
type DropQueryIndexOptions struct {
	IgnoreIfNotExists bool

	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// DropIndex drops a specific index by name.
func (qm *QueryIndexManager) DropIndex(bucketName, indexName string, opts *DropQueryIndexOptions) error {
	if opts == nil {
		opts = &DropQueryIndexOptions{}
	}

	if indexName == "" {
		return invalidArgumentsError{
			message: "an invalid index name was specified",
		}
	}

	span := qm.tracer.StartSpan("DropIndex", nil).
		SetTag("couchbase.service", "query")
	defer span.Finish()

	return qm.dropIndex(
		span.Context(),
		bucketName,
		indexName,
		dropQueryIndexOptions{
			IgnoreIfNotExists: opts.IgnoreIfNotExists,
			Timeout:           opts.Timeout,
			RetryStrategy:     opts.RetryStrategy,
		})
}

// DropPrimaryQueryIndexOptions is the set of options available to the query indexes DropPrimaryIndex operation.
type DropPrimaryQueryIndexOptions struct {
	IgnoreIfNotExists bool
	CustomName        string

	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// DropPrimaryIndex drops the primary index.  Pass an empty customName for unnamed primary indexes.
func (qm *QueryIndexManager) DropPrimaryIndex(bucketName string, opts *DropPrimaryQueryIndexOptions) error {
	if opts == nil {
		opts = &DropPrimaryQueryIndexOptions{}
	}

	span := qm.tracer.StartSpan("DropPrimaryIndex", nil).
		SetTag("couchbase.service", "query")
	defer span.Finish()

	return qm.dropIndex(
		span.Context(),
		bucketName,
		opts.CustomName,
		dropQueryIndexOptions{
			IgnoreIfNotExists: opts.IgnoreIfNotExists,
			Timeout:           opts.Timeout,
			RetryStrategy:     opts.RetryStrategy,
		})
}

// GetAllQueryIndexesOptions is the set of options available to the query indexes GetAllIndexes operation.
type GetAllQueryIndexesOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// GetAllIndexes returns a list of all currently registered indexes.
func (qm *QueryIndexManager) GetAllIndexes(bucketName string, opts *GetAllQueryIndexesOptions) ([]QueryIndex, error) {
	if opts == nil {
		opts = &GetAllQueryIndexesOptions{}
	}

	span := qm.tracer.StartSpan("GetAllIndexes", nil).
		SetTag("couchbase.service", "query")
	defer span.Finish()

	return qm.getAllIndexes(span.Context(), bucketName, opts)
}

func (qm *QueryIndexManager) getAllIndexes(
	tracectx requestSpanContext,
	bucketName string,
	opts *GetAllQueryIndexesOptions,
) ([]QueryIndex, error) {
	q := "SELECT `indexes`.* FROM system:indexes WHERE keyspace_id=? AND `using`=\"gsi\""
	rows, err := qm.doQuery(q, &QueryOptions{
		PositionalParameters: []interface{}{bucketName},
		Readonly:             true,
		Timeout:              opts.Timeout,
		RetryStrategy:        opts.RetryStrategy,
		Adhoc:                true,
		parentSpan:           tracectx,
	})
	if err != nil {
		return nil, err
	}

	var indexes []QueryIndex
	for _, row := range rows {
		var jsonIdx jsonQueryIndex
		err := json.Unmarshal(row, &jsonIdx)
		if err != nil {
			return nil, err
		}

		var index QueryIndex
		err = index.fromData(jsonIdx)
		if err != nil {
			return nil, err
		}

		indexes = append(indexes, index)
	}

	return indexes, nil
}

// BuildDeferredQueryIndexOptions is the set of options available to the query indexes BuildDeferredIndexes operation.
type BuildDeferredQueryIndexOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// BuildDeferredIndexes builds all indexes which are currently in deferred state.
func (qm *QueryIndexManager) BuildDeferredIndexes(bucketName string, opts *BuildDeferredQueryIndexOptions) ([]string, error) {
	if opts == nil {
		opts = &BuildDeferredQueryIndexOptions{}
	}

	span := qm.tracer.StartSpan("BuildDeferredIndexes", nil).
		SetTag("couchbase.service", "query")
	defer span.Finish()

	indexList, err := qm.getAllIndexes(
		span.Context(),
		bucketName,
		&GetAllQueryIndexesOptions{
			Timeout:       opts.Timeout,
			RetryStrategy: opts.RetryStrategy,
		})
	if err != nil {
		return nil, err
	}

	var deferredList []string
	for i := 0; i < len(indexList); i++ {
		var index = indexList[i]
		if index.State == "deferred" || index.State == "pending" {
			deferredList = append(deferredList, index.Name)
		}
	}

	if len(deferredList) == 0 {
		// Don't try to build an empty index list
		return nil, nil
	}

	var qs string
	qs += "BUILD INDEX ON `" + bucketName + "`("
	for i := 0; i < len(deferredList); i++ {
		if i > 0 {
			qs += ", "
		}
		qs += "`" + deferredList[i] + "`"
	}
	qs += ")"

	_, err = qm.doQuery(qs, &QueryOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		Adhoc:         true,
		parentSpan:    span,
	})
	if err != nil {
		return nil, err
	}

	return deferredList, nil
}

func checkIndexesActive(indexes []QueryIndex, checkList []string) (bool, error) {
	var checkIndexes []QueryIndex
	for i := 0; i < len(checkList); i++ {
		indexName := checkList[i]

		for j := 0; j < len(indexes); j++ {
			if indexes[j].Name == indexName {
				checkIndexes = append(checkIndexes, indexes[j])
				break
			}
		}
	}

	if len(checkIndexes) != len(checkList) {
		return false, ErrIndexNotFound
	}

	for i := 0; i < len(checkIndexes); i++ {
		if checkIndexes[i].State != "online" {
			return false, nil
		}
	}
	return true, nil
}

// WatchQueryIndexOptions is the set of options available to the query indexes Watch operation.
type WatchQueryIndexOptions struct {
	WatchPrimary bool

	RetryStrategy RetryStrategy
}

// WatchIndexes waits for a set of indexes to come online.
func (qm *QueryIndexManager) WatchIndexes(bucketName string, watchList []string, timeout time.Duration, opts *WatchQueryIndexOptions) error {
	if opts == nil {
		opts = &WatchQueryIndexOptions{}
	}

	span := qm.tracer.StartSpan("WatchIndexes", nil).
		SetTag("couchbase.service", "query")
	defer span.Finish()

	if opts.WatchPrimary {
		watchList = append(watchList, "#primary")
	}

	deadline := time.Now().Add(timeout)

	curInterval := 50 * time.Millisecond
	for {
		if deadline.Before(time.Now()) {
			return ErrUnambiguousTimeout
		}

		indexes, err := qm.getAllIndexes(
			span.Context(),
			bucketName,
			&GetAllQueryIndexesOptions{
				Timeout:       time.Until(deadline),
				RetryStrategy: opts.RetryStrategy,
			})
		if err != nil {
			return err
		}

		allOnline, err := checkIndexesActive(indexes, watchList)
		if err != nil {
			return err
		}

		if allOnline {
			break
		}

		curInterval += 500 * time.Millisecond
		if curInterval > 1000 {
			curInterval = 1000
		}

		// Make sure we don't sleep past our overall deadline, if we adjust the
		// deadline then it will be caught at the top of this loop as a timeout.
		sleepDeadline := time.Now().Add(curInterval)
		if sleepDeadline.After(deadline) {
			sleepDeadline = deadline
		}

		// wait till our next poll interval
		time.Sleep(time.Until(sleepDeadline))
	}

	return nil
}
