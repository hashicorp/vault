package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/fatih/structs"
	"github.com/hashicorp/go-uuid"

	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
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

	PasswordPolicy string `json:"password_policy" structs:"password_policy" mapstructure:"password_policy"`
}

func (c *DatabaseConfig) SupportsCredentialType(credentialType v5.CredentialType) bool {
	credTypes, ok := c.ConnectionDetails[v5.SupportedCredentialTypesKey].([]interface{})
	if !ok {
		// Default to supporting CredentialTypePassword for database plugins that
		// don't specify supported credential types in the initialization response
		return credentialType == v5.CredentialTypePassword
	}

	for _, ct := range credTypes {
		if ct == credentialType.String() {
			return true
		}
	}
	return false
}

// pathResetConnection configures a path to reset a plugin.
func pathResetConnection(b *databaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("reset/%s", framework.GenericNameRegex("name")),
		Fields: map[string]*framework.FieldSchema{
			"name": {
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
			"name": {
				Type:        framework.TypeString,
				Description: "Name of this database connection",
			},

			"plugin_name": {
				Type: framework.TypeString,
				Description: `The name of a builtin or previously registered
				plugin known to vault. This endpoint will create an instance of
				that plugin type.`,
			},

			"verify_connection": {
				Type:    framework.TypeBool,
				Default: true,
				Description: `If true, the connection details are verified by
				actually connecting to the database. Defaults to true.`,
			},

			"allowed_roles": {
				Type: framework.TypeCommaStringSlice,
				Description: `Comma separated string or array of the role names
				allowed to get creds from this database connection. If empty no
				roles are allowed. If "*" all roles are allowed.`,
			},

			"root_rotation_statements": {
				Type: framework.TypeStringSlice,
				Description: `Specifies the database statements to be executed
				to rotate the root user's credentials. See the plugin's API 
				page for more information on support and formatting for this 
				parameter.`,
			},
			"password_policy": {
				Type:        framework.TypeString,
				Description: `Password policy to use when generating passwords.`,
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
			return false, fmt.Errorf("failed to read connection configuration: %w", err)
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
			return nil, fmt.Errorf("failed to read connection configuration: %w", err)
		}
		if entry == nil {
			return nil, nil
		}

		var config DatabaseConfig
		if err := entry.DecodeJSON(&config); err != nil {
			return nil, err
		}

		// Ensure that we only ever include a redacted valid URL in the response.
		if connURLRaw, ok := config.ConnectionDetails["connection_url"]; ok {
			if p, err := url.Parse(connURLRaw.(string)); err == nil {
				config.ConnectionDetails["connection_url"] = p.Redacted()
			}
		}

		delete(config.ConnectionDetails, "password")
		delete(config.ConnectionDetails, "private_key")

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
			return nil, fmt.Errorf("failed to delete connection configuration: %w", err)
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
			return nil, fmt.Errorf("failed to read connection configuration: %w", err)
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

		if passwordPolicyRaw, ok := data.GetOk("password_policy"); ok {
			config.PasswordPolicy = passwordPolicyRaw.(string)
		}

		// Remove these entries from the data before we store it keyed under
		// ConnectionDetails.
		delete(data.Raw, "name")
		delete(data.Raw, "plugin_name")
		delete(data.Raw, "allowed_roles")
		delete(data.Raw, "verify_connection")
		delete(data.Raw, "root_rotation_statements")
		delete(data.Raw, "password_policy")

		id, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
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

		// Create a database plugin and initialize it.
		dbw, err := newDatabaseWrapper(ctx, config.PluginName, b.System(), b.logger)
		if err != nil {
			return logical.ErrorResponse("error creating database object: %s", err), nil
		}

		initReq := v5.InitializeRequest{
			Config:           config.ConnectionDetails,
			VerifyConnection: verifyConnection,
		}
		initResp, err := dbw.Initialize(ctx, initReq)
		if err != nil {
			dbw.Close()
			return logical.ErrorResponse("error creating database object: %s", err), nil
		}
		config.ConnectionDetails = initResp.Config

		b.Logger().Debug("created database object", "name", name, "plugin_name", config.PluginName)

		b.Lock()
		defer b.Unlock()

		// Close and remove the old connection
		b.clearConnection(name)

		b.connections[name] = &dbPluginInstance{
			database: dbw,
			name:     name,
			id:       id,
		}

		err = storeConfig(ctx, req.Storage, name, config)
		if err != nil {
			return nil, err
		}

		resp := &logical.Response{}

		// This is a simple test to check for passwords in the connection_url parameter. If one exists,
		// warn the user to use templated url string
		if connURLRaw, ok := config.ConnectionDetails["connection_url"]; ok {
			if connURL, err := url.Parse(connURLRaw.(string)); err == nil {
				if _, ok := connURL.User.Password(); ok {
					resp.AddWarning("Password found in connection_url, use a templated url to enable root rotation and prevent read access to password information.")
				}
			}
		}

		// If using a legacy DB plugin and set the `password_policy` field, send a warning to the user indicating
		// the `password_policy` will not be used
		if dbw.isV4() && config.PasswordPolicy != "" {
			resp.AddWarning(fmt.Sprintf("%s does not support password policies - upgrade to the latest version of "+
				"Vault (or the sdk if using a custom plugin) to gain password policy support", config.PluginName))
		}

		if len(resp.Warnings) == 0 {
			return nil, nil
		}
		return resp, nil
	}
}

func storeConfig(ctx context.Context, storage logical.Storage, name string, config *DatabaseConfig) error {
	entry, err := logical.StorageEntryJSON(fmt.Sprintf("config/%s", name), config)
	if err != nil {
		return fmt.Errorf("unable to marshal object to JSON: %w", err)
	}

	err = storage.Put(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to save object: %w", err)
	}
	return nil
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
