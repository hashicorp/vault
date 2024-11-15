// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ OrganizationMemberships = (*organizationMemberships)(nil)

// OrganizationMemberships describes all the organization membership related methods that
// the Terraform Enterprise API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/organization-memberships
type OrganizationMemberships interface {
	// List all the organization memberships of the given organization.
	List(ctx context.Context, organization string, options *OrganizationMembershipListOptions) (*OrganizationMembershipList, error)

	// Create a new organization membership with the given options.
	Create(ctx context.Context, organization string, options OrganizationMembershipCreateOptions) (*OrganizationMembership, error)

	// Read an organization membership by ID
	Read(ctx context.Context, organizationMembershipID string) (*OrganizationMembership, error)

	// Read an organization membership by ID with options
	ReadWithOptions(ctx context.Context, organizationMembershipID string, options OrganizationMembershipReadOptions) (*OrganizationMembership, error)

	// Delete an organization membership by its ID.
	Delete(ctx context.Context, organizationMembershipID string) error
}

// organizationMemberships implements OrganizationMemberships.
type organizationMemberships struct {
	client *Client
}

// OrganizationMembershipStatus represents an organization membership status.
type OrganizationMembershipStatus string

const (
	OrganizationMembershipActive  OrganizationMembershipStatus = "active"
	OrganizationMembershipInvited OrganizationMembershipStatus = "invited"
)

// OrganizationMembershipList represents a list of organization memberships.
type OrganizationMembershipList struct {
	*Pagination
	Items []*OrganizationMembership
}

// OrganizationMembership represents a Terraform Enterprise organization membership.
type OrganizationMembership struct {
	ID     string                       `jsonapi:"primary,organization-memberships"`
	Status OrganizationMembershipStatus `jsonapi:"attr,status"`
	Email  string                       `jsonapi:"attr,email"`

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`
	User         *User         `jsonapi:"relation,user"`
	Teams        []*Team       `jsonapi:"relation,teams"`
}

// OrgMembershipIncludeOpt represents the available options for include query params.
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/organization-memberships#available-related-resources
type OrgMembershipIncludeOpt string

const (
	OrgMembershipUser OrgMembershipIncludeOpt = "user"
	OrgMembershipTeam OrgMembershipIncludeOpt = "teams"
)

// OrganizationMembershipListOptions represents the options for listing organization memberships.
type OrganizationMembershipListOptions struct {
	ListOptions
	// Optional: A list of relations to include. See available resources
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/organization-memberships#available-related-resources
	Include []OrgMembershipIncludeOpt `url:"include,omitempty"`

	// Optional: A list of organization member emails to filter by.
	Emails []string `url:"filter[email],omitempty"`

	// Optional: If specified, restricts results to those matching status value.
	Status OrganizationMembershipStatus `url:"filter[status],omitempty"`

	// Optional: A query string to search organization memberships by user name
	// and email.
	Query string `url:"q,omitempty"`
}

// OrganizationMembershipCreateOptions represents the options for creating an organization membership.
type OrganizationMembershipCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,organization-memberships"`

	// Required: User's email address.
	Email *string `jsonapi:"attr,email"`

	// Optional: A list of teams in the organization to add the user to
	Teams []*Team `jsonapi:"relation,teams,omitempty"`
}

// OrganizationMembershipReadOptions represents the options for reading organization memberships.
type OrganizationMembershipReadOptions struct {
	// Optional: A list of relations to include. See available resources
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/organization-memberships#available-related-resources
	Include []OrgMembershipIncludeOpt `url:"include,omitempty"`
}

// List all the organization memberships of the given organization.
func (s *organizationMemberships) List(ctx context.Context, organization string, options *OrganizationMembershipListOptions) (*OrganizationMembershipList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/organization-memberships", url.PathEscape(organization))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	ml := &OrganizationMembershipList{}
	err = req.Do(ctx, ml)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

// Create an organization membership with the given options.
func (s *organizationMemberships) Create(ctx context.Context, organization string, options OrganizationMembershipCreateOptions) (*OrganizationMembership, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/organization-memberships", url.PathEscape(organization))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	m := &OrganizationMembership{}
	err = req.Do(ctx, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Read an organization membership by its ID.
func (s *organizationMemberships) Read(ctx context.Context, organizationMembershipID string) (*OrganizationMembership, error) {
	return s.ReadWithOptions(ctx, organizationMembershipID, OrganizationMembershipReadOptions{})
}

// Read an organization membership by ID with options
func (s *organizationMemberships) ReadWithOptions(ctx context.Context, organizationMembershipID string, options OrganizationMembershipReadOptions) (*OrganizationMembership, error) {
	if !validStringID(&organizationMembershipID) {
		return nil, ErrInvalidMembership
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organization-memberships/%s", url.PathEscape(organizationMembershipID))
	req, err := s.client.NewRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	mem := &OrganizationMembership{}
	err = req.Do(ctx, mem)
	if err != nil {
		return nil, err
	}

	return mem, nil
}

// Delete an organization membership by its ID.
func (s *organizationMemberships) Delete(ctx context.Context, organizationMembershipID string) error {
	if !validStringID(&organizationMembershipID) {
		return ErrInvalidMembership
	}

	u := fmt.Sprintf("organization-memberships/%s", url.PathEscape(organizationMembershipID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o OrganizationMembershipCreateOptions) valid() error {
	if o.Email == nil {
		return ErrRequiredEmail
	}
	return nil
}

func (o *OrganizationMembershipListOptions) valid() error {
	if o == nil {
		return nil
	}

	if err := validateOrgMembershipEmailParams(o.Emails); err != nil {
		return err
	}

	return nil
}

func (o OrganizationMembershipReadOptions) valid() error {
	return nil
}

func validateOrgMembershipEmailParams(emails []string) error {
	for _, email := range emails {
		if !validEmail(email) {
			return ErrInvalidEmail
		}
	}

	return nil
}
