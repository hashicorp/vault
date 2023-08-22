// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/oklog/run"
)

const (
	// lockRetryInterval is the interval we re-attempt to acquire the
	// HA lock if an error is encountered
	lockRetryInterval = 10 * time.Second

	// leaderCheckInterval is how often a standby checks for a new leader
	leaderCheckInterval = 2500 * time.Millisecond

	// keyRotateCheckInterval is how often a standby checks for a key
	// rotation taking place.
	keyRotateCheckInterval = 10 * time.Second

	// leaderPrefixCleanDelay is how long to wait between deletions
	// of orphaned leader keys, to prevent slamming the backend.
	leaderPrefixCleanDelay = 200 * time.Millisecond
)

var (
	addEnterpriseHaActors func(*Core, *run.Group) chan func()            = addEnterpriseHaActorsNoop
	interruptPerfStandby  func(chan func(), chan struct{}) chan struct{} = interruptPerfStandbyNoop
)

func addEnterpriseHaActorsNoop(*Core, *run.Group) chan func() { return nil }
func interruptPerfStandbyNoop(chan func(), chan struct{}) chan struct{} {
	return make(chan struct{})
}

// Standby checks if the Vault is in standby mode
func (c *Core) Standby() (bool, error) {
	c.stateLock.RLock()
	standby := c.standby
	c.stateLock.RUnlock()
	return standby, nil
}

// PerfStandby checks if the vault is a performance standby
// This function cannot be used during request handling
// because this causes a deadlock with the statelock.
func (c *Core) PerfStandby() bool {
	c.stateLock.RLock()
	perfStandby := c.perfStandby
	c.stateLock.RUnlock()
	return perfStandby
}

func (c *Core) ActiveTime() time.Time {
	c.stateLock.RLock()
	activeTime := c.activeTime
	c.stateLock.RUnlock()
	return activeTime
}

// StandbyStates is meant as a way to avoid some extra locking on the very
// common sys/health check.
func (c *Core) StandbyStates() (standby, perfStandby bool) {
	c.stateLock.RLock()
	standby = c.standby
	perfStandby = c.perfStandby
	c.stateLock.RUnlock()
	return
}

// getHAMembers retrieves cluster membership that doesn't depend on raft. This should only ever be called by the
// active node.
func (c *Core) getHAMembers() ([]HAStatusNode, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	leader := HAStatusNode{
		Hostname:       hostname,
		APIAddress:     c.redirectAddr,
		ClusterAddress: c.ClusterAddr(),
		ActiveNode:     true,
		Version:        c.effectiveSDKVersion,
	}

	if rb := c.getRaftBackend(); rb != nil {
		leader.UpgradeVersion = rb.EffectiveVersion()
		leader.RedundancyZone = rb.RedundancyZone()
	}

	nodes := []HAStatusNode{leader}

	for _, peerNode := range c.GetHAPeerNodesCached() {
		lastEcho := peerNode.LastEcho
		nodes = append(nodes, HAStatusNode{
			Hostname:       peerNode.Hostname,
			APIAddress:     peerNode.APIAddress,
			ClusterAddress: peerNode.ClusterAddress,
			LastEcho:       &lastEcho,
			Version:        peerNode.Version,
			UpgradeVersion: peerNode.UpgradeVersion,
			RedundancyZone: peerNode.RedundancyZone,
		})
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].APIAddress < nodes[j].APIAddress
	})

	return nodes, nil
}

