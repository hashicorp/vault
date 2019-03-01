package vault

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"sync/atomic"

	proto "github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault/seal"
)

// barrierTypeUpgradeCheck checks for backwards compat on barrier type, not
// applicable in the OSS side
var barrierTypeUpgradeCheck = func(_ string, _ *SealConfig) {}

// autoSeal is a Seal implementation that contains logic for encrypting and
// decrypting stored keys via an underlying AutoSealAccess implementation, as
// well as logic related to recovery keys and barrier config.
type autoSeal struct {
	seal.Access

	barrierConfig  atomic.Value
	recoveryConfig atomic.Value
	core           *Core
}

// Ensure we are implementing the Seal interface
var _ Seal = (*autoSeal)(nil)

func NewAutoSeal(lowLevel seal.Access) Seal {
	ret := &autoSeal{
		Access: lowLevel,
	}
	ret.barrierConfig.Store((*SealConfig)(nil))
	ret.recoveryConfig.Store((*SealConfig)(nil))
	return ret
}

func (d *autoSeal) checkCore() error {
	if d.core == nil {
		return fmt.Errorf("seal does not have a core set")
	}
	return nil
}

func (d *autoSeal) SetCore(core *Core) {
	d.core = core
}

func (d *autoSeal) Init(ctx context.Context) error {
	return d.Access.Init(ctx)
}

func (d *autoSeal) Finalize(ctx context.Context) error {
	return d.Access.Finalize(ctx)
}

func (d *autoSeal) BarrierType() string {
	return d.SealType()
}

func (d *autoSeal) StoredKeysSupported() bool {
	return true
}

func (d *autoSeal) RecoveryKeySupported() bool {
	return true
}

// SetStoredKeys uses the autoSeal.Access.Encrypts method to wrap the keys. The stored entry
// does not need to be seal wrapped in this case.
func (d *autoSeal) SetStoredKeys(ctx context.Context, keys [][]byte) error {
	if keys == nil {
		return fmt.Errorf("keys were nil")
	}
	if len(keys) == 0 {
		return fmt.Errorf("no keys provided")
	}

	buf, err := json.Marshal(keys)
	if err != nil {
		return errwrap.Wrapf("failed to encode keys for storage: {{err}}", err)
	}

	// Encrypt and marshal the keys
	blobInfo, err := d.Encrypt(ctx, buf)
	if err != nil {
		return errwrap.Wrapf("failed to encrypt keys for storage: {{err}}", err)
	}

	value, err := proto.Marshal(blobInfo)
	if err != nil {
		return errwrap.Wrapf("failed to marshal value for storage: {{err}}", err)
	}

	// Store the seal configuration.
	pe := &physical.Entry{
		Key:   StoredBarrierKeysPath,
		Value: value,
	}

	if err := d.core.physical.Put(ctx, pe); err != nil {
		return errwrap.Wrapf("failed to write keys to storage: {{err}}", err)
	}

	return nil
}

