package seal

import (
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal/transit"
)

func configureTransitSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	transitSeal := transit.NewSeal(logger)
	sealInfo, err := transitSeal.SetConfig(configSeal.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, err
		}
	}
	autoseal := vault.NewAutoSeal(transitSeal)
	if sealInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "Transit Address", "Transit Mount Path", "Transit Key Name")
		(*info)["Seal Type"] = configSeal.Type
		(*info)["Transit Address"] = sealInfo["address"]
		(*info)["Transit Mount Path"] = sealInfo["mount_path"]
		(*info)["Transit Key Name"] = sealInfo["key_name"]
		if namespace, ok := sealInfo["namespace"]; ok {
			*infoKeys = append(*infoKeys, "Transit Namespace")
			(*info)["Transit Namespace"] = namespace
		}
	}
	return autoseal, nil
}
