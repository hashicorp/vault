package linodego

import (
	"context"
	"encoding/json"

	"github.com/linode/linodego/pkg/errors"
)

// NetworkProtocol enum type
type NetworkProtocol string

// NetworkProtocol enum values
const (
	TCP  NetworkProtocol = "TCP"
	UDP  NetworkProtocol = "UDP"
	ICMP NetworkProtocol = "ICMP"
)

// NetworkAddresses are arrays of ipv4 and v6 addresses
type NetworkAddresses struct {
	IPv4 []string `json:"ipv4"`
	IPv6 []string `json:"ipv6"`
}

// A FirewallRule is a whitelist of ports, protocols, and addresses for which traffic should be allowed.
type FirewallRule struct {
	Ports     string           `json:"ports,omitempty"`
	Protocol  NetworkProtocol  `json:"protocol"`
	Addresses NetworkAddresses `json:"addresses"`
}

// FirewallRuleSet is a pair of inbound and outbound rules that specify what network traffic should be allowed.
type FirewallRuleSet struct {
	Inbound  []FirewallRule `json:"inbound,omitempty"`
	Outbound []FirewallRule `json:"outbound,omitempty"`
}

// GetFirewallRules gets the FirewallRuleSet for the given Firewall.
func (c *Client) GetFirewallRules(ctx context.Context, firewallID int) (*FirewallRuleSet, error) {
	e, err := c.FirewallRules.endpointWithID(firewallID)
	if err != nil {
		return nil, err
	}

	r, err := errors.CoupleAPIErrors(c.R(ctx).SetResult(&FirewallRuleSet{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*FirewallRuleSet), nil
}

// UpdateFirewallRules updates the FirewallRuleSet for the given Firewall
func (c *Client) UpdateFirewallRules(ctx context.Context, firewallID int, rules FirewallRuleSet) (*FirewallRuleSet, error) {
	e, err := c.FirewallRules.endpointWithID(firewallID)
	if err != nil {
		return nil, err
	}

	var body string
	req := c.R(ctx).SetResult(&FirewallRuleSet{})
	if bodyData, err := json.Marshal(rules); err == nil {
		body = string(bodyData)
	} else {
		return nil, errors.New(err)
	}

	r, err := errors.CoupleAPIErrors(req.SetBody(body).Put(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*FirewallRuleSet), nil
}
