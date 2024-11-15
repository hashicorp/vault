package gocbcore

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// N1QLRowReader providers access to the rows of a n1ql query
type N1QLRowReader struct {
	streamer   *queryStreamer
	endpoint   string
	statement  string
	statusCode int
}

// NextRow reads the next rows bytes from the stream
func (q *N1QLRowReader) NextRow() []byte {
	return q.streamer.NextRow()
}

// Err returns any errors that occurred during streaming.
func (q N1QLRowReader) Err() error {
	err := q.streamer.Err()
	if err != nil {
		return err
	}

	meta, metaErr := q.streamer.MetaData()
	if metaErr != nil {
		return metaErr
	}

	raw, descs, err := parseN1QLError(meta)
	if err != nil {
		return &N1QLError{
			InnerError:       err,
			Errors:           descs,
			ErrorText:        raw,
			Statement:        q.statement,
			HTTPResponseCode: q.statusCode,
		}
	}
	if len(descs) > 0 {
		return &N1QLError{
			InnerError:       errors.New("query error"),
			Errors:           descs,
			ErrorText:        raw,
			Statement:        q.statement,
			HTTPResponseCode: q.statusCode,
		}
	}

	return nil
}

// MetaData fetches the non-row bytes streamed in the response.
func (q *N1QLRowReader) MetaData() ([]byte, error) {
	return q.streamer.MetaData()
}

// Close immediately shuts down the connection
func (q *N1QLRowReader) Close() error {
	return q.streamer.Close()
}

// PreparedName returns the name of the prepared statement created when using enhanced prepared statements.
// If the prepared name has not been seen on the stream then this will return an error.
// Volatile: This API is subject to change.
func (q N1QLRowReader) PreparedName() (string, error) {
	val := q.streamer.EarlyMetadata("prepared")
	if val == nil {
		return "", wrapN1QLError(nil, "", errors.New("prepared name not found in metadata"), "", 0)
	}

	var name string
	err := json.Unmarshal(val, &name)
	if err != nil {
		return "", wrapN1QLError(nil, "", errors.New("failed to parse prepared name"), "", 0)
	}

	return name, nil
}

// Endpoint returns the address that this query was run against.
// Internal: This should never be used and is not supported.
func (q *N1QLRowReader) Endpoint() string {
	return q.endpoint
}

// N1QLQueryOptions represents the various options available for a n1ql query.
type N1QLQueryOptions struct {
	Payload       []byte
	RetryStrategy RetryStrategy
	Deadline      time.Time

	// Internal: This should never be used and is not supported.
	User string
	// Internal: This should never be used and is not supported.
	Endpoint string

	TraceContext RequestSpanContext
}

func wrapN1QLError(req *httpRequest, statement string, err error, errBody string, statusCode int) *N1QLError {
	if err == nil {
		err = errors.New("query error")
	}

	ierr := &N1QLError{
		InnerError: err,
	}

	if req != nil {
		ierr.Endpoint = req.Endpoint
		ierr.ClientContextID = req.UniqueID
		ierr.RetryAttempts = req.RetryAttempts()
		ierr.RetryReasons = req.RetryReasons()
	}

	ierr.ErrorText = errBody
	ierr.Statement = statement
	ierr.HTTPResponseCode = statusCode

	return ierr
}

type jsonN1QLError struct {
	Code   uint32                 `json:"code"`
	Msg    string                 `json:"msg"`
	Reason map[string]interface{} `json:"reason"`
	Retry  bool                   `json:"retry"`
}

type jsonN1QLErrorResponse struct {
	Errors json.RawMessage
}

func extractN1QL12009Error(err N1QLErrorDesc) error {
	if len(err.Reason) > 0 {
		if code, ok := err.Reason["code"]; ok {
			// sad panda
			code = int(code.(float64))
			if code == 12033 {
				return errCasMismatch
			} else if code == 17014 {
				return errDocumentNotFound
			} else if code == 17012 {
				return errDocumentExists
			}
		}

		return errDMLFailure
	}

	if strings.Contains(strings.ToLower(err.Message), "cas mismatch") {
		return errCasMismatch
	}
	return errDMLFailure

}

