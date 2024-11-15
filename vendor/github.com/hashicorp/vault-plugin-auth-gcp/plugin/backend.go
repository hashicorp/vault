// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpauth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/google/externalaccount"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"

	"github.com/hashicorp/vault-plugin-auth-gcp/plugin/cache"
)

const (
	userAgentPluginName = "auth-gcp"

	// operationPrefixGoogleCloud is used as a prefix for OpenAPI operation id's.
	operationPrefixGoogleCloud = "google-cloud"
)

// cacheTime is the duration for which to cache clients and credentials. This
// must be less than 60 minutes.
var cacheTime = 30 * time.Minute

type GcpAuthBackend struct {
	*framework.Backend

	// cache is the internal client/object cache. Callers should never access the
	// cache directly.
	cache *cache.Cache

	// pluginEnv contains Vault version information. It is used in user-agent headers.
	pluginEnv *logical.PluginEnvironment
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
		InitializeFunc: b.initialize,
		Invalidate:     b.invalidate,
	}
	return b
}

func (b *GcpAuthBackend) initialize(ctx context.Context, _ *logical.InitializationRequest) error {
	pluginEnv, err := b.System().PluginEnv(ctx)
	if err != nil {
		return fmt.Errorf("failed to read plugin environment: %w", err)
	}
	b.pluginEnv = pluginEnv

	return nil
}

// IAMClient returns a new IAM client. This client talks to the IAM endpoint,
// for all things that are not signing JWTs. The SignJWT method in the IAM
// client has been deprecated, but other methods are still valid and supported.
//
// See: https://pkg.go.dev/google.golang.org/api@v0.45.0/iam/v1 and:
// https://cloud.google.com/iam/docs/migrating-to-credentials-api#iam-sign-jwt-go
//
// The client is cached.
func (b *GcpAuthBackend) IAMClient(ctx context.Context, s logical.Storage) (*iam.Service, error) {
	cfg, err := b.config(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("failed to get config while creating IAM client: %w", err)
	}

	opts, err := b.clientOptions(ctx, s, cfg.IAMCustomEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM client options: %w", err)
	}

	client, err := b.cache.Fetch("iam", cacheTime, func() (interface{}, error) {
		client, err := iam.NewService(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create IAM client: %w", err)
		}
		client.UserAgent = useragent.PluginString(b.pluginEnv, userAgentPluginName)

		return client, nil
	})
	if err != nil {
		return nil, err
	}

	return client.(*iam.Service), nil
}

// ComputeClient returns a new Compute client. The client is cached.
func (b *GcpAuthBackend) ComputeClient(ctx context.Context, s logical.Storage) (*compute.Service, error) {
	cfg, err := b.config(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("failed to get config while creating Compute client: %w", err)
	}

	opts, err := b.clientOptions(ctx, s, cfg.ComputeCustomEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create Compute client options: %w", err)
	}

	client, err := b.cache.Fetch("compute", cacheTime, func() (interface{}, error) {
		client, err := compute.NewService(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create Compute client: %w", err)
		}
		client.UserAgent = useragent.PluginString(b.pluginEnv, userAgentPluginName)

		return client, nil
	})
	if err != nil {
		return nil, err
	}

	return client.(*compute.Service), nil
}

// CRMClient returns a new Cloud Resource Manager client. The client is cached.
func (b *GcpAuthBackend) CRMClient(ctx context.Context, s logical.Storage) (*cloudresourcemanager.Service, error) {
	cfg, err := b.config(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("failed to get config while creating Cloud Resource Manager client: %w", err)
	}

	opts, err := b.clientOptions(ctx, s, cfg.CRMCustomEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Resource Manager client options: %w", err)
	}

	client, err := b.cache.Fetch("crm", cacheTime, func() (interface{}, error) {
		client, err := cloudresourcemanager.NewService(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create Cloud Resource Manager client: %w", err)
		}
		client.UserAgent = useragent.PluginString(b.pluginEnv, userAgentPluginName)

		return client, nil
	})
	if err != nil {
		return nil, err
	}

	return client.(*cloudresourcemanager.Service), nil
}

