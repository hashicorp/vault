package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/vault"
)

// RemountCommand is a Command that remounts a mounted secret backend
// to a new endpoint.
type RemountCommand struct {
	Meta
}

func (c *RemountCommand) Run(args []string) int {
	var defaultLeaseTTL, maxLeaseTTL string
	flags := c.Meta.FlagSet("remount", FlagSetDefault)
	flags.StringVar(&defaultLeaseTTL, "default-lease-ttl", "", "")
	flags.StringVar(&maxLeaseTTL, "max-lease-ttl", "", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 2 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\nRemount expects two arguments: the from and to path"))
		return 1
	}

	from := args[0]
	to := args[1]

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

	if err := client.Sys().Remount(from, to, mountConfig); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Unmount error: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf(
		"Successfully remounted from '%s' to '%s'!", from, to))

	return 0
}

func (c *RemountCommand) Synopsis() string {
	return "Remount a secret backend to a new path"
}

func (c *RemountCommand) Help() string {
	helpText := `
Usage: vault remount [options] from to

  Remount a mounted secret backend to a new path.

  This command remounts a secret backend that is already mounted to
  a new path. All the secrets from the old path will be revoked, but
  the Vault data associated with the backend will be preserved (such
  as configuration data).

  If the 'from' and 'to' values of the same, performs an in-place
  remount. This allows you to change mount options.

  Example: vault remount secret/ generic/

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
