// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package useragent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserAgent(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	act := String()

	exp := "Vault/1.2.3 (+https://vault-test.com; go5.0)"
	require.Equal(t, exp, act)
}

// TestUserAgent_VaultAgent tests the AgentString() function works
// as expected
func TestUserAgent_VaultAgent(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	act := AgentString()

	exp := "Vault Agent/1.2.3 (+https://vault-test.com; go5.0)"
	require.Equal(t, exp, act)
}

// TestUserAgent_VaultAgentTemplating tests the AgentTemplatingString() function works
// as expected
func TestUserAgent_VaultAgentTemplating(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	act := AgentTemplatingString()

	exp := "Vault Agent Templating/1.2.3 (+https://vault-test.com; go5.0)"
	require.Equal(t, exp, act)
}

// TestUserAgent_VaultAgentProxy tests the AgentProxyString() function works
// as expected
func TestUserAgent_VaultAgentProxy(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	act := AgentProxyString()

	exp := "Vault Agent API Proxy/1.2.3 (+https://vault-test.com; go5.0)"
	require.Equal(t, exp, act)
}

// TestUserAgent_VaultAgentProxyWithProxiedUserAgent tests the AgentProxyStringWithProxiedUserAgent()
// function works as expected
func TestUserAgent_VaultAgentProxyWithProxiedUserAgent(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }
	userAgent := "my-user-agent"

	act := AgentProxyStringWithProxiedUserAgent(userAgent)

	exp := "Vault Agent API Proxy/1.2.3 (+https://vault-test.com; go5.0); my-user-agent"
	require.Equal(t, exp, act)
}

// TestUserAgent_VaultAgentAutoAuth tests the AgentAutoAuthString() function works
// as expected
func TestUserAgent_VaultAgentAutoAuth(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	act := AgentAutoAuthString()

	exp := "Vault Agent Auto-Auth/1.2.3 (+https://vault-test.com; go5.0)"
	require.Equal(t, exp, act)
}
