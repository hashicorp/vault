package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/copystructure"
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

// ListingVisibilityType represents the types for listing visibility
type ListingVisibilityType string

const (
	// ListingVisibilityDefault is the default value for listing visibility
	ListingVisibilityDefault ListingVisibilityType = ""
	// ListingVisibilityHidden is the hidden type for listing visibility
	ListingVisibilityHidden ListingVisibilityType = "hidden"
	// ListingVisibilityUnauth is the unauth type for listing visibility
	ListingVisibilityUnauth ListingVisibilityType = "unauth"

	systemMountPath    = "sys/"
	identityMountPath  = "identity/"
	cubbyholeMountPath = "cubbyhole/"

	systemMountType    = "system"
	identityMountType  = "identity"
	cubbyholeMountType = "cubbyhole"
	pluginMountType    = "plugin"

	MountTableUpdateStorage   = true
	MountTableNoUpdateStorage = false
)

var (
	// loadMountsFailed if loadMounts encounters an error
	errLoadMountsFailed = errors.New("failed to setup mount table")

	// protectedMounts cannot be remounted
	protectedMounts = []string{
		"audit/",
		"auth/",
		systemMountPath,
		cubbyholeMountPath,
		identityMountPath,
	}

	untunableMounts = []string{
		cubbyholeMountPath,
		systemMountPath,
		"audit/",
		identityMountPath,
	}

	// singletonMounts can only exist in one location and are
	// loaded by default. These are types, not paths.
	singletonMounts = []string{
		cubbyholeMountType,
		systemMountType,
		"token",
		identityMountType,
	}

	// mountAliases maps old backend names to new backend names, allowing us
	// to move/rename backends but maintain backwards compatibility
	mountAliases = map[string]string{"generic": "kv"}
)

func (c *Core) generateMountAccessor(entryType string) (string, error) {
	var accessor string
	for {
		randBytes, err := uuid.GenerateRandomBytes(4)
		if err != nil {
			return "", err
		}
		accessor = fmt.Sprintf("%s_%s", entryType, fmt.Sprintf("%08x", randBytes[0:4]))
		if entry := c.router.MatchingMountByAccessor(accessor); entry == nil {
			break
		}
	}

	return accessor, nil
}

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

// setTaint is used to set the taint on given entry Accepts either the mount
// entry's path or namespace + path, i.e. <ns-path>/secret/ or <ns-path>/token/
func (t *MountTable) setTaint(ctx context.Context, path string, value bool) (*MountEntry, error) {
	n := len(t.Entries)
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		if entry := t.Entries[i]; entry.Path == path && entry.Namespace().ID == ns.ID {
			t.Entries[i].Tainted = value
			return t.Entries[i], nil
		}
	}
	return nil, nil
}

// remove is used to remove a given path entry; returns the entry that was
// removed
func (t *MountTable) remove(ctx context.Context, path string) (*MountEntry, error) {
	n := len(t.Entries)
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		if entry := t.Entries[i]; entry.Path == path && entry.Namespace().ID == ns.ID {
			t.Entries[i], t.Entries[n-1] = t.Entries[n-1], nil
			t.Entries = t.Entries[:n-1]
			return entry, nil
		}
	}
	return nil, nil
}

// sortEntriesByPath sorts the entries in the table by path and returns the
// table; this is useful for tests
func (t *MountTable) sortEntriesByPath() *MountTable {
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].Path < t.Entries[j].Path
	})
	return t
}

// sortEntriesByPath sorts the entries in the table by path and returns the
// table; this is useful for tests
func (t *MountTable) sortEntriesByPathDepth() *MountTable {
	sort.Slice(t.Entries, func(i, j int) bool {
		return len(strings.Split(t.Entries[i].Namespace().Path+t.Entries[i].Path, "/")) < len(strings.Split(t.Entries[j].Namespace().Path+t.Entries[j].Path, "/"))
	})
	return t
}

// MountEntry is used to represent a mount table entry
type MountEntry struct {
	Table            string            `json:"table"`              // The table it belongs to
	Path             string            `json:"path"`               // Mount Path
	Type             string            `json:"type"`               // Logical backend Type
	Description      string            `json:"description"`        // User-provided description
	UUID             string            `json:"uuid"`               // Barrier view UUID
	BackendAwareUUID string            `json:"backend_aware_uuid"` // UUID that can be used by the backend as a helper when a consistent value is needed outside of storage.
	Accessor         string            `json:"accessor"`           // Unique but more human-friendly ID. Does not change, not used for any sensitive things (like as a salt, which the UUID sometimes is).
	Config           MountConfig       `json:"config"`             // Configuration related to this mount (but not backend-derived)
	Options          map[string]string `json:"options"`            // Backend options
	Local            bool              `json:"local"`              // Local mounts are not replicated or affected by replication
	SealWrap         bool              `json:"seal_wrap"`          // Whether to wrap CSPs
	Tainted          bool              `json:"tainted,omitempty"`  // Set as a Write-Ahead flag for unmount/remount
	NamespaceID      string            `json:"namespace_id"`

	// namespace contains the populated namespace
	namespace *namespace.Namespace

	// synthesizedConfigCache is used to cache configuration values. These
	// particular values are cached since we want to get them at a point-in-time
	// without separately managing their locks individually. See SyncCache() for
	// the specific values that are being cached.
	synthesizedConfigCache sync.Map
}

