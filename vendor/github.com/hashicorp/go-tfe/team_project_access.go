// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ TeamProjectAccesses = (*teamProjectAccesses)(nil)

// TeamProjectAccesses describes all the team project access related methods that the Terraform
// Enterprise API supports
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/project-team-access
type TeamProjectAccesses interface {
	// List all project accesses for a given project.
	List(ctx context.Context, options TeamProjectAccessListOptions) (*TeamProjectAccessList, error)

	// Add team access for a project.
	Add(ctx context.Context, options TeamProjectAccessAddOptions) (*TeamProjectAccess, error)

	// Read team access by project ID.
	Read(ctx context.Context, teamProjectAccessID string) (*TeamProjectAccess, error)

	// Update team access on a project.
	Update(ctx context.Context, teamProjectAccessID string, options TeamProjectAccessUpdateOptions) (*TeamProjectAccess, error)

	// Remove team access from a project.
	Remove(ctx context.Context, teamProjectAccessID string) error
}

// teamProjectAccesses implements TeamProjectAccesses
type teamProjectAccesses struct {
	client *Client
}

// TeamProjectAccessType represents a team project access type.
type TeamProjectAccessType string

const (
	TeamProjectAccessAdmin    TeamProjectAccessType = "admin"
	TeamProjectAccessMaintain TeamProjectAccessType = "maintain"
	TeamProjectAccessWrite    TeamProjectAccessType = "write"
	TeamProjectAccessRead     TeamProjectAccessType = "read"
	TeamProjectAccessCustom   TeamProjectAccessType = "custom"
)

// TeamProjectAccessList represents a list of team project accesses
type TeamProjectAccessList struct {
	*Pagination
	Items []*TeamProjectAccess
}

// TeamProjectAccess represents a project access for a team
type TeamProjectAccess struct {
	ID              string                                 `jsonapi:"primary,team-projects"`
	Access          TeamProjectAccessType                  `jsonapi:"attr,access"`
	ProjectAccess   *TeamProjectAccessProjectPermissions   `jsonapi:"attr,project-access"`
	WorkspaceAccess *TeamProjectAccessWorkspacePermissions `jsonapi:"attr,workspace-access"`

	// Relations
	Team    *Team    `jsonapi:"relation,team"`
	Project *Project `jsonapi:"relation,project"`
}

// ProjectPermissions represents the team's permissions on its project
type TeamProjectAccessProjectPermissions struct {
	ProjectSettingsPermission ProjectSettingsPermissionType `jsonapi:"attr,settings"`
	ProjectTeamsPermission    ProjectTeamsPermissionType    `jsonapi:"attr,teams"`
}

// WorkspacePermissions represents the team's permission on all workspaces in its project
type TeamProjectAccessWorkspacePermissions struct {
	WorkspaceRunsPermission          WorkspaceRunsPermissionType          `jsonapi:"attr,runs"`
	WorkspaceSentinelMocksPermission WorkspaceSentinelMocksPermissionType `jsonapi:"attr,sentinel-mocks"`
	WorkspaceStateVersionsPermission WorkspaceStateVersionsPermissionType `jsonapi:"attr,state-versions"`
	WorkspaceVariablesPermission     WorkspaceVariablesPermissionType     `jsonapi:"attr,variables"`
	WorkspaceCreatePermission        bool                                 `jsonapi:"attr,create"`
	WorkspaceLockingPermission       bool                                 `jsonapi:"attr,locking"`
	WorkspaceMovePermission          bool                                 `jsonapi:"attr,move"`
	WorkspaceDeletePermission        bool                                 `jsonapi:"attr,delete"`
	WorkspaceRunTasksPermission      bool                                 `jsonapi:"attr,run-tasks"`
}

// ProjectSettingsPermissionType represents the permissiontype to a project's settings
type ProjectSettingsPermissionType string

const (
	ProjectSettingsPermissionRead   ProjectSettingsPermissionType = "read"
	ProjectSettingsPermissionUpdate ProjectSettingsPermissionType = "update"
	ProjectSettingsPermissionDelete ProjectSettingsPermissionType = "delete"
)