// Leader is used to get information about the current active leader in relation to the current node (core).
// It utilizes a state lock on the Core by attempting to acquire a read lock. Care should be taken not to
// call this method if a read lock on this Core's state lock is currently held, as this can cause deadlock.
// e.g. if called from within request handling.
func (c *Core) Leader() (isLeader bool, leaderAddr, clusterAddr string, err error) {
	// Check if HA enabled. We don't need the lock for this check as it's set
	// on startup and never modified
	if c.ha == nil {
		return false, "", "", ErrHANotEnabled
	}

	// Check if sealed
	if c.Sealed() {
		return false, "", "", consts.ErrSealed
	}

	c.stateLock.RLock()

	// Check if we are the leader
	if !c.standby {
		c.stateLock.RUnlock()
		return true, c.redirectAddr, c.ClusterAddr(), nil
	}

	// Initialize a lock
	lock, err := c.ha.LockWith(CoreLockPath, "read")
	if err != nil {
		c.stateLock.RUnlock()
		return false, "", "", err
	}

	// Read the value
	held, leaderUUID, err := lock.Value()
	if err != nil {
		c.stateLock.RUnlock()
		return false, "", "", err
	}
	if !held {
		c.stateLock.RUnlock()
		return false, "", "", nil
	}

	var localLeaderUUID, localRedirectAddr, localClusterAddr string
	clusterLeaderParams := c.clusterLeaderParams.Load().(*ClusterLeaderParams)
	if clusterLeaderParams != nil {
		localLeaderUUID = clusterLeaderParams.LeaderUUID
		localRedirectAddr = clusterLeaderParams.LeaderRedirectAddr
		localClusterAddr = clusterLeaderParams.LeaderClusterAddr
	}

	// If the leader hasn't changed, return the cached value; nothing changes
	// mid-leadership, and the barrier caches anyways
	if leaderUUID == localLeaderUUID && localRedirectAddr != "" {
		c.stateLock.RUnlock()
		return false, localRedirectAddr, localClusterAddr, nil
	}

	c.logger.Trace("found new active node information, refreshing")

	defer c.stateLock.RUnlock()
	c.leaderParamsLock.Lock()
	defer c.leaderParamsLock.Unlock()

	// Validate base conditions again
	clusterLeaderParams = c.clusterLeaderParams.Load().(*ClusterLeaderParams)
	if clusterLeaderParams != nil {
		localLeaderUUID = clusterLeaderParams.LeaderUUID
		localRedirectAddr = clusterLeaderParams.LeaderRedirectAddr
		localClusterAddr = clusterLeaderParams.LeaderClusterAddr
	} else {
		localLeaderUUID = ""
		localRedirectAddr = ""
		localClusterAddr = ""
	}

	if leaderUUID == localLeaderUUID && localRedirectAddr != "" {
		return false, localRedirectAddr, localClusterAddr, nil
	}

	key := coreLeaderPrefix + leaderUUID
	// Use background because postUnseal isn't run on standby
	entry, err := c.barrier.Get(context.Background(), key)
	if err != nil {
		return false, "", "", err
	}
	if entry == nil {
		return false, "", "", nil
	}

	var oldAdv bool

	var adv activeAdvertisement
	err = jsonutil.DecodeJSON(entry.Value, &adv)
	if err != nil {
		// Fall back to pre-struct handling
		adv.RedirectAddr = string(entry.Value)
		c.logger.Debug("parsed redirect addr for new active node", "redirect_addr", adv.RedirectAddr)
		oldAdv = true
	}

	// At the top of this function we return early when we're the active node.
	// If we're not the active node, and there's a stale advertisement pointing
	// to ourself, there's no point in paying any attention to it.  And by
	// disregarding it, we can avoid a panic in raft tests using the Inmem network
	// layer when we try to connect back to ourself.
	if adv.ClusterAddr == c.ClusterAddr() && adv.RedirectAddr == c.redirectAddr && c.getRaftBackend() != nil {
		return false, "", "", nil
	}

	if !oldAdv {
		c.logger.Debug("parsing information for new active node", "active_cluster_addr", adv.ClusterAddr, "active_redirect_addr", adv.RedirectAddr)

		// Ensure we are using current values
		err = c.loadLocalClusterTLS(adv)
		if err != nil {
			return false, "", "", err
		}

		// This will ensure that we both have a connection at the ready and that
		// the address is the current known value
		// Since this is standby, we don't use the active context. Later we may
		// use a process-scoped context
		err = c.refreshRequestForwardingConnection(context.Background(), adv.ClusterAddr)
		if err != nil {
			return false, "", "", err
		}
	}

	// Don't set these until everything has been parsed successfully or we'll
	// never try again
	c.clusterLeaderParams.Store(&ClusterLeaderParams{
		LeaderUUID:         leaderUUID,
		LeaderRedirectAddr: adv.RedirectAddr,
		LeaderClusterAddr:  adv.ClusterAddr,
	})

	return false, adv.RedirectAddr, adv.ClusterAddr, nil
}

