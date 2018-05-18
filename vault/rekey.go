package vault

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/shamir"
)

const (
	// coreUnsealKeysBackupPath is the path used to backup encrypted unseal
	// keys if specified during a rekey operation. This is outside of the
	// barrier.
	coreBarrierUnsealKeysBackupPath = "core/unseal-keys-backup"

	// coreRecoveryUnsealKeysBackupPath is the path used to backup encrypted
	// recovery keys if specified during a rekey operation. This is outside of
	// the barrier.
	coreRecoveryUnsealKeysBackupPath = "core/recovery-keys-backup"
)

// RekeyResult is used to provide the key parts back after
// they are generated as part of the rekey.
type RekeyResult struct {
	SecretShares         [][]byte
	PGPFingerprints      []string
	Backup               bool
	RecoveryKey          bool
	VerificationRequired bool
	VerificationNonce    string
}

type RekeyVerifyResult struct {
	Nonce string
}

// RekeyBackup stores the backup copy of PGP-encrypted keys
type RekeyBackup struct {
	Nonce string
	Keys  map[string][]string
}

// RekeyThreshold returns the secret threshold for the current seal
// config. This threshold can either be the barrier key threshold or
// the recovery key threshold, depending on whether rekey is being
// performed on the recovery key, or whether the seal supports
// recovery keys.
func (c *Core) RekeyThreshold(ctx context.Context, recovery bool) (int, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return 0, consts.ErrSealed
	}
	if c.standby {
		return 0, consts.ErrStandby
	}

	c.rekeyLock.RLock()
	defer c.rekeyLock.RUnlock()

	var config *SealConfig
	var err error
	// If we are rekeying the recovery key, or if the seal supports
	// recovery keys and we are rekeying the barrier key, we use the
	// recovery config as the threshold instead.
	if recovery || c.seal.RecoveryKeySupported() {
		config, err = c.seal.RecoveryConfig(ctx)
	} else {
		config, err = c.seal.BarrierConfig(ctx)
	}
	if err != nil {
		return 0, err
	}

	return config.SecretThreshold, nil
}

// RekeyProgress is used to return the rekey progress (num shares).
func (c *Core) RekeyProgress(recovery bool) (int, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return 0, consts.ErrSealed
	}
	if c.standby {
		return 0, consts.ErrStandby
	}

	c.rekeyLock.RLock()
	defer c.rekeyLock.RUnlock()

	if recovery {
		return len(c.recoveryRekeyProgress), nil
	}
	return len(c.barrierRekeyProgress), nil
}

// RekeyVerifyProgress is used to return the rekey progress (num shares) during
// verification.
func (c *Core) RekeyVerifyProgress(recovery bool) (int, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return 0, consts.ErrSealed
	}
	if c.standby {
		return 0, consts.ErrStandby
	}

	c.rekeyLock.RLock()
	defer c.rekeyLock.RUnlock()

	if recovery {
		return len(c.recoveryRekeyVerifyProgress), nil
	}
	return len(c.barrierRekeyVerifyProgress), nil
}

// RekeyConfig is used to read the rekey configuration
func (c *Core) RekeyConfig(recovery bool) (*SealConfig, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, consts.ErrSealed
	}
	if c.standby {
		return nil, consts.ErrStandby
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Copy the seal config if any
	var conf *SealConfig
	if recovery {
		if c.recoveryRekeyConfig != nil {
			conf = c.recoveryRekeyConfig.Clone()
		}
	} else {
		if c.barrierRekeyConfig != nil {
			conf = c.barrierRekeyConfig.Clone()
		}
	}

	return conf, nil
}

// RekeyInit will either initialize the rekey of barrier or recovery key.
// recovery determines whether this is a rekey on the barrier or recovery key.
func (c *Core) RekeyInit(config *SealConfig, recovery bool) error {
	if recovery {
		return c.RecoveryRekeyInit(config)
	}
	return c.BarrierRekeyInit(config)
}

