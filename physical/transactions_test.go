package physical

import (
	"fmt"
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
