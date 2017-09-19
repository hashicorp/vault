package gcpauth

import (
	"errors"
	"fmt"
	"github.com/hashicorp/vault-plugin-auth-gcp/plugin/util"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"strings"
	"time"
)

const (
	// Role types
	iamRoleType = "iam"
	gceRoleType = "gce"

	// Errors
	errEmptyRoleName           = "role name is required"
	errEmptyRoleType           = "role type cannot be empty"
	errEmptyProjectId          = "project id cannot be empty"
	errEmptyIamServiceAccounts = "IAM role type must have at least one service account"

	errTemplateEditListWrongType   = "role is type '%s', cannot edit attribute '%s' (expected role type: '%s')"
	errTemplateInvalidRoleTypeArgs = "invalid args found for role of type %s: %s"

	// Other
	serviceAccountsWildcard = "*"

	// Default duration that JWT tokens must expire within to be accepted (currently only IAM)
	defaultIamMaxJwtExpMinutes int = 15

	// Max allowed duration that all JWT tokens must expire within to be accepted
	maxJwtExpMaxMinutes int = 60
)

var baseRoleFieldSchema map[string]*framework.FieldSchema = map[string]*framework.FieldSchema{
	"name": {
		Type:        framework.TypeString,
		Description: "Name of the role.",
	},
	"type": {
		Type:        framework.TypeString,
		Description: "Type of the role. Currently supported: iam, gce",
	},
	"policies": {
		Type:        framework.TypeCommaStringSlice,
		Description: "Policies to be set on tokens issued using this role.",
	},
	// Token Limits
	"ttl": {
		Type:    framework.TypeDurationSecond,
		Default: 0,
		Description: `
	Duration in seconds after which the issued token should expire. Defaults to 0,
	in which case the value will fallback to the system/mount defaults.`,
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
	// -- GCP Information
	"project_id": {
		Type:        framework.TypeString,
		Description: `The id of the project that authorized instances must belong to for this role.`,
	},
	"bound_service_accounts": {
		Type: framework.TypeCommaStringSlice,
		Description: `
	Can be set for both 'iam' and 'gce' roles (required for 'iam'). A comma-seperated list of authorized service accounts.
	If the single value "*" is given, this is assumed to be all service accounts under the role's project. If this
	is set on a GCE role, the inferred service account from the instance metadata token will be used.`,
	},
	"service_accounts": {
		Type:        framework.TypeCommaStringSlice,
		Description: `Deprecated, use bound_service_accounts instead.`,
	},
}

var iamOnlyFieldSchema map[string]*framework.FieldSchema = map[string]*framework.FieldSchema{
	"max_jwt_exp": {
		Type:        framework.TypeDurationSecond,
		Default:     defaultIamMaxJwtExpMinutes * 60,
		Description: `Currently enabled for 'iam' only. Duration in seconds from time of validation that a JWT must expire within.`,
	},
	"allow_gce_inference": {
		Type:        framework.TypeBool,
		Default:     true,
		Description: `'iam' roles only. If false, Vault will not not allow GCE instances to login in against this role`,
	},
}

var gceOnlyFieldSchema map[string]*framework.FieldSchema = map[string]*framework.FieldSchema{
	"bound_zone": {
		Type: framework.TypeString,
		Description: `
"gce" roles only. If set, determines the zone that a GCE instance must belong to. If a group is provided, it is assumed
to be a zonal group and the group must belong to this zone.`,
	},
	"bound_region": {
		Type: framework.TypeString,
		Description: `
"gce" roles only. If set, determines the region that a GCE instance must belong to. If a group is provided, it is
assumed to be a regional group and the group must belong to this region. If zone is provided, region will be ignored`,
	},
	"bound_instance_group": {
		Type:        framework.TypeString,
		Description: `"gce" roles only. If set, determines the instance group that an authorized instance must belong to.`,
	},
	"bound_labels": {
		Type: framework.TypeCommaStringSlice,
		Description: `
"gce" roles only. A comma-separated list of Google Cloud Platform labels formatted as "$key:$value" strings that are
required for authorized GCE instances.`,
	},
}

// pathsRole creates paths for listing roles and CRUD operations.
func pathsRole(b *GcpAuthBackend) []*framework.Path {
	roleFieldSchema := map[string]*framework.FieldSchema{}
	for k, v := range baseRoleFieldSchema {
		roleFieldSchema[k] = v
	}
	for k, v := range iamOnlyFieldSchema {
		roleFieldSchema[k] = v
	}
	for k, v := range gceOnlyFieldSchema {
		roleFieldSchema[k] = v
	}

	paths := []*framework.Path{
		{
			Pattern:        fmt.Sprintf("role/%s", framework.GenericNameRegex("name")),
			Fields:         roleFieldSchema,
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

		// Edit service accounts on an IAM role
		{
			Pattern: fmt.Sprintf("role/%s/service-accounts", framework.GenericNameRegex("name")),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"add": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Service-account emails or IDs to add.",
				},
				"remove": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Service-account emails or IDs to remove.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleEditIamServiceAccounts,
			},
			HelpSynopsis:    "Add or remove service accounts for an existing `iam` role",
			HelpDescription: "Add or remove service accounts from the list bound to an existing `iam` role",
		},

		// Edit labels on an GCE role
		{
			Pattern: fmt.Sprintf("role/%s/labels", framework.GenericNameRegex("name")),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"add": {
					Type:        framework.TypeCommaStringSlice,
					Description: "BoundLabels to add (in $key:$value)",
				},
				"remove": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Label key values to remove",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleEditGceLabels,
			},
			HelpSynopsis: "Add or remove labels for an existing 'gce' role",
			HelpDescription: `Add or remove labels for an existing 'gce' role. 'add' labels should be
			of format '$key:$value' and 'remove' labels should be a list of keys to remove.`,
		},
	}

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
		"role_type":              role.RoleType,
		"project_id":             role.ProjectId,
		"policies":               role.Policies,
		"ttl":                    int64(role.TTL / time.Second),
		"max_ttl":                int64(role.MaxTTL / time.Second),
		"period":                 int64(role.Period / time.Second),
		"bound_service_accounts": role.BoundServiceAccounts,
	}

	switch role.RoleType {
	case iamRoleType:
		roleMap["max_jwt_exp"] = int64(role.MaxJwtExp / time.Second)
		roleMap["allow_gce_inference"] = role.AllowGCEInference
	case gceRoleType:
		roleMap["bound_zone"] = role.BoundZone
		roleMap["bound_region"] = role.BoundRegion
		roleMap["bound_instance_group"] = role.BoundInstanceGroup
		// Ensure values are not nil to avoid errors during plugin RPC conversions.
		if role.BoundLabels != nil && len(role.BoundLabels) > 0 {
			roleMap["bound_labels"] = role.BoundLabels
		} else {
			roleMap["bound_labels"] = ""
		}
	}

	return &logical.Response{
		Data: roleMap,
	}, nil
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

	if err := role.updateRole(b.System(), req.Operation, data); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	return b.storeRole(req.Storage, name, role)
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