// BarrierRekeyInit is used to initialize the rekey settings for the barrier key
func (c *Core) BarrierRekeyInit(config *SealConfig) error {
	if config.StoredShares > 0 {
		if !c.seal.StoredKeysSupported() {
			return fmt.Errorf("storing keys not supported by barrier seal")
		}
		if len(config.PGPKeys) > 0 {
			return fmt.Errorf("PGP key encryption not supported when using stored keys")
		}
		if config.Backup {
			return fmt.Errorf("key backup not supported when using stored keys")
		}
	}

	if c.seal.RecoveryKeySupported() && c.seal.RecoveryType() == config.Type {
		c.logger.Debug("using recovery seal configuration to rekey barrier key")
		if config.VerificationRequired {
			return fmt.Errorf("requiring verification not supported when rekeying the barrier key with recovery keys")
		}
	}

	// Check if the seal configuration is valid
	if err := config.Validate(); err != nil {
		c.logger.Error("invalid rekey seal configuration", "error", err)
		return errwrap.Wrapf("invalid rekey seal configuration: {{err}}", err)
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return consts.ErrSealed
	}
	if c.standby {
		return consts.ErrStandby
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Prevent multiple concurrent re-keys
	if c.barrierRekeyConfig != nil {
		return fmt.Errorf("rekey already in progress")
	}

	// Copy the configuration
	c.barrierRekeyConfig = config.Clone()

	// Initialize the nonce
	nonce, err := uuid.GenerateUUID()
	if err != nil {
		c.barrierRekeyConfig = nil
		return err
	}
	c.barrierRekeyConfig.Nonce = nonce

	if c.logger.IsInfo() {
		c.logger.Info("rekey initialized", "nonce", c.barrierRekeyConfig.Nonce, "shares", c.barrierRekeyConfig.SecretShares, "threshold", c.barrierRekeyConfig.SecretThreshold, "validation_required", c.barrierRekeyConfig.VerificationRequired)
	}
	return nil
}

// RecoveryRekeyInit is used to initialize the rekey settings for the recovery key
func (c *Core) RecoveryRekeyInit(config *SealConfig) error {
	if config.StoredShares > 0 {
		return fmt.Errorf("stored shares not supported by recovery key")
	}

	// Check if the seal configuration is valid
	if err := config.Validate(); err != nil {
		c.logger.Error("invalid recovery configuration", "error", err)
		return errwrap.Wrapf("invalid recovery configuration: {{err}}", err)
	}

	if !c.seal.RecoveryKeySupported() {
		return fmt.Errorf("recovery keys not supported")
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return consts.ErrSealed
	}
	if c.standby {
		return consts.ErrStandby
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Prevent multiple concurrent re-keys
	if c.recoveryRekeyConfig != nil {
		return fmt.Errorf("rekey already in progress")
	}

	// Copy the configuration
	c.recoveryRekeyConfig = config.Clone()

	// Initialize the nonce
	nonce, err := uuid.GenerateUUID()
	if err != nil {
		c.recoveryRekeyConfig = nil
		return err
	}
	c.recoveryRekeyConfig.Nonce = nonce

	if c.logger.IsInfo() {
		c.logger.Info("rekey initialized", "nonce", c.recoveryRekeyConfig.Nonce, "shares", c.recoveryRekeyConfig.SecretShares, "threshold", c.recoveryRekeyConfig.SecretThreshold, "validation_required", c.recoveryRekeyConfig.VerificationRequired)
	}
	return nil
}

// RekeyUpdate is used to provide a new key part for the barrier or recovery key.
func (c *Core) RekeyUpdate(ctx context.Context, key []byte, nonce string, recovery bool) (*RekeyResult, error) {
	if recovery {
		return c.RecoveryRekeyUpdate(ctx, key, nonce)
	}
	return c.BarrierRekeyUpdate(ctx, key, nonce)
}

// BarrierRekeyUpdate is used to provide a new key part. Barrier rekey can be done
// with unseal keys, or recovery keys if that's supported and we are storing the barrier
// key.
//
// N.B.: If recovery keys are used to rekey, the new barrier key shares are not returned.
func (c *Core) BarrierRekeyUpdate(ctx context.Context, key []byte, nonce string) (*RekeyResult, error) {
	// Ensure we are already unsealed
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, consts.ErrSealed
	}
	if c.standby {
		return nil, consts.ErrStandby
	}

	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		return nil, &ErrInvalidKey{fmt.Sprintf("key is shorter than minimum %d bytes", min)}
	}
	if len(key) > max {
		return nil, &ErrInvalidKey{fmt.Sprintf("key is longer than maximum %d bytes", max)}
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Get the seal configuration
	var existingConfig *SealConfig
	var err error
	var useRecovery bool // Determines whether recovery key is being used to rekey the master key
	if c.seal.StoredKeysSupported() && c.seal.RecoveryKeySupported() {
		existingConfig, err = c.seal.RecoveryConfig(ctx)
		useRecovery = true
	} else {
		existingConfig, err = c.seal.BarrierConfig(ctx)
	}
	if err != nil {
		return nil, err
	}

	// Ensure the barrier is initialized
	if existingConfig == nil {
		return nil, ErrNotInit
	}

	// Ensure a rekey is in progress
	if c.barrierRekeyConfig == nil {
		return nil, fmt.Errorf("no rekey in progress")
	}

	if nonce != c.barrierRekeyConfig.Nonce {
		return nil, fmt.Errorf("incorrect nonce supplied; nonce for this rekey operation is %q", c.barrierRekeyConfig.Nonce)
	}

	// Check if we already have this piece
	for _, existing := range c.barrierRekeyProgress {
		if bytes.Equal(existing, key) {
			return nil, fmt.Errorf("given key has already been provided during this generation operation")
		}
	}

	// Store this key
	c.barrierRekeyProgress = append(c.barrierRekeyProgress, key)

	// Check if we don't have enough keys to unlock
	if len(c.barrierRekeyProgress) < existingConfig.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("cannot rekey yet, not enough keys", "keys", len(c.barrierRekeyProgress), "threshold", existingConfig.SecretThreshold)
		}
		return nil, nil
	}

	// Schedule the rekey progress for forgetting
	defer func() {
		c.barrierRekeyProgress = nil
	}()

	// Recover the master key or recovery key
	var recoveredKey []byte
	if existingConfig.SecretThreshold == 1 {
		recoveredKey = c.barrierRekeyProgress[0]
	} else {
		recoveredKey, err = shamir.Combine(c.barrierRekeyProgress)
		if err != nil {
			return nil, errwrap.Wrapf("failed to compute master key: {{err}}", err)
		}
	}

	if useRecovery {
		if err := c.seal.VerifyRecoveryKey(ctx, recoveredKey); err != nil {
			c.logger.Error("rekey aborted, recovery key verification failed", "error", err)
			return nil, err
		}
	} else {
		if err := c.barrier.VerifyMaster(recoveredKey); err != nil {
			c.logger.Error("rekey aborted, master key verification failed", "error", err)
			return nil, err
		}
	}

	// Generate a new master key
	newMasterKey, err := c.barrier.GenerateKey()
	if err != nil {
		c.logger.Error("failed to generate master key", "error", err)
		return nil, errwrap.Wrapf("master key generation failed: {{err}}", err)
	}

	results := &RekeyResult{
		Backup: c.barrierRekeyConfig.Backup,
	}
	// Set result.SecretShares to the master key if only a single key
	// part is used -- no Shamir split required.
	if c.barrierRekeyConfig.SecretShares == 1 {
		results.SecretShares = append(results.SecretShares, newMasterKey)
	} else {
		// Split the master key using the Shamir algorithm
		shares, err := shamir.Split(newMasterKey, c.barrierRekeyConfig.SecretShares, c.barrierRekeyConfig.SecretThreshold)
		if err != nil {
			c.logger.Error("failed to generate shares", "error", err)
			return nil, errwrap.Wrapf("failed to generate shares: {{err}}", err)
		}
		results.SecretShares = shares
	}

	// If we are storing any shares, add them to the shares to store and remove
	// from the returned keys
	var keysToStore [][]byte
	if c.seal.StoredKeysSupported() && c.barrierRekeyConfig.StoredShares > 0 {
		for i := 0; i < c.barrierRekeyConfig.StoredShares; i++ {
			keysToStore = append(keysToStore, results.SecretShares[0])
			results.SecretShares = results.SecretShares[1:]
		}
	}

	// If PGP keys are passed in, encrypt shares with corresponding PGP keys.
	if len(c.barrierRekeyConfig.PGPKeys) > 0 {
		hexEncodedShares := make([][]byte, len(results.SecretShares))
		for i, _ := range results.SecretShares {
			hexEncodedShares[i] = []byte(hex.EncodeToString(results.SecretShares[i]))
		}
		results.PGPFingerprints, results.SecretShares, err = pgpkeys.EncryptShares(hexEncodedShares, c.barrierRekeyConfig.PGPKeys)
		if err != nil {
			return nil, err
		}

		// If backup is enabled, store backup info in vault.coreBarrierUnsealKeysBackupPath
		if c.barrierRekeyConfig.Backup {
			backupInfo := map[string][]string{}
			for i := 0; i < len(results.PGPFingerprints); i++ {
				encShare := bytes.NewBuffer(results.SecretShares[i])
				if backupInfo[results.PGPFingerprints[i]] == nil {
					backupInfo[results.PGPFingerprints[i]] = []string{hex.EncodeToString(encShare.Bytes())}
				} else {
					backupInfo[results.PGPFingerprints[i]] = append(backupInfo[results.PGPFingerprints[i]], hex.EncodeToString(encShare.Bytes()))
				}
			}

			backupVals := &RekeyBackup{
				Nonce: c.barrierRekeyConfig.Nonce,
				Keys:  backupInfo,
			}
			buf, err := json.Marshal(backupVals)
			if err != nil {
				c.logger.Error("failed to marshal unseal key backup", "error", err)
				return nil, errwrap.Wrapf("failed to marshal unseal key backup: {{err}}", err)
			}
			pe := &physical.Entry{
				Key:   coreBarrierUnsealKeysBackupPath,
				Value: buf,
			}
			if err = c.physical.Put(ctx, pe); err != nil {
				c.logger.Error("failed to save unseal key backup", "error", err)
				return nil, errwrap.Wrapf("failed to save unseal key backup: {{err}}", err)
			}
		}
	}

	if keysToStore != nil {
		if err := c.seal.SetStoredKeys(ctx, keysToStore); err != nil {
			c.logger.Error("failed to store keys", "error", err)
			return nil, errwrap.Wrapf("failed to store keys: {{err}}", err)
		}
	}

	// If we are requiring validation, return now; otherwise rekey the barrier
	if c.barrierRekeyConfig.VerificationRequired {
		nonce, err := uuid.GenerateUUID()
		if err != nil {
			c.barrierRekeyConfig = nil
			return nil, err
		}
		c.barrierRekeyConfig.VerificationNonce = nonce
		c.barrierRekeyConfig.VerificationKey = newMasterKey

		results.VerificationRequired = true
		results.VerificationNonce = nonce
		return results, nil
	}

	if err := c.performBarrierRekey(ctx, newMasterKey); err != nil {
		return nil, err
	}

	c.barrierRekeyConfig = nil
	return results, nil
}

