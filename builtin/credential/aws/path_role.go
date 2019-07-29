package awsauth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/copystructure"
)

var (
	currentRoleStorageVersion = 3
)

func pathRole(b *backend) *framework.Path {
	p := &framework.Path{
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
				Type: framework.TypeCommaStringSlice,
				Description: `If set, defines a constraint on the EC2 instances that they should be
using one of the AMI IDs specified by this parameter. This is only applicable
when auth_type is ec2 or inferred_entity_type is ec2_instance.`,
			},
			"bound_account_id": {
				Type: framework.TypeCommaStringSlice,
				Description: `If set, defines a constraint on the EC2 instances that the account ID
in its identity document to match one of the IDs specified by this parameter.
This is only applicable when auth_type is ec2 or inferred_entity_type is
ec2_instance.`,
			},
			"bound_iam_principal_arn": {
				Type: framework.TypeCommaStringSlice,
				Description: `ARN of the IAM principals to bind to this role. Only applicable when
auth_type is iam.`,
			},
			"bound_region": {
				Type: framework.TypeCommaStringSlice,
				Description: `If set, defines a constraint on the EC2 instances that the region in
its identity document match one of the regions specified by this parameter. This is only
applicable when auth_type is ec2.`,
			},
			"bound_iam_role_arn": {
				Type: framework.TypeCommaStringSlice,
				Description: `If set, defines a constraint on the authenticating EC2 instance
that it must match one of the IAM role ARNs specified by this parameter.
The value is prefix-matched (as though it were a glob ending in
'*').  The configured IAM user or EC2 instance role must be allowed
to execute the 'iam:GetInstanceProfile' action if this is specified. This is
only applicable when auth_type is ec2 or inferred_entity_type is
ec2_instance.`,
			},
			"bound_iam_instance_profile_arn": {
				Type: framework.TypeCommaStringSlice,
				Description: `If set, defines a constraint on the EC2 instances to be associated
with an IAM instance profile ARN which has a prefix that matches
one of the values specified by this parameter. The value is prefix-matched
(as though it were a glob ending in '*'). This is only applicable when
auth_type is ec2 or inferred_entity_type is ec2_instance.`,
			},
			"bound_ec2_instance_id": {
				Type: framework.TypeCommaStringSlice,
				Description: `If set, defines a constraint on the EC2 instances to have one of the
given instance IDs. Can be a list or comma-separated string of EC2 instance
IDs. This is only applicable when auth_type is ec2 or inferred_entity_type is
ec2_instance.`,
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
				Type: framework.TypeCommaStringSlice,
				Description: `
If set, defines a constraint on the EC2 instance to be associated with a VPC
ID that matches one of the value specified by this parameter. This is only
applicable when auth_type is ec2 or inferred_entity_type is ec2_instance.`,
			},
			"bound_subnet_id": {
				Type: framework.TypeCommaStringSlice,
				Description: `
If set, defines a constraint on the EC2 instance to be associated with the
subnet ID that matches one of the values specified by this parameter. This is
only applicable when auth_type is ec2 or inferred_entity_type is
ec2_instance.`,
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
				Type:        framework.TypeDurationSecond,
				Description: tokenutil.DeprecationText("token_period"),
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
			"policies": {
				Type:        framework.TypeCommaStringSlice,
				Description: tokenutil.DeprecationText("token_policies"),
				Deprecated:  true,
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

	tokenutil.AddTokenFields(p.Fields)
	return p
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
func (b *backend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.role(ctx, req.Storage, strings.ToLower(data.Get("role").(string)))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// role fetches the role entry from cache, or loads from disk if necessary
func (b *backend) role(ctx context.Context, s logical.Storage, roleName string) (*awsRoleEntry, error) {
	if roleName == "" {
		return nil, fmt.Errorf("missing role name")
	}

	roleEntryRaw, found := b.roleCache.Get(roleName)
	if found && roleEntryRaw != nil {
		roleEntry, ok := roleEntryRaw.(*awsRoleEntry)
		if !ok {
			return nil, errors.New("could not convert role entry internally")
		}
		if roleEntry == nil {
			return nil, errors.New("converted role entry is nil")
		}
	}

	// Not found, or was nil
	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	return b.roleInternal(ctx, s, roleName)
}

// roleInternal does not perform locking, and rechecks the cache, going to disk if necessar
func (b *backend) roleInternal(ctx context.Context, s logical.Storage, roleName string) (*awsRoleEntry, error) {
	// Check cache again now that we have the lock
	roleEntryRaw, found := b.roleCache.Get(roleName)
	if found && roleEntryRaw != nil {
		roleEntry, ok := roleEntryRaw.(*awsRoleEntry)
		if !ok {
			return nil, errors.New("could not convert role entry internally")
		}
		if roleEntry == nil {
			return nil, errors.New("converted role entry is nil")
		}
	}

	// Fetch from storage
	entry, err := s.Get(ctx, "role/"+strings.ToLower(roleName))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	result := new(awsRoleEntry)
	if err := entry.DecodeJSON(result); err != nil {
		return nil, err
	}

	needUpgrade, err := b.upgradeRole(ctx, s, result)
	if err != nil {
		return nil, errwrap.Wrapf("error upgrading roleEntry: {{err}}", err)
	}
	if needUpgrade && (b.System().LocalMount() || !b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationPerformanceStandby)) {
		if err = b.setRole(ctx, s, roleName, result); err != nil {
			return nil, errwrap.Wrapf("error saving upgraded roleEntry: {{err}}", err)
		}
	}

	b.roleCache.SetDefault(roleName, result)

	return result, nil
}

// setRole creates or updates a role in the storage. The caller must hold
// the write lock.
func (b *backend) setRole(ctx context.Context, s logical.Storage, roleName string,
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

	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	b.roleCache.SetDefault(roleName, roleEntry)

	return nil
}

// initialize is used to initialize the AWS roles
func (b *backend) initialize(ctx context.Context, req *logical.InitializationRequest) error {

	// on standbys and DR secondaries we do not want to run any kind of upgrade logic
	if b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby | consts.ReplicationDRSecondary) {
		return nil
	}

	// Initialize only if we are either:
	//   (1) A local mount.
	//   (2) Are _NOT_ a replicated performance secondary
	if b.System().LocalMount() || !b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary) {

		s := req.Storage

		logger := b.Logger().Named("initialize")
		logger.Debug("starting initialization")

		var upgradeCtx context.Context
		upgradeCtx, b.upgradeCancelFunc = context.WithCancel(context.Background())

		go func() {
			// The vault will become unsealed while this goroutine is running,
			// so we could see some role requests block until the lock is
			// released.  However we'd rather see those requests block (and
			// potentially start timing out) than allow a non-upgraded role to
			// be fetched.
			b.roleMutex.Lock()
			defer b.roleMutex.Unlock()

			upgraded, err := b.upgrade(upgradeCtx, s)
			if err != nil {
				logger.Error("error running initialization", "error", err)
				return
			}
			if upgraded {
				logger.Info("an upgrade was performed during initialization")
			}
		}()

	}

	return nil
}

// awsVersion stores info about the the latest aws version that we have
// upgraded to.
type awsVersion struct {
	Version int `json:"version"`
}

// currentAwsVersion stores the latest version that we have upgraded to.
// Note that this is tracked independently from currentRoleStorageVersion.
const currentAwsVersion = 1

// upgrade does an upgrade, if necessary
func (b *backend) upgrade(ctx context.Context, s logical.Storage) (bool, error) {
	entry, err := s.Get(ctx, "config/version")
	if err != nil {
		return false, err
	}
	var version awsVersion
	if entry != nil {
		err = entry.DecodeJSON(&version)
		if err != nil {
			return false, err
		}
	}

	upgraded := version.Version < currentAwsVersion
	switch version.Version {
	case 0:
		// Read all the role names.
		roleNames, err := s.List(ctx, "role/")
		if err != nil {
			return false, err
		}

		// Upgrade the roles as necessary.
		for _, roleName := range roleNames {
			// make sure the context hasn't been canceled
			if ctx.Err() != nil {
				return false, err
			}
			_, err := b.roleInternal(ctx, s, roleName)
			if err != nil {
				return false, err
			}
		}
		fallthrough

	case currentAwsVersion:
		version.Version = currentAwsVersion

	default:
		return false, fmt.Errorf("unrecognized role version: %d", version.Version)
	}

	// save the current version
	if upgraded {
		entry, err = logical.StorageEntryJSON("config/version", &version)
		if err != nil {
			return false, err
		}
		err = s.Put(ctx, entry)
		if err != nil {
			return false, err
		}
	}

	return upgraded, nil
}

// If needed, updates the role entry and returns a bool indicating if it was updated
// (and thus needs to be persisted)
func (b *backend) upgradeRole(ctx context.Context, s logical.Storage, roleEntry *awsRoleEntry) (bool, error) {
	if roleEntry == nil {
		return false, fmt.Errorf("received nil roleEntry")
	}
	upgraded := roleEntry.Version < currentRoleStorageVersion
	switch roleEntry.Version {
	case 0:
		// Check if the value held by role ARN field is actually an instance profile ARN
		if roleEntry.BoundIamRoleARN != "" && strings.Contains(roleEntry.BoundIamRoleARN, ":instance-profile/") {
			// If yes, move it to the correct field
			roleEntry.BoundIamInstanceProfileARN = roleEntry.BoundIamRoleARN

			// Reset the old field
			roleEntry.BoundIamRoleARN = ""
		}

		// Check if there was no pre-existing AuthType set (from older versions)
		if roleEntry.AuthType == "" {
			// then default to the original behavior of ec2
			roleEntry.AuthType = ec2AuthType
		}

		// Check if we need to resolve the unique ID on the role
		if roleEntry.AuthType == iamAuthType &&
			roleEntry.ResolveAWSUniqueIDs &&
			roleEntry.BoundIamPrincipalARN != "" &&
			roleEntry.BoundIamPrincipalID == "" &&
			!strings.HasSuffix(roleEntry.BoundIamPrincipalARN, "*") {
			principalId, err := b.resolveArnToUniqueIDFunc(ctx, s, roleEntry.BoundIamPrincipalARN)
			if err != nil {
				return false, err
			}
			roleEntry.BoundIamPrincipalID = principalId
			// Not setting roleEntry.BoundIamPrincipalARN to "" here so that clients can see the original
			// ARN that the role was bound to
		}

		// Check if we need to convert individual string values to lists
		if roleEntry.BoundAmiID != "" {
			roleEntry.BoundAmiIDs = []string{roleEntry.BoundAmiID}
			roleEntry.BoundAmiID = ""
		}
		if roleEntry.BoundAccountID != "" {
			roleEntry.BoundAccountIDs = []string{roleEntry.BoundAccountID}
			roleEntry.BoundAccountID = ""
		}
		if roleEntry.BoundIamPrincipalARN != "" {
			roleEntry.BoundIamPrincipalARNs = []string{roleEntry.BoundIamPrincipalARN}
			roleEntry.BoundIamPrincipalARN = ""
		}
		if roleEntry.BoundIamPrincipalID != "" {
			roleEntry.BoundIamPrincipalIDs = []string{roleEntry.BoundIamPrincipalID}
			roleEntry.BoundIamPrincipalID = ""
		}
		if roleEntry.BoundIamRoleARN != "" {
			roleEntry.BoundIamRoleARNs = []string{roleEntry.BoundIamRoleARN}
			roleEntry.BoundIamRoleARN = ""
		}
		if roleEntry.BoundIamInstanceProfileARN != "" {
			roleEntry.BoundIamInstanceProfileARNs = []string{roleEntry.BoundIamInstanceProfileARN}
			roleEntry.BoundIamInstanceProfileARN = ""
		}
		if roleEntry.BoundRegion != "" {
			roleEntry.BoundRegions = []string{roleEntry.BoundRegion}
			roleEntry.BoundRegion = ""
		}
		if roleEntry.BoundSubnetID != "" {
			roleEntry.BoundSubnetIDs = []string{roleEntry.BoundSubnetID}
			roleEntry.BoundSubnetID = ""
		}
		if roleEntry.BoundVpcID != "" {
			roleEntry.BoundVpcIDs = []string{roleEntry.BoundVpcID}
			roleEntry.BoundVpcID = ""
		}
		fallthrough

	case 1:
		// Make BoundIamRoleARNs and BoundIamInstanceProfileARNs explicitly prefix-matched
		for i, arn := range roleEntry.BoundIamRoleARNs {
			roleEntry.BoundIamRoleARNs[i] = fmt.Sprintf("%s*", arn)
		}
		for i, arn := range roleEntry.BoundIamInstanceProfileARNs {
			roleEntry.BoundIamInstanceProfileARNs[i] = fmt.Sprintf("%s*", arn)
		}
		fallthrough

	case 2:
		roleID, err := uuid.GenerateUUID()
		if err != nil {
			return false, err
		}
		roleEntry.RoleID = roleID
		fallthrough

	case currentRoleStorageVersion:
		roleEntry.Version = currentRoleStorageVersion

	default:
		return false, fmt.Errorf("unrecognized role version: %q", roleEntry.Version)
	}

	// Add tokenutil upgrades. These don't need to be persisted, they're fine
	// being upgraded each time until changed.
	if roleEntry.TokenTTL == 0 && roleEntry.TTL > 0 {
		roleEntry.TokenTTL = roleEntry.TTL
	}
	if roleEntry.TokenMaxTTL == 0 && roleEntry.MaxTTL > 0 {
		roleEntry.TokenMaxTTL = roleEntry.MaxTTL
	}
	if roleEntry.TokenPeriod == 0 && roleEntry.Period > 0 {
		roleEntry.TokenPeriod = roleEntry.Period
	}
	if len(roleEntry.TokenPolicies) == 0 && len(roleEntry.Policies) > 0 {
		roleEntry.TokenPolicies = roleEntry.Policies
	}

	return upgraded, nil
}

// pathRoleDelete is used to delete the information registered for a given AMI ID.
func (b *backend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	err := req.Storage.Delete(ctx, "role/"+strings.ToLower(roleName))
	if err != nil {
		return nil, errwrap.Wrapf("error deleting role: {{err}}", err)
	}

	b.roleCache.Delete(roleName)

	return nil, nil
}

// pathRoleList is used to list all the AMI IDs registered with Vault.
func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(roles), nil
}

