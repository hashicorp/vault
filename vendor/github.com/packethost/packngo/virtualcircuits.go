package packngo

import "path"

const (
	virtualCircuitBasePath = "/virtual-circuits"
	vcStatusActive         = "active"
	vcStatusWaiting        = "waiting_on_customer_vlan"
	//vcStatusActivating         = "activating"
	//vcStatusDeactivating       = "deactivating"
	vcStatusActivationFailed   = "activation_failed"
	vcStatusDeactivationFailed = "dactivation_failed"
)

type VirtualCircuitService interface {
	Create(string, string, string, *VCCreateRequest, *GetOptions) (*VirtualCircuit, *Response, error)
	Get(string, *GetOptions) (*VirtualCircuit, *Response, error)
	Events(string, *GetOptions) ([]Event, *Response, error)
	Delete(string) (*Response, error)
	Update(string, *VCUpdateRequest, *GetOptions) (*VirtualCircuit, *Response, error)
}

type VCUpdateRequest struct {
	VirtualNetworkID *string `json:"vnid"`
}

type VCCreateRequest struct {
	VirtualNetworkID string `json:"vnid"`
	NniVLAN          int    `json:"nni_vlan,omitempty"`
	Name             string `json:"name,omitempty"`
}

type VirtualCircuitServiceOp struct {
	client *Client
}

type virtualCircuitsRoot struct {
	VirtualCircuits []VirtualCircuit `json:"virtual_circuits"`
	Meta            meta             `json:"meta"`
}

type VirtualCircuit struct {
	ID             string          `json:"id"`
	Name           string          `json:"name,omitempty"`
	Status         string          `json:"status,omitempty"`
	VNID           int             `json:"vnid,omitempty"`
	NniVNID        int             `json:"nni_vnid,omitempty"`
	NniVLAN        int             `json:"nni_vlan,omitempty"`
	Project        *Project        `json:"project,omitempty"`
	VirtualNetwork *VirtualNetwork `json:"virtual_network,omitempty"`
}

func (s *VirtualCircuitServiceOp) do(method, apiPathQuery string, req interface{}) (*VirtualCircuit, *Response, error) {
	vc := new(VirtualCircuit)
	resp, err := s.client.DoRequest(method, apiPathQuery, req, vc)
	if err != nil {
		return nil, resp, err
	}
	return vc, resp, err
}

func (s *VirtualCircuitServiceOp) Update(vcID string, req *VCUpdateRequest, opts *GetOptions) (*VirtualCircuit, *Response, error) {
	endpointPath := path.Join(virtualCircuitBasePath, vcID)
	apiPathQuery := opts.WithQuery(endpointPath)
	return s.do("PUT", apiPathQuery, req)
}

func (s *VirtualCircuitServiceOp) Events(id string, opts *GetOptions) ([]Event, *Response, error) {
	apiPath := path.Join(virtualCircuitBasePath, id, eventBasePath)
	return listEvents(s.client, apiPath, opts)
}

func (s *VirtualCircuitServiceOp) Get(id string, opts *GetOptions) (*VirtualCircuit, *Response, error) {
	endpointPath := path.Join(virtualCircuitBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	return s.do("GET", apiPathQuery, nil)
}

func (s *VirtualCircuitServiceOp) Delete(id string) (*Response, error) {
	apiPath := path.Join(virtualCircuitBasePath, id)
	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}

func (s *VirtualCircuitServiceOp) Create(projectID, connID, portID string, request *VCCreateRequest, opts *GetOptions) (*VirtualCircuit, *Response, error) {
	endpointPath := path.Join(projectBasePath, projectID, connectionBasePath, connID, portBasePath, portID, virtualCircuitBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	return s.do("POST", apiPathQuery, request)
}
