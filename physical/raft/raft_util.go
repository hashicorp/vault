// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package raft

import (
	"context"
	"errors"

	"github.com/hashicorp/raft"
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

func (d *Delegate) autopilotServerExtNonVoter(_ *FollowerState) interface{} {
	return nil
}

func (d *Delegate) metaWithoutZone(_ *FollowerState) map[string]string {
	return nil
}

func (d *Delegate) knownServer(id raft.ServerID, name, version string, state FollowerState, isVoter bool, cachedNodeType autopilot.NodeType) *autopilot.Server {
	nodeType := cachedNodeType
	if !isVoter {
		nodeType = ""
	} else if cachedNodeType == "" {
		nodeType = autopilot.NodeVoter
	}
	return &autopilot.Server{
		ID:          id,
		Name:        name,
		RaftVersion: raft.ProtocolVersionMax,
		Meta:        d.meta(&state),
		Version:     version,
		NodeType:    nodeType,
		Ext:         d.autopilotServerExtNonVoter(&state),
	}
}
