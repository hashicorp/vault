package vault

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"sync/atomic"

	proto "github.com/golang/protobuf/proto"
	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
)

// barrierTypeUpgradeCheck checks for backwards compat on barrier type, not
// applicable in the OSS side
var barrierTypeUpgradeCheck = func(_ string, _ *SealConfig) {}

// autoSeal is a Seal implementation that contains logic for encrypting and
// decrypting stored keys via an underlying AutoSealAccess implementation, as
// well as logic related to recovery keys and barrier config.
type autoSeal struct {
	*seal.Access

	barrierConfig  atomic.Value
	recoveryConfig atomic.Value
	core           *Core
	logger         log.Logger
}

// Ensure we are implementing the Seal interface
var _ Seal = (*autoSeal)(nil)

func NewAutoSeal(lowLevel *seal.Access) *autoSeal {
	ret := &autoSeal{
		Access: lowLevel,
	}
	ret.barrierConfig.Store((*SealConfig)(nil))
	ret.recoveryConfig.Store((*SealConfig)(nil))
	return ret
}

func (d *autoSeal) SealWrapable() bool {
	return true
}

func (d *autoSeal) GetAccess() *seal.Access {
	return d.Access
}

func (d *autoSeal) checkCore() error {
	if d.core == nil {
		return fmt.Errorf("seal does not have a core set")
	}
	return nil
}

func (d *autoSeal) SetCore(core *Core) {
	d.core = core
	if d.logger == nil {
		d.logger = d.core.Logger().Named("autoseal")
		d.core.AddLogger(d.logger)
	}
}

func (d *autoSeal) Init(ctx context.Context) error {
	return d.Access.Init(ctx)
}

func (d *autoSeal) Finalize(ctx context.Context) error {
	return d.Access.Finalize(ctx)
}

func (d *autoSeal) BarrierType() string {
	return d.Type()
}

func (d *autoSeal) StoredKeysSupported() seal.StoredKeysSupport {
	return seal.StoredKeysSupportedGeneric
}

func (d *autoSeal) RecoveryKeySupported() bool {
	return true
}

// SetStoredKeys uses the autoSeal.Access.Encrypts method to wrap the keys. The stored entry
// does not need to be seal wrapped in this case.
func (d *autoSeal) SetStoredKeys(ctx context.Context, keys [][]byte) error {
	return writeStoredKeys(ctx, d.core.physical, d.Access, keys)
}

// GetStoredKeys retrieves the key shares by unwrapping the encrypted key using the
// autoseal.
func (d *autoSeal) GetStoredKeys(ctx context.Context) ([][]byte, error) {
	return readStoredKeys(ctx, d.core.physical, d.Access)
}

func (d *autoSeal) upgradeStoredKeys(ctx context.Context) error {
	pe, err := d.core.physical.Get(ctx, StoredBarrierKeysPath)
	if err != nil {
		return fmt.Errorf("failed to fetch stored keys: %w", err)
	}
	if pe == nil {
		return fmt.Errorf("no stored keys found")
	}

	blobInfo := &wrapping.EncryptedBlobInfo{}
	if err := proto.Unmarshal(pe.Value, blobInfo); err != nil {
		return fmt.Errorf("failed to proto decode stored keys: %w", err)
	}

	if blobInfo.KeyInfo != nil && blobInfo.KeyInfo.KeyID != d.Access.KeyID() {
		d.logger.Info("upgrading stored keys")

		pt, err := d.Decrypt(ctx, blobInfo, nil)
		if err != nil {
			return fmt.Errorf("failed to decrypt encrypted stored keys: %w", err)
		}

		// Decode the barrier entry
		var keys [][]byte
		if err := json.Unmarshal(pt, &keys); err != nil {
			return fmt.Errorf("failed to decode stored keys: %w", err)
		}

		if err := d.SetStoredKeys(ctx, keys); err != nil {
			return fmt.Errorf("failed to save upgraded stored keys: %w", err)
		}
	}
	return nil
}

