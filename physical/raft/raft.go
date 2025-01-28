// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/golang/protobuf/proto"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-raftchunking"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/permitpool"
	"github.com/hashicorp/go-secure-stdlib/tlsutil"
	"github.com/hashicorp/raft"
	autopilot "github.com/hashicorp/raft-autopilot"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	snapshot "github.com/hashicorp/raft-snapshot"
	raftwal "github.com/hashicorp/raft-wal"
	walmetrics "github.com/hashicorp/raft-wal/metrics"
	"github.com/hashicorp/raft-wal/verifier"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/hashicorp/vault/version"
	bolt "go.etcd.io/bbolt"
)

const (
	// EnvVaultRaftNodeID is used to fetch the Raft node ID from the environment.
	EnvVaultRaftNodeID = "VAULT_RAFT_NODE_ID"

	// EnvVaultRaftPath is used to fetch the path where Raft data is stored from the environment.
	EnvVaultRaftPath = "VAULT_RAFT_PATH"

	// EnvVaultRaftNonVoter is used to override the non_voter config option, telling Vault to join as a non-voter (i.e. read replica).
	EnvVaultRaftNonVoter  = "VAULT_RAFT_RETRY_JOIN_AS_NON_VOTER"
	raftNonVoterConfigKey = "retry_join_as_non_voter"

	// EnvVaultRaftMaxBatchEntries is used to override the default maxBatchEntries
	// limit.
	EnvVaultRaftMaxBatchEntries = "VAULT_RAFT_MAX_BATCH_ENTRIES"

	// EnvVaultRaftMaxBatchSizeBytes is used to override the default maxBatchSize
	// limit.
	EnvVaultRaftMaxBatchSizeBytes = "VAULT_RAFT_MAX_BATCH_SIZE_BYTES"

	// defaultMaxBatchEntries is the default maxBatchEntries limit. This was
	// derived from performance testing. It is effectively high enough never to be
	// a real limit for realistic Vault operation sizes and the size limit
	// provides the actual limit since that amount of data stored is more relevant
	// that the specific number of operations.
	defaultMaxBatchEntries = 4096

	// defaultMaxBatchSize is the default maxBatchSize limit. This was derived
	// from performance testing.
	defaultMaxBatchSize = 128 * 1024
)

var (
	getMmapFlags     = func(string) int { return 0 }
	usingMapPopulate = func(int) bool { return false }
)

// Verify RaftBackend satisfies the correct interfaces
var (
	_ physical.Backend                   = (*RaftBackend)(nil)
	_ physical.Transactional             = (*RaftBackend)(nil)
	_ physical.TransactionalLimits       = (*RaftBackend)(nil)
	_ physical.HABackend                 = (*RaftBackend)(nil)
	_ physical.MountTableLimitingBackend = (*RaftBackend)(nil)
	_ physical.RemovableNodeHABackend    = (*RaftBackend)(nil)
	_ physical.Lock                      = (*RaftLock)(nil)
)

var (
	// raftLogCacheSize is the maximum number of logs to cache in-memory.
	// This is used to reduce disk I/O for the recently committed entries.
	raftLogCacheSize = 512

	raftState                          = "raft/"
	raftWalDir                         = "wal/"
	peersFileName                      = "peers.json"
	restoreOpDelayDuration             = 5 * time.Second
	defaultMaxEntrySize                = uint64(2 * raftchunking.ChunkSize)
	defaultRaftLogVerificationInterval = 60 * time.Second
	minimumRaftLogVerificationInterval = 10 * time.Second

	GetInTxnDisabledError = errors.New("get operations inside transactions are disabled in raft backend")
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

	// raftInitCh is used to block during HA lock acquisition if raft
	// has not been initialized yet, which can occur if raft is being
	// used for HA-only.
	raftInitCh chan struct{}

	// raftNotifyCh is used to receive updates about leadership changes
	// regarding this node.
	raftNotifyCh chan bool

	// streamLayer is the network layer used to connect the nodes in the raft
	// cluster.
	streamLayer *raftLayer

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

	// closers is a list of managed resource (such as stores above or wrapper
	// layers around them). That should have Close called on them when the backend
	// is closed. We need to take care that each distinct object is closed only
	// once which might involve knowing how wrappers to stores work. For example
	// raft wal verifier wraps LogStore and is an io.Closer but it also closes the
	// underlying LogStore so if we add it here we shouldn't also add the actual
	// LogStore or StableStore if it's the same underlying instance. We could use
	// a map[io.Closer]bool to prevent double registrations, but that doesn't
	// solve the problem of "knowing" whether or not calling Close on some wrapper
	// also calls "Close" on it's underlying.
	closers []io.Closer

	// dataDir is the location on the local filesystem that raft and FSM data
	// will be stored.
	dataDir string

	// localID is the ID for this node. This can either be configured in the
	// config file, via a file on disk, or is otherwise randomly generated.
	localID string

	// serverAddressProvider is used to map server IDs to addresses.
	serverAddressProvider raft.ServerAddressProvider

	// permitPool is used to limit the number of concurrent storage calls.
	permitPool *permitpool.Pool

	// maxEntrySize imposes a size limit (in bytes) on a raft entry (put or transaction).
	// It is suggested to use a value of 2x the Raft chunking size for optimal
	// performance.
	maxEntrySize uint64

	// maxMountAndNamespaceEntrySize imposes a size limit (in bytes) on a raft
	// entry (put or transaction) for paths related to mount table and namespace
	// metadata. The Raft storage doesn't "know" about these paths but Vault can
	// call RegisterMountTablePath to inform it so that it can apply this
	// alternate limit if one is configured.
	maxMountAndNamespaceEntrySize uint64

	// maxBatchEntries is the number of operation entries in each batch. It is set
	// by default to a value we've tested to work well but may be overridden by
	// Environment variable VAULT_RAFT_MAX_BATCH_ENTRIES.
	maxBatchEntries int

	// maxBatchSize is the maximum combined key and value size of operation
	// entries in each batch. It is set by default to a value we've tested to work
	// well but may be overridden by Environment variable
	// VAULT_RAFT_MAX_BATCH_SIZE_BYTES.
	maxBatchSize int

	// autopilot is the instance of raft-autopilot library implementation of the
	// autopilot features. This will be instantiated in both leader and followers.
	// However, only active node will have a "running" autopilot.
	autopilot *autopilot.Autopilot

	// autopilotConfig represents the configuration required to instantiate autopilot.
	autopilotConfig *AutopilotConfig

	// followerStates represents the information about all the peers of the raft
	// leader. This is used to track some state of the peers and as well as used
	// to see if the peers are "alive" using the heartbeat received from them.
	followerStates *FollowerStates

	// followerHeartbeatTicker is used to compute dead servers using follower
	// state heartbeats.
	followerHeartbeatTicker *time.Ticker

	// disableAutopilot if set will not put autopilot implementation to use. The
	// fallback will be to interact with the raft instance directly. This can only
	// be set during startup via the environment variable
	// VAULT_RAFT_AUTOPILOT_DISABLE during startup and can't be updated once the
	// node is up and running.
	disableAutopilot bool

	// autopilotReconcileInterval is how long between rounds of performing promotions, demotions
	// and leadership transfers.
	autopilotReconcileInterval time.Duration

	// autopilotUpdateInterval is the time between the periodic state updates. These periodic
	// state updates take in known servers from the delegate, request Raft stats be
	// fetched and pull in other inputs such as the Raft configuration to create
	// an updated view of the Autopilot State.
	autopilotUpdateInterval time.Duration

	// upgradeVersion is used to override the Vault SDK version when performing an autopilot automated upgrade.
	upgradeVersion string

	// redundancyZone specifies a redundancy zone for autopilot.
	redundancyZone string

	// nonVoter specifies whether the node should join the cluster as a non-voter. Non-voters get
	// replicated to and can serve reads, but do not take part in leader elections.
	nonVoter bool

	effectiveSDKVersion string
	failGetInTxn        *uint32

	// raftLogVerifierEnabled and raftLogVerificationInterval control enabling the raft log verifier and how often
	// it writes checkpoints.
	raftLogVerifierEnabled      bool
	raftLogVerificationInterval time.Duration

	// specialPathLimits is a map of special paths to their configured entrySize
	// limits.
	specialPathLimits map[string]uint64

	removed              *atomic.Bool
	removedCallback      func()
	removedServerCleanup func(context.Context, string) (bool, error)
}

func (b *RaftBackend) IsNodeRemoved(ctx context.Context, nodeID string) (bool, error) {
	conf, err := b.GetConfiguration(ctx)
	if err != nil {
		return false, err
	}
	for _, srv := range conf.Servers {
		if srv.NodeID == nodeID {
			return false, nil
		}
	}
	return true, nil
}

