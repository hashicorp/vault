package vault

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/shamir"
)

const (
	// coreUnsealKeysBackupPath is the path used to back upencrypted unseal
	// keys if specified during a rekey operation. This is outside of the
	// barrier.
	coreBarrierUnsealKeysBackupPath = "core/unseal-keys-backup"

	// coreRecoveryUnsealKeysBackupPath is the path used to back upencrypted
	// recovery keys if specified during a rekey operation. This is outside of
	// the barrier.
	coreRecoveryUnsealKeysBackupPath = "core/recovery-keys-backup"
)

// RekeyResult is used to provide the key parts back after
// they are generated as part of the rekey.
type RekeyResult struct {
	SecretShares    [][]byte
	PGPFingerprints []string
	Backup          bool
	RecoveryKey     bool
}

// RekeyBackup stores the backup copy of PGP-encrypted keys
type RekeyBackup struct {
	Nonce string
	Keys  map[string][]string
}

func (c *Core) RekeyThreshold(recovery bool) (int, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return 0, ErrSealed
	}
	if c.standby {
		return 0, ErrStandby
	}

	c.rekeyLock.RLock()
	defer c.rekeyLock.RUnlock()

	var config *SealConfig
	var err error
	if recovery {
		config, err = c.seal.RecoveryConfig()
	} else {
		config, err = c.seal.BarrierConfig()
	}
	if err != nil {
		return 0, err
	}

	return config.SecretThreshold, nil
}

// RekeyProgress is used to return the rekey progress (num shares)
func (c *Core) RekeyProgress(recovery bool) (int, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return 0, ErrSealed
	}
	if c.standby {
		return 0, ErrStandby
	}

	c.rekeyLock.RLock()
	defer c.rekeyLock.RUnlock()

	if recovery {
		return len(c.recoveryRekeyProgress), nil
	}
	return len(c.barrierRekeyProgress), nil
}

