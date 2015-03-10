package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/hashicorp/vault/physical"
)

const (
	// coreSealConfigPath is the path used to store our seal configuration.
	// This value is stored in plaintext, since we must be able to read
	// it even with the Vault sealed. This is required so that we know
	// how many secret parts must be used to reconstruct the master key.
	coreSealConfigPath = "core/seal-config"
)

var (
	// ErrSealed is returned if an operation is performed on
	// a sealed barrier. No operation is expected to succeed before unsealing
	ErrSealed = errors.New("Vault is sealed")

	// ErrAlreadyInit is returned if the core is already
	// initialized. This prevents a re-initialization.
	ErrAlreadyInit = errors.New("Vault is already initialized")

	// ErrNotInit is returned if a non-initialized barrier
	// is attempted to be unsealed.
	ErrNotInit = errors.New("Vault is not initialized")
)

// SealConfig is used to describe the seal configuration
type SealConfig struct {
	// SecretParts is the number of parts the secret is
	// split into. This is the N value of Shamir
	SecretParts int `json:"secret_parts"`

	// SecretThreshold is the number of parts required
	// to open the vault. This is the T value of Shamir
	SecretThreshold int `json:"secret_threshold"`

	// SecretProgress is the number of parts already provided.
	// Once the SecretThreshold is reached, an unseal attempt
	// is made.
	SecretProgress int `json:"secret_progress"`
}

// Validate is used to sanity check the seal configuration
func (s *SealConfig) Validate() error {
	if s.SecretParts <= 0 {
		return fmt.Errorf("must have a positive number for secret parts")
	}
	if s.SecretThreshold <= 0 {
		return fmt.Errorf("must have a positive number for secret threshold")
	}
	if s.SecretThreshold > s.SecretParts {
		return fmt.Errorf("secret threshold cannot be larger than secret parts")
	}
	return nil
}

// InitResult is used to provide the key parts back after
// they are generated as part of the initialization.
type InitResult struct {
	SecretParts [][]byte
}

// Core is used as the central manager of Vault activity. It is the primary point of
// interface for API handlers and is responsible for managing the logical and physical
// backends, router, security barrier, and audit trails.
type Core struct {
	// physical backend is the un-trusted backend with durable data
	physical physical.Backend

	// barrier is the security barrier wrapping the physical backend
	barrier SecurityBarrier

	// router is responsible for managing the mount points for logical backends.
	router *Router

	// stateLock protects mutable state
	stateLock sync.RWMutex
	sealed    bool

	// unlockParts has the keys provided to Unseal until
	// the threshold number of parts is available.
	unlockParts [][]byte

	logger *log.Logger
}

// NewCore is used to construct a new core
func NewCore(physical physical.Backend) (*Core, error) {
	// Construct a new AES-GCM barrier
	barrier, err := NewAESGCMBarrier(physical)
	if err != nil {
		return nil, fmt.Errorf("barrier setup failed: %v", err)
	}

	// Setup the core
	c := &Core{
		physical: physical,
		barrier:  barrier,
		router:   NewRouter(),
		sealed:   true,
	}

	// Create and mount the system backend
	sys := &SystemBackend{
		core: c,
	}
	c.router.Mount(sys, "system", "sys/", nil)

	return c, nil
}

// HandleRequest is used to handle a new incoming request
func (c *Core) HandleRequest(req *Request) (*Response, error) {
	// TODO:
	return c.router.Route(req)
}

// Initialized checks if the Vault is already initialized
func (c *Core) Initialized() (bool, error) {
	// Check the barrier first
	init, err := c.barrier.Initialized()
	if err != nil || !init {
		c.logger.Printf("[ERR] core: barrier init check failed: %v", err)
		return false, err
	}
	if !init {
		c.logger.Printf("[INFO] core: security barrier not initialized")
		return false, nil
	}

	// Verify the seal configuration
	sealConf, err := c.SealConfig()
	if err != nil {
		return false, err
	}
	if sealConf == nil {
		return false, nil
	}
	return true, nil
}

// Initialize is used to initialize the Vault with the given
// configurations.
func (c *Core) Initialize(config *SealConfig) (*InitResult, error) {
	// Check if the seal configuraiton is valid
	if err := config.Validate(); err != nil {
		c.logger.Printf("[ERR] core: invalid seal configuration: %v", err)
		return nil, fmt.Errorf("invalid seal configuration: %v", err)
	}

	// TODO: Remove restrict
	if config.SecretParts != 1 {
		return nil, fmt.Errorf("Unsupported configuration")
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

	// Encode the seal configuration
	config.SecretProgress = 0
	buf, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to encode seal configuration: %v", err)
	}

	// Store the seal configuration
	pe := &physical.Entry{
		Key:   coreSealConfigPath,
		Value: buf,
	}
	if err := c.physical.Put(pe); err != nil {
		c.logger.Printf("[ERR] core: failed to read seal configuration: %v", err)
		return nil, fmt.Errorf("failed to check seal configuration: %v", err)
	}

	// Generate a master key
	masterKey, err := c.barrier.GenerateKey()
	if err != nil {
		c.logger.Printf("[ERR] core: failed to generate master key: %v", err)
		return nil, fmt.Errorf("master key generation failed: %v", err)
	}

	// Initialize the barrier
	if err := c.barrier.Initialize(masterKey); err != nil {
		c.logger.Printf("[ERR] core: failed to initialize barrier: %v", err)
		return nil, fmt.Errorf("failed to initialize barrier: %v", err)
	}

	// Return the master key if only a single key part is used
	results := new(InitResult)
	if config.SecretParts == 1 {
		results.SecretParts = append(results.SecretParts, masterKey)

	} else {
		// TODO: Support multiple parts
		panic("unsupported")
	}

	return results, nil
}

// Sealed checks if the Vault is current sealed
func (c *Core) Sealed() (bool, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	return c.sealed, nil
}

// SealConfiguration is used to return information
// about the configuration of the Vault and it's current
// status.
func (c *Core) SealConfig() (*SealConfig, error) {
	// Fetch the core configuration
	pe, err := c.physical.Get(coreSealConfigPath)
	if err != nil {
		c.logger.Printf("[ERR] core: failed to read seal configuration: %v", err)
		return nil, fmt.Errorf("failed to check seal configuration: %v", err)
	}

	// If the seal configuration is missing, we are not initialized
	if pe == nil {
		c.logger.Printf("[INFO] core: seal configuration missing, not initialized")
		return nil, nil
	}

	// Decode the barrier entry
	var conf SealConfig
	if err := json.Unmarshal(pe.Value, &conf); err != nil {
		c.logger.Printf("[ERR] core: failed to decode seal configuration: %v", err)
		return nil, fmt.Errorf("failed to decode seal configuration: %v", err)
	}

	// Check for a valid seal configuration
	if err := conf.Validate(); err != nil {
		c.logger.Printf("[ERR] core: invalid seal configuration: %v", err)
		return nil, fmt.Errorf("seal validation failed: %v", err)
	}
	return &conf, nil
}

// Unseal is used to provide one of the key parts to
// unseal the Vault.
func (c *Core) Unseal(key []byte) (bool, error) {
	c.stateLock.Lock()
	defer c.stateLock.Unlock()

	// TODO
	return false, nil
}

// Seal is used to re-seal the Vault. This requires the Vaultto
// be unsealed again to perform any further operations.
func (c *Core) Seal() error {
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	if c.sealed {
		return nil
	}
	c.logger.Printf("[INFO] core: vault is being sealed")
	c.sealed = true
	return c.barrier.Seal()
}
