package godo

import (
	"context"
	"net/http"
	"time"
)

const vpcPeeringsPath = "/v2/vpc_peerings"

type vpcPeeringRoot struct {
	VPCPeering *VPCPeering `json:"vpc_peering"`
}

type vpcPeeringsRoot struct {
	VPCPeerings []*VPCPeering `json:"vpc_peerings"`
	Links       *Links        `json:"links"`
	Meta        *Meta         `json:"meta"`
}

// VPCPeering represents a DigitalOcean Virtual Private Cloud Peering configuration.
type VPCPeering struct {
	// ID is the generated ID of the VPC Peering
	ID string `json:"id"`
	// Name is the name of the VPC Peering
	Name string `json:"name"`
	// VPCIDs is the IDs of the pair of VPCs between which a peering is created
	VPCIDs []string `json:"vpc_ids"`
	// CreatedAt is time when this VPC Peering was first created
	CreatedAt time.Time `json:"created_at"`
	// Status is the status of the VPC Peering
	Status string `json:"status"`
}

// VPCPeeringCreateRequest represents a request to create a Virtual Private Cloud Peering
// for a list of associated VPC IDs.
type VPCPeeringCreateRequest struct {
	// Name is the name of the VPC Peering
	Name string `json:"name"`
	// VPCIDs is the IDs of the pair of VPCs between which a peering is created
	VPCIDs []string `json:"vpc_ids"`
}

// VPCPeeringUpdateRequest represents a request to update a Virtual Private Cloud Peering.
type VPCPeeringUpdateRequest struct {
	// Name is the name of the VPC Peering
	Name string `json:"name"`
}

// VPCPeeringCreateRequestByVPCID represents a request to create a Virtual Private Cloud Peering
// for an associated VPC ID.
type VPCPeeringCreateRequestByVPCID struct {
	// Name is the name of the VPC Peering
	Name string `json:"name"`
	// VPCID is the ID of one of the VPCs with which the peering has to be created
	VPCID string `json:"vpc_id"`
}

// CreateVPCPeering creates a new Virtual Private Cloud Peering.
func (v *VPCsServiceOp) CreateVPCPeering(ctx context.Context, create *VPCPeeringCreateRequest) (*VPCPeering, *Response, error) {
	path := vpcPeeringsPath
	req, err := v.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, nil, err
	}

	root := new(vpcPeeringRoot)
	resp, err := v.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.VPCPeering, resp, nil
}

// GetVPCPeering retrieves a Virtual Private Cloud Peering.
func (v *VPCsServiceOp) GetVPCPeering(ctx context.Context, id string) (*VPCPeering, *Response, error) {
	path := vpcPeeringsPath + "/" + id
	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(vpcPeeringRoot)
	resp, err := v.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.VPCPeering, resp, nil
}

// ListVPCPeerings lists all Virtual Private Cloud Peerings.
func (v *VPCsServiceOp) ListVPCPeerings(ctx context.Context, opt *ListOptions) ([]*VPCPeering, *Response, error) {
	path, err := addOptions(vpcPeeringsPath, opt)
	if err != nil {
		return nil, nil, err
	}
	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(vpcPeeringsRoot)
	resp, err := v.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}
	return root.VPCPeerings, resp, nil
}

// UpdateVPCPeering updates a Virtual Private Cloud Peering.
func (v *VPCsServiceOp) UpdateVPCPeering(ctx context.Context, id string, update *VPCPeeringUpdateRequest) (*VPCPeering, *Response, error) {
	path := vpcPeeringsPath + "/" + id
	req, err := v.client.NewRequest(ctx, http.MethodPatch, path, update)
	if err != nil {
		return nil, nil, err
	}

	root := new(vpcPeeringRoot)
	resp, err := v.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.VPCPeering, resp, nil
}

// DeleteVPCPeering deletes a Virtual Private Cloud Peering.
func (v *VPCsServiceOp) DeleteVPCPeering(ctx context.Context, id string) (*Response, error) {
	path := vpcPeeringsPath + "/" + id
	req, err := v.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := v.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// CreateVPCPeeringByVPCID creates a new Virtual Private Cloud Peering for requested VPC ID.
func (v *VPCsServiceOp) CreateVPCPeeringByVPCID(ctx context.Context, id string, create *VPCPeeringCreateRequestByVPCID) (*VPCPeering, *Response, error) {
	path := vpcsBasePath + "/" + id + "/peerings"
	req, err := v.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, nil, err
	}

	root := new(vpcPeeringRoot)
	resp, err := v.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.VPCPeering, resp, nil
}

// ListVPCPeeringsByVPCID lists all Virtual Private Cloud Peerings for requested VPC ID.
func (v *VPCsServiceOp) ListVPCPeeringsByVPCID(ctx context.Context, id string, opt *ListOptions) ([]*VPCPeering, *Response, error) {
	path, err := addOptions(vpcsBasePath+"/"+id+"/peerings", opt)
	req, err := v.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(vpcPeeringsRoot)
	resp, err := v.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}
	return root.VPCPeerings, resp, nil
}

// UpdateVPCPeeringByVPCID updates a Virtual Private Cloud Peering for requested VPC ID.
func (v *VPCsServiceOp) UpdateVPCPeeringByVPCID(ctx context.Context, vpcID, peerID string, update *VPCPeeringUpdateRequest) (*VPCPeering, *Response, error) {
	path := vpcsBasePath + "/" + vpcID + "/peerings" + "/" + peerID
	req, err := v.client.NewRequest(ctx, http.MethodPatch, path, update)
	if err != nil {
		return nil, nil, err
	}

	root := new(vpcPeeringRoot)
	resp, err := v.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.VPCPeering, resp, nil
}
