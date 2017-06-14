package plugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/logical"
)

// BackendPlugin is the plugin.Plugin implementation
type BackendPlugin struct {
	Factory func() (logical.Backend, error)
}

// Server gets called when on plugin.Serve()
func (b *BackendPlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	backend, err := b.Factory()
	if err != nil {
		return nil, err
	}
	return &backendPluginServer{backend: backend, broker: broker}, nil
}

// Client gets called on plugin.NewClient()
func (b BackendPlugin) Client(broker *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &backendPluginClient{client: c, broker: broker}, nil
}
