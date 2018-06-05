package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*AuthHelpCommand)(nil)
var _ cli.CommandAutocomplete = (*AuthHelpCommand)(nil)

type AuthHelpCommand struct {
	*BaseCommand

	Handlers map[string]LoginHandler
}

func (c *AuthHelpCommand) Synopsis() string {
	return "Prints usage for an auth method"
}

func (c *AuthHelpCommand) Help() string {
	helpText := `
Usage: vault auth help [options] TYPE | PATH

  Prints usage and help for an auth method.

    - If given a TYPE, this command prints the default help for the
      auth method of that type.

    - If given a PATH, this command prints the help output for the
      auth method enabled at that path. This path must already
      exist.

  Get usage instructions for the userpass auth method:

      $ vault auth help userpass

  Print usage for the auth method enabled at my-method/:

      $ vault auth help my-method/

  Each auth method produces its own help output.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuthHelpCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *AuthHelpCommand) AutocompleteArgs() complete.Predictor {
	handlers := make([]string, 0, len(c.Handlers))
	for k := range c.Handlers {
		handlers = append(handlers, k)
	}
	return complete.PredictSet(handlers...)
}

func (c *AuthHelpCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthHelpCommand) Run(args []string) int {
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

	// Start with the assumption that we have an auth type, not a path.
	authType := strings.TrimSpace(args[0])

	authHandler, ok := c.Handlers[authType]
	if !ok {
		// There was no auth type by that name, see if it's a mount
		auths, err := client.Sys().ListAuth()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error listing auth methods: %s", err))
			return 2
		}

		authPath := ensureTrailingSlash(sanitizePath(args[0]))
		auth, ok := auths[authPath]
		if !ok {
			c.UI.Warn(fmt.Sprintf(
				"No auth method available on the server at %q", authPath))
			return 1
		}

		authHandler, ok = c.Handlers[auth.Type]
		if !ok {
			c.UI.Warn(wrapAtLength(fmt.Sprintf(
				"No method-specific CLI handler available for auth method %q",
				authType)))
			return 2
		}
	}

	c.UI.Output(authHandler.Help())
	return 0
}
