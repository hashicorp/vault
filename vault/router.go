package vault

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	metrics "github.com/armon/go-metrics"
	radix "github.com/armon/go-radix"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
)

var (
	deniedPassthroughRequestHeaders = []string{
		consts.AuthHeaderName,
	}
)

// Router is used to do prefix based routing of a request to a logical backend
type Router struct {
	l                  sync.RWMutex
	root               *radix.Tree
	mountUUIDCache     *radix.Tree
	mountAccessorCache *radix.Tree
	tokenStoreSaltFunc func(context.Context) (*salt.Salt, error)
	// storagePrefix maps the prefix used for storage (ala the BarrierView)
	// to the backend. This is used to map a key back into the backend that owns it.
	// For example, logical/uuid1/foobar -> secrets/ (kv backend) + foobar
	storagePrefix *radix.Tree
}

// NewRouter returns a new router
func NewRouter() *Router {
	r := &Router{
		root:               radix.New(),
		storagePrefix:      radix.New(),
		mountUUIDCache:     radix.New(),
		mountAccessorCache: radix.New(),
	}
	return r
}

// routeEntry is used to represent a mount point in the router
type routeEntry struct {
	tainted       bool
	backend       logical.Backend
	mountEntry    *MountEntry
	storageView   logical.Storage
	storagePrefix string
	rootPaths     atomic.Value
	loginPaths    atomic.Value
	l             sync.RWMutex
}

type validateMountResponse struct {
	MountType     string `json:"mount_type" structs:"mount_type" mapstructure:"mount_type"`
	MountAccessor string `json:"mount_accessor" structs:"mount_accessor" mapstructure:"mount_accessor"`
	MountPath     string `json:"mount_path" structs:"mount_path" mapstructure:"mount_path"`
	MountLocal    bool   `json:"mount_local" structs:"mount_local" mapstructure:"mount_local"`
}

// validateMountByAccessor returns the mount type and ID for a given mount
// accessor
func (r *Router) validateMountByAccessor(accessor string) *validateMountResponse {
	if accessor == "" {
		return nil
	}

	mountEntry := r.MatchingMountByAccessor(accessor)
	if mountEntry == nil {
		return nil
	}

	mountPath := mountEntry.Path
	if mountEntry.Table == credentialTableType {
		mountPath = credentialRoutePrefix + mountPath
	}

	return &validateMountResponse{
		MountAccessor: mountEntry.Accessor,
		MountType:     mountEntry.Type,
		MountPath:     mountPath,
		MountLocal:    mountEntry.Local,
	}
}

// SaltID is used to apply a salt and hash to an ID to make sure its not reversible
func (re *routeEntry) SaltID(id string) string {
	return salt.SaltID(re.mountEntry.UUID, id, salt.SHA1Hash)
}

// Mount is used to expose a logical backend at a given prefix, using a unique salt,
// and the barrier view for that path.
func (r *Router) Mount(backend logical.Backend, prefix string, mountEntry *MountEntry, storageView *BarrierView) error {
	r.l.Lock()
	defer r.l.Unlock()

	// prepend namespace
	prefix = mountEntry.Namespace().Path + prefix

	// Check if this is a nested mount
	if existing, _, ok := r.root.LongestPrefix(prefix); ok && existing != "" {
		return fmt.Errorf("cannot mount under existing mount %q", existing)
	}

	// Build the paths
	paths := new(logical.Paths)
	if backend != nil {
		specialPaths := backend.SpecialPaths()
		if specialPaths != nil {
			paths = specialPaths
		}
	}

	// Create a mount entry
	re := &routeEntry{
		tainted:       false,
		backend:       backend,
		mountEntry:    mountEntry,
		storagePrefix: storageView.Prefix(),
		storageView:   storageView,
	}
	re.rootPaths.Store(pathsToRadix(paths.Root))
	re.loginPaths.Store(pathsToRadix(paths.Unauthenticated))

	switch {
	case prefix == "":
		return fmt.Errorf("missing prefix to be used for router entry; mount_path: %q, mount_type: %q", re.mountEntry.Path, re.mountEntry.Type)
	case re.storagePrefix == "":
		return fmt.Errorf("missing storage view prefix; mount_path: %q, mount_type: %q", re.mountEntry.Path, re.mountEntry.Type)
	case re.mountEntry.UUID == "":
		return fmt.Errorf("missing mount identifier; mount_path: %q, mount_type: %q", re.mountEntry.Path, re.mountEntry.Type)
	case re.mountEntry.Accessor == "":
		return fmt.Errorf("missing mount accessor; mount_path: %q, mount_type: %q", re.mountEntry.Path, re.mountEntry.Type)
	}

	r.root.Insert(prefix, re)
	r.storagePrefix.Insert(re.storagePrefix, re)
	r.mountUUIDCache.Insert(re.mountEntry.UUID, re.mountEntry)
	r.mountAccessorCache.Insert(re.mountEntry.Accessor, re.mountEntry)

	return nil
}

