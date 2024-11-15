// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/hashicorp/go-multierror"
)

const (
	lockLeaseRenewalFactor = 0.7
	lockRetryBackoffFactor = 1.1

	// DefaultLockTTL is the default value used to maintain a lock before it needs to
	// be renewed. The actual value comes from the experience with Consul.
	DefaultLockTTL = 15 * time.Second

	// DefaultLockDelay is the default a lock will be blocked after the TTL
	// went by without any renews. It is intended to prevent split brain situations.
	// The actual value comes from the experience with Consul.
	DefaultLockDelay = 15 * time.Second
)

var (
	// ErrLockConflict is returned in case a lock operation can't be performed
	// because the caller is not the current holder of the lock.
	ErrLockConflict = errors.New("conflicting operation over lock")

	//LockNoPathErr is returned when no path is provided in the variable to be
	// used for the lease mechanism
	LockNoPathErr = errors.New("variable's path can't be empty")
)

// Locks returns a new handle on a lock for the given variable.
func (c *Client) Locks(wo WriteOptions, v Variable, opts ...LocksOption) (*Locks, error) {

	if v.Path == "" {
		return nil, LockNoPathErr
	}

	ttl, err := time.ParseDuration(v.Lock.TTL)
	if err != nil {
		return nil, err
	}

	l := &Locks{
		c:            c,
		WriteOptions: wo,
		variable:     v,
		ttl:          ttl,
		ro: retryOptions{
			maxToLastCall: ttl,
			maxRetries:    defaultNumberOfRetries,
		},
	}

	for _, opt := range opts {
		opt(l)
	}

	l.c.configureRetries(&l.ro)

	return l, nil
}

// Locks is used to maintain all the resources necessary to operate over a lock.
// It makes the calls to the http using an exponential retry mechanism that will
// try until it either reaches 5 attempts or the ttl of the lock expires.
// The variable doesn't need to exist, one will be created internally
// but a path most be provided.
//
// Important: It will be on the user to remove the variable created for the lock.
type Locks struct {
	c        *Client
	variable Variable
	ttl      time.Duration
	ro       retryOptions

	WriteOptions
}

type LocksOption = func(l *Locks)

// LocksOptionWithMaxRetries allows access to configure the number of max retries the lock
// handler will perform in case of an expected response while interacting with the
// locks endpoint.
func LocksOptionWithMaxRetries(maxRetries int64) LocksOption {
	return func(l *Locks) {
		l.ro.maxRetries = maxRetries
	}
}

//	Acquire will make the actual call to acquire the lock over the variable using
//	the ttl in the Locks to create the VariableLock. It will return the
//	path of the variable holding the lock.
//
// Acquire returns the path to the variable holding the lock.
func (l *Locks) Acquire(ctx context.Context) (string, error) {

	var out Variable

	_, err := l.c.retryPut(ctx, "/v1/var/"+l.variable.Path+"?lock-acquire", l.variable, &out, &l.WriteOptions)
	if err != nil {
		callErr, ok := err.(UnexpectedResponseError)

		// http.StatusConflict means the lock is already held. This will happen
		// under the normal execution if multiple instances are fighting for the same lock and
		// doesn't disrupt the flow.
		if ok && callErr.statusCode == http.StatusConflict {
			return "", fmt.Errorf("acquire conflict %w", ErrLockConflict)
		}

		return "", err
	}

	l.variable.Lock = out.Lock

	return l.variable.Path, nil
}

// Release makes the call to release the lock over a variable, even if the ttl
// has not yet passed.
// In case of a call to release a non held lock, Release returns ErrLockConflict.
func (l *Locks) Release(ctx context.Context) error {
	var out Variable

	rv := &Variable{
		Lock: &VariableLock{
			ID: l.variable.LockID(),
		},
	}

	_, err := l.c.retryPut(ctx, "/v1/var/"+l.variable.Path+"?lock-release", rv,
		&out, &l.WriteOptions)
	if err != nil {
		callErr, ok := err.(UnexpectedResponseError)

		if ok && callErr.statusCode == http.StatusConflict {
			return fmt.Errorf("release conflict %w", ErrLockConflict)
		}
		return err
	}

	return nil
}

// Renew is used to extend the ttl of a lock. It can be used as a heartbeat or a
// lease to maintain the hold over the lock for longer periods or as a sync
// mechanism among multiple instances looking to acquire the same lock.
// Renew will return true if the renewal was successful.
//
// In case of a call to renew a non held lock, Renew returns ErrLockConflict.
func (l *Locks) Renew(ctx context.Context) error {
	var out VariableMetadata

	_, err := l.c.retryPut(ctx, "/v1/var/"+l.variable.Path+"?lock-renew", l.variable, &out, &l.WriteOptions)
	if err != nil {
		callErr, ok := err.(UnexpectedResponseError)

		if ok && callErr.statusCode == http.StatusConflict {
			return fmt.Errorf("renew conflict %w", ErrLockConflict)
		}

		return err
	}
	return nil
}

func (l *Locks) LockTTL() time.Duration {
	return l.ttl
}

