package gocb

import "github.com/couchbase/gocbcore/v10"

func translateCoreRetryReasons(reasons []gocbcore.RetryReason) []RetryReason {
	var reasonsOut []RetryReason

	for _, retryReason := range reasons {
		gocbReason, ok := retryReason.(RetryReason)
		if !ok {
			logErrorf("Failed to assert gocbcore retry reason to gocb retry reason: %v", retryReason)
			continue
		}
		reasonsOut = append(reasonsOut, gocbReason)
	}

	return reasonsOut
}

type wrappedCoreRetryRequest struct {
	req gocbcore.RetryRequest
}

func (req *wrappedCoreRetryRequest) RetryAttempts() uint32 {
	return req.req.RetryAttempts()
}

func (req *wrappedCoreRetryRequest) Identifier() string {
	return req.req.Identifier()
}

func (req *wrappedCoreRetryRequest) Idempotent() bool {
	return req.req.Idempotent()
}

func (req *wrappedCoreRetryRequest) RetryReasons() []RetryReason {
	return translateCoreRetryReasons(req.req.RetryReasons())
}

func newCoreRetryStrategyWrapper(strategy RetryStrategy) *coreRetryStrategyWrapper {
	return &coreRetryStrategyWrapper{
		wrapped: strategy,
	}
}

type coreRetryStrategyWrapper struct {
	wrapped RetryStrategy
}

// RetryAfter calculates and returns a RetryAction describing how long to wait before retrying an operation.
func (rs *coreRetryStrategyWrapper) RetryAfter(req gocbcore.RetryRequest, reason gocbcore.RetryReason) gocbcore.RetryAction {
	wreq := &wrappedCoreRetryRequest{
		req: req,
	}
	wrappedAction := rs.wrapped.RetryAfter(wreq, RetryReason(reason))
	return gocbcore.RetryAction(wrappedAction)
}
