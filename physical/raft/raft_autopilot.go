package raft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"go.uber.org/atomic"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/raft"
	autopilot "github.com/hashicorp/raft-autopilot"
	"github.com/mitchellh/mapstructure"
)

type CleanupDeadServersValue int

const (
	CleanupDeadServersUnset CleanupDeadServersValue = 0
	CleanupDeadServersTrue  CleanupDeadServersValue = 1
	CleanupDeadServersFalse CleanupDeadServersValue = 2
)

func (c CleanupDeadServersValue) Value() bool {
	switch c {
	case CleanupDeadServersTrue:
		return true
	default:
		return false
	}
}

// AutopilotConfig is used for querying/setting the Autopilot configuration.
type AutopilotConfig struct {
	// CleanupDeadServers controls whether to remove dead servers from the Raft
	// peer list periodically or when a new server joins.
	CleanupDeadServers bool `mapstructure:"cleanup_dead_servers"`

	// CleanupDeadServersValue is used to shadow the CleanupDeadServers field in
	// storage. Having it as an int helps in knowing if the value was set explicitly
	// using the API or not.
	CleanupDeadServersValue CleanupDeadServersValue `mapstructure:"cleanup_dead_servers_value"`

	// LastContactThreshold is the limit on the amount of time a server can go
	// without leader contact before being considered unhealthy.
	LastContactThreshold time.Duration `mapstructure:"-"`

	// DeadServerLastContactThreshold is the limit on the amount of time a server
	// can go without leader contact before being considered failed. This takes
	// effect only when CleanupDeadServers is set.
	DeadServerLastContactThreshold time.Duration `mapstructure:"-"`

	// MaxTrailingLogs is the amount of entries in the Raft Log that a server can
	// be behind before being considered unhealthy.
	MaxTrailingLogs uint64 `mapstructure:"max_trailing_logs"`

	// MinQuorum sets the minimum number of servers allowed in a cluster before
	// autopilot can prune dead servers.
	MinQuorum uint `mapstructure:"min_quorum"`

	// ServerStabilizationTime is the minimum amount of time a server must be in a
	// stable, healthy state before it can be added to the cluster. Only applicable
	// with Raft protocol version 3 or higher.
	ServerStabilizationTime time.Duration `mapstructure:"-"`
}

// Merge combines the supplied config with the receiver. Supplied ones take
// priority.
func (to *AutopilotConfig) Merge(from *AutopilotConfig) {
	if from == nil {
		return
	}
	if from.CleanupDeadServersValue != CleanupDeadServersUnset {
		to.CleanupDeadServers = from.CleanupDeadServersValue.Value()
	}
	if from.MinQuorum != 0 {
		to.MinQuorum = from.MinQuorum
	}
	if from.LastContactThreshold != 0 {
		to.LastContactThreshold = from.LastContactThreshold
	}
	if from.DeadServerLastContactThreshold != 0 {
		to.DeadServerLastContactThreshold = from.DeadServerLastContactThreshold
	}
	if from.MaxTrailingLogs != 0 {
		to.MaxTrailingLogs = from.MaxTrailingLogs
	}
	if from.ServerStabilizationTime != 0 {
		to.ServerStabilizationTime = from.ServerStabilizationTime
	}
}

// Clone returns a duplicate instance of AutopilotConfig with the exact same values.
func (ac *AutopilotConfig) Clone() *AutopilotConfig {
	if ac == nil {
		return nil
	}
	return &AutopilotConfig{
		CleanupDeadServers:             ac.CleanupDeadServers,
		LastContactThreshold:           ac.LastContactThreshold,
		DeadServerLastContactThreshold: ac.DeadServerLastContactThreshold,
		MaxTrailingLogs:                ac.MaxTrailingLogs,
		MinQuorum:                      ac.MinQuorum,
		ServerStabilizationTime:        ac.ServerStabilizationTime,
	}
}

// MarshalJSON makes the autopilot config fields JSON compatible
func (ac *AutopilotConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"cleanup_dead_servers":               ac.CleanupDeadServers,
		"cleanup_dead_servers_value":         ac.CleanupDeadServersValue,
		"last_contact_threshold":             ac.LastContactThreshold.String(),
		"dead_server_last_contact_threshold": ac.DeadServerLastContactThreshold.String(),
		"max_trailing_logs":                  ac.MaxTrailingLogs,
		"min_quorum":                         ac.MinQuorum,
		"server_stabilization_time":          ac.ServerStabilizationTime.String(),
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
	if ac.DeadServerLastContactThreshold, err = parseutil.ParseDurationSecond(conf["dead_server_last_contact_threshold"]); err != nil {
		return err
	}
	if ac.ServerStabilizationTime, err = parseutil.ParseDurationSecond(conf["server_stabilization_time"]); err != nil {
		return err
	}

	return nil
}

// FollowerState represents the information about peer that the leader tracks.
type FollowerState struct {
	AppliedIndex    uint64
	LastHeartbeat   time.Time
	LastTerm        uint64
	IsDead          *atomic.Bool
	DesiredSuffrage string
}

