// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package physical

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
)

const (
	// DefaultErrorPercent is used to determin how often we error
	DefaultErrorPercent = 20
)

// ErrorInjector is used to add errors into underlying physical requests
type ErrorInjector struct {
	backend      Backend
	errorPercent int
	randomLock   *sync.Mutex
	random       *rand.Rand
}

// TransactionalErrorInjector is the transactional version of the error
// injector
type TransactionalErrorInjector struct {
	*ErrorInjector
	Transactional
}

// Verify ErrorInjector satisfies the correct interfaces
var (
	_ Backend       = (*ErrorInjector)(nil)
	_ Transactional = (*TransactionalErrorInjector)(nil)
)

// NewErrorInjector returns a wrapped physical backend to inject error
func NewErrorInjector(b Backend, errorPercent int, logger log.Logger) *ErrorInjector {
	if errorPercent < 0 || errorPercent > 100 {
		errorPercent = DefaultErrorPercent
	}
	logger.Info("creating error injector")

	return &ErrorInjector{
		backend:      b,
		errorPercent: errorPercent,
		randomLock:   new(sync.Mutex),
		random:       rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
	}
}

// NewTransactionalErrorInjector creates a new transactional ErrorInjector
func NewTransactionalErrorInjector(b Backend, errorPercent int, logger log.Logger) *TransactionalErrorInjector {
	return &TransactionalErrorInjector{
		ErrorInjector: NewErrorInjector(b, errorPercent, logger),
		Transactional: b.(Transactional),
	}
}

func (e *ErrorInjector) SetErrorPercentage(p int) {
	e.errorPercent = p
}

func (e *ErrorInjector) addError() error {
	e.randomLock.Lock()
	roll := e.random.Intn(100)
	e.randomLock.Unlock()
	if roll < e.errorPercent {
		return errors.New("random error")
	}

	return nil
}

func (e *ErrorInjector) Put(ctx context.Context, entry *Entry) error {
	if err := e.addError(); err != nil {
		return err
	}
	return e.backend.Put(ctx, entry)
}

func (e *ErrorInjector) Get(ctx context.Context, key string) (*Entry, error) {
	if err := e.addError(); err != nil {
		return nil, err
	}
	return e.backend.Get(ctx, key)
}

func (e *ErrorInjector) Delete(ctx context.Context, key string) error {
	if err := e.addError(); err != nil {
		return err
	}
	return e.backend.Delete(ctx, key)
}

func (e *ErrorInjector) List(ctx context.Context, prefix string) ([]string, error) {
	if err := e.addError(); err != nil {
		return nil, err
	}
	return e.backend.List(ctx, prefix)
}

func (e *TransactionalErrorInjector) Transaction(ctx context.Context, txns []*TxnEntry) error {
	if err := e.addError(); err != nil {
		return err
	}
	return e.Transactional.Transaction(ctx, txns)
}

// TransactionLimits implements physical.TransactionalLimits
func (e *TransactionalErrorInjector) TransactionLimits() (int, int) {
	if tl, ok := e.Transactional.(TransactionalLimits); ok {
		return tl.TransactionLimits()
	}
	// We don't have any specific limits of our own so return zeros to signal that
	// the caller should use whatever reasonable defaults it would if it used a
	// non-TransactionalLimits backend.
	return 0, 0
}
