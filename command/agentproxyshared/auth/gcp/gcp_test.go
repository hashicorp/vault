// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package gcp

import (
	"testing"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
)

func baseConfig(extra map[string]interface{}) *auth.AuthConfig {
	cfg := map[string]interface{}{
		"type": "iam",
		"role": "my-role",
	}
	for k, v := range extra {
		cfg[k] = v
	}
	return &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/gcp",
		Config:    cfg,
	}
}

// TestNewGCPAuthMethod_serviceAccountDefault verifies that when no
// service_account is provided the field is left empty, so that the IAM
// branch can derive it from credentials.ClientEmail instead of sending
// the literal string "default" to the GCP IAM API.
func TestNewGCPAuthMethod_serviceAccountDefault(t *testing.T) {
	m, err := NewGCPAuthMethod(baseConfig(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := m.(*gcpMethod)
	if g.serviceAccount != "" {
		t.Errorf("expected empty serviceAccount when not configured, got %q", g.serviceAccount)
	}
}

// TestNewGCPAuthMethod_explicitServiceAccount verifies that an explicitly
// configured service_account is preserved as-is.
func TestNewGCPAuthMethod_explicitServiceAccount(t *testing.T) {
	const want = "my-sa@my-project.iam.gserviceaccount.com"
	m, err := NewGCPAuthMethod(baseConfig(map[string]interface{}{
		"service_account": want,
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := m.(*gcpMethod)
	if g.serviceAccount != want {
		t.Errorf("expected serviceAccount %q, got %q", want, g.serviceAccount)
	}
}

// TestNewGCPAuthMethod_GCETypeWithNoServiceAccount verifies that a GCE-type
// method can be constructed without a service_account; the "default" alias
// is applied inside Authenticate only for the GCE flow.
func TestNewGCPAuthMethod_GCETypeWithNoServiceAccount(t *testing.T) {
	m, err := NewGCPAuthMethod(&auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/gcp",
		Config:    map[string]interface{}{"type": "gce", "role": "my-role"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := m.(*gcpMethod)
	if g.serviceAccount != "" {
		t.Errorf("expected empty serviceAccount for GCE without explicit config, got %q", g.serviceAccount)
	}
}

// TestNewGCPAuthMethod_missingType verifies that missing type returns an error.
func TestNewGCPAuthMethod_missingType(t *testing.T) {
	_, err := NewGCPAuthMethod(&auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/gcp",
		Config:    map[string]interface{}{"role": "my-role"},
	})
	if err == nil {
		t.Fatal("expected error for missing 'type', got nil")
	}
}

// TestNewGCPAuthMethod_missingRole verifies that missing role returns an error.
func TestNewGCPAuthMethod_missingRole(t *testing.T) {
	_, err := NewGCPAuthMethod(&auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/gcp",
		Config:    map[string]interface{}{"type": "iam"},
	})
	if err == nil {
		t.Fatal("expected error for missing 'role', got nil")
	}
}

// TestNewGCPAuthMethod_nilConfig verifies that a nil config returns an error.
func TestNewGCPAuthMethod_nilConfig(t *testing.T) {
	_, err := NewGCPAuthMethod(nil)
	if err == nil {
		t.Fatal("expected error for nil config, got nil")
	}
}
