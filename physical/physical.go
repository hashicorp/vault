package physical

import "fmt"

// Backend is the interface required for a physical
// backend. A physical backend is used to durably store
// data outside of Vault. As such, it is completely untrusted,
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

// HABackend is an extentions to the standard physical
// backend to support high-availability. Vault only expects to
// use mutual exclusion to allow multiple instances to act as a
// hot standby for a leader that services all requests.
type HABackend interface {
	// LockWith is used for mutual exclusion based on the given key.
	LockWith(key, value string) (Lock, error)
}

// AdvertiseDetect is an optional interface that an HABackend
// can implement. If they do, an advertise address can be automatically
// detected.
type AdvertiseDetect interface {
	// DetectHostAddr is used to detect the host address
	DetectHostAddr() (string, error)
}

type Lock interface {
	// Lock is used to acquire the given lock
	// The stopCh is optional and if closed should interrupt the lock
	// acquisition attempt. The return struct should be closed when
	// leadership is lost.
	Lock(stopCh <-chan struct{}) (<-chan struct{}, error)

	// Unlock is used to release the lock
	Unlock() error

	// Returns the value of the lock and if it is held
	Value() (bool, string, error)
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
	"consul":    newConsulBackend,
	"zookeeper": newZookeeperBackend,
	"file":      newFileBackend,
	"s3":        newS3Backend,
	"etcd":      newEtcdBackend,
	"mysql":     newMySQLBackend,
}
