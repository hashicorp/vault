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
	"time"
)

const federationSettingsIdentityProviderBasePath = "api/atlas/v1.0/federationSettings/%s/identityProviders"

// A Resource describes a specific resource the Role will allow operating on.

// FederatedSettings represents a FederatedSettings List.
type FederatedSettingsIdentityProviders struct {
	Links      []*Link                             `json:"links,omitempty"`
	Results    []FederatedSettingsIdentityProvider `json:"results,omitempty"`
	TotalCount int                                 `json:"totalCount,omitempty"`
}

type FederatedSettingsIdentityProvider struct {
	AcsURL                     string            `json:"acsUrl,omitempty"`
	AssociatedDomains          []string          `json:"associatedDomains,omitempty"`
	AssociatedOrgs             []*AssociatedOrgs `json:"associatedOrgs,omitempty"`
	AudienceURI                string            `json:"audienceUri,omitempty"`
	DisplayName                string            `json:"displayName,omitempty"`
	IssuerURI                  string            `json:"issuerUri,omitempty"`
	OktaIdpID                  string            `json:"oktaIdpId,omitempty"`
	PemFileInfo                *PemFileInfo      `json:"pemFileInfo,omitempty"`
	RequestBinding             string            `json:"requestBinding,omitempty"`
	ResponseSignatureAlgorithm string            `json:"responseSignatureAlgorithm,omitempty"`
	SsoDebugEnabled            *bool             `json:"ssoDebugEnabled,omitempty"`
	SsoURL                     string            `json:"ssoUrl,omitempty"`
	Status                     string            `json:"status,omitempty"`
}

type AssociatedOrgs struct {
	DomainAllowList          []string        `json:"domainAllowList,omitempty"`
	DomainRestrictionEnabled *bool           `json:"domainRestrictionEnabled,omitempty"`
	IdentityProviderID       string          `json:"identityProviderId,omitempty"`
	OrgID                    string          `json:"orgId,omitempty"`
	PostAuthRoleGrants       []string        `json:"postAuthRoleGrants,omitempty"`
	RoleMappings             []*RoleMappings `json:"roleMappings,omitempty"`
	UserConflicts            *UserConflicts  `json:"userConflicts,omitempty"`
}
type PemFileInfo struct {
	Certificates []*Certificates `json:"certificates,omitempty"`
	FileName     string          `json:"fileName,omitempty"`
}
type Certificates struct {
	NotAfter  time.Time `json:"notAfter,omitempty"`
	NotBefore time.Time `json:"notBefore,omitempty"`
}

type UserConflicts []struct {
	EmailAddress         string `json:"emailAddress,omitempty"`
	FederationSettingsID string `json:"federationSettingsId,omitempty"`
	FirstName            string `json:"firstName,omitempty"`
	LastName             string `json:"lastName,omitempty"`
	UserID               string `json:"userId,omitempty"`
}

// ListIdentityProviders gets all Federated Settings Identity Providers for an organization.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/identity-provider-return-all/
func (s *FederatedSettingsServiceOp) ListIdentityProviders(ctx context.Context, federationSettingsID string, opts *ListOptions) ([]FederatedSettingsIdentityProvider, *Response, error) {
	if federationSettingsID == "" {
		return nil, nil, NewArgError("federationSettingsID", "must be set")
	}

	basePath := fmt.Sprintf(federationSettingsIdentityProviderBasePath, federationSettingsID)
	path, err := setListOptions(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(FederatedSettingsIdentityProviders)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// GetIdentityProvider gets Federated Settings Identity Providers for an organization.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/identity-provider-return-one/
func (s *FederatedSettingsServiceOp) GetIdentityProvider(ctx context.Context, federationSettingsID, idpID string) (*FederatedSettingsIdentityProvider, *Response, error) {
	if federationSettingsID == "" {
		return nil, nil, NewArgError("federationSettingsID", "must be set")
	}

	if idpID == "" {
		return nil, nil, NewArgError("idpID", "must be set")
	}

	basePath := fmt.Sprintf(federationSettingsIdentityProviderBasePath, federationSettingsID)
	path := fmt.Sprintf("%s/%s", basePath, idpID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(FederatedSettingsIdentityProvider)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// UpdateIdentityProvider updates Federated Settings Identity Providers for an organization.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/identity-provider-update-one/
func (s *FederatedSettingsServiceOp) UpdateIdentityProvider(ctx context.Context, federationSettingsID, idpID string, updateRequest *FederatedSettingsIdentityProvider) (*FederatedSettingsIdentityProvider, *Response, error) {
	if federationSettingsID == "" {
		return nil, nil, NewArgError("federationSettingsID", "must be set")
	}

	if idpID == "" {
		return nil, nil, NewArgError("idpID", "must be set")
	}

	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(federationSettingsIdentityProviderBasePath, federationSettingsID)
	path := fmt.Sprintf("%s/%s", basePath, idpID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(FederatedSettingsIdentityProvider)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
