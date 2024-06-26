// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/gammazero/workerpool"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	RollbackDefaultNumWorkers = 256
	RollbackWorkersEnvVar     = "VAULT_ROLLBACK_WORKERS"
)

var rollbackCanceled = errors.New("rollback attempt canceled")

// RollbackManager is responsible for performing rollbacks of partial
// secrets within logical backends.
//
// During normal operations, it is possible for logical backends to
// error partially through an operation. These are called "partial secrets":
// they are never sent back to a user, but they do need to be cleaned up.
// This manager handles that by periodically (on a timer) requesting that the
// backends clean up.
//
// The RollbackManager periodically initiates a logical.RollbackOperation
// on every mounted logical backend. It ensures that only one rollback operation
// is in-flight at any given time within a single seal/unseal phase.
type RollbackManager struct {
	logger log.Logger

	// This gives the current mount table of both logical and credential backends,
	// plus a RWMutex that is locked for reading. It is up to the caller to RUnlock
	// it when done with the mount table.
	backends func() []*MountEntry

	router *Router
	period time.Duration

	rollbackMetricsMountName bool
	inflightAll              sync.WaitGroup
	inflight                 map[string]*rollbackState
	inflightLock             sync.RWMutex

	doneCh          chan struct{}
	shutdown        bool
	shutdownCh      chan struct{}
	shutdownLock    sync.Mutex
	stopTicker      chan struct{}
	tickerIsStopped bool
	quitContext     context.Context
	runner          *workerpool.WorkerPool
	core            *Core
	// This channel is used for testing
	rollbacksDoneCh chan struct{}
}

// rollbackState is used to track the state of a single rollback attempt
type rollbackState struct {
	lastError error
	sync.WaitGroup
	cancelLockGrabCtx       context.Context
	cancelLockGrabCtxCancel context.CancelFunc
	// scheduled is the time that this job was created and submitted to the
	// rollbackRunner
	scheduled  time.Time
	isRunning  chan struct{}
	isCanceled chan struct{}
}

// NewRollbackManager is used to create a new rollback manager
func NewRollbackManager(ctx context.Context, logger log.Logger, backendsFunc func() []*MountEntry, router *Router, core *Core) *RollbackManager {
	r := &RollbackManager{
		logger:                   logger,
		backends:                 backendsFunc,
		router:                   router,
		period:                   core.rollbackPeriod,
		inflight:                 make(map[string]*rollbackState),
		doneCh:                   make(chan struct{}),
		shutdownCh:               make(chan struct{}),
		stopTicker:               make(chan struct{}),
		quitContext:              ctx,
		core:                     core,
		rollbackMetricsMountName: core.rollbackMountPathMetrics,
		rollbacksDoneCh:          make(chan struct{}),
	}
	numWorkers := r.numRollbackWorkers()
	r.logger.Info(fmt.Sprintf("Starting the rollback manager with %d workers", numWorkers))
	r.runner = workerpool.New(numWorkers)
	return r
}

func (m *RollbackManager) numRollbackWorkers() int {
	numWorkers := m.core.numRollbackWorkers
	envOverride := os.Getenv(RollbackWorkersEnvVar)
	if envOverride != "" {
		envVarWorkers, err := strconv.Atoi(envOverride)
		if err != nil || envVarWorkers < 1 {
			m.logger.Warn(fmt.Sprintf("%s must be a positive integer, but was %s", RollbackWorkersEnvVar, envOverride))
		} else {
			numWorkers = envVarWorkers
		}
	}
	return numWorkers
}

// Start starts the rollback manager
func (m *RollbackManager) Start() {
	go m.run()
}

// Stop stops the running manager. This will wait for any in-flight
// rollbacks to complete.
func (m *RollbackManager) Stop() {
	m.shutdownLock.Lock()
	defer m.shutdownLock.Unlock()
	if !m.shutdown {
		m.shutdown = true
		close(m.shutdownCh)
		<-m.doneCh
	}
	m.runner.StopWait()
}

// StopTicker stops the automatic Rollback manager's ticker, causing us
// to not do automatic rollbacks. This is useful for testing plugin's
// periodic function's behavior, without trying to race against the
// rollback manager proper.
//
// THIS SHOULD ONLY BE CALLED FROM TEST HELPERS.
func (m *RollbackManager) StopTicker() {
	if !m.tickerIsStopped {
		close(m.stopTicker)
		m.tickerIsStopped = true
	}
}