// clientOptions returns a new set of client options containing an http.Client and optional
// custom endpoint. The http.Client is authenticated using the provided credentials. The
// underlying http.Client is cached among all clients.
func (b *GcpAuthBackend) clientOptions(ctx context.Context, s logical.Storage, endpoint string) ([]option.ClientOption, error) {
	creds, err := b.credentials(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 http client: %w", err)
	}

	client, err := b.cache.Fetch("HTTPClient", cacheTime, func() (interface{}, error) {
		b.Logger().Debug("creating oauth2 http client")
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, cleanhttp.DefaultClient())
		return oauth2.NewClient(ctx, creds.TokenSource), nil
	})
	if err != nil {
		return nil, err
	}

	opts := []option.ClientOption{option.WithHTTPClient(client.(*http.Client))}
	if endpoint != "" {
		opts = append(opts, option.WithEndpoint(endpoint))
	}

	return opts, nil
}

// credentials returns the credentials which were specified in the
// configuration. If no credentials were given during configuration, this uses
// default application credentials. If no default application credentials are
// found, this function returns an error. The credentials are cached in-memory
// for performance.
func (b *GcpAuthBackend) credentials(ctx context.Context, s logical.Storage) (*google.Credentials, error) {
	creds, err := b.cache.Fetch("credentials", cacheTime, func() (interface{}, error) {
		b.Logger().Debug("loading credentials")

		config, err := b.config(ctx, s)
		if err != nil {
			return nil, err
		}

		// Get creds from the config
		credBytes, err := config.formatAndMarshalCredentials()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal credential JSON: %w", err)
		}

		// If credentials were provided, use those. Otherwise fall back to the
		// default application credentials.
		var creds *google.Credentials
		if len(credBytes) > 0 {
			creds, err = google.CredentialsFromJSON(ctx, credBytes, iam.CloudPlatformScope)
			if err != nil {
				return nil, fmt.Errorf("failed to parse credentials: %w", err)
			}
		} else if config.IdentityTokenAudience != "" {
			ts := &PluginIdentityTokenSupplier{
				sys:      b.System(),
				logger:   b.Logger(),
				audience: config.IdentityTokenAudience,
				ttl:      config.IdentityTokenTTL,
			}

			creds, err = b.GetExternalAccountConfig(config, ts).GetExternalAccountCredentials(ctx)
		} else {
			creds, err = google.FindDefaultCredentials(ctx, iam.CloudPlatformScope)
			if err != nil {
				return nil, fmt.Errorf("failed to get default credentials: %w", err)
			}
		}

		return creds, err
	})
	if err != nil {
		return nil, err
	}
	return creds.(*google.Credentials), nil
}

func (b *GcpAuthBackend) GetExternalAccountConfig(c *gcpConfig, ts *PluginIdentityTokenSupplier) *gcputil.ExternalAccountConfig {
	b.Logger().Debug("adding web identity token fetcher")
	cfg := &gcputil.ExternalAccountConfig{
		ServiceAccountEmail: c.ServiceAccountEmail,
		Audience:            c.IdentityTokenAudience,
		TTL:                 c.IdentityTokenTTL,
		TokenSupplier:       ts,
	}

	return cfg
}

type PluginIdentityTokenSupplier struct {
	sys      logical.SystemView
	logger   hclog.Logger
	audience string
	ttl      time.Duration
}

var _ externalaccount.SubjectTokenSupplier = (*PluginIdentityTokenSupplier)(nil)

func (p *PluginIdentityTokenSupplier) SubjectToken(ctx context.Context, opts externalaccount.SupplierOptions) (string, error) {
	p.logger.Info("fetching new plugin identity token")
	resp, err := p.sys.GenerateIdentityToken(ctx, &pluginutil.IdentityTokenRequest{
		Audience: p.audience,
		TTL:      p.ttl,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate plugin identity token: %w", err)
	}

	if resp.TTL < p.ttl {
		p.logger.Debug("generated plugin identity token has shorter TTL than requested",
			"requested", p.ttl.Seconds(), "actual", resp.TTL)
	}

	return resp.Token.Token(), nil
}

// ClearCaches deletes all cached clients and credentials.
func (b *GcpAuthBackend) ClearCaches() {
	b.cache.Clear()
}

// invalidate resets the plugin. This is called when a key is updated via
// replication.
func (b *GcpAuthBackend) invalidate(_ context.Context, key string) {
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
