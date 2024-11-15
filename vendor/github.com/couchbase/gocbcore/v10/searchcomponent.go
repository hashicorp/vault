package gocbcore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"sync"
	"time"
)

// SearchRowReader providers access to the rows of a view query
type SearchRowReader struct {
	streamer *queryStreamer
}

// NextRow reads the next rows bytes from the stream
func (q *SearchRowReader) NextRow() []byte {
	return q.streamer.NextRow()
}

// Err returns any errors that occurred during streaming.
func (q SearchRowReader) Err() error {
	return q.streamer.Err()
}

// MetaData fetches the non-row bytes streamed in the response.
func (q *SearchRowReader) MetaData() ([]byte, error) {
	return q.streamer.MetaData()
}

// Close immediately shuts down the connection
func (q *SearchRowReader) Close() error {
	return q.streamer.Close()
}

// SearchQueryOptions represents the various options available for a search query.
type SearchQueryOptions struct {
	BucketName    string
	ScopeName     string
	IndexName     string
	Payload       []byte
	RetryStrategy RetryStrategy
	Deadline      time.Time

	// Internal: This should never be used and is not supported.
	User string

	TraceContext RequestSpanContext
}

type jsonSearchErrorResponse struct {
	Error string
}

func wrapSearchError(req *httpRequest, indexName string, query interface{}, err error, statusCode int) *SearchError {
	if err == nil {
		err = errors.New("search error")
	}

	ierr := &SearchError{
		InnerError: err,
	}

	if req != nil {
		ierr.Endpoint = req.Endpoint
		ierr.RetryAttempts = req.RetryAttempts()
		ierr.RetryReasons = req.RetryReasons()
	}

	ierr.HTTPResponseCode = statusCode
	ierr.IndexName = indexName
	ierr.Query = query

	return ierr
}

func parseSearchError(req *httpRequest, indexName string, query interface{}, resp *HTTPResponse) *SearchError {
	var err error
	var errMsg string

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr == nil {
		var respParse jsonSearchErrorResponse
		parseErr := json.Unmarshal(respBody, &respParse)
		if parseErr == nil {
			errMsg = respParse.Error
		}
	}

	if resp.StatusCode == 500 {
		err = errInternalServerFailure
	}
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		err = errAuthenticationFailure
	}
	if resp.StatusCode == 400 && strings.Contains(errMsg, "index not found") {
		err = errIndexNotFound
	}
	if resp.StatusCode == 429 {
		if strings.Contains(errMsg, "num_concurrent_requests") {
			err = errRateLimitedFailure
		} else if strings.Contains(errMsg, "num_queries_per_min") {
			err = errRateLimitedFailure
		} else if strings.Contains(errMsg, "ingress_mib_per_min") {
			err = errRateLimitedFailure
		} else if strings.Contains(errMsg, "egress_mib_per_min") {
			err = errRateLimitedFailure
		}
	}

	errOut := wrapSearchError(req, indexName, query, err, resp.StatusCode)
	errOut.ErrorText = errMsg
	return errOut
}

type SearchCapability uint32

const (
	SearchCapabilityScopedIndexes SearchCapability = iota
	SearchCapabilityVectorSearch
)

type searchQueryComponent struct {
	httpComponent *httpComponent
	cfgMgr        configManager
	tracer        *tracerComponent

	caps     map[SearchCapability]CapabilityStatus
	capsLock sync.RWMutex
}

func newSearchQueryComponent(httpComponent *httpComponent, cfgMgr configManager, tracer *tracerComponent) *searchQueryComponent {
	sqc := &searchQueryComponent{
		httpComponent: httpComponent,
		cfgMgr:        cfgMgr,
		tracer:        tracer,

		caps: map[SearchCapability]CapabilityStatus{
			SearchCapabilityVectorSearch:  CapabilityStatusUnknown,
			SearchCapabilityScopedIndexes: CapabilityStatusUnknown,
		},
	}
	cfgMgr.AddConfigWatcher(sqc)

	return sqc
}

func (sqc *searchQueryComponent) OnNewRouteConfig(cfg *routeConfig) {
	sqc.capsLock.Lock()
	defer sqc.capsLock.Unlock()

	if cfg.ContainsClusterCapability(1, "search", "vectorSearch") {
		sqc.caps[SearchCapabilityVectorSearch] = CapabilityStatusSupported
	} else {
		sqc.caps[SearchCapabilityVectorSearch] = CapabilityStatusUnsupported
	}

	if cfg.ContainsClusterCapability(1, "search", "scopedSearchIndex") {
		sqc.caps[SearchCapabilityScopedIndexes] = CapabilityStatusSupported
	} else {
		sqc.caps[SearchCapabilityScopedIndexes] = CapabilityStatusUnsupported
	}
}

func (sqc *searchQueryComponent) capabilityStatus(cap SearchCapability) CapabilityStatus {
	sqc.capsLock.RLock()
	defer sqc.capsLock.RUnlock()

	status, ok := sqc.caps[cap]
	if !ok {
		return CapabilityStatusUnsupported
	}

	return status
}

