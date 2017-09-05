package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*MountTuneCommand)(nil)
var _ cli.CommandAutocomplete = (*MountTuneCommand)(nil)

// MountTuneCommand is a Command that remounts a mounted secret backend
// to a new endpoint.
type MountTuneCommand struct {
	*BaseCommand

	flagDefaultLeaseTTL time.Duration
	flagMaxLeaseTTL     time.Duration
}

func (c *MountTuneCommand) Synopsis() string {
	return "Tunes an existing mount's configuration"
}

func (c *MountTuneCommand) Help() string {
	helpText := `
Usage: vault mount-tune [options] PATH

  Tune the configuration options for a mounted secret backend at the given
  path. The argument corresponds to the PATH of the mount, not the TYPE!

  Tune the default lease for the PKI secret backend:

      $ vault mount-tune -default-lease-ttl=72h pki/

  For a full list of examples and paths, please see the documentation that
  corresponds to the secret backend in use.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *MountTuneCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.DurationVar(&DurationVar{
		Name:       "default-lease-ttl",
		Target:     &c.flagDefaultLeaseTTL,
		Default:    0,
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "The default lease TTL for this backend. If unspecified, this " +
			"defaults to the Vault server's globally configured default lease TTL, " +
			"or a previously configured value for the backend.",
	})

	f.DurationVar(&DurationVar{
		Name:       "max-lease-ttl",
		Target:     &c.flagMaxLeaseTTL,
		Default:    0,
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "The maximum lease TTL for this backend. If unspecified, this " +
			"defaults to the Vault server's globally configured maximum lease TTL, " +
			"or a previously configured value for the backend.",
	})

	return set
}

func (c *MountTuneCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultMounts()
}

func (c *MountTuneCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *MountTuneCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	mountPath, remaining, err := extractPath(args)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if len(remaining) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// Append a trailing slash to indicate it's a path in output
	mountPath = ensureTrailingSlash(mountPath)

	mountConfig := api.MountConfigInput{
		DefaultLeaseTTL: c.flagDefaultLeaseTTL.String(),
		MaxLeaseTTL:     c.flagMaxLeaseTTL.String(),
	}

	if err := client.Sys().TuneMount(mountPath, mountConfig); err != nil {
		c.UI.Error(fmt.Sprintf("Error tuning mount %s: %s", mountPath, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Tuned the mount at: %s", mountPath))
	return 0
}
