// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"crypto/subtle"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
)

// barrierTypeUpgradeCheck checks for backwards compat on barrier type, not
// applicable in the OSS side
var (
	barrierTypeUpgradeCheck     = func(_ SealConfigType, _ *SealConfig) {}
	autoSealUnavailableDuration = []string{"seal", "unreachable", "time"}
)

// autoSeal is a Seal implementation that contains logic for encrypting and
// decrypting stored keys via an underlying AutoSealAccess implementation, as
// well as logic related to recovery keys and barrier config.
type autoSeal struct {
	seal.Access

	barrierSealConfigType SealConfigType
	barrierConfig         atomic.Value
	recoveryConfig        atomic.Value
	core                  *Core
	logger                log.Logger

	allSealsHealthy bool
	hcLock          sync.RWMutex
	healthCheckStop chan struct{}
}

// Ensure we are implementing the Seal interface
var _ Seal = (*autoSeal)(nil)

func NewAutoSeal(lowLevel seal.Access) *autoSeal {
	ret := &autoSeal{
		Access: lowLevel,
	}
	ret.barrierConfig.Store((*SealConfig)(nil))
	ret.recoveryConfig.Store((*SealConfig)(nil))

	// See SealConfigType for the rules about computing the type.
	if len(lowLevel.GetSealGenerationInfo().Seals) > 1 {
		ret.barrierSealConfigType = SealConfigTypeMultiseal
	} else {
		// Note that the Access constructors guarantee that there is at least one KMS config
		ret.barrierSealConfigType = SealConfigType(lowLevel.GetSealGenerationInfo().Seals[0].Type)
	}

	return ret
}

func (d *autoSeal) Healthy() bool {
	d.hcLock.RLock()
	defer d.hcLock.RUnlock()
	return d.allSealsHealthy
}

func (d *autoSeal) SealWrapable() bool {
	return true
}

func (d *autoSeal) GetAccess() seal.Access {
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
	}
}

func (d *autoSeal) Init(ctx context.Context) error {
	return d.Access.Init(ctx)
}

func (d *autoSeal) Finalize(ctx context.Context) error {
	return d.Access.Finalize(ctx)
}

func (d *autoSeal) BarrierSealConfigType() SealConfigType {
	return d.barrierSealConfigType
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

	wrappedEntryValue, err := UnmarshalSealWrappedValue(pe.Value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal stored keys: %w", err)
	}
	uptodate, err := d.Access.IsUpToDate(ctx, wrappedEntryValue.getValue(), true)
	if err != nil {
		return fmt.Errorf("failed to check if stored keys are up-to-date: %w", err)
	}
	if !uptodate {
		d.logger.Info("upgrading stored keys")

		keys, err := UnsealWrapStoredBarrierKeys(ctx, d.GetAccess(), pe)
		if err != nil {
			return fmt.Errorf("failed to decrypt encrypted stored keys: %w", err)
		}

		if err := d.SetStoredKeys(ctx, keys); err != nil {
			return fmt.Errorf("failed to save upgraded stored keys: %w", err)
		}
	}
	return nil
}

// UpgradeKeys re-encrypts and saves the stored keys and the recovery key
// with the current key if the current KeyId is different from the KeyId
// the stored keys and the recovery key are encrypted with. The provided
// Context must be non-nil.
func (d *autoSeal) UpgradeKeys(ctx context.Context) error {
	if err := d.upgradeRecoveryKey(ctx); err != nil { // re-encrypts the recovery key
		return err
	}
	if err := d.upgradeStoredKeys(ctx); err != nil { // re-encrypts the root key
		return err
	}
	return nil
}

