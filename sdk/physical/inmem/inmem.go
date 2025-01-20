// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package inmem

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-radix"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
	uberAtomic "go.uber.org/atomic"
)

// Verify interfaces are satisfied
var (
	_ physical.Backend             = (*InmemBackend)(nil)
	_ physical.HABackend           = (*InmemHABackend)(nil)
	_ physical.HABackend           = (*TransactionalInmemHABackend)(nil)
	_ physical.Lock                = (*InmemLock)(nil)
	_ physical.Transactional       = (*TransactionalInmemBackend)(nil)
	_ physical.Transactional       = (*TransactionalInmemHABackend)(nil)
	_ physical.TransactionalLimits = (*TransactionalInmemBackend)(nil)
)

var (
	PutDisabledError      = errors.New("put operations disabled in inmem backend")
	GetDisabledError      = errors.New("get operations disabled in inmem backend")
	DeleteDisabledError   = errors.New("delete operations disabled in inmem backend")
	ListDisabledError     = errors.New("list operations disabled in inmem backend")
	GetInTxnDisabledError = errors.New("get operations inside transactions are disabled in inmem backend")
)

// InmemBackend is an in-memory only physical backend. It is useful
// for testing and development situations where the data is not
// expected to be durable.
type InmemBackend struct {
	sync.RWMutex
	root         *radix.Tree
	permitPool   *physical.PermitPool
	logger       log.Logger
	failGet      *uint32
	failPut      *uint32
	failDelete   *uint32
	failList     *uint32
	failGetInTxn *uint32
	logOps       bool
	maxValueSize int
	writeLatency time.Duration
}

type TransactionalInmemBackend struct {
	InmemBackend

	// Using Uber atomic because our SemGrep rules don't like the old pointer
	// trick we used above any more even though it's fine. The newer sync/atomic
	// types are almost the same, but lack was to initialize them cleanly in New*
	// functions so sticking with what SemGrep likes for now.
	maxBatchEntries *uberAtomic.Int32
	maxBatchSize    *uberAtomic.Int32

	largestBatchLen  *uberAtomic.Uint64
	largestBatchSize *uberAtomic.Uint64
}

// NewInmem constructs a new in-memory backend
func NewInmem(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	maxValueSize := 0
	maxValueSizeStr, ok := conf["max_value_size"]
	if ok {
		var err error
		maxValueSize, err = strconv.Atoi(maxValueSizeStr)
		if err != nil {
			return nil, err
		}
	}

	return &InmemBackend{
		root:         radix.New(),
		permitPool:   physical.NewPermitPool(physical.DefaultParallelOperations),
		logger:       logger,
		failGet:      new(uint32),
		failPut:      new(uint32),
		failDelete:   new(uint32),
		failList:     new(uint32),
		failGetInTxn: new(uint32),
		logOps:       os.Getenv("VAULT_INMEM_LOG_ALL_OPS") != "",
		maxValueSize: maxValueSize,
	}, nil
}

// Basically for now just creates a permit pool of size 1 so only one operation
// can run at a time
func NewTransactionalInmem(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	maxValueSize := 0
	maxValueSizeStr, ok := conf["max_value_size"]
	if ok {
		var err error
		maxValueSize, err = strconv.Atoi(maxValueSizeStr)
		if err != nil {
			return nil, err
		}
	}

	return &TransactionalInmemBackend{
		InmemBackend: InmemBackend{
			root:         radix.New(),
			permitPool:   physical.NewPermitPool(1),
			logger:       logger,
			failGet:      new(uint32),
			failPut:      new(uint32),
			failDelete:   new(uint32),
			failList:     new(uint32),
			failGetInTxn: new(uint32),
			logOps:       os.Getenv("VAULT_INMEM_LOG_ALL_OPS") != "",
			maxValueSize: maxValueSize,
		},

		maxBatchEntries:  uberAtomic.NewInt32(64),
		maxBatchSize:     uberAtomic.NewInt32(128 * 1024),
		largestBatchLen:  uberAtomic.NewUint64(0),
		largestBatchSize: uberAtomic.NewUint64(0),
	}, nil
}

