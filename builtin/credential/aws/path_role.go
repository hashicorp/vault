package aws

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRole(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("role_name"),
		Fields: map[string]*framework.FieldSchema{
			"role_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},

			"bound_ami_id": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `If set, defines a constraint that the EC2 instances that are trying to
login, should be using the AMI ID specified by this parameter.
`,
			},

			"role_tag": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "If set, enables the RoleTag for this AMI. The value set for this field should be the 'key' of the tag on the EC2 instance. The 'value' of the tag should be generated using 'role/<role_name>/tag' endpoint. Defaults to empty string.",
			},

			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Default:     0,
				Description: "The maximum allowed lease duration.",
			},

			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "default",
				Description: "Policies to be associated with the role.",
			},

			"allow_instance_migration": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set, allows migration of the underlying instance where the client resides. This keys off of pendingTime in the metadata document, so essentially, this disables the client nonce check whenever the instance is migrated to a new host and pendingTime is newer than the previously-remembered time. Use with caution.",
			},

			"disallow_reauthentication": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set, only allows a single token to be granted per instance ID. In order to perform a fresh login, the entry in whitelist for the instance ID needs to be cleared using 'auth/aws/whitelist/identity/<instance_id>' endpoint.",
			},
		},

		ExistenceCheck: b.pathRoleExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathRoleCreateUpdate,
			logical.UpdateOperation: b.pathRoleCreateUpdate,
			logical.ReadOperation:   b.pathRoleRead,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRoleSyn,
		HelpDescription: pathRoleDesc,
	}
}

// pathListRoles creates a path that enables listing of all the AMIs that are
// registered with Vault.
func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathRoleExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := awsRole(req.Storage, strings.ToLower(data.Get("role_name").(string)))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// awsRole is used to get the information registered for the given AMI ID.
func awsRole(s logical.Storage, role string) (*awsRoleEntry, error) {
	entry, err := s.Get("role/" + strings.ToLower(role))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result awsRoleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// pathRoleDelete is used to delete the information registered for a given AMI ID.
func (b *backend) pathRoleDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	return nil, req.Storage.Delete("role/" + strings.ToLower(roleName))
}

// pathRoleList is used to list all the AMI IDs registered with Vault.
func (b *backend) pathRoleList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List("role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

// pathRoleRead is used to view the information registered for a given AMI ID.
func (b *backend) pathRoleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleEntry, err := awsRole(req.Storage, strings.ToLower(data.Get("role_name").(string)))
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return nil, nil
	}

	// Prepare the map of all the entries in the roleEntry.
	respData := structs.New(roleEntry).Map()

	// HMAC key belonging to the role should NOT be exported.
	delete(respData, "hmac_key")

	// Display the max_ttl in seconds.
	respData["max_ttl"] = roleEntry.MaxTTL / time.Second

	return &logical.Response{
		Data: respData,
	}, nil
}

