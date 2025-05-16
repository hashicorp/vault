// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/grpc/metadata"
)

// pbMetadataCtxToLogicalCtx extracts the snapshot ID key from an incoming GRPC
// context and adds the logical context key to the returned context
func pbMetadataCtxToLogicalCtx(ctx context.Context) context.Context {
	var snapshotID string
	snapshotIDs := metadata.ValueFromIncomingContext(ctx, snapshotIDCtxKey)
	if len(snapshotIDs) > 0 {
		snapshotID = snapshotIDs[0]
	}
	return logical.CreateContextWithSnapshotID(ctx, snapshotID)
}

// logicalCtxToPBMetadataCtx extracts the logical context snapshot ID key from
// the context and appends it to an outgoing GRPC context
func logicalCtxToPBMetadataCtx(ctx context.Context) context.Context {
	snapshotID, ok := logical.ContextSnapshotIDValue(ctx)
	if !ok {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, snapshotIDCtxKey, snapshotID)
}

const snapshotIDCtxKey string = "snapshot_id"
