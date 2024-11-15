// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/go-cleanhttp"
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

	// The IAM endpoint to use; if not set will use the default
	IAMEndpoint string

	// The STS endpoint to use; if not set will use the default
	STSEndpoint string

	// If specified, the region will be provided to the config of the
	// EC2RoleProvider's client. This may be useful if you want to e.g. reuse
	// the client elsewhere. If not specified, the region will be determined
	// by the environment variables "AWS_REGION" or "AWS_DEFAULT_REGION".
	// Otherwise the default value is us-east-1.
	Region string

	// The filename for the shared credentials provider, if being used
	Filename string

	// The profile for the shared credentials provider, if being used
	Profile string

	// The role arn to use when creating either a web identity role provider
	// or a ec2-instance role provider.
	RoleARN string

	// The role session name to use when creating either a web identity role provider
	// or a ec2-instance role provider.
	RoleSessionName string

	// The role external ID to use when creating a ec2-instance role provider.
	RoleExternalId string

	// The role tags to use when creating a ec2-instance role provider.
	RoleTags map[string]string

	// The web identity token file to use if using the web identity token provider
	WebIdentityTokenFile string

	// The web identity token (contents, not the file path) to use with the web
	// identity token provider
	WebIdentityToken string

	// The web identity token fetcher to use with the web identity token provider
	WebIdentityTokenFetcher stscreds.TokenFetcher

	// The http.Client to use, or nil for the client to use its default
	HTTPClient *http.Client

	// The max retries to set on the client. This is a pointer because the zero
	// value has meaning. A nil pointer will use the default value.
	MaxRetries *int

	// The logger to use for credential acquisition debugging
	Logger hclog.Logger
}

// NewCredentialsConfig creates a CredentialsConfig with the provided options.
//
// Supported options: WithAccessKey, WithSecretKey, WithLogger, WithStsEndpoint,
// WithIamEndpoint, WithMaxRetries, WithRegion, WithHttpClient, WithRoleArn,
// WithRoleSessionName, WithRoleExternalId, WithRoleTags, WithWebIdentityTokenFile,
// WithWebIdentityToken, WithWebIdentityTokenFetcher.
func NewCredentialsConfig(opt ...Option) (*CredentialsConfig, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("error reading options in NewCredentialsConfig: %w", err)
	}

	c := &CredentialsConfig{
		AccessKey:      opts.withAccessKey,
		SecretKey:      opts.withSecretKey,
		Logger:         opts.withLogger,
		STSEndpoint:    opts.withStsEndpoint,
		IAMEndpoint:    opts.withIamEndpoint,
		MaxRetries:     opts.withMaxRetries,
		RoleExternalId: opts.withRoleExternalId,
		RoleTags:       opts.withRoleTags,
	}

	c.Region = opts.withRegion
	if c.Region == "" {
		c.Region = os.Getenv("AWS_REGION")
		if c.Region == "" {
			c.Region = os.Getenv("AWS_DEFAULT_REGION")
			if c.Region == "" {
				c.Region = "us-east-1"
			}
		}
	}
	c.RoleARN = opts.withRoleArn
	if c.RoleARN == "" {
		c.RoleARN = os.Getenv("AWS_ROLE_ARN")
	}
	c.RoleSessionName = opts.withRoleSessionName
	if c.RoleSessionName == "" {
		c.RoleSessionName = os.Getenv("AWS_ROLE_SESSION_NAME")
	}
	c.WebIdentityTokenFile = opts.withWebIdentityTokenFile
	if c.WebIdentityTokenFile == "" {
		c.WebIdentityTokenFile = os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE")
	}
	c.WebIdentityToken = opts.withWebIdentityToken

	c.WebIdentityTokenFetcher = opts.withWebIdentityTokenFetcher

	if c.RoleARN == "" {
		if c.RoleSessionName != "" {
			return nil, fmt.Errorf("role session name specified without role ARN")
		}
		if c.RoleExternalId != "" {
			return nil, fmt.Errorf("role external ID specified without role ARN")
		}
		if len(c.RoleTags) > 0 {
			return nil, fmt.Errorf("role tags specified without role ARN")
		}
		if c.WebIdentityTokenFile != "" {
			return nil, fmt.Errorf("web identity token file specified without role ARN")
		}
		if len(c.WebIdentityToken) > 0 {
			return nil, fmt.Errorf("web identity token specified without role ARN")
		}

		if c.WebIdentityTokenFetcher != nil {
			return nil, fmt.Errorf("web identity token fetcher specified without role ARN")
		}
	}

	c.HTTPClient = opts.withHttpClient
	if c.HTTPClient == nil {
		c.HTTPClient = cleanhttp.DefaultClient()
	}

	return c, nil
}

// Make sure the logger isn't nil before logging
func (c *CredentialsConfig) log(level hclog.Level, msg string, args ...interface{}) {
	if c.Logger != nil {
		c.Logger.Log(level, msg, args...)
	}
}

