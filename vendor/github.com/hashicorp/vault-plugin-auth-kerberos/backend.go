package kerberos

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	configPath string = "config"
)

type backend struct {
	*framework.Backend
}

func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	b := &backend{}

	b.Backend = &framework.Backend{
		BackendType: logical.TypeCredential,
		Help:        backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{"login"},
			SealWrapStorage: []string{configPath},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				b.pathConfig(),
				b.pathConfigLdap(),
				b.pathLogin(),
				b.pathGroups(),
				b.pathGroupsList(),
			},
		),
	}

	return b
}

func (b *backend) config(ctx context.Context, s logical.Storage) (*kerberosConfig, error) {
	raw, err := s.Get(ctx, configPath)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	conf := &kerberosConfig{}
	if err := json.Unmarshal(raw.Value, conf); err != nil {
		return nil, err
	}

	return conf, nil
}

var backendHelp string = `
The Kerberos Auth Backend allows authentication via Kerberos SPNEGO.
`
