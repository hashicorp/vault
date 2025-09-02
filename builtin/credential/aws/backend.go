// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	cache "github.com/patrickmn/go-cache"
)

const (
	amzHeaderPrefix    = "X-Amz-"
	amzSignedHeaders   = "X-Amz-SignedHeaders"
	operationPrefixAWS = "aws"
)

var defaultAllowedSTSRequestHeaders = []string{
	"X-Amz-Algorithm",
	"X-Amz-Content-Sha256",
	"X-Amz-Credential",
	"X-Amz-Date",
	"X-Amz-Security-Token",
	"X-Amz-Signature",
	amzSignedHeaders,
	"X-Amz-User-Agent",
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend(conf)
	if err != nil {
		return nil, err
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

type backend struct {
	*framework.Backend

	// Lock to make changes to any of the backend's configuration endpoints.
	configMutex sync.RWMutex

	// Lock to make changes to role entries
	roleMutex sync.Mutex

	// Lock to make changes to the deny list entries
	denyListMutex sync.RWMutex

	// Guards the deny list/access list tidy functions
	tidyDenyListCASGuard   *uint32
	tidyAccessListCASGuard *uint32

	// Duration after which the periodic function of the backend needs to
	// tidy the deny list and access list entries.
	tidyCooldownPeriod time.Duration

	// nextTidyTime holds the time at which the periodic func should initiate
	// the tidy operations. This is set by the periodicFunc based on the value
	// of tidyCooldownPeriod.
	nextTidyTime time.Time

	// Map to hold the EC2 client objects indexed by region and STS role.
	// This avoids the overhead of creating a client object for every login request.
	// When the credentials are modified or deleted, all the cached client objects
	// will be flushed. The empty STS role signifies the master account
	EC2ClientsMap map[string]map[string]*ec2.EC2

	// Map to hold the IAM client objects indexed by region and STS role.
	// This avoids the overhead of creating a client object for every login request.
	// When the credentials are modified or deleted, all the cached client objects
	// will be flushed. The empty STS role signifies the master account
	IAMClientsMap map[string]map[string]*iam.IAM

	// Map to associate a partition to a random region in that partition. Users of
	// this don't care what region in the partition they use, but there is some client
	// cache efficiency gain if we keep the mapping stable, hence caching a single copy.
	partitionToRegionMap map[string]*endpoints.Region

	// Map of AWS unique IDs to the full ARN corresponding to that unique ID
	// This avoids the overhead of an AWS API hit for every login request
	// using the IAM auth method when bound_iam_principal_arn contains a wildcard
	iamUserIdToArnCache *cache.Cache

	// AWS Account ID of the "default" AWS credentials
	// This cache avoids the need to call GetCallerIdentity repeatedly to learn it
	// We can't store this because, in certain pathological cases, it could change
	// out from under us, such as a standby and active Vault server in different AWS
	// accounts using their IAM instance profile to get their credentials.
	defaultAWSAccountID string

	// roleCache caches role entries to avoid locking headaches
	roleCache *cache.Cache

	resolveArnToUniqueIDFunc func(context.Context, logical.Storage, string) (string, error)

	// upgradeCancelFunc is used to cancel the context used in the upgrade
	// function
	upgradeCancelFunc context.CancelFunc

	// deprecatedTerms is used to downgrade preferred terminology (e.g. accesslist)
	// to the legacy term. This allows for consolidated aliasing of the affected
	// endpoints until the legacy terms are removed.
	deprecatedTerms *strings.Replacer
}

func Backend(_ *logical.BackendConfig) (*backend, error) {
	b := &backend{
		// Setting the periodic func to be run once in an hour.
		// If there is a real need, this can be made configurable.
		tidyCooldownPeriod:     time.Hour,
		EC2ClientsMap:          make(map[string]map[string]*ec2.EC2),
		IAMClientsMap:          make(map[string]map[string]*iam.IAM),
		iamUserIdToArnCache:    cache.New(7*24*time.Hour, 24*time.Hour),
		tidyDenyListCASGuard:   new(uint32),
		tidyAccessListCASGuard: new(uint32),
		roleCache:              cache.New(cache.NoExpiration, cache.NoExpiration),

		deprecatedTerms: strings.NewReplacer(
			"accesslist", "whitelist",
			"access-list", "whitelist",
			"denylist", "blacklist",
			"deny-list", "blacklist",
		),
	}

	b.resolveArnToUniqueIDFunc = b.resolveArnToRealUniqueId

	b.Backend = &framework.Backend{
		PeriodicFunc: b.periodicFunc,
		AuthRenew:    b.pathLoginRenew,
		Help:         backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
			LocalStorage: []string{
				identityAccessListStorage,
			},
			SealWrapStorage: []string{
				"config/client",
			},
		},
		Paths: []*framework.Path{
			b.pathLogin(),
			b.pathListRole(),
			b.pathListRoles(),
			b.pathRole(),
			b.pathRoleTag(),
			b.pathConfigClient(),
			b.pathConfigCertificate(),
			b.pathConfigIdentity(),
			b.pathConfigRotateRoot(),
			b.pathConfigSts(),
			b.pathListSts(),
			b.pathListCertificates(),

			// The following pairs of functions are path aliases. The first is the
			// primary endpoint, and the second is version using deprecated language,
			// for backwards compatibility. The functionality is identical between the two.
			b.pathConfigTidyRoletagDenyList(),
			b.genDeprecatedPath(b.pathConfigTidyRoletagDenyList()),

			b.pathConfigTidyIdentityAccessList(),
			b.genDeprecatedPath(b.pathConfigTidyIdentityAccessList()),

			b.pathListRoletagDenyList(),
			b.genDeprecatedPath(b.pathListRoletagDenyList()),

			b.pathRoletagDenyList(),
			b.genDeprecatedPath(b.pathRoletagDenyList()),

			b.pathTidyRoletagDenyList(),
			b.genDeprecatedPath(b.pathTidyRoletagDenyList()),

			b.pathListIdentityAccessList(),
			b.genDeprecatedPath(b.pathListIdentityAccessList()),

			b.pathIdentityAccessList(),
			b.genDeprecatedPath(b.pathIdentityAccessList()),

			b.pathTidyIdentityAccessList(),
			b.genDeprecatedPath(b.pathTidyIdentityAccessList()),
		},
		Invalidate:     b.invalidate,
		InitializeFunc: b.initialize,
		BackendType:    logical.TypeCredential,
		Clean:          b.cleanup,
		RotateCredential: func(ctx context.Context, request *logical.Request) error {
			_, err := b.rotateRoot(ctx, request)
			return err
		},
	}

	b.partitionToRegionMap = generatePartitionToRegionMap()

	return b, nil
}

// periodicFunc performs the tasks that the backend wishes to do periodically.
// Currently this will be triggered once in a minute by the RollbackManager.
//
// The tasks being done currently by this function are to cleanup the expired
// entries of both deny list role tags and access list identities. Tidying is done
// not once in a minute, but once in an hour, controlled by 'tidyCooldownPeriod'.
// Tidying of deny list and access list are by default enabled. This can be
// changed using `config/tidy/roletags` and `config/tidy/identities` endpoints.
func (b *backend) periodicFunc(ctx context.Context, req *logical.Request) error {
	// Run the tidy operations for the first time. Then run it when current
	// time matches the nextTidyTime.
	if b.nextTidyTime.IsZero() || !time.Now().Before(b.nextTidyTime) {
		if b.System().LocalMount() || !b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationPerformanceStandby) {
			// safetyBuffer defaults to 180 days for roletag deny list
			safetyBuffer := 15552000
			tidyBlacklistConfigEntry, err := b.lockedConfigTidyRoleTags(ctx, req.Storage)
			if err != nil {
				return err
			}
			skipBlacklistTidy := false
			// check if tidying of role tags was configured
			if tidyBlacklistConfigEntry != nil {
				// check if periodic tidying of role tags was disabled
				if tidyBlacklistConfigEntry.DisablePeriodicTidy {
					skipBlacklistTidy = true
				}
				// overwrite the default safetyBuffer with the configured value
				safetyBuffer = tidyBlacklistConfigEntry.SafetyBuffer
			}
			// tidy role tags if explicitly not disabled
			if !skipBlacklistTidy {
				b.tidyDenyListRoleTag(ctx, req, safetyBuffer)
			}
		}

		// We don't check for replication state for access list identities as
		// these are locally stored

		safety_buffer := 259200
		tidyWhitelistConfigEntry, err := b.lockedConfigTidyIdentities(ctx, req.Storage)
		if err != nil {
			return err
		}
		skipWhitelistTidy := false
		// check if tidying of identities was configured
		if tidyWhitelistConfigEntry != nil {
			// check if periodic tidying of identities was disabled
			if tidyWhitelistConfigEntry.DisablePeriodicTidy {
				skipWhitelistTidy = true
			}
			// overwrite the default safety_buffer with the configured value
			safety_buffer = tidyWhitelistConfigEntry.SafetyBuffer
		}
		// tidy identities if explicitly not disabled
		if !skipWhitelistTidy {
			b.tidyAccessListIdentity(ctx, req, safety_buffer)
		}

		// Update the time at which to run the tidy functions again.
		b.nextTidyTime = time.Now().Add(b.tidyCooldownPeriod)
	}
	return nil
}

