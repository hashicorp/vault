package gcpauth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-gcp-common/gcputil"
	vaultconsts "github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	// Role types
	iamRoleType = "iam"
	gceRoleType = "gce"

	// Errors
	errEmptyRoleName           = "role name is required"
	errEmptyRoleType           = "role type cannot be empty"
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

var baseRoleFieldSchema = map[string]*framework.FieldSchema{
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
	"bound_projects": {
		Type:        framework.TypeCommaStringSlice,
		Description: `GCP Projects that authenticating entities must belong to.`,
	},
	"bound_service_accounts": {
		Type: framework.TypeCommaStringSlice,
		Description: `
	Can be set for both 'iam' and 'gce' roles (required for 'iam'). A comma-seperated list of authorized service accounts.
	If the single value "*" is given, this is assumed to be all service accounts under the role's project. If this
	is set on a GCE role, the inferred service account from the instance metadata token will be used.`,
	},
	"add_group_aliases": {
		Type:    framework.TypeBool,
		Default: false,
		Description: "If true, will add group aliases to auth tokens generated under this role. " +
			"This will add the full list of ancestors (projects, folders, organizations) " +
			"for the given entity's project. Requires IAM permission `resourcemanager.projects.get` " +
			"on this project.",
	},
}

var iamOnlyFieldSchema = map[string]*framework.FieldSchema{
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

var gceOnlyFieldSchema = map[string]*framework.FieldSchema{
	"bound_zones": {
		Type: framework.TypeCommaStringSlice,
		Description: "Comma-separated list of permitted zones to which the GCE " +
			"instance must belong. If a group is provided, it is assumed to be a " +
			"zonal group. This can be a self-link or zone name. This option only " +
			"applies to \"gce\" roles.",
	},

	"bound_regions": {
		Type: framework.TypeCommaStringSlice,
		Description: "Comma-separated list of permitted regions to which the GCE " +
			"instance must belong. If a group is provided, it is assumed to be a " +
			"regional group. If \"zone\" is provided, this option is ignored. This " +
			"can be a self-link or region name. This option only applies to \"gce\" roles.",
	},

	"bound_instance_groups": {
		Type: framework.TypeCommaStringSlice,
		Description: "Comma-separated list of permitted instance groups to which " +
			"the GCE instance must belong. This option only applies to \"gce\" roles.",
	},

	"bound_labels": {
		Type: framework.TypeCommaStringSlice,
		Description: "Comma-separated list of GCP labels formatted as" +
			"\"key:value\" strings that must be present on the GCE instance " +
			"in order to authenticate. This option only applies to \"gce\" roles.",
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
	for k, v := range deprecatedFieldSchema {
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

func (b *GcpAuthBackend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.role(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *GcpAuthBackend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse(errEmptyRoleName), nil
	}

	if err := req.Storage.Delete(ctx, fmt.Sprintf("role/%s", name)); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *GcpAuthBackend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse(errEmptyRoleName), nil
	}

	role, err := b.role(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	} else if role == nil {
		return nil, nil
	}

	role.Period /= time.Second
	role.TTL /= time.Second
	role.MaxTTL /= time.Second
	role.MaxJwtExp /= time.Second

	resp := make(map[string]interface{})

	if role.RoleType != "" {
		resp["type"] = role.RoleType
	}
	if len(role.Policies) > 0 {
		resp["policies"] = role.Policies
	}
	if role.TTL != 0 {
		resp["ttl"] = role.TTL
	}
	if role.MaxTTL != 0 {
		resp["max_ttl"] = role.MaxTTL
	}
	if role.Period != 0 {
		resp["period"] = role.Period
	}
	if len(role.BoundServiceAccounts) > 0 {
		resp["bound_service_accounts"] = role.BoundServiceAccounts
	}
	if len(role.BoundProjects) > 0 {
		resp["bound_projects"] = role.BoundProjects
	}
	resp["add_group_aliases"] = role.AddGroupAliases

	switch role.RoleType {
	case iamRoleType:
		if role.MaxJwtExp != 0 {
			resp["max_jwt_exp"] = role.MaxJwtExp
		}
		resp["allow_gce_inference"] = role.AllowGCEInference
	case gceRoleType:
		if len(role.BoundRegions) > 0 {
			resp["bound_regions"] = role.BoundRegions
		}
		if len(role.BoundZones) > 0 {
			resp["bound_zones"] = role.BoundZones
		}
		if len(role.BoundInstanceGroups) > 0 {
			resp["bound_instance_groups"] = role.BoundInstanceGroups
		}
		if len(role.BoundLabels) > 0 {
			resp["bound_labels"] = role.BoundLabels
		}
	}

	return &logical.Response{
		Data: resp,
	}, nil
}

func (b *GcpAuthBackend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Validate we didn't get extraneous fields
	if err := validateFields(req, data); err != nil {
		return nil, logical.CodedError(422, err.Error())
	}

	name := strings.ToLower(data.Get("name").(string))
	if name == "" {
		return logical.ErrorResponse(errEmptyRoleName), nil
	}

	role, err := b.role(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		role = &gcpRole{}
	}

	warnings, err := role.updateRole(b.System(), req.Operation, data)
	if err != nil {
		resp := logical.ErrorResponse(err.Error())
		for _, w := range warnings {
			resp.AddWarning(w)
		}
		return resp, nil
	}
	return b.storeRole(ctx, req.Storage, name, role, warnings)
}

func (b *GcpAuthBackend) pathRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List(ctx, "role/")
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

func (b *GcpAuthBackend) pathRoleEditIamServiceAccounts(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Validate we didn't get extraneous fields
	if err := validateFields(req, data); err != nil {
		return nil, logical.CodedError(422, err.Error())
	}

	var warnings []string

	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse(errEmptyRoleName), nil
	}

	toAdd := data.Get("add").([]string)
	toRemove := data.Get("remove").([]string)
	if len(toAdd) == 0 && len(toRemove) == 0 {
		return logical.ErrorResponse("must provide at least one value to add or remove"), nil
	}

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	if role.RoleType != iamRoleType {
		return logical.ErrorResponse(fmt.Sprintf(errTemplateEditListWrongType, role.RoleType, "service_accounts", iamRoleType)), nil
	}
	role.BoundServiceAccounts = editStringValues(role.BoundServiceAccounts, toAdd, toRemove)

	return b.storeRole(ctx, req.Storage, roleName, role, warnings)
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

func (b *GcpAuthBackend) pathRoleEditGceLabels(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Validate we didn't get extraneous fields
	if err := validateFields(req, data); err != nil {
		return nil, logical.CodedError(422, err.Error())
	}

	var warnings []string

	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse(errEmptyRoleName), nil
	}

	toAdd := data.Get("add").([]string)
	toRemove := data.Get("remove").([]string)
	if len(toAdd) == 0 && len(toRemove) == 0 {
		return logical.ErrorResponse("must provide at least one value to add or remove"), nil
	}

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	if role.RoleType != gceRoleType {
		return logical.ErrorResponse(fmt.Sprintf(errTemplateEditListWrongType, role.RoleType, "labels", gceRoleType)), nil
	}

	labelsToAdd, invalidLabels := gcputil.ParseGcpLabels(toAdd)
	if len(invalidLabels) > 0 {
		return logical.ErrorResponse(fmt.Sprintf("given invalid labels to add: %q", invalidLabels)), nil
	}
	for k, v := range labelsToAdd {
		if role.BoundLabels == nil {
			role.BoundLabels = make(map[string]string, len(labelsToAdd))
		}
		role.BoundLabels[k] = v
	}

	for _, k := range toRemove {
		delete(role.BoundLabels, k)
	}

	return b.storeRole(ctx, req.Storage, roleName, role, warnings)
}

// role reads a gcpRole from storage. This assumes the caller has already obtained the role lock.
func (b *GcpAuthBackend) role(ctx context.Context, s logical.Storage, name string) (*gcpRole, error) {
	name = strings.ToLower(name)

	entry, err := s.Get(ctx, "role/"+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var role gcpRole
	if err := entry.DecodeJSON(&role); err != nil {
		return nil, err
	}

	// Keep track of whether we do an in-place upgrade of old fields
	modified := false

	// Move old bindings to new fields.
	if role.ProjectId != "" && len(role.BoundProjects) == 0 {
		role.BoundProjects = []string{role.ProjectId}
		role.ProjectId = ""
		modified = true
	}

	if role.BoundRegion != "" && len(role.BoundRegions) == 0 {
		role.BoundRegions = []string{role.BoundRegion}
		role.BoundRegion = ""
		modified = true
	}
	if role.BoundZone != "" && len(role.BoundZones) == 0 {
		role.BoundZones = []string{role.BoundZone}
		role.BoundZone = ""
		modified = true
	}
	if role.BoundInstanceGroup != "" && len(role.BoundInstanceGroups) == 0 {
		role.BoundInstanceGroups = []string{role.BoundInstanceGroup}
		role.BoundInstanceGroup = ""
		modified = true
	}

	if modified && (b.System().LocalMount() || !b.System().ReplicationState().HasState(vaultconsts.ReplicationPerformanceSecondary)) {
		b.Logger().Info("upgrading role to new schema",
			"role", name)

		updatedRole, err := logical.StorageEntryJSON("role/"+name, &role)
		if err != nil {
			return nil, err
		}
		if err := s.Put(ctx, updatedRole); err != nil {
			// Only perform upgrades on replication primary
			if !strings.Contains(err.Error(), logical.ErrReadOnly.Error()) {
				return nil, err
			}
		}
	}

	return &role, nil
}

// storeRole saves the gcpRole to storage.
// The returned response may contain either warnings or an error response,
// but will be nil if error is not nil
func (b *GcpAuthBackend) storeRole(ctx context.Context, s logical.Storage, roleName string, role *gcpRole, warnings []string) (*logical.Response, error) {
	var resp logical.Response
	for _, w := range warnings {
		resp.AddWarning(w)
	}

	validateWarnings, err := role.validate(b.System())
	for _, w := range validateWarnings {
		resp.AddWarning(w)
	}
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	entry, err := logical.StorageEntryJSON(fmt.Sprintf("role/%s", roleName), role)
	if err != nil {
		return nil, err
	}
	if err := s.Put(ctx, entry); err != nil {
		return nil, err
	}

	return &resp, nil
}

type gcpRole struct {
	// Type of this role. See path_role constants for currently supported types.
	RoleType string `json:"role_type,omitempty"`

	// Policies for Vault to assign to authorized entities.
	Policies []string `json:"policies,omitempty"`

	// TTL of Vault auth leases under this role.
	TTL time.Duration `json:"ttl,omitempty"`

	// Max total TTL including renewals, of Vault auth leases under this role.
	MaxTTL time.Duration `json:"max_ttl,omitempty"`

	// Period, If set, indicates that this token should not expire and
	// should be automatically renewed within this time period
	// with TTL equal to this value.
	Period time.Duration `json:"period,omitempty"`

	// Projects that entities must belong to
	BoundProjects []string `json:"bound_projects,omitempty"`

	// Service accounts allowed to login under this role.
	BoundServiceAccounts []string `json:"bound_service_accounts,omitempty"`

	// --| IAM-only attributes |--
	// MaxJwtExp is the duration from time of authentication that a JWT used to authenticate to role must expire within.
	// TODO(emilymye): Allow this to be updated for GCE roles once 'exp' parameter has been allowed for GCE metadata.
	MaxJwtExp time.Duration `json:"max_jwt_exp,omitempty"`

	// AllowGCEInference, if false, does not allow a GCE instance to login under this 'iam' role. If true (default),
	// a service account is inferred from the instance metadata and used as the authenticating instance.
	AllowGCEInference bool `json:"allow_gce_inference,omitempty"`

	// --| GCE-only attributes |--
	// BoundRegions that instances must belong to in order to login under this role.
	BoundRegions []string `json:"bound_regions,omitempty"`

	// BoundZones that instances must belong to in order to login under this role.
	BoundZones []string `json:"bound_zones,omitempty"`

	// BoundInstanceGroups are the instance group that instances must belong to in order to login under this role.
	BoundInstanceGroups []string `json:"bound_instance_groups,omitempty"`

	// BoundLabels that instances must currently have set in order to login under this role.
	BoundLabels map[string]string `json:"bound_labels,omitempty"`

	AddGroupAliases bool `json:"add_group_aliases,omitempty"`

	// Deprecated fields
	// TODO: Remove in 0.5.0+
	ProjectId          string `json:"project_id,omitempty"`
	BoundRegion        string `json:"bound_region,omitempty"`
	BoundZone          string `json:"bound_zone,omitempty"`
	BoundInstanceGroup string `json:"bound_instance_group,omitempty"`
}

// Update updates the given role with values parsed/validated from given FieldData.
// Exactly one of the response and error will be nil. The response is only used to pass back warnings.
// This method does not validate the role. Validation is done before storage.
func (role *gcpRole) updateRole(sys logical.SystemView, op logical.Operation, data *framework.FieldData) (warnings []string, err error) {
	// Set role type
	if rt, ok := data.GetOk("type"); ok {
		roleType := rt.(string)
		if role.RoleType != roleType && op == logical.UpdateOperation {
			err = errors.New("role type cannot be changed for an existing role")
			return
		}
		role.RoleType = roleType
	} else if op == logical.CreateOperation {
		err = errors.New(errEmptyRoleType)
		return
	}

	// Update policies
	if policies, ok := data.GetOk("policies"); ok {
		role.Policies = policyutil.ParsePolicies(policies)
	} else if op == logical.CreateOperation {
		// Force default policy
		role.Policies = policyutil.ParsePolicies(nil)
	}

	// Update token TTL.
	if ttl, ok := data.GetOk("ttl"); ok {
		role.TTL = time.Duration(ttl.(int)) * time.Second

		def := sys.DefaultLeaseTTL()
		if role.TTL > def {
			warnings = append(warnings, fmt.Sprintf(`Given "ttl" of %q is greater `+
				`than the maximum system/mount TTL of %q. The TTL will be capped at `+
				`%q during login.`, role.TTL, def, def))
		}
	} else if op == logical.CreateOperation {
		role.TTL = time.Duration(data.Get("ttl").(int)) * time.Second
	}

	// Update token Max TTL.
	if maxTTL, ok := data.GetOk("max_ttl"); ok {
		role.MaxTTL = time.Duration(maxTTL.(int)) * time.Second

		def := sys.MaxLeaseTTL()
		if role.MaxTTL > def {
			warnings = append(warnings, fmt.Sprintf(`Given "max_ttl" of %q is greater `+
				`than the maximum system/mount MaxTTL of %q. The MaxTTL will be `+
				`capped at %q during login.`, role.MaxTTL, def, def))
		}
	} else if op == logical.CreateOperation {
		role.MaxTTL = time.Duration(data.Get("max_ttl").(int)) * time.Second
	}

	// Update token period.
	if period, ok := data.GetOk("period"); ok {
		role.Period = time.Duration(period.(int)) * time.Second

		def := sys.MaxLeaseTTL()
		if role.Period > def {
			warnings = append(warnings, fmt.Sprintf(`Given "period" of %q is greater `+
				`than the maximum system/mount period of %q. The period will be `+
				`capped at %q during login.`, role.Period, def, def))
		}
	} else if op == logical.CreateOperation {
		role.Period = time.Duration(data.Get("period").(int)) * time.Second
	}

	// Update bound GCP service accounts.
	if sa, ok := data.GetOk("bound_service_accounts"); ok {
		role.BoundServiceAccounts = sa.([]string)
	} else {
		// Check for older version of param name
		if sa, ok := data.GetOk("service_accounts"); ok {
			warnings = append(warnings, `The "service_accounts" field is deprecated. `+
				`Please use "bound_service_accounts" instead. The "service_accounts" `+
				`field will be removed in a later release, so please update accordingly.`)
			role.BoundServiceAccounts = sa.([]string)
		}
	}
	if len(role.BoundServiceAccounts) > 0 {
		role.BoundServiceAccounts = strutil.TrimStrings(role.BoundServiceAccounts)
		role.BoundServiceAccounts = strutil.RemoveDuplicates(role.BoundServiceAccounts, false)
	}

	// Update bound GCP projects.
	boundProjects, givenBoundProj := data.GetOk("bound_projects")
	if givenBoundProj {
		role.BoundProjects = boundProjects.([]string)
	}
	if projectId, ok := data.GetOk("project_id"); ok {
		if givenBoundProj {
			return warnings, errors.New("only one of 'bound_projects' or 'project_id' can be given")
		}
		warnings = append(warnings,
			`The "project_id" (singular) field is deprecated. `+
				`Please use plural "bound_projects" instead to bind required GCP projects. `+
				`The "project_id" field will be removed in a later release, so please update accordingly.`)
		role.BoundProjects = []string{projectId.(string)}
	}
	if len(role.BoundProjects) > 0 {
		role.BoundProjects = strutil.TrimStrings(role.BoundProjects)
		role.BoundProjects = strutil.RemoveDuplicates(role.BoundProjects, false)
	}

	// Update bound GCP projects.
	addGroupAliases, ok := data.GetOk("add_group_aliases")
	if ok {
		role.AddGroupAliases = addGroupAliases.(bool)
	}

	// Update fields specific to this type
	switch role.RoleType {
	case iamRoleType:
		if err = checkInvalidRoleTypeArgs(data, gceOnlyFieldSchema); err != nil {
			return
		}
		if warnings, err = role.updateIamFields(data, op); err != nil {
			return
		}
	case gceRoleType:
		if err = checkInvalidRoleTypeArgs(data, iamOnlyFieldSchema); err != nil {
			return
		}
		if warnings, err = role.updateGceFields(data, op); err != nil {
			return
		}
	}

	return
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
func (role *gcpRole) updateIamFields(data *framework.FieldData, op logical.Operation) (warnings []string, err error) {
	if allowGCEInference, ok := data.GetOk("allow_gce_inference"); ok {
		role.AllowGCEInference = allowGCEInference.(bool)
	} else if op == logical.CreateOperation {
		role.AllowGCEInference = data.Get("allow_gce_inference").(bool)
	}

	if maxJwtExp, ok := data.GetOk("max_jwt_exp"); ok {
		role.MaxJwtExp = time.Duration(maxJwtExp.(int)) * time.Second
	} else if op == logical.CreateOperation {
		role.MaxJwtExp = time.Duration(defaultIamMaxJwtExpMinutes) * time.Minute
	}

	return
}

// updateGceFields updates GCE-only fields for a role.
func (role *gcpRole) updateGceFields(data *framework.FieldData, op logical.Operation) (warnings []string, err error) {
	if regions, ok := data.GetOk("bound_regions"); ok {
		role.BoundRegions = regions.([]string)
	} else if op == logical.CreateOperation {
		role.BoundRegions = data.Get("bound_regions").([]string)
	}

	if zones, ok := data.GetOk("bound_zones"); ok {
		role.BoundZones = zones.([]string)
	} else if op == logical.CreateOperation {
		role.BoundZones = data.Get("bound_zones").([]string)
	}

	if instanceGroups, ok := data.GetOk("bound_instance_groups"); ok {
		role.BoundInstanceGroups = instanceGroups.([]string)
	} else if op == logical.CreateOperation {
		role.BoundInstanceGroups = data.Get("bound_instance_groups").([]string)
	}

	if boundRegion, ok := data.GetOk("bound_region"); ok {
		if _, ok := data.GetOk("bound_regions"); ok {
			err = fmt.Errorf(`cannot specify both "bound_region" and "bound_regions"`)
			return
		}

		warnings = append(warnings, `The "bound_region" field is deprecated. `+
			`Please use "bound_regions" (plural) instead. You can still specify a `+
			`single region, but multiple regions are also now supported. The `+
			`"bound_region" field will be removed in a later release, so please `+
			`update accordingly.`)
		role.BoundRegions = append(role.BoundRegions, boundRegion.(string))
	}

	if boundZone, ok := data.GetOk("bound_zone"); ok {
		if _, ok := data.GetOk("bound_zones"); ok {
			err = fmt.Errorf(`cannot specify both "bound_zone" and "bound_zones"`)
			return
		}

		warnings = append(warnings, `The "bound_zone" field is deprecated. `+
			`Please use "bound_zones" (plural) instead. You can still specify a `+
			`single zone, but multiple zones are also now supported. The `+
			`"bound_zone" field will be removed in a later release, so please `+
			`update accordingly.`)
		role.BoundZones = append(role.BoundZones, boundZone.(string))
	}

	if boundInstanceGroup, ok := data.GetOk("bound_instance_group"); ok {
		if _, ok := data.GetOk("bound_instance_groups"); ok {
			err = fmt.Errorf(`cannot specify both "bound_instance_group" and "bound_instance_groups"`)
			return
		}

		warnings = append(warnings, `The "bound_instance_group" field is deprecated. `+
			`Please use "bound_instance_groups" (plural) instead. You can still specify a `+
			`single instance group, but multiple instance groups are also now supported. The `+
			`"bound_instance_group" field will be removed in a later release, so please `+
			`update accordingly.`)
		role.BoundInstanceGroups = append(role.BoundInstanceGroups, boundInstanceGroup.(string))
	}

	if labelsRaw, ok := data.GetOk("bound_labels"); ok {
		labels, invalidLabels := gcputil.ParseGcpLabels(labelsRaw.([]string))
		if len(invalidLabels) > 0 {
			err = fmt.Errorf("invalid labels given: %q", invalidLabels)
			return
		}
		role.BoundLabels = labels
	}

	if len(role.Policies) > 0 {
		role.Policies = strutil.TrimStrings(role.Policies)
		role.Policies = strutil.RemoveDuplicates(role.Policies, false)
	}

	if len(role.BoundRegions) > 0 {
		role.BoundRegions = strutil.TrimStrings(role.BoundRegions)
		role.BoundRegions = strutil.RemoveDuplicates(role.BoundRegions, false)
	}

	if len(role.BoundZones) > 0 {
		role.BoundZones = strutil.TrimStrings(role.BoundZones)
		role.BoundZones = strutil.RemoveDuplicates(role.BoundZones, false)
	}

	if len(role.BoundInstanceGroups) > 0 {
		role.BoundInstanceGroups = strutil.TrimStrings(role.BoundInstanceGroups)
		role.BoundInstanceGroups = strutil.RemoveDuplicates(role.BoundInstanceGroups, false)
	}

	return
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

	hasRegion := len(role.BoundRegions) > 0
	hasZone := len(role.BoundZones) > 0
	hasRegionOrZone := hasRegion || hasZone

	hasInstanceGroup := len(role.BoundInstanceGroups) > 0

	if hasInstanceGroup && !hasRegionOrZone {
		return warnings, errors.New(`region or zone information must be specified if an instance group is given`)
	}

	if hasRegion && hasZone {
		warnings = append(warnings, `Given both "bound_regions" and "bound_zones" `+
			`fields for role type "gce", "bound_regions" will be ignored in favor `+
			`of the more specific "bound_zones" field. To fix this warning, update `+
			`the role to remove either the "bound_regions" or "bound_zones" field.`)
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

// deprecatedFieldSchema contains the deprecated role attributes
var deprecatedFieldSchema = map[string]*framework.FieldSchema{
	"service_accounts": {
		Type:        framework.TypeCommaStringSlice,
		Description: "Deprecated: use \"bound_service_accounts\" instead.",
	},
	"project_id": {
		Type:        framework.TypeString,
		Description: "Deprecated: use \"bound_projects\" instead",
	},
	"bound_zone": {
		Type:        framework.TypeString,
		Description: "Deprecated: use \"bound_zones\" instead.",
	},
	"bound_region": {
		Type:        framework.TypeString,
		Description: "Deprecated: use \"bound_regions\" instead.",
	},
	"bound_instance_group": {
		Type:        framework.TypeString,
		Description: "Deprecated: use \"bound_instance_groups\" instead.",
	},
}
