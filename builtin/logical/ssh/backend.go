package ssh

import (
	"strings"

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
	salt, err := salt.NewSalt(conf.View, nil)
	if err != nil {
		return nil, err
	}

	var b backend
	b.salt = salt
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			Root: []string{
				"config/*",
				"keys/*",
			},
			Unauthenticated: []string{
				"verify",
			},
		},

		Paths: []*framework.Path{
			pathConfigLease(&b),
			pathKeys(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
			pathLookup(&b),
			pathVerify(&b),
		},

		Secrets: []*framework.Secret{
			secretDynamicKey(&b),
			secretOTP(&b),
		},
	}
	return b.Backend, nil
}

type backend struct {
	*framework.Backend
	salt *salt.Salt
}

const backendHelp = `
The SSH backend generates keys to eatablish SSH connection
with remote hosts. There are two options to create the keys:
long lived dynamic key, one time password. 

Long lived dynamic key is a rsa private key which can be used
to login to remote host using the publickey authentication.
There is no additional change required in the remote hosts to
support this type of keys. But the keys generated will be valid
as long as the lease of the key is valid. Also, logins to remote
hosts will not be audited in vault server.

One Time Password (OTP), on the other hand is a randomly generated
UUID that is used to login to remote host using the keyboard-
interactive challenge response authentication. A vault agent
has to be installed at the remote host to support OTP. Upon 
request, vault server generates and provides the key to the
user. During login, vault agent receives the key and verifies
the correctness with the vault server (and hence audited). The
server after verifying the key for the first time, deletes the
same (and hence one-time).

Both type of keys have a configurable lease set and are automatically
revoked at the end of the lease.

After mounting this backend, before generating the keys, configure
the lease using the 'config/lease' endpoint and create roles using
the 'roles/' endpoint.
`
