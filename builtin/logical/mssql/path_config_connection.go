package mssql

import (
	"database/sql"
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigConnection(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/connection",
		Fields: map[string]*framework.FieldSchema{
			"connection_params": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "DB connection parameters",
			},
			"max_open_connections": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Maximum number of open connections to database",
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
	connParams := data.Get("connection_params").(string)

	maxOpenConns := data.Get("max_open_connections").(int)
	if maxOpenConns == 0 {
		maxOpenConns = 2
	}

	// Verify the string
	db, err := sql.Open("mssql", BuildDsn(connParams))

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
		ConnectionParams:   connParams,
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
	ConnectionParams   string `json:"connection_params"`
	MaxOpenConnections int    `json:"max_open_connections"`
}

const pathConfigConnectionHelpSyn = `
Configure the connection string to talk to Microsoft Sql Server.
`

const pathConfigConnectionHelpDesc = `
This path configures the connection string used to connect to Sql Server.
The value of the string is a Data Source Name (DSN). An example is
using "server=<hostname>;port=<port>;user id=<username>;password=<password>;database=<database>;"

When configuring the connection string, the backend will verify its validity.
`
