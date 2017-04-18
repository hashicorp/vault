package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/helper/password"
	"github.com/hashicorp/vault/meta"
)

// UnsealCommand is a Command that unseals the vault.
type UnsealCommand struct {
	meta.Meta

	// Key can be used to pre-seed the key. If it is set, it will not
	// be asked with the `password` helper.
	Key string
}

func (c *UnsealCommand) Run(args []string) int {
	var reset bool
	flags := c.Meta.FlagSet("unseal", meta.FlagSetDefault)
	flags.BoolVar(&reset, "reset", false, "")
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

	sealStatus, err := client.Sys().SealStatus()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error checking seal status: %s", err))
		return 2
	}

	if !sealStatus.Sealed {
		c.Ui.Output("Vault is already unsealed.")
		return 0
	}

	args = flags.Args()
	if reset {
		sealStatus, err = client.Sys().ResetUnsealProcess()
	} else {
		value := c.Key
		if len(args) > 0 {
			value = args[0]
		}
		if value == "" {
			fmt.Printf("Key (will be hidden): ")
			value, err = password.Read(os.Stdin)
			fmt.Printf("\n")
			if err != nil {
				c.Ui.Error(fmt.Sprintf(
					"Error attempting to ask for password. The raw error message\n"+
						"is shown below, but the most common reason for this error is\n"+
						"that you attempted to pipe a value into unseal or you're\n"+
						"executing `vault unseal` from outside of a terminal.\n\n"+
						"You should use `vault unseal` from a terminal for maximum\n"+
						"security. If this isn't an option, the unseal key can be passed\n"+
						"in using the first parameter.\n\n"+
						"Raw error: %s", err))
				return 1
			}
		}
		sealStatus, err = client.Sys().Unseal(strings.TrimSpace(value))
	}

	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf(
		"Sealed: %v\n"+
			"Key Shares: %d\n"+
			"Key Threshold: %d\n"+
			"Unseal Progress: %d\n"+
			"Unseal Nonce: %v",
		sealStatus.Sealed,
		sealStatus.N,
		sealStatus.T,
		sealStatus.Progress,
		sealStatus.Nonce,
	))

	return 0
}

func (c *UnsealCommand) Synopsis() string {
	return "Unseals the Vault server"
}

func (c *UnsealCommand) Help() string {
	helpText := `
Usage: vault unseal [options] [key]

  Unseal the vault by entering a portion of the master key. Once all
  portions are entered, the vault will be unsealed.

  Every Vault server initially starts as sealed. It cannot perform any
  operation except unsealing until it is sealed. Secrets cannot be accessed
  in any way until the vault is unsealed. This command allows you to enter
  a portion of the master key to unseal the vault.

  The unseal key can be specified via the command line, but this is
  not recommended. The key may then live in your terminal history. This
  only exists to assist in scripting.

General Options:
` + meta.GeneralOptionsUsage() + `
Unseal Options:

  -reset                  Reset the unsealing process by throwing away
                          prior keys in process to unseal the vault.

`
	return strings.TrimSpace(helpText)
}
