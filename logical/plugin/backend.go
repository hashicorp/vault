package plugin

import (
	"context"
	"net/rpc"
	"sync/atomic"

	"google.golang.org/grpc"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin/pb"
)

var _ plugin.Plugin = (*BackendPlugin)(nil)
var _ plugin.GRPCPlugin = (*BackendPlugin)(nil)
var _ plugin.Plugin = (*GRPCBackendPlugin)(nil)
var _ plugin.GRPCPlugin = (*GRPCBackendPlugin)(nil)

// BackendPlugin is the plugin.Plugin implementation
type BackendPlugin struct {
	*GRPCBackendPlugin
}

// GRPCBackendPlugin is the plugin.Plugin implementation that only supports GRPC
// transport
type GRPCBackendPlugin struct {
	Factory      logical.Factory
	MetadataMode bool
	Logger       log.Logger

	// Embeding this will disable the netRPC protocol
	plugin.NetRPCUnsupportedPlugin
}

// Server gets called when on plugin.Serve()
func (b *BackendPlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &backendPluginServer{
		factory: b.Factory,
		broker:  broker,
		// We pass the logger down into the backend so go-plugin will forward
		// logs for us.
		logger: b.Logger,
	}, nil
}

// Client gets called on plugin.NewClient()
func (b BackendPlugin) Client(broker *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &backendPluginClient{
		client:       c,
		broker:       broker,
		metadataMode: b.MetadataMode,
	}, nil
}

func (b GRPCBackendPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterBackendServer(s, &backendGRPCPluginServer{
		broker:  broker,
		factory: b.Factory,
		// We pass the logger down into the backend so go-plugin will forward
		// logs for us.
		logger: b.Logger,
	})
	return nil
}

func (b *GRPCBackendPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	ret := &backendGRPCPluginClient{
		client:       pb.NewBackendClient(c),
		clientConn:   c,
		broker:       broker,
		cleanupCh:    make(chan struct{}),
		doneCtx:      ctx,
		metadataMode: b.MetadataMode,
	}

	// Create the value and set the type
	ret.server = new(atomic.Value)
	ret.server.Store((*grpc.Server)(nil))

	return ret, nil
}
