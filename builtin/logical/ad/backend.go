package ad

import (
	"context"

	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/builtin/logical/ad/creds"
	"github.com/hashicorp/vault/builtin/logical/ad/roles"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {

	confManager, err := config.NewManager(ctx, conf)
	if err != nil {
		return nil, err
	}

	roleManager := roles.NewManager(conf.Logger, confManager)

	credsManager := creds.NewManager(conf.Logger, confManager, roleManager)

	invalidator := newInvalidationMgr(confManager.Invalidate, roleManager.Invalidate, credsManager.Invalidate)

	b := &framework.Backend{
		Paths: []*framework.Path{
			confManager.Path(),
			roleManager.Path(),
			credsManager.Path(),
		},
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				config.BackendPath,
			},
		},
		Invalidate:  invalidator.invalidate,
		BackendType: logical.TypeLogical,
	}

	b.Setup(ctx, conf)

	return b, nil
}

func newInvalidationMgr(invalidationFuncs ...framework.InvalidateFunc) *invalidationMgr {
	return &invalidationMgr{invalidationFuncs}
}

type invalidationMgr struct {
	toCall []framework.InvalidateFunc
}

func (v *invalidationMgr) invalidate(ctx context.Context, key string) {
	for _, f := range v.toCall {
		f(ctx, key)
	}
}
