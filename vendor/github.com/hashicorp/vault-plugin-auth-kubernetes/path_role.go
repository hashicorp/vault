package kubeauth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathsRole returns the path configurations for the CRUD operations on roles
func pathsRole(b *kubeAuthBackend) []*framework.Path {
	p := []*framework.Path{
		&framework.Path{
			Pattern: "role/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathRoleList(),
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-list"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-list"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"bound_service_account_names": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `List of service account names able to access this role. If set to "*" all names
are allowed, both this and bound_service_account_namespaces can not be "*"`,
				},
				"bound_service_account_namespaces": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `List of namespaces allowed to access this role. If set to "*" all namespaces
are allowed, both this and bound_service_account_names can not be set to "*"`,
				},
				"policies": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: tokenutil.DeprecationText("token_policies"),
					Deprecated:  true,
				},
				"num_uses": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: tokenutil.DeprecationText("token_num_uses"),
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
			ExistenceCheck: b.pathRoleExistenceCheck(),
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathRoleCreateUpdate(),
				logical.UpdateOperation: b.pathRoleCreateUpdate(),
				logical.ReadOperation:   b.pathRoleRead(),
				logical.DeleteOperation: b.pathRoleDelete(),
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role"][1]),
		},
	}

	tokenutil.AddTokenFields(p[1].Fields)
	return p
}

// pathRoleExistenceCheck returns whether the role with the given name exists or not.
func (b *kubeAuthBackend) pathRoleExistenceCheck() framework.ExistenceFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
		b.l.RLock()
		defer b.l.RUnlock()

		role, err := b.role(ctx, req.Storage, data.Get("name").(string))
		if err != nil {
			return false, err
		}
		return role != nil, nil
	}
}

// pathRoleList is used to list all the Roles registered with the backend.
func (b *kubeAuthBackend) pathRoleList() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		b.l.RLock()
		defer b.l.RUnlock()

		roles, err := req.Storage.List(ctx, "role/")
		if err != nil {
			return nil, err
		}
		return logical.ListResponse(roles), nil
	}
}

// pathRoleRead grabs a read lock and reads the options set on the role from the storage
func (b *kubeAuthBackend) pathRoleRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
			"bound_service_account_names":      role.ServiceAccountNames,
			"bound_service_account_namespaces": role.ServiceAccountNamespaces,
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

		return &logical.Response{
			Data: d,
		}, nil
	}
}

// pathRoleDelete removes the role from storage
func (b *kubeAuthBackend) pathRoleDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
}

// pathRoleCreateUpdate registers a new role with the backend or updates the options
// of an existing role
func (b *kubeAuthBackend) pathRoleCreateUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
			return logical.ErrorResponse("\"bound_service_account_names\" can not be empty"), nil
		}
		// Verify * was not set with other data
		if len(role.ServiceAccountNames) > 1 && strutil.StrListContains(role.ServiceAccountNames, "*") {
			return logical.ErrorResponse("can not mix \"*\" with values"), nil
		}

		if namespaces, ok := data.GetOk("bound_service_account_namespaces"); ok {
			role.ServiceAccountNamespaces = namespaces.([]string)
		} else if req.Operation == logical.CreateOperation {
			role.ServiceAccountNamespaces = data.Get("bound_service_account_namespaces").([]string)
		}
		// Verify namespaces is not empty
		if len(role.ServiceAccountNamespaces) == 0 {
			return logical.ErrorResponse("\"bound_service_account_namespaces\" can not be empty"), nil
		}
		// Verify * was not set with other data
		if len(role.ServiceAccountNamespaces) > 1 && strutil.StrListContains(role.ServiceAccountNamespaces, "*") {
			return logical.ErrorResponse("can not mix \"*\" with values"), nil
		}

		// Verify that both names and namespaces are not set to "*"
		if strutil.StrListContains(role.ServiceAccountNames, "*") && strutil.StrListContains(role.ServiceAccountNamespaces, "*") {
			return logical.ErrorResponse("service_account_names and service_account_namespaces can not both be \"*\""), nil
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
