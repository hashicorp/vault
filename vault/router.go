package vault

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/armon/go-radix"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
)

// Router is used to do prefix based routing of a request to a logical backend
type Router struct {
	l              sync.RWMutex
	root           *radix.Tree
	tokenStoreSalt *salt.Salt

	// storagePrefix maps the prefix used for storage (ala the BarrierView)
	// to the backend. This is used to map a key back into the backend that owns it.
	// For example, logical/uuid1/foobar -> secrets/ (generic backend) + foobar
	storagePrefix *radix.Tree
}

// NewRouter returns a new router
func NewRouter() *Router {
	r := &Router{
		root:          radix.New(),
		storagePrefix: radix.New(),
	}
	return r
}

// routeEntry is used to represent a mount point in the router
type routeEntry struct {
	tainted     bool
	backend     logical.Backend
	mountEntry  *MountEntry
	storageView *BarrierView
	rootPaths   *radix.Tree
	loginPaths  *radix.Tree
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

	// Check if this is a nested mount
	if existing, _, ok := r.root.LongestPrefix(prefix); ok && existing != "" {
		return fmt.Errorf("cannot mount under existing mount '%s'", existing)
	}

	// Build the paths
	paths := backend.SpecialPaths()
	if paths == nil {
		paths = new(logical.Paths)
	}

	// Create a mount entry
	re := &routeEntry{
		tainted:     false,
		backend:     backend,
		mountEntry:  mountEntry,
		storageView: storageView,
		rootPaths:   pathsToRadix(paths.Root),
		loginPaths:  pathsToRadix(paths.Unauthenticated),
	}
	r.root.Insert(prefix, re)
	r.storagePrefix.Insert(storageView.prefix, re)

	return nil
}

// Unmount is used to remove a logical backend from a given prefix
func (r *Router) Unmount(prefix string) error {
	r.l.Lock()
	defer r.l.Unlock()

	// Fast-path out if the backend doesn't exist
	raw, ok := r.root.Get(prefix)
	if !ok {
		return nil
	}

	// Call backend's Cleanup routine
	re := raw.(*routeEntry)
	re.backend.Cleanup()

	// Purge from the radix trees
	r.root.Delete(prefix)
	r.storagePrefix.Delete(re.storageView.prefix)
	return nil
}

// Remount is used to change the mount location of a logical backend
func (r *Router) Remount(src, dst string) error {
	r.l.Lock()
	defer r.l.Unlock()

	// Check for existing mount
	raw, ok := r.root.Get(src)
	if !ok {
		return fmt.Errorf("no mount at '%s'", src)
	}

	// Update the mount point
	r.root.Delete(src)
	r.root.Insert(dst, raw)
	return nil
}

// Taint is used to mark a path as tainted. This means only RollbackOperation
// RevokeOperation requests are allowed to proceed
func (r *Router) Taint(path string) error {
	r.l.Lock()
	defer r.l.Unlock()
	_, raw, ok := r.root.LongestPrefix(path)
	if ok {
		raw.(*routeEntry).tainted = true
	}
	return nil
}

// Untaint is used to unmark a path as tainted.
func (r *Router) Untaint(path string) error {
	r.l.Lock()
	defer r.l.Unlock()
	_, raw, ok := r.root.LongestPrefix(path)
	if ok {
		raw.(*routeEntry).tainted = false
	}
	return nil
}

// MatchingMount returns the mount prefix that would be used for a path
func (r *Router) MatchingMount(path string) string {
	r.l.RLock()
	mount, _, ok := r.root.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return ""
	}
	return mount
}

// MatchingView returns the view used for a path
func (r *Router) MatchingStorageView(path string) *BarrierView {
	r.l.RLock()
	_, raw, ok := r.root.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return nil
	}
	return raw.(*routeEntry).storageView
}

// MatchingMountEntry returns the MountEntry used for a path
func (r *Router) MatchingMountEntry(path string) *MountEntry {
	r.l.RLock()
	_, raw, ok := r.root.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return nil
	}
	return raw.(*routeEntry).mountEntry
}

// MatchingMountEntry returns the MountEntry used for a path
func (r *Router) MatchingBackend(path string) logical.Backend {
	r.l.RLock()
	_, raw, ok := r.root.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return nil
	}
	return raw.(*routeEntry).backend
}