// run is a long running routine to periodically invoke rollback
func (m *RollbackManager) run() {
	m.logger.Info("starting rollback manager")
	tick := time.NewTicker(m.period)
	logTestStopOnce := false
	defer tick.Stop()
	defer close(m.doneCh)
	for {
		select {
		case <-tick.C:
			m.triggerRollbacks()
		case <-m.shutdownCh:
			m.logger.Info("stopping rollback manager")
			return

		case <-m.stopTicker:
			if !logTestStopOnce {
				m.logger.Info("stopping rollback manager ticker for tests")
				logTestStopOnce = true
			}
			tick.Stop()
		}
	}
}

// triggerRollbacks is used to trigger the rollbacks across all the backends
func (m *RollbackManager) triggerRollbacks() {
	backends := m.backends()

	for _, e := range backends {
		path := e.Path
		if e.Table == credentialTableType {
			path = credentialRoutePrefix + path
		}

		// When the mount is filtered, the backend will be nil
		ctx := namespace.ContextWithNamespace(m.quitContext, e.namespace)
		backend := m.router.MatchingBackend(ctx, path)
		if backend == nil {
			continue
		}
		fullPath := e.namespace.Path + path

		// Start a rollback if necessary
		m.startOrLookupRollback(ctx, fullPath, true)
	}
}

// lookupRollbackLocked checks if there's an inflight rollback with the given
// path. Callers must have the inflightLock. The function also reports metrics,
// since it is regularly called as part of the scheduled rollbacks.
func (m *RollbackManager) lookupRollbackLocked(fullPath string) *rollbackState {
	defer metrics.SetGauge([]string{"rollback", "queued"}, float32(m.runner.WaitingQueueSize()))
	defer metrics.SetGauge([]string{"rollback", "inflight"}, float32(len(m.inflight)))
	rsInflight := m.inflight[fullPath]
	return rsInflight
}

// newRollbackLocked creates a new rollback state and adds it to the inflight
// rollback map. Callers must have the inflightLock
func (m *RollbackManager) newRollbackLocked(fullPath string) *rollbackState {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	rs := &rollbackState{
		cancelLockGrabCtx:       cancelCtx,
		cancelLockGrabCtxCancel: cancelFunc,
		isRunning:               make(chan struct{}),
		isCanceled:              make(chan struct{}),
		scheduled:               time.Now(),
	}
	m.inflight[fullPath] = rs
	rs.Add(1)
	m.inflightAll.Add(1)
	return rs
}

// startOrLookupRollback is used to start an async rollback attempt.
func (m *RollbackManager) startOrLookupRollback(ctx context.Context, fullPath string, grabStatelock bool) *rollbackState {
	m.inflightLock.Lock()
	defer m.inflightLock.Unlock()
	rs := m.lookupRollbackLocked(fullPath)
	if rs != nil {
		return rs
	}

	// If no inflight rollback is already running, kick one off
	rs = m.newRollbackLocked(fullPath)

	select {
	case <-m.doneCh:
		// if we've already shut down, then don't submit the task to avoid a panic
		// we should still call finishRollback for the rollback state in order to remove
		// it from the map and decrement the waitgroup.

		// we already have the inflight lock, so we can't grab it here
		m.finishRollback(rs, errors.New("rollback manager is stopped"), fullPath, false)
	default:
		m.runner.Submit(func() {
			m.attemptRollback(ctx, fullPath, rs, grabStatelock)
			select {
			case m.rollbacksDoneCh <- struct{}{}:
			default:
			}
		})

	}
	return rs
}

func (m *RollbackManager) finishRollback(rs *rollbackState, err error, fullPath string, grabInflightLock bool) {
	rs.lastError = err
	rs.Done()
	m.inflightAll.Done()
	if grabInflightLock {
		m.inflightLock.Lock()
		defer m.inflightLock.Unlock()
	}
	if _, ok := m.inflight[fullPath]; ok {
		delete(m.inflight, fullPath)
	}
}

