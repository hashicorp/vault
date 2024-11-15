// Copyright 2023 MongoDB Inc
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
	dataFederationBasePath           = "api/atlas/v1.0/groups/%s/dataFederation"
	dataFederationQueryLimitBasePath = "api/atlas/v1.0/groups/%s/dataFederation/%s/limits"
)

// DataFederationService is an interface for interfacing with the Data Federation endpoints of the MongoDB Atlas API.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation
type DataFederationService interface {
	List(context.Context, string) ([]*DataFederationInstance, *Response, error)
	Get(context.Context, string, string) (*DataFederationInstance, *Response, error)
	Create(context.Context, string, *DataFederationInstance) (*DataFederationInstance, *Response, error)
	Update(context.Context, string, string, *DataFederationInstance, *DataFederationUpdateOptions) (*DataFederationInstance, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
	ListQueryLimits(context.Context, string, string) ([]*DataFederationQueryLimit, *Response, error)
	GetQueryLimit(context.Context, string, string, string) (*DataFederationQueryLimit, *Response, error)
	ConfigureQueryLimit(context.Context, string, string, string, *DataFederationQueryLimit) (*DataFederationQueryLimit, *Response, error)
	DeleteQueryLimit(context.Context, string, string, string) (*Response, error)
}

// DataFederationServiceOp handles communication with the DataFederationService related methods of the
// MongoDB Atlas API.
type DataFederationServiceOp service

var _ DataFederationService = &DataFederationServiceOp{}

// DataFederationInstance is the data federation configuration.
type DataFederationInstance struct {
	CloudProviderConfig *CloudProviderConfig   `json:"cloudProviderConfig,omitempty"`
	DataProcessRegion   *DataProcessRegion     `json:"dataProcessRegion,omitempty"`
	Storage             *DataFederationStorage `json:"storage,omitempty"`
	Name                string                 `json:"name,omitempty"`
	State               string                 `json:"state,omitempty"`
	Hostnames           []string               `json:"hostnames,omitempty"`
}

// DataFederationStorage represents the storage configuration for a data lake.
type DataFederationStorage struct {
	Databases []*DataFederationDatabase `json:"databases,omitempty"`
	Stores    []*DataFederationStore    `json:"stores,omitempty"`
}

// DataFederationDatabase represents queryable databases and collections for this data federation.
type DataFederationDatabase struct {
	Collections            []*DataFederationCollection   `json:"collections,omitempty"`
	Views                  []*DataFederationDatabaseView `json:"views,omitempty"`
	MaxWildcardCollections int32                         `json:"maxWildcardCollections,omitempty"`
	Name                   string                        `json:"name,omitempty"`
}

// DataFederationCollection represents queryable collections for this data federation.
type DataFederationCollection struct {
	DataSources []*DataFederationDataSource `json:"dataSources,omitempty"`
	Name        string                      `json:"name,omitempty"`
}

// DataFederationDataSource represents data stores that map to a collection for this data federation.
type DataFederationDataSource struct {
	AllowInsecure       *bool     `json:"allowInsecure,omitempty"`
	Collection          string    `json:"collection,omitempty"`
	CollectionRegex     string    `json:"collectionRegex,omitempty"`
	Database            string    `json:"database,omitempty"`
	DatabaseRegex       string    `json:"databaseRegex,omitempty"`
	DefaultFormat       string    `json:"defaultFormat,omitempty"`
	Path                string    `json:"path,omitempty"`
	ProvenanceFieldName string    `json:"provenanceFieldName,omitempty"`
	StoreName           string    `json:"storeName,omitempty"`
	Urls                []*string `json:"urls,omitempty"`
}

// DataFederationDatabaseView represents any view under a DataFederationDatabase.
type DataFederationDatabaseView struct {
	Name     string `json:"name,omitempty"`
	Source   string `json:"source,omitempty"`
	Pipeline string `json:"pipeline,omitempty"`
}

// DataFederationStore represents data stores for the data federation.
type DataFederationStore struct {
	ReadPreference           *ReadPreference `json:"readPreference,omitempty"`
	AdditionalStorageClasses []*string       `json:"additionalStorageClasses,omitempty"`
	Urls                     []*string       `json:"urls,omitempty"`
	Name                     string          `json:"name,omitempty"`
	Provider                 string          `json:"provider,omitempty"`
	ClusterName              string          `json:"clusterName,omitempty"`
	ClusterID                string          `json:"clusterId,omitempty"`
	Region                   string          `json:"region,omitempty"`
	Bucket                   string          `json:"bucket,omitempty"`
	Prefix                   string          `json:"prefix,omitempty"`
	Delimiter                string          `json:"delimiter,omitempty"`
	ProjectID                string          `json:"projectId,omitempty"`
	DefaultFormat            string          `json:"defaultFormat,omitempty"`
	IncludeTags              *bool           `json:"includeTags,omitempty"`
	Public                   *bool           `json:"public,omitempty"`
	AllowInsecure            *bool           `json:"allowInsecure,omitempty"`
}

// ReadPreference describes how to route read requests to the cluster.
type ReadPreference struct {
	MaxStalenessSeconds int32     `json:"maxStalenessSeconds,omitempty"`
	Mode                string    `json:"mode,omitempty"`
	TagSets             []*TagSet `json:"tagSets,omitempty"`
}

// TagSet describes a tag specification document.
type TagSet struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// DataFederationUpdateOptions specifies the optional parameters to Update method.
type DataFederationUpdateOptions struct {
	// Flag that indicates whether this request should check if the requesting IAM role can read from the S3 bucket.
	// AWS checks if the role can list the objects in the bucket before writing to it.
	// Some IAM roles only need write permissions. This flag allows you to skip that check.
	SkipRoleValidation bool `url:"skipRoleValidation"`
}

// DataFederationQueryLimit Details of a tenant-level query limit for Data Federation.
type DataFederationQueryLimit struct {
	CurrentUsage     int64  `json:"currentUsage,omitempty"`
	DefaultLimit     int64  `json:"defaultLimit,omitempty"`
	LastModifiedDate string `json:"lastModifiedDate,omitempty"`
	MaximumLimit     int64  `json:"maximumLimit,omitempty"`
	Name             string `json:"name,omitempty"`
	OverrunPolicy    string `json:"overrunPolicy,omitempty"`
	TenantName       string `json:"tenantName,omitempty"`
	Value            int64  `json:"value,omitempty"`
}

// List gets the details of all federated database instances in the specified project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/listFederatedDatabases
func (s *DataFederationServiceOp) List(ctx context.Context, groupID string) ([]*DataFederationInstance, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	path := fmt.Sprintf(dataFederationBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []*DataFederationInstance
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Get gets the details of one federated database instance within the specified project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/getFederatedDatabase
func (s *DataFederationServiceOp) Get(ctx context.Context, groupID, name string) (*DataFederationInstance, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}

	basePath := fmt.Sprintf(dataFederationBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, name)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataFederationInstance)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates one federated database instance in the specified project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createFederatedDatabase
func (s *DataFederationServiceOp) Create(ctx context.Context, groupID string, createRequest *DataFederationInstance) (*DataFederationInstance, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "must be set")
	}

	path := fmt.Sprintf(dataFederationBasePath, groupID)
	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataFederationInstance)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update updates the details of one federated database instance in the specified project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/updateFederatedDatabase
