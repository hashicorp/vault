package inmem

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	radix "github.com/armon/go-radix"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"
)

type faultyPseudo struct {
	underlying  InmemBackend
	faultyPaths map[string]struct{}
}

func (f *faultyPseudo) Get(key string) (*physical.Entry, error) {
	return f.underlying.Get(key)
}

func (f *faultyPseudo) Put(entry *physical.Entry) error {
	return f.underlying.Put(entry)
}

func (f *faultyPseudo) Delete(key string) error {
	return f.underlying.Delete(key)
}

func (f *faultyPseudo) GetInternal(key string) (*physical.Entry, error) {
	if _, ok := f.faultyPaths[key]; ok {
		return nil, fmt.Errorf("fault")
	}
	return f.underlying.GetInternal(key)
}

func (f *faultyPseudo) PutInternal(entry *physical.Entry) error {
	if _, ok := f.faultyPaths[entry.Key]; ok {
		return fmt.Errorf("fault")
	}
	return f.underlying.PutInternal(entry)
}

func (f *faultyPseudo) DeleteInternal(key string) error {
	if _, ok := f.faultyPaths[key]; ok {
		return fmt.Errorf("fault")
	}
	return f.underlying.DeleteInternal(key)
}

func (f *faultyPseudo) List(prefix string) ([]string, error) {
	return f.underlying.List(prefix)
}

func (f *faultyPseudo) Transaction(txns []*physical.TxnEntry) error {
	f.underlying.permitPool.Acquire()
	defer f.underlying.permitPool.Release()

	f.underlying.Lock()
	defer f.underlying.Unlock()

	return physical.GenericTransactionHandler(f, txns)
}

func newFaultyPseudo(logger log.Logger, faultyPaths []string) *faultyPseudo {
	out := &faultyPseudo{
		underlying: InmemBackend{
			root:       radix.New(),
			permitPool: physical.NewPermitPool(1),
			logger:     logger,
		},
		faultyPaths: make(map[string]struct{}, len(faultyPaths)),
	}
	for _, v := range faultyPaths {
		out.faultyPaths[v] = struct{}{}
	}
	return out
}

func TestPseudo_Basic(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)
	p := newFaultyPseudo(logger, nil)
	physical.ExerciseBackend(t, p)
	physical.ExerciseBackend_ListPrefix(t, p)
}

func TestPseudo_SuccessfulTransaction(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)
	p := newFaultyPseudo(logger, nil)

	physical.ExerciseTransactionalBackend(t, p)
}

func TestPseudo_FailedTransaction(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)
	p := newFaultyPseudo(logger, []string{"zip"})

	txns := physical.SetupTestingTransactions(t, p)
	if err := p.Transaction(txns); err == nil {
		t.Fatal("expected error during transaction")
	}

	keys, err := p.List("")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"foo", "zip", "deleteme", "deleteme2"}

	sort.Strings(keys)
	sort.Strings(expected)
	if !reflect.DeepEqual(keys, expected) {
		t.Fatalf("mismatch: expected\n%#v\ngot\n%#v\n", expected, keys)
	}

	entry, err := p.Get("foo")
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

	entry, err = p.Get("zip")
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
