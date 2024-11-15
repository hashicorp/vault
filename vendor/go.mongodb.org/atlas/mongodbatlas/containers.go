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

const containersPath = "api/atlas/v1.0/groups/%s/containers"

// ContainersService provides access to the network peering containers related functions in the Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc/
type ContainersService interface {
	List(context.Context, string, *ContainersListOptions) ([]Container, *Response, error)
	ListAll(context.Context, string, *ListOptions) ([]Container, *Response, error)
	Get(context.Context, string, string) (*Container, *Response, error)
	Create(context.Context, string, *Container) (*Container, *Response, error)
	Update(context.Context, string, string, *Container) (*Container, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

// ContainersServiceOp handles communication with the Network Peering Container related methods
// of the MongoDB Atlas API.
type ContainersServiceOp service

var _ ContainersService = &ContainersServiceOp{}

// ContainersListOptions filtering options for containers.
type ContainersListOptions struct {
	ProviderName string `url:"providerName,omitempty"`
	ListOptions
}

// Container represents MongoDB network peering containter.
type Container struct {
	AtlasCIDRBlock      string   `json:"atlasCidrBlock,omitempty"`
	AzureSubscriptionID string   `json:"azureSubscriptionId,omitempty"`
	GCPProjectID        string   `json:"gcpProjectId,omitempty"`
	ID                  string   `json:"id,omitempty"`
	NetworkName         string   `json:"networkName,omitempty"`
	ProviderName        string   `json:"providerName,omitempty"`
	Provisioned         *bool    `json:"provisioned,omitempty"`
	Region              string   `json:"region,omitempty"`     // Region is available for AZURE
	Regions             []string `json:"regions,omitempty"`    // Regions are available for GCP
	RegionName          string   `json:"regionName,omitempty"` // RegionName is available for AWS
	VNetName            string   `json:"vnetName,omitempty"`
	VPCID               string   `json:"vpcId,omitempty"`
}

// containersResponse is the response from the ContainersService.List.
type containersResponse struct {
	Links      []*Link     `json:"links,omitempty"`
	Results    []Container `json:"results,omitempty"`
	TotalCount int         `json:"totalCount,omitempty"`
}

// List gets details for all network peering containers in an Atlas project for a single cloud provider.
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc-get-containers-list/
func (s *ContainersServiceOp) List(ctx context.Context, groupID string, listOptions *ContainersListOptions) ([]Container, *Response, error) {
	path := fmt.Sprintf(containersPath, groupID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(containersResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// ListAll gets details for all network peering containers in an Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc-get-containers-list-all/
func (s *ContainersServiceOp) ListAll(ctx context.Context, groupID string, listOptions *ListOptions) ([]Container, *Response, error) {
	basePath := fmt.Sprintf(containersPath, groupID)
	path := fmt.Sprintf("%s/all", basePath)
	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(containersResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// Get gets details for one network peering container in an Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc-get-container/
func (s *ContainersServiceOp) Get(ctx context.Context, groupID, containerID string) (*Container, *Response, error) {
	if containerID == "" {
		return nil, nil, NewArgError("perrID", "must be set")
	}

	basePath := fmt.Sprintf(containersPath, groupID)
	escapedEntry := url.PathEscape(containerID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Container)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create adds a network peering container to the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc-create-container/
func (s *ContainersServiceOp) Create(ctx context.Context, groupID string, createRequest *Container) (*Container, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(containersPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Container)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update a network peering container in the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc-update-container/
func (s *ContainersServiceOp) Update(ctx context.Context, groupID, containerID string, updateRequest *Container) (*Container, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(containersPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, containerID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Container)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete the network peering container specified to {CONTAINER-ID} from the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc-delete-one-container/
func (s *ContainersServiceOp) Delete(ctx context.Context, groupID, containerID string) (*Response, error) {
	if containerID == "" {
		return nil, NewArgError("containerID", "must be set")
	}

	basePath := fmt.Sprintf(containersPath, groupID)
	escapedEntry := url.PathEscape(containerID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	// To avoid API Issues
	req.Header.Del("Content-Type")

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
