package packngo

import (
	"path"
)

type MetalGatewayState string

const (
	metalGatewayBasePath                   = "/metal-gateways"
	MetalGatewayActive   MetalGatewayState = "active"
	MetalGatewayReady    MetalGatewayState = "ready"
	MetalGatewayDeleting MetalGatewayState = "deleting"
)

type MetalGatewayService interface {
	List(projectID string, opts *ListOptions) ([]MetalGateway, *Response, error)
	Create(projectID string, input *MetalGatewayCreateRequest) (*MetalGateway, *Response, error)
	Get(metalGatewayID string, opts *GetOptions) (*MetalGateway, *Response, error)
	Delete(metalGatewayID string) (*Response, error)
}

type MetalGateway struct {
	ID             string                `json:"id"`
	State          MetalGatewayState     `json:"state"`
	Project        *Project              `json:"project,omitempty"`
	VirtualNetwork *VirtualNetwork       `json:"virtual_network,omitempty"`
	IPReservation  *IPAddressReservation `json:"ip_reservation,omitempty"`
	Href           string                `json:"href"`
	CreatedAt      string                `json:"created_at,omitempty"`
	UpdatedAt      string                `json:"updated_at,omitempty"`
	VRF            *VRF                  `json:"vrf,omitempty"`
}

// MetalGatewayLite struct representation of a Metal Gateway
type MetalGatewayLite struct {
	ID string `json:"id,omitempty"`
	// The current state of the Metal Gateway. 'Ready' indicates the gateway record has been configured, but is currently not active on the network. 'Active' indicates the gateway has been configured on the network. 'Deleting' is a temporary state used to indicate that the gateway is in the process of being un-configured from the network, after which the gateway record will be deleted.
	State     string     `json:"state,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
	UpdatedAt *Timestamp `json:"updated_at,omitempty"`
	// The gateway address with subnet CIDR value for this Metal Gateway. For example, a Metal Gateway using an IP reservation with block 10.1.2.0/27 would have a gateway address of 10.1.2.1/27.
	GatewayAddress string `json:"gateway_address,omitempty"`
	// The VLAN id of the Virtual Network record associated to this Metal Gateway. Example: 1001.
	VLAN int    `json:"vlan,omitempty"`
	Href string `json:"href,omitempty"`
}

type MetalGatewayServiceOp struct {
	client *Client
}

func (s *MetalGatewayServiceOp) List(projectID string, opts *ListOptions) (metalGateways []MetalGateway, resp *Response, err error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	type metalGatewaysRoot struct {
		MetalGateways []MetalGateway `json:"metal_gateways"`
		Meta          meta           `json:"meta"`
	}

	endpointPath := path.Join(projectBasePath, projectID, metalGatewayBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(metalGatewaysRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		metalGateways = append(metalGateways, subset.MetalGateways...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

type MetalGatewayCreateRequest struct {
	// VirtualNetworkID virtual network UUID.
	VirtualNetworkID string `json:"virtual_network_id"`

	// IPReservationID (optional) IP Reservation UUID (Public or VRF). Required for VRF.
	IPReservationID string `json:"ip_reservation_id,omitempty"`

	// PrivateIPv4SubnetSize (optional) Power of 2 between 8 and 128 (8, 16, 32, 64, 128). Invalid for VRF.
	PrivateIPv4SubnetSize int `json:"private_ipv4_subnet_size,omitempty"`
}

func (s *MetalGatewayServiceOp) Get(metalGatewayID string, opts *GetOptions) (*MetalGateway, *Response, error) {
	if validateErr := ValidateUUID(metalGatewayID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(metalGatewayBasePath, metalGatewayID)
	apiPathQuery := opts.WithQuery(endpointPath)
	metalGateway := new(MetalGateway)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, metalGateway)
	if err != nil {
		return nil, resp, err
	}

	return metalGateway, resp, err
}

func (s *MetalGatewayServiceOp) Create(projectID string, input *MetalGatewayCreateRequest) (*MetalGateway, *Response, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(projectBasePath, projectID, metalGatewayBasePath)
	output := new(MetalGateway)

	resp, err := s.client.DoRequest("POST", apiPath, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

func (s *MetalGatewayServiceOp) Delete(metalGatewayID string) (*Response, error) {
	if validateErr := ValidateUUID(metalGatewayID); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(metalGatewayBasePath, metalGatewayID)

	resp, err := s.client.DoRequest("DELETE", apiPath, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
