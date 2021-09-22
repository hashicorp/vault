package autopilot

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/raft"
)

// aliveServers will filter the input map of servers and output one with all of the
// servers in a Left state removed.
func aliveServers(servers map[raft.ServerID]*Server) map[raft.ServerID]*Server {
	serverMap := make(map[raft.ServerID]*Server)
	for _, server := range servers {
		if server.NodeStatus == NodeLeft {
			continue
		}

		serverMap[server.ID] = server
	}

	return serverMap
}

// nextStateInputs is the collection of values that can influence
// creation of the next State.
type nextStateInputs struct {
	Now          time.Time
	StartTime    time.Time
	Config       *Config
	RaftConfig   *raft.Configuration
	KnownServers map[raft.ServerID]*Server
	LatestIndex  uint64
	LastTerm     uint64
	FetchedStats map[raft.ServerID]*ServerStats
	LeaderID     raft.ServerID
}

// gatherNextStateInputs gathers all the information that would be used to
// create the new updated state from.
//
// - Time Providers current time.
// - Autopilot Config (needed to determine if the stats should indicate unhealthiness)
// - Current state
// - Raft Configuration
// - Known Servers
// - Latest raft index (gathered right before the remote server stats so that they should
//   be from about the same point in time)
// - Stats for all non-left servers
func (a *Autopilot) gatherNextStateInputs(ctx context.Context) (*nextStateInputs, error) {
	// there are a lot of inputs to computing the next state so they get put into a
	// struct so that we don't have to return 8 values.
	inputs := &nextStateInputs{
		Now:       a.time.Now(),
		StartTime: a.startTime,
	}

	// grab the latest autopilot configuration
	config := a.delegate.AutopilotConfig()
	if config == nil {
		return nil, fmt.Errorf("delegate did not return an Autopilot configuration")
	}
	inputs.Config = config

	// retrieve the raft configuration
	raftConfig, err := a.getRaftConfiguration()
	if err != nil {
		return nil, fmt.Errorf("failed to get the Raft configuration: %w", err)
	}
	inputs.RaftConfig = raftConfig

	// get the known servers which may include left/failed ones
	inputs.KnownServers = a.delegate.KnownServers()

	// Try to retrieve leader id from the delegate.
	for id, srv := range inputs.KnownServers {
		if srv.IsLeader {
			inputs.LeaderID = id
			break
		}
	}

	// Delegate setting the leader information is optional. If leader detection is
	// not successful, fallback on raft config to do the same.
	if inputs.LeaderID == "" {
		leader := a.raft.Leader()
		for _, s := range inputs.RaftConfig.Servers {
			if s.Address == leader {
				inputs.LeaderID = s.ID
				break
			}
		}
		if inputs.LeaderID == "" {
			return nil, fmt.Errorf("cannot detect the current leader server id from its address: %s", leader)
		}
	}

	// get the latest Raft index - this should be kept close to the call to
	// fetch the statistics so that the index values are as close in time as
	// possible to make the best decision regarding an individual servers
	// healthiness.
	inputs.LatestIndex = a.raft.LastIndex()

	term, err := a.lastTerm()
	if err != nil {
		return nil, fmt.Errorf("failed to determine the last Raft term: %w", err)
	}
	inputs.LastTerm = term

	// getting the raft configuration could block for a while so now is a good
	// time to check for context cancellation
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// in most cases getting the known servers should be quick but as we cannot
	// account for every potential delegate and prevent them from making
	// blocking network requests we should probably check the context again.
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// we only allow the fetch to take place for up to half the health interval
	// the next health interval will attempt to fetch the stats again but if
	// we do not see responses within this time then we can assume they are
	// unhealthy
	d := inputs.Now.Add(a.updateInterval / 2)
	fetchCtx, cancel := context.WithDeadline(ctx, d)
	defer cancel()

	inputs.FetchedStats = a.delegate.FetchServerStats(fetchCtx, aliveServers(inputs.KnownServers))

	// it might be nil but we propagate the ctx.Err just in case our context was
	// cancelled since the last time we checked.
	return inputs, ctx.Err()
}

// nextState will gather many inputs about the current state of servers from the
// delegate, raft and time provider among other sources and then compute the
// next Autopilot state.
func (a *Autopilot) nextState(ctx context.Context) (*State, error) {
	inputs, err := a.gatherNextStateInputs(ctx)
	if err != nil {
		return nil, err
	}

	state := a.nextStateWithInputs(inputs)
	if state.Leader == "" {
		return nil, fmt.Errorf("Unabled to detect the leader server")
	}
	return state, nil
}

