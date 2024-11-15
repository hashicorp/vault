package packngo

import (
	"path"
)

const virtualNetworkBasePath = "/virtual-networks"

// DevicePortService handles operations on a port which belongs to a particular device
type ProjectVirtualNetworkService interface {
	List(projectID string, opts *ListOptions) (*VirtualNetworkListResponse, *Response, error)
	Create(*VirtualNetworkCreateRequest) (*VirtualNetwork, *Response, error)
	Get(string, *GetOptions) (*VirtualNetwork, *Response, error)
	Delete(virtualNetworkID string) (*Response, error)
}

type VirtualNetwork struct {
	ID           string    `json:"id"`
	Description  string    `json:"description,omitempty"` // TODO: field can be null
	VXLAN        int       `json:"vxlan,omitempty"`
	FacilityCode string    `json:"facility_code,omitempty"`
	MetroCode    string    `json:"metro_code,omitempty"`
	CreatedAt    string    `json:"created_at,omitempty"`
	Href         string    `json:"href"`
	Project      *Project  `json:"assigned_to,omitempty"`
	Facility     *Facility `json:"facility,omitempty"`
	Metro        *Metro    `json:"metro,omitempty"`
	Instances    []*Device `json:"instances,omitempty"`
}

type ProjectVirtualNetworkServiceOp struct {
	client *Client
}

type VirtualNetworkListResponse struct {
	VirtualNetworks []VirtualNetwork `json:"virtual_networks"`
}

func (i *ProjectVirtualNetworkServiceOp) List(projectID string, opts *ListOptions) (*VirtualNetworkListResponse, *Response, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(projectBasePath, projectID, virtualNetworkBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	output := new(VirtualNetworkListResponse)

	resp, err := i.client.DoRequest("GET", apiPathQuery, nil, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

type VirtualNetworkCreateRequest struct {
	// ProjectID of the project where the VLAN will be made available.
	ProjectID string `json:"project_id"`

	// Description is a user supplied description of the VLAN.
	Description string `json:"description"`

	// TODO: default Description is null when not specified. Permitting *string here would require changing VirtualNetwork.Description to *string too.

	// Facility in which to create the VLAN. Mutually exclusive with Metro.
	Facility string `json:"facility,omitempty"`

	// Metro in which to create the VLAN. Mutually exclusive with Facility.
	Metro string `json:"metro,omitempty"`

	// VXLAN is the VLAN ID. VXLAN may be specified when Metro is defined. It is remotely incremented otherwise. Must be unique per Metro.
	VXLAN int `json:"vxlan,omitempty"`
}

func (i *ProjectVirtualNetworkServiceOp) Get(vlanID string, opts *GetOptions) (*VirtualNetwork, *Response, error) {
	if validateErr := ValidateUUID(vlanID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(virtualNetworkBasePath, vlanID)
	apiPathQuery := opts.WithQuery(endpointPath)
	vlan := new(VirtualNetwork)

	resp, err := i.client.DoRequest("GET", apiPathQuery, nil, vlan)
	if err != nil {
		return nil, resp, err
	}

	return vlan, resp, err
}

func (i *ProjectVirtualNetworkServiceOp) Create(input *VirtualNetworkCreateRequest) (*VirtualNetwork, *Response, error) {
	// TODO: May need to add timestamp to output from 'post' request
	// for the 'created_at' attribute of VirtualNetwork struct since
	// API response doesn't include it
	apiPath := path.Join(projectBasePath, input.ProjectID, virtualNetworkBasePath)
	output := new(VirtualNetwork)

	resp, err := i.client.DoRequest("POST", apiPath, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

func (i *ProjectVirtualNetworkServiceOp) Delete(virtualNetworkID string) (*Response, error) {
	if validateErr := ValidateUUID(virtualNetworkID); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(virtualNetworkBasePath, virtualNetworkID)

	resp, err := i.client.DoRequest("DELETE", apiPath, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
