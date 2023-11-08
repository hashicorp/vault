// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*SecretsEnableCommand)(nil)
	_ cli.CommandAutocomplete = (*SecretsEnableCommand)(nil)
)

type SecretsEnableCommand struct {
	*BaseCommand

	flagDescription               string
	flagPath                      string
	flagDefaultLeaseTTL           time.Duration
	flagMaxLeaseTTL               time.Duration
	flagAuditNonHMACRequestKeys   []string
	flagAuditNonHMACResponseKeys  []string
	flagListingVisibility         string
	flagPassthroughRequestHeaders []string
	flagAllowedResponseHeaders    []string
	flagForceNoCache              bool
	flagPluginName                string
	flagPluginVersion             string
	flagOptions                   map[string]string
	flagLocal                     bool
	flagSealWrap                  bool
	flagExternalEntropyAccess     bool
	flagVersion                   int
	flagAllowedManagedKeys        []string
}

func (c *SecretsEnableCommand) Synopsis() string {
	return "Enable a secrets engine"
}

func (c *SecretsEnableCommand) Help() string {
	helpText := `
Usage: vault secrets enable [options] TYPE

  Enables a secrets engine. By default, secrets engines are enabled at the path
  corresponding to their TYPE, but users can customize the path using the
  -path option.

  Once enabled, Vault will route all requests which begin with the path to the
  secrets engine.

  Enable the AWS secrets engine at aws/:

      $ vault secrets enable aws

  Enable the SSH secrets engine at ssh-prod/:

      $ vault secrets enable -path=ssh-prod ssh

  Enable the database secrets engine with an explicit maximum TTL of 30m:

      $ vault secrets enable -max-lease-ttl=30m database

  Enable a custom plugin (after it is registered in the plugin registry):

      $ vault secrets enable -path=my-secrets -plugin-name=my-plugin plugin

  OR (preferred way):

      $ vault secrets enable -path=my-secrets my-plugin

  For a full list of secrets engines and examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *SecretsEnableCommand) Flags() *FlagSets {
	set := c.FlagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "description",
		Target:     &c.flagDescription,
		Completion: complete.PredictAnything,
		Usage:      "Human-friendly description for the purpose of this engine.",
	})

	f.StringVar(&StringVar{
		Name:       "path",
		Target:     &c.flagPath,
		Default:    "", // The default is complex, so we have to manually document
		Completion: complete.PredictAnything,
		Usage: "Place where the secrets engine will be accessible. This must be " +
			"unique cross all secrets engines. This defaults to the \"type\" of the " +
			"secrets engine.",
	})

	f.DurationVar(&DurationVar{
		Name:       "default-lease-ttl",
		Target:     &c.flagDefaultLeaseTTL,
		Completion: complete.PredictAnything,
		Usage: "The default lease TTL for this secrets engine. If unspecified, " +
			"this defaults to the Vault server's globally configured default lease " +
			"TTL.",
	})

	f.DurationVar(&DurationVar{
		Name:       "max-lease-ttl",
		Target:     &c.flagMaxLeaseTTL,
		Completion: complete.PredictAnything,
		Usage: "The maximum lease TTL for this secrets engine. If unspecified, " +
			"this defaults to the Vault server's globally configured maximum lease " +
			"TTL.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNameAuditNonHMACRequestKeys,
		Target: &c.flagAuditNonHMACRequestKeys,
		Usage: "Key that will not be HMAC'd by audit devices in the request data object. " +
			"To specify multiple values, specify this flag multiple times.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNameAuditNonHMACResponseKeys,
		Target: &c.flagAuditNonHMACResponseKeys,
		Usage: "Key that will not be HMAC'd by audit devices in the response data object. " +
			"To specify multiple values, specify this flag multiple times.",
	})

	f.StringVar(&StringVar{
		Name:   flagNameListingVisibility,
		Target: &c.flagListingVisibility,
		Usage:  "Determines the visibility of the mount in the UI-specific listing endpoint.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNamePassthroughRequestHeaders,
		Target: &c.flagPassthroughRequestHeaders,
		Usage: "Request header value that will be sent to the plugins. To specify multiple " +
			"values, specify this flag multiple times.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNameAllowedResponseHeaders,
		Target: &c.flagAllowedResponseHeaders,
		Usage: "Response header value that plugins will be allowed to set. To specify multiple " +
			"values, specify this flag multiple times.",
	})

	f.BoolVar(&BoolVar{
		Name:    "force-no-cache",
		Target:  &c.flagForceNoCache,
		Default: false,
		Usage: "Force the secrets engine to disable caching. If unspecified, this " +
			"defaults to the Vault server's globally configured cache settings. " +
			"This does not affect caching of the underlying encrypted data storage.",
	})

	f.StringVar(&StringVar{
		Name:       "plugin-name",
		Target:     &c.flagPluginName,
		Completion: c.PredictVaultPlugins(api.PluginTypeSecrets, api.PluginTypeDatabase),
		Usage: "Name of the secrets engine plugin. This plugin name must already " +
			"exist in Vault's plugin catalog.",
	})

	f.StringVar(&StringVar{
		Name:    flagNamePluginVersion,
		Target:  &c.flagPluginVersion,
		Default: "",
		Usage:   "Select the semantic version of the plugin to enable.",
	})

	f.StringMapVar(&StringMapVar{
		Name:       "options",
		Target:     &c.flagOptions,
		Completion: complete.PredictAnything,
		Usage: "Key-value pair provided as key=value for the mount options. " +
			"This can be specified multiple times.",
	})

	f.BoolVar(&BoolVar{
		Name:    "local",
		Target:  &c.flagLocal,
		Default: false,
		Usage: "Mark the secrets engine as local-only. Local engines are not " +
			"replicated or removed by replication.",
	})

	f.BoolVar(&BoolVar{
		Name:    "seal-wrap",
		Target:  &c.flagSealWrap,
		Default: false,
		Usage:   "Enable seal wrapping of critical values in the secrets engine.",
	})

	f.BoolVar(&BoolVar{
		Name:    "external-entropy-access",
		Target:  &c.flagExternalEntropyAccess,
		Default: false,
		Usage:   "Enable secrets engine to access Vault's external entropy source.",
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

	return set
}

func (c *SecretsEnableCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultAvailableMounts()
}

func (c *SecretsEnableCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *SecretsEnableCommand) Run(args []string) int {
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

	// Get the engine type type (first arg)
	engineType := strings.TrimSpace(args[0])
	if engineType == "plugin" {
		engineType = c.flagPluginName
	}

	// If no path is specified, we default the path to the backend type
	// or use the plugin name if it's a plugin backend
	mountPath := c.flagPath
	if mountPath == "" {
		if engineType == "plugin" {
			mountPath = c.flagPluginName
		} else {
			mountPath = engineType
		}
	}

	if c.flagVersion > 0 {
		if c.flagOptions == nil {
			c.flagOptions = make(map[string]string)
		}
		c.flagOptions["version"] = strconv.Itoa(c.flagVersion)
	}

	// Append a trailing slash to indicate it's a path in output
	mountPath = ensureTrailingSlash(mountPath)

	// Build mount input
	mountInput := &api.MountInput{
		Type:                  engineType,
		Description:           c.flagDescription,
		Local:                 c.flagLocal,
		SealWrap:              c.flagSealWrap,
		ExternalEntropyAccess: c.flagExternalEntropyAccess,
		Config: api.MountConfigInput{
			DefaultLeaseTTL: c.flagDefaultLeaseTTL.String(),
			MaxLeaseTTL:     c.flagMaxLeaseTTL.String(),
			ForceNoCache:    c.flagForceNoCache,
		},
		Options: c.flagOptions,
	}

	// Set these values only if they are provided in the CLI
	f.Visit(func(fl *flag.Flag) {
		if fl.Name == flagNameAuditNonHMACRequestKeys {
			mountInput.Config.AuditNonHMACRequestKeys = c.flagAuditNonHMACRequestKeys
		}

		if fl.Name == flagNameAuditNonHMACResponseKeys {
			mountInput.Config.AuditNonHMACResponseKeys = c.flagAuditNonHMACResponseKeys
		}

		if fl.Name == flagNameListingVisibility {
			mountInput.Config.ListingVisibility = c.flagListingVisibility
		}

		if fl.Name == flagNamePassthroughRequestHeaders {
			mountInput.Config.PassthroughRequestHeaders = c.flagPassthroughRequestHeaders
		}

		if fl.Name == flagNameAllowedResponseHeaders {
			mountInput.Config.AllowedResponseHeaders = c.flagAllowedResponseHeaders
		}

		if fl.Name == flagNameAllowedManagedKeys {
			mountInput.Config.AllowedManagedKeys = c.flagAllowedManagedKeys
		}

		if fl.Name == flagNamePluginVersion {
			mountInput.Config.PluginVersion = c.flagPluginVersion
		}
	})

	if err := client.Sys().Mount(mountPath, mountInput); err != nil {
		c.UI.Error(fmt.Sprintf("Error enabling: %s", err))
		return 2
	}

	thing := engineType + " secrets engine"
	if engineType == "plugin" {
		thing = c.flagPluginName + " plugin"
	}
	if c.flagPluginVersion != "" {
		thing += " version " + c.flagPluginVersion
	}
	c.UI.Output(fmt.Sprintf("Success! Enabled the %s at: %s", thing, mountPath))
	return 0
}
