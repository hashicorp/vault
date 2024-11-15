// Copyright 2022 MongoDB Inc
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
)

const federationSettingsBasePath = "api/atlas/v1.0/orgs/%s/federationSettings"
const federationSettingsDeleteBasePath = "api/atlas/v1.0/federationSettings"

// FederatedSettingsService is an interface for working with the Federation Settings
// endpoints of the MongoDB Atlas API.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/
type FederatedSettingsService interface {
	Get(context.Context, string) (*FederatedSettings, *Response, error)
	Delete(context.Context, string) (*Response, error)
	ListConnectedOrgs(context.Context, string, *ListOptions) (*FederatedSettingsConnectedOrganizations, *Response, error)
	GetConnectedOrg(context.Context, string, string) (*FederatedSettingsConnectedOrganization, *Response, error)
	UpdateConnectedOrg(context.Context, string, string, *FederatedSettingsConnectedOrganization) (*FederatedSettingsConnectedOrganization, *Response, error)
	DeleteConnectedOrg(context.Context, string, string) (*Response, error)
	ListRoleMappings(context.Context, string, string, *ListOptions) (*FederatedSettingsOrganizationRoleMappings, *Response, error)
	GetRoleMapping(context.Context, string, string, string) (*FederatedSettingsOrganizationRoleMapping, *Response, error)
	CreateRoleMapping(context.Context, string, string, *FederatedSettingsOrganizationRoleMapping) (*FederatedSettingsOrganizationRoleMapping, *Response, error)
	UpdateRoleMapping(context.Context, string, string, string, *FederatedSettingsOrganizationRoleMapping) (*FederatedSettingsOrganizationRoleMapping, *Response, error)
	DeleteRoleMapping(context.Context, string, string, string) (*Response, error)
	ListIdentityProviders(context.Context, string, *ListOptions) ([]FederatedSettingsIdentityProvider, *Response, error)
	GetIdentityProvider(context.Context, string, string) (*FederatedSettingsIdentityProvider, *Response, error)
	UpdateIdentityProvider(context.Context, string, string, *FederatedSettingsIdentityProvider) (*FederatedSettingsIdentityProvider, *Response, error)
}

// FederatedSettingsServiceOp handles communication with the FederatedSettings related methods of the
// MongoDB Atlas API.
type FederatedSettingsServiceOp service

var _ FederatedSettingsService = &FederatedSettingsServiceOp{}

// A Resource describes a specific resource the Role will allow operating on.

// FederatedSettings represents a FederatedSettings List.
type FederatedSettings struct {
	FederatedDomains       []string `json:"federatedDomains,omitempty"`
	HasRoleMappings        *bool    `json:"hasRoleMappings,omitempty"`
	ID                     string   `json:"id,omitempty"`
	IdentityProviderID     string   `json:"identityProviderId,omitempty"`
	IdentityProviderStatus string   `json:"identityProviderStatus,omitempty"`
}

// Get gets Federated Settings for an organization.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/org-get-federation-settings/#std-label-atlas-org-get-federation-settings/

func (s *FederatedSettingsServiceOp) Get(ctx context.Context, orgID string) (*FederatedSettings, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf(federationSettingsBasePath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(FederatedSettings)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes federation setting.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/federation-delete-one/
func (s *FederatedSettingsServiceOp) Delete(ctx context.Context, federationSettingsID string) (*Response, error) {
	if federationSettingsID == "" {
		return nil, NewArgError("federationSettingsID", "must be set")
	}

	basePath := federationSettingsDeleteBasePath
	path := fmt.Sprintf("%s/%s", basePath, federationSettingsID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.Client.Do(ctx, req, nil)
}
