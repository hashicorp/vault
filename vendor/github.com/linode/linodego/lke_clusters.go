package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// LKEClusterStatus represents the status of an LKECluster
type LKEClusterStatus string

// LKEClusterStatus enums start with LKECluster
const (
	LKEClusterReady    LKEClusterStatus = "ready"
	LKEClusterNotReady LKEClusterStatus = "not_ready"
)

// LKECluster represents a LKECluster object
type LKECluster struct {
	ID           int                    `json:"id"`
	Created      *time.Time             `json:"-"`
	Updated      *time.Time             `json:"-"`
	Label        string                 `json:"label"`
	Region       string                 `json:"region"`
	Status       LKEClusterStatus       `json:"status"`
	K8sVersion   string                 `json:"k8s_version"`
	Tags         []string               `json:"tags"`
	ControlPlane LKEClusterControlPlane `json:"control_plane"`
}

// LKEClusterCreateOptions fields are those accepted by CreateLKECluster
type LKEClusterCreateOptions struct {
	NodePools    []LKENodePoolCreateOptions     `json:"node_pools"`
	Label        string                         `json:"label"`
	Region       string                         `json:"region"`
	K8sVersion   string                         `json:"k8s_version"`
	Tags         []string                       `json:"tags,omitempty"`
	ControlPlane *LKEClusterControlPlaneOptions `json:"control_plane,omitempty"`
}

// LKEClusterUpdateOptions fields are those accepted by UpdateLKECluster
type LKEClusterUpdateOptions struct {
	K8sVersion   string                         `json:"k8s_version,omitempty"`
	Label        string                         `json:"label,omitempty"`
	Tags         *[]string                      `json:"tags,omitempty"`
	ControlPlane *LKEClusterControlPlaneOptions `json:"control_plane,omitempty"`
}

// LKEClusterAPIEndpoint fields are those returned by ListLKEClusterAPIEndpoints
type LKEClusterAPIEndpoint struct {
	Endpoint string `json:"endpoint"`
}

// LKEClusterKubeconfig fields are those returned by GetLKEClusterKubeconfig
type LKEClusterKubeconfig struct {
	KubeConfig string `json:"kubeconfig"`
}

// LKEClusterDashboard fields are those returned by GetLKEClusterDashboard
type LKEClusterDashboard struct {
	URL string `json:"url"`
}

// LKEVersion fields are those returned by GetLKEVersion
type LKEVersion struct {
	ID string `json:"id"`
}

