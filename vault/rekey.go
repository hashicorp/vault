package vault

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/helper/salt"
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
	SecretShares         [][]byte
	PGPFingerprints      []string
	Backup               bool
	RecoveryKey          bool
	SecretSharesMetadata []*KeyShareMetadata
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
		return 0, consts.ErrSealed
	}
	if c.standby {
		return 0, consts.ErrStandby
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
		c.logger.Error("core: invalid rekey seal configuration", "error", err)
		return fmt.Errorf("invalid rekey seal configuration: %v", err)
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
		c.logger.Info("core: rekey initialized", "nonce", c.barrierRekeyConfig.Nonce, "shares", c.barrierRekeyConfig.SecretShares, "threshold", c.barrierRekeyConfig.SecretThreshold)
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
		c.logger.Error("core: invalid recovery configuration", "error", err)
		return fmt.Errorf("invalid recovery configuration: %v", err)
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
		c.logger.Info("core: rekey initialized", "nonce", c.recoveryRekeyConfig.Nonce, "shares", c.recoveryRekeyConfig.SecretShares, "threshold", c.recoveryRekeyConfig.SecretThreshold)
	}
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
			return nil, fmt.Errorf("given key has already been provided during this generation operation")
		}
	}

	// Store this key
	c.barrierRekeyProgress = append(c.barrierRekeyProgress, key)

	// Check if we don't have enough keys to unlock
	if len(c.barrierRekeyProgress) < existingConfig.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("core: cannot rekey yet, not enough keys", "keys", len(c.barrierRekeyProgress), "threshold", existingConfig.SecretThreshold)
		}
		return nil, nil
	}

	// Recover the master key
	var masterKey []byte
	if existingConfig.SecretThreshold == 1 {
		masterKey = c.barrierRekeyProgress[0]
	} else {
		masterKey, err = shamir.Combine(c.barrierRekeyProgress)
		if err != nil {
			return nil, fmt.Errorf("failed to compute master key: %v", err)
		}
	}

	if err := c.barrier.VerifyMaster(masterKey); err != nil {
		c.logger.Error("core: rekey aborted, master key verification failed", "error", err)
		return nil, err
	}

	// Fetch the unseal keys metadata and log which of the unseal key holders
	// performed the rekey operation
	keySharesMetadataEntry, err := c.barrier.Get(coreSecretSharesMetadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch unseal key shares metadata entry: %v", err)
	}

	var secretSharesMetadataValue keySharesMetadataStorageValue

	// Log only if the metadata is available
	if keySharesMetadataEntry != nil {
		if err = jsonutil.DecodeJSON(keySharesMetadataEntry.Value, &secretSharesMetadataValue); err != nil {
			return nil, fmt.Errorf("failed to decode unseal key shares metadata entry: %v", err)
		}

		for _, unlockPart := range c.barrierRekeyProgress {
			// Fetch the metadata associated to the unseal key share
			secretShareMetadata, ok := secretSharesMetadataValue.Data[base64.StdEncoding.EncodeToString(salt.SHA256Hash(unlockPart))]

			// If the storage entry is successfully read, metadata associated
			// with all the unseal keys must be available.
			if !ok || secretShareMetadata == nil {
				c.logger.Error("core: failed to fetch unseal key shares metadata")
				return nil, fmt.Errorf("failed to fetch unseal key shares metadata")
			}

			switch {
			case secretShareMetadata.ID != "" && secretShareMetadata.Name != "":
				c.logger.Info(fmt.Sprintf("core: unseal key share with identifier %q with name %q supplied for rekeying", secretShareMetadata.ID, secretShareMetadata.Name))
			case secretShareMetadata.ID != "":
				c.logger.Info(fmt.Sprintf("core: unseal key share with identifier %q supplied for rekeying", secretShareMetadata.ID))
			default:
				c.logger.Error("core: missing unseal key share metadata while rekeying")
				return nil, fmt.Errorf("missing unseal key share metadata while rekeying")
			}
		}
	}

	c.barrierRekeyProgress = nil

	// Generate a new master key
	newMasterKey, err := c.barrier.GenerateKey()
	if err != nil {
		c.logger.Error("core: failed to generate master key", "error", err)
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
			c.logger.Error("core: failed to generate shares", "error", err)
			return nil, fmt.Errorf("failed to generate shares: %v", err)
		}
		results.SecretShares = shares
	}

	// Cache the unencrypted key shares to associate metadata to each
	keyShares := results.SecretShares

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
				c.logger.Error("core: failed to marshal unseal key backup", "error", err)
				return nil, fmt.Errorf("failed to marshal unseal key backup: %v", err)
			}
			pe := &physical.Entry{
				Key:   coreBarrierUnsealKeysBackupPath,
				Value: buf,
			}
			if err = c.physical.Put(pe); err != nil {
				c.logger.Error("core: failed to save unseal key backup", "error", err)
				return nil, fmt.Errorf("failed to save unseal key backup: %v", err)
			}
		}
	}

	var secretSharesMetadataJSON []byte
	// Associate metadata for all the unseal key shares
	secretSharesMetadataJSON, results.SecretSharesMetadata, err = c.prepareKeySharesMetadata(keyShares, c.barrierRekeyConfig.SecretSharesIdentifierNames)
	if err != nil {
		c.logger.Error("core: failed to prepare unseal key shares metadata during rekey", "error", err)
		return nil, fmt.Errorf("failed to prepare unseal key shares metadata during rekey")
	}

	err = c.barrier.Put(&Entry{
		Key:   coreSecretSharesMetadataPath,
		Value: secretSharesMetadataJSON,
	})
	if err != nil {
		c.logger.Error("core: failed to store unseal key shares metadata", "error", err)
		return nil, err
	}

	// If we are storing any shares, add them to the shares to store and remove
	// from the returned keys
	var keysToStore [][]byte
	if c.barrierRekeyConfig.StoredShares > 0 {
		if len(c.barrierRekeyConfig.PGPKeys) > 0 {
			c.logger.Error("core: PGP keys not supported when storing shares")
			return nil, fmt.Errorf("PGP keys not supported when storing shares")
		}

		// Note that results.SecretShares will always be unencrypted here
		for i := 0; i < c.barrierRekeyConfig.StoredShares; i++ {
			keysToStore = append(keysToStore, results.SecretShares[0])
			results.SecretShares = results.SecretShares[1:]
		}

		if err := c.seal.SetStoredKeys(keysToStore); err != nil {
			c.logger.Error("core: failed to store keys", "error", err)
			return nil, fmt.Errorf("failed to store keys: %v", err)
		}
	}

	// Rekey the barrier
	if err := c.barrier.Rekey(newMasterKey); err != nil {
		c.logger.Error("core: failed to rekey barrier", "error", err)
		return nil, fmt.Errorf("failed to rekey barrier: %v", err)
	}
	if c.logger.IsInfo() {
		c.logger.Info("core: security barrier rekeyed", "shares", c.barrierRekeyConfig.SecretShares, "threshold", c.barrierRekeyConfig.SecretThreshold)
	}
	if err := c.seal.SetBarrierConfig(c.barrierRekeyConfig); err != nil {
		c.logger.Error("core: error saving rekey seal configuration", "error", err)
		return nil, fmt.Errorf("failed to save rekey seal configuration: %v", err)
	}

	// Write to the canary path, which will force a synchronous truing during
	// replication
	if err := c.barrier.Put(&Entry{
		Key:   coreKeyringCanaryPath,
		Value: []byte(c.barrierRekeyConfig.Nonce),
	}); err != nil {
		c.logger.Error("core: error saving keyring canary", "error", err)
		return nil, fmt.Errorf("failed to save keyring canary: %v", err)
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
		if c.logger.IsDebug() {
			c.logger.Debug("core: cannot rekey yet, not enough keys", "keys", len(c.recoveryRekeyProgress), "threshold", existingConfig.SecretThreshold)
		}
		return nil, nil
	}

	// Recover the master key
	var masterKey []byte
	if existingConfig.SecretThreshold == 1 {
		masterKey = c.recoveryRekeyProgress[0]
	} else {
		masterKey, err = shamir.Combine(c.recoveryRekeyProgress)
		if err != nil {
			return nil, fmt.Errorf("failed to compute recovery key: %v", err)
		}
	}

	// Verify the recovery key
	if err := c.seal.VerifyRecoveryKey(masterKey); err != nil {
		c.logger.Error("core: rekey aborted, recovery key verification failed", "error", err)
		return nil, err
	}

	// Fetch the recovery keys metadata and log which of the recovery key holders
	// performed the recovery rekey operation
	keySharesMetadataEntry, err := c.barrier.Get(coreRecoverySharesMetadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recovery key shares metadata entry: %v", err)
	}

	var recoverySharesMetadataValue keySharesMetadataStorageValue

	// Log only if the metadata is available
	if keySharesMetadataEntry != nil {
		if err = jsonutil.DecodeJSON(keySharesMetadataEntry.Value, &recoverySharesMetadataValue); err != nil {
			return nil, fmt.Errorf("failed to decode recovery key shares metadata entry: %v", err)
		}

		for _, recoveryKeyPart := range c.recoveryRekeyProgress {
			// Fetch the metadata associated to the recovery key share
			recoveryShareMetadata, ok := recoverySharesMetadataValue.Data[base64.StdEncoding.EncodeToString(salt.SHA256Hash(recoveryKeyPart))]

			// If the storage entry is successfully read, metadata associated
			// with all the key shares must be available.
			if !ok || recoveryShareMetadata == nil {
				c.logger.Error("core: failed to fetch recovery key shares metadata")
				return nil, fmt.Errorf("failed to fetch recovery key shares metadata")
			}

			switch {
			case recoveryShareMetadata.ID != "" && recoveryShareMetadata.Name != "":
				c.logger.Info(fmt.Sprintf("core: recovery key share identifier %q with name %q supplied for rekeying", recoveryShareMetadata.ID, recoveryShareMetadata.Name))
			case recoveryShareMetadata.ID != "":
				c.logger.Info(fmt.Sprintf("core: recovery key share with identifier %q supplied for rekeying", recoveryShareMetadata.ID))
			default:
				c.logger.Error("core: missing recovery key share metadata while rekeying")
				return nil, fmt.Errorf("missing recovery key share metadata while rekeying")
			}
		}
	}

	c.recoveryRekeyProgress = nil

	// Generate a new master key
	newMasterKey, err := c.barrier.GenerateKey()
	if err != nil {
		c.logger.Error("core: failed to generate recovery key", "error", err)
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
			c.logger.Error("core: failed to generate shares", "error", err)
			return nil, fmt.Errorf("failed to generate shares: %v", err)
		}
		results.SecretShares = shares
	}

	// Cache the unencrypted key shares to associate metadata to each
	keyShares := results.SecretShares

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
				c.logger.Error("core: failed to marshal recovery key backup", "error", err)
				return nil, fmt.Errorf("failed to marshal recovery key backup: %v", err)
			}
			pe := &physical.Entry{
				Key:   coreRecoveryUnsealKeysBackupPath,
				Value: buf,
			}
			if err = c.physical.Put(pe); err != nil {
				c.logger.Error("core: failed to save unseal key backup", "error", err)
				return nil, fmt.Errorf("failed to save unseal key backup: %v", err)
			}
		}
	}

	var recoverySharesMetadataJSON []byte
	// Associate metadata for all the recovery key shares
	recoverySharesMetadataJSON, results.SecretSharesMetadata, err = c.prepareKeySharesMetadata(keyShares, c.recoveryRekeyConfig.SecretSharesIdentifierNames)
	if err != nil {
		c.logger.Error("core: failed to prepare recovery key shares metadata during rekey", "error", err)
		return nil, fmt.Errorf("failed to prepare recovery key shares metadata during rekey")
	}

	err = c.barrier.Put(&Entry{
		Key:   coreRecoverySharesMetadataPath,
		Value: recoverySharesMetadataJSON,
	})
	if err != nil {
		c.logger.Error("core: failed to store rekey key shares metadata", "error", err)
		return nil, err
	}

	if err := c.seal.SetRecoveryKey(newMasterKey); err != nil {
		c.logger.Error("core: failed to set recovery key", "error", err)
		return nil, fmt.Errorf("failed to set recovery key: %v", err)
	}

	if err := c.seal.SetRecoveryConfig(c.recoveryRekeyConfig); err != nil {
		c.logger.Error("core: error saving rekey seal configuration", "error", err)
		return nil, fmt.Errorf("failed to save rekey seal configuration: %v", err)
	}

	// Write to the canary path, which will force a synchronous truing during
	// replication
	if err := c.barrier.Put(&Entry{
		Key:   coreKeyringCanaryPath,
		Value: []byte(c.recoveryRekeyConfig.Nonce),
	}); err != nil {
		c.logger.Error("core: error saving keyring canary", "error", err)
		return nil, fmt.Errorf("failed to save keyring canary: %v", err)
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

// RekeyRetrieveBackup is used to retrieve any backed-up PGP-encrypted unseal
// keys
func (c *Core) RekeyRetrieveBackup(recovery bool) (*RekeyBackup, error) {
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
		return consts.ErrSealed
	}
	if c.standby {
		return consts.ErrStandby
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()

	if recovery {
		return c.physical.Delete(coreRecoveryUnsealKeysBackupPath)
	}
	return c.physical.Delete(coreBarrierUnsealKeysBackupPath)
}
