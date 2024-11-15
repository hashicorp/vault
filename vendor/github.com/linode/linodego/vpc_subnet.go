package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// VPCSubnetLinodeInterface represents an interface on a Linode that is currently
// assigned to this VPC subnet.
type VPCSubnetLinodeInterface struct {
	ID     int  `json:"id"`
	Active bool `json:"active"`
}

// VPCSubnetLinode represents a Linode currently assigned to a VPC subnet.
type VPCSubnetLinode struct {
	ID         int                        `json:"id"`
	Interfaces []VPCSubnetLinodeInterface `json:"interfaces"`
}

type VPCSubnet struct {
	ID      int               `json:"id"`
	Label   string            `json:"label"`
	IPv4    string            `json:"ipv4"`
	Linodes []VPCSubnetLinode `json:"linodes"`
	Created *time.Time        `json:"-"`
	Updated *time.Time        `json:"-"`
}

type VPCSubnetCreateOptions struct {
	Label string `json:"label"`
	IPv4  string `json:"ipv4"`
}

type VPCSubnetUpdateOptions struct {
	Label string `json:"label"`
}

func (v *VPCSubnet) UnmarshalJSON(b []byte) error {
	type Mask VPCSubnet
	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(v),
	}
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	v.Created = (*time.Time)(p.Created)
	v.Updated = (*time.Time)(p.Updated)

	return nil
}

func (v VPCSubnet) GetCreateOptions() VPCSubnetCreateOptions {
	return VPCSubnetCreateOptions{
		Label: v.Label,
		IPv4:  v.IPv4,
	}
}

func (v VPCSubnet) GetUpdateOptions() VPCSubnetUpdateOptions {
	return VPCSubnetUpdateOptions{Label: v.Label}
}

func (c *Client) CreateVPCSubnet(
	ctx context.Context,
	opts VPCSubnetCreateOptions,
	vpcID int,
) (*VPCSubnet, error) {
	e := formatAPIPath("vpcs/%d/subnets", vpcID)
	response, err := doPOSTRequest[VPCSubnet](ctx, c, e, opts)
	return response, err
}

func (c *Client) GetVPCSubnet(
	ctx context.Context,
	vpcID int,
	subnetID int,
) (*VPCSubnet, error) {
	e := formatAPIPath("vpcs/%d/subnets/%d", vpcID, subnetID)
	response, err := doGETRequest[VPCSubnet](ctx, c, e)
	return response, err
}

func (c *Client) ListVPCSubnets(
	ctx context.Context,
	vpcID int,
	opts *ListOptions,
) ([]VPCSubnet, error) {
	response, err := getPaginatedResults[VPCSubnet](ctx, c, formatAPIPath("vpcs/%d/subnets", vpcID), opts)
	return response, err
}

func (c *Client) UpdateVPCSubnet(
	ctx context.Context,
	vpcID int,
	subnetID int,
	opts VPCSubnetUpdateOptions,
) (*VPCSubnet, error) {
	e := formatAPIPath("vpcs/%d/subnets/%d", vpcID, subnetID)
	response, err := doPUTRequest[VPCSubnet](ctx, c, e, opts)
	return response, err
}

func (c *Client) DeleteVPCSubnet(ctx context.Context, vpcID int, subnetID int) error {
	e := formatAPIPath("vpcs/%d/subnets/%d", vpcID, subnetID)
	err := doDELETERequest(ctx, c, e)
	return err
}
