package gocbcore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"
)

// ViewQueryRowReader providers access to the rows of a view query
type ViewQueryRowReader struct {
	streamer *queryStreamer
}

// NextRow reads the next rows bytes from the stream
func (q *ViewQueryRowReader) NextRow() []byte {
	return q.streamer.NextRow()
}

// Err returns any errors that occurred during streaming.
func (q ViewQueryRowReader) Err() error {
	return q.streamer.Err()
}

// MetaData fetches the non-row bytes streamed in the response.
func (q *ViewQueryRowReader) MetaData() ([]byte, error) {
	return q.streamer.MetaData()
}

// Close immediately shuts down the connection
func (q *ViewQueryRowReader) Close() error {
	return q.streamer.Close()
}

// ViewQueryOptions represents the various options available for a view query.
type ViewQueryOptions struct {
	DesignDocumentName string
	ViewType           string
	ViewName           string
	Options            url.Values
	RetryStrategy      RetryStrategy
	Deadline           time.Time

	// Internal: This should never be used and is not supported.
	User string

	TraceContext RequestSpanContext
}

func wrapViewQueryError(req *httpRequest, ddoc, view string, err error, errBody string, statusCode int) *ViewError {
	if err == nil {
		err = errors.New("view error")
	}

	ierr := &ViewError{
		InnerError: err,
	}

	if req != nil {
		ierr.Endpoint = req.Endpoint
		ierr.RetryAttempts = req.RetryAttempts()
		ierr.RetryReasons = req.RetryReasons()
	}

	ierr.ErrorText = errBody
	ierr.HTTPResponseCode = statusCode
	ierr.DesignDocumentName = ddoc
	ierr.ViewName = view

	return ierr
}

func parseViewQueryError(req *httpRequest, ddoc, view string, resp *HTTPResponse) *ViewError {
	var err error
	var errorDescs []ViewQueryErrorDesc

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr == nil {
		var errsMap map[string]string
		var errsArr []string

		if err := json.Unmarshal(respBody, &errsArr); err == nil {
			errorDescs = make([]ViewQueryErrorDesc, len(errsArr))
			for errIdx, errMessage := range errsArr {
				errorDescs[errIdx] = ViewQueryErrorDesc{
					SourceNode: "",
					Message:    errMessage,
				}
			}
		} else if err := json.Unmarshal(respBody, &errsMap); err == nil {
			for errNode, errMessage := range errsMap {
				errorDescs = append(errorDescs, ViewQueryErrorDesc{
					SourceNode: errNode,
					Message:    errMessage,
				})
			}
		}
	}

	if resp.StatusCode == 401 {
		err = errAuthenticationFailure
	} else if resp.StatusCode == 404 {
		err = errViewNotFound
	}

	if len(errorDescs) >= 1 {
		firstErrMsg := errorDescs[0].Message

		if strings.Contains(firstErrMsg, "not_found") {
			err = errViewNotFound
		}
	}

	var errText string
	if err == nil {
		errText = string(respBody)
	}

	errOut := wrapViewQueryError(req, ddoc, view, err, errText, resp.StatusCode)
	errOut.Errors = errorDescs
	return errOut
}

type viewQueryComponent struct {
	httpComponent *httpComponent
	tracer        *tracerComponent
}

func newViewQueryComponent(httpComponent *httpComponent, tracer *tracerComponent) *viewQueryComponent {
	return &viewQueryComponent{
		httpComponent: httpComponent,
		tracer:        tracer,
	}
}

// ViewQuery executes a view query
func (vqc *viewQueryComponent) ViewQuery(opts ViewQueryOptions, cb ViewQueryCallback) (PendingOp, error) {
	tracer := vqc.tracer.StartTelemeteryHandler(metricValueServiceViewsValue, "ViewQuery", opts.TraceContext)

	reqURI := fmt.Sprintf("/_design/%s/%s/%s?%s",
		opts.DesignDocumentName, opts.ViewType, opts.ViewName, opts.Options.Encode())

	ctx, cancel := context.WithCancel(context.Background())
	ireq := &httpRequest{
		Service:          CapiService,
		Method:           "GET",
		Path:             reqURI,
		IsIdempotent:     true,
		Deadline:         opts.Deadline,
		RetryStrategy:    opts.RetryStrategy,
		RootTraceContext: tracer.RootContext(),
		Context:          ctx,
		CancelFunc:       cancel,
		User:             opts.User,
	}

	ddoc := opts.DesignDocumentName
	view := opts.ViewName

	go func() {
		res, err := vqc.viewQuery(ireq, ddoc, view)
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

func (vqc *viewQueryComponent) viewQuery(ireq *httpRequest, ddoc, view string) (*ViewQueryRowReader, error) {
	resp, err := vqc.httpComponent.DoInternalHTTPRequest(ireq, false)
	if err != nil {
		if errors.Is(err, ErrRequestCanceled) {
			return nil, err
		}
		// execHTTPRequest will handle retrying due to in-flight socket close based
		// on whether or not IsIdempotent is set on the httpRequest
		return nil, wrapViewQueryError(ireq, ddoc, view, err, "", 0)
	}

	if resp.StatusCode != 200 {
		viewErr := parseViewQueryError(ireq, ddoc, view, resp)

		// viewErr is already wrapped here
		return nil, viewErr
	}

	streamer, err := newQueryStreamer(resp.Body, "rows")
	if err != nil {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			logDebugf("Failed to read response body: %v", readErr)
		}
		return nil, wrapViewQueryError(ireq, ddoc, view, err, string(respBody), resp.StatusCode)
	}

	return &ViewQueryRowReader{
		streamer: streamer,
	}, nil
}
