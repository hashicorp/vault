package awsec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/awsutil"
	"github.com/hashicorp/vault/logical"
)

// getClientConfig creates a aws-sdk-go config, which is used to create client
// that can interact with AWS API. This builds credentials in the following
// order of preference:
//
// * Static credentials from 'config/client'
// * Environment variables
// * Instance metadata role
func (b *backend) getClientConfig(s logical.Storage, region string) (*aws.Config, error) {
	credsConfig := &awsutil.CredentialsConfig{
		Region: region,
	}

	// Read the configured secret key and access key
	config, err := b.nonLockedClientConfigEntry(s)
	if err != nil {
		return nil, err
	}

	endpoint := aws.String("")
	if config != nil {
		// Override the default endpoint with the configured endpoint.
		if config.Endpoint != "" {
			endpoint = aws.String(config.Endpoint)
		}

		credsConfig.AccessKey = config.AccessKey
		credsConfig.SecretKey = config.SecretKey
	}

	credsConfig.HTTPClient = cleanhttp.DefaultClient()

	creds, err := credsConfig.GenerateCredentialChain()
	if err != nil {
		return nil, err
	}
	if creds == nil {
		return nil, fmt.Errorf("could not compile valid credential providers from static config, environemnt, shared, or instance metadata")
	}

	// Create a config that can be used to make the API calls.
	return &aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
		HTTPClient:  cleanhttp.DefaultClient(),
		Endpoint:    endpoint,
	}, nil
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

// clientEC2 creates a client to interact with AWS EC2 API
func (b *backend) clientEC2(s logical.Storage, region string) (*ec2.EC2, error) {
	b.configMutex.RLock()
	if b.EC2ClientsMap[region] != nil {
		defer b.configMutex.RUnlock()
		// If the client object was already created, return it
		return b.EC2ClientsMap[region], nil
	}

	// Release the read lock and acquire the write lock
	b.configMutex.RUnlock()
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	// If the client gets created while switching the locks, return it
	if b.EC2ClientsMap[region] != nil {
		return b.EC2ClientsMap[region], nil
	}

	// Create an AWS config object using a chain of providers
	awsConfig, err := b.getClientConfig(s, region)
	if err != nil {
		return nil, err
	}

	// Create a new EC2 client object, cache it and return the same
	b.EC2ClientsMap[region] = ec2.New(session.New(awsConfig))
	return b.EC2ClientsMap[region], nil
}

// clientIAM creates a client to interact with AWS IAM API
func (b *backend) clientIAM(s logical.Storage, region string) (*iam.IAM, error) {
	b.configMutex.RLock()
	if b.IAMClientsMap[region] != nil {
		defer b.configMutex.RUnlock()
		// If the client object was already created, return it
		return b.IAMClientsMap[region], nil
	}

	// Release the read lock and acquire the write lock
	b.configMutex.RUnlock()
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	// If the client gets created while switching the locks, return it
	if b.IAMClientsMap[region] != nil {
		return b.IAMClientsMap[region], nil
	}

	// Create an AWS config object using a chain of providers
	awsConfig, err := b.getClientConfig(s, region)
	if err != nil {
		return nil, err
	}

	// Create a new IAM client object, cache it and return the same
	b.IAMClientsMap[region] = iam.New(session.New(awsConfig))
	return b.IAMClientsMap[region], nil
}
