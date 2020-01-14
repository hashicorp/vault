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

var getAWSKMSFunc = func(opts *wrapping.WrapperOptions, config map[string]string) (wrapping.Wrapper, map[string]string, error) {
	kms := awskms.NewWrapper(nil)
	kmsInfo, err := kms.SetConfig(config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	return kms, kmsInfo, nil
}

func configureAWSKMSSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger hclog.Logger, inseal vault.Seal) (vault.Seal, error) {
	kms, kmsInfo, err := getAWSKMSFunc(nil, configSeal.Config)
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
