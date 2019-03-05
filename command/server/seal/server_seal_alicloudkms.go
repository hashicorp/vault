package seal

import (
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal/alicloudkms"
)

func configureAliCloudKMSSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	kms := alicloudkms.NewSeal(logger)
	kmsInfo, err := kms.SetConfig(configSeal.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, err
		}
	}
	autoseal := vault.NewAutoSeal(kms)
	if kmsInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "AliCloud KMS Region", "AliCloud KMS KeyID")
		(*info)["Seal Type"] = configSeal.Type
		(*info)["AliCloud KMS Region"] = kmsInfo["region"]
		(*info)["AliCloud KMS KeyID"] = kmsInfo["kms_key_id"]
		if domain, ok := kmsInfo["domain"]; ok {
			*infoKeys = append(*infoKeys, "AliCloud KMS Domain")
			(*info)["AliCloud KMS Domain"] = domain
		}
	}
	return autoseal, nil
}
