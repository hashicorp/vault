// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package ibmcloud

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	roleStoragePrefix = "role/"

	// Role types
	roleTypeAPIKey         = "apikey"
	roleTypeTrustedProfile = "trusted_profile"
)

// ibmCloudRole represents a role configuration
type ibmCloudRole struct {
	tokenutil.TokenParams

	// Type of authentication: "apikey" or "trusted_profile"
	Type string `json:"type"`

	// API Key type fields
	BoundAccountID      string   `json:"bound_account_id"`
	BoundServiceIDs     []string `json:"bound_service_ids"`
	BoundUserIDs        []string `json:"bound_user_ids"`
	BoundAccessGroupIDs []string `json:"bound_access_group_ids"`

	// Trusted Profile type fields
	BoundTrustedProfileID    string   `json:"bound_trusted_profile_id"`
	BoundComputeResourceCRN  []string `json:"bound_compute_resource_crn"`
	BoundComputeResourceType []string `json:"bound_compute_resource_type"`
}

func (b *backend) pathRole() *framework.Path {
	p := &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("name"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixIBMCloud,
			OperationSuffix: "role",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the role",
			},
			"type": {
				Type:        framework.TypeString,
				Description: "Type of authentication: 'apikey' or 'trusted_profile'",
			},
			"bound_account_id": {
				Type:        framework.TypeString,
				Description: "(API Key type) IBM Cloud account ID that the authenticating entity must belong to",
			},
			"bound_service_ids": {
				Type:        framework.TypeCommaStringSlice,
				Description: "(API Key type) List of service IDs that are allowed to authenticate",
			},
			"bound_user_ids": {
				Type:        framework.TypeCommaStringSlice,
				Description: "(API Key type) List of user IDs that are allowed to authenticate",
			},
			"bound_access_group_ids": {
				Type:        framework.TypeCommaStringSlice,
				Description: "(API Key type) List of access group IDs that the authenticating entity must belong to",
			},
			"bound_trusted_profile_id": {
				Type:        framework.TypeString,
				Description: "(Trusted Profile type) Trusted profile ID that must be used for authentication",
			},
			"bound_compute_resource_crn": {
				Type:        framework.TypeCommaStringSlice,
				Description: "(Trusted Profile type) CRN patterns for compute resources (supports wildcards)",
			},
			"bound_compute_resource_type": {
				Type:        framework.TypeCommaStringSlice,
				Description: "(Trusted Profile type) Resource types (e.g., 'vsi', 'iks')",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathRoleRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "read",
				},
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathRoleCreateUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "create",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathRoleCreateUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "update",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathRoleDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "delete",
				},
			},
		},

		ExistenceCheck: b.pathRoleExistenceCheck,

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}

	tokenutil.AddTokenFields(p.Fields)
	return p
}

// pathRoleExistenceCheck checks if a role exists
func (b *backend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	name := data.Get("name").(string)
	if name == "" {
		return false, nil
	}

	role, err := b.role(ctx, req.Storage, name)
	if err != nil {
		return false, err
	}

	return role != nil, nil
}

func (b *backend) pathListRole() *framework.Path {
	return &framework.Path{
		Pattern: "role/?$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixIBMCloud,
			OperationSuffix: "roles",
			Navigation:      true,
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathRoleList,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "list",
				},
			},
		},

		HelpSynopsis:    pathRoleListHelpSyn,
		HelpDescription: pathRoleListHelpDesc,
	}
}

func (b *backend) pathListRoles() *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixIBMCloud,
			OperationSuffix: "roles",
			Navigation:      true,
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathRoleList,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "list",
				},
			},
		},

		HelpSynopsis:    pathRoleListHelpSyn,
		HelpDescription: pathRoleListHelpDesc,
	}
}

func (b *backend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, nil
	}

	respData := map[string]interface{}{
		"type":                        role.Type,
		"bound_account_id":            role.BoundAccountID,
		"bound_service_ids":           role.BoundServiceIDs,
		"bound_user_ids":              role.BoundUserIDs,
		"bound_access_group_ids":      role.BoundAccessGroupIDs,
		"bound_trusted_profile_id":    role.BoundTrustedProfileID,
		"bound_compute_resource_crn":  role.BoundComputeResourceCRN,
		"bound_compute_resource_type": role.BoundComputeResourceType,
	}

	role.PopulateTokenData(respData)

	return &logical.Response{
		Data: respData,
	}, nil
}

