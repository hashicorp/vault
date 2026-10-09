// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	sdktpm "github.com/hashicorp/vault/sdk/helper/tpm"
	"github.com/posener/complete"
)

var _ cli.Command = (*TPMAttestCommand)(nil)

type TPMAttestCommand struct {
	*base.BaseCommand

	flagMountPath string
	flagRoleName  string
	flagSubjectCN string
	flagSubjectO  string
	flagSubjectOU string
	flagSubjectC  string
	flagSubjectST string
	flagSubjectL  string
}

func (c *TPMAttestCommand) Synopsis() string {
	return "Go through Attestation and update state directory on success"
}

func (c *TPMAttestCommand) Help() string {
	helpText := `
Usage: vault tpm attest [options]

  Note: Attestation requires the TPM to have been registered in the auth mount's namespace.

  Performs TPM attestation against a Vault tpm auth method mount. The command
  opens the local TPM device, runs the two-phase challenge/response workflow
  (begin + finish), and on success writes the issued certificate and associated
  key material to the state directory. Additionally, it stores the private key in the TPM:

      client.crt       - the issued certificate (PEM)
      ca_chain.pem     - the CA chain (PEM)
      client-key.json  - key reference metadata (JSON)
      ak.blob          - attestation key blob
      app.blob         - application key blob

  Attest using the default TPM device and write output to the current directory:

      $ vault tpm attest -role-name=my-role

  Specify a custom TPM device path and state directory:

      $ vault tpm attest -role-name=my-role \
          -tpm-device-path=/dev/tpm0 \
          -tpm-state-dir=/etc/vault/tpm

  Include certificate subject fields:

      $ vault tpm attest -role-name=my-role \
          -cert-subject-CN=my-host.example.com \
          -cert-subject-O="My Org" \
          -cert-subject-C=US

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

// Flags registers each flag, usage text, and tab completion behavior
func (c *TPMAttestCommand) Flags() *base.FlagSets {
	set := c.FlagSetForBit(base.FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&base.StringVar{
		Name:       "mount-path",
		Target:     &c.flagMountPath,
		Default:    "auth/tpm",
		Completion: complete.PredictAnything,
		Usage:      "Mount path for the TPM auth method, relative to any specified namespace.",
	})

	f.StringVar(&base.StringVar{
		Name:    "role-name",
		Target:  &c.flagRoleName,
		Default: "",
		Usage:   "Name of the Vault TPM auth role to attest against. Required.",
	})

	f.StringVar(&base.StringVar{
		Name:    "cert-subject-CN",
		Target:  &c.flagSubjectCN,
		Default: "",
		Usage:   "Common name (CN) for the certificate subject. Defaults to the system hostname.",
	})

	f.StringVar(&base.StringVar{
		Name:    "cert-subject-O",
		Target:  &c.flagSubjectO,
		Default: "",
		Usage:   "Organization (O) for the certificate subject.",
	})

	f.StringVar(&base.StringVar{
		Name:    "cert-subject-OU",
		Target:  &c.flagSubjectOU,
		Default: "",
		Usage:   "Organizational unit (OU) for the certificate subject.",
	})

	f.StringVar(&base.StringVar{
		Name:    "cert-subject-C",
		Target:  &c.flagSubjectC,
		Default: "",
		Usage:   "Country (C) for the certificate subject (two-letter ISO code).",
	})

	f.StringVar(&base.StringVar{
		Name:    "cert-subject-ST",
		Target:  &c.flagSubjectST,
		Default: "",
		Usage:   "State or province (ST) for the certificate subject.",
	})

	f.StringVar(&base.StringVar{
		Name:    "cert-subject-L",
		Target:  &c.flagSubjectL,
		Default: "",
		Usage:   "Locality (L) for the certificate subject.",
	})

	return set
}

// Run goes through both phases of attestation. Function creates HTTP client, creates state directory, builds the attestation config,
// goes through attestation process, and saves the output to the state directory on success.
func (c *TPMAttestCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if strings.TrimSpace(c.flagRoleName) == "" {
		c.UI.Error("flag -role-name is required")
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("error creating client: %s", err))
		return 2
	}

	if err := os.MkdirAll(c.FlagTPMStateDir, 0o700); err != nil {
		c.UI.Error(fmt.Sprintf("error creating state directory %q: %s", c.FlagTPMStateDir, err))
		return 2
	}

	cfg := &sdktpm.AttestationConfig{
		MountPath: c.flagMountPath,
		RoleName:  c.flagRoleName,
		TPMConfig: sdktpm.TPMConfig{
			DevicePath: c.FlagTPMDevicePath,
		},
		CSRParams: &sdktpm.CSRParams{
			CommonName: c.flagSubjectCN,
			Country:    c.flagSubjectC,
			State:      c.flagSubjectST,
			Locality:   c.flagSubjectL,
			Org:        c.flagSubjectO,
			OrgUnit:    c.flagSubjectOU,
		},
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

	c.UI.Output(fmt.Sprintf("Attestation successful. Certificate valid until: %s",
		result.NotAfter.UTC().Format(time.RFC1123)))
	return 0
}
