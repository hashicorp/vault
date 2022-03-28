package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*KVDeleteCommand)(nil)
	_ cli.CommandAutocomplete = (*KVDeleteCommand)(nil)
)

type KVDeleteCommand struct {
	*BaseCommand

	flagVersions []string
	flagMount    string
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

      $ vault kv delete -mount=secret foo

  A more path-like syntax can also be used, but note that for KV v2, this is not the full API path to the secret (secret/data/foo): 
  
      $ vault kv delete secret/foo

  To delete version 3 of key foo:

      $ vault kv delete -mount=secret -versions=3 foo

  To delete all versions and metadata, see the "vault kv metadata" subcommand.

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *KVDeleteCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)
	// Common Options
	f := set.NewFlagSet("Common Options")

	f.StringSliceVar(&StringSliceVar{
		Name:    "versions",
		Target:  &c.flagVersions,
		Default: nil,
		Usage:   `Specifies the version numbers to delete.`,
	})

	f.StringVar(&StringVar{
		Name:    "mount",
		Target:  &c.flagMount,
		Default: "", // no default, because the handling of the next arg is determined by whether this flag has a value
		Usage: `Specifies the path where the KV backend is mounted. If specified, 
		the next argument will be interpreted as the secret path. If this flag is 
		not specified, the next argument will be interpreted as the combined mount 
		path and secret path, with /data/ automatically appended between KV 
		v2 secrets.`,
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

	var secret *api.Secret
	if v2 {
		secret, err = c.deleteV2(path, mountPath, client)
	} else {
		secret, err = client.Logical().Delete(path)
	}

	if err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting %s: %s", path, err))
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}

	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data deleted (if it existed) at: %s", path))
		}
		return 0
	}

	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, secret)
}

func (c *KVDeleteCommand) deleteV2(path, mountPath string, client *api.Client) (*api.Secret, error) {
	var err error
	var secret *api.Secret
	switch {
	case len(c.flagVersions) > 0:
		path = addPrefixToKVPath(path, mountPath, "delete")
		if err != nil {
			return nil, err
		}

		data := map[string]interface{}{
			"versions": kvParseVersionsFlags(c.flagVersions),
		}

		secret, err = client.Logical().Write(path, data)
	default:

		path = addPrefixToKVPath(path, mountPath, "data")
		if err != nil {
			return nil, err
		}

		secret, err = client.Logical().Delete(path)
	}

	return secret, err
}
