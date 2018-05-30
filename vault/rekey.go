package vault

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/logical"
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
	Complete bool
	Nonce    string
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
func (c *Core) RekeyThreshold(ctx context.Context, recovery bool) (int, logical.HTTPCodedError) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return 0, logical.CodedError(http.StatusServiceUnavailable, consts.ErrSealed.Error())
	}
	if c.standby {
		return 0, logical.CodedError(http.StatusBadRequest, consts.ErrStandby.Error())
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
		return 0, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("unable to look up config: {{err}}", err).Error())
	}
	if config == nil {
		return 0, logical.CodedError(http.StatusBadRequest, ErrNotInit.Error())
	}

	return config.SecretThreshold, nil
}

// RekeyProgress is used to return the rekey progress (num shares).
func (c *Core) RekeyProgress(recovery, verification bool) (bool, int, logical.HTTPCodedError) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return false, 0, logical.CodedError(http.StatusServiceUnavailable, consts.ErrSealed.Error())
	}
	if c.standby {
		return false, 0, logical.CodedError(http.StatusBadRequest, consts.ErrStandby.Error())
	}

	c.rekeyLock.RLock()
	defer c.rekeyLock.RUnlock()

	var conf *SealConfig
	if recovery {
		conf = c.recoveryRekeyConfig
	} else {
		conf = c.barrierRekeyConfig
	}

	if conf == nil {
		return false, 0, logical.CodedError(http.StatusBadRequest, "rekey operation not in progress")
	}

	if verification {
		return len(conf.VerificationKey) > 0, len(conf.VerificationProgress), nil
	}
	return true, len(conf.RekeyProgress), nil
}

// RekeyConfig is used to read the rekey configuration
func (c *Core) RekeyConfig(recovery bool) (*SealConfig, logical.HTTPCodedError) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, logical.CodedError(http.StatusServiceUnavailable, consts.ErrSealed.Error())
	}
	if c.standby {
		return nil, logical.CodedError(http.StatusBadRequest, consts.ErrStandby.Error())
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
func (c *Core) RekeyInit(config *SealConfig, recovery bool) logical.HTTPCodedError {
	if config.SecretThreshold > config.SecretShares {
		return logical.CodedError(http.StatusBadRequest, "provided threshold greater than the total shares")
	}

	if recovery {
		return c.RecoveryRekeyInit(config)
	}
	return c.BarrierRekeyInit(config)
}

// BarrierRekeyInit is used to initialize the rekey settings for the barrier key
func (c *Core) BarrierRekeyInit(config *SealConfig) logical.HTTPCodedError {
	if c.seal.StoredKeysSupported() {
		c.logger.Warn("stored keys supported, forcing rekey shares/threshold to 1")
		config.SecretShares = 1
		config.SecretThreshold = 1
		config.StoredShares = 1
	}

	if config.StoredShares > 0 {
		if !c.seal.StoredKeysSupported() {
			return logical.CodedError(http.StatusBadRequest, "storing keys not supported by barrier seal")
		}
		if len(config.PGPKeys) > 0 {
			return logical.CodedError(http.StatusBadRequest, "PGP key encryption not supported when using stored keys")
		}
		if config.Backup {
			return logical.CodedError(http.StatusBadRequest, "key backup not supported when using stored keys")
		}

		if c.seal.RecoveryKeySupported() {
			if config.VerificationRequired {
				return logical.CodedError(http.StatusBadRequest, "requiring verification not supported when rekeying the barrier key with recovery keys")
			}
			c.logger.Debug("using recovery seal configuration to rekey barrier key")
		}
	}

	// Check if the seal configuration is valid
	if err := config.Validate(); err != nil {
		c.logger.Error("invalid rekey seal configuration", "error", err)
		return logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("invalid rekey seal configuration: {{err}}", err).Error())
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return logical.CodedError(http.StatusServiceUnavailable, consts.ErrSealed.Error())
	}
	if c.standby {
		return logical.CodedError(http.StatusBadRequest, consts.ErrStandby.Error())
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Prevent multiple concurrent re-keys
	if c.barrierRekeyConfig != nil {
		return logical.CodedError(http.StatusBadRequest, "rekey already in progress")
	}

	// Copy the configuration
	c.barrierRekeyConfig = config.Clone()

	// Initialize the nonce
	nonce, err := uuid.GenerateUUID()
	if err != nil {
		c.barrierRekeyConfig = nil
		return logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("error generating nonce for procedure: {{err}}", err).Error())
	}
	c.barrierRekeyConfig.Nonce = nonce

	if c.logger.IsInfo() {
		c.logger.Info("rekey initialized", "nonce", c.barrierRekeyConfig.Nonce, "shares", c.barrierRekeyConfig.SecretShares, "threshold", c.barrierRekeyConfig.SecretThreshold, "validation_required", c.barrierRekeyConfig.VerificationRequired)
	}
	return nil
}

