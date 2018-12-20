package raft

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/hashicorp/consul/lib"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/vault/physical/raft/logstore"

	"github.com/hashicorp/vault/physical"
)

// Verify RaftBackend satisfies the correct interfaces
var _ physical.Backend = (*RaftBackend)(nil)
var _ physical.Transactional = (*RaftBackend)(nil)
var _ physical.Unsealable = (*RaftBackend)(nil)

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

	fsm          *FSM
	raft         *raft.Raft
	raftNotifyCh chan bool
}

// NewRaftBackend constructs a RaftBackend using the given directory
func NewRaftBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// Create the FSM.
	var err error
	fsm, err := NewFSM(conf, logger)
	if err != nil {
		return nil, err
	}

	return &RaftBackend{
		logger: logger,
		fsm:    fsm,
		conf:   conf,
	}, nil
}

func (b *RaftBackend) Unseal(ctx context.Context, encryptor physical.EncryptorHook) error {
	b.l.Lock()
	defer b.l.Unlock()

	// We are already unsealed
	if b.raft != nil {
		return nil
	}

	path, ok := b.conf["path"]
	if !ok {
		return fmt.Errorf("'path' must be set")
	}

	raftConfig := raft.DefaultConfig()
	raftConfig.SnapshotThreshold = 1

	/*var serverAddressProvider raft.ServerAddressProvider = nil
	if s.config.RaftConfig.ProtocolVersion >= 3 { //ServerAddressProvider needs server ids to work correctly, which is only supported in protocol version 3 or higher
		serverAddressProvider = s.serverLookup
	}*/

	// Create a transport layer.
	trans, err := raft.NewTCPTransport("127.0.0.1:8202", nil, 3, 10*time.Second, nil)
	if err != nil {
		return err
	}

	/*	transConfig := &raft.NetworkTransportConfig{
			Stream:  s.raftLayer,
			MaxPool: 3,
			Timeout: 10 * time.Second,
			//	ServerAddressProvider: serverAddressProvider,
		}

		trans := raft.NewNetworkTransportWithConfig(transConfig)*/
	//	s.raftTransport = trans

	// Make sure we set the LogOutput.
	//	s.config.RaftConfig.LogOutput = s.config.LogOutput
	//	s.config.RaftConfig.Logger = s.logger

	// Versions of the Raft protocol below 3 require the LocalID to match the network
	// address of the transport.
	raftConfig.LocalID = raft.ServerID(trans.LocalAddr())

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
		if err := lib.EnsurePath(path, true); err != nil {
			return err
		}

		// Create the backend raft store for logs and stable storage.
		store, err := raftboltdb.NewBoltStore(filepath.Join(path, "raft.db"), encryptor)
		if err != nil {
			return err
		}
		stable = store

		// Wrap the store in a LogCache to improve performance.
		cacheStore, err := raft.NewLogCache(raftLogCacheSize, store)
		if err != nil {
			return err
		}
		log = cacheStore

		// Create the snapshot store.
		snapshots, err := raft.NewFileSnapshotStore(path, snapshotsRetained, nil)
		if err != nil {
			return err
		}
		snap = snapshots
	}

	// If we are in bootstrap or dev mode and the state is clean then we can
	// bootstrap now.
	//		if s.config.Bootstrap || s.config.DevMode {
	hasState, err := raft.HasExistingState(log, stable, snap)
	if err != nil {
		return err
	}
	if !hasState {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				raft.Server{
					ID:      raftConfig.LocalID,
					Address: trans.LocalAddr(),
				},
			},
		}
		if err := raft.BootstrapCluster(raftConfig, log, stable, snap, trans, configuration); err != nil {
			return err
		}
	}
	//		}

	// Set up a channel for reliable leader notifications.
	raftNotifyCh := make(chan bool, 1)
	raftConfig.NotifyCh = raftNotifyCh

	// Setup the Raft store.
	raftObj, err := raft.NewRaft(raftConfig, b.fsm, log, stable, snap, trans)
	if err != nil {
		return err
	}
	b.raft = raftObj
	b.raftNotifyCh = raftNotifyCh

	return nil
}

func (b *RaftBackend) Seal(ctx context.Context) error {
	b.l.Lock()
	future := b.raft.Shutdown()
	b.raft = nil
	b.l.Unlock()

	return future.Error()
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
