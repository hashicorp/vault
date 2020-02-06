package command

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"

	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/vault"
	vaultseal "github.com/hashicorp/vault/vault/seal"
	"github.com/pkg/errors"
)

var (
	onEnterprise                 = false
	createSecureRandomReaderFunc = createSecureRandomReader
)

func createSecureRandomReader(config *server.Config, seal *vault.Seal) (io.Reader, error) {
	return rand.Reader, nil
}

func adjustCoreForSealMigration(logger log.Logger, core *vault.Core, barrierSeal, unwrapSeal vault.Seal) error {
	existBarrierSealConfig, existRecoverySealConfig, err := core.PhysicalSealConfigs(context.Background())
	if err != nil {
		return fmt.Errorf("Error checking for existing seal: %s", err)
	}

	// If we don't have an existing config or if it's the deprecated auto seal
	// which needs an upgrade, skip out
	if existBarrierSealConfig == nil || existBarrierSealConfig.Type == wrapping.HSMAutoDeprecated {
		return nil
	}

	if unwrapSeal == nil {
		// We have the same barrier type and the unwrap seal is nil so we're not
		// migrating from same to same, IOW we assume it's not a migration
		if existBarrierSealConfig.Type == barrierSeal.BarrierType() {
			return nil
		}

		// If we're not coming from Shamir, and the existing type doesn't match
		// the barrier type, we need both the migration seal and the new seal
		if existBarrierSealConfig.Type != wrapping.Shamir && barrierSeal.BarrierType() != wrapping.Shamir {
			return errors.New(`Trying to migrate from auto-seal to auto-seal but no "disabled" seal stanza found`)
		}
	} else {
		if unwrapSeal.BarrierType() == wrapping.Shamir {
			return errors.New("Shamir seals cannot be set disabled (they should simply not be set)")
		}
	}

	if existBarrierSealConfig.Type == barrierSeal.BarrierType() {
		// In this case our migration seal is set so we are using it
		// (potentially) for unwrapping. Set it on core for that purpose then
		// exit.
		core.SetSealsForMigration(nil, nil, unwrapSeal)
		return nil
	}

	if existBarrierSealConfig.Type != wrapping.Shamir && existRecoverySealConfig == nil {
		return errors.New(`Recovery seal configuration not found for existing seal`)
	}

	if onEnterprise && barrierSeal.BarrierType() == wrapping.Shamir {
		return errors.New("Migrating from autoseal to Shamir seal is not currently supported on Vault Enterprise")
	}

	var migrationSeal vault.Seal
	var newSeal vault.Seal

	// Determine the migrationSeal. This is either going to be an instance of
	// shamir or the unwrapSeal.
	switch existBarrierSealConfig.Type {
	case wrapping.Shamir:
		// The value reflected in config is what we're going to
		migrationSeal = vault.NewDefaultSeal(&vaultseal.Access{
			Wrapper: aeadwrapper.NewWrapper(&wrapping.WrapperOptions{
				Logger: logger.Named("shamir"),
			}),
		})

	default:
		// If we're not coming from Shamir we expect the previous seal to be
		// in the config and disabled.
		migrationSeal = unwrapSeal
	}

	// newSeal will be the barrierSeal
	newSeal = barrierSeal

	// Set the appropriate barrier and recovery configs.
	switch {
	case migrationSeal.RecoveryKeySupported() && newSeal.RecoveryKeySupported():
		// Migrating from auto->auto, copy the configs over
		newSeal.SetCachedBarrierConfig(existBarrierSealConfig)
		newSeal.SetCachedRecoveryConfig(existRecoverySealConfig)
	case migrationSeal.RecoveryKeySupported():
		// Migrating from auto->shamir, clone auto's recovery config and set
		// stored keys to 1.
		newSealConfig := existRecoverySealConfig.Clone()
		newSealConfig.StoredShares = 1
		newSeal.SetCachedBarrierConfig(newSealConfig)
	case newSeal.RecoveryKeySupported():
		// Migrating from shamir->auto, set a new barrier config and set
		// recovery config to a clone of shamir's barrier config with stored
		// keys set to 0.
		newBarrierSealConfig := &vault.SealConfig{
			Type:            newSeal.BarrierType(),
			SecretShares:    1,
			SecretThreshold: 1,
			StoredShares:    1,
		}
		newSeal.SetCachedBarrierConfig(newBarrierSealConfig)

		newRecoveryConfig := existBarrierSealConfig.Clone()
		newRecoveryConfig.StoredShares = 0
		newSeal.SetCachedRecoveryConfig(newRecoveryConfig)
	}

	core.SetSealsForMigration(migrationSeal, newSeal, unwrapSeal)

	return nil
}
