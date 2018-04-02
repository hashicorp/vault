package kv

import (
	"context"
	"strings"

	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathDestroy returns the path configuration for the destroy endpoint
func pathDestroy(b *versionedKVBackend) *framework.Path {
	return &framework.Path{
		Pattern: "destroy/.*",
		Fields: map[string]*framework.FieldSchema{
			"versions": {
				Type:        framework.TypeCommaIntSlice,
				Description: "The versions to destroy. Their data will be permanently deleted.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.upgradeCheck(b.pathDestroyWrite()),
			logical.CreateOperation: b.upgradeCheck(b.pathDestroyWrite()),
		},

		HelpSynopsis:    destroyHelpSyn,
		HelpDescription: destroyHelpDesc,
	}
}

func (b *versionedKVBackend) pathDestroyWrite() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		key := strings.TrimPrefix(req.Path, "destroy/")

		versions := data.Get("versions").([]int)
		if len(versions) == 0 {
			return logical.ErrorResponse("no version number provided"), logical.ErrInvalidRequest
		}

		lock := locksutil.LockForKey(b.locks, key)
		lock.Lock()
		defer lock.Unlock()

		meta, err := b.getKeyMetadata(ctx, req.Storage, key)
		if err != nil {
			return nil, err
		}
		if meta == nil {
			return nil, nil
		}

		for _, verNum := range versions {
			// If there is no version, or the version is already destroyed,
			// continue
			lv := meta.Versions[uint64(verNum)]
			if lv == nil || lv.Destroyed {
				continue
			}

			lv.Destroyed = true
		}

		// Write the metadata key before deleting the versions
		err = b.writeKeyMetadata(ctx, req.Storage, meta)
		if err != nil {
			return nil, err
		}

		for _, verNum := range versions {
			// Delete versioned data
			versionKey, err := b.getVersionKey(ctx, key, uint64(verNum), req.Storage)
			if err != nil {
				return nil, err
			}

			err = req.Storage.Delete(ctx, versionKey)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	}
}

const destroyHelpSyn = `Permanently removes one or more versions in the KV store`
const destroyHelpDesc = `
Permanently removes the specified version data for the provided key and version
numbers from the key-value store.
`
