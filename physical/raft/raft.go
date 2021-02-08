package raft

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/armon/go-metrics"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/go-raftchunking"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/raft"
	autopilot "github.com/hashicorp/raft-autopilot"
	snapshot "github.com/hashicorp/raft-snapshot"
	raftboltdb "github.com/hashicorp/vault/physical/raft/logstore"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/hashicorp/vault/vault/seal"
)

// EnvVaultRaftNodeID is used to fetch the Raft node ID from the environment.
const EnvVaultRaftNodeID = "VAULT_RAFT_NODE_ID"

// EnvVaultRaftPath is used to fetch the path where Raft data is stored from the environment.
const EnvVaultRaftPath = "VAULT_RAFT_PATH"

// Verify RaftBackend satisfies the correct interfaces
var _ physical.Backend = (*RaftBackend)(nil)
var _ physical.Transactional = (*RaftBackend)(nil)
var _ physical.HABackend = (*RaftBackend)(nil)
var _ physical.Lock = (*RaftLock)(nil)

var (
	// raftLogCacheSize is the maximum number of logs to cache in-memory.
	// This is used to reduce disk I/O for the recently committed entries.
	raftLogCacheSize = 512

	raftState     = "raft/"
	peersFileName = "peers.json"

	restoreOpDelayDuration = 5 * time.Second

	defaultMaxEntrySize = uint64(2 * raftchunking.ChunkSize)
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

	// dataDir is the location on the local filesystem that raft and FSM data
	// will be stored.
	dataDir string

	// localID is the ID for this node. This can either be configured in the
	// config file, via a file on disk, or is otherwise randomly generated.
	localID string

	// serverAddressProvider is used to map server IDs to addresses.
	serverAddressProvider raft.ServerAddressProvider

	// permitPool is used to limit the number of concurrent storage calls.
	permitPool *physical.PermitPool

	// maxEntrySize imposes a size limit (in bytes) on a raft entry (put or transaction).
	// It is suggested to use a value of 2x the Raft chunking size for optimal
	// performance.
	maxEntrySize uint64

	// autopilot is the instance of raft-autopilot library implementation of the
	// autopilot features. This will be instantiated in both leader and followers.
	// However, only active node will have a "running" autopilot.
	autopilot *autopilot.Autopilot

	// disableAutopilot if set will not put autopilot implementation to use. The
	// fallback will be to interact with the raft instance directly. This is useful
	// to retain state of operations on existing clusters which do not have
	// autopilot enabled yet. For new clusters however, the default will be to have
	// autopilot enabled from get-go.
	disableAutopilot bool

	// autopilotConfig represents the configuration required to instantiate autopilot.
	autopilotConfig *AutopilotConfig

	// followerStates represents the information about all the peers of the raft
	// leader. This is used to track some state of the peers and as well as used
	// to see if the peers are "alive" using the heartbeat received from them.
	followerStates *FollowerStates

	// followerHeartbeatTrackerStopCh is used to terminate the goroutine that
	// keeps track of follower heartbeats.
	followerHeartbeatTrackerStopCh chan struct{}
}

// AutopilotConfiguration is used for querying/setting the Autopilot configuration.
// Autopilot helps manage operator tasks related to Vault servers like removing
// failed servers from the Raft quorum.
type AutopilotConfig struct {
	// CleanupDeadServers controls whether to remove dead servers from the Raft
	// peer list periodically or when a new server joins
	CleanupDeadServers bool `json:"cleanup_dead_servers" mapstructure:"cleanup_dead_servers"`

	// LastContactThreshold is the limit on the amount of time a server can go
	// without leader contact before being considered unhealthy.
	LastContactThreshold time.Duration `json:"last_contact_threshold" mapstructure:"-"`

	// MaxTrailingLogs is the amount of entries in the Raft Log that a server can
	// be behind before being considered unhealthy.
	MaxTrailingLogs uint64 `json:"max_trailing_logs" mapstructure:"max_trailing_logs"`

	// MinQuorum sets the minimum number of servers allowed in a cluster before
	// autopilot can prune dead servers.
	MinQuorum uint `json:"min_quorum" mapstructure:"min_quorum"`

	// ServerStabilizationTime is the minimum amount of time a server must be
	// in a stable, healthy state before it can be added to the cluster. Only
	// applicable with Raft protocol version 3 or higher.
	ServerStabilizationTime time.Duration `json:"server_stabilization_time" mapstructure:"-"`
}

