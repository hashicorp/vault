package gocb

import (
	"context"
	"io"
	"strings"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v10"
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
	Endpoint     string

	Timeout       time.Duration
	RetryStrategy RetryStrategy

	parentSpanCtx RequestSpanContext
}

type mgmtResponse struct {
	Endpoint   string
	StatusCode uint32
	Body       io.ReadCloser
}

type mgmtProvider interface {
	executeMgmtRequest(ctx context.Context, req mgmtRequest) (*mgmtResponse, error)
}

type mgmtProviderCore struct {
	provider             httpProvider
	mgmtTimeout          time.Duration
	retryStrategyWrapper *coreRetryStrategyWrapper
}

func (mpc *mgmtProviderCore) executeMgmtRequest(ctx context.Context, req mgmtRequest) (mgmtRespOut *mgmtResponse, errOut error) {
	timeout := req.Timeout
	if timeout == 0 {
		timeout = mpc.mgmtTimeout
	}

	retryStrategy := mpc.retryStrategyWrapper
	if req.RetryStrategy != nil {
		retryStrategy = newCoreRetryStrategyWrapper(req.RetryStrategy)
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
		TraceContext:  req.parentSpanCtx,
		Endpoint:      req.Endpoint,
	}

	coreresp, err := mpc.provider.DoHTTPRequest(ctx, corereq)
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

func checkForRateLimitError(statusCode uint32, errMsg string) error {
	if statusCode != 429 {
		return nil
	}

	errMsg = strings.ToLower(errMsg)
	var err error
	if strings.Contains(errMsg, "limit(s) exceeded") {
		err = ErrRateLimitedFailure
	} else if strings.Contains(errMsg, "maximum number of collections has been reached for scope") {
		err = ErrQuotaLimitedFailure
	}

	return err
}
