package packngo

import (
	"path"
	"strconv"
)

const ipBasePath = "/ips"

const (
	// PublicIPv4 fixed string representation of public ipv4
	PublicIPv4 IPReservationType = "public_ipv4"
	// PrivateIPv4 fixed string representation of private ipv4
	PrivateIPv4 IPReservationType = "private_ipv4"
	// GlobalIPv4 fixed string representation of global ipv4
	GlobalIPv4 IPReservationType = "global_ipv4"
	// PublicIPv6 fixed string representation of public ipv6
	PublicIPv6 IPReservationType = "public_ipv6"
	// PrivateIPv6 fixed string representation of private ipv6
	PrivateIPv6 IPReservationType = "private_ipv6"
	// GlobalIPv6 fixed string representation of global ipv6
	GlobalIPv6 IPReservationType = "global_ipv6"
	// VRFIPRange fixed string representation of vrf (virtual routing and forwarding). This may be any VRF supported range, including public and RFC-1918 IPv4 and IPv6 ranges.
	VRFIPRange IPReservationType = "vrf"
)

type IPReservationType string

// DeviceIPService handles assignment of addresses from reserved blocks to instances in a project.
type DeviceIPService interface {
	Assign(deviceID string, assignRequest *AddressStruct) (*IPAddressAssignment, *Response, error)
	Unassign(assignmentID string) (*Response, error)
	Get(assignmentID string, getOpt *GetOptions) (*IPAddressAssignment, *Response, error)
	List(deviceID string, opts *ListOptions) ([]IPAddressAssignment, *Response, error)
}

// ProjectIPService handles reservation of IP address blocks for a project.
type ProjectIPService interface {
	Get(reservationID string, getOpt *GetOptions) (*IPAddressReservation, *Response, error)
	List(projectID string, opts *ListOptions) ([]IPAddressReservation, *Response, error)
	// Deprecated Use Create instead of Request
	Request(projectID string, ipReservationReq *IPReservationRequest) (*IPAddressReservation, *Response, error)
	Create(projectID string, ipReservationReq *IPReservationCreateRequest) (*IPAddressReservation, *Response, error)
	Delete(ipReservationID string) (*Response, error)
	// Deprecated Use Delete instead of Remove
	Remove(ipReservationID string) (*Response, error)
	Update(assignmentID string, updateRequest *IPAddressUpdateRequest, opt *GetOptions) (*IPAddressReservation, *Response, error)
	AvailableAddresses(ipReservationID string, r *AvailableRequest) ([]string, *Response, error)
}

var (
	_ DeviceIPService  = (*DeviceIPServiceOp)(nil)
	_ ProjectIPService = (*ProjectIPServiceOp)(nil)
)

type IpAddressCommon struct { //nolint:golint
	ID            string            `json:"id"`
	Address       string            `json:"address"`
	Gateway       string            `json:"gateway"`
	Network       string            `json:"network"`
	AddressFamily int               `json:"address_family"`
	Netmask       string            `json:"netmask"`
	Public        bool              `json:"public"`
	CIDR          int               `json:"cidr"`
	Created       string            `json:"created_at,omitempty"`
	Updated       string            `json:"updated_at,omitempty"`
	Href          string            `json:"href"`
	Management    bool              `json:"management"`
	Manageable    bool              `json:"manageable"`
	Metro         *Metro            `json:"metro,omitempty"`
	Project       Href              `json:"project"`
	Global        bool              `json:"global_ip"`
	Tags          []string          `json:"tags,omitempty"`
	ParentBlock   *ParentBlock      `json:"parent_block,omitempty"`
	CustomData    interface{}       `json:"customdata,omitempty"`
	Type          IPReservationType `json:"type"`
	VRF           *VRF              `json:"vrf,omitempty"`
}

// ParentBlock is the network block for the parent of an IP address
type ParentBlock struct {
	Network string  `json:"network"`
	Netmask string  `json:"netmask"`
	CIDR    int     `json:"cidr"`
	Href    *string `json:"href,omitempty"`
}

type IPReservationState string

const (
	// IPReservationStatePending fixed string representation of pending
	IPReservationStatePending IPReservationState = "pending"

	// IPReservationStateCreated fixed string representation of created
	IPReservationStateCreated IPReservationState = "created"

	// IPReservationStateDenied fixed string representation of denied
	IPReservationStateDenied IPReservationState = "denied"
)

// IPAddressReservation is created when user sends IP reservation request for a project (considering it's within quota).
type IPAddressReservation struct {
	IpAddressCommon
	Assignments  []*IPAddressAssignment `json:"assignments"`
	Facility     *Facility              `json:"facility,omitempty"`
	Available    string                 `json:"available"`
	Addon        bool                   `json:"addon"`
	Bill         bool                   `json:"bill"`
	State        IPReservationState     `json:"state"`
	Description  *string                `json:"details"`
	Enabled      bool                   `json:"enabled"`
	MetalGateway *MetalGatewayLite      `json:"metal_gateway,omitempty"`
	RequestedBy  *UserLite              `json:"requested_by,omitempty"`
}

