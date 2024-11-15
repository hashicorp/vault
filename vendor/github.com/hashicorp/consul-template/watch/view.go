// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package watch

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"sync"
	"time"

	dep "github.com/hashicorp/consul-template/dependency"
)

var errLookup = fmt.Errorf("lookup error")

// View is a representation of a Dependency and the most recent data it has
// received from Consul.
type View struct {
	// dependency is the dependency that is associated with this View
	dependency dep.Dependency

	// clients is the list of clients to communicate upstream. This is passed
	// directly to the dependency.
	clients *dep.ClientSet

	// data is the most-recently-received data from Consul for this View. It is
	// accompanied by a series of locks and booleans to ensure consistency.
	dataLock     sync.RWMutex
	data         interface{}
	receivedData bool
	lastIndex    uint64

	// blockQueryWaitTime is amount of time in seconds to do a blocking query for
	blockQueryWaitTime time.Duration

	// maxStale is the maximum amount of time to allow a query to be stale.
	maxStale time.Duration

	// once determines if this view should receive data exactly once.
	once bool
	// failLookupErrors triggers error when a dependency Fetch fails to
	// return data after the first pass.
	failLookupErrors bool

	// retryFunc is the function to invoke on failure to determine if a retry
	// should be attempted.
	retryFunc RetryFunc

	// stopCh is used to stop polling on this View
	stopCh chan struct{}
}

// NewViewInput is used as input to the NewView function.
type NewViewInput struct {
	// Dependency is the dependency to associate with the new view.
	Dependency dep.Dependency

	// Clients is the list of clients to communicate upstream. This is passed
	// directly to the dependency.
	Clients *dep.ClientSet

	// BlockQueryWaitTime is amount of time in seconds to do a blocking query for
	BlockQueryWaitTime time.Duration

	// MaxStale is the maximum amount a time a query response is allowed to be
	// stale before forcing a read from the leader.
	MaxStale time.Duration

	// Once indicates this view should poll for data exactly one time.
	Once bool

	// FailLookupErrors triggers error when a dependency Fetch fails to
	// return data after the first pass.
	FailLookupErrors bool

	// RetryFunc is a function which dictates how this view should retry on
	// upstream errors.
	RetryFunc RetryFunc
}

// NewView constructs a new view with the given inputs.
func NewView(i *NewViewInput) (*View, error) {
	return &View{
		dependency:         i.Dependency,
		clients:            i.Clients,
		blockQueryWaitTime: i.BlockQueryWaitTime,
		maxStale:           i.MaxStale,
		once:               i.Once,
		failLookupErrors:   i.FailLookupErrors,
		retryFunc:          i.RetryFunc,
		stopCh:             make(chan struct{}, 1),
	}, nil
}

// Dependency returns the dependency attached to this View.
func (v *View) Dependency() dep.Dependency {
	return v.dependency
}

// Data returns the most-recently-received data from Consul for this View.
func (v *View) Data() interface{} {
	v.dataLock.RLock()
	defer v.dataLock.RUnlock()
	return v.data
}

// DataAndLastIndex returns the most-recently-received data from Consul for
// this view, along with the last index. This is atomic so you will get the
// index that goes with the data you are fetching.
func (v *View) DataAndLastIndex() (interface{}, uint64) {
	v.dataLock.RLock()
	defer v.dataLock.RUnlock()
	return v.data, v.lastIndex
}

// poll queries the Consul instance for data using the fetch function, but also
// accounts for interrupts on the interrupt channel. This allows the poll
// function to be fired in a goroutine, but then halted even if the fetch
// function is in the middle of a blocking query.
func (v *View) poll(viewCh chan<- *View, errCh chan<- error, serverErrCh chan<- error) {
	var retries int

	for {
		doneCh := make(chan struct{}, 1)
		successCh := make(chan struct{}, 1)
		fetchErrCh := make(chan error, 1)
		go v.fetch(doneCh, successCh, fetchErrCh)

	WAIT:
		select {
		case <-doneCh:
			// Reset the retry to avoid exponentially incrementing retries when we
			// have some successful requests
			retries = 0

			log.Printf("[TRACE] (view) %s received data", v.dependency)
			select {
			case <-v.stopCh:
				return
			case viewCh <- v:
			}

			// If we are operating in once mode, do not loop - we received data at
			// least once which is the API promise here.
			if v.once {
				return
			}
		case <-successCh:
			// We successfully received a non-error response from the server. This
			// does not mean we have data (that's dataCh's job), but rather this
			// just resets the counter indicating we communicated successfully. For
			// example, Consul make have an outage, but when it returns, the view
			// is unchanged. We have to reset the counter retries, but not update the
			// actual template.
			log.Printf("[TRACE] (view) %s successful contact, resetting retries", v.dependency)
			retries = 0
			goto WAIT
		case err := <-fetchErrCh:
			if !errors.Is(err, errLookup) && v.retryFunc != nil {
				retry, sleep := v.retryFunc(retries)
				serverErrCh <- err
				if retry {
					log.Printf("[WARN] (view) %s (retry attempt %d after %q)",
						err, retries+1, sleep)
					select {
					case <-time.After(sleep):
						retries++
						continue
					case <-v.stopCh:
						return
					}
				}
				log.Printf("[ERR] (view) %s (exceeded maximum retries)", err)
			}

			// Push the error back up to the watcher
			select {
			case <-v.stopCh:
				return
			case errCh <- err:
				return
			}
		case <-v.stopCh:
			log.Printf("[TRACE] (view) %s stopping poll (received on view stopCh)", v.dependency)
			return
		}
	}
}

