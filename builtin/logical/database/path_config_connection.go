// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"errors"
	"fmt"

	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/automatedrotationutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	respErrEmptyPluginName = "empty plugin name"
	respErrEmptyName       = "empty name attribute given"
)

// DatabaseConfig is used by the Factory function to configure a Database
// object.
type DatabaseConfig struct {
	EntDatabaseConfig `mapstructure:",squash"`

	PluginName           string `json:"plugin_name" structs:"plugin_name" mapstructure:"plugin_name"`
	PluginVersion        string `json:"plugin_version" structs:"plugin_version" mapstructure:"plugin_version"`
	RunningPluginVersion string `json:"running_plugin_version,omitempty" structs:"running_plugin_version,omitempty" mapstructure:"running_plugin_version,omitempty"`
	// ConnectionDetails stores the database specific connection settings needed
	// by each database type.
	ConnectionDetails map[string]interface{} `json:"connection_details" structs:"connection_details" mapstructure:"connection_details"`
	AllowedRoles      []string               `json:"allowed_roles" structs:"allowed_roles" mapstructure:"allowed_roles"`

	RootCredentialsRotateStatements []string `json:"root_credentials_rotate_statements" structs:"root_credentials_rotate_statements" mapstructure:"root_credentials_rotate_statements"`

	PasswordPolicy   string `json:"password_policy" structs:"password_policy" mapstructure:"password_policy"`
	VerifyConnection bool   `json:"verify_connection" structs:"verify_connection" mapstructure:"verify_connection"`

	// SkipStaticRoleImportRotation is a flag to toggle wether or not a given
	// static account's password should be rotated on creation of the static
	// roles associated with this DB config. This can be overridden at the
	// role-level by the role's skip_import_rotation field. The default is
	// false. Enterprise only.
	SkipStaticRoleImportRotation bool `json:"skip_static_role_import_rotation" structs:"skip_static_role_import_rotation" mapstructure:"skip_static_role_import_rotation"`

	automatedrotationutil.AutomatedRotationParams
}

// ConnectionDetails represents the DatabaseConfig.ConnectionDetails map as a
// struct
type ConnectionDetails struct {
	SelfManaged bool `json:"self_managed" structs:"self_managed" mapstructure:"self_managed"`
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

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixDatabase,
			OperationVerb:   "reset",
			OperationSuffix: "connection",
		},

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

		if err := b.reloadConnection(ctx, req.Storage, name); err != nil {
			return nil, err
		}

		b.dbEvent(ctx, "reset", req.Path, name, false)
		recordDatabaseObservation(ctx, b, req, name, ObservationTypeDatabaseConnectionReset)
		return nil, nil
	}
}

func (b *databaseBackend) reloadConnection(ctx context.Context, storage logical.Storage, name string) error {
	// Close plugin and delete the entry in the connections cache.
	if err := b.ClearConnection(name); err != nil {
		return err
	}

	// Execute plugin again, we don't need the object so throw away.
	if _, err := b.GetConnection(ctx, storage, name); err != nil {
		return err
	}

	return nil
}

// pathReloadPlugin reloads all connections using a named plugin.
func pathReloadPlugin(b *databaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("reload/%s", framework.GenericNameRegex("plugin_name")),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixDatabase,
			OperationVerb:   "reload",
			OperationSuffix: "plugin",
		},

		Fields: map[string]*framework.FieldSchema{
			"plugin_name": {
				Type:        framework.TypeString,
				Description: "Name of the database plugin",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.reloadPlugin(),
		},

		HelpSynopsis:    pathReloadPluginHelpSyn,
		HelpDescription: pathReloadPluginHelpDesc,
	}
}

