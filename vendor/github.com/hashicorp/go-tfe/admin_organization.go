// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ AdminOrganizations = (*adminOrganizations)(nil)

// AdminOrganizations describes all of the admin organization related methods that the Terraform
// Enterprise API supports. Note that admin settings are only available in Terraform Enterprise.
//
// TFE API docs: https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/organizations
type AdminOrganizations interface {
	// List all the organizations visible to the current user.
	List(ctx context.Context, options *AdminOrganizationListOptions) (*AdminOrganizationList, error)

	// Read attributes of an existing organization via admin API.
	Read(ctx context.Context, organization string) (*AdminOrganization, error)

	// Update attributes of an existing organization via admin API.
	Update(ctx context.Context, organization string, options AdminOrganizationUpdateOptions) (*AdminOrganization, error)

	// Delete an organization by its name via admin API
	Delete(ctx context.Context, organization string) error

	// ListModuleConsumers lists specific organizations in the Terraform Enterprise installation that have permission to use an organization's modules.
	ListModuleConsumers(ctx context.Context, organization string, options *AdminOrganizationListModuleConsumersOptions) (*AdminOrganizationList, error)

	// UpdateModuleConsumers specifies a list of organizations that can use modules from the sharing organization's private registry. Setting a list of module consumers will turn off global module sharing for an organization.
	UpdateModuleConsumers(ctx context.Context, organization string, consumerOrganizations []string) error
}

// adminOrganizations implements AdminOrganizations.
type adminOrganizations struct {
	client *Client
}

// AdminOrganization represents a Terraform Enterprise organization returned from the Admin API.
type AdminOrganization struct {
	Name                             string `jsonapi:"primary,organizations"`
	AccessBetaTools                  bool   `jsonapi:"attr,access-beta-tools"`
	ExternalID                       string `jsonapi:"attr,external-id"`
	GlobalModuleSharing              *bool  `jsonapi:"attr,global-module-sharing"`
	GlobalProviderSharing            *bool  `jsonapi:"attr,global-provider-sharing"`
	IsDisabled                       bool   `jsonapi:"attr,is-disabled"`
	NotificationEmail                string `jsonapi:"attr,notification-email"`
	SsoEnabled                       bool   `jsonapi:"attr,sso-enabled"`
	TerraformBuildWorkerApplyTimeout string `jsonapi:"attr,terraform-build-worker-apply-timeout"`
	TerraformBuildWorkerPlanTimeout  string `jsonapi:"attr,terraform-build-worker-plan-timeout"`
	ApplyTimeout                     string `jsonapi:"attr,apply-timeout"`
	PlanTimeout                      string `jsonapi:"attr,plan-timeout"`
	TerraformWorkerSudoEnabled       bool   `jsonapi:"attr,terraform-worker-sudo-enabled"`
	WorkspaceLimit                   *int   `jsonapi:"attr,workspace-limit"`

	// Relations
	Owners []*User `jsonapi:"relation,owners"`
}

// AdminOrganizationUpdateOptions represents the admin options for updating an organization.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/organizations#request-body
type AdminOrganizationUpdateOptions struct {
	AccessBetaTools                  *bool   `jsonapi:"attr,access-beta-tools,omitempty"`
	GlobalModuleSharing              *bool   `jsonapi:"attr,global-module-sharing,omitempty"`
	GlobalProviderSharing            *bool   `jsonapi:"attr,global-provider-sharing,omitempty"`
	IsDisabled                       *bool   `jsonapi:"attr,is-disabled,omitempty"`
	TerraformBuildWorkerApplyTimeout *string `jsonapi:"attr,terraform-build-worker-apply-timeout,omitempty"`
	TerraformBuildWorkerPlanTimeout  *string `jsonapi:"attr,terraform-build-worker-plan-timeout,omitempty"`
	ApplyTimeout                     *string `jsonapi:"attr,apply-timeout,omitempty"`
	PlanTimeout                      *string `jsonapi:"attr,plan-timeout,omitempty"`
	TerraformWorkerSudoEnabled       bool    `jsonapi:"attr,terraform-worker-sudo-enabled,omitempty"`
	WorkspaceLimit                   *int    `jsonapi:"attr,workspace-limit,omitempty"`
}

