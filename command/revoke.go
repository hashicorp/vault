package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*ReadCommand)(nil)
var _ cli.CommandAutocomplete = (*ReadCommand)(nil)

// RevokeCommand is a Command that mounts a new mount.
type RevokeCommand struct {
	*BaseCommand

	flagForce  bool
	flagPrefix bool
}

func (c *RevokeCommand) Synopsis() string {
	return "Revokes leases and secrets"
}

func (c *RevokeCommand) Help() string {
	helpText := `
Usage: vault revoke [options] ID

  Revokes secrets by their lease ID. This command can revoke a single secret
  or multiple secrets based on a path-matched prefix.

  Revoke a single lease:

      $ vault revoke database/creds/readonly/2f6a614c...

  Revoke all leases for a role:

      $ vault revoke -prefix aws/creds/deploy

  Force delete leases from Vault even if backend revocation fails:

      $ vault revoke -force -prefix consul/creds

  For a full list of examples and paths, please see the documentation that
  corresponds to the secret backend in use.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *RevokeCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "force",
		Aliases: []string{"f"},
		Target:  &c.flagForce,
		Default: false,
		Usage: "Delete the lease from Vault even if the backend revocation " +
			"fails. This is meant for recovery situations where the secret " +
			"in the backend was manually removed. If this flag is specified, " +
			"-prefix is also required.",
	})

	f.BoolVar(&BoolVar{
		Name:    "prefix",
		Target:  &c.flagPrefix,
		Default: false,
		Usage: "Treat the ID as a prefix instead of an exact lease ID. This can " +
			"revoke multiple leases simultaneously.",
	})

	return set
}

func (c *RevokeCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *RevokeCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RevokeCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	leaseID, remaining, err := extractID(args)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if len(remaining) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	if c.flagForce && !c.flagPrefix {
		c.UI.Error("Specifying -force requires also specifying -prefix")
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	switch {
	case c.flagForce && c.flagPrefix:
		c.UI.Warn(wrapAtLength("Warning! Force-removing leases can cause Vault " +
			"to become out of sync with credential backends!"))
		if err := client.Sys().RevokeForce(leaseID); err != nil {
			c.UI.Error(fmt.Sprintf("Error force revoking leases with prefix %s: %s", leaseID, err))
			return 2
		}
		c.UI.Output(fmt.Sprintf("Success! Force revoked any leases with prefix: %s", leaseID))
		return 0
	case c.flagPrefix:
		if err := client.Sys().RevokePrefix(leaseID); err != nil {
			c.UI.Error(fmt.Sprintf("Error revoking leases with prefix %s: %s", leaseID, err))
			return 2
		}
		c.UI.Output(fmt.Sprintf("Success! Revoked any leases with prefix: %s", leaseID))
		return 0
	default:
		if err := client.Sys().Revoke(leaseID); err != nil {
			c.UI.Error(fmt.Sprintf("Error revoking lease %s: %s", leaseID, err))
			return 2
		}
		c.UI.Output(fmt.Sprintf("Success! Revoked lease: %s", leaseID))
		return 0
	}
}
