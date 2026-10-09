// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
)

var _ cli.Command = (*TPMCleanupCommand)(nil)

// tpmStateFiles is the fixed set of files written by SaveToDirectory.
var tpmStateFiles = []string{
	"client-key.json",
	"ak.blob",
	"app.blob",
	"ca_chain.pem",
	"client.crt",
}

type TPMCleanupCommand struct {
	*base.BaseCommand
}

func (c *TPMCleanupCommand) Synopsis() string {
	return "Remove attestation artifacts from the state directory"
}

func (c *TPMCleanupCommand) Help() string {
	helpText := `
Usage: vault tpm cleanup [options]

  Removes the certificate and key material written to the state directory by
  "vault tpm attest". The directory itself is preserved. Any other files in
  the directory are not affected.

  The following files are removed:

      client.crt       - the issued certificate (PEM)
      ca_chain.pem     - the CA chain (PEM)
      client-key.json  - key reference metadata (JSON)
      ak.blob          - attestation key blob
      app.blob         - application key blob

  Clean up the state in the current directory:

      $ vault tpm cleanup

  Clean up a specific state directory:

      $ vault tpm cleanup -tpm-state-dir=/etc/vault/tpm

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *TPMCleanupCommand) Flags() *base.FlagSets {
	return c.FlagSetForBit(base.FlagSetHTTP)
}

// Run accesses the state directory, removes each fixed file from the State Directory, and outputs missing and removed files.
func (c *TPMCleanupCommand) Run(args []string) int {
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	info, err := os.Stat(c.FlagTPMStateDir)
	if err != nil {
		c.UI.Error(fmt.Sprintf("error accessing state directory %q: %s", c.FlagTPMStateDir, err))
		return 2
	}
	if !info.IsDir() {
		c.UI.Error(fmt.Sprintf("state directory path %q is not a directory", c.FlagTPMStateDir))
		return 2
	}

	var removed, missing []string
	for _, name := range tpmStateFiles {
		path := filepath.Join(c.FlagTPMStateDir, name)
		err := os.Remove(path)
		switch {
		case err == nil:
			removed = append(removed, name)
		case os.IsNotExist(err):
			missing = append(missing, name)
		default:
			c.UI.Error(fmt.Sprintf("error removing %s: %s", name, err))
			return 2
		}
	}

	for _, name := range removed {
		c.UI.Output(fmt.Sprintf("Removed %s", name))
	}
	if len(missing) > 0 {
		c.UI.Warn(fmt.Sprintf("Not found (already removed?): %s", strings.Join(missing, ", ")))
	}

	c.UI.Output(fmt.Sprintf("Cleanup complete. %d file(s) removed from %s", len(removed), c.FlagTPMStateDir))
	return 0
}
