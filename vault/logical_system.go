package vault

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/helper/duration"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

var (
	// protectedPaths cannot be accessed via the raw APIs.
	// This is both for security and to prevent disrupting Vault.
	protectedPaths = []string{
		"core",
	}
)

func NewSystemBackend(core *Core, config *logical.BackendConfig) logical.Backend {
	b := &SystemBackend{
		Core: core,
	}

	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(sysHelpRoot),

		PathsSpecial: &logical.Paths{
			Root: []string{
				"auth/*",
				"remount",
				"revoke-prefix/*",
				"audit",
				"audit/*",
				"raw/*",
				"rotate",
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
						Type:        framework.TypeString,
						Description: "Path on which capabilities are being queried.",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.handleCapabilitiesAccessor,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["capabilities_accessor"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["capabilities_accessor"][1]),
			},

			&framework.Path{
				Pattern: "capabilities$",

				Fields: map[string]*framework.FieldSchema{
					"token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token for which capabilities are being queried.",
					},
					"path": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Path on which capabilities are being queried.",
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
						Type:        framework.TypeString,
						Description: "Path on which capabilities are being queried.",
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
				Pattern: "renew/(?P<lease_id>.+)",

				Fields: map[string]*framework.FieldSchema{
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
				Pattern: "revoke/(?P<lease_id>.+)",

				Fields: map[string]*framework.FieldSchema{
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
				Pattern: "revoke-force/(?P<prefix>.+)",

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
				Pattern: "revoke-prefix/(?P<prefix>.+)",

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
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.handleEnableAuth,
					logical.DeleteOperation: b.handleDisableAuth,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["auth"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["auth"][1]),
			},

			&framework.Path{
				Pattern: "policy$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.handlePolicyList,
					logical.ListOperation: b.handlePolicyList,
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
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation:   b.handlePolicyRead,
					logical.UpdateOperation: b.handlePolicySet,
					logical.DeleteOperation: b.handlePolicyDelete,
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
						Type:        framework.TypeMap,
						Description: strings.TrimSpace(sysHelp["audit_opts"][0]),
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
				Pattern: "raw/(?P<path>.+)",

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
				},
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
		},
	}

	b.Backend.Setup(config)

	return b.Backend
}

// SystemBackend implements logical.Backend and is used to interact with
// the core of the system. This backend is hardcoded to exist at the "sys"
// prefix. Conceptually it is similar to procfs on Linux.
type SystemBackend struct {
	Core    *Core
	Backend *framework.Backend
}

// handleCapabilitiesreturns the ACL capabilities of the token for a given path
func (b *SystemBackend) handleCapabilities(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	capabilities, err := b.Core.Capabilities(d.Get("token").(string), d.Get("path").(string))
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"capabilities": capabilities,
		},
	}, nil
}

// handleCapabilitiesAccessor returns the ACL capabilities of the token associted
// with the given accessor for a given path.
func (b *SystemBackend) handleCapabilitiesAccessor(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	accessor := d.Get("accessor").(string)
	if accessor == "" {
		return logical.ErrorResponse("missing accessor"), nil
	}

	token, err := b.Core.tokenStore.lookupByAccessor(accessor)
	if err != nil {
		return nil, err
	}

	capabilities, err := b.Core.Capabilities(token, d.Get("path").(string))
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"capabilities": capabilities,
		},
	}, nil
}

// handleRekeyRetrieve returns backed-up, PGP-encrypted unseal keys from a
// rekey operation
func (b *SystemBackend) handleRekeyRetrieve(
	req *logical.Request,
	data *framework.FieldData,
	recovery bool) (*logical.Response, error) {
	backup, err := b.Core.RekeyRetrieveBackup(recovery)
	if err != nil {
		return nil, fmt.Errorf("unable to look up backed-up keys: %v", err)
	}
	if backup == nil {
		return logical.ErrorResponse("no backed-up keys found"), nil
	}

	// Format the status
	resp := &logical.Response{
		Data: map[string]interface{}{
			"nonce": backup.Nonce,
			"keys":  backup.Keys,
		},
	}

	return resp, nil
}

func (b *SystemBackend) handleRekeyRetrieveBarrier(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRekeyRetrieve(req, data, false)
}

func (b *SystemBackend) handleRekeyRetrieveRecovery(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRekeyRetrieve(req, data, true)
}

