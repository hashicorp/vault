package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// PolicyWriteCommand is a Command that enables a new endpoint.
type PolicyWriteCommand struct {
	meta.Meta
}

func (c *PolicyWriteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policy-write", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 2 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\npolicy-write expects exactly two arguments"))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	// Policies are normalized to lowercase
	name := strings.ToLower(args[0])
	path := args[1]

	// Read the policy
	var f io.Reader = os.Stdin
	if path != "-" {
		file, err := os.Open(path)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error opening file: %s", err))
			return 1
		}
		defer file.Close()
		f = file
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error reading file: %s", err))
		return 1
	}
	rules := buf.String()

	if err := client.Sys().PutPolicy(name, rules); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Policy '%s' written.", name))
	return 0
}

func (c *PolicyWriteCommand) Synopsis() string {
	return "Write a policy to the server"
}

func (c *PolicyWriteCommand) Help() string {
	helpText := `
Usage: vault policy-write [options] name path

  Write a policy with the given name from the contents of a file or stdin.

  If the path is "-", the policy is read from stdin. Otherwise, it is
  loaded from the file at the given path.

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
