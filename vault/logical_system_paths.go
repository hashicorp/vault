package vault

import (
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *SystemBackend) configPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "config/cors$",

			Fields: map[string]*framework.FieldSchema{
				"enable": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "Enables or disables CORS headers on requests.",
				},
				"allowed_origins": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "A comma-separated string or array of strings indicating origins that may make cross-origin requests.",
				},
				"allowed_headers": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "A comma-separated string or array of strings indicating headers that are allowed on cross-origin requests.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    b.handleCORSRead,
					Summary:     "Return the current CORS settings.",
					Description: "",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback:    b.handleCORSUpdate,
					Summary:     "Configure the CORS settings.",
					Description: "",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleCORSDelete,
					Summary:  "Remove any CORS settings.",
				},
			},

			HelpDescription: strings.TrimSpace(sysHelp["config/cors"][0]),
			HelpSynopsis:    strings.TrimSpace(sysHelp["config/cors"][1]),
		},

		{
			Pattern: "config/ui/headers/" + framework.GenericNameRegex("header"),

			Fields: map[string]*framework.FieldSchema{
				"header": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The name of the header.",
				},
				"values": &framework.FieldSchema{
					Type:        framework.TypeStringSlice,
					Description: "The values to set the header.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleConfigUIHeadersRead,
					Summary:  "Return the given UI header's configuration",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleConfigUIHeadersUpdate,
					Summary:  "Configure the values to be returned for the UI header.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleConfigUIHeadersDelete,
					Summary:  "Remove a UI header.",
				},
			},

			HelpDescription: strings.TrimSpace(sysHelp["config/ui/headers"][0]),
			HelpSynopsis:    strings.TrimSpace(sysHelp["config/ui/headers"][1]),
		},

		{
			Pattern: "config/ui/headers/$",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleConfigUIHeadersList,
					Summary:  "Return a list of configured UI headers.",
				},
			},

			HelpDescription: strings.TrimSpace(sysHelp["config/ui/headers"][0]),
			HelpSynopsis:    strings.TrimSpace(sysHelp["config/ui/headers"][1]),
		},

		{
			Pattern: "generate-root(/attempt)?$",
			Fields: map[string]*framework.FieldSchema{
				"pgp_key": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Specifies a base64-encoded PGP public key.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Summary: "Read the configuration and progress of the current root generation attempt.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Summary:     "Initializes a new root generation attempt.",
					Description: "Only a single root generation attempt can take place at a time. One (and only one) of otp or pgp_key are required.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Summary: "Cancels any in-progress root generation attempt.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["generate-root"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["generate-root"][1]),
		},
		{
			Pattern: "generate-root/update$",
			Fields: map[string]*framework.FieldSchema{
				"key": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Specifies a single master key share.",
				},
				"nonce": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Specifies the nonce of the attempt.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Summary:     "Enter a single master key share to progress the root generation attempt.",
					Description: "If the threshold number of master key shares is reached, Vault will complete the root generation and issue the new token. Otherwise, this API must be called multiple times until that threshold is met. The attempt nonce must be provided with each call.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["generate-root"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["generate-root"][1]),
		},
		{
			Pattern: "health$",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
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
		},

		{
			Pattern: "init$",
			Fields: map[string]*framework.FieldSchema{
				"pgp_keys": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as `secret_shares`.",
				},
				"root_token_pgp_key": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Specifies a PGP public key used to encrypt the initial root token. The key must be base64-encoded from its original binary representation.",
				},
				"secret_shares": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Specifies the number of shares to split the master key into.",
				},
				"secret_threshold": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Specifies the number of shares required to reconstruct the master key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as `secret_shares`.",
				},
				"stored_shares": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Specifies the number of shares that should be encrypted by the HSM and stored for auto-unsealing. Currently must be the same as `secret_shares`.",
				},
				"recovery_shares": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Specifies the number of shares to split the recovery key into.",
				},
				"recovery_threshold": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: " Specifies the number of shares required to reconstruct the recovery key. This must be less than or equal to `recovery_shares`.",
				},
				"recovery_pgp_keys": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "Specifies an array of PGP public keys used to encrypt the output recovery keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as `recovery_shares`.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Summary: "Returns the initialization status of Vault.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Summary:     "Initialize a new Vault.",
					Description: "The Vault must not have been previously initialized. The recovery options, as well as the stored shares option, are only available when using Vault HSM.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["init"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["init"][1]),
		},
		{
			Pattern: "leader$",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Summary: "Returns the high availability status and current leader instance of Vault.",
				},
			},

			HelpSynopsis: "Check the high availability status and current leader of Vault",
		},
		{
			Pattern: "step-down$",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Summary:     "Cause the node to give up active status.",
					Description: "This endpoint forces the node to give up active status. If the node does not have active status, this endpoint does nothing. Note that the node will sleep for ten seconds before attempting to grab the active lock again, but if no standby nodes grab the active lock in the interim, the same node may become the active node again.",
					Responses: map[int][]framework.Response{
						204: {{Description: "empty body"}},
					},
				},
			},
		},
	}
}

