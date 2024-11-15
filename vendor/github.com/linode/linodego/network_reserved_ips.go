package linodego

import (
	"context"
)

// ReserveIPOptions represents the options for reserving an IP address
// NOTE: Reserved IP feature may not currently be available to all users.
type ReserveIPOptions struct {
	Region string `json:"region"`
}

// ListReservedIPAddresses retrieves a list of reserved IP addresses
// NOTE: Reserved IP feature may not currently be available to all users.
func (c *Client) ListReservedIPAddresses(ctx context.Context, opts *ListOptions) ([]InstanceIP, error) {
	e := formatAPIPath("networking/reserved/ips")
	response, err := getPaginatedResults[InstanceIP](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetReservedIPAddress retrieves details of a specific reserved IP address
// NOTE: Reserved IP feature may not currently be available to all users.
func (c *Client) GetReservedIPAddress(ctx context.Context, ipAddress string) (*InstanceIP, error) {
	e := formatAPIPath("networking/reserved/ips/%s", ipAddress)
	response, err := doGETRequest[InstanceIP](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ReserveIPAddress reserves a new IP address
// NOTE: Reserved IP feature may not currently be available to all users.
func (c *Client) ReserveIPAddress(ctx context.Context, opts ReserveIPOptions) (*InstanceIP, error) {
	e := "networking/reserved/ips"
	response, err := doPOSTRequest[InstanceIP](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteReservedIPAddress deletes a reserved IP address
// NOTE: Reserved IP feature may not currently be available to all users.
func (c *Client) DeleteReservedIPAddress(ctx context.Context, ipAddress string) error {
	e := formatAPIPath("networking/reserved/ips/%s", ipAddress)
	return doDELETERequest(ctx, c, e)
}
