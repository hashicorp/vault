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

const projectIPAccessListPath = "api/atlas/v1.0/groups/%s/accessList"

// ProjectIPAccessListService provides access to the project access list related functions in the Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/organizations/
type ProjectIPAccessListService interface {
	List(context.Context, string, *ListOptions) (*ProjectIPAccessLists, *Response, error)
	Get(context.Context, string, string) (*ProjectIPAccessList, *Response, error)
	Create(context.Context, string, []*ProjectIPAccessList) (*ProjectIPAccessLists, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

// ProjectIPAccessListServiceOp provides an implementation of the ProjectIPAccessListService interface.
type ProjectIPAccessListServiceOp service

var _ ProjectIPAccessListService = &ProjectIPAccessListServiceOp{}

// ProjectIPAccessList represents MongoDB project's IP access list.
type ProjectIPAccessList struct {
	AwsSecurityGroup string `json:"awsSecurityGroup,omitempty"` // Unique identifier of AWS security group in this access list entry.
	CIDRBlock        string `json:"cidrBlock,omitempty"`        // Range of IP addresses in CIDR notation in this access list entry.
	Comment          string `json:"comment,omitempty"`          // Comment associated with this access list entry.
	DeleteAfterDate  string `json:"deleteAfterDate,omitempty"`  // Timestamp in ISO 8601 date and time format in UTC after which Atlas deletes the temporary access list entry. Atlas returns this field if you specified an expiration date when creating this access list entry.
	GroupID          string `json:"groupId,omitempty"`          // Unique identifier of the project to which this access list entry applies.
	IPAddress        string `json:"ipAddress,omitempty"`        // Entry using an IP address in this access list entry.
}

// ProjectIPAccessLists is the response from the ProjectIPAccessListService.List.
type ProjectIPAccessLists struct {
	Links      []*Link               `json:"links"`
	Results    []ProjectIPAccessList `json:"results"`
	TotalCount int                   `json:"totalCount"`
}

// List all access list entries in the project associated to {PROJECT-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/ip-access-list/get-all-access-list-entries/
func (s *ProjectIPAccessListServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) (*ProjectIPAccessLists, *Response, error) {
	path := fmt.Sprintf(projectIPAccessListPath, groupID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProjectIPAccessLists)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Get the access list entry specified to {ACCESS-LIST-ENTRY} from the project associated to {PROJECT-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/ip-access-list/get-one-access-list-entry/
func (s *ProjectIPAccessListServiceOp) Get(ctx context.Context, groupID, entry string) (*ProjectIPAccessList, *Response, error) {
	if entry == "" {
		return nil, nil, NewArgError("entry", "must be set")
	}

	basePath := fmt.Sprintf(projectIPAccessListPath, groupID)
	escapedEntry := url.PathEscape(entry)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProjectIPAccessList)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create adds one or more access list entries to the project associated to {PROJECT-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/ip-access-list/add-entries-to-access-list/
func (s *ProjectIPAccessListServiceOp) Create(ctx context.Context, groupID string, createRequest []*ProjectIPAccessList) (*ProjectIPAccessLists, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(projectIPAccessListPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProjectIPAccessLists)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, err
}

// Delete the access list entry specified to {ACCESS-LIST-ENTRY} from the project associated to {PROJECT-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/ip-access-list/delete-one-access-list-entry/
func (s *ProjectIPAccessListServiceOp) Delete(ctx context.Context, groupID, entry string) (*Response, error) {
	if entry == "" {
		return nil, NewArgError("entry", "must be set")
	}

	basePath := fmt.Sprintf(projectIPAccessListPath, groupID)
	escapedEntry := url.PathEscape(entry)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
