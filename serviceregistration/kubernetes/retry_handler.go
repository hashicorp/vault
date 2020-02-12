package kubernetes

import (
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/serviceregistration/kubernetes/client"
)

// How often to retry sending a state update if it fails.
var retryFreq = 5 * time.Second

// retryHandler executes retries.
// It is thread-safe.
type retryHandler struct {
	// These don't need a mutex because they're never mutated.
	logger             hclog.Logger
	namespace, podName string

	// To synchronize setInitialState and patchesToRetry.
	lock sync.Mutex

	// setInitialState will be nil if this has been done successfully.
	// It must be done before any patches are retried.
	setInitialState func() error

	// The map holds the path to the label being updated. It will only either
	// not hold a particular label, or hold _the last_ state we were aware of.
	// These should only be updated after initial state has been set.
	patchesToRetry map[string]*client.Patch
}

func (r *retryHandler) SetInitialState(setInitialState func() error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if err := setInitialState(); err != nil {
		if r.logger.IsWarn() {
			r.logger.Warn("unable to set initial state due to %s, will retry", err.Error())
		}
		r.setInitialState = setInitialState
	}
}

// Run must be called for retries to be started.
func (r *retryHandler) Run(shutdownCh <-chan struct{}, wait *sync.WaitGroup, c *client.Client) {
	// Run this in a go func so this call doesn't block.
	go func() {
		// Make sure Vault will give us time to finish up here.
		wait.Add(1)
		defer wait.Done()

		retry := time.NewTicker(retryFreq)
		defer retry.Stop()
		for {
			select {
			case <-shutdownCh:
				return
			case <-retry.C:
				r.retry(c)
			}
		}
	}()
}

// Notify adds a patch to be retried until it's either completed without
// error, or no longer needed.
func (r *retryHandler) Notify(c *client.Client, patch *client.Patch) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Initial state must be set first, or subsequent notifications we've
	// received could get smashed by a late-arriving initial state.
	// We will store this to retry it when appropriate.
	if r.setInitialState != nil {
		if r.logger.IsWarn() {
			r.logger.Warn("cannot notify of present state because initial state is unset", patch.Path)
		}
		r.patchesToRetry[patch.Path] = patch
		return
	}

	// Initial state has been sent, so it's OK to attempt a patch immediately.
	if err := c.PatchPod(r.namespace, r.podName, patch); err != nil {
		if r.logger.IsWarn() {
			r.logger.Warn("unable to update state due to %s, will retry", patch.Path, err.Error())
		}
		r.patchesToRetry[patch.Path] = patch
	}
}

func (r *retryHandler) retry(c *client.Client) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Initial state must be set first, or subsequent notifications we've
	// received could get smashed by a late-arriving initial state.
	if r.setInitialState != nil {
		if err := r.setInitialState(); err != nil {
			if r.logger.IsWarn() {
				r.logger.Warn("unable to set initial state due to %s, will retry", err.Error())
			}
			// On failure, we leave the initial state func populated for
			// the next retry.
			return
		}
		// On success, we set it to nil and allow the logic to continue.
		r.setInitialState = nil
	}

	if len(r.patchesToRetry) == 0 {
		// Nothing further to do here.
		return
	}

	patches := make([]*client.Patch, len(r.patchesToRetry))
	i := 0
	for _, patch := range r.patchesToRetry {
		patches[i] = patch
		i++
	}

	if err := c.PatchPod(r.namespace, r.podName, patches...); err != nil {
		if r.logger.IsWarn() {
			r.logger.Warn("unable to update state due to %s, will retry", err.Error())
		}
		return
	}
	r.patchesToRetry = make(map[string]*client.Patch)
}
