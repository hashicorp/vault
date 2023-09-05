// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-discover"
	discoverk8s "github.com/hashicorp/go-discover/provider/k8s"
	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/go-secure-stdlib/tlsutil"
	"github.com/hashicorp/go-uuid"
	goversion "github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/http2"
)

const (
	// undoLogMonitorInterval is how often the leader checks to see
	// if all the cluster members it knows about are new enough to support
	// undo logs.
	undoLogMonitorInterval = time.Second

	// undoLogSafeVersion is the minimum version Vault must be at in order
	// for undo logs to be turned on.
	undoLogSafeVersion = "1.12.0"
)

var (
	raftTLSStoragePath    = "core/raft/tls"
	raftTLSRotationPeriod = 24 * time.Hour

	raftAutopilotConfigurationStoragePath = "core/raft/autopilot/configuration"

	// TestingUpdateClusterAddr is used in tests to override the cluster address
	TestingUpdateClusterAddr uint32

	ErrJoinWithoutAutoloading = errors.New("attempt to join a cluster using autoloaded licenses while not using autoloading ourself")
)

// GetRaftNodeID returns the raft node ID if there is one, or an empty string if there's not
func (c *Core) GetRaftNodeID() string {
	rb := c.getRaftBackend()

	if rb == nil {
		return ""
	} else {
		return rb.NodeID()
	}
}

func (c *Core) GetRaftIndexes() (committed uint64, applied uint64) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()

	raftStorage, ok := c.underlyingPhysical.(*raft.RaftBackend)
	if !ok {
		return 0, 0
	}

	return raftStorage.CommittedIndex(), raftStorage.AppliedIndex()
}

// startRaftBackend will call SetupCluster in the raft backend which starts raft
// up and enables the cluster handler.
func (c *Core) startRaftBackend(ctx context.Context) (retErr error) {
	raftBackend := c.getRaftBackend()
	if raftBackend == nil {
		return nil
	}

	var creating bool
	var raftTLS *raft.TLSKeyring
	if !raftBackend.Initialized() {
		// Retrieve the raft TLS information
		raftTLSEntry, err := c.barrier.Get(ctx, raftTLSStoragePath)
		if err != nil {
			return err
		}

		switch raftTLSEntry {
		case nil:
			// If this is HA-only and no TLS keyring is found, that means the
			// cluster has not been bootstrapped or joined. We return early here in
			// this case. If we return here, the raft object has not been instantiated,
			// and a bootstrap call should be made.
			if c.isRaftHAOnly() {
				c.logger.Trace("skipping raft backend setup during unseal, no bootstrap operation has been started yet")
				return nil
			}

			// If we did not find a TLS keyring we will attempt to create one here.
			// This happens after a storage migration process. This node is also
			// marked to start as leader so we can write the new TLS Key. This is an
			// error condition if there are already multiple nodes in the cluster,
			// and the below storage write will fail. If the cluster is somehow in
			// this state the unseal will fail and a cluster recovery will need to
			// be done.
			creating = true
			raftTLSKey, err := raft.GenerateTLSKey(c.secureRandomReader)
			if err != nil {
				return err
			}

			raftTLS = &raft.TLSKeyring{
				Keys:        []*raft.TLSKey{raftTLSKey},
				ActiveKeyID: raftTLSKey.ID,
			}
		default:
			raftTLS = new(raft.TLSKeyring)
			if err := raftTLSEntry.DecodeJSON(raftTLS); err != nil {
				return err
			}
		}

		hasState, err := raftBackend.HasState()
		if err != nil {
			return err
		}

		// This can be hit on follower nodes that got their config updated to use
		// raft for HA-only before they are joined to the cluster. Since followers
		// in this case use shared storage, it doesn't return early from the TLS
		// case above, but there's not raft state yet for the backend to call
		// raft.SetupCluster.
		if !hasState {
			c.logger.Trace("skipping raft backend setup during unseal, no raft state found")
			return nil
		}

		raftBackend.SetRestoreCallback(c.raftSnapshotRestoreCallback(true, true))

		if err := raftBackend.SetupCluster(ctx, raft.SetupOpts{
			TLSKeyring:      raftTLS,
			ClusterListener: c.getClusterListener(),
			StartAsLeader:   creating,
		}); err != nil {
			return err
		}
	}

	defer func() {
		if retErr != nil {
			c.logger.Info("stopping raft server")
			if err := raftBackend.TeardownCluster(c.getClusterListener()); err != nil {
				c.logger.Error("failed to stop raft server", "error", err)
			}
		}
	}()

	// If we are in need of creating the TLS keyring then we should write it out
	// to storage here. If we fail it may mean we couldn't become leader and we
	// should error out.
	if creating {
		c.logger.Info("writing raft TLS keyring to storage")
		entry, err := logical.StorageEntryJSON(raftTLSStoragePath, raftTLS)
		if err != nil {
			c.logger.Error("error marshaling raft TLS keyring", "error", err)
			return err
		}
		if err := c.barrier.Put(ctx, entry); err != nil {
			c.logger.Error("error writing raft TLS keyring", "error", err)
			return err
		}
	}

	return nil
}

