// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/cli"
	hcpvlib "github.com/hashicorp/vault-hcp-lib"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/cliconfig"
	"github.com/hashicorp/vault/api/tokenhelper"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/mattn/go-isatty"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/posener/complete"
)

const (
	// maxLineLength is the maximum width of any line.
	maxLineLength int = 78

	// notSetValue is a flag value for a not-set value
	notSetValue = "(not set)"
)

// reRemoveWhitespace is a regular expression for stripping whitespace from
// a string.
var reRemoveWhitespace = regexp.MustCompile(`[\s]+`)

type BaseCommand struct {
	UI cli.Ui

	flags     *FlagSets
	flagsOnce sync.Once

	// Flag* fields hold parsed flag values. They are exported so the commands
	// in command/client and command/server can read them.
	FlagAddress           string
	FlagAgentProxyAddress string
	FlagCACert            string
	FlagCAPath            string
	FlagClientCert        string
	FlagClientKey         string
	FlagNamespace         string
	FlagNS                string
	FlagPolicyOverride    bool
	FlagTLSServerName     string
	FlagTLSSkipVerify     bool
	FlagDisableRedirects  bool
	FlagWrapTTL           time.Duration
	FlagUnlockKey         string
	FlagTPMStateDir       string
	FlagTPMDevicePath     string

	FlagFormat           string
	FlagField            string
	FlagDetailed         bool
	FlagOutputCurlString bool
	FlagOutputPolicy     bool
	FlagNonInteractive   bool
	addrWarning          string

	FlagMFA []string

	FlagHeader map[string]string

	tokenHelper    tokenhelper.TokenHelper
	hcpTokenHelper hcpvlib.HCPTokenHelper
	stdin          io.Reader

	FlagSnapshotID string

	client *api.Client
}

// Client returns the HTTP API client. The client is cached on the command to
// save performance on future calls.
func (c *BaseCommand) Client() (*api.Client, error) {
	// Read the test client if present
	if c.client != nil {
		// Ignoring homedir errors here and moving on to avoid
		// spamming user with warnings/errors that homedir isn't set.
		path, err := homedir.Dir()
		if err == nil {
			if err := c.applyHCPConfig(path); err != nil {
				return nil, err
			}
		}

		return c.client, nil
	}

	cfg := api.DefaultConfig()

	if err := cfg.ReadEnvironment(); err != nil {
		return nil, errors.Wrap(err, "failed to read environment")
	}

	if c.FlagAddress != "" {
		cfg.Address = c.FlagAddress
	}
	if c.FlagAgentProxyAddress != "" {
		cfg.AgentAddress = c.FlagAgentProxyAddress
	}

	if c.FlagOutputCurlString {
		cfg.OutputCurlString = c.FlagOutputCurlString
	}
	if c.FlagOutputPolicy {
		cfg.OutputPolicy = c.FlagOutputPolicy
	}

	// If we need custom TLS configuration, then set it
	if c.FlagCACert != "" || c.FlagCAPath != "" || c.FlagClientCert != "" ||
		c.FlagClientKey != "" || c.FlagTLSServerName != "" || c.FlagTLSSkipVerify {
		t := &api.TLSConfig{
			CACert:        c.FlagCACert,
			CAPath:        c.FlagCAPath,
			ClientCert:    c.FlagClientCert,
			ClientKey:     c.FlagClientKey,
			TLSServerName: c.FlagTLSServerName,
			Insecure:      c.FlagTLSSkipVerify,
		}

		// Setup TLS config
		if err := cfg.ConfigureTLS(t); err != nil {
			return nil, errors.Wrap(err, "failed to setup TLS config")
		}
	}

	// Build the client
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client")
	}

	// Turn off retries on the CLI
	if os.Getenv(api.EnvVaultMaxRetries) == "" {
		client.SetMaxRetries(0)
	}

	// Set the wrapping function
	client.SetWrappingLookupFunc(c.DefaultWrappingLookupFunc)

	// Get the token if it came in from the environment
	token := client.Token()

	// If we don't have a token, check the token helper
	if token == "" {
		helper, err := c.TokenHelper()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get token helper")
		}
		token, err = helper.Get()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get token from token helper")
		}
	}

	// Set the token
	if token != "" {
		client.SetToken(token)
	}

	client.SetMFACreds(c.FlagMFA)

	// FlagNS takes precedence over FlagNamespace. After resolution, point both
	// flags to the same value to be able to use them interchangeably anywhere.
	if c.FlagNS != notSetValue {
		c.FlagNamespace = c.FlagNS
	}
	if c.FlagNamespace != notSetValue {
		client.SetNamespace(namespace.Canonicalize(c.FlagNamespace))
	}
	if c.FlagPolicyOverride {
		client.SetPolicyOverride(c.FlagPolicyOverride)
	}

	if c.FlagHeader != nil {
		var forbiddenHeaders []string
		for key, val := range c.FlagHeader {
			if strings.HasPrefix(key, "X-Vault-") {
				forbiddenHeaders = append(forbiddenHeaders, key)
				continue
			}
			client.AddHeader(key, val)
		}

		if len(forbiddenHeaders) > 0 {
			return nil, fmt.Errorf("failed to setup Headers[%s]: Header starting by 'X-Vault-' are for internal usage only", strings.Join(forbiddenHeaders, ", "))
		}
	}

	c.client = client

	// Ignoring homedir errors here and moving on to avoid
	// spamming user with warnings/errors that homedir isn't set.
	path, err := homedir.Dir()
	if err == nil {
		if err := c.applyHCPConfig(path); err != nil {
			return nil, err
		}
	}

	if c.addrWarning != "" && c.UI != nil {
		if os.Getenv("VAULT_ADDR") == "" && !c.flags.HadAddressFlag {
			if !c.FlagNonInteractive && isatty.IsTerminal(os.Stdin.Fd()) {
				c.UI.Warn(WrapAtLength(c.addrWarning))
			}
		}
	}

	return client, nil
}

