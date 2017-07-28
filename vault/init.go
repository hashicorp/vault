package vault

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/shamir"
)

// keySharesMetadataStorageValue holds metadata about all the shamir shares of
// a key. This informaion is stored during the initialization and updated
// during rekey. It is referred during unsealing, during rekeying and during
// generation of a root token. Essentially in all the places where the key
// shares are put to use.
type keySharesMetadataStorageValue struct {
	// Data is a map from each of the key share to its respective identifier
	Data map[string]*KeyShareMetadata `json:"data" structs:"data" mapstructure:"data"`
}

// KeyShareMetadata holds metadata associated with each key share
type KeyShareMetadata struct {
	// Name is a human readable name optionally provided by the caller to be
	// associated with the identifier of the key share.
	Name string `json:"name" structs:"name" mapstructure:"name"`

	// ID is the UUID associated with the unseal key share
	ID string `json:"id" structs:"id" mapstructure:"id"`
}

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
	SecretShares           [][]byte
	RecoveryShares         [][]byte
	RootToken              string
	SecretSharesMetadata   []*KeyShareMetadata
	RecoverySharesMetadata []*KeyShareMetadata
}

// InitKeySharesIdentifiersResponse contains the UUID identifiers associated with the
// key shares.
type InitKeySharesIdentifiersResponse struct {
	KeyIdentifiers []*KeyShareMetadata
}

// GenerateSharesResult is used to hold the key and its shamir key shares. If
// PGP keys are used to encrypt the key shares this will also hold the
// encrypted key shares and the PGP key fingerprint of the respective key that
// encrypted each share.
type GenerateSharesResult struct {
	Key                   []byte
	KeyShares             [][]byte
	PGPKeyFingerprints    []string
	PGPEncryptedKeyShares [][]byte
}

// KeySharesIdentifiers returns the unique UUIDs associated with each of the
// key shares. This enables key share holders to verify that the identifier of
// the unseal/recovery key share given to them is actually valid.
func (c *Core) KeySharesIdentifiers(recovery bool) (*InitKeySharesIdentifiersResponse, error) {
	if c.sealed {
		return nil, consts.ErrSealed
	}
	if c.standby {
		return nil, consts.ErrStandby
	}

	var path string
	if recovery {
		path = coreRecoverySharesMetadataPath
	} else {
		path = coreSecretSharesMetadataPath
	}

	entry, err := c.barrier.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch key shares metadata: %v", err)
	}
	if entry == nil {
		return nil, nil
	}

	var keySharesMetadataValue keySharesMetadataStorageValue
	if err = jsonutil.DecodeJSON(entry.Value, &keySharesMetadataValue); err != nil {
		return nil, fmt.Errorf("failed to decode unseal metadata entry: %v", err)
	}

	response := &InitKeySharesIdentifiersResponse{}
	for _, keyShareMetadata := range keySharesMetadataValue.Data {
		if keyShareMetadata == nil || keyShareMetadata.ID == "" {
			return nil, fmt.Errorf("invalid unseal metadata entry in storage")
		}
		response.KeyIdentifiers = append(response.KeyIdentifiers, keyShareMetadata)
	}

	return response, nil
}

// Initialized checks if the Vault is already initialized
func (c *Core) Initialized() (bool, error) {
	// Check the barrier first
	init, err := c.barrier.Initialized()
	if err != nil {
		c.logger.Error("core: barrier init check failed", "error", err)
		return false, err
	}
	if !init {
		c.logger.Info("core: security barrier not initialized")
		return false, nil
	}

	// Verify the seal configuration
	sealConf, err := c.seal.BarrierConfig()
	if err != nil {
		return false, err
	}
	if sealConf == nil {
		return false, fmt.Errorf("core: barrier reports initialized but no seal configuration found")
	}

	return true, nil
}

