package command

import (
	"fmt"
	"strings"

	semver "github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PluginDeregisterCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginDeregisterCommand)(nil)
)

type PluginDeregisterCommand struct {
	*BaseCommand
}

func (c *PluginDeregisterCommand) Synopsis() string {
	return "Deregister an existing plugin in the catalog"
}

func (c *PluginDeregisterCommand) Help() string {
	helpText := `
Usage: vault plugin deregister [options] TYPE NAME

  Deregister an existing plugin in the catalog. If the plugin does not exist,
  no action is taken (the command is idempotent). The argument of type
  takes "auth", "database", or "secret".

  Deregister the plugin named my-custom-plugin:

      $ vault plugin deregister auth my-custom-plugin [version]

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginDeregisterCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *PluginDeregisterCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultPlugins(consts.PluginTypeUnknown)
}

func (c *PluginDeregisterCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginDeregisterCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	var pluginNameRaw, pluginTypeRaw, pluginVersionRaw string
	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, 2, or 3, got %d)", len(args)))
		return 1
	case len(args) > 3:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, 2, or 3, got %d)", len(args)))
		return 1

	// These cases should come after invalid cases have been checked
	case len(args) == 1:
		pluginTypeRaw = "unknown"
		pluginNameRaw = args[0]
	case len(args) == 2:
		pluginTypeRaw = args[0]
		pluginNameRaw = args[1]
	case len(args) == 3:
		pluginTypeRaw = args[0]
		pluginNameRaw = args[1]
		pluginVersionRaw = args[2]
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	pluginType, err := consts.ParsePluginType(strings.TrimSpace(pluginTypeRaw))
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}
	pluginName := strings.TrimSpace(pluginNameRaw)
	pluginVersion := strings.TrimSpace(pluginVersionRaw)
	if pluginVersion != "" {
		semanticVersion, err := semver.NewSemver(pluginVersion)
		if err != nil {
			c.UI.Error(fmt.Sprintf("version %q is not a valid semantic version: %v", pluginVersionRaw, err))
			return 2
		}

		// Canonicalize the version string.
		// Add the 'v' back in, since semantic version strips it out, and we want to be consistent with internal plugins.
		pluginVersion = "v" + semanticVersion.String()
	}

	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    pluginName,
		Type:    pluginType,
		Version: pluginVersion,
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error deregistering plugin named %s: %s", pluginName, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Deregistered plugin (if it was registered): %s", pluginName))
	return 0
}