func (c *BaseCommand) applyHCPConfig(path string) error {
	if c.hcpTokenHelper == nil {
		c.hcpTokenHelper = c.HCPTokenHelper()
	}

	hcpToken, err := c.hcpTokenHelper.GetHCPToken(path)
	if err != nil {
		return err
	}

	if hcpToken != nil {
		cookie := &http.Cookie{
			Name:    "hcp_access_token",
			Value:   hcpToken.AccessToken,
			Expires: hcpToken.AccessTokenExpiry,
		}

		if err := c.client.SetHCPCookie(cookie); err != nil {
			return fmt.Errorf("unable to correctly connect to the HCP Vault cluster; please reconnect to HCP: %w", err)
		}

		if err := c.client.SetAddress(hcpToken.ProxyAddr); err != nil {
			return fmt.Errorf("unable to correctly set the HCP address: %w", err)
		}

		// remove address warning since address was set to HCP's address
		c.addrWarning = ""
	}

	return nil
}

// Format returns the output format flag value (e.g. "json", "table").
func (c *BaseCommand) Format() string {
	return c.FlagFormat
}

// SetAddress sets the token helper on the command; useful for the demo server and other outside cases.
func (c *BaseCommand) SetAddress(addr string) {
	c.FlagAddress = addr
}

// ClientDirect returns the cached API client without re-initializing it.
// This is safe to call only after Client() has been called at least once.
func (c *BaseCommand) ClientDirect() *api.Client {
	return c.client
}

// SetTokenHelper sets the token helper on the command.
func (c *BaseCommand) SetTokenHelper(th tokenhelper.TokenHelper) {
	c.tokenHelper = th
}

// SetHCPTokenHelper sets the HCP token helper on the command.
func (c *BaseCommand) SetHCPTokenHelper(h hcpvlib.HCPTokenHelper) {
	c.hcpTokenHelper = h
}

// SetClient sets a pre-built API client on the command.
func (c *BaseCommand) SetClient(cl *api.Client) {
	c.client = cl
}

// Stdin returns the input stream for the command.
func (c *BaseCommand) Stdin() io.Reader {
	if c.stdin != nil {
		return c.stdin
	}
	return os.Stdin
}

