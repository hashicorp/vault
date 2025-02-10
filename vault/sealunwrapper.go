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
		underlying: underlying,
		logger:     logger,
		locks:      locksutil.CreateLocks(),
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
	allowUnwraps atomic.Bool
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

// unwrap gets an entry from underlying storage and tries to unwrap it.
// - If the entry is not wrapped: the entry will be returned unchanged and wasWrapped will be false
// - If the entry is wrapped and encrypted: an error is returned.
// - If the entry is wrapped but not encrypted: the entry will be unwrapped and returned. wasWrapped will be true.
func (d *sealUnwrapper) unwrap(ctx context.Context, key string) (unwrappedEntry *physical.Entry, wasWrapped bool, err error) {
	entry, err := d.underlying.Get(ctx, key)
	if err != nil {
		return nil, false, err
	}
	if entry == nil {
		return nil, false, nil
	}

	wrappedEntryValue, unmarshaled := UnmarshalSealWrappedValueWithCanary(entry.Value)
	switch {
	case !unmarshaled:
		// Entry is not wrapped
		return entry, false, nil
	case wrappedEntryValue.isEncrypted():
		// Entry is wrapped and encrypted
		return nil, true, fmt.Errorf("cannot decode sealwrapped storage entry %q", entry.Key)
	default:
		// Entry is wrapped and not encrypted
		pt, err := wrappedEntryValue.getPlaintextValue()
		if err != nil {
			return nil, true, err
		}
		return &physical.Entry{
			Key:   entry.Key,
			Value: pt,
		}, true, nil
	}
}

func (d *sealUnwrapper) Get(ctx context.Context, key string) (*physical.Entry, error) {
	entry, wasWrapped, err := d.unwrap(ctx, key)
	switch {
	case err != nil: // Failed to get entry
		return nil, err
	case entry == nil: // Entry doesn't exist
		return nil, nil
	case !wasWrapped || !d.allowUnwraps.Load(): // Entry was not wrapped or unwrapping not allowed
		return entry, nil
	}

	// Entry was wrapped, we need to replace it with the unwrapped value

	// Grab locks because we are performing a write
	locksutil.LockForKey(d.locks, key).Lock()
	defer locksutil.LockForKey(d.locks, key).Unlock()

	// Read entry again in case it was changed while we were waiting for the lock
	entry, wasWrapped, err = d.unwrap(ctx, key)
	switch {
	case err != nil: // Failed to get entry
		return nil, err
	case entry == nil: // Entry doesn't exist
		return nil, nil
	case !wasWrapped || !d.allowUnwraps.Load(): // Entry was not wrapped or unwrapping not allowed
		return entry, nil
	}

	// Write out the unwrapped value
	return entry, d.underlying.Put(ctx, entry)
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
	d.allowUnwraps.Store(false)
}

func (d *sealUnwrapper) runUnwraps() {
	// Allow key unwraps on key gets. This gets set only when running on the
	// active node to prevent standbys from changing data underneath the
	// primary
	d.allowUnwraps.Store(true)
}