// StepDown is used to step down from leadership
func (c *Core) StepDown(httpCtx context.Context, req *logical.Request) (retErr error) {
	defer metrics.MeasureSince([]string{"core", "step_down"}, time.Now())

	if req == nil {
		return errors.New("nil request to step-down")
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()

	if c.Sealed() {
		return nil
	}
	if c.ha == nil || c.standby {
		return nil
	}

	ctx, cancel := context.WithCancel(namespace.RootContext(nil))
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case <-httpCtx.Done():
			cancel()
		}
	}()

	err := c.PopulateTokenEntry(ctx, req)
	if err != nil {
		if errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
			return logical.ErrPermissionDenied
		}
		return logical.ErrInvalidRequest
	}
	acl, te, entity, identityPolicies, err := c.fetchACLTokenEntryAndEntity(ctx, req)
	if err != nil {
		return err
	}

	// Audit-log the request before going any further
	auth := &logical.Auth{
		ClientToken: req.ClientToken,
		Accessor:    req.ClientTokenAccessor,
	}
	if te != nil {
		auth.IdentityPolicies = identityPolicies[te.NamespaceID]
		delete(identityPolicies, te.NamespaceID)
		auth.ExternalNamespacePolicies = identityPolicies
		auth.TokenPolicies = te.Policies
		auth.Policies = append(te.Policies, identityPolicies[te.NamespaceID]...)
		auth.Metadata = te.Meta
		auth.DisplayName = te.DisplayName
		auth.EntityID = te.EntityID
		auth.TokenType = te.Type
	}

	logInput := &logical.LogInput{
		Auth:    auth,
		Request: req,
	}
	if err := c.auditBroker.LogRequest(ctx, logInput, c.auditedHeaders); err != nil {
		c.logger.Error("failed to audit request", "request_path", req.Path, "error", err)
		return errors.New("failed to audit request, cannot continue")
	}

	if entity != nil && entity.Disabled {
		c.logger.Warn("permission denied as the entity on the token is disabled")
		return logical.ErrPermissionDenied
	}

	if te != nil && te.EntityID != "" && entity == nil {
		c.logger.Warn("permission denied as the entity on the token is invalid")
		return logical.ErrPermissionDenied
	}

	// Attempt to use the token (decrement num_uses)
	if te != nil {
		te, err = c.tokenStore.UseToken(ctx, te)
		if err != nil {
			c.logger.Error("failed to use token", "error", err)
			return ErrInternalError
		}
		if te == nil {
			// Token has been revoked
			return logical.ErrPermissionDenied
		}
	}

	// Verify that this operation is allowed
	authResults := c.performPolicyChecks(ctx, acl, te, req, entity, &PolicyCheckOpts{
		RootPrivsRequired: true,
	})
	if !authResults.Allowed {
		retErr = multierror.Append(retErr, authResults.Error)
		if authResults.Error.ErrorOrNil() == nil || authResults.DeniedError {
			retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
		}
		return retErr
	}

	if te != nil && te.NumUses == tokenRevocationPending {
		// Token needs to be revoked. We do this immediately here because
		// we won't have a token store after sealing.
		leaseID, err := c.expiration.CreateOrFetchRevocationLeaseByToken(c.activeContext, te)
		if err == nil {
			err = c.expiration.Revoke(c.activeContext, leaseID)
		}
		if err != nil {
			c.logger.Error("token needed revocation before step-down but failed to revoke", "error", err)
			retErr = multierror.Append(retErr, ErrInternalError)
		}
	}

	select {
	case c.manualStepDownCh <- struct{}{}:
	default:
		c.logger.Warn("manual step-down operation already queued")
	}

	return retErr
}

// runStandby is a long running process that manages a number of the HA
// subsystems.
func (c *Core) runStandby(doneCh, manualStepDownCh, stopCh chan struct{}) {
	defer close(doneCh)
	defer close(manualStepDownCh)
	c.logger.Info("entering standby mode")

	var g run.Group
	newLeaderCh := addEnterpriseHaActors(c, &g)
	{
		// This will cause all the other actors to close when the stop channel
		// is closed.
		g.Add(func() error {
			<-stopCh
			return nil
		}, func(error) {})
	}
	{
		// Monitor for key rotations
		keyRotateStop := make(chan struct{})

		g.Add(func() error {
			c.periodicCheckKeyUpgrades(context.Background(), keyRotateStop)
			return nil
		}, func(error) {
			close(keyRotateStop)
			c.logger.Debug("shutting down periodic key rotation checker")
		})
	}
	{
		// Monitor for new leadership
		checkLeaderStop := make(chan struct{})

		g.Add(func() error {
			c.periodicLeaderRefresh(newLeaderCh, checkLeaderStop)
			return nil
		}, func(error) {
			close(checkLeaderStop)
			c.logger.Debug("shutting down periodic leader refresh")
		})
	}
	{
		metricsStop := make(chan struct{})

		g.Add(func() error {
			c.metricsLoop(metricsStop)
			return nil
		}, func(error) {
			close(metricsStop)
			c.logger.Debug("shutting down periodic metrics")
		})
	}
	{
		// Wait for leadership
		leaderStopCh := make(chan struct{})

		g.Add(func() error {
			c.waitForLeadership(newLeaderCh, manualStepDownCh, leaderStopCh)
			return nil
		}, func(error) {
			close(leaderStopCh)
			c.logger.Debug("shutting down leader elections")
		})
	}

	// Start all the actors
	g.Run()
}

