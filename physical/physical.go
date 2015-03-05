package physical

import "fmt"

// Backend is the interface required for a physical
// backend. A physical backend is used to durably store
// datd outside of Vault. As such, it is completely untrusted,
// and is only accessed via a security barrier. The backends
// must represent keys in a hierarchical manner. All methods
// are expected to be thread safe.
type Backend interface {
	// Put is used to insert or update an entry
	Put(entry *Entry) error

	// Get is used to fetch an entry
	Get(key string) (*Entry, error)

	// Delete is used to permanently delete an entry
	Delete(key string) error

	// List is used ot list all the keys under a given
	// prefix, up to the next prefix.
	List(prefix string) ([]string, error)
}

// Entry is used to represent data stored by the physical backend
type Entry struct {
	Key   string
	Value []byte
}

// Factory is the factory function to create a physical backend.
type Factory func(map[string]string) (Backend, error)

// NewBackend returns a new Bckend with the given type and configuration.
// The backend is looked up in the BuiltinBackends variable.
func NewBackend(t string, conf map[string]string) (Backend, error) {
	f, ok := BuiltinBackends[t]
	if !ok {
		return nil, fmt.Errorf("unknown physical backend type: %s", t)
	}
	return f(conf)
}

// BuiltinBackends is the list of built-in physical backends that can
// be used with NewBackend.
var BuiltinBackends = map[string]Factory{
	"inmem": func(map[string]string) (Backend, error) {
		return NewInmem(), nil
	},
	"consul": newConsulBackend,
	"file":   newFileBackend,
}
