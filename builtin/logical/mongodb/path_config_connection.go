package mongodb

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"gopkg.in/mgo.v2"
)

func pathConfigConnection(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/connection",
		Fields: map[string]*framework.FieldSchema{
			"uri": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "MongoDB standard connection string (URI)",
			},
			"verify_connection": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     true,
				Description: `If set, uri is verified by actually connecting to the database`,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConnectionWrite,
		},
		HelpSynopsis:    pathConfigConnectionHelpSyn,
		HelpDescription: pathConfigConnectionHelpDesc,
	}
}

func (b *backend) pathConnectionWrite(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	uri := data.Get("uri").(string)
	if uri == "" {
		return logical.ErrorResponse("uri parameter must be supplied"), nil
	}

	dialInfo, err := parseMongoURI(uri)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid uri: %s", err)), nil
	}

	// Don't check the config if verification is disabled
	verifyConnection := data.Get("verify_connection").(bool)
	if verifyConnection {
		// Verify the config
		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"Error validating connection info: %s", err)), nil
		}
		defer session.Close()
		if err := session.Ping(); err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"Error validating connection info: %s", err)), nil
		}
	}

	// Store it
	entry, err := logical.StorageEntryJSON("config/connection", connectionConfig{
		URI: uri,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	// Reset the Session
	b.ResetSession()

	return nil, nil
}

type connectionConfig struct {
	URI string `json:"uri"`
}

const pathConfigConnectionHelpSyn = `
Configure the connection string to talk to MongoDB.
`

const pathConfigConnectionHelpDesc = `
This path configures the standard connection string (URI) used to connect to MongoDB.

A MongoDB URI looks like:
"mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]"

See https://docs.mongodb.org/manual/reference/connection-string/ for detailed documentation of the URI format.

When configuring the connection string, the backend will verify its validity.
`
