package command

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/kv-builder"
	"github.com/hashicorp/vault/helper/password"
	"github.com/hashicorp/vault/meta"
	"github.com/mitchellh/mapstructure"
	"github.com/posener/complete"
	"github.com/ryanuber/columnize"
)

// AuthHandler is the interface that any auth handlers must implement
// to enable auth via the CLI.
type AuthHandler interface {
	Auth(*api.Client, map[string]string) (*api.Secret, error)
	Help() string
}

// AuthCommand is a Command that handles authentication.
type AuthCommand struct {
	meta.Meta

	Handlers map[string]AuthHandler

	// The fields below can be overwritten for tests
	testStdin io.Reader
}

func (c *AuthCommand) Run(args []string) int {
	var method, authPath string
	var methods, methodHelp, noVerify, noStore, tokenOnly bool
	flags := c.Meta.FlagSet("auth", meta.FlagSetDefault)
	flags.BoolVar(&methods, "methods", false, "")
	flags.BoolVar(&methodHelp, "method-help", false, "")
	flags.BoolVar(&noVerify, "no-verify", false, "")
	flags.BoolVar(&noStore, "no-store", false, "")
	flags.BoolVar(&tokenOnly, "token-only", false, "")
	flags.StringVar(&method, "method", "", "method")
	flags.StringVar(&authPath, "path", "", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	if methods {
		return c.listMethods()
	}

	args = flags.Args()

	tokenHelper, err := c.TokenHelper()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing token helper: %s\n\n"+
				"Please verify that the token helper is available and properly\n"+
				"configured for your system. Please refer to the documentation\n"+
				"on token helpers for more information.",
			err))
		return 1
	}

	// token is where the final token will go
	handler := c.Handlers[method]

	// Read token from stdin if first arg is exactly "-"
	var stdin io.Reader = os.Stdin
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	if len(args) > 0 && args[0] == "-" {
		stdinR := bufio.NewReader(stdin)
		args[0], err = stdinR.ReadString('\n')
		if err != nil && err != io.EOF {
			c.Ui.Error(fmt.Sprintf("Error reading from stdin: %s", err))
			return 1
		}
		args[0] = strings.TrimSpace(args[0])
	}

	if method == "" {
		token := ""
		if len(args) > 0 {
			token = args[0]
		}

		handler = &tokenAuthHandler{Token: token}
		args = nil

		switch authPath {
		case "", "auth/token":
		default:
			c.Ui.Error("Token authentication does not support custom paths")
			return 1
		}
	}

	if handler == nil {
		methods := make([]string, 0, len(c.Handlers))
		for k := range c.Handlers {
			methods = append(methods, k)
		}
		sort.Strings(methods)

		c.Ui.Error(fmt.Sprintf(
			"Unknown authentication method: %s\n\n"+
				"Please use a supported authentication method. The list of supported\n"+
				"authentication methods is shown below. Note that this list may not\n"+
				"be exhaustive: Vault may support other auth methods. For auth methods\n"+
				"unsupported by the CLI, please use the HTTP API.\n\n"+
				"%s",
			method,
			strings.Join(methods, ", ")))
		return 1
	}

	if methodHelp {
		c.Ui.Output(handler.Help())
		return 0
	}

	// Warn if the VAULT_TOKEN environment variable is set, as that will take
	// precedence. Don't output on token-only since we're likely piping output.
	if os.Getenv("VAULT_TOKEN") != "" && !tokenOnly {
		c.Ui.Output("==> WARNING: VAULT_TOKEN environment variable set!\n")
		c.Ui.Output("  The environment variable takes precedence over the value")
		c.Ui.Output("  set by the auth command. Either update the value of the")
		c.Ui.Output("  environment variable or unset it to use the new token.\n")
	}

	var vars map[string]string
	if len(args) > 0 {
		builder := kvbuilder.Builder{Stdin: os.Stdin}
		if err := builder.Add(args...); err != nil {
			c.Ui.Error(err.Error())
			return 1
		}

		if err := mapstructure.Decode(builder.Map(), &vars); err != nil {
			c.Ui.Error(fmt.Sprintf("Error parsing options: %s", err))
			return 1
		}
	} else {
		vars = make(map[string]string)
	}

	// Build the client so we can auth
	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client to auth: %s", err))
		return 1
	}

	if authPath != "" {
		vars["mount"] = authPath
	}

	// Authenticate
	secret, err := handler.Auth(client, vars)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}
	if secret == nil {
		c.Ui.Error("Empty response from auth helper")
		return 1
	}

	// If we had requested a wrapped token, we want to unset that request
	// before performing further functions
	client.SetWrappingLookupFunc(func(string, string) string {
		return ""
	})

