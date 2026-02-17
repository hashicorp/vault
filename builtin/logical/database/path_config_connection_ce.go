// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package database

import (
	"context"
	"fmt"
	"net/url"
	"sort"

	"github.com/fatih/structs"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/helper/versions"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/rotation"
)

// EntDatabaseConfig is an empty struct for community edition
type EntDatabaseConfig struct{}

// AddConnectionFieldsEnt is a no-op for community edition
func AddConnectionFieldsEnt(fields map[string]*framework.FieldSchema) {
	// no-op
}

// connectionWriteHandler returns a handler function for creating and updating
// both builtin and plugin database types for community edition
func (b *databaseBackend) connectionWriteHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse(respErrEmptyName), nil
		}

		// Baseline
		config := &DatabaseConfig{
			VerifyConnection: true,
		}

		entry, err := req.Storage.Get(ctx, fmt.Sprintf("config/%s", name))
		if err != nil {
			return nil, fmt.Errorf("failed to read connection configuration: %w", err)
		}
		if entry != nil {
			if err := entry.DecodeJSON(config); err != nil {
				return nil, err
			}
		}

		// If this value was provided as part of the request we want to set it to this value
		if verifyConnectionRaw, ok := data.GetOk("verify_connection"); ok {
			config.VerifyConnection = verifyConnectionRaw.(bool)
		} else if req.Operation == logical.CreateOperation {
			config.VerifyConnection = data.Get("verify_connection").(bool)
		}

		if pluginNameRaw, ok := data.GetOk("plugin_name"); ok {
			config.PluginName = pluginNameRaw.(string)
		} else if req.Operation == logical.CreateOperation {
			config.PluginName = data.Get("plugin_name").(string)
		}
		if config.PluginName == "" {
			return logical.ErrorResponse(respErrEmptyPluginName), nil
		}

		pluginVersion, respErr, err := b.selectPluginVersion(ctx, config, data, req.Operation)
		if respErr != nil || err != nil {
			return respErr, err
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

		if skipImportRotationRaw, ok := data.GetOk("skip_static_role_import_rotation"); ok {
			config.SkipStaticRoleImportRotation = skipImportRotationRaw.(bool)
		}

		if err := config.ParseAutomatedRotationFields(data); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		// Remove these entries from the data before we store it keyed under
		// ConnectionDetails.
		delete(data.Raw, "name")
		delete(data.Raw, "plugin_name")
		delete(data.Raw, "plugin_version")
		delete(data.Raw, "allowed_roles")
		delete(data.Raw, "verify_connection")
		delete(data.Raw, "root_rotation_statements")
		delete(data.Raw, "password_policy")
		delete(data.Raw, "skip_static_role_import_rotation")
		delete(data.Raw, "rotation_schedule")
		delete(data.Raw, "rotation_window")
		delete(data.Raw, "rotation_period")
		delete(data.Raw, "disable_automated_rotation")
		delete(data.Raw, "EntDatabaseConfig")

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
		dbw, err := newDatabaseWrapper(ctx, config.PluginName, pluginVersion, b.System(), b.logger)
		if err != nil {
			return logical.ErrorResponse("error creating database object: %s", err), nil
		}

		initReq := v5.InitializeRequest{
			Config:           config.ConnectionDetails,
			VerifyConnection: config.VerifyConnection,
		}
		initResp, err := dbw.Initialize(ctx, initReq)
		if err != nil {
			dbw.Close()
			return logical.ErrorResponse("error creating database object: %s", err), nil
		}
		config.ConnectionDetails = initResp.Config

		b.Logger().Debug("created database object", "name", name, "plugin_name", config.PluginName)

		// Close and remove the old connection
		oldConn := b.connections.Put(name, &dbPluginInstance{
			database:             dbw,
			name:                 name,
			id:                   id,
			runningPluginVersion: pluginVersion,
		})
		if oldConn != nil {
			oldConn.Close()
		}

		var performedRotationManagerOpern string
		if config.ShouldDeregisterRotationJob() {
			performedRotationManagerOpern = rotation.PerformedDeregistration
			// Disable Automated Rotation and Deregister credentials if required
			deregisterReq := &rotation.RotationJobDeregisterRequest{
				MountPoint: req.MountPoint,
				ReqPath:    req.Path,
			}

			b.Logger().Debug("Deregistering rotation job", "mount", req.MountPoint+req.Path)
			if err := b.System().DeregisterRotationJob(ctx, deregisterReq); err != nil {
				return logical.ErrorResponse("error deregistering rotation job: %s", err), nil
			}
		} else if config.ShouldRegisterRotationJob() {
			performedRotationManagerOpern = rotation.PerformedRegistration
			// Register the rotation job if it's required.
			cfgReq := &rotation.RotationJobConfigureRequest{
				MountPoint:       req.MountPoint,
				ReqPath:          req.Path,
				RotationSchedule: config.RotationSchedule,
				RotationWindow:   config.RotationWindow,
				RotationPeriod:   config.RotationPeriod,
			}

			b.Logger().Debug("Registering rotation job", "mount", req.MountPoint+req.Path)
			if _, err = b.System().RegisterRotationJob(ctx, cfgReq); err != nil {
				return logical.ErrorResponse("error registering rotation job: %s", err), nil
			}
		}

		// 1.12.0 and 1.12.1 stored builtin plugins in storage, but 1.12.2 reverted
		// that, so clean up any pre-existing stored builtin versions on write.
		if versions.IsBuiltinVersion(config.PluginVersion) {
			config.PluginVersion = ""
		}
		err = storeConfig(ctx, req.Storage, name, config)
		if err != nil {
			wrappedError := err
			if performedRotationManagerOpern != "" {
				b.Logger().Error("write to storage failed but the rotation manager still succeeded.",
					"operation", performedRotationManagerOpern, "mount", req.MountPoint, "path", req.Path)
				wrappedError = fmt.Errorf("write to storage failed but the rotation manager still succeeded; "+
					"operation=%s, mount=%s, path=%s, storageError=%s", performedRotationManagerOpern, req.MountPoint, req.Path, err)
			}
			return nil, wrappedError
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

		// We can ignore the error at this point since we're simply adding a warning.
		dbType, _ := dbw.Type()
		if dbType == "snowflake" && config.ConnectionDetails["password"] != nil {
			resp.AddWarning(`[DEPRECATED] Single-factor password authentication is deprecated in Snowflake and will
be removed by November 2025. Key pair authentication will be required after this date. Please
see the Vault documentation for details on the removal of this feature. More information is
available at https://www.snowflake.com/en/blog/blocking-single-factor-password-authentification`)
		}

		var rotationPeriodString string
		if config.RotationPeriod != 0 {
			rotationPeriodString = config.RotationPeriod.String()
		}
		b.dbEvent(ctx, "config-write", req.Path, name, true)
		recordDatabaseObservation(ctx, b, req, name, ObservationTypeDatabaseConfigWrite,
			AdditionalDatabaseMetadata{key: "root_rotation_period", value: rotationPeriodString},
			AdditionalDatabaseMetadata{key: "root_rotation_schedule", value: config.RotationSchedule})

		if len(resp.Warnings) == 0 {
			return nil, nil
		}
		return resp, nil
	}
}

