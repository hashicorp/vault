package kv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

// pathConfig returns the path configuration for CRUD operations on the backend
// configuration.
func pathData(b *versionedKVBackend) *framework.Path {
	return &framework.Path{
		Pattern: "data/.*",
		Fields: map[string]*framework.FieldSchema{
			"version": {
				Type:        framework.TypeInt,
				Description: "If provided during a read, the value at the version number will be returned",
			},
			"options": {
				Type: framework.TypeMap,
				Description: `Options for writing a KV entry. 

Set the "cas" value to use a Check-And-Set operation. If not set the write will
be allowed. If set to 0 a write will only be allowed if the key doesn’t exist.
If the index is non-zero the write will only be allowed if the key’s current
version matches the version specified in the cas parameter.`,
			},
			"data": {
				Type:        framework.TypeMap,
				Description: "The contents of the data map will be stored and returned on read.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.upgradeCheck(b.pathDataWrite()),
			logical.CreateOperation: b.upgradeCheck(b.pathDataWrite()),
			logical.ReadOperation:   b.upgradeCheck(b.pathDataRead()),
			logical.DeleteOperation: b.upgradeCheck(b.pathDataDelete()),
		},

		HelpSynopsis:    dataHelpSyn,
		HelpDescription: dataHelpDesc,
	}
}

// pathDataRead handles read commands to a kv entry
func (b *versionedKVBackend) pathDataRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		key := strings.TrimPrefix(req.Path, "data/")

		lock := locksutil.LockForKey(b.locks, key)
		lock.RLock()
		defer lock.RUnlock()

		meta, err := b.getKeyMetadata(ctx, req.Storage, key)
		if err != nil {
			return nil, err
		}
		if meta == nil {
			return nil, nil
		}

		verNum := meta.CurrentVersion
		verParam := data.Get("version").(int)
		if verParam > 0 {
			verNum = uint64(verParam)
		}

		// If there is no version with that number, return
		vm := meta.Versions[verNum]
		if vm == nil {
			return nil, nil
		}

		resp := &logical.Response{
			Data: map[string]interface{}{
				"data": nil,
				"metadata": map[string]interface{}{
					"version":       verNum,
					"created_time":  ptypesTimestampToString(vm.CreatedTime),
					"deletion_time": ptypesTimestampToString(vm.DeletionTime),
					"destroyed":     vm.Destroyed,
				},
			},
		}

		// If the version has been deleted return metadata with a 404
		if vm.DeletionTime != nil {
			deletionTime, err := ptypes.Timestamp(vm.DeletionTime)
			if err != nil {
				return nil, err
			}

			if deletionTime.Before(time.Now()) {
				return logical.RespondWithStatusCode(resp, req, http.StatusNotFound)

			}
		}

		// If the version has been destroyed return metadata with a 404
		if vm.Destroyed {
			return logical.RespondWithStatusCode(resp, req, http.StatusNotFound)

		}

		versionKey, err := b.getVersionKey(ctx, key, verNum, req.Storage)
		if err != nil {
			return nil, err
		}

		raw, err := req.Storage.Get(ctx, versionKey)
		if err != nil {
			return nil, err
		}
		if raw == nil {
			return nil, errors.New("could not find version data")
		}

		version := &Version{}
		if err := proto.Unmarshal(raw.Value, version); err != nil {
			return nil, err
		}

		vData := map[string]interface{}{}
		if err := json.Unmarshal(version.Data, &vData); err != nil {
			return nil, err
		}

		resp.Data["data"] = vData

		return resp, nil
	}
}