CHECK_TOKEN:
	var token string
	switch {
	case secret == nil:
		c.Ui.Error("Empty response from auth helper")
		return 1

	case secret.Auth != nil:
		token = secret.Auth.ClientToken

	case secret.WrapInfo != nil:
		if secret.WrapInfo.WrappedAccessor == "" {
			c.Ui.Error("Got a wrapped response from Vault but wrapped reply does not seem to contain a token")
			return 1
		}
		if tokenOnly {
			c.Ui.Output(secret.WrapInfo.Token)
			return 0
		}
		if noStore {
			return OutputSecret(c.Ui, "table", secret)
		}
		client.SetToken(secret.WrapInfo.Token)
		secret, err = client.Logical().Unwrap("")
		goto CHECK_TOKEN

	default:
		c.Ui.Error("No auth or wrapping info in auth helper response")
		return 1
	}

	// Cache the previous token so that it can be restored if authentication fails
	var previousToken string
	if previousToken, err = tokenHelper.Get(); err != nil {
		c.Ui.Error(fmt.Sprintf("Error caching the previous token: %s\n\n", err))
		return 1
	}

	if tokenOnly {
		c.Ui.Output(token)
		return 0
	}

	// Store the token!
	if !noStore {
		if err := tokenHelper.Store(token); err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error storing token: %s\n\n"+
					"Authentication was not successful and did not persist.\n"+
					"Please reauthenticate, or fix the issue above if possible.",
				err))
			return 1
		}
	}

	if noVerify {
		c.Ui.Output(fmt.Sprintf(
			"Authenticated - no token verification has been performed.",
		))

		if noStore {
			if err := tokenHelper.Erase(); err != nil {
				c.Ui.Error(fmt.Sprintf(
					"Error removing prior token: %s\n\n"+
						"Authentication was successful, but unable to remove the\n"+
						"previous token.",
					err))
				return 1
			}
		}
		return 0
	}

	// Build the client again so it can read the token we just wrote
	client, err = c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client to verify the token: %s", err))
		if !noStore {
			if err := tokenHelper.Store(previousToken); err != nil {
				c.Ui.Error(fmt.Sprintf(
					"Error restoring the previous token: %s\n\n"+
						"Please reauthenticate with a valid token.",
					err))
			}
		}
		return 1
	}
	client.SetWrappingLookupFunc(func(string, string) string {
		return ""
	})

	// If in no-store mode it won't have read the token from a token-helper (or
	// will read an old one) so set it explicitly
	if noStore {
		client.SetToken(token)
	}

	// Verify the token
	secret, err = client.Auth().Token().LookupSelf()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error validating token: %s", err))
		if err := tokenHelper.Store(previousToken); err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error restoring the previous token: %s\n\n"+
					"Please reauthenticate with a valid token.",
				err))
		}
		return 1
	}
	if secret == nil && !noStore {
		c.Ui.Error(fmt.Sprintf("Error: Invalid token"))
		if err := tokenHelper.Store(previousToken); err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error restoring the previous token: %s\n\n"+
					"Please reauthenticate with a valid token.",
				err))
		}
		return 1
	}

	if noStore {
		if err := tokenHelper.Erase(); err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error removing prior token: %s\n\n"+
					"Authentication was successful, but unable to remove the\n"+
					"previous token.",
				err))
			return 1
		}
	}

	// Get the policies we have
	policiesRaw, ok := secret.Data["policies"]
	if !ok || policiesRaw == nil {
		policiesRaw = []interface{}{"unknown"}
	}
	var policies []string
	for _, v := range policiesRaw.([]interface{}) {
		policies = append(policies, v.(string))
	}

	output := "Successfully authenticated! You are now logged in."
	if noStore {
		output += "\nThe token has not been stored to the configured token helper."
	}
	if method != "" {
		output += "\nThe token below is already saved in the session. You do not"
		output += "\nneed to \"vault auth\" again with the token."
	}
	output += fmt.Sprintf("\ntoken: %s", secret.Data["id"])
	output += fmt.Sprintf("\ntoken_duration: %s", secret.Data["ttl"].(json.Number).String())
	if len(policies) > 0 {
		output += fmt.Sprintf("\ntoken_policies: %v", policies)
	}

	c.Ui.Output(output)

	return 0

}

