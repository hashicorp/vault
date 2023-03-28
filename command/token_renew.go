// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*TokenRenewCommand)(nil)
	_ cli.CommandAutocomplete = (*TokenRenewCommand)(nil)
)

type TokenRenewCommand struct {
	*BaseCommand

	flagAccessor  bool
	flagIncrement time.Duration
}

func (c *TokenRenewCommand) Synopsis() string {
	return "Renew a token lease"
}

func (c *TokenRenewCommand) Help() string {
	helpText := `
Usage: vault token renew [options] [TOKEN]

  Renews a token's lease, extending the amount of time it can be used. If a
  TOKEN is not provided, the locally authenticated token is used. A token
  accessor can be used as well. Lease renewal will fail if the token is not
  renewable, the token has already been revoked, or if the token has already
  reached its maximum TTL.

  Renew a token (this uses the /auth/token/renew endpoint and permission):

      $ vault token renew 96ddf4bc-d217-f3ba-f9bd-017055595017

  Renew the currently authenticated token (this uses the /auth/token/renew-self
  endpoint and permission):

      $ vault token renew

  Renew a token requesting a specific increment value:

      $ vault token renew -increment=30m 96ddf4bc-d217-f3ba-f9bd-017055595017

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TokenRenewCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:       "accessor",
		Target:     &c.flagAccessor,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage: "Treat the argument as an accessor instead of a token. When " +
			"this option is selected, the output will NOT include the token.",
	})

	f.DurationVar(&DurationVar{
		Name:       "increment",
		Aliases:    []string{"i"},
		Target:     &c.flagIncrement,
		Default:    0,
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "Request a specific increment for renewal. This increment may " +
			"not be honored, for instance in the case of periodic tokens. If not " +
			"supplied, Vault will use the default TTL. This is specified as a " +
			"numeric string with suffix like \"30s\" or \"5m\".",
	})

	return set
}

func (c *TokenRenewCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *TokenRenewCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TokenRenewCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	token := ""
	increment := c.flagIncrement

	args = f.Args()
	switch len(args) {
	case 0:
		// Use the local token
	case 1:
		token = strings.TrimSpace(args[0])
	default:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	var secret *api.Secret
	inc := truncateToSeconds(increment)
	switch {
	case token == "":
		secret, err = client.Auth().Token().RenewSelf(inc)
	case c.flagAccessor:
		secret, err = client.Auth().Token().RenewAccessor(token, inc)
	default:
		secret, err = client.Auth().Token().Renew(token, inc)
	}
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error renewing token: %s", err))
		return 2
	}

	return OutputSecret(c.UI, secret)
}
