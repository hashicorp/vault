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
	config, err := clientConfigEntry(s)
	if err != nil {
		return nil, err
	}

	var providers []credentials.Provider

	if config != nil {
		switch {
		case config.AccessKey != "" && config.SecretKey != "":
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

	providers = append(providers, &credentials.EnvProvider{})

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
	}, nil
}

// flushCachedEC2Clients deletes all the cached ec2 client objects from the backend.
func (b *backend) flushCachedEC2Clients() {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	for region, _ := range b.EC2ClientsMap {
		delete(b.EC2ClientsMap, region)
	}
}

// clientEC2 creates a client to interact with AWS EC2 API.
func (b *backend) clientEC2(s logical.Storage, region string, recreate bool) (*ec2.EC2, error) {
	if !recreate {
		b.configMutex.RLock()
		if b.EC2ClientsMap[region] != nil {
			defer b.configMutex.RUnlock()
			return b.EC2ClientsMap[region], nil
		}
		b.configMutex.RUnlock()
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	awsConfig, err := b.getClientConfig(s, region)
	if err != nil {
		return nil, err
	}

	b.EC2ClientsMap[region] = ec2.New(session.New(awsConfig))
	return b.EC2ClientsMap[region], nil
}