// RecoveryRekeyInit is used to initialize the rekey settings for the recovery key
func (c *Core) RecoveryRekeyInit(config *SealConfig) logical.HTTPCodedError {
	if config.StoredShares > 0 {
		return logical.CodedError(http.StatusBadRequest, "stored shares not supported by recovery key")
	}

	// Check if the seal configuration is valid
	if err := config.Validate(); err != nil {
		c.logger.Error("invalid recovery configuration", "error", err)
		return logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("invalid recovery configuration: {{err}}", err).Error())
	}

	if !c.seal.RecoveryKeySupported() {
		return logical.CodedError(http.StatusBadRequest, "recovery keys not supported")
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return logical.CodedError(http.StatusServiceUnavailable, consts.ErrSealed.Error())
	}
	if c.standby {
		return logical.CodedError(http.StatusBadRequest, consts.ErrStandby.Error())
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Prevent multiple concurrent re-keys
	if c.recoveryRekeyConfig != nil {
		return logical.CodedError(http.StatusBadRequest, "rekey already in progress")
	}

	// Copy the configuration
	c.recoveryRekeyConfig = config.Clone()

	// Initialize the nonce
	nonce, err := uuid.GenerateUUID()
	if err != nil {
		c.recoveryRekeyConfig = nil
		return logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("error generating nonce for procedure: {{err}}", err).Error())
	}
	c.recoveryRekeyConfig.Nonce = nonce

	if c.logger.IsInfo() {
		c.logger.Info("rekey initialized", "nonce", c.recoveryRekeyConfig.Nonce, "shares", c.recoveryRekeyConfig.SecretShares, "threshold", c.recoveryRekeyConfig.SecretThreshold, "validation_required", c.recoveryRekeyConfig.VerificationRequired)
	}
	return nil
}

