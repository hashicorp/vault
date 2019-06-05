package raft

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	proto "github.com/golang/protobuf/proto"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/vault/physical/raft/logstore"
	"github.com/hashicorp/vault/sdk/helper/consts"

	physicalstd "github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/sdk/physical"
)

// Verify RaftBackend satisfies the correct interfaces
var _ physical.Backend = (*RaftBackend)(nil)
var _ physical.Transactional = (*RaftBackend)(nil)
var _ physicalstd.Clustered = (*RaftBackend)(nil)

var (
	// raftLogCacheSize is the maximum number of logs to cache in-memory.
	// This is used to reduce disk I/O for the recently committed entries.
	raftLogCacheSize = 512

	raftState         = "raft/"
	snapshotsRetained = 2
)

// RaftBackend implements the backend interfaces and uses the raft protocol to
// persist writes to the FSM.
type RaftBackend struct {
	logger log.Logger
	conf   map[string]string
	l      sync.RWMutex

	// fsm is the state store for vault's data
	fsm *FSM

	// raft is the instance of raft we will operate on.
	raft *raft.Raft

	// raftNotifyCh is used to receive updates about leadership changes
	// regarding this node.
	raftNotifyCh chan bool

	// raftLayer is the network layer used to connect the nodes in the raft
	// cluster.
	raftLayer *raftLayer

	// raftTransport is the transport layer that the raft library uses for RPC
	// communication.
	raftTransport raft.Transport

	// snapStore is our snapshot mechanism.
	snapStore raft.SnapshotStore

	// logStore is used by the raft library to store the raft logs in durable
	// storage.
	logStore raft.LogStore

	// stableStore is used by the raft library to store additional metadata in
	// durable storage.
	stableStore raft.StableStore

	// bootstrapConfig is only set when this node needs to be bootstrapped upon
	// startup.
	bootstrapConfig *raft.Configuration

	// dataDir is the location on the local filesystem that raft and FSM data
	// will be stored.
	dataDir string

	// localID is the ID for this node. This can either be configured in the
	// config file, via a file on disk, or is otherwise randomly generated.
	localID string
}

// EnsurePath is used to make sure a path exists
func EnsurePath(path string, dir bool) error {
	if !dir {
		path = filepath.Dir(path)
	}
	return os.MkdirAll(path, 0755)
}

// NewRaftBackend constructs a RaftBackend using the given directory
func NewRaftBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// Create the FSM.
	var err error
	fsm, err := NewFSM(conf, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create fsm: %v", err)
	}

	path, ok := conf["path"]
	if !ok {
		return nil, fmt.Errorf("'path' must be set")
	}

	/*var serverAddressProvider raft.ServerAddressProvider = nil
	if s.config.RaftConfig.ProtocolVersion >= 3 { //ServerAddressProvider needs server ids to work correctly, which is only supported in protocol version 3 or higher
		serverAddressProvider = s.serverLookup
	}*/

	// Build an all in-memory setup for dev mode, otherwise prepare a full
	// disk-based setup.
	var log raft.LogStore
	var stable raft.StableStore
	var snap raft.SnapshotStore
	var devMode bool
	if devMode {
		store := raft.NewInmemStore()
		//raftInmem = store
		stable = store
		log = store
		snap = raft.NewInmemSnapshotStore()
	} else {
		// Create the base raft path.
		path := filepath.Join(path, raftState)
		if err := EnsurePath(path, true); err != nil {
			return nil, err
		}

		// Create the backend raft store for logs and stable storage.
		store, err := raftboltdb.NewBoltStore(filepath.Join(path, "raft.db"))
		if err != nil {
			return nil, err
		}
		stable = store

		// Wrap the store in a LogCache to improve performance.
		cacheStore, err := raft.NewLogCache(raftLogCacheSize, store)
		if err != nil {
			return nil, err
		}
		log = cacheStore

		// Create the snapshot store.
		snapshots, err := NewBoltSnapshotStore(path, snapshotsRetained, logger.Named("snapshot"), fsm)
		if err != nil {
			return nil, err
		}
		snap = snapshots
	}

	var localID string
	{
		// Determine the local node ID
		localID = conf["node_id"]

		// If not set in the config check the "node-id" file.
		if len(localID) == 0 {
			localIDRaw, err := ioutil.ReadFile(filepath.Join(path, "node-id"))
			switch {
			case err == nil:
				if len(localIDRaw) > 0 {
					localID = string(localIDRaw)
				}
			case os.IsNotExist(err):
			default:
				return nil, err
			}
		}

		// If the file didn't exist generate a UUID and persist it to tne
		// "node-id" file.
		if len(localID) == 0 {
			id, err := uuid.GenerateUUID()
			if err != nil {
				return nil, err
			}

			if err := ioutil.WriteFile(filepath.Join(path, "node-id"), []byte(id), 0600); err != nil {
				return nil, err
			}

			localID = id
		}
	}

	return &RaftBackend{
		logger:      logger,
		fsm:         fsm,
		conf:        conf,
		logStore:    log,
		stableStore: stable,
		snapStore:   snap,
		dataDir:     path,
		localID:     localID,
	}, nil
}

