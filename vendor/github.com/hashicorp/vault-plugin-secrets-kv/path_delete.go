// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kv

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathsDelete returns the path configuration for the delete and undelete paths
func pathsDelete(b *versionedKVBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "delete/" + framework.MatchAllRegex("path"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixKVv2,
				OperationVerb:   "delete",
				OperationSuffix: "versions",
			},

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: "Location of the secret.",
				},
				"versions": {
					Type:        framework.TypeCommaIntSlice,
					Description: "The versions to be archived. The versioned data will not be deleted, but it will no longer be returned in normal get requests.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.upgradeCheck(b.pathDeleteWrite()),
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: http.StatusText(http.StatusNoContent),
						}},
					},
				},
			},

			HelpSynopsis:    deleteHelpSyn,
			HelpDescription: deleteHelpDesc,
		},
		{
			Pattern: "undelete/" + framework.MatchAllRegex("path"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixKVv2,
				OperationVerb:   "undelete",
				OperationSuffix: "versions",
			},

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: "Location of the secret.",
				},
				"versions": {
					Type:        framework.TypeCommaIntSlice,
					Description: "The versions to unarchive. The versions will be restored and their data will be returned on normal get requests.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.upgradeCheck(b.pathUndeleteWrite()),
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: http.StatusText(http.StatusNoContent),
						}},
					},
				},
			},

			HelpSynopsis:    undeleteHelpSyn,
			HelpDescription: undeleteHelpDesc,
		},
	}
}

// pathUndeleteWrite is used to undelete a set of versions
func (b *versionedKVBackend) pathUndeleteWrite() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		key := data.Get("path").(string)

		versions := data.Get("versions").([]int)
		if len(versions) == 0 {
			return logical.ErrorResponse("No version number provided"), logical.ErrInvalidRequest
		}

		config, err := b.config(ctx, req.Storage)
		if err != nil {
			return nil, err
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
			// If there is no version or the version is destroyed continue
			lv := meta.Versions[uint64(verNum)]
			if lv == nil || lv.Destroyed {
				continue
			}
			lv.DeletionTime = nil

			if !config.IsDeleteVersionAfterDisabled() {
				if dtime, ok := deletionTime(time.Now(), deleteVersionAfter(config), deleteVersionAfter(meta)); ok {
					dt, err := ptypes.TimestampProto(dtime)
					if err != nil {
						return logical.ErrorResponse("error setting deletion_time: converting %v to protobuf: %v", dtime, err), logical.ErrInvalidRequest
					}
					lv.DeletionTime = dt
				}
			}
		}
		err = b.writeKeyMetadata(ctx, req.Storage, meta)
		if err != nil {
			return nil, err
		}
		marshaledVersions, err := json.Marshal(&versions)
		if err != nil {
			return nil, err
		}
		kvEvent(ctx, b.Backend, "undelete", "undelete/"+key, "data/"+key, true, 2,
			"current_version", fmt.Sprintf("%d", meta.CurrentVersion),
			"oldest_version", fmt.Sprintf("%d", meta.OldestVersion),
			"undeleted_versions", string(marshaledVersions),
		)
		return nil, nil
	}
}

// pathDeleteWrite is used to delete a set of versions.
func (b *versionedKVBackend) pathDeleteWrite() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		key := data.Get("path").(string)

		versions := data.Get("versions").([]int)
		if len(versions) == 0 {
			return logical.ErrorResponse("No version number provided"), logical.ErrInvalidRequest
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
			// If there is no latest version, or the latest version is already
			// deleted or destroyed continue
			lv := meta.Versions[uint64(verNum)]
			if lv == nil || lv.Destroyed {
				continue
			}

			if lv.DeletionTime != nil {
				deletionTime, err := ptypes.Timestamp(lv.DeletionTime)
				if err != nil {
					return nil, err
				}

				if deletionTime.Before(time.Now()) {
					continue
				}
			}

			lv.DeletionTime = ptypes.TimestampNow()
		}

		err = b.writeKeyMetadata(ctx, req.Storage, meta)
		if err != nil {
			return nil, err
		}
		marshaledVersions, err := json.Marshal(&versions)
		if err != nil {
			return nil, err
		}
		kvEvent(ctx, b.Backend, "delete", "delete/"+key, "", true, 2,
			"current_version", fmt.Sprintf("%d", meta.CurrentVersion),
			"oldest_version", fmt.Sprintf("%d", meta.OldestVersion),
			"deleted_versions", string(marshaledVersions),
		)
		return nil, nil
	}
}

const deleteHelpSyn = `Marks one or more versions as deleted in the KV store.`
const deleteHelpDesc = `
Deletes the data for the provided version and path in the key-value store. The
versioned data will not be fully removed, but marked as deleted and will no
longer be returned in normal get requests. This operation can be undone.
`

const undeleteHelpSyn = `Undeletes one or more versions from the KV store.`
const undeleteHelpDesc = `
Undeletes the data for the provided version and path in the key-value store.
This restores the data, allowing it to be returned on get requests.
`
