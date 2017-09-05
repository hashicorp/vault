package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*RenewCommand)(nil)
var _ cli.CommandAutocomplete = (*RenewCommand)(nil)

// RenewCommand is a Command that mounts a new mount.
type RenewCommand struct {
	*BaseCommand

	flagIncrement time.Duration
}

func (c *RenewCommand) Synopsis() string {
	return "Renews the lease of a secret"
}

func (c *RenewCommand) Help() string {
	helpText := `
Usage: vault renew [options] ID

  Renews the lease on a secret, extending the time that it can be used before
  it is revoked by Vault.

  Every secret in Vault has a lease associated with it. If the owner of the
  secret wants to use it longer than the lease, then it must be renewed.
  Renewing the lease does not change the contents of the secret. The ID is the
  full path lease ID.

  Renew a secret:

      $ vault renew database/creds/readonly/2f6a614c-4aa2-7b19-24b9-ad944a8d4de6

  Lease renewal will fail if the secret is not renewable, the secret has already
  been revoked, or if the secret has already reached its maximum TTL.

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *RenewCommand) Flags() *FlagSets {
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

func (c *RenewCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *RenewCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RenewCommand) Run(args []string) int {
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
	case 2:
		// Deprecation
		// TODO: remove in 0.9.0
		c.UI.Warn(wrapAtLength(
			"WARNING! Specifying INCREMENT as a second argument is deprecated. " +
				"Please use -increment instead. This will be removed in the next " +
				"major release of Vault."))

		leaseID = strings.TrimSpace(args[0])
		parsed, err := time.ParseDuration(appendDurationSuffix(args[1]))
		if err != nil {
			c.UI.Error(fmt.Sprintf("Invalid increment: %s", err))
			return 1
		}
		increment = parsed
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

	return OutputSecret(c.UI, c.flagFormat, secret)
}