// pathRoleCreateUpdate is used to associate Vault policies to a given AMI ID.
func (b *backend) pathRoleCreateUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	roleName := strings.ToLower(data.Get("role_name").(string))
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	roleEntry, err := awsRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		roleEntry = &awsRoleEntry{}
	}

	// Set the bound parameters only if they are supplied.
	// There are no default values for bound parameters.
	boundAmiIDStr, ok := data.GetOk("bound_ami_id")
	if ok {
		roleEntry.BoundAmiID = boundAmiIDStr.(string)
	}

	// At least one bound parameter should be set. Currently, only
	// 'bound_ami_id' is supported. Check if that is set.
	if roleEntry.BoundAmiID == "" {
		return logical.ErrorResponse("role is not bounded to any resource; set bound_ami_id"), nil
	}

	policiesStr, ok := data.GetOk("policies")
	if ok {
		roleEntry.Policies = policyutil.ParsePolicies(policiesStr.(string))
	} else if req.Operation == logical.CreateOperation {
		roleEntry.Policies = []string{"default"}
	}

	disallowReauthenticationBool, ok := data.GetOk("disallow_reauthentication")
	if ok {
		roleEntry.DisallowReauthentication = disallowReauthenticationBool.(bool)
	} else if req.Operation == logical.CreateOperation {
		roleEntry.DisallowReauthentication = data.Get("disallow_reauthentication").(bool)
	}

	allowInstanceMigrationBool, ok := data.GetOk("allow_instance_migration")
	if ok {
		roleEntry.AllowInstanceMigration = allowInstanceMigrationBool.(bool)
	} else if req.Operation == logical.CreateOperation {
		roleEntry.AllowInstanceMigration = data.Get("allow_instance_migration").(bool)
	}

	maxTTLInt, ok := data.GetOk("max_ttl")
	if ok {
		maxTTL := time.Duration(maxTTLInt.(int)) * time.Second
		systemMaxTTL := b.System().MaxLeaseTTL()
		if maxTTL > systemMaxTTL {
			return logical.ErrorResponse(fmt.Sprintf("Given TTL of %d seconds greater than current mount/system default of %d seconds", maxTTL/time.Second, systemMaxTTL/time.Second)), nil
		}

		if maxTTL < time.Duration(0) {
			return logical.ErrorResponse("max_ttl cannot be negative"), nil
		}

		roleEntry.MaxTTL = maxTTL
	} else if req.Operation == logical.CreateOperation {
		roleEntry.MaxTTL = time.Duration(data.Get("max_ttl").(int)) * time.Second
	}

	roleTagStr, ok := data.GetOk("role_tag")
	if ok {
		roleEntry.RoleTag = roleTagStr.(string)
		// There is a limit of 127 characters on the tag key for AWS EC2 instances.
		// Complying to that requirement, do not allow the value of 'key' to be more than that.
		if len(roleEntry.RoleTag) > 127 {
			return logical.ErrorResponse("role tag 'key' is exceeding the limit of 127 characters"), nil
		}
	} else if req.Operation == logical.CreateOperation {
		roleEntry.RoleTag = data.Get("role_tag").(string)
	}

	roleEntry.HMACKey, err = uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate uuid HMAC key: %v", err)
	}

	entry, err := logical.StorageEntryJSON("role/"+roleName, roleEntry)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

// Struct to hold the information associated with an AMI ID in Vault.
type awsRoleEntry struct {
	BoundAmiID               string        `json:"bound_ami_id" structs:"bound_ami_id" mapstructure:"bound_ami_id"`
	RoleTag                  string        `json:"role_tag" structs:"role_tag" mapstructure:"role_tag"`
	AllowInstanceMigration   bool          `json:"allow_instance_migration" structs:"allow_instance_migration" mapstructure:"allow_instance_migration"`
	MaxTTL                   time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
	Policies                 []string      `json:"policies" structs:"policies" mapstructure:"policies"`
	DisallowReauthentication bool          `json:"disallow_reauthentication" structs:"disallow_reauthentication" mapstructure:"disallow_reauthentication"`
	HMACKey                  string        `json:"hmac_key" structs:"hmac_key" mapstructure:"hmac_key"`
}

const pathRoleSyn = `
Create a role and associate policies to it.
`

const pathRoleDesc = `
A precondition for login is that a role should be created in the backend.
The login endpoint takes in the role name against which the instance
should be validated. After authenticating the instance, the authorization
for the instance to access Vault's resources is determined by the policies
that are associated to the role though this endpoint.

When the instances require only a subset of policies on the role, then
'role_tag' option on the role can be enabled to create a role tag via the
endpoint 'role/<role_name>/tag'. This tag then needs to be applied on the
instance before it attempts a login. The policies on the tag should be a
subset of policies that are associated to the role. In order to enable
login using tags, 'role_tag' option should be set while creating a role.

Also, a 'max_ttl' can be configured in this endpoint that determines the maximum
duration for which a login can be renewed. Note that the 'max_ttl' has a upper
limit of the 'max_ttl' value on the backend's mount.
`

const pathListRolesHelpSyn = `
Lists all the roles that are registered with Vault.
`

const pathListRolesHelpDesc = `
Roles will be listed by their respective role names.
`
