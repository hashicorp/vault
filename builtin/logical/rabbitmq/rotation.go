// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rabbitmq

import (
	"context"
	"time"
	"strconv"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// Default interval to check the queue for items needing rotation
	defaultQueueTickSeconds = 5

	// Config key to set an alternate interval
	queueTickIntervalKey = "rotation_queue_tick_interval"

	// WAL storage key used for static account rotations
	staticWALKey = "staticRotationKey"
)

// initQueue preforms the necessary checks and initializations needed to perform
// automatic credential rotation for roles associated with static accounts. This
// method verifies if a queue is needed (primary server or local mount), and if
// so initializes the queue and launches a go-routine to periodically invoke a
// method to preform the rotations.
//
// initQueue is invoked by the Factory method in a go-routine. The Factory does
// not wait for success or failure of it's tasks before continuing. This is to
// avoid blocking the mount process while loading and evaluating existing roles,
// etc.
func (b *backend) initQueue(ctx context.Context, conf *logical.BackendConfig, replicationState consts.ReplicationState) {
	// Verify this mount is on the primary server, or is a local mount. If not, do
	// not create a queue or launch a ticker. Both processing the WAL list and
	// populating the queue are done sequentially and before launching a
	// go-routine to run the periodic ticker.
	if (conf.System.LocalMount() || !replicationState.HasState(consts.ReplicationPerformanceSecondary)) &&
		!replicationState.HasState(consts.ReplicationDRSecondary) &&
		!replicationState.HasState(consts.ReplicationPerformanceStandby) {
		b.Logger().Info("initializing database rotation queue")

		// Poll for a PutWAL call that does not return a "read-only storage" error.
		// This ensures the startup phases of loading WAL entries from any possible
		// failed rotations can complete without error when deleting from storage.
	READONLY_LOOP:
		for {
			select {
			case <-ctx.Done():
				b.Logger().Info("queue initialization canceled")
				return
			default:
			}

			walID, err := framework.PutWAL(ctx, conf.StorageView, staticWALKey, &setCredentialsWAL{RoleName: "vault-readonlytest"})
			if walID != "" && err == nil {
				defer framework.DeleteWAL(ctx, conf.StorageView, walID)
			}
			switch {
			case err == nil:
				break READONLY_LOOP
			case err.Error() == logical.ErrSetupReadOnly.Error():
				time.Sleep(10 * time.Millisecond)
			default:
				b.Logger().Error("deleting nil key resulted in error", "error", err)
				return
			}
		}

		// Load roles and populate queue with static accounts
		b.populateQueue(ctx, conf.StorageView)

		// Launch ticker
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

func (b *backend) populateQueue(ctx context.Context, s logical.Storage) {
	log := b.Logger()
	log.Info("populating role rotation queue")
	// TODO implement
}

// runTicker kicks off a periodic ticker that invoke the automatic credential
// rotation method at a determined interval. The default interval is 5 seconds.
func (b *backend) runTicker(ctx context.Context, queueTickInterval time.Duration, s logical.Storage) {
	b.Logger().Info("starting periodic ticker")
	tick := time.NewTicker(queueTickInterval)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			b.rotateCredentials(ctx, s)

		case <-ctx.Done():
			b.Logger().Info("stopping periodic ticker")
			return
		}
	}
}

func (b *backend) rotateCredentials(ctx context.Context, s logical.Storage) {
	for b.rotateCredential(ctx, s) {
	}
}

func (b *backend) rotateCredential(ctx context.Context, s logical.Storage) bool {
	// TODO implement
	return false
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
