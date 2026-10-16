// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

// Environment variables and flag name constants used across CLI commands.
// Both command/client and command/server import these from here.
const (
	// EnvVaultCLINoColor is an env var that toggles colored UI output.
	EnvVaultCLINoColor = `VAULT_CLI_NO_COLOR`
	// EnvVaultFormat is the output format
	EnvVaultFormat = `VAULT_FORMAT`
	// EnvVaultLicense is an env var used in Vault Enterprise to provide a license blob
	EnvVaultLicense = "VAULT_LICENSE"
	// EnvVaultLicensePath is an env var used in Vault Enterprise to provide a
	// path to a license file on disk
	EnvVaultLicensePath = "VAULT_LICENSE_PATH"
	// EnvVaultDetailed is to output detailed information (e.g., ListResponseWithInfo).
	EnvVaultDetailed = `VAULT_DETAILED`
	// EnvVaultLogFormat is used to specify the log format. Supported values are "standard" and "json"
	EnvVaultLogFormat = "VAULT_LOG_FORMAT"
	// EnvVaultLogLevel is used to specify the log level applied to logging
	// Supported log levels: Trace, Debug, Error, Warn, Info
	EnvVaultLogLevel = "VAULT_LOG_LEVEL"
	// EnvVaultExperiments defines the experiments to enable for a server as a
	// comma separated list. See experiments.ValidExperiments() for the list of
	// valid experiments. Not mutable or persisted in storage, only read and
	// logged at startup _per node_. This was initially introduced for the events
	// system being developed over multiple release cycles.
	EnvVaultExperiments = "VAULT_EXPERIMENTS"
	// EnvVaultPluginTmpdir sets the folder to use for Unix sockets when setting
	// up containerized plugins.
	EnvVaultPluginTmpdir = "VAULT_PLUGIN_TMPDIR"

	// FlagNameAddress is the flag used in the base command to read in the
	// address of the Vault server.
	FlagNameAddress = "address"
	// FlagNameCACert is the flag used in the base command to read in the CA
	// cert.
	FlagNameCACert = "ca-cert"
	// FlagNameCAPath is the flag used in the base command to read in the CA
	// cert path.
	FlagNameCAPath = "ca-path"
	// FlagNameClientKey is the flag used in the base command to read in the
	// client key
	FlagNameClientKey = "client-key"
	// FlagNameClientCert is the flag used in the base command to read in the
	// client cert
	FlagNameClientCert = "client-cert"
	// FlagNameTLSSkipVerify is the flag used in the base command to read in
	// the option to ignore TLS certificate verification.
	FlagNameTLSSkipVerify = "tls-skip-verify"
	// FlagNameTLSServerName is the flag used in the base command to read in
	// the TLS server name.
	FlagNameTLSServerName = "tls-server-name"
	// FlagNameAuditNonHMACRequestKeys is the flag name used for auth/secrets enable
	FlagNameAuditNonHMACRequestKeys = "audit-non-hmac-request-keys"
	// FlagNameAuditNonHMACResponseKeys is the flag name used for auth/secrets enable
	FlagNameAuditNonHMACResponseKeys = "audit-non-hmac-response-keys"
	// FlagNameDescription is the flag name used for tuning the secret and auth mount description parameter
	FlagNameDescription = "description"
	// FlagNameListingVisibility is the flag to toggle whether to show the mount in the UI-specific listing endpoint
	FlagNameListingVisibility = "listing-visibility"
	// FlagNamePassthroughRequestHeaders is the flag name used to set passthrough request headers to the backend
	FlagNamePassthroughRequestHeaders = "passthrough-request-headers"
	// FlagNameAllowedResponseHeaders is used to set allowed response headers from a plugin
	FlagNameAllowedResponseHeaders = "allowed-response-headers"
	// FlagNameTokenType is the flag name used to force a specific token type
	FlagNameTokenType = "token-type"
	// FlagNameAllowedManagedKeys is the flag name used for auth/secrets enable
	FlagNameAllowedManagedKeys = "allowed-managed-keys"
	// FlagNameSealWrap determines if CSPs are seal wrapped
	FlagNameSealWrap = "seal-wrap"
	// FlagNamePluginVersion selects what version of a plugin should be used.
	FlagNamePluginVersion = "plugin-version"
	// FlagNameOverridePinnedVersion is the flag name used for allowing plugin-version to override the pinned plugin version
	FlagNameOverridePinnedVersion = "override-pinned-version"
	// FlagNameIdentityTokenKey selects the key used to sign plugin identity tokens
	FlagNameIdentityTokenKey = "identity-token-key"
	// FlagNameTrimRequestTrailingSlashes selects the key used to determine whether to trim trailing slashes
	FlagNameTrimRequestTrailingSlashes = "trim-request-trailing-slashes"
	// FlagNameUserLockoutThreshold is the flag name used for tuning the auth mount lockout threshold parameter
	FlagNameUserLockoutThreshold = "user-lockout-threshold"
	// FlagNameUserLockoutDuration is the flag name used for tuning the auth mount lockout duration parameter
	FlagNameUserLockoutDuration = "user-lockout-duration"
	// FlagNameUserLockoutCounterResetDuration is the flag name used for tuning the auth mount lockout counter reset parameter
	FlagNameUserLockoutCounterResetDuration = "user-lockout-counter-reset-duration"
	// FlagNameUserLockoutDisable is the flag name used for tuning the auth mount disable lockout parameter
	FlagNameUserLockoutDisable = "user-lockout-disable"
	// FlagNameDisableRedirects is used to prevent the client from honoring a single redirect as a response to a request
	FlagNameDisableRedirects = "disable-redirects"
	// FlagNameCombineLogs is used to specify whether log output should be combined and sent to stdout
	FlagNameCombineLogs = "combine-logs"
	// FlagNameDisableGatedLogs is used to disable gated logs and immediately show the vault logs as they become available.
	FlagNameDisableGatedLogs = "disable-gated-logs"
	// FlagNameLogFile is used to specify the path to the log file that Vault should use for logging
	FlagNameLogFile = "log-file"
	// FlagNameLogRotateBytes is the flag used to specify the number of bytes a log file should be before it is rotated.
	FlagNameLogRotateBytes = "log-rotate-bytes"
	// FlagNameLogRotateDuration is the flag used to specify the duration after which a log file should be rotated.
	FlagNameLogRotateDuration = "log-rotate-duration"
	// FlagNameLogRotateMaxFiles is the flag used to specify the maximum number of older/archived log files to keep.
	FlagNameLogRotateMaxFiles = "log-rotate-max-files"
	// FlagNameLogFormat is the flag used to specify the log format. Supported values are "standard" and "json"
	FlagNameLogFormat = "log-format"
	// FlagNameLogLevel is used to specify the log level applied to logging
	// Supported log levels: Trace, Debug, Error, Warn, Info
	FlagNameLogLevel = "log-level"
	// FlagNameDelegatedAuthAccessors allows operators to specify the allowed mount accessors a backend can delegate
	// authentication
	FlagNameDelegatedAuthAccessors = "delegated-auth-accessors"
)
