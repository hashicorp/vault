// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kv

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathDestroy returns the path configuration for the destroy endpoint
func pathDestroy(b *versionedKVBackend) *framework.Path {
	return &framework.Path{
		Pattern: "destroy/" + framework.MatchAllRegex("path"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixKVv2,
			OperationVerb:   "destroy",
			OperationSuffix: "versions",
		},

		Fields: map[string]*framework.FieldSchema{
			"path": {
				Type:        framework.TypeString,
				Description: "Location of the secret.",
			},
			"versions": {
				Type:        framework.TypeCommaIntSlice,
				Description: "The versions to destroy. Their data will be permanently deleted.",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.upgradeCheck(b.pathDestroyWrite()),
				Responses: map[int][]framework.Response{
					http.StatusNoContent: {{
						Description: http.StatusText(http.StatusNoContent),
					}},
				},
			},
		},

		HelpSynopsis:    destroyHelpSyn,
		HelpDescription: destroyHelpDesc,
	}
}

func (b *versionedKVBackend) pathDestroyWrite() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		key := data.Get("path").(string)

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
		marshaledVersions, err := json.Marshal(&versions)
		if err != nil {
			return nil, err
		}
		kvEvent(ctx, b.Backend, "destroy", "destroy/"+key, "", true, 2,
			"current_version", fmt.Sprintf("%d", meta.CurrentVersion),
			"oldest_version", fmt.Sprintf("%d", meta.OldestVersion),
			"destroyed_versions", string(marshaledVersions),
		)
		return nil, nil
	}
}

const destroyHelpSyn = `Permanently removes one or more versions in the KV store`
const destroyHelpDesc = `
Permanently removes the specified version data for the provided key and version
numbers from the key-value store.
`
