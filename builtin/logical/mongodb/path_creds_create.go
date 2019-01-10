package mongodb

import (
	"context"
	"fmt"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathCredsCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the role to generate credentials for.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathCredsCreateRead,
		},

		HelpSynopsis:    pathCredsCreateReadHelpSyn,
		HelpDescription: pathCredsCreateReadHelpDesc,
	}
}

func (b *backend) pathCredsCreateRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	// Get the role
	role, err := b.Role(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", name)), nil
	}

	// Determine if we have a lease configuration
	leaseConfig, err := b.LeaseConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if leaseConfig == nil {
		leaseConfig = &configLease{}
	}

	// Generate the username and password
	displayName := req.DisplayName
	if displayName != "" {
		displayName += "-"
	}

	userUUID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	username := fmt.Sprintf("vault-%s%s", displayName, userUUID)

	password, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	// Build the user creation command
	createUserCmd := createUserCommand{
		Username: username,
		Password: password,
		Roles:    role.MongoDBRoles.toStandardRolesArray(),
	}

	// Get our connection
	session, err := b.Session(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	// Create the user
	err = session.DB(role.DB).Run(createUserCmd, nil)
	if err != nil {
		return nil, err
	}

	// Return the secret
	resp := b.Secret(SecretCredsType).Response(map[string]interface{}{
		"db":       role.DB,
		"username": username,
		"password": password,
	}, map[string]interface{}{
		"username": username,
		"db":       role.DB,
	})
	resp.Secret.TTL = leaseConfig.TTL
	resp.Secret.MaxTTL = leaseConfig.MaxTTL

	return resp, nil
}

type createUserCommand struct {
	Username string        `bson:"createUser"`
	Password string        `bson:"pwd"`
	Roles    []interface{} `bson:"roles"`
}

const pathCredsCreateReadHelpSyn = `
Request MongoDB database credentials for a particular role.
`

const pathCredsCreateReadHelpDesc = `
This path reads generates MongoDB database credentials for
a particular role. The database credentials will be
generated on demand and will be automatically revoked when
the lease is up.
`
