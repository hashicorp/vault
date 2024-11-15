// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azureauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathsRole returns the path configurations for the CRUD operations on roles
func pathsRole(b *azureAuthBackend) []*framework.Path {
	p := []*framework.Path{
		{
			Pattern: "role/?",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixAzure,
				OperationSuffix: "auth-roles",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.pathRoleList,
				},
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-list"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-list"][1]),
		},
		{
			Pattern: "role/" + framework.GenericNameRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixAzure,
				OperationSuffix: "auth-role",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"policies": {
					Type:        framework.TypeCommaStringSlice,
					Description: tokenutil.DeprecationText("token_policies"),
					Deprecated:  true,
				},
				"num_uses": {
					Type:        framework.TypeInt,
					Description: tokenutil.DeprecationText("token_num_uses"),
					Deprecated:  true,
				},
				"ttl": {
					Type:        framework.TypeDurationSecond,
					Description: tokenutil.DeprecationText("token_ttl"),
					Deprecated:  true,
				},
				"max_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: tokenutil.DeprecationText("token_max_ttl"),
					Deprecated:  true,
				},
				"period": {
					Type:        framework.TypeDurationSecond,
					Description: tokenutil.DeprecationText("token_period"),
					Deprecated:  true,
				},
				"bound_subscription_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: `Comma-separated list of subscription ids that login is restricted to.`,
				},
				"bound_resource_groups": {
					Type:        framework.TypeCommaStringSlice,
					Description: `Comma-separated list of resource groups that login is restricted to.`,
				},
				"bound_group_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: `Comma-separated list of group ids that login is restricted to.`,
				},
				"bound_service_principal_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: `Comma-separated list of service principal ids that login is restricted to.`,
				},
				"bound_locations": {
					Type:        framework.TypeCommaStringSlice,
					Description: `Comma-separated list of locations that login is restricted to.`,
				},
				"bound_scale_sets": {
					Type:        framework.TypeCommaStringSlice,
					Description: `Comma-separated list of scale sets that login is restricted to.`,
				},
			},
			ExistenceCheck: b.pathRoleExistenceCheck,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathRoleCreateUpdate,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathRoleCreateUpdate,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathRoleRead,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.pathRoleDelete,
				},
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role"][1]),
		},
	}

	tokenutil.AddTokenFields(p[1].Fields)

	return p
}

type azureRole struct {
	tokenutil.TokenParams

	// Policies that are to be required by the token to access this role
	Policies []string `json:"policies"`

	// TokenNumUses defines the number of allowed uses of the token issued
	NumUses int `json:"num_uses"`

	// Duration before which an issued token must be renewed
	TTL time.Duration `json:"ttl"`

	// Duration after which an issued token should not be allowed to be renewed
	MaxTTL time.Duration `json:"max_ttl"`

	// Period, if set, indicates that the token generated using this role
	// should never expire. The token should be renewed within the duration
	// specified by this value. The renewal duration will be fixed if the
	// value is not modified on the role. If the `Period` in the role is modified,
	// a token will pick up the new value during its next renewal.
	Period time.Duration `json:"period"`

	// Role binding properties
	BoundServicePrincipalIDs []string `json:"bound_service_principal_ids"`
	BoundGroupIDs            []string `json:"bound_group_ids"`
	BoundResourceGroups      []string `json:"bound_resource_groups"`
	BoundSubscriptionsIDs    []string `json:"bound_subscription_ids"`
	BoundLocations           []string `json:"bound_locations"`
	BoundScaleSets           []string `json:"bound_scale_sets"`
}

// role takes a storage backend and the name and returns the role's storage
// entryÍ
func (b *azureAuthBackend) role(ctx context.Context, s logical.Storage, name string) (*azureRole, error) {
	raw, err := s.Get(ctx, "role/"+strings.ToLower(name))
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	role := new(azureRole)
	if err := json.Unmarshal(raw.Value, role); err != nil {
		return nil, err
	}

	if role.TokenTTL == 0 && role.TTL > 0 {
		role.TokenTTL = role.TTL
	}
	if role.TokenMaxTTL == 0 && role.MaxTTL > 0 {
		role.TokenMaxTTL = role.MaxTTL
	}
	if role.TokenPeriod == 0 && role.Period > 0 {
		role.TokenPeriod = role.Period
	}
	if role.TokenNumUses == 0 && role.NumUses > 0 {
		role.TokenNumUses = role.NumUses
	}
	if len(role.TokenPolicies) == 0 && len(role.Policies) > 0 {
		role.TokenPolicies = role.Policies
	}

	return role, nil
}

// pathRoleExistenceCheck returns whether the role with the given name exists or not.
func (b *azureAuthBackend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	role, err := b.role(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return false, err
	}
	return role != nil, nil
}