// RekeyUpdate is used to provide a new key part for the barrier or recovery key.
func (c *Core) RekeyUpdate(ctx context.Context, key []byte, nonce string, recovery bool) (*RekeyResult, logical.HTTPCodedError) {
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
func (c *Core) BarrierRekeyUpdate(ctx context.Context, key []byte, nonce string) (*RekeyResult, logical.HTTPCodedError) {
	// Ensure we are already unsealed
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, logical.CodedError(http.StatusServiceUnavailable, consts.ErrSealed.Error())
	}
	if c.standby {
		return nil, logical.CodedError(http.StatusBadRequest, consts.ErrStandby.Error())
	}

	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		return nil, logical.CodedError(http.StatusBadRequest, fmt.Sprintf("key is shorter than minimum %d bytes", min))
	}
	if len(key) > max {
		return nil, logical.CodedError(http.StatusBadRequest, fmt.Sprintf("key is longer than maximum %d bytes", max))
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
		return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to fetch existing config: {{err}}", err).Error())
	}
	// Ensure the barrier is initialized
	if existingConfig == nil {
		return nil, logical.CodedError(http.StatusBadRequest, ErrNotInit.Error())
	}

	// Ensure a rekey is in progress
	if c.barrierRekeyConfig == nil {
		return nil, logical.CodedError(http.StatusBadRequest, "no barrier rekey in progress")
	}

	if len(c.barrierRekeyConfig.VerificationKey) > 0 {
		return nil, logical.CodedError(http.StatusBadRequest, fmt.Sprintf("rekey operation already finished; verification must be performed; nonce for the verification operation is %q", c.barrierRekeyConfig.VerificationNonce))
	}

	if nonce != c.barrierRekeyConfig.Nonce {
		return nil, logical.CodedError(http.StatusBadRequest, fmt.Sprintf("incorrect nonce supplied; nonce for this rekey operation is %q", c.barrierRekeyConfig.Nonce))
	}

	// Check if we already have this piece
	for _, existing := range c.barrierRekeyConfig.RekeyProgress {
		if subtle.ConstantTimeCompare(existing, key) == 1 {
			return nil, logical.CodedError(http.StatusBadRequest, "given key has already been provided during this generation operation")
		}
	}

	// Store this key
	c.barrierRekeyConfig.RekeyProgress = append(c.barrierRekeyConfig.RekeyProgress, key)

	// Check if we don't have enough keys to unlock
	if len(c.barrierRekeyConfig.RekeyProgress) < existingConfig.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("cannot rekey yet, not enough keys", "keys", len(c.barrierRekeyConfig.RekeyProgress), "threshold", existingConfig.SecretThreshold)
		}
		return nil, nil
	}

	// Recover the master key or recovery key
	var recoveredKey []byte
	if existingConfig.SecretThreshold == 1 {
		recoveredKey = c.barrierRekeyConfig.RekeyProgress[0]
	} else {
		recoveredKey, err = shamir.Combine(c.barrierRekeyConfig.RekeyProgress)
		if err != nil {
			return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to compute master key: {{err}}", err).Error())
		}
	}

	if useRecovery {
		if err := c.seal.VerifyRecoveryKey(ctx, recoveredKey); err != nil {
			c.logger.Error("rekey recovery key verification failed", "error", err)
			return nil, logical.CodedError(http.StatusBadRequest, errwrap.Wrapf("recovery key verification failed: {{err}}", err).Error())
		}
	} else {
		if err := c.barrier.VerifyMaster(recoveredKey); err != nil {
			c.logger.Error("master key verification failed", "error", err)
			return nil, logical.CodedError(http.StatusBadRequest, errwrap.Wrapf("master key verification failed: {{err}}", err).Error())
		}
	}

	// Generate a new master key
	newMasterKey, err := c.barrier.GenerateKey()
	if err != nil {
		c.logger.Error("failed to generate master key", "error", err)
		return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("master key generation failed: {{err}}", err).Error())
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
			return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to generate shares: {{err}}", err).Error())
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
			return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to encrypt shares: {{err}}", err).Error())
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
				return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to marshal unseal key backup: {{err}}", err).Error())
			}
			pe := &physical.Entry{
				Key:   coreBarrierUnsealKeysBackupPath,
				Value: buf,
			}
			if err = c.physical.Put(ctx, pe); err != nil {
				c.logger.Error("failed to save unseal key backup", "error", err)
				return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to save unseal key backup: {{err}}", err).Error())
			}
		}
	}

	if keysToStore != nil {
		if err := c.seal.SetStoredKeys(ctx, keysToStore); err != nil {
			c.logger.Error("failed to store keys", "error", err)
			return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to store keys: {{err}}", err).Error())
		}
	}

	// If we are requiring validation, return now; otherwise rekey the barrier
	if c.barrierRekeyConfig.VerificationRequired {
		nonce, err := uuid.GenerateUUID()
		if err != nil {
			c.barrierRekeyConfig = nil
			return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to generate verification nonce: {{err}}", err).Error())
		}
		c.barrierRekeyConfig.VerificationNonce = nonce
		c.barrierRekeyConfig.VerificationKey = newMasterKey

		results.VerificationRequired = true
		results.VerificationNonce = nonce
		return results, nil
	}

	if err := c.performBarrierRekey(ctx, newMasterKey); err != nil {
		return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to perform barrier rekey: {{err}}", err).Error())
	}

	c.barrierRekeyConfig = nil
	return results, nil
}

func (c *Core) performBarrierRekey(ctx context.Context, newMasterKey []byte) logical.HTTPCodedError {
	// Rekey the barrier
	if err := c.barrier.Rekey(ctx, newMasterKey); err != nil {
		c.logger.Error("failed to rekey barrier", "error", err)
		return logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to rekey barrier: {{err}}", err).Error())
	}
	if c.logger.IsInfo() {
		c.logger.Info("security barrier rekeyed", "shares", c.barrierRekeyConfig.SecretShares, "threshold", c.barrierRekeyConfig.SecretThreshold)
	}

	c.barrierRekeyConfig.VerificationKey = nil

	if err := c.seal.SetBarrierConfig(ctx, c.barrierRekeyConfig); err != nil {
		c.logger.Error("error saving rekey seal configuration", "error", err)
		return logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to save rekey seal configuration: {{err}}", err).Error())
	}

	// Write to the canary path, which will force a synchronous truing during
	// replication
	if err := c.barrier.Put(ctx, &Entry{
		Key:   coreKeyringCanaryPath,
		Value: []byte(c.barrierRekeyConfig.Nonce),
	}); err != nil {
		c.logger.Error("error saving keyring canary", "error", err)
		return logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to save keyring canary: {{err}}", err).Error())
	}

	c.barrierRekeyConfig.RekeyProgress = nil

	return nil
}

