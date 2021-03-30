package autopilot

import (
	"fmt"
	"sort"

	"github.com/hashicorp/raft"
)

// reconcile calculates and then applies promotions and demotions
func (a *Autopilot) reconcile() error {
	conf := a.delegate.AutopilotConfig()
	if conf == nil {
		return nil
	}

	// grab the current state while locked
	a.stateLock.Lock()
	state := a.state
	a.stateLock.Unlock()

	if state == nil || state.Leader == "" {
		return fmt.Errorf("Cannote reconcile Raft server voting rights without a valid autopilot state")
	}

	// have the promoter calculate the required Raft changeset.
	changes := a.promoter.CalculatePromotionsAndDemotions(conf, state)

	// apply the promotions, if we did apply any then stop here
	// as we do not want to apply the demotions at the same time
	// as a means of preventing cluster instability.
	if done, err := a.applyPromotions(state, changes); done {
		return err
	}

	// apply the demotions, if we did apply any then stop here
	// as we do not want to transition leadership and do demotions
	// at the same time. This is a preventative measure to maintain
	// cluster stability.
	if done, err := a.applyDemotions(state, changes); done {
		return err
	}

	// if no leadership transfer is desired then we can exit the method now.
	if changes.Leader == "" || changes.Leader == state.Leader {
		return nil
	}

	// lookup the server we want to transfer leadership to
	srv, ok := state.Servers[changes.Leader]
	if !ok {
		return fmt.Errorf("cannot transfer leadership to an unknown server with ID %s", changes.Leader)
	}

	// perform the leadership transfer
	return a.leadershipTransfer(changes.Leader, srv.Server.Address)
}

// applyPromotions will apply all the promotions in the RaftChanges parameter.
//
// IDs in the change set will be ignored if:
// * The server isn't tracked in the provided state
// * The server already has voting rights
// * The server is not healthy
//
// If any servers were promoted this function returns true for the bool value.
func (a *Autopilot) applyPromotions(state *State, changes RaftChanges) (bool, error) {
	promoted := false
	for _, change := range changes.Promotions {
		srv, found := state.Servers[change]
		if !found {
			a.logger.Debug("Ignoring promotion of server as it is not in the autopilot state", "id", change)
			// this shouldn't be able to happen but is a nice safety measure against the
			// delegate doing something less than desirable
			continue
		}

		if srv.HasVotingRights() {
			// There is no need to promote as this server is already a voter.
			// No logging is needed here as this could be a very common case
			// where the promoter just returns a lists of server ids that should
			// be voters and non-voters without caring about which ones currently
			// already are in that state.
			a.logger.Debug("Not promoting server that already has voting rights", "id", change)
			continue
		}

		if !srv.Health.Healthy {
			// do not promote unhealthy servers
			a.logger.Debug("Ignoring promotion of unhealthy server", "id", change)
			continue
		}

		a.logger.Info("Promoting server", "id", srv.Server.ID, "address", srv.Server.Address, "name", srv.Server.Name)

		if err := a.addVoter(srv.Server.ID, srv.Server.Address); err != nil {
			return true, fmt.Errorf("failed promoting server %s: %v", srv.Server.ID, err)
		}

		promoted = true
	}

	// when we promoted anything we return true to indicate that the promotion/demotion applying
	// process is finished to prevent promotions and demotions in the same round. This is what
	// autopilot within Consul used to do so I am keeping the behavior the same for now.
	return promoted, nil
}

// applyDemotions will apply all the demotions in the RaftChanges parameter.
//
// IDs in the change set will be ignored if:
// * The server isn't tracked in the provided state
// * The server does not have voting rights
//
// If any servers were demoted this function returns true for the bool value.
func (a *Autopilot) applyDemotions(state *State, changes RaftChanges) (bool, error) {
	demoted := false
	for _, change := range changes.Demotions {
		srv, found := state.Servers[change]
		if !found {
			a.logger.Debug("Ignoring demotion of server as it is not in the autopilot state", "id", change)
			// this shouldn't be able to happen but is a nice safety measure against the
			// delegate doing something less than desirable
			continue
		}

		if srv.State == RaftNonVoter {
			// There is no need to demote as this server is already a non-voter.
			// No logging is needed here as this could be a very common case
			// where the promoter just returns a lists of server ids that should
			// be voters and non-voters without caring about which ones currently
			// already are in that state.
			a.logger.Debug("Ignoring demotion of server that is already a non-voter", "id", change)
			continue
		}

		a.logger.Info("Demoting server", "id", srv.Server.ID, "address", srv.Server.Address, "name", srv.Server.Name)

		if err := a.demoteVoter(srv.Server.ID); err != nil {
			return true, fmt.Errorf("failed demoting server %s: %v", srv.Server.ID, err)
		}

		demoted = true
	}

	// similarly to applyPromotions here we want to stop the process and prevent leadership
	// transfer when any demotions took place. Basically we want to ensure the cluster is
	// stable before doing the transfer
	return demoted, nil
}

