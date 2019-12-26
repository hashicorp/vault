package kv

import (
	"context"
	"path"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathConfig returns the path configuration for CRUD operations on the backend
// configuration.
func pathConfig(b *versionedKVBackend) *framework.Path {
	return &framework.Path{
		Pattern: "config$",
		Fields: map[string]*framework.FieldSchema{
			"max_versions": {
				Type:        framework.TypeInt,
				Description: "The number of versions to keep for each key. Defaults to 10",
			},
			"cas_required": {
				Type:        framework.TypeBool,
				Description: "If true, the backend will require the cas parameter to be set for each write",
			},
			"delete_version_after": {
				Type: framework.TypeSignedDurationSecond,
				Description: `
If set, the length of time before a version is deleted. A negative duration
disables the use of delete_version_after on all keys. A zero duration
clears the current setting. Accepts a Go duration format string.`,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.upgradeCheck(b.pathConfigWrite()),
				Summary:  "Configure backend level settings that are applied to every key in the key-value store.",
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.upgradeCheck(b.pathConfigWrite()),
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.upgradeCheck(b.pathConfigRead()),
				Summary:  "Read the backend level settings.",
			},
		},

		HelpSynopsis:    confHelpSyn,
		HelpDescription: confHelpDesc,
	}
}

// pathConfigWrite handles create and update commands to the config
func (b *versionedKVBackend) pathConfigRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		config, err := b.config(ctx, req.Storage)
		if err != nil {
			return nil, err
		}

		rdata := map[string]interface{}{
			"max_versions": config.MaxVersions,
			"cas_required": config.CasRequired,
		}

		var deleteVersionAfter time.Duration
		if config.GetDeleteVersionAfter() != nil {
			deleteVersionAfter, err = ptypes.Duration(config.GetDeleteVersionAfter())
			if err != nil {
				return nil, err
			}
			rdata["delete_version_after"] = deleteVersionAfter.String()
		}

		return &logical.Response{
			Data: rdata,
		}, nil
	}
}

// pathConfigWrite handles create and update commands to the config
func (b *versionedKVBackend) pathConfigWrite() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		maxRaw, mOk := data.GetOk("max_versions")
		casRaw, cOk := data.GetOk("cas_required")
		dvaRaw, dvaOk := data.GetOk("delete_version_after")

		// Fast path validation
		if !mOk && !cOk && !dvaOk {
			return nil, nil
		}

		config, err := b.config(ctx, req.Storage)
		if err != nil {
			return nil, err
		}

		if mOk {
			config.MaxVersions = uint32(maxRaw.(int))
		}
		if cOk {
			config.CasRequired = casRaw.(bool)
		}

		if dvaOk {
			dva := dvaRaw.(int)
			switch {
			case dva < 0:
				config.DisableDeleteVersionAfter()
			case dva == 0:
				config.ResetDeleteVersionAfter()
			default:
				config.DeleteVersionAfter = ptypes.DurationProto(time.Duration(dva) * time.Second)
			}
		}

		bytes, err := proto.Marshal(config)
		if err != nil {
			return nil, err
		}

		err = req.Storage.Put(ctx, &logical.StorageEntry{
			Key:   path.Join(b.storagePrefix, configPath),
			Value: bytes,
		})
		if err != nil {
			return nil, err
		}

		b.globalConfigLock.Lock()
		defer b.globalConfigLock.Unlock()

		b.globalConfig = config

		return nil, nil
	}
}

const confHelpSyn = `Configures settings for the KV store`
const confHelpDesc = `
This path configures backend level settings that are applied to every key in the
key-value store. This parameter accetps:

	* max_versions (int) - The number of versions to keep for each key. Defaults
	  to 10

	* cas_required (bool) - If true, the backend will require the cas parameter
	  to be set for each write

	* delete_version_after (duration) - If set, the length of time before a
	  version is deleted. A negative duration disables the use of
	  delete_version_after on all keys. A zero duration clears the current
	  setting. Accepts a Go duration format string.
`
