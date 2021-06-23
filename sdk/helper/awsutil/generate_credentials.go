package awsutil

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
)

const iamServerIdHeader = "X-Vault-AWS-IAM-Server-ID"

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

	// The logger to use for credential acquisition debugging
	Logger hclog.Logger
}

// Make sure the logger isn't nil before logging
func (c *CredentialsConfig) log(level hclog.Level, msg string, args ...interface{}) {
	if c.Logger != nil {
		c.Logger.Log(level, msg, args...)
	}
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
			},
		})
		c.log(hclog.Debug, "added static credential provider", "AccessKey", c.AccessKey)

	case c.AccessKey == "" && c.SecretKey == "":
		// Attempt to get credentials from the IAM instance role below

	default: // Have one or the other but not both and not neither
		return nil, fmt.Errorf(
			"static AWS client credentials haven't been properly configured (the access key or secret key were provided but not both)")
	}

	roleARN := os.Getenv("AWS_ROLE_ARN")
	tokenPath := os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE")
	sessionName := os.Getenv("AWS_ROLE_SESSION_NAME")
	if roleARN != "" && tokenPath != "" {
		// this session is only created to create the WebIdentityRoleProvider, as the env variables are already there
		// this automatically assumes the role, but the provider needs to be added to the chain
		c.log(hclog.Debug, "adding web identity provider", "roleARN", roleARN)
		sess, err := session.NewSession()
		if err != nil {
			return nil, errors.Wrap(err, "error creating a new session to create a WebIdentityRoleProvider")
		}
		webIdentityProvider := stscreds.NewWebIdentityRoleProvider(sts.New(sess), roleARN, sessionName, tokenPath)

		// Check if the webIdentityProvider can successfully retrieve
		// credentials (via sts:AssumeRole), and warn if there's a problem.
		if _, err := webIdentityProvider.Retrieve(); err != nil {
			c.log(hclog.Warn, "error assuming role", "roleARN", roleARN, "tokenPath", tokenPath, "sessionName", sessionName, "err", err)
		}

		// Add the web identity role credential provider
		providers = append(providers, webIdentityProvider)
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
		_, checkFullURI := os.LookupEnv("AWS_CONTAINER_CREDENTIALS_FULL_URI")
		_, checkRelativeURI := os.LookupEnv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI")
		if !checkFullURI && !checkRelativeURI {
			// match the sdk defaults from https://github.com/aws/aws-sdk-go/pull/3066
			def.Config.HTTPClient.Timeout = 1 * time.Second
			def.Config.MaxRetries = aws.Int(2)
		}
	}

	providers = append(providers, defaults.RemoteCredProvider(*def.Config, def.Handlers))

	// Create the credentials required to access the API.
	creds := credentials.NewChainCredentials(providers)
	if creds == nil {
		return nil, fmt.Errorf("could not compile valid credential providers from static config, environment, shared, web identity or instance metadata")
	}

	return creds, nil
}

func RetrieveCreds(accessKey, secretKey, sessionToken string, logger hclog.Logger) (*credentials.Credentials, error) {
	credConfig := CredentialsConfig{
		AccessKey:    accessKey,
		SecretKey:    secretKey,
		SessionToken: sessionToken,
		Logger:       logger,
	}
	creds, err := credConfig.GenerateCredentialChain()
	if err != nil {
		return nil, err
	}
	if creds == nil {
		return nil, fmt.Errorf("could not compile valid credential providers from static config, environment, shared, or instance metadata")
	}

	_, err = creds.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve credentials from credential chain: %w", err)
	}
	return creds, nil
}

// GenerateLoginData populates the necessary data to send to the Vault server for generating a token
// This is useful for other API clients to use
func GenerateLoginData(creds *credentials.Credentials, headerValue, configuredRegion string, logger hclog.Logger) (map[string]interface{}, error) {
	loginData := make(map[string]interface{})

	// Use the credentials we've found to construct an STS session
	region, err := GetRegion(configuredRegion)
	if err != nil {
		logger.Warn(fmt.Sprintf("defaulting region to %q due to %s", DefaultRegion, err.Error()))
		region = DefaultRegion
	}
	stsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials:      creds,
			Region:           &region,
			EndpointResolver: endpoints.ResolverFunc(stsSigningResolver),
		},
	})
	if err != nil {
		return nil, err
	}

	var params *sts.GetCallerIdentityInput
	svc := sts.New(stsSession)
	stsRequest, _ := svc.GetCallerIdentityRequest(params)

	// Inject the required auth header value, if supplied, and then sign the request including that header
	if headerValue != "" {
		stsRequest.HTTPRequest.Header.Add(iamServerIdHeader, headerValue)
	}
	stsRequest.Sign()

	// Now extract out the relevant parts of the request
	headersJson, err := json.Marshal(stsRequest.HTTPRequest.Header)
	if err != nil {
		return nil, err
	}
	requestBody, err := ioutil.ReadAll(stsRequest.HTTPRequest.Body)
	if err != nil {
		return nil, err
	}
	loginData["iam_http_request_method"] = stsRequest.HTTPRequest.Method
	loginData["iam_request_url"] = base64.StdEncoding.EncodeToString([]byte(stsRequest.HTTPRequest.URL.String()))
	loginData["iam_request_headers"] = base64.StdEncoding.EncodeToString(headersJson)
	loginData["iam_request_body"] = base64.StdEncoding.EncodeToString(requestBody)

	return loginData, nil
}

// STS is a really weird service that used to only have global endpoints but now has regional endpoints as well.
// For backwards compatibility, even if you request a region other than us-east-1, it'll still sign for us-east-1.
// See, e.g., https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_enable-regions.html#id_credentials_temp_enable-regions_writing_code
// So we have to shim in this EndpointResolver to force it to sign for the right region
func stsSigningResolver(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
	defaultEndpoint, err := endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	if err != nil {
		return defaultEndpoint, err
	}
	defaultEndpoint.SigningRegion = region
	return defaultEndpoint, nil
}
