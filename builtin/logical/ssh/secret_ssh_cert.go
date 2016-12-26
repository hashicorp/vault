package ssh

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// SecretCertsType is the name used to identify this type
const SecretCertsType = "secret_ssh_ca"

func secretCerts(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretCertsType,
		Fields: map[string]*framework.FieldSchema{
			"signed_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The signd certificate.",
			},
			"serial_number": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The serial number of the certificate, for handy
reference`,
			},
		},

		Revoke: b.secretCredsRevoke,
	}
}

func (b *backend) secretCredsRevoke(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// this backend doesn't support CRL, so there's nothing that can be done when a certificate is revoked
	return &logical.Response{}, nil
}
