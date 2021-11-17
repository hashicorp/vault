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

const dbUsersBasePath = "api/atlas/v1.0/groups/%s/databaseUsers"

var adminX509Type = map[string]struct{}{
	"MANAGED":  {},
	"CUSTOMER": {},
}

var awsIAMType = map[string]struct{}{
	"USER": {},
	"ROLE": {},
}

// DatabaseUsersService is an interface for interfacing with the Database Users
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/database-users/index.html
type DatabaseUsersService interface {
	List(context.Context, string, *ListOptions) ([]DatabaseUser, *Response, error)
	Get(context.Context, string, string, string) (*DatabaseUser, *Response, error)
	Create(context.Context, string, *DatabaseUser) (*DatabaseUser, *Response, error)
	Update(context.Context, string, string, *DatabaseUser) (*DatabaseUser, *Response, error)
	Delete(context.Context, string, string, string) (*Response, error)
}

// DatabaseUsersServiceOp handles communication with the DatabaseUsers related methods of the
// MongoDB Atlas API.
type DatabaseUsersServiceOp service

var _ DatabaseUsersService = &DatabaseUsersServiceOp{}

// Role allows the user to perform particular actions on the specified database.
// A role on the admin database can include privileges that apply to the other databases as well.
type Role struct {
	RoleName       string `json:"roleName,omitempty"`
	DatabaseName   string `json:"databaseName,omitempty"`
	CollectionName string `json:"collectionName,omitempty"`
}

// DatabaseUser represents MongoDB users in your cluster.
type DatabaseUser struct {
	DatabaseName    string  `json:"databaseName,omitempty"`
	DeleteAfterDate string  `json:"deleteAfterDate,omitempty"`
	Labels          []Label `json:"labels,omitempty"`
	LDAPAuthType    string  `json:"ldapAuthType,omitempty"`
	X509Type        string  `json:"x509Type,omitempty"`
	AWSIAMType      string  `json:"awsIAMType,omitempty"`
	GroupID         string  `json:"groupId,omitempty"`
	Roles           []Role  `json:"roles,omitempty"`
	Scopes          []Scope `json:"scopes"`
	Password        string  `json:"password,omitempty"`
	Username        string  `json:"username,omitempty"`
}

// GetAuthDB determines the authentication database based on the type of user.
// LDAP, X509 and AWSIAM should all use $external.
// SCRAM-SHA should use admin.
func (user *DatabaseUser) GetAuthDB() (name string) {
	// base documentation https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/database_user
	name = "admin"
	_, isX509 := adminX509Type[user.X509Type]
	_, isIAM := awsIAMType[user.AWSIAMType]

	// just USER is external
	isLDAP := len(user.LDAPAuthType) > 0 && user.LDAPAuthType == "USER"

	if isX509 || isIAM || isLDAP {
		name = "$external"
	}

	return
}

// Scope if presents a database user only have access to the indicated resource
// if none is given then it has access to all.
type Scope struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Label containing key-value pairs that tag and categorize the database user.
type Label struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// databaseUserListResponse is the response from the DatabaseUserService.List.
type databaseUsers struct {
	Links      []*Link        `json:"links"`
	Results    []DatabaseUser `json:"results"`
	TotalCount int            `json:"totalCount"`
}

// List gets all users in the project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/database-users-get-all-users/
func (s *DatabaseUsersServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]DatabaseUser, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	path := fmt.Sprintf(dbUsersBasePath, groupID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(databaseUsers)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// Get gets a single user in the project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/database-users-get-single-user/
func (s *DatabaseUsersServiceOp) Get(ctx context.Context, databaseName, groupID, username string) (*DatabaseUser, *Response, error) {
	if databaseName == "" {
		return nil, nil, NewArgError("databaseName", "must be set")
	}
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if username == "" {
		return nil, nil, NewArgError("username", "must be set")
	}

	basePath := fmt.Sprintf(dbUsersBasePath, groupID)
	escapedEntry := url.PathEscape(username)
	path := fmt.Sprintf("%s/%s/%s", basePath, databaseName, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(DatabaseUser)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates a user for the project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/database-users-create-a-user/
func (s *DatabaseUsersServiceOp) Create(ctx context.Context, groupID string, createRequest *DatabaseUser) (*DatabaseUser, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(dbUsersBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(DatabaseUser)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update updates a user for the project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/database-users-update-a-user/
func (s *DatabaseUsersServiceOp) Update(ctx context.Context, groupID, username string, updateRequest *DatabaseUser) (*DatabaseUser, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if username == "" {
		return nil, nil, NewArgError("username", "must be set")
	}
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(dbUsersBasePath, groupID)

	escapedEntry := url.PathEscape(username)

	path := fmt.Sprintf("%s/%s/%s", basePath, updateRequest.GetAuthDB(), escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(DatabaseUser)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes a user for the project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/database-users-delete-a-user/
func (s *DatabaseUsersServiceOp) Delete(ctx context.Context, databaseName, groupID, username string) (*Response, error) {
	if databaseName == "" {
		return nil, NewArgError("databaseName", "must be set")
	}
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}
	if username == "" {
		return nil, NewArgError("username", "must be set")
	}

	basePath := fmt.Sprintf(dbUsersBasePath, groupID)
	escapedEntry := url.PathEscape(username)
	path := fmt.Sprintf("%s/%s/%s", basePath, databaseName, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
