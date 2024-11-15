// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"time"
)

// OrganizationsService provides access to the organization related functions
// in the GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/orgs/
type OrganizationsService service

// Organization represents a GitHub organization account.
type Organization struct {
	Login             *string    `json:"login,omitempty"`
	ID                *int64     `json:"id,omitempty"`
	NodeID            *string    `json:"node_id,omitempty"`
	AvatarURL         *string    `json:"avatar_url,omitempty"`
	HTMLURL           *string    `json:"html_url,omitempty"`
	Name              *string    `json:"name,omitempty"`
	Company           *string    `json:"company,omitempty"`
	Blog              *string    `json:"blog,omitempty"`
	Location          *string    `json:"location,omitempty"`
	Email             *string    `json:"email,omitempty"`
	Description       *string    `json:"description,omitempty"`
	PublicRepos       *int       `json:"public_repos,omitempty"`
	PublicGists       *int       `json:"public_gists,omitempty"`
	Followers         *int       `json:"followers,omitempty"`
	Following         *int       `json:"following,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	TotalPrivateRepos *int       `json:"total_private_repos,omitempty"`
	OwnedPrivateRepos *int       `json:"owned_private_repos,omitempty"`
	PrivateGists      *int       `json:"private_gists,omitempty"`
	DiskUsage         *int       `json:"disk_usage,omitempty"`
	Collaborators     *int       `json:"collaborators,omitempty"`
	BillingEmail      *string    `json:"billing_email,omitempty"`
	Type              *string    `json:"type,omitempty"`
	Plan              *Plan      `json:"plan,omitempty"`

	// API URLs
	URL              *string `json:"url,omitempty"`
	EventsURL        *string `json:"events_url,omitempty"`
	HooksURL         *string `json:"hooks_url,omitempty"`
	IssuesURL        *string `json:"issues_url,omitempty"`
	MembersURL       *string `json:"members_url,omitempty"`
	PublicMembersURL *string `json:"public_members_url,omitempty"`
	ReposURL         *string `json:"repos_url,omitempty"`
}

func (o Organization) String() string {
	return Stringify(o)
}

// Plan represents the payment plan for an account. See plans at https://github.com/plans.
type Plan struct {
	Name          *string `json:"name,omitempty"`
	Space         *int    `json:"space,omitempty"`
	Collaborators *int    `json:"collaborators,omitempty"`
	PrivateRepos  *int    `json:"private_repos,omitempty"`
}

func (p Plan) String() string {
	return Stringify(p)
}

// OrganizationsListOptions specifies the optional parameters to the
// OrganizationsService.ListAll method.
type OrganizationsListOptions struct {
	// Since filters Organizations by ID.
	Since int64 `url:"since,omitempty"`

	// Note: Pagination is powered exclusively by the Since parameter,
	// ListOptions.Page has no effect.
	// ListOptions.PerPage controls an undocumented GitHub API parameter.
	ListOptions
}

// ListAll lists all organizations, in the order that they were created on GitHub.
//
// Note: Pagination is powered exclusively by the since parameter. To continue
// listing the next set of organizations, use the ID of the last-returned organization
// as the opts.Since parameter for the next call.
//
// GitHub API docs: https://developer.github.com/v3/orgs/#list-all-organizations
func (s *OrganizationsService) ListAll(ctx context.Context, opt *OrganizationsListOptions) ([]*Organization, *Response, error) {
	u, err := addOptions("organizations", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	orgs := []*Organization{}
	resp, err := s.client.Do(ctx, req, &orgs)
	if err != nil {
		return nil, resp, err
	}
	return orgs, resp, nil
}

// List the organizations for a user. Passing the empty string will list
// organizations for the authenticated user.
//
// GitHub API docs: https://developer.github.com/v3/orgs/#list-user-organizations
func (s *OrganizationsService) List(ctx context.Context, user string, opt *ListOptions) ([]*Organization, *Response, error) {
	var u string
	if user != "" {
		u = fmt.Sprintf("users/%v/orgs", user)
	} else {
		u = "user/orgs"
	}
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var orgs []*Organization
	resp, err := s.client.Do(ctx, req, &orgs)
	if err != nil {
		return nil, resp, err
	}

	return orgs, resp, nil
}

// Get fetches an organization by name.
//
// GitHub API docs: https://developer.github.com/v3/orgs/#get-an-organization
func (s *OrganizationsService) Get(ctx context.Context, org string) (*Organization, *Response, error) {
	u := fmt.Sprintf("orgs/%v", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	organization := new(Organization)
	resp, err := s.client.Do(ctx, req, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, nil
}

// GetByID fetches an organization.
//
// Note: GetByID uses the undocumented GitHub API endpoint /organizations/:id.
func (s *OrganizationsService) GetByID(ctx context.Context, id int64) (*Organization, *Response, error) {
	u := fmt.Sprintf("organizations/%d", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	organization := new(Organization)
	resp, err := s.client.Do(ctx, req, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, nil
}

// Edit an organization.
//
// GitHub API docs: https://developer.github.com/v3/orgs/#edit-an-organization
func (s *OrganizationsService) Edit(ctx context.Context, name string, org *Organization) (*Organization, *Response, error) {
	u := fmt.Sprintf("orgs/%v", name)
	req, err := s.client.NewRequest("PATCH", u, org)
	if err != nil {
		return nil, nil, err
	}

	o := new(Organization)
	resp, err := s.client.Do(ctx, req, o)
	if err != nil {
		return nil, resp, err
	}

	return o, resp, nil
}