func (b *backend) cleanup(ctx context.Context) {
	if b.upgradeCancelFunc != nil {
		b.upgradeCancelFunc()
	}
}

func (b *backend) invalidate(ctx context.Context, key string) {
	switch {
	case key == "config/client":
		b.configMutex.Lock()
		defer b.configMutex.Unlock()
		b.flushCachedEC2Clients()
		b.flushCachedIAMClients()
		b.defaultAWSAccountID = ""
	case strings.HasPrefix(key, "role"):
		// TODO: We could make this better
		b.roleCache.Flush()
	}
}

// Putting this here so we can inject a fake resolver into the backend for unit testing
// purposes
func (b *backend) resolveArnToRealUniqueId(ctx context.Context, s logical.Storage, arn string) (string, error) {
	entity, err := parseIamArn(arn)
	if err != nil {
		return "", err
	}
	// This odd-looking code is here because IAM is an inherently global service. IAM and STS ARNs
	// don't have regions in them, and there is only a single global endpoint for IAM; see
	// http://docs.aws.amazon.com/general/latest/gr/rande.html#iam_region
	// However, the ARNs do have a partition in them, because the GovCloud and China partitions DO
	// have their own separate endpoints, and the partition is encoded in the ARN. If Amazon's Go SDK
	// would allow us to pass a partition back to the IAM client, it would be much simpler. But it
	// doesn't appear that's possible, so in order to properly support GovCloud and China, we do a
	// circular dance of extracting the partition from the ARN, finding any arbitrary region in the
	// partition, and passing that region back back to the SDK, so that the SDK can figure out the
	// proper partition from the arbitrary region we passed in to look up the endpoint.
	// Sigh
	region := b.partitionToRegionMap[entity.Partition]
	if region == nil {
		return "", fmt.Errorf("unable to resolve partition %q to a region", entity.Partition)
	}
	iamClient, err := b.clientIAM(ctx, s, region.ID(), entity.AccountNumber)
	if err != nil {
		return "", awsutil.AppendAWSError(err)
	}

	switch entity.Type {
	case "user":
		userInfo, err := iamClient.GetUserWithContext(ctx, &iam.GetUserInput{UserName: &entity.FriendlyName})
		if err != nil {
			return "", awsutil.AppendAWSError(err)
		}
		if userInfo == nil {
			return "", fmt.Errorf("got nil result from GetUser")
		}
		return *userInfo.User.UserId, nil
	case "role":
		roleInfo, err := iamClient.GetRoleWithContext(ctx, &iam.GetRoleInput{RoleName: &entity.FriendlyName})
		if err != nil {
			return "", awsutil.AppendAWSError(err)
		}
		if roleInfo == nil {
			return "", fmt.Errorf("got nil result from GetRole")
		}
		return *roleInfo.Role.RoleId, nil
	case "instance-profile":
		profileInfo, err := iamClient.GetInstanceProfileWithContext(ctx, &iam.GetInstanceProfileInput{InstanceProfileName: &entity.FriendlyName})
		if err != nil {
			return "", awsutil.AppendAWSError(err)
		}
		if profileInfo == nil {
			return "", fmt.Errorf("got nil result from GetInstanceProfile")
		}
		return *profileInfo.InstanceProfile.InstanceProfileId, nil
	default:
		return "", fmt.Errorf("unrecognized error type %#v", entity.Type)
	}
}

