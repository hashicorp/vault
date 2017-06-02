package vault

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/shamir"
)

// GenerateShareConfig holds the configuration for a share generation
// command.
type GenerateShareConfig struct {
	PGPKey         string
	PGPFingerprint string
}

// GenerateShareResult holds the result of a share generation update
// command
type GenerateShareResult struct {
	Progress       int
	Required       int
	Key            string
	PGPFingerprint string
}

// GenerateShare is used to return the share generation progress (num shares)
func (c *Core) GenerateShareProgress() (progress int, err error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		progress = 0
		err = consts.ErrSealed
		return
	}
	if c.standby {
		progress = 0
		err = consts.ErrStandby
		return
	}

	c.generateShareLock.Lock()
	defer c.generateShareLock.Unlock()

	progress = len(c.generateShareProgress)
	return
}

// GenerateShareConfig is used to read the share generation configuration
func (c *Core) GenerateShareConfiguration() (conf *GenerateShareConfig, err error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, consts.ErrSealed
	}
	if c.standby {
		return nil, consts.ErrStandby
	}

	c.generateShareLock.Lock()
	defer c.generateShareLock.Unlock()

	// Copy the config if any
	if c.generateShareConfig != nil {
		conf = new(GenerateShareConfig)
		*conf = *c.generateShareConfig
	}
	return
}

// GenerateShareInit is used to initialize the share generation settings
func (c *Core) GenerateShareInit(pgpKey string) error {

	// Get the seal configuration
	var config *SealConfig
	var err error
	if c.seal.RecoveryKeySupported() {
		config, err = c.seal.RecoveryConfig()
		if err != nil {
			return err
		}
	} else {
		config, err = c.seal.BarrierConfig()
		if err != nil {
			return err
		}
	}

	// Ensure the barrier is initialized
	if config == nil {
		return ErrNotInit
	}

	// Ensure key threshold is greater than 1
	if config.SecretThreshold == 1 {
		return fmt.Errorf("key threshold must be greater than 1 to generate additional shares")
	}

	var fingerprint string
	switch {
	case len(pgpKey) > 0:
		fingerprints, err := pgpkeys.GetFingerprints([]string{pgpKey}, nil)
		if err != nil {
			return fmt.Errorf("error parsing PGP key: %s", err)
		}
		if len(fingerprints) != 1 || fingerprints[0] == "" {
			return fmt.Errorf("could not acquire PGP key entity")
		}
		fingerprint = fingerprints[0]
	case len(pgpKey) == 0:
		break
	default:
		return fmt.Errorf("unreachable condition")
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return consts.ErrSealed
	}
	if c.standby {
		return consts.ErrStandby
	}

	c.generateShareLock.Lock()
	defer c.generateShareLock.Unlock()

	// Prevent multiple concurrent share generations
	if c.generateShareConfig != nil {
		return fmt.Errorf("share generation already in progress")
	}

	c.generateShareConfig = &GenerateShareConfig{
		PGPKey:         pgpKey,
		PGPFingerprint: fingerprint,
	}

	if c.logger.IsInfo() {
		c.logger.Info("core: share generation initialized")
	}
	return nil
}

