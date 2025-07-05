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

func NewPluginRegisterCommand(baseCommand *BaseCommand) cli.Command {
	return &PluginRegisterCommand{
		BaseCommand: baseCommand,
	}
}

type PluginRegisterCommand struct {
	*BaseCommand

	flagArgs     []string
	flagCommand  string
	flagSHA256   string
	flagVersion  string
	flagOCIImage string
	flagRuntime  string
	flagEnv      []string
	flagDownload bool
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

  Register the plugin named my-custom-plugin:

      $ vault plugin register -sha256=d3f0a8b... -version=v1.0.0 auth my-custom-plugin

  Register a plugin with custom arguments:

      $ vault plugin register \
          -sha256=d3f0a8b... \
          -version=v1.0.0 \
          -args=--with-glibc,--with-cgo \
          auth my-custom-plugin

  Register a plugin with -download (enterprise only):

      $ vault plugin register \
          -version=v0.17.0+ent \
          -download=true \
          secret vault-plugin-secrets-keymgmt
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
		Usage: "Command to spawn the plugin. If -sha256 is provided to register with a plugin binary, " +
			"this defaults to the name of the plugin if both oci_image and command are unspecified. " +
			"Otherwise, if -sha256 is not provided, a plugin artifact is expected for registration, and " +
			"this will be ignored because the run command is known.",
	})

	f.StringVar(&StringVar{
		Name:       "sha256",
		Target:     &c.flagSHA256,
		Completion: complete.PredictAnything,
		Usage: "SHA256 of the plugin binary or the OCI image provided. " +
			"This is required to register with a plugin binary but should not be " +
			"specified when registering with a plugin artifact.",
	})

	f.StringVar(&StringVar{
		Name:       "version",
		Target:     &c.flagVersion,
		Completion: complete.PredictAnything,
		Usage: "Semantic version of the plugin. Used as the tag when specifying oci_image, but with any leading 'v' trimmed. " +
			"This is required to register with a plugin artifact but optional when registering with a plugin binary.",
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
		Name:       "download",
		Target:     &c.flagDownload,
		Completion: complete.PredictAnything,
		Usage: "Enterprise only. If set, Vault will automatically download plugins from" +
			"releases.hashicorp.com",
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
	case c.flagSHA256 == "" && c.flagVersion == "":
		c.UI.Error("One of -sha256 or -version is required. " +
			"If registering with a binary, please provide at least -sha256 (-version optional)." +
			"If registering with an artifact, please provide -version only.")
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
	if c.flagSHA256 != "" && (command == "" && c.flagOCIImage == "") {
		command = pluginName
	}

	resp, err := client.Sys().RegisterPluginDetailed(&api.RegisterPluginInput{
		Name:     pluginName,
		Type:     pluginType,
		Args:     c.flagArgs,
		Command:  command,
		SHA256:   c.flagSHA256,
		Version:  c.flagVersion,
		OCIImage: c.flagOCIImage,
		Runtime:  c.flagRuntime,
		Env:      c.flagEnv,
		Download: c.flagDownload,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error registering plugin %s: %s", pluginName, err))
		return 2
	}

	if resp != nil && len(resp.Warnings) > 0 {
		c.UI.Warn(wrapAtLength(fmt.Sprintf(
			"Warnings while registering plugin %s: %s",
			pluginName,
			strings.Join(resp.Warnings, "\n\n"),
		)) + "\n")
	}

	c.UI.Output(fmt.Sprintf("Success! Registered plugin: %s", pluginName))
	return 0
}
