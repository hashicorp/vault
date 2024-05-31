// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

func NewServiceRegistration(config map[string]string, logger hclog.Logger, state sr.State) (sr.ServiceRegistration, error) {
	namespace, err := getRequiredField(logger, config, client.EnvVarKubernetesNamespace, "namespace")
	if err != nil {
		return nil, err
	}
	podName, err := getRequiredField(logger, config, client.EnvVarKubernetesPodName, "pod_name")
	if err != nil {
		return nil, err
	}

	c, err := client.New(logger)
	if err != nil {
		return nil, err
	}

	// The Vault version must be sanitized because it can contain special
	// characters like "+" which aren't acceptable by the Kube API.
	state.VaultVersion = client.Sanitize(state.VaultVersion)
	return &serviceRegistration{
		logger:    logger,
		namespace: namespace,
		podName:   podName,
		retryHandler: &retryHandler{
			logger:         logger,
			namespace:      namespace,
			podName:        podName,
			initialState:   state,
			patchesToRetry: make(map[string]*client.Patch),
			client:         c,
		},
	}, nil
}

type serviceRegistration struct {
	logger             hclog.Logger
	namespace, podName string
	retryHandler       *retryHandler
}

func (r *serviceRegistration) Run(shutdownCh <-chan struct{}, wait *sync.WaitGroup, _ string) error {
	r.retryHandler.Run(shutdownCh, wait)
	return nil
}

func (r *serviceRegistration) NotifyActiveStateChange(isActive bool) error {
	r.retryHandler.Notify(&client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelActive,
		Value:     strconv.FormatBool(isActive),
	})
	return nil
}

func (r *serviceRegistration) NotifySealedStateChange(isSealed bool) error {
	r.retryHandler.Notify(&client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelSealed,
		Value:     strconv.FormatBool(isSealed),
	})
	return nil
}

func (r *serviceRegistration) NotifyPerformanceStandbyStateChange(isStandby bool) error {
	r.retryHandler.Notify(&client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelPerfStandby,
		Value:     strconv.FormatBool(isStandby),
	})
	return nil
}

func (r *serviceRegistration) NotifyInitializedStateChange(isInitialized bool) error {
	r.retryHandler.Notify(&client.Patch{
		Operation: client.Replace,
		Path:      pathToLabels + labelInitialized,
		Value:     strconv.FormatBool(isInitialized),
	})
	return nil
}

func (c *serviceRegistration) NotifyConfigurationReload(conf *map[string]string) error {
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
