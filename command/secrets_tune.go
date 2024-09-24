// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*SecretsTuneCommand)(nil)
	_ cli.CommandAutocomplete = (*SecretsTuneCommand)(nil)
)

type SecretsTuneCommand struct {
	*BaseCommand

	flagAuditNonHMACRequestKeys   []string
	flagAuditNonHMACResponseKeys  []string
	flagDefaultLeaseTTL           time.Duration
	flagDescription               string
	flagListingVisibility         string
	flagMaxLeaseTTL               time.Duration
	flagPassthroughRequestHeaders []string
	flagAllowedResponseHeaders    []string
	flagOptions                   map[string]string
	flagVersion                   int
	flagPluginVersion             string
	flagAllowedManagedKeys        []string
	flagDelegatedAuthAccessors    []string
	flagIdentityTokenKey          string
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

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNameAuditNonHMACRequestKeys,
		Target: &c.flagAuditNonHMACRequestKeys,
		Usage: "Key that will not be HMAC'd by audit devices in the request data " +
			"object. To specify multiple values, specify this flag multiple times.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNameAuditNonHMACResponseKeys,
		Target: &c.flagAuditNonHMACResponseKeys,
		Usage: "Key that will not be HMAC'd by audit devices in the response data " +
			"object. To specify multiple values, specify this flag multiple times.",
	})

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

	f.StringVar(&StringVar{
		Name:   flagNameDescription,
		Target: &c.flagDescription,
		Usage: "Human-friendly description of this secret engine. This overrides the " +
			"current stored value, if any.",
	})

	f.StringVar(&StringVar{
		Name:   flagNameListingVisibility,
		Target: &c.flagListingVisibility,
		Usage: "Determines the visibility of the mount in the UI-specific listing " +
			"endpoint.",
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

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNamePassthroughRequestHeaders,
		Target: &c.flagPassthroughRequestHeaders,
		Usage: "Request header value that will be sent to the plugin. To specify " +
			"multiple values, specify this flag multiple times.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNameAllowedResponseHeaders,
		Target: &c.flagAllowedResponseHeaders,
		Usage: "Response header value that plugins will be allowed to set. To " +
			"specify multiple values, specify this flag multiple times.",
	})

	f.StringMapVar(&StringMapVar{
		Name:       "options",
		Target:     &c.flagOptions,
		Completion: complete.PredictAnything,
		Usage: "Key-value pair provided as key=value for the mount options. " +
			"This can be specified multiple times.",
	})

	f.IntVar(&IntVar{
		Name:    "version",
		Target:  &c.flagVersion,
		Default: 0,
		Usage:   "Select the version of the engine to run. Not supported by all engines.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNameAllowedManagedKeys,
		Target: &c.flagAllowedManagedKeys,
		Usage: "Managed key name(s) that the mount in question is allowed to access. " +
			"Note that multiple keys may be specified by providing this option multiple times, " +
			"each time with 1 key.",
	})

	f.StringVar(&StringVar{
		Name:    flagNamePluginVersion,
		Target:  &c.flagPluginVersion,
		Default: "",
		Usage: "Select the semantic version of the plugin to run. The new version must be registered in " +
			"the plugin catalog, and will not start running until the plugin is reloaded.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNameDelegatedAuthAccessors,
		Target: &c.flagDelegatedAuthAccessors,
		Usage: "A list of permitted authentication accessors this backend can delegate authentication to. " +
			"Note that multiple values may be specified by providing this option multiple times, " +
			"each time with 1 accessor.",
	})

	f.StringVar(&StringVar{
		Name:    flagNameIdentityTokenKey,
		Target:  &c.flagIdentityTokenKey,
		Default: "default",
		Usage:   "Select the key used to sign plugin identity tokens.",
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

	if c.flagVersion > 0 {
		if c.flagOptions == nil {
			c.flagOptions = make(map[string]string)
		}
		c.flagOptions["version"] = strconv.Itoa(c.flagVersion)
	}

	// Append a trailing slash to indicate it's a path in output
	mountPath := ensureTrailingSlash(sanitizePath(args[0]))

	mountConfigInput := api.MountConfigInput{
		DefaultLeaseTTL: ttlToAPI(c.flagDefaultLeaseTTL),
		MaxLeaseTTL:     ttlToAPI(c.flagMaxLeaseTTL),
		Options:         c.flagOptions,
	}

	// Set these values only if they are provided in the CLI
	f.Visit(func(fl *flag.Flag) {
		if fl.Name == flagNameAuditNonHMACRequestKeys {
			mountConfigInput.AuditNonHMACRequestKeys = c.flagAuditNonHMACRequestKeys
		}

		if fl.Name == flagNameAuditNonHMACResponseKeys {
			mountConfigInput.AuditNonHMACResponseKeys = c.flagAuditNonHMACResponseKeys
		}

		if fl.Name == flagNameDescription {
			mountConfigInput.Description = &c.flagDescription
		}

		if fl.Name == flagNameListingVisibility {
			mountConfigInput.ListingVisibility = c.flagListingVisibility
		}

		if fl.Name == flagNamePassthroughRequestHeaders {
			mountConfigInput.PassthroughRequestHeaders = c.flagPassthroughRequestHeaders
		}

		if fl.Name == flagNameAllowedResponseHeaders {
			mountConfigInput.AllowedResponseHeaders = c.flagAllowedResponseHeaders
		}

		if fl.Name == flagNameAllowedManagedKeys {
			mountConfigInput.AllowedManagedKeys = c.flagAllowedManagedKeys
		}

		if fl.Name == flagNamePluginVersion {
			mountConfigInput.PluginVersion = c.flagPluginVersion
		}

		if fl.Name == flagNameDelegatedAuthAccessors {
			mountConfigInput.DelegatedAuthAccessors = c.flagDelegatedAuthAccessors
		}

		if fl.Name == flagNameIdentityTokenKey {
			mountConfigInput.IdentityTokenKey = c.flagIdentityTokenKey
		}
	})

	if err := client.Sys().TuneMount(mountPath, mountConfigInput); err != nil {
		c.UI.Error(fmt.Sprintf("Error tuning secrets engine %s: %s", mountPath, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Tuned the secrets engine at: %s", mountPath))
	return 0
}
