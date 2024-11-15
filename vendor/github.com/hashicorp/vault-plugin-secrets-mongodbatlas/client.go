// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongodbatlas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mongodb-forks/digest"
	"go.mongodb.org/atlas/mongodbatlas"
)

const userAgentPluginName = "mongodbatlas-secrets"

func (b *Backend) clientMongo(ctx context.Context, s logical.Storage) (*mongodbatlas.Client, error) {
	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	// if the client is already created, just return it
	if b.client != nil {
		return b.client, nil
	}

	client, err := nonCachedClient(ctx, s)
	if err != nil {
		return nil, err
	}

	pluginEnv, err := b.system.PluginEnv(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin environment: %w", err)
	}
	client.UserAgent = useragent.PluginString(pluginEnv, userAgentPluginName)

	b.client = client

	return b.client, nil
}

func nonCachedClient(ctx context.Context, s logical.Storage) (*mongodbatlas.Client, error) {
	config, err := getRootConfig(ctx, s)
	if err != nil {
		return nil, err
	}

	transport := digest.NewTransport(config.PublicKey, config.PrivateKey)

	client, err := transport.Client()
	if err != nil {
		return nil, err
	}

	return mongodbatlas.NewClient(client), nil
}

func getRootConfig(ctx context.Context, s logical.Storage) (*config, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}
	if entry != nil {
		var config config
		if err := entry.DecodeJSON(&config); err != nil {
			return nil, errwrap.Wrapf("error reading root configuration: {{err}}", err)
		}

		// return the config, we are done
		return &config, nil

	}

	return nil, errors.New("empty config entry")
}
