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

const serverlessInstancesPath = "api/atlas/v1.0/groups/%s/serverless"

// ServerlessInstancesService is an interface for interfacing with the Serverless Instances endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/serverless/return-one-serverless-instance/
type ServerlessInstancesService interface {
	List(context.Context, string, *ListOptions) (*ClustersResponse, *Response, error)
	Get(context.Context, string, string) (*Cluster, *Response, error)
	Create(context.Context, string, *ServerlessCreateRequestParams) (*Cluster, *Response, error)
	Update(context.Context, string, string, *ServerlessUpdateRequestParams) (*Cluster, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

// ServerlessInstancesServiceOp handles communication with the Serverless Instances related methods of the MongoDB Atlas API.
type ServerlessInstancesServiceOp service

var _ ServerlessInstancesService = &ServerlessInstancesServiceOp{}

// ClustersResponse represents the response of ServerlessInstancesService.List.
type ClustersResponse struct {
	Links      []*Link    `json:"links,omitempty"`
	Results    []*Cluster `json:"results,omitempty"`
	TotalCount int        `json:"totalCount,omitempty"`
}

// ServerlessCreateRequestParams represents the Request Body Parameters of ServerlessInstancesService.Create.
type ServerlessCreateRequestParams struct {
	Name                         string                      `json:"name,omitempty"`
	ProviderSettings             *ServerlessProviderSettings `json:"providerSettings,omitempty"`
	ServerlessBackupOptions      *ServerlessBackupOptions    `json:"serverlessBackupOptions,omitempty"`
	TerminationProtectionEnabled *bool                       `json:"terminationProtectionEnabled,omitempty"`
	Tag                          *[]*Tag                     `json:"tags,omitempty"`
}

type ServerlessUpdateRequestParams struct {
	ServerlessBackupOptions      *ServerlessBackupOptions `json:"serverlessBackupOptions"`
	TerminationProtectionEnabled *bool                    `json:"terminationProtectionEnabled,omitempty"`
	Tag                          *[]*Tag                  `json:"tags"`
}

// ServerlessProviderSettings represents the Provider Settings of serverless instances.
type ServerlessProviderSettings struct {
	BackingProviderName string `json:"backingProviderName,omitempty"`
	ProviderName        string `json:"providerName,omitempty"`
	RegionName          string `json:"regionName,omitempty"`
}

// ServerlessBackupOptions Serverless Continuous Backup.
type ServerlessBackupOptions struct {
	ServerlessContinuousBackupEnabled *bool `json:"serverlessContinuousBackupEnabled,omitempty"`
}

// List gets all serverless instances in the specified project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/serverless/return-all-serverless-instances/
func (s *ServerlessInstancesServiceOp) List(ctx context.Context, projectID string, listOptions *ListOptions) (*ClustersResponse, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}

	path := fmt.Sprintf(serverlessInstancesPath, projectID)
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ClustersResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Get retrieves one serverless instance in the specified project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/serverless/return-one-serverless-instance/
func (s *ServerlessInstancesServiceOp) Get(ctx context.Context, projectID, instanceName string) (*Cluster, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}

	if instanceName == "" {
		return nil, nil, NewArgError("instanceName", "must be set")
	}

	basePath := fmt.Sprintf(serverlessInstancesPath, projectID)
	path := fmt.Sprintf("%s/%s", basePath, instanceName)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Cluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates one serverless instance in the specified project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/serverless/create-one-serverless-instance/
func (s *ServerlessInstancesServiceOp) Create(ctx context.Context, projectID string, bodyParams *ServerlessCreateRequestParams) (*Cluster, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}

	path := fmt.Sprintf(serverlessInstancesPath, projectID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, bodyParams)
	if err != nil {
		return nil, nil, err
	}

	root := new(Cluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update one serverless instance in the specified project..
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/serverless/update-one-serverless-instance/
func (s *ServerlessInstancesServiceOp) Update(ctx context.Context, projectID, instanceName string, bodyParams *ServerlessUpdateRequestParams) (*Cluster, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if instanceName == "" {
		return nil, nil, NewArgError("instanceName", "must be set")
	}

	basePath := fmt.Sprintf(serverlessInstancesPath, projectID)
	path := fmt.Sprintf("%s/%s", basePath, instanceName)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, bodyParams)
	if err != nil {
		return nil, nil, err
	}

	root := new(Cluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes one serverless instance in the specified project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/serverless/remove-one-serverless-instance/
func (s *ServerlessInstancesServiceOp) Delete(ctx context.Context, projectID, instanceName string) (*Response, error) {
	if projectID == "" {
		return nil, NewArgError("projectID", "must be set")
	}
	if instanceName == "" {
		return nil, NewArgError("instanceName", "must be set")
	}

	basePath := fmt.Sprintf(serverlessInstancesPath, projectID)
	path := fmt.Sprintf("%s/%s", basePath, instanceName)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
