package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathsRole returns the path configurations for the CRUD operations on roles
func pathsRole(b *azureAuthBackend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "role/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathRoleList,
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
				"bound_subscription_ids": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `Comma-separated list of subscription ids that login 
is restricted to.`,
				},
				"bound_resource_groups": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `Comma-separated list of resource groups that login 
is restricted to.`,
				},
				"bound_group_ids": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `Comma-separated list of group ids that login 
is restricted to.`,
				},
				"bound_service_principal_ids": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `Comma-separated list of service principal ids that login 
is restricted to.`,
				},
				"bound_locations": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `Comma-separated list of locations that login 
is restricted to.`,
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
		},
	}
}

type azureRole struct {
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
func (b *azureAuthBackend) pathRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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

	// Convert the 'time.Duration' values to second.
	role.TTL /= time.Second
	role.MaxTTL /= time.Second
	role.Period /= time.Second

	// Create a map of data to be returned
	resp := &logical.Response{
		Data: map[string]interface{}{
			"max_ttl":  role.MaxTTL,
			"num_uses": role.NumUses,
			"policies": role.Policies,
			"period":   role.Period,
			"ttl":      role.TTL,
			"bound_service_principal_ids": role.BoundServicePrincipalIDs,
			"bound_group_ids":             role.BoundGroupIDs,
			"bound_subscription_ids":      role.BoundSubscriptionsIDs,
			"bound_resource_groups":       role.BoundResourceGroups,
			"bound_locations":             role.BoundLocations,
		},
	}

	return resp, nil
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
