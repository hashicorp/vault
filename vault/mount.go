package vault

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
)

const (
	// coreMountConfigPath is used to store the mount configuration.
	// Mounts are protected within the Vault itself, which means they
	// can only be viewed or modified after an unseal.
	coreMountConfigPath = "core/mounts"

	// coreLocalMountConfigPath is used to store mount configuration for local
	// (non-replicated) mounts
	coreLocalMountConfigPath = "core/local-mounts"

	// backendBarrierPrefix is the prefix to the UUID used in the
	// barrier view for the backends.
	backendBarrierPrefix = "logical/"

	// systemBarrierPrefix is the prefix used for the
	// system logical backend.
	systemBarrierPrefix = "sys/"

	// mountTableType is the value we expect to find for the mount table and
	// corresponding entries
	mountTableType = "mounts"
)

var (
	// loadMountsFailed if loadMounts encounters an error
	errLoadMountsFailed = errors.New("failed to setup mount table")

	// protectedMounts cannot be remounted
	protectedMounts = []string{
		"audit/",
		"auth/",
		"sys/",
		"cubbyhole/",
	}

	untunableMounts = []string{
		"cubbyhole/",
		"sys/",
		"audit/",
	}

	// singletonMounts can only exist in one location and are
	// loaded by default. These are types, not paths.
	singletonMounts = []string{
		"cubbyhole",
		"system",
		"token",
	}
)

// MountTable is used to represent the internal mount table
type MountTable struct {
	Type    string        `json:"type"`
	Entries []*MountEntry `json:"entries"`
}

// shallowClone returns a copy of the mount table that
// keeps the MountEntry locations, so as not to invalidate
// other locations holding pointers. Care needs to be taken
// if modifying entries rather than modifying the table itself
func (t *MountTable) shallowClone() *MountTable {
	mt := &MountTable{
		Type:    t.Type,
		Entries: make([]*MountEntry, len(t.Entries)),
	}
	for i, e := range t.Entries {
		mt.Entries[i] = e
	}
	return mt
}

// Hash is used to generate a hash value for the mount table
func (t *MountTable) Hash() ([]byte, error) {
	buf, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	hash := sha1.Sum(buf)
	return hash[:], nil
}

// setTaint is used to set the taint on given entry
func (t *MountTable) setTaint(path string, value bool) *MountEntry {
	n := len(t.Entries)
	for i := 0; i < n; i++ {
		if t.Entries[i].Path == path {
			t.Entries[i].Tainted = value
			return t.Entries[i]
		}
	}
	return nil
}

// remove is used to remove a given path entry; returns the entry that was
// removed
func (t *MountTable) remove(path string) *MountEntry {
	n := len(t.Entries)
	for i := 0; i < n; i++ {
		if entry := t.Entries[i]; entry.Path == path {
			t.Entries[i], t.Entries[n-1] = t.Entries[n-1], nil
			t.Entries = t.Entries[:n-1]
			return entry
		}
	}
	return nil
}

// sortEntriesByPath sorts the entries in the table by path and returns the
// table; this is useful for tests
func (t *MountTable) sortEntriesByPath() *MountTable {
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].Path < t.Entries[j].Path
	})
	return t
}

// MountEntry is used to represent a mount table entry
type MountEntry struct {
	Table       string            `json:"table"`             // The table it belongs to
	Path        string            `json:"path"`              // Mount Path
	Type        string            `json:"type"`              // Logical backend Type
	Description string            `json:"description"`       // User-provided description
	UUID        string            `json:"uuid"`              // Barrier view UUID
	Config      MountConfig       `json:"config"`            // Configuration related to this mount (but not backend-derived)
	Options     map[string]string `json:"options"`           // Backend options
	Local       bool              `json:"local"`             // Local mounts are not replicated or affected by replication
	Tainted     bool              `json:"tainted,omitempty"` // Set as a Write-Ahead flag for unmount/remount
}