// NewBaseCommand constructs a BaseCommand wired from a RunOptions value.
func NewBaseCommand(ui cli.Ui, runOpts *RunOptions) *BaseCommand {
	bc := &BaseCommand{UI: ui}
	if runOpts != nil {
		bc.FlagAddress = runOpts.Address
		bc.tokenHelper = runOpts.TokenHelper
		bc.client = runOpts.Client
		bc.hcpTokenHelper = runOpts.HCPTokenHelper
		bc.stdin = runOpts.Stdin
	}
	return bc
}

// TokenHelper returns the token helper attached to the command.
func (c *BaseCommand) TokenHelper() (tokenhelper.TokenHelper, error) {
	if c.tokenHelper != nil {
		return c.tokenHelper, nil
	}

	helper, err := cliconfig.DefaultTokenHelper()
	if err != nil {
		return nil, err
	}
	return helper, nil
}

// HCPTokenHelper returns the HCPToken helper attached to the command.
func (c *BaseCommand) HCPTokenHelper() hcpvlib.HCPTokenHelper {
	if c.hcpTokenHelper != nil {
		return c.hcpTokenHelper
	}
	return &hcpvlib.InternalHCPTokenHelper{}
}

// DefaultWrappingLookupFunc is the default wrapping function based on the
// CLI flag.
func (c *BaseCommand) DefaultWrappingLookupFunc(operation, path string) string {
	if c.FlagWrapTTL != 0 {
		return c.FlagWrapTTL.String()
	}

	return api.DefaultWrappingLookupFunc(operation, path)
}

// GetMFAValidationRequired checks to see if the secret exists and has an MFA
// requirement. If MFA is required and the number of constraints is greater than
// 1, we can assert that interactive validation is not required.
func (c *BaseCommand) GetMFAValidationRequired(secret *api.Secret) bool {
	if secret != nil && secret.Auth != nil && secret.Auth.MFARequirement != nil {
		if c.FlagMFA == nil && len(secret.Auth.MFARequirement.MFAConstraints) == 1 {
			return true
		} else if len(secret.Auth.MFARequirement.MFAConstraints) > 1 {
			return true
		}
	}

	return false
}

// GetInteractiveMFAMethodInfo returns MFA method information only if operating
// in interactive mode and one MFA method is configured.
func (c *BaseCommand) GetInteractiveMFAMethodInfo(secret *api.Secret) *MFAMethodInfo {
	if secret == nil || secret.Auth == nil || secret.Auth.MFARequirement == nil {
		return nil
	}

	mfaConstraints := secret.Auth.MFARequirement.MFAConstraints
	if c.FlagNonInteractive || len(mfaConstraints) != 1 || !isatty.IsTerminal(os.Stdin.Fd()) {
		return nil
	}

	for _, mfaConstraint := range mfaConstraints {
		if len(mfaConstraint.Any) != 1 {
			return nil
		}

		return &MFAMethodInfo{
			MethodType:  mfaConstraint.Any[0].Type,
			MethodID:    mfaConstraint.Any[0].ID,
			UsePasscode: mfaConstraint.Any[0].UsesPasscode,
		}
	}

	return nil
}

// ValidateMFA performs interactive MFA validation.
func (c *BaseCommand) ValidateMFA(reqID string, methodInfo MFAMethodInfo) (*api.Secret, error) {
	var passcode string
	var err error
	if methodInfo.UsePasscode {
		passcode, err = c.UI.AskSecret(fmt.Sprintf("Enter the passphrase for methodID %q of type %q:", methodInfo.MethodID, methodInfo.MethodType))
		if err != nil {
			return nil, fmt.Errorf("failed to read passphrase: %w. please validate the login by sending a request to sys/mfa/validate", err)
		}
	} else {
		c.UI.Warn("Asking Vault to perform MFA validation with upstream service. " +
			"You should receive a push notification in your authenticator app shortly")
	}

	// passcode could be an empty string
	mfaPayload := map[string]interface{}{
		methodInfo.MethodID: []string{passcode},
	}

	client, err := c.Client()
	if err != nil {
		return nil, err
	}

	return client.Sys().MFAValidate(reqID, mfaPayload)
}

