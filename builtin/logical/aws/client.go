// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	awsutil "github.com/hashicorp/go-secure-stdlib/awsutil/v2"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// identityTokenFetchTimeout bounds how long GetIdentityToken waits on the plugin
// identity token service, since the stscreds.IdentityTokenRetriever interface
// provides no context to carry a caller deadline.
const identityTokenFetchTimeout = 30 * time.Second

// getRootIAMConfig creates an *aws.Config for Vault to connect to IAM.
func (b *backend) getRootIAMConfig(ctx context.Context, s logical.Storage, logger hclog.Logger) (*aws.Config, error) {
	credsConfig := &awsutil.CredentialsConfig{
		HTTPClient: cleanhttp.DefaultClient(),
		Logger:     logger,
	}
	var endpoint string

	entry, err := s.Get(ctx, "config/root")
	if err != nil {
		return nil, err
	}

	var config rootConfig
	if entry != nil {
		if err := entry.DecodeJSON(&config); err != nil {
			return nil, fmt.Errorf("error reading root configuration: %w", err)
		}

		credsConfig.AccessKey = config.AccessKey
		credsConfig.SecretKey = config.SecretKey
		credsConfig.Region = config.Region
		if config.MaxRetries >= 0 {
			credsConfig.MaxRetries = aws.Int(config.MaxRetries)
		}

		if config.IAMEndpoint != "" {
			endpoint = config.IAMEndpoint
		}
	}

	// Apply the region fallback once for all paths (WIF and non-WIF).
	if credsConfig.Region == "" {
		credsConfig.Region = getFallbackRegion()
	}

	if entry != nil && config.IdentityTokenAudience != "" {
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

		// Build a base config so we can construct the STS client used
		// by WebIdentityRoleProvider.
		baseCfg, err := credsConfig.GenerateCredentialChain(ctx)
		if err != nil {
			return nil, err
		}
		if endpoint != "" {
			baseCfg.BaseEndpoint = aws.String(endpoint)
		}

		// Wire the fetcher as a live token retriever so credentials refresh automatically.
		attachWebIdentityProvider(baseCfg, config.RoleARN, sessionSuffix, fetcher)
		return baseCfg, nil
	}

	// When no static credentials are configured, disable forcing the shared
	// "default" profile. Forcing it makes the v2 SDK short-circuit to the shared
	// profile and skip the environment credential provider (and IMDS/ECS);
	// disabling it lets the default chain (env vars, shared config, IMDS/ECS)
	// resolve naturally, matching the v1 SDK behavior.
	iamOpts := make([]awsutil.Option, 0)
	if credsConfig.AccessKey == "" && credsConfig.SecretKey == "" {
		iamOpts = append(iamOpts, awsutil.WithSharedCredentials(false))
	}

	awsConfig, err := credsConfig.GenerateCredentialChain(ctx, iamOpts...)
	if err != nil {
		return nil, err
	}

	if endpoint != "" {
		awsConfig.BaseEndpoint = aws.String(endpoint)
	}

	return awsConfig, nil
}

