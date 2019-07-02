package pcf

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault-plugin-auth-pcf/models"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const roleStoragePrefix = "roles/"

func (b *backend) pathListRoles() *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",
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
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("role"),
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
			"bound_cidrs": {
				Type: framework.TypeCommaStringSlice,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Bound CIDRs",
					Value: "192.168.100.14/24",
				},
				Description: `Comma separated string or list of CIDR blocks. If set, specifies the blocks of
IP addresses which can perform the login operation.`,
			},
			"policies": {
				Type:    framework.TypeCommaStringSlice,
				Default: "default",
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Policies",
					Value: "default",
				},
				Description: "Comma separated list of policies on the role.",
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
			"ttl": {
				Type: framework.TypeDurationSecond,
				Description: `Duration in seconds after which the issued token should expire. Defaults
to 0, in which case the value will fallback to the system/mount defaults.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "TTL",
				},
			},
			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "The maximum allowed lifetime of tokens issued using this role.",
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Max TTL",
				},
			},
			"period": {
				Type:    framework.TypeDurationSecond,
				Default: 0,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Period",
					Value: "0",
				},
				Description: `If set, indicates that the token generated using this role
should never expire. The token should be renewed within the
duration specified by this value. At each renewal, the token's
TTL will be set to the value of this parameter.`,
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
	if raw, ok := data.GetOk("bound_cidrs"); ok {
		parsedCIDRs, err := parseutil.ParseAddrs(raw)
		if err != nil {
			return nil, err
		}
		role.BoundCIDRs = parsedCIDRs
	}
	if raw, ok := data.GetOk("policies"); ok {
		role.Policies = raw.([]string)
	}
	if raw, ok := data.GetOk("disable_ip_matching"); ok {
		role.DisableIPMatching = raw.(bool)
	}
	if raw, ok := data.GetOk("ttl"); ok {
		role.TTL = time.Duration(raw.(int)) * time.Second
	}
	if raw, ok := data.GetOk("max_ttl"); ok {
		role.MaxTTL = time.Duration(raw.(int)) * time.Second
	}
	if raw, ok := data.GetOk("period"); ok {
		role.Period = time.Duration(raw.(int)) * time.Second
	}

	if role.MaxTTL > 0 && role.TTL > role.MaxTTL {
		return logical.ErrorResponse("ttl exceeds max_ttl"), nil
	}

	entry, err := logical.StorageEntryJSON(roleStoragePrefix+roleName, role)
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
	return &logical.Response{
		Data: map[string]interface{}{
			"bound_application_ids":  role.BoundAppIDs,
			"bound_space_ids":        role.BoundSpaceIDs,
			"bound_organization_ids": role.BoundOrgIDs,
			"bound_instance_ids":     role.BoundInstanceIDs,
			"bound_cidrs":            cidrs,
			"policies":               role.Policies,
			"disable_ip_matching":    role.DisableIPMatching,
			"ttl":                    role.TTL / time.Second,
			"max_ttl":                role.MaxTTL / time.Second,
			"period":                 role.Period / time.Second,
		},
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
	r := &models.RoleEntry{}
	entry, err := storage.Get(ctx, roleStoragePrefix+roleName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	if err := entry.DecodeJSON(r); err != nil {
		return nil, err
	}
	return r, nil
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