type FlagSetBit uint

const (
	FlagSetNone FlagSetBit = 1 << iota
	FlagSetHTTP
	FlagSetOutputField
	FlagSetOutputFormat
	FlagSetOutputDetailed
	FlagSetSnapshot
)

// FlagSetForBit creates the flags for this command, cached for performance.
func (c *BaseCommand) FlagSetForBit(bit FlagSetBit) *FlagSets {
	return c.flagSet(bit)
}

// FlagSetsResult returns the parsed FlagSets from the last Parse call.
func (c *BaseCommand) FlagSetsResult() *FlagSets {
	return c.flags
}

// flagSet creates the flags for this command. The result is cached on the
// command to save performance on future calls.
func (c *BaseCommand) flagSet(bit FlagSetBit) *FlagSets {
	c.flagsOnce.Do(func() {
		set := NewFlagSets(c.UI)

		// These flag sets will apply to all leaf subcommands.
		// TODO: Optional, but FlagSetHTTP can be safely removed from the individual
		// Flags() subcommands.
		bit = bit | FlagSetHTTP

		if bit&FlagSetHTTP != 0 {
			f := set.NewFlagSet("HTTP Options")

			addrStringVar := &StringVar{
				Name:        FlagNameAddress,
				Target:      &c.FlagAddress,
				EnvVar:      api.EnvVaultAddress,
				Completion:  complete.PredictAnything,
				Normalizers: []func(string) string{configutil.NormalizeAddr},
				Usage:       "Address of the Vault server.",
			}

			if c.FlagAddress != "" {
				addrStringVar.Default = c.FlagAddress
			} else {
				addrStringVar.Default = "https://127.0.0.1:8200"
				c.addrWarning = fmt.Sprintf("WARNING! VAULT_ADDR and -address unset. Defaulting to %s.", addrStringVar.Default)
			}
			f.StringVar(addrStringVar)

			agentAddrStringVar := &StringVar{
				Name:        "agent-address",
				Target:      &c.FlagAgentProxyAddress,
				EnvVar:      api.EnvVaultAgentAddr,
				Completion:  complete.PredictAnything,
				Normalizers: []func(string) string{configutil.NormalizeAddr},
				Usage:       "Address of the Agent.",
			}
			f.StringVar(agentAddrStringVar)

			f.StringVar(&StringVar{
				Name:       FlagNameCACert,
				Target:     &c.FlagCACert,
				Default:    "",
				EnvVar:     api.EnvVaultCACert,
				Completion: complete.PredictFiles("*"),
				Usage: "Path on the local disk to a single PEM-encoded CA " +
					"certificate to verify the Vault server's SSL certificate. This " +
					"takes precedence over -ca-path.",
			})

			f.StringVar(&StringVar{
				Name:       FlagNameCAPath,
				Target:     &c.FlagCAPath,
				Default:    "",
				EnvVar:     api.EnvVaultCAPath,
				Completion: complete.PredictDirs("*"),
				Usage: "Path on the local disk to a directory of PEM-encoded CA " +
					"certificates to verify the Vault server's SSL certificate.",
			})

			f.StringVar(&StringVar{
				Name:       FlagNameClientCert,
				Target:     &c.FlagClientCert,
				Default:    "",
				EnvVar:     api.EnvVaultClientCert,
				Completion: complete.PredictFiles("*"),
				Usage: "Path on the local disk to a single PEM-encoded CA " +
					"certificate to use for TLS authentication to the Vault server. If " +
					"this flag is specified, -client-key is also required.",
			})

			f.StringVar(&StringVar{
				Name:       FlagNameClientKey,
				Target:     &c.FlagClientKey,
				Default:    "",
				EnvVar:     api.EnvVaultClientKey,
				Completion: complete.PredictFiles("*"),
				Usage: "Path on the local disk to a single PEM-encoded private key " +
					"matching the client certificate from -client-cert.",
			})

			f.StringVar(&StringVar{
				Name:       "namespace",
				Target:     &c.FlagNamespace,
				Default:    notSetValue, // this can never be a real value
				EnvVar:     api.EnvVaultNamespace,
				Completion: complete.PredictAnything,
				Usage: "The namespace to use for the command. Setting this is not " +
					"necessary but allows using relative paths. -ns can be used as " +
					"shortcut.",
			})

			f.StringVar(&StringVar{
				Name:       "ns",
				Target:     &c.FlagNS,
				Default:    notSetValue,
				Completion: complete.PredictAnything,
				Hidden:     true,
				Usage:      "Alias for -namespace. This takes precedence over -namespace.",
			})

			f.StringVar(&StringVar{
				Name:       FlagNameTLSServerName,
				Target:     &c.FlagTLSServerName,
				Default:    "",
				EnvVar:     api.EnvVaultTLSServerName,
				Completion: complete.PredictAnything,
				Usage: "Name to use as the SNI host when connecting to the Vault " +
					"server via TLS.",
			})

			f.BoolVar(&BoolVar{
				Name:    FlagNameTLSSkipVerify,
				Target:  &c.FlagTLSSkipVerify,
				Default: false,
				EnvVar:  api.EnvVaultSkipVerify,
				Usage: "Disable verification of TLS certificates. Using this option " +
					"is highly discouraged as it decreases the security of data " +
					"transmissions to and from the Vault server.",
			})

			f.BoolVar(&BoolVar{
				Name:    FlagNameDisableRedirects,
				Target:  &c.FlagDisableRedirects,
				Default: false,
				EnvVar:  api.EnvVaultDisableRedirects,
				Usage: "Disable the default client behavior, which honors a single " +
					"redirect response from a request",
			})

			f.BoolVar(&BoolVar{
				Name:    "policy-override",
				Target:  &c.FlagPolicyOverride,
				Default: false,
				Usage: "Override a Sentinel policy that has a soft-mandatory " +
					"enforcement_level specified",
			})

			f.DurationVar(&DurationVar{
				Name:       "wrap-ttl",
				Target:     &c.FlagWrapTTL,
				Default:    0,
				EnvVar:     api.EnvVaultWrapTTL,
				Completion: complete.PredictAnything,
				Usage: "Wraps the response in a cubbyhole token with the requested " +
					"TTL. The response is available via the \"vault unwrap\" command. " +
					"The TTL is specified as a numeric string with suffix like \"30s\" " +
					"or \"5m\".",
			})

			f.StringSliceVar(&StringSliceVar{
				Name:       "mfa",
				Target:     &c.FlagMFA,
				Default:    nil,
				EnvVar:     api.EnvVaultMFA,
				Completion: complete.PredictAnything,
				Usage:      "Supply MFA credentials as part of X-Vault-MFA header.",
			})

			f.BoolVar(&BoolVar{
				Name:    "output-curl-string",
				Target:  &c.FlagOutputCurlString,
				Default: false,
				Usage: "Instead of executing the request, print an equivalent cURL " +
					"command string and exit.",
			})

			f.BoolVar(&BoolVar{
				Name:    "output-policy",
				Target:  &c.FlagOutputPolicy,
				Default: false,
				Usage: "Instead of executing the request, print an example HCL " +
					"policy that would be required to run this command, and exit.",
			})

			f.StringVar(&StringVar{
				Name:       "unlock-key",
				Target:     &c.FlagUnlockKey,
				Default:    notSetValue,
				Completion: complete.PredictNothing,
				Usage:      "Key to unlock a namespace API lock.",
			})

			f.StringMapVar(&StringMapVar{
				Name:       "header",
				Target:     &c.FlagHeader,
				Completion: complete.PredictAnything,
				Usage: "Key-value pair provided as key=value to provide http header added to any request done by the CLI." +
					"Trying to add headers starting with 'X-Vault-' is forbidden and will make the command fail " +
					"This can be specified multiple times.",
			})

			f.BoolVar(&BoolVar{
				Name:    "non-interactive",
				Target:  &c.FlagNonInteractive,
				Default: false,
				Usage:   "When set true, prevents asking the user for input via the terminal.",
			})

			f.StringVar(&StringVar{
				Name:    "tpm-state-dir",
				Target:  &c.FlagTPMStateDir,
				Default: ".",
				Usage:   "Sets where to write the cert and associated files. Default is working directory",
			})

			f.StringVar(&StringVar{
				Name:    "tpm-device-path",
				Target:  &c.FlagTPMDevicePath,
				Default: "",
				Usage: "Path to the TPM device (e.g. /dev/tpm0 or a unix socket). " +
					"Autodetected by default.",
			})
		}

		if bit&(FlagSetOutputField|FlagSetOutputFormat|FlagSetOutputDetailed) != 0 {
			outputSet := set.NewFlagSet("Output Options")

			if bit&FlagSetOutputField != 0 {
				outputSet.StringVar(&StringVar{
					Name:       "field",
					Target:     &c.FlagField,
					Default:    "",
					Completion: complete.PredictAnything,
					Usage: "Print only the field with the given name. Specifying " +
						"this option will take precedence over other formatting " +
						"directives. The result will not have a trailing newline " +
						"making it ideal for piping to other processes.",
				})
			}

			if bit&FlagSetOutputFormat != 0 {
				outputSet.StringVar(&StringVar{
					Name:       "format",
					Target:     &c.FlagFormat,
					Default:    "table",
					EnvVar:     EnvVaultFormat,
					Completion: complete.PredictSet("table", "json", "yaml", "pretty", "raw"),
					Usage: `Print the output in the given format. Valid formats
						are "table", "json", "yaml", or "pretty". "raw" is allowed
						for 'vault read' operations only.`,
				})
			}

			if bit&FlagSetOutputDetailed != 0 {
				outputSet.BoolVar(&BoolVar{
					Name:    "detailed",
					Target:  &c.FlagDetailed,
					Default: false,
					EnvVar:  EnvVaultDetailed,
					Usage:   "Enables additional metadata during some operations",
				})
			}

			if bit&FlagSetSnapshot != 0 {
				outputSet.StringVar(&StringVar{
					Name:       "snapshot-id",
					Target:     &c.FlagSnapshotID,
					Default:    "",
					Completion: complete.PredictAnything,
					Usage:      "ID of the loaded snapshot that this command will use",
				})
			}
		}

		c.flags = set
	})

	return c.flags
}

