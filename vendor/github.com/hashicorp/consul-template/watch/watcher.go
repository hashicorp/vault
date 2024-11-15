// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package watch

import (
	"log"
	"sync"
	"time"

	dep "github.com/hashicorp/consul-template/dependency"
	"github.com/pkg/errors"
)

// dataBufferSize is the default number of views to process in a batch.
const dataBufferSize = 2048

type RetryFunc func(int) (bool, time.Duration)

// Watcher is a top-level manager for views that poll Consul for data.
type Watcher struct {
	sync.Mutex

	// clients is the collection of API clients to talk to upstreams.
	clients *dep.ClientSet

	// dataCh is the chan where Views will be published.
	dataCh chan *View

	// errCh is the chan where any errors will be published.
	errCh chan error

	// serverErrCh is the chan where response errors from the server will be published
	serverErrCh chan error

	// blockQueryWaitTime is amount of time in seconds to do a blocking query for
	blockQueryWaitTime time.Duration

	// failLookupErrors triggers error when a dependency Fetch fails to
	// return data after the first pass.
	failLookupErrors bool

	// depViewMap is a map of Templates to Views. Templates are keyed by
	// their string.
	depViewMap map[string]*View

	// maxStale specifies the maximum staleness of a query response.
	maxStale time.Duration

	// once signals if this watcher should tell views to retrieve data exactly
	// one time instead of polling infinitely.
	once bool

	// stopped signals if this watcher should stop adding any new dependencies
	stopped bool

	// retryFuncs specifies the different ways to retry based on the upstream.
	retryFuncConsul  RetryFunc
	retryFuncDefault RetryFunc
	retryFuncVault   RetryFunc
	retryFuncNomad   RetryFunc
}

type NewWatcherInput struct {
	// Clients is the client set to communicate with upstreams.
	Clients *dep.ClientSet

	// MaxStale is the maximum staleness of a query.
	MaxStale time.Duration

	// Once specifies this watcher should tell views to poll exactly once.
	Once bool

	// WaitTime is amount of time in seconds to do a blocking query for
	BlockQueryWaitTime time.Duration

	// FailLookupErrors triggers error when a dependency Fetch fails to
	// return data after the first pass.
	FailLookupErrors bool

	// RenewVault indicates if this watcher should renew Vault tokens.
	RenewVault bool

	// VaultToken is the vault token to renew.
	VaultToken string

	// VaultAgentTokenFile is the path to Vault Agent token file
	VaultAgentTokenFile string

	// RetryFuncs specify the different ways to retry based on the upstream.
	RetryFuncConsul  RetryFunc
	RetryFuncDefault RetryFunc
	RetryFuncVault   RetryFunc
	RetryFuncNomad   RetryFunc
}

// NewWatcher creates a new watcher using the given API client.
func NewWatcher(i *NewWatcherInput) *Watcher {
	w := &Watcher{
		clients:            i.Clients,
		depViewMap:         make(map[string]*View),
		dataCh:             make(chan *View, dataBufferSize),
		errCh:              make(chan error),
		serverErrCh:        make(chan error),
		maxStale:           i.MaxStale,
		once:               i.Once,
		blockQueryWaitTime: i.BlockQueryWaitTime,
		failLookupErrors:   i.FailLookupErrors,
		retryFuncConsul:    i.RetryFuncConsul,
		retryFuncDefault:   i.RetryFuncDefault,
		retryFuncVault:     i.RetryFuncVault,
		retryFuncNomad:     i.RetryFuncNomad,
	}
	return w
}

// DataCh returns a read-only channel of Views which is populated when a view
// receives data from its upstream.
func (w *Watcher) DataCh() <-chan *View {
	return w.dataCh
}

// ErrCh returns a read-only channel of errors returned by the upstream.
func (w *Watcher) ErrCh() <-chan error {
	if w == nil {
		return nil
	}
	return w.errCh
}

// ServerErrCh returns a read-only channel of errors returned by the server
// as a response to each consul-template instance
func (w *Watcher) ServerErrCh() <-chan error {
	if w == nil {
		return nil
	}
	return w.serverErrCh
}

