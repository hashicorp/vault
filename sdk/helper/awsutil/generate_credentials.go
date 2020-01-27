package awsutil

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/defaults"
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

	// The MFA token used for assuming role operation
	MFAToken string
}

func (c *CredentialsConfig) GenerateCredentialChain() (*credentials.Credentials, error) {
	if os.Getenv("AWS_SDK_LOAD_CONFIG") != "" {
		var assumeRoleTokenProvider func() (string, error)
		if c.MFAToken == "" {
			assumeRoleTokenProvider = stscreds.StdinTokenProvider
		} else {
			assumeRoleTokenProvider = func() (string, error) {
				return c.MFAToken, nil
			}
		}
		sess, err := session.NewSessionWithOptions(session.Options{
			AssumeRoleTokenProvider: assumeRoleTokenProvider,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create AWS session: %v", err)
		}
		return sess.Config.Credentials, nil
	}

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

	// Add the remote provider
	def := defaults.Get()
	if c.Region != "" {
		def.Config.Region = aws.String(c.Region)
	}
	if c.HTTPClient != nil {
		def.Config.HTTPClient = c.HTTPClient
	}

	providers = append(providers, defaults.RemoteCredProvider(*def.Config, def.Handlers))

	// Create the credentials required to access the API.
	creds := credentials.NewChainCredentials(providers)
	if creds == nil {
		return nil, fmt.Errorf("could not compile valid credential providers from static config, environment, shared, or instance metadata")
	}

	return creds, nil
}
