// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ TeamMembers = (*teamMembers)(nil)

// TeamMembers describes all the team member related methods that the
// Terraform Enterprise API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/team-members
type TeamMembers interface {
	// List returns all Users of a team calling ListUsers
	// See ListOrganizationMemberships for fetching memberships
	List(ctx context.Context, teamID string) ([]*User, error)

	// ListUsers returns the Users of this team.
	ListUsers(ctx context.Context, teamID string) ([]*User, error)

	// ListOrganizationMemberships returns the OrganizationMemberships of this team.
	ListOrganizationMemberships(ctx context.Context, teamID string) ([]*OrganizationMembership, error)

	// Add multiple users to a team.
	Add(ctx context.Context, teamID string, options TeamMemberAddOptions) error

	// Remove multiple users from a team.
	Remove(ctx context.Context, teamID string, options TeamMemberRemoveOptions) error
}

// teamMembers implements TeamMembers.
type teamMembers struct {
	client *Client
}

type teamMemberUser struct {
	Username string `jsonapi:"primary,users"`
}

type teamMemberOrgMembership struct {
	ID string `jsonapi:"primary,organization-memberships"`
}

// TeamMemberAddOptions represents the options for
// adding or removing team members.
type TeamMemberAddOptions struct {
	Usernames                 []string
	OrganizationMembershipIDs []string
}

// TeamMemberRemoveOptions represents the options for
// adding or removing team members.
type TeamMemberRemoveOptions struct {
	Usernames                 []string
	OrganizationMembershipIDs []string
}

// List returns all Users of a team calling ListUsers
// See ListOrganizationMemberships for fetching memberships
func (s *teamMembers) List(ctx context.Context, teamID string) ([]*User, error) {
	return s.ListUsers(ctx, teamID)
}

// ListUsers returns the Users of this team.
func (s *teamMembers) ListUsers(ctx context.Context, teamID string) ([]*User, error) {
	if !validStringID(&teamID) {
		return nil, ErrInvalidTeamID
	}

	options := struct {
		Include []TeamIncludeOpt `url:"include,omitempty"`
	}{
		Include: []TeamIncludeOpt{TeamUsers},
	}

	u := fmt.Sprintf("teams/%s", url.PathEscape(teamID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	t := &Team{}
	err = req.Do(ctx, t)
	if err != nil {
		return nil, err
	}

	return t.Users, nil
}

// ListOrganizationMemberships returns the OrganizationMemberships of this team.
func (s *teamMembers) ListOrganizationMemberships(ctx context.Context, teamID string) ([]*OrganizationMembership, error) {
	if !validStringID(&teamID) {
		return nil, ErrInvalidTeamID
	}

	options := struct {
		Include []TeamIncludeOpt `url:"include,omitempty"`
	}{
		Include: []TeamIncludeOpt{TeamOrganizationMemberships},
	}

	u := fmt.Sprintf("teams/%s", url.PathEscape(teamID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	t := &Team{}
	err = req.Do(ctx, t)
	if err != nil {
		return nil, err
	}

	return t.OrganizationMemberships, nil
}

// Add multiple users to a team.
func (s *teamMembers) Add(ctx context.Context, teamID string, options TeamMemberAddOptions) error {
	if !validStringID(&teamID) {
		return ErrInvalidTeamID
	}
	if err := options.valid(); err != nil {
		return err
	}

	usersOrMemberships := options.kind()
	u := fmt.Sprintf("teams/%s/relationships/%s", url.PathEscape(teamID), usersOrMemberships)

	var req *ClientRequest

	if usersOrMemberships == "users" {
		var err error
		var members []*teamMemberUser
		for _, name := range options.Usernames {
			members = append(members, &teamMemberUser{Username: name})
		}
		req, err = s.client.NewRequest("POST", u, members)
		if err != nil {
			return err
		}
	} else {
		var err error
		var members []*teamMemberOrgMembership
		for _, ID := range options.OrganizationMembershipIDs {
			members = append(members, &teamMemberOrgMembership{ID: ID})
		}
		req, err = s.client.NewRequest("POST", u, members)
		if err != nil {
			return err
		}
	}

	return req.Do(ctx, nil)
}

// Remove multiple users from a team.
func (s *teamMembers) Remove(ctx context.Context, teamID string, options TeamMemberRemoveOptions) error {
	if !validStringID(&teamID) {
		return ErrInvalidTeamID
	}
	if err := options.valid(); err != nil {
		return err
	}

	usersOrMemberships := options.kind()
	u := fmt.Sprintf("teams/%s/relationships/%s", url.PathEscape(teamID), usersOrMemberships)

	var req *ClientRequest

	if usersOrMemberships == "users" {
		var err error
		var members []*teamMemberUser
		for _, name := range options.Usernames {
			members = append(members, &teamMemberUser{Username: name})
		}
		req, err = s.client.NewRequest("DELETE", u, members)
		if err != nil {
			return err
		}
	} else {
		var err error
		var members []*teamMemberOrgMembership
		for _, ID := range options.OrganizationMembershipIDs {
			members = append(members, &teamMemberOrgMembership{ID: ID})
		}
		req, err = s.client.NewRequest("DELETE", u, members)
		if err != nil {
			return err
		}
	}

	return req.Do(ctx, nil)
}

// kind returns "users" or "organization-memberships"
// depending on which is defined
func (o *TeamMemberAddOptions) kind() string {
	if o.Usernames != nil && len(o.Usernames) != 0 {
		return "users"
	}
	return "organization-memberships"
}

// kind returns "users" or "organization-memberships"
// depending on which is defined
func (o *TeamMemberRemoveOptions) kind() string {
	if o.Usernames != nil && len(o.Usernames) != 0 {
		return "users"
	}
	return "organization-memberships"
}

func (o *TeamMemberAddOptions) valid() error {
	if o.Usernames == nil && o.OrganizationMembershipIDs == nil {
		return ErrRequiredUsernameOrMembershipIds
	}
	if o.Usernames != nil && o.OrganizationMembershipIDs != nil {
		return ErrRequiredOnlyOneField
	}
	if o.Usernames != nil && len(o.Usernames) == 0 {
		return ErrInvalidUsernames
	}
	if o.OrganizationMembershipIDs != nil && len(o.OrganizationMembershipIDs) == 0 {
		return ErrInvalidMembershipIDs
	}
	return nil
}

func (o *TeamMemberRemoveOptions) valid() error {
	if o.Usernames == nil && o.OrganizationMembershipIDs == nil {
		return ErrRequiredUsernameOrMembershipIds
	}
	if o.Usernames != nil && o.OrganizationMembershipIDs != nil {
		return ErrRequiredOnlyOneField
	}
	if o.Usernames != nil && len(o.Usernames) == 0 {
		return ErrInvalidUsernames
	}
	if o.OrganizationMembershipIDs != nil && len(o.OrganizationMembershipIDs) == 0 {
		return ErrInvalidMembershipIDs
	}
	return nil
}
