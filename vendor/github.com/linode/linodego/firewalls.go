package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
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

// ListFirewalls returns a paginated list of Cloud Firewalls
func (c *Client) ListFirewalls(ctx context.Context, opts *ListOptions) ([]Firewall, error) {
	response, err := getPaginatedResults[Firewall](ctx, c, "networking/firewalls", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateFirewall creates a single Firewall with at least one set of inbound or outbound rules
func (c *Client) CreateFirewall(ctx context.Context, opts FirewallCreateOptions) (*Firewall, error) {
	e := "networking/firewalls"
	response, err := doPOSTRequest[Firewall](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetFirewall gets a single Firewall with the provided ID
func (c *Client) GetFirewall(ctx context.Context, firewallID int) (*Firewall, error) {
	e := formatAPIPath("networking/firewalls/%d", firewallID)
	response, err := doGETRequest[Firewall](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateFirewall updates a Firewall with the given ID
func (c *Client) UpdateFirewall(ctx context.Context, firewallID int, opts FirewallUpdateOptions) (*Firewall, error) {
	e := formatAPIPath("networking/firewalls/%d", firewallID)
	response, err := doPUTRequest[Firewall](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteFirewall deletes a single Firewall with the provided ID
func (c *Client) DeleteFirewall(ctx context.Context, firewallID int) error {
	e := formatAPIPath("networking/firewalls/%d", firewallID)
	err := doDELETERequest(ctx, c, e)
	return err
}
