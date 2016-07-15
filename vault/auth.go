package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/jsonutil"
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

	// credentialTableType is the value we expect to find for the credential
	// table and corresponding entries
	credentialTableType = "auth"
)

var (
	// errLoadAuthFailed if loadCredentials encounters an error
	errLoadAuthFailed = errors.New("failed to setup auth table")
)

// enableCredential is used to enable a new credential backend
func (c *Core) enableCredential(entry *MountEntry) error {
	// Ensure we end the path in a slash
	if !strings.HasSuffix(entry.Path, "/") {
		entry.Path += "/"
	}

	// Ensure there is a name
	if entry.Path == "/" {
		return fmt.Errorf("backend path must be specified")
	}

	c.authLock.Lock()
	defer c.authLock.Unlock()

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
	entryUUID, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	entry.UUID = entryUUID
	view := NewBarrierView(c.barrier, credentialBarrierPrefix+entry.UUID+"/")

	// Create the new backend
	backend, err := c.newCredentialBackend(entry.Type, c.mountEntrySysView(entry), view, nil)
	if err != nil {
		return err
	}

	// Update the auth table
	newTable := c.auth.ShallowClone()
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
	view := c.router.MatchingStorageView(fullPath)
	if view == nil {
		return fmt.Errorf("no matching backend")
	}

	c.authLock.Lock()
	defer c.authLock.Unlock()

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
	newTable := c.auth.ShallowClone()
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
	// We do this on the original since setting the taint operates
	// on the entries which a shallow clone shares anyways
	found := c.auth.SetTaint(path, true)

	// Ensure there was a match
	if !found {
		return fmt.Errorf("no matching backend")
	}

	// Update the auth table
	if err := c.persistAuth(c.auth); err != nil {
		return errors.New("failed to update auth table")
	}

	return nil
}

// loadCredentials is invoked as part of postUnseal to load the auth table
func (c *Core) loadCredentials() error {
	authTable := &MountTable{}
	// Load the existing mount table
	raw, err := c.barrier.Get(coreAuthConfigPath)
	if err != nil {
		c.logger.Printf("[ERR] core: failed to read auth table: %v", err)
		return errLoadAuthFailed
	}

	c.authLock.Lock()
	defer c.authLock.Unlock()

	if raw != nil {
		if err := jsonutil.DecodeJSON(raw.Value, authTable); err != nil {
			c.logger.Printf("[ERR] core: failed to decode auth table: %v", err)
			return errLoadAuthFailed
		}
		c.auth = authTable
	}

	// Done if we have restored the auth table
	if c.auth != nil {
		needPersist := false

		// Upgrade to typed auth table
		if c.auth.Type == "" {
			c.auth.Type = credentialTableType
			needPersist = true
		}

		// Upgrade to table-scoped entries
		for _, entry := range c.auth.Entries {
			// The auth backend "aws-ec2" was named "aws" in the master.
			// This is to support upgrade procedure from "aws" to "aws-ec2".
			if entry.Type == "aws" {
				entry.Type = "aws-ec2"
				needPersist = true
			}
			if entry.Table == "" {
				entry.Table = c.auth.Type
				needPersist = true
			}
		}

		if needPersist {
			return c.persistAuth(c.auth)
		}

		return nil
	}

	// Create and persist the default auth table
	c.auth = defaultAuthTable()
	if err := c.persistAuth(c.auth); err != nil {
		c.logger.Printf("[ERR] core: failed to persist auth table: %v", err)
		return errLoadAuthFailed
	}
	return nil
}

// persistAuth is used to persist the auth table after modification
func (c *Core) persistAuth(table *MountTable) error {
	if table.Type != credentialTableType {
		c.logger.Printf(
			"[ERR] core: given table to persist has type %s but need type %s",
			table.Type,
			credentialTableType)
		return fmt.Errorf("invalid table type given, not persisting")
	}

	for _, entry := range table.Entries {
		if entry.Table != table.Type {
			c.logger.Printf(
				"[ERR] core: entry in auth table with path %s has table value %s but is in table %s, refusing to persist",
				entry.Path,
				entry.Table,
				table.Type)
			return fmt.Errorf("invalid auth entry found, not persisting")
		}
	}

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
	var persistNeeded bool

	c.authLock.Lock()
	defer c.authLock.Unlock()

	for _, entry := range c.auth.Entries {
		// Work around some problematic code that existed in master for a while
		if strings.HasPrefix(entry.Path, credentialRoutePrefix) {
			entry.Path = strings.TrimPrefix(entry.Path, credentialRoutePrefix)
			persistNeeded = true
		}

		// Create a barrier view using the UUID
		view = NewBarrierView(c.barrier, credentialBarrierPrefix+entry.UUID+"/")

		// Initialize the backend
		backend, err = c.newCredentialBackend(entry.Type, c.mountEntrySysView(entry), view, nil)
		if err != nil {
			c.logger.Printf(
				"[ERR] core: failed to create credential entry %s: %v",
				entry.Path, err)
			return errLoadAuthFailed
		}

		// Mount the backend
		path := credentialRoutePrefix + entry.Path
		err = c.router.Mount(backend, path, entry, view)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to mount auth entry %s: %v", entry.Path, err)
			return errLoadAuthFailed
		}

		// Ensure the path is tainted if set in the mount table
		if entry.Tainted {
			c.router.Taint(path)
		}

		// Check if this is the token store
		if entry.Type == "token" {
			c.tokenStore = backend.(*TokenStore)

			// this is loaded *after* the normal mounts, including cubbyhole
			c.router.tokenStoreSalt = c.tokenStore.salt
			c.tokenStore.cubbyholeBackend = c.router.MatchingBackend("cubbyhole/").(*CubbyholeBackend)
		}
	}

	if persistNeeded {
		return c.persistAuth(c.auth)
	}

	return nil
}

// teardownCredentials is used before we seal the vault to reset the credential
// backends to their unloaded state. This is reversed by loadCredentials.
func (c *Core) teardownCredentials() error {
	c.authLock.Lock()
	defer c.authLock.Unlock()

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
		StorageView: view,
		Logger:      c.logger,
		Config:      conf,
		System:      sysView,
	}

	b, err := f(config)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// defaultAuthTable creates a default auth table
func defaultAuthTable() *MountTable {
	table := &MountTable{
		Type: credentialTableType,
	}
	tokenUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(fmt.Sprintf("could not generate UUID for default auth table token entry: %v", err))
	}
	tokenAuth := &MountEntry{
		Table:       credentialTableType,
		Path:        "token/",
		Type:        "token",
		Description: "token based credentials",
		UUID:        tokenUUID,
	}
	table.Entries = append(table.Entries, tokenAuth)
	return table
}