func (b *RaftBackend) IsRemoved() bool {
	return b.removed.Load()
}

var removedKey = []byte("removed")

func (b *RaftBackend) RemoveSelf() error {
	b.removed.Store(true)
	return b.stableStore.SetUint64(removedKey, 1)
}

func (b *RaftBackend) SetRemovedServerCleanupFunc(f func(context.Context, string) (bool, error)) {
	b.l.Lock()
	b.removedServerCleanup = f
	b.l.Unlock()
}

func (b *RaftBackend) RemovedServerCleanup(ctx context.Context, nodeID string) (bool, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.removedServerCleanup != nil {
		return b.removedServerCleanup(ctx, nodeID)
	}

	return false, nil
}

// LeaderJoinInfo contains information required by a node to join itself as a
// follower to an existing raft cluster
type LeaderJoinInfo struct {
	// AutoJoin defines any cloud auto-join metadata. If supplied, Vault will
	// attempt to automatically discover peers in addition to what can be provided
	// via 'leader_api_addr'.
	AutoJoin string `json:"auto_join"`

	// AutoJoinScheme defines the optional URI protocol scheme for addresses
	// discovered via auto-join.
	AutoJoinScheme string `json:"auto_join_scheme"`

	// AutoJoinPort defines the optional port used for addressed discovered via
	// auto-join.
	AutoJoinPort uint `json:"auto_join_port"`

	// LeaderAPIAddr is the address of the leader node to connect to
	LeaderAPIAddr string `json:"leader_api_addr"`

	// LeaderCACert is the CA cert of the leader node
	LeaderCACert string `json:"leader_ca_cert"`

	// LeaderClientCert is the client certificate for the follower node to
	// establish client authentication during TLS
	LeaderClientCert string `json:"leader_client_cert"`

	// LeaderClientKey is the client key for the follower node to establish
	// client authentication during TLS.
	LeaderClientKey string `json:"leader_client_key"`

	// LeaderCACertFile is the path on disk to the CA cert file of the
	// leader node. This should only be provided via Vault's configuration file.
	LeaderCACertFile string `json:"leader_ca_cert_file"`

	// LeaderClientCertFile is the path on disk to the client certificate file
	// for the follower node to establish client authentication during TLS. This
	// should only be provided via Vault's configuration file.
	LeaderClientCertFile string `json:"leader_client_cert_file"`

	// LeaderClientKeyFile is the path on disk to the client key file for the
	// follower node to establish client authentication during TLS. This should
	// only be provided via Vault's configuration file.
	LeaderClientKeyFile string `json:"leader_client_key_file"`

	// LeaderTLSServerName is the optional ServerName to expect in the leader's
	// certificate, instead of the host/IP we're actually connecting to.
	LeaderTLSServerName string `json:"leader_tls_servername"`

	// Retry indicates if the join process should automatically be retried
	Retry bool `json:"-"`

	// TLSConfig for the API client to use when communicating with the leader node
	TLSConfig *tls.Config `json:"-"`
}

// JoinConfig returns a list of information about possible leader nodes that
// this node can join as a follower
func (b *RaftBackend) JoinConfig() ([]*LeaderJoinInfo, error) {
	config := b.conf["retry_join"]
	if config == "" {
		return nil, nil
	}

	var leaderInfos []*LeaderJoinInfo
	err := jsonutil.DecodeJSON([]byte(config), &leaderInfos)
	if err != nil {
		return nil, fmt.Errorf("failed to decode retry_join config: %w", err)
	}

	if len(leaderInfos) == 0 {
		return nil, errors.New("invalid retry_join config")
	}

	for i, info := range leaderInfos {
		if len(info.AutoJoin) != 0 && len(info.LeaderAPIAddr) != 0 {
			return nil, errors.New("cannot provide both a leader_api_addr and auto_join")
		}

		if info.AutoJoinScheme != "" && (info.AutoJoinScheme != "http" && info.AutoJoinScheme != "https") {
			return nil, fmt.Errorf("invalid scheme %q; must either be http or https", info.AutoJoinScheme)
		}

		info.Retry = true
		info.TLSConfig, err = parseTLSInfo(info)
		if err != nil {
			return nil, fmt.Errorf("failed to create tls config to communicate with leader node (retry_join index: %d): %w", i, err)
		}
	}

	return leaderInfos, nil
}

// parseTLSInfo is a helper for parses the TLS information, preferring file
// paths over raw certificate content.
func parseTLSInfo(leaderInfo *LeaderJoinInfo) (*tls.Config, error) {
	var tlsConfig *tls.Config
	var err error
	if len(leaderInfo.LeaderCACertFile) != 0 || len(leaderInfo.LeaderClientCertFile) != 0 || len(leaderInfo.LeaderClientKeyFile) != 0 {
		tlsConfig, err = tlsutil.LoadClientTLSConfig(leaderInfo.LeaderCACertFile, leaderInfo.LeaderClientCertFile, leaderInfo.LeaderClientKeyFile)
		if err != nil {
			return nil, err
		}
	} else if len(leaderInfo.LeaderCACert) != 0 || len(leaderInfo.LeaderClientCert) != 0 || len(leaderInfo.LeaderClientKey) != 0 {
		tlsConfig, err = tlsutil.ClientTLSConfig([]byte(leaderInfo.LeaderCACert), []byte(leaderInfo.LeaderClientCert), []byte(leaderInfo.LeaderClientKey))
		if err != nil {
			return nil, err
		}
	}
	if tlsConfig != nil {
		tlsConfig.ServerName = leaderInfo.LeaderTLSServerName
	}

	return tlsConfig, nil
}

// EnsurePath is used to make sure a path exists
func EnsurePath(path string, dir bool) error {
	if !dir {
		path = filepath.Dir(path)
	}
	return os.MkdirAll(path, 0o700)
}

func NewClusterAddrBridge() *ClusterAddrBridge {
	return &ClusterAddrBridge{
		clusterAddressByNodeID: make(map[string]string),
	}
}

type ClusterAddrBridge struct {
	l                      sync.RWMutex
	clusterAddressByNodeID map[string]string
}

func (c *ClusterAddrBridge) UpdateClusterAddr(nodeId string, clusterAddr string) {
	c.l.Lock()
	defer c.l.Unlock()
	cu, _ := url.Parse(clusterAddr)
	c.clusterAddressByNodeID[nodeId] = cu.Host
}

func (c *ClusterAddrBridge) ServerAddr(id raft.ServerID) (raft.ServerAddress, error) {
	c.l.RLock()
	defer c.l.RUnlock()
	if addr, ok := c.clusterAddressByNodeID[string(id)]; ok {
		return raft.ServerAddress(addr), nil
	}
	return "", fmt.Errorf("could not find cluster addr for id=%s", id)
}

func batchLimitsFromEnv(logger log.Logger) (int, int) {
	maxBatchEntries := defaultMaxBatchEntries
	if envVal := os.Getenv(EnvVaultRaftMaxBatchEntries); envVal != "" {
		if i, err := strconv.Atoi(envVal); err == nil && i > 0 {
			maxBatchEntries = i
		} else {
			logger.Warn("failed to parse VAULT_RAFT_MAX_BATCH_ENTRIES as an integer > 0. Using default value.",
				"env_val", envVal, "default_used", maxBatchEntries)
		}
	}

	maxBatchSize := defaultMaxBatchSize
	if envVal := os.Getenv(EnvVaultRaftMaxBatchSizeBytes); envVal != "" {
		if i, err := strconv.Atoi(envVal); err == nil && i > 0 {
			maxBatchSize = i
		} else {
			logger.Warn("failed to parse VAULT_RAFT_MAX_BATCH_SIZE_BYTES as an integer > 0. Using default value.",
				"env_val", envVal, "default_used", maxBatchSize)
		}
	}

	return maxBatchEntries, maxBatchSize
}

