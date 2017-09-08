package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*AuthDisableCommand)(nil)
var _ cli.CommandAutocomplete = (*AuthDisableCommand)(nil)

type AuthDisableCommand struct {
	*BaseCommand
}

func (c *AuthDisableCommand) Synopsis() string {
	return "Disables an auth method"
}

func (c *AuthDisableCommand) Help() string {
	helpText := `
Usage: vault auth disable [options] PATH

  Disables an existing auth method at the given PATH. The argument corresponds
  to the PATH of the mount, not the TYPE!. Once the auth method is disabled its
  path can no longer be used to authenticate.

  All access tokens generated via the disabled auth method are immediately
  revoked. This command will block until all tokens are revoked.

  Disable the auth method at userpass/:

      $ vault auth disable userpass/

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuthDisableCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *AuthDisableCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultAuths()
}

func (c *AuthDisableCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthDisableCommand) Run(args []string) int {
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

	path := ensureTrailingSlash(sanitizePath(args[0]))

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if err := client.Sys().DisableAuth(path); err != nil {
		c.UI.Error(fmt.Sprintf("Error disabling auth method at %s: %s", path, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Disabled the auth method (if it existed) at: %s", path))
	return 0
}
