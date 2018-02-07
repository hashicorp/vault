package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*SecretsTuneCommand)(nil)
var _ cli.CommandAutocomplete = (*SecretsTuneCommand)(nil)

type SecretsTuneCommand struct {
	*BaseCommand

	flagDefaultLeaseTTL time.Duration
	flagMaxLeaseTTL     time.Duration
}

func (c *SecretsTuneCommand) Synopsis() string {
	return "Tune a secrets engine configuration"
}

func (c *SecretsTuneCommand) Help() string {
	helpText := `
Usage: vault secrets tune [options] PATH

  Tunes the configuration options for the secrets engine at the given PATH.
  The argument corresponds to the PATH where the secrets engine is enabled,
  not the TYPE!

  Tune the default lease for the PKI secrets engine:

      $ vault secrets tune -default-lease-ttl=72h pki/

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *SecretsTuneCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.DurationVar(&DurationVar{
		Name:       "default-lease-ttl",
		Target:     &c.flagDefaultLeaseTTL,
		Default:    0,
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "The default lease TTL for this secrets engine. If unspecified, " +
			"this defaults to the Vault server's globally configured default lease " +
			"TTL, or a previously configured value for the secrets engine.",
	})

	f.DurationVar(&DurationVar{
		Name:       "max-lease-ttl",
		Target:     &c.flagMaxLeaseTTL,
		Default:    0,
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "The maximum lease TTL for this secrets engine. If unspecified, " +
			"this defaults to the Vault server's globally configured maximum lease " +
			"TTL, or a previously configured value for the secrets engine.",
	})

	return set
}

func (c *SecretsTuneCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultMounts()
}

func (c *SecretsTuneCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *SecretsTuneCommand) Run(args []string) int {
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

	// Append a trailing slash to indicate it's a path in output
	mountPath := ensureTrailingSlash(sanitizePath(args[0]))

	if err := client.Sys().TuneMount(mountPath, api.MountConfigInput{
		DefaultLeaseTTL: ttlToAPI(c.flagDefaultLeaseTTL),
		MaxLeaseTTL:     ttlToAPI(c.flagMaxLeaseTTL),
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error tuning secrets engine %s: %s", mountPath, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Tuned the secrets engine at: %s", mountPath))
	return 0
}