// NewRaftBackend constructs a RaftBackend using the given directory
func NewRaftBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// parse the incoming map into a proper config struct
	backendConfig, err := parseRaftBackendConfig(conf, logger)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	// Create the FSM.
	fsm, err := NewFSM(backendConfig.Path, backendConfig.NodeId, logger.Named("fsm"))
	if err != nil {
		return nil, fmt.Errorf("failed to create fsm: %v", err)
	}

	if backendConfig.ApplyDelay > 0 {
		fsm.applyCallback = func() {
			time.Sleep(backendConfig.ApplyDelay)
		}
	}

	// Create the log store.
	// Build an all in-memory setup for dev mode, otherwise prepare a full
	// disk-based setup.
	var logStore raft.LogStore
	var stableStore raft.StableStore
	var snapStore raft.SnapshotStore
	var closers []io.Closer

	var devMode bool
	if devMode {
		store := raft.NewInmemStore()
		stableStore = store
		logStore = store
		snapStore = raft.NewInmemSnapshotStore()
	} else {
		// Create the base raft path.
		raftBasePath := filepath.Join(backendConfig.Path, raftState)
		if err := EnsurePath(raftBasePath, true); err != nil {
			return nil, err
		}
		dbPath := filepath.Join(raftBasePath, "raft.db")

		// If the existing raft db exists from a previous use of BoltDB, warn about this and continue to use BoltDB
		raftDbExists, err := fileExists(dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to check if raft.db already exists: %w", err)
		}

		if backendConfig.RaftWal && raftDbExists {
			logger.Warn("raft is configured to use raft-wal for storage but existing raft.db detected. raft-wal config will be ignored.")
			backendConfig.RaftWal = false
		}

		if backendConfig.RaftWal {
			raftWalPath := filepath.Join(raftBasePath, raftWalDir)
			if err := EnsurePath(raftWalPath, true); err != nil {
				return nil, err
			}

			mc := walmetrics.NewGoMetricsCollector([]string{"raft", "wal"}, nil, nil)
			wal, err := raftwal.Open(raftWalPath, raftwal.WithMetricsCollector(mc))
			if err != nil {
				return nil, fmt.Errorf("fail to open write-ahead-log: %w", err)
			}
			// We need to Close the store but don't register it in closers yet because
			// if we are going to wrap it with a verifier we need to close through
			// that instead.

			stableStore = wal
			logStore = wal
		} else {
			// use the traditional BoltDB setup
			opts := boltOptions(dbPath)
			raftOptions := raftboltdb.Options{
				Path:                    dbPath,
				BoltOptions:             opts,
				MsgpackUseNewTimeFormat: true,
			}

			if runtime.GOOS == "linux" && raftDbExists && !usingMapPopulate(opts.MmapFlags) {
				logger.Warn("the MAP_POPULATE mmap flag has not been set before opening the log store database. This may be due to the database file being larger than the available memory on the system, or due to the VAULT_RAFT_DISABLE_MAP_POPULATE environment variable being set. As a result, Vault may be slower to start up.")
			}

			store, err := raftboltdb.New(raftOptions)
			if err != nil {
				return nil, err
			}
			// We need to Close the store but don't register it in closers yet because
			// if we are going to wrap it with a verifier we need to close through
			// that instead.

			stableStore = store
			logStore = store
		}

		// Create the snapshot store.
		snapshots, err := NewBoltSnapshotStore(raftBasePath, logger.Named("snapshot"), fsm)
		if err != nil {
			return nil, err
		}
		snapStore = snapshots
	}

	// Hook up the verifier if it's enabled
	if backendConfig.RaftLogVerifierEnabled {
		mc := walmetrics.NewGoMetricsCollector([]string{"raft", "logstore", "verifier"}, nil, nil)
		reportFn := makeLogVerifyReportFn(logger.Named("raft.logstore.verifier"))
		v := verifier.NewLogStore(logStore, isLogVerifyCheckpoint, reportFn, mc)
		logStore = v
	}

	// Register the logStore as a closer whether or not it's wrapped in a verifier
	// (which is a closer). We do this before the LogCache since that is _not_ an
	// io.Closer.
	if closer, ok := logStore.(io.Closer); ok {
		closers = append(closers, closer)
	}
	// Note that we DON'T register the stableStore as a closer because right now
	// we always use the same underlying object as the logStore and we don't want
	// to call close on it twice. If we ever support different stable store and
	// log store then this logic will get even more complex! We don't register
	// snapStore because none of our snapshot stores are io.Closers.

	// Close the FSM
	closers = append(closers, fsm)

	// Wrap the store in a LogCache to improve performance.
	cacheStore, err := raft.NewLogCache(raftLogCacheSize, logStore)
	if err != nil {
		return nil, err
	}
	logStore = cacheStore

	if backendConfig.SnapshotDelay > 0 {
		snapStore = newSnapshotStoreDelay(snapStore, backendConfig.SnapshotDelay, logger)
	}

	isRemoved := new(atomic.Bool)
	removedVal, err := stableStore.GetUint64(removedKey)
	if err != nil {
		logger.Error("error checking if this node is removed. continuing under the assumption that it's not", "error", err)
	}
	if removedVal == 1 {
		isRemoved.Store(true)
	}
	return &RaftBackend{
		logger:                        logger,
		fsm:                           fsm,
		raftInitCh:                    make(chan struct{}),
		conf:                          conf,
		logStore:                      logStore,
		stableStore:                   stableStore,
		snapStore:                     snapStore,
		closers:                       closers,
		dataDir:                       backendConfig.Path,
		localID:                       backendConfig.NodeId,
		permitPool:                    permitpool.New(physical.DefaultParallelOperations),
		maxEntrySize:                  backendConfig.MaxEntrySize,
		maxMountAndNamespaceEntrySize: backendConfig.MaxMountAndNamespaceTableEntrySize,
		maxBatchEntries:               backendConfig.MaxBatchEntries,
		maxBatchSize:                  backendConfig.MaxBatchSize,
		followerHeartbeatTicker:       time.NewTicker(time.Second),
		autopilotReconcileInterval:    backendConfig.AutopilotReconcileInterval,
		autopilotUpdateInterval:       backendConfig.AutopilotUpdateInterval,
		redundancyZone:                backendConfig.AutopilotRedundancyZone,
		nonVoter:                      backendConfig.RaftNonVoter,
		upgradeVersion:                backendConfig.AutopilotUpgradeVersion,
		failGetInTxn:                  new(uint32),
		raftLogVerifierEnabled:        backendConfig.RaftLogVerifierEnabled,
		raftLogVerificationInterval:   backendConfig.RaftLogVerificationInterval,
		effectiveSDKVersion:           version.GetVersion().Version,
		removed:                       isRemoved,
	}, nil
}

// RegisterMountTablePath informs the Backend that the given path represents
// part of the mount tables or related metadata. This allows the backend to
// apply different limits for this entry if configured to do so.
func (b *RaftBackend) RegisterMountTablePath(path string) {
	// We don't need to lock here because this is only called during startup

	if b.maxMountAndNamespaceEntrySize > 0 {
		// Set up the limit for this special path.
		if b.specialPathLimits == nil {
			b.specialPathLimits = make(map[string]uint64)
		}
		b.specialPathLimits[path] = b.maxMountAndNamespaceEntrySize
	}
}

// GetSpecialPathLimits returns any paths registered with special entry size
// limits. It's really only used to make integration testing of the plumbing for
// these paths simpler.
func (b *RaftBackend) GetSpecialPathLimits() map[string]uint64 {
	return b.specialPathLimits
}

type snapshotStoreDelay struct {
	logger  log.Logger
	wrapped raft.SnapshotStore
	delay   time.Duration
}

func (s snapshotStoreDelay) Create(version raft.SnapshotVersion, index, term uint64, configuration raft.Configuration, configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
	s.logger.Trace("delaying before creating snapshot", "delay", s.delay)
	time.Sleep(s.delay)
	return s.wrapped.Create(version, index, term, configuration, configurationIndex, trans)
}

func (s snapshotStoreDelay) List() ([]*raft.SnapshotMeta, error) {
	return s.wrapped.List()
}

func (s snapshotStoreDelay) Open(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	return s.wrapped.Open(id)
}

var _ raft.SnapshotStore = &snapshotStoreDelay{}

func newSnapshotStoreDelay(snap raft.SnapshotStore, delay time.Duration, logger log.Logger) *snapshotStoreDelay {
	return &snapshotStoreDelay{
		logger:  logger,
		wrapped: snap,
		delay:   delay,
	}
}

