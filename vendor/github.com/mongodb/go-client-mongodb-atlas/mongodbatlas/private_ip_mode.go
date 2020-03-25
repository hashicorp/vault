package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const privateIpModePath = "groups/%s/privateIpMode"

//PrivateIpModeService is an interface for interfacing with the PrivateIpMode
// endpoints of the MongoDB Atlas API.
//See more: https://docs.atlas.mongodb.com/reference/api/get-private-ip-mode-for-project/
type PrivateIpModeService interface {
	Get(context.Context, string) (*PrivateIPMode, *Response, error)
	Update(context.Context, string, *PrivateIPMode) (*PrivateIPMode, *Response, error)
}

//PrivateIpModeServiceOp handles communication with the Private IP Mode related methods
// of the MongoDB Atlas API
type PrivateIpModeServiceOp struct {
	client *Client
}

var _ PrivateIpModeService = &PrivateIpModeServiceOp{}

// PrivateIPMode represents MongoDB Private IP Mode Configutation.
type PrivateIPMode struct {
	Enabled *bool `json:"enabled,omitempty"`
}

//Get Verify Connect via Peering Only Mode from the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/get-private-ip-mode-for-project/
func (s *PrivateIpModeServiceOp) Get(ctx context.Context, groupID string) (*PrivateIPMode, *Response, error) {
	path := fmt.Sprintf(privateIpModePath, groupID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(PrivateIPMode)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Update connection via Peering Only Mode in the project associated to {GROUP-ID}
//See more: https://docs.atlas.mongodb.com/reference/api/set-private-ip-mode-for-project/
func (s *PrivateIpModeServiceOp) Update(ctx context.Context, groupID string, updateRequest *PrivateIPMode) (*PrivateIPMode, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	path := fmt.Sprintf(privateIpModePath, groupID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(PrivateIPMode)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
