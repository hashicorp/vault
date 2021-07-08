package vault

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	memdb "github.com/hashicorp/go-memdb"
	"github.com/hashicorp/go-multierror"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/hostutil"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/monitor"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/random"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/mitchellh/mapstructure"
)

const (
	maxBytes    = 128 * 1024
	globalScope = "global"
)

func systemBackendMemDBSchema() *memdb.DBSchema {
	systemSchema := &memdb.DBSchema{
		Tables: make(map[string]*memdb.TableSchema),
	}

	schemas := getSystemSchemas()

	for _, schemaFunc := range schemas {
		schema := schemaFunc()
		if _, ok := systemSchema.Tables[schema.Name]; ok {
			panic(fmt.Sprintf("duplicate table name: %s", schema.Name))
		}
		systemSchema.Tables[schema.Name] = schema
	}

	return systemSchema
}

func NewSystemBackend(core *Core, logger log.Logger) *SystemBackend {
	db, _ := memdb.NewMemDB(systemBackendMemDBSchema())

	b := &SystemBackend{
		Core:      core,
		db:        db,
		logger:    logger,
		mfaLogger: core.baseLogger.Named("mfa"),
		mfaLock:   &sync.RWMutex{},
	}

	core.AddLogger(b.mfaLogger)

	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(sysHelpRoot),

		PathsSpecial: &logical.Paths{
			Root: []string{
				"auth/*",
				"remount",
				"audit",
				"audit/*",
				"raw",
				"raw/*",
				"replication/primary/secondary-token",
				"replication/performance/primary/secondary-token",
				"replication/dr/primary/secondary-token",
				"replication/reindex",
				"replication/dr/reindex",
				"replication/performance/reindex",
				"rotate",
				"config/cors",
				"config/auditing/*",
				"config/ui/headers/*",
				"plugins/catalog/*",
				"revoke-prefix/*",
				"revoke-force/*",
				"leases/revoke-prefix/*",
				"leases/revoke-force/*",
				"leases/lookup/*",
				"storage/raft/snapshot-auto/config/*",
				"leases",
			},

			Unauthenticated: []string{
				"wrapping/lookup",
				"wrapping/pubkey",
				"replication/status",
				"internal/specs/openapi",
				"internal/ui/mounts",
				"internal/ui/mounts/*",
				"internal/ui/namespaces",
				"replication/performance/status",
				"replication/dr/status",
				"replication/dr/secondary/promote",
				"replication/dr/secondary/disable",
				"replication/dr/secondary/recover",
				"replication/dr/secondary/update-primary",
				"replication/dr/secondary/operation-token/delete",
				"replication/dr/secondary/license",
				"replication/dr/secondary/license/signed",
				"replication/dr/secondary/license/status",
				"replication/dr/secondary/sys/config/reload/license",
				"replication/dr/secondary/reindex",
				"storage/raft/bootstrap/challenge",
				"storage/raft/bootstrap/answer",
				"init",
				"seal-status",
				"unseal",
				"leader",
				"health",
				"generate-root/attempt",
				"generate-root/update",
				"rekey/init",
				"rekey/update",
				"rekey/verify",
				"rekey-recovery-key/init",
				"rekey-recovery-key/update",
				"rekey-recovery-key/verify",
			},

			LocalStorage: []string{
				expirationSubPath,
				countersSubPath,
			},
		},
	}

	b.Backend.Paths = append(b.Backend.Paths, entPaths(b)...)
	b.Backend.Paths = append(b.Backend.Paths, b.configPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.rekeyPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.sealPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.statusPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.pluginsCatalogListPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.pluginsCatalogCRUDPath())
	b.Backend.Paths = append(b.Backend.Paths, b.pluginsReloadPath())
	b.Backend.Paths = append(b.Backend.Paths, b.auditPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.mountPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.authPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.leasePaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.policyPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.wrappingPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.toolsPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.capabilitiesPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.internalPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.pprofPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.remountPath())
	b.Backend.Paths = append(b.Backend.Paths, b.metricsPath())
	b.Backend.Paths = append(b.Backend.Paths, b.monitorPath())
	b.Backend.Paths = append(b.Backend.Paths, b.hostInfoPath())
	b.Backend.Paths = append(b.Backend.Paths, b.quotasPaths()...)
	b.Backend.Paths = append(b.Backend.Paths, b.rootActivityPaths()...)

	if core.rawEnabled {
		b.Backend.Paths = append(b.Backend.Paths, b.rawPaths()...)
	}

	if backend := core.getRaftBackend(); backend != nil {
		b.Backend.Paths = append(b.Backend.Paths, b.raftStoragePaths()...)
	}

	// If the node is in a DR secondary cluster, gate some raft operations by
	// the DR operation token.
	if core.IsDRSecondary() {
		b.Backend.PathsSpecial.Unauthenticated = append(b.Backend.PathsSpecial.Unauthenticated, "storage/raft/autopilot/configuration")
		b.Backend.PathsSpecial.Unauthenticated = append(b.Backend.PathsSpecial.Unauthenticated, "storage/raft/autopilot/state")
		b.Backend.PathsSpecial.Unauthenticated = append(b.Backend.PathsSpecial.Unauthenticated, "storage/raft/configuration")
		b.Backend.PathsSpecial.Unauthenticated = append(b.Backend.PathsSpecial.Unauthenticated, "storage/raft/remove-peer")
	}

	b.Backend.Invalidate = sysInvalidate(b)
	return b
}

func (b *SystemBackend) rawPaths() []*framework.Path {
	r := &RawBackend{
		barrier: b.Core.barrier,
		logger:  b.logger,
		checkRaw: func(path string) error {
			return checkRaw(b, path)
		},
	}
	return rawPaths("", r)
}

// SystemBackend implements logical.Backend and is used to interact with
// the core of the system. This backend is hardcoded to exist at the "sys"
// prefix. Conceptually it is similar to procfs on Linux.
type SystemBackend struct {
	*framework.Backend
	Core      *Core
	db        *memdb.MemDB
	mfaLock   *sync.RWMutex
	mfaLogger log.Logger
	logger    log.Logger
}

// handleConfigStateSanitized returns the current configuration state. The configuration
// data that it returns is a sanitized version of the combined configuration
// file(s) provided.
func (b *SystemBackend) handleConfigStateSanitized(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config := b.Core.SanitizedConfig()
	resp := &logical.Response{
		Data: config,
	}
	return resp, nil
}

// handleConfigReload handles reloading specific pieces of the configuration.
func (b *SystemBackend) handleConfigReload(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	subsystem := data.Get("subsystem").(string)

	switch subsystem {
	case "license":
		return handleLicenseReload(b)(ctx, req, data)
	}

	return nil, logical.ErrUnsupportedPath
}

// handleCORSRead returns the current CORS configuration
func (b *SystemBackend) handleCORSRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	corsConf := b.Core.corsConfig

	enabled := corsConf.IsEnabled()

	resp := &logical.Response{
		Data: map[string]interface{}{
			"enabled": enabled,
		},
	}

	if enabled {
		corsConf.RLock()
		resp.Data["allowed_origins"] = corsConf.AllowedOrigins
		resp.Data["allowed_headers"] = corsConf.AllowedHeaders
		corsConf.RUnlock()
	}

	return resp, nil
}

// handleCORSUpdate sets the list of origins that are allowed to make
// cross-origin requests and sets the CORS enabled flag to true
func (b *SystemBackend) handleCORSUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	origins := d.Get("allowed_origins").([]string)
	headers := d.Get("allowed_headers").([]string)

	return nil, b.Core.corsConfig.Enable(ctx, origins, headers)
}

// handleCORSDelete sets the CORS enabled flag to false and clears the list of
// allowed origins & headers.
func (b *SystemBackend) handleCORSDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return nil, b.Core.corsConfig.Disable(ctx)
}

func (b *SystemBackend) handleTidyLeases(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		tidyCtx := namespace.ContextWithNamespace(b.Core.activeContext, ns)
		err := b.Core.expiration.Tidy(tidyCtx)
		if err != nil {
			b.Backend.Logger().Error("failed to tidy leases", "error", err)
			return
		}
	}()

	resp := &logical.Response{}
	resp.AddWarning("Tidy operation successfully started. Any information from the operation will be printed to Vault's server logs.")
	return logical.RespondWithStatusCode(resp, req, http.StatusAccepted)
}

func (b *SystemBackend) handleLeaseCount(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	typeRaw, ok := d.GetOk("type")
	if !ok || strings.ToLower(typeRaw.(string)) != "irrevocable" {
		return nil, nil
	}

	includeChildNamespacesRaw, ok := d.GetOk("include_child_namespaces")
	includeChildNamespaces := ok && includeChildNamespacesRaw.(bool)

	resp, err := b.Core.expiration.getIrrevocableLeaseCounts(ctx, includeChildNamespaces)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: resp,
	}, nil
}

func processLimit(d *framework.FieldData) (bool, int, error) {
	limitStr := ""
	limitRaw, ok := d.GetOk("limit")
	if ok {
		limitStr = limitRaw.(string)
	}

	includeAll := false
	maxResults := MaxIrrevocableLeasesToReturn
	if limitStr == "" {
		// use the defaults
	} else if strings.ToLower(limitStr) == "none" {
		includeAll = true
	} else {
		// not having a valid, positive int here is an error
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return false, 0, fmt.Errorf("invalid 'limit' provided: %w", err)
		}

		if limit < 1 {
			return false, 0, fmt.Errorf("limit must be 'none' or a positive integer")
		}

		maxResults = limit
	}

	return includeAll, maxResults, nil
}

func (b *SystemBackend) handleLeaseList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	typeRaw, ok := d.GetOk("type")
	if !ok || strings.ToLower(typeRaw.(string)) != "irrevocable" {
		return nil, nil
	}

	includeChildNamespacesRaw, ok := d.GetOk("include_child_namespaces")
	includeChildNamespaces := ok && includeChildNamespacesRaw.(bool)

	includeAll, maxResults, err := processLimit(d)
	if err != nil {
		return nil, err
	}

	leases, warning, err := b.Core.expiration.listIrrevocableLeases(ctx, includeChildNamespaces, includeAll, maxResults)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: leases,
	}
	if warning != "" {
		resp.AddWarning(warning)
	}

	return resp, nil
}

func (b *SystemBackend) handlePluginCatalogTypedList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	pluginType, err := consts.ParsePluginType(d.Get("type").(string))
	if err != nil {
		return nil, err
	}

	plugins, err := b.Core.pluginCatalog.List(ctx, pluginType)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(plugins), nil
}

func (b *SystemBackend) handlePluginCatalogUntypedList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	pluginsByType := make(map[string]interface{})
	for _, pluginType := range consts.PluginTypes {
		plugins, err := b.Core.pluginCatalog.List(ctx, pluginType)
		if err != nil {
			return nil, err
		}
		if len(plugins) > 0 {
			sort.Strings(plugins)
			pluginsByType[pluginType.String()] = plugins
		}
	}
	return &logical.Response{
		Data: pluginsByType,
	}, nil
}

func (b *SystemBackend) handlePluginCatalogUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	pluginName := d.Get("name").(string)
	if pluginName == "" {
		return logical.ErrorResponse("missing plugin name"), nil
	}

	pluginTypeStr := d.Get("type").(string)
	if pluginTypeStr == "" {
		// If the plugin type is not provided, list it as unknown so that we
		// add it to the catalog and UpdatePlugins later will sort it.
		pluginTypeStr = "unknown"
	}
	pluginType, err := consts.ParsePluginType(pluginTypeStr)
	if err != nil {
		return nil, err
	}

	sha256 := d.Get("sha256").(string)
	if sha256 == "" {
		sha256 = d.Get("sha_256").(string)
		if sha256 == "" {
			return logical.ErrorResponse("missing SHA-256 value"), nil
		}
	}

	command := d.Get("command").(string)
	if command == "" {
		return logical.ErrorResponse("missing command value"), nil
	}

	// For backwards compatibility, also accept args as part of command. Don't
	// accepts args in both command and args.
	args := d.Get("args").([]string)
	parts := strings.Split(command, " ")
	if len(parts) <= 0 {
		return logical.ErrorResponse("missing command value"), nil
	} else if len(parts) > 1 && len(args) > 0 {
		return logical.ErrorResponse("must not specify args in command and args field"), nil
	} else if len(parts) > 1 {
		args = parts[1:]
	}

	env := d.Get("env").([]string)

	sha256Bytes, err := hex.DecodeString(sha256)
	if err != nil {
		return logical.ErrorResponse("Could not decode SHA-256 value from Hex"), err
	}

	err = b.Core.pluginCatalog.Set(ctx, pluginName, pluginType, parts[0], args, env, sha256Bytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *SystemBackend) handlePluginCatalogRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	pluginName := d.Get("name").(string)
	if pluginName == "" {
		return logical.ErrorResponse("missing plugin name"), nil
	}

	pluginTypeStr := d.Get("type").(string)
	if pluginTypeStr == "" {
		// If the plugin type is not provided (i.e. the old
		// sys/plugins/catalog/:name endpoint is being requested) short-circuit here
		// and return a warning
		resp := &logical.Response{}
		resp.AddWarning(fmt.Sprintf("Deprecated API endpoint, cannot read plugin information from catalog for %q", pluginName))
		return resp, nil
	}

	pluginType, err := consts.ParsePluginType(pluginTypeStr)
	if err != nil {
		return nil, err
	}

	plugin, err := b.Core.pluginCatalog.Get(ctx, pluginName, pluginType)
	if err != nil {
		return nil, err
	}
	if plugin == nil {
		return nil, nil
	}

	command := ""
	if !plugin.Builtin {
		command, err = filepath.Rel(b.Core.pluginCatalog.directory, plugin.Command)
		if err != nil {
			return nil, err
		}
	}

	data := map[string]interface{}{
		"name":    plugin.Name,
		"args":    plugin.Args,
		"command": command,
		"sha256":  hex.EncodeToString(plugin.Sha256),
		"builtin": plugin.Builtin,
	}

	return &logical.Response{
		Data: data,
	}, nil
}

func (b *SystemBackend) handlePluginCatalogDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	pluginName := d.Get("name").(string)
	if pluginName == "" {
		return logical.ErrorResponse("missing plugin name"), nil
	}

	var resp *logical.Response
	pluginTypeStr := d.Get("type").(string)
	if pluginTypeStr == "" {
		// If the plugin type is not provided (i.e. the old
		// sys/plugins/catalog/:name endpoint is being requested), set type to
		// unknown and let pluginCatalog.Delete proceed. It should handle
		// deregistering out of the old storage path (root of core/plugin-catalog)
		resp = new(logical.Response)
		resp.AddWarning(fmt.Sprintf("Deprecated API endpoint, cannot deregister plugin from catalog for %q", pluginName))
		pluginTypeStr = "unknown"
	}

	pluginType, err := consts.ParsePluginType(pluginTypeStr)
	if err != nil {
		return nil, err
	}
	if err := b.Core.pluginCatalog.Delete(ctx, pluginName, pluginType); err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *SystemBackend) handlePluginReloadUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	pluginName := d.Get("plugin").(string)
	pluginMounts := d.Get("mounts").([]string)
	scope := d.Get("scope").(string)

	if scope != "" && scope != globalScope {
		return logical.ErrorResponse("reload scope must be omitted or 'global'"), nil
	}

	if pluginName != "" && len(pluginMounts) > 0 {
		return logical.ErrorResponse("plugin and mounts cannot be set at the same time"), nil
	}
	if pluginName == "" && len(pluginMounts) == 0 {
		return logical.ErrorResponse("plugin or mounts must be provided"), nil
	}

	if pluginName != "" {
		err := b.Core.reloadMatchingPlugin(ctx, pluginName)
		if err != nil {
			return nil, err
		}
	} else if len(pluginMounts) > 0 {
		err := b.Core.reloadMatchingPluginMounts(ctx, pluginMounts)
		if err != nil {
			return nil, err
		}
	}

	r := logical.Response{
		Data: map[string]interface{}{
			"reload_id": req.ID,
		},
	}

	if scope == globalScope {
		err := handleGlobalPluginReload(ctx, b.Core, req.ID, pluginName, pluginMounts)
		if err != nil {
			return nil, err
		}
		return logical.RespondWithStatusCode(&r, req, http.StatusAccepted)
	}
	return &r, nil
}

// handleAuditedHeaderUpdate creates or overwrites a header entry
func (b *SystemBackend) handleAuditedHeaderUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	header := d.Get("header").(string)
	hmac := d.Get("hmac").(bool)
	if header == "" {
		return logical.ErrorResponse("missing header name"), nil
	}

	headerConfig := b.Core.AuditedHeadersConfig()
	err := headerConfig.add(ctx, header, hmac)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// handleAuditedHeaderDelete deletes the header with the given name
func (b *SystemBackend) handleAuditedHeaderDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	header := d.Get("header").(string)
	if header == "" {
		return logical.ErrorResponse("missing header name"), nil
	}

	headerConfig := b.Core.AuditedHeadersConfig()
	err := headerConfig.remove(ctx, header)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// handleAuditedHeaderRead returns the header configuration for the given header name
