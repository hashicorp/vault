package aws

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend(conf)
	if err != nil {
		return nil, err
	}
	return b.Setup(conf)
}

type backend struct {
	*framework.Backend
	Salt *salt.Salt

	// Lock to make changes to any of the backend's configuration endpoints.
	configMutex sync.RWMutex

	// Lock to make changes to the blacklist entries
	blacklistMutex sync.RWMutex

	// Guards the blacklist/whitelist tidy functions
	tidyBlacklistCASGuard uint32
	tidyWhitelistCASGuard uint32

	// Duration after which the periodic function of the backend needs to
	// tidy the blacklist and whitelist entries.
	tidyCooldownPeriod time.Duration

	// nextTidyTime holds the time at which the periodic func should initiatite
	// the tidy operations. This is set by the periodicFunc based on the value
	// of tidyCooldownPeriod.
	nextTidyTime time.Time

	// Map to hold the EC2 client objects indexed by region. This avoids the
	// overhead of creating a client object for every login request. When
	// the credentials are modified or deleted, all the cached client objects
	// will be flushed.
	EC2ClientsMap map[string]*ec2.EC2
}

func Backend(conf *logical.BackendConfig) (*framework.Backend, error) {
	salt, err := salt.NewSalt(conf.StorageView, &salt.Config{
		HashFunc: salt.SHA256Hash,
	})
	if err != nil {
		return nil, err
	}

	b := &backend{
		// Setting the periodic func to be run once in an hour.
		// If there is a real need, this can be made configurable.
		tidyCooldownPeriod: time.Hour,
		Salt:               salt,
		EC2ClientsMap:      make(map[string]*ec2.EC2),
	}

	b.Backend = &framework.Backend{
		PeriodicFunc: b.periodicFunc,
		AuthRenew:    b.pathLoginRenew,
		Help:         backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
		},
		Paths: []*framework.Path{
			pathLogin(b),
			pathRole(b),
			pathListRoles(b),
			pathRoleTag(b),
			pathConfigClient(b),
			pathConfigCertificate(b),
			pathConfigTidyRoleTags(b),
			pathConfigTidyIdentities(b),
			pathListCertificates(b),
			pathListBlacklistRoleTags(b),
			pathBlacklistRoleTag(b),
			pathTidyRoleTags(b),
			pathListWhitelistIdentities(b),
			pathWhitelistIdentity(b),
			pathTidyIdentities(b),
		},
	}

	return b.Backend, nil
}

// periodicFunc performs the tasks that the backend wishes to do periodically.
// Currently this will be triggered once in a minute by the RollbackManager.
//
// The tasks being done currently by this function are to cleanup the expired
// entries of both blacklist role tags and whitelist identities. Tidying is done
// not once in a minute, but once in an hour, controlled by 'tidyCooldownPeriod'.
// Tidying of blacklist and whitelist are by default enabled. This can be
// changed using `config/tidy/roletags` and `config/tidy/identities` endpoints.
func (b *backend) periodicFunc(req *logical.Request) error {

	// Run the tidy operations for the first time. Then run it when current
	// time matches the nextTidyTime.
	if b.nextTidyTime.IsZero() || !time.Now().UTC().Before(b.nextTidyTime) {
		// safety_buffer defaults to 180 days for roletag blacklist
		safety_buffer := 15552000
		tidyBlacklistConfigEntry, err := b.configTidyRoleTags(req.Storage)
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
			b.tidyBlacklistRoleTag(req.Storage, safety_buffer)
		}

		// reset the safety_buffer to 72h
		safety_buffer = 259200
		tidyWhitelistConfigEntry, err := b.configTidyIdentities(req.Storage)
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
			b.tidyWhitelistIdentity(req.Storage, safety_buffer)
		}

		// Update the time at which to run the tidy functions again.
		b.nextTidyTime = time.Now().UTC().Add(b.tidyCooldownPeriod)
	}
	return nil
}

const backendHelp = `
AWS auth backend takes in PKCS#7 signature of an AWS EC2 instance and a client
created nonce to authenticates the EC2 instance with Vault.

Authentication is backed by a preconfigured role in the backend. The role
represents the authorization of resources by containing Vault's policies.
Role can be created using 'role/<role_name>' endpoint.

If there is need to further restrict the policies set on the role, 'role_tag' option
can be enabled on the role, and a tag can be generated using 'role/<role_name>/tag'
endpoint. This tag represents the subset of capabilities set on the role. When the
'role_tag' option is enabled on the role, the login operation requires that a respective
role tag is attached to the EC2 instance that is performing the login.
`
