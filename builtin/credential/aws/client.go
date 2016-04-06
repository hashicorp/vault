package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/logical"
)

// getClientConfig creates a aws-sdk-go config, which is used to create
// client that can interact with AWS API. This reads out the secret key
// and access key that was configured via 'config/client' endpoint and
// uses them to create credentials required to make the AWS API calls.
func getClientConfig(s logical.Storage) (*aws.Config, error) {
	// Read the configured secret key and access key
	config, err := clientConfigEntry(s)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, fmt.Errorf(
			"client credentials haven't been configured. Please configure\n" +
				"them at the 'config/client' endpoint")
	}

	// Create the credentials required to access the API.
	creds := credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, "")

	// Create a config that can be used to make the API calls.
	return &aws.Config{
		Credentials: creds,
		Region:      aws.String(config.Region),
		HTTPClient:  cleanhttp.DefaultClient(),
	}, nil
}

// clientEC2 creates a client to interact with AWS EC2 API.
func clientEC2(s logical.Storage) (*ec2.EC2, error) {
	awsConfig, err := getClientConfig(s)
	if err != nil {
		return nil, err
	}
	return ec2.New(session.New(awsConfig)), nil
}
