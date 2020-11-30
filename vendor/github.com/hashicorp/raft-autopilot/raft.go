package autopilot

//
// The methods in this file are all mainly to provide synchronous methods
// for Raft operations that would normally return futures.
//

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/raft"
)

func requiredQuorum(voters int) int {
	return (voters / 2) + 1
}

// NumVoters is a helper for calculating the number of voting peers in the
// current raft configuration. This function ignores any autopilot state
// and will make the calculation based on a newly retrieved Raft configuration.
func (a *Autopilot) NumVoters() (int, error) {
	cfg, err := a.getRaftConfiguration()
	if err != nil {
		return 0, err
	}

	var numVoters int
	for _, server := range cfg.Servers {
		if server.Suffrage == raft.Voter {
			numVoters++
		}
	}

	return numVoters, nil
}

// AddServer is a helper for adding a new server to the raft configuration.
// This may remove servers with duplicate addresses or ids first and after
// its all done will trigger autopilot to remove dead servers if there
// are any. Servers added by this method will start in a non-voting
// state and later on autopilot will promote them to voting status
// if desired by the configured promoter. If too many removals would
// be required that would cause leadership loss then an error is returned
// instead of performing any Raft configuration changes.
func (a *Autopilot) AddServer(s *Server) error {
	cfg, err := a.getRaftConfiguration()
	if err != nil {
		a.logger.Error("failed to get raft configuration", "error", err)
		return err
	}

	var existingVoter bool
	var voterRemovals []raft.ServerID
	var nonVoterRemovals []raft.ServerID
	var numVoters int
	for _, server := range cfg.Servers {
		if server.Suffrage == raft.Voter {
			numVoters++
		}

		if server.Address == s.Address && server.ID == s.ID {
			// nothing to be done as the addr and ID both already match
			return nil
		} else if server.ID == s.ID {
			// special case for address updates only. In this case we should be
			// able to update the configuration without have to first remove the server
			if server.Suffrage == raft.Voter || server.Suffrage == raft.Staging {
				existingVoter = true
			}
		} else if server.Address == s.Address {
			if server.Suffrage == raft.Voter {
				voterRemovals = append(voterRemovals, server.ID)
			} else {
				nonVoterRemovals = append(nonVoterRemovals, server.ID)
			}
		}
	}

	requiredVoters := requiredQuorum(numVoters)
	if len(voterRemovals) > numVoters-requiredVoters {
		return fmt.Errorf("Preventing server addition that would require removal of too many servers and cause cluster instability")
	}

	for _, id := range voterRemovals {
		if err := a.removeServer(id); err != nil {
			return fmt.Errorf("error removing server %q with duplicate address %q: %w", id, s.Address, err)
		}
		a.logger.Info("removed server with duplicate address", "address", s.Address)
	}

	for _, id := range nonVoterRemovals {
		if err := a.removeServer(id); err != nil {
			return fmt.Errorf("error removing server %q with duplicate address %q: %w", id, s.Address, err)
		}
		a.logger.Info("removed server with duplicate address", "address", s.Address)
	}

	if existingVoter {
		if err := a.addVoter(s.ID, s.Address); err != nil {
			return err
		}
	} else {
		if err := a.addNonVoter(s.ID, s.Address); err != nil {
			return err
		}
	}

	// Trigger a check to remove dead servers
	a.RemoveDeadServers()
	return nil
}

// RemoveServer is a helper to remove a server from Raft if it
// exists in the latest Raft configuration
func (a *Autopilot) RemoveServer(id raft.ServerID) error {
	cfg, err := a.getRaftConfiguration()
	if err != nil {
		a.logger.Error("failed to get raft configuration", "error", err)
		return err
	}

	// only remove servers currently in the configuration
	for _, server := range cfg.Servers {
		if server.ID == id {
			return a.removeServer(server.ID)
		}
	}

	return nil
}

// addNonVoter is a wrapper around calling the AddNonVoter method on the Raft
// interface object provided to Autopilot
func (a *Autopilot) addNonVoter(id raft.ServerID, addr raft.ServerAddress) error {
	addFuture := a.raft.AddNonvoter(id, addr, 0, 0)
	if err := addFuture.Error(); err != nil {
		a.logger.Error("failed to add raft non-voting peer", "id", id, "address", addr, "error", err)
		return err
	}
	return nil
}

// addVoter is a wrapper around calling the AddVoter method on the Raft
// interface object provided to Autopilot
func (a *Autopilot) addVoter(id raft.ServerID, addr raft.ServerAddress) error {
	addFuture := a.raft.AddVoter(id, addr, 0, 0)
	if err := addFuture.Error(); err != nil {
		a.logger.Error("failed to add raft voting peer", "id", id, "address", addr, "error", err)
		return err
	}
	return nil
}

func (a *Autopilot) demoteVoter(id raft.ServerID) error {
	removeFuture := a.raft.DemoteVoter(id, 0, 0)
	if err := removeFuture.Error(); err != nil {
		a.logger.Error("failed to demote raft peer", "id", id, "error", err)
		return err
	}
	return nil
}

// removeServer is a wrapper around calling the RemoveServer method on the
// Raft interface object provided to Autopilot
func (a *Autopilot) removeServer(id raft.ServerID) error {
	a.logger.Debug("removing server by ID", "id", id)
	future := a.raft.RemoveServer(id, 0, 0)
	if err := future.Error(); err != nil {
		a.logger.Error("failed to remove raft server",
			"id", id,
			"error", err,
		)
		return err
	}
	a.logger.Info("removed server", "id", id)
	return nil
}

// getRaftConfiguration a wrapper arond calling the GetConfiguration method
// on the Raft interface object provided to Autopilot
func (a *Autopilot) getRaftConfiguration() (*raft.Configuration, error) {
	configFuture := a.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return nil, err
	}
	cfg := configFuture.Configuration()
	return &cfg, nil
}

// lastTerm will retrieve the raft stats and then pull the last term value out of it
func (a *Autopilot) lastTerm() (uint64, error) {
	return strconv.ParseUint(a.raft.Stats()["last_log_term"], 10, 64)
}

// leadershipTransfer will transfer leadership to the server with the specified id and address
func (a *Autopilot) leadershipTransfer(id raft.ServerID, address raft.ServerAddress) error {
	a.logger.Info("Transferring leadership to new server", "id", id, "address", address)
	future := a.raft.LeadershipTransferToServer(id, address)
	return future.Error()
}
