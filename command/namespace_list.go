// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*NamespaceListCommand)(nil)
	_ cli.CommandAutocomplete = (*NamespaceListCommand)(nil)
)

type NamespaceListCommand struct {
	*BaseCommand
}

func (c *NamespaceListCommand) Synopsis() string {
	return "List child namespaces"
}

func (c *NamespaceListCommand) Help() string {
	helpText := `
Usage: vault namespace list [options]

  Lists the enabled child namespaces.

  List all enabled child namespaces:

      $ vault namespace list

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *NamespaceListCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "detailed",
		Target:  &c.flagDetailed,
		Default: false,
		Usage:   "Print detailed information such as namespace ID.",
	})

	return set
}

func (c *NamespaceListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *NamespaceListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *NamespaceListCommand) Run(args []string) int {
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

	secret, err := client.Logical().List("sys/namespaces")
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing namespaces: %s", err))
		return 2
	}

	_, ok := extractListData(secret)
	if Format(c.UI) != "table" {
		if secret == nil || secret.Data == nil || !ok {
			OutputData(c.UI, map[string]interface{}{})
			return 2
		}
	}

	if secret == nil {
		c.UI.Error("No namespaces found")
		return 2
	}

	// There could be e.g. warnings
	if secret.Data == nil {
		return OutputSecret(c.UI, secret)
	}

	if secret.WrapInfo != nil && secret.WrapInfo.TTL != 0 {
		return OutputSecret(c.UI, secret)
	}

	if !ok {
		c.UI.Error("No entries found")
		return 2
	}

	if c.flagDetailed && Format(c.UI) != "table" {
		return OutputData(c.UI, secret.Data["key_info"])
	}

	return OutputList(c.UI, secret)
}
