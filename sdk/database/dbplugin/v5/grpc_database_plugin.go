// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/grpc"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var HandshakeConfig = plugin.HandshakeConfig{
	MagicCookieKey:   "VAULT_DATABASE_PLUGIN",
	MagicCookieValue: "926a0820-aea2-be28-51d6-83cdf00e8edb",
}

// Factory is the factory function to create a dbplugin Database.
type Factory func() (interface{}, error)

type GRPCDatabasePlugin struct {
	FactoryFunc Factory
	Impl        Database

	// Embeding this will disable the netRPC protocol
	plugin.NetRPCUnsupportedPlugin
}

var (
	_ plugin.Plugin     = &GRPCDatabasePlugin{}
	_ plugin.GRPCPlugin = &GRPCDatabasePlugin{}
)

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

func (GRPCDatabasePlugin) GRPCClient(doneCtx context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	client := gRPCClient{
		client:        proto.NewDatabaseClient(c),
		versionClient: logical.NewPluginVersionClient(c),
		doneCtx:       doneCtx,
	}
	return client, nil
}
