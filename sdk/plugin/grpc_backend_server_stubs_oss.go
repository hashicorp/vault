// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package plugin

import (
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	"google.golang.org/grpc"
)

func newGRPCSystemViewFromSetupArgs(conn *grpc.ClientConn, _ *pb.SetupArgs) logical.SystemView {
	return newGRPCSystemView(conn)
}
