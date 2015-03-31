package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/helper/password"
)

// AuthCommand is a Command that handles authentication.
type AuthCommand struct {
	Meta
}

func (c *AuthCommand) Run(args []string) int {
	var method string
	flags := c.Meta.FlagSet("auth", FlagSetDefault)
	flags.StringVar(&method, "method", "", "method")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) > 1 {
		flags.Usage()
		c.Ui.Error("\nError: auth expects at most one argument")
		return 1
	}
	if method != "" && len(args) > 0 {
		flags.Usage()
		c.Ui.Error("\nError: auth expects no arguments if -method is specified")
		return 1
	}

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
	var token string
	if method == "" {
		if len(args) > 0 {
			token = args[0]

			// TODO(mitchellh): stdin
		} else {
			// No arguments given, read the token from user input
			fmt.Printf("Token (will be hidden): ")
			token, err = password.Read(os.Stdin)
			if err != nil {
				c.Ui.Error(fmt.Sprintf(
					"Error attempting to ask for token. The raw error message\n"+
						"is shown below, but the most common reason for this error is\n"+
						"that you attempted to pipe a value into auth. If you want to\n"+
						"pipe the token, please pass '-' as the token argument.\n\n"+
						"Raw error: %s", err))
				return 1
			}
		}

		if token == "" {
			c.Ui.Error(fmt.Sprintf(
				"A token must be passed to auth. Please view the help\n" +
					"for more information."))
			return 1
		}
	} else {
		// TODO(mitchellh): other auth methods
	}

	// Store the token!
	if err := tokenHelper.Store(token); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error storing token: %s\n\n"+
				"Authentication was not successful and did not persist.\n"+
				"Please reauthenticate, or fix the issue above if possible.",
			err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf(
		"Successfully authenticated!"))

	return 0
}

func (c *AuthCommand) Synopsis() string {
	return "Prints information about how to authenticate with Vault"
}

func (c *AuthCommand) Help() string {
	helpText := `
Usage: vault auth [options] [token]

  Authenticate with Vault with the given token or via any supported
  authentication backend.

  If no -method is specified, then the token is expected. If it is not
  given on the command-line, it will be asked via user input. If the
  token is "-", it will be read from stdin.

  By specifying -method, alternate authentication methods can be used
  such as OAuth or TLS certificates. For these, additional -var flags
  may be required. It is an error to specify a token with -method.

General Options:

  -address=TODO           The address of the Vault server.

  -ca-cert=path           Path to a PEM encoded CA cert file to use to
                          verify the Vault server SSL certificate.

  -ca-path=path           Path to a directory of PEM encoded CA cert files
                          to verify the Vault server SSL certificate. If both
                          -ca-cert and -ca-path are specified, -ca-path is used.

  -insecure               Do not verify TLS certificate. This is highly
                          not recommended.

Auth Options:

  -method=name    Outputs help for the authentication method with the given
                  name for the remote server. If this authentication method
                  is not available, exit with code 1.
`
	return strings.TrimSpace(helpText)
}
