// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ TeamAccesses = (*teamAccesses)(nil)

// TeamAccesses describes all the team access related methods that the
// Terraform Enterprise API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/team-access
type TeamAccesses interface {
	// List all the team accesses for a given workspace.
	List(ctx context.Context, options *TeamAccessListOptions) (*TeamAccessList, error)

	// Add team access for a workspace.
	Add(ctx context.Context, options TeamAccessAddOptions) (*TeamAccess, error)

	// Read a team access by its ID.
	Read(ctx context.Context, teamAccessID string) (*TeamAccess, error)

	// Update a team access by its ID.
	Update(ctx context.Context, teamAccessID string, options TeamAccessUpdateOptions) (*TeamAccess, error)

	// Remove team access from a workspace.
	Remove(ctx context.Context, teamAccessID string) error
}

// teamAccesses implements TeamAccesses.
type teamAccesses struct {
	client *Client
}

// AccessType represents a team access type.
type AccessType string

const (
	AccessAdmin  AccessType = "admin"
	AccessPlan   AccessType = "plan"
	AccessRead   AccessType = "read"
	AccessWrite  AccessType = "write"
	AccessCustom AccessType = "custom"
)

// RunsPermissionType represents the permissiontype to a workspace's runs.
type RunsPermissionType string

const (
	RunsPermissionRead  RunsPermissionType = "read"
	RunsPermissionPlan  RunsPermissionType = "plan"
	RunsPermissionApply RunsPermissionType = "apply"
)

// VariablesPermissionType represents the permissiontype to a workspace's variables.
type VariablesPermissionType string

const (
	VariablesPermissionNone  VariablesPermissionType = "none"
	VariablesPermissionRead  VariablesPermissionType = "read"
	VariablesPermissionWrite VariablesPermissionType = "write"
)

// StateVersionsPermissionType represents the permissiontype to a workspace's state versions.
type StateVersionsPermissionType string

const (
	StateVersionsPermissionNone        StateVersionsPermissionType = "none"
	StateVersionsPermissionReadOutputs StateVersionsPermissionType = "read-outputs"
	StateVersionsPermissionRead        StateVersionsPermissionType = "read"
	StateVersionsPermissionWrite       StateVersionsPermissionType = "write"
)

// SentinelMocksPermissionType represents the permissiontype to a workspace's Sentinel mocks.
type SentinelMocksPermissionType string

const (
	SentinelMocksPermissionNone SentinelMocksPermissionType = "none"
	SentinelMocksPermissionRead SentinelMocksPermissionType = "read"
)

// TeamAccessList represents a list of team accesses.
type TeamAccessList struct {
	*Pagination
	Items []*TeamAccess
}

// TeamAccess represents the workspace access for a team.
type TeamAccess struct {
	ID               string                      `jsonapi:"primary,team-workspaces"`
	Access           AccessType                  `jsonapi:"attr,access"`
	Runs             RunsPermissionType          `jsonapi:"attr,runs"`
	Variables        VariablesPermissionType     `jsonapi:"attr,variables"`
	StateVersions    StateVersionsPermissionType `jsonapi:"attr,state-versions"`
	SentinelMocks    SentinelMocksPermissionType `jsonapi:"attr,sentinel-mocks"`
	WorkspaceLocking bool                        `jsonapi:"attr,workspace-locking"`
	RunTasks         bool                        `jsonapi:"attr,run-tasks"`

	// Relations
	Team      *Team      `jsonapi:"relation,team"`
	Workspace *Workspace `jsonapi:"relation,workspace"`
}

// TeamAccessListOptions represents the options for listing team accesses.
type TeamAccessListOptions struct {
	ListOptions
	WorkspaceID string `url:"filter[workspace][id]"`
}

