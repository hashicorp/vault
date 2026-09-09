// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	internal_awsutilv2 "github.com/hashicorp/vault/internal/awsutil/v2"
)

// TestNewAWSAuthMethod_NilConfig verifies that NewAWSAuthMethod returns an error
// when called with a nil config.
func TestNewAWSAuthMethod_NilConfig(t *testing.T) {
	_, err := NewAWSAuthMethod(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
	if err.Error() != "empty config" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestNewAWSAuthMethod_NilConfigData verifies that NewAWSAuthMethod returns an error
// when the config has nil config data.
func TestNewAWSAuthMethod_NilConfigData(t *testing.T) {
	conf := &auth.AuthConfig{
		Logger: hclog.NewNullLogger(),
	}
	_, err := NewAWSAuthMethod(conf)
	if err == nil {
		t.Fatal("expected error for nil config data")
	}
	if err.Error() != "empty config data" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestNewAWSAuthMethod_MissingType verifies that NewAWSAuthMethod returns an error
// when the 'type' field is missing from the config.
func TestNewAWSAuthMethod_MissingType(t *testing.T) {
	conf := &auth.AuthConfig{
		Logger: hclog.NewNullLogger(),
		Config: map[string]interface{}{
			"role": "test-role",
		},
	}
	_, err := NewAWSAuthMethod(conf)
	if err == nil {
		t.Fatal("expected error for missing type")
	}
	if err.Error() != "missing 'type' value" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestNewAWSAuthMethod_MissingRole verifies that NewAWSAuthMethod returns an error
// when the 'role' field is missing from the config.
func TestNewAWSAuthMethod_MissingRole(t *testing.T) {
	conf := &auth.AuthConfig{
		Logger: hclog.NewNullLogger(),
		Config: map[string]interface{}{
			"type": "iam",
		},
	}
	_, err := NewAWSAuthMethod(conf)
	if err == nil {
		t.Fatal("expected error for missing role")
	}
	if err.Error() != "missing 'role' value" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestNewAWSAuthMethod_EmptyRole verifies that NewAWSAuthMethod returns an error
// when the 'role' field is empty.
func TestNewAWSAuthMethod_EmptyRole(t *testing.T) {
	conf := &auth.AuthConfig{
		Logger: hclog.NewNullLogger(),
		Config: map[string]interface{}{
			"type": "iam",
			"role": "",
		},
	}
	_, err := NewAWSAuthMethod(conf)
	if err == nil {
		t.Fatal("expected error for empty role")
	}
	if err.Error() != "'role' value is empty" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestNewAWSAuthMethod_EmptyType verifies that NewAWSAuthMethod returns an error
// when the 'type' field is empty.
func TestNewAWSAuthMethod_EmptyType(t *testing.T) {
	conf := &auth.AuthConfig{
		Logger: hclog.NewNullLogger(),
		Config: map[string]interface{}{
			"type": "",
			"role": "test-role",
		},
	}
	_, err := NewAWSAuthMethod(conf)
	if err == nil {
		t.Fatal("expected error for empty type")
	}
	if err.Error() != "'type' value is empty" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestNewAWSAuthMethod_InvalidType verifies that NewAWSAuthMethod returns an error
// when the 'type' field contains an invalid value.
func TestNewAWSAuthMethod_InvalidType(t *testing.T) {
	conf := &auth.AuthConfig{
		Logger: hclog.NewNullLogger(),
		Config: map[string]interface{}{
			"type": "invalid",
			"role": "test-role",
		},
	}
	_, err := NewAWSAuthMethod(conf)
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
	if err.Error() != "'type' value is invalid" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestNewAWSAuthMethod_EC2Type verifies that NewAWSAuthMethod successfully creates
// an AWS auth method with EC2 type and validates the returned method's properties.
func TestNewAWSAuthMethod_EC2Type(t *testing.T) {
	conf := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/aws",
		Config: map[string]interface{}{
			"type": "ec2",
			"role": "test-role",
		},
	}
	method, err := NewAWSAuthMethod(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method == nil {
		t.Fatal("expected non-nil method")
	}

	awsMethod, ok := method.(*awsMethod)
	if !ok {
		t.Fatal("expected awsMethod type")
	}
	if awsMethod.authType != typeEC2 {
		t.Fatalf("expected authType %s, got %s", typeEC2, awsMethod.authType)
	}
	if awsMethod.role != "test-role" {
		t.Fatalf("expected role test-role, got %s", awsMethod.role)
	}
}

// TestNewAWSAuthMethod_IAMType verifies that NewAWSAuthMethod successfully creates
// an AWS auth method with IAM type and initializes its AWS config.
// Fake credentials are injected via environment variables so the fail-fast probe
// succeeds without requiring a real AWS environment.
func TestNewAWSAuthMethod_IAMType(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIDTEST")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "SECRETTEST")
	t.Setenv("AWS_SESSION_TOKEN", "")

	conf := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/aws",
		Config: map[string]interface{}{
			"type": "iam",
			"role": "test-role",
		},
	}

	method, err := NewAWSAuthMethod(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method == nil {
		t.Fatal("expected non-nil method")
	}

	awsMethod, ok := method.(*awsMethod)
	if !ok {
		t.Fatal("expected awsMethod type")
	}
	defer awsMethod.Shutdown()

	if awsMethod.authType != typeIAM {
		t.Fatalf("expected authType %s, got %s", typeIAM, awsMethod.authType)
	}
	if awsMethod.role != "test-role" {
		t.Fatalf("expected role test-role, got %s", awsMethod.role)
	}
	if awsMethod.lastCfg == nil {
		t.Fatal("expected non-nil aws config for iam auth type")
	}
}

// TestNewAWSAuthMethod_IAMDynamicCredsFailFast verifies that NewAWSAuthMethod returns
// an error immediately when no credentials are available via the dynamic chain (env
// vars, shared config, IMDS/ECS). This restores the v1 SDK fail-fast behavior so
// a misconfigured instance is caught at startup rather than on the first Authenticate
// call.
func TestNewAWSAuthMethod_IAMDynamicCredsFailFast(t *testing.T) {
	// Remove all credential-bearing env vars and disable IMDS/ECS so the SDK
	// credential chain has no valid provider to fall back to.
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	t.Setenv("AWS_SESSION_TOKEN", "")
	t.Setenv("AWS_PROFILE", "__no_such_profile__")
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", "/dev/null/no-such-file")
	t.Setenv("AWS_CONFIG_FILE", "/dev/null/no-such-file")
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	t.Setenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", "")
	t.Setenv("AWS_CONTAINER_CREDENTIALS_FULL_URI", "")

	conf := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/aws",
		Config: map[string]interface{}{
			"type": "iam",
			"role": "test-role",
		},
	}

	_, err := NewAWSAuthMethod(conf)
	if err == nil {
		t.Fatal("expected error when no credentials are available")
	}
}

// TestNewAWSAuthMethod_IAMWithStaticCreds verifies that NewAWSAuthMethod succeeds
// with explicit access_key/secret_key credentials. Static credentials bypass the
// dynamic-chain probe, so this must work even when IMDS is disabled and no ambient
// credentials exist.
func TestNewAWSAuthMethod_IAMWithStaticCreds(t *testing.T) {
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", "/dev/null/no-such-file")
	t.Setenv("AWS_CONFIG_FILE", "/dev/null/no-such-file")

	conf := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/aws",
		Config: map[string]interface{}{
			"type":       "iam",
			"role":       "test-role",
			"access_key": "AKIDSTATIC",
			"secret_key": "SECRETSTATIC",
		},
	}

	method, err := NewAWSAuthMethod(conf)
	if err != nil {
		t.Fatalf("unexpected error with static credentials: %v", err)
	}
	if method == nil {
		t.Fatal("expected non-nil method")
	}

	awsMethod, ok := method.(*awsMethod)
	if !ok {
		t.Fatal("expected awsMethod type")
	}
	defer awsMethod.Shutdown()

	if awsMethod.lastCfg == nil {
		t.Fatal("expected non-nil aws config")
	}
	if awsMethod.lastCfg.Credentials == nil {
		t.Fatal("expected non-nil credentials on aws config")
	}
}

// TestNewAWSAuthMethod_IAMTypeWithInvalidCredentialPollInterval verifies that
// NewAWSAuthMethod returns an error when the credential_poll_interval is negative.
func TestNewAWSAuthMethod_IAMTypeWithInvalidCredentialPollInterval(t *testing.T) {
	conf := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/aws",
		Config: map[string]interface{}{
			"type":                     "iam",
			"role":                     "test-role",
			"credential_poll_interval": -1,
		},
	}
	_, err := NewAWSAuthMethod(conf)
	if err == nil {
		t.Fatal("expected error for invalid credential_poll_interval")
	}
	if err.Error() != "could not convert 'credential_poll_interval' into positive int" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestNewAWSAuthMethod_WithOptionalFields verifies that NewAWSAuthMethod correctly
// handles optional configuration fields like nonce, header_value, and region.
func TestNewAWSAuthMethod_WithOptionalFields(t *testing.T) {
	conf := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/aws",
		Config: map[string]interface{}{
			"type":         "ec2",
			"role":         "test-role",
			"nonce":        "test-nonce",
			"header_value": "test-header",
			"region":       "us-west-2",
		},
	}
	method, err := NewAWSAuthMethod(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	awsMethod, ok := method.(*awsMethod)
	if !ok {
		t.Fatal("expected awsMethod type")
	}
	if awsMethod.nonce != "test-nonce" {
		t.Fatalf("expected nonce test-nonce, got %s", awsMethod.nonce)
	}
	if awsMethod.headerValue != "test-header" {
		t.Fatalf("expected headerValue test-header, got %s", awsMethod.headerValue)
	}
	if awsMethod.region != "us-west-2" {
		t.Fatalf("expected region us-west-2, got %s", awsMethod.region)
	}
}

// TestAWSMethod_NewCreds verifies that the NewCreds method returns a non-nil channel
// for receiving new credentials.
func TestAWSMethod_NewCreds(t *testing.T) {
	conf := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/aws",
		Config: map[string]interface{}{
			"type": "ec2",
			"role": "test-role",
		},
	}
	method, err := NewAWSAuthMethod(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	awsMethod, ok := method.(*awsMethod)
	if !ok {
		t.Fatal("expected awsMethod type")
	}

	ch := awsMethod.NewCreds()
	if ch == nil {
		t.Fatal("expected non-nil channel")
	}
}

// TestAWSMethod_Shutdown verifies that the Shutdown method can be called without
// causing a panic.
func TestAWSMethod_Shutdown(t *testing.T) {
	conf := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/aws",
		Config: map[string]interface{}{
			"type": "ec2",
			"role": "test-role",
		},
	}
	method, err := NewAWSAuthMethod(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	awsMethod, ok := method.(*awsMethod)
	if !ok {
		t.Fatal("expected awsMethod type")
	}

	// Should not panic
	awsMethod.Shutdown()
}

// TestAWSMethod_CredSuccess verifies that the CredSuccess method can be called
// without causing a panic.
func TestAWSMethod_CredSuccess(t *testing.T) {
	conf := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/aws",
		Config: map[string]interface{}{
			"type": "ec2",
			"role": "test-role",
		},
	}
	method, err := NewAWSAuthMethod(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	awsMethod, ok := method.(*awsMethod)
	if !ok {
		t.Fatal("expected awsMethod type")
	}

	// Should not panic
	awsMethod.CredSuccess()
}

// TestGenerateLoginDataV2_NilConfig verifies that GenerateLoginDataV2 returns an error
// when called with a nil AWS config.
func TestGenerateLoginDataV2_NilConfig(t *testing.T) {
	ctx := context.Background()
	logger := hclog.NewNullLogger()
	_, err := internal_awsutilv2.GenerateLoginDataV2(ctx, nil, "", "us-east-1", logger)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
	if err.Error() != "aws config must not be nil" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// staticCredentialsProvider provides deterministic credentials for unit tests.
type staticCredentialsProvider struct{}

func (staticCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     "AKIDEXAMPLE",
		SecretAccessKey: "SECRETEXAMPLE",
		SessionToken:    "TOKEN",
		Source:          "unit-test",
	}, nil
}

// TestGenerateLoginDataV2_ValidConfig verifies that GenerateLoginDataV2 returns
// a payload containing the expected login fields when provided static credentials.
func TestGenerateLoginDataV2_ValidConfig(t *testing.T) {
	ctx := context.Background()
	logger := hclog.NewNullLogger()

	// Create a minimal AWS config
	cfg := &aws.Config{
		Region:      "us-east-1",
		Credentials: aws.NewCredentialsCache(staticCredentialsProvider{}),
	}

	loginData, err := internal_awsutilv2.GenerateLoginDataV2(ctx, cfg, "", "us-east-1", logger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range []string{"iam_http_request_method", "iam_request_url", "iam_request_headers", "iam_request_body"} {
		if _, ok := loginData[k]; !ok {
			t.Fatalf("expected %q in login data", k)
		}
	}
}
