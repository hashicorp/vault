package physical

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	radix "github.com/armon/go-radix"
	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

type faultyPseudo struct {
	underlying  InmemBackend
	faultyPaths map[string]struct{}
}

func (f *faultyPseudo) Get(key string) (*Entry, error) {
	return f.underlying.Get(key)
}

func (f *faultyPseudo) Put(entry *Entry) error {
	return f.underlying.Put(entry)
}

func (f *faultyPseudo) Delete(key string) error {
	return f.underlying.Delete(key)
}

func (f *faultyPseudo) GetInternal(key string) (*Entry, error) {
	if _, ok := f.faultyPaths[key]; ok {
		return nil, fmt.Errorf("fault")
	}
	return f.underlying.GetInternal(key)
}

func (f *faultyPseudo) PutInternal(entry *Entry) error {
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

func (f *faultyPseudo) Transaction(txns []TxnEntry) error {
	f.underlying.permitPool.Acquire()
	defer f.underlying.permitPool.Release()

	f.underlying.Lock()
	defer f.underlying.Unlock()

	return genericTransactionHandler(f, txns)
}

func newFaultyPseudo(logger log.Logger, faultyPaths []string) *faultyPseudo {
	out := &faultyPseudo{
		underlying: InmemBackend{
			root:       radix.New(),
			permitPool: NewPermitPool(1),
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
	testBackend(t, p)
	testBackend_ListPrefix(t, p)
}

func TestPseudo_SuccessfulTransaction(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)
	p := newFaultyPseudo(logger, nil)

	txns := setupPseudo(p, t)

	if err := p.Transaction(txns); err != nil {
		t.Fatal(err)
	}

	keys, err := p.List("")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"foo", "zip"}

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
	if string(entry.Value) != "bar3" {
		t.Fatal("updates did not apply correctly")
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
	if string(entry.Value) != "zap3" {
		t.Fatal("updates did not apply correctly")
	}
}

func TestPseudo_FailedTransaction(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)
	p := newFaultyPseudo(logger, []string{"zip"})

	txns := setupPseudo(p, t)

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

func setupPseudo(p *faultyPseudo, t *testing.T) []TxnEntry {
	// Add a few keys so that we test rollback with deletion
	if err := p.Put(&Entry{
		Key:   "foo",
		Value: []byte("bar"),
	}); err != nil {
		t.Fatal(err)
	}
	if err := p.Put(&Entry{
		Key:   "zip",
		Value: []byte("zap"),
	}); err != nil {
		t.Fatal(err)
	}
	if err := p.Put(&Entry{
		Key: "deleteme",
	}); err != nil {
		t.Fatal(err)
	}
	if err := p.Put(&Entry{
		Key: "deleteme2",
	}); err != nil {
		t.Fatal(err)
	}

	txns := []TxnEntry{
		TxnEntry{
			Operation: PutOperation,
			Entry: &Entry{
				Key:   "foo",
				Value: []byte("bar2"),
			},
		},
		TxnEntry{
			Operation: DeleteOperation,
			Entry: &Entry{
				Key: "deleteme",
			},
		},
		TxnEntry{
			Operation: PutOperation,
			Entry: &Entry{
				Key:   "foo",
				Value: []byte("bar3"),
			},
		},
		TxnEntry{
			Operation: DeleteOperation,
			Entry: &Entry{
				Key: "deleteme2",
			},
		},
		TxnEntry{
			Operation: PutOperation,
			Entry: &Entry{
				Key:   "zip",
				Value: []byte("zap3"),
			},
		},
	}

	return txns
}