// reloadPlugin reloads all instances of a named plugin by closing the existing
// instances and creating new ones.
func (b *databaseBackend) reloadPlugin() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		pluginName := data.Get("plugin_name").(string)
		if pluginName == "" {
			return logical.ErrorResponse(respErrEmptyPluginName), nil
		}

		connNames, err := req.Storage.List(ctx, "config/")
		if err != nil {
			return nil, err
		}
		reloaded := []string{}
		reloadFailed := []string{}
		for _, connName := range connNames {
			entry, err := req.Storage.Get(ctx, fmt.Sprintf("config/%s", connName))
			if err != nil {
				return nil, fmt.Errorf("failed to read connection configuration: %w", err)
			}
			if entry == nil {
				continue
			}

			var config DatabaseConfig
			if err := entry.DecodeJSON(&config); err != nil {
				return nil, err
			}
			if config.PluginName == pluginName {
				if err := b.reloadConnection(ctx, req.Storage, connName); err != nil {
					b.Logger().Error("failed to reload connection", "name", connName, "error", err)
					b.dbEvent(ctx, "reload-connection-fail", req.Path, "", false, "name", connName)
					recordDatabaseObservation(ctx, b, req, connName, ObservationTypeDatabaseReloadFail,
						AdditionalDatabaseMetadata{key: "plugin_name", value: pluginName})
					reloadFailed = append(reloadFailed, connName)
				} else {
					b.Logger().Debug("reloaded connection", "name", connName)
					b.dbEvent(ctx, "reload-connection", req.Path, "", true, "name", connName)
					recordDatabaseObservation(ctx, b, req, connName, ObservationTypeDatabaseReloadSuccess,
						AdditionalDatabaseMetadata{key: "plugin_name", value: pluginName})
					reloaded = append(reloaded, connName)
				}
			}
		}

		recordDatabaseObservation(ctx, b, req, "", ObservationTypeDatabaseReloadPlugin,
			AdditionalDatabaseMetadata{key: "plugin_name", value: pluginName},
			AdditionalDatabaseMetadata{key: "reloaded", value: reloaded},
			AdditionalDatabaseMetadata{key: "reload_failed", value: reloadFailed})

		resp := &logical.Response{
			Data: map[string]interface{}{
				"connections": reloaded,
				"count":       len(reloaded),
			},
		}

		if len(reloaded) > 0 {
			b.dbEvent(ctx, "reload", req.Path, "", true, "plugin_name", pluginName)
		} else if len(reloaded) == 0 && len(reloadFailed) == 0 {
			b.Logger().Debug("no connections were found", "plugin_name", pluginName)
		}

		return resp, nil
	}
}

// pathConfigurePluginConnection returns a configured framework.Path setup to
// operate on plugins.
func pathConfigurePluginConnection(b *databaseBackend) *framework.Path {
	fields := map[string]*framework.FieldSchema{
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

		"plugin_version": {
			Type:        framework.TypeString,
			Description: `The version of the plugin to use.`,
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
	}
	AddConnectionFieldsEnt(fields)
	automatedrotationutil.AddAutomatedRotationFields(fields)

	return &framework.Path{
		Pattern: fmt.Sprintf("config/%s", framework.GenericNameRegex("name")),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixDatabase,
		},

		Fields: fields,

		ExistenceCheck: b.connectionExistenceCheck(),

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.connectionWriteHandler(),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "connection",
				},
				ForwardPerformanceSecondary: true,
				ForwardPerformanceStandby:   true,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.connectionWriteHandler(),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "connection",
				},
				ForwardPerformanceSecondary: true,
				ForwardPerformanceStandby:   true,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.connectionReadHandler(),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "read",
					OperationSuffix: "connection-configuration",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.connectionDeleteHandler(),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "delete",
					OperationSuffix: "connection-configuration",
				},
			},
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

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixDatabase,
			OperationSuffix: "connections",
		},

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

		b.dbEvent(ctx, "config-delete", req.Path, name, true)
		recordDatabaseObservation(ctx, b, req, name, ObservationTypeDatabaseConfigDelete)
		return nil, nil
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

func (b *databaseBackend) getPinnedVersion(ctx context.Context, pluginName string) (string, error) {
	extendedSys, ok := b.System().(logical.ExtendedSystemView)
	if !ok {
		return "", fmt.Errorf("database backend does not support running as an external plugin")
	}

	pin, err := extendedSys.GetPinnedPluginVersion(ctx, consts.PluginTypeDatabase, pluginName)
	if errors.Is(err, pluginutil.ErrPinnedVersionNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return pin.Version, nil
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

const pathReloadPluginHelpSyn = `
Reloads all connections using a named database plugin.
`

const pathReloadPluginHelpDesc = `
This path resets each database connection using a named plugin by closing each
existing database plugin instance and running a new one.
`
