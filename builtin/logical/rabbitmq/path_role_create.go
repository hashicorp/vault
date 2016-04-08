package rabbitmq

import (
	"fmt"

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
	name := data.Get("name").(string)

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
		lease = &configLease{}
	}

	displayName := req.DisplayName
	if len(displayName) > 26 {
		displayName = displayName[:26]
	}
	userUUID := uuid.GenerateUUID()
	username := fmt.Sprintf("%s-%s", displayName, userUUID)
	if len(username) > 63 {
		username = username[:63]
	}
	password := uuid.GenerateUUID()

	// Get our connection
	client, err := b.Client(req.Storage)
	if err != nil {
		return nil, err
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
			return nil, err
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
