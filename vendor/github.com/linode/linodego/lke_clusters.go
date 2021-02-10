package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
	"github.com/linode/linodego/pkg/errors"
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
	ID         int              `json:"id"`
	Created    *time.Time       `json:"-"`
	Updated    *time.Time       `json:"-"`
	Label      string           `json:"label"`
	Region     string           `json:"region"`
	Status     LKEClusterStatus `json:"status"`
	K8sVersion string           `json:"k8s_version"`
	Tags       []string         `json:"tags"`
}

// LKEClusterCreateOptions fields are those accepted by CreateLKECluster
type LKEClusterCreateOptions struct {
	NodePools  []LKEClusterPoolCreateOptions `json:"node_pools"`
	Label      string                        `json:"label"`
	Region     string                        `json:"region"`
	K8sVersion string                        `json:"k8s_version"`
	Tags       []string                      `json:"tags,omitempty"`
}

// LKEClusterUpdateOptions fields are those accepted by UpdateLKECluster
type LKEClusterUpdateOptions struct {
	Label string    `json:"label,omitempty"`
	Tags  *[]string `json:"tags,omitempty"`
}

// LKEClusterAPIEndpoint fields are those returned by ListLKEClusterAPIEndpoints
type LKEClusterAPIEndpoint struct {
	Endpoint string `json:"endpoint"`
}

// LKEClusterKubeconfig fields are those returned by GetLKEClusterKubeconfig
type LKEClusterKubeconfig struct {
	KubeConfig string `json:"kubeconfig"`
}

// LKEVersion fields are those returned by GetLKEVersion
type LKEVersion struct {
	ID string `json:"id"`
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
	// @TODO copy NodePools?
	return
}

// GetUpdateOptions converts a LKECluster to LKEClusterUpdateOptions for use in UpdateLKECluster
func (i LKECluster) GetUpdateOptions() (o LKEClusterUpdateOptions) {
	o.Label = i.Label
	o.Tags = &i.Tags
	return
}

// LKEClustersPagedResponse represents a paginated LKECluster API response
type LKEClustersPagedResponse struct {
	*PageOptions
	Data []LKECluster `json:"data"`
}

// LKEClusterAPIEndpointsPagedResponse represents a paginated LKEClusterAPIEndpoints API response
type LKEClusterAPIEndpointsPagedResponse struct {
	*PageOptions
	Data []LKEClusterAPIEndpoint `json:"data"`
}

// LKEVersionsPagedResponse represents a paginated LKEVersion API response
type LKEVersionsPagedResponse struct {
	*PageOptions
	Data []LKEVersion `json:"data"`
}

// endpoint gets the endpoint URL for LKECluster
func (LKEClustersPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.LKEClusters.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends LKEClusters when processing paginated LKECluster responses
func (resp *LKEClustersPagedResponse) appendData(r *LKEClustersPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// endpoint gets the endpoint URL for LKEVersion
func (LKEVersionsPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.LKEVersions.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// endpoint gets the endpoint URL for LKEClusterAPIEndpoints
func (LKEClusterAPIEndpointsPagedResponse) endpointWithID(c *Client, id int) string {
	endpoint, err := c.LKEClusterAPIEndpoints.endpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends LKEClusterAPIEndpoints when processing paginated LKEClusterAPIEndpoints responses
func (resp *LKEClusterAPIEndpointsPagedResponse) appendData(r *LKEClusterAPIEndpointsPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// appendData appends LKEVersions when processing paginated LKEVersion responses
func (resp *LKEVersionsPagedResponse) appendData(r *LKEVersionsPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListLKEClusters lists LKEClusters
func (c *Client) ListLKEClusters(ctx context.Context, opts *ListOptions) ([]LKECluster, error) {
	response := LKEClustersPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetLKECluster gets the lkeCluster with the provided ID
func (c *Client) GetLKECluster(ctx context.Context, id int) (*LKECluster, error) {
	e, err := c.LKEClusters.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)
	r, err := errors.CoupleAPIErrors(c.R(ctx).SetResult(&LKECluster{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*LKECluster), nil
}

// CreateLKECluster creates a LKECluster
func (c *Client) CreateLKECluster(ctx context.Context, createOpts LKEClusterCreateOptions) (*LKECluster, error) {
	var body string
	e, err := c.LKEClusters.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&LKECluster{})

	if bodyData, err := json.Marshal(createOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, errors.New(err)
	}

	r, err := errors.CoupleAPIErrors(req.
		SetBody(body).
		Post(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*LKECluster), nil
}

// UpdateLKECluster updates the LKECluster with the specified id
func (c *Client) UpdateLKECluster(ctx context.Context, id int, updateOpts LKEClusterUpdateOptions) (*LKECluster, error) {
	var body string
	e, err := c.LKEClusters.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)

	req := c.R(ctx).SetResult(&LKECluster{})

	if bodyData, err := json.Marshal(updateOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, errors.New(err)
	}

	r, err := errors.CoupleAPIErrors(req.
		SetBody(body).
		Put(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*LKECluster), nil
}

// DeleteLKECluster deletes the LKECluster with the specified id
func (c *Client) DeleteLKECluster(ctx context.Context, id int) error {
	e, err := c.LKEClusters.Endpoint()
	if err != nil {
		return err
	}
	e = fmt.Sprintf("%s/%d", e, id)

	_, err = errors.CoupleAPIErrors(c.R(ctx).Delete(e))
	return err
}

// ListLKEClusterAPIEndpoints gets the API Endpoint for the LKE Cluster specified
func (c *Client) ListLKEClusterAPIEndpoints(ctx context.Context, clusterID int, opts *ListOptions) ([]LKEClusterAPIEndpoint, error) {
	response := LKEClusterAPIEndpointsPagedResponse{}
	err := c.listHelperWithID(ctx, &response, clusterID, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetLKEClusterKubeconfig gets the Kubeconfig for the LKE Cluster specified
func (c *Client) GetLKEClusterKubeconfig(ctx context.Context, id int) (*LKEClusterKubeconfig, error) {
	e, err := c.LKEClusters.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d/kubeconfig", e, id)
	r, err := errors.CoupleAPIErrors(c.R(ctx).SetResult(&LKEClusterKubeconfig{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*LKEClusterKubeconfig), nil
}

// GetLKEVersion gets details about a specific LKE Version
func (c *Client) GetLKEVersion(ctx context.Context, version string) (*LKEVersion, error) {
	e, err := c.LKEVersions.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, version)
	r, err := errors.CoupleAPIErrors(c.R(ctx).SetResult(&LKEVersion{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*LKEVersion), nil
}

// ListLKEVersions lists the Kubernetes versions available through LKE
func (c *Client) ListLKEVersions(ctx context.Context, opts *ListOptions) ([]LKEVersion, error) {
	response := LKEVersionsPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
