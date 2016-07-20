package mysql

import (
	"database/sql"
	"fmt"

	"github.com/fatih/structs"
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
				Description: `DB connection string. Use 'connection_url' instead.
This name is deprecated.`,
			},
			"max_open_connections": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Maximum number of open connections to database",
			},
			"max_idle_connections": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Maximum number of idle connections to the database; a zero uses the value of max_open_connections and a negative value disables idle connections. If larger than max_open_connections it will be reduced to the same size.",
			},
			"verify_connection": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     true,
				Description: "If set, connection_url is verified by actually connecting to the database",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConnectionWrite,
			logical.ReadOperation:   b.pathConnectionRead,
		},

		HelpSynopsis:    pathConfigConnectionHelpSyn,
		HelpDescription: pathConfigConnectionHelpDesc,
	}
}

// pathConnectionRead reads out the connection configuration
func (b *backend) pathConnectionRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get("config/connection")
	if err != nil {
		return nil, fmt.Errorf("failed to read connection configuration")
	}
	if entry == nil {
		return nil, nil
	}

	var config connectionConfig
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: structs.New(config).Map(),
	}, nil
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

	maxIdleConns := data.Get("max_idle_connections").(int)
	if maxIdleConns == 0 {
		maxIdleConns = maxOpenConns
	}
	if maxIdleConns > maxOpenConns {
		maxIdleConns = maxOpenConns
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
		MaxIdleConnections: maxIdleConns,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	// Reset the DB connection
	b.ResetDB()

	resp := &logical.Response{}
	resp.AddWarning("Read access to this endpoint should be controlled via ACLs as it will return the connection URL as it is, including passwords, if any.")

	return resp, nil
}

type connectionConfig struct {
	ConnectionURL string `json:"connection_url" structs:"connection_url" mapstructure:"connection_url"`
	// Deprecate "value" in coming releases
	ConnectionString   string `json:"value" structs:"value" mapstructure:"value"`
	MaxOpenConnections int    `json:"max_open_connections" structs:"max_open_connections" mapstructure:"max_open_connections"`
	MaxIdleConnections int    `json:"max_idle_connections" structs:"max_idle_connections" mapstructure:"max_idle_connections"`
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
