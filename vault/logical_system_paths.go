// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"net/http"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *SystemBackend) configPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "config/cors$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "cors",
			},

			Fields: map[string]*framework.FieldSchema{
				"enable": {
					Type:        framework.TypeBool,
					Description: "Enables or disables CORS headers on requests.",
				},
				"allowed_origins": {
					Type:        framework.TypeCommaStringSlice,
					Description: "A comma-separated string or array of strings indicating origins that may make cross-origin requests.",
				},
				"allowed_headers": {
					Type:        framework.TypeCommaStringSlice,
					Description: "A comma-separated string or array of strings indicating headers that are allowed on cross-origin requests.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleCORSRead,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationSuffix: "configuration",
					},
					Summary:     "Return the current CORS settings.",
					Description: "",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"enabled": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"allowed_origins": {
									Type:     framework.TypeCommaStringSlice,
									Required: false,
								},
								"allowed_headers": {
									Type:     framework.TypeCommaStringSlice,
									Required: false,
								},
							},
						}},
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleCORSUpdate,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "configure",
					},
					Summary:     "Configure the CORS settings.",
					Description: "",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleCORSDelete,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "delete",
						OperationSuffix: "configuration",
					},
					Summary: "Remove any CORS settings.",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
			},

			HelpDescription: strings.TrimSpace(sysHelp["config/cors"][0]),
			HelpSynopsis:    strings.TrimSpace(sysHelp["config/cors"][1]),
		},

		{
			Pattern: "config/state/sanitized$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleConfigStateSanitized,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "sanitized-configuration-state",
					},
					Summary:     "Return a sanitized version of the Vault server configuration.",
					Description: "The sanitized output strips configuration values in the storage, HA storage, and seals stanzas, which may contain sensitive values such as API tokens. It also removes any token or secret fields in other stanzas, such as the circonus_api_token from telemetry.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							// response has dynamic keys
							Fields: map[string]*framework.FieldSchema{},
						}},
					},
				},
			},
		},

		{
			Pattern: "config/reload/(?P<subsystem>.+)",
			Fields: map[string]*framework.FieldSchema{
				"subsystem": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["config/reload"][0]),
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleConfigReload,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "reload",
						OperationSuffix: "subsystem",
					},
					Summary:     "Reload the given subsystem",
					Description: "",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
			},
		},

		{
			Pattern: "config/ui/headers/" + framework.GenericNameRegex("header"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "ui-headers",
			},

			Fields: map[string]*framework.FieldSchema{
				"header": {
					Type:        framework.TypeString,
					Description: "The name of the header.",
				},
				"values": {
					Type:        framework.TypeStringSlice,
					Description: "The values to set the header.",
				},
				"multivalue": {
					Type:        framework.TypeBool,
					Description: "Returns multiple values if true",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleConfigUIHeadersRead,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "configuration",
					},
					Summary: "Return the given UI header's configuration",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"value": {
									Type:        framework.TypeString,
									Required:    false,
									Description: "returns the first header value when `multivalue` request parameter is false",
								},
								"values": {
									Type:        framework.TypeCommaStringSlice,
									Required:    false,
									Description: "returns all header values when `multivalue` request parameter is true",
								},
							},
						}},
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleConfigUIHeadersUpdate,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "configure",
					},
					Summary: "Configure the values to be returned for the UI header.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							// returns 200 with null `data`
							Description: "OK",
						}},
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleConfigUIHeadersDelete,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "delete",
						OperationSuffix: "configuration",
					},
					Summary: "Remove a UI header.",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
			},

			HelpDescription: strings.TrimSpace(sysHelp["config/ui/headers"][0]),
			HelpSynopsis:    strings.TrimSpace(sysHelp["config/ui/headers"][1]),
		},

		{
			Pattern: "config/ui/headers/?$",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleConfigUIHeadersList,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationPrefix: "ui-headers",
						OperationVerb:   "list",
					},
					Summary: "Return a list of configured UI headers.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Fields: map[string]*framework.FieldSchema{
								"keys": {
									Type:        framework.TypeCommaStringSlice,
									Description: "Lists of configured UI headers. Omitted if list is empty",
									Required:    false,
								},
							},
						}},
					},
				},
			},

			HelpDescription: strings.TrimSpace(sysHelp["config/ui/headers"][0]),
			HelpSynopsis:    strings.TrimSpace(sysHelp["config/ui/headers"][1]),
		},

		{
			Pattern: "generate-root(/attempt)?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "root-token-generation",
			},

			Fields: map[string]*framework.FieldSchema{
				"pgp_key": {
					Type:        framework.TypeString,
					Description: "Specifies a base64-encoded PGP public key.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "progress2|progress",
					},
					Summary: "Read the configuration and progress of the current root generation attempt.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"nonce": {
									Type:     framework.TypeString,
									Required: true,
								},
								"started": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"progress": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"required": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"complete": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"encoded_token": {
									Type:     framework.TypeString,
									Required: true,
								},
								"encoded_root_token": {
									Type:     framework.TypeString,
									Required: true,
								},
								"pgp_fingerprint": {
									Type:     framework.TypeString,
									Required: true,
								},
								"otp": {
									Type:     framework.TypeString,
									Required: true,
								},
								"otp_length": {
									Type:     framework.TypeInt,
									Required: true,
								},
							},
						}},
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Summary: "Initializes a new root generation attempt.",
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "initialize",
						OperationSuffix: "2|",
					},
					Description: "Only a single root generation attempt can take place at a time. One (and only one) of otp or pgp_key are required.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"nonce": {
									Type:     framework.TypeString,
									Required: true,
								},
								"started": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"progress": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"required": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"complete": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"encoded_token": {
									Type:     framework.TypeString,
									Required: true,
								},
								"encoded_root_token": {
									Type:     framework.TypeString,
									Required: true,
								},
								"pgp_fingerprint": {
									Type:     framework.TypeString,
									Required: true,
								},
								"otp": {
									Type:     framework.TypeString,
									Required: true,
								},
								"otp_length": {
									Type:     framework.TypeInt,
									Required: true,
								},
							},
						}},
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "cancel",
						OperationSuffix: "2|",
					},
					Summary: "Cancels any in-progress root generation attempt.",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["generate-root"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["generate-root"][1]),
		},
		{
			Pattern: "generate-root/update$",
			Fields: map[string]*framework.FieldSchema{
				"key": {
					Type:        framework.TypeString,
					Description: "Specifies a single unseal key share.",
				},
				"nonce": {
					Type:        framework.TypeString,
					Description: "Specifies the nonce of the attempt.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationPrefix: "root-token-generation",
						OperationVerb:   "update",
					},
					Summary:     "Enter a single unseal key share to progress the root generation attempt.",
					Description: "If the threshold number of unseal key shares is reached, Vault will complete the root generation and issue the new token. Otherwise, this API must be called multiple times until that threshold is met. The attempt nonce must be provided with each call.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"nonce": {
									Type:     framework.TypeString,
									Required: true,
								},
								"started": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"progress": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"required": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"complete": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"encoded_token": {
									Type:     framework.TypeString,
									Required: true,
								},
								"encoded_root_token": {
									Type:     framework.TypeString,
									Required: true,
								},
								"pgp_fingerprint": {
									Type:     framework.TypeString,
									Required: true,
								},
								"otp": {
									Type:     framework.TypeString,
									Required: true,
								},
								"otp_length": {
									Type:     framework.TypeInt,
									Required: true,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["generate-root"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["generate-root"][1]),
		},
		{
			Pattern: "decode-token$",
			Fields: map[string]*framework.FieldSchema{
				"encoded_token": {
					Type:        framework.TypeString,
					Description: "Specifies the encoded token (result from generate-root).",
				},
				"otp": {
					Type:        framework.TypeString,
					Description: "Specifies the otp code for decode.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleGenerateRootDecodeTokenUpdate,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "decode",
					},
					Summary: "Decodes the encoded token with the otp.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{Description: "OK"}},
					},
				},
			},
		},

		{
			Pattern: "health$",
			Fields: map[string]*framework.FieldSchema{
				"standbyok": {
					Type:        framework.TypeBool,
					Description: "Specifies if being a standby should still return the active status code.",
				},
				"perfstandbyok": {
					Type:        framework.TypeBool,
					Description: "Specifies if being a performance standby should still return the active status code.",
				},
				"activecode": {
					Type:        framework.TypeInt,
					Description: "Specifies the status code for an active node.",
				},
				"standbycode": {
					Type:        framework.TypeInt,
					Description: "Specifies the status code for a standby node.",
				},
				"drsecondarycode": {
					Type:        framework.TypeInt,
					Description: "Specifies the status code for a DR secondary node.",
				},
				"performancestandbycode": {
					Type:        framework.TypeInt,
					Description: "Specifies the status code for a performance standby node.",
				},
				"sealedcode": {
					Type:        framework.TypeInt,
					Description: "Specifies the status code for a sealed node.",
				},
				"uninitcode": {
					Type:        framework.TypeInt,
					Description: "Specifies the status code for an uninitialized node.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "health-status",
					},
					Summary: "Returns the health status of Vault.",
					Responses: map[int][]framework.Response{
						200: {{Description: "initialized, unsealed, and active"}},
						429: {{Description: "unsealed and standby"}},
						472: {{Description: "data recovery mode replication secondary and active"}},
						501: {{Description: "not initialized"}},
						503: {{Description: "sealed"}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["health"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["health"][1]),
		},

		{
			Pattern: "init$",
			Fields: map[string]*framework.FieldSchema{
				"pgp_keys": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as `secret_shares`.",
				},
				"root_token_pgp_key": {
					Type:        framework.TypeString,
					Description: "Specifies a PGP public key used to encrypt the initial root token. The key must be base64-encoded from its original binary representation.",
				},
				"secret_shares": {
					Type:        framework.TypeInt,
					Description: "Specifies the number of shares to split the unseal key into.",
				},
				"secret_threshold": {
					Type:        framework.TypeInt,
					Description: "Specifies the number of shares required to reconstruct the unseal key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as `secret_shares`.",
				},
				"stored_shares": {
					Type:        framework.TypeInt,
					Description: "Specifies the number of shares that should be encrypted by the HSM and stored for auto-unsealing. Currently must be the same as `secret_shares`.",
				},
				"recovery_shares": {
					Type:        framework.TypeInt,
					Description: "Specifies the number of shares to split the recovery key into.",
				},
				"recovery_threshold": {
					Type:        framework.TypeInt,
					Description: " Specifies the number of shares required to reconstruct the recovery key. This must be less than or equal to `recovery_shares`.",
				},
				"recovery_pgp_keys": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Specifies an array of PGP public keys used to encrypt the output recovery keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as `recovery_shares`.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "initialization-status",
					},
					Summary: "Returns the initialization status of Vault.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "initialize",
					},
					Summary:     "Initialize a new Vault.",
					Description: "The Vault must not have been previously initialized. The recovery options, as well as the stored shares option, are only available when using Vault HSM.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["init"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["init"][1]),
		},
		{
			Pattern: "step-down$",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "step-down",
						OperationSuffix: "leader",
					},
					Summary:     "Cause the node to give up active status.",
					Description: "This endpoint forces the node to give up active status. If the node does not have active status, this endpoint does nothing. Note that the node will sleep for ten seconds before attempting to grab the active lock again, but if no standby nodes grab the active lock in the interim, the same node may become the active node again.",
					Responses: map[int][]framework.Response{
						204: {{Description: "empty body"}},
					},
				},
			},
		},
		{
			Pattern: "loggers$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "loggers",
			},
			Fields: map[string]*framework.FieldSchema{
				"level": {
					Type: framework.TypeString,
					Description: "Log verbosity level. Supported values (in order of detail) are " +
						"\"trace\", \"debug\", \"info\", \"warn\", and \"error\".",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleLoggersRead,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "verbosity-level",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
						}},
					},
					Summary: "Read the log level for all existing loggers.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleLoggersWrite,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "update",
						OperationSuffix: "verbosity-level",
					},
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					Summary: "Modify the log level for all existing loggers.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleLoggersDelete,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "revert",
						OperationSuffix: "verbosity-level",
					},
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					Summary: "Revert the all loggers to use log level provided in config.",
				},
			},
		},
		{
			Pattern: "loggers/" + framework.MatchAllRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "loggers",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "The name of the logger to be modified.",
				},
				"level": {
					Type: framework.TypeString,
					Description: "Log verbosity level. Supported values (in order of detail) are " +
						"\"trace\", \"debug\", \"info\", \"warn\", and \"error\".",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleLoggersByNameRead,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "verbosity-level-for",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
						}},
					},
					Summary: "Read the log level for a single logger.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleLoggersByNameWrite,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "update",
						OperationSuffix: "verbosity-level-for",
					},
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					Summary: "Modify the log level of a single logger.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleLoggersByNameDelete,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "revert",
						OperationSuffix: "verbosity-level-for",
					},
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					Summary: "Revert a single logger to use log level provided in config.",
				},
			},
		},
	}
}

