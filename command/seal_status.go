package command

import (
	"fmt"
	"strings"
)

// SealStatusCommand is a Command that outputs the status of whether
// Vault is sealed or not.
type SealStatusCommand struct {
	Meta
}

func (c *SealStatusCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("seal-status", FlagSetDefault)
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

	status, err := client.Sys().SealStatus()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error checking seal status: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf(
		"Sealed: %v\n"+
			"Key Shares: %d\n"+
			"Key Threshold: %d\n"+
			"Unseal Progress: %d",
		status.Sealed,
		status.N,
		status.T,
		status.Progress,
	))

	if status.Sealed {
		return 1
	} else {
		return 0
	}
}

func (c *SealStatusCommand) Synopsis() string {
	return "Outputs status of whether Vault is sealed"
}

func (c *SealStatusCommand) Help() string {
	helpText := `
Usage: vault seal-status [options]

  Outputs the state of the Vault, sealed or unsealed.

  This command outputs whether or not the Vault is sealed. The exit
  code also reflects the seal status (0 unsealed, 1 sealed, 2+ error).

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