// Unmount is used to remove a logical backend from a given prefix
func (r *Router) Unmount(ctx context.Context, prefix string) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	prefix = ns.Path + prefix

	r.l.Lock()
	defer r.l.Unlock()

	// Fast-path out if the backend doesn't exist
	raw, ok := r.root.Get(prefix)
	if !ok {
		return nil
	}

	// Call backend's Cleanup routine
	re := raw.(*routeEntry)
	if re.backend != nil {
		re.backend.Cleanup(ctx)
	}

	// Purge from the radix trees
	r.root.Delete(prefix)
	r.storagePrefix.Delete(re.storagePrefix)
	r.mountUUIDCache.Delete(re.mountEntry.UUID)
	r.mountAccessorCache.Delete(re.mountEntry.Accessor)

	return nil
}

// Remount is used to change the mount location of a logical backend
func (r *Router) Remount(ctx context.Context, src, dst string) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	src = ns.Path + src
	dst = ns.Path + dst

	r.l.Lock()
	defer r.l.Unlock()

	// Check for existing mount
	raw, ok := r.root.Get(src)
	if !ok {
		return fmt.Errorf("no mount at %q", src)
	}

	// Update the mount point
	r.root.Delete(src)
	r.root.Insert(dst, raw)
	return nil
}

// Taint is used to mark a path as tainted. This means only RollbackOperation
// RevokeOperation requests are allowed to proceed
func (r *Router) Taint(ctx context.Context, path string) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	path = ns.Path + path

	r.l.Lock()
	defer r.l.Unlock()
	_, raw, ok := r.root.LongestPrefix(path)
	if ok {
		raw.(*routeEntry).tainted = true
	}
	return nil
}

// Untaint is used to unmark a path as tainted.
func (r *Router) Untaint(ctx context.Context, path string) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	path = ns.Path + path

	r.l.Lock()
	defer r.l.Unlock()
	_, raw, ok := r.root.LongestPrefix(path)
	if ok {
		raw.(*routeEntry).tainted = false
	}
	return nil
}

func (r *Router) MatchingMountByUUID(mountID string) *MountEntry {
	if mountID == "" {
		return nil
	}

	r.l.RLock()

	_, raw, ok := r.mountUUIDCache.LongestPrefix(mountID)
	if !ok {
		r.l.RUnlock()
		return nil
	}

	r.l.RUnlock()
	return raw.(*MountEntry)
}

// MatchingMountByAccessor returns the MountEntry by accessor lookup
func (r *Router) MatchingMountByAccessor(mountAccessor string) *MountEntry {
	if mountAccessor == "" {
		return nil
	}

	r.l.RLock()

	_, raw, ok := r.mountAccessorCache.LongestPrefix(mountAccessor)
	if !ok {
		r.l.RUnlock()
		return nil
	}

	r.l.RUnlock()
	return raw.(*MountEntry)
}

// MatchingMount returns the mount prefix that would be used for a path
func (r *Router) MatchingMount(ctx context.Context, path string) string {
	r.l.RLock()
	mount := r.matchingMountInternal(ctx, path)
	r.l.RUnlock()
	return mount
}

func (r *Router) matchingMountInternal(ctx context.Context, path string) string {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return ""
	}
	path = ns.Path + path

	mount, _, ok := r.root.LongestPrefix(path)
	if !ok {
		return ""
	}
	return mount
}