func (b *SystemBackend) handleAuditedHeaderRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	header := d.Get("header").(string)
	if header == "" {
		return logical.ErrorResponse("missing header name"), nil
	}

	headerConfig := b.Core.AuditedHeadersConfig()
	settings, ok := headerConfig.Headers[strings.ToLower(header)]
	if !ok {
		return logical.ErrorResponse("Could not find header in config"), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			header: settings,
		},
	}, nil
}

// handleAuditedHeadersRead returns the whole audited headers config
func (b *SystemBackend) handleAuditedHeadersRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	headerConfig := b.Core.AuditedHeadersConfig()

	return &logical.Response{
		Data: map[string]interface{}{
			"headers": headerConfig.Headers,
		},
	}, nil
}

// handleCapabilitiesAccessor returns the ACL capabilities of the
// token associated with the given accessor for a given path.
func (b *SystemBackend) handleCapabilitiesAccessor(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	accessor := d.Get("accessor").(string)
	if accessor == "" {
		return logical.ErrorResponse("missing accessor"), nil
	}

	aEntry, err := b.Core.tokenStore.lookupByAccessor(ctx, accessor, false, false)
	if err != nil {
		return nil, err
	}

	d.Raw["token"] = aEntry.TokenID
	return b.handleCapabilities(ctx, req, d)
}

// handleCapabilities returns the ACL capabilities of the token for a given path
func (b *SystemBackend) handleCapabilities(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var token string
	if strings.HasSuffix(req.Path, "capabilities-self") {
		token = req.ClientToken
	} else {
		tokenRaw, ok := d.Raw["token"]
		if ok {
			token, _ = tokenRaw.(string)
		}
	}
	if token == "" {
		return nil, fmt.Errorf("no token found")
	}

	ret := &logical.Response{
		Data: map[string]interface{}{},
	}

	paths := d.Get("paths").([]string)
	if len(paths) == 0 {
		// Read from the deprecated field
		paths = d.Get("path").([]string)
	}

	if len(paths) == 0 {
		return logical.ErrorResponse("paths must be supplied"), nil
	}

	for _, path := range paths {
		pathCap, err := b.Core.Capabilities(ctx, token, path)
		if err != nil {
			if !strings.HasSuffix(req.Path, "capabilities-self") && errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
				return nil, &logical.StatusBadRequest{Err: "invalid token"}
			}
			return nil, err
		}
		ret.Data[path] = pathCap
	}

	// This is only here for backwards compatibility
	if len(paths) == 1 {
		ret.Data["capabilities"] = ret.Data[paths[0]]
	}

	return ret, nil
}

// handleRekeyRetrieve returns backed-up, PGP-encrypted unseal keys from a
// rekey operation
func (b *SystemBackend) handleRekeyRetrieve(
	ctx context.Context,
	req *logical.Request,
	data *framework.FieldData,
	recovery bool) (*logical.Response, error) {
	backup, err := b.Core.RekeyRetrieveBackup(ctx, recovery)
	if err != nil {
		return nil, fmt.Errorf("unable to look up backed-up keys: %w", err)
	}
	if backup == nil {
		return logical.ErrorResponse("no backed-up keys found"), nil
	}

	keysB64 := map[string][]string{}
	for k, v := range backup.Keys {
		for _, j := range v {
			currB64Keys := keysB64[k]
			if currB64Keys == nil {
				currB64Keys = []string{}
			}
			key, err := hex.DecodeString(j)
			if err != nil {
				return nil, fmt.Errorf("error decoding hex-encoded backup key: %w", err)
			}
			currB64Keys = append(currB64Keys, base64.StdEncoding.EncodeToString(key))
			keysB64[k] = currB64Keys
		}
	}

	// Format the status
	resp := &logical.Response{
		Data: map[string]interface{}{
			"nonce":       backup.Nonce,
			"keys":        backup.Keys,
			"keys_base64": keysB64,
		},
	}

	return resp, nil
}

func (b *SystemBackend) handleRekeyRetrieveBarrier(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRekeyRetrieve(ctx, req, data, false)
}

func (b *SystemBackend) handleRekeyRetrieveRecovery(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRekeyRetrieve(ctx, req, data, true)
}

// handleRekeyDelete deletes backed-up, PGP-encrypted unseal keys from a rekey
// operation
func (b *SystemBackend) handleRekeyDelete(
	ctx context.Context,
	req *logical.Request,
	data *framework.FieldData,
	recovery bool) (*logical.Response, error) {
	err := b.Core.RekeyDeleteBackup(ctx, recovery)
	if err != nil {
		return nil, fmt.Errorf("error during deletion of backed-up keys: %w", err)
	}

	return nil, nil
}

func (b *SystemBackend) handleRekeyDeleteBarrier(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRekeyDelete(ctx, req, data, false)
}

func (b *SystemBackend) handleRekeyDeleteRecovery(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRekeyDelete(ctx, req, data, true)
}

func mountInfo(entry *MountEntry) map[string]interface{} {
	info := map[string]interface{}{
		"type":                    entry.Type,
		"description":             entry.Description,
		"accessor":                entry.Accessor,
		"local":                   entry.Local,
		"seal_wrap":               entry.SealWrap,
		"external_entropy_access": entry.ExternalEntropyAccess,
		"options":                 entry.Options,
		"uuid":                    entry.UUID,
	}
	entryConfig := map[string]interface{}{
		"default_lease_ttl": int64(entry.Config.DefaultLeaseTTL.Seconds()),
		"max_lease_ttl":     int64(entry.Config.MaxLeaseTTL.Seconds()),
		"force_no_cache":    entry.Config.ForceNoCache,
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_request_keys"); ok {
		entryConfig["audit_non_hmac_request_keys"] = rawVal.([]string)
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_response_keys"); ok {
		entryConfig["audit_non_hmac_response_keys"] = rawVal.([]string)
	}
	// Even though empty value is valid for ListingVisibility, we can ignore
	// this case during mount since there's nothing to unset/hide.
	if len(entry.Config.ListingVisibility) > 0 {
		entryConfig["listing_visibility"] = entry.Config.ListingVisibility
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("passthrough_request_headers"); ok {
		entryConfig["passthrough_request_headers"] = rawVal.([]string)
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("allowed_response_headers"); ok {
		entryConfig["allowed_response_headers"] = rawVal.([]string)
	}
	if entry.Table == credentialTableType {
		entryConfig["token_type"] = entry.Config.TokenType.String()
	}

	info["config"] = entryConfig

	return info
}

// handleMountTable handles the "mounts" endpoint to provide the mount table
func (b *SystemBackend) handleMountTable(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	b.Core.mountsLock.RLock()
	defer b.Core.mountsLock.RUnlock()

	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}

	for _, entry := range b.Core.mounts.Entries {
		// Only show entries for current namespace
		if entry.Namespace().Path != ns.Path {
			continue
		}

		cont, err := b.Core.checkReplicatedFiltering(ctx, entry, "")
		if err != nil {
			return nil, err
		}
		if cont {
			continue
		}

		// Populate mount info
		info := mountInfo(entry)
		resp.Data[entry.Path] = info
	}

	return resp, nil
}

// handleMount is used to mount a new path
func (b *SystemBackend) handleMount(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	repState := b.Core.ReplicationState()

	local := data.Get("local").(bool)
	// If we are a performance secondary cluster we should forward the request
	// to the primary. We fail early here since the view in use isn't marked as
	// readonly
	if !local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return nil, logical.ErrReadOnly
	}

	// Get all the options
	path := data.Get("path").(string)
	path = sanitizePath(path)

	logicalType := data.Get("type").(string)
	description := data.Get("description").(string)
	pluginName := data.Get("plugin_name").(string)
	sealWrap := data.Get("seal_wrap").(bool)
	externalEntropyAccess := data.Get("external_entropy_access").(bool)
	options := data.Get("options").(map[string]string)

	var config MountConfig
	var apiConfig APIMountConfig

	configMap := data.Get("config").(map[string]interface{})
	if configMap != nil && len(configMap) != 0 {
		err := mapstructure.Decode(configMap, &apiConfig)
		if err != nil {
			return logical.ErrorResponse(
					"unable to convert given mount config information"),
				logical.ErrInvalidRequest
		}
	}

	switch apiConfig.DefaultLeaseTTL {
	case "":
	case "system":
	default:
		tmpDef, err := parseutil.ParseDurationSecond(apiConfig.DefaultLeaseTTL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
					"unable to parse default TTL of %s: %s", apiConfig.DefaultLeaseTTL, err)),
				logical.ErrInvalidRequest
		}
		config.DefaultLeaseTTL = tmpDef
	}

	switch apiConfig.MaxLeaseTTL {
	case "":
	case "system":
	default:
		tmpMax, err := parseutil.ParseDurationSecond(apiConfig.MaxLeaseTTL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
					"unable to parse max TTL of %s: %s", apiConfig.MaxLeaseTTL, err)),
				logical.ErrInvalidRequest
		}
		config.MaxLeaseTTL = tmpMax
	}

	if config.MaxLeaseTTL != 0 && config.DefaultLeaseTTL > config.MaxLeaseTTL {
		return logical.ErrorResponse(
				"given default lease TTL greater than given max lease TTL"),
			logical.ErrInvalidRequest
	}

	if config.DefaultLeaseTTL > b.Core.maxLeaseTTL && config.MaxLeaseTTL == 0 {
		return logical.ErrorResponse(fmt.Sprintf(
				"given default lease TTL greater than system max lease TTL of %d", int(b.Core.maxLeaseTTL.Seconds()))),
			logical.ErrInvalidRequest
	}

	switch logicalType {
	case "":
		return logical.ErrorResponse(
				"backend type must be specified as a string"),
			logical.ErrInvalidRequest
	case "plugin":
		// Only set plugin-name if mount is of type plugin, with apiConfig.PluginName
		// option taking precedence.
		switch {
		case apiConfig.PluginName != "":
			logicalType = apiConfig.PluginName
		case pluginName != "":
			logicalType = pluginName
		default:
			return logical.ErrorResponse(
					"plugin_name must be provided for plugin backend"),
				logical.ErrInvalidRequest
		}
	}

	switch logicalType {
	case "kv":
	case "kv-v1":
		// Alias KV v1
		logicalType = "kv"
		if options == nil {
			options = map[string]string{}
		}
		options["version"] = "1"

	case "kv-v2":
		// Alias KV v2
		logicalType = "kv"
		if options == nil {
			options = map[string]string{}
		}
		options["version"] = "2"

	default:
		if options != nil && options["version"] != "" {
			return logical.ErrorResponse(fmt.Sprintf(
					"secrets engine %q does not allow setting a version", logicalType)),
				logical.ErrInvalidRequest
		}
	}

	// Copy over the force no cache if set
	if apiConfig.ForceNoCache {
		config.ForceNoCache = true
	}

	if err := checkListingVisibility(apiConfig.ListingVisibility); err != nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid listing_visibility %s", apiConfig.ListingVisibility)), nil
	}
	config.ListingVisibility = apiConfig.ListingVisibility

	if len(apiConfig.AuditNonHMACRequestKeys) > 0 {
		config.AuditNonHMACRequestKeys = apiConfig.AuditNonHMACRequestKeys
	}
	if len(apiConfig.AuditNonHMACResponseKeys) > 0 {
		config.AuditNonHMACResponseKeys = apiConfig.AuditNonHMACResponseKeys
	}
	if len(apiConfig.PassthroughRequestHeaders) > 0 {
		config.PassthroughRequestHeaders = apiConfig.PassthroughRequestHeaders
	}
	if len(apiConfig.AllowedResponseHeaders) > 0 {
		config.AllowedResponseHeaders = apiConfig.AllowedResponseHeaders
	}

	// Create the mount entry
	me := &MountEntry{
		Table:                 mountTableType,
		Path:                  path,
		Type:                  logicalType,
		Description:           description,
		Config:                config,
		Local:                 local,
		SealWrap:              sealWrap,
		ExternalEntropyAccess: externalEntropyAccess,
		Options:               options,
	}

	// Attempt mount
	if err := b.Core.mount(ctx, me); err != nil {
		b.Backend.Logger().Error("error occurred during enable mount", "path", me.Path, "error", err)
		return handleError(err)
	}

	return nil, nil
}

// used to intercept an HTTPCodedError so it goes back to callee
func handleError(
	err error) (*logical.Response, error) {
	if strings.Contains(err.Error(), logical.ErrReadOnly.Error()) {
		return logical.ErrorResponse(err.Error()), err
	}
	switch err.(type) {
	case logical.HTTPCodedError:
		return logical.ErrorResponse(err.Error()), err
	default:
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
}

// Performs a similar function to handleError, but upon seeing a ReadOnlyError
// will actually strip it out to prevent forwarding
func handleErrorNoReadOnlyForward(
	err error) (*logical.Response, error) {
	if strings.Contains(err.Error(), logical.ErrReadOnly.Error()) {
		return nil, fmt.Errorf("operation could not be completed as storage is read-only")
	}
	switch err.(type) {
	case logical.HTTPCodedError:
		return logical.ErrorResponse(err.Error()), err
	default:
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
}

// handleUnmount is used to unmount a path
func (b *SystemBackend) handleUnmount(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	path = sanitizePath(path)

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	repState := b.Core.ReplicationState()
	entry := b.Core.router.MatchingMountEntry(ctx, path)

	// If we are a performance secondary cluster we should forward the request
	// to the primary. We fail early here since the view in use isn't marked as
	// readonly
	if entry != nil && !entry.Local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return nil, logical.ErrReadOnly
	}

	// We return success when the mount does not exist to not expose if the
	// mount existed or not
	match := b.Core.router.MatchingMount(ctx, path)
	if match == "" || ns.Path+path != match {
		return nil, nil
	}

	_, found := b.Core.router.MatchingStoragePrefixByAPIPath(ctx, path)
	if !found {
		b.Backend.Logger().Error("unable to find storage for path", "path", path)
		return handleError(fmt.Errorf("unable to find storage for path: %q", path))
	}

	// Attempt unmount
	if err := b.Core.unmount(ctx, path); err != nil {
		b.Backend.Logger().Error("unmount failed", "path", path, "error", err)
		return handleError(err)
	}

	// Get the view path if available
	var viewPath string
	if entry != nil {
		viewPath = entry.ViewPath()
	}

	// Remove from filtered mounts
	if err := b.Core.removePathFromFilteredPaths(ctx, ns.Path+path, viewPath); err != nil {
		b.Backend.Logger().Error("filtered path removal failed", path, "error", err)
		return handleError(err)
	}

	return nil, nil
}

func validateMountPath(p string) error {
	hasSuffix := strings.HasSuffix(p, "/")
	s := path.Clean(p)
	// Retain the trailing slash if it was provided
	if hasSuffix {
		s = s + "/"
	}
	if p != s {
		return fmt.Errorf("path '%v' does not match cleaned path '%v'", p, s)
	}

	// Check URL path for non-printable characters
	idx := strings.IndexFunc(p, func(c rune) bool {
		return !unicode.IsPrint(c)
	})

	if idx != -1 {
		return errors.New("path cannot contain non-printable characters")
	}

	return nil
}

// handleRemount is used to remount a path
func (b *SystemBackend) handleRemount(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	repState := b.Core.ReplicationState()

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get the paths
	fromPath := data.Get("from").(string)
	toPath := data.Get("to").(string)
	if fromPath == "" || toPath == "" {
		return logical.ErrorResponse(
				"both 'from' and 'to' path must be specified as a string"),
			logical.ErrInvalidRequest
	}

	if err = validateMountPath(toPath); err != nil {
		return handleError(fmt.Errorf("'to' %v", err))
	}

	entry := b.Core.router.MatchingMountEntry(ctx, fromPath)
	// If we are a performance secondary cluster we should forward the request
	// to the primary. We fail early here since the view in use isn't marked as
	// readonly
	if entry != nil && !entry.Local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return nil, logical.ErrReadOnly
	}

	// Attempt remount
	if err := b.Core.remount(ctx, fromPath, toPath, !b.Core.perfStandby); err != nil {
		b.Backend.Logger().Error("remount failed", "from_path", fromPath, "to_path", toPath, "error", err)
		return handleError(err)
	}

	// Get the view path if available
	var viewPath string
	if entry != nil {
		viewPath = entry.ViewPath()
	}

	// Remove from filtered mounts and restart evaluation process
	if err := b.Core.removePathFromFilteredPaths(ctx, ns.Path+fromPath, viewPath); err != nil {
		b.Backend.Logger().Error("filtered path removal failed", fromPath, "error", err)
		return handleError(err)
	}

	// Update quotas with the new path
	if err := b.Core.quotaManager.HandleRemount(ctx, ns.Path, sanitizePath(fromPath), sanitizePath(toPath)); err != nil {
		b.Core.logger.Error("failed to update quotas after remount", "ns_path", ns.Path, "from_path", fromPath, "to_path", toPath, "error", err)
		return handleError(err)
	}

	return nil, nil
}