func (b *SystemBackend) rekeyPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "rekey/init",

			Fields: map[string]*framework.FieldSchema{
				"secret_shares": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Specifies the number of shares to split the master key into.",
				},
				"secret_threshold": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Specifies the number of shares required to reconstruct the master key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as secret_shares.",
				},
				"pgp_keys": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as secret_shares.",
				},
				"backup": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "Specifies if using PGP-encrypted keys, whether Vault should also store a plaintext backup of the PGP-encrypted keys.",
				},
				"require_verification": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "Turns on verification functionality",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Summary: "Reads the configuration and progress of the current rekey attempt.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Summary:     "Initializes a new rekey attempt.",
					Description: "Only a single rekey attempt can take place at a time, and changing the parameters of a rekey requires canceling and starting a new rekey, which will also provide a new nonce.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Summary:     "Cancels any in-progress rekey.",
					Description: "This clears the rekey settings as well as any progress made. This must be called to change the parameters of the rekey. Note: verification is still a part of a rekey. If rekeying is canceled during the verification flow, the current unseal keys remain valid.",
				},
			},
		},
		{
			Pattern: "rekey/backup$",

			Fields: map[string]*framework.FieldSchema{},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRekeyRetrieveBarrier,
					Summary:  "Return the backup copy of PGP-encrypted unseal keys.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleRekeyDeleteBarrier,
					Summary:  "Delete the backup copy of PGP-encrypted unseal keys.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["rekey_backup"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["rekey_backup"][0]),
		},

		{
			Pattern: "rekey/recovery-key-backup$",

			Fields: map[string]*framework.FieldSchema{},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.handleRekeyRetrieveRecovery,
				logical.DeleteOperation: b.handleRekeyDeleteRecovery,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["rekey_backup"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["rekey_backup"][0]),
		},
		{
			Pattern: "rekey/update",

			Fields: map[string]*framework.FieldSchema{
				"key": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Specifies a single master key share.",
				},
				"nonce": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Specifies the nonce of the rekey attempt.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Summary: "Enter a single master key share to progress the rekey of the Vault.",
				},
			},
		},
		{
			Pattern: "rekey/verify",

			Fields: map[string]*framework.FieldSchema{
				"key": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Specifies a single master share key from the new set of shares.",
				},
				"nonce": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Specifies the nonce of the rekey verification operation.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Summary: "Read the configuration and progress of the current rekey verification attempt.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Summary:     "Cancel any in-progress rekey verification operation.",
					Description: "This clears any progress made and resets the nonce. Unlike a `DELETE` against `sys/rekey/init`, this only resets the current verification operation, not the entire rekey atttempt.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Summary: "Enter a single new key share to progress the rekey verification operation.",
				},
			},
		},

		{
			Pattern: "seal-status$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Summary: "Check the seal status of a Vault.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["seal-status"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["seal-status"][1]),
		},

		{
			Pattern: "seal$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Summary: "Seal the Vault.",
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["seal"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["seal"][1]),
		},

		{
			Pattern: "unseal$",
			Fields: map[string]*framework.FieldSchema{
				"key": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Specifies a single master key share. This is required unless reset is true.",
				},
				"reset": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "Specifies if previously-provided unseal keys are discarded and the unseal process is reset.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Summary: "Unseal the Vault.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["unseal"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["unseal"][1]),
		},
	}
}

