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
	"net/url"
)

const apiKeysPath = "api/atlas/v1.0/orgs/%s/apiKeys" //nolint:gosec // This is a path

// APIKeysService is an interface for interfacing with the APIKeys
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys/
type APIKeysService interface {
	List(context.Context, string, *ListOptions) ([]APIKey, *Response, error)
	Get(context.Context, string, string) (*APIKey, *Response, error)
	Create(context.Context, string, *APIKeyInput) (*APIKey, *Response, error)
	Update(context.Context, string, string, *APIKeyInput) (*APIKey, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

// APIKeysServiceOp handles communication with the APIKey related methods
// of the MongoDB Atlas API.
type APIKeysServiceOp service

var _ APIKeysService = &APIKeysServiceOp{}

// APIKeyInput represents MongoDB API key input request for Create.
type APIKeyInput struct {
	Desc  string   `json:"desc,omitempty"`
	Roles []string `json:"roles,omitempty"`
}

// APIKey represents MongoDB API Key.
type APIKey struct {
	ID         string      `json:"id,omitempty"`
	Desc       string      `json:"desc,omitempty"`
	Roles      []AtlasRole `json:"roles,omitempty"`
	PrivateKey string      `json:"privateKey,omitempty"`
	PublicKey  string      `json:"publicKey,omitempty"`
}

// AtlasRole represents a role name of API key.
type AtlasRole struct {
	GroupID  string `json:"groupId,omitempty"`
	OrgID    string `json:"orgId,omitempty"`
	RoleName string `json:"roleName,omitempty"`
}

// apiKeysResponse is the response from the APIKeysService.List.
type apiKeysResponse struct {
	Links      []*Link  `json:"links,omitempty"`
	Results    []APIKey `json:"results,omitempty"`
	TotalCount int      `json:"totalCount,omitempty"`
}

// List all API-KEY in the organization associated to {ORG-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-get-all/
func (s *APIKeysServiceOp) List(ctx context.Context, orgID string, listOptions *ListOptions) ([]APIKey, *Response, error) {
	path := fmt.Sprintf(apiKeysPath, orgID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(apiKeysResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// Get gets the APIKey specified to {API-KEY-ID} from the organization associated to {ORG-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-get-one/
func (s *APIKeysServiceOp) Get(ctx context.Context, orgID, apiKeyID string) (*APIKey, *Response, error) {
	if apiKeyID == "" {
		return nil, nil, NewArgError("name", "must be set")
	}

	basePath := fmt.Sprintf(apiKeysPath, orgID)
	escapedEntry := url.PathEscape(apiKeyID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(APIKey)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create an API Key by the {ORG-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-create-one/
func (s *APIKeysServiceOp) Create(ctx context.Context, orgID string, createRequest *APIKeyInput) (*APIKey, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(apiKeysPath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(APIKey)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update a API Key in the organization associated to {ORG-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-update-one/
func (s *APIKeysServiceOp) Update(ctx context.Context, orgID, apiKeyID string, updateRequest *APIKeyInput) (*APIKey, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(apiKeysPath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, apiKeyID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(APIKey)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete the API Key specified to {API-KEY-ID} from the organization associated to {ORG-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/apiKey-delete-one-apiKey/
func (s *APIKeysServiceOp) Delete(ctx context.Context, orgID, apiKeyID string) (*Response, error) {
	if apiKeyID == "" {
		return nil, NewArgError("apiKeyID", "must be set")
	}

	basePath := fmt.Sprintf(apiKeysPath, orgID)
	escapedEntry := url.PathEscape(apiKeyID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
