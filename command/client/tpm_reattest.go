// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	sdktpm "github.com/hashicorp/vault/sdk/helper/tpm"
	"github.com/posener/complete"
)

var _ cli.Command = (*TPMReattestCommand)(nil)

type TPMReattestCommand struct {
	*base.BaseCommand

	flagMountPath string
	flagRoleName  string
}

func (c *TPMReattestCommand) Synopsis() string {
	return "Renew a TPM certificate using subject fields from the existing certificate"
}

func (c *TPMReattestCommand) Help() string {
	helpText := `
Usage: vault tpm reattest [options]

  Performs a fresh TPM attestation, reusing the certificate subject fields
  stored in the existing certificate in the state directory. This is useful
  when a certificate has expired or is approaching expiry and needs to be
  renewed without re-specifying all subject flags.

  The existing client.crt in the state directory is read to recover the
  subject (CN, O, OU, C, ST, L).

  Reattest using the existing certificate subject (role read from cert SAN):

      $ vault tpm reattest

  Reattest with an explicit role name and custom state directory:

      $ vault tpm reattest -role-name=my-role \
          -tpm-state-dir=/etc/vault/tpm

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

// Flags registers each flag with their usage instructions, default values, and tab completion behavior
func (c *TPMReattestCommand) Flags() *base.FlagSets {
	set := c.FlagSetForBit(base.FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&base.StringVar{
		Name:       "mount-path",
		Target:     &c.flagMountPath,
		Default:    "auth/tpm",
		Completion: complete.PredictAnything,
		Usage:      "Mount path for the TPM auth method.",
	})

	f.StringVar(&base.StringVar{
		Name:    "role-name",
		Target:  &c.flagRoleName,
		Default: "",
		Usage:   "Name of the Vault TPM auth role to attest against. Optional if the existing certificate encodes the role name as a SAN.",
	})

	return set
}

// Run reads existing certificate, checks for existing role name, creates HTTP client, builds attestation configuration,
// runs both phases of attestation, and saves result to state directory on success.
func (c *TPMReattestCommand) Run(args []string) int {
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()

	if err := os.MkdirAll(c.FlagTPMStateDir, 0o700); err != nil {
		c.UI.Error(fmt.Sprintf("error creating state directory %q: %s", c.FlagTPMStateDir, err))
		return 2
	}

	// Read the existing certificate to recover the subject fields and,
	// if -role-name was not provided, the role name from the cert SAN.
	csrParams, err := c.csrParamsFromExistingCert()
	if err != nil {
		c.UI.Error(fmt.Sprintf("error reading existing certificate: %s", err))
		return 2
	}

	if strings.TrimSpace(c.flagRoleName) == "" {
		roleName, err := c.roleNameFromExistingCert()
		if err != nil || strings.TrimSpace(roleName) == "" {
			c.UI.Error("flag -role-name is required (the existing certificate does not encode a role name SAN)")
			return 1
		}
		c.flagRoleName = roleName
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("error creating client: %s", err))
		return 2
	}

	cfg := &sdktpm.AttestationConfig{
		MountPath: c.flagMountPath,
		RoleName:  c.flagRoleName,
		TPMConfig: sdktpm.TPMConfig{
			DevicePath: c.FlagTPMDevicePath,
		},
		CSRParams: csrParams,
	}

	result, err := sdktpm.Attest(context.Background(), client, cfg)
	if err != nil {
		c.UI.Error(fmt.Sprintf("attestation failed: %s", err))
		return 2
	}

	if err := result.SaveToDirectory(c.FlagTPMStateDir, 0); err != nil {
		c.UI.Error(fmt.Sprintf("error writing attestation output: %s", err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Reattestation successful. Certificate valid until: %s",
		result.NotAfter.UTC().Format(time.RFC1123)))
	return 0
}

// csrParamsFromExistingCert reads client.crt from the state directory and
// returns a CSRParams populated from its subject.
func (c *TPMReattestCommand) csrParamsFromExistingCert() (*sdktpm.CSRParams, error) {
	certPath := filepath.Join(c.FlagTPMStateDir, "client.crt")
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", certPath, err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode certificate PEM from %s", certPath)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	params := &sdktpm.CSRParams{
		CommonName: cert.Subject.CommonName,
	}
	if len(cert.Subject.Organization) > 0 {
		params.Org = cert.Subject.Organization[0]
	}
	if len(cert.Subject.OrganizationalUnit) > 0 {
		params.OrgUnit = cert.Subject.OrganizationalUnit[0]
	}
	if len(cert.Subject.Country) > 0 {
		params.Country = cert.Subject.Country[0]
	}
	if len(cert.Subject.Province) > 0 {
		params.State = cert.Subject.Province[0]
	}
	if len(cert.Subject.Locality) > 0 {
		params.Locality = cert.Subject.Locality[0]
	}

	return params, nil
}

// roleNameFromExistingCert reads client.crt from the state directory and
// returns the role name encoded in its OtherName SAN, if present.
func (c *TPMReattestCommand) roleNameFromExistingCert() (string, error) {
	certPath := filepath.Join(c.FlagTPMStateDir, "client.crt")
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", certPath, err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("failed to decode certificate PEM from %s", certPath)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %w", err)
	}

	othersans, err := certutil.GetOtherSANsFromX509Extensions(cert.Extensions)
	if err != nil {
		return "", nil
	}

	for _, other := range othersans {
		if other.Oid == consts.TPMAuthOIDRoleName {
			return other.Value, nil
		}
	}

	return "", nil
}
