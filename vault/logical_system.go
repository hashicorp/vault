package vault

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

var (
	// protectedPaths cannot be accessed via the raw APIs.
	// This is both for security and to prevent disrupting Vault.
	protectedPaths = []string{
		barrierInitPath,
		keyringPath,
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
				Pattern: "rekey/backup$",

				Fields: map[string]*framework.FieldSchema{},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation:   b.handleRekeyRetrieve,
					logical.DeleteOperation: b.handleRekeyDelete,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["mount_tune"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["mount_tune"][1]),
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
						Description: strings.TrimSpace(sysHelp["remount_from"][0]),
					},
					"to": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["remount_to"][0]),
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

// handleRekeyRetrieve returns backed-up, PGP-encrypted unseal keys from a
// rekey operation
func (b *SystemBackend) handleRekeyRetrieve(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	backup, err := b.Core.RekeyRetrieveBackup()
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

// handleRekeyDelete deletes backed-up, PGP-encrypted unseal keys from a rekey
// operation
func (b *SystemBackend) handleRekeyDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := b.Core.RekeyDeleteBackup()
	if err != nil {
		return nil, fmt.Errorf("error during deletion of backed-up keys: %v", err)
	}

	return nil, nil
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
				"default_lease_ttl": int(entry.Config.DefaultLeaseTTL.Seconds()),
				"max_lease_ttl":     int(entry.Config.MaxLeaseTTL.Seconds()),
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
		tmpDef, err := time.ParseDuration(apiConfig.DefaultLeaseTTL)
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
		tmpMax, err := time.ParseDuration(apiConfig.MaxLeaseTTL)
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

	// Attempt remount
	if err := b.Core.remount(fromPath, toPath); err != nil {
		b.Backend.Logger().Printf("[ERR] sys: remount '%s' to '%s' failed: %v", fromPath, toPath, err)
		return handleError(err)
	}

	return nil, nil
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

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

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

// handleMountTuneWrite is used to set config settings on a backend
func (b *SystemBackend) handleMountTuneWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path == "" {
		return logical.ErrorResponse(
				"path must be specified as a string"),
			logical.ErrInvalidRequest
	}

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

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
			tmpDef, err := time.ParseDuration(defTTL)
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
			tmpMax, err := time.ParseDuration(maxTTL)
			if err != nil {
				return handleError(err)
			}
			newMax = &tmpMax
		}

		if newDefault != nil || newMax != nil {
			b.Core.mountsLock.Lock()
			defer b.Core.mountsLock.Unlock()
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
	// Get all the options
	prefix := data.Get("prefix").(string)

	// Invoke the expiration manager directly
	if err := b.Core.expiration.RevokePrefix(prefix); err != nil {
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
		info := map[string]string{
			"type":        entry.Type,
			"description": entry.Description,
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

	// Create the mount entry
	me := &MountEntry{
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
	return logical.ListResponse(policies), err
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

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

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

const sysHelpRoot = `
The system backend is built-in to Vault and cannot be remounted or
unmounted. It contains the paths that are used to configure Vault itself
as well as perform core operations.
`

// sysHelp is all the help text for the sys backend.
var sysHelp = map[string][2]string{
	"mounts": {
		"List the currently mounted backends.",
		`
List the currently mounted backends: the mount path, the type of the backend,
and a user friendly description of the purpose for the mount.
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
Change the mount point of an already-mounted backend.
		`,
	},

	"remount_from": {
		"",
		"",
	},

	"remount_to": {
		"",
		"",
	},

	"mount_tune": {
		"Tune backend configuration parameters for this mount.",
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

	"auth-table": {
		"List the currently enabled credential backends.",
		`
List the currently enabled credential backends: the name, the type of the backend,
and a user friendly description of the purpose for the credential backend.
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
List the names of the configured access control policies. Policies are associated
with client tokens to limit access to keys in the Vault.
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
List the currently enabled audit backends: the name, the type of the backend,
a user friendly description of the audit backend, and it's configuration options.
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
}
