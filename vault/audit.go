package vault

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/mgutz/logxi/v1"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-multierror"
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
	backend, err := c.newAuditBackend(entry, view, entry.Options)
	if err != nil {
		return err
	}

	newTable := c.audit.shallowClone()
	newTable.Entries = append(newTable.Entries, entry)
	if err := c.persistAudit(newTable); err != nil {
		return errors.New("failed to update audit table")
	}

	c.audit = newTable

	// Register the backend
	c.auditBroker.Register(entry.Path, backend, view)
	if c.logger.IsInfo() {
		c.logger.Info("core: enabled audit backend", "path", entry.Path, "type", entry.Type)
	}
	return nil
}

// disableAudit is used to disable an existing audit backend
func (c *Core) disableAudit(path string) (bool, error) {
	// Ensure we end the path in a slash
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// Remove the entry from the mount table
	c.auditLock.Lock()
	defer c.auditLock.Unlock()

	newTable := c.audit.shallowClone()
	entry := newTable.remove(path)

	// Ensure there was a match
	if entry == nil {
		return false, fmt.Errorf("no matching backend")
	}

	c.removeAuditReloadFunc(entry)

	// Update the audit table
	if err := c.persistAudit(newTable); err != nil {
		return true, errors.New("failed to update audit table")
	}

	c.audit = newTable

	// Unmount the backend
	c.auditBroker.Deregister(path)
	if c.logger.IsInfo() {
		c.logger.Info("core: disabled audit backend", "path", path)
	}
	return true, nil
}

