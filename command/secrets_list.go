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

var _ cli.Command = (*SecretsListCommand)(nil)
var _ cli.CommandAutocomplete = (*SecretsListCommand)(nil)

type SecretsListCommand struct {
	*BaseCommand

	flagDetailed bool
}

func (c *SecretsListCommand) Synopsis() string {
	return "List enabled secrets engines"
}

func (c *SecretsListCommand) Help() string {
	helpText := `
Usage: vault secrets list [options]

  Lists the enabled secret engines on the Vault server. This command also
  outputs information about the enabled path including configured TTLs and
  human-friendly descriptions. A TTL of "system" indicates that the system
  default is in use.

  List all enabled secrets engines:

      $ vault secrets list

  List all enabled secrets engines with detailed output:

      $ vault secrets list -detailed

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *SecretsListCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "detailed",
		Target:  &c.flagDetailed,
		Default: false,
		Usage: "Print detailed information such as TTLs and replication status " +
			"about each secrets engine.",
	})

	return set
}

func (c *SecretsListCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *SecretsListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *SecretsListCommand) Run(args []string) int {
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

	mounts, err := client.Sys().ListMounts()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing secrets engines: %s", err))
		return 2
	}

	switch Format(c.UI) {
	case "table":
		if c.flagDetailed {
			c.UI.Output(tableOutput(c.detailedMounts(mounts), nil))
			return 0
		}
		c.UI.Output(tableOutput(c.simpleMounts(mounts), nil))
		return 0
	default:
		return OutputData(c.UI, mounts)
	}
}

func (c *SecretsListCommand) simpleMounts(mounts map[string]*api.MountOutput) []string {
	paths := make([]string, 0, len(mounts))
	for path := range mounts {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	out := []string{"Path | Type | Accessor | Description"}
	for _, path := range paths {
		mount := mounts[path]
		out = append(out, fmt.Sprintf("%s | %s | %s | %s", path, mount.Type, mount.Accessor, mount.Description))
	}

	return out
}

func (c *SecretsListCommand) detailedMounts(mounts map[string]*api.MountOutput) []string {
	paths := make([]string, 0, len(mounts))
	for path := range mounts {
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

	out := []string{"Path | Plugin | Accessor | Default TTL | Max TTL | Force No Cache | Replication | Seal Wrap | Options | Description"}
	for _, path := range paths {
		mount := mounts[path]

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

		out = append(out, fmt.Sprintf("%s | %s | %s | %s | %s | %t | %s | %t | %v | %s",
			path,
			pluginName,
			mount.Accessor,
			defaultTTL,
			maxTTL,
			mount.Config.ForceNoCache,
			replication,
			mount.SealWrap,
			mount.Options,
			mount.Description,
		))
	}

	return out
}