func parseN1QLErrorResp(req *httpRequest, statement string, resp *HTTPResponse) *N1QLError {
	var errorDescs []N1QLErrorDesc
	var err error
	var raw string
	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr == nil {
		raw, errorDescs, err = parseN1QLError(respBody)
	}
	errOut := wrapN1QLError(req, statement, err, raw, resp.StatusCode)
	errOut.Errors = errorDescs
	return errOut
}

func parseN1QLError(respBody []byte) (string, []N1QLErrorDesc, error) {
	var err error
	var errorDescs []N1QLErrorDesc

	var rawRespParse jsonN1QLErrorResponse
	parseErr := json.Unmarshal(respBody, &rawRespParse)
	if parseErr != nil {
		return "", nil, nil
	}

	var respParse []jsonN1QLError
	parseErr = json.Unmarshal(rawRespParse.Errors, &respParse)
	if parseErr == nil {
		for _, jsonErr := range respParse {
			errorDescs = append(errorDescs, N1QLErrorDesc{
				Code:    jsonErr.Code,
				Message: jsonErr.Msg,
				Reason:  jsonErr.Reason,
				Retry:   jsonErr.Retry,
			})
		}
	}

	if len(errorDescs) >= 1 {
		firstErr := errorDescs[0]
		errCode := firstErr.Code
		errCodeGroup := errCode / 1000
		msgLower := strings.ToLower(firstErr.Message)

		switch errCodeGroup {
		case 1:
			switch errCode {
			case 1065:
				if strings.Contains(msgLower, "query_context") {
					err = wrapError(errFeatureNotAvailable, "this server requires that a query context be used for queries")
				} else if strings.Contains(msgLower, "preserve_expiry") {
					err = wrapError(errFeatureNotAvailable, "this server does not support preserve expiry")
				} else if strings.Contains(msgLower, "use_replica") {
					err = wrapError(errFeatureNotAvailable, "this server does not support use replica")
				}
			case 1080:
				// This can happen when the server starts streaming responses - at this point our timeout is already
				// canceled. But then the streaming takes longer than the configured timeout, in which case the query
				// engine will proactively send us a timeout and we need to convert it.
				err = errUnambiguousTimeout
			case 1191:
				err = errRateLimitedFailure
			case 1192:
				err = errRateLimitedFailure
			case 1193:
				err = errRateLimitedFailure
			case 1194:
				err = errRateLimitedFailure
			case 1197:
				err = wrapError(errFeatureNotAvailable, "this server requires that a query context be used for queries")
			}
		case 3:
			switch errCode {
			case 3000:
				err = errParsingFailure
			case 3230:
				if strings.Contains(msgLower, "advisor") || strings.Contains(msgLower, "advise") {
					err = wrapError(errFeatureNotAvailable, "query index advisor is not supported on community edition")
				} else if strings.Contains(msgLower, "query window functions") {
					err = wrapError(errFeatureNotAvailable, "query window functions are not supported on community edition")
				}
			}
		case 4:
			switch errCode {
			case 4040:
				err = errPreparedStatementFailure
			case 4050:
				err = errPreparedStatementFailure
			case 4060:
				err = errPreparedStatementFailure
			case 4070:
				err = errPreparedStatementFailure
			case 4080:
				err = errPreparedStatementFailure
			case 4090:
				err = errPreparedStatementFailure
			case 4300:
				err = errIndexExists
			default:
				err = errPlanningFailure
			}
		case 5:
			switch errCode {
			case 5000:
				if match, matchErr := regexp.MatchString(".*?ndex .*? not found.*", msgLower); matchErr == nil && match {
					err = errIndexNotFound
				} else if match, matchErr := regexp.MatchString(".*?ndex does not exist.*", msgLower); matchErr == nil && match {
					err = errIndexNotFound
				} else if match, matchErr := regexp.MatchString(".*?ndex .*? already exist.*", msgLower); matchErr == nil && match {
					err = errIndexExists
				} else if strings.Contains(msgLower,
					"limit for number of indexes that can be created per scope has been reached") {
					err = errQuotaLimitedFailure
				} else {
					err = errInternalServerFailure
				}
			default:
				err = errInternalServerFailure
			}
		case 10:
			err = errAuthenticationFailure
		case 12:
			switch errCode {
			case 12004:
				err = errIndexNotFound
			case 12016:
				err = errIndexNotFound
			case 12009:
				err = extractN1QL12009Error(firstErr)
			default:
				err = errIndexFailure
			}
		case 13:
			switch errCode {
			case 13014:
				err = errAuthenticationFailure
			}
		case 14:
			err = errIndexFailure
		}
	}
	var rawErrors string
	if err == nil && len(rawRespParse.Errors) > 0 {
		// Only populate if this is an error that we don't recognise.
		rawErrors = string(rawRespParse.Errors)
	}

	return rawErrors, errorDescs, err
}