// matchingPrefixInternal returns a mount prefix that a path may be a part of
func (r *Router) matchingPrefixInternal(ctx context.Context, path string) string {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return ""
	}
	path = ns.Path + path

	var existing string
	fn := func(existingPath string, v interface{}) bool {
		if strings.HasPrefix(existingPath, path) {
			existing = existingPath
			return true
		}
		return false
	}
	r.root.WalkPrefix(path, fn)
	return existing
}

// MountConflict determines if there are potential path conflicts
func (r *Router) MountConflict(ctx context.Context, path string) string {
	r.l.RLock()
	defer r.l.RUnlock()
	if exactMatch := r.matchingMountInternal(ctx, path); exactMatch != "" {
		return exactMatch
	}
	if prefixMatch := r.matchingPrefixInternal(ctx, path); prefixMatch != "" {
		return prefixMatch
	}
	return ""
}

// MatchingStorageByAPIPath/StoragePath returns the storage used for
// API/Storage paths respectively
func (r *Router) MatchingStorageByAPIPath(ctx context.Context, path string) logical.Storage {
	return r.matchingStorage(ctx, path, true)
}
func (r *Router) MatchingStorageByStoragePath(ctx context.Context, path string) logical.Storage {
	return r.matchingStorage(ctx, path, false)
}
func (r *Router) matchingStorage(ctx context.Context, path string, apiPath bool) logical.Storage {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil
	}
	path = ns.Path + path

	var raw interface{}
	var ok bool
	r.l.RLock()
	if apiPath {
		_, raw, ok = r.root.LongestPrefix(path)
	} else {
		_, raw, ok = r.storagePrefix.LongestPrefix(path)
	}
	r.l.RUnlock()
	if !ok {
		return nil
	}
	return raw.(*routeEntry).storageView
}

// MatchingMountEntry returns the MountEntry used for a path
func (r *Router) MatchingMountEntry(ctx context.Context, path string) *MountEntry {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil
	}
	path = ns.Path + path

	r.l.RLock()
	_, raw, ok := r.root.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return nil
	}
	return raw.(*routeEntry).mountEntry
}

// MatchingBackend returns the backend used for a path
func (r *Router) MatchingBackend(ctx context.Context, path string) logical.Backend {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil
	}
	path = ns.Path + path

	r.l.RLock()
	_, raw, ok := r.root.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return nil
	}
	return raw.(*routeEntry).backend
}

// MatchingSystemView returns the SystemView used for a path
func (r *Router) MatchingSystemView(ctx context.Context, path string) logical.SystemView {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil
	}
	path = ns.Path + path

	r.l.RLock()
	_, raw, ok := r.root.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return nil
	}
	return raw.(*routeEntry).backend.System()
}

// MatchingStoragePrefixByAPIPath the storage prefix for the given api path
func (r *Router) MatchingStoragePrefixByAPIPath(ctx context.Context, path string) (string, bool) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return "", false
	}
	path = ns.Path + path

	_, prefix, found := r.matchingMountEntryByPath(ctx, path, true)
	return prefix, found
}

// MatchingAPIPrefixByStoragePath the api path information for the given storage path
func (r *Router) MatchingAPIPrefixByStoragePath(ctx context.Context, path string) (*namespace.Namespace, string, string, bool) {
	me, prefix, found := r.matchingMountEntryByPath(ctx, path, false)
	if !found {
		return nil, "", "", found
	}

	mountPath := me.Path
	// Add back the prefix for credential backends
	if strings.HasPrefix(path, credentialBarrierPrefix) {
		mountPath = credentialRoutePrefix + mountPath
	}

	return me.Namespace(), mountPath, prefix, found
}

func (r *Router) matchingMountEntryByPath(ctx context.Context, path string, apiPath bool) (*MountEntry, string, bool) {
	var raw interface{}
	var ok bool
	r.l.RLock()
	if apiPath {
		_, raw, ok = r.root.LongestPrefix(path)
	} else {
		_, raw, ok = r.storagePrefix.LongestPrefix(path)
	}
	r.l.RUnlock()
	if !ok {
		return nil, "", false
	}

	// Extract the mount path and storage prefix
	re := raw.(*routeEntry)
	prefix := re.storagePrefix

	return re.mountEntry, prefix, true
}

