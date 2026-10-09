// Copyright IBM Corp. 2016, 2025
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
	base "github.com/hashicorp/vault/command/base"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*AuthTuneCommand)(nil)
	_ cli.CommandAutocomplete = (*AuthTuneCommand)(nil)
)

type AuthTuneCommand struct {
	*base.BaseCommand

	flagAuditNonHMACRequestKeys         []string
	flagAuditNonHMACResponseKeys        []string
	flagDefaultLeaseTTL                 time.Duration
	flagDescription                     string
	flagListingVisibility               string
	flagMaxLeaseTTL                     time.Duration
	flagPassthroughRequestHeaders       []string
	flagAllowedResponseHeaders          []string
	flagOptions                         map[string]string
	flagTokenType                       string
	flagVersion                         int
	flagPluginVersion                   string
	flagOverridePinnedVersion           base.BoolPtr
	flagUserLockoutThreshold            uint
	flagUserLockoutDuration             time.Duration
	flagUserLockoutCounterResetDuration time.Duration
	flagUserLockoutDisable              bool
	flagIdentityTokenKey                string
	flagTrimRequestTrailingSlashes      base.BoolPtr
	flagSealWrap                        base.BoolPtr
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

func (c *AuthTuneCommand) Flags() *base.FlagSets {
	set := c.FlagSetForBit(base.FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringSliceVar(&base.StringSliceVar{
		Name:   base.FlagNameAuditNonHMACRequestKeys,
		Target: &c.flagAuditNonHMACRequestKeys,
		Usage: "Key that will not be HMAC'd by audit devices in the request data " +
			"object. To specify multiple values, specify this flag multiple times.",
	})

	f.StringSliceVar(&base.StringSliceVar{
		Name:   base.FlagNameAuditNonHMACResponseKeys,
		Target: &c.flagAuditNonHMACResponseKeys,
		Usage: "Key that will not be HMAC'd by audit devices in the response data " +
			"object. To specify multiple values, specify this flag multiple times.",
	})

	f.DurationVar(&base.DurationVar{
		Name:       "default-lease-ttl",
		Target:     &c.flagDefaultLeaseTTL,
		Default:    0,
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "The default lease TTL for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured default lease TTL, " +
			"or a previously configured value for the auth method.",
	})

	f.StringVar(&base.StringVar{
		Name:   base.FlagNameDescription,
		Target: &c.flagDescription,
		Usage: "Human-friendly description of the this auth method. This overrides " +
			"the current stored value, if any.",
	})

	f.StringVar(&base.StringVar{
		Name:   base.FlagNameListingVisibility,
		Target: &c.flagListingVisibility,
		Usage: "Determines the visibility of the mount in the UI-specific listing " +
			"endpoint.",
	})

	f.DurationVar(&base.DurationVar{
		Name:       "max-lease-ttl",
		Target:     &c.flagMaxLeaseTTL,
		Default:    0,
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "The maximum lease TTL for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured maximum lease TTL, " +
			"or a previously configured value for the auth method.",
	})

	f.StringSliceVar(&base.StringSliceVar{
		Name:   base.FlagNamePassthroughRequestHeaders,
		Target: &c.flagPassthroughRequestHeaders,
		Usage: "Request header value that will be sent to the plugin. To specify " +
			"multiple values, specify this flag multiple times.",
	})

	f.StringSliceVar(&base.StringSliceVar{
		Name:   base.FlagNameAllowedResponseHeaders,
		Target: &c.flagAllowedResponseHeaders,
		Usage: "Response header value that plugins will be allowed to set. To specify " +
			"multiple values, specify this flag multiple times.",
	})

	f.StringMapVar(&base.StringMapVar{
		Name:       "options",
		Target:     &c.flagOptions,
		Completion: complete.PredictAnything,
		Usage: "Key-value pair provided as key=value for the mount options. " +
			"This can be specified multiple times.",
	})

	f.StringVar(&base.StringVar{
		Name:   base.FlagNameTokenType,
		Target: &c.flagTokenType,
		Usage:  "Sets a forced token type for the mount.",
	})

	f.IntVar(&base.IntVar{
		Name:    "version",
		Target:  &c.flagVersion,
		Default: 0,
		Usage:   "Select the version of the auth method to run. Not supported by all auth methods.",
	})

	f.UintVar(&base.UintVar{
		Name:   base.FlagNameUserLockoutThreshold,
		Target: &c.flagUserLockoutThreshold,
		Usage: "The threshold for user lockout for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured user lockout threshold, " +
			"or a previously configured value for the auth method.",
	})

	f.DurationVar(&base.DurationVar{
		Name:       base.FlagNameUserLockoutDuration,
		Target:     &c.flagUserLockoutDuration,
		Completion: complete.PredictAnything,
		Usage: "The user lockout duration for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured user lockout duration, " +
			"or a previously configured value for the auth method.",
	})

	f.DurationVar(&base.DurationVar{
		Name:       base.FlagNameUserLockoutCounterResetDuration,
		Target:     &c.flagUserLockoutCounterResetDuration,
		Completion: complete.PredictAnything,
		Usage: "The user lockout counter reset duration for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured user lockout counter reset duration, " +
			"or a previously configured value for the auth method.",
	})

	f.BoolVar(&base.BoolVar{
		Name:    base.FlagNameUserLockoutDisable,
		Target:  &c.flagUserLockoutDisable,
		Default: false,
		Usage: "Disable user lockout for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured user lockout disable, " +
			"or a previously configured value for the auth method.",
	})

	f.StringVar(&base.StringVar{
		Name:    base.FlagNamePluginVersion,
		Target:  &c.flagPluginVersion,
		Default: "",
		Usage: "Select the semantic version of the plugin to run. The new version must be registered in " +
			"the plugin catalog, and will not start running until the plugin is reloaded.",
	})

	f.BoolPtrVar(&base.BoolPtrVar{
		Name:   base.FlagNameOverridePinnedVersion,
		Target: &c.flagOverridePinnedVersion,
		Usage:  "Whether to override the pinned version for this mount",
	})

	f.BoolPtrVar(&base.BoolPtrVar{
		Name:   base.FlagNameTrimRequestTrailingSlashes,
		Target: &c.flagTrimRequestTrailingSlashes,
		Usage:  "Whether to trim trailing slashes for incoming requests to this mount",
	})

	f.StringVar(&base.StringVar{
		Name:    base.FlagNameIdentityTokenKey,
		Target:  &c.flagIdentityTokenKey,
		Default: "default",
		Usage:   "Select the key used to sign plugin identity tokens.",
	})

	f.BoolPtrVar(&base.BoolPtrVar{
		Name:   base.FlagNameSealWrap,
		Target: &c.flagSealWrap,
		Usage:  "Enterprise only. Whether critical security parameters (CSPs) are seal wrapped in this mount",
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

	if c.flagVersion > 0 {
		if c.flagOptions == nil {
			c.flagOptions = make(map[string]string)
		}
		c.flagOptions["version"] = strconv.Itoa(c.flagVersion)
	}

	defaultLeaseTtl := base.TtlToAPI(c.flagDefaultLeaseTTL)
	maxLeaseTtl := base.TtlToAPI(c.flagMaxLeaseTTL)
	options := c.flagOptions
	tuneMountInput := api.TuneMountConfigInput{
		DefaultLeaseTTL: &defaultLeaseTtl,
		MaxLeaseTTL:     &maxLeaseTtl,
		Options:         &options,
	}

	userLockoutConfig := api.TuneUserLockoutConfigInput{}
	userLockoutConfigSet := false

	// Set these values only if they are provided in the CLI
	f.Visit(func(fl *flag.Flag) {
		if fl.Name == base.FlagNameAuditNonHMACRequestKeys {
			if len(c.flagAuditNonHMACRequestKeys) == 1 && c.flagAuditNonHMACRequestKeys[0] == "" {
				emptyList := []string{}
				tuneMountInput.AuditNonHMACRequestKeys = &emptyList
			} else {
				tuneMountInput.AuditNonHMACRequestKeys = &c.flagAuditNonHMACRequestKeys
			}
		}

		if fl.Name == base.FlagNameAuditNonHMACResponseKeys {
			if len(c.flagAuditNonHMACResponseKeys) == 1 && c.flagAuditNonHMACResponseKeys[0] == "" {
				emptyList := []string{}
				tuneMountInput.AuditNonHMACResponseKeys = &emptyList
			} else {
				tuneMountInput.AuditNonHMACResponseKeys = &c.flagAuditNonHMACResponseKeys
			}
		}

		if fl.Name == base.FlagNameDescription {
			tuneMountInput.Description = &c.flagDescription
		}

		if fl.Name == base.FlagNameListingVisibility {
			tuneMountInput.ListingVisibility = &c.flagListingVisibility
		}

		if fl.Name == base.FlagNamePassthroughRequestHeaders {
			if len(c.flagPassthroughRequestHeaders) == 1 && c.flagPassthroughRequestHeaders[0] == "" {
				emptyList := []string{}
				tuneMountInput.PassthroughRequestHeaders = &emptyList
			} else {
				tuneMountInput.PassthroughRequestHeaders = &c.flagPassthroughRequestHeaders
			}
		}

		if fl.Name == base.FlagNameAllowedResponseHeaders {
			if len(c.flagAllowedResponseHeaders) == 1 && c.flagAllowedResponseHeaders[0] == "" {
				emptyList := []string{}
				tuneMountInput.AllowedResponseHeaders = &emptyList
			} else {
				tuneMountInput.AllowedResponseHeaders = &c.flagAllowedResponseHeaders
			}
		}

		if fl.Name == base.FlagNameTokenType {
			tuneMountInput.TokenType = &c.flagTokenType
		}

		switch fl.Name {
		case base.FlagNameUserLockoutThreshold, base.FlagNameUserLockoutDuration, base.FlagNameUserLockoutCounterResetDuration, base.FlagNameUserLockoutDisable:
			userLockoutConfigSet = true
		}
		if fl.Name == base.FlagNameUserLockoutThreshold {
			lockoutThreshold := strconv.FormatUint(uint64(c.flagUserLockoutThreshold), 10)
			userLockoutConfig.LockoutThreshold = &lockoutThreshold
		}
		if fl.Name == base.FlagNameUserLockoutDuration {
			lockoutDuration := base.TtlToAPI(c.flagUserLockoutDuration)
			userLockoutConfig.LockoutDuration = &lockoutDuration
		}
		if fl.Name == base.FlagNameUserLockoutCounterResetDuration {
			lockoutCounterResetDuration := base.TtlToAPI(c.flagUserLockoutCounterResetDuration)
			userLockoutConfig.LockoutCounterResetDuration = &lockoutCounterResetDuration
		}
		if fl.Name == base.FlagNameUserLockoutDisable {
			userLockoutConfig.DisableLockout = &c.flagUserLockoutDisable
		}

		if fl.Name == base.FlagNamePluginVersion {
			tuneMountInput.PluginVersion = &c.flagPluginVersion
		}

		if fl.Name == base.FlagNameSealWrap {
			_sealVal := c.flagSealWrap.Get()
			tuneMountInput.SealWrap = &_sealVal
		}

		if fl.Name == base.FlagNameOverridePinnedVersion && c.flagOverridePinnedVersion.IsSet() {
			val := c.flagOverridePinnedVersion.Get()
			tuneMountInput.OverridePinnedVersion = &val
		}

		if fl.Name == base.FlagNameIdentityTokenKey {
			tuneMountInput.IdentityTokenKey = &c.flagIdentityTokenKey
		}

		if fl.Name == base.FlagNameTrimRequestTrailingSlashes && c.flagTrimRequestTrailingSlashes.IsSet() {
			val := c.flagTrimRequestTrailingSlashes.Get()
			tuneMountInput.TrimRequestTrailingSlashes = &val
		}
	})

	if userLockoutConfigSet {
		tuneMountInput.UserLockoutConfig = &userLockoutConfig
	}

	// Append /auth (since that's where auths live) and a trailing slash to
	// indicate it's a path in output
	mountPath := base.EnsureTrailingSlash(base.SanitizePath(args[0]))

	if err := client.Sys().TuneMountAllowNil("/auth/"+mountPath, tuneMountInput); err != nil {
		c.UI.Error(fmt.Sprintf("Error tuning auth method %s: %s", mountPath, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Tuned the auth method at: %s", mountPath))
	return 0
}
