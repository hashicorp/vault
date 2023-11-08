// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*AuthMoveCommand)(nil)
	_ cli.CommandAutocomplete = (*AuthMoveCommand)(nil)
)

type AuthMoveCommand struct {
	*BaseCommand
}

func (c *AuthMoveCommand) Synopsis() string {
	return "Move an auth method to a new path"
}

func (c *AuthMoveCommand) Help() string {
	helpText := `
Usage: vault auth move [options] SOURCE DESTINATION

  Moves an existing auth method to a new path. Any leases from the old
  auth method are revoked, but all configuration associated with the method
  is preserved. It initiates the migration and intermittently polls its status,
  exiting if a final state is reached.

  This command works within or across namespaces, both source and destination paths
  can be prefixed with a namespace heirarchy relative to the current namespace.

  WARNING! Moving an auth method will revoke any leases from the
  old method.

  Move the auth method at approle/ to generic/:

      $ vault auth move approle/ generic/

  Move the auth method at ns1/approle/ across namespaces to ns2/generic/, 
  where ns1 and ns2 are child namespaces of the current namespace:

  $ vault auth move ns1/approle/ ns2/generic/

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuthMoveCommand) Flags() *FlagSets {
	return c.FlagSet(FlagSetHTTP)
}

func (c *AuthMoveCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultMounts()
}

func (c *AuthMoveCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthMoveCommand) Run(args []string) int {
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

	// Grab the source and destination
	source := ensureTrailingSlash(args[0])
	destination := ensureTrailingSlash(args[1])

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	remountResp, err := client.Sys().StartRemount(source, destination)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error moving auth method %s to %s: %s", source, destination, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Started moving auth method %s to %s, with migration ID %s", source, destination, remountResp.MigrationID))

	// Poll the status endpoint with the returned migration ID
	// Exit if a terminal status is reached, else wait and retry
	for {
		remountStatusResp, err := client.Sys().RemountStatus(remountResp.MigrationID)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error checking migration status of auth method %s to %s: %s", source, destination, err))
			return 2
		}
		if remountStatusResp.MigrationInfo.MigrationStatus == MountMigrationStatusSuccess {
			c.UI.Output(fmt.Sprintf("Success! Finished moving auth method %s to %s, with migration ID %s", source, destination, remountResp.MigrationID))
			return 0
		}
		if remountStatusResp.MigrationInfo.MigrationStatus == MountMigrationStatusFailure {
			c.UI.Error(fmt.Sprintf("Failure! Error encountered moving auth method %s to %s, with migration ID %s", source, destination, remountResp.MigrationID))
			return 0
		}
		c.UI.Output(fmt.Sprintf("Waiting for terminal status in migration of auth method %s to %s, with migration ID %s", source, destination, remountResp.MigrationID))
		time.Sleep(10 * time.Second)
	}

	return 0
}
