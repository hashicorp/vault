package rabbitmq

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/michaelklishin/rabbit-hole"
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
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathConnectionWrite,
		},

		HelpSynopsis:    pathConfigConnectionHelpSyn,
		HelpDescription: pathConfigConnectionHelpDesc,
	}
}

func (b *backend) pathConnectionWrite(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	uri := data.Get("connection_uri").(string)
	username := data.Get("username").(string)
	password := data.Get("password").(string)

	if uri == "" {
		return logical.ErrorResponse(fmt.Sprintf(
			"'connection_uri' is a required parameter.")), nil
	}

	if username == "" {
		return logical.ErrorResponse(fmt.Sprintf(
			"'username' is a required parameter.")), nil
	}

	if password == "" {
		return logical.ErrorResponse(fmt.Sprintf(
			"'password' is a required parameter.")), nil
	}

	// Verify the string
	client, err := rabbithole.NewClient(uri, username, password)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error validating connection info: %s", err)), nil
	}

	_, err = client.ListUsers()
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error validating connection info by listing users: %s", err)), nil
	}

	// Store it
	entry, err := logical.StorageEntryJSON("config/connection", connectionConfig{
		URI:      uri,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	// Reset the client connection
	b.ResetClient()

	return nil, nil
}

type connectionConfig struct {
	URI      string `json:"connection_uri"`
	Username string `json:"username"`
	Password string `json:"password"`
}

const pathConfigConnectionHelpSyn = `
Configure the connection URI, username, and password to talk to RabbitMQ management HTTP API.
`

const pathConfigConnectionHelpDesc = `
This path configures the connection properties used to connect to RabbitMQ management HTTP API.
The "connection_uri" parameter is a string that is be used to connect to the API. The "username"
and "password" parameters are strings and used as credentials to the API.

The URI looks like:
"http://localhost:15672"

When configuring the connection URI, username, and password, the backend will verify their validity.
`
