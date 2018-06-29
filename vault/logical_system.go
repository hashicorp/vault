package vault

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/compressutil"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/helper/wrapping"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

var (
	// protectedPaths cannot be accessed via the raw APIs.
	// This is both for security and to prevent disrupting Vault.
	protectedPaths = []string{
		keyringPath,
		coreLocalClusterInfoPath,
	}

	replicationPaths = func(b *SystemBackend) []*framework.Path {
		return []*framework.Path{
			&framework.Path{
				Pattern: "replication/status",
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
						resp := &logical.Response{
							Data: map[string]interface{}{
								"mode": "disabled",
							},
						}
						return resp, nil
					},
				},
			},
		}
	}
)

func NewSystemBackend(core *Core, logger log.Logger) *SystemBackend {
	b := &SystemBackend{
		Core:   core,
		logger: logger,
	}

	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(sysHelpRoot),

		PathsSpecial: &logical.Paths{
			Root: []string{
				"auth/*",
				"remount",
				"audit",
				"audit/*",
				"raw",
				"raw/*",
				"replication/primary/secondary-token",
				"replication/reindex",
				"rotate",
				"config/cors",
				"config/auditing/*",
				"config/ui/headers/*",
				"plugins/catalog/*",
				"revoke-prefix/*",
				"revoke-force/*",
				"leases/revoke-prefix/*",
				"leases/revoke-force/*",
				"leases/lookup/*",
			},

			Unauthenticated: []string{
				"wrapping/lookup",
				"wrapping/pubkey",
				"replication/status",
				"internal/ui/mounts",
				"internal/ui/mounts/*",
			},
		},

		Paths: []*framework.Path{
			&framework.Path{
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

			&framework.Path{
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

			&framework.Path{
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

			&framework.Path{
				Pattern: "config/ui/headers/$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ListOperation: b.handleConfigUIHeadersList,
				},

				HelpDescription: strings.TrimSpace(sysHelp["config/ui/headers"][0]),
				HelpSynopsis:    strings.TrimSpace(sysHelp["config/ui/headers"][1]),
			},

			&framework.Path{
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

			&framework.Path{
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

			&framework.Path{
				Pattern:         "generate-root(/attempt)?$",
				HelpSynopsis:    strings.TrimSpace(sysHelp["generate-root"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["generate-root"][1]),
			},

			&framework.Path{
				Pattern:         "init$",
				HelpSynopsis:    strings.TrimSpace(sysHelp["init"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["init"][1]),
			},

			&framework.Path{
				Pattern: "rekey/backup$",

				Fields: map[string]*framework.FieldSchema{},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation:   b.handleRekeyRetrieveBarrier,
					logical.DeleteOperation: b.handleRekeyDeleteBarrier,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["rekey_backup"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["rekey_backup"][0]),
			},

			&framework.Path{
				Pattern: "rekey/recovery-key-backup$",

				Fields: map[string]*framework.FieldSchema{},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation:   b.handleRekeyRetrieveRecovery,
					logical.DeleteOperation: b.handleRekeyDeleteRecovery,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["rekey_backup"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["rekey_backup"][0]),
			},

			&framework.Path{
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

			&framework.Path{
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

			&framework.Path{
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

			&framework.Path{
				Pattern: "mounts$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.handleMountTable,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["mounts"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["mounts"][1]),
			},

			&framework.Path{
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
			},

			&framework.Path{
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

			&framework.Path{
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

			&framework.Path{
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

			&framework.Path{
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
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.handleRevoke,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["revoke"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["revoke"][1]),
			},

			&framework.Path{
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

			&framework.Path{
				Pattern: "(leases/)?revoke-prefix/(?P<prefix>.+)",

				Fields: map[string]*framework.FieldSchema{
					"prefix": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["revoke-prefix-path"][0]),
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.handleRevokePrefix,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["revoke-prefix"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["revoke-prefix"][1]),
			},

			&framework.Path{
				Pattern: "leases/tidy$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.handleTidyLeases,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["tidy_leases"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["tidy_leases"][1]),
			},

			&framework.Path{
				Pattern: "auth$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.handleAuthTable,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["auth-table"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["auth-table"][1]),
			},

			&framework.Path{
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

			&framework.Path{
				Pattern: "policy/?$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.handlePoliciesList(PolicyTypeACL),
					logical.ListOperation: b.handlePoliciesList(PolicyTypeACL),
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["policy-list"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["policy-list"][1]),
			},

			&framework.Path{
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

			&framework.Path{
				Pattern: "policies/acl/?$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ListOperation: b.handlePoliciesList(PolicyTypeACL),
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["policy-list"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["policy-list"][1]),
			},

			&framework.Path{
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

			&framework.Path{
				Pattern:         "seal-status$",
				HelpSynopsis:    strings.TrimSpace(sysHelp["seal-status"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["seal-status"][1]),
			},

			&framework.Path{
				Pattern:         "seal$",
				HelpSynopsis:    strings.TrimSpace(sysHelp["seal"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["seal"][1]),
			},

			&framework.Path{
				Pattern:         "unseal$",
				HelpSynopsis:    strings.TrimSpace(sysHelp["unseal"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["unseal"][1]),
			},

			&framework.Path{
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

			&framework.Path{
				Pattern: "audit$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.handleAuditTable,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["audit-table"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["audit-table"][1]),
			},

			&framework.Path{
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

			&framework.Path{
				Pattern: "key-status$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.handleKeyStatus,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["key-status"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["key-status"][1]),
			},

			&framework.Path{
				Pattern: "rotate$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.handleRotate,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["rotate"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["rotate"][1]),
			},

			&framework.Path{
				Pattern: "wrapping/wrap$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.handleWrappingWrap,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["wrap"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["wrap"][1]),
			},

			&framework.Path{
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

			&framework.Path{
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

			&framework.Path{
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

			&framework.Path{
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

			&framework.Path{
				Pattern: "config/auditing/request-headers$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.handleAuditedHeadersRead,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["audited-headers"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["audited-headers"][1]),
			},

			&framework.Path{
				Pattern: "plugins/catalog/?$",

				Fields: map[string]*framework.FieldSchema{},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ListOperation: b.handlePluginCatalogList,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["plugin-catalog"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["plugin-catalog"][1]),
			},

			&framework.Path{
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
			},
			&framework.Path{
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
			},
			&framework.Path{
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

			&framework.Path{
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
			&framework.Path{
				Pattern: "internal/ui/mounts",
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.pathInternalUIMountsRead,
				},
				HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-mounts"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["internal-ui-mounts"][1]),
			},
			&framework.Path{
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
			&framework.Path{
				Pattern: "internal/ui/resultant-acl",
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.pathInternalUIResultantACL,
				},
				HelpSynopsis:    strings.TrimSpace(sysHelp["internal-ui-resultant-acl"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["internal-ui-resultant-acl"][1]),
			},
		},
	}

	b.Backend.Paths = append(b.Backend.Paths, replicationPaths(b)...)

	if core.rawEnabled {
		b.Backend.Paths = append(b.Backend.Paths, &framework.Path{
			Pattern: "(raw/?$|raw/(?P<path>.+))",

			Fields: map[string]*framework.FieldSchema{
				"path": &framework.FieldSchema{
					Type: framework.TypeString,
				},
				"value": &framework.FieldSchema{
					Type: framework.TypeString,
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.handleRawRead,
				logical.UpdateOperation: b.handleRawWrite,
				logical.DeleteOperation: b.handleRawDelete,
				logical.ListOperation:   b.handleRawList,
			},
			HelpSynopsis:    strings.TrimSpace(sysHelp["raw"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["raw"][1]),
		})
	}

	b.Backend.Invalidate = b.invalidate

	return b
}

// SystemBackend implements logical.Backend and is used to interact with
// the core of the system. This backend is hardcoded to exist at the "sys"
// prefix. Conceptually it is similar to procfs on Linux.
type SystemBackend struct {
	*framework.Backend
	Core   *Core
	logger log.Logger
}

// handleCORSRead returns the current CORS configuration
func (b *SystemBackend) handleCORSRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	corsConf := b.Core.corsConfig

	enabled := corsConf.IsEnabled()

	resp := &logical.Response{
		Data: map[string]interface{}{
			"enabled": enabled,
		},
	}

	if enabled {
		corsConf.RLock()
		resp.Data["allowed_origins"] = corsConf.AllowedOrigins
		resp.Data["allowed_headers"] = corsConf.AllowedHeaders
		corsConf.RUnlock()
	}

	return resp, nil
}

// handleCORSUpdate sets the list of origins that are allowed to make
// cross-origin requests and sets the CORS enabled flag to true
func (b *SystemBackend) handleCORSUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	origins := d.Get("allowed_origins").([]string)
	headers := d.Get("allowed_headers").([]string)

	return nil, b.Core.corsConfig.Enable(ctx, origins, headers)
}

// handleCORSDelete sets the CORS enabled flag to false and clears the list of
// allowed origins & headers.
func (b *SystemBackend) handleCORSDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return nil, b.Core.corsConfig.Disable(ctx)
}

func (b *SystemBackend) handleTidyLeases(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	go func() {
		err := b.Core.expiration.Tidy()
		if err != nil {
			b.Backend.Logger().Error("failed to tidy leases", "error", err)
			return
		}
	}()

	resp := &logical.Response{}
	resp.AddWarning("Tidy operation successfully started. Any information from the operation will be printed to Vault's server logs.")
	return resp, nil
}

func (b *SystemBackend) invalidate(ctx context.Context, key string) {
	/*
		if b.Core.logger.IsTrace() {
			b.Core.logger.Trace("invalidating key", "key", key)
		}
	*/
	switch {
	case strings.HasPrefix(key, policyACLSubPath):
		b.Core.stateLock.RLock()
		defer b.Core.stateLock.RUnlock()
		if b.Core.policyStore != nil {
			b.Core.policyStore.invalidate(ctx, strings.TrimPrefix(key, policyACLSubPath), PolicyTypeACL)
		}
	case strings.HasPrefix(key, tokenSubPath):
		b.Core.stateLock.RLock()
		defer b.Core.stateLock.RUnlock()
		if b.Core.tokenStore != nil {
			b.Core.tokenStore.Invalidate(ctx, key)
		}
	}
}

func (b *SystemBackend) handlePluginCatalogList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	plugins, err := b.Core.pluginCatalog.List(ctx)
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(plugins), nil
}

func (b *SystemBackend) handlePluginCatalogUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	pluginName := d.Get("name").(string)
	if pluginName == "" {
		return logical.ErrorResponse("missing plugin name"), nil
	}

	sha256 := d.Get("sha256").(string)
	if sha256 == "" {
		sha256 = d.Get("sha_256").(string)
		if sha256 == "" {
			return logical.ErrorResponse("missing SHA-256 value"), nil
		}
	}

	command := d.Get("command").(string)
	if command == "" {
		return logical.ErrorResponse("missing command value"), nil
	}

	// For backwards compatibility, also accept args as part of command.  Don't
	// accepts args in both command and args.
	args := d.Get("args").([]string)
	parts := strings.Split(command, " ")
	if len(parts) <= 0 {
		return logical.ErrorResponse("missing command value"), nil
	} else if len(parts) > 1 && len(args) > 0 {
		return logical.ErrorResponse("must not specify args in command and args field"), nil
	} else if len(parts) > 1 {
		args = parts[1:]
	}

	sha256Bytes, err := hex.DecodeString(sha256)
	if err != nil {
		return logical.ErrorResponse("Could not decode SHA-256 value from Hex"), err
	}

	err = b.Core.pluginCatalog.Set(ctx, pluginName, parts[0], args, sha256Bytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *SystemBackend) handlePluginCatalogRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	pluginName := d.Get("name").(string)
	if pluginName == "" {
		return logical.ErrorResponse("missing plugin name"), nil
	}
	plugin, err := b.Core.pluginCatalog.Get(ctx, pluginName)
	if err != nil {
		return nil, err
	}
	if plugin == nil {
		return nil, nil
	}

	command := ""
	if !plugin.Builtin {
		command, err = filepath.Rel(b.Core.pluginCatalog.directory, plugin.Command)
		if err != nil {
			return nil, err
		}
	}

	data := map[string]interface{}{
		"name":    plugin.Name,
		"args":    plugin.Args,
		"command": command,
		"sha256":  hex.EncodeToString(plugin.Sha256),
		"builtin": plugin.Builtin,
	}

	return &logical.Response{
		Data: data,
	}, nil
}

func (b *SystemBackend) handlePluginCatalogDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	pluginName := d.Get("name").(string)
	if pluginName == "" {
		return logical.ErrorResponse("missing plugin name"), nil
	}
	err := b.Core.pluginCatalog.Delete(ctx, pluginName)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *SystemBackend) handlePluginReloadUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	pluginName := d.Get("plugin").(string)
	pluginMounts := d.Get("mounts").([]string)

	if pluginName != "" && len(pluginMounts) > 0 {
		return logical.ErrorResponse("plugin and mounts cannot be set at the same time"), nil
	}
	if pluginName == "" && len(pluginMounts) == 0 {
		return logical.ErrorResponse("plugin or mounts must be provided"), nil
	}

	if pluginName != "" {
		err := b.Core.reloadMatchingPlugin(ctx, pluginName)
		if err != nil {
			return nil, err
		}
	} else if len(pluginMounts) > 0 {
		err := b.Core.reloadMatchingPluginMounts(ctx, pluginMounts)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// handleAuditedHeaderUpdate creates or overwrites a header entry
func (b *SystemBackend) handleAuditedHeaderUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	header := d.Get("header").(string)
	hmac := d.Get("hmac").(bool)
	if header == "" {
		return logical.ErrorResponse("missing header name"), nil
	}

	headerConfig := b.Core.AuditedHeadersConfig()
	err := headerConfig.add(ctx, header, hmac)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// handleAuditedHeaderDelete deletes the header with the given name
func (b *SystemBackend) handleAuditedHeaderDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	header := d.Get("header").(string)
	if header == "" {
		return logical.ErrorResponse("missing header name"), nil
	}

	headerConfig := b.Core.AuditedHeadersConfig()
	err := headerConfig.remove(ctx, header)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// handleAuditedHeaderRead returns the header configuration for the given header name
func (b *SystemBackend) handleAuditedHeaderRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	header := d.Get("header").(string)
	if header == "" {
		return logical.ErrorResponse("missing header name"), nil
	}

	headerConfig := b.Core.AuditedHeadersConfig()
	settings, ok := headerConfig.Headers[strings.ToLower(header)]
	if !ok {
		return logical.ErrorResponse("Could not find header in config"), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			header: settings,
		},
	}, nil
}

// handleAuditedHeadersRead returns the whole audited headers config
func (b *SystemBackend) handleAuditedHeadersRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	headerConfig := b.Core.AuditedHeadersConfig()

	return &logical.Response{
		Data: map[string]interface{}{
			"headers": headerConfig.Headers,
		},
	}, nil
}

// handleCapabilitiesAccessor returns the ACL capabilities of the
// token associated with the given accessor for a given path.
func (b *SystemBackend) handleCapabilitiesAccessor(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	accessor := d.Get("accessor").(string)
	if accessor == "" {
		return logical.ErrorResponse("missing accessor"), nil
	}

	aEntry, err := b.Core.tokenStore.lookupByAccessor(ctx, accessor, false)
	if err != nil {
		return nil, err
	}

	d.Raw["token"] = aEntry.TokenID
	return b.handleCapabilities(ctx, req, d)
}

// handleCapabilities returns the ACL capabilities of the token for a given path
func (b *SystemBackend) handleCapabilities(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var token string
	if strings.HasSuffix(req.Path, "capabilities-self") {
		token = req.ClientToken
	} else {
		tokenRaw, ok := d.Raw["token"]
		if ok {
			token, _ = tokenRaw.(string)
		}
	}
	if token == "" {
		return nil, fmt.Errorf("no token found")
	}

	ret := &logical.Response{
		Data: map[string]interface{}{},
	}

	paths := d.Get("paths").([]string)
	if len(paths) == 0 {
		// Read from the deprecated field
		paths = d.Get("path").([]string)
	}

	if len(paths) == 0 {
		return logical.ErrorResponse("paths must be supplied"), nil
	}

	for _, path := range paths {
		pathCap, err := b.Core.Capabilities(ctx, token, path)
		if err != nil {
			if !strings.HasSuffix(req.Path, "capabilities-self") && errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
				return nil, &logical.StatusBadRequest{Err: "invalid token"}
			}
			return nil, err
		}
		ret.Data[path] = pathCap
	}

	// This is only here for backwards compatibility
	if len(paths) == 1 {
		ret.Data["capabilities"] = ret.Data[paths[0]]
	}

	return ret, nil
}

// handleRekeyRetrieve returns backed-up, PGP-encrypted unseal keys from a
// rekey operation
func (b *SystemBackend) handleRekeyRetrieve(
	ctx context.Context,
	req *logical.Request,
	data *framework.FieldData,
	recovery bool) (*logical.Response, error) {
	backup, err := b.Core.RekeyRetrieveBackup(ctx, recovery)
	if err != nil {
		return nil, errwrap.Wrapf("unable to look up backed-up keys: {{err}}", err)
	}
	if backup == nil {
		return logical.ErrorResponse("no backed-up keys found"), nil
	}

	keysB64 := map[string][]string{}
	for k, v := range backup.Keys {
		for _, j := range v {
			currB64Keys := keysB64[k]
			if currB64Keys == nil {
				currB64Keys = []string{}
			}
			key, err := hex.DecodeString(j)
			if err != nil {
				return nil, errwrap.Wrapf("error decoding hex-encoded backup key: {{err}}", err)
			}
			currB64Keys = append(currB64Keys, base64.StdEncoding.EncodeToString(key))
			keysB64[k] = currB64Keys
		}
	}

	// Format the status
	resp := &logical.Response{
		Data: map[string]interface{}{
			"nonce":       backup.Nonce,
			"keys":        backup.Keys,
			"keys_base64": keysB64,
		},
	}

	return resp, nil
}

func (b *SystemBackend) handleRekeyRetrieveBarrier(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRekeyRetrieve(ctx, req, data, false)
}

func (b *SystemBackend) handleRekeyRetrieveRecovery(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRekeyRetrieve(ctx, req, data, true)
}

// handleRekeyDelete deletes backed-up, PGP-encrypted unseal keys from a rekey
// operation
func (b *SystemBackend) handleRekeyDelete(
	ctx context.Context,
	req *logical.Request,
	data *framework.FieldData,
	recovery bool) (*logical.Response, error) {
	err := b.Core.RekeyDeleteBackup(ctx, recovery)
	if err != nil {
		return nil, errwrap.Wrapf("error during deletion of backed-up keys: {{err}}", err)
	}

	return nil, nil
}

func (b *SystemBackend) handleRekeyDeleteBarrier(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRekeyDelete(ctx, req, data, false)
}

func (b *SystemBackend) handleRekeyDeleteRecovery(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRekeyDelete(ctx, req, data, true)
}

func mountInfo(entry *MountEntry) map[string]interface{} {
	info := map[string]interface{}{
		"type":        entry.Type,
		"description": entry.Description,
		"accessor":    entry.Accessor,
		"local":       entry.Local,
		"seal_wrap":   entry.SealWrap,
		"options":     entry.Options,
	}
	entryConfig := map[string]interface{}{
		"default_lease_ttl": int64(entry.Config.DefaultLeaseTTL.Seconds()),
		"max_lease_ttl":     int64(entry.Config.MaxLeaseTTL.Seconds()),
		"force_no_cache":    entry.Config.ForceNoCache,
		"plugin_name":       entry.Config.PluginName,
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_request_keys"); ok {
		entryConfig["audit_non_hmac_request_keys"] = rawVal.([]string)
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_response_keys"); ok {
		entryConfig["audit_non_hmac_response_keys"] = rawVal.([]string)
	}
	// Even though empty value is valid for ListingVisibility, we can ignore
	// this case during mount since there's nothing to unset/hide.
	if len(entry.Config.ListingVisibility) > 0 {
		entryConfig["listing_visibility"] = entry.Config.ListingVisibility
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("passthrough_request_headers"); ok {
		entryConfig["passthrough_request_headers"] = rawVal.([]string)
	}

	info["config"] = entryConfig

	return info
}

// handleMountTable handles the "mounts" endpoint to provide the mount table
func (b *SystemBackend) handleMountTable(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Core.mountsLock.RLock()
	defer b.Core.mountsLock.RUnlock()

	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}

	for _, entry := range b.Core.mounts.Entries {
		// Populate mount info
		info := mountInfo(entry)
		resp.Data[entry.Path] = info
	}

	return resp, nil
}

// handleMount is used to mount a new path
func (b *SystemBackend) handleMount(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	repState := b.Core.ReplicationState()

	local := data.Get("local").(bool)
	if !local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return logical.ErrorResponse("cannot add a non-local mount to a replication secondary"), nil
	}

	// Get all the options
	path := data.Get("path").(string)
	path = sanitizeMountPath(path)

	logicalType := data.Get("type").(string)
	description := data.Get("description").(string)
	pluginName := data.Get("plugin_name").(string)
	sealWrap := data.Get("seal_wrap").(bool)
	options := data.Get("options").(map[string]string)

	var config MountConfig
	var apiConfig APIMountConfig

	configMap := data.Get("config").(map[string]interface{})
	if configMap != nil && len(configMap) != 0 {
		err := mapstructure.Decode(configMap, &apiConfig)
		if err != nil {
			return logical.ErrorResponse(
					"unable to convert given mount config information"),
				logical.ErrInvalidRequest
		}
	}

	switch apiConfig.DefaultLeaseTTL {
	case "":
	case "system":
	default:
		tmpDef, err := parseutil.ParseDurationSecond(apiConfig.DefaultLeaseTTL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
					"unable to parse default TTL of %s: %s", apiConfig.DefaultLeaseTTL, err)),
				logical.ErrInvalidRequest
		}
		config.DefaultLeaseTTL = tmpDef
	}

	switch apiConfig.MaxLeaseTTL {
	case "":
	case "system":
	default:
		tmpMax, err := parseutil.ParseDurationSecond(apiConfig.MaxLeaseTTL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
					"unable to parse max TTL of %s: %s", apiConfig.MaxLeaseTTL, err)),
				logical.ErrInvalidRequest
		}
		config.MaxLeaseTTL = tmpMax
	}

	if config.MaxLeaseTTL != 0 && config.DefaultLeaseTTL > config.MaxLeaseTTL {
		return logical.ErrorResponse(
				"given default lease TTL greater than given max lease TTL"),
			logical.ErrInvalidRequest
	}

	if config.DefaultLeaseTTL > b.Core.maxLeaseTTL && config.MaxLeaseTTL == 0 {
		return logical.ErrorResponse(fmt.Sprintf(
				"given default lease TTL greater than system max lease TTL of %d", int(b.Core.maxLeaseTTL.Seconds()))),
			logical.ErrInvalidRequest
	}

	switch logicalType {
	case "":
		return logical.ErrorResponse(
				"backend type must be specified as a string"),
			logical.ErrInvalidRequest

	case "plugin":
		// Only set plugin-name if mount is of type plugin, with apiConfig.PluginName
		// option taking precedence.
		switch {
		case apiConfig.PluginName != "":
			config.PluginName = apiConfig.PluginName
		case pluginName != "":
			config.PluginName = pluginName
		default:
			return logical.ErrorResponse(
					"plugin_name must be provided for plugin backend"),
				logical.ErrInvalidRequest
		}
	}

	switch logicalType {
	case "kv":
	case "kv-v1":
		// Alias KV v1
		logicalType = "kv"
		if options == nil {
			options = map[string]string{}
		}
		options["version"] = "1"

	case "kv-v2":
		// Alias KV v2
		logicalType = "kv"
		if options == nil {
			options = map[string]string{}
		}
		options["version"] = "2"

	default:
		if options != nil && options["version"] != "" {
			return logical.ErrorResponse(fmt.Sprintf(
					"secrets engine %q does not allow setting a version", logicalType)),
				logical.ErrInvalidRequest
		}
	}

	// Copy over the force no cache if set
	if apiConfig.ForceNoCache {
		config.ForceNoCache = true
	}

	if err := checkListingVisibility(apiConfig.ListingVisibility); err != nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid listing_visibility %s", apiConfig.ListingVisibility)), nil
	}
	config.ListingVisibility = apiConfig.ListingVisibility

	if len(apiConfig.AuditNonHMACRequestKeys) > 0 {
		config.AuditNonHMACRequestKeys = apiConfig.AuditNonHMACRequestKeys
	}
	if len(apiConfig.AuditNonHMACResponseKeys) > 0 {
		config.AuditNonHMACResponseKeys = apiConfig.AuditNonHMACResponseKeys
	}
	if len(apiConfig.PassthroughRequestHeaders) > 0 {
		config.PassthroughRequestHeaders = apiConfig.PassthroughRequestHeaders
	}

	// Create the mount entry
	me := &MountEntry{
		Table:       mountTableType,
		Path:        path,
		Type:        logicalType,
		Description: description,
		Config:      config,
		Local:       local,
		SealWrap:    sealWrap,
		Options:     options,
	}

	// Attempt mount
	if err := b.Core.mount(ctx, me); err != nil {
		b.Backend.Logger().Error("mount failed", "path", me.Path, "error", err)
		return handleError(err)
	}

	return nil, nil
}

// used to intercept an HTTPCodedError so it goes back to callee
func handleError(
	err error) (*logical.Response, error) {
	if strings.Contains(err.Error(), logical.ErrReadOnly.Error()) {
		return logical.ErrorResponse(err.Error()), err
	}
	switch err.(type) {
	case logical.HTTPCodedError:
		return logical.ErrorResponse(err.Error()), err
	default:
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
}

// Performs a similar function to handleError, but upon seeing a ReadOnlyError
// will actually strip it out to prevent forwarding
func handleErrorNoReadOnlyForward(
	err error) (*logical.Response, error) {
	if strings.Contains(err.Error(), logical.ErrReadOnly.Error()) {
		return nil, fmt.Errorf("operation could not be completed as storage is read-only")
	}
	switch err.(type) {
	case logical.HTTPCodedError:
		return logical.ErrorResponse(err.Error()), err
	default:
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
}

// handleUnmount is used to unmount a path
func (b *SystemBackend) handleUnmount(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	path = sanitizeMountPath(path)

	repState := b.Core.ReplicationState()
	entry := b.Core.router.MatchingMountEntry(path)
	if entry != nil && !entry.Local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return logical.ErrorResponse("cannot unmount a non-local mount on a replication secondary"), nil
	}

	// We return success when the mount does not exists to not expose if the
	// mount existed or not
	match := b.Core.router.MatchingMount(path)
	if match == "" || path != match {
		return nil, nil
	}

	// Attempt unmount
	if err := b.Core.unmount(ctx, path); err != nil {
		b.Backend.Logger().Error("unmount failed", "path", path, "error", err)
		return handleError(err)
	}

	return nil, nil
}

// handleRemount is used to remount a path
func (b *SystemBackend) handleRemount(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	repState := b.Core.ReplicationState()

	// Get the paths
	fromPath := data.Get("from").(string)
	toPath := data.Get("to").(string)
	if fromPath == "" || toPath == "" {
		return logical.ErrorResponse(
				"both 'from' and 'to' path must be specified as a string"),
			logical.ErrInvalidRequest
	}

	fromPath = sanitizeMountPath(fromPath)
	toPath = sanitizeMountPath(toPath)

	entry := b.Core.router.MatchingMountEntry(fromPath)
	if entry != nil && !entry.Local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return logical.ErrorResponse("cannot remount a non-local mount on a replication secondary"), nil
	}

	// Attempt remount
	if err := b.Core.remount(ctx, fromPath, toPath); err != nil {
		b.Backend.Logger().Error("remount failed", "from_path", fromPath, "to_path", toPath, "error", err)
		return handleError(err)
	}

	return nil, nil
}

// handleAuthTuneRead is used to get config settings on a auth path
func (b *SystemBackend) handleAuthTuneRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse(
				"path must be specified as a string"),
			logical.ErrInvalidRequest
	}
	return b.handleTuneReadCommon("auth/" + path)
}

// handleMountTuneRead is used to get config settings on a backend
func (b *SystemBackend) handleMountTuneRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse(
				"path must be specified as a string"),
			logical.ErrInvalidRequest
	}

	// This call will read both logical backend's configuration as well as auth methods'.
	// Retaining this behavior for backward compatibility. If this behavior is not desired,
	// an error can be returned if path has a prefix of "auth/".
	return b.handleTuneReadCommon(path)
}

// handleTuneReadCommon returns the config settings of a path
func (b *SystemBackend) handleTuneReadCommon(path string) (*logical.Response, error) {
	path = sanitizeMountPath(path)

	sysView := b.Core.router.MatchingSystemView(path)
	if sysView == nil {
		b.Backend.Logger().Error("cannot fetch sysview", "path", path)
		return handleError(fmt.Errorf("sys: cannot fetch sysview for path %q", path))
	}

	mountEntry := b.Core.router.MatchingMountEntry(path)
	if mountEntry == nil {
		b.Backend.Logger().Error("cannot fetch mount entry", "path", path)
		return handleError(fmt.Errorf("sys: cannot fetch mount entry for path %q", path))
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"default_lease_ttl": int(sysView.DefaultLeaseTTL().Seconds()),
			"max_lease_ttl":     int(sysView.MaxLeaseTTL().Seconds()),
			"force_no_cache":    mountEntry.Config.ForceNoCache,
		},
	}

	if rawVal, ok := mountEntry.synthesizedConfigCache.Load("audit_non_hmac_request_keys"); ok {
		resp.Data["audit_non_hmac_request_keys"] = rawVal.([]string)
	}

	if rawVal, ok := mountEntry.synthesizedConfigCache.Load("audit_non_hmac_response_keys"); ok {
		resp.Data["audit_non_hmac_response_keys"] = rawVal.([]string)
	}

	if len(mountEntry.Config.ListingVisibility) > 0 {
		resp.Data["listing_visibility"] = mountEntry.Config.ListingVisibility
	}

	if rawVal, ok := mountEntry.synthesizedConfigCache.Load("passthrough_request_headers"); ok {
		resp.Data["passthrough_request_headers"] = rawVal.([]string)
	}

	if len(mountEntry.Options) > 0 {
		resp.Data["options"] = mountEntry.Options
	}

	return resp, nil
}

// handleAuthTuneWrite is used to set config settings on an auth path
func (b *SystemBackend) handleAuthTuneWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse("path must be specified as a string"),
			logical.ErrInvalidRequest
	}
	return b.handleTuneWriteCommon(ctx, "auth/"+path, data)
}

// handleMountTuneWrite is used to set config settings on a backend
func (b *SystemBackend) handleMountTuneWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse("path must be specified as a string"),
			logical.ErrInvalidRequest
	}
	// This call will write both logical backend's configuration as well as auth methods'.
	// Retaining this behavior for backward compatibility. If this behavior is not desired,
	// an error can be returned if path has a prefix of "auth/".
	return b.handleTuneWriteCommon(ctx, path, data)
}

// handleTuneWriteCommon is used to set config settings on a path
func (b *SystemBackend) handleTuneWriteCommon(ctx context.Context, path string, data *framework.FieldData) (*logical.Response, error) {
	repState := b.Core.ReplicationState()

	path = sanitizeMountPath(path)

	// Prevent protected paths from being changed
	for _, p := range untunableMounts {
		if strings.HasPrefix(path, p) {
			b.Backend.Logger().Error("cannot tune this mount", "path", path)
			return handleError(fmt.Errorf("cannot tune %q", path))
		}
	}

	mountEntry := b.Core.router.MatchingMountEntry(path)
	if mountEntry == nil {
		b.Backend.Logger().Error("tune failed: no mount entry found", "path", path)
		return handleError(fmt.Errorf("tune of path %q failed: no mount entry found", path))
	}
	if mountEntry != nil && !mountEntry.Local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return logical.ErrorResponse("cannot tune a non-local mount on a replication secondary"), nil
	}

	var lock *sync.RWMutex
	switch {
	case strings.HasPrefix(path, credentialRoutePrefix):
		lock = &b.Core.authLock
	default:
		lock = &b.Core.mountsLock
	}

	lock.Lock()
	defer lock.Unlock()

	// Check again after grabbing the lock
	mountEntry = b.Core.router.MatchingMountEntry(path)
	if mountEntry == nil {
		b.Backend.Logger().Error("tune failed: no mount entry found", "path", path)
		return handleError(fmt.Errorf("tune of path %q failed: no mount entry found", path))
	}
	if mountEntry != nil && !mountEntry.Local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return logical.ErrorResponse("cannot tune a non-local mount on a replication secondary"), nil
	}

	// Timing configuration parameters
	{
		var newDefault, newMax time.Duration
		defTTL := data.Get("default_lease_ttl").(string)
		switch defTTL {
		case "":
			newDefault = mountEntry.Config.DefaultLeaseTTL
		case "system":
			newDefault = time.Duration(0)
		default:
			tmpDef, err := parseutil.ParseDurationSecond(defTTL)
			if err != nil {
				return handleError(err)
			}
			newDefault = tmpDef
		}

		maxTTL := data.Get("max_lease_ttl").(string)
		switch maxTTL {
		case "":
			newMax = mountEntry.Config.MaxLeaseTTL
		case "system":
			newMax = time.Duration(0)
		default:
			tmpMax, err := parseutil.ParseDurationSecond(maxTTL)
			if err != nil {
				return handleError(err)
			}
			newMax = tmpMax
		}

		if newDefault != mountEntry.Config.DefaultLeaseTTL ||
			newMax != mountEntry.Config.MaxLeaseTTL {

			if err := b.tuneMountTTLs(ctx, path, mountEntry, newDefault, newMax); err != nil {
				b.Backend.Logger().Error("tuning failed", "path", path, "error", err)
				return handleError(err)
			}
		}
	}

	description := data.Get("description").(string)
	if description != "" {
		oldDesc := mountEntry.Description
		mountEntry.Description = description

		// Update the mount table
		var err error
		switch {
		case strings.HasPrefix(path, "auth/"):
			err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
		default:
			err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
		}
		if err != nil {
			mountEntry.Description = oldDesc
			return handleError(err)
		}
		if b.Core.logger.IsInfo() {
			b.Core.logger.Info("mount tuning of description successful", "path", path)
		}
	}

	if rawVal, ok := data.GetOk("audit_non_hmac_request_keys"); ok {
		auditNonHMACRequestKeys := rawVal.([]string)

		oldVal := mountEntry.Config.AuditNonHMACRequestKeys
		mountEntry.Config.AuditNonHMACRequestKeys = auditNonHMACRequestKeys

		// Update the mount table
		var err error
		switch {
		case strings.HasPrefix(path, "auth/"):
			err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
		default:
			err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
		}
		if err != nil {
			mountEntry.Config.AuditNonHMACRequestKeys = oldVal
			return handleError(err)
		}

		mountEntry.SyncCache()

		if b.Core.logger.IsInfo() {
			b.Core.logger.Info("mount tuning of audit_non_hmac_request_keys successful", "path", path)
		}
	}

	if rawVal, ok := data.GetOk("audit_non_hmac_response_keys"); ok {
		auditNonHMACResponseKeys := rawVal.([]string)

		oldVal := mountEntry.Config.AuditNonHMACResponseKeys
		mountEntry.Config.AuditNonHMACResponseKeys = auditNonHMACResponseKeys

		// Update the mount table
		var err error
		switch {
		case strings.HasPrefix(path, "auth/"):
			err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
		default:
			err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
		}
		if err != nil {
			mountEntry.Config.AuditNonHMACResponseKeys = oldVal
			return handleError(err)
		}

		mountEntry.SyncCache()

		if b.Core.logger.IsInfo() {
			b.Core.logger.Info("mount tuning of audit_non_hmac_response_keys successful", "path", path)
		}
	}

	if rawVal, ok := data.GetOk("listing_visibility"); ok {
		lvString := rawVal.(string)
		listingVisibility := ListingVisibilityType(lvString)

		if err := checkListingVisibility(listingVisibility); err != nil {
			return logical.ErrorResponse(fmt.Sprintf("invalid listing_visibility %s", listingVisibility)), nil
		}

		oldVal := mountEntry.Config.ListingVisibility
		mountEntry.Config.ListingVisibility = listingVisibility

		// Update the mount table
		var err error
		switch {
		case strings.HasPrefix(path, "auth/"):
			err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
		default:
			err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
		}
		if err != nil {
			mountEntry.Config.ListingVisibility = oldVal
			return handleError(err)
		}

		if b.Core.logger.IsInfo() {
			b.Core.logger.Info("mount tuning of listing_visibility successful", "path", path)
		}
	}

	if rawVal, ok := data.GetOk("passthrough_request_headers"); ok {
		headers := rawVal.([]string)

		oldVal := mountEntry.Config.PassthroughRequestHeaders
		mountEntry.Config.PassthroughRequestHeaders = headers

		// Update the mount table
		var err error
		switch {
		case strings.HasPrefix(path, "auth/"):
			err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
		default:
			err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
		}
		if err != nil {
			mountEntry.Config.PassthroughRequestHeaders = oldVal
			return handleError(err)
		}

		mountEntry.SyncCache()

		if b.Core.logger.IsInfo() {
			b.Core.logger.Info("mount tuning of passthrough_request_headers successful", "path", path)
		}
	}

	var err error
	var resp *logical.Response
	var options map[string]string
	if optionsRaw, ok := data.GetOk("options"); ok {
		options = optionsRaw.(map[string]string)
	}
	if len(options) > 0 {
		b.Core.logger.Info("mount tuning of options", "path", path, "options", options)

		var changed bool
		var numBuiltIn int
		if v, ok := options["version"]; ok {
			changed = true
			numBuiltIn++
			// Special case to make sure we can not disable versioning once it's
			// enabeled. If the vkv backend suports downgrading this can be removed.
			meVersion, err := parseutil.ParseInt(mountEntry.Options["version"])
			if err != nil {
				return nil, errwrap.Wrapf("unable to parse mount entry: {{err}}", err)
			}
			optVersion, err := parseutil.ParseInt(v)
			if err != nil {
				return handleError(errwrap.Wrapf("unable to parse options: {{err}}", err))
			}
			if meVersion > optVersion {
				return logical.ErrorResponse(fmt.Sprintf("cannot downgrade mount from version %d", meVersion)), logical.ErrInvalidRequest
			}
			if meVersion < optVersion {
				resp = &logical.Response{}
				resp.AddWarning(fmt.Sprintf("Upgrading mount from version %d to version %d. This mount will be unavailable for a brief period and will resume service shortly.", meVersion, optVersion))
			}
		}
		if options != nil {
			// For anything we don't recognize and provide special handling,
			// always write
			if len(options) > numBuiltIn {
				changed = true
			}
		}

		if changed {
			oldVal := mountEntry.Options
			mountEntry.Options = options
			// Update the mount table
			switch {
			case strings.HasPrefix(path, "auth/"):
				err = b.Core.persistAuth(ctx, b.Core.auth, &mountEntry.Local)
			default:
				err = b.Core.persistMounts(ctx, b.Core.mounts, &mountEntry.Local)
			}
			if err != nil {
				mountEntry.Options = oldVal
				return handleError(err)
			}

			// Reload the backend to kick off the upgrade process.
			b.Core.reloadBackendCommon(ctx, mountEntry, strings.HasPrefix(path, credentialRoutePrefix))
		}
	}

	return resp, nil
}

// handleLease is use to view the metadata for a given LeaseID
func (b *SystemBackend) handleLeaseLookup(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	leaseID := data.Get("lease_id").(string)
	if leaseID == "" {
		return logical.ErrorResponse("lease_id must be specified"),
			logical.ErrInvalidRequest
	}

	leaseTimes, err := b.Core.expiration.FetchLeaseTimes(leaseID)
	if err != nil {
		b.Backend.Logger().Error("error retrieving lease", "lease_id", leaseID, "error", err)
		return handleError(err)
	}
	if leaseTimes == nil {
		return logical.ErrorResponse("invalid lease"), logical.ErrInvalidRequest
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"id":           leaseID,
			"issue_time":   leaseTimes.IssueTime,
			"expire_time":  nil,
			"last_renewal": nil,
			"ttl":          int64(0),
		},
	}
	renewable, _ := leaseTimes.renewable()
	resp.Data["renewable"] = renewable

	if !leaseTimes.LastRenewalTime.IsZero() {
		resp.Data["last_renewal"] = leaseTimes.LastRenewalTime
	}
	if !leaseTimes.ExpireTime.IsZero() {
		resp.Data["expire_time"] = leaseTimes.ExpireTime
		resp.Data["ttl"] = leaseTimes.ttl()
	}
	return resp, nil
}

func (b *SystemBackend) handleLeaseLookupList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	prefix := data.Get("prefix").(string)
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	keys, err := b.Core.expiration.idView.List(ctx, prefix)
	if err != nil {
		b.Backend.Logger().Error("error listing leases", "prefix", prefix, "error", err)
		return handleErrorNoReadOnlyForward(err)
	}
	return logical.ListResponse(keys), nil
}