func (c *Core) monitorUndoLogs() error {
	logger := c.logger.Named("undo-log-watcher")
	logger.Debug("starting undo log watcher")
	ctx := c.activeContext
	raftBackend := c.getRaftBackend()

	// First check storage and bail early if we already know undo logs are safe
	persisted, err := c.UndoLogsPersisted()
	if err != nil {
		return fmt.Errorf("error checking for undo logs persistence: %w", err)
	}

	if persisted {
		logger.Debug("undo logs are safe, no need to check any more")
		return nil
	}

	// If undo logs have been explicitly enabled, likely via VAULT_REPLICATION_USE_UNDO_LOGS, then exit, as presumably
	// we don't want to be checking for safety if we already know it's safe.
	if c.UndoLogsEnabled() {
		logger.Debug("undo logs have been explicitly enabled. exiting monitor.")
		return nil
	}

	minimumVersion, err := goversion.NewSemver(undoLogSafeVersion)
	if err != nil {
		return fmt.Errorf("minimum undo log version (%q) won't parse: %w", undoLogSafeVersion, err)
	}

	go func() {
		ticker := time.NewTicker(undoLogMonitorInterval)
		defer ticker.Stop()

		logger.Debug("undo logs have not been enabled yet, possibly due to a recent upgrade. starting a periodic check")

		for {
			select {
			case <-ticker.C:
			case <-ctx.Done():
				return
			}

			// Check the raft configuration for expected servers
			config, err := raftBackend.GetConfiguration(ctx)
			if err != nil {
				logger.Error("couldn't read raft config", "error", err)
				continue
			}

			// This tracks which servers we expect to find in the cluster, from the raft config.
			expectedServers := make(map[string]struct{})
			for _, server := range config.Servers {
				expectedServers[server.Address] = struct{}{}
			}

			// Retrieve all the nodes in our cluster
			nodes, err := c.getHAMembers()
			if err != nil {
				logger.Error("error getting HA members", "error", err)
				continue
			}

			// Check the versions of all of the cluster members. If they're all >= 1.12, undo logs are safe to enable.
			// If any are < 1.12, undo logs should remain off, regardless of how it was initially configured.
			enable := true
			for _, node := range nodes {
				nodeVersion, err := goversion.NewSemver(node.Version)
				if err != nil {
					logger.Error("error parsing node version", "node version", node.Version, "node", node.ClusterAddress, "error", err)
					break
				}

				if nodeVersion.LessThan(minimumVersion) {
					logger.Debug("node version is less than the minimum, disabling undo logs", "node", node.ClusterAddress, "version", node.Version)
					enable = false
					break
				} else {
					// Raft nodes have their address listed without a scheme, e.g. 127.0.0.1:8201. HA nodes that we get from
					// getHAMembers() have their address listed with a scheme, e.g. https://127.0.0.1:8201. So we need to
					// parse the HA node address and reconstruct it ourselves to get a matching hash key.
					clusterAddr, err := url.Parse(node.ClusterAddress)
					if err != nil {
						logger.Error("error parsing node cluster address", "node", node.ClusterAddress, "error", err)
						break
					}

					// Deleting from expectedServers means the node in question is running a Vault version greater than
					// or equal to the minimum.
					delete(expectedServers, clusterAddr.Host)
				}
			}

			// If expectedServers still has nodes in it, that means either we broke from the above loop because some
			// node's version was too low or because some member of the cluster hasn't sent an echo yet. Either way,
			// it means we can't enable undo logs.
			if len(expectedServers) != 0 {
				enable = false
			}

			if enable {
				logger.Debug("undo logs can be safely enabled now")
				err := c.PersistUndoLogs()
				if err != nil {
					logger.Error("error persisting undo logs safety", "error", err)
					continue
				}

				logger.Debug("undo logs have been enabled and this has been persisted to storage. shutting down the checker.")
				return
			}
		}
	}()

	return nil
}

func (c *Core) setupRaftActiveNode(ctx context.Context) error {
	raftBackend := c.getRaftBackend()
	if raftBackend == nil {
		return nil
	}

	c.logger.Info("starting raft active node")
	raftBackend.SetEffectiveSDKVersion(c.effectiveSDKVersion)

	autopilotConfig, err := c.loadAutopilotConfiguration(ctx)
	if err != nil {
		c.logger.Error("failed to load autopilot config from storage when setting up cluster; continuing since autopilot falls back to default config", "error", err)
	}
	disableAutopilot := c.disableAutopilot

	raftBackend.SetupAutopilot(c.activeContext, autopilotConfig, c.raftFollowerStates, disableAutopilot)

	c.pendingRaftPeers = &sync.Map{}

	// Reload the raft TLS keys to ensure we are using the latest version.
	if err := c.checkRaftTLSKeyUpgrades(ctx); err != nil {
		return err
	}

	// We always want to start this watcher - if undo logs are safe to be enabled, it will exit quickly. If not, it
	// will monitor for safety until they are enabled.
	if err := c.monitorUndoLogs(); err != nil {
		return err
	}
	return c.startPeriodicRaftTLSRotate(ctx)
}

func (c *Core) stopRaftActiveNode() {
	raftBackend := c.getRaftBackend()
	if raftBackend == nil {
		return
	}

	c.logger.Info("stopping raft active node")

	if !raftBackend.AutopilotDisabled() {
		raftBackend.StopAutopilot()
	}

	c.pendingRaftPeers = nil
	c.stopPeriodicRaftTLSRotate()
}

func (c *Core) startPeriodicRaftTLSRotate(ctx context.Context) error {
	raftBackend := c.getRaftBackend()

	// No-op if raft is not being used
	if raftBackend == nil {
		return nil
	}

	c.raftTLSRotationStopCh = make(chan struct{})
	logger := c.logger.Named("raft")

	if c.isRaftHAOnly() {
		return c.raftTLSRotateDirect(ctx, logger, c.raftTLSRotationStopCh)
	}

	return c.raftTLSRotatePhased(ctx, logger, raftBackend, c.raftTLSRotationStopCh)
}

