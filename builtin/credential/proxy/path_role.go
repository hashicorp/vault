package proxy

import (
	"context"
	"fmt"
	"net/textproto"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	roleStoragePathPrefix = "role/"

	roleNameField            = "name"
	roleAllowedUsersField    = "allowed_users"
	roleRequiredHeadersField = "required_headers"
	rolePoliciesField        = "policies"
	roleTTLField             = "ttl"
	roleMaxTTLField          = "max_ttl"
	rolePeriodField          = "period"

	roleHelpSyn  = `Manage the roles that are registered with the backend.`
	roleHelpDesc = `All logins with the "proxy" credential provider must be made against a ` +
		`preconfigured role.  A role configuration defines constraints around what ` +
		`conditions must be in place for a user to authenticate using that role, along with ` +
		`properties to apply to tokens that are generated when a successful login occurs.`
)

func pathRoleList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/?",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathRolesList,
				Summary:  `List the registered roles`,
			},
		},
		HelpSynopsis: roleHelpSyn,
	}
}

func pathRole(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			roleNameField: &framework.FieldSchema{
				Type:        framework.TypeLowerCaseString,
				Description: "The name of the role",
			},
			roleAllowedUsersField: &framework.FieldSchema{
				Type:     framework.TypeCommaStringSlice,
				Required: true,
				Description: `A comma-separated list of names.  Supports globbing.  At least ` +
					`one glob must match the username`,
			},
			roleRequiredHeadersField: &framework.FieldSchema{
				Type: framework.TypeKVPairs,
				Description: `Mapping of headers (key) and values.  If set, login attempts` +
					`will only be successful if the required case insensitive ` +
					`headers are set to the specified case sensitive value`,
			},
			rolePoliciesField: &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: `Comma-separated list of policies.`,
			},
			roleTTLField: &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `TTL for tokens issued by this backend.  Defaults to ` +
					`system/backend default TTL time.`,
			},
			roleMaxTTLField: &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `Duration in either an integer number of seconds (3600) or an ` +
					`integer time unit (60m) after which the issued token can no ` +
					`longer be renewed.`,
			},
			rolePeriodField: &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `If set, indicates that the token generated using this role ` +
					`should never expire. The token should be renewed within the ` +
					`duration specified by this value. At each renewal, the token's ` +
					`TTL will be set to the value of this parameter.`,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathRoleRead,
				Summary:  "Read an existing role.",
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathRoleCreateUpdate,
				Summary:  "Register a role with the backend.",
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathRoleCreateUpdate,
				Summary:  "Update an existing role.",
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathRoleDelete,
				Summary:  "Delete an existing role.",
			},
		},

		ExistenceCheck: b.pathRoleExistenceCheck,

		HelpSynopsis:    roleHelpSyn,
		HelpDescription: roleHelpDesc,
	}
}

type proxyRole struct {
	AllowedUsers    []string          `json:"allowed_users"`
	RequiredHeaders map[string]string `json:"required_headers"`
	Policies        []string          `json:"policies"`
	Period          time.Duration     `json:"period"`
	TTL             time.Duration     `json:"ttl"`
	MaxTTL          time.Duration     `json:"max_ttl"`
}

func (b *backend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			roleAllowedUsersField:    role.AllowedUsers,
			roleRequiredHeadersField: role.RequiredHeaders,
			rolePoliciesField:        role.Policies,
			rolePeriodField:          int64(role.Period.Seconds()),
			roleTTLField:             int64(role.TTL.Seconds()),
			roleMaxTTLField:          int64(role.MaxTTL.Seconds()),
		},
	}

	return resp, nil
}

