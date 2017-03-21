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
func (c *Core) GenerateShareProgress() (int, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return 0, consts.ErrSealed
	}
	if c.standby {
		return 0, consts.ErrStandby
	}

	c.generateShareLock.Lock()
	defer c.generateShareLock.Unlock()

	return len(c.generateShareProgress), nil
}

// GenerateShareConfig is used to read the share generation configuration
func (c *Core) GenerateShareConfiguration() (*GenerateShareConfig, error) {
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
	var conf *GenerateShareConfig
	if c.generateShareConfig != nil {
		conf = new(GenerateShareConfig)
		*conf = *c.generateShareConfig
	}
	return conf, nil
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
func (c *Core) GenerateShareUpdate(key []byte) (*GenerateShareResult, error) {
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
	var config *SealConfig
	var err error
	if c.seal.RecoveryKeySupported() {
		config, err = c.seal.RecoveryConfig()
		if err != nil {
			return nil, err
		}
	} else {
		config, err = c.seal.BarrierConfig()
		if err != nil {
			return nil, err
		}
	}

	// Ensure the barrier is initialized
	if config == nil {
		return nil, ErrNotInit
	}

	// Ensure we are already unsealed
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

	// Ensure a generateShare is in progress
	if c.generateShareConfig == nil {
		return nil, fmt.Errorf("no share generation in progress")
	}

	// Check if we already have this piece
	for _, existing := range c.generateShareProgress {
		if bytes.Equal(existing, key) {
			return nil, fmt.Errorf("given key has already been provided during this generation operation")
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
		return &GenerateShareResult{
			Progress:       progress,
			Required:       config.SecretThreshold,
			PGPFingerprint: c.generateShareConfig.PGPFingerprint,
		}, nil
	}

	if config.SecretThreshold == 1 {
		return nil, fmt.Errorf("key threshold must be greater than 1 to generate additional shares")
	}

	// Recover the master key
	var masterKey []byte
	masterKey, err = shamir.Combine(c.generateShareProgress)
	if err != nil {
		return nil, fmt.Errorf("failed to compute master key: %v", err)
	}

	// TODO: Don't do anything if secret threshold is only 1.

	// Verify the master key
	if c.seal.RecoveryKeySupported() {
		if err := c.seal.VerifyRecoveryKey(masterKey); err != nil {
			c.logger.Error("core: share generation aborted, recovery key verification failed", "error", err)
			return nil, err
		}
	} else {
		if err := c.barrier.VerifyMaster(masterKey); err != nil {
			c.logger.Error("core: share generation aborted, master key verification failed", "error", err)
			return nil, err
		}
	}

	// Generate the new share at position N + 1
	var newShareBytes []byte
	newShareBytes, err = shamir.GetShare(c.generateShareProgress, uint8(config.SecretShares+1))
	if err != nil {
		c.logger.Error("core: share generation aborted, share generated failed", "error", err)
		return nil, err
	}

	// Update barrier config with greater number of shares
	newConfig := config.Clone()
	newConfig.SecretShares++
	if c.seal.SetBarrierConfig(newConfig) != nil {
		c.logger.Error("core: unable to set barrier config", "error", err)
		return nil, err
	}

	// Encrypt the share if a PGP key was given
	if len(c.generateShareConfig.PGPKey) > 0 {
		hexEncodedShares := make([][]byte, 1)
		hexEncodedShares[0] = []byte(base64.StdEncoding.EncodeToString(newShareBytes))
		_, keyBytesArr, err := pgpkeys.EncryptShares(hexEncodedShares, []string{c.generateShareConfig.PGPKey})
		if err != nil {
			c.logger.Error("core: error encrypting new master key share", "error", err)
			return nil, err
		}
		newShareBytes = keyBytesArr[0]
	}

	results := &GenerateShareResult{
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
	return results, nil
}

// GenerateShareCancel is used to cancel an in-progress master key share generation
func (c *Core) GenerateShareCancel() error {
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

	// Clear any progress or config
	c.generateShareConfig = nil
	c.generateShareProgress = nil
	return nil
}
