// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rabbitmq

import (
	"fmt"
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	// Default interval to check the queue for items needing rotation
	defaultQueueTickSeconds = 5

	// Config key to set an alternate interval
	queueTickIntervalKey = "rotation_queue_tick_interval"

	// WAL storage key used for static account rotations
	staticWALKey = "staticRotationKey"
)

// rotateExpiredStaticCreds will pop expired credentials (credentials whose priority
// represents a time before the present), rotate the associated credential, and push
// them back onto the queue with the new priority.
func (b *backend) rotateExpiredStaticCreds(ctx context.Context, req *logical.Request) error {
	var errs *multierror.Error

	for {
		keepGoing, err := b.rotateCredential(ctx, req.Storage)
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

func (b *backend) rotateCredential(ctx context.Context, storage logical.Storage) (rotated bool, err error) {
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

	err = b.createStaticCredential(ctx, storage, cfg, true)
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

func (b *backend) createStaticCredential(ctx context.Context, storage logical.Storage, cfg staticRoleEntry, shouldLockStorage bool) error {
	return nil
}

func (b *backend) deleteStaticCredential(ctx context.Context, storage logical.Storage, cfg staticRoleEntry, shouldLockStorage bool) error {
	return nil
}

type setCredentialsWAL struct {
	NewPassword    string            `json:"new_password"`
	RoleName       string            `json:"role_name"`
	Username       string            `json:"username"`

	LastVaultRotation time.Time `json:"last_vault_rotation"`

	// Private fields which will not be included in json.Marshal/Unmarshal.
	walID        string
	walCreatedAt int64 // Unix time at which the WAL was created.
}
