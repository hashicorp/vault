package azuresecrets

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-01-01-preview/authorization"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	rolesStoragePath = "roles"

	credentialTypeSP = 0
)

// roleEntry is a Vault role construct that maps to Azure roles or Applications
type roleEntry struct {
	CredentialType      int           `json:"credential_type"` // Reserved. Always SP at this time.
	AzureRoles          []*AzureRole  `json:"azure_roles"`
	AzureGroups         []*AzureGroup `json:"azure_groups"`
	ApplicationID       string        `json:"application_id"`
	ApplicationObjectID string        `json:"application_object_id"`
	TTL                 time.Duration `json:"ttl"`
	MaxTTL              time.Duration `json:"max_ttl"`
}

// AzureRole is an Azure Role (https://docs.microsoft.com/en-us/azure/role-based-access-control/overview) applied
// to a scope. RoleName and RoleID are both traits of the role. RoleID is the unique identifier, but RoleName is
// more useful to a human (thought it is not unique).
type AzureRole struct {
	RoleName string `json:"role_name"` // e.g. Owner
	RoleID   string `json:"role_id"`   // e.g. /subscriptions/e0a207b2-.../providers/Microsoft.Authorization/roleDefinitions/de139f84-...
	Scope    string `json:"scope"`     // e.g. /subscriptions/e0a207b2-...
}

// AzureGroup is an Azure Active Directory Group
// (https://docs.microsoft.com/en-us/azure/role-based-access-control/overview).
// GroupName and ObjectID are both traits of the group. ObjectID is the unique
// identifier, but GroupName is more useful to a human (though it is not
// unique).
type AzureGroup struct {
	GroupName string `json:"group_name"` // e.g. MyGroup
	ObjectID  string `json:"object_id"`  // e.g. 90820a30-352d-400f-89e5-2ca74ac14333
}

func pathsRole(b *azureSecretBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "roles/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the role.",
				},
				"application_object_id": {
					Type:        framework.TypeString,
					Description: "Application Object ID to use for static service principal credentials.",
				},
				"azure_roles": {
					Type:        framework.TypeString,
					Description: "JSON list of Azure roles to assign.",
				},
				"azure_groups": {
					Type:        framework.TypeString,
					Description: "JSON list of Azure groups to add the service principal to.",
				},
				"ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Default lease for generated credentials. If not set or set to 0, will use system default.",
				},
				"max_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Maximum time a service principal. If not set or set to 0, will use system default.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.pathRoleRead,
				logical.CreateOperation: b.pathRoleUpdate,
				logical.UpdateOperation: b.pathRoleUpdate,
				logical.DeleteOperation: b.pathRoleDelete,
			},
			HelpSynopsis:    roleHelpSyn,
			HelpDescription: roleHelpDesc,
			ExistenceCheck:  b.pathRoleExistenceCheck,
		},
		{
			Pattern: "roles/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathRoleList,
			},
			HelpSynopsis:    roleListHelpSyn,
			HelpDescription: roleListHelpDesc,
		},
	}

}

// pathRoleUpdate creates or updates Vault roles.
//
// Basic validity check are made to verify that the provided fields meet requirements
// for the given credential type.
//
// Dynamic Service Principal:
//   Azure roles are checked for existence. The Azure role lookup step will allow the
//   operator to provide a role name or ID. ID is unambigious and will be used if provided.
//   Given just role name, a search will be performed and if exactly one match is found,
//   that role will be used.

