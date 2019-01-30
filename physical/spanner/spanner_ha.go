package spanner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/spanner"
	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/physical"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

// Verify Backend satisfies the correct interfaces
var _ physical.HABackend = (*Backend)(nil)
var _ physical.Lock = (*Lock)(nil)

const (
	// LockRenewInterval is the time to wait between lock renewals.
	LockRenewInterval = 5 * time.Second

	// LockRetryInterval is the amount of time to wait if the lock fails before
	// trying again.
	LockRetryInterval = 5 * time.Second

	// LockTTL is the default lock TTL.
	LockTTL = 15 * time.Second

	// LockWatchRetryInterval is the amount of time to wait if a watch fails
	// before trying again.
	LockWatchRetryInterval = 5 * time.Second

	// LockWatchRetryMax is the number of times to retry a failed watch before
	// signaling that leadership is lost.
	LockWatchRetryMax = 5
)

var (
	// metricLockUnlock is the metric to register for a lock delete.
	metricLockUnlock = []string{"spanner", "lock", "unlock"}

	// metricLockGet is the metric to register for a lock get.
	metricLockLock = []string{"spanner", "lock", "lock"}

	// metricLockValue is the metric to register for a lock create/update.
	metricLockValue = []string{"spanner", "lock", "value"}
)

// Lock is the HA lock.
type Lock struct {
	// backend is the underlying physical backend.
	backend *Backend

	// key is the name of the key. value is the value of the key.
	key, value string

	// held is a boolean indicating if the lock is currently held.
	held bool

	// identity is the internal identity of this key (unique to this server
	// instance).
	identity string

	// lock is an internal lock
	lock sync.Mutex

	// stopCh is the channel that stops all operations. It may be closed in the
	// event of a leader loss or graceful shutdown. stopped is a boolean
	// indicating if we are stopped - it exists to prevent double closing the
	// channel. stopLock is a mutex around the locks.
	stopCh   chan struct{}
	stopped  bool
	stopLock sync.Mutex

	// Allow modifying the Lock durations for ease of unit testing.
	renewInterval      time.Duration
	retryInterval      time.Duration
	ttl                time.Duration
	watchRetryInterval time.Duration
	watchRetryMax      int
}

// LockRecord is the struct that corresponds to a lock.
type LockRecord struct {
	Key       string
	Value     string
	Identity  string
	Timestamp time.Time
}

// HAEnabled implements HABackend and indicates that this backend supports high
// availability.
func (b *Backend) HAEnabled() bool {
	return b.haEnabled
}

// LockWith acquires a mutual exclusion based on the given key.
func (b *Backend) LockWith(key, value string) (physical.Lock, error) {
	identity, err := uuid.GenerateUUID()
	if err != nil {
		return nil, errwrap.Wrapf("lock with: {{err}}", err)
	}
	return &Lock{
		backend:  b,
		key:      key,
		value:    value,
		identity: identity,
		stopped:  true,

		renewInterval:      LockRenewInterval,
		retryInterval:      LockRetryInterval,
		ttl:                LockTTL,
		watchRetryInterval: LockWatchRetryInterval,
		watchRetryMax:      LockWatchRetryMax,
	}, nil
}

// Lock acquires the given lock. The stopCh is optional. If closed, it
// interrupts the lock acquisition attempt. The returned channel should be
// closed when leadership is lost.
func (l *Lock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	defer metrics.MeasureSince(metricLockLock, time.Now())

	l.lock.Lock()
	defer l.lock.Unlock()
	if l.held {
		return nil, errors.New("lock already held")
	}

	// Attempt to lock - this function blocks until a lock is acquired or an error
	// occurs.
	acquired, err := l.attemptLock(stopCh)
	if err != nil {
		return nil, errwrap.Wrapf("lock: {{err}}", err)
	}
	if !acquired {
		return nil, nil
	}

	// We have the lock now
	l.held = true

	// Build the locks
	l.stopLock.Lock()
	l.stopCh = make(chan struct{})
	l.stopped = false
	l.stopLock.Unlock()

	// Periodically renew and watch the lock
	go l.renewLock()
	go l.watchLock()

	return l.stopCh, nil
}

// Unlock releases the lock.
func (l *Lock) Unlock() error {
	defer metrics.MeasureSince(metricLockUnlock, time.Now())

	l.lock.Lock()
	defer l.lock.Unlock()
	if !l.held {
		return nil
	}

	// Stop any existing locking or renewal attempts
	l.stopLock.Lock()
	if !l.stopped {
		l.stopped = true
		close(l.stopCh)
	}
	l.stopLock.Unlock()

	// Pooling
	l.backend.permitPool.Acquire()
	defer l.backend.permitPool.Release()

	// Delete
	ctx := context.Background()
	if _, err := l.backend.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		row, err := txn.ReadRow(ctx, l.backend.haTable, spanner.Key{l.key}, []string{"Identity"})
		if err != nil {
			if spanner.ErrCode(err) != codes.NotFound {
				return nil
			}
			return err
		}

		var r LockRecord
		if derr := row.ToStruct(&r); derr != nil {
			return errwrap.Wrapf("failed to decode to struct: {{err}}", derr)
		}

		// If the identity is different, that means that between the time that after
		// we stopped acquisition, the TTL expired and someone else grabbed the
		// lock. We do not want to delete a lock that is not our own.
		if r.Identity != l.identity {
			return nil
		}

		return txn.BufferWrite([]*spanner.Mutation{
			spanner.Delete(l.backend.haTable, spanner.Key{l.key}),
		})
	}); err != nil {
		return errwrap.Wrapf("unlock: {{err}}", err)
	}

	// We are no longer holding the lock
	l.held = false

	return nil
}

