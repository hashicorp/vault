package vault

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/shamir"
)

// RekeyResult is used to provide the key parts back after
// they are generated as part of the rekey.
type RekeyResult struct {
	SecretShares    [][]byte
	PGPFingerprints []string
	Backup          bool
}

// RekeyBackup stores the backup copy of PGP-encrypted keys
type RekeyBackup struct {
	Nonce string
	Keys  map[string][]string
}

// RekeyProgress is used to return the rekey progress (num shares)
func (c *Core) RekeyProgress() (int, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return 0, ErrSealed
	}
	if c.standby {
		return 0, ErrStandby
	}

	c.rekeyLock.Lock()
	defer c.rekeyLock.Unlock()
	return len(c.rekeyProgress), nil
}

// RekeyConfig is used to read the rekey configuration
func (c *Core) RekeyConfig() (*SealConfig, error) {
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
	if c.rekeyConfig != nil {
		conf = new(SealConfig)
		*conf = *c.rekeyConfig
	}
	return conf, nil
}

// RekeyInit is used to initialize the rekey settings
func (c *Core) RekeyInit(config *SealConfig) error {
	// Check if the seal configuraiton is valid
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

	// Prevent multiple concurrent re-keys
	if c.rekeyConfig != nil {
		return fmt.Errorf("rekey already in progress")
	}

	// Copy the configuration
	c.rekeyConfig = new(SealConfig)
	*c.rekeyConfig = *config

	// Initialize the nonce
	nonce, err := uuid.GenerateUUID()
	if err != nil {
		c.rekeyConfig = nil
		return err
	}
	c.rekeyConfig.Nonce = nonce

	c.logger.Printf("[INFO] core: rekey initialized (nonce: %s, shares: %d, threshold: %d)",
		c.rekeyConfig.Nonce, c.rekeyConfig.SecretShares, c.rekeyConfig.SecretThreshold)
	return nil
}

// RekeyUpdate is used to provide a new key part
func (c *Core) RekeyUpdate(key []byte, nonce string) (*RekeyResult, error) {
	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		return nil, &ErrInvalidKey{fmt.Sprintf("key is shorter than minimum %d bytes", min)}
	}
	if len(key) > max {
		return nil, &ErrInvalidKey{fmt.Sprintf("key is longer than maximum %d bytes", max)}
	}

	// Get the seal configuration
	config, err := c.SealConfig()
	if err != nil {
		return nil, err
	}

	// Ensure the barrier is initialized
	if config == nil {
		return nil, ErrNotInit
	}

	// Ensure we are already unsealed
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

	// Ensure a rekey is in progress
	if c.rekeyConfig == nil {
		return nil, fmt.Errorf("no rekey in progress")
	}

	if nonce != c.rekeyConfig.Nonce {
		return nil, fmt.Errorf("incorrect nonce supplied; nonce for this rekey operation is %s", c.rekeyConfig.Nonce)
	}

	// Check if we already have this piece
	for _, existing := range c.rekeyProgress {
		if bytes.Equal(existing, key) {
			return nil, nil
		}
	}

	// Store this key
	c.rekeyProgress = append(c.rekeyProgress, key)

	// Check if we don't have enough keys to unlock
	if len(c.rekeyProgress) < config.SecretThreshold {
		c.logger.Printf("[DEBUG] core: cannot rekey, have %d of %d keys",
			len(c.rekeyProgress), config.SecretThreshold)
		return nil, nil
	}

	// Recover the master key
	var masterKey []byte
	if config.SecretThreshold == 1 {
		masterKey = c.rekeyProgress[0]
		c.rekeyProgress = nil
	} else {
		masterKey, err = shamir.Combine(c.rekeyProgress)
		c.rekeyProgress = nil
		if err != nil {
			return nil, fmt.Errorf("failed to compute master key: %v", err)
		}
	}

	// Verify the master key
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
		Backup: c.rekeyConfig.Backup,
	}

	if c.rekeyConfig.SecretShares == 1 {
		results.SecretShares = append(results.SecretShares, newMasterKey)
	} else {
		// Split the master key using the Shamir algorithm
		shares, err := shamir.Split(newMasterKey, c.rekeyConfig.SecretShares, c.rekeyConfig.SecretThreshold)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to generate shares: %v", err)
			return nil, fmt.Errorf("failed to generate shares: %v", err)
		}
		results.SecretShares = shares
	}

	if len(c.rekeyConfig.PGPKeys) > 0 {
		hexEncodedShares := make([][]byte, len(results.SecretShares))
		for i, _ := range results.SecretShares {
			hexEncodedShares[i] = []byte(hex.EncodeToString(results.SecretShares[i]))
		}
		results.PGPFingerprints, results.SecretShares, err = pgpkeys.EncryptShares(hexEncodedShares, c.rekeyConfig.PGPKeys)
		if err != nil {
			return nil, err
		}

		if c.rekeyConfig.Backup {
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
				Nonce: c.rekeyConfig.Nonce,
				Keys:  backupInfo,
			}
			buf, err := json.Marshal(backupVals)
			if err != nil {
				c.logger.Printf("[ERR] core: failed to marshal unseal key backup: %v", err)
				return nil, fmt.Errorf("failed to marshal unseal key backup: %v", err)
			}
			pe := &physical.Entry{
				Key:   coreUnsealKeysBackupPath,
				Value: buf,
			}
			if err = c.physical.Put(pe); err != nil {
				c.logger.Printf("[ERR] core: failed to save unseal key backup: %v", err)
				return nil, fmt.Errorf("failed to save unseal key backup: %v", err)
			}
		}
	}

	// Encode the seal configuration
	buf, err := json.Marshal(c.rekeyConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode seal configuration: %v", err)
	}

	// Rekey the barrier
	if err := c.barrier.Rekey(newMasterKey); err != nil {
		c.logger.Printf("[ERR] core: failed to rekey barrier: %v", err)
		return nil, fmt.Errorf("failed to rekey barrier: %v", err)
	}
	c.logger.Printf("[INFO] core: security barrier rekeyed (shares: %d, threshold: %d)",
		c.rekeyConfig.SecretShares, c.rekeyConfig.SecretThreshold)

	// Store the seal configuration
	pe := &physical.Entry{
		Key:   coreSealConfigPath,
		Value: buf,
	}
	if err := c.physical.Put(pe); err != nil {
		c.logger.Printf("[ERR] core: failed to update seal configuration: %v", err)
		return nil, fmt.Errorf("failed to update seal configuration: %v", err)
	}

	// Done!
	c.rekeyProgress = nil
	c.rekeyConfig = nil
	return results, nil
}

// RekeyCancel is used to cancel an inprogress rekey
func (c *Core) RekeyCancel() error {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return ErrSealed
	}
	if c.standby {
		return ErrStandby
	}

	// Clear any progress or config
	c.rekeyConfig = nil
	c.rekeyProgress = nil
	return nil
}

// RekeyRetrieveBackup is used to retrieve any backed-up PGP-encrypted unseal
// keys
func (c *Core) RekeyRetrieveBackup() (*RekeyBackup, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, ErrSealed
	}
	if c.standby {
		return nil, ErrStandby
	}

	entry, err := c.physical.Get(coreUnsealKeysBackupPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	ret := &RekeyBackup{}
	err = json.Unmarshal(entry.Value, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// RekeyDeleteBackup is used to delete any backed-up PGP-encrypted unseal keys
func (c *Core) RekeyDeleteBackup() error {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return ErrSealed
	}
	if c.standby {
		return ErrStandby
	}

	return c.physical.Delete(coreUnsealKeysBackupPath)
}
