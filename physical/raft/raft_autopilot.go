package raft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/helper/strutil"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/raft"
	autopilot "github.com/hashicorp/raft-autopilot"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"

	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/mitchellh/mapstructure"
)

// AutopilotConfig is used for querying/setting the Autopilot configuration.
// Autopilot helps manage operator tasks related to Vault servers like removing
// failed servers from the Raft quorum.
type AutopilotConfig struct {
	// CleanupDeadServers controls whether to remove dead servers from the Raft
	// peer list periodically or when a new server joins
	CleanupDeadServers bool `mapstructure:"cleanup_dead_servers"`

	// LastContactThreshold is the limit on the amount of time a server can go
	// without leader contact before being considered unhealthy.
	LastContactThreshold time.Duration `mapstructure:"-"`

	// MaxTrailingLogs is the amount of entries in the Raft Log that a server can
	// be behind before being considered unhealthy.
	MaxTrailingLogs uint64 `mapstructure:"max_trailing_logs"`

	// MinQuorum sets the minimum number of servers allowed in a cluster before
	// autopilot can prune dead servers.
	MinQuorum uint `mapstructure:"min_quorum"`

	// ServerStabilizationTime is the minimum amount of time a server must be
	// in a stable, healthy state before it can be added to the cluster. Only
	// applicable with Raft protocol version 3 or higher.
	ServerStabilizationTime time.Duration `mapstructure:"-"`
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

// FollowerState represents the information about peer that the leader tracks.
type FollowerState struct {
	AppliedIndex  uint64
	LastHeartbeat time.Time
	LastTerm      uint64
	IsDead        bool
	NonVoter      bool
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
func (s *FollowerStates) Update(nodeID string, appliedIndex uint64, term uint64, nonVoter bool) {
	state := FollowerState{
		AppliedIndex: appliedIndex,
		LastTerm:     term,
		NonVoter:     nonVoter,
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

// Get retrieves information about a peer from the follower states.
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
	config := &autopilot.Config{
		CleanupDeadServers:      false,
		LastContactThreshold:    d.autopilotConfig.LastContactThreshold,
		MaxTrailingLogs:         d.autopilotConfig.MaxTrailingLogs,
		MinQuorum:               d.autopilotConfig.MinQuorum,
		ServerStabilizationTime: d.autopilotConfig.ServerStabilizationTime,
		Ext:                     d.autopilotConfigExt(),
	}
	d.l.RUnlock()
	return config
}

// NotifyState is called by the autopilot library whenever there is a state
// change. We update a few metrics when this happens.
func (d *Delegate) NotifyState(state *autopilot.State) {
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
	d.l.RLock()
	future := d.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		d.logger.Error("failed to get raft configuration when computing known servers", "error", err)
	}

	servers := future.Configuration().Servers
	serverIDs := make([]string, 0, len(servers))
	for _, server := range servers {
		serverIDs = append(serverIDs, string(server.ID))
	}
	d.l.RUnlock()

	followerStates := d.RaftBackend.followerStates
	followerStates.l.RLock()
	defer followerStates.l.RUnlock()

	ret := make(map[raft.ServerID]*autopilot.Server)
	for id, state := range d.RaftBackend.followerStates.followers {
		// If the server is not in raft configuration, even if we received a follower
		// heartbeat, it shouldn't be a known server for autopilot.
		if !strutil.StrListContains(serverIDs, id) {
			continue
		}

		server := &autopilot.Server{
			ID:          raft.ServerID(id),
			Name:        id,
			RaftVersion: raft.ProtocolVersionMax,
			Ext:         d.autopilotServerExt(state.NonVoter),
		}

		switch state.IsDead {
		case true:
			d.logger.Info("informing autopilot that the node left", "id", id)
			server.NodeStatus = autopilot.NodeLeft
		default:
			server.NodeStatus = autopilot.NodeAlive
		}

		ret[raft.ServerID(id)] = server
	}

	// Add the leader
	ret[raft.ServerID(d.localID)] = &autopilot.Server{
		ID:          raft.ServerID(d.localID),
		Name:        d.localID,
		RaftVersion: raft.ProtocolVersionMax,
		NodeStatus:  autopilot.NodeAlive,
		Ext:         d.autopilotServerExt(false),
	}

	return ret
}

func (d *Delegate) RemoveFailedServer(server *autopilot.Server) {
	// TODO: implement me
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
		CleanupDeadServers:      false,
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

// StartAutopilot puts autopilot subsystem to work. This should only be called
// on the active node.
func (b *RaftBackend) StartAutopilot(ctx context.Context) {
	if b.autopilot == nil {
		return
	}
	b.autopilot.Start(ctx)
}

// StopAutopilot stops a running autopilot instance. This should only be called
// on the active node.
func (b *RaftBackend) StopAutopilot() {
	if b.autopilot == nil {
		return
	}
	b.autopilot.Stop()
}

// AutopilotState represents the health information retrieved from autopilot.
type AutopilotState struct {
	Healthy                    bool `json:"healthy"`
	FailureTolerance           int  `json:"failure_tolerance"`
	OptimisticFailureTolerance int  `json:"optimistic_failure_tolerance"`

	Servers   map[string]AutopilotServer `json:"servers"`
	Leader    string                     `json:"leader"`
	Voters    []string                   `json:"voters"`
	NonVoters []string                   `json:"non_voters,omitempty"`
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
	Status      string            `json:"status"`
	Meta        map[string]string `json:"meta"`
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

func autopilotToAPIHealth(state *autopilot.State) *AutopilotState {
	out := &AutopilotState{
		Healthy:          state.Healthy,
		FailureTolerance: state.FailureTolerance,
		Leader:           string(state.Leader),
		Voters:           stringIDs(state.Voters),
		Servers:          make(map[string]AutopilotServer),
	}

	for id, srv := range state.Servers {
		out.Servers[string(id)] = autopilotToAPIServer(srv)
	}

	autopilotToAPIStateEnterprise(state, out)

	return out
}

func autopilotToAPIServer(srv *autopilot.ServerState) AutopilotServer {
	apiSrv := AutopilotServer{
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
	}

	autopilotToAPIServerEnterprise(srv, &apiSrv)

	return apiSrv
}

// GetAutopilotServerState retrieves raft cluster health from autopilot to
// return over the API.
func (b *RaftBackend) GetAutopilotServerState(ctx context.Context) (*AutopilotState, error) {
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

func (b *RaftBackend) setupAutopilot(opts SetupOpts) error {
	// Start with a default config
	b.autopilotConfig = b.defaultAutopilotConfig()

	// Check if the config was present in storage
	switch opts.AutopilotConfig {
	case nil:
		// Autopilot config wasn't found in storage. Check if autopilot settings were part of config file.
		conf, err := b.autopilotConf()
		if err != nil {
			return err
		}
		if conf != nil {
			b.logger.Info("setting autopilot configuration retrieved from config file")
			b.autopilotConfig = conf
		}
	default:
		b.logger.Info("setting autopilot configuration retrieved from storage")
		b.autopilotConfig = opts.AutopilotConfig
	}

	// Create the autopilot instance
	b.logger.Info("setting up autopilot", "config", b.autopilotConfig)
	b.autopilot = autopilot.New(b.raft, &Delegate{b}, autopilot.WithLogger(b.logger), autopilot.WithPromoter(b.autopilotPromoter()))

	return nil
}
