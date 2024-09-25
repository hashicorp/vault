// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
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

// loadAuditFailed if loading audit tables encounters an error
var errLoadAuditFailed = errors.New("failed to setup audit table")

func (c *Core) generateAuditTestProbe() (*logical.LogInput, error) {
	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	return &logical.LogInput{
		Type: "request",
		Auth: nil,
		Request: &logical.Request{
			ID:        requestId,
			Operation: "update",
			Path:      "sys/audit/test",
		},
		Response: nil,
		OuterErr: nil,
	}, nil
}

// enableAudit is used to enable a new audit backend that didn't exist in storage beforehand.
func (c *Core) enableAudit(ctx context.Context, entry *MountEntry, updateStorage bool) error {
	// Check ahead of time if the type of audit device we're trying to enable is configured in Vault.
	if _, ok := c.auditBackends[entry.Type]; !ok {
		return fmt.Errorf("unknown backend type: %q: %w", entry.Type, audit.ErrExternalOptions)
	}

	// We can check early to ensure that non-Enterprise versions aren't trying to supply Enterprise only options.
	if audit.HasInvalidOptions(entry.Options) {
		return fmt.Errorf("enterprise-only options supplied: %w", audit.ErrExternalOptions)
	}

	// Ensure we end the path in a slash
	if !strings.HasSuffix(entry.Path, "/") {
		entry.Path += "/"
	}

	// Ensure there is a name
	if entry.Path == "/" {
		return fmt.Errorf("backend path must be specified: %w", audit.ErrExternalOptions)
	}

	if skipTestRaw, ok := entry.Options["skip_test"]; ok {
		skipTest, err := parseutil.ParseBool(skipTestRaw)
		if err != nil {
			return fmt.Errorf("cannot parse supplied 'skip_test' setting: %w", audit.ErrExternalOptions)
		}

		// Reassigning the value means we can ensure that the formatting
		// of it as a string is consistent for future comparisons.
		entry.Options["skip_test"] = strconv.FormatBool(skipTest)
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", audit.ErrInternal, err)
	}
	entry.NamespaceID = ns.ID
	entry.namespace = ns

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
			return fmt.Errorf("path already in use: %w", audit.ErrExternalOptions)
		}
	}

	// Generate a new UUID and view
	if entry.UUID == "" {
		entryUUID, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("%w: %w", audit.ErrInternal, err)
		}
		entry.UUID = entryUUID
	}
	if entry.Accessor == "" {
		accessor, err := c.generateMountAccessor("audit_" + entry.Type)
		if err != nil {
			return fmt.Errorf("%w: %w", audit.ErrInternal, err)
		}
		entry.Accessor = accessor
	}
	viewPath := entry.ViewPath()
	view := NewBarrierView(c.barrier, viewPath)

	// Mark the view as read-only until the mounting is complete and
	// ensure that it is reset after. This ensures that there will be no
	// writes during the construction of the backend.
	defer view.setReadOnlyErr(view.getReadOnlyErr())
	view.setReadOnlyErr(logical.ErrSetupReadOnly)

	// Lookup the new backend
	backend, err := c.newAuditBackend(entry, view, entry.Options)
	if err != nil {
		return err
	}
	if backend == nil {
		return fmt.Errorf("nil audit backend of type %q: %w", entry.Type, audit.ErrInternal)
	}

	if entry.Options["skip_test"] != "true" {
		// Test the new audit device and report failure if it doesn't work.
		testProbe, err := c.generateAuditTestProbe()
		if err != nil {
			return fmt.Errorf("error generating test probe: %w: %w", audit.ErrInternal, err)
		}
		err = backend.LogTestMessage(ctx, testProbe)
		if err != nil {
			c.logger.Error("new audit backend failed test", "path", entry.Path, "type", entry.Type, "error", err)
			return fmt.Errorf("audit backend failed test message: %w", err)
		}
	}

	// Now that we're happy that the backend has been created correctly, we can
	// try to register the backend with the audit broker.
	// An error at this point doesn't cause issues anywhere outside the broker.
	err = c.auditBroker.Register(backend, entry.Local)
	if err != nil {
		return fmt.Errorf("failed to register %q audit backend %q: %w", entry.Type, entry.Path, err)
	}

	// Update a copy of the mount table, we will swap out the oneCore has access to
	// for this one unless there is a problem persisting storage (if required).
	newTable := c.audit.shallowClone()
	newTable.Entries = append(newTable.Entries, entry)
	if updateStorage {
		if err := c.persistAudit(ctx, newTable, entry.Local); err != nil {
			return fmt.Errorf("failed to update audit table: %w: %w", audit.ErrInternal, err)
		}
	}

	// Wrap things up, add the path checker, swap the table, log our success.
	addAuditPathChecker(c, entry, view, viewPath)
	c.audit = newTable
	c.logger.Info("enabled audit backend", "path", entry.Path, "type", entry.Type)
	return nil
}

