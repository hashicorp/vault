// Copyright 2019 MongoDB Inc
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
	orgsBasePath = "orgs"
)

// OrganizationsService provides access to the organization related functions in the Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organizations/
type OrganizationsService interface {
	List(context.Context, *ListOptions) (*Organizations, *Response, error)
	Get(context.Context, string) (*Organization, *Response, error)
	Projects(context.Context, string, *ListOptions) (*Projects, *Response, error)
	Users(context.Context, string, *ListOptions) (*AtlasUsersResponse, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

// OrganizationsServiceOp provides an implementation of the OrganizationsService interface
type OrganizationsServiceOp service

var _ OrganizationsService = &OrganizationsServiceOp{}

// Organization represents the structure of an organization.
type Organization struct {
	ID        string  `json:"id,omitempty"`
	IsDeleted *bool   `json:"isDeleted,omitempty"`
	Links     []*Link `json:"links,omitempty"`
	Name      string  `json:"name,omitempty"`
}

// Organizations represents an array of organization
type Organizations struct {
	Links      []*Link         `json:"links"`
	Results    []*Organization `json:"results"`
	TotalCount int             `json:"totalCount"`
}

// List gets all organizations.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organization-get-all/
func (s *OrganizationsServiceOp) List(ctx context.Context, opts *ListOptions) (*Organizations, *Response, error) {
	path, err := setListOptions(orgsBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Organizations)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Get gets a single organization.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organization-get-one/
func (s *OrganizationsServiceOp) Get(ctx context.Context, orgID string) (*Organization, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", orgsBasePath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Organization)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Projects gets all projects for the given organization ID.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organization-get-all-projects/
func (s *OrganizationsServiceOp) Projects(ctx context.Context, orgID string, opts *ListOptions) (*Projects, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	basePath := fmt.Sprintf("%s/%s/groups", orgsBasePath, orgID)

	path, err := setListOptions(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Projects)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Users gets all users for the given organization ID.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organization-users-get-all-users/
func (s *OrganizationsServiceOp) Users(ctx context.Context, orgID string, opts *ListOptions) (*AtlasUsersResponse, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	basePath := fmt.Sprintf("%s/%s/users", orgsBasePath, orgID)

	path, err := setListOptions(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AtlasUsersResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes an organization.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organization-delete-one/
func (s *OrganizationsServiceOp) Delete(ctx context.Context, orgID string) (*Response, error) {
	if orgID == "" {
		return nil, NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf("%s/%s", orgsBasePath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, basePath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
