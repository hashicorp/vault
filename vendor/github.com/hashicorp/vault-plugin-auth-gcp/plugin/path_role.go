package gcpauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	vaultconsts "github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
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

func baseRoleFieldSchema() map[string]*framework.FieldSchema {
	d := map[string]*framework.FieldSchema{
		"name": {
			Type:        framework.TypeString,
			Description: "Name of the role.",
		},
		"type": {
			Type:        framework.TypeString,
			Description: "Type of the role. Currently supported: iam, gce",
		},
		// Token Limits
		"policies": {
			Type:        framework.TypeCommaStringSlice,
			Description: tokenutil.DeprecationText("token_policies"),
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
	tokenutil.AddTokenFields(d)
	return d
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
	for k, v := range baseRoleFieldSchema() {
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

	respData := make(map[string]interface{})
	role.PopulateTokenData(respData)

	respData["role_id"] = role.RoleID

	if role.RoleType != "" {
		respData["type"] = role.RoleType
	}
	if len(role.BoundServiceAccounts) > 0 {
		respData["bound_service_accounts"] = role.BoundServiceAccounts
	}
	if len(role.BoundProjects) > 0 {
		respData["bound_projects"] = role.BoundProjects
	}
	respData["add_group_aliases"] = role.AddGroupAliases

	switch role.RoleType {
	case iamRoleType:
		if role.MaxJwtExp != 0 {
			respData["max_jwt_exp"] = int64(role.MaxJwtExp.Seconds())
		}
		respData["allow_gce_inference"] = role.AllowGCEInference
	case gceRoleType:
		if len(role.BoundRegions) > 0 {
			respData["bound_regions"] = role.BoundRegions
		}
		if len(role.BoundZones) > 0 {
			respData["bound_zones"] = role.BoundZones
		}
		if len(role.BoundInstanceGroups) > 0 {
			respData["bound_instance_groups"] = role.BoundInstanceGroups
		}
		if len(role.BoundLabels) > 0 {
			respData["bound_labels"] = role.BoundLabels
		}
	}

	// Upgrade vals
	if len(role.Policies) > 0 {
		respData["policies"] = respData["token_policies"]
	}
	if role.TTL > 0 {
		respData["ttl"] = int64(role.TTL.Seconds())
	}
	if role.MaxTTL > 0 {
		respData["max_ttl"] = int64(role.MaxTTL.Seconds())
	}
	if role.Period > 0 {
		respData["period"] = int64(role.Period.Seconds())
	}

	resp := &logical.Response{
		Data: respData,
	}
	return resp, nil
}

func (b *GcpAuthBackend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Validate we didn't get extraneous fields
	if err := validateFields(req, data); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
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

	if role.RoleID == "" {
		roleID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, logical.CodedError(http.StatusInternalServerError, fmt.Sprintf("unable to generate roleID: %s", err))
		}
		role.RoleID = roleID
	}

	warnings, err := role.updateRole(b.System(), req, data)
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
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
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
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
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

// role from storage. This assumes the caller has already obtained the role lock.
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

	// Upgrade token role params
	if role.TokenTTL == 0 && role.TTL > 0 {
		role.TokenTTL = role.TTL
		modified = true
	}
	if role.TokenMaxTTL == 0 && role.MaxTTL > 0 {
		role.TokenMaxTTL = role.MaxTTL
		modified = true
	}
	if role.TokenPeriod == 0 && role.Period > 0 {
		role.TokenPeriod = role.Period
		modified = true
	}
	if len(role.TokenPolicies) == 0 && len(role.Policies) > 0 {
		role.TokenPolicies = role.Policies
		modified = true
	}

	// Ensure the role has a RoleID
	if role.RoleID == "" {
		roleID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, logical.CodedError(http.StatusInternalServerError, fmt.Sprintf("unable to generate roleID for role missing an ID: %s", err))
		}
		role.RoleID = roleID
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

// checkInvalidRoleTypeArgs checks that the data provided does not contain arguments
// for a different role type. If it does find some, it will return an error with the
// invalid args.
func checkInvalidRoleTypeArgs(data *framework.FieldData, invalidSchema map[string]*framework.FieldSchema) error {
	invalidArgs := []string{}

	for k := range data.Raw {
		if _, ok := baseRoleFieldSchema()[k]; ok {
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