// handleRenew is used to renew a lease with a given LeaseID
func (b *SystemBackend) handleRenew(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	leaseID := data.Get("lease_id").(string)
	if leaseID == "" {
		leaseID = data.Get("url_lease_id").(string)
	}
	if leaseID == "" {
		return logical.ErrorResponse("lease_id must be specified"),
			logical.ErrInvalidRequest
	}
	incrementRaw := data.Get("increment").(int)

	// Convert the increment
	increment := time.Duration(incrementRaw) * time.Second

	// Invoke the expiration manager directly
	resp, err := b.Core.expiration.Renew(leaseID, increment)
	if err != nil {
		b.Backend.Logger().Error("lease renewal failed", "lease_id", leaseID, "error", err)
		return handleErrorNoReadOnlyForward(err)
	}
	return resp, err
}

// handleRevoke is used to revoke a given LeaseID
func (b *SystemBackend) handleRevoke(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	leaseID := data.Get("lease_id").(string)
	if leaseID == "" {
		leaseID = data.Get("url_lease_id").(string)
	}
	if leaseID == "" {
		return logical.ErrorResponse("lease_id must be specified"),
			logical.ErrInvalidRequest
	}

	// Invoke the expiration manager directly
	if err := b.Core.expiration.Revoke(leaseID); err != nil {
		b.Backend.Logger().Error("lease revocation failed", "lease_id", leaseID, "error", err)
		return handleErrorNoReadOnlyForward(err)
	}
	return nil, nil
}

