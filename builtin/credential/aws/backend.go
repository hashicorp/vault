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
		Paths: append([]*framework.Path{
			pathLogin(b),
			pathImage(b),
			pathListImages(b),
			pathImageTag(b),
			pathConfigClient(b),
			pathConfigCertificate(b),
			pathConfigTidyRoleTags(b),
			pathConfigTidyIdentities(b),
			pathListCertificates(b),
			pathBlacklistRoleTag(b),
			pathListBlacklistRoleTags(b),
			pathTidyRoleTags(b),
			pathWhitelistIdentity(b),
			pathTidyIdentities(b),
			pathListWhitelistIdentities(b),
		}),
	}

	return b.Backend, nil
}

type backend struct {
	*framework.Backend
	Salt *salt.Salt

	// Lock to make changes to any of the backend's configuration endpoints.
	configMutex sync.RWMutex

	// Duration after which the periodic function of the backend needs to be
	// executed.
	tidyCooldownPeriod time.Duration

	// Var that holds the time at which the periodic func should initiatite
	// the tidy operations.
	nextTidyTime time.Time

	// Map to hold the EC2 client objects indexed by region. This avoids the
	// overhead of creating a client object for every login request.
	EC2ClientsMap map[string]*ec2.EC2
}

// periodicFunc performs the tasks that the backend wishes to do periodically.
// Currently this will be triggered once in a minute by the RollbackManager.
//
// The tasks being done are to cleanup the expired entries of both blacklist
// and whitelist. Tidying is done not once in a minute, but once in an hour.
// Tidying of blacklist and whitelist are by default enabled. This can be
// changed using `config/tidy/roletags` and `config/tidy/identities` endpoints.
func (b *backend) periodicFunc(req *logical.Request) error {
	if b.nextTidyTime.IsZero() || !time.Now().UTC().Before(b.nextTidyTime) {
		// safety_buffer defaults to 72h
		safety_buffer := 259200
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
			tidyBlacklistRoleTag(req.Storage, safety_buffer)
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
			tidyWhitelistIdentity(req.Storage, safety_buffer)
		}

		// Update the nextTidyTime
		b.nextTidyTime = time.Now().UTC().Add(b.tidyCooldownPeriod)
	}
	return nil
}

const backendHelp = `
AWS auth backend takes in PKCS#7 signature of an AWS EC2 instance and a client
created nonce to authenticates the EC2 instance with Vault.

Authentication is backed by a preconfigured association of AMIs to Vault's policies
through 'image/<ami_id>' endpoint. All the instances that are using this AMI will
get the policies configured on the AMI.

If there is need to further restrict the policies set on the AMI, 'role_tag' option
can be enabled on the AMI and a tag can be generated using 'image/<ami_id>/roletag'
endpoint. This tag represents the subset of capabilities set on the AMI. When the
'role_tag' option is enabled on the AMI, the login operation requires that a respective
role tag is attached to the EC2 instance that is performing the login.
`