// TeamAccessAddOptions represents the options for adding team access.
type TeamAccessAddOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,team-workspaces"`

	// The type of access to grant.
	Access *AccessType `jsonapi:"attr,access"`

	// Custom workspace access permissions. These can only be edited when Access is 'custom'; otherwise, they are
	// read-only and reflect the Access level's implicit permissions.
	Runs             *RunsPermissionType          `jsonapi:"attr,runs,omitempty"`
	Variables        *VariablesPermissionType     `jsonapi:"attr,variables,omitempty"`
	StateVersions    *StateVersionsPermissionType `jsonapi:"attr,state-versions,omitempty"`
	SentinelMocks    *SentinelMocksPermissionType `jsonapi:"attr,sentinel-mocks,omitempty"`
	WorkspaceLocking *bool                        `jsonapi:"attr,workspace-locking,omitempty"`
	RunTasks         *bool                        `jsonapi:"attr,run-tasks,omitempty"`

	// The team to add to the workspace
	Team *Team `jsonapi:"relation,team"`

	// The workspace to which the team is to be added.
	Workspace *Workspace `jsonapi:"relation,workspace"`
}

// TeamAccessUpdateOptions represents the options for updating team access.
type TeamAccessUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,team-workspaces"`

	// The type of access to grant.
	Access *AccessType `jsonapi:"attr,access,omitempty"`

	// Custom workspace access permissions. These can only be edited when Access is 'custom'; otherwise, they are
	// read-only and reflect the Access level's implicit permissions.
	Runs             *RunsPermissionType          `jsonapi:"attr,runs,omitempty"`
	Variables        *VariablesPermissionType     `jsonapi:"attr,variables,omitempty"`
	StateVersions    *StateVersionsPermissionType `jsonapi:"attr,state-versions,omitempty"`
	SentinelMocks    *SentinelMocksPermissionType `jsonapi:"attr,sentinel-mocks,omitempty"`
	WorkspaceLocking *bool                        `jsonapi:"attr,workspace-locking,omitempty"`
	RunTasks         *bool                        `jsonapi:"attr,run-tasks,omitempty"`
}

// List all the team accesses for a given workspace.
func (s *teamAccesses) List(ctx context.Context, options *TeamAccessListOptions) (*TeamAccessList, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", "team-workspaces", options)
	if err != nil {
		return nil, err
	}

	tal := &TeamAccessList{}
	err = req.Do(ctx, tal)
	if err != nil {
		return nil, err
	}

	return tal, nil
}

// Add team access for a workspace.
func (s *teamAccesses) Add(ctx context.Context, options TeamAccessAddOptions) (*TeamAccess, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", "team-workspaces", &options)
	if err != nil {
		return nil, err
	}

	ta := &TeamAccess{}
	err = req.Do(ctx, ta)
	if err != nil {
		return nil, err
	}

	return ta, nil
}

// Read a team access by its ID.
func (s *teamAccesses) Read(ctx context.Context, teamAccessID string) (*TeamAccess, error) {
	if !validStringID(&teamAccessID) {
		return nil, ErrInvalidAccessTeamID
	}

	u := fmt.Sprintf("team-workspaces/%s", url.PathEscape(teamAccessID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	ta := &TeamAccess{}
	err = req.Do(ctx, ta)
	if err != nil {
		return nil, err
	}

	return ta, nil
}

// Update team access for a workspace
func (s *teamAccesses) Update(ctx context.Context, teamAccessID string, options TeamAccessUpdateOptions) (*TeamAccess, error) {
	if !validStringID(&teamAccessID) {
		return nil, ErrInvalidAccessTeamID
	}

	u := fmt.Sprintf("team-workspaces/%s", url.PathEscape(teamAccessID))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	ta := &TeamAccess{}
	err = req.Do(ctx, ta)
	if err != nil {
		return nil, err
	}

	return ta, err
}

// Remove team access from a workspace.
func (s *teamAccesses) Remove(ctx context.Context, teamAccessID string) error {
	if !validStringID(&teamAccessID) {
		return ErrInvalidAccessTeamID
	}

	u := fmt.Sprintf("team-workspaces/%s", url.PathEscape(teamAccessID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o *TeamAccessListOptions) valid() error {
	if o == nil {
		return ErrRequiredTeamAccessListOps
	}
	if !validString(&o.WorkspaceID) {
		return ErrRequiredWorkspaceID
	}
	if !validStringID(&o.WorkspaceID) {
		return ErrInvalidWorkspaceID
	}

	return nil
}

func (o TeamAccessAddOptions) valid() error {
	if o.Access == nil {
		return ErrRequiredAccess
	}
	if o.Team == nil {
		return ErrRequiredTeam
	}
	if o.Workspace == nil {
		return ErrRequiredWorkspace
	}
	return nil
}
