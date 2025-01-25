// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package physical

import (
	"context"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/permitpool"
)

const DefaultParallelOperations = 128

// The operation type
type Operation string

const (
	DeleteOperation Operation = "delete"
	GetOperation              = "get"
	ListOperation             = "list"
	PutOperation              = "put"
)

const (
	ErrValueTooLarge = "put failed due to value being too large"
	ErrKeyTooLarge   = "put failed due to key being too large"
)

// Backend is the interface required for a physical
// backend. A physical backend is used to durably store
// data outside of Vault. As such, it is completely untrusted,
// and is only accessed via a security barrier. The backends
// must represent keys in a hierarchical manner. All methods
// are expected to be thread safe.
type Backend interface {
	// Put is used to insert or update an entry
	Put(ctx context.Context, entry *Entry) error

	// Get is used to fetch an entry
	Get(ctx context.Context, key string) (*Entry, error)

	// Delete is used to permanently delete an entry
	Delete(ctx context.Context, key string) error

	// List is used to list all the keys under a given
	// prefix, up to the next prefix.
	List(ctx context.Context, prefix string) ([]string, error)
}

// HABackend is an extensions to the standard physical
// backend to support high-availability. Vault only expects to
// use mutual exclusion to allow multiple instances to act as a
// hot standby for a leader that services all requests.
type HABackend interface {
	// LockWith is used for mutual exclusion based on the given key.
	LockWith(key, value string) (Lock, error)

	// Whether or not HA functionality is enabled
	HAEnabled() bool
}

// RemovableNodeHABackend is used for HA backends that can remove nodes from
// their cluster
type RemovableNodeHABackend interface {
	HABackend

	// IsNodeRemoved checks if the node with the given ID has been removed.
	// This will only be called on the active node.
	IsNodeRemoved(ctx context.Context, nodeID string) (bool, error)

	// NodeID returns the ID for this node
	NodeID() string

	// IsRemoved checks if this node has been removed
	IsRemoved() bool

	// RemoveSelf marks this node as being removed
	RemoveSelf() error
}

// FencingHABackend is an HABackend which provides the additional guarantee that
// each Lock it returns from LockWith is also a FencingLock. A FencingLock
// provides a mechanism to retrieve a fencing token that can be included by
// future writes by the backend to ensure that it is still the current lock
// holder at the time the write commits. Without this timing might allow a lock
// holder not to notice it's no longer the active node for long enough for it to
// write data to storage even while a new active node is writing causing
// corruption. For Consul backend the fencing token is the session id which is
// submitted with `check-session` operation on each write to ensure the write
// only completes if the session is still holding the lock. For raft backend
// this isn't needed because our in-process raft library is unable to write if
// it's not the leader anyway.
//
// If you implement this, Vault will call RegisterActiveNodeLock with the Lock
// instance returned by LockWith after it successfully locks it. This keeps the
// backend oblivious to the specific key we use for active node locks and allows
// potential future usage of locks for other purposes in the future.
//
// Note that all implementations must support writing to storage before
// RegisterActiveNodeLock is called to support initialization of a new cluster.
// They must also skip fencing writes if the write's Context contains a special
// value. This is necessary to allow Vault to clear and re-initialise secondary
// clusters even though there is already an active node with a specific lock
// session since we clear the cluster while Vault is sealed and clearing the
// data might remove the lock in some storages (e.g. Consul). As noted above
// it's not generally safe to allow unfenced writes after a lock so instead we
// special case just a few types of writes that only happen rarely while the
// cluster is sealed. See the IsUnfencedWrite helper function.
type FencingHABackend interface {
	HABackend

	RegisterActiveNodeLock(l Lock) error
}

// unfencedWriteContextKeyType is a special type to identify context values to
// disable fencing. It's a separate type per the best-practice in Context.Value
// docs to avoid collisions even if the key might match.
type unfencedWriteContextKeyType string

const (
	// unfencedWriteContextKey is the context key we pass the option to bypass
	// fencing through to a FencingHABackend. Note that this is not an ideal use
	// of context values and violates the "do not use it for optional arguments"
	// guidance but has been agreed as a pragmatic option for this case rather
	// than needing to specialize every physical.Backend to understand this
	// option.
	unfencedWriteContextKey unfencedWriteContextKeyType = "vault-disable-fencing"
)

// UnfencedWriteCtx adds metadata to a ctx such that any writes performed
// directly on a FencingHABackend using that context will _not_ add a fencing
// token.
func UnfencedWriteCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, unfencedWriteContextKey, true)
}

// IsUnfencedWrite returns whether or not the context passed has the unfenced
// flag value set.
func IsUnfencedWrite(ctx context.Context) bool {
	isUnfenced, ok := ctx.Value(unfencedWriteContextKey).(bool)
	return ok && isUnfenced
}

// ToggleablePurgemonster is an interface for backends that can toggle on or
// off special functionality and/or support purging. This is only used for the
// cache, don't use it for other things.
type ToggleablePurgemonster interface {
	Purge(ctx context.Context)
	SetEnabled(bool)
}

// RedirectDetect is an optional interface that an HABackend
// can implement. If they do, a redirect address can be automatically
// detected.
type RedirectDetect interface {
	// DetectHostAddr is used to detect the host address
	DetectHostAddr() (string, error)
}

// MountTableLimitingBackend is an optional interface a Backend can implement
// that allows it to support different entry size limits for mount-table-related
// paths. It will only be called in Vault Enterprise.
type MountTableLimitingBackend interface {
	// RegisterMountTablePath informs the Backend that the given path represents
	// part of the mount tables or related metadata. This allows the backend to
	// apply different limits for this entry if configured to do so.
	RegisterMountTablePath(path string)
}

type Lock interface {
	// Lock is used to acquire the given lock
	// The stopCh is optional and if closed should interrupt the lock
	// acquisition attempt. The return struct should be closed when
	// leadership is lost.
	Lock(stopCh <-chan struct{}) (<-chan struct{}, error)

	// Unlock is used to release the lock
	Unlock() error

	// Returns the value of the lock and if it is held by _any_ node
	Value() (bool, string, error)
}

// Factory is the factory function to create a physical backend.
type Factory func(config map[string]string, logger log.Logger) (Backend, error)

// PermitPool is used to limit maximum outstanding requests
// Deprecated: use permitpool.Pool from go-secure-stdlib.
type PermitPool struct {
	*permitpool.Pool
}

// NewPermitPool returns a new permit pool with the provided
// number of permits.
// Deprecated: use permitpool.New from go-secure-stdlib.
func NewPermitPool(permits int) *PermitPool {
	return &PermitPool{
		Pool: permitpool.New(permits),
	}
}

// Acquire returns when a permit has been acquired
// Deprecated: use permitpool.Acquire from go-secure-stdlib.
func (c *PermitPool) Acquire() {
	_ = c.Pool.Acquire(context.Background())
}

// Prefixes is a shared helper function returns all parent 'folders' for a
// given vault key.
// e.g. for 'foo/bar/baz', it returns ['foo', 'foo/bar']
func Prefixes(s string) []string {
	components := strings.Split(s, "/")
	result := []string{}
	for i := 1; i < len(components); i++ {
		result = append(result, strings.Join(components[:i], "/"))
	}
	return result
}
