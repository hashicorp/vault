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
	_ cli.Command             = (*AuthTuneCommand)(nil)
	_ cli.CommandAutocomplete = (*AuthTuneCommand)(nil)
)

type AuthTuneCommand struct {
	*BaseCommand

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
	flagUserLockoutThreshold            uint
	flagUserLockoutDuration             time.Duration
	flagUserLockoutCounterResetDuration time.Duration
	flagUserLockoutDisable              bool
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
		Usage: "The default lease TTL for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured default lease TTL, " +
			"or a previously configured value for the auth method.",
	})

	f.StringVar(&StringVar{
		Name:   flagNameDescription,
		Target: &c.flagDescription,
		Usage: "Human-friendly description of the this auth method. This overrides " +
			"the current stored value, if any.",
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
		Usage: "The maximum lease TTL for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured maximum lease TTL, " +
			"or a previously configured value for the auth method.",
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
		Usage: "Response header value that plugins will be allowed to set. To specify " +
			"multiple values, specify this flag multiple times.",
	})

	f.StringMapVar(&StringMapVar{
		Name:       "options",
		Target:     &c.flagOptions,
		Completion: complete.PredictAnything,
		Usage: "Key-value pair provided as key=value for the mount options. " +
			"This can be specified multiple times.",
	})

	f.StringVar(&StringVar{
		Name:   flagNameTokenType,
		Target: &c.flagTokenType,
		Usage:  "Sets a forced token type for the mount.",
	})

	f.IntVar(&IntVar{
		Name:    "version",
		Target:  &c.flagVersion,
		Default: 0,
		Usage:   "Select the version of the auth method to run. Not supported by all auth methods.",
	})

	f.UintVar(&UintVar{
		Name:   flagNameUserLockoutThreshold,
		Target: &c.flagUserLockoutThreshold,
		Usage: "The threshold for user lockout for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured user lockout threshold, " +
			"or a previously configured value for the auth method.",
	})

	f.DurationVar(&DurationVar{
		Name:       flagNameUserLockoutDuration,
		Target:     &c.flagUserLockoutDuration,
		Completion: complete.PredictAnything,
		Usage: "The user lockout duration for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured user lockout duration, " +
			"or a previously configured value for the auth method.",
	})

	f.DurationVar(&DurationVar{
		Name:       flagNameUserLockoutCounterResetDuration,
		Target:     &c.flagUserLockoutCounterResetDuration,
		Completion: complete.PredictAnything,
		Usage: "The user lockout counter reset duration for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured user lockout counter reset duration, " +
			"or a previously configured value for the auth method.",
	})

	f.BoolVar(&BoolVar{
		Name:    flagNameUserLockoutDisable,
		Target:  &c.flagUserLockoutDisable,
		Default: false,
		Usage: "Disable user lockout for this auth method. If unspecified, this " +
			"defaults to the Vault server's globally configured user lockout disable, " +
			"or a previously configured value for the auth method.",
	})

	f.StringVar(&StringVar{
		Name:    flagNamePluginVersion,
		Target:  &c.flagPluginVersion,
		Default: "",
		Usage: "Select the semantic version of the plugin to run. The new version must be registered in " +
			"the plugin catalog, and will not start running until the plugin is reloaded.",
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

		if fl.Name == flagNameTokenType {
			mountConfigInput.TokenType = c.flagTokenType
		}
		switch fl.Name {
		case flagNameUserLockoutThreshold, flagNameUserLockoutDuration, flagNameUserLockoutCounterResetDuration, flagNameUserLockoutDisable:
			if mountConfigInput.UserLockoutConfig == nil {
				mountConfigInput.UserLockoutConfig = &api.UserLockoutConfigInput{}
			}
		}
		if fl.Name == flagNameUserLockoutThreshold {
			mountConfigInput.UserLockoutConfig.LockoutThreshold = strconv.FormatUint(uint64(c.flagUserLockoutThreshold), 10)
		}
		if fl.Name == flagNameUserLockoutDuration {
			mountConfigInput.UserLockoutConfig.LockoutDuration = ttlToAPI(c.flagUserLockoutDuration)
		}
		if fl.Name == flagNameUserLockoutCounterResetDuration {
			mountConfigInput.UserLockoutConfig.LockoutCounterResetDuration = ttlToAPI(c.flagUserLockoutCounterResetDuration)
		}
		if fl.Name == flagNameUserLockoutDisable {
			mountConfigInput.UserLockoutConfig.DisableLockout = &c.flagUserLockoutDisable
		}

		if fl.Name == flagNamePluginVersion {
			mountConfigInput.PluginVersion = c.flagPluginVersion
		}
	})

	// Append /auth (since that's where auths live) and a trailing slash to
	// indicate it's a path in output
	mountPath := ensureTrailingSlash(sanitizePath(args[0]))

	if err := client.Sys().TuneMount("/auth/"+mountPath, mountConfigInput); err != nil {
		c.UI.Error(fmt.Sprintf("Error tuning auth method %s: %s", mountPath, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Tuned the auth method at: %s", mountPath))
	return 0
}
