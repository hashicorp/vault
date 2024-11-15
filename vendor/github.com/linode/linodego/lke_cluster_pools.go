package linodego

import (
	"context"
)

// Deprecated: LKEClusterPoolDisk represents a Node disk in an LKEClusterPool object
type LKEClusterPoolDisk = LKENodePoolDisk

// Deprecated: LKEClusterPoolAutoscaler represents an AutoScaler configuration
type LKEClusterPoolAutoscaler = LKENodePoolAutoscaler

// Deprecated: LKEClusterPoolLinode represents a LKEClusterPoolLinode object
type LKEClusterPoolLinode = LKENodePoolLinode

// Deprecated: LKEClusterPool represents a LKEClusterPool object
type LKEClusterPool = LKENodePool

// Deprecated: LKEClusterPoolCreateOptions fields are those accepted by CreateLKEClusterPool
type LKEClusterPoolCreateOptions = LKENodePoolCreateOptions

// Deprecated: LKEClusterPoolUpdateOptions fields are those accepted by UpdateLKEClusterPool
type LKEClusterPoolUpdateOptions = LKENodePoolUpdateOptions

// Deprecated: ListLKEClusterPools lists LKEClusterPools
func (c *Client) ListLKEClusterPools(ctx context.Context, clusterID int, opts *ListOptions) ([]LKEClusterPool, error) {
	return c.ListLKENodePools(ctx, clusterID, opts)
}

// Deprecated: GetLKEClusterPool gets the lkeClusterPool with the provided ID
func (c *Client) GetLKEClusterPool(ctx context.Context, clusterID, id int) (*LKEClusterPool, error) {
	return c.GetLKENodePool(ctx, clusterID, id)
}

// Deprecated: CreateLKEClusterPool creates a LKEClusterPool
func (c *Client) CreateLKEClusterPool(ctx context.Context, clusterID int, createOpts LKEClusterPoolCreateOptions) (*LKEClusterPool, error) {
	return c.CreateLKENodePool(ctx, clusterID, createOpts)
}

// Deprecated: UpdateLKEClusterPool updates the LKEClusterPool with the specified id
func (c *Client) UpdateLKEClusterPool(ctx context.Context, clusterID, id int, updateOpts LKEClusterPoolUpdateOptions) (*LKEClusterPool, error) {
	return c.UpdateLKENodePool(ctx, clusterID, id, updateOpts)
}

// Deprecated: DeleteLKEClusterPool deletes the LKEClusterPool with the specified id
func (c *Client) DeleteLKEClusterPool(ctx context.Context, clusterID, id int) error {
	return c.DeleteLKENodePool(ctx, clusterID, id)
}

// Deprecated: DeleteLKEClusterPoolNode deletes a given node from a cluster pool
func (c *Client) DeleteLKEClusterPoolNode(ctx context.Context, clusterID int, id string) error {
	return c.DeleteLKENodePoolNode(ctx, clusterID, id)
}