// FlagSets is a group of flag sets.
type FlagSets struct {
	flagSets    []*FlagSet
	mainSet     *flag.FlagSet
	hiddens     map[string]struct{}
	completions complete.Flags
	ui          cli.Ui
	// HadAddressFlag signals if the FlagSet had an -address
	// flag set, for the purposes of warning (see also:
	// BaseCommand::addrWarning).
	HadAddressFlag bool
}

// NewFlagSets creates a new flag sets.
func NewFlagSets(ui cli.Ui) *FlagSets {
	mainSet := flag.NewFlagSet("", flag.ContinueOnError)

	// Errors and usage are controlled by the CLI.
	mainSet.Usage = func() {}
	mainSet.SetOutput(io.Discard)

	return &FlagSets{
		flagSets:    make([]*FlagSet, 0, 6),
		mainSet:     mainSet,
		hiddens:     make(map[string]struct{}),
		completions: complete.Flags{},
		ui:          ui,
	}
}

// NewFlagSet creates a new flag set from the given flag sets.
func (f *FlagSets) NewFlagSet(name string) *FlagSet {
	flagSet := NewFlagSet(name)
	flagSet.mainSet = f.mainSet
	flagSet.completions = f.completions
	f.flagSets = append(f.flagSets, flagSet)
	return flagSet
}

