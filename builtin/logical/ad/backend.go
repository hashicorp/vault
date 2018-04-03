package ad

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(_ context.Context, conf *logical.BackendConfig) (logical.Backend, error) {

	confHandler := &configHandler{conf.Logger}

	return &framework.Backend{
		Paths: []*framework.Path{
			confHandler.Path(),
		},
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config",
			},
		},
		BackendType: logical.TypeLogical,
	}, nil
}
