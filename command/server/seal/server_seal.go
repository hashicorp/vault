package seal

import (
	"fmt"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

var (
	ConfigureSeal = configureSeal
)

func configureSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger log.Logger, inseal vault.Seal) (outseal vault.Seal, err error) {
	switch configSeal.Type {
	case seal.AliCloudKMS:
		return configureAliCloudKMSSeal(configSeal, infoKeys, info, logger, inseal)

	case seal.AWSKMS:
		return configureAWSKMSSeal(configSeal, infoKeys, info, logger, inseal)

	case seal.GCPCKMS:
		return configureGCPCKMSSeal(configSeal, infoKeys, info, logger, inseal)

	case seal.AzureKeyVault:
		return configureAzureKeyVaultSeal(configSeal, infoKeys, info, logger, inseal)

	case seal.Transit:
		return configureTransitSeal(configSeal, infoKeys, info, logger, inseal)

	case seal.PKCS11:
		return nil, fmt.Errorf("Seal type 'pkcs11' requires the Vault Enterprise HSM binary")

	case seal.Shamir:
		return inseal, nil

	default:
		return nil, fmt.Errorf("Unknown seal type %q", configSeal.Type)
	}
}
