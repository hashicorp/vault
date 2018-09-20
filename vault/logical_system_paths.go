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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.handleCORSRead,
				logical.UpdateOperation: b.handleCORSUpdate,
				logical.DeleteOperation: b.handleCORSDelete,
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.handleConfigUIHeadersRead,
				logical.UpdateOperation: b.handleConfigUIHeadersUpdate,
				logical.DeleteOperation: b.handleConfigUIHeadersDelete,
			},

			HelpDescription: strings.TrimSpace(sysHelp["config/ui/headers"][0]),
			HelpSynopsis:    strings.TrimSpace(sysHelp["config/ui/headers"][1]),
		},

		{
			Pattern: "config/ui/headers/$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.handleConfigUIHeadersList,
			},

			HelpDescription: strings.TrimSpace(sysHelp["config/ui/headers"][0]),
			HelpSynopsis:    strings.TrimSpace(sysHelp["config/ui/headers"][1]),
		},

		{
			Pattern:         "generate-root(/attempt)?$",
			HelpSynopsis:    strings.TrimSpace(sysHelp["generate-root"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["generate-root"][1]),
		},

		{
			Pattern:         "init$",
			HelpSynopsis:    strings.TrimSpace(sysHelp["init"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["init"][1]),
		},
	}
}

func (b *SystemBackend) rekeyPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "rekey/backup$",

			Fields: map[string]*framework.FieldSchema{},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.handleRekeyRetrieveBarrier,
				logical.DeleteOperation: b.handleRekeyDeleteBarrier,
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
			Pattern:         "seal-status$",
			HelpSynopsis:    strings.TrimSpace(sysHelp["seal-status"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["seal-status"][1]),
		},

		{
			Pattern:         "seal$",
			HelpSynopsis:    strings.TrimSpace(sysHelp["seal"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["seal"][1]),
		},

		{
			Pattern:         "unseal$",
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.handleAuditTable,
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleEnableAudit,
				logical.DeleteOperation: b.handleDisableAudit,
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleAuditedHeaderUpdate,
				logical.DeleteOperation: b.handleAuditedHeaderDelete,
				logical.ReadOperation:   b.handleAuditedHeaderRead,
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["audited-headers-name"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["audited-headers-name"][1]),
		},

		{
			Pattern: "config/auditing/request-headers$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.handleAuditedHeadersRead,
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

func (b *SystemBackend) pluginsCatalogPath() *framework.Path {
	return &framework.Path{
		Pattern: "plugins/catalog/(?P<name>.+)",

		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: strings.TrimSpace(sysHelp["plugin-catalog_name"][0]),
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
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handlePluginCatalogUpdate,
			logical.DeleteOperation: b.handlePluginCatalogDelete,
			logical.ReadOperation:   b.handlePluginCatalogRead,
		},

		HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-catalog"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["plugin-catalog"][1]),
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

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.handlePluginReloadUpdate,
		},

		HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-reload"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["plugin-reload"][1]),
	}
}

func (b *SystemBackend) pluginsCatalogListPath() *framework.Path {
	return &framework.Path{
		Pattern: "plugins/catalog/?$",

		Fields: map[string]*framework.FieldSchema{},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.handlePluginCatalogList,
		},

		HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-catalog"][0]),
		HelpDescription: strings.TrimSpace(sysHelp["plugin-catalog"][1]),
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

func (b *SystemBackend) internalUIPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "internal/ui/mounts",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathInternalUIMountsRead,
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
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathInternalUIMountRead,
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-mounts"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-ui-mounts"][1]),
		},
		{
			Pattern: "internal/ui/namespaces",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: pathInternalUINamespacesRead(b),
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-namespaces"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-ui-namespaces"][1]),
		},
		{
			Pattern: "internal/ui/resultant-acl",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathInternalUIResultantACL,
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-resultant-acl"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["internal-ui-resultant-acl"][1]),
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.handleLeaseLookupList,
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleLeaseLookup,
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleRenew,
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleRevoke,
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleRevokeForce,
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleRevokePrefix,
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
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.handleAuthTuneRead,
				logical.UpdateOperation: b.handleAuthTuneWrite,
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
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleEnableAuth,
				logical.DeleteOperation: b.handleDisableAuth,
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
				},
				"policy": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: strings.TrimSpace(sysHelp["policy-rules"][0]),
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.handlePoliciesRead(PolicyTypeACL),
				logical.UpdateOperation: b.handlePoliciesSet(PolicyTypeACL),
				logical.DeleteOperation: b.handlePoliciesDelete(PolicyTypeACL),
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.handlePoliciesRead(PolicyTypeACL),
				logical.UpdateOperation: b.handlePoliciesSet(PolicyTypeACL),
				logical.DeleteOperation: b.handlePoliciesDelete(PolicyTypeACL),
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleWrappingLookup,
				logical.ReadOperation:   b.handleWrappingLookup,
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

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.handleMount,
				logical.DeleteOperation: b.handleUnmount,
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
