package awsauth

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/awsutil"
	"github.com/hashicorp/vault/logical"
)

// getRawClientConfig creates a aws-sdk-go config, which is used to create client
// that can interact with AWS API. This builds credentials in the following
// order of preference:
//
// * Static credentials from 'config/client'
// * Environment variables
// * Instance metadata role
func (b *backend) getRawClientConfig(ctx context.Context, s logical.Storage, region, clientType string) (*aws.Config, error) {
	credsConfig := &awsutil.CredentialsConfig{
		Region: region,
	}

	// Read the configured secret key and access key
	config, err := b.nonLockedClientConfigEntry(ctx, s)
	if err != nil {
		return nil, err
	}

	endpoint := aws.String("")
	var maxRetries int = aws.UseServiceDefaultRetries
	if config != nil {
		// Override the default endpoint with the configured endpoint.
		switch {
		case clientType == "ec2" && config.Endpoint != "":
			endpoint = aws.String(config.Endpoint)
		case clientType == "iam" && config.IAMEndpoint != "":
			endpoint = aws.String(config.IAMEndpoint)
		case clientType == "sts" && config.STSEndpoint != "":
			endpoint = aws.String(config.STSEndpoint)
		}

		credsConfig.AccessKey = config.AccessKey
		credsConfig.SecretKey = config.SecretKey
		maxRetries = config.MaxRetries
	}

	credsConfig.HTTPClient = cleanhttp.DefaultClient()

	creds, err := credsConfig.GenerateCredentialChain()
	if err != nil {
		return nil, err
	}
	if creds == nil {
		return nil, fmt.Errorf("could not compile valid credential providers from static config, environment, shared, or instance metadata")
	}

	// Create a config that can be used to make the API calls.
	return &aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
		HTTPClient:  cleanhttp.DefaultClient(),
		Endpoint:    endpoint,
		MaxRetries:  aws.Int(maxRetries),
	}, nil
}

