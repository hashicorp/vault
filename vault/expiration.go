package vault

import "time"

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
	return "", nil
}