// handleRevokePrefix is used to revoke a prefix with many LeaseIDs
func (b *SystemBackend) handleRevokePrefix(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRevokePrefixCommon(req, data, false)
}

// handleRevokeForce is used to revoke a prefix with many LeaseIDs, ignoring errors
func (b *SystemBackend) handleRevokeForce(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRevokePrefixCommon(req, data, true)
}

// handleRevokePrefixCommon is used to revoke a prefix with many LeaseIDs
func (b *SystemBackend) handleRevokePrefixCommon(
	req *logical.Request, data *framework.FieldData, force bool) (*logical.Response, error) {
	// Get all the options
	prefix := data.Get("prefix").(string)

	// Invoke the expiration manager directly
	var err error
	if force {
		err = b.Core.expiration.RevokeForce(prefix)
	} else {
		err = b.Core.expiration.RevokePrefix(prefix)
	}
	if err != nil {
		b.Backend.Logger().Error("revoke prefix failed", "prefix", prefix, "error", err)
		return handleErrorNoReadOnlyForward(err)
	}
	return nil, nil
}

// handleAuthTable handles the "auth" endpoint to provide the auth table
func (b *SystemBackend) handleAuthTable(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Core.authLock.RLock()
	defer b.Core.authLock.RUnlock()

	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}
	for _, entry := range b.Core.auth.Entries {
		info := map[string]interface{}{
			"type":        entry.Type,
			"description": entry.Description,
			"accessor":    entry.Accessor,
			"local":       entry.Local,
			"seal_wrap":   entry.SealWrap,
			"options":     entry.Options,
		}
		entryConfig := map[string]interface{}{
			"default_lease_ttl": int64(entry.Config.DefaultLeaseTTL.Seconds()),
			"max_lease_ttl":     int64(entry.Config.MaxLeaseTTL.Seconds()),
			"plugin_name":       entry.Config.PluginName,
		}
		if rawVal, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_request_keys"); ok {
			entryConfig["audit_non_hmac_request_keys"] = rawVal.([]string)
		}
		if rawVal, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_response_keys"); ok {
			entryConfig["audit_non_hmac_response_keys"] = rawVal.([]string)
		}
		// Even though empty value is valid for ListingVisibility, we can ignore
		// this case during mount since there's nothing to unset/hide.
		if len(entry.Config.ListingVisibility) > 0 {
			entryConfig["listing_visibility"] = entry.Config.ListingVisibility
		}
		if rawVal, ok := entry.synthesizedConfigCache.Load("passthrough_request_headers"); ok {
			entryConfig["passthrough_request_headers"] = rawVal.([]string)
		}

		info["config"] = entryConfig
		resp.Data[entry.Path] = info
	}
	return resp, nil
}