func (b *SystemBackend) auditPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "audit-hash/(?P<path>.+)",

			Fields: map[string]*framework.FieldSchema{
				"path": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["audit_path"][0]),
				},

				"input": &framework.FieldSchema{
					Type: framework.TypeString,
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleAuditHash,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["audit-hash"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["audit-hash"][1]),
		},

		{
			Pattern: "audit$",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleAuditTable,
					Summary:  "List the enabled audit devices.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["audit-table"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["audit-table"][1]),
		},

		{
			Pattern: "audit/(?P<path>.+)",

			Fields: map[string]*framework.FieldSchema{
				"path": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["audit_path"][0]),
				},
				"type": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["audit_type"][0]),
				},
				"description": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["audit_desc"][0]),
				},
				"options": &framework.FieldSchema{
					Type:        framework.TypeKVPairs,
					Description: strings.TrimSpace(sysHelp["audit_opts"][0]),
				},
				"local": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Default:     false,
					Description: strings.TrimSpace(sysHelp["mount_local"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleEnableAudit,
					Summary:  "Enable a new audit device at the supplied path.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDisableAudit,
					Summary:  "Disable the audit device at the given path.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["audit"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["audit"][1]),
		},

		{
			Pattern: "config/auditing/request-headers/(?P<header>.+)",

			Fields: map[string]*framework.FieldSchema{
				"header": &framework.FieldSchema{
					Type: framework.TypeString,
				},
				"hmac": &framework.FieldSchema{
					Type: framework.TypeBool,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleAuditedHeaderUpdate,
					Summary:  "Enable auditing of a header.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleAuditedHeaderDelete,
					Summary:  "Disable auditing of the given request header.",
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleAuditedHeaderRead,
					Summary:  "List the information for the given request header.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["audited-headers-name"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["audited-headers-name"][1]),
		},

		{
			Pattern: "config/auditing/request-headers$",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleAuditedHeadersRead,
					Summary:  "List the request headers that are configured to be audited.",
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.handleKeyStatus,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["key-status"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["key-status"][1]),
		},

		{
			Pattern: "rotate$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleRotate,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["rotate"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["rotate"][1]),
		},
	}
}

func (b *SystemBackend) pluginsCatalogCRUDPath() *framework.Path {
	return &framework.Path{
		Pattern: "plugins/catalog(/(?P<type>auth|database|secret))?/(?P<name>.+)",

		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_name"][0]),
			},
			"type": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_type"][0]),
			},
			"sha256": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_sha-256"][0]),
			},
			"sha_256": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_sha-256"][0]),
			},
			"command": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_command"][0]),
			},
			"args": &framework.FieldSchema{
				Type:        framework.TypeStringSlice,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_args"][0]),
			},
			"env": &framework.FieldSchema{
				Type:        framework.TypeStringSlice,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_env"][0]),
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.handlePluginCatalogUpdate,
				Summary:  "Register a new plugin, or updates an existing one with the supplied name.",
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.handlePluginCatalogDelete,
				Summary:  "Remove the plugin with the given name.",
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handlePluginCatalogRead,
				Summary:  "Return the configuration data for the plugin with the given name.",
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

			Fields: map[string]*framework.FieldSchema{
				"type": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["plugin-catalog_type"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handlePluginCatalogTypedList,
					Summary:  "List the plugins in the catalog.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-catalog"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["plugin-catalog"][1]),
		},
		{
			Pattern: "plugins/catalog/?$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.handlePluginCatalogUntypedList,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-catalog-list-all"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["plugin-catalog-list-all"][1]),
		},
	}
}