// RaftServer has information about a server in the Raft configuration
type RaftServer struct {
	// NodeID is the name of the server
	NodeID string `json:"node_id"`

	// Address is the IP:port of the server, used for Raft communications
	Address string `json:"address"`

	// Leader is true if this server is the current cluster leader
	Leader bool `json:"leader"`

	// Protocol version is the raft protocol version used by the server
	ProtocolVersion string `json:"protocol_version"`

	// Voter is true if this server has a vote in the cluster. This might
	// be false if the server is staging and still coming online.
	Voter bool `json:"voter"`
}

// RaftConfigurationResponse is returned when querying for the current Raft
// configuration.
type RaftConfigurationResponse struct {
	// Servers has the list of servers in the Raft configuration.
	Servers []*RaftServer `json:"servers"`

	// Index has the Raft index of this configuration.
	Index uint64 `json:"index"`
}

// Peer defines the ID and Adress for a given member of the raft cluster.
type Peer struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

// NodeID returns the identifier of the node
func (b *RaftBackend) NodeID() string {
	return b.localID
}

// Initialized tells if raft is running or not
func (b *RaftBackend) Initialized() bool {
	return b.raft != nil
}

// Bootstrap prepares the given peers to be part of the raft cluster
func (b *RaftBackend) Bootstrap(ctx context.Context, peers []Peer) error {
	b.l.Lock()
	defer b.l.Unlock()

	hasState, err := raft.HasExistingState(b.logStore, b.stableStore, b.snapStore)
	if err != nil {
		return err
	}

	if hasState {
		return errors.New("error bootstrapping cluster: cluster already has state")
	}

	raftConfig := &raft.Configuration{
		Servers: make([]raft.Server, len(peers)),
	}

	for i, p := range peers {
		raftConfig.Servers[i] = raft.Server{
			ID:      raft.ServerID(p.ID),
			Address: raft.ServerAddress(p.Address),
		}
	}

	// Store the config for later use
	b.bootstrapConfig = raftConfig
	return nil
}

