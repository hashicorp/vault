// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/roottoken"
	"github.com/hashicorp/vault/shamir"
)

const coreDROperationTokenPath = "core/dr-operation-token"

var (
	// GenerateStandardRootTokenStrategy is the strategy used to generate a
	// typical root token
	GenerateStandardRootTokenStrategy GenerateRootStrategy = generateStandardRootToken{}

	// GenerateDROperationTokenStrategy is the strategy used to generate a
	// DR operation token
	GenerateDROperationTokenStrategy GenerateRootStrategy = generateStandardRootToken{}
)

// GenerateRootStrategy allows us to swap out the strategy we want to use to
// create a token upon completion of the generate root process.
type GenerateRootStrategy interface {
	generate(context.Context, *Core) (string, func(), error)
	authenticate(context.Context, *Core, []byte) error
}

// generateStandardRootToken implements the GenerateRootStrategy and is in
// charge of creating standard root tokens.
type generateStandardRootToken struct{}

func (g generateStandardRootToken) authenticate(ctx context.Context, c *Core, combinedKey []byte) error {
	rootKey, err := c.unsealKeyToRootKeyPostUnseal(ctx, combinedKey)
	if err != nil {
		return fmt.Errorf("unable to authenticate: %w", err)
	}
	if err := c.barrier.VerifyRoot(rootKey); err != nil {
		return fmt.Errorf("root key verification failed: %w", err)
	}

	return nil
}

func (g generateStandardRootToken) generate(ctx context.Context, c *Core) (string, func(), error) {
	te, err := c.tokenStore.rootToken(ctx)
	if err != nil {
		c.logger.Error("root token generation failed", "error", err)
		return "", nil, err
	}
	if te == nil {
		c.logger.Error("got nil token entry back from root generation")
		return "", nil, fmt.Errorf("got nil token entry back from root generation")
	}

	cleanupFunc := func() {
		c.tokenStore.revokeOrphan(ctx, te.ID)
	}

	return te.ExternalID, cleanupFunc, nil
}

// GenerateRootConfig holds the configuration for a root generation
// command.
type GenerateRootConfig struct {
	Nonce          string
	PGPKey         string
	PGPFingerprint string
	OTP            string
	Strategy       GenerateRootStrategy
}

// GenerateRootResult holds the result of a root generation update
// command
type GenerateRootResult struct {
	Progress       int
	Required       int
	EncodedToken   string
	PGPFingerprint string
}

