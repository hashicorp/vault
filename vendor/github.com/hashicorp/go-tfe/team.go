// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ Teams = (*teams)(nil)

// Teams describes all the team related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/teams
type Teams interface {
	// List all the teams of the given organization.
	List(ctx context.Context, organization string, options *TeamListOptions) (*TeamList, error)

	// Create a new team with the given options.
	Create(ctx context.Context, organization string, options TeamCreateOptions) (*Team, error)

	// Read a team by its ID.
	Read(ctx context.Context, teamID string) (*Team, error)

	// Update a team by its ID.
	Update(ctx context.Context, teamID string, options TeamUpdateOptions) (*Team, error)

	// Delete a team by its ID.
	Delete(ctx context.Context, teamID string) error
}

// teams implements Teams.
type teams struct {
	client *Client
}

// TeamList represents a list of teams.
type TeamList struct {
	*Pagination
	Items []*Team
}

// Team represents a Terraform Enterprise team.
type Team struct {
	ID                 string              `jsonapi:"primary,teams"`
	IsUnified          bool                `jsonapi:"attr,is-unified"`
	Name               string              `jsonapi:"attr,name"`
	OrganizationAccess *OrganizationAccess `jsonapi:"attr,organization-access"`
	Visibility         string              `jsonapi:"attr,visibility"`
	Permissions        *TeamPermissions    `jsonapi:"attr,permissions"`
	UserCount          int                 `jsonapi:"attr,users-count"`
	SSOTeamID          string              `jsonapi:"attr,sso-team-id"`
	// AllowMemberTokenManagement is false for TFE versions older than v202408
	AllowMemberTokenManagement bool `jsonapi:"attr,allow-member-token-management"`

	// Relations
	Users                   []*User                   `jsonapi:"relation,users"`
	OrganizationMemberships []*OrganizationMembership `jsonapi:"relation,organization-memberships"`
}

// OrganizationAccess represents the team's permissions on its organization
type OrganizationAccess struct {
	ManagePolicies           bool `jsonapi:"attr,manage-policies"`
	ManagePolicyOverrides    bool `jsonapi:"attr,manage-policy-overrides"`
	ManageWorkspaces         bool `jsonapi:"attr,manage-workspaces"`
	ManageVCSSettings        bool `jsonapi:"attr,manage-vcs-settings"`
	ManageProviders          bool `jsonapi:"attr,manage-providers"`
	ManageModules            bool `jsonapi:"attr,manage-modules"`
	ManageRunTasks           bool `jsonapi:"attr,manage-run-tasks"`
	ManageProjects           bool `jsonapi:"attr,manage-projects"`
	ReadWorkspaces           bool `jsonapi:"attr,read-workspaces"`
	ReadProjects             bool `jsonapi:"attr,read-projects"`
	ManageMembership         bool `jsonapi:"attr,manage-membership"`
	ManageTeams              bool `jsonapi:"attr,manage-teams"`
	ManageOrganizationAccess bool `jsonapi:"attr,manage-organization-access"`
	AccessSecretTeams        bool `jsonapi:"attr,access-secret-teams"`
	ManageAgentPools         bool `jsonapi:"attr,manage-agent-pools"`
}

// TeamPermissions represents the current user's permissions on the team.
type TeamPermissions struct {
	CanDestroy          bool `jsonapi:"attr,can-destroy"`
	CanUpdateMembership bool `jsonapi:"attr,can-update-membership"`
}

// TeamIncludeOpt represents the available options for include query params.
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/teams#available-related-resources
type TeamIncludeOpt string

const (
	TeamUsers                   TeamIncludeOpt = "users"
	TeamOrganizationMemberships TeamIncludeOpt = "organization-memberships"
)

// TeamListOptions represents the options for listing teams.
type TeamListOptions struct {
	ListOptions
	// Optional: A list of relations to include.
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/teams#available-related-resources
	Include []TeamIncludeOpt `url:"include,omitempty"`

	// Optional: A list of team names to filter by.
	Names []string `url:"filter[names],omitempty"`

	// Optional: A query string to search teams by names.
	Query string `url:"q,omitempty"`
}