// MountConfig is used to hold settable options
type MountConfig struct {
	DefaultLeaseTTL time.Duration `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"` // Override for global default
	MaxLeaseTTL     time.Duration `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`             // Override for global default
	ForceNoCache    bool          `json:"force_no_cache" structs:"force_no_cache" mapstructure:"force_no_cache"`          // Override for global default
}

// Returns a deep copy of the mount entry
func (e *MountEntry) Clone() *MountEntry {
	optClone := make(map[string]string)
	for k, v := range e.Options {
		optClone[k] = v
	}
	return &MountEntry{
		Table:       e.Table,
		Path:        e.Path,
		Type:        e.Type,
		Description: e.Description,
		UUID:        e.UUID,
		Config:      e.Config,
		Options:     optClone,
		Local:       e.Local,
		Tainted:     e.Tainted,
	}
}

// Mount is used to mount a new backend to the mount table.
func (c *Core) mount(entry *MountEntry) error {
	// Ensure we end the path in a slash
	if !strings.HasSuffix(entry.Path, "/") {
		entry.Path += "/"
	}

	// Prevent protected paths from being mounted
	for _, p := range protectedMounts {
		if strings.HasPrefix(entry.Path, p) {
			return logical.CodedError(403, fmt.Sprintf("cannot mount '%s'", entry.Path))
		}
	}

	// Do not allow more than one instance of a singleton mount
	for _, p := range singletonMounts {
		if entry.Type == p {
			return logical.CodedError(403, fmt.Sprintf("Cannot mount more than one instance of '%s'", entry.Type))
		}
	}

	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// Verify there is no conflicting mount
	if match := c.router.MatchingMount(entry.Path); match != "" {
		return logical.CodedError(409, fmt.Sprintf("existing mount at %s", match))
	}

	// Generate a new UUID and view
	if entry.UUID == "" {
		entryUUID, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}
		entry.UUID = entryUUID
	}
	viewPath := backendBarrierPrefix + entry.UUID + "/"
	view := NewBarrierView(c.barrier, viewPath)
	sysView := c.mountEntrySysView(entry)

	backend, err := c.newLogicalBackend(entry.Type, sysView, view, nil)
	if err != nil {
		return err
	}
	if backend == nil {
		return fmt.Errorf("nil backend of type %q returned from creation function", entry.Type)
	}

	// Call initialize; this takes care of init tasks that must be run after
	// the ignore paths are collected
	if err := backend.Initialize(); err != nil {
		return err
	}

	newTable := c.mounts.shallowClone()
	newTable.Entries = append(newTable.Entries, entry)
	if err := c.persistMounts(newTable, entry.Local); err != nil {
		c.logger.Error("core: failed to update mount table", "error", err)
		return logical.CodedError(500, "failed to update mount table")
	}
	c.mounts = newTable

	if err := c.router.Mount(backend, entry.Path, entry, view); err != nil {
		return err
	}

	if c.logger.IsInfo() {
		c.logger.Info("core: successful mount", "path", entry.Path, "type", entry.Type)
	}
	return nil
}

// Unmount is used to unmount a path. The boolean indicates whether the mount
// was found.
func (c *Core) unmount(path string) (bool, error) {
	// Ensure we end the path in a slash
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// Prevent protected paths from being unmounted
	for _, p := range protectedMounts {
		if strings.HasPrefix(path, p) {
			return true, fmt.Errorf("cannot unmount '%s'", path)
		}
	}

	// Verify exact match of the route
	match := c.router.MatchingMount(path)
	if match == "" || path != match {
		return false, fmt.Errorf("no matching mount")
	}

	// Get the view for this backend
	view := c.router.MatchingStorageView(path)

	// Mark the entry as tainted
	if err := c.taintMountEntry(path); err != nil {
		return true, err
	}

	// Taint the router path to prevent routing. Note that in-flight requests
	// are uncertain, right now.
	if err := c.router.Taint(path); err != nil {
		return true, err
	}

	// Invoke the rollback manager a final time
	if err := c.rollback.Rollback(path); err != nil {
		return true, err
	}

	// Revoke all the dynamic keys
	if err := c.expiration.RevokePrefix(path); err != nil {
		return true, err
	}

	// Call cleanup function if it exists
	backend := c.router.MatchingBackend(path)
	if backend != nil {
		backend.Cleanup()
	}

	// Unmount the backend entirely
	if err := c.router.Unmount(path); err != nil {
		return true, err
	}

	// Clear the data in the view
	if err := logical.ClearView(view); err != nil {
		return true, err
	}

	// Remove the mount table entry
	if err := c.removeMountEntry(path); err != nil {
		return true, err
	}
	if c.logger.IsInfo() {
		c.logger.Info("core: successfully unmounted", "path", path)
	}
	return true, nil
}

