package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// SealCommand is a Command that seals the vault.
type SealCommand struct {
	meta.Meta
}

func (c *SealCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("seal", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	if err := client.Sys().Seal(); err != nil {
		c.Ui.Error(fmt.Sprintf("Error sealing: %s", err))
		return 1
	}

	c.Ui.Output("Vault is now sealed.")
	return 0
}

func (c *SealCommand) Synopsis() string {
	return "Seals the Vault server"
}

func (c *SealCommand) Help() string {
	helpText := `
Usage: vault seal [options]

  Seal the vault.

  Sealing a vault tells the Vault server to stop responding to any
  access operations until it is unsealed again. A sealed vault throws away
  its master key to unlock the data, so it is physically blocked from
  responding to operations again until the vault is unsealed with
  the "unseal" command or via the API.

  This command is idempotent, if the vault is already sealed it does nothing.

  If an unseal has started, sealing the vault will reset the unsealing
  process. You'll have to re-enter every portion of the master key again.
  This is the same as running "vault unseal -reset".

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
