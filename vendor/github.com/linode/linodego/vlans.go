package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type VLAN struct {
	Label   string     `json:"label"`
	Linodes []int      `json:"linodes"`
	Region  string     `json:"region"`
	Created *time.Time `json:"-"`
}

// UnmarshalJSON for VLAN responses
func (v *VLAN) UnmarshalJSON(b []byte) error {
	type Mask VLAN

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
	}{
		Mask: (*Mask)(v),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	v.Created = (*time.Time)(p.Created)
	return nil
}

// ListVLANs returns a paginated list of VLANs
func (c *Client) ListVLANs(ctx context.Context, opts *ListOptions) ([]VLAN, error) {
	response, err := getPaginatedResults[VLAN](ctx, c, "networking/vlans", opts)
	return response, err
}

// GetVLANIPAMAddress returns the IPAM Address for a given VLAN Label as a string (10.0.0.1/24)
func (c *Client) GetVLANIPAMAddress(ctx context.Context, linodeID int, vlanLabel string) (string, error) {
	f := Filter{}
	f.AddField(Eq, "interfaces", vlanLabel)
	vlanFilter, err := f.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("Unable to convert VLAN label: %s to a filterable object: %w", vlanLabel, err)
	}

	cfgs, err := c.ListInstanceConfigs(ctx, linodeID, &ListOptions{Filter: string(vlanFilter)})
	if err != nil {
		return "", fmt.Errorf("Fetching configs for instance %v failed: %w", linodeID, err)
	}

	interfaces := cfgs[0].Interfaces
	for _, face := range interfaces {
		if face.Label == vlanLabel {
			return face.IPAMAddress, nil
		}
	}

	return "", fmt.Errorf("Failed to find IPAMAddress for VLAN: %s", vlanLabel)
}
