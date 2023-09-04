// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PluginRuntimeDeregisterCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginRuntimeDeregisterCommand)(nil)
)

type PluginRuntimeDeregisterCommand struct {
	*BaseCommand

	flagType string
}

func (c *PluginRuntimeDeregisterCommand) Synopsis() string {
	return "Deregister an existing plugin runtime in the catalog"
}

func (c *PluginRuntimeDeregisterCommand) Help() string {
	helpText := `
Usage: vault plugin runtime deregister [options] TYPE NAME

  Deregister an existing plugin runtime in the catalog with the given name. If
  any registered plugin references the plugin runtime, an error is returned. If
  the plugin runtime does not exist, an error is returned. The argument of type
  takes "container".

  Deregister a plugin runtime:

      $ vault plugin runtime deregister -type=container my-plugin-runtime

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginRuntimeDeregisterCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "type",
		Target:     &c.flagType,
		Completion: complete.PredictAnything,
		Usage:      "Plugin runtime type. Vault currently only supports \"container\" runtime type.",
	})

	return set
}

func (c *PluginRuntimeDeregisterCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *PluginRuntimeDeregisterCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginRuntimeDeregisterCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	runtimeTyeRaw := strings.TrimSpace(c.flagType)
	if len(runtimeTyeRaw) == 0 {
		c.UI.Error("-type is required for plugin runtime registration")
		return 1
	}

	runtimeType, err := api.ParsePluginRuntimeType(runtimeTyeRaw)
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	var runtimeNameRaw string
	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1

	// This case should come after invalid cases have been checked
	case len(args) == 1:
		runtimeNameRaw = args[0]
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	runtimeName := strings.TrimSpace(runtimeNameRaw)
	if err = client.Sys().DeregisterPluginRuntime(context.Background(), &api.DeregisterPluginRuntimeInput{
		Name: runtimeName,
		Type: runtimeType,
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error deregistering plugin runtime named %s: %s", runtimeName, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Deregistered plugin runtime (if it was registered): %s", runtimeName))
	return 0
}
