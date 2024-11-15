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
	privateEndpointsPath = "api/atlas/v1.0/groups/%s/privateEndpoint"
	regionalModePath     = privateEndpointsPath + "/regionalMode"
)

// PrivateEndpointsService is an interface for interfacing with the Private Endpoints
// of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoints/
type PrivateEndpointsService interface {
	Create(context.Context, string, *PrivateEndpointConnection) (*PrivateEndpointConnection, *Response, error)
	Get(context.Context, string, string, string) (*PrivateEndpointConnection, *Response, error)
	List(context.Context, string, string, *ListOptions) ([]PrivateEndpointConnection, *Response, error)
	Delete(context.Context, string, string, string) (*Response, error)
	AddOnePrivateEndpoint(context.Context, string, string, string, *InterfaceEndpointConnection) (*InterfaceEndpointConnection, *Response, error)
	GetOnePrivateEndpoint(context.Context, string, string, string, string) (*InterfaceEndpointConnection, *Response, error)
	DeleteOnePrivateEndpoint(context.Context, string, string, string, string) (*Response, error)
	UpdateRegionalizedPrivateEndpointSetting(context.Context, string, bool) (*RegionalizedPrivateEndpointSetting, *Response, error)
	GetRegionalizedPrivateEndpointSetting(context.Context, string) (*RegionalizedPrivateEndpointSetting, *Response, error)
}

// PrivateEndpointsServiceOp handles communication with the PrivateEndpoints related methods
// of the MongoDB Atlas API.
type PrivateEndpointsServiceOp service

var _ PrivateEndpointsService = &PrivateEndpointsServiceOp{}

// PrivateEndpointConnection represents MongoDB Private Endpoint Connection.
type PrivateEndpointConnection struct {
	ID                           string   `json:"id,omitempty"`                           // Unique identifier of the AWS PrivateLink connection or Azure Private Link Service.
	ProviderName                 string   `json:"providerName,omitempty"`                 // Name of the cloud provider for which you want to create the private endpoint service. Atlas accepts AWS, AZURE or GCP.
	Region                       string   `json:"region,omitempty"`                       // Cloud provider region for which you want to create the private endpoint service.
	EndpointServiceName          string   `json:"endpointServiceName,omitempty"`          // Name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.
	ErrorMessage                 string   `json:"errorMessage,omitempty"`                 // Error message pertaining to the AWS PrivateLink connection or Azure Private Link Service or GCP Private Service Connect. Returns null if there are no errors.
	InterfaceEndpoints           []string `json:"interfaceEndpoints,omitempty"`           // Unique identifiers of the interface endpoints in your VPC that you added to the AWS PrivateLink connection.
	PrivateEndpoints             []string `json:"privateEndpoints,omitempty"`             // All private endpoints that you have added to this Azure Private Link Service.
	PrivateLinkServiceName       string   `json:"privateLinkServiceName,omitempty"`       // Name of the Azure Private Link Service that Atlas manages.
	PrivateLinkServiceResourceID string   `json:"privateLinkServiceResourceId,omitempty"` // Resource ID of the Azure Private Link Service that Atlas manages.
	Status                       string   `json:"status,omitempty"`                       // Status of the AWS, Azure OR GCP PrivateLink connection: INITIATING, WAITING_FOR_USER, FAILED, DELETING, AVAILABLE.
	EndpointGroupNames           []string `json:"endpointGroupNames,omitempty"`           // GCP network endpoint groups corresponding to the Private Service Connect endpoint service.
	RegionName                   string   `json:"regionName,omitempty"`                   // GCP region name for the Private Service Connect endpoint service.
	ServiceAttachmentNames       []string `json:"serviceAttachmentNames,omitempty"`       // Unique alphanumeric and special character strings that identify the service attachments associated with the GCP Private Service Connect endpoint service. Returns an empty list while Atlas creates the service attachments.
}

