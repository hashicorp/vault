package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/meta"
)

// MountCommand is a Command that mounts a new mount.
type MountCommand struct {
	meta.Meta
}

func (c *MountCommand) Run(args []string) int {
	var description, path, defaultLeaseTTL, maxLeaseTTL string
	var local, forceNoCache bool
	flags := c.Meta.FlagSet("mount", meta.FlagSetDefault)
	flags.StringVar(&description, "description", "", "")
	flags.StringVar(&path, "path", "", "")
	flags.StringVar(&defaultLeaseTTL, "default-lease-ttl", "", "")
	flags.StringVar(&maxLeaseTTL, "max-lease-ttl", "", "")
	flags.BoolVar(&forceNoCache, "force-no-cache", false, "")
	flags.BoolVar(&local, "local", false, "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\nmount expects one argument: the type to mount."))
		return 1
	}

	mountType := args[0]

	// If no path is specified, we default the path to the backend type
	if path == "" {
		path = mountType
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	mountInfo := &api.MountInput{
		Type:        mountType,
		Description: description,
		Config: api.MountConfigInput{
			DefaultLeaseTTL: defaultLeaseTTL,
			MaxLeaseTTL:     maxLeaseTTL,
			ForceNoCache:    forceNoCache,
		},
		Local: local,
	}

	if err := client.Sys().Mount(path, mountInfo); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Mount error: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf(
		"Successfully mounted '%s' at '%s'!",
		mountType, path))

	return 0
}

func (c *MountCommand) Synopsis() string {
	return "Mount a logical backend"
}

func (c *MountCommand) Help() string {
	helpText := `
Usage: vault mount [options] type

  Mount a logical backend.

  This command mounts a logical backend for storing and/or generating
  secrets.

General Options:
` + meta.GeneralOptionsUsage() + `
Mount Options:

  -description=<desc>            Human-friendly description of the purpose for
                                 the mount. This shows up in the mounts command.

  -path=<path>                   Mount point for the logical backend. This
                                 defaults to the type of the mount.

  -default-lease-ttl=<duration>  Default lease time-to-live for this backend.
                                 If not specified, uses the global default, or
                                 the previously set value. Set to '0' to
                                 explicitly set it to use the global default.

  -max-lease-ttl=<duration>      Max lease time-to-live for this backend.
                                 If not specified, uses the global default, or
                                 the previously set value. Set to '0' to
                                 explicitly set it to use the global default.

  -force-no-cache                Forces the backend to disable caching. If not
                                 specified, uses the global default. This does
                                 not affect caching of the underlying encrypted
                                 data storage.

  -local                         Mark the mount as a local mount. Local mounts
                                 are not replicated nor (if a secondary)
                                 removed by replication.

`
	return strings.TrimSpace(helpText)
}