func (c *Core) performBarrierRekey(ctx context.Context, newMasterKey []byte) error {
	// Rekey the barrier
	if err := c.barrier.Rekey(ctx, newMasterKey); err != nil {
		c.logger.Error("failed to rekey barrier", "error", err)
		return errwrap.Wrapf("failed to rekey barrier: {{err}}", err)
	}
	if c.logger.IsInfo() {
		c.logger.Info("security barrier rekeyed", "shares", c.barrierRekeyConfig.SecretShares, "threshold", c.barrierRekeyConfig.SecretThreshold)
	}

	c.barrierRekeyConfig.VerificationKey = nil

	if err := c.seal.SetBarrierConfig(ctx, c.barrierRekeyConfig); err != nil {
		c.logger.Error("error saving rekey seal configuration", "error", err)
		return errwrap.Wrapf("failed to save rekey seal configuration: {{err}}", err)
	}

	// Write to the canary path, which will force a synchronous truing during
	// replication
	if err := c.barrier.Put(ctx, &Entry{
		Key:   coreKeyringCanaryPath,
		Value: []byte(c.barrierRekeyConfig.Nonce),
	}); err != nil {
		c.logger.Error("error saving keyring canary", "error", err)
		return errwrap.Wrapf("failed to save keyring canary: {{err}}", err)
	}

	return nil
}

