package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// AuthEnableCommand is a Command that enables a new endpoint.
type AuthEnableCommand struct {
	meta.Meta
}

func (c *AuthEnableCommand) Run(args []string) int {
	var description, path string
	flags := c.Meta.FlagSet("auth-enable", meta.FlagSetDefault)
	flags.StringVar(&description, "description", "", "")
	flags.StringVar(&path, "path", "", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\nauth-enable expects one argument: the type to enable."))
		return 1
	}

	authType := args[0]

	// If no path is specified, we default the path to the backend type
	if path == "" {
		path = authType
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	if err := client.Sys().EnableAuth(path, authType, description); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf(
		"Successfully enabled '%s' at '%s'!",
		authType, path))

	return 0
}

func (c *AuthEnableCommand) Synopsis() string {
	return "Enable a new auth provider"
}

func (c *AuthEnableCommand) Help() string {
	helpText := `
Usage: vault auth-enable [options] type

  Enable a new auth provider.

  This command enables a new auth provider. An auth provider is responsible
  for authenticating a user and assigning them policies with which they can
  access Vault.

General Options:
` + meta.GeneralOptionsUsage() + `
Auth Enable Options:

  -description=<desc>     Human-friendly description of the purpose for the
                          auth provider. This shows up in the auth -methods command.

  -path=<path>            Mount point for the auth provider. This defaults
                          to the type of the mount. This will make the auth
                          provider available at "/auth/<path>"

`
	return strings.TrimSpace(helpText)
}
