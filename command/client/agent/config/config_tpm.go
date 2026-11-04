// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/command/client/agentproxyshared/tpmconfig"
)

// validateAutoAuthTPMCertReferences validates that tpm_cert references in auto_auth
// method config are valid and don't conflict with explicit cert paths.
func (c *Config) validateAutoAuthTPMCertReferences() error {
	if c.AutoAuth == nil || c.AutoAuth.Method == nil {
		return nil
	}
	return tpmconfig.ValidateAutoAuthReference(c.AutoAuth.Method.Type, c.AutoAuth.Method.Config, c.TPMCerts)
}

// configureTPMAutoAuth updates the TPM autoauth config stanza to include information
// derived from the tpm_cert stanza.  It is an error to use TPM autoauth without
// referencing a tpm_cert stanza, though that is handled in a validation routine
// rather than here.
func (c *Config) configureTPMAutoAuth() error {
	if c.AutoAuth == nil || c.AutoAuth.Method == nil {
		return nil
	}
	return tpmconfig.ConfigureAutoAuth(c.AutoAuth.Method.Type, c.AutoAuth.Method.Config, c.TPMCerts)
}

func (c *Config) validateTPMCertConfig() error {
	inheritedNamespace := ""
	if c.Vault != nil {
		inheritedNamespace = c.Vault.Namespace
	}
	return tpmconfig.Validate(c.TPMCerts, inheritedNamespace)
}

func parseTPMCert(result *Config, list *ast.ObjectList) error {
	entries, err := tpmconfig.Parse(list)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return nil
	}

	result.TPMCerts = entries
	return nil
}