const (
	// nolint: unused,deadcode,varcheck
	useReplicaSupportLevelUnknown = uint32(iota)
	useReplicaSupportLevelUnsupported
	useReplicaSupportLevelSupported
)

type n1qlQueryComponent struct {
	httpComponent httpComponentInterface
	cfgMgr        configManager
	tracer        *tracerComponent

	queryCache *n1qlQueryCache

	enhancedPreparedSupported uint32
	useReplicaSupported       uint32
}

type n1qlQueryCache struct {
	cache     map[n1qlQueryCacheStatementContext]*n1qlQueryCacheEntry
	cacheLock sync.RWMutex
}

type n1qlQueryCacheStatementContext struct {
	Statement string
	Context   string
}

func newN1qlQueryCache() *n1qlQueryCache {
	return &n1qlQueryCache{
		cache: make(map[n1qlQueryCacheStatementContext]*n1qlQueryCacheEntry),
	}
}

func (cache *n1qlQueryCache) Invalidate() {
	cache.cacheLock.Lock()
	cache.cache = make(map[n1qlQueryCacheStatementContext]*n1qlQueryCacheEntry)
	cache.cacheLock.Unlock()
}

func (cache *n1qlQueryCache) Put(statement n1qlQueryCacheStatementContext, entry *n1qlQueryCacheEntry) {
	cache.cacheLock.Lock()
	cache.cache[statement] = entry
	cache.cacheLock.Unlock()
}

func (cache *n1qlQueryCache) Delete(statement n1qlQueryCacheStatementContext) {
	cache.cacheLock.Lock()
	delete(cache.cache, statement)
	cache.cacheLock.Unlock()
}

func (cache *n1qlQueryCache) Get(statement n1qlQueryCacheStatementContext) *n1qlQueryCacheEntry {
	cache.cacheLock.RLock()
	entry := cache.cache[statement]
	if entry == nil {
		cache.cacheLock.RUnlock()
		return nil
	}
	cached := *entry
	cache.cacheLock.RUnlock()

	return &cached
}

type n1qlQueryCacheEntry struct {
	name        string
	encodedPlan string
}

type n1qlJSONPrepData struct {
	EncodedPlan string `json:"encoded_plan"`
	Name        string `json:"name"`
}

func newN1QLQueryComponent(httpComponent httpComponentInterface, cfgMgr configManager, tracer *tracerComponent) *n1qlQueryComponent {
	nqc := &n1qlQueryComponent{
		httpComponent: httpComponent,
		cfgMgr:        cfgMgr,
		queryCache:    newN1qlQueryCache(),
		tracer:        tracer,
	}
	cfgMgr.AddConfigWatcher(nqc)

	return nqc
}

