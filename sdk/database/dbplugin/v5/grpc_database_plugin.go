package dbplugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var HandshakeConfig = plugin.HandshakeConfig{
	MagicCookieKey:   "VAULT_DATABASE_PLUGIN",
	MagicCookieValue: "926a0820-aea2-be28-51d6-83cdf00e8edb",
}

const multiplexingCtxKey string = "multiplex_id"

// Factory is the factory function to create a dbplugin Database.
type Factory func() (Database, error)

type GRPCDatabasePlugin struct {
	FactoryFunc         Factory
	Impl                Database
	multiplexingSupport bool

	// Embeding this will disable the netRPC protocol
	plugin.NetRPCUnsupportedPlugin
}

var (
	_ plugin.Plugin     = &GRPCDatabasePlugin{}
	_ plugin.GRPCPlugin = &GRPCDatabasePlugin{}
)

func (d GRPCDatabasePlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	server := gRPCServer{factoryFunc: d.FactoryFunc, instances: make(map[string]Database)}

	proto.RegisterDatabaseServer(s, server)
	return nil
}

func (GRPCDatabasePlugin) GRPCClient(doneCtx context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	client := gRPCClient{
		client:  proto.NewDatabaseClient(c),
		doneCtx: doneCtx,
	}
	return client, nil
}

type databaseClientConn struct {
	*grpc.ClientConn
	id string
}

var _ grpc.ClientConnInterface = &databaseClientConn{}

func (d *databaseClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	// Inject ID to the context
	md := metadata.Pairs(multiplexingCtxKey, d.id)
	idCtx := metadata.NewOutgoingContext(ctx, md)

	return d.ClientConn.Invoke(idCtx, method, args, reply, opts...)
}

func (d *databaseClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// Inject ID to the context
	md := metadata.Pairs(multiplexingCtxKey, d.id)
	idCtx := metadata.NewOutgoingContext(ctx, md)

	return d.ClientConn.NewStream(idCtx, desc, method, opts...)
}
