package awsauth

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
			"auth_type": {
				Type: framework.TypeString,
				Description: `The auth_type permitted to authenticate to this role. Must be one of
iam or ec2 and cannot be changed after role creation.`,
			},
			"bound_ami_id": {
				Type: framework.TypeString,
				Description: `If set, defines a constraint on the EC2 instances that they should be
using the AMI ID specified by this parameter. This is only applicable when auth_type is ec2
or inferred_entity_type is ec2_instance.`,
			},
			"bound_account_id": {
				Type: framework.TypeString,
				Description: `If set, defines a constraint on the EC2 instances that the account ID
in its identity document to match the one specified by this parameter. This is only
applicable when auth_type is ec2 or inferred_entity_type is ec2_instance.`,
			},
			"bound_iam_principal_arn": {
				Type: framework.TypeString,
				Description: `ARN of the IAM principal to bind to this role. Only applicable when
auth_type is iam.`,
			},
			"bound_region": {
				Type: framework.TypeString,
				Description: `If set, defines a constraint on the EC2 instances that the region in
its identity document to match the one specified by this parameter. This is only
applicable when auth_type is ec2.`,
			},
			"bound_iam_role_arn": {
				Type: framework.TypeString,
				Description: `If set, defines a constraint on the authenticating EC2 instance
that it must match the IAM role ARN specified by this parameter.
The value is prefix-matched (as though it were a glob ending in
'*').  The configured IAM user or EC2 instance role must be allowed
to execute the 'iam:GetInstanceProfile' action if this is specified. This is
only applicable when auth_type is ec2 or inferred_entity_type is
ec2_instance.`,
			},
			"bound_iam_instance_profile_arn": {
				Type: framework.TypeString,
				Description: `If set, defines a constraint on the EC2 instances to be associated
with an IAM instance profile ARN which has a prefix that matches
the value specified by this parameter. The value is prefix-matched
(as though it were a glob ending in '*'). This is only applicable when
auth_type is ec2 or inferred_entity_type is ec2_instance.`,
			},
			"resolve_aws_unique_ids": {
				Type:    framework.TypeBool,
				Default: true,
				Description: `If set, resolve all AWS IAM ARNs into AWS's internal unique IDs.
When an IAM entity (e.g., user, role, or instance profile) is deleted, then all references
to it within the role will be invalidated, which prevents a new IAM entity from being created
with the same name and matching the role's IAM binds. Once set, this cannot be unset.`,
			},
			"inferred_entity_type": {
				Type: framework.TypeString,
				Description: `When auth_type is iam, the
AWS entity type to infer from the authenticated principal. The only supported
value is ec2_instance, which will extract the EC2 instance ID from the
authenticated role and apply the following restrictions specific to EC2
instances: bound_ami_id, bound_account_id, bound_iam_role_arn,
bound_iam_instance_profile_arn, bound_vpc_id, bound_subnet_id. The configured
EC2 client must be able to find the inferred instance ID in the results, and the
instance must be running. If unable to determine the EC2 instance ID or unable
to find the EC2 instance ID among running instances, then authentication will
fail.`,
			},
			"inferred_aws_region": {
				Type: framework.TypeString,
				Description: `When auth_type is iam and
inferred_entity_type is set, the region to assume the inferred entity exists in.`,
			},
			"bound_vpc_id": {
				Type: framework.TypeString,
				Description: `
If set, defines a constraint on the EC2 instance to be associated with the VPC
ID that matches the value specified by this parameter. This is only applicable
when auth_type is ec2 or inferred_entity_type is ec2_instance.`,
			},
			"bound_subnet_id": {
				Type: framework.TypeString,
				Description: `
If set, defines a constraint on the EC2 instance to be associated with the
subnet ID that matches the value specified by this parameter. This is only
applicable when auth_type is ec2 or inferred_entity_type is ec2_instance.`,
			},
			"role_tag": {
				Type:    framework.TypeString,
				Default: "",
				Description: `If set, enables the role tags for this role. The value set for this
field should be the 'key' of the tag on the EC2 instance. The 'value'
of the tag should be generated using 'role/<role>/tag' endpoint.
Defaults to an empty string, meaning that role tags are disabled. This
is only allowed if auth_type is ec2.`,
			},
			"period": &framework.FieldSchema{
				Type:    framework.TypeDurationSecond,
				Default: 0,
				Description: `
If set, indicates that the token generated using this role should never expire.
The token should be renewed within the duration specified by this value. At
each renewal, the token's TTL will be set to the value of this parameter.`,
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
				Type:        framework.TypeCommaStringSlice,
				Default:     "default",
				Description: "Policies to be set on tokens issued using this role.",
			},
			"allow_instance_migration": {
				Type:    framework.TypeBool,
				Default: false,
				Description: `If set, allows migration of the underlying instance where the client
resides. This keys off of pendingTime in the metadata document, so
essentially, this disables the client nonce check whenever the
instance is migrated to a new host and pendingTime is newer than the
previously-remembered time. Use with caution. This is only checked when
auth_type is ec2.`,
			},
			"disallow_reauthentication": {
				Type:    framework.TypeBool,
				Default: false,
				Description: `If set, only allows a single token to be granted per
        instance ID. In order to perform a fresh login, the entry in whitelist
        for the instance ID needs to be cleared using
        'auth/aws-ec2/identity-whitelist/<instance_id>' endpoint. This is only
        applicable when auth_type is ec2.`,
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
	roleEntry, err := b.nonLockedAWSRole(s, roleName)
	// we manually unlock rather than defer the unlock because we might need to grab
	// a read/write lock in the upgrade path
	b.roleMutex.RUnlock()
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return nil, nil
	}
	needUpgrade, err := b.upgradeRoleEntry(s, roleEntry)
	if err != nil {
		return nil, fmt.Errorf("error upgrading roleEntry: %v", err)
	}
	if needUpgrade {
		b.roleMutex.Lock()
		defer b.roleMutex.Unlock()
		// Now that we have a R/W lock, we need to re-read the role entry in case it was
		// written to between releasing the read lock and acquiring the write lock
		roleEntry, err = b.nonLockedAWSRole(s, roleName)
		if err != nil {
			return nil, err
		}
		// somebody deleted the role, so no use in putting it back
		if roleEntry == nil {
			return nil, nil
		}
		// now re-check to see if we need to upgrade
		if needUpgrade, err = b.upgradeRoleEntry(s, roleEntry); err != nil {
			return nil, fmt.Errorf("error upgrading roleEntry: %v", err)
		}
		if needUpgrade {
			if err = b.nonLockedSetAWSRole(s, roleName, roleEntry); err != nil {
				return nil, fmt.Errorf("error saving upgraded roleEntry: %v", err)
			}
		}
	}
	return roleEntry, nil
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

