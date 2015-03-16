package vault

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/logical"
)

const (
	// expirationSubPath is the sub-path used for the expiration manager
	// view. This is nested under the system view.
	expirationSubPath = "expire/"

	// maxRevokeAttempts limits how many revoke attempts are made
	maxRevokeAttempts = 6

	// revokeRetryBase is a baseline retry time
	revokeRetryBase = 10 * time.Second

	// minRevokeDelay is used to prevent an instant revoke on restore
	minRevokeDelay = 5 * time.Second
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
	pendingLock sync.Mutex
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
	m.pendingLock.Lock()
	defer m.pendingLock.Unlock()

	// Accumulate existing leases
	var existing []string
	cb := func(path string) {
		existing = append(existing, path)
	}

	// Scan for all the leases
	if err := ScanView(m.view, cb); err != nil {
		return fmt.Errorf("failed to scan for leases: %v", err)
	}

	// Restore each key
	for _, vaultID := range existing {
		// Load the entry
		le, err := m.loadEntry(vaultID)
		if err != nil {
			return err
		}

		// If there is no entry, nothing to restore
		if le == nil {
			continue
		}

		// Determine the remaining time to expiration
		expires := time.Now().UTC().Sub(le.ExpireTime)
		if expires <= 0 {
			expires = minRevokeDelay
		}

		// Setup revocation timer
		m.pending[le.VaultID] = time.AfterFunc(expires, func() {
			m.expireID(le.VaultID)
		})
	}
	m.logger.Printf("[INFO] expire: restored %d leases", len(m.pending))
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
	// Load the entry
	le, err := m.loadEntry(vaultID)
	if err != nil {
		return err
	}

	// If there is no entry, nothing to revoke
	if le == nil {
		return nil
	}

	// Revoke the entry
	if err := m.revokeEntry(le); err != nil {
		return err
	}

	// Delete the entry
	if err := m.deleteEntry(vaultID); err != nil {
		return err
	}

	// Clear the expiration handler
	m.pendingLock.Lock()
	if timer, ok := m.pending[vaultID]; ok {
		timer.Stop()
		delete(m.pending, vaultID)
	}
	m.pendingLock.Unlock()
	return nil
}

// RevokePrefix is used to revoke all secrets with a given prefix.
// The prefix maps to that of the mount table to make this simpler
// to reason about.
func (m *ExpirationManager) RevokePrefix(prefix string) error {
	// Accumulate existing leases
	var existing []string
	cb := func(path string) {
		existing = append(existing, path)
	}

	// Ensure there is a trailing slash
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	// Scan for all the leases in the prefix
	if err := ScanView(m.view.SubView(prefix), cb); err != nil {
		return fmt.Errorf("failed to scan for leases: %v", err)
	}

	// Revoke all the keys
	for idx, vaultID := range existing {
		if err := m.Revoke(vaultID); err != nil {
			return fmt.Errorf("failed to revoke '%s' (%d / %d): %v",
				vaultID, idx+1, len(existing), err)
		}
	}
	return nil
}

// Renew is used to renew a secret using the given vaultID
// and a renew interval. The increment may be ignored.
func (m *ExpirationManager) Renew(vaultID string, increment time.Duration) (*logical.Response, error) {
	// Load the entry
	le, err := m.loadEntry(vaultID)
	if err != nil {
		return nil, err
	}

	// If there is no entry, cannot review
	if le == nil {
		return nil, fmt.Errorf("lease not found")
	}

	// Determine if the lease is expired
	if le.ExpireTime.Before(time.Now().UTC()) {
		return nil, fmt.Errorf("lease expired")
	}

	// Attempt to renew the entry
	resp, err := m.renewEntry(le)
	if err != nil {
		return nil, err
	}

	// Fast-path if there is no lease
	if resp == nil || resp.Lease == nil || !resp.IsSecret {
		return resp, nil
	}

	// Validate the lease
	if err := resp.Lease.Validate(); err != nil {
		return nil, err
	}

	// Update the lease entry
	le.Data = resp.Data
	le.Lease = resp.Lease
	le.ExpireTime = time.Now().UTC().Add(resp.Lease.Duration)
	if err := m.persistEntry(le); err != nil {
		return nil, err
	}

	// Update the expiration time
	m.pendingLock.Lock()
	if timer, ok := m.pending[vaultID]; ok {
		timer.Reset(resp.Lease.Duration)
	}
	m.pendingLock.Unlock()

	// Return the response
	return resp, nil
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
		VaultID:    path.Join(req.Path, generateUUID()),
		Path:       req.Path,
		Data:       resp.Data,
		Lease:      resp.Lease,
		IssueTime:  now,
		ExpireTime: now.Add(resp.Lease.Duration),
	}

	// Encode the entry
	if err := m.persistEntry(&le); err != nil {
		return "", err
	}

	// Setup revocation timer
	m.pendingLock.Lock()
	m.pending[le.VaultID] = time.AfterFunc(resp.Lease.Duration, func() {
		m.expireID(le.VaultID)
	})
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

	for attempt := uint(0); attempt < maxRevokeAttempts; attempt++ {
		err := m.Revoke(vaultID)
		if err == nil {
			m.logger.Printf("[INFO] expire: revoked '%s'", vaultID)
			return
		}
		m.logger.Printf("[ERR] expire: failed to revoke '%s': %v", vaultID, err)
		time.Sleep((1 << attempt) * revokeRetryBase)
	}
	m.logger.Printf("[ERR] expire: maximum revoke attempts for '%s' reached", vaultID)
}

// revokeEntry is used to attempt revocation of an internal entry
func (m *ExpirationManager) revokeEntry(le *leaseEntry) error {
	req := &logical.Request{
		Operation: logical.RevokeOperation,
		Path:      le.Path,
		Data:      le.Data,
	}
	_, err := m.router.Route(req)
	if err != nil {
		return fmt.Errorf("failed to revoke entry: %v", err)
	}
	return nil
}

// renewEntry is used to attempt renew of an internal entry
func (m *ExpirationManager) renewEntry(le *leaseEntry) (*logical.Response, error) {
	req := &logical.Request{
		Operation: logical.RenewOperation,
		Path:      le.Path,
		Data:      le.Data,
	}
	resp, err := m.router.Route(req)
	if err != nil {
		return nil, fmt.Errorf("failed to renew entry: %v", err)
	}
	return resp, nil
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
	VaultID    string                 `json:"vault_id"`
	Path       string                 `json:"path"`
	Data       map[string]interface{} `json:"data"`
	Lease      *logical.Lease         `json:"lease"`
	IssueTime  time.Time              `json:"issue_time"`
	ExpireTime time.Time              `json:"expire_time"`
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
