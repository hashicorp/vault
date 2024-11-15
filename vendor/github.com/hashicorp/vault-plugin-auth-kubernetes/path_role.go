// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubeauth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathsRole returns the path configurations for the CRUD operations on roles
func pathsRole(b *kubeAuthBackend) []*framework.Path {
	p := []*framework.Path{
		{
			Pattern: "role/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathRoleList,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-list"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-list"][1]),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixKubernetes,
				OperationSuffix: "auth-roles",
				Navigation:      true,
				ItemType:        "Role",
			},
		},
		{
			Pattern: "role/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"bound_service_account_names": {
					Type: framework.TypeCommaStringSlice,
					Description: `List of service account names able to access this role. If set to "*" all names
are allowed.`,
				},
				"bound_service_account_namespaces": {
					Type: framework.TypeCommaStringSlice,
					Description: `List of namespaces allowed to access this role. If set to "*" all namespaces
are allowed.`,
				},
				"bound_service_account_namespace_selector": {
					Type: framework.TypeString,
					Description: `A label selector for Kubernetes namespaces which are allowed to access this role.
Accepts either a JSON or YAML object. If set with bound_service_account_namespaces,
the conditions are ORed.`,
				},
				"audience": {
					Type:        framework.TypeString,
					Description: "Optional Audience claim to verify in the jwt.",
				},
				"alias_name_source": {
					Type: framework.TypeString,
					Description: fmt.Sprintf(`Source to use when deriving the Alias name.
valid choices:
	%q : <token.uid> e.g. 474b11b5-0f20-4f9d-8ca5-65715ab325e0 (most secure choice)
	%q : <namespace>/<serviceaccount> e.g. vault/vault-agent
default: %q
`, aliasNameSourceSAUid, aliasNameSourceSAName, aliasNameSourceDefault),
					Default: aliasNameSourceDefault,
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
				"bound_cidrs": {
					Type:        framework.TypeCommaStringSlice,
					Description: tokenutil.DeprecationText("token_bound_cidrs"),
					Deprecated:  true,
				},
			},
			ExistenceCheck: b.pathRoleExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathRoleCreateUpdate,
				logical.UpdateOperation: b.pathRoleCreateUpdate,
				logical.ReadOperation:   b.pathRoleRead,
				logical.DeleteOperation: b.pathRoleDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role"][1]),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixKubernetes,
				OperationSuffix: "auth-role",
				ItemType:        "Role",
				Action:          "Create",
			},
		},
	}

	tokenutil.AddTokenFields(p[1].Fields)
	return p
}

// pathRoleExistenceCheck returns whether the role with the given name exists or not.
func (b *kubeAuthBackend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	role, err := b.role(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return false, err
	}
	return role != nil, nil
}

// pathRoleList is used to list all the Roles registered with the backend.
func (b *kubeAuthBackend) pathRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	roles, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