//   Azure groups are checked for existence. The Azure groups lookup step will allow the
//   operator to provide a groups name or ID. ID is unambigious and will be used if provided.
//   Given just group name, a search will be performed and if exactly one match is found,
//   that group will be used.
//
// Static Service Principal:
//   The provided Application Object ID is checked for existence.
func (b *azureSecretBackend) pathRoleUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var resp *logical.Response

	client, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	// load or create role
	name := d.Get("name").(string)
	role, err := getRole(ctx, name, req.Storage)
	if err != nil {
		return nil, errwrap.Wrapf("error reading role: {{err}}", err)
	}

	if role == nil {
		if req.Operation == logical.UpdateOperation {
			return nil, errors.New("role entry not found during update operation")
		}
		role = &roleEntry{
			CredentialType: credentialTypeSP,
		}
	}

	// load and validate TTLs
	if ttlRaw, ok := d.GetOk("ttl"); ok {
		role.TTL = time.Duration(ttlRaw.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		role.TTL = time.Duration(d.Get("ttl").(int)) * time.Second
	}

	if maxTTLRaw, ok := d.GetOk("max_ttl"); ok {
		role.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		role.MaxTTL = time.Duration(d.Get("max_ttl").(int)) * time.Second
	}

	if role.MaxTTL != 0 && role.TTL > role.MaxTTL {
		return logical.ErrorResponse("ttl cannot be greater than max_ttl"), nil
	}

	// update and verify Application Object ID if provided
	if appObjectID, ok := d.GetOk("application_object_id"); ok {
		role.ApplicationObjectID = appObjectID.(string)
	}

	if role.ApplicationObjectID != "" {
		app, err := client.provider.GetApplication(ctx, role.ApplicationObjectID)
		if err != nil {
			return nil, errwrap.Wrapf("error loading Application: {{err}}", err)
		}
		role.ApplicationID = to.String(app.AppID)
	}

	// Parse the Azure roles
	if roles, ok := d.GetOk("azure_roles"); ok {
		parsedRoles := make([]*AzureRole, 0) // non-nil to avoid a "missing roles" error later

		err := jsonutil.DecodeJSON([]byte(roles.(string)), &parsedRoles)
		if err != nil {
			return logical.ErrorResponse("error parsing Azure roles '%s': %s", roles.(string), err.Error()), nil
		}
		role.AzureRoles = parsedRoles
	}

	// Parse the Azure groups
	if groups, ok := d.GetOk("azure_groups"); ok {
		parsedGroups := make([]*AzureGroup, 0) // non-nil to avoid a "missing groups" error later

		err := jsonutil.DecodeJSON([]byte(groups.(string)), &parsedGroups)
		if err != nil {
			return logical.ErrorResponse("error parsing Azure groups '%s': %s", groups.(string), err.Error()), nil
		}
		role.AzureGroups = parsedGroups
	}

	// update and verify Azure roles, including looking up each role by ID or name.
	roleSet := make(map[string]bool)
	for _, r := range role.AzureRoles {
		var roleDef authorization.RoleDefinition
		if r.RoleID != "" {
			roleDef, err = client.provider.GetRoleByID(ctx, r.RoleID)
			if err != nil {
				if strings.Contains(err.Error(), "RoleDefinitionDoesNotExist") {
					return logical.ErrorResponse("no role found for role_id: '%s'", r.RoleID), nil
				}
				return nil, errwrap.Wrapf("unable to lookup Azure role: {{err}}", err)
			}
		} else {
			defs, err := client.findRoles(ctx, r.RoleName)
			if err != nil {
				return nil, errwrap.Wrapf("unable to lookup Azure role: {{err}}", err)
			}
			if l := len(defs); l == 0 {
				return logical.ErrorResponse("no role found for role_name: '%s'", r.RoleName), nil
			} else if l > 1 {
				return logical.ErrorResponse("multiple matches found for role_name: '%s'. Specify role by ID instead.", r.RoleName), nil
			}
			roleDef = defs[0]
		}

		roleDefID := to.String(roleDef.ID)
		roleDefName := to.String(roleDef.RoleName)

		r.RoleName, r.RoleID = roleDefName, roleDefID

		rsKey := r.RoleID + "||" + r.Scope
		if roleSet[rsKey] {
			return logical.ErrorResponse("duplicate role_id and scope: '%s', '%s'", r.RoleID, r.Scope), nil
		}
		roleSet[rsKey] = true
	}

	// update and verify Azure groups, including looking up each group by ID or name.
	groupSet := make(map[string]bool)
	for _, r := range role.AzureGroups {
		var groupDef graphrbac.ADGroup
		if r.ObjectID != "" {
			groupDef, err = client.provider.GetGroup(ctx, r.ObjectID)
			if err != nil {
				if strings.Contains(err.Error(), "Request_ResourceNotFound") {
					return logical.ErrorResponse("no group found for object_id: '%s'", r.ObjectID), nil
				}
				return nil, errwrap.Wrapf("unable to lookup Azure group: {{err}}", err)
			}
		} else {
			defs, err := client.findGroups(ctx, r.GroupName)
			if err != nil {
				return nil, errwrap.Wrapf("unable to lookup Azure group: {{err}}", err)
			}
			if l := len(defs); l == 0 {
				return logical.ErrorResponse("no group found for group_name: '%s'", r.GroupName), nil
			} else if l > 1 {
				return logical.ErrorResponse("multiple matches found for group_name: '%s'. Specify group by ObjectID instead.", r.GroupName), nil
			}
			groupDef = defs[0]
		}

		groupDefID := to.String(groupDef.ObjectID)
		groupDefName := to.String(groupDef.DisplayName)
		r.GroupName, r.ObjectID = groupDefName, groupDefID

		if groupSet[r.ObjectID] {
			return logical.ErrorResponse("duplicate object_id '%s'", r.ObjectID), nil
		}
		groupSet[r.ObjectID] = true
	}

	if role.ApplicationObjectID == "" && len(role.AzureRoles) == 0 && len(role.AzureGroups) == 0 {
		return logical.ErrorResponse("either Azure role definitions, group definitions, or an Application Object ID must be provided"), nil
	}

	// save role
	err = saveRole(ctx, req.Storage, role, name)
	if err != nil {
		return nil, errwrap.Wrapf("error storing role: {{err}}", err)
	}

	return resp, nil
}