// Clone returns a duplicate instance of AutopilotConfig with the exact same values.
func (ac *AutopilotConfig) Clone() *AutopilotConfig {
	return &AutopilotConfig{
		CleanupDeadServers:      ac.CleanupDeadServers,
		LastContactThreshold:    ac.LastContactThreshold,
		MaxTrailingLogs:         ac.MaxTrailingLogs,
		MinQuorum:               ac.MinQuorum,
		ServerStabilizationTime: ac.ServerStabilizationTime,
	}

}

// MarshalJSON makes the autopilot config fields JSON compatible
func (ac *AutopilotConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"cleanup_dead_servers":      ac.CleanupDeadServers,
		"last_contact_threshold":    ac.LastContactThreshold.String(),
		"max_trailing_logs":         ac.MaxTrailingLogs,
		"min_quorum":                ac.MinQuorum,
		"server_stabilization_time": ac.ServerStabilizationTime.String(),
	})
}

// UnmarshalJSON parses the autopilot config JSON blob
func (ac *AutopilotConfig) UnmarshalJSON(b []byte) error {
	var data interface{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	conf := data.(map[string]interface{})
	if err = mapstructure.WeakDecode(conf, ac); err != nil {
		return err
	}
	if ac.LastContactThreshold, err = parseutil.ParseDurationSecond(conf["last_contact_threshold"]); err != nil {
		return err
	}
	if ac.ServerStabilizationTime, err = parseutil.ParseDurationSecond(conf["server_stabilization_time"]); err != nil {
		return err
	}

	return nil
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

	// LeaderCACertFile is the path on disk to the the CA cert file of the
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

// FollowerState represents the information about peer that the leader tracks.
type FollowerState struct {
	AppliedIndex  uint64
	LastHeartbeat time.Time
	LastTerm      uint64
	IsDead        bool
}

// FollowerStates holds information about all the peers in the raft cluster
// tracked by the leader.
type FollowerStates struct {
	l         sync.RWMutex
	followers map[string]FollowerState
}

// NewFollowerStates creates a new FollowerStates object
func NewFollowerStates() *FollowerStates {
	return &FollowerStates{
		followers: make(map[string]FollowerState),
	}
}

// MarkFollowerAsDead marks the node represented by the nodeID as dead. This
// implies that Vault will indicate that the node "left" the cluster, the next
// time autopilot asks for known servers.
func (s *FollowerStates) MarkFollowerAsDead(nodeID string) {
	s.l.Lock()
	defer s.l.Unlock()

	state, ok := s.followers[nodeID]
	if !ok {
		return
	}
	s.followers[nodeID] = FollowerState{
		LastHeartbeat: state.LastHeartbeat,
		LastTerm:      state.LastTerm,
		AppliedIndex:  state.AppliedIndex,
		IsDead:        true,
	}
}

// Update the peer information in the follower states
func (s *FollowerStates) Update(nodeID string, appliedIndex uint64, term uint64) {
	state := FollowerState{
		AppliedIndex: appliedIndex,
		LastTerm:     term,
	}
	if appliedIndex > 0 {
		state.LastHeartbeat = time.Now()
	}

	s.l.Lock()
	s.followers[nodeID] = state
	s.l.Unlock()
}

// Clear wipes all the information regarding peers in the follower states.
func (s *FollowerStates) Clear() {
	s.l.Lock()
	for i := range s.followers {
		delete(s.followers, i)
	}
	s.l.Unlock()
}

// Delete the entry of a peer represented by the nodeID from follower states.
func (s *FollowerStates) Delete(nodeID string) {
	s.l.Lock()
	delete(s.followers, nodeID)
	s.l.Unlock()
}

// Get retrives information about a peer from the followerstates.
func (s *FollowerStates) Get(nodeID string) FollowerState {
	s.l.RLock()
	state := s.followers[nodeID]
	s.l.RUnlock()
	return state
}

// MinIndex returns the minimum raft index applied in the raft cluster.
func (s *FollowerStates) MinIndex() uint64 {
	var min uint64 = math.MaxUint64
	minFunc := func(a, b uint64) uint64 {
		if a > b {
			return b
		}
		return a
	}

	s.l.RLock()
	for _, state := range s.followers {
		min = minFunc(min, state.AppliedIndex)
	}
	s.l.RUnlock()

	if min == math.MaxUint64 {
		return 0
	}

	return min
}

// Ensure that the Delegate implements the ApplicationIntegration interface
var _ autopilot.ApplicationIntegration = (*Delegate)(nil)

// Delegate is an implementation of autopilot.ApplicationIntegration interface.
// This is used by the autopilot library to retrieve information and to have
// application specific tasks performed.
type Delegate struct {
	*RaftBackend
}

// AutopilotConfig is called by the autopilot library to know the desired
// autopilot configuration.
func (d *Delegate) AutopilotConfig() *autopilot.Config {
	d.l.RLock()
	defer d.l.RUnlock()

	return &autopilot.Config{
		CleanupDeadServers:      d.autopilotConfig.CleanupDeadServers,
		LastContactThreshold:    d.autopilotConfig.LastContactThreshold,
		MaxTrailingLogs:         d.autopilotConfig.MaxTrailingLogs,
		MinQuorum:               d.autopilotConfig.MinQuorum,
		ServerStabilizationTime: d.autopilotConfig.ServerStabilizationTime,
	}
}

// NotifyState is called by the autopilot library whenever there is a state
// change. We update a few metrics when this happens.
func (d *Delegate) NotifyState(state *autopilot.State) {
	d.logger.Info("notify state called")
	if d.raft.State() == raft.Leader {
		metrics.SetGauge([]string{"autopilot", "failure_tolerance"}, float32(state.FailureTolerance))
		if state.Healthy {
			metrics.SetGauge([]string{"autopilot", "healthy"}, 1)
		} else {
			metrics.SetGauge([]string{"autopilot", "healthy"}, 0)
		}
	}
}

// FetchServerStats is called by the autopilot library to retrieve information
// about all the nodes in the raft cluster.
func (d *Delegate) FetchServerStats(ctx context.Context, servers map[raft.ServerID]*autopilot.Server) map[raft.ServerID]*autopilot.ServerStats {
	ret := make(map[raft.ServerID]*autopilot.ServerStats)

	followerStates := d.RaftBackend.followerStates
	followerStates.l.RLock()
	defer followerStates.l.RUnlock()

	now := time.Now()
	for id, followerState := range followerStates.followers {
		ret[raft.ServerID(id)] = &autopilot.ServerStats{
			LastContact: now.Sub(followerState.LastHeartbeat),
			LastTerm:    followerState.LastTerm,
			LastIndex:   followerState.AppliedIndex,
		}
	}

	leaderState, _ := d.fsm.LatestState()
	ret[raft.ServerID(d.localID)] = &autopilot.ServerStats{
		LastTerm:  leaderState.Term,
		LastIndex: leaderState.Index,
	}

	return ret
}

// KnownServers is called by the autopilot library to know the status of each
// node in the raft cluster. If the application thinks that certain nodes left,
// it is here that we let the autopilot library know of the same.
func (d *Delegate) KnownServers() map[raft.ServerID]*autopilot.Server {
	followerStates := d.RaftBackend.followerStates
	followerStates.l.RLock()
	defer followerStates.l.RUnlock()

	d.logger.Trace("returning known servers")
	ret := make(map[raft.ServerID]*autopilot.Server)
	for id, state := range d.RaftBackend.followerStates.followers {
		server := &autopilot.Server{
			ID:          raft.ServerID(id),
			Name:        id,
			RaftVersion: raft.ProtocolVersionMax,
		}

		switch state.IsDead {
		case true:
			d.logger.Info("informing autopilot that the node left", "id", id)
			server.NodeStatus = autopilot.NodeLeft
		default:
			server.NodeStatus = autopilot.NodeAlive
		}

		ret[raft.ServerID(id)] = server
		d.logger.Trace("returning known server", "id", id)
	}

	// Add the leader
	ret[raft.ServerID(d.localID)] = &autopilot.Server{
		ID:          raft.ServerID(d.localID),
		Name:        d.localID,
		RaftVersion: raft.ProtocolVersionMax,
		NodeStatus:  autopilot.NodeAlive,
	}
	d.logger.Trace("returning known server", "id", d.localID)

	return ret
}

// RemoveFailedServer is called by the autopilot library when it desires a node
// to be removed from the raft configuration. This function removes the node
// from the raft cluster and stops tracking its information in follower states.
// This function needs to return quickly. Hence removal is performed in a
// goroutine.
func (d *Delegate) RemoveFailedServer(server *autopilot.Server) {
	go func() {
		d.logger.Info("removing dead server from raft configuration", "id", server.ID)
		if future := d.raft.RemoveServer(server.ID, 0, 0); future.Error() != nil {
			d.logger.Error("failed to remove server", "server_id", server.ID, "server_address", server.Address, "server_name", server.Name, "error", future.Error())
		}
		d.followerStates.Delete(string(server.ID))
	}()
}

// SetFollowerStates sets the followerStates field in the backend to track peers
// in the raft cluster.
func (b *RaftBackend) SetFollowerStates(states *FollowerStates) {
	b.l.Lock()
	b.followerStates = states
	b.l.Unlock()
}

// SetAutopilotConfig updates the autopilot configuration in the backend.
func (b *RaftBackend) SetAutopilotConfig(config *AutopilotConfig) {
	b.l.Lock()
	b.autopilotConfig = config
	b.l.Unlock()
}

// AutopilotConfig returns the autopilot configuration in the backend.
func (b *RaftBackend) AutopilotConfig() *AutopilotConfig {
	b.l.RLock()
	conf := b.autopilotConfig
	b.l.RUnlock()
	return conf
}

func (b *RaftBackend) defaultAutopilotConfig() *AutopilotConfig {
	return &AutopilotConfig{
		CleanupDeadServers:      true,
		LastContactThreshold:    10 * time.Second,
		MaxTrailingLogs:         1000,
		MinQuorum:               3,
		ServerStabilizationTime: 10 * time.Second,
	}
}

func (b *RaftBackend) autopilotConf() (*AutopilotConfig, error) {
	config := b.conf["autopilot"]
	if config == "" {
		return nil, nil
	}

	// TODO: Find out why we are getting a list instead of a single item
	var configs []*AutopilotConfig
	err := jsonutil.DecodeJSON([]byte(config), &configs)
	if err != nil {
		return nil, errwrap.Wrapf("failed to decode autopilot config: {{err}}", err)
	}
	if len(configs) != 1 {
		return nil, fmt.Errorf("expected a single block of autopilot config")
	}

	return configs[0], nil
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
		return nil, errwrap.Wrapf("failed to decode retry_join config: {{err}}", err)
	}

	if len(leaderInfos) == 0 {
		return nil, errors.New("invalid retry_join config")
	}

	for i, info := range leaderInfos {
		if len(info.AutoJoin) != 0 && len(info.LeaderAPIAddr) != 0 {
			return nil, errors.New("cannot provide both a leader_api_addr and auto_join")
		}

		if info.AutoJoinScheme != "" && (info.AutoJoinScheme != "http" && info.AutoJoinScheme != "https") {
			return nil, fmt.Errorf("invalid scheme '%s'; must either be http or https", info.AutoJoinScheme)
		}

		info.Retry = true
		info.TLSConfig, err = parseTLSInfo(info)
		if err != nil {
			return nil, errwrap.Wrapf(fmt.Sprintf("failed to create tls config to communicate with leader node (retry_join index: %d): {{err}}", i), err)
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
	return os.MkdirAll(path, 0755)
}

// NewRaftBackend constructs a RaftBackend using the given directory
func NewRaftBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {

	path := os.Getenv(EnvVaultRaftPath)
	if path == "" {
		pathFromConfig, ok := conf["path"]
		if !ok {
			return nil, fmt.Errorf("'path' must be set")
		}
		path = pathFromConfig
	}

	// Create the FSM.
	fsm, err := NewFSM(path, logger.Named("fsm"))
	if err != nil {
		return nil, fmt.Errorf("failed to create fsm: %v", err)
	}

	// Build an all in-memory setup for dev mode, otherwise prepare a full
	// disk-based setup.
	var log raft.LogStore
	var stable raft.StableStore
	var snap raft.SnapshotStore

	var devMode bool
	if devMode {
		store := raft.NewInmemStore()
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
		snapshots, err := NewBoltSnapshotStore(path, logger.Named("snapshot"), fsm)
		if err != nil {
			return nil, err
		}
		snap = snapshots
	}

	if delayRaw, ok := conf["snapshot_delay"]; ok {
		delay, err := time.ParseDuration(delayRaw)
		if err != nil {
			return nil, fmt.Errorf("snapshot_delay does not parse as a duration: %w", err)
		}
		snap = newSnapshotStoreDelay(snap, delay)
	}

	var localID string
	{
		// Determine the local node ID from the environment.
		if raftNodeID := os.Getenv(EnvVaultRaftNodeID); raftNodeID != "" {
			localID = raftNodeID
		}

		// If not set in the environment check the configuration file.
		if len(localID) == 0 {
			localID = conf["node_id"]
		}

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

		// If all of the above fails generate a UUID and persist it to the
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

	maxEntrySize := defaultMaxEntrySize
	if maxEntrySizeCfg := conf["max_entry_size"]; len(maxEntrySizeCfg) != 0 {
		i, err := strconv.Atoi(maxEntrySizeCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse 'max_entry_size': %w", err)
		}

		maxEntrySize = uint64(i)
	}

	return &RaftBackend{
		logger:       logger,
		fsm:          fsm,
		raftInitCh:   make(chan struct{}),
		conf:         conf,
		logStore:     log,
		stableStore:  stable,
		snapStore:    snap,
		dataDir:      path,
		localID:      localID,
		permitPool:   physical.NewPermitPool(physical.DefaultParallelOperations),
		maxEntrySize: maxEntrySize,
	}, nil
}

type snapshotStoreDelay struct {
	wrapped raft.SnapshotStore
	delay   time.Duration
}

func (s snapshotStoreDelay) Create(version raft.SnapshotVersion, index, term uint64, configuration raft.Configuration, configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
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

func newSnapshotStoreDelay(snap raft.SnapshotStore, delay time.Duration) *snapshotStoreDelay {
	return &snapshotStoreDelay{
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

	if err := b.fsm.db.Close(); err != nil {
		return err
	}

	if err := b.stableStore.(*raftboltdb.BoltStore).Close(); err != nil {
		return err
	}

	return nil
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
	ID      string `json:"id"`
	Address string `json:"address"`
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
			ID:      raft.ServerID(p.ID),
			Address: raft.ServerAddress(p.Address),
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
	config.ElectionTimeout = config.ElectionTimeout * time.Duration(multiplier)
	config.HeartbeatTimeout = config.HeartbeatTimeout * time.Duration(multiplier)
	config.LeaderLeaseTimeout = config.LeaderLeaseTimeout * time.Duration(multiplier)

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

	config.NoSnapshotRestoreOnStart = true
	config.MaxAppendEntries = 64

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

	AutopilotConfig *AutopilotConfig
}

func (b *RaftBackend) StartRecoveryCluster(ctx context.Context, peer Peer) error {
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
		// Set the local address and localID in the streaming layer and the raft config.
		streamLayer, err := NewRaftLayer(b.logger.Named("stream"), opts.TLSKeyring, opts.ClusterListener)
		if err != nil {
			return err
		}
		transConfig := &raft.NetworkTransportConfig{
			Stream:                streamLayer,
			MaxPool:               3,
			Timeout:               10 * time.Second,
			ServerAddressProvider: b.serverAddressProvider,
			Logger:                b.logger.Named("raft-net"),
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
			return errwrap.Wrapf("raft recovery failed to parse peers.json: {{err}}", err)
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
			return errwrap.Wrapf("raft recovery failed: {{err}}", err)
		}

		err = os.Remove(peersFile)
		if err != nil {
			return errwrap.Wrapf("raft recovery failed to delete peers.json; please delete manually: {{err}}", err)
		}
		b.logger.Info("raft recovery deleted peers.json")
	}

	if opts.RecoveryModeConfig != nil {
		err = raft.RecoverCluster(raftConfig, b.fsm, b.logStore, b.stableStore, b.snapStore, b.raftTransport, *opts.RecoveryModeConfig)
		if err != nil {
			return errwrap.Wrapf("recovering raft cluster failed: {{err}}", err)
		}
	}

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
		for {
			if raftObj.State() == raft.Leader {
				break
			}

			select {
			case <-ctx.Done():
				future := raftObj.Shutdown()
				if future.Error() != nil {
					return errwrap.Wrapf("shutdown while waiting for leadership: {{err}}", future.Error())
				}

				return errors.New("shutdown while waiting for leadership")
			case <-time.After(10 * time.Millisecond):
			}
		}
	}

	b.raft = raftObj
	b.raftNotifyCh = raftNotifyCh

	// If there is a autopilot config in storage that takes precedence
	if opts.AutopilotConfig != nil {
		// Use the config present in storage
		b.logger.Info("setting autopilot configuration retrieved from storage")
		b.autopilotConfig = opts.AutopilotConfig
	}

	// Autopilot config wasn't found in storage
	if b.autopilotConfig == nil {
		// Check if autopilot settings were part of config file
		conf, err := b.autopilotConf()
		if err != nil {
			return err
		}

		switch {
		case conf != nil:
			b.logger.Info("setting autopilot configuration retrieved from config file")
			// Use the config present in config file
			b.autopilotConfig = conf
		default:
			// Autopilot config is not in both storage and config file. This is the case for
			// existing customers who have not yet enabled autopilot. Disable autopilot.
			b.logger.Info("autopilot configuration not found; disabling autopilot")
			b.disableAutopilot = true
		}
	}

	if !b.disableAutopilot {
		// Create the autopilot instance
		b.logger.Info("setting up autopilot", "config", b.autopilotConfig)
		b.autopilot = autopilot.New(b.raft, &Delegate{b}, autopilot.WithLogger(b.logger), autopilot.WithPromoter(autopilot.DefaultPromoter()))
		b.followerHeartbeatTrackerStopCh = make(chan struct{})
	}

	if b.streamLayer != nil {
		// Add Handler to the cluster.
		opts.ClusterListener.AddHandler(consts.RaftStorageALPN, b.streamLayer)

		// Add Client to the cluster.
		opts.ClusterListener.AddClient(consts.RaftStorageALPN, b.streamLayer)
	}

	// Close the init channel to signal setup has been completed
	close(b.raftInitCh)

	b.logger.Trace("finished setting up raft cluster")
	return nil
}

func (b *RaftBackend) startFollowerHeartbeatTracker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-b.followerHeartbeatTrackerStopCh:
			return
		case <-ticker.C:
			for id, state := range b.followerStates.followers {
				// TODO: Plumb last_contact_failure_threshold setting here
				if !state.LastHeartbeat.IsZero() && time.Now().After(state.LastHeartbeat.Add(5*time.Second)) {
					b.followerStates.MarkFollowerAsDead(id)
				}
			}
		}
	}
}

// StartAutopilot puts autopilot subsystem to work. This should only be called
// on the active node.
func (b *RaftBackend) StartAutopilot(ctx context.Context) {
	if b.autopilot == nil {
		return
	}
	b.autopilot.Start(ctx)

	go b.startFollowerHeartbeatTracker()
}

// StopAutopilot stops a running autopilot instance. This should only be called
// on the active node.
func (b *RaftBackend) StopAutopilot() {
	if b.autopilot == nil {
		return
	}
	b.autopilot.Stop()
	close(b.followerHeartbeatTrackerStopCh)
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

	if b.raft == nil {
		return errors.New("raft storage is not initialized")
	}

	switch b.disableAutopilot {
	case true:
		future := b.raft.RemoveServer(raft.ServerID(peerID), 0, 0)
		return future.Error()
	default:
		return b.autopilot.RemoveServer(raft.ServerID(peerID))
	}
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

// AutopilotHealth represents the health information retrieved from autopilot.
type AutopilotHealth struct {
	Healthy                    bool `json:"healthy"`
	FailureTolerance           int  `json:"failure_tolerance"`
	OptimisticFailureTolerance int  `json:"optimistic_failure_tolerance"`

	Servers      map[string]AutopilotServer `json:"servers"`
	Leader       string                     `json:"leader"`
	Voters       []string                   `json:"voters"`
	ReadReplicas []string                   `json:"read_replicas,omitempty"`
}

// AutopilotServer represents the health information of individual server node
// retrieved from autopilot.
type AutopilotServer struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Address     string            `json:"address"`
	NodeStatus  string            `json:"node_status"`
	LastContact *ReadableDuration `json:"last_contact"`
	LastTerm    uint64            `json:"last_term"`
	LastIndex   uint64            `json:"last_index"`
	Healthy     bool              `json:"healthy"`
	StableSince time.Time         `json:"stable_since"`
	ReadReplica bool              `json:"read_replica"`
	Status      string            `json:"status"`
	Meta        map[string]string `json:"meta"`
	NodeType    string            `json:"node_type"`
}

// ReadableDuration is a duration type that is serialized to JSON in human readable format.
type ReadableDuration time.Duration

func NewReadableDuration(dur time.Duration) *ReadableDuration {
	d := ReadableDuration(dur)
	return &d
}

func (d *ReadableDuration) String() string {
	return d.Duration().String()
}

func (d *ReadableDuration) Duration() time.Duration {
	if d == nil {
		return time.Duration(0)
	}
	return time.Duration(*d)
}

func (d *ReadableDuration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.Duration().String())), nil
}

