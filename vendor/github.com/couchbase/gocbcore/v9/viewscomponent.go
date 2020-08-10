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

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

func wrapViewQueryError(req *httpRequest, ddoc, view string, err error) *ViewError {
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

		if err := json.Unmarshal(respBody, &errsArr); err != nil {
			errorDescs = make([]ViewQueryErrorDesc, len(errsArr))
			for errIdx, errMessage := range errsArr {
				errorDescs[errIdx] = ViewQueryErrorDesc{
					SourceNode: "",
					Message:    errMessage,
				}
			}
		} else if err := json.Unmarshal(respBody, &errsMap); err != nil {
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

	errOut := wrapViewQueryError(req, ddoc, view, err)
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
	tracer := vqc.tracer.CreateOpTrace("ViewQuery", opts.TraceContext)
	defer tracer.Finish()

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
	}

	ddoc := opts.DesignDocumentName
	view := opts.ViewName

	go func() {
		resp, err := vqc.httpComponent.DoInternalHTTPRequest(ireq, false)
		if err != nil {
			cancel()
			// execHTTPRequest will handle retrying due to in-flight socket close based
			// on whether or not IsIdempotent is set on the httpRequest
			cb(nil, wrapViewQueryError(ireq, ddoc, view, err))
			return
		}

		if resp.StatusCode != 200 {
			viewErr := parseViewQueryError(ireq, ddoc, view, resp)

			cancel()
			// viewErr is already wrapped here
			cb(nil, viewErr)
			return
		}

		streamer, err := newQueryStreamer(resp.Body, "rows")
		if err != nil {
			cancel()
			cb(nil, wrapViewQueryError(ireq, ddoc, view, err))
			return
		}

		cb(&ViewQueryRowReader{
			streamer: streamer,
		}, nil)
	}()

	return ireq, nil
}
