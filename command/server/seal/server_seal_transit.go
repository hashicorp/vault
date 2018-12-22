package seal

import (
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal/transit"
)

func configureTransitSeal(config *server.Config, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	transitSeal := transit.NewSeal(logger)
	sealInfo, err := transitSeal.SetConfig(config.Seal.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, err
		}
	}
	autoseal := vault.NewAutoSeal(transitSeal)
	if sealInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "Transit Seal Address", "Transit Seal Mount Path", "Transit Seal Key Name")
		(*info)["Seal Type"] = config.Seal.Type
		(*info)["Transit Seal Address"] = sealInfo["address"]
		(*info)["Transit Seal Mount Path"] = sealInfo["mount_path"]
		(*info)["Transit Seal Key Name"] = sealInfo["key_name"]
		if endpoint, ok := sealInfo["namespace"]; ok {
			*infoKeys = append(*infoKeys, "Transit Seal Namespace")
			(*info)["AWS KMS Endpoint"] = endpoint
		}
	}
	return autoseal, nil
}
