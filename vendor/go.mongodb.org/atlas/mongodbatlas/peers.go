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

const peersPath = "api/atlas/v1.0/groups/%s/peers"

// PeersService is an interface for interfacing with the Peers
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc/
type PeersService interface {
	List(context.Context, string, *ContainersListOptions) ([]Peer, *Response, error)
	Get(context.Context, string, string) (*Peer, *Response, error)
	Create(context.Context, string, *Peer) (*Peer, *Response, error)
	Update(context.Context, string, string, *Peer) (*Peer, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

// PeersServiceOp handles communication with the Network Peering Connection related methods
// of the MongoDB Atlas API.
type PeersServiceOp service

var _ PeersService = &PeersServiceOp{}

// Peer represents MongoDB peer connection.
type Peer struct {
	AccepterRegionName  string `json:"accepterRegionName,omitempty"`
	AWSAccountID        string `json:"awsAccountId,omitempty"`
	ConnectionID        string `json:"connectionId,omitempty"`
	ContainerID         string `json:"containerId,omitempty"`
	ErrorStateName      string `json:"errorStateName,omitempty"`
	ID                  string `json:"id,omitempty"`
	ProviderName        string `json:"providerName,omitempty"`
	RouteTableCIDRBlock string `json:"routeTableCidrBlock,omitempty"`
	StatusName          string `json:"statusName,omitempty"`
	VpcID               string `json:"vpcId,omitempty"`
	AtlasCIDRBlock      string `json:"atlasCidrBlock,omitempty"`
	AzureDirectoryID    string `json:"azureDirectoryId,omitempty"`
	AzureSubscriptionID string `json:"azureSubscriptionId,omitempty"`
	ResourceGroupName   string `json:"resourceGroupName,omitempty"`
	VNetName            string `json:"vnetName,omitempty"`
	ErrorState          string `json:"errorState,omitempty"`
	Status              string `json:"status,omitempty"`
	GCPProjectID        string `json:"gcpProjectId,omitempty"`
	NetworkName         string `json:"networkName,omitempty"`
	ErrorMessage        string `json:"errorMessage,omitempty"`
}

// peersResponse is the response from the PeersService.List.
type peersResponse struct {
	Links      []*Link `json:"links,omitempty"`
	Results    []Peer  `json:"results,omitempty"`
	TotalCount int     `json:"totalCount,omitempty"`
}

// List all peers in the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc-get-connections-list/
func (s *PeersServiceOp) List(ctx context.Context, groupID string, listOptions *ContainersListOptions) ([]Peer, *Response, error) {
	path := fmt.Sprintf(peersPath, groupID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(peersResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// Get gets the netwprk peering connection specified to {PEER-ID} from the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc-get-connection/
func (s *PeersServiceOp) Get(ctx context.Context, groupID, peerID string) (*Peer, *Response, error) {
	if peerID == "" {
		return nil, nil, NewArgError("perrID", "must be set")
	}

	basePath := fmt.Sprintf(peersPath, groupID)
	escapedEntry := url.PathEscape(peerID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Peer)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create a peer connection to the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc-create-peering-connection/
func (s *PeersServiceOp) Create(ctx context.Context, groupID string, createRequest *Peer) (*Peer, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(peersPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Peer)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update a peer connection in the project associated to {GROUP-ID}
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc-update-peering-connection/
func (s *PeersServiceOp) Update(ctx context.Context, groupID, peerID string, updateRequest *Peer) (*Peer, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(peersPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, peerID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Peer)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete the peer connection specified to {PEER-ID} from the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/vpc-delete-peering-connection/
func (s *PeersServiceOp) Delete(ctx context.Context, groupID, peerID string) (*Response, error) {
	if peerID == "" {
		return nil, NewArgError("peerID", "must be set")
	}

	basePath := fmt.Sprintf(peersPath, groupID)
	escapedEntry := url.PathEscape(peerID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