func (b *SystemBackend) rekeyPaths() []*framework.Path {
	respFields := map[string]*framework.FieldSchema{
		"nounce": {
			Type:     framework.TypeString,
			Required: true,
		},
		"started": {
			Type:     framework.TypeString,
			Required: true,
		},
		"t": {
			Type:     framework.TypeInt,
			Required: true,
		},
		"n": {
			Type:     framework.TypeInt,
			Required: true,
		},
		"progress": {
			Type:     framework.TypeInt,
			Required: true,
		},
		"required": {
			Type:     framework.TypeInt,
			Required: true,
		},
		"verification_required": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"verification_nonce": {
			Type:     framework.TypeString,
			Required: true,
		},
		"backup": {
			Type: framework.TypeBool,
		},
		"pgp_fingerprints": {
			Type: framework.TypeCommaStringSlice,
		},
	}

	return []*framework.Path{
		{
			Pattern: "rekey/init",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "rekey-attempt",
			},

			Fields: map[string]*framework.FieldSchema{
				"secret_shares": {
					Type:        framework.TypeInt,
					Description: "Specifies the number of shares to split the unseal key into.",
				},
				"secret_threshold": {
					Type:        framework.TypeInt,
					Description: "Specifies the number of shares required to reconstruct the unseal key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as secret_shares.",
				},
				"pgp_keys": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as secret_shares.",
				},
				"backup": {
					Type:        framework.TypeBool,
					Description: "Specifies if using PGP-encrypted keys, whether Vault should also store a plaintext backup of the PGP-encrypted keys.",
				},
				"require_verification": {
					Type:        framework.TypeBool,
					Description: "Turns on verification functionality",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "progress",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields:      respFields,
						}},
					},
					Summary: "Reads the configuration and progress of the current rekey attempt.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "initialize",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields:      respFields,
						}},
					},
					Summary:     "Initializes a new rekey attempt.",
					Description: "Only a single rekey attempt can take place at a time, and changing the parameters of a rekey requires canceling and starting a new rekey, which will also provide a new nonce.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "cancel",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
						}},
					},
					Summary:     "Cancels any in-progress rekey.",
					Description: "This clears the rekey settings as well as any progress made. This must be called to change the parameters of the rekey. Note: verification is still a part of a rekey. If rekeying is canceled during the verification flow, the current unseal keys remain valid.",
				},
			},
		},
		{
			Pattern: "rekey/backup$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "rekey",
			},

			Fields: map[string]*framework.FieldSchema{},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRekeyRetrieveBarrier,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "backup-key",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"nonce": {
									Type:     framework.TypeString,
									Required: true,
								},
								"keys": {
									Type:     framework.TypeMap,
									Required: true,
								},
								"keys_base64": {
									Type:     framework.TypeMap,
									Required: true,
								},
							},
						}},
					},
					Summary: "Return the backup copy of PGP-encrypted unseal keys.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleRekeyDeleteBarrier,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "delete",
						OperationSuffix: "backup-key",
					},
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					Summary: "Delete the backup copy of PGP-encrypted unseal keys.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["rekey_backup"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["rekey_backup"][0]),
		},

		{
			Pattern: "rekey/recovery-key-backup$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "rekey",
			},

			Fields: map[string]*framework.FieldSchema{},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRekeyRetrieveRecovery,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "backup-recovery-key",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"nonce": {
									Type:     framework.TypeString,
									Required: true,
								},
								"keys": {
									Type:     framework.TypeMap,
									Required: true,
								},
								"keys_base64": {
									Type:     framework.TypeMap,
									Required: true,
								},
							},
						}},
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleRekeyDeleteRecovery,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "delete",
						OperationSuffix: "backup-recovery-key",
					},
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["rekey_backup"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["rekey_backup"][0]),
		},
		{
			Pattern: "rekey/update",

			Fields: map[string]*framework.FieldSchema{
				"key": {
					Type:        framework.TypeString,
					Description: "Specifies a single unseal key share.",
				},
				"nonce": {
					Type:        framework.TypeString,
					Description: "Specifies the nonce of the rekey attempt.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationPrefix: "rekey-attempt",
						OperationVerb:   "update",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"nounce": {
									Type:     framework.TypeString,
									Required: true,
								},
								"complete": {
									Type: framework.TypeBool,
								},
								"started": {
									Type: framework.TypeString,
								},
								"t": {
									Type: framework.TypeInt,
								},
								"n": {
									Type: framework.TypeInt,
								},
								"progress": {
									Type: framework.TypeInt,
								},
								"required": {
									Type: framework.TypeInt,
								},
								"keys": {
									Type: framework.TypeCommaStringSlice,
								},
								"keys_base64": {
									Type: framework.TypeCommaStringSlice,
								},
								"verification_required": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"verification_nonce": {
									Type:     framework.TypeString,
									Required: true,
								},
								"backup": {
									Type: framework.TypeBool,
								},
								"pgp_fingerprints": {
									Type: framework.TypeCommaStringSlice,
								},
							},
						}},
					},
					Summary: "Enter a single unseal key share to progress the rekey of the Vault.",
				},
			},
		},
		{
			Pattern: "rekey/verify",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "rekey-verification",
			},

			Fields: map[string]*framework.FieldSchema{
				"key": {
					Type:        framework.TypeString,
					Description: "Specifies a single unseal share key from the new set of shares.",
				},
				"nonce": {
					Type:        framework.TypeString,
					Description: "Specifies the nonce of the rekey verification operation.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "progress",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"nounce": {
									Type:     framework.TypeString,
									Required: true,
								},
								"started": {
									Type:     framework.TypeString,
									Required: true,
								},
								"t": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"n": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"progress": {
									Type:     framework.TypeInt,
									Required: true,
								},
							},
						}},
					},
					Summary: "Read the configuration and progress of the current rekey verification attempt.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "cancel",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"nounce": {
									Type:     framework.TypeString,
									Required: true,
								},
								"started": {
									Type:     framework.TypeString,
									Required: true,
								},
								"t": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"n": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"progress": {
									Type:     framework.TypeInt,
									Required: true,
								},
							},
						}},
					},
					Summary:     "Cancel any in-progress rekey verification operation.",
					Description: "This clears any progress made and resets the nonce. Unlike a `DELETE` against `sys/rekey/init`, this only resets the current verification operation, not the entire rekey atttempt.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "update",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"nounce": {
									Type:     framework.TypeString,
									Required: true,
								},
								"complete": {
									Type: framework.TypeBool,
								},
							},
						}},
					},
					Summary: "Enter a single new key share to progress the rekey verification operation.",
				},
			},
		},

		{
			Pattern: "seal$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb: "seal",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Summary: "Seal the Vault.",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["seal"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["seal"][1]),
		},

		{
			Pattern: "unseal$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb: "unseal",
			},

			Fields: map[string]*framework.FieldSchema{
				"key": {
					Type:        framework.TypeString,
					Description: "Specifies a single unseal key share. This is required unless reset is true.",
				},
				"reset": {
					Type:        framework.TypeBool,
					Description: "Specifies if previously-provided unseal keys are discarded and the unseal process is reset.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Summary: "Unseal the Vault.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							// unseal returns `vault.SealStatusResponse` struct
							Fields: map[string]*framework.FieldSchema{
								"type": {
									Type:     framework.TypeString,
									Required: true,
								},
								"initialized": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"sealed": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"t": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"n": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"progress": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"nonce": {
									Type:     framework.TypeString,
									Required: true,
								},
								"version": {
									Type:     framework.TypeString,
									Required: true,
								},
								"build_date": {
									Type:     framework.TypeString,
									Required: true,
								},
								"migration": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"cluster_name": {
									Type:     framework.TypeString,
									Required: false,
								},
								"cluster_id": {
									Type:     framework.TypeString,
									Required: false,
								},
								"recovery_seal": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"storage_type": {
									Type:     framework.TypeString,
									Required: false,
								},
								"hcp_link_status": {
									Type:     framework.TypeString,
									Required: false,
								},
								"hcp_link_resource_ID": {
									Type:     framework.TypeString,
									Required: false,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["unseal"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["unseal"][1]),
		},
	}
}

func (b *SystemBackend) statusPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "leader$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "leader",
				OperationVerb:   "status",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleLeaderStatus,
					Summary:  "Returns the high availability status and current leader instance of Vault.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							// returns `vault.LeaderResponse` struct
							Fields: map[string]*framework.FieldSchema{
								"ha_enabled": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"is_self": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"active_time": {
									Type: framework.TypeTime,
									// active_time has 'omitempty' tag, but its not a pointer so never "empty"
									Required: true,
								},
								"leader_address": {
									Type:     framework.TypeString,
									Required: true,
								},
								"leader_cluster_address": {
									Type:     framework.TypeString,
									Required: true,
								},
								"performance_standby": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"performance_standby_last_remote_wal": {
									Type:     framework.TypeInt64,
									Required: true,
								},
								"last_wal": {
									Type:     framework.TypeInt64,
									Required: false,
								},
								"raft_committed_index": {
									Type:     framework.TypeInt64,
									Required: false,
								},
								"raft_applied_index": {
									Type:     framework.TypeInt64,
									Required: false,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis: "Check the high availability status and current leader of Vault",
		},
		{
			Pattern: "seal-status$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "seal",
				OperationVerb:   "status",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleSealStatus,
					Summary:  "Check the seal status of a Vault.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							// unseal returns `vault.SealStatusResponse` struct
							Fields: map[string]*framework.FieldSchema{
								"type": {
									Type:     framework.TypeString,
									Required: true,
								},
								"initialized": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"sealed": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"t": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"n": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"progress": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"nonce": {
									Type:     framework.TypeString,
									Required: true,
								},
								"version": {
									Type:     framework.TypeString,
									Required: true,
								},
								"build_date": {
									Type:     framework.TypeString,
									Required: true,
								},
								"migration": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"cluster_name": {
									Type:     framework.TypeString,
									Required: false,
								},
								"cluster_id": {
									Type:     framework.TypeString,
									Required: false,
								},
								"recovery_seal": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"storage_type": {
									Type:     framework.TypeString,
									Required: false,
								},
								"hcp_link_status": {
									Type:     framework.TypeString,
									Required: false,
								},
								"hcp_link_resource_ID": {
									Type:     framework.TypeString,
									Required: false,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["seal-status"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["seal-status"][1]),
		},
		{
			Pattern: "ha-status$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "ha",
				OperationVerb:   "status",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleHAStatus,
					Summary:  "Check the HA status of a Vault cluster",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"nodes": {
									Type:     framework.TypeSlice,
									Required: true,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["ha-status"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["ha-status"][1]),
		},
		{
			Pattern: "version-history/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb: "version-history",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleVersionHistoryList,
					Summary:  "Returns map of historical version change entries",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"keys": {
									Type:     framework.TypeCommaStringSlice,
									Required: true,
								},
								"key_info": {
									Type:     framework.TypeKVPairs,
									Required: true,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["version-history"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["version-history"][1]),
		},
	}
}

func (b *SystemBackend) auditHashPath() *framework.Path {
	return &framework.Path{
		Pattern: "audit-hash/(?P<path>.+)",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: "auditing",
			OperationVerb:   "calculate",
			OperationSuffix: "hash",
		},

		Fields: map[string]*framework.FieldSchema{
			"path": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["audit_path"][0]),
			},

			"input": {
				Type: framework.TypeString,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.handleAuditHash,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"hash": {
								Type:     framework.TypeString,
								Required: true,
							},
						},
					}},
				},
			},
		},

		HelpSynopsis:    strings.TrimSpace(sysHelp["audit-hash"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["audit-hash"][1]),
	}
}

