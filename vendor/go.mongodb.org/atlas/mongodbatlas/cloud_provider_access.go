// Copyright 2021 MongoDB Inc
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

const cloudProviderAccessPath = "api/atlas/v1.0/groups/%s/cloudProviderAccess"

// CloudProviderAccessService provides access to the cloud provider access functions in the Atlas API.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Cloud-Provider-Access
type CloudProviderAccessService interface {
	ListRoles(context.Context, string) (*CloudProviderAccessRoles, *Response, error)
	GetRole(context.Context, string, string) (*CloudProviderAccessRole, *Response, error)
	CreateRole(context.Context, string, *CloudProviderAccessRoleRequest) (*CloudProviderAccessRole, *Response, error)
	AuthorizeRole(context.Context, string, string, *CloudProviderAccessRoleRequest) (*CloudProviderAccessRole, *Response, error)
	DeauthorizeRole(context.Context, *CloudProviderDeauthorizationRequest) (*Response, error)
}

// CloudProviderAccessServiceOp provides an implementation of the CloudProviderAccessService interface.
type CloudProviderAccessServiceOp service

var _ CloudProviderAccessService = &CloudProviderAccessServiceOp{}

// CloudProviderAccessRoles an array of awsIamRoles objects.
type CloudProviderAccessRoles struct {
	AWSIAMRoles            []CloudProviderAccessRole `json:"awsIamRoles,omitempty"`            // Unique identifier of AWS security group in this access list entry.
	AzureServicePrincipals []CloudProviderAccessRole `json:"azureServicePrincipals,omitempty"` // Unique identifier of Azure security group in this access list entry.
}

// CloudProviderAccessRole is the response from the CloudProviderAccessService.ListRoles.
type CloudProviderAccessRole struct {
	AtlasAWSAccountARN         string          `json:"atlasAWSAccountArn,omitempty"`         // ARN associated with the Atlas AWS account used to assume IAM roles in your AWS account.
	AtlasAssumedRoleExternalID string          `json:"atlasAssumedRoleExternalId,omitempty"` // Unique external ID Atlas uses when assuming the IAM role in your AWS account.
	AuthorizedDate             string          `json:"authorizedDate,omitempty"`             //	Date on which this role was authorized.
	CreatedDate                string          `json:"createdDate,omitempty"`                // Date on which this role was created.
	FeatureUsages              []*FeatureUsage `json:"featureUsages,omitempty"`              // Atlas features this AWS IAM role is linked to.
	IAMAssumedRoleARN          string          `json:"iamAssumedRoleArn,omitempty"`          // ARN of the IAM Role that Atlas assumes when accessing resources in your AWS account.
	ProviderName               string          `json:"providerName,omitempty"`               // Name of the cloud provider. Currently limited to AWS.
	RoleID                     string          `json:"roleId,omitempty"`                     // Unique 24-hexadecimal digit string that identifies the role.
	AzureID                    *string         `json:"_id,omitempty"`                        // Unique 24-hexadecimal digit string that identifies the Azure Service Principal in Atlas.
	AtlasAzureAppID            *string         `json:"atlasAzureAppId,omitempty"`            // Azure Active Directory Application ID of Atlas.
	LastUpdatedDate            string          `json:"lastUpdatedDate,omitempty"`            // UUID string that identifies the Azure Service Principal.
	AzureServicePrincipalID    *string         `json:"servicePrincipalId,omitempty"`         // Unique ID of this role.
	AzureTenantID              *string         `json:"tenantId,omitempty"`                   // UUID String that identifies the Azure Active Directory Tenant ID.
}

// FeatureUsage represents where the role sis being used.
type FeatureUsage struct {
	FeatureType string      `json:"featureType,omitempty"`
	FeatureID   interface{} `json:"featureId,omitempty"`
}

// CloudProviderAccessRoleRequest represent a new role creation.
type CloudProviderAccessRoleRequest struct {
	ProviderName            string  `json:"providerName"`                 // Human-readable label that identifies the cloud provider of the role.
	IAMAssumedRoleARN       *string `json:"iamAssumedRoleArn,omitempty"`  // Amazon Resource Name (ARN) that identifies the Amazon Web Services (AWS) Identity and Access Management (IAM) role that MongoDB Cloud assumes when it accesses resources in your AWS account.
	AtlasAzureAppID         *string `json:"atlasAzureAppId,omitempty"`    // Date and time when this Azure Service Principal was last updated. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	AzureServicePrincipalID *string `json:"servicePrincipalId,omitempty"` // Unique AzureID of this role.
	AzureTenantID           *string `json:"tenantId,omitempty"`           // UUID String that identifies the Azure Active Directory Tenant AzureID.

}

// CloudProviderDeauthorizationRequest represents a request to remove authorization.
type CloudProviderDeauthorizationRequest struct {
	ProviderName string
	GroupID      string
	RoleID       string
}

// GetRole Returns the Amazon Web Services (AWS) Identity and Access Management (IAM) role
// with the specified id and with access to the specified project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Cloud-Provider-Access/operation/getCloudProviderAccessRole
func (s *CloudProviderAccessServiceOp) GetRole(ctx context.Context, groupID, roleID string) (*CloudProviderAccessRole, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if roleID == "" {
		return nil, nil, NewArgError("roleID", "must be set")
	}

	basePath := fmt.Sprintf(cloudProviderAccessPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, roleID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderAccessRole)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// ListRoles retrieves existing AWS IAM roles.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Cloud-Provider-Access/operation/listCloudProviderAccessRoles
func (s *CloudProviderAccessServiceOp) ListRoles(ctx context.Context, groupID string) (*CloudProviderAccessRoles, *Response, error) {
	path := fmt.Sprintf(cloudProviderAccessPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderAccessRoles)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// CreateRole creates an AWS IAM role.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Cloud-Provider-Access/operation/createCloudProviderAccessRole
func (s *CloudProviderAccessServiceOp) CreateRole(ctx context.Context, groupID string, request *CloudProviderAccessRoleRequest) (*CloudProviderAccessRole, *Response, error) {
	if request == nil {
		return nil, nil, NewArgError("request", "must be set")
	}

	path := fmt.Sprintf(cloudProviderAccessPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderAccessRole)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// AuthorizeRole authorizes and configure an AWS Assumed IAM role.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Cloud-Provider-Access/operation/authorizeCloudProviderAccessRole
func (s *CloudProviderAccessServiceOp) AuthorizeRole(ctx context.Context, groupID, roleID string, request *CloudProviderAccessRoleRequest) (*CloudProviderAccessRole, *Response, error) {
	if roleID == "" {
		return nil, nil, NewArgError("roleID", "must be set")
	}

	if request == nil {
		return nil, nil, NewArgError("request", "must be set")
	}

	basePath := fmt.Sprintf(cloudProviderAccessPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, roleID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, request)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderAccessRole)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// DeauthorizeRole deauthorizes an AWS Assumed IAM role.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Cloud-Provider-Access/operation/deauthorizeCloudProviderAccessRole
func (s *CloudProviderAccessServiceOp) DeauthorizeRole(ctx context.Context, request *CloudProviderDeauthorizationRequest) (*Response, error) {
	if request.RoleID == "" {
		return nil, NewArgError("roleID", "must be set")
	}

	basePath := fmt.Sprintf(cloudProviderAccessPath, request.GroupID)
	path := fmt.Sprintf("%s/%s/%s", basePath, request.ProviderName, request.RoleID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
