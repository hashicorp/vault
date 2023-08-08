// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"fmt"
	"path"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*KVDestroyCommand)(nil)
	_ cli.CommandAutocomplete = (*KVDestroyCommand)(nil)
)

type KVDestroyCommand struct {
	*BaseCommand

	flagVersions []string
	flagMount    string
}

func (c *KVDestroyCommand) Synopsis() string {
	return "Permanently removes one or more versions in the KV store"
}

func (c *KVDestroyCommand) Help() string {
	helpText := `
Usage: vault kv destroy [options] KEY

  Permanently removes the specified versions' data from the key-value store. If
  no key exists at the path, no action is taken.

  To destroy version 3 of key foo:

      $ vault kv destroy -mount=secret -versions=3 foo

  The deprecated path-like syntax can also be used, but this should be avoided 
  for KV v2, as the fact that it is not actually the full API path to 
  the secret (secret/data/foo) can cause confusion: 
  
      $ vault kv destroy -versions=3 secret/foo

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
		c.UI.Error("Destroy not supported on KV Version 1")
		return 1
	}
	destroyPath := addPrefixToKVPath(partialPath, mountPath, "destroy", false)
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	data := map[string]interface{}{
		"versions": kvParseVersionsFlags(c.flagVersions),
	}

	secret, err := client.Logical().Write(destroyPath, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %s", destroyPath, err))
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}
	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", destroyPath))
		}
		return 0
	}

	return OutputSecret(c.UI, secret)
}
