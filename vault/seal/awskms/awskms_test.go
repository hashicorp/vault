package awskms

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
)

func TestAWSKMSSeal(t *testing.T) {
	s := NewSeal(logging.NewVaultLogger(log.Trace))
	s.client = &mockAWSKMSSealClient{
		keyID: aws.String(awsTestKeyID),
	}

	_, err := s.SetConfig(nil)
	if err == nil {
		t.Fatal("expected error when AWSKMSSeal key ID is not provided")
	}

	// Set the key
	os.Setenv(EnvAWSKMSSealKeyID, awsTestKeyID)
	_, err = s.SetConfig(nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAWSKMSSeal_Lifecycle(t *testing.T) {
	s := NewSeal(logging.NewVaultLogger(log.Trace))
	s.client = &mockAWSKMSSealClient{
		keyID: aws.String(awsTestKeyID),
	}

	os.Setenv(EnvAWSKMSSealKeyID, awsTestKeyID)
	_, err := s.SetConfig(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Test Encrypt and Decrypt calls
	input := []byte("foo")
	swi, err := s.Encrypt(context.Background(), input)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	pt, err := s.Decrypt(context.Background(), swi)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	if !reflect.DeepEqual(input, pt) {
		t.Fatalf("expected %s, got %s", input, pt)
	}
}

func TestAWSKMSSeal_custom_endpoint(t *testing.T) {
	customEndpoint := "https://custom.endpoint"
	customEndpoint2 := "https://custom.endpoint.2"
	endpointENV := "AWS_KMS_ENDPOINT"

	// unset at end of test
	os.Setenv(EnvAWSKMSSealKeyID, awsTestKeyID)
	defer func() {
		if err := os.Unsetenv(EnvAWSKMSSealKeyID); err != nil {
			t.Fatal(err)
		}
	}()

	cfg := make(map[string]string)
	cfg["endpoint"] = customEndpoint

	testCases := []struct {
		Title    string
		Env      string
		Config   map[string]string
		Expected *string
	}{
		{
			// Default will have nil for the config endpoint, and be looked up
			// dynamically by the SDK
			Title: "Default",
		},
		{
			Title:    "Environment",
			Env:      customEndpoint,
			Expected: aws.String(customEndpoint),
		},
		{
			Title:    "Config",
			Config:   cfg,
			Expected: aws.String(customEndpoint),
		},
		{
			// Expect environment to take precedence over configuration
			Title:    "Env-Config",
			Env:      customEndpoint2,
			Config:   cfg,
			Expected: aws.String(customEndpoint2),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			s := NewSeal(logging.NewVaultLogger(log.Trace))

			s.client = &mockAWSKMSSealClient{
				keyID: aws.String(awsTestKeyID),
			}

			if tc.Env != "" {
				if err := os.Setenv(endpointENV, tc.Env); err != nil {
					t.Fatal(err)
				}
			}

			// cfg starts as nil, and takes a test case value if given. If not,
			// SetConfig is called with nil and creates it's own config
			var cfg map[string]string
			if tc.Config != nil {
				cfg = tc.Config
			}
			if _, err := s.SetConfig(cfg); err != nil {
				t.Fatalf("error setting config: %s", err)
			}

			// call getAWSKMSClient() to get the configured client and verify it's
			// endpoint
			k, err := s.getAWSKMSClient()
			if err != nil {
				t.Fatal(err)
			}

			if tc.Expected == nil && k.Config.Endpoint != nil {
				t.Fatalf("Expected nil endpoint, got: (%s)", *k.Config.Endpoint)
			}

			if tc.Expected != nil {
				if k.Config.Endpoint == nil {
					t.Fatal("expected custom endpoint, but config was nil")
				}
				if *k.Config.Endpoint != *tc.Expected {
					t.Fatalf("expected custom endpoint (%s), got: (%s)", *tc.Expected, *k.Config.Endpoint)
				}
			}

			// clear endpoint env after each test
			if err := os.Unsetenv(endpointENV); err != nil {
				t.Fatal(err)
			}
		})
	}

}
