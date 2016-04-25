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

	configMutex        sync.RWMutex
	tidyCooldownPeriod time.Duration
	nextTidyTime       time.Time

	EC2ClientsMap map[string]*ec2.EC2
}

// periodicFunc performs the tasks that the backend wishes to do periodically.
// Currently this will be triggered once in a minute by the RollbackManager.
// The tasks being done are to cleanup the expired entries of both blacklist
// and whitelist.
func (b *backend) periodicFunc(req *logical.Request) error {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()
	if b.nextTidyTime.IsZero() || !time.Now().Before(b.nextTidyTime) {
		// safety_buffer defaults to 72h
		safety_buffer := 259200
		tidyBlacklistConfigEntry, err := configTidyBlacklistRoleTag(req.Storage)
		if err != nil {
			return err
		}
		skipBlacklistTidy := false
		if tidyBlacklistConfigEntry != nil {
			if tidyBlacklistConfigEntry.DisablePeriodicTidy {
				skipBlacklistTidy = true
			}
			safety_buffer = tidyBlacklistConfigEntry.SafetyBuffer
		}
		if !skipBlacklistTidy {
			tidyBlacklistRoleTag(req.Storage, safety_buffer)
		}

		// reset the safety_buffer to 72h
		safety_buffer = 259200
		tidyWhitelistConfigEntry, err := configTidyWhitelistIdentity(req.Storage)
		if err != nil {
			return err
		}
		skipWhitelistTidy := false
		if tidyWhitelistConfigEntry != nil {
			if tidyWhitelistConfigEntry.DisablePeriodicTidy {
				skipWhitelistTidy = true
			}
			safety_buffer = tidyWhitelistConfigEntry.SafetyBuffer
		}
		if !skipWhitelistTidy {
			tidyWhitelistIdentity(req.Storage, safety_buffer)
		}

		// Update the lastTidyTime
		b.nextTidyTime = time.Now().Add(b.tidyCooldownPeriod)
	}
	return nil
}

const backendHelp = `
AWS auth backend takes in a AWS EC2 instance identity document, its PKCS#7 signature
and a client created nonce to authenticates the instance with Vault.

Authentication is backed by a preconfigured association of AMIs to Vault's policies
through 'image/<ami_id>' endpoint. For instances that share an AMI, an instance tag can
be created through 'image/<ami_id>/tag'. This tag should be attached to the EC2 instance
before the instance attempts to login to Vault.
`