// ProjectTeamsPermissionType represents the permissiontype to a project's teams
type ProjectTeamsPermissionType string

const (
	ProjectTeamsPermissionNone   ProjectTeamsPermissionType = "none"
	ProjectTeamsPermissionRead   ProjectTeamsPermissionType = "read"
	ProjectTeamsPermissionManage ProjectTeamsPermissionType = "manage"
)

// WorkspaceRunsPermissionType represents the permissiontype to project workspaces' runs
type WorkspaceRunsPermissionType string

const (
	WorkspaceRunsPermissionRead  WorkspaceRunsPermissionType = "read"
	WorkspaceRunsPermissionPlan  WorkspaceRunsPermissionType = "plan"
	WorkspaceRunsPermissionApply WorkspaceRunsPermissionType = "apply"
)

// WorkspaceSentinelMocksPermissionType represents the permissiontype to project workspaces' sentinel-mocks
type WorkspaceSentinelMocksPermissionType string

const (
	WorkspaceSentinelMocksPermissionNone WorkspaceSentinelMocksPermissionType = "none"
	WorkspaceSentinelMocksPermissionRead WorkspaceSentinelMocksPermissionType = "read"
)

// WorkspaceStateVersionsPermissionType represents the permissiontype to project workspaces' state-versions
type WorkspaceStateVersionsPermissionType string

const (
	WorkspaceStateVersionsPermissionNone        WorkspaceStateVersionsPermissionType = "none"
	WorkspaceStateVersionsPermissionReadOutputs WorkspaceStateVersionsPermissionType = "read-outputs"
	WorkspaceStateVersionsPermissionRead        WorkspaceStateVersionsPermissionType = "read"
	WorkspaceStateVersionsPermissionWrite       WorkspaceStateVersionsPermissionType = "write"
)

// WorkspaceVariablesPermissionType represents the permissiontype to project workspaces' variables
type WorkspaceVariablesPermissionType string

const (
	WorkspaceVariablesPermissionNone  WorkspaceVariablesPermissionType = "none"
	WorkspaceVariablesPermissionRead  WorkspaceVariablesPermissionType = "read"
	WorkspaceVariablesPermissionWrite WorkspaceVariablesPermissionType = "write"
)

type TeamProjectAccessProjectPermissionsOptions struct {
	Settings *ProjectSettingsPermissionType `json:"settings,omitempty"`
	Teams    *ProjectTeamsPermissionType    `json:"teams,omitempty"`
}

type TeamProjectAccessWorkspacePermissionsOptions struct {
	Runs          *WorkspaceRunsPermissionType          `json:"runs,omitempty"`
	SentinelMocks *WorkspaceSentinelMocksPermissionType `json:"sentinel-mocks,omitempty"`
	StateVersions *WorkspaceStateVersionsPermissionType `json:"state-versions,omitempty"`
	Variables     *WorkspaceVariablesPermissionType     `json:"variables,omitempty"`
	Create        *bool                                 `json:"create,omitempty"`
	Locking       *bool                                 `json:"locking,omitempty"`
	Move          *bool                                 `json:"move,omitempty"`
	Delete        *bool                                 `json:"delete,omitempty"`
	RunTasks      *bool                                 `json:"run-tasks,omitempty"`
}

// TeamProjectAccessListOptions represents the options for listing team project accesses
type TeamProjectAccessListOptions struct {
	ListOptions
	ProjectID string `url:"filter[project][id]"`
}

// TeamProjectAccessAddOptions represents the options for adding team access for a project
type TeamProjectAccessAddOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,team-projects"`
	// The type of access to grant.
	Access TeamProjectAccessType `jsonapi:"attr,access"`
	// The levels that project and workspace permissions grant
	ProjectAccess   *TeamProjectAccessProjectPermissionsOptions   `jsonapi:"attr,project-access,omitempty"`
	WorkspaceAccess *TeamProjectAccessWorkspacePermissionsOptions `jsonapi:"attr,workspace-access,omitempty"`

	// The team to add to the project
	Team *Team `jsonapi:"relation,team"`
	// The project to which the team is to be added.
	Project *Project `jsonapi:"relation,project"`
}

