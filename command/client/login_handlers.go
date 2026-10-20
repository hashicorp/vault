// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	credAliCloud "github.com/hashicorp/vault-plugin-auth-alicloud"
	credCF "github.com/hashicorp/vault-plugin-auth-cf"
	credGcp "github.com/hashicorp/vault-plugin-auth-gcp/plugin"
	credOIDC "github.com/hashicorp/vault-plugin-auth-jwt"
	credKerb "github.com/hashicorp/vault-plugin-auth-kerberos"
	credOCI "github.com/hashicorp/vault-plugin-auth-oci"
	base "github.com/hashicorp/vault/command/base"
	"github.com/hashicorp/vault/command/client/login/aws"
	"github.com/hashicorp/vault/command/client/login/cert"
	"github.com/hashicorp/vault/command/client/login/github"
	"github.com/hashicorp/vault/command/client/login/ldap"
	"github.com/hashicorp/vault/command/client/login/okta"
	"github.com/hashicorp/vault/command/client/login/spiffe"
	"github.com/hashicorp/vault/command/client/login/token"
	"github.com/hashicorp/vault/command/client/login/tpm"
	"github.com/hashicorp/vault/command/client/login/userpass"
)

// LoginHandlers returns the CLI handlers used by `login` and `auth help`,
// keyed by auth method type. Both vault and vault-client use this set.
func LoginHandlers() map[string]base.LoginHandler {
	handlers := map[string]base.LoginHandler{
		"alicloud": &credAliCloud.CLIHandler{},
		"aws":      &aws.CLIHandler{},
		"cert":     &cert.CLIHandler{},
		"cf":       &credCF.CLIHandler{},
		"gcp":      &credGcp.CLIHandler{},
		"github":   &github.CLIHandler{},
		"kerberos": &credKerb.CLIHandler{},
		"ldap":     &ldap.CLIHandler{},
		"oci":      &credOCI.CLIHandler{},
		"oidc":     &credOIDC.CLIHandler{},
		"okta":     &okta.CLIHandler{},
		"pcf":      &credCF.CLIHandler{}, // Deprecated.
		"radius": &userpass.CLIHandler{
			DefaultMount: "radius",
		},
		"spiffe": &spiffe.CLIHandler{},
		"token":  &token.CLIHandler{},
		"tpm":    &tpm.CLIHandler{},
		"userpass": &userpass.CLIHandler{
			DefaultMount: "userpass",
		},
	}
	entExtendLoginHandlers(handlers)

	return handlers
}
