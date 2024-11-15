package gocb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (qpc *queryProviderCore) CreatePrimaryIndex(c *Collection, bucketName string, opts *CreatePrimaryQueryIndexOptions) error {
	return qpc.createIndex(c, bucketName, opts.CustomName, nil, &CreateQueryIndexOptions{
		IgnoreIfExists: opts.IgnoreIfExists,
		Deferred:       opts.Deferred,
		NumReplicas:    opts.NumReplicas,
		Timeout:        opts.Timeout,
		RetryStrategy:  opts.RetryStrategy,
		ParentSpan:     opts.ParentSpan,
		ScopeName:      opts.ScopeName,
		CollectionName: opts.CollectionName,
		Context:        opts.Context,
	})
}

func (qpc *queryProviderCore) CreateIndex(c *Collection, bucketName, indexName string, fields []string, opts *CreateQueryIndexOptions) error {
	return qpc.createIndex(c, bucketName, indexName, fields, opts)
}

func (qpc *queryProviderCore) createIndex(c *Collection, bucketName, indexName string, fields []string, opts *CreateQueryIndexOptions) error {
	start := time.Now()
	spanName := "manager_query_create_index"

	var qs string
	if len(fields) == 0 {
		spanName = "manager_query_create_primary_index"
		qs += "CREATE PRIMARY INDEX"
	} else {
		qs += "CREATE INDEX"
	}
	if indexName != "" {
		qs += " `" + indexName + "`"
	}
	qs += " ON " + qpc.makeKeyspace(c, bucketName, opts.ScopeName, opts.CollectionName)
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

	var with []string
	if opts.Deferred {
		with = append(with, `"defer_build":true`)
	}
	if opts.NumReplicas > 0 {
		with = append(with, `"num_replica":`+strconv.Itoa(opts.NumReplicas))
	}

	if len(with) > 0 {
		withStr := strings.Join(with, ",")
		qs += " WITH {" + withStr + "}"
	}

	defer qpc.meter.ValueRecord(meterValueServiceManagement, spanName, start)

	span := createSpan(qpc.tracer, opts.ParentSpan, spanName, "management")
	defer span.End()

	_, err := qpc.doQuery(c, qs, &QueryOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		Adhoc:         true,
		ParentSpan:    span,
		Context:       opts.Context,
	})
	if err == nil {
		return nil
	}

	if opts.IgnoreIfExists && errors.Is(err, ErrIndexExists) {
		return nil
	}

	return err
}

func (qpc *queryProviderCore) DropPrimaryIndex(c *Collection, bucketName string, opts *DropPrimaryQueryIndexOptions) error {
	return qpc.dropIndex(c, bucketName, opts.CustomName, &DropQueryIndexOptions{
		IgnoreIfNotExists: opts.IgnoreIfNotExists,
		Timeout:           opts.Timeout,
		RetryStrategy:     opts.RetryStrategy,
		ParentSpan:        opts.ParentSpan,
		ScopeName:         opts.ScopeName,
		CollectionName:    opts.CollectionName,
		Context:           opts.Context,
	})
}

func (qpc *queryProviderCore) DropIndex(c *Collection, bucketName, indexName string, opts *DropQueryIndexOptions) error {
	return qpc.dropIndex(c, bucketName, indexName, opts)
}

func (qpc *queryProviderCore) dropIndex(c *Collection, bucketName, indexName string, opts *DropQueryIndexOptions) error {
	start := time.Now()
	spanName := "manager_query_drop_index"
	var qs string

	keyspace := qpc.makeKeyspace(c, bucketName, opts.ScopeName, opts.CollectionName)
	if indexName == "" {
		spanName = "manager_query_drop_primary_index"
		qs += "DROP PRIMARY INDEX ON " + keyspace
	} else {
		if c == nil && opts.ScopeName == "" && opts.CollectionName == "" {
			qs += "DROP INDEX " + keyspace + ".`" + indexName + "`"
		} else {
			qs += "DROP INDEX `" + indexName + "` ON " + keyspace
		}
	}
	defer qpc.meter.ValueRecord(meterValueServiceManagement, spanName, start)

	span := createSpan(qpc.tracer, opts.ParentSpan, spanName, "management")
	defer span.End()

	_, err := qpc.doQuery(c, qs, &QueryOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		Adhoc:         true,
		ParentSpan:    span,
		Context:       opts.Context,
	})
	if err == nil {
		return nil
	}

	if opts.IgnoreIfNotExists && errors.Is(err, ErrIndexNotFound) {
		return nil
	}

	return err
}

