package token

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/password"
)

type CLIHandler struct {
	// for tests
	testStdin  io.Reader
	testStdout io.Writer
}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	token, ok := m["token"]
	if !ok {
		// Override the output
		stdout := h.testStdout
		if stdout == nil {
			stdout = os.Stdout
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

			return nil, fmt.Errorf("An error occurred attempting to "+
				"ask for a token. The raw error message is shown below, but usually "+
				"this is because you attempted to pipe a value into the command or "+
				"you are executing outside of a terminal (tty). If you want to pipe "+
				"the value, pass \"-\" as the argument to read from stdin. The raw "+
				"error was: %s", err)
		}
	}

	// Remove any whitespace, etc.
	token = strings.TrimSpace(token)

	if token == "" {
		return nil, fmt.Errorf(
			"A token must be passed to auth. Please view the help for more " +
				"information.")
	}

	return &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken: token,
		},
	}, nil
}

func (h *CLIHandler) Help() string {
	help := `
Usage: vault auth TOKEN [CONFIG K=V...]

  The token authentication provider allows logging in directly with a token.
  This can be a token from the "token-create" command or API. There are no
  configuration options for this authentication provider.

  Authenticate using a token:

      $ vault auth 96ddf4bc-d217-f3ba-f9bd-017055595017

  This token usually comes from a different source such as the API or via the
  built-in "vault token-create" command.

Configuration:

  token=<string>
      The token to use for authentication. This is usually provided directly
      via the "vault auth" command.

`

	return strings.TrimSpace(help)
}
