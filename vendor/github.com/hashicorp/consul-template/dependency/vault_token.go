package dependency

import (
	"log"
	"time"

	"github.com/hashicorp/vault/api"
)

var (
	// Ensure implements
	_ Dependency = (*VaultTokenQuery)(nil)
)

// VaultTokenQuery is the dependency to Vault for a secret
type VaultTokenQuery struct {
	stopCh      chan struct{}
	secret      *Secret
	vaultSecret *api.Secret
}

// NewVaultTokenQuery creates a new dependency.
func NewVaultTokenQuery(token string) (*VaultTokenQuery, error) {
	vaultSecret := &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken:   token,
			Renewable:     true,
			LeaseDuration: 1,
		},
	}
	return &VaultTokenQuery{
		stopCh:      make(chan struct{}, 1),
		vaultSecret: vaultSecret,
		secret:      transformSecret(vaultSecret),
	}, nil
}

// Fetch queries the Vault API
func (d *VaultTokenQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	if vaultSecretRenewable(d.secret) {
		renewSecret(clients, d)
	}

	// The secret isn't renewable, probably the generic secret backend.
	// TODO This is incorrect when given a non-renewable template. We should
	// instead to a lookup self to determine the lease duration.
	opts = opts.Merge(&QueryOptions{})
	dur := leaseCheckWait(d.secret)
	if dur < opts.VaultGrace {
		dur = opts.VaultGrace
	}

	log.Printf("[TRACE] %s: token is not renewable, sleeping for %s", d, dur)
	select {
	case <-time.After(dur):
	case <-d.stopCh:
		return nil, nil, ErrStopped
	}

	return nil, nil, ErrLeaseExpired
}

func (d *VaultTokenQuery) stopChan() chan struct{} {
	return d.stopCh
}

func (d *VaultTokenQuery) secrets() (*Secret, *api.Secret) {
	return d.secret, d.vaultSecret
}

// CanShare returns if this dependency is shareable.
func (d *VaultTokenQuery) CanShare() bool {
	return false
}

// Stop halts the dependency's fetch function.
func (d *VaultTokenQuery) Stop() {
	close(d.stopCh)
}

// String returns the human-friendly version of this dependency.
func (d *VaultTokenQuery) String() string {
	return "vault.token"
}

// Type returns the type of this dependency.
func (d *VaultTokenQuery) Type() Type {
	return TypeVault
}
