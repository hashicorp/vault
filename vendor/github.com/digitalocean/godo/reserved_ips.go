package godo

import (
	"context"
	"fmt"
	"net/http"
)

const resourceType = "ReservedIP"
const reservedIPsBasePath = "v2/reserved_ips"

// ReservedIPsService is an interface for interfacing with the reserved IPs
// endpoints of the Digital Ocean API.
// See: https://docs.digitalocean.com/reference/api/api-reference/#tag/Reserved-IPs
type ReservedIPsService interface {
	List(context.Context, *ListOptions) ([]ReservedIP, *Response, error)
	Get(context.Context, string) (*ReservedIP, *Response, error)
	Create(context.Context, *ReservedIPCreateRequest) (*ReservedIP, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

// ReservedIPsServiceOp handles communication with the reserved IPs related methods of the
// DigitalOcean API.
type ReservedIPsServiceOp struct {
	client *Client
}

var _ ReservedIPsService = &ReservedIPsServiceOp{}

// ReservedIP represents a Digital Ocean reserved IP.
type ReservedIP struct {
	Region    *Region  `json:"region"`
	Droplet   *Droplet `json:"droplet"`
	IP        string   `json:"ip"`
	ProjectID string   `json:"project_id"`
	Locked    bool     `json:"locked"`
}

func (f ReservedIP) String() string {
	return Stringify(f)
}

// URN returns the reserved IP in a valid DO API URN form.
func (f ReservedIP) URN() string {
	return ToURN(resourceType, f.IP)
}

type reservedIPsRoot struct {
	ReservedIPs []ReservedIP `json:"reserved_ips"`
	Links       *Links       `json:"links"`
	Meta        *Meta        `json:"meta"`
}

type reservedIPRoot struct {
	ReservedIP *ReservedIP `json:"reserved_ip"`
	Links      *Links      `json:"links,omitempty"`
}

// ReservedIPCreateRequest represents a request to create a reserved IP.
// Specify DropletID to assign the reserved IP to a Droplet or Region
// to reserve it to the region.
type ReservedIPCreateRequest struct {
	Region    string `json:"region,omitempty"`
	DropletID int    `json:"droplet_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

// List all reserved IPs.
func (r *ReservedIPsServiceOp) List(ctx context.Context, opt *ListOptions) ([]ReservedIP, *Response, error) {
	path := reservedIPsBasePath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := r.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(reservedIPsRoot)
	resp, err := r.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.ReservedIPs, resp, err
}

// Get an individual reserved IP.
func (r *ReservedIPsServiceOp) Get(ctx context.Context, ip string) (*ReservedIP, *Response, error) {
	path := fmt.Sprintf("%s/%s", reservedIPsBasePath, ip)

	req, err := r.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(reservedIPRoot)
	resp, err := r.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.ReservedIP, resp, err
}

// Create a reserved IP. If the DropletID field of the request is not empty,
// the reserved IP will also be assigned to the droplet.
func (r *ReservedIPsServiceOp) Create(ctx context.Context, createRequest *ReservedIPCreateRequest) (*ReservedIP, *Response, error) {
	path := reservedIPsBasePath

	req, err := r.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(reservedIPRoot)
	resp, err := r.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.ReservedIP, resp, err
}

// Delete a reserved IP.
func (r *ReservedIPsServiceOp) Delete(ctx context.Context, ip string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", reservedIPsBasePath, ip)

	req, err := r.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(ctx, req, nil)

	return resp, err
}
