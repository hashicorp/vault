package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

// LoginHandler is the interface that any auth handlers must implement to enable
// auth via the CLI.
type LoginHandler interface {
	Auth(*api.Client, map[string]string) (*api.Secret, error)
	Help() string
}

type LoginCommand struct {
	*BaseCommand

	Handlers map[string]LoginHandler

	flagMethod    string
	flagPath      string
	flagNoStore   bool
	flagNoPrint   bool
	flagTokenOnly bool

	testStdin io.Reader // for tests
}

func (c *LoginCommand) Synopsis() string {
	return "Authenticate locally"
}

func (c *LoginCommand) Help() string {
	helpText := `
Usage: vault login [options] [AUTH K=V...]

  Authenticates users or machines to Vault using the provided arguments. A
  successful authentication results in a Vault token - conceptually similar to
  a session token on a website. By default, this token is cached on the local
  machine for future requests.

  The default auth method is "token". If not supplied via the CLI,
  Vault will prompt for input. If the argument is "-", the values are read
  from stdin.

  The -method flag allows using other auth methods, such as userpass, github, or
  cert. For these, additional "K=V" pairs may be required. For example, to
  authenticate to the userpass auth method:

      $ vault login -method=userpass username=my-username

  For more information about the list of configuration parameters available for
  a given auth method, use the "vault auth help TYPE". You can also use "vault
  auth list" to see the list of enabled auth methods.

  If an auth method is enabled at a non-standard path, the -method flag still
  refers to the canonical type, but the -path flag refers to the enabled path.
  If a github auth method was enabled at "github-ent", authenticate like this:

      $ vault login -method=github -path=github-prod

  If the authentication is requested with response wrapping (via -wrap-ttl),
  the returned token is automatically unwrapped unless:

    - The -token-only flag is used, in which case this command will output
      the wrapping token.

    - The -no-store flag is used, in which case this command will output the
      details of the wrapping token.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *LoginCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "method",
		Target:     &c.flagMethod,
		Default:    "token",
		Completion: c.PredictVaultAvailableAuths(),
		Usage: "Type of authentication to use such as \"userpass\" or " +
			"\"ldap\". Note this corresponds to the TYPE, not the enabled path. " +
			"Use -path to specify the path where the authentication is enabled.",
	})

	f.StringVar(&StringVar{
		Name:       "path",
		Target:     &c.flagPath,
		Default:    "",
		Completion: c.PredictVaultAuths(),
		Usage: "Remote path in Vault where the auth method is enabled. " +
			"This defaults to the TYPE of method (e.g. userpass -> userpass/).",
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
		Name:    "no-print",
		Target:  &c.flagNoPrint,
		Default: false,
		Usage: "Do not display the token. The token will be still be stored to the " +
			"configured token helper.",
	})

	f.BoolVar(&BoolVar{
		Name:    "token-only",
		Target:  &c.flagTokenOnly,
		Default: false,
		Usage: "Output only the token with no verification. This flag is a " +
			"shortcut for \"-field=token -no-store\". Setting those flags to other " +
			"values will have no affect.",
	})

	return set
}

func (c *LoginCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *LoginCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *LoginCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()

	// Set the right flags if the user requested token-only - this overrides
	// any previously configured values, as documented.
	if c.flagTokenOnly {
		c.flagNoStore = true
		c.flagField = "token"
	}

	if c.flagNoStore && c.flagNoPrint {
		c.UI.Error(wrapAtLength(
			"-no-store and -no-print cannot be used together"))
		return 1
	}

	// Get the auth method
	authMethod := sanitizePath(c.flagMethod)
	if authMethod == "" {
		authMethod = "token"
	}

	// If no path is specified, we default the path to the method type
	// or use the plugin name if it's a plugin
	authPath := c.flagPath
	if authPath == "" {
		authPath = ensureTrailingSlash(authMethod)
	}

	// Get the handler function
	authHandler, ok := c.Handlers[authMethod]
	if !ok {
		c.UI.Error(wrapAtLength(fmt.Sprintf(
			"Unknown auth method: %s. Use \"vault auth list\" to see the "+
				"complete list of auth methods. Additionally, some "+
				"auth methods are only available via the HTTP API.",
			authMethod)))
		return 1
	}

	// Pull our fake stdin if needed
	stdin := (io.Reader)(os.Stdin)
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	// If the user provided a token, pass it along to the auth provider.
	if authMethod == "token" && len(args) > 0 && !strings.Contains(args[0], "=") {
		args = append([]string{"token=" + args[0]}, args[1:]...)
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
		c.UI.Error(fmt.Sprintf("Error authenticating: %s", err))
		return 2
	}

	// Unset any previous token wrapping functionality. If the original request
	// was for a wrapped token, we don't want future requests to be wrapped.
	client.SetWrappingLookupFunc(func(string, string) string { return "" })

	// Recursively extract the token, handling wrapping
	unwrap := !c.flagTokenOnly && !c.flagNoStore
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
		if c.flagTokenOnly {
			return PrintRawField(c.UI, secret, "wrapping_token")
		}
		if c.flagNoStore {
			return OutputSecret(c.UI, secret)
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

	if !c.flagNoStore {
		// Grab the token helper so we can store
		tokenHelper, err := c.TokenHelper()
		if err != nil {
			c.UI.Error(wrapAtLength(fmt.Sprintf(
				"Error initializing token helper. Please verify that the token "+
					"helper is available and properly configured for your system. The "+
					"error was: %s", err)))
			return 1
		}

		// Store the token in the local client
		if err := tokenHelper.Store(token); err != nil {
			c.UI.Error(fmt.Sprintf("Error storing token: %s", err))
			c.UI.Error(wrapAtLength(
				"Authentication was successful, but the token was not persisted. The "+
					"resulting token is shown below for your records.") + "\n")
			OutputSecret(c.UI, secret)
			return 2
		}

		// Warn if the VAULT_TOKEN environment variable is set, as that will take
		// precedence. We output as a warning, so piping should still work since it
		// will be on a different stream.
		if os.Getenv("VAULT_TOKEN") != "" {
			c.UI.Warn(wrapAtLength("WARNING! The VAULT_TOKEN environment variable "+
				"is set! This takes precedence over the value set by this command. To "+
				"use the value set by this command, unset the VAULT_TOKEN environment "+
				"variable or set it to the token displayed below.") + "\n")
		}
	} else if !c.flagTokenOnly {
		// If token-only the user knows it won't be stored, so don't warn
		c.UI.Warn(wrapAtLength(
			"The token was not stored in token helper. Set the VAULT_TOKEN "+
				"environment variable or pass the token below with each request to "+
				"Vault.") + "\n")
	}

	if c.flagNoPrint {
		return 0
	}

	// If the user requested a particular field, print that out now since we
	// are likely piping to another process.
	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	// Print some yay! text, but only in table mode.
	if Format(c.UI) == "table" {
		c.UI.Output(wrapAtLength(
			"Success! You are now authenticated. The token information displayed "+
				"below is already stored in the token helper. You do NOT need to run "+
				"\"vault login\" again. Future Vault requests will automatically use "+
				"this token.") + "\n")
	}

	return OutputSecret(c.UI, secret)
}

// extractToken extracts the token from the given secret, automatically
// unwrapping responses and handling error conditions if unwrap is true. The
// result also returns whether it was a wrapped response that was not unwrapped.
func (c *LoginCommand) extractToken(client *api.Client, secret *api.Secret, unwrap bool) (*api.Secret, bool, error) {
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