// removeMountEntry is used to remove an entry from the mount table
func (c *Core) removeMountEntry(path string) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// Remove the entry from the mount table
	newTable := c.mounts.shallowClone()
	entry := newTable.remove(path)
	if entry == nil {
		c.logger.Error("core: nil entry found removing entry in mounts table", "path", path)
		return logical.CodedError(500, "failed to remove entry in mounts table")
	}

	// When unmounting all entries the JSON code will load back up from storage
	// as a nil slice, which kills tests...just set it nil explicitly
	if len(newTable.Entries) == 0 {
		newTable.Entries = nil
	}

	// Update the mount table
	if err := c.persistMounts(newTable, entry.Local); err != nil {
		c.logger.Error("core: failed to remove entry from mounts table", "error", err)
		return logical.CodedError(500, "failed to remove entry from mounts table")
	}

	c.mounts = newTable
	return nil
}

// taintMountEntry is used to mark an entry in the mount table as tainted
func (c *Core) taintMountEntry(path string) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// As modifying the taint of an entry affects shallow clones,
	// we simply use the original
	entry := c.mounts.setTaint(path, true)
	if entry == nil {
		c.logger.Error("core: nil entry found tainting entry in mounts table", "path", path)
		return logical.CodedError(500, "failed to taint entry in mounts table")
	}

	// Update the mount table
	if err := c.persistMounts(c.mounts, entry.Local); err != nil {
		c.logger.Error("core: failed to taint entry in mounts table", "error", err)
		return logical.CodedError(500, "failed to taint entry in mounts table")
	}

	return nil
}

// Remount is used to remount a path at a new mount point.
func (c *Core) remount(src, dst string) error {
	// Ensure we end the path in a slash
	if !strings.HasSuffix(src, "/") {
		src += "/"
	}
	if !strings.HasSuffix(dst, "/") {
		dst += "/"
	}

	// Prevent protected paths from being remounted
	for _, p := range protectedMounts {
		if strings.HasPrefix(src, p) {
			return fmt.Errorf("cannot remount '%s'", src)
		}
	}

	// Verify exact match of the route
	match := c.router.MatchingMount(src)
	if match == "" || src != match {
		return fmt.Errorf("no matching mount at '%s'", src)
	}

	if match := c.router.MatchingMount(dst); match != "" {
		return fmt.Errorf("existing mount at '%s'", match)
	}

	// Mark the entry as tainted
	if err := c.taintMountEntry(src); err != nil {
		return err
	}

	// Taint the router path to prevent routing
	if err := c.router.Taint(src); err != nil {
		return err
	}

	// Invoke the rollback manager a final time
	if err := c.rollback.Rollback(src); err != nil {
		return err
	}

	// Revoke all the dynamic keys
	if err := c.expiration.RevokePrefix(src); err != nil {
		return err
	}

	c.mountsLock.Lock()
	var ent *MountEntry
	for _, ent = range c.mounts.Entries {
		if ent.Path == src {
			ent.Path = dst
			ent.Tainted = false
			break
		}
	}

	if ent == nil {
		c.logger.Error("core: failed to find entry in mounts table")
		return logical.CodedError(500, "failed to find entry in mounts table")
	}

	// Update the mount table
	if err := c.persistMounts(c.mounts, ent.Local); err != nil {
		ent.Path = src
		ent.Tainted = true
		c.mountsLock.Unlock()
		c.logger.Error("core: failed to update mounts table", "error", err)
		return logical.CodedError(500, "failed to update mounts table")
	}
	c.mountsLock.Unlock()

	// Remount the backend
	if err := c.router.Remount(src, dst); err != nil {
		return err
	}

	// Un-taint the path
	if err := c.router.Untaint(dst); err != nil {
		return err
	}

	if c.logger.IsInfo() {
		c.logger.Info("core: successful remount", "old_path", src, "new_path", dst)
	}
	return nil
}

