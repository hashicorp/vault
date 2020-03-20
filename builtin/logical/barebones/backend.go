package barebones

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := &backend{}
	b.Backend = &framework.Backend{
		Paths: []*framework.Path{
			b.somePath(),
		},
		BackendType: logical.TypeLogical,
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

type backend struct {
	*framework.Backend
}

func (b *backend) somePath() *framework.Path {
	return &framework.Path{
		Pattern: "empty-call",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.emptyRead,
		},
	}
}

func (b *backend) emptyRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return &logical.Response{
		Data: map[string]interface{}{
			"hello": "world",
		},
	}, nil
}
