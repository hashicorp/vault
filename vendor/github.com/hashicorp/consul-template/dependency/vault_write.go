package dependency

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

var (
	// Ensure implements
	_ Dependency = (*VaultWriteQuery)(nil)
)

// VaultWriteQuery is the dependency to Vault for a secret
type VaultWriteQuery struct {
	stopCh  chan struct{}
	sleepCh chan time.Duration

	path     string
	data     map[string]interface{}
	dataHash string
	secret   *Secret

	// vaultSecret is the actual Vault secret which we are renewing
	vaultSecret *api.Secret
}

// NewVaultWriteQuery creates a new datacenter dependency.
func NewVaultWriteQuery(s string, d map[string]interface{}) (*VaultWriteQuery, error) {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "/")
	if s == "" {
		return nil, fmt.Errorf("vault.write: invalid format: %q", s)
	}

	return &VaultWriteQuery{
		stopCh:   make(chan struct{}, 1),
		sleepCh:  make(chan time.Duration, 1),
		path:     s,
		data:     d,
		dataHash: sha1Map(d),
	}, nil
}

// Fetch queries the Vault API
func (d *VaultWriteQuery) Fetch(clients *ClientSet, opts *QueryOptions,
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

	opts = opts.Merge(&QueryOptions{})
	vaultSecret, err := d.writeSecret(clients, opts)
	if err != nil {
		return nil, nil, errors.Wrap(err, d.String())
	}

	// vaultSecret == nil when writing to KVv1 engines
	if vaultSecret == nil {
		return respWithMetadata(d.secret)
	}

	printVaultWarnings(d, vaultSecret.Warnings)
	d.vaultSecret = vaultSecret
	// cloned secret which will be exposed to the template
	d.secret = transformSecret(vaultSecret)

	if !vaultSecretRenewable(d.secret) {
		dur := leaseCheckWait(d.secret)
		log.Printf("[TRACE] %s: non-renewable secret, set sleep for %s", d, dur)
		d.sleepCh <- dur
	}

	return respWithMetadata(d.secret)
}

// meet renewer interface
func (d *VaultWriteQuery) stopChan() chan struct{} {
	return d.stopCh
}

func (d *VaultWriteQuery) secrets() (*Secret, *api.Secret) {
	return d.secret, d.vaultSecret
}

// CanShare returns if this dependency is shareable.
func (d *VaultWriteQuery) CanShare() bool {
	return false
}

// Stop halts the given dependency's fetch.
func (d *VaultWriteQuery) Stop() {
	close(d.stopCh)
}

// String returns the human-friendly version of this dependency.
func (d *VaultWriteQuery) String() string {
	return fmt.Sprintf("vault.write(%s -> %s)", d.path, d.dataHash)
}

// Type returns the type of this dependency.
func (d *VaultWriteQuery) Type() Type {
	return TypeVault
}

// sha1Map returns the sha1 hash of the data in the map. The reason this data is
// hashed is because it appears in the output and could contain sensitive
// information.
func sha1Map(m map[string]interface{}) string {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha1.New()
	for _, k := range keys {
		io.WriteString(h, fmt.Sprintf("%s=%q", k, m[k]))
	}

	return fmt.Sprintf("%.4x", h.Sum(nil))
}

func (d *VaultWriteQuery) printWarnings(warnings []string) {
	for _, w := range warnings {
		log.Printf("[WARN] %s: %s", d, w)
	}
}

func (d *VaultWriteQuery) writeSecret(clients *ClientSet, opts *QueryOptions) (*api.Secret, error) {
	log.Printf("[TRACE] %s: PUT %s", d, &url.URL{
		Path:     "/v1/" + d.path,
		RawQuery: opts.String(),
	})

	data := d.data

	_, isv2, _ := isKVv2(clients.Vault(), d.path)
	if isv2 {
		data = map[string]interface{}{"data": d.data}
	}

	vaultSecret, err := clients.Vault().Logical().Write(d.path, data)
	if err != nil {
		return nil, errors.Wrap(err, d.String())
	}
	// vaultSecret is always nil when KVv1 engine (isv2==false)
	if isv2 && vaultSecret == nil {
		return nil, fmt.Errorf("no secret exists at %s", d.path)
	}

	return vaultSecret, nil
}
