package plugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/logical"
)

// BackendPlugin is the plugin.Plugin implementation
type BackendPlugin struct {
	Factory      func(*logical.BackendConfig) (logical.Backend, error)
	metadataMode bool
}

// Server gets called when on plugin.Serve()
func (b *BackendPlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &backendPluginServer{factory: b.Factory, broker: broker}, nil
}

// Client gets called on plugin.NewClient()
func (b BackendPlugin) Client(broker *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &backendPluginClient{client: c, broker: broker, metadataMode: b.metadataMode}, nil
}