// getFailedServers aggregates all of the information about servers that the consuming application believes are in
// a failed/left state (indicated by the NodeStatus field on the Server type) as well as stale servers that are
// in the raft configuration but not know to the consuming application. This function will do nothing with
// that information and is purely to collect the data.
func (a *Autopilot) getFailedServers() (*FailedServers, int, error) {
	staleRaftServers := make(map[raft.ServerID]raft.Server)
	raftConfig, err := a.getRaftConfiguration()
	if err != nil {
		return nil, 0, err
	}

	// Populate a map of all the raft servers. We will
	// remove some later on from the map leaving us with
	// just the stale servers.
	var voters int
	for _, server := range raftConfig.Servers {
		staleRaftServers[server.ID] = server

		if server.Suffrage == raft.Voter {
			voters++
		}
	}

	var failed FailedServers
	for id, srv := range a.delegate.KnownServers() {
		raftSrv, found := staleRaftServers[id]
		if found {
			delete(staleRaftServers, id)
		}

		if srv.NodeStatus != NodeAlive {
			if found && raftSrv.Suffrage == raft.Voter {
				failed.FailedVoters = append(failed.FailedVoters, srv)
			} else if found {
				failed.FailedNonVoters = append(failed.FailedNonVoters, srv)
			}
		}
	}

	for id, srv := range staleRaftServers {
		if srv.Suffrage == raft.Voter {
			failed.StaleVoters = append(failed.StaleVoters, id)
		} else {
			failed.StaleNonVoters = append(failed.StaleNonVoters, id)
		}
	}

	sort.Slice(failed.StaleNonVoters, func(i, j int) bool {
		return failed.StaleNonVoters[i] < failed.StaleNonVoters[j]
	})
	sort.Slice(failed.StaleVoters, func(i, j int) bool {
		return failed.StaleVoters[i] < failed.StaleVoters[j]
	})
	sort.Slice(failed.FailedNonVoters, func(i, j int) bool {
		return failed.FailedNonVoters[i].ID < failed.FailedNonVoters[j].ID
	})
	sort.Slice(failed.FailedVoters, func(i, j int) bool {
		return failed.FailedVoters[i].ID < failed.FailedVoters[j].ID
	})

	return &failed, voters, nil
}

// pruneDeadServers will find stale raft servers and failed servers as indicated by the consuming application
// and remove them. For stale raft servers this means removing them from the Raft configuration. For failed
// servers this means issuing RemoveFailedNode calls to the delegate. All stale/failed non-voters will be
// removed first. Then stale voters and finally failed servers. For servers with voting rights we will
// cap the number removed so that we do not remove too many at a time and do not remove nodes to the
// point where the number of voters would be below the MinQuorum value from the autopilot config.
// Additionally the delegate will be consulted to determine if all of the removals should be done and
// can filter the failed servers listings if need be.
func (a *Autopilot) pruneDeadServers() error {
	conf := a.delegate.AutopilotConfig()
	if conf == nil || !conf.CleanupDeadServers {
		return nil
	}

	state := a.GetState()

	failed, voters, err := a.getFailedServers()
	if err != nil || failed == nil {
		return err
	}

	failed = a.promoter.FilterFailedServerRemovals(conf, state, failed)

	// remove failed non voting servers
	for _, srv := range failed.FailedNonVoters {
		a.logger.Info("Attempting removal of failed server node", "id", srv.ID, "name", srv.Name, "address", srv.Address)
		a.delegate.RemoveFailedServer(srv)
	}

	// remove stale non voters
	for _, id := range failed.StaleNonVoters {
		a.logger.Debug("removing stale raft server from configuration", "id", id)
		if err := a.removeServer(id); err != nil {
			return err
		}
	}

	maxRemoval := (voters - 1) / 2

	for _, id := range failed.StaleVoters {
		if voters-1 < int(conf.MinQuorum) {
			a.logger.Debug("will not remove server as it would leave less voters than the minimum number allowed", "id", id, "min", conf.MinQuorum)
		} else if maxRemoval < 1 {
			a.logger.Debug("will not remove server as removal of a majority or servers is not safe", "id", id)
		} else if err := a.removeServer(id); err != nil {
			return err
		} else {
			maxRemoval--
			voters--
		}
	}

	for _, srv := range failed.FailedVoters {
		if voters-1 < int(conf.MinQuorum) {
			a.logger.Debug("will not remove server as it would leave less voters than the minimum number allowed", "id", srv.ID, "min", conf.MinQuorum)
		} else if maxRemoval < 1 {
			a.logger.Debug("will not remove server as its removal would be unsafe due to affectingas removal of a majority or servers is not safe", "id", srv.ID)
		} else {
			a.logger.Info("Attempting removal of failed server node", "id", srv.ID, "name", srv.Name, "address", srv.Address)
			a.delegate.RemoveFailedServer(srv)
			maxRemoval--
			voters--
		}
	}

	return nil
}
