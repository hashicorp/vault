package packngo

import (
	"path"
)

// PortService handles operations on a port
type PortService interface {
	Assign(string, string) (*Port, *Response, error)
	Unassign(string, string) (*Port, *Response, error)
	AssignNative(string, string) (*Port, *Response, error)
	UnassignNative(string) (*Port, *Response, error)
	Bond(string, bool) (*Port, *Response, error)
	Disbond(string, bool) (*Port, *Response, error)
	ConvertToLayerTwo(string) (*Port, *Response, error)
	ConvertToLayerThree(string, []AddressRequest) (*Port, *Response, error)
	Get(string, *GetOptions) (*Port, *Response, error)
}

type PortServiceOp struct {
	client requestDoer
}

var _ PortService = (*PortServiceOp)(nil)

type PortData struct {
	// MAC address is set for NetworkPort ports
	MAC string `json:"mac,omitempty"`

	// Bonded is true for NetworkPort ports in a bond and NetworkBondPort ports
	// that are active
	Bonded bool `json:"bonded"`
}

type BondData struct {
	// ID is the Port.ID of the bonding port
	ID string `json:"id"`

	// Name of the port interface for the bond ("bond0")
	Name string `json:"name"`
}

// Port is a hardware port associated with a reserved or instantiated hardware
// device.
type Port struct {
	*Href `json:",inline"`

	// ID of the Port
	ID string `json:"id"`

	// Type is either "NetworkBondPort" for bond ports or "NetworkPort" for
	// bondable ethernet ports
	Type string `json:"type"`

	// Name of the interface for this port (such as "bond0" or "eth0")
	Name string `json:"name"`

	// Data about the port
	Data PortData `json:"data,omitempty"`

	// Indicates whether or not the bond can be broken on the port (when applicable).
	DisbondOperationSupported bool `json:"disbond_operation_supported,omitempty"`

	// NetworkType is either of layer2-bonded, layer2-individual, layer3,
	// hybrid, hybrid-bonded
	NetworkType string `json:"network_type,omitempty"`

	// NativeVirtualNetwork is the Native VLAN attached to the port
	// <https://metal.equinix.com/developers/docs/layer2-networking/native-vlan>
	NativeVirtualNetwork *VirtualNetwork `json:"native_virtual_network,omitempty"`

	// VLANs attached to the port
	AttachedVirtualNetworks []VirtualNetwork `json:"virtual_networks,omitempty"`

	// Bond details for ports with a NetworkPort type
	Bond *BondData `json:"bond,omitempty"`
}

type AddressRequest struct {
	AddressFamily int  `json:"address_family"`
	Public        bool `json:"public"`
}

type BackToL3Request struct {
	RequestIPs []AddressRequest `json:"request_ips"`
}

type PortAssignRequest struct {
	// PortID of the Port
	//
	// Deprecated: this is redundant to the portID parameter in request
	// functions. This is kept for use by deprecated DevicePortServiceOps
	// methods.
	PortID string `json:"id,omitempty"`

	VirtualNetworkID string `json:"vnid"`
}

type BondRequest struct {
	BulkEnable bool `json:"bulk_enable"`
}

type DisbondRequest struct {
	BulkDisable bool `json:"bulk_disable"`
}

// Assign adds a VLAN to a port
func (i *PortServiceOp) Assign(portID, vlanID string) (*Port, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(vlanID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(portBasePath, portID, "assign")
	par := &PortAssignRequest{VirtualNetworkID: vlanID}

	return i.portAction(apiPath, par)
}

// AssignNative assigns a virtual network to the port as a "native VLAN"
// The VLAN being assigned MUST first be added as a vlan using Assign() before 
// you may assign it as the native VLAN
func (i *PortServiceOp) AssignNative(portID, vlanID string) (*Port, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(vlanID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(portBasePath, portID, "native-vlan")
	par := &PortAssignRequest{VirtualNetworkID: vlanID}
	return i.portAction(apiPath, par)
}

// UnassignNative removes native VLAN from the supplied port
func (i *PortServiceOp) UnassignNative(portID string) (*Port, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(portBasePath, portID, "native-vlan")
	port := new(Port)

	resp, err := i.client.DoRequest("DELETE", apiPath, nil, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

// Unassign removes a VLAN from the port
func (i *PortServiceOp) Unassign(portID, vlanID string) (*Port, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(vlanID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(portBasePath, portID, "unassign")
	par := &PortAssignRequest{VirtualNetworkID: vlanID}

	return i.portAction(apiPath, par)
}

// Bond enables bonding for one or all ports
func (i *PortServiceOp) Bond(portID string, bulkEnable bool) (*Port, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	br := &BondRequest{BulkEnable: bulkEnable}
	apiPath := path.Join(portBasePath, portID, "bond")
	return i.portAction(apiPath, br)
}

// Disbond disables bonding for one or all ports
func (i *PortServiceOp) Disbond(portID string, bulkEnable bool) (*Port, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	dr := &DisbondRequest{BulkDisable: bulkEnable}
	apiPath := path.Join(portBasePath, portID, "disbond")
	return i.portAction(apiPath, dr)
}

func (i *PortServiceOp) portAction(apiPath string, req interface{}) (*Port, *Response, error) {
	port := new(Port)

	resp, err := i.client.DoRequest("POST", apiPath, req, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

// ConvertToLayerTwo converts a bond port to Layer 2. IP assignments of the port will be removed.
//
// portID is the UUID of a Bonding Port
func (i *PortServiceOp) ConvertToLayerTwo(portID string) (*Port, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(portBasePath, portID, "convert", "layer-2")
	port := new(Port)

	resp, err := i.client.DoRequest("POST", apiPath, nil, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

// ConvertToLayerThree converts a bond port to Layer 3. VLANs must first be unassigned.
func (i *PortServiceOp) ConvertToLayerThree(portID string, ips []AddressRequest) (*Port, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(portBasePath, portID, "convert", "layer-3")
	port := new(Port)

	req := BackToL3Request{
		RequestIPs: ips,
	}

	resp, err := i.client.DoRequest("POST", apiPath, &req, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

// Get returns a port by id
func (s *PortServiceOp) Get(portID string, opts *GetOptions) (*Port, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(portBasePath, portID)
	apiPathQuery := opts.WithQuery(endpointPath)
	port := new(Port)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, port)
	if err != nil {
		return nil, resp, err
	}
	return port, resp, err
}
