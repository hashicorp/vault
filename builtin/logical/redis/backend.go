package redis

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil
}

func Backend() *backend {
	var b backend

	b.Backend = &framework.Backend{
		Paths: framework.PathAppend(
			pathConfig(&b),
			[]*framework.Path{
				pathRole(&b),
				pathCreds(&b),
			},
		),

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},

		Invalidate:  b.invalidate,
		BackendType: logical.TypeLogical,
	}

	return &b
}

type backend struct {
	*framework.Backend

	client *redis.Client
	lock   sync.Mutex
}

func (b *backend) Client(ctx context.Context, s logical.Storage) (*redis.Client, error) {
	client := b.client
	if client != nil {
		return client, nil
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	conf, err := getConfig(ctx, s)
	if err != nil {
		return nil, err
	}

	b.client, err = conf.Client()

	return b.client, err
}

func (b *backend) invalidate(_ context.Context, key string) {
	if key == "config" {
		b.lock.Lock()
		defer b.lock.Unlock()

		b.client = nil
	}
}
