package aws

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				"config/root",
			},
		},

		Paths: []*framework.Path{
			pathConfigRoot(&b),
			pathConfigRotateRoot(&b),
			pathConfigLease(&b),
			pathRoles(&b),
			pathListRoles(&b),
			pathUser(&b),
		},

		Secrets: []*framework.Secret{
			secretAccessKeys(&b),
		},

		WALRollback:       b.walRollback,
		WALRollbackMinAge: minAwsUserRollbackAge,
		BackendType:       logical.TypeLogical,
	}

	return &b
}

type backend struct {
	*framework.Backend

	// Mutex to protect access to reading and writing policies
	roleMutex sync.RWMutex

	// Mutex to protect access to iam/sts clients and client configs
	clientMutex sync.RWMutex

	// iamClient and stsClient hold configured iam and sts clients for reuse, and
	// to enable mocking with AWS iface for tests
	iamClient iamiface.IAMAPI
	stsClient stsiface.STSAPI
}

const backendHelp = `
The AWS backend dynamically generates AWS access keys for a set of
IAM policies. The AWS access keys have a configurable lease set and
are automatically revoked at the end of the lease.

After mounting this backend, credentials to generate IAM keys must
be configured with the "root" path and policies must be written using
the "roles/" endpoints before any access keys can be generated.
`

// clientIAM returns the configured IAM client. If nil, it constructs a new one
// and returns it, setting it the internal variable
func (b *backend) clientIAM(ctx context.Context, s logical.Storage) (iamiface.IAMAPI, error) {
	b.clientMutex.RLock()
	if b.iamClient != nil {
		b.clientMutex.RUnlock()
		return b.iamClient, nil
	}

	// Upgrade the lock for writing
	b.clientMutex.RUnlock()
	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	// check client again, in the event that a client was being created while we
	// waited for Lock()
	if b.iamClient != nil {
		return b.iamClient, nil
	}

	iamClient, err := nonCachedClientIAM(ctx, s)
	if err != nil {
		return nil, err
	}
	b.iamClient = iamClient

	return b.iamClient, nil
}

func (b *backend) clientSTS(ctx context.Context, s logical.Storage) (stsiface.STSAPI, error) {
	b.clientMutex.RLock()
	if b.stsClient != nil {
		b.clientMutex.RUnlock()
		return b.stsClient, nil
	}

	// Upgrade the lock for writing
	b.clientMutex.RUnlock()
	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	// check client again, in the event that a client was being created while we
	// waited for Lock()
	if b.stsClient != nil {
		return b.stsClient, nil
	}

	stsClient, err := nonCachedClientSTS(ctx, s)
	if err != nil {
		return nil, err
	}
	b.stsClient = stsClient

	return b.stsClient, nil
}

const minAwsUserRollbackAge = 5 * time.Minute
