package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/password"
)

// RekeyCommand is a Command that rekeys the vault.
type RekeyCommand struct {
	Meta

	// Key can be used to pre-seed the key. If it is set, it will not
	// be asked with the `password` helper.
	Key string
}

func (c *RekeyCommand) Run(args []string) int {
	var init, cancel, status bool
	var shares, threshold int
	flags := c.Meta.FlagSet("rekey", FlagSetDefault)
	flags.BoolVar(&init, "init", false, "")
	flags.BoolVar(&cancel, "cancel", false, "")
	flags.BoolVar(&status, "status", false, "")
	flags.IntVar(&shares, "key-shares", 5, "")
	flags.IntVar(&threshold, "key-threshold", 3, "")
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

	// Check if we are running doing any restricted variants
	if init {
		return c.initRekey(client, shares, threshold)
	} else if cancel {
		return c.cancelRekey(client)
	} else if status {
		return c.rekeyStatus(client)
	}

	// Check if the rekey is started
	rekeyStatus, err := client.Sys().RekeyStatus()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading rekey status: %s", err))
		return 1
	}

	// Start the rekey process if not started
	if !rekeyStatus.Started {
		err := client.Sys().RekeyInit(&api.RekeyInitRequest{
			SecretShares:    shares,
			SecretThreshold: threshold,
		})
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error initializing rekey: %s", err))
			return 1
		}
	} else {
		shares = rekeyStatus.N
		threshold = rekeyStatus.T
		c.Ui.Output(fmt.Sprintf(
			"Rekey already in progress\n"+
				"Key Shares: %d\n"+
				"Key Threshold: %d\n",
			shares,
			threshold,
		))
	}

	// Get the unseal key
	args = flags.Args()
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
					"executing `vault rekey` from outside of a terminal.\n\n"+
					"You should use `vault rekey` from a terminal for maximum\n"+
					"security. If this isn't an option, the unseal key can be passed\n"+
					"in using the first parameter.\n\n"+
					"Raw error: %s", err))
			return 1
		}
	}

	// Provide the key, this may potentially complete the update
	result, err := client.Sys().RekeyUpdate(strings.TrimSpace(value))
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error attempting rekey update: %s", err))
		return 1
	}

	// If we are not complete, then dump the status
	if !result.Complete {
		return c.rekeyStatus(client)
	}

	// Provide the keys
	for i, key := range result.Keys {
		c.Ui.Output(fmt.Sprintf("Key %d: %s", i+1, key))
	}

	c.Ui.Output(fmt.Sprintf(
		"\n"+
			"Vault rekeyed with %d keys and a key threshold of %d. Please\n"+
			"securely distribute the above keys. When the Vault is re-sealed,\n"+
			"restarted, or stopped, you must provide at least %d of these keys\n"+
			"to unseal it again.\n\n"+
			"Vault does not store the master key. Without at least %d keys,\n"+
			"your Vault will remain permanently sealed.",
		shares,
		threshold,
		threshold,
		threshold,
	))

	return 0
}

// initRekey is used to start the rekey process
func (c *RekeyCommand) initRekey(client *api.Client, shares, threshold int) int {
	// Start the rekey
	err := client.Sys().RekeyInit(&api.RekeyInitRequest{
		SecretShares:    shares,
		SecretThreshold: threshold,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing rekey: %s", err))
		return 1
	}

	// Provide the current status
	return c.rekeyStatus(client)
}

// cancelRekey is used to abort the rekey process
func (c *RekeyCommand) cancelRekey(client *api.Client) int {
	err := client.Sys().RekeyCancel()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to cancel rekey: %s", err))
		return 1
	}
	c.Ui.Output("Rekey canceled.")
	return 0
}

// rekeyStatus is used just to fetch and dump the statu
func (c *RekeyCommand) rekeyStatus(client *api.Client) int {
	// Check the status
	status, err := client.Sys().RekeyStatus()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading rekey status: %s", err))
		return 1
	}

	// Dump the status
	c.Ui.Output(fmt.Sprintf(
		"Started: %v\n"+
			"Key Shares: %d\n"+
			"Key Threshold: %d\n"+
			"Rekey Progress: %d\n"+
			"Required Keys: %d",
		status.Started,
		status.N,
		status.T,
		status.Progress,
		status.Required,
	))
	return 0
}

func (c *RekeyCommand) Synopsis() string {
	return "Rekeys Vault to generate new unseal keys"
}

func (c *RekeyCommand) Help() string {
	helpText := `
Usage: vault rekey [options] [key]

  Rekey is used to change the unseal keys. This can be done to generate
  a new set of unseal keys or to change the number of shares and the
  required threshold.

  Rekey can only be done when the Vault is already unsealed. The operation
  is done online, but requires that a threshold of the current unseal
  keys be provided.

General Options:

  ` + generalOptionsUsage() + `

Unseal Options:

  -init                   Initialize the rekey operation by setting the desired
                          number of shares and the key threshold. This can only be
                          done if no rekey is already initiated.

  -cancel				  Reset the rekey process by throwing away
                          prior keys and the rekey configuration.

  -status                 Prints the status of the current rekey operation.
                          This can be used to see the status without attempting
                          to provide an unseal key.

  -key-shares=5           The number of key shares to split the master key
                          into.

  -key-threshold=3        The number of key shares required to reconstruct
                          the master key.
`
	return strings.TrimSpace(helpText)
}
