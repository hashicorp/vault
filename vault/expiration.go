package vault

import (
	"encoding/json"
	"fmt"
	"path"
	"time"
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
}

// NewExpirationManager creates a new ExpirationManager that is backed
// using a given view, and uses the provided router for revocation.
func NewExpirationManager(router *Router, view *BarrierView) *ExpirationManager {
	exp := &ExpirationManager{
		router: router,
		view:   view,
	}
	return exp
}

// setupExpiration is invoked after we've loaded the mount table to
// initialize the expiration manager
func (c *Core) setupExpiration() error {
	// Create a sub-view
	view := c.systemView.SubView(expirationSubPath)

	// Create the manager
	mgr := NewExpirationManager(c.router, view)
	c.expiration = mgr
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
func (m *ExpirationManager) Renew(vaultID string, increment time.Duration) (*Lease, error) {
	return nil, nil
}

// Register is used to take a request and response with an associated
// lease. The secret gets assigned a vaultId and the management of
// of lease is assumed by the expiration manager.
func (m *ExpirationManager) Register(req *Request, resp *Response) (string, error) {
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
	le := leaseEntry{
		VaultID:   path.Join(req.Path, generateUUID()),
		Path:      req.Path,
		Data:      resp.Data,
		Lease:     resp.Lease,
		IssueTime: time.Now().UTC(),
	}

	// Encode the entry
	buf, err := le.encode()
	if err != nil {
		return "", fmt.Errorf("failed to encode lease entry: %v", err)
	}

	// Write out to the view
	ent := Entry{
		Key:   le.VaultID,
		Value: buf,
	}
	if err := m.view.Put(&ent); err != nil {
		return "", fmt.Errorf("failed to persist lease entry: %v", err)
	}

	// TODO: Automatic revoke timer...

	// Done
	return le.VaultID, nil
}

// leaseEntry is used to structure the values the expiration
// manager stores. This is used to handle renew and revocation.
type leaseEntry struct {
	VaultID   string
	Path      string
	Data      map[string]interface{}
	Lease     *Lease
	IssueTime time.Time
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
