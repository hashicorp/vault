package vault

import (
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func NewSystemBackend(core *Core) logical.Backend {
	b := &SystemBackend{Core: core}

	return &framework.Backend{
		PathsRoot: []string{
			"mounts/*",
			"remount",
			"revoke-prefix/*",
		},

		Paths: []*framework.Path{
			&framework.Path{
				Pattern: "mounts$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.handleMountTable,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["mounts"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["mounts"][1]),
			},

			&framework.Path{
				Pattern: "mounts/(?P<path>.+)",

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
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation:  b.handleMount,
					logical.DeleteOperation: b.handleUnmount,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["mount"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["mount"][1]),
			},

			&framework.Path{
				Pattern: "remount",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation: b.handleRemount,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["remount"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["remount"][1]),
			},

			&framework.Path{
				Pattern: "renew/(?P<vault_id>.+)",

				Fields: map[string]*framework.FieldSchema{
					"vault_id": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["vault_id"][0]),
					},
					"increment": &framework.FieldSchema{
						Type:        framework.TypeInt,
						Description: strings.TrimSpace(sysHelp["increment"][0]),
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation: b.handleRenew,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["renew"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["renew"][1]),
			},

			&framework.Path{
				Pattern: "revoke/(?P<vault_id>.+)",

				Fields: map[string]*framework.FieldSchema{
					"vault_id": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: strings.TrimSpace(sysHelp["vault_id"][0]),
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation: b.handleRevoke,
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
					logical.WriteOperation: b.handleRevokePrefix,
				},

				HelpSynopsis:    strings.TrimSpace(sysHelp["revoke-prefix"][0]),
				HelpDescription: strings.TrimSpace(sysHelp["revoke-prefix"][1]),
			},
		},
	}
}

// SystemBackend implements logical.Backend and is used to interact with
// the core of the system. This backend is hardcoded to exist at the "sys"
// prefix. Conceptually it is similar to procfs on Linux.
type SystemBackend struct {
	Core *Core
}

// handleMountTable handles the "mounts" endpoint to provide the mount table
func (b *SystemBackend) handleMountTable(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Core.mounts.Lock()
	defer b.Core.mounts.Unlock()

	resp := &logical.Response{
		IsSecret: false,
		Data:     make(map[string]interface{}),
	}
	for _, entry := range b.Core.mounts.Entries {
		info := map[string]string{
			"type":        entry.Type,
			"description": entry.Description,
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

	// Attempt mount
	if err := b.Core.mount(me); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleUnmount is used to unmount a path
//
// TODO: Run a rollback operation
// TODO: Return an error if things go bad.
// TODO: Think through error scenario: umount should clean up everything
//   and should fail if it can't. Perhaps add a "force" flag to YOLO it.
func (b *SystemBackend) handleUnmount(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	suffix := strings.TrimPrefix(req.Path, "mounts/")
	if len(suffix) == 0 {
		return logical.ErrorResponse("path cannot be blank"), logical.ErrInvalidRequest
	}

	// Attempt unmount
	if err := b.Core.unmount(suffix); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	return nil, nil
}

// handleRemount is used to remount a path
func (b *SystemBackend) handleRemount(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Only accept write operations
	switch req.Operation {
	case logical.WriteOperation:
	default:
		return nil, logical.ErrUnsupportedOperation
	}

	// Get the paths
	fromPath := req.GetString("from")
	toPath := req.GetString("to")
	if fromPath == "" || toPath == "" {
		return logical.ErrorResponse(
				"both 'from' and 'to' path must be specified as a string"),
			logical.ErrInvalidRequest
	}

	// Attempt remount
	if err := b.Core.remount(fromPath, toPath); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	return nil, nil
}

// handleRenew is used to renew a lease with a given VaultID
func (b *SystemBackend) handleRenew(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	vaultID := data.Get("vault_id").(string)
	incrementRaw := data.Get("increment").(int)

	// Convert the increment
	increment := time.Duration(incrementRaw) * time.Second

	// Invoke the expiration manager directly
	resp, err := b.Core.expiration.Renew(vaultID, increment)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return resp, err
}

// handleRevoke is used to revoke a given VaultID
func (b *SystemBackend) handleRevoke(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	vaultID := data.Get("vault_id").(string)

	// Invoke the expiration manager directly
	if err := b.Core.expiration.Revoke(vaultID); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleRevokePrefix is used to revoke a prefix with many VaultIDs
func (b *SystemBackend) handleRevokePrefix(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all the options
	prefix := data.Get("prefix").(string)

	// Invoke the expiration manager directly
	if err := b.Core.expiration.RevokePrefix(prefix); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

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

	"remount": {
		"Move the mount point of an already-mounted backend.",
		`
Change the mount point of an already-mounted backend.
		`,
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

	"vault_id": {
		"The vault identifier to renew. This is included with a lease.",
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
used to revoke the secret with the given Vault ID.
		`,
	},

	"revoke-prefix": {
		"Revoke all secrets generated in a given prefix",
		`
Revokes all the secrets generated under a given mount prefix. As
an example, "prod/aws/" might be the AWS logical backend, and due to
a change in the "ops" policy, we may want to invalidate all the secrets
generated. We can do a revoke prefix at "prod/aws/ops" to revoke all
the ops secrets. This does a prefix match on the Vault IDs and revokes
all matching leases.
		`,
	},

	"revoke-prefix-path": {
		`The path to revoke keys under. Example: "prod/aws/ops"`,
		"",
	},
}
