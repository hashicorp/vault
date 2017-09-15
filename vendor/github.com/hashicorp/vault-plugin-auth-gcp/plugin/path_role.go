package gcpauth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	iamRoleType            = "iam"
	serviceAccountAttr     = "service-accounts"
	serviceAccountWildcard = "*"

	errEmptyRoleName           = "role name is required"
	errEmptyIamServiceAccounts = "IAM role type must have at least one service account"
)

// pathsRole creates paths for listing roles and CRUD operations.
func pathsRole(b *GcpAuthBackend) []*framework.Path {
	paths := []*framework.Path{
		{
			Pattern: fmt.Sprintf("role/%s", framework.GenericNameRegex("name")),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"type": {
					Type:        framework.TypeString,
					Description: "Type of the role. Currently supported: iam",
				},
				"policies": {
					Type:        framework.TypeString,
					Description: "Policies to be set on tokens issued using this role.",
				},
				"project_id": {
					Type:        framework.TypeString,
					Description: `The id of the project for service accounts allowed to authenticate to this role.`,
				},
				"max_jwt_exp": {
					Type:        framework.TypeDurationSecond,
					Default:     defaultMaxJwtExpMin * 3600,
					Description: `Duration in seconds from time of validation that a JWT must expire within.`,
				},
				// Token Limits
				"ttl": {
					Type:        framework.TypeDurationSecond,
					Default:     0,
					Description: `Duration in seconds after which the issued token should expire. Defaults to 0, in which case the value will fallback to the system/mount defaults.`,
				},
				"max_ttl": {
					Type:        framework.TypeDurationSecond,
					Default:     0,
					Description: "The maximum allowed lifetime of tokens issued using this role.",
				},
				"period": {
					Type:    framework.TypeDurationSecond,
					Default: 0,
					Description: `
If set, indicates that the token generated using this role should never expire. The token should be renewed within the
 duration specified by this value. At each renewal, the token's TTL will be set to the value of this parameter.`,
				},
				// IAM Role Domain
				"service_accounts": {
					Type: framework.TypeString,
					Description: `
A comma-seperated list of service accounts to allow to login as this role. If the single value "*" is given, this
is assumed to be all service accounts under the role's project.`,
				},
			},

			ExistenceCheck: b.pathRoleExistenceCheck,

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.DeleteOperation: b.pathRoleDelete,
				logical.ReadOperation:   b.pathRoleRead,
				logical.CreateOperation: b.pathRoleCreateUpdate,
				logical.UpdateOperation: b.pathRoleCreateUpdate,
			},
			HelpSynopsis:    pathRoleHelpSyn,
			HelpDescription: pathRoleHelpDesc,
		},
		// Paths for listing roles
		{
			Pattern: "role/?",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathRoleList,
			},

			HelpSynopsis:    pathListRolesHelpSyn,
			HelpDescription: pathListRolesHelpDesc,
		},
		{
			Pattern: "roles/?",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathRoleList,
			},

			HelpSynopsis:    pathListRolesHelpSyn,
			HelpDescription: pathListRolesHelpDesc,
		},
	}

	paths = append(paths, b.pathRoleEditListAttr(serviceAccountAttr)...)
	return paths
}

func (b *GcpAuthBackend) pathRoleExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.role(req.Storage, data.Get("name").(string))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *GcpAuthBackend) pathRoleDelete(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse(errEmptyRoleName), nil
	}

	if err := req.Storage.Delete(fmt.Sprintf("role/%s", name)); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *GcpAuthBackend) pathRoleRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse(errEmptyRoleName), nil
	}

	role, err := b.role(req.Storage, name)
	if err != nil {
		return nil, err
	} else if role == nil {
		return nil, nil
	}

	roleMap := map[string]interface{}{
		"role_type":   role.RoleType,
		"project_id":  role.ProjectId,
		"policies":    role.Policies,
		"max_jwt_exp": int64(role.MaxJwtExp / time.Second),
		"ttl":         int64(role.TTL / time.Second),
		"max_ttl":     int64(role.MaxTTL / time.Second),
		"period":      int64(role.Period / time.Second),
	}

	switch role.RoleType {
	case iamRoleType:
		roleMap["service_accounts"] = role.ServiceAccounts
	}

	return &logical.Response{
		Data: roleMap}, nil
}

func (b *GcpAuthBackend) pathRoleCreateUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := strings.ToLower(data.Get("name").(string))
	if name == "" {
		return logical.ErrorResponse(errEmptyRoleName), nil
	}

	role, err := b.role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		role = &gcpRole{}
	}

	warnResp, err := role.updateRole(b.System(), req.Operation, data)

	errResp, err := b.storeRole(req.Storage, name, role)
	if err != nil {
		return nil, err
	} else if errResp != nil {
		return errResp, err
	}

	return warnResp, nil
}

