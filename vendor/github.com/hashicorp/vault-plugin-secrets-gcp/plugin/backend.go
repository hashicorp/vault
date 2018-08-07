package gcpsecrets

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/vault-plugin-secrets-gcp/plugin/iamutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/oauth2"
	"google.golang.org/api/iam/v1"
)

type backend struct {
	*framework.Backend

	iamResources iamutil.IamResourceParser

	rolesetLock sync.Mutex
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b = backend{
		iamResources: iamutil.GetEnabledIamResources(),
	}

	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				"config",
			},
		},

		Paths: framework.PathAppend(
			pathsRoleSet(&b),
			[]*framework.Path{
				pathConfig(&b),
				pathSecretAccessToken(&b),
				pathSecretServiceAccountKey(&b),
			},
		),
		Secrets: []*framework.Secret{
			secretAccessToken(&b),
			secretServiceAccountKey(&b),
		},

		BackendType:       logical.TypeLogical,
		WALRollback:       b.walRollback,
		WALRollbackMinAge: 5 * time.Minute,
	}

	return &b
}

func newHttpClient(ctx context.Context, s logical.Storage, scopes ...string) (*http.Client, error) {
	if len(scopes) == 0 {
		scopes = []string{"https://www.googleapis.com/auth/cloud-platform"}
	}

	cfg, err := getConfig(ctx, s)
	if err != nil {
		return nil, err
	}
	credsJSON := ""
	if cfg != nil {
		credsJSON = cfg.CredentialsRaw
	}

	_, tokenSource, err := gcputil.FindCredentials(credsJSON, ctx, scopes...)
	if err != nil {
		return nil, err
	}

	tc := cleanhttp.DefaultClient()
	return oauth2.NewClient(
		context.WithValue(ctx, oauth2.HTTPClient, tc),
		tokenSource), nil
}

func newIamAdmin(ctx context.Context, s logical.Storage) (*iam.Service, error) {
	c, err := newHttpClient(ctx, s, iam.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	return iam.New(c)
}

const backendHelp = `
The GCP secrets backend dynamically generates GCP IAM service
account keys with a given set of IAM policies. The service
account keys have a configurable lease set and are automatically
revoked at the end of the lease.

After mounting this backend, credentials to generate IAM keys must
be configured with the "config/" endpoints and policies must be
written using the "roles/" endpoints before any keys can be generated.
`
