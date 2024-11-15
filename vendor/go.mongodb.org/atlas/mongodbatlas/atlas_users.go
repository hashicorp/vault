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
	usersGroupBasePath = "api/atlas/v1.0/groups/%s/users"
	usersBasePath      = "api/atlas/v1.0/users"
)

// AtlasUsersService is an interface for interfacing with the AtlasUsers
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/users/
type AtlasUsersService interface {
	List(context.Context, string, *ListOptions) ([]AtlasUser, *Response, error)
	Get(context.Context, string) (*AtlasUser, *Response, error)
	GetByName(context.Context, string) (*AtlasUser, *Response, error)
	Create(context.Context, *AtlasUser) (*AtlasUser, *Response, error)
}

// AtlasUsersServiceOp handles communication with the AtlasUsers related methods of the
// MongoDB Atlas API.
type AtlasUsersServiceOp service

var _ AtlasUsersService = &AtlasUsersServiceOp{}

// AtlasUsersResponse represents a array of users.
type AtlasUsersResponse struct {
	Links      []*Link     `json:"links"`
	Results    []AtlasUser `json:"results"`
	TotalCount int         `json:"totalCount"`
}

// AtlasUser represents a user.
type AtlasUser struct {
	EmailAddress string      `json:"emailAddress"`
	FirstName    string      `json:"firstName"`
	ID           string      `json:"id,omitempty"`
	LastName     string      `json:"lastName"`
	Roles        []AtlasRole `json:"roles"`
	TeamIds      []string    `json:"teamIds,omitempty"`
	Username     string      `json:"username"`
	MobileNumber string      `json:"mobileNumber"`
	Password     string      `json:"password"`
	Country      string      `json:"country"`
}

// List gets all users.
//
// See more: https://docs.atlas.mongodb.com/reference/api/user-get-all/
func (s *AtlasUsersServiceOp) List(ctx context.Context, orgID string, listOptions *ListOptions) ([]AtlasUser, *Response, error) {
	path := fmt.Sprintf(usersGroupBasePath, orgID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
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

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// Get gets a single atlas user.
//
// See more: https://docs.atlas.mongodb.com/reference/api/user-get-by-id/
func (s *AtlasUsersServiceOp) Get(ctx context.Context, userID string) (*AtlasUser, *Response, error) {
	if userID == "" {
		return nil, nil, NewArgError("userID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", usersBasePath, userID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AtlasUser)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetByName gets a single atlas user by name.
//
// See more: https://docs.atlas.mongodb.com/reference/api/user-get-one-by-name/
func (s *AtlasUsersServiceOp) GetByName(ctx context.Context, username string) (*AtlasUser, *Response, error) {
	if username == "" {
		return nil, nil, NewArgError("username", "must be set")
	}

	path := fmt.Sprintf("%s/byName/%s", usersBasePath, username)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AtlasUser)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates an Atlas User.
//
// See more: https://docs.atlas.mongodb.com/reference/api/user-create/
func (s *AtlasUsersServiceOp) Create(ctx context.Context, createRequest *AtlasUser) (*AtlasUser, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPost, usersBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(AtlasUser)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
