package command

import (
	"fmt"
	"strings"
)

// HelpCommand is a Command that lists the mounts.
type HelpCommand struct {
	Meta
}

func (c *HelpCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("help", FlagSetDefault)
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
		c.Ui.Error(fmt.Sprintf(
			"Error reading help: %s", err))
		return 1
	}

	c.Ui.Output(help.Help)
	return 0
}

func (c *HelpCommand) Synopsis() string {
	return "Look up the help for a path"
}

func (c *HelpCommand) Help() string {
	helpText := `
Usage: vault help [options] path

  Look up the help for a path.

  All endpoints in Vault from system paths, secret paths, and credential
  providers provide built-in help. This command looks up and outputs that
  help.

General Options:

  -address=TODO           The address of the Vault server.

  -ca-cert=path           Path to a PEM encoded CA cert file to use to
                          verify the Vault server SSL certificate.

  -ca-path=path           Path to a directory of PEM encoded CA cert files
                          to verify the Vault server SSL certificate. If both
                          -ca-cert and -ca-path are specified, -ca-path is used.

  -insecure               Do not verify TLS certificate. This is highly
                          not recommended. This is especially not recommended
                          for unsealing a vault.

`
	return strings.TrimSpace(helpText)
}
