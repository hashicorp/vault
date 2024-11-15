// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ AdminWorkspaces = (*adminWorkspaces)(nil)

// AdminWorkspaces describes all the admin workspace related methods that the Terraform Enterprise API supports.
// Note that admin settings are only available in Terraform Enterprise.
//
// TFE API docs: https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/workspaces
type AdminWorkspaces interface {
	// List all the workspaces within a workspace.
	List(ctx context.Context, options *AdminWorkspaceListOptions) (*AdminWorkspaceList, error)

	// Read a workspace by its ID.
	Read(ctx context.Context, workspaceID string) (*AdminWorkspace, error)

	// Delete a workspace by its ID.
	Delete(ctx context.Context, workspaceID string) error
}

// adminWorkspaces implements AdminWorkspaces interface.
type adminWorkspaces struct {
	client *Client
}

// AdminVCSRepo represents a VCS repository
type AdminVCSRepo struct {
	Identifier string `jsonapi:"attr,identifier"`
}

// AdminWorkspaces represents a Terraform Enterprise admin workspace.
type AdminWorkspace struct {
	ID      string        `jsonapi:"primary,workspaces"`
	Name    string        `jsonapi:"attr,name"`
	Locked  bool          `jsonapi:"attr,locked"`
	VCSRepo *AdminVCSRepo `jsonapi:"attr,vcs-repo"`

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`
	CurrentRun   *Run          `jsonapi:"relation,current-run"`
}

// AdminWorkspaceIncludeOpt represents the available options for include query params.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/workspaces#available-related-resources
type AdminWorkspaceIncludeOpt string

const (
	AdminWorkspaceOrg        AdminWorkspaceIncludeOpt = "organization"
	AdminWorkspaceCurrentRun AdminWorkspaceIncludeOpt = "current_run"
	AdminWorkspaceOrgOwners  AdminWorkspaceIncludeOpt = "organization.owners"
)

// AdminWorkspaceListOptions represents the options for listing workspaces.
type AdminWorkspaceListOptions struct {
	ListOptions

	// A query string (partial workspace name) used to filter the results.
	// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/workspaces#query-parameters
	Query string `url:"q,omitempty"`

	// Optional: A list of relations to include. See available resources
	// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/workspaces#available-related-resources
	Include []AdminWorkspaceIncludeOpt `url:"include,omitempty"`

	// Optional: A comma-separated list of Run statuses to restrict results. See available resources
	// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/workspaces#query-parameters
	Filter string `url:"filter[current_run][status],omitempty"`

	// Optional: May sort on "name" (the default) and "current-run.created-at" (which sorts by the time of the current run)
	// Prepending a hyphen to the sort parameter will reverse the order (e.g. "-name" to reverse the default order)
	Sort string `url:"sort,omitempty"`
}

// AdminWorkspaceList represents a list of workspaces.
type AdminWorkspaceList struct {
	*Pagination
	Items []*AdminWorkspace
}

// List all the workspaces within a workspace.
func (s *adminWorkspaces) List(ctx context.Context, options *AdminWorkspaceListOptions) (*AdminWorkspaceList, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := "admin/workspaces"
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	awl := &AdminWorkspaceList{}
	err = req.Do(ctx, awl)
	if err != nil {
		return nil, err
	}

	return awl, nil
}

// Read a workspace by its ID.
func (s *adminWorkspaces) Read(ctx context.Context, workspaceID string) (*AdminWorkspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceValue
	}

	u := fmt.Sprintf("admin/workspaces/%s", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	aw := &AdminWorkspace{}
	err = req.Do(ctx, aw)
	if err != nil {
		return nil, err
	}

	return aw, nil
}

// Delete a workspace by its ID.
func (s *adminWorkspaces) Delete(ctx context.Context, workspaceID string) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceValue
	}

	u := fmt.Sprintf("admin/workspaces/%s", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o *AdminWorkspaceListOptions) valid() error {
	return nil
}
