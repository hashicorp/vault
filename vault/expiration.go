package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/base62"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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

	//maxLeaseThreshold is the maximum lease count before generating log warning
	maxLeaseThreshold = 256000
)

type pendingInfo struct {
	exportLeaseTimes *leaseEntry
	timer            *time.Timer
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

	pending     map[string]pendingInfo
	pendingLock sync.RWMutex

	tidyLock *int32

	restoreMode        *int32
	restoreModeLock    sync.RWMutex
	restoreRequestLock sync.RWMutex
	restoreLocks       []*locksutil.LockEntry
	restoreLoaded      sync.Map
	quitCh             chan struct{}

	coreStateLock     *sync.RWMutex
	quitContext       context.Context
	leaseCheckCounter *uint32

	logLeaseExpirations bool
	expireFunc          ExpireLeaseStrategy
}

type ExpireLeaseStrategy func(context.Context, *ExpirationManager, *leaseEntry)

// revokeIDFunc is invoked when a given ID is expired
func expireLeaseStrategyRevoke(ctx context.Context, m *ExpirationManager, le *leaseEntry) {
	for attempt := uint(0); attempt < maxRevokeAttempts; attempt++ {
		revokeCtx, cancel := context.WithTimeout(ctx, DefaultMaxRequestDuration)
		revokeCtx = namespace.ContextWithNamespace(revokeCtx, le.namespace)

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
			m.logger.Error("shutting down, not attempting further revocation of lease", "lease_id", le.LeaseID)
			cancel()
			return
		case <-m.quitContext.Done():
			m.logger.Error("core context canceled, not attempting further revocation of lease", "lease_id", le.LeaseID)
			cancel()
			return
		default:
		}

		m.coreStateLock.RLock()
		err := m.Revoke(revokeCtx, le.LeaseID)
		m.coreStateLock.RUnlock()
		cancel()
		if err == nil {
			return
		}

		m.logger.Error("failed to revoke lease", "lease_id", le.LeaseID, "error", err)
		time.Sleep((1 << attempt) * revokeRetryBase)
	}
	m.logger.Error("maximum revoke attempts reached", "lease_id", le.LeaseID)
}

