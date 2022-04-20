package main

import (
	"os"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {

	if err := plugin.ServeMultiplex(&plugin.ServeOpts{
		BackendFactoryFunc: userpass.Factory,
	}); err != nil {
		logger := hclog.New(&hclog.LoggerOptions{})

		logger.Error("plugin shutting down", "error", err)
		os.Exit(1)
	}
}
