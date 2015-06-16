package ssh

import (
	"log"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretOneTimeKeyType = "one_time_key"

func secretOneTimeKey(b *backend) *framework.Secret {
	log.Printf("Vishal: ssh.secretPrivateKey\n")
	return &framework.Secret{
		Type: SecretOneTimeKeyType,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username in host",
			},
			"ip": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "ip address of host",
			},
			"one_time_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "SSH one-time-key for host",
			},
		},
		DefaultDuration:    1 * time.Hour,
		DefaultGracePeriod: 10 * time.Minute,
		Renew:              framework.LeaseExtend(1*time.Hour, 0),
		Revoke:             b.secretPrivateKeyRevoke,
	}
}

func (b *backend) secretPrivateKeyRevoke(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.secretPrivateKeyRevoke\n")
	return nil, nil
}