func (nqc *n1qlQueryComponent) OnNewRouteConfig(cfg *routeConfig) {
	if atomic.LoadUint32(&nqc.enhancedPreparedSupported) == 0 &&
		cfg.ContainsClusterCapability(1, "n1ql", "enhancedPreparedStatements") {
		logDebugf("Enabling enhanced prepared statement support")
		// Once supported this can't be unsupported
		nqc.queryCache.Invalidate()
		atomic.StoreUint32(&nqc.enhancedPreparedSupported, 1)
	}
	if cfg.ContainsClusterCapability(1, "n1ql", "readFromReplica") {
		atomic.StoreUint32(&nqc.useReplicaSupported, useReplicaSupportLevelSupported)
	} else {
		atomic.StoreUint32(&nqc.useReplicaSupported, useReplicaSupportLevelUnsupported)
	}
}

// N1QLQuery executes a N1QL query
func (nqc *n1qlQueryComponent) N1QLQuery(opts N1QLQueryOptions, cb N1QLQueryCallback) (PendingOp, error) {
	tracer := nqc.tracer.StartTelemeteryHandler(metricValueServiceQueryValue, "N1QLQuery",
		opts.TraceContext)

	var payloadMap map[string]interface{}
	err := json.Unmarshal(opts.Payload, &payloadMap)
	if err != nil {
		tracer.Finish()
		return nil, wrapN1QLError(nil, "", wrapError(err, "expected a JSON payload"), "", 0)
	}

	statement := getMapValueString(payloadMap, "statement", "")
	clientContextID := getMapValueString(payloadMap, "client_context_id", "")
	readOnly := getMapValueBool(payloadMap, "readonly", false)
	if _, ok := payloadMap["use_replica"]; ok {
		if atomic.LoadUint32(&nqc.useReplicaSupported) == useReplicaSupportLevelUnsupported {
			return nil, wrapN1QLError(nil, "", wrapError(errFeatureNotAvailable, "use replica is not supported by this cluster version"), "", 0)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	ireq := &httpRequest{
		Service:          N1qlService,
		Method:           "POST",
		Path:             "/query/service",
		IsIdempotent:     readOnly,
		UniqueID:         clientContextID,
		Deadline:         opts.Deadline,
		RetryStrategy:    opts.RetryStrategy,
		RootTraceContext: tracer.RootContext(),
		Context:          ctx,
		CancelFunc:       cancel,
		User:             opts.User,
		Endpoint:         opts.Endpoint,
	}

	go func() {
		resp, err := nqc.execute(ireq, payloadMap, statement, time.Now())
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		tracer.Finish()
		cb(resp, nil)
	}()

	return ireq, nil
}

// PreparedN1QLQuery executes a prepared N1QL query
func (nqc *n1qlQueryComponent) PreparedN1QLQuery(opts N1QLQueryOptions, cb N1QLQueryCallback) (PendingOp, error) {
	tracer := nqc.tracer.StartTelemeteryHandler(metricValueServiceQueryValue, "PreparedN1QLQuery", opts.TraceContext)

	ctx, cancel := context.WithCancel(context.Background())
	parentReqForCancel := &httpRequest{
		Context:    ctx,
		CancelFunc: cancel,
	}

	go func() {
		res, err := nqc.executePrepared(ctx, cancel, tracer.RootContext(), opts)
		if err != nil {
			cancel()
			tracer.Finish()
			cb(nil, err)
			return
		}

		tracer.Finish()
		cb(res, nil)
	}()

	return parentReqForCancel, nil
}

func (nqc *n1qlQueryComponent) executePrepared(ctx context.Context, cancel context.CancelFunc,
	traceCtx RequestSpanContext, opts N1QLQueryOptions) (*N1QLRowReader, error) {
	start := time.Now()
	var payloadMap map[string]interface{}
	err := json.Unmarshal(opts.Payload, &payloadMap)
	if err != nil {
		return nil, wrapN1QLError(nil, "", wrapError(err, "expected a JSON payload"), "", 0)
	}

	statement := getMapValueString(payloadMap, "statement", "")
	clientContextID := getMapValueString(payloadMap, "client_context_id", "")
	readOnly := getMapValueBool(payloadMap, "readonly", false)
	if _, ok := payloadMap["use_replica"]; ok {
		if atomic.LoadUint32(&nqc.useReplicaSupported) == useReplicaSupportLevelUnsupported {
			return nil, wrapN1QLError(nil, "", wrapError(errFeatureNotAvailable, "use replica is not supported by this cluster version"), "", 0)
		}
	}
	queryCtx := getMapValueString(payloadMap, "query_context", "")
	statementCtx := n1qlQueryCacheStatementContext{
		Statement: statement,
		Context:   queryCtx,
	}

	cachedStmt := nqc.queryCache.Get(statementCtx)

	enhanced := atomic.LoadUint32(&nqc.enhancedPreparedSupported) == 1

	var req *httpRequest
	if cachedStmt != nil {
		// Attempt to execute our cached query plan
		delete(payloadMap, "statement")
		payloadMap["prepared"] = cachedStmt.name
		if cachedStmt.encodedPlan != "" {
			payloadMap["encoded_plan"] = cachedStmt.encodedPlan
		}

		req = &httpRequest{
			Service:          N1qlService,
			Method:           "POST",
			Path:             "/query/service",
			IsIdempotent:     readOnly,
			UniqueID:         clientContextID,
			Deadline:         opts.Deadline,
			RetryStrategy:    opts.RetryStrategy,
			RootTraceContext: traceCtx,
			Context:          ctx,
			CancelFunc:       cancel,
			User:             opts.User,
			Endpoint:         opts.Endpoint,
		}

		results, err := nqc.execute(req, payloadMap, statement, start)
		if err == nil {
			return results, nil
		}

		retryErr := nqc.preparedStatementMaybeEvictAndRetry(req, err, start, statementCtx)
		if retryErr != nil {
			return nil, retryErr
		}

		logDebugf("Prepared statement execution failed, will attempt reprepare: %v", err)
	}

	delete(payloadMap, "prepared")
	delete(payloadMap, "encoded_plan")
	payloadMap["statement"] = "PREPARE " + statement

	if enhanced {
		payloadMap["auto_execute"] = true
	} else {
		delete(payloadMap, "auto_execute")
	}

	if req == nil {
		req = &httpRequest{
			Service:          N1qlService,
			Method:           "POST",
			Path:             "/query/service",
			IsIdempotent:     readOnly,
			UniqueID:         clientContextID,
			Deadline:         opts.Deadline,
			RetryStrategy:    opts.RetryStrategy,
			RootTraceContext: traceCtx,
			Context:          ctx,
			CancelFunc:       cancel,
			User:             opts.User,
			Endpoint:         opts.Endpoint,
		}
	}

	for {
		var res *N1QLRowReader
		var err error
		if enhanced {
			res, err = nqc.executeEnhPrepared(req, payloadMap, statementCtx, start)
		} else {
			res, err = nqc.executeOldPrepared(req, payloadMap, statementCtx, start)
		}
		if err == nil {
			return res, nil
		}

		err = nqc.preparedStatementMaybeEvictAndRetry(req, err, start, statementCtx)
		if err != nil {
			return nil, err
		}
	}
}

func (nqc *n1qlQueryComponent) preparedStatementMaybeEvictAndRetry(req *httpRequest, originalErr error, start time.Time,
	statementCtx n1qlQueryCacheStatementContext) error {
	var err *N1QLError
	if !errors.As(originalErr, &err) {
		return originalErr
	}

	var retryReason RetryReason
	if len(err.Errors) >= 1 {
		firstErrDesc := err.Errors[0]

		if firstErrDesc.Code == 4040 || firstErrDesc.Code == 4050 || firstErrDesc.Code == 4060 ||
			firstErrDesc.Code == 4070 || firstErrDesc.Code == 4080 || firstErrDesc.Code == 4090 {
			retryReason = QueryPreparedStatementFailureRetryReason

			// If the error is because of a prepared statement issue then we need to evict the cache entry and reprepare.
			nqc.queryCache.Delete(statementCtx)
		}

		if retryReason == nil {
			// n1qlErr is already wrapped here
			return originalErr
		}

		shouldRetry, retryTime := retryOrchMaybeRetry(req, retryReason)
		if !shouldRetry {
			// n1qlErr is already wrapped here
			return originalErr
		}

		select {
		case <-time.After(time.Until(req.Deadline)):
			err := &TimeoutError{
				InnerError:       errUnambiguousTimeout,
				OperationID:      "N1QLQuery",
				Opaque:           req.Identifier(),
				TimeObserved:     time.Since(start),
				RetryReasons:     req.retryReasons,
				RetryAttempts:    req.retryCount,
				LastDispatchedTo: req.Endpoint,
			}
			return wrapN1QLError(req, statementCtx.Statement, err, "", 0)
		case <-time.After(time.Until(retryTime)):
			return nil
		}
	}

	return originalErr
}

func (nqc *n1qlQueryComponent) executeEnhPrepared(ireq *httpRequest, payloadMap map[string]interface{},
	statementCtx n1qlQueryCacheStatementContext, start time.Time) (*N1QLRowReader, error) {
	cacheRes, err := nqc.execute(ireq, payloadMap, statementCtx.Statement, start)
	if err != nil {
		return nil, err
	}

	preparedName, err := cacheRes.PreparedName()
	if err != nil {
		logWarnf("Failed to read prepared name from result: %s", err)
		return cacheRes, nil
	}

	cachedStmt := &n1qlQueryCacheEntry{}
	cachedStmt.name = preparedName

	nqc.queryCache.Put(statementCtx, cachedStmt)

	return cacheRes, nil
}

func (nqc *n1qlQueryComponent) executeOldPrepared(ireq *httpRequest, payloadMap map[string]interface{}, statementCtx n1qlQueryCacheStatementContext,
	start time.Time) (*N1QLRowReader, error) {
	delete(payloadMap, "prepared")
	delete(payloadMap, "encoded_plan")
	delete(payloadMap, "auto_execute")
	prepStatement := "PREPARE " + statementCtx.Statement
	payloadMap["statement"] = prepStatement

	cacheRes, err := nqc.execute(ireq, payloadMap, statementCtx.Statement, start)
	if err != nil {
		return nil, err
	}

	b := cacheRes.NextRow()
	if b == nil {
		var n1qlError *N1QLError
		meta, metaErr := cacheRes.MetaData()
		if metaErr == nil {
			raw, descs, err := parseN1QLError(meta)
			if err != nil {
				n1qlError = wrapN1QLError(ireq, statementCtx.Statement, err, raw, 0)
				n1qlError.Errors = descs
			} else if len(descs) > 0 {
				n1qlError = wrapN1QLError(ireq, statementCtx.Statement, nil, raw, 0)
				n1qlError.Errors = descs
			}
		}
		if n1qlError == nil {
			n1qlError = wrapN1QLError(ireq, statementCtx.Statement, errCliInternalError, "", 0)
		}

		return nil, n1qlError
	}

	var prepData n1qlJSONPrepData
	err = json.Unmarshal(b, &prepData)
	if err != nil {
		return nil, wrapN1QLError(ireq, statementCtx.Statement, err, "", 0)
	}

	cachedStmt := &n1qlQueryCacheEntry{}
	cachedStmt.name = prepData.Name
	cachedStmt.encodedPlan = prepData.EncodedPlan

	nqc.queryCache.Put(statementCtx, cachedStmt)

	// Attempt to execute our cached query plan
	delete(payloadMap, "statement")
	payloadMap["prepared"] = cachedStmt.name
	payloadMap["encoded_plan"] = cachedStmt.encodedPlan

	resp, err := nqc.execute(ireq, payloadMap, statementCtx.Statement, start)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (nqc *n1qlQueryComponent) execute(ireq *httpRequest, payloadMap map[string]interface{}, statementForErr string,
	start time.Time) (*N1QLRowReader, error) {
	for {
		{
			if !ireq.Deadline.IsZero() {
				// Produce an updated payload with the appropriate timeout
				timeoutLeft := time.Until(ireq.Deadline)
				if timeoutLeft <= 0 {
					err := &TimeoutError{
						InnerError:       errUnambiguousTimeout,
						OperationID:      "N1QLQuery",
						Opaque:           ireq.Identifier(),
						TimeObserved:     time.Since(start),
						RetryReasons:     ireq.retryReasons,
						RetryAttempts:    ireq.retryCount,
						LastDispatchedTo: ireq.Endpoint,
					}
					return nil, wrapN1QLError(ireq, statementForErr, err, "", 0)
				}
				payloadMap["timeout"] = timeoutLeft.String()
			}

			newPayload, err := json.Marshal(payloadMap)
			if err != nil {
				return nil, wrapN1QLError(nil, "", wrapError(err, "failed to produce payload"), "", 0)
			}
			ireq.Body = newPayload
		}

		resp, err := nqc.httpComponent.DoInternalHTTPRequest(ireq, false)
		if err != nil {
			if errors.Is(err, ErrRequestCanceled) {
				return nil, err
			}
			// execHTTPRequest will handle retrying due to in-flight socket close based
			// on whether or not IsIdempotent is set on the httpRequest
			return nil, wrapN1QLError(ireq, statementForErr, err, "", 0)
		}

		if resp.StatusCode != 200 {
			n1qlErr := parseN1QLErrorResp(ireq, statementForErr, resp)

			// Note that prepared statement error code retries are handled higher up.
			var retryReason RetryReason
			if len(n1qlErr.Errors) >= 1 {
				firstErrDesc := n1qlErr.Errors[0]

				// See MB-50643 for why this code check is here.
				if firstErrDesc.Retry && firstErrDesc.Code != 12016 {
					retryReason = QueryErrorRetryable
				} else if strings.Contains(firstErrDesc.Message, "queryport.indexNotFound") {
					retryReason = QueryIndexNotFoundRetryReason
				}
			}

			if retryReason == nil {
				// n1qlErr is already wrapped here
				return nil, n1qlErr
			}

			shouldRetry, retryTime := retryOrchMaybeRetry(ireq, retryReason)
			if !shouldRetry {
				// n1qlErr is already wrapped here
				return nil, n1qlErr
			}

			select {
			case <-time.After(time.Until(retryTime)):
				continue
			case <-time.After(time.Until(ireq.Deadline)):
				err := &TimeoutError{
					InnerError:       errUnambiguousTimeout,
					OperationID:      "N1QLQuery",
					Opaque:           ireq.Identifier(),
					TimeObserved:     time.Since(start),
					RetryReasons:     ireq.retryReasons,
					RetryAttempts:    ireq.retryCount,
					LastDispatchedTo: ireq.Endpoint,
				}
				return nil, wrapN1QLError(ireq, statementForErr, err, "", 0)
			}
		}

		streamer, err := newQueryStreamer(resp.Body, "results")
		if err != nil {
			respBody, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				logDebugf("Failed to read response body: %v", readErr)
			}
			return nil, wrapN1QLError(ireq, statementForErr, err, string(respBody), resp.StatusCode)
		}

		return &N1QLRowReader{
			streamer:   streamer,
			endpoint:   resp.Endpoint,
			statement:  statementForErr,
			statusCode: resp.StatusCode,
		}, nil
	}
}
