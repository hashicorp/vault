package mssql

import (
	"database/sql"
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigConnection(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/connection",
		Fields: map[string]*framework.FieldSchema{
			"connection_string": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "DB connection parameters",
			},
			"max_open_connections": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Maximum number of open connections to database",
			},
			"verify_connection": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     true,
				Description: "If set, connection_string is verified by actually connecting to the database",
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

// pathConnectionWrite stores the connection configuration
func (b *backend) pathConnectionWrite(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	connString := data.Get("connection_string").(string)

	maxOpenConns := data.Get("max_open_connections").(int)
	if maxOpenConns == 0 {
		maxOpenConns = 2
	}

	// Don't check the connection_string if verification is disabled
	verifyConnection := data.Get("verify_connection").(bool)
	if verifyConnection {
		// Verify the string
		db, err := sql.Open("mssql", connString)

		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"Error validating connection info: %s", err)), nil
		}
		defer db.Close()
		if err := db.Ping(); err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"Error validating connection info: %s", err)), nil
		}
	}

	// Store it
	entry, err := logical.StorageEntryJSON("config/connection", connectionConfig{
		ConnectionString:   connString,
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

	resp := &logical.Response{}
	resp.AddWarning("Read access to this endpoint should be controlled via ACLs as it will return the connection string as it is, including passwords, if any.")

	return resp, nil
}

type connectionConfig struct {
	ConnectionString   string `json:"connection_string" structs:"connection_string" mapstructure:"connection_string"`
	MaxOpenConnections int    `json:"max_open_connections" structs:"max_open_connections" mapstructure:"max_open_connections"`
}

const pathConfigConnectionHelpSyn = `
Configure the connection string to talk to Microsoft Sql Server.
`

const pathConfigConnectionHelpDesc = `
This path configures the connection string used to connect to Sql Server.
The value of the string is a Data Source Name (DSN). An example is
using "server=<hostname>;port=<port>;user id=<username>;password=<password>;database=<database>;app name=vault;"

When configuring the connection string, the backend will verify its validity.
If the database is not available when setting the connection string, set the
"verify_connection" option to false.
`