// generateShares takes in a seal configuration and creates a barrier key. The
// key will then be split into the specified number of shares. If PGP keys are
// supplied, each key share will be encrypted with respective PGP keys.
func (c *Core) generateShares(sc *SealConfig) (*GenerateSharesResult, error) {
	// Generate a barrier key
	keyBytes, err := c.barrier.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("key generation failed: %v", err)
	}

	// Return the barrier key if only a single key part is used
	var keyShares [][]byte
	if sc.SecretShares == 1 {
		keyShares = append(keyShares, keyBytes)
	} else {
		// Split the master key using the Shamir algorithm
		keyShares, err = shamir.Split(keyBytes, sc.SecretShares, sc.SecretThreshold)
		if err != nil {
			return nil, fmt.Errorf("failed to generate barrier shares: %v", err)
		}
	}

	// If PGP keys are supplied, encrypt the key shares with respective PGP key
	var pgpEncryptedKeyShares [][]byte
	var pgpKeyFingerprints []string
	if len(sc.PGPKeys) > 0 {
		hexEncodedShares := make([][]byte, len(keyShares))
		for i, _ := range keyShares {
			hexEncodedShares[i] = []byte(hex.EncodeToString(keyShares[i]))
		}
		pgpKeyFingerprints, pgpEncryptedKeyShares, err = pgpkeys.EncryptShares(hexEncodedShares, sc.PGPKeys)
		if err != nil {
			return nil, err
		}
	}

	return &GenerateSharesResult{
		Key:                   keyBytes,
		KeyShares:             keyShares,
		PGPKeyFingerprints:    pgpKeyFingerprints,
		PGPEncryptedKeyShares: pgpEncryptedKeyShares,
	}, nil
}

// prepareKeySharesMetadata takes in the unseal key shares, both
// encrypted and unencrypted, associates identifiers for each key share and
// JSON encodes it. Identifier for unencrypted key shares will be UUIDs.
func (c *Core) prepareKeySharesMetadata(keyShares [][]byte, keyIdentifierNames string) ([]byte, []*KeyShareMetadata, error) {
	// If keyIdentifierNames are supplied, parse them
	var identifierNames []string
	if keyIdentifierNames != "" {
		identifierNames = strutil.ParseDedupAndSortStrings(keyIdentifierNames, ",")

		if len(identifierNames) != len(keyShares) {
			c.logger.Error("core: number of key identifier names not matching the number of key shares")
			return nil, nil, fmt.Errorf("number of key identifier names not matching the number of key shares")
		}
	}

	keySharesMetadataValue := &keySharesMetadataStorageValue{
		Data: make(map[string]*KeyShareMetadata),
	}

	var keySharesMetadata []*KeyShareMetadata

	// Associate a UUID for each key share
	for i, keyShare := range keyShares {
		metadata := &KeyShareMetadata{}
		keyUUID, err := uuid.GenerateUUID()
		if err != nil {
			c.logger.Error("core: failed to generate key share identifier", "error", err)
			return nil, nil, fmt.Errorf("failed to generate key share identifier: %v", err)
		}
		metadata.ID = keyUUID

		// Attach the name for the identifier if supplied
		if len(identifierNames) > 0 {
			metadata.Name = identifierNames[i]
		}

		keySharesMetadataValue.Data[base64.StdEncoding.EncodeToString(salt.SHA256Hash(keyShare))] = metadata
		keySharesMetadata = append(keySharesMetadata, metadata)
	}

	// JSON encode the key shares matadata
	keySharesMetadataJSON, err := jsonutil.EncodeJSON(keySharesMetadataValue)
	if err != nil {
		c.logger.Error("core: failed to encode key shares metadata", "error", err)
		return nil, nil, err
	}

	return keySharesMetadataJSON, keySharesMetadata, nil
}

