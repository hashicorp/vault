// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package dbplugin

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// GRPCClient (Vault CE edition) initializes and returns a gRPCClient with Database and
// PluginVersion gRPC clients. It implements GRPCClient() defined
// by GRPCPlugin interface in go-plugin/plugin.go
func (GRPCDatabasePlugin) GRPCClient(doneCtx context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	client := gRPCClient{
		client:        proto.NewDatabaseClient(c),
		versionClient: logical.NewPluginVersionClient(c),
		doneCtx:       doneCtx,
	}
	return client, nil
}

// GRPCServer (Vault CE edition) registers multiplexing server if the plugin supports it, and
// registers the Database and PluginVersion gRPC servers. It implements GRPCServer() defined
// by GRPCPlugin interface in go-plugin/plugin.go
func (d GRPCDatabasePlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	var server gRPCServer

	if d.Impl != nil {
		server = gRPCServer{singleImpl: d.Impl}
	} else {
		// multiplexing is supported
		server = gRPCServer{
			factoryFunc: d.FactoryFunc,
			instances:   make(map[string]Database),
		}

		// Multiplexing is enabled for this plugin, register the server so we
		// can tell the client in Vault.
		pluginutil.RegisterPluginMultiplexingServer(s, pluginutil.PluginMultiplexingServerImpl{
			Supported: true,
		})
	}

	proto.RegisterDatabaseServer(s, &server)
	logical.RegisterPluginVersionServer(s, &server)
	return nil
}
