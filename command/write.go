package command

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// DefaultDataKey is the key used in the write as a default for data.
const DefaultDataKey = "value"

// WriteCommand is a Command that puts data into the Vault.
type WriteCommand struct {
	Meta

	// The fields below can be overwritten for tests
	testStdin io.Reader
}

func (c *WriteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("write", FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 2 {
		c.Ui.Error("write expects two arguments")
		flags.Usage()
		return 1
	}

	path := args[0]
	if path[0] == '/' {
		path = path[1:]
	}
	var data map[string]interface{}
	if args[1] == "-" {
		var stdin io.Reader = os.Stdin
		if c.testStdin != nil {
			stdin = c.testStdin
		}

		dec := json.NewDecoder(stdin)
		if err := dec.Decode(&data); err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error decoding JSON of stdin: %s", err))
			return 1
		}
	} else {
		data = map[string]interface{}{DefaultDataKey: args[1]}
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	if err := client.Logical().Write(path, data); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error writing data to %s: %s", path, err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Success! Data written to: %s", path))
	return 0
}

func (c *WriteCommand) Synopsis() string {
	return "Write secrets or configuration into Vault"
}

func (c *WriteCommand) Help() string {
	helpText := `
Usage: vault write [options] path data

  Write data (secrets or configuration) into Vault.

  Write sends data into Vault at the given path. The behavior of the write
  is determined by the backend at the given path. For example, writing
  to "aws/policy/ops" will create an "ops" IAM policy for the AWS backend
  (configuration), but writing to "consul/foo" will write a value directly
  into Consul at that key. Check the documentation of the logical backend
  you're using for more information on key structure.

  If data is "-" then the data will be ready from stdin. To write a literal
  "-", you'll have to pipe that value in from stdin. To write data from a
  file, pipe the file contents in via stdin and set data to "-".

  If data is a string, it will be sent with the key of "value".

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