// handleRekeyDelete deletes backed-up, PGP-encrypted unseal keys from a rekey
// operation
func (b *SystemBackend) handleRekeyDelete(
	req *logical.Request,
	data *framework.FieldData,
	recovery bool) (*logical.Response, error) {
	err := b.Core.RekeyDeleteBackup(recovery)
	if err != nil {
		return nil, fmt.Errorf("error during deletion of backed-up keys: %v", err)
	}

	return nil, nil
}
func (b *SystemBackend) handleRekeyDeleteBarrier(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRekeyDelete(req, data, false)
}

func (b *SystemBackend) handleRekeyDeleteRecovery(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRekeyDelete(req, data, true)
}

// handleMountTable handles the "mounts" endpoint to provide the mount table
func (b *SystemBackend) handleMountTable(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Core.mountsLock.RLock()
	defer b.Core.mountsLock.RUnlock()

	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}

	for _, entry := range b.Core.mounts.Entries {
		info := map[string]interface{}{
			"type":        entry.Type,
			"description": entry.Description,
			"config": map[string]interface{}{
				"default_lease_ttl": int64(entry.Config.DefaultLeaseTTL.Seconds()),
				"max_lease_ttl":     int64(entry.Config.MaxLeaseTTL.Seconds()),
			},
		}

		resp.Data[entry.Path] = info
	}

	return resp, nil
}

// handleMount is used to mount a new path
func (b *SystemBackend) handleMount(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	path := data.Get("path").(string)
	logicalType := data.Get("type").(string)
	description := data.Get("description").(string)

	path = sanitizeMountPath(path)

	var config MountConfig

	var apiConfig struct {
		DefaultLeaseTTL string `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"`
		MaxLeaseTTL     string `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`
	}
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
		tmpDef, err := duration.ParseDurationSecond(apiConfig.DefaultLeaseTTL)
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
		tmpMax, err := duration.ParseDurationSecond(apiConfig.MaxLeaseTTL)
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

	if config.DefaultLeaseTTL > b.Core.maxLeaseTTL {
		return logical.ErrorResponse(fmt.Sprintf(
				"given default lease TTL greater than system max lease TTL of %d", int(b.Core.maxLeaseTTL.Seconds()))),
			logical.ErrInvalidRequest
	}

	if logicalType == "" {
		return logical.ErrorResponse(
				"backend type must be specified as a string"),
			logical.ErrInvalidRequest
	}

	// Create the mount entry
	me := &MountEntry{
		Table:       mountTableType,
		Path:        path,
		Type:        logicalType,
		Description: description,
		Config:      config,
	}

	// Attempt mount
	if err := b.Core.mount(me); err != nil {
		b.Backend.Logger().Printf("[ERR] sys: mount %s failed: %v", me.Path, err)
		return handleError(err)
	}

	return nil, nil
}

// used to intercept an HTTPCodedError so it goes back to callee
func handleError(
	err error) (*logical.Response, error) {
	switch err.(type) {
	case logical.HTTPCodedError:
		return logical.ErrorResponse(err.Error()), err
	default:
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
}

// handleUnmount is used to unmount a path
func (b *SystemBackend) handleUnmount(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	suffix := strings.TrimPrefix(req.Path, "mounts/")
	if len(suffix) == 0 {
		return logical.ErrorResponse("path cannot be blank"), logical.ErrInvalidRequest
	}

	suffix = sanitizeMountPath(suffix)

	// Attempt unmount
	if err := b.Core.unmount(suffix); err != nil {
		b.Backend.Logger().Printf("[ERR] sys: unmount '%s' failed: %v", suffix, err)
		return handleError(err)
	}

	return nil, nil
}

// handleRemount is used to remount a path
func (b *SystemBackend) handleRemount(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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

	// Attempt remount
	if err := b.Core.remount(fromPath, toPath); err != nil {
		b.Backend.Logger().Printf("[ERR] sys: remount '%s' to '%s' failed: %v", fromPath, toPath, err)
		return handleError(err)
	}

	return nil, nil
}

// handleAuthTuneRead is used to get config settings on a auth path
func (b *SystemBackend) handleAuthTuneRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse(
				"path must be specified as a string"),
			logical.ErrInvalidRequest
	}
	return b.handleTuneReadCommon("auth/" + path)
}

// handleMountTuneRead is used to get config settings on a backend
func (b *SystemBackend) handleMountTuneRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse(
				"path must be specified as a string"),
			logical.ErrInvalidRequest
	}

	// This call will read both logical backend's configuration as well as auth backends'.
	// Retaining this behavior for backward compatibility. If this behavior is not desired,
	// an error can be returned if path has a prefix of "auth/".
	return b.handleTuneReadCommon(path)
}