// Route is used to route a given request
func (r *Router) Route(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	resp, _, _, err := r.routeCommon(ctx, req, false)
	return resp, err
}

// RouteExistenceCheck is used to route a given existence check request
func (r *Router) RouteExistenceCheck(ctx context.Context, req *logical.Request) (*logical.Response, bool, bool, error) {
	resp, ok, exists, err := r.routeCommon(ctx, req, true)
	return resp, ok, exists, err
}

func (r *Router) routeCommon(ctx context.Context, req *logical.Request, existenceCheck bool) (*logical.Response, bool, bool, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, false, false, err
	}

	// Find the mount point
	r.l.RLock()
	adjustedPath := req.Path
	mount, raw, ok := r.root.LongestPrefix(ns.Path + adjustedPath)
	if !ok && !strings.HasSuffix(adjustedPath, "/") {
		// Re-check for a backend by appending a slash. This lets "foo" mean
		// "foo/" at the root level which is almost always what we want.
		adjustedPath += "/"
		mount, raw, ok = r.root.LongestPrefix(ns.Path + adjustedPath)
	}
	r.l.RUnlock()
	if !ok {
		return logical.ErrorResponse(fmt.Sprintf("no handler for route '%s'", req.Path)), false, false, logical.ErrUnsupportedPath
	}
	req.Path = adjustedPath
	defer metrics.MeasureSince([]string{"route", string(req.Operation),
		strings.Replace(mount, "/", "-", -1)}, time.Now())
	re := raw.(*routeEntry)

	// Grab a read lock on the route entry, this protects against the backend
	// being reloaded during a request.
	re.l.RLock()
	defer re.l.RUnlock()

	// Filtered mounts will have a nil backend
	if re.backend == nil {
		return logical.ErrorResponse(fmt.Sprintf("no handler for route '%s'", req.Path)), false, false, logical.ErrUnsupportedPath
	}

	// If the path is tainted, we reject any operation except for
	// Rollback and Revoke
	if re.tainted {
		switch req.Operation {
		case logical.RevokeOperation, logical.RollbackOperation:
		default:
			return logical.ErrorResponse(fmt.Sprintf("no handler for route '%s'", req.Path)), false, false, logical.ErrUnsupportedPath
		}
	}

	// Adjust the path to exclude the routing prefix
	originalPath := req.Path
	req.Path = strings.TrimPrefix(ns.Path+req.Path, mount)
	req.MountPoint = mount
	req.MountType = re.mountEntry.Type
	if req.Path == "/" {
		req.Path = ""
	}

	// Attach the storage view for the request
	req.Storage = re.storageView

	originalEntityID := req.EntityID

	// Hash the request token unless the request is being routed to the token
	// or system backend.
	clientToken := req.ClientToken
	switch {
	case strings.HasPrefix(originalPath, "auth/token/"):
	case strings.HasPrefix(originalPath, "sys/"):
	case strings.HasPrefix(originalPath, cubbyholeMountPath):
		if req.Operation == logical.RollbackOperation {
			// Backend doesn't support this and it can't properly look up a
			// cubbyhole ID so just return here
			return nil, false, false, nil
		}

		te := req.TokenEntry()

		if te == nil {
			return nil, false, false, fmt.Errorf("nil token entry")
		}

		if te.Type != logical.TokenTypeService {
			return logical.ErrorResponse(`cubbyhole operations are only supported by "service" type tokens`), false, false, nil
		}

		switch {
		case te.NamespaceID == namespace.RootNamespaceID && !strings.HasPrefix(req.ClientToken, "s."):
			// In order for the token store to revoke later, we need to have the same
			// salted ID, so we double-salt what's going to the cubbyhole backend
			salt, err := r.tokenStoreSaltFunc(ctx)
			if err != nil {
				return nil, false, false, err
			}
			req.ClientToken = re.SaltID(salt.SaltID(req.ClientToken))

		default:
			if te.CubbyholeID == "" {
				return nil, false, false, fmt.Errorf("empty cubbyhole id")
			}
			req.ClientToken = te.CubbyholeID
		}

	default:
		req.ClientToken = re.SaltID(req.ClientToken)
	}

	// Cache the pointer to the original connection object
	originalConn := req.Connection

	// Cache the identifier of the request
	originalReqID := req.ID

	// Cache the client token's number of uses in the request
	originalClientTokenRemainingUses := req.ClientTokenRemainingUses
	req.ClientTokenRemainingUses = 0

	originalMFACreds := req.MFACreds
	req.MFACreds = nil

	originalControlGroup := req.ControlGroup
	req.ControlGroup = nil

	// Cache the headers
	headers := req.Headers
	req.Headers = nil

	// Filter and add passthrough headers to the backend
	var passthroughRequestHeaders []string
	if rawVal, ok := re.mountEntry.synthesizedConfigCache.Load("passthrough_request_headers"); ok {
		passthroughRequestHeaders = rawVal.([]string)
	}
	var allowedResponseHeaders []string
	if rawVal, ok := re.mountEntry.synthesizedConfigCache.Load("allowed_response_headers"); ok {
		allowedResponseHeaders = rawVal.([]string)
	}

	if len(passthroughRequestHeaders) > 0 {
		req.Headers = filteredHeaders(headers, passthroughRequestHeaders, deniedPassthroughRequestHeaders)
	}

	// Cache the wrap info of the request
	var wrapInfo *logical.RequestWrapInfo
	if req.WrapInfo != nil {
		wrapInfo = &logical.RequestWrapInfo{
			TTL:      req.WrapInfo.TTL,
			Format:   req.WrapInfo.Format,
			SealWrap: req.WrapInfo.SealWrap,
		}
	}

	originalPolicyOverride := req.PolicyOverride
	reqTokenEntry := req.TokenEntry()
	req.SetTokenEntry(nil)

	// Reset the request before returning
	defer func() {
		req.Path = originalPath
		req.MountPoint = mount
		req.MountType = re.mountEntry.Type
		req.Connection = originalConn
		req.ID = originalReqID
		req.Storage = nil
		req.ClientToken = clientToken
		req.ClientTokenRemainingUses = originalClientTokenRemainingUses
		req.WrapInfo = wrapInfo
		req.Headers = headers
		req.PolicyOverride = originalPolicyOverride
		// This is only set in one place, after routing, so should never be set
		// by a backend
		req.SetLastRemoteWAL(0)

		// This will be used for attaching the mount accessor for the identities
		// returned by the authentication backends
		req.MountAccessor = re.mountEntry.Accessor

		req.EntityID = originalEntityID

		req.MFACreds = originalMFACreds

		req.SetTokenEntry(reqTokenEntry)
		req.ControlGroup = originalControlGroup
	}()

	// Invoke the backend
	if existenceCheck {
		ok, exists, err := re.backend.HandleExistenceCheck(ctx, req)
		return nil, ok, exists, err
	} else {
		resp, err := re.backend.HandleRequest(ctx, req)
		if resp != nil {
			if len(allowedResponseHeaders) > 0 {
				resp.Headers = filteredHeaders(resp.Headers, allowedResponseHeaders, nil)
			} else {
				resp.Headers = nil
			}

			if resp.Auth != nil {
				// When a token gets renewed, the request hits this path and
				// reaches token store. Token store delegates the renewal to the
				// expiration manager. Expiration manager in-turn creates a
				// different logical request and forwards the request to the auth
				// backend that had initially authenticated the login request. The
				// forwarding to auth backend will make this code path hit for the
				// second time for the same renewal request. The accessors in the
				// Alias structs should be of the auth backend and not of the token
				// store. Therefore, avoiding the overwriting of accessors by
				// having a check for path prefix having "renew". This gets applied
				// for "renew" and "renew-self" requests.
				if !strings.HasPrefix(req.Path, "renew") {
					if resp.Auth.Alias != nil {
						resp.Auth.Alias.MountAccessor = re.mountEntry.Accessor
					}
					for _, alias := range resp.Auth.GroupAliases {
						alias.MountAccessor = re.mountEntry.Accessor
					}
				}

				switch re.mountEntry.Type {
				case "token", "ns_token":
					// Nothing; we respect what the token store is telling us and
					// we don't allow tuning
				default:
					switch re.mountEntry.Config.TokenType {
					case logical.TokenTypeService, logical.TokenTypeBatch:
						resp.Auth.TokenType = re.mountEntry.Config.TokenType
					case logical.TokenTypeDefault, logical.TokenTypeDefaultService:
						if resp.Auth.TokenType == logical.TokenTypeDefault {
							resp.Auth.TokenType = logical.TokenTypeService
						}
					case logical.TokenTypeDefaultBatch:
						if resp.Auth.TokenType == logical.TokenTypeDefault {
							resp.Auth.TokenType = logical.TokenTypeBatch
						}
					}
				}
			}
		}

		return resp, false, false, err
	}
}