// AvailableResponse is a type for listing of available addresses from a reserved block.
type AvailableResponse struct {
	Available []string `json:"available"`
}

// AvailableRequest is a type for listing available addresses from a reserved block.
type AvailableRequest struct {
	CIDR int `json:"cidr"`
}

// IPAddressAssignment is created when an IP address from reservation block is assigned to a device.
type IPAddressAssignment struct {
	IpAddressCommon
	AssignedTo Href `json:"assigned_to"`
}

// IPAddressUpdateRequest represents the body of an IPAddress patch
type IPAddressUpdateRequest struct {
	Tags        *[]string   `json:"tags,omitempty"`
	Description *string     `json:"details,omitempty"`
	CustomData  interface{} `json:"customdata,omitempty"`
}

// IPReservationRequest represents the body of an IP reservation request
// Deprecated: use IPReservationCreateRequest
type IPReservationRequest = IPReservationCreateRequest

// IPReservationCreateRequest represents the body of an IP reservation request.
type IPReservationCreateRequest struct {
	// Type of IP reservation.
	Type        IPReservationType `json:"type"`
	Quantity    int               `json:"quantity"`
	Description string            `json:"details,omitempty"`
	Facility    *string           `json:"facility,omitempty"`
	Metro       *string           `json:"metro,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	CustomData  interface{}       `json:"customdata,omitempty"`
	// FailOnApprovalRequired if the IP request cannot be approved
	// automatically, rather than sending to the longer Equinix Metal approval
	// process, fail immediately with a 422 error
	FailOnApprovalRequired bool `json:"fail_on_approval_required,omitempty"`

	// Comments in support of the request for additional addresses when the
	// request must be manually approved.
	Comments string `json:"comments,omitempty"`

	// VRFID is the ID of the VRF to associate and draw the IP range from.
	// * Required when Type is VRFIPRange, not valid otherwise
	// * Network and CIDR are required when set
	// * Metro and Facility are not required when set
	VRFID string `json:"vrf_id,omitempty"`

	// Network an unreserved network address from an existing VRF ip_range.
	// * Required when Type is VRFIPRange, not valid otherwise
	Network string `json:"network,omitempty"`

	// CIDR the size of the network to reserve from an existing VRF ip_range.
	// * Required when Type is VRFIPRange, not valid otherwise
	// * Minimum range is 22-29, with 30-31 supported and necessary for virtual-circuits
	CIDR int `json:"cidr,omitempty"`
}

// AddressStruct is a helper type for request/response with dict like {"address": ... }
type AddressStruct struct {
	Address string `json:"address"`
}

func deleteFromIP(client *Client, resourceID string) (*Response, error) {
	if validateErr := ValidateUUID(resourceID); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(ipBasePath, resourceID)

	return client.DoRequest("DELETE", apiPath, nil, nil)
}

func (i IPAddressReservation) String() string {
	return Stringify(i)
}

func (i IPAddressAssignment) String() string {
	return Stringify(i)
}

// DeviceIPServiceOp is interface for IP-address assignment methods.
type DeviceIPServiceOp struct {
	client *Client
}

// Unassign unassigns an IP address from the device to which it is currently assigned.
// This will remove the relationship between an IP and the device and will make the IP
// address available to be assigned to another device.
func (i *DeviceIPServiceOp) Unassign(assignmentID string) (*Response, error) {
	if validateErr := ValidateUUID(assignmentID); validateErr != nil {
		return nil, validateErr
	}
	return deleteFromIP(i.client, assignmentID)
}

// Assign assigns an IP address to a device.
// The IP address must be in one of the IP ranges assigned to the device’s project.
func (i *DeviceIPServiceOp) Assign(deviceID string, assignRequest *AddressStruct) (*IPAddressAssignment, *Response, error) {
	if validateErr := ValidateUUID(deviceID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(deviceBasePath, deviceID, ipBasePath)
	ipa := new(IPAddressAssignment)

	resp, err := i.client.DoRequest("POST", apiPath, assignRequest, ipa)
	if err != nil {
		return nil, resp, err
	}

	return ipa, resp, err
}

// Get returns assignment by ID.
func (i *DeviceIPServiceOp) Get(assignmentID string, opts *GetOptions) (*IPAddressAssignment, *Response, error) {
	if validateErr := ValidateUUID(assignmentID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(ipBasePath, assignmentID)
	apiPathQuery := opts.WithQuery(endpointPath)
	ipa := new(IPAddressAssignment)

	resp, err := i.client.DoRequest("GET", apiPathQuery, nil, ipa)
	if err != nil {
		return nil, resp, err
	}

	return ipa, resp, err
}

// List list all of the IP address assignments on a device
func (i *DeviceIPServiceOp) List(deviceID string, opts *ListOptions) ([]IPAddressAssignment, *Response, error) {
	if validateErr := ValidateUUID(deviceID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(deviceBasePath, deviceID, ipBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	// ipList represents collection of IP Address reservations
	type ipList struct {
		IPs []IPAddressAssignment `json:"ip_addresses,omitempty"`
	}

	ips := new(ipList)

	resp, err := i.client.DoRequest("GET", apiPathQuery, nil, ips)
	if err != nil {
		return nil, resp, err
	}

	return ips.IPs, resp, err
}

// ProjectIPServiceOp is interface for IP assignment methods.
type ProjectIPServiceOp struct {
	client *Client
}

// Get returns reservation by ID.
func (i *ProjectIPServiceOp) Get(reservationID string, opts *GetOptions) (*IPAddressReservation, *Response, error) {
	if validateErr := ValidateUUID(reservationID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(ipBasePath, reservationID)
	apiPathQuery := opts.WithQuery(endpointPath)
	ipr := new(IPAddressReservation)

	resp, err := i.client.DoRequest("GET", apiPathQuery, nil, ipr)
	if err != nil {
		return nil, resp, err
	}

	return ipr, resp, err
}

// List provides a list of IP resevations for a single project.
// opts be filtered to limit the type of reservations returned:
// opts.Filter("type", "vrf")
func (i *ProjectIPServiceOp) List(projectID string, opts *ListOptions) ([]IPAddressReservation, *Response, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(projectBasePath, projectID, ipBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	reservations := new(struct {
		Reservations []IPAddressReservation `json:"ip_addresses"`
	})

	resp, err := i.client.DoRequest("GET", apiPathQuery, nil, reservations)
	if err != nil {
		return nil, resp, err
	}
	return reservations.Reservations, resp, nil
}

// Create creates a request for more IP space for a project in order to have
// additional IP addresses to assign to devices.
func (i *ProjectIPServiceOp) Create(projectID string, ipReservationReq *IPReservationCreateRequest) (*IPAddressReservation, *Response, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(projectBasePath, projectID, ipBasePath)
	ipr := new(IPAddressReservation)

	resp, err := i.client.DoRequest("POST", apiPath, ipReservationReq, ipr)
	if err != nil {
		return nil, resp, err
	}
	return ipr, resp, err
}

// Request requests more IP space for a project in order to have additional IP
// addresses to assign to devices.
//
// Deprecated: Use Create instead.
func (i *ProjectIPServiceOp) Request(projectID string, ipReservationReq *IPReservationRequest) (*IPAddressReservation, *Response, error) {
	return i.Create(projectID, ipReservationReq)
}

// Update updates an existing IP reservation.
func (i *ProjectIPServiceOp) Update(reservationID string, updateRequest *IPAddressUpdateRequest, opts *GetOptions) (*IPAddressReservation, *Response, error) {
	if validateErr := ValidateUUID(reservationID); validateErr != nil {
		return nil, nil, validateErr
	}
	if opts == nil {
		opts = &GetOptions{}
	}
	endpointPath := path.Join(ipBasePath, reservationID)
	apiPathQuery := opts.WithQuery(endpointPath)
	ipr := new(IPAddressReservation)

	resp, err := i.client.DoRequest("PATCH", apiPathQuery, updateRequest, ipr)
	if err != nil {
		return nil, resp, err
	}

	return ipr, resp, err
}

// Delete removes the requests for specific IP within a project
func (i *ProjectIPServiceOp) Delete(ipReservationID string) (*Response, error) {
	if validateErr := ValidateUUID(ipReservationID); validateErr != nil {
		return nil, validateErr
	}
	return deleteFromIP(i.client, ipReservationID)
}

// Remove removes an IP reservation from the project.
// Deprecated: Use Delete instead.
func (i *ProjectIPServiceOp) Remove(ipReservationID string) (*Response, error) {
	return i.Delete(ipReservationID)
}

// AvailableAddresses lists addresses available from a reserved block
func (i *ProjectIPServiceOp) AvailableAddresses(ipReservationID string, r *AvailableRequest) ([]string, *Response, error) {
	if validateErr := ValidateUUID(ipReservationID); validateErr != nil {
		return nil, nil, validateErr
	}

	opts := &GetOptions{}
	opts = opts.Filter("cidr", strconv.Itoa(r.CIDR))

	endpointPath := path.Join(ipBasePath, ipReservationID, "available")
	apiPathQuery := opts.WithQuery(endpointPath)

	ar := new(AvailableResponse)

	resp, err := i.client.DoRequest("GET", apiPathQuery, r, ar)
	if err != nil {
		return nil, resp, err
	}
	return ar.Available, resp, nil
}
