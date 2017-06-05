package awsauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/awsutil"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (string, error) {
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

	// Grab any supplied credentials off the command line
	// Ensure we're able to fall back to the SDK default credential providers
	credConfig := &awsutil.CredentialsConfig{
		AccessKey:    m["aws_access_key_id"],
		SecretKey:    m["aws_secret_access_key"],
		SessionToken: m["aws_security_token"],
	}
	creds, err := credConfig.GenerateCredentialChain()
	if err != nil {
		return "", err
	}
	if creds == nil {
		return "", fmt.Errorf("could not compile valid credential providers from static config, environment, shared, or instance metadata")
	}

	// Use the credentials we've found to construct an STS session
	stsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Credentials: creds},
	})
	if err != nil {
		return "", err
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
		return "", err
	}
	requestBody, err := ioutil.ReadAll(stsRequest.HTTPRequest.Body)
	if err != nil {
		return "", err
	}
	method := stsRequest.HTTPRequest.Method
	targetUrl := base64.StdEncoding.EncodeToString([]byte(stsRequest.HTTPRequest.URL.String()))
	headers := base64.StdEncoding.EncodeToString(headersJson)
	body := base64.StdEncoding.EncodeToString(requestBody)

	// And pass them on to the Vault server
	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := c.Logical().Write(path, map[string]interface{}{
		"iam_http_request_method": method,
		"iam_request_url":         targetUrl,
		"iam_request_headers":     headers,
		"iam_request_body":        body,
		"role":                    role,
	})

	if err != nil {
		return "", err
	}
	if secret == nil {
		return "", fmt.Errorf("empty response from credential provider")
	}

	return secret.Auth.ClientToken, nil
}

func (h *CLIHandler) Help() string {
	help := `
The AWS credential provider allows you to authenticate with
AWS IAM credentials. To use it, you specify valid AWS IAM credentials
in one of a number of ways. They can be specified explicitly on the
command line (which in general you should not do), via the standard AWS
environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and
AWS_SECURITY_TOKEN), via the ~/.aws/credentials file, or via an EC2
instance profile (in that order).

  Example: vault auth -method=aws

If you need to explicitly pass in credentials, you would do it like this:
  Example: vault auth -method=aws aws_access_key_id=<access key> aws_secret_access_key=<secret key> aws_security_token=<token>

Key/Value Pairs:

  mount=aws                           The mountpoint for the AWS credential provider.
                                      Defaults to "aws"
  aws_access_key_id=<access key>      Explicitly specified AWS access key
  aws_secret_access_key=<secret key>  Explicitly specified AWS secret key
  aws_security_token=<token>          Security token for temporary credentials
  header_value                        The Value of the X-Vault-AWS-IAM-Server-ID header.
  role                                The name of the role you're requesting a token for
  `

	return strings.TrimSpace(help)
}
