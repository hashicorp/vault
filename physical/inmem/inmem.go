package inmem

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"

	"github.com/armon/go-radix"
)

// Verify interfaces are satisfied
var _ physical.Backend = (*InmemBackend)(nil)
var _ physical.HABackend = (*InmemHABackend)(nil)
var _ physical.HABackend = (*TransactionalInmemHABackend)(nil)
var _ physical.Lock = (*InmemLock)(nil)
var _ physical.Transactional = (*TransactionalInmemBackend)(nil)
var _ physical.Transactional = (*TransactionalInmemHABackend)(nil)

var (
	PutDisabledError    = errors.New("put operations disabled in inmem backend")
	GetDisabledError    = errors.New("get operations disabled in inmem backend")
	DeleteDisabledError = errors.New("delete operations disabled in inmem backend")
	ListDisabledError   = errors.New("list operations disabled in inmem backend")
)

// InmemBackend is an in-memory only physical backend. It is useful
// for testing and development situations where the data is not
// expected to be durable.
type InmemBackend struct {
	sync.RWMutex
	root       *radix.Tree
	permitPool *physical.PermitPool
	logger     log.Logger
	FailGet    *uint32
	FailPut    *uint32
	FailDelete *uint32
	FailList   *uint32
}

type TransactionalInmemBackend struct {
	InmemBackend
}

// NewInmem constructs a new in-memory backend
func NewInmem(_ map[string]string, logger log.Logger) (physical.Backend, error) {
	in := &InmemBackend{
		root:       radix.New(),
		permitPool: physical.NewPermitPool(physical.DefaultParallelOperations),
		logger:     logger,
		FailGet:    new(uint32),
		FailPut:    new(uint32),
		FailDelete: new(uint32),
		FailList:   new(uint32),
	}
	return in, nil
}

// Basically for now just creates a permit pool of size 1 so only one operation
// can run at a time
func NewTransactionalInmem(_ map[string]string, logger log.Logger) (physical.Backend, error) {
	in := &TransactionalInmemBackend{
		InmemBackend: InmemBackend{
			root:       radix.New(),
			permitPool: physical.NewPermitPool(1),
			logger:     logger,
			FailGet:    new(uint32),
			FailPut:    new(uint32),
			FailDelete: new(uint32),
			FailList:   new(uint32),
		},
	}
	return in, nil
}

// Put is used to insert or update an entry
func (i *InmemBackend) Put(ctx context.Context, entry *physical.Entry) error {
	i.permitPool.Acquire()
	defer i.permitPool.Release()

	i.Lock()
	defer i.Unlock()

	return i.PutInternal(ctx, entry)
}

func (i *InmemBackend) PutInternal(ctx context.Context, entry *physical.Entry) error {
	if i.FailPut != nil && atomic.LoadUint32(i.FailPut) != 0 {
		return PutDisabledError
	}

	i.root.Insert(entry.Key, entry.Value)
	return nil
}

// Get is used to fetch an entry
func (i *InmemBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	i.permitPool.Acquire()
	defer i.permitPool.Release()

	i.RLock()
	defer i.RUnlock()

	return i.GetInternal(ctx, key)
}

func (i *InmemBackend) GetInternal(ctx context.Context, key string) (*physical.Entry, error) {
	if i.FailGet != nil && atomic.LoadUint32(i.FailGet) != 0 {
		return nil, GetDisabledError
	}

	if raw, ok := i.root.Get(key); ok {
		return &physical.Entry{
			Key:   key,
			Value: raw.([]byte),
		}, nil
	}
	return nil, nil
}

// Delete is used to permanently delete an entry
func (i *InmemBackend) Delete(ctx context.Context, key string) error {
	i.permitPool.Acquire()
	defer i.permitPool.Release()

	i.Lock()
	defer i.Unlock()

	return i.DeleteInternal(ctx, key)
}

func (i *InmemBackend) DeleteInternal(ctx context.Context, key string) error {
	if i.FailDelete != nil && atomic.LoadUint32(i.FailDelete) != 0 {
		return DeleteDisabledError
	}

	i.root.Delete(key)
	return nil
}

// List is used ot list all the keys under a given
// prefix, up to the next prefix.
func (i *InmemBackend) List(ctx context.Context, prefix string) ([]string, error) {
	i.permitPool.Acquire()
	defer i.permitPool.Release()

	i.RLock()
	defer i.RUnlock()

	return i.ListInternal(prefix)
}

func (i *InmemBackend) ListInternal(prefix string) ([]string, error) {
	if i.FailList != nil && atomic.LoadUint32(i.FailList) != 0 {
		return nil, ListDisabledError
	}

	var out []string
	seen := make(map[string]interface{})
	walkFn := func(s string, v interface{}) bool {
		trimmed := strings.TrimPrefix(s, prefix)
		sep := strings.Index(trimmed, "/")
		if sep == -1 {
			out = append(out, trimmed)
		} else {
			trimmed = trimmed[:sep+1]
			if _, ok := seen[trimmed]; !ok {
				out = append(out, trimmed)
				seen[trimmed] = struct{}{}
			}
		}
		return false
	}
	i.root.WalkPrefix(prefix, walkFn)

	return out, nil
}

// Implements the transaction interface
func (t *TransactionalInmemBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	t.permitPool.Acquire()
	defer t.permitPool.Release()

	t.Lock()
	defer t.Unlock()

	return physical.GenericTransactionHandler(ctx, t, txns)
}