// Completions returns the completions for this flag set.
func (f *FlagSets) Completions() complete.Flags {
	return f.completions
}

// Lookup returns the named flag from the underlying flag set, or nil if not
// found. This satisfies the FlagLookup interface used by LoginHandlerWithFlags.
func (f *FlagSets) Lookup(name string) *flag.Flag {
	return f.mainSet.Lookup(name)
}

type (
	ParseOptions              interface{}
	ParseOptionAllowRawFormat bool
	DisableDisplayFlagWarning bool
)

// Parse parses the given flags, returning any errors.
// Warnings, if any, regarding the arguments format are sent to stdout
func (f *FlagSets) Parse(args []string, opts ...ParseOptions) error {
	// Before parsing, check to see if we have an address flag, for the
	// purposes of warning later. This must be done now, as the argument
	// will be removed during parsing.
	for _, arg := range args {
		if strings.HasPrefix(arg, "-address") {
			f.HadAddressFlag = true
		}
	}

	err := f.mainSet.Parse(args)

	displayFlagWarningsDisabled := false
	for _, opt := range opts {
		if value, ok := opt.(DisableDisplayFlagWarning); ok {
			displayFlagWarningsDisabled = bool(value)
		}
	}
	if !displayFlagWarningsDisabled {
		warnings := GenerateFlagWarnings(f.Args())
		if warnings != "" && Format(f.ui) == "table" {
			f.ui.Warn(warnings)
		}
	}

	if err != nil {
		return err
	}

	// Now surface any other errors.
	return GenerateFlagErrors(f, opts...)
}

