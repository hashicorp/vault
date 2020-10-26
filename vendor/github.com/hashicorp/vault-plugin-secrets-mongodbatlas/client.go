package mongodbatlas

import (
	"context"
	"errors"

	"github.com/Sectorbob/mlab-ns2/gae/ns/digest"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

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