// selectPluginVersion returns the appropriate plugin version for community edition
func (b *databaseBackend) selectPluginVersion(ctx context.Context, config *DatabaseConfig, data *framework.FieldData, op logical.Operation) (string, *logical.Response, error) {
	pinnedVersion, err := b.getPinnedVersion(ctx, config.PluginName)
	if err != nil {
		return "", nil, err
	}
	pluginVersionRaw, ok := data.GetOk("plugin_version")

	switch {
	case ok && pinnedVersion != "":
		return "", logical.ErrorResponse("cannot specify plugin_version for plugin %q as it is pinned (%s)", config.PluginName, pinnedVersion), nil
	case pinnedVersion != "":
		return pinnedVersion, nil, nil
	case ok:
		config.PluginVersion = pluginVersionRaw.(string)
	}

	var builtinShadowed bool
	if unversionedPlugin, err := b.System().LookupPlugin(ctx, config.PluginName, consts.PluginTypeDatabase); err == nil && !unversionedPlugin.Builtin {
		builtinShadowed = true
	}
	switch {
	case config.PluginVersion != "":
		semanticVersion, err := version.NewVersion(config.PluginVersion)
		if err != nil {
			return "", logical.ErrorResponse("version %q is not a valid semantic version: %s", config.PluginVersion, err), nil
		}

		// Canonicalize the version.
		config.PluginVersion = "v" + semanticVersion.String()

		if config.PluginVersion == versions.GetBuiltinVersion(consts.PluginTypeDatabase, config.PluginName) {
			if builtinShadowed {
				return "", logical.ErrorResponse("database plugin %q, version %s not found, as it is"+
					" overridden by an unversioned plugin of the same name. Omit `plugin_version` to use the unversioned plugin", config.PluginName, config.PluginVersion), nil
			}

			config.PluginVersion = ""
		}
	case builtinShadowed:
		// We'll select the unversioned plugin that's been registered.
	case op == logical.CreateOperation:
		// No version provided and no unversioned plugin of that name available.
		// Pin to the current latest version if any versioned plugins are registered.
		plugins, err := b.System().ListVersionedPlugins(ctx, consts.PluginTypeDatabase)
		if err != nil {
			return "", nil, err
		}

		var versionedCandidates []pluginutil.VersionedPlugin
		for _, plugin := range plugins {
			if !plugin.Builtin && plugin.Name == config.PluginName && plugin.Version != "" {
				versionedCandidates = append(versionedCandidates, plugin)
			}
		}

		if len(versionedCandidates) != 0 {
			// Sort in reverse order.
			sort.SliceStable(versionedCandidates, func(i, j int) bool {
				return versionedCandidates[i].SemanticVersion.GreaterThan(versionedCandidates[j].SemanticVersion)
			})

			config.PluginVersion = "v" + versionedCandidates[0].SemanticVersion.String()
			b.logger.Debug(fmt.Sprintf("pinning %q database plugin version %q from candidates %v", config.PluginName, config.PluginVersion, versionedCandidates))
		}
	}

	return config.PluginVersion, nil, nil
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

		if versions.IsBuiltinVersion(config.PluginVersion) {
			// This gets treated as though it's empty when mounting, and will get
			// overwritten to be empty when the config is next written. See #18051.
			config.PluginVersion = ""
		}

		delete(config.ConnectionDetails, "password")
		delete(config.ConnectionDetails, "private_key")
		delete(config.ConnectionDetails, "service_account_json")

		resp := &logical.Response{}
		if dbi, err := b.GetConnectionSkipVerify(ctx, req.Storage, name); err == nil {
			config.RunningPluginVersion = dbi.runningPluginVersion
			if config.PluginVersion != "" && config.PluginVersion != config.RunningPluginVersion {
				warning := fmt.Sprintf("Plugin version is configured as %q, but running %q", config.PluginVersion, config.RunningPluginVersion)
				if pinnedVersion, _ := b.getPinnedVersion(ctx, config.PluginName); pinnedVersion == config.RunningPluginVersion {
					warning += " because that version is pinned"
				} else {
					warning += " either due to a pinned version or because the plugin was upgraded and not yet reloaded"
				}
				resp.AddWarning(warning)
			}
		}

		resp.Data = structs.New(config).Map()
		config.PopulateAutomatedRotationData(resp.Data)
		// remove extra nested AutomatedRotationParams key
		// before returning response
		delete(resp.Data, "AutomatedRotationParams")

		// remove nested EntDatabaseConfig key before returning response
		delete(resp.Data, "EntDatabaseConfig")

		recordDatabaseObservation(ctx, b, req, name, ObservationTypeDatabaseConfigRead)

		return resp, nil
	}
}
