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
	streamer   *queryStreamer
	statement  string
	statusCode int
}

// NextRow reads the next rows bytes from the stream
func (q *AnalyticsRowReader) NextRow() []byte {
	return q.streamer.NextRow()
}

// Err returns any errors that occurred during streaming.
func (q AnalyticsRowReader) Err() error {
	err := q.streamer.Err()
	if err != nil {
		return err
	}

	meta, metaErr := q.streamer.MetaData()
	if metaErr != nil {
		return metaErr
	}

	raw, descs, err := parseAnalyticsError(meta)
	if err != nil {
		return &AnalyticsError{
			InnerError:       err,
			Errors:           descs,
			ErrorText:        raw,
			Statement:        q.statement,
			HTTPResponseCode: q.statusCode,
		}
	}
	if len(descs) > 0 {
		return &AnalyticsError{
			InnerError:       errors.New("analytics error"),
			Errors:           descs,
			ErrorText:        raw,
			Statement:        q.statement,
			HTTPResponseCode: q.statusCode,
		}
	}

	return nil
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

	// Internal: This should never be used and is not supported.
	User string

	TraceContext RequestSpanContext
}

func wrapAnalyticsError(req *httpRequest, statement string, err error, errBody string, statusCode int) *AnalyticsError {
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

	ierr.ErrorText = errBody
	ierr.Statement = statement
	ierr.HTTPResponseCode = statusCode

	return ierr
}

type jsonAnalyticsError struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

type jsonAnalyticsErrorResponse struct {
	Errors json.RawMessage
}

func parseAnalyticsErrorResp(req *httpRequest, statement string, resp *HTTPResponse) *AnalyticsError {
	var errorDescs []AnalyticsErrorDesc
	var err error
	var raw string
	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr == nil {
		raw, errorDescs, err = parseAnalyticsError(respBody)
	}
	errOut := wrapAnalyticsError(req, statement, err, raw, resp.StatusCode)
	errOut.Errors = errorDescs
	return errOut
}

func parseAnalyticsError(respBody []byte) (string, []AnalyticsErrorDesc, error) {
	var err error
	var errorDescs []AnalyticsErrorDesc

	var rawRespParse jsonAnalyticsErrorResponse
	parseErr := json.Unmarshal(respBody, &rawRespParse)
	if parseErr != nil {
		return "", nil, nil
	}

	var respParse []jsonAnalyticsError
	parseErr = json.Unmarshal(rawRespParse.Errors, &respParse)
	if parseErr == nil {
		for _, jsonErr := range respParse {
			errorDescs = append(errorDescs, AnalyticsErrorDesc{
				Code:    jsonErr.Code,
				Message: jsonErr.Msg,
			})
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
	var rawErrors string
	if err == nil && len(rawRespParse.Errors) > 0 {
		// Only populate if this is an error that we don't recognise.
		rawErrors = string(rawRespParse.Errors)
	}

	return rawErrors, errorDescs, err
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
	tracer := aqc.tracer.StartTelemeteryHandler(metricValueServiceAnalyticsValue, "AnalyticsQuery", opts.TraceContext)

	var payloadMap map[string]interface{}
	err := json.Unmarshal(opts.Payload, &payloadMap)
	if err != nil {
		tracer.Finish()
		return nil, wrapAnalyticsError(nil, "", wrapError(err, "expected a JSON payload"), "", 0)
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
		User:             opts.User,
	}

	go func() {
		res, err := aqc.analyticsQuery(ireq, payloadMap, statement, tracer.StartTime())
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

func (aqc *analyticsQueryComponent) analyticsQuery(ireq *httpRequest, payloadMap map[string]interface{},
	statement string, startTime time.Time) (*AnalyticsRowReader, error) {
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
					return nil, wrapAnalyticsError(ireq, statement, err, "", 0)
				}
				payloadMap["timeout"] = timeoutLeft.String()
			}

			newPayload, err := json.Marshal(payloadMap)
			if err != nil {
				return nil, wrapAnalyticsError(nil, statement, wrapError(err, "failed to produce payload"), "", 0)
			}
			ireq.Body = newPayload
		}

		resp, err := aqc.httpComponent.DoInternalHTTPRequest(ireq, false)
		if err != nil {
			if errors.Is(err, ErrRequestCanceled) {
				return nil, err
			}
			// execHTTPRequest will handle retrying due to in-flight socket close based
			// on whether or not IsIdempotent is set on the httpRequest
			return nil, wrapAnalyticsError(ireq, statement, err, "", 0)
		}

		if resp.StatusCode != 200 {
			analyticsErr := parseAnalyticsErrorResp(ireq, statement, resp)

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
				// analyticsErr is already wrapped here
				return nil, analyticsErr
			}

			shouldRetry, retryTime := retryOrchMaybeRetry(ireq, retryReason)
			if !shouldRetry {
				// analyticsErr is already wrapped here
				return nil, analyticsErr
			}

			select {
			case <-time.After(time.Until(retryTime)):
				continue
			case <-time.After(time.Until(ireq.Deadline)):
				err := &TimeoutError{
					InnerError:       errUnambiguousTimeout,
					OperationID:      "http",
					Opaque:           ireq.Identifier(),
					TimeObserved:     time.Since(startTime),
					RetryReasons:     ireq.retryReasons,
					RetryAttempts:    ireq.retryCount,
					LastDispatchedTo: ireq.Endpoint,
				}
				return nil, wrapAnalyticsError(ireq, statement, err, "", 0)
			}
		}

		streamer, err := newQueryStreamer(resp.Body, "results")
		if err != nil {
			respBody, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				logDebugf("Failed to read response body: %v", readErr)
			}
			return nil, wrapAnalyticsError(ireq, statement, err, string(respBody), resp.StatusCode)
		}

		return &AnalyticsRowReader{
			streamer: streamer,
		}, nil
	}
}
