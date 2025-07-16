// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package aws

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
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
	creds                  aws.CredentialsProvider
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
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(a.region))
		if err != nil {
			return nil, fmt.Errorf("error loading AWS config: %w", err)
		}
		metadataSvc := imds.NewFromConfig(cfg)

		var path string
		switch a.signatureType {
		case pkcs7Type:
			path = "/instance-identity/pkcs7"
		case identityType:
			path = "/instance-identity/document"
		case rsa2048Type:
			path = "/instance-identity/rsa2048"
		default:
			return nil, fmt.Errorf("unknown signature type: %s", a.signatureType)
		}

		resp, err := metadataSvc.GetDynamicData(ctx, &imds.GetDynamicDataInput{Path: path})
		if err != nil {
			return nil, fmt.Errorf("unable to get identity data: %w", err)
		}
		defer resp.Content.Close()
		body, err := io.ReadAll(resp.Content)
		if err != nil {
			return nil, fmt.Errorf("error reading identity data: %w", err)
		}
		pkcs7 := strings.TrimSpace(string(body))
		loginData["pkcs7"] = pkcs7

		if a.nonce == "" {
			uuid, err := uuid.GenerateUUID()
			if err != nil {
				return nil, fmt.Errorf("error generating uuid: %w", err)
			}
			a.nonce = uuid
		}
		loginData["nonce"] = a.nonce
	case iamType:
		if a.creds == nil {
			cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(a.region))
			if err != nil {
				return nil, fmt.Errorf("unable to load AWS config: %w", err)
			}
			a.creds = cfg.Credentials
		}

		credsVal, err := a.creds.Retrieve(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve credentials: %w", err)
		}

		const iamBody = "Action=GetCallerIdentity&Version=2011-06-15"
		req, err := http.NewRequest("POST", "https://sts.amazonaws.com/", strings.NewReader(iamBody))
		if err != nil {
			return nil, fmt.Errorf("failed to construct STS request: %w", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		hash := sha256.Sum256([]byte(iamBody))
		payloadHash := hex.EncodeToString(hash[:])

		signer := v4.NewSigner()
		err = signer.SignHTTP(ctx, credsVal, req, payloadHash, "sts", a.region, time.Now().UTC())
		if err != nil {
			return nil, fmt.Errorf("failed to sign STS request: %w", err)
		}

		headersData, _ := json.Marshal(req.Header)

		loginData["iam_http_request_method"] = "POST"
		loginData["iam_request_url"] = base64.StdEncoding.EncodeToString([]byte(req.URL.String()))
		loginData["iam_request_body"] = base64.StdEncoding.EncodeToString([]byte(iamBody))
		loginData["iam_request_headers"] = base64.StdEncoding.EncodeToString(headersData)
	}

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
