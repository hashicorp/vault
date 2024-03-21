// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"context"

	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type pluginClientConn struct {
	*grpc.ClientConn
	id string
}

var _ grpc.ClientConnInterface = &pluginClientConn{}

func (d *pluginClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	// Inject ID to the context
	md := metadata.Pairs(pluginutil.MultiplexingCtxKey, d.id)
	idCtx := metadata.NewOutgoingContext(ctx, md)

	return d.ClientConn.Invoke(idCtx, method, args, reply, opts...)
}

func (d *pluginClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// Inject ID to the context
	md := metadata.Pairs(pluginutil.MultiplexingCtxKey, d.id)
	idCtx := metadata.NewOutgoingContext(ctx, md)

	return d.ClientConn.NewStream(idCtx, desc, method, opts...)
}
