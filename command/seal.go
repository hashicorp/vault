package command

import (
	"strings"
)

// SealCommand is a Command that seals the vault.
type SealCommand struct {
	Meta
}

func (c *SealCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("unseal", FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	return 0
}

func (c *SealCommand) Synopsis() string {
	return "Seals the vault server"
}

func (c *SealCommand) Help() string {
	helpText := `
Usage: vault seal [options]

  Seal the vault.

  Sealing a vault tells the Vault server to stop responding to any
  access operations until it is unsealed again. A sealed vault throws away
  its master key to unlock the data, so it physically is blocked from
  responding to operations again until the Vault is unsealed again with
  the "unseal" command or via the API.

  This command is idempotent, if the vault is already sealed it does nothing.

  If an unseal has started, sealing the vault will reset the unsealing
  process. You'll have to re-enter every portion of the master key again.
  This is the same as running "vault unseal -reset".

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