// RecoveryRekeyUpdate is used to provide a new key part
func (c *Core) RecoveryRekeyUpdate(ctx context.Context, key []byte, nonce string) (*RekeyResult, logical.HTTPCodedError) {
	// Ensure we are already unsealed
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, logical.CodedError(http.StatusServiceUnavailable, consts.ErrSealed.Error())
	}
	if c.standby {
		return nil, logical.CodedError(http.StatusBadRequest, consts.ErrStandby.Error())
	}

	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		return nil, logical.CodedError(http.StatusBadRequest, fmt.Sprintf("key is shorter than minimum %d bytes", min))
	}
	if len(key) > max {
		return nil, logical.CodedError(http.StatusBadRequest, fmt.Sprintf("key is longer than maximum %d bytes", max))
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Get the seal configuration
	existingConfig, err := c.seal.RecoveryConfig(ctx)
	if err != nil {
		return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to fetch existing recovery config: {{err}}", err).Error())
	}
	// Ensure the seal is initialized
	if existingConfig == nil {
		return nil, logical.CodedError(http.StatusBadRequest, ErrNotInit.Error())
	}

	// Ensure a rekey is in progress
	if c.recoveryRekeyConfig == nil {
		return nil, logical.CodedError(http.StatusBadRequest, "no recovery rekey in progress")
	}

	if len(c.recoveryRekeyConfig.VerificationKey) > 0 {
		return nil, logical.CodedError(http.StatusBadRequest, fmt.Sprintf("rekey operation already finished; verification must be performed; nonce for the verification operation is %q", c.recoveryRekeyConfig.VerificationNonce))
	}

	if nonce != c.recoveryRekeyConfig.Nonce {
		return nil, logical.CodedError(http.StatusBadRequest, fmt.Sprintf("incorrect nonce supplied; nonce for this rekey operation is %q", c.recoveryRekeyConfig.Nonce))
	}

	// Check if we already have this piece
	for _, existing := range c.recoveryRekeyConfig.RekeyProgress {
		if subtle.ConstantTimeCompare(existing, key) == 1 {
			return nil, logical.CodedError(http.StatusBadRequest, "given key has already been provided during this rekey operation")
		}
	}

	// Store this key
	c.recoveryRekeyConfig.RekeyProgress = append(c.recoveryRekeyConfig.RekeyProgress, key)

	// Check if we don't have enough keys to unlock
	if len(c.recoveryRekeyConfig.RekeyProgress) < existingConfig.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("cannot rekey yet, not enough keys", "keys", len(c.recoveryRekeyConfig.RekeyProgress), "threshold", existingConfig.SecretThreshold)
		}
		return nil, nil
	}

	// Recover the master key
	var recoveryKey []byte
	if existingConfig.SecretThreshold == 1 {
		recoveryKey = c.recoveryRekeyConfig.RekeyProgress[0]
	} else {
		recoveryKey, err = shamir.Combine(c.recoveryRekeyConfig.RekeyProgress)
		if err != nil {
			return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to compute recovery key: {{err}}", err).Error())
		}
	}

	// Verify the recovery key
	if err := c.seal.VerifyRecoveryKey(ctx, recoveryKey); err != nil {
		c.logger.Error("recovery key verification failed", "error", err)
		return nil, logical.CodedError(http.StatusBadRequest, errwrap.Wrapf("recovery key verification failed: {{err}}", err).Error())
	}

	// Generate a new master key
	newMasterKey, err := c.barrier.GenerateKey()
	if err != nil {
		c.logger.Error("failed to generate recovery key", "error", err)
		return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("recovery key generation failed: {{err}}", err).Error())
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
			return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to generate shares: {{err}}", err).Error())
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
			return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to encrypt shares: {{err}}", err).Error())
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
				return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to marshal recovery key backup: {{err}}", err).Error())
			}
			pe := &physical.Entry{
				Key:   coreRecoveryUnsealKeysBackupPath,
				Value: buf,
			}
			if err = c.physical.Put(ctx, pe); err != nil {
				c.logger.Error("failed to save unseal key backup", "error", err)
				return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to save unseal key backup: {{err}}", err).Error())
			}
		}
	}

	// If we are requiring validation, return now; otherwise save the recovery
	// key
	if c.recoveryRekeyConfig.VerificationRequired {
		nonce, err := uuid.GenerateUUID()
		if err != nil {
			c.recoveryRekeyConfig = nil
			return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to generate verification nonce: {{err}}", err).Error())
		}
		c.recoveryRekeyConfig.VerificationNonce = nonce
		c.recoveryRekeyConfig.VerificationKey = newMasterKey

		results.VerificationRequired = true
		results.VerificationNonce = nonce
		return results, nil
	}

	if err := c.performRecoveryRekey(ctx, newMasterKey); err != nil {
		return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to perform recovery rekey: {{err}}", err).Error())
	}

	c.recoveryRekeyConfig = nil
	return results, nil
}

