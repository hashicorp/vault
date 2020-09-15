package newdbplugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/database/newdbplugin/proto"
	"google.golang.org/grpc"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  5,
	MagicCookieKey:   "VAULT_DATABASE_PLUGIN",
	MagicCookieValue: "926a0820-aea2-be28-51d6-83cdf00e8edb",
}

type GRPCDatabasePlugin struct {
	Impl Database

	// Embeding this will disable the netRPC protocol
	plugin.NetRPCUnsupportedPlugin
}

var _ plugin.Plugin = &GRPCDatabasePlugin{}
var _ plugin.GRPCPlugin = &GRPCDatabasePlugin{}

func (d GRPCDatabasePlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterDatabaseServer(s, gRPCServer{impl: d.Impl})
	return nil
}

func (GRPCDatabasePlugin) GRPCClient(doneCtx context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	client := gRPCClient{
		client:  proto.NewDatabaseClient(c),
		doneCtx: doneCtx,
	}
	return client, nil
}
