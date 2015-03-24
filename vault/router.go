package vault

import (
	"fmt"
	"strings"
	"sync"

	"github.com/armon/go-radix"
	"github.com/hashicorp/vault/credential"
	"github.com/hashicorp/vault/logical"
)

// Router is used to do prefix based routing of a request to a logical backend
type Router struct {
	l    sync.RWMutex
	root *radix.Tree
}

// NewRouter returns a new router
func NewRouter() *Router {
	r := &Router{
		root: radix.New(),
	}
	return r
}

// mountEntry is used to represent a mount point
type mountEntry struct {
	backend    logical.Backend
	view       *BarrierView
	rootPaths  *radix.Tree
	loginPaths *radix.Tree
}

// Mount is used to expose a logical backend at a given prefix
func (r *Router) Mount(backend logical.Backend, prefix string, view *BarrierView) error {
	r.l.Lock()
	defer r.l.Unlock()

	// Check if this is a nested mount
	if existing, _, ok := r.root.LongestPrefix(prefix); ok && existing != "" {
		return fmt.Errorf("cannot mount under existing mount '%s'", existing)
	}

	// Get the root paths
	rootPaths := pathsToRadix(backend.RootPaths())

	// Check if this is a credential backend, calculate the login paths
	var loginPaths *radix.Tree
	if cred, ok := backend.(credential.Backend); ok {
		loginPaths = pathsToRadix(cred.LoginPaths())
	}

	// Create a mount entry
	me := &mountEntry{
		backend:    backend,
		view:       view,
		rootPaths:  rootPaths,
		loginPaths: loginPaths,
	}
	r.root.Insert(prefix, me)
	return nil
}

// Unmount is used to remove a logical backend from a given prefix
func (r *Router) Unmount(prefix string) error {
	r.l.Lock()
	defer r.l.Unlock()
	r.root.Delete(prefix)
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

// MatchingMOunt returns the mount prefix that would be used for a path
func (r *Router) MatchingMount(path string) string {
	r.l.RLock()
	mount, _, ok := r.root.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return ""
	}
	return mount
}

// Route is used to route a given request
func (r *Router) Route(req *logical.Request) (*logical.Response, error) {
	// Find the mount point
	r.l.RLock()
	mount, raw, ok := r.root.LongestPrefix(req.Path)
	r.l.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no handler for route '%s'", req.Path)
	}
	me := raw.(*mountEntry)

	// Adjust the path to exclude the routing prefix
	original := req.Path
	req.Path = strings.TrimPrefix(req.Path, mount)

	// Attach the storage view for the request
	req.Storage = me.view

	// Clear the request token unless this is the token backend
	clientToken := req.ClientToken
	if !strings.HasPrefix(original, "auth/token/") {
		req.ClientToken = ""
	}

	// Reset the request before returning
	defer func() {
		req.Path = original
		req.Storage = nil
		req.ClientToken = clientToken
	}()

	// Invoke the backend
	return me.backend.HandleRequest(req)
}

// RouteLogin is used to route a given login request
func (r *Router) RouteLogin(req *credential.Request) (*credential.Response, error) {
	// Ensure this is a login path
	if !r.LoginPath(req.Path) {
		return nil, fmt.Errorf("invalid login route '%s'", req.Path)
	}

	// Find the mount point
	r.l.RLock()
	mount, raw, ok := r.root.LongestPrefix(req.Path)
	r.l.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no handler for route '%s'", req.Path)
	}
	me := raw.(*mountEntry)

	// Adjust the path, attach the barrier view
	original := req.Path
	req.Path = strings.TrimPrefix(req.Path, mount)
	req.Storage = me.view

	// Reset the request before returning
	defer func() {
		req.Path = original
		req.Storage = nil
	}()

	// Convert to a credential backend
	cred, ok := me.backend.(credential.Backend)
	if !ok {
		return nil, fmt.Errorf("invalid login route '%s'", req.Path)
	}

	// Invoke the backend
	return cred.HandleLogin(req)
}

// RootPath checks if the given path requires root privileges
func (r *Router) RootPath(path string) bool {
	r.l.RLock()
	mount, raw, ok := r.root.LongestPrefix(path)
	r.l.RUnlock()
	if !ok {
		return false
	}
	me := raw.(*mountEntry)

	// Trim to get remaining path
	remain := strings.TrimPrefix(path, mount)

	// Check the rootPaths of this backend
	match, raw, ok := me.rootPaths.LongestPrefix(remain)
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
	me := raw.(*mountEntry)

	// Trim to get remaining path
	remain := strings.TrimPrefix(path, mount)

	// Check the loginPaths of this backend
	match, raw, ok := me.loginPaths.LongestPrefix(remain)
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

// pathsToRadix converts a list of paths potentially ending with
// a wildcard expansion "*" into a radix tree.
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
