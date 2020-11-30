package autopilot

import (
	"time"

	"github.com/hashicorp/raft"
)

func DefaultPromoter() Promoter {
	return new(StablePromoter)
}

type StablePromoter struct{}

func (_ *StablePromoter) GetServerExt(_ *Config, srv *ServerState) interface{} {
	return nil
}

func (_ *StablePromoter) GetStateExt(_ *Config, _ *State) interface{} {
	return nil
}

func (_ *StablePromoter) GetNodeTypes(_ *Config, s *State) map[raft.ServerID]NodeType {
	types := make(map[raft.ServerID]NodeType)
	for id := range s.Servers {
		// this basic implementation has all nodes be of the "voter" type regardless of
		// any other settings. That means that in a healthy state all nodes in the cluster
		// will be a voter.
		types[id] = NodeVoter
	}
	return types
}

func (_ *StablePromoter) FilterFailedServerRemovals(_ *Config, _ *State, failed *FailedServers) *FailedServers {
	return failed
}

// CalculatePromotionsAndDemotions will return a list of all promotions and demotions to be done as well as the server id of
// the desired leader. This particular interface implementation maintains a stable leader and will promote healthy servers
// to voting status. It will never change the leader ID nor will it perform demotions.
func (_ *StablePromoter) CalculatePromotionsAndDemotions(c *Config, s *State) RaftChanges {
	var changes RaftChanges

	now := time.Now()
	minStableDuration := s.ServerStabilizationTime(c)
	for id, server := range s.Servers {
		// ignore staging state as they are not ready yet
		if server.State == RaftNonVoter && server.Health.IsStable(now, minStableDuration) {
			changes.Promotions = append(changes.Promotions, id)
		}
	}

	return changes
}
