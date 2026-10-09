// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	sdktpm "github.com/hashicorp/vault/sdk/helper/tpm"
)

var _ cli.Command = (*TPMInspectCommand)(nil)

type TPMInspectCommand struct {
	*base.BaseCommand
}

func (c *TPMInspectCommand) Synopsis() string {
	return "Reports the current status of attestation"
}

func (c *TPMInspectCommand) Help() string {
	helpText := `
Usage: vault tpm inspect [options]

  Reads the attestation state directory written by "vault tpm attest" and
  displays information about the stored certificate and key material without
  contacting Vault or opening the TPM device:

      client.crt       - the issued certificate (PEM)
      ca_chain.pem     - the CA chain (PEM)
      client-key.json  - key reference metadata (JSON)
      ak.blob          - attestation key blob
      app.blob         - application key blob

  Inspect the attestation state in the current directory:

      $ vault tpm inspect

  Inspect a state directory at a custom path:

      $ vault tpm inspect -tpm-state-dir=/etc/vault/tpm
	` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *TPMInspectCommand) Flags() *base.FlagSets {
	return c.FlagSetForBit(base.FlagSetHTTP | base.FlagSetOutputFormat)
}

// Run executes main inspect logic. The function reads the state directory, decodes the cert PEM, parses the certificate, checks time until expiration,
// builds the pipe deliminated output, checks for existing role name, and outputs the columnized output
func (c *TPMInspectCommand) Run(args []string) int {
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	result, err := sdktpm.LoadFromDirectory(c.FlagTPMStateDir)
	if err != nil {
		c.UI.Error(fmt.Sprintf("error reading state directory %q: %s", c.FlagTPMStateDir, err))
		return 2
	}

	block, _ := pem.Decode(result.CertPEM)
	if block == nil {
		c.UI.Error("error parsing certificate: no PEM block found")
		return 2
	}
	serverCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		c.UI.Error(fmt.Sprintf("error parsing certificate: %s", err))
		return 2
	}

	remaining := max(time.Until(result.NotAfter).Round(time.Second), 0)

	data := map[string]any{
		"common_name":       serverCert.Subject.CommonName,
		"public_key_sha256": result.KeyMaterial.PublicKeySHA256,
		"valid_from":        result.NotBefore.UTC().Format(time.RFC1123),
		"valid_until":       result.NotAfter.UTC().Format(time.RFC1123),
		"time_until_expiry": remaining.String(),
	}

	if len(serverCert.Subject.Organization) > 0 {
		data["organization"] = serverCert.Subject.Organization
	}
	if len(serverCert.Subject.OrganizationalUnit) > 0 {
		data["org_unit"] = serverCert.Subject.OrganizationalUnit
	}
	if len(serverCert.Subject.Country) > 0 {
		data["country"] = serverCert.Subject.Country
	}
	if len(serverCert.Subject.Province) > 0 {
		data["state"] = serverCert.Subject.Province
	}
	if len(serverCert.Subject.Locality) > 0 {
		data["locality"] = serverCert.Subject.Locality
	}

	othersans, err := certutil.GetOtherSANsFromX509Extensions(serverCert.Extensions)
	if err != nil {
		c.UI.Error(fmt.Sprintf("error reading certificate SANs: %s", err))
		return 2
	}

	var tpmID string
	for _, other := range othersans {
		if other.Oid == consts.TPMAuthOIDTPMID {
			tpmID = other.Value
			break
		}
	}
	if tpmID == "" {
		c.UI.Error("certificate is missing TPM ID SAN")
		return 2
	}
	data["tpm_id"] = tpmID

	var roleName string
	for _, other := range othersans {
		if other.Oid == consts.TPMAuthOIDRoleName {
			roleName = other.Value
			break
		}
	}
	if roleName == "" {
		c.UI.Error("certificate is missing role name SAN")
		return 2
	}

	data["role"] = roleName

	return base.OutputData(c.UI, data)
}
