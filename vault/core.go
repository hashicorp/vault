package vault

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/shamir"
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

	// ErrInternalError is returned when we don't want to leak
	// any information about an internal error
	ErrInternalError = errors.New("internal error")
)

// SealConfig is used to describe the seal configuration
type SealConfig struct {
	// SecretShares is the number of shares the secret is
	// split into. This is the N value of Shamir
	SecretShares int `json:"secret_shares"`

	// SecretThreshold is the number of parts required
	// to open the vault. This is the T value of Shamir
	SecretThreshold int `json:"secret_threshold"`
}

// Validate is used to sanity check the seal configuration
func (s *SealConfig) Validate() error {
	if s.SecretShares < 1 {
		return fmt.Errorf("secret shares must be at least one")
	}
	if s.SecretThreshold < 1 {
		return fmt.Errorf("secret threshold must be at least one")
	}
	if s.SecretShares > 255 {
		return fmt.Errorf("secret shares must be less than 256")
	}
	if s.SecretThreshold > 255 {
		return fmt.Errorf("secret threshold must be less than 256")
	}
	if s.SecretThreshold > s.SecretShares {
		return fmt.Errorf("secret threshold cannot be larger than secret shares")
	}
	return nil
}

// InitResult is used to provide the key parts back after
// they are generated as part of the initialization.
type InitResult struct {
	SecretShares [][]byte
}

// ErrInvalidKey is returned if there is an error with a
// provided unseal key.
type ErrInvalidKey struct {
	Reason string
}

func (e *ErrInvalidKey) Error() string {
	return fmt.Sprintf("invalid key: %v", e.Reason)
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

	// backends is the mapping of backends to use for this core
	backends map[string]logical.Factory

	// stateLock protects mutable state
	stateLock sync.RWMutex
	sealed    bool

	// unlockParts has the keys provided to Unseal until
	// the threshold number of parts is available.
	unlockParts [][]byte

	// mounts is loaded after unseal since it is a protected
	// configuration
	mounts     *MountTable
	mountsLock sync.RWMutex

	// systemView is the barrier view for the system backend
	systemView *BarrierView

	// expiration manager is used for managing vaultIDs,
	// renewal, expiration and revocation
	expiration *ExpirationManager

	logger *log.Logger
}

// CoreConfig is used to parameterize a core
type CoreConfig struct {
	Backends map[string]logical.Factory
	Physical physical.Backend
	Logger   *log.Logger
}

// NewCore isk used to construct a new core
func NewCore(conf *CoreConfig) (*Core, error) {
	// Construct a new AES-GCM barrier
	barrier, err := NewAESGCMBarrier(conf.Physical)
	if err != nil {
		return nil, fmt.Errorf("barrier setup failed: %v", err)
	}

	// Make a default logger if not provided
	if conf.Logger == nil {
		conf.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}

	// Setup the core
	c := &Core{
		physical: conf.Physical,
		barrier:  barrier,
		router:   NewRouter(),
		sealed:   true,
		logger:   conf.Logger,
	}

	// Setup the backends
	backends := make(map[string]logical.Factory)
	for k, f := range conf.Backends {
		backends[k] = f
	}
	backends["generic"] = PassthroughBackendFactory
	backends["system"] = func(map[string]string) (logical.Backend, error) {
		return NewSystemBackend(c), nil
	}

	c.backends = backends
	return c, nil
}

// HandleRequest is used to handle a new incoming request
func (c *Core) HandleRequest(req *logical.Request) (*logical.Response, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, ErrSealed
	}

	// TODO: Enforce ACLs

	// Route the request
	resp, err := c.router.Route(req)

	// Check if there is a lease, we must register this
	if resp != nil && resp.IsSecret && resp.Lease != nil {
		vaultID, err := c.expiration.Register(req, resp)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to register lease (request: %#v, response: %#v): %v", req, resp, err)
			return nil, ErrInternalError
		}
		resp.Lease.VaultID = vaultID
	}

	// Return the response and error
	return resp, err
}

