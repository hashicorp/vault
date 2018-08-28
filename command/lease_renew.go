package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*LeaseRenewCommand)(nil)
var _ cli.CommandAutocomplete = (*LeaseRenewCommand)(nil)

type LeaseRenewCommand struct {
	*BaseCommand

	flagIncrement time.Duration
}

func (c *LeaseRenewCommand) Synopsis() string {
	return "Renews the lease of a secret"
}

func (c *LeaseRenewCommand) Help() string {
	helpText := `
Usage: vault lease renew [options] ID

  Renews the lease on a secret, extending the time that it can be used before
  it is revoked by Vault.

  Every secret in Vault has a lease associated with it. If the owner of the
  secret wants to use it longer than the lease, then it must be renewed.
  Renewing the lease does not change the contents of the secret. The ID is the
  full path lease ID.

  Renew a secret:

      $ vault lease renew database/creds/readonly/2f6a614c...

  Lease renewal will fail if the secret is not renewable, the secret has already
  been revoked, or if the secret has already reached its maximum TTL.

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *LeaseRenewCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")

	f.DurationVar(&DurationVar{
		Name:       "increment",
		Target:     &c.flagIncrement,
		Default:    0,
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "Request a specific increment in seconds. Vault is not required " +
			"to honor this request.",
	})

	return set
}

func (c *LeaseRenewCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *LeaseRenewCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *LeaseRenewCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	leaseID := ""
	increment := c.flagIncrement

	args = f.Args()
	switch len(args) {
	case 0:
		c.UI.Error("Missing ID!")
		return 1
	case 1:
		leaseID = strings.TrimSpace(args[0])
	default:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1-2, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	secret, err := client.Sys().Renew(leaseID, truncateToSeconds(increment))
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error renewing %s: %s", leaseID, err))
		return 2
	}

	return OutputSecret(c.UI, secret)
}