// handleTuneReadCommon returns the config settings of a path
func (b *SystemBackend) handleTuneReadCommon(path string) (*logical.Response, error) {
	path = sanitizeMountPath(path)

	sysView := b.Core.router.MatchingSystemView(path)
	if sysView == nil {
		err := fmt.Errorf("[ERR] sys: cannot fetch sysview for path %s", path)
		b.Backend.Logger().Print(err)
		return handleError(err)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"default_lease_ttl": int(sysView.DefaultLeaseTTL().Seconds()),
			"max_lease_ttl":     int(sysView.MaxLeaseTTL().Seconds()),
		},
	}

	return resp, nil
}

// handleAuthTuneWrite is used to set config settings on an auth path
func (b *SystemBackend) handleAuthTuneWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse("path must be specified as a string"),
			logical.ErrInvalidRequest
	}
	return b.handleTuneWriteCommon("auth/"+path, data)
}

// handleMountTuneWrite is used to set config settings on a backend
func (b *SystemBackend) handleMountTuneWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse("path must be specified as a string"),
			logical.ErrInvalidRequest
	}
	// This call will write both logical backend's configuration as well as auth backends'.
	// Retaining this behavior for backward compatibility. If this behavior is not desired,
	// an error can be returned if path has a prefix of "auth/".
	return b.handleTuneWriteCommon(path, data)
}

// handleTuneWriteCommon is used to set config settings on a path
func (b *SystemBackend) handleTuneWriteCommon(
	path string, data *framework.FieldData) (*logical.Response, error) {
	path = sanitizeMountPath(path)

	// Prevent protected paths from being changed
	for _, p := range untunableMounts {
		if strings.HasPrefix(path, p) {
			err := fmt.Errorf("[ERR] core: cannot tune '%s'", path)
			b.Backend.Logger().Print(err)
			return handleError(err)
		}
	}

	mountEntry := b.Core.router.MatchingMountEntry(path)
	if mountEntry == nil {
		err := fmt.Errorf("[ERR] sys: tune of path '%s' failed: no mount entry found", path)
		b.Backend.Logger().Print(err)
		return handleError(err)
	}

	var lock *sync.RWMutex
	switch {
	case strings.HasPrefix(path, "auth/"):
		lock = &b.Core.authLock
	default:
		lock = &b.Core.mountsLock
	}

	// Timing configuration parameters
	{
		var newDefault, newMax *time.Duration
		defTTL := data.Get("default_lease_ttl").(string)
		switch defTTL {
		case "":
		case "system":
			tmpDef := time.Duration(0)
			newDefault = &tmpDef
		default:
			tmpDef, err := duration.ParseDurationSecond(defTTL)
			if err != nil {
				return handleError(err)
			}
			newDefault = &tmpDef
		}

		maxTTL := data.Get("max_lease_ttl").(string)
		switch maxTTL {
		case "":
		case "system":
			tmpMax := time.Duration(0)
			newMax = &tmpMax
		default:
			tmpMax, err := duration.ParseDurationSecond(maxTTL)
			if err != nil {
				return handleError(err)
			}
			newMax = &tmpMax
		}

		if newDefault != nil || newMax != nil {
			lock.Lock()
			defer lock.Unlock()

			if err := b.tuneMountTTLs(path, &mountEntry.Config, newDefault, newMax); err != nil {
				b.Backend.Logger().Printf("[ERR] sys: tune of path '%s' failed: %v", path, err)
				return handleError(err)
			}
		}
	}

	return nil, nil
}

// handleRenew is used to renew a lease with a given LeaseID
func (b *SystemBackend) handleRenew(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	leaseID := data.Get("lease_id").(string)
	incrementRaw := data.Get("increment").(int)

	// Convert the increment
	increment := time.Duration(incrementRaw) * time.Second

	// Invoke the expiration manager directly
	resp, err := b.Core.expiration.Renew(leaseID, increment)
	if err != nil {
		b.Backend.Logger().Printf("[ERR] sys: renew '%s' failed: %v", leaseID, err)
		return handleError(err)
	}
	return resp, err
}

// handleRevoke is used to revoke a given LeaseID
func (b *SystemBackend) handleRevoke(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	leaseID := data.Get("lease_id").(string)

	// Invoke the expiration manager directly
	if err := b.Core.expiration.Revoke(leaseID); err != nil {
		b.Backend.Logger().Printf("[ERR] sys: revoke '%s' failed: %v", leaseID, err)
		return handleError(err)
	}
	return nil, nil
}

