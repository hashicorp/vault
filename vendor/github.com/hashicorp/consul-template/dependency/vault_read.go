package dependency

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

var (
	// Ensure implements
	_ Dependency = (*VaultReadQuery)(nil)
)

// VaultReadQuery is the dependency to Vault for a secret
type VaultReadQuery struct {
	stopCh  chan struct{}
	sleepCh chan time.Duration

	rawPath     string
	queryValues url.Values
	secret      *Secret
	isKVv2      *bool
	secretPath  string

	// vaultSecret is the actual Vault secret which we are renewing
	vaultSecret *api.Secret
}

// NewVaultReadQuery creates a new datacenter dependency.
func NewVaultReadQuery(s string) (*VaultReadQuery, error) {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "/")
	if s == "" {
		return nil, fmt.Errorf("vault.read: invalid format: %q", s)
	}

	secretURL, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	return &VaultReadQuery{
		stopCh:      make(chan struct{}, 1),
		sleepCh:     make(chan time.Duration, 1),
		rawPath:     secretURL.Path,
		queryValues: secretURL.Query(),
	}, nil
}

// Fetch queries the Vault API
func (d *VaultReadQuery) Fetch(clients *ClientSet, opts *QueryOptions,
) (interface{}, *ResponseMetadata, error) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}
	select {
	case dur := <-d.sleepCh:
		time.Sleep(dur)
	default:
	}

	firstRun := d.secret == nil

	if !firstRun && vaultSecretRenewable(d.secret) {
		err := renewSecret(clients, d)
		if err != nil {
			return nil, nil, errors.Wrap(err, d.String())
		}
	}

	err := d.fetchSecret(clients, opts)
	if err != nil {
		return nil, nil, errors.Wrap(err, d.String())
	}

	if !vaultSecretRenewable(d.secret) {
		dur := leaseCheckWait(d.secret)
		log.Printf("[TRACE] %s: non-renewable secret, set sleep for %s", d, dur)
		d.sleepCh <- dur
	}

	return respWithMetadata(d.secret)
}

func (d *VaultReadQuery) fetchSecret(clients *ClientSet, opts *QueryOptions,
) error {
	opts = opts.Merge(&QueryOptions{})
	vaultSecret, err := d.readSecret(clients, opts)
	if err == nil {
		printVaultWarnings(d, vaultSecret.Warnings)
		d.vaultSecret = vaultSecret
		// the cloned secret which will be exposed to the template
		d.secret = transformSecret(vaultSecret)
	}
	return err
}

func (d *VaultReadQuery) stopChan() chan struct{} {
	return d.stopCh
}

func (d *VaultReadQuery) secrets() (*Secret, *api.Secret) {
	return d.secret, d.vaultSecret
}

// CanShare returns if this dependency is shareable.
func (d *VaultReadQuery) CanShare() bool {
	return false
}

// Stop halts the given dependency's fetch.
func (d *VaultReadQuery) Stop() {
	close(d.stopCh)
}

// String returns the human-friendly version of this dependency.
func (d *VaultReadQuery) String() string {
	return fmt.Sprintf("vault.read(%s)", d.rawPath)
}

// Type returns the type of this dependency.
func (d *VaultReadQuery) Type() Type {
	return TypeVault
}

func (d *VaultReadQuery) readSecret(clients *ClientSet, opts *QueryOptions) (*api.Secret, error) {
	vaultClient := clients.Vault()

	// Check whether this secret refers to a KV v2 entry if we haven't yet.
	if d.isKVv2 == nil {
		mountPath, isKVv2, err := isKVv2(vaultClient, d.rawPath)
		if err != nil {
			log.Printf("[WARN] %s: failed to check if %s is KVv2, "+
				"assume not: %s", d, d.rawPath, err)
			isKVv2 = false
			d.secretPath = d.rawPath
		} else if isKVv2 {
			d.secretPath = addPrefixToVKVPath(d.rawPath, mountPath, "data")
		} else {
			d.secretPath = d.rawPath
		}
		d.isKVv2 = &isKVv2
	}

	queryString := d.queryValues.Encode()
	log.Printf("[TRACE] %s: GET %s", d, &url.URL{
		Path:     "/v1/" + d.secretPath,
		RawQuery: queryString,
	})
	vaultSecret, err := vaultClient.Logical().ReadWithData(d.secretPath,
		d.queryValues)

	if err != nil {
		return nil, errors.Wrap(err, d.String())
	}
	if vaultSecret == nil || deletedKVv2(vaultSecret) {
		return nil, fmt.Errorf("no secret exists at %s", d.secretPath)
	}
	return vaultSecret, nil
}

func deletedKVv2(s *api.Secret) bool {
	switch md := s.Data["metadata"].(type) {
	case map[string]interface{}:
		return md["deletion_time"] != ""
	}
	return false
}
