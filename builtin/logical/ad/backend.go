package ad

import (
	"context"

	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {

	confManager, err := config.NewManager(ctx, conf)
	if err != nil {
		return nil, err
	}

	b := &framework.Backend{
		Paths: []*framework.Path{
			confManager.Path(),
		},
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				config.BackendPath,
			},
		},
		Invalidate:  confManager.Invalidate,
		BackendType: logical.TypeLogical,
	}

	b.Setup(ctx, conf)

	return b, nil
}
