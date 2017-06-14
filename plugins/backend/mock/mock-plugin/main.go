package main

import (
	"os"

	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical/plugin"
	"github.com/hashicorp/vault/plugins/backend/mock"
)

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

	plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: mock.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
}
