package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*KVEnableVersioningCommand)(nil)
var _ cli.CommandAutocomplete = (*KVEnableVersioningCommand)(nil)

type KVEnableVersioningCommand struct {
	*BaseCommand
}

func (c *KVEnableVersioningCommand) Synopsis() string {
	return "Turns on versioning for a KV store"
}

func (c *KVEnableVersioningCommand) Help() string {
	helpText := `
Usage: vault kv enable-versions [options] KEY

  This command turns on versioning for the backend at the provided path.

      $ vault kv enable-versions secret

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVEnableVersioningCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	return set
}

func (c *KVEnableVersioningCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KVEnableVersioningCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVEnableVersioningCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// Append a trailing slash to indicate it's a path in output
	mountPath := ensureTrailingSlash(sanitizePath(args[0]))

	if err := client.Sys().TuneMount(mountPath, api.MountConfigInput{
		Options: map[string]string{
			"version": "2",
		},
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error tuning secrets engine %s: %s", mountPath, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Tuned the secrets engine at: %s", mountPath))
	return 0
}