// handleEnableAuth is used to enable a new credential backend
func (b *SystemBackend) handleEnableAuth(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	repState := b.Core.ReplicationState()
	local := data.Get("local").(bool)
	if !local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return logical.ErrorResponse("cannot add a non-local mount to a replication secondary"), nil
	}

	// Get all the options
	path := data.Get("path").(string)
	path = sanitizeMountPath(path)
	logicalType := data.Get("type").(string)
	description := data.Get("description").(string)
	pluginName := data.Get("plugin_name").(string)
	sealWrap := data.Get("seal_wrap").(bool)
	options := data.Get("options").(map[string]string)

	var config MountConfig
	var apiConfig APIMountConfig

	configMap := data.Get("config").(map[string]interface{})
	if configMap != nil && len(configMap) != 0 {
		err := mapstructure.Decode(configMap, &apiConfig)
		if err != nil {
			return logical.ErrorResponse(
					"unable to convert given auth config information"),
				logical.ErrInvalidRequest
		}
	}

	switch apiConfig.DefaultLeaseTTL {
	case "":
	case "system":
	default:
		tmpDef, err := parseutil.ParseDurationSecond(apiConfig.DefaultLeaseTTL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
					"unable to parse default TTL of %s: %s", apiConfig.DefaultLeaseTTL, err)),
				logical.ErrInvalidRequest
		}
		config.DefaultLeaseTTL = tmpDef
	}

	switch apiConfig.MaxLeaseTTL {
	case "":
	case "system":
	default:
		tmpMax, err := parseutil.ParseDurationSecond(apiConfig.MaxLeaseTTL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
					"unable to parse max TTL of %s: %s", apiConfig.MaxLeaseTTL, err)),
				logical.ErrInvalidRequest
		}
		config.MaxLeaseTTL = tmpMax
	}

	if config.MaxLeaseTTL != 0 && config.DefaultLeaseTTL > config.MaxLeaseTTL {
		return logical.ErrorResponse(
				"given default lease TTL greater than given max lease TTL"),
			logical.ErrInvalidRequest
	}

	if config.DefaultLeaseTTL > b.Core.maxLeaseTTL && config.MaxLeaseTTL == 0 {
		return logical.ErrorResponse(fmt.Sprintf(
				"given default lease TTL greater than system max lease TTL of %d", int(b.Core.maxLeaseTTL.Seconds()))),
			logical.ErrInvalidRequest
	}

	switch logicalType {
	case "":
		return logical.ErrorResponse(
				"backend type must be specified as a string"),
			logical.ErrInvalidRequest

	case "plugin":
		// Only set plugin name if mount is of type plugin, with apiConfig.PluginName
		// option taking precedence.
		switch {
		case apiConfig.PluginName != "":
			config.PluginName = apiConfig.PluginName
		case pluginName != "":
			config.PluginName = pluginName
		default:
			return logical.ErrorResponse(
					"plugin_name must be provided for plugin backend"),
				logical.ErrInvalidRequest
		}
	}

	if options != nil && options["version"] != "" {
		return logical.ErrorResponse(fmt.Sprintf(
				"auth method %q does not allow setting a version", logicalType)),
			logical.ErrInvalidRequest
	}

	if err := checkListingVisibility(apiConfig.ListingVisibility); err != nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid listing_visibility %s", apiConfig.ListingVisibility)), nil
	}
	config.ListingVisibility = apiConfig.ListingVisibility

	if len(apiConfig.AuditNonHMACRequestKeys) > 0 {
		config.AuditNonHMACRequestKeys = apiConfig.AuditNonHMACRequestKeys
	}
	if len(apiConfig.AuditNonHMACResponseKeys) > 0 {
		config.AuditNonHMACResponseKeys = apiConfig.AuditNonHMACResponseKeys
	}
	if len(apiConfig.PassthroughRequestHeaders) > 0 {
		config.PassthroughRequestHeaders = apiConfig.PassthroughRequestHeaders
	}

	// Create the mount entry
	me := &MountEntry{
		Table:       credentialTableType,
		Path:        path,
		Type:        logicalType,
		Description: description,
		Config:      config,
		Local:       local,
		SealWrap:    sealWrap,
		Options:     options,
	}

	// Attempt enabling
	if err := b.Core.enableCredential(ctx, me); err != nil {
		b.Backend.Logger().Error("enable auth mount failed", "path", me.Path, "error", err)
		return handleError(err)
	}
	return nil, nil
}