// NewExpirationManager creates a new ExpirationManager that is backed
// using a given view, and uses the provided router for revocation.
func NewExpirationManager(c *Core, view *BarrierView, e ExpireLeaseStrategy, logger log.Logger) *ExpirationManager {
	exp := &ExpirationManager{
		core:       c,
		router:     c.router,
		idView:     view.SubView(leaseViewPrefix),
		tokenView:  view.SubView(tokenViewPrefix),
		tokenStore: c.tokenStore,
		logger:     logger,
		pending:    make(map[string]pendingInfo),
		tidyLock:   new(int32),

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
	}
	*exp.restoreMode = 1

	if exp.logger == nil {
		opts := log.LoggerOptions{Name: "expiration_manager"}
		exp.logger = log.New(&opts)
	}

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
		// Clear from the pending expiration
		leaseID := strings.TrimPrefix(key, leaseViewPrefix)
		m.pendingLock.Lock()
		if pending, ok := m.pending[leaseID]; ok {
			pending.timer.Stop()
			delete(m.pending, leaseID)
		}
		m.pendingLock.Unlock()
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
			tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf(fmt.Sprintf("failed to load the lease ID %q: {{err}}", leaseID), err))
			return
		}

		if le == nil {
			tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf(fmt.Sprintf("nil entry for lease ID %q: {{err}}", leaseID), err))
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
				tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf("failed to lookup token: {{err}}", err))
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
			err = m.revokeCommon(ctx, leaseID, true, true)
			if err != nil {
				tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf(fmt.Sprintf("failed to revoke an invalid lease with ID %q: {{err}}", leaseID), err))
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
			m.logger.Warn("context cancled while restoring leases, stopping lease loading")
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

	// Do this before stopping pending timers to avoid potential races with
	// expiring timers
	close(m.quitCh)

	m.pendingLock.Lock()
	for _, pending := range m.pending {
		pending.timer.Stop()
	}
	m.pending = make(map[string]pendingInfo)
	m.pendingLock.Unlock()

	if m.inRestoreMode() {
		for {
			if !m.inRestoreMode() {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}

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
	{
		m.pendingLock.Lock()
		if err := m.persistEntry(ctx, le); err != nil {
			m.pendingLock.Unlock()
			return err
		}

		m.updatePendingInternal(le, 0)
		m.pendingLock.Unlock()
	}

	return nil
}

// revokeCommon does the heavy lifting. If force is true, we ignore a problem
// during revocation and still remove entries/index/lease timers
func (m *ExpirationManager) revokeCommon(ctx context.Context, leaseID string, force, skipToken bool) error {
	defer metrics.MeasureSince([]string{"expire", "revoke-common"}, time.Now())

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

	// Delete the secondary index, but only if it's a leased secret (not auth)
	if le.Secret != nil {
		if err := m.removeIndexByToken(ctx, le); err != nil {
			return err
		}
	}

	// Clear the expiration handler
	m.pendingLock.Lock()
	if pending, ok := m.pending[leaseID]; ok {
		pending.timer.Stop()
		delete(m.pending, leaseID)
	}
	m.pendingLock.Unlock()

	if m.logger.IsInfo() && !skipToken && m.logLeaseExpirations {
		m.logger.Info("revoked lease", "lease_id", leaseID)
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
// token store's revokeSalted function.
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
		return errwrap.Wrapf("failed to scan for leases: {{err}}", err)
	}

	// Revoke all the keys
	for _, leaseID := range existing {
		// Load the entry
		le, err := m.loadEntry(ctx, leaseID)
		if err != nil {
			return err
		}

		// If there's a lease, set expiration to now, persist, and call
		// updatePending to hand off revocation to the expiration manager's pending
		// timer map
		if le != nil {
			le.ExpireTime = time.Now()

			{
				m.pendingLock.Lock()
				if err := m.persistEntry(ctx, le); err != nil {
					m.pendingLock.Unlock()
					return err
				}

				m.updatePendingInternal(le, 0)
				m.pendingLock.Unlock()
			}
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
					return errwrap.Wrapf(fmt.Sprintf("failed to revoke %q: {{err}}", prefix), err)
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
		return errwrap.Wrapf("failed to scan for leases: {{err}}", err)
	}

	// Revoke all the keys
	for idx, suffix := range existing {
		leaseID := prefix + suffix
		switch {
		case sync:
			if err := m.revokeCommon(ctx, leaseID, force, false); err != nil {
				return errwrap.Wrapf(fmt.Sprintf("failed to revoke %q (%d / %d): {{err}}", leaseID, idx+1, len(existing)), err)
			}
		default:
			if err := m.LazyRevoke(ctx, leaseID); err != nil {
				return errwrap.Wrapf(fmt.Sprintf("failed to revoke %q (%d / %d): {{err}}", leaseID, idx+1, len(existing)), err)
			}
		}
	}

	return nil
}

// Renew is used to renew a secret using the given leaseID
// and a renew interval. The increment may be ignored.
func (m *ExpirationManager) Renew(ctx context.Context, leaseID string, increment time.Duration) (*logical.Response, error) {
	defer metrics.MeasureSince([]string{"expire", "renew"}, time.Now())

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
			return logical.ErrorResponse("tokens cannot be renewed through this endpoint"), logical.ErrPermissionDenied
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
		if le.ExpireTime.After(tokenLeaseTimes.ExpireTime) {
			resp.Secret.TTL = tokenLeaseTimes.ExpireTime.Sub(le.LastRenewalTime)
			le.ExpireTime = tokenLeaseTimes.ExpireTime
		}
	}

	{
		m.pendingLock.Lock()
		if err := m.persistEntry(ctx, le); err != nil {
			m.pendingLock.Unlock()
			return nil, err
		}

		// Update the expiration time
		m.updatePendingInternal(le, resp.Secret.LeaseTotal())
		m.pendingLock.Unlock()
	}

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
	if resp.Auth.EntityID != "" &&
		resp.Auth.GroupAliases != nil &&
		m.core.identityStore != nil {
		validAliases, err := m.core.identityStore.refreshExternalGroupMembershipsByEntityID(resp.Auth.EntityID, resp.Auth.GroupAliases)
		if err != nil {
			return nil, err
		}
		resp.Auth.GroupAliases = validAliases
	}

	// Update the lease entry
	le.Auth = resp.Auth
	le.ExpireTime = resp.Auth.ExpirationTime()
	le.LastRenewalTime = time.Now()

	{
		m.pendingLock.Lock()
		if err := m.persistEntry(ctx, le); err != nil {
			m.pendingLock.Unlock()
			return nil, err
		}

		// Update the expiration time
		m.updatePendingInternal(le, resp.Auth.LeaseTotal())
		m.pendingLock.Unlock()
	}

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
				retErr = multierror.Append(retErr, errwrap.Wrapf("an additional internal error was encountered revoking the newly-generated secret: {{err}}", err))
			} else if revResp != nil && revResp.IsError() {
				retErr = multierror.Append(retErr, errwrap.Wrapf("an additional error was encountered revoking the newly-generated secret: {{err}}", revResp.Error()))
			}

			if err := m.deleteEntry(ctx, le); err != nil {
				retErr = multierror.Append(retErr, errwrap.Wrapf("an additional error was encountered deleting any lease associated with the newly-generated secret: {{err}}", err))
			}

			if err := m.removeIndexByToken(ctx, le); err != nil {
				retErr = multierror.Append(retErr, errwrap.Wrapf("an additional error was encountered removing lease indexes associated with the newly-generated secret: {{err}}", err))
			}
		}
	}()

	// If the token is a batch token, we want to constrain the maximum lifetime
	// by the token's lifetime
	if te.Type == logical.TokenTypeBatch {
		tokenLeaseTimes, err := m.FetchLeaseTimesByToken(ctx, te)
		if err != nil {
			return "", err
		}
		if le.ExpireTime.After(tokenLeaseTimes.ExpireTime) {
			le.ExpireTime = tokenLeaseTimes.ExpireTime
		}
	}

	// Encode the entry
	if err := m.persistEntry(ctx, le); err != nil {
		return "", err
	}

	// Maintain secondary index by token, except for orphan batch tokens
	switch {
	case te.Type != logical.TokenTypeBatch:
		if err := m.createIndexByToken(ctx, le, le.ClientToken); err != nil {
			return "", err
		}
	case te.Parent != "":
		// If it's a non-orphan batch token, assign the secondary index to its
		// parent
		if err := m.createIndexByToken(ctx, le, te.Parent); err != nil {
			return "", err
		}
	}

	// Setup revocation timer if there is a lease
	m.updatePending(le, resp.Secret.LeaseTotal())

	// Done
	return le.LeaseID, nil
}