func (c *Core) performRecoveryRekey(ctx context.Context, newMasterKey []byte) logical.HTTPCodedError {
	if err := c.seal.SetRecoveryKey(ctx, newMasterKey); err != nil {
		c.logger.Error("failed to set recovery key", "error", err)
		return logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to set recovery key: {{err}}", err).Error())
	}

	c.recoveryRekeyConfig.VerificationKey = nil

	if err := c.seal.SetRecoveryConfig(ctx, c.recoveryRekeyConfig); err != nil {
		c.logger.Error("error saving rekey seal configuration", "error", err)
		return logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to save rekey seal configuration: {{err}}", err).Error())
	}

	// Write to the canary path, which will force a synchronous truing during
	// replication
	if err := c.barrier.Put(ctx, &Entry{
		Key:   coreKeyringCanaryPath,
		Value: []byte(c.recoveryRekeyConfig.Nonce),
	}); err != nil {
		c.logger.Error("error saving keyring canary", "error", err)
		return logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to save keyring canary: {{err}}", err).Error())
	}

	c.recoveryRekeyConfig.RekeyProgress = nil

	return nil
}

func (c *Core) RekeyVerify(ctx context.Context, key []byte, nonce string, recovery bool) (ret *RekeyVerifyResult, retErr logical.HTTPCodedError) {
	// Ensure we are already unsealed
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, logical.CodedError(http.StatusServiceUnavailable, consts.ErrSealed.Error())
	}
	if c.standby {
		return nil, logical.CodedError(http.StatusBadRequest, consts.ErrStandby.Error())
	}

	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		return nil, logical.CodedError(http.StatusBadRequest, fmt.Sprintf("key is shorter than minimum %d bytes", min))
	}
	if len(key) > max {
		return nil, logical.CodedError(http.StatusBadRequest, fmt.Sprintf("key is longer than maximum %d bytes", max))
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	config := c.barrierRekeyConfig
	if recovery {
		config = c.recoveryRekeyConfig
	}

	// Ensure a rekey is in progress
	if config == nil {
		return nil, logical.CodedError(http.StatusBadRequest, "no rekey in progress")
	}

	if len(config.VerificationKey) == 0 {
		return nil, logical.CodedError(http.StatusBadRequest, "no rekey verification in progress")
	}

	if nonce != config.VerificationNonce {
		return nil, logical.CodedError(http.StatusBadRequest, fmt.Sprintf("incorrect nonce supplied; nonce for this verify operation is %q", config.VerificationNonce))
	}

	// Check if we already have this piece
	for _, existing := range config.VerificationProgress {
		if subtle.ConstantTimeCompare(existing, key) == 1 {
			return nil, logical.CodedError(http.StatusBadRequest, "given key has already been provided during this verify operation")
		}
	}

	// Store this key
	config.VerificationProgress = append(config.VerificationProgress, key)

	// Check if we don't have enough keys to unlock
	if len(config.VerificationProgress) < config.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("cannot verify yet, not enough keys", "keys", len(config.VerificationProgress), "threshold", config.SecretThreshold)
		}
		return nil, nil
	}

	// Schedule the progress for forgetting and rotate the nonce if possible
	defer func() {
		config.VerificationProgress = nil
		if ret != nil && ret.Complete {
			return
		}
		// Not complete, so rotate nonce
		nonce, err := uuid.GenerateUUID()
		if err == nil {
			config.VerificationNonce = nonce
			if ret != nil {
				ret.Nonce = nonce
			}
		}
	}()

	// Recover the master key or recovery key
	var recoveredKey []byte
	if config.SecretThreshold == 1 {
		recoveredKey = config.VerificationProgress[0]
	} else {
		var err error
		recoveredKey, err = shamir.Combine(config.VerificationProgress)
		if err != nil {
			return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to compute key for verification: {{err}}", err).Error())
		}
	}

	if subtle.ConstantTimeCompare(recoveredKey, config.VerificationKey) != 1 {
		c.logger.Error("rekey verification failed")
		return nil, logical.CodedError(http.StatusBadRequest, "rekey verification failed; incorrect key shares supplied")
	}

	switch recovery {
	case false:
		if err := c.performBarrierRekey(ctx, recoveredKey); err != nil {
			return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to perform rekey: {{err}}", err).Error())
		}
		c.barrierRekeyConfig = nil
	default:
		if err := c.performRecoveryRekey(ctx, recoveredKey); err != nil {
			return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("failed to perform recovery key rekey: {{err}}", err).Error())
		}
		c.recoveryRekeyConfig = nil
	}

	res := &RekeyVerifyResult{
		Nonce:    config.VerificationNonce,
		Complete: true,
	}

	return res, nil
}

