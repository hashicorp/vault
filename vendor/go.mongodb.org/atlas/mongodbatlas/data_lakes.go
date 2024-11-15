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

const (
	dataLakesBasePath                = "api/atlas/v1.0/groups"
	privateLinkEndpointsDataLakePath = "api/atlas/v1.0/groups/%s/privateNetworkSettings/endpointIds"
)

// DataLakeService is an interface for interfacing with the Data Lake endpoints of the MongoDB Atlas API.
//
// See more: https://docs.mongodb.com/datalake/reference/api/datalakes-api
type DataLakeService interface {
	List(context.Context, string) ([]DataLake, *Response, error)
	Get(context.Context, string, string) (*DataLake, *Response, error)
	Create(context.Context, string, *DataLakeCreateRequest) (*DataLake, *Response, error)
	Update(context.Context, string, string, *DataLakeUpdateRequest) (*DataLake, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
	CreatePrivateLinkEndpoint(context.Context, string, *PrivateLinkEndpointDataLake) (*PrivateLinkEndpointDataLakeResponse, *Response, error)
	GetPrivateLinkEndpoint(context.Context, string, string) (*PrivateLinkEndpointDataLake, *Response, error)
	ListPrivateLinkEndpoint(context.Context, string) (*PrivateLinkEndpointDataLakeResponse, *Response, error)
	DeletePrivateLinkEndpoint(context.Context, string, string) (*Response, error)
}

// DataLakeServiceOp handles communication with the DataLakeService related methods of the
// MongoDB Atlas API.
type DataLakeServiceOp service

var _ DataLakeService = &DataLakeServiceOp{}

// AwsCloudProviderConfig is the data lake configuration for AWS.
type AwsCloudProviderConfig struct {
	ExternalID        string `json:"externalId,omitempty"`
	IAMAssumedRoleARN string `json:"iamAssumedRoleARN,omitempty"`
	IAMUserARN        string `json:"iamUserARN,omitempty"`
	RoleID            string `json:"roleId,omitempty"`
	TestS3Bucket      string `json:"testS3Bucket,omitempty"`
}

// CloudProviderConfig represents the configuration for all supported cloud providers.
type CloudProviderConfig struct {
	AWSConfig AwsCloudProviderConfig `json:"aws,omitempty"`
}

// DataProcessRegion represents the region where a data lake is processed.
type DataProcessRegion struct {
	CloudProvider string `json:"cloudProvider,omitempty"`
	Region        string `json:"region,omitempty"`
}

// DataLakeStore represents a store of data lake data. Docs: https://docs.mongodb.com/datalake/reference/format/data-lake-configuration/#stores
type DataLakeStore struct {
	Name                     string   `json:"name,omitempty"`
	Provider                 string   `json:"provider,omitempty"`
	Region                   string   `json:"region,omitempty"`
	Bucket                   string   `json:"bucket,omitempty"`
	Prefix                   string   `json:"prefix,omitempty"`
	Delimiter                string   `json:"delimiter,omitempty"`
	IncludeTags              *bool    `json:"includeTags,omitempty"`
	AdditionalStorageClasses []string `json:"additionalStorageClasses,omitempty"`
}

// DataLakeDataSource represents the data source of a data lake.
type DataLakeDataSource struct {
	StoreName     string `json:"storeName,omitempty"`
	DefaultFormat string `json:"defaultFormat,omitempty"`
	Path          string `json:"path,omitempty"`
}

// DataLakeCollection represents collections under a DataLakeDatabase.
type DataLakeCollection struct {
	Name        string               `json:"name,omitempty"`
	DataSources []DataLakeDataSource `json:"dataSources,omitempty"`
}

// DataLakeDatabaseView represents any view under a DataLakeDatabase.
type DataLakeDatabaseView struct {
	Name     string `json:"name,omitempty"`
	Source   string `json:"source,omitempty"`
	Pipeline string `json:"pipeline,omitempty"`
}

// DataLakeDatabase represents the mapping of a data lake to a database. Docs: https://docs.mongodb.com/datalake/reference/format/data-lake-configuration/#databases
type DataLakeDatabase struct {
	Name                   string                 `json:"name,omitempty"`
	Collections            []DataLakeCollection   `json:"collections,omitempty"`
	Views                  []DataLakeDatabaseView `json:"views,omitempty"`
	MaxWildcardCollections *int64                 `json:"maxWildcardCollections,omitempty"`
}

// Storage represents the storage configuration for a data lake.
type Storage struct {
	Databases []DataLakeDatabase `json:"databases,omitempty"`
	Stores    []DataLakeStore    `json:"stores,omitempty"`
}

// DataLake represents a data lake.
type DataLake struct {
	CloudProviderConfig CloudProviderConfig `json:"cloudProviderConfig,omitempty"` // Configuration for the cloud service where Data Lake source data is stored.
	DataProcessRegion   DataProcessRegion   `json:"dataProcessRegion,omitempty"`   // Cloud provider region which clients are routed to for data processing.
	GroupID             string              `json:"groupId,omitempty"`             // Unique identifier for the project.
	Hostnames           []string            `json:"hostnames,omitempty"`           // List of hostnames for the data lake.
	Name                string              `json:"name,omitempty"`                // Name of the data lake.
	State               string              `json:"state,omitempty"`               // Current state of the data lake.
	Storage             Storage             `json:"storage,omitempty"`             // Configuration for each data store and its mapping to MongoDB collections / databases.
}

// DataLakeUpdateRequest represents all possible fields that can be updated in a data lake.
type DataLakeUpdateRequest struct {
	CloudProviderConfig *CloudProviderConfig `json:"cloudProviderConfig,omitempty"`
	DataProcessRegion   *DataProcessRegion   `json:"dataProcessRegion,omitempty"`
}

// DataLakeCreateRequest represents the required fields to create a new data lake.
type DataLakeCreateRequest struct {
	Name                string               `json:"name,omitempty"`
	CloudProviderConfig *CloudProviderConfig `json:"cloudProviderConfig,omitempty"`
}

// PrivateLinkEndpointDataLakeResponse represents MongoDB Private Endpoint Connection to DataLake.
type PrivateLinkEndpointDataLakeResponse struct {
	Links      []*Link                        `json:"links,omitempty"`
	Results    []*PrivateLinkEndpointDataLake `json:"results"`
	TotalCount int                            `json:"totalCount"`
}

// PrivateLinkEndpointDataLake represents the private link result for data lake.
type PrivateLinkEndpointDataLake struct {
	Comment    string `json:"comment,omitempty"`
	EndpointID string `json:"endpointId,omitempty"`
	Provider   string `json:"provider,omitempty"`
	Type       string `json:"type,omitempty"`
}

// List gets all data lakes for the specified group.
//
// See more: https://docs.mongodb.com/datalake/reference/api/dataLakes-get-all-tenants
func (s *DataLakeServiceOp) List(ctx context.Context, groupID string) ([]DataLake, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf("%s/%s/dataLakes", dataLakesBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []DataLake
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Get gets the data laked associated with a specific name.
//
// See more: https://docs.mongodb.com/datalake/reference/api/dataLakes-get-one-tenant/
func (s *DataLakeServiceOp) Get(ctx context.Context, groupID, name string) (*DataLake, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}

	path := fmt.Sprintf("%s/%s/dataLakes/%s", dataLakesBasePath, groupID, name)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataLake)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates a new Data Lake.
//
// See more: https://docs.mongodb.com/datalake/reference/api/dataLakes-create-one-tenant/
func (s *DataLakeServiceOp) Create(ctx context.Context, groupID string, createRequest *DataLakeCreateRequest) (*DataLake, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "must be set")
	}

	path := fmt.Sprintf("%s/%s/dataLakes", dataLakesBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataLake)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update updates an existing Data Lake.
//
// See more: https://docs.mongodb.com/datalake/reference/api/dataLakes-update-one-tenant/
func (s *DataLakeServiceOp) Update(ctx context.Context, groupID, name string, updateRequest *DataLakeUpdateRequest) (*DataLake, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/dataLakes/%s", dataLakesBasePath, groupID, name)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataLake)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes the Data Lake with a given name.
//
// See more: https://docs.mongodb.com/datalake/reference/api/dataLakes-delete-one-tenant/
func (s *DataLakeServiceOp) Delete(ctx context.Context, groupID, name string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if name == "" {
		return nil, NewArgError("name", "must be set")
	}

	path := fmt.Sprintf("%s/%s/dataLakes/%s", dataLakesBasePath, groupID, name)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// CreatePrivateLinkEndpoint creates one private link endpoint in Data Lake Atlas project.
//
// See more: https://docs.mongodb.com/datalake/reference/api/dataLakes-private-link-create-one/#std-label-api-pvt-link-create-one
func (s *DataLakeServiceOp) CreatePrivateLinkEndpoint(ctx context.Context, groupID string, createRequest *PrivateLinkEndpointDataLake) (*PrivateLinkEndpointDataLakeResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "must be set")
	}

	path := fmt.Sprintf(privateLinkEndpointsDataLakePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(PrivateLinkEndpointDataLakeResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// DeletePrivateLinkEndpoint deletes the Data Lake private link endpoint with a given endpoint id.
//
// See more: https://docs.mongodb.com/datalake/reference/api/dataLakes-private-link-delete-one/#std-label-api-pvt-link-delete-one
func (s *DataLakeServiceOp) DeletePrivateLinkEndpoint(ctx context.Context, groupID, endpointID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if endpointID == "" {
		return nil, NewArgError("endpointID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", fmt.Sprintf(privateLinkEndpointsDataLakePath, groupID), endpointID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// ListPrivateLinkEndpoint gets all private link endpoints for data lake for the specified group.
//
// See more: https://docs.atlas.mongodb.com/reference/api/online-archive-private-link-get-all/#std-label-api-online-archive-pvt-link-get-all
func (s *DataLakeServiceOp) ListPrivateLinkEndpoint(ctx context.Context, groupID string) (*PrivateLinkEndpointDataLakeResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(privateLinkEndpointsDataLakePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root = new(PrivateLinkEndpointDataLakeResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return root, resp, err
	}

	return root, resp, nil
}

// GetPrivateLinkEndpoint gets the data lake private link endpoint associated with a specific group and endpointID.
//
// See more: https://docs.mongodb.com/datalake/reference/api/dataLakes-private-link-get-one/#std-label-api-pvt-link-get-one
func (s *DataLakeServiceOp) GetPrivateLinkEndpoint(ctx context.Context, groupID, endpointID string) (*PrivateLinkEndpointDataLake, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if endpointID == "" {
		return nil, nil, NewArgError("endpointID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", fmt.Sprintf(privateLinkEndpointsDataLakePath, groupID), endpointID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(PrivateLinkEndpointDataLake)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
