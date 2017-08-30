package main

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin"
	"github.com/hashicorp/vault/logical/plugin/mock"
)

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:]) // Ignore command, strictly parse flags

	// Set tlsProviderFunc if -metadata not passed in
	var tlsProviderFunc func() (*tls.Config, error)
	if !apiClientMeta.FetchMetadata() {
		tlsConfig := apiClientMeta.GetTLSConfig()
		tlsProviderFunc = pluginutil.VaultPluginTLSProvider(tlsConfig)
	}

	factoryFunc := mock.FactoryType(logical.TypeLogical)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: factoryFunc,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