// waitForLeadership is a long running routine that is used when an HA backend
// is enabled. It waits until we are leader and switches this Vault to
// active.
func (c *Core) waitForLeadership(newLeaderCh chan func(), manualStepDownCh, stopCh chan struct{}) {
	var manualStepDown bool
	firstIteration := true
	for {
		// Check for a shutdown
		select {
		case <-stopCh:
			c.logger.Debug("stop channel triggered in runStandby")
			return
		default:
			// If we've just down, we could instantly grab the lock again. Give
			// the other nodes a chance.
			if manualStepDown {
				time.Sleep(manualStepDownSleepPeriod)
				manualStepDown = false
			} else if !firstIteration {
				// If we restarted the for loop due to an error, wait a second
				// so that we don't busy loop if the error persists.
				time.Sleep(1 * time.Second)
			}
		}
		firstIteration = false

		// Create a lock
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			c.logger.Error("failed to generate uuid", "error", err)
			continue
		}
		lock, err := c.ha.LockWith(CoreLockPath, uuid)
		if err != nil {
			c.logger.Error("failed to create lock", "error", err)
			continue
		}

		// Attempt the acquisition
		leaderLostCh := c.acquireLock(lock, stopCh)

		// Bail if we are being shutdown
		if leaderLostCh == nil {
			return
		}

		if atomic.LoadUint32(c.neverBecomeActive) == 1 {
			c.heldHALock = nil
			lock.Unlock()
			c.logger.Info("marked never become active, giving up active state")
			continue
		}

		c.logger.Info("acquired lock, enabling active operation")

		// This is used later to log a metrics event; this can be helpful to
		// detect flapping
		activeTime := time.Now()

		continueCh := interruptPerfStandby(newLeaderCh, stopCh)

		// Grab the statelock or stop
		l := newLockGrabber(c.stateLock.Lock, c.stateLock.Unlock, stopCh)
		go l.grab()
		if stopped := l.lockOrStop(); stopped {
			lock.Unlock()
			close(continueCh)
			metrics.MeasureSince([]string{"core", "leadership_setup_failed"}, activeTime)
			return
		}

		if c.Sealed() {
			c.logger.Warn("grabbed HA lock but already sealed, exiting")
			lock.Unlock()
			close(continueCh)
			c.stateLock.Unlock()
			metrics.MeasureSince([]string{"core", "leadership_setup_failed"}, activeTime)
			return
		}

		// Store the lock so that we can manually clear it later if needed
		c.heldHALock = lock

		// Create the active context
		activeCtx, activeCtxCancel := context.WithCancel(namespace.RootContext(nil))
		c.activeContext = activeCtx
		c.activeContextCancelFunc.Store(activeCtxCancel)

		// Perform seal migration
		if err := c.migrateSeal(c.activeContext); err != nil {
			c.logger.Error("seal migration error", "error", err)
			c.barrier.Seal()
			c.logger.Warn("vault is sealed")
			c.heldHALock = nil
			lock.Unlock()
			close(continueCh)
			c.stateLock.Unlock()
			return
		}

		// This block is used to wipe barrier/seal state and verify that
		// everything is sane. If we have no sanity in the barrier, we actually
		// seal, as there's little we can do.
		{
			c.seal.ClearBarrierConfig(activeCtx)
			if c.seal.RecoveryKeySupported() {
				c.seal.ClearRecoveryConfig(activeCtx)
			}

			if err := c.performKeyUpgrades(activeCtx); err != nil {
				c.logger.Error("error performing key upgrades", "error", err)

				// If we fail due to anything other than a context canceled
				// error we should shutdown as we may have the incorrect Keys.
				if !strings.Contains(err.Error(), context.Canceled.Error()) {
					// We call this in a goroutine so that we can give up the
					// statelock and have this shut us down; sealInternal has a
					// workflow where it watches for the stopCh to close so we want
					// to return from here
					go c.Shutdown()
				}

				c.heldHALock = nil
				lock.Unlock()
				close(continueCh)
				c.stateLock.Unlock()
				metrics.MeasureSince([]string{"core", "leadership_setup_failed"}, activeTime)

				// If we are shutting down we should return from this function,
				// otherwise continue
				if !strings.Contains(err.Error(), context.Canceled.Error()) {
					continue
				} else {
					return
				}
			}
		}

		{
			// Clear previous local cluster cert info so we generate new. Since the
			// UUID will have changed, standbys will know to look for new info
			c.localClusterParsedCert.Store((*x509.Certificate)(nil))
			c.localClusterCert.Store(([]byte)(nil))
			c.localClusterPrivateKey.Store((*ecdsa.PrivateKey)(nil))

			if err := c.setupCluster(activeCtx); err != nil {
				c.heldHALock = nil
				lock.Unlock()
				close(continueCh)
				c.stateLock.Unlock()
				c.logger.Error("cluster setup failed", "error", err)
				metrics.MeasureSince([]string{"core", "leadership_setup_failed"}, activeTime)
				continue
			}

		}
		// Advertise as leader
		if err := c.advertiseLeader(activeCtx, uuid, leaderLostCh); err != nil {
			c.heldHALock = nil
			lock.Unlock()
			close(continueCh)
			c.stateLock.Unlock()
			c.logger.Error("leader advertisement setup failed", "error", err)
			metrics.MeasureSince([]string{"core", "leadership_setup_failed"}, activeTime)
			continue
		}

		// Attempt the post-unseal process
		err = c.postUnseal(activeCtx, activeCtxCancel, standardUnsealStrategy{})
		if err == nil {
			c.standby = false
			c.leaderUUID = uuid
			c.metricSink.SetGaugeWithLabels([]string{"core", "active"}, 1, nil)
		}

		close(continueCh)
		c.stateLock.Unlock()

		// Handle a failure to unseal
		if err != nil {
			c.logger.Error("post-unseal setup failed", "error", err)
			lock.Unlock()
			metrics.MeasureSince([]string{"core", "leadership_setup_failed"}, activeTime)
			continue
		}

		// Monitor a loss of leadership
		select {
		case <-leaderLostCh:
			c.logger.Warn("leadership lost, stopping active operation")
		case <-stopCh:
		case <-manualStepDownCh:
			manualStepDown = true
			c.logger.Warn("stepping down from active operation to standby")
		}

		// Stop Active Duty
		{
			// Spawn this in a go routine so we can cancel the context and
			// unblock any inflight requests that are holding the statelock.
			go func() {
				timer := time.NewTimer(DefaultMaxRequestDuration)
				select {
				case <-activeCtx.Done():
					timer.Stop()
					// Attempt to drain any inflight requests
				case <-timer.C:
					activeCtxCancel()
				}
			}()

			// Grab lock if we are not stopped
			l := newLockGrabber(c.stateLock.Lock, c.stateLock.Unlock, stopCh)
			go l.grab()
			stopped := l.lockOrStop()

			// Cancel the context incase the above go routine hasn't done it
			// yet
			activeCtxCancel()
			metrics.MeasureSince([]string{"core", "leadership_lost"}, activeTime)

			// Mark as standby
			c.standby = true
			c.leaderUUID = ""
			c.metricSink.SetGaugeWithLabels([]string{"core", "active"}, 0, nil)

			// Seal
			if err := c.preSeal(); err != nil {
				c.logger.Error("pre-seal teardown failed", "error", err)
			}

			// If we are not meant to keep the HA lock, clear it
			if atomic.LoadUint32(c.keepHALockOnStepDown) == 0 {
				if err := c.clearLeader(uuid); err != nil {
					c.logger.Error("clearing leader advertisement failed", "error", err)
				}

				if err := c.heldHALock.Unlock(); err != nil {
					c.logger.Error("unlocking HA lock failed", "error", err)
				}
				c.heldHALock = nil
			}

			// Advertise ourselves as a standby.
			if c.serviceRegistration != nil {
				if err := c.serviceRegistration.NotifyActiveStateChange(false); err != nil {
					c.logger.Warn("failed to notify standby status", "error", err)
				}
			}

			// If we are stopped return, otherwise unlock the statelock
			if stopped {
				return
			}
			c.stateLock.Unlock()
		}
	}
}