// Value returns the value of the lock and if it is held.
func (l *Lock) Value() (bool, string, error) {
	defer metrics.MeasureSince(metricLockValue, time.Now())

	r, err := l.get(context.Background())
	if err != nil {
		return false, "", err
	}
	if r == nil {
		return false, "", err
	}
	return true, string(r.Value), nil
}

// attemptLock attempts to acquire a lock. If the given channel is closed, the
// acquisition attempt stops. This function returns when a lock is acquired or
// an error occurs.
func (l *Lock) attemptLock(stopCh <-chan struct{}) (bool, error) {
	ticker := time.NewTicker(l.retryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			acquired, err := l.writeLock()
			if err != nil {
				return false, errwrap.Wrapf("attempt lock: {{err}}", err)
			}
			if !acquired {
				continue
			}

			return true, nil
		case <-stopCh:
			return false, nil
		}
	}
}

// renewLock renews the given lock until the channel is closed.
func (l *Lock) renewLock() {
	ticker := time.NewTicker(l.renewInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.writeLock()
		case <-l.stopCh:
			return
		}
	}
}

// watchLock checks whether the lock has changed in the table and closes the
// leader channel accordingly. If an error occurs during the check, watchLock
// will retry the operation and then close the leader channel if it can't
// succeed after retries.
func (l *Lock) watchLock() {
	retries := 0
	ticker := time.NewTicker(l.watchRetryInterval)

OUTER:
	for {
		// Check if the channel is already closed
		select {
		case <-l.stopCh:
			break OUTER
		default:
		}

		// Check if we've exceeded retries
		if retries >= l.watchRetryMax-1 {
			break OUTER
		}

		// Wait for the timer
		select {
		case <-ticker.C:
		case <-l.stopCh:
			break OUTER
		}

		// Attempt to read the key
		r, err := l.get(context.Background())
		if err != nil {
			retries++
			continue
		}

		// Verify the identity is the same
		if r == nil || r.Identity != l.identity {
			break OUTER
		}
	}

	l.stopLock.Lock()
	defer l.stopLock.Unlock()
	if !l.stopped {
		l.stopped = true
		close(l.stopCh)
	}
}

// writeLock writes the given lock using the following algorithm:
//
// - lock does not exist
//   - write the lock
// - lock exists
//   - if key is empty or identity is the same or timestamp exceeds TTL
//     - update the lock to self
func (l *Lock) writeLock() (bool, error) {
	// Pooling
	l.backend.permitPool.Acquire()
	defer l.backend.permitPool.Release()

	// Keep track of whether the lock was written
	lockWritten := false

	// Create a transaction to read and the update (maybe)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// The transaction will be retried, and it could sit in a queue behind, say,
	// the delete operation. To stop the transaction, we close the context when
	// the associated stopCh is received.
	go func() {
		select {
		case <-l.stopCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	_, err := l.backend.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		row, err := txn.ReadRow(ctx, l.backend.haTable, spanner.Key{l.key}, []string{"Key", "Identity", "Timestamp"})
		if err != nil && spanner.ErrCode(err) != codes.NotFound {
			return err
		}

		// If there was a record, verify that the record is still trustable.
		if row != nil {
			var r LockRecord
			if derr := row.ToStruct(&r); derr != nil {
				return errwrap.Wrapf("failed to decode to struct: {{err}}", derr)
			}

			// If the key is empty or the identity is ours or the ttl expired, we can
			// write. Otherwise, return now because we cannot.
			if r.Key != "" && r.Identity != l.identity && time.Now().UTC().Sub(r.Timestamp) < l.ttl {
				return nil
			}
		}

		m, err := spanner.InsertOrUpdateStruct(l.backend.haTable, &LockRecord{
			Key:       l.key,
			Value:     l.value,
			Identity:  l.identity,
			Timestamp: time.Now().UTC(),
		})
		if err != nil {
			return errwrap.Wrapf("failed to generate struct: {{err}}", err)
		}
		if err := txn.BufferWrite([]*spanner.Mutation{m}); err != nil {
			return errwrap.Wrapf("failed to write: {{err}}", err)
		}

		// Mark that the lock was acquired
		lockWritten = true

		return nil
	})
	if err != nil {
		return false, errwrap.Wrapf("write lock: {{err}}", err)
	}

	return lockWritten, nil
}

// get retrieves the value for the lock.
func (l *Lock) get(ctx context.Context) (*LockRecord, error) {
	// Pooling
	l.backend.permitPool.Acquire()
	defer l.backend.permitPool.Release()

	// Read
	row, err := l.backend.client.Single().ReadRow(ctx, l.backend.haTable, spanner.Key{l.key}, []string{"Key", "Value", "Timestamp", "Identity"})
	if spanner.ErrCode(err) == codes.NotFound {
		return nil, nil
	}
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to read value for %q: {{err}}", l.key), err)
	}

	var r LockRecord
	if err := row.ToStruct(&r); err != nil {
		return nil, errwrap.Wrapf("failed to decode lock: {{err}}", err)
	}
	return &r, nil
}