// RootPath checks if the given path requires root privileges
func (r *Router) RootPath(ctx context.Context, path string) bool {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return false
	}

	adjustedPath := ns.Path + path

	r.l.RLock()
	mount, raw, ok := r.root.LongestPrefix(adjustedPath)
	r.l.RUnlock()
	if !ok {
		return false
	}
	re := raw.(*routeEntry)

	// Trim to get remaining path
	remain := strings.TrimPrefix(adjustedPath, mount)

	// Check the rootPaths of this backend
	rootPaths := re.rootPaths.Load().(*radix.Tree)
	match, raw, ok := rootPaths.LongestPrefix(remain)
	if !ok {
		return false
	}
	prefixMatch := raw.(bool)

	// Handle the prefix match case
	if prefixMatch {
		return strings.HasPrefix(remain, match)
	}

	// Handle the exact match case
	return match == remain
}

// LoginPath checks if the given path is used for logins
func (r *Router) LoginPath(ctx context.Context, path string) bool {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return false
	}

	adjustedPath := ns.Path + path

	r.l.RLock()
	mount, raw, ok := r.root.LongestPrefix(adjustedPath)
	r.l.RUnlock()
	if !ok {
		return false
	}
	re := raw.(*routeEntry)

	// Trim to get remaining path
	remain := strings.TrimPrefix(adjustedPath, mount)

	// Check the loginPaths of this backend
	loginPaths := re.loginPaths.Load().(*radix.Tree)
	match, raw, ok := loginPaths.LongestPrefix(remain)
	if !ok {
		return false
	}
	prefixMatch := raw.(bool)

	// Handle the prefix match case
	if prefixMatch {
		return strings.HasPrefix(remain, match)
	}

	// Handle the exact match case
	return match == remain
}

