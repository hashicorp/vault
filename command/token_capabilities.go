package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*TokenCapabilitiesCommand)(nil)
var _ cli.CommandAutocomplete = (*TokenCapabilitiesCommand)(nil)

type TokenCapabilitiesCommand struct {
	*BaseCommand
}

func (c *TokenCapabilitiesCommand) Synopsis() string {
	return "Print capabilities of a token on a path"
}

func (c *TokenCapabilitiesCommand) Help() string {
	helpText := `
Usage: vault token capabilities [options] [TOKEN] PATH

  Fetches the capabilities of a token for a given path. If a TOKEN is provided
  as an argument, the "/sys/capabilities" endpoint and permission is used. If
  no TOKEN is provided, the "/sys/capabilities-self" endpoint and permission
  is used with the locally authenticated token.

  List capabilities for the local token on the "secret/foo" path:

      $ vault token capabilities secret/foo

  List capabilities for a token on the "cubbyhole/foo" path:

      $ vault token capabilities 96ddf4bc-d217-f3ba-f9bd-017055595017 cubbyhole/foo

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TokenCapabilitiesCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *TokenCapabilitiesCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TokenCapabilitiesCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TokenCapabilitiesCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	token := ""
	path := ""
	args = f.Args()
	switch len(args) {
	case 0:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1-2, got 0)"))
		return 1
	case 1:
		path = args[0]
	case 2:
		token, path = args[0], args[1]
	default:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1-2, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	var capabilities []string
	if token == "" {
		capabilities, err = client.Sys().CapabilitiesSelf(path)
	} else {
		capabilities, err = client.Sys().Capabilities(token, path)
	}
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing capabilities: %s", err))
		return 2
	}
	if capabilities == nil {
		c.UI.Error(fmt.Sprintf("No capabilities found"))
		return 1
	}

	switch Format(c.UI) {
	case "table":
		sort.Strings(capabilities)
		c.UI.Output(strings.Join(capabilities, ", "))
		return 0
	default:
		return OutputData(c.UI, capabilities)
	}
}