// loadMounts is invoked as part of postUnseal to load the mount table
func (c *Core) loadMounts() error {
	mountTable := &MountTable{}
	localMountTable := &MountTable{}
	// Load the existing mount table
	raw, err := c.barrier.Get(coreMountConfigPath)
	if err != nil {
		c.logger.Error("core: failed to read mount table", "error", err)
		return errLoadMountsFailed
	}
	rawLocal, err := c.barrier.Get(coreLocalMountConfigPath)
	if err != nil {
		c.logger.Error("core: failed to read local mount table", "error", err)
		return errLoadMountsFailed
	}

	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	if raw != nil {
		// Check if the persisted value has canary in the beginning. If
		// yes, decompress the table and then JSON decode it. If not,
		// simply JSON decode it.
		if err := jsonutil.DecodeJSON(raw.Value, mountTable); err != nil {
			c.logger.Error("core: failed to decompress and/or decode the mount table", "error", err)
			return err
		}
		c.mounts = mountTable
	}
	if rawLocal != nil {
		if err := jsonutil.DecodeJSON(rawLocal.Value, localMountTable); err != nil {
			c.logger.Error("core: failed to decompress and/or decode the local mount table", "error", err)
			return err
		}
		c.mounts.Entries = append(c.mounts.Entries, localMountTable.Entries...)
	}

	// Ensure that required entries are loaded, or new ones
	// added may never get loaded at all. Note that this
	// is only designed to work with singletons, as it checks
	// by type only.
	if c.mounts != nil {
		needPersist := false

		// Upgrade to typed mount table
		if c.mounts.Type == "" {
			c.mounts.Type = mountTableType
			needPersist = true
		}

		for _, requiredMount := range requiredMountTable().Entries {
			foundRequired := false
			for _, coreMount := range c.mounts.Entries {
				if coreMount.Type == requiredMount.Type {
					foundRequired = true
					break
				}
			}
			// In a replication scenario we will let sync invalidation take
			// care of creating a new required mount that doesn't exist yet.
			// This should only happen in the upgrade case where a new one is
			// introduced on the primary; otherwise initial bootstrapping will
			// ensure this comes over. If we upgrade first, we simply don't
			// create the mount, so we won't conflict when we sync. If this is
			// local (e.g. cubbyhole) we do still add it.
			if !foundRequired && (c.replicationState != consts.ReplicationSecondary || requiredMount.Local) {
				c.mounts.Entries = append(c.mounts.Entries, requiredMount)
				needPersist = true
			}
		}

		// Upgrade to table-scoped entries
		for _, entry := range c.mounts.Entries {
			if entry.Type == "cubbyhole" && !entry.Local {
				entry.Local = true
				needPersist = true
			}
			if entry.Table == "" {
				entry.Table = c.mounts.Type
				needPersist = true
			}
		}

		// Done if we have restored the mount table and we don't need
		// to persist
		if !needPersist {
			return nil
		}
	} else {
		// Create and persist the default mount table
		c.mounts = defaultMountTable()
	}

	if err := c.persistMounts(c.mounts, false); err != nil {
		c.logger.Error("core: failed to persist mount table", "error", err)
		return errLoadMountsFailed
	}
	return nil
}