// SetupCluster starts the raft cluster and enables the networking needed for
// the raft nodes to communicate.
func (b *RaftBackend) SetupCluster(ctx context.Context, networkConfig *physicalstd.NetworkConfig, clusterListener physicalstd.ClusterHook) error {
	b.logger.Trace("setting up raft cluster")

	b.l.Lock()
	defer b.l.Unlock()

	// We are already unsealed
	if b.raft != nil {
		b.logger.Debug("raft already started, not setting up cluster")
		return nil
	}

	if len(b.localID) == 0 {
		return errors.New("no local node id configured")
	}

	// Setup the raft config
	raftConfig := raft.DefaultConfig()
	raftConfig.Logger = b.logger
	raftConfig.SnapshotThreshold = 100
	raftConfig.TrailingLogs = 100
	raftConfig.ElectionTimeout = raftConfig.ElectionTimeout * 5
	raftConfig.HeartbeatTimeout = raftConfig.HeartbeatTimeout * 5
	raftConfig.LeaderLeaseTimeout = raftConfig.LeaderLeaseTimeout * 5

	switch {
	case networkConfig == nil:
		// If we don't have a provided network config we use an in-memory one.
		// This allows us to bootstrap a node without bringing up a cluster
		// network. This will be true during bootstrap and dev modes.
		_, b.raftTransport = raft.NewInmemTransport(raft.ServerAddress(clusterListener.Addr().String()))
	default:
		// Load the base TLS config from the cluster listener.
		baseTLSConfig, err := clusterListener.TLSConfig(ctx)
		if err != nil {
			return err
		}

		// Set the local address and localID in the streaming layer and the raft config.
		raftLayer, err := NewRaftLayer(b.logger.Named("stream"), networkConfig, baseTLSConfig)
		if err != nil {
			return err
		}
		transConfig := &raft.NetworkTransportConfig{
			Stream:  raftLayer,
			MaxPool: 3,
			Timeout: 10 * time.Second,
			//	ServerAddressProvider: serverAddressProvider,
		}
		transport := raft.NewNetworkTransportWithConfig(transConfig)

		b.raftLayer = raftLayer
		b.raftTransport = transport
	}

	raftConfig.LocalID = raft.ServerID(b.localID)

	// Set up a channel for reliable leader notifications.
	raftNotifyCh := make(chan bool, 1)
	raftConfig.NotifyCh = raftNotifyCh

	// If we have a bootstrapConfig set we should bootstrap now.
	if b.bootstrapConfig != nil {
		bootstrapConfig := b.bootstrapConfig
		// Unset the bootstrap config
		b.bootstrapConfig = nil

		// Bootstrap raft with our known cluster members.
		if err := raft.BootstrapCluster(raftConfig, b.logStore, b.stableStore, b.snapStore, b.raftTransport, *bootstrapConfig); err != nil {
			return err
		}
		// If we are the only node we should start as the leader.
		if len(bootstrapConfig.Servers) == 1 {
			raftConfig.StartAsLeader = true
		}
	}

	// Setup the Raft store.
	b.fsm.SetNoopRestore(true)
	raftObj, err := raft.NewRaft(raftConfig, b.fsm, b.logStore, b.stableStore, b.snapStore, b.raftTransport)
	b.fsm.SetNoopRestore(false)
	if err != nil {
		return err
	}
	b.raft = raftObj
	b.raftNotifyCh = raftNotifyCh

	if b.raftLayer != nil {
		// Add Handler to the cluster.
		clusterListener.AddHandler(consts.RaftStorageALPN, b.raftLayer)

		// Add Client to the cluster.
		clusterListener.AddClient(consts.RaftStorageALPN, b.raftLayer)
	}

	return nil
}

// TeardownCluster shuts down the raft cluster
func (b *RaftBackend) TeardownCluster(clusterListener physicalstd.ClusterHook) error {
	clusterListener.StopHandler(consts.RaftStorageALPN)
	clusterListener.RemoveClient(consts.RaftStorageALPN)
	b.l.Lock()
	future := b.raft.Shutdown()
	b.raft = nil
	b.l.Unlock()

	return future.Error()
}

// RemovePeer removes the given peer ID from the raft cluster. If the node is
// ourselves we will give up leadership.
func (b *RaftBackend) RemovePeer(ctx context.Context, peerID string) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage is not initialized")
	}

	future := b.raft.RemoveServer(raft.ServerID(peerID), 0, 0)

	return future.Error()
}

func (b *RaftBackend) GetConfiguration(ctx context.Context) (*RaftConfigurationResponse, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return nil, errors.New("raft storage is not initialized")
	}

	future := b.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		return nil, err
	}

	config := &RaftConfigurationResponse{
		Index: future.Index(),
	}

	for _, server := range future.Configuration().Servers {
		entry := &RaftServer{
			NodeID:          string(server.ID),
			Address:         string(server.Address),
			Leader:          server.Address == b.raft.Leader(),
			Voter:           server.Suffrage == raft.Voter,
			ProtocolVersion: string(raft.ProtocolVersionMax),
		}
		config.Servers = append(config.Servers, entry)
	}

	return config, nil
}

// AddPeer adds a new server to the raft cluster
func (b *RaftBackend) AddPeer(ctx context.Context, peerID, clusterAddr string) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage is not initialized")
	}

	future := b.raft.AddVoter(raft.ServerID(peerID), raft.ServerAddress(clusterAddr), 0, 0)

	return future.Error()
}

// Peers returns all the servers present in the raft cluster
func (b *RaftBackend) Peers(ctx context.Context) ([]Peer, error) {
	if b.raft == nil {
		return nil, errors.New("raft storage backend is not initialized")
	}

	future := b.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		return nil, err
	}

	ret := make([]Peer, len(future.Configuration().Servers))
	for i, s := range future.Configuration().Servers {
		ret[i] = Peer{
			ID:      string(s.ID),
			Address: string(s.Address),
		}
	}

	return ret, nil
}

// Delete inserts an entry in the log to delete the given path
func (b *RaftBackend) Delete(ctx context.Context, path string) error {
	command := &LogData{
		Operations: []*LogOperation{
			&LogOperation{
				OpType: deleteOp,
				Key:    path,
			},
		},
	}

	return b.applyLog(ctx, command)
}