// pathDataWrite handles create and update commands to a kv entry
func (b *versionedKVBackend) pathDataWrite() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		key := strings.TrimPrefix(req.Path, "data/")

		config, err := b.config(ctx, req.Storage)
		if err != nil {
			return nil, err
		}

		// Parse data, this can happen before the lock so we can fail early if
		// not set.
		var marshaledData []byte
		{
			dataRaw, ok := data.GetOk("data")
			if !ok {
				return logical.ErrorResponse("no data provided"), logical.ErrInvalidRequest
			}
			marshaledData, err = json.Marshal(dataRaw.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
		}

		lock := locksutil.LockForKey(b.locks, key)
		lock.Lock()
		defer lock.Unlock()

		meta, err := b.getKeyMetadata(ctx, req.Storage, key)
		if err != nil {
			return nil, err
		}
		if meta == nil {
			meta = &KeyMetadata{
				Key:      key,
				Versions: map[uint64]*VersionMetadata{},
			}
		}

		// Parse options
		{
			var casRaw interface{}
			var casOk bool
			optionsRaw, ok := data.GetOk("options")
			if ok {
				options := optionsRaw.(map[string]interface{})

				// Verify the CAS parameter is valid.
				casRaw, casOk = options["cas"]
			}

			switch {
			case casOk:
				var cas int
				if err := mapstructure.WeakDecode(casRaw, &cas); err != nil {
					return logical.ErrorResponse("error parsing check-and-set parameter"), logical.ErrInvalidRequest
				}
				if uint64(cas) != meta.CurrentVersion {
					return logical.ErrorResponse("check-and-set parameter did not match the current version"), logical.ErrInvalidRequest
				}
			case config.CasRequired, meta.CasRequired:
				return logical.ErrorResponse("check-and-set parameter required for this call"), logical.ErrInvalidRequest
			}
		}

		// Create a version key for the new version
		versionKey, err := b.getVersionKey(ctx, key, meta.CurrentVersion+1, req.Storage)
		if err != nil {
			return nil, err
		}
		version := &Version{
			Data:        marshaledData,
			CreatedTime: ptypes.TimestampNow(),
		}

		buf, err := proto.Marshal(version)
		if err != nil {
			return nil, err
		}

		// Write the new version
		if err := req.Storage.Put(ctx, &logical.StorageEntry{
			Key:   versionKey,
			Value: buf,
		}); err != nil {
			return nil, err
		}

		vm, versionToDelete := meta.AddVersion(version.CreatedTime, nil, config.MaxVersions)
		err = b.writeKeyMetadata(ctx, req.Storage, meta)
		if err != nil {
			return nil, err
		}

		// We create the response here so we can add warnings to it below.
		resp := &logical.Response{
			Data: map[string]interface{}{
				"version":       meta.CurrentVersion,
				"created_time":  ptypesTimestampToString(vm.CreatedTime),
				"deletion_time": ptypesTimestampToString(vm.DeletionTime),
				"destroyed":     vm.Destroyed,
			},
		}

		// Cleanup the version data that is past max version.
		if versionToDelete > 0 {

			// Create a list of version keys to delete. We will delete from the
			// back of the array so we can delete the oldest versions
			// first. If there is an error deleting one of the keys we can
			// ensure the rest will be deleted on the next go around.
			var versionKeysToDelete []string

			for i := versionToDelete; i > 0; i-- {
				versionKey, err := b.getVersionKey(ctx, key, i, req.Storage)
				if err != nil {
					resp.AddWarning(fmt.Sprintf("Error occured when cleaning up old versions, these will be cleaned up on next write: %s", err))
					return resp, nil
				}

				// We intentionally do not return these errors here. If the get
				// or delete fail they will be cleaned up on the next write.
				v, err := req.Storage.Get(ctx, versionKey)
				if err != nil {
					resp.AddWarning(fmt.Sprintf("Error occured when cleaning up old versions, these will be cleaned up on next write: %s", err))
					return resp, nil
				}

				if v == nil {
					break
				}

				// append to the end of the list
				versionKeysToDelete = append(versionKeysToDelete, versionKey)
			}

			// Walk the list backwards deleting the oldest versions first. This
			// allows us to continue the cleanup on next write if an error
			// occurs during one of the deletes.
			for i := len(versionKeysToDelete) - 1; i >= 0; i-- {
				err := req.Storage.Delete(ctx, versionKeysToDelete[i])
				if err != nil {
					resp.AddWarning(fmt.Sprintf("Error occured when cleaning up old versions, these will be cleaned up on next write: %s", err))
					break
				}
			}

		}

		return resp, nil
	}
}

func (b *versionedKVBackend) pathDataDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		key := strings.TrimPrefix(req.Path, "data/")

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

		// If there is no latest version, or the latest version is already
		// deleted or destroyed return
		lv := meta.Versions[meta.CurrentVersion]
		if lv == nil || lv.Destroyed {
			return nil, nil
		}

		if lv.DeletionTime != nil {
			deletionTime, err := ptypes.Timestamp(lv.DeletionTime)
			if err != nil {
				return nil, err
			}

			if deletionTime.Before(time.Now()) {
				return nil, nil
			}
		}

		lv.DeletionTime = ptypes.TimestampNow()

		err = b.writeKeyMetadata(ctx, req.Storage, meta)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// AddVersion adds a version to the key metadata and moves the sliding window of
// max versions. It returns the newly added version and the version to delete
// from storage.
func (k *KeyMetadata) AddVersion(createdTime, deletionTime *timestamp.Timestamp, configMaxVersions uint32) (*VersionMetadata, uint64) {
	if k.Versions == nil {
		k.Versions = map[uint64]*VersionMetadata{}
	}

	vm := &VersionMetadata{
		CreatedTime:  createdTime,
		DeletionTime: deletionTime,
	}

	k.CurrentVersion++
	k.Versions[k.CurrentVersion] = vm
	k.UpdatedTime = createdTime
	if k.CreatedTime == nil {
		k.CreatedTime = createdTime
	}

	var maxVersions uint32
	switch {
	case max(k.MaxVersions, configMaxVersions) > 0:
		maxVersions = max(k.MaxVersions, configMaxVersions)
	default:
		maxVersions = defaultMaxVersions
	}

	if uint32(k.CurrentVersion-k.OldestVersion) >= maxVersions {
		versionToDelete := k.CurrentVersion - uint64(maxVersions)
		// We need to do a loop here in the event that max versions has
		// changed and we need to delete more than one entry.
		for i := k.OldestVersion; i < versionToDelete+1; i++ {
			delete(k.Versions, i)
		}

		k.OldestVersion = versionToDelete + 1

		return vm, versionToDelete
	}

	return vm, 0
}

func max(a, b uint32) uint32 {
	if b > a {
		return b
	}

	return a
}

const dataHelpSyn = `Write, Read, and Delete data in the Key-Value Store.`
const dataHelpDesc = `
This path takes a key name and based on the opperation stores, retreives or
deletes versions of data.

If a write operation is used the endpoint takes an options object and a data
object. The options object is used to pass some options to the write command and
the data object is encrypted and stored in the storage backend. Each write
operation for a key creates a new version and does not overwrite the previous
data.

A read operation will return the latest version for a key unless the "version"
parameter is set, then it returns the version at that number.

Delete operations are a soft delete. They will mark the latest version as
deleted, but the underlying data will not be fully removed. Delete operations
can be undone.
`