func (d *ReadableDuration) UnmarshalJSON(raw []byte) (err error) {
	if d == nil {
		return fmt.Errorf("cannot unmarshal to nil pointer")
	}

	var dur time.Duration
	str := string(raw)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		// quoted string
		dur, err = time.ParseDuration(str[1 : len(str)-1])
		if err != nil {
			return err
		}
	} else {
		// no quotes, not a string
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		dur = time.Duration(v)
	}

	*d = ReadableDuration(dur)
	return nil
}

func stringIDs(ids []raft.ServerID) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = string(id)
	}
	return out
}

func autopilotToAPIHealth(state *autopilot.State) *AutopilotHealth {
	out := &AutopilotHealth{
		Healthy:          state.Healthy,
		FailureTolerance: state.FailureTolerance,
		Leader:           string(state.Leader),
		Voters:           stringIDs(state.Voters),
		Servers:          make(map[string]AutopilotServer),
	}

	for id, srv := range state.Servers {
		out.Servers[string(id)] = autopilotToAPIServer(srv)
	}

	return out
}

func autopilotToAPIServer(srv *autopilot.ServerState) AutopilotServer {
	return AutopilotServer{
		ID:          string(srv.Server.ID),
		Name:        srv.Server.Name,
		Address:     string(srv.Server.Address),
		NodeStatus:  string(srv.Server.NodeStatus),
		LastContact: NewReadableDuration(srv.Stats.LastContact),
		LastTerm:    srv.Stats.LastTerm,
		LastIndex:   srv.Stats.LastIndex,
		Healthy:     srv.Health.Healthy,
		StableSince: srv.Health.StableSince,
		Status:      string(srv.State),
		Meta:        srv.Server.Meta,
		NodeType:    string(srv.Server.NodeType),
	}
}