func (b *SystemBackend) auditPaths() []*framework.Path {
	return []*framework.Path{
		b.auditHashPath(),

		{
			Pattern: "audit$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "auditing",
				OperationVerb:   "list",
				OperationSuffix: "enabled-devices",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleAuditTable,
					Summary:  "List the enabled audit devices.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							// this response has dynamic keys
							Description: "OK",
							Fields:      nil,
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["audit-table"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["audit-table"][1]),
		},

		{
			Pattern: "audit/(?P<path>.+)",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "auditing",
			},

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["audit_path"][0]),
				},
				"type": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["audit_type"][0]),
				},
				"description": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["audit_desc"][0]),
				},
				"options": {
					Type:        framework.TypeKVPairs,
					Description: strings.TrimSpace(sysHelp["audit_opts"][0]),
				},
				"local": {
					Type:        framework.TypeBool,
					Default:     false,
					Description: strings.TrimSpace(sysHelp["mount_local"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleEnableAudit,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "enable",
						OperationSuffix: "device",
					},
					Summary: "Enable a new audit device at the supplied path.",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDisableAudit,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "disable",
						OperationSuffix: "device",
					},
					Summary: "Disable the audit device at the given path.",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["audit"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["audit"][1]),
		},

		{
			Pattern: "config/auditing/request-headers/(?P<header>.+)",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "auditing",
			},

			Fields: map[string]*framework.FieldSchema{
				"header": {
					Type: framework.TypeString,
				},
				"hmac": {
					Type: framework.TypeBool,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleAuditedHeaderUpdate,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "enable",
						OperationSuffix: "request-header",
					},
					Summary: "Enable auditing of a header.",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleAuditedHeaderDelete,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "disable",
						OperationSuffix: "request-header",
					},
					Summary: "Disable auditing of the given request header.",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleAuditedHeaderRead,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "request-header-information",
					},
					Summary: "List the information for the given request header.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							// the response keys are dynamic
							Fields: nil,
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["audited-headers-name"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["audited-headers-name"][1]),
		},

		{
			Pattern: "config/auditing/request-headers$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "auditing",
				OperationVerb:   "list",
				OperationSuffix: "request-headers",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleAuditedHeadersRead,
					Summary:  "List the request headers that are configured to be audited.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"headers": {
									Type:     framework.TypeMap,
									Required: true,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["audited-headers"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["audited-headers"][1]),
		},
	}
}

func (b *SystemBackend) sealPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "key-status$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "encryption-key",
				OperationVerb:   "status",
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.handleKeyStatus,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["key-status"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["key-status"][1]),
		},

		{
			Pattern: "rotate/config$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "encryption-key",
			},

			Fields: map[string]*framework.FieldSchema{
				"enabled": {
					Type:        framework.TypeBool,
					Description: strings.TrimSpace(sysHelp["rotation-enabled"][0]),
				},
				"max_operations": {
					Type:        framework.TypeInt64,
					Description: strings.TrimSpace(sysHelp["rotation-max-operations"][0]),
				},
				"interval": {
					Type:        framework.TypeDurationSecond,
					Description: strings.TrimSpace(sysHelp["rotation-interval"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleKeyRotationConfigRead,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "rotation-configuration",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"max_operations": {
									Type:     framework.TypeInt64,
									Required: true,
								},
								"enabled": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"interval": {
									Type:     framework.TypeDurationSecond,
									Required: true,
								},
							},
						}},
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleKeyRotationConfigUpdate,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "configure",
						OperationSuffix: "rotation",
					},
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					ForwardPerformanceSecondary: true,
					ForwardPerformanceStandby:   true,
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["rotate-config"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["rotate-config"][1]),
		},

		{
			Pattern: "rotate$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "encryption-key",
				OperationVerb:   "rotate",
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleRotate,
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRotate,
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["rotate"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["rotate"][1]),
		},
	}
}

