package vault

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/shamir"
)

// RootGenerationConfig holds the configuration for a root generation
// command.
type RootGenerationConfig struct {
	Nonce          string
	PGPKey         string
	PGPFingerprint string
	OTP            string
}

// RootGenerationResult holds the result of a root generation update
// command
type RootGenerationResult struct {
	Progress         int
	Required         int
	EncodedRootToken string
	PGPFingerprint   string
}

// RootGeneration is used to return the root generation progress (num shares)
func (c *Core) RootGenerationProgress() (int, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return 0, ErrSealed
	}
	if c.standby {
		return 0, ErrStandby
	}

	c.rootGenerationLock.Lock()
	defer c.rootGenerationLock.Unlock()

	return len(c.rootGenerationProgress), nil
}

// RootGenerationConfig is used to read the root generation configuration
// It stubbornly refuses to return the OTP if one is there.
func (c *Core) RootGenerationConfiguration() (*RootGenerationConfig, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, ErrSealed
	}
	if c.standby {
		return nil, ErrStandby
	}

	c.rootGenerationLock.Lock()
	defer c.rootGenerationLock.Unlock()

	// Copy the config if any
	var conf *RootGenerationConfig
	if c.rootGenerationConfig != nil {
		conf = new(RootGenerationConfig)
		*conf = *c.rootGenerationConfig
		conf.OTP = ""
	}
	return conf, nil
}

// RootGenerationInit is used to initialize the root generation settings
func (c *Core) RootGenerationInit(otp, pgpKey string) error {
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

	c.rootGenerationLock.Lock()
	defer c.rootGenerationLock.Unlock()

	// Prevent multiple concurrent root generations
	if c.rootGenerationConfig != nil {
		return fmt.Errorf("root generation already in progress")
	}

	// Copy the configuration
	generationNonce, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}

	c.rootGenerationConfig = &RootGenerationConfig{
		Nonce:          generationNonce,
		OTP:            otp,
		PGPKey:         pgpKey,
		PGPFingerprint: fingerprint,
	}

	c.logger.Printf("[INFO] core: root generation initialized (nonce: %s)",
		c.rootGenerationConfig.Nonce)
	return nil
}

// RootGenerationUpdate is used to provide a new key part
func (c *Core) RootGenerationUpdate(key []byte, nonce string) (*RootGenerationResult, error) {
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

	c.rootGenerationLock.Lock()
	defer c.rootGenerationLock.Unlock()

	// Ensure a rootGeneration is in progress
	if c.rootGenerationConfig == nil {
		return nil, fmt.Errorf("no root generation in progress")
	}

	if nonce != c.rootGenerationConfig.Nonce {
		return nil, fmt.Errorf("incorrect nonce supplied; nonce for this root generation operation is %s", c.rootGenerationConfig.Nonce)
	}

	// Check if we already have this piece
	for _, existing := range c.rootGenerationProgress {
		if bytes.Equal(existing, key) {
			return nil, nil
		}
	}

	// Store this key
	c.rootGenerationProgress = append(c.rootGenerationProgress, key)
	progress := len(c.rootGenerationProgress)

	// Check if we don't have enough keys to unlock
	if len(c.rootGenerationProgress) < config.SecretThreshold {
		c.logger.Printf("[DEBUG] core: cannot generate root, have %d of %d keys",
			progress, config.SecretThreshold)
		return &RootGenerationResult{
			Progress:       progress,
			Required:       config.SecretThreshold,
			PGPFingerprint: c.rootGenerationConfig.PGPFingerprint,
		}, nil
	}

	// Recover the master key
	var masterKey []byte
	if config.SecretThreshold == 1 {
		masterKey = c.rootGenerationProgress[0]
		c.rootGenerationProgress = nil
	} else {
		masterKey, err = shamir.Combine(c.rootGenerationProgress)
		c.rootGenerationProgress = nil
		if err != nil {
			return nil, fmt.Errorf("failed to compute master key: %v", err)
		}
	}

	// Verify the master key
	if err := c.barrier.VerifyMaster(masterKey); err != nil {
		c.logger.Printf("[ERR] core: root generation aborted, master key verification failed: %v", err)
		return nil, err
	}

	// Generate the raw token bytes
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		c.logger.Printf("failed to read random bytes: %v", err)
		return nil, err
	}

	uuidStr, err := uuid.FormatUUID(buf)
	if err != nil {
		c.logger.Printf("error formatting token: %v", err)
		return nil, err
	}

	var tokenBytes []byte
	// Get the encoded value first so that if there is an error we don't create
	// the root token.
	switch {
	case len(c.rootGenerationConfig.OTP) > 0:
		// This function performs decoding checks so rather than decode the OTP,
		// just encode the value we're passing in.
		tokenBytes, err = xor.XORBase64(c.rootGenerationConfig.OTP, base64.StdEncoding.EncodeToString(buf))
		if err != nil {
			c.logger.Printf("[ERR] core: xor of root token failed: %v", err)
			return nil, err
		}

	case len(c.rootGenerationConfig.PGPKey) > 0:
		_, tokenBytesArr, err := pgpkeys.EncryptShares([][]byte{[]byte(uuidStr)}, []string{c.rootGenerationConfig.PGPKey})
		if err != nil {
			c.logger.Printf("[ERR] core: error encrypting new root token: %v", err)
			return nil, err
		}
		tokenBytes = tokenBytesArr[0]

	default:
		return nil, fmt.Errorf("unreachable condition")
	}

	_, err = c.tokenStore.rootToken(uuidStr)
	if err != nil {
		c.logger.Printf("[ERR] core: root token generation failed: %v", err)
		return nil, err
	}

	results := &RootGenerationResult{
		Progress:         progress,
		Required:         config.SecretThreshold,
		EncodedRootToken: base64.StdEncoding.EncodeToString(tokenBytes),
		PGPFingerprint:   c.rootGenerationConfig.PGPFingerprint,
	}

	c.logger.Printf("[INFO] core: root generation finished (nonce: %s)",
		c.rootGenerationConfig.Nonce)

	c.rootGenerationProgress = nil
	c.rootGenerationConfig = nil
	return results, nil
}

// RootGenerationCancel is used to cancel an in-progress root generation
func (c *Core) RootGenerationCancel() error {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return ErrSealed
	}
	if c.standby {
		return ErrStandby
	}

	c.rootGenerationLock.Lock()
	defer c.rootGenerationLock.Unlock()

	// Clear any progress or config
	c.rootGenerationConfig = nil
	c.rootGenerationProgress = nil
	return nil
}
