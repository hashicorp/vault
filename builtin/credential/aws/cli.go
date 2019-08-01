package awsauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/errwrap"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/awsutil"
)

type CLIHandler struct{}

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

// GenerateLoginData populates the necessary data to send to the Vault server for generating a token
// This is useful for other API clients to use
func GenerateLoginData(creds *credentials.Credentials, headerValue, configuredRegion string) (map[string]interface{}, error) {
	loginData := make(map[string]interface{})

	// Use the credentials we've found to construct an STS session
	region := awsutil.GetOrDefaultRegion(hclog.Default(), configuredRegion)
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

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "aws"
	}

	role, ok := m["role"]
	if !ok {
		role = ""
	}

	headerValue, ok := m["header_value"]
	if !ok {
		headerValue = ""
	}

	creds, err := RetrieveCreds(m["aws_access_key_id"], m["aws_secret_access_key"], m["aws_security_token"])
	if err != nil {
		return nil, err
	}

	loginData, err := GenerateLoginData(creds, headerValue, m["region"])
	if err != nil {
		return nil, err
	}
	if loginData == nil {
		return nil, fmt.Errorf("got nil response from GenerateLoginData")
	}
	loginData["role"] = role
	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := c.Logical().Write(path, loginData)

	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("empty response from credential provider")
	}

	return secret, nil
}

func RetrieveCreds(accessKey, secretKey, sessionToken string) (*credentials.Credentials, error) {
	credConfig := &awsutil.CredentialsConfig{
		AccessKey:    accessKey,
		SecretKey:    secretKey,
		SessionToken: sessionToken,
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
		return nil, errwrap.Wrapf("failed to retrieve credentials from credential chain: {{err}}", err)
	}
	return creds, nil
}

func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=aws [CONFIG K=V...]

  The AWS auth method allows users to authenticate with AWS IAM
  credentials. The AWS IAM credentials may be specified in a number of ways,
  listed in order of precedence below:

    1. Explicitly via the command line (not recommended)

    2. Via the standard AWS environment variables (AWS_ACCESS_KEY, etc.)

    3. Via the ~/.aws/credentials file

    4. Via EC2 instance profile

  Authenticate using locally stored credentials:

      $ vault login -method=aws

  Authenticate by passing keys:

      $ vault login -method=aws aws_access_key_id=... aws_secret_access_key=...

Configuration:

  aws_access_key_id=<string>
      Explicit AWS access key ID

  aws_secret_access_key=<string>
      Explicit AWS secret access key

  aws_security_token=<string>
      Explicit AWS security token for temporary credentials

  header_value=<string>
      Value for the x-vault-aws-iam-server-id header in requests

  mount=<string>
      Path where the AWS credential method is mounted. This is usually provided
      via the -path flag in the "vault login" command, but it can be specified
      here as well. If specified here, it takes precedence over the value for
      -path. The default value is "aws".

  role=<string>
      Name of the role to request a token against
`

	return strings.TrimSpace(help)
}
