package seal

import (
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal/local"
)

func configureLocalSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (vault.Seal, error) {
	localSeal := local.NewSeal(logger)
	sealInfo, err := localSeal.SetConfig(configSeal.Config)
	if err != nil {
		// If the error is any other than logical.KeyNotFoundError, return the error
		if !errwrap.ContainsType(err, new(logical.KeyNotFoundError)) {
			return nil, err
		}
	}

	autoseal := vault.NewAutoSeal(localSeal)
	if sealInfo != nil {
		// set data about our seal that's written to the server output
		// on startup
		*infoKeys = append(*infoKeys, "Seal Type", "Key Glob")
		(*info)["Seal Type"] = configSeal.Type
		(*info)["Key Glob"] = sealInfo["key_glob"]
	}
	return autoseal, nil
}