// InterfaceEndpointConnection represents MongoDB Interface Endpoint Connection.
type InterfaceEndpointConnection struct {
	ID                            string         `json:"id,omitempty"`                            // Unique identifier of the private endpoint you created in your AWS VPC or Azure VNet.
	InterfaceEndpointID           string         `json:"interfaceEndpointId,omitempty"`           // Unique identifier of the interface endpoint.
	PrivateEndpointConnectionName string         `json:"privateEndpointConnectionName,omitempty"` // Name of the connection for this private endpoint that Atlas generates.
	PrivateEndpointIPAddress      string         `json:"privateEndpointIPAddress,omitempty"`      // Private IP address of the private endpoint network interface.
	PrivateEndpointResourceID     string         `json:"privateEndpointResourceId,omitempty"`     // Unique identifier of the private endpoint.
	DeleteRequested               *bool          `json:"deleteRequested,omitempty"`               // Indicates if Atlas received a request to remove the interface endpoint from the private endpoint connection.
	ErrorMessage                  string         `json:"errorMessage,omitempty"`                  // Error message pertaining to the interface endpoint. Returns null if there are no errors.
	AWSConnectionStatus           string         `json:"connectionStatus,omitempty"`              // Status of the interface endpoint: NONE, PENDING_ACCEPTANCE, PENDING, AVAILABLE, REJECTED, DELETING.
	Status                        string         `json:"status,omitempty"`                        // Status of the interface endpoint AZURE: INITIATING, AVAILABLE, FAILED, DELETING. GCP: INITIATING, AVAILABLE, FAILED, DELETING, VERIFIED
	EndpointGroupName             string         `json:"endpointGroupName,omitempty"`             // Unique identifier of the endpoint group. The endpoint group encompasses all of the endpoints that you created in GCP.
	GCPProjectID                  string         `json:"gcpProjectId,omitempty"`                  // Unique identifier of the GCP project in which you created your endpoints.
	Endpoints                     []*GCPEndpoint `json:"endpoints,omitempty"`                     // Collection of individual private endpoints that comprise your endpoint group.
}

// RegionalizedPrivateEndpointSetting represents MongoDB Regionalized private Endpoint Setting.
type RegionalizedPrivateEndpointSetting struct {
	Enabled bool `json:"enabled"` // Flag that indicates whether the regionalized private endpoint setting is enabled for one Atlas project.
}

// GCPEndpoint represents MongoDB GCP endpoint group.
type GCPEndpoint struct {
	IPAddress             string `json:"ipAddress,omitempty"`             // Private IP address of the endpoint you created in GCP.
	EndpointName          string `json:"endpointName,omitempty"`          // Forwarding rule that corresponds to the endpoint you created in GCP.
	Status                string `json:"status,omitempty"`                // Status of the endpoint. Atlas returns one of the values shown above.
	ServiceAttachmentName string `json:"serviceAttachmentName,omitempty"` // Unique alphanumeric and special character strings that identify the service attachment associated with the endpoint.
}

