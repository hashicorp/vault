package vault

import (
	"crypto/sha256"
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

	// coreLocalAuditConfigPath is used to store audit information for local
	// (non-replicated) mounts
	coreLocalAuditConfigPath = "core/local-audit"

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
	if entry.UUID == "" {
		entryUUID, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}
		entry.UUID = entryUUID
	}
	if entry.Accessor == "" {
		accessor, err := c.generateMountAccessor("audit_" + entry.Type)
		if err != nil {
			return err
		}
		entry.Accessor = accessor
	}
	viewPath := auditBarrierPrefix + entry.UUID + "/"
	view := NewBarrierView(c.barrier, viewPath)

	// Lookup the new backend
	backend, err := c.newAuditBackend(entry, view, entry.Options)
	if err != nil {
		return err
	}
	if backend == nil {
		return fmt.Errorf("nil audit backend of type %q returned from factory", entry.Type)
	}

	newTable := c.audit.shallowClone()
	newTable.Entries = append(newTable.Entries, entry)
	if err := c.persistAudit(newTable, entry.Local); err != nil {
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

	// When unmounting all entries the JSON code will load back up from storage
	// as a nil slice, which kills tests...just set it nil explicitly
	if len(newTable.Entries) == 0 {
		newTable.Entries = nil
	}

	// Update the audit table
	if err := c.persistAudit(newTable, entry.Local); err != nil {
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
	localAuditTable := &MountTable{}

	// Load the existing audit table
	raw, err := c.barrier.Get(coreAuditConfigPath)
	if err != nil {
		c.logger.Error("core: failed to read audit table", "error", err)
		return errLoadAuditFailed
	}
	rawLocal, err := c.barrier.Get(coreLocalAuditConfigPath)
	if err != nil {
		c.logger.Error("core: failed to read local audit table", "error", err)
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

	var needPersist bool
	if c.audit == nil {
		c.audit = defaultAuditTable()
		needPersist = true
	}

	if rawLocal != nil {
		if err := jsonutil.DecodeJSON(rawLocal.Value, localAuditTable); err != nil {
			c.logger.Error("core: failed to decode local audit table", "error", err)
			return errLoadAuditFailed
		}
		if localAuditTable != nil && len(localAuditTable.Entries) > 0 {
			c.audit.Entries = append(c.audit.Entries, localAuditTable.Entries...)
		}
	}

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
		if entry.Accessor == "" {
			accessor, err := c.generateMountAccessor("audit_" + entry.Type)
			if err != nil {
				return err
			}
			entry.Accessor = accessor
			needPersist = true
		}
	}

	if !needPersist {
		return nil
	}

	if err := c.persistAudit(c.audit, false); err != nil {
		return errLoadAuditFailed
	}
	return nil
}

// persistAudit is used to persist the audit table after modification
func (c *Core) persistAudit(table *MountTable, localOnly bool) error {
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

	nonLocalAudit := &MountTable{
		Type: auditTableType,
	}

	localAudit := &MountTable{
		Type: auditTableType,
	}

	for _, entry := range table.Entries {
		if entry.Local {
			localAudit.Entries = append(localAudit.Entries, entry)
		} else {
			nonLocalAudit.Entries = append(nonLocalAudit.Entries, entry)
		}
	}

	if !localOnly {
		// Marshal the table
		compressedBytes, err := jsonutil.EncodeJSONAndCompress(nonLocalAudit, nil)
		if err != nil {
			c.logger.Error("core: failed to encode and/or compress audit table", "error", err)
			return err
		}

		// Create an entry
		entry := &Entry{
			Key:   coreAuditConfigPath,
			Value: compressedBytes,
		}

		// Write to the physical backend
		if err := c.barrier.Put(entry); err != nil {
			c.logger.Error("core: failed to persist audit table", "error", err)
			return err
		}
	}

	// Repeat with local audit
	compressedBytes, err := jsonutil.EncodeJSONAndCompress(localAudit, nil)
	if err != nil {
		c.logger.Error("core: failed to encode and/or compress local audit table", "error", err)
		return err
	}

	entry := &Entry{
		Key:   coreLocalAuditConfigPath,
		Value: compressedBytes,
	}

	if err := c.barrier.Put(entry); err != nil {
		c.logger.Error("core: failed to persist local audit table", "error", err)
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
		viewPath := auditBarrierPrefix + entry.UUID + "/"
		view := NewBarrierView(c.barrier, viewPath)

		// Initialize the backend
		backend, err := c.newAuditBackend(entry, view, entry.Options)
		if err != nil {
			c.logger.Error("core: failed to create audit entry", "path", entry.Path, "error", err)
			continue
		}
		if backend == nil {
			c.logger.Error("core: created audit entry was nil", "path", entry.Path, "type", entry.Type)
			continue
		}

		// Mount the backend
		broker.Register(entry.Path, backend, view)

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
	saltConfig := &salt.Config{
		HMAC:     sha256.New,
		HMACType: "hmac-sha256",
		Location: salt.DefaultLocation,
	}

	be, err := f(&audit.BackendConfig{
		SaltView:   view,
		SaltConfig: saltConfig,
		Config:     conf,
	})
	if err != nil {
		return nil, err
	}
	if be == nil {
		return nil, fmt.Errorf("nil backend returned from %q factory function", entry.Type)
	}

	switch entry.Type {
	case "file":
		key := "audit_file|" + entry.Path

		c.reloadFuncsLock.Lock()

		if c.logger.IsDebug() {
			c.logger.Debug("audit: adding reload function", "path", entry.Path)
			if entry.Options != nil {
				c.logger.Debug("audit: file backend options", "path", entry.Path, "file_path", entry.Options["file_path"])
			}
		}

		c.reloadFuncs[key] = append(c.reloadFuncs[key], func(map[string]interface{}) error {
			if c.logger.IsInfo() {
				c.logger.Info("audit: reloading file audit backend", "path", entry.Path)
			}
			return be.Reload()
		})

		c.reloadFuncsLock.Unlock()
	case "socket":
		if c.logger.IsDebug() {
			if entry.Options != nil {
				c.logger.Debug("audit: socket backend options", "path", entry.Path, "address", entry.Options["address"], "socket type", entry.Options["socket_type"])
			}
		}
	case "syslog":
		if c.logger.IsDebug() {
			if entry.Options != nil {
				c.logger.Debug("audit: syslog backend options", "path", entry.Path, "facility", entry.Options["facility"], "tag", entry.Options["tag"])
			}
		}
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

	return be.backend.GetHash(input)
}

// LogRequest is used to ensure all the audit backends have an opportunity to
// log the given request and that *at least one* succeeds.
func (a *AuditBroker) LogRequest(auth *logical.Auth, req *logical.Request, headersConfig *AuditedHeadersConfig, outerErr error) (ret error) {
	defer metrics.MeasureSince([]string{"audit", "log_request"}, time.Now())
	a.RLock()
	defer a.RUnlock()

	var retErr *multierror.Error

	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("audit: panic during logging", "request_path", req.Path, "error", r)
			retErr = multierror.Append(retErr, fmt.Errorf("panic generating audit log"))
		}

		ret = retErr.ErrorOrNil()

		if ret != nil {
			metrics.IncrCounter([]string{"audit", "log_request_failure"}, 1.0)
		}
	}()

	// All logged requests must have an identifier
	//if req.ID == "" {
	//	a.logger.Error("audit: missing identifier in request object", "request_path", req.Path)
	//	retErr = multierror.Append(retErr, fmt.Errorf("missing identifier in request object: %s", req.Path))
	//	return
	//}

	headers := req.Headers
	defer func() {
		req.Headers = headers
	}()

	// Ensure at least one backend logs
	anyLogged := false
	for name, be := range a.backends {
		req.Headers = nil
		transHeaders, thErr := headersConfig.ApplyConfig(headers, be.backend.GetHash)
		if thErr != nil {
			a.logger.Error("audit: backend failed to include headers", "backend", name, "error", thErr)
			continue
		}
		req.Headers = transHeaders

		start := time.Now()
		lrErr := be.backend.LogRequest(auth, req, outerErr)
		metrics.MeasureSince([]string{"audit", name, "log_request"}, start)
		if lrErr != nil {
			a.logger.Error("audit: backend failed to log request", "backend", name, "error", lrErr)
		} else {
			anyLogged = true
		}
	}
	if !anyLogged && len(a.backends) > 0 {
		retErr = multierror.Append(retErr, fmt.Errorf("no audit backend succeeded in logging the request"))
	}

	return retErr.ErrorOrNil()
}

// LogResponse is used to ensure all the audit backends have an opportunity to
// log the given response and that *at least one* succeeds.
func (a *AuditBroker) LogResponse(auth *logical.Auth, req *logical.Request,
	resp *logical.Response, headersConfig *AuditedHeadersConfig, err error) (ret error) {
	defer metrics.MeasureSince([]string{"audit", "log_response"}, time.Now())
	a.RLock()
	defer a.RUnlock()

	var retErr *multierror.Error

	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("audit: panic during logging", "request_path", req.Path, "error", r)
			retErr = multierror.Append(retErr, fmt.Errorf("panic generating audit log"))
		}

		ret = retErr.ErrorOrNil()

		if ret != nil {
			metrics.IncrCounter([]string{"audit", "log_response_failure"}, 1.0)
		}
	}()

	headers := req.Headers
	defer func() {
		req.Headers = headers
	}()

	// Ensure at least one backend logs
	anyLogged := false
	for name, be := range a.backends {
		req.Headers = nil
		transHeaders, thErr := headersConfig.ApplyConfig(headers, be.backend.GetHash)
		if thErr != nil {
			a.logger.Error("audit: backend failed to include headers", "backend", name, "error", thErr)
			continue
		}
		req.Headers = transHeaders

		start := time.Now()
		lrErr := be.backend.LogResponse(auth, req, resp, err)
		metrics.MeasureSince([]string{"audit", name, "log_response"}, start)
		if lrErr != nil {
			a.logger.Error("audit: backend failed to log response", "backend", name, "error", lrErr)
		} else {
			anyLogged = true
		}
	}
	if !anyLogged && len(a.backends) > 0 {
		retErr = multierror.Append(retErr, fmt.Errorf("no audit backend succeeded in logging the response"))
	}

	return retErr.ErrorOrNil()
}

func (a *AuditBroker) Invalidate(key string) {
	// For now we ignore the key as this would only apply to salts. We just
	// sort of brute force it on each one.
	a.Lock()
	defer a.Unlock()
	for _, be := range a.backends {
		be.backend.Invalidate()
	}
}