// Get returns the value corresponding to the given path from the fsm
func (b *RaftBackend) Get(ctx context.Context, path string) (*physical.Entry, error) {
	if b.fsm == nil {
		return nil, errors.New("raft: fsm not configured")
	}

	return b.fsm.Get(ctx, path)
}

// Put inserts an entry in the log for the put operation
func (b *RaftBackend) Put(ctx context.Context, entry *physical.Entry) error {
	command := &LogData{
		Operations: []*LogOperation{
			&LogOperation{
				OpType: putOp,
				Key:    entry.Key,
				Value:  entry.Value,
			},
		},
	}

	return b.applyLog(ctx, command)
}

// List enumerates all the items under the prefix from the fsm
func (b *RaftBackend) List(ctx context.Context, prefix string) ([]string, error) {
	if b.fsm == nil {
		return nil, errors.New("raft: fsm not configured")
	}

	return b.fsm.List(ctx, prefix)
}

// Transaction applies all the given operations into a single log and
// applies it.
func (b *RaftBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	command := &LogData{
		Operations: make([]*LogOperation, len(txns)),
	}
	for i, txn := range txns {
		op := &LogOperation{}
		switch txn.Operation {
		case physical.PutOperation:
			op.OpType = putOp
			op.Key = txn.Entry.Key
			op.Value = txn.Entry.Value
		case physical.DeleteOperation:
			op.OpType = deleteOp
			op.Key = txn.Entry.Key
		default:
			return fmt.Errorf("%q is not a supported transaction operation", txn.Operation)
		}

		command.Operations[i] = op
	}

	return b.applyLog(ctx, command)
}

// applyLog will take a given log command and apply it to the raft log. applyLog
// doesn't return until the log has been applied to a quorum of servers and is
// persisted to the local FSM.
func (b *RaftBackend) applyLog(ctx context.Context, command *LogData) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage backend is not initialized")
	}

	commandBytes, err := proto.Marshal(command)
	if err != nil {
		return err
	}

	applyFuture := b.raft.Apply(commandBytes, 0)
	err = applyFuture.Error()
	if err != nil {
		return err
	}

	if resp, ok := applyFuture.Response().(*FSMApplyResponse); !ok || !resp.Success {
		return errors.New("could not apply data")
	}

	return nil
}

// HAEnabled is the implemention of the HABackend interface
func (b *RaftBackend) HAEnabled() bool { return true }

// HAEnabled is the implemention of the HABackend interface
func (b *RaftBackend) LockWith(key, value string) (physical.Lock, error) {
	return &RaftLock{
		key:   key,
		value: []byte(value),
		b:     b,
	}, nil
}

// RaftLock implements the physical Lock interface and enables HA for this
// backend. The Lock uses the raftNotifyCh for receiving leadership edge
// triggers. Vault's active duty matches raft's leadership.
type RaftLock struct {
	key   string
	value []byte

	b *RaftBackend
}

// monitorLeadership waits until we receive an update on the raftNotifyCh and
// closes the leaderLost channel.
func (l *RaftLock) monitorLeadership(stopCh <-chan struct{}) <-chan struct{} {
	leaderLost := make(chan struct{})
	go func() {
		select {
		case <-l.b.raftNotifyCh:
			close(leaderLost)
		case <-stopCh:
		}
	}()
	return leaderLost
}

// Lock blocks until we become leader or are shutdown. It returns a channel that
// is closed when we detect a loss of leadership.
func (l *RaftLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	for {
		select {
		case isLeader := <-l.b.raftNotifyCh:
			if isLeader {
				// We are leader, set the key
				err := l.b.applyLog(context.Background(), &LogData{
					Operations: []*LogOperation{
						&LogOperation{
							OpType: putOp,
							Key:    l.key,
							Value:  l.value,
						},
					},
				})
				if err != nil {
					return nil, err
				}

				return l.monitorLeadership(stopCh), nil
			}
		case <-stopCh:
			return nil, nil
		}
	}

	return nil, nil
}

// Unlock gives up leadership.
func (l *RaftLock) Unlock() error {
	return l.b.raft.LeadershipTransfer().Error()
}

// Value reads the value of the lock. This informs us who is currently leader.
func (l *RaftLock) Value() (bool, string, error) {
	e, err := l.b.Get(context.Background(), l.key)
	if err != nil {
		return false, "", err
	}
	if e == nil {
		return false, "", nil
	}

	value := string(e.Value)
	// TODO: how to tell if held?
	return true, value, nil
}