// RecoveryRekeyUpdate is used to provide a new key part
func (c *Core) RecoveryRekeyUpdate(ctx context.Context, key []byte, nonce string) (*RekeyResult, error) {
	// Ensure we are already unsealed
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, consts.ErrSealed
	}
	if c.standby {
		return nil, consts.ErrStandby
	}

	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		return nil, &ErrInvalidKey{fmt.Sprintf("key is shorter than minimum %d bytes", min)}
	}
	if len(key) > max {
		return nil, &ErrInvalidKey{fmt.Sprintf("key is longer than maximum %d bytes", max)}
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Get the seal configuration
	existingConfig, err := c.seal.RecoveryConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Ensure the seal is initialized
	if existingConfig == nil {
		return nil, ErrNotInit
	}

	// Ensure a rekey is in progress
	if c.recoveryRekeyConfig == nil {
		return nil, fmt.Errorf("no rekey in progress")
	}

	if nonce != c.recoveryRekeyConfig.Nonce {
		return nil, fmt.Errorf("incorrect nonce supplied; nonce for this rekey operation is %q", c.recoveryRekeyConfig.Nonce)
	}

	// Check if we already have this piece
	for _, existing := range c.recoveryRekeyProgress {
		if bytes.Equal(existing, key) {
			return nil, nil
		}
	}

	// Store this key
	c.recoveryRekeyProgress = append(c.recoveryRekeyProgress, key)

	// Check if we don't have enough keys to unlock
	if len(c.recoveryRekeyProgress) < existingConfig.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("cannot rekey yet, not enough keys", "keys", len(c.recoveryRekeyProgress), "threshold", existingConfig.SecretThreshold)
		}
		return nil, nil
	}

	// Schedule the rekey progress for forgetting
	defer func() {
		c.recoveryRekeyProgress = nil
	}()

	// Recover the master key
	var recoveryKey []byte
	if existingConfig.SecretThreshold == 1 {
		recoveryKey = c.recoveryRekeyProgress[0]
	} else {
		recoveryKey, err = shamir.Combine(c.recoveryRekeyProgress)
		if err != nil {
			return nil, errwrap.Wrapf("failed to compute recovery key: {{err}}", err)
		}
	}

	// Verify the recovery key
	if err := c.seal.VerifyRecoveryKey(ctx, recoveryKey); err != nil {
		c.logger.Error("rekey aborted, recovery key verification failed", "error", err)
		return nil, err
	}

	// Generate a new master key
	newMasterKey, err := c.barrier.GenerateKey()
	if err != nil {
		c.logger.Error("failed to generate recovery key", "error", err)
		return nil, errwrap.Wrapf("recovery key generation failed: {{err}}", err)
	}

	// Return the master key if only a single key part is used
	results := &RekeyResult{
		Backup: c.recoveryRekeyConfig.Backup,
	}

	if c.recoveryRekeyConfig.SecretShares == 1 {
		results.SecretShares = append(results.SecretShares, newMasterKey)
	} else {
		// Split the master key using the Shamir algorithm
		shares, err := shamir.Split(newMasterKey, c.recoveryRekeyConfig.SecretShares, c.recoveryRekeyConfig.SecretThreshold)
		if err != nil {
			c.logger.Error("failed to generate shares", "error", err)
			return nil, errwrap.Wrapf("failed to generate shares: {{err}}", err)
		}
		results.SecretShares = shares
	}

	if len(c.recoveryRekeyConfig.PGPKeys) > 0 {
		hexEncodedShares := make([][]byte, len(results.SecretShares))
		for i, _ := range results.SecretShares {
			hexEncodedShares[i] = []byte(hex.EncodeToString(results.SecretShares[i]))
		}
		results.PGPFingerprints, results.SecretShares, err = pgpkeys.EncryptShares(hexEncodedShares, c.recoveryRekeyConfig.PGPKeys)
		if err != nil {
			return nil, err
		}

		if c.recoveryRekeyConfig.Backup {
			backupInfo := map[string][]string{}
			for i := 0; i < len(results.PGPFingerprints); i++ {
				encShare := bytes.NewBuffer(results.SecretShares[i])
				if backupInfo[results.PGPFingerprints[i]] == nil {
					backupInfo[results.PGPFingerprints[i]] = []string{hex.EncodeToString(encShare.Bytes())}
				} else {
					backupInfo[results.PGPFingerprints[i]] = append(backupInfo[results.PGPFingerprints[i]], hex.EncodeToString(encShare.Bytes()))
				}
			}

			backupVals := &RekeyBackup{
				Nonce: c.recoveryRekeyConfig.Nonce,
				Keys:  backupInfo,
			}
			buf, err := json.Marshal(backupVals)
			if err != nil {
				c.logger.Error("failed to marshal recovery key backup", "error", err)
				return nil, errwrap.Wrapf("failed to marshal recovery key backup: {{err}}", err)
			}
			pe := &physical.Entry{
				Key:   coreRecoveryUnsealKeysBackupPath,
				Value: buf,
			}
			if err = c.physical.Put(ctx, pe); err != nil {
				c.logger.Error("failed to save unseal key backup", "error", err)
				return nil, errwrap.Wrapf("failed to save unseal key backup: {{err}}", err)
			}
		}
	}

	// If we are requiring validation, return now; otherwise save the recovery
	// key
	if c.recoveryRekeyConfig.VerificationRequired {
		nonce, err := uuid.GenerateUUID()
		if err != nil {
			c.recoveryRekeyConfig = nil
			return nil, err
		}
		c.recoveryRekeyConfig.VerificationNonce = nonce
		c.recoveryRekeyConfig.VerificationKey = newMasterKey

		results.VerificationRequired = true
		results.VerificationNonce = nonce
		return results, nil
	}

	if err := c.performRecoveryRekey(ctx, newMasterKey); err != nil {
		return nil, err
	}

	c.recoveryRekeyConfig = nil
	return results, nil
}

