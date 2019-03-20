package seal

import (
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal/azurekeyvault"
)

func configureAzureKeyVaultSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	kv := azurekeyvault.NewSeal(logger)
	kvInfo, err := kv.SetConfig(configSeal.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, err
		}
	}
	autoseal := vault.NewAutoSeal(kv)
	if kvInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "Azure Environment", "Azure Vault Name", "Azure Key Name")
		(*info)["Seal Type"] = configSeal.Type
		(*info)["Azure Environment"] = kvInfo["environment"]
		(*info)["Azure Vault Name"] = kvInfo["vault_name"]
		(*info)["Azure Key Name"] = kvInfo["key_name"]
	}
	return autoseal, nil
}
