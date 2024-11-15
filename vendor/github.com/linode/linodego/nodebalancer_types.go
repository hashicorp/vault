package linodego

import (
	"context"
)

// NodeBalancerType represents a single valid NodeBalancer type.
type NodeBalancerType struct {
	baseType[NodeBalancerTypePrice, NodeBalancerTypeRegionPrice]
}

// NodeBalancerTypePrice represents the base hourly and monthly prices
// for a NodeBalancer type entry.
type NodeBalancerTypePrice struct {
	baseTypePrice
}

// NodeBalancerTypeRegionPrice represents the regional hourly and monthly prices
// for a NodeBalancer type entry.
type NodeBalancerTypeRegionPrice struct {
	baseTypeRegionPrice
}

// ListNodeBalancerTypes lists NodeBalancer types. This endpoint is cached by default.
func (c *Client) ListNodeBalancerTypes(ctx context.Context, opts *ListOptions) ([]NodeBalancerType, error) {
	e := "nodebalancers/types"

	endpoint, err := generateListCacheURL(e, opts)
	if err != nil {
		return nil, err
	}

	if result := c.getCachedResponse(endpoint); result != nil {
		return result.([]NodeBalancerType), nil
	}

	response, err := getPaginatedResults[NodeBalancerType](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(endpoint, response, &cacheExpiryTime)

	return response, nil
}
