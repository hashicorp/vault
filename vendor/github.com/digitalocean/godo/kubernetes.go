package godo

import (
	"bytes"
	"context"
	"encoding"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	kubernetesBasePath     = "/v2/kubernetes"
	kubernetesClustersPath = kubernetesBasePath + "/clusters"
	kubernetesOptionsPath  = kubernetesBasePath + "/options"
)

// KubernetesService is an interface for interfacing with the kubernetes endpoints
// of the DigitalOcean API.
// See: https://developers.digitalocean.com/documentation/v2#kubernetes
type KubernetesService interface {
	Create(context.Context, *KubernetesClusterCreateRequest) (*KubernetesCluster, *Response, error)
	Get(context.Context, string) (*KubernetesCluster, *Response, error)
	GetKubeConfig(context.Context, string) (*KubernetesClusterConfig, *Response, error)
	List(context.Context, *ListOptions) ([]*KubernetesCluster, *Response, error)
	Update(context.Context, string, *KubernetesClusterUpdateRequest) (*KubernetesCluster, *Response, error)
	Delete(context.Context, string) (*Response, error)

	CreateNodePool(ctx context.Context, clusterID string, req *KubernetesNodePoolCreateRequest) (*KubernetesNodePool, *Response, error)
	GetNodePool(ctx context.Context, clusterID, poolID string) (*KubernetesNodePool, *Response, error)
	ListNodePools(ctx context.Context, clusterID string, opts *ListOptions) ([]*KubernetesNodePool, *Response, error)
	UpdateNodePool(ctx context.Context, clusterID, poolID string, req *KubernetesNodePoolUpdateRequest) (*KubernetesNodePool, *Response, error)
	RecycleNodePoolNodes(ctx context.Context, clusterID, poolID string, req *KubernetesNodePoolRecycleNodesRequest) (*Response, error)
	DeleteNodePool(ctx context.Context, clusterID, poolID string) (*Response, error)

	GetOptions(context.Context) (*KubernetesOptions, *Response, error)
}

var _ KubernetesService = &KubernetesServiceOp{}

// KubernetesServiceOp handles communication with Kubernetes methods of the DigitalOcean API.
type KubernetesServiceOp struct {
	client *Client
}

// KubernetesClusterCreateRequest represents a request to create a Kubernetes cluster.
type KubernetesClusterCreateRequest struct {
	Name        string   `json:"name,omitempty"`
	RegionSlug  string   `json:"region,omitempty"`
	VersionSlug string   `json:"version,omitempty"`
	Tags        []string `json:"tags,omitempty"`

	NodePools []*KubernetesNodePoolCreateRequest `json:"node_pools,omitempty"`
}

