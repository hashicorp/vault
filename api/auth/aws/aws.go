// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package aws

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/go-hclog"
	awsutil "github.com/hashicorp/go-secure-stdlib/awsutil/v2"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
)

type AWSAuth struct {
	// If not provided with the WithRole login option, the Vault server will look for a role
	// with the friendly name of the IAM principal if using the IAM auth type,
	// or the name of the EC2 instance's AMI ID if using the EC2 auth type.
	// If no matching role is found, login will fail.
	roleName  string
	mountPath string
	// Can be "iam" or "ec2". Defaults to "iam".
	authType string
	// Can be "pkcs7", "identity", or "rsa2048". Defaults to "pkcs7".
	signatureType          string
	region                 string
	iamServerIDHeaderValue string
	awsConfig              *awsv2.Config
	nonce                  string
}

var _ api.AuthMethod = (*AWSAuth)(nil)

type LoginOption func(a *AWSAuth) error

const (
	iamType              = "iam"
	ec2Type              = "ec2"
	pkcs7Type            = "pkcs7"
	identityType         = "identity"
	rsa2048Type          = "rsa2048"
	defaultMountPath     = "aws"
	defaultAuthType      = iamType
	defaultRegion        = "us-east-1"
	defaultSignatureType = pkcs7Type
)

// NewAWSAuth initializes a new AWS auth method interface to be
// passed as a parameter to the client.Auth().Login method.
//
// Supported options: WithRole, WithMountPath, WithIAMAuth, WithEC2Auth,
// WithPKCS7Signature, WithIdentitySignature, WithRSA2048Signature, WithIAMServerIDHeader, WithNonce, WithRegion
func NewAWSAuth(opts ...LoginOption) (*AWSAuth, error) {
	a := &AWSAuth{
		mountPath:     defaultMountPath,
		authType:      defaultAuthType,
		region:        defaultRegion,
		signatureType: defaultSignatureType,
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *AWSAuth as the argument
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("error with login option: %w", err)
		}
	}

	// return the modified auth struct instance
	return a, nil
}

// Login sets up the required request body for the AWS auth method's /login
// endpoint, and performs a write to it. This method defaults to the "iam"
// auth type unless NewAWSAuth is called with WithEC2Auth().
//
// The Vault client will set its credentials to the values of the
// AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_SESSION environment
// variables. To specify a path to a credentials file on disk instead, set
// the environment variable AWS_SHARED_CREDENTIALS_FILE.
func (a *AWSAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	loginData := make(map[string]interface{})
	switch a.authType {
	case ec2Type:
		metadataSvc := imds.New(imds.Options{})

		getDynamicData := func(path string) (string, error) {
			out, err := metadataSvc.GetDynamicData(ctx, &imds.GetDynamicDataInput{Path: path})
			if err != nil {
				return "", err
			}
			defer out.Content.Close()
			data, err := io.ReadAll(out.Content)
			if err != nil {
				return "", err
			}
			return string(data), nil
		}

		if a.signatureType == pkcs7Type {
			// fetch PKCS #7 signature
			resp, err := getDynamicData("instance-identity/pkcs7")
			if err != nil {
				return nil, fmt.Errorf("unable to get PKCS 7 data from metadata service: %w", err)
			}
			loginData["pkcs7"] = strings.TrimSpace(resp)
		} else if a.signatureType == identityType {
			// fetch signature from identity document
			doc, err := getDynamicData("instance-identity/document")
			if err != nil {
				return nil, fmt.Errorf("error requesting instance identity doc: %w", err)
			}
			loginData["identity"] = base64.StdEncoding.EncodeToString([]byte(doc))

			signature, err := getDynamicData("instance-identity/signature")
			if err != nil {
				return nil, fmt.Errorf("error requesting signature: %w", err)
			}
			loginData["signature"] = signature
		} else if a.signatureType == rsa2048Type {
			// fetch RSA 2048 signature, which is also a PKCS#7 signature
			resp, err := getDynamicData("instance-identity/rsa2048")
			if err != nil {
				return nil, fmt.Errorf("unable to get PKCS 7 data from metadata service: %w", err)
			}
			loginData["pkcs7"] = strings.TrimSpace(resp)
		} else {
			return nil, fmt.Errorf("unknown signature type: %s", a.signatureType)
		}

		// Add the reauthentication value, if we have one
		if a.nonce == "" {
			uid, err := uuid.GenerateUUID()
			if err != nil {
				return nil, fmt.Errorf("error generating uuid for reauthentication value: %w", err)
			}
			a.nonce = uid
		}
		loginData["nonce"] = a.nonce
	case iamType:
		logger := hclog.Default()
		if a.awsConfig == nil {
			credsConfig := awsutil.CredentialsConfig{
				AccessKey:    os.Getenv("AWS_ACCESS_KEY_ID"),
				SecretKey:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
				SessionToken: os.Getenv("AWS_SESSION_TOKEN"),
				Logger:       logger,
			}

			// the env vars above will take precedence if they are set, as
			// they will be added to the ChainProvider stack first
			var hasCredsFile bool
			credsFilePath := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
			if credsFilePath != "" {
				hasCredsFile = true
				credsConfig.Filename = credsFilePath
			}

			cfg, err := credsConfig.GenerateCredentialChain(ctx, awsutil.WithSharedCredentials(hasCredsFile))
			if err != nil {
				return nil, err
			}

			if _, err := cfg.Credentials.Retrieve(ctx); err != nil {
				return nil, fmt.Errorf("failed to retrieve credentials from credential chain: %w", err)
			}

			a.awsConfig = cfg
		}

		data, err := GenerateLoginDataV2(ctx, a.awsConfig, a.iamServerIDHeaderValue, a.region, logger)
		if err != nil {
			return nil, fmt.Errorf("unable to generate login data for AWS auth endpoint: %w", err)
		}
		loginData = data
	}

	// Add role if we have one. If not, Vault will infer the role name based
	// on the IAM friendly name (iam auth type) or EC2 instance's
	// AMI ID (ec2 auth type).
	if a.roleName != "" {
		loginData["role"] = a.roleName
	}

	if a.iamServerIDHeaderValue != "" {
		client.AddHeader("iam_server_id_header_value", a.iamServerIDHeaderValue)
	}

	path := fmt.Sprintf("auth/%s/login", a.mountPath)
	resp, err := client.Logical().WriteWithContext(ctx, path, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with AWS auth: %w", err)
	}

	return resp, nil
}