func (d *autoSeal) BarrierConfig(ctx context.Context) (*SealConfig, error) {
	if cfg := d.barrierConfig.Load().(*SealConfig); cfg != nil {
		return cfg.Clone(), nil
	}

	if err := d.checkCore(); err != nil {
		return nil, err
	}

	// Fetch the core configuration
	conf, err := d.core.PhysicalBarrierSealConfig(ctx)
	if err != nil {
		d.logger.Error("failed to read seal configuration", "error", err)
		return nil, fmt.Errorf("failed to read seal configuration: %w", err)
	}

	// If the seal configuration is missing, we are not initialized
	if conf == nil {
		d.logger.Info("seal configuration missing, not initialized")
		return nil, nil
	}

	barrierTypeUpgradeCheck(d.BarrierSealConfigType(), conf)

	if conf.Type != d.BarrierSealConfigType().String() {
		d.logger.Error("barrier seal type does not match loaded type", "seal_type", conf.Type, "loaded_type", d.BarrierSealConfigType())
		return nil, fmt.Errorf("barrier seal type of %q does not match loaded type of %q", conf.Type, d.BarrierSealConfigType())
	}

	d.SetCachedBarrierConfig(conf)
	return conf.Clone(), nil
}

func (d *autoSeal) ClearBarrierConfig(ctx context.Context) error {
	return d.SetBarrierConfig(ctx, nil)
}

func (d *autoSeal) SetBarrierConfig(ctx context.Context, conf *SealConfig) error {
	if err := d.checkCore(); err != nil {
		return err
	}

	if conf == nil {
		d.barrierConfig.Store((*SealConfig)(nil))
		return nil
	}

	conf.Type = d.BarrierSealConfigType().String()

	err := d.core.SetPhysicalBarrierSealConfig(ctx, conf)
	if err != nil {
		return err
	}

	d.SetCachedBarrierConfig(conf.Clone())

	return nil
}

func (d *autoSeal) SetCachedBarrierConfig(config *SealConfig) {
	d.barrierConfig.Store(config)
}

func (d *autoSeal) RecoverySealConfigType() SealConfigType {
	return SealConfigTypeRecovery
}

// RecoveryConfig returns the recovery config on recoverySealConfigPlaintextPath.
func (d *autoSeal) RecoveryConfig(ctx context.Context) (*SealConfig, error) {
	if cfg := d.recoveryConfig.Load().(*SealConfig); cfg != nil {
		return cfg.Clone(), nil
	}

	if err := d.checkCore(); err != nil {
		return nil, err
	}

	conf, err := d.core.PhysicalRecoverySealConfig(ctx)
	if err != nil {
		d.logger.Error("failed to read recovery seal configuration", "error", err)
		return nil, fmt.Errorf("failed to read recovery seal configuration: %w", err)
	}

	if conf == nil {
		if d.core.Sealed() {
			d.logger.Info("recovery seal configuration missing, but cannot check old path as core is sealed")
			return nil, nil
		}

		// Check the old recovery seal config path so an upgraded standby will
		// return the correct seal config
		conf, err := d.core.PhysicalRecoverySealConfigOldPath(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to read old recovery seal configuration: %w", err)
		}

		// If the seal configuration is missing, then we are not initialized.
		if conf == nil {
			d.logger.Info("recovery seal configuration missing, not initialized")
			return nil, nil
		}
	}

	if !d.RecoverySealConfigType().IsSameAs(conf.Type) {
		d.logger.Error("recovery seal type does not match loaded type", "seal_type", conf.Type, "loaded_type", d.RecoverySealConfigType())
		return nil, fmt.Errorf("recovery seal type of %q does not match loaded type of %q", conf.Type, d.RecoverySealConfigType())
	}

	d.recoveryConfig.Store(conf)
	return conf.Clone(), nil
}

