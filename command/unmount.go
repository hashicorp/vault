package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*UnmountCommand)(nil)
var _ cli.CommandAutocomplete = (*UnmountCommand)(nil)

// UnmountCommand is a Command that mounts a new mount.
type UnmountCommand struct {
	*BaseCommand
}

func (c *UnmountCommand) Synopsis() string {
	return "Unmounts a secret backend"
}

func (c *UnmountCommand) Help() string {
	helpText := `
Usage: vault unmount [options] PATH

  Unmounts a secret backend at the given PATH. The argument corresponds to
  the PATH of the mount, not the TYPE! All secrets created by this backend
  are revoked and its Vault data is removed.

  If no mount exists at the given path, the command will still return as
  successful because unmounting is an idempotent operation.

  Unmount the secret backend mounted at aws/:

      $ vault unmount aws/

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *UnmountCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *UnmountCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultMounts()
}

func (c *UnmountCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *UnmountCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	mountPath, remaining, err := extractPath(args)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if len(remaining) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// Append a trailing slash to indicate it's a path in output
	mountPath = ensureTrailingSlash(mountPath)

	if err := client.Sys().Unmount(mountPath); err != nil {
		c.UI.Error(fmt.Sprintf("Error unmounting %s: %s", mountPath, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Unmounted the secret backend (if it existed) at: %s", mountPath))
	return 0
}
