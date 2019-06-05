package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*KVUndeleteCommand)(nil)
var _ cli.CommandAutocomplete = (*KVUndeleteCommand)(nil)

type KVUndeleteCommand struct {
	*BaseCommand

	flagVersions           []string
	flagDeleteVersionAfter time.Duration
}

func (c *KVUndeleteCommand) Synopsis() string {
	return "Undeletes versions in the KV store"
}

func (c *KVUndeleteCommand) Help() string {
	helpText := `
Usage: vault kv undelete [options] KEY

  Undeletes the data for the provided version and path in the key-value store.
  This restores the data, allowing it to be returned on get requests.

  To undelete version 3 of key "foo":

      $ vault kv undelete -versions=3 secret/foo

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVUndeleteCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.StringSliceVar(&StringSliceVar{
		Name:    "versions",
		Target:  &c.flagVersions,
		Default: nil,
		Usage:   `Specifies the version numbers to undelete.`,
	})

	f.DurationVar(&DurationVar{
		Name:       "delete-version-after",
		Target:     &c.flagDeleteVersionAfter,
		Default:    0,
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: `Specifies the length of time before these versions will be
		deleted. If not set, the metadata's delete-version-after is used.
		Cannot be greater than the metadata's delete-version-after. The
		delete-version-after is specified as a numeric string with a suffix
		like "30s" or
		"3h25m19s".`,
	})

	return set
}

func (c *KVUndeleteCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KVUndeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVUndeleteCommand) Run(args []string) int {
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
		c.UI.Error("No versions provided, use the \"-versions\" flag to specify the version to undelete.")
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
		c.UI.Error("Undelete not supported on KV Version 1")
		return 1
	}

	path = addPrefixToVKVPath(path, mountPath, "undelete")
	data := map[string]interface{}{
		"versions": kvParseVersionsFlags(c.flagVersions),
	}

	if c.flagDeleteVersionAfter > 0 {
		data["delete_version_after"] = c.flagDeleteVersionAfter.String()
	}

	secret, err := client.Logical().Write(path, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %s", path, err))
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
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
