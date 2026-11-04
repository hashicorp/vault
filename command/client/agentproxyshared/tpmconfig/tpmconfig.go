// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

// Package tpmconfig holds the tpm_cert configuration types and the parsing,
// validation, and auto-auth resolution logic shared by Vault Agent and Vault
// Proxy.
package tpmconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// AutoAuthMethodType is the auto_auth method type that uses tpm_cert blocks.
	AutoAuthMethodType = "tpm"

	defaultPercentRenewBeforeExpiry = 20
	defaultUmask                    = "077"

	blockName = "tpm_cert"
	blockFmt  = "tpm_cert[%d]"

	// umaskBaseOctal and umaskBitSize are used when parsing umask values.
	umaskBaseOctal = 8
	umaskBitSize   = 32
)

// TPMCert is a parsed tpm_cert block.
type TPMCert struct {
	Name string `hcl:"-"`

	Namespace string      `hcl:"namespace,optional"`
	MountPath string      `hcl:"mount_path"`
	RoleName  string      `hcl:"role_name"`
	CSR       *TPMCertCSR `hcl:"csr,optional"`

	TPM         *TPMCertTPM         `hcl:"tpm,optional"`
	Renew       *TPMCertRenew       `hcl:"renew,optional"`
	Destination *TPMCertDestination `hcl:"destination"`

	LineNumber int
}

type TPMCertCSR struct {
	CommonName string `hcl:"CN,optional"`
	Country    string `hcl:"C,optional"`
	State      string `hcl:"ST,optional"`
	Locality   string `hcl:"L,optional"`
	Org        string `hcl:"O,optional"`
	OrgUnit    string `hcl:"OU,optional"`
}

type TPMCertTPM struct {
	DevicePath string `hcl:"device_path,optional"`
}

type TPMCertRenew struct {
	PercentRenewBeforeExpiry *int   `hcl:"percent_renew_before_expiry,optional"`
	MinBackoff               string `hcl:"min_backoff,optional"`
	MaxBackoff               string `hcl:"max_backoff,optional"`
}

type TPMCertDestination struct {
	Path  string `hcl:"path"`
	Umask string `hcl:"umask,optional"`
}

// Parse decodes every tpm_cert block found in list.
func Parse(list *ast.ObjectList) ([]*TPMCert, error) {
	tpmCertList := list.Filter(blockName)
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

// seenBlockName tracks the array index and source line of a previously seen
// block name so duplicate-name errors can reference both locations.
type seenBlockName struct {
	idx  int
	line int
}

// Validate checks the tpm_cert entries and applies defaults to optional
// fields. inheritedNamespace is used for entries that don't set a namespace.
func Validate(entries []*TPMCert, inheritedNamespace string) error {
	seenDestinationPaths := make(map[string]int)
	seenNames := make(map[string]seenBlockName)

	for idx, entry := range entries {
		if entry == nil {
			return fmt.Errorf("tpm_cert[%d]: configuration is empty", idx)
		}

		blockName := fmt.Sprintf(blockFmt, idx)
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

		if err := validateRenew(blockName, entry); err != nil {
			return err
		}
		if err := validateDestination(blockName, idx, seenDestinationPaths, entry); err != nil {
			return err
		}
	}

	return nil
}

func validateRenew(blockName string, entry *TPMCert) error {
	if entry.Renew == nil {
		entry.Renew = &TPMCertRenew{}
	}

	if entry.Renew.PercentRenewBeforeExpiry == nil {
		defaultPercent := defaultPercentRenewBeforeExpiry
		entry.Renew.PercentRenewBeforeExpiry = &defaultPercent
	} else if *entry.Renew.PercentRenewBeforeExpiry <= 0 || *entry.Renew.PercentRenewBeforeExpiry >= 100 {
		return fmt.Errorf("%s.renew.percent_renew_before_expiry must be between 1 and 99", blockName)
	}

	return nil
}

func validateDestination(blockName string, idx int, seenDestinationPaths map[string]int, entry *TPMCert) error {
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
		entry.Destination.Umask = defaultUmask
	} else {
		if _, err := strconv.ParseUint(entry.Destination.Umask, umaskBaseOctal, umaskBitSize); err != nil {
			return fmt.Errorf("%s.destination.umask invalid value %q: %w", blockName, entry.Destination.Umask, err)
		}
	}

	return nil
}

// normalizeDestinationPath converts path to an absolute, cleaned path and
// resolves symlinks when possible. If Abs or EvalSymlinks fails, it falls back
// to returning the best cleaned path available.
func normalizeDestinationPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	absPath = filepath.Clean(absPath)
	if resolved, err := filepath.EvalSymlinks(absPath); err == nil {
		return resolved
	}
	return absPath
}

// AutoAuthReference returns the tpm_cert block referenced by an auto_auth
// method. It returns logical.ErrNotFound if the method isn't a TPM method.
func AutoAuthReference(methodType string, methodConfig map[string]interface{}, certs []*TPMCert) (*TPMCert, error) {
	if methodType != AutoAuthMethodType {
		return nil, logical.ErrNotFound
	}

	tpmCertName, hasTpmCert := methodConfig["tpm_cert"].(string)
	if !hasTpmCert {
		return nil, fmt.Errorf("missing tpm_cert config")
	}

	for _, tc := range certs {
		if tc.Name == tpmCertName {
			return tc, nil
		}
	}

	return nil, fmt.Errorf("auto_auth.method.config references tpm_cert %q but no such tpm_cert block exists", tpmCertName)
}

// ValidateAutoAuthReference validates that a TPM auto_auth method references
// an existing tpm_cert block. Non-TPM methods are ignored.
func ValidateAutoAuthReference(methodType string, methodConfig map[string]interface{}, certs []*TPMCert) error {
	_, err := AutoAuthReference(methodType, methodConfig, certs)
	if errors.Is(err, logical.ErrNotFound) {
		return nil
	}
	return err
}

// ConfigureAutoAuth updates a TPM auto_auth method's config map to include
// information derived from the referenced tpm_cert block: tpm_cert_destination,
// tpm_cert_device, and role_name. It is a no-op for non-TPM methods.
func ConfigureAutoAuth(methodType string, methodConfig map[string]interface{}, certs []*TPMCert) error {
	tpmCert, err := AutoAuthReference(methodType, methodConfig, certs)
	if err != nil {
		if errors.Is(err, logical.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("error resolving auto_auth tpm cert reference: %w", err)
	}

	if tpmCert.Destination == nil {
		return fmt.Errorf("tpm_cert %q: destination block must be specified", tpmCert.Name)
	}

	methodConfig["tpm_cert_destination"] = tpmCert.Destination.Path
	if tpmCert.TPM != nil {
		methodConfig["tpm_cert_device"] = tpmCert.TPM.DevicePath
	}
	// Since most auto_auth methods allow configuring the role, we'll tolerate
	// it here so long as it matches the value in tpm_cert.  We can't check this
	// in validation because it would already have been overwritten.
	existingRoleName, ok := methodConfig["role_name"].(string)
	if ok && existingRoleName != tpmCert.RoleName {
		return fmt.Errorf("role_name in autoauth config stanza (%q) conflicts with role_name in tpm_cert %q: %q",
			methodConfig["role_name"], tpmCert.Name, tpmCert.RoleName)
	}
	methodConfig["role_name"] = tpmCert.RoleName

	return nil
}
