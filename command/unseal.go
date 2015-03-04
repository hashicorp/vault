package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/helper/password"
)

// UnsealCommand is a Command that unseals the vault.
type UnsealCommand struct {
	Meta
}

func (c *UnsealCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("unseal", FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	value, err := password.Read(os.Stdin)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error attempting to ask for password. The raw error message\n"+
				"is shown below, but the most common reason for this error is\n"+
				"that you attempted to pipe a value into unseal. This is not\n"+
				"allowed. The value must be typed directly into the command\n"+
				"after it is executed.\n\n"+
				"Raw error: %s", err))
		return 1
	}

	c.Ui.Output(value)
	return 0
}

func (c *UnsealCommand) Synopsis() string {
	return "Unseals the vault server"
}

func (c *UnsealCommand) Help() string {
	helpText := `
Usage: vault unseal [options]

  Unseal the vault by entering a portion of the master key. Once all
  portions are entered, the Vault will be unsealed.

  Every Vault server initially starts as sealed. It cannot perform any
  operation except unsealing until it is sealed. Secrets cannot be accessed
  in any way until the vault is unsealed. This command allows you to enter
  a portion of the master key to unseal the vault.

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
