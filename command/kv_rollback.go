// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"flag"
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*KVRollbackCommand)(nil)
	_ cli.CommandAutocomplete = (*KVRollbackCommand)(nil)
)

type KVRollbackCommand struct {
	*BaseCommand

	flagVersion int
	flagMount   string
}

func (c *KVRollbackCommand) Synopsis() string {
	return "Rolls back to a previous version of data"
}

func (c *KVRollbackCommand) Help() string {
	helpText := `
Usage: vault kv rollback [options] KEY

  *NOTE*: This is only supported for KV v2 engine mounts.

  Restores a given previous version to the current version at the given path.
  The value is written as a new version; for instance, if the current version
  is 5 and the rollback version is 2, the data from version 2 will become
  version 6.

      $ vault kv rollback -mount=secret -version=2 foo

  The deprecated path-like syntax can also be used, but this should be avoided, 
  as the fact that it is not actually the full API path to 
  the secret (secret/data/foo) can cause confusion: 
  
      $ vault kv rollback -version=2 secret/foo

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVRollbackCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.IntVar(&IntVar{
		Name:   "version",
		Target: &c.flagVersion,
		Usage:  `Specifies the version number that should be made current again.`,
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

func (c *KVRollbackCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KVRollbackCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVRollbackCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	var version *int
	f.Visit(func(fl *flag.Flag) {
		if fl.Name == "version" {
			version = &c.flagVersion
		}
	})

	args = f.Args()

	switch {
	case len(args) != 1:
		c.UI.Error(fmt.Sprintf("Invalid number of arguments (expected 1, got %d)", len(args)))
		return 1
	case version == nil:
		c.UI.Error("Version flag must be specified")
		return 1
	case c.flagVersion <= 0:
		c.UI.Error(fmt.Sprintf("Invalid value %d for the version flag", c.flagVersion))
		return 1
	}

	var err error

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// If true, we're working with "-mount=secret foo" syntax.
	// If false, we're using "secret/foo" syntax.
	mountFlagSyntax := c.flagMount != ""

	var (
		mountPath   string
		partialPath string
		v2          bool
	)

	// Parse the paths and grab the KV version
	if mountFlagSyntax {
		// In this case, this arg is the secret path (e.g. "foo").
		partialPath = sanitizePath(args[0])
		mountPath, v2, err = isKVv2(sanitizePath(c.flagMount), client)
		if err != nil {
			c.UI.Error(err.Error())
			return 2
		}

		if v2 {
			partialPath = path.Join(mountPath, partialPath)
		}
	} else {
		// In this case, this arg is a path-like combination of mountPath/secretPath.
		// (e.g. "secret/foo")
		partialPath = sanitizePath(args[0])
		mountPath, v2, err = isKVv2(partialPath, client)
		if err != nil {
			c.UI.Error(err.Error())
			return 2
		}
	}

	if !v2 {
		c.UI.Error("KV engine mount must be version 2 for rollback support")
		return 2
	}

	fullPath := addPrefixToKVPath(partialPath, mountPath, "data", false)
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// First, do a read to get the current version for check-and-set
	var meta map[string]interface{}
	{
		secret, err := kvReadRequest(client, fullPath, nil)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error doing pre-read at %s: %s", fullPath, err))
			return 2
		}

		// Make sure a value already exists
		if secret == nil || secret.Data == nil {
			c.UI.Error(fmt.Sprintf("No value found at %s", fullPath))
			return 2
		}

		// Verify metadata found
		rawMeta, ok := secret.Data["metadata"]
		if !ok || rawMeta == nil {
			c.UI.Error(fmt.Sprintf("No metadata found at %s; rollback only works on existing data", fullPath))
			return 2
		}
		meta, ok = rawMeta.(map[string]interface{})
		if !ok {
			c.UI.Error(fmt.Sprintf("Metadata found at %s is not the expected type (JSON object)", fullPath))
			return 2
		}
		if meta == nil {
			c.UI.Error(fmt.Sprintf("No metadata found at %s; rollback only works on existing data", fullPath))
			return 2
		}
	}

	casVersion := meta["version"]

	// Set the version parameter
	versionParam := map[string]string{
		"version": fmt.Sprintf("%d", c.flagVersion),
	}

	// Now run it again and read the version we want to roll back to
	var data map[string]interface{}
	{
		secret, err := kvReadRequest(client, fullPath, versionParam)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error doing pre-read at %s: %s", fullPath, err))
			return 2
		}

		// Make sure a value already exists
		if secret == nil || secret.Data == nil {
			c.UI.Error(fmt.Sprintf("No value found at %s", fullPath))
			return 2
		}

		// Verify metadata found
		rawMeta, ok := secret.Data["metadata"]
		if !ok || rawMeta == nil {
			c.UI.Error(fmt.Sprintf("No metadata found at %s; rollback only works on existing data", fullPath))
			return 2
		}
		meta, ok := rawMeta.(map[string]interface{})
		if !ok {
			c.UI.Error(fmt.Sprintf("Metadata found at %s is not the expected type (JSON object)", fullPath))
			return 2
		}
		if meta == nil {
			c.UI.Error(fmt.Sprintf("No metadata found at %s; rollback only works on existing data", fullPath))
			return 2
		}

		// Verify it hasn't been deleted
		if meta["deletion_time"] != nil && meta["deletion_time"].(string) != "" {
			c.UI.Error("Cannot roll back to a version that has been deleted")
			return 2
		}

		if meta["destroyed"] != nil && meta["destroyed"].(bool) {
			c.UI.Error("Cannot roll back to a version that has been destroyed")
			return 2
		}

		// Verify old data found
		rawData, ok := secret.Data["data"]
		if !ok || rawData == nil {
			c.UI.Error(fmt.Sprintf("No data found at %s; rollback only works on existing data", fullPath))
			return 2
		}
		data, ok = rawData.(map[string]interface{})
		if !ok {
			c.UI.Error(fmt.Sprintf("Data found at %s is not the expected type (JSON object)", fullPath))
			return 2
		}
		if data == nil {
			c.UI.Error(fmt.Sprintf("No data found at %s; rollback only works on existing data", fullPath))
			return 2
		}
	}

	secret, err := client.Logical().Write(fullPath, map[string]interface{}{
		"data": data,
		"options": map[string]interface{}{
			"cas": casVersion,
		},
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %s", fullPath, err))
		return 2
	}
	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", fullPath))
		}
		return 0
	}

	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, secret)
}
