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
// TFE API docs: https://www.terraform.io/docs/cloud/api/admin/workspaces.html
type AdminWorkspaces interface {
	// List all the workspaces within a workspace.
	List(ctx context.Context, options AdminWorkspaceListOptions) (*AdminWorkspaceList, error)

	// Read a workspace by its ID.
	Read(ctx context.Context, workspaceID string) (*AdminWorkspace, error)

	// Delete a workspace by its ID.
	Delete(ctx context.Context, workspaceID string) error
}

// adminWorkspaces implements AdminWorkspaces.
type adminWorkspaces struct {
	client *Client
}

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

// AdminWorkspaceListOptions represents the options for listing workspaces.
type AdminWorkspaceListOptions struct {
	ListOptions

	// A query string (partial workspace name) used to filter the results.
	// https://www.terraform.io/docs/cloud/api/admin/workspaces.html#query-parameters
	Query *string `url:"q,omitempty"`

	// A list of relations to include. See available resources
	// https://www.terraform.io/docs/cloud/api/admin/workspaces.html#available-related-resources
	Include *string `url:"include"`
}

// AdminWorkspaceList represents a list of workspaces.
type AdminWorkspaceList struct {
	*Pagination
	Items []*AdminWorkspace
}

// List all the workspaces within a worksapce.
func (s *adminWorkspaces) List(ctx context.Context, options AdminWorkspaceListOptions) (*AdminWorkspaceList, error) {
	u := "admin/workspaces"
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	awl := &AdminWorkspaceList{}
	err = s.client.do(ctx, req, awl)
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

	u := fmt.Sprintf("admin/workspaces/%s", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	aw := &AdminWorkspace{}
	err = s.client.do(ctx, req, aw)
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

	u := fmt.Sprintf("admin/workspaces/%s", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