// LKEClusterRegenerateOptions fields are those accepted by RegenerateLKECluster
type LKEClusterRegenerateOptions struct {
	KubeConfig   bool `json:"kubeconfig"`
	ServiceToken bool `json:"servicetoken"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *LKECluster) UnmarshalJSON(b []byte) error {
	type Mask LKECluster

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.Updated = (*time.Time)(p.Updated)

	return nil
}

// GetCreateOptions converts a LKECluster to LKEClusterCreateOptions for use in CreateLKECluster
func (i LKECluster) GetCreateOptions() (o LKEClusterCreateOptions) {
	o.Label = i.Label
	o.Region = i.Region
	o.K8sVersion = i.K8sVersion
	o.Tags = i.Tags

	isHA := i.ControlPlane.HighAvailability

	o.ControlPlane = &LKEClusterControlPlaneOptions{
		HighAvailability: &isHA,
		// ACL will not be populated in the control plane response
	}

	// @TODO copy NodePools?
	return
}

// GetUpdateOptions converts a LKECluster to LKEClusterUpdateOptions for use in UpdateLKECluster
func (i LKECluster) GetUpdateOptions() (o LKEClusterUpdateOptions) {
	o.K8sVersion = i.K8sVersion
	o.Label = i.Label
	o.Tags = &i.Tags

	isHA := i.ControlPlane.HighAvailability

	o.ControlPlane = &LKEClusterControlPlaneOptions{
		HighAvailability: &isHA,
		// ACL will not be populated in the control plane response
	}

	return
}

// ListLKEVersions lists the Kubernetes versions available through LKE. This endpoint is cached by default.
func (c *Client) ListLKEVersions(ctx context.Context, opts *ListOptions) ([]LKEVersion, error) {
	e := "lke/versions"

	endpoint, err := generateListCacheURL(e, opts)
	if err != nil {
		return nil, err
	}

	if result := c.getCachedResponse(endpoint); result != nil {
		return result.([]LKEVersion), nil
	}

	response, err := getPaginatedResults[LKEVersion](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(endpoint, response, &cacheExpiryTime)

	return response, nil
}

// GetLKEVersion gets details about a specific LKE Version. This endpoint is cached by default.
func (c *Client) GetLKEVersion(ctx context.Context, version string) (*LKEVersion, error) {
	e := formatAPIPath("lke/versions/%s", version)

	if result := c.getCachedResponse(e); result != nil {
		result := result.(LKEVersion)
		return &result, nil
	}

	response, err := doGETRequest[LKEVersion](ctx, c, e)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(e, response, &cacheExpiryTime)

	return response, nil
}

// ListLKEClusterAPIEndpoints gets the API Endpoint for the LKE Cluster specified
func (c *Client) ListLKEClusterAPIEndpoints(ctx context.Context, clusterID int, opts *ListOptions) ([]LKEClusterAPIEndpoint, error) {
	response, err := getPaginatedResults[LKEClusterAPIEndpoint](ctx, c, formatAPIPath("lke/clusters/%d/api-endpoints", clusterID), opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ListLKEClusters lists LKEClusters
func (c *Client) ListLKEClusters(ctx context.Context, opts *ListOptions) ([]LKECluster, error) {
	response, err := getPaginatedResults[LKECluster](ctx, c, "lke/clusters", opts)
	return response, err
}

// GetLKECluster gets the lkeCluster with the provided ID
func (c *Client) GetLKECluster(ctx context.Context, clusterID int) (*LKECluster, error) {
	e := formatAPIPath("lke/clusters/%d", clusterID)
	response, err := doGETRequest[LKECluster](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateLKECluster creates a LKECluster
func (c *Client) CreateLKECluster(ctx context.Context, opts LKEClusterCreateOptions) (*LKECluster, error) {
	e := "lke/clusters"
	response, err := doPOSTRequest[LKECluster](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateLKECluster updates the LKECluster with the specified id
func (c *Client) UpdateLKECluster(ctx context.Context, clusterID int, opts LKEClusterUpdateOptions) (*LKECluster, error) {
	e := formatAPIPath("lke/clusters/%d", clusterID)
	response, err := doPUTRequest[LKECluster](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteLKECluster deletes the LKECluster with the specified id
func (c *Client) DeleteLKECluster(ctx context.Context, clusterID int) error {
	e := formatAPIPath("lke/clusters/%d", clusterID)
	err := doDELETERequest(ctx, c, e)
	return err
}

// GetLKEClusterKubeconfig gets the Kubeconfig for the LKE Cluster specified
func (c *Client) GetLKEClusterKubeconfig(ctx context.Context, clusterID int) (*LKEClusterKubeconfig, error) {
	e := formatAPIPath("lke/clusters/%d/kubeconfig", clusterID)
	response, err := doGETRequest[LKEClusterKubeconfig](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetLKEClusterDashboard gets information about the dashboard for an LKE cluster
func (c *Client) GetLKEClusterDashboard(ctx context.Context, clusterID int) (*LKEClusterDashboard, error) {
	e := formatAPIPath("lke/clusters/%d/dashboard", clusterID)
	response, err := doGETRequest[LKEClusterDashboard](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// RecycleLKEClusterNodes recycles all nodes in all pools of the specified LKE Cluster.
func (c *Client) RecycleLKEClusterNodes(ctx context.Context, clusterID int) error {
	e := formatAPIPath("lke/clusters/%d/recycle", clusterID)
	_, err := doPOSTRequest[LKECluster, any](ctx, c, e)
	return err
}

// RegenerateLKECluster regenerates the Kubeconfig file and/or the service account token for the specified LKE Cluster.
func (c *Client) RegenerateLKECluster(ctx context.Context, clusterID int, opts LKEClusterRegenerateOptions) (*LKECluster, error) {
	e := formatAPIPath("lke/clusters/%d/regenerate", clusterID)
	response, err := doPOSTRequest[LKECluster](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteLKEClusterServiceToken deletes and regenerate the service account token for a Cluster.
func (c *Client) DeleteLKEClusterServiceToken(ctx context.Context, clusterID int) error {
	e := formatAPIPath("lke/clusters/%d/servicetoken", clusterID)
	err := doDELETERequest(ctx, c, e)
	return err
}
