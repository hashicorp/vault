package linodego

import (
	"context"
)

// Region represents a linode region object
type RegionAvailability struct {
	Region    string `json:"region"`
	Plan      string `json:"plan"`
	Available bool   `json:"available"`
}

// ListRegionsAvailability lists Regions. This endpoint is cached by default.
func (c *Client) ListRegionsAvailability(ctx context.Context, opts *ListOptions) ([]RegionAvailability, error) {
	e := "regions/availability"

	endpoint, err := generateListCacheURL(e, opts)
	if err != nil {
		return nil, err
	}

	if result := c.getCachedResponse(endpoint); result != nil {
		return result.([]RegionAvailability), nil
	}

	response, err := getPaginatedResults[RegionAvailability](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(endpoint, response, &cacheExpiryTime)

	return response, nil
}

// GetRegionAvailability gets the template with the provided ID. This endpoint is cached by default.
func (c *Client) GetRegionAvailability(ctx context.Context, regionID string) (*RegionAvailability, error) {
	e := formatAPIPath("regions/%s/availability", regionID)

	if result := c.getCachedResponse(e); result != nil {
		result := result.(RegionAvailability)
		return &result, nil
	}

	response, err := doGETRequest[RegionAvailability](ctx, c, e)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(e, response, &cacheExpiryTime)

	return response, nil
}
