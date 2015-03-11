package vault

import (
	"encoding/json"
	"errors"
)

const (
	// coreMountConfigPath is used to store the mount configuration.
	// Mounts are protected within the Vault itself, which means they
	// can only be viewed or modified after an unseal.
	coreMountConfigPath = "core/mounts"

	// backendBarrierPrefix is the prefix to the UUID used in the
	// barrier view for the backends.
	backendBarrierPrefix = "logical/"
)

var (
	// loadMountsFailed if loadMounts encounters an error
	loadMountsFailed = errors.New("failed to setup mount table")
)

// MountTable is used to represent the internal mount table
type MountTable struct {
	Entries []*MountEntry `json:"entries"`
}

// MountEntry is used to represent a mount table entry
type MountEntry struct {
	Path        string `json:"path"`        // Mount Path
	Type        string `json:"type"`        // Logical backend Type
	Description string `json:"description"` // User-provided description
	UUID        string `json:"uuid"`        // Barrier view UUID
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
	if err := c.persistMounts(); err != nil {
		return loadMountsFailed
	}
	return nil
}

// persistMounts is used to persist the mount table after modification
func (c *Core) persistMounts() error {
	// Marshal the table
	raw, err := json.Marshal(c.mounts)
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
	var backend LogicalBackend
	var err error
	for _, entry := range c.mounts.Entries {
		// Initialize the backend, special casing for system
		if entry.Type == "system" {
			backend = &SystemBackend{core: c}
		} else {
			backend, err = NewBackend(entry.Type, nil)
			if err != nil {
				c.logger.Printf("[ERR] core: failed to create mount entry %#v: %v", entry, err)
				return loadMountsFailed
			}
		}

		// Create a barrier view using the UUID
		view := NewBarrierView(c.barrier, backendBarrierPrefix+entry.UUID+"/")

		// Mount the backend
		if err := c.router.Mount(backend, entry.Type, entry.Path, view); err != nil {
			c.logger.Printf("[ERR] core: failed to mount entry %#v: %v", entry, err)
			return loadMountsFailed
		}
	}
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
