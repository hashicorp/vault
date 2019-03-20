package vault

import (
	"context"
	"errors"
	"sync"

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
	storage         *logical.StorageView
	readOnlyErr     error
	readOnlyErrLock sync.RWMutex
	iCheck          interface{}
}

// NewBarrierView takes an underlying security barrier and returns
// a view of it that can only operate with the given prefix.
func NewBarrierView(barrier logical.Storage, prefix string) *BarrierView {
	return &BarrierView{
		storage: logical.NewStorageView(barrier, prefix),
	}
}

func (v *BarrierView) setICheck(iCheck interface{}) {
	v.iCheck = iCheck
}

func (v *BarrierView) setReadOnlyErr(readOnlyErr error) {
	v.readOnlyErrLock.Lock()
	defer v.readOnlyErrLock.Unlock()
	v.readOnlyErr = readOnlyErr
}

func (v *BarrierView) getReadOnlyErr() error {
	v.readOnlyErrLock.RLock()
	defer v.readOnlyErrLock.RUnlock()
	return v.readOnlyErr
}

func (v *BarrierView) Prefix() string {
	return v.storage.Prefix()
}

func (v *BarrierView) List(ctx context.Context, prefix string) ([]string, error) {
	return v.storage.List(ctx, prefix)
}

func (v *BarrierView) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	return v.storage.Get(ctx, key)
}

// Put differs from List/Get because it checks read-only errors
func (v *BarrierView) Put(ctx context.Context, entry *logical.StorageEntry) error {
	if entry == nil {
		return errors.New("cannot write nil entry")
	}

	expandedKey := v.storage.ExpandKey(entry.Key)

	roErr := v.getReadOnlyErr()
	if roErr != nil {
		if runICheck(v, expandedKey, roErr) {
			return roErr
		}
	}

	return v.storage.Put(ctx, entry)
}

// logical.Storage impl.
func (v *BarrierView) Delete(ctx context.Context, key string) error {
	expandedKey := v.storage.ExpandKey(key)

	roErr := v.getReadOnlyErr()
	if roErr != nil {
		if runICheck(v, expandedKey, roErr) {
			return roErr
		}
	}

	return v.storage.Delete(ctx, key)
}

// SubView constructs a nested sub-view using the given prefix
func (v *BarrierView) SubView(prefix string) *BarrierView {
	return &BarrierView{
		storage:     v.storage.SubView(prefix),
		readOnlyErr: v.getReadOnlyErr(),
		iCheck:      v.iCheck,
	}
}
