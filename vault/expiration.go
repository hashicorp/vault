package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/fairshare"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/base62"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/quotas"
	uberAtomic "go.uber.org/atomic"
)

const (
	// expirationSubPath is the sub-path used for the expiration manager
	// view. This is nested under the system view.
	expirationSubPath = "expire/"

	// leaseViewPrefix is the prefix used for the ID based lookup of leases.
	leaseViewPrefix = "id/"

	// tokenViewPrefix is the prefix used for the token based lookup of leases.
	tokenViewPrefix = "token/"

	// maxRevokeAttempts limits how many revoke attempts are made
	maxRevokeAttempts = 6

	// revokeRetryBase is a baseline retry time
	revokeRetryBase = 10 * time.Second

	// maxLeaseDuration is the default maximum lease duration
	maxLeaseTTL = 32 * 24 * time.Hour

	// defaultLeaseDuration is the default lease duration used when no lease is specified
	defaultLeaseTTL = maxLeaseTTL

	// maxLeaseThreshold is the maximum lease count before generating log warning
	maxLeaseThreshold = 256000

	// numExpirationWorkersDefault is the maximum amount of workers working on lease expiration
	numExpirationWorkersDefault = 200

	// number of workers to use for general purpose testing
	numExpirationWorkersTest = 10

	fairshareWorkersOverrideVar = "VAULT_LEASE_REVOCATION_WORKERS"

	// limit irrevocable error messages to 240 characters to be respectful of
	// storage/memory
	maxIrrevocableErrorLength = 240

	genericIrrevocableErrorMessage = "unknown"

	outOfRetriesMessage = "out of retries"

	// maximum number of irrevocable leases we return to the irrevocable lease
	// list API **without** the `force` flag set
	MaxIrrevocableLeasesToReturn = 10000

	MaxIrrevocableLeasesWarning = "Command halted because many irrevocable leases were found. To emit the entire list, re-run the command with force set true."
)

type pendingInfo struct {
	// A subset of the lease entry, cached in memory
	cachedLeaseInfo  *leaseEntry
	timer            *time.Timer
	revokesAttempted uint8
}

// ExpirationManager is used by the Core to manage leases. Secrets
// can provide a lease, meaning that they can be renewed or revoked.
// If a secret is not renewed in timely manner, it may be expired, and
// the ExpirationManager will handle doing automatic revocation.
type ExpirationManager struct {
	core       *Core
	router     *Router
	idView     *BarrierView
	tokenView  *BarrierView
	tokenStore *TokenStore
	logger     log.Logger

	// Although the data structure itself is atomic,
	// pendingLock should be held to ensure lease modifications
	// are atomic (with respect to storage, expiration time,
	// and particularly the lease count.)
	// The nonexpiring map holds entries for root tokens with
	// TTL zero, which we want to count but have no timer associated.
	pending     sync.Map
	nonexpiring sync.Map
	leaseCount  int
	pendingLock sync.RWMutex

	// A sync.Lock for every active leaseID
	lockPerLease sync.Map
	// Track expired leases that have been determined to be irrevocable (without
	// manual intervention). We retain a subset of the lease info in memory
	irrevocable sync.Map

	// Track count for metrics reporting
	// This value is protected by pendingLock
	irrevocableLeaseCount int

	// The uniquePolicies map holds policy sets, so they can
	// be deduplicated. It is periodically emptied to prevent
	// unbounded growth.
	uniquePolicies      map[string][]string
	emptyUniquePolicies *time.Ticker

	tidyLock *int32

	restoreMode        *int32
	restoreModeLock    sync.RWMutex
	restoreRequestLock sync.RWMutex
	restoreLocks       []*locksutil.LockEntry
	restoreLoaded      sync.Map
	quitCh             chan struct{}

	// do not hold coreStateLock in any API handler code - it is already held
	coreStateLock     *DeadlockRWMutex
	quitContext       context.Context
	leaseCheckCounter *uint32

	logLeaseExpirations bool
	expireFunc          ExpireLeaseStrategy

	revokePermitPool *physical.PermitPool

	// testRegisterAuthFailure, if set to true, triggers an explicit failure on
	// RegisterAuth to simulate a partial failure during a token creation
	// request. This value should only be set by tests.
	testRegisterAuthFailure uberAtomic.Bool

	jobManager *fairshare.JobManager
}

type ExpireLeaseStrategy func(context.Context, *ExpirationManager, string, *namespace.Namespace)

// revocationJob should only be created through newRevocationJob()
type revocationJob struct {
	leaseID   string
	ns        *namespace.Namespace
	m         *ExpirationManager
	nsCtx     context.Context
	startTime time.Time
}

func newRevocationJob(nsCtx context.Context, leaseID string, ns *namespace.Namespace, m *ExpirationManager) (*revocationJob, error) {
	if leaseID == "" {
		return nil, fmt.Errorf("cannot have empty lease id")
	}
	if m == nil {
		return nil, fmt.Errorf("cannot have nil expiration manager")
	}
	if nsCtx == nil {
		return nil, fmt.Errorf("cannot have nil namespace context.Context")
	}

	return &revocationJob{
		leaseID:   leaseID,
		ns:        ns,
		m:         m,
		nsCtx:     nsCtx,
		startTime: time.Now(),
	}, nil
}

// errIsUnrecoverable returns true if the logical error is unlikely to resolve
// automatically or with additional retries
func errIsUnrecoverable(err error) bool {
	switch {
	case errors.Is(err, logical.ErrUnrecoverable),
		errors.Is(err, logical.ErrUnsupportedOperation),
		errors.Is(err, logical.ErrUnsupportedPath),
		errors.Is(err, logical.ErrInvalidRequest):
		return true
	}

	return false
}

func (r *revocationJob) Execute() error {
	r.m.core.metricSink.IncrCounterWithLabels([]string{"expire", "lease_expiration"}, 1, []metrics.Label{metricsutil.NamespaceLabel(r.ns)})
	r.m.core.metricSink.MeasureSinceWithLabels([]string{"expire", "lease_expiration", "time_in_queue"}, r.startTime, []metrics.Label{metricsutil.NamespaceLabel(r.ns)})

	// don't start the timer until the revocation is being executed
	revokeCtx, cancel := context.WithTimeout(r.nsCtx, DefaultMaxRequestDuration)
	defer cancel()

	go func() {
		select {
		case <-r.m.quitCh:
			cancel()
		case <-revokeCtx.Done():
		}
	}()

	select {
	case <-r.m.quitCh:
		r.m.logger.Error("shutting down, not attempting further revocation of lease", "lease_id", r.leaseID)
		return nil
	case <-r.m.quitContext.Done():
		r.m.logger.Error("core context canceled, not attempting further revocation of lease", "lease_id", r.leaseID)
		return nil
	default:
	}

	r.m.coreStateLock.RLock()
	err := r.m.Revoke(revokeCtx, r.leaseID)
	r.m.coreStateLock.RUnlock()

	return err
}

func (r *revocationJob) OnFailure(err error) {
	r.m.core.metricSink.IncrCounterWithLabels([]string{"expire", "lease_expiration", "error"}, 1, []metrics.Label{metricsutil.NamespaceLabel(r.ns)})
	r.m.logger.Error("failed to revoke lease", "lease_id", r.leaseID, "error", err)

	r.m.pendingLock.Lock()
	defer r.m.pendingLock.Unlock()
	pendingRaw, ok := r.m.pending.Load(r.leaseID)
	if !ok {
		r.m.logger.Warn("failed to find lease in pending map for revocation retry", "lease_id", r.leaseID)
		return
	}

	pending := pendingRaw.(pendingInfo)
	pending.revokesAttempted++
	if pending.revokesAttempted >= maxRevokeAttempts || errIsUnrecoverable(err) {
		r.m.logger.Trace("marking lease as irrevocable", "lease_id", r.leaseID, "error", err)
		if pending.revokesAttempted >= maxRevokeAttempts {
			r.m.logger.Trace("lease has consumed all retry attempts", "lease_id", r.leaseID)
			err = fmt.Errorf("%v: %w", outOfRetriesMessage, err)
		}

		le, loadErr := r.m.loadEntry(r.nsCtx, r.leaseID)
		if loadErr != nil {
			r.m.logger.Warn("failed to mark lease as irrevocable - failed to load", "lease_id", r.leaseID, "err", loadErr)
			return
		}
		if le == nil {
			r.m.logger.Warn("failed to mark lease as irrevocable - nil lease", "lease_id", r.leaseID)
			return
		}

		r.m.markLeaseIrrevocable(r.nsCtx, le, err)
		return
	}

	pending.timer.Reset(revokeExponentialBackoff(pending.revokesAttempted))
	r.m.pending.Store(r.leaseID, pending)
}

func expireLeaseStrategyFairsharing(ctx context.Context, m *ExpirationManager, leaseID string, ns *namespace.Namespace) {
	nsCtx := namespace.ContextWithNamespace(ctx, ns)

	mountAccessor := m.getLeaseMountAccessorLocked(ctx, leaseID)

	job, err := newRevocationJob(nsCtx, leaseID, ns, m)
	if err != nil {
		m.logger.Warn("error creating revocation job", "error", err)
		return
	}

	m.jobManager.AddJob(job, mountAccessor)
}

func revokeExponentialBackoff(attempt uint8) time.Duration {
	exp := (1 << attempt) * revokeRetryBase
	randomDelta := 0.5 * float64(exp)

	// Allow backoff time to be a random value between exp +/- (0.5*exp)
	backoffTime := (float64(exp) - randomDelta) + (rand.Float64() * (2 * randomDelta))
	return time.Duration(backoffTime)
}

// revokeIDFunc is invoked when a given ID is expired
func expireLeaseStrategyRevoke(ctx context.Context, m *ExpirationManager, leaseID string, ns *namespace.Namespace) {
	for attempt := uint(0); attempt < maxRevokeAttempts; attempt++ {
		releasePermit := func() {}
		if m.revokePermitPool != nil {
			m.logger.Trace("expiring lease; waiting for permit pool")
			m.revokePermitPool.Acquire()
			releasePermit = m.revokePermitPool.Release
			m.logger.Trace("expiring lease; got permit pool")
		}

		metrics.IncrCounterWithLabels([]string{"expire", "lease_expiration"}, 1, []metrics.Label{{"namespace", ns.ID}})

		revokeCtx, cancel := context.WithTimeout(ctx, DefaultMaxRequestDuration)
		revokeCtx = namespace.ContextWithNamespace(revokeCtx, ns)

		go func() {
			select {
			case <-ctx.Done():
			case <-m.quitCh:
				cancel()
			case <-revokeCtx.Done():
			}
		}()

		select {
		case <-m.quitCh:
			m.logger.Error("shutting down, not attempting further revocation of lease", "lease_id", leaseID)
			releasePermit()
			cancel()
			return
		case <-m.quitContext.Done():
			m.logger.Error("core context canceled, not attempting further revocation of lease", "lease_id", leaseID)
			releasePermit()
			cancel()
			return
		default:
		}

		m.coreStateLock.RLock()
		err := m.Revoke(revokeCtx, leaseID)
		m.coreStateLock.RUnlock()
		releasePermit()
		cancel()
		if err == nil {
			return
		}

		metrics.IncrCounterWithLabels([]string{"expire", "lease_expiration", "error"}, 1, []metrics.Label{{"namespace", ns.ID}})

		m.logger.Error("failed to revoke lease", "lease_id", leaseID, "error", err)
		time.Sleep((1 << attempt) * revokeRetryBase)
	}
	m.logger.Error("maximum revoke attempts reached", "lease_id", leaseID)
}

