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
		return "", fmt.Errorf("could not compile valid credential providers from static config, environemnt, shared, or instance metadata")
	}

	stsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Credentials: creds},
	})
	if err != nil {
		return "", err
	}

	var params *sts.GetCallerIdentityInput
	svc := sts.New(stsSession)
	stsRequest, _ := svc.GetCallerIdentityRequest(params)
	if headerValue != "" {
		stsRequest.HTTPRequest.Header.Add(magicVaultHeader, headerValue)
	}
	stsRequest.Sign()
	headersJson, err := json.Marshal(stsRequest.HTTPRequest.Header)
	if err != nil {
		return "", err
	}
	requestBody, err := ioutil.ReadAll(stsRequest.HTTPRequest.Body)
	if err != nil {
		return "", err
	}
	method := stsRequest.HTTPRequest.Method
	targetUrl := stsRequest.HTTPRequest.URL.String()
	headers := base64.StdEncoding.EncodeToString(headersJson)
	body := base64.StdEncoding.EncodeToString(requestBody)

	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := c.Logical().Write(path, map[string]interface{}{
		"auth_type":       "iam",
		"request_method":  method,
		"request_url":     targetUrl,
		"request_headers": headers,
		"request_body":    body,
		"role":            role,
	})

	if err != nil {
		return "", err
	}
	if secret == nil {
		return "", fmt.Errorf("empty response from credential provider")
	}

	return secret.Auth.ClientToken, nil

	return "", nil
}

func (h *CLIHandler) Help() string {
	help := `
The AWS credentaial provider allows you to authenticate with
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
                                      Defaults to "aws-iam"
  aws_access_key_id=<access key>      Explicitly specified AWS access key
  aws_secret_access_key=<secret key>  Explicitly specified AWS secret key
  aws_security_token=<token>          Security token for temporary credentials
  header_value                        The Value of the X-Vault-AWSIAM-Server-ID header.
  role                                The name of the role you're requesting a token for
  `

	return strings.TrimSpace(help)
}
