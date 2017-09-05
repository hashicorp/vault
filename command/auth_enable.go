package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*AuthEnableCommand)(nil)
var _ cli.CommandAutocomplete = (*AuthEnableCommand)(nil)

// AuthEnableCommand is a Command that enables a new endpoint.
type AuthEnableCommand struct {
	*BaseCommand

	flagDescription string
	flagPath        string
	flagPluginName  string
	flagLocal       bool
}

func (c *AuthEnableCommand) Synopsis() string {
	return "Enables a new auth provider"
}

func (c *AuthEnableCommand) Help() string {
	helpText := `
Usage: vault auth-enable [options] TYPE

  Enables a new authentication provider. An authentication provider is
  responsible for authenticating users or machiens and assigning them
  policies with which they can access Vault.

  Enable the userpass auth provider at userpass/:

      $ vault auth-enable userpass

  Enable the LDAP auth provider at auth-prod/:

      $ vault auth-enable -path=auth-prod ldap

  Enable a custom auth plugin (after it is registered in the plugin registry):

      $ vault auth-enable -path=my-auth -plugin-name=my-auth-plugin plugin

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuthEnableCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "description",
		Target:     &c.flagDescription,
		Completion: complete.PredictAnything,
		Usage: "Human-friendly description for the purpose of this " +
			"authentication provider.",
	})

	f.StringVar(&StringVar{
		Name:       "path",
		Target:     &c.flagPath,
		Default:    "", // The default is complex, so we have to manually document
		Completion: complete.PredictAnything,
		Usage: "Place where the auth provider will be accessible. This must be " +
			"unique across all auth providers. This defaults to the \"type\" of " +
			"the mount. The auth provider will be accessible at \"/auth/<path>\".",
	})

	f.StringVar(&StringVar{
		Name:       "plugin-name",
		Target:     &c.flagPluginName,
		Completion: complete.PredictAnything,
		Usage: "Name of the auth provider plugin. This plugin name must already " +
			"exist in the Vault server's plugin catalog.",
	})

	f.BoolVar(&BoolVar{
		Name:    "local",
		Target:  &c.flagLocal,
		Default: false,
		Usage: "Mark the auth provider as local-only. Local auth providers are " +
			"not replicated nor removed by replication.",
	})

	return set
}

func (c *AuthEnableCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultAvailableAuths()
}

func (c *AuthEnableCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthEnableCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	authType := strings.TrimSpace(args[0])

	// If no path is specified, we default the path to the backend type
	// or use the plugin name if it's a plugin backend
	authPath := c.flagPath
	if authPath == "" {
		if authType == "plugin" {
			authPath = c.flagPluginName
		} else {
			authPath = authType
		}
	}

	// Append a trailing slash to indicate it's a path in output
	authPath = ensureTrailingSlash(authPath)

	if err := client.Sys().EnableAuthWithOptions(authPath, &api.EnableAuthOptions{
		Type:        authType,
		Description: c.flagDescription,
		Config: api.AuthConfigInput{
			PluginName: c.flagPluginName,
		},
		Local: c.flagLocal,
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error enabling %s auth: %s", authType, err))
		return 2
	}

	authThing := authType + " auth provider"
	if authType == "plugin" {
		authThing = c.flagPluginName + " plugin"
	}

	c.UI.Output(fmt.Sprintf("Success! Enabled %s at: %s", authThing, authPath))
	return 0
}