// RekeyConfig is used to read the rekey configuration
func (c *Core) RekeyConfig(recovery bool) (*SealConfig, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, ErrSealed
	}
	if c.standby {
		return nil, ErrStandby
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

	// Check if the seal configuration is valid
	if err := config.Validate(); err != nil {
		c.logger.Printf("[ERR] core: invalid rekey seal configuration: %v", err)
		return fmt.Errorf("invalid rekey seal configuration: %v", err)
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return ErrSealed
	}
	if c.standby {
		return ErrStandby
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

	c.logger.Printf("[INFO] core: rekey initialized (nonce: %s, shares: %d, threshold: %d)",
		c.barrierRekeyConfig.Nonce, c.barrierRekeyConfig.SecretShares, c.barrierRekeyConfig.SecretThreshold)
	return nil
}

// RecoveryRekeyInit is used to initialize the rekey settings for the recovery key
func (c *Core) RecoveryRekeyInit(config *SealConfig) error {
	if config.StoredShares > 0 {
		return fmt.Errorf("stored shares not supported by recovery key")
	}

	// Check if the seal configuration is valid
	if err := config.Validate(); err != nil {
		c.logger.Printf("[ERR] core: invalid recovery configuration: %v", err)
		return fmt.Errorf("invalid recovery configuration: %v", err)
	}

	if !c.seal.RecoveryKeySupported() {
		return fmt.Errorf("recovery keys not supported")
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return ErrSealed
	}
	if c.standby {
		return ErrStandby
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

	c.logger.Printf("[INFO] core: rekey initialized (nonce: %s, shares: %d, threshold: %d)",
		c.recoveryRekeyConfig.Nonce, c.recoveryRekeyConfig.SecretShares, c.recoveryRekeyConfig.SecretThreshold)
	return nil
}

func (c *Core) RekeyUpdate(key []byte, nonce string, recovery bool) (*RekeyResult, error) {
	if recovery {
		return c.RecoveryRekeyUpdate(key, nonce)
	}
	return c.BarrierRekeyUpdate(key, nonce)
}

// BarrierRekeyUpdate is used to provide a new key part
func (c *Core) BarrierRekeyUpdate(key []byte, nonce string) (*RekeyResult, error) {
	// Ensure we are already unsealed
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, ErrSealed
	}
	if c.standby {
		return nil, ErrStandby
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
	existingConfig, err := c.seal.BarrierConfig()
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
		return nil, fmt.Errorf("incorrect nonce supplied; nonce for this rekey operation is %s", c.barrierRekeyConfig.Nonce)
	}

	// Check if we already have this piece
	for _, existing := range c.barrierRekeyProgress {
		if bytes.Equal(existing, key) {
			return nil, nil
		}
	}

	// Store this key
	c.barrierRekeyProgress = append(c.barrierRekeyProgress, key)

	// Check if we don't have enough keys to unlock
	if len(c.barrierRekeyProgress) < existingConfig.SecretThreshold {
		c.logger.Printf("[DEBUG] core: cannot rekey, have %d of %d keys",
			len(c.barrierRekeyProgress), existingConfig.SecretThreshold)
		return nil, nil
	}

	// Recover the master key
	var masterKey []byte
	if existingConfig.SecretThreshold == 1 {
		masterKey = c.barrierRekeyProgress[0]
		c.barrierRekeyProgress = nil
	} else {
		masterKey, err = shamir.Combine(c.barrierRekeyProgress)
		c.barrierRekeyProgress = nil
		if err != nil {
			return nil, fmt.Errorf("failed to compute master key: %v", err)
		}
	}

	if err := c.barrier.VerifyMaster(masterKey); err != nil {
		c.logger.Printf("[ERR] core: rekey aborted, master key verification failed: %v", err)
		return nil, err
	}

	// Generate a new master key
	newMasterKey, err := c.barrier.GenerateKey()
	if err != nil {
		c.logger.Printf("[ERR] core: failed to generate master key: %v", err)
		return nil, fmt.Errorf("master key generation failed: %v", err)
	}

	// Return the master key if only a single key part is used
	results := &RekeyResult{
		Backup: c.barrierRekeyConfig.Backup,
	}

	if c.barrierRekeyConfig.SecretShares == 1 {
		results.SecretShares = append(results.SecretShares, newMasterKey)
	} else {
		// Split the master key using the Shamir algorithm
		shares, err := shamir.Split(newMasterKey, c.barrierRekeyConfig.SecretShares, c.barrierRekeyConfig.SecretThreshold)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to generate shares: %v", err)
			return nil, fmt.Errorf("failed to generate shares: %v", err)
		}
		results.SecretShares = shares
	}

	// If we are storing any shares, add them to the shares to store and remove
	// from the returned keys
	var keysToStore [][]byte
	if c.barrierRekeyConfig.StoredShares > 0 {
		for i := 0; i < c.barrierRekeyConfig.StoredShares; i++ {
			keysToStore = append(keysToStore, results.SecretShares[0])
			results.SecretShares = results.SecretShares[1:]
		}
	}

	if len(c.barrierRekeyConfig.PGPKeys) > 0 {
		hexEncodedShares := make([][]byte, len(results.SecretShares))
		for i, _ := range results.SecretShares {
			hexEncodedShares[i] = []byte(hex.EncodeToString(results.SecretShares[i]))
		}
		results.PGPFingerprints, results.SecretShares, err = pgpkeys.EncryptShares(hexEncodedShares, c.barrierRekeyConfig.PGPKeys)
		if err != nil {
			return nil, err
		}

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
				c.logger.Printf("[ERR] core: failed to marshal unseal key backup: %v", err)
				return nil, fmt.Errorf("failed to marshal unseal key backup: %v", err)
			}
			pe := &physical.Entry{
				Key:   coreBarrierUnsealKeysBackupPath,
				Value: buf,
			}
			if err = c.physical.Put(pe); err != nil {
				c.logger.Printf("[ERR] core: failed to save unseal key backup: %v", err)
				return nil, fmt.Errorf("failed to save unseal key backup: %v", err)
			}
		}
	}

	if keysToStore != nil {
		if err := c.seal.SetStoredKeys(keysToStore); err != nil {
			c.logger.Printf("[ERR] core: failed to store keys: %v", err)
			return nil, fmt.Errorf("failed to store keys: %v", err)
		}
	}

	// Rekey the barrier
	if err := c.barrier.Rekey(newMasterKey); err != nil {
		c.logger.Printf("[ERR] core: failed to rekey barrier: %v", err)
		return nil, fmt.Errorf("failed to rekey barrier: %v", err)
	}
	c.logger.Printf("[INFO] core: security barrier rekeyed (shares: %d, threshold: %d)",
		c.barrierRekeyConfig.SecretShares, c.barrierRekeyConfig.SecretThreshold)

	if err := c.seal.SetBarrierConfig(c.barrierRekeyConfig); err != nil {
		c.logger.Printf("[ERR] core: error saving rekey seal configuration: %v", err)
		return nil, fmt.Errorf("failed to save rekey seal configuration: %v", err)
	}

	// Done!
	c.barrierRekeyProgress = nil
	c.barrierRekeyConfig = nil
	return results, nil
}

