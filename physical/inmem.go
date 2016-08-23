package physical

import (
	"strings"
	"sync"

	log "github.com/mgutz/logxi/v1"

	"github.com/armon/go-radix"
)

// InmemBackend is an in-memory only physical backend. It is useful
// for testing and development situations where the data is not
// expected to be durable.
type InmemBackend struct {
	root       *radix.Tree
	l          sync.RWMutex
	permitPool *PermitPool
	logger     log.Logger
}

// NewInmem constructs a new in-memory backend
func NewInmem(logger log.Logger) *InmemBackend {
	in := &InmemBackend{
		root:       radix.New(),
		permitPool: NewPermitPool(DefaultParallelOperations),
	}
	return in
}

// Put is used to insert or update an entry
func (i *InmemBackend) Put(entry *Entry) error {
	i.permitPool.Acquire()
	defer i.permitPool.Release()

	i.l.Lock()
	defer i.l.Unlock()

	i.root.Insert(entry.Key, entry)
	return nil
}

// Get is used to fetch an entry
func (i *InmemBackend) Get(key string) (*Entry, error) {
	i.permitPool.Acquire()
	defer i.permitPool.Release()

	i.l.RLock()
	defer i.l.RUnlock()

	if raw, ok := i.root.Get(key); ok {
		return raw.(*Entry), nil
	}
	return nil, nil
}

// Delete is used to permanently delete an entry
func (i *InmemBackend) Delete(key string) error {
	i.permitPool.Acquire()
	defer i.permitPool.Release()

	i.l.Lock()
	defer i.l.Unlock()

	i.root.Delete(key)
	return nil
}

// List is used ot list all the keys under a given
// prefix, up to the next prefix.
func (i *InmemBackend) List(prefix string) ([]string, error) {
	i.permitPool.Acquire()
	defer i.permitPool.Release()

	i.l.RLock()
	defer i.l.RUnlock()

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