func WithRole(roleName string) LoginOption {
	return func(a *AWSAuth) error {
		a.roleName = roleName
		return nil
	}
}

func WithMountPath(mountPath string) LoginOption {
	return func(a *AWSAuth) error {
		a.mountPath = mountPath
		return nil
	}
}

func WithEC2Auth() LoginOption {
	return func(a *AWSAuth) error {
		a.authType = ec2Type
		return nil
	}
}

func WithIAMAuth() LoginOption {
	return func(a *AWSAuth) error {
		a.authType = iamType
		return nil
	}
}

// WithIdentitySignature will have the client send the cryptographic identity
// document signature to verify EC2 auth logins. Only used by EC2 auth type.
// If this option is not provided, will default to using the PKCS #7 signature.
// The signature type used should match the type of the public AWS cert Vault
// has been configured with to verify EC2 instance identity.
// https://developer.hashicorp.com/vault/api-docs/auth/aws#create-certificate-configuration
func WithIdentitySignature() LoginOption {
	return func(a *AWSAuth) error {
		a.signatureType = identityType
		return nil
	}
}

// WithPKCS7Signature will explicitly tell the client to send the PKCS #7
// signature to verify EC2 auth logins. Only used by EC2 auth type.
// PKCS #7 is the default, but this method is provided for additional clarity.
// The signature type used should match the type of the public AWS cert Vault
// has been configured with to verify EC2 instance identity.
// https://developer.hashicorp.com/vault/api-docs/auth/aws#create-certificate-configuration
func WithPKCS7Signature() LoginOption {
	return func(a *AWSAuth) error {
		a.signatureType = pkcs7Type
		return nil
	}
}

// WithRSA2048Signature will explicitly tell the client to send the RSA2048
// signature to verify EC2 auth logins. Only used by EC2 auth type.
// If this option is not provided, will default to using the PKCS #7 signature.
// The signature type used should match the type of the public AWS cert Vault
// has been configured with to verify EC2 instance identity.
// https://www.vaultproject.io/api/auth/aws#create-certificate-configuration
func WithRSA2048Signature() LoginOption {
	return func(a *AWSAuth) error {
		a.signatureType = rsa2048Type
		return nil
	}
}

func WithIAMServerIDHeader(headerValue string) LoginOption {
	return func(a *AWSAuth) error {
		a.iamServerIDHeaderValue = headerValue
		return nil
	}
}

// WithNonce can be used to specify a named nonce for the ec2 auth login
// method. If not provided, an automatically-generated uuid will be used
// instead.
func WithNonce(nonce string) LoginOption {
	return func(a *AWSAuth) error {
		a.nonce = nonce
		return nil
	}
}