func (b *GcpAuthBackend) pathRoleEditIamServiceAccounts(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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

	if role.RoleType != iamRoleType {
		return logical.ErrorResponse(fmt.Sprintf(errTemplateEditListWrongType, role.RoleType, "service_accounts", iamRoleType)), nil
	}
	role.BoundServiceAccounts = editStringValues(role.BoundServiceAccounts, toAdd, toRemove)

	return b.storeRole(req.Storage, roleName, role)
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

func (b *GcpAuthBackend) pathRoleEditGceLabels(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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

	if role.RoleType != gceRoleType {
		return logical.ErrorResponse(fmt.Sprintf(errTemplateEditListWrongType, role.RoleType, "labels", gceRoleType)), nil
	}

	labelsToAdd, invalidLabels := util.ParseGcpLabels(toAdd)
	if len(invalidLabels) > 0 {
		return logical.ErrorResponse(fmt.Sprintf("given invalid labels to add: %q", invalidLabels)), nil
	}
	for k, v := range labelsToAdd {
		role.BoundLabels[k] = v
	}

	for _, k := range toRemove {
		delete(role.BoundLabels, k)
	}

	return b.storeRole(req.Storage, roleName, role)
}

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
// The returned response may contain either warnings or an error response,
// but will be nil if error is not nil
func (b *GcpAuthBackend) storeRole(s logical.Storage, roleName string, role *gcpRole) (*logical.Response, error) {
	var resp *logical.Response
	warnings, err := role.validate(b.System())

	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	} else if len(warnings) > 0 {
		resp = &logical.Response{
			Warnings: warnings,
		}
	}

	entry, err := logical.StorageEntryJSON(fmt.Sprintf("role/%s", roleName), role)
	if err != nil {
		return nil, err
	}

	if err := s.Put(entry); err != nil {
		return nil, err
	}

	return resp, nil
}

