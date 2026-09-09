// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/cli"
	semver "github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PluginDeregisterCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginDeregisterCommand)(nil)
)

type PluginDeregisterCommand struct {
	*BaseCommand

	flagPluginVersion string
}

func (c *PluginDeregisterCommand) Synopsis() string {
	return "Deregister an existing plugin in the catalog"
}

func (c *PluginDeregisterCommand) Help() string {
	helpText := `
Usage: vault plugin deregister [options] TYPE NAME

  Deregister an existing plugin in the catalog. If the plugin does not exist,
  no action is taken (the command is idempotent). The TYPE argument
  takes "auth", "database", or "secret".

  Deregister the unversioned auth plugin named my-custom-plugin:

      $ vault plugin deregister auth my-custom-plugin

  Deregister the auth plugin named my-custom-plugin, version 1.0.0:

      $ vault plugin deregister -version=v1.0.0 auth my-custom-plugin

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginDeregisterCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "version",
		Target:     &c.flagPluginVersion,
		Completion: complete.PredictAnything,
		Usage: "Semantic version of the plugin to deregister. If unset, " +
			"only an unversioned plugin may be deregistered.",
	})

	return set
}

func (c *PluginDeregisterCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultPlugins(api.PluginTypeUnknown)
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

	var pluginNameRaw, pluginTypeRaw string
	args = f.Args()
	positionalArgsCount := len(args)
	switch positionalArgsCount {
	case 0, 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2, got %d)", positionalArgsCount))
		return 1
	case 2:
		pluginTypeRaw = args[0]
		pluginNameRaw = args[1]
	default:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 2, got %d)", positionalArgsCount))
		return 1
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
	if c.flagPluginVersion != "" {
		_, err := semver.NewSemver(c.flagPluginVersion)
		if err != nil {
			c.UI.Error(fmt.Sprintf("version %q is not a valid semantic version: %v", c.flagPluginVersion, err))
			return 2
		}
	}

	// The deregister endpoint returns 200 if the plugin doesn't exist, so first
	// try fetching the plugin to help improve info printed to the user.
	// 404 / 400 => Return early with a descriptive message.
	// Auto-selected warning present => versioned-only plugin; no unversioned entry
	//                                  to deregister, so treat as "does not exist".
	// Other error => Continue attempting to deregister the plugin anyway.
	// Plugin exists but is builtin => Error early.
	// Otherwise => If deregister succeeds, we can report that the plugin really
	//              was deregistered (and not just already absent).
	var pluginExists bool
	// Build the catalog path for the raw Logical read so we get back Warnings,
	// which the typed GetPlugin API strips out.
	catalogPath := fmt.Sprintf("sys/plugins/catalog/%s/%s", pluginType, pluginName)
	req := client.NewRequest(http.MethodGet, "/v1/"+catalogPath)
	if c.flagPluginVersion != "" {
		req.Params.Set("version", c.flagPluginVersion)
	}
	if secret, err := client.RawRequestWithContext(context.Background(), req); err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && (respErr.StatusCode == http.StatusNotFound || respErr.StatusCode == http.StatusBadRequest) {
			// Surface the server's error message directly — it already explains
			// whether the plugin is missing entirely or multiple versions exist.
			if len(respErr.Errors) > 0 {
				c.UI.Error(respErr.Errors[0])
			} else {
				c.UI.Output(fmt.Sprintf("Plugin %q (type: %q, version %q) does not exist in the catalog", pluginName, pluginType, c.flagPluginVersion))
			}
			return 0
		}
		// Best-effort check, continue trying to deregister.
	} else {
		defer secret.Body.Close()
		var result api.Secret
		if err := secret.DecodeJSON(&result); err == nil {
			// If the server auto-selected a version because none was specified,
			// there is no unversioned entry to deregister. Tell the user which
			// version they need to specify explicitly.
			for _, w := range result.Warnings {
				if strings.Contains(w, "no version was specified") {
					// Use the version from the response data — more reliable
					// than parsing the warning string.
					autoVersion, _ := result.Data["version"].(string)
					if autoVersion != "" {
						c.UI.Output(fmt.Sprintf("Plugin %q (type: %q, version %q) does not exist in the catalog as an unversioned entry; did you mean -version=%s?", pluginName, pluginType, autoVersion, autoVersion))
					} else {
						c.UI.Output(fmt.Sprintf("Plugin %q (type: %q) does not exist in the catalog as an unversioned entry; use -version to specify a version", pluginName, pluginType))
					}
					return 0
				}
			}
			if builtin, _ := result.Data["builtin"].(bool); builtin {
				c.UI.Error(fmt.Sprintf("Plugin %q (type: %q) is a builtin plugin and cannot be deregistered", pluginName, pluginType))
				return 2
			}
			pluginExists = true
		}
	}

	if err := client.Sys().DeregisterPluginWithContext(context.Background(), &api.DeregisterPluginInput{
		Name:    pluginName,
		Type:    pluginType,
		Version: c.flagPluginVersion,
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error deregistering plugin named %s: %s", pluginName, err))
		return 2
	}

	if pluginExists {
		c.UI.Output(fmt.Sprintf("Success! Deregistered %s plugin: %s", pluginType, pluginName))
	} else {
		c.UI.Output(fmt.Sprintf("Success! Deregistered %s plugin (if it was registered): %s", pluginType, pluginName))
	}
	return 0
}
