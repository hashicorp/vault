package token

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/password"
)

type CLIHandler struct {
	// for tests
	testStdin  io.Reader
	testStdout io.Writer
}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	// Parse "lookup" first - we want to return an early error if the user
	// supplied an invalid value here before we prompt them for a token. It would
	// be annoying to type your token and then be told you supplied an invalid
	// value that we could have known in advance.
	lookup := true
	if x, ok := m["lookup"]; ok {
		parsed, err := strconv.ParseBool(x)
		if err != nil {
			return nil, errwrap.Wrapf("Failed to parse \"lookup\" as boolean: {{err}}", err)
		}
		lookup = parsed
	}

	// Parse the token.
	token, ok := m["token"]
	if !ok {
		// Override the output
		stdout := h.testStdout
		if stdout == nil {
			stdout = os.Stderr
		}

		// No arguments given, read the token from user input
		fmt.Fprintf(stdout, "Token (will be hidden): ")
		var err error
		token, err = password.Read(os.Stdin)
		fmt.Fprintf(stdout, "\n")

		if err != nil {
			if err == password.ErrInterrupted {
				return nil, fmt.Errorf("user interrupted")
			}

			return nil, errwrap.Wrapf("An error occurred attempting to "+
				"ask for a token. The raw error message is shown below, but usually "+
				"this is because you attempted to pipe a value into the command or "+
				"you are executing outside of a terminal (tty). If you want to pipe "+
				"the value, pass \"-\" as the argument to read from stdin. The raw "+
				"error was: {{err}}", err)
		}
	}

	// Remove any whitespace, etc.
	token = strings.TrimSpace(token)

	if token == "" {
		return nil, fmt.Errorf(
			"a token must be passed to auth, please view the help for more " +
				"information")
	}

	// If the user declined verification, return now. Note that we will not have
	// a lot of information about the token.
	if !lookup {
		return &api.Secret{
			Auth: &api.SecretAuth{
				ClientToken: token,
			},
		}, nil
	}

	// If we got this far, we want to lookup and lookup the token and pull it's
	// list of policies an metadata.
	c.SetToken(token)
	c.SetWrappingLookupFunc(func(string, string) string { return "" })

	secret, err := c.Auth().Token().LookupSelf()
	if err != nil {
		return nil, errwrap.Wrapf("error looking up token: {{err}}", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("empty response from lookup-self")
	}

	// Return an auth struct that "looks" like the response from an auth method.
	// lookup and lookup-self return their data in data, not auth. We try to
	// mirror that data here.
	id, err := secret.TokenID()
	if err != nil {
		return nil, errwrap.Wrapf("error accessing token ID: {{err}}", err)
	}
	accessor, err := secret.TokenAccessor()
	if err != nil {
		return nil, errwrap.Wrapf("error accessing token accessor: {{err}}", err)
	}
	// This populates secret.Auth
	_, err = secret.TokenPolicies()
	if err != nil {
		return nil, errwrap.Wrapf("error accessing token policies: {{err}}", err)
	}
	metadata, err := secret.TokenMetadata()
	if err != nil {
		return nil, errwrap.Wrapf("error accessing token metadata: {{err}}", err)
	}
	dur, err := secret.TokenTTL()
	if err != nil {
		return nil, errwrap.Wrapf("error converting token TTL: {{err}}", err)
	}
	renewable, err := secret.TokenIsRenewable()
	if err != nil {
		return nil, errwrap.Wrapf("error checking if token is renewable: {{err}}", err)
	}
	return &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken:      id,
			Accessor:         accessor,
			Policies:         secret.Auth.Policies,
			TokenPolicies:    secret.Auth.TokenPolicies,
			IdentityPolicies: secret.Auth.IdentityPolicies,
			Metadata:         metadata,

			LeaseDuration: int(dur.Seconds()),
			Renewable:     renewable,
		},
	}, nil

}

func (h *CLIHandler) Help() string {
	help := `
Usage: vault login TOKEN [CONFIG K=V...]

  The token auth method allows logging in directly with a token. This
  can be a token from the "token-create" command or API. There are no
  configuration options for this auth method.

  Authenticate using a token:

      $ vault login 96ddf4bc-d217-f3ba-f9bd-017055595017

  Authenticate but do not lookup information about the token:

      $ vault login token=96ddf4bc-d217-f3ba-f9bd-017055595017 lookup=false

  This token usually comes from a different source such as the API or via the
  built-in "vault token create" command.

Configuration:

  token=<string>
      The token to use for authentication. This is usually provided directly
      via the "vault login" command.

  lookup=<bool>
      Perform a lookup of the token's metadata and policies.
`

	return strings.TrimSpace(help)
}
