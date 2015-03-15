package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/logical"
)

// TEMPORARY!

// BuiltinBackends contains all of the available backends
var BuiltinBackends = map[string]logical.Factory{
	"generic": PassthroughBackendFactory,
}

// NewBackend returns a new logical Backend with the given type and configuration.
// The backend is looked up in the BuiltinBackends variable.
func NewBackend(t string, conf map[string]string) (logical.Backend, error) {
	f, ok := BuiltinBackends[t]
	if !ok {
		return nil, fmt.Errorf("unknown logical backend type: %s", t)
	}
	return f(conf)
}

// TEMPORARY!

const (
	// coreMountConfigPath is used to store the mount configuration.
	// Mounts are protected within the Vault itself, which means they
	// can only be viewed or modified after an unseal.
	coreMountConfigPath = "core/mounts"

	// backendBarrierPrefix is the prefix to the UUID used in the
	// barrier view for the backends.
	backendBarrierPrefix = "logical/"

	// systemBarrierPrefix is sthe prefix used for the
	// system logical backend.
	systemBarrierPrefix = "sys/"
)

var (
	// loadMountsFailed if loadMounts encounters an error
	loadMountsFailed = errors.New("failed to setup mount table")

	// protectedMounts cannot be remounted
	protectedMounts = []string{
		"sys/",
	}
)

// MountTable is used to represent the internal mount table
type MountTable struct {
	Entries []*MountEntry `json:"entries"`
}

// Returns a deep copy of the mount table
func (t *MountTable) Clone() *MountTable {
	mt := &MountTable{
		Entries: make([]*MountEntry, len(t.Entries)),
	}
	for i, e := range t.Entries {
		mt.Entries[i] = e.Clone()
	}
	return mt
}

// MountEntry is used to represent a mount table entry
type MountEntry struct {
	Path        string `json:"path"`        // Mount Path
	Type        string `json:"type"`        // Logical backend Type
	Description string `json:"description"` // User-provided description
	UUID        string `json:"uuid"`        // Barrier view UUID
}

// Returns a deep copy of the mount entry
func (e *MountEntry) Clone() *MountEntry {
	return &MountEntry{
		Path:        e.Path,
		Type:        e.Type,
		Description: e.Description,
		UUID:        e.UUID,
	}
}

// Mount is used to mount a new backend to the mount table.
func (c *Core) mount(me *MountEntry) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// Ensure we end the path in a slash
	if !strings.HasSuffix(me.Path, "/") {
		me.Path += "/"
	}

	// Verify there is no conflicting mount
	if match := c.router.MatchingMount(me.Path); match != "" {
		return fmt.Errorf("existing mount at '%s'", match)
	}

	// Lookup the new backend
	backend, err := NewBackend(me.Type, nil)
	if err != nil {
		return err
	}

	// Generate a new UUID and view
	me.UUID = generateUUID()
	view := NewBarrierView(c.barrier, backendBarrierPrefix+me.UUID+"/")

	// Update the mount table
	newTable := c.mounts.Clone()
	newTable.Entries = append(newTable.Entries, me)
	if err := c.persistMounts(newTable); err != nil {
		return errors.New("failed to update mount table")
	}
	c.mounts = newTable

	// Mount the backend
	if err := c.router.Mount(backend, me.Type, me.Path, view); err != nil {
		return err
	}
	c.logger.Printf("[INFO] core: mounted '%s' type: %s", me.Path, me.Type)
	return nil
}

// Unmount is used to unmount a path.
//
// TODO: document what happens to all secrets currently out for this path.
func (c *Core) unmount(path string) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// Ensure we end the path in a slash
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// Verify exact match of the route
	match := c.router.MatchingMount(path)
	if match == "" || path != match {
		return fmt.Errorf("no matching mount")
	}

	// Remove the entry from the mount table
	newTable := c.mounts.Clone()
	n := len(newTable.Entries)
	for i := 0; i < n; i++ {
		if newTable.Entries[i].Path == path {
			newTable.Entries[i], newTable.Entries[n-1] = newTable.Entries[n-1], nil
			newTable.Entries = newTable.Entries[:n-1]
			break
		}
	}

	// Update the mount table
	if err := c.persistMounts(newTable); err != nil {
		return errors.New("failed to update mount table")
	}
	c.mounts = newTable

	// Unmount the backend
	if err := c.router.Unmount(path); err != nil {
		return err
	}

	// TODO: Delete data in view?
	// TODO: Handle revocation?
	c.logger.Printf("[INFO] core: unmounted '%s'", path)
	return nil
}

