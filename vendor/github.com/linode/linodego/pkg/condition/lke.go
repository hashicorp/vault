package condition

import (
	"context"
	"errors"
	"fmt"

	"github.com/linode/linodego/internal/kubernetes"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterConditionFunc represents a function that tests a condition against an LKE cluster,
// returns true if the condition has been reached, false if it has not yet been reached.
type ClusterConditionFunc func(context.Context, kubernetes.Clientset) (bool, error)

// ClusterHasReadyNode is a ClusterConditionFunc which polls for at least one node to have the
// condition NodeReady=True.
func ClusterHasReadyNode(ctx context.Context, clientset kubernetes.Clientset) (bool, error) {
	nodes, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get nodes for cluster: %s", err)
	}

	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
	}

	return false, errors.New("no nodes in cluster are ready")
}
