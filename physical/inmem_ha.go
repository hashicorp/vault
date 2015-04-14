package physical

import (
	"fmt"
	"sync"
)

type InmemHABackend struct {
	InmemBackend
	locks map[string]*sync.Mutex
	l     sync.Mutex
}

// NewInmemHA constructs a new in-memory HA backend. This is only for testing.
func NewInmemHA() *InmemHABackend {
	in := &InmemHABackend{
		InmemBackend: *NewInmem(),
		locks:        make(map[string]*sync.Mutex),
	}
	return in
}

// LockWith is used for mutual exclusion based on the given key.
func (i *InmemHABackend) LockWith(key string) (Lock, error) {
	i.l.Lock()
	defer i.l.Unlock()

	mutex, ok := i.locks[key]
	if !ok {
		mutex = new(sync.Mutex)
		i.locks[key] = mutex
	}

	return &InmemLock{mutex: mutex}, nil
}

// InmemLock is an in-memory Lock implementation for the HABackend
type InmemLock struct {
	// mutex is the underlying mutex, may be shared between
	// instances of InmemLock
	mutex *sync.Mutex

	held     bool
	leaderCh chan struct{}
	l        sync.Mutex
}

func (i *InmemLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	i.l.Lock()
	defer i.l.Unlock()
	if i.held {
		return nil, fmt.Errorf("lock already held")
	}

	// Attempt an async acquisition
	didLock := make(chan struct{})
	releaseCh := make(chan bool, 1)
	go func() {
		i.mutex.Lock()
		close(didLock)

		// Handle an early abort
		release := <-releaseCh
		if release {
			i.mutex.Unlock()
		}
	}()

	// Wait for lock acquisition or shutdown
	select {
	case <-didLock:
		releaseCh <- false
	case <-stopCh:
		releaseCh <- true
		return nil, nil
	}

	// Create the leader channel
	i.held = true
	i.leaderCh = make(chan struct{})
	return i.leaderCh, nil
}

func (i *InmemLock) Unlock() error {
	i.l.Lock()
	defer i.l.Unlock()

	if !i.held {
		return nil
	}

	close(i.leaderCh)
	i.leaderCh = nil
	i.held = false
	i.mutex.Unlock()
	return nil
}
