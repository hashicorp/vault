package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*AuthTuneCommand)(nil)
var _ cli.CommandAutocomplete = (*AuthTuneCommand)(nil)

type AuthTuneCommand struct {
	*BaseCommand

	flagDefaultLeaseTTL time.Duration
	flagMaxLeaseTTL     time.Duration
}

func (c *AuthTuneCommand) Synopsis() string {
	return "Tunes an auth method configuration"
}

func (c *AuthTuneCommand) Help() string {
	helpText := `
Usage: vault auth tune [options] PATH

  Tunes the configuration options for the auth method at the given PATH. The
  argument corresponds to the PATH where the auth method is enabled, not the
  TYPE!

  Tune the default lease for the github auth method:

      $ vault auth tune -default-lease-ttl=72h github/

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuthTuneCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.DurationVar(&DurationVar{
		Name:       "default-lease-ttl",
		Target:     &c.flagDefaultLeaseTTL,
		Default:    0,
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "The default lease TTL for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured default lease TTL, " +
			"or a previously configured value for the auth method.",
	})

	f.DurationVar(&DurationVar{
		Name:       "max-lease-ttl",
		Target:     &c.flagMaxLeaseTTL,
		Default:    0,
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "The maximum lease TTL for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured maximum lease TTL, " +
			"or a previously configured value for the auth method.",
	})

	return set
}

func (c *AuthTuneCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultAuths()
}

func (c *AuthTuneCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthTuneCommand) Run(args []string) int {
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

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// Append /auth (since that's where auths live) and a trailing slash to
	// indicate it's a path in output
	mountPath := ensureTrailingSlash(sanitizePath(args[0]))

	if err := client.Sys().TuneMount("/auth/"+mountPath, api.MountConfigInput{
		DefaultLeaseTTL: ttlToAPI(c.flagDefaultLeaseTTL),
		MaxLeaseTTL:     ttlToAPI(c.flagMaxLeaseTTL),
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error tuning auth method %s: %s", mountPath, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Tuned the auth method at: %s", mountPath))
	return 0
}
