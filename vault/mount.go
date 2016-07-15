package vault

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
)

const (
	// coreMountConfigPath is used to store the mount configuration.
	// Mounts are protected within the Vault itself, which means they
	// can only be viewed or modified after an unseal.
	coreMountConfigPath = "core/mounts"

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
	}
)

// MountTable is used to represent the internal mount table
type MountTable struct {
	Type    string        `json:"type"`
	Entries []*MountEntry `json:"entries"`
}

// ShallowClone returns a copy of the mount table that
// keeps the MountEntry locations, so as not to invalidate
// other locations holding pointers. Care needs to be taken
// if modifying entries rather than modifying the table itself
func (t *MountTable) ShallowClone() *MountTable {
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

// Find is used to lookup an entry
func (t *MountTable) Find(path string) *MountEntry {
	n := len(t.Entries)
	for i := 0; i < n; i++ {
		if t.Entries[i].Path == path {
			return t.Entries[i]
		}
	}
	return nil
}

// SetTaint is used to set the taint on given entry
func (t *MountTable) SetTaint(path string, value bool) bool {
	n := len(t.Entries)
	for i := 0; i < n; i++ {
		if t.Entries[i].Path == path {
			t.Entries[i].Tainted = value
			return true
		}
	}
	return false
}

// Remove is used to remove a given path entry
func (t *MountTable) Remove(path string) bool {
	n := len(t.Entries)
	for i := 0; i < n; i++ {
		if t.Entries[i].Path == path {
			t.Entries[i], t.Entries[n-1] = t.Entries[n-1], nil
			t.Entries = t.Entries[:n-1]
			return true
		}
	}
	return false
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
	Tainted     bool              `json:"tainted,omitempty"` // Set as a Write-Ahead flag for unmount/remount
}

// MountConfig is used to hold settable options
type MountConfig struct {
	DefaultLeaseTTL time.Duration `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"` // Override for global default
	MaxLeaseTTL     time.Duration `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`             // Override for global default
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
	}
}

// Mount is used to mount a new backend to the mount table.
func (c *Core) mount(me *MountEntry) error {
	// Ensure we end the path in a slash
	if !strings.HasSuffix(me.Path, "/") {
		me.Path += "/"
	}

	// Prevent protected paths from being mounted
	for _, p := range protectedMounts {
		if strings.HasPrefix(me.Path, p) {
			return logical.CodedError(403, fmt.Sprintf("cannot mount '%s'", me.Path))
		}
	}

	// Do not allow more than one instance of a singleton mount
	for _, p := range singletonMounts {
		if me.Type == p {
			return logical.CodedError(403, fmt.Sprintf("Cannot mount more than one instance of '%s'", me.Type))
		}
	}

	// Verify there is no conflicting mount
	if match := c.router.MatchingMount(me.Path); match != "" {
		return logical.CodedError(409, fmt.Sprintf("existing mount at %s", match))
	}

	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// Generate a new UUID and view
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	me.UUID = meUUID
	view := NewBarrierView(c.barrier, backendBarrierPrefix+me.UUID+"/")

	backend, err := c.newLogicalBackend(me.Type, c.mountEntrySysView(me), view, nil)
	if err != nil {
		return err
	}

	// Update the mount table
	newTable := c.mounts.ShallowClone()
	newTable.Entries = append(newTable.Entries, me)
	if err := c.persistMounts(newTable); err != nil {
		return errors.New("failed to update mount table")
	}
	c.mounts = newTable

	// Mount the backend
	if err := c.router.Mount(backend, me.Path, me, view); err != nil {
		return err
	}
	c.logger.Printf("[INFO] core: mounted '%s' type: %s", me.Path, me.Type)
	return nil
}

// Unmount is used to unmount a path.
func (c *Core) unmount(path string) error {
	// Ensure we end the path in a slash
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// Prevent protected paths from being unmounted
	for _, p := range protectedMounts {
		if strings.HasPrefix(path, p) {
			return fmt.Errorf("cannot unmount '%s'", path)
		}
	}

	// Verify exact match of the route
	match := c.router.MatchingMount(path)
	if match == "" || path != match {
		return fmt.Errorf("no matching mount")
	}

	// Store the view for this backend
	view := c.router.MatchingStorageView(path)

	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// Mark the entry as tainted
	if err := c.taintMountEntry(path); err != nil {
		return err
	}

	// Taint the router path to prevent routing
	if err := c.router.Taint(path); err != nil {
		return err
	}

	// Invoke the rollback manager a final time
	if err := c.rollback.Rollback(path); err != nil {
		return err
	}

	// Revoke all the dynamic keys
	if err := c.expiration.RevokePrefix(path); err != nil {
		return err
	}

	// Unmount the backend entirely
	if err := c.router.Unmount(path); err != nil {
		return err
	}

	// Clear the data in the view
	if err := ClearView(view); err != nil {
		return err
	}

	// Remove the mount table entry
	if err := c.removeMountEntry(path); err != nil {
		return err
	}
	c.logger.Printf("[INFO] core: unmounted '%s'", path)
	return nil
}

