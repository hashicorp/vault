package ssh

import (
	"log"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretSshHostKeyType = "secret_ssh_host_key_type"

func secretOneTimeKey(b *backend) *framework.Secret {
	log.Printf("Vishal: ssh.secretPrivateKey\n")
	return &framework.Secret{
		Type: SecretSshHostKeyType,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username in host",
			},
			"ip": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "ip address of host",
			},
		},
		DefaultDuration:    1 * time.Hour,
		DefaultGracePeriod: 10 * time.Minute,
		Renew:              b.secretPrivateKeyRenew,
		Revoke:             b.secretPrivateKeyRevoke,
	}
}

func (b *backend) secretPrivateKeyRenew(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.secretPrivateKeyRenew\n")
	lease, err := b.Lease(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{Lease: 1 * time.Hour}
	}
	f := framework.LeaseExtend(lease.Lease, lease.LeaseMax, false)
	return f(req, d)
}

func (b *backend) secretPrivateKeyRevoke(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.secretPrivateKeyRevoke\n")
	//TODO: implement here
	return nil, nil
}
