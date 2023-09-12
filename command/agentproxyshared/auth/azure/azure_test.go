// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package azure

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
)

// TestAzureAuthMethod tests that NewAzureAuthMethod succeeds
// with valid config.
func TestAzureAuthMethod(t *testing.T) {
	t.Parallel()
	config := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth-test",
		Config: map[string]interface{}{
			"resource":                      "test",
			"client_id":                     "test",
			"role":                          "test",
			"scope":                         "test",
			"authenticate_from_environment": true,
		},
	}

	_, err := NewAzureAuthMethod(config)
	if err != nil {
		t.Fatal(err)
	}
}

// TestAzureAuthMethod_StringAuthFromEnvironment tests that NewAzureAuthMethod succeeds
// with valid config, where authenticate_from_environment is a string literal.
func TestAzureAuthMethod_StringAuthFromEnvironment(t *testing.T) {
	t.Parallel()
	config := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth-test",
		Config: map[string]interface{}{
			"resource":                      "test",
			"client_id":                     "test",
			"role":                          "test",
			"scope":                         "test",
			"authenticate_from_environment": "true",
		},
	}

	_, err := NewAzureAuthMethod(config)
	if err != nil {
		t.Fatal(err)
	}
}

// TestAzureAuthMethod_BadConfig tests that NewAzureAuthMethod fails with
// an invalid config.
func TestAzureAuthMethod_BadConfig(t *testing.T) {
	t.Parallel()
	config := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth-test",
		Config: map[string]interface{}{
			"bad_value": "abc",
		},
	}

	_, err := NewAzureAuthMethod(config)
	if err == nil {
		t.Fatal("Expected error, got none.")
	}
}

// TestAzureAuthMethod_BadAuthFromEnvironment tests that NewAzureAuthMethod fails
// with otherwise valid config, but where authenticate_from_environment is
// an invalid string literal.
func TestAzureAuthMethod_BadAuthFromEnvironment(t *testing.T) {
	t.Parallel()
	config := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth-test",
		Config: map[string]interface{}{
			"resource":                      "test",
			"client_id":                     "test",
			"role":                          "test",
			"scope":                         "test",
			"authenticate_from_environment": "bad_value",
		},
	}

	_, err := NewAzureAuthMethod(config)
	if err == nil {
		t.Fatal("Expected error, got none.")
	}
}
