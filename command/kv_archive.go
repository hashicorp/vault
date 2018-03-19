package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*KVArchiveCommand)(nil)
var _ cli.CommandAutocomplete = (*KVArchiveCommand)(nil)

type KVArchiveCommand struct {
	*BaseCommand

	flagVersions []string
}

func (c *KVArchiveCommand) Synopsis() string {
	return "Archives versions in the KV store"
}

func (c *KVArchiveCommand) Help() string {
	helpText := `
Usage: vault kv archive [options] KEY

  Archives the data for the provided version and path in the key-value store.
  This marks the data as archived, but will not delete the underlying data.

  To archive version 3 of key foo:

      $ vault kv archive -versions=3 secret/foo

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVArchiveCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.StringSliceVar(&StringSliceVar{
		Name:    "versions",
		Target:  &c.flagVersions,
		Default: nil,
		Usage:   `Specifies the version numbers to archive.`,
	})

	return set
}

func (c *KVArchiveCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KVArchiveCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVArchiveCommand) Run(args []string) int {
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
		c.UI.Error("No versions provided, use the \"-versions\" flag to specify the version to archive.")
		return 1
	}
	var err error
	path := sanitizePath(args[0])
	path, err = addPrefixToVKVPath(path, "archive")
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
