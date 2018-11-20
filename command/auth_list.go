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

var _ cli.Command = (*AuthListCommand)(nil)
var _ cli.CommandAutocomplete = (*AuthListCommand)(nil)

type AuthListCommand struct {
	*BaseCommand

	flagDetailed bool
}

func (c *AuthListCommand) Synopsis() string {
	return "Lists enabled auth methods"
}

func (c *AuthListCommand) Help() string {
	helpText := `
Usage: vault auth list [options]

  Lists the enabled auth methods on the Vault server. This command also outputs
  information about the method including configuration and human-friendly
  descriptions. A TTL of "system" indicates that the system default is in use.

  List all enabled auth methods:

      $ vault auth list

  List all enabled auth methods with detailed output:

      $ vault auth list -detailed

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuthListCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "detailed",
		Target:  &c.flagDetailed,
		Default: false,
		Usage: "Print detailed information such as configuration and replication " +
			"status about each auth method. This option is only applicable to " +
			"table-formatted output.",
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

	switch Format(c.UI) {
	case "table":
		if c.flagDetailed {
			c.UI.Output(tableOutput(c.detailedMounts(auths), nil))
			return 0
		}
		c.UI.Output(tableOutput(c.simpleMounts(auths), nil))
		return 0
	default:
		return OutputData(c.UI, auths)
	}
}

func (c *AuthListCommand) simpleMounts(auths map[string]*api.AuthMount) []string {
	paths := make([]string, 0, len(auths))
	for path := range auths {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	out := []string{"Path | Type | Accessor | Description"}
	for _, path := range paths {
		mount := auths[path]
		out = append(out, fmt.Sprintf("%s | %s | %s | %s", path, mount.Type, mount.Accessor, mount.Description))
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

	out := []string{"Path | Plugin | Accessor | Default TTL | Max TTL | Token Type | Replication | Seal Wrap | Options | Description"}
	for _, path := range paths {
		mount := auths[path]

		defaultTTL := calcTTL(mount.Type, mount.Config.DefaultLeaseTTL)
		maxTTL := calcTTL(mount.Type, mount.Config.MaxLeaseTTL)

		replication := "replicated"
		if mount.Local {
			replication = "local"
		}

		pluginName := mount.Type
		if pluginName == "plugin" {
			pluginName = mount.Config.PluginName
		}

		out = append(out, fmt.Sprintf("%s | %s | %s | %s | %s | %s | %s | %t | %v | %s",
			path,
			pluginName,
			mount.Accessor,
			defaultTTL,
			maxTTL,
			mount.Config.TokenType,
			replication,
			mount.SealWrap,
			mount.Options,
			mount.Description,
		))
	}

	return out
}
