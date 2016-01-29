package pki

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// SecretCertsType is the name used to identify this type
const SecretCertsType = "pki"

func secretCerts(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretCertsType,
		Fields: map[string]*framework.FieldSchema{
			"certificate": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The PEM-encoded concatenated certificate and
issuing certificate authority`,
			},
			"private_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The PEM-encoded private key for the certificate",
			},
			"serial": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The serial number of the certificate, for handy
reference`,
			},
		},

		Revoke: b.secretCredsRevoke,
	}
}

func (b *backend) secretCredsRevoke(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if req.Secret == nil {
		return nil, fmt.Errorf("Secret is nil in request")
	}

	serialInt, ok := req.Secret.InternalData["serial_number"]
	if !ok {
		return nil, fmt.Errorf("Could not find serial in internal secret data")
	}

	serial := strings.Replace(strings.ToLower(serialInt.(string)), "-", ":", -1)

	b.revokeStorageLock.Lock()
	defer b.revokeStorageLock.Unlock()

	return revokeCert(b, req, serial)
}
