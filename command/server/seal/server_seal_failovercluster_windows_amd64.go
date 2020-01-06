package seal

import (
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal/failovercluster"
)

func configureFailoverClusterSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	kv := failovercluster.NewSeal(logger)
	kvInfo, err := kv.SetConfig(configSeal.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, err
		}
	}
	autoseal := vault.NewAutoSeal(kv)
	if kvInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "MSFC Resource Name")
		(*info)["Seal Type"] = configSeal.Type
		(*info)["MSFC Resource Name"] = kvInfo["resource_name"]
	}
	return autoseal, nil
}
