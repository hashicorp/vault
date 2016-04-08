package aws

import (
	"sync"

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

	var b backend
	b.Salt = salt
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
		},

		Paths: append([]*framework.Path{
			pathLogin(&b),
			pathImage(&b),
			pathListImages(&b),
			pathImageTag(&b),
			pathConfigClient(&b),
			pathConfigCertificate(&b),
			pathBlacklistRoleTag(&b),
			pathListBlacklistRoleTags(&b),
			pathBlacklistRoleTagTidy(&b),
			pathWhitelistIdentity(&b),
			pathWhitelistIdentityTidy(&b),
			pathListWhitelistIdentities(&b),
		}),

		AuthRenew: b.pathLoginRenew,
	}

	return b.Backend, nil
}

type backend struct {
	*framework.Backend
	Salt *salt.Salt

	configMutex sync.RWMutex

	ec2Client *ec2.EC2
}

const backendHelp = `
AWS auth backend takes in a AWS EC2 instance identity document, its PKCS#7 signature
and a client created nonce to authenticates the instance with Vault.

Authentication is backed by a preconfigured association of AMIs to Vault's policies
through 'image/<ami_id>' endpoint. For instances that share an AMI, an instance tag can
be created through 'image/<ami_id>/tag'. This tag should be attached to the EC2 instance
before the instance attempts to login to Vault.
`
