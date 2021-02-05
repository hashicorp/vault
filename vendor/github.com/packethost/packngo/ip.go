package packngo

import (
	"fmt"
	"path"
)

const ipBasePath = "/ips"

const (
	// PublicIPv4 fixed string representation of public ipv4
	PublicIPv4 = "public_ipv4"
	// PrivateIPv4 fixed string representation of private ipv4
	PrivateIPv4 = "private_ipv4"
	// GlobalIPv4 fixed string representation of global ipv4
	GlobalIPv4 = "global_ipv4"
	// PublicIPv6 fixed string representation of public ipv6
	PublicIPv6 = "public_ipv6"
	// PrivateIPv6 fixed string representation of private ipv6
	PrivateIPv6 = "private_ipv6"
	// GlobalIPv6 fixed string representation of global ipv6
	GlobalIPv6 = "global_ipv6"
)

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
	Request(projectID string, ipReservationReq *IPReservationRequest) (*IPAddressReservation, *Response, error)
	Remove(ipReservationID string) (*Response, error)
	AvailableAddresses(ipReservationID string, r *AvailableRequest) ([]string, *Response, error)
}

type IpAddressCommon struct { //nolint:golint
	ID            string      `json:"id"`
	Address       string      `json:"address"`
	Gateway       string      `json:"gateway"`
	Network       string      `json:"network"`
	AddressFamily int         `json:"address_family"`
	Netmask       string      `json:"netmask"`
	Public        bool        `json:"public"`
	CIDR          int         `json:"cidr"`
	Created       string      `json:"created_at,omitempty"`
	Updated       string      `json:"updated_at,omitempty"`
	Href          string      `json:"href"`
	Management    bool        `json:"management"`
	Manageable    bool        `json:"manageable"`
	Project       Href        `json:"project"`
	Global        *bool       `json:"global_ip"`
	Tags          []string    `json:"tags,omitempty"`
	CustomData    interface{} `json:"customdata,omitempty"`
}

// IPAddressReservation is created when user sends IP reservation request for a project (considering it's within quota).
type IPAddressReservation struct {
	IpAddressCommon
	Assignments []*IPAddressAssignment `json:"assignments"`
	Facility    *Facility              `json:"facility,omitempty"`
	Available   string                 `json:"available"`
	Addon       bool                   `json:"addon"`
	Bill        bool                   `json:"bill"`
	Description *string                `json:"details"`
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

// IPReservationRequest represents the body of a reservation request.
type IPReservationRequest struct {
	Type        string      `json:"type"`
	Quantity    int         `json:"quantity"`
	Description string      `json:"details,omitempty"`
	Facility    *string     `json:"facility,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	CustomData  interface{} `json:"customdata,omitempty"`
	// FailOnApprovalRequired if the IP request cannot be approved automatically, rather than sending to
	// the longer Equinix Metal approval process, fail immediately with a 422 error
	FailOnApprovalRequired bool `json:"fail_on_approval_required,omitempty"`
}

// AddressStruct is a helper type for request/response with dict like {"address": ... }
type AddressStruct struct {
	Address string `json:"address"`
}

func deleteFromIP(client *Client, resourceID string) (*Response, error) {
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
	return deleteFromIP(i.client, assignmentID)
}

// Assign assigns an IP address to a device.
// The IP address must be in one of the IP ranges assigned to the device’s project.
func (i *DeviceIPServiceOp) Assign(deviceID string, assignRequest *AddressStruct) (*IPAddressAssignment, *Response, error) {
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
	endpointPath := path.Join(deviceBasePath, deviceID, ipBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	//ipList represents collection of IP Address reservations
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
func (i *ProjectIPServiceOp) List(projectID string, opts *ListOptions) ([]IPAddressReservation, *Response, error) {
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

// Request requests more IP space for a project in order to have additional IP addresses to assign to devices.
func (i *ProjectIPServiceOp) Request(projectID string, ipReservationReq *IPReservationRequest) (*IPAddressReservation, *Response, error) {
	apiPath := path.Join(projectBasePath, projectID, ipBasePath)
	ipr := new(IPAddressReservation)

	resp, err := i.client.DoRequest("POST", apiPath, ipReservationReq, ipr)
	if err != nil {
		return nil, resp, err
	}
	return ipr, resp, err
}

// Remove removes an IP reservation from the project.
func (i *ProjectIPServiceOp) Remove(ipReservationID string) (*Response, error) {
	return deleteFromIP(i.client, ipReservationID)
}

// AvailableAddresses lists addresses available from a reserved block
func (i *ProjectIPServiceOp) AvailableAddresses(ipReservationID string, r *AvailableRequest) ([]string, *Response, error) {
	apiPathQuery := fmt.Sprintf("%s/%s/available?cidr=%d", ipBasePath, ipReservationID, r.CIDR)
	ar := new(AvailableResponse)

	resp, err := i.client.DoRequest("GET", apiPathQuery, r, ar)
	if err != nil {
		return nil, resp, err
	}
	return ar.Available, resp, nil

}
