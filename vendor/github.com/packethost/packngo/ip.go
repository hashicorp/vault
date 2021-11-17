package packngo

import (
	"fmt"
)

const ipBasePath = "/ips"

// DeviceIPService handles assignment of addresses from reserved blocks to instances in a project.
type DeviceIPService interface {
	Assign(deviceID string, assignRequest *AddressStruct) (*IPAddressAssignment, *Response, error)
	Unassign(assignmentID string) (*Response, error)
	Get(assignmentID string) (*IPAddressAssignment, *Response, error)
}

// ProjectIPService handles reservation of IP address blocks for a project.
type ProjectIPService interface {
	Get(reservationID string) (*IPAddressReservation, *Response, error)
	List(projectID string) ([]IPAddressReservation, *Response, error)
	Request(projectID string, ipReservationReq *IPReservationRequest) (*IPAddressReservation, *Response, error)
	Remove(ipReservationID string) (*Response, error)
	AvailableAddresses(ipReservationID string, r *AvailableRequest) ([]string, *Response, error)
}

type ipAddressCommon struct {
	ID            string `json:"id"`
	Address       string `json:"address"`
	Gateway       string `json:"gateway"`
	Network       string `json:"network"`
	AddressFamily int    `json:"address_family"`
	Netmask       string `json:"netmask"`
	Public        bool   `json:"public"`
	CIDR          int    `json:"cidr"`
	Created       string `json:"created_at,omitempty"`
	Updated       string `json:"updated_at,omitempty"`
	Href          string `json:"href"`
	Management    bool   `json:"management"`
	Manageable    bool   `json:"manageable"`
	Project       Href   `json:"project"`
}

// IPAddressReservation is created when user sends IP reservation request for a project (considering it's within quota).
type IPAddressReservation struct {
	ipAddressCommon
	Assignments []Href   `json:"assignments"`
	Facility    Facility `json:"facility,omitempty"`
	Available   string   `json:"available"`
	Addon       bool     `json:"addon"`
	Bill        bool     `json:"bill"`
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
	ipAddressCommon
	AssignedTo Href `json:"assigned_to"`
}

// IPReservationRequest represents the body of a reservation request.
type IPReservationRequest struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
	Comments string `json:"comments"`
	Facility string `json:"facility"`
}

// AddressStruct is a helper type for request/response with dict like {"address": ... }
type AddressStruct struct {
	Address string `json:"address"`
}

func deleteFromIP(client *Client, resourceID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", ipBasePath, resourceID)

	return client.DoRequest("DELETE", path, nil, nil)
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
	path := fmt.Sprintf("%s/%s%s", deviceBasePath, deviceID, ipBasePath)
	ipa := new(IPAddressAssignment)

	resp, err := i.client.DoRequest("POST", path, assignRequest, ipa)
	if err != nil {
		return nil, resp, err
	}

	return ipa, resp, err
}

// Get returns assignment by ID.
func (i *DeviceIPServiceOp) Get(assignmentID string) (*IPAddressAssignment, *Response, error) {
	path := fmt.Sprintf("%s/%s", ipBasePath, assignmentID)
	ipa := new(IPAddressAssignment)

	resp, err := i.client.DoRequest("GET", path, nil, ipa)
	if err != nil {
		return nil, resp, err
	}

	return ipa, resp, err
}

// ProjectIPServiceOp is interface for IP assignment methods.
type ProjectIPServiceOp struct {
	client *Client
}

// Get returns reservation by ID.
func (i *ProjectIPServiceOp) Get(reservationID string) (*IPAddressReservation, *Response, error) {
	path := fmt.Sprintf("%s/%s", ipBasePath, reservationID)
	ipr := new(IPAddressReservation)

	resp, err := i.client.DoRequest("GET", path, nil, ipr)
	if err != nil {
		return nil, resp, err
	}

	return ipr, resp, err
}

// List provides a list of IP resevations for a single project.
func (i *ProjectIPServiceOp) List(projectID string) ([]IPAddressReservation, *Response, error) {
	path := fmt.Sprintf("%s/%s%s", projectBasePath, projectID, ipBasePath)
	reservations := new(struct {
		Reservations []IPAddressReservation `json:"ip_addresses"`
	})

	resp, err := i.client.DoRequest("GET", path, nil, reservations)
	if err != nil {
		return nil, resp, err
	}
	return reservations.Reservations, resp, nil
}

// Request requests more IP space for a project in order to have additional IP addresses to assign to devices.
func (i *ProjectIPServiceOp) Request(projectID string, ipReservationReq *IPReservationRequest) (*IPAddressReservation, *Response, error) {
	path := fmt.Sprintf("%s/%s%s", projectBasePath, projectID, ipBasePath)
	ipr := new(IPAddressReservation)

	resp, err := i.client.DoRequest("POST", path, ipReservationReq, ipr)
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
	path := fmt.Sprintf("%s/%s/available?cidr=%d", ipBasePath, ipReservationID, r.CIDR)
	ar := new(AvailableResponse)

	resp, err := i.client.DoRequest("GET", path, r, ar)
	if err != nil {
		return nil, resp, err
	}
	return ar.Available, resp, nil

}
