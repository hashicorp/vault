package linodego

import (
	"context"
	"fmt"

	"github.com/linode/linodego/pkg/errors"
)

// ObjectStorageCluster represents a linode object storage cluster object
type ObjectStorageCluster struct {
	ID               string `json:"id"`
	Domain           string `json:"domain"`
	Status           string `json:"status"`
	Region           string `json:"region"`
	StaticSiteDomain string `json:"static_site_domain"`
}

// ObjectStorageClustersPagedResponse represents a linode API response for listing
type ObjectStorageClustersPagedResponse struct {
	*PageOptions
	Data []ObjectStorageCluster `json:"data"`
}

// endpoint gets the endpoint URL for ObjectStorageCluster
func (ObjectStorageClustersPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.ObjectStorageClusters.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends ObjectStorageClusters when processing paginated ObjectStorageCluster responses
func (resp *ObjectStorageClustersPagedResponse) appendData(r *ObjectStorageClustersPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListObjectStorageClusters lists ObjectStorageClusters
func (c *Client) ListObjectStorageClusters(ctx context.Context, opts *ListOptions) ([]ObjectStorageCluster, error) {
	response := ObjectStorageClustersPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetObjectStorageCluster gets the template with the provided ID
func (c *Client) GetObjectStorageCluster(ctx context.Context, id string) (*ObjectStorageCluster, error) {
	e, err := c.ObjectStorageClusters.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := errors.CoupleAPIErrors(c.R(ctx).SetResult(&ObjectStorageCluster{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*ObjectStorageCluster), nil
}
