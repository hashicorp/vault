// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rabbitmq

import (
	"context"
	"strings"
	"sync"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

const (
	operationPrefixRabbitMQ = "rabbit-mq"
	rabbitMQRolePath        = "role/"
	rabbitMQStaticRolePath  = "static-role/"
)

// Factory creates and configures the backend
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	b.credRotationQueue = queue.New()
	go b.initQueue(b.queueCtx, conf, conf.System.ReplicationState())
	return b, nil
}

// Creates a new backend with all the paths and secrets belonging to it
func Backend(conf *logical.BackendConfig) *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config/connection",
			},
		},

		Paths: framework.PathAppend(
			[]*framework.Path{
				pathConfigConnection(&b),
				pathConfigLease(&b),
			},
			pathCreds(&b),
			pathListRoles(&b),
			pathRoles(&b),
		),

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},

		Clean:      b.resetClient,
		Invalidate: b.invalidate,
		PeriodicFunc: func(ctx context.Context, req *logical.Request) error {
			repState := conf.System.ReplicationState()
			if (conf.System.LocalMount() ||
				!repState.HasState(consts.ReplicationPerformanceSecondary)) &&
				!repState.HasState(consts.ReplicationDRSecondary) &&
				!repState.HasState(consts.ReplicationPerformanceStandby) {
				return b.rotateExpiredStaticCreds(ctx, req.Storage)
			}
			return nil
		},
		BackendType: logical.TypeLogical,
	}
	b.queueCtx, b.cancelQueueCtx = context.WithCancel(context.Background())
	return &b
}

type backend struct {
	*framework.Backend

	client *rabbithole.Client
	lock   sync.RWMutex

	credRotationQueue *queue.PriorityQueue
	queueCtx          context.Context
	cancelQueueCtx    context.CancelFunc
}

// DB returns the database connection.
func (b *backend) Client(ctx context.Context, s logical.Storage) (*rabbithole.Client, error) {
	b.lock.RLock()

	// If we already have a client, return it
	if b.client != nil {
		b.lock.RUnlock()
		return b.client, nil
	}

	b.lock.RUnlock()

	// Otherwise, attempt to make connection
	connConfig, err := readConfig(ctx, s)
	if err != nil {
		return nil, err
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	// If the client was created during the lock switch, return it
	if b.client != nil {
		return b.client, nil
	}

	b.client, err = rabbithole.NewClient(connConfig.URI, connConfig.Username, connConfig.Password)
	if err != nil {
		return nil, err
	}
	// Use a default pooled transport so there would be no leaked file descriptors
	b.client.SetTransport(cleanhttp.DefaultPooledTransport())

	return b.client, nil
}

// resetClient forces a connection next time Client() is called.
func (b *backend) resetClient(_ context.Context) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.client = nil
}

func (b *backend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config/connection":
		b.resetClient(ctx)
	}
}

func (b *backend) runTicker(ctx context.Context, queueTickInterval time.Duration, s logical.Storage) {
	b.Logger().Info("starting periodic ticker")
	tick := time.NewTicker(queueTickInterval)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			b.rotateExpiredStaticCreds(ctx, s)

		case <-ctx.Done():
			b.Logger().Info("stopping periodic ticker")
			return
		}
	}
}

// Lease returns the lease information
func (b *backend) Lease(ctx context.Context, s logical.Storage) (*configLease, error) {
	entry, err := s.Get(ctx, "config/lease")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result configLease
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

const backendHelp = `
The RabbitMQ backend dynamically generates RabbitMQ users.

After mounting this backend, configure it using the endpoints within
the "config/" path.
`
