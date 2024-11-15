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

const federationSettingsOrganizationRoleMappingBasePath = "api/atlas/v1.0/federationSettings/%s/connectedOrgConfigs/%s/roleMappings"

// A Resource describes a specific resource the Role will allow operating on.

// FederatedSettings represents a FederatedSettings Organization Connection..
type FederatedSettingsOrganizationRoleMappings struct {
	Links      []*Link                                     `json:"links,omitempty"`
	Results    []*FederatedSettingsOrganizationRoleMapping `json:"results,omitempty"`
	TotalCount int                                         `json:"totalCount,omitempty"`
}

type FederatedSettingsOrganizationRoleMapping struct {
	ExternalGroupName string             `json:"externalGroupName,omitempty"`
	ID                string             `json:"id,omitempty"`
	RoleAssignments   []*RoleAssignments `json:"roleAssignments,omitempty"`
}

type RoleAssignments struct {
	GroupID string `json:"groupId,omitempty"`
	OrgID   string `json:"orgId,omitempty"`
	Role    string `json:"role,omitempty"`
}

// ListRoleMappings gets all Federated Settings Role Mappings for an organization.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/role-mapping-return-all/
func (s *FederatedSettingsServiceOp) ListRoleMappings(ctx context.Context, federationSettingsID, orgID string, opts *ListOptions) (*FederatedSettingsOrganizationRoleMappings, *Response, error) {
	if federationSettingsID == "" {
		return nil, nil, NewArgError("federationSettingsID", "must be set")
	}

	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf(federationSettingsOrganizationRoleMappingBasePath, federationSettingsID, orgID)
	path, err := setListOptions(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(FederatedSettingsOrganizationRoleMappings)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// GetRoleMapping gets Federated Settings Role Mapping for an organization.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/role-mapping-return-one/
func (s *FederatedSettingsServiceOp) GetRoleMapping(ctx context.Context, federationSettingsID, orgID, roleMappingID string) (*FederatedSettingsOrganizationRoleMapping, *Response, error) {
	if federationSettingsID == "" {
		return nil, nil, NewArgError("federationSettingsID", "must be set")
	}

	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf(federationSettingsOrganizationRoleMappingBasePath, federationSettingsID, orgID)
	path := fmt.Sprintf("%s/%s", basePath, roleMappingID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(FederatedSettingsOrganizationRoleMapping)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// CreateRoleMapping creates one new Federated Settings Role Mapping for an organization.
//
// See more: https://docs.atlas.mongodb.com/reference/api/live-migration/create-one-migration/
func (s *FederatedSettingsServiceOp) CreateRoleMapping(ctx context.Context, federationSettingsID, orgID string, body *FederatedSettingsOrganizationRoleMapping) (*FederatedSettingsOrganizationRoleMapping, *Response, error) {
	if federationSettingsID == "" {
		return nil, nil, NewArgError("federationSettingsID", "must be set")
	}

	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf(federationSettingsOrganizationRoleMappingBasePath, federationSettingsID, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, nil, err
	}

	root := new(FederatedSettingsOrganizationRoleMapping)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// UpdateRoleMapping updates Federated Settings Federated Settings Role Mapping for an organization
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/role-mapping-create-one/
func (s *FederatedSettingsServiceOp) UpdateRoleMapping(ctx context.Context, federationSettingsID, orgID, roleMappingID string, updateRequest *FederatedSettingsOrganizationRoleMapping) (*FederatedSettingsOrganizationRoleMapping, *Response, error) {
	if federationSettingsID == "" {
		return nil, nil, NewArgError("federationSettingsID", "must be set")
	}

	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(federationSettingsOrganizationRoleMappingBasePath, federationSettingsID, orgID)
	path := fmt.Sprintf("%s/%s", basePath, roleMappingID)

	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(FederatedSettingsOrganizationRoleMapping)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// DeleteRoleMapping deletes Federated Settings Role Mapping for an organization.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/role-mapping-delete-one/
func (s *FederatedSettingsServiceOp) DeleteRoleMapping(ctx context.Context, federationSettingsID, orgID, roleMappingID string) (*Response, error) {
	if federationSettingsID == "" {
		return nil, NewArgError("federationSettingsID", "must be set")
	}

	if orgID == "" {
		return nil, NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf(federationSettingsOrganizationRoleMappingBasePath, federationSettingsID, orgID)
	path := fmt.Sprintf("%s/%s", basePath, roleMappingID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