// FollowerStates holds information about all the followers in the raft cluster
// tracked by the leader.
type FollowerStates struct {
	l         sync.RWMutex
	followers map[string]*FollowerState
}

// NewFollowerStates creates a new FollowerStates object
func NewFollowerStates() *FollowerStates {
	return &FollowerStates{
		followers: make(map[string]*FollowerState),
	}
}

// Update the peer information in the follower states
func (s *FollowerStates) Update(nodeID string, appliedIndex uint64, term uint64, desiredSuffrage string) {
	s.l.Lock()
	defer s.l.Unlock()

	state, ok := s.followers[nodeID]
	if !ok {
		state = &FollowerState{
			IsDead: atomic.NewBool(false),
		}
		s.followers[nodeID] = state
	}

	state.IsDead.Store(false)
	state.AppliedIndex = appliedIndex
	state.LastTerm = term
	state.DesiredSuffrage = desiredSuffrage
	state.LastHeartbeat = time.Now()
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

	// dl is a lock dedicated for guarding delegate's fields
	dl               sync.RWMutex
	inflightRemovals map[raft.ServerID]bool
}

func newDelegate(b *RaftBackend) *Delegate {
	return &Delegate{
		RaftBackend:      b,
		inflightRemovals: make(map[raft.ServerID]bool),
	}
}

// AutopilotConfig is called by the autopilot library to know the desired
// autopilot configuration.
func (d *Delegate) AutopilotConfig() *autopilot.Config {
	d.l.RLock()
	config := &autopilot.Config{
		CleanupDeadServers:      d.autopilotConfig.CleanupDeadServers,
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

		for id, state := range state.Servers {
			labels := []metrics.Label{
				{
					Name:  "node_id",
					Value: string(id),
				},
			}
			if state.Health.Healthy {
				metrics.SetGaugeWithLabels([]string{"autopilot", "node", "healthy"}, 1, labels)
			} else {
				metrics.SetGaugeWithLabels([]string{"autopilot", "node", "healthy"}, 0, labels)
			}
		}
	}
}

// FetchServerStats is called by the autopilot library to retrieve information
// about all the nodes in the raft cluster.
func (d *Delegate) FetchServerStats(ctx context.Context, servers map[raft.ServerID]*autopilot.Server) map[raft.ServerID]*autopilot.ServerStats {
	ret := make(map[raft.ServerID]*autopilot.ServerStats)

	d.l.RLock()
	followerStates := d.followerStates
	d.l.RUnlock()

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
	defer d.l.RUnlock()
	future := d.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		d.logger.Error("failed to get raft configuration when computing known servers", "error", err)
		return nil
	}

	servers := future.Configuration().Servers
	serverIDs := make([]string, 0, len(servers))
	for _, server := range servers {
		serverIDs = append(serverIDs, string(server.ID))
	}

	d.followerStates.l.RLock()
	defer d.followerStates.l.RUnlock()

	ret := make(map[raft.ServerID]*autopilot.Server)
	for id, state := range d.followerStates.followers {
		// If the server is not in raft configuration, even if we received a follower
		// heartbeat, it shouldn't be a known server for autopilot.
		if !strutil.StrListContains(serverIDs, id) {
			continue
		}

		server := &autopilot.Server{
			ID:          raft.ServerID(id),
			Name:        id,
			RaftVersion: raft.ProtocolVersionMax,
			Ext:         d.autopilotServerExt(state.DesiredSuffrage),
		}

		switch state.IsDead.Load() {
		case true:
			d.logger.Debug("informing autopilot that the node left", "id", id)
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
		Ext:         d.autopilotServerExt("voter"),
		IsLeader:    true,
	}

	return ret
}

