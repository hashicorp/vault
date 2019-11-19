package seal

import (
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/go-kms-wrapping/wrappers/awskms"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func configureAWSKMSSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger hclog.Logger, inseal vault.Seal) (vault.Seal, error) {
	kms := awskms.NewWrapper(&wrapping.WrapperOptions{
		Logger: logger,
	})
	kmsInfo, err := kms.SetConfig(configSeal.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, err
		}
	}
	autoseal := vault.NewAutoSeal(&seal.Access{
		Wrapper: kms,
	})
	if kmsInfo != nil {
		*infoKeys = append(*infoKeys, "Seal Type", "AWS KMS Region", "AWS KMS KeyID")
		(*info)["Seal Type"] = configSeal.Type
		(*info)["AWS KMS Region"] = kmsInfo["region"]
		(*info)["AWS KMS KeyID"] = kmsInfo["kms_key_id"]
		if endpoint, ok := kmsInfo["endpoint"]; ok {
			*infoKeys = append(*infoKeys, "AWS KMS Endpoint")
			(*info)["AWS KMS Endpoint"] = endpoint
		}
	}
	return autoseal, nil
}