// Close is used to gracefully close all file resources.  N.B. This method
// should only be called if you are sure the RaftBackend will never be used
// again.
func (b *RaftBackend) Close() error {
	b.l.Lock()
	defer b.l.Unlock()

	for _, cl := range b.closers {
		if err := cl.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (b *RaftBackend) FailGetInTxn(fail bool) {
	var val uint32
	if fail {
		val = 1
	}
	atomic.StoreUint32(b.failGetInTxn, val)
}

func (b *RaftBackend) RedundancyZone() string {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.redundancyZone
}

func (b *RaftBackend) NonVoter() bool {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.nonVoter
}

// UpgradeVersion returns the string that should be used by autopilot during automated upgrades. We return the
// specified upgradeVersion if it's present. If it's not, we fall back to effectiveSDKVersion, which is
// Vault's binary version (though that can be overridden for tests).
func (b *RaftBackend) UpgradeVersion() string {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.upgradeVersion != "" {
		return b.upgradeVersion
	}

	return b.effectiveSDKVersion
}

func (b *RaftBackend) SDKVersion() string {
	b.l.RLock()
	defer b.l.RUnlock()
	return b.effectiveSDKVersion
}

func (b *RaftBackend) verificationInterval() time.Duration {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.raftLogVerificationInterval
}

func (b *RaftBackend) verifierEnabled() bool {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.raftLogVerifierEnabled
}

// DisableUpgradeMigration returns the state of the DisableUpgradeMigration config flag and whether it was set or not
func (b *RaftBackend) DisableUpgradeMigration() (bool, bool) {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.autopilotConfig == nil {
		return false, false
	}

	return b.autopilotConfig.DisableUpgradeMigration, true
}

// StartRaftWalVerifier runs a raft log store verifier in the background, if configured to do so.
// This periodically writes out special raft logs to verify that the log store is not corrupting data.
// This is only safe to run on the raft leader.
func (b *RaftBackend) StartRaftWalVerifier(ctx context.Context) {
	if !b.verifierEnabled() {
		return
	}

	go func() {
		ticker := time.NewTicker(b.verificationInterval())
		defer ticker.Stop()

		logger := b.logger.Named("raft-wal-verifier")

		for {
			select {
			case <-ticker.C:
				err := b.applyVerifierCheckpoint()
				if err != nil {
					logger.Error("error applying verification checkpoint", "error", err)
				}
				logger.Debug("sent verification checkpoint")
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (b *RaftBackend) applyVerifierCheckpoint() error {
	data := make([]byte, 1)
	data[0] = byte(verifierCheckpointOp)

	if err := b.permitPool.Acquire(context.Background()); err != nil {
		return err
	}
	b.l.RLock()

	var err error
	applyFuture := b.raft.Apply(data, 0)
	if e := applyFuture.Error(); e != nil {
		err = e
	}

	b.l.RUnlock()
	b.permitPool.Release()

	return err
}

// isLogVerifyCheckpoint is the verifier.IsCheckpointFn that can decode our raft logs for
// their type.
func isLogVerifyCheckpoint(l *raft.Log) (bool, error) {
	return isRaftLogVerifyCheckpoint(l), nil
}

func makeLogVerifyReportFn(logger log.Logger) verifier.ReportFn {
	return func(r verifier.VerificationReport) {
		if r.SkippedRange != nil {
			logger.Warn("verification skipped range, consider decreasing validation interval if this is frequent",
				"rangeStart", int64(r.SkippedRange.Start),
				"rangeEnd", int64(r.SkippedRange.End),
			)
		}

		l2 := logger.With(
			"rangeStart", int64(r.Range.Start),
			"rangeEnd", int64(r.Range.End),
			"leaderChecksum", fmt.Sprintf("%08x", r.ExpectedSum),
			"elapsed", r.Elapsed,
		)

		if r.Err == nil {
			l2.Info("verification checksum OK",
				"readChecksum", fmt.Sprintf("%08x", r.ReadSum),
			)
			return
		}

		if errors.Is(r.Err, verifier.ErrRangeMismatch) {
			l2.Warn("verification checksum skipped as we don't have all logs in range")
			return
		}

		var csErr verifier.ErrChecksumMismatch
		if errors.As(r.Err, &csErr) {
			if r.WrittenSum > 0 && r.WrittenSum != r.ExpectedSum {
				// The failure occurred before the follower wrote to the log so it
				// must be corrupted in flight from the leader!
				l2.Error("verification checksum FAILED: in-flight corruption",
					"followerWriteChecksum", fmt.Sprintf("%08x", r.WrittenSum),
					"readChecksum", fmt.Sprintf("%08x", r.ReadSum),
				)
			} else {
				l2.Error("verification checksum FAILED: storage corruption",
					"followerWriteChecksum", fmt.Sprintf("%08x", r.WrittenSum),
					"readChecksum", fmt.Sprintf("%08x", r.ReadSum),
				)
			}
			return
		}

		// Some other unknown error occurred
		l2.Error(r.Err.Error())
	}
}

func (b *RaftBackend) CollectMetrics(sink *metricsutil.ClusterMetricSink) {
	var stats map[string]string
	var logStoreStats *bolt.Stats

	b.l.RLock()
	if boltStore, ok := b.stableStore.(*raftboltdb.BoltStore); ok {
		bss := boltStore.Stats()
		logStoreStats = &bss
	}

	if b.raft != nil {
		stats = b.raft.Stats()
	}

	fsmStats := b.fsm.Stats()
	b.l.RUnlock()

	if logStoreStats != nil {
		b.collectMetricsWithStats(*logStoreStats, sink, "logstore")
	}

	b.collectMetricsWithStats(fsmStats, sink, "fsm")
	labels := []metrics.Label{
		{
			Name:  "peer_id",
			Value: b.localID,
		},
	}

	if stats != nil {
		for _, key := range []string{"term", "commit_index", "applied_index", "fsm_pending"} {
			n, err := strconv.ParseUint(stats[key], 10, 64)
			if err == nil {
				sink.SetGaugeWithLabels([]string{"raft_storage", "stats", key}, float32(n), labels)
			}
		}
	}
}

func (b *RaftBackend) collectMetricsWithStats(stats bolt.Stats, sink *metricsutil.ClusterMetricSink, database string) {
	txstats := stats.TxStats
	labels := []metricsutil.Label{{"database", database}}
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "freelist", "free_pages"}, float32(stats.FreePageN), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "freelist", "pending_pages"}, float32(stats.PendingPageN), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "freelist", "allocated_bytes"}, float32(stats.FreeAlloc), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "freelist", "used_bytes"}, float32(stats.FreelistInuse), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "transaction", "started_read_transactions"}, float32(stats.TxN), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "transaction", "currently_open_read_transactions"}, float32(stats.OpenTxN), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "page", "count"}, float32(txstats.GetPageCount()), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "page", "bytes_allocated"}, float32(txstats.GetPageAlloc()), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "cursor", "count"}, float32(txstats.GetCursorCount()), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "node", "count"}, float32(txstats.GetNodeCount()), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "node", "dereferences"}, float32(txstats.GetNodeDeref()), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "rebalance", "count"}, float32(txstats.GetRebalance()), labels)
	sink.AddSampleWithLabels([]string{"raft_storage", "bolt", "rebalance", "time"}, float32(txstats.GetRebalanceTime().Milliseconds()), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "split", "count"}, float32(txstats.GetSplit()), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "spill", "count"}, float32(txstats.GetSpill()), labels)
	sink.AddSampleWithLabels([]string{"raft_storage", "bolt", "spill", "time"}, float32(txstats.GetSpillTime().Milliseconds()), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "write", "count"}, float32(txstats.GetWrite()), labels)
	sink.IncrCounterWithLabels([]string{"raft_storage", "bolt", "write", "time"}, float32(txstats.GetWriteTime().Milliseconds()), labels)
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

// Peer defines the ID and Address for a given member of the raft cluster.
type Peer struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Suffrage int    `json:"suffrage"`
}

// NodeID returns the identifier of the node
func (b *RaftBackend) NodeID() string {
	return b.localID
}

// Initialized tells if raft is running or not
func (b *RaftBackend) Initialized() bool {
	b.l.RLock()
	init := b.raft != nil
	b.l.RUnlock()
	return init
}

// SetTLSKeyring is used to install a new keyring. If the active key has changed
// it will also close any network connections or streams forcing a reconnect
// with the new key.
func (b *RaftBackend) SetTLSKeyring(keyring *TLSKeyring) error {
	b.l.RLock()
	err := b.streamLayer.setTLSKeyring(keyring)
	b.l.RUnlock()

	return err
}

// SetServerAddressProvider sets a the address provider for determining the raft
// node addresses. This is currently only used in tests.
func (b *RaftBackend) SetServerAddressProvider(provider raft.ServerAddressProvider) {
	b.l.Lock()
	b.serverAddressProvider = provider
	b.l.Unlock()
}

// Bootstrap prepares the given peers to be part of the raft cluster
func (b *RaftBackend) Bootstrap(peers []Peer) error {
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
			ID:       raft.ServerID(p.ID),
			Address:  raft.ServerAddress(p.Address),
			Suffrage: raft.ServerSuffrage(p.Suffrage),
		}
	}

	// Store the config for later use
	b.bootstrapConfig = raftConfig
	return nil
}

