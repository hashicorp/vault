package kubeauth

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathsRole returns the path configurations for the CRUD operations on roles
func pathsRole(b *kubeAuthBackend) []*framework.Path {
	return []*framework.Path{
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
					Description: "List of policies on the role.",
				},
				"num_uses": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: `Number of times issued tokens can be used`,
				},
				"ttl": &framework.FieldSchema{
					Type: framework.TypeDurationSecond,
					Description: `Duration in seconds after which the issued token should expire. Defaults
to 0, in which case the value will fall back to the system/mount defaults.`,
				},
				"max_ttl": &framework.FieldSchema{
					Type: framework.TypeDurationSecond,
					Description: `Duration in seconds after which the issued token should not be allowed to
be renewed. Defaults to 0, in which case the value will fall back to the system/mount defaults.`,
				},
				"period": &framework.FieldSchema{
					Type:    framework.TypeDurationSecond,
					Default: 0,
					Description: `If set, indicates that the token generated using this role
should never expire. The token should be renewed within the
duration specified by this value. At each renewal, the token's
TTL will be set to the value of this parameter.`,
				},
			},
			ExistenceCheck: b.pathRoleExistenceCheck,
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
}

// pathRoleExistenceCheck returns whether the role with the given name exists or not.
func (b *kubeAuthBackend) pathRoleExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	role, err := b.role(req.Storage, data.Get("name").(string))
	if err != nil {
		return false, err
	}
	return role != nil, nil
}

// pathRoleList is used to list all the Roles registered with the backend.
func (b *kubeAuthBackend) pathRoleList() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		b.l.RLock()
		defer b.l.RUnlock()

		roles, err := req.Storage.List("role/")
		if err != nil {
			return nil, err
		}
		return logical.ListResponse(roles), nil
	}
}

// pathRoleRead grabs a read lock and reads the options set on the role from the storage
func (b *kubeAuthBackend) pathRoleRead() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		roleName := data.Get("name").(string)
		if roleName == "" {
			return logical.ErrorResponse("missing name"), nil
		}

		b.l.RLock()
		defer b.l.RUnlock()

		role, err := b.role(req.Storage, roleName)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, nil
		}

		// Convert the 'time.Duration' values to second.
		role.TTL /= time.Second
		role.MaxTTL /= time.Second
		role.Period /= time.Second

		// Create a map of data to be returned
		resp := &logical.Response{
			Data: map[string]interface{}{
				"bound_service_account_names":      role.ServiceAccountNames,
				"bound_service_account_namespaces": role.ServiceAccountNamespaces,
				"max_ttl":                          role.MaxTTL,
				"num_uses":                         role.NumUses,
				"policies":                         role.Policies,
				"period":                           role.Period,
				"ttl":                              role.TTL,
			},
		}

		return resp, nil
	}
}

// pathRoleDelete removes the role from storage
func (b *kubeAuthBackend) pathRoleDelete() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		roleName := data.Get("name").(string)
		if roleName == "" {
			return logical.ErrorResponse("missing role name"), nil
		}

		// Acquire the lock before deleting the role.
		b.l.Lock()
		defer b.l.Unlock()

		// Delete the role itself
		if err := req.Storage.Delete("role/" + strings.ToLower(roleName)); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// pathRoleCreateUpdate registers a new role with the backend or updates the options
// of an existing role
func (b *kubeAuthBackend) pathRoleCreateUpdate() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		roleName := data.Get("name").(string)
		if roleName == "" {
			return logical.ErrorResponse("missing role name"), nil
		}

		b.l.Lock()
		defer b.l.Unlock()

		// Check if the role already exists
		role, err := b.role(req.Storage, roleName)
		if err != nil {
			return nil, err
		}

		// Create a new entry object if this is a CreateOperation
		if role == nil && req.Operation == logical.CreateOperation {
			role = &roleStorageEntry{}
		} else if role == nil {
			return nil, fmt.Errorf("role entry not found during update operation")
		}

		if policiesRaw, ok := data.GetOk("policies"); ok {
			role.Policies = policyutil.ParsePolicies(policiesRaw)
		}

		periodRaw, ok := data.GetOk("period")
		if ok {
			role.Period = time.Second * time.Duration(periodRaw.(int))
		} else if req.Operation == logical.CreateOperation {
			role.Period = time.Second * time.Duration(data.Get("period").(int))
		}
		if role.Period > b.System().MaxLeaseTTL() {
			return logical.ErrorResponse(fmt.Sprintf("'period' of '%q' is greater than the backend's maximum lease TTL of '%q'", role.Period.String(), b.System().MaxLeaseTTL().String())), nil
		}

		if tokenNumUsesRaw, ok := data.GetOk("num_uses"); ok {
			role.NumUses = tokenNumUsesRaw.(int)
		} else if req.Operation == logical.CreateOperation {
			role.NumUses = data.Get("num_uses").(int)
		}
		if role.NumUses < 0 {
			return logical.ErrorResponse("num_uses cannot be negative"), nil
		}

		if tokenTTLRaw, ok := data.GetOk("ttl"); ok {
			role.TTL = time.Second * time.Duration(tokenTTLRaw.(int))
		} else if req.Operation == logical.CreateOperation {
			role.TTL = time.Second * time.Duration(data.Get("ttl").(int))
		}

		if tokenMaxTTLRaw, ok := data.GetOk("max_ttl"); ok {
			role.MaxTTL = time.Second * time.Duration(tokenMaxTTLRaw.(int))
		} else if req.Operation == logical.CreateOperation {
			role.MaxTTL = time.Second * time.Duration(data.Get("max_ttl").(int))
		}

		// Check that the TTL value provided is less than the MaxTTL.
		// Sanitizing the TTL and MaxTTL is not required now and can be performed
		// at credential issue time.
		if role.MaxTTL > time.Duration(0) && role.TTL > role.MaxTTL {
			return logical.ErrorResponse("ttl should not be greater than max_ttl"), nil
		}

		var resp *logical.Response
		if role.MaxTTL > b.System().MaxLeaseTTL() {
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
		if err = req.Storage.Put(entry); err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// roleStorageEntry stores all the options that are set on an role
type roleStorageEntry struct {
	// Policies that are to be required by the token to access this role
	Policies []string `json:"policies" structs:"policies" mapstructure:"policies"`

	// TokenNumUses defines the number of allowed uses of the token issued
	NumUses int `json:"num_uses" mapstructure:"num_uses" structs:"num_uses"`

	// Duration before which an issued token must be renewed
	TTL time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`

	// Duration after which an issued token should not be allowed to be renewed
	MaxTTL time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`

	// Period, if set, indicates that the token generated using this role
	// should never expire. The token should be renewed within the duration
	// specified by this value. The renewal duration will be fixed if the
	// value is not modified on the role. If the `Period` in the role is modified,
	// a token will pick up the new value during its next renewal.
	Period time.Duration `json:"period" mapstructure:"period" structs:"period"`

	// ServiceAccountNames is the array of service accounts able to
	// access this role.
	ServiceAccountNames []string `json:"bound_service_account_names" mapstructure:"bound_service_account_names" structs:"bound_service_account_names"`

	// ServiceAccountNamespaces is the array of namespaces able to access this
	// role.
	ServiceAccountNamespaces []string `json:"bound_service_account_namespaces" mapstructure:"bound_service_account_namespaces" structs:"bound_service_account_namespaces"`
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