// MountConfig is used to hold settable options
type MountConfig struct {
	DefaultLeaseTTL           time.Duration         `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"` // Override for global default
	MaxLeaseTTL               time.Duration         `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`             // Override for global default
	ForceNoCache              bool                  `json:"force_no_cache" structs:"force_no_cache" mapstructure:"force_no_cache"`          // Override for global default
	AuditNonHMACRequestKeys   []string              `json:"audit_non_hmac_request_keys,omitempty" structs:"audit_non_hmac_request_keys" mapstructure:"audit_non_hmac_request_keys"`
	AuditNonHMACResponseKeys  []string              `json:"audit_non_hmac_response_keys,omitempty" structs:"audit_non_hmac_response_keys" mapstructure:"audit_non_hmac_response_keys"`
	ListingVisibility         ListingVisibilityType `json:"listing_visibility,omitempty" structs:"listing_visibility" mapstructure:"listing_visibility"`
	PassthroughRequestHeaders []string              `json:"passthrough_request_headers,omitempty" structs:"passthrough_request_headers" mapstructure:"passthrough_request_headers"`
	AllowedResponseHeaders    []string              `json:"allowed_response_headers,omitempty" structs:"allowed_response_headers" mapstructure:"allowed_response_headers"`
	TokenType                 logical.TokenType     `json:"token_type" structs:"token_type" mapstructure:"token_type"`

	// PluginName is the name of the plugin registered in the catalog.
	//
	// Deprecated: MountEntry.Type should be used instead for Vault 1.0.0 and beyond.
	PluginName string `json:"plugin_name,omitempty" structs:"plugin_name,omitempty" mapstructure:"plugin_name"`
}

// APIMountConfig is an embedded struct of api.MountConfigInput
type APIMountConfig struct {
	DefaultLeaseTTL           string                `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"`
	MaxLeaseTTL               string                `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`
	ForceNoCache              bool                  `json:"force_no_cache" structs:"force_no_cache" mapstructure:"force_no_cache"`
	AuditNonHMACRequestKeys   []string              `json:"audit_non_hmac_request_keys,omitempty" structs:"audit_non_hmac_request_keys" mapstructure:"audit_non_hmac_request_keys"`
	AuditNonHMACResponseKeys  []string              `json:"audit_non_hmac_response_keys,omitempty" structs:"audit_non_hmac_response_keys" mapstructure:"audit_non_hmac_response_keys"`
	ListingVisibility         ListingVisibilityType `json:"listing_visibility,omitempty" structs:"listing_visibility" mapstructure:"listing_visibility"`
	PassthroughRequestHeaders []string              `json:"passthrough_request_headers,omitempty" structs:"passthrough_request_headers" mapstructure:"passthrough_request_headers"`
	AllowedResponseHeaders    []string              `json:"allowed_response_headers,omitempty" structs:"allowed_response_headers" mapstructure:"allowed_response_headers"`
	TokenType                 string                `json:"token_type" structs:"token_type" mapstructure:"token_type"`

	// PluginName is the name of the plugin registered in the catalog.
	//
	// Deprecated: MountEntry.Type should be used instead for Vault 1.0.0 and beyond.
	PluginName string `json:"plugin_name,omitempty" structs:"plugin_name,omitempty" mapstructure:"plugin_name"`
}

// Clone returns a deep copy of the mount entry
func (e *MountEntry) Clone() (*MountEntry, error) {
	cp, err := copystructure.Copy(e)
	if err != nil {
		return nil, err
	}
	return cp.(*MountEntry), nil
}

// Namespace returns the namespace for the mount entry
func (e *MountEntry) Namespace() *namespace.Namespace {
	return e.namespace
}

// APIPath returns the full API Path for the given mount entry
func (e *MountEntry) APIPath() string {
	path := e.Path
	if e.Table == credentialTableType {
		path = credentialRoutePrefix + path
	}
	return e.namespace.Path + path
}

// SyncCache syncs tunable configuration values to the cache. In the case of
// cached values, they should be retrieved via synthesizedConfigCache.Load()
// instead of accessing them directly through MountConfig.
func (e *MountEntry) SyncCache() {
	if len(e.Config.AuditNonHMACRequestKeys) == 0 {
		e.synthesizedConfigCache.Delete("audit_non_hmac_request_keys")
	} else {
		e.synthesizedConfigCache.Store("audit_non_hmac_request_keys", e.Config.AuditNonHMACRequestKeys)
	}

	if len(e.Config.AuditNonHMACResponseKeys) == 0 {
		e.synthesizedConfigCache.Delete("audit_non_hmac_response_keys")
	} else {
		e.synthesizedConfigCache.Store("audit_non_hmac_response_keys", e.Config.AuditNonHMACResponseKeys)
	}

	if len(e.Config.PassthroughRequestHeaders) == 0 {
		e.synthesizedConfigCache.Delete("passthrough_request_headers")
	} else {
		e.synthesizedConfigCache.Store("passthrough_request_headers", e.Config.PassthroughRequestHeaders)
	}

	if len(e.Config.AllowedResponseHeaders) == 0 {
		e.synthesizedConfigCache.Delete("allowed_response_headers")
	} else {
		e.synthesizedConfigCache.Store("allowed_response_headers", e.Config.AllowedResponseHeaders)
	}
}