func (c *Core) performRecoveryRekey(ctx context.Context, newMasterKey []byte) error {
	if err := c.seal.SetRecoveryKey(ctx, newMasterKey); err != nil {
		c.logger.Error("failed to set recovery key", "error", err)
		return errwrap.Wrapf("failed to set recovery key: {{err}}", err)
	}

	c.recoveryRekeyConfig.VerificationKey = nil

	if err := c.seal.SetRecoveryConfig(ctx, c.recoveryRekeyConfig); err != nil {
		c.logger.Error("error saving rekey seal configuration", "error", err)
		return errwrap.Wrapf("failed to save rekey seal configuration: {{err}}", err)
	}

	// Write to the canary path, which will force a synchronous truing during
	// replication
	if err := c.barrier.Put(ctx, &Entry{
		Key:   coreKeyringCanaryPath,
		Value: []byte(c.recoveryRekeyConfig.Nonce),
	}); err != nil {
		c.logger.Error("error saving keyring canary", "error", err)
		return errwrap.Wrapf("failed to save keyring canary: {{err}}", err)
	}

	return nil
}

func (c *Core) RekeyVerify(ctx context.Context, key []byte, recovery bool) (*RekeyVerifyResult, error) {
	if recovery {
		//return c.RecoveryRekeyVerify(ctx, key)
	}
	return c.BarrierRekeyVerify(ctx, key)
}