// grabLockOrStop returns stopped=false if the lock is acquired. Returns
// stopped=true if the lock is not acquired, because stopCh was closed. If the
// lock was acquired (stopped=false) then it's up to the caller to unlock. If
// the lock was not acquired (stopped=true), the caller does not hold the lock and
// should not call unlock.
// It's probably better to inline the body of grabLockOrStop into your function
// instead of calling it. If multiple functions call grabLockOrStop, when a deadlock
// occurs, we have no way of knowing who launched the grab goroutine, complicating
// investigation.
func grabLockOrStop(lockFunc, unlockFunc func(), stopCh chan struct{}) (stopped bool) {
	l := newLockGrabber(lockFunc, unlockFunc, stopCh)
	go l.grab()
	return l.lockOrStop()
}

type lockGrabber struct {
	// stopCh provides a way to interrupt the grab-or-stop
	stopCh chan struct{}
	// doneCh is closed when the child goroutine is done.
	doneCh     chan struct{}
	lockFunc   func()
	unlockFunc func()
	// lock protects these variables which are shared by parent and child.
	lock          sync.Mutex
	parentWaiting bool
	locked        bool
}

func newLockGrabber(lockFunc, unlockFunc func(), stopCh chan struct{}) *lockGrabber {
	return &lockGrabber{
		doneCh:        make(chan struct{}),
		lockFunc:      lockFunc,
		unlockFunc:    unlockFunc,
		parentWaiting: true,
		stopCh:        stopCh,
	}
}