func (c *Core) decodeMountTable(ctx context.Context, raw []byte) (*MountTable, error) {
	// Decode into mount table
	mountTable := new(MountTable)
	if err := jsonutil.DecodeJSON(raw, mountTable); err != nil {
		return nil, err
	}

	// Populate the namespace in memory
	var mountEntries []*MountEntry
	for _, entry := range mountTable.Entries {
		if entry.NamespaceID == "" {
			entry.NamespaceID = namespace.RootNamespaceID
		}
		ns, err := NamespaceByID(ctx, entry.NamespaceID, c)
		if err != nil {
			return nil, err
		}
		if ns == nil {
			c.logger.Error("namespace on mount entry not found", "namespace_id", entry.NamespaceID, "mount_path", entry.Path, "mount_description", entry.Description)
			continue
		}

		entry.namespace = ns
		mountEntries = append(mountEntries, entry)
	}

	return &MountTable{
		Type:    mountTable.Type,
		Entries: mountEntries,
	}, nil
}

// Mount is used to mount a new backend to the mount table.
func (c *Core) mount(ctx context.Context, entry *MountEntry) error {
	// Ensure we end the path in a slash
	if !strings.HasSuffix(entry.Path, "/") {
		entry.Path += "/"
	}

	// Prevent protected paths from being mounted
	for _, p := range protectedMounts {
		if strings.HasPrefix(entry.Path, p) && entry.namespace == nil {
			return logical.CodedError(403, fmt.Sprintf("cannot mount %q", entry.Path))
		}
	}

	// Do not allow more than one instance of a singleton mount
	for _, p := range singletonMounts {
		if entry.Type == p {
			return logical.CodedError(403, fmt.Sprintf("mount type of %q is not mountable", entry.Type))
		}
	}
	return c.mountInternal(ctx, entry, MountTableUpdateStorage)
}

func (c *Core) mountInternal(ctx context.Context, entry *MountEntry, updateStorage bool) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	if err := verifyNamespace(c, ns, entry); err != nil {
		return err
	}

	entry.NamespaceID = ns.ID
	entry.namespace = ns

	// Ensure the cache is populated, don't need the result
	NamespaceByID(ctx, ns.ID, c)

	// Verify there are no conflicting mounts
	if match := c.router.MountConflict(ctx, entry.Path); match != "" {
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
	if entry.BackendAwareUUID == "" {
		bUUID, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}
		entry.BackendAwareUUID = bUUID
	}
	if entry.Accessor == "" {
		accessor, err := c.generateMountAccessor(entry.Type)
		if err != nil {
			return err
		}
		entry.Accessor = accessor
	}
	// Sync values to the cache
	entry.SyncCache()

	viewPath := entry.ViewPath()
	view := NewBarrierView(c.barrier, viewPath)

	// Singleton mounts cannot be filtered on a per-secondary basis
	// from replication
	if strutil.StrListContains(singletonMounts, entry.Type) {
		addFilterablePath(c, viewPath)
	}

	nilMount, err := preprocessMount(c, entry, view)
	if err != nil {
		return err
	}
	origReadOnlyErr := view.getReadOnlyErr()

	// Mark the view as read-only until the mounting is complete and
	// ensure that it is reset after. This ensures that there will be no
	// writes during the construction of the backend.
	view.setReadOnlyErr(logical.ErrSetupReadOnly)
	// We defer this because we're already up and running so we don't need to
	// time it for after postUnseal
	defer view.setReadOnlyErr(origReadOnlyErr)

	var backend logical.Backend
	sysView := c.mountEntrySysView(entry)

	backend, err = c.newLogicalBackend(ctx, entry, sysView, view)
	if err != nil {
		return err
	}
	if backend == nil {
		return fmt.Errorf("nil backend of type %q returned from creation function", entry.Type)
	}

	// Check for the correct backend type
	backendType := backend.Type()
	if backendType != logical.TypeLogical {
		if entry.Type != "kv" && entry.Type != "system" && entry.Type != "cubbyhole" {
			return fmt.Errorf(`unknown backend type: "%s"`, entry.Type)
		}
	}

	addPathCheckers(c, entry, backend, viewPath)

	c.setCoreBackend(entry, backend, view)

	// If the mount is filtered or we are on a DR secondary we don't want to
	// keep the actual backend running, so we clean it up and set it to nil
	// so the router does not have a pointer to the object.
	if nilMount {
		backend.Cleanup(ctx)
		backend = nil
	}

	newTable := c.mounts.shallowClone()
	newTable.Entries = append(newTable.Entries, entry)
	if updateStorage {
		if err := c.persistMounts(ctx, newTable, &entry.Local); err != nil {
			c.logger.Error("failed to update mount table", "error", err)
			if err == logical.ErrReadOnly && c.perfStandby {
				return err
			}

			return logical.CodedError(500, "failed to update mount table")
		}
	}
	c.mounts = newTable

	if err := c.router.Mount(backend, entry.Path, entry, view); err != nil {
		return err
	}

	if c.logger.IsInfo() {
		c.logger.Info("successful mount", "namespace", entry.Namespace().Path, "path", entry.Path, "type", entry.Type)
	}
	return nil
}

