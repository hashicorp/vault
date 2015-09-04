package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/uuid"
	"github.com/hashicorp/vault/logical"
)

const (
	// coreAuthConfigPath is used to store the auth configuration.
	// Auth configuration is protected within the Vault itself, which means it
	// can only be viewed or modified after an unseal.
	coreAuthConfigPath = "core/auth"

	// credentialBarrierPrefix is the prefix to the UUID used in the
	// barrier view for the credential backends.
	credentialBarrierPrefix = "auth/"

	// credentialRoutePrefix is the mount prefix used for the router
	credentialRoutePrefix = "auth/"
)

var (
	// loadAuthFailed if loadCreddentials encounters an error
	loadAuthFailed = errors.New("failed to setup auth table")
)

// enableCredential is used to enable a new credential backend
func (c *Core) enableCredential(entry *MountEntry) error {
	c.auth.Lock()
	defer c.auth.Unlock()

	// Ensure we end the path in a slash
	if !strings.HasSuffix(entry.Path, "/") {
		entry.Path += "/"
	}

	// Ensure there is a name
	if entry.Path == "/" {
		return fmt.Errorf("backend path must be specified")
	}

	// Look for matching name
	for _, ent := range c.auth.Entries {
		switch {
		// Existing is oauth/github/ new is oauth/ or
		// existing is oauth/ and new is oauth/github/
		case strings.HasPrefix(ent.Path, entry.Path):
			fallthrough
		case strings.HasPrefix(entry.Path, ent.Path):
			return logical.CodedError(409, "path is already in use")
		}
	}

	// Ensure the token backend is a singleton
	if entry.Type == "token" {
		return fmt.Errorf("token credential backend cannot be instantiated")
	}

	// Generate a new UUID and view
	entry.UUID = uuid.GenerateUUID()
	view := NewBarrierView(c.barrier, credentialBarrierPrefix+entry.UUID+"/")

	// Create the new backend
	backend, err := c.newCredentialBackend(entry.Type, c.mountEntrySysView(entry), view, nil)
	if err != nil {
		return err
	}

	// Update the auth table
	newTable := c.auth.Clone()
	newTable.Entries = append(newTable.Entries, entry)
	if err := c.persistAuth(newTable); err != nil {
		return errors.New("failed to update auth table")
	}
	c.auth = newTable

	// Mount the backend
	path := credentialRoutePrefix + entry.Path
	if err := c.router.Mount(backend, path, entry, view); err != nil {
		return err
	}
	c.logger.Printf("[INFO] core: enabled credential backend '%s' type: %s",
		entry.Path, entry.Type)
	return nil
}

// disableCredential is used to disable an existing credential backend
func (c *Core) disableCredential(path string) error {
	c.auth.Lock()
	defer c.auth.Unlock()

	// Ensure we end the path in a slash
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// Ensure the token backend is not affected
	if path == "token/" {
		return fmt.Errorf("token credential backend cannot be disabled")
	}

	// Store the view for this backend
	fullPath := credentialRoutePrefix + path
	view := c.router.MatchingView(fullPath)
	if view == nil {
		return fmt.Errorf("no matching backend")
	}

	// Mark the entry as tainted
	if err := c.taintCredEntry(path); err != nil {
		return err
	}

	// Taint the router path to prevent routing
	if err := c.router.Taint(fullPath); err != nil {
		return err
	}

	// Revoke credentials from this path
	if err := c.expiration.RevokePrefix(fullPath); err != nil {
		return err
	}

	// Unmount the backend
	if err := c.router.Unmount(fullPath); err != nil {
		return err
	}

	// Clear the data in the view
	if view != nil {
		if err := ClearView(view); err != nil {
			return err
		}
	}

	// Remove the mount table entry
	if err := c.removeCredEntry(path); err != nil {
		return err
	}
	c.logger.Printf("[INFO] core: disabled credential backend '%s'", path)
	return nil
}

// removeCredEntry is used to remove an entry in the auth table
func (c *Core) removeCredEntry(path string) error {
	// Taint the entry from the auth table
	newTable := c.auth.Clone()
	newTable.Remove(path)

	// Update the auth table
	if err := c.persistAuth(newTable); err != nil {
		return errors.New("failed to update auth table")
	}
	c.auth = newTable
	return nil
}

// taintCredEntry is used to mark an entry in the auth table as tainted
func (c *Core) taintCredEntry(path string) error {
	// Taint the entry from the auth table
	newTable := c.auth.Clone()
	found := newTable.SetTaint(path, true)

	// Ensure there was a match
	if !found {
		return fmt.Errorf("no matching backend")
	}

	// Update the auth table
	if err := c.persistAuth(newTable); err != nil {
		return errors.New("failed to update auth table")
	}
	c.auth = newTable
	return nil
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
		c.auth = &MountTable{}
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
func (c *Core) persistAuth(table *MountTable) error {
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
	var backend logical.Backend
	var view *BarrierView
	var err error
	for _, entry := range c.auth.Entries {
		// Create a barrier view using the UUID
		view = NewBarrierView(c.barrier, credentialBarrierPrefix+entry.UUID+"/")

		// Initialize the backend
		backend, err = c.newCredentialBackend(entry.Type, c.mountEntrySysView(entry), view, nil)
		if err != nil {
			c.logger.Printf(
				"[ERR] core: failed to create credential entry %#v: %v",
				entry, err)
			return loadAuthFailed
		}

		// Mount the backend
		path := credentialRoutePrefix + entry.Path
		err = c.router.Mount(backend, path, entry, view)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to mount auth entry %#v: %v", entry, err)
			return loadAuthFailed
		}

		// Ensure the path is tainted if set in the mount table
		if entry.Tainted {
			c.router.Taint(path)
		}

		// Check if this is the token store
		if entry.Type == "token" {
			c.tokenStore = backend.(*TokenStore)
		}
	}
	return nil
}

// teardownCredentials is used before we seal the vault to reset the credential
// backends to their unloaded state. This is reversed by loadCredentials.
func (c *Core) teardownCredentials() error {
	c.auth = nil
	c.tokenStore = nil
	return nil
}

// newCredentialBackend is used to create and configure a new credential backend by name
func (c *Core) newCredentialBackend(
	t string, sysView logical.SystemView, view logical.Storage, conf map[string]string) (logical.Backend, error) {
	f, ok := c.credentialBackends[t]
	if !ok {
		return nil, fmt.Errorf("unknown backend type: %s", t)
	}

	config := &logical.BackendConfig{
		View:   view,
		Logger: c.logger,
		Config: conf,
		System: sysView,
	}

	b, err := f(config)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// defaultAuthTable creates a default auth table
func defaultAuthTable() *MountTable {
	table := &MountTable{}
	tokenAuth := &MountEntry{
		Path:        "token/",
		Type:        "token",
		Description: "token based credentials",
		UUID:        uuid.GenerateUUID(),
	}
	table.Entries = append(table.Entries, tokenAuth)
	return table
}
