package autopilot

import (
	"context"
	"time"

	"github.com/hashicorp/raft"
)

//go:generate mockery --all --case snake --inpackage

// RaftState is the status of a single server in the Raft cluster.
type RaftState string

const (
	RaftNone     RaftState = "none"
	RaftLeader   RaftState = "leader"
	RaftVoter    RaftState = "voter"
	RaftNonVoter RaftState = "non-voter"
	RaftStaging  RaftState = "staging"
)

func (s RaftState) IsPotentialVoter() bool {
	switch s {
	case RaftVoter, RaftStaging, RaftLeader:
		return true
	default:
		return false
	}
}

// NodeStatus represents the health of a server as know to the autopilot consumer.
// This should not take into account Raft health and the server being on a new enough
// term and index.
type NodeStatus string

const (
	NodeUnknown NodeStatus = "unknown"
	NodeAlive   NodeStatus = "alive"
	NodeFailed  NodeStatus = "failed"
	NodeLeft    NodeStatus = "left"
)

type NodeType string

const (
	NodeVoter NodeType = "voter"
)

// Config represents all the tunables of autopilot
type Config struct {
	// CleanupDeadServers controls whether to remove dead servers when a new
	// server is added to the Raft peers.
	CleanupDeadServers bool

	// LastContactThreshold is the limit on the amount of time a server can go
	// without leader contact before being considered unhealthy.
	LastContactThreshold time.Duration

	// MaxTrailingLogs is the amount of entries in the Raft Log that a server can
	// be behind before being considered unhealthy.
	MaxTrailingLogs uint64

	// MinQuorum sets the minimum number of servers required in a cluster
	// before autopilot can prune dead servers.
	MinQuorum uint

	// ServerStabilizationTime is the minimum amount of time a server must be
	// in a stable, healthy state before it can be added to the cluster. Only
	// applicable with Raft protocol version 3 or higher.
	ServerStabilizationTime time.Duration

	Ext interface{}
}

// Server represents one Raft server
type Server struct {
	// This first set of fields are those that the autopilot delegate
	// needs to fill in

	ID          raft.ServerID
	Name        string
	Address     raft.ServerAddress
	NodeStatus  NodeStatus
	Version     string
	Meta        map[string]string
	RaftVersion int
	IsLeader    bool

	// The remaining fields are those that the promoter
	// will fill in

	NodeType NodeType
	Ext      interface{}
}

type ServerState struct {
	Server Server
	State  RaftState
	Stats  ServerStats
	Health ServerHealth
}

func (s *ServerState) HasVotingRights() bool {
	return s.State == RaftVoter || s.State == RaftLeader
}

// isHealthy determines whether this ServerState is considered healthy
// based on the given Autopilot config
func (s *ServerState) isHealthy(lastTerm uint64, leaderLastIndex uint64, conf *Config) bool {
	// Raft hasn't been bootstrapped yet so nothing is healthy
	if leaderLastIndex == 0 || lastTerm == 0 {
		return false
	}

	// Check that the application still thinks the server is alive and well.
	if s.Server.NodeStatus != NodeAlive {
		return false
	}

	// Check to ensure that the server was contacted recently enough.
	if s.Stats.LastContact > conf.LastContactThreshold || s.Stats.LastContact < 0 {
		return false
	}

	// Check if the server has a different Raft term from the leader
	if s.Stats.LastTerm != lastTerm {
		return false
	}

	// Check if the server has fallen behind more than the configured max trailing logs value
	if s.Stats.LastIndex+conf.MaxTrailingLogs < leaderLastIndex {
		return false
	}

	return true
}

type ServerHealth struct {
	// Healthy is whether the server is healthy according to the current
	// Autopilot config.
	Healthy bool

	// StableSince is the last time this server's Healthy value changed.
	StableSince time.Time
}