// Unmount is used to unmount a path. The boolean indicates whether the mount
// was found.
func (c *Core) unmount(ctx context.Context, path string) error {
	// Ensure we end the path in a slash
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// Prevent protected paths from being unmounted
	for _, p := range protectedMounts {
		if strings.HasPrefix(path, p) {
			return fmt.Errorf("cannot unmount %q", path)
		}
	}
	return c.unmountInternal(ctx, path, MountTableUpdateStorage)
}

func (c *Core) unmountInternal(ctx context.Context, path string, updateStorage bool) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	// Verify exact match of the route
	match := c.router.MatchingMount(ctx, path)
	if match == "" || ns.Path+path != match {
		return fmt.Errorf("no matching mount")
	}

	// Get the view for this backend
	view := c.router.MatchingStorageByAPIPath(ctx, path)

	// Get the backend/mount entry for this path, used to remove ignored
	// replication prefixes
	backend := c.router.MatchingBackend(ctx, path)
	entry := c.router.MatchingMountEntry(ctx, path)

	// Mark the entry as tainted
	if err := c.taintMountEntry(ctx, path, updateStorage); err != nil {
		c.logger.Error("failed to taint mount entry for path being unmounted", "error", err, "path", path)
		return err
	}

	// Taint the router path to prevent routing. Note that in-flight requests
	// are uncertain, right now.
	if err := c.router.Taint(ctx, path); err != nil {
		return err
	}

	rCtx := namespace.ContextWithNamespace(c.activeContext, ns)
	if backend != nil && c.rollback != nil {
		// Invoke the rollback manager a final time
		if err := c.rollback.Rollback(rCtx, path); err != nil {
			return err
		}
	}
	if backend != nil && c.expiration != nil && updateStorage {
		// Revoke all the dynamic keys
		if err := c.expiration.RevokePrefix(rCtx, path, true); err != nil {
			return err
		}
	}

	if backend != nil {
		// Call cleanup function if it exists
		backend.Cleanup(ctx)
	}

	// Unmount the backend entirely
	if err := c.router.Unmount(ctx, path); err != nil {
		return err
	}

	viewPath := entry.ViewPath()
	switch {
	case !updateStorage:
		// Don't attempt to clear data, replication will handle this
	case c.IsDRSecondary(), entry.Local, !c.ReplicationState().HasState(consts.ReplicationPerformanceSecondary):
		// Have writable storage, remove the whole thing
		if err := logical.ClearView(ctx, view); err != nil {
			c.logger.Error("failed to clear view for path being unmounted", "error", err, "path", path)
			return err
		}

	case !entry.Local && c.ReplicationState().HasState(consts.ReplicationPerformanceSecondary):
		if err := clearIgnoredPaths(ctx, c, backend, viewPath); err != nil {
			return err
		}
	}
	// Remove the mount table entry
	if err := c.removeMountEntry(ctx, path, updateStorage); err != nil {
		c.logger.Error("failed to remove mount entry for path being unmounted", "error", err, "path", path)
		return err
	}

	removePathCheckers(c, entry, viewPath)

	if c.logger.IsInfo() {
		c.logger.Info("successfully unmounted", "path", path, "namespace", ns.Path)
	}

	return nil
}

// removeMountEntry is used to remove an entry from the mount table
func (c *Core) removeMountEntry(ctx context.Context, path string, updateStorage bool) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// Remove the entry from the mount table
	newTable := c.mounts.shallowClone()
	entry, err := newTable.remove(ctx, path)
	if err != nil {
		return err
	}
	if entry == nil {
		c.logger.Error("nil entry found removing entry in mounts table", "path", path)
		return logical.CodedError(500, "failed to remove entry in mounts table")
	}

	// When unmounting all entries the JSON code will load back up from storage
	// as a nil slice, which kills tests...just set it nil explicitly
	if len(newTable.Entries) == 0 {
		newTable.Entries = nil
	}

	if updateStorage {
		// Update the mount table
		if err := c.persistMounts(ctx, newTable, &entry.Local); err != nil {
			c.logger.Error("failed to remove entry from mounts table", "error", err)
			return logical.CodedError(500, "failed to remove entry from mounts table")
		}
	}

	c.mounts = newTable
	return nil
}

