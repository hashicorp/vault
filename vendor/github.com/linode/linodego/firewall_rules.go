package linodego

import (
	"context"
)

// NetworkProtocol enum type
type NetworkProtocol string

// NetworkProtocol enum values
const (
	TCP     NetworkProtocol = "TCP"
	UDP     NetworkProtocol = "UDP"
	ICMP    NetworkProtocol = "ICMP"
	IPENCAP NetworkProtocol = "IPENCAP"
)

// NetworkAddresses are arrays of ipv4 and v6 addresses
type NetworkAddresses struct {
	IPv4 *[]string `json:"ipv4,omitempty"`
	IPv6 *[]string `json:"ipv6,omitempty"`
}

// A FirewallRule is a whitelist of ports, protocols, and addresses for which traffic should be allowed.
type FirewallRule struct {
	Action      string           `json:"action"`
	Label       string           `json:"label"`
	Description string           `json:"description,omitempty"`
	Ports       string           `json:"ports,omitempty"`
	Protocol    NetworkProtocol  `json:"protocol"`
	Addresses   NetworkAddresses `json:"addresses"`
}

// FirewallRuleSet is a pair of inbound and outbound rules that specify what network traffic should be allowed.
type FirewallRuleSet struct {
	Inbound        []FirewallRule `json:"inbound"`
	InboundPolicy  string         `json:"inbound_policy"`
	Outbound       []FirewallRule `json:"outbound"`
	OutboundPolicy string         `json:"outbound_policy"`
}

// GetFirewallRules gets the FirewallRuleSet for the given Firewall.
func (c *Client) GetFirewallRules(ctx context.Context, firewallID int) (*FirewallRuleSet, error) {
	e := formatAPIPath("networking/firewalls/%d/rules", firewallID)
	response, err := doGETRequest[FirewallRuleSet](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateFirewallRules updates the FirewallRuleSet for the given Firewall
func (c *Client) UpdateFirewallRules(ctx context.Context, firewallID int, rules FirewallRuleSet) (*FirewallRuleSet, error) {
	e := formatAPIPath("networking/firewalls/%d/rules", firewallID)
	response, err := doPUTRequest[FirewallRuleSet](ctx, c, e, rules)
	if err != nil {
		return nil, err
	}

	return response, nil
}
