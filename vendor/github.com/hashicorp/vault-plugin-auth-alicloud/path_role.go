// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package alicloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathRole(b *backend) *framework.Path {
	p := &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("role"),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAliCloud,
			OperationSuffix: "auth-role",
		},
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeLowerCaseString,
				Description: "The name of the role as it should appear in Vault.",
			},
			"arn": {
				Type:        framework.TypeString,
				Description: "ARN of the RAM to bind to this role.",
			},
			"policies": {
				Type:        framework.TypeCommaStringSlice,
				Description: tokenutil.DeprecationText("token_policies"),
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
			"bound_cidrs": {
				Type:        framework.TypeCommaStringSlice,
				Description: tokenutil.DeprecationText("token_bound_cidrs"),
				Deprecated:  true,
			},
		},
		ExistenceCheck: b.operationRoleExistenceCheck,
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.operationRoleCreateUpdate,
			logical.UpdateOperation: b.operationRoleCreateUpdate,
			logical.ReadOperation:   b.operationRoleRead,
			logical.DeleteOperation: b.operationRoleDelete,
		},
		HelpSynopsis:    pathRoleSyn,
		HelpDescription: pathRoleDesc,
	}

	tokenutil.AddTokenFields(p.Fields)
	return p
}

func pathListRole(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/?",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAliCloud,
			OperationVerb:   "list",
			OperationSuffix: "auth-roles",
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.operationRoleList,
		},
		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAliCloud,
			OperationVerb:   "list",
			OperationSuffix: "auth-roles2",
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.operationRoleList,
		},
		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) operationRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := readRole(ctx, req.Storage, data.Get("role").(string))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *backend) operationRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	role, err := readRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil && req.Operation == logical.UpdateOperation {
		return nil, fmt.Errorf("no role found to update for %s", roleName)
	} else if role == nil {
		role = &roleEntry{}
	}

	if raw, ok := data.GetOk("arn"); ok {
		arn, err := parseARN(raw.(string))
		if err != nil {
			return nil, fmt.Errorf("unable to parse arn %s: %w", arn, err)
		}
		if arn.Type != arnTypeRole {
			return nil, fmt.Errorf("only role arn types are supported at this time, but %s was provided", role.ARN)
		}
		role.ARN = arn
	} else if req.Operation == logical.CreateOperation {
		return nil, errors.New("the arn is required to create a role")
	}

	// Get tokenutil fields
	if err := role.ParseTokenFields(req, data); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Handle upgrade cases
	{
		if err := tokenutil.UpgradeValue(data, "policies", "token_policies", &role.Policies, &role.TokenPolicies); err != nil {
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

		if err := tokenutil.UpgradeValue(data, "bound_cidrs", "token_bound_cidrs", &role.BoundCIDRs, &role.TokenBoundCIDRs); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	if role.TokenMaxTTL > 0 && role.TokenTTL > role.TokenMaxTTL {
		return nil, errors.New("ttl exceeds max ttl")
	}

	entry, err := logical.StorageEntryJSON("role/"+roleName, role)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	if role.TTL > b.System().MaxLeaseTTL() {
		resp := &logical.Response{}
		resp.AddWarning(fmt.Sprintf("ttl of %d exceeds the system max ttl of %d, the latter will be used during login", role.TTL, b.System().MaxLeaseTTL()))
		return resp, nil
	}
	return nil, nil
}

func (b *backend) operationRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	role, err := readRole(ctx, req.Storage, data.Get("role").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}
	return &logical.Response{
		Data: role.ToResponseData(),
	}, nil
}

func (b *backend) operationRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, "role/"+data.Get("role").(string)); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) operationRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleNames, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roleNames), nil
}

func readRole(ctx context.Context, s logical.Storage, roleName string) (*roleEntry, error) {
	role, err := s.Get(ctx, "role/"+roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}
	result := &roleEntry{}
	if err := role.DecodeJSON(result); err != nil {
		return nil, err
	}

	if result.TokenTTL == 0 && result.TTL > 0 {
		result.TokenTTL = result.TTL
	}
	if result.TokenMaxTTL == 0 && result.MaxTTL > 0 {
		result.TokenMaxTTL = result.MaxTTL
	}
	if result.TokenPeriod == 0 && result.Period > 0 {
		result.TokenPeriod = result.Period
	}
	if len(result.TokenPolicies) == 0 && len(result.Policies) > 0 {
		result.TokenPolicies = result.Policies
	}
	if len(result.TokenBoundCIDRs) == 0 && len(result.BoundCIDRs) > 0 {
		result.TokenBoundCIDRs = result.BoundCIDRs
	}

	return result, nil
}

const pathRoleSyn = `
Create a role and associate policies to it.
`

const pathRoleDesc = `
A precondition for login is that a role should be created in the backend.
The login endpoint takes in the role name against which the instance
should be validated. After authenticating the instance, the authorization
for the instance to access Vault's resources is determined by the policies
that are associated to the role though this endpoint.

Also, a 'max_ttl' can be configured in this endpoint that determines the maximum
duration for which a login can be renewed. Note that the 'max_ttl' has an upper
limit of the 'max_ttl' value on the backend's mount. The same applies to the 'ttl'.
`

const pathListRolesHelpSyn = `
Lists all the roles that are registered with Vault.
`

const pathListRolesHelpDesc = `
Roles will be listed by their respective role names.
`
