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
	_ cli.Command             = (*PluginRuntimeInfoCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginRuntimeInfoCommand)(nil)
)

type PluginRuntimeInfoCommand struct {
	*BaseCommand

	flagType string
}

func (c *PluginRuntimeInfoCommand) Synopsis() string {
	return "Read information about a plugin runtime in the catalog"
}

func (c *PluginRuntimeInfoCommand) Help() string {
	helpText := `
Usage: vault plugin runtime info [options] NAME

  Displays information about a plugin runtime in the catalog with the given name. If
  the plugin runtime does not exist, an error is returned. The -type flag
  currently only accepts "container".

  Get info about a plugin runtime:

      $ vault plugin runtime info -type=container my-plugin-runtime

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginRuntimeInfoCommand) Flags() *FlagSets {
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

func (c *PluginRuntimeInfoCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *PluginRuntimeInfoCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginRuntimeInfoCommand) Run(args []string) int {
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
	resp, err := client.Sys().GetPluginRuntime(context.Background(), &api.GetPluginRuntimeInput{
		Name: runtimeName,
		Type: runtimeType,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading plugin runtime named %s: %s", runtimeName, err))
		return 2
	}

	if resp == nil {
		c.UI.Error(fmt.Sprintf("No value found for plugin runtime %q", runtimeName))
		return 2
	}

	data := map[string]interface{}{
		"name":          resp.Name,
		"type":          resp.Type,
		"oci_runtime":   resp.OCIRuntime,
		"cgroup_parent": resp.CgroupParent,
		"cpu_nanos":     resp.CPU,
		"memory_bytes":  resp.Memory,
	}

	if c.flagField != "" {
		return PrintRawField(c.UI, data, c.flagField)
	}
	return OutputData(c.UI, data)
}
