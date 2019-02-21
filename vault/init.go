package vault

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/shamir"
)

// InitParams keeps the init function from being littered with too many
// params, that's it!
type InitParams struct {
	BarrierConfig   *SealConfig
	RecoveryConfig  *SealConfig
	RootTokenPGPKey string
}

// InitResult is used to provide the key parts back after
// they are generated as part of the initialization.
type InitResult struct {
	SecretShares   [][]byte
	RecoveryShares [][]byte
	RootToken      string
}

var (
	initPTFunc     = func(c *Core) func() { return nil }
	initInProgress uint32
)

// Initialized checks if the Vault is already initialized
func (c *Core) Initialized(ctx context.Context) (bool, error) {
	// Check the barrier first
	init, err := c.barrier.Initialized(ctx)
	if err != nil {
		c.logger.Error("barrier init check failed", "error", err)
		return false, err
	}
	if !init {
		c.logger.Info("security barrier not initialized")
		return false, nil
	}

	// Verify the seal configuration
	sealConf, err := c.seal.BarrierConfig(ctx)
	if err != nil {
		return false, err
	}
	if sealConf == nil {
		return false, fmt.Errorf("core: barrier reports initialized but no seal configuration found")
	}

	return true, nil
}

func (c *Core) generateShares(sc *SealConfig) ([]byte, [][]byte, error) {
	// Generate a master key
	masterKey, err := c.barrier.GenerateKey()
	if err != nil {
		return nil, nil, errwrap.Wrapf("key generation failed: {{err}}", err)
	}

	// Return the master key if only a single key part is used
	var unsealKeys [][]byte
	if sc.SecretShares == 1 {
		unsealKeys = append(unsealKeys, masterKey)
	} else {
		// Split the master key using the Shamir algorithm
		shares, err := shamir.Split(masterKey, sc.SecretShares, sc.SecretThreshold)
		if err != nil {
			return nil, nil, errwrap.Wrapf("failed to generate barrier shares: {{err}}", err)
		}
		unsealKeys = shares
	}

	// If we have PGP keys, perform the encryption
	if len(sc.PGPKeys) > 0 {
		hexEncodedShares := make([][]byte, len(unsealKeys))
		for i, _ := range unsealKeys {
			hexEncodedShares[i] = []byte(hex.EncodeToString(unsealKeys[i]))
		}
		_, encryptedShares, err := pgpkeys.EncryptShares(hexEncodedShares, sc.PGPKeys)
		if err != nil {
			return nil, nil, err
		}
		unsealKeys = encryptedShares
	}

	return masterKey, unsealKeys, nil
}