// UpgradeKeys re-encrypts and saves the stored keys and the recovery key
// with the current key if the current KeyID is different from the KeyID
// the stored keys and the recovery key are encrypted with. The provided
// Context must be non-nil.
func (d *autoSeal) UpgradeKeys(ctx context.Context) error {
	// Many of the seals update their keys to the latest KeyID when Encrypt
	// is called.
	if _, err := d.Encrypt(ctx, []byte("a"), nil); err != nil {
		return err
	}

	if err := d.upgradeRecoveryKey(ctx); err != nil {
		return err
	}
	if err := d.upgradeStoredKeys(ctx); err != nil {
		return err
	}
	return nil
}

func (d *autoSeal) BarrierConfig(ctx context.Context) (*SealConfig, error) {
	if d.barrierConfig.Load().(*SealConfig) != nil {
		return d.barrierConfig.Load().(*SealConfig).Clone(), nil
	}

	if err := d.checkCore(); err != nil {
		return nil, err
	}

	sealType := "barrier"

	entry, err := d.core.physical.Get(ctx, barrierSealConfigPath)
	if err != nil {
		d.logger.Error("failed to read seal configuration", "seal_type", sealType, "error", err)
		return nil, fmt.Errorf("failed to read %q seal configuration: %w", sealType, err)
	}

	// If the seal configuration is missing, we are not initialized
	if entry == nil {
		if d.logger.IsInfo() {
			d.logger.Info("seal configuration missing, not initialized", "seal_type", sealType)
		}
		return nil, nil
	}

	conf := &SealConfig{}
	err = json.Unmarshal(entry.Value, conf)
	if err != nil {
		d.logger.Error("failed to decode seal configuration", "seal_type", sealType, "error", err)
		return nil, fmt.Errorf("failed to decode %q seal configuration: %w", sealType, err)
	}

	// Check for a valid seal configuration
	if err := conf.Validate(); err != nil {
		d.logger.Error("invalid seal configuration", "seal_type", sealType, "error", err)
		return nil, fmt.Errorf("%q seal validation failed: %w", sealType, err)
	}

	barrierTypeUpgradeCheck(d.BarrierType(), conf)

	if conf.Type != d.BarrierType() {
		d.logger.Error("barrier seal type does not match loaded type", "seal_type", conf.Type, "loaded_type", d.BarrierType())
		return nil, fmt.Errorf("barrier seal type of %q does not match loaded type of %q", conf.Type, d.BarrierType())
	}

	d.SetCachedBarrierConfig(conf)
	return conf.Clone(), nil
}

func (d *autoSeal) SetBarrierConfig(ctx context.Context, conf *SealConfig) error {
	if err := d.checkCore(); err != nil {
		return err
	}

	if conf == nil {
		d.barrierConfig.Store((*SealConfig)(nil))
		return nil
	}

	conf.Type = d.BarrierType()

	// Encode the seal configuration
	buf, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("failed to encode barrier seal configuration: %w", err)
	}

	// Store the seal configuration
	pe := &physical.Entry{
		Key:   barrierSealConfigPath,
		Value: buf,
	}

	if err := d.core.physical.Put(ctx, pe); err != nil {
		d.logger.Error("failed to write barrier seal configuration", "error", err)
		return fmt.Errorf("failed to write barrier seal configuration: %w", err)
	}

	d.SetCachedBarrierConfig(conf.Clone())

	return nil
}

func (d *autoSeal) SetCachedBarrierConfig(config *SealConfig) {
	d.barrierConfig.Store(config)
}

func (d *autoSeal) RecoveryType() string {
	return RecoveryTypeShamir
}