// pathsToRadix converts a list of special paths to a radix tree.
func pathsToRadix(paths []string) *radix.Tree {
	tree := radix.New()
	for _, path := range paths {
		// Check if this is a prefix or exact match
		prefixMatch := len(path) >= 1 && path[len(path)-1] == '*'
		if prefixMatch {
			path = path[:len(path)-1]
		}

		tree.Insert(path, prefixMatch)
	}

	return tree
}

// filteredHeaders returns a headers map[string][]string that
// contains the filtered values contained in candidateHeaders. Filtering of
// candidateHeaders from the origHeaders is done is a case-insensitive manner.
// Headers that match values from deniedHeaders will be ignored.
func filteredHeaders(origHeaders map[string][]string, candidateHeaders, deniedHeaders []string) map[string][]string {
	// Short-circuit if there's nothing to filter
	if len(candidateHeaders) == 0 {
		return nil
	}

	retHeaders := make(map[string][]string, len(origHeaders))

	// Filter candidateHeaders values through deniedHeaders first. Returns the
	// lowercased complement set. We call even if no denied headers to get the
	// values lowercased.
	allowedCandidateHeaders := strutil.Difference(candidateHeaders, deniedHeaders, true)

	// Create a map that uses lowercased header values as the key and the original
	// header naming as the value for comparison down below.
	lowerOrigHeaderKeys := make(map[string]string, len(origHeaders))
	for key := range origHeaders {
		lowerOrigHeaderKeys[strings.ToLower(key)] = key
	}

	// Case-insensitive compare of passthrough headers against originating
	// headers. The returned headers will be the same casing as the originating
	// header name.
	for _, ch := range allowedCandidateHeaders {
		if header, ok := lowerOrigHeaderKeys[ch]; ok {
			retHeaders[header] = origHeaders[header]
		}
	}

	return retHeaders
}