// If needed, updates the role entry and returns a bool indicating if it was updated
// (and thus needs to be persisted)
func (b *backend) upgradeRoleEntry(s logical.Storage, roleEntry *awsRoleEntry) (bool, error) {
	if roleEntry == nil {
		return false, fmt.Errorf("received nil roleEntry")
	}
	var upgraded bool
	// Check if the value held by role ARN field is actually an instance profile ARN
	if roleEntry.BoundIamRoleARN != "" && strings.Contains(roleEntry.BoundIamRoleARN, ":instance-profile/") {
		// If yes, move it to the correct field
		roleEntry.BoundIamInstanceProfileARN = roleEntry.BoundIamRoleARN

		// Reset the old field
		roleEntry.BoundIamRoleARN = ""

		upgraded = true
	}

	// Check if there was no pre-existing AuthType set (from older versions)
	if roleEntry.AuthType == "" {
		// then default to the original behavior of ec2
		roleEntry.AuthType = ec2AuthType
		upgraded = true
	}

	if roleEntry.AuthType == iamAuthType &&
		roleEntry.ResolveAWSUniqueIDs &&
		roleEntry.BoundIamPrincipalARN != "" &&
		roleEntry.BoundIamPrincipalID == "" &&
		!strings.HasSuffix(roleEntry.BoundIamPrincipalARN, "*") {
		principalId, err := b.resolveArnToUniqueIDFunc(s, roleEntry.BoundIamPrincipalARN)
		if err != nil {
			return false, err
		}
		roleEntry.BoundIamPrincipalID = principalId
		upgraded = true
	}

	return upgraded, nil

}

