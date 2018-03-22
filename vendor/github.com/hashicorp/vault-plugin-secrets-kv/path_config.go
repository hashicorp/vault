package kv

import (
	"context"
	"path"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.upgradeCheck(b.pathConfigWrite()),
			logical.CreateOperation: b.upgradeCheck(b.pathConfigWrite()),
			logical.ReadOperation:   b.upgradeCheck(b.pathConfigRead()),
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
		if config == nil {
			return nil, nil
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"max_versions": config.MaxVersions,
				"cas_required": config.CasRequired,
			},
		}, nil
	}
}

// pathConfigWrite handles create and update commands to the config
func (b *versionedKVBackend) pathConfigWrite() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		maxRaw, mOk := data.GetOk("max_versions")
		casRaw, cOk := data.GetOk("cas_required")

		// Fast path validation
		if !mOk && !cOk {
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
`
