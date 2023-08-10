// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*LeaseLookupCommand)(nil)
	_ cli.CommandAutocomplete = (*LeaseLookupCommand)(nil)
)

type LeaseLookupCommand struct {
	*BaseCommand
}

func (c *LeaseLookupCommand) Synopsis() string {
	return "Lookup the lease of a secret"
}

func (c *LeaseLookupCommand) Help() string {
	helpText := `
Usage: vault lease lookup ID

  Lookup the lease information of a secret.

  Every secret in Vault has a lease associated with it. Users can look up
  information on the lease by referencing the lease ID.

  Lookup lease of a secret:

      $ vault lease lookup database/creds/readonly/2f6a614c...

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *LeaseLookupCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	return set
}

func (c *LeaseLookupCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *LeaseLookupCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *LeaseLookupCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	leaseID := ""

	args = f.Args()
	switch len(args) {
	case 0:
		c.UI.Error("Missing ID!")
		return 1
	case 1:
		leaseID = strings.TrimSpace(args[0])
	default:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	secret, err := client.Sys().Lookup(leaseID)
	if err != nil {
		c.UI.Error(fmt.Sprintf("error looking up lease id %s: %s", leaseID, err))
		return 2
	}

	return OutputSecret(c.UI, secret)
}
