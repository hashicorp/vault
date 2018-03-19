package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*KVUnarchiveCommand)(nil)
var _ cli.CommandAutocomplete = (*KVUnarchiveCommand)(nil)

type KVUnarchiveCommand struct {
	*BaseCommand

	flagVersions []string
}

func (c *KVUnarchiveCommand) Synopsis() string {
	return "Unarchives versions in the KV store"
}

func (c *KVUnarchiveCommand) Help() string {
	helpText := `
Usage: vault kv unarchive [options] KEY

  Unarchives the data for the provided version and path in the key-value store.
  This restores the data, allowing it to be returned on get requests.

  To unarchive version 3 of key "foo":
  
      $ vault kv unarchive -versions=3 secret/foo

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVUnarchiveCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.StringSliceVar(&StringSliceVar{
		Name:    "versions",
		Target:  &c.flagVersions,
		Default: nil,
		Usage:   `Specifies the version numbers to unarchive.`,
	})

	return set
}

func (c *KVUnarchiveCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KVUnarchiveCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVUnarchiveCommand) Run(args []string) int {
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

	if len(c.flagVersions) == 0 {
		c.UI.Error("No versions provided, use the \"-versions\" flag to specify the version to unarchive.")
		return 1
	}
	var err error
	path := sanitizePath(args[0])
	path, err = addPrefixToVKVPath(path, "unarchive")
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	data := map[string]interface{}{
		"versions": c.flagVersions,
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	secret, err := client.Logical().Write(path, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %s", path, err))
		return 2
	}
	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", path))
		}
		return 0
	}

	return OutputSecret(c.UI, secret)
}
