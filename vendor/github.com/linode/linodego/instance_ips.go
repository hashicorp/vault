package linodego

import (
	"context"
)

// InstanceIPAddressResponse contains the IPv4 and IPv6 details for an Instance
type InstanceIPAddressResponse struct {
	IPv4 *InstanceIPv4Response `json:"ipv4"`
	IPv6 *InstanceIPv6Response `json:"ipv6"`
}

// InstanceIPv4Response contains the details of all IPv4 addresses associated with an Instance
type InstanceIPv4Response struct {
	Public   []*InstanceIP `json:"public"`
	Private  []*InstanceIP `json:"private"`
	Shared   []*InstanceIP `json:"shared"`
	Reserved []*InstanceIP `json:"reserved"`
	VPC      []*VPCIP      `json:"vpc"`
}

// InstanceIP represents an Instance IP with additional DNS and networking details
type InstanceIP struct {
	Address    string             `json:"address"`
	Gateway    string             `json:"gateway"`
	SubnetMask string             `json:"subnet_mask"`
	Prefix     int                `json:"prefix"`
	Type       InstanceIPType     `json:"type"`
	Public     bool               `json:"public"`
	RDNS       string             `json:"rdns"`
	LinodeID   int                `json:"linode_id"`
	Region     string             `json:"region"`
	VPCNAT1To1 *InstanceIPNAT1To1 `json:"vpc_nat_1_1"`
	Reserved   bool               `json:"reserved"`
}

// VPCIP represents a private IP address in a VPC subnet with additional networking details
type VPCIP struct {
	Address      *string `json:"address"`
	AddressRange *string `json:"address_range"`
	Gateway      string  `json:"gateway"`
	SubnetMask   string  `json:"subnet_mask"`
	Prefix       int     `json:"prefix"`
	LinodeID     int     `json:"linode_id"`
	Region       string  `json:"region"`
	Active       bool    `json:"active"`
	NAT1To1      *string `json:"nat_1_1"`
	VPCID        int     `json:"vpc_id"`
	SubnetID     int     `json:"subnet_id"`
	ConfigID     int     `json:"config_id"`
	InterfaceID  int     `json:"interface_id"`
}

// InstanceIPv6Response contains the IPv6 addresses and ranges for an Instance
type InstanceIPv6Response struct {
	LinkLocal *InstanceIP `json:"link_local"`
	SLAAC     *InstanceIP `json:"slaac"`
	Global    []IPv6Range `json:"global"`
}

// InstanceIPNAT1To1 contains information about the NAT 1:1 mapping
// of a public IP address to a VPC subnet.
type InstanceIPNAT1To1 struct {
	Address  string `json:"address"`
	SubnetID int    `json:"subnet_id"`
	VPCID    int    `json:"vpc_id"`
}

// IPv6Range represents a range of IPv6 addresses routed to a single Linode in a given Region
type IPv6Range struct {
	Range  string `json:"range"`
	Region string `json:"region"`
	Prefix int    `json:"prefix"`

	RouteTarget string `json:"route_target"`

	// These fields are only returned by GetIPv6Range(...)
	IsBGP   bool  `json:"is_bgp"`
	Linodes []int `json:"linodes"`
}

type InstanceReserveIPOptions struct {
	Type    string `json:"type"`
	Public  bool   `json:"public"`
	Address string `json:"address"`
}

// InstanceIPType constants start with IPType and include Linode Instance IP Types
type InstanceIPType string

// InstanceIPType constants represent the IP types an Instance IP may be
const (
	IPTypeIPv4      InstanceIPType = "ipv4"
	IPTypeIPv6      InstanceIPType = "ipv6"
	IPTypeIPv6Pool  InstanceIPType = "ipv6/pool"
	IPTypeIPv6Range InstanceIPType = "ipv6/range"
)

// GetInstanceIPAddresses gets the IPAddresses for a Linode instance
func (c *Client) GetInstanceIPAddresses(ctx context.Context, linodeID int) (*InstanceIPAddressResponse, error) {
	e := formatAPIPath("linode/instances/%d/ips", linodeID)
	response, err := doGETRequest[InstanceIPAddressResponse](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetInstanceIPAddress gets the IPAddress for a Linode instance matching a supplied IP address
func (c *Client) GetInstanceIPAddress(ctx context.Context, linodeID int, ipaddress string) (*InstanceIP, error) {
	e := formatAPIPath("linode/instances/%d/ips/%s", linodeID, ipaddress)
	response, err := doGETRequest[InstanceIP](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// AddInstanceIPAddress adds a public or private IP to a Linode instance
func (c *Client) AddInstanceIPAddress(ctx context.Context, linodeID int, public bool) (*InstanceIP, error) {
	instanceipRequest := struct {
		Type   string `json:"type"`
		Public bool   `json:"public"`
	}{"ipv4", public}

	e := formatAPIPath("linode/instances/%d/ips", linodeID)
	response, err := doPOSTRequest[InstanceIP](ctx, c, e, instanceipRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateInstanceIPAddress updates the IPAddress with the specified instance id and IP address
func (c *Client) UpdateInstanceIPAddress(ctx context.Context, linodeID int, ipAddress string, opts IPAddressUpdateOptions) (*InstanceIP, error) {
	e := formatAPIPath("linode/instances/%d/ips/%s", linodeID, ipAddress)
	response, err := doPUTRequest[InstanceIP](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) DeleteInstanceIPAddress(ctx context.Context, linodeID int, ipAddress string) error {
	e := formatAPIPath("linode/instances/%d/ips/%s", linodeID, ipAddress)
	err := doDELETERequest(ctx, c, e)
	return err
}

// Function to add additional reserved IPV4 addresses to an existing linode
func (c *Client) AssignInstanceReservedIP(ctx context.Context, linodeID int, opts InstanceReserveIPOptions) (*InstanceIP, error) {
	endpoint := formatAPIPath("linode/instances/%d/ips", linodeID)
	response, err := doPOSTRequest[InstanceIP](ctx, c, endpoint, opts)
	if err != nil {
		return nil, err
	}
	return response, nil
}