func (b *SystemBackend) pluginsCatalogCRUDPath() *framework.Path {
	return &framework.Path{
		Pattern: "plugins/catalog(/(?P<type>auth|database|secret))?/(?P<name>.+)",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: "plugins-catalog",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_name"][0]),
			},
			"type": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_type"][0]),
			},
			"sha256": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_sha-256"][0]),
			},
			"sha_256": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_sha-256"][0]),
			},
			"oci_image": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_oci-image"][0]),
			},
			"runtime": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_runtime"][0]),
			},
			"command": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_command"][0]),
			},
			"args": {
				Type:        framework.TypeStringSlice,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_args"][0]),
			},
			"env": {
				Type:        framework.TypeStringSlice,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_env"][0]),
			},
			"version": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.handlePluginCatalogUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "register",
					OperationSuffix: "plugin|plugin-with-type|plugin-with-type-and-name",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
					}},
				},
				Summary: "Register a new plugin, or updates an existing one with the supplied name.",
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.handlePluginCatalogDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "remove",
					OperationSuffix: "plugin|plugin-with-type|plugin-with-type-and-name",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields:      map[string]*framework.FieldSchema{},
					}},
				},
				Summary: "Remove the plugin with the given name.",
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handlePluginCatalogRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "read",
					OperationSuffix: "plugin-configuration|plugin-configuration-with-type|plugin-configuration-with-type-and-name",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"name": {
								Type:        framework.TypeString,
								Description: strings.TrimSpace(sysHelp["plugin-catalog_name"][0]),
								Required:    true,
							},
							"sha256": {
								Type:        framework.TypeString,
								Description: strings.TrimSpace(sysHelp["plugin-catalog_sha-256"][0]),
								Required:    true,
							},
							"oci_image": {
								Type:        framework.TypeString,
								Description: strings.TrimSpace(sysHelp["plugin-catalog_oci-image"][0]),
							},
							"runtime": {
								Type:        framework.TypeString,
								Description: strings.TrimSpace(sysHelp["plugin-catalog_runtime"][0]),
							},
							"command": {
								Type:        framework.TypeString,
								Description: strings.TrimSpace(sysHelp["plugin-catalog_command"][0]),
								Required:    true,
							},
							"args": {
								Type:        framework.TypeStringSlice,
								Description: strings.TrimSpace(sysHelp["plugin-catalog_args"][0]),
								Required:    true,
							},
							"version": {
								Type:        framework.TypeString,
								Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
								Required:    true,
							},
							"builtin": {
								Type:     framework.TypeBool,
								Required: true,
							},
							"deprecation_status": {
								Type:     framework.TypeString,
								Required: false,
							},
						},
					}},
				},
				Summary: "Return the configuration data for the plugin with the given name.",
			},
		},

		HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-catalog"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["plugin-catalog"][1]),
	}
}

func (b *SystemBackend) pluginsCatalogListPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "plugins/catalog/(?P<type>auth|database|secret)/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "plugins-catalog",
				OperationVerb:   "list",
				OperationSuffix: "plugins-with-type",
			},

			Fields: map[string]*framework.FieldSchema{
				"type": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["plugin-catalog_type"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handlePluginCatalogTypedList,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"keys": {
									Type:        framework.TypeStringSlice,
									Description: "List of plugin names in the catalog",
									Required:    true,
								},
							},
						}},
					},
					Summary: "List the plugins in the catalog.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-catalog"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["plugin-catalog"][1]),
		},
		{
			Pattern: "plugins/catalog/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "plugins-catalog",
				OperationVerb:   "list",
				OperationSuffix: "plugins",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handlePluginCatalogUntypedList,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"detailed": {
									Type:     framework.TypeMap,
									Required: false,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-catalog-list-all"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["plugin-catalog-list-all"][1]),
		},
	}
}

func (b *SystemBackend) pluginsReloadPath() *framework.Path {
	return &framework.Path{
		Pattern: "plugins/reload/backend$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: "plugins",
			OperationVerb:   "reload",
			OperationSuffix: "backends",
		},

		Fields: map[string]*framework.FieldSchema{
			"plugin": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-backend-reload-plugin"][0]),
			},
			"mounts": {
				Type:        framework.TypeCommaStringSlice,
				Description: strings.TrimSpace(sysHelp["plugin-backend-reload-mounts"][0]),
			},
			"scope": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-backend-reload-scope"][0]),
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.handlePluginReloadUpdate,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"reload_id": {
								Type:     framework.TypeString,
								Required: true,
							},
						},
					}},
					http.StatusAccepted: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"reload_id": {
								Type:     framework.TypeString,
								Required: true,
							},
						},
					}},
				},
				Summary:     "Reload mounted plugin backends.",
				Description: "Either the plugin name (`plugin`) or the desired plugin backend mounts (`mounts`) must be provided, but not both. In the case that the plugin name is provided, all mounted paths that use that plugin backend will be reloaded.  If (`scope`) is provided and is (`global`), the plugin(s) are reloaded globally.",
			},
		},

		HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-reload"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["plugin-reload"][1]),
	}
}

func (b *SystemBackend) pluginsRuntimesCatalogCRUDPath() *framework.Path {
	return &framework.Path{
		Pattern: "plugins/runtimes/catalog/(?P<type>container)/" + framework.GenericNameRegex("name"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: "plugins-runtimes-catalog",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-runtime-catalog_name"][0]),
			},
			"type": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-runtime-catalog_type"][0]),
			},
			"oci_runtime": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-runtime-catalog_oci-runtime"][0]),
			},
			"cgroup_parent": {
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-runtime-catalog_cgroup-parent"][0]),
			},
			"cpu_nanos": {
				Type:        framework.TypeInt64,
				Description: strings.TrimSpace(sysHelp["plugin-runtime-catalog_cpu-nanos"][0]),
			},
			"memory_bytes": {
				Type:        framework.TypeInt64,
				Description: strings.TrimSpace(sysHelp["plugin-runtime-catalog_memory-bytes"][0]),
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.handlePluginRuntimeCatalogUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "register",
					OperationSuffix: "plugin-runtime|plugin-runtime-with-type|plugin-runtime-with-type-and-name", // TODO
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
					}},
				},
				Summary: "Register a new plugin runtime, or updates an existing one with the supplied name.",
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.handlePluginRuntimeCatalogDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "remove",
					OperationSuffix: "plugin-runtime|plugin-runtime-with-type|plugin-runtime-with-type-and-name", // TODO
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
					}},
				},
				Summary: "Remove the plugin runtime with the given name.",
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handlePluginRuntimeCatalogRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "read",
					OperationSuffix: "plugin-runtime-configuration|plugin-runtime-configuration-with-type|plugin-runtime-configuration-with-type-and-name",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"name": {
								Type:        framework.TypeString,
								Description: strings.TrimSpace(sysHelp["plugin-runtime-catalog_name"][0]),
								Required:    true,
							},
							"type": {
								Type:        framework.TypeString,
								Description: strings.TrimSpace(sysHelp["plugin-runtime-catalog_type"][0]),
								Required:    true,
							},
							"oci_runtime": {
								Type:        framework.TypeString,
								Description: strings.TrimSpace(sysHelp["plugin-runtime-catalog_oci-runtime"][0]),
								Required:    true,
							},
							"cgroup_parent": {
								Type:        framework.TypeString,
								Description: strings.TrimSpace(sysHelp["plugin-runtime-catalog_cgroup-parent"][0]),
								Required:    true,
							},
							"cpu_nanos": {
								Type:        framework.TypeInt64,
								Description: strings.TrimSpace(sysHelp["plugin-runtime-catalog_cpu-nanos"][0]),
								Required:    true,
							},
							"memory_bytes": {
								Type:        framework.TypeInt64,
								Description: strings.TrimSpace(sysHelp["plugin-runtime-catalog_memory-bytes"][0]),
								Required:    true,
							},
						},
					}},
				},
				Summary: "Return the configuration data for the plugin runtime with the given name.",
			},
		},

		HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-runtime-catalog"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["plugin-runtime-catalog"][1]),
	}
}

func (b *SystemBackend) pluginsRuntimesCatalogListPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "plugins/runtimes/catalog/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "plugins-runtimes-catalog",
				OperationVerb:   "list",
				OperationSuffix: "plugins-runtimes",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handlePluginRuntimeCatalogList,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"runtimes": {
									Type:        framework.TypeSlice,
									Description: "List of all plugin runtimes in the catalog",
									Required:    true,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-runtime-catalog-list-all"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["plugin-runtime-catalog-list-all"][1]),
		},
	}
}