// taintMountEntry is used to mark an entry in the mount table as tainted
func (c *Core) taintMountEntry(ctx context.Context, path string, updateStorage bool) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	// As modifying the taint of an entry affects shallow clones,
	// we simply use the original
	entry, err := c.mounts.setTaint(ctx, path, true)
	if err != nil {
		return err
	}
	if entry == nil {
		c.logger.Error("nil entry found tainting entry in mounts table", "path", path)
		return logical.CodedError(500, "failed to taint entry in mounts table")
	}

	if updateStorage {
		// Update the mount table
		if err := c.persistMounts(ctx, c.mounts, &entry.Local); err != nil {
			if err == logical.ErrReadOnly && c.perfStandby {
				return err
			}

			c.logger.Error("failed to taint entry in mounts table", "error", err)
			return logical.CodedError(500, "failed to taint entry in mounts table")
		}
	}

	return nil
}

// remountForce takes a copy of the mount entry for the path and fully unmounts
// and remounts the backend to pick up any changes, such as filtered paths
func (c *Core) remountForce(ctx context.Context, path string) error {
	me := c.router.MatchingMountEntry(ctx, path)
	if me == nil {
		return fmt.Errorf("cannot find mount for path %q", path)
	}

	me, err := me.Clone()
	if err != nil {
		return err
	}

	if err := c.unmount(ctx, path); err != nil {
		return err
	}
	return c.mount(ctx, me)
}

// Remount is used to remount a path at a new mount point.
func (c *Core) remount(ctx context.Context, src, dst string) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

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
			return fmt.Errorf("cannot remount %q", src)
		}
	}

	// Verify exact match of the route
	srcMatch := c.router.MatchingMountEntry(ctx, src)
	if srcMatch == nil {
		return fmt.Errorf("no matching mount at %q", src)
	}
	if srcMatch.NamespaceID != ns.ID {
		return fmt.Errorf("source mount in a different namespace than request")
	}

	if err := verifyNamespace(c, ns, &MountEntry{Path: dst}); err != nil {
		return err
	}

	if match := c.router.MatchingMount(ctx, dst); match != "" {
		return fmt.Errorf("existing mount at %q", match)
	}

	// Mark the entry as tainted
	if err := c.taintMountEntry(ctx, src, true); err != nil {
		return err
	}

	// Taint the router path to prevent routing
	if err := c.router.Taint(ctx, src); err != nil {
		return err
	}

	if !c.IsDRSecondary() {
		// Invoke the rollback manager a final time
		rCtx := namespace.ContextWithNamespace(c.activeContext, ns)
		if err := c.rollback.Rollback(rCtx, src); err != nil {
			return err
		}

		entry := c.router.MatchingMountEntry(ctx, src)
		if entry == nil {
			return fmt.Errorf("no matching mount at %q", src)
		}

		// Revoke all the dynamic keys
		if err := c.expiration.RevokePrefix(rCtx, src, true); err != nil {
			return err
		}
	}

	c.mountsLock.Lock()
	var entry *MountEntry
	for _, mountEntry := range c.mounts.Entries {
		if mountEntry.Path == src && mountEntry.NamespaceID == ns.ID {
			entry = mountEntry
			entry.Path = dst
			entry.Tainted = false
			break
		}
	}

	if entry == nil {
		c.mountsLock.Unlock()
		c.logger.Error("failed to find entry in mounts table")
		return logical.CodedError(500, "failed to find entry in mounts table")
	}

	// Update the mount table
	if err := c.persistMounts(ctx, c.mounts, &entry.Local); err != nil {
		entry.Path = src
		entry.Tainted = true
		c.mountsLock.Unlock()
		if err == logical.ErrReadOnly && c.perfStandby {
			return err
		}

		c.logger.Error("failed to update mounts table", "error", err)
		return logical.CodedError(500, "failed to update mounts table")
	}
	c.mountsLock.Unlock()

	// Remount the backend
	if err := c.router.Remount(ctx, src, dst); err != nil {
		return err
	}

	// Un-taint the path
	if err := c.router.Untaint(ctx, dst); err != nil {
		return err
	}

	if c.logger.IsInfo() {
		c.logger.Info("successful remount", "old_path", src, "new_path", dst)
	}
	return nil
}

