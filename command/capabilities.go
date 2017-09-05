package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*CapabilitiesCommand)(nil)
var _ cli.CommandAutocomplete = (*CapabilitiesCommand)(nil)

// CapabilitiesCommand is a Command that enables a new endpoint.
type CapabilitiesCommand struct {
	*BaseCommand
}

func (c *CapabilitiesCommand) Synopsis() string {
	return "Fetchs the capabilities of a token"
}

func (c *CapabilitiesCommand) Help() string {
	helpText := `
Usage: vault capabilities [options] [TOKEN] PATH

  Fetches the capabilities of a token for a given path. If a TOKEN is provided
  as an argument, the "/sys/capabilities" endpoint and permission is used. If
  no TOKEN is  provided, the "/sys/capabilities-self" endpoint and permission
  is used with the locally authenticated token.

  List capabilities for the local token on the "secret/foo" path:

      $ vault capabilities secret/foo

  List capabilities for a token on the "cubbyhole/foo" path:

      $ vault capabilities 96ddf4bc-d217-f3ba-f9bd-017055595017 cubbyhole/foo

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *CapabilitiesCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *CapabilitiesCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *CapabilitiesCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *CapabilitiesCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	token := ""
	path := ""
	args = f.Args()
	switch {
	case len(args) == 1:
		path = args[0]
	case len(args) == 2:
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

	sort.Strings(capabilities)
	c.UI.Output(strings.Join(capabilities, ", "))
	return 0
}
