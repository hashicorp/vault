package cassandra

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// SecretCredsType is the type of creds issued from this backend
const SecretCredsType = "cassandra"

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
	roleRaw, ok := req.Secret.InternalData["role"]
	if !ok {
		return nil, fmt.Errorf("secret is missing role internal data")
	}
	roleName, ok := roleRaw.(string)
	if !ok {
		return nil, fmt.Errorf("error converting role internal data to string")
	}

	role, err := getRole(req.Storage, roleName)
	if err != nil {
		return nil, fmt.Errorf("unable to load role: %s", err)
	}

	return framework.LeaseExtend(role.Lease, 0, b.System())(req, d)
}

func (b *backend) secretCredsRevoke(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("error converting username internal data to string")
	}

	session, err := b.DB(req.Storage)
	if err != nil {
		return nil, fmt.Errorf("error getting session")
	}

	err = session.Query(fmt.Sprintf("DROP USER '%s'", username)).Exec()
	if err != nil {
		return nil, fmt.Errorf("error removing user %s", username)
	}

	return nil, nil
}