func WithRegion(region string) LoginOption {
	return func(a *AWSAuth) error {
		a.region = region
		return nil
	}
}

// iamServerIdHeader is the name of the header that the client signs and the
// server validates so that a signed GetCallerIdentity request cannot be
// replayed against an unintended Vault server.
const iamServerIdHeader = "X-Vault-AWS-IAM-Server-ID"

// GenerateLoginDataV2 builds the login payload for the AWS IAM auth method using
// the AWS SDK for Go v2. It signs a GetCallerIdentity STS request with the
// credentials resolved from cfg and returns the base64-encoded request
// components expected by the auth/aws login endpoint.
//
// This replaces the v1 awsutil.GenerateLoginData helper that was dropped during
// the SDK v2 upgrade. It is kept local to this module so the public AWS auth
// SDK does not depend on Vault's internal packages.
//
// NOTE: GenerateLoginDataV2, STSLoginEndpoint, and STSRegionalEndpoint are
// intentionally duplicated in internal/awsutil/v2/generate_credentials.go, which
// the server-side aws auth backend and login CLI use. This api module
// (github.com/hashicorp/vault/api) cannot import that package -- it is a
// separate module and Go's internal/ visibility rule forbids it -- so the two
// copies must be kept in sync until the logic is upstreamed to a shared library.
func GenerateLoginDataV2(ctx context.Context, cfg *awsv2.Config, headerValue, configuredRegion string, logger hclog.Logger) (map[string]interface{}, error) {
	if cfg == nil {
		return nil, fmt.Errorf("aws config must not be nil")
	}
	if cfg.Credentials == nil {
		return nil, fmt.Errorf("aws config credentials must not be nil")
	}
	loginData := make(map[string]interface{})

	region, err := awsutil.GetRegion(ctx, configuredRegion)
	if err != nil {
		logger.Warn(fmt.Sprintf("defaulting region to %q due to %s", awsutil.DefaultRegion, err.Error()))
		region = awsutil.DefaultRegion
	}

	body := []byte("Action=GetCallerIdentity&Version=2011-06-15")

	stsURL, err := STSLoginEndpoint(ctx, region)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, stsURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	if headerValue != "" {
		req.Header.Set(iamServerIdHeader, headerValue)
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, err
	}

	payloadHash := sha256.Sum256(body)
	signer := v4.NewSigner()
	if err := signer.SignHTTP(ctx, creds, req, hex.EncodeToString(payloadHash[:]), "sts", region, time.Now()); err != nil {
		return nil, fmt.Errorf("failed to sign STS request: %w", err)
	}

	headersJSON, err := json.Marshal(req.Header)
	if err != nil {
		return nil, err
	}

	loginData["iam_http_request_method"] = req.Method
	loginData["iam_request_url"] = base64.StdEncoding.EncodeToString([]byte(req.URL.String()))
	loginData["iam_request_headers"] = base64.StdEncoding.EncodeToString(headersJSON)
	loginData["iam_request_body"] = base64.StdEncoding.EncodeToString(body)

	return loginData, nil
}

// STSLoginEndpoint returns the STS endpoint URL whose host matches the region
// the login request is signed for. The global endpoint is retained for the
// default region; any other region resolves to a regional endpoint so the
// signed Host matches the request region (this also handles non-default
// partitions such as AWS China and GovCloud).
func STSLoginEndpoint(ctx context.Context, region string) (string, error) {
	if region == awsutil.DefaultRegion {
		return "https://sts.amazonaws.com/", nil
	}
	regional, err := STSRegionalEndpoint(ctx, region)
	if err != nil {
		return "", err
	}
	// The SDK resolver may or may not return a trailing slash; trim it before
	// re-appending so the endpoint never ends up with a double slash.
	return strings.TrimRight(regional, "/") + "/", nil
}

// STSRegionalEndpoint resolves the regional STS endpoint URL for the given
// region using the AWS SDK v2 endpoint resolver, accounting for non-default
// partitions.
func STSRegionalEndpoint(ctx context.Context, region string) (string, error) {
	resolver := sts.NewDefaultEndpointResolverV2()
	resolvedEndpoint, err := resolver.ResolveEndpoint(ctx, sts.EndpointParameters{
		Region:            awsv2.String(region),
		UseDualStack:      awsv2.Bool(false),
		UseFIPS:           awsv2.Bool(false),
		UseGlobalEndpoint: awsv2.Bool(false),
	})
	if err != nil {
		return "", fmt.Errorf("unable to get regional STS endpoint for region: %v", region)
	}
	return resolvedEndpoint.URI.String(), nil
}
