package kubernetes

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/transport"
)

// Clientset is an alias to k8s.io/client-go/kubernetes.Interface
type Clientset kubernetes.Interface

// NewClientsetFromBytes builds a Clientset from a given Kubeconfig.
//
// Takes an optional transport.WrapperFunc to add request/response middleware to
// api-server requests.
func BuildClientsetFromConfig(
	kubeconfigBytes []byte,
	transportWrapper transport.WrapperFunc,
) (Clientset, error) {
	config, err := clientcmd.NewClientConfigFromBytes(kubeconfigBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LKE cluster kubeconfig: %s", err)
	}

	restClientConfig, err := config.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get REST client config: %s", err)
	}

	if transportWrapper != nil {
		restClientConfig.Wrap(transportWrapper)
	}

	clientset, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build k8s client from LKE cluster kubeconfig: %s", err)
	}
	return clientset, nil
}