// getClientConfig returns an aws-sdk-go config, with optionally assumed credentials
// It uses getRawClientConfig to obtain config for the runtime environment, and if
// stsRole is a non-empty string, it will use AssumeRole to obtain a set of assumed
// credentials. The credentials will expire after 15 minutes but will auto-refresh.
func (b *backend) getClientConfig(ctx context.Context, s logical.Storage, region, stsRole, accountID, clientType string) (*aws.Config, error) {

	config, err := b.getRawClientConfig(ctx, s, region, clientType)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, fmt.Errorf("could not compile valid credentials through the default provider chain")
	}

	stsConfig, err := b.getRawClientConfig(ctx, s, region, "sts")
	if stsConfig == nil {
		return nil, fmt.Errorf("could not configure STS client")
	}
	if err != nil {
		return nil, err
	}
	if stsRole != "" {
		assumedCredentials := stscreds.NewCredentials(session.New(stsConfig), stsRole)
		// Test that we actually have permissions to assume the role
		if _, err = assumedCredentials.Get(); err != nil {
			return nil, err
		}
		config.Credentials = assumedCredentials
	} else {
		if b.defaultAWSAccountID == "" {
			client := sts.New(session.New(stsConfig))
			if client == nil {
				return nil, errwrap.Wrapf("could not obtain sts client: {{err}}", err)
			}
			inputParams := &sts.GetCallerIdentityInput{}
			identity, err := client.GetCallerIdentity(inputParams)
			if err != nil {
				return nil, errwrap.Wrapf("unable to fetch current caller: {{err}}", err)
			}
			if identity == nil {
				return nil, fmt.Errorf("got nil result from GetCallerIdentity")
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
	// deleting items in map during iteration is safe
	for region, _ := range b.EC2ClientsMap {
		delete(b.EC2ClientsMap, region)
	}
}

// flushCachedIAMClients deletes all the cached iam client objects from the
// backend. If the client credentials configuration is deleted or updated in
// the backend, all the cached IAM client objects will be flushed. Config mutex
// lock should be acquired for write operation before calling this method.
func (b *backend) flushCachedIAMClients() {
	// deleting items in map during iteration is safe
	for region, _ := range b.IAMClientsMap {
		delete(b.IAMClientsMap, region)
	}
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

func (b *backend) stsRoleForAccount(ctx context.Context, s logical.Storage, accountID string) (string, error) {
	// Check if an STS configuration exists for the AWS account
	sts, err := b.lockedAwsStsEntry(ctx, s, accountID)
	if err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf("error fetching STS config for account ID %q: {{err}}", accountID), err)
	}
	// An empty STS role signifies the master account
	if sts != nil {
		return sts.StsRole, nil
	}
	return "", nil
}

// clientEC2 creates a client to interact with AWS EC2 API
func (b *backend) clientEC2(ctx context.Context, s logical.Storage, region, accountID string) (*ec2.EC2, error) {
	stsRole, err := b.stsRoleForAccount(ctx, s, accountID)
	if err != nil {
		return nil, err
	}
	b.configMutex.RLock()
	if b.EC2ClientsMap[region] != nil && b.EC2ClientsMap[region][stsRole] != nil {
		defer b.configMutex.RUnlock()
		// If the client object was already created, return it
		return b.EC2ClientsMap[region][stsRole], nil
	}

	// Release the read lock and acquire the write lock
	b.configMutex.RUnlock()
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	// If the client gets created while switching the locks, return it
	if b.EC2ClientsMap[region] != nil && b.EC2ClientsMap[region][stsRole] != nil {
		return b.EC2ClientsMap[region][stsRole], nil
	}

	// Create an AWS config object using a chain of providers
	var awsConfig *aws.Config
	awsConfig, err = b.getClientConfig(ctx, s, region, stsRole, accountID, "ec2")

	if err != nil {
		return nil, err
	}

	if awsConfig == nil {
		return nil, fmt.Errorf("could not retrieve valid assumed credentials")
	}

	// Create a new EC2 client object, cache it and return the same
	client := ec2.New(session.New(awsConfig))
	if client == nil {
		return nil, fmt.Errorf("could not obtain ec2 client")
	}
	if _, ok := b.EC2ClientsMap[region]; !ok {
		b.EC2ClientsMap[region] = map[string]*ec2.EC2{stsRole: client}
	} else {
		b.EC2ClientsMap[region][stsRole] = client
	}

	return b.EC2ClientsMap[region][stsRole], nil
}

// clientIAM creates a client to interact with AWS IAM API
func (b *backend) clientIAM(ctx context.Context, s logical.Storage, region, accountID string) (*iam.IAM, error) {
	stsRole, err := b.stsRoleForAccount(ctx, s, accountID)
	if err != nil {
		return nil, err
	}
	b.configMutex.RLock()
	if b.IAMClientsMap[region] != nil && b.IAMClientsMap[region][stsRole] != nil {
		defer b.configMutex.RUnlock()
		// If the client object was already created, return it
		return b.IAMClientsMap[region][stsRole], nil
	}

	// Release the read lock and acquire the write lock
	b.configMutex.RUnlock()
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	// If the client gets created while switching the locks, return it
	if b.IAMClientsMap[region] != nil && b.IAMClientsMap[region][stsRole] != nil {
		return b.IAMClientsMap[region][stsRole], nil
	}

	// Create an AWS config object using a chain of providers
	var awsConfig *aws.Config
	awsConfig, err = b.getClientConfig(ctx, s, region, stsRole, accountID, "iam")

	if err != nil {
		return nil, err
	}

	if awsConfig == nil {
		return nil, fmt.Errorf("could not retrieve valid assumed credentials")
	}

	// Create a new IAM client object, cache it and return the same
	client := iam.New(session.New(awsConfig))
	if client == nil {
		return nil, fmt.Errorf("could not obtain iam client")
	}
	if _, ok := b.IAMClientsMap[region]; !ok {
		b.IAMClientsMap[region] = map[string]*iam.IAM{stsRole: client}
	} else {
		b.IAMClientsMap[region][stsRole] = client
	}
	return b.IAMClientsMap[region][stsRole], nil
}
