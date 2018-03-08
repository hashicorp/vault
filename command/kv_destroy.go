package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*KVDestroyCommand)(nil)
var _ cli.CommandAutocomplete = (*KVDestroyCommand)(nil)

type KVDestroyCommand struct {
	*BaseCommand

	flagVersions []string
}

func (c *KVDestroyCommand) Synopsis() string {
	return "Destroys a version in the KV store"
}

func (c *KVDestroyCommand) Help() string {
	helpText := `
Usage: vault kv destroy [options] KEY

  Writes the data to the given path in the key-value store. The data can be of
  any type.

      $ vault kv put secret/foo bar=baz

  The data can also be consumed from a file on disk by prefixing with the "@"
  symbol. For example:

      $ vault kv put secret/foo @data.json

  Or it can be read from stdin using the "-" symbol:

      $ echo "abcd1234" | vault kv put secret/foo bar=-

  To perform a Check-And-Set operation, specify the -cas flag with the
  appropriate version numer corresponding to the key you want to perform
  the CAS operation on:

      $ vault kv put -cas=1 secret/foo bar=baz

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVDestroyCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.StringSliceVar(&StringSliceVar{
		Name:    "versions",
		Target:  &c.flagVersions,
		Default: nil,
		Usage:   `Specifies the version numbers to destroy.`,
	})

	return set
}

func (c *KVDestroyCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KVDestroyCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVDestroyCommand) Run(args []string) int {
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
		c.UI.Error("No versions provided, use the \"-versions\" flag to specify the version to destroy.")
		return 1
	}
	var err error
	path := sanitizePath(args[0])
	path, err = addPrefixToVKVPath(path, "destroy")
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
