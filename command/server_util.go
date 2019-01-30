package command

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/vault"
	vaultseal "github.com/hashicorp/vault/vault/seal"
	"github.com/pkg/errors"
)

var (
	onEnterprise = false
)

func adjustCoreForSealMigration(ctx context.Context, core *vault.Core, coreConfig *vault.CoreConfig, seal vault.Seal, config *server.Config) error {
	existBarrierSealConfig, existRecoverySealConfig, err := core.PhysicalSealConfigs(context.Background())
	if err != nil {
		return fmt.Errorf("Error checking for existing seal: %s", err)
	}
	var existSeal vault.Seal
	var newSeal vault.Seal
	if existBarrierSealConfig != nil && existBarrierSealConfig.Type != vaultseal.HSMAutoDeprecated &&
		(existBarrierSealConfig.Type != seal.BarrierType() ||
			config.Seal != nil && config.Seal.Disabled) {
		// If the existing seal is not Shamir, we're going to Shamir, which
		// means we require them setting "disabled" to true in their
		// configuration as a sanity check.
		if (config.Seal == nil || !config.Seal.Disabled) && existBarrierSealConfig.Type != vaultseal.Shamir {
			return errors.New(`Seal migration requires specifying "disabled" as "true" in the "seal" block of Vault's configuration file"`)
		}

		// Conversely, if they are going from Shamir to auto, we want to
		// ensure disabled is *not* set
		if existBarrierSealConfig.Type == vaultseal.Shamir && config.Seal != nil && config.Seal.Disabled {
			coreConfig.Logger.Warn(`when not migrating, Vault's config should not specify "disabled" as "true" in the "seal" block of Vault's configuration file`)
			return nil
		}

		if existBarrierSealConfig.Type != vaultseal.Shamir && existRecoverySealConfig == nil {
			return errors.New(`Recovery seal configuration not found for existing seal`)
		}

		switch existBarrierSealConfig.Type {
		case vaultseal.Shamir:
			// The value reflected in config is what we're going to
			existSeal = vault.NewDefaultSeal()
			existSeal.SetCore(core)
			newSeal = seal
			newBarrierSealConfig := &vault.SealConfig{
				Type:            newSeal.BarrierType(),
				SecretShares:    1,
				SecretThreshold: 1,
				StoredShares:    1,
			}
			newSeal.SetCachedBarrierConfig(newBarrierSealConfig)
			newSeal.SetCachedRecoveryConfig(existBarrierSealConfig)

		default:
			if onEnterprise {
				return errors.New("Migrating from autoseal to Shamir seal is not supported on Vault Enterprise")
			}

			// The disabled value reflected in config is what we're going from
			existSeal = coreConfig.Seal
			newSeal = vault.NewDefaultSeal()
			newSeal.SetCore(core)
			newSeal.SetCachedBarrierConfig(existRecoverySealConfig)
		}

		core.SetSealsForMigration(existSeal, newSeal)
	}

	return nil
}