// handleAuthTuneRead is used to get config settings on a auth path
func (b *SystemBackend) handleAuthTuneRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse(
				"path must be specified as a string"),
			logical.ErrInvalidRequest
	}
	return b.handleTuneReadCommon(ctx, "auth/"+path)
}

// handleMountTuneRead is used to get config settings on a backend
func (b *SystemBackend) handleMountTuneRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse(
				"path must be specified as a string"),
			logical.ErrInvalidRequest
	}

	// This call will read both logical backend's configuration as well as auth methods'.
	// Retaining this behavior for backward compatibility. If this behavior is not desired,
	// an error can be returned if path has a prefix of "auth/".
	return b.handleTuneReadCommon(ctx, path)
}

// handleTuneReadCommon returns the config settings of a path
func (b *SystemBackend) handleTuneReadCommon(ctx context.Context, path string) (*logical.Response, error) {
	path = sanitizePath(path)

	sysView := b.Core.router.MatchingSystemView(ctx, path)
	if sysView == nil {
		b.Backend.Logger().Error("cannot fetch sysview", "path", path)
		return handleError(fmt.Errorf("cannot fetch sysview for path %q", path))
	}

	mountEntry := b.Core.router.MatchingMountEntry(ctx, path)
	if mountEntry == nil {
		b.Backend.Logger().Error("cannot fetch mount entry", "path", path)
		return handleError(fmt.Errorf("cannot fetch mount entry for path %q", path))
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"description":       mountEntry.Description,
			"default_lease_ttl": int(sysView.DefaultLeaseTTL().Seconds()),
			"max_lease_ttl":     int(sysView.MaxLeaseTTL().Seconds()),
			"force_no_cache":    mountEntry.Config.ForceNoCache,
		},
	}

	// not tunable so doesn't need to be stored/loaded through synthesizedConfigCache
	if mountEntry.ExternalEntropyAccess {
		resp.Data["external_entropy_access"] = true
	}

	if mountEntry.Table == credentialTableType {
		resp.Data["token_type"] = mountEntry.Config.TokenType.String()
	}

	if rawVal, ok := mountEntry.synthesizedConfigCache.Load("audit_non_hmac_request_keys"); ok {
		resp.Data["audit_non_hmac_request_keys"] = rawVal.([]string)
	}

	if rawVal, ok := mountEntry.synthesizedConfigCache.Load("audit_non_hmac_response_keys"); ok {
		resp.Data["audit_non_hmac_response_keys"] = rawVal.([]string)
	}

	if len(mountEntry.Config.ListingVisibility) > 0 {
		resp.Data["listing_visibility"] = mountEntry.Config.ListingVisibility
	}

	if rawVal, ok := mountEntry.synthesizedConfigCache.Load("passthrough_request_headers"); ok {
		resp.Data["passthrough_request_headers"] = rawVal.([]string)
	}

	if rawVal, ok := mountEntry.synthesizedConfigCache.Load("allowed_response_headers"); ok {
		resp.Data["allowed_response_headers"] = rawVal.([]string)
	}

	if len(mountEntry.Options) > 0 {
		resp.Data["options"] = mountEntry.Options
	}

	return resp, nil
}

// handleAuthTuneWrite is used to set config settings on an auth path
func (b *SystemBackend) handleAuthTuneWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse("missing path"), nil
	}

	return b.handleTuneWriteCommon(ctx, "auth/"+path, data)
}

// handleMountTuneWrite is used to set config settings on a backend
func (b *SystemBackend) handleMountTuneWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse("missing path"), nil
	}

	// This call will write both logical backend's configuration as well as auth methods'.
	// Retaining this behavior for backward compatibility. If this behavior is not desired,
	// an error can be returned if path has a prefix of "auth/".
	return b.handleTuneWriteCommon(ctx, path, data)
}

// handleTuneWriteCommon is used to set config settings on a path
func (b *SystemBackend) handleTuneWriteCommon(ctx context.Context, path string, data *framework.FieldData) (*logical.Response, error) {
	repState := b.Core.ReplicationState()

	path = sanitizePath(path)

	// Prevent protected paths from being changed
	for _, p := range untunableMounts {
		if strings.HasPrefix(path, p) {
			b.Backend.Logger().Error("cannot tune this mount", "path", path)
			return handleError(fmt.Errorf("cannot tune %q", path))
		}
	}

	mountEntry := b.Core.router.MatchingMountEntry(ctx, path)
	if mountEntry == nil {
		b.Backend.Logger().Error("tune failed", "error", "no mount entry found", "path", path)
		return handleError(fmt.Errorf("tune of path %q failed: no mount entry found", path))
	}
	if mountEntry != nil && !mountEntry.Local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return nil, logical.ErrReadOnly
	}

	var lock *sync.RWMutex
	switch {
	case strings.HasPrefix(path, credentialRoutePrefix):
		lock = &b.Core.authLock
	default:
		lock = &b.Core.mountsLock
	}

	lock.Lock()
	defer lock.Unlock()

	// Check again after grabbing the lock
	mountEntry = b.Core.router.MatchingMountEntry(ctx, path)
	if mountEntry == nil {
		b.Backend.Logger().Error("tune failed", "error", "no mount entry found", "path", path)
		return handleError(fmt.Errorf("tune of path %q failed: no mount entry found", path))
	}
	if mountEntry != nil && !mountEntry.Local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return nil, logical.ErrReadOnly
	}

	// Timing configuration parameters
	{
		var newDefault, newMax time.Duration
		defTTL := data.Get("default_lease_ttl").(string)
		switch defTTL {
		case "":
			newDefault = mountEntry.Config.DefaultLeaseTTL
		case "system":
			newDefault = time.Duration(0)
		default:
			tmpDef, err := parseutil.ParseDurationSecond(defTTL)
			if err != nil {
				return handleError(err)
			}
			newDefault = tmpDef
		}

		maxTTL := data.Get("max_lease_ttl").(string)
		switch maxTTL {
		case "":
			newMax = mountEntry.Config.MaxLeaseTTL
		case "system":
			newMax = time.Duration(0)
		default:
			tmpMax, err := parseutil.ParseDurationSecond(maxTTL)
			if err != nil {
				return handleError(err)
			}
			newMax = tmpMax
		}

		if newDefault != mountEntry.Config.DefaultLeaseTTL ||
			newMax != mountEntry.Config.MaxLeaseTTL {

			if err := b.tuneMountTTLs(ctx, path, mountEntry, newDefault, newMax); err != nil {
				b.Backend.Logger().Error("tuning failed", "path", path, "error", err)
				return handleError(err)
			}
		}
	}

	if rawVal, ok := data.GetOk("description"); ok {
		description := rawVal.(string)

		oldDesc := mountEntry.Description
		mountEntry.Description = description

		// Update the mount table
		var err error
		switch {
		case strings.HasPrefix(path, "auth/"):
			err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
		default:
			err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
		}
		if err != nil {
			mountEntry.Description = oldDesc
			return handleError(err)
		}
		if b.Core.logger.IsInfo() {
			b.Core.logger.Info("mount tuning of description successful", "path", path, "description", description)
		}
	}

	if rawVal, ok := data.GetOk("audit_non_hmac_request_keys"); ok {
		auditNonHMACRequestKeys := rawVal.([]string)

		oldVal := mountEntry.Config.AuditNonHMACRequestKeys
		mountEntry.Config.AuditNonHMACRequestKeys = auditNonHMACRequestKeys

		// Update the mount table
		var err error
		switch {
		case strings.HasPrefix(path, "auth/"):
			err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
		default:
			err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
		}
		if err != nil {
			mountEntry.Config.AuditNonHMACRequestKeys = oldVal
			return handleError(err)
		}

		mountEntry.SyncCache()

		if b.Core.logger.IsInfo() {
			b.Core.logger.Info("mount tuning of audit_non_hmac_request_keys successful", "path", path)
		}
	}

	if rawVal, ok := data.GetOk("audit_non_hmac_response_keys"); ok {
		auditNonHMACResponseKeys := rawVal.([]string)

		oldVal := mountEntry.Config.AuditNonHMACResponseKeys
		mountEntry.Config.AuditNonHMACResponseKeys = auditNonHMACResponseKeys

		// Update the mount table
		var err error
		switch {
		case strings.HasPrefix(path, "auth/"):
			err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
		default:
			err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
		}
		if err != nil {
			mountEntry.Config.AuditNonHMACResponseKeys = oldVal
			return handleError(err)
		}

		mountEntry.SyncCache()

		if b.Core.logger.IsInfo() {
			b.Core.logger.Info("mount tuning of audit_non_hmac_response_keys successful", "path", path)
		}
	}

	if rawVal, ok := data.GetOk("listing_visibility"); ok {
		lvString := rawVal.(string)
		listingVisibility := ListingVisibilityType(lvString)

		if err := checkListingVisibility(listingVisibility); err != nil {
			return logical.ErrorResponse(fmt.Sprintf("invalid listing_visibility %s", listingVisibility)), nil
		}

		oldVal := mountEntry.Config.ListingVisibility
		mountEntry.Config.ListingVisibility = listingVisibility

		// Update the mount table
		var err error
		switch {
		case strings.HasPrefix(path, "auth/"):
			err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
		default:
			err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
		}
		if err != nil {
			mountEntry.Config.ListingVisibility = oldVal
			return handleError(err)
		}

		if b.Core.logger.IsInfo() {
			b.Core.logger.Info("mount tuning of listing_visibility successful", "path", path)
		}
	}

	if rawVal, ok := data.GetOk("token_type"); ok {
		if !strings.HasPrefix(path, "auth/") {
			return logical.ErrorResponse(fmt.Sprintf("'token_type' can only be modified on auth mounts")), logical.ErrInvalidRequest
		}
		if mountEntry.Type == "token" || mountEntry.Type == "ns_token" {
			return logical.ErrorResponse(fmt.Sprintf("'token_type' cannot be set for 'token' or 'ns_token' auth mounts")), logical.ErrInvalidRequest
		}

		tokenType := logical.TokenTypeDefaultService
		ttString := rawVal.(string)

		switch ttString {
		case "", "default-service":
		case "default-batch":
			tokenType = logical.TokenTypeDefaultBatch
		case "service":
			tokenType = logical.TokenTypeService
		case "batch":
			tokenType = logical.TokenTypeBatch
		default:
			return logical.ErrorResponse(fmt.Sprintf(
				"invalid value for 'token_type'")), logical.ErrInvalidRequest
		}

		oldVal := mountEntry.Config.TokenType
		mountEntry.Config.TokenType = tokenType

		// Update the mount table
		if err := b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local); err != nil {
			mountEntry.Config.TokenType = oldVal
			return handleError(err)
		}

		if b.Core.logger.IsInfo() {
			b.Core.logger.Info("mount tuning of token_type successful", "path", path, "token_type", ttString)
		}
	}

	if rawVal, ok := data.GetOk("passthrough_request_headers"); ok {
		headers := rawVal.([]string)

		oldVal := mountEntry.Config.PassthroughRequestHeaders
		mountEntry.Config.PassthroughRequestHeaders = headers

		// Update the mount table
		var err error
		switch {
		case strings.HasPrefix(path, "auth/"):
			err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
		default:
			err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
		}
		if err != nil {
			mountEntry.Config.PassthroughRequestHeaders = oldVal
			return handleError(err)
		}

		mountEntry.SyncCache()

		if b.Core.logger.IsInfo() {
			b.Core.logger.Info("mount tuning of passthrough_request_headers successful", "path", path)
		}
	}

	if rawVal, ok := data.GetOk("allowed_response_headers"); ok {
		headers := rawVal.([]string)

		oldVal := mountEntry.Config.AllowedResponseHeaders
		mountEntry.Config.AllowedResponseHeaders = headers

		// Update the mount table
		var err error
		switch {
		case strings.HasPrefix(path, "auth/"):
			err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
		default:
			err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
		}
		if err != nil {
			mountEntry.Config.AllowedResponseHeaders = oldVal
			return handleError(err)
		}

		mountEntry.SyncCache()

		if b.Core.logger.IsInfo() {
			b.Core.logger.Info("mount tuning of allowed_response_headers successful", "path", path)
		}
	}

	var err error
	var resp *logical.Response
	var options map[string]string
	if optionsRaw, ok := data.GetOk("options"); ok {
		options = optionsRaw.(map[string]string)
	}

	if len(options) > 0 {
		b.Core.logger.Info("mount tuning of options", "path", path, "options", options)
		newOptions := make(map[string]string)
		var kvUpgraded bool

		// The version options should only apply to the KV mount, check that first
		if v, ok := options["version"]; ok {
			// Special case to make sure we can not disable versioning once it's
			// enabled. If the vkv backend suports downgrading this can be removed.
			meVersion, err := parseutil.ParseInt(mountEntry.Options["version"])
			if err != nil {
				return nil, fmt.Errorf("unable to parse mount entry: %w", err)
			}
			optVersion, err := parseutil.ParseInt(v)
			if err != nil {
				return handleError(fmt.Errorf("unable to parse options: %w", err))
			}

			// Only accept valid versions
			switch optVersion {
			case 1:
			case 2:
			default:
				return logical.ErrorResponse(fmt.Sprintf("invalid version provided: %d", optVersion)), logical.ErrInvalidRequest
			}

			if meVersion > optVersion {
				// Return early if version option asks for a downgrade
				return logical.ErrorResponse(fmt.Sprintf("cannot downgrade mount from version %d", meVersion)), logical.ErrInvalidRequest
			}
			if meVersion < optVersion {
				kvUpgraded = true
				resp = &logical.Response{}
				resp.AddWarning(fmt.Sprintf("Upgrading mount from version %d to version %d. This mount will be unavailable for a brief period and will resume service shortly.", meVersion, optVersion))
			}
		}

		// Upsert options value to a copy of the existing mountEntry's options
		for k, v := range mountEntry.Options {
			newOptions[k] = v
		}
		for k, v := range options {
			// If the value of the provided option is empty, delete the key We
			// special-case the version value here to guard against KV downgrades, but
			// this piece could potentially be refactored in the future to be non-KV
			// specific.
			if len(v) == 0 && k != "version" {
				delete(newOptions, k)
			} else {
				newOptions[k] = v
			}
		}

		// Update the mount table
		oldVal := mountEntry.Options
		mountEntry.Options = newOptions
		switch {
		case strings.HasPrefix(path, "auth/"):
			err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
		default:
			err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
		}
		if err != nil {
			mountEntry.Options = oldVal
			return handleError(err)
		}

		// Reload the backend to kick off the upgrade process. It should only apply to KV backend so we
		// trigger based on the version logic above.
		if kvUpgraded {
			err = b.Core.reloadBackendCommon(ctx, mountEntry, strings.HasPrefix(path, credentialRoutePrefix))
			if err != nil {
				b.Core.logger.Error("mount tuning of options: could not reload backend", "error", err, "path", path, "options", options)
			}

		}
	}

	return resp, nil
}

// handleLease is use to view the metadata for a given LeaseID
func (b *SystemBackend) handleLeaseLookup(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	leaseID := data.Get("lease_id").(string)
	if leaseID == "" {
		return logical.ErrorResponse("lease_id must be specified"),
			logical.ErrInvalidRequest
	}

	leaseTimes, err := b.Core.expiration.FetchLeaseTimes(ctx, leaseID)
	if err != nil {
		b.Backend.Logger().Error("error retrieving lease", "lease_id", leaseID, "error", err)
		return handleError(err)
	}
	if leaseTimes == nil {
		return logical.ErrorResponse("invalid lease"), logical.ErrInvalidRequest
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"id":           leaseID,
			"issue_time":   leaseTimes.IssueTime,
			"expire_time":  nil,
			"last_renewal": nil,
			"ttl":          int64(0),
		},
	}
	renewable, _ := leaseTimes.renewable()
	resp.Data["renewable"] = renewable

	if !leaseTimes.LastRenewalTime.IsZero() {
		resp.Data["last_renewal"] = leaseTimes.LastRenewalTime
	}
	if !leaseTimes.ExpireTime.IsZero() {
		resp.Data["expire_time"] = leaseTimes.ExpireTime
		resp.Data["ttl"] = leaseTimes.ttl()
	}
	return resp, nil
}

func (b *SystemBackend) handleLeaseLookupList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	prefix := data.Get("prefix").(string)
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	view := b.Core.expiration.leaseView(ns)
	keys, err := view.List(ctx, prefix)
	if err != nil {
		b.Backend.Logger().Error("error listing leases", "prefix", prefix, "error", err)
		return handleErrorNoReadOnlyForward(err)
	}
	return logical.ListResponse(keys), nil
}