// AdminOrganizationList represents a list of organizations via Admin API.
type AdminOrganizationList struct {
	*Pagination
	Items []*AdminOrganization
}

// AdminOrgIncludeOpt represents the available options for include query params.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/organizations#available-related-resources
type AdminOrgIncludeOpt string

const AdminOrgOwners AdminOrgIncludeOpt = "owners"

// AdminOrganizationListOptions represents the options for listing organizations via Admin API.
type AdminOrganizationListOptions struct {
	ListOptions

	// Optional: A query string used to filter organizations.
	// Any organizations with a name or notification email partially matching this value will be returned.
	Query string `url:"q,omitempty"`
	// Optional: A list of relations to include. See available resources
	// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/organizations#available-related-resources
	Include []AdminOrgIncludeOpt `url:"include,omitempty"`
}

// AdminOrganizationListModuleConsumersOptions represents the options for listing organization module consumers through the Admin API
type AdminOrganizationListModuleConsumersOptions struct {
	ListOptions
}

type AdminOrganizationID struct {
	ID string `jsonapi:"primary,organizations"`
}

// List all the organizations visible to the current user.
func (s *adminOrganizations) List(ctx context.Context, options *AdminOrganizationListOptions) (*AdminOrganizationList, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}
	u := "admin/organizations"
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	orgl := &AdminOrganizationList{}
	err = req.Do(ctx, orgl)
	if err != nil {
		return nil, err
	}

	return orgl, nil
}

// ListModuleConsumers lists specific organizations in the Terraform Enterprise installation that have permission to use an organization's modules.
func (s *adminOrganizations) ListModuleConsumers(ctx context.Context, organization string, options *AdminOrganizationListModuleConsumersOptions) (*AdminOrganizationList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("admin/organizations/%s/relationships/module-consumers", url.PathEscape(organization))

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	orgl := &AdminOrganizationList{}
	err = req.Do(ctx, orgl)
	if err != nil {
		return nil, err
	}

	return orgl, nil
}

// Read an organization by its name.
func (s *adminOrganizations) Read(ctx context.Context, organization string) (*AdminOrganization, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("admin/organizations/%s", url.PathEscape(organization))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	org := &AdminOrganization{}
	err = req.Do(ctx, org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// Update an organization by its name.
func (s *adminOrganizations) Update(ctx context.Context, organization string, options AdminOrganizationUpdateOptions) (*AdminOrganization, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("admin/organizations/%s", url.PathEscape(organization))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	org := &AdminOrganization{}
	err = req.Do(ctx, org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// UpdateModuleConsumers updates an organization to specify a list of organizations that can use modules from the sharing organization's private registry.
func (s *adminOrganizations) UpdateModuleConsumers(ctx context.Context, organization string, consumerOrganizationIDs []string) error {
	if !validStringID(&organization) {
		return ErrInvalidOrg
	}

	u := fmt.Sprintf("admin/organizations/%s/relationships/module-consumers", url.PathEscape(organization))

	var organizations []*AdminOrganizationID
	for _, id := range consumerOrganizationIDs {
		if !validStringID(&id) {
			return ErrInvalidOrg
		}
		organizations = append(organizations, &AdminOrganizationID{ID: id})
	}

	req, err := s.client.NewRequest("PATCH", u, organizations)
	if err != nil {
		return err
	}

	err = req.Do(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

// Delete an organization by its name.
func (s *adminOrganizations) Delete(ctx context.Context, organization string) error {
	if !validStringID(&organization) {
		return ErrInvalidOrg
	}

	u := fmt.Sprintf("admin/organizations/%s", url.PathEscape(organization))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o *AdminOrganizationListOptions) valid() error {
	return nil
}
