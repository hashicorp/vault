// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cf

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault-plugin-auth-cf/models"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const roleStoragePrefix = "roles/"

func (b *backend) pathListRoles() *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixCloudFoundry,
			OperationVerb:   "list",
			OperationSuffix: "roles",
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.operationRolesList,
			},
		},
		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

func (b *backend) operationRolesList(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, roleStoragePrefix)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(entries), nil
}

func (b *backend) pathRoles() *framework.Path {
	p := &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("role"),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixCloudFoundry,
			OperationSuffix: "role",
		},
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeLowerCaseString,
				Required:    true,
				Description: "The name of the role.",
			},
			"bound_application_ids": {
				Type: framework.TypeCommaStringSlice,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Bound Application IDs",
					Value: "6b814521-5f08-4b1a-8c4e-fbe7c5f3a169",
				},
				Description: "Require that the client certificate presented has at least one of these app IDs.",
			},
			"bound_space_ids": {
				Type: framework.TypeCommaStringSlice,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Bound Space IDs",
					Value: "3d2eba6b-ef19-44d5-91dd-1975b0db5cc9",
				},
				Description: "Require that the client certificate presented has at least one of these space IDs.",
			},
			"bound_organization_ids": {
				Type: framework.TypeCommaStringSlice,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Bound Organization IDs",
					Value: "34a878d0-c2f9-4521-ba73-a9f664e82c7b",
				},
				Description: "Require that the client certificate presented has at least one of these org IDs.",
			},
			"bound_instance_ids": {
				Type: framework.TypeCommaStringSlice,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Bound Instance IDs",
					Value: "8a886b31-ccf7-480d-54d8-cc28",
				},
				Description: "Require that the client certificate presented has at least one of these instance IDs.",
			},
			"disable_ip_matching": {
				Type:    framework.TypeBool,
				Default: false,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Disable IP Address Matching",
					Value: "false",
				},
				Description: `If set to true, disables the default behavior that logging in must be performed from 
an acceptable IP address described by the certificate presented.`,
			},
			"policies": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: tokenutil.DeprecationText("token_policies"),
				Deprecated:  true,
			},
			"ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: tokenutil.DeprecationText("token_ttl"),
				Deprecated:  true,
			},
			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: tokenutil.DeprecationText("token_max_ttl"),
				Deprecated:  true,
			},
			"period": &framework.FieldSchema{
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
		ExistenceCheck: b.operationRolesExistenceCheck,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.operationRolesCreateUpdate,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.operationRolesCreateUpdate,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.operationRolesRead,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.operationRolesDelete,
			},
		},
		HelpSynopsis:    pathRolesHelpSyn,
		HelpDescription: pathRolesHelpDesc,
	}

	tokenutil.AddTokenFields(p.Fields)
	return p
}

func (b *backend) operationRolesExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := req.Storage.Get(ctx, roleStoragePrefix+data.Get("role").(string))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *backend) operationRolesCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	role := &models.RoleEntry{}
	if req.Operation == logical.UpdateOperation {
		storedRole, err := getRole(ctx, req.Storage, roleName)
		if err != nil {
			return nil, err
		}
		if storedRole != nil {
			role = storedRole
		}
	}
	if raw, ok := data.GetOk("bound_application_ids"); ok {
		role.BoundAppIDs = raw.([]string)
	}
	if raw, ok := data.GetOk("bound_space_ids"); ok {
		role.BoundSpaceIDs = raw.([]string)
	}
	if raw, ok := data.GetOk("bound_organization_ids"); ok {
		role.BoundOrgIDs = raw.([]string)
	}
	if raw, ok := data.GetOk("bound_instance_ids"); ok {
		role.BoundInstanceIDs = raw.([]string)
	}
	if raw, ok := data.GetOk("disable_ip_matching"); ok {
		role.DisableIPMatching = raw.(bool)
	}

	if err := role.ParseTokenFields(req, data); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Handle upgrade cases
	{
		if err := tokenutil.UpgradeValue(data, "policies", "token_policies", &role.Policies, &role.TokenPolicies); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "bound_cidrs", "token_bound_cidrs", &role.BoundCIDRs, &role.TokenBoundCIDRs); err != nil {
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

	if role.TokenMaxTTL > 0 && role.TokenTTL > role.TokenMaxTTL {
		return logical.ErrorResponse("ttl exceeds max ttl"), nil
	}

	entry, err := logical.StorageEntryJSON(roleStoragePrefix+roleName, role)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	if role.TokenTTL > b.System().MaxLeaseTTL() {
		resp := &logical.Response{}
		resp.AddWarning(fmt.Sprintf("ttl of %d exceeds the system max ttl of %d, the latter will be used during login", role.TokenTTL, b.System().MaxLeaseTTL()))
		return resp, nil
	}
	return nil, nil
}

func (b *backend) operationRolesRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)
	role, err := getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}
	cidrs := make([]string, len(role.BoundCIDRs))
	for i, cidr := range role.BoundCIDRs {
		cidrs[i] = cidr.String()
	}

	d := map[string]interface{}{
		"bound_application_ids":  role.BoundAppIDs,
		"bound_space_ids":        role.BoundSpaceIDs,
		"bound_organization_ids": role.BoundOrgIDs,
		"bound_instance_ids":     role.BoundInstanceIDs,
		"disable_ip_matching":    role.DisableIPMatching,
	}

	role.PopulateTokenData(d)

	if len(role.Policies) > 0 {
		d["policies"] = d["token_policies"]
	}
	if len(role.BoundCIDRs) > 0 {
		d["bound_cidrs"] = d["token_bound_cidrs"]
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

	return &logical.Response{
		Data: d,
	}, nil
}

func (b *backend) operationRolesDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)
	if err := req.Storage.Delete(ctx, roleStoragePrefix+roleName); err != nil {
		return nil, err
	}
	return nil, nil
}

func getRole(ctx context.Context, storage logical.Storage, roleName string) (*models.RoleEntry, error) {
	role := &models.RoleEntry{}
	entry, err := storage.Get(ctx, roleStoragePrefix+roleName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	if err := entry.DecodeJSON(role); err != nil {
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
	if len(role.TokenPolicies) == 0 && len(role.Policies) > 0 {
		role.TokenPolicies = role.Policies
	}
	if len(role.TokenBoundCIDRs) == 0 && len(role.BoundCIDRs) > 0 {
		role.TokenBoundCIDRs = role.BoundCIDRs
	}

	return role, nil
}

const pathListRolesHelpSyn = "List the existing roles in this backend."

const pathListRolesHelpDesc = "Roles will be listed by the role name."

const pathRolesHelpSyn = `
Read, write and reference policies and roles that tokens can be made for.
`

const pathRolesHelpDesc = `
This path allows you to read and write roles that are used to
create Vault tokens.
Once configured, credentials will be able to be obtained using this role name
if the caller can successfully provide a client certificate, and sign it
using a valid secret key. The client certificate provided must have been issued
by the configured certficate authority. Its parameters must also match anything
you've listed as "bound".
`
