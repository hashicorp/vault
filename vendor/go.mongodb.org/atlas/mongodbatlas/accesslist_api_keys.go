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

const accessListAPIKeysPath = "api/atlas/v1.0/orgs/%s/apiKeys/%s/accessList" //nolint:gosec // This is a path

// AccessListAPIKeysService is an interface for interfacing with the AccessList API Keys
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/apiKeys#organization-api-key-access-list-endpoints
type AccessListAPIKeysService interface {
	List(context.Context, string, string, *ListOptions) (*AccessListAPIKeys, *Response, error)
	Get(context.Context, string, string, string) (*AccessListAPIKey, *Response, error)
	Create(context.Context, string, string, []*AccessListAPIKeysReq) (*AccessListAPIKeys, *Response, error)
	Delete(context.Context, string, string, string) (*Response, error)
}

// AccessListAPIKeysServiceOp handles communication with the AccessList API keys related methods of the
// MongoDB Atlas API.
type AccessListAPIKeysServiceOp service

var _ AccessListAPIKeysService = &AccessListAPIKeysServiceOp{}

// AccessListAPIKey represents a AccessList API key.
type AccessListAPIKey struct {
	CidrBlock       string  `json:"cidrBlock,omitempty"`       // CIDR-notated range of permitted IP addresses.
	Count           int     `json:"count,omitempty"`           // Total number of requests that have originated from this IP address.
	Created         string  `json:"created,omitempty"`         // Date this IP address was added to the access list.
	IPAddress       string  `json:"ipAddress,omitempty"`       // IP address in the API access list.
	LastUsed        string  `json:"lastUsed,omitempty"`        // Timestamp in ISO 8601 date and time format in UTC when the most recent request that originated from this IP address. This parameter only appears if at least one request has originated from this IP address, and is only updated when a permitted resource is accessed.
	LastUsedAddress string  `json:"lastUsedAddress,omitempty"` // IP address from which the last call to the API was issued. This field only appears if at least one request has originated from this IP address.
	Links           []*Link `json:"links,omitempty"`           // An array of documents, representing a link to one or more sub-resources and/or related resources such as list pagination. See Linking for more information.}
}

// AccessListAPIKeys represents all AccessList API keys.
type AccessListAPIKeys struct {
	Results    []*AccessListAPIKey `json:"results,omitempty"`    // Includes one AccessListAPIKey object for each item detailed in the results array section.
	Links      []*Link             `json:"links,omitempty"`      // One or more links to sub-resources and/or related resources.
	TotalCount int                 `json:"totalCount,omitempty"` // Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.
}

// AccessListAPIKeysReq represents the request to the mehtod create.
type AccessListAPIKeysReq struct {
	IPAddress string `json:"ipAddress,omitempty"` // IP address to be added to the access list for the API key.
	CidrBlock string `json:"cidrBlock,omitempty"` // CIDR-notation block of IP addresses to be added to the access list for the API key.
}

// List gets all AccessList API keys.
//
// See more: https://docs.atlas.mongodb.com/reference/api/api-access-list/get-all-api-access-entries/
func (s *AccessListAPIKeysServiceOp) List(ctx context.Context, orgID, apiKeyID string, listOptions *ListOptions) (*AccessListAPIKeys, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, nil, NewArgError("apiKeyID", "must be set")
	}

	path := fmt.Sprintf(accessListAPIKeysPath, orgID, apiKeyID)
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AccessListAPIKeys)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Get retrieves information on a single API Key access list entry using the unique identifier for the API Key and desired permitted address.
//
// See more: https://docs.atlas.mongodb.com/reference/api/api-access-list/get-one-api-access-entry/
func (s *AccessListAPIKeysServiceOp) Get(ctx context.Context, orgID, apiKeyID, ipAddress string) (*AccessListAPIKey, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, nil, NewArgError("apiKeyID", "must be set")
	}
	if ipAddress == "" {
		return nil, nil, NewArgError("ipAddress", "must be set")
	}

	path := fmt.Sprintf(accessListAPIKeysPath+"/%s", orgID, apiKeyID, ipAddress)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AccessListAPIKey)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates one or more new access list entries for the specified API Key.
//
// See more: https://docs.atlas.mongodb.com/reference/api/api-access-list/create-api-access-entries/
func (s *AccessListAPIKeysServiceOp) Create(ctx context.Context, orgID, apiKeyID string, createRequest []*AccessListAPIKeysReq) (*AccessListAPIKeys, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, nil, NewArgError("apiKeyID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(accessListAPIKeysPath, orgID, apiKeyID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(AccessListAPIKeys)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes the AccessList API keys.
//
// See more: https://docs.atlas.mongodb.com/reference/api/api-access-list/delete-one-api-access-entry/
func (s *AccessListAPIKeysServiceOp) Delete(ctx context.Context, orgID, apiKeyID, ipAddress string) (*Response, error) {
	if orgID == "" {
		return nil, NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, NewArgError("apiKeyID", "must be set")
	}
	if ipAddress == "" {
		return nil, NewArgError("ipAddress", "must be set")
	}

	path := fmt.Sprintf(accessListAPIKeysPath+"/%s", orgID, apiKeyID, ipAddress)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
