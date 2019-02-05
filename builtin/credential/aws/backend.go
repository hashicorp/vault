package awsauth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/vault/helper/awsutil"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	cache "github.com/patrickmn/go-cache"
)

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
	roleMutex sync.RWMutex

	// Lock to make changes to the blacklist entries
	blacklistMutex sync.RWMutex

	// Guards the blacklist/whitelist tidy functions
	tidyBlacklistCASGuard *uint32
	tidyWhitelistCASGuard *uint32

	// Duration after which the periodic function of the backend needs to
	// tidy the blacklist and whitelist entries.
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

	resolveArnToUniqueIDFunc func(context.Context, logical.Storage, string) (string, error)
}

func Backend(conf *logical.BackendConfig) (*backend, error) {
	b := &backend{
		// Setting the periodic func to be run once in an hour.
		// If there is a real need, this can be made configurable.
		tidyCooldownPeriod:    time.Hour,
		EC2ClientsMap:         make(map[string]map[string]*ec2.EC2),
		IAMClientsMap:         make(map[string]map[string]*iam.IAM),
		iamUserIdToArnCache:   cache.New(7*24*time.Hour, 24*time.Hour),
		tidyBlacklistCASGuard: new(uint32),
		tidyWhitelistCASGuard: new(uint32),
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
				"whitelist/identity/",
			},
			SealWrapStorage: []string{
				"config/client",
			},
		},
		Paths: []*framework.Path{
			pathLogin(b),
			pathListRole(b),
			pathListRoles(b),
			pathRole(b),
			pathRoleTag(b),
			pathConfigClient(b),
			pathConfigCertificate(b),
			pathConfigIdentity(b),
			pathConfigSts(b),
			pathListSts(b),
			pathConfigTidyRoletagBlacklist(b),
			pathConfigTidyIdentityWhitelist(b),
			pathListCertificates(b),
			pathListRoletagBlacklist(b),
			pathRoletagBlacklist(b),
			pathTidyRoletagBlacklist(b),
			pathListIdentityWhitelist(b),
			pathIdentityWhitelist(b),
			pathTidyIdentityWhitelist(b),
		},
		Invalidate:  b.invalidate,
		BackendType: logical.TypeCredential,
	}

	return b, nil
}

// periodicFunc performs the tasks that the backend wishes to do periodically.
// Currently this will be triggered once in a minute by the RollbackManager.
//
// The tasks being done currently by this function are to cleanup the expired
// entries of both blacklist role tags and whitelist identities. Tidying is done
// not once in a minute, but once in an hour, controlled by 'tidyCooldownPeriod'.
// Tidying of blacklist and whitelist are by default enabled. This can be
// changed using `config/tidy/roletags` and `config/tidy/identities` endpoints.
func (b *backend) periodicFunc(ctx context.Context, req *logical.Request) error {
	// Run the tidy operations for the first time. Then run it when current
	// time matches the nextTidyTime.
	if b.nextTidyTime.IsZero() || !time.Now().Before(b.nextTidyTime) {
		if b.System().LocalMount() || !b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationPerformanceStandby) {
			// safety_buffer defaults to 180 days for roletag blacklist
			safety_buffer := 15552000
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
				// overwrite the default safety_buffer with the configured value
				safety_buffer = tidyBlacklistConfigEntry.SafetyBuffer
			}
			// tidy role tags if explicitly not disabled
			if !skipBlacklistTidy {
				b.tidyBlacklistRoleTag(ctx, req, safety_buffer)
			}
		}

		// We don't check for replication state for whitelist identities as
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
			b.tidyWhitelistIdentity(ctx, req, safety_buffer)
		}

		// Update the time at which to run the tidy functions again.
		b.nextTidyTime = time.Now().Add(b.tidyCooldownPeriod)
	}
	return nil
}

func (b *backend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config/client":
		b.configMutex.Lock()
		defer b.configMutex.Unlock()
		b.flushCachedEC2Clients()
		b.flushCachedIAMClients()
		b.defaultAWSAccountID = ""
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
	region := getAnyRegionForAwsPartition(entity.Partition)
	if region == nil {
		return "", fmt.Errorf("unable to resolve partition %q to a region", entity.Partition)
	}
	iamClient, err := b.clientIAM(ctx, s, region.ID(), entity.AccountNumber)
	if err != nil {
		return "", awsutil.AppendLogicalError(err)
	}

	switch entity.Type {
	case "user":
		userInfo, err := iamClient.GetUser(&iam.GetUserInput{UserName: &entity.FriendlyName})
		if err != nil {
			return "", awsutil.AppendLogicalError(err)
		}
		if userInfo == nil {
			return "", fmt.Errorf("got nil result from GetUser")
		}
		return *userInfo.User.UserId, nil
	case "role":
		roleInfo, err := iamClient.GetRole(&iam.GetRoleInput{RoleName: &entity.FriendlyName})
		if err != nil {
			return "", awsutil.AppendLogicalError(err)
		}
		if roleInfo == nil {
			return "", fmt.Errorf("got nil result from GetRole")
		}
		return *roleInfo.Role.RoleId, nil
	case "instance-profile":
		profileInfo, err := iamClient.GetInstanceProfile(&iam.GetInstanceProfileInput{InstanceProfileName: &entity.FriendlyName})
		if err != nil {
			return "", awsutil.AppendLogicalError(err)
		}
		if profileInfo == nil {
			return "", fmt.Errorf("got nil result from GetInstanceProfile")
		}
		return *profileInfo.InstanceProfile.InstanceProfileId, nil
	default:
		return "", fmt.Errorf("unrecognized error type %#v", entity.Type)
	}
}

// Adapted from https://docs.aws.amazon.com/sdk-for-go/api/aws/endpoints/
// the "Enumerating Regions and Endpoint Metadata" section
func getAnyRegionForAwsPartition(partitionId string) *endpoints.Region {
	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()

	for _, p := range partitions {
		if p.ID() == partitionId {
			for _, r := range p.Regions() {
				return &r
			}
		}
	}
	return nil
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
