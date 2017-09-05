package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

// AuthHandler is the interface that any auth handlers must implement
// to enable auth via the CLI.
type AuthHandler interface {
	Auth(*api.Client, map[string]string) (*api.Secret, error)
	Help() string
}

// AuthCommand is a Command that handles authentication.
type AuthCommand struct {
	*BaseCommand

	Handlers map[string]AuthHandler

	flagMethod    string
	flagPath      string
	flagNoVerify  bool
	flagNoStore   bool
	flagOnlyToken bool

	// Deprecations
	// TODO: remove in 0.9.0
	flagTokenOnly  bool
	flagMethods    bool
	flagMethodHelp bool

	testStdin io.Reader // for tests
}

func (c *AuthCommand) Synopsis() string {
	return "Authenticates users or machines"
}

func (c *AuthCommand) Help() string {
	helpText := `
Usage: vault auth [options] [AUTH K=V...]

  Authenticates users or machines to Vault using the provided arguments. By
  default, the authentication method is "token". If not supplied via the CLI,
  Vault will prompt for input. If argument is "-", the configuration options
  are read from stdin.

  The -method flag allows alternative authentication providers to be used,
  such as userpass, github, or cert. For these, additional "key=value" pairs
  may be required. For example, to authenticate to the userpass auth backend:

      $ vault auth -method=userpass username=my-username

  Use "vault auth-help TYPE" for more information about the list of
  configuration parameters and examples for a particular provider. Use the
  "vault auth-list" command to see a list of enabled authentication providers.

  If an authentication provider is mounted at a different path, the -method
  flag should by the canonical type, and the -path flag should be set to the
  mount path. If a github authentication provider was mounted at "github-ent",
  you would authenticate to this backend like this:

      $ vault auth -method=github -path=github-prod

  If the authentication is requested with response wrapping (via -wrap-ttl),
  the returned token is automatically unwrapped unless:

    - The -only-token flag is used, in which case this command will output
      the wrapping token

    - The -no-store flag is used, in which case this command will output
      the details of the wrapping token.

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuthCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "method",
		Target:     &c.flagMethod,
		Default:    "token",
		Completion: c.PredictVaultAvailableAuths(),
		Usage: "Type of authentication to use such as \"userpass\" or " +
			"\"ldap\". Note this corresponds to the TYPE, not the mount path. Use " +
			"-path to specify the path where the authentication is mounted.",
	})

	f.StringVar(&StringVar{
		Name:       "path",
		Target:     &c.flagPath,
		Default:    "",
		Completion: c.PredictVaultAuths(),
		Usage: "Mount point where the auth backend is enabled. This defaults to " +
			"the TYPE of backend (e.g. userpass -> userpass/).",
	})

	f.BoolVar(&BoolVar{
		Name:    "no-verify",
		Target:  &c.flagNoVerify,
		Default: false,
		Usage: "Do not verify the token after authentication. By default, Vault " +
			"issue a request to get more metdata about the token. This request " +
			"against the use-limit of the token. Set this to true to disable the " +
			"post-authenication lookup.",
	})

	f.BoolVar(&BoolVar{
		Name:    "no-store",
		Target:  &c.flagNoStore,
		Default: false,
		Usage: "Do not persist the token to the token helper (usually the " +
			"local filesystem) after authentication for use in future requests. " +
			"The token will only be displayed in the command output.",
	})

	f.BoolVar(&BoolVar{
		Name:    "only-token",
		Target:  &c.flagOnlyToken,
		Default: false,
		Usage: "Output only the token with no verification. This flag is a " +
			"shortcut for \"-field=token -no-store -no-verify\". Setting those " +
			"flags to other values will have no affect.",
	})

	// Deprecations
	// TODO: remove in Vault 0.9.0

	f.BoolVar(&BoolVar{
		Name:    "token-only", // Prefer only-token
		Target:  &c.flagTokenOnly,
		Default: false,
		Hidden:  true,
	})

	f.BoolVar(&BoolVar{
		Name:    "methods", // Prefer auth-list
		Target:  &c.flagMethods,
		Default: false,
		Hidden:  true,
	})

	f.BoolVar(&BoolVar{
		Name:    "method-help", // Prefer auth-help
		Target:  &c.flagMethodHelp,
		Default: false,
		Hidden:  true,
	})

	return set
}

func (c *AuthCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *AuthCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()

	// Deprecations - do this before any argument validations
	// TODO: remove in 0.9.0
	switch {
	case c.flagMethods:
		c.UI.Warn(wrapAtLength(
			"WARNING! The -methods flag is deprecated. Please use " +
				"\"vault auth-list\". This flag will be removed in the next major " +
				"release of Vault."))
		cmd := &AuthListCommand{
			BaseCommand: &BaseCommand{
				UI:     c.UI,
				client: c.client,
			},
		}
		return cmd.Run(nil)
	case c.flagMethodHelp:
		c.UI.Warn(wrapAtLength(
			"WARNING! The -method-help flag is deprecated. Please use " +
				"\"vault auth-help\". This flag will be removed in the next major " +
				"release of Vault."))
		cmd := &AuthHelpCommand{
			BaseCommand: &BaseCommand{
				UI:     c.UI,
				client: c.client,
			},
			Handlers: c.Handlers,
		}
		return cmd.Run([]string{c.flagMethod})
	}

	// TODO: remove in 0.9.0
	if c.flagTokenOnly {
		c.UI.Warn(wrapAtLength(
			"WARNING! The -token-only flag is deprecated. Plase use -only-token " +
				"instead. This flag will be removed in the next major release of " +
				"Vault."))
		c.flagOnlyToken = c.flagTokenOnly
	}

	// Set the right flags if the user requested only-token - this overrides
	// any previously configured values, as documented.
	if c.flagOnlyToken {
		c.flagNoStore = true
		c.flagNoVerify = true
		c.flagField = "token"
	}

	// Get the auth method
	authMethod := sanitizePath(c.flagMethod)
	if authMethod == "" {
		authMethod = "token"
	}

	// If no path is specified, we default the path to the backend type
	// or use the plugin name if it's a plugin backend
	authPath := c.flagPath
	if authPath == "" {
		authPath = ensureTrailingSlash(authMethod)
	}

	// Get the handler function
	authHandler, ok := c.Handlers[authMethod]
	if !ok {
		c.UI.Error(wrapAtLength(fmt.Sprintf(
			"Unknown authentication method: %s. Use \"vault auth-list\" to see the "+
				"complete list of authentication providers. Additionally, some "+
				"authentication providers are only available via the HTTP API.",
			authMethod)))
		return 1
	}

	// Pull our fake stdin if needed
	stdin := (io.Reader)(os.Stdin)
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	// If the user provided a token, pass it along to the auth provier.
	if authMethod == "token" && len(args) == 1 {
		args = []string{"token=" + args[0]}
	}

	config, err := parseArgsDataString(stdin, args)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error parsing configuration: %s", err))
		return 1
	}

	// If the user did not specify a mount path, use the provided mount path.
	if config["mount"] == "" && authPath != "" {
		config["mount"] = authPath
	}

	// Create the client
	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// Authenticate delegation to the auth handler
	secret, err := authHandler.Auth(client, config)
	if err != nil {
		c.UI.Error(wrapAtLength(fmt.Sprintf(
			"Error authenticating: %s", err)))
		return 2
	}

	// Unset any previous token wrapping functionality. If the original request
	// was for a wrapped token, we don't want future requests to be wrapped.
	client.SetWrappingLookupFunc(func(string, string) string { return "" })

	// Recursively extract the token, handling wrapping
	unwrap := !c.flagOnlyToken && !c.flagNoStore
	secret, isWrapped, err := c.extractToken(client, secret, unwrap)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error extracting token: %s", err))
		return 2
	}
	if secret == nil {
		c.UI.Error("Vault returned an empty secret")
		return 2
	}

	// Handle special cases if the token was wrapped
	if isWrapped {
		if c.flagOnlyToken {
			return PrintRawField(c.UI, secret, "wrapping_token")
		}
		if c.flagNoStore {
			return OutputSecret(c.UI, c.flagFormat, secret)
		}
	}

	// If we got this far, verify we have authentication data before continuing
	if secret.Auth == nil {
		c.UI.Error(wrapAtLength(
			"Vault returned a secret, but the secret has no authentication " +
				"information attached. This should never happen and is likely a " +
				"bug."))
		return 2
	}

	// Pull the token itself out, since we don't need the rest of the auth
	// information anymore/.
	token := secret.Auth.ClientToken

	tokenHelper, err := c.TokenHelper()
	if err != nil {
		c.UI.Error(fmt.Sprintf(
			"Error initializing token helper: %s\n\n"+
				"Please verify that the token helper is available and properly\n"+
				"configured for your system. Please refer to the documentation\n"+
				"on token helpers for more information.",
			err))
		return 1
	}

	if !c.flagNoVerify {
		// Verify the token and pull it's list of policies
		client.SetToken(token)
		client.SetWrappingLookupFunc(func(string, string) string { return "" })

		secret, err = client.Auth().Token().LookupSelf()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error verifying token: %s", err))
			return 2
		}
		if secret == nil {
			c.UI.Error("Empty response from lookup-self")
			return 2
		}
	}

	if !c.flagNoStore {
		// Store the token in the local client
		if err := tokenHelper.Store(token); err != nil {
			c.UI.Error(fmt.Sprintf("Error storing token: %s", err))
			c.UI.Error(wrapAtLength(
				"Authentication was successful, but the token was not persisted. The " +
					"resulting token is shown below for your records."))
			OutputSecret(c.UI, c.flagFormat, secret)
			return 2
		}

		// Warn if the VAULT_TOKEN environment variable is set, as that will take
		// precedence. Don't output on token-only since we're likely piping output.
		if c.flagField == "" && c.flagFormat == "table" {
			if os.Getenv("VAULT_TOKEN") != "" {
				c.UI.Warn(wrapAtLength("WARNING! The VAULT_TOKEN environment variable " +
					"is set! This takes precedence over the value set by this command. To " +
					"use the value set by this command, unset the VAULT_TOKEN environment " +
					"variable or set it to the token displayed below."))
			}
		}
	}

	// If the user requested a particular field, print that out now since we
	// are likely piping to another process.
	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	// Output the secret as json or yaml if requested. We have to maintain
	// backwards compatiability
	if c.flagFormat != "table" {
		return OutputSecret(c.UI, c.flagFormat, secret)
	}

	output := "Success! You are now authenticated. "
	if c.flagNoVerify {
		output += "The token was not verified for validity. "
	}
	if c.flagNoStore {
		output += "The token was not stored in the token helper. "
	} else {
		output += "The token information displayed below is already stored in " +
			"the token helper. You do NOT need to run \"vault auth\" again."
	}
	c.UI.Output(wrapAtLength(output) + "\n")

	// TODO make this consistent with other printed token secrets.
	c.UI.Output(fmt.Sprintf("token: %s", secret.TokenID()))
	c.UI.Output(fmt.Sprintf("accessor: %s", secret.TokenAccessor()))

	if ttl := secret.TokenTTL(); ttl != 0 {
		c.UI.Output(fmt.Sprintf("duration: %s", ttl))
	}

	c.UI.Output(fmt.Sprintf("renewable: %t", secret.TokenIsRenewable()))

	if policies := secret.TokenPolicies(); len(policies) > 0 {
		c.UI.Output(fmt.Sprintf("policies: %s", policies))
	}

	return 0
}

// extractToken extracts the token from the given secret, automatically
// unwrapping responses and handling error conditions if unwrap is true. The
// result also returns whether it was a wrapped resonse that was not unwrapped.
func (c *AuthCommand) extractToken(client *api.Client, secret *api.Secret, unwrap bool) (*api.Secret, bool, error) {
	switch {
	case secret == nil:
		return nil, false, fmt.Errorf("empty response from auth helper")

	case secret.Auth != nil:
		return secret, false, nil

	case secret.WrapInfo != nil:
		if secret.WrapInfo.WrappedAccessor == "" {
			return nil, false, fmt.Errorf("wrapped response does not contain a token")
		}

		if !unwrap {
			return secret, true, nil
		}

		client.SetToken(secret.WrapInfo.Token)
		secret, err := client.Logical().Unwrap("")
		if err != nil {
			return nil, false, err
		}
		return c.extractToken(client, secret, unwrap)

	default:
		return nil, false, fmt.Errorf("no auth or wrapping info in response")
	}
}
