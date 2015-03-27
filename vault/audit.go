package vault

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/audit"
)

const (
	// coreAuditConfigPath is used to store the audit configuration.
	// Audit configuration is protected within the Vault itself, which means it
	// can only be viewed or modified after an unseal.
	coreAuditConfigPath = "core/audit"
)

var (
	// loadAuditFailed if loading audit tables encounters an error
	loadAuditFailed = errors.New("failed to setup audit table")
)

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
	for _, entry := range c.audit.Entries {
		// Initialize the backend
		_, err := c.newAuditBackend(entry.Type, nil)
		if err != nil {
			c.logger.Printf(
				"[ERR] core: failed to create audit entry %#v: %v",
				entry, err)
			return loadAuditFailed
		}
		// TODO: Do something with backend
	}
	return nil
}

// teardownAudit is used before we seal the vault to reset the audit
// backends to their unloaded state. This is reversed by loadAudits.
func (c *Core) teardownAudits() error {
	c.audit = nil
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
