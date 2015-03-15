package vault

import (
	"strings"

	"github.com/hashicorp/vault/logical"
)

// BarrierView wraps a SecurityBarrier and ensures all access is automatically
// prefixed. This is used to prevent anyone with access to the view to access
// any data in the durable storage outside of their prefix. Conceptually this
// is like a "chroot" into the barrier.
//
// BarrierView implements logical.Storage so it can be passed in as the
// durable storage mechanism for logical views.
type BarrierView struct {
	barrier SecurityBarrier
	prefix  string
}

// NewBarrierView takes an underlying security barrier and returns
// a view of it that can only operate with the given prefix.
func NewBarrierView(barrier SecurityBarrier, prefix string) *BarrierView {
	return &BarrierView{
		barrier: barrier,
		prefix:  prefix,
	}
}

// logical.Storage impl.
func (v *BarrierView) List(prefix string) ([]string, error) {
	return v.barrier.List(v.expandKey(prefix))
}

// logical.Storage impl.
func (v *BarrierView) Get(key string) (*logical.StorageEntry, error) {
	entry, err := v.barrier.Get(v.expandKey(key))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	if entry != nil {
		entry.Key = v.truncateKey(entry.Key)
	}

	return &logical.StorageEntry{
		Key:   entry.Key,
		Value: entry.Value,
	}, nil
}

// logical.Storage impl.
func (v *BarrierView) Put(entry *logical.StorageEntry) error {
	nested := &Entry{
		Key:   v.expandKey(entry.Key),
		Value: entry.Value,
	}
	return v.barrier.Put(nested)
}

// logical.Storage impl.
func (v *BarrierView) Delete(key string) error {
	return v.barrier.Delete(v.expandKey(key))
}

// SubView constructs a nested sub-view using the given prefix
func (v *BarrierView) SubView(prefix string) *BarrierView {
	sub := v.expandKey(prefix)
	return &BarrierView{barrier: v.barrier, prefix: sub}
}

// expandKey is used to expand to the full key path with the prefix
func (v *BarrierView) expandKey(suffix string) string {
	return v.prefix + suffix
}

// truncateKey is used to remove the prefix of the key
func (v *BarrierView) truncateKey(full string) string {
	return strings.TrimPrefix(full, v.prefix)
}
