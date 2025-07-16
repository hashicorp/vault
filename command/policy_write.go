// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PolicyWriteCommand)(nil)
	_ cli.CommandAutocomplete = (*PolicyWriteCommand)(nil)
)

type PolicyWriteCommand struct {
	*BaseCommand

	testStdin io.Reader // for tests
}

func (c *PolicyWriteCommand) Synopsis() string {
	return "Uploads a named policy from a file"
}

func (c *PolicyWriteCommand) Help() string {
	helpText := `
Usage: vault policy write [options] NAME PATH

  Uploads a policy with name NAME from the contents of a local file PATH or
  stdin. If PATH is "-", the policy is read from stdin. Otherwise, it is
  loaded from the file at the given path on the local disk.

  Upload a policy named "my-policy" from "/tmp/policy.hcl" on the local disk:

      $ vault policy write my-policy /tmp/policy.hcl

  Upload a policy from stdin:

      $ cat my-policy.hcl | vault policy write my-policy -

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PolicyWriteCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *PolicyWriteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictFunc(func(args complete.Args) []string {
		// Predict the LAST argument hcl files - we don't want to predict the
		// name argument as a filepath.
		if len(args.All) == 3 {
			return complete.PredictFiles("*.hcl").Predict(args)
		}
		return nil
	})
}

func (c *PolicyWriteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PolicyWriteCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 2:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2, got %d)", len(args)))
		return 1
	case len(args) > 2:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 2, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// Policies are normalized to lowercase
	policyName := args[0]
	formattedName := strings.TrimSpace(strings.ToLower(policyName))
	path := strings.TrimSpace(args[1])

	// Get the policy contents, either from stdin of a file
	var reader io.Reader
	if path == "-" {
		reader = os.Stdin
		if c.testStdin != nil {
			reader = c.testStdin
		}
	} else {
		file, err := os.Open(path)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error opening policy file: %s", err))
			return 2
		}
		defer file.Close()
		reader = file
	}

	// Read the policy
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		c.UI.Error(fmt.Sprintf("Error reading policy: %s", err))
		return 2
	}
	rules := buf.String()

	if err := client.Sys().PutPolicy(formattedName, rules); err != nil {
		c.UI.Error(fmt.Sprintf("Error uploading policy: %s", err))
		return 2
	}

	if policyName != formattedName {
		c.UI.Warn(fmt.Sprintf("Policy name was converted from \"%s\" to \"%s\"", policyName, formattedName))
	}

	c.UI.Output(fmt.Sprintf("Success! Uploaded policy: %s", formattedName))
	return 0
}
