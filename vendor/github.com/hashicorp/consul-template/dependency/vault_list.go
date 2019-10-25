package dependency

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	// Ensure implements
	_ Dependency = (*VaultListQuery)(nil)
)

// VaultListQuery is the dependency to Vault for a secret
type VaultListQuery struct {
	stopCh chan struct{}

	path string
}

// NewVaultListQuery creates a new datacenter dependency.
func NewVaultListQuery(s string) (*VaultListQuery, error) {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "/")
	if s == "" {
		return nil, fmt.Errorf("vault.list: invalid format: %q", s)
	}

	return &VaultListQuery{
		stopCh: make(chan struct{}, 1),
		path:   s,
	}, nil
}

// Fetch queries the Vault API
func (d *VaultListQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	opts = opts.Merge(&QueryOptions{})

	// If this is not the first query, poll to simulate blocking-queries.
	if opts.WaitIndex != 0 {
		dur := VaultDefaultLeaseDuration
		log.Printf("[TRACE] %s: long polling for %s", d, dur)

		select {
		case <-d.stopCh:
			return nil, nil, ErrStopped
		case <-time.After(dur):
		}
	}

	// If we got this far, we either didn't have a secret to renew, the secret was
	// not renewable, or the renewal failed, so attempt a fresh list.
	log.Printf("[TRACE] %s: LIST %s", d, &url.URL{
		Path:     "/v1/" + d.path,
		RawQuery: opts.String(),
	})
	secret, err := clients.Vault().Logical().List(d.path)
	if err != nil {
		return nil, nil, errors.Wrap(err, d.String())
	}

	var result []string

	// The secret could be nil if it does not exist.
	if secret == nil || secret.Data == nil {
		log.Printf("[TRACE] %s: no data", d)
		return respWithMetadata(result)
	}

	// This is a weird thing that happened once...
	keys, ok := secret.Data["keys"]
	if !ok {
		log.Printf("[TRACE] %s: no keys", d)
		return respWithMetadata(result)
	}

	list, ok := keys.([]interface{})
	if !ok {
		log.Printf("[TRACE] %s: not list", d)
		return nil, nil, fmt.Errorf("%s: unexpected response", d)
	}

	for _, v := range list {
		typed, ok := v.(string)
		if !ok {
			return nil, nil, fmt.Errorf("%s: non-string in list", d)
		}
		result = append(result, typed)
	}
	sort.Strings(result)

	log.Printf("[TRACE] %s: returned %d results", d, len(result))

	return respWithMetadata(result)
}

// CanShare returns if this dependency is shareable.
func (d *VaultListQuery) CanShare() bool {
	return false
}

// Stop halts the given dependency's fetch.
func (d *VaultListQuery) Stop() {
	close(d.stopCh)
}

// String returns the human-friendly version of this dependency.
func (d *VaultListQuery) String() string {
	return fmt.Sprintf("vault.list(%s)", d.path)
}

// Type returns the type of this dependency.
func (d *VaultListQuery) Type() Type {
	return TypeVault
}
