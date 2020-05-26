package command

import (
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/vault"
)

var (
	adjustCoreConfigForEnt = adjustCoreConfigForEntNoop
)

func adjustCoreConfigForEntNoop(config *server.Config, coreConfig *vault.CoreConfig) {
}
