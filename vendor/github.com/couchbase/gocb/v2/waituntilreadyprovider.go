package gocb

import (
	"context"
	"errors"
	"time"

	"github.com/couchbase/gocbcoreps"

	"github.com/couchbase/gocbcore/v10"
)

type waitUntilReadyProvider interface {
	WaitUntilReady(ctx context.Context, deadline time.Time, opts *WaitUntilReadyOptions) error
}

type gocbcoreWaitUntilReadyProvider interface {
	WaitUntilReady(deadline time.Time, opts gocbcore.WaitUntilReadyOptions,
		cb gocbcore.WaitUntilReadyCallback) (gocbcore.PendingOp, error)
}

type waitUntilReadyProviderCore struct {
	retryStrategyWrapper *coreRetryStrategyWrapper
	provider             gocbcoreWaitUntilReadyProvider
}

func (wpw *waitUntilReadyProviderCore) WaitUntilReady(ctx context.Context, deadline time.Time,
	opts *WaitUntilReadyOptions) error {
	desiredState := opts.DesiredState
	if desiredState == 0 {
		desiredState = ClusterStateOnline
	}

	gocbcoreServices := make([]gocbcore.ServiceType, len(opts.ServiceTypes))
	for i, svc := range opts.ServiceTypes {
		gocbcoreServices[i] = gocbcore.ServiceType(svc)
	}

	wrapper := wpw.retryStrategyWrapper
	if opts.RetryStrategy != nil {
		wrapper = newCoreRetryStrategyWrapper(opts.RetryStrategy)
	}

	coreOpts := gocbcore.WaitUntilReadyOptions{
		DesiredState:  gocbcore.ClusterState(desiredState),
		ServiceTypes:  gocbcoreServices,
		RetryStrategy: wrapper,
	}

	var errOut error
	opm := newAsyncOpManager(ctx)
	err := opm.Wait(wpw.provider.WaitUntilReady(deadline, coreOpts, func(res *gocbcore.WaitUntilReadyResult, err error) {
		if err != nil {
			errOut = maybeEnhanceCoreErr(err)
			opm.Reject()
			return
		}

		opm.Resolve()
	}))
	if err != nil {
		errOut = maybeEnhanceCoreErr(err)
	}

	return errOut
}

type waitUntilreadyRequestPs struct {
	// reasons is effectively a set, so we can't just use len(reasons) for num attempts.
	reasons  []RetryReason
	attempts uint32

	strategy RetryStrategy
}

func (w *waitUntilreadyRequestPs) RetryAttempts() uint32 {
	return w.attempts
}

func (w *waitUntilreadyRequestPs) Identifier() string {
	return "WaitUntilReady"
}

func (w *waitUntilreadyRequestPs) Idempotent() bool {
	return true
}

func (w *waitUntilreadyRequestPs) RetryReasons() []RetryReason {
	return w.reasons
}

func (w *waitUntilreadyRequestPs) retryStrategy() RetryStrategy {
	return w.strategy
}

func (w *waitUntilreadyRequestPs) recordRetryAttempt(reason RetryReason) {
	w.attempts++
	found := false
	for i := 0; i < len(w.reasons); i++ {
		if w.reasons[i] == reason {
			found = true
			break
		}
	}

	// if idx is out of the range of retryReasons then it wasn't found.
	if !found {
		w.reasons = append(w.reasons, reason)
	}
}

type waitUntilReadyProviderPs struct {
	defaultRetryStrategy RetryStrategy
	client               *gocbcoreps.RoutingClient
}

func (wpw *waitUntilReadyProviderPs) WaitUntilReady(ctx context.Context, deadline time.Time,
	opts *WaitUntilReadyOptions) error {
	start := time.Now()
	desiredState := opts.DesiredState
	if desiredState == ClusterStateOffline {
		return makeInvalidArgumentsError("cannot use offline as a desired state")
	}
	if desiredState == 0 {
		desiredState = ClusterStateOnline
	}

	retryStrategy := wpw.defaultRetryStrategy
	if opts.RetryStrategy != nil {
		retryStrategy = opts.RetryStrategy
	}

	retryRequest := &waitUntilreadyRequestPs{
		strategy: retryStrategy,
	}
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	for {
		state := wpw.client.ConnectionState()
		var gocbState ClusterState
		switch state {
		case gocbcoreps.ConnStateOffline:
			gocbState = ClusterStateOffline
		case gocbcoreps.ConnStateOnline:
			gocbState = ClusterStateOnline
		case gocbcoreps.ConnStateDegraded:
			gocbState = ClusterStateDegraded
		}

		if gocbState == desiredState {
			return nil
		}

		shouldRetry, retryAfter := retryOrchMaybeRetry(retryRequest, NotReadyRetryReason)
		if !shouldRetry {
			// This should never actually happen - not ready is always retry.
			return ErrRequestCanceled
		}

		select {
		case <-ctx.Done():
			err := ctx.Err()
			if errors.Is(err, context.DeadlineExceeded) {
				return &TimeoutError{
					InnerError:    ErrUnambiguousTimeout,
					TimeObserved:  time.Since(start),
					RetryReasons:  retryRequest.RetryReasons(),
					RetryAttempts: retryRequest.RetryAttempts(),
				}
			}
		case <-time.After(time.Until(retryAfter)):
		}
	}
}