// GenerateCredentialChain uses the config to generate a credential chain
// suitable for creating AWS sessions and clients.
//
// Supported options: WithEnvironmentCredentials, WithSharedCredentials,
// WithSkipWebIdentityValidity
func (c *CredentialsConfig) GenerateCredentialChain(opt ...Option) (*credentials.Credentials, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("error reading options in GenerateCredentialChain: %w", err)
	}

	var providers []credentials.Provider

	// Have one or the other but not both and not neither
	if (c.AccessKey != "" && c.SecretKey == "") || (c.AccessKey == "" && c.SecretKey != "") {
		return nil, fmt.Errorf("static AWS client credentials haven't been properly configured (the access key or secret key were provided but not both)")
	}
	// Add the static credential provider
	if c.AccessKey != "" && c.SecretKey != "" {
		providers = append(providers, &credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     c.AccessKey,
				SecretAccessKey: c.SecretKey,
				SessionToken:    c.SessionToken,
			},
		})
		c.log(hclog.Debug, "added static credential provider", "AccessKey", c.AccessKey)
	}

	// Add the environment credential provider
	if opts.withEnvironmentCredentials {
		providers = append(providers, &credentials.EnvProvider{})
		c.log(hclog.Debug, "added environment variable credential provider")
	}

	// Add the shared credentials provider
	if opts.withSharedCredentials {
		profile := os.Getenv("AWS_PROFILE")
		if profile != "" {
			c.Profile = profile
		}
		if c.Profile == "" {
			c.Profile = "default"
		}
		providers = append(providers, &credentials.SharedCredentialsProvider{
			Filename: c.Filename,
			Profile:  c.Profile,
		})
		c.log(hclog.Debug, "added shared credential provider")
	}

	// Add the assume role provider
	roleARN := c.RoleARN
	if roleARN == "" {
		roleARN = os.Getenv("AWS_ROLE_ARN")
	}
	tokenPath := c.WebIdentityTokenFile
	if tokenPath == "" {
		tokenPath = os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE")
	}
	roleSessionName := c.RoleSessionName
	if roleSessionName == "" {
		roleSessionName = os.Getenv("AWS_ROLE_SESSION_NAME")
	}
	if roleARN != "" {
		if tokenPath != "" || c.WebIdentityToken != "" || c.WebIdentityTokenFetcher != nil {
			sess, err := session.NewSession()
			if err != nil {
				return nil, errors.Wrap(err, "error creating a new session to create a WebIdentityRoleProvider")
			}

			var webIdentityProvider *stscreds.WebIdentityRoleProvider
			switch {
			case tokenPath != "":
				c.log(hclog.Debug, "adding web identity provider", "roleARN", roleARN)
				webIdentityProvider = stscreds.NewWebIdentityRoleProvider(sts.New(sess), roleARN, roleSessionName, tokenPath)
			case c.WebIdentityToken != "":
				c.log(hclog.Debug, "adding web identity provider with token", "roleARN", roleARN)
				webIdentityProvider = stscreds.NewWebIdentityRoleProviderWithToken(sts.New(sess), roleARN, roleSessionName, FetchTokenContents(c.WebIdentityToken))
			case c.WebIdentityTokenFetcher != nil:
				c.log(hclog.Debug, "adding web identity provider with token fetcher", "roleARN", roleARN)
				webIdentityProvider = stscreds.NewWebIdentityRoleProviderWithToken(sts.New(sess), roleARN, roleSessionName, c.WebIdentityTokenFetcher)
			}

			if opts.withSkipWebIdentityValidity {
				// Add the web identity role credential provider without
				// generating credentials to check validity first
				providers = append(providers, webIdentityProvider)
			} else {
				// Check if the webIdentityProvider can successfully retrieve
				// credentials (via sts:AssumeRole), and warn if there's a problem.
				if _, err := webIdentityProvider.Retrieve(); err != nil {
					c.log(hclog.Warn, "error assuming role with web identity", "roleARN", roleARN, "sessionName", roleSessionName, "err", err)
				} else {
					// Add the web identity role credential provider
					providers = append(providers, webIdentityProvider)
				}
			}
		} else {
			// this session is only created to create the AssumeRoleProvider, variables used to
			// assume a role are pulled from values provided in options. If the option values are
			// not set, then the provider will default to using the environment variables.
			c.log(hclog.Debug, "adding ec2-instance role provider", "roleARN", roleARN)
			sess, err := session.NewSession()
			if err != nil {
				return nil, errors.Wrap(err, "error creating a new session for ec2 instance role credentials")
			}
			assumedRoleCredentials := stscreds.NewCredentials(sess, roleARN, func(p *stscreds.AssumeRoleProvider) {
				p.RoleSessionName = roleSessionName
				if c.RoleExternalId != "" {
					p.ExternalID = aws.String(c.RoleExternalId)
				}
				if len(c.RoleTags) != 0 {
					p.Tags = []*sts.Tag{}
					for k, v := range c.RoleTags {
						p.Tags = append(p.Tags, &sts.Tag{
							Key:   aws.String(k),
							Value: aws.String(v),
						})
					}
				}
			})
			// Check if the credentials are successfully retrieved
			// (via sts:AssumeRole), and warn if there's a problem.
			creds, err := assumedRoleCredentials.Get()
			if err != nil {
				c.log(hclog.Warn, "error assuming role", "roleARN", roleARN, "sessionName", roleSessionName, "err", err)
			} else {
				providers = append(providers, &credentials.StaticProvider{
					Value: creds,
				})
			}
		}
	}

	// Add the remote provider
	def := defaults.Get()
	if c.Region != "" {
		def.Config.Region = aws.String(c.Region)
	}
	// We are taking care of this in the New() function but for legacy reasons
	// we also set this here
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

// FetchTokenContents allows the use of the content of a token in the
// WebIdentityProvider, instead of the path to a token. Useful with a
// serviceaccount token requested directly from the EKS/K8s API, for example.
type FetchTokenContents []byte

var _ stscreds.TokenFetcher = (*FetchTokenContents)(nil)

func (f FetchTokenContents) FetchToken(_ aws.Context) ([]byte, error) {
	return f, nil
}
