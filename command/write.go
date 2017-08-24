package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/helper/kv-builder"
	"github.com/hashicorp/vault/meta"
	"github.com/posener/complete"
)

// WriteCommand is a Command that puts data into the Vault.
type WriteCommand struct {
	meta.Meta

	// The fields below can be overwritten for tests
	testStdin io.Reader
}

func (c *WriteCommand) Run(args []string) int {
	var field, format string
	var force bool
	flags := c.Meta.FlagSet("write", meta.FlagSetDefault)
	flags.StringVar(&format, "format", "table", "")
	flags.StringVar(&field, "field", "", "")
	flags.BoolVar(&force, "force", false, "")
	flags.BoolVar(&force, "f", false, "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) < 1 {
		c.Ui.Error("write requires a path")
		flags.Usage()
		return 1
	}

	if len(args) < 2 && !force {
		c.Ui.Error("write expects at least two arguments; use -f to perform the write anyways")
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

	secret, err := client.Logical().Write(path, data)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error writing data to %s: %s", path, err))
		return 1
	}

	if secret == nil {
		// Don't output anything if people aren't using the "human" output
		if format == "table" {
			c.Ui.Output(fmt.Sprintf("Success! Data written to: %s", path))
		}
		return 0
	}

	// Handle single field output
	if field != "" {
		return PrintRawField(c.Ui, secret, field)
	}

	return OutputSecret(c.Ui, format, secret)
}

func (c *WriteCommand) parseData(args []string) (map[string]interface{}, error) {
	var stdin io.Reader = os.Stdin
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	builder := &kvbuilder.Builder{Stdin: stdin}
	if err := builder.Add(args...); err != nil {
		return nil, err
	}

	return builder.Map(), nil
}

func (c *WriteCommand) Synopsis() string {
	return "Write secrets or configuration into Vault"
}

func (c *WriteCommand) Help() string {
	helpText := `
Usage: vault write [options] path [data]

  Write data (secrets or configuration) into Vault.

  Write sends data into Vault at the given path. The behavior of the write is
  determined by the backend at the given path. For example, writing to
  "aws/policy/ops" will create an "ops" IAM policy for the AWS backend
  (configuration), but writing to "consul/foo" will write a value directly into
  Consul at that key. Check the documentation of the logical backend you're
  using for more information on key structure.

  Data is sent via additional arguments in "key=value" pairs. If value begins
  with an "@", then it is loaded from a file. Write expects data in the file to
  be in JSON format. If you want to start the value with a literal "@", then
  prefix the "@" with a slash: "\@".

General Options:
` + meta.GeneralOptionsUsage() + `
Write Options:

  -f | -force             Force the write to continue without any data values
                          specified. This allows writing to keys that do not
                          need or expect any fields to be specified.

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json or yaml.

  -field=field            If included, the raw value of the specified field
                          will be output raw to stdout.

`
	return strings.TrimSpace(helpText)
}

func (c *WriteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *WriteCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{
		"-force":  complete.PredictNothing,
		"-format": predictFormat,
		"-field":  complete.PredictNothing,
	}
}
