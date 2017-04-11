package database

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathResetConnection configures a path to reset a plugin.
func pathResetConnection(b *databaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("reset/%s", framework.GenericNameRegex("name")),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of this DB type",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConnectionReset(),
		},

		HelpSynopsis:    pathResetConnectionHelpSyn,
		HelpDescription: pathResetConnectionHelpDesc,
	}
}

// pathConnectionReset resets a plugin by closing the existing instance and
// creating a new one.
func (b *databaseBackend) pathConnectionReset() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse("Empty name attribute given"), nil
		}

		// Grab the mutex lock
		b.Lock()
		defer b.Unlock()

		b.clearConnection(name)

		_, err := b.getOrCreateDBObj(req.Storage, name)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// pathConfigurePluginConnection returns a configured framework.Path setup to
// operate on plugins.
func pathConfigurePluginConnection(b *databaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("config/%s", framework.GenericNameRegex("name")),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of this DB type",
			},

			"verify_connection": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: true,
				Description: `If set, the connection details are verified by
							actually connecting to the database`,
			},

			"plugin_name": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The name of a builtin or previously registered
							plugin known to vault. This endpoint will create an instance of
							that plugin type.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.connectionWriteHandler(),
			logical.ReadOperation:   b.connectionReadHandler(),
			logical.DeleteOperation: b.connectionDeleteHandler(),
		},

		HelpSynopsis:    pathConfigConnectionHelpSyn,
		HelpDescription: pathConfigConnectionHelpDesc,
	}
}

// pathConnectionRead reads out the connection configuration
func (b *databaseBackend) connectionReadHandler() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)

		entry, err := req.Storage.Get(fmt.Sprintf("dbs/%s", name))
		if err != nil {
			return nil, fmt.Errorf("failed to read connection configuration")
		}
		if entry == nil {
			return nil, nil
		}

		var config DatabaseConfig
		if err := entry.DecodeJSON(&config); err != nil {
			return nil, err
		}
		return &logical.Response{
			Data: structs.New(config).Map(),
		}, nil
	}
}

// connectionDeleteHandler deletes the connection configuration
func (b *databaseBackend) connectionDeleteHandler() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse("Empty name attribute given"), nil
		}

		err := req.Storage.Delete(fmt.Sprintf("dbs/%s", name))
		if err != nil {
			return nil, fmt.Errorf("failed to delete connection configuration")
		}

		b.Lock()
		defer b.Unlock()

		if _, ok := b.connections[name]; ok {
			err = b.connections[name].Close()
			if err != nil {
				return nil, err
			}
		}

		delete(b.connections, name)

		return nil, nil
	}
}

// connectionWriteHandler returns a handler function for creating and updating
// both builtin and plugin database types.
func (b *databaseBackend) connectionWriteHandler() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

		config := &DatabaseConfig{
			ConnectionDetails: data.Raw,
			PluginName:        data.Get("plugin_name").(string),
		}

		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse("Empty name attribute given"), nil
		}

		verifyConnection := data.Get("verify_connection").(bool)

		// Grab the mutex lock
		b.Lock()
		defer b.Unlock()

		db, err := dbplugin.PluginFactory(config.PluginName, b.System(), b.logger)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Error creating database object: %s", err)), nil
		}

		err = db.Initialize(config.ConnectionDetails, verifyConnection)
		if err != nil {
			db.Close()
			return logical.ErrorResponse(fmt.Sprintf("Error creating database object: %s", err)), nil
		}

		if _, ok := b.connections[name]; ok {
			// Close and remove the old connection
			err := b.connections[name].Close()
			if err != nil {
				db.Close()
				return nil, err
			}

			delete(b.connections, name)
		}

		// Save the new connection
		b.connections[name] = db

		// Store it
		entry, err := logical.StorageEntryJSON(fmt.Sprintf("dbs/%s", name), config)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(entry); err != nil {
			return nil, err
		}

		resp := &logical.Response{}
		resp.AddWarning("Read access to this endpoint should be controlled via ACLs as it will return the connection string or URL as it is, including passwords, if any.")

		return resp, nil
	}
}

const pathConfigConnectionHelpSyn = `
Configure connection details to a database plugin.
`

const pathConfigConnectionHelpDesc = `
This path configures the connection details used to connect to a particular
database. This path runs the provided plugin name and passes the configured
connection details to the plugin. See the documentation for the plugin specified
for a full list of accepted connection details. 

In addition to the database specific connection details, this endpoing also
accepts:

	* "plugin_name" (required) - The name of a builtin or previously registered
	   plugin known to vault. This endpoint will create an instance of that
	   plugin type.

	* "verify_connection" - A boolean value denoting if the plugin should verify
	   it is able to connect to the database using the provided connection
       details.
`

const pathResetConnectionHelpSyn = `
Resets a database plugin.
`

const pathResetConnectionHelpDesc = `
This path resets the database connection by closing the existing database plugin
instance and running a new one.
`
