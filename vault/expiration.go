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

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/uuid"
	"github.com/hashicorp/vault/logical"
)

const (
	// expirationSubPath is the sub-path used for the expiration manager
	// view. This is nested under the system view.
	expirationSubPath = "expire/"

	// leaseViewPrefix is the prefix used for the ID based lookup of leases.
	leaseViewPrefix = "id/"

	// tokenViewPrefix is the prefix used for the token based lookup of leases.
	tokenViewPrefix = "token/"

	// maxRevokeAttempts limits how many revoke attempts are made
	maxRevokeAttempts = 6

	// revokeRetryBase is a baseline retry time
	revokeRetryBase = 10 * time.Second

	// minRevokeDelay is used to prevent an instant revoke on restore
	minRevokeDelay = 5 * time.Second

	// maxLeaseDuration is the default maximum lease duration
	maxLeaseTTL = 30 * 24 * time.Hour

	// defaultLeaseDuration is the default lease duration used when no lease is specified
	defaultLeaseTTL = maxLeaseTTL
)

// ExpirationManager is used by the Core to manage leases. Secrets
// can provide a lease, meaning that they can be renewed or revoked.
// If a secret is not renewed in timely manner, it may be expired, and
// the ExpirationManager will handle doing automatic revocation.
type ExpirationManager struct {
	router     *Router
	idView     *BarrierView
	tokenView  *BarrierView
	tokenStore *TokenStore
	logger     *log.Logger

	pending     map[string]*time.Timer
	pendingLock sync.Mutex
}

