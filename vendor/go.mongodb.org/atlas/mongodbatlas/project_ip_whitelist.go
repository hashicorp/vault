package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const projectIPWhitelistPath = "groups/%s/whitelist"

// ProjectIPWhitelistService is an interface for interfacing with the Project IP Whitelist
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/whitelist/
type ProjectIPWhitelistService interface {
	List(context.Context, string, *ListOptions) ([]ProjectIPWhitelist, *Response, error)
	Get(context.Context, string, string) (*ProjectIPWhitelist, *Response, error)
	Create(context.Context, string, []*ProjectIPWhitelist) ([]ProjectIPWhitelist, *Response, error)
	Update(context.Context, string, []*ProjectIPWhitelist) ([]ProjectIPWhitelist, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

// ProjectIPWhitelistServiceOp handles communication with the ProjectIPWhitelist related methods
// of the MongoDB Atlas API
type ProjectIPWhitelistServiceOp service

var _ ProjectIPWhitelistService = &ProjectIPWhitelistServiceOp{}

// ProjectIPWhitelist represents MongoDB project's IP whitelist.
type ProjectIPWhitelist struct {
	GroupID          string `json:"groupId,omitempty"`          // The unique identifier for the project for which you want to update one or more whitelist entries.
	AwsSecurityGroup string `json:"awsSecurityGroup,omitempty"` // ID of the whitelisted AWS security group to update. Mutually exclusive with cidrBlock and ipAddress.
	CIDRBlock        string `json:"cidrBlock,omitempty"`        // Whitelist entry in Classless Inter-Domain Routing (CIDR) notation to update. Mutually exclusive with awsSecurityGroup and ipAddress.
	IPAddress        string `json:"ipAddress,omitempty"`        // Whitelisted IP address to update. Mutually exclusive with awsSecurityGroup and cidrBlock.
	Comment          string `json:"comment,omitempty"`          // Optional The comment associated with the whitelist entry. Specify an empty string "" to delete the comment associated to an IP address.
	DeleteAfterDate  string `json:"deleteAfterDate,omitempty"`  // Optional The ISO-8601-formatted UTC date after which Atlas removes the entry from the whitelist. The specified date must be in the future and within one week of the time you make the API request. To update a temporary whitelist entry to be permanent, set the value of this field to null
}

// projectIPWhitelistsResponse is the response from the ProjectIPWhitelistService.List.
type projectIPWhitelistsResponse struct {
	Links      []*Link              `json:"links"`
	Results    []ProjectIPWhitelist `json:"results"`
	TotalCount int                  `json:"totalCount"`
}

// Create adds one or more whitelist entries to the project associated to {GROUP-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/database-users-create-a-user/
func (s *ProjectIPWhitelistServiceOp) Create(ctx context.Context, groupID string, createRequest []*ProjectIPWhitelist) ([]ProjectIPWhitelist, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(projectIPWhitelistPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(projectIPWhitelistsResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, err
}

// Get gets the whitelist entry specified to {WHITELIST-ENTRY} from the project associated to {GROUP-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/whitelist-get-one-entry/
func (s *ProjectIPWhitelistServiceOp) Get(ctx context.Context, groupID, whiteListEntry string) (*ProjectIPWhitelist, *Response, error) {
	if whiteListEntry == "" {
		return nil, nil, NewArgError("whiteListEntry", "must be set")
	}

	basePath := fmt.Sprintf(projectIPWhitelistPath, groupID)
	escapedEntry := url.PathEscape(whiteListEntry)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProjectIPWhitelist)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// List all whitelist entries in the project associated to {GROUP-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/whitelist-get-all/
func (s *ProjectIPWhitelistServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]ProjectIPWhitelist, *Response, error) {
	path := fmt.Sprintf(projectIPWhitelistPath, groupID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(projectIPWhitelistsResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// Update one or more whitelist entries in the project associated to {GROUP-ID}
// See more: https://docs.atlas.mongodb.com/reference/api/whitelist-update-one/
func (s *ProjectIPWhitelistServiceOp) Update(ctx context.Context, groupID string, updateRequest []*ProjectIPWhitelist) ([]ProjectIPWhitelist, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	path := fmt.Sprintf(projectIPWhitelistPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(projectIPWhitelistsResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, err
}

// Delete the whitelist entry specified to {WHITELIST-ENTRY} from the project associated to {GROUP-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/whitelist-delete-one/
func (s *ProjectIPWhitelistServiceOp) Delete(ctx context.Context, groupID, whitelistEntry string) (*Response, error) {
	if whitelistEntry == "" {
		return nil, NewArgError("whitelistEntry", "must be set")
	}

	basePath := fmt.Sprintf(projectIPWhitelistPath, groupID)
	escapedEntry := url.PathEscape(whitelistEntry)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