// raftTLSRotateDirect will spawn a go routine in charge of periodically
// rotating the TLS certs and keys used for raft traffic.
//
// The logic for updating the TLS keyring is through direct storage update. This
// is called whenever raft is used for HA-only, which means that the underlying
// storage is a shared physical object, thus requiring no additional
// coordination.
func (c *Core) raftTLSRotateDirect(ctx context.Context, logger hclog.Logger, stopCh chan struct{}) error {
	logger.Info("creating new raft TLS config")

	rotateKeyring := func() (time.Time, error) {
		// Create a new key
		raftTLSKey, err := raft.GenerateTLSKey(c.secureRandomReader)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to generate new raft TLS key: %w", err)
		}

		// Read the existing keyring
		keyring, err := c.raftReadTLSKeyring(ctx)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to read raft TLS keyring: %w", err)
		}

		// Advance the term and store the new key, replacing the old one.
		// Unlike phased rotation, we don't need to update AppliedIndex since
		// we don't rely on it to check whether the followers got the key. A
		// shared storage means that followers will have the key as soon as it's
		// written to storage.
		keyring.Term += 1
		keyring.Keys[0] = raftTLSKey
		keyring.ActiveKeyID = raftTLSKey.ID
		entry, err := logical.StorageEntryJSON(raftTLSStoragePath, keyring)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to json encode keyring: %w", err)
		}
		if err := c.barrier.Put(ctx, entry); err != nil {
			return time.Time{}, fmt.Errorf("failed to write keyring: %w", err)
		}

		logger.Info("wrote new raft TLS config")

		// Schedule the next rotation
		return raftTLSKey.CreatedTime.Add(raftTLSRotationPeriod), nil
	}

	// Read the keyring to calculate the time of next rotation.
	keyring, err := c.raftReadTLSKeyring(ctx)
	if err != nil {
		return err
	}

	activeKey := keyring.GetActive()
	if activeKey == nil {
		return errors.New("no active raft TLS key found")
	}

	go func() {
		nextRotationTime := activeKey.CreatedTime.Add(raftTLSRotationPeriod)

		var backoff bool
		for {
			// If we encountered and error we should try to create the key
			// again.
			if backoff {
				nextRotationTime = time.Now().Add(10 * time.Second)
				backoff = false
			}

			timer := time.NewTimer(time.Until(nextRotationTime))
			select {
			case <-timer.C:
				// It's time to rotate the keys
				next, err := rotateKeyring()
				if err != nil {
					logger.Error("failed to rotate TLS key", "error", err)
					backoff = true
					continue
				}

				nextRotationTime = next

			case <-stopCh:
				timer.Stop()
				return
			}
		}
	}()

	return nil
}

// raftTLSRotatePhased will spawn a go routine in charge of periodically
// rotating the TLS certs and keys used for raft traffic.
//
// The logic for updating the TLS certificate uses a pseudo two phase commit
// using the known applied indexes from standby nodes. When writing a new Key
// it will be appended to the end of the keyring. Standbys can start accepting
// connections with this key as soon as they see the update. Then it will write
// the keyring a second time indicating the applied index for this key update.
//
// The active node will wait until it sees all standby nodes are at or past the
// applied index for this update. At that point it will delete the older key
// and make the new key active. The key isn't officially in use until this
// happens. The dual write ensures the standby at least gets the first update
// containing the key before the active node switches over to using it.
//
// If a standby is shut down then it cannot advance the key term until it
// receives the update. This ensures a standby node isn't left behind and unable
// to reconnect with the cluster. Additionally, only one outstanding key
// is allowed for this same reason (max keyring size of 2).
func (c *Core) raftTLSRotatePhased(ctx context.Context, logger hclog.Logger, raftBackend *raft.RaftBackend, stopCh chan struct{}) error {
	followerStates := c.raftFollowerStates
	followerStates.Clear()

	// Pre-populate the follower list with the set of peers.
	raftConfig, err := raftBackend.GetConfiguration(ctx)
	if err != nil {
		return err
	}
	for _, server := range raftConfig.Servers {
		if server.NodeID != raftBackend.NodeID() {
			followerStates.Update(&raft.EchoRequestUpdate{
				NodeID:          server.NodeID,
				AppliedIndex:    0,
				Term:            0,
				DesiredSuffrage: "voter",
			})
		}
	}

	// rotateKeyring writes new key data to the keyring and adds an applied
	// index that is used to verify it has been committed. The keys written in
	// this function can be used on standbys but the active node doesn't start
	// using it yet.
	rotateKeyring := func() (time.Time, error) {
		// Read the existing keyring
		keyring, err := c.raftReadTLSKeyring(ctx)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to read raft TLS keyring: %w", err)
		}

		switch {
		case len(keyring.Keys) == 2 && keyring.Keys[1].AppliedIndex == 0:
			// If this case is hit then the second write to add the applied
			// index failed. Attempt to write it again.
			keyring.Keys[1].AppliedIndex = raftBackend.AppliedIndex()
			keyring.AppliedIndex = raftBackend.AppliedIndex()
			entry, err := logical.StorageEntryJSON(raftTLSStoragePath, keyring)
			if err != nil {
				return time.Time{}, fmt.Errorf("failed to json encode keyring: %w", err)
			}
			if err := c.barrier.Put(ctx, entry); err != nil {
				return time.Time{}, fmt.Errorf("failed to write keyring: %w", err)
			}

		case len(keyring.Keys) > 1:
			// If there already exists a pending key update then the update
			// hasn't replicated down to all standby nodes yet. Don't allow any
			// new keys to be created until all standbys have seen this previous
			// rotation. As a backoff strategy, another rotation attempt is
			// scheduled for 5 minutes from now.
			logger.Warn("skipping new raft TLS config creation, keys are pending")
			return time.Now().Add(time.Minute * 5), nil
		}

		logger.Info("creating new raft TLS config")

		// Create a new key
		raftTLSKey, err := raft.GenerateTLSKey(c.secureRandomReader)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to generate new raft TLS key: %w", err)
		}

		// Advance the term and store the new key
		keyring.Term += 1
		keyring.Keys = append(keyring.Keys, raftTLSKey)
		entry, err := logical.StorageEntryJSON(raftTLSStoragePath, keyring)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to json encode keyring: %w", err)
		}
		if err := c.barrier.Put(ctx, entry); err != nil {
			return time.Time{}, fmt.Errorf("failed to write keyring: %w", err)
		}

		// Write the keyring again with the new applied index. This allows us to
		// track if standby nodes received the update.
		keyring.Keys[1].AppliedIndex = raftBackend.AppliedIndex()
		keyring.AppliedIndex = raftBackend.AppliedIndex()
		entry, err = logical.StorageEntryJSON(raftTLSStoragePath, keyring)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to json encode keyring: %w", err)
		}
		if err := c.barrier.Put(ctx, entry); err != nil {
			return time.Time{}, fmt.Errorf("failed to write keyring: %w", err)
		}

		logger.Info("wrote new raft TLS config")
		// Schedule the next rotation
		return raftTLSKey.CreatedTime.Add(raftTLSRotationPeriod), nil
	}

	// checkCommitted verifies key updates have been applied to all nodes and
	// finalizes the rotation by deleting the old keys and updating the raft
	// backend.
	checkCommitted := func() error {
		keyring, err := c.raftReadTLSKeyring(ctx)
		if err != nil {
			return fmt.Errorf("failed to read raft TLS keyring: %w", err)
		}

		switch {
		case len(keyring.Keys) == 1:
			// No Keys to apply
			return nil
		case keyring.Keys[1].AppliedIndex != keyring.AppliedIndex:
			// We haven't fully committed the new key, continue here
			return nil
		case followerStates.MinIndex() < keyring.AppliedIndex:
			// Not all the followers have applied the latest key
			return nil
		}

		// Upgrade to the new key
		keyring.Keys = keyring.Keys[1:]
		keyring.ActiveKeyID = keyring.Keys[0].ID
		keyring.Term += 1
		entry, err := logical.StorageEntryJSON(raftTLSStoragePath, keyring)
		if err != nil {
			return fmt.Errorf("failed to json encode keyring: %w", err)
		}
		if err := c.barrier.Put(ctx, entry); err != nil {
			return fmt.Errorf("failed to write keyring: %w", err)
		}

		// Update the TLS Key in the backend
		if err := raftBackend.SetTLSKeyring(keyring); err != nil {
			return fmt.Errorf("failed to install keyring: %w", err)
		}

		logger.Info("installed new raft TLS key", "term", keyring.Term)
		return nil
	}

	// Read the keyring to calculate the time of next rotation.
	keyring, err := c.raftReadTLSKeyring(ctx)
	if err != nil {
		return err
	}
	activeKey := keyring.GetActive()
	if activeKey == nil {
		return errors.New("no active raft TLS key found")
	}

	getNextRotationTime := func(next time.Time) time.Time {
		now := time.Now()

		// active key's CreatedTime + raftTLSRotationPeriod might be in
		// the past (meaning it is ready to be rotated) which will cause
		// NewTicker to panic when used with time.Until, prevent this by
		// pushing out rotation time into very near future
		if next.Before(now) {
			return now.Add(1 * time.Minute)
		}

		// push out to ensure proposed time does not elapse
		return next.Add(10 * time.Second)
	}

	// Start the process in a go routine
	go func() {
		nextRotationTime := getNextRotationTime(activeKey.CreatedTime.Add(raftTLSRotationPeriod))

		keyCheckInterval := time.NewTicker(1 * time.Minute)
		defer keyCheckInterval.Stop()

		// ticker is used to prevent memory leak of using time.After in
		// for - select pattern.
		ticker := time.NewTicker(time.Until(nextRotationTime))
		defer ticker.Stop()
		for {
			select {
			case <-keyCheckInterval.C:
				err := checkCommitted()
				if err != nil {
					logger.Error("failed to activate TLS key", "error", err)
				}
			case <-ticker.C:
				// It's time to rotate the keys
				next, err := rotateKeyring()
				if err != nil {
					logger.Error("failed to rotate TLS key", "error", err)
					nextRotationTime = time.Now().Add(10 * time.Second)
				} else {
					nextRotationTime = getNextRotationTime(next)
				}

				ticker.Reset(time.Until(nextRotationTime))

			case <-stopCh:
				return
			}
		}
	}()

	return nil
}