// IsStable returns true if the ServerState shows a stable, passing state
// according to the given AutopilotConfig
func (h *ServerHealth) IsStable(now time.Time, minStableDuration time.Duration) bool {
	if h == nil {
		return false
	}

	if !h.Healthy {
		return false
	}

	if now.Sub(h.StableSince) < minStableDuration {
		return false
	}

	return true
}

// ServerStats holds miscellaneous Raft metrics for a server
type ServerStats struct {
	// LastContact is the time since this node's last contact with the leader.
	LastContact time.Duration

	// LastTerm is the highest leader term this server has a record of in its Raft log.
	LastTerm uint64

	// LastIndex is the last log index this server has a record of in its Raft log.
	LastIndex uint64
}

type State struct {
	firstStateTime   time.Time
	Healthy          bool
	FailureTolerance int
	Servers          map[raft.ServerID]*ServerState
	Leader           raft.ServerID
	Voters           []raft.ServerID
	Ext              interface{}
}

func (s *State) ServerStabilizationTime(c *Config) time.Duration {
	// Only use the configured stabilization time when autopilot has
	// been running for at least as long as when the first state was
	// generated. If it hasn't been running that long then we would
	// guarantee that all checks against the stabilization time will
	// fail which will result in excessive leader elections.
	if time.Since(s.firstStateTime) > c.ServerStabilizationTime {
		return c.ServerStabilizationTime
	}

	// ignore stabilization time if autopilot hasn't been running long enough
	// to be tracking any server long enough to meet that requirement
	return 0
}

// Raft is the interface of all the methods on the Raft type that autopilot needs to function. Autopilot will
// take in an interface for Raft instead of a concrete type to allow for dependency injection in tests.
type Raft interface {
	AddNonvoter(id raft.ServerID, address raft.ServerAddress, prevIndex uint64, timeout time.Duration) raft.IndexFuture
	AddVoter(id raft.ServerID, address raft.ServerAddress, prevIndex uint64, timeout time.Duration) raft.IndexFuture
	DemoteVoter(id raft.ServerID, prevIndex uint64, timeout time.Duration) raft.IndexFuture
	LastIndex() uint64
	Leader() raft.ServerAddress
	GetConfiguration() raft.ConfigurationFuture
	RemoveServer(id raft.ServerID, prevIndex uint64, timeout time.Duration) raft.IndexFuture
	Stats() map[string]string
	LeadershipTransferToServer(id raft.ServerID, address raft.ServerAddress) raft.Future
	State() raft.RaftState
}

type ApplicationIntegration interface {
	// AutopilotConfig is used to retrieve the latest configuration from the delegate
	AutopilotConfig() *Config

	// NotifyState will be called when the autopilot state is updated. The application may choose to emit metrics
	// or perform other actions based on this information.
	NotifyState(*State)

	// FetchServerStats will be called to request the application fetch the ServerStats out of band. Usually this
	// will require an RPC to each server.
	FetchServerStats(context.Context, map[raft.ServerID]*Server) map[raft.ServerID]*ServerStats

	// KnownServers fetchs the list of servers as known to the application
	KnownServers() map[raft.ServerID]*Server

	// RemoveFailedServer notifies the application to forcefully remove the server in the failed state
	// It is expected that this returns nearly immediately so if a longer running operation needs to be
	// performed then the Delegate implementation should spawn a go routine itself.
	RemoveFailedServer(*Server)
}

type RaftChanges struct {
	Promotions []raft.ServerID
	Demotions  []raft.ServerID
	Leader     raft.ServerID
}

type FailedServers struct {
	// StaleNonVoters are the ids of those server in the raft configuration as non-voters
	// that are not present in the delegates view of what servers should be available
	StaleNonVoters []raft.ServerID
	// StaleVoters are the ids of those servers in the raft configuration as voters that
	// are not present in the delegates view of what servers should be available
	StaleVoters []raft.ServerID

	// FailedNonVoters are the servers without voting rights in the cluster that the
	// delegate has indicated are in a failed state
	FailedNonVoters []*Server
	// FailedVoters are the servers without voting rights in the cluster that the
	// delegate has indicated are in a failed state
	FailedVoters []*Server
}