func (b *SystemBackend) toolsPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "tools/hash" + framework.OptionalParamRegex("urlalgorithm"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb:   "generate",
				OperationSuffix: "hash|hash-with-algorithm",
			},

			Fields: map[string]*framework.FieldSchema{
				"input": {
					Type:        framework.TypeString,
					Description: "The base64-encoded input data",
				},

				"algorithm": {
					Type:    framework.TypeString,
					Default: "sha2-256",
					Description: `Algorithm to use (POST body parameter). Valid values are:

			* sha2-224
			* sha2-256
			* sha2-384
			* sha2-512

			Defaults to "sha2-256".`,
				},

				"urlalgorithm": {
					Type:        framework.TypeString,
					Description: `Algorithm to use (POST URL parameter)`,
				},

				"format": {
					Type:        framework.TypeString,
					Default:     "hex",
					Description: `Encoding format to use. Can be "hex" or "base64". Defaults to "hex".`,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathHashWrite,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"sum": {
									Type:     framework.TypeString,
									Required: true,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["hash"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["hash"][1]),
		},

		{
			Pattern: "tools/random(/" + framework.GenericNameRegex("source") + ")?" + framework.OptionalParamRegex("urlbytes"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb:   "generate",
				OperationSuffix: "random|random-with-source|random-with-bytes|random-with-source-and-bytes",
			},

			Fields: map[string]*framework.FieldSchema{
				"urlbytes": {
					Type:        framework.TypeString,
					Description: "The number of bytes to generate (POST URL parameter)",
				},

				"bytes": {
					Type:        framework.TypeInt,
					Default:     32,
					Description: "The number of bytes to generate (POST body parameter). Defaults to 32 (256 bits).",
				},

				"format": {
					Type:        framework.TypeString,
					Default:     "base64",
					Description: `Encoding format to use. Can be "hex" or "base64". Defaults to "base64".`,
				},

				"source": {
					Type:        framework.TypeString,
					Default:     "platform",
					Description: `Which system to source random data from, ether "platform", "seal", or "all".`,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathRandomWrite,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"random_bytes": {
									Type:     framework.TypeString,
									Required: true,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["random"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["random"][1]),
		},
	}
}

func (b *SystemBackend) internalPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "internal/specs/openapi",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "internal",
				OperationVerb:   "generate",
			},

			Fields: map[string]*framework.FieldSchema{
				"context": {
					Type:        framework.TypeString,
					Description: "Context string appended to every operationId",
					Query:       true,
				},
				"generic_mount_paths": {
					Type:        framework.TypeBool,
					Description: "Use generic mount paths",
					Query:       true,
					Default:     false,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathInternalOpenAPI,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationSuffix: "open-api-document",
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathInternalOpenAPI,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationSuffix: "open-api-document-with-parameters",
					},
				},
			},

			HelpSynopsis: "Generate an OpenAPI 3 document of all mounted paths.",
		},
		{
			Pattern: "internal/ui/feature-flags",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "internal-ui",
				OperationVerb:   "list",
				OperationSuffix: "enabled-feature-flags",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					// callback is absent because this is an unauthenticated method
					Summary: "Lists enabled feature flags.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"feature_flags": {
									Type:     framework.TypeCommaStringSlice,
									Required: true,
								},
							},
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-feature-flags"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-ui-feature-flags"][1]),
		},
		{
			Pattern: "internal/ui/mounts",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "internal-ui",
				OperationVerb:   "list",
				OperationSuffix: "enabled-visible-mounts",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathInternalUIMountsRead,
					Summary:  "Lists all enabled and visible auth and secrets mounts.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"secret": {
									Description: "secret mounts",
									Type:        framework.TypeMap,
									Required:    true,
								},
								"auth": {
									Description: "auth mounts",
									Type:        framework.TypeMap,
									Required:    true,
								},
							},
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-mounts"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-ui-mounts"][1]),
		},
		{
			Pattern: "internal/ui/mounts/(?P<path>.+)",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "internal-ui",
				OperationVerb:   "read",
				OperationSuffix: "mount-information",
			},

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: "The path of the mount.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathInternalUIMountRead,
					Summary:  "Return information about the given mount.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"type": {
									Type:     framework.TypeString,
									Required: true,
								},
								"description": {
									Type:     framework.TypeString,
									Required: true,
								},
								"accessor": {
									Type:     framework.TypeString,
									Required: true,
								},
								"local": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"seal_wrap": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"external_entropy_access": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"options": {
									Type:     framework.TypeMap,
									Required: true,
								},
								"uuid": {
									Type:     framework.TypeString,
									Required: true,
								},
								"plugin_version": {
									Type:     framework.TypeString,
									Required: true,
								},
								"running_plugin_version": {
									Type:     framework.TypeString,
									Required: true,
								},
								"running_sha256": {
									Type:     framework.TypeString,
									Required: true,
								},
								"path": {
									Type:     framework.TypeString,
									Required: true,
								},
								"config": {
									Type:     framework.TypeMap,
									Required: true,
								},
							},
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-mounts"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-ui-mounts"][1]),
		},
		{
			Pattern: "internal/ui/namespaces",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "internal-ui",
				OperationVerb:   "list",
				OperationSuffix: "namespaces",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: pathInternalUINamespacesRead(b),
					Summary:  "Backwards compatibility is not guaranteed for this API",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"keys": {
									Type:        framework.TypeCommaStringSlice,
									Description: "field is only returned if there are one or more namespaces",
								},
							},
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-namespaces"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-ui-namespaces"][1]),
		},
		{
			Pattern: "internal/ui/resultant-acl",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "internal-ui",
				OperationVerb:   "read",
				OperationSuffix: "resultant-acl",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathInternalUIResultantACL,
					Summary:  "Backwards compatibility is not guaranteed for this API",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "empty response returned if no client token",
							Fields:      nil,
						}},
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"root": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"exact_paths": {
									Type:     framework.TypeMap,
									Required: false,
								},
								"glob_paths": {
									Type:     framework.TypeMap,
									Required: false,
								},
							},
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-resultant-acl"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-ui-resultant-acl"][1]),
		},
		{
			Pattern: "internal/counters/requests",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "internal",
				OperationVerb:   "count",
				OperationSuffix: "requests",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:   b.pathInternalCountersRequests,
					Deprecated: true,
					Summary:    "Backwards compatibility is not guaranteed for this API",
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-counters-requests"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-counters-requests"][1]),
		},
		{
			Pattern: "internal/counters/tokens",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "internal",
				OperationVerb:   "count",
				OperationSuffix: "tokens",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathInternalCountersTokens,
					Summary:  "Backwards compatibility is not guaranteed for this API",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"counters": {
									Type:     framework.TypeMap,
									Required: true,
								},
							},
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-counters-tokens"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-counters-tokens"][1]),
		},
		{
			Pattern: "internal/counters/entities",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "internal",
				OperationVerb:   "count",
				OperationSuffix: "entities",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathInternalCountersEntities,
					Summary:  "Backwards compatibility is not guaranteed for this API",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"counters": {
									Type:     framework.TypeMap,
									Required: true,
								},
							},
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-counters-entities"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-counters-entities"][1]),
		},
	}
}

func (b *SystemBackend) introspectionPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "internal/inspect/router/" + framework.GenericNameRegex("tag"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "internal",
				OperationVerb:   "inspect",
				OperationSuffix: "router",
			},
			Fields: map[string]*framework.FieldSchema{
				"tag": {
					Type:        framework.TypeString,
					Description: "Name of subtree being observed",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathInternalInspectRouter,
					Summary:  "Expose the route entry and mount entry tables present in the router",
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-inspect-router"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-inspect-router"][1]),
		},
	}
}

func (b *SystemBackend) capabilitiesPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "capabilities-accessor$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb:   "query",
				OperationSuffix: "token-accessor-capabilities",
			},

			Fields: map[string]*framework.FieldSchema{
				"accessor": {
					Type:        framework.TypeString,
					Description: "Accessor of the token for which capabilities are being queried.",
				},
				"path": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Use 'paths' instead.",
					Deprecated:  true,
				},
				"paths": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Paths on which capabilities are being queried.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleCapabilitiesAccessor,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							// response keys are dynamic
							Fields: nil,
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["capabilities_accessor"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["capabilities_accessor"][1]),
		},

		{
			Pattern: "capabilities$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb:   "query",
				OperationSuffix: "token-capabilities",
			},

			Fields: map[string]*framework.FieldSchema{
				"token": {
					Type:        framework.TypeString,
					Description: "Token for which capabilities are being queried.",
				},
				"path": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Use 'paths' instead.",
					Deprecated:  true,
				},
				"paths": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Paths on which capabilities are being queried.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleCapabilities,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							// response keys are dynamic
							Fields: nil,
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["capabilities"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["capabilities"][1]),
		},

		{
			Pattern: "capabilities-self$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb:   "query",
				OperationSuffix: "token-self-capabilities",
			},

			Fields: map[string]*framework.FieldSchema{
				"token": {
					Type:        framework.TypeString,
					Description: "Token for which capabilities are being queried.",
				},
				"path": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Use 'paths' instead.",
					Deprecated:  true,
				},
				"paths": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Paths on which capabilities are being queried.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleCapabilities,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							// response keys are dynamic
							Fields: nil,
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["capabilities_self"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["capabilities_self"][1]),
		},
	}
}

