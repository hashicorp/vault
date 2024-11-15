package linodego

import (
	"context"
)

// LKEType represents a single valid LKE type.
// NOTE: This typically corresponds to the availability of a cluster's
// control plane.
type LKEType struct {
	baseType[LKETypePrice, LKETypeRegionPrice]
}

// LKETypePrice represents the base hourly and monthly prices
// for an LKE type entry.
type LKETypePrice struct {
	baseTypePrice
}

// LKETypeRegionPrice represents the regional hourly and monthly prices
// for an LKE type entry.
type LKETypeRegionPrice struct {
	baseTypeRegionPrice
}

// ListLKETypes lists LKE types. This endpoint is cached by default.
func (c *Client) ListLKETypes(ctx context.Context, opts *ListOptions) ([]LKEType, error) {
	e := "lke/types"

	endpoint, err := generateListCacheURL(e, opts)
	if err != nil {
		return nil, err
	}

	if result := c.getCachedResponse(endpoint); result != nil {
		return result.([]LKEType), nil
	}

	response, err := getPaginatedResults[LKEType](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(endpoint, response, &cacheExpiryTime)

	return response, nil
}