// persistMounts is used to persist the mount table after modification
func (c *Core) persistMounts(table *MountTable, localOnly bool) error {
	if table.Type != mountTableType {
		c.logger.Error("core: given table to persist has wrong type", "actual_type", table.Type, "expected_type", mountTableType)
		return fmt.Errorf("invalid table type given, not persisting")
	}

	for _, entry := range table.Entries {
		if entry.Table != table.Type {
			c.logger.Error("core: given entry to persist in mount table has wrong table value", "path", entry.Path, "entry_table_type", entry.Table, "actual_type", table.Type)
			return fmt.Errorf("invalid mount entry found, not persisting")
		}
	}

	nonLocalMounts := &MountTable{
		Type: mountTableType,
	}

	localMounts := &MountTable{
		Type: mountTableType,
	}

	for _, entry := range table.Entries {
		if entry.Local {
			localMounts.Entries = append(localMounts.Entries, entry)
		} else {
			nonLocalMounts.Entries = append(nonLocalMounts.Entries, entry)
		}
	}

	if !localOnly {
		// Encode the mount table into JSON and compress it (lzw).
		compressedBytes, err := jsonutil.EncodeJSONAndCompress(nonLocalMounts, nil)
		if err != nil {
			c.logger.Error("core: failed to encode and/or compress the mount table", "error", err)
			return err
		}

		// Create an entry
		entry := &Entry{
			Key:   coreMountConfigPath,
			Value: compressedBytes,
		}

		// Write to the physical backend
		if err := c.barrier.Put(entry); err != nil {
			c.logger.Error("core: failed to persist mount table", "error", err)
			return err
		}
	}

	// Repeat with local mounts
	compressedBytes, err := jsonutil.EncodeJSONAndCompress(localMounts, nil)
	if err != nil {
		c.logger.Error("core: failed to encode and/or compress the local mount table", "error", err)
		return err
	}

	entry := &Entry{
		Key:   coreLocalMountConfigPath,
		Value: compressedBytes,
	}

	if err := c.barrier.Put(entry); err != nil {
		c.logger.Error("core: failed to persist local mount table", "error", err)
		return err
	}

	return nil
}

// setupMounts is invoked after we've loaded the mount table to
// initialize the logical backends and setup the router
func (c *Core) setupMounts() error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	var backend logical.Backend
	var view *BarrierView
	var err error

	for _, entry := range c.mounts.Entries {
		// Initialize the backend, special casing for system
		barrierPath := backendBarrierPrefix + entry.UUID + "/"
		if entry.Type == "system" {
			barrierPath = systemBarrierPrefix
		}

		// Create a barrier view using the UUID
		view = NewBarrierView(c.barrier, barrierPath)
		sysView := c.mountEntrySysView(entry)
		// Initialize the backend
		// Create the new backend
		backend, err = c.newLogicalBackend(entry.Type, sysView, view, nil)
		if err != nil {
			c.logger.Error("core: failed to create mount entry", "path", entry.Path, "error", err)
			return errLoadMountsFailed
		}
		if backend == nil {
			return fmt.Errorf("created mount entry of type %q is nil", entry.Type)
		}

		if err := backend.Initialize(); err != nil {
			return err
		}

		switch entry.Type {
		case "system":
			c.systemBarrierView = view
		case "cubbyhole":
			ch := backend.(*CubbyholeBackend)
			ch.saltUUID = entry.UUID
			ch.storageView = view
		}

		// Mount the backend
		err = c.router.Mount(backend, entry.Path, entry, view)
		if err != nil {
			c.logger.Error("core: failed to mount entry", "path", entry.Path, "error", err)
			return errLoadMountsFailed
		} else {
			if c.logger.IsInfo() {
				c.logger.Info("core: successfully mounted backend", "type", entry.Type, "path", entry.Path)
			}
		}

		// Ensure the path is tainted if set in the mount table
		if entry.Tainted {
			c.router.Taint(entry.Path)
		}
	}
	return nil
}

