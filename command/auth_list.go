package command

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*AuthListCommand)(nil)
var _ cli.CommandAutocomplete = (*AuthListCommand)(nil)

// AuthListCommand is a Command that lists the enabled authentication methods
// and data about them.
type AuthListCommand struct {
	*BaseCommand

	flagDetailed bool
}

func (c *AuthListCommand) Synopsis() string {
	return "Lists enabled auth providers"
}

func (c *AuthListCommand) Help() string {
	helpText := `
Usage: vault auth-methods [options]

  Lists the enabled authentication providers on the Vault server. This command
  also outputs information about the provider including configuration and
  human-friendly descriptions. A TTL of "system" indicates that the system
  default is in use.

  List all enabled authentication providers:

      $ vault auth-list

  List all enabled authentication providers with detailed output:

      $ vault auth-list -detailed

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuthListCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "detailed",
		Target:  &c.flagDetailed,
		Default: false,
		Usage: "Print detailed information such as configuration and replication " +
			"status about each authentication provider.",
	})

	return set
}

func (c *AuthListCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *AuthListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthListCommand) Run(args []string) int {
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

	auths, err := client.Sys().ListAuth()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing enabled authentications: %s", err))
		return 2
	}

	if c.flagDetailed {
		c.UI.Output(tableOutput(c.detailedMounts(auths)))
		return 0
	}

	c.UI.Output(tableOutput(c.simpleMounts(auths)))
	return 0
}

func (c *AuthListCommand) simpleMounts(auths map[string]*api.AuthMount) []string {
	paths := make([]string, 0, len(auths))
	for path := range auths {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	out := []string{"Path | Type | Description"}
	for _, path := range paths {
		mount := auths[path]
		out = append(out, fmt.Sprintf("%s | %s | %s", path, mount.Type, mount.Description))
	}

	return out
}

func (c *AuthListCommand) detailedMounts(auths map[string]*api.AuthMount) []string {
	paths := make([]string, 0, len(auths))
	for path := range auths {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	calcTTL := func(typ string, ttl int) string {
		switch {
		case typ == "system", typ == "cubbyhole":
			return ""
		case ttl != 0:
			return strconv.Itoa(ttl)
		default:
			return "system"
		}
	}

	out := []string{"Path | Type | Accessor | Plugin | Default TTL | Max TTL | Replication | Description"}
	for _, path := range paths {
		mount := auths[path]

		defaultTTL := calcTTL(mount.Type, mount.Config.DefaultLeaseTTL)
		maxTTL := calcTTL(mount.Type, mount.Config.MaxLeaseTTL)

		replication := "replicated"
		if mount.Local {
			replication = "local"
		}

		out = append(out, fmt.Sprintf("%s | %s | %s | %s | %s | %s | %v | %s",
			path,
			mount.Type,
			mount.Accessor,
			mount.Config.PluginName,
			defaultTTL,
			maxTTL,
			replication,
			mount.Description,
		))
	}

	return out
}
