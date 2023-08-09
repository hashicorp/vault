// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubernetes

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	sr "github.com/hashicorp/vault/serviceregistration"
	"github.com/hashicorp/vault/serviceregistration/kubernetes/client"
	"github.com/oklog/run"
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

	// initialStateSet determines whether an initial state has been set
	// successfully or whether a state already exists.
	initialStateSet bool

	// State stores an initial state to be set
	initialState sr.State

	// The map holds the path to the label being updated. It will only either
	// not hold a particular label, or hold _the last_ state we were aware of.
	// These should only be updated after initial state has been set.
	patchesToRetry map[string]*client.Patch

	// client is the Client to use when making API calls against kubernetes
	client *client.Client
}

// Run must be called for retries to be started.
func (r *retryHandler) Run(shutdownCh <-chan struct{}, wait *sync.WaitGroup) {
	r.setInitialState(shutdownCh)

	// Run this in a go func so this call doesn't block.
	wait.Add(1)
	go func() {
		// Make sure Vault will give us time to finish up here.
		defer wait.Done()

		var g run.Group

		// This run group watches for the shutdownCh
		shutdownActorStop := make(chan struct{})
		g.Add(func() error {
			select {
			case <-shutdownCh:
			case <-shutdownActorStop:
			}
			return nil
		}, func(error) {
			close(shutdownActorStop)
		})

		checkUpdateStateStop := make(chan struct{})
		g.Add(func() error {
			r.periodicUpdateState(checkUpdateStateStop)
			return nil
		}, func(error) {
			close(checkUpdateStateStop)
			r.client.Shutdown()
		})

		if err := g.Run(); err != nil {
			r.logger.Error("error encountered during periodic state update", "error", err)
		}
	}()
}

func (r *retryHandler) setInitialState(shutdownCh <-chan struct{}) {
	r.lock.Lock()
	defer r.lock.Unlock()

	doneCh := make(chan struct{})

	go func() {
		if err := r.setInitialStateInternal(); err != nil {
			if r.logger.IsWarn() {
				r.logger.Warn(fmt.Sprintf("unable to set initial state due to %s, will retry", err.Error()))
			}
		}
		close(doneCh)
	}()

	// Wait until the state is set or shutdown happens
	select {
	case <-doneCh:
	case <-shutdownCh:
	}
}

// Notify adds a patch to be retried until it's either completed without
// error, or no longer needed.
func (r *retryHandler) Notify(patch *client.Patch) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Initial state must be set first, or subsequent notifications we've
	// received could get smashed by a late-arriving initial state.
	// We will store this to retry it when appropriate.
	if !r.initialStateSet {
		if r.logger.IsWarn() {
			r.logger.Warn(fmt.Sprintf("cannot notify of present state for %s because initial state is unset", patch.Path))
		}
		r.patchesToRetry[patch.Path] = patch
		return
	}

	// Initial state has been sent, so it's OK to attempt a patch immediately.
	if err := r.client.PatchPod(r.namespace, r.podName, patch); err != nil {
		if r.logger.IsWarn() {
			r.logger.Warn(fmt.Sprintf("unable to update state for %s due to %s, will retry", patch.Path, err.Error()))
		}
		r.patchesToRetry[patch.Path] = patch
	}
}

// setInitialStateInternal sets the initial state remotely. This should be
// called with the lock held.
func (r *retryHandler) setInitialStateInternal() error {
	// If this is set, we return immediately
	if r.initialStateSet {
		return nil
	}

	// Verify that the pod exists and our configuration looks good.
	pod, err := r.client.GetPod(r.namespace, r.podName)
	if err != nil {
		return err
	}

	// Now to initially label our pod.
	if pod.Metadata == nil {
		// This should never happen IRL, just being defensive.
		return fmt.Errorf("no pod metadata on %+v", pod)
	}
	if pod.Metadata.Labels == nil {
		// Notify the labels field, and the labels as part of that one call.
		// The reason we must take a different approach to adding them is discussed here:
		// https://stackoverflow.com/questions/57480205/error-while-applying-json-patch-to-kubernetes-custom-resource
		if err := r.client.PatchPod(r.namespace, r.podName, &client.Patch{
			Operation: client.Add,
			Path:      "/metadata/labels",
			Value: map[string]string{
				labelVaultVersion: r.initialState.VaultVersion,
				labelActive:       strconv.FormatBool(r.initialState.IsActive),
				labelSealed:       strconv.FormatBool(r.initialState.IsSealed),
				labelPerfStandby:  strconv.FormatBool(r.initialState.IsPerformanceStandby),
				labelInitialized:  strconv.FormatBool(r.initialState.IsInitialized),
			},
		}); err != nil {
			return err
		}
	} else {
		// Create the labels through a patch to each individual field.
		patches := []*client.Patch{
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelVaultVersion,
				Value:     r.initialState.VaultVersion,
			},
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelActive,
				Value:     strconv.FormatBool(r.initialState.IsActive),
			},
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelSealed,
				Value:     strconv.FormatBool(r.initialState.IsSealed),
			},
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelPerfStandby,
				Value:     strconv.FormatBool(r.initialState.IsPerformanceStandby),
			},
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelInitialized,
				Value:     strconv.FormatBool(r.initialState.IsInitialized),
			},
		}
		if err := r.client.PatchPod(r.namespace, r.podName, patches...); err != nil {
			return err
		}
	}
	r.initialStateSet = true
	return nil
}

func (r *retryHandler) periodicUpdateState(stopCh chan struct{}) {
	retry := time.NewTicker(retryFreq)
	defer retry.Stop()

	for {
		// Call updateState immediately so we don't wait for the first tick
		// if setting the initial state
		r.updateState()

		select {
		case <-stopCh:
			return
		case <-retry.C:
		}
	}
}

func (r *retryHandler) updateState() {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Initial state must be set first, or subsequent notifications we've
	// received could get smashed by a late-arriving initial state.
	// If the state is already set, this is a no-op.
	if err := r.setInitialStateInternal(); err != nil {
		if r.logger.IsWarn() {
			r.logger.Warn(fmt.Sprintf("unable to set initial state due to %s, will retry", err.Error()))
		}
		// On failure, we leave the initial state func populated for
		// the next retry.
		return
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

	if err := r.client.PatchPod(r.namespace, r.podName, patches...); err != nil {
		if r.logger.IsWarn() {
			r.logger.Warn(fmt.Sprintf("unable to update state for due to %s, will retry", err.Error()))
		}
		return
	}
	r.patchesToRetry = make(map[string]*client.Patch)
}
