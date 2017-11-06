package physical

import (
	"strings"
	"sync"

	log "github.com/mgutz/logxi/v1"
)

const DefaultParallelOperations = 128

// The operation type
type Operation string

const (
	DeleteOperation Operation = "delete"
	GetOperation              = "get"
	ListOperation             = "list"
	PutOperation              = "put"
)

// ShutdownSignal
type ShutdownChannel chan struct{}

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

// HABackend is an extensions to the standard physical
// backend to support high-availability. Vault only expects to
// use mutual exclusion to allow multiple instances to act as a
// hot standby for a leader that services all requests.
type HABackend interface {
	// LockWith is used for mutual exclusion based on the given key.
	LockWith(key, value string) (Lock, error)

	// Whether or not HA functionality is enabled
	HAEnabled() bool
}

// Purgable is an optional interface for backends that support
// purging of their caches.
type Purgable interface {
	Purge()
}

// RedirectDetect is an optional interface that an HABackend
// can implement. If they do, a redirect address can be automatically
// detected.
type RedirectDetect interface {
	// DetectHostAddr is used to detect the host address
	DetectHostAddr() (string, error)
}

// Callback signatures for RunServiceDiscovery
type ActiveFunction func() bool
type SealedFunction func() bool

// ServiceDiscovery is an optional interface that an HABackend can implement.
// If they do, the state of a backend is advertised to the service discovery
// network.
type ServiceDiscovery interface {
	// NotifyActiveStateChange is used by Core to notify a backend
	// capable of ServiceDiscovery that this Vault instance has changed
	// its status to active or standby.
	NotifyActiveStateChange() error

	// NotifySealedStateChange is used by Core to notify a backend
	// capable of ServiceDiscovery that Vault has changed its Sealed
	// status to sealed or unsealed.
	NotifySealedStateChange() error

	// Run executes any background service discovery tasks until the
	// shutdown channel is closed.
	RunServiceDiscovery(waitGroup *sync.WaitGroup, shutdownCh ShutdownChannel, redirectAddr string, activeFunc ActiveFunction, sealedFunc SealedFunction) error
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
	Key      string
	Value    []byte
	SealWrap bool `json:"seal_wrap,omitempty"`
}

// Factory is the factory function to create a physical backend.
type Factory func(config map[string]string, logger log.Logger) (Backend, error)

// PermitPool is used to limit maximum outstanding requests
type PermitPool struct {
	sem chan int
}

// NewPermitPool returns a new permit pool with the provided
// number of permits
func NewPermitPool(permits int) *PermitPool {
	if permits < 1 {
		permits = DefaultParallelOperations
	}
	return &PermitPool{
		sem: make(chan int, permits),
	}
}

// Acquire returns when a permit has been acquired
func (c *PermitPool) Acquire() {
	c.sem <- 1
}

// Release returns a permit to the pool
func (c *PermitPool) Release() {
	<-c.sem
}

// Prefixes is a shared helper function returns all parent 'folders' for a
// given vault key.
// e.g. for 'foo/bar/baz', it returns ['foo', 'foo/bar']
func Prefixes(s string) []string {
	components := strings.Split(s, "/")
	result := []string{}
	for i := 1; i < len(components); i++ {
		result = append(result, strings.Join(components[:i], "/"))
	}
	return result
}