// disableAudit is used to disable an existing audit backend.
// It returns whether the entry existed, and any errors in disabling it.
func (c *Core) disableAudit(ctx context.Context, path string, updateStorage bool) (bool, error) {
	var existed bool

	// Ensure we end the path in a slash
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// Ensure there is a name
	if path == "/" {
		return existed, fmt.Errorf("backend path must be specified")
	}

	// Remove the entry from the mount table
	c.auditLock.Lock()
	defer c.auditLock.Unlock()

	newTable := c.audit.shallowClone()
	entry, err := newTable.remove(ctx, path)
	if err != nil {
		return existed, err
	}

	// Ensure there was a match
	if entry == nil {
		return existed, fmt.Errorf("no matching backend")
	}

	c.removeAuditReloadFunc(entry)

	// We're satisfied that the entry exists now.
	existed = true

	// When unmounting all entries the JSON code will load back up from storage
	// as a nil slice, which kills tests...just set it nil explicitly
	if len(newTable.Entries) == 0 {
		newTable.Entries = nil
	}

	if updateStorage {
		// Update the audit table
		if err := c.persistAudit(ctx, newTable, entry.Local); err != nil {
			return existed, fmt.Errorf("failed to update audit table: %w: %w", audit.ErrInternal, err)
		}
	}

	c.audit = newTable

	// Unmount the backend, any returned error can be ignored since the
	// Backend will already have been removed from the AuditBroker's map.
	err = c.auditBroker.Deregister(ctx, path)
	if err != nil {
		return existed, fmt.Errorf("failed to deregister %q audit backend %q: %w", entry.Type, entry.Path, err)
	}

	c.logger.Info("disabled audit backend", "path", path)

	removeAuditPathChecker(c, entry)

	return existed, nil
}

// loadAudits is invoked as part of postUnseal to load the audit table
func (c *Core) loadAudits(ctx context.Context) error {
	auditTable := &MountTable{}
	localAuditTable := &MountTable{}

	// Load the existing audit table
	raw, err := c.barrier.Get(ctx, coreAuditConfigPath)
	if err != nil {
		c.logger.Error("failed to read audit table", "error", err)
		return errLoadAuditFailed
	}
	rawLocal, err := c.barrier.Get(ctx, coreLocalAuditConfigPath)
	if err != nil {
		c.logger.Error("failed to read local audit table", "error", err)
		return errLoadAuditFailed
	}

	c.auditLock.Lock()
	defer c.auditLock.Unlock()

	if raw != nil {
		if err := jsonutil.DecodeJSON(raw.Value, auditTable); err != nil {
			c.logger.Error("failed to decode audit table", "error", err)
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
			c.logger.Error("failed to decode local audit table", "error", err)
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

		if entry.NamespaceID == "" {
			entry.NamespaceID = namespace.RootNamespaceID
			needPersist = true
		}
		// Get the namespace from the namespace ID and load it in memory
		ns, err := NamespaceByID(ctx, entry.NamespaceID, c)
		if err != nil {
			return err
		}
		if ns == nil {
			return namespace.ErrNoNamespace
		}
		entry.namespace = ns
	}

	if !needPersist || c.perfStandby {
		return nil
	}

	if err := c.persistAudit(ctx, c.audit, false); err != nil {
		return errLoadAuditFailed
	}
	return nil
}

// persistAudit is used to persist the audit table after modification
func (c *Core) persistAudit(ctx context.Context, table *MountTable, localOnly bool) error {
	if table.Type != auditTableType {
		c.logger.Error("given table to persist has wrong type", "actual_type", table.Type, "expected_type", auditTableType)
		return fmt.Errorf("invalid table type given, not persisting")
	}

	nonLocalAudit := &MountTable{
		Type: auditTableType,
	}

	localAudit := &MountTable{
		Type: auditTableType,
	}

	for _, entry := range table.Entries {
		if entry.Table != table.Type {
			c.logger.Error("given entry to persist in audit table has wrong table value", "path", entry.Path, "entry_table_type", entry.Table, "actual_type", table.Type)
			return fmt.Errorf("invalid audit entry found, not persisting")
		}

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
			c.logger.Error("failed to encode and/or compress audit table", "error", err)
			return err
		}

		// Create an entry
		entry := &logical.StorageEntry{
			Key:   coreAuditConfigPath,
			Value: compressedBytes,
		}

		// Write to the physical backend
		if err := c.barrier.Put(ctx, entry); err != nil {
			c.logger.Error("failed to persist audit table", "error", err)
			return err
		}
	}

	// Repeat with local audit
	compressedBytes, err := jsonutil.EncodeJSONAndCompress(localAudit, nil)
	if err != nil {
		c.logger.Error("failed to encode and/or compress local audit table", "error", err)
		return err
	}

	entry := &logical.StorageEntry{
		Key:   coreLocalAuditConfigPath,
		Value: compressedBytes,
	}

	if err := c.barrier.Put(ctx, entry); err != nil {
		c.logger.Error("failed to persist local audit table", "error", err)
		return err
	}

	return nil
}