// pathRoleRead grabs a read lock and reads the options set on the role from the storage
func (b *kubeAuthBackend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	b.l.RLock()
	defer b.l.RUnlock()

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// Create a map of data to be returned
	d := map[string]interface{}{
		"bound_service_account_names":              role.ServiceAccountNames,
		"bound_service_account_namespaces":         role.ServiceAccountNamespaces,
		"bound_service_account_namespace_selector": role.ServiceAccountNamespaceSelector,
	}

	if role.Audience != "" {
		d["audience"] = role.Audience
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
	if role.NumUses > 0 {
		d["num_uses"] = role.NumUses
	}

	d["alias_name_source"] = role.AliasNameSource

	return &logical.Response{
		Data: d,
	}, nil
}

// pathRoleDelete removes the role from storage
func (b *kubeAuthBackend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	// Acquire the lock before deleting the role.
	b.l.Lock()
	defer b.l.Unlock()

	// Delete the role itself
	if err := req.Storage.Delete(ctx, "role/"+strings.ToLower(roleName)); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathRoleCreateUpdate registers a new role with the backend or updates the options
// of an existing role
func (b *kubeAuthBackend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	b.l.Lock()
	defer b.l.Unlock()

	// Check if the role already exists
	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	// Create a new entry object if this is a CreateOperation
	if role == nil && req.Operation == logical.CreateOperation {
		role = &roleStorageEntry{}
	} else if role == nil {
		return nil, fmt.Errorf("role entry not found during update operation")
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
		return logical.ErrorResponse(fmt.Sprintf("token period of '%q' is greater than the backend's maximum lease TTL of '%q'", role.TokenPeriod.String(), b.System().MaxLeaseTTL().String())), nil
	}

	// Check that the TTL value provided is less than the MaxTTL.
	// Sanitizing the TTL and MaxTTL is not required now and can be performed
	// at credential issue time.
	if role.TokenMaxTTL > time.Duration(0) && role.TokenTTL > role.TokenMaxTTL {
		return logical.ErrorResponse("token ttl should not be greater than token max ttl"), nil
	}

	var resp *logical.Response
	if role.TokenMaxTTL > b.System().MaxLeaseTTL() {
		resp = &logical.Response{}
		resp.AddWarning("max_ttl is greater than the system or backend mount's maximum TTL value; issued tokens' max TTL value will be truncated")
	}

	if serviceAccountUUIDs, ok := data.GetOk("bound_service_account_names"); ok {
		role.ServiceAccountNames = serviceAccountUUIDs.([]string)
	} else if req.Operation == logical.CreateOperation {
		role.ServiceAccountNames = data.Get("bound_service_account_names").([]string)
	}
	// Verify names was not empty
	if len(role.ServiceAccountNames) == 0 {
		return logical.ErrorResponse("%q can not be empty", "bound_service_account_names"), nil
	}
	// Verify * was not set with other data
	if len(role.ServiceAccountNames) > 1 && strutil.StrListContains(role.ServiceAccountNames, "*") {
		return logical.ErrorResponse("can not mix %q with values", "*"), nil
	}

	if namespaces, ok := data.GetOk("bound_service_account_namespaces"); ok {
		role.ServiceAccountNamespaces = namespaces.([]string)
	}

	role.ServiceAccountNamespaceSelector = data.Get("bound_service_account_namespace_selector").(string)

	// Verify namespaces is not empty unless selector is set
	saNamespaceLen := len(role.ServiceAccountNamespaces)
	if saNamespaceLen == 0 && role.ServiceAccountNamespaceSelector == "" {
		return logical.ErrorResponse("%q can not be empty if %q is not set",
			"bound_service_account_namespaces", "bound_service_account_namespace_selector"), nil
	}

	// Verify namespace selector is correct
	if role.ServiceAccountNamespaceSelector != "" {
		if _, err := makeNsLabelSelector(role.ServiceAccountNamespaceSelector); err != nil {
			return logical.ErrorResponse("invalid %q configured", "bound_service_account_namespace_selector"), nil
		}
	}

	// Verify * was not set with other data
	if saNamespaceLen > 1 && strutil.StrListContains(role.ServiceAccountNamespaces, "*") {
		return logical.ErrorResponse("can not mix %q with values", "*"), nil
	}

	// optional audience field
	if audience, ok := data.GetOk("audience"); ok {
		role.Audience = audience.(string)
	}

	if source, ok := data.GetOk("alias_name_source"); ok {
		// migrate the role.AliasNameSource to be the default
		// if both it and the field value are unset
		if role.AliasNameSource == aliasNameSourceUnset && source.(string) == aliasNameSourceUnset {
			role.AliasNameSource = data.GetDefaultOrZero("alias_name_source").(string)
		} else {
			role.AliasNameSource = source.(string)
		}
	} else if role.AliasNameSource == aliasNameSourceUnset {
		role.AliasNameSource = data.Get("alias_name_source").(string)
	}

	if err := validateAliasNameSource(role.AliasNameSource); err != nil {
		return logical.ErrorResponse(err.Error()), nil
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
type roleStorageEntry struct {
	tokenutil.TokenParams

	// ServiceAccountNames is the array of service accounts able to
	// access this role.
	ServiceAccountNames []string `json:"bound_service_account_names" mapstructure:"bound_service_account_names" structs:"bound_service_account_names"`

	// ServiceAccountNamespaces is the array of namespaces able to access this
	// role.
	ServiceAccountNamespaces []string `json:"bound_service_account_namespaces" mapstructure:"bound_service_account_namespaces" structs:"bound_service_account_namespaces"`

	// ServiceAccountNamespaceSelector is the label selector string of the
	// namespaces able to access this role.
	ServiceAccountNamespaceSelector string `json:"bound_service_account_namespace_selector" mapstructure:"bound_service_account_namespace_selector" structs:"bound_service_account_namespace_selector"`

	// Audience is an optional jwt claim to verify
	Audience string `json:"audience" mapstructure:"audience" structs:"audience"`

	// AliasNameSource used when deriving the Alias' name.
	AliasNameSource string `json:"alias_name_source" mapstructure:"alias_name_source" structs:"alias_name_source"`

	// Deprecated by TokenParams
	Policies   []string      `json:"policies" structs:"policies" mapstructure:"policies"`
	NumUses    int           `json:"num_uses" mapstructure:"num_uses" structs:"num_uses"`
	TTL        time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`
	MaxTTL     time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
	Period     time.Duration `json:"period" mapstructure:"period" structs:"period"`
	BoundCIDRs []*sockaddr.SockAddrMarshaler
}

var roleHelp = map[string][2]string{
	"role-list": {
		"Lists all the roles registered with the backend.",
		"The list will contain the names of the roles.",
	},
	"role": {
		"Register an role with the backend.",
		`A role is required to authenticate with this backend. The role binds
		kubernetes service account metadata with token policies and settings.
		The bindings, token polices and token settings can all be configured
		using this endpoint`,
	},
}
