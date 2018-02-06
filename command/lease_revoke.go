package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*LeaseRevokeCommand)(nil)
var _ cli.CommandAutocomplete = (*LeaseRevokeCommand)(nil)

type LeaseRevokeCommand struct {
	*BaseCommand

	flagForce  bool
	flagPrefix bool
}

func (c *LeaseRevokeCommand) Synopsis() string {
	return "Revokes leases and secrets"
}

func (c *LeaseRevokeCommand) Help() string {
	helpText := `
Usage: vault lease revoke [options] ID

  Revokes secrets by their lease ID. This command can revoke a single secret
  or multiple secrets based on a path-matched prefix.

  Revoke a single lease:

      $ vault lease revoke database/creds/readonly/2f6a614c...

  Revoke all leases for a role:

      $ vault lease revoke -prefix aws/creds/deploy

  Force delete leases from Vault even if secret engine revocation fails:

      $ vault lease revoke -force -prefix consul/creds

  For a full list of examples and paths, please see the documentation that
  corresponds to the secret engine in use.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *LeaseRevokeCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "force",
		Aliases: []string{"f"},
		Target:  &c.flagForce,
		Default: false,
		Usage: "Delete the lease from Vault even if the secret engine revocation " +
			"fails. This is meant for recovery situations where the secret " +
			"in the target secret engine was manually removed. If this flag is " +
			"specified, -prefix is also required.",
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

func (c *LeaseRevokeCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *LeaseRevokeCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *LeaseRevokeCommand) Run(args []string) int {
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

	if c.flagForce && !c.flagPrefix {
		c.UI.Error("Specifying -force requires also specifying -prefix")
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	leaseID := strings.TrimSpace(args[0])

	switch {
	case c.flagForce && c.flagPrefix:
		c.UI.Warn(wrapAtLength("Warning! Force-removing leases can cause Vault " +
			"to become out of sync with secret engines!"))
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