func (d *autoSeal) ClearRecoveryConfig(ctx context.Context) error {
	return d.SetRecoveryConfig(ctx, nil)
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

	conf.Type = d.RecoverySealConfigType().String()

	if err := d.core.SetPhysicalRecoverySealConfig(ctx, conf); err != nil {
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
	be, err := SealWrapRecoveryKey(ctx, d.Access, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt keys for storage: %w", err)
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

	pt, err := UnsealWrapRecoveryKey(ctx, d.Access, pe)
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

	wrappedEntryValue, err := UnmarshalSealWrappedValue(pe.Value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal recovery key: %w", err)
	}
	uptodate, err := d.Access.IsUpToDate(ctx, wrappedEntryValue.getValue(), true)
	if err != nil {
		return fmt.Errorf("failed to check if recovery key is up-to-date: %w", err)
	}

	if !uptodate {
		d.logger.Info("upgrading recovery key")
		pt, err := UnsealWrapRecoveryKey(ctx, d.Access, pe)
		if err != nil {
			return fmt.Errorf("failed to decrypt recovery key: %w", err)
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

// StartHealthCheck starts a goroutine that tests the health of the auto-unseal backend once every 10 minutes.
// If unhealthy, logs a warning on the condition and begins testing every one minute until healthy again.
func (d *autoSeal) StartHealthCheck() {
	d.StopHealthCheck()
	d.hcLock.Lock()
	defer d.hcLock.Unlock()

	healthCheck := time.NewTicker(seal.HealthTestIntervalNominal)
	d.healthCheckStop = make(chan struct{})
	healthCheckStop := d.healthCheckStop
	ctx := d.core.activeContext

	go func() {
		lastTestOk := true
		lastSeenOk := time.Now()

		check := func(now time.Time) {
			ctx, cancel := context.WithTimeout(ctx, seal.HealthTestTimeout)
			defer cancel()

			allHealthy := true
			allUnhealthy := true
			for _, sealWrapper := range d.Access.GetAllSealWrappersByPriority() {
				mLabels := []metrics.Label{{Name: "seal_wrapper_name", Value: sealWrapper.Name}}

				wasHealthy := sealWrapper.IsHealthy()
				if err := sealWrapper.CheckHealth(ctx, now); err != nil {
					// Seal wrapper is unhealthy
					d.logger.Warn("seal wrapper health check failed", "seal_name", sealWrapper.Name, "err", err)
					d.core.MetricSink().SetGaugeWithLabels(autoSealUnavailableDuration,
						float32(time.Since(sealWrapper.LastSeenHealthy).Milliseconds()), mLabels)
					allHealthy = false
				} else {
					// Seal wrapper is healthy
					if wasHealthy {
						d.logger.Debug("seal wrapper health test passed", "seal_name", sealWrapper.Name)
					} else {
						d.logger.Info("seal wrapper is now healthy again", "downtime", "seal_name", sealWrapper.Name,
							now.Sub(sealWrapper.LastSeenHealthy).String())
					}
					allUnhealthy = false
				}
			}
			if allHealthy {
				if !lastTestOk {
					d.logger.Info("seal backend is fully healthy again", "downtime", now.Sub(lastSeenOk).String())
				}
				lastTestOk = true
				lastSeenOk = now
				healthCheck.Reset(seal.HealthTestIntervalNominal)
				d.core.MetricSink().SetGauge(autoSealUnavailableDuration, 0)
			} else {
				if lastTestOk && allUnhealthy {
					d.logger.Info("seal backend is completely unhealthy (all seal wrappers all unhealthy)", "downtime", now.Sub(lastSeenOk).String())
				}
				lastTestOk = false
				healthCheck.Reset(seal.HealthTestIntervalUnhealthy)
				d.core.MetricSink().SetGauge(autoSealUnavailableDuration, float32(time.Since(lastSeenOk).Milliseconds()))
			}

			d.hcLock.Lock()
			defer d.hcLock.Unlock()
			d.allSealsHealthy = allHealthy
		}

		for {
			select {
			case <-healthCheckStop:
				if healthCheck != nil {
					healthCheck.Stop()
				}
				healthCheckStop = nil
				return
			case t := <-healthCheck.C:
				check(t)
			}
		}
	}()
}

func (d *autoSeal) StopHealthCheck() {
	d.hcLock.Lock()
	defer d.hcLock.Unlock()
	if d.healthCheckStop != nil {
		close(d.healthCheckStop)
		d.healthCheckStop = nil
	}
}
