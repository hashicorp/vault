// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/raft"
	autopilot "github.com/hashicorp/raft-autopilot"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/atomic"
)

type CleanupDeadServersValue int

const (
	CleanupDeadServersUnset    CleanupDeadServersValue = 0
	CleanupDeadServersTrue     CleanupDeadServersValue = 1
	CleanupDeadServersFalse    CleanupDeadServersValue = 2
	AutopilotUpgradeVersionTag string                  = "upgrade_version"
	AutopilotRedundancyZoneTag string                  = "redundancy_zone"
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

	// (Enterprise-only) DisableUpgradeMigration will disable Autopilot's upgrade migration
	// strategy of waiting until enough newer-versioned servers have been added to the
	// cluster before promoting them to voters.
	DisableUpgradeMigration bool `mapstructure:"disable_upgrade_migration"`

	// (Enterprise-only) RedundancyZoneTag is the node tag to use for separating
	// servers into zones for redundancy. If left blank, this feature will be disabled.
	RedundancyZoneTag string `mapstructure:"redundancy_zone_tag"`

	// (Enterprise-only) UpgradeVersionTag is the node tag to use for version info when
	// performing upgrade migrations. If left blank, the Consul version will be used.
	UpgradeVersionTag string `mapstructure:"upgrade_version_tag"`
}

func (ac *AutopilotConfig) String() string {
	s := "CleanupDeadServers:%t " +
		"LastContactThreshold:%s " +
		"DeadServerLastContactThreshold:%s " +
		"MaxTrailingLogs:%d " +
		"MinQuorum:%d " +
		"ServerStabilizationTime:%s " +
		"DisableUpgradeMigration:%t " +
		"RedundancyZoneTag:%s " +
		"UpgradeVersionTag:%s"
	return fmt.Sprintf(s, ac.CleanupDeadServers,
		ac.LastContactThreshold,
		ac.DeadServerLastContactThreshold,
		ac.MaxTrailingLogs,
		ac.MinQuorum,
		ac.ServerStabilizationTime,
		ac.DisableUpgradeMigration,
		ac.RedundancyZoneTag,
		ac.UpgradeVersionTag)
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

	// UpgradeVersionTag and RedundancyZoneTag are purposely not included here since those values aren't user
	// controllable and should never change.
	to.DisableUpgradeMigration = from.DisableUpgradeMigration
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
		UpgradeVersionTag:              ac.UpgradeVersionTag,
		RedundancyZoneTag:              ac.RedundancyZoneTag,
		DisableUpgradeMigration:        ac.DisableUpgradeMigration,
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
		"upgrade_version_tag":                ac.UpgradeVersionTag,
		"redundancy_zone_tag":                ac.RedundancyZoneTag,
		"disable_upgrade_migration":          ac.DisableUpgradeMigration,
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
	Version         string
	UpgradeVersion  string
	RedundancyZone  string
}

// partialCopy returns a partial copy of the follower state.
// This copy uses the same pointer to the IsDead
// atomic field. We need to do this to ensure that
// an update of the IsDead boolean will still be
// accessible in a copied state.
func (f *FollowerState) partialCopy() *FollowerState {
	return &FollowerState{
		AppliedIndex:    f.AppliedIndex,
		LastHeartbeat:   f.LastHeartbeat,
		LastTerm:        f.LastTerm,
		IsDead:          f.IsDead,
		DesiredSuffrage: f.DesiredSuffrage,
		Version:         f.Version,
		UpgradeVersion:  f.UpgradeVersion,
		RedundancyZone:  f.RedundancyZone,
	}
}

// PersistedFollowerState holds the information that gets persisted to storage
type PersistedFollowerState struct {
	Version        string `json:"version"`
	UpgradeVersion string `json:"upgrade_version"`
}

type PersistedFollowerStates struct {
	l      sync.RWMutex
	States map[string]PersistedFollowerState
}

