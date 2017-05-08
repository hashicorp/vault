package awsutil

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

type CredentialsConfig struct {
	// The access key if static credentials are being used
	AccessKey string

	// The secret key if static credentials are being used
	SecretKey string

	// The session token if it is being used
	SessionToken string

	// If specified, the region will be provided to the config of the
	// EC2RoleProvider's client. This may be useful if you want to e.g. reuse
	// the client elsewhere.
	Region string

	// The filename for the shared credentials provider, if being used
	Filename string

	// The profile for the shared credentials provider, if being used
	Profile string

	// The http.Client to use, or nil for the client to use its default
	HTTPClient *http.Client
}

func (c *CredentialsConfig) GenerateCredentialChain() (*credentials.Credentials, error) {
	var providers []credentials.Provider

	switch {
	case c.AccessKey != "" && c.SecretKey != "":
		// Add the static credential provider
		providers = append(providers, &credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     c.AccessKey,
				SecretAccessKey: c.SecretKey,
				SessionToken:    c.SessionToken,
			}})
	case c.AccessKey == "" && c.SecretKey == "":
		// Attempt to get credentials from the IAM instance role below

	default: // Have one or the other but not both and not neither
		return nil, fmt.Errorf(
			"static AWS client credentials haven't been properly configured (the access key or secret key were provided but not both)")
	}

	// Add the environment credential provider
	providers = append(providers, &credentials.EnvProvider{})

	// Add the shared credentials provider
	providers = append(providers, &credentials.SharedCredentialsProvider{
		Filename: c.Filename,
		Profile:  c.Profile,
	})

	// Add the instance metadata role provider
	providers = append(providers, &ec2rolecreds.EC2RoleProvider{
		Client: ec2metadata.New(session.New(&aws.Config{
			Region:     aws.String(c.Region),
			HTTPClient: c.HTTPClient,
		})),
		ExpiryWindow: 15,
	})

	// Create the credentials required to access the API.
	creds := credentials.NewChainCredentials(providers)
	if creds == nil {
		return nil, fmt.Errorf("could not compile valid credential providers from static config, environemnt, shared, or instance metadata")
	}

	return creds, nil
}
