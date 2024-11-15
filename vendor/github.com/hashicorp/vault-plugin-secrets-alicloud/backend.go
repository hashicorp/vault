// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package alicloud

import (
	"context"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// operationPrefixAliCloud is used as a prefix for OpenAPI operation id's.
const operationPrefixAliCloud = "ali-cloud"

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	sdkConfig := sdk.NewConfig()
	sdkConfig.Scheme = "https"
	b := newBackend(sdkConfig)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// newBackend allows us to pass in the sdkConfig for testing purposes.
func newBackend(sdkConfig *sdk.Config) logical.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config",
			},
		},
		Paths: []*framework.Path{
			b.pathConfig(),
			b.pathRole(),
			b.pathListRoles(),
			b.pathCreds(),
		},
		Secrets: []*framework.Secret{
			b.pathSecrets(),
		},
		BackendType: logical.TypeLogical,
	}
	b.sdkConfig = sdkConfig
	return b
}

type backend struct {
	*framework.Backend
	sdkConfig *sdk.Config
}

const backendHelp = `
The AliCloud backend dynamically generates AliCloud access keys for a set of
RAM policies. The AliCloud access keys have a configurable ttl set and
are automatically revoked at the end of the ttl.

After mounting this backend, credentials to generate RAM keys must
be configured and roles must be written using
the "role/" endpoints before any access keys can be generated.
`
