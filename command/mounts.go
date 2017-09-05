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
var _ cli.Command = (*MountsCommand)(nil)
var _ cli.CommandAutocomplete = (*MountsCommand)(nil)

// MountsCommand is a Command that lists the mounts.
type MountsCommand struct {
	*BaseCommand

	flagDetailed bool
}

func (c *MountsCommand) Synopsis() string {
	return "Lists mounted secret backends"
}

func (c *MountsCommand) Help() string {
	helpText := `
Usage: vault mounts [options]

  Lists the mounted secret backends on the Vault server. This command also
  outputs information about the mount point including configured TTLs and
  human-friendly descriptions. A TTL of "system" indicates that the system
  default is in use.

  List all mounts:

      $ vault mounts

  List all mounts with detailed output:

      $ vault mounts -detailed

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *MountsCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "detailed",
		Target:  &c.flagDetailed,
		Default: false,
		Usage: "Print detailed information such as TTLs and replication status " +
			"about each mount.",
	})

	return set
}

func (c *MountsCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *MountsCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *MountsCommand) Run(args []string) int {
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
		c.UI.Error(fmt.Sprintf("Error listing mounts: %s", err))
		return 2
	}

	if c.flagDetailed {
		c.UI.Output(tableOutput(c.detailedMounts(mounts)))
		return 0
	}

	c.UI.Output(tableOutput(c.simpleMounts(mounts)))
	return 0
}

func (c *MountsCommand) simpleMounts(mounts map[string]*api.MountOutput) []string {
	paths := make([]string, 0, len(mounts))
	for path := range mounts {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	out := []string{"Path | Type | Description"}
	for _, path := range paths {
		mount := mounts[path]
		out = append(out, fmt.Sprintf("%s | %s | %s", path, mount.Type, mount.Description))
	}

	return out
}

func (c *MountsCommand) detailedMounts(mounts map[string]*api.MountOutput) []string {
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

	out := []string{"Path | Type | Accessor | Plugin | Default TTL | Max TTL | Force No Cache | Replication | Description"}
	for _, path := range paths {
		mount := mounts[path]

		defaultTTL := calcTTL(mount.Type, mount.Config.DefaultLeaseTTL)
		maxTTL := calcTTL(mount.Type, mount.Config.MaxLeaseTTL)

		replication := "replicated"
		if mount.Local {
			replication = "local"
		}

		out = append(out, fmt.Sprintf("%s | %s | %s | %s | %s | %s | %v | %s | %s",
			path,
			mount.Type,
			mount.Accessor,
			mount.Config.PluginName,
			defaultTTL,
			maxTTL,
			mount.Config.ForceNoCache,
			replication,
			mount.Description,
		))
	}

	return out
}
