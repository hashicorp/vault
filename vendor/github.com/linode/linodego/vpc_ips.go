package linodego

import (
	"context"
	"fmt"
)

// ListAllVPCIPAddresses gets the list of all IP addresses of all VPCs in the Linode account.
func (c *Client) ListAllVPCIPAddresses(
	ctx context.Context, opts *ListOptions,
) ([]VPCIP, error) {
	return getPaginatedResults[VPCIP](ctx, c, "vpcs/ips", opts)
}

// ListVPCIPAddresses gets the list of all IP addresses of a specific VPC.
func (c *Client) ListVPCIPAddresses(
	ctx context.Context, vpcID int, opts *ListOptions,
) ([]VPCIP, error) {
	return getPaginatedResults[VPCIP](ctx, c, fmt.Sprintf("vpcs/%d/ips", vpcID), opts)
}
