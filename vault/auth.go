package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/vault/credential"
)

const (
	// coreAuthConfigPath is used to store the auth configuration.
	// Auth configuration is protected within the Vault itself, which means it
	// can only be viewed or modified after an unseal.
	coreAuthConfigPath = "core/auth"

	// credentialBarrierPrefix is the prefix to the UUID used in the
	// barrier view for the credential backends.
	credentialBarrierPrefix = "auth/"

	// credentialMountPrefix is the mount prefix used for the router
	credentialMountPrefix = "auth/"
)

var (
	// loadAuthFailed if loadCreddentials encounters an error
	loadAuthFailed = errors.New("failed to setup auth table")
)

// AuthTable is used to represent the internal auth table
type AuthTable struct {
	// This lock should be held whenever modifying the Entries field.
	sync.RWMutex
	Entries []*AuthEntry `json:"entries"`
}

// Returns a deep copy of the auth table
func (t *AuthTable) Clone() *AuthTable {
	at := &AuthTable{
		Entries: make([]*AuthEntry, len(t.Entries)),
	}
	for i, e := range t.Entries {
		at.Entries[i] = e.Clone()
	}
	return at
}

// AuthEntry is used to represent an auth table entry
type AuthEntry struct {
	Name        string `json:"name "`       // Backend name (e.g. "github")
	Type        string `json:"type"`        // Credential backend Type (e.g. "oauth")
	Description string `json:"description"` // User-provided description
	UUID        string `json:"uuid"`        // Barrier view UUID
}

// Returns a deep copy of the auth entry
func (a *AuthEntry) Clone() *AuthEntry {
	return &AuthEntry{
		Name:        a.Name,
		Type:        a.Type,
		Description: a.Description,
		UUID:        a.UUID,
	}
}

// loadCredentials is invoked as part of postUnseal to load the auth table
func (c *Core) loadCredentials() error {
	// Load the existing mount table
	raw, err := c.barrier.Get(coreAuthConfigPath)
	if err != nil {
		c.logger.Printf("[ERR] core: failed to read auth table: %v", err)
		return loadAuthFailed
	}
	if raw != nil {
		c.auth = &AuthTable{}
		if err := json.Unmarshal(raw.Value, c.auth); err != nil {
			c.logger.Printf("[ERR] core: failed to decode auth table: %v", err)
			return loadAuthFailed
		}
	}

	// Done if we have restored the auth table
	if c.auth != nil {
		return nil
	}

	// Create and persist the default auth table
	c.auth = defaultAuthTable()
	if err := c.persistAuth(c.auth); err != nil {
		return loadAuthFailed
	}
	return nil
}

// persistAuth is used to persist the auth table after modification
func (c *Core) persistAuth(table *AuthTable) error {
	// Marshal the table
	raw, err := json.Marshal(table)
	if err != nil {
		c.logger.Printf("[ERR] core: failed to encode auth table: %v", err)
		return err
	}

	// Create an entry
	entry := &Entry{
		Key:   coreAuthConfigPath,
		Value: raw,
	}

	// Write to the physical backend
	if err := c.barrier.Put(entry); err != nil {
		c.logger.Printf("[ERR] core: failed to persist auth table: %v", err)
		return err
	}
	return nil
}

// setupCredentials is invoked after we've loaded the auth table to
// initialize the credential backends and setup the router
func (c *Core) setupCredentials() error {
	var backend credential.Backend
	var view *BarrierView
	var err error
	for _, entry := range c.auth.Entries {
		// Initialize the backend
		backend, err = c.newCredentialBackend(entry.Type, nil)
		if err != nil {
			c.logger.Printf(
				"[ERR] core: failed to create credential entry %#v: %v",
				entry, err)
			return loadAuthFailed
		}

		// Create a barrier view using the UUID
		view = NewBarrierView(c.barrier, credentialBarrierPrefix+entry.UUID+"/")

		// Mount the backend
		path := credentialMountPrefix + entry.Name
		err = c.router.Mount(backend, path, view)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to mount auth entry %#v: %v", entry, err)
			return loadAuthFailed
		}
	}
	return nil
}

// teardownCredentials is used before we seal the vault to reset the credential
// backends to their unloaded state. This is reversed by loadCredentials.
func (c *Core) teardownCredentials() error {
	c.auth = nil
	return nil
}

// newCredentialBackend is used to create and configure a new credential backend by name
func (c *Core) newCredentialBackend(t string, conf map[string]string) (credential.Backend, error) {
	f, ok := c.credentialBackends[t]
	if !ok {
		return nil, fmt.Errorf("unknown backend type: %s", t)
	}
	return f(conf)
}

// defaultAuthTable creates a default auth table
func defaultAuthTable() *AuthTable {
	table := &AuthTable{}
	tokenAuth := &AuthEntry{
		Name:        "token/",
		Type:        "token",
		Description: "token based credentials",
		UUID:        generateUUID(),
	}
	table.Entries = append(table.Entries, tokenAuth)
	return table
}
