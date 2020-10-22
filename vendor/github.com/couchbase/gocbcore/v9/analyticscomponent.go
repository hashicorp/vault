package gocbcore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

// AnalyticsRowReader providers access to the rows of a analytics query
type AnalyticsRowReader struct {
	streamer *queryStreamer
}

// NextRow reads the next rows bytes from the stream
func (q *AnalyticsRowReader) NextRow() []byte {
	return q.streamer.NextRow()
}

// Err returns any errors that occurred during streaming.
func (q AnalyticsRowReader) Err() error {
	return q.streamer.Err()
}

// MetaData fetches the non-row bytes streamed in the response.
func (q *AnalyticsRowReader) MetaData() ([]byte, error) {
	return q.streamer.MetaData()
}

// Close immediately shuts down the connection
func (q *AnalyticsRowReader) Close() error {
	return q.streamer.Close()
}

// AnalyticsQueryOptions represents the various options available for an analytics query.
type AnalyticsQueryOptions struct {
	Payload       []byte
	Priority      int
	RetryStrategy RetryStrategy
	Deadline      time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

func wrapAnalyticsError(req *httpRequest, statement string, err error) *AnalyticsError {
	if err == nil {
		err = errors.New("analytics error")
	}

	ierr := &AnalyticsError{
		InnerError: err,
	}

	if req != nil {
		ierr.Endpoint = req.Endpoint
		ierr.ClientContextID = req.UniqueID
		ierr.RetryAttempts = req.RetryAttempts()
		ierr.RetryReasons = req.RetryReasons()
	}

	ierr.Statement = statement

	return ierr
}

type jsonAnalyticsError struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

type jsonAnalyticsErrorResponse struct {
	Errors []jsonAnalyticsError
}

func parseAnalyticsError(req *httpRequest, statement string, resp *HTTPResponse) *AnalyticsError {
	var err error
	var errorDescs []AnalyticsErrorDesc

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr == nil {
		var respParse jsonAnalyticsErrorResponse
		parseErr := json.Unmarshal(respBody, &respParse)
		if parseErr == nil {

			for _, jsonErr := range respParse.Errors {
				errorDescs = append(errorDescs, AnalyticsErrorDesc{
					Code:    jsonErr.Code,
					Message: jsonErr.Msg,
				})
			}
		}
	}

	if len(errorDescs) >= 1 {
		firstErr := errorDescs[0]
		errCode := firstErr.Code
		errCodeGroup := errCode / 1000

		if errCodeGroup == 25 {
			err = errInternalServerFailure
		}
		if errCodeGroup == 20 {
			err = errAuthenticationFailure
		}
		if errCodeGroup == 24 {
			err = errCompilationFailure
		}
		if errCode == 23000 || errCode == 23003 {
			err = errTemporaryFailure
		}
		if errCode == 24000 {
			err = errParsingFailure
		}
		if errCode == 24047 {
			err = errIndexNotFound
		}
		if errCode == 24048 {
			err = errIndexExists
		}

		if errCode == 23007 {
			err = errJobQueueFull
		}
		if errCode == 24025 || errCode == 24044 || errCode == 24045 {
			err = errDatasetNotFound
		}
		if errCode == 24034 {
			err = errDataverseNotFound
		}
		if errCode == 24040 {
			err = errDatasetExists
		}
		if errCode == 24039 {
			err = errDataverseExists
		}
		if errCode == 24006 {
			err = errLinkNotFound
		}
	}

	errOut := wrapAnalyticsError(req, statement, err)
	errOut.Errors = errorDescs
	return errOut
}

type analyticsQueryComponent struct {
	httpComponent *httpComponent
	tracer        *tracerComponent
}

func newAnalyticsQueryComponent(httpComponent *httpComponent, tracer *tracerComponent) *analyticsQueryComponent {
	return &analyticsQueryComponent{
		httpComponent: httpComponent,
		tracer:        tracer,
	}
}

// AnalyticsQuery executes an analytics query
func (aqc *analyticsQueryComponent) AnalyticsQuery(opts AnalyticsQueryOptions, cb AnalyticsQueryCallback) (PendingOp, error) {
	tracer := aqc.tracer.CreateOpTrace("AnalyticsQuery", opts.TraceContext)
	defer tracer.Finish()

	var payloadMap map[string]interface{}
	err := json.Unmarshal(opts.Payload, &payloadMap)
	if err != nil {
		return nil, wrapAnalyticsError(nil, "", wrapError(err, "expected a JSON payload"))
	}

	statement := getMapValueString(payloadMap, "statement", "")
	clientContextID := getMapValueString(payloadMap, "client_context_id", "")
	readOnly := getMapValueBool(payloadMap, "readonly", false)

	ctx, cancel := context.WithCancel(context.Background())
	ireq := &httpRequest{
		Service: CbasService,
		Method:  "POST",
		Path:    "/query/service",
		Headers: map[string]string{
			"Analytics-Priority": fmt.Sprintf("%d", opts.Priority),
		},
		Body:             opts.Payload,
		IsIdempotent:     readOnly,
		UniqueID:         clientContextID,
		Deadline:         opts.Deadline,
		RetryStrategy:    opts.RetryStrategy,
		RootTraceContext: tracer.RootContext(),
		Context:          ctx,
		CancelFunc:       cancel,
	}
	start := time.Now()

	go func() {
	ExecuteLoop:
		for {
			{ // Produce an updated payload with the appropriate timeout
				timeoutLeft := time.Until(ireq.Deadline)
				payloadMap["timeout"] = timeoutLeft.String()

				newPayload, err := json.Marshal(payloadMap)
				if err != nil {
					cancel()
					cb(nil, wrapAnalyticsError(nil, "", wrapError(err, "failed to produce payload")))
					return
				}
				ireq.Body = newPayload
			}

			resp, err := aqc.httpComponent.DoInternalHTTPRequest(ireq, false)
			if err != nil {
				cancel()
				// execHTTPRequest will handle retrying due to in-flight socket close based
				// on whether or not IsIdempotent is set on the httpRequest
				cb(nil, wrapAnalyticsError(ireq, statement, err))
				return
			}

			if resp.StatusCode != 200 {
				analyticsErr := parseAnalyticsError(ireq, statement, resp)

				var retryReason RetryReason
				if len(analyticsErr.Errors) >= 1 {
					firstErrDesc := analyticsErr.Errors[0]

					if firstErrDesc.Code == 23000 {
						retryReason = AnalyticsTemporaryFailureRetryReason
					} else if firstErrDesc.Code == 23003 {
						retryReason = AnalyticsTemporaryFailureRetryReason
					} else if firstErrDesc.Code == 23007 {
						retryReason = AnalyticsTemporaryFailureRetryReason
					}
				}

				if retryReason == nil {
					cancel()
					// analyticsErr is already wrapped here
					cb(nil, analyticsErr)
					return
				}

				shouldRetry, retryTime := retryOrchMaybeRetry(ireq, retryReason)
				if !shouldRetry {
					cancel()
					// analyticsErr is already wrapped here
					cb(nil, analyticsErr)
					return
				}

				select {
				case <-time.After(time.Until(retryTime)):
					continue ExecuteLoop
				case <-time.After(time.Until(ireq.Deadline)):
					cancel()
					err := &TimeoutError{
						InnerError:       errUnambiguousTimeout,
						OperationID:      "http",
						Opaque:           ireq.Identifier(),
						TimeObserved:     time.Since(start),
						RetryReasons:     ireq.retryReasons,
						RetryAttempts:    ireq.retryCount,
						LastDispatchedTo: ireq.Endpoint,
					}
					cb(nil, wrapAnalyticsError(ireq, statement, err))
					return
				}
			}

			streamer, err := newQueryStreamer(resp.Body, "results")
			if err != nil {
				cancel()
				cb(nil, wrapAnalyticsError(ireq, statement, err))
				return
			}

			cb(&AnalyticsRowReader{
				streamer: streamer,
			}, nil)
			return
		}
	}()

	return ireq, nil
}
