package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// UnmountCommand is a Command that mounts a new mount.
type UnmountCommand struct {
	meta.Meta
}

func (c *UnmountCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("mount", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\nunmount expects one argument: the path to unmount"))
		return 1
	}

	path := args[0]

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	if err := client.Sys().Unmount(path); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Unmount error: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf(
		"Successfully unmounted '%s' if it was mounted", path))

	return 0
}

func (c *UnmountCommand) Synopsis() string {
	return "Unmount a secret backend"
}

func (c *UnmountCommand) Help() string {
	helpText := `
Usage: vault unmount [options] path

  Unmount a secret backend.

  This command unmounts a secret backend. All the secrets created
  by this backend will be revoked and its Vault data will be deleted.

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
