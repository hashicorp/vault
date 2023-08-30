// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/hashicorp/vault/command/server"

	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/physical"

	"github.com/hashicorp/vault/vault/seal"
)

const (
	// barrierSealConfigPath is the path used to store our seal configuration.
	// This value is stored in plaintext, since we must be able to read it even
	// with the Vault sealed. This is required so that we know how many secret
	// parts must be used to reconstruct the unseal key.
	barrierSealConfigPath = "core/seal-config"

	// recoverySealConfigPath is the path to the recovery key seal
	// configuration. It lives inside the barrier.
	// DEPRECATED: Use recoverySealConfigPlaintextPath instead.
	recoverySealConfigPath = "core/recovery-seal-config"

	// recoverySealConfigPlaintextPath is the path to the recovery key seal
	// configuration. This is stored in plaintext so that we can perform
	// auto-unseal.
	recoverySealConfigPlaintextPath = "core/recovery-config"

	// recoveryKeyPath is the path to the recovery key
	recoveryKeyPath = "core/recovery-key"

	// StoredBarrierKeysPath is the path used for storing HSM-encrypted unseal keys
	StoredBarrierKeysPath = "core/hsm/barrier-unseal-keys"

	// hsmStoredIVPath is the path to the initialization vector for stored keys
	hsmStoredIVPath = "core/hsm/iv"

	// SealGenInfoPath is the path used to store our seal generation info.
	// This is required so that we can detect any seal config changes and introduce
	// a new generation which keeps track if a rewrap of all CSPs and seal wrapped
	// values has completed .
	SealGenInfoPath = "core/seal-gen-info"
)

type Seal interface {
	SetCore(*Core)
	Init(context.Context) error
	Finalize(context.Context) error
	StoredKeysSupported() seal.StoredKeysSupport
	SealWrapable() bool
	SetStoredKeys(context.Context, [][]byte) error
	GetStoredKeys(context.Context) ([][]byte, error)
	BarrierSealConfigType() SealConfigType
	BarrierConfig(context.Context) (*SealConfig, error)
	ClearBarrierConfig(context.Context) error
	SetBarrierConfig(context.Context, *SealConfig) error
	SetCachedBarrierConfig(*SealConfig)

	RecoveryKeySupported() bool
	// RecoveryType returns the SealConfigType for the recovery seal. We only ever
	// expect this to be one of SealConfigTypeRecovery or SealConfigTypeRecoveryUnsupported
	RecoverySealConfigType() SealConfigType
	RecoveryConfig(context.Context) (*SealConfig, error)
	RecoveryKey(context.Context) ([]byte, error)
	ClearRecoveryConfig(context.Context) error
	SetRecoveryConfig(context.Context, *SealConfig) error
	SetCachedRecoveryConfig(*SealConfig)
	SetRecoveryKey(context.Context, []byte) error
	VerifyRecoveryKey(context.Context, []byte) error
	GetAccess() seal.Access
	Healthy() bool
}

type defaultSeal struct {
	access seal.Access
	config atomic.Value
	core   *Core
}

var _ Seal = (*defaultSeal)(nil)

func NewDefaultSeal(lowLevel seal.Access) Seal {
	ret := &defaultSeal{
		access: lowLevel,
	}
	ret.config.Store((*SealConfig)(nil))
	return ret
}

func (d *defaultSeal) SealWrapable() bool {
	return false
}

func (d *defaultSeal) checkCore() error {
	if d.core == nil {
		return fmt.Errorf("seal does not have a core set")
	}
	return nil
}

func (d *defaultSeal) GetAccess() seal.Access {
	return d.access
}

func (d *defaultSeal) SetAccess(access seal.Access) {
	d.access = access
}

func (d *defaultSeal) SetCore(core *Core) {
	d.core = core
}

func (d *defaultSeal) Init(ctx context.Context) error {
	return nil
}

func (d *defaultSeal) Finalize(ctx context.Context) error {
	return nil
}

