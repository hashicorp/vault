package rabbitmq

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// SecretCredsType is the key for this backend's secrets.
const SecretCredsType = "creds"

func secretCreds(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretCredsType,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password",
			},
		},

		Renew:  b.secretCredsRenew,
		Revoke: b.secretCredsRevoke,
	}
}

func (b *backend) secretCredsRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the lease information
	lease, err := b.Lease(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}

	f := framework.LeaseExtend(lease.Lease, lease.LeaseMax, b.System())
	resp, err := f(req, d)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *backend) secretCredsRevoke(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username := usernameRaw.(string)

	// Get our connection
	client, err := b.Client(req.Storage)
	if err != nil {
		return nil, err
	}

	_, err = client.DeleteUser(username)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("could not delete user: %s", err)), nil
	}

	return nil, nil
}
