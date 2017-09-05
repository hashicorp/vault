package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/password"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*UnsealCommand)(nil)
var _ cli.CommandAutocomplete = (*UnsealCommand)(nil)

// UnsealCommand is a Command that unseals the vault.
type UnsealCommand struct {
	*BaseCommand

	flagReset bool

	testOutput io.Writer // for tests
}

func (c *UnsealCommand) Synopsis() string {
	return "Unseals the Vault server"
}

func (c *UnsealCommand) Help() string {
	helpText := `
Usage: vault unseal [options] [KEY]

  Provide a portion of the master key to unseal a Vault server. Vault starts
  in a sealed state. It cannot perform operations until it is unsealed. This
  command accepts a portion of the master key (an "unseal key").

  The unseal key can be supplied as an argument to the command, but this is
  not recommended as the unseal key will be available in your history:

      $ vault unseal IXyR0OJnSFobekZMMCKCoVEpT7wI6l+USMzE3IcyDyo=

  Instead, run the command with no arguments and it will prompt for the key:

      $ vault unseal
      Key (will be hidden): IXyR0OJnSFobekZMMCKCoVEpT7wI6l+USMzE3IcyDyo=

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *UnsealCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:       "reset",
		Aliases:    []string{},
		Target:     &c.flagReset,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage:      "Discard any previously entered keys to the unseal process.",
	})

	return set
}

func (c *UnsealCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *UnsealCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *UnsealCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	unsealKey := ""

	args = f.Args()
	switch len(args) {
	case 0:
		// We will prompt for the unsealKey later
	case 1:
		unsealKey = strings.TrimSpace(args[0])
	default:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if c.flagReset {
		status, err := client.Sys().ResetUnsealProcess()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error resetting unseal process: %s", err))
			return 2
		}
		c.prettySealStatus(status)
		return 0
	}

	if unsealKey == "" {
		// Override the output
		writer := (io.Writer)(os.Stdout)
		if c.testOutput != nil {
			writer = c.testOutput
		}

		fmt.Fprintf(writer, "Key (will be hidden): ")
		value, err := password.Read(os.Stdin)
		fmt.Fprintf(writer, "\n")
		if err != nil {
			c.UI.Error(wrapAtLength(fmt.Sprintf("An error occurred attempting to "+
				"ask for an unseal key. The raw error message is shown below, but "+
				"usually this is because you attempted to pipe a value into the "+
				"unseal command or you are executing outside of a terminal (tty). "+
				"You should run the unseal command from a terminal for maximum "+
				"security. If this is not an option, the unseal can be provided as "+
				"the first argument to the unseal command. The raw error "+
				"was:\n\n%s", err)))
			return 1
		}
		unsealKey = strings.TrimSpace(value)
	}

	status, err := client.Sys().Unseal(unsealKey)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error unsealing: %s", err))
		return 2
	}

	c.prettySealStatus(status)
	return 0
}

func (c *UnsealCommand) prettySealStatus(status *api.SealStatusResponse) {
	c.UI.Output(fmt.Sprintf("Sealed: %t", status.Sealed))
	c.UI.Output(fmt.Sprintf("Key Shares: %d", status.N))
	c.UI.Output(fmt.Sprintf("Key Threshold: %d", status.T))
	c.UI.Output(fmt.Sprintf("Unseal Progress: %d", status.Progress))
	if status.Nonce != "" {
		c.UI.Output(fmt.Sprintf("Unseal Nonce: %s", status.Nonce))
	}
}