// Add adds the given dependency to the list of monitored dependencies
// and start the associated view. If the dependency already exists, no action is
// taken.
//
// If the Dependency already existed or the watcher was concurrently stopped, it
// this function will return false. If the view was successfully created, it
// will return true. If an error occurs while creating the view, it will be
// returned here (but future errors returned by the view will happen on the
// channel).
func (w *Watcher) Add(d dep.Dependency) (bool, error) {
	w.Lock()
	defer w.Unlock()
	if w.stopped {
		log.Printf("[TRACE] (watcher) did not add %s because watcher is stopped", d)
		return false, nil
	}

	log.Printf("[DEBUG] (watcher) adding %s", d)

	if _, ok := w.depViewMap[d.String()]; ok {
		log.Printf("[TRACE] (watcher) %s already exists, skipping", d)
		return false, nil
	}

	// Choose the correct retry function based off of the dependency's type.
	var retryFunc RetryFunc
	switch d.Type() {
	case dep.TypeConsul:
		retryFunc = w.retryFuncConsul
	case dep.TypeVault:
		retryFunc = w.retryFuncVault
	case dep.TypeNomad:
		retryFunc = w.retryFuncNomad
	default:
		retryFunc = w.retryFuncDefault
	}

	v, err := NewView(&NewViewInput{
		Dependency:         d,
		Clients:            w.clients,
		MaxStale:           w.maxStale,
		BlockQueryWaitTime: w.blockQueryWaitTime,
		FailLookupErrors:   w.failLookupErrors,
		Once:               w.once,
		RetryFunc:          retryFunc,
	})
	if err != nil {
		return false, errors.Wrap(err, "watcher")
	}

	log.Printf("[TRACE] (watcher) %s starting", d)

	w.depViewMap[d.String()] = v
	go v.poll(w.dataCh, w.errCh, w.serverErrCh)

	return true, nil
}

// Watching determines if the given dependency is being watched.
func (w *Watcher) Watching(d dep.Dependency) bool {
	w.Lock()
	defer w.Unlock()

	_, ok := w.depViewMap[d.String()]
	return ok
}

// ForceWatching is used to force setting the internal state of watching
// a dependency. This is only used for unit testing purposes.
func (w *Watcher) ForceWatching(d dep.Dependency, enabled bool) {
	w.Lock()
	defer w.Unlock()

	if enabled {
		w.depViewMap[d.String()] = nil
	} else {
		delete(w.depViewMap, d.String())
	}
}

// Remove removes the given dependency from the list and stops the
// associated View. If a View for the given dependency does not exist, this
// function will return false. If the View does exist, this function will return
// true upon successful deletion.
func (w *Watcher) Remove(d dep.Dependency) bool {
	w.Lock()
	defer w.Unlock()

	log.Printf("[DEBUG] (watcher) removing %s", d)

	if view, ok := w.depViewMap[d.String()]; ok {
		log.Printf("[TRACE] (watcher) actually removing %s", d)
		view.stop()
		delete(w.depViewMap, d.String())
		return true
	}

	log.Printf("[TRACE] (watcher) %s did not exist, skipping", d)
	return false
}

// Size returns the number of views this watcher is watching.
func (w *Watcher) Size() int {
	w.Lock()
	defer w.Unlock()
	return len(w.depViewMap)
}

// Stop halts this watcher and any currently polling views immediately. If a
// view was in the middle of a poll, no data will be returned.
func (w *Watcher) Stop() {
	if w == nil {
		return
	}
	w.Lock()
	defer w.Unlock()

	log.Printf("[DEBUG] (watcher) stopping all views")

	for _, view := range w.depViewMap {
		if view == nil {
			continue
		}
		log.Printf("[TRACE] (watcher) stopping %s", view.Dependency())
		view.stop()
	}

	// Reset the map to have no views
	w.depViewMap = make(map[string]*View)

	// Close any idle TCP connections
	w.clients.Stop()
	w.stopped = true
}