// NewExpirationManager creates a new ExpirationManager that is backed
// using a given view, and uses the provided router for revocation.
func NewExpirationManager(router *Router, view *BarrierView, ts *TokenStore, logger *log.Logger) *ExpirationManager {
	if logger == nil {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	exp := &ExpirationManager{
		router:     router,
		idView:     view.SubView(leaseViewPrefix),
		tokenView:  view.SubView(tokenViewPrefix),
		tokenStore: ts,
		logger:     logger,
		pending:    make(map[string]*time.Timer),
	}
	return exp
}

// setupExpiration is invoked after we've loaded the mount table to
// initialize the expiration manager
func (c *Core) setupExpiration() error {
	// Create a sub-view
	view := c.systemBarrierView.SubView(expirationSubPath)

	// Create the manager
	mgr := NewExpirationManager(c.router, view, c.tokenStore, c.logger)
	c.expiration = mgr

	// Link the token store to this
	c.tokenStore.SetExpirationManager(mgr)

	// Restore the existing state
	if err := c.expiration.Restore(); err != nil {
		return fmt.Errorf("expiration state restore failed: %v", err)
	}
	return nil
}

// stopExpiration is used to stop the expiration manager before
// sealing the Vault.
func (c *Core) stopExpiration() error {
	if c.expiration != nil {
		if err := c.expiration.Stop(); err != nil {
			return err
		}
		c.expiration = nil
	}
	return nil
}

// Restore is used to recover the lease states when starting.
// This is used after starting the vault.
func (m *ExpirationManager) Restore() error {
	m.pendingLock.Lock()
	defer m.pendingLock.Unlock()

	// Accumulate existing leases
	existing, err := CollectKeys(m.idView)
	if err != nil {
		return fmt.Errorf("failed to scan for leases: %v", err)
	}

	// Restore each key
	for _, leaseID := range existing {
		// Load the entry
		le, err := m.loadEntry(leaseID)
		if err != nil {
			return err
		}

		// If there is no entry, nothing to restore
		if le == nil {
			continue
		}

		// If there is no expiry time, don't do anything
		if le.ExpireTime.IsZero() {
			continue
		}

		// Determine the remaining time to expiration
		expires := le.ExpireTime.Sub(time.Now().UTC())
		if expires <= 0 {
			expires = minRevokeDelay
		}

		// Setup revocation timer
		m.pending[le.LeaseID] = time.AfterFunc(expires, func() {
			m.expireID(le.LeaseID)
		})
	}
	if len(m.pending) > 0 {
		m.logger.Printf("[INFO] expire: restored %d leases", len(m.pending))
	}
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

// Revoke is used to revoke a secret named by the given LeaseID
func (m *ExpirationManager) Revoke(leaseID string) error {
	defer metrics.MeasureSince([]string{"expire", "revoke"}, time.Now())
	// Load the entry
	le, err := m.loadEntry(leaseID)
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
	if err := m.deleteEntry(leaseID); err != nil {
		return err
	}

	// Delete the secondary index
	if err := m.indexByToken(le.ClientToken, le.LeaseID); err != nil {
		return err
	}

	// Clear the expiration handler
	m.pendingLock.Lock()
	if timer, ok := m.pending[leaseID]; ok {
		timer.Stop()
		delete(m.pending, leaseID)
	}
	m.pendingLock.Unlock()
	return nil
}

// RevokePrefix is used to revoke all secrets with a given prefix.
// The prefix maps to that of the mount table to make this simpler
// to reason about.
func (m *ExpirationManager) RevokePrefix(prefix string) error {
	defer metrics.MeasureSince([]string{"expire", "revoke-prefix"}, time.Now())
	// Ensure there is a trailing slash
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	// Accumulate existing leases
	sub := m.idView.SubView(prefix)
	existing, err := CollectKeys(sub)
	if err != nil {
		return fmt.Errorf("failed to scan for leases: %v", err)
	}

	// Revoke all the keys
	for idx, suffix := range existing {
		leaseID := prefix + suffix
		if err := m.Revoke(leaseID); err != nil {
			return fmt.Errorf("failed to revoke '%s' (%d / %d): %v",
				leaseID, idx+1, len(existing), err)
		}
	}
	return nil
}

// RevokeByToken is used to revoke all the secrets issued with
// a given token. This is done by using the secondary index.
func (m *ExpirationManager) RevokeByToken(token string) error {
	defer metrics.MeasureSince([]string{"expire", "revoke-by-token"}, time.Now())
	// Lookup the leases
	existing, err := m.lookupByToken(token)
	if err != nil {
		return fmt.Errorf("failed to scan for leases: %v", err)
	}

	// Revoke all the keys
	for idx, leaseID := range existing {
		if err := m.Revoke(leaseID); err != nil {
			return fmt.Errorf("failed to revoke '%s' (%d / %d): %v",
				leaseID, idx+1, len(existing), err)
		}
	}
	return nil
}

// Renew is used to renew a secret using the given leaseID
// and a renew interval. The increment may be ignored.
func (m *ExpirationManager) Renew(leaseID string, increment time.Duration) (*logical.Response, error) {
	defer metrics.MeasureSince([]string{"expire", "renew"}, time.Now())
	// Load the entry
	le, err := m.loadEntry(leaseID)
	if err != nil {
		return nil, err
	}

	// Check if the lease is renewable
	if err := le.renewable(); err != nil {
		return nil, err
	}

	// Attempt to renew the entry
	resp, err := m.renewEntry(le, increment)
	if err != nil {
		return nil, err
	}

	// Fast-path if there is no lease
	if resp == nil || resp.Secret == nil || !resp.Secret.LeaseEnabled() {
		return resp, nil
	}

	// Validate the lease
	if err := resp.Secret.Validate(); err != nil {
		return nil, err
	}

	// Attach the LeaseID
	resp.Secret.LeaseID = leaseID

	// Update the lease entry
	le.Data = resp.Data
	le.Secret = resp.Secret
	le.ExpireTime = resp.Secret.ExpirationTime()
	if err := m.persistEntry(le); err != nil {
		return nil, err
	}

	// Update the expiration time
	m.updatePending(le, resp.Secret.LeaseTotal())

	// Return the response
	return resp, nil
}

// RenewToken is used to renew a token which does not need to
// invoke a logical backend.
func (m *ExpirationManager) RenewToken(source string, token string,
	increment time.Duration) (*logical.Auth, error) {
	defer metrics.MeasureSince([]string{"expire", "renew-token"}, time.Now())
	// Compute the Lease ID
	leaseID := path.Join(source, m.tokenStore.SaltID(token))

	// Load the entry
	le, err := m.loadEntry(leaseID)
	if err != nil {
		return nil, err
	}

	// Check if the lease is renewable
	if err := le.renewable(); err != nil {
		return nil, err
	}

	// Attempt to renew the auth entry
	resp, err := m.renewAuthEntry(le, increment)
	if err != nil {
		return nil, err
	}

	// Fast-path if there is no renewal
	if resp == nil {
		return nil, nil
	}
	if resp.Auth == nil || !resp.Auth.LeaseEnabled() {
		return resp.Auth, nil
	}

	// Attach the ClientToken
	resp.Auth.ClientToken = token
	resp.Auth.Increment = 0

	// Update the lease entry
	le.Auth = resp.Auth
	le.ExpireTime = resp.Auth.ExpirationTime()
	if err := m.persistEntry(le); err != nil {
		return nil, err
	}

	// Update the expiration time
	m.updatePending(le, resp.Auth.LeaseTotal())
	return resp.Auth, nil
}

// Register is used to take a request and response with an associated
// lease. The secret gets assigned a LeaseID and the management of
// of lease is assumed by the expiration manager.
func (m *ExpirationManager) Register(req *logical.Request, resp *logical.Response) (string, error) {
	defer metrics.MeasureSince([]string{"expire", "register"}, time.Now())
	// Ignore if there is no leased secret
	if resp == nil || resp.Secret == nil {
		return "", nil
	}

	// Validate the secret
	if err := resp.Secret.Validate(); err != nil {
		return "", err
	}

	// Create a lease entry
	le := leaseEntry{
		LeaseID:     path.Join(req.Path, uuid.GenerateUUID()),
		ClientToken: req.ClientToken,
		Path:        req.Path,
		Data:        resp.Data,
		Secret:      resp.Secret,
		IssueTime:   time.Now().UTC(),
		ExpireTime:  resp.Secret.ExpirationTime(),
	}

	// Encode the entry
	if err := m.persistEntry(&le); err != nil {
		return "", err
	}

	// Maintain secondary index by token
	if err := m.indexByToken(le.ClientToken, le.LeaseID); err != nil {
		return "", err
	}

	// Setup revocation timer if there is a lease
	m.updatePending(&le, resp.Secret.LeaseTotal())

	// Done
	return le.LeaseID, nil
}

// RegisterAuth is used to take an Auth response with an associated lease.
// The token does not get a LeaseID, but the lease management is handled by
// the expiration manager.
func (m *ExpirationManager) RegisterAuth(source string, auth *logical.Auth) error {
	defer metrics.MeasureSince([]string{"expire", "register-auth"}, time.Now())

	// Create a lease entry
	le := leaseEntry{
		LeaseID:     path.Join(source, m.tokenStore.SaltID(auth.ClientToken)),
		ClientToken: auth.ClientToken,
		Auth:        auth,
		Path:        source,
		IssueTime:   time.Now().UTC(),
		ExpireTime:  auth.ExpirationTime(),
	}

	// Encode the entry
	if err := m.persistEntry(&le); err != nil {
		return err
	}

	// Setup revocation timer
	m.updatePending(&le, auth.LeaseTotal())
	return nil
}

// updatePending is used to update a pending invocation for a lease
func (m *ExpirationManager) updatePending(le *leaseEntry, leaseTotal time.Duration) {
	m.pendingLock.Lock()
	defer m.pendingLock.Unlock()

	// Check for an existing timer
	timer, ok := m.pending[le.LeaseID]

	// Create entry if it does not exist
	if !ok && leaseTotal > 0 {
		timer := time.AfterFunc(leaseTotal, func() {
			m.expireID(le.LeaseID)
		})
		m.pending[le.LeaseID] = timer
		return
	}

	// Delete the timer if the expiration time is zero
	if ok && leaseTotal == 0 {
		timer.Stop()
		delete(m.pending, le.LeaseID)
		return
	}

	// Extend the timer by the lease total
	if ok && leaseTotal > 0 {
		timer.Reset(leaseTotal)
	}
}

// expireID is invoked when a given ID is expired
func (m *ExpirationManager) expireID(leaseID string) {
	// Clear from the pending expiration
	m.pendingLock.Lock()
	delete(m.pending, leaseID)
	m.pendingLock.Unlock()

	for attempt := uint(0); attempt < maxRevokeAttempts; attempt++ {
		err := m.Revoke(leaseID)
		if err == nil {
			m.logger.Printf("[INFO] expire: revoked '%s'", leaseID)
			return
		}
		m.logger.Printf("[ERR] expire: failed to revoke '%s': %v", leaseID, err)
		time.Sleep((1 << attempt) * revokeRetryBase)
	}
	m.logger.Printf("[ERR] expire: maximum revoke attempts for '%s' reached", leaseID)
}

// revokeEntry is used to attempt revocation of an internal entry
func (m *ExpirationManager) revokeEntry(le *leaseEntry) error {
	// Revocation of login tokens is special since we can by-pass the
	// backend and directly interact with the token store
	if le.Auth != nil {
		if err := m.tokenStore.RevokeTree(le.Auth.ClientToken); err != nil {
			return fmt.Errorf("failed to revoke token: %v", err)
		}
		return nil
	}

	// Handle standard revocation via backends
	_, err := m.router.Route(logical.RevokeRequest(
		le.Path, le.Secret, le.Data))
	if err != nil {
		return fmt.Errorf("failed to revoke entry: %v", err)
	}
	return nil
}

// renewEntry is used to attempt renew of an internal entry
func (m *ExpirationManager) renewEntry(le *leaseEntry, increment time.Duration) (*logical.Response, error) {
	secret := *le.Secret
	secret.IssueTime = le.IssueTime
	secret.Increment = increment
	secret.LeaseID = ""

	req := logical.RenewRequest(le.Path, &secret, le.Data)
	resp, err := m.router.Route(req)
	if err != nil {
		return nil, fmt.Errorf("failed to renew entry: %v", err)
	}
	return resp, nil
}

// renewAuthEntry is used to attempt renew of an auth entry
func (m *ExpirationManager) renewAuthEntry(le *leaseEntry, increment time.Duration) (*logical.Response, error) {
	auth := *le.Auth
	auth.IssueTime = le.IssueTime
	auth.Increment = increment
	auth.ClientToken = ""

	req := logical.RenewAuthRequest(le.Path, &auth, nil)
	resp, err := m.router.Route(req)
	if err != nil {
		return nil, fmt.Errorf("failed to renew entry: %v", err)
	}
	return resp, nil
}

// loadEntry is used to read a lease entry
func (m *ExpirationManager) loadEntry(leaseID string) (*leaseEntry, error) {
	out, err := m.idView.Get(leaseID)
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
		Key:   le.LeaseID,
		Value: buf,
	}
	if err := m.idView.Put(&ent); err != nil {
		return fmt.Errorf("failed to persist lease entry: %v", err)
	}
	return nil
}

// deleteEntry is used to delete a lease entry
func (m *ExpirationManager) deleteEntry(leaseID string) error {
	if err := m.idView.Delete(leaseID); err != nil {
		return fmt.Errorf("failed to delete lease entry: %v", err)
	}
	return nil
}

// indexByToken creates a secondary index from the token to a lease entry
func (m *ExpirationManager) indexByToken(token, leaseID string) error {
	ent := logical.StorageEntry{
		Key:   m.tokenStore.SaltID(token) + "/" + m.tokenStore.SaltID(leaseID),
		Value: []byte(leaseID),
	}
	if err := m.tokenView.Put(&ent); err != nil {
		return fmt.Errorf("failed to persist lease index entry: %v", err)
	}
	return nil
}

// removeIndexByToken removes the secondary index from the token to a lease entry
func (m *ExpirationManager) removeIndexByToken(token, leaseID string) error {
	key := m.tokenStore.SaltID(token) + "/" + m.tokenStore.SaltID(leaseID)
	if err := m.tokenView.Delete(key); err != nil {
		return fmt.Errorf("failed to delete lease index entry: %v", err)
	}
	return nil
}

// lookupByToken is used to lookup all the leaseID's via the
func (m *ExpirationManager) lookupByToken(token string) ([]string, error) {
	// Scan via the index for sub-leases
	prefix := m.tokenStore.SaltID(token) + "/"
	subKeys, err := m.tokenView.List(prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to list leases: %v", err)
	}

	// Read each index entry
	leaseIDs := make([]string, 0, len(subKeys))
	for _, sub := range subKeys {
		out, err := m.tokenView.Get(prefix + sub)
		if err != nil {
			return nil, fmt.Errorf("failed to read lease index: %v", err)
		}
		if out == nil {
			continue
		}
		leaseIDs = append(leaseIDs, string(out.Value))
	}
	return leaseIDs, nil
}

// emitMetrics is invoked periodically to emit statistics
func (m *ExpirationManager) emitMetrics() {
	m.pendingLock.Lock()
	num := len(m.pending)
	m.pendingLock.Unlock()
	metrics.SetGauge([]string{"expire", "num_leases"}, float32(num))
}

// leaseEntry is used to structure the values the expiration
// manager stores. This is used to handle renew and revocation.
type leaseEntry struct {
	LeaseID     string                 `json:"lease_id"`
	ClientToken string                 `json:"client_token"`
	Path        string                 `json:"path"`
	Data        map[string]interface{} `json:"data"`
	Secret      *logical.Secret        `json:"secret"`
	Auth        *logical.Auth          `json:"auth"`
	IssueTime   time.Time              `json:"issue_time"`
	ExpireTime  time.Time              `json:"expire_time"`
}

// encode is used to JSON encode the lease entry
func (l *leaseEntry) encode() ([]byte, error) {
	return json.Marshal(l)
}

func (le *leaseEntry) renewable() error {
	// If there is no entry, cannot review
	if le == nil || le.ExpireTime.IsZero() {
		return fmt.Errorf("lease not found or lease is not renewable")
	}

	// Determine if the lease is expired
	if le.ExpireTime.Before(time.Now().UTC()) {
		return fmt.Errorf("lease expired")
	}

	// Determine if the lease is renewable
	if le.Secret != nil && !le.Secret.Renewable {
		return fmt.Errorf("lease is not renewable")
	}
	if le.Auth != nil && !le.Auth.Renewable {
		return fmt.Errorf("lease is not renewable")
	}
	return nil
}

// decodeLeaseEntry is used to reverse encode and return a new entry
func decodeLeaseEntry(buf []byte) (*leaseEntry, error) {
	out := new(leaseEntry)
	return out, json.Unmarshal(buf, out)
}