// TeamProjectAccessUpdateOptions represents the options for updating a team project access
type TeamProjectAccessUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,team-projects"`
	// The type of access to grant.
	Access          *TeamProjectAccessType                        `jsonapi:"attr,access,omitempty"`
	ProjectAccess   *TeamProjectAccessProjectPermissionsOptions   `jsonapi:"attr,project-access,omitempty"`
	WorkspaceAccess *TeamProjectAccessWorkspacePermissionsOptions `jsonapi:"attr,workspace-access,omitempty"`
}

// List all team accesses for a given project.
func (s *teamProjectAccesses) List(ctx context.Context, options TeamProjectAccessListOptions) (*TeamProjectAccessList, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", "team-projects", &options)
	if err != nil {
		return nil, err
	}

	tpal := &TeamProjectAccessList{}
	err = req.Do(ctx, tpal)
	if err != nil {
		return nil, err
	}

	return tpal, nil
}

// Add team access for a project.
func (s *teamProjectAccesses) Add(ctx context.Context, options TeamProjectAccessAddOptions) (*TeamProjectAccess, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	if err := validateTeamProjectAccessType(options.Access); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", "team-projects", &options)
	if err != nil {
		return nil, err
	}

	tpa := &TeamProjectAccess{}
	err = req.Do(ctx, tpa)
	if err != nil {
		return nil, err
	}

	return tpa, nil
}

// Read a team project access by its ID.
func (s *teamProjectAccesses) Read(ctx context.Context, teamProjectAccessID string) (*TeamProjectAccess, error) {
	if !validStringID(&teamProjectAccessID) {
		return nil, ErrInvalidTeamProjectAccessID
	}

	u := fmt.Sprintf("team-projects/%s", url.PathEscape(teamProjectAccessID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	tpa := &TeamProjectAccess{}
	err = req.Do(ctx, tpa)
	if err != nil {
		return nil, err
	}

	return tpa, nil
}

// Update team access for a project.
func (s *teamProjectAccesses) Update(ctx context.Context, teamProjectAccessID string, options TeamProjectAccessUpdateOptions) (*TeamProjectAccess, error) {
	if !validStringID(&teamProjectAccessID) {
		return nil, ErrInvalidTeamProjectAccessID
	}

	if err := validateTeamProjectAccessType(*options.Access); err != nil {
		return nil, err
	}
	u := fmt.Sprintf("team-projects/%s", url.PathEscape(teamProjectAccessID))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	ta := &TeamProjectAccess{}
	err = req.Do(ctx, ta)
	if err != nil {
		return nil, err
	}

	return ta, err
}

// Remove team access from a project.
func (s *teamProjectAccesses) Remove(ctx context.Context, teamProjectAccessID string) error {
	if !validStringID(&teamProjectAccessID) {
		return ErrInvalidTeamProjectAccessID
	}

	u := fmt.Sprintf("team-projects/%s", url.PathEscape(teamProjectAccessID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o TeamProjectAccessListOptions) valid() error {
	if !validStringID(&o.ProjectID) {
		return ErrInvalidProjectID
	}

	return nil
}

func (o TeamProjectAccessAddOptions) valid() error {
	if err := validateTeamProjectAccessType(o.Access); err != nil {
		return err
	}
	if o.Team == nil {
		return ErrRequiredTeam
	}
	if o.Project == nil {
		return ErrRequiredProject
	}

	return nil
}

func validateTeamProjectAccessType(t TeamProjectAccessType) error {
	switch t {
	case TeamProjectAccessAdmin,
		TeamProjectAccessMaintain,
		TeamProjectAccessWrite,
		TeamProjectAccessRead,
		TeamProjectAccessCustom:
		// do nothing
	default:
		return ErrInvalidTeamProjectAccessType
	}
	return nil
}
