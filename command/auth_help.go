package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*AuthHelpCommand)(nil)
var _ cli.CommandAutocomplete = (*AuthHelpCommand)(nil)

// AuthHelpCommand is a Command that prints help output for a given auth
// provider
type AuthHelpCommand struct {
	*BaseCommand

	Handlers map[string]AuthHandler
}

func (c *AuthHelpCommand) Synopsis() string {
	return "Prints usage for an auth provider"
}

func (c *AuthHelpCommand) Help() string {
	helpText := `
Usage: vault path-help [options] TYPE | PATH

  Prints usage and help for an authentication provider. If provided a TYPE,
  this command retrieves the default help for the given authentication
  provider of that type. If given a PATH, this command returns the help
  output for the authentication provider mounted at that path. If given a
  PATH argument, the path must exist and be mounted.

  Get usage instructions for the userpass authentication provider:

      $ vault auth-help userpass

  Print usage for the authentication provider mounted at my-provider/

      $ vault auth-help my-provider/:

  Each authentication provider produces its own help output. For additional
  information, please view the online documentation.

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
			c.UI.Error(fmt.Sprintf("Error listing authentication providers: %s", err))
			return 2
		}

		authPath := ensureTrailingSlash(sanitizePath(args[0]))
		auth, ok := auths[authPath]
		if !ok {
			c.UI.Error(fmt.Sprintf(
				"Error retrieving help: unknown authentication provider: %s", args[0]))
			return 1
		}

		authHandler, ok = c.Handlers[auth.Type]
		if !ok {
			c.UI.Error(wrapAtLength(fmt.Sprintf(
				"INTERNAL ERROR! Found an authentication provider mounted at %s, but "+
					"its type %q is not registered in Vault. This is a bug and should "+
					"be reported. Please open an issue at github.com/hashicorp/vault.",
				authPath, authType)))
			return 2
		}
	}

	c.UI.Output(authHandler.Help())
	return 0
}