func getNumExpirationWorkers(c *Core, l log.Logger) int {
	numWorkers := c.numExpirationWorkers

	workerOverride := os.Getenv(fairshareWorkersOverrideVar)
	if workerOverride != "" {
		i, err := strconv.Atoi(workerOverride)
		if err != nil {
			l.Warn("vault lease revocation workers override must be an integer", "value", workerOverride)
		} else if i < 1 || i > 10000 {
			l.Warn("vault lease revocation workers override out of range", "value", i)
		} else {
			numWorkers = i
		}
	}

	return numWorkers
}

// NewExpirationManager creates a new ExpirationManager that is backed
// using a given view, and uses the provided router for revocation.
func NewExpirationManager(c *Core, view *BarrierView, e ExpireLeaseStrategy, logger log.Logger) *ExpirationManager {
	var permitPool *physical.PermitPool
	if os.Getenv("VAULT_16_REVOKE_PERMITPOOL") != "" {
		permitPoolSize := 50
		permitPoolSizeRaw, err := strconv.Atoi(os.Getenv("VAULT_16_REVOKE_PERMITPOOL"))
		if err == nil && permitPoolSizeRaw > 0 {
			permitPoolSize = permitPoolSizeRaw
		}

		permitPool = physical.NewPermitPool(permitPoolSize)

	}

	jobManager := fairshare.NewJobManager("expire", getNumExpirationWorkers(c, logger), logger.Named("job-manager"), c.metricSink)
	jobManager.Start()

	exp := &ExpirationManager{
		core:        c,
		router:      c.router,
		idView:      view.SubView(leaseViewPrefix),
		tokenView:   view.SubView(tokenViewPrefix),
		tokenStore:  c.tokenStore,
		logger:      logger,
		pending:     sync.Map{},
		nonexpiring: sync.Map{},
		leaseCount:  0,
		tidyLock:    new(int32),

		lockPerLease: sync.Map{},

		uniquePolicies:      make(map[string][]string),
		emptyUniquePolicies: time.NewTicker(7 * 24 * time.Hour),

		// new instances of the expiration manager will go immediately into
		// restore mode
		restoreMode:  new(int32),
		restoreLocks: locksutil.CreateLocks(),
		quitCh:       make(chan struct{}),

		coreStateLock:     &c.stateLock,
		quitContext:       c.activeContext,
		leaseCheckCounter: new(uint32),

		logLeaseExpirations: os.Getenv("VAULT_SKIP_LOGGING_LEASE_EXPIRATIONS") == "",
		expireFunc:          e,
		revokePermitPool:    permitPool,

		jobManager: jobManager,
	}
	*exp.restoreMode = 1

	if exp.logger == nil {
		opts := log.LoggerOptions{Name: "expiration_manager"}
		exp.logger = log.New(&opts)
	}

	go exp.uniquePoliciesGc()

	return exp
}

// setupExpiration is invoked after we've loaded the mount table to
// initialize the expiration manager
func (c *Core) setupExpiration(e ExpireLeaseStrategy) error {
	c.metricsMutex.Lock()
	defer c.metricsMutex.Unlock()
	// Create a sub-view
	view := c.systemBarrierView.SubView(expirationSubPath)

	// Create the manager
	expLogger := c.baseLogger.Named("expiration")
	c.AddLogger(expLogger)
	mgr := NewExpirationManager(c, view, e, expLogger)
	c.expiration = mgr

	// Link the token store to this
	c.tokenStore.SetExpirationManager(mgr)

	// Restore the existing state
	c.logger.Info("restoring leases")
	errorFunc := func() {
		c.logger.Error("shutting down")
		if err := c.Shutdown(); err != nil {
			c.logger.Error("error shutting down core", "error", err)
		}
	}
	go c.expiration.Restore(errorFunc)

	quit := c.expiration.quitCh
	go func() {
		t := time.NewTimer(24 * time.Hour)
		for {
			select {
			case <-quit:
				return
			case <-t.C:
				c.expiration.attemptIrrevocableLeasesRevoke()
				t.Reset(24 * time.Hour)
			}
		}
	}()

	return nil
}

// stopExpiration is used to stop the expiration manager before
// sealing the Vault.
func (c *Core) stopExpiration() error {
	if c.expiration != nil {
		if err := c.expiration.Stop(); err != nil {
			return err
		}
		c.metricsMutex.Lock()
		defer c.metricsMutex.Unlock()
		c.expiration = nil
	}
	return nil
}

// lockLease takes out a lock for a given lease ID
func (m *ExpirationManager) lockLease(leaseID string) {
	locksutil.LockForKey(m.restoreLocks, leaseID).Lock()
}

// unlockLease unlocks a given lease ID
func (m *ExpirationManager) unlockLease(leaseID string) {
	locksutil.LockForKey(m.restoreLocks, leaseID).Unlock()
}

// inRestoreMode returns if we are currently in restore mode
func (m *ExpirationManager) inRestoreMode() bool {
	return atomic.LoadInt32(m.restoreMode) == 1
}

func (m *ExpirationManager) invalidate(key string) {
	switch {
	case strings.HasPrefix(key, leaseViewPrefix):
		leaseID := strings.TrimPrefix(key, leaseViewPrefix)
		ctx := m.quitContext
		_, nsID := namespace.SplitIDFromString(leaseID)
		leaseNS := namespace.RootNamespace
		var err error
		if nsID != "" {
			leaseNS, err = NamespaceByID(ctx, nsID, m.core)
			if err != nil {
				m.logger.Error("failed to invalidate lease entry", "error", err)
				return
			}
		}

		le, err := m.loadEntryInternal(namespace.ContextWithNamespace(ctx, leaseNS), leaseID, false, false)
		if err != nil {
			m.logger.Error("failed to invalidate lease entry", "error", err)
			return
		}

		m.pendingLock.Lock()
		defer m.pendingLock.Unlock()
		info, ok := m.pending.Load(leaseID)
		switch {
		case ok:
			switch {
			case le == nil:
				// Handle lease deletion
				pending := info.(pendingInfo)
				pending.timer.Stop()
				m.pending.Delete(leaseID)
				m.leaseCount--

				if err := m.core.quotasHandleLeases(ctx, quotas.LeaseActionDeleted, []string{leaseID}); err != nil {
					m.logger.Error("failed to update quota on lease invalidation", "error", err)
					return
				}
			default:
				// Update the lease in memory
				m.updatePendingInternal(le)
			}
		default:
			if le == nil {
				// There is no entry in the pending map and the invalidation
				// resulted in a nil entry. Therefore we should clean up the
				// other maps, and update metrics/quotas if appropriate.
				m.nonexpiring.Delete(leaseID)

				if _, ok := m.irrevocable.Load(leaseID); ok {
					m.irrevocable.Delete(leaseID)
					m.irrevocableLeaseCount--

					m.leaseCount--
					if err := m.core.quotasHandleLeases(ctx, quotas.LeaseActionDeleted, []string{leaseID}); err != nil {
						m.logger.Error("failed to update quota on lease invalidation", "error", err)
						return
					}
				}
				return
			}
			// Handle lease update (if irrevocable) or creation (if pending)
			m.updatePendingInternal(le)
		}
	}
}

// Tidy cleans up the dangling storage entries for leases. It scans the storage
// view to find all the available leases, checks if the token embedded in it is
// either empty or invalid and in both the cases, it revokes them. It also uses
// a token cache to avoid multiple lookups of the same token ID. It is normally
// not required to use the API that invokes this. This is only intended to
// clean up the corrupt storage due to bugs.
func (m *ExpirationManager) Tidy(ctx context.Context) error {
	if m.inRestoreMode() {
		return errors.New("cannot run tidy while restoring leases")
	}

	var tidyErrors *multierror.Error

	logger := m.logger.Named("tidy")
	m.core.AddLogger(logger)

	if !atomic.CompareAndSwapInt32(m.tidyLock, 0, 1) {
		logger.Warn("tidy operation on leases is already in progress")
		return nil
	}

	defer atomic.CompareAndSwapInt32(m.tidyLock, 1, 0)

	logger.Info("beginning tidy operation on leases")
	defer logger.Info("finished tidy operation on leases")

	// Create a cache to keep track of looked up tokens
	tokenCache := make(map[string]bool)
	var countLease, revokedCount, deletedCountInvalidToken, deletedCountEmptyToken int64

	tidyFunc := func(leaseID string) {
		countLease++
		if countLease%500 == 0 {
			logger.Info("tidying leases", "progress", countLease)
		}

		le, err := m.loadEntry(ctx, leaseID)
		if err != nil {
			tidyErrors = multierror.Append(tidyErrors, fmt.Errorf("failed to load the lease ID %q: %w", leaseID, err))
			return
		}

		if le == nil {
			tidyErrors = multierror.Append(tidyErrors, fmt.Errorf("nil entry for lease ID %q: %w", leaseID, err))
			return
		}

		var isValid, ok bool
		revokeLease := false
		if le.ClientToken == "" {
			logger.Debug("revoking lease which has an empty token", "lease_id", leaseID)
			revokeLease = true
			deletedCountEmptyToken++
			goto REVOKE_CHECK
		}

		isValid, ok = tokenCache[le.ClientToken]
		if !ok {
			lock := locksutil.LockForKey(m.tokenStore.tokenLocks, le.ClientToken)
			lock.RLock()
			te, err := m.tokenStore.lookupInternal(ctx, le.ClientToken, false, true)
			lock.RUnlock()

			if err != nil {
				tidyErrors = multierror.Append(tidyErrors, fmt.Errorf("failed to lookup token: %w", err))
				return
			}

			if te == nil {
				logger.Debug("revoking lease which holds an invalid token", "lease_id", leaseID)
				revokeLease = true
				deletedCountInvalidToken++
				tokenCache[le.ClientToken] = false
			} else {
				tokenCache[le.ClientToken] = true
			}

			goto REVOKE_CHECK
		} else {
			if isValid {
				return
			}

			logger.Debug("revoking lease which contains an invalid token", "lease_id", leaseID)
			revokeLease = true
			deletedCountInvalidToken++
			goto REVOKE_CHECK
		}

	REVOKE_CHECK:
		if revokeLease {
			// Force the revocation and skip going through the token store
			// again

			leaseLock := m.lockForLeaseID(leaseID)
			leaseLock.Lock()
			err = m.revokeCommon(ctx, leaseID, true, true)
			leaseLock.Unlock()
			if err != nil {
				tidyErrors = multierror.Append(tidyErrors, fmt.Errorf("failed to revoke an invalid lease with ID %q: %w", leaseID, err))
				return
			}
			revokedCount++
		}
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	leaseView := m.leaseView(ns)
	if err := logical.ScanView(m.quitContext, leaseView, tidyFunc); err != nil {
		return err
	}

	logger.Info("number of leases scanned", "count", countLease)
	logger.Info("number of leases which had empty tokens", "count", deletedCountEmptyToken)
	logger.Info("number of leases which had invalid tokens", "count", deletedCountInvalidToken)
	logger.Info("number of leases successfully revoked", "count", revokedCount)

	return tidyErrors.ErrorOrNil()
}

// Restore is used to recover the lease states when starting.
// This is used after starting the vault.
func (m *ExpirationManager) Restore(errorFunc func()) (retErr error) {
	defer func() {
		// Turn off restore mode. We can do this safely without the lock because
		// if restore mode finished successfully, restore mode was already
		// disabled with the lock. In an error state, this will allow the
		// Stop() function to shut everything down.
		atomic.StoreInt32(m.restoreMode, 0)

		switch {
		case retErr == nil:
		case strings.Contains(retErr.Error(), context.Canceled.Error()):
			// Don't run error func because we lost leadership
			m.logger.Warn("context canceled while restoring leases, stopping lease loading")
			retErr = nil
		case errwrap.Contains(retErr, ErrBarrierSealed.Error()):
			// Don't run error func because we're likely already shutting down
			m.logger.Warn("barrier sealed while restoring leases, stopping lease loading")
			retErr = nil
		default:
			m.logger.Error("error restoring leases", "error", retErr)
			if errorFunc != nil {
				errorFunc()
			}
		}
	}()

	// Accumulate existing leases
	m.logger.Debug("collecting leases")
	existing, leaseCount, err := m.collectLeases()
	if err != nil {
		return err
	}
	m.logger.Debug("leases collected", "num_existing", leaseCount)

	// Make the channels used for the worker pool
	type lease struct {
		namespace *namespace.Namespace
		id        string
	}
	broker := make(chan *lease)
	quit := make(chan bool)
	// Buffer these channels to prevent deadlocks
	errs := make(chan error, len(existing))
	result := make(chan struct{}, len(existing))

	// Use a wait group
	wg := &sync.WaitGroup{}

	// Create 64 workers to distribute work to
	for i := 0; i < consts.ExpirationRestoreWorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case lease, ok := <-broker:
					// broker has been closed, we are done
					if !ok {
						return
					}

					ctx := namespace.ContextWithNamespace(m.quitContext, lease.namespace)
					err := m.processRestore(ctx, lease.id)
					if err != nil {
						errs <- err
						continue
					}

					// Send message that lease is done
					result <- struct{}{}

				// quit early
				case <-quit:
					return

				case <-m.quitCh:
					return
				}
			}
		}()
	}

	// Distribute the collected keys to the workers in a go routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for ns := range existing {
			for _, leaseID := range existing[ns] {
				i++
				if i%500 == 0 {
					m.logger.Debug("leases loading", "progress", i)
				}

				select {
				case <-quit:
					return

				case <-m.quitCh:
					return

				default:
					broker <- &lease{
						namespace: ns,
						id:        leaseID,
					}
				}
			}
		}

		// Close the broker, causing worker routines to exit
		close(broker)
	}()

	// Ensure all keys on the chan are processed
	for i := 0; i < leaseCount; i++ {
		select {
		case err := <-errs:
			// Close all go routines
			close(quit)
			return err

		case <-m.quitCh:
			close(quit)
			return nil

		case <-result:
		}
	}

	// Let all go routines finish
	wg.Wait()

	m.restoreModeLock.Lock()
	atomic.StoreInt32(m.restoreMode, 0)
	m.restoreLoaded.Range(func(k, v interface{}) bool {
		m.restoreLoaded.Delete(k)
		return true
	})
	m.restoreLocks = nil
	m.restoreModeLock.Unlock()

	m.logger.Info("lease restore complete")
	return nil
}

