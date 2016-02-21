package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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
This name is deprecated.`,
			},
			"max_open_connections": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Maximum number of open connections to database",
			},
			"verify_connection": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     true,
				Description: "If set, connection_url is verified by actually connecting to the database",
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
	connValue := data.Get("value").(string)
	connURL := data.Get("connection_url").(string)
	if connURL == "" {
		if connValue == "" {
			return logical.ErrorResponse("the connection_url parameter must be supplied"), nil
		} else {
			connURL = connValue
		}
	}

	maxOpenConns := data.Get("max_open_connections").(int)
	if maxOpenConns == 0 {
		maxOpenConns = 2
	}

	// Don't check the connection_url if verification is disabled
	verifyConnection := data.Get("verify_connection").(bool)
	if verifyConnection {
		// Verify the string
		db, err := sql.Open("mysql", connURL)

		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"error validating connection info: %s", err)), nil
		}
		defer db.Close()
		if err := db.Ping(); err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"error validating connection info: %s", err)), nil
		}
	}

	// Store it
	entry, err := logical.StorageEntryJSON("config/connection", connectionConfig{
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
Configure the connection string to talk to MySQL.
`

const pathConfigConnectionHelpDesc = `
This path configures the connection string used to connect to MySQL.  The value
of the string is a Data Source Name (DSN). An example is using
"username:password@protocol(address)/dbname?param=value"

For example, RDS may look like:
"id:password@tcp(your-amazonaws-uri.com:3306)/dbname"

When configuring the connection string, the backend will verify its validity.
If the database is not available when setting the connection URL, set the
"verify_connection" option to false.
`
