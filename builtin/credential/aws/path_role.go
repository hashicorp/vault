package awsec2

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
		Pattern: "role/" + framework.GenericNameRegex("role"),
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
			"bound_ami_id": {
				Type: framework.TypeString,
				Description: `If set, defines a constraint on the EC2 instances that they should be
using the AMI ID specified by this parameter.`,
			},
			"bound_account_id": {
				Type: framework.TypeString,
				Description: `If set, defines a constraint on the EC2 instances that the account ID
in its identity document to match the one specified by this parameter.`,
			},
			"bound_iam_role_arn": {
				Type: framework.TypeString,
				Description: `If set, defines a constraint on the authenticating EC2 instance
that it must match the IAM role ARN specified by this parameter.
The value is prefix-matched (as though it were a glob ending in
'*').  The configured IAM user or EC2 instance role must be allowed
to execute the 'iam:GetInstanceProfile' action if this is
specified.`,
			},
			"bound_iam_instance_profile_arn": {
				Type: framework.TypeString,
				Description: `If set, defines a constraint on the EC2 instances to be associated
with an IAM instance profile ARN which has a prefix that matches
the value specified by this parameter. The value is prefix-matched
(as though it were a glob ending in '*').`,
			},
			"role_tag": {
				Type:        framework.TypeString,
				Default:     "",
				Description: "If set, enables the role tags for this role. The value set for this field should be the 'key' of the tag on the EC2 instance. The 'value' of the tag should be generated using 'role/<role>/tag' endpoint. Defaults to an empty string, meaning that role tags are disabled.",
			},
			"ttl": {
				Type:    framework.TypeDurationSecond,
				Default: 0,
				Description: `Duration in seconds after which the issued token should expire. Defaults
to 0, in which case the value will fallback to the system/mount defaults.`,
			},
			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Default:     0,
				Description: "The maximum allowed lifetime of tokens issued using this role.",
			},
			"policies": {
				Type:        framework.TypeString,
				Default:     "default",
				Description: "Policies to be set on tokens issued using this role.",
			},
			"allow_instance_migration": {
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set, allows migration of the underlying instance where the client resides. This keys off of pendingTime in the metadata document, so essentially, this disables the client nonce check whenever the instance is migrated to a new host and pendingTime is newer than the previously-remembered time. Use with caution.",
			},
			"disallow_reauthentication": {
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set, only allows a single token to be granted per instance ID. In order to perform a fresh login, the entry in whitelist for the instance ID needs to be cleared using 'auth/aws-ec2/identity-whitelist/<instance_id>' endpoint.",
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

func pathListRole(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

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
	entry, err := b.lockedAWSRole(req.Storage, strings.ToLower(data.Get("role").(string)))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// lockedAWSRole returns the properties set on the given role. This method
// acquires the read lock before reading the role from the storage.
func (b *backend) lockedAWSRole(s logical.Storage, roleName string) (*awsRoleEntry, error) {
	if roleName == "" {
		return nil, fmt.Errorf("missing role name")
	}

	b.roleMutex.RLock()
	defer b.roleMutex.RUnlock()

	return b.nonLockedAWSRole(s, roleName)
}

// lockedSetAWSRole creates or updates a role in the storage. This method
// acquires the write lock before creating or updating the role at the storage.
func (b *backend) lockedSetAWSRole(s logical.Storage, roleName string, roleEntry *awsRoleEntry) error {
	if roleName == "" {
		return fmt.Errorf("missing role name")
	}

	if roleEntry == nil {
		return fmt.Errorf("nil role entry")
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	return b.nonLockedSetAWSRole(s, roleName, roleEntry)
}

// nonLockedSetAWSRole creates or updates a role in the storage. This method
// does not acquire the write lock before reading the role from the storage. If
// locking is desired, use lockedSetAWSRole instead.
func (b *backend) nonLockedSetAWSRole(s logical.Storage, roleName string,
	roleEntry *awsRoleEntry) error {
	if roleName == "" {
		return fmt.Errorf("missing role name")
	}

	if roleEntry == nil {
		return fmt.Errorf("nil role entry")
	}

	entry, err := logical.StorageEntryJSON("role/"+strings.ToLower(roleName), roleEntry)
	if err != nil {
		return err
	}

	if err := s.Put(entry); err != nil {
		return err
	}

	return nil
}

// nonLockedAWSRole returns the properties set on the given role. This method
// does not acquire the read lock before reading the role from the storage. If
// locking is desired, use lockedAWSRole instead.
func (b *backend) nonLockedAWSRole(s logical.Storage, roleName string) (*awsRoleEntry, error) {
	if roleName == "" {
		return nil, fmt.Errorf("missing role name")
	}

	entry, err := s.Get("role/" + strings.ToLower(roleName))
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

	// Check if the value held by role ARN field is actually an instance profile ARN
	if result.BoundIamRoleARN != "" && strings.Contains(result.BoundIamRoleARN, ":instance-profile/") {
		// If yes, move it to the correct field
		result.BoundIamInstanceProfileARN = result.BoundIamRoleARN

		// Reset the old field
		result.BoundIamRoleARN = ""

		// Save the update
		if err = b.nonLockedSetAWSRole(s, roleName, &result); err != nil {
			return nil, fmt.Errorf("failed to move instance profile ARN to bound_iam_instance_profile_arn field")
		}
	}

	return &result, nil
}

// pathRoleDelete is used to delete the information registered for a given AMI ID.
func (b *backend) pathRoleDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	return nil, req.Storage.Delete("role/" + strings.ToLower(roleName))
}

// pathRoleList is used to list all the AMI IDs registered with Vault.
func (b *backend) pathRoleList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.roleMutex.RLock()
	defer b.roleMutex.RUnlock()

	roles, err := req.Storage.List("role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

// pathRoleRead is used to view the information registered for a given AMI ID.
func (b *backend) pathRoleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleEntry, err := b.lockedAWSRole(req.Storage, strings.ToLower(data.Get("role").(string)))
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

	// Display the ttl in seconds.
	respData["ttl"] = roleEntry.TTL / time.Second
	// Display the max_ttl in seconds.
	respData["max_ttl"] = roleEntry.MaxTTL / time.Second

	return &logical.Response{
		Data: respData,
	}, nil
}

// pathRoleCreateUpdate is used to associate Vault policies to a given AMI ID.
func (b *backend) pathRoleCreateUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	roleName := strings.ToLower(data.Get("role").(string))
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	roleEntry, err := b.nonLockedAWSRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		roleEntry = &awsRoleEntry{}
	}

	// Fetch and set the bound parameters. There can't be default values
	// for these.
	if boundAmiIDRaw, ok := data.GetOk("bound_ami_id"); ok {
		roleEntry.BoundAmiID = boundAmiIDRaw.(string)
	}

	if boundAccountIDRaw, ok := data.GetOk("bound_account_id"); ok {
		roleEntry.BoundAccountID = boundAccountIDRaw.(string)
	}

	if boundIamRoleARNRaw, ok := data.GetOk("bound_iam_role_arn"); ok {
		roleEntry.BoundIamRoleARN = boundIamRoleARNRaw.(string)
	}

	if boundIamInstanceProfileARNRaw, ok := data.GetOk("bound_iam_instance_profile_arn"); ok {
		roleEntry.BoundIamInstanceProfileARN = boundIamInstanceProfileARNRaw.(string)
	}

	// Ensure that at least one bound is set on the role
	switch {
	case roleEntry.BoundAccountID != "":
	case roleEntry.BoundAmiID != "":
	case roleEntry.BoundIamInstanceProfileARN != "":
	case roleEntry.BoundIamRoleARN != "":
	default:

		return logical.ErrorResponse("at least be one bound parameter should be specified on the role"), nil
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

	var resp logical.Response

	ttlRaw, ok := data.GetOk("ttl")
	if ok {
		ttl := time.Duration(ttlRaw.(int)) * time.Second
		defaultLeaseTTL := b.System().DefaultLeaseTTL()
		if ttl > defaultLeaseTTL {
			resp.AddWarning(fmt.Sprintf("Given ttl of %d seconds greater than current mount/system default of %d seconds; ttl will be capped at login time", ttl/time.Second, defaultLeaseTTL/time.Second))
		}
		roleEntry.TTL = ttl
	} else if req.Operation == logical.CreateOperation {
		roleEntry.TTL = time.Duration(data.Get("ttl").(int)) * time.Second
	}

	maxTTLInt, ok := data.GetOk("max_ttl")
	if ok {
		maxTTL := time.Duration(maxTTLInt.(int)) * time.Second
		systemMaxTTL := b.System().MaxLeaseTTL()
		if maxTTL > systemMaxTTL {
			resp.AddWarning(fmt.Sprintf("Given max_ttl of %d seconds greater than current mount/system default of %d seconds; max_ttl will be capped at login time", maxTTL/time.Second, systemMaxTTL/time.Second))
		}

		if maxTTL < time.Duration(0) {
			return logical.ErrorResponse("max_ttl cannot be negative"), nil
		}

		roleEntry.MaxTTL = maxTTL
	} else if req.Operation == logical.CreateOperation {
		roleEntry.MaxTTL = time.Duration(data.Get("max_ttl").(int)) * time.Second
	}

	if roleEntry.MaxTTL != 0 && roleEntry.MaxTTL < roleEntry.TTL {
		return logical.ErrorResponse("ttl should be shorter than max_ttl"), nil
	}

	roleTagStr, ok := data.GetOk("role_tag")
	if ok {
		roleEntry.RoleTag = roleTagStr.(string)
		// There is a limit of 127 characters on the tag key for AWS EC2 instances.
		// Complying to that requirement, do not allow the value of 'key' to be more than that.
		if len(roleEntry.RoleTag) > 127 {
			return logical.ErrorResponse("length of role tag exceeds the EC2 key limit of 127 characters"), nil
		}
	} else if req.Operation == logical.CreateOperation {
		roleEntry.RoleTag = data.Get("role_tag").(string)
	}

	if roleEntry.HMACKey == "" {
		roleEntry.HMACKey, err = uuid.GenerateUUID()
		if err != nil {
			return nil, fmt.Errorf("failed to generate role HMAC key: %v", err)
		}
	}

	if err := b.nonLockedSetAWSRole(req.Storage, roleName, roleEntry); err != nil {
		return nil, err
	}

	if len(resp.Warnings()) == 0 {
		return nil, nil
	}

	return &resp, nil
}

// Struct to hold the information associated with an AMI ID in Vault.
type awsRoleEntry struct {
	BoundAmiID                 string        `json:"bound_ami_id" structs:"bound_ami_id" mapstructure:"bound_ami_id"`
	BoundAccountID             string        `json:"bound_account_id" structs:"bound_account_id" mapstructure:"bound_account_id"`
	BoundIamRoleARN            string        `json:"bound_iam_role_arn" structs:"bound_iam_role_arn" mapstructure:"bound_iam_role_arn"`
	BoundIamInstanceProfileARN string        `json:"bound_iam_instance_profile_arn" structs:"bound_iam_instance_profile_arn" mapstructure:"bound_iam_instance_profile_arn"`
	RoleTag                    string        `json:"role_tag" structs:"role_tag" mapstructure:"role_tag"`
	AllowInstanceMigration     bool          `json:"allow_instance_migration" structs:"allow_instance_migration" mapstructure:"allow_instance_migration"`
	TTL                        time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`
	MaxTTL                     time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
	Policies                   []string      `json:"policies" structs:"policies" mapstructure:"policies"`
	DisallowReauthentication   bool          `json:"disallow_reauthentication" structs:"disallow_reauthentication" mapstructure:"disallow_reauthentication"`
	HMACKey                    string        `json:"hmac_key" structs:"hmac_key" mapstructure:"hmac_key"`
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
endpoint 'role/<role>/tag'. This tag then needs to be applied on the
instance before it attempts a login. The policies on the tag should be a
subset of policies that are associated to the role. In order to enable
login using tags, 'role_tag' option should be set while creating a role.

Also, a 'max_ttl' can be configured in this endpoint that determines the maximum
duration for which a login can be renewed. Note that the 'max_ttl' has an upper
limit of the 'max_ttl' value on the backend's mount.
`

const pathListRolesHelpSyn = `
Lists all the roles that are registered with Vault.
`

const pathListRolesHelpDesc = `
Roles will be listed by their respective role names.
`
