package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// NodeBalancer represents a NodeBalancer object
type NodeBalancer struct {
	// This NodeBalancer's unique ID.
	ID int `json:"id"`
	// This NodeBalancer's label. These must be unique on your Account.
	Label *string `json:"label"`
	// The Region where this NodeBalancer is located. NodeBalancers only support backends in the same Region.
	Region string `json:"region"`
	// This NodeBalancer's hostname, ending with .nodebalancer.linode.com
	Hostname *string `json:"hostname"`
	// This NodeBalancer's public IPv4 address.
	IPv4 *string `json:"ipv4"`
	// This NodeBalancer's public IPv6 address.
	IPv6 *string `json:"ipv6"`
	// Throttle connections per second (0-20). Set to 0 (zero) to disable throttling.
	ClientConnThrottle int `json:"client_conn_throttle"`
	// Information about the amount of transfer this NodeBalancer has had so far this month.
	Transfer NodeBalancerTransfer `json:"transfer"`

	// An array of tags applied to this object. Tags are for organizational purposes only.
	Tags []string `json:"tags"`

	Created *time.Time `json:"-"`
	Updated *time.Time `json:"-"`
}

// NodeBalancerTransfer contains information about the amount of transfer a NodeBalancer has had in the current month
type NodeBalancerTransfer struct {
	// The total transfer, in MB, used by this NodeBalancer this month.
	Total *float64 `json:"total"`
	// The total inbound transfer, in MB, used for this NodeBalancer this month.
	Out *float64 `json:"out"`
	// The total outbound transfer, in MB, used for this NodeBalancer this month.
	In *float64 `json:"in"`
}

// NodeBalancerCreateOptions are the options permitted for CreateNodeBalancer
type NodeBalancerCreateOptions struct {
	Label              *string                            `json:"label,omitempty"`
	Region             string                             `json:"region,omitempty"`
	ClientConnThrottle *int                               `json:"client_conn_throttle,omitempty"`
	Configs            []*NodeBalancerConfigCreateOptions `json:"configs,omitempty"`
	Tags               []string                           `json:"tags"`
	FirewallID         int                                `json:"firewall_id,omitempty"`
}

// NodeBalancerUpdateOptions are the options permitted for UpdateNodeBalancer
type NodeBalancerUpdateOptions struct {
	Label              *string   `json:"label,omitempty"`
	ClientConnThrottle *int      `json:"client_conn_throttle,omitempty"`
	Tags               *[]string `json:"tags,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *NodeBalancer) UnmarshalJSON(b []byte) error {
	type Mask NodeBalancer

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

// GetCreateOptions converts a NodeBalancer to NodeBalancerCreateOptions for use in CreateNodeBalancer
func (i NodeBalancer) GetCreateOptions() NodeBalancerCreateOptions {
	return NodeBalancerCreateOptions{
		Label:              i.Label,
		Region:             i.Region,
		ClientConnThrottle: &i.ClientConnThrottle,
		Tags:               i.Tags,
	}
}

// GetUpdateOptions converts a NodeBalancer to NodeBalancerUpdateOptions for use in UpdateNodeBalancer
func (i NodeBalancer) GetUpdateOptions() NodeBalancerUpdateOptions {
	return NodeBalancerUpdateOptions{
		Label:              i.Label,
		ClientConnThrottle: &i.ClientConnThrottle,
		Tags:               &i.Tags,
	}
}

// ListNodeBalancers lists NodeBalancers
func (c *Client) ListNodeBalancers(ctx context.Context, opts *ListOptions) ([]NodeBalancer, error) {
	response, err := getPaginatedResults[NodeBalancer](ctx, c, "nodebalancers", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetNodeBalancer gets the NodeBalancer with the provided ID
func (c *Client) GetNodeBalancer(ctx context.Context, nodebalancerID int) (*NodeBalancer, error) {
	e := formatAPIPath("nodebalancers/%d", nodebalancerID)
	response, err := doGETRequest[NodeBalancer](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateNodeBalancer creates a NodeBalancer
func (c *Client) CreateNodeBalancer(ctx context.Context, opts NodeBalancerCreateOptions) (*NodeBalancer, error) {
	e := "nodebalancers"
	response, err := doPOSTRequest[NodeBalancer](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateNodeBalancer updates the NodeBalancer with the specified id
func (c *Client) UpdateNodeBalancer(ctx context.Context, nodebalancerID int, opts NodeBalancerUpdateOptions) (*NodeBalancer, error) {
	e := formatAPIPath("nodebalancers/%d", nodebalancerID)
	response, err := doPUTRequest[NodeBalancer](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteNodeBalancer deletes the NodeBalancer with the specified id
func (c *Client) DeleteNodeBalancer(ctx context.Context, nodebalancerID int) error {
	e := formatAPIPath("nodebalancers/%d", nodebalancerID)
	err := doDELETERequest(ctx, c, e)
	return err
}
