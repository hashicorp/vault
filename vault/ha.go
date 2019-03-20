package vault

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	multierror "github.com/hashicorp/go-multierror"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
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

	// keyRotateGracePeriod is how long we allow an upgrade path
	// for standby instances before we delete the upgrade keys
	keyRotateGracePeriod = 2 * time.Minute

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
func (c *Core) PerfStandby() bool {
	c.stateLock.RLock()
	perfStandby := c.perfStandby
	c.stateLock.RUnlock()
	return perfStandby
}

// Leader is used to get the current active leader
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
		return true, c.redirectAddr, c.clusterAddr, nil
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
		retErr = multierror.Append(retErr, errors.New("nil request to step-down"))
		return retErr
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

	acl, te, entity, identityPolicies, err := c.fetchACLTokenEntryAndEntity(ctx, req)
	if err != nil {
		retErr = multierror.Append(retErr, err)
		return retErr
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

	logInput := &audit.LogInput{
		Auth:    auth,
		Request: req,
	}
	if err := c.auditBroker.LogRequest(ctx, logInput, c.auditedHeaders); err != nil {
		c.logger.Error("failed to audit request", "request_path", req.Path, "error", err)
		retErr = multierror.Append(retErr, errors.New("failed to audit request, cannot continue"))
		return retErr
	}

	if entity != nil && entity.Disabled {
		c.logger.Warn("permission denied as the entity on the token is disabled")
		retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
		c.stateLock.RUnlock()
		return retErr
	}

	if te != nil && te.EntityID != "" && entity == nil {
		c.logger.Warn("permission denied as the entity on the token is invalid")
		retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
		c.stateLock.RUnlock()
		return retErr
	}

	// Attempt to use the token (decrement num_uses)
	if te != nil {
		te, err = c.tokenStore.UseToken(ctx, te)
		if err != nil {
			c.logger.Error("failed to use token", "error", err)
			retErr = multierror.Append(retErr, ErrInternalError)
			return retErr
		}
		if te == nil {
			// Token has been revoked
			retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
			return retErr
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
		// Monitor for key rotation
		keyRotateStop := make(chan struct{})

		g.Add(func() error {
			c.periodicCheckKeyUpgrade(context.Background(), keyRotateStop)
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
			}
		}

		// Create a lock
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			c.logger.Error("failed to generate uuid", "error", err)
			return
		}
		lock, err := c.ha.LockWith(CoreLockPath, uuid)
		if err != nil {
			c.logger.Error("failed to create lock", "error", err)
			return
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
		if stopped := grabLockOrStop(c.stateLock.Lock, c.stateLock.Unlock, stopCh); stopped {
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

		// This block is used to wipe barrier/seal state and verify that
		// everything is sane. If we have no sanity in the barrier, we actually
		// seal, as there's little we can do.
		{
			c.seal.SetBarrierConfig(activeCtx, nil)
			if c.seal.RecoveryKeySupported() {
				c.seal.SetRecoveryConfig(activeCtx, nil)
			}

			if err := c.performKeyUpgrades(activeCtx); err != nil {
				// We call this in a goroutine so that we can give up the
				// statelock and have this shut us down; sealInternal has a
				// workflow where it watches for the stopCh to close so we want
				// to return from here
				c.logger.Error("error performing key upgrades", "error", err)
				go c.Shutdown()
				c.heldHALock = nil
				lock.Unlock()
				close(continueCh)
				c.stateLock.Unlock()
				metrics.MeasureSince([]string{"core", "leadership_setup_failed"}, activeTime)
				return
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
				select {
				case <-activeCtx.Done():
				// Attempt to drain any inflight requests
				case <-time.After(DefaultMaxRequestDuration):
					activeCtxCancel()
				}
			}()

			// Grab lock if we are not stopped
			stopped := grabLockOrStop(c.stateLock.Lock, c.stateLock.Unlock, stopCh)

			// Cancel the context incase the above go routine hasn't done it
			// yet
			activeCtxCancel()
			metrics.MeasureSince([]string{"core", "leadership_lost"}, activeTime)

			// Mark as standby
			c.standby = true

			// Seal
			if err := c.preSeal(); err != nil {
				c.logger.Error("pre-seal teardown failed", "error", err)
			}

			// If we are not meant to keep the HA lock, clear it
			if atomic.LoadUint32(c.keepHALockOnStepDown) == 0 {
				if err := c.clearLeader(uuid); err != nil {
					c.logger.Error("clearing leader advertisement failed", "error", err)
				}

				c.heldHALock.Unlock()
				c.heldHALock = nil
			}

			// If we are stopped return, otherwise unlock the statelock
			if stopped {
				return
			}
			c.stateLock.Unlock()
		}
	}
}

// grabLockOrStop returns true if we failed to get the lock before stopCh
// was closed.  Returns false if the lock was obtained, in which case it's
// the caller's responsibility to unlock it.
func grabLockOrStop(lockFunc, unlockFunc func(), stopCh chan struct{}) (stopped bool) {
	// Grab the lock as we need it for cluster setup, which needs to happen
	// before advertising;
	lockGrabbedCh := make(chan struct{})
	go func() {
		// Grab the lock
		lockFunc()
		// If stopCh has been closed, which only happens while the
		// stateLock is held, we have actually terminated, so we just
		// instantly give up the lock, otherwise we notify that it's ready
		// for consumption
		select {
		case <-stopCh:
			unlockFunc()
		default:
			close(lockGrabbedCh)
		}
	}()

	select {
	case <-stopCh:
		return true
	case <-lockGrabbedCh:
		// We now have the lock and can use it
	}

	return false
}

// This checks the leader periodically to ensure that we switch RPC to a new
// leader pretty quickly. There is logic in Leader() already to not make this
// onerous and avoid more traffic than needed, so we just call that and ignore
// the result.
func (c *Core) periodicLeaderRefresh(newLeaderCh chan func(), stopCh chan struct{}) {
	opCount := new(int32)

	clusterAddr := ""
	for {
		select {
		case <-time.After(leaderCheckInterval):
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
			return
		}
	}
}

// periodicCheckKeyUpgrade is used to watch for key rotation events as a standby
func (c *Core) periodicCheckKeyUpgrade(ctx context.Context, stopCh chan struct{}) {
	opCount := new(int32)
	for {
		select {
		case <-time.After(keyRotateCheckInterval):
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
				if entry != nil && len(entry.Value) > 0 {
					c.logger.Warn("encryption keys have changed out from underneath us (possibly due to replication enabling), must be unsealed again")
					go c.Shutdown()
					atomic.AddInt32(lopCount, -1)
					return
				}

				if err := c.checkKeyUpgrades(ctx); err != nil {
					c.logger.Error("key rotation periodic upgrade check failed", "error", err)
				}

				atomic.AddInt32(lopCount, -1)
				return
			}()
		case <-stopCh:
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

func (c *Core) performKeyUpgrades(ctx context.Context) error {
	if err := c.checkKeyUpgrades(ctx); err != nil {
		return errwrap.Wrapf("error checking for key upgrades: {{err}}", err)
	}

	if err := c.barrier.ReloadMasterKey(ctx); err != nil {
		return errwrap.Wrapf("error reloading master key: {{err}}", err)
	}

	if err := c.barrier.ReloadKeyring(ctx); err != nil {
		return errwrap.Wrapf("error reloading keyring: {{err}}", err)
	}

	if err := c.scheduleUpgradeCleanup(ctx); err != nil {
		return errwrap.Wrapf("error scheduling upgrade cleanup: {{err}}", err)
	}

	return nil
}

// scheduleUpgradeCleanup is used to ensure that all the upgrade paths
// are cleaned up in a timely manner if a leader failover takes place
func (c *Core) scheduleUpgradeCleanup(ctx context.Context) error {
	// List the upgrades
	upgrades, err := c.barrier.List(ctx, keyringUpgradePrefix)
	if err != nil {
		return errwrap.Wrapf("failed to list upgrades: {{err}}", err)
	}

	// Nothing to do if no upgrades
	if len(upgrades) == 0 {
		return nil
	}

	// Schedule cleanup for all of them
	time.AfterFunc(keyRotateGracePeriod, func() {
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
		select {
		case <-time.After(lockRetryInterval):
		case <-stopCh:
			return nil
		}
	}
}

// advertiseLeader is used to advertise the current node as leader
func (c *Core) advertiseLeader(ctx context.Context, uuid string, leaderLostCh <-chan struct{}) error {
	go c.cleanLeaderPrefix(ctx, uuid, leaderLostCh)

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
		ClusterAddr:      c.clusterAddr,
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

	sd, ok := c.ha.(physical.ServiceDiscovery)
	if ok {
		if err := sd.NotifyActiveStateChange(); err != nil {
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
		select {
		case <-time.After(leaderPrefixCleanDelay):
			if keys[0] != uuid {
				c.barrier.Delete(ctx, coreLeaderPrefix+keys[0])
			}
			keys = keys[1:]
		case <-leaderLostCh:
			return
		}
	}
}

// clearLeader is used to clear our leadership entry
func (c *Core) clearLeader(uuid string) error {
	key := coreLeaderPrefix + uuid
	err := c.barrier.Delete(context.Background(), key)

	// Advertise ourselves as a standby
	sd, ok := c.ha.(physical.ServiceDiscovery)
	if ok {
		if err := sd.NotifyActiveStateChange(); err != nil {
			if c.logger.IsWarn() {
				c.logger.Warn("failed to notify standby status", "error", err)
			}
		}
	}

	return err
}

func (c *Core) SetNeverBecomeActive(on bool) {
	if on {
		atomic.StoreUint32(c.neverBecomeActive, 1)
	} else {
		atomic.StoreUint32(c.neverBecomeActive, 0)
	}
}
