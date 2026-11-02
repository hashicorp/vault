// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	defaultTPMCertPercentRenewBeforeExpiry = 20
	defaultTPMCertUmask                    = "077"

	tpmCertBlockName = "tpm_cert"
	tpmCertBlockFmt  = "tpm_cert[%d]"
)

// validateAutoAuthTPMCertReferences validates that tpm_cert references in auto_auth
// method config are valid and don't conflict with explicit cert paths.
func (c *Config) validateAutoAuthTPMCertReferences() error {
	_, err := c.getAutoAuthTPMCertReference()
	if errors.Is(err, logical.ErrNotFound) {
		return nil
	}
	return err
}

// configureTPMAutoAuth updates the TPM autoauth config stanza to include information
// derived from the tpm_cert stanza.  It is an error to use TPM autoauth without
// referencing a tpm_cert stanza, though that is handled in a validation routine
// rather than here.
func (c *Config) configureTPMAutoAuth() error {
	tpmCert, err := c.getAutoAuthTPMCertReference()
	if err != nil {
		if errors.Is(err, logical.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("error resolving auto_auth tpm cert reference: %w", err)
	}
	if tpmCert != nil {
		c.AutoAuth.Method.Config["tpm_cert_destination"] = tpmCert.Destination.Path
		if tpmCert.TPM != nil {
			c.AutoAuth.Method.Config["tpm_cert_device"] = tpmCert.TPM.DevicePath
		}
		// Since most auto_auth methods allow configuring the role, we'll tolerate
		// it here so long as it matches the value in tpm_cert.  We can't check this
		// in validation because it would already have been overwritten.
		existingRoleName, ok := c.AutoAuth.Method.Config["role_name"].(string)
		if ok && existingRoleName != tpmCert.RoleName {
			return fmt.Errorf("role_name in autoauth config stanza (%q) conflicts with role_name in tpm_cert %q: %q",
				c.AutoAuth.Method.Config["role_name"], tpmCert.Name, tpmCert.RoleName)
		}
		c.AutoAuth.Method.Config["role_name"] = tpmCert.RoleName
	}

	return nil
}

func (c *Config) getAutoAuthTPMCertReference() (*TPMCert, error) {
	if c.AutoAuth == nil || c.AutoAuth.Method == nil {
		return nil, logical.ErrNotFound
	}

	method := c.AutoAuth.Method
	if method.Type != "tpm" {
		return nil, logical.ErrNotFound
	}

	// Check if tpm_cert is specified in the config
	tpmCertName, hasTpmCert := method.Config["tpm_cert"].(string)
	if !hasTpmCert {
		return nil, fmt.Errorf("missing tpm_cert config")
	}

	// Find the referenced tpm_cert block
	var tpmCert *TPMCert
	for _, tc := range c.TPMCerts {
		if tc.Name == tpmCertName {
			tpmCert = tc
			break
		}
	}

	if tpmCert == nil {
		return nil, fmt.Errorf("auto_auth.method.config references tpm_cert %q but no such tpm_cert block exists", tpmCertName)
	}

	return tpmCert, nil
}

func (c *Config) validateTPMCertConfig(logger hclog.Logger) error {
	inheritedNamespace := ""
	if c.Vault != nil {
		inheritedNamespace = c.Vault.Namespace
	}
	return ValidateTPMCert(c.TPMCerts, inheritedNamespace, logger)
}

func ValidateTPMCert(entries []*TPMCert, inheritedNamespace string, _ hclog.Logger) error {
	seenDestinationPaths := make(map[string]int)
	seenNames := make(map[string]seenBlockName)

	for idx, entry := range entries {
		if entry == nil {
			return fmt.Errorf("tpm_cert[%d]: configuration is empty", idx)
		}

		blockName := fmt.Sprintf(tpmCertBlockFmt, idx)
		if entry.Name != "" {
			if prev, exists := seenNames[entry.Name]; exists {
				if prev.line > 0 && entry.LineNumber > 0 {
					return fmt.Errorf("%s (line %d) name %q duplicates tpm_cert[%d] (line %d)", blockName, entry.LineNumber, entry.Name, prev.idx, prev.line)
				}
				return fmt.Errorf("%s name %q duplicates tpm_cert[%d]", blockName, entry.Name, prev.idx)
			}
			seenNames[entry.Name] = seenBlockName{idx: idx, line: entry.LineNumber}
		}

		ns := namespace.Canonicalize(entry.Namespace)
		if ns == "" {
			ns = namespace.Canonicalize(inheritedNamespace)
		}
		entry.Namespace = ns

		if strings.TrimSpace(entry.MountPath) == "" {
			return fmt.Errorf("%s.mount_path must be specified", blockName)
		}

		if err := validateTPMCertRenew(blockName, entry); err != nil {
			return err
		}
		if err := validateTPMCertDestination(blockName, idx, seenDestinationPaths, entry); err != nil {
			return err
		}
	}

	return nil
}

func validateTPMCertRenew(blockName string, entry *TPMCert) error {
	if entry.Renew == nil {
		entry.Renew = &TPMCertRenew{}
	}

	if entry.Renew.PercentRenewBeforeExpiry == nil {
		defaultPercent := defaultTPMCertPercentRenewBeforeExpiry
		entry.Renew.PercentRenewBeforeExpiry = &defaultPercent
	} else if *entry.Renew.PercentRenewBeforeExpiry <= 0 || *entry.Renew.PercentRenewBeforeExpiry >= 100 {
		return fmt.Errorf("%s.renew.percent_renew_before_expiry must be between 1 and 99", blockName)
	}

	return nil
}

func validateTPMCertDestination(blockName string, idx int, seenDestinationPaths map[string]int, entry *TPMCert) error {
	if entry.Destination == nil {
		return fmt.Errorf("%s.destination block must be specified", blockName)
	}

	path := strings.TrimSpace(entry.Destination.Path)
	if path == "" {
		return fmt.Errorf("%s.destination.path must be specified", blockName)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s.destination.path must be an existing directory: %w", blockName, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s.destination.path must be an existing directory", blockName)
	}

	normalizedPath := normalizeDestinationPath(path)
	if prevIdx, exists := seenDestinationPaths[normalizedPath]; exists {
		return fmt.Errorf("%s.destination.path duplicates tpm_cert[%d].destination.path", blockName, prevIdx)
	}
	seenDestinationPaths[normalizedPath] = idx

	if strings.TrimSpace(entry.Destination.Umask) == "" {
		entry.Destination.Umask = defaultTPMCertUmask
	} else {
		if _, err := strconv.ParseUint(entry.Destination.Umask, umaskBaseOctal, umaskBitSize); err != nil {
			return fmt.Errorf("%s.destination.umask invalid value %q: %w", blockName, entry.Destination.Umask, err)
		}
	}

	return nil
}

func ParseTPMCert(list *ast.ObjectList) ([]*TPMCert, error) {
	tpmCertList := list.Filter(tpmCertBlockName)
	if len(tpmCertList.Items) == 0 {
		return []*TPMCert{}, nil
	}

	entries := make([]*TPMCert, 0, len(tpmCertList.Items))
	for _, item := range tpmCertList.Items {
		var entry TPMCert
		if err := hcl.DecodeObject(&entry, item.Val); err != nil {
			return nil, err
		}
		if len(item.Keys) != 1 {
			return nil, fmt.Errorf("tpm_cert: expected exactly one label, got %d", len(item.Keys))
		}
		entry.Name = strings.Trim(item.Keys[0].Token.Text, `"`)
		entry.LineNumber = item.Pos().Line
		entries = append(entries, &entry)
	}

	return entries, nil
}

func parseTPMCert(result *Config, list *ast.ObjectList) error {
	entries, err := ParseTPMCert(list)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return nil
	}

	result.TPMCerts = entries
	return nil
}
