package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/vault/audit"
)

const (
	// coreAuditConfigPath is used to store the audit configuration.
	// Audit configuration is protected within the Vault itself, which means it
	// can only be viewed or modified after an unseal.
	coreAuditConfigPath = "core/audit"

	// auditBarrierPrefix is the prefix to the UUID used in the
	// barrier view for the audit backends.
	auditBarrierPrefix = "audit/"
)

var (
	// loadAuditFailed if loading audit tables encounters an error
	loadAuditFailed = errors.New("failed to setup audit table")
)

// enableAudit is used to enable a new audit backend
func (c *Core) enableAudit(entry *MountEntry) error {
	c.audit.Lock()
	defer c.audit.Unlock()

	// Ensure there is a name
	if entry.Path == "" {
		return fmt.Errorf("backend path must be specified")
	}
	if strings.Contains(entry.Path, "/") {
		return fmt.Errorf("backend path cannot have a forward slash")
	}

	// Look for matching name
	for _, ent := range c.audit.Entries {
		if ent.Path == entry.Path {
			return fmt.Errorf("path already in use")
		}
	}

	// Lookup the new backend
	backend, err := c.newAuditBackend(entry.Type, entry.Options)
	if err != nil {
		return err
	}

	// Generate a new UUID and view
	entry.UUID = generateUUID()
	view := NewBarrierView(c.barrier, auditBarrierPrefix+entry.UUID+"/")

	// Update the audit table
	newTable := c.audit.Clone()
	newTable.Entries = append(newTable.Entries, entry)
	if err := c.persistAudit(newTable); err != nil {
		return errors.New("failed to update audit table")
	}
	c.audit = newTable

	// Register the backend
	c.auditBroker.Register(entry.Path, backend, view)
	c.logger.Printf("[INFO] core: enabled audit backend '%s' type: %s",
		entry.Path, entry.Type)
	return nil
}

// disableAudit is used to disable an existing audit backend
func (c *Core) disableAudit(path string) error {
	c.audit.Lock()
	defer c.audit.Unlock()

	// Remove the entry from the mount table
	found := false
	newTable := c.audit.Clone()
	n := len(newTable.Entries)
	for i := 0; i < n; i++ {
		if newTable.Entries[i].Path == path {
			newTable.Entries[i], newTable.Entries[n-1] = newTable.Entries[n-1], nil
			newTable.Entries = newTable.Entries[:n-1]
			found = true
			break
		}
	}

	// Ensure there was a match
	if !found {
		return fmt.Errorf("no matching backend")
	}

	// Update the audit table
	if err := c.persistAudit(newTable); err != nil {
		return errors.New("failed to update audit table")
	}
	c.audit = newTable

	// Unmount the backend
	c.auditBroker.Deregister(path)
	c.logger.Printf("[INFO] core: disabled audit backend '%s'", path)
	return nil
}

// loadAudits is invoked as part of postUnseal to load the audit table
func (c *Core) loadAudits() error {
	// Load the existing audit table
	raw, err := c.barrier.Get(coreAuditConfigPath)
	if err != nil {
		c.logger.Printf("[ERR] core: failed to read audit table: %v", err)
		return loadAuditFailed
	}
	if raw != nil {
		c.audit = &MountTable{}
		if err := json.Unmarshal(raw.Value, c.audit); err != nil {
			c.logger.Printf("[ERR] core: failed to decode audit table: %v", err)
			return loadAuditFailed
		}
	}

	// Done if we have restored the audit table
	if c.audit != nil {
		return nil
	}

	// Create and persist the default audit table
	c.audit = defaultAuditTable()
	if err := c.persistAudit(c.audit); err != nil {
		return loadAuditFailed
	}
	return nil
}

// persistAudit is used to persist the audit table after modification
func (c *Core) persistAudit(table *MountTable) error {
	// Marshal the table
	raw, err := json.Marshal(table)
	if err != nil {
		c.logger.Printf("[ERR] core: failed to encode audit table: %v", err)
		return err
	}

	// Create an entry
	entry := &Entry{
		Key:   coreAuditConfigPath,
		Value: raw,
	}

	// Write to the physical backend
	if err := c.barrier.Put(entry); err != nil {
		c.logger.Printf("[ERR] core: failed to persist audit table: %v", err)
		return err
	}
	return nil
}

// setupAudit is invoked after we've loaded the audit able to
// initialize the audit backends
func (c *Core) setupAudits() error {
	broker := NewAuditBroker()
	for _, entry := range c.audit.Entries {
		// Initialize the backend
		audit, err := c.newAuditBackend(entry.Type, entry.Options)
		if err != nil {
			c.logger.Printf(
				"[ERR] core: failed to create audit entry %#v: %v",
				entry, err)
			return loadAuditFailed
		}

		// Create a barrier view using the UUID
		view := NewBarrierView(c.barrier, auditBarrierPrefix+entry.UUID+"/")

		// Mount the backend
		broker.Register(entry.Path, audit, view)
	}
	c.auditBroker = broker
	return nil
}

// teardownAudit is used before we seal the vault to reset the audit
// backends to their unloaded state. This is reversed by loadAudits.
func (c *Core) teardownAudits() error {
	c.audit = nil
	c.auditBroker = nil
	return nil
}

// newAuditBackend is used to create and configure a new audit backend by name
func (c *Core) newAuditBackend(t string, conf map[string]string) (audit.Backend, error) {
	f, ok := c.auditBackends[t]
	if !ok {
		return nil, fmt.Errorf("unknown backend type: %s", t)
	}
	return f(conf)
}

// defaultAuditTable creates a default audit table
func defaultAuditTable() *MountTable {
	table := &MountTable{}
	return table
}

type backendEntry struct {
	backend audit.Backend
	view    *BarrierView
}

// AuditBroker is used to provide a single ingest interface to auditable
// events given that multiple backends may be configured.
type AuditBroker struct {
	l        sync.RWMutex
	backends map[string]backendEntry
}

// NewAuditBroker creates a new audit broker
func NewAuditBroker() *AuditBroker {
	b := &AuditBroker{
		backends: make(map[string]backendEntry),
	}
	return b
}

// Register is used to add new audit backend to the broker
func (a *AuditBroker) Register(name string, b audit.Backend, v *BarrierView) {
	a.l.Lock()
	defer a.l.Unlock()
	a.backends[name] = backendEntry{
		backend: b,
		view:    v,
	}
}

// Deregister is used to remove an audit backend from the broker
func (a *AuditBroker) Deregister(name string) {
	a.l.Lock()
	defer a.l.Unlock()
	delete(a.backends, name)
}

// IsRegistered is used to check if a given audit backend is registered
func (a *AuditBroker) IsRegistered(name string) bool {
	a.l.RLock()
	defer a.l.RUnlock()
	_, ok := a.backends[name]
	return ok
}