func (b *SystemBackend) pluginsReloadPath() *framework.Path {
	return &framework.Path{
		Pattern: "plugins/reload/backend$",

		Fields: map[string]*framework.FieldSchema{
			"plugin": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-backend-reload-plugin"][0]),
			},
			"mounts": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: strings.TrimSpace(sysHelp["plugin-backend-reload-mounts"][0]),
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:    b.handlePluginReloadUpdate,
				Summary:     "Reload mounted plugin backends.",
				Description: "Either the plugin name (`plugin`) or the desired plugin backend mounts (`mounts`) must be provided, but not both. In the case that the plugin name is provided, all mounted paths that use that plugin backend will be reloaded.",
			},
		},

		HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-reload"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["plugin-reload"][1]),
	}
}

func (b *SystemBackend) toolsPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "tools/hash" + framework.OptionalParamRegex("urlalgorithm"),
			Fields: map[string]*framework.FieldSchema{
				"input": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The base64-encoded input data",
				},

				"algorithm": &framework.FieldSchema{
					Type:    framework.TypeString,
					Default: "sha2-256",
					Description: `Algorithm to use (POST body parameter). Valid values are:

			* sha2-224
			* sha2-256
			* sha2-384
			* sha2-512

			Defaults to "sha2-256".`,
				},

				"urlalgorithm": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: `Algorithm to use (POST URL parameter)`,
				},

				"format": &framework.FieldSchema{
					Type:        framework.TypeString,
					Default:     "hex",
					Description: `Encoding format to use. Can be "hex" or "base64". Defaults to "hex".`,
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathHashWrite,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["hash"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["hash"][1]),
		},

		{
			Pattern: "tools/random" + framework.OptionalParamRegex("urlbytes"),
			Fields: map[string]*framework.FieldSchema{
				"urlbytes": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The number of bytes to generate (POST URL parameter)",
				},

				"bytes": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Default:     32,
					Description: "The number of bytes to generate (POST body parameter). Defaults to 32 (256 bits).",
				},

				"format": &framework.FieldSchema{
					Type:        framework.TypeString,
					Default:     "base64",
					Description: `Encoding format to use. Can be "hex" or "base64". Defaults to "base64".`,
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRandomWrite,
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
			Fields: map[string]*framework.FieldSchema{
				"context": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Context string appended to every operationId",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.pathInternalOpenAPI,
				logical.UpdateOperation: b.pathInternalOpenAPI,
			},
		},
		{
			Pattern: "internal/specs/openapi",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathInternalOpenAPI,
					Summary:  "Generate an OpenAPI 3 document of all mounted paths.",
				},
			},
		},
		{
			Pattern: "internal/ui/mounts",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathInternalUIMountsRead,
					Summary:  "Lists all enabled and visible auth and secrets mounts.",
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-mounts"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-ui-mounts"][1]),
		},
		{
			Pattern: "internal/ui/mounts/(?P<path>.+)",
			Fields: map[string]*framework.FieldSchema{
				"path": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The path of the mount.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathInternalUIMountRead,
					Summary:  "Return information about the given mount.",
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-mounts"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-ui-mounts"][1]),
		},
		{
			Pattern: "internal/ui/namespaces",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    pathInternalUINamespacesRead(b),
					Unpublished: true,
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-namespaces"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-ui-namespaces"][1]),
		},
		{
			Pattern: "internal/ui/resultant-acl",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    b.pathInternalUIResultantACL,
					Unpublished: true,
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-resultant-acl"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-ui-resultant-acl"][1]),
		},
		{
			Pattern: "internal/counters/requests",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    b.pathInternalCountersRequests,
					Unpublished: true,
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-counters-requests"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-counters-requests"][1]),
		},
	}
}

func (b *SystemBackend) capabilitiesPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "capabilities-accessor$",

			Fields: map[string]*framework.FieldSchema{
				"accessor": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Accessor of the token for which capabilities are being queried.",
				},
				"path": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "(DEPRECATED) Path on which capabilities are being queried. Use 'paths' instead.",
					Deprecated:  true,
				},
				"paths": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "Paths on which capabilities are being queried.",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleCapabilitiesAccessor,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["capabilities_accessor"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["capabilities_accessor"][1]),
		},

		{
			Pattern: "capabilities$",

			Fields: map[string]*framework.FieldSchema{
				"token": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Token for which capabilities are being queried.",
				},
				"path": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "(DEPRECATED) Path on which capabilities are being queried. Use 'paths' instead.",
					Deprecated:  true,
				},
				"paths": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "Paths on which capabilities are being queried.",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleCapabilities,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["capabilities"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["capabilities"][1]),
		},

		{
			Pattern: "capabilities-self$",

			Fields: map[string]*framework.FieldSchema{
				"token": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Token for which capabilities are being queried.",
				},
				"path": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "(DEPRECATED) Path on which capabilities are being queried. Use 'paths' instead.",
					Deprecated:  true,
				},
				"paths": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "Paths on which capabilities are being queried.",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleCapabilities,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["capabilities_self"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["capabilities_self"][1]),
		},
	}
}