func (d *defaultSeal) BarrierSealConfigType() SealConfigType {
	return SealConfigTypeShamir
}

func (d *defaultSeal) StoredKeysSupported() seal.StoredKeysSupport {
	switch {
	case d.LegacySeal():
		return seal.StoredKeysNotSupported
	default:
		return seal.StoredKeysSupportedShamirRoot
	}
}

func (d *defaultSeal) RecoveryKeySupported() bool {
	return false
}

func (d *defaultSeal) SetStoredKeys(ctx context.Context, keys [][]byte) error {
	if d.LegacySeal() {
		return fmt.Errorf("stored keys are not supported")
	}
	return writeStoredKeys(ctx, d.core.physical, d.access, keys)
}

func (d *defaultSeal) LegacySeal() bool {
	cfg := d.config.Load().(*SealConfig)
	if cfg == nil {
		return false
	}
	return cfg.StoredShares == 0
}

func (d *defaultSeal) GetStoredKeys(ctx context.Context) ([][]byte, error) {
	if d.LegacySeal() {
		return nil, fmt.Errorf("stored keys are not supported")
	}
	keys, err := readStoredKeys(ctx, d.core.physical, d.access)
	return keys, err
}

func (d *defaultSeal) BarrierConfig(ctx context.Context) (*SealConfig, error) {
	if cfg := d.config.Load().(*SealConfig); cfg != nil {
		return cfg.Clone(), nil
	}

	if err := d.checkCore(); err != nil {
		return nil, err
	}

	// Fetch the core configuration
	conf, err := d.core.PhysicalBarrierSealConfig(ctx)
	if err != nil {
		d.core.logger.Error("failed to read seal configuration", "error", err)
		return nil, fmt.Errorf("failed to check seal configuration: %w", err)
	}

	// If the seal configuration is missing, we are not initialized
	if conf == nil {
		d.core.logger.Info("seal configuration missing, not initialized")
		return nil, nil
	}

	switch conf.Type {
	case d.BarrierSealConfigType().String():
	default:
		d.core.logger.Error("barrier seal type does not match expected type", "barrier_seal_type", conf.Type, "loaded_seal_type", d.BarrierSealConfigType())
		return nil, fmt.Errorf("barrier seal type of %q does not match expected type of %q", conf.Type, d.BarrierSealConfigType())
	}

	d.SetCachedBarrierConfig(conf)
	return conf.Clone(), nil
}

func (d *defaultSeal) ClearBarrierConfig(ctx context.Context) error {
	return d.SetBarrierConfig(ctx, nil)
}

func (d *defaultSeal) SetBarrierConfig(ctx context.Context, config *SealConfig) error {
	if err := d.checkCore(); err != nil {
		return err
	}

	// Provide a way to wipe out the cached value (also prevents actually
	// saving a nil config)
	if config == nil {
		d.config.Store((*SealConfig)(nil))
		return nil
	}

	config.Type = d.BarrierSealConfigType().String()

	// If we are doing a raft unseal we do not want to persist the barrier config
	// because storage isn't setup yet.
	if d.core.isRaftUnseal() {
		d.config.Store(config.Clone())
		return nil
	}

	err := d.core.SetPhysicalBarrierSealConfig(ctx, config)
	if err != nil {
		return err
	}

	d.SetCachedBarrierConfig(config.Clone())

	return nil
}

func (d *defaultSeal) SetCachedBarrierConfig(config *SealConfig) {
	d.config.Store(config)
}

func (d *defaultSeal) RecoverySealConfigType() SealConfigType {
	return SealConfigTypeRecoveryUnsupported
}

func (d *defaultSeal) RecoveryConfig(ctx context.Context) (*SealConfig, error) {
	return nil, fmt.Errorf("recovery not supported")
}

func (d *defaultSeal) RecoveryKey(ctx context.Context) ([]byte, error) {
	return nil, fmt.Errorf("recovery not supported")
}

func (d *defaultSeal) ClearRecoveryConfig(ctx context.Context) error {
	return d.SetRecoveryConfig(ctx, nil)
}

