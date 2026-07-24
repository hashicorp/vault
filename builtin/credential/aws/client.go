// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	awsutil "github.com/hashicorp/go-secure-stdlib/awsutil/v2"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const useServiceDefaultRetries = -1

// getRawClientConfig creates an aws-sdk-go-v2 config used to create AWS clients.
// This builds credentials in the following order of preference:
//
// * Static credentials from 'config/client'
// * Environment variables
// * Instance metadata role
func (b *backend) getRawClientConfig(ctx context.Context, s logical.Storage, region, clientType string) (*awsv2.Config, error) {
	credsConfig := &awsutil.CredentialsConfig{
		Region: region,
		Logger: b.Logger(),
	}

	// Read the configured secret key and access key
	config, err := b.nonLockedClientConfigEntry(ctx, s)
	if err != nil {
		return nil, err
	}

	var endpoint *string
	if config != nil {
		// Override the defaults with configured values.
		switch {
		case clientType == "ec2" && config.Endpoint != "":
			endpoint = awsv2.String(config.Endpoint)
		case clientType == "iam" && config.IAMEndpoint != "":
			endpoint = awsv2.String(config.IAMEndpoint)
		case clientType == "sts":
			if config.STSEndpoint != "" {
				endpoint = awsv2.String(config.STSEndpoint)
			}
			if config.STSRegion != "" {
				region = config.STSRegion
				credsConfig.Region = region // v2 reads region from credsConfig.Region, so copy the sts_region override here.
			}
		}

		credsConfig.AccessKey = config.AccessKey
		credsConfig.SecretKey = config.SecretKey
		if config.MaxRetries >= 0 {
			credsConfig.MaxRetries = &config.MaxRetries
		}
	}

	credsConfig.HTTPClient = cleanhttp.DefaultClient()

	// When no static credentials are configured, don't force the shared "default"
	// profile. In SDK v2, setting the shared profile makes credential resolution
	// short-circuit to that profile and skip the environment and IMDS/ECS
	// providers; leaving it unset lets the default chain (env vars, shared config,
	// IMDS/ECS) resolve naturally, matching the v1 SDK behavior.
	var credChainOpts []awsutil.Option
	if credsConfig.AccessKey == "" && credsConfig.SecretKey == "" {
		credChainOpts = append(credChainOpts, awsutil.WithSharedCredentials(false))
	}

	awsConfig, err := credsConfig.GenerateCredentialChain(ctx, credChainOpts...)
	if err != nil {
		return nil, err
	}
	awsConfig.HTTPClient = credsConfig.HTTPClient
	awsConfig.BaseEndpoint = endpoint

	// SDK v2 removed WebIdentityTokenFetcher from CredentialsConfig. WIF logic now executes
	// after GenerateCredentialChain to use the resolved *aws.Config for building the STS client.
	if config != nil && config.IdentityTokenAudience != "" { // nil check guards against a nil-pointer dereference when reading IdentityTokenAudience.
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
		credsConfig.RoleSessionName = fmt.Sprintf("vault-aws-auth-%s", sessionSuffix)
		credsConfig.RoleARN = config.RoleARN
		stsClient := sts.NewFromConfig(*awsConfig)
		provider := stscreds.NewWebIdentityRoleProvider(stsClient, credsConfig.RoleARN, fetcher, func(o *stscreds.WebIdentityRoleOptions) {
			o.RoleSessionName = credsConfig.RoleSessionName
		})
		awsConfig.Credentials = awsv2.NewCredentialsCache(provider)
	}

	if awsConfig.Credentials == nil {
		return nil, fmt.Errorf("could not compile valid credential providers from static config, environment, shared, or instance metadata")
	}

	return awsConfig, nil
}

