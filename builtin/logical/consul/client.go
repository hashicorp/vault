package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) client(ctx context.Context, s logical.Storage) (*api.Client, error, error) {
	conf, userErr, intErr := b.readConfigAccess(ctx, s)
	if intErr != nil {
		return nil, nil, intErr
	}
	if userErr != nil {
		return nil, userErr, nil
	}
	if conf == nil {
		return nil, nil, fmt.Errorf("no error received but no configuration found")
	}

	client, err := conf.Client()
	return client, nil, err
}
