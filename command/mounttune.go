package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/vault"
)

// MountTuneCommand is a Command that remounts a mounted secret backend
// to a new endpoint.
type MountTuneCommand struct {
	Meta
}

func (c *MountTuneCommand) Run(args []string) int {
	var defaultLeaseTTL, maxLeaseTTL string
	flags := c.Meta.FlagSet("mount-tune", FlagSetDefault)
	flags.StringVar(&defaultLeaseTTL, "default-lease-ttl", "", "")
	flags.StringVar(&maxLeaseTTL, "max-lease-ttl", "", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\n'mount-tune' expects one arguments: the mount path"))
		return 1
	}

	path := args[0]

	mountConfig := vault.MountConfig{}
	if defaultLeaseTTL != "" {
		defTTL, err := time.ParseDuration(defaultLeaseTTL)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error parsing default lease TTL duration: %s", err))
			return 2
		}
		mountConfig.DefaultLeaseTTL = &defTTL
	}
	if maxLeaseTTL != "" {
		maxTTL, err := time.ParseDuration(maxLeaseTTL)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error parsing max lease TTL duration: %s", err))
			return 2
		}
		mountConfig.MaxLeaseTTL = &maxTTL
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	if err := client.Sys().TuneMount(path, mountConfig); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Mount tune error: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf(
		"Successfully tuned mount '%s'!", path))

	return 0
}

func (c *MountTuneCommand) Synopsis() string {
	return "Tune mount configuration parameters"
}

func (c *MountTuneCommand) Help() string {
	helpText := `
  Usage: vault mount-tune [options] path

  Tune configuration options for a mounted secret backend.

  Example: vault tune-mount -default-lease-ttl="24h" secret/

General Options:

  ` + generalOptionsUsage() + `

Mount Options:

  -default-lease-ttl=<duration>  Default lease time-to-live for this backend.
                                 If not specified, uses the global default, or
                                 the previously set value. Set to '0' to
                                 explicitly set it to use the global default.

  -max-lease-ttl=<duration>      Max lease time-to-live for this backend.
                                 If not specified, uses the global default, or
                                 the previously set value. Set to '0' to
                                 explicitly set it to use the global default.

`
	return strings.TrimSpace(helpText)
}
