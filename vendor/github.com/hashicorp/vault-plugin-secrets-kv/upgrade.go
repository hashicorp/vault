// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kv

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *versionedKVBackend) perfSecondaryCheck() bool {
	replState := b.System().ReplicationState()
	if (!b.System().LocalMount() && replState.HasState(consts.ReplicationPerformanceSecondary)) ||
		replState.HasState(consts.ReplicationPerformanceStandby) {
		return true
	}
	return false
}

func (b *versionedKVBackend) upgradeCheck(next framework.OperationFunc) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		if atomic.LoadUint32(b.upgrading) == 1 {
			// Sleep for a very short time before returning. This helps clients
			// that are trying to access a mount immediately upon enabling be
			// more likely to behave correctly since the operation should take
			// almost no time.
			time.Sleep(15 * time.Millisecond)

			if atomic.LoadUint32(b.upgrading) == 1 {
				if b.perfSecondaryCheck() {
					return logical.ErrorResponse("Waiting for the primary to upgrade from non-versioned to versioned data. This backend will be unavailable for a brief period and will resume service when the primary is finished."), logical.ErrInvalidRequest
				} else {
					return logical.ErrorResponse("Upgrading from non-versioned to versioned data. This backend will be unavailable for a brief period and will resume service shortly."), logical.ErrInvalidRequest
				}
			}
		}

		return next(ctx, req, data)
	}
}

func (b *versionedKVBackend) upgradeDone(ctx context.Context, s logical.Storage) (bool, error) {
	upgradeEntry, err := s.Get(ctx, path.Join(b.storagePrefix, "upgrading"))
	if err != nil {
		return false, err
	}

	var upgradeInfo UpgradeInfo
	if upgradeEntry != nil {
		err := proto.Unmarshal(upgradeEntry.Value, &upgradeInfo)
		if err != nil {
			return false, err
		}
	}

	return upgradeInfo.Done, nil
}

func (b *versionedKVBackend) Upgrade(ctx context.Context, s logical.Storage) error {
	replState := b.System().ReplicationState()

	// Don't run if the plugin is in metadata mode.
	if pluginutil.InMetadataMode() {
		b.Logger().Info("upgrade not running while plugin is in metadata mode")
		return nil
	}

	// Don't run while on a DR secondary.
	if replState.HasState(consts.ReplicationDRSecondary) {
		b.Logger().Info("upgrade not running on disaster recovery replication secondary")
		return nil
	}

	if !atomic.CompareAndSwapUint32(b.upgrading, 0, 1) {
		return errors.New("upgrade already in process")
	}

	// If we are a replication secondary or performance standby, wait until the primary has finished
	// upgrading.
	if b.perfSecondaryCheck() {
		b.Logger().Info("upgrade not running on performance replication secondary or performance standby")

		go func() {
			for {
				time.Sleep(time.Second)

				// If we failed because the context is closed we are
				// shutting down. Close this go routine and set the upgrade
				// flag back to 0 for good measure.
				if ctx.Err() != nil {
					atomic.StoreUint32(b.upgrading, 0)
					return
				}

				done, err := b.upgradeDone(ctx, s)
				if err != nil {
					b.Logger().Error("upgrading resulted in error", "error", err)
				}

				if done {
					break
				}
			}

			atomic.StoreUint32(b.upgrading, 0)
		}()

		return nil
	}

	// If we have 0 keys, it's either a new mount or one that's trivial to upgrade,
	// so we should do the upgrade synchronously
	upgradeSynchronously := false
	keys, err := logical.CollectKeys(ctx, s)
	if err != nil {
		b.Logger().Error("upgrading resulted in error", "error", err)
		return err
	}
	if len(keys) == 0 {
		upgradeSynchronously = true
	}

	upgradeInfo := &UpgradeInfo{
		StartedTime: ptypes.TimestampNow(),
	}

	// Encode the canary
	info, err := proto.Marshal(upgradeInfo)
	if err != nil {
		return err
	}

	// Because this is a long-running process we need a new context.
	ctx = context.Background()

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

	prepareUpgradeInfoDoneFunc := func() ([]byte, error) {
		upgradeInfo.Done = true
		info, err := proto.Marshal(upgradeInfo)
		if err != nil {
			b.Logger().Error("encoding upgrade info resulted in an error", "error", err)
			return nil, err
		}
		return info, nil
	}

	writeUpgradeInfoDoneFunc := func(info []byte) {
		for {
			err = s.Put(ctx, &logical.StorageEntry{
				Key:   path.Join(b.storagePrefix, "upgrading"),
				Value: info,
			})
			switch {
			case err == nil:
				return
			case err.Error() == logical.ErrSetupReadOnly.Error():
				time.Sleep(10 * time.Millisecond)
			default:
				b.Logger().Error("writing upgrade info resulted in an error, but all keys were successfully upgraded", "error", err)
				return
			}
		}
	}

	upgradeFunc := func() {
		// Write the canary value and if we are read only wait until the setup
		// process has finished.
	READONLY_LOOP:
		for {
			err := s.Put(ctx, &logical.StorageEntry{
				Key:   path.Join(b.storagePrefix, "upgrading"),
				Value: info,
			})
			switch {
			case err == nil:
				break READONLY_LOOP
			case err.Error() == logical.ErrSetupReadOnly.Error():
				time.Sleep(10 * time.Millisecond)
			default:
				b.Logger().Error("writing upgrade info resulted in an error", "error", err)
				return
			}
		}

		b.Logger().Info("collecting keys to upgrade")
		keys, err := logical.CollectKeys(ctx, s)
		if err != nil {
			b.Logger().Error("upgrading resulted in error", "error", err)
			return
		}

		b.Logger().Info("done collecting keys", "num_keys", len(keys))
		for i, key := range keys {
			if b.Logger().IsDebug() && i%500 == 0 {
				b.Logger().Debug("upgrading keys", "progress", fmt.Sprintf("%d/%d", i, len(keys)))
			}
			err := upgradeKey(key)
			if err != nil {
				b.Logger().Error("upgrading resulted in error", "error", err, "progress", fmt.Sprintf("%d/%d", i+1, len(keys)))
				return
			}
		}

		b.Logger().Info("upgrading keys finished")

		// We do this now so that we ensure it's written by the primary before
		// secondaries unblock
		b.l.Lock()
		if _, err = b.policy(ctx, s); err != nil {
			b.Logger().Error("error checking/creating policy after upgrade", "error", err)
		}
		b.l.Unlock()

		info, err := prepareUpgradeInfoDoneFunc()
		if err != nil {
			b.Logger().Error("error marshalling upgrade info after upgrade", "error", err)
			return
		}
		writeUpgradeInfoDoneFunc(info)
		atomic.StoreUint32(b.upgrading, 0)
	}

	if upgradeSynchronously {
		// Set us to having 'upgraded' before we insert the upgrade value, as the mount is ready to use now
		atomic.StoreUint32(b.upgrading, 0)
		info, err := prepareUpgradeInfoDoneFunc()
		if err != nil {
			return err
		}
		// We write the upgrade done info into storage in a goroutine, as a Vault mount is set to read only
		// during the mount process, so we cannot do it now
		go writeUpgradeInfoDoneFunc(info)
	} else {
		// We run the actual upgrade in a go routine, so we don't block the client on a
		// potentially long process.
		go upgradeFunc()
	}

	return nil
}
