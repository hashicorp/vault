package vault

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
)

const (
	// coreAuditConfigPath is used to store the audit configuration.
	// Audit configuration is protected within the Vault itself, which means it
	// can only be viewed or modified after an unseal.
	coreAuditConfigPath = "core/audit"

	// auditBarrierPrefix is the prefix to the UUID used in the
	// barrier view for the audit backends.
	auditBarrierPrefix = "audit/"

	// auditTableType is the value we expect to find for the audit table and
	// corresponding entries
	auditTableType = "audit"
)

var (
	// loadAuditFailed if loading audit tables encounters an error
	errLoadAuditFailed = errors.New("failed to setup audit table")
)

// enableAudit is used to enable a new audit backend
func (c *Core) enableAudit(entry *MountEntry) error {
	// Ensure we end the path in a slash
	if !strings.HasSuffix(entry.Path, "/") {
		entry.Path += "/"
	}

	// Ensure there is a name
	if entry.Path == "/" {
		return fmt.Errorf("backend path must be specified")
	}

	// Update the audit table
	c.auditLock.Lock()
	defer c.auditLock.Unlock()

	// Look for matching name
	for _, ent := range c.audit.Entries {
		switch {
		// Existing is sql/mysql/ new is sql/ or
		// existing is sql/ and new is sql/mysql/
		case strings.HasPrefix(ent.Path, entry.Path):
			fallthrough
		case strings.HasPrefix(entry.Path, ent.Path):
			return fmt.Errorf("path already in use")
		}
	}

	// Generate a new UUID and view
	entryUUID, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	entry.UUID = entryUUID
	view := NewBarrierView(c.barrier, auditBarrierPrefix+entry.UUID+"/")

	// Lookup the new backend
	backend, err := c.newAuditBackend(entry.Type, view, entry.Options)
	if err != nil {
		return err
	}

	newTable := c.audit.ShallowClone()
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
	// Ensure we end the path in a slash
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// Remove the entry from the mount table
	c.auditLock.Lock()
	defer c.auditLock.Unlock()

	newTable := c.audit.ShallowClone()
	found := newTable.Remove(path)

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
	auditTable := &MountTable{}

	// Load the existing audit table
	raw, err := c.barrier.Get(coreAuditConfigPath)
	if err != nil {
		c.logger.Printf("[ERR] core: failed to read audit table: %v", err)
		return errLoadAuditFailed
	}

	c.auditLock.Lock()
	defer c.auditLock.Unlock()

	if raw != nil {
		if err := jsonutil.DecodeJSON(raw.Value, auditTable); err != nil {
			c.logger.Printf("[ERR] core: failed to decode audit table: %v", err)
			return errLoadAuditFailed
		}
		c.audit = auditTable
	}

	// Done if we have restored the audit table
	if c.audit != nil {
		needPersist := false

		// Upgrade to typed auth table
		if c.audit.Type == "" {
			c.audit.Type = auditTableType
			needPersist = true
		}

		// Upgrade to table-scoped entries
		for _, entry := range c.audit.Entries {
			if entry.Table == "" {
				entry.Table = c.audit.Type
				needPersist = true
			}
		}

		if needPersist {
			return c.persistAudit(c.audit)
		}

		return nil
	}

	// Create and persist the default audit table
	c.audit = defaultAuditTable()
	if err := c.persistAudit(c.audit); err != nil {
		return errLoadAuditFailed
	}
	return nil
}

// persistAudit is used to persist the audit table after modification
func (c *Core) persistAudit(table *MountTable) error {
	if table.Type != auditTableType {
		c.logger.Printf(
			"[ERR] core: given table to persist has type %s but need type %s",
			table.Type,
			auditTableType)
		return fmt.Errorf("invalid table type given, not persisting")
	}

	for _, entry := range table.Entries {
		if entry.Table != table.Type {
			c.logger.Printf(
				"[ERR] core: entry in audit table with path %s has table value %s but is in table %s, refusing to persist",
				entry.Path,
				entry.Table,
				table.Type)
			return fmt.Errorf("invalid audit entry found, not persisting")
		}
	}

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
	broker := NewAuditBroker(c.logger)

	c.auditLock.Lock()
	defer c.auditLock.Unlock()

	for _, entry := range c.audit.Entries {
		// Create a barrier view using the UUID
		view := NewBarrierView(c.barrier, auditBarrierPrefix+entry.UUID+"/")

		// Initialize the backend
		audit, err := c.newAuditBackend(entry.Type, view, entry.Options)
		if err != nil {
			c.logger.Printf(
				"[ERR] core: failed to create audit entry %s: %v",
				entry.Path, err)
			return errLoadAuditFailed
		}

		// Mount the backend
		broker.Register(entry.Path, audit, view)
	}
	c.auditBroker = broker
	return nil
}