func (d *defaultSeal) SetRecoveryConfig(ctx context.Context, config *SealConfig) error {
	return fmt.Errorf("recovery not supported")
}

func (d *defaultSeal) SetCachedRecoveryConfig(config *SealConfig) {
}

func (d *defaultSeal) VerifyRecoveryKey(ctx context.Context, key []byte) error {
	return fmt.Errorf("recovery not supported")
}

func (d *defaultSeal) SetRecoveryKey(ctx context.Context, key []byte) error {
	return fmt.Errorf("recovery not supported")
}

func (d *defaultSeal) Healthy() bool {
	return true
}

type ErrEncrypt struct {
	Err error
}

var _ error = &ErrEncrypt{}

func (e *ErrEncrypt) Error() string {
	return e.Err.Error()
}

func (e *ErrEncrypt) Is(target error) bool {
	_, ok := target.(*ErrEncrypt)
	return ok || errors.Is(e.Err, target)
}

type ErrDecrypt struct {
	Err error
}

var _ error = &ErrDecrypt{}

func (e *ErrDecrypt) Error() string {
	return e.Err.Error()
}

func (e *ErrDecrypt) Is(target error) bool {
	_, ok := target.(*ErrDecrypt)
	return ok || errors.Is(e.Err, target)
}

func writeStoredKeys(ctx context.Context, storage physical.Backend, encryptor seal.Access, keys [][]byte) error {
	if keys == nil {
		return fmt.Errorf("keys were nil")
	}
	if len(keys) == 0 {
		return fmt.Errorf("no keys provided")
	}

	// Encrypt and marshal the keys
	pe, err := SealWrapStoredBarrierKeys(ctx, encryptor, keys)
	if err != nil {
		return fmt.Errorf("failed to marshal value for storage: %w", err)
	}

	// Store the seal configuration.
	if err := storage.Put(ctx, pe); err != nil {
		return fmt.Errorf("failed to write keys to storage: %w", err)
	}

	return nil
}

func readStoredKeys(ctx context.Context, storage physical.Backend, encryptor seal.Access) ([][]byte, error) {
	pe, err := storage.Get(ctx, StoredBarrierKeysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stored keys: %w", err)
	}

	// This is not strictly an error; we may not have any stored keys, for
	// instance, if we're not initialized
	if pe == nil {
		return nil, nil
	}

	return UnsealWrapStoredBarrierKeys(ctx, encryptor, pe)
}

func (c *Core) SetPhysicalSealGenInfo(ctx context.Context, sealGenInfo *seal.SealGenerationInfo) error {
	if enabled, err := server.IsSealHABetaEnabled(); err != nil {
		return err
	} else if !enabled {
		return nil
	}

	if sealGenInfo == nil {
		return errors.New("invalid seal generation information: generation is unknown")
	}
	// Encode the seal generation info
	buf, err := json.Marshal(sealGenInfo)
	if err != nil {
		return fmt.Errorf("failed to encode seal generation info: %w", err)
	}

	// Store the seal generation info
	pe := &physical.Entry{
		Key:   SealGenInfoPath,
		Value: buf,
	}

	if err := c.physical.Put(ctx, pe); err != nil {
		c.logger.Error("failed to write seal generation info", "error", err)
		return fmt.Errorf("failed to write seal generation info: %w", err)
	}

	return nil
}

func PhysicalSealGenInfo(ctx context.Context, storage physical.Backend) (*seal.SealGenerationInfo, error) {
	if enabled, err := server.IsSealHABetaEnabled(); err != nil {
		return nil, err
	} else if !enabled {
		return nil, nil
	}

	pe, err := storage.Get(ctx, SealGenInfoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch seal generation info: %w", err)
	}
	if pe == nil {
		return nil, nil
	}

	sealGenInfo := new(seal.SealGenerationInfo)

	if err := jsonutil.DecodeJSON(pe.Value, sealGenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode seal generation info: %w", err)
	}

	return sealGenInfo, nil
}