func (b *SystemBackend) leasePaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "leases/lookup/" + framework.MatchAllRegex("prefix"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "leases",
				OperationVerb:   "look-up",
			},

			Fields: map[string]*framework.FieldSchema{
				"prefix": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["leases-list-prefix"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleLeaseLookupList,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"keys": {
									Type:        framework.TypeCommaStringSlice,
									Description: "A list of lease ids",
									Required:    false,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["leases"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["leases"][1]),
		},

		{
			Pattern: "leases/lookup",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "leases",
				OperationVerb:   "read",
				OperationSuffix: "lease",
			},

			Fields: map[string]*framework.FieldSchema{
				"lease_id": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["lease_id"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleLeaseLookup,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"id": {
									Type:        framework.TypeString,
									Description: "Lease id",
									Required:    true,
								},
								"issue_time": {
									Type:        framework.TypeTime,
									Description: "Timestamp for the lease's issue time",
									Required:    true,
								},
								"renewable": {
									Type:        framework.TypeBool,
									Description: "True if the lease is able to be renewed",
									Required:    true,
								},
								"expire_time": {
									Type:        framework.TypeTime,
									Description: "Optional lease expiry time ",
									Required:    true,
								},
								"last_renewal": {
									Type:        framework.TypeTime,
									Description: "Optional Timestamp of the last time the lease was renewed",
									Required:    true,
								},
								"ttl": {
									Type:        framework.TypeInt,
									Description: "Time to Live set for the lease, returns 0 if unset",
									Required:    true,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["leases"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["leases"][1]),
		},

		{
			Pattern: "(leases/)?renew" + framework.OptionalParamRegex("url_lease_id"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "leases",
				OperationVerb:   "renew",
				OperationSuffix: "lease2|lease|lease-with-id2|lease-with-id",
			},

			Fields: map[string]*framework.FieldSchema{
				"url_lease_id": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["lease_id"][0]),
				},
				"lease_id": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["lease_id"][0]),
				},
				"increment": {
					Type:        framework.TypeDurationSecond,
					Description: strings.TrimSpace(sysHelp["increment"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRenew,
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					Summary: "Renews a lease, requesting to extend the lease.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["renew"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["renew"][1]),
		},

		{
			Pattern: "(leases/)?revoke" + framework.OptionalParamRegex("url_lease_id"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "leases",
				OperationVerb:   "revoke",
				OperationSuffix: "lease2|lease|lease-with-id2|lease-with-id",
			},

			Fields: map[string]*framework.FieldSchema{
				"url_lease_id": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["lease_id"][0]),
				},
				"lease_id": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["lease_id"][0]),
				},
				"sync": {
					Type:        framework.TypeBool,
					Default:     true,
					Description: strings.TrimSpace(sysHelp["revoke-sync"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRevoke,
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					Summary: "Revokes a lease immediately.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["revoke"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["revoke"][1]),
		},

		{
			Pattern: "(leases/)?revoke-force/(?P<prefix>.+)",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "leases",
				OperationVerb:   "force-revoke",
				OperationSuffix: "lease-with-prefix2|lease-with-prefix",
			},

			Fields: map[string]*framework.FieldSchema{
				"prefix": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["revoke-force-path"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRevokeForce,
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					Summary:     "Revokes all secrets or tokens generated under a given prefix immediately",
					Description: "Unlike `/sys/leases/revoke-prefix`, this path ignores backend errors encountered during revocation. This is potentially very dangerous and should only be used in specific emergency situations where errors in the backend or the connected backend service prevent normal revocation.\n\nBy ignoring these errors, Vault abdicates responsibility for ensuring that the issued credentials or secrets are properly revoked and/or cleaned up. Access to this endpoint should be tightly controlled.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["revoke-force"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["revoke-force"][1]),
		},

		{
			Pattern: "(leases/)?revoke-prefix/(?P<prefix>.+)",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "leases",
				OperationVerb:   "revoke",
				OperationSuffix: "lease-with-prefix2|lease-with-prefix",
			},

			Fields: map[string]*framework.FieldSchema{
				"prefix": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["revoke-prefix-path"][0]),
				},
				"sync": {
					Type:        framework.TypeBool,
					Default:     true,
					Description: strings.TrimSpace(sysHelp["revoke-sync"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRevokePrefix,
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					Summary: "Revokes all secrets (via a lease ID prefix) or tokens (via the tokens' path property) generated under a given prefix immediately.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["revoke-prefix"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["revoke-prefix"][1]),
		},

		{
			Pattern: "leases/tidy$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "leases",
				OperationVerb:   "tidy",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleTidyLeases,
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
							Fields:      map[string]*framework.FieldSchema{},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["tidy_leases"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["tidy_leases"][1]),
		},

		{
			Pattern: "leases/count$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "leases",
				OperationVerb:   "count",
			},

			Fields: map[string]*framework.FieldSchema{
				"type": {
					Type:        framework.TypeString,
					Required:    true,
					Description: "Type of leases to get counts for (currently only supporting irrevocable).",
				},
				"include_child_namespaces": {
					Type:        framework.TypeBool,
					Default:     false,
					Description: "Set true if you want counts for this namespace and its children.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					// currently only works for irrevocable leases with param: type=irrevocable
					Callback: b.handleLeaseCount,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"lease_count": {
									Type:        framework.TypeInt,
									Description: "Number of matching leases",
									Required:    true,
								},
								"counts": {
									Type:        framework.TypeInt,
									Description: "Number of matching leases per mount",
									Required:    true,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["count-leases"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["count-leases"][1]),
		},

		{
			Pattern: "leases$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "leases",
				OperationVerb:   "list",
			},

			Fields: map[string]*framework.FieldSchema{
				"type": {
					Type:        framework.TypeString,
					Required:    true,
					Description: "Type of leases to retrieve (currently only supporting irrevocable).",
				},
				"include_child_namespaces": {
					Type:        framework.TypeBool,
					Default:     false,
					Description: "Set true if you want leases for this namespace and its children.",
				},
				"limit": {
					Type:        framework.TypeString,
					Default:     "",
					Description: "Set to a positive integer of the maximum number of entries to return. If you want all results, set to 'none'. If not set, you will get a maximum of 10,000 results returned.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					// currently only works for irrevocable leases with param: type=irrevocable
					Callback: b.handleLeaseList,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"lease_count": {
									Type:        framework.TypeInt,
									Description: "Number of matching leases",
									Required:    true,
								},
								"counts": {
									Type:        framework.TypeInt,
									Description: "Number of matching leases per mount",
									Required:    true,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["list-leases"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["list-leases"][1]),
		},
	}
}

func (b *SystemBackend) remountPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "remount",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb: "remount",
			},

			Fields: map[string]*framework.FieldSchema{
				"from": {
					Type:        framework.TypeString,
					Description: "The previous mount point.",
				},
				"to": {
					Type:        framework.TypeString,
					Description: "The new mount point.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRemount,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"migration_id": {
									Type:     framework.TypeString,
									Required: true,
								},
							},
						}},
					},
					Summary: "Initiate a mount migration",
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["remount"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["remount"][1]),
		},
		{
			Pattern: "remount/status/(?P<migration_id>.+?)$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "remount",
				OperationVerb:   "status",
			},

			Fields: map[string]*framework.FieldSchema{
				"migration_id": {
					Type:        framework.TypeString,
					Description: "The ID of the migration operation",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRemountStatusCheck,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"migration_id": {
									Type:     framework.TypeString,
									Required: true,
								},
								"migration_info": {
									Type:     framework.TypeMap,
									Required: true,
								},
							},
						}},
					},
					Summary: "Check status of a mount migration",
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["remount-status"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["remount-status"][1]),
		},
	}
}

func (b *SystemBackend) metricsPath() *framework.Path {
	return &framework.Path{
		Pattern: "metrics",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationVerb: "metrics",
		},

		Fields: map[string]*framework.FieldSchema{
			"format": {
				Type:        framework.TypeString,
				Description: "Format to export metrics into. Currently accepts only \"prometheus\".",
				Query:       true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handleMetrics,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
					}},
				},
			},
		},
		HelpSynopsis:    strings.TrimSpace(sysHelp["metrics"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["metrics"][1]),
	}
}

func (b *SystemBackend) monitorPath() *framework.Path {
	return &framework.Path{
		Pattern: "monitor",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationVerb: "monitor",
		},

		Fields: map[string]*framework.FieldSchema{
			"log_level": {
				Type:        framework.TypeString,
				Description: "Log level to view system logs at. Currently supported values are \"trace\", \"debug\", \"info\", \"warn\", \"error\".",
				Query:       true,
			},
			"log_format": {
				Type:        framework.TypeString,
				Description: "Output format of logs. Supported values are \"standard\" and \"json\". The default is \"standard\".",
				Query:       true,
				Default:     "standard",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handleMonitor,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
					}},
				},
			},
		},
		HelpSynopsis:    strings.TrimSpace(sysHelp["monitor"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["monitor"][1]),
	}
}

func (b *SystemBackend) inFlightRequestPath() *framework.Path {
	return &framework.Path{
		Pattern: "in-flight-req",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationVerb:   "collect",
			OperationSuffix: "in-flight-request-information",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:    b.handleInFlightRequestData,
				Summary:     strings.TrimSpace(sysHelp["in-flight-req"][0]),
				Description: strings.TrimSpace(sysHelp["in-flight-req"][1]),
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields:      nil, // dynamic fields
					}},
				},
			},
		},
	}
}

func (b *SystemBackend) hostInfoPath() *framework.Path {
	return &framework.Path{
		Pattern: "host-info/?",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationVerb:   "collect",
			OperationSuffix: "host-information",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:    b.handleHostInfo,
				Summary:     strings.TrimSpace(sysHelp["host-info"][0]),
				Description: strings.TrimSpace(sysHelp["host-info"][1]),
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"timestamp": {
								Type:     framework.TypeTime,
								Required: true,
							},
							"cpu": {
								Type:     framework.TypeSlice,
								Required: false,
							},
							"cpu_times": {
								Type:     framework.TypeSlice,
								Required: false,
							},
							"disk": {
								Type:     framework.TypeSlice,
								Required: false,
							},
							"host": {
								Type:     framework.TypeMap,
								Required: false,
							},
							"memory": {
								Type:     framework.TypeMap,
								Required: false,
							},
						},
					}},
				},
			},
		},
		HelpSynopsis:    strings.TrimSpace(sysHelp["host-info"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["host-info"][1]),
	}
}