// Create one private endpoint service for AWS or Azure in an Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoints-service-create-one/
func (s *PrivateEndpointsServiceOp) Create(ctx context.Context, groupID string, createRequest *PrivateEndpointConnection) (*PrivateEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(privateEndpointsPath, groupID)
	path := fmt.Sprintf("%s/endpointService", basePath)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(PrivateEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Get retrieve details for one private endpoint service for AWS or Azure in an Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoints-service-get-one/
func (s *PrivateEndpointsServiceOp) Get(ctx context.Context, groupID, cloudProvider, endpointServiceID string) (*PrivateEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if endpointServiceID == "" {
		return nil, nil, NewArgError("endpointServiceID", "must be set")
	}
	if cloudProvider == "" {
		return nil, nil, NewArgError("cloudProvider", "must be set")
	}

	basePath := fmt.Sprintf(privateEndpointsPath, groupID)
	path := fmt.Sprintf("%s/%s/endpointService/%s", basePath, cloudProvider, endpointServiceID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(PrivateEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// List retrieve details for all private endpoint services for AWS or Azure in one Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoints-service-get-all/
func (s *PrivateEndpointsServiceOp) List(ctx context.Context, groupID, cloudProvider string, listOptions *ListOptions) ([]PrivateEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if cloudProvider == "" {
		return nil, nil, NewArgError("cloudProvider", "must be set")
	}

	basePath := fmt.Sprintf(privateEndpointsPath, groupID)
	path := fmt.Sprintf("%s/%s/endpointService", basePath, cloudProvider)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new([]PrivateEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return *root, resp, nil
}

// Delete one private endpoint service for AWS or Azure in an Atlas project.
//
// See more https://docs.atlas.mongodb.com/reference/api/private-endpoints-service-delete-one/
func (s *PrivateEndpointsServiceOp) Delete(ctx context.Context, groupID, cloudProvider, endpointServiceID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}
	if endpointServiceID == "" {
		return nil, NewArgError("endpointServiceID", "must be set")
	}
	if cloudProvider == "" {
		return nil, NewArgError("cloudProvider", "must be set")
	}

	basePath := fmt.Sprintf(privateEndpointsPath, groupID)
	path := fmt.Sprintf("%s/%s/endpointService/%s", basePath, cloudProvider, endpointServiceID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.Client.Do(ctx, req, nil)
}

// AddOnePrivateEndpoint Adds one private endpoint for AWS or Azure in an Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoints-endpoint-create-one/
func (s *PrivateEndpointsServiceOp) AddOnePrivateEndpoint(ctx context.Context, groupID, cloudProvider, endpointServiceID string, createRequest *InterfaceEndpointConnection) (*InterfaceEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if endpointServiceID == "" {
		return nil, nil, NewArgError("endpointServiceID", "must be set")
	}
	if cloudProvider == "" {
		return nil, nil, NewArgError("cloudProvider", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(privateEndpointsPath, groupID)
	path := fmt.Sprintf("%s/%s/endpointService/%s/endpoint", basePath, cloudProvider, endpointServiceID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(InterfaceEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetOnePrivateEndpoint retrieve details for one private endpoint for AWS or Azure in an Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoints-endpoint-get-one/
func (s *PrivateEndpointsServiceOp) GetOnePrivateEndpoint(ctx context.Context, groupID, cloudProvider, endpointServiceID, privateEndpointID string) (*InterfaceEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if endpointServiceID == "" {
		return nil, nil, NewArgError("endpointServiceID", "must be set")
	}
	if cloudProvider == "" {
		return nil, nil, NewArgError("cloudProvider", "must be set")
	}
	if privateEndpointID == "" {
		return nil, nil, NewArgError("privateEndpointID", "must be set")
	}

	basePath := fmt.Sprintf(privateEndpointsPath, groupID)
	path := fmt.Sprintf("%s/%s/endpointService/%s/endpoint/%s", basePath, cloudProvider, endpointServiceID, url.PathEscape(privateEndpointID))

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(InterfaceEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// DeleteOnePrivateEndpoint remove one private endpoint for AWS or Azure from an Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoints-endpoint-delete-one/
func (s *PrivateEndpointsServiceOp) DeleteOnePrivateEndpoint(ctx context.Context, groupID, cloudProvider, endpointServiceID, privateEndpointID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}
	if endpointServiceID == "" {
		return nil, NewArgError("endpointServiceID", "must be set")
	}
	if cloudProvider == "" {
		return nil, NewArgError("cloudProvider", "must be set")
	}
	if privateEndpointID == "" {
		return nil, NewArgError("privateEndpointID", "must be set")
	}

	basePath := fmt.Sprintf(privateEndpointsPath, groupID)
	path := fmt.Sprintf("%s/%s/endpointService/%s/endpoint/%s", basePath, cloudProvider, endpointServiceID, url.PathEscape(privateEndpointID))

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.Client.Do(ctx, req, nil)
}

// UpdateRegionalizedPrivateEndpointSetting updates the regionalized private endpoint setting for one Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoints-update-regional-mode
func (s *PrivateEndpointsServiceOp) UpdateRegionalizedPrivateEndpointSetting(ctx context.Context, groupID string, enabled bool) (*RegionalizedPrivateEndpointSetting, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(regionalModePath, groupID)
	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, enabled)
	if err != nil {
		return nil, nil, err
	}

	root := new(RegionalizedPrivateEndpointSetting)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetRegionalizedPrivateEndpointSetting updates the regionalized private endpoint setting for one Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoints-get-regional-mode
func (s *PrivateEndpointsServiceOp) GetRegionalizedPrivateEndpointSetting(ctx context.Context, groupID string) (*RegionalizedPrivateEndpointSetting, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(regionalModePath, groupID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(RegionalizedPrivateEndpointSetting)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