// nextStateWithInputs computes the next state given pre-gathered inputs
func (a *Autopilot) nextStateWithInputs(inputs *nextStateInputs) *State {
	nextServers := a.nextServers(inputs)

	newState := &State{
		startTime: inputs.StartTime,
		Healthy:   true,
		Servers:   nextServers,
	}

	voterCount := 0
	healthyVoters := 0

	// This loop will
	//   1. Determine the ID of the leader server and set it in the state
	//   2. Count the number of voters in the cluster
	//   3. Count the number of healthy voters in the cluster
	//   4. Detect unhealthy servers and mark the overall health as false
	for id, srv := range nextServers {
		if !srv.Health.Healthy {
			// any unhealthiness results in overall unhealthiness
			newState.Healthy = false
		}

		switch srv.State {
		case RaftLeader:
			newState.Leader = id
			fallthrough
		case RaftVoter:
			newState.Voters = append(newState.Voters, id)
			voterCount++

			if srv.Health.Healthy {
				healthyVoters++
			}
		}
	}

	// If we have extra healthy voters, update FailureTolerance from its
	// zero value in the struct.
	requiredQuorum := requiredQuorum(voterCount)
	if healthyVoters > requiredQuorum {
		newState.FailureTolerance = healthyVoters - requiredQuorum
	}

	// update any promoter specific overall state
	if newExt := a.promoter.GetStateExt(inputs.Config, newState); newExt != nil {
		newState.Ext = newExt
	}

	// update the node types - these are really informational for users to
	// know how autopilot and the associate promoter algorithms have classed
	// each server as some promotion algorithms may want to keep certain
	// servers as non-voters for reasons. The node type then can be used
	// to indicate why that might be happening.
	for id, typ := range a.promoter.GetNodeTypes(inputs.Config, newState) {
		if srv, ok := newState.Servers[id]; ok {
			srv.Server.NodeType = typ
		}
	}

	// Sort the voters list to keep the output stable. This is done near the end
	// as SortServers may use other parts of the state that were created in
	// this method and populated in the newState. Requiring output stability
	// helps make tests easier to manage and means that if you happen to be dumping
	// the state periodically you shouldn't see things change unless there
	// are real changes to server health or overall configuration.
	SortServers(newState.Voters, newState)

	return newState
}

// nextServers will build out the servers map for the next state to be created
// from the given inputs. This will take into account all the various sources
// of partial state (current state, raft config, application known servers etc.)
// and combine them into the final server map.
func (a *Autopilot) nextServers(inputs *nextStateInputs) map[raft.ServerID]*ServerState {
	newServers := make(map[raft.ServerID]*ServerState)

	for _, srv := range inputs.RaftConfig.Servers {
		state := a.buildServerState(inputs, srv)

		// update any promoter specific information. This isn't done within
		// buildServerState to keep that function "pure" and not require
		// mocking for tests
		if newExt := a.promoter.GetServerExt(inputs.Config, &state); newExt != nil {
			state.Server.Ext = newExt
		}

		newServers[srv.ID] = &state
	}

	return newServers
}

