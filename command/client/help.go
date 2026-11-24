// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"regexp"
	"strings"

	"github.com/hashicorp/cli"
	hcpvlib "github.com/hashicorp/vault-hcp-lib"
)

// TODO - VAULT-51873: Everything in this file is a temporary workaround that we
// should remove. vault-hcp-lib and the auth plugins (alicloud, cf/pcf, gcp,
// kerberos, oci, oidc, saml) hard-code "vault" in their help and output, so we
// wrap the HCP commands and rewrite their text with regexes, and `auth help`
// does the same for the plugins' help. That's brittle since we're matching text
// we don't own. As each external command or handler takes a BinaryName like our
// in-repo handlers do, we should set it directly and delete its rewrite along
// with its tests. https://hashicorp.atlassian.net/browse/VAULT-51873

var (
	_ cli.Command = (*hcpConnectCommand)(nil)
	_ cli.Command = (*hcpDisconnectCommand)(nil)
)

var (
	// vaultInvocation matches the "vault login" and "vault hcp" invocations
	// that external help hard-codes. It only matches where "vault" starts a
	// word, so "Vault", URLs, paths, VAULT_* variables, and vault-plugin-*
	// names stay as they are.
	vaultInvocation = regexp.MustCompile(`(^|[\s"'` + "`" + `])vault (login|hcp)\b`)

	// quotedVaultInvocation matches the same invocations only when they are
	// quoted, like 'vault login'.
	quotedVaultInvocation = regexp.MustCompile(`(["'` + "`" + `])vault (login|hcp)\b`)
)

// renameVaultInvocations swaps binaryName in for "vault" in the invocations
// that pattern matches. Text for the vault binary is returned as-is.
func renameVaultInvocations(pattern *regexp.Regexp, text, binaryName string) string {
	if binaryName == "" || binaryName == "vault" {
		return text
	}
	return pattern.ReplaceAllStringFunc(text, func(match string) string {
		prefix, command, _ := strings.Cut(match, "vault ")
		return prefix + binaryName + " " + command
	})
}

type hcpConnectCommand struct {
	*hcpvlib.HCPConnectCommand
	BinaryName string
}

func (c *hcpConnectCommand) Help() string {
	return renameVaultInvocations(vaultInvocation, c.HCPConnectCommand.Help(), c.BinaryName)
}

// hcpConnectUi renames quoted vault invocations, like 'vault login', in
// everything hcp connect prints. Unquoted text is left alone because upstream
// prints HCP data such as org and project names through the same calls, and
// users type those names back. Prompts pass through because they only ask for
// a choice and return input.
type hcpConnectUi struct {
	cli.Ui
	BinaryName string
}

func (u *hcpConnectUi) Output(message string) {
	u.Ui.Output(renameVaultInvocations(quotedVaultInvocation, message, u.BinaryName))
}

func (u *hcpConnectUi) Info(message string) {
	u.Ui.Info(renameVaultInvocations(quotedVaultInvocation, message, u.BinaryName))
}

func (u *hcpConnectUi) Warn(message string) {
	u.Ui.Warn(renameVaultInvocations(quotedVaultInvocation, message, u.BinaryName))
}

func (u *hcpConnectUi) Error(message string) {
	u.Ui.Error(renameVaultInvocations(quotedVaultInvocation, message, u.BinaryName))
}

type hcpDisconnectCommand struct {
	*hcpvlib.HCPDisconnectCommand
	BinaryName string
}

func (c *hcpDisconnectCommand) Help() string {
	return renameVaultInvocations(vaultInvocation, c.HCPDisconnectCommand.Help(), c.BinaryName)
}