// shouldUpdate checks if the persisted state contains the same servers as the
// current autopilot state. If grabLock is true, a read lock is acquired before
// accessing the map
func (p *PersistedFollowerStates) shouldUpdate(state *autopilot.State, grabLock bool) bool {
	if grabLock {
		p.l.RLock()
		defer p.l.RUnlock()
	}
	if len(state.Servers) != len(p.States) {
		return true
	}
	for id, server := range state.Servers {
		persistedServer, found := p.States[string(id)]
		if !found {
			return true
		}
		if server.Server.Version != persistedServer.Version ||
			server.Server.Meta[AutopilotUpgradeVersionTag] != persistedServer.UpgradeVersion {
			return true
		}
	}
	return false
}

// updatePersistedState checks if the persisted state matches the current
// autopilot state. If not, the state is replaced and persisted
func (d *Delegate) updatePersistedState(state *autopilot.State) error {
	if !d.persistedState.shouldUpdate(state, true) {
		return nil
	}
	newStates := make(map[string]PersistedFollowerState)
	for id, server := range state.Servers {
		newStates[string(id)] = PersistedFollowerState{
			Version:        server.Server.Version,
			UpgradeVersion: server.Server.Meta[AutopilotUpgradeVersionTag],
		}
	}
	d.persistedState.l.Lock()
	defer d.persistedState.l.Unlock()
	if !d.persistedState.shouldUpdate(state, false) {
		return nil
	}
	d.logger.Debug("updating autopilot persisted state")
	err := d.saveStateFn(newStates)
	if err != nil {
		return err
	}
	d.persistedState.States = newStates
	return nil
}

