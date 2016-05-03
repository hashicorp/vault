package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/go-cleanhttp"
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
	// Read the configured secret key and access key
	config, err := b.clientConfigEntry(s)
	if err != nil {
		return nil, err
	}

	var providers []credentials.Provider

	endpoint := aws.String("")
	if config != nil {
		// Override the default endpoint with the configured endpoint.
		if config.Endpoint != "" {
			endpoint = aws.String(config.Endpoint)
		}

		switch {
		case config.AccessKey != "" && config.SecretKey != "":
			// Add the static credential provider
			providers = append(providers, &credentials.StaticProvider{
				Value: credentials.Value{
					AccessKeyID:     config.AccessKey,
					SecretAccessKey: config.SecretKey,
				}})
		case config.AccessKey == "" && config.AccessKey == "":
			// Attempt to get credentials from the IAM instance role below
		default: // Have one or the other but not both and not neither
			return nil, fmt.Errorf(
				"static AWS client credentials haven't been properly configured (the access key or secret key were provided but not both); configure or remove them at the 'config/client' endpoint")
		}
	}

	// Add the environment credential provider
	providers = append(providers, &credentials.EnvProvider{})

	// Add the instance metadata role provider
	// Create the credentials required to access the API.
	providers = append(providers, &ec2rolecreds.EC2RoleProvider{
		Client: ec2metadata.New(session.New(&aws.Config{
			Region:     aws.String(region),
			HTTPClient: cleanhttp.DefaultClient(),
		})),
		ExpiryWindow: 15,
	})

	creds := credentials.NewChainCredentials(providers)
	if creds == nil {
		return nil, fmt.Errorf("could not compile valid credential providers from static config, environemnt, or instance metadata")
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
// the cached EC2 client objects will be flushed.
//
// Write lock should be acquired using b.configMutex.Lock() before calling this method
// and lock should be released using b.configMutex.Unlock() after the method returns.
func (b *backend) flushCachedEC2Clients() {
	// deleting items in map during iteration is safe.
	for region, _ := range b.EC2ClientsMap {
		delete(b.EC2ClientsMap, region)
	}
}

// clientEC2 creates a client to interact with AWS EC2 API.
func (b *backend) clientEC2(s logical.Storage, region string) (*ec2.EC2, error) {
	b.configMutex.RLock()
	if b.EC2ClientsMap[region] != nil {
		defer b.configMutex.RUnlock()
		// If the client object was already created, return it.
		return b.EC2ClientsMap[region], nil
	}

	// Release the read lock and acquire the write lock.
	b.configMutex.RUnlock()
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	// If the client gets created while switching the locks, return it.
	if b.EC2ClientsMap[region] != nil {
		return b.EC2ClientsMap[region], nil
	}

	// Create a AWS config object using a chain of providers.
	awsConfig, err := b.getClientConfig(s, region)
	if err != nil {
		return nil, err
	}

	// Create a new EC2 client object, cache it and return the same.
	b.EC2ClientsMap[region] = ec2.New(session.New(awsConfig))
	return b.EC2ClientsMap[region], nil
}
