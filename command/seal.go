package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*SealCommand)(nil)
var _ cli.CommandAutocomplete = (*SealCommand)(nil)

// SealCommand is a Command that seals the vault.
type SealCommand struct {
	*BaseCommand
}

func (c *SealCommand) Synopsis() string {
	return "Seals the Vault server"
}

func (c *SealCommand) Help() string {
	helpText := `
Usage: vault seal [options]

  Seals the Vault server. Sealing tells the Vault server to stop responding
  to any operations until it is unsealed. When sealed, the Vault server
  discards its in-memory master key to unlock the data, so it is physically
  blocked from responding to operations unsealed.

  If an unseal is in progress, sealing the Vault will reset the unsealing
  process. Users will have to re-enter their portions of the master key again.

  This command does nothing if the Vault server is already sealed.

  Seal the Vault server:

      $ vault seal

  For a full list of examples and why you might want to seal the Vault, please
  see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *SealCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *SealCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *SealCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *SealCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if err := client.Sys().Seal(); err != nil {
		c.UI.Error(fmt.Sprintf("Error sealing: %s", err))
		return 2
	}

	c.UI.Output("Success! Vault is sealed.")
	return 0
}