// handleRevokePrefix is used to revoke a prefix with many LeaseIDs
func (b *SystemBackend) handleRevokePrefix(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRevokePrefixCommon(req, data, false)
}

// handleRevokeForce is used to revoke a prefix with many LeaseIDs, ignoring errors
func (b *SystemBackend) handleRevokeForce(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
		b.Backend.Logger().Printf("[ERR] sys: revoke prefix '%s' failed: %v", prefix, err)
		return handleError(err)
	}
	return nil, nil
}

// handleAuthTable handles the "auth" endpoint to provide the auth table
func (b *SystemBackend) handleAuthTable(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Core.authLock.RLock()
	defer b.Core.authLock.RUnlock()

	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}
	for _, entry := range b.Core.auth.Entries {
		info := map[string]interface{}{
			"type":        entry.Type,
			"description": entry.Description,
			"config": map[string]interface{}{
				"default_lease_ttl": int64(entry.Config.DefaultLeaseTTL.Seconds()),
				"max_lease_ttl":     int64(entry.Config.MaxLeaseTTL.Seconds()),
			},
		}
		resp.Data[entry.Path] = info
	}
	return resp, nil
}

// handleEnableAuth is used to enable a new credential backend
func (b *SystemBackend) handleEnableAuth(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	path := data.Get("path").(string)
	logicalType := data.Get("type").(string)
	description := data.Get("description").(string)

	if logicalType == "" {
		return logical.ErrorResponse(
				"backend type must be specified as a string"),
			logical.ErrInvalidRequest
	}

	path = sanitizeMountPath(path)

	// Create the mount entry
	me := &MountEntry{
		Table:       credentialTableType,
		Path:        path,
		Type:        logicalType,
		Description: description,
	}

	// Attempt enabling
	if err := b.Core.enableCredential(me); err != nil {
		b.Backend.Logger().Printf("[ERR] sys: enable auth %s failed: %v", me.Path, err)
		return handleError(err)
	}
	return nil, nil
}

// handleDisableAuth is used to disable a credential backend
func (b *SystemBackend) handleDisableAuth(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	suffix := strings.TrimPrefix(req.Path, "auth/")
	if len(suffix) == 0 {
		return logical.ErrorResponse("path cannot be blank"), logical.ErrInvalidRequest
	}

	suffix = sanitizeMountPath(suffix)

	// Attempt disable
	if err := b.Core.disableCredential(suffix); err != nil {
		b.Backend.Logger().Printf("[ERR] sys: disable auth '%s' failed: %v", suffix, err)
		return handleError(err)
	}
	return nil, nil
}

// handlePolicyList handles the "policy" endpoint to provide the enabled policies
func (b *SystemBackend) handlePolicyList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the configured policies
	policies, err := b.Core.policyStore.ListPolicies()

	// Add the special "root" policy
	policies = append(policies, "root")
	resp := logical.ListResponse(policies)

	// Backwords compatibility
	resp.Data["policies"] = resp.Data["keys"]

	return resp, err
}

// handlePolicyRead handles the "policy/<name>" endpoint to read a policy
func (b *SystemBackend) handlePolicyRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	policy, err := b.Core.policyStore.GetPolicy(name)
	if err != nil {
		return handleError(err)
	}

	if policy == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"name":  name,
			"rules": policy.Raw,
		},
	}, nil
}

// handlePolicySet handles the "policy/<name>" endpoint to set a policy
func (b *SystemBackend) handlePolicySet(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	rules := data.Get("rules").(string)

	// Validate the rules parse
	parse, err := Parse(rules)
	if err != nil {
		return handleError(err)
	}

	// Override the name
	parse.Name = strings.ToLower(name)

	// Update the policy
	if err := b.Core.policyStore.SetPolicy(parse); err != nil {
		return handleError(err)
	}
	return nil, nil
}

// handlePolicyDelete handles the "policy/<name>" endpoint to delete a policy
func (b *SystemBackend) handlePolicyDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	if err := b.Core.policyStore.DeletePolicy(name); err != nil {
		return handleError(err)
	}
	return nil, nil
}

// handleAuditTable handles the "audit" endpoint to provide the audit table
func (b *SystemBackend) handleAuditTable(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
		}
		resp.Data[entry.Path] = info
	}
	return resp, nil
}

