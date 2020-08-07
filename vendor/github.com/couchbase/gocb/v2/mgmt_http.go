package gocb

import (
	"io"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v9"
)

type mgmtRequest struct {
	Service      ServiceType
	Method       string
	Path         string
	Body         []byte
	Headers      map[string]string
	ContentType  string
	IsIdempotent bool
	UniqueID     string

	Timeout       time.Duration
	RetryStrategy RetryStrategy

	parentSpan requestSpanContext
}

type mgmtResponse struct {
	Endpoint   string
	StatusCode uint32
	Body       io.ReadCloser
}

type mgmtProvider interface {
	executeMgmtRequest(req mgmtRequest) (*mgmtResponse, error)
}

func (c *Cluster) executeMgmtRequest(req mgmtRequest) (mgmtRespOut *mgmtResponse, errOut error) {
	timeout := req.Timeout
	if timeout == 0 {
		timeout = c.timeoutsConfig.ManagementTimeout
	}

	provider, err := c.getHTTPProvider()
	if err != nil {
		return nil, err
	}

	retryStrategy := c.retryStrategyWrapper
	if req.RetryStrategy != nil {
		retryStrategy = newRetryStrategyWrapper(req.RetryStrategy)
	}

	corereq := &gocbcore.HTTPRequest{
		Service:       gocbcore.ServiceType(req.Service),
		Method:        req.Method,
		Path:          req.Path,
		Body:          req.Body,
		Headers:       req.Headers,
		ContentType:   req.ContentType,
		IsIdempotent:  req.IsIdempotent,
		UniqueID:      req.UniqueID,
		Deadline:      time.Now().Add(timeout),
		RetryStrategy: retryStrategy,
		TraceContext:  req.parentSpan,
	}

	coreresp, err := provider.DoHTTPRequest(corereq)
	if err != nil {
		return nil, makeGenericHTTPError(err, corereq, coreresp)
	}

	resp := &mgmtResponse{
		Endpoint:   coreresp.Endpoint,
		StatusCode: uint32(coreresp.StatusCode),
		Body:       coreresp.Body,
	}
	return resp, nil
}

func (b *Bucket) executeMgmtRequest(req mgmtRequest) (mgmtRespOut *mgmtResponse, errOut error) {
	timeout := req.Timeout
	if timeout == 0 {
		timeout = b.timeoutsConfig.ManagementTimeout
	}

	provider, err := b.connectionManager.getHTTPProvider()
	if err != nil {
		return nil, err
	}

	retryStrategy := b.retryStrategyWrapper
	if req.RetryStrategy != nil {
		retryStrategy = newRetryStrategyWrapper(req.RetryStrategy)
	}

	corereq := &gocbcore.HTTPRequest{
		Service:       gocbcore.ServiceType(req.Service),
		Method:        req.Method,
		Path:          req.Path,
		Body:          req.Body,
		Headers:       req.Headers,
		ContentType:   req.ContentType,
		IsIdempotent:  req.IsIdempotent,
		UniqueID:      req.UniqueID,
		Deadline:      time.Now().Add(timeout),
		RetryStrategy: retryStrategy,
	}

	coreresp, err := provider.DoHTTPRequest(corereq)
	if err != nil {
		return nil, makeGenericHTTPError(err, corereq, coreresp)
	}

	resp := &mgmtResponse{
		Endpoint:   coreresp.Endpoint,
		StatusCode: uint32(coreresp.StatusCode),
		Body:       coreresp.Body,
	}
	return resp, nil
}

func ensureBodyClosed(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		logDebugf("Failed to close socket: %v", err)
	}
}