// loadAudits is invoked as part of postUnseal to load the audit table
func (c *Core) loadAudits() error {
	auditTable := &MountTable{}

	// Load the existing audit table
	raw, err := c.barrier.Get(coreAuditConfigPath)
	if err != nil {
		c.logger.Error("core: failed to read audit table", "error", err)
		return errLoadAuditFailed
	}

	c.auditLock.Lock()
	defer c.auditLock.Unlock()

	if raw != nil {
		if err := jsonutil.DecodeJSON(raw.Value, auditTable); err != nil {
			c.logger.Error("core: failed to decode audit table", "error", err)
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
		c.logger.Error("core: given table to persist has wrong type", "actual_type", table.Type, "expected_type", auditTableType)
		return fmt.Errorf("invalid table type given, not persisting")
	}

	for _, entry := range table.Entries {
		if entry.Table != table.Type {
			c.logger.Error("core: given entry to persist in audit table has wrong table value", "path", entry.Path, "entry_table_type", entry.Table, "actual_type", table.Type)
			return fmt.Errorf("invalid audit entry found, not persisting")
		}
	}

	// Marshal the table
	raw, err := json.Marshal(table)
	if err != nil {
		c.logger.Error("core: failed to encode audit table", "error", err)
		return err
	}

	// Create an entry
	entry := &Entry{
		Key:   coreAuditConfigPath,
		Value: raw,
	}

	// Write to the physical backend
	if err := c.barrier.Put(entry); err != nil {
		c.logger.Error("core: failed to persist audit table", "error", err)
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

	var successCount int

	for _, entry := range c.audit.Entries {
		// Create a barrier view using the UUID
		view := NewBarrierView(c.barrier, auditBarrierPrefix+entry.UUID+"/")

		// Initialize the backend
		audit, err := c.newAuditBackend(entry, view, entry.Options)
		if err != nil {
			c.logger.Error("core: failed to create audit entry", "path", entry.Path, "error", err)
			continue
		}

		// Mount the backend
		broker.Register(entry.Path, audit, view)

		successCount += 1
	}

	if len(c.audit.Entries) > 0 && successCount == 0 {
		return errLoadAuditFailed
	}

	c.auditBroker = broker
	return nil
}

// teardownAudit is used before we seal the vault to reset the audit
// backends to their unloaded state. This is reversed by loadAudits.
func (c *Core) teardownAudits() error {
	c.auditLock.Lock()
	defer c.auditLock.Unlock()

	if c.audit != nil {
		for _, entry := range c.audit.Entries {
			c.removeAuditReloadFunc(entry)
		}
	}

	c.audit = nil
	c.auditBroker = nil
	return nil
}

// removeAuditReloadFunc removes the reload func from the working set. The
// audit lock needs to be held before calling this.
func (c *Core) removeAuditReloadFunc(entry *MountEntry) {
	switch entry.Type {
	case "file":
		key := "audit_file|" + entry.Path
		c.reloadFuncsLock.Lock()

		if c.logger.IsDebug() {
			c.logger.Debug("audit: removing reload function", "path", entry.Path)
		}

		delete(c.reloadFuncs, key)

		c.reloadFuncsLock.Unlock()
	}
}

// newAuditBackend is used to create and configure a new audit backend by name
func (c *Core) newAuditBackend(entry *MountEntry, view logical.Storage, conf map[string]string) (audit.Backend, error) {
	f, ok := c.auditBackends[entry.Type]
	if !ok {
		return nil, fmt.Errorf("unknown backend type: %s", entry.Type)
	}
	salter, err := salt.NewSalt(view, &salt.Config{
		HMAC:     sha256.New,
		HMACType: "hmac-sha256",
	})
	if err != nil {
		return nil, fmt.Errorf("core: unable to generate salt: %v", err)
	}

	be, err := f(&audit.BackendConfig{
		Salt:   salter,
		Config: conf,
	})
	if err != nil {
		return nil, err
	}

	switch entry.Type {
	case "file":
		key := "audit_file|" + entry.Path

		c.reloadFuncsLock.Lock()

		if c.logger.IsDebug() {
			c.logger.Debug("audit: adding reload function", "path", entry.Path)
		}

		c.reloadFuncs[key] = append(c.reloadFuncs[key], func(map[string]string) error {
			if c.logger.IsInfo() {
				c.logger.Info("audit: reloading file audit backend", "path", entry.Path)
			}
			return be.Reload()
		})

		c.reloadFuncsLock.Unlock()
	}

	return be, err
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
	sync.RWMutex
	backends map[string]backendEntry
	logger   log.Logger
}

// NewAuditBroker creates a new audit broker
func NewAuditBroker(log log.Logger) *AuditBroker {
	b := &AuditBroker{
		backends: make(map[string]backendEntry),
		logger:   log,
	}
	return b
}

// Register is used to add new audit backend to the broker
func (a *AuditBroker) Register(name string, b audit.Backend, v *BarrierView) {
	a.Lock()
	defer a.Unlock()
	a.backends[name] = backendEntry{
		backend: b,
		view:    v,
	}
}

// Deregister is used to remove an audit backend from the broker
func (a *AuditBroker) Deregister(name string) {
	a.Lock()
	defer a.Unlock()
	delete(a.backends, name)
}

// IsRegistered is used to check if a given audit backend is registered
func (a *AuditBroker) IsRegistered(name string) bool {
	a.RLock()
	defer a.RUnlock()
	_, ok := a.backends[name]
	return ok
}

// GetHash returns a hash using the salt of the given backend
func (a *AuditBroker) GetHash(name string, input string) (string, error) {
	a.RLock()
	defer a.RUnlock()
	be, ok := a.backends[name]
	if !ok {
		return "", fmt.Errorf("unknown audit backend %s", name)
	}

	return be.backend.GetHash(input), nil
}

// LogRequest is used to ensure all the audit backends have an opportunity to
// log the given request and that *at least one* succeeds.
func (a *AuditBroker) LogRequest(auth *logical.Auth, req *logical.Request, outerErr error) (retErr error) {
	defer metrics.MeasureSince([]string{"audit", "log_request"}, time.Now())
	a.RLock()
	defer a.RUnlock()
	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("audit: panic during logging", "request_path", req.Path, "error", r)
			retErr = multierror.Append(retErr, fmt.Errorf("panic generating audit log"))
		}
	}()

	// All logged requests must have an identifier
	//if req.ID == "" {
	//	a.logger.Error("audit: missing identifier in request object", "request_path", req.Path)
	//	retErr = multierror.Append(retErr, fmt.Errorf("missing identifier in request object: %s", req.Path))
	//	return
	//}

	// Ensure at least one backend logs
	anyLogged := false
	for name, be := range a.backends {
		start := time.Now()
		err := be.backend.LogRequest(auth, req, outerErr)
		metrics.MeasureSince([]string{"audit", name, "log_request"}, start)
		if err != nil {
			a.logger.Error("audit: backend failed to log request", "backend", name, "error", err)
		} else {
			anyLogged = true
		}
	}
	if !anyLogged && len(a.backends) > 0 {
		retErr = multierror.Append(retErr, fmt.Errorf("no audit backend succeeded in logging the request"))
		return
	}
	return nil
}

// LogResponse is used to ensure all the audit backends have an opportunity to
// log the given response and that *at least one* succeeds.
func (a *AuditBroker) LogResponse(auth *logical.Auth, req *logical.Request,
	resp *logical.Response, err error) (reterr error) {
	defer metrics.MeasureSince([]string{"audit", "log_response"}, time.Now())
	a.RLock()
	defer a.RUnlock()
	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("audit: panic during logging", "request_path", req.Path, "error", r)
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
			a.logger.Error("audit: backend failed to log response", "backend", name, "error", err)
		} else {
			anyLogged = true
		}
	}
	if !anyLogged && len(a.backends) > 0 {
		return fmt.Errorf("no audit backend succeeded in logging the response")
	}
	return nil
}
