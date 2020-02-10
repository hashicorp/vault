package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const projectIPWhitelistPath = "groups/%s/whitelist"

//ProjectIPWhitelistService is an interface for interfacing with the Project IP Whitelist
// endpoints of the MongoDB Atlas API.
//See more: https://docs.atlas.mongodb.com/reference/api/whitelist/
type ProjectIPWhitelistService interface {
	List(context.Context, string, *ListOptions) ([]ProjectIPWhitelist, *Response, error)
	Get(context.Context, string, string) (*ProjectIPWhitelist, *Response, error)
	Create(context.Context, string, []*ProjectIPWhitelist) ([]ProjectIPWhitelist, *Response, error)
	Update(context.Context, string, string, []*ProjectIPWhitelist) ([]ProjectIPWhitelist, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

//ProjectIPWhitelistServiceOp handles communication with the ProjectIPWhitelist related methods
// of the MongoDB Atlas API
type ProjectIPWhitelistServiceOp struct {
	client *Client
}

var _ ProjectIPWhitelistService = &ProjectIPWhitelistServiceOp{}

// ProjectIPWhitelist represents MongoDB project's IP whitelist.
type ProjectIPWhitelist struct {
	Comment         string `json:"comment,omitempty"`
	GroupID         string `json:"groupId,omitempty"`
	CIDRBlock       string `json:"cidrBlock,omitempty"`
	IPAddress       string `json:"ipAddress,omitempty"`
	DeleteAfterDate string `json:"deleteAfterDate,omitempty"`
}

// projectIPWhitelistsResponse is the response from the ProjectIPWhitelistService.List.
type projectIPWhitelistsResponse struct {
	Links      []*Link              `json:"links"`
	Results    []ProjectIPWhitelist `json:"results"`
	TotalCount int                  `json:"totalCount"`
}

//List all whitelist entries in the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/whitelist-get-all/
func (s *ProjectIPWhitelistServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]ProjectIPWhitelist, *Response, error) {
	path := fmt.Sprintf(projectIPWhitelistPath, groupID)

	//Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(projectIPWhitelistsResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

//Get gets the whitelist entry specified to {WHITELIST-ENTRY} from the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/whitelist-get-one-entry/
func (s *ProjectIPWhitelistServiceOp) Get(ctx context.Context, groupID string, whiteListEntry string) (*ProjectIPWhitelist, *Response, error) {
	if whiteListEntry == "" {
		return nil, nil, NewArgError("whiteListEntry", "must be set")
	}

	basePath := fmt.Sprintf(projectIPWhitelistPath, groupID)
	escapedEntry := url.PathEscape(whiteListEntry)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProjectIPWhitelist)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Add one or more whitelist entries to the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/database-users-create-a-user/
func (s *ProjectIPWhitelistServiceOp) Create(ctx context.Context, groupID string, createRequest []*ProjectIPWhitelist) ([]ProjectIPWhitelist, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(projectIPWhitelistPath, groupID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(projectIPWhitelistsResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, err
}

//Update one or more whitelist entries in the project associated to {GROUP-ID}
//See more: https://docs.atlas.mongodb.com/reference/api/whitelist-update-one/
func (s *ProjectIPWhitelistServiceOp) Update(ctx context.Context, groupID string, whitelistEntry string, updateRequest []*ProjectIPWhitelist) ([]ProjectIPWhitelist, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(projectIPWhitelistPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, whitelistEntry)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(projectIPWhitelistsResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, err
}

//Delete the whitelist entry specified to {WHITELIST-ENTRY} from the project associated to {GROUP-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/whitelist-delete-one/
func (s *ProjectIPWhitelistServiceOp) Delete(ctx context.Context, groupID string, whitelistEntry string) (*Response, error) {
	if whitelistEntry == "" {
		return nil, NewArgError("whitelistEntry", "must be set")
	}

	basePath := fmt.Sprintf(projectIPWhitelistPath, groupID)
	escapedEntry := url.PathEscape(whitelistEntry)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}