func (c *AuthCommand) getMethods() (map[string]*api.AuthMount, error) {
	client, err := c.Client()
	if err != nil {
		return nil, err
	}
	client.SetWrappingLookupFunc(func(string, string) string {
		return ""
	})

	auth, err := client.Sys().ListAuth()
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func (c *AuthCommand) listMethods() int {
	auth, err := c.getMethods()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error reading auth table: %s", err))
		return 1
	}

	paths := make([]string, 0, len(auth))
	for path := range auth {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	columns := []string{"Path | Type | Accessor | Default TTL | Max TTL | Replication Behavior | Description"}
	for _, path := range paths {
		auth := auth[path]
		defTTL := "system"
		if auth.Config.DefaultLeaseTTL != 0 {
			defTTL = strconv.Itoa(auth.Config.DefaultLeaseTTL)
		}
		maxTTL := "system"
		if auth.Config.MaxLeaseTTL != 0 {
			maxTTL = strconv.Itoa(auth.Config.MaxLeaseTTL)
		}
		replicatedBehavior := "replicated"
		if auth.Local {
			replicatedBehavior = "local"
		}
		columns = append(columns, fmt.Sprintf(
			"%s | %s | %s | %s | %s | %s | %s", path, auth.Type, auth.Accessor, defTTL, maxTTL, replicatedBehavior, auth.Description))
	}

	c.Ui.Output(columnize.SimpleFormat(columns))
	return 0
}

func (c *AuthCommand) Synopsis() string {
	return "Prints information about how to authenticate with Vault"
}

func (c *AuthCommand) Help() string {
	helpText := `
Usage: vault auth [options] [auth-information]

  Authenticate with Vault using the given token or via any supported
  authentication backend.

  By default, the -method is assumed to be token. If not supplied via the
  command-line, a prompt for input will be shown. If the authentication
  information is "-", it will be read from stdin.

  The -method option allows alternative authentication methods to be used,
  such as userpass, GitHub, or TLS certificates. For these, additional
  values as "key=value" pairs may be required. For example, to authenticate
  to the userpass auth backend:

      $ vault auth -method=userpass username=my-username

  Use "-method-help" to get help for a specific method.

  If an auth backend is enabled at a different path, the "-method" flag
  should still point to the canonical name, and the "-path" flag should be
  used. If a GitHub auth backend was mounted as "github-private", one would
  authenticate to this backend via:

      $ vault auth -method=github -path=github-private

  The value of the "-path" flag is supplied to auth providers as the "mount"
  option in the payload to specify the mount point.

  If response wrapping is used (via -wrap-ttl), the returned token will be
  automatically unwrapped unless:
    * -token-only is used, in which case the wrapping token will be output
    * -no-store is used, in which case the details of the wrapping token
      will be printed

General Options:

  ` + meta.GeneralOptionsUsage() + `

Auth Options:

  -method=name      Use the method given here, which is a type of backend, not
                    the path. If this authentication method is not available,
                    exit with code 1.

  -method-help      If set, the help for the selected method will be shown.

  -methods          List the available auth methods.

  -no-verify        Do not verify the token after creation; avoids a use count
                    decrement.

  -no-store         Do not store the token after creation; it will only be
                    displayed in the command output.

  -token-only       Output only the token to stdout. This implies -no-verify
                    and -no-store.

  -path             The path at which the auth backend is enabled. If an auth
                    backend is mounted at multiple paths, this option can be
                    used to authenticate against specific paths.
`
	return strings.TrimSpace(helpText)
}

// tokenAuthHandler handles retrieving the token from the command-line.
type tokenAuthHandler struct {
	Token string
}

func (h *tokenAuthHandler) Auth(*api.Client, map[string]string) (*api.Secret, error) {
	token := h.Token
	if token == "" {
		var err error

		// No arguments given, read the token from user input
		fmt.Printf("Token (will be hidden): ")
		token, err = password.Read(os.Stdin)
		fmt.Printf("\n")
		if err != nil {
			return nil, fmt.Errorf(
				"Error attempting to ask for token. The raw error message\n"+
					"is shown below, but the most common reason for this error is\n"+
					"that you attempted to pipe a value into auth. If you want to\n"+
					"pipe the token, please pass '-' as the token argument.\n\n"+
					"Raw error: %s", err)
		}
	}

	if token == "" {
		return nil, fmt.Errorf(
			"A token must be passed to auth. Please view the help\n" +
				"for more information.")
	}

	return &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken: token,
		},
	}, nil
}

func (h *tokenAuthHandler) Help() string {
	help := `
No method selected with the "-method" flag, so the "auth" command assumes
you'll be using raw token authentication. For this, specify the token to
authenticate as the parameter to "vault auth". Example:

    vault auth 123456

The token used to authenticate must come from some other source. A root
token is created when Vault is first initialized. After that, subsequent
tokens are created via the API or command line interface (with the
"token"-prefixed commands).
`

	return strings.TrimSpace(help)
}

func (c *AuthCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *AuthCommand) AutocompleteFlags() complete.Flags {
	var predictFunc complete.PredictFunc = func(a complete.Args) []string {
		auths, err := c.getMethods()
		if err != nil {
			return []string{}
		}

		methods := make([]string, 0, len(auths))
		for _, auth := range auths {
			if strings.HasPrefix(auth.Type, a.Last) {
				methods = append(methods, auth.Type)
			}
		}

		return methods
	}

	return complete.Flags{
		"-method":      predictFunc,
		"-methods":     complete.PredictNothing,
		"-method-help": complete.PredictNothing,
		"-no-verify":   complete.PredictNothing,
		"-no-store":    complete.PredictNothing,
		"-token-only":  complete.PredictNothing,
		"-path":        complete.PredictNothing,
	}
}
