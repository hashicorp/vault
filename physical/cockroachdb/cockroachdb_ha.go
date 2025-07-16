// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cockroachdb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/physical"
)

const (
	// The lock TTL matches the default that Consul API uses, 15 seconds.
	// Used as part of SQL commands to set/extend lock expiry time relative to
	// database clock.
	CockroachDBLockTTLSeconds = 15

	// The amount of time to wait between the lock renewals
	CockroachDBLockRenewInterval = 5 * time.Second

	// CockroachDBLockRetryInterval is the amount of time to wait
	// if a lock fails before trying again.
	CockroachDBLockRetryInterval = time.Second
)

// Verify backend satisfies the correct interfaces.
var (
	_ physical.HABackend = (*CockroachDBBackend)(nil)
	_ physical.Lock      = (*CockroachDBLock)(nil)
)

type CockroachDBLock struct {
	backend  *CockroachDBBackend
	key      string
	value    string
	identity string
	lock     sync.Mutex

	renewTicker *time.Ticker

	// ttlSeconds is how long a lock is valid for.
	ttlSeconds int

	// renewInterval is how much time to wait between lock renewals.  must be << ttl.
	renewInterval time.Duration

	// retryInterval is how much time to wait between attempts to grab the lock.
	retryInterval time.Duration
}

func (c *CockroachDBBackend) HAEnabled() bool {
	return c.haEnabled
}

func (c *CockroachDBBackend) LockWith(key, value string) (physical.Lock, error) {
	identity, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	return &CockroachDBLock{
		backend:       c,
		key:           key,
		value:         value,
		identity:      identity,
		ttlSeconds:    CockroachDBLockTTLSeconds,
		renewInterval: CockroachDBLockRenewInterval,
		retryInterval: CockroachDBLockRetryInterval,
	}, nil
}

// Lock tries to acquire the lock by repeatedly trying to create a record in the
// CockroachDB table. It will block until either the stop channel is closed or
// the lock could be acquired successfully. The returned channel will be closed
// once the lock in the CockroachDB table cannot be renewed, either due to an
// error speaking to CockroachDB or because someone else has taken it.
func (l *CockroachDBLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	var (
		success = make(chan struct{})
		errors  = make(chan error, 1)
		leader  = make(chan struct{})
	)
	go l.tryToLock(stopCh, success, errors)

	select {
	case <-success:
		// After acquiring it successfully, we must renew the lock periodically.
		l.renewTicker = time.NewTicker(l.renewInterval)
		go l.periodicallyRenewLock(leader)
	case err := <-errors:
		return nil, err
	case <-stopCh:
		return nil, nil
	}

	return leader, nil
}

// Unlock releases the lock by deleting the lock record from the
// CockroachDB table.
func (l *CockroachDBLock) Unlock() error {
	c := l.backend
	if err := c.permitPool.Acquire(context.Background()); err != nil {
		return err
	}
	defer c.permitPool.Release()

	if l.renewTicker != nil {
		l.renewTicker.Stop()
	}

	_, err := c.haStatements["delete"].Exec(l.key)
	return err
}

// Value checks whether or not the lock is held by any instance of CockroachDBLock,
// including this one, and returns the current value.
func (l *CockroachDBLock) Value() (bool, string, error) {
	c := l.backend
	if err := c.permitPool.Acquire(context.Background()); err != nil {
		return false, "", err
	}
	defer c.permitPool.Release()
	var result string
	err := c.haStatements["get"].QueryRow(l.key).Scan(&result)

	switch err {
	case nil:
		return true, result, nil
	case sql.ErrNoRows:
		return false, "", nil
	default:
		return false, "", err

	}
}

// tryToLock tries to create a new item in CockroachDB every `retryInterval`.
// As long as the item cannot be created (because it already exists), it will
// be retried. If the operation fails due to an error, it is sent to the errors
// channel. When the lock could be acquired successfully, the success channel
// is closed.
func (l *CockroachDBLock) tryToLock(stop <-chan struct{}, success chan struct{}, errors chan error) {
	ticker := time.NewTicker(l.retryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			gotlock, err := l.writeItem()
			switch {
			case err != nil:
				// Send to the error channel and don't block if full.
				select {
				case errors <- err:
				default:
				}
				return
			case gotlock:
				close(success)
				return
			}
		}
	}
}

func (l *CockroachDBLock) periodicallyRenewLock(done chan struct{}) {
	for range l.renewTicker.C {
		gotlock, err := l.writeItem()
		if err != nil || !gotlock {
			close(done)
			l.renewTicker.Stop()
			return
		}
	}
}

// Attempts to put/update the CockroachDB item using condition expressions to
// evaluate the TTL.  Returns true if the lock was obtained, false if not.
// If false error may be nil or non-nil: nil indicates simply that someone
// else has the lock, whereas non-nil means that something unexpected happened.
func (l *CockroachDBLock) writeItem() (bool, error) {
	c := l.backend
	if err := c.permitPool.Acquire(context.Background()); err != nil {
		return false, err
	}
	defer c.permitPool.Release()

	sqlResult, err := c.haStatements["upsert"].Exec(l.identity, l.key, l.value, fmt.Sprintf("%d seconds", l.ttlSeconds))
	if err != nil {
		return false, err
	}
	if sqlResult == nil {
		return false, fmt.Errorf("empty SQL response received")
	}

	ar, err := sqlResult.RowsAffected()
	if err != nil {
		return false, err
	}
	return ar == 1, nil
}