// RecoveryConfig returns the recovery config on recoverySealConfigPlaintextPath.
func (d *autoSeal) RecoveryConfig(ctx context.Context) (*SealConfig, error) {
	if d.recoveryConfig.Load().(*SealConfig) != nil {
		return d.recoveryConfig.Load().(*SealConfig).Clone(), nil
	}

	if err := d.checkCore(); err != nil {
		return nil, err
	}

	sealType := "recovery"

	var entry *physical.Entry
	var err error
	entry, err = d.core.physical.Get(ctx, recoverySealConfigPlaintextPath)
	if err != nil {
		d.logger.Error("failed to read seal configuration", "seal_type", sealType, "error", err)
		return nil, fmt.Errorf("failed to read %q seal configuration: %w", sealType, err)
	}

	if entry == nil {
		if d.core.Sealed() {
			d.logger.Info("seal configuration missing, but cannot check old path as core is sealed", "seal_type", sealType)
			return nil, nil
		}

		// Check the old recovery seal config path so an upgraded standby will
		// return the correct seal config
		be, err := d.core.barrier.Get(ctx, recoverySealConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read old recovery seal configuration: %w", err)
		}

		// If the seal configuration is missing, then we are not initialized.
		if be == nil {
			if d.logger.IsInfo() {
				d.logger.Info("seal configuration missing, not initialized", "seal_type", sealType)
			}
			return nil, nil
		}

		// Reconstruct the physical entry
		entry = &physical.Entry{
			Key:   be.Key,
			Value: be.Value,
		}
	}

	conf := &SealConfig{}
	if err := json.Unmarshal(entry.Value, conf); err != nil {
		d.logger.Error("failed to decode seal configuration", "seal_type", sealType, "error", err)
		return nil, fmt.Errorf("failed to decode %q seal configuration: %w", sealType, err)
	}

	// Check for a valid seal configuration
	if err := conf.Validate(); err != nil {
		d.logger.Error("invalid seal configuration", "seal_type", sealType, "error", err)
		return nil, fmt.Errorf("%q seal validation failed: %w", sealType, err)
	}

	if conf.Type != d.RecoveryType() {
		d.logger.Error("recovery seal type does not match loaded type", "seal_type", conf.Type, "loaded_type", d.RecoveryType())
		return nil, fmt.Errorf("recovery seal type of %q does not match loaded type of %q", conf.Type, d.RecoveryType())
	}

	d.recoveryConfig.Store(conf)
	return conf.Clone(), nil
}

// SetRecoveryConfig writes the recovery configuration to the physical storage
// and sets it as the seal's recoveryConfig.
func (d *autoSeal) SetRecoveryConfig(ctx context.Context, conf *SealConfig) error {
	if err := d.checkCore(); err != nil {
		return err
	}

	// Perform migration if applicable
	if err := d.migrateRecoveryConfig(ctx); err != nil {
		return err
	}

	if conf == nil {
		d.recoveryConfig.Store((*SealConfig)(nil))
		return nil
	}

	conf.Type = d.RecoveryType()

	// Encode the seal configuration
	buf, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("failed to encode recovery seal configuration: %w", err)
	}

	// Store the seal configuration directly in the physical storage
	pe := &physical.Entry{
		Key:   recoverySealConfigPlaintextPath,
		Value: buf,
	}

	if err := d.core.physical.Put(ctx, pe); err != nil {
		d.logger.Error("failed to write recovery seal configuration", "error", err)
		return fmt.Errorf("failed to write recovery seal configuration: %w", err)
	}

	d.recoveryConfig.Store(conf.Clone())

	return nil
}

func (d *autoSeal) SetCachedRecoveryConfig(config *SealConfig) {
	d.recoveryConfig.Store(config)
}

func (d *autoSeal) VerifyRecoveryKey(ctx context.Context, key []byte) error {
	if key == nil {
		return fmt.Errorf("recovery key to verify is nil")
	}

	pt, err := d.getRecoveryKeyInternal(ctx)
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare(key, pt) != 1 {
		return fmt.Errorf("recovery key does not match submitted values")
	}

	return nil
}

