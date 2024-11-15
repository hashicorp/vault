package linodego

import (
	"context"
)

// Domain represents a Domain object
type Domain struct {
	//	This Domain's unique ID
	ID int `json:"id"`

	// The domain this Domain represents. These must be unique in our system; you cannot have two Domains representing the same domain.
	Domain string `json:"domain"`

	// If this Domain represents the authoritative source of information for the domain it describes, or if it is a read-only copy of a master (also called a slave).
	Type DomainType `json:"type"` // Enum:"master" "slave"

	// Deprecated: The group this Domain belongs to. This is for display purposes only.
	Group string `json:"group"`

	// Used to control whether this Domain is currently being rendered.
	Status DomainStatus `json:"status"` // Enum:"disabled" "active" "edit_mode" "has_errors"

	// A description for this Domain. This is for display purposes only.
	Description string `json:"description"`

	// Start of Authority email address. This is required for master Domains.
	SOAEmail string `json:"soa_email"`

	// The interval, in seconds, at which a failed refresh should be retried.
	// Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.
	RetrySec int `json:"retry_sec"`

	// The IP addresses representing the master DNS for this Domain.
	MasterIPs []string `json:"master_ips"`

	// The list of IPs that may perform a zone transfer for this Domain. This is potentially dangerous, and should be set to an empty list unless you intend to use it.
	AXfrIPs []string `json:"axfr_ips"`

	// An array of tags applied to this object. Tags are for organizational purposes only.
	Tags []string `json:"tags"`

	// The amount of time in seconds that may pass before this Domain is no longer authoritative. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.
	ExpireSec int `json:"expire_sec"`

	// The amount of time in seconds before this Domain should be refreshed. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.
	RefreshSec int `json:"refresh_sec"`

	// "Time to Live" - the amount of time in seconds that this Domain's records may be cached by resolvers or other domain servers. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.
	TTLSec int `json:"ttl_sec"`
}

// DomainZoneFile represents the Zone File of a Domain
type DomainZoneFile struct {
	ZoneFile []string `json:"zone_file"`
}

// DomainCreateOptions fields are those accepted by CreateDomain
type DomainCreateOptions struct {
	// The domain this Domain represents. These must be unique in our system; you cannot have two Domains representing the same domain.
	Domain string `json:"domain"`

	// If this Domain represents the authoritative source of information for the domain it describes, or if it is a read-only copy of a master (also called a slave).
	// Enum:"master" "slave"
	Type DomainType `json:"type"`

	// Deprecated: The group this Domain belongs to. This is for display purposes only.
	Group string `json:"group,omitempty"`

	// Used to control whether this Domain is currently being rendered.
	// Enum:"disabled" "active" "edit_mode" "has_errors"
	Status DomainStatus `json:"status,omitempty"`

	// A description for this Domain. This is for display purposes only.
	Description string `json:"description,omitempty"`

	// Start of Authority email address. This is required for master Domains.
	SOAEmail string `json:"soa_email,omitempty"`

	// The interval, in seconds, at which a failed refresh should be retried.
	// Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.
	RetrySec int `json:"retry_sec,omitempty"`

	// The IP addresses representing the master DNS for this Domain.
	MasterIPs []string `json:"master_ips"`

	// The list of IPs that may perform a zone transfer for this Domain. This is potentially dangerous, and should be set to an empty list unless you intend to use it.
	AXfrIPs []string `json:"axfr_ips"`

	// An array of tags applied to this object. Tags are for organizational purposes only.
	Tags []string `json:"tags"`

	// The amount of time in seconds that may pass before this Domain is no longer authoritative. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.
	ExpireSec int `json:"expire_sec,omitempty"`

	// The amount of time in seconds before this Domain should be refreshed. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.
	RefreshSec int `json:"refresh_sec,omitempty"`

	// "Time to Live" - the amount of time in seconds that this Domain's records may be cached by resolvers or other domain servers. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.
	TTLSec int `json:"ttl_sec,omitempty"`
}

