// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package physical

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
)

// PathErrorInjector wraps a physical backend and injects errors on storage
// operations whose keys match configured path prefixes. Unlike ErrorInjector,
// it does not support a global error rate — errors are only injected for
// paths that have been explicitly configured via SetErrorPercentageForPath.
//
// This is useful for tests that need surgical failure injection on specific
// storage paths (e.g. packer buckets) without affecting unrelated operations.
type PathErrorInjector struct {
	backend    Backend
	pathErrors map[string]int
	mu         sync.RWMutex
	randMu     sync.Mutex
	random     *rand.Rand
}

// TransactionalPathErrorInjector is the transactional version of
// PathErrorInjector.
type TransactionalPathErrorInjector struct {
	*PathErrorInjector
	Transactional
}

var (
	_ Backend       = (*PathErrorInjector)(nil)
	_ Transactional = (*TransactionalPathErrorInjector)(nil)
)

// NewPathErrorInjector creates a new PathErrorInjector wrapping the given
// backend. No errors are injected until SetErrorPercentageForPath is called.
func NewPathErrorInjector(b Backend, logger log.Logger) *PathErrorInjector {
	logger.Info("creating path error injector")
	return &PathErrorInjector{
		backend:    b,
		pathErrors: make(map[string]int),
		random:     rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
	}
}

// NewTransactionalPathErrorInjector creates a new transactional
// PathErrorInjector wrapping the given backend.
func NewTransactionalPathErrorInjector(b Backend, logger log.Logger) *TransactionalPathErrorInjector {
	return &TransactionalPathErrorInjector{
		PathErrorInjector: NewPathErrorInjector(b, logger),
		Transactional:     b.(Transactional),
	}
}

// SetErrorPercentageForPath sets an error injection percentage for a specific
// path prefix. Any storage operation whose key starts with the given prefix
// will fail with the configured probability (0-100). The longest matching
// prefix wins when multiple prefixes match a key. This method is safe for
// concurrent use.
func (e *PathErrorInjector) SetErrorPercentageForPath(path string, p int) {
	e.mu.Lock()
	e.pathErrors[path] = p
	e.mu.Unlock()
}

func (e *PathErrorInjector) addError(key string) error {
	e.mu.RLock()
	percent := 0
	longestMatch := 0
	for prefix, p := range e.pathErrors {
		if strings.HasPrefix(key, prefix) && len(prefix) > longestMatch {
			longestMatch = len(prefix)
			percent = p
		}
	}
	e.mu.RUnlock()

	if percent == 0 {
		return nil
	}

	e.randMu.Lock()
	roll := e.random.Intn(100)
	e.randMu.Unlock()

	if roll < percent {
		return errors.New("random error")
	}
	return nil
}

func (e *PathErrorInjector) Put(ctx context.Context, entry *Entry) error {
	if err := e.addError(entry.Key); err != nil {
		return err
	}
	return e.backend.Put(ctx, entry)
}

func (e *PathErrorInjector) Get(ctx context.Context, key string) (*Entry, error) {
	if err := e.addError(key); err != nil {
		return nil, err
	}
	return e.backend.Get(ctx, key)
}

func (e *PathErrorInjector) Delete(ctx context.Context, key string) error {
	if err := e.addError(key); err != nil {
		return err
	}
	return e.backend.Delete(ctx, key)
}

func (e *PathErrorInjector) List(ctx context.Context, prefix string) ([]string, error) {
	if err := e.addError(prefix); err != nil {
		return nil, err
	}
	return e.backend.List(ctx, prefix)
}

func (e *TransactionalPathErrorInjector) Transaction(ctx context.Context, txns []*TxnEntry) error {
	for _, txn := range txns {
		if txn != nil {
			if err := e.addError(txn.Entry.Key); err != nil {
				return err
			}
		}
	}
	return e.Transactional.Transaction(ctx, txns)
}

// TransactionLimits implements physical.TransactionalLimits
func (e *TransactionalPathErrorInjector) TransactionLimits() (int, int) {
	if tl, ok := e.Transactional.(TransactionalLimits); ok {
		return tl.TransactionLimits()
	}
	return 0, 0
}
