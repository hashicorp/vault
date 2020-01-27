package kubernetes

import (
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/go-hclog"
	sr "github.com/hashicorp/vault/serviceregistration"
	"github.com/hashicorp/vault/serviceregistration/kubernetes/client"
)

const (
	// Labels are placed in a pod's metadata.
	labelVaultVersion = "vault-version"
	labelActive       = "vault-active"
	labelSealed       = "vault-sealed"
	labelPerfStandby  = "vault-perf-standby"
	labelInitialized  = "vault-initialized"

	// This is the path to where these labels are applied.
	pathToLabels = "/metadata/labels/"
)

func NewServiceRegistration(config map[string]string, logger hclog.Logger, state sr.State, _ string) (sr.ServiceRegistration, error) {
	namespace, err := getRequiredField(logger, config, client.EnvVarKubernetesNamespace, "namespace")
	if err != nil {
		return nil, err
	}
	podName, err := getRequiredField(logger, config, client.EnvVarKubernetesPodName, "pod_name")
	if err != nil {
		return nil, err
	}
	return &serviceRegistration{
		logger:       logger,
		namespace:    namespace,
		podName:      podName,
		initialState: state,
		retryHandler: &retryHandler{
			logger:         logger,
			namespace:      namespace,
			podName:        podName,
			patchesToRetry: make([]*client.Patch, 0),
		},
	}, nil
}

type serviceRegistration struct {
	logger             hclog.Logger
	namespace, podName string
	client             *client.Client
	initialState       sr.State
	retryHandler       *retryHandler
}

func (r *serviceRegistration) Run(shutdownCh <-chan struct{}, wait *sync.WaitGroup) error {
	c, err := client.New(r.logger, shutdownCh)
	if err != nil {
		return err
	}
	r.client = c
	r.retryHandler.client = c

	// Verify that the pod exists and our configuration looks good.
	pod, err := c.GetPod(r.namespace, r.podName)
	if err != nil {
		return err
	}

	// Now to initially label our pod.
	if pod.Metadata == nil {
		// This should never happen IRL, just being defensive.
		return fmt.Errorf("no pod metadata on %+v", pod)
	}
	if pod.Metadata.Labels == nil {
		// Add the labels field, and the labels as part of that one call.
		// The reason we must take a different approach to adding them is discussed here:
		// https://stackoverflow.com/questions/57480205/error-while-applying-json-patch-to-kubernetes-custom-resource
		if err := c.PatchPod(r.namespace, r.podName, &client.Patch{
			Operation: client.Add,
			Path:      "/metadata/labels",
			Value: map[string]string{
				labelVaultVersion: r.initialState.VaultVersion,
				labelActive:       toString(r.initialState.IsActive),
				labelSealed:       toString(r.initialState.IsSealed),
				labelPerfStandby:  toString(r.initialState.IsPerformanceStandby),
				labelInitialized:  toString(r.initialState.IsInitialized),
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
				Value:     toString(r.initialState.IsActive),
			},
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelSealed,
				Value:     toString(r.initialState.IsSealed),
			},
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelPerfStandby,
				Value:     toString(r.initialState.IsPerformanceStandby),
			},
			{
				Operation: client.Replace,
				Path:      pathToLabels + labelInitialized,
				Value:     toString(r.initialState.IsInitialized),
			},
		}
		if err := c.PatchPod(r.namespace, r.podName, patches...); err != nil {
			return err
		}
	}

	// Run a service that retries errored-out notifications if they occur.
	go r.retryHandler.Run(shutdownCh, wait)

	// Run a background goroutine to leave labels in the final state we'd like
	// when Vault shuts down.
	go r.onShutdown(shutdownCh, wait)

	return nil
}

func (r *serviceRegistration) NotifyActiveStateChange(isActive bool) error {
	return r.notifyOrRetry(&client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelActive,
		Value:     toString(isActive),
	})
}

func (r *serviceRegistration) NotifySealedStateChange(isSealed bool) error {
	return r.notifyOrRetry(&client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelSealed,
		Value:     toString(isSealed),
	})
}

func (r *serviceRegistration) NotifyPerformanceStandbyStateChange(isStandby bool) error {
	return r.notifyOrRetry(&client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelPerfStandby,
		Value:     toString(isStandby),
	})
}

func (r *serviceRegistration) NotifyInitializedStateChange(isInitialized bool) error {
	return r.notifyOrRetry(&client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelInitialized,
		Value:     toString(isInitialized),
	})
}

func (r *serviceRegistration) onShutdown(shutdownCh <-chan struct{}, wait *sync.WaitGroup) {
	// Make sure Vault will give us time to finish up here.
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
		if err := r.retryHandler.Add(patch); err != nil {
			return err
		}
	}
	return nil
}

func getRequiredField(logger hclog.Logger, config map[string]string, envVar, configParam string) (string, error) {
	value := ""
	switch {
	case os.Getenv(envVar) != "":
		value = os.Getenv(envVar)
	case config[configParam] != "":
		value = config[configParam]
	default:
		return "", fmt.Errorf(`%s must be provided via %q or the %q config parameter`, configParam, envVar, configParam)
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("%q: %q", configParam, value))
	}
	return value, nil
}

// Converts a bool to "true" or "false".
func toString(b bool) string {
	return fmt.Sprintf("%t", b)
}