// RemoveFailedServer is called by the autopilot library when it desires a node
// to be removed from the raft configuration. This function removes the node
// from the raft cluster and stops tracking its information in follower states.
// This function needs to return quickly. Hence removal is performed in a
// goroutine.
func (d *Delegate) RemoveFailedServer(server *autopilot.Server) {
	go func() {
		added := false
		defer func() {
			if added {
				d.dl.Lock()
				delete(d.inflightRemovals, server.ID)
				d.dl.Unlock()
			}
		}()

		d.dl.Lock()
		_, ok := d.inflightRemovals[server.ID]
		if ok {
			d.logger.Info("removal of dead server is already initiated", "id", server.ID)
			d.dl.Unlock()
			return
		}

		added = true
		d.inflightRemovals[server.ID] = true
		d.dl.Unlock()

		d.logger.Info("removing dead server from raft configuration", "id", server.ID)
		if future := d.raft.RemoveServer(server.ID, 0, 0); future.Error() != nil {
			d.logger.Error("failed to remove server", "server_id", server.ID, "server_address", server.Address, "server_name", server.Name, "error", future.Error())
			return
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
	b.logger.Info("updated autopilot configuration", "config", b.autopilotConfig)
	b.l.Unlock()
}

// AutopilotConfig returns the autopilot configuration in the backend.
func (b *RaftBackend) AutopilotConfig() *AutopilotConfig {
	b.l.RLock()
	defer b.l.RUnlock()
	return b.autopilotConfig.Clone()
}

func (b *RaftBackend) defaultAutopilotConfig() *AutopilotConfig {
	return &AutopilotConfig{
		CleanupDeadServers:             false,
		LastContactThreshold:           10 * time.Second,
		DeadServerLastContactThreshold: 24 * time.Hour,
		MaxTrailingLogs:                1000,
		ServerStabilizationTime:        10 * time.Second,
	}
}

func (b *RaftBackend) AutopilotDisabled() bool {
	b.l.RLock()
	disabled := b.disableAutopilot
	b.l.RUnlock()
	return disabled
}

func (b *RaftBackend) startFollowerHeartbeatTracker() {
	b.l.RLock()
	tickerCh := b.followerHeartbeatTicker.C
	b.l.RUnlock()

	for range tickerCh {
		b.l.RLock()
		if b.autopilotConfig.CleanupDeadServers && b.autopilotConfig.DeadServerLastContactThreshold != 0 {
			b.followerStates.l.RLock()
			for _, state := range b.followerStates.followers {
				if state.LastHeartbeat.IsZero() || state.IsDead.Load() {
					continue
				}
				now := time.Now()
				if now.After(state.LastHeartbeat.Add(b.autopilotConfig.DeadServerLastContactThreshold)) {
					state.IsDead.Store(true)
				}
			}
			b.followerStates.l.RUnlock()
		}
		b.l.RUnlock()
	}
}

// StopAutopilot stops a running autopilot instance. This should only be called
// on the active node.
func (b *RaftBackend) StopAutopilot() {
	b.l.Lock()
	defer b.l.Unlock()

	if b.autopilot == nil {
		return
	}
	b.autopilot.Stop()
	b.followerHeartbeatTicker.Stop()
}

// AutopilotState represents the health information retrieved from autopilot.
type AutopilotState struct {
	Healthy          bool `json:"healthy"`
	FailureTolerance int  `json:"failure_tolerance"`

	Servers   map[string]*AutopilotServer `json:"servers"`
	Leader    string                      `json:"leader"`
	Voters    []string                    `json:"voters"`
	NonVoters []string                    `json:"non_voters,omitempty"`
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

func autopilotToAPIState(state *autopilot.State) (*AutopilotState, error) {
	out := &AutopilotState{
		Healthy:          state.Healthy,
		FailureTolerance: state.FailureTolerance,
		Leader:           string(state.Leader),
		Voters:           stringIDs(state.Voters),
		Servers:          make(map[string]*AutopilotServer),
	}

	for id, srv := range state.Servers {
		out.Servers[string(id)] = autopilotToAPIServer(srv)
	}

	return out, nil
}

func autopilotToAPIServer(srv *autopilot.ServerState) *AutopilotServer {
	apiSrv := &AutopilotServer{
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

	autopilotToAPIServerEnterprise(srv, apiSrv)

	return apiSrv
}

// GetAutopilotServerState retrieves raft cluster state from autopilot to
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

	apState := b.autopilot.GetState()
	if apState == nil {
		return nil, nil
	}

	return autopilotToAPIState(apState)
}

func (b *RaftBackend) DisableAutopilot() {
	b.l.Lock()
	b.disableAutopilot = true
	b.l.Unlock()
}

// SetupAutopilot gathers information required to configure autopilot and starts
// it. If autopilot is disabled, this function does nothing.
func (b *RaftBackend) SetupAutopilot(ctx context.Context, storageConfig *AutopilotConfig, followerStates *FollowerStates, disable bool) {
	b.l.Lock()
	if disable || os.Getenv("VAULT_RAFT_AUTOPILOT_DISABLE") != "" {
		b.disableAutopilot = true
	}

	if b.disableAutopilot {
		b.logger.Info("disabling autopilot")
		b.l.Unlock()
		return
	}

	// Start with a default config
	b.autopilotConfig = b.defaultAutopilotConfig()

	// Merge the setting provided over the API
	b.autopilotConfig.Merge(storageConfig)

	// Create the autopilot instance
	options := []autopilot.Option{
		autopilot.WithLogger(b.logger),
		autopilot.WithPromoter(b.autopilotPromoter()),
	}
	if b.autopilotReconcileInterval != 0 {
		options = append(options, autopilot.WithReconcileInterval(b.autopilotReconcileInterval))
	}
	b.autopilot = autopilot.New(b.raft, newDelegate(b), options...)
	b.followerStates = followerStates
	b.followerHeartbeatTicker = time.NewTicker(1 * time.Second)

	b.l.Unlock()

	b.logger.Info("starting autopilot", "config", b.autopilotConfig, "reconcile_interval", b.autopilotReconcileInterval)
	b.autopilot.Start(ctx)

	go b.startFollowerHeartbeatTracker()
}
