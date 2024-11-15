package linodego

import (
	"context"
)

// VolumeType represents a single valid Volume type.
type VolumeType struct {
	baseType[VolumeTypePrice, VolumeTypeRegionPrice]
}

// VolumeTypePrice represents the base hourly and monthly prices
// for a volume type entry.
type VolumeTypePrice struct {
	baseTypePrice
}

// VolumeTypeRegionPrice represents the regional hourly and monthly prices
// for a volume type entry.
type VolumeTypeRegionPrice struct {
	baseTypeRegionPrice
}

// ListVolumeTypes lists Volume types. This endpoint is cached by default.
func (c *Client) ListVolumeTypes(ctx context.Context, opts *ListOptions) ([]VolumeType, error) {
	e := "volumes/types"

	endpoint, err := generateListCacheURL(e, opts)
	if err != nil {
		return nil, err
	}

	if result := c.getCachedResponse(endpoint); result != nil {
		return result.([]VolumeType), nil
	}

	response, err := getPaginatedResults[VolumeType](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(endpoint, response, &cacheExpiryTime)

	return response, nil
}
