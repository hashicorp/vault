package linodego

import (
	"context"
	"net/url"
)

// LinodeType represents a linode type object
type LinodeType struct {
	ID           string              `json:"id"`
	Disk         int                 `json:"disk"`
	Class        LinodeTypeClass     `json:"class"` // enum: nanode, standard, highmem, dedicated, gpu
	Price        *LinodePrice        `json:"price"`
	Label        string              `json:"label"`
	Addons       *LinodeAddons       `json:"addons"`
	RegionPrices []LinodeRegionPrice `json:"region_prices"`
	NetworkOut   int                 `json:"network_out"`
	Memory       int                 `json:"memory"`
	Transfer     int                 `json:"transfer"`
	VCPUs        int                 `json:"vcpus"`
	GPUs         int                 `json:"gpus"`
	Successor    string              `json:"successor"`
}

// LinodePrice represents a linode type price object
type LinodePrice struct {
	Hourly  float32 `json:"hourly"`
	Monthly float32 `json:"monthly"`
}

// LinodeBackupsAddon represents a linode backups addon object
type LinodeBackupsAddon struct {
	Price        *LinodePrice        `json:"price"`
	RegionPrices []LinodeRegionPrice `json:"region_prices"`
}

// LinodeAddons represent the linode addons object
type LinodeAddons struct {
	Backups *LinodeBackupsAddon `json:"backups"`
}

// LinodeRegionPrice represents an individual type or addon
// price exception for a region.
type LinodeRegionPrice struct {
	ID      string  `json:"id"`
	Hourly  float32 `json:"hourly"`
	Monthly float32 `json:"monthly"`
}

// LinodeTypeClass constants start with Class and include Linode API Instance Type Classes
type LinodeTypeClass string

// LinodeTypeClass contants are the Instance Type Classes that an Instance Type can be assigned
const (
	ClassNanode    LinodeTypeClass = "nanode"
	ClassStandard  LinodeTypeClass = "standard"
	ClassHighmem   LinodeTypeClass = "highmem"
	ClassDedicated LinodeTypeClass = "dedicated"
	ClassGPU       LinodeTypeClass = "gpu"
)

// ListTypes lists linode types. This endpoint is cached by default.
func (c *Client) ListTypes(ctx context.Context, opts *ListOptions) ([]LinodeType, error) {
	e := "linode/types"

	endpoint, err := generateListCacheURL(e, opts)
	if err != nil {
		return nil, err
	}

	if result := c.getCachedResponse(endpoint); result != nil {
		return result.([]LinodeType), nil
	}

	response, err := getPaginatedResults[LinodeType](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(endpoint, response, &cacheExpiryTime)

	return response, nil
}

// GetType gets the type with the provided ID. This endpoint is cached by default.
func (c *Client) GetType(ctx context.Context, typeID string) (*LinodeType, error) {
	e := formatAPIPath("linode/types/%s", url.PathEscape(typeID))

	if result := c.getCachedResponse(e); result != nil {
		result := result.(LinodeType)
		return &result, nil
	}

	response, err := doGETRequest[LinodeType](ctx, c, e)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(e, response, &cacheExpiryTime)

	return response, nil
}