// Initialized checks if the Vault is already initialized
func (c *Core) Initialized() (bool, error) {
	// Check the barrier first
	init, err := c.barrier.Initialized()
	if err != nil {
		c.logger.Printf("[ERR] core: barrier init check failed: %v", err)
		return false, err
	}
	if !init {
		return false, nil
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
	if config.SecretShares == 1 {
		results.SecretShares = append(results.SecretShares, masterKey)

	} else {
		// Split the master key using the Shamir algorithm
		shares, err := shamir.Split(masterKey, config.SecretShares, config.SecretThreshold)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to generate shares: %v", err)
			return nil, fmt.Errorf("failed to generate shares: %v", err)
		}
		results.SecretShares = shares
	}

	c.logger.Printf("[INFO] core: security barrier initialized")
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

// SecretProgress returns the number of keys provided so far
func (c *Core) SecretProgress() int {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	return len(c.unlockParts)
}

// Unseal is used to provide one of the key parts to unseal the Vault.
//
// They key given as a parameter will automatically be zerod after
// this method is done with it. If you want to keep the key around, a copy
// should be made.
func (c *Core) Unseal(key []byte) (bool, error) {
	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		return false, &ErrInvalidKey{fmt.Sprintf("key is shorter than minimum %d bytes", min)}
	}
	if len(key) > max {
		return false, &ErrInvalidKey{fmt.Sprintf("key is longer than maximum %d bytes", max)}
	}

	// Get the seal configuration
	config, err := c.SealConfig()
	if err != nil {
		return false, err
	}

	// Ensure the barrier is initialized
	if config == nil {
		return false, ErrNotInit
	}

	c.stateLock.Lock()
	defer c.stateLock.Unlock()

	// Check if already unsealed
	if !c.sealed {
		return true, nil
	}

	// Check if we already have this piece
	for _, existing := range c.unlockParts {
		if bytes.Equal(existing, key) {
			return false, nil
		}
	}

	// Store this key
	c.unlockParts = append(c.unlockParts, key)

	// Check if we don't have enough keys to unlock
	if len(c.unlockParts) < config.SecretThreshold {
		c.logger.Printf("[DEBUG] core: cannot unseal, have %d of %d keys",
			len(c.unlockParts), config.SecretThreshold)
		return false, nil
	}

	// Recover the master key
	var masterKey []byte
	if config.SecretThreshold == 1 {
		masterKey = c.unlockParts[0]
		c.unlockParts = nil
	} else {
		masterKey, err = shamir.Combine(c.unlockParts)
		c.unlockParts = nil
		if err != nil {
			return false, fmt.Errorf("failed to compute master key: %v", err)
		}
	}
	defer memzero(masterKey)

	// Attempt to unlock
	if err := c.barrier.Unseal(masterKey); err != nil {
		return false, err
	}
	c.logger.Printf("[INFO] core: vault is unsealed")

	// Do post-unseal setup
	c.logger.Printf("[INFO] core: post-unseal setup starting")
	if err := c.postUnseal(); err != nil {
		c.logger.Printf("[ERR] core: post-unseal setup failed: %v", err)
		c.barrier.Seal()
		c.logger.Printf("[WARN] core: vault is sealed")
		return false, err
	}
	c.logger.Printf("[INFO] core: post-unseal setup complete")

	// Success!
	c.sealed = false
	return true, nil
}

// Seal is used to re-seal the Vault. This requires the Vault to
// be unsealed again to perform any further operations.
func (c *Core) Seal() error {
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	if c.sealed {
		return nil
	}
	c.sealed = true

	// Do pre-seal teardown
	c.logger.Printf("[INFO] core: pre-seal teardown starting")
	if err := c.preSeal(); err != nil {
		c.logger.Printf("[ERR] core: pre-seal teardown failed: %v", err)
		return fmt.Errorf("internal error")
	}
	c.logger.Printf("[INFO] core: pre-seal teardown complete")

	if err := c.barrier.Seal(); err != nil {
		return err
	}
	c.logger.Printf("[INFO] core: vault is sealed")
	return nil
}

// postUnseal is invoked after the barrier is unsealed, but before
// allowing any user operations. This allows us to setup any state that
// requires the Vault to be unsealed such as mount tables, logical backends,
// credential stores, etc.
func (c *Core) postUnseal() error {
	if err := c.loadMounts(); err != nil {
		return err
	}
	if err := c.setupMounts(); err != nil {
		return err
	}
	if err := c.setupExpiration(); err != nil {
		return err
	}
	return nil
}

// preSeal is invoked before the barrier is sealed, allowing
// for any state teardown required.
func (c *Core) preSeal() error {
	if err := c.stopExpiration(); err != nil {
		return err
	}
	if err := c.unloadMounts(); err != nil {
		return err
	}
	return nil
}