// EchoRequestUpdate is here to avoid 1) the list of arguments to Update() getting huge 2) an import cycle on the vault package
type EchoRequestUpdate struct {
	NodeID          string
	AppliedIndex    uint64
	Term            uint64
	DesiredSuffrage string
	UpgradeVersion  string
	SDKVersion      string
	RedundancyZone  string
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

// Update the peer information in the follower states. Note that this function
// runs on the active node. Returns true if a new entry was added, as opposed
// to modifying one already present.
func (s *FollowerStates) Update(req *EchoRequestUpdate) bool {
	s.l.Lock()
	defer s.l.Unlock()

	state, present := s.followers[req.NodeID]
	if !present {
		state = &FollowerState{
			IsDead: atomic.NewBool(false),
		}
		s.followers[req.NodeID] = state
	}

	state.IsDead.Store(false)
	state.AppliedIndex = req.AppliedIndex
	state.LastTerm = req.Term
	state.DesiredSuffrage = req.DesiredSuffrage
	state.LastHeartbeat = time.Now()
	state.Version = req.SDKVersion
	state.UpgradeVersion = req.UpgradeVersion
	state.RedundancyZone = req.RedundancyZone

	return !present
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
	emptyVersionLogs map[raft.ServerID]struct{}
	persistedState   *PersistedFollowerStates
	saveStateFn      func(p map[string]PersistedFollowerState) error
}

func NewDelegate(b *RaftBackend, persistedStates map[string]PersistedFollowerState, savePersistedStates func(p map[string]PersistedFollowerState) error) *Delegate {
	return &Delegate{
		RaftBackend:      b,
		inflightRemovals: make(map[raft.ServerID]bool),
		emptyVersionLogs: make(map[raft.ServerID]struct{}),
		persistedState:   &PersistedFollowerStates{States: persistedStates},
		saveStateFn:      savePersistedStates,
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

		// if there is a change in versions or membership, we should update
		// our persisted state
		err := d.updatePersistedState(state)
		if err != nil {
			d.logger.Error("failed to persist autopilot state", "error", err)
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

	apServerStates := d.autopilot.GetState().Servers
	servers := future.Configuration().Servers
	serverIDs := make([]string, 0, len(servers))
	for _, server := range servers {
		serverIDs = append(serverIDs, string(server.ID))
	}

	d.followerStates.l.RLock()
	defer d.followerStates.l.RUnlock()

	d.persistedState.l.RLock()
	defer d.persistedState.l.RUnlock()

	ret := make(map[raft.ServerID]*autopilot.Server)
	for id, state := range d.followerStates.followers {
		// If the server is not in raft configuration, even if we received a follower
		// heartbeat, it shouldn't be a known server for autopilot.
		if !strutil.StrListContains(serverIDs, id) {
			continue
		}

		currentServerID := raft.ServerID(id)
		followerVersion, upgradeVersion := d.determineFollowerVersions(id, state)
		if state.UpgradeVersion != upgradeVersion {
			// we only have a read lock on state, so we can't modify it
			// safely. Instead, copy it to override the upgrade version
			state = state.partialCopy()
			state.UpgradeVersion = upgradeVersion
		}

		server := &autopilot.Server{
			ID:          currentServerID,
			Name:        id,
			RaftVersion: raft.ProtocolVersionMax,
			Meta:        d.meta(state),
			Version:     followerVersion,
			Ext:         d.autopilotServerExt(state),
		}

		// As KnownServers is a delegate called by autopilot let's check if we already
		// had this data in the correct format and use it. If we don't (which sounds a
		// bit sad, unless this ISN'T a voter) then as a fail-safe, let's try what we've
		// done elsewhere in code to check the desired suffrage and manually set NodeType
		// based on whether that's a voter or not. If we don't  do either of these
		// things, NodeType isn't set which means technically it's not a voter.
		// It shouldn't be a voter and end up in this state.
		if apServerState, found := apServerStates[raft.ServerID(id)]; found && apServerState.Server.NodeType != "" {
			server.NodeType = apServerState.Server.NodeType
		} else if state.DesiredSuffrage == "voter" {
			server.NodeType = autopilot.NodeVoter
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
		NodeType:    autopilot.NodeVoter, // The leader must be a voter
		Meta: d.meta(&FollowerState{
			UpgradeVersion: d.UpgradeVersion(),
			RedundancyZone: d.RedundancyZone(),
		}),
		Version:  d.effectiveSDKVersion,
		Ext:      d.autopilotServerExt(nil),
		IsLeader: true,
	}

	return ret
}

// determineFollowerVersions uses the following logic:
//   - if the version and upgrade version are present in the follower state,
//     return those.
//   - if the persisted states map is empty, it means that persisted states
//     don't exist. This happens on an upgrade to 1.18. Use the leader node's
//     versions.
//   - use the versions in the persisted states map
//
// This function must be called with a lock on d.followerStates
// and d.persistedStates.
func (d *Delegate) determineFollowerVersions(id string, state *FollowerState) (version string, upgradeVersion string) {
	// if we have both versions in follower states, use those
	if state.Version != "" && state.UpgradeVersion != "" {
		return state.Version, state.UpgradeVersion
	}

	version = state.Version
	upgradeVersion = state.UpgradeVersion

	// the persistedState map should only be empty on upgrades
	// to 1.18.x. This is the only case where we'll stub with
	// the leader's versions
	if len(d.persistedState.States) == 0 {
		if version == "" {
			version = d.effectiveSDKVersion
			d.logger.Debug("no persisted state, using leader version", "id", id, "version", version)
		}
		if upgradeVersion == "" {
			upgradeVersion = d.upgradeVersion
			d.logger.Debug("no persisted state, using leader upgrade version version", "id", id, "upgrade_version", upgradeVersion)
		}
		return version, upgradeVersion
	}

	// Use the persistedStates map to fill in the sdk
	// and upgrade versions
	pState := d.persistedState.States[id]
	if version == "" {
		version = pState.Version
		d.logger.Debug("using follower version from persisted states", "id", id, "version", version)
	}
	if upgradeVersion == "" {
		upgradeVersion = pState.UpgradeVersion
		d.logger.Debug("using upgrade version from persisted states", "id", id, "upgrade_version", upgradeVersion)
	}
	return version, upgradeVersion
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
		_, err := d.RemovedServerCleanup(context.Background(), string(server.ID))
		if err != nil {
			d.logger.Error("failed to run cleanup", "error", err)
		}
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
		DisableUpgradeMigration:        false,
		UpgradeVersionTag:              AutopilotUpgradeVersionTag,
		RedundancyZoneTag:              AutopilotRedundancyZoneTag,
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

	followerGauge := func(peerID string, suffix string, value float32) {
		labels := []metrics.Label{
			{
				Name:  "peer_id",
				Value: peerID,
			},
		}
		metrics.SetGaugeWithLabels([]string{"raft_storage", "follower", suffix}, value, labels)
	}
	for range tickerCh {
		b.l.RLock()
		if b.raft == nil {
			// We could be racing with teardown, which will stop the ticker
			// but that doesn't guarantee that we won't reach this line with a nil
			// b.raft.
			b.l.RUnlock()
			return
		}
		b.followerStates.l.RLock()
		myAppliedIndex := b.raft.AppliedIndex()
		for peerID, state := range b.followerStates.followers {
			timeSinceLastHeartbeat := time.Now().Sub(state.LastHeartbeat) / time.Millisecond
			followerGauge(peerID, "last_heartbeat_ms", float32(timeSinceLastHeartbeat))
			followerGauge(peerID, "applied_index_delta", float32(myAppliedIndex-state.AppliedIndex))

			if b.autopilotConfig.CleanupDeadServers && b.autopilotConfig.DeadServerLastContactThreshold != 0 {
				if state.LastHeartbeat.IsZero() || state.IsDead.Load() {
					continue
				}
				now := time.Now()
				if now.After(state.LastHeartbeat.Add(b.autopilotConfig.DeadServerLastContactThreshold)) {
					state.IsDead.Store(true)
				}
			}
		}
		b.followerStates.l.RUnlock()
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
	stopCh := b.autopilot.Stop()
	<-stopCh
	b.autopilot = nil
	b.followerHeartbeatTicker.Stop()
}

// AutopilotState represents the health information retrieved from autopilot.
type AutopilotState struct {
	Healthy                    bool                        `json:"healthy" mapstructure:"healthy"`
	FailureTolerance           int                         `json:"failure_tolerance" mapstructure:"failure_tolerance"`
	Servers                    map[string]*AutopilotServer `json:"servers" mapstructure:"servers"`
	Leader                     string                      `json:"leader" mapstructure:"leader"`
	Voters                     []string                    `json:"voters" mapstructure:"voters"`
	NonVoters                  []string                    `json:"non_voters,omitempty" mapstructure:"non_voters,omitempty"`
	RedundancyZones            map[string]AutopilotZone    `json:"redundancy_zones,omitempty" mapstructure:"redundancy_zones,omitempty"`
	Upgrade                    *AutopilotUpgrade           `json:"upgrade_info,omitempty" mapstructure:"upgrade_info,omitempty"`
	OptimisticFailureTolerance int                         `json:"optimistic_failure_tolerance,omitempty" mapstructure:"optimistic_failure_tolerance,omitempty"`
}

// AutopilotServer represents the health information of individual server node
// retrieved from autopilot.
type AutopilotServer struct {
	ID             string            `json:"id" mapstructure:"id"`
	Name           string            `json:"name" mapstructure:"name"`
	Address        string            `json:"address" mapstructure:"address"`
	NodeStatus     string            `json:"node_status" mapstructure:"node_status"`
	LastContact    *ReadableDuration `json:"last_contact" mapstructure:"last_contact"`
	LastTerm       uint64            `json:"last_term" mapstructure:"last_term"`
	LastIndex      uint64            `json:"last_index" mapstructure:"last_index"`
	Healthy        bool              `json:"healthy" mapstructure:"healthy"`
	StableSince    time.Time         `json:"stable_since" mapstructure:"stable_since"`
	Status         string            `json:"status" mapstructure:"status"`
	Version        string            `json:"version" mapstructure:"version"`
	RedundancyZone string            `json:"redundancy_zone,omitempty" mapstructure:"redundancy_zone,omitempty"`
	UpgradeVersion string            `json:"upgrade_version,omitempty" mapstructure:"upgrade_version,omitempty"`
	ReadReplica    bool              `json:"read_replica,omitempty" mapstructure:"read_replica,omitempty"`
	NodeType       string            `json:"node_type,omitempty" mapstructure:"node_type,omitempty"`
}

type AutopilotZone struct {
	Servers          []string `json:"servers,omitempty" mapstructure:"servers,omitempty"`
	Voters           []string `json:"voters,omitempty" mapstructure:"voters,omitempty"`
	FailureTolerance int      `json:"failure_tolerance,omitempty" mapstructure:"failure_tolerance,omitempty"`
}

type AutopilotUpgrade struct {
	Status                    string                                  `json:"status" mapstructure:"status"`
	TargetVersion             string                                  `json:"target_version,omitempty" mapstructure:"target_version,omitempty"`
	TargetVersionVoters       []string                                `json:"target_version_voters,omitempty" mapstructure:"target_version_voters,omitempty"`
	TargetVersionNonVoters    []string                                `json:"target_version_non_voters,omitempty" mapstructure:"target_version_non_voters,omitempty"`
	TargetVersionReadReplicas []string                                `json:"target_version_read_replicas,omitempty" mapstructure:"target_version_read_replicas,omitempty"`
	OtherVersionVoters        []string                                `json:"other_version_voters,omitempty" mapstructure:"other_version_voters,omitempty"`
	OtherVersionNonVoters     []string                                `json:"other_version_non_voters,omitempty" mapstructure:"other_version_non_voters,omitempty"`
	OtherVersionReadReplicas  []string                                `json:"other_version_read_replicas,omitempty" mapstructure:"other_version_read_replicas,omitempty"`
	RedundancyZones           map[string]AutopilotZoneUpgradeVersions `json:"redundancy_zones,omitempty" mapstructure:"redundancy_zones,omitempty"`
}

type AutopilotZoneUpgradeVersions struct {
	TargetVersionVoters    []string `json:"target_version_voters,omitempty" mapstructure:"target_version_voters,omitempty"`
	TargetVersionNonVoters []string `json:"target_version_non_voters,omitempty" mapstructure:"target_version_non_voters,omitempty"`
	OtherVersionVoters     []string `json:"other_version_voters,omitempty" mapstructure:"other_version_voters,omitempty"`
	OtherVersionNonVoters  []string `json:"other_version_non_voters,omitempty" mapstructure:"other_version_non_voters,omitempty"`
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
		dur, err = parseutil.ParseDurationSecond(str[1 : len(str)-1])
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
		aps, err := autopilotToAPIServer(srv)
		if err != nil {
			return nil, err
		}
		out.Servers[string(id)] = aps
	}

	err := autopilotToAPIStateEnterprise(state, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func autopilotToAPIServer(srv *autopilot.ServerState) (*AutopilotServer, error) {
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
		Version:     srv.Server.Version,
		NodeType:    string(srv.Server.NodeType),
	}

	err := autopilotToAPIServerEnterprise(&srv.Server, apiSrv)
	if err != nil {
		return nil, err
	}

	return apiSrv, nil
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

type AutopilotSetupOptions struct {
	StorageConfig       *AutopilotConfig
	FollowerStates      *FollowerStates
	Disable             bool
	PersistedStates     map[string]PersistedFollowerState
	SavePersistedStates func(p map[string]PersistedFollowerState) error
}

// SetupAutopilot gathers information required to configure autopilot and starts
// it. If autopilot is disabled, this function does nothing.
func (b *RaftBackend) SetupAutopilot(ctx context.Context, opts *AutopilotSetupOptions) {
	b.l.Lock()
	if opts.Disable || os.Getenv("VAULT_RAFT_AUTOPILOT_DISABLE") != "" {
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
	b.autopilotConfig.Merge(opts.StorageConfig)

	infoArgs := []interface{}{"config", b.autopilotConfig}

	// Create the autopilot instance
	options := []autopilot.Option{
		autopilot.WithLogger(b.logger),
		autopilot.WithPromoter(b.autopilotPromoter()),
	}
	if b.autopilotReconcileInterval != 0 {
		options = append(options, autopilot.WithReconcileInterval(b.autopilotReconcileInterval))
		infoArgs = append(infoArgs, []interface{}{"reconcile_interval", b.autopilotReconcileInterval}...)
	}
	if b.autopilotUpdateInterval != 0 {
		options = append(options, autopilot.WithUpdateInterval(b.autopilotUpdateInterval))
		infoArgs = append(infoArgs, []interface{}{"update_interval", b.autopilotUpdateInterval}...)
	}
	delegate := NewDelegate(b, opts.PersistedStates, opts.SavePersistedStates)
	b.autopilot = autopilot.New(b.raft, delegate, options...)
	b.followerStates = opts.FollowerStates
	b.followerHeartbeatTicker = time.NewTicker(1 * time.Second)
	b.l.Unlock()

	b.logger.Info("starting autopilot", infoArgs...)
	b.autopilot.Start(ctx)

	go b.startFollowerHeartbeatTracker()
}