// Initialize is used to initialize the Vault with the given
// configurations.
func (c *Core) Initialize(initParams *InitParams) (*InitResult, error) {
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
			c.logger.Error("core: invalid recovery configuration", "error", err)
			return nil, fmt.Errorf("invalid recovery configuration: %v", err)
		}
	}

	// Check if the seal configuration is valid
	if err := barrierConfig.Validate(); err != nil {
		c.logger.Error("core: invalid seal configuration", "error", err)
		return nil, fmt.Errorf("invalid seal configuration: %v", err)
	}

	// Avoid an initialization race
	c.stateLock.Lock()
	defer c.stateLock.Unlock()

	// Check if we are initialized
	init, err := c.Initialized()
	if err != nil {
		return nil, err
	}
	if init {
		return nil, ErrAlreadyInit
	}

	err = c.seal.Init()
	if err != nil {
		c.logger.Error("core: failed to initialize seal", "error", err)
		return nil, fmt.Errorf("error initializing seal: %v", err)
	}

	err = c.seal.SetBarrierConfig(barrierConfig)
	if err != nil {
		c.logger.Error("core: failed to save barrier configuration", "error", err)
		return nil, fmt.Errorf("barrier configuration saving failed: %v", err)
	}

	barrierShares, err := c.generateShares(barrierConfig)
	if err != nil || barrierShares == nil {
		c.logger.Error("core: error generating barrier shares", "error", err)
		return nil, err
	}

	results := &InitResult{}

	var secretSharesMetadataJSON []byte

	//
	// Prepare metadata for each of the unseal key share generated. Associate
	// the metatada with plaintext unseal key shares and not the PGP encrypted
	// key shares. Metadata should be created for all the key shares and hence
	// this should be done before processing stored keys.
	//

	// Associate metadata for all the unseal key shares
	secretSharesMetadataJSON, results.SecretSharesMetadata, err = c.prepareKeySharesMetadata(barrierShares.KeyShares, barrierConfig.SecretSharesIdentifierNames)
	if err != nil {
		c.logger.Error("core: failed to prepare unseal key shares metadata", "error", err)
		return nil, fmt.Errorf("failed to prepare unseal key shares metadata: %v", err)
	}

	// Determine whether to return plaintext unseal key shares or its PGP
	// encrypted counterparts
	var returnedKeys [][]byte
	switch {
	case len(barrierShares.PGPEncryptedKeyShares) > 0:
		returnedKeys = barrierShares.PGPEncryptedKeyShares
	default:
		returnedKeys = barrierShares.KeyShares
	}

	// If we are storing shares, pop them out of the returned results and push
	// them through the seal
	if barrierConfig.StoredShares > 0 {
		if len(barrierConfig.PGPKeys) > 0 {
			c.logger.Error("core: PGP keys not supported when storing shares")
			return nil, fmt.Errorf("PGP keys not supported when storing shares")
		}

		// Note that returnedKeys will always be unencrypted here
		var keysToStore [][]byte
		for i := 0; i < barrierConfig.StoredShares; i++ {
			keysToStore = append(keysToStore, returnedKeys[0])
			returnedKeys = returnedKeys[1:]
		}
		if err := c.seal.SetStoredKeys(keysToStore); err != nil {
			c.logger.Error("core: failed to store keys", "error", err)
			return nil, fmt.Errorf("failed to store keys: %v", err)
		}
	}

	results.SecretShares = returnedKeys

	// Initialize the barrier
	if err := c.barrier.Initialize(barrierShares.Key); err != nil {
		c.logger.Error("core: failed to initialize barrier", "error", err)
		return nil, fmt.Errorf("failed to initialize barrier: %v", err)
	}
	if c.logger.IsInfo() {
		c.logger.Info("core: security barrier initialized", "shares", barrierConfig.SecretShares, "threshold", barrierConfig.SecretThreshold)
	}

	// Unseal the barrier
	if err := c.barrier.Unseal(barrierShares.Key); err != nil {
		c.logger.Error("core: failed to unseal barrier", "error", err)
		return nil, fmt.Errorf("failed to unseal barrier: %v", err)
	}

	// Ensure the barrier is re-sealed
	defer func() {
		// Defers are LIFO so we need to run this here too to ensure the stop
		// happens before sealing. preSeal also stops, so we just make the
		// stopping safe against multiple calls.
		if err := c.barrier.Seal(); err != nil {
			c.logger.Error("core: failed to seal barrier", "error", err)
		}
	}()

	// Now that the barrier is unsealed, persist the unseal shares metadata
	err = c.barrier.Put(&Entry{
		Key:   coreSecretSharesMetadataPath,
		Value: secretSharesMetadataJSON,
	})
	if err != nil {
		c.logger.Error("core: failed to store unseal metadata", "error", err)
		return nil, err
	}

	// Perform initial setup
	if err := c.setupCluster(); err != nil {
		c.logger.Error("core: cluster setup failed during init", "error", err)
		return nil, err
	}
	if err := c.postUnseal(); err != nil {
		c.logger.Error("core: post-unseal setup failed during init", "error", err)
		return nil, err
	}

	// Save the configuration regardless, but only generate a key if it's not
	// disabled. When using recovery keys they are stored in the barrier, so
	// this must happen post-unseal.
	if c.seal.RecoveryKeySupported() {
		err = c.seal.SetRecoveryConfig(recoveryConfig)
		if err != nil {
			c.logger.Error("core: failed to save recovery configuration", "error", err)
			return nil, fmt.Errorf("recovery configuration saving failed: %v", err)
		}

		if recoveryConfig.SecretShares > 0 {
			recoveryShares, err := c.generateShares(recoveryConfig)
			if err != nil || recoveryShares == nil {
				c.logger.Error("core: failed to generate recovery shares", "error", err)
				return nil, err
			}

			err = c.seal.SetRecoveryKey(recoveryShares.Key)
			if err != nil {
				return nil, err
			}

			// Associate metadata for all the recovery key shares
			var recoverySharesMetadataJSON []byte
			recoverySharesMetadataJSON, results.RecoverySharesMetadata, err = c.prepareKeySharesMetadata(recoveryShares.KeyShares, recoveryConfig.SecretSharesIdentifierNames)
			if err != nil {
				c.logger.Error("core: failed to prepare recovery key shares metadata", "error", err)
				return nil, fmt.Errorf("failed to prepare recovery key shares metadata: %v", err)
			}

			err = c.barrier.Put(&Entry{
				Key:   coreRecoverySharesMetadataPath,
				Value: recoverySharesMetadataJSON,
			})
			if err != nil {
				c.logger.Error("core: failed to store recovery shares metadata", "error", err)
				return nil, err
			}

			switch {
			case len(recoveryShares.PGPEncryptedKeyShares) > 0:
				results.RecoveryShares = recoveryShares.PGPEncryptedKeyShares
			default:
				results.RecoveryShares = recoveryShares.KeyShares
			}

		}
	}

	// Generate a new root token
	rootToken, err := c.tokenStore.rootToken()
	if err != nil {
		c.logger.Error("core: root token generation failed", "error", err)
		return nil, err
	}
	results.RootToken = rootToken.ID
	c.logger.Info("core: root token generated")

	if initParams.RootTokenPGPKey != "" {
		_, encryptedVals, err := pgpkeys.EncryptShares([][]byte{[]byte(results.RootToken)}, []string{initParams.RootTokenPGPKey})
		if err != nil {
			c.logger.Error("core: root token encryption failed", "error", err)
			return nil, err
		}
		results.RootToken = base64.StdEncoding.EncodeToString(encryptedVals[0])
	}

	// Prepare to re-seal
	if err := c.preSeal(); err != nil {
		c.logger.Error("core: pre-seal teardown failed", "error", err)
		return nil, err
	}

	return results, nil
}

