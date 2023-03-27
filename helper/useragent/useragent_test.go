// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package useragent

import (
	"testing"
)

func TestUserAgent(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	act := String()

	exp := "Vault/1.2.3 (+https://vault-test.com; go5.0)"
	if exp != act {
		t.Errorf("expected %q to be %q", act, exp)
	}
}

// TestUserAgentVaultAgent tests the AgentString() function works
// as expected
func TestUserAgentVaultAgent(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	act := AgentString()

	exp := "Vault Agent/1.2.3 (+https://vault-test.com; go5.0)"
	if exp != act {
		t.Errorf("expected %q to be %q", act, exp)
	}
}

// TestUserAgentVaultAgentTemplating tests the AgentTemplatingString() function works
// as expected
func TestUserAgentVaultAgentTemplating(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	act := AgentTemplatingString()

	exp := "Vault Agent Templating/1.2.3 (+https://vault-test.com; go5.0)"
	if exp != act {
		t.Errorf("expected %q to be %q", act, exp)
	}
}

// TestUserAgentVaultAgentProxy tests the AgentProxyString() function works
// as expected
func TestUserAgentVaultAgentProxy(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	act := AgentProxyString()

	exp := "Vault Agent API Proxy/1.2.3 (+https://vault-test.com; go5.0)"
	if exp != act {
		t.Errorf("expected %q to be %q", act, exp)
	}
}

// TestUserAgentVaultAgentProxyWithProxiedUserAgent tests the AgentProxyStringWithProxiedUserAgent()
// function works as expected
func TestUserAgentVaultAgentProxyWithProxiedUserAgent(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }
	userAgent := "my-user-agent"

	act := AgentProxyStringWithProxiedUserAgent(userAgent)

	exp := "Vault Agent API Proxy/1.2.3 (+https://vault-test.com; go5.0); my-user-agent"
	if exp != act {
		t.Errorf("expected %q to be %q", act, exp)
	}
}

// TestUserAgentVaultAgentAutoAuth tests the AgentAutoAuthString() function works
// as expected
func TestUserAgentVaultAgentAutoAuth(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	act := AgentAutoAuthString()

	exp := "Vault Agent Auto-Auth/1.2.3 (+https://vault-test.com; go5.0)"
	if exp != act {
		t.Errorf("expected %q to be %q", act, exp)
	}
}