// getClientConfig returns an aws-sdk-go-v2 config, with optionally assumed credentials.
// It uses getRawClientConfig to obtain config for the runtime environment, and if
// stsRole is a non-empty string, it will use AssumeRole to obtain a set of assumed
// credentials. The credentials will expire after 15 minutes but will auto-refresh.
func (b *backend) getClientConfig(ctx context.Context, s logical.Storage, region, stsRole, externalID, accountID, clientType string) (*awsv2.Config, error) {
	config, err := b.getRawClientConfig(ctx, s, region, clientType)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, fmt.Errorf("could not compile valid credentials through the default provider chain")
	}

	stsConfig, err := b.getRawClientConfig(ctx, s, region, "sts")
	if err != nil {
		return nil, err
	}
	if stsConfig == nil {
		return nil, fmt.Errorf("could not configure STS client")
	}
	if stsRole != "" {
		stsClient := sts.NewFromConfig(*stsConfig)
		provider := stscreds.NewAssumeRoleProvider(stsClient, stsRole, func(o *stscreds.AssumeRoleOptions) {
			if externalID != "" {
				o.ExternalID = awsv2.String(externalID)
			}
		})
		assumedCredentials := awsv2.NewCredentialsCache(provider)
		// Test that we actually have permissions to assume the role
		if _, err = assumedCredentials.Retrieve(ctx); err != nil {
			return nil, err
		}
		config.Credentials = assumedCredentials
	} else {
		if b.defaultAWSAccountID == "" {
			client := sts.NewFromConfig(*stsConfig)
			identity, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
			if err != nil {
				return nil, fmt.Errorf("unable to fetch current caller: %w", err)
			}
			if identity == nil {
				return nil, fmt.Errorf("got nil result from GetCallerIdentity")
			}
			if identity.Account == nil {
				return nil, fmt.Errorf("got nil account from GetCallerIdentity")
			}
			b.defaultAWSAccountID = *identity.Account
		}
		if b.defaultAWSAccountID != accountID {
			return nil, fmt.Errorf("unable to fetch client for account ID %q -- default client is for account %q", accountID, b.defaultAWSAccountID)
		}
	}

	return config, nil
}

// flushCachedEC2Clients deletes all the cached ec2 client objects from the backend.
// If the client credentials configuration is deleted or updated in the backend, all
// the cached EC2 client objects will be flushed. Config mutex lock should be
// acquired for write operation before calling this method.
func (b *backend) flushCachedEC2Clients() {
	b.EC2ClientsMap = make(map[clientKey]*ec2.Client)
}

// flushCachedIAMClients deletes all the cached iam client objects from the
// backend. If the client credentials configuration is deleted or updated in
// the backend, all the cached IAM client objects will be flushed. Config mutex
// lock should be acquired for write operation before calling this method.
func (b *backend) flushCachedIAMClients() {
	b.IAMClientsMap = make(map[clientKey]*iam.Client)
}

// Gets an entry out of the user ID cache
func (b *backend) getCachedUserId(userId string) string {
	if userId == "" {
		return ""
	}
	if entry, ok := b.iamUserIdToArnCache.Get(userId); ok {
		b.iamUserIdToArnCache.SetDefault(userId, entry)
		return entry.(string)
	}
	return ""
}

// Sets an entry in the user ID cache
func (b *backend) setCachedUserId(userId, arn string) {
	if userId != "" {
		b.iamUserIdToArnCache.SetDefault(userId, arn)
	}
}

func (b *backend) stsRoleForAccount(ctx context.Context, s logical.Storage, accountID string) (string, string, error) {
	// Check if an STS configuration exists for the AWS account
	sts, err := b.lockedAwsStsEntry(ctx, s, accountID)
	if err != nil {
		return "", "", fmt.Errorf("error fetching STS config for account ID %q: %w", accountID, err)
	}
	// An empty STS role signifies the master account
	if sts != nil {
		return sts.StsRole, sts.ExternalID, nil
	}

	// Return an error if there's no STS config for an account which is not the default one
	if b.defaultAWSAccountID != "" && b.defaultAWSAccountID != accountID {
		return "", "", fmt.Errorf("no STS configuration found for account ID %q", accountID)
	}
	return "", "", nil
}