// GenerateRootProgress is used to return the root generation progress (num shares)
func (c *Core) GenerateRootProgress() (int, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.Sealed() && !c.recoveryMode {
		return 0, consts.ErrSealed
	}
	if c.standby && !c.recoveryMode {
		return 0, consts.ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	return len(c.generateRootProgress), nil
}

// GenerateRootConfiguration is used to read the root generation configuration
// It stubbornly refuses to return the OTP if one is there.
func (c *Core) GenerateRootConfiguration() (*GenerateRootConfig, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.Sealed() && !c.recoveryMode {
		return nil, consts.ErrSealed
	}
	if c.standby && !c.recoveryMode {
		return nil, consts.ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	// Copy the config if any
	var conf *GenerateRootConfig
	if c.generateRootConfig != nil {
		conf = new(GenerateRootConfig)
		*conf = *c.generateRootConfig
		conf.OTP = ""
		conf.Strategy = nil
	}
	return conf, nil
}

// GenerateRootInit is used to initialize the root generation settings
func (c *Core) GenerateRootInit(otp, pgpKey string, strategy GenerateRootStrategy) error {
	var fingerprint string
	switch {
	case len(otp) > 0:
		if (len(otp) != TokenLength+TokenPrefixLength && !c.DisableSSCTokens()) ||
			(len(otp) != TokenLength+OldTokenPrefixLength && c.DisableSSCTokens()) {
			return fmt.Errorf("OTP string is wrong length")
		}

	case len(pgpKey) > 0:
		fingerprints, err := pgpkeys.GetFingerprints([]string{pgpKey}, nil)
		if err != nil {
			return fmt.Errorf("error parsing PGP key: %w", err)
		}
		if len(fingerprints) != 1 || fingerprints[0] == "" {
			return fmt.Errorf("could not acquire PGP key entity")
		}
		fingerprint = fingerprints[0]

	default:
		return fmt.Errorf("otp or pgp_key parameter must be provided")
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.Sealed() && !c.recoveryMode {
		return consts.ErrSealed
	}
	barrierSealed, err := c.barrier.Sealed()
	if err != nil {
		return errors.New("unable to check barrier seal status")
	}
	if !barrierSealed && c.recoveryMode {
		return errors.New("attempt to generate recovery token when already unsealed")
	}
	if c.standby && !c.recoveryMode {
		return consts.ErrStandby
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
		Strategy:       strategy,
	}

	if c.logger.IsInfo() {
		switch strategy.(type) {
		case generateStandardRootToken:
			c.logger.Info("root generation initialized", "nonce", c.generateRootConfig.Nonce)
		case *generateRecoveryToken:
			c.logger.Info("recovery token generation initialized", "nonce", c.generateRootConfig.Nonce)
		default:
			c.logger.Info("dr operation token generation initialized", "nonce", c.generateRootConfig.Nonce)
		}
	}

	return nil
}

// GenerateRootUpdate is used to provide a new key part
func (c *Core) GenerateRootUpdate(ctx context.Context, key []byte, nonce string, strategy GenerateRootStrategy) (*GenerateRootResult, error) {
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
		config, err = c.seal.RecoveryConfig(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		config, err = c.seal.BarrierConfig(ctx)
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
	if c.Sealed() && !c.recoveryMode {
		return nil, consts.ErrSealed
	}

	barrierSealed, err := c.barrier.Sealed()
	if err != nil {
		return nil, errors.New("unable to check barrier seal status")
	}
	if !barrierSealed && c.recoveryMode {
		return nil, errors.New("attempt to generate recovery token when already unsealed")
	}

	if c.standby && !c.recoveryMode {
		return nil, consts.ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	// Ensure a generateRoot is in progress
	if c.generateRootConfig == nil {
		return nil, fmt.Errorf("no root generation in progress")
	}

	if nonce != c.generateRootConfig.Nonce {
		return nil, fmt.Errorf("incorrect nonce supplied; nonce for this root generation operation is %q", c.generateRootConfig.Nonce)
	}

	if strategy != c.generateRootConfig.Strategy {
		return nil, fmt.Errorf("incorrect strategy supplied; a generate root operation of another type is already in progress")
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
			c.logger.Debug("cannot generate root, not enough keys", "keys", progress, "threshold", config.SecretThreshold)
		}
		return &GenerateRootResult{
			Progress:       progress,
			Required:       config.SecretThreshold,
			PGPFingerprint: c.generateRootConfig.PGPFingerprint,
		}, nil
	}

	// Combine the key parts
	var combinedKey []byte
	if config.SecretThreshold == 1 {
		combinedKey = c.generateRootProgress[0]
		c.generateRootProgress = nil
	} else {
		combinedKey, err = shamir.Combine(c.generateRootProgress)
		c.generateRootProgress = nil
		if err != nil {
			return nil, fmt.Errorf("failed to compute root key: %w", err)
		}
	}

	if err := strategy.authenticate(ctx, c, combinedKey); err != nil {
		c.logger.Error("root generation aborted", "error", err.Error())
		return nil, fmt.Errorf("root generation aborted: %w", err)
	}

	// Run the generate strategy
	token, cleanupFunc, err := strategy.generate(ctx, c)
	if err != nil {
		return nil, err
	}

	var encodedToken string

	switch {
	case len(c.generateRootConfig.OTP) > 0:
		encodedToken, err = roottoken.EncodeToken(token, c.generateRootConfig.OTP)
	case len(c.generateRootConfig.PGPKey) > 0:
		var tokenBytesArr [][]byte
		_, tokenBytesArr, err = pgpkeys.EncryptShares([][]byte{[]byte(token)}, []string{c.generateRootConfig.PGPKey})
		encodedToken = base64.StdEncoding.EncodeToString(tokenBytesArr[0])
	default:
		err = fmt.Errorf("unreachable condition")
	}

	if err != nil {
		cleanupFunc()
		return nil, err
	}

	results := &GenerateRootResult{
		Progress:       progress,
		Required:       config.SecretThreshold,
		EncodedToken:   encodedToken,
		PGPFingerprint: c.generateRootConfig.PGPFingerprint,
	}

	switch strategy.(type) {
	case generateStandardRootToken:
		c.logger.Info("root generation finished", "nonce", c.generateRootConfig.Nonce)
	case *generateRecoveryToken:
		c.logger.Info("recovery token generation finished", "nonce", c.generateRootConfig.Nonce)
	default:
		c.logger.Info("dr operation token generation finished", "nonce", c.generateRootConfig.Nonce)
	}

	c.generateRootProgress = nil
	c.generateRootConfig = nil
	return results, nil
}

// GenerateRootCancel is used to cancel an in-progress root generation
func (c *Core) GenerateRootCancel() error {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.Sealed() && !c.recoveryMode {
		return consts.ErrSealed
	}
	if c.standby && !c.recoveryMode {
		return consts.ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	// Clear any progress or config
	c.generateRootConfig = nil
	c.generateRootProgress = nil
	return nil
}