// handleDisableAuth is used to disable a credential backend
func (b *SystemBackend) handleDisableAuth(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	path = sanitizeMountPath(path)

	fullPath := credentialRoutePrefix + path

	repState := b.Core.ReplicationState()
	entry := b.Core.router.MatchingMountEntry(fullPath)
	if entry != nil && !entry.Local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return logical.ErrorResponse("cannot unmount a non-local mount on a replication secondary"), nil
	}

	// We return success when the mount does not exists to not expose if the
	// mount existed or not
	match := b.Core.router.MatchingMount(fullPath)
	if match == "" || fullPath != match {
		return nil, nil
	}

	// Attempt disable
	if err := b.Core.disableCredential(ctx, path); err != nil {
		b.Backend.Logger().Error("disable auth mount failed", "path", path, "error", err)
		return handleError(err)
	}
	return nil, nil
}

// handlePoliciesList handles /sys/policy/ and /sys/policies/<type> endpoints to provide the enabled policies
func (b *SystemBackend) handlePoliciesList(policyType PolicyType) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		policies, err := b.Core.policyStore.ListPolicies(ctx, policyType)
		if err != nil {
			return nil, err
		}

		switch policyType {
		case PolicyTypeACL:
			// Add the special "root" policy if not egp
			policies = append(policies, "root")
			resp := logical.ListResponse(policies)

			// If the request is from sys/policy/ we handle backwards compatibility
			if strings.HasPrefix(req.Path, "policy") {
				resp.Data["policies"] = resp.Data["keys"]
			}

			return resp, nil
		}

		return logical.ErrorResponse("unknown policy type"), nil
	}
}

// handlePoliciesRead handles the "/sys/policy/<name>" and "/sys/policies/<type>/<name>" endpoints to read a policy
func (b *SystemBackend) handlePoliciesRead(policyType PolicyType) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)

		policy, err := b.Core.policyStore.GetPolicy(ctx, name, policyType)
		if err != nil {
			return handleError(err)
		}

		if policy == nil {
			return nil, nil
		}

		// If the request is from sys/policy/ we handle backwards compatibility
		var respDataPolicyName string
		if policyType == PolicyTypeACL && strings.HasPrefix(req.Path, "policy") {
			respDataPolicyName = "rules"
		} else {
			respDataPolicyName = "policy"
		}

		resp := &logical.Response{
			Data: map[string]interface{}{
				"name":             policy.Name,
				respDataPolicyName: policy.Raw,
			},
		}

		return resp, nil
	}
}

// handlePoliciesSet handles the "/sys/policy/<name>" and "/sys/policies/<type>/<name>" endpoints to set a policy
func (b *SystemBackend) handlePoliciesSet(policyType PolicyType) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		var resp *logical.Response

		policy := &Policy{
			Name: strings.ToLower(data.Get("name").(string)),
			Type: policyType,
		}
		if policy.Name == "" {
			return logical.ErrorResponse("policy name must be provided in the URL"), nil
		}

		policy.Raw = data.Get("policy").(string)
		if policy.Raw == "" && policyType == PolicyTypeACL && strings.HasPrefix(req.Path, "policy") {
			policy.Raw = data.Get("rules").(string)
			if resp == nil {
				resp = &logical.Response{}
			}
			resp.AddWarning("'rules' is deprecated, please use 'policy' instead")
		}
		if policy.Raw == "" {
			return logical.ErrorResponse("'policy' parameter not supplied or empty"), nil
		}

		if polBytes, err := base64.StdEncoding.DecodeString(policy.Raw); err == nil {
			policy.Raw = string(polBytes)
		}

		switch policyType {
		case PolicyTypeACL:
			p, err := ParseACLPolicy(policy.Raw)
			if err != nil {
				return handleError(err)
			}
			policy.Paths = p.Paths

		default:
			return logical.ErrorResponse("unknown policy type"), nil
		}

		// Update the policy
		if err := b.Core.policyStore.SetPolicy(ctx, policy); err != nil {
			return handleError(err)
		}
		return resp, nil
	}
}

func (b *SystemBackend) handlePoliciesDelete(policyType PolicyType) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)

		if err := b.Core.policyStore.DeletePolicy(ctx, name, policyType); err != nil {
			return handleError(err)
		}
		return nil, nil
	}
}

// handleAuditTable handles the "audit" endpoint to provide the audit table
func (b *SystemBackend) handleAuditTable(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Core.auditLock.RLock()
	defer b.Core.auditLock.RUnlock()

	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}
	for _, entry := range b.Core.audit.Entries {
		info := map[string]interface{}{
			"path":        entry.Path,
			"type":        entry.Type,
			"description": entry.Description,
			"options":     entry.Options,
			"local":       entry.Local,
		}
		resp.Data[entry.Path] = info
	}
	return resp, nil
}

// handleAuditHash is used to fetch the hash of the given input data with the
// specified audit backend's salt
func (b *SystemBackend) handleAuditHash(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	input := data.Get("input").(string)
	if input == "" {
		return logical.ErrorResponse("the \"input\" parameter is empty"), nil
	}

	path = sanitizeMountPath(path)

	hash, err := b.Core.auditBroker.GetHash(ctx, path, input)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"hash": hash,
		},
	}, nil
}

// handleEnableAudit is used to enable a new audit backend
func (b *SystemBackend) handleEnableAudit(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	repState := b.Core.ReplicationState()

	local := data.Get("local").(bool)
	if !local && repState.HasState(consts.ReplicationPerformanceSecondary) {
		return logical.ErrorResponse("cannot add a non-local mount to a replication secondary"), nil
	}

	// Get all the options
	path := data.Get("path").(string)
	backendType := data.Get("type").(string)
	description := data.Get("description").(string)
	options := data.Get("options").(map[string]string)

	// Create the mount entry
	me := &MountEntry{
		Table:       auditTableType,
		Path:        path,
		Type:        backendType,
		Description: description,
		Options:     options,
		Local:       local,
	}

	// Attempt enabling
	if err := b.Core.enableAudit(ctx, me); err != nil {
		b.Backend.Logger().Error("enable audit mount failed", "path", me.Path, "error", err)
		return handleError(err)
	}
	return nil, nil
}

// handleDisableAudit is used to disable an audit backend
func (b *SystemBackend) handleDisableAudit(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)

	// Attempt disable
	if existed, err := b.Core.disableAudit(ctx, path); existed && err != nil {
		b.Backend.Logger().Error("disable audit mount failed", "path", path, "error", err)
		return handleError(err)
	}
	return nil, nil
}

func (b *SystemBackend) handleConfigUIHeadersRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	header := data.Get("header").(string)

	value, err := b.Core.uiConfig.GetHeader(ctx, header)
	if err != nil {
		return nil, err
	}
	if value == "" {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"value": value,
		},
	}, nil
}

func (b *SystemBackend) handleConfigUIHeadersList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	headers, err := b.Core.uiConfig.HeaderKeys(ctx)
	if err != nil {
		return nil, err
	}
	if len(headers) == 0 {
		return nil, nil
	}

	return logical.ListResponse(headers), nil
}

func (b *SystemBackend) handleConfigUIHeadersUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	header := data.Get("header").(string)
	values := data.Get("values").([]string)
	if header == "" || len(values) == 0 {
		return logical.ErrorResponse("header and values must be specified"), logical.ErrInvalidRequest
	}

	lowerHeader := strings.ToLower(header)
	if strings.HasPrefix(lowerHeader, "x-vault-") {
		return logical.ErrorResponse("X-Vault headers cannot be set"), logical.ErrInvalidRequest
	}

	// Translate the list of values to the valid header string
	value := http.Header{}
	for _, v := range values {
		value.Add(header, v)
	}
	err := b.Core.uiConfig.SetHeader(ctx, header, value.Get(header))
	if err != nil {
		return nil, err
	}

	// Warn when overriding the CSP
	resp := &logical.Response{}
	if lowerHeader == "content-security-policy" {
		resp.AddWarning("overriding default Content-Security-Policy which is secure by default, proceed with caution")
	}

	return resp, nil
}

func (b *SystemBackend) handleConfigUIHeadersDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	header := data.Get("header").(string)
	err := b.Core.uiConfig.DeleteHeader(ctx, header)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// handleRawRead is used to read directly from the barrier