// Return a slice of *aws.Config, based on descending configuration priority. STS endpoints are the only place this is used.
// NOTE: The caller is required to ensure that b.clientMutex is at least read locked
func (b *backend) getRootSTSConfigs(ctx context.Context, s logical.Storage, logger hclog.Logger) ([]*aws.Config, error) {
	// set fallback region (we can overwrite later)
	fallbackRegion := getFallbackRegion()

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
		// No config means no static credentials; disable forcing the shared
		// "default" profile so the v2 SDK does not short-circuit to it and skip the
		// env var / IMDS/ECS providers. This lets the default chain resolve
		// naturally, matching the v1 SDK behavior.
		awsConfig, err := ccfg.GenerateCredentialChain(ctx, awsutil.WithSharedCredentials(false))
		if err != nil {
			return nil, err
		}
		configs = append(configs, awsConfig)

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
	if config.MaxRetries >= 0 {
		credsConfig.MaxRetries = aws.Int(config.MaxRetries)
	}

	if config.Region != "" {
		regions = append(regions, config.Region)
	}

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

	// wifFetcher and wifSessionSuffix are set when WIF is configured and used
	// inside the per-config loop to attach a live-refresh WebIdentityRoleProvider.
	var wifFetcher *PluginIdentityTokenFetcher
	var wifSessionSuffix string

	opts := make([]awsutil.Option, 0)
	if config.IdentityTokenAudience != "" {
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get namespace from context: %w", err)
		}

		wifFetcher = &PluginIdentityTokenFetcher{
			sys:      b.System(),
			logger:   b.Logger(),
			ns:       ns,
			audience: config.IdentityTokenAudience,
			ttl:      config.IdentityTokenTTL,
		}
		wifSessionSuffix = strconv.FormatInt(time.Now().UnixNano(), 10)

		// explicitly disable shared credential providers when using a web identity
		// token, enabling WIF usage in environments that may use AWS Profiles for
		// other use-cases
		opts = append(opts, awsutil.WithSharedCredentials(false))
	}

	// When no static credentials are configured (and not using WIF), disable
	// forcing the shared "default" profile. Forcing it makes the v2 SDK
	// short-circuit to the shared profile and skip the environment credential
	// provider (and IMDS/ECS); disabling it lets the default chain (env vars,
	// shared config, IMDS/ECS) resolve naturally, matching the v1 SDK.
	if config.IdentityTokenAudience == "" && config.AccessKey == "" && config.SecretKey == "" {
		opts = append(opts, awsutil.WithSharedCredentials(false))
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
		awsConfig, err := credsConfig.GenerateCredentialChain(ctx, opts...)
		if err != nil {
			return nil, err
		}
		if endpoints[i] != "" {
			awsConfig.BaseEndpoint = aws.String(endpoints[i])
		}

		// Wire the fetcher as a live token retriever so credentials refresh automatically.
		if wifFetcher != nil {
			attachWebIdentityProvider(awsConfig, config.RoleARN, wifSessionSuffix, wifFetcher)
		}

		configs = append(configs, awsConfig)
	}

	return configs, nil
}

// attachWebIdentityProvider wires a live-refresh WebIdentityRoleProvider onto cfg so
// that STS credentials are automatically renewed using the plugin identity token
// fetcher. It is shared by the IAM and STS root config paths to avoid duplication.
func attachWebIdentityProvider(cfg *aws.Config, roleARN, sessionSuffix string, fetcher *PluginIdentityTokenFetcher) {
	stsClient := sts.NewFromConfig(*cfg)
	provider := stscreds.NewWebIdentityRoleProvider(
		stsClient,
		roleARN,
		fetcher,
		func(o *stscreds.WebIdentityRoleOptions) {
			o.RoleSessionName = fmt.Sprintf("vault-aws-secrets-%s", sessionSuffix)
		},
	)
	cfg.Credentials = aws.NewCredentialsCache(provider)
}

func (b *backend) nonCachedClientIAM(ctx context.Context, s logical.Storage, logger hclog.Logger, entry *staticRoleEntry) (iamAPI, error) {
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

	// iam.NewFromConfig always returns a non-nil client, so there is no error
	// case to handle here.
	return iam.NewFromConfig(*awsConfig), nil
}

func (b *backend) nonCachedClientSTS(ctx context.Context, s logical.Storage, logger hclog.Logger) (stsAPI, error) {
	awsConfigs, err := b.getRootSTSConfigs(ctx, s, logger)
	if err != nil {
		return nil, err
	}

	for _, cfg := range awsConfigs {
		client := sts.NewFromConfig(*cfg)

		// ping the client - we only care about errors
		_, err = client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err == nil {
			return client, nil
		}

		endpoint := ""
		if cfg.BaseEndpoint != nil {
			endpoint = *cfg.BaseEndpoint
		}
		b.Logger().Debug("couldn't connect with config trying next", "failed endpoint", endpoint, "failed region", cfg.Region)
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

func (f PluginIdentityTokenFetcher) FetchToken(ctx context.Context) ([]byte, error) {
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

// GetIdentityToken implements stscreds.IdentityTokenRetriever for use with stscreds.NewWebIdentityRoleProvider.
// The interface provides no context, so a bounded timeout is applied to avoid
// blocking indefinitely if the plugin identity token service is unresponsive.
func (f PluginIdentityTokenFetcher) GetIdentityToken() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), identityTokenFetchTimeout)
	defer cancel()
	return f.FetchToken(ctx)
}
