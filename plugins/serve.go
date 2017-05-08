package plugins

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/pluginutil"
)

// Serve is used to start a plugin's RPC server. It takes an interface that must
// implement a known plugin interface to vault and an optional api.TLSConfig for
// use during the inital unwrap request to vault. The api config is particulary
// useful when vault is setup to require client cert checking.
func Serve(plugin interface{}, tlsConfig *api.TLSConfig) {
	tlsProvider := pluginutil.VaultPluginTLSProvider(tlsConfig)

	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		fmt.Println(err)
		return
	}

	switch p := plugin.(type) {
	case dbplugin.Database:
		dbplugin.Serve(p, tlsProvider)
	default:
		fmt.Println("Unsupported plugin type")
	}

}
