package kv

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *versionedKVBackend) upgradeCheck(next framework.OperationFunc) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		if atomic.LoadUint32(b.upgrading) == 1 {
			return logical.ErrorResponse("Can not handle request while upgrade is in process"), logical.ErrInvalidRequest
		}

		return next(ctx, req, data)
	}
}

func (b *versionedKVBackend) Upgrade(ctx context.Context, s logical.Storage) error {
	if !b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary) {
		b.Logger().Info("versioned k/v: upgrade not running on performace replication secondary")
		return nil
	}

	if !atomic.CompareAndSwapUint32(b.upgrading, 0, 1) {
		return errors.New("upgrade already in process")
	}

	// Write upgrade canary
	info, err := proto.Marshal(&UpgradeInfo{
		StartedTime: ptypes.TimestampNow(),
	})
	if err != nil {
		return err
	}

	err = s.Put(ctx, &logical.StorageEntry{
		Key:   path.Join(b.storagePrefix, "upgrading"),
		Value: info,
	})
	if err != nil {
		return err
	}

	upgradeKey := func(key string) error {
		if strings.HasPrefix(key, b.storagePrefix) {
			return nil
		}

		// Read the old data
		data, err := s.Get(ctx, key)
		if err != nil {
			return err
		}

		locksutil.LockForKey(b.locks, key).Lock()
		defer locksutil.LockForKey(b.locks, key).Unlock()

		meta := &KeyMetadata{
			Key:      key,
			Versions: map[uint64]*VersionMetadata{},
		}

		versionKey, err := b.getVersionKey(ctx, key, 1, s)
		if err != nil {
			return err
		}

		version := &Version{
			Data:        data.Value,
			CreatedTime: ptypes.TimestampNow(),
		}

		buf, err := proto.Marshal(version)
		if err != nil {
			return err
		}

		// Store the version data
		if err := s.Put(ctx, &logical.StorageEntry{
			Key:   versionKey,
			Value: buf,
		}); err != nil {
			return err
		}

		// Store the metadata
		meta.AddVersion(version.CreatedTime, nil, 1)
		err = b.writeKeyMetadata(ctx, s, meta)
		if err != nil {
			return err
		}

		// delete the old key
		err = s.Delete(ctx, key)
		if err != nil {
			return err
		}

		return nil
	}

	// Run the actual upgrade in a go routine so we don't block the client on a
	// potentially long process.
	go func() {

		b.Logger().Info("versioned k/v: collecting keys")
		keys, err := logical.CollectKeys(ctx, s)
		if err != nil {
			b.Logger().Error("versioned k/v: upgrading resulted in error", "error", err)
			return
		}

		b.Logger().Info("versioned k/v: done collecting keys", "num_keys", len(keys))
		for i, key := range keys {
			if b.Logger().IsTrace() && i%500 == 0 {
				b.Logger().Trace("versioned k/v: upgrading keys", "progress", fmt.Sprintf("%d/%d", i, len(keys)))
			}
			err := upgradeKey(key)
			if err != nil {
				b.Logger().Error("versioned k/v: upgrading resulted in error", "error", err, "progress", fmt.Sprintf("%d/%d", i+1, len(keys)))
				return
			}
		}

		b.Logger().Info("versioned k/v: upgrading keys finished")

		// Remove the upgrading canary
		err = s.Delete(ctx, path.Join(b.storagePrefix, "upgrading"))
		if err != nil {
			b.Logger().Error("versioned k/v: removing upgrade canary resulted in an error", "error", err)
			return
		}

		atomic.StoreUint32(b.upgrading, 0)
	}()
	return nil
}
