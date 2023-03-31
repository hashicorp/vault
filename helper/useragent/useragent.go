// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package useragent

import (
	"fmt"
	"runtime"

	"github.com/hashicorp/vault/version"
)

var (
	// projectURL is the project URL.
	projectURL = "https://www.vaultproject.io/"

	// rt is the runtime - variable for tests.
	rt = runtime.Version()

	// versionFunc is the func that returns the current version. This is a
	// function to take into account the different build processes and distinguish
	// between enterprise and oss builds.
	versionFunc = func() string {
		return version.GetVersion().VersionNumber()
	}
)

// String returns the consistent user-agent string for Vault.
//
// e.g. Vault/0.10.4 (+https://www.vaultproject.io/; go1.10.1)
func String() string {
	return fmt.Sprintf("Vault/%s (+%s; %s)",
		versionFunc(), projectURL, rt)
}

// AgentString returns the consistent user-agent string for Vault Agent.
//
// e.g. Vault Agent/0.10.4 (+https://www.vaultproject.io/; go1.10.1)
func AgentString() string {
	return fmt.Sprintf("Vault Agent/%s (+%s; %s)",
		versionFunc(), projectURL, rt)
}

// AgentTemplatingString returns the consistent user-agent string for Vault Agent Templating.
//
// e.g. Vault Agent Templating/0.10.4 (+https://www.vaultproject.io/; go1.10.1)
func AgentTemplatingString() string {
	return fmt.Sprintf("Vault Agent Templating/%s (+%s; %s)",
		versionFunc(), projectURL, rt)
}

// AgentProxyString returns the consistent user-agent string for Vault Agent API Proxying.
//
// e.g. Vault Agent API Proxy/0.10.4 (+https://www.vaultproject.io/; go1.10.1)
func AgentProxyString() string {
	return fmt.Sprintf("Vault Agent API Proxy/%s (+%s; %s)",
		versionFunc(), projectURL, rt)
}

// AgentProxyStringWithProxiedUserAgent returns the consistent user-agent
// string for Vault Agent API Proxying, keeping the User-Agent of the proxied
// client as an extension to this UserAgent
//
// e.g. Vault Agent API Proxy/0.10.4 (+https://www.vaultproject.io/; go1.10.1); proxiedUserAgent
func AgentProxyStringWithProxiedUserAgent(proxiedUserAgent string) string {
	return fmt.Sprintf("Vault Agent API Proxy/%s (+%s; %s); %s",
		versionFunc(), projectURL, rt, proxiedUserAgent)
}

// AgentAutoAuthString returns the consistent user-agent string for Vault Agent Auto-Auth.
//
// e.g. Vault Agent Auto Auth/0.10.4 (+https://www.vaultproject.io/; go1.10.1)
func AgentAutoAuthString() string {
	return fmt.Sprintf("Vault Agent Auto-Auth/%s (+%s; %s)",
		versionFunc(), projectURL, rt)
}