// handleRenew is used to renew a lease with a given LeaseID
func (b *SystemBackend) handleRenew(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	leaseID := data.Get("lease_id").(string)
	if leaseID == "" {
		leaseID = data.Get("url_lease_id").(string)
	}
	if leaseID == "" {
		return logical.ErrorResponse("lease_id must be specified"),
			logical.ErrInvalidRequest
	}
	incrementRaw := data.Get("increment").(int)

	// Convert the increment
	increment := time.Duration(incrementRaw) * time.Second

	// Invoke the expiration manager directly
	resp, err := b.Core.expiration.Renew(ctx, leaseID, increment)
	if err != nil {
		b.Backend.Logger().Error("lease renewal failed", "lease_id", leaseID, "error", err)
		return handleErrorNoReadOnlyForward(err)
	}
	return resp, err
}

// handleRevoke is used to revoke a given LeaseID
func (b *SystemBackend) handleRevoke(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	leaseID := data.Get("lease_id").(string)
	if leaseID == "" {
		leaseID = data.Get("url_lease_id").(string)
	}
	if leaseID == "" {
		return logical.ErrorResponse("lease_id must be specified"),
			logical.ErrInvalidRequest
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	revokeCtx := namespace.ContextWithNamespace(b.Core.activeContext, ns)
	if data.Get("sync").(bool) {
		// Invoke the expiration manager directly
		if err := b.Core.expiration.Revoke(revokeCtx, leaseID); err != nil {
			b.Backend.Logger().Error("lease revocation failed", "lease_id", leaseID, "error", err)
			return handleErrorNoReadOnlyForward(err)
		}

		return nil, nil
	}

	if err := b.Core.expiration.LazyRevoke(revokeCtx, leaseID); err != nil {
		b.Backend.Logger().Error("lease revocation failed", "lease_id", leaseID, "error", err)
		return handleErrorNoReadOnlyForward(err)
	}

	return logical.RespondWithStatusCode(nil, nil, http.StatusAccepted)
}

// handleRevokePrefix is used to revoke a prefix with many LeaseIDs
func (b *SystemBackend) handleRevokePrefix(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRevokePrefixCommon(ctx, req, data, false, data.Get("sync").(bool))
}

// handleRevokeForce is used to revoke a prefix with many LeaseIDs, ignoring errors
func (b *SystemBackend) handleRevokeForce(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRevokePrefixCommon(ctx, req, data, true, true)
}

// handleRevokePrefixCommon is used to revoke a prefix with many LeaseIDs
func (b *SystemBackend) handleRevokePrefixCommon(ctx context.Context,
	req *logical.Request, data *framework.FieldData, force, sync bool) (*logical.Response, error) {
	// Get all the options
	prefix := data.Get("prefix").(string)

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Invoke the expiration manager directly
	revokeCtx := namespace.ContextWithNamespace(b.Core.activeContext, ns)
	if force {
		err = b.Core.expiration.RevokeForce(revokeCtx, prefix)
	} else {
		err = b.Core.expiration.RevokePrefix(revokeCtx, prefix, sync)
	}
	if err != nil {
		b.Backend.Logger().Error("revoke prefix failed", "prefix", prefix, "error", err)
		return handleErrorNoReadOnlyForward(err)
	}

	if sync {
		return nil, nil
	}

	return logical.RespondWithStatusCode(nil, nil, http.StatusAccepted)
}

// handleAuthTable handles the "auth" endpoint to provide the auth table
func (b *SystemBackend) handleAuthTable(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	b.Core.authLock.RLock()
	defer b.Core.authLock.RUnlock()

	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}

	for _, entry := range b.Core.auth.Entries {
		// Only show entries for current namespace
		if entry.Namespace().Path != ns.Path {
			continue
		}

		cont, err := b.Core.checkReplicatedFiltering(ctx, entry, credentialRoutePrefix)
		if err != nil {
			return nil, err
		}
		if cont {
			continue
		}

		info := mountInfo(entry)
		resp.Data[entry.Path] = info
	}

	return resp, nil
}

// handleEnableAuth is used to enable a new credential backend
func (b *SystemBackend) handleEnableAuth(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	repState := b.Core.ReplicationState()
	local := data.Get("local").(bool)

	// If we are a performance secondary cluster we should forward the request
	// to the primary. We fail early here since the view in use isn't marked as
	// readonly
	if !local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return nil, logical.ErrReadOnly
	}

	// Get all the options
	path := data.Get("path").(string)
	path = sanitizePath(path)
	logicalType := data.Get("type").(string)
	description := data.Get("description").(string)
	pluginName := data.Get("plugin_name").(string)
	sealWrap := data.Get("seal_wrap").(bool)
	externalEntropyAccess := data.Get("external_entropy_access").(bool)
	options := data.Get("options").(map[string]string)

	var config MountConfig
	var apiConfig APIMountConfig

	configMap := data.Get("config").(map[string]interface{})
	if configMap != nil && len(configMap) != 0 {
		err := mapstructure.Decode(configMap, &apiConfig)
		if err != nil {
			return logical.ErrorResponse(
					"unable to convert given auth config information"),
				logical.ErrInvalidRequest
		}
	}

	switch apiConfig.DefaultLeaseTTL {
	case "":
	case "system":
	default:
		tmpDef, err := parseutil.ParseDurationSecond(apiConfig.DefaultLeaseTTL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
					"unable to parse default TTL of %s: %s", apiConfig.DefaultLeaseTTL, err)),
				logical.ErrInvalidRequest
		}
		config.DefaultLeaseTTL = tmpDef
	}

	switch apiConfig.MaxLeaseTTL {
	case "":
	case "system":
	default:
		tmpMax, err := parseutil.ParseDurationSecond(apiConfig.MaxLeaseTTL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
					"unable to parse max TTL of %s: %s", apiConfig.MaxLeaseTTL, err)),
				logical.ErrInvalidRequest
		}
		config.MaxLeaseTTL = tmpMax
	}

	if config.MaxLeaseTTL != 0 && config.DefaultLeaseTTL > config.MaxLeaseTTL {
		return logical.ErrorResponse(
				"given default lease TTL greater than given max lease TTL"),
			logical.ErrInvalidRequest
	}

	if config.DefaultLeaseTTL > b.Core.maxLeaseTTL && config.MaxLeaseTTL == 0 {
		return logical.ErrorResponse(fmt.Sprintf(
				"given default lease TTL greater than system max lease TTL of %d", int(b.Core.maxLeaseTTL.Seconds()))),
			logical.ErrInvalidRequest
	}

	switch apiConfig.TokenType {
	case "", "default-service":
		config.TokenType = logical.TokenTypeDefaultService
	case "default-batch":
		config.TokenType = logical.TokenTypeDefaultBatch
	case "service":
		config.TokenType = logical.TokenTypeService
	case "batch":
		config.TokenType = logical.TokenTypeBatch
	default:
		return logical.ErrorResponse(fmt.Sprintf(
			"invalid value for 'token_type'")), logical.ErrInvalidRequest
	}

	switch logicalType {
	case "":
		return logical.ErrorResponse(
				"backend type must be specified as a string"),
			logical.ErrInvalidRequest
	case "plugin":
		// Only set plugin name if mount is of type plugin, with apiConfig.PluginName
		// option taking precedence.
		switch {
		case apiConfig.PluginName != "":
			logicalType = apiConfig.PluginName
		case pluginName != "":
			logicalType = pluginName
		default:
			return logical.ErrorResponse(
					"plugin_name must be provided for plugin backend"),
				logical.ErrInvalidRequest
		}
	}

	if options != nil && options["version"] != "" {
		return logical.ErrorResponse(fmt.Sprintf(
				"auth method %q does not allow setting a version", logicalType)),
			logical.ErrInvalidRequest
	}

	if err := checkListingVisibility(apiConfig.ListingVisibility); err != nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid listing_visibility %s", apiConfig.ListingVisibility)), nil
	}
	config.ListingVisibility = apiConfig.ListingVisibility

	if len(apiConfig.AuditNonHMACRequestKeys) > 0 {
		config.AuditNonHMACRequestKeys = apiConfig.AuditNonHMACRequestKeys
	}
	if len(apiConfig.AuditNonHMACResponseKeys) > 0 {
		config.AuditNonHMACResponseKeys = apiConfig.AuditNonHMACResponseKeys
	}
	if len(apiConfig.PassthroughRequestHeaders) > 0 {
		config.PassthroughRequestHeaders = apiConfig.PassthroughRequestHeaders
	}
	if len(apiConfig.AllowedResponseHeaders) > 0 {
		config.AllowedResponseHeaders = apiConfig.AllowedResponseHeaders
	}

	// Create the mount entry
	me := &MountEntry{
		Table:                 credentialTableType,
		Path:                  path,
		Type:                  logicalType,
		Description:           description,
		Config:                config,
		Local:                 local,
		SealWrap:              sealWrap,
		ExternalEntropyAccess: externalEntropyAccess,
		Options:               options,
	}

	// Attempt enabling
	if err := b.Core.enableCredential(ctx, me); err != nil {
		b.Backend.Logger().Error("error occurred during enable credential", "path", me.Path, "error", err)
		return handleError(err)
	}
	return nil, nil
}

// handleDisableAuth is used to disable a credential backend
func (b *SystemBackend) handleDisableAuth(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	path = sanitizePath(path)

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	fullPath := credentialRoutePrefix + path

	repState := b.Core.ReplicationState()
	entry := b.Core.router.MatchingMountEntry(ctx, fullPath)

	// If we are a performance secondary cluster we should forward the request
	// to the primary. We fail early here since the view in use isn't marked as
	// readonly
	if entry != nil && !entry.Local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return nil, logical.ErrReadOnly
	}

	// We return success when the mount does not exist to not expose if the
	// mount existed or not
	match := b.Core.router.MatchingMount(ctx, fullPath)
	if match == "" || ns.Path+fullPath != match {
		return nil, nil
	}

	_, found := b.Core.router.MatchingStoragePrefixByAPIPath(ctx, fullPath)
	if !found {
		b.Backend.Logger().Error("unable to find storage for path", "path", fullPath)
		return handleError(fmt.Errorf("unable to find storage for path: %q", fullPath))
	}

	// Attempt disable
	if err := b.Core.disableCredential(ctx, path); err != nil {
		b.Backend.Logger().Error("disable auth mount failed", "path", path, "error", err)
		return handleError(err)
	}

	// Get the view path if available
	var viewPath string
	if entry != nil {
		viewPath = entry.ViewPath()
	}

	// Remove from filtered mounts
	if err := b.Core.removePathFromFilteredPaths(ctx, fullPath, viewPath); err != nil {
		b.Backend.Logger().Error("filtered path removal failed", path, "error", err)
		return handleError(err)
	}

	return nil, nil
}

// handlePoliciesList handles /sys/policy/ and /sys/policies/<type> endpoints to provide the enabled policies
func (b *SystemBackend) handlePoliciesList(policyType PolicyType) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		policies, err := b.Core.policyStore.ListPolicies(ctx, policyType)
		if err != nil {
			return nil, err
		}

		switch policyType {
		case PolicyTypeACL:
			// Add the special "root" policy if not egp and we are at the root namespace
			if ns.ID == namespace.RootNamespaceID {
				policies = append(policies, "root")
			}
			resp := logical.ListResponse(policies)

			// If the request is from sys/policy/ we handle backwards compatibility
			if strings.HasPrefix(req.Path, "policy") {
				resp.Data["policies"] = resp.Data["keys"]
			}
			return resp, nil

		case PolicyTypeRGP:
			return logical.ListResponse(policies), nil

		case PolicyTypeEGP:
			nsScopedKeyInfo := getEGPListResponseKeyInfo(b, ns)
			return &logical.Response{
				Data: map[string]interface{}{
					"keys":     policies,
					"key_info": nsScopedKeyInfo,
				},
			}, nil
		}

		return logical.ErrorResponse("unknown policy type"), nil
	}
}

// handlePoliciesRead handles the "/sys/policy/<name>" and "/sys/policies/<type>/<name>" endpoints to read a policy
func (b *SystemBackend) handlePoliciesRead(policyType PolicyType) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)

		policy, err := b.Core.policyStore.GetPolicy(ctx, name, policyType)
		if err != nil {
			return handleError(err)
		}

		if policy == nil {
			return nil, nil
		}

		// If the request is from sys/policy/ we handle backwards compatibility
		var respDataPolicyName string
		if policyType == PolicyTypeACL && strings.HasPrefix(req.Path, "policy") {
			respDataPolicyName = "rules"
		} else {
			respDataPolicyName = "policy"
		}

		resp := &logical.Response{
			Data: map[string]interface{}{
				"name":             policy.Name,
				respDataPolicyName: policy.Raw,
			},
		}

		switch policy.Type {
		case PolicyTypeRGP, PolicyTypeEGP:
			addSentinelPolicyData(resp.Data, policy)
		}

		return resp, nil
	}
}

// handlePoliciesSet handles the "/sys/policy/<name>" and "/sys/policies/<type>/<name>" endpoints to set a policy
func (b *SystemBackend) handlePoliciesSet(policyType PolicyType) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		var resp *logical.Response

		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		policy := &Policy{
			Name:      strings.ToLower(data.Get("name").(string)),
			Type:      policyType,
			namespace: ns,
		}
		if policy.Name == "" {
			return logical.ErrorResponse("policy name must be provided in the URL"), nil
		}

		policy.Raw = data.Get("policy").(string)
		if policy.Raw == "" && policyType == PolicyTypeACL && strings.HasPrefix(req.Path, "policy") {
			policy.Raw = data.Get("rules").(string)
			if resp == nil {
				resp = &logical.Response{}
			}
			resp.AddWarning("'rules' is deprecated, please use 'policy' instead")
		}
		if policy.Raw == "" {
			return logical.ErrorResponse("'policy' parameter not supplied or empty"), nil
		}

		if polBytes, err := base64.StdEncoding.DecodeString(policy.Raw); err == nil {
			policy.Raw = string(polBytes)
		}

		switch policyType {
		case PolicyTypeACL:
			p, err := ParseACLPolicy(ns, policy.Raw)
			if err != nil {
				return handleError(err)
			}
			policy.Paths = p.Paths
			policy.Templated = p.Templated

		case PolicyTypeRGP, PolicyTypeEGP:

		default:
			return logical.ErrorResponse("unknown policy type"), nil
		}

		if policy.Type == PolicyTypeRGP || policy.Type == PolicyTypeEGP {
			if errResp := inputSentinelPolicyData(data, policy); errResp != nil {
				return errResp, nil
			}
		}

		// Update the policy
		if err := b.Core.policyStore.SetPolicy(ctx, policy); err != nil {
			return handleError(err)
		}
		return resp, nil
	}
}

func (b *SystemBackend) handlePoliciesDelete(policyType PolicyType) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)

		if err := b.Core.policyStore.DeletePolicy(ctx, name, policyType); err != nil {
			return handleError(err)
		}
		return nil, nil
	}
}

type passwordPolicyConfig struct {
	HCLPolicy string `json:"policy"`
}

func getPasswordPolicyKey(policyName string) string {
	return fmt.Sprintf("password_policy/%s", policyName)
}

const (
	minPasswordLength = 4
	maxPasswordLength = 100
)

// handlePoliciesPasswordSet saves/updates password policies
func (*SystemBackend) handlePoliciesPasswordSet(ctx context.Context, req *logical.Request, data *framework.FieldData) (resp *logical.Response, err error) {
	policyName := data.Get("name").(string)
	if policyName == "" {
		return nil, logical.CodedError(http.StatusBadRequest, "missing policy name")
	}

	rawPolicy := data.Get("policy").(string)
	if rawPolicy == "" {
		return nil, logical.CodedError(http.StatusBadRequest, "missing policy")
	}

	// Optionally decode base64 string
	decodedPolicy, err := base64.StdEncoding.DecodeString(rawPolicy)
	if err == nil {
		rawPolicy = string(decodedPolicy)
	}

	// Parse the policy to ensure that it's valid
	policy, err := random.ParsePolicy(rawPolicy)
	if err != nil {
		return nil, logical.CodedError(http.StatusBadRequest, fmt.Sprintf("invalid password policy: %s", err))
	}

	if policy.Length > maxPasswordLength || policy.Length < minPasswordLength {
		return nil, logical.CodedError(http.StatusBadRequest,
			fmt.Sprintf("passwords must be between %d and %d characters", minPasswordLength, maxPasswordLength))
	}

	// Generate some passwords to ensure that we're confident that the policy isn't impossible
	timeout := 1 * time.Second
	genCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	attempts := 10
	failed := 0
	for i := 0; i < attempts; i++ {
		_, err = policy.Generate(genCtx, nil)
		if err != nil {
			failed++
		}
	}

	if failed == attempts {
		return nil, logical.CodedError(http.StatusBadRequest,
			fmt.Sprintf("unable to generate password from provided policy in %s: are the rules impossible?", timeout))
	}

	if failed > 0 {
		return nil, logical.CodedError(http.StatusBadRequest,
			fmt.Sprintf("failed to generate passwords %d times out of %d attempts in %s - is the policy too restrictive?", failed, attempts, timeout))
	}

	cfg := passwordPolicyConfig{
		HCLPolicy: rawPolicy,
	}
	entry, err := logical.StorageEntryJSON(getPasswordPolicyKey(policyName), cfg)
	if err != nil {
		return nil, logical.CodedError(http.StatusInternalServerError, fmt.Sprintf("unable to save password policy: %s", err))
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, logical.CodedError(http.StatusInternalServerError,
			fmt.Sprintf("failed to save policy to storage backend: %s", err))
	}

	return logical.RespondWithStatusCode(nil, req, http.StatusNoContent)
}