// genDeprecatedPath will return a deprecated version of a framework.Path. The
// path pattern and display attributes (if any) will contain deprecated terms,
// and the path will be marked as deprecated.
func (b *backend) genDeprecatedPath(path *framework.Path) *framework.Path {
	pathDeprecated := *path
	pathDeprecated.Pattern = b.deprecatedTerms.Replace(path.Pattern)
	pathDeprecated.Deprecated = true

	if path.DisplayAttrs != nil {
		deprecatedDisplayAttrs := *path.DisplayAttrs
		deprecatedDisplayAttrs.OperationPrefix = b.deprecatedTerms.Replace(path.DisplayAttrs.OperationPrefix)
		deprecatedDisplayAttrs.OperationVerb = b.deprecatedTerms.Replace(path.DisplayAttrs.OperationVerb)
		deprecatedDisplayAttrs.OperationSuffix = b.deprecatedTerms.Replace(path.DisplayAttrs.OperationSuffix)
		pathDeprecated.DisplayAttrs = &deprecatedDisplayAttrs
	}

	for i, op := range path.Operations {
		if op.Properties().DisplayAttrs != nil {
			deprecatedDisplayAttrs := *op.Properties().DisplayAttrs
			deprecatedDisplayAttrs.OperationPrefix = b.deprecatedTerms.Replace(op.Properties().DisplayAttrs.OperationPrefix)
			deprecatedDisplayAttrs.OperationVerb = b.deprecatedTerms.Replace(op.Properties().DisplayAttrs.OperationVerb)
			deprecatedDisplayAttrs.OperationSuffix = b.deprecatedTerms.Replace(op.Properties().DisplayAttrs.OperationSuffix)
			deprecatedProperties := pathDeprecated.Operations[i].(*framework.PathOperation)
			deprecatedProperties.DisplayAttrs = &deprecatedDisplayAttrs
		}
	}

	return &pathDeprecated
}

