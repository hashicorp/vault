package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	_ "github.com/lib/pq"
)

func pathConfigConnection(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/connection",
		Fields: map[string]*framework.FieldSchema{
			"connection_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "DB connection string",
			},
			"value": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				DB connection string. Use 'connection_url' instead.
				This will be deprecated.`,
			},
			"max_open_connections": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Maximum number of open connections to the database",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConnectionWrite,
		},

		HelpSynopsis:    pathConfigConnectionHelpSyn,
		HelpDescription: pathConfigConnectionHelpDesc,
	}
}

func (b *backend) pathConnectionWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	connString := data.Get("value").(string)
	connURL := data.Get("connection_url").(string)

	maxOpenConns := data.Get("max_open_connections").(int)
	if maxOpenConns == 0 {
		maxOpenConns = 2
	}

	// Verify the string
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error validating connection info: %s", err)), nil
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error validating connection info: %s", err)), nil
	}

	// Store it
	entry, err := logical.StorageEntryJSON("config/connection", connectionConfig{
		ConnectionString:   connString,
		ConnectionURL:      connURL,
		MaxOpenConnections: maxOpenConns,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	// Reset the DB connection
	b.ResetDB()

	return nil, nil
}

type connectionConfig struct {
	ConnectionURL string `json:"connection_url"`
	// Deprecate "value" in coming releases
	ConnectionString   string `json:"value"`
	MaxOpenConnections int    `json:"max_open_connections"`
}

const pathConfigConnectionHelpSyn = `
Configure the connection string to talk to PostgreSQL.
`

const pathConfigConnectionHelpDesc = `
This path configures the connection string used to connect to PostgreSQL.
The value of the string can be a URL, or a PG style string in the
format of "user=foo host=bar" etc.

The URL looks like:
"postgresql://user:pass@host:port/dbname"

When configuring the connection string, the backend will verify its validity.
`
