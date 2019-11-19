package seal

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/vault"
)

var (
	ConfigureSeal = configureSeal
)

func configureSeal(configSeal *server.Seal, infoKeys *[]string, info *map[string]string, logger hclog.Logger, inseal vault.Seal) (outseal vault.Seal, err error) {
	switch configSeal.Type {
	case wrapping.AliCloudKMS:
		return configureAliCloudKMSSeal(configSeal, infoKeys, info, logger, inseal)

	case wrapping.AWSKMS:
		return configureAWSKMSSeal(configSeal, infoKeys, info, logger, inseal)

	case wrapping.AzureKeyVault:
		return configureAzureKeyVaultSeal(configSeal, infoKeys, info, logger, inseal)

		/*
			case seal.GCPCKMS:
				return configureGCPCKMSSeal(configSeal, infoKeys, info, logger, inseal)

			case seal.OCIKMS:
				return configureOCIKMSSeal(configSeal, infoKeys, info, logger, inseal)

		*/
	case wrapping.Transit:
		return configureTransitSeal(configSeal, infoKeys, info, logger, inseal)

		/*
			case seal.PKCS11:
				return nil, fmt.Errorf("Seal type 'pkcs11' requires the Vault Enterprise HSM binary")

		*/
	case wrapping.Shamir:
		return inseal, nil

	default:
		return nil, fmt.Errorf("Unknown seal type %q", configSeal.Type)
	}
}
