package linodego

import (
	"context"
	"time"
)

// This is an enumeration of Capabilities Linode offers that can be referenced
// through the user-facing parts of the application.
// Defined as strings rather than a custom type to avoid breaking change.
// Can be changed in the potential v2 version.
const (
	CapabilityLinodes                string = "Linodes"
	CapabilityNodeBalancers          string = "NodeBalancers"
	CapabilityBlockStorage           string = "Block Storage"
	CapabilityObjectStorage          string = "Object Storage"
	CapabilityObjectStorageRegions   string = "Object Storage Access Key Regions"
	CapabilityLKE                    string = "Kubernetes"
	CapabilityLkeHaControlPlanes     string = "LKE HA Control Planes"
	CapabilityCloudFirewall          string = "Cloud Firewall"
	CapabilityGPU                    string = "GPU Linodes"
	CapabilityVlans                  string = "Vlans"
	CapabilityVPCs                   string = "VPCs"
	CapabilityVPCsExtra              string = "VPCs Extra"
	CapabilityMachineImages          string = "Machine Images"
	CapabilityBareMetal              string = "Bare Metal"
	CapabilityDBAAS                  string = "Managed Databases"
	CapabilityBlockStorageMigrations string = "Block Storage Migrations"
	CapabilityMetadata               string = "Metadata"
	CapabilityPremiumPlans           string = "Premium Plans"
	CapabilityEdgePlans              string = "Edge Plans"
	CapabilityLKEControlPlaneACL     string = "LKE Network Access Control List (IP ACL)"
	CapabilityACLB                   string = "Akamai Cloud Load Balancer"
	CapabilitySupportTicketSeverity  string = "Support Ticket Severity"
	CapabilityBackups                string = "Backups"
	CapabilityPlacementGroup         string = "Placement Group"
	CapabilityDiskEncryption         string = "Disk Encryption"
	CapabilityBlockStorageEncryption string = "Block Storage Encryption"
)

// Region-related endpoints have a custom expiry time as the
// `status` field may update for database outages.
var cacheExpiryTime = time.Minute

// Region represents a linode region object
type Region struct {
	ID      string `json:"id"`
	Country string `json:"country"`

	// A List of enums from the above constants
	Capabilities []string `json:"capabilities"`

	Status   string `json:"status"`
	Label    string `json:"label"`
	SiteType string `json:"site_type"`

	Resolvers            RegionResolvers             `json:"resolvers"`
	PlacementGroupLimits *RegionPlacementGroupLimits `json:"placement_group_limits"`
}

// RegionResolvers contains the DNS resolvers of a region
type RegionResolvers struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

// RegionPlacementGroupLimits contains information about the
// placement group limits for the current user in the current region.
type RegionPlacementGroupLimits struct {
	MaximumPGsPerCustomer int `json:"maximum_pgs_per_customer"`
	MaximumLinodesPerPG   int `json:"maximum_linodes_per_pg"`
}

// ListRegions lists Regions. This endpoint is cached by default.
func (c *Client) ListRegions(ctx context.Context, opts *ListOptions) ([]Region, error) {
	endpoint, err := generateListCacheURL("regions", opts)
	if err != nil {
		return nil, err
	}

	if result := c.getCachedResponse(endpoint); result != nil {
		return result.([]Region), nil
	}

	response, err := getPaginatedResults[Region](ctx, c, "regions", opts)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(endpoint, response, &cacheExpiryTime)

	return response, nil
}

// GetRegion gets the template with the provided ID. This endpoint is cached by default.
func (c *Client) GetRegion(ctx context.Context, regionID string) (*Region, error) {
	e := formatAPIPath("regions/%s", regionID)

	if result := c.getCachedResponse(e); result != nil {
		result := result.(Region)
		return &result, nil
	}

	response, err := doGETRequest[Region](ctx, c, e)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(e, response, &cacheExpiryTime)

	return response, nil
}