// SetRestoreCallback sets the callback to be used when a restoreCallbackOp is
// processed through the FSM.
func (b *RaftBackend) SetRestoreCallback(restoreCb restoreCallback) {
	b.fsm.l.Lock()
	b.fsm.restoreCb = restoreCb
	b.fsm.l.Unlock()
}

func (b *RaftBackend) SetRemovedCallback(cb func()) {
	b.l.Lock()
	defer b.l.Unlock()
	b.removedCallback = cb
}

func (b *RaftBackend) applyConfigSettings(config *raft.Config) error {
	config.Logger = b.logger
	multiplierRaw, ok := b.conf["performance_multiplier"]
	multiplier := 5
	if ok {
		var err error
		multiplier, err = strconv.Atoi(multiplierRaw)
		if err != nil {
			return err
		}
	}
	config.ElectionTimeout *= time.Duration(multiplier)
	config.HeartbeatTimeout *= time.Duration(multiplier)
	config.LeaderLeaseTimeout *= time.Duration(multiplier)

	snapThresholdRaw, ok := b.conf["snapshot_threshold"]
	if ok {
		var err error
		snapThreshold, err := strconv.Atoi(snapThresholdRaw)
		if err != nil {
			return err
		}
		config.SnapshotThreshold = uint64(snapThreshold)
	}

	trailingLogsRaw, ok := b.conf["trailing_logs"]
	if ok {
		var err error
		trailingLogs, err := strconv.Atoi(trailingLogsRaw)
		if err != nil {
			return err
		}
		config.TrailingLogs = uint64(trailingLogs)
	}
	snapshotIntervalRaw, ok := b.conf["snapshot_interval"]
	if ok {
		var err error
		snapshotInterval, err := parseutil.ParseDurationSecond(snapshotIntervalRaw)
		if err != nil {
			return err
		}
		config.SnapshotInterval = snapshotInterval
	}

	config.NoSnapshotRestoreOnStart = true
	config.MaxAppendEntries = 64

	// Setting BatchApplyCh allows the raft library to enqueue up to
	// MaxAppendEntries into each raft apply rather than relying on the
	// scheduler.
	config.BatchApplyCh = true

	b.logger.Trace("applying raft config", "inputs", b.conf)
	return nil
}

// SetupOpts are used to pass options to the raft setup function.
type SetupOpts struct {
	// TLSKeyring is the keyring to use for the cluster traffic.
	TLSKeyring *TLSKeyring

	// ClusterListener is the cluster hook used to register the raft handler and
	// client with core's cluster listeners.
	ClusterListener cluster.ClusterHook

	// StartAsLeader is used to specify this node should start as leader and
	// bypass the leader election. This should be used with caution.
	StartAsLeader bool

	// RecoveryModeConfig is the configuration for the raft cluster in recovery
	// mode.
	RecoveryModeConfig *raft.Configuration

	// EffectiveSDKVersion is typically the version string baked into the binary.
	// We pass it in though because it can be overridden in tests or via ENV in
	// core.
	EffectiveSDKVersion string

	// RemovedCallback is the function to call when the node has been removed
	RemovedCallback func()
}

func (b *RaftBackend) StartRecoveryCluster(ctx context.Context, peer Peer, removedCallback func()) error {
	recoveryModeConfig := &raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(peer.ID),
				Address: raft.ServerAddress(peer.Address),
			},
		},
	}

	return b.SetupCluster(context.Background(), SetupOpts{
		StartAsLeader:      true,
		RecoveryModeConfig: recoveryModeConfig,
		RemovedCallback:    removedCallback,
	})
}

func (b *RaftBackend) HasState() (bool, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	return raft.HasExistingState(b.logStore, b.stableStore, b.snapStore)
}

// SetupCluster starts the raft cluster and enables the networking needed for
// the raft nodes to communicate.
func (b *RaftBackend) SetupCluster(ctx context.Context, opts SetupOpts) error {
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

	if opts.EffectiveSDKVersion != "" {
		// Override the SDK version
		b.effectiveSDKVersion = opts.EffectiveSDKVersion
	}

	// Setup the raft config
	raftConfig := raft.DefaultConfig()
	if err := b.applyConfigSettings(raftConfig); err != nil {
		return err
	}

	listenerIsNil := func(cl cluster.ClusterHook) bool {
		switch {
		case opts.ClusterListener == nil:
			return true
		default:
			// Concrete type checks
			switch cl.(type) {
			case *cluster.Listener:
				return cl.(*cluster.Listener) == nil
			}
		}
		return false
	}

	var initialTimeoutMultiplier time.Duration
	switch {
	case opts.TLSKeyring == nil && listenerIsNil(opts.ClusterListener):
		// If we don't have a provided network we use an in-memory one.
		// This allows us to bootstrap a node without bringing up a cluster
		// network. This will be true during bootstrap, tests and dev modes.
		_, b.raftTransport = raft.NewInmemTransportWithTimeout(raft.ServerAddress(b.localID), time.Second)
	case opts.TLSKeyring == nil:
		return errors.New("no keyring provided")
	case listenerIsNil(opts.ClusterListener):
		return errors.New("no cluster listener provided")
	default:
		initialTimeoutMultiplier = 3
		if !opts.StartAsLeader {
			electionTimeout, heartbeatTimeout := raftConfig.ElectionTimeout, raftConfig.HeartbeatTimeout
			// Use bigger values for first election
			raftConfig.ElectionTimeout *= initialTimeoutMultiplier
			raftConfig.HeartbeatTimeout *= initialTimeoutMultiplier
			b.logger.Trace("using larger timeouts for raft at startup",
				"initial_election_timeout", raftConfig.ElectionTimeout,
				"initial_heartbeat_timeout", raftConfig.HeartbeatTimeout,
				"normal_election_timeout", electionTimeout,
				"normal_heartbeat_timeout", heartbeatTimeout)
		}

		// Set the local address and localID in the streaming layer and the raft config.
		streamLayer, err := NewRaftLayer(b.logger.Named("stream"), opts.TLSKeyring, opts.ClusterListener)
		if err != nil {
			return err
		}
		transConfig := &raft.NetworkTransportConfig{
			Stream:                  streamLayer,
			MaxPool:                 3,
			Timeout:                 10 * time.Second,
			ServerAddressProvider:   b.serverAddressProvider,
			Logger:                  b.logger.Named("raft-net"),
			MsgpackUseNewTimeFormat: true,
		}
		transport := raft.NewNetworkTransportWithConfig(transConfig)

		b.streamLayer = streamLayer
		b.raftTransport = transport
	}

	raftConfig.LocalID = raft.ServerID(b.localID)

	// Set up a channel for reliable leader notifications.
	raftNotifyCh := make(chan bool, 10)
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
	}

	// Setup the Raft store.
	b.fsm.SetNoopRestore(true)

	raftPath := filepath.Join(b.dataDir, raftState)
	peersFile := filepath.Join(raftPath, peersFileName)
	_, err := os.Stat(peersFile)
	if err == nil {
		b.logger.Info("raft recovery initiated", "recovery_file", peersFileName)

		recoveryConfig, err := raft.ReadConfigJSON(peersFile)
		if err != nil {
			return fmt.Errorf("raft recovery failed to parse peers.json: %w", err)
		}

		// Non-voting servers are only allowed in enterprise. If Suffrage is disabled,
		// error out to indicate that it isn't allowed.
		for idx := range recoveryConfig.Servers {
			if !nonVotersAllowed && recoveryConfig.Servers[idx].Suffrage == raft.Nonvoter {
				return fmt.Errorf("raft recovery failed to parse configuration for node %q: setting `non_voter` is only supported in enterprise", recoveryConfig.Servers[idx].ID)
			}
		}

		b.logger.Info("raft recovery found new config", "config", recoveryConfig)

		err = raft.RecoverCluster(raftConfig, b.fsm, b.logStore, b.stableStore, b.snapStore, b.raftTransport, recoveryConfig)
		if err != nil {
			return fmt.Errorf("raft recovery failed: %w", err)
		}

		err = os.Remove(peersFile)
		if err != nil {
			return fmt.Errorf("raft recovery failed to delete peers.json; please delete manually: %w", err)
		}
		b.logger.Info("raft recovery deleted peers.json")
	}

	if opts.RecoveryModeConfig != nil {
		err = raft.RecoverCluster(raftConfig, b.fsm, b.logStore, b.stableStore, b.snapStore, b.raftTransport, *opts.RecoveryModeConfig)
		if err != nil {
			return fmt.Errorf("recovering raft cluster failed: %w", err)
		}
	}

	b.logger.Info("creating Raft", "config", fmt.Sprintf("%#v", raftConfig))
	raftObj, err := raft.NewRaft(raftConfig, b.fsm.chunker, b.logStore, b.stableStore, b.snapStore, b.raftTransport)
	b.fsm.SetNoopRestore(false)
	if err != nil {
		return err
	}

	// If we are expecting to start as leader wait until we win the election.
	// This should happen quickly since there is only one node in the cluster.
	// StartAsLeader is only set during init, recovery mode, storage migration,
	// and tests.
	if opts.StartAsLeader {
		// ticker is used to prevent memory leak of using time.After in
		// for - select pattern.
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for {
			if raftObj.State() == raft.Leader {
				break
			}

			ticker.Reset(10 * time.Millisecond)
			select {
			case <-ctx.Done():
				future := raftObj.Shutdown()
				if future.Error() != nil {
					return fmt.Errorf("shutdown while waiting for leadership: %w", future.Error())
				}

				return errors.New("shutdown while waiting for leadership")
			case <-ticker.C:
			}
		}
	}

	b.raft = raftObj
	b.raftNotifyCh = raftNotifyCh

	if err := b.fsm.upgradeLocalNodeConfig(); err != nil {
		b.logger.Error("failed to upgrade local node configuration")
		return err
	}

	if b.streamLayer != nil {
		// Add Handler to the cluster.
		opts.ClusterListener.AddHandler(consts.RaftStorageALPN, b.streamLayer)

		// Add Client to the cluster.
		opts.ClusterListener.AddClient(consts.RaftStorageALPN, b.streamLayer)
	}

	// Close the init channel to signal setup has been completed
	close(b.raftInitCh)

	reloadConfig := func() {
		newCfg := raft.ReloadableConfig{
			TrailingLogs:      raftConfig.TrailingLogs,
			SnapshotInterval:  raftConfig.SnapshotInterval,
			SnapshotThreshold: raftConfig.SnapshotThreshold,
			HeartbeatTimeout:  raftConfig.HeartbeatTimeout / initialTimeoutMultiplier,
			ElectionTimeout:   raftConfig.ElectionTimeout / initialTimeoutMultiplier,
		}
		err := raftObj.ReloadConfig(newCfg)
		if err != nil {
			b.logger.Error("failed to reload raft config to set lower timeouts", "error", err)
		} else {
			b.logger.Trace("reloaded raft config to set lower timeouts", "config", fmt.Sprintf("%#v", newCfg))
		}
	}
	confFuture := raftObj.GetConfiguration()
	numServers := 0
	if err := confFuture.Error(); err != nil {
		// This should probably never happen, but just in case we'll log the error.
		// We'll default in this case to the multi-node behaviour.
		b.logger.Error("failed to read raft configuration", "error", err)
	} else {
		clusterConf := confFuture.Configuration()
		numServers = len(clusterConf.Servers)
	}
	if initialTimeoutMultiplier != 0 {
		if numServers == 1 {
			reloadConfig()
		} else {
			go func() {
				ticker := time.NewTicker(50 * time.Millisecond)
				// Emulate the random timeout used in Raft lib, to ensure that
				// if all nodes are brought up simultaneously, they don't all
				// call for an election at once.
				extra := time.Duration(rand.Int63()) % raftConfig.HeartbeatTimeout
				timeout := time.NewTimer(raftConfig.HeartbeatTimeout + extra)
				for {
					select {
					case <-ticker.C:
						switch raftObj.State() {
						case raft.Candidate, raft.Leader:
							b.logger.Trace("triggering raft config reload due to being candidate or leader")
							reloadConfig()
							return
						case raft.Shutdown:
							return
						}
					case <-timeout.C:
						b.logger.Trace("triggering raft config reload due to initial timeout")
						reloadConfig()
						return
					}
				}
			}()
		}
	}

	if opts.RemovedCallback != nil {
		b.removedCallback = opts.RemovedCallback
	}
	b.StartRemovedChecker(ctx)

	b.logger.Trace("finished setting up raft cluster")
	return nil
}