func (c *Core) raftReadTLSKeyring(ctx context.Context) (*raft.TLSKeyring, error) {
	tlsKeyringEntry, err := c.barrier.Get(ctx, raftTLSStoragePath)
	if err != nil {
		return nil, err
	}
	if tlsKeyringEntry == nil {
		return nil, errors.New("no keyring found")
	}
	var keyring raft.TLSKeyring
	if err := tlsKeyringEntry.DecodeJSON(&keyring); err != nil {
		return nil, err
	}

	return &keyring, nil
}

// raftCreateTLSKeyring creates the initial TLS key and the TLS Keyring for raft
// use. If a keyring entry is already present in storage, it will return an
// error.
func (c *Core) raftCreateTLSKeyring(ctx context.Context) (*raft.TLSKeyring, error) {
	if raftBackend := c.getRaftBackend(); raftBackend == nil {
		return nil, fmt.Errorf("raft backend not in use")
	}

	// Check if the keyring is already present
	raftTLSEntry, err := c.barrier.Get(ctx, raftTLSStoragePath)
	if err != nil {
		return nil, err
	}

	if raftTLSEntry != nil {
		return nil, fmt.Errorf("TLS keyring already present")
	}

	raftTLS, err := raft.GenerateTLSKey(c.secureRandomReader)
	if err != nil {
		return nil, err
	}

	keyring := &raft.TLSKeyring{
		Keys:        []*raft.TLSKey{raftTLS},
		ActiveKeyID: raftTLS.ID,
	}

	entry, err := logical.StorageEntryJSON(raftTLSStoragePath, keyring)
	if err != nil {
		return nil, err
	}
	if err := c.barrier.Put(ctx, entry); err != nil {
		return nil, err
	}
	return keyring, nil
}

func (c *Core) stopPeriodicRaftTLSRotate() {
	if c.raftTLSRotationStopCh != nil {
		close(c.raftTLSRotationStopCh)
	}
	c.raftTLSRotationStopCh = nil
	c.raftFollowerStates.Clear()
}

func (c *Core) checkRaftTLSKeyUpgrades(ctx context.Context) error {
	raftBackend := c.getRaftBackend()
	if raftBackend == nil {
		return nil
	}

	tlsKeyringEntry, err := c.barrier.Get(ctx, raftTLSStoragePath)
	if err != nil {
		return err
	}
	if tlsKeyringEntry == nil {
		return nil
	}

	var keyring raft.TLSKeyring
	if err := tlsKeyringEntry.DecodeJSON(&keyring); err != nil {
		return err
	}

	if err := raftBackend.SetTLSKeyring(&keyring); err != nil {
		return err
	}

	return nil
}

