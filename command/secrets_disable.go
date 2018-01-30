package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*SecretsDisableCommand)(nil)
var _ cli.CommandAutocomplete = (*SecretsDisableCommand)(nil)

type SecretsDisableCommand struct {
	*BaseCommand
}

func (c *SecretsDisableCommand) Synopsis() string {
	return "Disable a secret engine"
}

func (c *SecretsDisableCommand) Help() string {
	helpText := `
Usage: vault secrets disable [options] PATH

  Disables a secrets engine at the given PATH. The argument corresponds to
  the enabled PATH of the engine, not the TYPE! All secrets created by this
  engine are revoked and its Vault data is removed.

  Disable the secrets engine enabled at aws/:

      $ vault secrets disable aws/

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *SecretsDisableCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *SecretsDisableCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultMounts()
}

func (c *SecretsDisableCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *SecretsDisableCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	path := ensureTrailingSlash(sanitizePath(args[0]))

	if err := client.Sys().Unmount(path); err != nil {
		c.UI.Error(fmt.Sprintf("Error disabling secrets engine at %s: %s", path, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Disabled the secrets engine (if it existed) at: %s", path))
	return 0
}
