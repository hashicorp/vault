package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
	"github.com/linode/linodego/pkg/errors"
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

// FirewallDevicesPagedResponse represents a Linode API response for FirewallDevices
type FirewallDevicesPagedResponse struct {
	*PageOptions
	Data []FirewallDevice `json:"data"`
}

// endpointWithID gets the endpoint URL for FirewallDevices of a given Firewall
func (FirewallDevicesPagedResponse) endpointWithID(c *Client, id int) string {
	endpoint, err := c.FirewallDevices.endpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (resp *FirewallDevicesPagedResponse) appendData(r *FirewallDevicesPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListFirewallDevices get devices associated with a given Firewall
func (c *Client) ListFirewallDevices(ctx context.Context, firewallID int, opts *ListOptions) ([]FirewallDevice, error) {
	response := FirewallDevicesPagedResponse{}
	err := c.listHelperWithID(ctx, &response, firewallID, opts)

	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetFirewallDevice gets a FirewallDevice given an ID
func (c *Client) GetFirewallDevice(ctx context.Context, firewallID, deviceID int) (*FirewallDevice, error) {
	e, err := c.FirewallDevices.endpointWithID(firewallID)
	if err != nil {
		return nil, err
	}

	e = fmt.Sprintf("%s/%d", e, deviceID)
	r, err := errors.CoupleAPIErrors(c.R(ctx).SetResult(&FirewallDevice{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*FirewallDevice), nil
}

// AddFirewallDevice associates a Device with a given Firewall
func (c *Client) CreateFirewallDevice(ctx context.Context, firewallID int, createOpts FirewallDeviceCreateOptions) (*FirewallDevice, error) {
	var body string
	e, err := c.FirewallDevices.endpointWithID(firewallID)
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&FirewallDevice{})
	if bodyData, err := json.Marshal(createOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, errors.New(err)
	}

	r, err := errors.CoupleAPIErrors(req.SetBody(body).Post(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*FirewallDevice), nil
}

// DeleteFirewallDevice disassociates a Device with a given Firewall
func (c *Client) DeleteFirewallDevice(ctx context.Context, firewallID, deviceID int) error {
	e, err := c.FirewallDevices.endpointWithID(firewallID)
	if err != nil {
		return err
	}

	e = fmt.Sprintf("%s/%d", e, deviceID)
	_, err = errors.CoupleAPIErrors(c.R(ctx).Delete(e))
	return err
}
