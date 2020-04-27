package seal

import (
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-kms-wrapping/wrappers/tencentcloudkms"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func configureTencentCloudKMSSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	kms := tencentcloudkms.NewWrapper(nil)
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
		*infoKeys = append(*infoKeys, "Seal Type", "TencentCloud KMS Region", "TencentCloud KMS KeyID")
		(*info)["Seal Type"] = configSeal.Type
		(*info)["TencentCloud KMS Region"] = kmsInfo["region"]
		(*info)["TencentCloud KMS KeyID"] = kmsInfo["kms_key_id"]
	}
	return autoseal, nil
}