func (c *Core) BarrierRekeyVerify(ctx context.Context, key []byte) (*RekeyVerifyResult, error) {
	// Ensure we are already unsealed
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, consts.ErrSealed
	}
	if c.standby {
		return nil, consts.ErrStandby
	}

	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		return nil, &ErrInvalidKey{fmt.Sprintf("key is shorter than minimum %d bytes", min)}
	}
	if len(key) > max {
		return nil, &ErrInvalidKey{fmt.Sprintf("key is longer than maximum %d bytes", max)}
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Ensure a rekey is in progress
	if c.barrierRekeyConfig == nil {
		return nil, fmt.Errorf("no rekey in progress")
	}

	// Check if we already have this piece
	for _, existing := range c.barrierRekeyVerifyProgress {
		if bytes.Equal(existing, key) {
			return nil, fmt.Errorf("given key has already been provided during this verify operation")
		}
	}

	// Store this key
	c.barrierRekeyVerifyProgress = append(c.barrierRekeyVerifyProgress, key)

	// Check if we don't have enough keys to unlock
	if len(c.barrierRekeyVerifyProgress) < c.barrierRekeyConfig.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("cannot verify yet, not enough keys", "keys", len(c.barrierRekeyVerifyProgress), "threshold", c.barrierRekeyConfig.SecretThreshold)
		}
		return nil, nil
	}

	// Schedule the progress for forgetting and rotate the nonce if possible
	defer func() {
		c.barrierRekeyVerifyProgress = nil
		if c.barrierRekeyConfig != nil {
			nonce, err := uuid.GenerateUUID()
			if err == nil {
				c.barrierRekeyConfig.VerificationNonce = nonce
			}
		}
	}()

	// Recover the master key or recovery key
	var recoveredKey []byte
	var err error
	if c.barrierRekeyConfig.SecretThreshold == 1 {
		recoveredKey = c.barrierRekeyVerifyProgress[0]
	} else {
		recoveredKey, err = shamir.Combine(c.barrierRekeyVerifyProgress)
		if err != nil {
			return nil, errwrap.Wrapf("failed to compute master key for verification: {{err}}", err)
		}
	}

	if subtle.ConstantTimeCompare(recoveredKey, c.barrierRekeyConfig.VerificationKey) != 1 {
		c.logger.Error("rekey verification failed")
		return nil, errors.New("rekey verification failed")
	}

	if err := c.performBarrierRekey(ctx, recoveredKey); err != nil {
		return nil, err
	}

	res := &RekeyVerifyResult{
		Nonce: c.barrierRekeyConfig.VerificationNonce,
	}

	c.barrierRekeyConfig = nil
	return res, nil
}

