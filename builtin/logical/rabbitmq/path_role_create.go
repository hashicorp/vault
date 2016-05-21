package rabbitmq

import (
	"fmt"
	"time"

	"github.com/hashicorp/uuid"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/michaelklishin/rabbit-hole"
)

func pathRoleCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathRoleCreateRead,
		},

		HelpSynopsis:    pathRoleCreateReadHelpSyn,
		HelpDescription: pathRoleCreateReadHelpDesc,
	}
}

func (b *backend) pathRoleCreateRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Validate name
	name, err := validateName(data)
	if err != nil {
		return nil, err
	}

	// Get the role
	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", name)), nil
	}

	// Determine if we have a lease
	lease, err := b.Lease(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{Lease: 1 * time.Hour}
	}

	// Ensure username is unique
	username := fmt.Sprintf("%s-%s", req.DisplayName, uuid.GenerateUUID())
	password := uuid.GenerateUUID()

	// Get our connection
	client, err := b.Client(req.Storage)
	if err != nil {
		return nil, err
	}

	if client == nil {
		return logical.ErrorResponse("unable to get client"), nil
	}

	// Create the user
	_, err = client.PutUser(username, rabbithole.UserSettings{
		Password: password,
		Tags:     role.Tags,
	})

	if err != nil {
		return nil, err
	}

	for vhost, permission := range role.VHosts {
		_, err := client.UpdatePermissionsIn(vhost, username, rabbithole.Permissions{
			Configure: permission.Configure,
			Write:     permission.Write,
			Read:      permission.Read,
		})

		if err != nil {
			// Delete the user because it's in an unknown state
			_, rmErr := client.DeleteUser(username)
			if rmErr != nil {
				return logical.ErrorResponse(fmt.Sprintf("failed to update user: %s, failed to delete user: %s, user: %s", err, rmErr, username)), rmErr
			}
			return logical.ErrorResponse(fmt.Sprintf("failed to update user: %s, 	  user: %s", err, username)), err
		}
	}

	// Return the secret
	resp := b.Secret(SecretCredsType).Response(map[string]interface{}{
		"username": username,
		"password": password,
	}, map[string]interface{}{
		"username": username,
	})
	resp.Secret.TTL = lease.Lease
	return resp, nil
}

const pathRoleCreateReadHelpSyn = `
Request RabbitMQ credentials for a certain role.
`

const pathRoleCreateReadHelpDesc = `
This path reads RabbitMQ credentials for a certain role. The
RabbitMQ credentials will be generated on demand and will be automatically
revoked when the lease is up.
`
