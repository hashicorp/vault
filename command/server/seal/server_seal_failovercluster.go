// +build !windows !amd64

package seal

import (
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func configureFailoverClusterSeal(*server.Seal, *[]string, *map[string]string, log.Logger, vault.Seal) (vault.Seal, error) {
	return nil, logical.ErrUnsupportedOperation
}
