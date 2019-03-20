package seal

import (
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal/gcpckms"
)

func configureGCPCKMSSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	kms := gcpckms.NewSeal(logger)
	kmsInfo, err := kms.SetConfig(configSeal.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, err
		}
	}
	autoseal := vault.NewAutoSeal(kms)
	if kmsInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "GCP KMS Project", "GCP KMS Region", "GCP KMS Key Ring", "GCP KMS Crypto Key")
		(*info)["Seal Type"] = configSeal.Type
		(*info)["GCP KMS Project"] = kmsInfo["project"]
		(*info)["GCP KMS Region"] = kmsInfo["region"]
		(*info)["GCP KMS Key Ring"] = kmsInfo["key_ring"]
		(*info)["GCP KMS Crypto Key"] = kmsInfo["crypto_key"]
	}
	return autoseal, nil
}