// GetAutopilotServerHealth retrieves raft cluster health from autopilot to
// return over the API.
func (b *RaftBackend) GetAutopilotServerHealth(ctx context.Context) (*AutopilotHealth, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return nil, errors.New("raft storage is not initialized")
	}

	if b.autopilot == nil {
		return nil, nil
	}

	return autopilotToAPIHealth(b.autopilot.GetState()), nil
}

// AddPeer adds a new server to the raft cluster
func (b *RaftBackend) AddPeer(ctx context.Context, peerID, clusterAddr string) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage is not initialized")
	}

	b.logger.Debug("adding raft peer", "node_id", peerID, "cluster_addr", clusterAddr)

	switch b.disableAutopilot {
	case true:
		future := b.raft.AddVoter(raft.ServerID(peerID), raft.ServerAddress(clusterAddr), 0, 0)
		return future.Error()
	default:
		return b.autopilot.AddServer(&autopilot.Server{
			ID:          raft.ServerID(peerID),
			Name:        peerID,
			Address:     raft.ServerAddress(clusterAddr),
			RaftVersion: raft.ProtocolVersionMax,
			NodeType:    autopilot.NodeVoter,
		})
	}
}

// Peers returns all the servers present in the raft cluster
func (b *RaftBackend) Peers(ctx context.Context) ([]Peer, error) {
	b.l.RLock()
	defer b.l.RUnlock()

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

// SnapshotHTTP is a wrapper for Snapshot that sends the snapshot as an HTTP
// response.
func (b *RaftBackend) SnapshotHTTP(out *logical.HTTPResponseWriter, access *seal.Access) error {
	out.Header().Add("Content-Disposition", "attachment")
	out.Header().Add("Content-Type", "application/gzip")

	return b.Snapshot(out, access)
}

// Snapshot takes a raft snapshot, packages it into a archive file and writes it
// to the provided writer. Seal access is used to encrypt the SHASUM file so we
// can validate the snapshot was taken using the same master keys or not.
func (b *RaftBackend) Snapshot(out io.Writer, access *seal.Access) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage backend is sealed")
	}

	// If we have access to the seal create a sealer object
	var s snapshot.Sealer
	if access != nil {
		s = &sealer{
			access: access,
		}
	}

	return snapshot.Write(b.logger.Named("snapshot"), b.raft, s, out)
}