func (d *autoSeal) SetRecoveryKey(ctx context.Context, key []byte) error {
	if err := d.checkCore(); err != nil {
		return err
	}

	if key == nil {
		return fmt.Errorf("recovery key to store is nil")
	}

	// Encrypt and marshal the keys
	blobInfo, err := d.Encrypt(ctx, key, nil)
	if err != nil {
		return fmt.Errorf("failed to encrypt keys for storage: %w", err)
	}

	value, err := proto.Marshal(blobInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal value for storage: %w", err)
	}

	be := &physical.Entry{
		Key:   recoveryKeyPath,
		Value: value,
	}

	if err := d.core.physical.Put(ctx, be); err != nil {
		d.logger.Error("failed to write recovery key", "error", err)
		return fmt.Errorf("failed to write recovery key: %w", err)
	}

	return nil
}

func (d *autoSeal) RecoveryKey(ctx context.Context) ([]byte, error) {
	return d.getRecoveryKeyInternal(ctx)
}

func (d *autoSeal) getRecoveryKeyInternal(ctx context.Context) ([]byte, error) {
	pe, err := d.core.physical.Get(ctx, recoveryKeyPath)
	if err != nil {
		d.logger.Error("failed to read recovery key", "error", err)
		return nil, fmt.Errorf("failed to read recovery key: %w", err)
	}
	if pe == nil {
		d.logger.Warn("no recovery key found")
		return nil, fmt.Errorf("no recovery key found")
	}

	blobInfo := &wrapping.EncryptedBlobInfo{}
	if err := proto.Unmarshal(pe.Value, blobInfo); err != nil {
		return nil, fmt.Errorf("failed to proto decode stored keys: %w", err)
	}

	pt, err := d.Decrypt(ctx, blobInfo, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt encrypted stored keys: %w", err)
	}

	return pt, nil
}

func (d *autoSeal) upgradeRecoveryKey(ctx context.Context) error {
	pe, err := d.core.physical.Get(ctx, recoveryKeyPath)
	if err != nil {
		return fmt.Errorf("failed to fetch recovery key: %w", err)
	}
	if pe == nil {
		return fmt.Errorf("no recovery key found")
	}

	blobInfo := &wrapping.EncryptedBlobInfo{}
	if err := proto.Unmarshal(pe.Value, blobInfo); err != nil {
		return fmt.Errorf("failed to proto decode recovery key: %w", err)
	}

	if blobInfo.KeyInfo != nil && blobInfo.KeyInfo.KeyID != d.Access.KeyID() {
		d.logger.Info("upgrading recovery key")

		pt, err := d.Decrypt(ctx, blobInfo, nil)
		if err != nil {
			return fmt.Errorf("failed to decrypt encrypted recovery key: %w", err)
		}
		if err := d.SetRecoveryKey(ctx, pt); err != nil {
			return fmt.Errorf("failed to save upgraded recovery key: %w", err)
		}
	}
	return nil
}

// migrateRecoveryConfig is a helper func to migrate the recovery config to
// live outside the barrier. This is called from SetRecoveryConfig which is
// always called with the stateLock.
func (d *autoSeal) migrateRecoveryConfig(ctx context.Context) error {
	// Get config from the old recoverySealConfigPath path
	be, err := d.core.barrier.Get(ctx, recoverySealConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read old recovery seal configuration during migration: %w", err)
	}

	// If this entry is nil, then skip migration
	if be == nil {
		return nil
	}

	// Only log if we are performing the migration
	d.logger.Debug("migrating recovery seal configuration")
	defer d.logger.Debug("done migrating recovery seal configuration")

	// Perform migration
	pe := &physical.Entry{
		Key:   recoverySealConfigPlaintextPath,
		Value: be.Value,
	}

	if err := d.core.physical.Put(ctx, pe); err != nil {
		return fmt.Errorf("failed to write recovery seal configuration during migration: %w", err)
	}

	// Perform deletion of the old entry
	if err := d.core.barrier.Delete(ctx, recoverySealConfigPath); err != nil {
		return fmt.Errorf("failed to delete old recovery seal configuration during migration: %w", err)
	}

	return nil
}