func (b *SystemBackend) leasePaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "leases/lookup/(?P<prefix>.+?)?",

			Fields: map[string]*framework.FieldSchema{
				"prefix": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["leases-list-prefix"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleLeaseLookupList,
					Summary:  "Returns a list of lease ids.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["leases"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["leases"][1]),
		},

		{
			Pattern: "leases/lookup",

			Fields: map[string]*framework.FieldSchema{
				"lease_id": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["lease_id"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleLeaseLookup,
					Summary:  "Retrieve lease metadata.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["leases"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["leases"][1]),
		},

		{
			Pattern: "(leases/)?renew" + framework.OptionalParamRegex("url_lease_id"),

			Fields: map[string]*framework.FieldSchema{
				"url_lease_id": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["lease_id"][0]),
				},
				"lease_id": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["lease_id"][0]),
				},
				"increment": &framework.FieldSchema{
					Type:        framework.TypeDurationSecond,
					Description: strings.TrimSpace(sysHelp["increment"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRenew,
					Summary:  "Renews a lease, requesting to extend the lease.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["renew"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["renew"][1]),
		},

		{
			Pattern: "(leases/)?revoke" + framework.OptionalParamRegex("url_lease_id"),

			Fields: map[string]*framework.FieldSchema{
				"url_lease_id": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["lease_id"][0]),
				},
				"lease_id": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["lease_id"][0]),
				},
				"sync": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Default:     true,
					Description: strings.TrimSpace(sysHelp["revoke-sync"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRevoke,
					Summary:  "Revokes a lease immediately.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["revoke"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["revoke"][1]),
		},

		{
			Pattern: "(leases/)?revoke-force/(?P<prefix>.+)",

			Fields: map[string]*framework.FieldSchema{
				"prefix": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["revoke-force-path"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:    b.handleRevokeForce,
					Summary:     "Revokes all secrets or tokens generated under a given prefix immediately",
					Description: "Unlike `/sys/leases/revoke-prefix`, this path ignores backend errors encountered during revocation. This is potentially very dangerous and should only be used in specific emergency situations where errors in the backend or the connected backend service prevent normal revocation.\n\nBy ignoring these errors, Vault abdicates responsibility for ensuring that the issued credentials or secrets are properly revoked and/or cleaned up. Access to this endpoint should be tightly controlled.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["revoke-force"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["revoke-force"][1]),
		},

		{
			Pattern: "(leases/)?revoke-prefix/(?P<prefix>.+)",

			Fields: map[string]*framework.FieldSchema{
				"prefix": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["revoke-prefix-path"][0]),
				},
				"sync": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Default:     true,
					Description: strings.TrimSpace(sysHelp["revoke-sync"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRevokePrefix,
					Summary:  "Revokes all secrets (via a lease ID prefix) or tokens (via the tokens' path property) generated under a given prefix immediately.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["revoke-prefix"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["revoke-prefix"][1]),
		},

		{
			Pattern: "leases/tidy$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleTidyLeases,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["tidy_leases"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["tidy_leases"][1]),
		},
	}
}

func (b *SystemBackend) remountPath() *framework.Path {
	return &framework.Path{
		Pattern: "remount",

		Fields: map[string]*framework.FieldSchema{
			"from": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The previous mount point.",
			},
			"to": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The new mount point.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handleRemount,
		},

		HelpSynopsis:    strings.TrimSpace(sysHelp["remount"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["remount"][1]),
	}
}

func (b *SystemBackend) metricsPath() *framework.Path {
	return &framework.Path{
		Pattern: "metrics",
		Fields: map[string]*framework.FieldSchema{
			"format": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Format to export metrics into. Currently accepts only \"prometheus\".",
				Query:       true,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.handleMetrics,
		},
		HelpSynopsis:    strings.TrimSpace(sysHelp["metrics"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["metrics"][1]),
	}

}

func (b *SystemBackend) authPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "auth$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.handleAuthTable,
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["auth-table"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["auth-table"][1]),
		},
		{
			Pattern: "auth/(?P<path>.+?)/tune$",
			Fields: map[string]*framework.FieldSchema{
				"path": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_tune"][0]),
				},
				"default_lease_ttl": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["tune_default_lease_ttl"][0]),
				},
				"max_lease_ttl": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["tune_max_lease_ttl"][0]),
				},
				"description": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_desc"][0]),
				},
				"audit_non_hmac_request_keys": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["tune_audit_non_hmac_request_keys"][0]),
				},
				"audit_non_hmac_response_keys": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["tune_audit_non_hmac_response_keys"][0]),
				},
				"options": &framework.FieldSchema{
					Type:        framework.TypeKVPairs,
					Description: strings.TrimSpace(sysHelp["tune_mount_options"][0]),
				},
				"listing_visibility": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["listing_visibility"][0]),
				},
				"passthrough_request_headers": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["passthrough_request_headers"][0]),
				},
				"allowed_response_headers": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["allowed_response_headers"][0]),
				},
				"token_type": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["token_type"][0]),
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    b.handleAuthTuneRead,
					Summary:     "Reads the given auth path's configuration.",
					Description: "This endpoint requires sudo capability on the final path, but the same functionality can be achieved without sudo via `sys/mounts/auth/[auth-path]/tune`.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback:    b.handleAuthTuneWrite,
					Summary:     "Tune configuration parameters for a given auth path.",
					Description: "This endpoint requires sudo capability on the final path, but the same functionality can be achieved without sudo via `sys/mounts/auth/[auth-path]/tune`.",
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["auth_tune"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["auth_tune"][1]),
		},
		{
			Pattern: "auth/(?P<path>.+)",
			Fields: map[string]*framework.FieldSchema{
				"path": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_path"][0]),
				},
				"type": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_type"][0]),
				},
				"description": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_desc"][0]),
				},
				"config": &framework.FieldSchema{
					Type:        framework.TypeMap,
					Description: strings.TrimSpace(sysHelp["auth_config"][0]),
				},
				"local": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Default:     false,
					Description: strings.TrimSpace(sysHelp["mount_local"][0]),
				},
				"seal_wrap": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Default:     false,
					Description: strings.TrimSpace(sysHelp["seal_wrap"][0]),
				},
				"plugin_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_plugin"][0]),
				},
				"options": &framework.FieldSchema{
					Type:        framework.TypeKVPairs,
					Description: strings.TrimSpace(sysHelp["auth_options"][0]),
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleEnableAuth,
					Summary:  "Enables a new auth method.",
					Description: `After enabling, the auth method can be accessed and configured via the auth path specified as part of the URL. This auth path will be nested under the auth prefix.

For example, enable the "foo" auth method will make it accessible at /auth/foo.`,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDisableAuth,
					Summary:  "Disable the auth method at the given auth path",
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.handlePoliciesList(PolicyTypeACL),
				logical.ListOperation: b.handlePoliciesList(PolicyTypeACL),
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["policy-list"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["policy-list"][1]),
		},

		{
			Pattern: "policy/(?P<name>.+)",

			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["policy-name"][0]),
				},
				"rules": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["policy-rules"][0]),
					Deprecated:  true,
				},
				"policy": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["policy-rules"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handlePoliciesRead(PolicyTypeACL),
					Summary:  "Retrieve the policy body for the named policy.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handlePoliciesSet(PolicyTypeACL),
					Summary:  "Add a new or update an existing policy.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handlePoliciesDelete(PolicyTypeACL),
					Summary:  "Delete the policy with the given name.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["policy"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["policy"][1]),
		},

		{
			Pattern: "policies/acl/?$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.handlePoliciesList(PolicyTypeACL),
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["policy-list"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["policy-list"][1]),
		},

		{
			Pattern: "policies/acl/(?P<name>.+)",

			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["policy-name"][0]),
				},
				"policy": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["policy-rules"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handlePoliciesRead(PolicyTypeACL),
					Summary:  "Retrieve information about the named ACL policy.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handlePoliciesSet(PolicyTypeACL),
					Summary:  "Add a new or update an existing ACL policy.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handlePoliciesDelete(PolicyTypeACL),
					Summary:  "Delete the ACL policy with the given name.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["policy"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["policy"][1]),
		},
	}
}

