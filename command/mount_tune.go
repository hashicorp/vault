package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/meta"
)

// MountTuneCommand is a Command that remounts a mounted secret backend
// to a new endpoint.
type MountTuneCommand struct {
	meta.Meta
}

func (c *MountTuneCommand) Run(args []string) int {
	var defaultLeaseTTL, maxLeaseTTL string
	flags := c.Meta.FlagSet("mount-tune", meta.FlagSetDefault)
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
			"\nmount-tune expects one arguments: the mount path"))
		return 1
	}

	path := args[0]

	mountConfig := api.MountConfigInput{
		DefaultLeaseTTL: defaultLeaseTTL,
		MaxLeaseTTL:     maxLeaseTTL,
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

  Example: vault mount-tune -default-lease-ttl="24h" secret

General Options:
` + meta.GeneralOptionsUsage() + `
Mount Options:

  -default-lease-ttl=<duration>  Default lease time-to-live for this backend.
                                 If not specified, uses the system default, or
                                 the previously set value. Set to 'system' to
                                 explicitly set it to use the system default.

  -max-lease-ttl=<duration>      Max lease time-to-live for this backend.
                                 If not specified, uses the system default, or
                                 the previously set value. Set to 'system' to
                                 explicitly set it to use the system default.

`
	return strings.TrimSpace(helpText)
}
