// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	initializationRetryDelay    = 30 * time.Second
	initializationRetryMaxDelay = 5 * time.Minute
	initializationRetryTimeout  = 30 * time.Second
	initializationRetryLimit    = 3
)

// Bound retry traffic across all PKI mounts in this process, including the
// storage and seal-wrap operations needed to inspect or migrate an issuer.
// Periodic calls never wait for a permit or create a retry goroutine.
var initializationRetrySlots = make(chan struct{}, 2)

// A storage error during Initialize must not permanently leave a backend using
// its constructor's default storage version. Recovery is armed by an actual
// initialization failure, not by legacy mode: a replication secondary may
// legitimately continue using legacy storage until its primary is upgraded.
type initializationState struct {
	sync.Mutex
	pending atomic.Bool

	request   logical.InitializationRequest
	retryAt   time.Time
	delay     time.Duration
	attempts  int
	retryable bool

	shutdownContext context.Context
	shutdown        context.CancelFunc
}

func (b *backend) initialize(ctx context.Context, request *logical.InitializationRequest) error {
	state := &b.initialization
	state.Lock()
	defer state.Unlock()

	if err := state.shutdownContext.Err(); err != nil {
		return err
	}
	state.request = *request
	state.delay = 0
	state.attempts = 0
	state.pending.Store(true)

	// Cleanup cancels in-flight initialization even when the caller supplied a
	// context which is not itself canceled by unmount or plugin reload.
	initContext, cancel := context.WithCancel(ctx)
	stop := context.AfterFunc(state.shutdownContext, cancel)
	defer stop()
	defer cancel()
	retryable, err := b.initializePrerequisites(initContext)
	if err == nil {
		err = initContext.Err()
	}
	if err != nil {
		b.finishInitialization(err, retryable)
		return err
	}
	err = b.initializeRemaining(ctx, request)
	b.finishInitialization(err, false)
	return err
}

// retryInitialization is called by the existing periodic lifecycle. It does no
// work after a successful initialization and never overlaps another attempt.
// Failures use a finite budget and jittered backoff. Only prerequisites that
// run before workers and ancillary initialization are retried. Each attempt is
// bounded by the current periodic context, its own deadline, and cleanup.
func (b *backend) retryInitialization(ctx context.Context) (bool, error) {
	state := &b.initialization
	if !state.pending.Load() {
		return true, nil
	}
	if !state.TryLock() {
		return false, nil
	}
	defer state.Unlock()
	if !state.pending.Load() {
		return true, nil
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if err := state.shutdownContext.Err(); err != nil {
		return false, err
	}
	if !state.retryable || state.attempts >= initializationRetryLimit || time.Now().Before(state.retryAt) {
		return false, nil
	}

	select {
	case initializationRetrySlots <- struct{}{}:
	default:
		return false, nil
	}
	// Count an attempt before starting I/O, including failed or canceled attempts.
	state.attempts++
	retryable, err := func() (bool, error) {
		defer func() { <-initializationRetrySlots }()
		retryContext, cancel := context.WithTimeout(ctx, initializationRetryTimeout)
		stop := context.AfterFunc(state.shutdownContext, cancel)
		defer stop()
		defer cancel()
		retryable, err := b.initializePrerequisites(retryContext)
		if err == nil {
			err = retryContext.Err()
		}
		return retryable, err
	}()
	if err != nil {
		b.finishInitialization(err, retryable)
		return false, fmt.Errorf("retrying PKI initialization: %w", err)
	}

	// The remaining startup runs once, with the caller's normal context rather
	// than the short prerequisite deadline. Certificate-count errors remain
	// intentionally nonfatal, and an enterprise startup error is not retried.
	err = b.initializeRemaining(ctx, &state.request)
	b.finishInitialization(err, false)
	if err != nil {
		return false, fmt.Errorf("retrying PKI initialization: %w", err)
	}
	b.Logger().Info("PKI backend initialization recovered")
	return true, nil
}

// finishInitialization is called with initializationState locked.
func (b *backend) finishInitialization(err error, retryable bool) {
	state := &b.initialization
	state.retryable = retryable
	if err == nil {
		state.delay = 0
		state.retryAt = time.Time{}
		state.pending.Store(false)
		return
	}
	state.pending.Store(true)
	if !retryable || state.attempts >= initializationRetryLimit {
		b.Logger().Error("PKI backend initialization is incomplete; automatic recovery stopped, inspect the error before reloading the mount to retry", "error", err, "retry_attempts", state.attempts)
		return
	}
	if state.delay == 0 {
		state.delay = initializationRetryDelay
	} else {
		state.delay = min(2*state.delay, initializationRetryMaxDelay)
	}
	// Add up to 50% jitter so mounts failing together do not retry in lockstep.
	delay := min(state.delay+time.Duration(rand.Int64N(int64(state.delay/2)+1)), initializationRetryMaxDelay)
	state.retryAt = time.Now().Add(delay)
	b.Logger().Warn("PKI backend initialization is incomplete", "error", err, "retry_delay", delay, "retry_attempts", state.attempts)
}

// Background startup work must preserve namespace and other context values
// without inheriting a request or prerequisite deadline. Cleanup owns its
// lifetime, including on unmount, reload, or active-node teardown.
func (b *backend) initializationLifetimeContext(ctx context.Context) context.Context {
	lifetime, cancel := context.WithCancel(context.WithoutCancel(ctx))
	context.AfterFunc(b.initialization.shutdownContext, cancel)
	return lifetime
}

// Requests must not read the legacy CA representation after initialization
// failed. In particular, raw CA endpoints would otherwise hide the error in an
// empty 204 response. Rollbacks must still reach PeriodicFunc to allow recovery.
func (b *backend) HandleRequest(ctx context.Context, request *logical.Request) (*logical.Response, error) {
	if b.initialization.pending.Load() && request.Operation != logical.RollbackOperation {
		state := b.System().ReplicationState()
		// Issuance also decides whether to forward dynamically, after reading
		// the role. An unready performance standby cannot safely make that
		// decision, so forward all its requests to the cluster's active node.
		if state.HasState(consts.ReplicationPerformanceStandby) {
			return nil, logical.ErrReadOnly
		}
		// Preserve explicit forwarding to the primary for non-local secondary
		// mounts before gating operations which must be handled locally.
		if path := b.Route(request.Path); path != nil {
			if operation := path.Operations[request.Operation]; operation != nil {
				properties := operation.Properties()
				if properties.ForwardPerformanceSecondary && !b.System().LocalMount() && state.HasState(consts.ReplicationPerformanceSecondary) {
					return nil, logical.ErrReadOnly
				}
			}
		}
		return nil, logical.CodedError(http.StatusServiceUnavailable, "PKI backend initialization has not completed; check server logs and reload the mount if automatic recovery has stopped")
	}
	return b.Backend.HandleRequest(ctx, request)
}