// RekeyCancel is used to cancel an inprogress rekey
func (c *Core) RekeyCancel(recovery bool) error {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return consts.ErrSealed
	}
	if c.standby {
		return consts.ErrStandby
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Clear any progress or config
	if recovery {
		c.recoveryRekeyConfig = nil
		c.recoveryRekeyProgress = nil
	} else {
		c.barrierRekeyConfig = nil
		c.barrierRekeyProgress = nil
	}
	return nil
}

// RekeyVerifyCancel is used to start the verification process over
func (c *Core) RekeyVerifyRestart(recovery bool) error {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return consts.ErrSealed
	}
	if c.standby {
		return consts.ErrStandby
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Attempt to generate a new nonce, but don't bail if it doesn't succeed
	// (which is extraordinarily unlikely)
	nonce, nonceErr := uuid.GenerateUUID()

	// Clear any progress or config
	if recovery {
		c.recoveryRekeyVerifyProgress = nil
		if nonceErr == nil {
			c.recoveryRekeyConfig.VerificationNonce = nonce
		}
	} else {
		c.barrierRekeyVerifyProgress = nil
		if nonceErr == nil {
			c.barrierRekeyConfig.VerificationNonce = nonce
		}
	}

	return nil
}

// RekeyRetrieveBackup is used to retrieve any backed-up PGP-encrypted unseal
// keys
func (c *Core) RekeyRetrieveBackup(ctx context.Context, recovery bool) (*RekeyBackup, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, consts.ErrSealed
	}
	if c.standby {
		return nil, consts.ErrStandby
	}

	c.rekeyLock.RLock()
	defer c.rekeyLock.RUnlock()

	var entry *physical.Entry
	var err error
	if recovery {
		entry, err = c.physical.Get(ctx, coreRecoveryUnsealKeysBackupPath)
	} else {
		entry, err = c.physical.Get(ctx, coreBarrierUnsealKeysBackupPath)
	}
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	ret := &RekeyBackup{}
	err = jsonutil.DecodeJSON(entry.Value, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// RekeyDeleteBackup is used to delete any backed-up PGP-encrypted unseal keys
func (c *Core) RekeyDeleteBackup(ctx context.Context, recovery bool) error {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return consts.ErrSealed
	}
	if c.standby {
		return consts.ErrStandby
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	if recovery {
		return c.physical.Delete(ctx, coreRecoveryUnsealKeysBackupPath)
	}
	return c.physical.Delete(ctx, coreBarrierUnsealKeysBackupPath)
}
