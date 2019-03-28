package command

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*AuthEnableCommand)(nil)
var _ cli.CommandAutocomplete = (*AuthEnableCommand)(nil)

type AuthEnableCommand struct {
	*BaseCommand

	flagDescription               string
	flagPath                      string
	flagDefaultLeaseTTL           time.Duration
	flagMaxLeaseTTL               time.Duration
	flagAuditNonHMACRequestKeys   []string
	flagAuditNonHMACResponseKeys  []string
	flagListingVisibility         string
	flagPluginName                string
	flagPassthroughRequestHeaders []string
	flagAllowedResponseHeaders    []string
	flagOptions                   map[string]string
	flagLocal                     bool
	flagSealWrap                  bool
	flagTokenType                 string
	flagVersion                   int
}

func (c *AuthEnableCommand) Synopsis() string {
	return "Enables a new auth method"
}

func (c *AuthEnableCommand) Help() string {
	helpText := `
Usage: vault auth enable [options] TYPE

  Enables a new auth method. An auth method is responsible for authenticating
  users or machines and assigning them policies with which they can access
  Vault.

  Enable the userpass auth method at userpass/:

      $ vault auth enable userpass

  Enable the LDAP auth method at auth-prod/:

      $ vault auth enable -path=auth-prod ldap

  Enable a custom auth plugin (after it's registered in the plugin registry):

      $ vault auth enable -path=my-auth -plugin-name=my-auth-plugin plugin

      OR (preferred way):

      $ vault auth enable -path=my-auth my-auth-plugin

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuthEnableCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "description",
		Target:     &c.flagDescription,
		Completion: complete.PredictAnything,
		Usage: "Human-friendly description for the purpose of this " +
			"auth method.",
	})

	f.StringVar(&StringVar{
		Name:       "path",
		Target:     &c.flagPath,
		Default:    "", // The default is complex, so we have to manually document
		Completion: complete.PredictAnything,
		Usage: "Place where the auth method will be accessible. This must be " +
			"unique across all auth methods. This defaults to the \"type\" of " +
			"the auth method. The auth method will be accessible at " +
			"\"/auth/<path>\".",
	})

	f.DurationVar(&DurationVar{
		Name:       "default-lease-ttl",
		Target:     &c.flagDefaultLeaseTTL,
		Completion: complete.PredictAnything,
		Usage: "The default lease TTL for this auth method. If unspecified, " +
			"this defaults to the Vault server's globally configured default lease " +
			"TTL.",
	})

	f.DurationVar(&DurationVar{
		Name:       "max-lease-ttl",
		Target:     &c.flagMaxLeaseTTL,
		Completion: complete.PredictAnything,
		Usage: "The maximum lease TTL for this auth method. If unspecified, " +
			"this defaults to the Vault server's globally configured maximum lease " +
			"TTL.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNameAuditNonHMACRequestKeys,
		Target: &c.flagAuditNonHMACRequestKeys,
		Usage: "Comma-separated string or list of keys that will not be HMAC'd by audit " +
			"devices in the request data object.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNameAuditNonHMACResponseKeys,
		Target: &c.flagAuditNonHMACResponseKeys,
		Usage: "Comma-separated string or list of keys that will not be HMAC'd by audit " +
			"devices in the response data object.",
	})

	f.StringVar(&StringVar{
		Name:   flagNameListingVisibility,
		Target: &c.flagListingVisibility,
		Usage:  "Determines the visibility of the mount in the UI-specific listing endpoint.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNamePassthroughRequestHeaders,
		Target: &c.flagPassthroughRequestHeaders,
		Usage: "Comma-separated string or list of request header values that " +
			"will be sent to the plugin",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   flagNameAllowedResponseHeaders,
		Target: &c.flagAllowedResponseHeaders,
		Usage: "Comma-separated string or list of response header values that " +
			"plugins will be allowed to set",
	})

	f.StringVar(&StringVar{
		Name:       "plugin-name",
		Target:     &c.flagPluginName,
		Completion: c.PredictVaultPlugins(consts.PluginTypeCredential),
		Usage: "Name of the auth method plugin. This plugin name must already " +
			"exist in the Vault server's plugin catalog.",
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
		Usage: "Mark the auth method as local-only. Local auth methods are " +
			"not replicated nor removed by replication.",
	})

	f.BoolVar(&BoolVar{
		Name:    "seal-wrap",
		Target:  &c.flagSealWrap,
		Default: false,
		Usage:   "Enable seal wrapping of critical values in the secrets engine.",
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

	return set
}

func (c *AuthEnableCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultAvailableAuths()
}

func (c *AuthEnableCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthEnableCommand) Run(args []string) int {
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

	authType := strings.TrimSpace(args[0])
	if authType == "plugin" {
		authType = c.flagPluginName
	}

	// If no path is specified, we default the path to the backend type
	// or use the plugin name if it's a plugin backend
	authPath := c.flagPath
	if authPath == "" {
		if authType == "plugin" {
			authPath = c.flagPluginName
		} else {
			authPath = authType
		}
	}

	// Append a trailing slash to indicate it's a path in output
	authPath = ensureTrailingSlash(authPath)

	if c.flagVersion > 0 {
		if c.flagOptions == nil {
			c.flagOptions = make(map[string]string)
		}
		c.flagOptions["version"] = strconv.Itoa(c.flagVersion)
	}

	authOpts := &api.EnableAuthOptions{
		Type:        authType,
		Description: c.flagDescription,
		Local:       c.flagLocal,
		SealWrap:    c.flagSealWrap,
		Config: api.AuthConfigInput{
			DefaultLeaseTTL: c.flagDefaultLeaseTTL.String(),
			MaxLeaseTTL:     c.flagMaxLeaseTTL.String(),
		},
		Options: c.flagOptions,
	}

	// Set these values only if they are provided in the CLI
	f.Visit(func(fl *flag.Flag) {
		if fl.Name == flagNameAuditNonHMACRequestKeys {
			authOpts.Config.AuditNonHMACRequestKeys = c.flagAuditNonHMACRequestKeys
		}

		if fl.Name == flagNameAuditNonHMACResponseKeys {
			authOpts.Config.AuditNonHMACResponseKeys = c.flagAuditNonHMACResponseKeys
		}

		if fl.Name == flagNameListingVisibility {
			authOpts.Config.ListingVisibility = c.flagListingVisibility
		}

		if fl.Name == flagNamePassthroughRequestHeaders {
			authOpts.Config.PassthroughRequestHeaders = c.flagPassthroughRequestHeaders
		}

		if fl.Name == flagNameAllowedResponseHeaders {
			authOpts.Config.AllowedResponseHeaders = c.flagAllowedResponseHeaders
		}

		if fl.Name == flagNameTokenType {
			authOpts.Config.TokenType = c.flagTokenType
		}
	})

	if err := client.Sys().EnableAuthWithOptions(authPath, authOpts); err != nil {
		c.UI.Error(fmt.Sprintf("Error enabling %s auth: %s", authType, err))
		return 2
	}

	authThing := authType + " auth method"
	if authType == "plugin" {
		authThing = c.flagPluginName + " plugin"
	}
	c.UI.Output(fmt.Sprintf("Success! Enabled %s at: %s", authThing, authPath))
	return 0
}