// setupAudits is invoked after we've loaded the audit table to initialize the audit backends
func (c *Core) setupAudits(ctx context.Context) error {
	c.auditLock.Lock()
	defer c.auditLock.Unlock()

	brokerLogger := c.baseLogger.Named("audit")

	broker, err := audit.NewBroker(brokerLogger)
	if err != nil {
		return err
	}
	c.auditBroker = broker

	var successCount int

	for _, entry := range c.audit.Entries {
		// Create a barrier view using the UUID
		viewPath := entry.ViewPath()
		view := NewBarrierView(c.barrier, viewPath)

		// Mark the view as read-only until the mounting is complete and
		// ensure that it is reset after. This ensures that there will be no
		// writes during the construction of the backend.
		view.setReadOnlyErr(logical.ErrSetupReadOnly)
		c.postUnsealFuncs = append(c.postUnsealFuncs, func() {
			view.setReadOnlyErr(view.getReadOnlyErr())
		})

		// Initialize the backend
		backend, err := c.newAuditBackend(entry, view, entry.Options)
		if err != nil {
			c.logger.Error("failed to create audit entry", "path", entry.Path, "error", err)
			continue
		}
		if backend == nil {
			c.logger.Error("created audit entry was nil", "path", entry.Path, "type", entry.Type)
			continue
		}

		// Mount the backend
		err = broker.Register(backend, entry.Local)
		if err != nil {
			c.logger.Error("failed to setup audit backed", "path", entry.Path, "type", entry.Type, "error", err)
			continue
		}

		addAuditPathChecker(c, entry, view, viewPath)
		successCount++
	}

	if len(c.audit.Entries) > 0 && successCount == 0 {
		return errLoadAuditFailed
	}

	c.AddLogger(brokerLogger)
	return nil
}

// teardownAudit is used before we seal the vault to reset the audit
// backends to their unloaded state. This is reversed by loadAudits.
func (c *Core) teardownAudits() error {
	c.auditLock.Lock()
	defer c.auditLock.Unlock()

	if c.audit != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for _, entry := range c.audit.Entries {
			c.removeAuditReloadFunc(entry)
			removeAuditPathChecker(c, entry)
			err := c.auditBroker.Deregister(ctx, entry.Path)
			if err != nil {
				c.logger.Error("unable to deregister audit during teardown", "path", entry.Path, "type", entry.Type, "error", err)
			}
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
	case audit.TypeFile:
		key := "audit_file|" + entry.Path
		c.reloadFuncsLock.Lock()

		if c.logger.IsDebug() {
			c.baseLogger.Named("audit").Debug("removing reload function", "path", entry.Path)
		}

		delete(c.reloadFuncs, key)

		c.reloadFuncsLock.Unlock()
	}
}

