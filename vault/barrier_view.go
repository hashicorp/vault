package vault

import (
	"errors"
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
	barrier  BarrierStorage
	prefix   string
	readonly bool
}

var (
	ErrRelativePath = errors.New("relative paths not supported")
)

// NewBarrierView takes an underlying security barrier and returns
// a view of it that can only operate with the given prefix.
func NewBarrierView(barrier BarrierStorage, prefix string) *BarrierView {
	return &BarrierView{
		barrier: barrier,
		prefix:  prefix,
	}
}

// sanityCheck is used to perform a sanity check on a key
func (v *BarrierView) sanityCheck(key string) error {
	if strings.Contains(key, "..") {
		return ErrRelativePath
	}
	return nil
}

// logical.Storage impl.
func (v *BarrierView) List(prefix string) ([]string, error) {
	if err := v.sanityCheck(prefix); err != nil {
		return nil, err
	}
	return v.barrier.List(v.expandKey(prefix))
}

// logical.Storage impl.
func (v *BarrierView) Get(key string) (*logical.StorageEntry, error) {
	if err := v.sanityCheck(key); err != nil {
		return nil, err
	}
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
		Key:      entry.Key,
		Value:    entry.Value,
		SealWrap: entry.SealWrap,
	}, nil
}

// logical.Storage impl.
func (v *BarrierView) Put(entry *logical.StorageEntry) error {
	if err := v.sanityCheck(entry.Key); err != nil {
		return err
	}

	expandedKey := v.expandKey(entry.Key)

	if v.readonly {
		return logical.ErrReadOnly
	}

	nested := &Entry{
		Key:      expandedKey,
		Value:    entry.Value,
		SealWrap: entry.SealWrap,
	}
	return v.barrier.Put(nested)
}

// logical.Storage impl.
func (v *BarrierView) Delete(key string) error {
	if err := v.sanityCheck(key); err != nil {
		return err
	}

	expandedKey := v.expandKey(key)

	if v.readonly {
		return logical.ErrReadOnly
	}

	return v.barrier.Delete(expandedKey)
}

// SubView constructs a nested sub-view using the given prefix
func (v *BarrierView) SubView(prefix string) *BarrierView {
	sub := v.expandKey(prefix)
	return &BarrierView{barrier: v.barrier, prefix: sub, readonly: v.readonly}
}

// expandKey is used to expand to the full key path with the prefix
func (v *BarrierView) expandKey(suffix string) string {
	return v.prefix + suffix
}

// truncateKey is used to remove the prefix of the key
func (v *BarrierView) truncateKey(full string) string {
	return strings.TrimPrefix(full, v.prefix)
}