// handlePoliciesPasswordGet retrieves a password policy if it exists
func (*SystemBackend) handlePoliciesPasswordGet(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	policyName := data.Get("name").(string)
	if policyName == "" {
		return nil, logical.CodedError(http.StatusBadRequest, "missing policy name")
	}

	cfg, err := retrievePasswordPolicy(ctx, req.Storage, policyName)
	if err != nil {
		return nil, logical.CodedError(http.StatusInternalServerError, "failed to retrieve password policy")
	}
	if cfg == nil {
		return nil, logical.CodedError(http.StatusNotFound, "policy does not exist")
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"policy": cfg.HCLPolicy,
		},
	}

	return resp, nil
}

// retrievePasswordPolicy retrieves a password policy from the logical storage
func retrievePasswordPolicy(ctx context.Context, storage logical.Storage, policyName string) (policyCfg *passwordPolicyConfig, err error) {
	entry, err := storage.Get(ctx, getPasswordPolicyKey(policyName))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	policyCfg = &passwordPolicyConfig{}
	err = json.Unmarshal(entry.Value, &policyCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal stored data: %w", err)
	}

	return policyCfg, nil
}

// handlePoliciesPasswordDelete deletes a password policy if it exists
func (*SystemBackend) handlePoliciesPasswordDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	policyName := data.Get("name").(string)
	if policyName == "" {
		return nil, logical.CodedError(http.StatusBadRequest, "missing policy name")
	}

	err := req.Storage.Delete(ctx, getPasswordPolicyKey(policyName))
	if err != nil {
		return nil, logical.CodedError(http.StatusInternalServerError,
			fmt.Sprintf("failed to delete password policy: %s", err))
	}

	return nil, nil
}

// handlePoliciesPasswordGenerate generates a password from the specified password policy
func (*SystemBackend) handlePoliciesPasswordGenerate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	policyName := data.Get("name").(string)
	if policyName == "" {
		return nil, logical.CodedError(http.StatusBadRequest, "missing policy name")
	}

	cfg, err := retrievePasswordPolicy(ctx, req.Storage, policyName)
	if err != nil {
		return nil, logical.CodedError(http.StatusInternalServerError, "failed to retrieve password policy")
	}
	if cfg == nil {
		return nil, logical.CodedError(http.StatusNotFound, "policy does not exist")
	}

	policy, err := random.ParsePolicy(cfg.HCLPolicy)
	if err != nil {
		return nil, logical.CodedError(http.StatusInternalServerError,
			"stored password policy configuration failed to parse")
	}

	password, err := policy.Generate(ctx, nil)
	if err != nil {
		return nil, logical.CodedError(http.StatusInternalServerError,
			fmt.Sprintf("failed to generate password from policy: %s", err))
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"password": password,
		},
	}
	return resp, nil
}

// handleAuditTable handles the "audit" endpoint to provide the audit table
func (b *SystemBackend) handleAuditTable(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Core.auditLock.RLock()
	defer b.Core.auditLock.RUnlock()

	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}
	for _, entry := range b.Core.audit.Entries {
		info := map[string]interface{}{
			"path":        entry.Path,
			"type":        entry.Type,
			"description": entry.Description,
			"options":     entry.Options,
			"local":       entry.Local,
		}
		resp.Data[entry.Path] = info
	}
	return resp, nil
}

// handleAuditHash is used to fetch the hash of the given input data with the
// specified audit backend's salt
func (b *SystemBackend) handleAuditHash(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	input := data.Get("input").(string)
	if input == "" {
		return logical.ErrorResponse("the \"input\" parameter is empty"), nil
	}

	path = sanitizePath(path)

	hash, err := b.Core.auditBroker.GetHash(ctx, path, input)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"hash": hash,
		},
	}, nil
}

// handleEnableAudit is used to enable a new audit backend
func (b *SystemBackend) handleEnableAudit(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	repState := b.Core.ReplicationState()

	local := data.Get("local").(bool)
	// If we are a performance secondary cluster we should forward the request
	// to the primary. We fail early here since the view in use isn't marked as
	// readonly
	if !local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return nil, logical.ErrReadOnly
	}

	// Get all the options
	path := data.Get("path").(string)
	backendType := data.Get("type").(string)
	description := data.Get("description").(string)
	options := data.Get("options").(map[string]string)

	// Create the mount entry
	me := &MountEntry{
		Table:       auditTableType,
		Path:        path,
		Type:        backendType,
		Description: description,
		Options:     options,
		Local:       local,
	}

	// Attempt enabling
	if err := b.Core.enableAudit(ctx, me, true); err != nil {
		b.Backend.Logger().Error("enable audit mount failed", "path", me.Path, "error", err)
		return handleError(err)
	}
	return nil, nil
}

// handleDisableAudit is used to disable an audit backend
func (b *SystemBackend) handleDisableAudit(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	if path == "/" {
		return handleError(errors.New("audit device path must be specified"))
	}

	b.Core.auditLock.RLock()
	table := b.Core.audit.shallowClone()
	entry, err := table.find(ctx, path)
	b.Core.auditLock.RUnlock()

	if err != nil {
		return handleError(err)
	}
	if entry == nil {
		return nil, nil
	}

	repState := b.Core.ReplicationState()

	// If we are a performance secondary cluster we should forward the request
	// to the primary. We fail early here since the view in use isn't marked as
	// readonly
	if !entry.Local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return nil, logical.ErrReadOnly
	}

	// Attempt disable
	if existed, err := b.Core.disableAudit(ctx, path, true); existed && err != nil {
		b.Backend.Logger().Error("disable audit mount failed", "path", path, "error", err)
		return handleError(err)
	}
	return nil, nil
}

func (b *SystemBackend) handleConfigUIHeadersRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	header := data.Get("header").(string)
	multivalue := data.Get("multivalue").(bool)

	values, err := b.Core.uiConfig.GetHeader(ctx, header)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}

	// Return multiple values if specified
	if multivalue {
		return &logical.Response{
			Data: map[string]interface{}{
				"values": values,
			},
		}, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"value": values[0],
		},
	}, nil
}

func (b *SystemBackend) handleConfigUIHeadersList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	headers, err := b.Core.uiConfig.HeaderKeys(ctx)
	if err != nil {
		return nil, err
	}
	if len(headers) == 0 {
		return nil, nil
	}

	return logical.ListResponse(headers), nil
}

func (b *SystemBackend) handleConfigUIHeadersUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	header := data.Get("header").(string)
	values := data.Get("values").([]string)
	if header == "" || len(values) == 0 {
		return logical.ErrorResponse("header and values must be specified"), logical.ErrInvalidRequest
	}

	lowerHeader := strings.ToLower(header)
	if strings.HasPrefix(lowerHeader, "x-vault-") {
		return logical.ErrorResponse("X-Vault headers cannot be set"), logical.ErrInvalidRequest
	}

	// Translate the list of values to the valid header string
	value := http.Header{}
	for _, v := range values {
		value.Add(header, v)
	}
	err := b.Core.uiConfig.SetHeader(ctx, header, value.Values(header))
	if err != nil {
		return nil, err
	}

	// Warn when overriding the CSP
	resp := &logical.Response{}
	if lowerHeader == "content-security-policy" {
		resp.AddWarning("overriding default Content-Security-Policy which is secure by default, proceed with caution")
	}

	return resp, nil
}

func (b *SystemBackend) handleConfigUIHeadersDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	header := data.Get("header").(string)
	err := b.Core.uiConfig.DeleteHeader(ctx, header)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// handleKeyStatus returns status information about the backend key
func (b *SystemBackend) handleKeyStatus(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get the key info
	info, err := b.Core.barrier.ActiveKeyInfo()
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"term":         info.Term,
			"install_time": info.InstallTime.Format(time.RFC3339Nano),
			"encryptions":  info.Encryptions,
		},
	}
	return resp, nil
}

// handleKeyRotationConfigRead returns the barrier key rotation config
func (b *SystemBackend) handleKeyRotationConfigRead(_ context.Context, _ *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	// Get the key info
	rotConfig, err := b.Core.barrier.RotationConfig()
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"max_operations": rotConfig.MaxOperations,
			"enabled":        !rotConfig.Disabled,
		},
	}
	if rotConfig.Interval > 0 {
		resp.Data["interval"] = rotConfig.Interval.String()
	} else {
		resp.Data["interval"] = 0
	}
	return resp, nil
}

// handleKeyRotationConfigRead returns the barrier key rotation config
func (b *SystemBackend) handleKeyRotationConfigUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	rotConfig, err := b.Core.barrier.RotationConfig()
	if err != nil {
		return nil, err
	}
	maxOps, ok, err := data.GetOkErr("max_operations")
	if err != nil {
		return nil, err
	}
	if ok {
		rotConfig.MaxOperations = maxOps.(int64)
	}
	interval, ok, err := data.GetOkErr("interval")
	if err != nil {
		return nil, err
	}
	if ok {
		rotConfig.Interval = time.Second * time.Duration(interval.(int))
	}

	enabled, ok, err := data.GetOkErr("enabled")
	if err != nil {
		return nil, err
	}
	if ok {
		rotConfig.Disabled = !enabled.(bool)
	}

	// Reject out of range settings
	if rotConfig.Interval < minimumRotationInterval && rotConfig.Interval != 0 {
		return logical.ErrorResponse("interval must be greater or equal to %s", minimumRotationInterval.String()), logical.ErrInvalidRequest
	}

	if rotConfig.MaxOperations < absoluteOperationMinimum || rotConfig.MaxOperations > absoluteOperationMaximum {
		return logical.ErrorResponse("max_operations must be in the range [%d,%d]", absoluteOperationMinimum, absoluteOperationMaximum), logical.ErrInvalidRequest
	}

	// Store the rotation config
	b.Core.barrier.SetRotationConfig(ctx, rotConfig)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// handleRotate is used to trigger a key rotation
func (b *SystemBackend) handleRotate(ctx context.Context, _ *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	repState := b.Core.ReplicationState()
	if repState.HasState(consts.ReplicationPerformanceSecondary) {
		return logical.ErrorResponse("cannot rotate on a replication secondary"), nil
	}

	if err := b.rotateBarrierKey(ctx); err != nil {
		b.Backend.Logger().Error("error handling key rotation", "error", err)
		return handleError(err)
	}
	return nil, nil
}

func (b *SystemBackend) handleWrappingPubkey(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	x, _ := b.Core.wrappingJWTKey.X.MarshalText()
	y, _ := b.Core.wrappingJWTKey.Y.MarshalText()
	return &logical.Response{
		Data: map[string]interface{}{
			"jwt_x":     string(x),
			"jwt_y":     string(y),
			"jwt_curve": corePrivateKeyTypeP521,
		},
	}, nil
}

func (b *SystemBackend) handleWrappingWrap(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.WrapInfo == nil || req.WrapInfo.TTL == 0 {
		return logical.ErrorResponse("endpoint requires response wrapping to be used"), logical.ErrInvalidRequest
	}

	// N.B.: Do *NOT* allow JWT wrapping tokens to be created through this
	// endpoint. JWTs are signed so if we don't allow users to create wrapping
	// tokens using them we can ensure that an operator can't spoof a legit JWT
	// wrapped token, which makes certain init/rekey/generate-root cases have
	// better properties.
	req.WrapInfo.Format = "uuid"

	return &logical.Response{
		Data: data.Raw,
	}, nil
}

// handleWrappingUnwrap will unwrap a response wrapping token or complete a
// request that required a control group.
func (b *SystemBackend) handleWrappingUnwrap(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// If a third party is unwrapping (rather than the calling token being the
	// wrapping token) we detect this so that we can revoke the original
	// wrapping token after reading it
	var thirdParty bool

	token := data.Get("token").(string)
	if token != "" {
		thirdParty = true
	} else {
		token = req.ClientToken
	}

	// Get the policies so we can determine if this is a normal response
	// wrapping request or a control group token.
	//
	// We use lookupTainted here because the token might have already been used
	// by handleRequest(), this happens when it's a normal response wrapping
	// request and the token was provided "first party". We want to inspect the
	// token policies but will not use this token entry for anything else.
	te, err := b.Core.tokenStore.lookupTainted(ctx, token)
	if err != nil {
		return nil, err
	}
	if te == nil {
		return nil, nil
	}
	if len(te.Policies) != 1 {
		return nil, errors.New("token is not a valid unwrap token")
	}

	unwrapNS, err := NamespaceByID(ctx, te.NamespaceID, b.Core)
	if err != nil {
		return nil, err
	}
	unwrapCtx := namespace.ContextWithNamespace(ctx, unwrapNS)

	var response string
	switch te.Policies[0] {
	case controlGroupPolicyName:
		response, err = controlGroupUnwrap(unwrapCtx, b, token, thirdParty)
	case responseWrappingPolicyName:
		response, err = b.responseWrappingUnwrap(unwrapCtx, te, thirdParty)
	}
	if err != nil {
		var respErr *logical.Response
		if len(response) > 0 {
			respErr = logical.ErrorResponse(response)
		}

		return respErr, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{},
	}

	if len(response) == 0 {
		resp.Data[logical.HTTPStatusCode] = 204
		return resp, nil
	}

	// Most of the time we want to just send over the marshalled HTTP bytes.
	// However there is a sad separate case: if the original response was using
	// bare values we need to use those or else what comes back is garbled.
	httpResp := &logical.HTTPResponse{}
	err = jsonutil.DecodeJSON([]byte(response), httpResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding wrapped response: %w", err)
	}
	if httpResp.Data != nil &&
		(httpResp.Data[logical.HTTPStatusCode] != nil ||
			httpResp.Data[logical.HTTPRawBody] != nil ||
			httpResp.Data[logical.HTTPContentType] != nil) {
		if httpResp.Data[logical.HTTPStatusCode] != nil {
			resp.Data[logical.HTTPStatusCode] = httpResp.Data[logical.HTTPStatusCode]
		}
		if httpResp.Data[logical.HTTPContentType] != nil {
			resp.Data[logical.HTTPContentType] = httpResp.Data[logical.HTTPContentType]
		}

		rawBody := httpResp.Data[logical.HTTPRawBody]
		if rawBody != nil {
			// Decode here so that we can audit properly
			switch rawBody.(type) {
			case string:
				// Best effort decoding; if this works, the original value was
				// probably a []byte instead of a string, but was marshaled
				// when the value was saved, so this restores it as it was
				decBytes, err := base64.StdEncoding.DecodeString(rawBody.(string))
				if err == nil {
					// We end up with []byte, will not be HMAC'd
					resp.Data[logical.HTTPRawBody] = decBytes
				} else {
					// We end up with string, will be HMAC'd
					resp.Data[logical.HTTPRawBody] = rawBody
				}
			default:
				b.Core.Logger().Error("unexpected type of raw body when decoding wrapped token", "type", fmt.Sprintf("%T", rawBody))
			}

			resp.Data[logical.HTTPRawBodyAlreadyJSONDecoded] = true
		}

		return resp, nil
	}

	resp.Data[logical.HTTPStatusCode] = 200
	resp.Data[logical.HTTPRawBody] = []byte(response)
	resp.Data[logical.HTTPContentType] = "application/json"

	return resp, nil
}

// responseWrappingUnwrap will read the stored response in the cubbyhole and
// return the raw HTTP response.
func (b *SystemBackend) responseWrappingUnwrap(ctx context.Context, te *logical.TokenEntry, thirdParty bool) (string, error) {
	tokenID := te.ID
	if thirdParty {
		// Use the token to decrement the use count to avoid a second operation on the token.
		_, err := b.Core.tokenStore.UseTokenByID(ctx, tokenID)
		if err != nil {
			return "", fmt.Errorf("error decrementing wrapping token's use-count: %w", err)
		}

		defer b.Core.tokenStore.revokeOrphan(ctx, tokenID)
	}

	cubbyReq := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "cubbyhole/response",
		ClientToken: tokenID,
	}
	cubbyReq.SetTokenEntry(te)
	cubbyResp, err := b.Core.router.Route(ctx, cubbyReq)
	if err != nil {
		return "", fmt.Errorf("error looking up wrapping information: %w", err)
	}
	if cubbyResp == nil {
		return "no information found; wrapping token may be from a previous Vault version", ErrInternalError
	}
	if cubbyResp != nil && cubbyResp.IsError() {
		return cubbyResp.Error().Error(), nil
	}
	if cubbyResp.Data == nil {
		return "wrapping information was nil; wrapping token may be from a previous Vault version", ErrInternalError
	}

	responseRaw := cubbyResp.Data["response"]
	if responseRaw == nil {
		return "", fmt.Errorf("no response found inside the cubbyhole")
	}
	response, ok := responseRaw.(string)
	if !ok {
		return "", fmt.Errorf("could not decode response inside the cubbyhole")
	}

	return response, nil
}