// unloadMounts is used before we seal the vault to reset the mounts to
// their unloaded state, calling Cleanup if defined. This is reversed by load and setup mounts.
func (c *Core) unloadMounts() error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	if c.mounts != nil {
		mountTable := c.mounts.shallowClone()
		for _, e := range mountTable.Entries {
			backend := c.router.MatchingBackend(e.Path)
			if backend != nil {
				backend.Cleanup()
			}
		}
	}

	c.mounts = nil
	c.router = NewRouter()
	c.systemBarrierView = nil
	return nil
}

// newLogicalBackend is used to create and configure a new logical backend by name
func (c *Core) newLogicalBackend(t string, sysView logical.SystemView, view logical.Storage, conf map[string]string) (logical.Backend, error) {
	f, ok := c.logicalBackends[t]
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
	if b == nil {
		return nil, fmt.Errorf("nil backend of type %q returned from factory", t)
	}
	return b, nil
}

// mountEntrySysView creates a logical.SystemView from global and
// mount-specific entries; because this should be called when setting
// up a mountEntry, it doesn't check to ensure that me is not nil
func (c *Core) mountEntrySysView(entry *MountEntry) logical.SystemView {
	return dynamicSystemView{
		core:       c,
		mountEntry: entry,
	}
}

// defaultMountTable creates a default mount table
func defaultMountTable() *MountTable {
	table := &MountTable{
		Type: mountTableType,
	}
	mountUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(fmt.Sprintf("could not create default mount table UUID: %v", err))
	}
	genericMount := &MountEntry{
		Table:       mountTableType,
		Path:        "secret/",
		Type:        "generic",
		Description: "generic secret storage",
		UUID:        mountUUID,
	}
	table.Entries = append(table.Entries, genericMount)
	table.Entries = append(table.Entries, requiredMountTable().Entries...)
	return table
}

// requiredMountTable() creates a mount table with entries required
// to be available
func requiredMountTable() *MountTable {
	table := &MountTable{
		Type: mountTableType,
	}
	cubbyholeUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(fmt.Sprintf("could not create cubbyhole UUID: %v", err))
	}
	cubbyholeMount := &MountEntry{
		Table:       mountTableType,
		Path:        "cubbyhole/",
		Type:        "cubbyhole",
		Description: "per-token private secret storage",
		UUID:        cubbyholeUUID,
		Local:       true,
	}

	sysUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(fmt.Sprintf("could not create sys UUID: %v", err))
	}
	sysMount := &MountEntry{
		Table:       mountTableType,
		Path:        "sys/",
		Type:        "system",
		Description: "system endpoints used for control, policy and debugging",
		UUID:        sysUUID,
	}
	table.Entries = append(table.Entries, cubbyholeMount)
	table.Entries = append(table.Entries, sysMount)
	return table
}

// This function returns tables that are singletons. The main usage of this is
// for replication, so we can send over mount info (especially, UUIDs of
// mounts, which are used for salts) for mounts that may not be able to be
// handled normally. After saving these values on the secondary, we let normal
// sync invalidation do its thing. Because of its use for replication, we
// exclude local mounts.
func (c *Core) singletonMountTables() (mounts, auth *MountTable) {
	mounts = &MountTable{}
	auth = &MountTable{}

	c.mountsLock.RLock()
	for _, entry := range c.mounts.Entries {
		if strutil.StrListContains(singletonMounts, entry.Type) && !entry.Local {
			mounts.Entries = append(mounts.Entries, entry)
		}
	}
	c.mountsLock.RUnlock()

	c.authLock.RLock()
	for _, entry := range c.auth.Entries {
		if strutil.StrListContains(singletonMounts, entry.Type) && !entry.Local {
			auth.Entries = append(auth.Entries, entry)
		}
	}
	c.authLock.RUnlock()

	return
}