type gcpRole struct {
	// Type of this role. See path_role constants for currently supported types.
	RoleType string `json:"role_type" structs:"role_type" mapstructure:"role_type"`

	// Project ID in GCP for authorized entities.
	ProjectId string `json:"project_id" structs:"project_id" mapstructure:"project_id"`

	// Policies for Vault to assign to authorized entities.
	Policies []string `json:"policies" structs:"policies" mapstructure:"policies"`

	// TTL of Vault auth leases under this role.
	TTL time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`

	// Max total TTL including renewals, of Vault auth leases under this role.
	MaxTTL time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`

	// Period, If set, indicates that this token should not expire and
	// should be automatically renewed within this time period
	// with TTL equal to this value.
	Period time.Duration `json:"period" structs:"period" mapstructure:"period"`

	// Service accounts allowed to login under this role.
	BoundServiceAccounts []string `json:"bound_service_accounts" structs:"bound_service_accounts" mapstructure:"bound_service_accounts"`

	// --| IAM-only attributes |--
	// MaxJwtExp is the duration from time of authentication that a JWT used to authenticate to role must expire within.
	// TODO(emilymye): Allow this to be updated for GCE roles once 'exp' parameter has been allowed for GCE metadata.
	MaxJwtExp time.Duration `json:"max_jwt_exp" structs:"max_jwt_exp" mapstructure:"max_jwt_exp"`

	// AllowGCEInference, if false, does not allow a GCE instance to login under this 'iam' role. If true (default),
	// a service account is inferred from the instance metadata and used as the authenticating instance.
	AllowGCEInference bool `json:"allow_gce_inference" structs:"allow_gce_inference" mapstructure:"allow_gce_inference"`

	// --| GCE-only attributes |--
	// BoundRegion that instances must belong to in order to login under this role.
	BoundRegion string `json:"bound_region" structs:"bound_region" mapstructure:"bound_region"`

	// BoundZone that instances must belong to in order to login under this role.
	BoundZone string `json:"bound_zone" structs:"bound_zone" mapstructure:"bound_zone"`

	// Instance group that instances must belong to in order to login under this role.
	BoundInstanceGroup string `json:"bound_instance_group" structs:"bound_instance_group" mapstructure:"bound_instance_group"`

	// BoundLabels that instances must currently have set in order to login under this role.
	BoundLabels map[string]string `json:"bound_labels" structs:"bound_labels" mapstructure:"bound_labels"`
}

