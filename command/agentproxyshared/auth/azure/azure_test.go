// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package azure

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
)

func TestAzureAuthMethod(t *testing.T) {
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

func TestAzureAuthMethod_StringAuthFromEnvironment(t *testing.T) {
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

func TestAzureAuthMethod_BadConfig(t *testing.T) {
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

func TestAzureAuthMethod_BadAuthFromEnvironment(t *testing.T) {
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