func (b *SystemBackend) handleMetrics(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	format := data.Get("format").(string)
	if format == "" {
		format = metricsutil.FormatFromRequest(req)
	}
	return b.Core.metricsHelper.ResponseForFormat(format), nil
}

func (b *SystemBackend) handleMonitor(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	ll := data.Get("log_level").(string)
	w := req.ResponseWriter

	if ll == "" {
		ll = "info"
	}
	logLevel := log.LevelFromString(ll)

	if logLevel == log.NoLevel {
		return logical.ErrorResponse("unknown log level"), nil
	}

	flusher, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		return logical.ErrorResponse("streaming not supported"), nil
	}

	isJson := b.Core.LogFormat() == "json"
	logger := b.Core.Logger().(log.InterceptLogger)

	mon, err := monitor.NewMonitor(512, logger, &log.LoggerOptions{
		Level:      logLevel,
		JSONFormat: isJson,
	})
	if err != nil {
		return nil, err
	}

	logCh := mon.Start()
	defer mon.Stop()

	if logCh == nil {
		return nil, fmt.Errorf("error trying to start a monitor that's already been started")
	}

	w.WriteHeader(http.StatusOK)

	// 0 byte write is needed before the Flush call so that if we are using
	// a gzip stream it will go ahead and write out the HTTP response header
	_, err = w.Write([]byte(""))
	if err != nil {
		return nil, fmt.Errorf("error seeding flusher: %w", err)
	}

	flusher.Flush()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Stream logs until the connection is closed.
	for {
		select {
		// Periodically check for the seal status and return if core gets
		// marked as sealed.
		case <-ticker.C:
			if b.Core.Sealed() {
				// We still return the error, but this will be ignored upstream
				// due to the fact that we've already sent a response by
				// writing the header and flushing the writer above.
				_, err = fmt.Fprint(w, "core received sealed state change, ending monitor session")
				if err != nil {
					return nil, fmt.Errorf("error checking seal state: %w", err)
				}
			}
		case <-ctx.Done():
			return nil, nil
		case l := <-logCh:
			// We still return the error, but this will be ignored upstream
			// due to the fact that we've already sent a response by
			// writing the header and flushing the writer above.
			_, err = fmt.Fprint(w, string(l))
			if err != nil {
				return nil, fmt.Errorf("error streaming monitor output: %w", err)
			}

			flusher.Flush()
		}
	}
}

// handleHostInfo collects and returns host-related information, which includes
// system information, cpu, disk, and memory usage. Any capture-related errors
// returned by the collection method will be returned as response warnings.
func (b *SystemBackend) handleHostInfo(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	resp := &logical.Response{}
	info, err := hostutil.CollectHostInfo(ctx)
	if err != nil {
		// If the error is a HostInfoError, we return them as response warnings
		if errs, ok := err.(*multierror.Error); ok {
			var warnings []string
			for _, mErr := range errs.Errors {
				if errwrap.ContainsType(mErr, new(hostutil.HostInfoError)) {
					warnings = append(warnings, mErr.Error())
				} else {
					// If the error is a multierror, it should only be for
					// HostInfoError, but if it's not for any reason, we return
					// it as an error to avoid it being swallowed.
					return nil, err
				}
			}
			resp.Warnings = warnings
		} else {
			return nil, err
		}
	}

	if info == nil {
		return nil, errors.New("unable to collect host information: nil HostInfo")
	}

	respData := map[string]interface{}{
		"timestamp": info.Timestamp,
	}
	if info.CPU != nil {
		respData["cpu"] = info.CPU
	}
	if info.CPUTimes != nil {
		respData["cpu_times"] = info.CPUTimes
	}
	if info.Disk != nil {
		respData["disk"] = info.Disk
	}
	if info.Host != nil {
		respData["host"] = info.Host
	}
	if info.Memory != nil {
		respData["memory"] = info.Memory
	}
	resp.Data = respData

	return resp, nil
}

func (b *SystemBackend) handleWrappingLookup(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// This ordering of lookups has been validated already in the wrapping
	// validation func, we're just doing this for a safety check
	token := data.Get("token").(string)
	if token == "" {
		token = req.ClientToken
		if token == "" {
			return logical.ErrorResponse("missing \"token\" value in input"), logical.ErrInvalidRequest
		}
	}

	te, err := b.Core.tokenStore.lookupTainted(ctx, token)
	if err != nil {
		return nil, err
	}
	if te == nil {
		return nil, nil
	}
	if len(te.Policies) != 1 {
		return nil, errors.New("token is not a valid unwrap token")
	}

	cubbyReq := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "cubbyhole/wrapinfo",
		ClientToken: token,
	}
	cubbyReq.SetTokenEntry(te)
	cubbyResp, err := b.Core.router.Route(ctx, cubbyReq)
	if err != nil {
		return nil, fmt.Errorf("error looking up wrapping information: %w", err)
	}
	if cubbyResp == nil {
		return logical.ErrorResponse("no information found; wrapping token may be from a previous Vault version"), nil
	}
	if cubbyResp != nil && cubbyResp.IsError() {
		return cubbyResp, nil
	}
	if cubbyResp.Data == nil {
		return logical.ErrorResponse("wrapping information was nil; wrapping token may be from a previous Vault version"), nil
	}

	creationTTLRaw := cubbyResp.Data["creation_ttl"]
	creationTime := cubbyResp.Data["creation_time"]
	creationPath := cubbyResp.Data["creation_path"]

	resp := &logical.Response{
		Data: map[string]interface{}{},
	}
	if creationTTLRaw != nil {
		creationTTL, err := creationTTLRaw.(json.Number).Int64()
		if err != nil {
			return nil, fmt.Errorf("error reading creation_ttl value from wrapping information: %w", err)
		}
		resp.Data["creation_ttl"] = time.Duration(creationTTL).Seconds()
	}
	if creationTime != nil {
		// This was JSON marshaled so it's already a string in RFC3339 format
		resp.Data["creation_time"] = cubbyResp.Data["creation_time"]
	}
	if creationPath != nil {
		resp.Data["creation_path"] = cubbyResp.Data["creation_path"]
	}

	return resp, nil
}

func (b *SystemBackend) handleWrappingRewrap(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// If a third party is rewrapping (rather than the calling token being the
	// wrapping token) we detect this so that we can revoke the original
	// wrapping token after reading it. Right now wrapped tokens can't unwrap
	// themselves, but in case we change it, this will be ready to do the right
	// thing.
	var thirdParty bool

	token := data.Get("token").(string)
	if token != "" {
		thirdParty = true
	} else {
		token = req.ClientToken
	}

	te, err := b.Core.tokenStore.lookupTainted(ctx, token)
	if err != nil {
		return nil, err
	}
	if te == nil {
		return nil, nil
	}
	if len(te.Policies) != 1 {
		return nil, errors.New("token is not a valid unwrap token")
	}

	if thirdParty {
		// Use the token to decrement the use count to avoid a second operation on the token.
		_, err := b.Core.tokenStore.UseTokenByID(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("error decrementing wrapping token's use-count: %w", err)
		}
		defer b.Core.tokenStore.revokeOrphan(ctx, token)
	}

	// Fetch the original TTL
	cubbyReq := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "cubbyhole/wrapinfo",
		ClientToken: token,
	}
	cubbyReq.SetTokenEntry(te)
	cubbyResp, err := b.Core.router.Route(ctx, cubbyReq)
	if err != nil {
		return nil, fmt.Errorf("error looking up wrapping information: %w", err)
	}
	if cubbyResp == nil {
		return logical.ErrorResponse("no information found; wrapping token may be from a previous Vault version"), nil
	}
	if cubbyResp != nil && cubbyResp.IsError() {
		return cubbyResp, nil
	}
	if cubbyResp.Data == nil {
		return logical.ErrorResponse("wrapping information was nil; wrapping token may be from a previous Vault version"), nil
	}

	// Set the creation TTL on the request
	creationTTLRaw := cubbyResp.Data["creation_ttl"]
	if creationTTLRaw == nil {
		return nil, fmt.Errorf("creation_ttl value in wrapping information was nil")
	}
	creationTTL, err := cubbyResp.Data["creation_ttl"].(json.Number).Int64()
	if err != nil {
		return nil, fmt.Errorf("error reading creation_ttl value from wrapping information: %w", err)
	}

	// Get creation_path to return as the response later
	creationPathRaw := cubbyResp.Data["creation_path"]
	if creationPathRaw == nil {
		return nil, fmt.Errorf("creation_path value in wrapping information was nil")
	}
	creationPath := creationPathRaw.(string)

	// Fetch the original response and return it as the data for the new response
	cubbyReq = &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "cubbyhole/response",
		ClientToken: token,
	}
	cubbyReq.SetTokenEntry(te)
	cubbyResp, err = b.Core.router.Route(ctx, cubbyReq)
	if err != nil {
		return nil, fmt.Errorf("error looking up response: %w", err)
	}
	if cubbyResp == nil {
		return logical.ErrorResponse("no information found; wrapping token may be from a previous Vault version"), nil
	}
	if cubbyResp != nil && cubbyResp.IsError() {
		return cubbyResp, nil
	}
	if cubbyResp.Data == nil {
		return logical.ErrorResponse("wrapping information was nil; wrapping token may be from a previous Vault version"), nil
	}

	response := cubbyResp.Data["response"]
	if response == nil {
		return nil, fmt.Errorf("no response found inside the cubbyhole")
	}

	// Return response in "response"; wrapping code will detect the rewrap and
	// slot in instead of nesting
	return &logical.Response{
		Data: map[string]interface{}{
			"response": response,
		},
		WrapInfo: &wrapping.ResponseWrapInfo{
			TTL:          time.Duration(creationTTL),
			CreationPath: creationPath,
		},
	}, nil
}

func (b *SystemBackend) pathHashWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	inputB64 := d.Get("input").(string)
	format := d.Get("format").(string)
	algorithm := d.Get("urlalgorithm").(string)
	if algorithm == "" {
		algorithm = d.Get("algorithm").(string)
	}

	input, err := base64.StdEncoding.DecodeString(inputB64)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("unable to decode input as base64: %s", err)), logical.ErrInvalidRequest
	}

	switch format {
	case "hex":
	case "base64":
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported encoding format %s; must be \"hex\" or \"base64\"", format)), nil
	}

	var hf hash.Hash
	switch algorithm {
	case "sha2-224":
		hf = sha256.New224()
	case "sha2-256":
		hf = sha256.New()
	case "sha2-384":
		hf = sha512.New384()
	case "sha2-512":
		hf = sha512.New()
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported algorithm %s", algorithm)), nil
	}
	hf.Write(input)
	retBytes := hf.Sum(nil)

	var retStr string
	switch format {
	case "hex":
		retStr = hex.EncodeToString(retBytes)
	case "base64":
		retStr = base64.StdEncoding.EncodeToString(retBytes)
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"sum": retStr,
		},
	}
	return resp, nil
}

func (b *SystemBackend) pathRandomWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	bytes := 0
	var err error
	strBytes := d.Get("urlbytes").(string)
	if strBytes != "" {
		bytes, err = strconv.Atoi(strBytes)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("error parsing url-set byte count: %s", err)), nil
		}
	} else {
		bytes = d.Get("bytes").(int)
	}
	format := d.Get("format").(string)

	if bytes < 1 {
		return logical.ErrorResponse(`"bytes" cannot be less than 1`), nil
	}

	if bytes > maxBytes {
		return logical.ErrorResponse(`"bytes" should be less than %d`, maxBytes), nil
	}

	switch format {
	case "hex":
	case "base64":
	default:
		return logical.ErrorResponse("unsupported encoding format %q; must be \"hex\" or \"base64\"", format), nil
	}

	randBytes, err := uuid.GenerateRandomBytes(bytes)
	if err != nil {
		return nil, err
	}

	var retStr string
	switch format {
	case "hex":
		retStr = hex.EncodeToString(randBytes)
	case "base64":
		retStr = base64.StdEncoding.EncodeToString(randBytes)
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"random_bytes": retStr,
		},
	}
	return resp, nil
}

func hasMountAccess(ctx context.Context, acl *ACL, path string) bool {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return false
	}

	// If a policy is giving us direct access to the mount path then we can do
	// a fast return.
	capabilities := acl.Capabilities(ctx, ns.TrimmedPath(path))
	if !strutil.StrListContains(capabilities, DenyCapability) {
		return true
	}

	var aclCapabilitiesGiven bool
	walkFn := func(s string, v interface{}) bool {
		if v == nil {
			return false
		}

		perms := v.(*ACLPermissions)

		switch {
		case perms.CapabilitiesBitmap&DenyCapabilityInt > 0:
			return false

		case perms.CapabilitiesBitmap&CreateCapabilityInt > 0,
			perms.CapabilitiesBitmap&DeleteCapabilityInt > 0,
			perms.CapabilitiesBitmap&ListCapabilityInt > 0,
			perms.CapabilitiesBitmap&ReadCapabilityInt > 0,
			perms.CapabilitiesBitmap&SudoCapabilityInt > 0,
			perms.CapabilitiesBitmap&UpdateCapabilityInt > 0:

			aclCapabilitiesGiven = true

			return true
		}

		return false
	}

	acl.exactRules.WalkPrefix(path, walkFn)
	if !aclCapabilitiesGiven {
		acl.prefixRules.WalkPrefix(path, walkFn)
	}

	if !aclCapabilitiesGiven {
		if perms := acl.CheckAllowedFromNonExactPaths(path, true); perms != nil {
			return true
		}
	}

	return aclCapabilitiesGiven
}

func (b *SystemBackend) pathInternalUIMountsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}

	secretMounts := make(map[string]interface{})
	authMounts := make(map[string]interface{})
	resp.Data["secret"] = secretMounts
	resp.Data["auth"] = authMounts

	var acl *ACL
	var isAuthed bool
	if req.ClientToken != "" {
		isAuthed = true

		var entity *identity.Entity
		var te *logical.TokenEntry
		// Load the ACL policies so we can walk the prefix for this mount
		acl, te, entity, _, err = b.Core.fetchACLTokenEntryAndEntity(ctx, req)
		if err != nil {
			return nil, err
		}
		if entity != nil && entity.Disabled {
			b.logger.Warn("permission denied as the entity on the token is disabled")
			return nil, logical.ErrPermissionDenied
		}
		if te != nil && te.EntityID != "" && entity == nil {
			b.logger.Warn("permission denied as the entity on the token is invalid")
			return nil, logical.ErrPermissionDenied
		}
	}

	hasAccess := func(ctx context.Context, me *MountEntry) bool {
		if me.Config.ListingVisibility == ListingVisibilityUnauth {
			return true
		}

		if isAuthed {
			return hasMountAccess(ctx, acl, me.Namespace().Path+me.Path)
		}

		return false
	}

	b.Core.mountsLock.RLock()
	for _, entry := range b.Core.mounts.Entries {
		ctxWithNamespace := namespace.ContextWithNamespace(ctx, entry.Namespace())
		filtered, err := b.Core.checkReplicatedFiltering(ctxWithNamespace, entry, "")
		if err != nil {
			b.Core.mountsLock.RUnlock()
			return nil, err
		}
		if filtered {
			continue
		}

		if ns.ID == entry.NamespaceID && hasAccess(ctx, entry) {
			if isAuthed {
				// If this is an authed request return all the mount info
				secretMounts[entry.Path] = mountInfo(entry)
			} else {
				secretMounts[entry.Path] = map[string]interface{}{
					"type":        entry.Type,
					"description": entry.Description,
					"options":     entry.Options,
				}
			}
		}
	}
	b.Core.mountsLock.RUnlock()

	b.Core.authLock.RLock()
	for _, entry := range b.Core.auth.Entries {
		ctxWithNamespace := namespace.ContextWithNamespace(ctx, entry.Namespace())
		filtered, err := b.Core.checkReplicatedFiltering(ctxWithNamespace, entry, credentialRoutePrefix)
		if err != nil {
			b.Core.authLock.RUnlock()
			return nil, err
		}
		if filtered {
			continue
		}

		if ns.ID == entry.NamespaceID && hasAccess(ctx, entry) {
			if isAuthed {
				// If this is an authed request return all the mount info
				authMounts[entry.Path] = mountInfo(entry)
			} else {
				authMounts[entry.Path] = map[string]interface{}{
					"type":        entry.Type,
					"description": entry.Description,
					"options":     entry.Options,
				}
			}
		}
	}
	b.Core.authLock.RUnlock()

	return resp, nil
}