func (b *SystemBackend) handleRawRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)

	// Prevent access of protected paths
	for _, p := range protectedPaths {
		if strings.HasPrefix(path, p) {
			err := fmt.Sprintf("cannot read '%s'", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	entry, err := b.Core.barrier.Get(ctx, path)
	if err != nil {
		return handleErrorNoReadOnlyForward(err)
	}
	if entry == nil {
		return nil, nil
	}

	// Run this through the decompression helper to see if it's been compressed.
	// If the input contained the compression canary, `outputBytes` will hold
	// the decompressed data. If the input was not compressed, then `outputBytes`
	// will be nil.
	outputBytes, _, err := compressutil.Decompress(entry.Value)
	if err != nil {
		return handleErrorNoReadOnlyForward(err)
	}

	// `outputBytes` is nil if the input is uncompressed. In that case set it to the original input.
	if outputBytes == nil {
		outputBytes = entry.Value
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"value": string(outputBytes),
		},
	}
	return resp, nil
}

// handleRawWrite is used to write directly to the barrier
func (b *SystemBackend) handleRawWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)

	// Prevent access of protected paths
	for _, p := range protectedPaths {
		if strings.HasPrefix(path, p) {
			err := fmt.Sprintf("cannot write '%s'", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	value := data.Get("value").(string)
	entry := &Entry{
		Key:   path,
		Value: []byte(value),
	}
	if err := b.Core.barrier.Put(ctx, entry); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleRawDelete is used to delete directly from the barrier
func (b *SystemBackend) handleRawDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)

	// Prevent access of protected paths
	for _, p := range protectedPaths {
		if strings.HasPrefix(path, p) {
			err := fmt.Sprintf("cannot delete '%s'", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	if err := b.Core.barrier.Delete(ctx, path); err != nil {
		return handleErrorNoReadOnlyForward(err)
	}
	return nil, nil
}

// handleRawList is used to list directly from the barrier
func (b *SystemBackend) handleRawList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path != "" && !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// Prevent access of protected paths
	for _, p := range protectedPaths {
		if strings.HasPrefix(path, p) {
			err := fmt.Sprintf("cannot list '%s'", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	keys, err := b.Core.barrier.List(ctx, path)
	if err != nil {
		return handleErrorNoReadOnlyForward(err)
	}
	return logical.ListResponse(keys), nil
}

// handleKeyStatus returns status information about the backend key
func (b *SystemBackend) handleKeyStatus(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get the key info
	info, err := b.Core.barrier.ActiveKeyInfo()
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"term":         info.Term,
			"install_time": info.InstallTime.Format(time.RFC3339Nano),
		},
	}
	return resp, nil
}

// handleRotate is used to trigger a key rotation
func (b *SystemBackend) handleRotate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	repState := b.Core.ReplicationState()
	if repState.HasState(consts.ReplicationPerformanceSecondary) {
		return logical.ErrorResponse("cannot rotate on a replication secondary"), nil
	}

	// Rotate to the new term
	newTerm, err := b.Core.barrier.Rotate(ctx)
	if err != nil {
		b.Backend.Logger().Error("failed to create new encryption key", "error", err)
		return handleError(err)
	}
	b.Backend.Logger().Info("installed new encryption key")

	// In HA mode, we need to an upgrade path for the standby instances
	if b.Core.ha != nil {
		// Create the upgrade path to the new term
		if err := b.Core.barrier.CreateUpgrade(ctx, newTerm); err != nil {
			b.Backend.Logger().Error("failed to create new upgrade", "term", newTerm, "error", err)
		}

		// Schedule the destroy of the upgrade path
		time.AfterFunc(keyRotateGracePeriod, func() {
			if err := b.Core.barrier.DestroyUpgrade(ctx, newTerm); err != nil {
				b.Backend.Logger().Error("failed to destroy upgrade", "term", newTerm, "error", err)
			}
		})
	}

	// Write to the canary path, which will force a synchronous truing during
	// replication
	if err := b.Core.barrier.Put(ctx, &Entry{
		Key:   coreKeyringCanaryPath,
		Value: []byte(fmt.Sprintf("new-rotation-term-%d", newTerm)),
	}); err != nil {
		b.Core.logger.Error("error saving keyring canary", "error", err)
		return nil, errwrap.Wrapf("failed to save keyring canary: {{err}}", err)
	}

	return nil, nil
}

func (b *SystemBackend) handleWrappingPubkey(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	x, _ := b.Core.wrappingJWTKey.X.MarshalText()
	y, _ := b.Core.wrappingJWTKey.Y.MarshalText()
	return &logical.Response{
		Data: map[string]interface{}{
			"jwt_x":     string(x),
			"jwt_y":     string(y),
			"jwt_curve": corePrivateKeyTypeP521,
		},
	}, nil
}

func (b *SystemBackend) handleWrappingWrap(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.WrapInfo == nil || req.WrapInfo.TTL == 0 {
		return logical.ErrorResponse("endpoint requires response wrapping to be used"), logical.ErrInvalidRequest
	}

	// N.B.: Do *NOT* allow JWT wrapping tokens to be created through this
	// endpoint. JWTs are signed so if we don't allow users to create wrapping
	// tokens using them we can ensure that an operator can't spoof a legit JWT
	// wrapped token, which makes certain init/rekey/generate-root cases have
	// better properties.
	req.WrapInfo.Format = "uuid"

	return &logical.Response{
		Data: data.Raw,
	}, nil
}

func (b *SystemBackend) handleWrappingUnwrap(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// If a third party is unwrapping (rather than the calling token being the
	// wrapping token) we detect this so that we can revoke the original
	// wrapping token after reading it
	var thirdParty bool

	token := data.Get("token").(string)
	if token != "" {
		thirdParty = true
	} else {
		token = req.ClientToken
	}

	// Get the policies so we can determine if this is a normal response
	// wrapping request or a control group token.
	//
	// We use lookupTainted here because the token might have already been used
	// by handleRequest(), this happens when it's a normal response wrapping
	// request and the token was provided "first party". We want to inspect the
	// token policies but will not use this token entry for anything else.
	te, err := b.Core.tokenStore.lookupTainted(ctx, token)
	if err != nil {
		return nil, err
	}
	if te == nil {
		return nil, errors.New("could not find token")
	}
	if len(te.Policies) != 1 {
		return nil, errors.New("token is not a valid unwrap token")
	}

	var response string
	switch te.Policies[0] {
	case responseWrappingPolicyName:
		response, err = b.responseWrappingUnwrap(ctx, token, thirdParty)
	}
	if err != nil {
		var respErr *logical.Response
		if len(response) > 0 {
			respErr = logical.ErrorResponse(response)
		}

		return respErr, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{},
	}

	// Most of the time we want to just send over the marshalled HTTP bytes.
	// However there is a sad separate case: if the original response was using
	// bare values we need to use those or else what comes back is garbled.
	httpResp := &logical.HTTPResponse{}
	err = jsonutil.DecodeJSON([]byte(response), httpResp)
	if err != nil {
		return nil, errwrap.Wrapf("error decoding wrapped response: {{err}}", err)
	}
	if httpResp.Data != nil &&
		(httpResp.Data[logical.HTTPStatusCode] != nil ||
			httpResp.Data[logical.HTTPRawBody] != nil ||
			httpResp.Data[logical.HTTPContentType] != nil) {
		if httpResp.Data[logical.HTTPStatusCode] != nil {
			resp.Data[logical.HTTPStatusCode] = httpResp.Data[logical.HTTPStatusCode]
		}
		if httpResp.Data[logical.HTTPContentType] != nil {
			resp.Data[logical.HTTPContentType] = httpResp.Data[logical.HTTPContentType]
		}

		rawBody := httpResp.Data[logical.HTTPRawBody]
		if rawBody != nil {
			// Decode here so that we can audit properly
			switch rawBody.(type) {
			case string:
				// Best effort decoding; if this works, the original value was
				// probably a []byte instead of a string, but was marshaled
				// when the value was saved, so this restores it as it was
				decBytes, err := base64.StdEncoding.DecodeString(rawBody.(string))
				if err == nil {
					// We end up with []byte, will not be HMAC'd
					resp.Data[logical.HTTPRawBody] = decBytes
				} else {
					// We end up with string, will be HMAC'd
					resp.Data[logical.HTTPRawBody] = rawBody
				}
			default:
				b.Core.Logger().Error("unexpected type of raw body when decoding wrapped token", "type", fmt.Sprintf("%T", rawBody))
			}

			resp.Data[logical.HTTPRawBodyAlreadyJSONDecoded] = true
		}

		return resp, nil
	}

	if len(response) == 0 {
		resp.Data[logical.HTTPStatusCode] = 204
	} else {
		resp.Data[logical.HTTPStatusCode] = 200
		resp.Data[logical.HTTPRawBody] = []byte(response)
		resp.Data[logical.HTTPContentType] = "application/json"
	}

	return resp, nil
}

// responseWrappingUnwrap will read the stored response in the cubbyhole and
// return the raw HTTP response.
func (b *SystemBackend) responseWrappingUnwrap(ctx context.Context, token string, thirdParty bool) (string, error) {
	if thirdParty {
		// Use the token to decrement the use count to avoid a second operation on the token.
		_, err := b.Core.tokenStore.UseTokenByID(ctx, token)
		if err != nil {
			return "", errwrap.Wrapf("error decrementing wrapping token's use-count: {{err}}", err)
		}

		defer b.Core.tokenStore.revokeOrphan(ctx, token)
	}

	cubbyReq := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "cubbyhole/response",
		ClientToken: token,
	}
	cubbyResp, err := b.Core.router.Route(ctx, cubbyReq)
	if err != nil {
		return "", errwrap.Wrapf("error looking up wrapping information: {{err}}", err)
	}
	if cubbyResp == nil {
		return "no information found; wrapping token may be from a previous Vault version", ErrInternalError
	}
	if cubbyResp != nil && cubbyResp.IsError() {
		return cubbyResp.Error().Error(), nil
	}
	if cubbyResp.Data == nil {
		return "wrapping information was nil; wrapping token may be from a previous Vault version", ErrInternalError
	}

	responseRaw := cubbyResp.Data["response"]
	if responseRaw == nil {
		return "", fmt.Errorf("no response found inside the cubbyhole")
	}
	response, ok := responseRaw.(string)
	if !ok {
		return "", fmt.Errorf("could not decode response inside the cubbyhole")
	}

	return response, nil
}

func (b *SystemBackend) handleWrappingLookup(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// This ordering of lookups has been validated already in the wrapping
	// validation func, we're just doing this for a safety check
	token := data.Get("token").(string)
	if token == "" {
		token = req.ClientToken
		if token == "" {
			return logical.ErrorResponse("missing \"token\" value in input"), logical.ErrInvalidRequest
		}
	}

	cubbyReq := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "cubbyhole/wrapinfo",
		ClientToken: token,
	}
	cubbyResp, err := b.Core.router.Route(ctx, cubbyReq)
	if err != nil {
		return nil, errwrap.Wrapf("error looking up wrapping information: {{err}}", err)
	}
	if cubbyResp == nil {
		return logical.ErrorResponse("no information found; wrapping token may be from a previous Vault version"), nil
	}
	if cubbyResp != nil && cubbyResp.IsError() {
		return cubbyResp, nil
	}
	if cubbyResp.Data == nil {
		return logical.ErrorResponse("wrapping information was nil; wrapping token may be from a previous Vault version"), nil
	}

	creationTTLRaw := cubbyResp.Data["creation_ttl"]
	creationTime := cubbyResp.Data["creation_time"]
	creationPath := cubbyResp.Data["creation_path"]

	resp := &logical.Response{
		Data: map[string]interface{}{},
	}
	if creationTTLRaw != nil {
		creationTTL, err := creationTTLRaw.(json.Number).Int64()
		if err != nil {
			return nil, errwrap.Wrapf("error reading creation_ttl value from wrapping information: {{err}}", err)
		}
		resp.Data["creation_ttl"] = time.Duration(creationTTL).Seconds()
	}
	if creationTime != nil {
		// This was JSON marshaled so it's already a string in RFC3339 format
		resp.Data["creation_time"] = cubbyResp.Data["creation_time"]
	}
	if creationPath != nil {
		resp.Data["creation_path"] = cubbyResp.Data["creation_path"]
	}

	return resp, nil
}

func (b *SystemBackend) handleWrappingRewrap(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// If a third party is rewrapping (rather than the calling token being the
	// wrapping token) we detect this so that we can revoke the original
	// wrapping token after reading it. Right now wrapped tokens can't unwrap
	// themselves, but in case we change it, this will be ready to do the right
	// thing.
	var thirdParty bool

	token := data.Get("token").(string)
	if token != "" {
		thirdParty = true
	} else {
		token = req.ClientToken
	}

	if thirdParty {
		// Use the token to decrement the use count to avoid a second operation on the token.
		_, err := b.Core.tokenStore.UseTokenByID(ctx, token)
		if err != nil {
			return nil, errwrap.Wrapf("error decrementing wrapping token's use-count: {{err}}", err)
		}
		defer b.Core.tokenStore.revokeOrphan(ctx, token)
	}

	// Fetch the original TTL
	cubbyReq := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "cubbyhole/wrapinfo",
		ClientToken: token,
	}
	cubbyResp, err := b.Core.router.Route(ctx, cubbyReq)
	if err != nil {
		return nil, errwrap.Wrapf("error looking up wrapping information: {{err}}", err)
	}
	if cubbyResp == nil {
		return logical.ErrorResponse("no information found; wrapping token may be from a previous Vault version"), nil
	}
	if cubbyResp != nil && cubbyResp.IsError() {
		return cubbyResp, nil
	}
	if cubbyResp.Data == nil {
		return logical.ErrorResponse("wrapping information was nil; wrapping token may be from a previous Vault version"), nil
	}

	// Set the creation TTL on the request
	creationTTLRaw := cubbyResp.Data["creation_ttl"]
	if creationTTLRaw == nil {
		return nil, fmt.Errorf("creation_ttl value in wrapping information was nil")
	}
	creationTTL, err := cubbyResp.Data["creation_ttl"].(json.Number).Int64()
	if err != nil {
		return nil, errwrap.Wrapf("error reading creation_ttl value from wrapping information: {{err}}", err)
	}

	// Get creation_path to return as the response later
	creationPathRaw := cubbyResp.Data["creation_path"]
	if creationPathRaw == nil {
		return nil, fmt.Errorf("creation_path value in wrapping information was nil")
	}
	creationPath := creationPathRaw.(string)

	// Fetch the original response and return it as the data for the new response
	cubbyReq = &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "cubbyhole/response",
		ClientToken: token,
	}
	cubbyResp, err = b.Core.router.Route(ctx, cubbyReq)
	if err != nil {
		return nil, errwrap.Wrapf("error looking up response: {{err}}", err)
	}
	if cubbyResp == nil {
		return logical.ErrorResponse("no information found; wrapping token may be from a previous Vault version"), nil
	}
	if cubbyResp != nil && cubbyResp.IsError() {
		return cubbyResp, nil
	}
	if cubbyResp.Data == nil {
		return logical.ErrorResponse("wrapping information was nil; wrapping token may be from a previous Vault version"), nil
	}

	response := cubbyResp.Data["response"]
	if response == nil {
		return nil, fmt.Errorf("no response found inside the cubbyhole")
	}

	// Return response in "response"; wrapping code will detect the rewrap and
	// slot in instead of nesting
	return &logical.Response{
		Data: map[string]interface{}{
			"response": response,
		},
		WrapInfo: &wrapping.ResponseWrapInfo{
			TTL:          time.Duration(creationTTL),
			CreationPath: creationPath,
		},
	}, nil
}

func (b *SystemBackend) pathHashWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	inputB64 := d.Get("input").(string)
	format := d.Get("format").(string)
	algorithm := d.Get("urlalgorithm").(string)
	if algorithm == "" {
		algorithm = d.Get("algorithm").(string)
	}

	input, err := base64.StdEncoding.DecodeString(inputB64)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("unable to decode input as base64: %s", err)), logical.ErrInvalidRequest
	}

	switch format {
	case "hex":
	case "base64":
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported encoding format %s; must be \"hex\" or \"base64\"", format)), nil
	}

	var hf hash.Hash
	switch algorithm {
	case "sha2-224":
		hf = sha256.New224()
	case "sha2-256":
		hf = sha256.New()
	case "sha2-384":
		hf = sha512.New384()
	case "sha2-512":
		hf = sha512.New()
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported algorithm %s", algorithm)), nil
	}
	hf.Write(input)
	retBytes := hf.Sum(nil)

	var retStr string
	switch format {
	case "hex":
		retStr = hex.EncodeToString(retBytes)
	case "base64":
		retStr = base64.StdEncoding.EncodeToString(retBytes)
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"sum": retStr,
		},
	}
	return resp, nil
}

