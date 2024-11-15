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
	"net/http"
)

const rootPath = "api/atlas/v1.0"

// RootService is an interface for interfacing with the Root
// endpoints of the MongoDB Atlas API.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Root/operation/getSystemStatus
type RootService interface {
	List(context.Context, *ListOptions) (*Root, *Response, error)
}

// RootServiceOp handles communication with the APIKey related methods
// of the MongoDB Atlas API.
type RootServiceOp service
type Root struct {
	APIKey struct {
		AccessList []struct {
			CIDRBlock string `json:"cidrBlock"`
			IPAddress string `json:"ipAddress"`
		} `json:"accessList"`
		ID        string      `json:"id"`
		PublicKey string      `json:"publicKey"`
		Roles     []AtlasRole `json:"roles,omitempty"`
	} `json:"apiKey"`
	AppName    string  `json:"appName"`
	Build      string  `json:"build"`
	Links      []*Link `json:"links"`
	Throttling bool    `json:"throttling"`
}

var _ RootService = &RootServiceOp{}

// List all API-KEY related data
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Root/operation/getSystemStatus
func (s *RootServiceOp) List(ctx context.Context, listOptions *ListOptions) (*Root, *Response, error) {
	path := rootPath

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Root)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