// newAuditBackend is used to create and configure a new audit backend by name
func (c *Core) newAuditBackend(entry *MountEntry, view logical.Storage, conf map[string]string) (audit.Backend, error) {
	// Ensure that non-Enterprise versions aren't trying to supply Enterprise only options.
	if audit.HasInvalidOptions(entry.Options) {
		return nil, fmt.Errorf("enterprise-only options supplied: %w", audit.ErrInvalidParameter)
	}

	f, ok := c.auditBackends[entry.Type]
	if !ok {
		return nil, fmt.Errorf("unknown backend type: %q: %w", entry.Type, audit.ErrInvalidParameter)
	}
	saltConfig := &salt.Config{
		HMAC:     sha256.New,
		HMACType: "hmac-sha256",
		Location: salt.DefaultLocation,
	}
	auditLogger := c.baseLogger.Named("audit")

	be, err := f(&audit.BackendConfig{
		SaltView:   view,
		SaltConfig: saltConfig,
		Config:     conf,
		MountPath:  entry.Path,
		Logger:     auditLogger,
	}, c.auditedHeaders)
	if err != nil {
		return nil, fmt.Errorf("unable to create new audit backend: %w", err)
	}
	if be == nil {
		return nil, fmt.Errorf("nil backend returned from %q factory function: %w", entry.Type, audit.ErrInternal)
	}

	switch entry.Type {
	case audit.TypeFile:
		key := "audit_file|" + entry.Path

		c.reloadFuncsLock.Lock()

		if auditLogger.IsDebug() {
			auditLogger.Debug("adding reload function", "path", entry.Path)
			if entry.Options != nil {
				auditLogger.Debug("file backend options", "path", entry.Path, "file_path", entry.Options["file_path"])
			}
		}

		c.reloadFuncs[key] = append(c.reloadFuncs[key], func() error {
			auditLogger.Info("reloading file audit backend", "path", entry.Path)
			return be.Reload()
		})

		c.reloadFuncsLock.Unlock()
	case audit.TypeSocket:
		if auditLogger.IsDebug() && entry.Options != nil {
			auditLogger.Debug("socket backend options", "path", entry.Path, "address", entry.Options["address"], "socket type", entry.Options["socket_type"])
		}
	case audit.TypeSyslog:
		if auditLogger.IsDebug() && entry.Options != nil {
			auditLogger.Debug("syslog backend options", "path", entry.Path, "facility", entry.Options["facility"], "tag", entry.Options["tag"])
		}
	}

	c.AddLogger(auditLogger)
	return be, err
}

// defaultAuditTable creates a default audit table
func defaultAuditTable() *MountTable {
	table := &MountTable{
		Type: auditTableType,
	}
	return table
}

type AuditLogger interface {
	AuditRequest(ctx context.Context, input *logical.LogInput) error
	AuditResponse(ctx context.Context, input *logical.LogInput) error
}

type basicAuditor struct {
	c *Core
}

func (b *basicAuditor) AuditRequest(ctx context.Context, input *logical.LogInput) error {
	if b.c.auditBroker == nil {
		return consts.ErrSealed
	}
	return b.c.auditBroker.LogRequest(ctx, input)
}

func (b *basicAuditor) AuditResponse(ctx context.Context, input *logical.LogInput) error {
	if b.c.auditBroker == nil {
		return consts.ErrSealed
	}
	return b.c.auditBroker.LogResponse(ctx, input)
}

type genericAuditor struct {
	c         *Core
	mountType string
	namespace *namespace.Namespace
}

func (g genericAuditor) AuditRequest(ctx context.Context, input *logical.LogInput) error {
	ctx = namespace.ContextWithNamespace(ctx, g.namespace)
	logInput := *input
	logInput.Type = g.mountType + "-request"
	return g.c.auditBroker.LogRequest(ctx, &logInput)
}

func (g genericAuditor) AuditResponse(ctx context.Context, input *logical.LogInput) error {
	ctx = namespace.ContextWithNamespace(ctx, g.namespace)
	logInput := *input
	logInput.Type = g.mountType + "-response"
	return g.c.auditBroker.LogResponse(ctx, &logInput)
}