func buildGetAllIndexesWhereClause(c *Collection, bucketName, scopeName, collectionName string) (string, map[string]interface{}) {
	if c == nil {
		bucketCond := "bucket_id = $bucketName"
		scopeCond := "(" + bucketCond + " AND scope_id = $scopeName)"
		collectionCond := "(" + scopeCond + " AND keyspace_id = $collectionName)"
		params := map[string]interface{}{
			"bucketName": bucketName,
		}

		var where string
		if collectionName != "" {
			where = collectionCond
			params["scopeName"] = scopeName
			params["collectionName"] = collectionName
		} else if scopeName != "" {
			where = scopeCond
			params["scopeName"] = scopeName
		} else {
			where = bucketCond
		}

		if collectionName == "_default" || collectionName == "" {
			defaultColCond := "(bucket_id IS MISSING AND keyspace_id = $bucketName)"
			where = "(" + where + " OR " + defaultColCond + ")"
		}

		return where, params
	}

	scope, collection := normaliseQueryCollectionKeyspace(c)
	var where string
	if scope == "_default" && collection == "_default" {
		where = "((bucket_id=$bucketName AND scope_id=$scopeName AND keyspace_id=$collectionName) OR (bucket_id IS MISSING and keyspace_id=$bucketName)) "
	} else {
		where = "(bucket_id=$bucketName AND scope_id=$scopeName AND keyspace_id=$collectionName)"
	}

	return where, map[string]interface{}{
		"bucketName":     c.bucketName(),
		"scopeName":      scope,
		"collectionName": collection,
	}

}

func (qpc *queryProviderCore) GetAllIndexes(c *Collection, bucketName string, opts *GetAllQueryIndexesOptions) ([]QueryIndex, error) {
	start := time.Now()
	defer qpc.meter.ValueRecord(meterValueServiceManagement, "manager_query_get_all_indexes", start)

	return qpc.getAllIndexes(c, bucketName, opts)
}

