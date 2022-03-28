package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*KVMetadataDeleteCommand)(nil)
	_ cli.CommandAutocomplete = (*KVMetadataDeleteCommand)(nil)
)

type KVMetadataDeleteCommand struct {
	*BaseCommand
	flagMount string
}

func (c *KVMetadataDeleteCommand) Synopsis() string {
	return "Deletes all versions and metadata for a key in the KV store"
}

func (c *KVMetadataDeleteCommand) Help() string {
	helpText := `
Usage: vault kv metadata delete [options] PATH

  Deletes all versions and metadata for the provided key. 

      $ vault kv metadata delete -mount=secret foo

  A more path-like syntax can also be used, but note that for KV v2, this is not the full API path to the secret (secret/metadata/foo): 
  
      $ vault kv metadata delete secret/foo

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *KVMetadataDeleteCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.StringVar(&StringVar{
		Name:    "mount",
		Target:  &c.flagMount,
		Default: "", // no default, because the handling of the next arg is determined by whether this flag has a value
		Usage: `Specifies the path where the KV backend is mounted. If specified, 
		the next argument will be interpreted as the secret path. If this flag is 
		not specified, the next argument will be interpreted as the combined mount 
		path and secret path, with /metadata/ automatically appended between KV 
		v2 secrets.`,
	})

	return set
}

func (c *KVMetadataDeleteCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *KVMetadataDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVMetadataDeleteCommand) Run(args []string) int {
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

	path := sanitizePath(args[0])
	mountPath, v2, err := isKVv2(path, client)
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}
	if !v2 {
		c.UI.Error("Metadata not supported on KV Version 1")
		return 1
	}

	path = addPrefixToKVPath(path, mountPath, "metadata")
	if secret, err := client.Logical().Delete(path); err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting %s: %s", path, err))
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}

	c.UI.Info(fmt.Sprintf("Success! Data deleted (if it existed) at: %s", path))
	return 0
}