// pathRoleList is used to list all the Roles registered with the backend.
func (b *azureAuthBackend) pathRoleList(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

// pathRoleRead grabs a read lock and reads the options set on the role from the storage
func (b *azureAuthBackend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	d := map[string]interface{}{
		"bound_service_principal_ids": role.BoundServicePrincipalIDs,
		"bound_group_ids":             role.BoundGroupIDs,
		"bound_subscription_ids":      role.BoundSubscriptionsIDs,
		"bound_resource_groups":       role.BoundResourceGroups,
		"bound_locations":             role.BoundLocations,
		"bound_scale_sets":            role.BoundScaleSets,
	}

	role.PopulateTokenData(d)

	if len(role.Policies) > 0 {
		d["policies"] = d["token_policies"]
	}
	if role.TTL > 0 {
		d["ttl"] = int64(role.TTL.Seconds())
	}
	if role.MaxTTL > 0 {
		d["max_ttl"] = int64(role.MaxTTL.Seconds())
	}
	if role.Period > 0 {
		d["period"] = int64(role.Period.Seconds())
	}
	if role.NumUses > 0 {
		d["num_uses"] = role.NumUses
	}

	return &logical.Response{
		Data: d,
	}, nil
}

// pathRoleDelete removes the role from storage
func (b *azureAuthBackend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("role name required"), nil
	}

	// Delete the role itself
	if err := req.Storage.Delete(ctx, "role/"+strings.ToLower(roleName)); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathRoleCreateUpdate registers a new role with the backend or updates the options
// of an existing role
func (b *azureAuthBackend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	// Check if the role already exists
	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	// Create a new entry object if this is a CreateOperation
	if role == nil {
		if req.Operation == logical.UpdateOperation {
			return nil, errors.New("role entry not found during update operation")
		}
		role = new(azureRole)
	}

	if err := role.ParseTokenFields(req, data); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Handle upgrade cases
	{
		if err := tokenutil.UpgradeValue(data, "policies", "token_policies", &role.Policies, &role.TokenPolicies); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "num_uses", "token_num_uses", &role.NumUses, &role.TokenNumUses); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "ttl", "token_ttl", &role.TTL, &role.TokenTTL); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "max_ttl", "token_max_ttl", &role.MaxTTL, &role.TokenMaxTTL); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "period", "token_period", &role.Period, &role.TokenPeriod); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	if role.TokenPeriod > b.System().MaxLeaseTTL() {
		return logical.ErrorResponse("token period of %q is greater than the backend's maximum lease TTL of %q", role.TokenPeriod.String(), b.System().MaxLeaseTTL().String()), nil
	}

	if role.TokenNumUses < 0 {
		return logical.ErrorResponse("token num uses cannot be negative"), nil
	}

	if boundServicePrincipalIDs, ok := data.GetOk("bound_service_principal_ids"); ok {
		role.BoundServicePrincipalIDs = boundServicePrincipalIDs.([]string)
	}

	if boundGroupIDs, ok := data.GetOk("bound_group_ids"); ok {
		role.BoundGroupIDs = boundGroupIDs.([]string)
	}

	if boundSubscriptionsIDs, ok := data.GetOk("bound_subscription_ids"); ok {
		role.BoundSubscriptionsIDs = boundSubscriptionsIDs.([]string)
	}

	if boundResourceGroups, ok := data.GetOk("bound_resource_groups"); ok {
		role.BoundResourceGroups = boundResourceGroups.([]string)
	}

	if boundLocations, ok := data.GetOk("bound_locations"); ok {
		role.BoundLocations = boundLocations.([]string)
	}

	if boundScaleSets, ok := data.GetOk("bound_scale_sets"); ok {
		role.BoundScaleSets = boundScaleSets.([]string)
	}

	if len(role.BoundServicePrincipalIDs) == 0 &&
		len(role.BoundGroupIDs) == 0 &&
		len(role.BoundSubscriptionsIDs) == 0 &&
		len(role.BoundResourceGroups) == 0 &&
		len(role.BoundLocations) == 0 &&
		len(role.BoundScaleSets) == 0 {
		return logical.ErrorResponse("must have at least one bound constraint when creating/updating a role"), nil
	}

	// Check that the TTL value provided is less than the MaxTTL.
	// Sanitizing the TTL and MaxTTL is not required now and can be performed
	// at credential issue time.
	if role.TokenMaxTTL > 0 && role.TokenTTL > role.TokenMaxTTL {
		return logical.ErrorResponse("ttl should not be greater than max_ttl"), nil
	}

	var resp *logical.Response
	if role.TokenMaxTTL > b.System().MaxLeaseTTL() {
		resp = &logical.Response{}
		resp.AddWarning("max_ttl is greater than the system or backend mount's maximum TTL value; issued tokens' max TTL value will be truncated")
	}

	// Store the entry.
	entry, err := logical.StorageEntryJSON("role/"+strings.ToLower(roleName), role)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("failed to create storage entry for role %s", roleName)
	}
	if err = req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return resp, nil
}

// roleStorageEntry stores all the options that are set on an role
var roleHelp = map[string][2]string{
	"role-list": {
		"Lists all the roles registered with the backend.",
		"The list will contain the names of the roles.",
	},
	"role": {
		"Register an role with the backend.",
		`A role is required to authenticate with this backend. The role binds
		Azure instance metadata with token policies and settings.
		The bindings, token polices and token settings can all be configured
		using this endpoint`,
	},
}
