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

func TestNewGCPAuthMethod_nilConfig(t *testing.T) {
	_, err := NewGCPAuthMethod(nil)
	if err == nil {
		t.Fatal("expected error for nil config, got nil")
	}
}
