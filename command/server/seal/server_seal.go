package seal

import (
	"fmt"
	"os"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

var (
	ConfigureSeal func(*server.Config, *[]string, *map[string]string, log.Logger, vault.Seal) (vault.Seal, error) = configureSeal
)

func configureSeal(config *server.Config, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (outseal vault.Seal, err error) {
	if config.Seal != nil || os.Getenv("VAULT_SEAL_TYPE") != "" {
		if config.Seal == nil {
			config.Seal = &server.Seal{
				Type: os.Getenv("VAULT_SEAL_TYPE"),
			}
		}
		switch config.Seal.Type {
		case seal.AliCloudKMS:
			return configureAliCloudKMSSeal(config, infoKeys, info, logger, inseal)

		case seal.AWSKMS:
			return configureAWSKMSSeal(config, infoKeys, info, logger, inseal)

		case seal.GCPCKMS:
			return configureGCPCKMSSeal(config, infoKeys, info, logger, inseal)

		case seal.AzureKeyVault:
			return configureAzureKeyVaultSeal(config, infoKeys, info, logger, inseal)

		case seal.PKCS11:
			return nil, fmt.Errorf("Seal type 'pkcs11' requires the Vault Enterprise HSM binary")

		default:
			return nil, fmt.Errorf("Unknown seal type %q", config.Seal.Type)
		}
	}

	return inseal, nil
}
