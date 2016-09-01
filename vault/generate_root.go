package vault

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/shamir"
)

// GenerateRootConfig holds the configuration for a root generation
// command.
type GenerateRootConfig struct {
	Nonce          string
	PGPKey         string
	PGPFingerprint string
	OTP            string
}

// GenerateRootResult holds the result of a root generation update
// command
type GenerateRootResult struct {
	Progress         int
	Required         int
	EncodedRootToken string
	PGPFingerprint   string
}

// GenerateRoot is used to return the root generation progress (num shares)
func (c *Core) GenerateRootProgress() (int, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return 0, ErrSealed
	}
	if c.standby {
		return 0, ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	return len(c.generateRootProgress), nil
}

// GenerateRootConfig is used to read the root generation configuration
// It stubbornly refuses to return the OTP if one is there.
func (c *Core) GenerateRootConfiguration() (*GenerateRootConfig, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, ErrSealed
	}
	if c.standby {
		return nil, ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	// Copy the config if any
	var conf *GenerateRootConfig
	if c.generateRootConfig != nil {
		conf = new(GenerateRootConfig)
		*conf = *c.generateRootConfig
		conf.OTP = ""
	}
	return conf, nil
}

// GenerateRootInit is used to initialize the root generation settings
func (c *Core) GenerateRootInit(otp, pgpKey string) error {
	var fingerprint string
	switch {
	case len(otp) > 0:
		otpBytes, err := base64.StdEncoding.DecodeString(otp)
		if err != nil {
			return fmt.Errorf("error decoding base64 OTP value: %s", err)
		}
		if otpBytes == nil || len(otpBytes) != 16 {
			return fmt.Errorf("decoded OTP value is invalid or wrong length")
		}

	case len(pgpKey) > 0:
		fingerprints, err := pgpkeys.GetFingerprints([]string{pgpKey}, nil)
		if err != nil {
			return fmt.Errorf("error parsing PGP key: %s", err)
		}
		if len(fingerprints) != 1 || fingerprints[0] == "" {
			return fmt.Errorf("could not acquire PGP key entity")
		}
		fingerprint = fingerprints[0]

	default:
		return fmt.Errorf("unreachable condition")
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return ErrSealed
	}
	if c.standby {
		return ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	// Prevent multiple concurrent root generations
	if c.generateRootConfig != nil {
		return fmt.Errorf("root generation already in progress")
	}

	// Copy the configuration
	generationNonce, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}

	c.generateRootConfig = &GenerateRootConfig{
		Nonce:          generationNonce,
		OTP:            otp,
		PGPKey:         pgpKey,
		PGPFingerprint: fingerprint,
	}

	if c.logger.IsInfo() {
		c.logger.Info("core: root generation initialized", "nonce", c.generateRootConfig.Nonce)
	}
	return nil
}

// GenerateRootUpdate is used to provide a new key part
func (c *Core) GenerateRootUpdate(key []byte, nonce string) (*GenerateRootResult, error) {
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
		return nil, ErrSealed
	}
	if c.standby {
		return nil, ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	// Ensure a generateRoot is in progress
	if c.generateRootConfig == nil {
		return nil, fmt.Errorf("no root generation in progress")
	}

	if nonce != c.generateRootConfig.Nonce {
		return nil, fmt.Errorf("incorrect nonce supplied; nonce for this root generation operation is %s", c.generateRootConfig.Nonce)
	}

	// Check if we already have this piece
	for _, existing := range c.generateRootProgress {
		if bytes.Equal(existing, key) {
			return nil, fmt.Errorf("given key has already been provided during this generation operation")
		}
	}

	// Store this key
	c.generateRootProgress = append(c.generateRootProgress, key)
	progress := len(c.generateRootProgress)

	// Check if we don't have enough keys to unlock
	if len(c.generateRootProgress) < config.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("core: cannot generate root, not enough keys", "keys", progress, "threshold", config.SecretThreshold)
		}
		return &GenerateRootResult{
			Progress:       progress,
			Required:       config.SecretThreshold,
			PGPFingerprint: c.generateRootConfig.PGPFingerprint,
		}, nil
	}

	// Recover the master key
	var masterKey []byte
	if config.SecretThreshold == 1 {
		masterKey = c.generateRootProgress[0]
		c.generateRootProgress = nil
	} else {
		masterKey, err = shamir.Combine(c.generateRootProgress)
		c.generateRootProgress = nil
		if err != nil {
			return nil, fmt.Errorf("failed to compute master key: %v", err)
		}
	}

	// Verify the master key
	if c.seal.RecoveryKeySupported() {
		if err := c.seal.VerifyRecoveryKey(masterKey); err != nil {
			c.logger.Error("core: root generation aborted, recovery key verification failed", "error", err)
			return nil, err
		}
	} else {
		if err := c.barrier.VerifyMaster(masterKey); err != nil {
			c.logger.Error("core: root generation aborted, master key verification failed", "error", err)
			return nil, err
		}
	}

	te, err := c.tokenStore.rootToken()
	if err != nil {
		c.logger.Error("core: root token generation failed", "error", err)
		return nil, err
	}
	if te == nil {
		c.logger.Error("core: got nil token entry back from root generation")
		return nil, fmt.Errorf("got nil token entry back from root generation")
	}

	uuidBytes, err := uuid.ParseUUID(te.ID)
	if err != nil {
		c.tokenStore.Revoke(te.ID)
		c.logger.Error("core: error getting generated token bytes", "error", err)
		return nil, err
	}
	if uuidBytes == nil {
		c.tokenStore.Revoke(te.ID)
		c.logger.Error("core: got nil parsed UUID bytes")
		return nil, fmt.Errorf("got nil parsed UUID bytes")
	}

	var tokenBytes []byte
	// Get the encoded value first so that if there is an error we don't create
	// the root token.
	switch {
	case len(c.generateRootConfig.OTP) > 0:
		// This function performs decoding checks so rather than decode the OTP,
		// just encode the value we're passing in.
		tokenBytes, err = xor.XORBase64(c.generateRootConfig.OTP, base64.StdEncoding.EncodeToString(uuidBytes))
		if err != nil {
			c.tokenStore.Revoke(te.ID)
			c.logger.Error("core: xor of root token failed", "error", err)
			return nil, err
		}

	case len(c.generateRootConfig.PGPKey) > 0:
		_, tokenBytesArr, err := pgpkeys.EncryptShares([][]byte{[]byte(te.ID)}, []string{c.generateRootConfig.PGPKey})
		if err != nil {
			c.tokenStore.Revoke(te.ID)
			c.logger.Error("core: error encrypting new root token", "error", err)
			return nil, err
		}
		tokenBytes = tokenBytesArr[0]

	default:
		c.tokenStore.Revoke(te.ID)
		return nil, fmt.Errorf("unreachable condition")
	}

	results := &GenerateRootResult{
		Progress:         progress,
		Required:         config.SecretThreshold,
		EncodedRootToken: base64.StdEncoding.EncodeToString(tokenBytes),
		PGPFingerprint:   c.generateRootConfig.PGPFingerprint,
	}

	if c.logger.IsInfo() {
		c.logger.Info("core: root generation finished", "nonce", c.generateRootConfig.Nonce)
	}

	c.generateRootProgress = nil
	c.generateRootConfig = nil
	return results, nil
}

// GenerateRootCancel is used to cancel an in-progress root generation
func (c *Core) GenerateRootCancel() error {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return ErrSealed
	}
	if c.standby {
		return ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	// Clear any progress or config
	c.generateRootConfig = nil
	c.generateRootProgress = nil
	return nil
}