// WriteSnapshotToTemp reads a snapshot archive off the provided reader,
// extracts the data and writes the snapshot to a temporary file. The seal
// access is used to decrypt the SHASUM file in the archive to ensure this
// snapshot has the same master key as the running instance. If the provided
// access is nil then it will skip that validation.
func (b *RaftBackend) WriteSnapshotToTemp(in io.ReadCloser, access *seal.Access) (*os.File, func(), raft.SnapshotMeta, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	var metadata raft.SnapshotMeta
	if b.raft == nil {
		return nil, nil, metadata, errors.New("raft storage backend is sealed")
	}

	// If we have access to the seal create a sealer object
	var s snapshot.Sealer
	if access != nil {
		s = &sealer{
			access: access,
		}
	}

	snap, cleanup, err := snapshot.WriteToTempFileWithSealer(b.logger.Named("snapshot"), in, &metadata, s)
	return snap, cleanup, metadata, err
}

// RestoreSnapshot applies the provided snapshot metadata and snapshot data to
// raft.
func (b *RaftBackend) RestoreSnapshot(ctx context.Context, metadata raft.SnapshotMeta, snap io.Reader) error {
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
			&LogOperation{
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
	command := &LogData{
		Operations: []*LogOperation{
			&LogOperation{
				OpType: deleteOp,
				Key:    path,
			},
		},
	}
	b.permitPool.Acquire()
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

	b.permitPool.Acquire()
	defer b.permitPool.Release()

	entry, err := b.fsm.Get(ctx, path)
	if entry != nil {
		valueLen := len(entry.Value)
		if uint64(valueLen) > b.maxEntrySize {
			b.logger.Warn(
				"retrieved entry value is too large; see https://www.vaultproject.io/docs/configuration/storage/raft#raft-parameters",
				"size", valueLen, "suggested", b.maxEntrySize,
			)
		}
	}

	return entry, err
}

// Put inserts an entry in the log for the put operation. It will return an
// error if the resulting entry encoding exceeds the configured max_entry_size
// or if the call to applyLog fails.
func (b *RaftBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"raft-storage", "put"}, time.Now())
	command := &LogData{
		Operations: []*LogOperation{
			&LogOperation{
				OpType: putOp,
				Key:    entry.Key,
				Value:  entry.Value,
			},
		},
	}

	b.permitPool.Acquire()
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

	b.permitPool.Acquire()
	defer b.permitPool.Release()

	return b.fsm.List(ctx, prefix)
}

