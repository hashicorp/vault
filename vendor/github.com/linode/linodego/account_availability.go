package linodego

import (
	"context"
)

// AccountAvailability returns the resources availability in a region to an account.
type AccountAvailability struct {
	// region id
	Region string `json:"region"`

	// the unavailable resources in a region to the customer
	Unavailable []string `json:"unavailable"`

	// the available resources in a region to the customer
	Available []string `json:"available"`
}

// ListAccountAvailabilities lists all regions and the resource availabilities to the account.
func (c *Client) ListAccountAvailabilities(ctx context.Context, opts *ListOptions) ([]AccountAvailability, error) {
	response, err := getPaginatedResults[AccountAvailability](ctx, c, "account/availability", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetAccountAvailability gets the resources availability in a region to the customer.
func (c *Client) GetAccountAvailability(ctx context.Context, regionID string) (*AccountAvailability, error) {
	b := formatAPIPath("account/availability/%s", regionID)
	response, err := doGETRequest[AccountAvailability](ctx, c, b)
	if err != nil {
		return nil, err
	}

	return response, nil
}
