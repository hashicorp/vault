// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*AgentInitConfigCommand)(nil)
	_ cli.CommandAutocomplete = (*AgentInitConfigCommand)(nil)
)

type AgentInitConfigCommand struct {
	*BaseCommand

	flagConfigType string
	flagPaths      []string
}

func (c *AgentInitConfigCommand) Synopsis() string {
	return "Create a vault agent configuration file."
}

func (c *AgentInitConfigCommand) Help() string {
	helpText := `
Usage: vault agent init-config [options]
` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *AgentInitConfigCommand) Flags() *FlagSets {
	set := NewFlagSets(c.UI)

	// Common Options
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:    "config-type",
		Target:  &c.flagConfigType,
		Default: "env-template",
		Usage:   "The type of configuration file to generate, currently only 'env-template' is supported.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   "path",
		Target: &c.flagPaths,
		Usage:  "Path to a KV v1/v2 secret (e.g. secret/data/foo, secret/prefix/*).",
	})

	return set
}

func (c *AgentInitConfigCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFolders()
}

func (c *AgentInitConfigCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AgentInitConfigCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	// Pull our fake stdin if needed
	stdin := (io.Reader)(os.Stdin)
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected >1, got %d)", len(args)))
		return 1
	case len(args) == 1:
		c.UI.Error("Must supply data")
		return 1
	}

	var err error

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	data, err := parseArgsData(stdin, args[1:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
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

	// Add /data to v2 paths only
	var fullPath string
	if v2 {
		fullPath = addPrefixToKVPath(partialPath, mountPath, "data")
		data = map[string]interface{}{
			"data":    data,
			"options": map[string]interface{}{},
		}

		if c.flagCAS > -1 {
			data["options"].(map[string]interface{})["cas"] = c.flagCAS
		}
	} else {
		// v1
		if mountFlagSyntax {
			fullPath = path.Join(mountPath, partialPath)
		} else {
			fullPath = partialPath
		}
	}

	secret, err := client.Logical().Write(fullPath, data)
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

	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	if Format(c.UI) == "table" {
		outputPath(c.UI, fullPath, "Secret Path")
		metadata := secret.Data
		c.UI.Info(getHeaderForMap("Metadata", metadata))
		return OutputData(c.UI, metadata)
	}

	return OutputSecret(c.UI, secret)
}
