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
	flags.StringVar(&defaultLeaseTTL, "default_lease_ttl", "", "")
	flags.StringVar(&maxLeaseTTL, "max_lease_ttl", "", "")
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

	mountConfig := &vault.MountConfig{}
	var err error
	var passConfig bool
	if defaultLeaseTTL != "" {
		mountConfig.DefaultLeaseTTL, err = time.ParseDuration(defaultLeaseTTL)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error parsing default lease TTL duration: %s", err))
			return 1
		}
		passConfig = true
	}
	if maxLeaseTTL != "" {
		mountConfig.MaxLeaseTTL, err = time.ParseDuration(maxLeaseTTL)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error parsing max lease TTL duration: %s", err))
			return 1
		}
		passConfig = true
	}

	if !passConfig {
		mountConfig = nil
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

  Example: vault remount secret/ generic/

General Options:

  ` + generalOptionsUsage()
	return strings.TrimSpace(helpText)
}
