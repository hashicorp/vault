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
	_ cli.Command             = (*PluginRuntimeListCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginRuntimeListCommand)(nil)
)

type PluginRuntimeListCommand struct {
	*BaseCommand

	flagType string
}

func (c *PluginRuntimeListCommand) Synopsis() string {
	return "Lists available plugin runtimes"
}

func (c *PluginRuntimeListCommand) Help() string {
	helpText := `
Usage: vault plugin runtime list [options]

  Lists available plugin runtimes registered in the catalog. This does not list whether
  plugin runtimes are in use, but rather just their availability.

  List all available plugin runtimes in the catalog:

      $ vault plugin runtime list

  List all available container plugin runtimes in the catalog:

      $ vault plugin runtime list -type=container

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginRuntimeListCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "type",
		Target:     &c.flagType,
		Completion: complete.PredictAnything,
		Usage:      "Plugin runtime type. Vault currently only supports \"container\" runtime type.",
	})

	return set
}

func (c *PluginRuntimeListCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *PluginRuntimeListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginRuntimeListCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if len(f.Args()) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	var input *api.ListPluginRuntimesInput
	runtimeTyeRaw := strings.TrimSpace(c.flagType)
	if len(runtimeTyeRaw) > 0 {
		runtimeType, err := api.ParsePluginRuntimeType(runtimeTyeRaw)
		if err != nil {
			c.UI.Error(err.Error())
			return 2
		}
		input = &api.ListPluginRuntimesInput{Type: runtimeType}
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	resp, err := client.Sys().ListPluginRuntimes(context.Background(), input)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing available plugin runtimes: %s", err))
		return 2
	}
	if resp == nil {
		c.UI.Error("No tableResponse from server when listing plugin runtimes")
		return 2
	}

	switch Format(c.UI) {
	case "table":
		c.UI.Output(tableOutput(c.tableResponse(resp), nil))
		return 0
	default:
		return OutputData(c.UI, resp.Runtimes)
	}
}

func (c *PluginRuntimeListCommand) tableResponse(response *api.ListPluginRuntimesResponse) []string {
	out := []string{"Name | Type | OCI Runtime | Parent Cgroup | CPU Nanos | Memory Bytes"}
	for _, runtime := range response.Runtimes {
		out = append(out, fmt.Sprintf("%s | %s | %s | %s | %d | %d",
			runtime.Name, runtime.Type, runtime.OCIRuntime, runtime.CgroupParent, runtime.CPU, runtime.Memory))
	}

	return out
}