// GetStoredKeys retrieves the key shares by unwrapping the encrypted key using the
// autoseal.
func (d *autoSeal) GetStoredKeys(ctx context.Context) ([][]byte, error) {
	pe, err := d.core.physical.Get(ctx, StoredBarrierKeysPath)
	if err != nil {
		return nil, errwrap.Wrapf("failed to fetch stored keys: {{err}}", err)
	}

	// This is not strictly an error; we may not have any stored keys, for
	// instance, if we're not initialized
	if pe == nil {
		return nil, nil
	}

	blobInfo := &physical.EncryptedBlobInfo{}
	if err := proto.Unmarshal(pe.Value, blobInfo); err != nil {
		return nil, errwrap.Wrapf("failed to proto decode stored keys: {{err}}", err)
	}

	pt, err := d.Decrypt(ctx, blobInfo)
	if err != nil {
		return nil, errwrap.Wrapf("failed to decrypt encrypted stored keys: {{err}}", err)
	}

	// Decode the barrier entry
	var keys [][]byte
	if err := json.Unmarshal(pt, &keys); err != nil {
		return nil, fmt.Errorf("failed to decode stored keys: %v", err)
	}

	return keys, nil
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
		d.core.logger.Error("autoseal: failed to read seal configuration", "seal_type", sealType, "error", err)
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to read %q seal configuration: {{err}}", sealType), err)
	}

	// If the seal configuration is missing, we are not initialized
	if entry == nil {
		if d.core.logger.IsInfo() {
			d.core.logger.Info("autoseal: seal configuration missing, not initialized", "seal_type", sealType)
		}
		return nil, nil
	}

	conf := &SealConfig{}
	err = json.Unmarshal(entry.Value, conf)
	if err != nil {
		d.core.logger.Error("autoseal: failed to decode seal configuration", "seal_type", sealType, "error", err)
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to decode %q seal configuration: {{err}}", sealType), err)
	}

	// Check for a valid seal configuration
	if err := conf.Validate(); err != nil {
		d.core.logger.Error("autoseal: invalid seal configuration", "seal_type", sealType, "error", err)
		return nil, errwrap.Wrapf(fmt.Sprintf("%q seal validation failed: {{err}}", sealType), err)
	}

	barrierTypeUpgradeCheck(d.BarrierType(), conf)

	if conf.Type != d.BarrierType() {
		d.core.logger.Error("autoseal: barrier seal type does not match loaded type", "seal_type", conf.Type, "loaded_type", d.BarrierType())
		return nil, fmt.Errorf("barrier seal type of %q does not match loaded type of %q", conf.Type, d.BarrierType())
	}

	d.barrierConfig.Store(conf)
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
		return errwrap.Wrapf("failed to encode barrier seal configuration: {{err}}", err)
	}

	// Store the seal configuration
	pe := &physical.Entry{
		Key:   barrierSealConfigPath,
		Value: buf,
	}

	if err := d.core.physical.Put(ctx, pe); err != nil {
		d.core.logger.Error("autoseal: failed to write barrier seal configuration", "error", err)
		return errwrap.Wrapf("failed to write barrier seal configuration: {{err}}", err)
	}

	d.barrierConfig.Store(conf.Clone())

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
		d.core.logger.Error("autoseal: failed to read seal configuration", "seal_type", sealType, "error", err)
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to read %q seal configuration: {{err}}", sealType), err)
	}

	if entry == nil {
		if d.core.Sealed() {
			d.core.logger.Info("autoseal: seal configuration missing, but cannot check old path as core is sealed", "seal_type", sealType)
			return nil, nil
		}

		// Check the old recovery seal config path so an upgraded standby will
		// return the correct seal config
		be, err := d.core.barrier.Get(ctx, recoverySealConfigPath)
		if err != nil {
			return nil, errwrap.Wrapf("failed to read old recovery seal configuration: {{err}}", err)
		}

		// If the seal configuration is missing, then we are not initialized.
		if be == nil {
			if d.core.logger.IsInfo() {
				d.core.logger.Info("autoseal: seal configuration missing, not initialized", "seal_type", sealType)
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
		d.core.logger.Error("autoseal: failed to decode seal configuration", "seal_type", sealType, "error", err)
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to decode %q seal configuration: {{err}}", sealType), err)
	}

	// Check for a valid seal configuration
	if err := conf.Validate(); err != nil {
		d.core.logger.Error("autoseal: invalid seal configuration", "seal_type", sealType, "error", err)
		return nil, errwrap.Wrapf(fmt.Sprintf("%q seal validation failed: {{err}}", sealType), err)
	}

	if conf.Type != d.RecoveryType() {
		d.core.logger.Error("autoseal: recovery seal type does not match loaded type", "seal_type", conf.Type, "loaded_type", d.RecoveryType())
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
		return errwrap.Wrapf("failed to encode recovery seal configuration: {{err}}", err)
	}

	// Store the seal configuration directly in the physical storage
	pe := &physical.Entry{
		Key:   recoverySealConfigPlaintextPath,
		Value: buf,
	}

	if err := d.core.physical.Put(ctx, pe); err != nil {
		d.core.logger.Error("autoseal: failed to write recovery seal configuration", "error", err)
		return errwrap.Wrapf("failed to write recovery seal configuration: {{err}}", err)
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
	blobInfo, err := d.Encrypt(ctx, key)
	if err != nil {
		return errwrap.Wrapf("failed to encrypt keys for storage: {{err}}", err)
	}

	value, err := proto.Marshal(blobInfo)
	if err != nil {
		return errwrap.Wrapf("failed to marshal value for storage: {{err}}", err)
	}

	be := &physical.Entry{
		Key:   recoveryKeyPath,
		Value: value,
	}

	if err := d.core.physical.Put(ctx, be); err != nil {
		d.core.logger.Error("autoseal: failed to write recovery key", "error", err)
		return errwrap.Wrapf("failed to write recovery key: {{err}}", err)
	}

	return nil
}

func (d *autoSeal) RecoveryKey(ctx context.Context) ([]byte, error) {
	return d.getRecoveryKeyInternal(ctx)
}

func (d *autoSeal) getRecoveryKeyInternal(ctx context.Context) ([]byte, error) {
	pe, err := d.core.physical.Get(ctx, recoveryKeyPath)
	if err != nil {
		d.core.logger.Error("autoseal: failed to read recovery key", "error", err)
		return nil, errwrap.Wrapf("failed to read recovery key: {{err}}", err)
	}
	if pe == nil {
		d.core.logger.Warn("autoseal: no recovery key found")
		return nil, fmt.Errorf("no recovery key found")
	}

	blobInfo := &physical.EncryptedBlobInfo{}
	if err := proto.Unmarshal(pe.Value, blobInfo); err != nil {
		return nil, errwrap.Wrapf("failed to proto decode stored keys: {{err}}", err)
	}

	pt, err := d.Decrypt(ctx, blobInfo)
	if err != nil {
		return nil, errwrap.Wrapf("failed to decrypt encrypted stored keys: {{err}}", err)
	}

	return pt, nil
}

// migrateRecoveryConfig is a helper func to migrate the recovery config to
// live outside the barrier. This is called from SetRecoveryConfig which is
// always called with the stateLock.
func (d *autoSeal) migrateRecoveryConfig(ctx context.Context) error {
	// Get config from the old recoverySealConfigPath path
	be, err := d.core.barrier.Get(ctx, recoverySealConfigPath)
	if err != nil {
		return errwrap.Wrapf("failed to read old recovery seal configuration during migration: {{err}}", err)
	}

	// If this entry is nil, then skip migration
	if be == nil {
		return nil
	}

	// Only log if we are performing the migration
	d.core.logger.Debug("migrating recovery seal configuration")
	defer d.core.logger.Debug("done migrating recovery seal configuration")

	// Perform migration
	pe := &physical.Entry{
		Key:   recoverySealConfigPlaintextPath,
		Value: be.Value,
	}

	if err := d.core.physical.Put(ctx, pe); err != nil {
		return errwrap.Wrapf("failed to write recovery seal configuration during migration: {{err}}", err)
	}

	// Perform deletion of the old entry
	if err := d.core.barrier.Delete(ctx, recoverySealConfigPath); err != nil {
		return errwrap.Wrapf("failed to delete old recovery seal configuration during migration: {{err}}", err)
	}

	return nil
}
