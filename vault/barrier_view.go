package vault

import (
	"fmt"
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
	barrier BarrierStorage
	prefix  string
}

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
		return fmt.Errorf("key cannot be relative path")
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
		Key:   entry.Key,
		Value: entry.Value,
	}, nil
}

// logical.Storage impl.
func (v *BarrierView) Put(entry *logical.StorageEntry) error {
	if err := v.sanityCheck(entry.Key); err != nil {
		return err
	}
	nested := &Entry{
		Key:   v.expandKey(entry.Key),
		Value: entry.Value,
	}
	return v.barrier.Put(nested)
}

// logical.Storage impl.
func (v *BarrierView) Delete(key string) error {
	if err := v.sanityCheck(key); err != nil {
		return err
	}
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

// ScanView is used to scan all the keys in a view iteratively
func ScanView(view *BarrierView, cb func(path string)) error {
	frontier := []string{""}
	for len(frontier) > 0 {
		n := len(frontier)
		current := frontier[n-1]
		frontier = frontier[:n-1]

		// List the contents
		contents, err := view.List(current)
		if err != nil {
			return fmt.Errorf("list failed at path '%s': %v", current, err)
		}

		// Handle the contents in the directory
		for _, c := range contents {
			fullPath := current + c
			if strings.HasSuffix(c, "/") {
				frontier = append(frontier, fullPath)
			} else {
				cb(fullPath)
			}
		}
	}
	return nil
}

// CollectKeys is used to collect all the keys in a view
func CollectKeys(view *BarrierView) ([]string, error) {
	// Accumulate the keys
	var existing []string
	cb := func(path string) {
		existing = append(existing, path)
	}

	// Scan for all the keys
	if err := ScanView(view, cb); err != nil {
		return nil, err
	}
	return existing, nil
}

// ClearView is used to delete all the keys in a view
func ClearView(view *BarrierView) error {
	// Collect all the keys
	keys, err := CollectKeys(view)
	if err != nil {
		return err
	}

	// Delete all the keys
	for _, key := range keys {
		if err := view.Delete(key); err != nil {
			return err
		}
	}
	return nil
}
