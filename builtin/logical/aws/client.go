// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// getRootIAMConfig creates an *aws.Config for Vault to connect to IAM.
func (b *backend) getRootIAMConfig(ctx context.Context, s logical.Storage, logger hclog.Logger) (*aws.Config, error) {
	credsConfig := &awsutil.CredentialsConfig{}
	var endpoint string
	var maxRetries int = aws.UseServiceDefaultRetries

	entry, err := s.Get(ctx, "config/root")
	if err != nil {
		return nil, err
	}
	if entry != nil {
		var config rootConfig
		if err := entry.DecodeJSON(&config); err != nil {
			return nil, fmt.Errorf("error reading root configuration: %w", err)
		}

		credsConfig.AccessKey = config.AccessKey
		credsConfig.SecretKey = config.SecretKey
		credsConfig.Region = config.Region
		maxRetries = config.MaxRetries

		if config.IAMEndpoint != "" {
			endpoint = *aws.String(config.IAMEndpoint)
		}

		if config.IdentityTokenAudience != "" {
			ns, err := namespace.FromContext(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get namespace from context: %w", err)
			}

			fetcher := &PluginIdentityTokenFetcher{
				sys:      b.System(),
				logger:   b.Logger(),
				ns:       ns,
				audience: config.IdentityTokenAudience,
				ttl:      config.IdentityTokenTTL,
			}

			sessionSuffix := strconv.FormatInt(time.Now().UnixNano(), 10)
			credsConfig.RoleSessionName = fmt.Sprintf("vault-aws-secrets-%s", sessionSuffix)
			credsConfig.WebIdentityTokenFetcher = fetcher
			credsConfig.RoleARN = config.RoleARN
		}
	}

	if credsConfig.Region == "" {
		credsConfig.Region = getFallbackRegion()
	}

	credsConfig.HTTPClient = cleanhttp.DefaultClient()

	credsConfig.Logger = logger

	creds, err := credsConfig.GenerateCredentialChain()
	if err != nil {
		return nil, err
	}

	return &aws.Config{
		Credentials: creds,
		Region:      aws.String(credsConfig.Region),
		Endpoint:    &endpoint,
		HTTPClient:  cleanhttp.DefaultClient(),
		MaxRetries:  aws.Int(maxRetries),
	}, nil
}

// Return a slice of *aws.Config, based on descending configuration priority. STS endpoints are the only place this is used.
// NOTE: The caller is required to ensure that b.clientMutex is at least read locked
func (b *backend) getRootSTSConfigs(ctx context.Context, s logical.Storage, logger hclog.Logger) ([]*aws.Config, error) {
	// set fallback region (we can overwrite later)
	fallbackRegion := getFallbackRegion()

	maxRetries := aws.UseServiceDefaultRetries

	entry, err := s.Get(ctx, "config/root")
	if err != nil {
		return nil, err
	}
	var configs []*aws.Config

	// ensure the nil case uses defaults
	if entry == nil {
		ccfg := awsutil.CredentialsConfig{
			HTTPClient: cleanhttp.DefaultClient(),
			Logger:     logger,
			Region:     fallbackRegion,
		}
		creds, err := ccfg.GenerateCredentialChain()
		if err != nil {
			return nil, err
		}
		configs = append(configs, &aws.Config{
			Credentials: creds,
			Region:      aws.String(fallbackRegion),
			Endpoint:    aws.String(""),
			MaxRetries:  aws.Int(maxRetries),
		})

		return configs, nil
	}

	var config rootConfig
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, fmt.Errorf("error reading root configuration: %w", err)
	}

	var endpoints []string
	var regions []string
	credsConfig := &awsutil.CredentialsConfig{}

	credsConfig.AccessKey = config.AccessKey
	credsConfig.SecretKey = config.SecretKey
	credsConfig.HTTPClient = cleanhttp.DefaultClient()
	credsConfig.Logger = logger

	if config.Region != "" {
		regions = append(regions, config.Region)
	}

	maxRetries = config.MaxRetries
	if config.STSEndpoint != "" {
		endpoints = append(endpoints, config.STSEndpoint)
		if config.STSRegion != "" {
			// this retains original logic, where sts region was only used if sts endpoint was set
			regions = []string{config.STSRegion} // override to be "only" region if set
		}

		if len(config.STSFallbackEndpoints) > 0 {
			endpoints = append(endpoints, config.STSFallbackEndpoints...)
		}

		if len(config.STSFallbackRegions) > 0 {
			regions = append(regions, config.STSFallbackRegions...)
		}
	}

	opts := make([]awsutil.Option, 0)
	if config.IdentityTokenAudience != "" {
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get namespace from context: %w", err)
		}

		fetcher := &PluginIdentityTokenFetcher{
			sys:      b.System(),
			logger:   b.Logger(),
			ns:       ns,
			audience: config.IdentityTokenAudience,
			ttl:      config.IdentityTokenTTL,
		}

		sessionSuffix := strconv.FormatInt(time.Now().UnixNano(), 10)
		credsConfig.RoleSessionName = fmt.Sprintf("vault-aws-secrets-%s", sessionSuffix)
		credsConfig.WebIdentityTokenFetcher = fetcher
		credsConfig.RoleARN = config.RoleARN

		// explicitly disable environment and shared credential providers when using Web Identity Token Fetcher
		// enables WIF usage in environments that may use AWS Profiles or environment variables for other use-cases
		opts = append(opts, awsutil.WithEnvironmentCredentials(false), awsutil.WithSharedCredentials(false))
	}

	// at this point, in the IAM case,
	// - regions contains config.Region, if it was set.
	// - endpoints contains iam_endpoint, if it was set.
	// in the sts case,
	// - regions contains sts_region, if it was set, then sts_fallback_regions in order, if they were set.
	// - endpoints contains sts_endpoint, if it was set, then sts_fallback_endpoints in order, if they were set.

	// case in which nothing was supplied
	if len(regions) == 0 {
		// fallback region is in descending order, AWS_REGION, or AWS_DEFAULT_REGION, or us-east-1
		regions = append(regions, fallbackRegion)
	}

	if len(endpoints) == 0 {
		for _, v := range regions {
			endpoints = append(endpoints, matchingSTSEndpoint(v))
		}
	}

	// for this approach of using parallel arrays to part out the configs, we want equal numbers of regions and endpoints
	if len(regions) != len(endpoints) {
		return nil, errors.New("number of regions does not match number of endpoints")
	}

	for i := 0; i < len(endpoints); i++ {
		if len(regions) > i {
			credsConfig.Region = regions[i]
		} else {
			credsConfig.Region = fallbackRegion
		}
		creds, err := credsConfig.GenerateCredentialChain(opts...)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &aws.Config{
			Credentials: creds,
			Region:      aws.String(credsConfig.Region),
			Endpoint:    aws.String(endpoints[i]),
			MaxRetries:  aws.Int(maxRetries),
			HTTPClient:  cleanhttp.DefaultClient(),
		})
	}

	return configs, nil
}

