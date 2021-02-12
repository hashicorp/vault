// +build !enterprise

package raft

import (
	"context"
	"errors"

	autopilot "github.com/hashicorp/raft-autopilot"
)

const nonVotersAllowed = false

func (b *RaftBackend) autopilotPromoter() autopilot.Promoter {
	return autopilot.DefaultPromoter()
}

// AddPeer adds a new server to the raft cluster
func (b *RaftBackend) AddNonVotingPeer(ctx context.Context, peerID, clusterAddr string) error {
	return errors.New("not implemented")
}

func autopilotToAPIServerEnterprise(_ *autopilot.ServerState, _ *AutopilotServer) {
	// noop in oss
}

func autopilotToAPIStateEnterprise(state *autopilot.State, apiState *AutopilotHealth) {
	// without the enterprise features there is no different between these two and we don't want to
	// alarm anyone by leaving this as the zero value.
	apiState.OptimisticFailureTolerance = state.FailureTolerance
}

func (d *Delegate) autopilotConfigExt() interface{} {
	return nil
}
