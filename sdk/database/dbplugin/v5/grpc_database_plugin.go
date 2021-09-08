package dbplugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"google.golang.org/grpc"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
// TODO: Does this need to change?
var handshakeConfig = plugin.HandshakeConfig{
	// TODO: Can this be removed since we're using versioned plugins?
	// ProtocolVersion:  5,
	MagicCookieKey:   "VAULT_DATABASE_PLUGIN",
	MagicCookieValue: "926a0820-aea2-be28-51d6-83cdf00e8edb",
}

type GRPCDatabasePlugin struct {
	DBFactory func() (Database, error)

	// Embeding this will disable the netRPC protocol
	plugin.NetRPCUnsupportedPlugin
}

var (
	_ plugin.Plugin     = &GRPCDatabasePlugin{}
	_ plugin.GRPCPlugin = &GRPCDatabasePlugin{}
)

func (d GRPCDatabasePlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterDatabaseServer(s, &gRPCServer{dbFactory: d.DBFactory})
	return nil
}

func (GRPCDatabasePlugin) GRPCClient(doneCtx context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	client := &gRPCClient{
		client:  proto.NewDatabaseClient(c),
		doneCtx: doneCtx,
	}
	return client, nil
}
