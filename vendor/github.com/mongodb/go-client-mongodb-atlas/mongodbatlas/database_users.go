package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const dbUsersBasePath = "groups/%s/databaseUsers"

//DatabaseUsersService is an interface for interfacing with the Database Users
// endpoints of the MongoDB Atlas API.
//See more: https://docs.atlas.mongodb.com/reference/api/database-users/index.html
type DatabaseUsersService interface {
	List(context.Context, string, *ListOptions) ([]DatabaseUser, *Response, error)
	Get(context.Context, string, string) (*DatabaseUser, *Response, error)
	Create(context.Context, string, *DatabaseUser) (*DatabaseUser, *Response, error)
	Update(context.Context, string, string, *DatabaseUser) (*DatabaseUser, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

//DatabaseUsersServiceOp handles communication with the DatabaseUsers related methos of the
//MongoDB Atlas API
type DatabaseUsersServiceOp struct {
	client *Client
}

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
	Roles           []Role `json:"roles,omitempty"`
	GroupID         string `json:"groupId,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	DatabaseName    string `json:"databaseName,omitempty"`
	LDAPAuthType    string `json:"ldapAuthType,omitempty"`
	DeleteAfterDate string `json:"deleteAfterDate,omitempty"`
}

// databaseUserListResponse is the response from the DatabaseUserService.List.
type databaseUsers struct {
	Links      []*Link        `json:"links"`
	Results    []DatabaseUser `json:"results"`
	TotalCount int            `json:"totalCount"`
}

//List gets all users in the project.
//See more: https://docs.atlas.mongodb.com/reference/api/database-users-get-all-users/
func (s *DatabaseUsersServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]DatabaseUser, *Response, error) {
	path := fmt.Sprintf(dbUsersBasePath, groupID)

	//Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(databaseUsers)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

//Get gets a single user in the project.
//See more: https://docs.atlas.mongodb.com/reference/api/database-users-get-single-user/
func (s *DatabaseUsersServiceOp) Get(ctx context.Context, groupID string, username string) (*DatabaseUser, *Response, error) {
	if username == "" {
		return nil, nil, NewArgError("username", "must be set")
	}

	basePath := fmt.Sprintf(dbUsersBasePath, groupID)
	path := fmt.Sprintf("%s/admin/%s", basePath, username)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(DatabaseUser)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Create creates a user for the project.
//See more: https://docs.atlas.mongodb.com/reference/api/database-users-create-a-user/
func (s *DatabaseUsersServiceOp) Create(ctx context.Context, groupID string, createRequest *DatabaseUser) (*DatabaseUser, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(dbUsersBasePath, groupID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(DatabaseUser)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Update updates a user for the project.
//See more: https://docs.atlas.mongodb.com/reference/api/database-users-update-a-user/
func (s *DatabaseUsersServiceOp) Update(ctx context.Context, groupID string, username string, updateRequest *DatabaseUser) (*DatabaseUser, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(dbUsersBasePath, groupID)
	path := fmt.Sprintf("%s/admin/%s", basePath, username)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(DatabaseUser)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Delete deletes a user for the project.
// See more: https://docs.atlas.mongodb.com/reference/api/database-users-delete-a-user/
func (s *DatabaseUsersServiceOp) Delete(ctx context.Context, groupID string, username string) (*Response, error) {
	if username == "" {
		return nil, NewArgError("username", "must be set")
	}

	basePath := fmt.Sprintf(dbUsersBasePath, groupID)
	path := fmt.Sprintf("%s/admin/%s", basePath, username)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}