// Remount is used to remount a path at a new mount point.
func (c *Core) remount(src, dst string) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// Ensure we end the path in a slash
	if !strings.HasSuffix(src, "/") {
		src += "/"
	}
	if !strings.HasSuffix(dst, "/") {
		dst += "/"
	}

	// Prevent sys/ from being remounted
	for _, p := range protectedMounts {
		if src == p {
			return fmt.Errorf("cannot remount '%s'", p)
		}
	}

	// Verify exact match of the route
	match := c.router.MatchingMount(src)
	if match == "" || src != match {
		return fmt.Errorf("no matching mount at '%s'", src)
	}

	// Verify there is no conflicting mount
	if match := c.router.MatchingMount(dst); match != "" {
		return fmt.Errorf("existing mount at '%s'", match)
	}

	// Update the entry in the mount table
	newTable := c.mounts.Clone()
	for _, ent := range newTable.Entries {
		if ent.Path == src {
			ent.Path = dst
			break
		}
	}

	// Update the mount table
	if err := c.persistMounts(newTable); err != nil {
		return errors.New("failed to update mount table")
	}
	c.mounts = newTable

	// Remount the backend
	if err := c.router.Remount(src, dst); err != nil {
		return err
	}
	c.logger.Printf("[INFO] core: remounted '%s' to '%s'", src, dst)
	return nil
}

// loadMounts is invoked as part of postUnseal to load the mount table
func (c *Core) loadMounts() error {
	// Load the existing mount table
	raw, err := c.barrier.Get(coreMountConfigPath)
	if err != nil {
		c.logger.Printf("[ERR] core: failed to read mount table: %v", err)
		return loadMountsFailed
	}
	if raw != nil {
		c.mounts = &MountTable{}
		if err := json.Unmarshal(raw.Value, c.mounts); err != nil {
			c.logger.Printf("[ERR] core: failed to decode mount table: %v", err)
			return loadMountsFailed
		}
	}

	// Done if we have restored the mount table
	if c.mounts != nil {
		return nil
	}

	// Create and persist the default mount table
	c.mounts = defaultMountTable()
	if err := c.persistMounts(c.mounts); err != nil {
		return loadMountsFailed
	}
	return nil
}

// persistMounts is used to persist the mount table after modification
func (c *Core) persistMounts(table *MountTable) error {
	// Marshal the table
	raw, err := json.Marshal(table)
	if err != nil {
		c.logger.Printf("[ERR] core: failed to encode mount table: %v", err)
		return err
	}

	// Create an entry
	entry := &Entry{
		Key:   coreMountConfigPath,
		Value: raw,
	}

	// Write to the physical backend
	if err := c.barrier.Put(entry); err != nil {
		c.logger.Printf("[ERR] core: failed to persist mount table: %v", err)
		return err
	}
	return nil
}

// setupMounts is invoked after we've loaded the mount table to
// initialize the logical backends and setup the router
func (c *Core) setupMounts() error {
	var backend logical.Backend
	var view *BarrierView
	var err error
	for _, entry := range c.mounts.Entries {
		// Initialize the backend, special casing for system
		if entry.Type == "system" {
			backend = &SystemBackend2{Core: c}
			view = NewBarrierView(c.barrier, systemBarrierPrefix+entry.UUID+"/")
			c.systemView = view

		} else {
			backend, err = NewBackend(entry.Type, nil)
			if err != nil {
				c.logger.Printf("[ERR] core: failed to create mount entry %#v: %v", entry, err)
				return loadMountsFailed
			}

			// Create a barrier view using the UUID
			view = NewBarrierView(c.barrier, backendBarrierPrefix+entry.UUID+"/")
		}

		// Mount the backend
		if err := c.router.Mount(backend, entry.Type, entry.Path, view); err != nil {
			c.logger.Printf("[ERR] core: failed to mount entry %#v: %v", entry, err)
			return loadMountsFailed
		}
	}
	return nil
}

// unloadMounts is used before we seal the vault to reset the mounts to
// their unloaded state. This is reversed by load and setup mounts.
func (c *Core) unloadMounts() error {
	c.mounts = nil
	c.router = NewRouter()
	c.systemView = nil
	return nil
}

// defaultMountTable creates a default mount table
func defaultMountTable() *MountTable {
	table := &MountTable{}
	genericMount := &MountEntry{
		Path:        "secret/",
		Type:        "generic",
		Description: "generic secret storage",
		UUID:        generateUUID(),
	}
	sysMount := &MountEntry{
		Path:        "sys/",
		Type:        "system",
		Description: "system endpoints used for control, policy and debugging",
		UUID:        generateUUID(),
	}
	table.Entries = append(table.Entries, genericMount)
	table.Entries = append(table.Entries, sysMount)
	return table
}
