package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// LongviewClient represents a LongviewClient object
type LongviewClient struct {
	ID          int        `json:"id"`
	APIKey      string     `json:"api_key"`
	Created     *time.Time `json:"-"`
	InstallCode string     `json:"install_code"`
	Label       string     `json:"label"`
	Updated     *time.Time `json:"-"`
	Apps        struct {
		Apache any `json:"apache"`
		MySQL  any `json:"mysql"`
		NginX  any `json:"nginx"`
	} `json:"apps"`
}

// LongviewClientCreateOptions is an options struct used when Creating a Longview Client
type LongviewClientCreateOptions struct {
	Label string `json:"label"`
}

// LongviewClientCreateOptions is an options struct used when Updating a Longview Client
type LongviewClientUpdateOptions struct {
	Label string `json:"label"`
}

// LongviewPlan represents a Longview Plan object
type LongviewPlan struct {
	ID              string `json:"id"`
	Label           string `json:"label"`
	ClientsIncluded int    `json:"clients_included"`
	Price           struct {
		Hourly  float64 `json:"hourly"`
		Monthly float64 `json:"monthly"`
	} `json:"price"`
}

// LongviewPlanUpdateOptions is an options struct used when Updating a Longview Plan
type LongviewPlanUpdateOptions struct {
	LongviewSubscription string `json:"longview_subscription"`
}

// ListLongviewClients lists LongviewClients
func (c *Client) ListLongviewClients(ctx context.Context, opts *ListOptions) ([]LongviewClient, error) {
	response, err := getPaginatedResults[LongviewClient](ctx, c, "longview/clients", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetLongviewClient gets the template with the provided ID
func (c *Client) GetLongviewClient(ctx context.Context, clientID int) (*LongviewClient, error) {
	e := formatAPIPath("longview/clients/%d", clientID)
	response, err := doGETRequest[LongviewClient](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateLongviewClient creates a Longview Client
func (c *Client) CreateLongviewClient(ctx context.Context, opts LongviewClientCreateOptions) (*LongviewClient, error) {
	e := "longview/clients"
	response, err := doPOSTRequest[LongviewClient](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteLongviewClient deletes a Longview Client
func (c *Client) DeleteLongviewClient(ctx context.Context, clientID int) error {
	e := formatAPIPath("longview/clients/%d", clientID)
	err := doDELETERequest(ctx, c, e)
	return err
}

// UpdateLongviewClient updates a Longview Client
func (c *Client) UpdateLongviewClient(ctx context.Context, clientID int, opts LongviewClientUpdateOptions) (*LongviewClient, error) {
	e := formatAPIPath("longview/clients/%d", clientID)
	response, err := doPUTRequest[LongviewClient](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetLongviewPlan gets the template with the provided ID
func (c *Client) GetLongviewPlan(ctx context.Context) (*LongviewPlan, error) {
	e := "longview/plan"
	response, err := doGETRequest[LongviewPlan](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateLongviewPlan updates a Longview Plan
func (c *Client) UpdateLongviewPlan(ctx context.Context, opts LongviewPlanUpdateOptions) (*LongviewPlan, error) {
	e := "longview/plan"
	response, err := doPUTRequest[LongviewPlan](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *LongviewClient) UnmarshalJSON(b []byte) error {
	type Mask LongviewClient

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.Updated = (*time.Time)(p.Updated)

	return nil
}
