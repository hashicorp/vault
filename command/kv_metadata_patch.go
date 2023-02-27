// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*KVMetadataPutCommand)(nil)
	_ cli.CommandAutocomplete = (*KVMetadataPutCommand)(nil)
)

type KVMetadataPatchCommand struct {
	*BaseCommand

	flagMaxVersions          int
	flagCASRequired          BoolPtr
	flagDeleteVersionAfter   time.Duration
	flagCustomMetadata       map[string]string
	flagRemoveCustomMetadata []string
	flagMount                string
	testStdin                io.Reader // for tests
}

func (c *KVMetadataPatchCommand) Synopsis() string {
	return "Patches key settings in the KV store"
}

func (c *KVMetadataPatchCommand) Help() string {
	helpText := `
Usage: vault kv metadata patch [options] KEY

  This command can be used to create a blank key in the key-value store or to
  update key configuration for a specified key.

  Create a key in the key-value store with no data:

      $ vault kv metadata patch -mount=secret foo

  The deprecated path-like syntax can also be used, but this should be avoided 
  for KV v2, as the fact that it is not actually the full API path to 
  the secret (secret/metadata/foo) can cause confusion: 
  
      $ vault kv metadata patch secret/foo

  Set a max versions setting on the key:

      $ vault kv metadata patch -mount=secret -max-versions=5 foo

  Set delete-version-after on the key:

      $ vault kv metadata patch -mount=secret -delete-version-after=3h25m19s foo

  Require Check-and-Set for this key:

      $ vault kv metadata patch -mount=secret -cas-required foo

  Set custom metadata on the key:

      $ vault kv metadata patch -mount=secret -custom-metadata=foo=abc -custom-metadata=bar=123 foo

  To remove custom meta data from the corresponding path in the key-value store, kv metadata patch can be used.

      $ vault kv metadata patch -mount=secret -remove-custom-metadata=bar foo

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVMetadataPatchCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.IntVar(&IntVar{
		Name:    "max-versions",
		Target:  &c.flagMaxVersions,
		Default: -1,
		Usage:   `The number of versions to keep. If not set, the backend’s configured max version is used.`,
	})

	f.BoolPtrVar(&BoolPtrVar{
		Name:   "cas-required",
		Target: &c.flagCASRequired,
		Usage:  `If true the key will require the cas parameter to be set on all write requests. If false, the backend’s configuration will be used.`,
	})

	f.DurationVar(&DurationVar{
		Name:       "delete-version-after",
		Target:     &c.flagDeleteVersionAfter,
		Default:    -1,
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: `Specifies the length of time before a version is deleted.
		If not set, the backend's configured delete-version-after is used. Cannot be
		greater than the backend's delete-version-after. The delete-version-after is
		specified as a numeric string with a suffix like "30s" or
		"3h25m19s".`,
	})

	f.StringMapVar(&StringMapVar{
		Name:    "custom-metadata",
		Target:  &c.flagCustomMetadata,
		Default: map[string]string{},
		Usage: `Specifies arbitrary version-agnostic key=value metadata meant to describe a secret.
		This can be specified multiple times to add multiple pieces of metadata.`,
	})

	f.StringSliceVar(&StringSliceVar{
		Name:    "remove-custom-metadata",
		Target:  &c.flagRemoveCustomMetadata,
		Default: []string{},
		Usage:   "Key to remove from custom metadata. To specify multiple values, specify this flag multiple times.",
	})

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

func (c *KVMetadataPatchCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KVMetadataPatchCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVMetadataPatchCommand) Run(args []string) int {
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
		c.UI.Error("Metadata not supported on KV Version 1")
		return 1
	}

	fullPath := addPrefixToKVPath(partialPath, mountPath, "metadata")

	data := make(map[string]interface{}, 0)

	if c.flagMaxVersions >= 0 {
		data["max_versions"] = c.flagMaxVersions
	}

	if c.flagCASRequired.IsSet() {
		data["cas_required"] = c.flagCASRequired.Get()
	}

	if c.flagDeleteVersionAfter >= 0 {
		data["delete_version_after"] = c.flagDeleteVersionAfter.String()
	}

	customMetadata := make(map[string]interface{})

	for key, value := range c.flagCustomMetadata {
		customMetadata[key] = value
	}

	for _, key := range c.flagRemoveCustomMetadata {
		// A null in a JSON merge patch payload will remove the associated key
		customMetadata[key] = nil
	}

	data["custom_metadata"] = customMetadata

	secret, err := client.Logical().JSONMergePatch(context.Background(), fullPath, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %s", fullPath, err))

		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}

	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", fullPath))
		}
		return 0
	}

	return OutputSecret(c.UI, secret)
}