// lockOrStop waits for grab to get a lock or give up, see grabLockOrStop for how to use it.
func (l *lockGrabber) lockOrStop() (stopped bool) {
	stop := false
	select {
	case <-l.stopCh:
		stop = true
	case <-l.doneCh:
	}

	// The child goroutine may not have acquired the lock yet.
	l.lock.Lock()
	defer l.lock.Unlock()
	l.parentWaiting = false
	if stop {
		if l.locked {
			l.unlockFunc()
		}
		return true
	}
	return false
}

// grab tries to get a lock, see grabLockOrStop for how to use it.
func (l *lockGrabber) grab() {
	defer close(l.doneCh)
	l.lockFunc()

	// The parent goroutine may or may not be waiting.
	l.lock.Lock()
	defer l.lock.Unlock()
	if !l.parentWaiting {
		l.unlockFunc()
	} else {
		l.locked = true
	}
}

// This checks the leader periodically to ensure that we switch RPC to a new
// leader pretty quickly. There is logic in Leader() already to not make this
// onerous and avoid more traffic than needed, so we just call that and ignore
// the result.
func (c *Core) periodicLeaderRefresh(newLeaderCh chan func(), stopCh chan struct{}) {
	opCount := new(int32)

	clusterAddr := ""
	for {
		timer := time.NewTimer(leaderCheckInterval)
		select {
		case <-timer.C:
			count := atomic.AddInt32(opCount, 1)
			if count > 1 {
				atomic.AddInt32(opCount, -1)
				continue
			}
			// We do this in a goroutine because otherwise if this refresh is
			// called while we're shutting down the call to Leader() can
			// deadlock, which then means stopCh can never been seen and we can
			// block shutdown
			go func() {
				// Bind locally, as the race detector is tripping here
				lopCount := opCount
				isLeader, _, newClusterAddr, _ := c.Leader()

				// If we are the leader reset the clusterAddr since the next
				// failover might go to the node that was previously active.
				if isLeader {
					clusterAddr = ""
				}

				if !isLeader && newClusterAddr != clusterAddr && newLeaderCh != nil {
					select {
					case newLeaderCh <- nil:
						c.logger.Debug("new leader found, triggering new leader channel")
						clusterAddr = newClusterAddr
					default:
						c.logger.Debug("new leader found, but still processing previous leader change")
					}
				}
				atomic.AddInt32(lopCount, -1)
			}()
		case <-stopCh:
			timer.Stop()
			return
		}
	}
}

