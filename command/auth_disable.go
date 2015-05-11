package command

import (
	"fmt"
	"strings"
)

// AuthDisableCommand is a Command that enables a new endpoint.
type AuthDisableCommand struct {
	Meta
}

func (c *AuthDisableCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("auth-disable", FlagSetDefault)
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
		"Disabled auth provider at path '%s'!", path))

	return 0
}

func (c *AuthDisableCommand) Synopsis() string {
	return "Disable an auth provider"
}

func (c *AuthDisableCommand) Help() string {
	helpText := `
Usage: vault auth-disable [options] path

  Disable an already-enabled auth provider.

  Once the auth provider is disabled, that path cannot be used anymore
  to authenticate. All access tokens generated via the disabled auth provider
  will be revoked. This command will block until all tokens are revoked.
  If the command is exited early, the tokens will still be revoked.

General Options:

  -address=addr           The address of the Vault server.

  -ca-cert=path           Path to a PEM encoded CA cert file to use to
                          verify the Vault server SSL certificate.

  -ca-path=path           Path to a directory of PEM encoded CA cert files
                          to verify the Vault server SSL certificate. If both
                          -ca-cert and -ca-path are specified, -ca-path is used.

  -tls-skip-verify        Do not verify TLS certificate. This is highly
                          not recommended. This is especially not recommended
                          for unsealing a vault.

`
	return strings.TrimSpace(helpText)
}