func (b *backend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	if role == nil {
		role = &ibmCloudRole{}
	}

	// Update type
	if roleType, ok := data.GetOk("type"); ok {
		role.Type = roleType.(string)
	}

	// Validate type
	if role.Type != roleTypeAPIKey && role.Type != roleTypeTrustedProfile {
		return logical.ErrorResponse("type must be 'apikey' or 'trusted_profile'"), nil
	}

	// Update API Key type fields
	if boundAccountID, ok := data.GetOk("bound_account_id"); ok {
		role.BoundAccountID = boundAccountID.(string)
	}

	if boundServiceIDs, ok := data.GetOk("bound_service_ids"); ok {
		role.BoundServiceIDs = boundServiceIDs.([]string)
	}

	if boundUserIDs, ok := data.GetOk("bound_user_ids"); ok {
		role.BoundUserIDs = boundUserIDs.([]string)
	}

	if boundAccessGroupIDs, ok := data.GetOk("bound_access_group_ids"); ok {
		role.BoundAccessGroupIDs = boundAccessGroupIDs.([]string)
	}

	// Update Trusted Profile type fields
	if boundTrustedProfileID, ok := data.GetOk("bound_trusted_profile_id"); ok {
		role.BoundTrustedProfileID = boundTrustedProfileID.(string)
	}

	if boundComputeResourceCRN, ok := data.GetOk("bound_compute_resource_crn"); ok {
		role.BoundComputeResourceCRN = boundComputeResourceCRN.([]string)
	}

	if boundComputeResourceType, ok := data.GetOk("bound_compute_resource_type"); ok {
		role.BoundComputeResourceType = boundComputeResourceType.([]string)
	}

	// Validate role configuration
	if err := role.validate(); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Update token parameters
	if err := role.ParseTokenFields(req, data); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Store the role
	entry, err := logical.StorageEntryJSON(roleStoragePrefix+roleName, role)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	if err := req.Storage.Delete(ctx, roleStoragePrefix+roleName); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List(ctx, roleStoragePrefix)
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(roles), nil
}

// role reads a role from storage
func (b *backend) role(ctx context.Context, s logical.Storage, name string) (*ibmCloudRole, error) {
	entry, err := s.Get(ctx, roleStoragePrefix+name)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	role := &ibmCloudRole{}
	if err := entry.DecodeJSON(role); err != nil {
		return nil, fmt.Errorf("failed to decode role: %w", err)
	}

	return role, nil
}

// validate checks if the role configuration is valid
func (r *ibmCloudRole) validate() error {
	if r.Type == "" {
		return fmt.Errorf("type is required")
	}

	switch r.Type {
	case roleTypeAPIKey:
		if r.BoundAccountID == "" {
			return fmt.Errorf("bound_account_id is required for apikey type")
		}
	case roleTypeTrustedProfile:
		if r.BoundTrustedProfileID == "" {
			return fmt.Errorf("bound_trusted_profile_id is required for trusted_profile type")
		}
	default:
		return fmt.Errorf("invalid type: must be 'apikey' or 'trusted_profile'")
	}

	return nil
}

// matchesIdentity checks if the identity matches the role bindings
func (r *ibmCloudRole) matchesIdentity(identity *IdentityInfo) error {
	switch r.Type {
	case roleTypeAPIKey:
		return r.matchesAPIKeyIdentity(identity)
	case roleTypeTrustedProfile:
		return r.matchesTrustedProfileIdentity(identity)
	default:
		return fmt.Errorf("invalid role type: %s", r.Type)
	}
}