func (c *Core) UnsealWithStoredKeys() error {
	if !c.seal.StoredKeysSupported() {
		return nil
	}

	sealed, err := c.Sealed()
	if err != nil {
		c.logger.Error("core: error checking sealed status in auto-unseal", "error", err)
		return fmt.Errorf("error checking sealed status in auto-unseal: %s", err)
	}
	if !sealed {
		return nil
	}

	c.logger.Info("core: stored unseal keys supported, attempting fetch")
	keys, err := c.seal.GetStoredKeys()
	if err != nil {
		c.logger.Error("core: fetching stored unseal keys failed", "error", err)
		return &NonFatalError{Err: fmt.Errorf("fetching stored unseal keys failed: %v", err)}
	}
	if len(keys) == 0 {
		c.logger.Warn("core: stored unseal key(s) supported but none found")
	} else {
		unsealed := false
		keysUsed := 0
		for _, key := range keys {
			unsealed, err = c.Unseal(key)
			if err != nil {
				c.logger.Error("core: unseal with stored unseal key failed", "error", err)
				return &NonFatalError{Err: fmt.Errorf("unseal with stored key failed: %v", err)}
			}
			keysUsed += 1
			if unsealed {
				break
			}
		}
		if !unsealed {
			if c.logger.IsWarn() {
				c.logger.Warn("core: stored unseal key(s) used but Vault not unsealed yet", "stored_keys_used", keysUsed)
			}
		} else {
			if c.logger.IsInfo() {
				c.logger.Info("core: successfully unsealed with stored key(s)", "stored_keys_used", keysUsed)
			}
		}
	}

	return nil
}
