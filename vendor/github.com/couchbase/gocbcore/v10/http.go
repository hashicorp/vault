package gocbcore

import (
	"context"
	"errors"
	"io"
	"sort"
	"sync/atomic"
	"time"
)

type httpRequest struct {
	Service          ServiceType
	Endpoint         string
	Method           string
	Path             string
	Username         string
	Password         string
	Headers          map[string]string
	ContentType      string
	Body             []byte
	IsIdempotent     bool
	UniqueID         string
	Deadline         time.Time
	RetryStrategy    RetryStrategy
	RootTraceContext RequestSpanContext
	// Whilst the http component will handle deadlines itself this context can be use from places like Ping which
	// need to also be able to cancel the context for other reasons.
	Context    context.Context
	CancelFunc context.CancelFunc

	User string

	retryCount   uint32
	retryReasons []RetryReason
}

func (hr *httpRequest) retryStrategy() RetryStrategy {
	return hr.RetryStrategy
}

func (hr *httpRequest) Cancel() {
	if hr.CancelFunc != nil {
		hr.CancelFunc()
	}
}

func (hr *httpRequest) RetryAttempts() uint32 {
	return atomic.LoadUint32(&hr.retryCount)
}

func (hr *httpRequest) Identifier() string {
	return hr.UniqueID
}

func (hr *httpRequest) Idempotent() bool {
	return hr.IsIdempotent
}

func (hr *httpRequest) RetryReasons() []RetryReason {
	return hr.retryReasons
}

func (hr *httpRequest) recordRetryAttempt(reason RetryReason) {
	atomic.AddUint32(&hr.retryCount, 1)
	idx := sort.Search(len(hr.retryReasons), func(i int) bool {
		return hr.retryReasons[i] == reason
	})

	// if idx is out of the range of retryReasons then it wasn't found.
	if idx > len(hr.retryReasons)-1 {
		hr.retryReasons = append(hr.retryReasons, reason)
	}
}

// HTTPRequest contains the description of an HTTP request to perform.
type HTTPRequest struct {
	Service       ServiceType
	Method        string
	Endpoint      string
	Path          string
	Username      string
	Password      string
	Body          []byte
	Headers       map[string]string
	ContentType   string
	IsIdempotent  bool
	UniqueID      string
	Deadline      time.Time
	RetryStrategy RetryStrategy

	// Internal: This should never be used and is not supported.
	User string

	TraceContext RequestSpanContext
}

// HTTPResponse encapsulates the response from an HTTP request.
type HTTPResponse struct {
	Endpoint      string
	StatusCode    int
	ContentLength int64
	Body          io.ReadCloser
}

func wrapHTTPError(req *httpRequest, err error) HTTPError {
	if err == nil {
		err = errors.New("http error")
	}

	ierr := HTTPError{
		InnerError: err,
	}

	if req != nil {
		ierr.Endpoint = req.Endpoint
		ierr.UniqueID = req.UniqueID
		ierr.RetryAttempts = req.RetryAttempts()
		ierr.RetryReasons = req.RetryReasons()
	}

	return ierr
}
