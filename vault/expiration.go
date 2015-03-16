package vault

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/hashicorp/vault/logical"
)

const (
	// expirationSubPath is the sub-path used for the expiration manager
	// view. This is nested under the system view.
	expirationSubPath = "expire/"
)

// ExpirationManager is used by the Core to manage leases. Secrets
// can provide a lease, meaning that they can be renewed or revoked.
// If a secret is not renewed in timely manner, it may be expired, and
// the ExpirationManager will handle doing automatic revocation.
type ExpirationManager struct {
	router *Router
	view   *BarrierView
	logger *log.Logger

	pending     map[string]*time.Timer
	pendingLock sync.RWMutex
}

// NewExpirationManager creates a new ExpirationManager that is backed
// using a given view, and uses the provided router for revocation.
func NewExpirationManager(router *Router, view *BarrierView, logger *log.Logger) *ExpirationManager {
	if logger == nil {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	exp := &ExpirationManager{
		router:  router,
		view:    view,
		logger:  logger,
		pending: make(map[string]*time.Timer),
	}
	return exp
}

// setupExpiration is invoked after we've loaded the mount table to
// initialize the expiration manager
func (c *Core) setupExpiration() error {
	// Create a sub-view
	view := c.systemView.SubView(expirationSubPath)

	// Create the manager
	mgr := NewExpirationManager(c.router, view, c.logger)
	c.expiration = mgr

	// Restore the existing state
	if err := c.expiration.Restore(); err != nil {
		return fmt.Errorf("expiration state restore failed: %v", err)
	}
	return nil
}

// stopExpiration is used to stop the expiration manager before
// sealing the Vault.
func (c *Core) stopExpiration() error {
	if err := c.expiration.Stop(); err != nil {
		return err
	}
	c.expiration = nil
	return nil
}

// Restore is used to recover the lease states when starting.
// This is used after starting the vault.
func (m *ExpirationManager) Restore() error {
	// TODO: Restore...
	return nil
}

// Stop is used to prevent further automatic revocations.
// This must be called before sealing the view.
func (m *ExpirationManager) Stop() error {
	// Stop all the pending expiration timers
	m.pendingLock.Lock()
	for _, timer := range m.pending {
		timer.Stop()
	}
	m.pending = make(map[string]*time.Timer)
	m.pendingLock.Unlock()
	return nil
}

// Revoke is used to revoke a secret named by the given vaultID
func (m *ExpirationManager) Revoke(vaultID string) error {
	return nil
}

// RevokePrefix is used to revoke all secrets with a given prefix.
// The prefix maps to that of the mount table to make this simpler
// to reason about.
func (m *ExpirationManager) RevokePrefix(prefix string) error {
	return nil
}

// Renew is used to renew a secret using the given vaultID
// and a renew interval. The increment may be ignored.
func (m *ExpirationManager) Renew(vaultID string, increment time.Duration) (*logical.Lease, error) {
	return nil, nil
}

// Register is used to take a request and response with an associated
// lease. The secret gets assigned a vaultId and the management of
// of lease is assumed by the expiration manager.
func (m *ExpirationManager) Register(req *logical.Request, resp *logical.Response) (string, error) {
	// Ignore if there is no lease
	if resp == nil || resp.Lease == nil {
		return "", nil
	}

	// Validate the lease
	if err := resp.Lease.Validate(); err != nil {
		return "", err
	}

	// Cannot register a non-secret (e.g. a policy or configuration key)
	if !resp.IsSecret {
		return "", fmt.Errorf("cannot attach lease to non-secret")
	}

	// Create a lease entry
	now := time.Now().UTC()
	le := leaseEntry{
		VaultID:   path.Join(req.Path, generateUUID()),
		Path:      req.Path,
		Data:      resp.Data,
		Lease:     resp.Lease,
		IssueTime: now,
		RenewTime: now,
	}

	// Encode the entry
	if err := m.persistEntry(&le); err != nil {
		return "", err
	}

	// Setup revocation timer
	m.pendingLock.Lock()
	timer := time.AfterFunc(resp.Lease.Duration, func() {
		m.expireID(le.VaultID)
	})
	m.pending[le.VaultID] = timer
	m.pendingLock.Unlock()

	// Done
	return le.VaultID, nil
}

// expireID is invoked when a given ID is expired
func (m *ExpirationManager) expireID(vaultID string) {
	// Clear from the pending expiration
	m.pendingLock.Lock()
	delete(m.pending, vaultID)
	m.pendingLock.Unlock()

	// Load the entry
	le, err := m.loadEntry(vaultID)
	if err != nil {
		m.logger.Printf("[ERR] expire: failed to read entry '%s': %v", vaultID, err)
	}

	// Revoke the entry
	if err := m.revokeEntry(le); err != nil {
		m.logger.Printf("[ERR] expire: failed to revoke entry '%s': %v", vaultID, err)
	}

	// Delete the entry
	if err := m.deleteEntry(vaultID); err != nil {
		m.logger.Printf("[ERR] expire: failed to delete entry '%s': %v", vaultID, err)
	}
	m.logger.Printf("[INFO] expire: revoked '%s'", vaultID)
}

// revokeEntry is used to attempt revocation of an internal entry
func (m *ExpirationManager) revokeEntry(le *leaseEntry) error {
	req := &Request{
		Operation: RevokeOperation,
		Path:      le.Path,
		Data:      le.Data,
	}
	_, err := m.router.Route(req)
	return err
}

// loadEntry is used to read a lease entry
func (m *ExpirationManager) loadEntry(vaultID string) (*leaseEntry, error) {
	out, err := m.view.Get(vaultID)
	if err != nil {
		return nil, fmt.Errorf("failed to read lease entry: %v", err)
	}
	if out == nil {
		return nil, nil
	}
	le, err := decodeLeaseEntry(out.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode lease entry: %v", err)
	}
	return le, nil
}

// persistEntry is used to persist a lease entry
func (m *ExpirationManager) persistEntry(le *leaseEntry) error {
	// Encode the entry
	buf, err := le.encode()
	if err != nil {
		return fmt.Errorf("failed to encode lease entry: %v", err)
	}

	// Write out to the view
	ent := logical.StorageEntry{
		Key:   le.VaultID,
		Value: buf,
	}
	if err := m.view.Put(&ent); err != nil {
		return fmt.Errorf("failed to persist lease entry: %v", err)
	}
	return nil
}

// deleteEntry is used to delete a lease entry
func (m *ExpirationManager) deleteEntry(vaultID string) error {
	if err := m.view.Delete(vaultID); err != nil {
		return fmt.Errorf("failed to delete lease entry: %v", err)
	}
	return nil
}

// leaseEntry is used to structure the values the expiration
// manager stores. This is used to handle renew and revocation.
type leaseEntry struct {
	VaultID        string                 `json:"vault_id"`
	Path           string                 `json:"path"`
	Data           map[string]interface{} `json:"data"`
	Lease          *logical.Lease         `json:"lease"`
	IssueTime      time.Time              `json:"issue_time"`
	RenewTime      time.Time              `json:"renew_time"`
	RevokeAttempts int                    `json:"renew_attempts"`
}

// encode is used to JSON encode the lease entry
func (l *leaseEntry) encode() ([]byte, error) {
	return json.Marshal(l)
}

// decodeLeaseEntry is used to reverse encode and return a new entry
func decodeLeaseEntry(buf []byte) (*leaseEntry, error) {
	out := new(leaseEntry)
	return out, json.Unmarshal(buf, out)
}
