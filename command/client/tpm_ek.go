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
	"github.com/hashicorp/vault/internalshared/tpmutil"
	sdktpm "github.com/hashicorp/vault/sdk/helper/tpm"
)

var _ cli.Command = (*TPMEkCommand)(nil)

type TPMEkCommand struct {
	*base.BaseCommand
}

func (c *TPMEkCommand) Synopsis() string {
	return "Print the TPM endorsement key public key as PEM and the TPM ID"
}

func (c *TPMEkCommand) Help() string {
	helpText := `
Usage: vault tpm ek

  Reads the endorsement key (EK) from the local TPM device and prints its
  public key in PEM format (type "PUBLIC KEY") along with the TPM ID. 

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *TPMEkCommand) Flags() *base.FlagSets {
	return c.FlagSetForBit(base.FlagSetHTTP | base.FlagSetOutputFormat | base.FlagSetOutputField)
}

func (c *TPMEkCommand) Run(args []string) int {
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

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

	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})

	tpmID := tpmutil.IDFromDER(pubDER)

	dataMap := map[string]interface{}{
		"tpm_ek_public_key": string(pemBytes),
		"tpm_id":            tpmID,
	}

	if c.FlagField != "" {
		return base.PrintRawField(c.UI, dataMap, c.FlagField)
	}

	return base.OutputData(c.UI, dataMap)
}
