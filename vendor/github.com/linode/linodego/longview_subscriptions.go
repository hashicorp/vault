package linodego

import (
	"context"
	"fmt"
)

// LongviewSubscription represents a LongviewSubscription object
type LongviewSubscription struct {
	ID              string       `json:"id"`
	Label           string       `json:"label"`
	ClientsIncluded int          `json:"clients_included"`
	Price           *LinodePrice `json:"price"`
	// UpdatedStr string `json:"updated"`
	// Updated *time.Time `json:"-"`
}

// LongviewSubscriptionsPagedResponse represents a paginated LongviewSubscription API response
type LongviewSubscriptionsPagedResponse struct {
	*PageOptions
	Data []LongviewSubscription `json:"data"`
}

// endpoint gets the endpoint URL for LongviewSubscription
func (LongviewSubscriptionsPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.LongviewSubscriptions.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends LongviewSubscriptions when processing paginated LongviewSubscription responses
func (resp *LongviewSubscriptionsPagedResponse) appendData(r *LongviewSubscriptionsPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListLongviewSubscriptions lists LongviewSubscriptions
func (c *Client) ListLongviewSubscriptions(ctx context.Context, opts *ListOptions) ([]LongviewSubscription, error) {
	response := LongviewSubscriptionsPagedResponse{}
	err := c.listHelper(ctx, &response, opts)

	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetLongviewSubscription gets the template with the provided ID
func (c *Client) GetLongviewSubscription(ctx context.Context, id string) (*LongviewSubscription, error) {
	e, err := c.LongviewSubscriptions.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := c.R(ctx).SetResult(&LongviewSubscription{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*LongviewSubscription), nil
}