func (b *azureSecretBackend) pathRoleRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var data = make(map[string]interface{})

	name := d.Get("name").(string)

	r, err := getRole(ctx, name, req.Storage)
	if err != nil {
		return nil, errwrap.Wrapf("error reading role: {{err}}", err)
	}

	if r == nil {
		return nil, nil
	}

	data["ttl"] = r.TTL / time.Second
	data["max_ttl"] = r.MaxTTL / time.Second
	data["azure_roles"] = r.AzureRoles
	data["azure_groups"] = r.AzureGroups
	data["application_object_id"] = r.ApplicationObjectID

	return &logical.Response{
		Data: data,
	}, nil
}

func (b *azureSecretBackend) pathRoleList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List(ctx, rolesStoragePath+"/")
	if err != nil {
		return nil, errwrap.Wrapf("error listing roles: {{err}}", err)
	}

	return logical.ListResponse(roles), nil
}

func (b *azureSecretBackend) pathRoleDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	err := req.Storage.Delete(ctx, fmt.Sprintf("%s/%s", rolesStoragePath, name))
	if err != nil {
		return nil, errwrap.Wrapf("error deleting role: {{err}}", err)
	}

	return nil, nil
}

func (b *azureSecretBackend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)

	role, err := getRole(ctx, name, req.Storage)
	if err != nil {
		return false, errwrap.Wrapf("error reading role: {{err}}", err)
	}

	return role != nil, nil
}

func saveRole(ctx context.Context, s logical.Storage, c *roleEntry, name string) error {
	entry, err := logical.StorageEntryJSON(fmt.Sprintf("%s/%s", rolesStoragePath, name), c)
	if err != nil {
		return err
	}

	return s.Put(ctx, entry)
}

func getRole(ctx context.Context, name string, s logical.Storage) (*roleEntry, error) {
	entry, err := s.Get(ctx, fmt.Sprintf("%s/%s", rolesStoragePath, name))
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	role := new(roleEntry)
	if err := entry.DecodeJSON(role); err != nil {
		return nil, err
	}
	return role, nil
}

const roleHelpSyn = "Manage the Vault roles used to generate Azure credentials."
const roleHelpDesc = `
This path allows you to read and write roles that are used to generate Azure login
credentials. These roles are associated with either an existing Application, or a set
of Azure roles, which are used to control permissions to Azure resources.

If the backend is mounted at "azure", you would create a Vault role at "azure/roles/my_role",
and request credentials from "azure/creds/my_role".

Each Vault role is configured with the standard ttl parameters and either an
Application Object ID or a combination of Azure groups to make the service
principal a member of, and Azure roles and scopes to assign the service
principal to. During the Vault role creation, any set Azure role, group, or
Object ID will be fetched and verified, and therefore must exist for the request
to succeed. When a user requests credentials against the Vault role, a new
password will be created for the Application if an Application Object ID was
configured. Otherwise, a new service principal will be created and the
configured set of Azure roles are assigned to it and it will be added to the
configured groups.
`
const roleListHelpSyn = `List existing roles.`
const roleListHelpDesc = `List existing roles by name.`