// handleSnapshotRestore is for the raft backend to hook back into core after a
// snapshot is restored so we can clear the necessary caches and handle changing
// keyrings or root keys
func (c *Core) raftSnapshotRestoreCallback(grabLock bool, sealNode bool) func(context.Context) error {
	return func(ctx context.Context) (retErr error) {
		c.logger.Info("running post snapshot restore invalidations")

		if grabLock {
			// Grab statelock
			l := newLockGrabber(c.stateLock.Lock, c.stateLock.Unlock, c.standbyStopCh.Load().(chan struct{}))
			go l.grab()
			if stopped := l.lockOrStop(); stopped {
				c.logger.Error("did not apply snapshot; vault is shutting down")
				return errors.New("did not apply snapshot; vault is shutting down")
			}
			defer c.stateLock.Unlock()
		}

		if sealNode {
			// If we failed to restore the snapshot we should seal this node as
			// it's in an unknown state
			defer func() {
				if retErr != nil {
					if err := c.sealInternalWithOptions(false, false, true); err != nil {
						c.logger.Error("failed to seal node", "error", err)
					}
				}
			}()
		}

		// Purge the cache so we make sure we are operating on fresh data
		c.physicalCache.Purge(ctx)

		// Reload the keyring in case it changed. If this fails it's likely
		// we've changed root keys.
		err := c.performKeyUpgrades(ctx)
		if err != nil {
			// The snapshot contained a root key or keyring we couldn't
			// recover
			switch c.seal.BarrierSealConfigType() {
			case SealConfigTypeShamir:
				// If we are a shamir seal we can't do anything. Just
				// seal all nodes.

				// Seal ourselves
				c.logger.Info("failed to perform key upgrades, sealing", "error", err)
				return err

			default:
				// If we are using an auto-unseal we can try to use the seal to
				// unseal again. If the auto-unseal mechanism has changed then
				// there isn't anything we can do but seal.
				c.logger.Info("failed to perform key upgrades, reloading using auto seal")
				keys, err := c.seal.GetStoredKeys(ctx)
				if err != nil {
					c.logger.Error("raft snapshot restore failed to get stored keys", "error", err)
					return err
				}
				if err := c.barrier.Seal(); err != nil {
					c.logger.Error("raft snapshot restore failed to seal barrier", "error", err)
					return err
				}
				if err := c.barrier.Unseal(ctx, keys[0]); err != nil {
					c.logger.Error("raft snapshot restore failed to unseal barrier", "error", err)
					return err
				}
				c.logger.Info("done reloading root key using auto seal")
			}
		}

		// Refresh the raft TLS keys
		if err := c.checkRaftTLSKeyUpgrades(ctx); err != nil {
			c.logger.Info("failed to perform TLS key upgrades, sealing", "error", err)
			return err
		}

		return nil
	}
}

func (c *Core) InitiateRetryJoin(ctx context.Context) error {
	raftBackend := c.getRaftBackend()
	if raftBackend == nil {
		return nil
	}

	if raftBackend.Initialized() {
		return nil
	}

	leaderInfos, err := raftBackend.JoinConfig()
	if err != nil {
		return err
	}

	// Nothing to do if config wasn't supplied
	if len(leaderInfos) == 0 {
		return nil
	}

	c.logger.Info("raft retry join initiated")

	if _, err = c.JoinRaftCluster(ctx, leaderInfos, raftBackend.NonVoter()); err != nil {
		return err
	}

	return nil
}

// getRaftChallenge is a helper function used by the raft join process for adding a
// node to a cluster: it contacts the given node and initiates the bootstrap
// challenge, returning the result or an error.
func (c *Core) getRaftChallenge(leaderInfo *raft.LeaderJoinInfo) (*raftInformation, error) {
	if leaderInfo == nil {
		return nil, errors.New("raft leader information is nil")
	}
	if len(leaderInfo.LeaderAPIAddr) == 0 {
		return nil, errors.New("raft leader address not provided")
	}

	c.logger.Info("attempting to join possible raft leader node", "leader_addr", leaderInfo.LeaderAPIAddr)

	// Create an API client to interact with the leader node
	transport := cleanhttp.DefaultPooledTransport()

	var err error
	if leaderInfo.TLSConfig == nil && (len(leaderInfo.LeaderCACert) != 0 || len(leaderInfo.LeaderClientCert) != 0 || len(leaderInfo.LeaderClientKey) != 0) {
		leaderInfo.TLSConfig, err = tlsutil.ClientTLSConfig([]byte(leaderInfo.LeaderCACert), []byte(leaderInfo.LeaderClientCert), []byte(leaderInfo.LeaderClientKey))
		if err != nil {
			return nil, fmt.Errorf("failed to create TLS config: %w", err)
		}
		leaderInfo.TLSConfig.ServerName = leaderInfo.LeaderTLSServerName
	}
	if leaderInfo.TLSConfig == nil && leaderInfo.LeaderTLSServerName != "" {
		leaderInfo.TLSConfig, err = tlsutil.SetupTLSConfig(map[string]string{"address": leaderInfo.LeaderTLSServerName}, "")
		if err != nil {
			return nil, fmt.Errorf("failed to create TLS config: %w", err)
		}
	}

	if leaderInfo.TLSConfig != nil {
		transport.TLSClientConfig = leaderInfo.TLSConfig.Clone()
		if err := http2.ConfigureTransport(transport); err != nil {
			return nil, fmt.Errorf("failed to configure TLS: %w", err)
		}
	}

	client := &http.Client{
		Transport: transport,
	}

	config := api.DefaultConfig()
	if config.Error != nil {
		return nil, fmt.Errorf("failed to create api client: %w", config.Error)
	}

	config.Address = leaderInfo.LeaderAPIAddr
	config.HttpClient = client
	config.MaxRetries = 0

	apiClient, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create api client: %w", err)
	}
	// Clearing namespace, as this client should only ever be using the root namespace
	apiClient.ClearNamespace()

	// Attempt to join the leader by requesting for the bootstrap challenge
	secret, err := apiClient.Logical().Write("sys/storage/raft/bootstrap/challenge", map[string]interface{}{
		"server_id": c.getRaftBackend().NodeID(),
	})
	if err != nil {
		return nil, fmt.Errorf("error during raft bootstrap init call: %w", err)
	}
	if secret == nil {
		return nil, errors.New("could not retrieve raft bootstrap package")
	}

	var sealConfig SealConfig
	err = mapstructure.Decode(secret.Data["seal_config"], &sealConfig)
	if err != nil {
		return nil, err
	}

	if sealConfig.Type != c.seal.BarrierSealConfigType().String() {
		return nil, fmt.Errorf("mismatching seal types between raft leader (%s) and follower (%s)", sealConfig.Type, c.seal.BarrierSealConfigType())
	}

	challengeB64, ok := secret.Data["challenge"]
	if !ok {
		return nil, errors.New("error during raft bootstrap call, no challenge given")
	}
	challengeRaw, err := base64.StdEncoding.DecodeString(challengeB64.(string))
	if err != nil {
		return nil, fmt.Errorf("error decoding raft bootstrap challenge: %w", err)
	}

	eBlob := &wrapping.BlobInfo{}
	if err := proto.Unmarshal(challengeRaw, eBlob); err != nil {
		return nil, fmt.Errorf("error decoding raft bootstrap challenge: %w", err)
	}

	return &raftInformation{
		challenge:           eBlob,
		leaderClient:        apiClient,
		leaderBarrierConfig: &sealConfig,
	}, nil
}

