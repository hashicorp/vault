package command

import (
	"strings"
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
