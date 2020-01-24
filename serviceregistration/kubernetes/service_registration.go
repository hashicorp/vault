package kubernetes

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	sr "github.com/hashicorp/vault/serviceregistration"
	"github.com/hashicorp/vault/serviceregistration/kubernetes/client"
)

const (
	// Labels are placed in a pod's metadata.
	labelVaultVersion = "vault-version"
	labelActive       = "vault-ha-active"
	labelSealed       = "vault-ha-sealed"
	labelPerfStandby  = "vault-ha-perf-standby"
	labelInitialized  = "vault-ha-initialized"

	// This is the path to where these labels are applied.
	pathToLabels = "/metadata/labels/"

	// How often to retry sending a state update if it fails.
	retryFreq = 5 * time.Second
)

func NewServiceRegistration(config map[string]string, logger log.Logger, state *sr.State, _ string) (sr.ServiceRegistration, error) {
	c, err := client.New(logger)
	if err != nil {
		return nil, err
	}

	namespace := ""
	switch {
	case os.Getenv(client.EnvVarKubernetesNamespace) != "":
		namespace = os.Getenv(client.EnvVarKubernetesNamespace)
	case config["namespace"] != "":
		namespace = config["namespace"]
	default:
		return nil, fmt.Errorf(`namespace must be provided via %q or the "namespace" config parameter`, client.EnvVarKubernetesNamespace)
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("namespace: %q", namespace))
	}

	podName := ""
	switch {
	case os.Getenv(client.EnvVarKubernetesPodName) != "":
		podName = os.Getenv(client.EnvVarKubernetesPodName)
	case config["pod_name"] != "":
		podName = config["pod_name"]
	default:
		return nil, fmt.Errorf(`pod name must be provided via %q or the "pod_name" config parameter`, client.EnvVarKubernetesPodName)
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("pod name: %q", podName))
	}

	// Verify that the pod exists and our configuration looks good.
	pod, err := c.GetPod(namespace, podName)
	if err != nil {
		return nil, err
	}

	// Now to initially label our pod.
	if pod.Metadata == nil {
		// This should never happen because it's required to add a pod
		// name to the metadata, and kubernetes adds some more as well,
		// just being defensive.
		return nil, fmt.Errorf("no pod metadata on %+v", pod)
	}
	if pod.Metadata.Labels == nil {
		// If this Kube pod doesn't already have a labels field, we won't
		// be able to add them. This is discussed here:
		// https://stackoverflow.com/questions/57480205/error-while-applying-json-patch-to-kubernetes-custom-resource
		// Create the labels as part of adding the labels field.
		if err := c.PatchPod(namespace, podName, &client.Patch{
			Operation: client.Add,
			Path:      "/metadata/labels",
			Value: map[string]string{
				labelVaultVersion: state.VaultVersion,
				labelActive:       toString(state.IsActive),
				labelSealed:       toString(state.IsSealed),
				labelPerfStandby:  toString(state.IsPerformanceStandby),
				labelInitialized:  toString(state.IsInitialized),
			},
		}); err != nil {
			return nil, err
		}
	} else {
		// Create the labels through a patch to each field.
		patches := []*client.Patch{
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelVaultVersion,
				Value:     state.VaultVersion,
			},
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelActive,
				Value:     toString(state.IsActive),
			},
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelSealed,
				Value:     toString(state.IsSealed),
			},
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelPerfStandby,
				Value:     toString(state.IsPerformanceStandby),
			},
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelInitialized,
				Value:     toString(state.IsInitialized),
			},
		}
		if err := c.PatchPod(namespace, podName, patches...); err != nil {
			return nil, err
		}
	}

	// Construct a registration to receive ongoing state updates.
	registration := &serviceRegistration{
		logger:    logger,
		namespace: namespace,
		podName:   podName,
		client:    c,
		retryer: &retryer{
			logger:    logger,
			namespace: namespace,
			podName:   podName,
			client:    c,
		},
	}
	return registration, nil
}

type serviceRegistration struct {
	logger             log.Logger
	namespace, podName string
	client             *client.Client
	retryer            *retryer
}

func (r *serviceRegistration) Run(shutdownCh <-chan struct{}, wait *sync.WaitGroup) error {
	// Run a background goroutine to leave labels in the final state we'd like
	// when Vault shuts down.
	go r.onShutdown(shutdownCh, wait)

	// Run a service that retries errored-out notifications if they occur.
	go r.retryer.Run(shutdownCh, wait)

	return nil
}

func (r *serviceRegistration) NotifyActiveStateChange(isActive bool) error {
	patch := &client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelActive,
		Value:     toString(isActive),
	}
	return r.notifyOrRetry(patch)
}

func (r *serviceRegistration) NotifySealedStateChange(isSealed bool) error {
	patch := &client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelSealed,
		Value:     toString(isSealed),
	}
	return r.notifyOrRetry(patch)
}