// attemptRollback invokes a RollbackOperation for the given path
func (m *RollbackManager) attemptRollback(ctx context.Context, fullPath string, rs *rollbackState, grabStatelock bool) (err error) {
	close(rs.isRunning)
	defer m.finishRollback(rs, err, fullPath, true)
	select {
	case <-rs.isCanceled:
		return rollbackCanceled
	default:
	}

	metrics.MeasureSince([]string{"rollback", "waiting"}, rs.scheduled)
	metricName := []string{"rollback", "attempt"}
	if m.rollbackMetricsMountName {
		metricName = append(metricName, strings.ReplaceAll(fullPath, "/", "-"))
	}
	defer metrics.MeasureSince(metricName, time.Now())

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		m.logger.Error("rollback failed to derive namespace from context", "path", fullPath)
		return err
	}
	if ns == nil {
		m.logger.Error("rollback found no namespace", "path", fullPath)
		return namespace.ErrNoNamespace
	}

	// Invoke a RollbackOperation
	req := &logical.Request{
		Operation: logical.RollbackOperation,
		Path:      ns.TrimmedPath(fullPath),
	}

	releaseLock := true
	if grabStatelock {
		doneCh := make(chan struct{})
		defer close(doneCh)

		stopCh := make(chan struct{})
		go func() {
			defer close(stopCh)

			select {
			case <-m.shutdownCh:
			case <-rs.cancelLockGrabCtx.Done():
			case <-doneCh:
			case <-rs.isCanceled:
			}
		}()

		// Grab the statelock or stop
		l := newLockGrabber(m.core.stateLock.RLock, m.core.stateLock.RUnlock, stopCh)
		go l.grab()
		if stopped := l.lockOrStop(); stopped {
			// If we stopped due to shutdown, return. Otherwise another thread
			// is holding the lock for us, continue on.
			select {
			case <-m.shutdownCh:
				return errors.New("rollback shutting down")
			case <-rs.isCanceled:
				return rollbackCanceled
			default:
				releaseLock = false
			}
		}
	}

	var cancelFunc context.CancelFunc
	ctx, cancelFunc = context.WithTimeout(ctx, DefaultMaxRequestDuration)
	resp, err := m.router.Route(ctx, req)
	if grabStatelock && releaseLock {
		m.core.stateLock.RUnlock()
	}
	cancelFunc()

	// If the error is an unsupported operation, then it doesn't
	// matter, the backend doesn't support it.
	if err == logical.ErrUnsupportedOperation {
		err = nil
	}
	// If we failed due to read-only storage, we can't do anything; ignore
	if (err != nil && strings.Contains(err.Error(), logical.ErrReadOnly.Error())) ||
		(resp.IsError() && strings.Contains(resp.Error().Error(), logical.ErrReadOnly.Error())) {
		err = nil
	}
	if err != nil {
		m.logger.Error("error rolling back", "path", fullPath, "error", err)
	}
	return
}

// Rollback is used to trigger an immediate rollback of the path,
// or to join an existing rollback operation if in flight. Caller should have
// core's statelock held (write OR read). If an already inflight rollback is
// happening this function will simply wait for it to complete
func (m *RollbackManager) Rollback(ctx context.Context, path string) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	fullPath := ns.Path + path

	m.inflightLock.Lock()
	rs := m.lookupRollbackLocked(fullPath)
	if rs != nil {
		// Since we have the statelock held, tell any inflight rollback to give up
		// trying to acquire it. This will prevent deadlocks in the case where we
		// have the write lock. In the case where it was waiting to grab
		// a read lock it will then simply continue with the rollback
		// operation under the protection of our write lock.
		rs.cancelLockGrabCtxCancel()

		select {
		case <-rs.isRunning:
			// if the rollback has started then we should wait for it to complete
			m.inflightLock.Unlock()
			rs.Wait()
			return rs.lastError
		default:
		}

		// if the rollback hasn't started and there's no capacity, we could
		// end up deadlocking. Cancel the existing rollback and start a new
		// one.
		close(rs.isCanceled)
	}
	rs = m.newRollbackLocked(fullPath)
	m.inflightLock.Unlock()

	// we can ignore the error, since it's going to be set in rs.lastError
	m.attemptRollback(ctx, fullPath, rs, false)

	rs.Wait()
	return rs.lastError
}

// The methods below are the hooks from core that are called pre/post seal.

// startRollback is used to start the rollback manager after unsealing
func (c *Core) startRollback() error {
	backendsFunc := func() []*MountEntry {
		ret := []*MountEntry{}
		c.mountsLock.RLock()
		defer c.mountsLock.RUnlock()
		// During teardown/setup after a leader change or unseal there could be
		// something racy here so make sure the table isn't nil
		if c.mounts != nil {
			for _, entry := range c.mounts.Entries {
				ret = append(ret, entry)
			}
		}
		c.authLock.RLock()
		defer c.authLock.RUnlock()
		// During teardown/setup after a leader change or unseal there could be
		// something racy here so make sure the table isn't nil
		if c.auth != nil {
			for _, entry := range c.auth.Entries {
				ret = append(ret, entry)
			}
		}
		return ret
	}
	rollbackLogger := c.baseLogger.Named("rollback")
	c.AddLogger(rollbackLogger)
	c.rollback = NewRollbackManager(c.activeContext, rollbackLogger, backendsFunc, c.router, c)
	c.rollback.Start()
	return nil
}

// stopRollback is used to stop running the rollback manager before sealing
func (c *Core) stopRollback() error {
	if c.rollback != nil {
		c.rollback.Stop()
		c.rollback = nil
	}
	return nil
}