func (c *Core) JoinRaftCluster(ctx context.Context, leaderInfos []*raft.LeaderJoinInfo, nonVoter bool) (bool, error) {
	raftBackend := c.getRaftBackend()
	if raftBackend == nil {
		return false, errors.New("raft backend not in use")
	}

	init, err := c.InitializedLocally(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check if core is initialized: %w", err)
	}

	isRaftHAOnly := c.isRaftHAOnly()
	// Prevent join from happening if we're using raft for storage and
	// it has already been initialized.
	if init && !isRaftHAOnly {
		return true, nil
	}

	// Check on seal status and storage type before proceeding:
	// If raft is used for storage, core needs to be sealed
	if !isRaftHAOnly && !c.Sealed() {
		c.logger.Error("node must be sealed before joining")
		return false, errors.New("node must be sealed before joining")
	}

	// If raft is used for ha-only, core needs to be unsealed
	if isRaftHAOnly && c.Sealed() {
		c.logger.Error("node must be unsealed before joining")
		return false, errors.New("node must be unsealed before joining")
	}

	// Disallow leader API address to be provided if we're using raft for HA-only.
	// The leader API address is obtained directly through storage. This serves
	// as a form of verification that this node is sharing the same physical
	// storage as the leader node.
	if isRaftHAOnly {
		for _, info := range leaderInfos {
			if info.LeaderAPIAddr != "" || info.AutoJoin != "" {
				return false, errors.New("leader API address and auto-join metadata must be unset when raft is used exclusively for HA")
			}
		}

		// Get the leader address from storage
		keys, err := c.barrier.List(ctx, coreLeaderPrefix)
		if err != nil {
			return false, err
		}

		if len(keys) == 0 || len(keys[0]) == 0 {
			return false, errors.New("unable to fetch leadership entry")
		}

		leadershipEntry := coreLeaderPrefix + keys[0]
		entry, err := c.barrier.Get(ctx, leadershipEntry)
		if err != nil {
			return false, err
		}
		if entry == nil {
			return false, errors.New("unable to read leadership entry")
		}

		var adv activeAdvertisement
		err = jsonutil.DecodeJSON(entry.Value, &adv)
		if err != nil {
			return false, fmt.Errorf("unable to decoded leader entry: %w", err)
		}

		leaderInfos[0].LeaderAPIAddr = adv.RedirectAddr
	}

	disco, err := newDiscover()
	if err != nil {
		return false, fmt.Errorf("failed to create auto-join discovery: %w", err)
	}

	retryFailures := leaderInfos[0].Retry
	// answerChallenge performs the second part of a raft join: after we've issued
	// the sys/storage/raft/bootstrap/challenge call to initiate the join, this
	// func uses the seal to compute an answer to the challenge and sends it
	// back to the server that provided the challenge.
	answerChallenge := func(ctx context.Context, raftInfo *raftInformation) error {
		// If we're using Shamir and using raft for both physical and HA, we
		// need to block until the node is unsealed, unless retry is set to
		// false.
		if c.seal.BarrierSealConfigType() == SealConfigTypeShamir && !c.isRaftHAOnly() {
			c.raftInfo.Store(raftInfo)
			if err := c.seal.SetBarrierConfig(ctx, raftInfo.leaderBarrierConfig); err != nil {
				return err
			}

			if !retryFailures {
				return nil
			}

			// Wait until unseal keys are supplied
			raftInfo.joinInProgress = true
			c.raftInfo.Store(raftInfo)
			if atomic.LoadUint32(c.postUnsealStarted) != 1 {
				return errors.New("waiting for unseal keys to be supplied")
			}
		}

		raftInfo.nonVoter = nonVoter
		if err := c.joinRaftSendAnswer(ctx, c.seal.GetAccess(), raftInfo); err != nil {
			return fmt.Errorf("failed to send answer to raft leader node: %w", err)
		}

		if c.seal.BarrierSealConfigType() == SealConfigTypeShamir && !isRaftHAOnly {
			// Reset the state
			c.raftInfo.Store((*raftInformation)(nil))

			// In case of Shamir unsealing, inform the unseal process that raft join is completed
			close(c.raftJoinDoneCh)
		}

		c.logger.Info("successfully joined the raft cluster", "leader_addr", raftInfo.leaderClient.Address())
		return nil
	}

	// join attempts to join to any of the leaders defined in leaderInfos,
	// using the first one that returns a challenge to our request.  If shamir
	// seal is in use, we must wait to get enough unseal keys to solve the
	// challenge.  If we're unable to get a challenge from any leader, or if
	// we fail to answer the challenge successfully, or if ctx times out,
	// an error is returned.
	join := func() error {
		init, err := c.InitializedLocally(ctx)
		if err != nil {
			return fmt.Errorf("failed to check if core is initialized: %w", err)
		}

		// InitializedLocally will return non-nil before HA backends are
		// initialized. c.Initialized(ctx) checks InitializedLocally first, so
		// we can't use that generically for both cases. Instead check
		// raftBackend.Initialized() directly for the HA-Only case.
		if (!isRaftHAOnly && init) || (isRaftHAOnly && raftBackend.Initialized()) {
			c.logger.Info("returning from raft join as the node is initialized")
			return nil
		}

		if err := raftBackend.SetDesiredSuffrage(nonVoter); err != nil {
			c.logger.Error("failed to set desired suffrage for this node", "error", err)
			return nil
		}

		challengeCh := make(chan *raftInformation)
		var expandedJoinInfos []*raft.LeaderJoinInfo
		for _, leaderInfo := range leaderInfos {
			joinInfos, err := c.raftLeaderInfo(leaderInfo, disco)
			if err != nil {
				c.logger.Error("error in retry_join stanza, will not use it for raft join", "error", err,
					"leader_api_addr", leaderInfo.LeaderAPIAddr, "auto_join", leaderInfo.AutoJoin != "")
				continue
			}
			expandedJoinInfos = append(expandedJoinInfos, joinInfos...)
		}
		if err != nil {
			return err
		}
		var wg sync.WaitGroup
		for i := range expandedJoinInfos {
			wg.Add(1)
			go func(joinInfo *raft.LeaderJoinInfo) {
				defer wg.Done()
				raftInfo, err := c.getRaftChallenge(joinInfo)
				if err != nil {
					c.Logger().Error("failed to get raft challenge", "leader_addr", joinInfo.LeaderAPIAddr, "error", err)
					return
				}
				challengeCh <- raftInfo
			}(expandedJoinInfos[i])
		}
		go func() {
			wg.Wait()
			close(challengeCh)
		}()

		select {
		case <-ctx.Done():
			err = ctx.Err()
		case raftInfo := <-challengeCh:
			if raftInfo != nil {
				err = answerChallenge(ctx, raftInfo)
				if err == nil {
					return nil
				}
			} else {
				// Return an error so we can retry_join
				err = fmt.Errorf("failed to get raft challenge")
			}
		}
		return err
	}

	switch retryFailures {
	case true:
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				err := join()
				if err == nil {
					return
				}
				c.logger.Error("failed to retry join raft cluster", "retry", "2s", "err", err)
				time.Sleep(2 * time.Second)
			}
		}()

		// Backgrounded so return false
		return false, nil
	default:
		if err := join(); err != nil {
			c.logger.Error("failed to join raft cluster", "error", err)
			return false, fmt.Errorf("failed to join raft cluster: %w", err)
		}
	}

	return true, nil
}

