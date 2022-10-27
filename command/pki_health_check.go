package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

const (
	pkiRetOK int = iota
	pkiRetUsage
	pkiRetInformational
	pkiRetWarning
	pkiRetCritical
	pkiRetInvalidVersion
	pkiRetInsufficientPermissions
)

const (
	oneDay   = 24 * time.Hour
	oneMonth = 30 * oneDay
	oneYear  = 365 * oneDay
)

var (
	_ cli.Command             = (*PKIHealthCheckCommand)(nil)
	_ cli.CommandAutocomplete = (*PKIHealthCheckCommand)(nil)
)

type PKIHealthCheckCommand struct {
	*BaseCommand

	flagConfig          string
	flagReturnIndicator string
	flagDefaultDisabled bool
	flagList            bool
}

func (c *PKIHealthCheckCommand) Synopsis() string {
	return "Check PKI Secrets Engine health and operational status"
}

func (c *PKIHealthCheckCommand) Help() string {
	helpText := `
Usage: vault pki health-check [options] MOUNT

  Reports status of the specified mount against best practices and pending
  failures. This is an informative command and not all recommendations will
  apply to all mounts; consider using a configuration file to tune the
  executed health checks.

  To check the pki-root mount with default configuration:

      $ vault pki health-check pki-root

  To specify a configuration:

      $ vault pki health-check -health-config=mycorp-root.json /pki-root

  Return codes indicate failure type:

      0 - Everything is good.
      1 - Usage error (check CLI parameters).
	  2 - Informational message from a health check.
	  3 - Warning message from a health check.
	  4 - Critical message from a health check.
	  5 - A version mismatch between health check and Vault Server occurred,
	      preventing one or more health checks from being run.
      6 - A permission denied message was returned from Vault Server for
	      one or more health checks.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKIHealthCheckCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:    "health-config",
		Target:  &c.flagConfig,
		Default: "",
		EnvVar:  "",
		Usage:   "Path to JSON configuration file to modify health check execution and parameters.",
	})

	f.StringVar(&StringVar{
		Name:    "return-indicator",
		Target:  &c.flagReturnIndicator,
		Default: "default",
		EnvVar:  "",
		Usage: `Behavior of the return value:
 - permission, for exiting with a non-zero code when the tool lacks
               permissions or has a version mismatch with the server;
 - critical, for exiting with a non-zero code when a check returns a
             critical status in addition to the above;
 - warning, for exiting with a non-zero status when a check returns a
            warning status in addition to the above;
 - informational, for exiting with a non-zero status when a check returns
                  an informational status in addition to the above;
 - default, for the default behavior based on severity of message and
            only returning a zero exit status when all checks have passed
			and no execution errors have occurred.
		`,
	})

	f.BoolVar(&BoolVar{
		Name:    "default-disabled",
		Target:  &c.flagDefaultDisabled,
		Default: false,
		EnvVar:  "",
		Usage: `When specified, results in all health checks being disabled by
default unless enabled by the configuration file explicitly.`,
	})

	f.BoolVar(&BoolVar{
		Name:    "list",
		Target:  &c.flagList,
		Default: false,
		EnvVar:  "",
		Usage: `When specified, no health checks are run, but all known health
checks are printed.`,
	})

	return set
}

func (c *PKIHealthCheckCommand) AutocompleteArgs() complete.Predictor {
	// Return an anything predictor here, similar to `vault write`. We
	// don't know what values are valid for the mount path.
	return complete.PredictAnything
}

func (c *PKIHealthCheckCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PKIHealthCheckCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return pkiRetUsage
	}

	args = f.Args()
	if len(args) < 1 {
		c.UI.Error("Not enough arguments (expected mount path, got nothing)")
		return pkiRetUsage
	} else if len(args) > 1 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected only mount path, got %d arguments)", len(args)))
		for _, arg := range args {
			if strings.HasPrefix(arg, "-") {
				c.UI.Warn(fmt.Sprintf("Options (%v) must be specified before positional arguments (%v)", arg, args[0]))
				break
			}
		}
		return pkiRetUsage
	}

	// mount := sanitizePath(args[0])

	return 0
}
