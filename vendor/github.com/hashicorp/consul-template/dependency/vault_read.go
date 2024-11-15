// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dependency

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

// Ensure implements
var _ Dependency = (*VaultReadQuery)(nil)

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
		select {
		case <-time.After(dur):
			break
		case <-d.stopCh:
			return nil, nil, ErrStopped
		}
	default:
	}

	firstRun := d.secret == nil

	if !firstRun && vaultSecretRenewable(d.secret) {
		err := renewSecret(clients, d)
		if err != nil {
			return nil, nil, errors.Wrap(err, d.String())
		}
	}

	err := d.fetchSecret(clients)
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

func (d *VaultReadQuery) fetchSecret(clients *ClientSet) error {
	vaultSecret, err := d.readSecret(clients)
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
	if v := d.queryValues["version"]; len(v) > 0 {
		return fmt.Sprintf("vault.read(%s.v%s)", d.rawPath, v[0])
	}
	return fmt.Sprintf("vault.read(%s)", d.rawPath)
}

// Type returns the type of this dependency.
func (d *VaultReadQuery) Type() Type {
	return TypeVault
}

func (d *VaultReadQuery) readSecret(clients *ClientSet) (*api.Secret, error) {
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
			d.secretPath = shimKVv2Path(d.rawPath, mountPath, clients.Vault().Namespace())
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
		deletionTime, ok := md["deletion_time"].(string)
		if !ok {
			// Key not present or not a string, end early
			return false
		}
		t, err := time.Parse(time.RFC3339, deletionTime)
		if err != nil {
			// Deletion time is either empty, or not a valid string.
			return false
		}

		// If now is 'after' the deletion time, then the secret
		// should be deleted.
		return time.Now().After(t)
	}
	return false
}

// shimKVv2Path aligns the supported legacy path to KV v2 specs by inserting
// /data/ into the path for reading secrets. Paths for metadata are not modified.
func shimKVv2Path(rawPath, mountPath, clientNamespace string) string {
	switch {
	case rawPath == mountPath, rawPath == strings.TrimSuffix(mountPath, "/"):
		return path.Join(mountPath, "data")
	default:

		// Canonicalize the client namespace path to always having a '/' suffix
		if !strings.HasSuffix(clientNamespace, "/") {
			clientNamespace += "/"
		}
		// Extract client namespace from mount path if it exists
		rawPathNsAndMountPath := strings.TrimPrefix(mountPath, clientNamespace)

		// Trim (mount path - client namespace) from the raw path
		p := strings.TrimPrefix(rawPath, rawPathNsAndMountPath)

		// Only add /data/ prefix to the path if neither /data/ or /metadata/ are
		// present.
		if strings.HasPrefix(p, "data/") || strings.HasPrefix(p, "metadata/") {
			return rawPath
		}

		return path.Join(rawPathNsAndMountPath, "data", p)
	}
}
