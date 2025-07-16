// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package inmem

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"testing"

	radix "github.com/armon/go-radix"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/permitpool"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

type faultyPseudo struct {
	underlying  InmemBackend
	faultyPaths map[string]struct{}
}

func (f *faultyPseudo) Get(ctx context.Context, key string) (*physical.Entry, error) {
	return f.underlying.Get(context.Background(), key)
}

func (f *faultyPseudo) Put(ctx context.Context, entry *physical.Entry) error {
	return f.underlying.Put(context.Background(), entry)
}

func (f *faultyPseudo) Delete(ctx context.Context, key string) error {
	return f.underlying.Delete(context.Background(), key)
}

func (f *faultyPseudo) GetInternal(ctx context.Context, key string) (*physical.Entry, error) {
	if _, ok := f.faultyPaths[key]; ok {
		return nil, fmt.Errorf("fault")
	}
	return f.underlying.GetInternal(context.Background(), key)
}

func (f *faultyPseudo) PutInternal(ctx context.Context, entry *physical.Entry) error {
	if _, ok := f.faultyPaths[entry.Key]; ok {
		return fmt.Errorf("fault")
	}
	return f.underlying.PutInternal(context.Background(), entry)
}

func (f *faultyPseudo) DeleteInternal(ctx context.Context, key string) error {
	if _, ok := f.faultyPaths[key]; ok {
		return fmt.Errorf("fault")
	}
	return f.underlying.DeleteInternal(context.Background(), key)
}

func (f *faultyPseudo) List(ctx context.Context, prefix string) ([]string, error) {
	return f.underlying.List(context.Background(), prefix)
}

func (f *faultyPseudo) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	if err := f.underlying.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer f.underlying.permitPool.Release()

	f.underlying.Lock()
	defer f.underlying.Unlock()

	return physical.GenericTransactionHandler(ctx, f, txns)
}

func newFaultyPseudo(logger log.Logger, faultyPaths []string) *faultyPseudo {
	out := &faultyPseudo{
		underlying: InmemBackend{
			root:       radix.New(),
			permitPool: permitpool.New(1),
			logger:     logger.Named("storage.inmembackend"),
			failGet:    new(uint32),
			failPut:    new(uint32),
			failDelete: new(uint32),
			failList:   new(uint32),
		},
		faultyPaths: make(map[string]struct{}, len(faultyPaths)),
	}
	for _, v := range faultyPaths {
		out.faultyPaths[v] = struct{}{}
	}
	return out
}

func TestPseudo_Basic(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)
	p := newFaultyPseudo(logger, nil)
	physical.ExerciseBackend(t, p)
	physical.ExerciseBackend_ListPrefix(t, p)
}

func TestPseudo_SuccessfulTransaction(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)
	p := newFaultyPseudo(logger, nil)

	physical.ExerciseTransactionalBackend(t, p)
}

func TestPseudo_FailedTransaction(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)
	p := newFaultyPseudo(logger, []string{"zip"})

	txns := physical.SetupTestingTransactions(t, p)
	if err := p.Transaction(context.Background(), txns); err == nil {
		t.Fatal("expected error during transaction")
	}

	keys, err := p.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"foo", "zip", "deleteme", "deleteme2"}

	sort.Strings(keys)
	sort.Strings(expected)
	if !reflect.DeepEqual(keys, expected) {
		t.Fatalf("mismatch: expected\n%#v\ngot\n%#v\n", expected, keys)
	}

	entry, err := p.Get(context.Background(), "foo")
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatal("got nil entry")
	}
	if entry.Value == nil {
		t.Fatal("got nil value")
	}
	if string(entry.Value) != "bar" {
		t.Fatal("values did not rollback correctly")
	}

	entry, err = p.Get(context.Background(), "zip")
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatal("got nil entry")
	}
	if entry.Value == nil {
		t.Fatal("got nil value")
	}
	if string(entry.Value) != "zap" {
		t.Fatal("values did not rollback correctly")
	}
}