// raftLeaderInfo uses go-discover to expand leaderInfo to include any auto-join results
func (c *Core) raftLeaderInfo(leaderInfo *raft.LeaderJoinInfo, disco *discover.Discover) ([]*raft.LeaderJoinInfo, error) {
	var ret []*raft.LeaderJoinInfo
	switch {
	case leaderInfo.LeaderAPIAddr != "" && leaderInfo.AutoJoin != "":
		return nil, errors.New("cannot provide both leader address and auto-join metadata")

	case leaderInfo.LeaderAPIAddr != "":
		ret = append(ret, leaderInfo)

	case leaderInfo.AutoJoin != "":
		scheme := leaderInfo.AutoJoinScheme
		if scheme == "" {
			// default to HTTPS when no scheme is provided
			scheme = "https"
		}
		port := leaderInfo.AutoJoinPort
		if port == 0 {
			// default to 8200 when no port is provided
			port = 8200
		}
		// Addrs returns either IPv4 or IPv6 address, without scheme or port
		clusterIPs, err := disco.Addrs(leaderInfo.AutoJoin, c.logger.StandardLogger(nil))
		if err != nil {
			return nil, fmt.Errorf("failed to parse addresses from auto-join metadata: %w", err)
		}
		for _, ip := range clusterIPs {
			if strings.Count(ip, ":") >= 2 && !strings.HasPrefix(ip, "[") {
				// An IPv6 address in implicit form, however we need it in explicit form to use in a URL.
				ip = fmt.Sprintf("[%s]", ip)
			}
			u := fmt.Sprintf("%s://%s:%d", scheme, ip, port)
			info := *leaderInfo
			info.LeaderAPIAddr = u
			ret = append(ret, &info)
		}

	default:
		return nil, errors.New("must provide leader address or auto-join metadata")
	}
	return ret, nil
}

// NewDelegateForCore creates a raft.Delegate for the specified core using its backend.
func NewDelegateForCore(c *Core) *raft.Delegate {
	return raft.NewDelegate(c.getRaftBackend())
}

// getRaftBackend returns the RaftBackend from the HA or physical backend,
// in that order of preference, or nil if not of type RaftBackend.
func (c *Core) getRaftBackend() *raft.RaftBackend {
	var raftBackend *raft.RaftBackend

	if raftHA, ok := c.ha.(*raft.RaftBackend); ok {
		raftBackend = raftHA
	}

	if raftStorage, ok := c.underlyingPhysical.(*raft.RaftBackend); ok {
		raftBackend = raftStorage
	}

	return raftBackend
}

