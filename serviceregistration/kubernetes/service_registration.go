package kubernetes

import (
	"fmt"
	"os"
	"strconv"
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
			patchesToRetry: make(map[string]*client.Patch),
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

	// Now that we've populated the client, we can begin using the retry handler.
	r.retryHandler.SetInitialState(r.setInitialState)
	r.retryHandler.Run(shutdownCh, wait, c)
	return nil
}

func (r *serviceRegistration) setInitialState() error {
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
	return nil
}

func (r *serviceRegistration) NotifyActiveStateChange(isActive bool) error {
	r.retryHandler.Notify(r.client, &client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelActive,
		Value:     strconv.FormatBool(isActive),
	})
	return nil
}

func (r *serviceRegistration) NotifySealedStateChange(isSealed bool) error {
	r.retryHandler.Notify(r.client, &client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelSealed,
		Value:     strconv.FormatBool(isSealed),
	})
	return nil
}

func (r *serviceRegistration) NotifyPerformanceStandbyStateChange(isStandby bool) error {
	r.retryHandler.Notify(r.client, &client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelPerfStandby,
		Value:     strconv.FormatBool(isStandby),
	})
	return nil
}

func (r *serviceRegistration) NotifyInitializedStateChange(isInitialized bool) error {
	r.retryHandler.Notify(r.client, &client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelInitialized,
		Value:     strconv.FormatBool(isInitialized),
	})
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
