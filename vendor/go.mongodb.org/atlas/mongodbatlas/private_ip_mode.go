// Copyright 2021 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const privateIPModePath = "api/atlas/v1.0/groups/%s/privateIpMode"

// PrivateIPModeService is an interface for interfacing with the PrivateIpMode
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/get-private-ip-mode-for-project/
type PrivateIPModeService interface {
	Get(context.Context, string) (*PrivateIPMode, *Response, error)
	Update(context.Context, string, *PrivateIPMode) (*PrivateIPMode, *Response, error)
}

// PrivateIPModeServiceOp handles communication with the Private IP Mode related methods
// of the MongoDB Atlas API.
type PrivateIPModeServiceOp service

var _ PrivateIPModeService = &PrivateIPModeServiceOp{}

// PrivateIPMode represents MongoDB Private IP Mode Configutation.
type PrivateIPMode struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// Get Verify Connect via Peering Only Mode from the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/get-private-ip-mode-for-project/
func (s *PrivateIPModeServiceOp) Get(ctx context.Context, groupID string) (*PrivateIPMode, *Response, error) {
	path := fmt.Sprintf(privateIPModePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(PrivateIPMode)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update connection via Peering Only Mode in the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/set-private-ip-mode-for-project/
func (s *PrivateIPModeServiceOp) Update(ctx context.Context, groupID string, updateRequest *PrivateIPMode) (*PrivateIPMode, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	path := fmt.Sprintf(privateIPModePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(PrivateIPMode)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
