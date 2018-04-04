package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*KVDeleteCommand)(nil)
var _ cli.CommandAutocomplete = (*KVDeleteCommand)(nil)

type KVDeleteCommand struct {
	*BaseCommand

	flagVersions []string
}

func (c *KVDeleteCommand) Synopsis() string {
	return "Deletes versions in the KV store"
}

func (c *KVDeleteCommand) Help() string {
	helpText := `
Usage: vault kv delete [options] PATH

  Deletes the data for the provided version and path in the key-value store. The
  versioned data will not be fully removed, but marked as deleted and will no
  longer be returned in normal get requests.

  To delete the latest version of the key "foo": 

      $ vault kv delete secret/foo

  To delete version 3 of key foo:

      $ vault kv delete -versions=3 secret/foo

  To delete all versions and metadata, see the "vault kv metadata" subcommand.

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *KVDeleteCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	// Common Options
	f := set.NewFlagSet("Common Options")

	f.StringSliceVar(&StringSliceVar{
		Name:    "versions",
		Target:  &c.flagVersions,
		Default: nil,
		Usage:   `Specifies the version numbers to delete.`,
	})

	return set
}

func (c *KVDeleteCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *KVDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVDeleteCommand) Run(args []string) int {
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

	path := sanitizePath(args[0])
	var err error
	if len(c.flagVersions) > 0 {
		err = c.deleteVersions(path, kvParseVersionsFlags(c.flagVersions))
	} else {
		err = c.deleteLatest(path)
	}
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting %s: %s", path, err))
		return 2
	}

	c.UI.Info(fmt.Sprintf("Success! Data deleted (if it existed) at: %s", path))
	return 0
}

func (c *KVDeleteCommand) deleteLatest(path string) error {
	var err error
	path, err = addPrefixToVKVPath(path, "data")
	if err != nil {
		return err
	}

	client, err := c.Client()
	if err != nil {
		return err
	}

	_, err = kvDeleteRequest(client, path)

	return err
}

func (c *KVDeleteCommand) deleteVersions(path string, versions []string) error {
	var err error
	path, err = addPrefixToVKVPath(path, "delete")
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"versions": versions,
	}

	client, err := c.Client()
	if err != nil {
		return err
	}

	_, err = kvWriteRequest(client, path, data)
	return err
}
