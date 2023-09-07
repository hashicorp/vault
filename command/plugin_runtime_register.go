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
	_ cli.Command             = (*PluginRuntimeRegisterCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginRuntimeRegisterCommand)(nil)
)

type PluginRuntimeRegisterCommand struct {
	*BaseCommand

	flagType         string
	flagOCIRuntime   string
	flagCgroupParent string
	flagCPUNanos     int64
	flagMemoryBytes  int64
}

func (c *PluginRuntimeRegisterCommand) Synopsis() string {
	return "Registers a new plugin runtime in the catalog"
}

func (c *PluginRuntimeRegisterCommand) Help() string {
	helpText := `
Usage: vault plugin runtime register [options] NAME

  Registers a new plugin runtime in the catalog. Currently, Vault only supports registering runtime of container type.
The OCI runtime must be available on Vault's host. If no OCI runtime is specified, Vault will use \"runsc\", gVisor's OCI runtime.

  Register the plugin runtime named my-custom-plugin-runtime:

      $ vault plugin runtime register -type=container -oci_runtime=my-oci-runtime my-custom-plugin-runtime

  Register a plugin runtime with custom arguments:

      $ vault plugin runtime register \
          -oci_runtime=my-oci-runtime \
          -type=container my-custom-plugin

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginRuntimeRegisterCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "type",
		Target:     &c.flagType,
		Completion: complete.PredictAnything,
		Usage:      "Plugin runtime type. Vault currently only supports \"container\" runtime type.",
	})

	f.StringVar(&StringVar{
		Name:       "oci_runtime",
		Target:     &c.flagOCIRuntime,
		Completion: complete.PredictAnything,
		Usage:      "OCI runtime. Default is \"runsc\", gVisor's OCI runtime.",
	})

	f.StringVar(&StringVar{
		Name:       "cgroup_parent",
		Target:     &c.flagCgroupParent,
		Completion: complete.PredictAnything,
		Usage:      "Parent cgroup for the container runtime",
	})

	f.Int64Var(&Int64Var{
		Name:       "cpu_nanos",
		Target:     &c.flagCPUNanos,
		Completion: complete.PredictAnything,
		Usage:      "CPU limit to set per container in nanos. Defaults to no limit.",
	})

	f.Int64Var(&Int64Var{
		Name:       "memory_bytes",
		Target:     &c.flagMemoryBytes,
		Completion: complete.PredictAnything,
		Usage:      "Memory limit to set per container in bytes. Defaults to no limit.",
	})

	return set
}

func (c *PluginRuntimeRegisterCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *PluginRuntimeRegisterCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginRuntimeRegisterCommand) Run(args []string) int {
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
	ociRuntime := strings.TrimSpace(c.flagOCIRuntime)
	cgroupParent := strings.TrimSpace(c.flagCgroupParent)

	if err := client.Sys().RegisterPluginRuntime(context.Background(), &api.RegisterPluginRuntimeInput{
		Name:         runtimeName,
		Type:         runtimeType,
		OCIRuntime:   ociRuntime,
		CgroupParent: cgroupParent,
		CPU:          c.flagCPUNanos,
		Memory:       c.flagMemoryBytes,
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error registering plugin runtime %s: %s", runtimeName, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Registered plugin runtime: %s", runtimeName))
	return 0
}