func (b *SystemBackend) pathRandomWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	bytes := 0
	var err error
	strBytes := d.Get("urlbytes").(string)
	if strBytes != "" {
		bytes, err = strconv.Atoi(strBytes)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("error parsing url-set byte count: %s", err)), nil
		}
	} else {
		bytes = d.Get("bytes").(int)
	}
	format := d.Get("format").(string)

	if bytes < 1 {
		return logical.ErrorResponse(`"bytes" cannot be less than 1`), nil
	}

	switch format {
	case "hex":
	case "base64":
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported encoding format %s; must be \"hex\" or \"base64\"", format)), nil
	}

	randBytes, err := uuid.GenerateRandomBytes(bytes)
	if err != nil {
		return nil, err
	}

	var retStr string
	switch format {
	case "hex":
		retStr = hex.EncodeToString(randBytes)
	case "base64":
		retStr = base64.StdEncoding.EncodeToString(randBytes)
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"random_bytes": retStr,
		},
	}
	return resp, nil
}

func hasMountAccess(acl *ACL, path string) bool {
	// If an ealier policy is giving us access to the mount path then we can do
	// a fast return.
	capabilities := acl.Capabilities(path)
	if !strutil.StrListContains(capabilities, DenyCapability) {
		return true
	}

	var aclCapabilitiesGiven bool
	walkFn := func(s string, v interface{}) bool {
		if v == nil {
			return false
		}

		perms := v.(*ACLPermissions)

		switch {
		case perms.CapabilitiesBitmap&DenyCapabilityInt > 0:
			return false

		case perms.CapabilitiesBitmap&CreateCapabilityInt > 0,
			perms.CapabilitiesBitmap&DeleteCapabilityInt > 0,
			perms.CapabilitiesBitmap&ListCapabilityInt > 0,
			perms.CapabilitiesBitmap&ReadCapabilityInt > 0,
			perms.CapabilitiesBitmap&SudoCapabilityInt > 0,
			perms.CapabilitiesBitmap&UpdateCapabilityInt > 0:

			aclCapabilitiesGiven = true
			return true
		}

		return false
	}

	acl.exactRules.WalkPrefix(path, walkFn)
	if !aclCapabilitiesGiven {
		acl.globRules.WalkPrefix(path, walkFn)
	}

	return aclCapabilitiesGiven
}

func (b *SystemBackend) pathInternalUIMountsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}

	secretMounts := make(map[string]interface{})
	authMounts := make(map[string]interface{})
	resp.Data["secret"] = secretMounts
	resp.Data["auth"] = authMounts

	var acl *ACL
	var isAuthed bool
	var err error
	if req.ClientToken != "" {
		isAuthed = true

		var entity *identity.Entity
		var te *logical.TokenEntry
		// Load the ACL policies so we can walk the prefix for this mount
		acl, te, entity, _, err = b.Core.fetchACLTokenEntryAndEntity(req)
		if err != nil {
			return nil, err
		}
		if entity != nil && entity.Disabled {
			b.logger.Warn("permission denied as the entity on the token is disabled")
			return nil, logical.ErrPermissionDenied
		}
		if te != nil && te.EntityID != "" && entity == nil {
			b.logger.Warn("permission denied as the entity on the token is invalid")
			return nil, logical.ErrPermissionDenied
		}
	}

	hasAccess := func(me *MountEntry) bool {
		if me.Config.ListingVisibility == ListingVisibilityUnauth {
			return true
		}

		if isAuthed {
			return hasMountAccess(acl, me.Path)
		}

		return false
	}

	b.Core.mountsLock.RLock()
	for _, entry := range b.Core.mounts.Entries {
		if hasAccess(entry) {
			if isAuthed {
				// If this is an authed request return all the mount info
				secretMounts[entry.Path] = mountInfo(entry)
			} else {
				secretMounts[entry.Path] = map[string]interface{}{
					"type":        entry.Type,
					"description": entry.Description,
					"options":     entry.Options,
				}
			}
		}
	}
	b.Core.mountsLock.RUnlock()

	b.Core.authLock.RLock()
	for _, entry := range b.Core.auth.Entries {
		if hasAccess(entry) {
			if isAuthed {
				// If this is an authed request return all the mount info
				authMounts[entry.Path] = mountInfo(entry)
			} else {
				authMounts[entry.Path] = map[string]interface{}{
					"type":        entry.Type,
					"description": entry.Description,
					"options":     entry.Options,
				}
			}
		}
	}
	b.Core.authLock.RUnlock()

	return resp, nil
}

func (b *SystemBackend) pathInternalUIMountRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	path := d.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse("path not set"), logical.ErrInvalidRequest
	}
	path = sanitizeMountPath(path)

	errResp := logical.ErrorResponse(fmt.Sprintf("Preflight capability check returned 403, please ensure client's policies grant access to path \"%s\"", path))

	me := b.Core.router.MatchingMountEntry(path)
	if me == nil {
		// Return a permission denied error here so this path cannot be used to
		// brute force a list of mounts.
		return errResp, logical.ErrPermissionDenied
	}

	resp := &logical.Response{
		Data: mountInfo(me),
	}
	resp.Data["path"] = me.Path

	// Load the ACL policies so we can walk the prefix for this mount
	acl, te, entity, _, err := b.Core.fetchACLTokenEntryAndEntity(req)
	if err != nil {
		return nil, err
	}
	if entity != nil && entity.Disabled {
		b.logger.Warn("permission denied as the entity on the token is disabled")
		return errResp, logical.ErrPermissionDenied
	}
	if te != nil && te.EntityID != "" && entity == nil {
		b.logger.Warn("permission denied as the entity on the token is invalid")
		return nil, logical.ErrPermissionDenied
	}

	if !hasMountAccess(acl, me.Path) {
		return errResp, logical.ErrPermissionDenied
	}

	return resp, nil
}

func (b *SystemBackend) pathInternalUIResultantACL(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		// 204 -- no ACL
		return nil, nil
	}

	acl, te, entity, _, err := b.Core.fetchACLTokenEntryAndEntity(req)
	if err != nil {
		return nil, err
	}

	if entity != nil && entity.Disabled {
		b.logger.Warn("permission denied as the entity on the token is disabled")
		return logical.ErrorResponse(logical.ErrPermissionDenied.Error()), nil
	}
	if te != nil && te.EntityID != "" && entity == nil {
		b.logger.Warn("permission denied as the entity on the token is invalid")
		return logical.ErrorResponse(logical.ErrPermissionDenied.Error()), nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"root": false,
		},
	}

	if acl.root {
		resp.Data["root"] = true
		return resp, nil
	}

	exact := map[string]interface{}{}
	glob := map[string]interface{}{}

	walkFn := func(pt map[string]interface{}, s string, v interface{}) {
		if v == nil {
			return
		}

		perms := v.(*ACLPermissions)
		capabilities := []string{}

		if perms.CapabilitiesBitmap&CreateCapabilityInt > 0 {
			capabilities = append(capabilities, CreateCapability)
		}
		if perms.CapabilitiesBitmap&DeleteCapabilityInt > 0 {
			capabilities = append(capabilities, DeleteCapability)
		}
		if perms.CapabilitiesBitmap&ListCapabilityInt > 0 {
			capabilities = append(capabilities, ListCapability)
		}
		if perms.CapabilitiesBitmap&ReadCapabilityInt > 0 {
			capabilities = append(capabilities, ReadCapability)
		}
		if perms.CapabilitiesBitmap&SudoCapabilityInt > 0 {
			capabilities = append(capabilities, SudoCapability)
		}
		if perms.CapabilitiesBitmap&UpdateCapabilityInt > 0 {
			capabilities = append(capabilities, UpdateCapability)
		}

		// If "deny" is explicitly set or if the path has no capabilities at all,
		// set the path capabilities to "deny"
		if perms.CapabilitiesBitmap&DenyCapabilityInt > 0 || len(capabilities) == 0 {
			capabilities = []string{DenyCapability}
		}

		res := map[string]interface{}{}
		if len(capabilities) > 0 {
			res["capabilities"] = capabilities
		}
		if perms.MinWrappingTTL != 0 {
			res["min_wrapping_ttl"] = int64(perms.MinWrappingTTL.Seconds())
		}
		if perms.MaxWrappingTTL != 0 {
			res["max_wrapping_ttl"] = int64(perms.MaxWrappingTTL.Seconds())
		}
		if len(perms.AllowedParameters) > 0 {
			res["allowed_parameters"] = perms.AllowedParameters
		}
		if len(perms.DeniedParameters) > 0 {
			res["denied_parameters"] = perms.DeniedParameters
		}
		if len(perms.RequiredParameters) > 0 {
			res["required_parameters"] = perms.RequiredParameters
		}

		pt[s] = res
	}

	exactWalkFn := func(s string, v interface{}) bool {
		walkFn(exact, s, v)
		return false
	}

	globWalkFn := func(s string, v interface{}) bool {
		walkFn(glob, s, v)
		return false
	}

	acl.exactRules.Walk(exactWalkFn)
	acl.globRules.Walk(globWalkFn)

	resp.Data["exact_paths"] = exact
	resp.Data["glob_paths"] = glob

	return resp, nil
}

func sanitizeMountPath(path string) string {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	return path
}

func checkListingVisibility(visibility ListingVisibilityType) error {
	switch visibility {
	case ListingVisibilityHidden:
	case ListingVisibilityUnauth:
	default:
		return fmt.Errorf("invalid listing visilibity type")
	}

	return nil
}

const sysHelpRoot = `
The system backend is built-in to Vault and cannot be remounted or
unmounted. It contains the paths that are used to configure Vault itself
as well as perform core operations.
`

