package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// AuthDisableCommand is a Command that enables a new endpoint.
type AuthDisableCommand struct {
	meta.Meta
}

func (c *AuthDisableCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("auth-disable", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\nauth-disable expects one argument: the path to disable."))
		return 1
	}

	path := args[0]

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	if err := client.Sys().DisableAuth(path); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf(
		"Disabled auth provider at path '%s' if it was enabled", path))

	return 0
}

func (c *AuthDisableCommand) Synopsis() string {
	return "Disable an auth provider"
}

func (c *AuthDisableCommand) Help() string {
	helpText := `
Usage: vault auth-disable [options] path

  Disable an already-enabled auth provider.

  Once the auth provider is disabled its path can no longer be used
  to authenticate. All access tokens generated via the disabled auth provider
  will be revoked. This command will block until all tokens are revoked.
  If the command is exited early the tokens will still be revoked.

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