// RecoveryRekeyUpdate is used to provide a new key part
func (c *Core) RecoveryRekeyUpdate(key []byte, nonce string) (*RekeyResult, error) {
	// Ensure we are already unsealed
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, ErrSealed
	}
	if c.standby {
		return nil, ErrStandby
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
	barrierConfig, err := c.seal.BarrierConfig()
	if err != nil {
		return nil, err
	}

	// Ensure the barrier is initialized
	if barrierConfig == nil {
		return nil, ErrNotInit
	}

	existingConfig, err := c.seal.RecoveryConfig()
	if err != nil {
		return nil, err
	}

	// Ensure a rekey is in progress
	if c.recoveryRekeyConfig == nil {
		return nil, fmt.Errorf("no rekey in progress")
	}

	if nonce != c.recoveryRekeyConfig.Nonce {
		return nil, fmt.Errorf("incorrect nonce supplied; nonce for this rekey operation is %s", c.recoveryRekeyConfig.Nonce)
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
		c.logger.Printf("[DEBUG] core: cannot rekey, have %d of %d keys",
			len(c.recoveryRekeyProgress), existingConfig.SecretThreshold)
		return nil, nil
	}

	// Recover the master key
	var masterKey []byte
	if existingConfig.SecretThreshold == 1 {
		masterKey = c.recoveryRekeyProgress[0]
		c.recoveryRekeyProgress = nil
	} else {
		masterKey, err = shamir.Combine(c.recoveryRekeyProgress)
		c.recoveryRekeyProgress = nil
		if err != nil {
			return nil, fmt.Errorf("failed to compute recovery key: %v", err)
		}
	}

	// Verify the recovery key
	if err := c.seal.VerifyRecoveryKey(masterKey); err != nil {
		c.logger.Printf("[ERR] core: rekey aborted, recovery key verification failed: %v", err)
		return nil, err
	}

	// Generate a new master key
	newMasterKey, err := c.barrier.GenerateKey()
	if err != nil {
		c.logger.Printf("[ERR] core: failed to generate recovery key: %v", err)
		return nil, fmt.Errorf("recovery key generation failed: %v", err)
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
			c.logger.Printf("[ERR] core: failed to generate shares: %v", err)
			return nil, fmt.Errorf("failed to generate shares: %v", err)
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
				c.logger.Printf("[ERR] core: failed to marshal recovery key backup: %v", err)
				return nil, fmt.Errorf("failed to marshal recovery key backup: %v", err)
			}
			pe := &physical.Entry{
				Key:   coreRecoveryUnsealKeysBackupPath,
				Value: buf,
			}
			if err = c.physical.Put(pe); err != nil {
				c.logger.Printf("[ERR] core: failed to save unseal key backup: %v", err)
				return nil, fmt.Errorf("failed to save unseal key backup: %v", err)
			}
		}
	}

	if err := c.seal.SetRecoveryKey(newMasterKey); err != nil {
		c.logger.Printf("[ERR] core: failed to set recovery key: %v", err)
		return nil, fmt.Errorf("failed to set recovery key: %v", err)
	}

	if err := c.seal.SetRecoveryConfig(c.recoveryRekeyConfig); err != nil {
		c.logger.Printf("[ERR] core: error saving rekey seal configuration: %v", err)
		return nil, fmt.Errorf("failed to save rekey seal configuration: %v", err)
	}

	// Done!
	c.recoveryRekeyProgress = nil
	c.recoveryRekeyConfig = nil
	return results, nil
}

// RekeyCancel is used to cancel an inprogress rekey
func (c *Core) RekeyCancel(recovery bool) error {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return ErrSealed
	}
	if c.standby {
		return ErrStandby
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

// RekeyRetrieveBackup is used to retrieve any backed-up PGP-encrypted unseal
// keys
func (c *Core) RekeyRetrieveBackup(recovery bool) (*RekeyBackup, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, ErrSealed
	}
	if c.standby {
		return nil, ErrStandby
	}

	c.rekeyLock.RLock()
	defer c.rekeyLock.RUnlock()

	var entry *physical.Entry
	var err error
	if recovery {
		entry, err = c.physical.Get(coreRecoveryUnsealKeysBackupPath)
	} else {
		entry, err = c.physical.Get(coreBarrierUnsealKeysBackupPath)
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
func (c *Core) RekeyDeleteBackup(recovery bool) error {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return ErrSealed
	}
	if c.standby {
		return ErrStandby
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	if recovery {
		return c.physical.Delete(coreRecoveryUnsealKeysBackupPath)
	}
	return c.physical.Delete(coreBarrierUnsealKeysBackupPath)
}