// TeamCreateOptions represents the options for creating a team.
type TeamCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,teams"`

	// Name of the team.
	Name *string `jsonapi:"attr,name"`

	// Optional: Unique Identifier to control team membership via SAML
	SSOTeamID *string `jsonapi:"attr,sso-team-id,omitempty"`

	// The team's organization access
	OrganizationAccess *OrganizationAccessOptions `jsonapi:"attr,organization-access,omitempty"`

	// The team's visibility ("secret", "organization")
	Visibility *string `jsonapi:"attr,visibility,omitempty"`

	// Optional: Used by Owners and users with "Manage Teams" permissions to control whether team members can manage team tokens
	AllowMemberTokenManagement *bool `jsonapi:"attr,allow-member-token-management,omitempty"`
}

// TeamUpdateOptions represents the options for updating a team.
type TeamUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,teams"`

	// Optional: New name for the team
	Name *string `jsonapi:"attr,name,omitempty"`

	// Optional: Unique Identifier to control team membership via SAML
	SSOTeamID *string `jsonapi:"attr,sso-team-id,omitempty"`

	// Optional: The team's organization access
	OrganizationAccess *OrganizationAccessOptions `jsonapi:"attr,organization-access,omitempty"`

	// Optional: The team's visibility ("secret", "organization")
	Visibility *string `jsonapi:"attr,visibility,omitempty"`

	// Optional: Used by Owners and users with "Manage Teams" permissions to control whether team members can manage team tokens
	AllowMemberTokenManagement *bool `jsonapi:"attr,allow-member-token-management,omitempty"`
}

// OrganizationAccessOptions represents the organization access options of a team.
type OrganizationAccessOptions struct {
	ManagePolicies           *bool `json:"manage-policies,omitempty"`
	ManagePolicyOverrides    *bool `json:"manage-policy-overrides,omitempty"`
	ManageWorkspaces         *bool `json:"manage-workspaces,omitempty"`
	ManageVCSSettings        *bool `json:"manage-vcs-settings,omitempty"`
	ManageProviders          *bool `json:"manage-providers,omitempty"`
	ManageModules            *bool `json:"manage-modules,omitempty"`
	ManageRunTasks           *bool `json:"manage-run-tasks,omitempty"`
	ManageProjects           *bool `json:"manage-projects,omitempty"`
	ReadWorkspaces           *bool `json:"read-workspaces,omitempty"`
	ReadProjects             *bool `json:"read-projects,omitempty"`
	ManageMembership         *bool `json:"manage-membership,omitempty"`
	ManageTeams              *bool `json:"manage-teams,omitempty"`
	ManageOrganizationAccess *bool `json:"manage-organization-access,omitempty"`
	AccessSecretTeams        *bool `json:"access-secret-teams,omitempty"`
	ManageAgentPools         *bool `json:"manage-agent-pools,omitempty"`
}

// List all the teams of the given organization.
func (s *teams) List(ctx context.Context, organization string, options *TeamListOptions) (*TeamList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}
	u := fmt.Sprintf("organizations/%s/teams", url.PathEscape(organization))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	tl := &TeamList{}
	err = req.Do(ctx, tl)
	if err != nil {
		return nil, err
	}

	return tl, nil
}

// Create a new team with the given options.
func (s *teams) Create(ctx context.Context, organization string, options TeamCreateOptions) (*Team, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/teams", url.PathEscape(organization))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	t := &Team{}
	err = req.Do(ctx, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Read a single team by its ID.
func (s *teams) Read(ctx context.Context, teamID string) (*Team, error) {
	if !validStringID(&teamID) {
		return nil, ErrInvalidTeamID
	}

	u := fmt.Sprintf("teams/%s", url.PathEscape(teamID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	t := &Team{}
	err = req.Do(ctx, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Update a team by its ID.
func (s *teams) Update(ctx context.Context, teamID string, options TeamUpdateOptions) (*Team, error) {
	if !validStringID(&teamID) {
		return nil, ErrInvalidTeamID
	}

	u := fmt.Sprintf("teams/%s", url.PathEscape(teamID))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	t := &Team{}
	err = req.Do(ctx, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Delete a team by its ID.
func (s *teams) Delete(ctx context.Context, teamID string) error {
	if !validStringID(&teamID) {
		return ErrInvalidTeamID
	}

	u := fmt.Sprintf("teams/%s", url.PathEscape(teamID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o TeamCreateOptions) valid() error {
	if !validString(o.Name) {
		return ErrRequiredName
	}
	return nil
}

func (o *TeamListOptions) valid() error {
	if o == nil {
		return nil // nothing to validate
	}

	if err := validateTeamNames(o.Names); err != nil {
		return err
	}

	return nil
}

func validateTeamNames(names []string) error {
	for _, name := range names {
		if name == "" {
			return ErrEmptyTeamName
		}
	}

	return nil
}
