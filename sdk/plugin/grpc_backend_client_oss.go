// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package plugin

import (
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	"google.golang.org/grpc"
)

// registerSystewViewServer (on Vault Community edition) registers the SystemView server
// to the gRPC service registrar
func registerSystewViewServer(s *grpc.Server, sysView logical.SystemView, _ *logical.BackendConfig) {
	pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
		impl: sysView,
	})
}