func (b *SystemBackend) pathInternalUIMountRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	path := d.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse("path not set"), logical.ErrInvalidRequest
	}
	path = sanitizePath(path)

	// Load the ACL policies so we can walk the prefix for this mount
	acl, te, entity, _, err := b.Core.fetchACLTokenEntryAndEntity(ctx, req)
	if err != nil {
		return nil, err
	}
	if entity != nil && entity.Disabled {
		b.logger.Warn("permission denied as the entity on the token is disabled")
		return nil, logical.ErrPermissionDenied
	}
	if te != nil && te.EntityID != "" && entity == nil {
		b.logger.Warn("permission denied as the entity on the token is invalid")
		return nil, logical.ErrPermissionDenied
	}

	errResp := logical.ErrorResponse(fmt.Sprintf("preflight capability check returned 403, please ensure client's policies grant access to path %q", path))

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	me := b.Core.router.MatchingMountEntry(ctx, path)
	if me == nil {
		// Return a permission denied error here so this path cannot be used to
		// brute force a list of mounts.
		return errResp, logical.ErrPermissionDenied
	}

	filtered, err := b.Core.checkReplicatedFiltering(ctx, me, "")
	if err != nil {
		return nil, err
	}
	if filtered {
		return errResp, logical.ErrPermissionDenied
	}
	resp := &logical.Response{
		Data: mountInfo(me),
	}
	resp.Data["path"] = me.Path

	fullMountPath := ns.Path + me.Path
	if ns.ID != me.Namespace().ID {
		resp.Data["path"] = me.Namespace().Path + me.Path
		fullMountPath = ns.Path + me.Namespace().Path + me.Path
	}

	if !hasMountAccess(ctx, acl, fullMountPath) {
		return errResp, logical.ErrPermissionDenied
	}

	return resp, nil
}

func (b *SystemBackend) pathInternalCountersRequests(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	counters, err := b.Core.loadAllRequestCounters(ctx, time.Now())
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"counters": counters,
		},
	}

	return resp, nil
}

func (b *SystemBackend) pathInternalCountersTokens(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	activeTokens, err := b.Core.countActiveTokens(ctx)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"counters": activeTokens,
		},
	}

	return resp, nil
}

func (b *SystemBackend) pathInternalCountersEntities(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	activeEntities, err := b.Core.countActiveEntities(ctx)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"counters": activeEntities,
		},
	}

	return resp, nil
}

func (b *SystemBackend) pathInternalUIResultantACL(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		// 204 -- no ACL
		return nil, nil
	}

	acl, te, entity, _, err := b.Core.fetchACLTokenEntryAndEntity(ctx, req)
	if err != nil {
		return nil, err
	}

	if entity != nil && entity.Disabled {
		b.logger.Warn("permission denied as the entity on the token is disabled")
		return logical.ErrorResponse(logical.ErrPermissionDenied.Error()), nil
	}
	if te != nil && te.EntityID != "" && entity == nil {
		b.logger.Warn("permission denied as the entity on the token is invalid")
		return logical.ErrorResponse(logical.ErrPermissionDenied.Error()), nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"root": false,
		},
	}

	if acl.root {
		resp.Data["root"] = true
		return resp, nil
	}

	exact := map[string]interface{}{}
	glob := map[string]interface{}{}

	walkFn := func(pt map[string]interface{}, s string, v interface{}) {
		if v == nil {
			return
		}

		perms := v.(*ACLPermissions)
		capabilities := []string{}

		if perms.CapabilitiesBitmap&CreateCapabilityInt > 0 {
			capabilities = append(capabilities, CreateCapability)
		}
		if perms.CapabilitiesBitmap&DeleteCapabilityInt > 0 {
			capabilities = append(capabilities, DeleteCapability)
		}
		if perms.CapabilitiesBitmap&ListCapabilityInt > 0 {
			capabilities = append(capabilities, ListCapability)
		}
		if perms.CapabilitiesBitmap&ReadCapabilityInt > 0 {
			capabilities = append(capabilities, ReadCapability)
		}
		if perms.CapabilitiesBitmap&SudoCapabilityInt > 0 {
			capabilities = append(capabilities, SudoCapability)
		}
		if perms.CapabilitiesBitmap&UpdateCapabilityInt > 0 {
			capabilities = append(capabilities, UpdateCapability)
		}

		// If "deny" is explicitly set or if the path has no capabilities at all,
		// set the path capabilities to "deny"
		if perms.CapabilitiesBitmap&DenyCapabilityInt > 0 || len(capabilities) == 0 {
			capabilities = []string{DenyCapability}
		}

		res := map[string]interface{}{}
		if len(capabilities) > 0 {
			res["capabilities"] = capabilities
		}
		if perms.MinWrappingTTL != 0 {
			res["min_wrapping_ttl"] = int64(perms.MinWrappingTTL.Seconds())
		}
		if perms.MaxWrappingTTL != 0 {
			res["max_wrapping_ttl"] = int64(perms.MaxWrappingTTL.Seconds())
		}
		if len(perms.AllowedParameters) > 0 {
			res["allowed_parameters"] = perms.AllowedParameters
		}
		if len(perms.DeniedParameters) > 0 {
			res["denied_parameters"] = perms.DeniedParameters
		}
		if len(perms.RequiredParameters) > 0 {
			res["required_parameters"] = perms.RequiredParameters
		}

		pt[s] = res
	}

	exactWalkFn := func(s string, v interface{}) bool {
		walkFn(exact, s, v)
		return false
	}

	globWalkFn := func(s string, v interface{}) bool {
		walkFn(glob, s, v)
		return false
	}

	acl.exactRules.Walk(exactWalkFn)
	acl.prefixRules.Walk(globWalkFn)

	resp.Data["exact_paths"] = exact
	resp.Data["glob_paths"] = glob

	return resp, nil
}

func (b *SystemBackend) pathInternalOpenAPI(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Limit output to authorized paths
	resp, err := b.pathInternalUIMountsRead(ctx, req, d)
	if err != nil {
		return nil, err
	}

	context := d.Get("context").(string)

	// Set up target document and convert to map[string]interface{} which is what will
	// be received from plugin backends.
	doc := framework.NewOASDocument()

	procMountGroup := func(group, mountPrefix string) error {
		for mount := range resp.Data[group].(map[string]interface{}) {
			backend := b.Core.router.MatchingBackend(ctx, mountPrefix+mount)

			if backend == nil {
				continue
			}

			req := &logical.Request{
				Operation: logical.HelpOperation,
				Storage:   req.Storage,
			}

			resp, err := backend.HandleRequest(ctx, req)
			if err != nil {
				return err
			}

			var backendDoc *framework.OASDocument

			// Normalize response type, which will be different if received
			// from an external plugin.
			switch v := resp.Data["openapi"].(type) {
			case *framework.OASDocument:
				backendDoc = v
			case map[string]interface{}:
				backendDoc, err = framework.NewOASDocumentFromMap(v)
				if err != nil {
					return err
				}
			default:
				continue
			}

			// Prepare to add tags to default builtins that are
			// type "unknown" and won't already be tagged.
			var tag string
			switch mountPrefix + mount {
			case "cubbyhole/", "secret/":
				tag = "secrets"
			case "sys/":
				tag = "system"
			case "auth/token/":
				tag = "auth"
			case "identity/":
				tag = "identity"
			}

			// Merge backend paths with existing document
			for path, obj := range backendDoc.Paths {
				path := strings.TrimPrefix(path, "/")

				// Add tags to all of the operations if necessary
				if tag != "" {
					for _, op := range []*framework.OASOperation{obj.Get, obj.Post, obj.Delete} {
						// TODO: a special override for identity is used used here because the backend
						// is currently categorized as "secret", which will likely change. Also of interest
						// is removing all tag handling here and providing the mount information to OpenAPI.
						if op != nil && (len(op.Tags) == 0 || tag == "identity") {
							op.Tags = []string{tag}
						}
					}
				}

				doc.Paths["/"+mountPrefix+mount+path] = obj
			}
		}
		return nil
	}

	if err := procMountGroup("secret", ""); err != nil {
		return nil, err
	}
	if err := procMountGroup("auth", "auth/"); err != nil {
		return nil, err
	}

	doc.CreateOperationIDs(context)

	buf, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	resp = &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPStatusCode:  200,
			logical.HTTPRawBody:     buf,
			logical.HTTPContentType: "application/json",
		},
	}

	return resp, nil
}

type SealStatusResponse struct {
	Type         string `json:"type"`
	Initialized  bool   `json:"initialized"`
	Sealed       bool   `json:"sealed"`
	T            int    `json:"t"`
	N            int    `json:"n"`
	Progress     int    `json:"progress"`
	Nonce        string `json:"nonce"`
	Version      string `json:"version"`
	Migration    bool   `json:"migration"`
	ClusterName  string `json:"cluster_name,omitempty"`
	ClusterID    string `json:"cluster_id,omitempty"`
	RecoverySeal bool   `json:"recovery_seal"`
	StorageType  string `json:"storage_type,omitempty"`
}

func (core *Core) GetSealStatus(ctx context.Context) (*SealStatusResponse, error) {
	sealed := core.Sealed()

	initialized, err := core.Initialized(ctx)
	if err != nil {
		return nil, err
	}

	var sealConfig *SealConfig
	if core.SealAccess().RecoveryKeySupported() {
		sealConfig, err = core.SealAccess().RecoveryConfig(ctx)
	} else {
		sealConfig, err = core.SealAccess().BarrierConfig(ctx)
	}
	if err != nil {
		return nil, err
	}

	if sealConfig == nil {
		return &SealStatusResponse{
			Type:         core.SealAccess().BarrierType(),
			Initialized:  initialized,
			Sealed:       true,
			RecoverySeal: core.SealAccess().RecoveryKeySupported(),
			StorageType:  core.StorageType(),
			Version:      version.GetVersion().VersionNumber(),
		}, nil
	}

	// Fetch the local cluster name and identifier
	var clusterName, clusterID string
	if !sealed {
		cluster, err := core.Cluster(ctx)
		if err != nil {
			return nil, err
		}
		if cluster == nil {
			return nil, fmt.Errorf("failed to fetch cluster details")
		}
		clusterName = cluster.Name
		clusterID = cluster.ID
	}

	progress, nonce := core.SecretProgress()

	return &SealStatusResponse{
		Type:         sealConfig.Type,
		Initialized:  initialized,
		Sealed:       sealed,
		T:            sealConfig.SecretThreshold,
		N:            sealConfig.SecretShares,
		Progress:     progress,
		Nonce:        nonce,
		Version:      version.GetVersion().VersionNumber(),
		Migration:    core.IsInSealMigrationMode() && !core.IsSealMigrated(),
		ClusterName:  clusterName,
		ClusterID:    clusterID,
		RecoverySeal: core.SealAccess().RecoveryKeySupported(),
		StorageType:  core.StorageType(),
	}, nil
}

type LeaderResponse struct {
	HAEnabled                bool      `json:"ha_enabled"`
	IsSelf                   bool      `json:"is_self"`
	ActiveTime               time.Time `json:"active_time,omitempty"`
	LeaderAddress            string    `json:"leader_address"`
	LeaderClusterAddress     string    `json:"leader_cluster_address"`
	PerfStandby              bool      `json:"performance_standby"`
	PerfStandbyLastRemoteWAL uint64    `json:"performance_standby_last_remote_wal"`
	LastWAL                  uint64    `json:"last_wal,omitempty"`

	// Raft Indexes for this node
	RaftCommittedIndex uint64 `json:"raft_committed_index,omitempty"`
	RaftAppliedIndex   uint64 `json:"raft_applied_index,omitempty"`
}

func (core *Core) GetLeaderStatus() (*LeaderResponse, error) {
	haEnabled := true
	isLeader, address, clusterAddr, err := core.Leader()
	if errwrap.Contains(err, ErrHANotEnabled.Error()) {
		haEnabled = false
		err = nil
	}
	if err != nil {
		return nil, err
	}

	resp := &LeaderResponse{
		HAEnabled:            haEnabled,
		IsSelf:               isLeader,
		LeaderAddress:        address,
		LeaderClusterAddress: clusterAddr,
		PerfStandby:          core.PerfStandby(),
	}
	if isLeader {
		resp.ActiveTime = core.ActiveTime()
	}
	if resp.PerfStandby {
		resp.PerfStandbyLastRemoteWAL = LastRemoteWAL(core)
	} else if isLeader || !haEnabled {
		resp.LastWAL = LastWAL(core)
	}

	resp.RaftCommittedIndex, resp.RaftAppliedIndex = core.GetRaftIndexes()
	return resp, nil
}

func (b *SystemBackend) handleSealStatus(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	status, err := b.Core.GetSealStatus(ctx)
	if err != nil {
		return nil, err
	}
	buf, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}
	httpResp := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPStatusCode:  200,
			logical.HTTPRawBody:     buf,
			logical.HTTPContentType: "application/json",
		},
	}
	return httpResp, nil
}

func (b *SystemBackend) handleLeaderStatus(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	status, err := b.Core.GetLeaderStatus()
	if err != nil {
		return nil, err
	}
	buf, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}
	httpResp := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPStatusCode:  200,
			logical.HTTPRawBody:     buf,
			logical.HTTPContentType: "application/json",
		},
	}
	return httpResp, nil
}

func (b *SystemBackend) verifyDROperationTokenOnSecondary(f framework.OperationFunc, lock bool) framework.OperationFunc {
	if b.Core.IsDRSecondary() {
		return b.verifyDROperationToken(f, lock)
	}
	return f
}

func (b *SystemBackend) rotateBarrierKey(ctx context.Context) error {
	// Rotate to the new term
	newTerm, err := b.Core.barrier.Rotate(ctx, b.Core.secureRandomReader)
	if err != nil {
		return errwrap.Wrap(errors.New("failed to create new encryption key"), err)
	}
	b.Backend.Logger().Info("installed new encryption key")

	// In HA mode, we need to an upgrade path for the standby instances
	if b.Core.ha != nil && b.Core.KeyRotateGracePeriod() > 0 {
		// Create the upgrade path to the new term
		if err := b.Core.barrier.CreateUpgrade(ctx, newTerm); err != nil {
			b.Backend.Logger().Error("failed to create new upgrade", "term", newTerm, "error", err)
		}

		// Schedule the destroy of the upgrade path
		time.AfterFunc(b.Core.KeyRotateGracePeriod(), func() {
			b.Backend.Logger().Debug("cleaning up upgrade keys", "waited", b.Core.KeyRotateGracePeriod())
			if err := b.Core.barrier.DestroyUpgrade(b.Core.activeContext, newTerm); err != nil {
				b.Backend.Logger().Error("failed to destroy upgrade", "term", newTerm, "error", err)
			}
		})
	}

	// Write to the canary path, which will force a synchronous truing during
	// replication
	if err := b.Core.barrier.Put(ctx, &logical.StorageEntry{
		Key:   coreKeyringCanaryPath,
		Value: []byte(fmt.Sprintf("new-rotation-term-%d", newTerm)),
	}); err != nil {
		b.Core.logger.Error("error saving keyring canary", "error", err)
		return errwrap.Wrap(errors.New("failed to save keyring canary"), err)
	}

	return nil
}

func sanitizePath(path string) string {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	return path
}

func checkListingVisibility(visibility ListingVisibilityType) error {
	switch visibility {
	case ListingVisibilityDefault:
	case ListingVisibilityHidden:
	case ListingVisibilityUnauth:
	default:
		return fmt.Errorf("invalid listing visilibity type")
	}

	return nil
}

const sysHelpRoot = `
The system backend is built-in to Vault and cannot be remounted or
unmounted. It contains the paths that are used to configure Vault itself
as well as perform core operations.
`

