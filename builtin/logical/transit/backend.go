package transit

import (
	"strings"

	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf)
	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend(conf *logical.BackendConfig) *backend {
	var b backend
	b.Backend = &framework.Backend{
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"archive/",
				"policy/",
			},
		},

		Paths: []*framework.Path{
			// Rotate/Config needs to come before Keys
			// as the handler is greedy
			b.pathConfig(),
			b.pathRotate(),
			b.pathRewrap(),
			b.pathKeys(),
			b.pathListKeys(),
			b.pathExportKeys(),
			b.pathEncrypt(),
			b.pathDecrypt(),
			b.pathDatakey(),
			b.pathRandom(),
			b.pathHash(),
			b.pathHMAC(),
			b.pathSign(),
			b.pathVerify(),
		},

		Secrets:     []*framework.Secret{},
		Invalidate:  b.invalidate,
		BackendType: logical.TypeLogical,
	}

	b.lm = keysutil.NewLockManager(conf.System.CachingDisabled())

	return &b
}

type backend struct {
	*framework.Backend
	lm *keysutil.LockManager
}

func (b *backend) invalidate(key string) {
	if b.Logger().IsTrace() {
		b.Logger().Trace("transit: invalidating key", "key", key)
	}
	switch {
	case strings.HasPrefix(key, "policy/"):
		name := strings.TrimPrefix(key, "policy/")
		b.lm.InvalidatePolicy(name)
	}
}