// fetch queries the Consul instance for the attached dependency. This API
// promises that either data will be written to doneCh or an error will be
// written to errCh. It is designed to be run in a goroutine that selects the
// result of doneCh and errCh. It is assumed that only one instance of fetch
// is running per View and therefore no locking or mutexes are used.
func (v *View) fetch(doneCh, successCh chan<- struct{}, errCh chan<- error) {
	log.Printf("[TRACE] (view) %s starting fetch", v.dependency)

	var allowStale bool
	if v.maxStale != 0 {
		allowStale = true
	}

	firstLoop := true // to disable rate limiting on first pass
	for {
		// If the view was stopped, short-circuit this loop. This prevents a bug
		// where a view can get "lost" in the event Consul Template is reloaded.
		select {
		case <-v.stopCh:
			return
		default:
		}

		start := time.Now() // for rateLimiter below

		data, rm, err := v.dependency.Fetch(v.clients, &dep.QueryOptions{
			AllowStale: allowStale,
			WaitTime:   v.blockQueryWaitTime,
			WaitIndex:  v.lastIndex,
		})
		if err != nil {
			if err == dep.ErrStopped {
				log.Printf("[TRACE] (view) %s reported stop", v.dependency)
			} else {
				errCh <- err
			}
			return
		}

		if rm == nil {
			errCh <- fmt.Errorf("received nil response metadata - this is a bug " +
				"and should be reported")
			return
		}

		// If we got this far, we received data successfully. That data might not
		// trigger a data update (because we could continue below), but we need to
		// inform the poller to reset the retry count.
		log.Printf("[TRACE] (view) %s marking successful data response", v.dependency)
		select {
		case successCh <- struct{}{}:
		default:
		}

		if allowStale && rm.LastContact > v.maxStale {
			allowStale = false
			log.Printf("[TRACE] (view) %s stale data (last contact exceeded max_stale)", v.dependency)
			continue
		}

		if v.maxStale != 0 {
			allowStale = true
		}

		if v.failLookupErrors && !firstLoop && !v.receivedData {
			errCh <- errLookup
			return
		}

		if dur := rateLimiter(start); dur > 1 && !firstLoop {
			time.Sleep(dur)
		}
		firstLoop = false

		// blocking queries that return due to block timeout
		// will have the same index
		if rm.LastIndex == v.lastIndex {
			log.Printf("[TRACE] (view) %s no new data (index was the same)", v.dependency)
			continue
		}

		v.dataLock.Lock()
		if rm.LastIndex < v.lastIndex {
			log.Printf("[TRACE] (view) %s had a lower index, resetting", v.dependency)
			v.lastIndex = 0
			v.dataLock.Unlock()
			continue
		}
		v.lastIndex = rm.LastIndex

		if v.receivedData && reflect.DeepEqual(data, v.data) {
			log.Printf("[TRACE] (view) %s no new data (contents were the same)", v.dependency)
			v.dataLock.Unlock()
			continue
		}

		// this is for queries that are blocking but return a nil value on
		// lookup failures, but you want the dependency to act like it is still
		// blocking and loop back and hit it again.
		if data == nil && rm.BlockOnNil {
			log.Printf("[TRACE] (view) %s asked for blocking query", v.dependency)
			v.dataLock.Unlock()
			continue
		}

		v.data = data
		v.receivedData = true
		v.dataLock.Unlock()

		close(doneCh)
		return
	}
}

const minDelayBetweenUpdates = time.Millisecond * 100

// return a duration to sleep to limit the frequency of upstream calls
func rateLimiter(start time.Time) time.Duration {
	remaining := minDelayBetweenUpdates - time.Since(start)
	if remaining > 0 {
		dither := time.Duration(rand.Int63n(20000000)) // 0-20ms
		return remaining + dither
	}
	return 0
}

// stop halts polling of this view.
func (v *View) stop() {
	v.dependency.Stop()
	close(v.stopCh)
}
