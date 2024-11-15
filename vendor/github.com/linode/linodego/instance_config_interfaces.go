package linodego

import (
	"context"
)

// InstanceConfigInterface contains information about a configuration's network interface
type InstanceConfigInterface struct {
	ID          int                    `json:"id"`
	IPAMAddress string                 `json:"ipam_address"`
	Label       string                 `json:"label"`
	Purpose     ConfigInterfacePurpose `json:"purpose"`
	Primary     bool                   `json:"primary"`
	Active      bool                   `json:"active"`
	VPCID       *int                   `json:"vpc_id"`
	SubnetID    *int                   `json:"subnet_id"`
	IPv4        *VPCIPv4               `json:"ipv4"`
	IPRanges    []string               `json:"ip_ranges"`
}

type VPCIPv4 struct {
	VPC     string  `json:"vpc,omitempty"`
	NAT1To1 *string `json:"nat_1_1,omitempty"`
}

type InstanceConfigInterfaceCreateOptions struct {
	IPAMAddress string                 `json:"ipam_address,omitempty"`
	Label       string                 `json:"label,omitempty"`
	Purpose     ConfigInterfacePurpose `json:"purpose,omitempty"`
	Primary     bool                   `json:"primary,omitempty"`
	SubnetID    *int                   `json:"subnet_id,omitempty"`
	IPv4        *VPCIPv4               `json:"ipv4,omitempty"`
	IPRanges    []string               `json:"ip_ranges,omitempty"`
}

type InstanceConfigInterfaceUpdateOptions struct {
	Primary  bool      `json:"primary,omitempty"`
	IPv4     *VPCIPv4  `json:"ipv4,omitempty"`
	IPRanges *[]string `json:"ip_ranges,omitempty"`
}

type InstanceConfigInterfacesReorderOptions struct {
	IDs []int `json:"ids"`
}

func getInstanceConfigInterfacesCreateOptionsList(
	interfaces []InstanceConfigInterface,
) []InstanceConfigInterfaceCreateOptions {
	interfaceOptsList := make([]InstanceConfigInterfaceCreateOptions, len(interfaces))
	for index, configInterface := range interfaces {
		interfaceOptsList[index] = configInterface.GetCreateOptions()
	}
	return interfaceOptsList
}

func (i InstanceConfigInterface) GetCreateOptions() InstanceConfigInterfaceCreateOptions {
	opts := InstanceConfigInterfaceCreateOptions{
		Label:    i.Label,
		Purpose:  i.Purpose,
		Primary:  i.Primary,
		SubnetID: i.SubnetID,
	}

	if len(i.IPRanges) > 0 {
		opts.IPRanges = i.IPRanges
	}

	if i.Purpose == InterfacePurposeVPC && i.IPv4 != nil {
		opts.IPv4 = &VPCIPv4{
			VPC:     i.IPv4.VPC,
			NAT1To1: i.IPv4.NAT1To1,
		}
	}

	opts.IPAMAddress = i.IPAMAddress

	return opts
}

func (i InstanceConfigInterface) GetUpdateOptions() InstanceConfigInterfaceUpdateOptions {
	opts := InstanceConfigInterfaceUpdateOptions{
		Primary: i.Primary,
	}

	if i.Purpose == InterfacePurposeVPC && i.IPv4 != nil {
		opts.IPv4 = &VPCIPv4{
			VPC:     i.IPv4.VPC,
			NAT1To1: i.IPv4.NAT1To1,
		}
	}

	if i.IPRanges != nil {
		// Copy the slice to prevent accidental
		// mutations
		copiedIPRanges := make([]string, len(i.IPRanges))
		copy(copiedIPRanges, i.IPRanges)

		opts.IPRanges = &copiedIPRanges
	}

	return opts
}

func (c *Client) AppendInstanceConfigInterface(
	ctx context.Context,
	linodeID int,
	configID int,
	opts InstanceConfigInterfaceCreateOptions,
) (*InstanceConfigInterface, error) {
	e := formatAPIPath("/linode/instances/%d/configs/%d/interfaces", linodeID, configID)
	response, err := doPOSTRequest[InstanceConfigInterface](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetInstanceConfigInterface(
	ctx context.Context,
	linodeID int,
	configID int,
	interfaceID int,
) (*InstanceConfigInterface, error) {
	e := formatAPIPath(
		"linode/instances/%d/configs/%d/interfaces/%d",
		linodeID,
		configID,
		interfaceID,
	)
	response, err := doGETRequest[InstanceConfigInterface](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ListInstanceConfigInterfaces(
	ctx context.Context,
	linodeID int,
	configID int,
) ([]InstanceConfigInterface, error) {
	e := formatAPIPath(
		"linode/instances/%d/configs/%d/interfaces",
		linodeID,
		configID,
	)
	response, err := doGETRequest[[]InstanceConfigInterface](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return *response, nil
}

func (c *Client) UpdateInstanceConfigInterface(
	ctx context.Context,
	linodeID int,
	configID int,
	interfaceID int,
	opts InstanceConfigInterfaceUpdateOptions,
) (*InstanceConfigInterface, error) {
	e := formatAPIPath(
		"linode/instances/%d/configs/%d/interfaces/%d",
		linodeID,
		configID,
		interfaceID,
	)
	response, err := doPUTRequest[InstanceConfigInterface](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) DeleteInstanceConfigInterface(
	ctx context.Context,
	linodeID int,
	configID int,
	interfaceID int,
) error {
	e := formatAPIPath(
		"linode/instances/%d/configs/%d/interfaces/%d",
		linodeID,
		configID,
		interfaceID,
	)
	err := doDELETERequest(ctx, c, e)
	return err
}

func (c *Client) ReorderInstanceConfigInterfaces(
	ctx context.Context,
	linodeID int,
	configID int,
	opts InstanceConfigInterfacesReorderOptions,
) error {
	e := formatAPIPath(
		"linode/instances/%d/configs/%d/interfaces/order",
		linodeID,
		configID,
	)
	_, err := doPOSTRequest[OAuthClient](ctx, c, e, opts)

	return err
}