// sysHelp is all the help text for the sys backend.
var sysHelp = map[string][2]string{
	"config/cors": {
		"Configures or returns the current configuration of CORS settings.",
		`
This path responds to the following HTTP methods.

    GET /
        Returns the configuration of the CORS setting.

    POST /
        Sets the comma-separated list of origins that can make cross-origin requests.

    DELETE /
        Clears the CORS configuration and disables acceptance of CORS requests.
		`,
	},
	"config/ui/headers": {
		"Configures response headers that should be returned from the UI.",
		`
This path responds to the following HTTP methods.
    GET /<header>
        Returns the header value.
    POST /<header>
        Sets the header value for the UI.
    DELETE /<header>
        Clears the header value for UI.
        
    LIST /
        List the headers configured for the UI.
        `,
	},
	"init": {
		"Initializes or returns the initialization status of the Vault.",
		`
This path responds to the following HTTP methods.

    GET /
        Returns the initialization status of the Vault.

    POST /
        Initializes a new vault.
		`,
	},
	"generate-root": {
		"Reads, generates, or deletes a root token regeneration process.",
		`
This path responds to multiple HTTP methods which change the behavior. Those
HTTP methods are listed below.

    GET /attempt
        Reads the configuration and progress of the current root generation
        attempt.

    POST /attempt
        Initializes a new root generation attempt. Only a single root generation
        attempt can take place at a time. One (and only one) of otp or pgp_key
        are required.

    DELETE /attempt
        Cancels any in-progress root generation attempt. This clears any
        progress made. This must be called to change the OTP or PGP key being
        used.
		`,
	},
	"seal-status": {
		"Returns the seal status of the Vault.",
		`
This path responds to the following HTTP methods.

    GET /
        Returns the seal status of the Vault. This is an unauthenticated
        endpoint.
		`,
	},
	"seal": {
		"Seals the Vault.",
		`
This path responds to the following HTTP methods.

    PUT /
        Seals the Vault.
		`,
	},
	"unseal": {
		"Unseals the Vault.",
		`
This path responds to the following HTTP methods.

    PUT /
        Unseals the Vault.
		`,
	},
	"mounts": {
		"List the currently mounted backends.",
		`
This path responds to the following HTTP methods.

    GET /
        Lists all the mounted secret backends.

    GET /<mount point>
        Get information about the mount at the specified path.

    POST /<mount point>
        Mount a new secret backend to the mount point in the URL.

    POST /<mount point>/tune
        Tune configuration parameters for the given mount point.

    DELETE /<mount point>
        Unmount the specified mount point.
		`,
	},

	"mount": {
		`Mount a new backend at a new path.`,
		`
Mount a backend at a new path. A backend can be mounted multiple times at
multiple paths in order to configure multiple separately configured backends.
Example: you might have an AWS backend for the east coast, and one for the
west coast.
		`,
	},

	"mount_path": {
		`The path to mount to. Example: "aws/east"`,
		"",
	},

	"mount_type": {
		`The type of the backend. Example: "passthrough"`,
		"",
	},

	"mount_desc": {
		`User-friendly description for this mount.`,
		"",
	},

	"mount_config": {
		`Configuration for this mount, such as default_lease_ttl
and max_lease_ttl.`,
	},

	"mount_local": {
		`Mark the mount as a local mount, which is not replicated
and is unaffected by replication.`,
	},

	"mount_plugin_name": {
		`Name of the plugin to mount based from the name registered
in the plugin catalog.`,
	},

	"mount_options": {
		`The options to pass into the backend. Should be a json object with string keys and values.`,
	},

	"seal_wrap": {
		`Whether to turn on seal wrapping for the mount.`,
	},

	"tune_default_lease_ttl": {
		`The default lease TTL for this mount.`,
	},

	"tune_max_lease_ttl": {
		`The max lease TTL for this mount.`,
	},

	"tune_audit_non_hmac_request_keys": {
		`The list of keys in the request data object that will not be HMAC'ed by audit devices.`,
	},

	"tune_audit_non_hmac_response_keys": {
		`The list of keys in the response data object that will not be HMAC'ed by audit devices.`,
	},

	"tune_mount_options": {
		`The options to pass into the backend. Should be a json object with string keys and values.`,
	},

	"remount": {
		"Move the mount point of an already-mounted backend.",
		`
This path responds to the following HTTP methods.

    POST /sys/remount
        Changes the mount point of an already-mounted backend.
		`,
	},

	"auth_tune": {
		"Tune the configuration parameters for an auth path.",
		`Read and write the 'default-lease-ttl' and 'max-lease-ttl' values of
the auth path.`,
	},

	"mount_tune": {
		"Tune backend configuration parameters for this mount.",
		`Read and write the 'default-lease-ttl' and 'max-lease-ttl' values of
the mount.`,
	},

	"renew": {
		"Renew a lease on a secret",
		`
When a secret is read, it may optionally include a lease interval
and a boolean indicating if renew is possible. For secrets that support
lease renewal, this endpoint is used to extend the validity of the
lease and to prevent an automatic revocation.
		`,
	},

	"lease_id": {
		"The lease identifier to renew. This is included with a lease.",
		"",
	},

	"increment": {
		"The desired increment in seconds to the lease",
		"",
	},

	"revoke": {
		"Revoke a leased secret immediately",
		`
When a secret is generated with a lease, it is automatically revoked
at the end of the lease period if not renewed. However, in some cases
you may want to force an immediate revocation. This endpoint can be
used to revoke the secret with the given Lease ID.
		`,
	},

	"revoke-prefix": {
		"Revoke all secrets generated in a given prefix",
		`
Revokes all the secrets generated under a given mount prefix. As
an example, "prod/aws/" might be the AWS logical backend, and due to
a change in the "ops" policy, we may want to invalidate all the secrets
generated. We can do a revoke prefix at "prod/aws/ops" to revoke all
the ops secrets. This does a prefix match on the Lease IDs and revokes
all matching leases.
		`,
	},

	"revoke-prefix-path": {
		`The path to revoke keys under. Example: "prod/aws/ops"`,
		"",
	},

	"revoke-force": {
		"Revoke all secrets generated in a given prefix, ignoring errors.",
		`
See the path help for 'revoke-prefix'; this behaves the same, except that it
ignores errors encountered during revocation. This can be used in certain
recovery situations; for instance, when you want to unmount a backend, but it
is impossible to fix revocation errors and these errors prevent the unmount
from proceeding. This is a DANGEROUS operation as it removes Vault's oversight
of external secrets. Access to this prefix should be tightly controlled.
		`,
	},

	"revoke-force-path": {
		`The path to revoke keys under. Example: "prod/aws/ops"`,
		"",
	},

	"auth-table": {
		"List the currently enabled credential backends.",
		`
This path responds to the following HTTP methods.

    GET /
        List the currently enabled credential backends: the name, the type of
        the backend, and a user friendly description of the purpose for the
        credential backend.

    POST /<mount point>
        Enable a new auth method.

    DELETE /<mount point>
        Disable the auth method at the given mount point.
		`,
	},

	"auth": {
		`Enable a new credential backend with a name.`,
		`
Enable a credential mechanism at a new path. A backend can be mounted multiple times at
multiple paths in order to configure multiple separately configured backends.
Example: you might have an OAuth backend for GitHub, and one for Google Apps.
		`,
	},

	"auth_path": {
		`The path to mount to. Cannot be delimited. Example: "user"`,
		"",
	},

	"auth_type": {
		`The type of the backend. Example: "userpass"`,
		"",
	},

	"auth_desc": {
		`User-friendly description for this credential backend.`,
		"",
	},

	"auth_config": {
		`Configuration for this mount, such as plugin_name.`,
	},

	"auth_plugin": {
		`Name of the auth plugin to use based from the name in the plugin catalog.`,
		"",
	},

	"auth_options": {
		`The options to pass into the backend. Should be a json object with string keys and values.`,
	},

	"policy-list": {
		`List the configured access control policies.`,
		`
This path responds to the following HTTP methods.

    GET /
        List the names of the configured access control policies.

    GET /<name>
        Retrieve the rules for the named policy.

    PUT /<name>
        Add or update a policy.

    DELETE /<name>
        Delete the policy with the given name.
		`,
	},

	"policy": {
		`Read, Modify, or Delete an access control policy.`,
		`
Read the rules of an existing policy, create or update the rules of a policy,
or delete a policy.
		`,
	},

	"policy-name": {
		`The name of the policy. Example: "ops"`,
		"",
	},

	"policy-rules": {
		`The rules of the policy. Either given in HCL or JSON format.`,
		"",
	},

	"audit-hash": {
		"The hash of the given string via the given audit backend",
		"",
	},

	"audit-table": {
		"List the currently enabled audit backends.",
		`
This path responds to the following HTTP methods.

    GET /
        List the currently enabled audit backends.

    PUT /<path>
        Enable an audit backend at the given path.

    DELETE /<path>
        Disable the given audit backend.
		`,
	},

	"audit_path": {
		`The name of the backend. Cannot be delimited. Example: "mysql"`,
		"",
	},

	"audit_type": {
		`The type of the backend. Example: "mysql"`,
		"",
	},

	"audit_desc": {
		`User-friendly description for this audit backend.`,
		"",
	},

	"audit_opts": {
		`Configuration options for the audit backend.`,
		"",
	},

	"audit": {
		`Enable or disable audit backends.`,
		`
Enable a new audit backend or disable an existing backend.
		`,
	},

	"key-status": {
		"Provides information about the backend encryption key.",
		`
		Provides the current backend encryption key term and installation time.
		`,
	},

	"rotate": {
		"Rotates the backend encryption key used to persist data.",
		`
		Rotate generates a new encryption key which is used to encrypt all
		data going to the storage backend. The old encryption keys are kept so
		that data encrypted using those keys can still be decrypted.
		`,
	},

	"rekey_backup": {
		"Allows fetching or deleting the backup of the rotated unseal keys.",
		"",
	},

	"capabilities": {
		"Fetches the capabilities of the given token on the given path.",
		`Returns the capabilities of the given token on the path.
		The path will be searched for a path match in all the policies associated with the token.`,
	},

	"capabilities_self": {
		"Fetches the capabilities of the given token on the given path.",
		`Returns the capabilities of the client token on the path.
		The path will be searched for a path match in all the policies associated with the client token.`,
	},

	"capabilities_accessor": {
		"Fetches the capabilities of the token associated with the given token, on the given path.",
		`When there is no access to the token, token accessor can be used to fetch the token's capabilities
		on a given path.`,
	},

	"tidy_leases": {
		`This endpoint performs cleanup tasks that can be run if certain error
conditions have occurred.`,
		`This endpoint performs cleanup tasks that can be run to clean up the
lease entries after certain error conditions. Usually running this is not
necessary, and is only required if upgrade notes or support personnel suggest
it.`,
	},

	"wrap": {
		"Response-wraps an arbitrary JSON object.",
		`Round trips the given input data into a response-wrapped token.`,
	},

	"wrappubkey": {
		"Returns pubkeys used in some wrapping formats.",
		"Returns pubkeys used in some wrapping formats.",
	},

	"unwrap": {
		"Unwraps a response-wrapped token.",
		`Unwraps a response-wrapped token. Unlike simply reading from cubbyhole/response,
		this provides additional validation on the token, and rather than a JSON-escaped
		string, the returned response is the exact same as the contained wrapped response.`,
	},

	"wraplookup": {
		"Looks up the properties of a response-wrapped token.",
		`Returns the creation TTL and creation time of a response-wrapped token.`,
	},

	"rewrap": {
		"Rotates a response-wrapped token.",
		`Rotates a response-wrapped token; the output is a new token with the same
		response wrapped inside and the same creation TTL. The original token is revoked.`,
	},
	"audited-headers-name": {
		"Configures the headers sent to the audit logs.",
		`
This path responds to the following HTTP methods.

	GET /<name>
		Returns the setting for the header with the given name.

	POST /<name>
		Enable auditing of the given header.

	DELETE /<path>
		Disable auditing of the given header.
		`,
	},
	"audited-headers": {
		"Lists the headers configured to be audited.",
		`Returns a list of headers that have been configured to be audited.`,
	},
	"plugin-catalog": {
		"Configures the plugins known to vault",
		`
This path responds to the following HTTP methods.
		LIST /
			Returns a list of names of configured plugins.

		GET /<name>
			Retrieve the metadata for the named plugin.

		PUT /<name>
			Add or update plugin.

		DELETE /<name>
			Delete the plugin with the given name.
		`,
	},
	"plugin-catalog_name": {
		"The name of the plugin",
		"",
	},
	"plugin-catalog_sha-256": {
		`The SHA256 sum of the executable used in the
command field. This should be HEX encoded.`,
		"",
	},
	"plugin-catalog_command": {
		`The command used to start the plugin. The
executable defined in this command must exist in vault's
plugin directory.`,
		"",
	},
	"plugin-catalog_args": {
		`The args passed to plugin command.`,
		"",
	},
	"leases": {
		`View or list lease metadata.`,
		`
This path responds to the following HTTP methods.

    PUT /
        Retrieve the metadata for the provided lease id.

    LIST /<prefix>
        Lists the leases for the named prefix.
		`,
	},

	"leases-list-prefix": {
		`The path to list leases under. Example: "aws/creds/deploy"`,
		"",
	},
	"plugin-reload": {
		"Reload mounts that use a particular backend plugin.",
		`Reload mounts that use a particular backend plugin. Either the plugin name
		or the desired plugin backend mounts must be provided, but not both. In the
		case that the plugin name is provided, all mounted paths that use that plugin
		backend will be reloaded.`,
	},
	"plugin-backend-reload-plugin": {
		`The name of the plugin to reload, as registered in the plugin catalog.`,
		"",
	},
	"plugin-backend-reload-mounts": {
		`The mount paths of the plugin backends to reload.`,
		"",
	},
	"hash": {
		"Generate a hash sum for input data",
		"Generates a hash sum of the given algorithm against the given input data.",
	},
	"random": {
		"Generate random bytes",
		"This function can be used to generate high-entropy random bytes.",
	},
	"listing_visibility": {
		"Determines the visibility of the mount in the UI-specific listing endpoint. Accepted value are 'unauth' and ''.",
		"",
	},
	"passthrough_request_headers": {
		"A list of headers to whitelist and pass from the request to the backend.",
		"",
	},
	"raw": {
		"Write, Read, and Delete data directly in the Storage backend.",
		"",
	},
	"internal-ui-mounts": {
		"Information about mounts returned according to their tuned visibility. Internal API; its location, inputs, and outputs may change.",
		"",
	},
	"internal-ui-resultant-acl": {
		"Information about a token's resultant ACL. Internal API; its location, inputs, and outputs may change.",
		"",
	},
}