// buildServerState takes all the nextStateInputs and builds out a ServerState
// for the given Raft server. This will take into account the raft configuration
// existing state, application known servers and recently fetched stats.
func (a *Autopilot) buildServerState(inputs *nextStateInputs, srv raft.Server) ServerState {
	// Note that the ordering of operations in this method are very important.
	// We are building up the ServerState from the least important sources
	// and overriding them with more up to date values.

	// build the basic state from the Raft server
	state := ServerState{
		Server: Server{
			ID:      srv.ID,
			Address: srv.Address,
		},
	}

	switch srv.Suffrage {
	case raft.Voter:
		state.State = RaftVoter
	case raft.Nonvoter:
		state.State = RaftNonVoter
	case raft.Staging:
		state.State = RaftStaging
	default:
		// should be impossible unless the constants in Raft were updated
		// to have a new state.
		// TODO (mkeeler) maybe a panic would be better here. The downside is
		// that it would be hard to catch that in tests when updating the Raft
		// version.
		state.State = RaftNone
	}

	// overwrite the raft state to mark the leader as such instead of just
	// a regular voter
	if srv.ID == inputs.LeaderID {
		state.State = RaftLeader
	}

	var previousHealthy *bool

	a.stateLock.RLock()
	// copy some state from an existing server into the new state - most of this
	// should be overridden soon but at this point we are just building the base.
	if existing, found := a.state.Servers[srv.ID]; found {
		state.Stats = existing.Stats
		state.Health = existing.Health
		previousHealthy = &state.Health.Healthy

		// it is is important to note that the map values we retrieved this from are
		// stored by value. Therefore we are modifying a copy of what is in the existing
		// state and not the actual state itself. We want to ensure that the Address
		// is what Raft will know about.
		state.Server = existing.Server
		state.Server.Address = srv.Address
	}
	a.stateLock.RUnlock()

	// pull in the latest information from the applications knowledge of the
	// server. Mainly we want the NodeStatus & Meta
	if known, found := inputs.KnownServers[srv.ID]; found {
		// it is important to note that we are modifying a copy of a Server as the
		// map we retrieved this from has a non-pointer type value. We definitely
		// do not want to modify the current known servers but we do want to ensure
		// that we do not overwrite the Address
		state.Server = *known
		state.Server.Address = srv.Address
	} else {
		// TODO (mkeeler) do we need a None state. In the previous autopilot code
		// we would have set this to serf.StatusNone
		state.Server.NodeStatus = NodeLeft
	}

	// override the Stats if any where in the fetched results
	if stats, found := inputs.FetchedStats[srv.ID]; found {
		state.Stats = *stats
	}

	// now populate the healthy field given the stats
	state.Health.Healthy = state.isHealthy(inputs.LastTerm, inputs.LatestIndex, inputs.Config)
	// overwrite the StableSince field if this is a new server or when
	// the health status changes. No need for an else as we previously set
	// it when we overwrote the whole Health structure when finding a
	// server in the existing state
	if previousHealthy == nil || *previousHealthy != state.Health.Healthy {
		state.Health.StableSince = inputs.Now
	}

	return state
}

// updateState will compute the nextState, set it on the Autopilot instance and
// then notify the delegate of the update.
func (a *Autopilot) updateState(ctx context.Context) {
	newState, err := a.nextState(ctx)
	if err != nil {
		a.logger.Error("Error when computing next state", "error", err)
		return
	}

	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	a.state = newState
	a.delegate.NotifyState(newState)
}

// SortServers will take a list of raft ServerIDs and sort it using
// information from the State. See the ServerLessThan function for
// details about how two servers get compared.
func SortServers(ids []raft.ServerID, s *State) {
	sort.Slice(ids, func(i, j int) bool {
		return ServerLessThan(ids[i], ids[j], s)
	})
}

// ServerLessThan will lookup both servers in the given State and return
// true if the first id corresponds to a server that is logically less than
// lower than, better than etc. the second server. The following criteria
// are considered in order of most important to least important
//
// 1. A Leader server is always less than all others
// 2. A voter is less than non voters
// 3. Healthy servers are less than unhealthy servers
// 4. Servers that have been stable longer are consider less than.
func ServerLessThan(id1 raft.ServerID, id2 raft.ServerID, s *State) bool {
	srvI := s.Servers[id1]
	srvJ := s.Servers[id2]

	// the leader always comes first
	if srvI.State == RaftLeader {
		return true
	} else if srvJ.State == RaftLeader {
		return false
	}

	// voters come before non-voters & staging
	if srvI.State == RaftVoter && srvJ.State != RaftVoter {
		return true
	} else if srvI.State != RaftVoter && srvJ.State == RaftVoter {
		return false
	}

	// at this point we know that the raft state of both nodes is roughly
	// equivalent so we want to now sort based on health
	if srvI.Health.Healthy == srvJ.Health.Healthy {
		if srvI.Health.StableSince.Before(srvJ.Health.StableSince) {
			return srvI.Health.Healthy
		} else if srvJ.Health.StableSince.Before(srvI.Health.StableSince) {
			return !srvI.Health.Healthy
		}

		// with all else equal sort by the IDs
		return id1 < id2
	}

	// one of the two isn't healthy. We consider the healthy one as less than
	// the other. So we return true if server I is healthy and false if it isn't
	// as we know that server J is healthy and thus should come before server I.
	return srvI.Health.Healthy
}
