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

const (
	serverlessPrivateEndpointsPath = "api/atlas/v1.0/groups/%s/privateEndpoint/serverless/instance/%s/endpoint"
)

// ServerlessPrivateEndpointsService is an interface for interfacing with the Private Endpoints
// of the MongoDB Atlas API.
//
// See more: See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Serverless-Private-Endpoints
type ServerlessPrivateEndpointsService interface {
	List(context.Context, string, string, *ListOptions) ([]ServerlessPrivateEndpointConnection, *Response, error)
	Create(context.Context, string, string, *ServerlessPrivateEndpointConnection) (*ServerlessPrivateEndpointConnection, *Response, error)
	Get(context.Context, string, string, string) (*ServerlessPrivateEndpointConnection, *Response, error)
	Delete(context.Context, string, string, string) (*Response, error)
	Update(context.Context, string, string, string, *ServerlessPrivateEndpointConnection) (*ServerlessPrivateEndpointConnection, *Response, error)
}

// PrivateServerlessEndpointsServiceOp handles communication with the PrivateServerlessEndpoints related methods
// of the MongoDB Atlas API.
type ServerlessPrivateEndpointsServiceOp service

var _ ServerlessPrivateEndpointsService = &ServerlessPrivateEndpointsServiceOp{}

// PrivateEndpointServerlessConnection represents MongoDB Private Endpoint Connection.
type ServerlessPrivateEndpointConnection struct {
	ID                           string `json:"_id,omitempty"` // Unique identifier of the Serverless PrivateLink Service.
	CloudProviderEndpointID      string `json:"cloudProviderEndpointId,omitempty"`
	Comment                      string `json:"comment,omitempty"`
	EndpointServiceName          string `json:"endpointServiceName,omitempty"`          // Name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.
	ErrorMessage                 string `json:"errorMessage,omitempty"`                 // Error message pertaining to the AWS Service Connect. Returns null if there are no errors.
	Status                       string `json:"status,omitempty"`                       // Status of the AWS Serverless PrivateLink connection: INITIATING, WAITING_FOR_USER, FAILED, DELETING, AVAILABLE.
	ProviderName                 string `json:"providerName,omitempty"`                 // Human-readable label that identifies the cloud provider. Values include AWS or AZURE.
	PrivateEndpointIPAddress     string `json:"privateEndpointIpAddress,omitempty"`     // IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
	PrivateLinkServiceResourceID string `json:"privateLinkServiceResourceId,omitempty"` // Root-relative path that identifies the Azure Private Link Service that MongoDB Cloud manages. MongoDB Cloud returns null while it creates the endpoint service.
}

// List retrieve details for all private Serverless endpoint services in one Atlas project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#operation/returnAllPrivateEndpointsForOneServerlessInstance
func (s *ServerlessPrivateEndpointsServiceOp) List(ctx context.Context, groupID, instanceName string, listOptions *ListOptions) ([]ServerlessPrivateEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if instanceName == "" {
		return nil, nil, NewArgError("instanceID", "must be set")
	}

	path := fmt.Sprintf(serverlessPrivateEndpointsPath, groupID, instanceName) // Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new([]ServerlessPrivateEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return *root, resp, nil
}

// Delete one private serverless endpoint service in an Atlas project.
//
// See more https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#operation/removeOnePrivateEndpointFromOneServerlessInstance
func (s *ServerlessPrivateEndpointsServiceOp) Delete(ctx context.Context, groupID, instanceName, privateEndpointID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}
	if privateEndpointID == "" {
		return nil, NewArgError("PrivateEndpointID", "must be set")
	}
	if instanceName == "" {
		return nil, NewArgError("instanceName", "must be set")
	}

	basePath := fmt.Sprintf(serverlessPrivateEndpointsPath, groupID, instanceName)
	path := fmt.Sprintf("%s/%s", basePath, url.PathEscape(privateEndpointID))

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.Client.Do(ctx, req, nil)
}

// Create Adds one serverless  private endpoint in an Atlas project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#operation/createOnePrivateEndpointForOneServerlessInstance
func (s *ServerlessPrivateEndpointsServiceOp) Create(ctx context.Context, groupID, instanceName string, createRequest *ServerlessPrivateEndpointConnection) (*ServerlessPrivateEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if instanceName == "" {
		return nil, nil, NewArgError("instanceName", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(serverlessPrivateEndpointsPath, groupID, instanceName)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(ServerlessPrivateEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Get retrieve details for one private serverless endpoint in an Atlas project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#operation/returnOnePrivateEndpointForOneServerlessInstance
func (s *ServerlessPrivateEndpointsServiceOp) Get(ctx context.Context, groupID, instanceName, privateEndpointID string) (*ServerlessPrivateEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if instanceName == "" {
		return nil, nil, NewArgError("instanceName", "must be set")
	}
	if privateEndpointID == "" {
		return nil, nil, NewArgError("privateEndpointID", "must be set")
	}

	basePath := fmt.Sprintf(serverlessPrivateEndpointsPath, groupID, instanceName)
	path := fmt.Sprintf("%s/%s", basePath, url.PathEscape(privateEndpointID))

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ServerlessPrivateEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update updates the private serverless endpoint setting for one Atlas project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#operation/updateOnePrivateEndpointForOneServerlessInstance
func (s *ServerlessPrivateEndpointsServiceOp) Update(ctx context.Context, groupID, instanceName, privateEndpointID string, updateRequest *ServerlessPrivateEndpointConnection) (*ServerlessPrivateEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if instanceName == "" {
		return nil, nil, NewArgError("instanceName", "must be set")
	}
	if privateEndpointID == "" {
		return nil, nil, NewArgError("privateEndpointID", "must be set")
	}

	basePath := fmt.Sprintf(serverlessPrivateEndpointsPath, groupID, instanceName)
	path := fmt.Sprintf("%s/%s", basePath, url.PathEscape(privateEndpointID))
	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(ServerlessPrivateEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
