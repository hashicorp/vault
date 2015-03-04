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

	return 0
}

func (c *AuthCommand) Synopsis() string {
	return "Prints information about how to authenticate with Vault"
}

func (c *AuthCommand) Help() string {
	helpText := `
Usage: vault auth [options]

  Outputs instructions for authenticating with vault.

  Vault authentication is always done via environmental variables. The
  specific environmental variables and the meaning for each environmental
  variable varies depending on the auth mechanism that Vault is using.
  This command outputs the mechanism vault is using along with documentation
  for how to authenticate.

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