// RegisterAuth is used to take an Auth response with an associated lease.
// The token does not get a LeaseID, but the lease management is handled by
// the expiration manager.
func (m *ExpirationManager) RegisterAuth(ctx context.Context, te *logical.TokenEntry, auth *logical.Auth) error {
	defer metrics.MeasureSince([]string{"expire", "register-auth"}, time.Now())

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
		ExpireTime:  auth.ExpirationTime(),
		namespace:   tokenNS,
	}

	// Encode the entry
	if err := m.persistEntry(ctx, &le); err != nil {
		return err
	}

	// Setup revocation timer
	m.updatePending(&le, auth.LeaseTotal())

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

	m.pendingLock.RLock()
	val := m.pending[leaseID]
	m.pendingLock.RUnlock()

	if val.exportLeaseTimes != nil {
		return val.exportLeaseTimes, nil
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

// updatePending is used to update a pending invocation for a lease
func (m *ExpirationManager) updatePending(le *leaseEntry, leaseTotal time.Duration) {
	m.pendingLock.Lock()
	defer m.pendingLock.Unlock()

	m.updatePendingInternal(le, leaseTotal)
}

// updatePendingInternal is the locked version of updatePending; do not call
// this without a write lock on m.pending
func (m *ExpirationManager) updatePendingInternal(le *leaseEntry, leaseTotal time.Duration) {
	// Check for an existing timer
	pending, ok := m.pending[le.LeaseID]

	// If there is no expiry time, don't do anything
	if le.ExpireTime.IsZero() {
		// if the timer happened to exist, stop the time and delete it from the
		// pending timers.
		if ok {
			pending.timer.Stop()
			delete(m.pending, le.LeaseID)
		}
		return
	}

	// Create entry if it does not exist or reset if it does
	if ok {
		pending.timer.Reset(leaseTotal)
	} else {
		timer := time.AfterFunc(leaseTotal, func() {
			m.expireFunc(m.quitContext, m, le)
		})
		pending = pendingInfo{
			timer: timer,
		}
	}

	// Extend the timer by the lease total
	pending.exportLeaseTimes = m.leaseTimesForExport(le)

	m.pending[le.LeaseID] = pending
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
			return errwrap.Wrapf("failed to revoke token: {{err}}", err)
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
		return errwrap.Wrapf(fmt.Sprintf("failed to revoke entry: resp: %#v err: {{err}}", resp), err)
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
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to renew entry: resp: %#v err: {{err}}", resp), err)
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
		return nil, errwrap.Wrapf("failed to renew entry: {{err}}", err)
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
	return m.loadEntryInternal(ctx, leaseID, restoreMode, true)
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
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to read lease entry %s: {{err}}", leaseID), err)
	}
	if out == nil {
		return nil, nil
	}
	le, err := decodeLeaseEntry(out.Value)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to decode lease entry %s: {{err}}", leaseID), err)
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
		m.updatePending(le, le.ExpireTime.Sub(time.Now()))
	}
	return le, nil
}