// Locker is the interface that wraps the lock handler. It is used by the lock
// leaser to handle all lock operations.
type Locker interface {
	// Acquire will make the actual call to acquire the lock over the variable using
	// the ttl in the Locks to create the VariableLock.
	//
	// Acquire returns the path to the variable holding the lock.
	Acquire(ctx context.Context) (string, error)
	// Release makes the call to release the lock over a variable, even if the ttl
	// has not yet passed.
	Release(ctx context.Context) error
	// Renew is used to extend the ttl of a lock. It can be used as a heartbeat or a
	// lease to maintain the hold over the lock for longer periods or as a sync
	// mechanism among multiple instances looking to acquire the same lock.
	Renew(ctx context.Context) error

	// LockTTL returns the expiration time of the underlying lock.
	LockTTL() time.Duration
}

// LockLeaser is a helper used to run a protected function that should only be
// active if the instance that runs it is currently holding the lock.
// Can be used to provide synchrony among multiple independent instances.
//
// It includes the lease renewal mechanism and tracking in case the protected
// function returns an error. Internally it uses an exponential retry mechanism
// for the api calls.
type LockLeaser struct {
	Name          string
	renewalPeriod time.Duration
	waitPeriod    time.Duration
	randomDelay   time.Duration
	earlyReturn   bool
	locked        bool

	locker Locker
}

type LockLeaserOption = func(l *LockLeaser)

// LockLeaserOptionWithEarlyReturn informs the leaser to return after the lock
// acquire fails and to not wait to attempt again.
func LockLeaserOptionWithEarlyReturn(er bool) LockLeaserOption {
	return func(l *LockLeaser) {
		l.earlyReturn = er
	}
}

// LockLeaserOptionWithWaitPeriod is used to set a back off period between
// calls to attempt to acquire the lock. By default it is set to 1.1 * TTLs.
func LockLeaserOptionWithWaitPeriod(wp time.Duration) LockLeaserOption {
	return func(l *LockLeaser) {
		l.waitPeriod = wp
	}
}

// NewLockLeaser returns an instance of LockLeaser. callerID
// is optional, in case they it is not provided, internal one will be created.
func (c *Client) NewLockLeaser(l Locker, opts ...LockLeaserOption) *LockLeaser {

	rn := rand.New(rand.NewSource(time.Now().Unix())).Intn(100)

	ll := &LockLeaser{
		renewalPeriod: time.Duration(float64(l.LockTTL()) * lockLeaseRenewalFactor),
		waitPeriod:    time.Duration(float64(l.LockTTL()) * lockRetryBackoffFactor),
		randomDelay:   time.Duration(rn) * time.Millisecond,
		locker:        l,
		earlyReturn:   false,
	}

	for _, opt := range opts {
		opt(ll)
	}

	return ll
}

// Start wraps the start function in charge of executing the protected
// function and maintain the lease but is in charge of releasing the
// lock before exiting. It is a blocking function.
func (ll *LockLeaser) Start(ctx context.Context, protectedFuncs ...func(ctx context.Context) error) error {
	var mErr multierror.Error

	err := ll.start(ctx, protectedFuncs...)
	if err != nil {
		mErr.Errors = append(mErr.Errors, err)
	}

	if ll.locked {
		err = ll.locker.Release(ctx)
		if err != nil {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("lock release: %w", err))
		}
	}

	return mErr.ErrorOrNil()
}

// start starts the process of maintaining the lease and executes the protected
// function on an independent go routine. It is a blocking function, it
// will return once the protected function is done or an execution error
// arises.
func (ll *LockLeaser) start(ctx context.Context, protectedFuncs ...func(ctx context.Context) error) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// errChannel is used track execution errors
	errChannel := make(chan error, 1)
	defer close(errChannel)

	// To avoid collisions if all the instances start at the same time, wait
	// a random time before making the first call.
	waitWithContext(ctx, ll.randomDelay)

	waitTicker := time.NewTicker(ll.waitPeriod)
	defer waitTicker.Stop()

	for {
		lockID, err := ll.locker.Acquire(ctx)
		if err != nil {

			if errors.Is(err, ErrLockConflict) && ll.earlyReturn {

				return nil
			}

			if !errors.Is(err, ErrLockConflict) {
				errChannel <- err
			}
		}

		if lockID != "" {
			ll.locked = true

			funcCtx, funcCancel := context.WithCancel(ctx)
			defer funcCancel()

			// Execute the lock protected function.
			go func() {
				defer funcCancel()
				for _, f := range protectedFuncs {
					err := f(funcCtx)
					if err != nil {
						errChannel <- fmt.Errorf("error executing protected function %w", err)
						return
					}
					cancel()
				}
			}()

			// Maintain lease is a blocking function, it will return if there is
			// an error maintaining the lease or the protected function returned.
			err = ll.maintainLease(funcCtx)
			if err != nil && !errors.Is(err, ErrLockConflict) {
				errChannel <- fmt.Errorf("error renewing the lease: %w", err)
			}
		}

		waitTicker.Stop()
		waitTicker = time.NewTicker(ll.waitPeriod)
		select {
		case <-ctx.Done():
			return nil

		case err := <-errChannel:
			return fmt.Errorf("locks: %w", err)

		case <-waitTicker.C:
		}
	}
}

func (ll *LockLeaser) maintainLease(ctx context.Context) error {
	renewTicker := time.NewTicker(ll.renewalPeriod)
	defer renewTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-renewTicker.C:
			err := ll.locker.Renew(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func waitWithContext(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
	case <-t.C:
	}
}
