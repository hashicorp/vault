package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// FirewallDeviceType represents the different kinds of devices governable by a Firewall
type FirewallDeviceType string

// FirewallDeviceType constants start with FirewallDevice
const (
	FirewallDeviceLinode       FirewallDeviceType = "linode"
	FirewallDeviceNodeBalancer FirewallDeviceType = "nodebalancer"
)

// FirewallDevice represents a device governed by a Firewall
type FirewallDevice struct {
	ID      int                  `json:"id"`
	Entity  FirewallDeviceEntity `json:"entity"`
	Created *time.Time           `json:"-"`
	Updated *time.Time           `json:"-"`
}

// FirewallDeviceCreateOptions fields are those accepted by CreateFirewallDevice
type FirewallDeviceCreateOptions struct {
	ID   int                `json:"id"`
	Type FirewallDeviceType `json:"type"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (device *FirewallDevice) UnmarshalJSON(b []byte) error {
	type Mask FirewallDevice

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(device),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	device.Created = (*time.Time)(p.Created)
	device.Updated = (*time.Time)(p.Updated)
	return nil
}

// FirewallDeviceEntity contains information about a device associated with a Firewall
type FirewallDeviceEntity struct {
	ID    int                `json:"id"`
	Type  FirewallDeviceType `json:"type"`
	Label string             `json:"label"`
	URL   string             `json:"url"`
}

// ListFirewallDevices get devices associated with a given Firewall
func (c *Client) ListFirewallDevices(ctx context.Context, firewallID int, opts *ListOptions) ([]FirewallDevice, error) {
	response, err := getPaginatedResults[FirewallDevice](ctx, c, formatAPIPath("networking/firewalls/%d/devices", firewallID), opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetFirewallDevice gets a FirewallDevice given an ID
func (c *Client) GetFirewallDevice(ctx context.Context, firewallID, deviceID int) (*FirewallDevice, error) {
	e := formatAPIPath("networking/firewalls/%d/devices/%d", firewallID, deviceID)
	response, err := doGETRequest[FirewallDevice](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// AddFirewallDevice associates a Device with a given Firewall
func (c *Client) CreateFirewallDevice(ctx context.Context, firewallID int, opts FirewallDeviceCreateOptions) (*FirewallDevice, error) {
	e := formatAPIPath("networking/firewalls/%d/devices", firewallID)
	response, err := doPOSTRequest[FirewallDevice](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteFirewallDevice disassociates a Device with a given Firewall
func (c *Client) DeleteFirewallDevice(ctx context.Context, firewallID, deviceID int) error {
	e := formatAPIPath("networking/firewalls/%d/devices/%d", firewallID, deviceID)
	err := doDELETERequest(ctx, c, e)
	return err
}
