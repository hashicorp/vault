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
	client             *client.Client

	// This gets mutated on multiple threads.
	// The map holds the path to the label being updated. It will only either
	// not hold a particular label, or hold _the last_ state we were aware of.
	patchesToRetry     map[string]*client.Patch
	patchesToRetryLock sync.Mutex
}

// Run runs at an interval, checking if anything has failed and if so,
// attempting to send them again.
func (r *retryHandler) Run(shutdownCh <-chan struct{}, wait *sync.WaitGroup) {
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
			r.retry()
		}
	}
}

// Add adds a patch to be retried until it's either completed without
// error, or no longer needed.
func (r *retryHandler) Add(patch *client.Patch) {
	r.patchesToRetryLock.Lock()
	defer r.patchesToRetryLock.Unlock()
	r.patchesToRetry[patch.Path] = patch
}

func (r *retryHandler) retry() {
	r.patchesToRetryLock.Lock()
	defer r.patchesToRetryLock.Unlock()

	if len(r.patchesToRetry) == 0 {
		// Nothing to do here.
		return
	}

	patches := make([]*client.Patch, len(r.patchesToRetry))
	i := 0
	for _, patch := range r.patchesToRetry {
		patches[i] = patch
		i++
	}

	if err := r.client.PatchPod(r.namespace, r.podName, patches...); err != nil {
		if r.logger.IsWarn() {
			r.logger.Warn("unable to update state due to %s, will retry", err.Error())
		}
		return
	}
	r.patchesToRetry = make(map[string]*client.Patch)
}