// Update updates the given role with values parsed/validated from given FieldData.
// Exactly one of the response and error will be nil. The response is only used to pass back warnings.
// This method does not validate the role. Validation is done before storage.
func (role *gcpRole) updateRole(sys logical.SystemView, op logical.Operation, data *framework.FieldData) error {
	// Set role type
	roleTypeRaw, ok := data.GetOk("type")
	if ok {
		if op == logical.UpdateOperation {
			return errors.New("role type cannot be changed for an existing role")
		}
		role.RoleType = roleTypeRaw.(string)
	} else if op == logical.CreateOperation {
		return errors.New(errEmptyRoleType)
	}

	// Update policies.
	policies, ok := data.GetOk("policies")
	if ok {
		role.Policies = policyutil.ParsePolicies(policies)
	} else if op == logical.CreateOperation {
		role.Policies = policyutil.ParsePolicies(data.Get("policies"))
	}

	// Update GCP project id.
	projectIdRaw, ok := data.GetOk("project_id")
	if ok {
		role.ProjectId = projectIdRaw.(string)
	}

	// Update token TTL.
	ttlRaw, ok := data.GetOk("ttl")
	if ok {
		role.TTL = time.Duration(ttlRaw.(int)) * time.Second

	} else if op == logical.CreateOperation {
		role.TTL = time.Duration(data.Get("ttl").(int)) * time.Second
	}

	// Update token Max TTL.
	maxTTLRaw, ok := data.GetOk("max_ttl")
	if ok {
		role.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
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

	// Update bound GCP service accounts.
	serviceAccountsRaw, ok := data.GetOk("bound_service_accounts")
	if ok {
		role.BoundServiceAccounts = serviceAccountsRaw.([]string)
	} else {
		// Check for older version of param name
		serviceAccountsRaw, ok := data.GetOk("service_accounts")
		if ok {
			role.BoundServiceAccounts = serviceAccountsRaw.([]string)
		}
	}

	// Update fields specific to this type
	switch role.RoleType {
	case iamRoleType:
		if err := checkInvalidRoleTypeArgs(data, gceOnlyFieldSchema); err != nil {
			return err
		}
		if err := role.updateIamFields(data, op); err != nil {
			return err
		}
	case gceRoleType:
		if err := checkInvalidRoleTypeArgs(data, iamOnlyFieldSchema); err != nil {
			return err
		}
		if err := role.updateGceFields(data, op); err != nil {
			return err
		}
	}

	return nil
}

func (role *gcpRole) validate(sys logical.SystemView) (warnings []string, err error) {
	warnings = []string{}

	switch role.RoleType {
	case iamRoleType:
		if warnings, err = role.validateForIAM(); err != nil {
			return warnings, err
		}
	case gceRoleType:
		if warnings, err = role.validateForGCE(); err != nil {
			return warnings, err
		}
	case "":
		return warnings, errors.New(errEmptyRoleType)
	default:
		return warnings, fmt.Errorf("role type '%s' is invalid", role.RoleType)
	}

	if role.ProjectId == "" {
		return warnings, errors.New(errEmptyProjectId)
	}

	defaultLeaseTTL := sys.DefaultLeaseTTL()
	if role.TTL > defaultLeaseTTL {
		warnings = append(warnings, fmt.Sprintf(
			"Given ttl of %d seconds greater than current mount/system default of %d seconds; ttl will be capped at login time",
			role.TTL/time.Second, defaultLeaseTTL/time.Second))
	}

	defaultMaxTTL := sys.MaxLeaseTTL()
	if role.MaxTTL > defaultMaxTTL {
		warnings = append(warnings, fmt.Sprintf(
			"Given max_ttl of %d seconds greater than current mount/system default of %d seconds; max_ttl will be capped at login time",
			role.MaxTTL/time.Second, defaultMaxTTL/time.Second))
	}
	if role.MaxTTL < time.Duration(0) {
		return warnings, errors.New("max_ttl cannot be negative")
	}
	if role.MaxTTL != 0 && role.MaxTTL < role.TTL {
		return warnings, errors.New("ttl should be shorter than max_ttl")
	}

	if role.Period > sys.MaxLeaseTTL() {
		return warnings, fmt.Errorf("'period' of '%s' is greater than the backend's maximum lease TTL of '%s'", role.Period.String(), sys.MaxLeaseTTL().String())
	}

	return warnings, nil
}

// updateIamFields updates IAM-only fields for a role.
func (role *gcpRole) updateIamFields(data *framework.FieldData, op logical.Operation) error {
	allowGCEInference, ok := data.GetOk("allow_gce_inference")
	if ok {
		role.AllowGCEInference = allowGCEInference.(bool)
	} else if op == logical.CreateOperation {
		role.AllowGCEInference = data.Get("allow_gce_inference").(bool)
	}

	maxJwtExp, ok := data.GetOk("max_jwt_exp")
	if ok {
		role.MaxJwtExp = time.Duration(maxJwtExp.(int)) * time.Second
	} else if op == logical.CreateOperation {
		role.MaxJwtExp = time.Duration(defaultIamMaxJwtExpMinutes) * time.Minute
	}

	return nil
}

// updateGceFields updates GCE-only fields for a role.
func (role *gcpRole) updateGceFields(data *framework.FieldData, op logical.Operation) error {
	region, hasRegion := data.GetOk("bound_region")
	if hasRegion {
		role.BoundRegion = region.(string)
	}

	zone, hasZone := data.GetOk("bound_zone")
	if hasZone {
		role.BoundZone = zone.(string)
	}

	instanceGroup, ok := data.GetOk("bound_instance_group")
	if ok {
		role.BoundInstanceGroup = instanceGroup.(string)
	}

	labels, ok := data.GetOk("bound_labels")
	if ok {
		var invalidLabels []string
		role.BoundLabels, invalidLabels = util.ParseGcpLabels(labels.([]string))
		if len(invalidLabels) > 0 {
			return fmt.Errorf("invalid labels given: %q", invalidLabels)
		}
	}

	return nil
}

// validateIamFields validates the IAM-only fields for a role.
func (role *gcpRole) validateForIAM() (warnings []string, err error) {
	if len(role.BoundServiceAccounts) == 0 {
		return []string{}, errors.New(errEmptyIamServiceAccounts)
	}

	if len(role.BoundServiceAccounts) > 1 && strutil.StrListContains(role.BoundServiceAccounts, serviceAccountsWildcard) {
		return []string{}, fmt.Errorf("cannot provide IAM service account wildcard '%s' (for all service accounts) with other service accounts", serviceAccountsWildcard)
	}

	maxMaxJwtExp := time.Duration(maxJwtExpMaxMinutes) * time.Minute
	if role.MaxJwtExp > maxMaxJwtExp {
		return warnings, fmt.Errorf("max_jwt_exp cannot be more than %d minutes", maxJwtExpMaxMinutes)
	}

	return []string{}, nil
}

// validateGceFields validates the GCE-only fields for a role.
func (role *gcpRole) validateForGCE() (warnings []string, err error) {
	warnings = []string{}

	hasRegion := len(role.BoundRegion) > 0
	hasZone := len(role.BoundZone) > 0
	hasRegionOrZone := hasRegion || hasZone

	hasInstanceGroup := len(role.BoundInstanceGroup) > 0

	if hasInstanceGroup && !hasRegionOrZone {
		return warnings, errors.New(`region or zone information must be specified if a group is given`)
	}

	if hasRegion && hasZone {
		warnings = append(warnings, "Given both region and zone for role of type 'gce' - region will be ignored.")
	}

	return warnings, nil
}

// checkInvalidRoleTypeArgs checks that the data provided does not contain arguments
// for a different role type. If it does find some, it will return an error with the
// invalid args.
func checkInvalidRoleTypeArgs(data *framework.FieldData, invalidSchema map[string]*framework.FieldSchema) error {
	invalidArgs := []string{}

	for k := range data.Raw {
		if _, ok := baseRoleFieldSchema[k]; ok {
			continue
		}
		if _, ok := invalidSchema[k]; ok {
			invalidArgs = append(invalidArgs, k)
		}
	}

	if len(invalidArgs) > 0 {
		return fmt.Errorf(errTemplateInvalidRoleTypeArgs, data.Get("type"), strings.Join(invalidArgs, ","))
	}
	return nil
}
