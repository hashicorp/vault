// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PatchCommand)(nil)
	_ cli.CommandAutocomplete = (*PatchCommand)(nil)
)

// PatchCommand is a Command that puts data into the Vault.
type PatchCommand struct {
	*BaseCommand

	flagForce bool

	testStdin io.Reader // for tests
}

func (c *PatchCommand) Synopsis() string {
	return "Patch data, configuration, and secrets"
}

func (c *PatchCommand) Help() string {
	helpText := `
Usage: vault patch [options] PATH [DATA K=V...]

  Patches data in Vault at the given path. The data can be credentials, secrets,
  configuration, or arbitrary data. The specific behavior of this command is
  determined at the thing mounted at the path.

  Data is specified as "key=value" pairs. If the value begins with an "@", then
  it is loaded from a file. If the value is "-", Vault will read the value from
  stdin.

  Unlike write, patch will only modify specified fields.

  Persist data in the generic secrets engine without modifying any other fields:

      $ vault patch pki/roles/example allow_localhost=false

  The data can also be consumed from a file on disk by prefixing with the "@"
  symbol. For example:

      $ vault patch pki/roles/example @role.json

  Or it can be read from stdin using the "-" symbol:

      $ echo "example.com" | vault patch pki/roles/example allowed_domains=-

  For a full list of examples and paths, please see the documentation that
  corresponds to the secret engines in use.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PatchCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:       "force",
		Aliases:    []string{"f"},
		Target:     &c.flagForce,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage: "Allow the operation to continue with no key=value pairs. This " +
			"allows writing to keys that do not need or expect data.",
	})

	return set
}

func (c *PatchCommand) AutocompleteArgs() complete.Predictor {
	// Return an anything predictor here. Without a way to access help
	// information, we don't know what paths we could patch.
	return complete.PredictAnything
}

func (c *PatchCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PatchCommand) Run(args []string) int {
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
	case len(args) == 1 && !c.flagForce:
		c.UI.Error("Must supply data or use -force")
		return 1
	}

	// Pull our fake stdin if needed
	stdin := (io.Reader)(os.Stdin)
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	path := sanitizePath(args[0])

	data, err := parseArgsData(stdin, args[1:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	secret, err := client.Logical().JSONMergePatch(context.Background(), path, data)
	return handleWriteSecretOutput(c.BaseCommand, path, secret, err)
}