func (b *backend) nonCachedClientIAM(ctx context.Context, s logical.Storage, logger hclog.Logger, entry *staticRoleEntry) (*iam.IAM, error) {
	var awsConfig *aws.Config
	var err error

	if entry != nil && entry.AssumeRoleARN != "" {
		awsConfig, err = b.assumeRoleStatic(ctx, s, entry)
		if err != nil {
			return nil, fmt.Errorf("failed to assume role %q: %w", entry.AssumeRoleARN, err)
		}
	} else {
		awsConfig, err = b.getRootIAMConfig(ctx, s, logger)
		if err != nil {
			return nil, err
		}
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	client := iam.New(sess)
	if client == nil {
		return nil, fmt.Errorf("could not obtain IAM client")
	}
	return client, nil
}

func (b *backend) nonCachedClientSTS(ctx context.Context, s logical.Storage, logger hclog.Logger) (*sts.STS, error) {
	awsConfig, err := b.getRootSTSConfigs(ctx, s, logger)
	if err != nil {
		return nil, err
	}

	var client *sts.STS

	for _, cfg := range awsConfig {
		sess, err := session.NewSession(cfg)
		if err != nil {
			return nil, err
		}
		client = sts.New(sess)
		if client == nil {
			return nil, fmt.Errorf("could not obtain sts client")
		}

		// ping the client - we only care about errors
		_, err = client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
		if err == nil {
			return client, nil
		} else {
			b.Logger().Debug("couldn't connect with config trying next", "failed endpoint", *cfg.Endpoint, "failed region", *cfg.Region)
		}
	}

	return nil, fmt.Errorf("could not obtain sts client")
}

// matchingSTSEndpoint returns the endpoint for the supplied region, according to
// http://docs.aws.amazon.com/general/latest/gr/sts.html
func matchingSTSEndpoint(stsRegion string) string {
	return fmt.Sprintf("https://sts.%s.amazonaws.com", stsRegion)
}

// getFallbackRegion returns an aws region fallback. It will check in the AWS specified order:
// - AWS_REGION, then
// - AWS_DEFAULT_REGION, then
// - us-east-1
func getFallbackRegion() string {
	// set fallback region (we can overwrite later)
	fallbackRegion := os.Getenv("AWS_REGION")
	if fallbackRegion == "" {
		fallbackRegion = os.Getenv("AWS_DEFAULT_REGION")
	}
	if fallbackRegion == "" {
		fallbackRegion = "us-east-1"
	}

	return fallbackRegion
}

// PluginIdentityTokenFetcher fetches plugin identity tokens from Vault. It is provided
// to the AWS SDK client to keep assumed role credentials refreshed through expiration.
// When the client's STS credentials expire, it will use this interface to fetch a new
// plugin identity token and exchange it for new STS credentials.
type PluginIdentityTokenFetcher struct {
	sys      logical.SystemView
	logger   hclog.Logger
	audience string
	ns       *namespace.Namespace
	ttl      time.Duration
}

var _ stscreds.TokenFetcher = (*PluginIdentityTokenFetcher)(nil)

func (f PluginIdentityTokenFetcher) FetchToken(ctx aws.Context) ([]byte, error) {
	nsCtx := namespace.ContextWithNamespace(ctx, f.ns)
	resp, err := f.sys.GenerateIdentityToken(nsCtx, &pluginutil.IdentityTokenRequest{
		Audience: f.audience,
		TTL:      f.ttl,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate plugin identity token: %w", err)
	}
	f.logger.Info("fetched new plugin identity token")

	if resp.TTL < f.ttl {
		f.logger.Debug("generated plugin identity token has shorter TTL than requested",
			"requested", f.ttl, "actual", resp.TTL)
	}

	return []byte(resp.Token.Token()), nil
}