// Parsed reports whether the command-line flags have been parsed.
func (f *FlagSets) Parsed() bool {
	return f.mainSet.Parsed()
}

// Args returns the remaining args after parsing.
func (f *FlagSets) Args() []string {
	return f.mainSet.Args()
}

// Visit visits the flags in lexicographical order, calling fn for each. It
// visits only those flags that have been set.
func (f *FlagSets) Visit(fn func(*flag.Flag)) {
	f.mainSet.Visit(fn)
}

// Help builds custom help for this command, grouping by flag set.
func (f *FlagSets) Help() string {
	var out bytes.Buffer

	for _, set := range f.flagSets {
		printFlagTitle(&out, set.name+":")
		set.VisitAll(func(f *flag.Flag) {
			// Skip any hidden flags
			if v, ok := f.Value.(FlagVisibility); ok && v.Hidden() {
				return
			}
			printFlagDetail(&out, f)
		})
	}

	return strings.TrimRight(out.String(), "\n")
}

// FlagSet is a grouped wrapper around a real flag set and a grouped flag set.
type FlagSet struct {
	name        string
	flagSet     *flag.FlagSet
	mainSet     *flag.FlagSet
	completions complete.Flags
}

// NewFlagSet creates a new flag set.
func NewFlagSet(name string) *FlagSet {
	return &FlagSet{
		name:    name,
		flagSet: flag.NewFlagSet(name, flag.ContinueOnError),
	}
}

// Name returns the name of this flag set.
func (f *FlagSet) Name() string {
	return f.name
}

func (f *FlagSet) Visit(fn func(*flag.Flag)) {
	f.flagSet.Visit(fn)
}

func (f *FlagSet) VisitAll(fn func(*flag.Flag)) {
	f.flagSet.VisitAll(fn)
}

// printFlagTitle prints a consistently-formatted title to the given writer.
func printFlagTitle(w io.Writer, s string) {
	fmt.Fprintf(w, "%s\n\n", s)
}

// printFlagDetail prints a single flag to the given writer.
func printFlagDetail(w io.Writer, f *flag.Flag) {
	// Check if the flag is hidden - do not print any flag detail or help output
	// if it is hidden.
	if h, ok := f.Value.(FlagVisibility); ok && h.Hidden() {
		return
	}

	// Check for a detailed example
	example := ""
	if t, ok := f.Value.(FlagExample); ok {
		example = t.Example()
	}

	if example != "" {
		fmt.Fprintf(w, "  -%s=<%s>\n", f.Name, example)
	} else {
		fmt.Fprintf(w, "  -%s\n", f.Name)
	}

	usage := reRemoveWhitespace.ReplaceAllString(f.Usage, " ")
	indented := wrapAtLengthWithPadding(usage, 6)
	fmt.Fprintf(w, "%s\n\n", indented)
}
