// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PluginRegisterCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginRegisterCommand)(nil)
)

type PluginRegisterCommand struct {
	*BaseCommand

	flagArgs     []string
	flagCommand  string
	flagSHA256   string
	flagVersion  string
	flagOCIImage string
	flagRuntime  string
	flagEnv      []string
	flagManaged  bool
}

func (c *PluginRegisterCommand) Synopsis() string {
	return "Registers a new plugin in the catalog"
}

func (c *PluginRegisterCommand) Help() string {
	helpText := `
Usage: vault plugin register [options] TYPE NAME

  Registers a new plugin in the catalog. The plugin binary must exist in Vault's
  configured plugin directory. The argument of type takes "auth", "database",
  or "secret".

  Register and automatically download a plugin named my-custom-plugin:

  	  $ vault plugin register -managed -version=v1.0.0 auth my-custom-plugin

  Register the plugin named my-custom-plugin:

      $ vault plugin register -sha256=d3f0a8b... -version=v1.0.0 auth my-custom-plugin

  Register a plugin with custom arguments:

      $ vault plugin register \
          -sha256=d3f0a8b... \
          -version=v1.0.0 \
          -args=--with-glibc,--with-cgo \
          auth my-custom-plugin

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginRegisterCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringSliceVar(&StringSliceVar{
		Name:       "args",
		Target:     &c.flagArgs,
		Completion: complete.PredictAnything,
		Usage: "Argument to pass to the plugin when starting. This " +
			"flag can be specified multiple times to specify multiple args.",
	})

	f.StringVar(&StringVar{
		Name:       "command",
		Target:     &c.flagCommand,
		Completion: complete.PredictAnything,
		Usage: "Command to spawn the plugin. This defaults to the name of the " +
			"plugin if both oci_image and command are unspecified.",
	})

	f.StringVar(&StringVar{
		Name:       "sha256",
		Target:     &c.flagSHA256,
		Completion: complete.PredictAnything,
		Usage:      "SHA256 of the plugin binary or the oci_image provided. This is required for all plugins.",
	})

	f.StringVar(&StringVar{
		Name:       "version",
		Target:     &c.flagVersion,
		Completion: complete.PredictAnything,
		Usage:      "Semantic version of the plugin. Used as the tag when specifying oci_image, but with any leading 'v' trimmed. Optional.",
	})

	f.StringVar(&StringVar{
		Name:       "oci_image",
		Target:     &c.flagOCIImage,
		Completion: complete.PredictAnything,
		Usage: "OCI image to run. If specified, setting command, args, and env will update the " +
			"container's entrypoint, args, and environment variables (append-only) respectively.",
	})

	f.StringVar(&StringVar{
		Name:       "runtime",
		Target:     &c.flagRuntime,
		Completion: complete.PredictAnything,
		Usage:      "Vault plugin runtime to use if oci_image is specified.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:       "env",
		Target:     &c.flagEnv,
		Completion: complete.PredictAnything,
		Usage: "Environment variables to set for the plugin when starting. This " +
			"flag can be specified multiple times to specify multiple environment variables.",
	})

	f.BoolVar(&BoolVar{
		Name:   "managed",
		Target: &c.flagManaged,
		Usage:  "Automatically download and manage the lifecycle of the plugin binary.",
	})

	return set
}

func (c *PluginRegisterCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultPlugins(api.PluginTypeUnknown)
}

func (c *PluginRegisterCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginRegisterCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	var pluginNameRaw, pluginTypeRaw string
	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1 or 2, got %d)", len(args)))
		return 1
	case len(args) > 2:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1 or 2, got %d)", len(args)))
		return 1
	case c.flagSHA256 == "" && !c.flagManaged:
		c.UI.Error("-sha256 or -managed is required for all plugins")
		return 1
	case c.flagManaged && len(args) != 2:
		c.UI.Error("Must specify plugin type and name when using -managed flag")
		return 1

	// These cases should come after invalid cases have been checked
	case len(args) == 1:
		pluginTypeRaw = "unknown"
		pluginNameRaw = args[0]
	case len(args) == 2:
		pluginTypeRaw = args[0]
		pluginNameRaw = args[1]
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	pluginType, err := api.ParsePluginType(strings.TrimSpace(pluginTypeRaw))
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}
	pluginName := strings.TrimSpace(pluginNameRaw)

	command := c.flagCommand
	if command == "" {
		switch {
		case c.flagOCIImage != "":
			// Containerized plugins don't require a command.
		case c.flagManaged:
			c.UI.Info(fmt.Sprintf("Expecting plugin binary vault-plugin-%s-%s, use -command to override", pluginTypeRaw, pluginName))
			command = fmt.Sprintf("vault-plugin-%s-%s", pluginTypeRaw, pluginName)
		default:
			command = pluginName
		}
	}

	if err := client.Sys().RegisterPlugin(&api.RegisterPluginInput{
		Name:     pluginName,
		Type:     pluginType,
		Args:     c.flagArgs,
		Command:  command,
		SHA256:   c.flagSHA256,
		Version:  c.flagVersion,
		OCIImage: c.flagOCIImage,
		Runtime:  c.flagRuntime,
		Env:      c.flagEnv,
		Managed:  c.flagManaged,
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error registering plugin %s: %s", pluginName, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Registered plugin: %s", pluginName))
	return 0
}