// teardownAudit is used before we seal the vault to reset the audit
// backends to their unloaded state. This is reversed by loadAudits.
func (c *Core) teardownAudits() error {
	c.auditLock.Lock()
	defer c.auditLock.Unlock()

	c.audit = nil
	c.auditBroker = nil
	return nil
}

// newAuditBackend is used to create and configure a new audit backend by name
func (c *Core) newAuditBackend(t string, view logical.Storage, conf map[string]string) (audit.Backend, error) {
	f, ok := c.auditBackends[t]
	if !ok {
		return nil, fmt.Errorf("unknown backend type: %s", t)
	}
	salter, err := salt.NewSalt(view, &salt.Config{
		HMAC:     sha256.New,
		HMACType: "hmac-sha256",
	})
	if err != nil {
		return nil, fmt.Errorf("[ERR] core: unable to generate salt: %v", err)
	}
	return f(&audit.BackendConfig{
		Salt:   salter,
		Config: conf,
	})
}

// defaultAuditTable creates a default audit table
func defaultAuditTable() *MountTable {
	table := &MountTable{
		Type: auditTableType,
	}
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
	logger   *log.Logger
}

// NewAuditBroker creates a new audit broker
func NewAuditBroker(log *log.Logger) *AuditBroker {
	b := &AuditBroker{
		backends: make(map[string]backendEntry),
		logger:   log,
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

// GetHash returns a hash using the salt of the given backend
func (a *AuditBroker) GetHash(name string, input string) (string, error) {
	a.l.RLock()
	defer a.l.RUnlock()
	be, ok := a.backends[name]
	if !ok {
		return "", fmt.Errorf("unknown audit backend %s", name)
	}

	return be.backend.GetHash(input), nil
}

// LogRequest is used to ensure all the audit backends have an opportunity to
// log the given request and that *at least one* succeeds.
func (a *AuditBroker) LogRequest(auth *logical.Auth, req *logical.Request, outerErr error) (reterr error) {
	defer metrics.MeasureSince([]string{"audit", "log_request"}, time.Now())
	a.l.RLock()
	defer a.l.RUnlock()
	defer func() {
		if r := recover(); r != nil {
			a.logger.Printf("[ERR] audit: panic logging: req path: %s", req.Path)
			reterr = fmt.Errorf("panic generating audit log")
		}
	}()

	// Ensure at least one backend logs
	anyLogged := false
	for name, be := range a.backends {
		start := time.Now()
		err := be.backend.LogRequest(auth, req, outerErr)
		metrics.MeasureSince([]string{"audit", name, "log_request"}, start)
		if err != nil {
			a.logger.Printf("[ERR] audit: backend '%s' failed to log request: %v", name, err)
		} else {
			anyLogged = true
		}
	}
	if !anyLogged && len(a.backends) > 0 {
		return fmt.Errorf("no audit backend succeeded in logging the request")
	}
	return nil
}

// LogResponse is used to ensure all the audit backends have an opportunity to
// log the given response and that *at least one* succeeds.
func (a *AuditBroker) LogResponse(auth *logical.Auth, req *logical.Request,
	resp *logical.Response, err error) (reterr error) {
	defer metrics.MeasureSince([]string{"audit", "log_response"}, time.Now())
	a.l.RLock()
	defer a.l.RUnlock()
	defer func() {
		if r := recover(); r != nil {
			a.logger.Printf("[ERR] audit: panic logging: req path: %s: %v", req.Path, r)
			reterr = fmt.Errorf("panic generating audit log")
		}
	}()

	// Ensure at least one backend logs
	anyLogged := false
	for name, be := range a.backends {
		start := time.Now()
		err := be.backend.LogResponse(auth, req, resp, err)
		metrics.MeasureSince([]string{"audit", name, "log_response"}, start)
		if err != nil {
			a.logger.Printf("[ERR] audit: backend '%s' failed to log response: %v", name, err)
		} else {
			anyLogged = true
		}
	}
	if !anyLogged && len(a.backends) > 0 {
		return fmt.Errorf("no audit backend succeeded in logging the response")
	}
	return nil
}
