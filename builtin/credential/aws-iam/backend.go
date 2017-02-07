package awsiam

import (
	"sync"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

func Backend() *backend {
	var b backend

	b.Backend = &framework.Backend{
		Help: backendHelp,
		PathsSpecial: &logical.Paths{
			Root: []string{},
			Unauthenticated: []string{
				"login",
			},
		},
		Paths: []*framework.Path{
			pathLogin(&b),
			pathConfigClient(&b),
			pathRole(&b),
			pathListRole(&b),
		},
		AuthRenew: b.pathLoginRenew,
	}

	return &b
}

type backend struct {
	*framework.Backend

	// Lock to make changes to the client config
	configMutex sync.RWMutex

	// Lock to make changes to the roles
	roleMutex sync.RWMutex
}

const backendHelp = `
The aws-iam auth backend provides a secure introduction mechanism for AWS IAM
principals. This allows you to reduce the problem of securely introducing a
Vault token to the problem of securely introducing AWS IAM credentials, which
AWS has already solved in a number of use cases, such as via EC2 Instance
Profiles and IAM roles attached to Lambda functions. It also allows you to have
a consistent workflow between developers working on a local laptop with tools
such as Hologram.
`