// periodicCheckKeyUpgrade is used to watch for key rotation events as a standby
func (c *Core) periodicCheckKeyUpgrades(ctx context.Context, stopCh chan struct{}) {
	raftBackend := c.getRaftBackend()
	isRaft := raftBackend != nil

	opCount := new(int32)
	for {
		timer := time.NewTimer(keyRotateCheckInterval)
		select {
		case <-timer.C:
			count := atomic.AddInt32(opCount, 1)
			if count > 1 {
				atomic.AddInt32(opCount, -1)
				continue
			}

			go func() {
				// Bind locally, as the race detector is tripping here
				lopCount := opCount

				// Only check if we are a standby
				c.stateLock.RLock()
				standby := c.standby
				c.stateLock.RUnlock()
				if !standby {
					atomic.AddInt32(lopCount, -1)
					return
				}

				// Check for a poison pill. If we can read it, it means we have stale
				// keys (e.g. from replication being activated) and we need to seal to
				// be unsealed again.
				entry, _ := c.barrier.Get(ctx, poisonPillPath)
				entryDR, _ := c.barrier.Get(ctx, poisonPillDRPath)
				if (entry != nil && len(entry.Value) > 0) || (entryDR != nil && len(entryDR.Value) > 0) {
					c.logger.Warn("encryption keys have changed out from underneath us (possibly due to replication enabling), must be unsealed again")
					// If we are using raft storage we do not want to shut down
					// raft during replication secondary enablement. This will
					// allow us to keep making progress on the raft log.
					go c.sealInternalWithOptions(true, false, !isRaft)
					atomic.AddInt32(lopCount, -1)
					return
				}

				if err := c.checkKeyUpgrades(ctx); err != nil {
					c.logger.Error("key rotation periodic upgrade check failed", "error", err)
				}

				if isRaft {
					hasState, err := raftBackend.HasState()
					if err != nil {
						c.logger.Error("could not check raft state", "error", err)
					}

					if raftBackend.Initialized() && hasState {
						if err := c.checkRaftTLSKeyUpgrades(ctx); err != nil {
							c.logger.Error("raft tls periodic upgrade check failed", "error", err)
						}
					}
				}

				atomic.AddInt32(lopCount, -1)
				return
			}()
		case <-stopCh:
			timer.Stop()
			return
		}
	}
}

// checkKeyUpgrades is used to check if there have been any key rotations
// and if there is a chain of upgrades available
func (c *Core) checkKeyUpgrades(ctx context.Context) error {
	for {
		// Check for an upgrade
		didUpgrade, newTerm, err := c.barrier.CheckUpgrade(ctx)
		if err != nil {
			return err
		}

		// Nothing to do if no upgrade
		if !didUpgrade {
			break
		}
		if c.logger.IsInfo() {
			c.logger.Info("upgraded to new key term", "term", newTerm)
		}
	}
	return nil
}

func (c *Core) reloadRootKey(ctx context.Context) error {
	if err := c.barrier.ReloadRootKey(ctx); err != nil {
		return fmt.Errorf("error reloading root key: %w", err)
	}
	return nil
}

func (c *Core) reloadShamirKey(ctx context.Context) error {
	_ = c.seal.ClearBarrierConfig(ctx)

	cfg, _ := c.seal.BarrierConfig(ctx)
	if cfg == nil {
		return nil
	}

	var shamirKey []byte
	switch c.seal.StoredKeysSupported() {
	case seal.StoredKeysSupportedGeneric:
		return nil
	case seal.StoredKeysSupportedShamirRoot:
		entry, err := c.barrier.Get(ctx, shamirKekPath)
		if err != nil {
			return err
		}
		if entry == nil {
			return nil
		}
		shamirKey = entry.Value
	case seal.StoredKeysNotSupported:
		keyring, err := c.barrier.Keyring()
		if err != nil {
			return fmt.Errorf("failed to update seal access: %w", err)
		}
		shamirKey = keyring.rootKey
	}
	return c.seal.GetAccess().SetShamirSealKey(shamirKey)
}

func (c *Core) performKeyUpgrades(ctx context.Context) error {
	if err := c.checkKeyUpgrades(ctx); err != nil {
		return fmt.Errorf("error checking for key upgrades: %w", err)
	}

	if err := c.reloadRootKey(ctx); err != nil {
		return fmt.Errorf("error reloading root key: %w", err)
	}

	if err := c.barrier.ReloadKeyring(ctx); err != nil {
		return fmt.Errorf("error reloading keyring: %w", err)
	}

	if err := c.reloadShamirKey(ctx); err != nil {
		return fmt.Errorf("error reloading shamir kek key: %w", err)
	}

	if err := c.scheduleUpgradeCleanup(ctx); err != nil {
		return fmt.Errorf("error scheduling upgrade cleanup: %w", err)
	}

	return nil
}

