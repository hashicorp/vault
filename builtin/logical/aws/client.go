// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
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

// NOTE: The caller is required to ensure that b.clientMutex is at least read locked
func (b *backend) getRootConfig(ctx context.Context, s logical.Storage, clientType string, logger hclog.Logger) (*aws.Config, error) {
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
		switch {
		case clientType == "iam" && config.IAMEndpoint != "":
			endpoint = *aws.String(config.IAMEndpoint)
		case clientType == "sts" && config.STSEndpoint != "":
			endpoint = *aws.String(config.STSEndpoint)
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
		credsConfig.Region = os.Getenv("AWS_REGION")
		if credsConfig.Region == "" {
			credsConfig.Region = os.Getenv("AWS_DEFAULT_REGION")
			if credsConfig.Region == "" {
				credsConfig.Region = "us-east-1"
			}
		}
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

func (b *backend) nonCachedClientIAM(ctx context.Context, s logical.Storage, logger hclog.Logger) (*iam.IAM, error) {
	awsConfig, err := b.getRootConfig(ctx, s, "iam", logger)
	if err != nil {
		return nil, err
	}
	sess, err := session.NewSession(awsConfig)
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
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	client := sts.New(sess)
	if client == nil {
		return nil, fmt.Errorf("could not obtain sts client")
	}
	return client, nil
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
