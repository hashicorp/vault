// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*ReadCommand)(nil)
	_ cli.CommandAutocomplete = (*ReadCommand)(nil)
)

type ReadCommand struct {
	*base.BaseCommand
}

func (c *ReadCommand) Synopsis() string {
	return "Read data and retrieves secrets"
}

func (c *ReadCommand) Help() string {
	helpText := `
Usage: vault read [options] PATH

  Reads data from Vault at the given path. This can be used to read secrets,
  generate dynamic credentials, get configuration details, and more.

  Read details of your own token:

      $ vault read auth/token/lookup-self

  Read entity details of a given ID:

      $ vault read identity/entity/id/2f09126d-d161-abb8-2241-555886491d97

  Generate credentials for my-role in an AWS secrets engine:

      $ vault read aws/creds/my-role

  For a full list of examples and paths, please see the documentation that
  corresponds to the secrets engine in use.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ReadCommand) Flags() *base.FlagSets {
	return c.FlagSetForBit(base.FlagSetHTTP | base.FlagSetOutputField | base.FlagSetOutputFormat | base.FlagSetSnapshot)
}

func (c *ReadCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *ReadCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ReadCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args, base.ParseOptionAllowRawFormat(true)); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// client.ReadRaw* methods require a manual timeout override
	ctx, cancel := context.WithTimeout(context.Background(), client.ClientTimeout())
	defer cancel()

	// Pull our fake stdin if needed
	stdin := c.Stdin()

	path := base.SanitizePath(args[0])

	data, err := base.ParseArgsDataStringLists(stdin, args[1:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	if c.FlagSnapshotID != "" {
		if data == nil {
			data = make(map[string][]string)
		}
		data["read_snapshot_id"] = []string{c.FlagSnapshotID}
	}

	if base.Format(c.UI) != "raw" {
		secret, err := client.Logical().ReadWithDataWithContext(ctx, path, data)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error reading %s: %s", path, err))
			return 2
		}
		if secret == nil {
			c.UI.Error(fmt.Sprintf("No value found at %s", path))
			return 2
		}

		if c.FlagField != "" {
			return base.PrintRawField(c.UI, secret, c.FlagField)
		}

		return base.OutputSecret(c.UI, secret)
	}

	resp, err := client.Logical().ReadRawWithDataWithContext(ctx, path, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading: %s: %s", path, err))
		return 2
	}
	if resp == nil || resp.Body == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", path))
		return 2
	}
	defer resp.Body.Close()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading: %s: %s", path, err))
		return 2
	}

	return base.OutputData(c.UI, contents)
}
