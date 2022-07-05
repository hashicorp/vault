package command

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mattn/go-isatty"
	"github.com/mitchellh/cli"
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

	flagAddress        string
	flagAgentAddress   string
	flagCACert         string
	flagCAPath         string
	flagClientCert     string
	flagClientKey      string
	flagNamespace      string
	flagNS             string
	flagPolicyOverride bool
	flagTLSServerName  string
	flagTLSSkipVerify  bool
	flagWrapTTL        time.Duration
	flagUnlockKey      string

	flagFormat           string
	flagField            string
	flagDetailed         bool
	flagOutputCurlString bool
	flagOutputPolicy     bool
	flagNonInteractive   bool

	flagMFA []string

	flagHeader map[string]string

	tokenHelper token.TokenHelper

	client *api.Client
}

// Client returns the HTTP API client. The client is cached on the command to
// save performance on future calls.
func (c *BaseCommand) Client() (*api.Client, error) {
	// Read the test client if present
	if c.client != nil {
		return c.client, nil
	}

	config := api.DefaultConfig()

	if err := config.ReadEnvironment(); err != nil {
		return nil, errors.Wrap(err, "failed to read environment")
	}

	if c.flagAddress != "" {
		config.Address = c.flagAddress
	}
	if c.flagAgentAddress != "" {
		config.Address = c.flagAgentAddress
	}

	if c.flagOutputCurlString {
		config.OutputCurlString = c.flagOutputCurlString
	}
	if c.flagOutputPolicy {
		config.OutputPolicy = c.flagOutputPolicy
	}

	// If we need custom TLS configuration, then set it
	if c.flagCACert != "" || c.flagCAPath != "" || c.flagClientCert != "" ||
		c.flagClientKey != "" || c.flagTLSServerName != "" || c.flagTLSSkipVerify {
		t := &api.TLSConfig{
			CACert:        c.flagCACert,
			CAPath:        c.flagCAPath,
			ClientCert:    c.flagClientCert,
			ClientKey:     c.flagClientKey,
			TLSServerName: c.flagTLSServerName,
			Insecure:      c.flagTLSSkipVerify,
		}

		// Setup TLS config
		if err := config.ConfigureTLS(t); err != nil {
			return nil, errors.Wrap(err, "failed to setup TLS config")
		}
	}

	// Build the client
	client, err := api.NewClient(config)
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

	client.SetMFACreds(c.flagMFA)

	// flagNS takes precedence over flagNamespace. After resolution, point both
	// flags to the same value to be able to use them interchangeably anywhere.
	if c.flagNS != notSetValue {
		c.flagNamespace = c.flagNS
	}
	if c.flagNamespace != notSetValue {
		client.SetNamespace(namespace.Canonicalize(c.flagNamespace))
	}
	if c.flagPolicyOverride {
		client.SetPolicyOverride(c.flagPolicyOverride)
	}

	if c.flagHeader != nil {

		var forbiddenHeaders []string
		for key, val := range c.flagHeader {

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

	return client, nil
}

// SetAddress sets the token helper on the command; useful for the demo server and other outside cases.
func (c *BaseCommand) SetAddress(addr string) {
	c.flagAddress = addr
}

// SetTokenHelper sets the token helper on the command.
func (c *BaseCommand) SetTokenHelper(th token.TokenHelper) {
	c.tokenHelper = th
}

// TokenHelper returns the token helper attached to the command.
func (c *BaseCommand) TokenHelper() (token.TokenHelper, error) {
	if c.tokenHelper != nil {
		return c.tokenHelper, nil
	}

	helper, err := DefaultTokenHelper()
	if err != nil {
		return nil, err
	}
	return helper, nil
}

// DefaultWrappingLookupFunc is the default wrapping function based on the
// CLI flag.
func (c *BaseCommand) DefaultWrappingLookupFunc(operation, path string) string {
	if c.flagWrapTTL != 0 {
		return c.flagWrapTTL.String()
	}

	return api.DefaultWrappingLookupFunc(operation, path)
}

func (c *BaseCommand) isInteractiveEnabled(mfaConstraintLen int) bool {
	if mfaConstraintLen != 1 || !isatty.IsTerminal(os.Stdin.Fd()) {
		return false
	}

	if !c.flagNonInteractive {
		return true
	}

	return false
}

// getMFAMethodInfo returns MFA method information only if one MFA method is
// configured.
func (c *BaseCommand) getMFAMethodInfo(mfaConstraintAny map[string]*logical.MFAConstraintAny) MFAMethodInfo {
	for _, mfaConstraint := range mfaConstraintAny {
		if len(mfaConstraint.Any) != 1 {
			return MFAMethodInfo{}
		}

		return MFAMethodInfo{
			methodType:  mfaConstraint.Any[0].Type,
			methodID:    mfaConstraint.Any[0].ID,
			usePasscode: mfaConstraint.Any[0].UsesPasscode,
		}
	}

	return MFAMethodInfo{}
}

func (c *BaseCommand) validateMFA(reqID string, methodInfo MFAMethodInfo) int {
	var passcode string
	var err error
	if methodInfo.usePasscode {
		passcode, err = c.UI.AskSecret(fmt.Sprintf("Enter the passphrase for methodID %q of type %q:", methodInfo.methodID, methodInfo.methodType))
		if err != nil {
			c.UI.Error(fmt.Sprintf("failed to read the passphrase with error %q. please validate the login by sending a request to sys/mfa/validate", err.Error()))
			return 2
		}
	} else {
		c.UI.Warn("Asking Vault to perform MFA validation with upstream service. " +
			"You should receive a push notification in your authenticator app shortly")
	}

	// passcode could be an empty string
	mfaPayload := map[string]interface{}{
		methodInfo.methodID: []string{passcode},
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	secret, err := client.Sys().MFAValidate(reqID, mfaPayload)
	if err != nil {
		c.UI.Error(err.Error())
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}
	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info("Success! Data written to: sys/mfa/validate")
		}
		return 0
	}

	// Handle single field output
	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, secret)
}