// KubernetesClusterUpdateRequest represents a request to update a Kubernetes cluster.
type KubernetesClusterUpdateRequest struct {
	Name string   `json:"name,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

// KubernetesNodePoolCreateRequest represents a request to create a node pool for a
// Kubernetes cluster.
type KubernetesNodePoolCreateRequest struct {
	Name  string   `json:"name,omitempty"`
	Size  string   `json:"size,omitempty"`
	Count int      `json:"count,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

// KubernetesNodePoolUpdateRequest represents a request to update a node pool in a
// Kubernetes cluster.
type KubernetesNodePoolUpdateRequest struct {
	Name  string   `json:"name,omitempty"`
	Count int      `json:"count,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

// KubernetesNodePoolRecycleNodesRequest represents a request to recycle a set of
// nodes in a node pool. This will recycle the nodes by ID.
type KubernetesNodePoolRecycleNodesRequest struct {
	Nodes []string `json:"nodes,omitempty"`
}

// KubernetesCluster represents a Kubernetes cluster.
type KubernetesCluster struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	RegionSlug    string   `json:"region,omitempty"`
	VersionSlug   string   `json:"version,omitempty"`
	ClusterSubnet string   `json:"cluster_subnet,omitempty"`
	ServiceSubnet string   `json:"service_subnet,omitempty"`
	IPv4          string   `json:"ipv4,omitempty"`
	Endpoint      string   `json:"endpoint,omitempty"`
	Tags          []string `json:"tags,omitempty"`

	NodePools []*KubernetesNodePool `json:"node_pools,omitempty"`

	Status    *KubernetesClusterStatus `json:"status,omitempty"`
	CreatedAt time.Time                `json:"created_at,omitempty"`
	UpdatedAt time.Time                `json:"updated_at,omitempty"`
}

// Possible states for a cluster.
const (
	KubernetesClusterStatusProvisioning = KubernetesClusterStatusState("provisioning")
	KubernetesClusterStatusRunning      = KubernetesClusterStatusState("running")
	KubernetesClusterStatusDegraded     = KubernetesClusterStatusState("degraded")
	KubernetesClusterStatusError        = KubernetesClusterStatusState("error")
	KubernetesClusterStatusDeleted      = KubernetesClusterStatusState("deleted")
	KubernetesClusterStatusInvalid      = KubernetesClusterStatusState("invalid")
)

// KubernetesClusterStatusState represents states for a cluster.
type KubernetesClusterStatusState string

var _ encoding.TextUnmarshaler = (*KubernetesClusterStatusState)(nil)

// UnmarshalText unmarshals the state.
func (s *KubernetesClusterStatusState) UnmarshalText(text []byte) error {
	switch KubernetesClusterStatusState(strings.ToLower(string(text))) {
	case KubernetesClusterStatusProvisioning:
		*s = KubernetesClusterStatusProvisioning
	case KubernetesClusterStatusRunning:
		*s = KubernetesClusterStatusRunning
	case KubernetesClusterStatusDegraded:
		*s = KubernetesClusterStatusDegraded
	case KubernetesClusterStatusError:
		*s = KubernetesClusterStatusError
	case KubernetesClusterStatusDeleted:
		*s = KubernetesClusterStatusDeleted
	case "", KubernetesClusterStatusInvalid:
		*s = KubernetesClusterStatusInvalid
	default:
		return fmt.Errorf("unknown cluster state %q", string(text))
	}
	return nil
}

// KubernetesClusterStatus describes the status of a cluster.
type KubernetesClusterStatus struct {
	State   KubernetesClusterStatusState `json:"state,omitempty"`
	Message string                       `json:"message,omitempty"`
}

// KubernetesNodePool represents a node pool in a Kubernetes cluster.
type KubernetesNodePool struct {
	ID    string   `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Size  string   `json:"size,omitempty"`
	Count int      `json:"count,omitempty"`
	Tags  []string `json:"tags,omitempty"`

	Nodes []*KubernetesNode `json:"nodes,omitempty"`
}

// KubernetesNode represents a Node in a node pool in a Kubernetes cluster.
type KubernetesNode struct {
	ID     string                `json:"id,omitempty"`
	Name   string                `json:"name,omitempty"`
	Status *KubernetesNodeStatus `json:"status,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// KubernetesNodeStatus represents the status of a particular Node in a Kubernetes cluster.
type KubernetesNodeStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// KubernetesOptions represents options available for creating Kubernetes clusters.
type KubernetesOptions struct {
	Versions []*KubernetesVersion  `json:"versions,omitempty"`
	Regions  []*KubernetesRegion   `json:"regions,omitempty"`
	Sizes    []*KubernetesNodeSize `json:"sizes,omitempty"`
}

// KubernetesVersion is a DigitalOcean Kubernetes release.
type KubernetesVersion struct {
	Slug              string `json:"slug,omitempty"`
	KubernetesVersion string `json:"kubernetes_version,omitempty"`
}

// KubernetesNodeSize is a node sizes supported for Kubernetes clusters.
type KubernetesNodeSize struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// KubernetesRegion is a region usable by Kubernetes clusters.
type KubernetesRegion struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type kubernetesClustersRoot struct {
	Clusters []*KubernetesCluster `json:"kubernetes_clusters,omitempty"`
	Links    *Links               `json:"links,omitempty"`
}

type kubernetesClusterRoot struct {
	Cluster *KubernetesCluster `json:"kubernetes_cluster,omitempty"`
}

type kubernetesNodePoolRoot struct {
	NodePool *KubernetesNodePool `json:"node_pool,omitempty"`
}

type kubernetesNodePoolsRoot struct {
	NodePools []*KubernetesNodePool `json:"node_pools,omitempty"`
	Links     *Links                `json:"links,omitempty"`
}

// Get retrieves the details of a Kubernetes cluster.
func (svc *KubernetesServiceOp) Get(ctx context.Context, clusterID string) (*KubernetesCluster, *Response, error) {
	path := fmt.Sprintf("%s/%s", kubernetesClustersPath, clusterID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(kubernetesClusterRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Cluster, resp, nil
}

// Create creates a Kubernetes cluster.
func (svc *KubernetesServiceOp) Create(ctx context.Context, create *KubernetesClusterCreateRequest) (*KubernetesCluster, *Response, error) {
	path := kubernetesClustersPath
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, nil, err
	}
	root := new(kubernetesClusterRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Cluster, resp, nil
}

// Delete deletes a Kubernetes cluster. There is no way to recover a cluster
// once it has been destroyed.
func (svc *KubernetesServiceOp) Delete(ctx context.Context, clusterID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", kubernetesClustersPath, clusterID)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// List returns a list of the Kubernetes clusters visible with the caller's API token.
func (svc *KubernetesServiceOp) List(ctx context.Context, opts *ListOptions) ([]*KubernetesCluster, *Response, error) {
	path := kubernetesClustersPath
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(kubernetesClustersRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Clusters, resp, nil
}

// KubernetesClusterConfig is the content of a Kubernetes config file, which can be
// used to interact with your Kubernetes cluster using `kubectl`.
// See: https://kubernetes.io/docs/tasks/tools/install-kubectl/
type KubernetesClusterConfig struct {
	KubeconfigYAML []byte
}

// GetKubeConfig returns a Kubernetes config file for the specified cluster.
func (svc *KubernetesServiceOp) GetKubeConfig(ctx context.Context, clusterID string) (*KubernetesClusterConfig, *Response, error) {
	path := fmt.Sprintf("%s/%s/kubeconfig", kubernetesClustersPath, clusterID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	configBytes := bytes.NewBuffer(nil)
	resp, err := svc.client.Do(ctx, req, configBytes)
	if err != nil {
		return nil, resp, err
	}
	res := &KubernetesClusterConfig{
		KubeconfigYAML: configBytes.Bytes(),
	}
	return res, resp, nil
}

// Update updates a Kubernetes cluster's properties.
func (svc *KubernetesServiceOp) Update(ctx context.Context, clusterID string, update *KubernetesClusterUpdateRequest) (*KubernetesCluster, *Response, error) {
	path := fmt.Sprintf("%s/%s", kubernetesClustersPath, clusterID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, nil, err
	}
	root := new(kubernetesClusterRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Cluster, resp, nil
}

// CreateNodePool creates a new node pool in an existing Kubernetes cluster.
func (svc *KubernetesServiceOp) CreateNodePool(ctx context.Context, clusterID string, create *KubernetesNodePoolCreateRequest) (*KubernetesNodePool, *Response, error) {
	path := fmt.Sprintf("%s/%s/node_pools", kubernetesClustersPath, clusterID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, nil, err
	}
	root := new(kubernetesNodePoolRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.NodePool, resp, nil
}

// GetNodePool retrieves an existing node pool in a Kubernetes cluster.
func (svc *KubernetesServiceOp) GetNodePool(ctx context.Context, clusterID, poolID string) (*KubernetesNodePool, *Response, error) {
	path := fmt.Sprintf("%s/%s/node_pools/%s", kubernetesClustersPath, clusterID, poolID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(kubernetesNodePoolRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.NodePool, resp, nil
}

// ListNodePools lists all the node pools found in a Kubernetes cluster.
func (svc *KubernetesServiceOp) ListNodePools(ctx context.Context, clusterID string, opts *ListOptions) ([]*KubernetesNodePool, *Response, error) {
	path := fmt.Sprintf("%s/%s/node_pools", kubernetesClustersPath, clusterID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(kubernetesNodePoolsRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.NodePools, resp, nil
}

// UpdateNodePool updates the details of an existing node pool.
func (svc *KubernetesServiceOp) UpdateNodePool(ctx context.Context, clusterID, poolID string, update *KubernetesNodePoolUpdateRequest) (*KubernetesNodePool, *Response, error) {
	path := fmt.Sprintf("%s/%s/node_pools/%s", kubernetesClustersPath, clusterID, poolID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, nil, err
	}
	root := new(kubernetesNodePoolRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.NodePool, resp, nil
}

// RecycleNodePoolNodes schedules nodes in a node pool for recycling.
func (svc *KubernetesServiceOp) RecycleNodePoolNodes(ctx context.Context, clusterID, poolID string, recycle *KubernetesNodePoolRecycleNodesRequest) (*Response, error) {
	path := fmt.Sprintf("%s/%s/node_pools/%s/recycle", kubernetesClustersPath, clusterID, poolID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, recycle)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// DeleteNodePool deletes a node pool, and subsequently all the nodes in that pool.
func (svc *KubernetesServiceOp) DeleteNodePool(ctx context.Context, clusterID, poolID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s/node_pools/%s", kubernetesClustersPath, clusterID, poolID)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

type kubernetesOptionsRoot struct {
	Options *KubernetesOptions `json:"options,omitempty"`
	Links   *Links             `json:"links,omitempty"`
}

// GetOptions returns options about the Kubernetes service, such as the versions available for
// cluster creation.
func (svc *KubernetesServiceOp) GetOptions(ctx context.Context) (*KubernetesOptions, *Response, error) {
	path := kubernetesOptionsPath
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(kubernetesOptionsRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Options, resp, nil
}