// Transaction applies all the given operations into a single log and
// applies it.
func (b *RaftBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	defer metrics.MeasureSince([]string{"raft-storage", "transaction"}, time.Now())
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

	b.permitPool.Acquire()
	defer b.permitPool.Release()

	b.l.RLock()
	err := b.applyLog(ctx, command)
	b.l.RUnlock()
	return err
}

// applyLog will take a given log command and apply it to the raft log. applyLog
// doesn't return until the log has been applied to a quorum of servers and is
// persisted to the local FSM. Caller should hold the backend's read lock.
func (b *RaftBackend) applyLog(ctx context.Context, command *LogData) error {
	if b.raft == nil {
		return errors.New("raft storage backend is not initialized")
	}

	commandBytes, err := proto.Marshal(command)
	if err != nil {
		return err
	}

	cmdSize := len(commandBytes)
	if uint64(cmdSize) > b.maxEntrySize {
		return fmt.Errorf("%s; got %d bytes, max: %d bytes", physical.ErrValueTooLarge, cmdSize, b.maxEntrySize)
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

	if resp, ok := resp.(*FSMApplyResponse); !ok || !resp.Success {
		return errors.New("could not apply data")
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
				&LogOperation{
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
						&LogOperation{
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

	return nil, nil
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

// sealer implements the snapshot.Sealer interface and is used in the snapshot
// process for encrypting/decrypting the SHASUM file in snapshot archives.
type sealer struct {
	access *seal.Access
}

// Seal encrypts the data with using the seal access object.
func (s sealer) Seal(ctx context.Context, pt []byte) ([]byte, error) {
	if s.access == nil {
		return nil, errors.New("no seal access available")
	}
	eblob, err := s.access.Encrypt(ctx, pt, nil)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(eblob)
}

// Open decrypts the data using the seal access object.
func (s sealer) Open(ctx context.Context, ct []byte) ([]byte, error) {
	if s.access == nil {
		return nil, errors.New("no seal access available")
	}

	var eblob wrapping.EncryptedBlobInfo
	err := proto.Unmarshal(ct, &eblob)
	if err != nil {
		return nil, err
	}

	return s.access.Decrypt(ctx, &eblob, nil)
}
