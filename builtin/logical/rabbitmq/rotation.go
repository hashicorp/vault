// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rabbitmq

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

const (
	// Default interval to check the queue for items needing rotation
	defaultQueueTickSeconds = 5

	// Config key to set an alternate interval
	// TODO make this configurable in the backend
	queueTickIntervalKey = "rotation_queue_tick_interval"

	// WAL storage key used for static account rotations
	staticWALKey = "staticRotationKey"
)

// rotateExpiredStaticCreds will pop expired credentials (credentials whose priority
// represents a time before the present), rotate the associated credential, and push
// them back onto the queue with the new priority.
func (b *backend) rotateExpiredStaticCreds(ctx context.Context, s logical.Storage) error {
	var errs *multierror.Error

	for {
		keepGoing, err := b.rotateCredential(ctx, s)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		if !keepGoing {
			if errs.ErrorOrNil() != nil {
				return fmt.Errorf("error(s) occurred while rotating expired static credentials: %w", errs)
			} else {
				return nil
			}
		}
	}
}

func (b *backend) validateRotationPeriod(period time.Duration) error {
	if period < defaultQueueTickSeconds {
		return fmt.Errorf("role rotation period out of range: must be greater than %d seconds", defaultQueueTickSeconds)
	}
	return nil
}

// TODO check if not pop by key needed
func (b *backend) rotateCredential(ctx context.Context, storage logical.Storage) (rotated bool, err error) {
	select {
	case <-ctx.Done():
		return false, nil
	default:
	}
	item, err := b.credRotationQueue.Pop()
	if err != nil {
		if err == queue.ErrEmpty {
			return false, nil
		}
		return false, fmt.Errorf("failed to pop from queue for role %q: %w", item.Key, err)
	}
	if item.Priority > time.Now().Unix() {
		err = b.credRotationQueue.Push(item)
		if err != nil {
			return false, fmt.Errorf("failed to add item into the rotation queue for username %q: %w", item.Key, err)
		}
		return false, nil
	}

	cfg := item.Value.(staticRoleEntry)

	err = b.createStaticCredential(ctx, storage, &cfg, item.Key)
	if err != nil {
		return false, err
	}

	// set new priority and re-queue
	item.Priority = time.Now().Add(cfg.RotationPeriod).Unix()
	err = b.credRotationQueue.Push(item)
	if err != nil {
		return false, fmt.Errorf("failed to add item into the rotation queue for username %q: %w", cfg.Username, err)
	}

	return true, nil
}

// TODO add option to lock storage while performing updates
func (b *backend) createStaticCredential(ctx context.Context, storage logical.Storage, cfg *staticRoleEntry, entryName string) error {
	config, err := readConfig(ctx, storage)
	if err != nil {
		return fmt.Errorf("unable to read configuration: %w", err)
	}
	newPassword, err := b.generatePassword(ctx, config.PasswordPolicy)
	if err != nil {
		return err
	}
	cfg.Password = newPassword
	client, err := b.Client(ctx, storage)
	if err != nil {
		return err
	}
	if _, err = client.DeleteUser(cfg.Username); err != nil {
		return fmt.Errorf("could not delete user: %w", err)
	}
	_, err = client.PutUser(cfg.Username, rabbithole.UserSettings{
		Password: cfg.Password,
		Tags:     []string{cfg.RoleEntry.Tags},
	})
	if err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}
	// update storage with new password and new rotation
	cfg.LastVaultRotation = time.Now()
	entry, err := logical.StorageEntryJSON(rabbitMQStaticRolePath+entryName, cfg)
	if err != nil {
		return err
	}
	if err := storage.Put(ctx, entry); err != nil {
		return err
	}
	// TODO: refactor host and topic setting function from path_role_create and set permissions here
	return nil
}

func (b *backend) deleteStaticCredential(ctx context.Context, storage logical.Storage, cfg staticRoleEntry, shouldLockStorage bool) error {
	if cfg.RevokeUserOnDelete {
		client, err := b.Client(ctx, storage)
		if err != nil {
			return err
		}
		if _, err = client.DeleteUser(cfg.Username); err != nil {
			return fmt.Errorf("could not delete user: %w", err)
		}
	}
	return nil
}

func (b *backend) initQueue(ctx context.Context, conf *logical.BackendConfig, replicationState consts.ReplicationState) {
	if (conf.System.LocalMount() || !replicationState.HasState(consts.ReplicationPerformanceSecondary)) &&
		!replicationState.HasState(consts.ReplicationDRSecondary) &&
		!replicationState.HasState(consts.ReplicationPerformanceStandby) {
		b.Logger().Info("initializing rabbitmq rotation queue")
		queueTickerInterval := defaultQueueTickSeconds * time.Second
		if strVal, ok := conf.Config[queueTickIntervalKey]; ok {
			newVal, err := strconv.Atoi(strVal)
			if err == nil {
				queueTickerInterval = time.Duration(newVal) * time.Second
			} else {
				b.Logger().Error("bad value for %q option: %q", queueTickIntervalKey, strVal)
			}
		}
		go b.runTicker(ctx, queueTickerInterval, conf.StorageView)
	}
}

func (b *backend) popFromRotationQueueByKey(name string) (*queue.Item, error) {
	select {
	case <-b.queueCtx.Done():
	default:
		item, err := b.credRotationQueue.PopByKey(name)
		if err != nil {
			return nil, err
		}
		if item != nil {
			return item, nil
		}
	}
	return nil, queue.ErrEmpty
}
