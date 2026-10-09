// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	sdktpm "github.com/hashicorp/vault/sdk/helper/tpm"
)

var _ cli.Command = (*TPMEnrollCommand)(nil)

type TPMEnrollCommand struct {
	*base.BaseCommand

	flagEKPEM    string
	flagName     string
	flagMetadata map[string]string
}

func (c *TPMEnrollCommand) Synopsis() string {
	return "Enroll a TPM endorsement key with Vault"
}

func (c *TPMEnrollCommand) Help() string {
	helpText := `
Usage: vault tpm enroll [options]

  Registers the local TPM's endorsement key (EK) with the Vault identity
  backend so that the TPM can later be used for attestation. By default the EK
  is read directly from the local TPM device; alternatively, a PEM-encoded
  public key can be supplied via -ekpem to pre-enroll a TPM that is not
  attached to the current host.

  Enroll the local TPM, specifying a device path and metadata:

      $ vault tpm enroll -name=my-server \
          -tpm-device-path=/dev/tpm0 \
          -metadata=env=prod \
          -metadata=region=us-east-1

  Enroll a remote TPM by supplying its EK public key directly:

      $ vault tpm enroll -name=remote-server \
          -ekpem="$(cat ek-public.pem)" \
          -metadata=env=prod

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *TPMEnrollCommand) Flags() *base.FlagSets {
	set := c.FlagSetForBit(base.FlagSetHTTP | base.FlagSetOutputFormat | base.FlagSetOutputField)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&base.StringVar{
		Name:    "ekpem",
		Target:  &c.flagEKPEM,
		Default: "",
		Usage:   "Enrollment key pem string",
	})

	f.StringVar(&base.StringVar{
		Name:    "name",
		Target:  &c.flagName,
		Default: "",
		Usage:   "Name of the new Enrollment Key",
	})

	f.StringMapVar(&base.StringMapVar{
		Name:    "metadata",
		Target:  &c.flagMetadata,
		Default: map[string]string{},
		Usage:   "Metadata for the new EK in key value pairs",
	})

	return set
}

func (c *TPMEnrollCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if strings.TrimSpace(c.flagEKPEM) == "" {
		eks, err := sdktpm.EndorsementKeys(sdktpm.TPMConfig{
			DevicePath: c.FlagTPMDevicePath,
		})
		if err != nil {
			c.UI.Error(fmt.Sprintf("error reading endorsement keys: %s", err))
			return 2
		}
		if len(eks) == 0 {
			c.UI.Error("no endorsement keys found on TPM")
			return 2
		}

		pubDER, err := x509.MarshalPKIXPublicKey(eks[0].Public)
		if err != nil {
			c.UI.Error(fmt.Sprintf("error marshaling EK public key: %s", err))
			return 2
		}
		c.flagEKPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	data := map[string]interface{}{
		"tpm_ek_public_key": c.flagEKPEM,
	}

	if c.flagName != "" {
		data["name"] = c.flagName
	}
	if len(c.flagMetadata) > 0 {
		data["metadata"] = c.flagMetadata
	}

	secret, err := client.Logical().Write("identity/tpm", data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("error enrolling TPM: %s", err))
		return 2
	}

	if c.FlagField != "" {
		return base.PrintRawField(c.UI, secret, c.FlagField)
	}

	return base.OutputData(c.UI, secret.Data)
}
