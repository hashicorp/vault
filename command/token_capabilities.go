// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*TokenCapabilitiesCommand)(nil)
	_ cli.CommandAutocomplete = (*TokenCapabilitiesCommand)(nil)
)

type TokenCapabilitiesCommand struct {
	*BaseCommand

	flagAccessor bool
}

func (c *TokenCapabilitiesCommand) Synopsis() string {
	return "Print capabilities of a token on a path"
}

func (c *TokenCapabilitiesCommand) Help() string {
	helpText := `
Usage: vault token capabilities [options] [TOKEN | ACCESSOR] PATH

  Fetches the capabilities of a token or accessor for a given path. If a TOKEN
  is provided as an argument, the "/sys/capabilities" endpoint is used, which
  returns the capabilities of the provided TOKEN. If an ACCESSOR is provided
  as an argument along with the -accessor option, the "/sys/capabilities-accessor"
  endpoint is used, which returns the capabilities of the token referenced by
  ACCESSOR. If no TOKEN is provided, the "/sys/capabilities-self" endpoint
  is used, which returns the capabilities of the locally authenticated token.

  List capabilities for the local token on the "secret/foo" path:

      $ vault token capabilities secret/foo

  List capabilities for a token on the "cubbyhole/foo" path:

      $ vault token capabilities 96ddf4bc-d217-f3ba-f9bd-017055595017 cubbyhole/foo

  List capabilities for a token on the "cubbyhole/foo" path via its accessor:

      $ vault token capabilities -accessor 9793c9b3-e04a-46f3-e7b8-748d7da248da cubbyhole/foo

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TokenCapabilitiesCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:       "accessor",
		Target:     &c.flagAccessor,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage:      "Treat the argument as an accessor instead of a token.",
	})

	return set
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
	switch {
	case c.flagAccessor && len(args) < 2:
		c.UI.Error(fmt.Sprintf("Not enough arguments with -accessor (expected 2, got %d)", len(args)))
		return 1
	case c.flagAccessor && len(args) > 2:
		c.UI.Error(fmt.Sprintf("Too many arguments with -accessor (expected 2, got %d)", len(args)))
		return 1
	case len(args) == 0:
		c.UI.Error("Not enough arguments (expected 1-2, got 0)")
		return 1
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
	switch {
	case token == "":
		capabilities, err = client.Sys().CapabilitiesSelf(path)
	case c.flagAccessor:
		capabilities, err = client.Sys().CapabilitiesAccessor(token, path)
	default:
		capabilities, err = client.Sys().Capabilities(token, path)
	}

	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing capabilities: %s", err))
		return 2
	}
	if capabilities == nil {
		c.UI.Error("No capabilities found")
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
