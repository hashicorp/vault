package vault

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/http2"
)

var (
	raftTLSStoragePath    = "core/raft/tls"
	raftTLSRotationPeriod = 24 * time.Hour
)

type raftFollowerStates struct {
	l         sync.RWMutex
	followers map[string]uint64
}

func (s *raftFollowerStates) update(nodeID string, appliedIndex uint64) {
	s.l.Lock()
	s.followers[nodeID] = appliedIndex
	s.l.Unlock()
}
func (s *raftFollowerStates) delete(nodeID string) {
	s.l.RLock()
	delete(s.followers, nodeID)
	s.l.RUnlock()
}
func (s *raftFollowerStates) get(nodeID string) uint64 {
	s.l.RLock()
	index := s.followers[nodeID]
	s.l.RUnlock()
	return index
}
func (s *raftFollowerStates) minIndex() uint64 {
	var min uint64 = math.MaxUint64
	minFunc := func(a, b uint64) uint64 {
		if a > b {
			return b
		}
		return a
	}

	s.l.RLock()
	for _, i := range s.followers {
		min = minFunc(min, i)
	}
	s.l.RUnlock()

	if min == math.MaxUint64 {
		return 0
	}

	return min
}

// startRaftStorage will call SetupCluster in the raft backend which starts raft
// up and enables the cluster handler.
func (c *Core) startRaftStorage(ctx context.Context) (retErr error) {
	raftStorage, ok := c.underlyingPhysical.(*raft.RaftBackend)
	if !ok {
		return nil
	}
	if raftStorage.Initialized() {
		return nil
	}

	// Retrieve the raft TLS information
	raftTLSEntry, err := c.barrier.Get(ctx, raftTLSStoragePath)
	if err != nil {
		return err
	}

	var creating bool
	var raftTLS *raft.TLSKeyring
	switch raftTLSEntry {
	case nil:
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

	raftStorage.SetRestoreCallback(c.raftSnapshotRestoreCallback(true, true))
	if err := raftStorage.SetupCluster(ctx, raft.SetupOpts{
		TLSKeyring:      raftTLS,
		ClusterListener: c.getClusterListener(),
		StartAsLeader:   creating,
	}); err != nil {
		return err
	}

	defer func() {
		if retErr != nil {
			c.logger.Info("stopping raft server")
			if err := raftStorage.TeardownCluster(c.getClusterListener()); err != nil {
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

func (c *Core) setupRaftActiveNode(ctx context.Context) error {
	c.pendingRaftPeers = make(map[string][]byte)
	return c.startPeriodicRaftTLSRotate(ctx)
}

func (c *Core) stopRaftActiveNode() {
	c.pendingRaftPeers = nil
	c.stopPeriodicRaftTLSRotate()
}

// startPeriodicRaftTLSRotate will spawn a go routine in charge of periodically
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
func (c *Core) startPeriodicRaftTLSRotate(ctx context.Context) error {
	raftStorage, ok := c.underlyingPhysical.(*raft.RaftBackend)
	if !ok {
		return nil
	}

	stopCh := make(chan struct{})
	followerStates := &raftFollowerStates{
		followers: make(map[string]uint64),
	}

	// Pre-populate the follower list with the set of peers.
	raftConfig, err := raftStorage.GetConfiguration(ctx)
	if err != nil {
		return err
	}
	for _, server := range raftConfig.Servers {
		if server.NodeID != raftStorage.NodeID() {
			followerStates.update(server.NodeID, 0)
		}
	}

	logger := c.logger.Named("raft")
	c.raftTLSRotationStopCh = stopCh
	c.raftFollowerStates = followerStates

	readKeyring := func() (*raft.TLSKeyring, error) {
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

	// rotateKeyring writes new key data to the keyring and adds an applied
	// index that is used to verify it has been committed. The keys written in
	// this function can be used on standbys but the active node doesn't start
	// using it yet.
	rotateKeyring := func() (time.Time, error) {
		// Read the existing keyring
		keyring, err := readKeyring()
		if err != nil {
			return time.Time{}, errwrap.Wrapf("failed to read raft TLS keyring: {{err}}", err)
		}

		switch {
		case len(keyring.Keys) == 2 && keyring.Keys[1].AppliedIndex == 0:
			// If this case is hit then the second write to add the applied
			// index failed. Attempt to write it again.
			keyring.Keys[1].AppliedIndex = raftStorage.AppliedIndex()
			keyring.AppliedIndex = raftStorage.AppliedIndex()
			entry, err := logical.StorageEntryJSON(raftTLSStoragePath, keyring)
			if err != nil {
				return time.Time{}, errwrap.Wrapf("failed to json encode keyring: {{err}}", err)
			}
			if err := c.barrier.Put(ctx, entry); err != nil {
				return time.Time{}, errwrap.Wrapf("failed to write keyring: {{err}}", err)
			}

		case len(keyring.Keys) > 1:
			// If there already exists a pending key update then the update
			// hasn't replicated down to all standby nodes yet. Don't allow any
			// new keys to be created until all standbys have seen this previous
			// rotation. As a backoff strategy another rotation attempt is
			// scheduled for 5 minutes from now.
			logger.Warn("skipping new raft TLS config creation, keys are pending")
			return time.Now().Add(time.Minute * 5), nil
		}

		logger.Info("creating new raft TLS config")

		// Create a new key
		raftTLSKey, err := raft.GenerateTLSKey(c.secureRandomReader)
		if err != nil {
			return time.Time{}, errwrap.Wrapf("failed to generate new raft TLS key: {{err}}", err)
		}

		// Advance the term and store the new key
		keyring.Term += 1
		keyring.Keys = append(keyring.Keys, raftTLSKey)
		entry, err := logical.StorageEntryJSON(raftTLSStoragePath, keyring)
		if err != nil {
			return time.Time{}, errwrap.Wrapf("failed to json encode keyring: {{err}}", err)
		}
		if err := c.barrier.Put(ctx, entry); err != nil {
			return time.Time{}, errwrap.Wrapf("failed to write keyring: {{err}}", err)
		}

		// Write the keyring again with the new applied index. This allows us to
		// track if standby nodes receive the update.
		keyring.Keys[1].AppliedIndex = raftStorage.AppliedIndex()
		keyring.AppliedIndex = raftStorage.AppliedIndex()
		entry, err = logical.StorageEntryJSON(raftTLSStoragePath, keyring)
		if err != nil {
			return time.Time{}, errwrap.Wrapf("failed to json encode keyring: {{err}}", err)
		}
		if err := c.barrier.Put(ctx, entry); err != nil {
			return time.Time{}, errwrap.Wrapf("failed to write keyring: {{err}}", err)
		}

		logger.Info("wrote new raft TLS config")
		// Schedule the next rotation
		return raftTLSKey.CreatedTime.Add(raftTLSRotationPeriod), nil
	}

	// checkCommitted verifies key updates have been applied to all nodes and
	// finalizes the rotation by deleting the old keys and updating the raft
	// backend.
	checkCommitted := func() error {
		keyring, err := readKeyring()
		if err != nil {
			return errwrap.Wrapf("failed to read raft TLS keyring: {{err}}", err)
		}

		switch {
		case len(keyring.Keys) == 1:
			// No Keys to apply
			return nil
		case keyring.Keys[1].AppliedIndex != keyring.AppliedIndex:
			// We haven't fully committed the new key, continue here
			return nil
		case followerStates.minIndex() < keyring.AppliedIndex:
			// Not all the followers have applied the latest key
			return nil
		}

		// Upgrade to the new key
		keyring.Keys = keyring.Keys[1:]
		keyring.ActiveKeyID = keyring.Keys[0].ID
		keyring.Term += 1
		entry, err := logical.StorageEntryJSON(raftTLSStoragePath, keyring)
		if err != nil {
			return errwrap.Wrapf("failed to json encode keyring: {{err}}", err)
		}
		if err := c.barrier.Put(ctx, entry); err != nil {
			return errwrap.Wrapf("failed to write keyring: {{err}}", err)
		}

		// Update the TLS Key in the backend
		if err := raftStorage.SetTLSKeyring(keyring); err != nil {
			return errwrap.Wrapf("failed to install keyring: {{err}}", err)
		}

		logger.Info("installed new raft TLS key", "term", keyring.Term)
		return nil
	}

	// Read the keyring to calculate the time of next rotation.
	keyring, err := readKeyring()
	if err != nil {
		return err
	}
	activeKey := keyring.GetActive()
	if activeKey == nil {
		return errors.New("no active raft TLS key found")
	}

	// Start the process in a go routine
	go func() {
		nextRotationTime := activeKey.CreatedTime.Add(raftTLSRotationPeriod)

		keyCheckInterval := time.NewTicker(1 * time.Minute)
		defer keyCheckInterval.Stop()

		var backoff bool
		for {
			// If we encountered and error we should try to create the key
			// again.
			if backoff {
				nextRotationTime = time.Now().Add(10 * time.Second)
				backoff = false
			}

			select {
			case <-keyCheckInterval.C:
				err := checkCommitted()
				if err != nil {
					logger.Error("failed to activate TLS key", "error", err)
				}
			case <-time.After(time.Until(nextRotationTime)):
				// It's time to rotate the keys
				next, err := rotateKeyring()
				if err != nil {
					logger.Error("failed to rotate TLS key", "error", err)
					backoff = true
					continue
				}

				nextRotationTime = next

			case <-stopCh:
				return
			}
		}
	}()

	return nil
}

func (c *Core) createRaftTLSKeyring(ctx context.Context) error {
	if _, ok := c.underlyingPhysical.(*raft.RaftBackend); !ok {
		return nil
	}

	raftTLS, err := raft.GenerateTLSKey(c.secureRandomReader)
	if err != nil {
		return err
	}

	keyring := &raft.TLSKeyring{
		Keys:        []*raft.TLSKey{raftTLS},
		ActiveKeyID: raftTLS.ID,
	}

	entry, err := logical.StorageEntryJSON(raftTLSStoragePath, keyring)
	if err != nil {
		return err
	}
	if err := c.barrier.Put(ctx, entry); err != nil {
		return err
	}
	return nil
}

func (c *Core) stopPeriodicRaftTLSRotate() {
	if c.raftTLSRotationStopCh != nil {
		close(c.raftTLSRotationStopCh)
	}
	c.raftTLSRotationStopCh = nil
	c.raftFollowerStates = nil
}

func (c *Core) checkRaftTLSKeyUpgrades(ctx context.Context) error {
	raftStorage, ok := c.underlyingPhysical.(*raft.RaftBackend)
	if !ok {
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

	if err := raftStorage.SetTLSKeyring(&keyring); err != nil {
		return err
	}

	return nil
}

// handleSnapshotRestore is for the raft backend to hook back into core after a
// snapshot is restored so we can clear the necessary caches and handle changing
// keyrings or master keys
func (c *Core) raftSnapshotRestoreCallback(grabLock bool, sealNode bool) func(context.Context) error {
	return func(ctx context.Context) (retErr error) {
		c.logger.Info("running post snapshot restore invalidations")

		if grabLock {
			// Grab statelock
			if stopped := grabLockOrStop(c.stateLock.Lock, c.stateLock.Unlock, c.standbyStopCh); stopped {
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

		// Refresh the raft TLS keys
		if err := c.checkRaftTLSKeyUpgrades(ctx); err != nil {
			c.logger.Info("failed to perform TLS key upgrades, sealing", "error", err)
			return err
		}

		// Reload the keyring in case it changed. If this fails it's likely
		// we've changed master keys.
		err := c.performKeyUpgrades(ctx)
		if err != nil {
			// The snapshot contained a master key or keyring we couldn't
			// recover
			switch c.seal.BarrierType() {
			case wrapping.Shamir:
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
				c.logger.Info("done reloading master key using auto seal")
			}
		}

		return nil
	}
}

func (c *Core) InitiateRetryJoin(ctx context.Context) error {
	raftStorage, ok := c.underlyingPhysical.(*raft.RaftBackend)
	if !ok {
		return errors.New("raft storage not configured")
	}

	if raftStorage.Initialized() {
		return nil
	}

	leaderInfos, err := raftStorage.JoinConfig()
	if err != nil {
		return err
	}

	// Nothing to do if config wasn't supplied
	if len(leaderInfos) == 0 {
		return nil
	}

	c.logger.Info("raft retry join initiated")

	if _, err = c.JoinRaftCluster(ctx, leaderInfos, false); err != nil {
		return err
	}

	return nil
}

func (c *Core) JoinRaftCluster(ctx context.Context, leaderInfos []*raft.LeaderJoinInfo, nonVoter bool) (bool, error) {
	raftStorage, ok := c.underlyingPhysical.(*raft.RaftBackend)
	if !ok {
		return false, errors.New("raft storage not configured")
	}

	if raftStorage.Initialized() {
		return false, errors.New("raft storage is already initialized")
	}

	init, err := c.Initialized(ctx)
	if err != nil {
		return false, errwrap.Wrapf("failed to check if core is initialized: {{err}}", err)
	}
	if init {
		return true, nil
	}

	join := func(retry bool) error {
		joinLeader := func(leaderInfo *raft.LeaderJoinInfo) error {
			if leaderInfo == nil {
				return errors.New("raft leader information is nil")
			}
			if len(leaderInfo.LeaderAPIAddr) == 0 {
				return errors.New("raft leader address not provided")
			}

			init, err := c.Initialized(ctx)
			if err != nil {
				return errwrap.Wrapf("failed to check if core is initialized: {{err}}", err)
			}
			if init {
				c.logger.Info("returning from raft join as the node is initialized")
				return nil
			}
			if !c.Sealed() {
				c.logger.Info("returning from raft join as the node is unsealed")
				return nil
			}

			c.logger.Info("attempting to join possible raft leader node", "leader_addr", leaderInfo.LeaderAPIAddr)

			// Create an API client to interact with the leader node
			transport := cleanhttp.DefaultPooledTransport()

			if leaderInfo.TLSConfig == nil && (len(leaderInfo.LeaderCACert) != 0 || len(leaderInfo.LeaderClientCert) != 0 || len(leaderInfo.LeaderClientKey) != 0) {
				leaderInfo.TLSConfig, err = tlsutil.ClientTLSConfig([]byte(leaderInfo.LeaderCACert), []byte(leaderInfo.LeaderClientCert), []byte(leaderInfo.LeaderClientKey))
				if err != nil {
					return errwrap.Wrapf("failed to create TLS config: {{err}}", err)
				}
			}

			if leaderInfo.TLSConfig != nil {
				transport.TLSClientConfig = leaderInfo.TLSConfig.Clone()
				if err := http2.ConfigureTransport(transport); err != nil {
					return errwrap.Wrapf("failed to configure TLS: {{err}}", err)
				}
			}

			client := &http.Client{
				Transport: transport,
			}
			config := api.DefaultConfig()
			if config.Error != nil {
				return errwrap.Wrapf("failed to create api client: {{err}}", config.Error)
			}
			config.Address = leaderInfo.LeaderAPIAddr
			config.HttpClient = client
			config.MaxRetries = 0
			apiClient, err := api.NewClient(config)
			if err != nil {
				return errwrap.Wrapf("failed to create api client: {{err}}", err)
			}

			// Attempt to join the leader by requesting for the bootstrap challenge
			secret, err := apiClient.Logical().Write("sys/storage/raft/bootstrap/challenge", map[string]interface{}{
				"server_id": raftStorage.NodeID(),
			})
			if err != nil {
				return errwrap.Wrapf("error during raft bootstrap init call: {{err}}", err)
			}
			if secret == nil {
				return errors.New("could not retrieve raft bootstrap package")
			}

			var sealConfig SealConfig
			err = mapstructure.Decode(secret.Data["seal_config"], &sealConfig)
			if err != nil {
				return err
			}

			if sealConfig.Type != c.seal.BarrierType() {
				return fmt.Errorf("mismatching seal types between raft leader (%s) and follower (%s)", sealConfig.Type, c.seal.BarrierType())
			}

			challengeB64, ok := secret.Data["challenge"]
			if !ok {
				return errors.New("error during raft bootstrap call, no challenge given")
			}
			challengeRaw, err := base64.StdEncoding.DecodeString(challengeB64.(string))
			if err != nil {
				return errwrap.Wrapf("error decoding raft bootstrap challenge: {{err}}", err)
			}

			eBlob := &wrapping.EncryptedBlobInfo{}
			if err := proto.Unmarshal(challengeRaw, eBlob); err != nil {
				return errwrap.Wrapf("error decoding raft bootstrap challenge: {{err}}", err)
			}
			raftInfo := &raftInformation{
				challenge:           eBlob,
				leaderClient:        apiClient,
				leaderBarrierConfig: &sealConfig,
				nonVoter:            nonVoter,
			}

			if c.seal.BarrierType() == wrapping.Shamir {
				c.raftInfo = raftInfo
				if err := c.seal.SetBarrierConfig(ctx, &sealConfig); err != nil {
					return err
				}

				if !retry {
					return nil
				}

				// Wait until unseal keys are supplied
				c.raftInfo.joinInProgress = true
				if atomic.LoadUint32(c.postUnsealStarted) != 1 {
					return errors.New("waiting for unseal keys to be supplied")
				}
			}

			if err := c.joinRaftSendAnswer(ctx, c.seal.GetAccess(), raftInfo); err != nil {
				return errwrap.Wrapf("failed to send answer to raft leader node: {{err}}", err)
			}

			if c.seal.BarrierType() == wrapping.Shamir {
				// Reset the state
				c.raftInfo = nil

				// In case of Shamir unsealing, inform the unseal process that raft join is completed
				close(c.raftJoinDoneCh)
			}

			c.logger.Info("successfully joined the raft cluster", "leader_addr", leaderInfo.LeaderAPIAddr)
			return nil
		}

		// Each join try goes through all the possible leader nodes and attempts to join
		// them, until one of the attempt succeeds.
		for _, leaderInfo := range leaderInfos {
			err = joinLeader(leaderInfo)
			if err == nil {
				return nil
			}
			c.logger.Info("join attempt failed", "error", err)
		}

		return errors.New("failed to join any raft leader node")
	}

	switch leaderInfos[0].Retry {
	case true:
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				err := join(true)
				if err == nil {
					return
				}
				c.logger.Error("failed to retry join raft cluster", "retry", "2s")
				time.Sleep(2 * time.Second)
			}
		}()

		// Backgrounded so return false
		return false, nil
	default:
		if err := join(false); err != nil {
			c.logger.Error("failed to join raft cluster", "error", err)
			return false, errwrap.Wrapf("failed to join raft cluster: {{err}}", err)
		}
	}

	return true, nil
}

// This is used in tests to override the cluster address
var UpdateClusterAddrForTests uint32

func (c *Core) joinRaftSendAnswer(ctx context.Context, sealAccess *seal.Access, raftInfo *raftInformation) error {
	if raftInfo.challenge == nil {
		return errors.New("raft challenge is nil")
	}

	raftStorage, ok := c.underlyingPhysical.(*raft.RaftBackend)
	if !ok {
		return errors.New("raft storage not in use")
	}

	if raftStorage.Initialized() {
		return errors.New("raft is already initialized")
	}

	plaintext, err := sealAccess.Decrypt(ctx, raftInfo.challenge, nil)
	if err != nil {
		return errwrap.Wrapf("error decrypting challenge: {{err}}", err)
	}

	parsedClusterAddr, err := url.Parse(c.ClusterAddr())
	if err != nil {
		return errwrap.Wrapf("error parsing cluster address: {{err}}", err)
	}
	clusterAddr := parsedClusterAddr.Host
	if atomic.LoadUint32(&UpdateClusterAddrForTests) == 1 && strings.HasSuffix(clusterAddr, ":0") {
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
		"server_id":    raftStorage.NodeID(),
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

	raftStorage.Bootstrap(ctx, answerResp.Data.Peers)

	err = c.startClusterListener(ctx)
	if err != nil {
		return errwrap.Wrapf("error starting cluster: {{err}}", err)
	}

	raftStorage.SetRestoreCallback(c.raftSnapshotRestoreCallback(true, true))
	err = raftStorage.SetupCluster(ctx, raft.SetupOpts{
		TLSKeyring:      answerResp.Data.TLSKeyring,
		ClusterListener: c.getClusterListener(),
	})
	if err != nil {
		return errwrap.Wrapf("failed to setup raft cluster: {{err}}", err)
	}

	return nil
}

func (c *Core) isRaftUnseal() bool {
	return c.raftInfo != nil
}

type answerRespData struct {
	Data answerResp `json:"data"`
}

type answerResp struct {
	Peers      []raft.Peer      `json:"peers"`
	TLSKeyring *raft.TLSKeyring `json:"tls_keyring"`
}
