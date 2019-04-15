package gcpauth

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault-plugin-auth-gcp/plugin/cache"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iam/v1"
)

var (
	// cacheTime is the duration for which to cache clients and credentials. This
	// must be less than 60 minutes.
	cacheTime = 30 * time.Minute
)

type GcpAuthBackend struct {
	*framework.Backend

	// cache is the internal client/object cache. Callers should never access the
	// cache directly.
	cache *cache.Cache
}

// Factory returns a new backend as logical.Backend.
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *GcpAuthBackend {
	b := &GcpAuthBackend{
		cache: cache.New(),
	}

	b.Backend = &framework.Backend{
		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
		Help:        backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
			SealWrapStorage: []string{
				"config",
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathConfig(b),
				pathLogin(b),
			},
			pathsRole(b),
		),

		Invalidate: b.invalidate,
	}
	return b
}

// IAMClient returns a new IAM client. The client is cached.
func (b *GcpAuthBackend) IAMClient(s logical.Storage) (*iam.Service, error) {
	httpClient, err := b.httpClient(s)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create IAM HTTP client: {{err}}", err)
	}

	client, err := b.cache.Fetch("iam", cacheTime, func() (interface{}, error) {
		client, err := iam.New(httpClient)
		if err != nil {
			return nil, errwrap.Wrapf("failed to create IAM client: {{err}}", err)
		}
		client.UserAgent = useragent.String()

		return client, nil
	})
	if err != nil {
		return nil, err
	}

	return client.(*iam.Service), nil
}

// ComputeClient returns a new Compute client. The client is cached.
func (b *GcpAuthBackend) ComputeClient(s logical.Storage) (*compute.Service, error) {
	httpClient, err := b.httpClient(s)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create Compute HTTP client: {{err}}", err)
	}

	client, err := b.cache.Fetch("compute", cacheTime, func() (interface{}, error) {
		client, err := compute.New(httpClient)
		if err != nil {
			return nil, errwrap.Wrapf("failed to create Compute client: {{err}}", err)
		}
		client.UserAgent = useragent.String()

		return client, nil
	})
	if err != nil {
		return nil, err
	}

	return client.(*compute.Service), nil
}

// CRMClient returns a new Cloud Resource Manager client. The client is cached.
func (b *GcpAuthBackend) CRMClient(s logical.Storage) (*cloudresourcemanager.Service, error) {
	httpClient, err := b.httpClient(s)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create Cloud Resource Manager HTTP client: {{err}}", err)
	}

	client, err := b.cache.Fetch("crm", cacheTime, func() (interface{}, error) {
		client, err := cloudresourcemanager.New(httpClient)
		if err != nil {
			return nil, errwrap.Wrapf("failed to create Cloud Resource Manager client: {{err}}", err)
		}
		client.UserAgent = useragent.String()

		return client, nil
	})
	if err != nil {
		return nil, err
	}

	return client.(*cloudresourcemanager.Service), nil
}

// httpClient returns a new http.Client that is authenticated using the provided
// credentials. The underlying httpClient is cached among all clients.
func (b *GcpAuthBackend) httpClient(s logical.Storage) (*http.Client, error) {
	creds, err := b.credentials(s)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create oauth2 http client: {{err}}", err)
	}

	client, err := b.cache.Fetch("HTTPClient", cacheTime, func() (interface{}, error) {
		b.Logger().Debug("creating oauth2 http client")
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, cleanhttp.DefaultClient())
		return oauth2.NewClient(ctx, creds.TokenSource), nil
	})
	if err != nil {
		return nil, err
	}

	return client.(*http.Client), nil
}

// credentials returns the credentials which were specified in the
// configuration. If no credentials were given during configuration, this uses
// default application credentials. If no default application credentials are
// found, this function returns an error. The credentials are cached in-memory
// for performance.
func (b *GcpAuthBackend) credentials(s logical.Storage) (*google.Credentials, error) {
	creds, err := b.cache.Fetch("credentials", cacheTime, func() (interface{}, error) {
		b.Logger().Debug("loading credentials")

		ctx := context.Background()

		config, err := b.config(ctx, s)
		if err != nil {
			return nil, err
		}

		// Get creds from the config
		credBytes, err := config.formatAndMarshalCredentials()
		if err != nil {
			return nil, errwrap.Wrapf("failed to marshal credential JSON: {{err}}", err)
		}

		// If credentials were provided, use those. Otherwise fall back to the
		// default application credentials.
		var creds *google.Credentials
		if len(credBytes) > 0 {
			creds, err = google.CredentialsFromJSON(ctx, credBytes, iam.CloudPlatformScope)
			if err != nil {
				return nil, errwrap.Wrapf("failed to parse credentials: {{err}}", err)
			}
		} else {
			creds, err = google.FindDefaultCredentials(ctx, iam.CloudPlatformScope)
			if err != nil {
				return nil, errwrap.Wrapf("failed to get default credentials: {{err}}", err)
			}
		}

		return creds, err
	})
	if err != nil {
		return nil, err
	}
	return creds.(*google.Credentials), nil
}

// ClearCaches deletes all cached clients and credentials.
func (b *GcpAuthBackend) ClearCaches() {
	b.cache.Clear()
}

// invalidate resets the plugin. This is called when a key is updated via
// replication.
func (b *GcpAuthBackend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config":
		b.ClearCaches()
	}
}

const backendHelp = `
The GCP auth method allows machines to authenticate Google Cloud Platform
entities. It supports two modes of authentication:

- IAM service accounts: provides a signed JSON Web Token for a given
  service account key

- GCE VM metadata: provides a signed JSON Web Token using instance metadata
  obtained from the GCE instance metadata server
`
