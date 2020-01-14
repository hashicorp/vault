package seal

import (
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/go-kms-wrapping/wrappers/transit"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

var getTransitKMSFunc = func(opts *wrapping.WrapperOptions, config map[string]string, logger log.Logger) (wrapping.Wrapper, map[string]string, error) {
	transitSeal := transit.NewWrapper(opts)
	sealInfo, err := transitSeal.SetConfig(config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, nil, err
		}
	}
	return transitSeal, sealInfo, nil
}

func configureTransitSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	transitSeal, sealInfo, err := getTransitKMSFunc(
		&wrapping.WrapperOptions{
			Logger: logger.ResetNamed("seal-transit"),
		}, configSeal.Config, logger)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, err
		}
	}
	autoseal := vault.NewAutoSeal(&seal.Access{
		Wrapper: transitSeal,
	})
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