func (r *ibmCloudRole) matchesAPIKeyIdentity(identity *IdentityInfo) error {
	// Check account ID
	if identity.AccountID != r.BoundAccountID {
		return fmt.Errorf("account ID %s does not match bound account ID %s", identity.AccountID, r.BoundAccountID)
	}

	// Check service IDs if specified
	if len(r.BoundServiceIDs) > 0 {
		matched := false
		for _, serviceID := range r.BoundServiceIDs {
			if identity.IAMID == serviceID || identity.Subject == serviceID {
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("service ID %s not in allowed list", identity.IAMID)
		}
	}

	// Check user IDs if specified
	if len(r.BoundUserIDs) > 0 {
		matched := false
		for _, userID := range r.BoundUserIDs {
			if identity.IAMID == userID || identity.Subject == userID {
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("user ID %s not in allowed list", identity.IAMID)
		}
	}

	// Note: Access group membership is verified separately via API call in path_login.go
	// because it's not included in the token claims and requires a separate IAM API call.
	// See pathLoginUpdate for the actual access group verification.

	return nil
}

func (r *ibmCloudRole) matchesTrustedProfileIdentity(identity *IdentityInfo) error {
	// Check trusted profile ID
	if identity.TrustedProfileID != r.BoundTrustedProfileID {
		return fmt.Errorf("trusted profile ID %s does not match bound profile ID %s", identity.TrustedProfileID, r.BoundTrustedProfileID)
	}

	// Check compute resource CRN patterns if specified
	if len(r.BoundComputeResourceCRN) > 0 {
		if identity.CRN == "" {
			return fmt.Errorf("no CRN found in identity")
		}
		matched := false
		for _, pattern := range r.BoundComputeResourceCRN {
			if matchCRNPattern(identity.CRN, pattern) {
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("CRN %s does not match any bound patterns", identity.CRN)
		}
	}

	// Check compute resource type if specified
	if len(r.BoundComputeResourceType) > 0 {
		if identity.ComputeResourceID == "" {
			return fmt.Errorf("no compute resource ID found in identity")
		}
		// Extract resource type from CRN or compute resource ID
		resourceType := extractResourceType(identity.CRN, identity.ComputeResourceID)
		matched := false
		for _, allowedType := range r.BoundComputeResourceType {
			if resourceType == allowedType {
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("resource type %s not in allowed list", resourceType)
		}
	}

	return nil
}

// matchCRNPattern checks if a CRN matches a pattern (supports wildcards)
func matchCRNPattern(crn, pattern string) bool {
	if pattern == "*" {
		return true
	}

	// Split pattern by "*" and escape all non-wildcard parts to be regex-safe
	parts := strings.Split(pattern, "*")
	for i, part := range parts {
		parts[i] = regexp.QuoteMeta(part)
	}

	// Rejoin segments with ".*" and enforce exact start and end matching
	regexPattern := "^" + strings.Join(parts, ".*") + "$"
	matched, err := regexp.MatchString(regexPattern, crn)
	if err != nil {
		return false
	}
	return matched
}

// extractResourceType extracts the resource type from CRN or compute resource ID
func extractResourceType(crn, computeResourceID string) string {
	// CRN format: crn:version:cname:ctype:service-name:location:scope:service-instance:resource-type:resource
	// Example: crn:v1:bluemix:public:is:us-south:a/abc123::instance:xxx
	parts := strings.Split(crn, ":")
	if len(parts) >= 9 {
		return parts[8] // resource-type
	}
	return ""
}

const pathRoleHelpSyn = `
Create or update a role for IBM Cloud authentication.
`

const pathRoleHelpDesc = `
A role defines the authentication requirements and token properties for
IBM Cloud authentication.

Role Types:

  * apikey - Authenticate using IBM Cloud API keys
  * trusted_profile - Authenticate using IBM Cloud Trusted Profiles

API Key Type Parameters:

  * bound_account_id (required) - IBM Cloud account ID
  * bound_service_ids (optional) - List of allowed service IDs
  * bound_user_ids (optional) - List of allowed user IDs
  * bound_access_group_ids (optional) - List of required access group IDs

Trusted Profile Type Parameters:

  * bound_trusted_profile_id (required) - Trusted profile ID
  * bound_compute_resource_crn (optional) - CRN patterns (supports wildcards)
  * bound_compute_resource_type (optional) - Resource types (e.g., 'vsi', 'iks')

Token Parameters:

  All standard Vault token parameters are supported (token_ttl, token_max_ttl,
  token_policies, etc.)

Example:

  $ vault write auth/ibmcloud/role/my-role \
      type=apikey \
      bound_account_id="abc123" \
      bound_service_ids="ServiceId-xxx" \
      token_policies="default,my-policy" \
      token_ttl=1h
`

const pathRoleListHelpSyn = `
List all roles configured for IBM Cloud authentication.
`

const pathRoleListHelpDesc = `
This endpoint lists all roles that have been configured for IBM Cloud
authentication.
`

// Made with Bob