// processRestore takes a lease and restores it in the expiration manager if it has
// not already been seen
func (m *ExpirationManager) processRestore(ctx context.Context, leaseID string) error {
	m.restoreRequestLock.RLock()
	defer m.restoreRequestLock.RUnlock()

	// Check if the lease has been seen
	if _, ok := m.restoreLoaded.Load(leaseID); ok {
		return nil
	}

	m.lockLease(leaseID)
	defer m.unlockLease(leaseID)

	// Check again with the lease locked
	if _, ok := m.restoreLoaded.Load(leaseID); ok {
		return nil
	}

	// Load lease and restore expiration timer
	_, err := m.loadEntryInternal(ctx, leaseID, true, false)
	if err != nil {
		return err
	}

	return nil
}

// Stop is used to prevent further automatic revocations.
// This must be called before sealing the view.
func (m *ExpirationManager) Stop() error {
	// Stop all the pending expiration timers
	m.logger.Debug("stop triggered")
	defer m.logger.Debug("finished stopping")

	m.jobManager.Stop()

	// Do this before stopping pending timers to avoid potential races with
	// expiring timers
	close(m.quitCh)

	m.pendingLock.Lock()
	// Replacing the entire map would cause a race with
	// a simultaneous WalkTokens, which doesn't hold pendingLock.
	m.pending.Range(func(key, value interface{}) bool {
		info := value.(pendingInfo)
		info.timer.Stop()
		m.pending.Delete(key)
		return true
	})
	m.leaseCount = 0
	m.nonexpiring.Range(func(key, value interface{}) bool {
		m.nonexpiring.Delete(key)
		return true
	})
	m.uniquePolicies = make(map[string][]string)
	m.irrevocable.Range(func(key, _ interface{}) bool {
		m.irrevocable.Delete(key)
		return true
	})
	m.irrevocableLeaseCount = 0
	m.pendingLock.Unlock()

	if m.inRestoreMode() {
		for {
			if !m.inRestoreMode() {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}

	m.emptyUniquePolicies.Stop()

	return nil
}

// Revoke is used to revoke a secret named by the given LeaseID
func (m *ExpirationManager) Revoke(ctx context.Context, leaseID string) error {
	defer metrics.MeasureSince([]string{"expire", "revoke"}, time.Now())

	return m.revokeCommon(ctx, leaseID, false, false)
}

// LazyRevoke is used to queue revocation for a secret named by the given
// LeaseID. If the lease was not found it returns nil; if the lease was found
// it triggers a return of a 202.
func (m *ExpirationManager) LazyRevoke(ctx context.Context, leaseID string) error {
	defer metrics.MeasureSince([]string{"expire", "lazy-revoke"}, time.Now())
	return m.lazyRevokeInternal(ctx, leaseID)
}

// Mark a lease as expiring immediately
func (m *ExpirationManager) lazyRevokeInternal(ctx context.Context, leaseID string) error {
	leaseLock := m.lockForLeaseID(leaseID)
	leaseLock.Lock()
	defer leaseLock.Unlock()

	// Load the entry
	le, err := m.loadEntry(ctx, leaseID)
	if err != nil {
		return err
	}

	// If there is no entry, nothing to revoke
	if le == nil {
		return nil
	}

	le.ExpireTime = time.Now()
	if err := m.persistEntry(ctx, le); err != nil {
		return err
	}
	m.updatePending(le)

	return nil
}

// should be run on a schedule. something like once a day, maybe once a week
func (m *ExpirationManager) attemptIrrevocableLeasesRevoke() {
	m.irrevocable.Range(func(k, v interface{}) bool {
		leaseID := k.(string)
		le := v.(*leaseEntry)

		if le.ExpireTime.Add(time.Hour).Before(time.Now()) {
			// if we get an error (or no namespace) note it, but continue attempting
			// to revoke other leases
			leaseNS, err := m.getNamespaceFromLeaseID(m.core.activeContext, leaseID)
			if err != nil {
				m.logger.Debug("could not get lease namespace from ID", "error", err)
				return true
			}
			if leaseNS == nil {
				m.logger.Debug("could not get lease namespace from ID: nil namespace")
				return true
			}

			ctxWithNS := namespace.ContextWithNamespace(m.core.activeContext, leaseNS)
			ctxWithNSAndTimeout, _ := context.WithTimeout(ctxWithNS, time.Minute)
			if err := m.revokeCommon(ctxWithNSAndTimeout, leaseID, false, false); err != nil {
				// on failure, force some delay to mitigate resource spike while
				// this is running. if revocations succeed, we are okay with
				// the higher resource consumption.
				time.Sleep(10 * time.Millisecond)
			}
		}

		return true
	})
}

// revokeCommon does the heavy lifting. If force is true, we ignore a problem
// during revocation and still remove entries/index/lease timers
func (m *ExpirationManager) revokeCommon(ctx context.Context, leaseID string, force, skipToken bool) error {
	defer metrics.MeasureSince([]string{"expire", "revoke-common"}, time.Now())

	if !skipToken {
		// Acquire lock for this lease
		// If skipToken is true, then we're either being (1) called via RevokeByToken, so
		// probably the lock is already held, and if we re-acquire we get deadlock, or
		// (2) called by tidy, in which case the lock is held by the tidy thread.
		leaseLock := m.lockForLeaseID(leaseID)
		leaseLock.Lock()
		defer leaseLock.Unlock()
	}

	// Load the entry
	le, err := m.loadEntry(ctx, leaseID)
	if err != nil {
		return err
	}

	// If there is no entry, nothing to revoke
	if le == nil {
		return nil
	}

	// Revoke the entry
	if !skipToken || le.Auth == nil {
		if err := m.revokeEntry(ctx, le); err != nil {
			if !force {
				return err
			}

			if m.logger.IsWarn() {
				m.logger.Warn("revocation from the backend failed, but in force mode so ignoring", "error", err)
			}
		}
	}

	// Delete the entry
	if err := m.deleteEntry(ctx, le); err != nil {
		return err
	}

	// Lease has been removed, also remove the in-memory lock.
	m.deleteLockForLease(leaseID)

	// Delete the secondary index, but only if it's a leased secret (not auth)
	if le.Secret != nil {
		var indexToken string
		// Maintain secondary index by token, except for orphan batch tokens
		switch le.ClientTokenType {
		case logical.TokenTypeBatch:
			te, err := m.tokenStore.lookupBatchTokenInternal(ctx, le.ClientToken)
			if err != nil {
				return err
			}
			// If it's a non-orphan batch token, assign the secondary index to its
			// parent
			indexToken = te.Parent
		default:
			indexToken = le.ClientToken
		}
		if indexToken != "" {
			if err := m.removeIndexByToken(ctx, le, indexToken); err != nil {
				return err
			}
		}
	}

	// Clear the expiration handler
	m.pendingLock.Lock()
	m.removeFromPending(ctx, leaseID, true)
	m.nonexpiring.Delete(leaseID)

	if _, ok := m.irrevocable.Load(le.LeaseID); ok {
		m.irrevocable.Delete(leaseID)
		m.irrevocableLeaseCount--
	}
	m.pendingLock.Unlock()

	if m.logger.IsInfo() && !skipToken && m.logLeaseExpirations {
		m.logger.Info("revoked lease", "lease_id", leaseID)
	}
	if m.logger.IsWarn() && !skipToken && le.isIncorrectlyNonExpiring() {
		var accessor string
		if le.Auth != nil {
			accessor = le.Auth.Accessor
		}
		m.logger.Warn("finished revoking incorrectly non-expiring lease", "leaseID", le.LeaseID, "accessor", accessor)
	}
	return nil
}

// RevokeForce works similarly to RevokePrefix but continues in the case of a
// revocation error; this is mostly meant for recovery operations
func (m *ExpirationManager) RevokeForce(ctx context.Context, prefix string) error {
	defer metrics.MeasureSince([]string{"expire", "revoke-force"}, time.Now())

	return m.revokePrefixCommon(ctx, prefix, true, true)
}

// RevokePrefix is used to revoke all secrets with a given prefix.
// The prefix maps to that of the mount table to make this simpler
// to reason about.
func (m *ExpirationManager) RevokePrefix(ctx context.Context, prefix string, sync bool) error {
	defer metrics.MeasureSince([]string{"expire", "revoke-prefix"}, time.Now())

	return m.revokePrefixCommon(ctx, prefix, false, sync)
}

// RevokeByToken is used to revoke all the secrets issued with a given token.
// This is done by using the secondary index. It also removes the lease entry
// for the token itself. As a result it should *ONLY* ever be called from the
// token store's revokeInternal function.
// (NB: it's called by token tidy as well.)
func (m *ExpirationManager) RevokeByToken(ctx context.Context, te *logical.TokenEntry) error {
	defer metrics.MeasureSince([]string{"expire", "revoke-by-token"}, time.Now())
	tokenNS, err := NamespaceByID(ctx, te.NamespaceID, m.core)
	if err != nil {
		return err
	}
	if tokenNS == nil {
		return namespace.ErrNoNamespace
	}

	tokenCtx := namespace.ContextWithNamespace(ctx, tokenNS)
	// Lookup the leases
	existing, err := m.lookupLeasesByToken(tokenCtx, te)
	if err != nil {
		return fmt.Errorf("failed to scan for leases: %w", err)
	}

	// Revoke all the keys by marking them expired
	for _, leaseID := range existing {
		err := m.lazyRevokeInternal(ctx, leaseID)
		if err != nil {
			return err
		}
	}

	// te.Path should never be empty, but we check just in case
	if te.Path != "" {
		saltCtx := namespace.ContextWithNamespace(ctx, tokenNS)
		saltedID, err := m.tokenStore.SaltID(saltCtx, te.ID)
		if err != nil {
			return err
		}
		tokenLeaseID := path.Join(te.Path, saltedID)

		if tokenNS.ID != namespace.RootNamespaceID {
			tokenLeaseID = fmt.Sprintf("%s.%s", tokenLeaseID, tokenNS.ID)
		}

		// We want to skip the revokeEntry call as that will call back into
		// revocation logic in the token store, which is what is running this
		// function in the first place -- it'd be a deadlock loop. Since the only
		// place that this function is called is revokeSalted in the token store,
		// we're already revoking the token, so we just want to clean up the lease.
		// This avoids spurious revocations later in the log when the timer runs
		// out, and eases up resource usage.
		return m.revokeCommon(ctx, tokenLeaseID, false, true)
	}

	return nil
}

func (m *ExpirationManager) revokePrefixCommon(ctx context.Context, prefix string, force, sync bool) error {
	if m.inRestoreMode() {
		m.restoreRequestLock.Lock()
		defer m.restoreRequestLock.Unlock()
	}

	// Ensure there is a trailing slash; or, if there is no slash, see if there
	// is a matching specific ID
	if !strings.HasSuffix(prefix, "/") {
		le, err := m.loadEntry(ctx, prefix)
		if err == nil && le != nil {
			if sync {
				if err := m.revokeCommon(ctx, prefix, force, false); err != nil {
					return fmt.Errorf("failed to revoke %q: %w", prefix, err)
				}
				return nil
			}
			return m.LazyRevoke(ctx, prefix)
		}
		prefix = prefix + "/"
	}

	// Accumulate existing leases
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	view := m.leaseView(ns)
	sub := view.SubView(prefix)
	existing, err := logical.CollectKeys(ctx, sub)
	if err != nil {
		return fmt.Errorf("failed to scan for leases: %w", err)
	}

	// Revoke all the keys
	for idx, suffix := range existing {
		leaseID := prefix + suffix
		// No need to acquire per-lease lock here, one of these two will do it.
		switch {
		case sync:
			if err := m.revokeCommon(ctx, leaseID, force, false); err != nil {
				return fmt.Errorf("failed to revoke %q (%d / %d): %w", leaseID, idx+1, len(existing), err)
			}
		default:
			if err := m.LazyRevoke(ctx, leaseID); err != nil {
				return fmt.Errorf("failed to revoke %q (%d / %d): %w", leaseID, idx+1, len(existing), err)
			}
		}
	}

	return nil
}

// Renew is used to renew a secret using the given leaseID
// and a renew interval. The increment may be ignored.
func (m *ExpirationManager) Renew(ctx context.Context, leaseID string, increment time.Duration) (*logical.Response, error) {
	defer metrics.MeasureSince([]string{"expire", "renew"}, time.Now())

	// Acquire lock for this lease
	leaseLock := m.lockForLeaseID(leaseID)
	leaseLock.Lock()
	defer leaseLock.Unlock()

	// Load the entry
	le, err := m.loadEntry(ctx, leaseID)
	if err != nil {
		return nil, err
	}

	// Check if the lease is renewable
	if _, err := le.renewable(); err != nil {
		return nil, err
	}

	if le.Secret == nil {
		if le.Auth != nil {
			return logical.ErrorResponse("tokens cannot be renewed through this endpoint"), nil
		}
		return logical.ErrorResponse("lease does not correspond to a secret"), nil
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if ns.ID != le.namespace.ID {
		return nil, errors.New("cannot renew a lease across namespaces")
	}

	sysViewCtx := namespace.ContextWithNamespace(ctx, le.namespace)
	sysView := m.router.MatchingSystemView(sysViewCtx, le.Path)
	if sysView == nil {
		return nil, fmt.Errorf("unable to retrieve system view from router")
	}

	// Attempt to renew the entry
	resp, err := m.renewEntry(ctx, le, increment)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	if resp.IsError() {
		return &logical.Response{
			Data: resp.Data,
		}, nil
	}
	if resp.Secret == nil {
		return nil, nil
	}

	ttl, warnings, err := framework.CalculateTTL(sysView, increment, resp.Secret.TTL, 0, resp.Secret.MaxTTL, 0, le.IssueTime)
	if err != nil {
		return nil, err
	}
	for _, warning := range warnings {
		resp.AddWarning(warning)
	}
	resp.Secret.TTL = ttl

	// Attach the LeaseID
	resp.Secret.LeaseID = leaseID

	// Update the lease entry
	le.Data = resp.Data
	le.Secret = resp.Secret
	le.ExpireTime = resp.Secret.ExpirationTime()
	le.LastRenewalTime = time.Now()

	// If the token it's associated with is a batch token, constrain lease
	// times
	if le.ClientTokenType == logical.TokenTypeBatch {
		te, err := m.tokenStore.Lookup(ctx, le.ClientToken)
		if err != nil {
			return nil, err
		}
		if te == nil {
			return nil, errors.New("cannot renew lease, no valid associated token")
		}
		tokenLeaseTimes, err := m.FetchLeaseTimesByToken(ctx, te)
		if err != nil {
			return nil, err
		}

		if tokenLeaseTimes == nil {
			return nil, errors.New("failed to load batch token expiration time")
		}

		if le.ExpireTime.After(tokenLeaseTimes.ExpireTime) {
			resp.Secret.TTL = tokenLeaseTimes.ExpireTime.Sub(le.LastRenewalTime)
			le.ExpireTime = tokenLeaseTimes.ExpireTime
		}
	}

	if err := m.persistEntry(ctx, le); err != nil {
		return nil, err
	}

	// Update the expiration time
	m.updatePending(le)

	// Return the response
	return resp, nil
}

// RenewToken is used to renew a token which does not need to
// invoke a logical backend.
func (m *ExpirationManager) RenewToken(ctx context.Context, req *logical.Request, te *logical.TokenEntry,
	increment time.Duration) (*logical.Response, error) {
	defer metrics.MeasureSince([]string{"expire", "renew-token"}, time.Now())

	tokenNS, err := NamespaceByID(ctx, te.NamespaceID, m.core)
	if err != nil {
		return nil, err
	}
	if tokenNS == nil {
		return nil, namespace.ErrNoNamespace
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if ns.ID != tokenNS.ID {
		return nil, errors.New("cannot renew a token across namespaces")
	}

	// Compute the Lease ID
	saltedID, err := m.tokenStore.SaltID(ctx, te.ID)
	if err != nil {
		return nil, err
	}

	leaseID := path.Join(te.Path, saltedID)

	if ns.ID != namespace.RootNamespaceID {
		leaseID = fmt.Sprintf("%s.%s", leaseID, ns.ID)
	}

	// Acquire lock for this lease
	leaseLock := m.lockForLeaseID(leaseID)
	leaseLock.Lock()
	defer leaseLock.Unlock()

	// Load the entry
	le, err := m.loadEntry(ctx, leaseID)
	if err != nil {
		return nil, err
	}
	if le == nil {
		return logical.ErrorResponse("invalid lease ID"), logical.ErrInvalidRequest
	}

	// Check if the lease is renewable. Note that this also checks for a nil
	// lease and errors in that case as well.
	if _, err := le.renewable(); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Attempt to renew the auth entry
	resp, err := m.renewAuthEntry(ctx, req, le, increment)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	if resp.IsError() {
		return &logical.Response{
			Data: resp.Data,
		}, nil
	}
	if resp.Auth == nil {
		return nil, nil
	}

	sysViewCtx := namespace.ContextWithNamespace(ctx, le.namespace)
	sysView := m.router.MatchingSystemView(sysViewCtx, le.Path)
	if sysView == nil {
		return nil, fmt.Errorf("unable to retrieve system view from router")
	}

	ttl, warnings, err := framework.CalculateTTL(sysView, increment, resp.Auth.TTL, resp.Auth.Period, resp.Auth.MaxTTL, resp.Auth.ExplicitMaxTTL, le.IssueTime)
	if err != nil {
		return nil, err
	}
	retResp := &logical.Response{}
	for _, warning := range warnings {
		retResp.AddWarning(warning)
	}
	resp.Auth.TTL = ttl

	// Attach the ClientToken
	resp.Auth.ClientToken = te.ID

	// Refresh groups
	if resp.Auth.EntityID != "" && m.core.identityStore != nil {
		mountAccessor := ""
		if resp.Auth.Alias != nil {
			mountAccessor = resp.Auth.Alias.MountAccessor
		}
		validAliases, err := m.core.identityStore.refreshExternalGroupMembershipsByEntityID(ctx, resp.Auth.EntityID, resp.Auth.GroupAliases, mountAccessor)
		if err != nil {
			return nil, err
		}
		resp.Auth.GroupAliases = validAliases
	}

	// Update the lease entry
	le.Auth = resp.Auth
	le.ExpireTime = resp.Auth.ExpirationTime()
	le.LastRenewalTime = time.Now()

	if err := m.persistEntry(ctx, le); err != nil {
		return nil, err
	}
	m.updatePending(le)

	retResp.Auth = resp.Auth
	return retResp, nil
}

// Register is used to take a request and response with an associated
// lease. The secret gets assigned a LeaseID and the management of
// of lease is assumed by the expiration manager.
func (m *ExpirationManager) Register(ctx context.Context, req *logical.Request, resp *logical.Response) (id string, retErr error) {
	defer metrics.MeasureSince([]string{"expire", "register"}, time.Now())

	te := req.TokenEntry()
	if te == nil {
		return "", fmt.Errorf("cannot register a lease with an empty client token")
	}

	// Ignore if there is no leased secret
	if resp == nil || resp.Secret == nil {
		return "", nil
	}

	// Validate the secret
	if err := resp.Secret.Validate(); err != nil {
		return "", err
	}

	// Create a lease entry
	leaseRand, err := base62.Random(TokenLength)
	if err != nil {
		return "", err
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return "", err
	}

	leaseID := path.Join(req.Path, leaseRand)

	if ns.ID != namespace.RootNamespaceID {
		leaseID = fmt.Sprintf("%s.%s", leaseID, ns.ID)
	}

	le := &leaseEntry{
		LeaseID:         leaseID,
		ClientToken:     req.ClientToken,
		ClientTokenType: te.Type,
		Path:            req.Path,
		Data:            resp.Data,
		Secret:          resp.Secret,
		IssueTime:       time.Now(),
		ExpireTime:      resp.Secret.ExpirationTime(),
		namespace:       ns,
		Version:         1,
	}

	var indexToken string
	// Maintain secondary index by token, except for orphan batch tokens
	switch {
	case te.Type != logical.TokenTypeBatch:
		indexToken = le.ClientToken
	case te.Parent != "":
		// If it's a non-orphan batch token, assign the secondary index to its
		// parent
		indexToken = te.Parent
	}

	defer func() {
		// If there is an error we want to rollback as much as possible (note
		// that errors here are ignored to do as much cleanup as we can). We
		// want to revoke a generated secret (since an error means we may not
		// be successfully tracking it), remove indexes, and delete the entry.
		if retErr != nil {
			revokeCtx := namespace.ContextWithNamespace(m.quitContext, ns)
			revResp, err := m.router.Route(revokeCtx, logical.RevokeRequest(req.Path, resp.Secret, resp.Data))
			if err != nil {
				retErr = multierror.Append(retErr, fmt.Errorf("an additional internal error was encountered revoking the newly-generated secret: %w", err))
			} else if revResp != nil && revResp.IsError() {
				retErr = multierror.Append(retErr, fmt.Errorf("an additional error was encountered revoking the newly-generated secret: %w", revResp.Error()))
			}

			if err := m.deleteEntry(ctx, le); err != nil {
				retErr = multierror.Append(retErr, fmt.Errorf("an additional error was encountered deleting any lease associated with the newly-generated secret: %w", err))
			}

			if err := m.removeIndexByToken(ctx, le, indexToken); err != nil {
				retErr = multierror.Append(retErr, fmt.Errorf("an additional error was encountered removing lease indexes associated with the newly-generated secret: %w", err))
			}

			m.deleteLockForLease(leaseID)
		}
	}()

	// If the token is a batch token, we want to constrain the maximum lifetime
	// by the token's lifetime
	if te.Type == logical.TokenTypeBatch {
		tokenLeaseTimes, err := m.FetchLeaseTimesByToken(ctx, te)
		if err != nil {
			return "", err
		}
		if tokenLeaseTimes == nil {
			return "", errors.New("failed to load batch token expiration time")
		}
		if le.ExpireTime.After(tokenLeaseTimes.ExpireTime) {
			le.ExpireTime = tokenLeaseTimes.ExpireTime
		}
	}

	// Acquire the lock here so persistEntry and updatePending are atomic,
	// although it is *very unlikely* that anybody could grab the lease ID
	// before this function returns. (They could find it in an index, or
	// find it in a list.)
	leaseLock := m.lockForLeaseID(leaseID)
	leaseLock.Lock()
	defer leaseLock.Unlock()

	// Encode the entry
	if err := m.persistEntry(ctx, le); err != nil {
		return "", err
	}

	if indexToken != "" {
		if err := m.createIndexByToken(ctx, le, indexToken); err != nil {
			return "", err
		}
	}

	// Setup revocation timer if there is a lease
	m.updatePending(le)

	// We round here because the clock will have already started
	// ticking, so we'll end up always returning 299 instead of 300 or
	// 26399 instead of 26400, say, even if it's just a few
	// microseconds. This provides a nicer UX.
	resp.Secret.TTL = le.ExpireTime.Sub(time.Now()).Round(time.Second)

	// Done
	return le.LeaseID, nil
}

// RegisterAuth is used to take an Auth response with an associated lease.
// The token does not get a LeaseID, but the lease management is handled by
// the expiration manager.
func (m *ExpirationManager) RegisterAuth(ctx context.Context, te *logical.TokenEntry, auth *logical.Auth) error {
	defer metrics.MeasureSince([]string{"expire", "register-auth"}, time.Now())

	// Triggers failure of RegisterAuth. This should only be set and triggered
	// by tests to simulate partial failure during a token creation request.
	if m.testRegisterAuthFailure.Load() {
		return fmt.Errorf("failing explicitly on RegisterAuth")
	}

	authExpirationTime := auth.ExpirationTime()

	if te.TTL == 0 && authExpirationTime.IsZero() && (len(te.Policies) != 1 || te.Policies[0] != "root") {
		return errors.New("refusing to register a lease for a non-root token with no TTL")
	}

	if te.Type == logical.TokenTypeBatch {
		return errors.New("cannot register a lease for a batch token")
	}

	if auth.ClientToken == "" {
		return errors.New("cannot register an auth lease with an empty token")
	}

	if strings.Contains(te.Path, "..") {
		return consts.ErrPathContainsParentReferences
	}

	tokenNS, err := NamespaceByID(ctx, te.NamespaceID, m.core)
	if err != nil {
		return err
	}
	if tokenNS == nil {
		return namespace.ErrNoNamespace
	}

	saltCtx := namespace.ContextWithNamespace(ctx, tokenNS)
	saltedID, err := m.tokenStore.SaltID(saltCtx, auth.ClientToken)
	if err != nil {
		return err
	}

	leaseID := path.Join(te.Path, saltedID)
	if tokenNS.ID != namespace.RootNamespaceID {
		leaseID = fmt.Sprintf("%s.%s", leaseID, tokenNS.ID)
	}

	// Create a lease entry
	le := leaseEntry{
		LeaseID:     leaseID,
		ClientToken: auth.ClientToken,
		Auth:        auth,
		Path:        te.Path,
		IssueTime:   time.Now(),
		ExpireTime:  authExpirationTime,
		namespace:   tokenNS,
		Version:     1,
	}

	leaseLock := m.lockForLeaseID(leaseID)
	leaseLock.Lock()
	defer leaseLock.Unlock()

	// Encode the entry
	if err := m.persistEntry(ctx, &le); err != nil {
		return err
	}

	// Setup revocation timer
	m.updatePending(&le)

	return nil
}

// FetchLeaseTimesByToken is a helper function to use token values to compute
// the leaseID, rather than pushing that logic back into the token store.
// As a special case, for a batch token it simply returns the information
// encoded on it.
func (m *ExpirationManager) FetchLeaseTimesByToken(ctx context.Context, te *logical.TokenEntry) (*leaseEntry, error) {
	defer metrics.MeasureSince([]string{"expire", "fetch-lease-times-by-token"}, time.Now())

	if te == nil {
		return nil, errors.New("cannot fetch lease times for nil token")
	}

	if te.Type == logical.TokenTypeBatch {
		issueTime := time.Unix(te.CreationTime, 0)
		return &leaseEntry{
			IssueTime:       issueTime,
			ExpireTime:      issueTime.Add(te.TTL),
			ClientTokenType: logical.TokenTypeBatch,
		}, nil
	}

	tokenNS, err := NamespaceByID(ctx, te.NamespaceID, m.core)
	if err != nil {
		return nil, err
	}
	if tokenNS == nil {
		return nil, namespace.ErrNoNamespace
	}

	saltCtx := namespace.ContextWithNamespace(ctx, tokenNS)
	saltedID, err := m.tokenStore.SaltID(saltCtx, te.ID)
	if err != nil {
		return nil, err
	}

	leaseID := path.Join(te.Path, saltedID)

	if tokenNS.ID != namespace.RootNamespaceID {
		leaseID = fmt.Sprintf("%s.%s", leaseID, tokenNS.ID)
	}

	return m.FetchLeaseTimes(ctx, leaseID)
}

// FetchLeaseTimes is used to fetch the issue time, expiration time, and last
// renewed time of a lease entry. It returns a leaseEntry itself, but with only
// those values copied over.
func (m *ExpirationManager) FetchLeaseTimes(ctx context.Context, leaseID string) (*leaseEntry, error) {
	defer metrics.MeasureSince([]string{"expire", "fetch-lease-times"}, time.Now())

	info, ok := m.pending.Load(leaseID)
	if ok && info.(pendingInfo).cachedLeaseInfo != nil {
		return m.leaseTimesForExport(info.(pendingInfo).cachedLeaseInfo), nil
	}

	info, ok = m.irrevocable.Load(leaseID)
	if ok && info.(*leaseEntry) != nil {
		return m.leaseTimesForExport(info.(*leaseEntry)), nil
	}

	// Load the entry
	le, err := m.loadEntryInternal(ctx, leaseID, true, false)
	if err != nil {
		return nil, err
	}
	if le == nil {
		return nil, nil
	}

	return m.leaseTimesForExport(le), nil
}

// Returns lease times for outside callers based on the full leaseEntry passed in
func (m *ExpirationManager) leaseTimesForExport(le *leaseEntry) *leaseEntry {
	ret := &leaseEntry{
		IssueTime:       le.IssueTime,
		ExpireTime:      le.ExpireTime,
		LastRenewalTime: le.LastRenewalTime,
	}
	if le.Secret != nil {
		ret.Secret = &logical.Secret{}
		ret.Secret.Renewable = le.Secret.Renewable
		ret.Secret.TTL = le.Secret.TTL
	}
	if le.Auth != nil {
		ret.Auth = &logical.Auth{}
		ret.Auth.Renewable = le.Auth.Renewable
		ret.Auth.TTL = le.Auth.TTL
	}

	return ret
}

// Restricts lease entry stored in pendingInfo to a low-cost subset of the
// information.
func (m *ExpirationManager) inMemoryLeaseInfo(le *leaseEntry) *leaseEntry {
	ret := m.leaseTimesForExport(le)
	// Need to index:
	//   namespace -- derived from lease ID
	//   policies -- stored in Auth object
	//   auth method -- derived from lease.Path
	if le.Auth != nil {
		// Ensure that list of policies is not copied more than
		// once. This method is called with pendingLock held.

		// We could use hashstructure here to generate a key, but that
		// seems like it would be substantially slower?
		key := strings.Join(le.Auth.Policies, "\n")
		uniq, ok := m.uniquePolicies[key]
		if ok {
			ret.Auth.Policies = uniq
		} else {
			m.uniquePolicies[key] = le.Auth.Policies
			ret.Auth.Policies = le.Auth.Policies
		}
		ret.Path = le.Path
	}
	if le.isIrrevocable() {
		ret.RevokeErr = le.RevokeErr
	}
	return ret
}

func (m *ExpirationManager) uniquePoliciesGc() {
	for {
		<-m.emptyUniquePolicies.C

		// If the maximum lease is a month, and we blow away the unique
		// policy cache every week, the pessimal case is 4x larger space
		// utilization than keeping the cache indefinitely.
		m.pendingLock.Lock()
		m.uniquePolicies = make(map[string][]string)
		m.pendingLock.Unlock()
	}
}

// Placing a lock in pendingMap means that we need to work very hard on reload
// to only create one lock.  Instead, we'll create locks on-demand in an atomic fashion.
//
// Acquiring a lock from a leaseEntry is a bad idea because it could change
// between loading and acquiring the lock. So we only provide an ID-based map, and the
// locking discipline should be:
//    1. Lock lease
//    2. Load, or attempt to load, leaseEntry
//    3. Modify leaseEntry and pendingMap (atomic wrt operations on this lease)
//    4. Unlock lease
//
// The lock must be removed from the map when the lease is deleted, or is
// found to not exist in storage. loadEntry does this whenever it returns
// nil, but we should also do it in revokeCommon().
func (m *ExpirationManager) lockForLeaseID(id string) *sync.Mutex {
	mutex := &sync.Mutex{}
	lock, _ := m.lockPerLease.LoadOrStore(id, mutex)
	return lock.(*sync.Mutex)
}

func (m *ExpirationManager) deleteLockForLease(id string) {
	m.lockPerLease.Delete(id)
}

// updatePending is used to update a pending invocation for a lease
func (m *ExpirationManager) updatePending(le *leaseEntry) {
	m.pendingLock.Lock()
	defer m.pendingLock.Unlock()

	m.updatePendingInternal(le)
}

// updatePendingInternal is the locked version of updatePending; do not call
// this without a write lock on m.pending
func (m *ExpirationManager) updatePendingInternal(le *leaseEntry) {
	// Check for an existing timer
	info, leaseInPending := m.pending.Load(le.LeaseID)

	var pending pendingInfo

	if le.ExpireTime.IsZero() && le.nonexpiringToken() {
		// Store this in the nonexpiring map instead of pending.
		// There does not appear to be any cases where a token that had
		// a nonzero can be can be assigned a zero TTL, but that can be
		// handled by the next check
		pending.cachedLeaseInfo = m.inMemoryLeaseInfo(le)
		m.nonexpiring.Store(le.LeaseID, pending)

		// if the timer happened to exist, stop the time and delete it from the
		// pending timers.
		if leaseInPending {
			info.(pendingInfo).timer.Stop()
			m.pending.Delete(le.LeaseID)
			m.leaseCount--
			if err := m.core.quotasHandleLeases(m.quitContext, quotas.LeaseActionDeleted, []string{le.LeaseID}); err != nil {
				m.logger.Error("failed to update quota on lease deletion", "error", err)
				return
			}
		}
		return
	}

	leaseTotal := le.ExpireTime.Sub(time.Now())
	leaseCreated := false

	if le.isIrrevocable() {
		// It's possible this function is being called to update the in-memory state
		// for a lease from pending to irrevocable (we don't support the opposite).
		// If this is the case, we need to know if the lease was previously counted
		// so that we can maintain correct metric and quota lease counts.
		_, leaseInIrrevocable := m.irrevocable.Load(le.LeaseID)
		if !(leaseInPending || leaseInIrrevocable) {
			leaseCreated = true
		}

		m.removeFromPending(m.quitContext, le.LeaseID, false)
		m.irrevocable.Store(le.LeaseID, m.inMemoryLeaseInfo(le))

		// Increment count if the lease was not present in the irrevocable map
		// prior to being added to it above
		if !leaseInIrrevocable {
			m.irrevocableLeaseCount++
		}
	} else {
		// Create entry if it does not exist or reset if it does
		if leaseInPending {
			pending = info.(pendingInfo)
			pending.timer.Reset(leaseTotal)
			// No change to lease count in this case
		} else {
			leaseID, namespace := le.LeaseID, le.namespace
			// Extend the timer by the lease total
			timer := time.AfterFunc(leaseTotal, func() {
				m.expireFunc(m.quitContext, m, leaseID, namespace)
			})
			pending = pendingInfo{
				timer: timer,
			}

			leaseCreated = true
		}

		pending.cachedLeaseInfo = m.inMemoryLeaseInfo(le)
		m.pending.Store(le.LeaseID, pending)
	}

	if leaseCreated {
		m.leaseCount++
		if err := m.core.quotasHandleLeases(m.quitContext, quotas.LeaseActionCreated, []string{le.LeaseID}); err != nil {
			m.logger.Error("failed to update quota on lease creation", "error", err)
			return
		}
	}
}

// revokeEntry is used to attempt revocation of an internal entry
func (m *ExpirationManager) revokeEntry(ctx context.Context, le *leaseEntry) error {
	// Revocation of login tokens is special since we can by-pass the
	// backend and directly interact with the token store
	if le.Auth != nil {
		if le.ClientTokenType == logical.TokenTypeBatch {
			return errors.New("batch tokens cannot be revoked")
		}

		if err := m.tokenStore.revokeTree(ctx, le); err != nil {
			return fmt.Errorf("failed to revoke token: %w", err)
		}

		return nil
	}

	if le.Secret != nil {
		// not sure if this is really valid to have a leaseEntry with a nil Secret
		// (if there's a nil Secret, what are you really leasing?), but the tests
		// create one, and good to be defensive
		le.Secret.IssueTime = le.IssueTime
	}

	// Make sure we're operating in the right namespace
	nsCtx := namespace.ContextWithNamespace(ctx, le.namespace)

	// Handle standard revocation via backends
	resp, err := m.router.Route(nsCtx, logical.RevokeRequest(le.Path, le.Secret, le.Data))
	if err != nil || (resp != nil && resp.IsError()) {
		return fmt.Errorf("failed to revoke entry: resp: %#v err: %w", resp, err)
	}
	return nil
}

// renewEntry is used to attempt renew of an internal entry
func (m *ExpirationManager) renewEntry(ctx context.Context, le *leaseEntry, increment time.Duration) (*logical.Response, error) {
	secret := *le.Secret
	secret.IssueTime = le.IssueTime
	secret.Increment = increment
	secret.LeaseID = ""

	// Make sure we're operating in the right namespace
	nsCtx := namespace.ContextWithNamespace(ctx, le.namespace)

	req := logical.RenewRequest(le.Path, &secret, le.Data)
	resp, err := m.router.Route(nsCtx, req)
	if err != nil || (resp != nil && resp.IsError()) {
		return nil, fmt.Errorf("failed to renew entry: resp: %#v err: %w", resp, err)
	}
	return resp, nil
}

// renewAuthEntry is used to attempt renew of an auth entry. Only the token
// store should get the actual token ID intact.
func (m *ExpirationManager) renewAuthEntry(ctx context.Context, req *logical.Request, le *leaseEntry, increment time.Duration) (*logical.Response, error) {
	if le.ClientTokenType == logical.TokenTypeBatch {
		return logical.ErrorResponse("batch tokens cannot be renewed"), nil
	}

	auth := *le.Auth
	auth.IssueTime = le.IssueTime
	auth.Increment = increment
	if strings.HasPrefix(le.Path, "auth/token/") {
		auth.ClientToken = le.ClientToken
	} else {
		auth.ClientToken = ""
	}

	// Make sure we're operating in the right namespace
	nsCtx := namespace.ContextWithNamespace(ctx, le.namespace)

	authReq := logical.RenewAuthRequest(le.Path, &auth, nil)
	authReq.Connection = req.Connection
	resp, err := m.router.Route(nsCtx, authReq)
	if err != nil {
		return nil, fmt.Errorf("failed to renew entry: %w", err)
	}
	return resp, nil
}

// loadEntry is used to read a lease entry
func (m *ExpirationManager) loadEntry(ctx context.Context, leaseID string) (*leaseEntry, error) {
	// Take out the lease locks after we ensure we are in restore mode
	restoreMode := m.inRestoreMode()
	if restoreMode {
		m.restoreModeLock.RLock()
		defer m.restoreModeLock.RUnlock()

		restoreMode = m.inRestoreMode()
		if restoreMode {
			m.lockLease(leaseID)
			defer m.unlockLease(leaseID)
		}
	}

	_, nsID := namespace.SplitIDFromString(leaseID)
	if nsID != "" {
		leaseNS, err := NamespaceByID(ctx, nsID, m.core)
		if err != nil {
			return nil, err
		}
		if leaseNS != nil {
			ctx = namespace.ContextWithNamespace(ctx, leaseNS)
		}
	} else {
		ctx = namespace.ContextWithNamespace(ctx, namespace.RootNamespace)
	}

	// If a lease entry is nil, proactively delete the lease lock, in case we
	// created one erroneously.
	// If there was an error, we don't know whether the lease entry exists or not.
	leaseEntry, err := m.loadEntryInternal(ctx, leaseID, restoreMode, true)
	if err == nil && leaseEntry == nil {
		m.deleteLockForLease(leaseID)
	}
	return leaseEntry, err

}

// loadEntryInternal is used when you need to load an entry but also need to
// control the lifecycle of the restoreLock
func (m *ExpirationManager) loadEntryInternal(ctx context.Context, leaseID string, restoreMode bool, checkRestored bool) (*leaseEntry, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	view := m.leaseView(ns)
	out, err := view.Get(ctx, leaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to read lease entry %s: %w", leaseID, err)
	}
	if out == nil {
		return nil, nil
	}
	le, err := decodeLeaseEntry(out.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode lease entry %s: %w", leaseID, err)
	}
	le.namespace = ns

	if restoreMode {
		if checkRestored {
			// If we have already loaded this lease, we don't need to update on
			// load. In the case of renewal and revocation, updatePending will be
			// done after making the appropriate modifications to the lease.
			if _, ok := m.restoreLoaded.Load(leaseID); ok {
				return le, nil
			}
		}

		// Update the cache of restored leases, either synchronously or through
		// the lazy loaded restore process
		m.restoreLoaded.Store(le.LeaseID, struct{}{})

		// Setup revocation timer
		m.updatePending(le)
	}
	return le, nil
}

// persistEntry is used to persist a lease entry
func (m *ExpirationManager) persistEntry(ctx context.Context, le *leaseEntry) error {
	// Encode the entry
	buf, err := le.encode()
	if err != nil {
		return fmt.Errorf("failed to encode lease entry: %w", err)
	}

	// Write out to the view
	ent := logical.StorageEntry{
		Key:   le.LeaseID,
		Value: buf,
	}
	if le.Auth != nil && len(le.Auth.Policies) == 1 && le.Auth.Policies[0] == "root" {
		ent.SealWrap = true
	}

	view := m.leaseView(le.namespace)
	if err := view.Put(ctx, &ent); err != nil {
		return fmt.Errorf("failed to persist lease entry: %w", err)
	}
	return nil
}

// deleteEntry is used to delete a lease entry
func (m *ExpirationManager) deleteEntry(ctx context.Context, le *leaseEntry) error {
	view := m.leaseView(le.namespace)
	if err := view.Delete(ctx, le.LeaseID); err != nil {
		return fmt.Errorf("failed to delete lease entry: %w", err)
	}
	return nil
}

// createIndexByToken creates a secondary index from the token to a lease entry
func (m *ExpirationManager) createIndexByToken(ctx context.Context, le *leaseEntry, token string) error {
	tokenNS := namespace.RootNamespace
	saltCtx := namespace.ContextWithNamespace(ctx, namespace.RootNamespace)
	_, nsID := namespace.SplitIDFromString(token)
	if nsID != "" {
		var err error
		tokenNS, err = NamespaceByID(ctx, nsID, m.core)
		if err != nil {
			return err
		}
		if tokenNS != nil {
			saltCtx = namespace.ContextWithNamespace(ctx, tokenNS)
		}
	}

	saltedID, err := m.tokenStore.SaltID(saltCtx, token)
	if err != nil {
		return err
	}

	leaseSaltedID, err := m.tokenStore.SaltID(saltCtx, le.LeaseID)
	if err != nil {
		return err
	}

	ent := logical.StorageEntry{
		Key:   saltedID + "/" + leaseSaltedID,
		Value: []byte(le.LeaseID),
	}
	tokenView := m.tokenIndexView(tokenNS)
	if err := tokenView.Put(ctx, &ent); err != nil {
		return fmt.Errorf("failed to persist lease index entry: %w", err)
	}
	return nil
}

// indexByToken looks up the secondary index from the token to a lease entry
func (m *ExpirationManager) indexByToken(ctx context.Context, le *leaseEntry) (*logical.StorageEntry, error) {
	tokenNS := namespace.RootNamespace
	saltCtx := namespace.ContextWithNamespace(ctx, tokenNS)
	_, nsID := namespace.SplitIDFromString(le.ClientToken)
	if nsID != "" {
		var err error
		tokenNS, err = NamespaceByID(ctx, nsID, m.core)
		if err != nil {
			return nil, err
		}
		if tokenNS != nil {
			saltCtx = namespace.ContextWithNamespace(ctx, tokenNS)
		}
	}

	saltedID, err := m.tokenStore.SaltID(saltCtx, le.ClientToken)
	if err != nil {
		return nil, err
	}

	leaseSaltedID, err := m.tokenStore.SaltID(saltCtx, le.LeaseID)
	if err != nil {
		return nil, err
	}

	key := saltedID + "/" + leaseSaltedID
	tokenView := m.tokenIndexView(tokenNS)
	entry, err := tokenView.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to look up secondary index entry")
	}
	return entry, nil
}

// removeIndexByToken removes the secondary index from the token to a lease entry
func (m *ExpirationManager) removeIndexByToken(ctx context.Context, le *leaseEntry, token string) error {
	tokenNS := namespace.RootNamespace
	saltCtx := namespace.ContextWithNamespace(ctx, namespace.RootNamespace)
	_, nsID := namespace.SplitIDFromString(token)
	if nsID != "" {
		var err error
		tokenNS, err = NamespaceByID(ctx, nsID, m.core)
		if err != nil {
			return err
		}
		if tokenNS != nil {
			saltCtx = namespace.ContextWithNamespace(ctx, tokenNS)
		}

		// Downgrade logic for old-style (V0) namespace leases that had its
		// secondary index live in the root namespace. This reverts to the old
		// behavior of looking for the secondary index on these leases in the
		// root namespace to be cleaned up properly. We set it here because the
		// old behavior used the namespace's token store salt for its saltCtx.
		if le.Version < 1 {
			tokenNS = namespace.RootNamespace
		}
	}

	saltedID, err := m.tokenStore.SaltID(saltCtx, token)
	if err != nil {
		return err
	}

	leaseSaltedID, err := m.tokenStore.SaltID(saltCtx, le.LeaseID)
	if err != nil {
		return err
	}

	key := saltedID + "/" + leaseSaltedID
	tokenView := m.tokenIndexView(tokenNS)
	if err := tokenView.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete lease index entry: %w", err)
	}
	return nil
}

// CreateOrFetchRevocationLeaseByToken is used to create or fetch the matching
// leaseID for a particular token. The lease is set to expire immediately after
// it's created.
func (m *ExpirationManager) CreateOrFetchRevocationLeaseByToken(ctx context.Context, te *logical.TokenEntry) (string, error) {
	// Fetch the saltedID of the token and construct the leaseID
	tokenNS, err := NamespaceByID(ctx, te.NamespaceID, m.core)
	if err != nil {
		return "", err
	}
	if tokenNS == nil {
		return "", namespace.ErrNoNamespace
	}

	saltCtx := namespace.ContextWithNamespace(ctx, tokenNS)
	saltedID, err := m.tokenStore.SaltID(saltCtx, te.ID)
	if err != nil {
		return "", err
	}
	leaseID := path.Join(te.Path, saltedID)

	if tokenNS.ID != namespace.RootNamespaceID {
		leaseID = fmt.Sprintf("%s.%s", leaseID, tokenNS.ID)
	}

	// Load the entry
	le, err := m.loadEntry(ctx, leaseID)
	if err != nil {
		return "", err
	}

	// If there's no associated leaseEntry for the token, we create one
	if le == nil {

		// Acquire the lock here so persistEntry and updatePending are atomic,
		// although it is *very unlikely* that anybody could grab the lease ID
		// before this function returns. (They could find it in an index, or
		// find it in a list.)
		leaseLock := m.lockForLeaseID(leaseID)
		leaseLock.Lock()
		defer leaseLock.Unlock()

		auth := &logical.Auth{
			ClientToken: te.ID,
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Nanosecond,
			},
		}

		if strings.Contains(te.Path, "..") {
			return "", consts.ErrPathContainsParentReferences
		}

		// Create a lease entry
		now := time.Now()
		le = &leaseEntry{
			LeaseID:     leaseID,
			ClientToken: auth.ClientToken,
			Auth:        auth,
			Path:        te.Path,
			IssueTime:   now,
			ExpireTime:  now.Add(time.Nanosecond),
			namespace:   tokenNS,
			Version:     1,
		}

		// Encode the entry
		if err := m.persistEntry(ctx, le); err != nil {
			m.deleteLockForLease(leaseID)
			return "", err
		}
	}

	return le.LeaseID, nil
}

// lookupLeasesByToken is used to lookup all the leaseID's via the tokenID
func (m *ExpirationManager) lookupLeasesByToken(ctx context.Context, te *logical.TokenEntry) ([]string, error) {
	tokenNS, err := NamespaceByID(ctx, te.NamespaceID, m.core)
	if err != nil {
		return nil, err
	}
	if tokenNS == nil {
		return nil, namespace.ErrNoNamespace
	}

	saltCtx := namespace.ContextWithNamespace(ctx, tokenNS)
	saltedID, err := m.tokenStore.SaltID(saltCtx, te.ID)
	if err != nil {
		return nil, err
	}

	tokenView := m.tokenIndexView(tokenNS)

	// Scan via the index for sub-leases
	prefix := saltedID + "/"
	subKeys, err := tokenView.List(ctx, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to list leases: %w", err)
	}

	// Read each index entry
	leaseIDs := make([]string, 0, len(subKeys))
	for _, sub := range subKeys {
		out, err := tokenView.Get(ctx, prefix+sub)
		if err != nil {
			return nil, fmt.Errorf("failed to read lease index: %w", err)
		}
		if out == nil {
			continue
		}
		leaseIDs = append(leaseIDs, string(out.Value))
	}

	// Downgrade logic for old-style (V0) leases entries created by a namespace
	// token that lived in the root namespace.
	if tokenNS.ID != namespace.RootNamespaceID {
		tokenView := m.tokenIndexView(namespace.RootNamespace)

		// Scan via the index for sub-leases on the root namespace
		prefix := saltedID + "/"
		subKeys, err := tokenView.List(ctx, prefix)
		if err != nil {
			return nil, fmt.Errorf("failed to list leases on root namespace: %w", err)
		}

		for _, sub := range subKeys {
			out, err := tokenView.Get(ctx, prefix+sub)
			if err != nil {
				return nil, fmt.Errorf("failed to read lease index on root namespace: %w", err)
			}
			if out == nil {
				continue
			}
			leaseIDs = append(leaseIDs, string(out.Value))
		}
	}

	return leaseIDs, nil
}

// emitMetrics is invoked periodically to emit statistics
func (m *ExpirationManager) emitMetrics() {
	// All updates of these values are with the pendingLock held.
	m.pendingLock.RLock()
	allLeases := m.leaseCount
	irrevocableLeases := m.irrevocableLeaseCount
	m.pendingLock.RUnlock()

	metrics.SetGauge([]string{"expire", "num_leases"}, float32(allLeases))

	metrics.SetGauge([]string{"expire", "num_irrevocable_leases"}, float32(irrevocableLeases))
	// Check if lease count is greater than the threshold
	if allLeases > maxLeaseThreshold {
		if atomic.LoadUint32(m.leaseCheckCounter) > 59 {
			m.logger.Warn("lease count exceeds warning lease threshold", "have", allLeases, "threshold", maxLeaseThreshold)
			atomic.StoreUint32(m.leaseCheckCounter, 0)
		} else {
			atomic.AddUint32(m.leaseCheckCounter, 1)
		}
	}
}

func (m *ExpirationManager) leaseAggregationMetrics(ctx context.Context, consts metricsutil.TelemetryConstConfig) ([]metricsutil.GaugeLabelValues, error) {
	expiryTimes := make(map[metricsutil.LeaseExpiryLabel]int)
	leaseEpsilon := consts.LeaseMetricsEpsilon
	nsLabel := consts.LeaseMetricsNameSpaceLabels

	rollingWindow := time.Now().Add(time.Duration(consts.NumLeaseMetricsTimeBuckets) * leaseEpsilon)

	err := m.walkLeases(func(entryID string, expireTime time.Time) bool {
		select {
		// Abort and return empty collection if it's taking too much time, nonblocking check.
		case <-ctx.Done():
			return false
		default:
			if entryID == "" {
				return true
			}
			_, nsID := namespace.SplitIDFromString(entryID)
			if nsID == "" {
				nsID = "root" // this is what metricsutil.NamespaceLabel does
			}
			label := metricsutil.ExpiryBucket(expireTime, leaseEpsilon, rollingWindow, nsID, nsLabel)
			if label != nil {
				expiryTimes[*label] += 1
			}
			return true
		}
	})
	if err != nil {
		return []metricsutil.GaugeLabelValues{}, suppressRestoreModeError(err)
	}

	// If collection was cancelled, return an empty array.
	select {
	case <-ctx.Done():
		return []metricsutil.GaugeLabelValues{}, nil
	default:
		break
	}

	flattenedResults := make([]metricsutil.GaugeLabelValues, 0, len(expiryTimes))

	for bucket, count := range expiryTimes {
		if nsLabel {
			flattenedResults = append(flattenedResults,
				metricsutil.GaugeLabelValues{
					Labels: []metrics.Label{{"expiring", bucket.LabelName}, {"namespace", bucket.LabelNS}},
					Value:  float32(count),
				})
		} else {
			flattenedResults = append(flattenedResults,
				metricsutil.GaugeLabelValues{
					Labels: []metrics.Label{{"expiring", bucket.LabelName}},
					Value:  float32(count),
				})
		}
	}
	return flattenedResults, nil
}

// Callback function type to walk tokens referenced in the expiration
// manager. Don't want to use leaseEntry here because it's an unexported
// type (though most likely we would only call this from within the "vault" core package.)
type ExpirationWalkFunction = func(leaseID string, auth *logical.Auth, path string) bool

var ErrInRestoreMode = errors.New("expiration manager in restore mode")

// WalkTokens extracts the Auth structure from leases corresponding to tokens.
// Returning false from the walk function terminates the iteration.
func (m *ExpirationManager) WalkTokens(walkFn ExpirationWalkFunction) error {
	if m.inRestoreMode() {
		return ErrInRestoreMode
	}

	callback := func(key, value interface{}) bool {
		p := value.(pendingInfo)
		if p.cachedLeaseInfo == nil {
			return true
		}
		lease := p.cachedLeaseInfo
		if lease.Auth != nil {
			return walkFn(key.(string), lease.Auth, lease.Path)
		}
		return true
	}

	m.pending.Range(callback)
	m.nonexpiring.Range(callback)

	return nil
}

// leaseWalkFunction can only be used by the core package.
type leaseWalkFunction = func(leaseID string, expireTime time.Time) bool

func (m *ExpirationManager) walkLeases(walkFn leaseWalkFunction) error {
	if m.inRestoreMode() {
		return ErrInRestoreMode
	}

	callback := func(key, value interface{}) bool {
		p := value.(pendingInfo)
		if p.cachedLeaseInfo == nil {
			return true
		}
		lease := p.cachedLeaseInfo
		expireTime := lease.ExpireTime
		return walkFn(key.(string), expireTime)
	}

	m.pending.Range(callback)
	m.nonexpiring.Range(callback)

	return nil
}

// must be called with m.pendingLock held
// set decrementCounters true to decrement the lease count metric and quota
func (m *ExpirationManager) removeFromPending(ctx context.Context, leaseID string, decrementCounters bool) {
	if info, ok := m.pending.Load(leaseID); ok {
		pending := info.(pendingInfo)
		pending.timer.Stop()
		m.pending.Delete(leaseID)
		if decrementCounters {
			m.leaseCount--
			// Log but do not fail; unit tests (and maybe Tidy on production systems)
			if err := m.core.quotasHandleLeases(ctx, quotas.LeaseActionDeleted, []string{leaseID}); err != nil {
				m.logger.Error("failed to update quota on revocation", "error", err)
			}
		}
	}
}

// Marks a pending lease as irrevocable. Because the lease is being moved from
// pending to irrevocable, no total lease count metrics/quotas updates are needed.
// However, irrevocable lease count will need to be incremented
// note: must be called with pending lock held
func (m *ExpirationManager) markLeaseIrrevocable(ctx context.Context, le *leaseEntry, err error) {
	if le == nil {
		m.logger.Warn("attempted to mark nil lease as irrevocable")
		return
	}
	if le.isIrrevocable() {
		m.logger.Info("attempted to re-mark lease as irrevocable", "original_error", le.RevokeErr, "new_error", err.Error())
		return
	}

	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	if len(errStr) == 0 {
		errStr = genericIrrevocableErrorMessage
	}
	if len(errStr) > maxIrrevocableErrorLength {
		errStr = errStr[:maxIrrevocableErrorLength]
	}

	le.RevokeErr = errStr
	m.persistEntry(ctx, le)

	m.irrevocable.Store(le.LeaseID, m.inMemoryLeaseInfo(le))
	m.irrevocableLeaseCount++
	m.removeFromPending(ctx, le.LeaseID, false)
	m.nonexpiring.Delete(le.LeaseID)
}

func (m *ExpirationManager) getNamespaceFromLeaseID(ctx context.Context, leaseID string) (*namespace.Namespace, error) {
	_, nsID := namespace.SplitIDFromString(leaseID)

	// avoid re-declaring leaseNS and err with scope inside the if
	leaseNS := namespace.RootNamespace
	var err error
	if nsID != "" {
		leaseNS, err = NamespaceByID(ctx, nsID, m.core)
		if err != nil {
			return nil, err
		}
	}

	if leaseNS == nil {
		return nil, namespace.ErrNoNamespace
	}

	return leaseNS, nil
}

func (m *ExpirationManager) getLeaseMountAccessorLocked(ctx context.Context, leaseID string) string {
	m.coreStateLock.RLock()
	defer m.coreStateLock.RUnlock()
	return m.getLeaseMountAccessor(ctx, leaseID)
}

// note: this function must be called with m.coreStateLock held for read
func (m *ExpirationManager) getLeaseMountAccessor(ctx context.Context, leaseID string) string {
	mount := m.core.router.MatchingMountEntry(ctx, leaseID)

	var mountAccessor string
	if mount == nil {
		mountAccessor = "mount-accessor-not-found"
	} else {
		mountAccessor = mount.Accessor
	}

	return mountAccessor
}

func (m *ExpirationManager) getIrrevocableLeaseCounts(ctx context.Context, includeChildNamespaces bool) (map[string]interface{}, error) {
	requestNS, err := namespace.FromContext(ctx)
	if err != nil {
		m.logger.Error("could not get namespace from context", "error", err)
		return nil, err
	}

	numMatchingLeasesPerMount := make(map[string]int)
	numMatchingLeases := 0
	m.irrevocable.Range(func(k, v interface{}) bool {
		leaseID := k.(string)
		leaseNS, err := m.getNamespaceFromLeaseID(ctx, leaseID)
		if err != nil {
			// We should probably note that an error occured, but continue counting
			m.logger.Warn("could not get lease namespace from ID", "error", err)
			return true
		}

		leaseMatches := (leaseNS == requestNS) || (includeChildNamespaces && leaseNS.HasParent(requestNS))
		if !leaseMatches {
			// the lease doesn't meet our criteria, so keep looking
			return true
		}

		mountAccessor := m.getLeaseMountAccessor(ctx, leaseID)

		if _, ok := numMatchingLeasesPerMount[mountAccessor]; !ok {
			numMatchingLeasesPerMount[mountAccessor] = 0
		}

		numMatchingLeases++
		numMatchingLeasesPerMount[mountAccessor]++

		return true
	})

	resp := make(map[string]interface{})
	resp["lease_count"] = numMatchingLeases
	resp["counts"] = numMatchingLeasesPerMount

	return resp, nil
}

type leaseResponse struct {
	LeaseID    string `json:"lease_id"`
	MountID    string `json:"mount_id"`
	ErrMsg     string `json:"error"`
	expireTime time.Time
}

// returns a warning string, if applicable
// limit specifies how many results to return, and must be >0
// includeAll specifies if all results should be returned, regardless of limit
func (m *ExpirationManager) listIrrevocableLeases(ctx context.Context, includeChildNamespaces, returnAll bool, limit int) (map[string]interface{}, string, error) {
	requestNS, err := namespace.FromContext(ctx)
	if err != nil {
		m.logger.Error("could not get namespace from context", "error", err)
		return nil, "", err
	}

	// map of mount point : lease info
	matchingLeases := make([]*leaseResponse, 0)
	numMatchingLeases := 0
	var warning string
	m.irrevocable.Range(func(k, v interface{}) bool {
		leaseID := k.(string)
		leaseInfo := v.(*leaseEntry)

		leaseNS, err := m.getNamespaceFromLeaseID(ctx, leaseID)
		if err != nil {
			// We probably want to track that an error occured, but continue counting
			m.logger.Warn("could not get lease namespace from ID", "error", err)
			return true
		}

		leaseMatches := (leaseNS == requestNS) || (includeChildNamespaces && leaseNS.HasParent(requestNS))
		if !leaseMatches {
			// the lease doesn't meet our criteria, so keep looking
			return true
		}

		if !returnAll && (numMatchingLeases >= limit) {
			m.logger.Warn("hit max irrevocable leases without force flag set")
			warning = MaxIrrevocableLeasesWarning
			return false
		}

		mountAccessor := m.getLeaseMountAccessor(ctx, leaseID)

		numMatchingLeases++
		matchingLeases = append(matchingLeases, &leaseResponse{
			LeaseID:    leaseID,
			MountID:    mountAccessor,
			ErrMsg:     leaseInfo.RevokeErr,
			expireTime: leaseInfo.ExpireTime,
		})

		return true
	})

	// sort the results for consistent API response. we primarily sort on
	// increasing expire time, and break ties with increasing lease id
	sort.Slice(matchingLeases, func(i, j int) bool {
		if !matchingLeases[i].expireTime.Equal(matchingLeases[j].expireTime) {
			return matchingLeases[i].expireTime.Before(matchingLeases[j].expireTime)
		}

		return matchingLeases[i].LeaseID < matchingLeases[j].LeaseID
	})

	resp := make(map[string]interface{})
	resp["lease_count"] = numMatchingLeases
	resp["leases"] = matchingLeases

	return resp, warning, nil
}

// leaseEntry is used to structure the values the expiration
// manager stores. This is used to handle renew and revocation.
type leaseEntry struct {
	LeaseID         string                 `json:"lease_id"`
	ClientToken     string                 `json:"client_token"`
	ClientTokenType logical.TokenType      `json:"token_type"`
	Path            string                 `json:"path"`
	Data            map[string]interface{} `json:"data"`
	Secret          *logical.Secret        `json:"secret"`
	Auth            *logical.Auth          `json:"auth"`
	IssueTime       time.Time              `json:"issue_time"`
	ExpireTime      time.Time              `json:"expire_time"`
	LastRenewalTime time.Time              `json:"last_renewal_time"`

	// Version is used to track new different versions of leases. V0 (or
	// zero-value) had non-root namespaced secondary indexes live in the root
	// namespace, and V1 has secondary indexes live in the matching namespace.
	Version int `json:"version"`

	namespace *namespace.Namespace

	// RevokeErr tracks if a lease has failed revocation in a way that is
	// unlikely to be automatically resolved. The first time this happens,
	// RevokeErr will be set, thus marking this leaseEntry as irrevocable. From
	// there, it must be manually removed (force revoked).
	RevokeErr string `json:"revokeErr"`
}

// encode is used to JSON encode the lease entry
func (le *leaseEntry) encode() ([]byte, error) {
	return json.Marshal(le)
}

func (le *leaseEntry) renewable() (bool, error) {
	switch {
	// If there is no entry, cannot review to renew
	case le == nil:
		return false, fmt.Errorf("lease not found")

	case le.isIrrevocable():
		return false, fmt.Errorf("lease is expired and has failed previous revocation attempts")

	case le.ExpireTime.IsZero():
		return false, fmt.Errorf("lease is not renewable")

	case le.ClientTokenType == logical.TokenTypeBatch:
		return false, nil

	// Determine if the lease is expired
	case le.ExpireTime.Before(time.Now()):
		return false, fmt.Errorf("lease expired")

	// Determine if the lease is renewable
	case le.Secret != nil && !le.Secret.Renewable:
		return false, fmt.Errorf("lease is not renewable")

	case le.Auth != nil && !le.Auth.Renewable:
		return false, fmt.Errorf("lease is not renewable")
	}

	return true, nil
}

func (le *leaseEntry) ttl() int64 {
	return int64(le.ExpireTime.Sub(time.Now().Round(time.Second)).Seconds())
}

func (le *leaseEntry) nonexpiringToken() bool {
	if le.Auth == nil {
		return false
	}
	// Note that at this time the only non-expiring tokens are root tokens, this test is more involved as it is trying
	// to catch tokens created by the VAULT-1949 non-expiring tokens bug and ensure they become expiring.
	return !le.Auth.LeaseEnabled() && len(le.Auth.Policies) == 1 && le.Auth.Policies[0] == "root" && le.namespace != nil &&
		le.namespace.ID == namespace.RootNamespaceID
}

// TODO maybe lock RevokeErr once this goes in: https://github.com/hashicorp/vault/pull/11122
func (le *leaseEntry) isIrrevocable() bool {
	return le.RevokeErr != ""
}

func (le *leaseEntry) isIncorrectlyNonExpiring() bool {
	return le.ExpireTime.IsZero() && !le.nonexpiringToken()
}

// decodeLeaseEntry is used to reverse encode and return a new entry
func decodeLeaseEntry(buf []byte) (*leaseEntry, error) {
	out := new(leaseEntry)
	return out, jsonutil.DecodeJSON(buf, out)
}
