package awsutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-cleanhttp"
)

type AWSCredentialsConfig struct {
	AccessKey string
	SecretKey string
	Region    string
}

func GenerateCredentialChain(config *AWSCredentialsConfig) (*credentials.Credentials, error) {
	if config == nil {
		return nil, fmt.Errorf("nil configuration provided")
	}

	var providers []credentials.Provider

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
			"static AWS client credentials haven't been properly configured (the access key or secret key were provided but not both)")
	}

	// Add the environment credential provider
	providers = append(providers, &credentials.EnvProvider{})

	// Add the instance metadata role provider
	// Create the credentials required to access the API.
	providers = append(providers, &ec2rolecreds.EC2RoleProvider{
		Client: ec2metadata.New(session.New(&aws.Config{
			Region:     aws.String(config.Region),
			HTTPClient: cleanhttp.DefaultClient(),
		})),
		ExpiryWindow: 15,
	})

	creds := credentials.NewChainCredentials(providers)
	if creds == nil {
		return nil, fmt.Errorf("could not compile valid credential providers from static config, environemnt, or instance metadata")
	}

	return creds, nil
}