// RekeyCancel is used to cancel an in-progress rekey
func (c *Core) RekeyCancel(recovery bool) logical.HTTPCodedError {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return logical.CodedError(http.StatusServiceUnavailable, consts.ErrSealed.Error())
	}
	if c.standby {
		return logical.CodedError(http.StatusBadRequest, consts.ErrStandby.Error())
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Clear any progress or config
	if recovery {
		c.recoveryRekeyConfig = nil
	} else {
		c.barrierRekeyConfig = nil
	}
	return nil
}

// RekeyVerifyRestart is used to start the verification process over
func (c *Core) RekeyVerifyRestart(recovery bool) logical.HTTPCodedError {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return logical.CodedError(http.StatusServiceUnavailable, consts.ErrSealed.Error())
	}
	if c.standby {
		return logical.CodedError(http.StatusBadRequest, consts.ErrStandby.Error())
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	// Attempt to generate a new nonce, but don't bail if it doesn't succeed
	// (which is extraordinarily unlikely)
	nonce, nonceErr := uuid.GenerateUUID()

	// Clear any progress or config
	if recovery {
		c.recoveryRekeyConfig.VerificationProgress = nil
		if nonceErr == nil {
			c.recoveryRekeyConfig.VerificationNonce = nonce
		}
	} else {
		c.barrierRekeyConfig.VerificationProgress = nil
		if nonceErr == nil {
			c.barrierRekeyConfig.VerificationNonce = nonce
		}
	}

	return nil
}

// RekeyRetrieveBackup is used to retrieve any backed-up PGP-encrypted unseal
// keys
func (c *Core) RekeyRetrieveBackup(ctx context.Context, recovery bool) (*RekeyBackup, logical.HTTPCodedError) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, logical.CodedError(http.StatusServiceUnavailable, consts.ErrSealed.Error())
	}
	if c.standby {
		return nil, logical.CodedError(http.StatusBadRequest, consts.ErrStandby.Error())
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
		return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("error getting keys from backup: {{err}}", err).Error())
	}
	if entry == nil {
		return nil, nil
	}

	ret := &RekeyBackup{}
	err = jsonutil.DecodeJSON(entry.Value, ret)
	if err != nil {
		return nil, logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("error decoding backup keys: {{err}}", err).Error())
	}

	return ret, nil
}

// RekeyDeleteBackup is used to delete any backed-up PGP-encrypted unseal keys
func (c *Core) RekeyDeleteBackup(ctx context.Context, recovery bool) logical.HTTPCodedError {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return logical.CodedError(http.StatusServiceUnavailable, consts.ErrSealed.Error())
	}
	if c.standby {
		return logical.CodedError(http.StatusBadRequest, consts.ErrStandby.Error())
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	if recovery {
		err := c.physical.Delete(ctx, coreRecoveryUnsealKeysBackupPath)
		if err != nil {
			return logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("error deleting backup keys: {{err}}", err).Error())
		}
		return nil
	}
	err := c.physical.Delete(ctx, coreBarrierUnsealKeysBackupPath)
	if err != nil {
		return logical.CodedError(http.StatusInternalServerError, errwrap.Wrapf("error deleting backup keys: {{err}}", err).Error())
	}
	return nil
}
