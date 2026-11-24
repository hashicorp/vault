// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/grpc/metadata"
)

// pbMetadataCtxToLogicalCtx extracts the snapshot ID key from an incoming GRPC
// context and adds the logical context key to the returned context
func pbMetadataCtxToLogicalCtx(ctx context.Context) (context.Context, error) {
	var snapshotID string
	snapshotIDs := metadata.ValueFromIncomingContext(ctx, snapshotIDCtxKey)
	if len(snapshotIDs) > 0 {
		snapshotID = snapshotIDs[0]
		ctx = logical.CreateContextWithSnapshotID(ctx, snapshotID)
	}

	clusterID := metadata.ValueFromIncomingContext(ctx, indexStateCtxKeyClusterID)
	localRaw := metadata.ValueFromIncomingContext(ctx, indexStateCtxKeyLocal)
	replicatedRaw := metadata.ValueFromIncomingContext(ctx, indexStateCtxKeyReplicated)
	if len(clusterID) > 0 {
		local, err := strconv.ParseUint(localRaw[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing local index: %w", err)
		}
		replicated, err := strconv.ParseUint(replicatedRaw[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing replicated index: %w", err)
		}
		w := &logical.WALState{
			ClusterID:       clusterID[0],
			LocalIndex:      local,
			ReplicatedIndex: replicated,
		}
		ctx = logical.IndexStateContext(ctx, w)
	} else {
		ctx = logical.IndexStateContext(ctx, &logical.WALState{})
	}
	return ctx, nil
}

// logicalCtxToPBMetadataCtx extracts the logical context snapshot ID key from
// the context and appends it to an outgoing GRPC context
func logicalCtxToPBMetadataCtx(ctx context.Context) context.Context {
	var args []string
	if snapshotID, ok := logical.ContextSnapshotIDValue(ctx); ok {
		args = append(args, snapshotIDCtxKey, snapshotID)
	}
	if index := logical.IndexStateFromContext(ctx); index != nil {
		args = append(args, indexStateCtxKeyClusterID, index.ClusterID,
			indexStateCtxKeyLocal, fmt.Sprintf("%d", index.LocalIndex),
			indexStateCtxKeyReplicated, fmt.Sprintf("%d", index.ReplicatedIndex))
	}
	return metadata.AppendToOutgoingContext(ctx, args...)
}

const (
	snapshotIDCtxKey           string = "snapshot_id"
	indexStateCtxKey                  = "index_state"
	indexStateCtxKeyClusterID         = indexStateCtxKey + "_cluster_id"
	indexStateCtxKeyLocal             = indexStateCtxKey + "_local"
	indexStateCtxKeyReplicated        = indexStateCtxKey + "_replicated"
)