// MatchingSystemView returns the SystemView used for a path
func (r *Router) MatchingSystemView(path string) logical.SystemView {
	r.l.RLock()
	_, raw, ok := r.root.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return nil
	}
	return raw.(*routeEntry).backend.System()
}

// MatchingStoragePrefix returns the mount path matching and storage prefix
// matching the given path
func (r *Router) MatchingStoragePrefix(path string) (string, string, bool) {
	r.l.RLock()
	_, raw, ok := r.storagePrefix.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return "", "", false
	}

	// Extract the mount path and storage prefix
	re := raw.(*routeEntry)
	mountPath := re.mountEntry.Path
	prefix := re.storageView.prefix
	return mountPath, prefix, true
}

// Route is used to route a given request
func (r *Router) Route(req *logical.Request) (*logical.Response, error) {
	resp, _, _, err := r.routeCommon(req, false)
	return resp, err
}

// Route is used to route a given existence check request
func (r *Router) RouteExistenceCheck(req *logical.Request) (bool, bool, error) {
	_, ok, exists, err := r.routeCommon(req, true)
	return ok, exists, err
}

func (r *Router) routeCommon(req *logical.Request, existenceCheck bool) (*logical.Response, bool, bool, error) {
	// Find the mount point
	r.l.RLock()
	mount, raw, ok := r.root.LongestPrefix(req.Path)
	if !ok {
		// Re-check for a backend by appending a slash. This lets "foo" mean
		// "foo/" at the root level which is almost always what we want.
		req.Path += "/"
		mount, raw, ok = r.root.LongestPrefix(req.Path)
	}
	r.l.RUnlock()
	if !ok {
		return logical.ErrorResponse(fmt.Sprintf("no handler for route '%s'", req.Path)), false, false, logical.ErrUnsupportedPath
	}
	defer metrics.MeasureSince([]string{"route", string(req.Operation),
		strings.Replace(mount, "/", "-", -1)}, time.Now())
	re := raw.(*routeEntry)

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
	req.Path = strings.TrimPrefix(req.Path, mount)
	req.MountPoint = mount
	req.MountType = re.mountEntry.Type
	if req.Path == "/" {
		req.Path = ""
	}

	// Attach the storage view for the request
	req.Storage = re.storageView

	// Hash the request token unless this is the token backend
	clientToken := req.ClientToken
	switch {
	case strings.HasPrefix(originalPath, "auth/token/"):
	case strings.HasPrefix(originalPath, "sys/"):
	case strings.HasPrefix(originalPath, "cubbyhole/"):
		// In order for the token store to revoke later, we need to have the same
		// salted ID, so we double-salt what's going to the cubbyhole backend
		req.ClientToken = re.SaltID(r.tokenStoreSalt.SaltID(req.ClientToken))
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

	// Cache the headers and hide them from backends
	headers := req.Headers
	req.Headers = nil

	// Cache the wrap info of the request
	var wrapInfo *logical.RequestWrapInfo
	if req.WrapInfo != nil {
		wrapInfo = &logical.RequestWrapInfo{
			TTL:    req.WrapInfo.TTL,
			Format: req.WrapInfo.Format,
		}
	}

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
		// This is only set in one place, after routing, so should never be set
		// by a backend
		req.SetLastRemoteWAL(0)
	}()

	// Invoke the backend
	if existenceCheck {
		ok, exists, err := re.backend.HandleExistenceCheck(req)
		return nil, ok, exists, err
	} else {
		resp, err := re.backend.HandleRequest(req)
		return resp, false, false, err
	}
}

// RootPath checks if the given path requires root privileges
func (r *Router) RootPath(path string) bool {
	r.l.RLock()
	mount, raw, ok := r.root.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return false
	}
	re := raw.(*routeEntry)

	// Trim to get remaining path
	remain := strings.TrimPrefix(path, mount)

	// Check the rootPaths of this backend
	match, raw, ok := re.rootPaths.LongestPrefix(remain)
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
func (r *Router) LoginPath(path string) bool {
	r.l.RLock()
	mount, raw, ok := r.root.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return false
	}
	re := raw.(*routeEntry)

	// Trim to get remaining path
	remain := strings.TrimPrefix(path, mount)

	// Check the loginPaths of this backend
	match, raw, ok := re.loginPaths.LongestPrefix(remain)
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

// pathsToRadix converts a the mapping of special paths to a mapping
// of special paths to radix trees.
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
