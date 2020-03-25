// +build !enterprise

package raft

import (
	"context"
	"errors"
)

// AddPeer adds a new server to the raft cluster
func (b *RaftBackend) AddNonVotingPeer(ctx context.Context, peerID, clusterAddr string) error {
	return errors.New("not implemented")
}
