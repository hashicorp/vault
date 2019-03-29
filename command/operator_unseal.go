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

var _ cli.Command = (*OperatorUnsealCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorUnsealCommand)(nil)

type OperatorUnsealCommand struct {
	*BaseCommand

	flagReset   bool
	flagMigrate bool

	testOutput io.Writer // for tests
}

func (c *OperatorUnsealCommand) Synopsis() string {
	return "Unseals the Vault server"
}

func (c *OperatorUnsealCommand) Help() string {
	helpText := `
Usage: vault operator unseal [options] [KEY]

  Provide a portion of the master key to unseal a Vault server. Vault starts
  in a sealed state. It cannot perform operations until it is unsealed. This
  command accepts a portion of the master key (an "unseal key").

  The unseal key can be supplied as an argument to the command, but this is
  not recommended as the unseal key will be available in your history:

      $ vault operator unseal IXyR0OJnSFobekZMMCKCoVEpT7wI6l+USMzE3IcyDyo=

  Instead, run the command with no arguments and it will prompt for the key:

      $ vault operator unseal
      Key (will be hidden): IXyR0OJnSFobekZMMCKCoVEpT7wI6l+USMzE3IcyDyo=

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorUnsealCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

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

	f.BoolVar(&BoolVar{
		Name:       "migrate",
		Aliases:    []string{},
		Target:     &c.flagMigrate,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage:      "Indicate that this share is provided with the intent that it is part of a seal migration process.",
	})

	return set
}

func (c *OperatorUnsealCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorUnsealCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorUnsealCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	unsealKey := ""

	args = f.Args()
	switch len(args) {
	case 0:
		// We will prompt for the unseal key later
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
		return OutputSealStatus(c.UI, client, status)
	}

	if unsealKey == "" {
		// Override the output
		writer := (io.Writer)(os.Stdout)
		if c.testOutput != nil {
			writer = c.testOutput
		}

		fmt.Fprintf(writer, "Unseal Key (will be hidden): ")
		value, err := password.Read(os.Stdin)
		fmt.Fprintf(writer, "\n")
		if err != nil {
			c.UI.Error(wrapAtLength(fmt.Sprintf("An error occurred attempting to "+
				"ask for an unseal key. The raw error message is shown below, but "+
				"usually this is because you attempted to pipe a value into the "+
				"unseal command or you are executing outside of a terminal (tty). "+
				"You should run the unseal command from a terminal for maximum "+
				"security. If this is not an option, the unseal key can be provided "+
				"as the first argument to the unseal command. The raw error "+
				"was:\n\n%s", err)))
			return 1
		}
		unsealKey = strings.TrimSpace(value)
	}

	status, err := client.Sys().UnsealWithOptions(&api.UnsealOpts{
		Key:     unsealKey,
		Migrate: c.flagMigrate,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error unsealing: %s", err))
		return 2
	}

	return OutputSealStatus(c.UI, client, status)
}
