package packngo

import "path"

type VCStatus string

const (
	virtualCircuitBasePath = "/virtual-circuits"

	// VC is being create but not ready yet
	VCStatusPending VCStatus = "pending"

	// VC is ready with a VLAN
	VCStatusActive VCStatus = "active"

	// VC is ready without a VLAN
	VCStatusWaiting VCStatus = "waiting_on_customer_vlan"

	// VC is being deleted
	VCStatusDeleting VCStatus = "deleting"

	// not sure what the following states mean, or whether they exist
	// someone from the API side could check
	VCStatusActivating         VCStatus = "activating"
	VCStatusDeactivating       VCStatus = "deactivating"
	VCStatusActivationFailed   VCStatus = "activation_failed"
	VCStatusDeactivationFailed VCStatus = "dactivation_failed"
)

type VirtualCircuitService interface {
	Create(string, string, string, *VCCreateRequest, *GetOptions) (*VirtualCircuit, *Response, error)
	Get(string, *GetOptions) (*VirtualCircuit, *Response, error)
	Events(string, *GetOptions) ([]Event, *Response, error)
	Delete(string) (*Response, error)
	Update(string, *VCUpdateRequest, *GetOptions) (*VirtualCircuit, *Response, error)
}

type VCUpdateRequest struct {
	Name             *string   `json:"name,omitempty"`
	Tags             *[]string `json:"tags,omitempty"`
	Description      *string   `json:"description,omitempty"`
	VirtualNetworkID *string   `json:"vnid,omitempty"`

	// Speed is a bps representation of the VirtualCircuit throughput. This is informational only, the field is a user-controlled description of the speed. It may be presented as a whole number with a bps, mpbs, or gbps suffix (or the respective initial).
	Speed string `json:"speed,omitempty"`
}

type VCCreateRequest struct {
	// VirtualNetworkID of the Virtual Network to connect to the Virtual Circuit (required when VRFID is not specified)
	VirtualNetworkID string `json:"vnid,omitempty"`
	// VRFID of the VRF to connect to the Virtual Circuit (required when VirtualNetworkID is not specified)
	VRFID string `json:"vrf_id,omitempty"`
	// PeerASN (optional, required with VRFID) The BGP ASN of the device the switch will peer with. Can be the used across several VCs, but cannot be the same as the local_asn.
	PeerASN int `json:"peer_asn,omitempty"`
	// Subnet (Required for VRF) A subnet from one of the IP blocks associated with the VRF that we
	// will help create an IP reservation for. Can only be either a /30 or /31.
	//  * For a /31 block, it will only have two IP addresses, which will be used for the metal_ip and customer_ip.
	//  * For a /30 block, it will have four IP
	// addresses, but the first and last IP addresses are not usable. We will
	// default to the first usable IP address for the metal_ip.
	Subnet string `json:"subnet,omitempty"`
	// MetalIP (optional, required with VRFID) The IP address that’s set as “our” IP that is
	// configured on the rack_local_vlan SVI. Will default to the first usable
	// IP in the subnet.
	MetalIP string `json:"metal_ip,omitempty"`
	// CustomerIP (optional, requires VRFID) The IP address set as the customer IP which the CSR
	// switch will peer with. Will default to the other usable IP in the subnet.
	CustomerIP string `json:"customer_ip,omitempty"`
	// MD5 (optional, requires VRFID) The password that can be set for the VRF BGP peer
	MD5         string   `json:"md5,omitempty"`
	NniVLAN     int      `json:"nni_vlan,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`

	// Speed is a bps representation of the VirtualCircuit throughput. This is informational only, the field is a user-controlled description of the speed. It may be presented as a whole number with a bps, mpbs, or gbps suffix (or the respective initial).
	Speed string `json:"speed,omitempty"`
}

type VirtualCircuitServiceOp struct {
	client *Client
}

type virtualCircuitsRoot struct {
	VirtualCircuits []VirtualCircuit `json:"virtual_circuits"`
	Meta            meta             `json:"meta"`
}

type VirtualCircuit struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	// Speed of the VirtualCircuit in bps
	Speed          int             `json:"speed,omitempty"`
	Status         VCStatus        `json:"status,omitempty"`
	VNID           int             `json:"vnid,omitempty"`
	NniVNID        int             `json:"nni_vnid,omitempty"`
	NniVLAN        int             `json:"nni_vlan,omitempty"`
	Project        *Project        `json:"project,omitempty"`
	Port           *ConnectionPort `json:"port,omitempty"`
	VirtualNetwork *VirtualNetwork `json:"virtual_network,omitempty"`
	Tags           []string        `json:"tags,omitempty"`
	// VRF connected to the Virtual Circuit
	VRF *VRF `json:"vrf,omitempty"`

	// PeerASN (optional, required with VRFID) The BGP ASN of the device the switch will peer with. Can be the used across several VCs, but cannot be the same as the local_asn.
	PeerASN int `json:"peer_asn,omitempty"`

	// Subnet (returned with VRF) A subnet from one of the IP blocks associated with the VRF that we
	// will help create an IP reservation for. Can only be either a /30 or /31.
	//  * For a /31 block, it will only have two IP addresses, which will be used for the metal_ip and customer_ip.
	//  * For a /30 block, it will have four IP
	// addresses, but the first and last IP addresses are not usable. We will
	// default to the first usable IP address for the metal_ip.
	Subnet string `json:"subnet,omitempty"`

	// MetalIP (returned with VRF) The IP address that’s set as “our” IP that is
	// configured on the rack_local_vlan SVI. Will default to the first usable
	// IP in the subnet.
	MetalIP string `json:"metal_ip,omitempty"`

	// CustomerIP (returned with VRF) The IP address set as the customer IP which the CSR
	// switch will peer with. Will default to the other usable IP in the subnet.
	CustomerIP string `json:"customer_ip,omitempty"`

	// MD5 (returned with VRF) The password that can be set for the VRF BGP peer
	MD5 string `json:"md5,omitempty"`
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
	if validateErr := ValidateUUID(vcID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(virtualCircuitBasePath, vcID)
	apiPathQuery := opts.WithQuery(endpointPath)
	return s.do("PUT", apiPathQuery, req)
}

func (s *VirtualCircuitServiceOp) Events(id string, opts *GetOptions) ([]Event, *Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(virtualCircuitBasePath, id, eventBasePath)
	return listEvents(s.client, apiPath, opts)
}

func (s *VirtualCircuitServiceOp) Get(id string, opts *GetOptions) (*VirtualCircuit, *Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(virtualCircuitBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	return s.do("GET", apiPathQuery, nil)
}

func (s *VirtualCircuitServiceOp) Delete(id string) (*Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(virtualCircuitBasePath, id)
	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}

func (s *VirtualCircuitServiceOp) Create(projectID, connID, portID string, request *VCCreateRequest, opts *GetOptions) (*VirtualCircuit, *Response, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(connID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(projectBasePath, projectID, connectionBasePath, connID, portBasePath, portID, virtualCircuitBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	return s.do("POST", apiPathQuery, request)
}