// TeardownCluster shuts down the raft cluster
func (b *RaftBackend) TeardownCluster(clusterListener cluster.ClusterHook) error {
	if clusterListener != nil {
		clusterListener.StopHandler(consts.RaftStorageALPN)
		clusterListener.RemoveClient(consts.RaftStorageALPN)
	}

	b.l.Lock()

	// Perform shutdown only if the raft object is non-nil. The object could be nil
	// if the node is unsealed but has not joined the peer set.
	var future raft.Future
	if b.raft != nil {
		future = b.raft.Shutdown()
	}

	b.raft = nil

	// If we're tearing down, then we need to recreate the raftInitCh
	b.raftInitCh = make(chan struct{})
	b.l.Unlock()

	if future != nil {
		return future.Error()
	}

	return nil
}

func (b *RaftBackend) StartRemovedChecker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		hasBeenPresent := false

		logger := b.logger.Named("removed.checker")
		for {
			select {
			case <-ticker.C:
				// If the raft cluster has been torn down (which will happen on
				// seal) the raft backend will be uninitialized. We want to exit
				// the loop in that case. If the cluster unseals, we'll get a
				// new backend setup and that will have its own removed checker.

				// There is a ctx.Done() check below that will also exit, but
				// in most (if not all) places we pass in context.Background()
				// to this function. Checking initialization will prevent this
				// loop from continuing to run after the raft backend is stopped
				// regardless of the context.
				if !b.Initialized() {
					return
				}
				removed, err := b.IsNodeRemoved(ctx, b.localID)
				if err != nil {
					logger.Error("failed to check if node is removed", "node ID", b.localID, "error", err)
					continue
				}
				if !removed {
					hasBeenPresent = true
				}
				// the node must have been previously present in the config,
				// only then should we consider it removed and shutdown
				if removed && hasBeenPresent {
					err := b.RemoveSelf()
					if err != nil {
						logger.Error("failed to remove self", "node ID", b.localID, "error", err)
					}
					b.l.RLock()
					if b.removedCallback != nil {
						b.removedCallback()
					}
					b.l.RUnlock()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// CommittedIndex returns the latest index committed to stable storage
func (b *RaftBackend) CommittedIndex() uint64 {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return 0
	}

	return b.raft.LastIndex()
}

// AppliedIndex returns the latest index applied to the FSM
func (b *RaftBackend) AppliedIndex() uint64 {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.fsm == nil {
		return 0
	}

	// We use the latest index that the FSM has seen here, which may be behind
	// raft.AppliedIndex() due to the async nature of the raft library.
	indexState, _ := b.fsm.LatestState()
	return indexState.Index
}

// Term returns the raft term of this node.
func (b *RaftBackend) Term() uint64 {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.fsm == nil {
		return 0
	}

	// We use the latest index that the FSM has seen here, which may be behind
	// raft.AppliedIndex() due to the async nature of the raft library.
	indexState, _ := b.fsm.LatestState()
	return indexState.Term
}

// RemovePeer removes the given peer ID from the raft cluster. If the node is
// ourselves we will give up leadership.
func (b *RaftBackend) RemovePeer(ctx context.Context, peerID string) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	if b.disableAutopilot {
		if b.raft == nil {
			return errors.New("raft storage is not initialized")
		}
		b.logger.Trace("removing server from raft", "id", peerID)
		future := b.raft.RemoveServer(raft.ServerID(peerID), 0, 0)
		return future.Error()
	}

	if b.autopilot == nil {
		return errors.New("raft storage autopilot is not initialized")
	}

	b.logger.Trace("removing server from raft via autopilot", "id", peerID)
	return b.autopilot.RemoveServer(raft.ServerID(peerID))
}

// GetConfigurationOffline is used to read the stale, last known raft
// configuration to this node. It accesses the last state written into the
// FSM. When a server is online use GetConfiguration instead.
func (b *RaftBackend) GetConfigurationOffline() (*RaftConfigurationResponse, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft != nil {
		return nil, errors.New("raft storage is initialized, used GetConfiguration instead")
	}

	if b.fsm == nil {
		return nil, nil
	}

	state, configuration := b.fsm.LatestState()
	config := &RaftConfigurationResponse{
		Index: state.Index,
	}

	if configuration == nil || configuration.Servers == nil {
		return config, nil
	}

	for _, server := range configuration.Servers {
		entry := &RaftServer{
			NodeID:  server.Id,
			Address: server.Address,
			// Since we are offline no node is the leader.
			Leader: false,
			Voter:  raft.ServerSuffrage(server.Suffrage) == raft.Voter,
		}
		config.Servers = append(config.Servers, entry)
	}

	return config, nil
}

func (b *RaftBackend) GetConfiguration(ctx context.Context) (*RaftConfigurationResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

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
			NodeID:  string(server.ID),
			Address: string(server.Address),
			// Since we only service this request on the active node our node ID
			// denotes the raft leader.
			Leader:          string(server.ID) == b.NodeID(),
			Voter:           server.Suffrage == raft.Voter,
			ProtocolVersion: strconv.Itoa(raft.ProtocolVersionMax),
		}

		config.Servers = append(config.Servers, entry)
	}

	return config, nil
}

