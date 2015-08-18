package ssh

import (
	"strings"

	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type backend struct {
	*framework.Backend
	salt *salt.Salt
}

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

const backendHelp = `
The SSH backend generates credentials to establish SSH connection with remote hosts.
There are two types of credentials that could be generated: Dynamic and OTP. The
desired way of key creation should be chosen by using 'key_type' parameter of 'roles/'
endpoint. When a credential is requested for a particular role, Vault will generate
a credential accordingly and issue it.

Dynamic Key: is a RSA private key which can be used to establish SSH session using
publickey authentication. When the client receives a key and uses it to establish
connections with hosts, Vault server will have no way to know when and how many 
times the key will be used. So, these login attempts will not be audited by Vault.
To create a dynamic credential, Vault will use the shared private key registered
with the role. Named key should be created using 'keys/' endpoint and used with
'roles/' endpoint for Vault to know the shared key to use for installing the newly
generated key. Since Vault uses the shared key to install keys for other usernames,
shared key should have sudoer privileges in remote hosts and password prompts for
sudoers should be disabled. Also, dynamic keys are leased keys and gets revoked
in remote hosts by Vault after the expiry.

OTP Key: is a UUID which can be used to login using keyboard-interactive authentication.
All the hosts that intend to support OTP should have Vault SSH Agent installed in
them. This agent will receive the OTP from client and get it validated by Vault server.
And since Vault server has a role to play for each successful connection, all the
events will be audited. Vault server validates a key only once, hence it is a OTP.

After mounting this backend, before generating the keys, configure the lease using
'congig/lease' endpoint and create roles using 'roles/' endpoint.
`