// handleAuditHash is used to fetch the hash of the given input data with the
// specified audit backend's salt
func (b *SystemBackend) handleAuditHash(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	input := data.Get("input").(string)
	if input == "" {
		return logical.ErrorResponse("the \"input\" parameter is empty"), nil
	}

	path = sanitizeMountPath(path)

	hash, err := b.Core.auditBroker.GetHash(path, input)
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
func (b *SystemBackend) handleEnableAudit(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	path := data.Get("path").(string)
	backendType := data.Get("type").(string)
	description := data.Get("description").(string)
	options := data.Get("options").(map[string]interface{})

	optionMap := make(map[string]string)
	for k, v := range options {
		vStr, ok := v.(string)
		if !ok {
			return logical.ErrorResponse("options must be string valued"),
				logical.ErrInvalidRequest
		}
		optionMap[k] = vStr
	}

	// Create the mount entry
	me := &MountEntry{
		Table:       auditTableType,
		Path:        path,
		Type:        backendType,
		Description: description,
		Options:     optionMap,
	}

	// Attempt enabling
	if err := b.Core.enableAudit(me); err != nil {
		b.Backend.Logger().Printf("[ERR] sys: enable audit %s failed: %v", me.Path, err)
		return handleError(err)
	}
	return nil, nil
}

// handleDisableAudit is used to disable an audit backend
func (b *SystemBackend) handleDisableAudit(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)

	// Attempt disable
	if err := b.Core.disableAudit(path); err != nil {
		b.Backend.Logger().Printf("[ERR] sys: disable audit '%s' failed: %v", path, err)
		return handleError(err)
	}
	return nil, nil
}

// handleRawRead is used to read directly from the barrier
func (b *SystemBackend) handleRawRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)

	// Prevent access of protected paths
	for _, p := range protectedPaths {
		if strings.HasPrefix(path, p) {
			err := fmt.Sprintf("cannot read '%s'", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	entry, err := b.Core.barrier.Get(path)
	if err != nil {
		return handleError(err)
	}
	if entry == nil {
		return nil, nil
	}
	resp := &logical.Response{
		Data: map[string]interface{}{
			"value": string(entry.Value),
		},
	}
	return resp, nil
}

// handleRawWrite is used to write directly to the barrier
func (b *SystemBackend) handleRawWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
	if err := b.Core.barrier.Put(entry); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleRawDelete is used to delete directly from the barrier
func (b *SystemBackend) handleRawDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)

	// Prevent access of protected paths
	for _, p := range protectedPaths {
		if strings.HasPrefix(path, p) {
			err := fmt.Sprintf("cannot delete '%s'", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	if err := b.Core.barrier.Delete(path); err != nil {
		return handleError(err)
	}
	return nil, nil
}

// handleKeyStatus returns status information about the backend key
func (b *SystemBackend) handleKeyStatus(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get the key info
	info, err := b.Core.barrier.ActiveKeyInfo()
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"term":         info.Term,
			"install_time": info.InstallTime.Format(time.RFC3339),
		},
	}
	return resp, nil
}

// handleRotate is used to trigger a key rotation
func (b *SystemBackend) handleRotate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Rotate to the new term
	newTerm, err := b.Core.barrier.Rotate()
	if err != nil {
		b.Backend.Logger().Printf("[ERR] sys: failed to create new encryption key: %v", err)
		return handleError(err)
	}
	b.Backend.Logger().Printf("[INFO] sys: installed new encryption key")

	// In HA mode, we need to an upgrade path for the standby instances
	if b.Core.ha != nil {
		// Create the upgrade path to the new term
		if err := b.Core.barrier.CreateUpgrade(newTerm); err != nil {
			b.Backend.Logger().Printf("[ERR] sys: failed to create new upgrade for key term %d: %v", newTerm, err)
		}

		// Schedule the destroy of the upgrade path
		time.AfterFunc(keyRotateGracePeriod, func() {
			if err := b.Core.barrier.DestroyUpgrade(newTerm); err != nil {
				b.Backend.Logger().Printf("[ERR] sys: failed to destroy upgrade for key term %d: %v", newTerm, err)
			}
		})
	}
	return nil, nil
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

const sysHelpRoot = `
The system backend is built-in to Vault and cannot be remounted or
unmounted. It contains the paths that are used to configure Vault itself
as well as perform core operations.
`

// sysHelp is all the help text for the sys backend.
var sysHelp = map[string][2]string{
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

	"tune_default_lease_ttl": {
		`The default lease TTL for this mount.`,
	},

	"tune_max_lease_ttl": {
		`The max lease TTL for this mount.`,
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
        Enable a new auth backend.

    DELETE /<mount point>
        Disable the auth backend at the given mount point.
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
		`User-friendly description for this crential backend.`,
		"",
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
}
