package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*StatusCommand)(nil)
var _ cli.CommandAutocomplete = (*StatusCommand)(nil)

type StatusCommand struct {
	*BaseCommand
}

func (c *StatusCommand) Synopsis() string {
	return "Print seal and HA status"
}

func (c *StatusCommand) Help() string {
	helpText := `
Usage: vault status [options]

  Prints the current state of Vault including whether it is sealed and if HA
  mode is enabled. This command prints regardless of whether the Vault is
  sealed.

  The exit code reflects the seal status:

      - 0 - unsealed
      - 1 - error
      - 2 - sealed

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *StatusCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *StatusCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *StatusCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *StatusCommand) Run(args []string) int {
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
		// We return 2 everywhere else, but 2 is reserved for "sealed" here
		return 1
	}

	status, err := client.Sys().SealStatus()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error checking seal status: %s", err))
		return 1
	}

	// Do not return the int here yet, since we may want to return a custom error
	// code depending on the seal status.
	code := OutputSealStatus(c.UI, client, status)

	if status.Sealed {
		return 2
	}

	return code
}