func (qpc *queryProviderCore) getAllIndexes(c *Collection, bucketName string, opts *GetAllQueryIndexesOptions) ([]QueryIndex, error) {
	whereClause, params := buildGetAllIndexesWhereClause(c, bucketName, opts.ScopeName, opts.CollectionName)

	span := createSpan(qpc.tracer, opts.ParentSpan, "manager_query_get_all_indexes", "management")
	defer span.End()

	q := "SELECT `idx`.* FROM system:indexes AS idx WHERE " + whereClause + " AND `using` = \"gsi\" " +
		"ORDER BY is_primary DESC, name ASC"

	rows, err := qpc.doQuery(c, q, &QueryOptions{
		NamedParameters: params,
		Readonly:        true,
		Timeout:         opts.Timeout,
		RetryStrategy:   opts.RetryStrategy,
		Adhoc:           true,
		ParentSpan:      span,
		Context:         opts.Context,
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

func (qpc *queryProviderCore) BuildDeferredIndexes(c *Collection, bucketName string, opts *BuildDeferredQueryIndexOptions) ([]string, error) {
	start := time.Now()
	defer qpc.meter.ValueRecord(meterValueServiceManagement, "manager_query_build_deferred_indexes", start)

	span := createSpan(qpc.tracer, opts.ParentSpan, "manager_query_build_deferred_indexes", "management")
	defer span.End()

	var whereClause string
	params := make(map[string]interface{})
	if c == nil {
		if opts.CollectionName == "" {
			whereClause = "(keyspace_id = $bucketName AND bucket_id IS MISSING)"
			params["bucketName"] = bucketName
		} else {
			whereClause = "bucket_id = $bucketName AND scope_id = $scopeName AND keyspace_id = $collectionName"
			params["bucketName"] = bucketName
			params["scopeName"] = opts.ScopeName
			params["collectionName"] = opts.CollectionName
		}
	} else {
		scope, collection := normaliseQueryCollectionKeyspace(c)
		if scope == "_default" && collection == "_default" {
			whereClause = "((bucket_id=$bucketName AND scope_id=$scopeName AND keyspace_id=$collectionName) OR (bucket_id IS MISSING and keyspace_id=$bucketName)) "
		} else {
			whereClause = "(bucket_id=$bucketName AND scope_id=$scopeName AND keyspace_id=$collectionName)"
		}
		params = map[string]interface{}{
			"bucketName":     c.bucketName(),
			"scopeName":      scope,
			"collectionName": collection,
		}
	}

	query := "SELECT RAW name from system:indexes WHERE " + whereClause + " AND state = \"deferred\""

	indexesRes, err := qpc.doQuery(c, query, &QueryOptions{
		Timeout:         opts.Timeout,
		RetryStrategy:   opts.RetryStrategy,
		Adhoc:           true,
		ParentSpan:      span,
		Context:         opts.Context,
		NamedParameters: params,
	})
	if err != nil {
		return nil, err
	}

	var deferredList []string
	for _, row := range indexesRes {
		var name string
		err := json.Unmarshal(row, &name)
		if err != nil {
			return nil, err
		}

		deferredList = append(deferredList, name)
	}

	if len(deferredList) == 0 {
		// Don't try to build an empty index list
		return nil, nil
	}

	keyspace := qpc.makeKeyspace(c, bucketName, opts.ScopeName, opts.CollectionName)
	var qs string
	qs += "BUILD INDEX ON " + keyspace + "("
	for i := 0; i < len(deferredList); i++ {
		if i > 0 {
			qs += ", "
		}
		qs += "`" + deferredList[i] + "`"
	}
	qs += ")"

	_, err = qpc.doQuery(c, qs, &QueryOptions{
		Timeout:       opts.Timeout,
		RetryStrategy: opts.RetryStrategy,
		Adhoc:         true,
		ParentSpan:    span,
		Context:       opts.Context,
	})
	if err != nil {
		return nil, err
	}

	return deferredList, nil
}

func checkIndexesActiveCore(indexes []QueryIndex, checkList []string) (bool, error) {
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
		if checkIndexes[i].State != string(queryIndexStateOnline) {
			logDebugf("Index not online: %s is in state %s", checkIndexes[i].Name, checkIndexes[i].State)
			return false, nil
		}
	}
	return true, nil
}

func (qpc *queryProviderCore) WatchIndexes(c *Collection, bucketName string, watchList []string, timeout time.Duration, opts *WatchQueryIndexOptions,
) error {
	start := time.Now()
	defer qpc.meter.ValueRecord(meterValueServiceManagement, "manager_query_watch_indexes", start)

	span := createSpan(qpc.tracer, opts.ParentSpan, "manager_query_watch_indexes", "management")
	defer span.End()

	if opts.WatchPrimary {
		watchList = append(watchList, "#primary")
	}

	deadline := time.Now().Add(timeout)

	curInterval := 50 * time.Millisecond
	for {
		if deadline.Before(time.Now()) {
			return ErrUnambiguousTimeout
		}

		indexes, err := qpc.getAllIndexes(
			c,
			bucketName,
			&GetAllQueryIndexesOptions{
				Timeout:        time.Until(deadline),
				RetryStrategy:  opts.RetryStrategy,
				ParentSpan:     span,
				ScopeName:      opts.ScopeName,
				CollectionName: opts.CollectionName,
				Context:        opts.Context,
			})
		if err != nil {
			return err
		}

		allOnline, err := checkIndexesActiveCore(indexes, watchList)
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

func (qpc *queryProviderCore) doQuery(c *Collection, q string, opts *QueryOptions) ([][]byte, error) {
	if opts.Timeout == 0 {
		opts.Timeout = qpc.timeouts.ManagementTimeout
	}

	var scope *Scope
	if c != nil {
		// If we have a collection then we need to do a scope level query.
		scope = c.bucket.Scope(c.scope)
	}

	result, err := qpc.Query(q, scope, opts)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return rows, nil
}

type queryIndexState string

const (
	queryIndexStateDeferred  queryIndexState = "deferred"
	queryIndexStateBuilding  queryIndexState = "building"
	queryIndexStatePending   queryIndexState = "pending"
	queryIndexStateOnline    queryIndexState = "online"
	queryIndexStateOffline   queryIndexState = "offline"
	queryIndexStateAbridged  queryIndexState = "abridged"
	queryIndexStateScheduled queryIndexState = "scheduled for creation"
)

type jsonQueryIndex struct {
	Name      string          `json:"name"`
	IsPrimary bool            `json:"is_primary"`
	Type      QueryIndexType  `json:"using"`
	State     queryIndexState `json:"state"`
	Keyspace  string          `json:"keyspace_id"`
	Namespace string          `json:"namespace_id"`
	IndexKey  []string        `json:"index_key"`
	Condition string          `json:"condition"`
	Partition string          `json:"partition"`
	Scope     string          `json:"scope_id"`
	Bucket    string          `json:"bucket_id"`
}

// QueryIndex represents a Couchbase GSI index.
type QueryIndex struct {
	Name           string
	IsPrimary      bool
	Type           QueryIndexType
	State          string
	Keyspace       string
	Namespace      string
	IndexKey       []string
	Condition      string
	Partition      string
	CollectionName string
	ScopeName      string
	BucketName     string
}

func (index *QueryIndex) fromData(data jsonQueryIndex) error {
	index.Name = data.Name
	index.IsPrimary = data.IsPrimary
	index.Type = data.Type
	index.State = string(data.State)
	index.Keyspace = data.Keyspace
	index.Namespace = data.Namespace
	index.IndexKey = data.IndexKey
	index.Condition = data.Condition
	index.Partition = data.Partition
	index.ScopeName = data.Scope
	if data.Bucket == "" {
		index.BucketName = data.Keyspace
	} else {
		index.BucketName = data.Bucket
	}
	if data.Scope != "" {
		index.CollectionName = data.Keyspace
	}

	return nil
}

func normaliseQueryCollectionKeyspace(c *Collection) (string, string) {
	// Ensure scope and collection names are populated, if the DefaultX functions on bucket are
	// used then the names will be empty by default.
	scope := c.scope
	if scope == "" {
		scope = "_default"
	}
	collection := c.collectionName
	if collection == "" {
		collection = "_default"
	}

	return scope, collection
}

func (qpc *queryProviderCore) makeKeyspace(c *Collection, bucketName, scopeName, collectionName string) string {
	if c != nil {
		// If we have a collection then we need to build the namespace using it rather than options.
		scope, collection := normaliseQueryCollectionKeyspace(c)

		return fmt.Sprintf("`%s`.`%s`.`%s`", c.bucketName(), scope, collection)
	}

	if scopeName != "" && collectionName != "" {
		return fmt.Sprintf("`%s`.`%s`.`%s`", bucketName, scopeName, collectionName)
	} else if collectionName == "" && scopeName != "" {
		return fmt.Sprintf("`%s`.`%s`.`_default", bucketName, scopeName)
	} else if collectionName != "" && scopeName == "" {
		return fmt.Sprintf("`%s`.`_default`.`%s", bucketName, collectionName)
	}
	return "`" + bucketName + "`"
}