// scheduleUpgradeCleanup is used to ensure that all the upgrade paths
// are cleaned up in a timely manner if a leader failover takes place
func (c *Core) scheduleUpgradeCleanup(ctx context.Context) error {
	// List the upgrades
	upgrades, err := c.barrier.List(ctx, keyringUpgradePrefix)
	if err != nil {
		return fmt.Errorf("failed to list upgrades: %w", err)
	}

	// Nothing to do if no upgrades
	if len(upgrades) == 0 {
		return nil
	}

	// Schedule cleanup for all of them
	time.AfterFunc(c.KeyRotateGracePeriod(), func() {
		sealed, err := c.barrier.Sealed()
		if err != nil {
			c.logger.Warn("failed to check barrier status at upgrade cleanup time")
			return
		}
		if sealed {
			c.logger.Warn("barrier sealed at upgrade cleanup time")
			return
		}
		for _, upgrade := range upgrades {
			path := fmt.Sprintf("%s%s", keyringUpgradePrefix, upgrade)
			if err := c.barrier.Delete(ctx, path); err != nil {
				c.logger.Error("failed to cleanup upgrade", "path", path, "error", err)
			}
		}
	})
	return nil
}

// acquireLock blocks until the lock is acquired, returning the leaderLostCh
func (c *Core) acquireLock(lock physical.Lock, stopCh <-chan struct{}) <-chan struct{} {
	for {
		// Attempt lock acquisition
		leaderLostCh, err := lock.Lock(stopCh)
		if err == nil {
			return leaderLostCh
		}

		// Retry the acquisition
		c.logger.Error("failed to acquire lock", "error", err)
		timer := time.NewTimer(lockRetryInterval)
		select {
		case <-timer.C:
		case <-stopCh:
			timer.Stop()
			return nil
		}
	}
}

// advertiseLeader is used to advertise the current node as leader
func (c *Core) advertiseLeader(ctx context.Context, uuid string, leaderLostCh <-chan struct{}) error {
	if leaderLostCh != nil {
		go c.cleanLeaderPrefix(ctx, uuid, leaderLostCh)
	}

	var key *ecdsa.PrivateKey
	switch c.localClusterPrivateKey.Load().(type) {
	case *ecdsa.PrivateKey:
		key = c.localClusterPrivateKey.Load().(*ecdsa.PrivateKey)
	default:
		c.logger.Error("unknown cluster private key type", "key_type", fmt.Sprintf("%T", c.localClusterPrivateKey.Load()))
		return fmt.Errorf("unknown cluster private key type %T", c.localClusterPrivateKey.Load())
	}

	keyParams := &certutil.ClusterKeyParams{
		Type: corePrivateKeyTypeP521,
		X:    key.X,
		Y:    key.Y,
		D:    key.D,
	}

	locCert := c.localClusterCert.Load().([]byte)
	localCert := make([]byte, len(locCert))
	copy(localCert, locCert)
	adv := &activeAdvertisement{
		RedirectAddr:     c.redirectAddr,
		ClusterAddr:      c.ClusterAddr(),
		ClusterCert:      localCert,
		ClusterKeyParams: keyParams,
	}
	val, err := jsonutil.EncodeJSON(adv)
	if err != nil {
		return err
	}
	ent := &logical.StorageEntry{
		Key:   coreLeaderPrefix + uuid,
		Value: val,
	}
	err = c.barrier.Put(ctx, ent)
	if err != nil {
		return err
	}

	if c.serviceRegistration != nil {
		if err := c.serviceRegistration.NotifyActiveStateChange(true); err != nil {
			if c.logger.IsWarn() {
				c.logger.Warn("failed to notify active status", "error", err)
			}
		}
	}
	return nil
}

func (c *Core) cleanLeaderPrefix(ctx context.Context, uuid string, leaderLostCh <-chan struct{}) {
	keys, err := c.barrier.List(ctx, coreLeaderPrefix)
	if err != nil {
		c.logger.Error("failed to list entries in core/leader", "error", err)
		return
	}
	for len(keys) > 0 {
		timer := time.NewTimer(leaderPrefixCleanDelay)
		select {
		case <-timer.C:
			if keys[0] != uuid {
				c.barrier.Delete(ctx, coreLeaderPrefix+keys[0])
			}
			keys = keys[1:]
		case <-leaderLostCh:
			timer.Stop()
			return
		}
	}
}

// clearLeader is used to clear our leadership entry
func (c *Core) clearLeader(uuid string) error {
	key := coreLeaderPrefix + uuid
	return c.barrier.Delete(context.Background(), key)
}

func (c *Core) SetNeverBecomeActive(on bool) {
	if on {
		atomic.StoreUint32(c.neverBecomeActive, 1)
	} else {
		atomic.StoreUint32(c.neverBecomeActive, 0)
	}
}
