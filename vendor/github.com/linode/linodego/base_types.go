// This package contains various type-related base classes intended
// to be used in composition across type structures in this project.

package linodego

// baseType is a base struct containing the core fields of a resource type
// returned from the Linode API.
type baseType[PriceType any, RegionPriceType any] struct {
	ID           string            `json:"id"`
	Label        string            `json:"label"`
	Price        PriceType         `json:"price"`
	RegionPrices []RegionPriceType `json:"region_prices"`
	Transfer     int               `json:"transfer"`
}

// baseTypePrice is a base struct containing the core fields of a resource type's
// base price.
type baseTypePrice struct {
	Hourly  float64 `json:"hourly"`
	Monthly float64 `json:"monthly"`
}

// baseTypeRegionPrice is a base struct containing the core fields of a resource type's
// region-specific price.
type baseTypeRegionPrice struct {
	baseTypePrice

	ID string `json:"id"`
}
