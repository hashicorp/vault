// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/experiments"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/versions"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
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

	systemMountType      = "system"
	identityMountType    = "identity"
	cubbyholeMountType   = "cubbyhole"
	pluginMountType      = "plugin"
	mountTypeNSCubbyhole = "ns_cubbyhole"

	MountTableUpdateStorage   = true
	MountTableNoUpdateStorage = false
)

// DeprecationStatus errors
var (
	errMountDeprecated     = errors.New("mount entry associated with deprecated builtin")
	errMountPendingRemoval = errors.New("mount entry associated with pending removal builtin")
	errMountRemoved        = errors.New("mount entry associated with removed builtin")
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

type MountMigrationStatus int

const (
	MigrationInProgressStatus MountMigrationStatus = iota
	MigrationSuccessStatus
	MigrationFailureStatus
)

func (m MountMigrationStatus) String() string {
	switch m {
	case MigrationInProgressStatus:
		return "in-progress"
	case MigrationSuccessStatus:
		return "success"
	case MigrationFailureStatus:
		return "failure"
	}
	return "unknown"
}

type MountMigrationInfo struct {
	SourceMount     string `json:"source_mount"`
	TargetMount     string `json:"target_mount"`
	MigrationStatus string `json:"status"`
}

// tableMetrics is responsible for setting gauge metrics for
// mount table storage sizes (in bytes) and mount table num
// entries. It does this via setGaugeWithLabels. It then
// saves these metrics in a cache for regular reporting in
// a loop, via AddGaugeLoopMetric.

// Note that the reported storage sizes are pre-encryption
// sizes. Currently barrier uses aes-gcm for encryption, which
// preserves plaintext size, adding a constant of 30 bytes of
// padding, which is negligable and subject to change, and thus
// not accounted for.
func (c *Core) tableMetrics(entryCount int, isLocal bool, isAuth bool, compressedTable []byte) {
	if c.metricsHelper == nil {
		// do nothing if metrics are not initialized
		return
	}
	typeAuthLabelMap := map[bool]metrics.Label{
		true:  {Name: "type", Value: "auth"},
		false: {Name: "type", Value: "logical"},
	}

	typeLocalLabelMap := map[bool]metrics.Label{
		true:  {Name: "local", Value: "true"},
		false: {Name: "local", Value: "false"},
	}

	c.metricSink.SetGaugeWithLabels(metricsutil.LogicalTableSizeName,
		float32(entryCount), []metrics.Label{
			typeAuthLabelMap[isAuth],
			typeLocalLabelMap[isLocal],
		})

	c.metricsHelper.AddGaugeLoopMetric(metricsutil.LogicalTableSizeName,
		float32(entryCount), []metrics.Label{
			typeAuthLabelMap[isAuth],
			typeLocalLabelMap[isLocal],
		})

	c.metricSink.SetGaugeWithLabels(metricsutil.PhysicalTableSizeName,
		float32(len(compressedTable)), []metrics.Label{
			typeAuthLabelMap[isAuth],
			typeLocalLabelMap[isLocal],
		})

	c.metricsHelper.AddGaugeLoopMetric(metricsutil.PhysicalTableSizeName,
		float32(len(compressedTable)), []metrics.Label{
			typeAuthLabelMap[isAuth],
			typeLocalLabelMap[isLocal],
		})
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
func (t *MountTable) setTaint(nsID, path string, tainted bool, mountState string) (*MountEntry, error) {
	n := len(t.Entries)
	for i := 0; i < n; i++ {
		if entry := t.Entries[i]; entry.Path == path && entry.Namespace().ID == nsID {
			t.Entries[i].Tainted = tainted
			t.Entries[i].MountState = mountState
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

func (t *MountTable) find(ctx context.Context, path string) (*MountEntry, error) {
	n := len(t.Entries)
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		if entry := t.Entries[i]; entry.Path == path && entry.Namespace().ID == ns.ID {
			return entry, nil
		}
	}
	return nil, nil
}

func (t *MountTable) findByBackendUUID(ctx context.Context, backendUUID string) (*MountEntry, error) {
	n := len(t.Entries)
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		if entry := t.Entries[i]; entry.BackendAwareUUID == backendUUID && entry.Namespace().ID == ns.ID {
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

const mountStateUnmounting = "unmounting"

// MountEntry is used to represent a mount table entry
type MountEntry struct {
	Table                 string            `json:"table"`                             // The table it belongs to
	Path                  string            `json:"path"`                              // Mount Path
	Type                  string            `json:"type"`                              // Logical backend Type. NB: This is the plugin name, e.g. my-vault-plugin, NOT plugin type (e.g. auth).
	Description           string            `json:"description"`                       // User-provided description
	UUID                  string            `json:"uuid"`                              // Barrier view UUID
	BackendAwareUUID      string            `json:"backend_aware_uuid"`                // UUID that can be used by the backend as a helper when a consistent value is needed outside of storage.
	Accessor              string            `json:"accessor"`                          // Unique but more human-friendly ID. Does not change, not used for any sensitive things (like as a salt, which the UUID sometimes is).
	Config                MountConfig       `json:"config"`                            // Configuration related to this mount (but not backend-derived)
	Options               map[string]string `json:"options"`                           // Backend options
	Local                 bool              `json:"local"`                             // Local mounts are not replicated or affected by replication
	SealWrap              bool              `json:"seal_wrap"`                         // Whether to wrap CSPs
	ExternalEntropyAccess bool              `json:"external_entropy_access,omitempty"` // Whether to allow external entropy source access
	Tainted               bool              `json:"tainted,omitempty"`                 // Set as a Write-Ahead flag for unmount/remount
	MountState            string            `json:"mount_state,omitempty"`             // The current mount state.  The only non-empty mount state right now is "unmounting"
	NamespaceID           string            `json:"namespace_id"`

	// namespace contains the populated namespace
	namespace *namespace.Namespace

	// synthesizedConfigCache is used to cache configuration values. These
	// particular values are cached since we want to get them at a point-in-time
	// without separately managing their locks individually. See SyncCache() for
	// the specific values that are being cached.
	synthesizedConfigCache sync.Map

	// version info
	Version        string `json:"plugin_version,omitempty"`         // The semantic version of the mounted plugin, e.g. v1.2.3.
	RunningVersion string `json:"running_plugin_version,omitempty"` // The semantic version of the mounted plugin as reported by the plugin.
	RunningSha256  string `json:"running_sha256,omitempty"`
}

// MountConfig is used to hold settable options
type MountConfig struct {
	DefaultLeaseTTL           time.Duration         `json:"default_lease_ttl,omitempty" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"` // Override for global default
	MaxLeaseTTL               time.Duration         `json:"max_lease_ttl,omitempty" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`             // Override for global default
	ForceNoCache              bool                  `json:"force_no_cache,omitempty" structs:"force_no_cache" mapstructure:"force_no_cache"`          // Override for global default
	AuditNonHMACRequestKeys   []string              `json:"audit_non_hmac_request_keys,omitempty" structs:"audit_non_hmac_request_keys" mapstructure:"audit_non_hmac_request_keys"`
	AuditNonHMACResponseKeys  []string              `json:"audit_non_hmac_response_keys,omitempty" structs:"audit_non_hmac_response_keys" mapstructure:"audit_non_hmac_response_keys"`
	ListingVisibility         ListingVisibilityType `json:"listing_visibility,omitempty" structs:"listing_visibility" mapstructure:"listing_visibility"`
	PassthroughRequestHeaders []string              `json:"passthrough_request_headers,omitempty" structs:"passthrough_request_headers" mapstructure:"passthrough_request_headers"`
	AllowedResponseHeaders    []string              `json:"allowed_response_headers,omitempty" structs:"allowed_response_headers" mapstructure:"allowed_response_headers"`
	TokenType                 logical.TokenType     `json:"token_type,omitempty" structs:"token_type" mapstructure:"token_type"`
	AllowedManagedKeys        []string              `json:"allowed_managed_keys,omitempty" mapstructure:"allowed_managed_keys"`
	UserLockoutConfig         *UserLockoutConfig    `json:"user_lockout_config,omitempty" mapstructure:"user_lockout_config"`

	// PluginName is the name of the plugin registered in the catalog.
	//
	// Deprecated: MountEntry.Type should be used instead for Vault 1.0.0 and beyond.
	PluginName string `json:"plugin_name,omitempty" structs:"plugin_name,omitempty" mapstructure:"plugin_name"`
}

type UserLockoutConfig struct {
	LockoutThreshold    uint64        `json:"lockout_threshold,omitempty" structs:"lockout_threshold" mapstructure:"lockout_threshold"`
	LockoutDuration     time.Duration `json:"lockout_duration,omitempty" structs:"lockout_duration" mapstructure:"lockout_duration"`
	LockoutCounterReset time.Duration `json:"lockout_counter_reset,omitempty" structs:"lockout_counter_reset" mapstructure:"lockout_counter_reset"`
	DisableLockout      bool          `json:"disable_lockout,omitempty" structs:"disable_lockout" mapstructure:"disable_lockout"`
}

type APIUserLockoutConfig struct {
	LockoutThreshold            string `json:"lockout_threshold,omitempty" structs:"lockout_threshold" mapstructure:"lockout_threshold"`
	LockoutDuration             string `json:"lockout_duration,omitempty" structs:"lockout_duration" mapstructure:"lockout_duration"`
	LockoutCounterResetDuration string `json:"lockout_counter_reset_duration,omitempty" structs:"lockout_counter_reset_duration" mapstructure:"lockout_counter_reset_duration"`
	DisableLockout              *bool  `json:"lockout_disable,omitempty" structs:"lockout_disable" mapstructure:"lockout_disable"`
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
	AllowedManagedKeys        []string              `json:"allowed_managed_keys,omitempty" mapstructure:"allowed_managed_keys"`
	UserLockoutConfig         *UserLockoutConfig    `json:"user_lockout_config,omitempty" mapstructure:"user_lockout_config"`
	PluginVersion             string                `json:"plugin_version,omitempty" mapstructure:"plugin_version"`

	// PluginName is the name of the plugin registered in the catalog.
	//
	// Deprecated: MountEntry.Type should be used instead for Vault 1.0.0 and beyond.
	PluginName string `json:"plugin_name,omitempty" structs:"plugin_name,omitempty" mapstructure:"plugin_name"`
}

type FailedLoginUser struct {
	aliasName     string
	mountAccessor string
}

type FailedLoginInfo struct {
	count               uint
	lastFailedLoginTime int
}

// Clone returns a deep copy of the mount entry
func (e *MountEntry) Clone() (*MountEntry, error) {
	cp, err := copystructure.Copy(e)
	if err != nil {
		return nil, err
	}
	return cp.(*MountEntry), nil
}

// IsExternalPlugin returns whether the plugin is running externally
// if the RunningSha256 is non-empty, the builtin is external. Otherwise, it's builtin
func (e *MountEntry) IsExternalPlugin() bool {
	return e.RunningSha256 != ""
}

// MountClass returns the mount class based on Accessor and Path
func (e *MountEntry) MountClass() string {
	if e.Accessor == "" || strings.HasPrefix(e.Path, fmt.Sprintf("%s/", systemMountPath)) {
		return ""
	}

	if e.Table == credentialTableType {
		return consts.PluginTypeCredential.String()
	}

	return consts.PluginTypeSecrets.String()
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

// APIPathNoNamespace returns the API Path without the namespace for the given mount entry
func (e *MountEntry) APIPathNoNamespace() string {
	path := e.Path
	if e.Table == credentialTableType {
		path = credentialRoutePrefix + path
	}
	return path
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

	if len(e.Config.AllowedManagedKeys) == 0 {
		e.synthesizedConfigCache.Delete("allowed_managed_keys")
	} else {
		e.synthesizedConfigCache.Store("allowed_managed_keys", e.Config.AllowedManagedKeys)
	}
}

func (entry *MountEntry) Deserialize() map[string]interface{} {
	return map[string]interface{}{
		"mount_path":      entry.Path,
		"mount_namespace": entry.Namespace().Path,
		"uuid":            entry.UUID,
		"accessor":        entry.Accessor,
		"mount_type":      entry.Type,
	}
}

// DecodeMountTable is used for testing
func (c *Core) DecodeMountTable(ctx context.Context, raw []byte) (*MountTable, error) {
	return c.decodeMountTable(ctx, raw)
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

	// Mount internally
	if err := c.mountInternal(ctx, entry, MountTableUpdateStorage); err != nil {
		return err
	}

	return nil
}

func (c *Core) mountInternal(ctx context.Context, entry *MountEntry, updateStorage bool) error {
	c.mountsLock.Lock()
	c.authLock.Lock()
	locked := true
	unlock := func() {
		if locked {
			c.authLock.Unlock()
			c.mountsLock.Unlock()
			locked = false
		}
	}
	defer unlock()

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

	// Basic check for matching names
	for _, ent := range c.mounts.Entries {
		if ns.ID == ent.NamespaceID {
			switch {
			// Existing is oauth/github/ new is oauth/ or
			// existing is oauth/ and new is oauth/github/
			case strings.HasPrefix(ent.Path, entry.Path):
				fallthrough
			case strings.HasPrefix(entry.Path, ent.Path):
				return logical.CodedError(409, fmt.Sprintf("path is already in use at %s", ent.Path))
			}
		}
	}

	// Verify there are no conflicting mounts in the router
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

	// Resolution to absolute storage paths (versus uuid-relative) needs
	// to happen prior to calling into the forwarded writer. Thus we
	// intercept writes just before they hit barrier storage.
	forwarded, err := c.NewForwardedWriter(ctx, c.barrier, entry.Local)
	if err != nil {
		return fmt.Errorf("error creating forwarded writer: %v", err)
	}

	viewPath := entry.ViewPath()
	view := NewBarrierView(forwarded, viewPath)

	// Singleton mounts cannot be filtered manually on a per-secondary basis
	// from replication.
	if strutil.StrListContains(singletonMounts, entry.Type) {
		addFilterablePath(c, viewPath)
	}
	addKnownPath(c, viewPath)

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

	backend, entry.RunningSha256, err = c.newLogicalBackend(ctx, entry, sysView, view)
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

	// update the entry running version with the configured version, which was verified during registration.
	entry.RunningVersion = entry.Version
	if entry.RunningVersion == "" {
		// don't set the running version to a builtin if it is running as an external plugin
		if entry.RunningSha256 == "" {
			entry.RunningVersion = versions.GetBuiltinVersion(consts.PluginTypeSecrets, entry.Type)
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
	if err = c.entBuiltinPluginMetrics(ctx, entry, 1); err != nil {
		c.logger.Error("failed to emit enabled ent builtin plugin metrics", "error", err)
		return err
	}

	// Re-evaluate filtered paths
	if err := runFilteredPathsEvaluation(ctx, c, false); err != nil {
		c.logger.Error("failed to evaluate filtered paths", "error", err)

		unlock()
		// We failed to evaluate filtered paths so we are undoing the mount operation
		if unmountInternalErr := c.unmountInternal(ctx, entry.Path, MountTableUpdateStorage); unmountInternalErr != nil {
			c.logger.Error("failed to unmount", "error", unmountInternalErr)
		}
		return err
	}

	if !nilMount {
		// restore the original readOnlyErr, so we can write to the view in
		// Initialize() if necessary
		view.setReadOnlyErr(origReadOnlyErr)

		// initialize, using the core's active context.
		nsActiveContext := namespace.ContextWithNamespace(c.activeContext, ns)
		err := backend.Initialize(nsActiveContext, &logical.InitializationRequest{Storage: view})
		if err != nil {
			return err
		}
	}

	if c.logger.IsInfo() {
		c.logger.Info("successful mount", "namespace", entry.Namespace().Path, "path", entry.Path, "type", entry.Type, "version", entry.Version)
	}
	return nil
}

// builtinTypeFromMountEntry attempts to find a builtin PluginType associated
// with the specified MountEntry. Returns consts.PluginTypeUnknown if not found.
func (c *Core) builtinTypeFromMountEntry(ctx context.Context, entry *MountEntry) consts.PluginType {
	if c.builtinRegistry == nil || entry == nil {
		return consts.PluginTypeUnknown
	}

	if !versions.IsBuiltinVersion(entry.RunningVersion) {
		return consts.PluginTypeUnknown
	}

	builtinPluginType := func(name string, pluginType consts.PluginType) (consts.PluginType, bool) {
		plugin, err := c.pluginCatalog.Get(ctx, name, pluginType, entry.RunningVersion)
		if err == nil && plugin != nil && plugin.Builtin {
			return plugin.Type, true
		}
		return consts.PluginTypeUnknown, false
	}

	// auth plugins have their own dedicated mount table
	if pluginType, err := consts.ParsePluginType(entry.Table); err == nil {
		if builtinType, ok := builtinPluginType(entry.Type, pluginType); ok {
			return builtinType
		}
	}

	// Check for possible matches
	var builtinTypes []consts.PluginType
	for _, pluginType := range [...]consts.PluginType{consts.PluginTypeSecrets, consts.PluginTypeDatabase} {
		if builtinType, ok := builtinPluginType(entry.Type, pluginType); ok {
			builtinTypes = append(builtinTypes, builtinType)
		}
	}

	if len(builtinTypes) == 1 {
		return builtinTypes[0]
	}

	return consts.PluginTypeUnknown
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

	// Unmount mount internally
	if err := c.unmountInternal(ctx, path, MountTableUpdateStorage); err != nil {
		return err
	}

	// Re-evaluate filtered paths
	if err := runFilteredPathsEvaluation(ctx, c, true); err != nil {
		// Even we failed to evaluate filtered paths, the unmount operation was still successful
		c.logger.Error("failed to evaluate filtered paths", "error", err)
	}
	return nil
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
	if err := c.taintMountEntry(ctx, ns.ID, path, updateStorage, true); err != nil {
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
		// Invoke the rollback manager a final time. This is not fatal as
		// various periodic funcs (e.g., PKI) can legitimately error; the
		// periodic rollback manager logs these errors rather than failing
		// replication like returning this error would do.
		if err := c.rollback.Rollback(rCtx, path); err != nil {
			c.logger.Error("ignoring rollback error during unmount", "error", err, "path", path)
			err = nil
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

	viewPath := entry.ViewPath()
	switch {
	case !updateStorage:
		// Don't attempt to clear data, replication will handle this
	case c.IsDRSecondary():
		// If we are a dr secondary we want to clear the view, but the provided
		// view is marked as read only. We use the barrier here to get around
		// it.
		if err := logical.ClearViewWithLogging(ctx, NewBarrierView(c.barrier, viewPath), c.logger.Named("secrets.deletion").With("namespace", ns.ID, "path", path)); err != nil {
			c.logger.Error("failed to clear view for path being unmounted", "error", err, "path", path)
			return err
		}

	case entry.Local, !c.IsPerfSecondary():
		// Have writable storage, remove the whole thing
		if err := logical.ClearViewWithLogging(ctx, view, c.logger.Named("secrets.deletion").With("namespace", ns.ID, "path", path)); err != nil {
			c.logger.Error("failed to clear view for path being unmounted", "error", err, "path", path)
			return err
		}

	case !entry.Local && c.IsPerfSecondary():
		if err := clearIgnoredPaths(ctx, c, backend, viewPath); err != nil {
			return err
		}
	}

	// Remove the mount table entry
	if err := c.removeMountEntry(ctx, path, updateStorage); err != nil {
		c.logger.Error("failed to remove mount entry for path being unmounted", "error", err, "path", path)
		return err
	}

	// Unmount the backend entirely
	if err := c.router.Unmount(ctx, path); err != nil {
		return err
	}
	if err = c.entBuiltinPluginMetrics(ctx, entry, -1); err != nil {
		c.logger.Error("failed to emit disabled ent builtin plugin metrics", "error", err)
		return err
	}

	removePathCheckers(c, entry, viewPath)

	if c.quotaManager != nil {
		if err := c.quotaManager.HandleBackendDisabling(ctx, ns.Path, path); err != nil {
			c.logger.Error("failed to update quotas after disabling mount", "path", path, "error", err)
			return err
		}
	}

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
func (c *Core) taintMountEntry(ctx context.Context, nsID, mountPath string, updateStorage, unmounting bool) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	mountState := ""
	if unmounting {
		mountState = mountStateUnmounting
	}

	// As modifying the taint of an entry affects shallow clones,
	// we simply use the original
	entry, err := c.mounts.setTaint(nsID, mountPath, true, mountState)
	if err != nil {
		return err
	}
	if entry == nil {
		c.logger.Error("nil entry found tainting entry in mounts table", "path", mountPath)
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

// handleDeprecatedMountEntry handles the Deprecation Status of the specified
// mount entry's builtin engine. Warnings are appended to the returned response
// and logged. Errors are returned with a nil response to be processed by the
// caller.
func (c *Core) handleDeprecatedMountEntry(ctx context.Context, entry *MountEntry, pluginType consts.PluginType) (*logical.Response, error) {
	resp := &logical.Response{}

	if c.builtinRegistry == nil || entry == nil {
		return nil, nil
	}

	// Allow type to be determined from mount entry when not otherwise specified
	if pluginType == consts.PluginTypeUnknown {
		pluginType = c.builtinTypeFromMountEntry(ctx, entry)
	}

	// Handle aliases
	t := entry.Type
	if alias, ok := mountAliases[t]; ok {
		t = alias
	}

	status, ok := c.builtinRegistry.DeprecationStatus(t, pluginType)
	if ok {
		switch status {
		case consts.Deprecated:
			c.logger.Warn("mounting deprecated builtin", "name", t, "type", pluginType, "path", entry.Path)
			resp.AddWarning(errMountDeprecated.Error())
			return resp, nil

		case consts.PendingRemoval:
			if c.pendingRemovalMountsAllowed {
				c.Logger().Info("mount allowed by environment variable", "env", consts.EnvVaultAllowPendingRemovalMounts)
				resp.AddWarning(errMountPendingRemoval.Error())
				return resp, nil
			}
			return nil, errMountPendingRemoval

		case consts.Removed:
			return nil, errMountRemoved
		}
	}
	return nil, nil
}

// remountForceInternal takes a copy of the mount entry for the path and fully unmounts
// and remounts the backend to pick up any changes, such as filtered paths.
// Should be only used for internal usage.
func (c *Core) remountForceInternal(ctx context.Context, path string, updateStorage bool) error {
	me := c.router.MatchingMountEntry(ctx, path)
	if me == nil {
		return fmt.Errorf("cannot find mount for path %q", path)
	}

	me, err := me.Clone()
	if err != nil {
		return err
	}

	if err := c.unmountInternal(ctx, path, updateStorage); err != nil {
		return err
	}

	// Mount internally
	if err := c.mountInternal(ctx, me, updateStorage); err != nil {
		return err
	}

	return nil
}

func (c *Core) remountSecretsEngineCurrentNamespace(ctx context.Context, src, dst string, updateStorage bool) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	srcPathDetails := c.splitNamespaceAndMountFromPath(ns.Path, src)
	dstPathDetails := c.splitNamespaceAndMountFromPath(ns.Path, dst)
	return c.remountSecretsEngine(ctx, srcPathDetails, dstPathDetails, updateStorage)
}

// remountSecretsEngine is used to remount a path at a new mount point.
func (c *Core) remountSecretsEngine(ctx context.Context, src, dst namespace.MountPathDetails, updateStorage bool) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	// Prevent protected paths from being remounted, or target mounts being in protected paths
	for _, p := range protectedMounts {
		if strings.HasPrefix(src.MountPath, p) {
			return fmt.Errorf("cannot remount %q", src.MountPath)
		}

		if strings.HasPrefix(dst.MountPath, p) {
			return fmt.Errorf("cannot remount to destination %+v", dst)
		}
	}

	srcRelativePath := src.GetRelativePath(ns)
	dstRelativePath := dst.GetRelativePath(ns)

	// Verify exact match of the route
	srcMatch := c.router.MatchingMountEntry(ctx, srcRelativePath)
	if srcMatch == nil {
		return fmt.Errorf("no matching mount at %q", src.Namespace.Path+src.MountPath)
	}

	if match := c.router.MountConflict(ctx, dstRelativePath); match != "" {
		return fmt.Errorf("path in use at %q", match)
	}

	// Mark the entry as tainted
	if err := c.taintMountEntry(ctx, src.Namespace.ID, src.MountPath, updateStorage, false); err != nil {
		return err
	}

	// Taint the router path to prevent routing
	if err := c.router.Taint(ctx, srcRelativePath); err != nil {
		return err
	}

	if !c.IsDRSecondary() {
		// Invoke the rollback manager a final time. This is not fatal as
		// various periodic funcs (e.g., PKI) can legitimately error; the
		// periodic rollback manager logs these errors rather than failing
		// replication like returning this error would do.
		rCtx := namespace.ContextWithNamespace(c.activeContext, ns)
		if c.rollback != nil && c.router.MatchingBackend(ctx, srcRelativePath) != nil {
			if err := c.rollback.Rollback(rCtx, srcRelativePath); err != nil {
				c.logger.Error("ignoring rollback error during remount", "error", err, "path", src.Namespace.Path+src.MountPath)
				err = nil
			}
		}

		revokeCtx := namespace.ContextWithNamespace(ctx, src.Namespace)
		// Revoke all the dynamic keys
		if err := c.expiration.RevokePrefix(revokeCtx, src.MountPath, true); err != nil {
			return err
		}
	}

	c.mountsLock.Lock()
	if match := c.router.MountConflict(ctx, dstRelativePath); match != "" {
		c.mountsLock.Unlock()
		return fmt.Errorf("path in use at %q", match)
	}

	srcMatch.Tainted = false
	srcMatch.NamespaceID = dst.Namespace.ID
	srcMatch.namespace = dst.Namespace
	srcPath := srcMatch.Path
	srcMatch.Path = dst.MountPath

	// Update the mount table
	if err := c.persistMounts(ctx, c.mounts, &srcMatch.Local); err != nil {
		srcMatch.Path = srcPath
		srcMatch.Tainted = true
		c.mountsLock.Unlock()
		if err == logical.ErrReadOnly && c.perfStandby {
			return err
		}

		return fmt.Errorf("failed to update mount table with error %+v", err)
	}

	// Remount the backend
	if err := c.router.Remount(ctx, srcRelativePath, dstRelativePath); err != nil {
		c.mountsLock.Unlock()
		return err
	}
	c.mountsLock.Unlock()

	// Un-taint the path
	if err := c.router.Untaint(ctx, dstRelativePath); err != nil {
		return err
	}

	return nil
}

// From an input path that has a relative namespace hierarchy followed by a mount point, return the full
// namespace of the mount point, along with the mount point without the namespace related prefix.
// For example, in a hierarchy ns1/ns2/ns3/secret-mount, when currNs is ns1 and path is ns2/ns3/secret-mount,
// this returns the namespace object for ns1/ns2/ns3/, and the string "secret-mount"
func (c *Core) splitNamespaceAndMountFromPath(currNs, path string) namespace.MountPathDetails {
	fullPath := currNs + path
	fullNs := c.namespaceByPath(fullPath)

	mountPath := strings.TrimPrefix(fullPath, fullNs.Path)

	return namespace.MountPathDetails{
		Namespace: fullNs,
		MountPath: sanitizePath(mountPath),
	}
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
		c.tableMetrics(len(mountTable.Entries), false, false, raw.Value)
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
			c.tableMetrics(len(localMountTable.Entries), true, false, rawLocal.Value)
			c.mounts.Entries = append(c.mounts.Entries, localMountTable.Entries...)
		}
	}

	// If this node is a performance standby we do not want to attempt to
	// upgrade the mount table, this will be the active node's responsibility.
	if !c.perfStandby {
		err := c.runMountUpdates(ctx, needPersist)
		if err != nil {
			c.logger.Error("failed to run mount table upgrades", "error", err)
			return err
		}
	}

	for _, entry := range c.mounts.Entries {
		if entry.NamespaceID == "" {
			entry.NamespaceID = namespace.RootNamespaceID
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
	return nil
}

// Note that this is only designed to work with singletons, as it checks by
// type only.
func (c *Core) runMountUpdates(ctx context.Context, needPersist bool) error {
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
		if !foundRequired && (!c.IsPerfSecondary() || requiredMount.Local) {
			c.mounts.Entries = append(c.mounts.Entries, requiredMount)
			needPersist = true
		}
	}

	// Upgrade to table-scoped entries
	for _, entry := range c.mounts.Entries {
		if !c.PR1103disabled && entry.Type == mountTypeNSCubbyhole && !entry.Local && !c.ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationDRSecondary) {
			entry.Local = true
			needPersist = true
		}
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

		// Don't store built-in version in the mount table, to make upgrades smoother.
		if versions.IsBuiltinVersion(entry.Version) {
			entry.Version = ""
			needPersist = true
		}
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

	nonLocalMounts := &MountTable{
		Type: mountTableType,
	}

	localMounts := &MountTable{
		Type: mountTableType,
	}

	for _, entry := range table.Entries {
		if entry.Table != table.Type {
			c.logger.Error("given entry to persist in mount table has wrong table value", "path", entry.Path, "entry_table_type", entry.Table, "actual_type", table.Type)
			return fmt.Errorf("invalid mount entry found, not persisting")
		}

		if entry.Local {
			localMounts.Entries = append(localMounts.Entries, entry)
		} else {
			nonLocalMounts.Entries = append(nonLocalMounts.Entries, entry)
		}
	}

	writeTable := func(mt *MountTable, path string) ([]byte, error) {
		// Encode the mount table into JSON and compress it (lzw).
		compressedBytes, err := jsonutil.EncodeJSONAndCompress(mt, nil)
		if err != nil {
			c.logger.Error("failed to encode or compress mount table", "error", err)
			return nil, err
		}

		// Create an entry
		entry := &logical.StorageEntry{
			Key:   path,
			Value: compressedBytes,
		}

		// Write to the physical backend
		if err := c.barrier.Put(ctx, entry); err != nil {
			c.logger.Error("failed to persist mount table", "error", err)
			return nil, err
		}
		return compressedBytes, nil
	}

	var err error
	var compressedBytes []byte
	switch {
	case local == nil:
		// Write non-local mounts
		compressedBytes, err := writeTable(nonLocalMounts, coreMountConfigPath)
		if err != nil {
			return err
		}
		c.tableMetrics(len(nonLocalMounts.Entries), false, false, compressedBytes)

		// Write local mounts
		compressedBytes, err = writeTable(localMounts, coreLocalMountConfigPath)
		if err != nil {
			return err
		}
		c.tableMetrics(len(localMounts.Entries), true, false, compressedBytes)

	case *local:
		// Write local mounts
		compressedBytes, err = writeTable(localMounts, coreLocalMountConfigPath)
		if err != nil {
			return err
		}
		c.tableMetrics(len(localMounts.Entries), true, false, compressedBytes)
	default:
		// Write non-local mounts
		compressedBytes, err = writeTable(nonLocalMounts, coreMountConfigPath)
		if err != nil {
			return err
		}
		c.tableMetrics(len(nonLocalMounts.Entries), false, false, compressedBytes)
	}

	return nil
}

// setupMounts is invoked after we've loaded the mount table to
// initialize the logical backends and setup the router
func (c *Core) setupMounts(ctx context.Context) error {
	c.mountsLock.Lock()
	defer c.mountsLock.Unlock()

	for _, entry := range c.mounts.sortEntriesByPathDepth().Entries {
		// Initialize the backend, special casing for system
		barrierPath := entry.ViewPath()

		// Resolution to absolute storage paths (versus uuid-relative) needs
		// to happen prior to calling into the forwarded writer. Thus we
		// intercept writes just before they hit barrier storage.
		forwarded, err := c.NewForwardedWriter(ctx, c.barrier, entry.Local)
		if err != nil {
			return fmt.Errorf("error creating forwarded writer: %v", err)
		}

		// Create a barrier storage view using the UUID
		view := NewBarrierView(forwarded, barrierPath)

		// Singleton mounts cannot be filtered manually on a per-secondary basis
		// from replication
		if strutil.StrListContains(singletonMounts, entry.Type) {
			addFilterablePath(c, barrierPath)
		}
		addKnownPath(c, barrierPath)

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
		}

		var backend logical.Backend
		// Create the new backend
		sysView := c.mountEntrySysView(entry)
		backend, entry.RunningSha256, err = c.newLogicalBackend(ctx, entry, sysView, view)
		if err != nil {
			c.logger.Error("failed to create mount entry", "path", entry.Path, "error", err)

			if c.isMountable(ctx, entry, consts.PluginTypeSecrets) {
				c.logger.Warn("skipping plugin-based mount entry", "path", entry.Path)
				goto ROUTER_MOUNT
			}
			return errLoadMountsFailed
		}
		if backend == nil {
			return fmt.Errorf("created mount entry of type %q is nil", entry.Type)
		}

		// update the entry running version with the configured version, which was verified during registration.
		entry.RunningVersion = entry.Version
		if entry.RunningVersion == "" {
			// don't set the running version to a builtin if it is running as an external plugin
			if entry.RunningSha256 == "" {
				entry.RunningVersion = versions.GetBuiltinVersion(consts.PluginTypeSecrets, entry.Type)
			}
		}

		// Do not start up deprecated builtin plugins. If this is a major
		// upgrade, stop unsealing and shutdown. If we've already mounted this
		// plugin, proceed with unsealing and skip backend initialization.
		if versions.IsBuiltinVersion(entry.RunningVersion) {
			_, err := c.handleDeprecatedMountEntry(ctx, entry, consts.PluginTypeSecrets)
			if c.isMajorVersionFirstMount(ctx) && err != nil {
				go c.ShutdownCoreError(fmt.Errorf("could not mount %q: %w", entry.Type, err))
				return errLoadMountsFailed
			} else if err != nil {
				c.logger.Error("skipping deprecated mount entry", "name", entry.Type, "path", entry.Path, "error", err)
				backend.Cleanup(ctx)
				backend = nil
				goto ROUTER_MOUNT
			}
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

		// Initialize
		if !nilMount {
			// Bind locally
			localEntry := entry
			c.postUnsealFuncs = append(c.postUnsealFuncs, func() {
				postUnsealLogger := c.logger.With("type", localEntry.Type, "version", localEntry.RunningVersion, "path", localEntry.Path)
				if backend == nil {
					postUnsealLogger.Error("skipping initialization for nil backend", "path", localEntry.Path)
					return
				}
				if !strutil.StrListContains(singletonMounts, localEntry.Type) {
					view.setReadOnlyErr(origReadOnlyErr)
				}

				nsActiveContext := namespace.ContextWithNamespace(c.activeContext, localEntry.Namespace())
				err := backend.Initialize(nsActiveContext, &logical.InitializationRequest{Storage: view})
				if err != nil {
					postUnsealLogger.Error("failed to initialize mount backend", "error", err)
				}
			})
		}

		if c.logger.IsInfo() {
			c.logger.Info("successfully mounted", "type", entry.Type, "version", entry.RunningVersion, "path", entry.Path, "namespace", entry.Namespace())
		}

		// Ensure the path is tainted if set in the mount table
		if entry.Tainted {
			// Calculate any namespace prefixes here, because when Taint() is called, there won't be
			// a namespace to pull from the context. This is similar to what we do above in c.router.Mount().
			path := entry.Namespace().Path + entry.Path
			c.logger.Debug("tainting a mount due to it being marked as tainted in mount table", "entry.path", entry.Path, "entry.namespace.path", entry.Namespace().Path, "full_path", path)
			c.router.Taint(ctx, path)
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
	c.router.reset()
	c.systemBarrierView = nil
	return nil
}

// newLogicalBackend is used to create and configure a new logical backend by name.
// It also returns the SHA256 of the plugin, if available.
func (c *Core) newLogicalBackend(ctx context.Context, entry *MountEntry, sysView logical.SystemView, view logical.Storage) (logical.Backend, string, error) {
	t := entry.Type
	if alias, ok := mountAliases[t]; ok {
		t = alias
	}

	var runningSha string
	f, ok := c.logicalBackends[t]
	if !ok {
		plug, err := c.pluginCatalog.Get(ctx, t, consts.PluginTypeSecrets, entry.Version)
		if err != nil {
			return nil, "", err
		}
		if plug == nil {
			errContext := t
			if entry.Version != "" {
				errContext += fmt.Sprintf(", version=%s", entry.Version)
			}
			return nil, "", fmt.Errorf("%w: %s", ErrPluginNotFound, errContext)
		}
		if len(plug.Sha256) > 0 {
			runningSha = hex.EncodeToString(plug.Sha256)
		}

		f = plugin.Factory
		if !plug.Builtin {
			f = wrapFactoryCheckPerms(c, plugin.Factory)
		}
	}
	// Set up conf to pass in plugin_name
	conf := make(map[string]string)
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
	conf["plugin_version"] = entry.Version

	backendLogger := c.baseLogger.Named(fmt.Sprintf("secrets.%s.%s", t, entry.Accessor))
	pluginEventSender, err := c.events.WithPlugin(entry.namespace, &logical.EventPluginInfo{
		MountClass:    consts.PluginTypeSecrets.String(),
		MountAccessor: entry.Accessor,
		MountPath:     entry.Path,
		Plugin:        entry.Type,
		PluginVersion: entry.RunningVersion,
		Version:       entry.Version,
	})
	if err != nil {
		return nil, "", err
	}
	config := &logical.BackendConfig{
		StorageView: view,
		Logger:      backendLogger,
		Config:      conf,
		System:      sysView,
		BackendUUID: entry.BackendAwareUUID,
	}
	if c.IsExperimentEnabled(experiments.VaultExperimentEventsAlpha1) {
		config.EventsSender = pluginEventSender
	}

	ctx = namespace.ContextWithNamespace(ctx, entry.namespace)
	ctx = context.WithValue(ctx, "core_number", c.coreNumber)
	b, err := f(ctx, config)
	if err != nil {
		return nil, "", err
	}
	if b == nil {
		return nil, "", fmt.Errorf("nil backend of type %q returned from factory", t)
	}
	addLicenseCallback(c, b)

	return b, runningSha, nil
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
			RunningVersion: versions.GetBuiltinVersion(consts.PluginTypeSecrets, "kv"),
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
		RunningVersion:   versions.GetBuiltinVersion(consts.PluginTypeSecrets, "cubbyhole"),
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
		SealWrap:         true, // Enable SealWrap since SystemBackend utilizes SealWrapStorage, see factory in addExtraLogicalBackends().
		Config: MountConfig{
			PassthroughRequestHeaders: []string{"Accept"},
		},
		RunningVersion: versions.DefaultBuiltinVersion,
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
		Config: MountConfig{
			PassthroughRequestHeaders: []string{"Authorization"},
		},
		RunningVersion: versions.DefaultBuiltinVersion,
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

func (c *Core) createMigrationStatus(from, to namespace.MountPathDetails) (string, error) {
	migrationID, err := uuid.GenerateUUID()
	if err != nil {
		return "", fmt.Errorf("error generating uuid for mount move invocation: %w", err)
	}
	migrationInfo := MountMigrationInfo{
		SourceMount:     from.Namespace.Path + from.MountPath,
		TargetMount:     to.Namespace.Path + to.MountPath,
		MigrationStatus: MigrationInProgressStatus.String(),
	}
	c.mountMigrationTracker.Store(migrationID, migrationInfo)
	return migrationID, nil
}

func (c *Core) setMigrationStatus(migrationID string, migrationStatus MountMigrationStatus) error {
	migrationInfoRaw, ok := c.mountMigrationTracker.Load(migrationID)
	if !ok {
		return fmt.Errorf("Migration Tracker entry missing for ID %s", migrationID)
	}
	migrationInfo := migrationInfoRaw.(MountMigrationInfo)
	migrationInfo.MigrationStatus = migrationStatus.String()
	c.mountMigrationTracker.Store(migrationID, migrationInfo)
	return nil
}

func (c *Core) readMigrationStatus(migrationID string) *MountMigrationInfo {
	migrationInfoRaw, ok := c.mountMigrationTracker.Load(migrationID)
	if !ok {
		return nil
	}
	migrationInfo := migrationInfoRaw.(MountMigrationInfo)
	return &migrationInfo
}