type FlagSetBit uint

const (
	FlagSetNone FlagSetBit = 1 << iota
	FlagSetHTTP
	FlagSetOutputField
	FlagSetOutputFormat
	FlagSetOutputDetailed
)

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
				Name:       flagNameAddress,
				Target:     &c.flagAddress,
				EnvVar:     api.EnvVaultAddress,
				Completion: complete.PredictAnything,
				Usage:      "Address of the Vault server.",
			}
			if c.flagAddress != "" {
				addrStringVar.Default = c.flagAddress
			} else {
				addrStringVar.Default = "https://127.0.0.1:8200"
			}
			f.StringVar(addrStringVar)

			agentAddrStringVar := &StringVar{
				Name:       "agent-address",
				Target:     &c.flagAgentAddress,
				EnvVar:     api.EnvVaultAgentAddr,
				Completion: complete.PredictAnything,
				Usage:      "Address of the Agent.",
			}
			f.StringVar(agentAddrStringVar)

			f.StringVar(&StringVar{
				Name:       flagNameCACert,
				Target:     &c.flagCACert,
				Default:    "",
				EnvVar:     api.EnvVaultCACert,
				Completion: complete.PredictFiles("*"),
				Usage: "Path on the local disk to a single PEM-encoded CA " +
					"certificate to verify the Vault server's SSL certificate. This " +
					"takes precedence over -ca-path.",
			})

			f.StringVar(&StringVar{
				Name:       flagNameCAPath,
				Target:     &c.flagCAPath,
				Default:    "",
				EnvVar:     api.EnvVaultCAPath,
				Completion: complete.PredictDirs("*"),
				Usage: "Path on the local disk to a directory of PEM-encoded CA " +
					"certificates to verify the Vault server's SSL certificate.",
			})

			f.StringVar(&StringVar{
				Name:       flagNameClientCert,
				Target:     &c.flagClientCert,
				Default:    "",
				EnvVar:     api.EnvVaultClientCert,
				Completion: complete.PredictFiles("*"),
				Usage: "Path on the local disk to a single PEM-encoded CA " +
					"certificate to use for TLS authentication to the Vault server. If " +
					"this flag is specified, -client-key is also required.",
			})

			f.StringVar(&StringVar{
				Name:       flagNameClientKey,
				Target:     &c.flagClientKey,
				Default:    "",
				EnvVar:     api.EnvVaultClientKey,
				Completion: complete.PredictFiles("*"),
				Usage: "Path on the local disk to a single PEM-encoded private key " +
					"matching the client certificate from -client-cert.",
			})

			f.StringVar(&StringVar{
				Name:       "namespace",
				Target:     &c.flagNamespace,
				Default:    notSetValue, // this can never be a real value
				EnvVar:     api.EnvVaultNamespace,
				Completion: complete.PredictAnything,
				Usage: "The namespace to use for the command. Setting this is not " +
					"necessary but allows using relative paths. -ns can be used as " +
					"shortcut.",
			})

			f.StringVar(&StringVar{
				Name:       "ns",
				Target:     &c.flagNS,
				Default:    notSetValue, // this can never be a real value
				Completion: complete.PredictAnything,
				Hidden:     true,
				Usage:      "Alias for -namespace. This takes precedence over -namespace.",
			})

			f.StringVar(&StringVar{
				Name:       flagTLSServerName,
				Target:     &c.flagTLSServerName,
				Default:    "",
				EnvVar:     api.EnvVaultTLSServerName,
				Completion: complete.PredictAnything,
				Usage: "Name to use as the SNI host when connecting to the Vault " +
					"server via TLS.",
			})

			f.BoolVar(&BoolVar{
				Name:    flagNameTLSSkipVerify,
				Target:  &c.flagTLSSkipVerify,
				Default: false,
				EnvVar:  api.EnvVaultSkipVerify,
				Usage: "Disable verification of TLS certificates. Using this option " +
					"is highly discouraged as it decreases the security of data " +
					"transmissions to and from the Vault server.",
			})

			f.BoolVar(&BoolVar{
				Name:    "policy-override",
				Target:  &c.flagPolicyOverride,
				Default: false,
				Usage: "Override a Sentinel policy that has a soft-mandatory " +
					"enforcement_level specified",
			})

			f.DurationVar(&DurationVar{
				Name:       "wrap-ttl",
				Target:     &c.flagWrapTTL,
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
				Target:     &c.flagMFA,
				Default:    nil,
				EnvVar:     api.EnvVaultMFA,
				Completion: complete.PredictAnything,
				Usage:      "Supply MFA credentials as part of X-Vault-MFA header.",
			})

			f.BoolVar(&BoolVar{
				Name:    "output-curl-string",
				Target:  &c.flagOutputCurlString,
				Default: false,
				Usage: "Instead of executing the request, print an equivalent cURL " +
					"command string and exit.",
			})

			f.BoolVar(&BoolVar{
				Name:    "output-policy",
				Target:  &c.flagOutputPolicy,
				Default: false,
				Usage: "Instead of executing the request, print an example HCL " +
					"policy that would be required to run this command, and exit.",
			})

			f.StringVar(&StringVar{
				Name:       "unlock-key",
				Target:     &c.flagUnlockKey,
				Default:    notSetValue,
				Completion: complete.PredictNothing,
				Usage:      "Key to unlock a namespace API lock.",
			})

			f.StringMapVar(&StringMapVar{
				Name:       "header",
				Target:     &c.flagHeader,
				Completion: complete.PredictAnything,
				Usage: "Key-value pair provided as key=value to provide http header added to any request done by the CLI." +
					"Trying to add headers starting with 'X-Vault-' is forbidden and will make the command fail " +
					"This can be specified multiple times.",
			})

			f.BoolVar(&BoolVar{
				Name:    "non-interactive",
				Target:  &c.flagNonInteractive,
				Default: false,
				Usage:   "When set true, prevents asking the user for input via the terminal.",
			})

		}

		if bit&(FlagSetOutputField|FlagSetOutputFormat|FlagSetOutputDetailed) != 0 {
			outputSet := set.NewFlagSet("Output Options")

			if bit&FlagSetOutputField != 0 {
				outputSet.StringVar(&StringVar{
					Name:       "field",
					Target:     &c.flagField,
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
					Target:     &c.flagFormat,
					Default:    "table",
					EnvVar:     EnvVaultFormat,
					Completion: complete.PredictSet("table", "json", "yaml", "pretty"),
					Usage: `Print the output in the given format. Valid formats
						are "table", "json", "yaml", or "pretty".`,
				})
			}

			if bit&FlagSetOutputDetailed != 0 {
				outputSet.BoolVar(&BoolVar{
					Name:    "detailed",
					Target:  &c.flagDetailed,
					Default: false,
					EnvVar:  EnvVaultDetailed,
					Usage:   "Enables additional metadata during some operations",
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
}

// NewFlagSets creates a new flag sets.
func NewFlagSets(ui cli.Ui) *FlagSets {
	mainSet := flag.NewFlagSet("", flag.ContinueOnError)

	// Errors and usage are controlled by the CLI.
	mainSet.Usage = func() {}
	mainSet.SetOutput(ioutil.Discard)

	return &FlagSets{
		flagSets:    make([]*FlagSet, 0, 6),
		mainSet:     mainSet,
		hiddens:     make(map[string]struct{}),
		completions: complete.Flags{},
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

// Parse parses the given flags, returning any errors.
func (f *FlagSets) Parse(args []string) error {
	return f.mainSet.Parse(args)
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
func (fs *FlagSets) Help() string {
	var out bytes.Buffer

	for _, set := range fs.flagSets {
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