func (b *SystemBackend) authPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "auth$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "auth",
				OperationVerb:   "list",
				OperationSuffix: "enabled-methods",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleAuthTable,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							// response keys are dynamic
							Fields: nil,
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["auth-table"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["auth-table"][1]),
		},
		{
			Pattern: "auth/(?P<path>.+?)/tune$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "auth",
			},

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_tune"][0]),
				},
				"default_lease_ttl": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["tune_default_lease_ttl"][0]),
				},
				"max_lease_ttl": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["tune_max_lease_ttl"][0]),
				},
				"description": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_desc"][0]),
				},
				"audit_non_hmac_request_keys": {
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["tune_audit_non_hmac_request_keys"][0]),
				},
				"audit_non_hmac_response_keys": {
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["tune_audit_non_hmac_response_keys"][0]),
				},
				"options": {
					Type:        framework.TypeKVPairs,
					Description: strings.TrimSpace(sysHelp["tune_mount_options"][0]),
				},
				"listing_visibility": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["listing_visibility"][0]),
				},
				"passthrough_request_headers": {
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["passthrough_request_headers"][0]),
				},
				"allowed_response_headers": {
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["allowed_response_headers"][0]),
				},
				"token_type": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["token_type"][0]),
				},
				"user_lockout_config": {
					Type:        framework.TypeMap,
					Description: strings.TrimSpace(sysHelp["tune_user_lockout_config"][0]),
				},
				"plugin_version": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleAuthTuneRead,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "tuning-information",
					},
					Summary:     "Reads the given auth path's configuration.",
					Description: "This endpoint requires sudo capability on the final path, but the same functionality can be achieved without sudo via `sys/mounts/auth/[auth-path]/tune`.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"description": {
									Type:     framework.TypeString,
									Required: true,
								},
								"default_lease_ttl": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"max_lease_ttl": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"force_no_cache": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"external_entropy_access": {
									Type:     framework.TypeBool,
									Required: false,
								},
								"token_type": {
									Type:     framework.TypeString,
									Required: false,
								},
								"audit_non_hmac_request_keys": {
									Type:     framework.TypeCommaStringSlice,
									Required: false,
								},
								"audit_non_hmac_response_keys": {
									Type:     framework.TypeCommaStringSlice,
									Required: false,
								},
								"listing_visibility": {
									Type:     framework.TypeString,
									Required: false,
								},
								"passthrough_request_headers": {
									Type:     framework.TypeCommaStringSlice,
									Required: false,
								},
								"allowed_response_headers": {
									Type:     framework.TypeCommaStringSlice,
									Required: false,
								},
								"allowed_managed_keys": {
									Type:     framework.TypeCommaStringSlice,
									Required: false,
								},
								"user_lockout_counter_reset_duration": {
									Type:     framework.TypeInt64,
									Required: false,
								},
								"user_lockout_threshold": {
									Type:     framework.TypeInt64, // uint64
									Required: false,
								},
								"user_lockout_duration": {
									Type:     framework.TypeInt64,
									Required: false,
								},
								"user_lockout_disable": {
									Type:     framework.TypeBool,
									Required: false,
								},
								"options": {
									Type:     framework.TypeMap,
									Required: false,
								},
								"plugin_version": {
									Type:     framework.TypeString,
									Required: false,
								},
							},
						}},
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleAuthTuneWrite,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "tune",
						OperationSuffix: "configuration-parameters",
					},
					Summary:     "Tune configuration parameters for a given auth path.",
					Description: "This endpoint requires sudo capability on the final path, but the same functionality can be achieved without sudo via `sys/mounts/auth/[auth-path]/tune`.",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["auth_tune"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["auth_tune"][1]),
		},
		{
			Pattern: "auth/(?P<path>.+)",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "auth",
			},

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_path"][0]),
				},
				"type": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_type"][0]),
				},
				"description": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_desc"][0]),
				},
				"config": {
					Type:        framework.TypeMap,
					Description: strings.TrimSpace(sysHelp["auth_config"][0]),
				},
				"local": {
					Type:        framework.TypeBool,
					Default:     false,
					Description: strings.TrimSpace(sysHelp["mount_local"][0]),
				},
				"seal_wrap": {
					Type:        framework.TypeBool,
					Default:     false,
					Description: strings.TrimSpace(sysHelp["seal_wrap"][0]),
				},
				"external_entropy_access": {
					Type:        framework.TypeBool,
					Default:     false,
					Description: strings.TrimSpace(sysHelp["external_entropy_access"][0]),
				},
				"plugin_name": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_plugin"][0]),
				},
				"options": {
					Type:        framework.TypeKVPairs,
					Description: strings.TrimSpace(sysHelp["auth_options"][0]),
				},
				"plugin_version": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleReadAuth,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "configuration",
					},
					Summary: "Read the configuration of the auth engine at the given path.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"type": {
									Type:     framework.TypeString,
									Required: true,
								},
								"description": {
									Type:     framework.TypeString,
									Required: true,
								},
								"accessor": {
									Type:     framework.TypeString,
									Required: true,
								},
								"local": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"seal_wrap": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"external_entropy_access": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"options": {
									Type:     framework.TypeMap,
									Required: true,
								},
								"uuid": {
									Type:     framework.TypeString,
									Required: true,
								},
								"plugin_version": {
									Type:     framework.TypeString,
									Required: true,
								},
								"running_plugin_version": {
									Type:     framework.TypeString,
									Required: true,
								},
								"running_sha256": {
									Type:     framework.TypeString,
									Required: true,
								},
								"deprecation_status": {
									Type:     framework.TypeString,
									Required: false,
								},
								"config": {
									Type:     framework.TypeMap,
									Required: true,
								},
							},
						}},
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleEnableAuth,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "enable",
						OperationSuffix: "method",
					},
					Summary: "Enables a new auth method.",
					Description: `After enabling, the auth method can be accessed and configured via the auth path specified as part of the URL. This auth path will be nested under the auth prefix.

For example, enable the "foo" auth method will make it accessible at /auth/foo.`,
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDisableAuth,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "disable",
						OperationSuffix: "method",
					},
					Summary: "Disable the auth method at the given auth path",
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["auth"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["auth"][1]),
		},
	}
}

func (b *SystemBackend) policyPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "policy/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "policies",
				OperationVerb:   "list",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handlePoliciesList(PolicyTypeACL),
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"keys": {
									Type:     framework.TypeStringSlice,
									Required: true,
								},
								"policies": {
									Type: framework.TypeStringSlice,
								},
							},
						}},
					},
					DisplayAttrs: &framework.DisplayAttributes{
						OperationSuffix: "acl-policies2", // this endpoint duplicates sys/policies/acl
					},
				},
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handlePoliciesList(PolicyTypeACL),
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"keys": {
									Type:     framework.TypeStringSlice,
									Required: true,
								},
								"policies": {
									Type: framework.TypeStringSlice,
								},
							},
						}},
					},
					DisplayAttrs: &framework.DisplayAttributes{
						OperationSuffix: "acl-policies3", // this endpoint duplicates sys/policies/acl
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["policy-list"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["policy-list"][1]),
		},

		{
			Pattern: "policy/(?P<name>.+)",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "policies",
				OperationSuffix: "acl-policy2", // this endpoint duplicates /sys/policies/acl
			},

			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["policy-name"][0]),
				},
				"rules": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["policy-rules"][0]),
					Deprecated:  true,
				},
				"policy": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["policy-rules"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handlePoliciesRead(PolicyTypeACL),
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"name": {
									Type:     framework.TypeString,
									Required: true,
								},
								"rules": {
									Type:     framework.TypeString,
									Required: true,
								},
								"policy": {
									Type:     framework.TypeString,
									Required: false,
								},
							},
						}},
					},
					Summary: "Retrieve the policy body for the named policy.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handlePoliciesSet(PolicyTypeACL),
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
							Fields:      map[string]*framework.FieldSchema{},
						}},
					},
					Summary: "Add a new or update an existing policy.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handlePoliciesDelete(PolicyTypeACL),
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
							Fields:      map[string]*framework.FieldSchema{},
						}},
					},
					Summary: "Delete the policy with the given name.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["policy"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["policy"][1]),
		},

		{
			Pattern: "policies/acl/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "policies",
				OperationSuffix: "acl-policies",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handlePoliciesList(PolicyTypeACL),
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"keys": {
									Type:     framework.TypeStringSlice,
									Required: true,
								},
								"policies": {
									Type: framework.TypeStringSlice,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["policy-list"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["policy-list"][1]),
		},

		{
			Pattern: "policies/acl/(?P<name>.+)",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "policies",
				OperationSuffix: "acl-policy",
			},

			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["policy-name"][0]),
				},
				"policy": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["policy-rules"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handlePoliciesRead(PolicyTypeACL),
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"name": {
									Type:     framework.TypeString,
									Required: false,
								},
								"rules": {
									Type:     framework.TypeString,
									Required: false,
								},
								"policy": {
									Type:     framework.TypeString,
									Required: false,
								},
							},
						}},
					},
					Summary: "Retrieve information about the named ACL policy.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handlePoliciesSet(PolicyTypeACL),
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
							Fields:      map[string]*framework.FieldSchema{},
						}},
					},
					Summary: "Add a new or update an existing ACL policy.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handlePoliciesDelete(PolicyTypeACL),
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
							Fields:      map[string]*framework.FieldSchema{},
						}},
					},
					Summary: "Delete the ACL policy with the given name.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["policy"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["policy"][1]),
		},

		{
			Pattern: "policies/password/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "policies",
				OperationSuffix: "password-policies",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handlePoliciesPasswordList,
					Summary:  "List the existing password policies.",
				},
			},
		},

		{
			Pattern: "policies/password/(?P<name>.+)/generate$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "policies",
				OperationVerb:   "generate",
				OperationSuffix: "password-from-password-policy",
			},

			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "The name of the password policy.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handlePoliciesPasswordGenerate,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"password": {
									Type:     framework.TypeString,
									Required: true,
								},
							},
						}},
					},
					Summary: "Generate a password from an existing password policy.",
				},
			},

			HelpSynopsis:    "Generate a password from an existing password policy.",
			HelpDescription: "Generate a password from an existing password policy.",
		},

		{
			Pattern: "policies/password/(?P<name>.+)$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "policies",
				OperationSuffix: "password-policy",
			},

			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "The name of the password policy.",
				},
				"policy": {
					Type:        framework.TypeString,
					Description: "The password policy",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handlePoliciesPasswordSet,
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
							Fields:      map[string]*framework.FieldSchema{},
						}},
					},
					Summary: "Add a new or update an existing password policy.",
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handlePoliciesPasswordGet,
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"policy": {
									Type:     framework.TypeString,
									Required: true,
								},
							},
						}},
					},
					Summary: "Retrieve an existing password policy.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handlePoliciesPasswordDelete,
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
							Fields:      map[string]*framework.FieldSchema{},
						}},
					},
					Summary: "Delete a password policy.",
				},
			},

			HelpSynopsis: "Read, Modify, or Delete a password policy.",
			HelpDescription: "Read the rules of an existing password policy, create or update " +
				"the rules of a password policy, or delete a password policy.",
		},
	}
}

