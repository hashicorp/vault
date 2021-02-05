package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
	"github.com/linode/linodego/pkg/errors"
)

// FirewallStatus enum type
type FirewallStatus string

// FirewallStatus enums start with Firewall
const (
	FirewallEnabled  FirewallStatus = "enabled"
	FirewallDisabled FirewallStatus = "disabled"
	FirewallDeleted  FirewallStatus = "deleted"
)

// A Firewall is a set of networking rules (iptables) applied to Devices with which it is associated
type Firewall struct {
	ID      int             `json:"id"`
	Label   string          `json:"label"`
	Status  FirewallStatus  `json:"status"`
	Tags    []string        `json:"tags,omitempty"`
	Rules   FirewallRuleSet `json:"rules"`
	Created *time.Time      `json:"-"`
	Updated *time.Time      `json:"-"`
}

// DevicesCreationOptions fields are used when adding devices during the Firewall creation process.
type DevicesCreationOptions struct {
	Linodes       []int `json:"linodes,omitempty"`
	NodeBalancers []int `json:"nodebalancers,omitempty"`
}

// FirewallCreateOptions fields are those accepted by CreateFirewall
type FirewallCreateOptions struct {
	Label   string                 `json:"label,omitempty"`
	Rules   FirewallRuleSet        `json:"rules"`
	Tags    []string               `json:"tags,omitempty"`
	Devices DevicesCreationOptions `json:"devices,omitempty"`
}

// FirewallUpdateOptions is an options struct used when Updating a Firewall
type FirewallUpdateOptions struct {
	Label  string         `json:"label,omitempty"`
	Status FirewallStatus `json:"status,omitempty"`
	Tags   *[]string      `json:"tags,omitempty"`
}

// GetUpdateOptions converts a Firewall to FirewallUpdateOptions for use in Client.UpdateFirewall.
func (f *Firewall) GetUpdateOptions() FirewallUpdateOptions {
	return FirewallUpdateOptions{
		Label:  f.Label,
		Status: f.Status,
		Tags:   &f.Tags,
	}
}

// UnmarshalJSON for Firewall responses
func (f *Firewall) UnmarshalJSON(b []byte) error {
	type Mask Firewall

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(f),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	f.Created = (*time.Time)(p.Created)
	f.Updated = (*time.Time)(p.Updated)
	return nil
}

// FirewallsPagedResponse represents a Linode API response for listing of Cloud Firewalls
type FirewallsPagedResponse struct {
	*PageOptions
	Data []Firewall `json:"data"`
}

func (FirewallsPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.Firewalls.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (resp *FirewallsPagedResponse) appendData(r *FirewallsPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListFirewalls returns a paginated list of Cloud Firewalls
func (c *Client) ListFirewalls(ctx context.Context, opts *ListOptions) ([]Firewall, error) {
	response := FirewallsPagedResponse{}

	err := c.listHelper(ctx, &response, opts)

	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// CreateFirewall creates a single Firewall with at least one set of inbound or outbound rules
func (c *Client) CreateFirewall(ctx context.Context, createOpts FirewallCreateOptions) (*Firewall, error) {
	var body string
	e, err := c.Firewalls.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&Firewall{})

	if bodyData, err := json.Marshal(createOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, errors.New(err)
	}

	r, err := errors.CoupleAPIErrors(req.
		SetBody(body).
		Post(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*Firewall), nil
}

// GetFirewall gets a single Firewall with the provided ID
func (c *Client) GetFirewall(ctx context.Context, id int) (*Firewall, error) {
	e, err := c.Firewalls.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx)

	e = fmt.Sprintf("%s/%d", e, id)
	r, err := errors.CoupleAPIErrors(req.SetResult(&Firewall{}).Get(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*Firewall), nil
}

// UpdateFirewall updates a Firewall with the given ID
func (c *Client) UpdateFirewall(ctx context.Context, id int, updateOpts FirewallUpdateOptions) (*Firewall, error) {
	e, err := c.Firewalls.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&Firewall{})

	bodyData, err := json.Marshal(updateOpts)
	if err != nil {
		return nil, errors.New(err)
	}

	body := string(bodyData)

	e = fmt.Sprintf("%s/%d", e, id)
	r, err := errors.CoupleAPIErrors(req.SetBody(body).Put(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*Firewall), nil
}

// DeleteFirewall deletes a single Firewall with the provided ID
func (c *Client) DeleteFirewall(ctx context.Context, id int) error {
	e, err := c.Firewalls.Endpoint()
	if err != nil {
		return err
	}

	req := c.R(ctx)

	e = fmt.Sprintf("%s/%d", e, id)
	_, err = errors.CoupleAPIErrors(req.Delete(e))
	return err
}
