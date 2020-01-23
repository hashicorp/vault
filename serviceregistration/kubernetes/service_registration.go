package kubernetes

import (
	"fmt"
	"os"

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
)

func NewServiceRegistration(shutdownCh <-chan struct{}, config map[string]string, logger log.Logger, state *sr.State, _ string) (sr.ServiceRegistration, error) {
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
		podName:   podName,
		namespace: namespace,
		client:    c,
	}

	// Run a background goroutine to leave labels in the final state we'd like
	// when Vault shuts down.
	go registration.onShutdown(shutdownCh)
	return registration, nil
}

type serviceRegistration struct {
	logger             log.Logger
	namespace, podName string
	client             *client.Client
}

func (r *serviceRegistration) NotifyActiveStateChange(isActive bool) error {
	return r.client.PatchPod(r.namespace, r.podName, &client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelActive,
		Value:     toString(isActive),
	})
}

func (r *serviceRegistration) NotifySealedStateChange(isSealed bool) error {
	return r.client.PatchPod(r.namespace, r.podName, &client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelSealed,
		Value:     toString(isSealed),
	})
}

func (r *serviceRegistration) NotifyPerformanceStandbyStateChange(isStandby bool) error {
	return r.client.PatchPod(r.namespace, r.podName, &client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelPerfStandby,
		Value:     toString(isStandby),
	})
}

func (r *serviceRegistration) NotifyInitializedStateChange(isInitialized bool) error {
	return r.client.PatchPod(r.namespace, r.podName, &client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelInitialized,
		Value:     toString(isInitialized),
	})
}

func (r *serviceRegistration) onShutdown(shutdownCh <-chan struct{}) {
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

// Converts a bool to "true" or "false".
func toString(b bool) string {
	return fmt.Sprintf("%t", b)
}
