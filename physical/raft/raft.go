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

type RaftBackend struct {
	logger log.Logger
	conf   map[string]string
	l      sync.RWMutex

	fsm             *FSM
	raft            *raft.Raft
	raftNotifyCh    chan bool
	raftLayer       *raftLayer
	raftTransport   raft.Transport
	snapStore       raft.SnapshotStore
	logStore        raft.LogStore
	stableStore     raft.StableStore
	bootstrapConfig *raft.Configuration

	dataDir string
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
		return nil, err
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
		snapshots, err := raft.NewFileSnapshotStore(path, snapshotsRetained, nil)
		if err != nil {
			return nil, err
		}
		snap = snapshots
	}

	var localID string
	{
		// Determine the local node ID
		localID = conf["node_id"]

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

type Peer struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

func (b *RaftBackend) NodeID() string {
	return b.localID
}

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

	b.bootstrapConfig = raftConfig
	return nil
}

func (b *RaftBackend) SetupCluster(ctx context.Context, networkConfig *physicalstd.NetworkConfig, clusterListener physicalstd.ClusterHook) error {
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

	raftConfig := raft.DefaultConfig()
	// Make sure we set the LogOutput.
	//	s.config.RaftConfig.LogOutput = s.config.LogOutput
	//raftConfig.Logger = logger

	switch {
	case networkConfig == nil:
		_, b.raftTransport = raft.NewInmemTransport(raft.ServerAddress(clusterListener.Addr().String()))
	default:
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
		b.bootstrapConfig = nil

		if err := raft.BootstrapCluster(raftConfig, b.logStore, b.stableStore, b.snapStore, b.raftTransport, *bootstrapConfig); err != nil {
			return err
		}
		if len(bootstrapConfig.Servers) == 1 {
			raftConfig.StartAsLeader = true
		}
	}

	// Setup the Raft store.
	raftObj, err := raft.NewRaft(raftConfig, b.fsm, b.logStore, b.stableStore, b.snapStore, b.raftTransport)
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

func (b *RaftBackend) TeardownCluster(clusterListener physicalstd.ClusterHook) error {
	clusterListener.StopHandler(consts.RaftStorageALPN)
	clusterListener.RemoveClient(consts.RaftStorageALPN)
	b.l.Lock()
	future := b.raft.Shutdown()
	b.raft = nil
	b.l.Unlock()

	return future.Error()
}

func (b *RaftBackend) AddPeer(ctx context.Context, peerID, clusterAddr string) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage backend is sealed")
	}

	future := b.raft.AddVoter(raft.ServerID(peerID), raft.ServerAddress(clusterAddr), 0, 0)

	return future.Error()
}

func (b *RaftBackend) Peers(ctx context.Context) ([]Peer, error) {
	if b.raft == nil {
		return nil, errors.New("raft storage backend is sealed")
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

func (b *RaftBackend) Get(ctx context.Context, path string) (*physical.Entry, error) {
	if b.fsm == nil {
		return nil, errors.New("raft: fsm not configured")
	}

	return b.fsm.Get(ctx, path)
}

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

func (b *RaftBackend) List(ctx context.Context, prefix string) ([]string, error) {
	if b.fsm == nil {
		return nil, errors.New("raft: fsm not configured")
	}

	return b.fsm.List(ctx, prefix)
}

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

func (b *RaftBackend) applyLog(ctx context.Context, command *LogData) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage backend is sealed")
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

	if !applyFuture.Response().(*FSMApplyResponse).Success {
		return errors.New("could not apply data")
	}

	return nil
}

func (b *RaftBackend) HAEnabled() bool { return true }
func (b *RaftBackend) LockWith(key, value string) (physical.Lock, error) {
	return &RaftLock{
		key:   key,
		value: []byte(value),
		b:     b,
	}, nil
}

type RaftLock struct {
	key   string
	value []byte

	b *RaftBackend
}

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

func (l *RaftLock) Unlock() error {
	// TODO: how do you stepdown a node?
	return nil
}

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