// persistEntry is used to persist a lease entry
func (m *ExpirationManager) persistEntry(ctx context.Context, le *leaseEntry) error {
	// Encode the entry
	buf, err := le.encode()
	if err != nil {
		return errwrap.Wrapf("failed to encode lease entry: {{err}}", err)
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
		return errwrap.Wrapf("failed to persist lease entry: {{err}}", err)
	}
	return nil
}

// deleteEntry is used to delete a lease entry
func (m *ExpirationManager) deleteEntry(ctx context.Context, le *leaseEntry) error {
	view := m.leaseView(le.namespace)
	if err := view.Delete(ctx, le.LeaseID); err != nil {
		return errwrap.Wrapf("failed to delete lease entry: {{err}}", err)
	}
	return nil
}

// createIndexByToken creates a secondary index from the token to a lease entry
func (m *ExpirationManager) createIndexByToken(ctx context.Context, le *leaseEntry, token string) error {
	tokenNS := namespace.RootNamespace
	saltCtx := namespace.ContextWithNamespace(ctx, namespace.RootNamespace)
	_, nsID := namespace.SplitIDFromString(token)
	if nsID != "" {
		tokenNS, err := NamespaceByID(ctx, nsID, m.core)
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
		return errwrap.Wrapf("failed to persist lease index entry: {{err}}", err)
	}
	return nil
}

// indexByToken looks up the secondary index from the token to a lease entry
func (m *ExpirationManager) indexByToken(ctx context.Context, le *leaseEntry) (*logical.StorageEntry, error) {
	tokenNS := namespace.RootNamespace
	saltCtx := namespace.ContextWithNamespace(ctx, tokenNS)
	_, nsID := namespace.SplitIDFromString(le.ClientToken)
	if nsID != "" {
		tokenNS, err := NamespaceByID(ctx, nsID, m.core)
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
func (m *ExpirationManager) removeIndexByToken(ctx context.Context, le *leaseEntry) error {
	tokenNS := namespace.RootNamespace
	saltCtx := namespace.ContextWithNamespace(ctx, namespace.RootNamespace)
	_, nsID := namespace.SplitIDFromString(le.ClientToken)
	if nsID != "" {
		tokenNS, err := NamespaceByID(ctx, nsID, m.core)
		if err != nil {
			return err
		}
		if tokenNS != nil {
			saltCtx = namespace.ContextWithNamespace(ctx, tokenNS)
		}
	}

	saltedID, err := m.tokenStore.SaltID(saltCtx, le.ClientToken)
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
		return errwrap.Wrapf("failed to delete lease index entry: {{err}}", err)
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
		}

		// Encode the entry
		if err := m.persistEntry(ctx, le); err != nil {
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
		return nil, errwrap.Wrapf("failed to list leases: {{err}}", err)
	}

	// Read each index entry
	leaseIDs := make([]string, 0, len(subKeys))
	for _, sub := range subKeys {
		out, err := tokenView.Get(ctx, prefix+sub)
		if err != nil {
			return nil, errwrap.Wrapf("failed to read lease index: {{err}}", err)
		}
		if out == nil {
			continue
		}
		leaseIDs = append(leaseIDs, string(out.Value))
	}
	return leaseIDs, nil
}

// emitMetrics is invoked periodically to emit statistics
func (m *ExpirationManager) emitMetrics() {
	m.pendingLock.RLock()
	num := len(m.pending)
	m.pendingLock.RUnlock()
	metrics.SetGauge([]string{"expire", "num_leases"}, float32(num))
	// Check if lease count is greater than the threshold
	if num > maxLeaseThreshold {
		if atomic.LoadUint32(m.leaseCheckCounter) > 59 {
			m.logger.Warn("lease count exceeds warning lease threshold")
			atomic.StoreUint32(m.leaseCheckCounter, 0)
		} else {
			atomic.AddUint32(m.leaseCheckCounter, 1)
		}
	}
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

	namespace *namespace.Namespace
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

// decodeLeaseEntry is used to reverse encode and return a new entry
func decodeLeaseEntry(buf []byte) (*leaseEntry, error) {
	out := new(leaseEntry)
	return out, jsonutil.DecodeJSON(buf, out)
}