// SearchQuery executes a Search query
func (sqc *searchQueryComponent) SearchQuery(opts SearchQueryOptions, cb SearchQueryCallback) (PendingOp, error) {
	tracer := sqc.tracer.StartTelemeteryHandler(metricValueServiceSearchValue, "SearchQuery", opts.TraceContext)

	var payloadMap map[string]interface{}
	err := json.Unmarshal(opts.Payload, &payloadMap)
	if err != nil {
		tracer.Finish()
		return nil, wrapSearchError(nil, "", nil, wrapError(err, "expected a JSON payload"), 0)
	}

	var ctlMap map[string]interface{}
	if foundCtlMap, ok := payloadMap["ctl"]; ok {
		if coercedCtlMap, ok := foundCtlMap.(map[string]interface{}); ok {
			ctlMap = coercedCtlMap
		} else {
			tracer.Finish()
			return nil, wrapSearchError(nil, "", nil,
				wrapError(errInvalidArgument, "expected ctl to be a map"), 0)
		}
	} else {
		ctlMap = make(map[string]interface{})
	}

	if opts.BucketName != "" && opts.ScopeName != "" {
		if sqc.capabilityStatus(SearchCapabilityScopedIndexes) == CapabilityStatusUnsupported {
			return nil, wrapSearchError(nil, "", nil,
				wrapError(errFeatureNotAvailable, "scoped search indexes are not supported by this cluster version"), 0)
		}
	}

	if _, ok := payloadMap["knn"]; ok {
		if sqc.capabilityStatus(SearchCapabilityVectorSearch) == CapabilityStatusUnsupported {
			return nil, wrapSearchError(nil, "", nil,
				wrapError(errFeatureNotAvailable, "vector search is not supported by this cluster version"), 0)
		}
	}

	indexName := opts.IndexName
	query := payloadMap["query"]

	ctx, cancel := context.WithCancel(context.Background())
	var reqURI string
	if opts.BucketName != "" && opts.ScopeName != "" {
		reqURI = fmt.Sprintf("/api/bucket/%s/scope/%s/index/%s/query",
			url.PathEscape(opts.BucketName), url.PathEscape(opts.ScopeName), url.PathEscape(opts.IndexName))
	} else {
		reqURI = fmt.Sprintf("/api/index/%s/query", url.PathEscape(opts.IndexName))
	}
	ireq := &httpRequest{
		Service:          FtsService,
		Method:           "POST",
		Path:             reqURI,
		Body:             opts.Payload,
		IsIdempotent:     true,
		Deadline:         opts.Deadline,
		RetryStrategy:    opts.RetryStrategy,
		RootTraceContext: tracer.RootContext(),
		Context:          ctx,
		CancelFunc:       cancel,
		User:             opts.User,
	}

	go func() {
		res, err := sqc.searchQuery(ireq, indexName, query, payloadMap, ctlMap, tracer.StartTime())
		if err != nil {
			cancel()
			tracer.Finish()
			cb(nil, err)
			return
		}

		tracer.Finish()
		cb(res, nil)
	}()

	return ireq, nil
}

func (sqc *searchQueryComponent) searchQuery(ireq *httpRequest, indexName string, query interface{}, payloadMap map[string]interface{},
	ctlMap map[string]interface{}, startTime time.Time) (*SearchRowReader, error) {
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
						TimeObserved:     time.Since(startTime),
						RetryReasons:     ireq.retryReasons,
						RetryAttempts:    ireq.retryCount,
						LastDispatchedTo: ireq.Endpoint,
					}
					return nil, wrapSearchError(nil, indexName, query, err, 0)
				}
				ctlMap["timeout"] = timeoutLeft / time.Millisecond
				payloadMap["ctl"] = ctlMap
			}

			newPayload, err := json.Marshal(payloadMap)
			if err != nil {
				return nil, wrapSearchError(nil, indexName, query,
					wrapError(err, "failed to produce payload"), 0)
			}
			ireq.Body = newPayload
		}

		resp, err := sqc.httpComponent.DoInternalHTTPRequest(ireq, false)
		if err != nil {
			if errors.Is(err, ErrRequestCanceled) {
				return nil, err
			}
			// execHTTPRequest will handle retrying due to in-flight socket close based
			// on whether or not IsIdempotent is set on the httpRequest
			return nil, wrapSearchError(ireq, indexName, query, err, 0)
		}

		if resp.StatusCode != 200 {
			searchErr := parseSearchError(ireq, indexName, query, resp)

			var retryReason RetryReason
			if searchErr.HTTPResponseCode == 429 && !errors.Is(searchErr, ErrRateLimitedFailure) {
				retryReason = SearchTooManyRequestsRetryReason
			}

			if retryReason == nil {
				// searchErr is already wrapped here
				return nil, searchErr
			}

			shouldRetry, retryTime := retryOrchMaybeRetry(ireq, retryReason)
			if !shouldRetry {
				// searchErr is already wrapped here
				return nil, searchErr
			}

			select {
			case <-time.After(time.Until(retryTime)):
				continue
			case <-time.After(time.Until(ireq.Deadline)):
				err := &TimeoutError{
					InnerError:       errUnambiguousTimeout,
					OperationID:      "SearchQuery",
					Opaque:           ireq.Identifier(),
					TimeObserved:     time.Since(startTime),
					RetryReasons:     ireq.retryReasons,
					RetryAttempts:    ireq.retryCount,
					LastDispatchedTo: ireq.Endpoint,
				}
				return nil, wrapSearchError(ireq, indexName, query, err, 0)
			}
		}

		streamer, err := newQueryStreamer(resp.Body, "hits")
		if err != nil {
			respBody, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				logDebugf("Failed to read response body: %v", readErr)
			}
			sErr := wrapSearchError(ireq, indexName, query, err, resp.StatusCode)
			sErr.ErrorText = string(respBody)
			return nil, sErr
		}

		return &SearchRowReader{
			streamer: streamer,
		}, nil
	}
}
