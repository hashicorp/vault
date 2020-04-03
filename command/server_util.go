package command

import (
	"crypto/rand"
	"io"

	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/vault"
)

var (
	createSecureRandomReaderFunc = createSecureRandomReader
	adjustCoreConfigForEnt       = adjustCoreConfigForEntNoop
)

func adjustCoreConfigForEntNoop(config *server.Config, coreConfig *vault.CoreConfig) {
}

func createSecureRandomReader(config *server.Config, seal *vault.Seal) (io.Reader, error) {
	return rand.Reader, nil
}