func (b *GcpAuthBackend) pathRoleList(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List("role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

const pathRoleHelpSyn = `Create a GCP role with associated policies and required attributes.`
const pathRoleHelpDesc = `
A role is required to login under the GCP auth backend. A role binds Vault policies and has
required attributes that an authenticating entity must fulfill to login against this role.
After authenticating the instance, Vault uses the bound policies to determine which resources
the authorization token for the instance can access.
`

const pathListRolesHelpSyn = `Lists all the roles that are registered with Vault.`
const pathListRolesHelpDesc = `Lists all roles under the GCP backends by name.`

func (b *GcpAuthBackend) pathRoleEditListAttr(attr string) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: fmt.Sprintf("role/%s/%s", framework.GenericNameRegex("name"), attr),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"add": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Values to add.",
				},
				"remove": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Values to remove.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleEditListOperator(attr),
			},
			HelpSynopsis:    pathAddToListHelpSyn,
			HelpDescription: pathAddToListHelpDesc,
		},
	}
}

func (b *GcpAuthBackend) pathRoleEditListOperator(attr string) framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		roleName := data.Get("name").(string)
		if roleName == "" {
			return logical.ErrorResponse(errEmptyRoleName), nil
		}

		toAdd := data.Get("add").([]string)
		toRemove := data.Get("remove").([]string)
		if len(toAdd) == 0 && len(toRemove) == 0 {
			return logical.ErrorResponse("must provide at least one value to add or remove"), nil
		}

		role, err := b.role(req.Storage, roleName)
		if err != nil {
			return nil, err
		}

		switch attr {
		case serviceAccountAttr:
			if role.RoleType != iamRoleType {
				return logical.ErrorResponse(fmt.Sprintf("Cannot edit service accounts on non-IAM roles, role is type %s", role.RoleType)), nil
			}
			role.ServiceAccounts = editStringValues(role.ServiceAccounts, toAdd, toRemove)
		default:
			return nil, fmt.Errorf("unsupported attribute '%s'", attr)
		}

		return b.storeRole(req.Storage, roleName, role)
	}
}

func editStringValues(initial []string, toAdd []string, toRemove []string) []string {
	strMap := map[string]struct{}{}
	for _, name := range initial {
		strMap[name] = struct{}{}
	}

	for _, name := range toAdd {
		strMap[name] = struct{}{}
	}

	for _, name := range toRemove {
		delete(strMap, name)
	}

	updated := make([]string, len(strMap))

	i := 0
	for k := range strMap {
		updated[i] = k
		i++
	}

	return updated
}

const pathAddToListHelpSyn = `Add and remove values for an list attribute on an existing role`
const pathAddToListHelpDesc = `This path allows a user to add values to and/or remove from a list of values
for a given role's list attribute`

// role reads a gcpRole from storage. This assumes the caller has already obtained the role lock.
func (b *GcpAuthBackend) role(s logical.Storage, name string) (*gcpRole, error) {
	entry, err := s.Get(fmt.Sprintf("role/%s", strings.ToLower(name)))

	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	role := &gcpRole{}
	if err := entry.DecodeJSON(role); err != nil {
		return nil, err
	}

	return role, nil
}

// storeRole saves the gcpRole to storage.
func (b *GcpAuthBackend) storeRole(s logical.Storage, roleName string, role *gcpRole) (*logical.Response, error) {
	if err := role.validate(b.System()); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	entry, err := logical.StorageEntryJSON(fmt.Sprintf("role/%s", roleName), role)
	if err != nil {
		return nil, err
	}

	return nil, s.Put(entry)
}

type gcpRole struct {
	// Type of this role. See path_role constants for currently supported types.
	RoleType string `json:"role_type" structs:"role_type" mapstructure:"role_type"`

	// Project ID in GCP for authorized entities.
	ProjectId string `json:"project_id" structs:"project_id" mapstructure:"project_id"`

	// Policies for Vault to assign to authorized entities.
	Policies []string `json:"policies" structs:"policies" mapstructure:"policies"`

	// MaxJwtExp is the duration from time of authentication that a JWT used to authenticate to role must expire within.
	MaxJwtExp time.Duration `json:"max_jwt_exp" structs:"max_jwt_exp" mapstructure:"max_jwt_exp"`

	// TTL of Vault auth leases under this role.
	TTL time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`

	// Max total TTL including renewals, of Vault auth leases under this role.
	MaxTTL time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`

	// Period, If set, indicates that this token should not expire and
	// should be automatically renewed within this time period
	// with TTL equal to this value.
	Period time.Duration `json:"period" structs:"period" mapstructure:"period"`

	// IAM-specific attributes
	ServiceAccounts []string `json:"service_accounts" structs:"service_accounts" mapstructure:"service_accounts"`
}

