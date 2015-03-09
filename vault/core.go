package vault

import (
	"fmt"

	"github.com/hashicorp/vault/physical"
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
	return c.router.Route(req)
}

// Initialized checks if the Vault is already initialized
func (c *Core) Initialized() (bool, error) {
	return false, nil
}

// Initialize is used to initialize the Vault with the given
// configurations.
func (c *Core) Initialize(config *SealConfig) error {
	return nil
}

// Sealed checks if the Vault is current sealed
func (c *Core) Sealed() (bool, error) {
	return true, nil
}

// SealConfiguration is used to return information
// about the configuration of the Vault and it's current
// status.
func (c *Core) SealConfig() (*SealConfig, error) {
	return nil, nil
}

// Unseal is used to provide one of the key parts to
// unseal the Vault.
func (c *Core) Unseal(key []byte) (bool, error) {
	return false, nil
}