// DomainUpdateOptions converts a Domain to DomainUpdateOptions for use in UpdateDomain
type DomainUpdateOptions struct {
	// The domain this Domain represents. These must be unique in our system; you cannot have two Domains representing the same domain.
	Domain string `json:"domain,omitempty"`

	// If this Domain represents the authoritative source of information for the domain it describes, or if it is a read-only copy of a master (also called a slave).
	// Enum:"master" "slave"
	Type DomainType `json:"type,omitempty"`

	// Deprecated: The group this Domain belongs to. This is for display purposes only.
	Group string `json:"group,omitempty"`

	// Used to control whether this Domain is currently being rendered.
	// Enum:"disabled" "active" "edit_mode" "has_errors"
	Status DomainStatus `json:"status,omitempty"`

	// A description for this Domain. This is for display purposes only.
	Description string `json:"description,omitempty"`

	// Start of Authority email address. This is required for master Domains.
	SOAEmail string `json:"soa_email,omitempty"`

	// The interval, in seconds, at which a failed refresh should be retried.
	// Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.
	RetrySec int `json:"retry_sec,omitempty"`

	// The IP addresses representing the master DNS for this Domain.
	MasterIPs []string `json:"master_ips"`

	// The list of IPs that may perform a zone transfer for this Domain. This is potentially dangerous, and should be set to an empty list unless you intend to use it.
	AXfrIPs []string `json:"axfr_ips"`

	// An array of tags applied to this object. Tags are for organizational purposes only.
	Tags []string `json:"tags"`

	// The amount of time in seconds that may pass before this Domain is no longer authoritative. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.
	ExpireSec int `json:"expire_sec,omitempty"`

	// The amount of time in seconds before this Domain should be refreshed. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.
	RefreshSec int `json:"refresh_sec,omitempty"`

	// "Time to Live" - the amount of time in seconds that this Domain's records may be cached by resolvers or other domain servers. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.
	TTLSec int `json:"ttl_sec,omitempty"`
}

// DomainType constants start with DomainType and include Linode API Domain Type values
type DomainType string

// DomainType constants reflect the DNS zone type of a Domain
const (
	DomainTypeMaster DomainType = "master"
	DomainTypeSlave  DomainType = "slave"
)

// DomainStatus constants start with DomainStatus and include Linode API Domain Status values
type DomainStatus string

// DomainStatus constants reflect the current status of a Domain
const (
	DomainStatusDisabled  DomainStatus = "disabled"
	DomainStatusActive    DomainStatus = "active"
	DomainStatusEditMode  DomainStatus = "edit_mode"
	DomainStatusHasErrors DomainStatus = "has_errors"
)

// GetUpdateOptions converts a Domain to DomainUpdateOptions for use in UpdateDomain
func (d Domain) GetUpdateOptions() (du DomainUpdateOptions) {
	du.Domain = d.Domain
	du.Type = d.Type
	du.Group = d.Group
	du.Status = d.Status
	du.Description = d.Description
	du.SOAEmail = d.SOAEmail
	du.RetrySec = d.RetrySec
	du.MasterIPs = d.MasterIPs
	du.AXfrIPs = d.AXfrIPs
	du.Tags = d.Tags
	du.ExpireSec = d.ExpireSec
	du.RefreshSec = d.RefreshSec
	du.TTLSec = d.TTLSec

	return
}

// ListDomains lists Domains
func (c *Client) ListDomains(ctx context.Context, opts *ListOptions) ([]Domain, error) {
	response, err := getPaginatedResults[Domain](ctx, c, "domains", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetDomain gets the domain with the provided ID
func (c *Client) GetDomain(ctx context.Context, domainID int) (*Domain, error) {
	e := formatAPIPath("domains/%d", domainID)
	response, err := doGETRequest[Domain](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateDomain creates a Domain
func (c *Client) CreateDomain(ctx context.Context, opts DomainCreateOptions) (*Domain, error) {
	e := "domains"
	response, err := doPOSTRequest[Domain](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateDomain updates the Domain with the specified id
func (c *Client) UpdateDomain(ctx context.Context, domainID int, opts DomainUpdateOptions) (*Domain, error) {
	e := formatAPIPath("domains/%d", domainID)
	response, err := doPUTRequest[Domain](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteDomain deletes the Domain with the specified id
func (c *Client) DeleteDomain(ctx context.Context, domainID int) error {
	e := formatAPIPath("domains/%d", domainID)
	err := doDELETERequest(ctx, c, e)
	return err
}

// GetDomainZoneFile gets the zone file for the last rendered zone for the specified domain.
func (c *Client) GetDomainZoneFile(ctx context.Context, domainID int) (*DomainZoneFile, error) {
	e := formatAPIPath("domains/%d/zone-file", domainID)
	response, err := doGETRequest[DomainZoneFile](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}