// AddPeer adds a new server to the raft cluster
func (b *RaftBackend) AddPeer(ctx context.Context, peerID, clusterAddr string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	b.l.RLock()
	defer b.l.RUnlock()

	if b.disableAutopilot {
		if b.raft == nil {
			return errors.New("raft storage is not initialized")
		}
		b.logger.Trace("adding server to raft", "id", peerID, "addr", clusterAddr)
		future := b.raft.AddVoter(raft.ServerID(peerID), raft.ServerAddress(clusterAddr), 0, 0)
		return future.Error()
	}

	if b.autopilot == nil {
		return errors.New("raft storage autopilot is not initialized")
	}

	b.logger.Trace("adding server to raft via autopilot", "id", peerID, "addr", clusterAddr)
	return b.autopilot.AddServer(&autopilot.Server{
		ID:          raft.ServerID(peerID),
		Name:        peerID,
		Address:     raft.ServerAddress(clusterAddr),
		RaftVersion: raft.ProtocolVersionMax,
		NodeType:    autopilot.NodeVoter,
	})
}

// Peers returns all the servers present in the raft cluster
func (b *RaftBackend) Peers(ctx context.Context) ([]Peer, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return nil, errors.New("raft storage is not initialized")
	}

	future := b.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		return nil, err
	}

	ret := make([]Peer, len(future.Configuration().Servers))
	for i, s := range future.Configuration().Servers {
		ret[i] = Peer{
			ID:       string(s.ID),
			Address:  string(s.Address),
			Suffrage: int(s.Suffrage),
		}
	}

	return ret, nil
}

// SnapshotHTTP is a wrapper for Snapshot that sends the snapshot as an HTTP
// response.
func (b *RaftBackend) SnapshotHTTP(out *logical.HTTPResponseWriter, sealer snapshot.Sealer) error {
	out.Header().Add("Content-Disposition", "attachment")
	out.Header().Add("Content-Type", "application/gzip")

	return b.Snapshot(out, sealer)
}

// Snapshot takes a raft snapshot, packages it into a archive file and writes it
// to the provided writer. Seal access is used to encrypt the SHASUM file so we
// can validate the snapshot was taken using the same root keys or not.
func (b *RaftBackend) Snapshot(out io.Writer, sealer snapshot.Sealer) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage is sealed")
	}

	return snapshot.Write(b.logger.Named("snapshot"), b.raft, sealer, out)
}

// WriteSnapshotToTemp reads a snapshot archive off the provided reader,
// extracts the data and writes the snapshot to a temporary file. The seal
// access is used to decrypt the SHASUM file in the archive to ensure this
// snapshot has the same root key as the running instance. If the provided
// access is nil then it will skip that validation.
func (b *RaftBackend) WriteSnapshotToTemp(in io.ReadCloser, sealer snapshot.Sealer) (*os.File, func(), raft.SnapshotMeta, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	var metadata raft.SnapshotMeta
	if b.raft == nil {
		return nil, nil, metadata, errors.New("raft storage is sealed")
	}

	snap, cleanup, err := snapshot.WriteToTempFileWithSealer(b.logger.Named("snapshot"), in, &metadata, sealer)
	return snap, cleanup, metadata, err
}

// RestoreSnapshot applies the provided snapshot metadata and snapshot data to
// raft.
func (b *RaftBackend) RestoreSnapshot(ctx context.Context, metadata raft.SnapshotMeta, snap io.Reader) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage is not initialized")
	}

	if err := b.raft.Restore(&metadata, snap, 0); err != nil {
		b.logger.Named("snapshot").Error("failed to restore snapshot", "error", err)
		return err
	}

	// Apply a log that tells the follower nodes to run the restore callback
	// function. This is done after the restore call so we can be sure the
	// snapshot applied to a quorum of nodes.
	command := &LogData{
		Operations: []*LogOperation{
			{
				OpType: restoreCallbackOp,
			},
		},
	}

	err := b.applyLog(ctx, command)

	// Do a best-effort attempt to let the standbys apply the restoreCallbackOp
	// before we continue.
	time.Sleep(restoreOpDelayDuration)
	return err
}

// Delete inserts an entry in the log to delete the given path
func (b *RaftBackend) Delete(ctx context.Context, path string) error {
	defer metrics.MeasureSince([]string{"raft-storage", "delete"}, time.Now())

	if err := ctx.Err(); err != nil {
		return err
	}

	command := &LogData{
		Operations: []*LogOperation{
			{
				OpType: deleteOp,
				Key:    path,
			},
		},
	}
	if err := b.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer b.permitPool.Release()

	b.l.RLock()
	err := b.applyLog(ctx, command)
	b.l.RUnlock()
	return err
}

// Get returns the value corresponding to the given path from the fsm
func (b *RaftBackend) Get(ctx context.Context, path string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"raft-storage", "get"}, time.Now())
	if b.fsm == nil {
		return nil, errors.New("raft: fsm not configured")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := b.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer b.permitPool.Release()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	entry, err := b.fsm.Get(ctx, path)
	if entry != nil {
		valueLen := len(entry.Value)
		maxSize := b.entrySizeLimitForPath(path)
		if uint64(valueLen) > maxSize {
			b.logger.Warn("retrieved entry value is too large, has raft's max_entry_size or max_mount_and_namespace_table_entry_size been reduced?",
				"size", valueLen, "max_size", maxSize)
		}
	}

	return entry, err
}

// Put inserts an entry in the log for the put operation. It will return an
// error if the resulting entry encoding exceeds the configured max_entry_size
// or if the call to applyLog fails.
func (b *RaftBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"raft-storage", "put"}, time.Now())
	if len(entry.Key) > bolt.MaxKeySize {
		return fmt.Errorf("%s, max key size for integrated storage is %d", physical.ErrKeyTooLarge, bolt.MaxKeySize)
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	command := &LogData{
		Operations: []*LogOperation{
			{
				OpType: putOp,
				Key:    entry.Key,
				Value:  entry.Value,
			},
		},
	}

	if err := b.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer b.permitPool.Release()

	b.l.RLock()
	err := b.applyLog(ctx, command)
	b.l.RUnlock()
	return err
}

// List enumerates all the items under the prefix from the fsm
func (b *RaftBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"raft-storage", "list"}, time.Now())
	if b.fsm == nil {
		return nil, errors.New("raft: fsm not configured")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := b.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer b.permitPool.Release()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return b.fsm.List(ctx, prefix)
}

// Transaction applies all the given operations into a single log and
// applies it.
func (b *RaftBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	defer metrics.MeasureSince([]string{"raft-storage", "transaction"}, time.Now())

	if err := ctx.Err(); err != nil {
		return err
	}

	failGetInTxn := atomic.LoadUint32(b.failGetInTxn)
	for _, t := range txns {
		if t.Operation == physical.GetOperation && failGetInTxn != 0 {
			return GetInTxnDisabledError
		}
	}

	txnMap := make(map[string]*physical.TxnEntry)

	command := &LogData{
		Operations: make([]*LogOperation, len(txns)),
	}
	for i, txn := range txns {
		op := &LogOperation{}
		switch txn.Operation {
		case physical.PutOperation:
			if len(txn.Entry.Key) > bolt.MaxKeySize {
				return fmt.Errorf("%s, max key size for integrated storage is %d", physical.ErrKeyTooLarge, bolt.MaxKeySize)
			}
			op.OpType = putOp
			op.Key = txn.Entry.Key
			op.Value = txn.Entry.Value
		case physical.DeleteOperation:
			op.OpType = deleteOp
			op.Key = txn.Entry.Key
		case physical.GetOperation:
			op.OpType = getOp
			op.Key = txn.Entry.Key
			txnMap[op.Key] = txn
		default:
			return fmt.Errorf("%q is not a supported transaction operation", txn.Operation)
		}

		command.Operations[i] = op
	}

	if err := b.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer b.permitPool.Release()

	b.l.RLock()
	err := b.applyLog(ctx, command)
	b.l.RUnlock()

	// loop over results and update pointers to get operations
	for _, logOp := range command.Operations {
		if logOp.OpType == getOp {
			if txn, found := txnMap[logOp.Key]; found {
				txn.Entry.Value = logOp.Value
			}
		}
	}

	return err
}