// sysHelp is all the help text for the sys backend.
var sysHelp = map[string][2]string{
	"license": {
		"Sets the license of the server.",
		`
The path responds to the following HTTP methods.

    GET /
        Returns information on the installed license

    POST
        Sets the license for the server
	`,
	},
	"config/cors": {
		"Configures or returns the current configuration of CORS settings.",
		`
This path responds to the following HTTP methods.

    GET /
        Returns the configuration of the CORS setting.

    POST /
        Sets the comma-separated list of origins that can make cross-origin requests.

    DELETE /
        Clears the CORS configuration and disables acceptance of CORS requests.
		`,
	},
	"config/ui/headers": {
		"Configures response headers that should be returned from the UI.",
		`
This path responds to the following HTTP methods.
    GET /<header>
        Returns the header value.
    POST /<header>
        Sets the header value for the UI.
    DELETE /<header>
        Clears the header value for UI.

    LIST /
        List the headers configured for the UI.
        `,
	},
	"init": {
		"Initializes or returns the initialization status of the Vault.",
		`
This path responds to the following HTTP methods.

    GET /
        Returns the initialization status of the Vault.

    POST /
        Initializes a new vault.
		`,
	},
	"health": {
		"Checks the health status of the Vault.",
		`
This path responds to the following HTTP methods.

	GET /
		Returns health information about the Vault.
		`,
	},
	"generate-root": {
		"Reads, generates, or deletes a root token regeneration process.",
		`
This path responds to multiple HTTP methods which change the behavior. Those
HTTP methods are listed below.

    GET /attempt
        Reads the configuration and progress of the current root generation
        attempt.

    POST /attempt
        Initializes a new root generation attempt. Only a single root generation
        attempt can take place at a time. One (and only one) of otp or pgp_key
        are required.

    DELETE /attempt
        Cancels any in-progress root generation attempt. This clears any
        progress made. This must be called to change the OTP or PGP key being
        used.
		`,
	},
	"seal-status": {
		"Returns the seal status of the Vault.",
		`
This path responds to the following HTTP methods.

    GET /
        Returns the seal status of the Vault. This is an unauthenticated
        endpoint.
		`,
	},
	"seal": {
		"Seals the Vault.",
		`
This path responds to the following HTTP methods.

    PUT /
        Seals the Vault.
		`,
	},
	"unseal": {
		"Unseals the Vault.",
		`
This path responds to the following HTTP methods.

    PUT /
        Unseals the Vault.
		`,
	},
	"mounts": {
		"List the currently mounted backends.",
		`
This path responds to the following HTTP methods.

    GET /
        Lists all the mounted secret backends.

    GET /<mount point>
        Get information about the mount at the specified path.

    POST /<mount point>
        Mount a new secret backend to the mount point in the URL.

    POST /<mount point>/tune
        Tune configuration parameters for the given mount point.

    DELETE /<mount point>
        Unmount the specified mount point.
		`,
	},

	"mount": {
		`Mount a new backend at a new path.`,
		`
Mount a backend at a new path. A backend can be mounted multiple times at
multiple paths in order to configure multiple separately configured backends.
Example: you might have an AWS backend for the east coast, and one for the
west coast.
		`,
	},

	"mount_path": {
		`The path to mount to. Example: "aws/east"`,
		"",
	},

	"mount_type": {
		`The type of the backend. Example: "passthrough"`,
		"",
	},

	"mount_desc": {
		`User-friendly description for this mount.`,
		"",
	},

	"mount_config": {
		`Configuration for this mount, such as default_lease_ttl
and max_lease_ttl.`,
	},

	"mount_local": {
		`Mark the mount as a local mount, which is not replicated
and is unaffected by replication.`,
	},

	"mount_plugin_name": {
		`Name of the plugin to mount based from the name registered
in the plugin catalog.`,
	},

	"mount_options": {
		`The options to pass into the backend. Should be a json object with string keys and values.`,
	},

	"seal_wrap": {
		`Whether to turn on seal wrapping for the mount.`,
	},

	"external_entropy_access": {
		`Whether to give the mount access to Vault's external entropy.`,
	},

	"tune_default_lease_ttl": {
		`The default lease TTL for this mount.`,
	},

	"tune_max_lease_ttl": {
		`The max lease TTL for this mount.`,
	},

	"tune_audit_non_hmac_request_keys": {
		`The list of keys in the request data object that will not be HMAC'ed by audit devices.`,
	},

	"tune_audit_non_hmac_response_keys": {
		`The list of keys in the response data object that will not be HMAC'ed by audit devices.`,
	},

	"tune_mount_options": {
		`The options to pass into the backend. Should be a json object with string keys and values.`,
	},

	"remount": {
		"Move the mount point of an already-mounted backend.",
		`
This path responds to the following HTTP methods.

    POST /sys/remount
        Changes the mount point of an already-mounted backend.
		`,
	},

	"auth_tune": {
		"Tune the configuration parameters for an auth path.",
		`Read and write the 'default-lease-ttl' and 'max-lease-ttl' values of
the auth path.`,
	},

	"mount_tune": {
		"Tune backend configuration parameters for this mount.",
		`Read and write the 'default-lease-ttl' and 'max-lease-ttl' values of
the mount.`,
	},

	"renew": {
		"Renew a lease on a secret",
		`
When a secret is read, it may optionally include a lease interval
and a boolean indicating if renew is possible. For secrets that support
lease renewal, this endpoint is used to extend the validity of the
lease and to prevent an automatic revocation.
		`,
	},

	"lease_id": {
		"The lease identifier to renew. This is included with a lease.",
		"",
	},

	"increment": {
		"The desired increment in seconds to the lease",
		"",
	},

	"revoke": {
		"Revoke a leased secret immediately",
		`
When a secret is generated with a lease, it is automatically revoked
at the end of the lease period if not renewed. However, in some cases
you may want to force an immediate revocation. This endpoint can be
used to revoke the secret with the given Lease ID.
		`,
	},

	"revoke-sync": {
		"Whether or not to perform the revocation synchronously",
		`
If false, the call will return immediately and revocation will be queued; if it
fails, Vault will keep trying. If true, if the revocation fails, Vault will not
automatically try again and will return an error. For revoke-prefix, this
setting will apply to all leases being revoked. For revoke-force, since errors
are ignored, this setting is not supported.
`,
	},

	"revoke-prefix": {
		"Revoke all secrets generated in a given prefix",
		`
Revokes all the secrets generated under a given mount prefix. As
an example, "prod/aws/" might be the AWS logical backend, and due to
a change in the "ops" policy, we may want to invalidate all the secrets
generated. We can do a revoke prefix at "prod/aws/ops" to revoke all
the ops secrets. This does a prefix match on the Lease IDs and revokes
all matching leases.
		`,
	},

	"revoke-prefix-path": {
		`The path to revoke keys under. Example: "prod/aws/ops"`,
		"",
	},

	"revoke-force": {
		"Revoke all secrets generated in a given prefix, ignoring errors.",
		`
See the path help for 'revoke-prefix'; this behaves the same, except that it
ignores errors encountered during revocation. This can be used in certain
recovery situations; for instance, when you want to unmount a backend, but it
is impossible to fix revocation errors and these errors prevent the unmount
from proceeding. This is a DANGEROUS operation as it removes Vault's oversight
of external secrets. Access to this prefix should be tightly controlled.
		`,
	},

	"revoke-force-path": {
		`The path to revoke keys under. Example: "prod/aws/ops"`,
		"",
	},

	"auth-table": {
		"List the currently enabled credential backends.",
		`
This path responds to the following HTTP methods.

    GET /
        List the currently enabled credential backends: the name, the type of
        the backend, and a user friendly description of the purpose for the
        credential backend.

    POST /<mount point>
        Enable a new auth method.

    DELETE /<mount point>
        Disable the auth method at the given mount point.
		`,
	},

	"auth": {
		`Enable a new credential backend with a name.`,
		`
Enable a credential mechanism at a new path. A backend can be mounted multiple times at
multiple paths in order to configure multiple separately configured backends.
Example: you might have an OAuth backend for GitHub, and one for Google Apps.
		`,
	},

	"auth_path": {
		`The path to mount to. Cannot be delimited. Example: "user"`,
		"",
	},

	"auth_type": {
		`The type of the backend. Example: "userpass"`,
		"",
	},

	"auth_desc": {
		`User-friendly description for this credential backend.`,
		"",
	},

	"auth_config": {
		`Configuration for this mount, such as plugin_name.`,
	},

	"auth_plugin": {
		`Name of the auth plugin to use based from the name in the plugin catalog.`,
		"",
	},

	"auth_options": {
		`The options to pass into the backend. Should be a json object with string keys and values.`,
	},

	"policy-list": {
		`List the configured access control policies.`,
		`
This path responds to the following HTTP methods.

    GET /
        List the names of the configured access control policies.

    GET /<name>
        Retrieve the rules for the named policy.

    PUT /<name>
        Add or update a policy.

    DELETE /<name>
        Delete the policy with the given name.
		`,
	},

	"policy": {
		`Read, Modify, or Delete an access control policy.`,
		`
Read the rules of an existing policy, create or update the rules of a policy,
or delete a policy.
		`,
	},

	"policy-name": {
		`The name of the policy. Example: "ops"`,
		"",
	},

	"policy-rules": {
		`The rules of the policy.`,
		"",
	},

	"policy-paths": {
		`The paths on which the policy should be applied.`,
		"",
	},

	"policy-enforcement-level": {
		`The enforcement level to apply to the policy.`,
		"",
	},

	"password-policy-name": {
		`The name of the password policy.`,
		"",
	},

	"audit-hash": {
		"The hash of the given string via the given audit backend",
		"",
	},

	"audit-table": {
		"List the currently enabled audit backends.",
		`
This path responds to the following HTTP methods.

    GET /
        List the currently enabled audit backends.

    PUT /<path>
        Enable an audit backend at the given path.

    DELETE /<path>
        Disable the given audit backend.
		`,
	},

	"audit_path": {
		`The name of the backend. Cannot be delimited. Example: "mysql"`,
		"",
	},

	"audit_type": {
		`The type of the backend. Example: "mysql"`,
		"",
	},

	"audit_desc": {
		`User-friendly description for this audit backend.`,
		"",
	},

	"audit_opts": {
		`Configuration options for the audit backend.`,
		"",
	},

	"audit": {
		`Enable or disable audit backends.`,
		`
Enable a new audit backend or disable an existing backend.
		`,
	},

	"key-status": {
		"Provides information about the backend encryption key.",
		`
		Provides the current backend encryption key term and installation time.
		`,
	},

	"rotate-config": {
		"Configures settings related to the backend encryption key management.",
		`
		Configures settings related to the automatic rotation of the backend encryption key.
		`,
	},

	"rotation-enabled": {
		"Whether automatic rotation is enabled.",
		"",
	},
	"rotation-max-operations": {
		"The number of encryption operations performed before the barrier key is automatically rotated.",
		"",
	},
	"rotation-interval": {
		"How long after installation of an active key term that the key will be automatically rotated.",
		"",
	},
	"rotate": {
		"Rotates the backend encryption key used to persist data.",
		`
		Rotate generates a new encryption key which is used to encrypt all
		data going to the storage backend. The old encryption keys are kept so
		that data encrypted using those keys can still be decrypted.
		`,
	},

	"rekey_backup": {
		"Allows fetching or deleting the backup of the rotated unseal keys.",
		"",
	},

	"capabilities": {
		"Fetches the capabilities of the given token on the given path.",
		`Returns the capabilities of the given token on the path.
		The path will be searched for a path match in all the policies associated with the token.`,
	},

	"capabilities_self": {
		"Fetches the capabilities of the given token on the given path.",
		`Returns the capabilities of the client token on the path.
		The path will be searched for a path match in all the policies associated with the client token.`,
	},

	"capabilities_accessor": {
		"Fetches the capabilities of the token associated with the given token, on the given path.",
		`When there is no access to the token, token accessor can be used to fetch the token's capabilities
		on a given path.`,
	},

	"tidy_leases": {
		`This endpoint performs cleanup tasks that can be run if certain error
conditions have occurred.`,
		`This endpoint performs cleanup tasks that can be run to clean up the
lease entries after certain error conditions. Usually running this is not
necessary, and is only required if upgrade notes or support personnel suggest
it.`,
	},

	"wrap": {
		"Response-wraps an arbitrary JSON object.",
		`Round trips the given input data into a response-wrapped token.`,
	},

	"wrappubkey": {
		"Returns pubkeys used in some wrapping formats.",
		"Returns pubkeys used in some wrapping formats.",
	},

	"unwrap": {
		"Unwraps a response-wrapped token.",
		`Unwraps a response-wrapped token. Unlike simply reading from cubbyhole/response,
		this provides additional validation on the token, and rather than a JSON-escaped
		string, the returned response is the exact same as the contained wrapped response.`,
	},

	"wraplookup": {
		"Looks up the properties of a response-wrapped token.",
		`Returns the creation TTL and creation time of a response-wrapped token.`,
	},

	"rewrap": {
		"Rotates a response-wrapped token.",
		`Rotates a response-wrapped token; the output is a new token with the same
		response wrapped inside and the same creation TTL. The original token is revoked.`,
	},
	"audited-headers-name": {
		"Configures the headers sent to the audit logs.",
		`
This path responds to the following HTTP methods.

	GET /<name>
		Returns the setting for the header with the given name.

	POST /<name>
		Enable auditing of the given header.

	DELETE /<path>
		Disable auditing of the given header.
		`,
	},
	"audited-headers": {
		"Lists the headers configured to be audited.",
		`Returns a list of headers that have been configured to be audited.`,
	},
	"plugin-catalog-list-all": {
		"Lists all the plugins known to Vault",
		`
This path responds to the following HTTP methods.
		LIST /
			Returns a list of names of configured plugins.
		`,
	},
	"plugin-catalog": {
		"Configures the plugins known to Vault",
		`
This path responds to the following HTTP methods.
		LIST /
			Returns a list of names of configured plugins.

		GET /<name>
			Retrieve the metadata for the named plugin.

		PUT /<name>
			Add or update plugin.

		DELETE /<name>
			Delete the plugin with the given name.
		`,
	},
	"plugin-catalog_name": {
		"The name of the plugin",
		"",
	},
	"plugin-catalog_type": {
		"The type of the plugin, may be auth, secret, or database",
		"",
	},
	"plugin-catalog_sha-256": {
		`The SHA256 sum of the executable used in the
command field. This should be HEX encoded.`,
		"",
	},
	"plugin-catalog_command": {
		`The command used to start the plugin. The
executable defined in this command must exist in vault's
plugin directory.`,
		"",
	},
	"plugin-catalog_args": {
		`The args passed to plugin command.`,
		"",
	},
	"plugin-catalog_env": {
		`The environment variables passed to plugin command.
Each entry is of the form "key=value".`,
		"",
	},
	"leases": {
		`View or list lease metadata.`,
		`
This path responds to the following HTTP methods.

    PUT /
        Retrieve the metadata for the provided lease id.

    LIST /<prefix>
        Lists the leases for the named prefix.
		`,
	},

	"leases-list-prefix": {
		`The path to list leases under. Example: "aws/creds/deploy"`,
		"",
	},
	"plugin-reload": {
		"Reload mounts that use a particular backend plugin.",
		`Reload mounts that use a particular backend plugin. Either the plugin name
		or the desired plugin backend mounts must be provided, but not both. In the
		case that the plugin name is provided, all mounted paths that use that plugin
		backend will be reloaded.`,
	},
	"plugin-backend-reload-plugin": {
		`The name of the plugin to reload, as registered in the plugin catalog.`,
		"",
	},
	"plugin-backend-reload-mounts": {
		`The mount paths of the plugin backends to reload.`,
		"",
	},
	"hash": {
		"Generate a hash sum for input data",
		"Generates a hash sum of the given algorithm against the given input data.",
	},
	"random": {
		"Generate random bytes",
		"This function can be used to generate high-entropy random bytes.",
	},
	"listing_visibility": {
		"Determines the visibility of the mount in the UI-specific listing endpoint. Accepted value are 'unauth' and ''.",
		"",
	},
	"passthrough_request_headers": {
		"A list of headers to whitelist and pass from the request to the plugin.",
		"",
	},
	"allowed_response_headers": {
		"A list of headers to whitelist and allow a plugin to set on responses.",
		"",
	},
	"token_type": {
		"The type of token to issue (service or batch).",
		"",
	},
	"raw": {
		"Write, Read, and Delete data directly in the Storage backend.",
		"",
	},
	"internal-ui-feature-flags": {
		"Enabled feature flags. Internal API; its location, inputs, and outputs may change.",
		"",
	},
	"internal-ui-mounts": {
		"Information about mounts returned according to their tuned visibility. Internal API; its location, inputs, and outputs may change.",
		"",
	},
	"internal-ui-namespaces": {
		"Information about visible child namespaces. Internal API; its location, inputs, and outputs may change.",
		`Information about visible child namespaces returned starting from the request's
		context namespace and filtered based on access from the client token. Internal API;
		its location, inputs, and outputs may change.`,
	},
	"internal-ui-resultant-acl": {
		"Information about a token's resultant ACL. Internal API; its location, inputs, and outputs may change.",
		"",
	},
	"metrics": {
		"Export the metrics aggregated for telemetry purpose.",
		"",
	},
	"internal-counters-requests": {
		"Count of requests seen by this Vault cluster over time.",
		"Count of requests seen by this Vault cluster over time. Not included in count: health checks, UI asset requests, requests forwarded from another cluster.",
	},
	"internal-counters-tokens": {
		"Count of active tokens in this Vault cluster.",
		"Count of active tokens in this Vault cluster.",
	},
	"internal-counters-entities": {
		"Count of active entities in this Vault cluster.",
		"Count of active entities in this Vault cluster.",
	},
	"host-info": {
		"Information about the host instance that this Vault server is running on.",
		`Information about the host instance that this Vault server is running on.
		The information that gets collected includes host hardware information, and CPU,
		disk, and memory utilization`,
	},
	"activity-query": {
		"Query the historical count of clients.",
		"Query the historical count of clients.",
	},
	"activity-monthly": {
		"Count of active clients so far this month.",
		"Count of active clients so far this month.",
	},
	"activity-config": {
		"Control the collection and reporting of client counts.",
		"Control the collection and reporting of client counts.",
	},
	"count-leases": {
		"Count of leases associated with this Vault cluster",
		"Count of leases associated with this Vault cluster",
	},
	"list-leases": {
		"List leases associated with this Vault cluster",
		"Requires sudo capability. List leases associated with this Vault cluster",
	},
}