func (r *serviceRegistration) NotifyPerformanceStandbyStateChange(isStandby bool) error {
	patch := &client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelPerfStandby,
		Value:     toString(isStandby),
	}
	return r.notifyOrRetry(patch)
}

func (r *serviceRegistration) NotifyInitializedStateChange(isInitialized bool) error {
	patch := &client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelInitialized,
		Value:     toString(isInitialized),
	}
	return r.notifyOrRetry(patch)
}

func (r *serviceRegistration) onShutdown(shutdownCh <-chan struct{}, wait *sync.WaitGroup) {
	// Ensure Vault will allow us time to finish this code.
	wait.Add(1)
	defer wait.Done()

	// Start running this when we receive a shutdown.
	<-shutdownCh

	// Label the pod with the values we want to leave behind after shutdown.
	patches := []*client.Patch{
		{
			Operation: client.Replace,
			Path:      pathToLabels + labelActive,
			Value:     toString(false),
		},
		{
			Operation: client.Replace,
			Path:      pathToLabels + labelSealed,
			Value:     toString(true),
		},
		{
			Operation: client.Replace,
			Path:      pathToLabels + labelPerfStandby,
			Value:     toString(false),
		},
		{
			Operation: client.Replace,
			Path:      pathToLabels + labelInitialized,
			Value:     toString(false),
		},
	}
	if err := r.client.PatchPod(r.namespace, r.podName, patches...); err != nil {
		if r.logger.IsError() {
			r.logger.Error(fmt.Sprintf("unable to set final status on pod name %q in namespace %q on shutdown: %s", r.podName, r.namespace, err))
		}
		return
	}
}

func (r *serviceRegistration) notifyOrRetry(patch *client.Patch) error {
	if err := r.client.PatchPod(r.namespace, r.podName, patch); err != nil {
		if r.logger.IsWarn() {
			r.logger.Warn("unable to update state due to %s, will retry", err.Error())
		}
		if err := r.retryer.Add(); err != nil {
			return err
		}
	}
	return nil
}

type retryer struct {
	logger             log.Logger
	namespace, podName string
	client             *client.Client

	// To be populated by the Add method,
	// and subtracted from by successfully applying
	// them in Run.
	desiredPatches     map[string]*client.Patch
	desiredPatchesLock sync.Mutex
}

func (r *retryer) Run(shutdownCh <-chan struct{}, wait *sync.WaitGroup) {
	wait.Add(1)
	defer wait.Done()

	retry := time.NewTicker(retryFreq)
	for {
		select {
		case <-shutdownCh:
			return
		case <-retry.C:
			r.attemptDesiredPatches()
		}
	}
}

func (r *retryer) Add(patches ...*client.Patch) error {
	r.desiredPatchesLock.Lock()
	defer r.desiredPatchesLock.Unlock()

	for _, newPatch := range patches {
		operationAndPath := newPatch.Operation.String() + newPatch.Path
		prevPatch, ok := r.desiredPatches[operationAndPath]
		if !ok {
			// This is a new, unique patch.
			r.desiredPatches[operationAndPath] = newPatch
			continue
		}

		// Attempt to convert the value to a bool so we can see if the
		// patch we already have for this operation reverts the pre-existing
		// one or is a dupe.
		newPatchValStr, ok := newPatch.Value.(string)
		if !ok {
			return fmt.Errorf("all patches must have bool values but received %+x", newPatch)
		}
		newPatchVal, err := strconv.ParseBool(newPatchValStr)
		if err != nil {
			return err
		}

		// This was already verified to not be a bool string when it was added.
		prevPatchVal, _ := strconv.ParseBool(prevPatch.Value.(string))
		if newPatchVal != prevPatchVal {
			// The new patch cancels/reverts the need to apply the prev patch
			// because we've now returned to the original state.
			delete(r.desiredPatches, operationAndPath)
		}
		// If we arrive here, the new patch is a dupe of the previous patch.
		// We don't need to add it to the desired patches because it already
		// exists.
	}
	return nil
}

func (r *retryer) attemptDesiredPatches() {
	r.desiredPatchesLock.Lock()
	defer r.desiredPatchesLock.Unlock()

	if len(r.desiredPatches) == 0 {
		return
	}
	patches := make([]*client.Patch, len(r.desiredPatches))
	i := 0
	for _, patch := range r.desiredPatches {
		patches[i] = patch
		i++
	}
	if err := r.client.PatchPod(r.namespace, r.podName, patches...); err != nil {
		if r.logger.IsWarn() {
			r.logger.Warn("unable to update state due to %s, will retry", err.Error())
		}
		return
	}
	// We succeeded at applying the patches! We're back in a normal state.
	r.desiredPatches = make(map[string]*client.Patch)
}

// Converts a bool to "true" or "false".
func toString(b bool) string {
	return fmt.Sprintf("%t", b)
}