// Adapted from https://docs.aws.amazon.com/sdk-for-go/api/aws/endpoints/
// the "Enumerating Regions and Endpoint Metadata" section
func generatePartitionToRegionMap() map[string]*endpoints.Region {
	partitionToRegion := make(map[string]*endpoints.Region)

	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()

	for _, p := range partitions {
		// For most partitions, it's fine to choose a single region randomly.
		// However, there are a few exceptions:
		//
		//   For "aws", choose "us-east-1" because it is always enabled (and
		//   enabled for STS) by default.
		//
		//   For "aws-us-gov", choose "us-gov-west-1" because it is the only
		//   valid region for IAM operations.
		//   ref: https://github.com/aws/aws-sdk-go/blob/v1.34.25/aws/endpoints/defaults.go#L8176-L8194
		for _, r := range p.Regions() {
			if p.ID() == "aws" && r.ID() != "us-east-1" {
				continue
			}
			if p.ID() == "aws-us-gov" && r.ID() != "us-gov-west-1" {
				continue
			}
			partitionToRegion[p.ID()] = &r
			break
		}
	}

	return partitionToRegion
}

const backendHelp = `
The aws auth method uses either AWS IAM credentials or AWS-signed EC2 metadata
to authenticate clients, which are IAM principals or EC2 instances.

Authentication is backed by a preconfigured role in the backend. The role
represents the authorization of resources by containing Vault's policies.
Role can be created using 'role/<role>' endpoint.

Authentication of IAM principals, either IAM users or roles, is done using a
specifically signed AWS API request using clients' AWS IAM credentials. IAM
principals can then be assigned to roles within Vault. This is known as the
"iam" auth method.

Authentication of EC2 instances is done using either a signed PKCS#7 document
or a detached RSA signature of an AWS EC2 instance's identity document along
with a client-created nonce. This is known as the "ec2" auth method.

If there is need to further restrict the capabilities of the role on the instance
that is using the role, 'role_tag' option can be enabled on the role, and a tag
can be generated using 'role/<role>/tag' endpoint. This tag represents the
subset of capabilities set on the role. When the 'role_tag' option is enabled on
the role, the login operation requires that a respective role tag is attached to
the EC2 instance which performs the login.
`
