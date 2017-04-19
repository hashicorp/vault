package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// PathHelpCommand is a Command that lists the mounts.
type PathHelpCommand struct {
	meta.Meta
}

func (c *PathHelpCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("help", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		c.Ui.Error("\nhelp expects a single argument")
		return 1
	}

	path := args[0]

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	help, err := client.Help(path)
	if err != nil {
		if strings.Contains(err.Error(), "Vault is sealed") {
			c.Ui.Error(`Error: Vault is sealed.

The path-help command requires the vault to be unsealed so that
mount points of secret backends are known.`)
		} else {
			c.Ui.Error(fmt.Sprintf(
				"Error reading help: %s", err))
		}
		return 1
	}

	c.Ui.Output(help.Help)
	return 0
}

func (c *PathHelpCommand) Synopsis() string {
	return "Look up the help for a path"
}

func (c *PathHelpCommand) Help() string {
	helpText := `
Usage: vault path-help [options] path

  Look up the help for a path.

  All endpoints in Vault from system paths, secret paths, and credential
  providers provide built-in help. This command looks up and outputs that
  help.

  The command requires that the vault be unsealed, because otherwise
  the mount points of the backends are unknown.

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