func (b *backend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get(roleNameField).(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	if role == nil {
		role = &proxyRole{}
	}

	if allowedUsers, ok := data.GetOk(roleAllowedUsersField); ok {
		role.AllowedUsers = allowedUsers.([]string)
	}
	if len(role.AllowedUsers) == 0 {
		return logical.ErrorResponse(fmt.Sprintf("%s must be set", roleAllowedUsersField)), nil
	}

	if requiredHeaders, ok := data.GetOk(roleRequiredHeadersField); ok {
		role.RequiredHeaders = make(map[string]string)
		for k, v := range requiredHeaders.(map[string]string) {
			nKey := textproto.CanonicalMIMEHeaderKey(k)
			if _, ok := role.RequiredHeaders[nKey]; ok {
				return logical.ErrorResponse("canonical form of required header %q specified multiple times", k), nil
			}

			role.RequiredHeaders[nKey] = v
		}
	}

	if policies, ok := data.GetOk(rolePoliciesField); ok {
		role.Policies = policyutil.ParsePolicies(policies)
	}

	resp := logical.Response{}

	systemDefaultTTL := b.System().DefaultLeaseTTL()
	if ttl, ok := data.GetOk(roleTTLField); ok {
		role.TTL = time.Duration(ttl.(int)) * time.Second
		if role.TTL < time.Duration(0) {
			return logical.ErrorResponse(fmt.Sprintf("%q cannot be negative", roleTTLField)), nil
		}
		if role.TTL > systemDefaultTTL {
			resp.AddWarning(fmt.Sprintf("%q of %d seconds is greater than current mount/system default of %d seconds",
				roleTTLField, role.TTL/time.Second, systemDefaultTTL/time.Second))
		}
	}

	systemMaxTTL := b.System().MaxLeaseTTL()
	if maxTTL, ok := data.GetOk(roleMaxTTLField); ok {
		role.MaxTTL = time.Duration(maxTTL.(int)) * time.Second
		if role.MaxTTL < time.Duration(0) {
			return logical.ErrorResponse(fmt.Sprintf("%q cannot be negative", roleMaxTTLField)), nil
		}

		if role.MaxTTL > systemMaxTTL {
			resp.AddWarning(fmt.Sprintf("%q of %d seconds is greater than current mount/system default of %d seconds",
				roleMaxTTLField, role.MaxTTL/time.Second, systemMaxTTL/time.Second))
		}
	}

	if role.MaxTTL > 0 && role.TTL > role.MaxTTL {
		return logical.ErrorResponse(fmt.Sprintf("%q should not be greater than %s", roleTTLField, roleMaxTTLField)), nil
	}

	if period, ok := data.GetOk(rolePeriodField); ok {
		role.Period = time.Duration(period.(int)) * time.Second
		if role.Period < time.Duration(0) {
			return logical.ErrorResponse(fmt.Sprintf("%q cannot be negative", rolePeriodField)), nil
		}

		if role.Period > systemMaxTTL {
			resp.AddWarning(fmt.Sprintf("%q of %d seconds is greater than the backend's maximum TTL of %d seconds",
				rolePeriodField, role.Period/time.Second, systemMaxTTL/time.Second))
		}
	}

	entry, err := logical.StorageEntryJSON(b.getRoleStoragePath(roleName), role)
	if err != nil {
		return nil, err
	}
	if err = req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (b *backend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("role name required"), nil
	}

	if err := req.Storage.Delete(ctx, b.getRoleStoragePath(roleName)); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRolesList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List(ctx, roleStoragePathPrefix)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

// getRole takes a storage backend and the name and returns the role's storage entry
func (b *backend) getRole(ctx context.Context, s logical.Storage, name string) (*proxyRole, error) {
	raw, err := s.Get(ctx, b.getRoleStoragePath(name))
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	role := proxyRole{}
	if err := raw.DecodeJSON(&role); err != nil {
		return nil, errwrap.Wrapf("error reading role configuration: {{err}}", err)
	}

	return &role, nil
}

func (b *backend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	role, err := b.getRole(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return false, err
	}

	return role != nil, nil
}

func (b *backend) getRoleStoragePath(name string) string {
	return roleStoragePathPrefix + strings.ToLower(name)
}
