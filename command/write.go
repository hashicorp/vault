package command

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

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
	if len(args) < 2 {
		c.Ui.Error("write expects at least two arguments")
		flags.Usage()
		return 1
	}

	path := args[0]
	if path[0] == '/' {
		path = path[1:]
	}

	data, err := c.parseData(args[1:])
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error loading data: %s", err))
		return 1
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

func (c *WriteCommand) parseData(args []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for i, arg := range args {
		// If the arg is exactly "-" then we read from stdin and merge
		// the resulting structure into the result.
		if arg == "-" {
			var stdin io.Reader = os.Stdin
			if c.testStdin != nil {
				stdin = c.testStdin
			}

			dec := json.NewDecoder(stdin)
			if err := dec.Decode(&result); err != nil {
				return nil, fmt.Errorf(
					"Error loading data at index %d: %s", i, err)
			}

			continue
		}

		// If the arg begins with "@" then we read the file directly.
		if arg[0] == '@' {
			f, err := os.Open(arg[1:])
			if err != nil {
				return nil, fmt.Errorf(
					"Error loading data at index %d: %s", i, err)
			}

			dec := json.NewDecoder(f)
			err = dec.Decode(&result)
			f.Close()
			if err != nil {
				return nil, fmt.Errorf(
					"Error loading data at index %d: %s", i, err)
			}

			continue
		}

		// Split into key/value
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf(
				"Data at index %d is not in key=value format: %s",
				i, arg)
		}
		key, value := parts[0], parts[1]

		if value[0] == '@' {
			contents, err := ioutil.ReadFile(value[1:])
			if err != nil {
				return nil, fmt.Errorf(
					"Error reading file value for index %d: %s", i, err)
			}

			value = string(contents)
		} else if value[0] == '\\' && value[1] == '@' {
			value = value[1:]
		}

		result[key] = value
	}

	return result, nil
}

func (c *WriteCommand) Synopsis() string {
	return "Write secrets or configuration into Vault"
}

func (c *WriteCommand) Help() string {
	helpText := `
Usage: vault write [options] path [data]

  Write data (secrets or configuration) into Vault.

  Write sends data into Vault at the given path. The behavior of the write
  is determined by the backend at the given path. For example, writing
  to "aws/policy/ops" will create an "ops" IAM policy for the AWS backend
  (configuration), but writing to "consul/foo" will write a value directly
  into Consul at that key. Check the documentation of the logical backend
  you're using for more information on key structure.

  Data is sent via additional arguments in "key=value" pairs. If value
  begins with an "@", then it is loaded from a file. If you want to start
  the value with a literal "@", then prefix the "@" with a slash: "\@".

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
