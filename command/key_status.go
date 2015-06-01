package command

import (
	"fmt"
	"strings"
)

// KeyStatusCommand is a Command that provides information about the key status
type KeyStatusCommand struct {
	Meta
}

func (c *KeyStatusCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("key-status", FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	status, err := client.Sys().KeyStatus()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error reading audits: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf("Key Term: %d", status.Term))
	c.Ui.Output(fmt.Sprintf("Installation Time: %v", status.InstallTime))
	return 0
}

func (c *KeyStatusCommand) Synopsis() string {
	return "Provides information about the active encryption key"
}

func (c *KeyStatusCommand) Help() string {
	helpText := `
Usage: vault key-status [options]

  Provides information about the active encryption key. Specifically,
  the current key term and the key installation time.

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
