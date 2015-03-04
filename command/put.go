package command

import (
	"strings"
)

// PutCommand is a Command that puts data into the Vault.
type PutCommand struct {
	Meta
}

func (c *PutCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("put", FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	return 0
}

func (c *PutCommand) Synopsis() string {
	return "Put secrets or configuration into Vault"
}

func (c *PutCommand) Help() string {
	helpText := `
Usage: vault put [options] path data

  Write data (secrets or configuration) into Vault.

  Put sends data into Vault at the given path. The behavior of the write
  is determined by the backend at the given path. For example, writing
  to "aws/policy/ops" will create an "ops" IAM policy for the AWS backend
  (configuration), but writing to "consul/foo" will write a value directly
  into Consul at that key. Check the documentation of the logical backend
  you're using for more information on key structure.

  If data is "-" then the data will be ready from stdin. To write a literal
  "-", you'll have to pipe that value in from stdin.

General Options:

  -address=TODO           The address of the Vault server.

  -ca-cert=path           Path to a PEM encoded CA cert file to use to
                          verify the Vault server SSL certificate.

  -ca-path=path           Path to a directory of PEM encoded CA cert files
                          to verify the Vault server SSL certificate. If both
                          -ca-cert and -ca-path are specified, -ca-path is used.

  -insecure               Do not verify TLS certificate. This is highly
                          not recommended. This is especially not recommended
                          for unsealing a vault.

`
	return strings.TrimSpace(helpText)
}