// clientEC2 creates a client to interact with AWS EC2 API
func (b *backend) clientEC2(ctx context.Context, s logical.Storage, region, accountID string) (*ec2.Client, error) {
	stsRole, stsExternalID, err := b.stsRoleForAccount(ctx, s, accountID)
	if err != nil {
		return nil, err
	}
	b.configMutex.RLock()

	key := clientKey{
		AccountID: accountID,
		Region:    region,
		STSRole:   stsRole,
	}
	if cachedClient, ok := b.EC2ClientsMap[key]; ok {
		defer b.configMutex.RUnlock()
		// If the client object was already created, return it
		return cachedClient, nil
	}

	// Release the read lock and acquire the write lock
	b.configMutex.RUnlock()
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	// If the client gets created while switching the locks, return it
	if cachedClient, ok := b.EC2ClientsMap[key]; ok {
		return cachedClient, nil
	}

	// Create an AWS config object using a chain of providers
	var awsConfig *awsv2.Config
	awsConfig, err = b.getClientConfig(ctx, s, region, stsRole, stsExternalID, accountID, "ec2")
	if err != nil {
		return nil, err
	}

	if awsConfig == nil {
		return nil, fmt.Errorf("could not retrieve valid assumed credentials")
	}

	// Create a new EC2 client object, cache it and return the same
	client := ec2.NewFromConfig(*awsConfig)
	if client == nil {
		return nil, fmt.Errorf("could not obtain ec2 client")
	}

	b.EC2ClientsMap[key] = client
	return b.EC2ClientsMap[key], nil
}

// clientIAM creates a client to interact with AWS IAM API
func (b *backend) clientIAM(ctx context.Context, s logical.Storage, region, accountID string) (*iam.Client, error) {
	stsRole, stsExternalID, err := b.stsRoleForAccount(ctx, s, accountID)
	if err != nil {
		return nil, err
	}
	if stsRole == "" {
		b.Logger().Debug("no stsRole found for account", "accountID", accountID)
	} else {
		b.Logger().Debug("found stsRole for account", "stsRole", stsRole, "accountID", accountID)
	}
	b.configMutex.RLock()

	key := clientKey{
		AccountID: accountID,
		Region:    region,
		STSRole:   stsRole,
	}
	if cachedClient, ok := b.IAMClientsMap[key]; ok {
		defer b.configMutex.RUnlock()
		// If the client object was already created, return it
		b.Logger().Debug("returning cached client for key", "key", key)
		return cachedClient, nil
	}
	b.Logger().Debug("no cached client for key", "key", key)

	// Release the read lock and acquire the write lock
	b.configMutex.RUnlock()
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	// If the client gets created while switching the locks, return it
	if cachedClient, ok := b.IAMClientsMap[key]; ok {
		b.Logger().Debug("returning cached client for key", "key", key)
		return cachedClient, nil
	}

	// Create an AWS config object using a chain of providers
	var awsConfig *awsv2.Config
	awsConfig, err = b.getClientConfig(ctx, s, region, stsRole, stsExternalID, accountID, "iam")
	if err != nil {
		return nil, err
	}

	if awsConfig == nil {
		return nil, fmt.Errorf("could not retrieve valid assumed credentials")
	}

	// Create a new IAM client object, cache it and return the same
	client := iam.NewFromConfig(*awsConfig)
	if client == nil {
		return nil, fmt.Errorf("could not obtain iam client")
	}
	b.IAMClientsMap[key] = client
	return b.IAMClientsMap[key], nil
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

var _ stscreds.IdentityTokenRetriever = (*PluginIdentityTokenFetcher)(nil)

func (f PluginIdentityTokenFetcher) GetIdentityToken() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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