func (b *SystemBackend) wrappingPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "wrapping/wrap$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb: "wrap",
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleWrappingWrap,
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrappingWrap,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							// dynamic fields
							Fields: nil,
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["wrap"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["wrap"][1]),

			TakesArbitraryInput: true,
		},

		{
			Pattern: "wrapping/unwrap$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb: "unwrap",
			},

			Fields: map[string]*framework.FieldSchema{
				"token": {
					Type: framework.TypeString,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrappingUnwrap,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							// dynamic fields
							Fields: nil,
						}},
						http.StatusNoContent: {{
							Description: "No content",
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["unwrap"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["unwrap"][1]),
		},

		{
			Pattern: "wrapping/lookup$",

			Fields: map[string]*framework.FieldSchema{
				"token": {
					Type:  framework.TypeString,
					Query: true,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrappingLookup,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "wrapping-properties",
					},
					Summary: "Look up wrapping properties for the given token.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"creation_ttl": {
									Type:     framework.TypeDurationSecond,
									Required: false,
								},
								"creation_time": {
									Type:     framework.TypeTime,
									Required: false,
								},
								"creation_path": {
									Type:     framework.TypeString,
									Required: false,
								},
							},
						}},
					},
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleWrappingLookup,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "wrapping-properties2",
					},
					Summary: "Look up wrapping properties for the requester's token.",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"creation_ttl": {
									Type:     framework.TypeDurationSecond,
									Required: false,
								},
								"creation_time": {
									Type:     framework.TypeTime,
									Required: false,
								},
								"creation_path": {
									Type:     framework.TypeString,
									Required: false,
								},
							},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["wraplookup"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["wraplookup"][1]),
		},

		{
			Pattern: "wrapping/rewrap$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb: "rewrap",
			},

			Fields: map[string]*framework.FieldSchema{
				"token": {
					Type: framework.TypeString,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrappingRewrap,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							// dynamic fields
							Fields: nil,
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["rewrap"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["rewrap"][1]),
		},
	}
}

func (b *SystemBackend) mountPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "mounts/(?P<path>.+?)/tune$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mounts",
			},

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["mount_path"][0]),
				},
				"default_lease_ttl": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["tune_default_lease_ttl"][0]),
				},
				"max_lease_ttl": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["tune_max_lease_ttl"][0]),
				},
				"description": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_desc"][0]),
				},
				"audit_non_hmac_request_keys": {
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["tune_audit_non_hmac_request_keys"][0]),
				},
				"audit_non_hmac_response_keys": {
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["tune_audit_non_hmac_response_keys"][0]),
				},
				"options": {
					Type:        framework.TypeKVPairs,
					Description: strings.TrimSpace(sysHelp["tune_mount_options"][0]),
				},
				"listing_visibility": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["listing_visibility"][0]),
				},
				"passthrough_request_headers": {
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["passthrough_request_headers"][0]),
				},
				"allowed_response_headers": {
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["allowed_response_headers"][0]),
				},
				"token_type": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["token_type"][0]),
				},
				"allowed_managed_keys": {
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["tune_allowed_managed_keys"][0]),
				},
				"plugin_version": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
				},
				"user_lockout_config": {
					Type:        framework.TypeMap,
					Description: strings.TrimSpace(sysHelp["tune_user_lockout_config"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleMountTuneRead,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "tuning-information",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"max_lease_ttl": {
									Type:        framework.TypeInt,
									Description: strings.TrimSpace(sysHelp["tune_max_lease_ttl"][0]),
									Required:    true,
								},
								"description": {
									Type:        framework.TypeString,
									Description: strings.TrimSpace(sysHelp["auth_desc"][0]),
									Required:    true,
								},
								"default_lease_ttl": {
									Type:        framework.TypeInt,
									Description: strings.TrimSpace(sysHelp["tune_default_lease_ttl"][0]),
									Required:    true,
								},
								"force_no_cache": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"token_type": {
									Type:        framework.TypeString,
									Description: strings.TrimSpace(sysHelp["token_type"][0]),
									Required:    false,
								},
								"allowed_managed_keys": {
									Type:        framework.TypeCommaStringSlice,
									Description: strings.TrimSpace(sysHelp["tune_allowed_managed_keys"][0]),
									Required:    false,
								},
								"allowed_response_headers": {
									Type:        framework.TypeCommaStringSlice,
									Description: strings.TrimSpace(sysHelp["allowed_response_headers"][0]),
									Required:    false,
								},
								"options": {
									Type:        framework.TypeKVPairs,
									Description: strings.TrimSpace(sysHelp["tune_mount_options"][0]),
									Required:    false,
								},
								"plugin_version": {
									Type:        framework.TypeString,
									Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
									Required:    false,
								},
								"external_entropy_access": {
									Type:     framework.TypeBool,
									Required: false,
								},
								"audit_non_hmac_request_keys": {
									Type:     framework.TypeCommaStringSlice,
									Required: false,
								},
								"audit_non_hmac_response_keys": {
									Type:     framework.TypeCommaStringSlice,
									Required: false,
								},
								"listing_visibility": {
									Type:     framework.TypeString,
									Required: false,
								},
								"passthrough_request_headers": {
									Type:     framework.TypeCommaStringSlice,
									Required: false,
								},
								"user_lockout_counter_reset_duration": {
									Type:     framework.TypeInt64,
									Required: false,
								},
								"user_lockout_threshold": {
									Type:     framework.TypeInt64, // TODO this is actuall a Uint64 do we need a new type?
									Required: false,
								},
								"user_lockout_duration": {
									Type:     framework.TypeInt64,
									Required: false,
								},
								"user_lockout_disable": {
									Type:     framework.TypeBool,
									Required: false,
								},
							},
						}},
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleMountTuneWrite,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "tune",
						OperationSuffix: "configuration-parameters",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["mount_tune"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["mount_tune"][1]),
		},

		{
			Pattern: "mounts/(?P<path>.+?)",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mounts",
			},

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["mount_path"][0]),
				},
				"type": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["mount_type"][0]),
				},
				"description": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["mount_desc"][0]),
				},
				"config": {
					Type:        framework.TypeMap,
					Description: strings.TrimSpace(sysHelp["mount_config"][0]),
				},
				"local": {
					Type:        framework.TypeBool,
					Default:     false,
					Description: strings.TrimSpace(sysHelp["mount_local"][0]),
				},
				"seal_wrap": {
					Type:        framework.TypeBool,
					Default:     false,
					Description: strings.TrimSpace(sysHelp["seal_wrap"][0]),
				},
				"external_entropy_access": {
					Type:        framework.TypeBool,
					Default:     false,
					Description: strings.TrimSpace(sysHelp["external_entropy_access"][0]),
				},
				"plugin_name": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["mount_plugin_name"][0]),
				},
				"options": {
					Type:        framework.TypeKVPairs,
					Description: strings.TrimSpace(sysHelp["mount_options"][0]),
				},
				"plugin_version": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleReadMount,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "configuration",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"type": {
									Type:        framework.TypeString,
									Description: strings.TrimSpace(sysHelp["mount_type"][0]),
									Required:    true,
								},
								"description": {
									Type:        framework.TypeString,
									Description: strings.TrimSpace(sysHelp["mount_desc"][0]),
									Required:    true,
								},
								"accessor": {
									Type:     framework.TypeString,
									Required: true,
								},
								"local": {
									Type:        framework.TypeBool,
									Default:     false,
									Description: strings.TrimSpace(sysHelp["mount_local"][0]),
									Required:    true,
								},
								"seal_wrap": {
									Type:        framework.TypeBool,
									Default:     false,
									Description: strings.TrimSpace(sysHelp["seal_wrap"][0]),
									Required:    true,
								},
								"external_entropy_access": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"options": {
									Type:        framework.TypeKVPairs,
									Description: strings.TrimSpace(sysHelp["mount_options"][0]),
									Required:    true,
								},
								"plugin_version": {
									Type:        framework.TypeString,
									Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
									Required:    true,
								},
								"uuid": {
									Type:     framework.TypeString,
									Required: true,
								},
								"running_plugin_version": {
									Type:     framework.TypeString,
									Required: true,
								},
								"running_sha256": {
									Type:     framework.TypeString,
									Required: true,
								},
								"config": {
									Type:        framework.TypeMap,
									Description: strings.TrimSpace(sysHelp["mount_config"][0]),
									Required:    true,
								},
								"deprecation_status": {
									Type:     framework.TypeString,
									Required: false,
								},
							},
						}},
					},
					Summary: "Read the configuration of the secret engine at the given path.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleMount,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "enable",
						OperationSuffix: "secrets-engine",
					},
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					Summary: "Enable a new secrets engine at the given path.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleUnmount,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "disable",
						OperationSuffix: "secrets-engine",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
						}},
					},
					Summary: "Disable the mount point specified at the given path.",
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["mount"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["mount"][1]),
		},

		{
			Pattern: "mounts$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mounts",
				OperationVerb:   "list",
				OperationSuffix: "secrets-engines",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleMountTable,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields:      map[string]*framework.FieldSchema{},
						}},
					},
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["mounts"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["mounts"][1]),
		},
	}
}

func (b *SystemBackend) experimentPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "experiments$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb:   "list",
				OperationSuffix: "experimental-features",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleReadExperiments,
					Summary:  "Returns the available and enabled experiments",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["experiments"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["experiments"][1]),
		},
	}
}

func (b *SystemBackend) lockedUserPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "locked-users/(?P<mount_accessor>.+?)/unlock/(?P<alias_identifier>.+)",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "locked-users",
				OperationVerb:   "unlock",
			},

			Fields: map[string]*framework.FieldSchema{
				"mount_accessor": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["mount_accessor"][0]),
				},
				"alias_identifier": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["alias_identifier"][0]),
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleUnlockUser,
					Summary:  "Unlocks the user with given mount_accessor and alias_identifier",
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["unlock_user"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["unlock_user"][1]),
		},
		{
			Pattern: "locked-users",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "locked-users",
				OperationVerb:   "list",
			},

			Fields: map[string]*framework.FieldSchema{
				"mount_accessor": {
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["mount_accessor"][0]),
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleLockedUsersMetricQuery,
					Summary:  "Report the locked user count metrics, for this namespace and all child namespaces.",
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["locked_users"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["locked_users"][1]),
		},
	}
}
