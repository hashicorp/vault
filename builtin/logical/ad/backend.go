package ad

import (
	"context"

	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/builtin/logical/ad/creds"
	"github.com/hashicorp/vault/builtin/logical/ad/roles"
	"github.com/hashicorp/vault/builtin/logical/ad/util"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {

	configHandler := config.Handler(conf.Logger)

	roleHandler := roles.Handler(conf.Logger, configHandler)

	credsHandler := creds.Handler(conf.Logger, configHandler, roleHandler)

	roleHandler.AddDeleteWatcher(credsHandler)

	b := &framework.Backend{
		Paths: []*framework.Path{
			configHandler.Path(),
			roleHandler.Path(),
			credsHandler.Path(),
		},
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				config.BackendPath,
				creds.BackendPath,
			},
		},
		Invalidate:  util.Invalidator(configHandler.Invalidate, roleHandler.Invalidate, credsHandler.Invalidate).Invalidate,
		BackendType: logical.TypeLogical,
	}

	b.Setup(ctx, conf)

	return b, nil
}
