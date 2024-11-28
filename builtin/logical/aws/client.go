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

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-secure-stdlib/awsutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// Return a slice of *aws.Config, based on descending configuration priority. STS endpoints are the only place this is used.
// NOTE: The caller is required to ensure that b.clientMutex is at least read locked
func (b *backend) getRootConfig(ctx context.Context, s logical.Storage, clientType string, logger hclog.Logger) ([]*aws.Config, error) {
	// to make an aws config we need:
	// credsConfig (to call generateCredentialsChain)
	//   - accessKey (string)
	//   - secretKey (string)
	//   - Region (string)
	//   - RoleSessionName (string) (if identtyTokenAudience not empty)
	//   - WebIdentityTokenFactor
	//

	// set fallback region (we can overwrite later)
	fallbackRegion := os.Getenv("AWS_REGION")
	if fallbackRegion == "" {
		fallbackRegion = os.Getenv("AWS_DEFAULT_REGION")
	}
	if fallbackRegion == "" {
		fallbackRegion = "us-east-1"
	}

	maxRetries := aws.UseServiceDefaultRetries

	entry, err := s.Get(ctx, "config/root")
	if err != nil {
		return nil, err
	}
	var configs []*aws.Config

	// I'm not sure this is a valid scenario, but the previous code had it as a case.
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

	maxRetries = config.MaxRetries
	if clientType == "iam" && config.IAMEndpoint != "" {
		endpoints = append(endpoints, config.IAMEndpoint)
	} else if clientType == "sts" && len(config.STSEndpoint) > 0 {
		endpoints = append(endpoints, config.STSEndpoint...)
		regions = append(regions, config.STSRegion...)
	}
	// len(endpoints) is zero here if no specific endpoints were set
	// len(regions) is zero if no specific sts regions were set

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

	if len(regions) == 0 {
		regions = append(regions, fallbackRegion)
	}

	if len(regions) != len(endpoints) {
		// this probably can't happen, if the input was checked correctly
		return nil, errors.New("number of regions does not match number of endpoints")
	}

	for i := 0; i < len(regions); i++ {
		credsConfig.Region = regions[i]
		creds, err := credsConfig.GenerateCredentialChain()
		if err != nil {
			return nil, err
		}
		configs = append(configs, &aws.Config{
			Credentials: creds,
			Region:      aws.String(regions[i]),
			Endpoint:    aws.String(endpoints[i]),
			MaxRetries:  aws.Int(maxRetries),
			HTTPClient:  cleanhttp.DefaultClient(),
		})
	}

	return configs, nil

	//var endpoint string
	//maxRetries := aws.UseServiceDefaultRetries
	//
	//entry, err := s.Get(ctx, "config/root")
	//if err != nil {
	//	return nil, err
	//}
	//if entry != nil {
	//	var config rootConfig
	//	if err := entry.DecodeJSON(&config); err != nil {
	//		return nil, fmt.Errorf("error reading root configuration: %w", err)
	//	}
	//
	//	credsConfig.AccessKey = config.AccessKey
	//	credsConfig.SecretKey = config.SecretKey
	//	credsConfig.Region = config.Region
	//	maxRetries = config.MaxRetries
	//	fmt.Println(clientType)
	//	switch {
	//	case clientType == "iam" && config.IAMEndpoint != "":
	//		endpoint = *aws.String(config.IAMEndpoint)
	//	case clientType == "sts" && len(config.STSEndpoint) > 0:
	//		endpoint = *aws.String(config.STSEndpoint[0])
	//		if len(config.STSRegion) > 0 {
	//			credsConfig.Region = config.STSRegion[0]
	//		}
	//	}
	//
	//	if config.IdentityTokenAudience != "" {
	//		ns, err := namespace.FromContext(ctx)
	//		if err != nil {
	//			return nil, fmt.Errorf("failed to get namespace from context: %w", err)
	//		}
	//
	//		fetcher := &PluginIdentityTokenFetcher{
	//			sys:      b.System(),
	//			logger:   b.Logger(),
	//			ns:       ns,
	//			audience: config.IdentityTokenAudience,
	//			ttl:      config.IdentityTokenTTL,
	//		}
	//
	//		sessionSuffix := strconv.FormatInt(time.Now().UnixNano(), 10)
	//		credsConfig.RoleSessionName = fmt.Sprintf("vault-aws-secrets-%s", sessionSuffix)
	//		credsConfig.WebIdentityTokenFetcher = fetcher
	//		credsConfig.RoleARN = config.RoleARN
	//	}
	//}
	//
	//if credsConfig.Region == "" {
	//	credsConfig.Region = os.Getenv("AWS_REGION")
	//	if credsConfig.Region == "" {
	//		credsConfig.Region = os.Getenv("AWS_DEFAULT_REGION")
	//		if credsConfig.Region == "" {
	//			credsConfig.Region = "us-east-1"
	//		}
	//	}
	//}
	//
	//credsConfig.HTTPClient = cleanhttp.DefaultClient()
	//
	//credsConfig.Logger = logger
	//
	//creds, err := credsConfig.GenerateCredentialChain()
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &aws.Config{
	//	Credentials: creds,
	//	Region:      aws.String(credsConfig.Region),
	//	Endpoint:    &endpoint,
	//	HTTPClient:  cleanhttp.DefaultClient(),
	//	MaxRetries:  aws.Int(maxRetries),
	//}, nil
}

func (b *backend) nonCachedClientIAM(ctx context.Context, s logical.Storage, logger hclog.Logger) (*iam.IAM, error) {
	awsConfig, err := b.getRootConfig(ctx, s, "iam", logger)
	if err != nil {
		return nil, err
	}
	sess, err := session.NewSession(awsConfig[0])
	if err != nil {
		return nil, err
	}
	client := iam.New(sess)
	if client == nil {
		return nil, fmt.Errorf("could not obtain iam client")
	}
	return client, nil
}

func (b *backend) nonCachedClientSTS(ctx context.Context, s logical.Storage, logger hclog.Logger) (*sts.STS, error) {
	awsConfig, err := b.getRootConfig(ctx, s, "sts", logger)
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
			fmt.Println("couldn't connect with config, trying the next one")
		}
	}

	return nil, fmt.Errorf("could not obtain sts client")
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