// Update updates the given role with values parsed/validated from given FieldData.
// Exactly one of the response and error will be nil. The response is only used to pass back warnings.
// This method does not validate the role. Validation is done before storage.
func (role *gcpRole) updateRole(sys logical.SystemView, op logical.Operation, data *framework.FieldData) (*logical.Response, error) {
	warnResp := &logical.Response{}

	// Set role type
	roleTypeRaw, ok := data.GetOk("type")
	if ok {
		if op == logical.UpdateOperation {
			return nil, errors.New("role type cannot be changed for an existing role")
		}
		role.RoleType = roleTypeRaw.(string)
	} else if op == logical.CreateOperation {
		return nil, errors.New("role type must be provided for a new role")
	}

	//Update fields specific to this type
	switch role.RoleType {
	case iamRoleType:
		if err := role.updateIamFields(data, op); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("role type '%s' is not supported", role.RoleType)
	}

	// Update policies.
	policies, ok := data.GetOk("policies")
	if ok {
		role.Policies = policyutil.ParsePolicies(policies.(string))
	} else if op == logical.CreateOperation {
		role.Policies = policyutil.ParsePolicies(data.Get("policies").(string))
	}

	// Update GCP project id.
	projectIdRaw, ok := data.GetOk("project_id")
	if ok {
		role.ProjectId = projectIdRaw.(string)
	}

	// Update max JWT exp duration.
	maxJwtExp, ok := data.GetOk("max_jwt_exp")
	if ok {
		role.MaxJwtExp = time.Duration(maxJwtExp.(int)) * time.Second
	} else {
		role.MaxJwtExp = time.Duration(defaultMaxJwtExpMin) * time.Minute
	}

	// Update token TTL.
	ttlRaw, ok := data.GetOk("ttl")
	if ok {
		role.TTL = time.Duration(ttlRaw.(int)) * time.Second
		defaultLeaseTTL := sys.DefaultLeaseTTL()
		if role.TTL > defaultLeaseTTL {
			warnResp.AddWarning(fmt.Sprintf(
				"Given ttl of %d seconds greater than current mount/system default of %d seconds; ttl will be capped at login time",
				role.TTL/time.Second, defaultLeaseTTL/time.Second))
		}
	} else if op == logical.CreateOperation {
		role.TTL = time.Duration(data.Get("ttl").(int)) * time.Second
	}

	// Update token Max TTL.
	maxTTLRaw, ok := data.GetOk("max_ttl")
	if ok {
		role.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
		systemMaxTTL := sys.MaxLeaseTTL()
		if role.MaxTTL > systemMaxTTL {
			warnResp.AddWarning(fmt.Sprintf(
				"Given max_ttl of %d seconds greater than current mount/system default of %d seconds; max_ttl will be capped at login time",
				role.MaxTTL/time.Second, systemMaxTTL/time.Second))
		}
	} else if op == logical.CreateOperation {
		role.MaxTTL = time.Duration(data.Get("max_ttl").(int)) * time.Second
	}

	// Update token period.
	periodRaw, ok := data.GetOk("period")
	if ok {
		role.Period = time.Second * time.Duration(periodRaw.(int))
	} else if op == logical.CreateOperation {
		role.Period = time.Second * time.Duration(data.Get("period").(int))
	}

	if len(warnResp.Warnings) == 0 {
		warnResp = nil
	}
	return warnResp, nil
}

func (role *gcpRole) validate(sys logical.SystemView) error {
	switch role.RoleType {
	case iamRoleType:
		if err := role.validateIamFields(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("role type '%s' is invalid", role.RoleType)
	}

	if len(role.Policies) == 0 {
		return errors.New("role must have at least one bound policy")
	}

	if role.ProjectId == "" {
		return errors.New("role cannot have empty project_id")
	}

	if role.MaxJwtExp > time.Hour {
		return errors.New("max_jwt_exp cannot be more than one hour")
	}

	if role.MaxTTL < time.Duration(0) {
		return errors.New("max_ttl cannot be negative")
	}
	if role.MaxTTL != 0 && role.MaxTTL < role.TTL {
		return errors.New("ttl should be shorter than max_ttl")
	}

	if role.Period > sys.MaxLeaseTTL() {
		return fmt.Errorf("'period' of '%s' is greater than the backend's maximum lease TTL of '%s'", role.Period.String(), sys.MaxLeaseTTL().String())
	}

	return nil
}

// updateIamFields updates IAM-only fields for a role.
func (role *gcpRole) updateIamFields(data *framework.FieldData, op logical.Operation) error {
	serviceAccountsRaw, ok := data.GetOk("service_accounts")
	if ok {
		role.ServiceAccounts = strings.Split(serviceAccountsRaw.(string), ",")
	} else if op == logical.CreateOperation {
		return errors.New(errEmptyIamServiceAccounts)
	}

	return nil
}

// validateIamFields validates the IAM-only fields for a role.
func (role *gcpRole) validateIamFields() error {
	if len(role.ServiceAccounts) == 0 {
		return errors.New(errEmptyIamServiceAccounts)
	}

	if len(role.ServiceAccounts) > 1 && strutil.StrListContains(role.ServiceAccounts, serviceAccountWildcard) {
		return fmt.Errorf("cannot provide IAM service account wildcard '%s' (for all service accounts) with other service accounts", serviceAccountWildcard)
	}
	return nil
}