// removeMountEntry is used to remove an entry from the mount table
func (c *Core) removeMountEntry(path string) error {
	// Remove the entry from the mount table
	newTable := c.mounts.ShallowClone()
	newTable.Remove(path)

	// Update the mount table
	if err := c.persistMounts(newTable); err != nil {
		return errors.New("failed to update mount table")
	}

	c.mounts = newTable
	return nil
}

// taintMountEntry is used to mark an entry in the mount table as tainted
func (c *Core) taintMountEntry(path string) error {
	// As modifying the taint of an entry affects shallow clones,
	// we simply use the original
	c.mounts.SetTaint(path, true)

	// Update the mount table
	if err := c.persistMounts(c.mounts); err != nil {
		return errors.New("failed to update mount table")
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

	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

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

	var ent *MountEntry
	for _, ent = range c.mounts.Entries {
		if ent.Path == src {
			ent.Path = dst
			ent.Tainted = false
			break
		}
	}

	// Update the mount table
	if err := c.persistMounts(c.mounts); err != nil {
		ent.Path = src
		ent.Tainted = true
		return errors.New("failed to update mount table")
	}

	// Remount the backend
	if err := c.router.Remount(src, dst); err != nil {
		return err
	}

	// Un-taint the path
	if err := c.router.Untaint(dst); err != nil {
		return err
	}

	c.logger.Printf("[INFO] core: remounted '%s' to '%s'", src, dst)
	return nil
}

// loadMounts is invoked as part of postUnseal to load the mount table
func (c *Core) loadMounts() error {
	mountTable := &MountTable{}
	// Load the existing mount table
	raw, err := c.barrier.Get(coreMountConfigPath)
	if err != nil {
		c.logger.Printf("[ERR] core: failed to read mount table: %v", err)
		return errLoadMountsFailed
	}

	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	if raw != nil {
		if err := jsonutil.DecodeJSON(raw.Value, mountTable); err != nil {
			c.logger.Printf("[ERR] core: failed to decode mount table: %v", err)
			return errLoadMountsFailed
		}
		c.mounts = mountTable
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
			if !foundRequired {
				c.mounts.Entries = append(c.mounts.Entries, requiredMount)
				needPersist = true
			}
		}

		// Upgrade to table-scoped entries
		for _, entry := range c.mounts.Entries {
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

	if err := c.persistMounts(c.mounts); err != nil {
		return errLoadMountsFailed
	}
	return nil
}

// persistMounts is used to persist the mount table after modification
func (c *Core) persistMounts(table *MountTable) error {
	if table.Type != mountTableType {
		c.logger.Printf(
			"[ERR] core: given table to persist has type %s but need type %s",
			table.Type,
			mountTableType)
		return fmt.Errorf("invalid table type given, not persisting")
	}

	for _, entry := range table.Entries {
		if entry.Table != table.Type {
			c.logger.Printf(
				"[ERR] core: entry in mount table with path %s has table value %s but is in table %s, refusing to persist",
				entry.Path,
				entry.Table,
				table.Type)
			return fmt.Errorf("invalid mount entry found, not persisting")
		}
	}

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

		// Initialize the backend
		// Create the new backend
		backend, err = c.newLogicalBackend(entry.Type, c.mountEntrySysView(entry), view, nil)
		if err != nil {
			c.logger.Printf(
				"[ERR] core: failed to create mount entry %s: %v",
				entry.Path, err)
			return errLoadMountsFailed
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
			c.logger.Printf("[ERR] core: failed to mount entry %s: %v", entry.Path, err)
			return errLoadMountsFailed
		} else {
			c.logger.Printf("[INFO] core: mounted backend of type %s at %s", entry.Type, entry.Path)
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
		mountTable := c.mounts.ShallowClone()
		for _, e := range mountTable.Entries {
			prefix := e.Path
			b, ok := c.router.root.Get(prefix)
			if ok {
				b.(*routeEntry).backend.Cleanup()
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
	return b, nil
}

// mountEntrySysView creates a logical.SystemView from global and
// mount-specific entries; because this should be called when setting
// up a mountEntry, it doesn't check to ensure that me is not nil
func (c *Core) mountEntrySysView(me *MountEntry) logical.SystemView {
	return dynamicSystemView{
		core:       c,
		mountEntry: me,
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