// pathRoleRead is used to view the information registered for a given AMI ID.
func (b *backend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleEntry, err := b.role(ctx, req.Storage, strings.ToLower(data.Get("role").(string)))
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: roleEntry.ToResponseData(),
	}, nil
}

// pathRoleCreateUpdate is used to associate Vault policies to a given AMI ID.
func (b *backend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := strings.ToLower(data.Get("role").(string))
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	// We use the internal one here to ensure that we have fresh data and
	// nobody else is concurrently modifying. This will also call the upgrade
	// path on existing role entries.
	roleEntry, err := b.roleInternal(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		roleID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}
		roleEntry = &awsRoleEntry{
			RoleID:  roleID,
			Version: currentRoleStorageVersion,
		}
	} else {
		// We want to always use a copy so we aren't modifying items in the
		// version in the cache while other users may be looking it up (or if
		// we fail somewhere)
		cp, err := copystructure.Copy(roleEntry)
		if err != nil {
			return nil, err
		}
		roleEntry = cp.(*awsRoleEntry)
	}

	// Fetch and set the bound parameters. There can't be default values
	// for these.
	if boundAmiIDRaw, ok := data.GetOk("bound_ami_id"); ok {
		roleEntry.BoundAmiIDs = boundAmiIDRaw.([]string)
	}

	if boundAccountIDRaw, ok := data.GetOk("bound_account_id"); ok {
		roleEntry.BoundAccountIDs = boundAccountIDRaw.([]string)
	}

	if boundRegionRaw, ok := data.GetOk("bound_region"); ok {
		roleEntry.BoundRegions = boundRegionRaw.([]string)
	}

	if boundVpcIDRaw, ok := data.GetOk("bound_vpc_id"); ok {
		roleEntry.BoundVpcIDs = boundVpcIDRaw.([]string)
	}

	if boundSubnetIDRaw, ok := data.GetOk("bound_subnet_id"); ok {
		roleEntry.BoundSubnetIDs = boundSubnetIDRaw.([]string)
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
		roleEntry.BoundIamRoleARNs = boundIamRoleARNRaw.([]string)
	}

	if boundIamInstanceProfileARNRaw, ok := data.GetOk("bound_iam_instance_profile_arn"); ok {
		roleEntry.BoundIamInstanceProfileARNs = boundIamInstanceProfileARNRaw.([]string)
	}

	if boundEc2InstanceIDRaw, ok := data.GetOk("bound_ec2_instance_id"); ok {
		roleEntry.BoundEc2InstanceIDs = boundEc2InstanceIDRaw.([]string)
	}

	if boundIamPrincipalARNRaw, ok := data.GetOk("bound_iam_principal_arn"); ok {
		principalARNs := boundIamPrincipalARNRaw.([]string)
		roleEntry.BoundIamPrincipalARNs = principalARNs
		roleEntry.BoundIamPrincipalIDs = []string{}
	}
	if roleEntry.ResolveAWSUniqueIDs && len(roleEntry.BoundIamPrincipalIDs) == 0 {
		// we might be turning on resolution on this role, so ensure we update the IDs
		for _, principalARN := range roleEntry.BoundIamPrincipalARNs {
			if !strings.HasSuffix(principalARN, "*") {
				principalID, err := b.resolveArnToUniqueIDFunc(ctx, req.Storage, principalARN)
				if err != nil {
					return logical.ErrorResponse(fmt.Sprintf("unable to resolve ARN %#v to internal ID: %s", principalARN, err.Error())), nil
				}
				roleEntry.BoundIamPrincipalIDs = append(roleEntry.BoundIamPrincipalIDs, principalID)
			}
		}
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

	if len(roleEntry.BoundAccountIDs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_account_id but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundRegions) > 0 {
		if roleEntry.AuthType != ec2AuthType {
			return logical.ErrorResponse("specified bound_region but not specifying ec2 auth_type"), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundAmiIDs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_ami_id but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundIamInstanceProfileARNs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_iam_instance_profile_arn but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundEc2InstanceIDs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_ec2_instance_id but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundIamRoleARNs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_iam_role_arn but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundIamPrincipalARNs) > 0 {
		if roleEntry.AuthType != iamAuthType {
			return logical.ErrorResponse("specified bound_iam_principal_arn but not specifying iam auth_type"), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundVpcIDs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_vpc_id but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundSubnetIDs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_subnet_id but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if numBinds == 0 {
		return logical.ErrorResponse("at least one bound parameter should be specified on the role"), nil
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

	if err := roleEntry.ParseTokenFields(req, data); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Handle upgrade cases
	{
		if err := tokenutil.UpgradeValue(data, "policies", "token_policies", &roleEntry.Policies, &roleEntry.TokenPolicies); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "ttl", "token_ttl", &roleEntry.TTL, &roleEntry.TokenTTL); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
		// Special case here for old lease value
		_, ok := data.GetOk("token_ttl")
		if !ok {
			_, ok = data.GetOk("ttl")
			if !ok {
				ttlRaw, ok := data.GetOk("lease")
				if ok {
					roleEntry.TTL = time.Duration(ttlRaw.(int)) * time.Second
					roleEntry.TokenTTL = roleEntry.TTL
				}
			}
		}

		if err := tokenutil.UpgradeValue(data, "max_ttl", "token_max_ttl", &roleEntry.MaxTTL, &roleEntry.TokenMaxTTL); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "period", "token_period", &roleEntry.Period, &roleEntry.TokenPeriod); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	defaultLeaseTTL := b.System().DefaultLeaseTTL()
	systemMaxTTL := b.System().MaxLeaseTTL()
	if roleEntry.TokenTTL > defaultLeaseTTL {
		resp.AddWarning(fmt.Sprintf("Given ttl of %d seconds greater than current mount/system default of %d seconds; ttl will be capped at login time", roleEntry.TokenTTL/time.Second, defaultLeaseTTL/time.Second))
	}
	if roleEntry.TokenMaxTTL > systemMaxTTL {
		resp.AddWarning(fmt.Sprintf("Given max ttl of %d seconds greater than current mount/system default of %d seconds; max ttl will be capped at login time", roleEntry.TokenMaxTTL/time.Second, systemMaxTTL/time.Second))
	}
	if roleEntry.TokenMaxTTL != 0 && roleEntry.TokenMaxTTL < roleEntry.TokenTTL {
		return logical.ErrorResponse("ttl should be shorter than max ttl"), nil
	}
	if roleEntry.TokenPeriod > b.System().MaxLeaseTTL() {
		return logical.ErrorResponse(fmt.Sprintf("period of '%s' is greater than the backend's maximum lease TTL of '%s'", roleEntry.TokenPeriod.String(), b.System().MaxLeaseTTL().String())), nil
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
			return nil, errwrap.Wrapf("failed to generate role HMAC key: {{err}}", err)
		}
	}

	if err := b.setRole(ctx, req.Storage, roleName, roleEntry); err != nil {
		return nil, err
	}

	if len(resp.Warnings) == 0 {
		return nil, nil
	}

	return &resp, nil
}

// Struct to hold the information associated with a Vault role
type awsRoleEntry struct {
	tokenutil.TokenParams

	RoleID                      string   `json:"role_id"`
	AuthType                    string   `json:"auth_type"`
	BoundAmiIDs                 []string `json:"bound_ami_id_list"`
	BoundAccountIDs             []string `json:"bound_account_id_list"`
	BoundEc2InstanceIDs         []string `json:"bound_ec2_instance_id_list"`
	BoundIamPrincipalARNs       []string `json:"bound_iam_principal_arn_list"`
	BoundIamPrincipalIDs        []string `json:"bound_iam_principal_id_list"`
	BoundIamRoleARNs            []string `json:"bound_iam_role_arn_list"`
	BoundIamInstanceProfileARNs []string `json:"bound_iam_instance_profile_arn_list"`
	BoundRegions                []string `json:"bound_region_list"`
	BoundSubnetIDs              []string `json:"bound_subnet_id_list"`
	BoundVpcIDs                 []string `json:"bound_vpc_id_list"`
	InferredEntityType          string   `json:"inferred_entity_type"`
	InferredAWSRegion           string   `json:"inferred_aws_region"`
	ResolveAWSUniqueIDs         bool     `json:"resolve_aws_unique_ids"`
	RoleTag                     string   `json:"role_tag"`
	AllowInstanceMigration      bool     `json:"allow_instance_migration"`
	DisallowReauthentication    bool     `json:"disallow_reauthentication"`
	HMACKey                     string   `json:"hmac_key"`
	Version                     int      `json:"version"`

	// Deprecated: These are superceded by TokenUtil
	TTL      time.Duration `json:"ttl"`
	MaxTTL   time.Duration `json:"max_ttl"`
	Period   time.Duration `json:"period"`
	Policies []string      `json:"policies"`

	// DEPRECATED -- these are the old fields before we supported lists and exist for backwards compatibility
	BoundAmiID                 string `json:"bound_ami_id,omitempty" `
	BoundAccountID             string `json:"bound_account_id,omitempty"`
	BoundIamPrincipalARN       string `json:"bound_iam_principal_arn,omitempty"`
	BoundIamPrincipalID        string `json:"bound_iam_principal_id,omitempty"`
	BoundIamRoleARN            string `json:"bound_iam_role_arn,omitempty"`
	BoundIamInstanceProfileARN string `json:"bound_iam_instance_profile_arn,omitempty"`
	BoundRegion                string `json:"bound_region,omitempty"`
	BoundSubnetID              string `json:"bound_subnet_id,omitempty"`
	BoundVpcID                 string `json:"bound_vpc_id,omitempty"`
}

func (r *awsRoleEntry) ToResponseData() map[string]interface{} {
	responseData := map[string]interface{}{
		"auth_type":                      r.AuthType,
		"bound_ami_id":                   r.BoundAmiIDs,
		"bound_account_id":               r.BoundAccountIDs,
		"bound_ec2_instance_id":          r.BoundEc2InstanceIDs,
		"bound_iam_principal_arn":        r.BoundIamPrincipalARNs,
		"bound_iam_principal_id":         r.BoundIamPrincipalIDs,
		"bound_iam_role_arn":             r.BoundIamRoleARNs,
		"bound_iam_instance_profile_arn": r.BoundIamInstanceProfileARNs,
		"bound_region":                   r.BoundRegions,
		"bound_subnet_id":                r.BoundSubnetIDs,
		"bound_vpc_id":                   r.BoundVpcIDs,
		"inferred_entity_type":           r.InferredEntityType,
		"inferred_aws_region":            r.InferredAWSRegion,
		"resolve_aws_unique_ids":         r.ResolveAWSUniqueIDs,
		"role_id":                        r.RoleID,
		"role_tag":                       r.RoleTag,
		"allow_instance_migration":       r.AllowInstanceMigration,
		"disallow_reauthentication":      r.DisallowReauthentication,
	}

	r.PopulateTokenData(responseData)
	if r.TTL > 0 {
		responseData["ttl"] = int64(r.TTL.Seconds())
	}
	if r.MaxTTL > 0 {
		responseData["max_ttl"] = int64(r.MaxTTL.Seconds())
	}
	if r.Period > 0 {
		responseData["period"] = int64(r.Period.Seconds())
	}
	if len(r.Policies) > 0 {
		responseData["policies"] = responseData["token_policies"]
	}

	convertNilToEmptySlice := func(data map[string]interface{}, field string) {
		if data[field] == nil || len(data[field].([]string)) == 0 {
			data[field] = []string{}
		}
	}
	convertNilToEmptySlice(responseData, "bound_ami_id")
	convertNilToEmptySlice(responseData, "bound_account_id")
	convertNilToEmptySlice(responseData, "bound_iam_principal_arn")
	convertNilToEmptySlice(responseData, "bound_iam_principal_id")
	convertNilToEmptySlice(responseData, "bound_iam_role_arn")
	convertNilToEmptySlice(responseData, "bound_iam_instance_profile_arn")
	convertNilToEmptySlice(responseData, "bound_region")
	convertNilToEmptySlice(responseData, "bound_subnet_id")
	convertNilToEmptySlice(responseData, "bound_vpc_id")

	return responseData
}

const pathRoleSyn = `
Create a role and associate policies to it.
`

const pathRoleDesc = `
A precondition for login is that a role should be created in the backend.
The login endpoint takes in the role name against which the client
should be validated. After authenticating the client, the authorization
to access Vault's resources is determined by the policies that are
associated to the role though this endpoint.

When an EC2 instance requires only a subset of policies on the role, then
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