// Initialize is used to initialize the Vault with the given
// configurations.
func (c *Core) Initialize(ctx context.Context, initParams *InitParams) (*InitResult, error) {
	atomic.StoreUint32(&initInProgress, 1)
	defer atomic.StoreUint32(&initInProgress, 0)
	barrierConfig := initParams.BarrierConfig
	recoveryConfig := initParams.RecoveryConfig

	if c.seal.RecoveryKeySupported() {
		if recoveryConfig == nil {
			return nil, fmt.Errorf("recovery configuration must be supplied")
		}

		if recoveryConfig.SecretShares < 1 {
			return nil, fmt.Errorf("recovery configuration must specify a positive number of shares")
		}

		// Check if the seal configuration is valid
		if err := recoveryConfig.Validate(); err != nil {
			c.logger.Error("invalid recovery configuration", "error", err)
			return nil, errwrap.Wrapf("invalid recovery configuration: {{err}}", err)
		}
	}

	// Check if the seal configuration is valid
	if err := barrierConfig.Validate(); err != nil {
		c.logger.Error("invalid seal configuration", "error", err)
		return nil, errwrap.Wrapf("invalid seal configuration: {{err}}", err)
	}

	// Avoid an initialization race
	c.stateLock.Lock()
	defer c.stateLock.Unlock()

	// Check if we are initialized
	init, err := c.Initialized(ctx)
	if err != nil {
		return nil, err
	}
	if init {
		return nil, ErrAlreadyInit
	}

	err = c.seal.Init(ctx)
	if err != nil {
		c.logger.Error("failed to initialize seal", "error", err)
		return nil, errwrap.Wrapf("error initializing seal: {{err}}", err)
	}

	barrierKey, barrierUnsealKeys, err := c.generateShares(barrierConfig)
	if err != nil {
		c.logger.Error("error generating shares", "error", err)
		return nil, err
	}

	initPTCleanup := initPTFunc(c)
	if initPTCleanup != nil {
		defer initPTCleanup()
	}

	// Initialize the barrier
	if err := c.barrier.Initialize(ctx, barrierKey); err != nil {
		c.logger.Error("failed to initialize barrier", "error", err)
		return nil, errwrap.Wrapf("failed to initialize barrier: {{err}}", err)
	}
	if c.logger.IsInfo() {
		c.logger.Info("security barrier initialized", "shares", barrierConfig.SecretShares, "threshold", barrierConfig.SecretThreshold)
	}

	// Unseal the barrier
	if err := c.barrier.Unseal(ctx, barrierKey); err != nil {
		c.logger.Error("failed to unseal barrier", "error", err)
		return nil, errwrap.Wrapf("failed to unseal barrier: {{err}}", err)
	}

	// Ensure the barrier is re-sealed
	defer func() {
		// Defers are LIFO so we need to run this here too to ensure the stop
		// happens before sealing. preSeal also stops, so we just make the
		// stopping safe against multiple calls.
		if err := c.barrier.Seal(); err != nil {
			c.logger.Error("failed to seal barrier", "error", err)
		}
	}()

	err = c.seal.SetBarrierConfig(ctx, barrierConfig)
	if err != nil {
		c.logger.Error("failed to save barrier configuration", "error", err)
		return nil, errwrap.Wrapf("barrier configuration saving failed: {{err}}", err)
	}

	// If we are storing shares, pop them out of the returned results and push
	// them through the seal
	if barrierConfig.StoredShares > 0 {
		var keysToStore [][]byte
		for i := 0; i < barrierConfig.StoredShares; i++ {
			keysToStore = append(keysToStore, barrierUnsealKeys[0])
			barrierUnsealKeys = barrierUnsealKeys[1:]
		}
		if err := c.seal.SetStoredKeys(ctx, keysToStore); err != nil {
			c.logger.Error("failed to store keys", "error", err)
			return nil, errwrap.Wrapf("failed to store keys: {{err}}", err)
		}
	}

	results := &InitResult{
		SecretShares: barrierUnsealKeys,
	}

	// Perform initial setup
	if err := c.setupCluster(ctx); err != nil {
		c.logger.Error("cluster setup failed during init", "error", err)
		return nil, err
	}

	// Start tracking
	if initPTCleanup != nil {
		initPTCleanup()
	}

	activeCtx, ctxCancel := context.WithCancel(namespace.RootContext(nil))
	if err := c.postUnseal(activeCtx, ctxCancel, standardUnsealStrategy{}); err != nil {
		c.logger.Error("post-unseal setup failed during init", "error", err)
		return nil, err
	}

	// Save the configuration regardless, but only generate a key if it's not
	// disabled. When using recovery keys they are stored in the barrier, so
	// this must happen post-unseal.
	if c.seal.RecoveryKeySupported() {
		err = c.seal.SetRecoveryConfig(ctx, recoveryConfig)
		if err != nil {
			c.logger.Error("failed to save recovery configuration", "error", err)
			return nil, errwrap.Wrapf("recovery configuration saving failed: {{err}}", err)
		}

		if recoveryConfig.SecretShares > 0 {
			recoveryKey, recoveryUnsealKeys, err := c.generateShares(recoveryConfig)
			if err != nil {
				c.logger.Error("failed to generate recovery shares", "error", err)
				return nil, err
			}

			err = c.seal.SetRecoveryKey(ctx, recoveryKey)
			if err != nil {
				return nil, err
			}

			results.RecoveryShares = recoveryUnsealKeys
		}
	}

	// Generate a new root token
	rootToken, err := c.tokenStore.rootToken(ctx)
	if err != nil {
		c.logger.Error("root token generation failed", "error", err)
		return nil, err
	}
	results.RootToken = rootToken.ID
	c.logger.Info("root token generated")

	if initParams.RootTokenPGPKey != "" {
		_, encryptedVals, err := pgpkeys.EncryptShares([][]byte{[]byte(results.RootToken)}, []string{initParams.RootTokenPGPKey})
		if err != nil {
			c.logger.Error("root token encryption failed", "error", err)
			return nil, err
		}
		results.RootToken = base64.StdEncoding.EncodeToString(encryptedVals[0])
	}

	// Prepare to re-seal
	if err := c.preSeal(); err != nil {
		c.logger.Error("pre-seal teardown failed", "error", err)
		return nil, err
	}

	return results, nil
}

// UnsealWithStoredKeys performs auto-unseal using stored keys. An error
// return value of "nil" implies the Vault instance is unsealed.
//
// Callers should attempt to retry any NOnFatalErrors. Callers should
// not re-attempt fatal errors.
func (c *Core) UnsealWithStoredKeys(ctx context.Context) error {
	c.unsealWithStoredKeysLock.Lock()
	defer c.unsealWithStoredKeysLock.Unlock()

	if !c.seal.StoredKeysSupported() {
		return nil
	}

	// Disallow auto-unsealing when migrating
	if c.IsInSealMigration() {
		return NewNonFatalError(errors.New("cannot auto-unseal during seal migration"))
	}

	sealed := c.Sealed()
	if !sealed {
		c.Logger().Warn("attempted unseal with stored keys, but vault is already unsealed")
		return nil
	}

	c.Logger().Info("stored unseal keys supported, attempting fetch")
	keys, err := c.seal.GetStoredKeys(ctx)
	if err != nil {
		return NewNonFatalError(errwrap.Wrapf("fetching stored unseal keys failed: {{err}}", err))
	}

	// This usually happens when auto-unseal is configured, but the servers have
	// not been initialized yet.
	if len(keys) == 0 {
		return NewNonFatalError(errors.New("stored unseal keys are supported, but none were found"))
	}

	unsealed := false
	keysUsed := 0
	for _, key := range keys {
		unsealed, err = c.Unseal(key)
		if err != nil {
			return NewNonFatalError(errwrap.Wrapf("unseal with stored key failed: {{err}}", err))
		}
		keysUsed++
		if unsealed {
			break
		}
	}

	if !unsealed {
		// This most likely means that the user configured Vault to only store a
		// subset of the required threshold of keys. We still consider this a
		// "success", since trying again would yield the same result.
		c.Logger().Warn("vault still sealed after using stored unseal keys", "stored_keys_used", keysUsed)
	} else {
		c.Logger().Info("unsealed with stored keys", "stored_keys_used", keysUsed)
	}

	return nil
}