// nonLockedAWSRole returns the properties set on the given role. This method
// does not acquire the read lock before reading the role from the storage. If
// locking is desired, use lockedAWSRole instead.
// This method also does NOT check to see if a role upgrade is required. It is
// the responsibility of the caller to check if a role upgrade is required and,
// if so, to upgrade the role
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

	// Display all the durations in seconds
	respData["ttl"] = roleEntry.TTL / time.Second
	respData["max_ttl"] = roleEntry.MaxTTL / time.Second
	respData["period"] = roleEntry.Period / time.Second

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
	} else {
		needUpdate, err := b.upgradeRoleEntry(req.Storage, roleEntry)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to update roleEntry: %v", err)), nil
		}
		if needUpdate {
			err = b.nonLockedSetAWSRole(req.Storage, roleName, roleEntry)
			if err != nil {
				return logical.ErrorResponse(fmt.Sprintf("failed to save upgraded roleEntry: %v", err)), nil
			}
		}
	}

	// Fetch and set the bound parameters. There can't be default values
	// for these.
	if boundAmiIDRaw, ok := data.GetOk("bound_ami_id"); ok {
		roleEntry.BoundAmiID = boundAmiIDRaw.(string)
	}

	if boundAccountIDRaw, ok := data.GetOk("bound_account_id"); ok {
		roleEntry.BoundAccountID = boundAccountIDRaw.(string)
	}

	if boundRegionRaw, ok := data.GetOk("bound_region"); ok {
		roleEntry.BoundRegion = boundRegionRaw.(string)
	}

	if boundVpcIDRaw, ok := data.GetOk("bound_vpc_id"); ok {
		roleEntry.BoundVpcID = boundVpcIDRaw.(string)
	}

	if boundSubnetIDRaw, ok := data.GetOk("bound_subnet_id"); ok {
		roleEntry.BoundSubnetID = boundSubnetIDRaw.(string)
	}

	if resolveAWSUniqueIDsRaw, ok := data.GetOk("resolve_aws_unique_ids"); ok {
		switch {
		case req.Operation == logical.CreateOperation:
			roleEntry.ResolveAWSUniqueIDs = resolveAWSUniqueIDsRaw.(bool)
		case roleEntry.ResolveAWSUniqueIDs && !resolveAWSUniqueIDsRaw.(bool):
			return logical.ErrorResponse("changing resolve_aws_unique_ids from true to false is not allowed"), nil
		default:
			roleEntry.ResolveAWSUniqueIDs = resolveAWSUniqueIDsRaw.(bool)
		}
	} else if req.Operation == logical.CreateOperation {
		roleEntry.ResolveAWSUniqueIDs = data.Get("resolve_aws_unique_ids").(bool)
	}

	if boundIamRoleARNRaw, ok := data.GetOk("bound_iam_role_arn"); ok {
		roleEntry.BoundIamRoleARN = boundIamRoleARNRaw.(string)
	}

	if boundIamInstanceProfileARNRaw, ok := data.GetOk("bound_iam_instance_profile_arn"); ok {
		roleEntry.BoundIamInstanceProfileARN = boundIamInstanceProfileARNRaw.(string)
	}

	if boundIamPrincipalARNRaw, ok := data.GetOk("bound_iam_principal_arn"); ok {
		principalARN := boundIamPrincipalARNRaw.(string)
		roleEntry.BoundIamPrincipalARN = principalARN
		// Explicitly not checking to see if the user has changed the ARN under us
		// This allows the user to sumbit an update with the same ARN to force Vault
		// to re-resolve the ARN to the unique ID, in case an entity was deleted and
		// recreated
		if roleEntry.ResolveAWSUniqueIDs && !strings.HasSuffix(roleEntry.BoundIamPrincipalARN, "*") {
			principalID, err := b.resolveArnToUniqueIDFunc(req.Storage, principalARN)
			if err != nil {
				return logical.ErrorResponse(fmt.Sprintf("failed updating the unique ID of ARN %#v: %#v", principalARN, err)), nil
			}
			roleEntry.BoundIamPrincipalID = principalID
		} else {
			// Need to handle the case where we're switching from a non-wildcard principal to a wildcard principal
			roleEntry.BoundIamPrincipalID = ""
		}
	} else if roleEntry.ResolveAWSUniqueIDs && roleEntry.BoundIamPrincipalARN != "" && !strings.HasSuffix(roleEntry.BoundIamPrincipalARN, "*") {
		// we're turning on resolution on this role, so ensure we update it
		principalID, err := b.resolveArnToUniqueIDFunc(req.Storage, roleEntry.BoundIamPrincipalARN)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("unable to resolve ARN %#v to internal ID: %#v", roleEntry.BoundIamPrincipalARN, err)), nil
		}
		roleEntry.BoundIamPrincipalID = principalID
	}

	if inferRoleTypeRaw, ok := data.GetOk("inferred_entity_type"); ok {
		roleEntry.InferredEntityType = inferRoleTypeRaw.(string)
	}

	if inferredAWSRegionRaw, ok := data.GetOk("inferred_aws_region"); ok {
		roleEntry.InferredAWSRegion = inferredAWSRegionRaw.(string)
	}

	// auth_type is a special case as it's immutable and can't be changed once a role is created
	if authTypeRaw, ok := data.GetOk("auth_type"); ok {
		// roleEntry.AuthType should only be "" when it's a new role; existing roles without an
		// auth_type should have already been upgraded to have one before we get here
		if roleEntry.AuthType == "" {
			switch authTypeRaw.(string) {
			case ec2AuthType, iamAuthType:
				roleEntry.AuthType = authTypeRaw.(string)
			default:
				return logical.ErrorResponse(fmt.Sprintf("unrecognized auth_type: %v", authTypeRaw.(string))), nil
			}
		} else if authTypeRaw.(string) != roleEntry.AuthType {
			return logical.ErrorResponse("changing auth_type on a role is not allowed"), nil
		}
	} else if req.Operation == logical.CreateOperation {
		switch req.MountType {
		// maintain backwards compatibility for old aws-ec2 auth types
		case "aws-ec2":
			roleEntry.AuthType = ec2AuthType
		// but default to iamAuth for new mounts going forward
		case "aws":
			roleEntry.AuthType = iamAuthType
		default:
			roleEntry.AuthType = iamAuthType
		}
	}

	allowEc2Binds := roleEntry.AuthType == ec2AuthType

	if roleEntry.InferredEntityType != "" {
		switch {
		case roleEntry.AuthType != iamAuthType:
			return logical.ErrorResponse("specified inferred_entity_type but didn't allow iam auth_type"), nil
		case roleEntry.InferredEntityType != ec2EntityType:
			return logical.ErrorResponse(fmt.Sprintf("specified invalid inferred_entity_type: %s", roleEntry.InferredEntityType)), nil
		case roleEntry.InferredAWSRegion == "":
			return logical.ErrorResponse("specified inferred_entity_type but not inferred_aws_region"), nil
		}
		allowEc2Binds = true
	} else if roleEntry.InferredAWSRegion != "" {
		return logical.ErrorResponse("specified inferred_aws_region but not inferred_entity_type"), nil
	}

	numBinds := 0

	if roleEntry.BoundAccountID != "" {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_account_id but not allowing ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if roleEntry.BoundRegion != "" {
		if roleEntry.AuthType != ec2AuthType {
			return logical.ErrorResponse("specified bound_region but not allowing ec2 auth_type"), nil
		}
		numBinds++
	}

	if roleEntry.BoundAmiID != "" {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_ami_id but not allowing ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if roleEntry.BoundIamInstanceProfileARN != "" {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_iam_instance_profile_arn but not allowing ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if roleEntry.BoundIamRoleARN != "" {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_iam_role_arn but not allowing ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if roleEntry.BoundIamPrincipalARN != "" {
		if roleEntry.AuthType != iamAuthType {
			return logical.ErrorResponse("specified bound_iam_principal_arn but not allowing iam auth_type"), nil
		}
		numBinds++
	}

	if roleEntry.BoundVpcID != "" {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_vpc_id but not allowing ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if roleEntry.BoundSubnetID != "" {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_subnet_id but not allowing ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if numBinds == 0 {
		return logical.ErrorResponse("at least be one bound parameter should be specified on the role"), nil
	}

	policiesRaw, ok := data.GetOk("policies")
	if ok {
		roleEntry.Policies = policyutil.ParsePolicies(policiesRaw)
	} else if req.Operation == logical.CreateOperation {
		roleEntry.Policies = []string{}
	}

	disallowReauthenticationBool, ok := data.GetOk("disallow_reauthentication")
	if ok {
		if roleEntry.AuthType != ec2AuthType {
			return logical.ErrorResponse("specified disallow_reauthentication when not using ec2 auth type"), nil
		}
		roleEntry.DisallowReauthentication = disallowReauthenticationBool.(bool)
	} else if req.Operation == logical.CreateOperation && roleEntry.AuthType == ec2AuthType {
		roleEntry.DisallowReauthentication = data.Get("disallow_reauthentication").(bool)
	}

	allowInstanceMigrationBool, ok := data.GetOk("allow_instance_migration")
	if ok {
		if roleEntry.AuthType != ec2AuthType {
			return logical.ErrorResponse("specified allow_instance_migration when not using ec2 auth type"), nil
		}
		roleEntry.AllowInstanceMigration = allowInstanceMigrationBool.(bool)
	} else if req.Operation == logical.CreateOperation && roleEntry.AuthType == ec2AuthType {
		roleEntry.AllowInstanceMigration = data.Get("allow_instance_migration").(bool)
	}

	if roleEntry.AllowInstanceMigration && roleEntry.DisallowReauthentication {
		return logical.ErrorResponse("cannot specify both disallow_reauthentication=true and allow_instance_migration=true"), nil
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

	periodRaw, ok := data.GetOk("period")
	if ok {
		roleEntry.Period = time.Second * time.Duration(periodRaw.(int))
	} else if req.Operation == logical.CreateOperation {
		roleEntry.Period = time.Second * time.Duration(data.Get("period").(int))
	}

	if roleEntry.Period > b.System().MaxLeaseTTL() {
		return logical.ErrorResponse(fmt.Sprintf("'period' of '%s' is greater than the backend's maximum lease TTL of '%s'", roleEntry.Period.String(), b.System().MaxLeaseTTL().String())), nil
	}

	roleTagStr, ok := data.GetOk("role_tag")
	if ok {
		if roleEntry.AuthType != ec2AuthType {
			return logical.ErrorResponse("tried to enable role_tag when not using ec2 auth method"), nil
		}
		roleEntry.RoleTag = roleTagStr.(string)
		// There is a limit of 127 characters on the tag key for AWS EC2 instances.
		// Complying to that requirement, do not allow the value of 'key' to be more than that.
		if len(roleEntry.RoleTag) > 127 {
			return logical.ErrorResponse("length of role tag exceeds the EC2 key limit of 127 characters"), nil
		}
	} else if req.Operation == logical.CreateOperation && roleEntry.AuthType == ec2AuthType {
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

	if len(resp.Warnings) == 0 {
		return nil, nil
	}

	return &resp, nil
}

// Struct to hold the information associated with an AMI ID in Vault.
type awsRoleEntry struct {
	AuthType                   string        `json:"auth_type" structs:"auth_type" mapstructure:"auth_type"`
	BoundAmiID                 string        `json:"bound_ami_id" structs:"bound_ami_id" mapstructure:"bound_ami_id"`
	BoundAccountID             string        `json:"bound_account_id" structs:"bound_account_id" mapstructure:"bound_account_id"`
	BoundIamPrincipalARN       string        `json:"bound_iam_principal_arn" structs:"bound_iam_principal_arn" mapstructure:"bound_iam_principal_arn"`
	BoundIamPrincipalID        string        `json:"bound_iam_principal_id" structs:"bound_iam_principal_id" mapstructure:"bound_iam_principal_id"`
	BoundIamRoleARN            string        `json:"bound_iam_role_arn" structs:"bound_iam_role_arn" mapstructure:"bound_iam_role_arn"`
	BoundIamInstanceProfileARN string        `json:"bound_iam_instance_profile_arn" structs:"bound_iam_instance_profile_arn" mapstructure:"bound_iam_instance_profile_arn"`
	BoundRegion                string        `json:"bound_region" structs:"bound_region" mapstructure:"bound_region"`
	BoundSubnetID              string        `json:"bound_subnet_id" structs:"bound_subnet_id" mapstructure:"bound_subnet_id"`
	BoundVpcID                 string        `json:"bound_vpc_id" structs:"bound_vpc_id" mapstructure:"bound_vpc_id"`
	InferredEntityType         string        `json:"inferred_entity_type" structs:"inferred_entity_type" mapstructure:"inferred_entity_type"`
	InferredAWSRegion          string        `json:"inferred_aws_region" structs:"inferred_aws_region" mapstructure:"inferred_aws_region"`
	ResolveAWSUniqueIDs        bool          `json:"resolve_aws_unique_ids" structs:"resolve_aws_unique_ids" mapstructure:"resolve_aws_unique_ids"`
	RoleTag                    string        `json:"role_tag" structs:"role_tag" mapstructure:"role_tag"`
	AllowInstanceMigration     bool          `json:"allow_instance_migration" structs:"allow_instance_migration" mapstructure:"allow_instance_migration"`
	TTL                        time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`
	MaxTTL                     time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
	Policies                   []string      `json:"policies" structs:"policies" mapstructure:"policies"`
	DisallowReauthentication   bool          `json:"disallow_reauthentication" structs:"disallow_reauthentication" mapstructure:"disallow_reauthentication"`
	HMACKey                    string        `json:"hmac_key" structs:"hmac_key" mapstructure:"hmac_key"`
	Period                     time.Duration `json:"period" mapstructure:"period" structs:"period"`
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
This only applies when authenticating EC2 instances.

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