func (s *DataFederationServiceOp) Update(ctx context.Context, groupID, name string, updateRequest *DataFederationInstance, option *DataFederationUpdateOptions) (*DataFederationInstance, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(dataFederationBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, name)

	if option == nil {
		option = &DataFederationUpdateOptions{
			SkipRoleValidation: true,
		}
	}

	// Add query params from DataFederationUpdateOptions
	pathWithOptions, err := setListOptions(path, option)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, pathWithOptions, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataFederationInstance)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete removes one federated database instance from the specified project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/deleteFederatedDatabase
func (s *DataFederationServiceOp) Delete(ctx context.Context, groupID, name string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if name == "" {
		return nil, NewArgError("name", "must be set")
	}

	basePath := fmt.Sprintf(dataFederationBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, name)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// ConfigureQueryLimit Creates or updates one query limit for one federated database instance.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createOneDataFederationQueryLimit
func (s *DataFederationServiceOp) ConfigureQueryLimit(ctx context.Context, groupID, name, limitName string, queryLimit *DataFederationQueryLimit) (*DataFederationQueryLimit, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}
	if limitName == "" {
		return nil, nil, NewArgError("limitName", "must be set")
	}
	if queryLimit == nil {
		return nil, nil, NewArgError("queryLimit", "must be set")
	}

	basePath := fmt.Sprintf(dataFederationQueryLimitBasePath, groupID, name)
	path := fmt.Sprintf("%s/%s", basePath, limitName)
	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, queryLimit)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataFederationQueryLimit)
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

func (s *DataFederationServiceOp) DeleteQueryLimit(ctx context.Context, groupID, name, limitName string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if name == "" {
		return nil, NewArgError("name", "must be set")
	}
	if limitName == "" {
		return nil, NewArgError("limitName", "must be set")
	}

	basePath := fmt.Sprintf(dataFederationQueryLimitBasePath, groupID, name)
	path := fmt.Sprintf("%s/%s", basePath, limitName)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

func (s *DataFederationServiceOp) GetQueryLimit(ctx context.Context, groupID, name, limitName string) (*DataFederationQueryLimit, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}
	if limitName == "" {
		return nil, nil, NewArgError("limitName", "must be set")
	}

	basePath := fmt.Sprintf(dataFederationQueryLimitBasePath, groupID, name)
	path := fmt.Sprintf("%s/%s", basePath, limitName)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataFederationQueryLimit)
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

func (s *DataFederationServiceOp) ListQueryLimits(ctx context.Context, groupID, name string) ([]*DataFederationQueryLimit, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}

	path := fmt.Sprintf(dataFederationQueryLimitBasePath, groupID, name)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []*DataFederationQueryLimit
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
