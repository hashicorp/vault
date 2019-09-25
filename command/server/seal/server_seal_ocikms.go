// Copyright © 2019, Oracle and/or its affiliates.
package seal

import (
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal/ocikms"
)

func configureOCIKMSSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	kms := ocikms.NewSeal(logger)
	kmsInfo, err := kms.SetConfig(configSeal.Config)
	if err != nil {
		logger.Error("error on setting up config for OCI KMS", "error", err)
		return nil, err
	}
	autoseal := vault.NewAutoSeal(kms)
	if kmsInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "OCI KMS KeyID")
		(*info)["Seal Type"] = configSeal.Type
		(*info)["OCI KMS KeyID"] = kmsInfo["key_id"]
	}
	return autoseal, nil
}