func (b *SystemBackend) wrappingPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "wrapping/wrap$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleWrappingWrap,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["wrap"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["wrap"][1]),
		},

		{
			Pattern: "wrapping/unwrap$",

			Fields: map[string]*framework.FieldSchema{
				"token": &framework.FieldSchema{
					Type: framework.TypeString,
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleWrappingUnwrap,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["unwrap"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["unwrap"][1]),
		},

		{
			Pattern: "wrapping/lookup$",

			Fields: map[string]*framework.FieldSchema{
				"token": &framework.FieldSchema{
					Type: framework.TypeString,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrappingLookup,
					Summary:  "Look up wrapping properties for the given token.",
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleWrappingLookup,
					Summary:  "Look up wrapping properties for the requester's token.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["wraplookup"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["wraplookup"][1]),
		},

		{
			Pattern: "wrapping/rewrap$",

			Fields: map[string]*framework.FieldSchema{
				"token": &framework.FieldSchema{
					Type: framework.TypeString,
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleWrappingRewrap,
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

			Fields: map[string]*framework.FieldSchema{
				"path": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["mount_path"][0]),
				},
				"default_lease_ttl": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["tune_default_lease_ttl"][0]),
				},
				"max_lease_ttl": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["tune_max_lease_ttl"][0]),
				},
				"description": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["auth_desc"][0]),
				},
				"audit_non_hmac_request_keys": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["tune_audit_non_hmac_request_keys"][0]),
				},
				"audit_non_hmac_response_keys": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["tune_audit_non_hmac_response_keys"][0]),
				},
				"options": &framework.FieldSchema{
					Type:        framework.TypeKVPairs,
					Description: strings.TrimSpace(sysHelp["tune_mount_options"][0]),
				},
				"listing_visibility": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["listing_visibility"][0]),
				},
				"passthrough_request_headers": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["passthrough_request_headers"][0]),
				},
				"allowed_response_headers": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: strings.TrimSpace(sysHelp["allowed_response_headers"][0]),
				},
				"token_type": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["token_type"][0]),
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.handleMountTuneRead,
				logical.UpdateOperation: b.handleMountTuneWrite,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["mount_tune"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["mount_tune"][1]),
		},

		{
			Pattern: "mounts/(?P<path>.+?)",

			Fields: map[string]*framework.FieldSchema{
				"path": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["mount_path"][0]),
				},
				"type": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["mount_type"][0]),
				},
				"description": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["mount_desc"][0]),
				},
				"config": &framework.FieldSchema{
					Type:        framework.TypeMap,
					Description: strings.TrimSpace(sysHelp["mount_config"][0]),
				},
				"local": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Default:     false,
					Description: strings.TrimSpace(sysHelp["mount_local"][0]),
				},
				"seal_wrap": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Default:     false,
					Description: strings.TrimSpace(sysHelp["seal_wrap"][0]),
				},
				"plugin_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["mount_plugin_name"][0]),
				},
				"options": &framework.FieldSchema{
					Type:        framework.TypeKVPairs,
					Description: strings.TrimSpace(sysHelp["mount_options"][0]),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleMount,
					Summary:  "Enable a new secrets engine at the given path.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleUnmount,
					Summary:  "Disable the mount point specified at the given path.",
				},
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["mount"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["mount"][1]),
		},

		{
			Pattern: "mounts$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.handleMountTable,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["mounts"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["mounts"][1]),
		},
	}
}
