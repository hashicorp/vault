package pcf

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// These env vars are used frequently to pull the client certificate and private key
	// from PCF containers; thus are placed here for ease of discovery and use from
	// outside packages.
	EnvVarInstanceCertificate = "CF_INSTANCE_CERT"
	EnvVarInstanceKey         = "CF_INSTANCE_KEY"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := &backend{}
	b.Backend = &framework.Backend{
		AuthRenew: b.pathLoginRenew,
		Help:      backendHelp,
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{"config"},
			Unauthenticated: []string{"login"},
		},
		Paths: []*framework.Path{
			b.pathConfig(),
			b.pathListRoles(),
			b.pathRoles(),
			b.pathLogin(),
		},
		BackendType: logical.TypeCredential,
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The PCF auth backend supports logging in using PCF's identity service.
Once a CA certificate is configured, and Vault is configured to consume
PCF's API, PCF's instance identity credentials can be used to authenticate.'
`