// loadMounts is invoked as part of postUnseal to load the mount table
func (c *Core) loadMounts(ctx context.Context) error {
	// Load the existing mount table
	raw, err := c.barrier.Get(ctx, coreMountConfigPath)
	if err != nil {
		c.logger.Error("failed to read mount table", "error", err)
		return errLoadMountsFailed
	}
	rawLocal, err := c.barrier.Get(ctx, coreLocalMountConfigPath)
	if err != nil {
		c.logger.Error("failed to read local mount table", "error", err)
		return errLoadMountsFailed
	}

	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	if raw != nil {
		// Check if the persisted value has canary in the beginning. If
		// yes, decompress the table and then JSON decode it. If not,
		// simply JSON decode it.
		mountTable, err := c.decodeMountTable(ctx, raw.Value)
		if err != nil {
			c.logger.Error("failed to decompress and/or decode the mount table", "error", err)
			return err
		}
		c.mounts = mountTable
	}

	var needPersist bool
	if c.mounts == nil {
		c.logger.Info("no mounts; adding default mount table")
		c.mounts = c.defaultMountTable()
		needPersist = true
	}

	if rawLocal != nil {
		localMountTable, err := c.decodeMountTable(ctx, rawLocal.Value)
		if err != nil {
			c.logger.Error("failed to decompress and/or decode the local mount table", "error", err)
			return err
		}
		if localMountTable != nil && len(localMountTable.Entries) > 0 {
			c.mounts.Entries = append(c.mounts.Entries, localMountTable.Entries...)
		}
	}

	// Note that this is only designed to work with singletons, as it checks by
	// type only.

	// Upgrade to typed mount table
	if c.mounts.Type == "" {
		c.mounts.Type = mountTableType
		needPersist = true
	}

	for _, requiredMount := range c.requiredMountTable().Entries {
		foundRequired := false
		for _, coreMount := range c.mounts.Entries {
			if coreMount.Type == requiredMount.Type {
				foundRequired = true
				coreMount.Config = requiredMount.Config
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
		if !foundRequired && (!c.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) || requiredMount.Local) {
			c.mounts.Entries = append(c.mounts.Entries, requiredMount)
			needPersist = true
		}
	}

	// Upgrade to table-scoped entries
	for _, entry := range c.mounts.Entries {
		if entry.Type == cubbyholeMountType && !entry.Local {
			entry.Local = true
			needPersist = true
		}
		if entry.Table == "" {
			entry.Table = c.mounts.Type
			needPersist = true
		}
		if entry.Accessor == "" {
			accessor, err := c.generateMountAccessor(entry.Type)
			if err != nil {
				return err
			}
			entry.Accessor = accessor
			needPersist = true
		}
		if entry.BackendAwareUUID == "" {
			bUUID, err := uuid.GenerateUUID()
			if err != nil {
				return err
			}
			entry.BackendAwareUUID = bUUID
			needPersist = true
		}

		if entry.NamespaceID == "" {
			entry.NamespaceID = namespace.RootNamespaceID
			needPersist = true
		}
		ns, err := NamespaceByID(ctx, entry.NamespaceID, c)
		if err != nil {
			return err
		}
		if ns == nil {
			return namespace.ErrNoNamespace
		}
		entry.namespace = ns

		// Sync values to the cache
		entry.SyncCache()
	}

	// Done if we have restored the mount table and we don't need
	// to persist
	if !needPersist {
		return nil
	}

	// Persist both mount tables
	if err := c.persistMounts(ctx, c.mounts, nil); err != nil {
		c.logger.Error("failed to persist mount table", "error", err)
		return errLoadMountsFailed
	}
	return nil
}

// persistMounts is used to persist the mount table after modification
func (c *Core) persistMounts(ctx context.Context, table *MountTable, local *bool) error {
	if table.Type != mountTableType {
		c.logger.Error("given table to persist has wrong type", "actual_type", table.Type, "expected_type", mountTableType)
		return fmt.Errorf("invalid table type given, not persisting")
	}

	for _, entry := range table.Entries {
		if entry.Table != table.Type {
			c.logger.Error("given entry to persist in mount table has wrong table value", "path", entry.Path, "entry_table_type", entry.Table, "actual_type", table.Type)
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

	writeTable := func(mt *MountTable, path string) error {
		// Encode the mount table into JSON and compress it (lzw).
		compressedBytes, err := jsonutil.EncodeJSONAndCompress(mt, nil)
		if err != nil {
			c.logger.Error("failed to encode or compress mount table", "error", err)
			return err
		}

		// Create an entry
		entry := &logical.StorageEntry{
			Key:   path,
			Value: compressedBytes,
		}

		// Write to the physical backend
		if err := c.barrier.Put(ctx, entry); err != nil {
			c.logger.Error("failed to persist mount table", "error", err)
			return err
		}
		return nil
	}

	var err error
	switch {
	case local == nil:
		// Write non-local mounts
		err := writeTable(nonLocalMounts, coreMountConfigPath)
		if err != nil {
			return err
		}

		// Write local mounts
		err = writeTable(localMounts, coreLocalMountConfigPath)
		if err != nil {
			return err
		}
	case *local:
		// Write local mounts
		err = writeTable(localMounts, coreLocalMountConfigPath)
	default:
		// Write non-local mounts
		err = writeTable(nonLocalMounts, coreMountConfigPath)
	}

	return err
}

// setupMounts is invoked after we've loaded the mount table to
// initialize the logical backends and setup the router
func (c *Core) setupMounts(ctx context.Context) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	for _, entry := range c.mounts.sortEntriesByPathDepth().Entries {
		// Initialize the backend, special casing for system
		barrierPath := entry.ViewPath()

		// Create a barrier view using the UUID
		view := NewBarrierView(c.barrier, barrierPath)

		// Singleton mounts cannot be filtered on a per-secondary basis
		// from replication
		if strutil.StrListContains(singletonMounts, entry.Type) {
			addFilterablePath(c, barrierPath)
		}

		// Determining the replicated state of the mount
		nilMount, err := preprocessMount(c, entry, view)
		if err != nil {
			return err
		}
		origReadOnlyErr := view.getReadOnlyErr()

		// Mark the view as read-only until the mounting is complete and
		// ensure that it is reset after. This ensures that there will be no
		// writes during the construction of the backend.
		view.setReadOnlyErr(logical.ErrSetupReadOnly)
		if strutil.StrListContains(singletonMounts, entry.Type) {
			defer view.setReadOnlyErr(origReadOnlyErr)
		} else {
			c.postUnsealFuncs = append(c.postUnsealFuncs, func() {
				view.setReadOnlyErr(origReadOnlyErr)
			})
		}

		var backend logical.Backend
		// Create the new backend
		sysView := c.mountEntrySysView(entry)
		backend, err = c.newLogicalBackend(ctx, entry, sysView, view)
		if err != nil {
			c.logger.Error("failed to create mount entry", "path", entry.Path, "error", err)
			if !c.builtinRegistry.Contains(entry.Type, consts.PluginTypeSecrets) {
				// If we encounter an error instantiating the backend due to an error,
				// skip backend initialization but register the entry to the mount table
				// to preserve storage and path.
				c.logger.Warn("skipping plugin-based mount entry", "path", entry.Path)
				goto ROUTER_MOUNT
			}
			return errLoadMountsFailed
		}
		if backend == nil {
			return fmt.Errorf("created mount entry of type %q is nil", entry.Type)
		}

		{
			// Check for the correct backend type
			backendType := backend.Type()

			if backendType != logical.TypeLogical {
				if entry.Type != "kv" && entry.Type != "system" && entry.Type != "cubbyhole" {
					return fmt.Errorf(`unknown backend type: "%s"`, entry.Type)
				}
			}

			addPathCheckers(c, entry, backend, barrierPath)

			c.setCoreBackend(entry, backend, view)
		}

		// If the mount is filtered or we are on a DR secondary we don't want to
		// keep the actual backend running, so we clean it up and set it to nil
		// so the router does not have a pointer to the object.
		if nilMount {
			backend.Cleanup(ctx)
			backend = nil
		}

	ROUTER_MOUNT:
		// Mount the backend
		err = c.router.Mount(backend, entry.Path, entry, view)
		if err != nil {
			c.logger.Error("failed to mount entry", "path", entry.Path, "error", err)
			return errLoadMountsFailed
		}

		if c.logger.IsInfo() {
			c.logger.Info("successfully mounted backend", "type", entry.Type, "path", entry.Path)
		}

		// Ensure the path is tainted if set in the mount table
		if entry.Tainted {
			c.router.Taint(ctx, entry.Path)
		}

		// Ensure the cache is populated, don't need the result
		NamespaceByID(ctx, entry.NamespaceID, c)
	}
	return nil
}

// unloadMounts is used before we seal the vault to reset the mounts to
// their unloaded state, calling Cleanup if defined. This is reversed by load and setup mounts.
func (c *Core) unloadMounts(ctx context.Context) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	if c.mounts != nil {
		mountTable := c.mounts.shallowClone()
		for _, e := range mountTable.Entries {
			backend := c.router.MatchingBackend(namespace.ContextWithNamespace(ctx, e.namespace), e.Path)
			if backend != nil {
				backend.Cleanup(ctx)
			}

			viewPath := e.ViewPath()
			removePathCheckers(c, e, viewPath)
		}
	}

	c.mounts = nil
	c.router = NewRouter()
	c.systemBarrierView = nil
	return nil
}

// newLogicalBackend is used to create and configure a new logical backend by name
func (c *Core) newLogicalBackend(ctx context.Context, entry *MountEntry, sysView logical.SystemView, view logical.Storage) (logical.Backend, error) {
	t := entry.Type
	if alias, ok := mountAliases[t]; ok {
		t = alias
	}

	f, ok := c.logicalBackends[t]
	if !ok {
		f = plugin.Factory
	}

	// Set up conf to pass in plugin_name
	conf := make(map[string]string, len(entry.Options)+1)
	for k, v := range entry.Options {
		conf[k] = v
	}

	switch {
	case entry.Type == "plugin":
		conf["plugin_name"] = entry.Config.PluginName
	default:
		conf["plugin_name"] = t
	}

	conf["plugin_type"] = consts.PluginTypeSecrets.String()

	backendLogger := c.baseLogger.Named(fmt.Sprintf("secrets.%s.%s", t, entry.Accessor))
	c.AddLogger(backendLogger)
	config := &logical.BackendConfig{
		StorageView: view,
		Logger:      backendLogger,
		Config:      conf,
		System:      sysView,
		BackendUUID: entry.BackendAwareUUID,
	}

	b, err := f(ctx, config)
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
func (c *Core) defaultMountTable() *MountTable {
	table := &MountTable{
		Type: mountTableType,
	}
	table.Entries = append(table.Entries, c.requiredMountTable().Entries...)

	if os.Getenv("VAULT_INTERACTIVE_DEMO_SERVER") != "" {
		mountUUID, err := uuid.GenerateUUID()
		if err != nil {
			panic(fmt.Sprintf("could not create default secret mount UUID: %v", err))
		}
		mountAccessor, err := c.generateMountAccessor("kv")
		if err != nil {
			panic(fmt.Sprintf("could not generate default secret mount accessor: %v", err))
		}
		bUUID, err := uuid.GenerateUUID()
		if err != nil {
			panic(fmt.Sprintf("could not create default secret mount backend UUID: %v", err))
		}

		kvMount := &MountEntry{
			Table:            mountTableType,
			Path:             "secret/",
			Type:             "kv",
			Description:      "key/value secret storage",
			UUID:             mountUUID,
			Accessor:         mountAccessor,
			BackendAwareUUID: bUUID,
			Options: map[string]string{
				"version": "2",
			},
		}
		table.Entries = append(table.Entries, kvMount)
	}

	return table
}

// requiredMountTable() creates a mount table with entries required
// to be available
func (c *Core) requiredMountTable() *MountTable {
	table := &MountTable{
		Type: mountTableType,
	}
	cubbyholeUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(fmt.Sprintf("could not create cubbyhole UUID: %v", err))
	}
	cubbyholeAccessor, err := c.generateMountAccessor("cubbyhole")
	if err != nil {
		panic(fmt.Sprintf("could not generate cubbyhole accessor: %v", err))
	}
	cubbyholeBackendUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(fmt.Sprintf("could not create cubbyhole backend UUID: %v", err))
	}
	cubbyholeMount := &MountEntry{
		Table:            mountTableType,
		Path:             cubbyholeMountPath,
		Type:             cubbyholeMountType,
		Description:      "per-token private secret storage",
		UUID:             cubbyholeUUID,
		Accessor:         cubbyholeAccessor,
		Local:            true,
		BackendAwareUUID: cubbyholeBackendUUID,
	}

	sysUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(fmt.Sprintf("could not create sys UUID: %v", err))
	}
	sysAccessor, err := c.generateMountAccessor("system")
	if err != nil {
		panic(fmt.Sprintf("could not generate sys accessor: %v", err))
	}
	sysBackendUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(fmt.Sprintf("could not create sys backend UUID: %v", err))
	}
	sysMount := &MountEntry{
		Table:            mountTableType,
		Path:             "sys/",
		Type:             systemMountType,
		Description:      "system endpoints used for control, policy and debugging",
		UUID:             sysUUID,
		Accessor:         sysAccessor,
		BackendAwareUUID: sysBackendUUID,
		Config: MountConfig{
			PassthroughRequestHeaders: []string{"Accept"},
		},
	}

	identityUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(fmt.Sprintf("could not create identity mount entry UUID: %v", err))
	}
	identityAccessor, err := c.generateMountAccessor("identity")
	if err != nil {
		panic(fmt.Sprintf("could not generate identity accessor: %v", err))
	}
	identityBackendUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(fmt.Sprintf("could not create identity backend UUID: %v", err))
	}
	identityMount := &MountEntry{
		Table:            mountTableType,
		Path:             "identity/",
		Type:             "identity",
		Description:      "identity store",
		UUID:             identityUUID,
		Accessor:         identityAccessor,
		BackendAwareUUID: identityBackendUUID,
	}

	table.Entries = append(table.Entries, cubbyholeMount)
	table.Entries = append(table.Entries, sysMount)
	table.Entries = append(table.Entries, identityMount)

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
		if strutil.StrListContains(singletonMounts, entry.Type) && !entry.Local && entry.Namespace().ID == namespace.RootNamespaceID {
			mounts.Entries = append(mounts.Entries, entry)
		}
	}
	c.mountsLock.RUnlock()

	c.authLock.RLock()
	for _, entry := range c.auth.Entries {
		if strutil.StrListContains(singletonMounts, entry.Type) && !entry.Local && entry.Namespace().ID == namespace.RootNamespaceID {
			auth.Entries = append(auth.Entries, entry)
		}
	}
	c.authLock.RUnlock()

	return
}

func (c *Core) setCoreBackend(entry *MountEntry, backend logical.Backend, view *BarrierView) {
	switch entry.Type {
	case systemMountType:
		c.systemBackend = backend.(*SystemBackend)
		c.systemBarrierView = view
	case cubbyholeMountType:
		ch := backend.(*CubbyholeBackend)
		ch.saltUUID = entry.UUID
		ch.storageView = view
		c.cubbyholeBackend = ch
	case identityMountType:
		c.identityStore = backend.(*IdentityStore)
	}
}
