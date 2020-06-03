package rabbitmq

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	rabbithole "github.com/michaelklishin/rabbit-hole"
)

const (
	storageKey = "config/connection"
)

func pathConfigConnection(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/connection",
		Fields: map[string]*framework.FieldSchema{
			"connection_uri": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "RabbitMQ Management URI",
			},
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username of a RabbitMQ management administrator",
			},
			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password of the provided RabbitMQ management user",
			},
			"verify_connection": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     true,
				Description: `If set, connection_uri is verified by actually connecting to the RabbitMQ management API`,
			},
			"password_policy": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the password policy to use to generate passwords for dynamic credentials.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConnectionUpdate,
		},

		HelpSynopsis:    pathConfigConnectionHelpSyn,
		HelpDescription: pathConfigConnectionHelpDesc,
	}
}

func (b *backend) pathConnectionUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	uri := data.Get("connection_uri").(string)
	if uri == "" {
		return logical.ErrorResponse("missing connection_uri"), nil
	}

	username := data.Get("username").(string)
	if username == "" {
		return logical.ErrorResponse("missing username"), nil
	}

	password := data.Get("password").(string)
	if password == "" {
		return logical.ErrorResponse("missing password"), nil
	}

	passwordPolicy := data.Get("password_policy").(string)

	// Don't check the connection_url if verification is disabled
	verifyConnection := data.Get("verify_connection").(bool)
	if verifyConnection {
		// Create RabbitMQ management client
		client, err := rabbithole.NewClient(uri, username, password)
		if err != nil {
			return nil, errwrap.Wrapf("failed to create client: {{err}}", err)
		}

		// Verify that configured credentials is capable of listing
		if _, err = client.ListUsers(); err != nil {
			return nil, errwrap.Wrapf("failed to validate the connection: {{err}}", err)
		}
	}

	// Store it
	config := connectionConfig{
		URI:            uri,
		Username:       username,
		Password:       password,
		PasswordPolicy: passwordPolicy,
	}
	err := writeConfig(ctx, req.Storage, config)
	if err != nil {
		return nil, err
	}

	// Reset the client connection
	b.resetClient(ctx)

	return nil, nil
}

func readConfig(ctx context.Context, storage logical.Storage) (config connectionConfig, err error) {
	entry, err := storage.Get(ctx, storageKey)
	if err != nil {
		return connectionConfig{}, err
	}
	if entry == nil {
		return connectionConfig{}, nil
	}

	var connConfig connectionConfig
	if err := entry.DecodeJSON(&connConfig); err != nil {
		return connectionConfig{}, err
	}
	return connConfig, nil
}

func writeConfig(ctx context.Context, storage logical.Storage, config connectionConfig) (err error) {
	entry, err := logical.StorageEntryJSON(storageKey, config)
	if err != nil {
		return err
	}
	if err := storage.Put(ctx, entry); err != nil {
		return err
	}
	return nil
}

// connectionConfig contains the information required to make a connection to a RabbitMQ node
type connectionConfig struct {
	// URI of the RabbitMQ server
	URI string `json:"connection_uri"`

	// Username which has 'administrator' tag attached to it
	Username string `json:"username"`

	// Password for the Username
	Password string `json:"password"`

	// PasswordPolicy for generating passwords for dynamic credentials
	PasswordPolicy string `json:"password_policy"`
}

const pathConfigConnectionHelpSyn = `
Configure the connection URI, username, and password to talk to RabbitMQ management HTTP API.
`

const pathConfigConnectionHelpDesc = `
This path configures the connection properties used to connect to RabbitMQ management HTTP API.
The "connection_uri" parameter is a string that is used to connect to the API. The "username"
and "password" parameters are strings that are used as credentials to the API. The "verify_connection"
parameter is a boolean that is used to verify whether the provided connection URI, username, and password
are valid.

The URI looks like:
"http://localhost:15672"
`