// SetWriteLatency add a sleep to each Put/Delete operation (and each op in a
// transaction for a TransactionalInmemBackend). It's not so much to simulate
// real disk latency as much as to make the go runtime schedule things more like
// a real disk where concurrent write operations are more likely to interleave
// as each one blocks on disk IO. Set to 0 to disable again (the default).
func (i *InmemBackend) SetWriteLatency(latency time.Duration) {
	i.Lock()
	defer i.Unlock()
	i.writeLatency = latency
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
	if i.logOps {
		i.logger.Trace("put", "key", entry.Key)
	}
	if atomic.LoadUint32(i.failPut) != 0 {
		return PutDisabledError
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if i.maxValueSize > 0 && len(entry.Value) > i.maxValueSize {
		return fmt.Errorf("%s", physical.ErrValueTooLarge)
	}

	i.root.Insert(entry.Key, entry.Value)
	if i.writeLatency > 0 {
		time.Sleep(i.writeLatency)
	}
	return nil
}

func (i *InmemBackend) FailPut(fail bool) {
	var val uint32
	if fail {
		val = 1
	}
	atomic.StoreUint32(i.failPut, val)
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
	if i.logOps {
		i.logger.Trace("get", "key", key)
	}
	if atomic.LoadUint32(i.failGet) != 0 {
		return nil, GetDisabledError
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if raw, ok := i.root.Get(key); ok {
		return &physical.Entry{
			Key:   key,
			Value: raw.([]byte),
		}, nil
	}
	return nil, nil
}

func (i *InmemBackend) FailGet(fail bool) {
	var val uint32
	if fail {
		val = 1
	}
	atomic.StoreUint32(i.failGet, val)
}

func (i *InmemBackend) FailGetInTxn(fail bool) {
	var val uint32
	if fail {
		val = 1
	}
	atomic.StoreUint32(i.failGetInTxn, val)
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
	if i.logOps {
		i.logger.Trace("delete", "key", key)
	}
	if atomic.LoadUint32(i.failDelete) != 0 {
		return DeleteDisabledError
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	i.root.Delete(key)
	if i.writeLatency > 0 {
		time.Sleep(i.writeLatency)
	}
	return nil
}

func (i *InmemBackend) FailDelete(fail bool) {
	var val uint32
	if fail {
		val = 1
	}
	atomic.StoreUint32(i.failDelete, val)
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (i *InmemBackend) List(ctx context.Context, prefix string) ([]string, error) {
	i.permitPool.Acquire()
	defer i.permitPool.Release()

	i.RLock()
	defer i.RUnlock()

	return i.ListInternal(ctx, prefix)
}

func (i *InmemBackend) ListInternal(ctx context.Context, prefix string) ([]string, error) {
	if i.logOps {
		i.logger.Trace("list", "prefix", prefix)
	}
	if atomic.LoadUint32(i.failList) != 0 {
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

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return out, nil
}

func (i *InmemBackend) FailList(fail bool) {
	var val uint32
	if fail {
		val = 1
	}
	atomic.StoreUint32(i.failList, val)
}

// Transaction implements the transaction interface
func (t *TransactionalInmemBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	t.permitPool.Acquire()
	defer t.permitPool.Release()

	t.Lock()
	defer t.Unlock()

	failGetInTxn := atomic.LoadUint32(t.failGetInTxn)
	size := uint64(0)
	for _, t := range txns {
		// We use 2x key length to match the logic in WALBackend.persistWALs
		// presumably this is attempting to account for some amount of encoding
		// overhead.
		size += uint64(2*len(t.Entry.Key) + len(t.Entry.Value))
		if t.Operation == physical.GetOperation && failGetInTxn != 0 {
			return GetInTxnDisabledError
		}
	}

	if size > t.largestBatchSize.Load() {
		t.largestBatchSize.Store(size)
	}
	if len(txns) > int(t.largestBatchLen.Load()) {
		t.largestBatchLen.Store(uint64(len(txns)))
	}

	return physical.GenericTransactionHandler(ctx, t, txns)
}

func (t *TransactionalInmemBackend) SetMaxBatchEntries(entries int) {
	t.maxBatchEntries.Store(int32(entries))
}

func (t *TransactionalInmemBackend) SetMaxBatchSize(entries int) {
	t.maxBatchSize.Store(int32(entries))
}

func (t *TransactionalInmemBackend) TransactionLimits() (int, int) {
	return int(t.maxBatchEntries.Load()), int(t.maxBatchSize.Load())
}

func (t *TransactionalInmemBackend) BatchStats() (maxEntries uint64, maxSize uint64) {
	return t.largestBatchLen.Load(), t.largestBatchSize.Load()
}
