package main

import (
	gcpbackend "github.com/hashicorp/vault-plugin-auth-gcp/plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical/plugin"
	"log"
	"os"
)

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args)

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: gcpbackend.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