// isRaftHAOnly returns true if c.ha is raft and physical storage is non-raft
func (c *Core) isRaftHAOnly() bool {
	_, isRaftHA := c.ha.(*raft.RaftBackend)
	_, isRaftStorage := c.underlyingPhysical.(*raft.RaftBackend)

	return isRaftHA && !isRaftStorage
}

func (c *Core) joinRaftSendAnswer(ctx context.Context, sealAccess seal.Access, raftInfo *raftInformation) error {
	if raftInfo.challenge == nil {
		return errors.New("raft challenge is nil")
	}

	raftBackend := c.getRaftBackend()
	if raftBackend == nil {
		return errors.New("raft backend is not in use")
	}

	if raftBackend.Initialized() {
		return errors.New("raft is already initialized")
	}

	multiWrapValue := &seal.MultiWrapValue{
		Generation: sealAccess.Generation(),
		Slots:      []*wrapping.BlobInfo{raftInfo.challenge},
	}
	plaintext, _, err := sealAccess.Decrypt(ctx, multiWrapValue, nil)
	if err != nil {
		return fmt.Errorf("error decrypting challenge: %w", err)
	}

	parsedClusterAddr, err := url.Parse(c.ClusterAddr())
	if err != nil {
		return fmt.Errorf("error parsing cluster address: %w", err)
	}
	clusterAddr := parsedClusterAddr.Host
	if atomic.LoadUint32(&TestingUpdateClusterAddr) == 1 && strings.HasSuffix(clusterAddr, ":0") {
		// We are testing and have an address provider, so just create a random
		// addr, it will be overwritten later.
		var err error
		clusterAddr, err = uuid.GenerateUUID()
		if err != nil {
			return err
		}
	}

	answerReq := raftInfo.leaderClient.NewRequest("PUT", "/v1/sys/storage/raft/bootstrap/answer")
	if err := answerReq.SetJSONBody(map[string]interface{}{
		"answer":       base64.StdEncoding.EncodeToString(plaintext),
		"cluster_addr": clusterAddr,
		"server_id":    raftBackend.NodeID(),
		"non_voter":    raftInfo.nonVoter,
	}); err != nil {
		return err
	}

	answerRespJson, err := raftInfo.leaderClient.RawRequestWithContext(ctx, answerReq)
	if answerRespJson != nil {
		defer answerRespJson.Body.Close()
	}
	if err != nil {
		return err
	}

	var answerResp answerRespData
	if err := jsonutil.DecodeJSONFromReader(answerRespJson.Body, &answerResp); err != nil {
		return err
	}

	if answerResp.Data.AutoloadedLicense && !LicenseAutoloaded(c) {
		return ErrJoinWithoutAutoloading
	}
	if err := raftBackend.Bootstrap(answerResp.Data.Peers); err != nil {
		return err
	}

	err = c.startClusterListener(ctx)
	if err != nil {
		return fmt.Errorf("error starting cluster: %w", err)
	}

	raftBackend.SetRestoreCallback(c.raftSnapshotRestoreCallback(true, true))
	opts := raft.SetupOpts{
		TLSKeyring:      answerResp.Data.TLSKeyring,
		ClusterListener: c.getClusterListener(),
	}
	err = raftBackend.SetupCluster(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to setup raft cluster: %w", err)
	}

	return nil
}

func (c *Core) loadAutopilotConfiguration(ctx context.Context) (*raft.AutopilotConfig, error) {
	var autopilotConfig *raft.AutopilotConfig
	entry, err := c.barrier.Get(ctx, raftAutopilotConfigurationStoragePath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	if err := jsonutil.DecodeJSON(entry.Value, &autopilotConfig); err != nil {
		return nil, err
	}

	return autopilotConfig, nil
}

// RaftBootstrap performs bootstrapping of a raft cluster if core contains a raft
// backend. If raft is not part for the storage or HA storage backend, this
// call results in an error.
func (c *Core) RaftBootstrap(ctx context.Context, onInit bool) error {
	if c.logger.IsDebug() {
		c.logger.Debug("bootstrapping raft backend")
		defer c.logger.Debug("finished bootstrapping raft backend")
	}

	raftBackend := c.getRaftBackend()
	if raftBackend == nil {
		return errors.New("raft backend not in use")
	}

	parsedClusterAddr, err := url.Parse(c.ClusterAddr())
	if err != nil {
		return fmt.Errorf("error parsing cluster address: %w", err)
	}
	if err := raftBackend.Bootstrap([]raft.Peer{
		{
			ID:      raftBackend.NodeID(),
			Address: parsedClusterAddr.Host,
		},
	}); err != nil {
		return fmt.Errorf("could not bootstrap clustered storage: %w", err)
	}

	raftOpts := raft.SetupOpts{
		StartAsLeader: true,
	}

	if !onInit {
		// Generate the TLS Keyring info for SetupCluster to consume
		raftTLS, err := c.raftCreateTLSKeyring(ctx)
		if err != nil {
			return fmt.Errorf("could not generate TLS keyring during bootstrap: %w", err)
		}

		raftBackend.SetRestoreCallback(c.raftSnapshotRestoreCallback(true, true))
		raftOpts.ClusterListener = c.getClusterListener()

		raftOpts.TLSKeyring = raftTLS
	}

	if err := raftBackend.SetupCluster(ctx, raftOpts); err != nil {
		return fmt.Errorf("could not start clustered storage: %w", err)
	}

	return nil
}

func (c *Core) isRaftUnseal() bool {
	return c.raftInfo.Load().(*raftInformation) != nil
}

type answerRespData struct {
	Data answerResp `json:"data"`
}

type answerResp struct {
	Peers             []raft.Peer      `json:"peers"`
	TLSKeyring        *raft.TLSKeyring `json:"tls_keyring"`
	AutoloadedLicense bool             `json:"autoloaded_license"`
}

func newDiscover() (*discover.Discover, error) {
	providers := make(map[string]discover.Provider)
	for k, v := range discover.Providers {
		providers[k] = v
	}

	providers["k8s"] = &discoverk8s.Provider{}

	return discover.New(
		discover.WithProviders(providers),
	)
}
