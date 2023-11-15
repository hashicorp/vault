// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"fmt"
	"sync/atomic"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/physical"
)

// NewSealUnwrapper creates a new seal unwrapper
func NewSealUnwrapper(underlying physical.Backend, logger log.Logger) physical.Backend {
	ret := &sealUnwrapper{
		underlying:   underlying,
		logger:       logger,
		locks:        locksutil.CreateLocks(),
		allowUnwraps: new(uint32),
	}

	if underTxn, ok := underlying.(physical.Transactional); ok {
		return &transactionalSealUnwrapper{
			sealUnwrapper: ret,
			Transactional: underTxn,
		}
	}

	return ret
}

var (
	_ physical.Backend       = (*sealUnwrapper)(nil)
	_ physical.Transactional = (*transactionalSealUnwrapper)(nil)
)

type sealUnwrapper struct {
	underlying   physical.Backend
	logger       log.Logger
	locks        []*locksutil.LockEntry
	allowUnwraps *uint32
}

// transactionalSealUnwrapper is a seal unwrapper that wraps a physical that is transactional
type transactionalSealUnwrapper struct {
	*sealUnwrapper
	physical.Transactional
}

func (d *sealUnwrapper) Put(ctx context.Context, entry *physical.Entry) error {
	if entry == nil {
		return nil
	}

	locksutil.LockForKey(d.locks, entry.Key).Lock()
	defer locksutil.LockForKey(d.locks, entry.Key).Unlock()

	return d.underlying.Put(ctx, entry)
}

// unwrap gets an entry from underlying storage and tries to unwrap it. If the entry was not wrapped, return
// value unwrappedEntry will be nil. If the entry is wrapped and encrypted, an error is returned.
func (d *sealUnwrapper) unwrap(ctx context.Context, key string) (entry, unwrappedEntry *physical.Entry, err error) {
	entry, err = d.underlying.Get(ctx, key)
	if err != nil {
		return nil, nil, err
	}
	if entry == nil {
		return nil, nil, err
	}

	wrappedEntryValue, unmarshaled := UnmarshalSealWrappedValueWithCanary(entry.Value)
	switch {
	case !unmarshaled:
		unwrappedEntry = entry
	case wrappedEntryValue.isEncrypted():
		return nil, nil, fmt.Errorf("cannot decode sealwrapped storage entry %q", entry.Key)
	default:
		pt, err := wrappedEntryValue.getPlaintextValue()
		if err != nil {
			return nil, nil, err
		}
		unwrappedEntry = &physical.Entry{
			Key:   entry.Key,
			Value: pt,
		}
	}

	return entry, unwrappedEntry, nil
}

func (d *sealUnwrapper) Get(ctx context.Context, key string) (*physical.Entry, error) {
	entry, unwrappedEntry, err := d.unwrap(ctx, key)
	switch {
	case err != nil:
		return nil, err
	case entry == nil:
		return nil, nil
	case atomic.LoadUint32(d.allowUnwraps) != 1:
		return unwrappedEntry, nil
	}

	locksutil.LockForKey(d.locks, key).Lock()
	defer locksutil.LockForKey(d.locks, key).Unlock()

	// At this point we need to re-read and re-check
	entry, unwrappedEntry, err = d.unwrap(ctx, key)
	switch {
	case err != nil:
		return nil, err
	case entry == nil:
		return nil, nil
	case atomic.LoadUint32(d.allowUnwraps) != 1:
		return unwrappedEntry, nil
	}

	return unwrappedEntry, d.underlying.Put(ctx, unwrappedEntry)
}

func (d *sealUnwrapper) Delete(ctx context.Context, key string) error {
	locksutil.LockForKey(d.locks, key).Lock()
	defer locksutil.LockForKey(d.locks, key).Unlock()

	return d.underlying.Delete(ctx, key)
}

func (d *sealUnwrapper) List(ctx context.Context, prefix string) ([]string, error) {
	return d.underlying.List(ctx, prefix)
}

func (d *transactionalSealUnwrapper) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	// Collect keys that need to be locked
	var keys []string
	for _, curr := range txns {
		keys = append(keys, curr.Entry.Key)
	}
	// Lock the keys
	for _, l := range locksutil.LocksForKeys(d.locks, keys) {
		l.Lock()
		defer l.Unlock()
	}

	if err := d.Transactional.Transaction(ctx, txns); err != nil {
		return err
	}

	return nil
}

// This should only run during preSeal which ensures that it can't be run
// concurrently and that it will be run only by the active node
func (d *sealUnwrapper) stopUnwraps() {
	atomic.StoreUint32(d.allowUnwraps, 0)
}

func (d *sealUnwrapper) runUnwraps() {
	// Allow key unwraps on key gets. This gets set only when running on the
	// active node to prevent standbys from changing data underneath the
	// primary
	atomic.StoreUint32(d.allowUnwraps, 1)
}
