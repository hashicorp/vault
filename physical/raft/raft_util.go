// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

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

// AddNonVotingPeer adds a new server to the raft cluster
func (b *RaftBackend) AddNonVotingPeer(ctx context.Context, peerID, clusterAddr string) error {
	return errors.New("adding non voting peer is not allowed")
}

func autopilotToAPIServerEnterprise(_ *autopilot.Server, _ *AutopilotServer) error {
	return nil
}

func autopilotToAPIStateEnterprise(_ *autopilot.State, _ *AutopilotState) error {
	return nil
}

func (d *Delegate) autopilotConfigExt() interface{} {
	return nil
}

func (d *Delegate) autopilotServerExt(_ *FollowerState) interface{} {
	return nil
}

func (d *Delegate) meta(_ *FollowerState) map[string]string {
	return nil
}