func (b *RaftBackend) TransactionLimits() (int, int) {
	return b.maxBatchEntries, b.maxBatchSize
}

// validateCommandEntrySizes is a helper to check the size of each transaction
// value isn't larger than is allowed. It must take into account the path in
// case any special limits have been set for mount table paths. Finally it
// returns the largest limit it needed to use so that calling code can size the
// overall log entry check correctly. In other words if max_entry_size is 1MB
// and max_mount_and_namespace_table_entry_size is 2MB, we check each key
// against the right limit and then return 1MB unless there is at least one
// mount table key being written in which case we allow the larger limit of 2MB.
func (b *RaftBackend) validateCommandEntrySizes(command *LogData) (uint64, error) {
	largestEntryLimit := b.maxEntrySize

	for _, op := range command.Operations {
		if op.OpType == putOp {
			entrySize := b.entrySizeLimitForPath(op.Key)
			if len(op.Value) > int(entrySize) {
				return 0, fmt.Errorf("%s, max value size for key %s is %d, got %d", physical.ErrValueTooLarge, op.Key, entrySize, len(op.Value))
			}
			if entrySize > largestEntryLimit {
				largestEntryLimit = entrySize
			}
		}
	}
	return largestEntryLimit, nil
}

// applyLog will take a given log command and apply it to the raft log. applyLog
// doesn't return until the log has been applied to a quorum of servers and is
// persisted to the local FSM. Caller should hold the backend's read lock.
func (b *RaftBackend) applyLog(ctx context.Context, command *LogData) error {
	if b.raft == nil {
		return errors.New("raft storage is not initialized")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	totalLogSizeLimit, err := b.validateCommandEntrySizes(command)
	if err != nil {
		return err
	}

	commandBytes, err := proto.Marshal(command)
	if err != nil {
		return err
	}

	cmdSize := len(commandBytes)
	if uint64(cmdSize) > totalLogSizeLimit {
		return fmt.Errorf("%s; got %d bytes, max: %d bytes", physical.ErrValueTooLarge, cmdSize, totalLogSizeLimit)
	}

	defer metrics.AddSample([]string{"raft-storage", "entry_size"}, float32(cmdSize))

	var chunked bool
	var applyFuture raft.ApplyFuture
	switch {
	case len(commandBytes) <= raftchunking.ChunkSize:
		applyFuture = b.raft.Apply(commandBytes, 0)
	default:
		chunked = true
		applyFuture = raftchunking.ChunkingApply(commandBytes, nil, 0, b.raft.ApplyLog)
	}

	if err := applyFuture.Error(); err != nil {
		return err
	}

	resp := applyFuture.Response()

	if chunked {
		// In this case we didn't apply all chunks successfully, possibly due
		// to a term change
		if resp == nil {
			// This returns the error in the interface because the raft library
			// returns errors from the FSM via the future, not via err from the
			// apply function. Downstream client code expects to see any error
			// from the FSM (as opposed to the apply itself) and decide whether
			// it can retry in the future's response.
			return errors.New("applying chunking failed, please retry")
		}

		// We expect that this conversion should always work
		chunkedSuccess, ok := resp.(raftchunking.ChunkingSuccess)
		if !ok {
			return errors.New("unknown type of response back from chunking FSM")
		}

		// Replace the reply with the inner wrapped version
		resp = chunkedSuccess.Response
	}

	fsmar, ok := resp.(*FSMApplyResponse)
	if !ok || !fsmar.Success {
		return errors.New("could not apply data")
	}

	// populate command with our results
	if fsmar.EntrySlice == nil {
		return errors.New("entries on FSM response were empty")
	}

	for i, logOp := range command.Operations {
		if logOp.OpType == getOp {
			fsmEntry := fsmar.EntrySlice[i]

			// this should always be true because the entries in the slice were created in the same order as
			// the command operations.
			if logOp.Key == fsmEntry.Key {
				if len(fsmEntry.Value) > 0 {
					logOp.Value = fsmEntry.Value
				}
			} else {
				// this shouldn't happen
				return errors.New("entries in FSM response were out of order")
			}
		}
	}

	return nil
}

// HAEnabled is the implementation of the HABackend interface
func (b *RaftBackend) HAEnabled() bool { return true }

// HAEnabled is the implementation of the HABackend interface
func (b *RaftBackend) LockWith(key, value string) (physical.Lock, error) {
	return &RaftLock{
		key:   key,
		value: []byte(value),
		b:     b,
	}, nil
}

// SetDesiredSuffrage sets a field in the fsm indicating the suffrage intent for
// this node.
func (b *RaftBackend) SetDesiredSuffrage(nonVoter bool) error {
	b.l.Lock()
	defer b.l.Unlock()

	var desiredSuffrage string
	switch nonVoter {
	case true:
		desiredSuffrage = "non-voter"
	default:
		desiredSuffrage = "voter"
	}

	err := b.fsm.recordSuffrage(desiredSuffrage)
	if err != nil {
		return err
	}

	return nil
}

func (b *RaftBackend) DesiredSuffrage() string {
	return b.fsm.DesiredSuffrage()
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
func (l *RaftLock) monitorLeadership(stopCh <-chan struct{}, leaderNotifyCh <-chan bool) <-chan struct{} {
	leaderLost := make(chan struct{})
	go func() {
		for {
			select {
			case isLeader := <-leaderNotifyCh:
				// leaderNotifyCh may deliver a true value initially if this
				// server is already the leader prior to RaftLock.Lock call
				// (the true message was already queued). The next message is
				// always going to be false. The for loop should loop at most
				// twice.
				if !isLeader {
					close(leaderLost)
					return
				}
			case <-stopCh:
				return
			}
		}
	}()
	return leaderLost
}

// Lock blocks until we become leader or are shutdown. It returns a channel that
// is closed when we detect a loss of leadership.
func (l *RaftLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	// If not initialized, block until it is
	if !l.b.Initialized() {
		select {
		case <-l.b.raftInitCh:
		case <-stopCh:
			return nil, nil
		}
	}

	l.b.l.RLock()

	// Ensure that we still have a raft instance after grabbing the read lock
	if l.b.raft == nil {
		l.b.l.RUnlock()
		return nil, errors.New("attempted to grab a lock on a nil raft backend")
	}

	// Cache the notifyCh locally
	leaderNotifyCh := l.b.raftNotifyCh

	// Check to see if we are already leader.
	if l.b.raft.State() == raft.Leader {
		err := l.b.applyLog(context.Background(), &LogData{
			Operations: []*LogOperation{
				{
					OpType: putOp,
					Key:    l.key,
					Value:  l.value,
				},
			},
		})
		l.b.l.RUnlock()
		if err != nil {
			return nil, err
		}

		return l.monitorLeadership(stopCh, leaderNotifyCh), nil
	}
	l.b.l.RUnlock()

	for {
		select {
		case isLeader := <-leaderNotifyCh:
			if isLeader {
				// We are leader, set the key
				l.b.l.RLock()
				err := l.b.applyLog(context.Background(), &LogData{
					Operations: []*LogOperation{
						{
							OpType: putOp,
							Key:    l.key,
							Value:  l.value,
						},
					},
				})
				l.b.l.RUnlock()
				if err != nil {
					return nil, err
				}

				return l.monitorLeadership(stopCh, leaderNotifyCh), nil
			}
		case <-stopCh:
			return nil, nil
		}
	}
}

// Unlock gives up leadership.
func (l *RaftLock) Unlock() error {
	if l.b.raft == nil {
		return nil
	}

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

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		// File exists!
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	// We hit some other error trying to stat the file which leaves us in an
	// unknown state so we can't proceed.
	return false, err
}

func isRaftLogVerifyCheckpoint(l *raft.Log) bool {
	if !bytes.Equal(l.Data, []byte{byte(verifierCheckpointOp)}) {
		return false
	}

	// Single byte log with that byte value can only be a checkpoint or
	// the last byte of a chunked message. If it's chunked it will have
	// chunking metadata.
	if len(l.Extensions) == 0 {
		// No metadata, must be a checkpoint on the leader with no
		// verifier metadata yet.
		return true
	}

	if bytes.HasPrefix(l.Extensions, logVerifierMagicBytes[:]) {
		// Has verifier metadata so must be a replicated checkpoint on a follower
		return true
	}

	// Must be the last chunk of a chunked object that has chunking meta
	return false
}