func (f *FailedServers) getFailed(ids []raft.ServerID, isVoter bool) []*Server {
	var servers []*Server
	var result []*Server

	if isVoter {
		servers = f.FailedVoters
	} else {
		servers = f.FailedNonVoters
	}

	for _, id := range ids {
		for _, srv := range servers {
			if srv.ID == id {
				result = append(result, srv)
			}
		}
	}

	return result
}

// Promoter is an interface to provide promotion/demotion algorithms to the core autopilot type.
// The BasicPromoter satisfies this interface and will promote any stable servers but other
// algorithms could be implemented. The implementation of these methods shouldn't "block".
// While they are synchronous autopilot expects the algorithms to not make any network
// or other requests which way cause an indefinite amount of waiting to occur.
//
// Note that all parameters passed to these functions should be considered read-only and
// their modification could result in undefined behavior of the core autopilot routines
// including potential crashes.
type Promoter interface {
	// GetServerExt returns some object that should be stored in the Ext field of the Server
	// This value will not be used by the code in this repo but may be used by the other
	// Promoter methods and the application utilizing autopilot. If the value returned is
	// nil the extended state will not be updated.
	GetServerExt(*Config, *ServerState) interface{}

	// GetStateExt returns some object that should be stored in the Ext field of the State
	// This value will not be used by the code in this repo but may be used by the other
	// Promoter methods and the application utilizing autopilot. If the value returned is
	// nil the extended state will not be updated.
	GetStateExt(*Config, *State) interface{}

	// GetNodeTypes returns a map of ServerID to NodeType for all the servers which
	// should have their NodeType field updated
	GetNodeTypes(*Config, *State) map[raft.ServerID]NodeType

	// CalculatePromotionsAndDemotions
	CalculatePromotionsAndDemotions(*Config, *State) RaftChanges

	// FilterFailedServerRemovals takes in the current state and structure outlining all the
	// failed/stale servers and will return those failed servers which the promoter thinks
	// should be allowed to be removed.
	FilterFailedServerRemovals(*Config, *State, *FailedServers) *FailedServers

	// IsPotentialVoter takes a NodeType and returns whether that type represents
	// a potential voter, based on a predicate implemented by the promoter.
	IsPotentialVoter(NodeType) bool
}

// TimeProvider is an interface for getting a local time. This is mainly useful for testing
// to inject certain times so that output validation is easier.
type TimeProvider interface {
	Now() time.Time
}

type runtimeTimeProvider struct{}

func (_ *runtimeTimeProvider) Now() time.Time {
	return time.Now()
}

func (v *voterEligibility) isCurrentVoter() bool {
	return v.currentVoter
}

func (v *voterEligibility) isPotentialVoter() bool {
	return v.potentialVoter
}

func (v *voterEligibility) setPotentialVoter(isVoter bool) {
	v.potentialVoter = isVoter
}

// voterEligibility represents whether a node can currently vote,
// and if it could potentially vote in the future.
type voterEligibility struct {
	currentVoter   bool
	potentialVoter bool
}

type voterRegistry struct {
	eligibility map[raft.ServerID]*voterEligibility
}

func newVoterRegistry() *voterRegistry {
	var result voterRegistry
	result.eligibility = make(map[raft.ServerID]*voterEligibility)
	return &result
}

func (vr *voterRegistry) potentialVoters() int {
	potentialVoters := 0

	for _, v := range vr.eligibility {
		if v.isPotentialVoter() {
			potentialVoters++
		}
	}

	return potentialVoters
}

func (vr *voterRegistry) filter(ids []*Server) []raft.ServerID {
	var result []raft.ServerID

	for _, srv := range ids {
		if _, ok := vr.eligibility[srv.ID]; ok {
			result = append(result, srv.ID)
		}
	}

	return result
}

func (vr *voterRegistry) remove(ids ...raft.ServerID) *voterRegistry {
	for _, id := range ids {
		delete(vr.eligibility, id)
	}

	return vr
}
