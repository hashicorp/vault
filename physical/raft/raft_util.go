// +build !enterprise

package raft

import (
	"context"
	"errors"
)

const readReplicasAllowed = false

// AddReadReplicaPeer adds a new server to the raft cluster which does not have
// voting rights but gets all the data replicated to it.
func (b *RaftBackend) AddReadReplicaPeer(ctx context.Context, peerID, clusterAddr string) error {
	return errors.New("not implemented")
}