// GenerateShareUpdate is used to provide a new key part
func (c *Core) GenerateShareUpdate(key []byte) (shareResult *GenerateShareResult, err error) {
	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		err = &ErrInvalidKey{fmt.Sprintf("key is shorter than minimum %d bytes", min)}
		return
	}
	if len(key) > max {
		err = &ErrInvalidKey{fmt.Sprintf("key is longer than maximum %d bytes", max)}
		return
	}

	// Get the seal configuration
	var config *SealConfig
	if c.seal.RecoveryKeySupported() {
		config, err = c.seal.RecoveryConfig()
		if err != nil {
			return
		}
	} else {
		config, err = c.seal.BarrierConfig()
		if err != nil {
			return
		}
	}

	// Ensure the barrier is initialized
	if config == nil {
		err = ErrNotInit
		return
	}

	// Ensure we are already unsealed
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		err = consts.ErrSealed
		return
	}
	if c.standby {
		err = consts.ErrStandby
		return
	}

	c.generateShareLock.Lock()
	defer c.generateShareLock.Unlock()

	// Ensure a generateShare is in progress
	if c.generateShareConfig == nil {
		err = fmt.Errorf("no share generation in progress")
		return
	}

	// Check if we already have this piece
	for _, existing := range c.generateShareProgress {
		if bytes.Equal(existing, key) {
			err = fmt.Errorf("given key has already been provided during this generation operation")
			return
		}
	}

	// Store this key
	c.generateShareProgress = append(c.generateShareProgress, key)
	progress := len(c.generateShareProgress)

	// Check if we don't have enough keys to unlock
	if len(c.generateShareProgress) < config.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("core: cannot generate share, not enough keys", "keys", progress, "threshold", config.SecretThreshold)
		}
		shareResult = &GenerateShareResult{
			Progress:       progress,
			Required:       config.SecretThreshold,
			PGPFingerprint: c.generateShareConfig.PGPFingerprint,
		}
		return
	}

	if config.SecretThreshold == 1 {
		err = fmt.Errorf("key threshold must be greater than 1 to generate additional shares")
		return
	}

	// Recover the master key
	var masterKey []byte
	masterKey, err = shamir.Combine(c.generateShareProgress)
	if err != nil {
		err = fmt.Errorf("failed to compute master key: %v", err)
		return
	}

	// Verify the master key
	if c.seal.RecoveryKeySupported() {
		if err = c.seal.VerifyRecoveryKey(masterKey); err != nil {
			c.logger.Error("core: share generation aborted, recovery key verification failed", "error", err)
			return
		}
	} else {
		if err = c.barrier.VerifyMaster(masterKey); err != nil {
			c.logger.Error("core: share generation aborted, master key verification failed", "error", err)
			return
		}
	}

	// Generate the new share at position N + 1
	var newShareBytes []byte
	newShareBytes, err = shamir.GetShareAt(c.generateShareProgress, uint8(config.SecretShares+1))
	if err != nil {
		c.logger.Error("core: share generation aborted, share generated failed", "error", err)
		return
	}

	// Update barrier config with greater number of shares
	newConfig := config.Clone()
	newConfig.SecretShares++
	if c.seal.SetBarrierConfig(newConfig) != nil {
		c.logger.Error("core: unable to set barrier config", "error", err)
		return
	}

	// Encrypt the share if a PGP key was given
	if len(c.generateShareConfig.PGPKey) > 0 {
		hexEncodedShares := make([][]byte, 1)
		hexEncodedShares[0] = []byte(base64.StdEncoding.EncodeToString(newShareBytes))
		_, keyBytesArr, er := pgpkeys.EncryptShares(hexEncodedShares, []string{c.generateShareConfig.PGPKey})
		if er != nil {
			c.logger.Error("core: error encrypting new master key share", "error", er)
			err = er
			return
		}
		newShareBytes = keyBytesArr[0]
	}

	shareResult = &GenerateShareResult{
		Progress:       progress,
		Required:       config.SecretThreshold,
		Key:            base64.StdEncoding.EncodeToString(newShareBytes),
		PGPFingerprint: c.generateShareConfig.PGPFingerprint,
	}

	if c.logger.IsInfo() {
		c.logger.Info("core: master key share generation finished")
	}

	c.generateShareProgress = nil
	c.generateShareConfig = nil
	return
}

// GenerateShareCancel is used to cancel an in-progress master key share generation
func (c *Core) GenerateShareCancel() (err error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		err = consts.ErrSealed
		return
	}
	if c.standby {
		err = consts.ErrStandby
		return
	}

	c.generateShareLock.Lock()
	defer c.generateShareLock.Unlock()

	// Clear any progress or config
	c.generateShareConfig = nil
	c.generateShareProgress = nil
	return
}
