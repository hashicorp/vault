package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/fatih/structs"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

var (
	respErrEmptyPluginName = "empty plugin name"
	respErrEmptyName       = "empty name attribute given"
)

// DatabaseConfig is used by the Factory function to configure a Database
// object.
type DatabaseConfig struct {
	PluginName string `json:"plugin_name" structs:"plugin_name" mapstructure:"plugin_name"`
	// ConnectionDetails stores the database specific connection settings needed
	// by each database type.
	ConnectionDetails map[string]interface{} `json:"connection_details" structs:"connection_details" mapstructure:"connection_details"`
	AllowedRoles      []string               `json:"allowed_roles" structs:"allowed_roles" mapstructure:"allowed_roles"`

	RootCredentialsRotateStatements []string `json:"root_credentials_rotate_statements" structs:"root_credentials_rotate_statements" mapstructure:"root_credentials_rotate_statements"`
}

// pathResetConnection configures a path to reset a plugin.
func pathResetConnection(b *databaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("reset/%s", framework.GenericNameRegex("name")),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of this database connection",
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
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse(respErrEmptyName), nil
		}

		// Close plugin and delete the entry in the connections cache.
		if err := b.ClearConnection(name); err != nil {
			return nil, err
		}

		// Execute plugin again, we don't need the object so throw away.
		if _, err := b.GetConnection(ctx, req.Storage, name); err != nil {
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
				Description: "Name of this database connection",
			},

			"plugin_name": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The name of a builtin or previously registered
				plugin known to vault. This endpoint will create an instance of
				that plugin type.`,
			},

			"verify_connection": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: true,
				Description: `If true, the connection details are verified by
				actually connecting to the database. Defaults to true.`,
			},

			"allowed_roles": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `Comma separated string or array of the role names
				allowed to get creds from this database connection. If empty no
				roles are allowed. If "*" all roles are allowed.`,
			},

			"root_rotation_statements": &framework.FieldSchema{
				Type: framework.TypeStringSlice,
				Description: `Specifies the database statements to be executed
				to rotate the root user's credentials. See the plugin's API 
				page for more information on support and formatting for this 
				parameter.`,
			},
		},

		ExistenceCheck: b.connectionExistenceCheck(),
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.connectionWriteHandler(),
			logical.UpdateOperation: b.connectionWriteHandler(),
			logical.ReadOperation:   b.connectionReadHandler(),
			logical.DeleteOperation: b.connectionDeleteHandler(),
		},

		HelpSynopsis:    pathConfigConnectionHelpSyn,
		HelpDescription: pathConfigConnectionHelpDesc,
	}
}

func (b *databaseBackend) connectionExistenceCheck() framework.ExistenceFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
		name := data.Get("name").(string)
		if name == "" {
			return false, errors.New(`missing "name" parameter`)
		}

		entry, err := req.Storage.Get(ctx, fmt.Sprintf("config/%s", name))
		if err != nil {
			return false, errors.New("failed to read connection configuration")
		}

		return entry != nil, nil
	}
}

func pathListPluginConnection(b *databaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("config/?$"),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.connectionListHandler(),
		},

		HelpSynopsis:    pathConfigConnectionHelpSyn,
		HelpDescription: pathConfigConnectionHelpDesc,
	}
}

func (b *databaseBackend) connectionListHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		entries, err := req.Storage.List(ctx, "config/")
		if err != nil {
			return nil, err
		}

		return logical.ListResponse(entries), nil
	}
}

// connectionReadHandler reads out the connection configuration
func (b *databaseBackend) connectionReadHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse(respErrEmptyName), nil
		}

		entry, err := req.Storage.Get(ctx, fmt.Sprintf("config/%s", name))
		if err != nil {
			return nil, errors.New("failed to read connection configuration")
		}
		if entry == nil {
			return nil, nil
		}

		var config DatabaseConfig
		if err := entry.DecodeJSON(&config); err != nil {
			return nil, err
		}

		// Mask the password if it is in the url
		if connURLRaw, ok := config.ConnectionDetails["connection_url"]; ok {
			connURL := connURLRaw.(string)
			if conn, err := url.Parse(connURL); err == nil {
				if password, ok := conn.User.Password(); ok {
					config.ConnectionDetails["connection_url"] = strings.Replace(connURL, password, "*****", -1)
				}
			}
		}

		delete(config.ConnectionDetails, "password")

		return &logical.Response{
			Data: structs.New(config).Map(),
		}, nil
	}
}

// connectionDeleteHandler deletes the connection configuration
func (b *databaseBackend) connectionDeleteHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse(respErrEmptyName), nil
		}

		err := req.Storage.Delete(ctx, fmt.Sprintf("config/%s", name))
		if err != nil {
			return nil, errors.New("failed to delete connection configuration")
		}

		if err := b.ClearConnection(name); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// connectionWriteHandler returns a handler function for creating and updating
// both builtin and plugin database types.
func (b *databaseBackend) connectionWriteHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		verifyConnection := data.Get("verify_connection").(bool)

		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse(respErrEmptyName), nil
		}

		// Baseline
		config := &DatabaseConfig{}

		entry, err := req.Storage.Get(ctx, fmt.Sprintf("config/%s", name))
		if err != nil {
			return nil, errors.New("failed to read connection configuration")
		}
		if entry != nil {
			if err := entry.DecodeJSON(config); err != nil {
				return nil, err
			}
		}

		if pluginNameRaw, ok := data.GetOk("plugin_name"); ok {
			config.PluginName = pluginNameRaw.(string)
		} else if req.Operation == logical.CreateOperation {
			config.PluginName = data.Get("plugin_name").(string)
		}
		if config.PluginName == "" {
			return logical.ErrorResponse(respErrEmptyPluginName), nil
		}

		if allowedRolesRaw, ok := data.GetOk("allowed_roles"); ok {
			config.AllowedRoles = allowedRolesRaw.([]string)
		} else if req.Operation == logical.CreateOperation {
			config.AllowedRoles = data.Get("allowed_roles").([]string)
		}

		if rootRotationStatementsRaw, ok := data.GetOk("root_rotation_statements"); ok {
			config.RootCredentialsRotateStatements = rootRotationStatementsRaw.([]string)
		} else if req.Operation == logical.CreateOperation {
			config.RootCredentialsRotateStatements = data.Get("root_rotation_statements").([]string)
		}

		// Remove these entries from the data before we store it keyed under
		// ConnectionDetails.
		delete(data.Raw, "name")
		delete(data.Raw, "plugin_name")
		delete(data.Raw, "allowed_roles")
		delete(data.Raw, "verify_connection")
		delete(data.Raw, "root_rotation_statements")

		// Create a database plugin and initialize it.
		db, err := dbplugin.PluginFactory(ctx, config.PluginName, b.System(), b.logger)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("error creating database object: %s", err)), nil
		}

		// If this is an update, take any new values, overwrite what was there
		// before, and pass that in as the "new" set of values to the plugin,
		// then save what results
		if req.Operation == logical.CreateOperation {
			config.ConnectionDetails = data.Raw
		} else {
			if config.ConnectionDetails == nil {
				config.ConnectionDetails = make(map[string]interface{})
			}
			for k, v := range data.Raw {
				config.ConnectionDetails[k] = v
			}
		}

		config.ConnectionDetails, err = db.Init(ctx, config.ConnectionDetails, verifyConnection)
		if err != nil {
			db.Close()
			return logical.ErrorResponse(fmt.Sprintf("error creating database object: %s", err)), nil
		}

		b.Lock()
		defer b.Unlock()

		// Close and remove the old connection
		b.clearConnection(name)

		id, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}

		b.connections[name] = &dbPluginInstance{
			Database: db,
			name:     name,
			id:       id,
		}

		// Store it
		entry, err = logical.StorageEntryJSON(fmt.Sprintf("config/%s", name), config)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		resp := &logical.Response{}

		// This is a simple test to to check for passwords in the connection_url paramater. If one exists,
		// warn the user to use templated url string
		if connURLRaw, ok := config.ConnectionDetails["connection_url"]; ok {
			if connURL, err := url.Parse(connURLRaw.(string)); err == nil {
				if _, ok := connURL.User.Password(); ok {
					resp.AddWarning("Password found in connection_url, use a templated url to enable root rotation and prevent read access to password information.")
				}
			}
		}

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

In addition to the database specific connection details, this endpoint also
accepts:

	* "plugin_name" (required) - The name of a builtin or previously registered
	   plugin known to vault. This endpoint will create an instance of that
	   plugin type.

	* "verify_connection" (default: true) - A boolean value denoting if the plugin should verify
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
