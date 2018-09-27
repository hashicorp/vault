package seal

import (
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/vault"
)

var (
	ConfigureSeal func(*server.Config, *[]string, *map[string]string, log.Logger, vault.Seal) (vault.Seal, error) = configureSeal
)

func configureSeal(config *server.Config, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (seal vault.Seal, err error) {
	return inseal, nil
}
