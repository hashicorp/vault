package inmem

import (
	"fmt"
	"sync"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
)

type InmemHABackend struct {
	physical.Backend
	locks  map[string]string
	l      *sync.Mutex
	cond   *sync.Cond
	logger log.Logger
}

type TransactionalInmemHABackend struct {
	physical.Transactional
	InmemHABackend
}

// NewInmemHA constructs a new in-memory HA backend. This is only for testing.
func NewInmemHA(_ map[string]string, logger log.Logger) (physical.Backend, error) {
	be, err := NewInmem(nil, logger)
	if err != nil {
		return nil, err
	}

	in := &InmemHABackend{
		Backend: be,
		locks:   make(map[string]string),
		logger:  logger,
		l:       new(sync.Mutex),
	}
	in.cond = sync.NewCond(in.l)
	return in, nil
}

func NewTransactionalInmemHA(_ map[string]string, logger log.Logger) (physical.Backend, error) {
	transInmem, err := NewTransactionalInmem(nil, logger)
	if err != nil {
		return nil, err
	}
	inmemHA := InmemHABackend{
		Backend: transInmem,
		locks:   make(map[string]string),
		logger:  logger,
		l:       new(sync.Mutex),
	}

	in := &TransactionalInmemHABackend{
		InmemHABackend: inmemHA,
		Transactional:  transInmem.(physical.Transactional),
	}
	in.cond = sync.NewCond(in.l)
	return in, nil
}

// LockWith is used for mutual exclusion based on the given key.
func (i *InmemHABackend) LockWith(key, value string) (physical.Lock, error) {
	l := &InmemLock{
		in:    i,
		key:   key,
		value: value,
	}
	return l, nil
}

// LockMapSize is used in some tests to determine whether this backend has ever
// been used for HA purposes rather than simply for storage
func (i *InmemHABackend) LockMapSize() int {
	return len(i.locks)
}

// HAEnabled indicates whether the HA functionality should be exposed.
// Currently always returns true.
func (i *InmemHABackend) HAEnabled() bool {
	return true
}

// InmemLock is an in-memory Lock implementation for the HABackend
type InmemLock struct {
	in    *InmemHABackend
	key   string
	value string

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
		// Wait to acquire the lock
		i.in.l.Lock()
		_, ok := i.in.locks[i.key]
		for ok {
			i.in.cond.Wait()
			_, ok = i.in.locks[i.key]
		}
		i.in.locks[i.key] = i.value
		i.in.l.Unlock()

		// Signal that lock is held
		close(didLock)

		// Handle an early abort
		release := <-releaseCh
		if release {
			i.in.l.Lock()
			delete(i.in.locks, i.key)
			i.in.l.Unlock()
			i.in.cond.Broadcast()
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

	i.in.l.Lock()
	delete(i.in.locks, i.key)
	i.in.l.Unlock()
	i.in.cond.Broadcast()
	return nil
}

func (i *InmemLock) Value() (bool, string, error) {
	i.in.l.Lock()
	val, ok := i.in.locks[i.key]
	i.in.l.Unlock()
	return ok, val, nil
}
