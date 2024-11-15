package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type VPC struct {
	ID          int         `json:"id"`
	Label       string      `json:"label"`
	Description string      `json:"description"`
	Region      string      `json:"region"`
	Subnets     []VPCSubnet `json:"subnets"`
	Created     *time.Time  `json:"-"`
	Updated     *time.Time  `json:"-"`
}

type VPCCreateOptions struct {
	Label       string                   `json:"label"`
	Description string                   `json:"description,omitempty"`
	Region      string                   `json:"region"`
	Subnets     []VPCSubnetCreateOptions `json:"subnets,omitempty"`
}

type VPCUpdateOptions struct {
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
}

func (v VPC) GetCreateOptions() VPCCreateOptions {
	subnetCreations := make([]VPCSubnetCreateOptions, len(v.Subnets))
	for i, s := range v.Subnets {
		subnetCreations[i] = s.GetCreateOptions()
	}

	return VPCCreateOptions{
		Label:       v.Label,
		Description: v.Description,
		Region:      v.Region,
		Subnets:     subnetCreations,
	}
}

func (v VPC) GetUpdateOptions() VPCUpdateOptions {
	return VPCUpdateOptions{
		Label:       v.Label,
		Description: v.Description,
	}
}

func (v *VPC) UnmarshalJSON(b []byte) error {
	type Mask VPC
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

func (c *Client) CreateVPC(
	ctx context.Context,
	opts VPCCreateOptions,
) (*VPC, error) {
	e := "vpcs"
	response, err := doPOSTRequest[VPC](ctx, c, e, opts)
	return response, err
}

func (c *Client) GetVPC(ctx context.Context, vpcID int) (*VPC, error) {
	e := formatAPIPath("/vpcs/%d", vpcID)
	response, err := doGETRequest[VPC](ctx, c, e)
	return response, err
}

func (c *Client) ListVPCs(ctx context.Context, opts *ListOptions) ([]VPC, error) {
	response, err := getPaginatedResults[VPC](ctx, c, "vpcs", opts)
	return response, err
}

func (c *Client) UpdateVPC(
	ctx context.Context,
	vpcID int,
	opts VPCUpdateOptions,
) (*VPC, error) {
	e := formatAPIPath("vpcs/%d", vpcID)
	response, err := doPUTRequest[VPC](ctx, c, e, opts)
	return response, err
}

func (c *Client) DeleteVPC(ctx context.Context, vpcID int) error {
	e := formatAPIPath("vpcs/%d", vpcID)
	err := doDELETERequest(ctx, c, e)
	return err
}
